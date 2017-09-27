package bedrock

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/crystalmine/levels/nbt"
	"github.com/kr/pretty"
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

	// nbt.ReadAll(data[8:])
	name, res, err := nbt.Decode(data[8:], true)
	if err != nil {
		panic(err)
	}
	pretty.Println(name, res)
}

func readAll(data []byte) {
	d := nbt.NewDecoder(data, true)
	for {
		name, value, err := d.DecodeTag()
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		pretty.Println(name, value)
	}
}
