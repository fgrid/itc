package id

import (
	"fmt"
	"github.com/fgrid/itc/bit"
	"testing"
)

func TestIdAsLeaf(t *testing.T) {
	i := New()
	i.asLeaf(zero)
	if !i.IsLeaf {
		t.Error("id is not recognized as leaf after \"asLeaf()\"")
	}
	if i.Left != nil || i.Right != nil {
		t.Error("Left or Right is filled - should not in a leaf")
	}
}

func TestIdLeafStringer(t *testing.T) {
	idString := New().asLeaf(one).String()
	if idString != "1" {
		t.Errorf("leaf id did not serialize as expected %q", idString)
	}
}

func TestIdNodeStringer(t *testing.T) {
	id := New().asNode(one, two)
	idString := id.String()
	if idString != "(1, 2)" {
		t.Errorf("node id did not serialize as expected %q", idString)
	}
}

func ExampleSplitId0() {
	source := NewWithValue(zero)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split(0) = (0, 0)
}

func ExampleSplitId1() {
	source := NewWithValue(one)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split(1) = ((1, 0), (0, 1))
}

func ExampleSplitIdLeafNode() {
	source := New().asNode(zero, one)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split((0, 1)) = ((0, (1, 0)), (0, (0, 1)))
}

func ExampleSplitIdNodeLeaf() {
	source := New().asNode(one, zero)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split((1, 0)) = (((1, 0), 0), ((0, 1), 0))
}

func ExampleSplitIdNodeNode() {
	source := New().asNode(one, one)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split((1, 1)) = ((1, 0), (0, 1))
}

func ExampleNormalizeIdZero() {
	source := New().asNode(zero, zero)
	sourceString := source.String()
	fmt.Printf("norm(%s) = %s\n", sourceString, source.Norm())
	// Output:
	// norm((0, 0)) = 0
}

func ExampleNormalizeIdOne() {
	source := New().asNode(one, one)
	sourceString := source.String()
	fmt.Printf("norm(%s) = %s\n", sourceString, source.Norm())
	// Output:
	// norm((1, 1)) = 1
}

func ExampleSumIdLeaf() {
	i1 := NewWithValue(zero)
	i2 := NewWithValue(one)
	fmt.Printf("sum(%s, %s) = %s\n", i1, i2, New().Sum(i1, i2))
	fmt.Printf("sum(%s, %s) = %s\n", i2, i1, New().Sum(i2, i1))
	// Output:
	// sum(0, 1) = 1
	// sum(1, 0) = 1
}

func ExampleSumIdNode() {
	i1, i2 := New().Split()
	fmt.Printf("sum(%s, %s) = %s\n", i1, i2, New().Sum(i1, i2))
	// Output:
	// sum((1, 0), (0, 1)) = 1
}

func ExampleBitPack_EncodeId_Leaves() {
	source0 := New().asLeaf(zero)
	source1 := New().asLeaf(one)
	pack0 := bit.NewPack()
	pack1 := bit.NewPack()
	source0.Pack(pack0)
	source1.Pack(pack1)
	fmt.Printf("enc(%s) = %s\n", source0, pack0)
	fmt.Printf("enc(%s) = %s\n", source1, pack1)
	// Output:
	// enc(0) = <<0:2, 0:1>>
	// enc(1) = <<0:2, 1:1>>
}

func ExampleBitPack_EnocdeId_Nodes() {
	sourceL := New().asNode(zero, one)
	sourceR := New().asNode(one, zero)
	sourceB := New().asNode(one, one)
	packL := bit.NewPack()
	packR := bit.NewPack()
	packB := bit.NewPack()
	sourceL.Pack(packL)
	sourceR.Pack(packR)
	sourceB.Pack(packB)
	fmt.Printf("enc(%s) = %s\n", sourceL, packL)
	fmt.Printf("enc(%s) = %s\n", sourceR, packR)
	fmt.Printf("enc(%s) = %s\n", sourceB, packB)
	// Output:
	// enc((0, 1)) = <<1:2, 0:2, 1:1>>
	// enc((1, 0)) = <<2:2, 0:2, 1:1>>
	// enc((1, 1)) = <<3:2, 0:2, 1:1, 0:2, 1:1>>
}

func ExampleBitUnPack_decodeID_Leaves() {
	packer := bit.NewPack()
	New().asLeaf(zero).Pack(packer)
	unpacker := bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	New().asLeaf(one).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	// Output:
	// dec(<<0:2, 0:1>>) = 0
	// dec(<<0:2, 1:1>>) = 1
}

func ExampleEncDecIdNode() {
	packer := bit.NewPack()
	New().asNode(zero, one).Pack(packer)
	unpacker := bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	New().asNode(one, zero).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	New().asNode(one, one).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	// Output:
	// dec(<<1:2, 0:2, 1:1>>) = (0, 1)
	// dec(<<2:2, 0:2, 1:1>>) = (1, 0)
	// dec(<<3:2, 0:2, 1:1, 0:2, 1:1>>) = (1, 1)
}
