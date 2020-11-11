package iprange

import (
	"sort"
)

// IP octet bounds.
// uint16 is used to be able to use with IPv6.
type ipOctet struct {
	lo, hi uint16
}

// ipOctets represents IP address as list of bounds for every octet.
// singleRange IP 192.168.1.1 will be represented as [192 192], [168 168], [1 1], [1 1].
// IP with octets ranges 192.168.1,2.1-10 as [192 192], [168 168], [[1 1] [2 2]], [[1 10]]
type ipOctets [][]ipOctet

// min returns the lowest IP-address that satisfies the octet boundaries.
func (octs ipOctets) min() []uint16 {
	ip := make([]uint16, len(octs))
	for i, oct := range octs {
		ip[i] = oct[0].lo
	}
	return ip
}

// max returns the biggest IP-address that satisfies the octet boundaries
func (octs ipOctets) max() []uint16 {
	ip := make([]uint16, len(octs))
	for i, oct := range octs {
		ip[i] = oct[len(oct)-1].hi
	}
	return ip
}

// adds bound for i-th octet.
func (octs ipOctets) push(i int, lo, hi uint16) {
	if hi == 0 {
		hi = lo
	} else if lo > hi {
		lo, hi = hi, lo
	}

	octs[i] = append(octs[i], ipOctet{lo, hi})
}

// returns true if IP has octet ranges.
func (octs ipOctets) hasRanges() bool {
	for _, oct := range octs {
		if len(oct) > 1 || (len(oct) == 1 && oct[0].lo != oct[0].hi) {
			return true
		}
	}
	return false
}

// sorts octets bounds.
func (octs ipOctets) sort() {
	for _, oct := range octs {
		sort.SliceStable(oct, func(i, j int) bool {
			return oct[i].lo < oct[j].lo
		})
	}
}