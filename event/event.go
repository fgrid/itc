package event

import (
	"fmt"
	"github.com/fgrid/itc/bit"
)

type Event struct {
	Value       uint32
	Left, Right *Event
	IsLeaf      bool
}

const (
	zero  = uint32(0)
	one   = uint32(1)
	two   = uint32(2)
	three = uint32(3)
)

func New() *Event {
	return NewLeaf(zero)
}

func NewLeaf(value uint32) *Event {
	return &Event{Value: value, IsLeaf: true}
}

func NewEmptyNode(value uint32) *Event {
	return &Event{Value: value, IsLeaf: false, Left: New(), Right: New()}
}

func NewNode(value, left, right uint32) *Event {
	return &Event{Value: value, IsLeaf: false, Left: NewLeaf(left), Right: NewLeaf(right)}
}

func (e *Event) Clone() *Event {
	result := New()
	result.IsLeaf = e.IsLeaf
	result.Value = e.Value
	if e.Left != nil {
		result.Left = e.Left.Clone()
	}
	if e.Right != nil {
		result.Right = e.Right.Clone()
	}
	return result
}

func (e *Event) Equals(o *Event) bool {
	return (e == nil && o == nil) ||
		((e.IsLeaf == o.IsLeaf) &&
			(e.Value == o.Value) &&
			e.Left.Equals(o.Left) &&
			e.Right.Equals(o.Right))
}

func (e *Event) Norm() *Event {
	if e.IsLeaf {
		return e
	}
	if e.Left.IsLeaf && e.Right.IsLeaf && e.Left.Value == e.Right.Value {
		return NewLeaf(e.Value + e.Left.Value)
	}
	m := Min(e.Left.Min(), e.Right.Min())
	e.Left = e.Left.Norm().sink(m)
	e.Right = e.Right.Norm().sink(m)
	return e.lift(m)
}

func (e *Event) lift(value uint32) *Event {
	result := e.Clone()
	result.Value += value
	return result
}

func (e *Event) Max() uint32 {
	if e.IsLeaf {
		return e.Value
	}
	return e.Value + Max(e.Left.Max(), e.Right.Max())
}

func (e *Event) Min() uint32 {
	if e.IsLeaf {
		return e.Value
	}
	return e.Value + Min(e.Left.Min(), e.Right.Min())
}

func (e *Event) sink(value uint32) *Event {
	result := e.Clone()
	result.Value -= value
	return result
}

func (e *Event) String() string {
	if e.IsLeaf {
		return fmt.Sprintf("%d", e.Value)
	}
	return fmt.Sprintf("(%d, %s, %s)", e.Value, e.Left, e.Right)
}

// ----------
func Join(e1, e2 *Event) *Event {
	if e1.IsLeaf && e2.IsLeaf {
		return NewLeaf(Max(e1.Value, e2.Value))
	}
	if e1.IsLeaf {
		return Join(NewEmptyNode(e1.Value), e2)
	}
	if e2.IsLeaf {
		return Join(e1, NewEmptyNode(e2.Value))
	}
	if e1.Value > e2.Value {
		return Join(e2, e1)
	}
	e := NewEmptyNode(e1.Value)
	e.Left = Join(e1.Left, e2.Left.lift(e2.Value-e1.Value))
	e.Right = Join(e1.Right, e2.Right.lift(e2.Value-e1.Value))
	return e.Norm()
}

func LEQ(e1, e2 *Event) bool {
	if e1.IsLeaf {
		return e1.Value <= e2.Value
	}
	if e2.IsLeaf {
		return (e1.Value <= e2.Value) &&
			LEQ(e1.Left.lift(e1.Value), e2) &&
			LEQ(e1.Right.lift(e1.Value), e2)
	}
	return (e1.Value <= e2.Value) &&
		LEQ(e1.Left.lift(e1.Value), e2.Left.lift(e2.Value)) &&
		LEQ(e1.Right.lift(e1.Value), e2.Right.lift(e2.Value))
}

func Max(n1, n2 uint32) uint32 {
	if n1 > n2 {
		return n1
	}
	return n2
}

func Min(n1, n2 uint32) uint32 {
	if n1 < n2 {
		return n1
	}
	return n2
}

func (e Event) Pack(bp *bit.Pack) {
	if e.IsLeaf {
		bp.Push(one, one)
		bit.Enc(uint32(e.Value), two, bp)
		return
	}

	bp.Push(zero, one)
	if e.Value == 0 {
		if e.Left.IsLeaf && e.Left.Value == 0 {
			bp.Push(zero, two)
			e.Right.Pack(bp)
			return
		}
		if e.Right.IsLeaf && e.Right.Value == 0 {
			bp.Push(one, two)
			e.Left.Pack(bp)
			return
		}
		bp.Push(two, two)
		e.Left.Pack(bp)
		e.Right.Pack(bp)
		return
	}

	bp.Push(three, two)
	if e.Left.IsLeaf && e.Left.Value == 0 {
		bp.Push(zero, one)
		bp.Push(zero, one)
		NewLeaf(e.Value).Pack(bp)
		e.Right.Pack(bp)
		return
	}
	if e.Right.IsLeaf && e.Right.Value == 0 {
		bp.Push(zero, one)
		bp.Push(one, one)
		NewLeaf(e.Value).Pack(bp)
		e.Left.Pack(bp)
		return
	}
	bp.Push(one, one)
	NewLeaf(e.Value).Pack(bp)
	e.Left.Pack(bp)
	e.Right.Pack(bp)
	return
}

func UnPack(bup *bit.UnPack) *Event {
	if bup.Pop(one) == 1 {
		return NewLeaf(bit.Dec(two, bup))
	}
	e := NewNode(zero, zero, zero)
	switch bup.Pop(two) {
	case 0:
		e.Right = UnPack(bup)
	case 1:
		e.Left = UnPack(bup)
	case 2:
		e.Left = UnPack(bup)
		e.Right = UnPack(bup)
	case 3:
		if bup.Pop(one) == 0 {
			if bup.Pop(one) == 0 {
				bup.Pop(one)
				e.Value = bit.Dec(two, bup)
				e.Right = UnPack(bup)
			} else {
				bup.Pop(one)
				e.Value = bit.Dec(two, bup)
				e.Left = UnPack(bup)
			}
		} else {
			bup.Pop(one)
			e.Value = bit.Dec(two, bup)
			e.Left = UnPack(bup)
			e.Right = UnPack(bup)
		}
	}
	return e
}
