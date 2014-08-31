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
