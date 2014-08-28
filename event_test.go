package itc

import (
	"fmt"
	"testing"
)

func TestEventAsLeaf(t *testing.T) {
	e := newEvent().asLeaf(1)
	if !e.isLeaf {
		t.Error("event is not recognized as leaf after \"setAsLeaf()\"")
	}
	if e.left != nil || e.right != nil {
		t.Error("left or right is filled - should not in a leaf")
	}
}

func TestEventAsNode(t *testing.T) {
	e := newEvent().asNode(1, 1, 1)
	if e.isLeaf {
		t.Error("event is not recognized as node after \"setAsNode()\"")
	}
	if e.left == nil || e.right == nil {
		t.Error("left or right is not filled - should in a node")
	}
}

func TestEventLeafStringer(t *testing.T) {
	eventString := newEvent().asLeaf(0).String()
	if eventString != "0" {
		t.Errorf("leaf event did not serialize as expected %q", eventString)
	}
}

func TestEventNodeStringer(t *testing.T) {
	event := newEvent().asNode(0, 1, 2)
	eventString := event.String()
	if eventString != "(0, 1, 2)" {
		t.Errorf("node event did not serialize as expected %q", eventString)
	}
}

func ExampleLiftLeafEvent() {
	event := newEventWithValue(4)
	sourceString := event.String()
	fmt.Printf("lift(%s, 3) = %s", sourceString, event.lift(3))
	// Output:
	// lift(4, 3) = 7
}

func ExampleLiftNodeEvent() {
	event := newEvent().asNode(1, 2, 3)
	sourceString := event.String()
	fmt.Printf("lift(%s, 3) = %s", sourceString, event.lift(3))
	// Output:
	// lift((1, 2, 3), 3) = (4, 2, 3)
}

func ExampleSinkLeafEvent() {
	event := newEventWithValue(4)
	sourceString := event.String()
	fmt.Printf("sink(%s, 3) = %s", sourceString, event.sink(3))
	// Output:
	// sink(4, 3) = 1
}

func ExampleSinkNodeEvent() {
	event := newEvent().asNode(4, 2, 3)
	sourceString := event.String()
	fmt.Printf("sink(%s, 3) = %s", sourceString, event.sink(3))
	// Output:
	// sink((4, 2, 3), 3) = (1, 2, 3)
}

func ExampleNormLeafEvent() {
	event := newEventWithValue(4)
	sourceString := event.String()
	fmt.Printf("norm(%s) = %s", sourceString, event.norm())
	// Output:
	// norm(4) = 4
}

func ExampleNormNodeEventWithLeaves() {
	event := newEvent().asNode(2, 1, 1)
	sourceString := event.String()
	fmt.Printf("norm(%s) = %s", sourceString, event.norm())
	// Output:
	// norm((2, 1, 1)) = 3
}

func ExampleNormNodeEventWithNodes() {
	event := newEvent().asNode(2, 1, 1)
	event.left.asNode(2, 1, 0)
	event.right.asLeaf(3)
	sourceString := event.String()
	fmt.Printf("norm(%s) = %s", sourceString, event.norm())
	// Output:
	// norm((2, (2, 1, 0), 3)) = (4, (0, 1, 0), 1)
}

func ExampleMinOfLeafEvent() {
	event := newEventWithValue(4)
	fmt.Printf("min(%s) = %d", event, event.min())
	// Output:
	// min(4) = 4
}

func ExampleEncLeafEvent() {
	source0 := newEventWithValue(1)
	source1 := newEventWithValue(4)
	source2 := newEventWithValue(8)
	source3 := newEventWithValue(13)
	fmt.Printf("enc(%2s, 2) = %s\n", source0, source0.enc(NewBitPacker()))
	fmt.Printf("enc(%2s, 2) = %s\n", source1, source1.enc(NewBitPacker()))
	fmt.Printf("enc(%2s, 2) = %s\n", source2, source2.enc(NewBitPacker()))
	fmt.Printf("enc(%2s, 2) = %s\n", source3, source3.enc(NewBitPacker()))
	// Output:
	// enc( 1, 2) = <<1:1, 0:1, 1:2>>
	// enc( 4, 2) = <<1:1, 1:1, 0:1, 0:3>>
	// enc( 8, 2) = <<1:1, 1:1, 0:1, 4:3>>
	// enc(13, 2) = <<1:1, 1:1, 1:1, 0:1, 1:4>>
}

