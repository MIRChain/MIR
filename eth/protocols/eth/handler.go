// Copyright 2020 The go-ethereum Authors
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
	"fmt"
	"math/big"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/metrics"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

const (
	// softResponseLimit is the target maximum size of replies to data retrievals.
	softResponseLimit = 2 * 1024 * 1024

	// estHeaderSize is the approximate size of an RLP encoded block header.
	estHeaderSize = 500

	// maxHeadersServe is the maximum number of block headers to serve. This number
	// is there to limit the number of disk lookups.
	maxHeadersServe = 1024

	// maxBodiesServe is the maximum number of block bodies to serve. This number
	// is mostly there to limit the number of disk lookups. With 24KB block sizes
	// nowadays, the practical limit will always be softResponseLimit.
	maxBodiesServe = 1024

	// maxNodeDataServe is the maximum number of state trie nodes to serve. This
	// number is there to limit the number of disk lookups.
	maxNodeDataServe = 1024

	// maxReceiptsServe is the maximum number of block receipts to serve. This
	// number is mostly there to limit the number of disk lookups. With block
	// containing 200+ transactions nowadays, the practical limit will always
	// be softResponseLimit.
	maxReceiptsServe = 1024
)

// Handler is a callback to invoke from an outside runner after the boilerplate
// exchanges have passed.
type Handler[T crypto.PrivateKey, P crypto.PublicKey] func(peer *Peer[T,P]) error

// Backend defines the data retrieval methods to serve remote requests and the
// callback methods to invoke on remote deliveries.
type Backend [T crypto.PrivateKey, P crypto.PublicKey] interface {
	// Chain retrieves the blockchain object to serve data.
	Chain() *core.BlockChain[P]

	// StateBloom retrieves the bloom filter - if any - for state trie nodes.
	StateBloom() *trie.SyncBloom

	// TxPool retrieves the transaction pool object to serve data.
	TxPool() TxPool[P]

	// AcceptTxs retrieves whether transaction processing is enabled on the node
	// or if inbound transactions should simply be dropped.
	AcceptTxs() bool

	// RunPeer is invoked when a peer joins on the `eth` protocol. The handler
	// should do any peer maintenance work, handshakes and validations. If all
	// is passed, control should be given back to the `handler` to process the
	// inbound messages going forward.
	RunPeer(peer *Peer[T,P], handler Handler[T,P]) error

	// PeerInfo retrieves all known `eth` information about a peer.
	PeerInfo(id enode.ID) interface{}

	// Handle is a callback to be invoked when a data packet is received from
	// the remote peer. Only packets not consumed by the protocol handler will
	// be forwarded to the backend.
	Handle(peer *Peer[T,P], packet Packet) error
}

// TxPool defines the methods needed by the protocol handler to serve transactions.
type TxPool [P crypto.PublicKey] interface {
	// Get retrieves the the transaction from the local txpool with the given hash.
	Get(hash common.Hash) *types.Transaction[P]
}

// MakeProtocols constructs the P2P protocol definitions for `eth`.
func MakeProtocols[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], network uint64, dnsdisc enode.Iterator[P]) []p2p.Protocol[T,P] {
	protocols := make([]p2p.Protocol[T,P], len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		version := version // Closure

		protocols[i] = p2p.Protocol[T,P]{
			Name:    ProtocolName,
			Version: version,
			Length:  protocolLengths[version],
			Run: func(p *p2p.Peer[T,P], rw p2p.MsgReadWriter) error {
				peer := NewPeer(version, p, rw, backend.TxPool())
				defer peer.Close()

				return backend.RunPeer(peer, func(peer *Peer[T,P]) error {
					return Handle(backend, peer)
				})
			},
			NodeInfo: func() interface{} {
				return nodeInfo[T,P](backend.Chain(), network)
			},
			PeerInfo: func(id enode.ID) interface{} {
				return backend.PeerInfo(id)
			},
			Attributes:     []enr.Entry{currentENREntry(backend.Chain())},
			DialCandidates: dnsdisc,
		}
	}
	return protocols
}

// NodeInfo represents a short summary of the `eth` sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64              `json:"network"`    // Ethereum network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Difficulty *big.Int            `json:"difficulty"` // Total difficulty of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // Hex hash of the host's best owned block
}

// nodeInfo retrieves some `eth` protocol metadata about the running host node.
func nodeInfo[T crypto.PrivateKey, P crypto.PublicKey](chain *core.BlockChain[P], network uint64) *NodeInfo {
	head := chain.CurrentBlock()
	return &NodeInfo{
		Network:    network,
		Difficulty: chain.GetTd(head.Hash(), head.NumberU64()),
		Genesis:    chain.Genesis().Hash(),
		Config:     chain.Config(),
		Head:       head.Hash(),
	}
}

// Handle is invoked whenever an `eth` connection is made that successfully passes
// the protocol handshake. This method will keep processing messages until the
// connection is torn down.
func Handle[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], peer *Peer[T,P]) error {
	for {
		if err := handleMessage(backend, peer); err != nil {
			peer.Log().Debug("Message handling failed in `eth`", "err", err)
			return err
		}
	}
}

