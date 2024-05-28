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

package discover

import (
	"net"

	"github.com/MIRChain/MIR/common/mclock"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/log"
	"github.com/MIRChain/MIR/p2p/enode"
	"github.com/MIRChain/MIR/p2p/enr"
	"github.com/MIRChain/MIR/p2p/netutil"
)

// UDPConn is a network connection on which discovery can operate.
type UDPConn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}

// Config holds settings for the discovery listener.
type Config[T crypto.PrivateKey, P crypto.PublicKey] struct {
	// These settings are required and configure the UDP listener:
	PrivateKey T

	// These settings are optional:
	NetRestrict  *netutil.Netlist   // network whitelist
	Bootnodes    []*enode.Node[P]   // list of bootstrap nodes
	Unhandled    chan<- ReadPacket  // unhandled packets are sent on this channel
	Log          log.Logger         // if set, log messages go here
	ValidSchemes enr.IdentityScheme // allowed identity schemes
	Clock        mclock.Clock
}

func (cfg Config[T, P]) withDefaults() Config[T, P] {
	if cfg.Log == nil {
		cfg.Log = log.Root()
	}
	if cfg.ValidSchemes == nil {
		cfg.ValidSchemes = enr.SchemeMap{"v4": enode.V4ID[P]{}}
	}
	if cfg.Clock == nil {
		cfg.Clock = mclock.System{}
	}
	return cfg
}

// ListenUDP starts listening for discovery packets on the given UDP socket.
func ListenUDP[T crypto.PrivateKey, P crypto.PublicKey](c UDPConn, ln *enode.LocalNode[T, P], cfg Config[T, P]) (*UDPv4[T, P], error) {
	return ListenV4(c, ln, cfg)
}

// ReadPacket is a packet that couldn't be handled. Those packets are sent to the unhandled
// channel if configured.
type ReadPacket struct {
	Data []byte
	Addr *net.UDPAddr
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
