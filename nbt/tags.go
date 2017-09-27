//go:generate stringer --type=TagId
package nbt

type TagId uint8

const (
	Tag_END      TagId = iota // 0x00
	Tag_Byte                  // 0x01
	Tag_Short                 // 0x02
	Tag_Int                   // 0x03
	Tag_Long                  // 0x04
	Tag_Float                 // 0x05
	Tag_Double                // 0x06
	_                         // 0x07
	Tag_String                // 0x08
	Tag_List                  // 0x09
	Tag_Compound              // 0x0a
)
