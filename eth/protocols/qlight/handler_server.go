package qlight

import (
	"fmt"

	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/eth/protocols/eth"
	"github.com/MIRChain/MIR/p2p"
	"github.com/MIRChain/MIR/p2p/enode"
	"github.com/MIRChain/MIR/p2p/enr"
)

// MakeProtocolsServer constructs the P2P protocol definitions for `qlight` server.
func MakeProtocolsServer[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T, P], network uint64, dnsdisc enode.Iterator[P]) []p2p.Protocol[T, P] {
	protocols := make([]p2p.Protocol[T, P], 1)
	version := uint(QLIGHT65)
	protocols[0] = p2p.Protocol[T, P]{
		Name:    ProtocolName,
		Version: QLIGHT65,
		Length:  QLightProtocolLength,
		Run: func(p *p2p.Peer[T, P], rw p2p.MsgReadWriter) error {
			ethPeer := eth.NewPeerWithTxBroadcast(version, p, rw, backend.TxPool())
			peer := NewPeerWithBlockBroadcast(version, p, rw, ethPeer)
			defer ethPeer.Close()
			defer peer.Close()

			return backend.RunQPeer(peer, func(peer *Peer[T, P]) error {
				return HandleServer(backend, peer)
			})
		},
		NodeInfo: func() interface{} {
			return eth.NodeInfoFunc[T, P](backend.Chain(), network)
		},
		PeerInfo: func(id enode.ID) interface{} {
			return backend.PeerInfo(id)
		},
		Attributes:     []enr.Entry{eth.CurrentENREntry[T, P](backend.Chain())},
		DialCandidates: dnsdisc,
	}
	return protocols
}

func HandleServer[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T, P], peer *Peer[T, P]) error {
	for {
		if err := handleMessageServer(backend, peer); err != nil {
			return err
		}
	}
}

func handleMessageServer[T crypto.PrivateKey, P crypto.PublicKey](backend Backend[T, P], peer *Peer[T, P]) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := peer.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > maxMessageSize {
		return fmt.Errorf("%w: %v > %v", errMsgTooLarge, msg.Size, maxMessageSize)
	}
	defer msg.Discard()

	peer.Log().Info("QLight client message received", "msg", msg.Code)

	var handlers = map[uint64]eth.MsgHandler[T, P]{
		// old 64 messages
		eth.GetBlockHeadersMsg: eth.HandleGetBlockHeaders[T, P],
		eth.BlockHeadersMsg:    eth.HandleBlockHeaders[T, P],
		eth.GetBlockBodiesMsg:  eth.HandleGetBlockBodies[T, P],
		eth.BlockBodiesMsg:     eth.HandleBlockBodies[T, P],
		eth.NewBlockHashesMsg:  eth.HandleNewBlockhashes[T, P],
		eth.NewBlockMsg:        eth.HandleNewBlock[T, P],
		eth.TransactionsMsg:    eth.HandleTransactions[T, P],
		// New eth65 messages
		eth.NewPooledTransactionHashesMsg: eth.HandleNewPooledTransactionHashes[T, P],
		eth.GetPooledTransactionsMsg:      eth.HandleGetPooledTransactions[T, P],
		eth.PooledTransactionsMsg:         eth.HandlePooledTransactions[T, P],
	}

	switch msg.Code {
	case eth.BlockHeadersMsg:
		peer.Log().Info("QLight Block Headers message received. Ignoring.")
		return nil
	case eth.TransactionsMsg:
		peer.Log().Info("QLight Transactions message received. Ignoring.")
		return nil
	case QLightTokenUpdateMsg:
		peer.Log().Info("QLight Token update received.")
		res := new(qLightTokenUpdateData)
		if err := msg.Decode(res); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		peer.qlightToken = res.Token
		return nil
	case eth.GetBlockHeadersMsg:
		if handler := handlers[msg.Code]; handler != nil {
			return handler(backend, msg, peer.EthPeer)
		}
	case eth.GetBlockBodiesMsg:
		res := new(eth.GetBlockBodiesPacket)
		if err := msg.Decode(res); err != nil {
			return fmt.Errorf("%w: message %v: %v", errDecode, msg, err)
		}
		return backend.QHandle(peer, res)
	case eth.NewBlockHashesMsg:
		peer.Log().Info("QLight New Block Hashes message received. Ignoring.")
		return nil
	}
	peer.Log().Info("QLight Unable to find handler for received message", "msg", msg.Code)
	return fmt.Errorf("%w: %v", errInvalidMsgCode, msg.Code)
}
