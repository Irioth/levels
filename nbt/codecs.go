package nbt

import (
	"encoding/binary"
	"math"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type LittleEndian struct {
	b   []byte
	pos int
}

func NewLittleEndianCodec(data []byte) *LittleEndian {
	return &LittleEndian{data, 0}
}

func (r *LittleEndian) Type() TagId {
	if r.Remaining() == 0 {
		return Tag_END
	}
	return TagId(r.Byte())
}

func (r *LittleEndian) String() string {
	l := int(r.Int16())
	s := string(r.bytes()[:l])
	r.pos += l
	return s
}

func (r *LittleEndian) Int16() int16 {
	i := binary.LittleEndian.Uint16(r.bytes())
	r.pos += 2
	return int16(i)
}

func (r *LittleEndian) Int32() int32 {
	i := binary.LittleEndian.Uint32(r.bytes())
	r.pos += 4
	return int32(i)
}

func (r *LittleEndian) Float32() float32 {
	i := binary.LittleEndian.Uint32(r.bytes())
	r.pos += 4
	return math.Float32frombits(i)
}

func (r *LittleEndian) Float64() float64 {
	i := binary.LittleEndian.Uint64(r.bytes())
	r.pos += 8
	return math.Float64frombits(i)
}

func (r *LittleEndian) Int64() int64 {
	i := binary.LittleEndian.Uint64(r.bytes())
	r.pos += 8
	return int64(i)
}

func (r *LittleEndian) Byte() byte {
	b := r.bytes()[0]
	r.pos++
	return b
}

func (r *LittleEndian) Remaining() int {
	return len(r.bytes())
}

func (r *LittleEndian) bytes() []byte {
	return r.b[r.pos:]
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type BigEndian struct {
	b   []byte
	pos int
}

func NewBigEndianCodec(data []byte) *BigEndian {
	return &BigEndian{data, 0}
}

func (r *BigEndian) Type() TagId {
	if r.Remaining() == 0 {
		return Tag_END
	}
	return TagId(r.Byte())
}

func (r *BigEndian) String() string {
	l := int(r.Int16())
	s := string(r.bytes()[:l])
	r.pos += l
	return s
}

func (r *BigEndian) Int16() int16 {
	i := binary.BigEndian.Uint16(r.bytes())
	r.pos += 2
	return int16(i)
}

func (r *BigEndian) Int32() int32 {
	i := binary.BigEndian.Uint32(r.bytes())
	r.pos += 4
	return int32(i)
}

func (r *BigEndian) Float32() float32 {
	i := binary.BigEndian.Uint32(r.bytes())
	r.pos += 4
	return math.Float32frombits(i)
}

func (r *BigEndian) Float64() float64 {
	i := binary.BigEndian.Uint64(r.bytes())
	r.pos += 8
	return math.Float64frombits(i)
}

func (r *BigEndian) Int64() int64 {
	i := binary.BigEndian.Uint64(r.bytes())
	r.pos += 8
	return int64(i)
}

func (r *BigEndian) Byte() byte {
	b := r.bytes()[0]
	r.pos++
	return b
}

func (r *BigEndian) Remaining() int {
	return len(r.bytes())
}

func (r *BigEndian) bytes() []byte {
	return r.b[r.pos:]
}
