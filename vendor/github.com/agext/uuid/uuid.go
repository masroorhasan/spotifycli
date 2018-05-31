// Copyright 2015 ALRUX Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package uuid implements generation and manipulation of UUIDs (v1 defined in RFC 4122).

Version 1 UUIDs are time-based and include a node identifier that can be a MAC address or a random 48-bit value.

This package uses the random approach for the node identifier, setting both the 'multicast' and 'local' bits to make sure the value cannot be confused with a real IEEE 802 address (see section 4.5 of RFC 4122). The initial node identifier is a cryptographic-quality random 46-bit value. The first 30 bits can be set and retrieved with the `SetNodeId` and `NodeId` functions and method, so that they can be used as a hard-coded instance id. The remaining 16 bits are reserved for increasing the randomness of the UUIDs and to avoid collisions on clock sequence rollovers.

The basic generator `New` increments the clock sequence on every call and when the counter rolls over the last 16 bits of the node identifier are regenerated using a PRNG seeded at init()-time with the initial node identifier. This approach sacrifices cryptographic quality for speed and for avoiding depletion of the OS entropy pool (yes, it can and does happen).

The `NewCrypto` generator replaces the clock sequence and last 16 bits of the node identifier on each call with cryptographic-quality random values.
*/
package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mrand "math/rand"
	"strings"
	"sync"
	"time"
)

const (
	gregorianEpoch = 0x01B21DD213814000
)

// UUID is a byte-encoded sequence in the following form:
//
//    0                   1                   2                   3
//     0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//    |                          time_low                             |
//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//    |       time_mid                |         time_hi_and_version   |
//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//    |clk_seq_hi_res |  clk_seq_low  |         node (0-1)            |
//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//    |                         node (2-5)                            |
//    +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
//
//    Field                  Data Type     Octet  Note
//                                         #
//
//    time_low               unsigned 32   0-3    The low field of the
//                           bit integer          timestamp
//
//    time_mid               unsigned 16   4-5    The middle field of the
//                           bit integer          timestamp
//
//    time_hi_and_version    unsigned 16   6-7    The high field of the
//                           bit integer          timestamp multiplexed
//                                                with the version number
//
//    clock_seq_hi_res       unsigned 8    8      The high field of the
//                           bit integer          clock sequence
//                                                multiplexed with the
//                                                variant
//
//    clock_seq_low          unsigned 8    9      The low field of the
//                           bit integer          clock sequence
//
//    node                   unsigned 48   10-15  The spatially unique
//                           bit integer          node identifier
//
// see http://www.ietf.org/rfc/rfc4122.txt section 4.1.2
type UUID []byte

var (
	csanMutex       sync.RWMutex
	randBuf         []byte
	randBufCap      = 256
	randBufOffset   int
	clockSeqAndNode uint64
	clockSeq        uint16
	nodeRand        uint16
	// aliases to allow mocking in tests
	timeNow = time.Now
)

func init() {
	randBuf = make([]byte, 8, randBufCap)
	n, _ := rand.Read(randBuf)
	for i := 0; i < 8 && n < 8; i++ {
		n2, _ := rand.Read(randBuf[n:])
		n += n2
	}
	if n < 8 {
		panic(fmt.Sprintf("uuid.init: Could not generate %d random bytes (got %d)", 8, n))
	}
	// set the variant inside the clock sequence
	randBuf[0] = uint8(randBuf[0]&0x1f | /*variant*/ 1<<5)
	// set the 'local' and 'multicast' bits of the MAC replacement,
	// to avoid conflicts with real MAC addresses.
	randBuf[2] = uint8((randBuf[2] << 2) | 0x03)

	clockSeqAndNode = binary.BigEndian.Uint64(randBuf)
	clockSeq = uint16((clockSeqAndNode >> 48) & 0x1fff)
	nodeRand = uint16(clockSeqAndNode & 0xffff)
	mrand.Seed(int64(clockSeqAndNode))
}

