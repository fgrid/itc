package itc

import "fmt"

type Stamp struct {
	event *event
	id    *id
}

func NewStamp() *Stamp {
	return &Stamp{event: newEvent(), id: newId()}
}

func (s *Stamp) dec(unpacker *BitUnPacker) *Stamp {
	s.id.dec(unpacker)
	s.event.dec(unpacker)
	return s
}

func (s *Stamp) enc(packer *BitPacker) *BitPacker {
	s.id.enc(packer)
	s.event.enc(packer)
	return packer
}

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

func (s *Stamp) Join(s2 *Stamp) {
	s.id = newId().sum(s.id, s2.id)
	s.event = newEvent().join(s.event, s2.event)
}

func (s *Stamp) Leq(s2 *Stamp) bool {
	return leq(s.event, s2.event)
}

func (s *Stamp) Pack() []byte {
	return s.enc(NewBitPacker()).Pack()
}

func (s *Stamp) String() string {
	return fmt.Sprintf("(%s, %s)", s.id, s.event)
}

func (s *Stamp) UnPack(packed []byte) *Stamp {
	return s.dec(NewBitUnPacker(packed))
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
	} else {
		r := newEvent().asNode(e.value, 0, 0)
		r.left = e.left
		r.right = exr
		return r, cr + 1
	}
}
