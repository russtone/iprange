package iprange_test

import (
	"fmt"
	"math/big"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/russtone/iprange"
)

func TestContains(t *testing.T) {
	tests := []struct {
		ss       []string
		ip       string
		contains bool
	}{
		//
		// IPv4
		//

		// single
		{[]string{"104.244.42.65"}, "104.244.42.65", true},
		{[]string{"87.250.250.242"}, "74.125.131.139", false},

		// cidr
		{[]string{"87.240.129.133/24"}, "87.240.129.211", true},
		{[]string{"87.240.129.133/24"}, "87.240.130.211", false},
		{[]string{"140.82.118.4/22"}, "140.82.119.4", true},
		{[]string{"140.82.118.4/22"}, "140.83.119.4", false},

		// begin_end
		{[]string{"35.231.145.151_35.231.150.10"}, "35.231.150.10", true},
		{[]string{"35.231.145.151_35.231.150.10"}, "35.231.140.100", false},

		// octets
		{[]string{"104.16.99-100.52-55"}, "104.16.99.53", true},
		{[]string{"104.16.99.52-255"}, "104.16.99.255", true},
		{[]string{"104.16.99-100.52-55"}, "104.16.98.2", false},

		//
		// IPv6
		//

		// single
		{[]string{"2001:4860:0:2001::68"}, "2001:4860:0:2001::68", true},
		{[]string{"2001:4860:0000:2001:0000:0000:0000:0068"}, "2001:4860:0:2001::68", true},

		// cidr
		{[]string{"2001:db8::/48"}, "2001:db8::10", true},
		{[]string{"2001:db8::/48"}, "2001:ab8::", false},

		// begin_end
		{[]string{"2001:db8::_2001:db8::10"}, "2001:db8::5", true},
		{[]string{"2001:db8::_2001:db8::10"}, "2001:db8::abab", false},

		// octets
		{[]string{"2001:DB8:3C4D:7777::123-130"}, "2001:DB8:3C4D:7777::124", true},
		{[]string{"2001:DB8:3C4D:7777::123-130"}, "2001:DB8:3C4D:7777::dead", false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d/%s", i, strings.Join(tt.ss, ",")), func(t *testing.T) {
			rr := make(iprange.Ranges, 0)
			for _, s := range tt.ss {
				r := iprange.Parse(s)
				require.NotNil(t, r)
				rr = append(rr, r)
			}

			contains := rr.Contains(net.ParseIP(tt.ip))
			assert.Equal(t, tt.contains, contains)
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		ss    []string
		count int64
	}{
		//
		// IPv4
		//

		// single
		{[]string{"104.244.42.65"}, 1},

		// cidr
		{[]string{"87.240.129.133/24"}, 256},
		{[]string{"87.240.129.133/32"}, 1},
		{[]string{"87.240.129.133/0"}, 2 << 31},

		// begin_end
		{[]string{"35.231.145.151_35.231.145.200"}, 50},

		// octets
		{[]string{"104.16.99-100.52-55"}, 8},

		//
		// IPv6
		//

		// single
		{[]string{"2001:4860:0:2001::68"}, 1},

		// begin_end
		{[]string{"2001:db8::_2001:db8::10"}, 17},

		// cidr
		{[]string{"2001:4860:0:2001::68/128"}, 1},
		{[]string{"2001:4860:0:2001::68/127"}, 2},
		{[]string{"2001:4860:0:2001::68/100"}, 2 << 27},

		// octets
		{[]string{"2001:DB8:3C4D:7777::123-130"}, 14},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d/%s", i, strings.Join(tt.ss, ",")), func(t *testing.T) {
			rr := make(iprange.Ranges, 0)
			for _, s := range tt.ss {
				r := iprange.Parse(s)
				require.NotNil(t, r)
				rr = append(rr, r)
			}

			assert.Equal(t, big.NewInt(tt.count), rr.Count())
		})
	}
}