// SetNodeId sets the bits corresponding to the node id.
// Any unsigned 32-bit integer is accepted and the operation is always successful,
// but only the least significant 30 bits are used. An error is returned
// if the discarded, most significant 2 bits are non-zero.
func SetNodeId(nodeId uint32) error {
	csanMutex.Lock()
	// keep the clock sequence, node counter, and the 'local' and 'multicast' bits
	// of the MAC replacement, to avoid conflicts with real MAC addresses.
	clockSeqAndNode = (clockSeqAndNode & 0xffff03000000ffff) |
		(uint64(((nodeId&0x3f000000)<<2)|(nodeId&0x00ffffff)) << 16)
	csanMutex.Unlock()
	if nodeId>>30 != 0 {
		return fmt.Errorf("uuid.SetNodeId: discarded non-zero most significant 2 bits from nodeId %x", nodeId)
	}
	return nil
}

// NodeId returns the current node id used to generate UUIDs.
func NodeId() uint32 {
	csanMutex.Lock()
	nodeId := clockSeqAndNode >> 16
	csanMutex.Unlock()
	return uint32((nodeId & 0x00ffffff) | ((nodeId & 0xfc000000) >> 2))
}

// New creates a new UUID v1 from the current time, clock sequence and node identifier.
func New() UUID {
	uuid := make([]byte, 16)

	csanMutex.Lock()
	if clockSeq = (clockSeq + 1) & 0x1fff; clockSeq == 0 {
		nodeRand = uint16(mrand.Int31n(0x10000))
	}
	clockSeqAndNode = (clockSeqAndNode & 0xe000ffffffff0000) |
		((uint64(clockSeq)) << 48) | uint64(nodeRand)
	binary.BigEndian.PutUint64(uuid[8:], uint64(clockSeqAndNode))
	csanMutex.Unlock()

	ts := fromUnixNano(int64(timeNow().UTC().UnixNano()))
	// "timestamp" multiplexed with version
	binary.BigEndian.PutUint32(uuid[0:4], uint32(ts&0xffffffff))
	binary.BigEndian.PutUint16(uuid[4:6], uint16((ts>>32)&0xffff))
	binary.BigEndian.PutUint16(uuid[6:8], uint16((ts>>48)&0x0fff)| /*version*/ 1<<12)

	return UUID(uuid)
}

// NewCrypto creates a new UUID v1 from the current time, with cryptographic-quality random clock sequence and last 16 bits of the node identifier.
func NewCrypto() UUID {
	uuid := make([]byte, 16)

	csanMutex.Lock()
	n, _ := rand.Read(randBuf[randBufOffset : randBufOffset+4])
	if randBufOffset += n - 4; randBufOffset < 0 {
		randBufOffset = 0
	}
	val := binary.BigEndian.Uint32(randBuf[randBufOffset : randBufOffset+4])
	if randBufOffset += 4; randBufOffset > randBufCap-4 {
		randBufOffset = 0
	}
	clockSeq = uint16((val >> 16) & 0x1fff)
	nodeRand = uint16(val & 0xffff)
	clockSeqAndNode = (clockSeqAndNode & 0xe000ffffffff0000) |
		((uint64(clockSeq)) << 48) | uint64(nodeRand)
	binary.BigEndian.PutUint64(uuid[8:], uint64(clockSeqAndNode))
	csanMutex.Unlock()

	ts := fromUnixNano(int64(timeNow().UTC().UnixNano()))
	// "timestamp" multiplexed with version
	binary.BigEndian.PutUint32(uuid[0:4], uint32(ts&0xffffffff))
	binary.BigEndian.PutUint16(uuid[4:6], uint16((ts>>32)&0xffff))
	binary.BigEndian.PutUint16(uuid[6:8], uint16((ts>>48)&0x0fff)| /*version*/ 1<<12)

	return UUID(uuid)
}

// NewFromBytes creates a UUID from a slice of byte; mostly useful for copying UUIDs.
func NewFromBytes(b []byte) (UUID, error) {
	if len(b) != 16 {
		return nil, fmt.Errorf("uuid.NewFromBytes: Input length is wrong (%d instead of 16)", len(b))
	}
	uuid := make([]byte, 16)
	copy(uuid, b[:16])

	return UUID(uuid), nil
}

