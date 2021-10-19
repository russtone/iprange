# iprange

[![build](https://github.com/russtone/iprange/actions/workflows/build.yml/badge.svg)](https://github.com/russtone/iprange/actions/workflows/build.yml)
[![Go report](https://goreportcard.com/badge/github.com/russtone/iprange)](https://goreportcard.com/report/github.com/russtone/iprange)
[![codecov](https://codecov.io/gh/russtone/iprange/branch/master/graph/badge.svg?token=c8Fim47Tp0)](https://codecov.io/gh/russtone/iprange)
[![Documentation](https://pkg.go.dev/badge/github.com/russtone/iprange.svg)](https://pkg.go.dev/github.com/russtone/iprange)

Package iprange provides functions to work with different IP ranges.
Supports IPv4 and IPv6.

Supported ranges formats:

- Single address. Examples: `192.168.1.1`, `2001:db8:a0b:12f0::1`
- CIDR. Examples: `192.168.1.0/24`, `2001:db8:a0b:12f0::1`
- Begin_End. Examples: `192.168.1.10_192.168.2.20`, `2001:db8:a0b:12f0::1_2001:db8:a0b:12f0::10`
- Octets ranges: `192.168.1,3-5.1-10`, `2001:db8:a0b:12f0::1,1-10`

For more information see the docs.
