package itc

import "fmt"

func ExamplePackSimple() {
	bp := newBitPack()
	bp.push(2, 3)
	bp.push(0, 1)
	bp.push(1, 2)
	source, dest := exercisePack(bp)
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<2:3, 0:1, 1:2>>) = 010001
}

func ExamplePackLarger() {
	bp := newBitPack()
	bp.push(2, 2)
	bp.push(2, 2)
	bp.push(2, 2)
	bp.push(2, 2)
	bp.push(2, 2)
	source, dest := exercisePack(bp)
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<2:2, 2:2, 2:2, 2:2, 2:2>>) = 1010101010
}

func ExamplePackUneven() {
	bp := newBitPack()
	bp.push(6, 3)
	bp.push(6, 3)
	bp.push(6, 3)
	source, dest := exercisePack(bp)
	fmt.Printf("Pack(%s) = %s\n", source, dest)
	// Output:
	// Pack(<<6:3, 6:3, 6:3>>) = 110110110
}

func exercisePack(bp *bitPack) (string, string) {
	return bp.String(), bp.PackedString()
}

func ExampleBitPacke_EncodeId_Leaves() {
	source0 := newID().asLeaf(0)
	source1 := newID().asLeaf(1)
	fmt.Printf("enc(%s) = %s\n", source0, newBitPack().encodeID(source0))
	fmt.Printf("enc(%s) = %s\n", source1, newBitPack().encodeID(source1))
	// Output:
	// enc(0) = <<0:2, 0:1>>
	// enc(1) = <<0:2, 1:1>>
}

func ExampleBitPack_EnocdeId_Nodes() {
	sourceL := newID().asNode(0, 1)
	sourceR := newID().asNode(1, 0)
	sourceB := newID().asNode(1, 1)
	fmt.Printf("enc(%s) = %s\n", sourceL, newBitPack().encodeID(sourceL))
	fmt.Printf("enc(%s) = %s\n", sourceR, newBitPack().encodeID(sourceR))
	fmt.Printf("enc(%s) = %s\n", sourceB, newBitPack().encodeID(sourceB))
	// Output:
	// enc((0, 1)) = <<1:2, 0:2, 1:1>>
	// enc((1, 0)) = <<2:2, 0:2, 1:1>>
	// enc((1, 1)) = <<3:2, 0:2, 1:1, 0:2, 1:1>>
}

func ExampleBitPack_EncodeEvent_Leaves() {
	source0 := newLeafEvent(1)
	source1 := newLeafEvent(4)
	source2 := newLeafEvent(8)
	source3 := newLeafEvent(13)
	fmt.Printf("enc(%2s, 2) = %s\n", source0, newBitPack().encodeEvent(source0))
	fmt.Printf("enc(%2s, 2) = %s\n", source1, newBitPack().encodeEvent(source1))
	fmt.Printf("enc(%2s, 2) = %s\n", source2, newBitPack().encodeEvent(source2))
	fmt.Printf("enc(%2s, 2) = %s\n", source3, newBitPack().encodeEvent(source3))
	// Output:
	// enc( 1, 2) = <<1:1, 0:1, 1:2>>
	// enc( 4, 2) = <<1:1, 1:1, 0:1, 0:3>>
	// enc( 8, 2) = <<1:1, 1:1, 0:1, 4:3>>
	// enc(13, 2) = <<1:1, 1:1, 1:1, 0:1, 1:4>>
}

func ExampleBitPack_EncodeEvent_Nodes() {
	source0 := newNodeEvent(0, 0, 1)
	source1 := newNodeEvent(0, 1, 0)
	source2 := newNodeEvent(0, 1, 1)
	source3 := newNodeEvent(1, 0, 1)
	source4 := newNodeEvent(1, 1, 0)
	source5 := newNodeEvent(1, 1, 1)

	fmt.Printf("enc(%2s) = %s\n", source0, newBitPack().encodeEvent(source0))
	fmt.Printf("enc(%2s) = %s\n", source1, newBitPack().encodeEvent(source1))
	fmt.Printf("enc(%2s) = %s\n", source2, newBitPack().encodeEvent(source2))
	fmt.Printf("enc(%2s) = %s\n", source3, newBitPack().encodeEvent(source3))
	fmt.Printf("enc(%2s) = %s\n", source4, newBitPack().encodeEvent(source4))
	fmt.Printf("enc(%2s) = %s\n", source5, newBitPack().encodeEvent(source5))

	// Output:
	// enc((0, 0, 1)) = <<0:1, 0:2, 1:1, 0:1, 1:2>>
	// enc((0, 1, 0)) = <<0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((0, 1, 1)) = <<0:1, 2:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 0, 1)) = <<0:1, 3:2, 0:1, 0:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 1, 0)) = <<0:1, 3:2, 0:1, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 1, 1)) = <<0:1, 3:2, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
}
