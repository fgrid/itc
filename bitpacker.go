package itc

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type BitPackEntry struct {
	value, size uint32
}

type BitPacker struct {
	bitLength uint32
	entries   []BitPackEntry
	packed    []byte
}

func NewBitPacker() *BitPacker {
	return &BitPacker{entries: make([]BitPackEntry, 0), packed: make([]byte, 4)}
}

func (bp *BitPacker) Push(value, size uint32) *BitPacker {
	bp.entries = append(bp.entries, BitPackEntry{value: value, size: size})
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
	return bp
}

func (bp *BitPacker) Pack() []byte {
	return bp.packed
}

func (bp *BitPacker) PackedString() string {
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

func (bp *BitPacker) String() string {
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
