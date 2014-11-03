package id

import (
	"fmt"
	"github.com/fgrid/itc/bit"
	"log"
)

type ID struct {
	Value       uint32
	Left, Right *ID
	IsLeaf      bool
}

const (
	zero  = uint32(0)
	one   = uint32(1)
	two   = uint32(2)
	three = uint32(3)
)

func New() *ID {
	return NewWithValue(one)
}

func NewWithValue(value uint32) *ID {
	return &ID{Value: value, IsLeaf: true}
}

func (i *ID) asLeaf(value uint32) *ID {
	i.Value, i.IsLeaf, i.Left, i.Right = value, true, nil, nil
	return i
}

func (i *ID) asNode(left, right uint32) *ID {
	return i.asNodeWithIds(NewWithValue(left), NewWithValue(right))
}

func (i *ID) asNodeWithIds(left, right *ID) *ID {
	i.Value, i.IsLeaf, i.Left, i.Right = 0, false, left, right
	return i
}

// Normalize an ID as defined in section "5.2 Normal form"
func (i *ID) Norm() *ID {
	if i.IsLeaf || !i.Left.IsLeaf || !i.Right.IsLeaf || i.Left.Value != i.Right.Value {
		return i
	}
	return i.asLeaf(i.Left.Value)
}

// Split an ID as defined in section "5.3.2 Fork"
func (i *ID) Split() (i1, i2 *ID) {

	i1 = New()
	i2 = New()

	if i.IsLeaf && i.Value == 0 {
		// split(0) = (0, 0)
		i1.Value = 0
		i2.Value = 0
		return
	}
	if i.IsLeaf && i.Value == 1 {
		// split(1) = ((1,0), (0,1))
		i1.asNode(one, zero)
		i2.asNode(zero, one)
		return
	}
	if (i.Left.IsLeaf && i.Left.Value == 0) && (!i.Right.IsLeaf || i.Right.Value == 1) {
		// split((0, i)) = ((0, i1), (0, i2)), where (i1, i2) = split(i)
		r1, r2 := i.Right.Split()
		i1.asNodeWithIds(NewWithValue(zero), r1)
		i2.asNodeWithIds(NewWithValue(zero), r2)
		return
	}
	if (!i.Left.IsLeaf || i.Left.Value == 1) && (i.Right.IsLeaf && i.Right.Value == 0) {
		// split((i, 0)) = ((i1, 0), (i2, 0)), where (i1, i2) = split(i)
		l1, l2 := i.Left.Split()
		i1.asNodeWithIds(l1, NewWithValue(zero))
		i2.asNodeWithIds(l2, NewWithValue(zero))
		return
	}
	if (!i.Left.IsLeaf || i.Left.Value == 1) && (!i.Right.IsLeaf || i.Right.Value == 1) {
		// split((i1, i2)) = ((i1, 0), (0, i2))
		i1.asNodeWithIds(i.Left, NewWithValue(zero))
		i2.asNodeWithIds(NewWithValue(zero), i.Right)
		return
	}
	log.Fatalf("unable to split ID with unexpected setup: %s", i.String())
	return
}

func (i *ID) String() string {
	if i.IsLeaf {
		return fmt.Sprintf("%d", i.Value)
	}
	return fmt.Sprintf("(%s, %s)", i.Left, i.Right)
}

func (i *ID) Sum(i1, i2 *ID) *ID {
	if i1.IsLeaf && i1.Value == 0 {
		i.Value = i2.Value
		i.Left, i.Right = i2.Left, i2.Right
		i.IsLeaf = i2.IsLeaf
		return i
	}
	if i2.IsLeaf && i2.Value == 0 {
		i.Value = i1.Value
		i.Left, i.Right = i1.Left, i1.Right
		i.IsLeaf = i1.IsLeaf
		return i
	}
	return i.asNodeWithIds(New().Sum(i1.Left, i2.Left), New().Sum(i1.Right, i2.Right)).Norm()
}

func (i *ID) Pack(p *bit.Pack) {
	if i.IsLeaf {
		p.Push(zero, two)
		p.Push(i.Value, one)
		return
	}
	if i.Left.IsLeaf && i.Left.Value == 0 {
		p.Push(one, two)
		i.Right.Pack(p)
		return
	}
	if i.Right.IsLeaf && i.Right.Value == 0 {
		p.Push(two, two)
		i.Left.Pack(p)
		return
	}
	p.Push(three, two)
	i.Left.Pack(p)
	i.Right.Pack(p)
	return
}

func UnPack(bup *bit.UnPack) *ID {
	i := New()
	switch bup.Pop(two) {
	case 0:
		i.asLeaf(bup.Pop(one))
	case 1:
		newID := UnPack(bup)
		i.asNodeWithIds(NewWithValue(zero), newID)
	case 2:
		newID := UnPack(bup)
		i.asNodeWithIds(newID, NewWithValue(zero))
	case 3:
		newLeft := UnPack(bup)
		newRight := UnPack(bup)
		i.asNodeWithIds(newLeft, newRight)
	}
	return i
}
