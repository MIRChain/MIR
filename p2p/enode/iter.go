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

package enode

import (
	"sync"
	"time"

	"github.com/MIRChain/MIR/crypto"
)

// Iterator represents a sequence of nodes. The Next method moves to the next node in the
// sequence. It returns false when the sequence has ended or the iterator is closed. Close
// may be called concurrently with Next and Node, and interrupts Next if it is blocked.
type Iterator[P crypto.PublicKey] interface {
	Next() bool     // moves to next node
	Node() *Node[P] // returns current node
	Close()         // ends the iterator
}

// ReadNodes reads at most n nodes from the given iterator. The return value contains no
// duplicates and no nil values. To prevent looping indefinitely for small repeating node
// sequences, this function calls Next at most n times.
func ReadNodes[P crypto.PublicKey](it Iterator[P], n int) []*Node[P] {
	seen := make(map[ID]*Node[P], n)
	for i := 0; i < n && it.Next(); i++ {
		// Remove duplicates, keeping the node with higher seq.
		node := it.Node()
		prevNode, ok := seen[node.ID()]
		if ok && prevNode.Seq() > node.Seq() {
			continue
		}
		seen[node.ID()] = node
	}
	result := make([]*Node[P], 0, len(seen))
	for _, node := range seen {
		result = append(result, node)
	}
	return result
}

// IterNodes makes an iterator which runs through the given nodes once.
func IterNodes[P crypto.PublicKey](nodes []*Node[P]) Iterator[P] {
	return &sliceIter[P]{nodes: nodes, index: -1}
}

// CycleNodes makes an iterator which cycles through the given nodes indefinitely.
func CycleNodes[P crypto.PublicKey](nodes []*Node[P]) Iterator[P] {
	return &sliceIter[P]{nodes: nodes, index: -1, cycle: true}
}

type sliceIter[P crypto.PublicKey] struct {
	mu    sync.Mutex
	nodes []*Node[P]
	index int
	cycle bool
}

func (it *sliceIter[P]) Next() bool {
	it.mu.Lock()
	defer it.mu.Unlock()

	if len(it.nodes) == 0 {
		return false
	}
	it.index++
	if it.index == len(it.nodes) {
		if it.cycle {
			it.index = 0
		} else {
			it.nodes = nil
			return false
		}
	}
	return true
}

func (it *sliceIter[P]) Node() *Node[P] {
	it.mu.Lock()
	defer it.mu.Unlock()
	if len(it.nodes) == 0 {
		return nil
	}
	return it.nodes[it.index]
}

func (it *sliceIter[P]) Close() {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.nodes = nil
}

// Filter wraps an iterator such that Next only returns nodes for which
// the 'check' function returns true.
func Filter[P crypto.PublicKey](it Iterator[P], check func(*Node[P]) bool) Iterator[P] {
	return &filterIter[P]{it, check}
}

type filterIter[P crypto.PublicKey] struct {
	Iterator[P]
	check func(*Node[P]) bool
}

func (f *filterIter[P]) Next() bool {
	for f.Iterator.Next() {
		if f.check(f.Node()) {
			return true
		}
	}
	return false
}

// FairMix aggregates multiple node iterators. The mixer itself is an iterator which ends
// only when Close is called. Source iterators added via AddSource are removed from the
// mix when they end.
//
// The distribution of nodes returned by Next is approximately fair, i.e. FairMix
// attempts to draw from all sources equally often. However, if a certain source is slow
// and doesn't return a node within the configured timeout, a node from any other source
// will be returned.
//
// It's safe to call AddSource and Close concurrently with Next.
type FairMix[P crypto.PublicKey] struct {
	wg      sync.WaitGroup
	fromAny chan *Node[P]
	timeout time.Duration
	cur     *Node[P]

	mu      sync.Mutex
	closed  chan struct{}
	sources []*mixSource[P]
	last    int
}

type mixSource[P crypto.PublicKey] struct {
	it      Iterator[P]
	next    chan *Node[P]
	timeout time.Duration
}

// NewFairMix creates a mixer.
//
// The timeout specifies how long the mixer will wait for the next fairly-chosen source
// before giving up and taking a node from any other source. A good way to set the timeout
// is deciding how long you'd want to wait for a node on average. Passing a negative
// timeout makes the mixer completely fair.
func NewFairMix[P crypto.PublicKey](timeout time.Duration) *FairMix[P] {
	m := &FairMix[P]{
		fromAny: make(chan *Node[P]),
		closed:  make(chan struct{}),
		timeout: timeout,
	}
	return m
}

// AddSource adds a source of nodes.
func (m *FairMix[P]) AddSource(it Iterator[P]) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed == nil {
		return
	}
	m.wg.Add(1)
	source := &mixSource[P]{it, make(chan *Node[P]), m.timeout}
	m.sources = append(m.sources, source)
	go m.runSource(m.closed, source)
}

// Close shuts down the mixer and all current sources.
// Calling this is required to release resources associated with the mixer.
func (m *FairMix[P]) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed == nil {
		return
	}
	for _, s := range m.sources {
		s.it.Close()
	}
	close(m.closed)
	m.wg.Wait()
	close(m.fromAny)
	m.sources = nil
	m.closed = nil
}

// Next returns a node from a random source.
func (m *FairMix[P]) Next() bool {
	m.cur = nil

	var timeout <-chan time.Time
	if m.timeout >= 0 {
		timer := time.NewTimer(m.timeout)
		timeout = timer.C
		defer timer.Stop()
	}
	for {
		source := m.pickSource()
		if source == nil {
			return m.nextFromAny()
		}
		select {
		case n, ok := <-source.next:
			if ok {
				m.cur = n
				source.timeout = m.timeout
				return true
			}
			// This source has ended.
			m.deleteSource(source)
		case <-timeout:
			source.timeout /= 2
			return m.nextFromAny()
		}
	}
}

// Node returns the current node.
func (m *FairMix[P]) Node() *Node[P] {
	return m.cur
}

// nextFromAny is used when there are no sources or when the 'fair' choice
// doesn't turn up a node quickly enough.
func (m *FairMix[P]) nextFromAny() bool {
	n, ok := <-m.fromAny
	if ok {
		m.cur = n
	}
	return ok
}

// pickSource chooses the next source to read from, cycling through them in order.
func (m *FairMix[P]) pickSource() *mixSource[P] {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.sources) == 0 {
		return nil
	}
	m.last = (m.last + 1) % len(m.sources)
	return m.sources[m.last]
}

// deleteSource deletes a source.
func (m *FairMix[P]) deleteSource(s *mixSource[P]) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.sources {
		if m.sources[i] == s {
			copy(m.sources[i:], m.sources[i+1:])
			m.sources[len(m.sources)-1] = nil
			m.sources = m.sources[:len(m.sources)-1]
			break
		}
	}
}

// runSource reads a single source in a loop.
func (m *FairMix[P]) runSource(closed chan struct{}, s *mixSource[P]) {
	defer m.wg.Done()
	defer close(s.next)
	for s.it.Next() {
		n := s.it.Node()
		select {
		case s.next <- n:
		case m.fromAny <- n:
		case <-closed:
			return
		}
	}
}
