package bedrock

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/crystalmine/mapper/nbt"

	"github.com/syndtr/goleveldb/leveldb"
)

type Tag byte

const (
	TagData2D   Tag = 0x2d // 45
	TagSubChunk Tag = 0x2f // 47

	TagBlockEntity Tag = 0x31 // 49
	TagEntity      Tag = 0x32 // 50

	TagBiomeState     Tag = 0x35 // 53
	TagFinalizedState Tag = 0x36 // 54

	TagVersion Tag = 0x76 // 118 'v'
)

var (
	KnownKeys = map[string]bool{
		"AutonomousEntities": true,
		"BiomeData":          true,
		"Overworld":          true,
		"mVillages":          true,
		"portals":            true,
		"~local_player":      true,
	}
)

func Scan(fname string) {
	db, err := leveldb.OpenFile(fname, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		if !KnownKeys[string(key)] {
			if len(key) == 9 {

				tag := Tag(key[8])
				switch tag {
				case TagData2D:
				case TagVersion:
					x, z := unmarshalChunkPos(key)
					fmt.Println("Chunk", x, z, "Version", value)
				case TagBlockEntity:
					readAll(nbt.NewReader(value))
					// nbt.ReadAll(value)
				case TagEntity:
					readAll(nbt.NewReader(value))
					// nbt.ReadAll(value)
				case TagBiomeState:
				case TagFinalizedState:
				default:
					fmt.Printf("Unknwon chunk tag: %#02x\n", tag)
				}
			} else if len(key) == 10 && Tag(key[8]) == TagSubChunk {
				x, z := unmarshalChunkPos(key)
				fmt.Println("SubChunk", x, z, key[9], len(value), "Version", value[0])
			} else {
				fmt.Println("Unknown key", string(key), hex.EncodeToString(key), len(value))
			}

		}

	}

}

func unmarshalChunkPos(data []byte) (x, z int) {
	return int(int32(binary.LittleEndian.Uint32(data[:4]))), int(int32(binary.LittleEndian.Uint32(data[4:8])))
}
