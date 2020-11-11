package iprange_test

import (
	"fmt"
	"github.com/russtone/iprange"
	"net"
)

func ExampleRange_single() {
	r := iprange.Parse("192.168.1.1")

	fmt.Println(r.Count())
	fmt.Println(r.Contains(net.ParseIP("192.168.1.1")))
	fmt.Println(r.Contains(net.ParseIP("192.168.10.1")))

	// Output:
	// 1
	// true
	// false
}

func ExampleRange_cidr() {
	r := iprange.Parse("192.168.1.10/24")

	fmt.Println(r.Count())
	fmt.Println(r.Contains(net.ParseIP("192.168.1.50")))
	fmt.Println(r.Contains(net.ParseIP("192.168.10.1")))

	// Output:
	// 256
	// true
	// false
}

func ExampleRange_beginend() {
	r := iprange.Parse("192.168.1.10\\192.168.2.9")

	fmt.Println(r.Count())
	fmt.Println(r.Contains(net.ParseIP("192.168.2.1")))
	fmt.Println(r.Contains(net.ParseIP("192.168.10.1")))

	// Output:
	// 256
	// true
	// false
}

func ExampleRange_octets() {
	r := iprange.Parse("192.168.1-2,4.1-10")

	fmt.Println(r.Count())
	fmt.Println(r.Contains(net.ParseIP("192.168.2.5")))
	fmt.Println(r.Contains(net.ParseIP("192.168.1.15")))

	// Output:
	// 30
	// true
	// false
}

func ExampleIterator() {
	r := iprange.Parse("192.168.1.0/29")
	if r == nil {
		return
	}
	it := r.Iterator()
	var ip net.IP
	for it.Next(&ip) {
		fmt.Println(ip.String())
	}

	// Output:
	// 192.168.1.0
	// 192.168.1.1
	// 192.168.1.2
	// 192.168.1.3
	// 192.168.1.4
	// 192.168.1.5
	// 192.168.1.6
	// 192.168.1.7
}

func ExampleRanges() {
	rr := make(iprange.Ranges, 0)
	rr = append(rr, iprange.Parse("192.168.1.0/24"))
	rr = append(rr, iprange.Parse("192.168.2.0/24"))
	fmt.Println(rr.Count())
	fmt.Println(rr.Contains(net.ParseIP("192.168.1.10")))
	fmt.Println(rr.Contains(net.ParseIP("192.168.2.10")))

	// Output:
	// 512
	// true
	// true
}
