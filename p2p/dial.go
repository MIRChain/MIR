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

package p2p

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	mrand "math/rand"
	"net"
	"sync"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common/mclock"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/netutil"
)

const (
	// This is the amount of time spent waiting in between redialing a certain node. The
	// limit is a bit higher than inboundThrottleTime to prevent failing dials in small
	// private networks.
	dialHistoryExpiration = inboundThrottleTime + 5*time.Second

	// Config for the "Looking for peers" message.
	dialStatsLogInterval = 10 * time.Second // printed at most this often
	dialStatsPeerLimit   = 3                // but not if more than this many dialed peers

	// Endpoint resolution is throttled with bounded backoff.
	initialResolveDelay = 60 * time.Second
	maxResolveDelay     = time.Hour
)

// NodeDialer is used to connect to nodes in the network, typically by using
// an underlying net.Dialer but also using net.Pipe in tests.
type NodeDialer [P crypto.PublicKey]interface {
	Dial(context.Context, *enode.Node[P]) (net.Conn, error)
}

type nodeResolver [P crypto.PublicKey]interface {
	Resolve(*enode.Node[P]) *enode.Node[P]
}

// tcpDialer implements NodeDialer using real TCP connections.
type tcpDialer [P crypto.PublicKey]struct {
	d *net.Dialer
}

func (t tcpDialer[P]) Dial(ctx context.Context, dest *enode.Node[P]) (net.Conn, error) {
	return t.d.DialContext(ctx, "tcp", nodeAddr(dest).String())
}

func nodeAddr[P crypto.PublicKey](n *enode.Node[P]) net.Addr {
	return &net.TCPAddr{IP: n.IP(), Port: n.TCP()}
}

// checkDial errors:
var (
	errSelf             = errors.New("is self")
	errAlreadyDialing   = errors.New("already dialing")
	errAlreadyConnected = errors.New("already connected")
	errRecentlyDialed   = errors.New("recently dialed")
	errNotWhitelisted   = errors.New("not contained in netrestrict whitelist")
	errNoPort           = errors.New("node does not provide TCP port")
)

// dialer creates outbound connections and submits them into Server.
// Two types of peer connections can be created:
//
//  - static dials are pre-configured connections. The dialer attempts
//    keep these nodes connected at all times.
//
//  - dynamic dials are created from node discovery results. The dialer
//    continuously reads candidate nodes from its input iterator and attempts
//    to create peer connections to nodes arriving through the iterator.
//
type dialScheduler [T crypto.PrivateKey, P crypto.PublicKey] struct {
	dialConfig[P]
	setupFunc   dialSetupFunc[P]
	wg          sync.WaitGroup
	cancel      context.CancelFunc
	ctx         context.Context
	nodesIn     chan *enode.Node[P]
	doneCh      chan *dialTask[T,P]
	addStaticCh chan *enode.Node[P]
	remStaticCh chan *enode.Node[P]
	addPeerCh   chan *conn[T,P]
	remPeerCh   chan *conn[T,P]

	// Everything below here belongs to loop and
	// should only be accessed by code on the loop goroutine.
	dialing   map[enode.ID]*dialTask[T,P] // active tasks
	peers     map[enode.ID]connFlag  // all connected peers
	dialPeers int                    // current number of dialed peers

	// The static map tracks all static dial tasks. The subset of usable static dial tasks
	// (i.e. those passing checkDial) is kept in staticPool. The scheduler prefers
	// launching random static tasks from the pool over launching dynamic dials from the
	// iterator.
	static     map[enode.ID]*dialTask[T,P]
	staticPool []*dialTask[T,P]

	// The dial history keeps recently dialed nodes. Members of history are not dialed.
	history          expHeap
	historyTimer     mclock.Timer
	historyTimerTime mclock.AbsTime

	// for logStats
	lastStatsLog     mclock.AbsTime
	doneSinceLastLog int
}

type dialSetupFunc [P crypto.PublicKey] func(net.Conn, connFlag, *enode.Node[P]) error

type dialConfig [P crypto.PublicKey] struct {
	self           enode.ID         // our own ID
	maxDialPeers   int              // maximum number of dialed peers
	maxActiveDials int              // maximum number of active dials
	netRestrict    *netutil.Netlist // IP whitelist, disabled if nil
	resolver       nodeResolver[P]
	dialer         NodeDialer[P]
	log            log.Logger
	clock          mclock.Clock
	rand           *mrand.Rand
}

