package bedrock

import (
	"crystal/nbt"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

func ReadNbtFile(fname string) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	// fmt.Println(hex.Dump(data))

	d := binary.LittleEndian.Uint32(data)
	fmt.Println("File version", d)
	q := binary.LittleEndian.Uint32(data[4:])
	fmt.Println("Nbt Size", q)

	nbt.ReadAll(data[8:])
}
