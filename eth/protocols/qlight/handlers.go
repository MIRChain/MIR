package qlight

import (
	"fmt"

	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/eth/protocols/eth"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

func qlightClientHandleNewBlock[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg eth.Decoder, peer *Peer[T,P]) error {
	ann := new(eth.NewBlockPacket[P])
	if err := msg.Decode(ann); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	if hash := types.CalcUncleHash(ann.Block.Uncles()); hash != ann.Block.UncleHash() {
		log.Warn("Propagated block has invalid uncles", "have", hash, "exp", ann.Block.UncleHash())
		return nil // TODO(karalabe): return error eventually, but wait a few releases
	}
	if hash := types.DeriveSha(ann.Block.Transactions(), trie.NewStackTrie[P](nil)); hash != ann.Block.TxHash() {
		log.Warn("Propagated block has invalid body", "have", hash, "exp", ann.Block.TxHash())
		return nil // TODO(karalabe): return error eventually, but wait a few releases
	}
	if err := ann.Block.SanityCheck(); err != nil {
		return err
	}
	//TD at mainnet block #7753254 is 76 bits. If it becomes 100 million times
	// larger, it will still fit within 100 bits
	if tdlen := ann.TD.BitLen(); tdlen > 100 {
		return fmt.Errorf("too large block TD: bitlen %d", tdlen)
	}

	ann.Block.ReceivedAt = msg.Time()
	ann.Block.ReceivedFrom = peer

	// Mark the peer as owning the block
	peer.EthPeer.MarkBlock(ann.Block.Hash())

	return backend.QHandle(peer, ann)
}

func qlightClientHandleNewBlockPrivateData[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg eth.Decoder, peer *Peer[T,P]) error {
	res := new(BlockPrivateDataPacket)
	if err := msg.Decode(res); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	return backend.QHandle(peer, res)
}

func qlightClientHandleTransactions[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T,P], msg eth.Decoder, peer *Peer[T,P]) error {
	// Transactions arrived, make sure we have a valid and fresh chain to handle them
	if !backend.AcceptTxs() {
		return nil
	}
	// Transactions can be processed, parse all of them and deliver to the pool
	var txs eth.TransactionsPacket[P]
	if err := msg.Decode(&txs); err != nil {
		return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
	}
	for i, tx := range txs {
		// Validate and mark the remote transaction
		if tx == nil {
			return fmt.Errorf("%w: transaction %d is nil", errDecode, i)
		}
		peer.EthPeer.MarkTransaction(tx.Hash())
	}
	return backend.QHandle(peer, &txs)
}
