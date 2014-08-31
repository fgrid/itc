package itc

import "fmt"

type event struct {
	value       int
	left, right *event
	isLeaf      bool
}

func newEvent() *event {
	return newEventWithValue(0)
}

func newEventWithValue(value int) *event {
	return &event{value: value, isLeaf: true}
}

func (e *event) asLeaf(value int) *event {
	e.value, e.isLeaf, e.left, e.right = value, true, nil, nil
	return e
}

func (e *event) asNode(value, left, right int) *event {
	e.value, e.isLeaf, e.left, e.right = value, false, newEventWithValue(left), newEventWithValue(right)
	return e
}

func (e *event) clone() *event {
	result := newEvent()
	result.isLeaf = e.isLeaf
	result.value = e.value
	if e.left != nil {
		result.left = e.left.clone()
	}
	if e.right != nil {
		result.right = e.right.clone()
	}
	return result
}

func (e *event) equals(o *event) bool {
	return (e == nil && o == nil) ||
		((e.isLeaf == o.isLeaf) &&
			(e.value == o.value) &&
			e.left.equals(o.left) &&
			e.right.equals(o.right))
}

func (e *event) join(e1, e2 *event) *event {
	if e1.isLeaf && e2.isLeaf {
		return e.asLeaf(max(e1.value, e2.value))
	}
	if e1.isLeaf {
		return e.join(newEvent().asNode(e1.value, 0, 0), e2)
	}
	if e2.isLeaf {
		return e.join(e1, newEvent().asNode(e2.value, 0, 0))
	}
	if e1.value > e2.value {
		return e.join(e2, e1)
	}
	e.isLeaf = false
	e.value = e1.value
	e.left = newEvent().join(e1.left, e2.left.lift(e2.value-e1.value))
	e.right = newEvent().join(e1.right, e2.right.lift(e2.value-e1.value))
	return e.norm()
}

func (e *event) norm() *event {
	if e.isLeaf {
		return e
	}
	if e.left.isLeaf && e.right.isLeaf && e.left.value == e.right.value {
		return e.asLeaf(e.value + e.left.value)
	}
	m := min(e.left.min(), e.right.min())
	e.left = e.left.norm().sink(m)
	e.right = e.right.norm().sink(m)
	return e.lift(m)
}

func (e *event) lift(value int) *event {
	result := e.clone()
	result.value += value
	return result
}

func (e *event) max() int {
	if e.isLeaf {
		return e.value
	}
	return e.value + max(e.left.max(), e.right.max())
}

func (e *event) min() int {
	if e.isLeaf {
		return e.value
	}
	return e.value + min(e.left.min(), e.right.min())
}

func (e *event) sink(value int) *event {
	return e.lift(-1 * value)
}

func (e *event) String() string {
	if e.isLeaf {
		return fmt.Sprintf("%d", e.value)
	}
	return fmt.Sprintf("(%d, %s, %s)", e.value, e.left, e.right)
}

// ----------
func leq(e1, e2 *event) bool {
	if e1.isLeaf {
		return e1.value <= e2.value
	}
	if e2.isLeaf {
		return (e1.value <= e2.value) &&
			leq(e1.left.lift(e1.value), e2) &&
			leq(e1.right.lift(e1.value), e2)
	}
	return (e1.value <= e2.value) &&
		leq(e1.left.lift(e1.value), e2.left.lift(e2.value)) &&
		leq(e1.right.lift(e1.value), e2.right.lift(e2.value))
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}
