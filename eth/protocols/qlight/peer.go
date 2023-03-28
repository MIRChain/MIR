package qlight

import (
	"math/big"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/eth/protocols/eth"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/qlight"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

const (
	// maxKnownBlocks is the maximum block hashes to keep in the known list
	// before starting to randomly evict them.
	maxKnownBlocks = 1024

	// maxQueuedBlocks is the maximum number of block propagations to queue up before
	// dropping broadcasts. There's not much point in queueing stale blocks, so a few
	// that might cover uncles should be enough.
	maxQueuedBlocks = 4
)

// Peer is a collection of relevant information we have about a `snap` peer.
type Peer [T crypto.PrivateKey, P crypto.PublicKey] struct {
	id string // Unique ID for the peer, cached

	*p2p.Peer[T,P] // The embedded P2P package peer
	rw        p2p.MsgReadWriter
	version   uint // Protocol version negotiated

	logger log.Logger // Contextual logger with the peer id injected

	EthPeer *eth.Peer[T,P]

	knownBlocks     mapset.Set             // Set of block hashes known to be known by this peer
	queuedBlocks    chan *blockPropagation[P] // Queue of blocks to broadcast to the peer
	queuedBlockAnns chan *types.Block[P]      // Queue of blocks to announce to the peer

	term chan struct{} // Termination channel to stop the broadcasters
	lock sync.RWMutex  // Mutex protecting the internal fields

	qlightServer bool
	qlightPSI    string
	qlightToken  string

	QLightPeriodicAuthFunc func() error
}

// newPeer create a wrapper for a network connection and negotiated  protocol
// version.
func NewPeer [T crypto.PrivateKey, P crypto.PublicKey](version uint, p *p2p.Peer[T,P], rw p2p.MsgReadWriter, ethPeer *eth.Peer[T,P]) *Peer[T,P] {
	id := p.ID().String()
	return &Peer[T,P]{
		id:           id,
		Peer:         p,
		rw:           rw,
		version:      version,
		logger:       log.New("peer", id[:8]),
		EthPeer:      ethPeer,
		term:         make(chan struct{}),
		knownBlocks:  mapset.NewSet(),
		queuedBlocks: make(chan *blockPropagation[P], maxQueuedBlocks),
	}
}

func NewPeerWithBlockBroadcast[T crypto.PrivateKey, P crypto.PublicKey](version uint, p *p2p.Peer[T,P], rw p2p.MsgReadWriter, ethPeer *eth.Peer[T,P]) *Peer[T,P] {
	peer := NewPeer(version, p, rw, ethPeer)
	go peer.broadcastBlocksQLightServer()
	return peer
}

// ID retrieves the peer's unique identifier.
func (p *Peer[T,P]) ID() string {
	return p.id
}

// Version retrieves the peer's negoatiated `snap` protocol version.
func (p *Peer[T,P]) Version() uint {
	return p.version
}

// Log overrides the P2P logget with the higher level one containing only the id.
func (p *Peer[T,P]) Log() log.Logger {
	return p.logger
}

func (p *Peer[T,P]) QLightServer() bool {
	return p.qlightServer
}

func (p *Peer[T,P]) QLightPSI() string {
	return p.qlightPSI
}

func (p *Peer[T,P]) QLightToken() string {
	return p.qlightToken
}

func (p *Peer[T,P]) SendNewAuthToken(token string) error {
	return p2p.Send(p.rw, QLightTokenUpdateMsg, &qLightTokenUpdateData{
		Token: token,
	})
}

func (p *Peer[T,P]) SendNewBlock(block *types.Block[P], td *big.Int) error {
	// Mark all the block hash as known, but ensure we don't overflow our limits
	for p.knownBlocks.Cardinality() >= maxKnownBlocks {
		p.knownBlocks.Pop()
	}
	p.knownBlocks.Add(block.Hash())
	return p2p.Send(p.rw, eth.NewBlockMsg, &eth.NewBlockPacket[P]{
		Block: block,
		TD:    td,
	})
}

func (p *Peer[T,P]) SendBlockPrivateData(data []qlight.BlockPrivateData) error {
	// Mark all the block hash as known, but ensure we don't overflow our limits
	return p2p.Send(p.rw, QLightNewBlockPrivateDataMsg, data)
}

// AsyncSendNewBlock queues an entire block for propagation to a remote peer. If
// the peer's broadcast queue is full, the event is silently dropped.
func (p *Peer[T,P]) AsyncSendNewBlock(block *types.Block[P], td *big.Int, blockPrivateData *qlight.BlockPrivateData) {
	select {
	case p.queuedBlocks <- &blockPropagation[P]{block: block, td: td, blockPrivateData: blockPrivateData}:
		// Mark all the block hash as known, but ensure we don't overflow our limits
		for p.knownBlocks.Cardinality() >= maxKnownBlocks {
			p.knownBlocks.Pop()
		}
		p.knownBlocks.Add(block.Hash())
	default:
		p.Log().Debug("Dropping block propagation", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// KnownBlock returns whether peer is known to already have a block.
func (p *Peer[T,P]) KnownBlock(hash common.Hash) bool {
	return p.knownBlocks.Contains(hash)
}

func (p *Peer[T,P]) Close() {
	close(p.term)
}
