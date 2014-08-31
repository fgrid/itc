package itc

import "fmt"

func ExamplePopSimple() {
	bup := newBitUnPack([]byte{0x44, 0x00, 0x00, 0x01, 0x80, 0x00, 0x00, 0x00})
	fmt.Printf("Base = 010001\n")
	fmt.Printf("Pop(3) = %d\n", bup.Pop(3))
	fmt.Printf("Pop(1) = %d\n", bup.Pop(1))
	fmt.Printf("Pop(2) = %d\n", bup.Pop(2))
	fmt.Printf("Pop(25) = %d\n", bup.Pop(25))
	fmt.Printf("Pop(2) = %d\n", bup.Pop(2))
	// Output:
	// Base = 010001
	// Pop(3) = 2
	// Pop(1) = 0
	// Pop(2) = 1
	// Pop(25) = 0
	// Pop(2) = 3
}

func ExampleBitUnPack_decodeEvent_Leaves() {
	packer := newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEventWithValue(0)), newBitUnPack(packer.Pack()).decodeEvent())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEventWithValue(1)), newBitUnPack(packer.Pack()).decodeEvent())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEventWithValue(13)), newBitUnPack(packer.Pack()).decodeEvent())
	// Output:
	// dec(<<1:1, 0:1, 0:2>>) = 0
	// dec(<<1:1, 0:1, 1:2>>) = 1
	// dec(<<1:1, 1:1, 1:1, 0:1, 1:4>>) = 13
}

func ExampleBitUnPack_decodeEvent_Nodes() {
	packer := newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEvent().asNode(0, 0, 1)), newBitUnPack(packer.Pack()).decodeEvent())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEvent().asNode(0, 1, 0)), newBitUnPack(packer.Pack()).decodeEvent())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEvent().asNode(0, 1, 1)), newBitUnPack(packer.Pack()).decodeEvent())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEvent().asNode(1, 0, 1)), newBitUnPack(packer.Pack()).decodeEvent())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEvent().asNode(1, 1, 0)), newBitUnPack(packer.Pack()).decodeEvent())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeEvent(newEvent().asNode(1, 1, 1)), newBitUnPack(packer.Pack()).decodeEvent())
	// Output:
	// dec(<<0:1, 0:2, 1:1, 0:1, 1:2>>) = (0, 0, 1)
	// dec(<<0:1, 1:2, 1:1, 0:1, 1:2>>) = (0, 1, 0)
	// dec(<<0:1, 2:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (0, 1, 1)
	// dec(<<0:1, 3:2, 0:1, 0:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 0, 1)
	// dec(<<0:1, 3:2, 0:1, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 1, 0)
	// dec(<<0:1, 3:2, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 1, 1)
}

func ExampleBitUnPack_decodeId_Leaves() {
	packer := newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeId(newId().asLeaf(0)), newBitUnPack(packer.Pack()).decodeId())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeId(newId().asLeaf(1)), newBitUnPack(packer.Pack()).decodeId())
	// Output:
	// dec(<<0:2, 0:1>>) = 0
	// dec(<<0:2, 1:1>>) = 1
}

func ExampleEncDecIdNode() {
	packer := newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeId(newId().asNode(0, 1)), newBitUnPack(packer.Pack()).decodeId())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeId(newId().asNode(1, 0)), newBitUnPack(packer.Pack()).decodeId())
	packer = newBitPack()
	fmt.Printf("dec(%s) = %s\n", packer.encodeId(newId().asNode(1, 1)), newBitUnPack(packer.Pack()).decodeId())
	// Output:
	// dec(<<1:2, 0:2, 1:1>>) = (0, 1)
	// dec(<<2:2, 0:2, 1:1>>) = (1, 0)
	// dec(<<3:2, 0:2, 1:1, 0:2, 1:1>>) = (1, 1)
}
