// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package enode

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"
	"net"
	"strings"

	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
	"github.com/pavelkrolevets/MIR-pro/rlp"
)

var errMissingPrefix = errors.New("missing 'enr:' prefix for base64-encoded record")

// Node represents a host on the network.
type Node [P crypto.PublicKey] struct {
	r  enr.Record
	id ID
}

// New wraps a node record. The record must be valid according to the given
// identity scheme.
func New[P crypto.PublicKey] (validSchemes enr.IdentityScheme, r *enr.Record) (*Node[P], error) {
	if err := r.VerifySignature(validSchemes); err != nil {
		return nil, err
	}
	node := &Node[P]{r: *r}
	if n := copy(node.id[:], validSchemes.NodeAddr(&node.r)); n != len(ID{}) {
		return nil, fmt.Errorf("invalid node ID length %d, need %d", n, len(ID{}))
	}
	return node, nil
}

// MustParse parses a node record or enode:// URL. It panics if the input is invalid.
func MustParse[P crypto.PublicKey] (rawurl string) *Node[P] {
	n, err := Parse[P](enr.SchemeMap{"v4": V4ID[P]{}}, rawurl)
	if err != nil {
		panic("invalid node: " + err.Error())
	}
	return n
}

// Parse decodes and verifies a base64-encoded node record.
func Parse[P crypto.PublicKey] (validSchemes enr.IdentityScheme, input string) (*Node[P], error) {
	if strings.HasPrefix(input, "enode://") {
		return ParseV4[P](input)
	}
	if !strings.HasPrefix(input, "enr:") {
		return nil, errMissingPrefix
	}
	bin, err := base64.RawURLEncoding.DecodeString(input[4:])
	if err != nil {
		return nil, err
	}
	var r enr.Record
	if err := rlp.DecodeBytes(bin, &r); err != nil {
		return nil, err
	}
	return New[P](validSchemes, &r)
}

// ID returns the node identifier.
func (n *Node[P]) ID() ID {
	return n.id
}

// Seq returns the sequence number of the underlying record.
func (n *Node[P]) Seq() uint64 {
	return n.r.Seq()
}

// Quorum
// Incomplete returns true for nodes with no IP address and no hostname if with raftport.
func (n *Node[P]) Incomplete() bool {
	return n.IP() == nil && (!n.HasRaftPort() || (n.Host() == "" && n.HasRaftPort()))
}

// Load retrieves an entry from the underlying record.
func (n *Node[P]) Load(k enr.Entry) error {
	return n.r.Load(k)
}

// IP returns the IP address of the node.
//
// Quorum
// To support DNS lookup in node ip. The function performs hostname lookup if hostname is defined in enr.Hostname
// and falls back to enr.IP value in case of failure. It also makes sure the resolved IP is in IPv4 or IPv6 format
func (n *Node[P]) IP() net.IP {
	if n.Host() == "" {
		// no host is set, so use the IP directly
		return n.loadIP()
	}
	// attempt to look up IP addresses if host is a FQDN
	lookupIPs, err := net.LookupIP(n.Host())
	if err != nil {
		return n.loadIP()
	}
	// set to first ip by default & as Ethereum upstream
	ip := lookupIPs[0]
	// Ensure the IP is 4 bytes long for IPv4 addresses.
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}
	return ip
}

func (n *Node[P]) loadIP() net.IP {
	var (
		ip4 enr.IPv4
		ip6 enr.IPv6
	)
	if n.Load(&ip4) == nil {
		return net.IP(ip4)
	}
	if n.Load(&ip6) == nil {
		return net.IP(ip6)
	}
	return nil
}

// Quorum
func (n *Node[P]) Host() string {
	var hostname string
	n.Load((*enr.Hostname)(&hostname))
	return hostname
}

// End-Quorum

// UDP returns the UDP port of the node.
func (n *Node[P]) UDP() int {
	var port enr.UDP
	n.Load(&port)
	return int(port)
}

// used by Quorum RAFT - returns the Raft port of the node
func (n *Node[P]) RaftPort() int {
	var port enr.RaftPort
	err := n.Load(&port)
	if err != nil {
		return 0
	}
	return int(port)
}

func (n *Node[P]) HasRaftPort() bool {
	return n.RaftPort() > 0
}

// UDP returns the TCP port of the node.
func (n *Node[P]) TCP() int {
	var port enr.TCP
	n.Load(&port)
	return int(port)
}

// Pubkey returns the secp256k1 public key of the node, if present.
func (n *Node[P]) Pubkey() P {
	var key P
	switch p:= any(&key).(type){
	case *nist.PublicKey:
		if n.Load((*Secp256k1)(p)) != nil {
			return crypto.ZeroPublicKey[P]()
		}
	case *gost3410.PublicKey:
		if n.Load((*Gost3410)(p)) != nil {
			return crypto.ZeroPublicKey[P]()
		}
	}
	return key
}

