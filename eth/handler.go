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

package eth

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/clique"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/forkid"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/eth/fetcher"
	"github.com/pavelkrolevets/MIR-pro/eth/protocols/eth"
	"github.com/pavelkrolevets/MIR-pro/eth/protocols/snap"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/qlight"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

const (
	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096

	// Quorum
	protocolMaxMsgSize = 10 * 1024 * 1024 // Maximum cap on the size of a protocol message
)

var (
	syncChallengeTimeout = 15 * time.Second // Time allowance for a node to reply to the sync progress challenge

	// Quorum
	errMsgTooLarge = errors.New("message too long")
)

// txPool defines the methods needed from a transaction pool implementation to
// support all the operations needed by the Ethereum chain protocols.
type txPool [P crypto.PublicKey] interface {
	// Has returns an indicator whether txpool has a transaction
	// cached with the given hash.
	Has(hash common.Hash) bool

	// Get retrieves the transaction from local txpool with given
	// tx hash.
	Get(hash common.Hash) *types.Transaction[P]

	// AddRemotes should add the given transactions to the pool.
	AddRemotes([]*types.Transaction[P]) []error

	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending() (map[common.Address]types.Transactions[P], error)

	// SubscribeNewTxsEvent should return an event subscription of
	// NewTxsEvent and send events to the given channel.
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent[P]) event.Subscription
}

// handlerConfig is the collection of initialization parameters to create a full
// node network handler.
type handlerConfig [T crypto.PrivateKey, P crypto.PublicKey] struct {
	Database   ethdb.Database            // Database for direct sync insertions
	Chain      *core.BlockChain[P]          // Blockchain to serve data from
	TxPool     txPool[P]                    // Transaction pool to propagate from
	Network    uint64                    // Network identifier to adfvertise
	Sync       downloader.SyncMode       // Whether to fast or full sync
	BloomCache uint64                    // Megabytes to alloc for fast sync bloom
	EventMux   *event.TypeMux            // Legacy event mux, deprecate for `feed`
	Checkpoint *params.TrustedCheckpoint // Hard coded checkpoint for sync challenges

	// Quorum
	AuthorizationList map[uint64]common.Hash // Hard coded authorization list for sync challenged

	Engine   consensus.Engine[P]
	RaftMode bool

	// Quorum QLight
	// client
	psi                string
	privateClientCache qlight.PrivateClientCache
	tokenHolder        *qlight.TokenHolder[T,P]
	// server
	authProvider             qlight.AuthProvider
	privateBlockDataResolver qlight.PrivateBlockDataResolver[P]
}

type handler [T crypto.PrivateKey, P crypto.PublicKey] struct {
	networkID  uint64
	forkFilter forkid.Filter // Fork ID filter, constant across the lifetime of the node

	fastSync  uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	snapSync  uint32 // Flag whether fast sync should operate on top of the snap protocol
	acceptTxs uint32 // Flag whether we're considered synchronised (enables transaction processing)

	checkpointNumber uint64      // Block number for the sync progress validator to cross reference
	checkpointHash   common.Hash // Block hash for the sync progress validator to cross reference

	database ethdb.Database
	txpool   txPool[P]
	chain    *core.BlockChain[P]
	maxPeers int

	downloader   *downloader.Downloader[T,P]
	stateBloom   *trie.SyncBloom
	blockFetcher *fetcher.BlockFetcher[P]
	txFetcher    *fetcher.TxFetcher[P]
	peers        *peerSet[T,P]

	eventMux      *event.TypeMux
	txsCh         chan core.NewTxsEvent[P]
	txsSub        event.Subscription
	minedBlockSub *event.TypeMuxSubscription

	authorizationList map[uint64]common.Hash

	// channels for fetcher, syncer, txsyncLoop
	txsyncCh chan *txsync[T,P]
	quitSync chan struct{}

	chainSync *chainSyncer[T,P]
	wg        sync.WaitGroup
	peerWG    sync.WaitGroup

	// Quorum
	raftMode    bool
	engine      consensus.Engine[P]
	tokenHolder *qlight.TokenHolder[T,P]

	// Test fields or hooks
	broadcastTxAnnouncesOnly bool // Testing field, disable transaction propagation

	// Quorum QLight
	// client
	psi                string
	privateClientCache qlight.PrivateClientCache
	// server
	authProvider             qlight.AuthProvider
	privateBlockDataResolver qlight.PrivateBlockDataResolver[P]
}

