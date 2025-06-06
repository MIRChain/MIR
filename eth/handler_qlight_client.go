package eth

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/forkid"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/eth/fetcher"
	"github.com/pavelkrolevets/MIR-pro/eth/protocols/eth"
	qlightproto "github.com/pavelkrolevets/MIR-pro/eth/protocols/qlight"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

type qlightClientHandler[T crypto.PrivateKey, P crypto.PublicKey] ethHandler[T,P]

func (h *qlightClientHandler[T,P]) Chain() *core.BlockChain[P]     { return h.chain }
func (h *qlightClientHandler[T,P]) StateBloom() *trie.SyncBloom { return h.stateBloom }
func (h *qlightClientHandler[T,P]) TxPool() eth.TxPool[P]          { return h.txpool }

func (h *qlightClientHandler[T,P]) RunPeer(peer *eth.Peer[T,P], handler eth.Handler[T,P]) error {
	return nil
}
func (h *qlightClientHandler[T,P]) Handle(peer *eth.Peer[T,P], packet eth.Packet) error {
	return (*ethHandler[T,P])(h).Handle(peer, packet)
}

func (h *qlightClientHandler[T,P]) RunQPeer(peer *qlightproto.Peer[T,P], hand qlightproto.Handler[T,P]) error {
	return (*handler[T,P])(h).runQLightClientPeer(peer, hand)
}

// PeerInfo retrieves all known `eth` information about a peer.
func (h *qlightClientHandler[T,P]) PeerInfo(id enode.ID) interface{} {
	if p := h.peers.peer(id.String()); p != nil {
		return p.info()
	}
	return nil
}

// AcceptTxs retrieves whether transaction processing is enabled on the node
// or if inbound transactions should simply be dropped.
func (h *qlightClientHandler[T,P]) AcceptTxs() bool {
	return atomic.LoadUint32(&h.acceptTxs) == 1
}

