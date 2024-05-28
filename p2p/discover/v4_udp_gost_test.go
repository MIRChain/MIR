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
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/p2p/discover/v4wire"
	"github.com/MIRChain/MIR/p2p/enode"
	"github.com/MIRChain/MIR/p2p/enr"
)

var testTargetGost = v4wire.Pubkey[gost3410.PublicKey]{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}

func TestUDPv4_packetErrorsGost(t *testing.T) {
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	test.packetIn(errExpired, &v4wire.Ping{From: testRemote, To: testLocalAnnounced, Version: 4})
	test.packetIn(errUnsolicitedReply, &v4wire.Pong{ReplyTok: []byte{}, Expiration: futureExp})
	test.packetIn(errUnknownNode, &v4wire.Findnode[gost3410.PublicKey]{Expiration: futureExp})
	test.packetIn(errUnsolicitedReply, &v4wire.Neighbors[gost3410.PublicKey]{Expiration: futureExp})
}

func TestUDPv4_pingTimeoutGost(t *testing.T) {
	t.Parallel()
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	key := newkey[gost3410.PrivateKey]()
	toaddr := &net.UDPAddr{IP: net.ParseIP("1.2.3.4"), Port: 2222}
	node := enode.NewV4(*key.Public(), toaddr.IP, 0, toaddr.Port)
	if _, err := test.udp.ping(node); err != errTimeout {
		t.Error("expected timeout error, got", err)
	}
}

func TestUDPv4_responseTimeoutsGost(t *testing.T) {
	t.Parallel()
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	rand.Seed(time.Now().UnixNano())
	randomDuration := func(max time.Duration) time.Duration {
		return time.Duration(rand.Int63n(int64(max)))
	}

	var (
		nReqs      = 200
		nTimeouts  = 0                       // number of requests with ptype > 128
		nilErr     = make(chan error, nReqs) // for requests that get a reply
		timeoutErr = make(chan error, nReqs) // for requests that time out
	)
	for i := 0; i < nReqs; i++ {
		// Create a matcher for a random request in udp.loop. Requests
		// with ptype <= 128 will not get a reply and should time out.
		// For all other requests, a reply is scheduled to arrive
		// within the timeout window.
		p := &replyMatcher{
			ptype:    byte(rand.Intn(255)),
			callback: func(v4wire.Packet) (bool, bool) { return true, true },
		}
		binary.BigEndian.PutUint64(p.from[:], uint64(i))
		if p.ptype <= 128 {
			p.errc = timeoutErr
			test.udp.addReplyMatcher <- p
			nTimeouts++
		} else {
			p.errc = nilErr
			test.udp.addReplyMatcher <- p
			time.AfterFunc(randomDuration(60*time.Millisecond), func() {
				if !test.udp.handleReply(p.from, p.ip, testPacket(p.ptype)) {
					t.Logf("not matched: %v", p)
				}
			})
		}
		time.Sleep(randomDuration(30 * time.Millisecond))
	}

	// Check that all timeouts were delivered and that the rest got nil errors.
	// The replies must be delivered.
	var (
		recvDeadline        = time.After(20 * time.Second)
		nTimeoutsRecv, nNil = 0, 0
	)
	for i := 0; i < nReqs; i++ {
		select {
		case err := <-timeoutErr:
			if err != errTimeout {
				t.Fatalf("got non-timeout error on timeoutErr %d: %v", i, err)
			}
			nTimeoutsRecv++
		case err := <-nilErr:
			if err != nil {
				t.Fatalf("got non-nil error on nilErr %d: %v", i, err)
			}
			nNil++
		case <-recvDeadline:
			t.Fatalf("exceeded recv deadline")
		}
	}
	if nTimeoutsRecv != nTimeouts {
		t.Errorf("wrong number of timeout errors received: got %d, want %d", nTimeoutsRecv, nTimeouts)
	}
	if nNil != nReqs-nTimeouts {
		t.Errorf("wrong number of successful replies: got %d, want %d", nNil, nReqs-nTimeouts)
	}
}

func TestUDPv4_findnodeTimeoutGost(t *testing.T) {
	t.Parallel()
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	toaddr := &net.UDPAddr{IP: net.ParseIP("1.2.3.4"), Port: 2222}
	toid := enode.ID{1, 2, 3, 4}
	target := v4wire.Pubkey[gost3410.PublicKey]{4, 5, 6, 7}
	result, err := test.udp.findnode(toid, toaddr, target)
	if err != errTimeout {
		t.Error("expected timeout error, got", err)
	}
	if len(result) > 0 {
		t.Error("expected empty result, got", result)
	}
}

