//go:generate stringer --type=TagType
package nbt

import (
	"encoding/binary"
	// "encoding/hex"
	// "io"
	"fmt"
	"math"
)

type TagType byte

const (
	Tag_END TagType = iota
	Tag_Byte
	Tag_Short
	Tag_Int
	Tag_Long
	Tag_Float
	_
	_
	Tag_String
	Tag_List
	Tag_Compound
)

type Writer struct {
	b   []byte
	pos int
}

func NewWriter(data []byte) *Writer {
	return &Writer{data, 0}
}

func (w *Writer) Cut() []byte {
	return w.b[:w.pos]
}

func (w *Writer) bytes() []byte {
	return w.b[w.pos:]
}

func (w *Writer) Byte(v byte) {
	w.bytes()[0] = v
	w.pos++
}

func (w *Writer) Type(t TagType) {
	w.Byte(byte(t))
}

func (w *Writer) Int32(v int32) {
	n := binary.PutVarint(w.bytes(), int64(v))
	w.pos += n
}

func (w *Writer) Int16(v int16) {
	binary.LittleEndian.PutUint16(w.bytes(), uint16(v))
	w.pos += 2
}

func (w *Writer) Float32(v float32) {
	binary.LittleEndian.PutUint32(w.bytes(), math.Float32bits(v))
	w.pos += 4
}

func (w *Writer) Uint(v int) {
	n := binary.PutUvarint(w.bytes(), uint64(v))
	w.pos += n
}

func (w *Writer) String(s ByteString) {
	w.Uint(len(s))
	copy(w.bytes(), s)
	w.pos += len(s)
}

func ConvertToNet(data []byte) []byte {
	r := NewReader(data)
	w := NewWriter(make([]byte, len(data)))
	ConvertTag(r, w)
	return w.Cut()
}

func ConvertToNetAll(data []byte) []byte {
	r := NewReader(data)
	w := NewWriter(make([]byte, len(data)))
	for r.Remaining() > 0 {
		ConvertTag(r, w)
	}
	return w.Cut()
}

func ConvertTag(r *Reader, w *Writer) TagType {
	ttype := r.Type()
	w.Type(ttype)
	if ttype == Tag_END {
		return ttype
	}
	w.String(r.String())
	ConvertValue(ttype, r, w)
	return ttype
}

func ConvertValue(ttype TagType, r *Reader, w *Writer) {
	switch ttype {
	case Tag_Byte:
		w.Byte(r.Byte())
	case Tag_Short:
		w.Int16(r.Int16())
	case Tag_Int:
		w.Int32(r.Int32())
	// case Tag_Long:
	// w.Int64(r.Int64())
	case Tag_Float:
		w.Float32(r.Float32())
	case Tag_String:
		w.String(r.String())
	case Tag_List:
		tt := r.Type()
		w.Type(tt)
		l := r.Int32()
		w.Int32(l)
		for i := 0; i < int(l); i++ {
			ConvertValue(tt, r, w)
		}

	case Tag_Compound:
		for ConvertTag(r, w) != Tag_END {
		}
	default:
		panic(fmt.Errorf("unk tag %s", ttype))
	}
}
func ReadNet(data []byte) {
	r := NewNetReader(data)
	for r.Remaining() > 0 {
		fmt.Println("==============================")
		ReadTag(r)
	}
}

func ReadAll(data []byte) {
	r := NewReader(data)
	for r.Remaining() > 0 {
		fmt.Println("==============================")
		ReadTag(r)
	}
}

func Read(data []byte) {
	r := NewReader(data)
	ReadTag(r)
}

func ReadTag(r *Reader) TagType {
	ttype := r.Type()
	fmt.Println(ttype)
	if ttype == Tag_END {
		return ttype
	}
	name := r.String()
	fmt.Println(name)
	ReadValue(ttype, r)
	return ttype
}

func ReadValue(ttype TagType, r *Reader) {
	switch ttype {
	case Tag_Byte:
		fmt.Println(r.Byte())
	case Tag_Short:
		fmt.Println(r.Int16())
	case Tag_Int:
		fmt.Println(r.Int32())
	case Tag_Long:
		fmt.Println(r.Int64())
	case Tag_Float:
		fmt.Println(r.Float32())
	case Tag_String:
		fmt.Println(r.String())
	case Tag_List:
		tt := r.Type()
		l := int(r.Int32())
		fmt.Println(tt, l)
		for i := 0; i < l; i++ {
			ReadValue(tt, r)
		}
	case Tag_Compound:
		for ReadTag(r) != Tag_END {
		}
	default:
		panic("unk tag")
	}
}

type ByteString []byte

func (b ByteString) String() string {
	return string(b)
}

type Reader struct {
	b   []byte
	pos int
	net bool
}

func NewNetReader(data []byte) *Reader {
	return &Reader{data, 0, true}
}

func NewReader(data []byte) *Reader {
	return &Reader{data, 0, false}
}

func (r *Reader) Type() TagType {
	if r.Remaining() == 0 {
		return Tag_END
	}
	return TagType(r.Byte())
}

func (r *Reader) String() ByteString {
	l := int(r.Uint16())
	s := ByteString(r.bytes()[:l])
	r.pos += l
	return s
}

func (r *Reader) Uint16() uint16 {
	if r.net {
		i, n := binary.Uvarint(r.bytes())
		r.pos += n
		return uint16(i)
	} else {
		i := binary.LittleEndian.Uint16(r.bytes())
		r.pos += 2
		return i
	}
}

func (r *Reader) Int16() int16 {
	i := binary.LittleEndian.Uint16(r.bytes())
	r.pos += 2
	return int16(i)
}

func (r *Reader) Int32() int32 {
	if r.net {
		i, n := binary.Varint(r.bytes())
		r.pos += n
		return int32(i)
	} else {
		i := binary.LittleEndian.Uint32(r.bytes())
		r.pos += 4
		return int32(i)
	}
}

func (r *Reader) Float32() float32 {
	i := binary.LittleEndian.Uint32(r.bytes())
	r.pos += 4
	return math.Float32frombits(i)
}

func (r *Reader) Int64() int64 {
	i := binary.LittleEndian.Uint64(r.bytes())
	r.pos += 8
	return int64(i)
}

func (r *Reader) Byte() byte {
	b := r.bytes()[0]
	r.pos++
	return b
}

func (r *Reader) Skip(tag TagType) {
	switch tag {
	case Tag_Byte:
		r.Byte()
	case Tag_Short:
		r.Int16()
	case Tag_Int:
		r.Int32()
	case Tag_Long:
		r.Int64()
	case Tag_Float:
		r.Float32()
	case Tag_String:
		r.String()
	case Tag_List:
		tt := r.Type()
		l := int(r.Int32())
		for i := 0; i < l; i++ {
			r.Skip(tt)
		}
	case Tag_Compound:
		for {
			tt := r.Type()
			if tt == Tag_END {
				break
			}
			r.String()
			r.Skip(tt)
		}
	default:
		panic("unk tag")
	}
}

func (r *Reader) Remaining() int {
	return len(r.bytes())
}

func (r *Reader) bytes() []byte {
	return r.b[r.pos:]
}

func (r *Reader) Cut() []byte {
	d := r.b[:r.pos]
	r.b = r.bytes()
	r.pos = 0
	return d
}