type MsgHandler[T crypto.PrivateKey, P crypto.PublicKey] func(backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error
type Decoder interface {
	Decode(val interface{}) error
	Time() time.Time
}

// var eth65 = map[uint64]msgHandler{
// 	GetBlockHeadersMsg:            handleGetBlockHeaders,
// 	BlockHeadersMsg:               handleBlockHeaders,
// 	GetBlockBodiesMsg:             handleGetBlockBodies,
// 	BlockBodiesMsg:                handleBlockBodies,
// 	GetNodeDataMsg:                handleGetNodeData,
// 	NodeDataMsg:                   handleNodeData,
// 	GetReceiptsMsg:                handleGetReceipts,
// 	ReceiptsMsg:                   handleReceipts,
// 	NewBlockHashesMsg:             handleNewBlockhashes,
// 	NewBlockMsg:                   handleNewBlock,
// 	TransactionsMsg:               handleTransactions,
// 	NewPooledTransactionHashesMsg: handleNewPooledTransactionHashes,
// 	GetPooledTransactionsMsg:      handleGetPooledTransactions,
// 	PooledTransactionsMsg:         handlePooledTransactions,
// }

// var eth66 = map[uint64]msgHandler{
// 	NewBlockHashesMsg:             handleNewBlockhashes,
// 	NewBlockMsg:                   handleNewBlock,
// 	TransactionsMsg:               handleTransactions,
// 	NewPooledTransactionHashesMsg: handleNewPooledTransactionHashes,
// 	// eth66 messages with request-id
// 	GetBlockHeadersMsg:       handleGetBlockHeaders66,
// 	BlockHeadersMsg:          handleBlockHeaders66,
// 	GetBlockBodiesMsg:        handleGetBlockBodies66,
// 	BlockBodiesMsg:           handleBlockBodies66,
// 	GetNodeDataMsg:           handleGetNodeData66,
// 	NodeDataMsg:              handleNodeData66,
// 	GetReceiptsMsg:           handleGetReceipts66,
// 	ReceiptsMsg:              handleReceipts66,
// 	GetPooledTransactionsMsg: handleGetPooledTransactions66,
// 	PooledTransactionsMsg:    handlePooledTransactions66,
// }

// handleMessage is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func handleMessage[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], peer *Peer[T,P]) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := peer.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > maxMessageSize {
		return fmt.Errorf("%w: %v > %v", errMsgTooLarge, msg.Size, maxMessageSize)
	}
	defer msg.Discard()

	var handlers = map[uint64]MsgHandler[T,P]{
		GetBlockHeadersMsg:            HandleGetBlockHeaders[T,P],
		BlockHeadersMsg:               HandleBlockHeaders[T,P],
		GetBlockBodiesMsg:             HandleGetBlockBodies[T,P],
		BlockBodiesMsg:                HandleBlockBodies[T,P],
		GetNodeDataMsg:                HandleGetNodeData[T,P],
		NodeDataMsg:                   HandleNodeData[T,P],
		GetReceiptsMsg:                HandleGetReceipts[T,P],
		ReceiptsMsg:                   HandleReceipts[T,P],
		NewBlockHashesMsg:             HandleNewBlockhashes[T,P],
		NewBlockMsg:                   HandleNewBlock[T,P],
		TransactionsMsg:               HandleTransactions[T,P],
		NewPooledTransactionHashesMsg: HandleNewPooledTransactionHashes[T,P],
		GetPooledTransactionsMsg:      HandleGetPooledTransactions[T,P],
		PooledTransactionsMsg:         HandlePooledTransactions[T,P],
	}
	if peer.Version() >= ETH66 {
		handlers = map[uint64]MsgHandler[T,P]{
			NewBlockHashesMsg:             HandleNewBlockhashes[T,P],
			NewBlockMsg:                   HandleNewBlock[T,P],
			TransactionsMsg:               HandleTransactions[T,P],
			NewPooledTransactionHashesMsg: HandleNewPooledTransactionHashes[T,P],
			// eth66 messages with request-id
			GetBlockHeadersMsg:       HandleGetBlockHeaders66[T,P],
			BlockHeadersMsg:          HandleBlockHeaders66[T,P],
			GetBlockBodiesMsg:        HandleGetBlockBodies66[T,P],
			BlockBodiesMsg:           HandleBlockBodies66[T,P],
			GetNodeDataMsg:           HandleGetNodeData66[T,P],
			NodeDataMsg:              HandleNodeData66[T,P],
			GetReceiptsMsg:           HandleGetReceipts66[T,P],
			ReceiptsMsg:              HandleReceipts66[T,P],
			GetPooledTransactionsMsg: HandleGetPooledTransactions66[T,P],
			PooledTransactionsMsg:    HandlePooledTransactions66[T,P],
		}
	}
	// Track the amount of time it takes to serve the request and run the handler
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%d/%#02x", p2p.HandleHistName, ProtocolName, peer.Version(), msg.Code)
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.ResettingSample(
					metrics.NewExpDecaySample(1028, 0.015),
				)
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}
	if handler := handlers[msg.Code]; handler != nil {
		return handler(backend, msg, peer)
	}
	return fmt.Errorf("%w: %v", errInvalidMsgCode, msg.Code)
}
