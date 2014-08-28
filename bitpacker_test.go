package itc

import "fmt"

func ExamplePackSimple() {
	source, dest := exercisePack(NewBitPacker().Push(2, 3).Push(0, 1).Push(1, 2))
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<2:3, 0:1, 1:2>>) = 010001
}

func ExamplePackLarger() {
	source, dest := exercisePack(NewBitPacker().Push(2, 2).Push(2, 2).Push(2, 2).Push(2, 2).Push(2, 2))
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<2:2, 2:2, 2:2, 2:2, 2:2>>) = 1010101010
}

func ExamplePackUneven() {
	source, dest := exercisePack(NewBitPacker().Push(6, 3).Push(6, 3).Push(6, 3))
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<6:3, 6:3, 6:3>>) = 110110110
}

func exercisePack(bp *BitPacker) (string, string) {
	return bp.String(), bp.PackedString()
}
