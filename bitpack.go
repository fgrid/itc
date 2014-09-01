package itc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type bitPackEntry struct {
	value, size uint32
}

type bitPack struct {
	bitLength uint32
	entries   []bitPackEntry
	packed    []byte
}

func newBitPack() *bitPack {
	return &bitPack{entries: make([]bitPackEntry, 0), packed: make([]byte, 4)}
}

func (bp *bitPack) push(value, size uint32) {
	bp.entries = append(bp.entries, bitPackEntry{value: value, size: size})
	freeBits := uint32(8*len(bp.packed)) - bp.bitLength
	index := bp.bitLength / 32
	if freeBits >= size {
		shift := freeBits - size
		v := binary.BigEndian.Uint32(bp.packed[index:]) | (value << shift)
		binary.BigEndian.PutUint32(bp.packed[index:], v)
	} else {
		buf := make([]byte, 4)
		low := value
		if freeBits > 0 {
			high := value>>size - freeBits
			v := binary.BigEndian.Uint32(bp.packed[index:]) | high
			binary.BigEndian.PutUint32(bp.packed[index:], v)
			low = ((value << freeBits) >> freeBits)
		} else {
			shift := 32 - size
			low = (value << shift)
		}
		binary.BigEndian.PutUint32(buf, low)
		bp.packed = append(bp.packed[:], buf[:]...)
	}
	bp.bitLength += size
}

func (bp *bitPack) Pack() []byte {
	return bp.packed
}

func (bp *bitPack) PackedString() string {
	var b bytes.Buffer
	for i, w := range bp.packed {
		remainingBits := int32(bp.bitLength) - int32(i*8)
		if remainingBits < 0 {
			break
		}
		if remainingBits < 8 {
			b.WriteString(fmt.Sprintf("%0*b", int(remainingBits), w>>uint(8-remainingBits)))
		} else {
			b.WriteString(fmt.Sprintf("%0*b", int(8), w))
		}
	}
	return b.String()
}

func (bp *bitPack) String() string {
	var buf bytes.Buffer
	buf.WriteString("<<")
	for i, entry := range bp.entries {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%d:%d", entry.value, entry.size))
	}
	buf.WriteString(">>")
	return buf.String()
}

func (bp *bitPack) encodeEvent(e *event) *bitPack {
	if e.isLeaf {
		bp.push(1, 1)
		return enc(uint32(e.value), 2, bp)
	}

	bp.push(0, 1)
	if e.value == 0 {
		if e.left.isLeaf && e.left.value == 0 {
			bp.push(0, 2)
			return bp.encodeEvent(e.right)
		}
		if e.right.isLeaf && e.right.value == 0 {
			bp.push(1, 2)
			return bp.encodeEvent(e.left)
		}
		bp.push(2, 2)
		bp.encodeEvent(e.left)
		return bp.encodeEvent(e.right)
	}

	bp.push(3, 2)
	if e.left.isLeaf && e.left.value == 0 {
		bp.push(0, 1)
		bp.push(0, 1)
		bp.encodeEvent(newLeafEvent(e.value))
		return bp.encodeEvent(e.right)
	}
	if e.right.isLeaf && e.right.value == 0 {
		bp.push(0, 1)
		bp.push(1, 1)
		bp.encodeEvent(newLeafEvent(e.value))
		return bp.encodeEvent(e.left)
	}
	bp.push(1, 1)
	bp.encodeEvent(newLeafEvent(e.value))
	bp.encodeEvent(e.left)
	return bp.encodeEvent(e.right)
}

func enc(n, B uint32, packer *bitPack) *bitPack {
	max := uint32(1) << uint32(B)
	if n < max {
		packer.push(0, 1)
		packer.push(uint32(n), uint32(B))
	} else {
		packer.push(1, 1)
		enc(n-max, B+1, packer)
	}
	return packer
}

func (bp *bitPack) encodeID(i *id) *bitPack {
	if i.isLeaf {
		bp.push(0, 2)
		bp.push(uint32(i.value), 1)
	} else if i.left.isLeaf && i.left.value == 0 {
		bp.push(1, 2)
		bp.encodeID(i.right)
	} else if i.right.isLeaf && i.right.value == 0 {
		bp.push(2, 2)
		bp.encodeID(i.left)
	} else {
		bp.push(3, 2)
		bp.encodeID(i.left)
		bp.encodeID(i.right)
	}
	return bp
}

func (bp *bitPack) encodeStamp(s *Stamp) {
	bp.encodeID(s.id)
	bp.encodeEvent(s.event)
}
