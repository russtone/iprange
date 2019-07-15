package iprange_test

import (
	"net"
	"testing"

	"github.com/russtone/iprange"
)

var rangetests = []struct {
	rng string
	ip  string
	err bool
	res bool
}{
	// One IP
	{"104.244.42.65", "104.244.42.65", false, true},
	{"87.250.250.242", "74.125.131.139", false, false},

	// CIDR
	{"87.240.129.133/24", "87.240.129.211", false, true},
	{"87.240.129.133/24", "87.240.130.211", false, false},
	{"140.82.118.4/22", "140.82.119.4", false, true},
	{"140.82.118.4/22", "140.83.119.4", false, false},

	// IPv4 dash range
	{"35.231.145.151-35.231.150.10", "35.231.148.1", false, true},
	{"35.231.145.151-35.231.150.10", "35.231.150.10", false, true},
	{"35.231.145.151-35.231.150.10", "35.231.145.151", false, true},
	{"35.231.145.151-35.231.150.10", "35.231.140.100", false, false},
	{"35.231.145.151-35.231.150.10", "35.231.150.100", false, false},

	// IPv4 octet dash/asterisk range
	{"104.16.99-100.52-55", "104.16.99.53", false, true},
	{"104.16.99-100.52-55", "104.16.99.52", false, true},
	{"104.16.99-100.52-55", "104.16.99.55", false, true},
	{"104.16.99-100.52-55", "104.16.100.55", false, true},
	{"104.16.99-100.52-55", "104.16.100.52", false, true},
	{"104.16.99-100.52-55", "104.16.98.2", false, false},
	{"104.16.99-100.52-55", "104.16.99.50", false, false},
	{"104.16.99.*", "104.16.99.0", false, true},
	{"104.16.99.*", "104.16.99.255", false, true},
	{"104.16.99.*", "104.16.98.50", false, false},

	// Invalid ranges
	{"invalid", "104.16.98.50", true, false},
	{"104.16.99.*/24", "104.16.98.50", true, false},
	{"104.16-14.99.10", "104.16.98.50", true, false},
}

func TestContains(t *testing.T) {
	for _, tt := range rangetests {
		t.Run(tt.rng, func(t *testing.T) {
			r, err := iprange.Parse(tt.rng)

			if err != nil {
				if !tt.err {
					t.Error(err)
				} else {
					return
				}
			}

			res := r.Contains(net.ParseIP(tt.ip))
			if res != tt.res {
				t.Errorf("invalid result for %s", tt.ip)
			}
		})
	}
}