func (cfg dialConfig[P]) withDefaults() dialConfig[P] {
	if cfg.maxActiveDials == 0 {
		cfg.maxActiveDials = defaultMaxPendingPeers
	}
	if cfg.log == nil {
		cfg.log = log.Root()
	}
	if cfg.clock == nil {
		cfg.clock = mclock.System{}
	}
	if cfg.rand == nil {
		seedb := make([]byte, 8)
		crand.Read(seedb)
		seed := int64(binary.BigEndian.Uint64(seedb))
		cfg.rand = mrand.New(mrand.NewSource(seed))
	}
	return cfg
}

func newDialScheduler[T crypto.PrivateKey, P crypto.PublicKey](config dialConfig[P], it enode.Iterator[P], setupFunc dialSetupFunc[P]) *dialScheduler[T,P] {
	d := &dialScheduler[T,P]{
		dialConfig:  config.withDefaults(),
		setupFunc:   setupFunc,
		dialing:     make(map[enode.ID]*dialTask[T,P]),
		static:      make(map[enode.ID]*dialTask[T,P]),
		peers:       make(map[enode.ID]connFlag),
		doneCh:      make(chan *dialTask[T,P]),
		nodesIn:     make(chan *enode.Node[P]),
		addStaticCh: make(chan *enode.Node[P]),
		remStaticCh: make(chan *enode.Node[P]),
		addPeerCh:   make(chan *conn[T,P]),
		remPeerCh:   make(chan *conn[T,P]),
	}
	d.lastStatsLog = d.clock.Now()
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.wg.Add(2)
	go d.readNodes(it)
	go d.loop(it)
	return d
}

// stop shuts down the dialer, canceling all current dial tasks.
func (d *dialScheduler[T,P]) stop() {
	d.cancel()
	d.wg.Wait()
}

// addStatic adds a static dial candidate.
func (d *dialScheduler[T,P]) addStatic(n *enode.Node[P]) {
	select {
	case d.addStaticCh <- n:
	case <-d.ctx.Done():
	}
}

// removeStatic removes a static dial candidate.
func (d *dialScheduler[T,P]) removeStatic(n *enode.Node[P]) {
	select {
	case d.remStaticCh <- n:
	case <-d.ctx.Done():
	}
}

// peerAdded updates the peer set.
func (d *dialScheduler[T,P]) peerAdded(c *conn[T,P]) {
	select {
	case d.addPeerCh <- c:
	case <-d.ctx.Done():
	}
}

// peerRemoved updates the peer set.
func (d *dialScheduler[T,P]) peerRemoved(c *conn[T,P]) {
	select {
	case d.remPeerCh <- c:
	case <-d.ctx.Done():
	}
}

// loop is the main loop of the dialer.
func (d *dialScheduler[T,P]) loop(it enode.Iterator[P]) {
	var (
		nodesCh    chan *enode.Node[P]
		historyExp = make(chan struct{}, 1)
	)

loop:
	for {
		// Launch new dials if slots are available.
		slots := d.freeDialSlots()
		slots -= d.startStaticDials(slots)
		if slots > 0 {
			nodesCh = d.nodesIn
		} else {
			nodesCh = nil
		}
		d.rearmHistoryTimer(historyExp)
		d.logStats()

		select {
		case node := <-nodesCh:
			if err := d.checkDial(node); err != nil {
				d.log.Trace("Discarding dial candidate", "id", node.ID(), "ip", node.IP(), "reason", err)
			} else {
				d.startDial(newDialTask[T,P](node, dynDialedConn))
			}

		case task := <-d.doneCh:
			id := task.dest.ID()
			delete(d.dialing, id)
			d.updateStaticPool(id)
			d.doneSinceLastLog++

		case c := <-d.addPeerCh:
			if c.is(dynDialedConn) || c.is(staticDialedConn) {
				d.dialPeers++
			}
			id := c.node.ID()
			d.peers[id] = c.flags
			// Remove from static pool because the node is now connected.
			task := d.static[id]
			if task != nil && task.staticPoolIndex >= 0 {
				d.removeFromStaticPool(task.staticPoolIndex)
			}
			// TODO: cancel dials to connected peers

		case c := <-d.remPeerCh:
			if c.is(dynDialedConn) || c.is(staticDialedConn) {
				d.dialPeers--
			}
			delete(d.peers, c.node.ID())
			d.updateStaticPool(c.node.ID())

		case node := <-d.addStaticCh:
			id := node.ID()
			_, exists := d.static[id]
			d.log.Trace("Adding static node", "id", id, "ip", node.IP(), "added", !exists)
			if exists {
				continue loop
			}
			task := newDialTask[T,P](node, staticDialedConn)
			d.static[id] = task
			if d.checkDial(node) == nil {
				d.addToStaticPool(task)
			}

		case node := <-d.remStaticCh:
			id := node.ID()
			task := d.static[id]
			d.log.Trace("Removing static node", "id", id, "ok", task != nil)
			if task != nil {
				delete(d.static, id)
				if task.staticPoolIndex >= 0 {
					d.removeFromStaticPool(task.staticPoolIndex)
				}
			}

		case <-historyExp:
			d.expireHistory()

		case <-d.ctx.Done():
			it.Close()
			break loop
		}
	}

	d.stopHistoryTimer(historyExp)
	for range d.dialing {
		<-d.doneCh
	}
	d.wg.Done()
}

