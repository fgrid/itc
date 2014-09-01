package itc

import "encoding/binary"

type bitUnPack struct {
	index  uint32
	packed []byte
}

func newBitUnPack(packed []byte) *bitUnPack {
	return &bitUnPack{packed: packed}
}

func (bup *bitUnPack) Pop(size uint32) (value uint32) {
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

func (bup *bitUnPack) decodeEvent() *event {
	if bup.Pop(1) == 1 {
		return newLeafEvent(int(dec(2, bup)))
	}
	e := newNodeEvent(0, 0, 0)
	switch bup.Pop(2) {
	case 0:
		e.right = bup.decodeEvent()
	case 1:
		e.left = bup.decodeEvent()
	case 2:
		e.left = bup.decodeEvent()
		e.right = bup.decodeEvent()
	case 3:
		if bup.Pop(1) == 0 {
			if bup.Pop(1) == 0 {
				bup.Pop(1)
				e.value = int(dec(2, bup))
				e.right = bup.decodeEvent()
			} else {
				bup.Pop(1)
				e.value = int(dec(2, bup))
				e.left = bup.decodeEvent()
			}
		} else {
			bup.Pop(1)
			e.value = int(dec(2, bup))
			e.left = bup.decodeEvent()
			e.right = bup.decodeEvent()
		}
	}
	return e
}

func (bup *bitUnPack) decodeID() *id {
	i := newID()
	switch bup.Pop(2) {
	case 0:
		i.asLeaf(int(bup.Pop(1)))
	case 1:
		newID := bup.decodeID()
		i.asNodeWithIds(newIDWithValue(0), newID)
	case 2:
		newID := bup.decodeID()
		i.asNodeWithIds(newID, newIDWithValue(0))
	case 3:
		newLeft := bup.decodeID()
		newRight := bup.decodeID()
		i.asNodeWithIds(newLeft, newRight)
	}
	return i
}

func (bup *bitUnPack) decodeStamp(s *Stamp) {
	s.id = bup.decodeID()
	s.event = bup.decodeEvent()
}

func dec(B uint32, unpacker *bitUnPack) uint32 {
	max := uint32(1) << uint32(B)
	if unpacker.Pop(1) == 0 {
		return unpacker.Pop(B)
	}
	return max + dec(B+1, unpacker)
}
