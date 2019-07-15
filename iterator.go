package iprange

import (
	"net"
)

type Iterator interface {
	Next(net.IP) bool
}

type iterator struct {
	iter iterable
	ip   net.IP
}

type multiIterator struct {
	rngs []IPRange
	pos  int
	ip   net.IP
	cur  Iterator
}

func NewIterator(r IPRange) Iterator {
	if rr, ok := r.(IPRanges); ok {
		return &multiIterator{rr, 0, nil, nil}
	}

	// All ranges types except IPRanges satisfy iterable.
	it, _ := r.(iterable)

	return &iterator{it, nil}
}

func (it *iterator) Next(out net.IP) bool {
	if out.To16() != nil {
		out = out.To4()
	}

	it.ip = it.iter.next(it.ip)

	if it.ip == nil {
		return false
	}

	copy(out, it.ip)
	return true
}

func (it *multiIterator) Next(out net.IP) bool {

	for {
		if it.cur != nil {
			if it.cur.Next(out) {
				return true
			}
			it.cur = nil
		}

		if it.pos < len(it.rngs) {
			r := it.rngs[it.pos]
			it.cur = NewIterator(r)
			it.pos++
		} else {
			break
		}
	}

	return false
}
