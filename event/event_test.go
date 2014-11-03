package event

import (
	"fmt"
	"github.com/fgrid/itc/bit"
	"testing"
)

func TestEventAsLeaf(t *testing.T) {
	e := NewLeaf(one)
	if !e.IsLeaf {
		t.Error("event is not recognized as leaf after \"setAsLeaf()\"")
	}
	if e.Left != nil || e.Right != nil {
		t.Error("Left or Right is filled - should not in a leaf")
	}
}

func TestEventAsNode(t *testing.T) {
	e := NewNode(one, one, one)
	if e.IsLeaf {
		t.Error("event is not recognized as node after \"setAsNode()\"")
	}
	if e.Left == nil || e.Right == nil {
		t.Error("Left or Right is not filled - should in a node")
	}
}

func TestEventLeafStringer(t *testing.T) {
	eventString := NewLeaf(zero).String()
	if eventString != "0" {
		t.Errorf("leaf event did not serialize as expected %q", eventString)
	}
}

func TestEventNodeStringer(t *testing.T) {
	event := NewNode(zero, one, two)
	eventString := event.String()
	if eventString != "(0, 1, 2)" {
		t.Errorf("node event did not serialize as expected %q", eventString)
	}
}

func ExampleLiftLeafEvent() {
	event := NewLeaf(uint32(4))
	sourceString := event.String()
	fmt.Printf("lift(%s, 3) = %s", sourceString, event.lift(three))
	// Output:
	// lift(4, 3) = 7
}

func ExampleLiftNodeEvent() {
	event := NewNode(one, two, three)
	sourceString := event.String()
	fmt.Printf("lift(%s, 3) = %s", sourceString, event.lift(three))
	// Output:
	// lift((1, 2, 3), 3) = (4, 2, 3)
}

func ExampleSinkLeafEvent() {
	event := NewLeaf(uint32(4))
	sourceString := event.String()
	fmt.Printf("sink(%s, 3) = %s", sourceString, event.sink(three))
	// Output:
	// sink(4, 3) = 1
}

func ExampleSinkNodeEvent() {
	event := NewNode(uint32(4), two, three)
	sourceString := event.String()
	fmt.Printf("sink(%s, 3) = %s", sourceString, event.sink(three))
	// Output:
	// sink((4, 2, 3), 3) = (1, 2, 3)
}

func ExampleNormLeafEvent() {
	event := NewLeaf(uint32(4))
	sourceString := event.String()
	fmt.Printf("Norm(%s) = %s", sourceString, event.Norm())
	// Output:
	// Norm(4) = 4
}

func ExampleNormNodeEventWithLeaves() {
	event := NewNode(two, one, one)
	sourceString := event.String()
	fmt.Printf("Norm(%s) = %s", sourceString, event.Norm())
	// Output:
	// Norm((2, 1, 1)) = 3
}

func ExampleNormNodeEventWithNodes() {
	event := NewNode(two, one, one)
	event.Left = NewNode(two, one, zero)
	event.Right = NewLeaf(three)
	sourceString := event.String()
	fmt.Printf("Norm(%s) = %s", sourceString, event.Norm())
	// Output:
	// Norm((2, (2, 1, 0), 3)) = (4, (0, 1, 0), 1)
}

func ExampleMinOfLeafEvent() {
	event := NewLeaf(uint32(4))
	fmt.Printf("Min(%s) = %d", event, event.Min())
	// Output:
	// Min(4) = 4
}

func ExampleJoinLeafEvents() {
	e1 := NewLeaf(uint32(7))
	e2 := NewLeaf(uint32(9))
	fmt.Printf("Join(%s, %s) = %s\n", e1, e2, Join(e1, e2))
	// Output:
	// Join(7, 9) = 9
}

func ExampleJoinNodeEvents() {
	e1 := NewNode(one, two, three)
	e2 := NewNode(uint32(4), uint32(5), uint32(6))
	fmt.Printf("Join(%s, %s) = %s\n", e1, e2, Join(e1, e2))
	// Output:
	// Join((1, 2, 3), (4, 5, 6)) = (9, 0, 1)
}

func ExampleEqualsNilEvents() {
	var e1, e2 *Event
	fmt.Printf("e1 = %q, e2 = %q => e1.Equals(e2) = %t", e1, e2, e1.Equals(e2))
	// Output:
	// e1 = <nil>, e2 = <nil> => e1.Equals(e2) = true
}

func ExampleLeqLeafEvents() {
	e1 := NewLeaf(one)
	e2 := NewLeaf(two)
	s2 := NewLeaf(two)

	fmt.Printf("LEQ(%s, %s) = %t\n", e1, e2, LEQ(e1, e2))
	fmt.Printf("LEQ(%s, %s) = %t\n", e2, e1, LEQ(e2, e1))
	fmt.Printf("LEQ(%s, %s) = %t\n", e2, s2, LEQ(e2, s2))

	// Output:
	// LEQ(1, 2) = true
	// LEQ(2, 1) = false
	// LEQ(2, 2) = true
}

func ExampleLeqLeftLeaf() {
	e1 := NewLeaf(one)
	e2 := NewNode(two, one, one)
	e3 := NewLeaf(three)
	e4 := NewNode(one, one, one)

	fmt.Printf("LEQ(%s, %s) = %t\n", e1, e2, LEQ(e1, e2))
	fmt.Printf("LEQ(%s, %s) = %t\n", e1, e4, LEQ(e1, e4))
	fmt.Printf("LEQ(%s, %s) = %t\n", e3, e2, LEQ(e3, e2))

	// Output:
	// LEQ(1, (2, 1, 1)) = true
	// LEQ(1, (1, 1, 1)) = true
	// LEQ(3, (2, 1, 1)) = false
}

