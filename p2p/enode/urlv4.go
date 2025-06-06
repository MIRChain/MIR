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
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strconv"

	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
)

var (
	incompleteNodeURL = regexp.MustCompile("(?i)^(?:enode://)?([0-9a-f]+)$")
	lookupIPFunc      = net.LookupIP
)

// MustParseV4 parses a node URL. It panics if the URL is not valid.
func MustParseV4[P crypto.PublicKey](rawurl string) *Node[P] {
	n, err := ParseV4[P](rawurl)
	if err != nil {
		panic("invalid node URL: " + err.Error())
	}
	return n
}

// ParseV4 parses a node URL.
//
// There are two basic forms of node URLs:
//
//   - incomplete nodes, which only have the public key (node ID)
//   - complete nodes, which contain the public key and IP/Port information
//
// For incomplete nodes, the designator must look like one of these
//
//    enode://<hex node id>
//    <hex node id>
//
// For complete nodes, the node ID is encoded in the username portion
// of the URL, separated from the host by an @ sign. The hostname can
// only be given as an IP address or using DNS domain name.
// The port in the host name section is the TCP listening port. If the
// TCP and UDP (discovery) ports differ, the UDP port is specified as
// query parameter "discport".
//
// In the following example, the node URL describes
// a node with IP address 10.3.58.6, TCP listening port 30303
// and UDP discovery port 30301.
//
//    enode://<hex node id>@10.3.58.6:30303?discport=30301
func ParseV4[P crypto.PublicKey](rawurl string) (*Node[P], error) {
	if m := incompleteNodeURL.FindStringSubmatch(rawurl); m != nil {
		id, err := parsePubkey[P](m[1])
		if err != nil {
			return nil, fmt.Errorf("invalid public key (%v)", err)
		}
		return NewV4(id, nil, 0, 0), nil
	}
	return parseComplete[P](rawurl)
}

// NewV4 creates a node from discovery v4 node information. The record
// contained in the node has a zero-length signature.
func NewV4[P crypto.PublicKey] (pubkey P, ip net.IP, tcp, udp int) *Node[P] {
	var r enr.Record
	if len(ip) > 0 {
		r.Set(enr.IP(ip))
	}
	return newV4(pubkey, r, tcp, udp)
}

// broken out from `func NewV4` (above) same in upstream go-ethereum, but taken out
// to avoid code duplication b/t NewV4 and NewV4Hostname
func newV4[P crypto.PublicKey](pubkey P, r enr.Record, tcp, udp int) *Node[P] {
	if udp != 0 {
		r.Set(enr.UDP(udp))
	}
	if tcp != 0 {
		r.Set(enr.TCP(tcp))
	}
	signV4Compat(&r, pubkey)
	n, err := New[P](v4CompatID[P]{}, &r)
	if err != nil {
		panic(err)
	}
	return n
}

// isNewV4 returns true for nodes created by NewV4.
func isNewV4[P crypto.PublicKey](n *Node[P]) bool {
	var k s256raw
	return n.r.IdentityScheme() == "" && n.r.Load(&k) == nil && len(n.r.Signature()) == 0
}

// Quorum

// NewV4Hostname creates a node from discovery v4 node information. The record
// contained in the node has a zero-length signature. It sets the hostname or ip
// of the node depends on hostname context
func NewV4Hostname[P crypto.PublicKey] (pubkey P, hostname string, tcp, udp, raftPort int) *Node[P] {
	var r enr.Record

	if ip := net.ParseIP(hostname); ip == nil {
		r.Set(enr.Hostname(hostname))
	} else {
		r.Set(enr.IP(ip))
	}

	if raftPort != 0 {
		r.Set(enr.RaftPort(raftPort))
	}

	return newV4[P](pubkey, r, tcp, udp)
}

// End-Quorum

func parseComplete[P crypto.PublicKey](rawurl string) (*Node[P], error) {
	var (
		id               P
		ip               net.IP
		tcpPort, udpPort uint64
	)
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "enode" {
		return nil, errors.New("invalid URL scheme, want \"enode\"")
	}
	// Parse the Node ID from the user portion.
	if u.User == nil {
		return nil, errors.New("does not contain node ID")
	}
	if id, err = parsePubkey[P](u.User.String()); err != nil {
		return nil, fmt.Errorf("invalid public key (%v)", err)
	}
	qv := u.Query()
	// Parse the IP address.
	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		// Quorum: if IP look up fail don't return error for raft url
		if qv.Get("raftport") == "" {
			return nil, err
		}
	} else {
		ip = ips[0]
		// Ensure the IP is 4 bytes long for IPv4 addresses.
		if ipv4 := ip.To4(); ipv4 != nil {
			ip = ipv4
		}
	}
	// Parse the port numbers.
	if tcpPort, err = strconv.ParseUint(u.Port(), 10, 16); err != nil {
		return nil, errors.New("invalid port")
	}
	udpPort = tcpPort

	if qv.Get("discport") != "" {
		udpPort, err = strconv.ParseUint(qv.Get("discport"), 10, 16)
		if err != nil {
			return nil, errors.New("invalid discport in query")
		}
	}

	// Quorum
	if qv.Get("raftport") != "" {
		raftPort, err := strconv.ParseUint(qv.Get("raftport"), 10, 16)
		if err != nil {
			return nil, errors.New("invalid raftport in query")
		}
		if u.Hostname() == "" {
			return nil, errors.New("empty hostname in raft url")
		}
		return NewV4Hostname[P](id, u.Hostname(), int(tcpPort), int(udpPort), int(raftPort)), nil
	}
	// End-Quorum

	return NewV4[P](id, ip, int(tcpPort), int(udpPort)), nil
}

