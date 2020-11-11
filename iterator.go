package iprange

import (
	"math/big"
	"net"
)

// Iterator is an interface to iterate IP addresses
// from one of multiple ranges.
type Iterator interface {
	// Next returns true if there is at least one IP-address left in the iterator
	// and saves this address into the given pointer. If no addresses left returns false.
	Next(*net.IP) bool

	// Reset resets the iterator so it can be used again.
	Reset()

	// Count returns total number of IP-addresses in iterator.
	Count() *big.Int

	// Contains checks if the given IP is in one of the ranges of the iterator.
	Contains(net.IP) bool
}

//
// Single IP iterator.
//

type singleIterator struct {
	singleRange
	done bool
}

var _ Iterator = &singleIterator{}

func (it *singleIterator) Next(out *net.IP) bool {
	if it.done {
		return false
	}

	ip := make(net.IP, len(it.IP))
	copy(ip, it.IP)
	*out = ip

	it.done = true

	return true
}

func (it *singleIterator) Reset() {
	it.done = false
}

//
// Min max range iterator.
//

type minMaxIterator struct {
	minMaxRange
	current *big.Int
	last    *big.Int
}

var _ Iterator = &minMaxIterator{}

func (it *minMaxIterator) Next(out *net.IP) bool {
	if it.current.Cmp(it.last) == 1 {
		return false
	}

	ip := make(net.IP, len(it.min))
	copy(ip, it.current.Bytes())
	*out = ip

	it.current.Add(it.current, big.NewInt(1))

	return true
}

func (it *minMaxIterator) Reset() {
	it.current.SetBytes(it.min)
}

//
// Octets range iterator.
//

type octetsIterator struct {
	octetsRange
	done    bool
	indexes []int
	current []uint16
}

var _ Iterator = &octetsIterator{}

func (it *octetsIterator) Next(out *net.IP) bool {
	if it.done {
		return false
	}

	ip := octets2ip(it.current)
	*out = ip

	if octcmp(it.current, it.octets.max()) == 0 {
		it.done = true
		return true
	}

	var (
		oct uint16
		j   int
	)

	for i := len(it.octets) - 1; i >= 0; i-- {
		oct = it.current[i]
		oct++

		j = it.indexes[i]

		if oct >= it.octets[i][j].lo && oct <= it.octets[i][j].hi {
			it.current[i] = oct
			break
		} else if j+1 < len(it.octets[i]) {
			it.current[i] = it.octets[i][j+1].lo
			it.indexes[i] = j + 1
			break
		}

		it.current[i] = it.octets[i][0].lo
		it.indexes[i] = 0
	}

	return true
}

func (it *octetsIterator) Reset() {
	it.done = false
	for i := 0; i < len(it.octets); i++ {
		it.indexes[i] = 0
		it.current[i] = it.octets[i][0].lo
	}
}

//
// Multiple ranges iterator.
//

type rangesIterator struct {
	its []Iterator
	idx int
}

var _ Iterator = &rangesIterator{}

func (it *rangesIterator) Next(out *net.IP) bool {
	for i := it.idx; i < len(it.its); i++ {
		if it.its[i].Next(out) {
			return true
		}

		it.idx++
	}

	return false
}

func (it *rangesIterator) Reset() {
	it.idx = 0
	for _, i := range it.its {
		i.Reset()
	}
}

func (it *rangesIterator) Count() *big.Int {
	c := big.NewInt(0)

	for _, r := range it.its {
		c.Add(c, r.Count())
	}

	return c
}

func (it *rangesIterator) Contains(ip net.IP) bool {
	for _, i := range it.its {
		if i.Contains(ip) {
			return true
		}
	}
	return false
}
