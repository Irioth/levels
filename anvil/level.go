package anvil

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/crystalmine/levels/nbt"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type Level struct {
	Path    string
	Data    *LevelData
	Regions []RegionPos
}

func (l *Level) LoadRegion(pos RegionPos) (*Region, error) {
	return LoadRegionFile(filepath.Join(l.Path, "region", fmt.Sprintf("r.%d.%d.mca", pos.X, pos.Z)))
}

func LoadLevel(path string) (*Level, error) {
	data, err := loadLevelDat(filepath.Join(path, "level.dat"))
	if err != nil {
		return nil, fmt.Errorf("failed to load level.dat: %v", err)
	}

	regions, err := findRegions(filepath.Join(path, "region"))
	if err != nil {
		return nil, fmt.Errorf("failed to find regions: %v", err)
	}

	return &Level{
		Path:    path,
		Data:    data,
		Regions: regions,
	}, nil
}

type LevelData struct {
	Version          int32
	Player           interface{}
	GameRules        interface{}
	LevelName        string
	GeneratorName    string
	GeneratorVersion int32
	GeneratorOptions string
	LastPlayed       int64
	RandomSeed       int64
	Time             int64
	DayTime          int64
	SizeOnDisk       int64
	SpawnX           int32
	SpawnY           int32
	SpawnZ           int32
	RainTime         int32
	ThunderTime      int32
	GameType         int32
	Initialized      bool
	MapFeatures      bool
	AllowCommands    bool
	Hardcore         bool
	Raining          bool
	Thundering       bool
}

func loadLevelDat(fname string) (*LevelData, error) {
	raw, err := readGzipped(fname)
	if err != nil {
		return nil, err
	}

	// nbt.Generate(raw, false)
	_, res, err := nbt.Decode(raw, false)
	if err != nil {
		return nil, err
	}

	type LevelDataRoot struct{ Data LevelData }
	var data LevelDataRoot
	if err := mapDecode(res, &data); err != nil {
		return nil, fmt.Errorf("failed to decode level.dat structure: %v", err)
	}

	return &data.Data, nil
}

func mapDecode(in interface{}, v interface{}) error {
	if e := log.Debug(); e.Enabled() {
		var md mapstructure.Metadata
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Metadata:         &md,
			Result:           v,
			WeaklyTypedInput: true,
		})
		if err != nil {
			return err
		}
		if err := decoder.Decode(in); err != nil {
			return err
		}
		e.Str("struct", reflect.TypeOf(v).Elem().String()).Strs("unused", md.Unused).Msg("Decoding NBT to struct.")
		return nil
	}

	if err := mapstructure.Decode(in, v); err != nil {
		return err
	}
	return nil
}

func readGzipped(fname string) ([]byte, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	return raw, nil

}
