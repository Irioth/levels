package bedrock

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/crystalmine/mapper/nbt"
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
	name, res, err := readstruct(nbt.NewReader(data[8:]))
	if err != nil {
		panic(err)
	}
	pretty.Println(name, res)
}

func readAll(r *nbt.Reader) {
	for {
		name, value, err := readstruct(r)
		if err == io.EOF {
			return
		}
		if err != nil {
			panic(err)
		}
		pretty.Println(name, value)
	}
}

func readstruct(r *nbt.Reader) (string, interface{}, error) {
	tag := r.Type()
	if tag == nbt.Tag_END {
		return "", nil, io.EOF
	}
	name := r.String()
	value, err := readValue(tag, r)
	return string(name), value, err
}

func readValue(tag nbt.TagType, r *nbt.Reader) (interface{}, error) {
	switch tag {
	case nbt.Tag_Byte:
		return r.Byte(), nil
	case nbt.Tag_Short:
		return r.Int16(), nil
	case nbt.Tag_Int:
		return r.Int32(), nil
	case nbt.Tag_Long:
		return r.Int64(), nil
	case nbt.Tag_Float:
		return r.Float32(), nil
	case nbt.Tag_String:
		return string(r.String()), nil
	case nbt.Tag_List:
		listTag := r.Type()
		l := int(r.Int32())
		res := make([]interface{}, l)
		for i := 0; i < l; i++ {
			value, err := readValue(listTag, r)
			if err != nil {
				return nil, err
			}
			res[i] = value
		}
		return res, nil
	case nbt.Tag_Compound:
		res := make(map[string]interface{})
		for {
			name, value, err := readstruct(r)
			if err == io.EOF {
				return res, nil
			}
			if err != nil {
				return nil, err
			}
			res[name] = value
		}
	default:
		return nil, fmt.Errorf("Unknown tag %s", tag)
	}

}
