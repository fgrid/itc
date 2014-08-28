package itc

import "encoding/binary"

type BitUnPacker struct {
	index  uint32
	packed []byte
}

func NewBitUnPacker(packed []byte) *BitUnPacker {
	return &BitUnPacker{packed: packed}
}

func (bup *BitUnPacker) Pop(size uint32) (value uint32) {
	byteIndex := (bup.index / 8) / 4
	offset := bup.index % 32
	value = binary.BigEndian.Uint32(bup.packed[byteIndex : byteIndex+4])
	shift := 32 - size
	value = uint32(value<<(offset)) >> shift
	if offset > shift {
		remain := offset - shift
		low := binary.BigEndian.Uint32(bup.packed[byteIndex+4:byteIndex+8]) >> (32 - remain)
		value = value | low
	}
	bup.index += size
	return
}
