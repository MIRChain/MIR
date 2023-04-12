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
	"encoding/json"
	"fmt"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

// handleGetBlockHeaders handles Block header query, collect the requested headers and reply
func HandleGetBlockHeaders[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the complex header query
	var query GetBlockHeadersPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetBlockHeadersQuery(backend, &query, peer)
	return peer.SendBlockHeaders(response)
}

// handleGetBlockHeaders66 is the eth/66 version of handleGetBlockHeaders
func HandleGetBlockHeaders66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the complex header query
	var query GetBlockHeadersPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetBlockHeadersQuery(backend, query.GetBlockHeadersPacket, peer)
	return peer.ReplyBlockHeaders(query.RequestId, response)
}

func AnswerGetBlockHeadersQuery[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], query *GetBlockHeadersPacket, peer *Peer[T,P]) []*types.Header[P] {
	hashMode := query.Origin.Hash != (common.Hash{})
	first := true
	maxNonCanonical := uint64(100)

	// Gather headers until the fetch or network limits is reached
	var (
		bytes   common.StorageSize
		headers []*types.Header[P]
		unknown bool
		lookups int
	)
	for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit &&
		len(headers) < maxHeadersServe && lookups < 2*maxHeadersServe {
		lookups++
		// Retrieve the next header satisfying the query
		var origin *types.Header[P]
		if hashMode {
			if first {
				first = false
				origin = backend.Chain().GetHeaderByHash(query.Origin.Hash)
				if origin != nil {
					query.Origin.Number = origin.Number.Uint64()
				}
			} else {
				origin = backend.Chain().GetHeader(query.Origin.Hash, query.Origin.Number)
			}
		} else {
			origin = backend.Chain().GetHeaderByNumber(query.Origin.Number)
		}
		if origin == nil {
			break
		}
		headers = append(headers, origin)
		bytes += estHeaderSize

		// Advance to the next header of the query
		switch {
		case hashMode && query.Reverse:
			// Hash based traversal towards the genesis block
			ancestor := query.Skip + 1
			if ancestor == 0 {
				unknown = true
			} else {
				query.Origin.Hash, query.Origin.Number = backend.Chain().GetAncestor(query.Origin.Hash, query.Origin.Number, ancestor, &maxNonCanonical)
				unknown = (query.Origin.Hash == common.Hash{})
			}
		case hashMode && !query.Reverse:
			// Hash based traversal towards the leaf block
			var (
				current = origin.Number.Uint64()
				next    = current + query.Skip + 1
			)
			if next <= current {
				infos, _ := json.MarshalIndent(peer.Peer.Info(), "", "  ")
				peer.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
				unknown = true
			} else {
				if header := backend.Chain().GetHeaderByNumber(next); header != nil {
					nextHash := header.Hash()
					expOldHash, _ := backend.Chain().GetAncestor(nextHash, next, query.Skip+1, &maxNonCanonical)
					if expOldHash == query.Origin.Hash {
						query.Origin.Hash, query.Origin.Number = nextHash, next
					} else {
						unknown = true
					}
				} else {
					unknown = true
				}
			}
		case query.Reverse:
			// Number based traversal towards the genesis block
			if query.Origin.Number >= query.Skip+1 {
				query.Origin.Number -= query.Skip + 1
			} else {
				unknown = true
			}

		case !query.Reverse:
			// Number based traversal towards the leaf block
			query.Origin.Number += query.Skip + 1
		}
	}
	return headers
}

func HandleGetBlockBodies[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the block body retrieval message
	var query GetBlockBodiesPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetBlockBodiesQuery(backend, query, peer)
	return peer.SendBlockBodiesRLP(response)
}

func HandleGetBlockBodies66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the block body retrieval message
	var query GetBlockBodiesPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetBlockBodiesQuery(backend, query.GetBlockBodiesPacket, peer)
	return peer.ReplyBlockBodiesRLP(query.RequestId, response)
}

func AnswerGetBlockBodiesQuery[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], query GetBlockBodiesPacket, peer *Peer[T,P]) []rlp.RawValue {
	// Gather blocks until the fetch or network limits is reached
	var (
		bytes  int
		bodies []rlp.RawValue
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(bodies) >= maxBodiesServe ||
			lookups >= 2*maxBodiesServe {
			break
		}
		if data := backend.Chain().GetBodyRLP(hash); len(data) != 0 {
			bodies = append(bodies, data)
			bytes += len(data)
		}
	}
	return bodies
}

