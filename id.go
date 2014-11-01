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

func newID() *id {
	return newIDWithValue(1)
}

func newIDWithValue(value int) *id {
	return &id{value: value, isLeaf: true}
}

func (i *id) asLeaf(value int) *id {
	i.value, i.isLeaf, i.left, i.right = value, true, nil, nil
	return i
}

func (i *id) asNode(left, right int) *id {
	return i.asNodeWithIds(newIDWithValue(left), newIDWithValue(right))
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

// Split an ID as defined in section "5.3.2 Fork"
func (i *id) Split() (i1, i2 *id) {

	i1 = newID()
	i2 = newID()

	if i.isLeaf && i.value == 0 {
		// split(0) = (0, 0)
		i1.value = 0
		i2.value = 0
		return
	}
	if i.isLeaf && i.value == 1 {
		// split(1) = ((1,0), (0,1))
		i1.asNode(1, 0)
		i2.asNode(0, 1)
		return
	}
	if (i.left.isLeaf && i.left.value == 0) && (!i.right.isLeaf || i.right.value == 1) {
		// split((0, i)) = ((0, i1), (0, i2)), where (i1, i2) = split(i)
		r1, r2 := i.right.Split()
		i1.asNodeWithIds(newIDWithValue(0), r1)
		i2.asNodeWithIds(newIDWithValue(0), r2)
		return
	}
	if (!i.left.isLeaf || i.left.value == 1) && (i.right.isLeaf && i.right.value == 0) {
		// split((i, 0)) = ((i1, 0), (i2, 0)), where (i1, i2) = split(i)
		l1, l2 := i.left.Split()
		i1.asNodeWithIds(l1, newIDWithValue(0))
		i2.asNodeWithIds(l2, newIDWithValue(0))
		return
	}
	if (!i.left.isLeaf || i.left.value == 1) && (!i.right.isLeaf || i.right.value == 1) {
		// split((i1, i2)) = ((i1, 0), (0, i2))
		i1.asNodeWithIds(i.left, newIDWithValue(0))
		i2.asNodeWithIds(newIDWithValue(0), i.right)
		return
	}
	log.Fatalf("unable to split id with unexpected setup: %s", i.String())
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
	return i.asNodeWithIds(newID().sum(i1.left, i2.left), newID().sum(i1.right, i2.right)).Norm()
}