func TestUDPv4_findnodeGost(t *testing.T) {
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	// put a few nodes into the table. their exact
	// distribution shouldn't matter much, although we need to
	// take care not to overflow any bucket.
	nodes := &nodesByDistance[gost3410.PublicKey]{target: testTargetGost.ID()}
	live := make(map[enode.ID]bool)
	numCandidates := 2 * bucketSize
	for i := 0; i < numCandidates; i++ {
		key := newkey[gost3410.PrivateKey]()
		ip := net.IP{10, 13, 0, byte(i)}
		n := wrapNode(enode.NewV4(*key.Public(), ip, 0, 2000))
		// Ensure half of table content isn't verified live yet.
		if i > numCandidates/2 {
			n.livenessChecks = 1
			live[n.ID()] = true
		}
		nodes.push(n, numCandidates)
	}
	fillTable[gost3410.PublicKey](test.table, nodes.entries)

	// ensure there's a bond with the test node,
	// findnode won't be accepted otherwise.
	remoteID := v4wire.EncodePubkey(*test.remotekey.Public()).ID()
	test.table.db.UpdateLastPongReceived(remoteID, test.remoteaddr.IP, time.Now())

	// check that closest neighbors are returned.
	expected := test.table.findnodeByID(testTargetGost.ID(), bucketSize, true)
	test.packetIn(nil, &v4wire.Findnode[gost3410.PublicKey]{Target: testTargetGost, Expiration: futureExp})
	waitNeighbors := func(want []*node[gost3410.PublicKey]) {
		test.waitPacketOut(func(p *v4wire.Neighbors[gost3410.PublicKey], to *net.UDPAddr, hash []byte) {
			if len(p.Nodes) != len(want) {
				t.Errorf("wrong number of results: got %d, want %d", len(p.Nodes), bucketSize)
			}
			for i, n := range p.Nodes {
				if n.ID.ID() != want[i].ID() {
					t.Errorf("result mismatch at %d:\n  got:  %v\n  want: %v", i, n, expected.entries[i])
				}
				if !live[n.ID.ID()] {
					t.Errorf("result includes dead node %v", n.ID.ID())
				}
			}
		})
	}
	// Receive replies.
	want := expected.entries
	if len(want) > v4wire.MaxNeighbors {
		waitNeighbors(want[:v4wire.MaxNeighbors])
		want = want[v4wire.MaxNeighbors:]
	}
	waitNeighbors(want)
}

func TestUDPv4_findnodeMultiReplyGost(t *testing.T) {
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	rid := enode.PubkeyToIDV4(*test.remotekey.Public())
	test.table.db.UpdateLastPingReceived(rid, test.remoteaddr.IP, time.Now())

	// queue a pending findnode request
	resultc, errc := make(chan []*node[gost3410.PublicKey]), make(chan error)
	go func() {
		rid := encodePubkey(*test.remotekey.Public()).id()
		ns, err := test.udp.findnode(rid, test.remoteaddr, testTargetGost)
		if err != nil && len(ns) == 0 {
			errc <- err
		} else {
			resultc <- ns
		}
	}()

	// wait for the findnode to be sent.
	// after it is sent, the transport is waiting for a reply
	test.waitPacketOut(func(p *v4wire.Findnode[gost3410.PublicKey], to *net.UDPAddr, hash []byte) {
		if p.Target != testTargetGost {
			t.Errorf("wrong target: got %v, want %v", p.Target, testTargetGost)
		}
	})

	// send the reply as two packets.
	list := []*node[gost3410.PublicKey]{
		wrapNode(enode.MustParse[gost3410.PublicKey]("enode://ba6303fbb832ef2a24a5a2b8603bc18981073dd7634b9c6454aa8183ccf07c3b60aabbca2fd8601a29958dd7d7799eb725051d25244635417b9f5a34f1c29a28@10.0.1.16:30303?discport=30304")),
		wrapNode(enode.MustParse[gost3410.PublicKey]("enode://6143426d8cda5e8182def4a679ae4415f371d428d53dd6ee63f841b58791dff23149a8c5e029ca466f0e1c0829d4499e959dd30e08fe36fedc84db0f7f67b3ec@10.0.1.16:30303")),
		wrapNode(enode.MustParse[gost3410.PublicKey]("enode://152738b84a9c9b02cbd0be3cca8f89533cff5b4c81819d122f6c62e03713afd90add6783c4d92aa4b9530172c17a07c8b337771db37be3f2ab230365050476e7@10.0.1.36:30301?discport=17")),
		wrapNode(enode.MustParse[gost3410.PublicKey]("enode://161df3a5f868dc64c0f778b2fb3c3724311b9b7f4fdec9eccac3f9a40e065ffb04932ee860d7224fc2972db93ee8a4e43d50e0dca79953ebb0c849d5083f9ab3@10.0.1.16:30303")),
	}
	rpclist := make([]v4wire.Node[gost3410.PublicKey], len(list))
	for i := range list {
		rpclist[i] = nodeToRPC(list[i])
	}
	test.packetIn(nil, &v4wire.Neighbors[gost3410.PublicKey]{Expiration: futureExp, Nodes: rpclist[:2]})
	test.packetIn(nil, &v4wire.Neighbors[gost3410.PublicKey]{Expiration: futureExp, Nodes: rpclist[2:]})

	// check that the sent neighbors are all returned by findnode
	select {
	case result := <-resultc:
		want := append(list[:2], list[3:]...)
		if !reflect.DeepEqual(result, want) {
			t.Errorf("neighbors mismatch:\n  got:  %v\n  want: %v", result, want)
		}
	case err := <-errc:
		t.Errorf("findnode error: %v", err)
	case <-time.After(5 * time.Second):
		t.Error("findnode did not return within 5 seconds")
	}
}