// newHandler returns a handler for all Ethereum chain management protocol.
func newQLightClientHandler[T crypto.PrivateKey, P crypto.PublicKey](config *handlerConfig[T,P]) (*handler[T,P], error) {
	// Create the protocol manager with the base fields
	if config.EventMux == nil {
		config.EventMux = new(event.TypeMux) // Nicety initialization for tests
	}
	h := &handler[T,P]{
		networkID:          config.Network,
		forkFilter:         forkid.NewFilter[P](config.Chain),
		eventMux:           config.EventMux,
		database:           config.Database,
		txpool:             config.TxPool,
		chain:              config.Chain,
		peers:              newPeerSet[T,P](),
		authorizationList:  config.AuthorizationList,
		txsyncCh:           make(chan *txsync[T,P]),
		quitSync:           make(chan struct{}),
		raftMode:           config.RaftMode,
		engine:             config.Engine,
		psi:                config.psi,
		privateClientCache: config.privateClientCache,
		tokenHolder:        config.tokenHolder,
	}

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
	if atomic.LoadUint32(&h.fastSync) == 1 {
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
	h.blockFetcher = fetcher.NewBlockFetcher(false, nil, h.chain.GetBlockByHash, validator, h.BroadcastBlockQLightClient, heighter, nil, inserter, h.removePeer)

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
func (h *handler[T,P]) runQLightClientPeer(peer *qlightproto.Peer[T,P], handler qlightproto.Handler[T,P]) error {
	// If the peer has a `snap` extension, wait for it to connect so we can have
	// a uniform initialization/teardown mechanism
	snap, err := h.peers.waitSnapExtension(peer.EthPeer)
	if err != nil {
		peer.Log().Error("Snapshot extension barrier failed", "err", err)
		return err
	}
	// TODO(karalabe): Not sure why this is needed
	if !h.chainSync.handlePeerEvent(peer.EthPeer) {
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
	if err := peer.EthPeer.Handshake(h.networkID, td, hash, genesis.Hash(), forkID, h.forkFilter); err != nil {
		peer.Log().Debug("Ethereum handshake failed", "err", err)

		// Quorum
		// When the Handshake() returns an error, the Run method corresponding to `eth` protocol returns with the error, causing the peer to drop, signal subprotocol as well to exit the `Run` method
		peer.EthPeerDisconnected <- struct{}{}
		// End Quorum
		return err
	}

	log.Info("QLight attempting handshake")
	if err := peer.QLightHandshake(false, h.psi, h.tokenHolder.CurrentToken()); err != nil {
		peer.Log().Debug("QLight handshake failed", "err", err)
		log.Info("QLight handshake failed", "err", err)

		// Quorum
		// When the Handshake() returns an error, the Run method corresponding to `eth` protocol returns with the error, causing the peer to drop, signal subprotocol as well to exit the `Run` method
		peer.EthPeerDisconnected <- struct{}{}
		// End Quorum
		return err
	}

	peer.Log().Debug("QLight handshake result for peer", "peer", peer.ID(), "server", peer.QLightServer(), "psi", peer.QLightPSI(), "token", peer.QLightToken())
	log.Info("QLight handshake result for peer", "peer", peer.ID(), "server", peer.QLightServer(), "psi", peer.QLightPSI(), "token", peer.QLightToken())
	// if we're not connected to a qlight server - disconnect the peer
	if !peer.QLightServer() {
		peer.Log().Debug("QLight connected to a non server peer. Disconnecting.")

		// Quorum
		// When the Handshake() returns an error, the Run method corresponding to `eth` protocol returns with the error, causing the peer to drop, signal subprotocol as well to exit the `Run` method
		peer.EthPeerDisconnected <- struct{}{}
		// End Quorum
		return fmt.Errorf("connected to a non server peer")
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
	if err := h.peers.registerQPeer(peer); err != nil {
		peer.Log().Error("Ethereum peer registration failed", "err", err)

		// Quorum
		// When the Register() returns an error, the Run method corresponding to `eth` protocol returns with the error, causing the peer to drop, signal subprotocol as well to exit the `Run` method
		peer.EthPeerDisconnected <- struct{}{}
		// End Quorum

		return err
	}
	defer h.removePeer(peer.ID())

	p := h.peers.peer(peer.ID())
	if p == nil {
		return errors.New("peer dropped during handling")
	}
	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := h.downloader.RegisterPeer(peer.ID(), peer.Version(), peer.EthPeer); err != nil {
		peer.Log().Error("Failed to register peer in eth syncer", "err", err)
		return err
	}
	if snap != nil {
		if err := h.downloader.SnapSyncer.Register(snap); err != nil {
			peer.Log().Error("Failed to register peer in snap syncer", "err", err)
			return err
		}
	}
	h.chainSync.handlePeerEvent(peer.EthPeer)

	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	h.syncTransactions(peer.EthPeer)

	// If we have a trusted CHT, reject all peers below that (avoid fast sync eclipse)
	if h.checkpointHash != (common.Hash{}) {
		// Request the peer's checkpoint header for chain height/weight validation
		if err := peer.EthPeer.RequestHeadersByNumber(h.checkpointNumber, 1, 0, false); err != nil {
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
		if err := peer.EthPeer.RequestHeadersByNumber(number, 1, 0, false); err != nil {
			return err
		}
	}

	// Quorum notify other subprotocols that the eth peer is ready, and has been added to the peerset.
	p.EthPeerRegistered <- struct{}{}
	// Quorum

	// Handle incoming messages until the connection is torn down
	return handler(peer)
}

func (h *handler[T,P]) StartQLightClient() {
	h.maxPeers = 1
	// Quorum
	if h.raftMode {
		// We set this immediately in raft mode to make sure the miner never drops
		// incoming txes. Raft mode doesn't use the fetcher or downloader, and so
		// this would never be set otherwise.
		atomic.StoreUint32(&h.acceptTxs, 1)
	}
	// End Quorum

	// start sync handlers
	h.wg.Add(1)
	go h.chainSync.loop()
}

func (h *handler[T,P]) StopQLightClient() {
	if h == nil {
		return
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
	log.Info("QLight client protocol stopped")
}

// BroadcastBlock will either propagate a block to a subset of its peers, or
// will only announce its availability (depending what's requested).
func (h *handler[T,P]) BroadcastBlockQLightClient(block *types.Block[P], propagate bool) {
}

// Handle is invoked from a peer's message handler when it receives a new remote
// message that the handler couldn't consume and serve itself.
func (h *qlightClientHandler[T,P]) QHandle(peer *qlightproto.Peer[T,P], packet eth.Packet) error {
	// Consume any broadcasts and announces, forwarding the rest to the downloader
	switch packet := packet.(type) {
	case *eth.BlockHeadersPacket[P]:
		return (*ethHandler[T,P])(h).Handle(peer.EthPeer, packet)

	case *eth.BlockBodiesPacket[P]:
		txset, uncleset := packet.Unpack()
		h.handleBodiesQLight(txset)
		return (*ethHandler[T,P])(h).handleBodies(peer.EthPeer, txset, uncleset)

	case *eth.NewBlockHashesPacket:
		return (*ethHandler[T,P])(h).Handle(peer.EthPeer, packet)

	case *eth.NewBlockPacket[P]:
		h.updateCacheWithNonPartyTxData(packet.Block.Transactions())
		return (*ethHandler[T,P])(h).handleBlockBroadcast(peer.EthPeer, packet.Block, packet.TD)
	case *qlightproto.BlockPrivateDataPacket:
		return h.handleBlockPrivateData(packet)
	case *eth.NewPooledTransactionHashesPacket:
		return (*ethHandler[T,P])(h).Handle(peer.EthPeer, packet)
	case *eth.TransactionsPacket[P]:
		h.updateCacheWithNonPartyTxData(*packet)
		return (*ethHandler[T,P])(h).Handle(peer.EthPeer, packet)
	case *eth.PooledTransactionsPacket[P]:
		return (*ethHandler[T,P])(h).Handle(peer.EthPeer, packet)

	default:
		return fmt.Errorf("unexpected eth packet type: %T", packet)
	}
}

// handleBodies is invoked from a peer's message handler when it transmits a batch
// of block bodies for the local node to process.
func (h *qlightClientHandler[T,P]) handleBodiesQLight(txs [][]*types.Transaction[P]) {
	for _, txArray := range txs {
		h.updateCacheWithNonPartyTxData(txArray)
	}
}

func (h *qlightClientHandler[T,P]) updateCacheWithNonPartyTxData(transactions []*types.Transaction[P]) {
	for _, tx := range transactions {
		if tx.IsPrivate() || tx.IsPrivacyMarker() {
			txHash := common.BytesToEncryptedPayloadHash(tx.Data())
			h.privateClientCache.CheckAndAddEmptyEntry(txHash)
		}
	}
}

func (h *qlightClientHandler[T,P]) handleBlockPrivateData(blockPrivateData *qlightproto.BlockPrivateDataPacket) error {
	for _, b := range *blockPrivateData {
		if err := h.privateClientCache.AddPrivateBlock(b); err != nil {
			return fmt.Errorf("Unable to handle private block data: %v", err)
		}
	}
	return nil
}