func HexPubkey[P crypto.PublicKey](h string) (P, error) {
	k, err := parsePubkey[P](h)
	if err != nil {
		return crypto.ZeroPublicKey[P](), err
	}
	return k, err
}

// parsePubkey parses a hex-encoded secp256k1 public key.
func parsePubkey[P crypto.PublicKey](in string) (P, error) {
	b, err := hex.DecodeString(in)
	if err != nil {
		return crypto.ZeroPublicKey[P](), err
	} else if len(b) != 64 {
		return crypto.ZeroPublicKey[P](), fmt.Errorf("wrong length, want %d hex chars", 128)
	}
	b = append([]byte{0x4}, b...)
	return crypto.UnmarshalPubkey[P](b)
}

// used by Quorum RAFT - to derive enodeID
func (n *Node[P]) EnodeID() string {
	var (
		scheme enr.ID
		nodeid string
		key    P
	)
	n.Load(&scheme)
	switch p:=any(&key).(type){
	case *nist.PublicKey:
		n.Load((*Secp256k1)(p))
	case *gost3410.PublicKey:
		n.Load((*Gost3410)(p))
	}
	switch {
	case scheme == "v4" || !reflect.ValueOf(&key).IsZero():
		nodeid = fmt.Sprintf("%x", crypto.FromECDSAPub(key)[1:])
	default:
		nodeid = fmt.Sprintf("%s.%x", scheme, n.id[:])
	}
	return nodeid
}

func (n *Node[P]) URLv4() string {
	var (
		scheme enr.ID
		nodeid string
		key    P
	)
	n.Load(&scheme)
	switch p:=any(&key).(type){
	case *nist.PublicKey:
		n.Load((*Secp256k1)(p))
	case *gost3410.PublicKey:
		n.Load((*Gost3410)(p))
	}
	switch {
	case scheme == "v4" ||!reflect.ValueOf(&key).IsZero():
		nodeid = fmt.Sprintf("%x", crypto.FromECDSAPub(key)[1:])
	default:
		nodeid = fmt.Sprintf("%s.%x", scheme, n.id[:])
	}
	u := url.URL{Scheme: "enode"}
	if n.Incomplete() {
		u.Host = nodeid
	} else {
		u.User = url.User(nodeid)
		if n.Host() != "" && net.ParseIP(n.Host()) == nil {
			// Quorum
			u.Host = net.JoinHostPort(n.Host(), strconv.Itoa(n.TCP()))
		} else {
			addr := net.TCPAddr{IP: n.IP(), Port: n.TCP()}
			u.Host = addr.String()
		}
		if n.UDP() != n.TCP() {
			u.RawQuery = "discport=" + strconv.Itoa(n.UDP())
		}
		// Quorum
		if n.HasRaftPort() {
			raftQuery := "raftport=" + strconv.Itoa(n.RaftPort())
			if len(u.RawQuery) > 0 {
				u.RawQuery = u.RawQuery + "&" + raftQuery
			} else {
				u.RawQuery = raftQuery
			}
		}
	}
	return u.String()
}

// PubkeyToIDV4 derives the v4 node address from the given public key.
func PubkeyToIDV4[P crypto.PublicKey ] (key P) ID {
	switch pubkey := any(&key).(type) {
	case *nist.PublicKey:
		e := make([]byte, 64)
		math.ReadBits(pubkey.X, e[:len(e)/2])
		math.ReadBits(pubkey.Y, e[len(e)/2:])
		return ID(crypto.Keccak256Hash[P](e))
	case *gost3410.PublicKey:
		e := make([]byte, 64)
		math.ReadBits(pubkey.X, e[:len(e)/2])
		math.ReadBits(pubkey.Y, e[len(e)/2:])
		return ID(crypto.Keccak256Hash[P](e))
	default:
		panic("cant infer type of public key")
	}
}

func PubkeyToEnodeID[P crypto.PublicKey ](key P) EnodeID {
	switch pubkey := any(&key).(type) {
	case *nist.PublicKey:
		e := make([]byte, 64)
		math.ReadBits(pubkey.X, e[:len(e)/2])
		math.ReadBits(pubkey.Y, e[len(e)/2:])
		return EnodeID(e)
	case *gost3410.PublicKey:
		e := make([]byte, 64)
		math.ReadBits(pubkey.X, e[:len(e)/2])
		math.ReadBits(pubkey.Y, e[len(e)/2:])
		return EnodeID(e)
	default:
		panic("cant infer type of public key")
	}
}
