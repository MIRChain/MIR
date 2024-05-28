package qlight

import (
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/eth/protocols/eth"
)

// Handler is a callback to invoke from an outside runner after the boilerplate
// exchanges have passed.
type Handler[T crypto.PrivateKey, P crypto.PublicKey] func(peer *Peer[T, P]) error

// Backend defines the data retrieval methods to serve remote requests and the
// callback methods to invoke on remote deliveries.
type Backend[T crypto.PrivateKey, P crypto.PublicKey] interface {
	eth.Backend[T, P]
	// RunPeer is invoked when a peer joins on the `eth` protocol. The handler
	// should do any peer maintenance work, handshakes and validations. If all
	// is passed, control should be given back to the `handler` to process the
	// inbound messages going forward.
	RunQPeer(peer *Peer[T, P], handler Handler[T, P]) error
	// Handle is a callback to be invoked when a data packet is received from
	// the remote peer. Only packets not consumed by the protocol handler will
	// be forwarded to the backend.
	QHandle(peer *Peer[T, P], packet eth.Packet) error
}
