package bit

import (
	"fmt"
)

func ExamplePackSimple() {
	bp := NewPack()
	bp.Push(uint32(2), uint32(3))
	bp.Push(uint32(0), uint32(1))
	bp.Push(uint32(1), uint32(2))
	source, dest := exercisePack(bp)
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<2:3, 0:1, 1:2>>) = 010001
}

func ExamplePackLarger() {
	bp := NewPack()
	bp.Push(uint32(2), uint32(2))
	bp.Push(uint32(2), uint32(2))
	bp.Push(uint32(2), uint32(2))
	bp.Push(uint32(2), uint32(2))
	bp.Push(uint32(2), uint32(2))
	source, dest := exercisePack(bp)
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<2:2, 2:2, 2:2, 2:2, 2:2>>) = 1010101010
}

func ExamplePackUneven() {
	bp := NewPack()
	bp.Push(uint32(6), uint32(3))
	bp.Push(uint32(6), uint32(3))
	bp.Push(uint32(6), uint32(3))
	source, dest := exercisePack(bp)
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<6:3, 6:3, 6:3>>) = 110110110
}

func exercisePack(bp *Pack) (string, string) {
	return bp.String(), bp.PackedString()
}
