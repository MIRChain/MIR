package p2p

import (
	"crypto/tls"
	"net"

	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/log"
	"github.com/MIRChain/MIR/p2p/rlpx"
)

var qlightTLSConfig *tls.Config

func SetQLightTLSConfig(config *tls.Config) {
	qlightTLSConfig = config
}

type tlsErrorTransport[T crypto.PrivateKey, P crypto.PublicKey] struct {
	err error
}

func (tr *tlsErrorTransport[T, P]) doEncHandshake(prv T) (P, error) {
	return crypto.ZeroPublicKey[P](), tr.err
}
func (tr *tlsErrorTransport[T, P]) doProtoHandshake(our *protoHandshake) (*protoHandshake, error) {
	return nil, tr.err
}
func (tr *tlsErrorTransport[T, P]) ReadMsg() (Msg, error) { return Msg{}, tr.err }
func (tr *tlsErrorTransport[T, P]) WriteMsg(Msg) error    { return tr.err }
func (tr *tlsErrorTransport[T, P]) close(err error)       {}

func NewQlightClientTransport[T crypto.PrivateKey, P crypto.PublicKey](conn net.Conn, dialDest P) transport[T, P] {
	log.Info("Setting up qlight client transport")
	if qlightTLSConfig != nil {
		tlsConn := tls.Client(conn, qlightTLSConfig)
		err := tlsConn.Handshake()
		if err != nil {
			log.Error("Failure setting up qlight client transport", "err", err)
			return &tlsErrorTransport[T, P]{err}
		}
		log.Info("Qlight client tls transport established successfully")
		return &rlpxTransport[T, P]{conn: rlpx.NewConn[T, P](tlsConn, dialDest)}
	}
	return &rlpxTransport[T, P]{conn: rlpx.NewConn[T, P](conn, dialDest)}
}

func NewQlightServerTransport[T crypto.PrivateKey, P crypto.PublicKey](conn net.Conn, dialDest P) transport[T, P] {
	log.Info("Setting up qlight server transport")
	if qlightTLSConfig != nil {
		tlsConn := tls.Server(conn, qlightTLSConfig)
		err := tlsConn.Handshake()
		if err != nil {
			log.Error("Failure setting up qlight server transport", "err", err)
			return &tlsErrorTransport[T, P]{err}
		}
		log.Info("Qlight server tls transport established successfully")
		return &rlpxTransport[T, P]{conn: rlpx.NewConn[T, P](tlsConn, dialDest)}
	}
	return &rlpxTransport[T, P]{conn: rlpx.NewConn[T, P](conn, dialDest)}
}
