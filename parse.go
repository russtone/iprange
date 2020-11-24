package iprange

import (
	"net"
)

type parseFunc func(string) (ipOctets, int)

// Parse parses s as an IP addresses range (IPv4 or IPv6), returning the result.
// The string s can be in the following formats:
// single IP ("192.0.2.1", "2001:db8::68"), CIDR range ("192.168.1.0/24", "2001:db8::68/120"),
// begin_end range ("192.168.1.1_192.168.1.10", "2001:db8::68_2001:db8::80") or
// octets range ("192.168.1,3,5.1-10", "2001:db8::0,1:68-80").
// If s is not a valid textual representation of an IP addresses range,
// ParseIP returns nil.
func Parse(s string) Range {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return parse(s, parseIPv4, net.IPv4len)
		case ':':
			return parse(s, parseIPv6, net.IPv6len)
		}
	}
	return nil
}

// parse decides which kind of range is s: singleRange, cidrRange, minMaxRange or octetsRange.
func parse(s string, parseFn parseFunc, iplen int) Range {
	ip, c := parseFn(s)

	if ip == nil {
		return nil
	}

	s = s[c:]

	if len(s) > 0 && s[0] == '_' {
		// begin_end range.

		if ip.hasRanges() {
			// Already have octet ranges.
			return nil
		}

		s = s[1:]

		max, c := parseFn(s)

		if max == nil || max.hasRanges() || len(ip) != len(max) {
			// Invalid ip or have octet ranges.
			return nil
		}

		s = s[c:]

		// Must have used entire string.
		if len(s) != 0 {
			return nil
		}

		return &minMaxRange{octets2ip(ip.min()), octets2ip(max.min())}
	}

	if len(s) > 0 && s[0] == '/' {
		// CIDR range.

		if ip.hasRanges() {
			// Already have octet ranges.
			return nil
		}

		s = s[1:]

		// Decimal mask.
		n, c, ok := dtoi(s)
		if !ok || n < 0 || n > 8*iplen {
			return nil
		}

		s = s[c:]

		// Must have used entire string.
		if len(s) != 0 {
			return nil
		}

		mask := net.CIDRMask(n, iplen*8)
		min := octets2ip(ip.min()).Mask(mask)
		max := make(net.IP, len(min))

		for i, m := range mask {
			if m == 0xff {
				max[i] = min[i]
			} else {
				max[i] = min[i] | (m ^ 0xff)
			}
		}

		return &minMaxRange{min, max}
	}

	// Must have used entire string.
	if len(s) != 0 {
		return nil
	}

	if !ip.hasRanges() {
		// singleRange ip.
		return &singleRange{octets2ip(ip.min())}
	}

	// Sort ip2octets.
	ip.sort()

	return &octetsRange{ip}
}

// parseIPv4 parses s as IPv4, based on net.parseIPv4.
// Returns IP octets, characters consumed.
func parseIPv4(s string) (ip ipOctets, cc int) {
	ip = make(ipOctets, net.IPv4len)

	for i := 0; i < net.IPv4len; i++ {
		ip[i] = make([]ipOctet, 0)
	}

	var bb [2]uint16 // octet bounds: 0 - lo, 1 - hi

	i := 0 // octet idx
	k := 0 // bound idx: 0 - lo, 1 - hi

loop:
	for i < net.IPv4len {
		// Decimal number.
		n, c, ok := dtoi(s)
		if !ok || n > 0xFF {
			return nil, cc
		}

		// Save bound.
		bb[k] = uint16(n)

		// Stop at max of string.
		s = s[c:]
		cc += c
		if len(s) == 0 {
			ip.push(i, bb[0], bb[1])
			i++
			break
		}

		// Otherwise must be followed by dot, colon or dp.
		switch s[0] {
		case '.':
			fallthrough
		case ',':
			ip.push(i, bb[0], bb[1])
			bb[1] = 0
			k = 0
		case '-':
			if k == 1 {
				// To many dashes in one octet.
				return nil, cc
			}
			k++
		default:
			ip.push(i, bb[0], bb[1])
			i++
			break loop
		}

		if s[0] == '.' {
			i++
		}

		s = s[1:]
		cc++
	}

	if i < net.IPv4len {
		// Missing ip2octets.
		return nil, cc
	}

	return ip, cc
}

// parseIPv6 parses s as IPv6, based on net.parseIPv6.
// Returns IP octets, characters consumed.
func parseIPv6(s string) (ip ipOctets, cc int) {
	ip = make(ipOctets, net.IPv6len/2)

	for i := 0; i < net.IPv6len/2; i++ {
		ip[i] = make([]ipOctet, 0)
	}

	ellipsis := -1 // position of ellipsis in ip

	// Might have leading ellipsis
	if len(s) >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		s = s[2:]
		cc += 2

		// Might be only ellipsis
		if len(s) == 0 || s[0] == '_' || s[0] == '/' {
			for i := 0; i < net.IPv6len/2; i++ {
				ip.push(i, 0, 0)
			}
			return ip, cc
		}
	}

	var bb [2]uint16 // octet bounds
	i := 0           // octet idx
	k := 0           // bound idx: 0 - lo, 1 - hi

	// Loop, parsing hex numbers followed by colon.
loop:
	for i < net.IPv6len/2 {
		// Hex number.
		n, c, ok := xtoi(s)
		if !ok || n > 0xFFFF {
			return nil, cc
		}

		// If followed by dot, might be in trailing net.IPv4.
		if c < len(s) && s[c] == '.' {
			ip, n := parseIPv4(s)

			return ip, cc + n
		}

		// Save this 16-bit chunk.
		bb[k] = uint16(n)

		// Stop at max of string.
		s = s[c:]
		cc += c
		if len(s) == 0 {
			ip.push(i, bb[0], bb[1])
			i++
			break
		}

		switch s[0] {
		case ':':
			fallthrough
		case ',':
			ip.push(i, bb[0], bb[1])
			bb[1] = 0
			k = 0
		case '-':
			if k == 1 {
				// To many dashes in one octet.
				return nil, cc
			}
			k++
		default:
			ip.push(i, bb[0], bb[1])
			break loop
		}

		if s[0] == ':' {
			i++
		}

		s = s[1:]
		cc++

		// Look for ellipsis.
		if s[0] == ':' {
			if ellipsis >= 0 { // already have one
				return nil, cc
			}
			ellipsis = i
			s = s[1:]
			cc++
			if len(s) == 0 || s[0] == '_' || s[0] == '/' { // can be at end
				break
			}
		}
	}

	// If didn't parse enough, expand ellipsis.
	if i < net.IPv6len/2 {
		if ellipsis < 0 {
			return nil, cc
		}
		n := net.IPv6len/2 - i
		for j := i - 1; j >= ellipsis; j-- {
			ip[j+n] = ip[j]
		}
		for j := ellipsis + n - 1; j >= ellipsis; j-- {
			ip[j] = []ipOctet{{0, 0}}
		}
	} else if ellipsis >= 0 {
		// Ellipsis must represent at least one 0 group.
		return nil, cc
	}

	return ip, cc
}
