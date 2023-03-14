// Copyright 2015 The go-ethereum Authors
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

package discover

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/csp"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
)

// node represents a host on the network.
// The fields of Node may not be modified.
type node [P crypto.PublicKey] struct {
	enode.Node[P]
	addedAt        time.Time // time when the node was added to the table
	livenessChecks uint      // how often liveness was checked
}

type encPubkey [64]byte

func encodePubkey [P crypto.PublicKey] (key P) encPubkey {
	var e encPubkey 
	math.ReadBits(key.GetX(), e[:len(e)/2])
	math.ReadBits(key.GetY(), e[len(e)/2:])
	return e
}

func decodePubkey[P crypto.PublicKey](e []byte) (P, error) {
	var pub P
	switch p:= any(&pub).(type){
	case *nist.PublicKey:
		if len(e) != len(encPubkey{}) {
			return crypto.ZeroPublicKey[P](), errors.New("wrong size public key data")
		}
		k := &ecdsa.PublicKey{Curve: crypto.S256(), X: new(big.Int), Y: new(big.Int)}
		half := len(e) / 2
		k.X.SetBytes(e[:half])
		k.Y.SetBytes(e[half:])
		if !p.Curve.IsOnCurve(p.X, p.Y) {
			return crypto.ZeroPublicKey[P](), errors.New("invalid curve point")
		}
		*p = nist.PublicKey{k}
	case *gost3410.PublicKey:
		if len(e) != len(encPubkey{}) {
			return crypto.ZeroPublicKey[P](), errors.New("wrong size public key data")
		}
		k := &gost3410.PublicKey{C: gost3410.GostCurve, X: new(big.Int), Y: new(big.Int)}
		half := len(e) / 2
		k.X.SetBytes(e[:half])
		k.Y.SetBytes(e[half:])
		if !p.C.IsOnCurve(p.X, p.Y) {
			return crypto.ZeroPublicKey[P](), errors.New("invalid curve point")
		}
		*p = *k
	case *csp.PublicKey:
		if len(e) != len(encPubkey{}) {
			return crypto.ZeroPublicKey[P](), errors.New("wrong size public key data")
		}
		k := &csp.PublicKey{Curve: gost3410.CurveIdGostR34102001CryptoProAParamSet(), X: new(big.Int), Y: new(big.Int)}
		half := len(e) / 2
		k.X.SetBytes(e[:half])
		k.Y.SetBytes(e[half:])
		if !p.Curve.IsOnCurve(p.X, p.Y) {
			return crypto.ZeroPublicKey[P](), errors.New("invalid curve point")
		}
		*p = *k
	default:
		return crypto.ZeroPublicKey[P](), fmt.Errorf("cant infer public key")
	}
	return pub, nil
}

func (e encPubkey) id() enode.ID {
	return enode.ID(crypto.Keccak256Hash(e[:]))
}

func wrapNode[P crypto.PublicKey](n *enode.Node[P]) *node[P] {
	return &node[P]{Node: *n}
}

func wrapNodes[P crypto.PublicKey](ns []*enode.Node[P]) []*node[P] {
	result := make([]*node[P], len(ns))
	for i, n := range ns {
		result[i] = wrapNode[P](n)
	}
	return result
}

func unwrapNode[P crypto.PublicKey](n *node[P]) *enode.Node[P] {
	return &n.Node
}

func unwrapNodes[P crypto.PublicKey](ns []*node[P]) []*enode.Node[P] {
	result := make([]*enode.Node[P], len(ns))
	for i, n := range ns {
		result[i] = unwrapNode(n)
	}
	return result
}

func (n *node[P]) addr() *net.UDPAddr {
	return &net.UDPAddr{IP: n.IP(), Port: n.UDP()}
}

func (n *node[P]) String() string {
	return n.Node.String()
}
