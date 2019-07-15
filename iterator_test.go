package iprange_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/russtone/iprange"
)

var itertests = []struct {
	rng string
	res []string
	err bool
}{
	// One IP
	{"104.244.42.65", []string{"104.244.42.65"}, false},

	// CIDR
	{"87.240.129.133/30", []string{
		"87.240.129.132",
		"87.240.129.133",
		"87.240.129.134",
		"87.240.129.135",
	}, false},
	{"87.240.129.133/32", []string{
		"87.240.129.133",
	}, false},

	// IPv4 dash range
	{"35.231.145.10-35.231.145.13", []string{
		"35.231.145.10",
		"35.231.145.11",
		"35.231.145.12",
		"35.231.145.13",
	}, false},

	// IPv4 octet dash/asterisk range
	{"104.16.99-100.52-53", []string{
		"104.16.99.52",
		"104.16.99.53",
		"104.16.100.52",
		"104.16.100.53",
	}, false},
}

func TestIterator(t *testing.T) {
	for _, tt := range itertests {
		t.Run(tt.rng, func(t *testing.T) {
			r, err := iprange.Parse(tt.rng)

			if err != nil {
				if !tt.err {
					t.Error(err)
				} else {
					return
				}
			}

			it := iprange.NewIterator(r)

			res := make([]string, 0)
			ip := net.IPv4(0, 0, 0, 0)

			for it.Next(ip) {
				res = append(res, ip.String())
			}

			if !reflect.DeepEqual(res, tt.res) {
				t.Errorf("invalid result: %+v", res)
			}
		})
	}
}
