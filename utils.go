package iprange

import (
	"math/big"
	"net"
)

// net.IP to *big.Int
func ip2big(ip net.IP) *big.Int {
	return big.NewInt(0).SetBytes(ip)
}

// octcmp compares two slices of IP octets.
// Returns 1 if a > b, -1 if a < b and 0 if a = b.
func octcmp(a, b []uint16) int {
	for i := 0; i < len(a); i++ {
		if a[i] > b[i] {
			return 1
		} else if b[i] > a[i] {
			return -1
		}
	}
	return 0
}

// net.IP to octets as []uint16.
func ip2octets(ip net.IP) []uint16 {
	octs := make([]uint16, 0)

	if ip4 := ip.To4(); ip4 != nil {
		for _, b := range ip4 {
			octs = append(octs, uint16(b))
		}
		return octs
	}

	for i := 0; i < net.IPv6len/2; i++ {
		octs = append(octs, uint16(ip[2*i])<<8|uint16(ip[2*i+1]))
	}

	return octs
}

// []uint16 octets to net.IP.
func octets2ip(octs []uint16) net.IP {
	iplen := net.IPv4len
	if len(octs) != net.IPv4len {
		iplen = net.IPv6len
	}

	ip := make(net.IP, iplen)

	for i, oct := range octs {
		if iplen == net.IPv4len {
			ip[i] = byte(oct)
		} else {
			ip[2*i] = byte(oct >> 8)
			ip[2*i+1] = byte(oct)
		}
	}

	return ip
}

// Bigger than we need, not too big to worry about overflow
const toobig = 0xFFFFFF

// Decimal to integer.
// Returns number, characters consumed, success.
func dtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i]-'0')
		if n >= toobig {
			return toobig, i, false
		}
	}
	if i == 0 {
		return 0, 0, false
	}
	return n, i, true
}

// Hexadecimal to integer.
// Returns number, characters consumed, success.
func xtoi(s string) (n int, i int, ok bool) {
	n = 0
	for i = 0; i < len(s); i++ {
		if '0' <= s[i] && s[i] <= '9' {
			n *= 16
			n += int(s[i] - '0')
		} else if 'a' <= s[i] && s[i] <= 'f' {
			n *= 16
			n += int(s[i]-'a') + 10
		} else if 'A' <= s[i] && s[i] <= 'F' {
			n *= 16
			n += int(s[i]-'A') + 10
		} else {
			break
		}
		if n >= toobig {
			return 0, i, false
		}
	}
	if i == 0 {
		return 0, i, false
	}
	return n, i, true
}
