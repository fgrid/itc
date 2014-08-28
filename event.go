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

func (e *event) dec(unpacker *BitUnPacker) *event {
	if unpacker.Pop(1) == 1 {
		return e.asLeaf(int(dec(2, unpacker)))
	}
	e.asNode(0, 0, 0)
	switch unpacker.Pop(2) {
	case 0:
		e.right.dec(unpacker)
	case 1:
		e.left.dec(unpacker)
	case 2:
		e.left.dec(unpacker)
		e.right.dec(unpacker)
	case 3:
		if unpacker.Pop(1) == 0 {
			if unpacker.Pop(1) == 0 {
				unpacker.Pop(1)
				e.value = int(dec(2, unpacker))
				e.right.dec(unpacker)
			} else {
				unpacker.Pop(1)
				e.value = int(dec(2, unpacker))
				e.left.dec(unpacker)
			}
		} else {
			unpacker.Pop(1)
			e.value = int(dec(2, unpacker))
			e.left.dec(unpacker)
			e.right.dec(unpacker)
		}
	}
	return e
}

func (e *event) enc(packer *BitPacker) *BitPacker {
	if e.isLeaf {
		packer.Push(1, 1)
		return enc(uint32(e.value), 2, packer)
	}

	packer.Push(0, 1)
	if e.value == 0 {
		if e.left.isLeaf && e.left.value == 0 {
			packer.Push(0, 2)
			return e.right.enc(packer)
		}
		if e.right.isLeaf && e.right.value == 0 {
			packer.Push(1, 2)
			return e.left.enc(packer)
		}
		packer.Push(2, 2)
		e.left.enc(packer)
		return e.right.enc(packer)
	}

	packer.Push(3, 2)
	if e.left.isLeaf && e.left.value == 0 {
		packer.Push(0, 1)
		packer.Push(0, 1)
		newEventWithValue(e.value).enc(packer)
		return e.right.enc(packer)
	}
	if e.right.isLeaf && e.right.value == 0 {
		packer.Push(0, 1)
		packer.Push(1, 1)
		newEventWithValue(e.value).enc(packer)
		return e.left.enc(packer)
	}
	packer.Push(1, 1)
	newEventWithValue(e.value).enc(packer)
	e.left.enc(packer)
	return e.right.enc(packer)
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
func dec(B uint32, unpacker *BitUnPacker) uint32 {
	max := uint32(1) << uint32(B)
	if unpacker.Pop(1) == 0 {
		return unpacker.Pop(B)
	}
	return max + dec(B+1, unpacker)
}

func enc(n, B uint32, packer *BitPacker) *BitPacker {
	max := uint32(1) << uint32(B)
	if n < max {
		packer.Push(0, 1)
		packer.Push(uint32(n), uint32(B))
	} else {
		packer.Push(1, 1)
		enc(n-max, B+1, packer)
	}
	return packer
}

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
