package itc

import (
	"fmt"
	"log"
)

type id struct {
	value       int
	left, right *id
	isLeaf      bool
}

func newId() *id {
	return newIdWithValue(1)
}

func newIdWithValue(value int) *id {
	return &id{value: value, isLeaf: true}
}

func (i *id) asLeaf(value int) *id {
	i.value, i.isLeaf, i.left, i.right = value, true, nil, nil
	return i
}

func (i *id) asNode(left, right int) *id {
	return i.asNodeWithIds(newIdWithValue(left), newIdWithValue(right))
}

func (i *id) asNodeWithIds(left, right *id) *id {
	i.value, i.isLeaf, i.left, i.right = 0, false, left, right
	return i
}

// Normalize an ID as defined in section "5.2 Normal form"
func (i *id) Norm() *id {
	if i.isLeaf || !i.left.isLeaf || !i.right.isLeaf || i.left.value != i.right.value {
		return i
	}
	return i.asLeaf(i.left.value)
}

func (i *id) dec(unpacker *BitUnPacker) *id {
	switch unpacker.Pop(2) {
	case 0:
		i.asLeaf(int(unpacker.Pop(1)))
	case 1:
		newId := newId()
		newId.dec(unpacker)
		i.asNodeWithIds(newIdWithValue(0), newId)
	case 2:
		newId := newId()
		newId.dec(unpacker)
		i.asNodeWithIds(newId, newIdWithValue(0))
	case 3:
		newLeft := newId()
		newLeft.dec(unpacker)
		newRight := newId()
		newRight.dec(unpacker)
		i.asNodeWithIds(newLeft, newRight)
	}
	return i
}

func (i *id) enc(packer *BitPacker) *BitPacker {
	if i.isLeaf {
		packer.Push(0, 2)
		packer.Push(uint32(i.value), 1)
	} else if i.left.isLeaf && i.left.value == 0 {
		packer.Push(1, 2)
		i.right.enc(packer)
	} else if i.right.isLeaf && i.right.value == 0 {
		packer.Push(2, 2)
		i.left.enc(packer)
	} else {
		packer.Push(3, 2)
		i.left.enc(packer)
		i.right.enc(packer)
	}
	return packer
}

// Split an ID as defined in section "5.3.2 Fork"
func (i *id) Split() (i1, i2 *id) {

	i1 = newId()
	i2 = newId()

	if i.isLeaf && i.value == 0 {
		// split(0) = (0, 0)
		i1.value = 0
		i2.value = 0
	} else if i.isLeaf && i.value == 1 {
		// split(1) = ((1,0), (0,1))
		i1.asNode(1, 0)
		i2.asNode(0, 1)
	} else if (i.left.isLeaf && i.left.value == 0) && (!i.right.isLeaf || i.right.value == 1) {
		// split((0, i)) = ((0, i1), (0, i2)), where (i1, i2) = split(i)
		r1, r2 := i.right.Split()
		i1.asNodeWithIds(newIdWithValue(0), r1)
		i2.asNodeWithIds(newIdWithValue(0), r2)
	} else if (!i.left.isLeaf || i.left.value == 1) && (i.right.isLeaf && i.right.value == 0) {
		// split((i, 0)) = ((i1, 0), (i2, 0)), where (i1, i2) = split(i)
		l1, l2 := i.left.Split()
		i1.asNodeWithIds(l1, newIdWithValue(0))
		i2.asNodeWithIds(l2, newIdWithValue(0))
	} else if (!i.left.isLeaf || i.left.value == 1) && (!i.right.isLeaf || i.right.value == 1) {
		// split((i1, i2)) = ((i1, 0), (0, i2))
		i1.asNodeWithIds(i.left, newIdWithValue(0))
		i2.asNodeWithIds(newIdWithValue(0), i.right)
	} else {
		log.Fatalf("unable to split id with unexpected setup: %s", i.String())
	}
	return
}

func (i *id) String() string {
	if i.isLeaf {
		return fmt.Sprintf("%d", i.value)
	}
	return fmt.Sprintf("(%s, %s)", i.left, i.right)
}

func (i *id) sum(i1, i2 *id) *id {
	if i1.isLeaf && i1.value == 0 {
		i.value = i2.value
		i.left, i.right = i2.left, i2.right
		i.isLeaf = i2.isLeaf
		return i
	}
	if i2.isLeaf && i2.value == 0 {
		i.value = i1.value
		i.left, i.right = i1.left, i1.right
		i.isLeaf = i1.isLeaf
		return i
	}
	return i.asNodeWithIds(newId().sum(i1.left, i2.left), newId().sum(i1.right, i2.right)).Norm()
}
