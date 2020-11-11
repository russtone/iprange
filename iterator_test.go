package iprange_test

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/russtone/iprange"
)

func TestIterator(t *testing.T) {
	tests := []struct {
		s        []string
		res      []string
		ip       string
		contains bool
	}{
		//
		// IPv4
		//

		// single
		{
			[]string{"104.244.42.65"},
			[]string{"104.244.42.65"},
			"104.244.42.65",
			true,
		},

		// cidr
		{
			[]string{"87.240.129.133/30"},
			[]string{
				"87.240.129.132",
				"87.240.129.133",
				"87.240.129.134",
				"87.240.129.135",
			},
			"87.240.129.134",
			true,
		},
		{
			[]string{"87.240.129.133/32"},
			[]string{
				"87.240.129.133",
			},
			"87.240.129.133",
			true,
		},

		// begin\end
		{
			[]string{"35.231.145.10\\35.231.145.13"},
			[]string{
				"35.231.145.10",
				"35.231.145.11",
				"35.231.145.12",
				"35.231.145.13",
			},
			"35.231.145.14",
			false,
		},

		// Octet
		{
			[]string{"104.16.99-100.10,52-53"},
			[]string{
				"104.16.99.10",
				"104.16.99.52",
				"104.16.99.53",
				"104.16.100.10",
				"104.16.100.52",
				"104.16.100.53",
			},
			"104.16.101.53",
			false,
		},

		//
		// IPv6
		//

		// single
		{
			[]string{"2001:db8::1"},
			[]string{"2001:db8::1"},
			"2001:db8::1",
			true,
		},

		// CIDR
		{
			[]string{"2001:db8::1/126"},
			[]string{
				"2001:db8::",
				"2001:db8::1",
				"2001:db8::2",
				"2001:db8::3",
			},
			"2001:db8::",
			true,
		},

		// begin\end
		{
			[]string{"2001:db8::\\2001:db8::5"},
			[]string{
				"2001:db8::",
				"2001:db8::1",
				"2001:db8::2",
				"2001:db8::3",
				"2001:db8::4",
				"2001:db8::5",
			},
			"2001:db8::6",
			false,
		},

		// octets
		{
			[]string{"2001:db8:9,10-12::"},
			[]string{
				"2001:db8:9::",
				"2001:db8:10::",
				"2001:db8:11::",
				"2001:db8:12::",
			},
			"2001:db8:9::1",
			false,
		},

		// Multiple
		{
			[]string{"104.244.42.65", "87.240.129.133/30"},
			[]string{
				"104.244.42.65",
				"87.240.129.132",
				"87.240.129.133",
				"87.240.129.134",
				"87.240.129.135",
			},
			"87.240.129.133",
			true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d/%s", i, strings.Join(tt.s, ",")), func(t *testing.T) {
			// Parse.
			rr := make(iprange.Ranges, 0)
			for _, rng := range tt.s {
				r := iprange.Parse(rng)
				require.NotNil(t, r)

				rr = append(rr, r)
			}

			it := rr.Iterator()

			// Check count.
			assert.EqualValues(t, len(tt.res), it.Count().Int64())

			// Check contains.
			assert.EqualValues(t, tt.contains, it.Contains(net.ParseIP(tt.ip)))

			// Get results.
			res := make([]string, 0)
			var ip net.IP

			for it.Next(&ip) {
				res = append(res, ip.String())
			}
			assert.Equal(t, tt.res, res)

			// Reset and get results again.
			it.Reset()
			res = make([]string, 0)
			for it.Next(&ip) {
				res = append(res, ip.String())
			}
			assert.Equal(t, tt.res, res, "after reset")
		})
	}
}