// readNodes runs in its own goroutine and delivers nodes from
// the input iterator to the nodesIn channel.
func (d *dialScheduler[T,P]) readNodes(it enode.Iterator[P]) {
	defer d.wg.Done()

	for it.Next() {
		select {
		case d.nodesIn <- it.Node():
		case <-d.ctx.Done():
		}
	}
}

// logStats prints dialer statistics to the log. The message is suppressed when enough
// peers are connected because users should only see it while their client is starting up
// or comes back online.
func (d *dialScheduler[T,P]) logStats() {
	now := d.clock.Now()
	if d.lastStatsLog.Add(dialStatsLogInterval) > now {
		return
	}
	if d.dialPeers < dialStatsPeerLimit && d.dialPeers < d.maxDialPeers {
		d.log.Info("Looking for peers", "peercount", len(d.peers), "tried", d.doneSinceLastLog, "static", len(d.static))
	}
	d.doneSinceLastLog = 0
	d.lastStatsLog = now
}

// rearmHistoryTimer configures d.historyTimer to fire when the
// next item in d.history expires.
func (d *dialScheduler[T,P]) rearmHistoryTimer(ch chan struct{}) {
	if len(d.history) == 0 || d.historyTimerTime == d.history.nextExpiry() {
		return
	}
	d.stopHistoryTimer(ch)
	d.historyTimerTime = d.history.nextExpiry()
	timeout := time.Duration(d.historyTimerTime - d.clock.Now())
	d.historyTimer = d.clock.AfterFunc(timeout, func() { ch <- struct{}{} })
}

// stopHistoryTimer stops the timer and drains the channel it sends on.
func (d *dialScheduler[T,P]) stopHistoryTimer(ch chan struct{}) {
	if d.historyTimer != nil && !d.historyTimer.Stop() {
		<-ch
	}
}

// expireHistory removes expired items from d.history.
func (d *dialScheduler[T,P]) expireHistory() {
	d.historyTimer.Stop()
	d.historyTimer = nil
	d.historyTimerTime = 0
	d.history.expire(d.clock.Now(), func(hkey string) {
		var id enode.ID
		copy(id[:], hkey)
		d.updateStaticPool(id)
	})
}

// freeDialSlots returns the number of free dial slots. The result can be negative
// when peers are connected while their task is still running.
func (d *dialScheduler[T,P]) freeDialSlots() int {
	slots := (d.maxDialPeers - d.dialPeers) * 2
	if slots > d.maxActiveDials {
		slots = d.maxActiveDials
	}
	free := slots - len(d.dialing)
	return free
}

// checkDial returns an error if node n should not be dialed.
func (d *dialScheduler[T,P]) checkDial(n *enode.Node[P]) error {
	if n.ID() == d.self {
		return errSelf
	}
	if n.IP() != nil && n.TCP() == 0 {
		// This check can trigger if a non-TCP node is found
		// by discovery. If there is no IP, the node is a static
		// node and the actual endpoint will be resolved later in dialTask.
		return errNoPort
	}
	if _, ok := d.dialing[n.ID()]; ok {
		return errAlreadyDialing
	}
	if _, ok := d.peers[n.ID()]; ok {
		return errAlreadyConnected
	}
	if d.netRestrict != nil && !d.netRestrict.Contains(n.IP()) {
		return errNotWhitelisted
	}
	if d.history.contains(string(n.ID().Bytes())) {
		return errRecentlyDialed
	}
	return nil
}

// startStaticDials starts n static dial tasks.
func (d *dialScheduler[T,P]) startStaticDials(n int) (started int) {
	for started = 0; started < n && len(d.staticPool) > 0; started++ {
		idx := d.rand.Intn(len(d.staticPool))
		task := d.staticPool[idx]
		d.startDial(task)
		d.removeFromStaticPool(idx)
	}
	return started
}

// updateStaticPool attempts to move the given static dial back into staticPool.
func (d *dialScheduler[T,P]) updateStaticPool(id enode.ID) {
	task, ok := d.static[id]
	if ok && task.staticPoolIndex < 0 && d.checkDial(task.dest) == nil {
		d.addToStaticPool(task)
	}
}

func (d *dialScheduler[T,P]) addToStaticPool(task *dialTask[T,P]) {
	if task.staticPoolIndex >= 0 {
		panic("attempt to add task to staticPool twice")
	}
	d.staticPool = append(d.staticPool, task)
	task.staticPoolIndex = len(d.staticPool) - 1
}