func ExampleEncNodeEvent() {
	source0 := newEvent().asNode(0, 0, 1)
	source1 := newEvent().asNode(0, 1, 0)
	source2 := newEvent().asNode(0, 1, 1)
	source3 := newEvent().asNode(1, 0, 1)
	source4 := newEvent().asNode(1, 1, 0)
	source5 := newEvent().asNode(1, 1, 1)

	fmt.Printf("enc(%2s) = %s\n", source0, source0.enc(NewBitPacker()))
	fmt.Printf("enc(%2s) = %s\n", source1, source1.enc(NewBitPacker()))
	fmt.Printf("enc(%2s) = %s\n", source2, source2.enc(NewBitPacker()))
	fmt.Printf("enc(%2s) = %s\n", source3, source3.enc(NewBitPacker()))
	fmt.Printf("enc(%2s) = %s\n", source4, source4.enc(NewBitPacker()))
	fmt.Printf("enc(%2s) = %s\n", source5, source5.enc(NewBitPacker()))

	// Output:
	// enc((0, 0, 1)) = <<0:1, 0:2, 1:1, 0:1, 1:2>>
	// enc((0, 1, 0)) = <<0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((0, 1, 1)) = <<0:1, 2:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 0, 1)) = <<0:1, 3:2, 0:1, 0:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 1, 0)) = <<0:1, 3:2, 0:1, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 1, 1)) = <<0:1, 3:2, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
}

func ExampleDecLeafEvent() {
	packer := NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEventWithValue(0).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEventWithValue(1).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEventWithValue(13).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	// Output:
	// dec(<<1:1, 0:1, 0:2>>) = 0
	// dec(<<1:1, 0:1, 1:2>>) = 1
	// dec(<<1:1, 1:1, 1:1, 0:1, 1:4>>) = 13
}

func ExampleDevNodeEvent() {
	packer := NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEvent().asNode(0, 0, 1).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEvent().asNode(0, 1, 0).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEvent().asNode(0, 1, 1).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEvent().asNode(1, 0, 1).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEvent().asNode(1, 1, 0).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	packer = NewBitPacker()
	fmt.Printf("dec(%s) = %s\n", newEvent().asNode(1, 1, 1).enc(packer), newEvent().dec(NewBitUnPacker(packer.Pack())))
	// Output:
	// dec(<<0:1, 0:2, 1:1, 0:1, 1:2>>) = (0, 0, 1)
	// dec(<<0:1, 1:2, 1:1, 0:1, 1:2>>) = (0, 1, 0)
	// dec(<<0:1, 2:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (0, 1, 1)
	// dec(<<0:1, 3:2, 0:1, 0:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 0, 1)
	// dec(<<0:1, 3:2, 0:1, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 1, 0)
	// dec(<<0:1, 3:2, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 1, 1)
}

func ExampleJoinLeafEvents() {
	e1 := newEventWithValue(7)
	e2 := newEventWithValue(9)
	fmt.Printf("join(%s, %s) = %s\n", e1, e2, newEvent().join(e1, e2))
	// Output:
	// join(7, 9) = 9
}

func ExampleJoinNodeEvents() {
	e1 := newEvent().asNode(1, 2, 3)
	e2 := newEvent().asNode(4, 5, 6)
	fmt.Printf("join(%s, %s) = %s\n", e1, e2, newEvent().join(e1, e2))
	// Output:
	// join((1, 2, 3), (4, 5, 6)) = (9, 0, 1)
}

func ExampleEqualsNilEvents() {
	var e1, e2 *event
	fmt.Printf("e1 = %q, e2 = %q => e1.equals(e2) = %t", e1, e2, e1.equals(e2))
	// Output:
	// e1 = <nil>, e2 = <nil> => e1.equals(e2) = true
}

func ExampleLeqLeafEvents() {
	e1 := newEventWithValue(1)
	e2 := newEventWithValue(2)
	s2 := newEventWithValue(2)

	fmt.Printf("leq(%s, %s) = %t\n", e1, e2, leq(e1, e2))
	fmt.Printf("leq(%s, %s) = %t\n", e2, e1, leq(e2, e1))
	fmt.Printf("leq(%s, %s) = %t\n", e2, s2, leq(e2, s2))

	// Output:
	// leq(1, 2) = true
	// leq(2, 1) = false
	// leq(2, 2) = true
}

func ExampleLeqLeftLeaf() {
	e1 := newEventWithValue(1)
	e2 := newEvent().asNode(2, 1, 1)
	e3 := newEventWithValue(3)
	e4 := newEvent().asNode(1, 1, 1)

	fmt.Printf("leq(%s, %s) = %t\n", e1, e2, leq(e1, e2))
	fmt.Printf("leq(%s, %s) = %t\n", e1, e4, leq(e1, e4))
	fmt.Printf("leq(%s, %s) = %t\n", e3, e2, leq(e3, e2))

	// Output:
	// leq(1, (2, 1, 1)) = true
	// leq(1, (1, 1, 1)) = true
	// leq(3, (2, 1, 1)) = false
}

func ExampleLeqLeftNode() {
	e1 := newEvent().asNode(2, 1, 1)
	e2 := newEventWithValue(2)
	e3 := newEvent().asNode(2, 1, 1)

	fmt.Printf("leq(%s, %s) = %t\n", e1, e2, leq(e1, e2))
	fmt.Printf("leq(%s, %s) = %t\n", e1, e3, leq(e1, e3))

	// Output:
	// leq((2, 1, 1), 2) = false
	// leq((2, 1, 1), (2, 1, 1)) = true
}