// newHandler returns a handler for all Ethereum chain management protocol.
func newHandler[T crypto.PrivateKey, P crypto.PublicKey](config *handlerConfig[T,P]) (*handler[T,P], error) {
	// Create the protocol manager with the base fields
	if config.EventMux == nil {
		config.EventMux = new(event.TypeMux) // Nicety initialization for tests
	}
	h := &handler[T,P]{
		networkID:  config.Network,
		forkFilter: forkid.NewFilter[P](config.Chain),
		eventMux:   config.EventMux,
		database:   config.Database,
		txpool:     config.TxPool,
		chain:      config.Chain,
		peers:      newPeerSet[T,P](),
		txsyncCh:   make(chan *txsync[T,P]),
		quitSync:   make(chan struct{}),
		// Quorum
		authorizationList: config.AuthorizationList,
		raftMode:          config.RaftMode,
		engine:            config.Engine,
		tokenHolder:       config.tokenHolder,
	}

	// Quorum
	if handler, ok := h.engine.(consensus.Handler[P]); ok {
		handler.SetBroadcaster(h)
	}
	// /Quorum

	if config.Sync == downloader.FullSync {
		// The database seems empty as the current block is the genesis. Yet the fast
		// block is ahead, so fast sync was enabled for this node at a certain point.
		// The scenarios where this can happen is
		// * if the user manually (or via a bad block) rolled back a fast sync node
		//   below the sync point.
		// * the last fast sync is not finished while user specifies a full sync this
		//   time. But we don't have any recent state for full sync.
		// In these cases however it's safe to reenable fast sync.
		fullBlock, fastBlock := h.chain.CurrentBlock(), h.chain.CurrentFastBlock()
		if fullBlock.NumberU64() == 0 && fastBlock.NumberU64() > 0 {
			h.fastSync = uint32(1)
			log.Warn("Switch sync mode from full sync to fast sync")
		}
	} else {
		if h.chain.CurrentBlock().NumberU64() > 0 {
			// Print warning log if database is not empty to run fast sync.
			log.Warn("Switch sync mode from fast sync to full sync")
		} else {
			// If fast sync was requested and our database is empty, grant it
			h.fastSync = uint32(1)
			if config.Sync == downloader.SnapSync {
				h.snapSync = uint32(1)
			}
		}
	}
	// If we have trusted checkpoints, enforce them on the chain
	if config.Checkpoint != nil {
		h.checkpointNumber = (config.Checkpoint.SectionIndex+1)*params.CHTFrequency - 1
		h.checkpointHash = config.Checkpoint.SectionHead
	}
	// Construct the downloader (long sync) and its backing state bloom if fast
	// sync is requested. The downloader is responsible for deallocating the state
	// bloom when it's done.
	// Note: we don't enable it if snap-sync is performed, since it's very heavy
	// and the heal-portion of the snap sync is much lighter than fast. What we particularly
	// want to avoid, is a 90%-finished (but restarted) snap-sync to begin
	// indexing the entire trie
	if atomic.LoadUint32(&h.fastSync) == 1 && atomic.LoadUint32(&h.snapSync) == 0 {
		h.stateBloom = trie.NewSyncBloom(config.BloomCache, config.Database)
	}
	h.downloader = downloader.New[T,P](h.checkpointNumber, config.Database, h.stateBloom, h.eventMux, h.chain, nil, h.removePeer)

	// Construct the fetcher (short sync)
	validator := func(header *types.Header[P]) error {
		return h.chain.Engine().VerifyHeader(h.chain, header, true)
	}
	heighter := func() uint64 {
		return h.chain.CurrentBlock().NumberU64()
	}
	inserter := func(blocks types.Blocks[P]) (int, error) {
		// If sync hasn't reached the checkpoint yet, deny importing weird blocks.
		//
		// Ideally we would also compare the head block's timestamp and similarly reject
		// the propagated block if the head is too old. Unfortunately there is a corner
		// case when starting new networks, where the genesis might be ancient (0 unix)
		// which would prevent full nodes from accepting it.
		if h.chain.CurrentBlock().NumberU64() < h.checkpointNumber {
			log.Warn("Unsynced yet, discarded propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		// If fast sync is running, deny importing weird blocks. This is a problematic
		// clause when starting up a new network, because fast-syncing miners might not
		// accept each others' blocks until a restart. Unfortunately we haven't figured
		// out a way yet where nodes can decide unilaterally whether the network is new
		// or not. This should be fixed if we figure out a solution.
		if atomic.LoadUint32(&h.fastSync) == 1 {
			log.Warn("Fast syncing, discarded propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		n, err := h.chain.InsertChain(blocks)
		if err == nil {
			atomic.StoreUint32(&h.acceptTxs, 1) // Mark initial sync done on any fetcher import
		}
		return n, err
	}
	h.blockFetcher = fetcher.NewBlockFetcher[P](false, nil, h.chain.GetBlockByHash, validator, h.BroadcastBlock, heighter, nil, inserter, h.removePeer)

	fetchTx := func(peer string, hashes []common.Hash) error {
		p := h.peers.peer(peer)
		if p == nil {
			return errors.New("unknown peer")
		}
		return p.RequestTxs(hashes)
	}
	h.txFetcher = fetcher.NewTxFetcher(h.txpool.Has, h.txpool.AddRemotes, fetchTx)
	h.chainSync = newChainSyncer(h)
	return h, nil
}

// runEthPeer registers an eth peer into the joint eth/snap peerset, adds it to
// various subsistems and starts handling messages.
func (h *handler[T,P]) runEthPeer(peer *eth.Peer[T,P], handler eth.Handler[T,P]) error {
	// If the peer has a `snap` extension, wait for it to connect so we can have
	// a uniform initialization/teardown mechanism
	snap, err := h.peers.waitSnapExtension(peer)
	if err != nil {
		peer.Log().Error("Snapshot extension barrier failed", "err", err)
		return err
	}
	// TODO(karalabe): Not sure why this is needed
	if !h.chainSync.handlePeerEvent(peer) {
		return p2p.DiscQuitting
	}
	h.peerWG.Add(1)
	defer h.peerWG.Done()

	// Execute the Ethereum handshake
	var (
		genesis = h.chain.Genesis()
		head    = h.chain.CurrentHeader()
		hash    = head.Hash()
		number  = head.Number.Uint64()
		td      = h.chain.GetTd(hash, number)
	)
	forkID := forkid.NewID(h.chain.Config(), h.chain.Genesis().Hash(), h.chain.CurrentHeader().Number.Uint64())
	if err := peer.Handshake(h.networkID, td, hash, genesis.Hash(), forkID, h.forkFilter); err != nil {
		peer.Log().Debug("Ethereum handshake failed", "err", err)

		// Quorum
		// When the Handshake() returns an error, the Run method corresponding to `eth` protocol returns with the error, causing the peer to drop, signal subprotocol as well to exit the `Run` method
		peer.EthPeerDisconnected <- struct{}{}
		// End Quorum
		return err
	}
	reject := false // reserved peer slots
	if atomic.LoadUint32(&h.snapSync) == 1 {
		if snap == nil {
			// If we are running snap-sync, we want to reserve roughly half the peer
			// slots for peers supporting the snap protocol.
			// The logic here is; we only allow up to 5 more non-snap peers than snap-peers.
			if all, snp := h.peers.len(), h.peers.snapLen(); all-snp > snp+5 {
				reject = true
			}
		}
	}
	// Ignore maxPeers if this is a trusted peer
	if !peer.Peer.Info().Network.Trusted {
		if reject || h.peers.len() >= h.maxPeers {
			return p2p.DiscTooManyPeers
		}
	}
	peer.Log().Debug("Mir peer connected", "name", peer.Name())

	// Register the peer locally
	if err := h.peers.registerPeer(peer, snap); err != nil {
		peer.Log().Error("Mir peer registration failed", "err", err)

		// Quorum
		// When the Register() returns an error, the Run method corresponding to `eth` protocol returns with the error, causing the peer to drop, signal subprotocol as well to exit the `Run` method
		peer.EthPeerDisconnected <- struct{}{}
		// End Quorum
		return err
	}
	defer h.unregisterPeer(peer.ID()) // Quorum: changed by https://github.com/bnb-chain/bsc/pull/856

	p := h.peers.peer(peer.ID())
	if p == nil {
		return errors.New("peer dropped during handling")
	}
	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := h.downloader.RegisterPeer(peer.ID(), peer.Version(), peer); err != nil {
		peer.Log().Error("Failed to register peer in eth syncer", "err", err)
		return err
	}
	if snap != nil {
		if err := h.downloader.SnapSyncer.Register(snap); err != nil {
			peer.Log().Error("Failed to register peer in snap syncer", "err", err)
			return err
		}
	}
	h.chainSync.handlePeerEvent(peer)

	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	h.syncTransactions(peer)

	// If we have a trusted CHT, reject all peers below that (avoid fast sync eclipse)
	if h.checkpointHash != (common.Hash{}) {
		// Request the peer's checkpoint header for chain height/weight validation
		if err := peer.RequestHeadersByNumber(h.checkpointNumber, 1, 0, false); err != nil {
			return err
		}
		// Start a timer to disconnect if the peer doesn't reply in time
		p.syncDrop = time.AfterFunc(syncChallengeTimeout, func() {
			peer.Log().Warn("Checkpoint challenge timed out, dropping", "addr", peer.RemoteAddr(), "type", peer.Name())
			h.removePeer(peer.ID())
		})
		// Make sure it's cleaned up if the peer dies off
		defer func() {
			if p.syncDrop != nil {
				p.syncDrop.Stop()
				p.syncDrop = nil
			}
		}()
	}
	// If we have any explicit authorized block hashes, request them
	for number := range h.authorizationList {
		if err := peer.RequestHeadersByNumber(number, 1, 0, false); err != nil {
			return err
		}
	}

	// Quorum notify other subprotocols that the eth peer is ready, and has been added to the peerset.
	p.EthPeerRegistered <- struct{}{}
	// Quorum

	// Handle incoming messages until the connection is torn down
	return handler(peer)
}

// runSnapExtension registers a `snap` peer into the joint eth/snap peerset and
// starts handling inbound messages. As `snap` is only a satellite protocol to
// `eth`, all subsystem registrations and lifecycle management will be done by
// the main `eth` handler to prevent strange races.
func (h *handler[T,P]) runSnapExtension(peer *snap.Peer[T,P], handler snap.Handler[T,P]) error {
	h.peerWG.Add(1)
	defer h.peerWG.Done()

	if err := h.peers.registerSnapExtension(peer); err != nil {
		peer.Log().Error("Snapshot extension registration failed", "err", err)
		return err
	}
	return handler(peer)
}

// removePeer requests disconnection of a peer.
// Quorum: added by https://github.com/bnb-chain/bsc/pull/856
func (h *handler[T,P]) removePeer(id string) {
	peer := h.peers.peer(id)
	if peer != nil {
		// Hard disconnect at the networking layer. Handler will get an EOF and terminate the peer. defer unregisterPeer will do the cleanup task after then.
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

// unregisterPeer removes a peer from the downloader, fetchers and main peer set.
// Quorum: changed by https://github.com/bnb-chain/bsc/pull/856
func (h *handler[T,P]) unregisterPeer(id string) {
	// Create a custom logger to avoid printing the entire id
	var logger log.Logger
	if len(id) < 16 {
		// Tests use short IDs, don't choke on them
		logger = log.New("peer", id)
	} else {
		logger = log.New("peer", id[:8])
	}
	// Abort if the peer does not exist
	peer := h.peers.peer(id)
	if peer == nil {
		logger.Error("Ethereum peer removal failed", "err", errPeerNotRegistered)
		return
	}
	// Remove the `eth` peer if it exists
	logger.Debug("Removing Ethereum peer", "snap", peer.snapExt != nil)

	// Remove the `snap` extension if it exists
	if peer.snapExt != nil {
		h.downloader.SnapSyncer.Unregister(id)
	}
	h.downloader.UnregisterPeer(id)
	h.txFetcher.Drop(id)

	if err := h.peers.unregisterPeer(id); err != nil {
		logger.Error("Ethereum peer removal failed", "err", err)
	}
	// Hard disconnect at the networking layer
	// peer.Peer.Disconnect(p2p.DiscUselessPeer) // Quorum: removed by https://github.com/bnb-chain/bsc/pull/856
}

func (h *handler[T,P]) Start(maxPeers int) {
	h.maxPeers = maxPeers

	// broadcast transactions
	h.wg.Add(1)
	h.txsCh = make(chan core.NewTxsEvent[P], txChanSize)
	h.txsSub = h.txpool.SubscribeNewTxsEvent(h.txsCh)
	go h.txBroadcastLoop()

	// Quorum
	if !h.raftMode {
		// broadcast mined blocks
		h.wg.Add(1)
		h.minedBlockSub = h.eventMux.Subscribe(core.NewMinedBlockEvent[P]{})
		go h.minedBroadcastLoop()
	} else {
		// We set this immediately in raft mode to make sure the miner never drops
		// incoming txes. Raft mode doesn't use the fetcher or downloader, and so
		// this would never be set otherwise.
		atomic.StoreUint32(&h.acceptTxs, 1)
	}
	// End Quorum

	// start sync handlers
	h.wg.Add(2)
	go h.chainSync.loop()
	go h.txsyncLoop64() // TODO(karalabe): Legacy initial tx echange, drop with eth/64.
}

func (h *handler[T,P]) Stop() {
	h.txsSub.Unsubscribe() // quits txBroadcastLoop
	// quorum - ensure raft stops cleanly
	if h.minedBlockSub != nil {
		h.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop
	}

	// Quit chainSync and txsync64.
	// After this is done, no new peers will be accepted.
	close(h.quitSync)
	h.wg.Wait()

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to h.peers yet
	// will exit when they try to register.
	h.peers.close()
	h.peerWG.Wait()

	log.Info("Ethereum protocol stopped")
}

// Quorum
func (h *handler[T,P]) Enqueue(id string, block *types.Block[P]) {
	h.blockFetcher.Enqueue(id, block)
}

// BroadcastBlock will either propagate a block to a subset of its peers, or
// will only announce its availability (depending what's requested).
func (h *handler[T,P]) BroadcastBlock(block *types.Block[P], propagate bool) {
	hash := block.Hash()
	peers := h.peers.peersWithoutBlock(hash)

	// If propagation is requested, send to a subset of the peer
	if propagate {
		// Calculate the TD of the block (it's not imported yet, so block.Td is not valid)
		var td *big.Int
		if parent := h.chain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
			td = new(big.Int).Add(block.Difficulty(), h.chain.GetTd(block.ParentHash(), block.NumberU64()-1))
		} else {
			log.Error("Propagating dangling block", "number", block.Number(), "hash", hash)
			return
		}
		// Send the block to a subset of our peers
		transfer := peers[:int(math.Sqrt(float64(len(peers))))]
		for _, peer := range transfer {
			peer.AsyncSendNewBlock(block, td)
		}
		log.Trace("Propagated block", "hash", hash, "recipients", len(transfer), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
		return
	}
	// Otherwise if the block is indeed in out own chain, announce it
	if h.chain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.AsyncSendNewBlockHash(block)
		}
		log.Trace("Announced block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTransactions will propagate a batch of transactions
// - To a square root of all peers
// - And, separately, as announcements to all peers which are not known to
// already have the given transaction.
func (h *handler[T,P]) BroadcastTransactions(txs types.Transactions[P]) {
	var (
		annoCount   int // Count of announcements made
		annoPeers   int
		directCount int // Count of the txs sent directly to peers
		directPeers int // Count of the peers that were sent transactions directly

		txset = make(map[*ethPeer[T,P]][]common.Hash) // Set peer->hash to transfer directly
		annos = make(map[*ethPeer[T,P]][]common.Hash) // Set peer->hash to announce

	)
	// Broadcast transactions to a batch of peers not knowing about it
	// NOTE: Raft-based consensus currently assumes that geth broadcasts
	// transactions to all peers in the network. A previous comment here
	// indicated that this logic might change in the future to only send to a
	// subset of peers. If this change occurs upstream, a merge conflict should
	// arise here, and we should add logic to send to *all* peers in raft mode.

	for _, tx := range txs {
		peers := h.peers.peersWithoutTransaction(tx.Hash())
		// Send the tx unconditionally to a subset of our peers
		// Quorum changes for broadcasting to all peers not only Sqrt
		//numDirect := int(math.Sqrt(float64(len(peers))))
		for _, peer := range peers {
			txset[peer] = append(txset[peer], tx.Hash())
		}
		// For the remaining peers, send announcement only
		//for _, peer := range peers[numDirect:] {
		//	annos[peer] = append(annos[peer], tx.Hash())
		//}
		log.Trace("Broadcast transaction", "hash", tx.Hash(), "recipients", len(peers))
	}
	for peer, hashes := range txset {
		directPeers++
		directCount += len(hashes)
		peer.AsyncSendTransactions(hashes)
	}
	for peer, hashes := range txset {
		directPeers++
		directCount += len(hashes)
		peer.AsyncSendTransactions(hashes)
	}
	for peer, hashes := range annos {
		annoPeers++
		annoCount += len(hashes)
		peer.AsyncSendPooledTransactionHashes(hashes)
	}
	log.Debug("Transaction broadcast", "txs", len(txs),
		"announce packs", annoPeers, "announced hashes", annoCount,
		"tx packs", directPeers, "broadcast txs", directCount)
}

// minedBroadcastLoop sends mined blocks to connected peers.
func (h *handler[T,P]) minedBroadcastLoop() {
	defer h.wg.Done()

	for obj := range h.minedBlockSub.Chan() {
		if ev, ok := obj.Data.(core.NewMinedBlockEvent[P]); ok {
			h.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			h.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}

// txBroadcastLoop announces new transactions to connected peers.
func (h *handler[T,P]) txBroadcastLoop() {
	defer h.wg.Done()
	for {
		select {
		case event := <-h.txsCh:
			h.BroadcastTransactions(event.Txs)
		case <-h.txsSub.Err():
			return
		}
	}
}

// NodeInfo represents a short summary of the Ethereum sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64              `json:"network"`    // Ethereum network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Difficulty *big.Int            `json:"difficulty"` // Total difficulty of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // SHA3 hash of the host's best owned block
	Consensus  string              `json:"consensus"`  // Consensus mechanism in use
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (h *handler[T,P]) NodeInfo() *NodeInfo {
	currentBlock := h.chain.CurrentBlock()
	// //Quorum
	//
	// changes done to fetch maxCodeSize dynamically based on the
	// maxCodeSizeConfig changes
	// /Quorum
	chainConfig := h.chain.Config()
	chainConfig.MaxCodeSize = uint64(chainConfig.GetMaxCodeSize(h.chain.CurrentBlock().Number()) / 1024)

	return &NodeInfo{
		Network:    h.networkID,
		Difficulty: h.chain.GetTd(currentBlock.Hash(), currentBlock.NumberU64()),
		Genesis:    h.chain.Genesis().Hash(),
		Config:     chainConfig,
		Head:       currentBlock.Hash(),
		Consensus:  h.getConsensusAlgorithm(),
	}
}

// Quorum
func (h *handler[T,P]) getConsensusAlgorithm() string {
	var consensusAlgo string
	if h.raftMode { // raft does not use consensus interface
		consensusAlgo = "raft"
	} else {
		switch h.engine.(type) {
		case consensus.Istanbul[P]:
			consensusAlgo = "istanbul"
		case *clique.Clique[P]:
			consensusAlgo = "clique"
		case *ethash.Ethash[P]:
			consensusAlgo = "ethash"
		default:
			consensusAlgo = "unknown"
		}
	}
	return consensusAlgo
}

func (h *handler[T,P]) FindPeers(targets map[common.Address]bool) map[common.Address]consensus.Peer {
	m := make(map[common.Address]consensus.Peer)
	h.peers.lock.RLock()
	defer h.peers.lock.RUnlock()
	for _, p := range h.peers.peers {
		pubKey := p.Node().Pubkey()
		addr := crypto.PubkeyToAddress(pubKey)
		if targets[addr] {
			m[addr] = p
		}
	}
	return m
}

// makeQuorumConsensusProtocol is similar to eth/handler.go -> makeProtocol. Called from eth/handler.go -> Protocols.
// returns the supported subprotocol to the p2p server.
// The Run method starts the protocol and is called by the p2p server. The quorum consensus subprotocol,
// leverages the peer created and managed by the "eth" subprotocol.
// The quorum consensus protocol requires that the "eth" protocol is running as well.
func (h *handler[T,P]) makeQuorumConsensusProtocol(protoName string, version uint, length uint64, backend eth.Backend[T,P], network uint64, dnsdisc enode.Iterator[P]) p2p.Protocol[T,P] {
	log.Debug("registering qouorum protocol ", "protoName", protoName, "version", version)

	return p2p.Protocol[T,P]{
		Name:    protoName,
		Version: version,
		Length:  length,
		// no new peer created, uses the "eth" peer, so no peer management needed.
		Run: func(p *p2p.Peer[T,P], rw p2p.MsgReadWriter) error {
			/*
			* 1. wait for the eth protocol to create and register an eth peer.
			* 2. get the associate eth peer that was registered by he "eth" protocol.
			* 3. add the rw protocol for the quorum subprotocol to the eth peer.
			* 4. start listening for incoming messages.
			* 5. the incoming message will be sent on the quorum specific subprotocol, e.g. "istanbul/100".
			* 6. send messages to the consensus engine handler.
			* 7. messages to other to other peers listening to the subprotocol can be sent using the
			*    (eth)peer.ConsensusSend() which will write to the protoRW.
			 */
			// wait for the "eth" protocol to create and register the peer (added to peerset)
			select {
			case <-p.EthPeerRegistered:
				// the ethpeer should be registered, try to retrieve it and start the consensus handler.
				p2pPeerId := fmt.Sprintf("%x", p.ID().Bytes()[:8])
				ethPeer := h.peers.peer(p2pPeerId)
				if ethPeer == nil {
					p2pPeerId = fmt.Sprintf("%x", p.ID().Bytes()) //TODO:BBO
					ethPeer = h.peers.peer(p2pPeerId)
					log.Warn("full p2p peer", "id", p2pPeerId, "ethPeer", ethPeer)
				}
				if ethPeer != nil {
					p.Log().Debug("consensus subprotocol retrieved eth peer from peerset", "ethPeer.id", p2pPeerId, "ProtoName", protoName)
					// add the rw protocol for the quorum subprotocol to the eth peer.
					ethPeer.AddConsensusProtoRW(rw)
					peer := eth.NewPeer[T,P](version, p, rw, h.txpool)
					return h.handleConsensusLoop(peer, rw, nil)
				}
				p.Log().Error("consensus subprotocol retrieved nil eth peer from peerset", "ethPeer.id", p2pPeerId)
				return errEthPeerNil
			case <-p.EthPeerDisconnected:
				return errEthPeerNotRegistered
			}
		},
		NodeInfo: func() interface{} {
			return eth.NodeInfoFunc[T,P](backend.Chain(), network)
		},
		PeerInfo: func(id enode.ID) interface{} {
			return backend.PeerInfo(id)
		},
		Attributes:     []enr.Entry{eth.CurrentENREntry[T,P](backend.Chain())},
		DialCandidates: dnsdisc,
	}
}

func (h *handler[T,P]) handleConsensusLoop(p *eth.Peer[T,P], protoRW p2p.MsgReadWriter, fallThroughBackend eth.Backend[T,P]) error {
	// Handle incoming messages until the connection is torn down
	for {
		if err := h.handleConsensus(p, protoRW, fallThroughBackend); err != nil {
			// allow the P2P connection to remain active during sync (when the engine is stopped)
			if errors.Is(err, istanbul.ErrStoppedEngine) && h.downloader.Synchronising() {
				// should this be warn or debug
				p.Log().Debug("Ignoring `stopped engine` consensus error due to active sync.")
				continue
			}
			p.Log().Debug("Ethereum quorum message handling failed", "err", err)
			return err
		}
	}
}

// This is a no-op because the eth handleMsg main loop handle ibf message as well.
func (h *handler[T,P]) handleConsensus(p *eth.Peer[T,P], protoRW p2p.MsgReadWriter, fallThroughBackend eth.Backend[T,P]) error {
	// Read the next message from the remote peer (in protoRW), and ensure it's fully consumed
	msg, err := protoRW.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > protocolMaxMsgSize {
		return fmt.Errorf("%w: %v > %v", errMsgTooLarge, msg.Size, protocolMaxMsgSize)
	}
	defer msg.Discard()

	// See if the consensus engine protocol can handle this message, e.g. istanbul will check for message is
	// istanbulMsg = 0x11, and NewBlockMsg = 0x07.
	handled, err := h.handleConsensusMsg(p, msg)
	if handled {
		p.Log().Debug("consensus message was handled by consensus engine", "msg", msg.Code,
			"quorumConsensusProtocolName", quorumConsensusProtocolName, "err", err)
		return err
	}

	if fallThroughBackend != nil {
		var handlers = map[uint64]eth.MsgHandler[T,P]{
			// old 64 messages
			eth.GetBlockHeadersMsg: eth.HandleGetBlockHeaders[T,P],
			eth.BlockHeadersMsg:    eth.HandleBlockHeaders[T,P],
			eth.GetBlockBodiesMsg:  eth.HandleGetBlockBodies[T,P],
			eth.BlockBodiesMsg:     eth.HandleBlockBodies[T,P],
			eth.NewBlockHashesMsg:  eth.HandleNewBlockhashes[T,P],
			eth.NewBlockMsg:        eth.HandleNewBlock[T,P],
			eth.TransactionsMsg:    eth.HandleTransactions[T,P],
			// New eth65 messages
			eth.NewPooledTransactionHashesMsg: eth.HandleNewPooledTransactionHashes[T,P],
			eth.GetPooledTransactionsMsg:      eth.HandleGetPooledTransactions[T,P],
			eth.PooledTransactionsMsg:         eth.HandlePooledTransactions[T,P],
		}

		p.Log().Trace("Message not handled by legacy sub-protocol", "msg", msg.Code)

		if handler := handlers[msg.Code]; handler != nil {
			p.Log().Debug("Found eth handler for msg", "msg", msg.Code)
			return handler(fallThroughBackend, msg, p)
		}
	}

	return nil
}

func (h *handler[T,P]) handleConsensusMsg(p *eth.Peer[T,P], msg p2p.Msg) (bool, error) {
	if handler, ok := h.engine.(consensus.Handler[P]); ok {
		pubKey := p.Node().Pubkey()
		addr := crypto.PubkeyToAddress(pubKey)
		handled, err := handler.HandleMsg(addr, msg)
		return handled, err
	}
	return false, nil
}

// makeLegacyProtocol is basically a copy of the eth makeProtocol, but for legacy subprotocols, e.g. "istanbul/99" "istabnul/64"
// If support legacy subprotocols is removed, remove this and associated code as well.
// If quorum is using a legacy protocol then the "eth" subprotocol should not be available.
func (h *handler[T,P]) makeLegacyProtocol(protoName string, version uint, length uint64, backend eth.Backend[T,P], network uint64, dnsdisc enode.Iterator[P]) p2p.Protocol[T,P] {
	log.Debug("registering a legacy protocol ", "protoName", protoName, "version", version)
	return p2p.Protocol[T,P]{
		Name:    protoName,
		Version: version,
		Length:  length,
		Run: func(p *p2p.Peer[T,P], rw p2p.MsgReadWriter) error {
			peer := eth.NewPeer[T,P](version, p, rw, h.txpool)
			return h.runEthPeer(peer, func(peer *eth.Peer[T,P]) error {
				// We pass through the backend so that we can 'handle' messages that we can't handle
				return h.handleConsensusLoop(peer, rw, backend)
			})
		},
		NodeInfo: func() interface{} {
			return eth.NodeInfoFunc[T,P](backend.Chain(), network)
		},
		PeerInfo: func(id enode.ID) interface{} {
			return backend.PeerInfo(id)
		},
		Attributes:     []enr.Entry{eth.CurrentENREntry[T,P](backend.Chain())},
		DialCandidates: dnsdisc,
	}
}

// End Quorum
