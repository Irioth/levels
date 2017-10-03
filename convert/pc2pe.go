package convert

import (
	"crystal/mc"
	"crystal/mc/world"
	"encoding/hex"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/crystalmine/levels/anvil"
	"github.com/crystalmine/levels/bedrock"
)

func PC2PE(pcpath, pepath string) error {
	pc, err := anvil.LoadLevel(pcpath)
	if err != nil {
		return err
	}

	pe, err := bedrock.LoadLevel(pepath)
	if err != nil {
		return err
	}

	db, err := pe.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, pos := range pc.Regions {
		r, err := pc.LoadRegion(pos)
		log.Info().Interface("pos", pos).Msg("Processing region")
		if err != nil {
			log.Error().Err(err).Int("x", pos.X).Int("z", pos.Z).Msg("Failed to load region")
			continue
		}

		for _, c := range r.Chunks {
			if err := convChunk(db, c); err != nil {
				log.Error().Err(err).Int("x", c.X).Int("z", c.Z).Msg("Failed to convert chunk")
				continue
			}

		}
	}

	return nil
}

func FillChunk(db *bedrock.DB, pos bedrock.ChunkPos) error {
	key := make([]byte, 10)
	bedrock.MarshalChunkPos(key, pos)

	fmt.Println(hex.Dump(key))

	// Version
	key[8] = byte(bedrock.TagVersion)
	if err := db.Put(key[:9], []byte{7}); err != nil {
		return err
	}
	// Data2D
	data := make([]byte, 512+256)
	for i := 0; i < 256; i++ {
		data[2*i] = 255
	}
	// biomes
	for i := 0; i < 256; i++ {
		data[512+i] = 1
	}
	key[8] = byte(bedrock.TagData2D)
	if err := db.Put(key[:9], data); err != nil {
		return err
	}

	// FinalizedState
	key[8] = byte(bedrock.TagFinalizedState)
	if err := db.Put(key[:9], []byte{2, 0, 0, 0}); err != nil {
		return err
	}
	// // FinalizedState
	// key[8] = byte(bedrock.TagEntity)
	// if err := db.Put(key[:9], []byte{}); err != nil {
	// 	return err
	// }
	// // FinalizedState
	// key[8] = byte(bedrock.TagBlockEntity)
	// if err := db.Put(key[:9], []byte{}); err != nil {
	// 	return err
	// }

	// Sections
	for y := 0; y < 16; y++ {

		data := make([]byte, 1+4096+2048)
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				for z := 0; z < 16; z++ {
					data[(x*16+z)*16+y+1] = 7
				}
			}
		}

		key[8] = byte(bedrock.TagSection)
		key[9] = byte(y)
		if err := db.Put(key[:10], data); err != nil {
			return err
		}

	}

	return nil

}

func convChunk(db *bedrock.DB, c anvil.Chunk) error {
	key := make([]byte, 10)
	bedrock.MarshalChunkPos(key, bedrock.ChunkPos(c.ChunkPos))

	// fmt.Println(hex.Dump(key))

	// Version
	key[8] = byte(bedrock.TagVersion)
	if err := db.Put(key[:9], []byte{7}); err != nil {
		return err
	}
	// Data2D
	data := make([]byte, 512+256)
	for i := 0; i < 256; i++ {
		data[2*i] = byte(c.Data.HeightMap[i])
		// data[2*i] = 255
	}
	// biomes
	for i := 0; i < 256; i++ {
		data[512+i] = 1
		// if len(c.Data.Biomes) <= i {
		// 	// fmt.Println("biomes", len(c.Data.Biomes))
		// 	continue
		// }
		// data[512+i] = c.Data.Biomes[i]
	}
	key[8] = byte(bedrock.TagData2D)
	if err := db.Put(key[:9], data); err != nil {
		return err
	}
	// FinalizedState
	key[8] = byte(bedrock.TagFinalizedState)
	if err := db.Put(key[:9], []byte{2, 0, 0, 0}); err != nil {
		return err
	}

	// Sections
	for _, s := range c.Data.Sections {
		// for y := 0; y < 16; y++ {

		data := make([]byte, 1+4096+2048)
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				for z := 0; z < 16; z++ {
					// if s.Blocks[(y*16+z)*16+x] != 0 {
					// 	data[(x*16+z)*16+y+1] = 7
					// }
					data[(x*16+z)*16+y+1] = s.Blocks[(y*16+z)*16+x]
				}
			}
		}

		q := world.NibbleArray(data[1+4096:])
		qq := world.NibbleArray(s.Data)
		for x := 0; x < 16; x++ {
			for y := 0; y < 16; y++ {
				for z := 0; z < 16; z++ {
					q.Set(mc.BlockPos{x, y, z}, qq.Get(mc.BlockPos{y, x, z}))
				}
			}
		}
		// _, _ = q, qq

		// copy(data[1:], s.Blocks)
		// copy(data[1+4096:], s.Data)
		key[8] = byte(bedrock.TagSection)
		key[9] = s.Y
		// key[9] = byte(y)
		if err := db.Put(key[:10], data); err != nil {
			return err
		}

	}

	return nil

}
