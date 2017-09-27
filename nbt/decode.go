package nbt

import (
	"fmt"
	"io"
)

type Codec interface {
	Type() TagId
	Byte() byte
	Int16() int16
	Int32() int32
	Int64() int64
	Float32() float32
	Float64() float64
	String() string
}

type Decoder struct {
	r Codec
}

func NewDecoder(data []byte, le bool) *Decoder {
	if le {
		return &Decoder{r: NewLittleEndianCodec(data)}
	}
	return &Decoder{r: NewBigEndianCodec(data)}
}

func (d *Decoder) DecodeTag() (string, interface{}, error) {
	return d.readTag()
}

func Decode(data []byte, le bool) (string, interface{}, error) {
	return NewDecoder(data, le).DecodeTag()
}

func (d *Decoder) readTag() (string, interface{}, error) {
	r := d.r
	tag := r.Type()
	if tag == Tag_END {
		return "", nil, io.EOF
	}
	name := r.String()
	value, err := d.readValue(tag)
	return name, value, err
}

func (d *Decoder) readValue(tag TagId) (interface{}, error) {
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
	case Tag_Compound:
		res := make(map[string]interface{})
		for {
			name, value, err := d.readTag()
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
