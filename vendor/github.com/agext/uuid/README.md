# A Go package for generating and manipulating UUIDs

[![Release](https://img.shields.io/github/release/agext/uuid.svg?style=flat)](https://github.com/agext/uuid/releases/latest)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/agext/uuid) 
[![Build Status](https://travis-ci.org/agext/uuid.svg?branch=master&style=flat)](https://travis-ci.org/agext/uuid)
[![Coverage Status](https://coveralls.io/repos/github/agext/uuid/badge.svg?style=flat)](https://coveralls.io/github/agext/uuid)
[![Go Report Card](https://goreportcard.com/badge/github.com/agext/uuid?style=flat)](https://goreportcard.com/report/github.com/agext/uuid)


Generate, encode, and decode UUIDs v1, as defined in [RFC 4122](http://www.ietf.org/rfc/rfc4122.txt), in [Go](http://golang.org).

## Project Status

v1.0.1 Stable: Guaranteed no breaking changes to the API in future v1.x releases. Probably safe to use in production, though provided on "AS IS" basis.

This package is being actively maintained. If you encounter any problems or have any suggestions for improvement, please [open an issue](https://github.com/agext/uuid/issues). Pull requests are welcome.

## Overview

Package uuid implements generation and manipulation of UUIDs (v1 defined in RFC 4122).

Version 1 UUIDs are time-based and include a node identifier that can be a MAC address or a random 48-bit value.

This package uses the random approach for the node identifier, setting both the 'multicast' and 'local' bits to make sure the value cannot be confused with a real IEEE 802 address (see section 4.5 of RFC 4122). The initial node identifier is a cryptographic-quality random 46-bit value. The first 30 bits can be set and retrieved with the `SetNodeId` and `NodeId` functions and method, so that they can be used as a hard-coded instance id. The remaining 16 bits are reserved for increasing the randomness of the UUIDs and to avoid collisions on clock sequence rollovers.

The basic generator `New` increments the clock sequence on every call and when the counter rolls over the last 16 bits of the node identifier are regenerated using a PRNG seeded at init()-time with the initial node identifier. This approach sacrifices cryptographic quality for speed and for avoiding depletion of the OS entropy pool (yes, it can and does happen).

The `NewCrypto` generator replaces the clock sequence and last 16 bits of the node identifier on each call with cryptographic-quality random values.

## Installation

```
go get github.com/agext/uuid
```

## License

Package uuid is released under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
