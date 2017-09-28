package nbt

import (
	"fmt"
	"io"
)

type Generator struct {
	r Codec
}

func NewGenerator(data []byte, le bool) *Generator {
	if le {
		return &Generator{r: NewLittleEndianCodec(data)}
	}
	return &Generator{r: NewBigEndianCodec(data)}
}

func (d *Generator) Generate() (string, interface{}, error) {
	return d.readTag()
}

func Generate(data []byte, le bool) (string, interface{}, error) {
	return NewGenerator(data, le).Generate()
}

func (d *Generator) readTag() (string, interface{}, error) {
	r := d.r
	tag := r.Type()
	if tag == Tag_END {
		return "", nil, io.EOF
	}
	name := r.String()
	fmt.Print(name, " ")
	if tag != Tag_Compound && tag != Tag_List {
		fmt.Println(tag.GoType())
	}
	value, err := d.readValue(tag)
	return name, value, err
}

func (d *Generator) readValue(tag TagId) (interface{}, error) {
	r := d.r
	switch tag {
	case Tag_Byte:
		return r.Byte(), nil
	case Tag_Short:
		return r.Int16(), nil
	case Tag_Int:
		return r.Int32(), nil
	case Tag_Long:
		return r.Int64(), nil
	case Tag_Float:
		return r.Float32(), nil
	case Tag_Double:
		return r.Float64(), nil
	case Tag_String:
		return r.String(), nil
	case Tag_List:
		listTag := r.Type()
		fmt.Println("[]" + listTag.GoType())
		l := int(r.Int32())
		res := make([]interface{}, l)
		for i := 0; i < l; i++ {
			value, err := d.readValue(listTag)
			if err != nil {
				return nil, err
			}
			res[i] = value
		}
		return res, nil
	case Tag_Byte_Array:
		l := int(r.Int32())
		res := make([]byte, l)
		for i := 0; i < l; i++ {
			res[i] = r.Byte()
		}
		return res, nil
	case Tag_Int_Array:
		l := int(r.Int32())
		res := make([]int32, l)
		for i := 0; i < l; i++ {
			res[i] = r.Int32()
		}
		return res, nil
	case Tag_Compound:
		fmt.Println("struct {")
		res := make(map[string]interface{})
		for {
			name, value, err := d.readTag()
			if err == io.EOF {
				fmt.Println("}")
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

func (t TagId) GoType() string {
	switch t {
	case Tag_Byte:
		return "byte"
	case Tag_Short:
		return "int16"
	case Tag_Int:
		return "int32"
	case Tag_Long:
		return "int64"
	case Tag_Float:
		return "float32"
	case Tag_Double:
		return "float64"
	case Tag_String:
		return "string"
	case Tag_Byte_Array:
		return "[]byte"
	case Tag_Int_Array:
		return "[]int32"
	}
	return ""
}
