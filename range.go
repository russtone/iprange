package iprange

import (
	"bytes"
	"math/big"
	"net"
)

// Range is an interface to work IP-addresses range.
type Range interface {
	// Contains checks if the given IP-address is in the range.
	Contains(net.IP) bool

	// Count returns number of addresses in the range.
	Count() *big.Int

	// Iterator returns iterator for the range.
	Iterator() Iterator
}

//
// singleRange
//

type singleRange struct {
	net.IP
}

var _ Range = singleRange{}

func (r singleRange) Contains(ip net.IP) bool {
	return r.Equal(ip)
}

func (r singleRange) Count() *big.Int {
	return big.NewInt(1)
}

func (r singleRange) Iterator() Iterator {
	return &singleIterator{r, false}
}

//
// min max
//

type minMaxRange struct {
	min, max net.IP
}

var _ Range = minMaxRange{}

func (r minMaxRange) Contains(ip net.IP) bool {
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}

	min := bytes.Compare(ip, r.min)
	max := bytes.Compare(ip, r.max)
	return min == 0 || max == 0 || (min == 1 && max == -1)
}

func (r minMaxRange) Count() *big.Int {
	c := ip2big(r.max)
	c.Sub(c, ip2big(r.min))
	c.Add(c, big.NewInt(1))
	return c
}

func (r minMaxRange) Iterator() Iterator {
	return &minMaxIterator{r, ip2big(r.min), ip2big(r.max)}
}

//
// octetsRange
//

type octetsRange struct {
	octets ipOctets
}

var _ Range = &octetsRange{}

func (r octetsRange) Contains(ip net.IP) bool {
	var contains bool
	for i, oct := range ip2octets(ip) {
		contains = false
		for _, b := range r.octets[i] {
			if oct >= b.lo && oct <= b.hi {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

func (r octetsRange) Count() *big.Int {
	c := big.NewInt(1)

	for i := len(r.octets) - 1; i >= 0; i-- {
		s := big.NewInt(0)
		for _, b := range r.octets[i] {
			s.Add(s, big.NewInt(int64(b.hi - b.lo + 1)))
		}
		c.Mul(c, s)
	}

	return c
}

func (r octetsRange) Iterator() Iterator {
	indexes := make([]int, len(r.octets))
	min := r.octets.min()
	return &octetsIterator{r, false,indexes, min}
}

//
// Ranges
//

// Ranges allows to combine multiple ranges and use them as one.
type Ranges []Range

var _ Range = Ranges{}

// Contains allows Ranges to satisfy Range interface.
func (rr Ranges) Contains(ip net.IP) bool {
	for _, r := range rr {
		if r.Contains(ip) {
			return true
		}
	}
	return false
}

// Count allows Ranges to satisfy Range interface.
func (rr Ranges) Count() *big.Int {
	c := big.NewInt(0)
	for _, r := range rr {
		c.Add(c, r.Count())
	}
	return c
}

// Iterator allows Ranges to satisfy Range interface.
func (rr Ranges) Iterator() Iterator {
	its := make([]Iterator, len(rr))
	for i := 0; i < len(rr); i++ {
		its[i] = rr[i].Iterator()
	}

	return &rangesIterator{its, 0}
}