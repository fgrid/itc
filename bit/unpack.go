package bit

import (
	"encoding/binary"
)

type UnPack struct {
	index  uint32
	packed []byte
}

func NewUnPack(packed []byte) *UnPack {
	return &UnPack{packed: packed}
}

func (bup *UnPack) Pop(size uint32) (value uint32) {
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

func Dec(B uint32, unpacker *UnPack) uint32 {
	max := uint32(1) << uint32(B)
	if unpacker.Pop(uint32(1)) == 0 {
		return unpacker.Pop(B)
	}
	return max + Dec(B+1, unpacker)
}
