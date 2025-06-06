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

package discover

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"sync"

	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
)

var nullNode *enode.Node[nist.PublicKey]

func init() {
	var r enr.Record
	r.Set(enr.IP{0, 0, 0, 0})
	nullNode = enode.SignNull[nist.PrivateKey,nist.PublicKey](&r, enode.ID{})
}

func newTestTable(t transport[nist.PublicKey]) (*Table[nist.PublicKey], *enode.DB[nist.PublicKey]) {
	db, _ := enode.OpenDB[nist.PublicKey]("")
	tab, _ := newTable(t, db, nil, log.Root())
	go tab.loop()
	return tab, db
}

// nodeAtDistance creates a node for which enode.LogDist(base, n.id) == ld.
func nodeAtDistance[T crypto.PrivateKey, P crypto.PublicKey](base enode.ID, ld int, ip net.IP) *node[P] {
	var r enr.Record
	r.Set(enr.IP(ip))
	return wrapNode(enode.SignNull[T,P](&r, idAtDistance(base, ld)))
}

// nodesAtDistance creates n nodes for which enode.LogDist(base, node.ID()) == ld.
func nodesAtDistance[T crypto.PrivateKey, P crypto.PublicKey](base enode.ID, ld int, n int) []*enode.Node[P] {
	results := make([]*enode.Node[P], n)
	for i := range results {
		results[i] = unwrapNode(nodeAtDistance[T,P](base, ld, intIP(i)))
	}
	return results
}

func nodesToRecords(nodes []*enode.Node[nist.PublicKey]) []*enr.Record {
	records := make([]*enr.Record, len(nodes))
	for i := range nodes {
		records[i] = nodes[i].Record()
	}
	return records
}

// idAtDistance returns a random hash such that enode.LogDist(a, b) == n
func idAtDistance(a enode.ID, n int) (b enode.ID) {
	if n == 0 {
		return a
	}
	// flip bit at position n, fill the rest with random bits
	b = a
	pos := len(a) - n/8 - 1
	bit := byte(0x01) << (byte(n%8) - 1)
	if bit == 0 {
		pos++
		bit = 0x80
	}
	b[pos] = a[pos]&^bit | ^a[pos]&bit // TODO: randomize end bits
	for i := pos + 1; i < len(a); i++ {
		b[i] = byte(rand.Intn(255))
	}
	return b
}

func intIP(i int) net.IP {
	return net.IP{byte(i), 0, 2, byte(i)}
}

// fillBucket inserts nodes into the given bucket until it is full.
func fillBucket[T crypto.PrivateKey, P crypto.PublicKey](tab *Table[P], n *node[P]) (last *node[P]) {
	ld := enode.LogDist(tab.self().ID(), n.ID())
	b := tab.bucket(n.ID())
	for len(b.entries) < bucketSize {
		b.entries = append(b.entries, nodeAtDistance[T,P](tab.self().ID(), ld, intIP(ld)))
	}
	return b.entries[bucketSize-1]
}

// fillTable adds nodes the table to the end of their corresponding bucket
// if the bucket is not full. The caller must not hold tab.mutex.
func fillTable[P crypto.PublicKey](tab *Table[P], nodes []*node[P]) {
	for _, n := range nodes {
		tab.addSeenNode(n)
	}
}

type pingRecorder struct {
	mu           sync.Mutex
	dead, pinged map[enode.ID]bool
	records      map[enode.ID]*enode.Node[nist.PublicKey]
	n            *enode.Node[nist.PublicKey]
}

func newPingRecorder() *pingRecorder {
	var r enr.Record
	r.Set(enr.IP{0, 0, 0, 0})
	n := enode.SignNull[nist.PrivateKey,nist.PublicKey](&r, enode.ID{})

	return &pingRecorder{
		dead:    make(map[enode.ID]bool),
		pinged:  make(map[enode.ID]bool),
		records: make(map[enode.ID]*enode.Node[nist.PublicKey]),
		n:       n,
	}
}

// setRecord updates a node record. Future calls to ping and
// requestENR will return this record.
func (t *pingRecorder) updateRecord(n *enode.Node[nist.PublicKey]) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.records[n.ID()] = n
}

// Stubs to satisfy the transport interface.
func (t *pingRecorder) Self() *enode.Node[nist.PublicKey]           { return nullNode }
func (t *pingRecorder) lookupSelf() []*enode.Node[nist.PublicKey]   { return nil }
func (t *pingRecorder) lookupRandom() []*enode.Node[nist.PublicKey] { return nil }

// ping simulates a ping request.
func (t *pingRecorder) ping(n *enode.Node[nist.PublicKey]) (seq uint64, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.pinged[n.ID()] = true
	if t.dead[n.ID()] {
		return 0, errTimeout
	}
	if t.records[n.ID()] != nil {
		seq = t.records[n.ID()].Seq()
	}
	return seq, nil
}

// requestENR simulates an ENR request.
func (t *pingRecorder) RequestENR(n *enode.Node[nist.PublicKey]) (*enode.Node[nist.PublicKey], error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.dead[n.ID()] || t.records[n.ID()] == nil {
		return nil, errTimeout
	}
	return t.records[n.ID()], nil
}

func hasDuplicates[P crypto.PublicKey] (slice []*node[P]) bool {
	seen := make(map[enode.ID]bool)
	for i, e := range slice {
		if e == nil {
			panic(fmt.Sprintf("nil *Node at %d", i))
		}
		if seen[e.ID()] {
			return true
		}
		seen[e.ID()] = true
	}
	return false
}

// checkNodesEqual checks whether the two given node lists contain the same nodes.
func checkNodesEqual[P crypto.PublicKey](got, want []*enode.Node[P]) error {
	if len(got) == len(want) {
		for i := range got {
			if !nodeEqual(got[i], want[i]) {
				goto NotEqual
			}
		}
	}
	return nil

NotEqual:
	output := new(bytes.Buffer)
	fmt.Fprintf(output, "got %d nodes:\n", len(got))
	for _, n := range got {
		fmt.Fprintf(output, "  %v %v\n", n.ID(), n)
	}
	fmt.Fprintf(output, "want %d:\n", len(want))
	for _, n := range want {
		fmt.Fprintf(output, "  %v %v\n", n.ID(), n)
	}
	return errors.New(output.String())
}

func nodeEqual[P crypto.PublicKey](n1 *enode.Node[P], n2 *enode.Node[P]) bool {
	return n1.ID() == n2.ID() && n1.IP().Equal(n2.IP())
}

func sortByID[P crypto.PublicKey](nodes []*enode.Node[P]) {
	sort.Slice(nodes, func(i, j int) bool {
		return string(nodes[i].ID().Bytes()) < string(nodes[j].ID().Bytes())
	})
}

func sortedByDistanceTo[P crypto.PublicKey](distbase enode.ID, slice []*node[P]) bool {
	return sort.SliceIsSorted(slice, func(i, j int) bool {
		return enode.DistCmp(distbase, slice[i].ID(), slice[j].ID()) < 0
	})
}

// hexEncPrivkey decodes h as a private key.
func hexEncPrivkey[T crypto.PrivateKey](h string) T {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	key, err := crypto.ToECDSA[T](b)
	if err != nil {
		panic(err)
	}
	return key
}

// hexEncPubkey decodes h as a public key.
func hexEncPubkey[P crypto.PublicKey] (h string) (ret encPubkey[P]) {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	if len(b) != len(ret) {
		panic("invalid length")
	}
	copy(ret[:], b)
	return ret
}