// Record returns the node's record. The return value is a copy and may
// be modified by the caller.
func (n *Node[P]) Record() *enr.Record {
	cpy := n.r
	return &cpy
}

// ValidateComplete checks whether n has a valid IP and UDP port.
// Deprecated: don't use this method.
func (n *Node[P]) ValidateComplete() error {
	if n.Incomplete() {
		return errors.New("missing IP address")
	}
	if n.UDP() == 0 {
		return errors.New("missing UDP port")
	}
	ip := n.IP()
	if ip.IsMulticast() || ip.IsUnspecified() {
		return errors.New("invalid IP (multicast/unspecified)")
	}
	// Validate the node key (on curve, etc.).
	var pub P
	switch any(&pub).(type){
	case *nist.PublicKey:
		var key Secp256k1
		return n.Load(&key)
	case *gost3410.PublicKey:
		var key Gost3410
		return n.Load(&key)
	default:
		return fmt.Errorf("error validating the node key (on curve, etc.)")
	}
}

// String returns the text representation of the record.
func (n *Node[P]) String() string {
	if isNewV4(n) {
		return n.URLv4() // backwards-compatibility glue for NewV4 nodes
	}
	enc, _ := rlp.EncodeToBytes(&n.r) // always succeeds because record is valid
	b64 := base64.RawURLEncoding.EncodeToString(enc)
	return "enr:" + b64
}

// MarshalText implements encoding.TextMarshaler.
func (n *Node[P]) MarshalText() ([]byte, error) {
	return []byte(n.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *Node[P]) UnmarshalText(text []byte) error {
	dec, err := Parse[P](enr.SchemeMap{"v4": V4ID[P]{}}, string(text))
	if err == nil {
		*n = *dec
	}
	return err
}

// ID is a unique identifier for each node.
type ID [32]byte

// Bytes returns a byte slice representation of the ID
func (n ID) Bytes() []byte {
	return n[:]
}

// ID prints as a long hexadecimal number.
func (n ID) String() string {
	return fmt.Sprintf("%x", n[:])
}

// The Go syntax representation of a ID is a call to HexID.
func (n ID) GoString() string {
	return fmt.Sprintf("enode.HexID(\"%x\")", n[:])
}

// TerminalString returns a shortened hex string for terminal logging.
func (n ID) TerminalString() string {
	return hex.EncodeToString(n[:8])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n ID) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(n[:])), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (n *ID) UnmarshalText(text []byte) error {
	id, err := ParseID(string(text))
	if err != nil {
		return err
	}
	*n = id
	return nil
}

// ID is a unique identifier for each node used by RAFT
type EnodeID [64]byte

// ID prints as a long hexadecimal number.
func (n EnodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

// The Go syntax representation of a ID is a call to HexID.
func (n EnodeID) GoString() string {
	return fmt.Sprintf("enode.HexID(\"%x\")", n[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n EnodeID) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(n[:])), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (n *EnodeID) UnmarshalText(text []byte) error {
	id, err := RaftHexID(string(text))
	if err != nil {
		return err
	}
	*n = id
	return nil
}

// HexID converts a hex string to an ID.
// The string may be prefixed with 0x.
// It panics if the string is not a valid ID.
func HexID(in string) ID {
	id, err := ParseID(in)
	if err != nil {
		panic(err)
	}
	return id
}

// used by Quorum RAFT to derive 64 byte nodeId from 128 byte enodeID
func RaftHexID(in string) (EnodeID, error) {
	var id EnodeID
	b, err := hex.DecodeString(strings.TrimPrefix(in, "0x"))
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("wrong length, want %d hex chars", len(id)*2)
	}

	copy(id[:], b)
	return id, nil
}

func ParseID(in string) (ID, error) {
	var id ID
	b, err := hex.DecodeString(strings.TrimPrefix(in, "0x"))
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("wrong length, want %d hex chars", len(id)*2)
	}
	copy(id[:], b)
	return id, nil
}

// DistCmp compares the distances a->target and b->target.
// Returns -1 if a is closer to target, 1 if b is closer to target
// and 0 if they are equal.
func DistCmp(target, a, b ID) int {
	for i := range target {
		da := a[i] ^ target[i]
		db := b[i] ^ target[i]
		if da > db {
			return 1
		} else if da < db {
			return -1
		}
	}
	return 0
}

// LogDist returns the logarithmic distance between a and b, log2(a ^ b).
func LogDist(a, b ID) int {
	lz := 0
	for i := range a {
		x := a[i] ^ b[i]
		if x == 0 {
			lz += 8
		} else {
			lz += bits.LeadingZeros8(x)
			break
		}
	}
	return len(a)*8 - lz
}
