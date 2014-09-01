// Package itc implements the interval tree clock as described in the paper
// 'Interval Tree Clocks: A Logical Clock for Dynamic Systems' by Paulo Sergio Almeida,
// Carlos Baquero and Victor Fonte. (http://gsd.di.uminho.pt/members/cbm/ps/itc2008.pdf)
//
// Causality tracking mechanisms can be modeled by a set of core operations: fork; event and join, that
// act on stamps (logical clocks) whose structure is a pair (i, e), formed by an id and an event component
// that encodes causally known events.
package itc

import "fmt"

// Stamp declares the state of the clock for a given identity and a given stream of events.
type Stamp struct {
	event *event
	id    *id
}

// NewStamp creates a new so called seed-Stamp (represented as: (1, 0)).
func NewStamp() *Stamp {
	return &Stamp{event: newEvent(), id: newID()}
}

// Event adds a new event to the clock's event component, so that if (i, e') results from event((i, e))
// the causal ordering is such that e < e'.
func (s *Stamp) Event() {
	oldE := s.event.clone()
	newE := s.fill()
	if newE.equals(oldE) {
		s.event, _ = s.grow()
	} else {
		s.event = newE
	}
}

func (s *Stamp) fill() *event {
	return fill(s.id, s.event)
}

// Fork clones the causal past of a stamp, resulting in a pair of stamps that
// have identical copies of the event component and distinct IDs.
func (s *Stamp) Fork() *Stamp {
	st := NewStamp()
	id1, id2 := s.id.Split()
	s.id = id1
	st.id = id2
	st.event = s.event.clone()
	return st
}

func (s *Stamp) grow() (*event, int) {
	return grow(s.id, s.event)
}

// Join merges two stamps, producing a new one.
func (s *Stamp) Join(other *Stamp) {
	s.id = newID().sum(s.id, other.id)
	s.event = newEvent().join(s.event, other.event)
}

// Leq Compares the stamp with the given other stamp and returns 'true' if this stamp is less or equal (leq).
func (s *Stamp) Leq(other *Stamp) bool {
	return leq(s.event, other.event)
}

// MarshalBinary encodes the stamp s into a binary form and returns the result.
func (s *Stamp) MarshalBinary() ([]byte, error) {
	bp := newBitPack()
	bp.encodeStamp(s)
	return bp.Pack(), nil
}

// String returns the string corresponding to stamp s.
func (s *Stamp) String() string {
	return fmt.Sprintf("(%s, %s)", s.id, s.event)
}

// UnmarshalBinary decodes the stamp s from the given binary form data (created by MarshalBinary).
func (s *Stamp) UnmarshalBinary(data []byte) error {
	bup := newBitUnPack(data)
	bup.decodeStamp(s)
	return nil
}

func fill(i *id, e *event) *event {
	if i.isLeaf {
		if i.value == 0 {
			return e
		}
		return newEventWithValue(e.max())
	}
	if e.isLeaf {
		return e
	}
	r := newEvent().asNode(e.value, 0, 0)
	if i.left.isLeaf && i.left.value == 1 {
		r.right = fill(i.right, e.right)
		r.left = newEventWithValue(max(e.left.max(), r.right.min()))
	} else if i.right.isLeaf && i.right.value == 1 {
		r.left = fill(i.left, e.left)
		r.right = newEventWithValue(max(e.right.max(), r.left.min()))
	} else {
		r.left = fill(i.left, e.left)
		r.right = fill(i.right, e.right)
	}
	return r.norm()
}

func grow(i *id, e *event) (*event, int) {
	if e.isLeaf {
		if i.isLeaf && i.value == 1 {
			return newEventWithValue(e.value + 1), 0
		}
		ex, c := grow(i, newEvent().asNode(e.value, 0, 0))
		return ex, c + 99999
	}
	if i.left.isLeaf && i.left.value == 0 {
		exr, cr := grow(i.right, e.right)
		r := newEvent().asNode(e.value, 0, 0)
		r.left = e.left
		r.right = exr
		return r, cr + 1
	}
	if i.right.isLeaf && i.right.value == 0 {
		exl, cl := grow(i.left, e.left)
		r := newEvent().asNode(e.value, 0, 0)
		r.left = exl
		r.right = e.right
		return r, cl + 1
	}
	exl, cl := grow(i.left, e.left)
	exr, cr := grow(i.right, e.right)
	if cl < cr {
		r := newEvent().asNode(e.value, 0, 0)
		r.left = exl
		r.right = e.right
		return r, cl + 1
	}
	r := newEvent().asNode(e.value, 0, 0)
	r.left = e.left
	r.right = exr
	return r, cr + 1
}