func ExampleLeqLeftNode() {
	e1 := NewNode(two, one, one)
	e2 := NewLeaf(two)
	e3 := NewNode(two, one, one)

	fmt.Printf("LEQ(%s, %s) = %t\n", e1, e2, LEQ(e1, e2))
	fmt.Printf("LEQ(%s, %s) = %t\n", e1, e3, LEQ(e1, e3))

	// Output:
	// LEQ((2, 1, 1), 2) = false
	// LEQ((2, 1, 1), (2, 1, 1)) = true
}

func ExampleBitPack_EncodeEvent_Leaves() {
	source0 := NewLeaf(one)
	source1 := NewLeaf(uint32(4))
	source2 := NewLeaf(uint32(8))
	source3 := NewLeaf(uint32(13))
	pack0 := bit.NewPack()
	pack1 := bit.NewPack()
	pack2 := bit.NewPack()
	pack3 := bit.NewPack()
	source0.Pack(pack0)
	source1.Pack(pack1)
	source2.Pack(pack2)
	source3.Pack(pack3)
	fmt.Printf("enc(%2s, 2) = %s\n", source0, pack0)
	fmt.Printf("enc(%2s, 2) = %s\n", source1, pack1)
	fmt.Printf("enc(%2s, 2) = %s\n", source2, pack2)
	fmt.Printf("enc(%2s, 2) = %s\n", source3, pack3)
	// Output:
	// enc( 1, 2) = <<1:1, 0:1, 1:2>>
	// enc( 4, 2) = <<1:1, 1:1, 0:1, 0:3>>
	// enc( 8, 2) = <<1:1, 1:1, 0:1, 4:3>>
	// enc(13, 2) = <<1:1, 1:1, 1:1, 0:1, 1:4>>
}

func ExampleBitPack_EncodeEvent_Nodes() {
	source0 := NewNode(zero, zero, one)
	source1 := NewNode(zero, one, zero)
	source2 := NewNode(zero, one, one)
	source3 := NewNode(one, zero, one)
	source4 := NewNode(one, one, zero)
	source5 := NewNode(one, one, one)
	pack0 := bit.NewPack()
	pack1 := bit.NewPack()
	pack2 := bit.NewPack()
	pack3 := bit.NewPack()
	pack4 := bit.NewPack()
	pack5 := bit.NewPack()
	source0.Pack(pack0)
	source1.Pack(pack1)
	source2.Pack(pack2)
	source3.Pack(pack3)
	source4.Pack(pack4)
	source5.Pack(pack5)

	fmt.Printf("enc(%2s) = %s\n", source0, pack0)
	fmt.Printf("enc(%2s) = %s\n", source1, pack1)
	fmt.Printf("enc(%2s) = %s\n", source2, pack2)
	fmt.Printf("enc(%2s) = %s\n", source3, pack3)
	fmt.Printf("enc(%2s) = %s\n", source4, pack4)
	fmt.Printf("enc(%2s) = %s\n", source5, pack5)

	// Output:
	// enc((0, 0, 1)) = <<0:1, 0:2, 1:1, 0:1, 1:2>>
	// enc((0, 1, 0)) = <<0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((0, 1, 1)) = <<0:1, 2:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 0, 1)) = <<0:1, 3:2, 0:1, 0:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 1, 0)) = <<0:1, 3:2, 0:1, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
	// enc((1, 1, 1)) = <<0:1, 3:2, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>
}

func ExampleBitUnPack_decodeEvent_Leaves() {
	packer := bit.NewPack()
	NewLeaf(zero).Pack(packer)
	unpacker := bit.NewUnPack(packer.Pack())
	event0 := UnPack(unpacker)
	fmt.Printf("dec(%s) = %s\n", packer, event0)

	packer = bit.NewPack()
	NewLeaf(one).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	event0 = UnPack(unpacker)
	fmt.Printf("dec(%s) = %s\n", packer, event0)

	packer = bit.NewPack()
	NewLeaf(uint32(13)).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	event0 = UnPack(unpacker)
	fmt.Printf("dec(%s) = %s\n", packer, event0)

	// Output:
	// dec(<<1:1, 0:1, 0:2>>) = 0
	// dec(<<1:1, 0:1, 1:2>>) = 1
	// dec(<<1:1, 1:1, 1:1, 0:1, 1:4>>) = 13
}

func ExampleBitUnPack_decodeEvent_Nodes() {
	packer := bit.NewPack()
	NewNode(zero, zero, one).Pack(packer)
	unpacker := bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	NewNode(zero, one, zero).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	NewNode(zero, one, one).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	NewNode(one, zero, one).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	NewNode(one, one, zero).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	packer = bit.NewPack()
	NewNode(one, one, one).Pack(packer)
	unpacker = bit.NewUnPack(packer.Pack())
	fmt.Printf("dec(%s) = %s\n", packer, UnPack(unpacker))

	// Output:
	// dec(<<0:1, 0:2, 1:1, 0:1, 1:2>>) = (0, 0, 1)
	// dec(<<0:1, 1:2, 1:1, 0:1, 1:2>>) = (0, 1, 0)
	// dec(<<0:1, 2:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (0, 1, 1)
	// dec(<<0:1, 3:2, 0:1, 0:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 0, 1)
	// dec(<<0:1, 3:2, 0:1, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 1, 0)
	// dec(<<0:1, 3:2, 1:1, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2, 1:1, 0:1, 1:2>>) = (1, 1, 1)
}
