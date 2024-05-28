// Copyright 2019 The go-ethereum Authors
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

// Package v4wire implements the Discovery v4 Wire Protocol.
package v4wire

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/MIRChain/MIR/common/math"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/csp"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/nist"
	"github.com/MIRChain/MIR/p2p/enode"
	"github.com/MIRChain/MIR/p2p/enr"
	"github.com/MIRChain/MIR/rlp"
)

// RPC packet types
const (
	PingPacket = iota + 1 // zero is 'reserved'
	PongPacket
	FindnodePacket
	NeighborsPacket
	ENRRequestPacket
	ENRResponsePacket
)

// RPC request structures
type (
	Ping struct {
		Version    uint
		From, To   Endpoint
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// Pong is the reply to ping.
	Pong struct {
		// This field should mirror the UDP envelope address
		// of the ping packet, which provides a way to discover the
		// the external address (after NAT).
		To         Endpoint
		ReplyTok   []byte // This contains the hash of the ping packet.
		Expiration uint64 // Absolute timestamp at which the packet becomes invalid.
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// Findnode is a query for nodes close to the given target.
	Findnode[P crypto.PublicKey] struct {
		Target     Pubkey[P]
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// Neighbors is the reply to findnode.
	Neighbors[P crypto.PublicKey] struct {
		Nodes      []Node[P]
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// enrRequest queries for the remote node's record.
	ENRRequest struct {
		Expiration uint64
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}

	// enrResponse is the reply to enrRequest.
	ENRResponse struct {
		ReplyTok []byte // Hash of the enrRequest packet.
		Record   enr.Record
		// Ignore additional fields (for forward compatibility).
		Rest []rlp.RawValue `rlp:"tail"`
	}
)

// This number is the maximum number of neighbor nodes in a Neigbors packet.
const MaxNeighbors = 12

// This code computes the MaxNeighbors constant value.

// func init() {
// 	var maxNeighbors int
// 	p := Neighbors{Expiration: ^uint64(0)}
// 	maxSizeNode := Node{IP: make(net.IP, 16), UDP: ^uint16(0), TCP: ^uint16(0)}
// 	for n := 0; ; n++ {
// 		p.Nodes = append(p.Nodes, maxSizeNode)
// 		size, _, err := rlp.EncodeToReader(p)
// 		if err != nil {
// 			// If this ever happens, it will be caught by the unit tests.
// 			panic("cannot encode: " + err.Error())
// 		}
// 		if headSize+size+1 >= 1280 {
// 			maxNeighbors = n
// 			break
// 		}
// 	}
// 	fmt.Println("maxNeighbors", maxNeighbors)
// }

// Pubkey represents an encoded 64-byte public key.
type Pubkey[P crypto.PublicKey] [64]byte

// ID returns the node ID corresponding to the public key.
func (e Pubkey[P]) ID() enode.ID {
	return enode.ID(crypto.Keccak256Hash[P](e[:]))
}

// Node represents information about a node.
type Node[P crypto.PublicKey] struct {
	IP  net.IP // len 4 for IPv4 or 16 for IPv6
	UDP uint16 // for discovery protocol
	TCP uint16 // for RLPx protocol
	ID  Pubkey[P]
}

// Endpoint represents a network endpoint.
type Endpoint struct {
	IP  net.IP // len 4 for IPv4 or 16 for IPv6
	UDP uint16 // for discovery protocol
	TCP uint16 // for RLPx protocol
}

// NewEndpoint creates an endpoint.
func NewEndpoint(addr *net.UDPAddr, tcpPort uint16) Endpoint {
	ip := net.IP{}
	if ip4 := addr.IP.To4(); ip4 != nil {
		ip = ip4
	} else if ip6 := addr.IP.To16(); ip6 != nil {
		ip = ip6
	}
	return Endpoint{IP: ip, UDP: uint16(addr.Port), TCP: tcpPort}
}

type Packet interface {
	// packet name and type for logging purposes.
	Name() string
	Kind() byte
}

func (req *Ping) Name() string   { return "PING/v4" }
func (req *Ping) Kind() byte     { return PingPacket }
func (req *Ping) ENRSeq() uint64 { return seqFromTail(req.Rest) }

func (req *Pong) Name() string   { return "PONG/v4" }
func (req *Pong) Kind() byte     { return PongPacket }
func (req *Pong) ENRSeq() uint64 { return seqFromTail(req.Rest) }

func (req *Findnode[P]) Name() string { return "FINDNODE/v4" }
func (req *Findnode[P]) Kind() byte   { return FindnodePacket }

func (req *Neighbors[P]) Name() string { return "NEIGHBORS/v4" }
func (req *Neighbors[P]) Kind() byte   { return NeighborsPacket }

func (req *ENRRequest) Name() string { return "ENRREQUEST/v4" }
func (req *ENRRequest) Kind() byte   { return ENRRequestPacket }

func (req *ENRResponse) Name() string { return "ENRRESPONSE/v4" }
func (req *ENRResponse) Kind() byte   { return ENRResponsePacket }

// Expired checks whether the given UNIX time stamp is in the past.
func Expired(ts uint64) bool {
	return time.Unix(int64(ts), 0).Before(time.Now())
}

func seqFromTail(tail []rlp.RawValue) uint64 {
	if len(tail) == 0 {
		return 0
	}
	var seq uint64
	rlp.DecodeBytes(tail[0], &seq)
	return seq
}

// Encoder/decoder.

const (
	macSize  = 32
	sigSize  = crypto.SignatureLength
	headSize = macSize + sigSize // space of packet frame data
)

var (
	ErrPacketTooSmall = errors.New("too small")
	ErrBadHash        = errors.New("bad hash")
	ErrBadPoint       = errors.New("invalid curve point")
)

var headSpace = make([]byte, headSize)

// Decode reads a discovery v4 packet.
func Decode[P crypto.PublicKey](input []byte) (Packet, Pubkey[P], []byte, error) {
	if len(input) < headSize+1 {
		return nil, Pubkey[P]{}, nil, ErrPacketTooSmall
	}
	hash, sig, sigdata := input[:macSize], input[macSize:headSize], input[headSize:]
	shouldhash := crypto.Keccak256[P](input[macSize:])
	if !bytes.Equal(hash, shouldhash) {
		return nil, Pubkey[P]{}, nil, ErrBadHash
	}
	fromKey, err := recoverNodeKey[P](crypto.Keccak256[P](input[headSize:]), sig)
	if err != nil {
		return nil, fromKey, hash, err
	}

	var req Packet
	switch ptype := sigdata[0]; ptype {
	case PingPacket:
		req = new(Ping)
	case PongPacket:
		req = new(Pong)
	case FindnodePacket:
		req = new(Findnode[P])
	case NeighborsPacket:
		req = new(Neighbors[P])
	case ENRRequestPacket:
		req = new(ENRRequest)
	case ENRResponsePacket:
		req = new(ENRResponse)
	default:
		return nil, fromKey, hash, fmt.Errorf("unknown type: %d", ptype)
	}
	s := rlp.NewStream(bytes.NewReader(sigdata[1:]), 0)
	err = s.Decode(req)
	return req, fromKey, hash, err
}

// Encode encodes a discovery packet.
func Encode[T crypto.PrivateKey, P crypto.PublicKey](priv T, req Packet) (packet, hash []byte, err error) {
	b := new(bytes.Buffer)
	b.Write(headSpace)
	b.WriteByte(req.Kind())
	if err := rlp.Encode(b, req); err != nil {
		return nil, nil, err
	}
	packet = b.Bytes()
	sig, err := crypto.Sign(crypto.Keccak256[P](packet[headSize:]), priv)
	if err != nil {
		return nil, nil, err
	}
	copy(packet[macSize:], sig)
	// Add the hash to the front. Note: this doesn't protect the packet in any way.
	hash = crypto.Keccak256[P](packet[macSize:])
	copy(packet, hash)
	return packet, hash, nil
}

// recoverNodeKey computes the public key used to sign the given hash from the signature.
func recoverNodeKey[P crypto.PublicKey](hash, sig []byte) (key Pubkey[P], err error) {
	pubkey, err := crypto.Ecrecover[P](hash, sig)
	if err != nil {
		return key, err
	}
	copy(key[:], pubkey[1:])
	return key, nil
}

// EncodePubkey encodes a secp256k1 public key.
func EncodePubkey[P crypto.PublicKey](key P) Pubkey[P] {
	switch p := any(&key).(type) {
	case *nist.PublicKey:
		var e Pubkey[P]
		math.ReadBits(p.X, e[:len(e)/2])
		math.ReadBits(p.Y, e[len(e)/2:])
		return e
	case *gost3410.PublicKey:
		var e Pubkey[P]
		math.ReadBits(p.X, e[:len(e)/2])
		math.ReadBits(p.Y, e[len(e)/2:])
		return e
	case *csp.PublicKey:
		var e Pubkey[P]
		math.ReadBits(p.X, e[:len(e)/2])
		math.ReadBits(p.Y, e[len(e)/2:])
		return e
	default:
		panic("cant encode pub key")
	}
}

// DecodePubkey reads an encoded secp256k1 public key.
func DecodePubkey[P crypto.PublicKey](e Pubkey[P]) (P, error) {
	var pub P
	switch p := any(&pub).(type) {
	case *nist.PublicKey:
		k := &ecdsa.PublicKey{Curve: crypto.S256(), X: new(big.Int), Y: new(big.Int)}
		half := len(e) / 2
		k.X.SetBytes(e[:half])
		k.Y.SetBytes(e[half:])
		if !k.Curve.IsOnCurve(k.X, k.Y) {
			return crypto.ZeroPublicKey[P](), ErrBadPoint
		}
		*p = nist.PublicKey{k}
	case *gost3410.PublicKey:
		k := &gost3410.PublicKey{C: gost3410.GostCurve, X: new(big.Int), Y: new(big.Int)}
		half := len(e) / 2
		// crypto.Reverse(e[:])
		k.X.SetBytes(e[:half])
		k.Y.SetBytes(e[half:])
		// if !k.C.IsOnCurve(k.X, k.Y) {
		// 	return crypto.ZeroPublicKey[P](), ErrBadPoint
		// }
		*p = *k
	case *csp.PublicKey:
		k := &csp.PublicKey{Curve: gost3410.CurveIdGostR34102001CryptoProAParamSet(), X: new(big.Int), Y: new(big.Int)}
		half := len(e) / 2
		k.X.SetBytes(e[:half])
		k.Y.SetBytes(e[half:])
		if !k.Curve.IsOnCurve(k.X, k.Y) {
			return crypto.ZeroPublicKey[P](), ErrBadPoint
		}
		*p = *k
	default:
		return crypto.ZeroPublicKey[P](), fmt.Errorf("cant infer public key")
	}
	return pub, nil
}
