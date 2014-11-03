package bit

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type packEntry struct {
	value, size uint32
}

type Pack struct {
	bitLength uint32
	entries   []packEntry
	packed    []byte
}

func NewPack() *Pack {
	return &Pack{entries: make([]packEntry, 0), packed: make([]byte, 4)}
}

func (p *Pack) Push(value, size uint32) {
	p.entries = append(p.entries, packEntry{value: value, size: size})
	freeBits := uint32(8*len(p.packed)) - p.bitLength
	index := p.bitLength / 32
	if freeBits >= size {
		shift := freeBits - size
		v := binary.BigEndian.Uint32(p.packed[index:]) | (value << shift)
		binary.BigEndian.PutUint32(p.packed[index:], v)
	} else {
		buf := make([]byte, 4)
		low := value
		if freeBits > 0 {
			high := value>>size - freeBits
			v := binary.BigEndian.Uint32(p.packed[index:]) | high
			binary.BigEndian.PutUint32(p.packed[index:], v)
			low = ((value << freeBits) >> freeBits)
		} else {
			shift := 32 - size
			low = (value << shift)
		}
		binary.BigEndian.PutUint32(buf, low)
		p.packed = append(p.packed[:], buf[:]...)
	}
	p.bitLength += size
}

func (p *Pack) Pack() []byte {
	return p.packed
}

func (p *Pack) PackedString() string {
	var b bytes.Buffer
	for i, w := range p.packed {
		remainingBits := int32(p.bitLength) - int32(i*8)
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

func (p *Pack) String() string {
	var buf bytes.Buffer
	buf.WriteString("<<")
	for i, entry := range p.entries {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%d:%d", entry.value, entry.size))
	}
	buf.WriteString(">>")
	return buf.String()
}

func Enc(n, B uint32, packer *Pack) *Pack {
	max := uint32(1) << uint32(B)
	if n < max {
		packer.Push(uint32(0), uint32(1))
		packer.Push(uint32(n), uint32(B))
	} else {
		packer.Push(uint32(1), uint32(1))
		Enc(n-max, B+1, packer)
	}
	return packer
}