// removeFromStaticPool removes the task at idx from staticPool. It does that by moving the
// current last element of the pool to idx and then shortening the pool by one.
func (d *dialScheduler[T,P]) removeFromStaticPool(idx int) {
	task := d.staticPool[idx]
	end := len(d.staticPool) - 1
	d.staticPool[idx] = d.staticPool[end]
	d.staticPool[idx].staticPoolIndex = idx
	d.staticPool[end] = nil
	d.staticPool = d.staticPool[:end]
	task.staticPoolIndex = -1
}

// startDial runs the given dial task in a separate goroutine.
func (d *dialScheduler[T,P]) startDial(task *dialTask[T,P]) {
	d.log.Trace("Starting p2p dial", "id", task.dest.ID(), "ip", task.dest.IP(), "flag", task.flags)
	hkey := string(task.dest.ID().Bytes())
	d.history.add(hkey, d.clock.Now().Add(dialHistoryExpiration))
	d.dialing[task.dest.ID()] = task
	go func() {
		task.run(d)
		d.doneCh <- task
	}()
}

// A dialTask generated for each node that is dialed.
type dialTask [T crypto.PrivateKey, P crypto.PublicKey] struct {
	staticPoolIndex int
	flags           connFlag
	// These fields are private to the task and should not be
	// accessed by dialScheduler while the task is running.
	dest         *enode.Node[P]
	lastResolved mclock.AbsTime
	resolveDelay time.Duration
}

func newDialTask[T crypto.PrivateKey, P crypto.PublicKey](dest *enode.Node[P], flags connFlag) *dialTask[T,P] {
	return &dialTask[T,P]{dest: dest, flags: flags, staticPoolIndex: -1}
}

type dialError struct {
	error
}

func (t *dialTask[T,P]) run(d *dialScheduler[T,P]) {
	if t.needResolve() && !t.resolve(d) {
		return
	}

	err := t.dial(d, t.dest)
	if err != nil {
		// For static nodes, resolve one more time if dialing fails.
		if _, ok := err.(*dialError); ok && t.flags&staticDialedConn != 0 {
			if t.resolve(d) {
				t.dial(d, t.dest)
			}
		}
	}
}

func (t *dialTask[T,P]) needResolve() bool {
	return t.flags&staticDialedConn != 0 && t.dest.IP() == nil
}

// resolve attempts to find the current endpoint for the destination
// using discovery.
//
// Resolve operations are throttled with backoff to avoid flooding the
// discovery network with useless queries for nodes that don't exist.
// The backoff delay resets when the node is found.
func (t *dialTask[T,P]) resolve(d *dialScheduler[T,P]) bool {
	if d.resolver == nil {
		return false
	}
	if t.resolveDelay == 0 {
		t.resolveDelay = initialResolveDelay
	}
	if t.lastResolved > 0 && time.Duration(d.clock.Now()-t.lastResolved) < t.resolveDelay {
		return false
	}
	resolved := d.resolver.Resolve(t.dest)
	t.lastResolved = d.clock.Now()
	if resolved == nil {
		t.resolveDelay *= 2
		if t.resolveDelay > maxResolveDelay {
			t.resolveDelay = maxResolveDelay
		}
		d.log.Debug("Resolving node failed", "id", t.dest.ID(), "newdelay", t.resolveDelay)
		return false
	}
	// The node was found.
	t.resolveDelay = initialResolveDelay
	t.dest = resolved
	d.log.Debug("Resolved node", "id", t.dest.ID(), "addr", &net.TCPAddr{IP: t.dest.IP(), Port: t.dest.TCP()})
	return true
}

// dial performs the actual connection attempt.
func (t *dialTask[T,P]) dial(d *dialScheduler[T,P], dest *enode.Node[P]) error {
	fd, err := d.dialer.Dial(d.ctx, t.dest)
	if err != nil {
		d.log.Trace("Dial error", "id", t.dest.ID(), "addr", nodeAddr(t.dest), "conn", t.flags, "err", cleanupDialErr(err))
		return &dialError{err}
	}
	mfd := newMeteredConn(fd, false, &net.TCPAddr{IP: dest.IP(), Port: dest.TCP()})
	return d.setupFunc(mfd, t.flags, dest)
}

func (t *dialTask[T,P]) String() string {
	id := t.dest.ID()
	return fmt.Sprintf("%v %x %v:%d", t.flags, id[:8], t.dest.IP(), t.dest.TCP())
}

func cleanupDialErr(err error) error {
	if netErr, ok := err.(*net.OpError); ok && netErr.Op == "dial" {
		return netErr.Err
	}
	return err
}
