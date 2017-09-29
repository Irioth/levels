package convert

import (
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
		log.Info().Interface("pos", pos).Msg("Region loaded")
		break
	}

	return nil
}

func convChunk(db *bedrock.DB, c anvil.Chunk) error {
	key := make([]byte, 10)
	bedrock.MarshalChunkPos(key, bedrock.ChunkPos(c.ChunkPos))

	fmt.Println(hex.Dump(key))

	// Version
	key[8] = byte(bedrock.TagVersion)
	if err := db.Put(key[:9], []byte{7}); err != nil {
		return err
	}
	// Data2D
	data := make([]byte, 512+256)
	for i := 0; i < 256; i++ {
		data[2*i] = byte(c.Data.HeightMap[i])
		data[512+i] = c.Data.Biomes[i]
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

	// SubChunks
	for i, s := range c.Data.Sections {

	}

	return nil

}