// NewFromString creates a UUID from a dash-separated hex string
func NewFromString(s string) (UUID, error) {
	digits := strings.Replace(s, "-", "", -1)
	if hex.DecodedLen(len(digits)) != 16 {
		return nil, fmt.Errorf("uuid.NewFromString: %s is not a valid UUID", s)
	}
	uuid, err := hex.DecodeString(digits)
	if err != nil {
		return nil, fmt.Errorf("uuid.NewFromString: %v", err)
	}

	return UUID(uuid), nil
}

// Hex formats the receiver UUID as a hex string.
func (u UUID) Hex() string {
	return hex.EncodeToString([]byte(u))
}

// String formats the receiver UUID as a dash-separated hex string.
func (u UUID) String() string {
	h := u.Hex()
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

// NodeId extracts the node id from the receiver UUID.
func (u UUID) NodeId() uint32 {
	nodeId := binary.BigEndian.Uint64(u[8:16]) >> 16
	return uint32((nodeId & 0x00ffffff) | ((nodeId & 0xfc000000) >> 2))
}

// Time extracts the time from the receiver UUID as time.Time.
func (u UUID) Time() time.Time {
	timeLow := uint64(binary.BigEndian.Uint32(u[0:4]))
	timeMid := uint64(binary.BigEndian.Uint16(u[4:6]))
	timeHi := uint64((binary.BigEndian.Uint16(u[6:8]) & 0x0fff))
	nanosecs := toUnixNano(int64((timeLow) + (timeMid << 32) + (timeHi << 48)))

	return time.Unix(nanosecs/1e9, nanosecs%1e9).UTC()
}

// Version extracts the version of the receiver UUID.
//
// The following table lists the currently-defined versions for this
// UUID variant.
//
//    Msb0  Msb1  Msb2  Msb3   Version  Description
//
//     0     0     0     1        1     The time-based version
//                                      specified in this document.
//
//     0     0     1     0        2     DCE Security version, with
//                                      embedded POSIX UIDs.
//
//     0     0     1     1        3     The name-based version
//                                      specified in this document
//                                      that uses MD5 hashing.
//
//     0     1     0     0        4     The randomly or pseudo-
//                                      randomly generated version
//                                      specified in this document.
//
//     0     1     0     1        5     The name-based version
//                                      specified in this document
//                                      that uses SHA-1 hashing.
//
// see http://www.ietf.org/rfc/rfc4122.txt section 4.1.3
func (u UUID) Version() int {
	return int((binary.BigEndian.Uint16(u[6:8]) & 0xf000) >> 12)
}

// Variant extracts the variant of the receiver UUID.
//
// The following table lists the contents of the variant field, where
// the letter "x" indicates a "don't-care" value.
//
//    Msb0  Msb1  Msb2  Description
//
//     0     x     x    Reserved, NCS backward compatibility.
//
//     1     0     x    The variant specified in this document.
//
//     1     1     0    Reserved, Microsoft Corporation backward
//                      compatibility
//
//     1     1     1    Reserved for future definition.
//
// see http://www.ietf.org/rfc/rfc4122.txt section 4.1.1
func (u UUID) Variant() int {
	return int((binary.BigEndian.Uint16(u[8:10]) & 0xe000) >> 13)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (u *UUID) UnmarshalJSON(b []byte) error {
	var field string
	if err := json.Unmarshal(b, &field); err != nil {
		return err
	}

	uuid, err := NewFromString(field)
	if err != nil {
		return err
	}

	*u = uuid

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

// fromUnixNano converts a Unix Epoch timestamp of nanosecond precision to Gregorian Epoch.
func fromUnixNano(ts int64) int64 {
	return (ts / 100) + gregorianEpoch
}

// toUnixNano converts a Gregorian Epoch timestamp of nanosecond precision to Unix Epoch.
func toUnixNano(ts int64) int64 {
	return (ts - gregorianEpoch) * 100
}
