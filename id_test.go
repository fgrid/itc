package itc

import (
	"fmt"
	"testing"
)

func TestIdAsLeaf(t *testing.T) {
	i := newId()
	i.asLeaf(0)
	if !i.isLeaf {
		t.Error("id is not recognized as leaf after \"asLeaf()\"")
	}
	if i.left != nil || i.right != nil {
		t.Error("left or right is filled - should not in a leaf")
	}
}

func TestIdLeafStringer(t *testing.T) {
	idString := newId().asLeaf(1).String()
	if idString != "1" {
		t.Errorf("leaf id did not serialize as expected %q", idString)
	}
}

func TestIdNodeStringer(t *testing.T) {
	id := newId().asNode(1, 2)
	idString := id.String()
	if idString != "(1, 2)" {
		t.Errorf("node id did not serialize as expected %q", idString)
	}
}

func ExampleSplitId0() {
	source := newIdWithValue(0)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split(0) = (0, 0)
}

func ExampleSplitId1() {
	source := newIdWithValue(1)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split(1) = ((1, 0), (0, 1))
}

func ExampleSplitIdLeafNode() {
	source := newId().asNode(0, 1)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split((0, 1)) = ((0, (1, 0)), (0, (0, 1)))
}

func ExampleSplitIdNodeLeaf() {
	source := newId().asNode(1, 0)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split((1, 0)) = (((1, 0), 0), ((0, 1), 0))
}

func ExampleSplitIdNodeNode() {
	source := newId().asNode(1, 1)
	i1, i2 := source.Split()

	fmt.Printf("split(%s) = (%s, %s)\n", source, i1, i2)
	// Output:
	// split((1, 1)) = ((1, 0), (0, 1))
}

func ExampleNormalizeIdZero() {
	source := newId().asNode(0, 0)
	sourceString := source.String()
	fmt.Printf("norm(%s) = %s\n", sourceString, source.Norm())
	// Output:
	// norm((0, 0)) = 0
}

func ExampleNormalizeIdOne() {
	source := newId().asNode(1, 1)
	sourceString := source.String()
	fmt.Printf("norm(%s) = %s\n", sourceString, source.Norm())
	// Output:
	// norm((1, 1)) = 1
}

func ExampleEncodeIdLeaf() {
	source0 := newId().asLeaf(0)
	source1 := newId().asLeaf(1)
	fmt.Printf("enc(%s) = %s\n", source0, source0.enc(NewBitPacker()))
	fmt.Printf("enc(%s) = %s\n", source1, source1.enc(NewBitPacker()))
	// Output:
	// enc(0) = <<0:2, 0:1>>
	// enc(1) = <<0:2, 1:1>>
}

func ExampleEncodeIdNode() {
	sourceL := newId().asNode(0, 1)
	sourceR := newId().asNode(1, 0)
	sourceB := newId().asNode(1, 1)
	fmt.Printf("enc(%s) = %s\n", sourceL, sourceL.enc(NewBitPacker()))
	fmt.Printf("enc(%s) = %s\n", sourceR, sourceR.enc(NewBitPacker()))
	fmt.Printf("enc(%s) = %s\n", sourceB, sourceB.enc(NewBitPacker()))
	// Output:
	// enc((0, 1)) = <<1:2, 0:2, 1:1>>
	// enc((1, 0)) = <<2:2, 0:2, 1:1>>
	// enc((1, 1)) = <<3:2, 0:2, 1:1, 0:2, 1:1>>
}

func ExampleEncDecIdLeaf() {
	packer := NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newId().asLeaf(0).enc(packer), newId().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newId().asLeaf(1).enc(packer), newId().dec(NewBitUnPacker(packer.Pack())))
	// Output:
	// dec(<<0:2, 0:1>>) = 0
	// dec(<<0:2, 1:1>>) = 1
}

func ExampleEncDecIdNode() {
	packer := NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newId().asNode(0, 1).enc(packer), newId().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newId().asNode(1, 0).enc(packer), newId().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newId().asNode(1, 1).enc(packer), newId().dec(NewBitUnPacker(packer.Pack())))
	// Output:
	// dec(<<1:2, 0:2, 1:1>>) = (0, 1)
	// dec(<<2:2, 0:2, 1:1>>) = (1, 0)
	// dec(<<3:2, 0:2, 1:1, 0:2, 1:1>>) = (1, 1)
}

func ExampleSumIdLeaf() {
	i1 := newIdWithValue(0)
	i2 := newIdWithValue(1)
	fmt.Printf("sum(%s, %s) = %s\n", i1, i2, newId().sum(i1, i2))
	fmt.Printf("sum(%s, %s) = %s\n", i2, i1, newId().sum(i2, i1))
	// Output:
	// sum(0, 1) = 1
	// sum(1, 0) = 1
}

func ExampleSumIdNode() {
	i1, i2 := newId().Split()
	fmt.Printf("sum(%s, %s) = %s\n", i1, i2, newId().sum(i1, i2))
	// Output:
	// sum((1, 0), (0, 1)) = 1
}
