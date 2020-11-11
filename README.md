# iprange

[![Build and test](https://img.shields.io/github/workflow/status/russtone/iprange/Build%20and%20test)](https://github.com/russtone/iprange/actions?query=workflow%3A%22Build+and+test%22)
[![Go report](https://goreportcard.com/badge/github.com/russtone/iprange)](https://goreportcard.com/report/github.com/russtone/iprange)
[![Code coverage](https://img.shields.io/codecov/c/gh/russtone/iprange.svg)](https://codecov.io/gh/russtone/iprange)
[![Documentation](https://godoc.org/github.com/russtone/iprange?status.svg)](http://godoc.org/github.com/russtone/iprange)

Package iprange provides functions to work with different IP ranges.
Supports IPv4 and IPv6.

Supported ranges formats:

- Single address. Examples: `192.168.1.1`, `2001:db8:a0b:12f0::1`
- CIDR. Examples: `192.168.1.0/24`, `2001:db8:a0b:12f0::1`
- Begin\End. Examples: `192.168.1.10\192.168.2.20`, `2001:db8:a0b:12f0::1\2001:db8:a0b:12f0::10`
- Octets ranges: `192.168.1,3-5.1-10`, `2001:db8:a0b:12f0::1,1-10`

For more information see the docs.