func HandleGetNodeData[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the trie node data retrieval message
	var query GetNodeDataPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetNodeDataQuery(backend, query, peer)
	return peer.SendNodeData(response)
}

func HandleGetNodeData66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the trie node data retrieval message
	var query GetNodeDataPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetNodeDataQuery(backend, query.GetNodeDataPacket, peer)
	return peer.ReplyNodeData(query.RequestId, response)
}

func AnswerGetNodeDataQuery[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], query GetNodeDataPacket, peer *Peer[T,P]) [][]byte {
	// Gather state data until the fetch or network limits is reached
	var (
		bytes int
		nodes [][]byte
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(nodes) >= maxNodeDataServe ||
			lookups >= 2*maxNodeDataServe {
			break
		}
		// Retrieve the requested state entry
		if bloom := backend.StateBloom(); bloom != nil && !bloom.Contains(hash[:]) {
			// Only lookup the trie node if there's chance that we actually have it
			continue
		}
		entry, err := backend.Chain().TrieNode(hash)
		if len(entry) == 0 || err != nil {
			// Read the contract code with prefix only to save unnecessary lookups.
			entry, err = backend.Chain().ContractCodeWithPrefix(hash)
		}
		if err == nil && len(entry) > 0 {
			nodes = append(nodes, entry)
			bytes += len(entry)
		}
	}
	return nodes
}

func HandleGetReceipts[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the block receipts retrieval message
	var query GetReceiptsPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetReceiptsQuery(backend, query, peer)
	return peer.SendReceiptsRLP(response)
}

func HandleGetReceipts66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the block receipts retrieval message
	var query GetReceiptsPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	response := AnswerGetReceiptsQuery(backend, query.GetReceiptsPacket, peer)
	return peer.ReplyReceiptsRLP(query.RequestId, response)
}

func AnswerGetReceiptsQuery[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], query GetReceiptsPacket, peer *Peer[T,P]) []rlp.RawValue {
	// Gather state data until the fetch or network limits is reached
	var (
		bytes    int
		receipts []rlp.RawValue
	)
	for lookups, hash := range query {
		if bytes >= softResponseLimit || len(receipts) >= maxReceiptsServe ||
			lookups >= 2*maxReceiptsServe {
			break
		}
		// Retrieve the requested block's receipts
		results := backend.Chain().GetReceiptsByHash(hash)
		if results == nil {
			if header := backend.Chain().GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
				continue
			}
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(results); err != nil {
			log.Error("Failed to encode receipt", "err", err)
		} else {
			receipts = append(receipts, encoded)
			bytes += len(encoded)
		}
	}
	return receipts
}

func HandleNewBlockhashes[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of new block announcements just arrived
	ann := new(NewBlockHashesPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Mark the hashes as present at the remote node
	for _, block := range *ann {
		peer.markBlock(block.Hash)
	}
	// Deliver them all to the backend for queuing
	return backend.Handle(peer, ann)
}

func HandleNewBlock[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Retrieve and decode the propagated block
	ann := new(NewBlockPacket[P])
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if err := ann.sanityCheck(); err != nil {
		return err
	}
	if hash := types.CalcUncleHash(ann.Block.Uncles()); hash != ann.Block.UncleHash() {
		log.Warn("Propagated block has invalid uncles", "have", hash, "exp", ann.Block.UncleHash())
		return nil // TODO(karalabe): return error eventually, but wait a few releases
	}
	if hash := types.DeriveSha(ann.Block.Transactions(), trie.NewStackTrie[P](nil)); hash != ann.Block.TxHash() {
		log.Warn("Propagated block has invalid body", "have", hash, "exp", ann.Block.TxHash())
		return nil // TODO(karalabe): return error eventually, but wait a few releases
	}
	ann.Block.ReceivedAt = msg.Time()
	ann.Block.ReceivedFrom = peer

	// Mark the peer as owning the block
	peer.markBlock(ann.Block.Hash())

	return backend.Handle(peer, ann)
}

func HandleBlockHeaders[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of headers arrived to one of our previous requests
	res := new(BlockHeadersPacket[P])
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	return backend.Handle(peer, res)
}

func HandleBlockHeaders66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of headers arrived to one of our previous requests
	res := new(BlockHeadersPacket66[P])
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	requestTracker.Fulfil(peer.id, peer.version, BlockHeadersMsg, res.RequestId)

	return backend.Handle(peer, &res.BlockHeadersPacket)
}

func HandleBlockBodies[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of block bodies arrived to one of our previous requests
	res := new(BlockBodiesPacket[P])
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	return backend.Handle(peer, res)
}

