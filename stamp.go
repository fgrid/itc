// Package itc implements the interval tree clock as described in the paper
// 'Interval Tree Clocks: A Logical Clock for Dynamic Systems' by Paulo Sergio Almeida,
// Carlos Baquero and Victor Fonte. (http://gsd.di.uminho.pt/members/cbm/ps/itc2008.pdf)
//
// Causality tracking mechanisms can be modeled by a set of core operations: fork; event and join, that
// act on stamps (logical clocks) whose structure is a pair (i, e), formed by an id and an event component
// that encodes causally known events.
package itc

import (
	"fmt"
	"github.com/fgrid/itc/bit"
	"github.com/fgrid/itc/event"
	"github.com/fgrid/itc/id"
)

// Stamp declares the state of the clock for a given identity and a given stream of events.
type Stamp struct {
	event *event.Event
	id    *id.ID
}

// NewStamp creates a new so called seed-Stamp (represented as: (1, 0)).
func NewStamp() *Stamp {
	return &Stamp{event: event.New(), id: id.New()}
}

// Event adds a new event to the clock's event component, so that if (i, e') results from event((i, e))
// the causal ordering is such that e < e'.
func (s *Stamp) Event() {
	oldE := s.event.Clone()
	newE := s.fill()
	if newE.Equals(oldE) {
		s.event, _ = s.grow()
	} else {
		s.event = newE
	}
}

func (s *Stamp) fill() *event.Event {
	return fill(s.id, s.event)
}

// Fork clones the causal past of a stamp, resulting in a pair of stamps that
// have identical copies of the event component and distinct IDs.
func (s *Stamp) Fork() *Stamp {
	st := NewStamp()
	id1, id2 := s.id.Split()
	s.id = id1
	st.id = id2
	st.event = s.event.Clone()
	return st
}

func (s *Stamp) grow() (*event.Event, int) {
	return grow(s.id, s.event)
}

// Join merges two stamps, producing a new one.
func (s *Stamp) Join(other *Stamp) {
	s.id = id.New().Sum(s.id, other.id)
	s.event = event.Join(s.event, other.event)
}

// LEQ Compares the stamp with the given other stamp and returns 'true' if this stamp is less or equal (LEQ).
func (s *Stamp) LEQ(other *Stamp) bool {
	return event.LEQ(s.event, other.event)
}

// MarshalBinary encodes the stamp s into a binary form and returns the result.
func (s *Stamp) MarshalBinary() ([]byte, error) {
	bp := bit.NewPack()
	s.Pack(bp)
	return bp.Pack(), nil
}

// String returns the string corresponding to stamp s.
func (s *Stamp) String() string {
	return fmt.Sprintf("(%s, %s)", s.id, s.event)
}

// UnmarshalBinary decodes the stamp s from the given binary form data (created by MarshalBinary).
func (s *Stamp) UnmarshalBinary(data []byte) error {
	bup := bit.NewUnPack(data)
	s.UnPack(bup)
	return nil
}

func fill(i *id.ID, e *event.Event) *event.Event {
	if i.IsLeaf {
		if i.Value == 0 {
			return e
		}
		return event.NewLeaf(e.Max())
	}
	if e.IsLeaf {
		return e
	}
	r := event.NewEmptyNode(e.Value)
	if i.Left.IsLeaf && i.Left.Value == 1 {
		r.Right = fill(i.Right, e.Right)
		r.Left = event.NewLeaf(event.Max(e.Left.Max(), r.Right.Min()))
	} else if i.Right.IsLeaf && i.Right.Value == 1 {
		r.Left = fill(i.Left, e.Left)
		r.Right = event.NewLeaf(event.Max(e.Right.Max(), r.Left.Min()))
	} else {
		r.Left = fill(i.Left, e.Left)
		r.Right = fill(i.Right, e.Right)
	}
	return r.Norm()
}

func grow(i *id.ID, e *event.Event) (*event.Event, int) {
	if e.IsLeaf {
		if i.IsLeaf && i.Value == 1 {
			return event.NewLeaf(e.Value + 1), 0
		}
		ex, c := grow(i, event.NewEmptyNode(e.Value))
		return ex, c + 99999
	}
	if i.Left.IsLeaf && i.Left.Value == 0 {
		exr, cr := grow(i.Right, e.Right)
		r := event.NewEmptyNode(e.Value)
		r.Left = e.Left
		r.Right = exr
		return r, cr + 1
	}
	if i.Right.IsLeaf && i.Right.Value == 0 {
		exl, cl := grow(i.Left, e.Left)
		r := event.NewEmptyNode(e.Value)
		r.Left = exl
		r.Right = e.Right
		return r, cl + 1
	}
	exl, cl := grow(i.Left, e.Left)
	exr, cr := grow(i.Right, e.Right)
	if cl < cr {
		r := event.NewEmptyNode(e.Value)
		r.Left = exl
		r.Right = e.Right
		return r, cl + 1
	}
	r := event.NewEmptyNode(e.Value)
	r.Left = e.Left
	r.Right = exr
	return r, cr + 1
}

func (s *Stamp) Pack(p *bit.Pack) {
	s.id.Pack(p)
	s.event.Pack(p)
}

func (s *Stamp) UnPack(bup *bit.UnPack) {
	s.id = id.UnPack(bup)
	s.event = event.UnPack(bup)
}
