package iprange

import (
	"net"
)

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func incEx(ip, lower, upper net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++

		if ip[i] <= upper[i] {
			break
		} else {
			ip[i] = lower[i]
			continue
		}
	}
}