func HandleBlockBodies66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of block bodies arrived to one of our previous requests
	res := new(BlockBodiesPacket66[P])
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	requestTracker.Fulfil(peer.id, peer.version, BlockBodiesMsg, res.RequestId)

	return backend.Handle(peer, &res.BlockBodiesPacket)
}

func HandleNodeData[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of node state data arrived to one of our previous requests
	res := new(NodeDataPacket)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	return backend.Handle(peer, res)
}

func HandleNodeData66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of node state data arrived to one of our previous requests
	res := new(NodeDataPacket66)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	requestTracker.Fulfil(peer.id, peer.version, NodeDataMsg, res.RequestId)

	return backend.Handle(peer, &res.NodeDataPacket)
}

func HandleReceipts[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of receipts arrived to one of our previous requests
	res := new(ReceiptsPacket[P])
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	return backend.Handle(peer, res)
}

func HandleReceipts66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// A batch of receipts arrived to one of our previous requests
	res := new(ReceiptsPacket66[P])
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	requestTracker.Fulfil(peer.id, peer.version, ReceiptsMsg, res.RequestId)

	return backend.Handle(peer, &res.ReceiptsPacket)
}

func HandleNewPooledTransactionHashes[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// New transaction announcement arrived, make sure we have
	// a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}
	ann := new(NewPooledTransactionHashesPacket)
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	// Schedule all the unknown hashes for retrieval
	for _, hash := range *ann {
		peer.markTransaction(hash)
	}
	return backend.Handle(peer, ann)
}

func HandleGetPooledTransactions[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the pooled transactions retrieval message
	var query GetPooledTransactionsPacket
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	hashes, txs := AnswerGetPooledTransactions(backend, query, peer)
	return peer.SendPooledTransactionsRLP(hashes, txs)
}

func HandleGetPooledTransactions66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Decode the pooled transactions retrieval message
	var query GetPooledTransactionsPacket66
	if err := msg.Decode(&query); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	hashes, txs := AnswerGetPooledTransactions(backend, query.GetPooledTransactionsPacket, peer)
	return peer.ReplyPooledTransactionsRLP(query.RequestId, hashes, txs)
}

func AnswerGetPooledTransactions[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], query GetPooledTransactionsPacket, peer *Peer[T,P]) ([]common.Hash, []rlp.RawValue) {
	// Gather transactions until the fetch or network limits is reached
	var (
		bytes  int
		hashes []common.Hash
		txs    []rlp.RawValue
	)
	for _, hash := range query {
		if bytes >= softResponseLimit {
			break
		}
		// Retrieve the requested transaction, skipping if unknown to us
		tx := backend.TxPool().Get(hash)
		if tx == nil {
			continue
		}
		// If known, encode and queue for response packet
		if encoded, err := rlp.EncodeToBytes(tx); err != nil {
			log.Error("Failed to encode transaction", "err", err)
		} else {
			hashes = append(hashes, hash)
			txs = append(txs, encoded)
			bytes += len(encoded)
		}
	}
	return hashes, txs
}

func HandleTransactions[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs TransactionsPacket[P]
	if err := msg.Decode(&txs); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	for i, tx := range txs {
		// Validate and mark the remote transaction
		if tx == nil {
			return fmt.Errorf("%w: transaction %d is nil", errDecode, i)
		}
		peer.markTransaction(tx.Hash())
	}
	return backend.Handle(peer, &txs)
}

func HandlePooledTransactions[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs PooledTransactionsPacket[P]
	if err := msg.Decode(&txs); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	for i, tx := range txs {
		// Validate and mark the remote transaction
		if tx == nil {
			return fmt.Errorf("%w: transaction %d is nil", errDecode, i)
		}
		peer.markTransaction(tx.Hash())
	}
	return backend.Handle(peer, &txs)
}

func HandlePooledTransactions66[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg Decoder, peer *Peer[T,P]) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs PooledTransactionsPacket66[P]
	if err := msg.Decode(&txs); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	for i, tx := range txs.PooledTransactionsPacket {
		// Validate and mark the remote transaction
		if tx == nil {
			return fmt.Errorf("%w: transaction %d is nil", errDecode, i)
		}
		peer.markTransaction(tx.Hash())
	}
	requestTracker.Fulfil(peer.id, peer.version, PooledTransactionsMsg, txs.RequestId)

	return backend.Handle(peer, &txs.PooledTransactionsPacket)
}
