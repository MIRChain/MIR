package eth

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/p2p"
)

func CurrentENREntry[T crypto.PrivateKey, P crypto.PublicKey](chain *core.BlockChain[P]) *enrEntry {
	return currentENREntry(chain)
}

func NodeInfoFunc[T crypto.PrivateKey, P crypto.PublicKey](chain *core.BlockChain[P], network uint64) *NodeInfo {
	return nodeInfo[T,P](chain, network)
}

// var ETH_65_FULL_SYNC = map[uint64]MsgHandler{
// 	// old 64 messages
// 	GetBlockHeadersMsg: handleGetBlockHeaders,
// 	BlockHeadersMsg:    handleBlockHeaders,
// 	GetBlockBodiesMsg:  handleGetBlockBodies,
// 	BlockBodiesMsg:     handleBlockBodies,
// 	NewBlockHashesMsg:  handleNewBlockhashes,
// 	NewBlockMsg:        handleNewBlock,
// 	TransactionsMsg:    handleTransactions,
// 	// New eth65 messages
// 	NewPooledTransactionHashesMsg: handleNewPooledTransactionHashes,
// 	GetPooledTransactionsMsg:      handleGetPooledTransactions,
// 	PooledTransactionsMsg:         handlePooledTransactions,
// }

func NewPeerWithTxBroadcast[T crypto.PrivateKey, P crypto.PublicKey](version uint, p *p2p.Peer[T,P], rw p2p.MsgReadWriter, txpool TxPool[P]) *Peer[T,P] {
	peer := NewPeerNoBroadcast(version, p, rw, txpool)
	// Start up all the broadcasters
	go peer.broadcastTransactions()
	if version >= ETH65 {
		go peer.announceTransactions()
	}
	return peer
}

func NewPeerNoBroadcast[T crypto.PrivateKey, P crypto.PublicKey](version uint, p *p2p.Peer[T,P], rw p2p.MsgReadWriter, txpool TxPool[P]) *Peer[T,P] {
	peer := &Peer[T,P]{
		id:              p.ID().String(),
		Peer:            p,
		rw:              rw,
		version:         version,
		knownTxs:        mapset.NewSet(),
		knownBlocks:     mapset.NewSet(),
		queuedBlocks:    make(chan *blockPropagation[P], maxQueuedBlocks),
		queuedBlockAnns: make(chan *types.Block[P], maxQueuedBlockAnns),
		txBroadcast:     make(chan []common.Hash),
		txAnnounce:      make(chan []common.Hash),
		txpool:          txpool,
		term:            make(chan struct{}),
	}
	return peer
}

func (p *Peer[T,P]) MarkBlock(hash common.Hash) {
	p.markBlock(hash)
}

func (p *Peer[T,P]) MarkTransaction(hash common.Hash) {
	p.markTransaction(hash)
}