// This test checks that reply matching of pong verifies the ping hash.
func TestUDPv4_pingMatchGost(t *testing.T) {
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	randToken := make([]byte, 32)
	crand.Read(randToken)

	test.packetIn(nil, &v4wire.Ping{From: testRemote, To: testLocalAnnounced, Version: 4, Expiration: futureExp})
	test.waitPacketOut(func(*v4wire.Pong, *net.UDPAddr, []byte) {})
	test.waitPacketOut(func(*v4wire.Ping, *net.UDPAddr, []byte) {})
	test.packetIn(errUnsolicitedReply, &v4wire.Pong{ReplyTok: randToken, To: testLocalAnnounced, Expiration: futureExp})
}

// This test checks that reply matching of pong verifies the sender IP address.
func TestUDPv4_pingMatchIPGost(t *testing.T) {
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	test.packetIn(nil, &v4wire.Ping{From: testRemote, To: testLocalAnnounced, Version: 4, Expiration: futureExp})
	test.waitPacketOut(func(*v4wire.Pong, *net.UDPAddr, []byte) {})

	test.waitPacketOut(func(p *v4wire.Ping, to *net.UDPAddr, hash []byte) {
		wrongAddr := &net.UDPAddr{IP: net.IP{33, 44, 1, 2}, Port: 30000}
		test.packetInFrom(errUnsolicitedReply, test.remotekey, wrongAddr, &v4wire.Pong{
			ReplyTok:   hash,
			To:         testLocalAnnounced,
			Expiration: futureExp,
		})
	})
}

func TestUDPv4_successfulPingGost(t *testing.T) {
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	added := make(chan *node[gost3410.PublicKey], 1)
	test.table.nodeAddedHook = func(n *node[gost3410.PublicKey]) { added <- n }
	defer test.close()

	// The remote side sends a ping packet to initiate the exchange.
	go test.packetIn(nil, &v4wire.Ping{From: testRemote, To: testLocalAnnounced, Version: 4, Expiration: futureExp})

	// The ping is replied to.
	test.waitPacketOut(func(p *v4wire.Pong, to *net.UDPAddr, hash []byte) {
		pinghash := test.sent[0][:32]
		if !bytes.Equal(p.ReplyTok, pinghash) {
			t.Errorf("got pong.ReplyTok %x, want %x", p.ReplyTok, pinghash)
		}
		wantTo := v4wire.Endpoint{
			// The mirrored UDP address is the UDP packet sender
			IP: test.remoteaddr.IP, UDP: uint16(test.remoteaddr.Port),
			// The mirrored TCP port is the one from the ping packet
			TCP: testRemote.TCP,
		}
		if !reflect.DeepEqual(p.To, wantTo) {
			t.Errorf("got pong.To %v, want %v", p.To, wantTo)
		}
	})

	// Remote is unknown, the table pings back.
	test.waitPacketOut(func(p *v4wire.Ping, to *net.UDPAddr, hash []byte) {
		if !reflect.DeepEqual(p.From, test.udp.ourEndpoint()) {
			t.Errorf("got ping.From %#v, want %#v", p.From, test.udp.ourEndpoint())
		}
		wantTo := v4wire.Endpoint{
			// The mirrored UDP address is the UDP packet sender.
			IP:  test.remoteaddr.IP,
			UDP: uint16(test.remoteaddr.Port),
			TCP: 0,
		}
		if !reflect.DeepEqual(p.To, wantTo) {
			t.Errorf("got ping.To %v, want %v", p.To, wantTo)
		}
		test.packetIn(nil, &v4wire.Pong{ReplyTok: hash, Expiration: futureExp})
	})

	// The node should be added to the table shortly after getting the
	// pong packet.
	select {
	case n := <-added:
		rid := encodePubkey(*test.remotekey.Public()).id()
		if n.ID() != rid {
			t.Errorf("node has wrong ID: got %v, want %v", n.ID(), rid)
		}
		if !n.IP().Equal(test.remoteaddr.IP) {
			t.Errorf("node has wrong IP: got %v, want: %v", n.IP(), test.remoteaddr.IP)
		}
		if n.UDP() != test.remoteaddr.Port {
			t.Errorf("node has wrong UDP port: got %v, want: %v", n.UDP(), test.remoteaddr.Port)
		}
		if n.TCP() != int(testRemote.TCP) {
			t.Errorf("node has wrong TCP port: got %v, want: %v", n.TCP(), testRemote.TCP)
		}
	case <-time.After(2 * time.Second):
		t.Errorf("node was not added within 2 seconds")
	}
}

