package iprange_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/russtone/iprange"
)

func TestInvalid(t *testing.T) {
	tests := []string{
		//
		// Misc
		//

		"",
		"invalid",

		//
		// IPv4
		//

		// sigle
		"192.168.1.",
		"192.168.300.1",
		"192.168.1.1.1",
		".168.1.1",
		"192.168.1",

		// cidr
		"192.168.1.1/33",
		"192.168.1.1/",
		"192.168.1-2.1/24",
		"192.168.1.1-2/24",
		"192.168.1.1/24/2",

		// begin_end
		"192.168.1.1-192.168",
		"192.168.1.1-192.168.1.10/24",
		"192.168.1.1-1:2::1",

		// Octet
		"192.168.1.1-300",
		"192.168.1.1-2-3",
		"192.168.1,2,.1",

		//
		// IPv6
		//

		// single
		":100:abab:dead",
		":100:abab:10000",
		"::100:abab:dead::",
		"1:2:3:4:5:6:7:8:9",
		"1:2:3:4",
		"1:2:3:4::5:6:7:8",

		// cidr
		"1::abab:dead/1/1",
		"1::abab:dead/129",
		"1::1-2:dead/129",

		// begin_end
		"1::abab:dead_1:2:3",
		"1-2::abab:dead_1:2:3",
		"1::abab:dead_1:2::3/24",

		// octets
		"1-2-3::abab:dead",
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d/%s", i, tt), func(t *testing.T) {
			r := iprange.Parse(tt)
			assert.Nil(t, r)
		})
	}
}

func TestValid(t *testing.T) {
	tests := []string{

		//
		// IPv4
		//

		// single
		"192.168.1.1",

		// cidr
		"192.168.1.1/24",

		// begin_end
		"192.168.1.1_192.168.2.10",

		// octets
		"192.168.1,2-5.1,2,3",

		//
		// IPv6
		//

		// single
		"::",
		"1:2:3:4::abab:dead",
		"::1:2:3:4:abab:dead",
		"1:2:3:4:abab:dead::",
		"::ffff:192.168.1.1",

		// cidr
		"1:2:3:4::abab:dead/120",
		"::/120",
		"1:2:3:4::/120",

		// begin_end
		"1:2:3:4::abab:1_1:2:3:4::abab:10",
		"1:2:3:4::_1:2:3:4::5",

		// octets
		"1:2:3:4::1-10:1,2,ffff",
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d/%s", i, tt), func(t *testing.T) {
			r := iprange.Parse(tt)
			assert.NotNil(t, r)
		})
	}
}
