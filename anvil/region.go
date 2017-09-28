package anvil

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/crystalmine/levels/nbt"
	"github.com/rs/zerolog/log"
)

type RegionPos struct {
	X, Z int
}

func findRegions(path string) ([]RegionPos, error) {
	fi, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var regions []RegionPos
	for i := range fi {
		name := fi[i].Name()
		if filepath.Ext(name) != ".mca" {
			log.Warn().Str("filename", name).Msg("Unexpected file extension")
			continue
		}
		info, err := parseRegionName(name)
		if err != nil {
			log.Warn().Err(err).Str("filename", name).Msg("Can't parse region coords")
			continue
		}
		regions = append(regions, info)
	}
	return regions, nil

}

func parseRegionName(fname string) (RegionPos, error) {

	parts := strings.Split(fname, ".")
	if len(parts) != 4 {
		return RegionPos{}, fmt.Errorf("invalid file name format")
	}
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return RegionPos{}, err
	}
	z, err := strconv.Atoi(parts[2])
	if err != nil {
		return RegionPos{}, err
	}
	return RegionPos{x, z}, nil
}

type Region struct {
	RegionPos
	Chunks map[ChunkPos]Chunk
}

func LoadRegionFile(fname string) (*Region, error) {
	pos, err := parseRegionName(fname)
	if err != nil {
		return nil, err
	}

	fileraw, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	chunks := make(map[ChunkPos]Chunk)

	off := 0
	for z := 0; z < 32; z++ {
		for x := 0; x < 32; x, off = x+1, off+4 {
			loc := binary.BigEndian.Uint32(fileraw[off:])
			if loc == 0 {
				continue
			}
			timestamp := binary.BigEndian.Uint32(fileraw[4096+off:])
			size := loc & 0xff
			offset := loc >> 8
			cd, err := loadChunkData(fileraw[4096*offset : 4096*(offset+size)])
			if err != nil {
				log.Error().Str("region", filepath.Base(fname)).Err(err).Int("x", x).Int("z", z).Msg("Failed chunk loading")
				continue
			}
			cp := ChunkPos{X: pos.X*32 + x, Z: pos.Z*32 + z}
			chunks[cp] = Chunk{
				ChunkPos:  cp,
				Timestamp: timestamp,
				Data:      cd,
			}
		}
	}

	return &Region{RegionPos: pos, Chunks: chunks}, nil

}

func loadChunkData(raw []byte) (*ChunkData, error) {
	size := binary.BigEndian.Uint32(raw)
	compression := raw[4]
	if compression != 2 {
		return nil, fmt.Errorf("Unexpected compression type %d", compression)
	}
	data := raw[5 : 5+size-1]

	unzipped, err := unzip(data)
	if err != nil {
		return nil, err
	}

	// nbt.Generate(unzipped, false)
	_, d, err := nbt.Decode(unzipped, false)
	if err != nil {
		return nil, err
	}

	type ChunkDataRoot struct{ Level ChunkData }
	var x ChunkDataRoot
	if err := mapDecode(d, &x); err != nil {
		return nil, err
	}

	return &x.Level, nil
}

func unzip(data []byte) ([]byte, error) {
	z, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(z)
}

type ChunkPos struct {
	X, Z int
}

type Chunk struct {
	ChunkPos
	Timestamp uint32
	Data      interface{}
}

type ChunkData struct {
	XPos             int32
	ZPos             int32
	LastUpdate       int64
	TerrainPopulated bool
	HeightMap        []int32
	Biomes           []byte
	Entities         []byte
	TileEntities     []byte
	Sections         []struct {
		Y          byte
		Blocks     []byte
		Data       []byte
		SkyLight   []byte
		BlockLight []byte
	}
}