// This test checks that EIP-868 requests work.
func TestUDPv4_EIP868Gost(t *testing.T) {
	test := newUDPTest[gost3410.PrivateKey, gost3410.PublicKey](t)
	defer test.close()

	test.udp.localNode.Set(enr.WithEntry("foo", "bar"))
	wantNode := test.udp.localNode.Node()

	// ENR requests aren't allowed before endpoint proof.
	test.packetIn(errUnknownNode, &v4wire.ENRRequest{Expiration: futureExp})

	// Perform endpoint proof and check for sequence number in packet tail.
	test.packetIn(nil, &v4wire.Ping{Expiration: futureExp})
	test.waitPacketOut(func(p *v4wire.Pong, addr *net.UDPAddr, hash []byte) {
		if p.ENRSeq() != wantNode.Seq() {
			t.Errorf("wrong sequence number in pong: %d, want %d", p.ENRSeq(), wantNode.Seq())
		}
	})
	test.waitPacketOut(func(p *v4wire.Ping, addr *net.UDPAddr, hash []byte) {
		if p.ENRSeq() != wantNode.Seq() {
			t.Errorf("wrong sequence number in ping: %d, want %d", p.ENRSeq(), wantNode.Seq())
		}
		test.packetIn(nil, &v4wire.Pong{Expiration: futureExp, ReplyTok: hash})
	})

	// Request should work now.
	test.packetIn(nil, &v4wire.ENRRequest{Expiration: futureExp})
	test.waitPacketOut(func(p *v4wire.ENRResponse, addr *net.UDPAddr, hash []byte) {
		n, err := enode.New[gost3410.PublicKey](enr.SchemeMap{"v4": enode.V4ID[gost3410.PublicKey]{}}, &p.Record)
		if err != nil {
			t.Fatalf("invalid record: %v", err)
		}
		if !reflect.DeepEqual(n, wantNode) {
			t.Fatalf("wrong node in enrResponse: %v", n)
		}
	})
}

// This test verifies that a small network of nodes can boot up into a healthy state.
func TestUDPv4_smallNetConvergenceGost(t *testing.T) {
	t.Parallel()

	// Start the network.
	nodes := make([]*UDPv4[gost3410.PrivateKey, gost3410.PublicKey], 4)
	for i := range nodes {
		var cfg Config[gost3410.PrivateKey, gost3410.PublicKey]
		if i > 0 {
			bn := nodes[0].Self()
			cfg.Bootnodes = []*enode.Node[gost3410.PublicKey]{bn}
		}
		nodes[i] = startLocalhostV4(t, cfg)
		defer nodes[i].Close()
	}

	// Run through the iterator on all nodes until
	// they have all found each other.
	status := make(chan error, len(nodes))
	for i := range nodes {
		node := nodes[i]
		go func() {
			found := make(map[enode.ID]bool, len(nodes))
			it := node.RandomNodes()
			for it.Next() {
				found[it.Node().ID()] = true
				if len(found) == len(nodes) {
					status <- nil
					return
				}
			}
			status <- fmt.Errorf("node %s didn't find all nodes", node.Self().ID().TerminalString())
		}()
	}

	// Wait for all status reports.
	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()
	for received := 0; received < len(nodes); {
		select {
		case <-timeout.C:
			for _, node := range nodes {
				node.Close()
			}
		case err := <-status:
			received++
			if err != nil {
				t.Error("ERROR:", err)
				return
			}
		}
	}
}
