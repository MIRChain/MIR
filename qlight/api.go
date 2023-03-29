package qlight

import (
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

type RunningPeerAuthUpdater interface {
	UpdateTokenForRunningQPeers(token string) error
}

type PrivateQLightAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	tokenHolder *TokenHolder[T,P]
	peerUpdater RunningPeerAuthUpdater
	rpcClient   *rpc.Client
}

// NewPublicEthereumAPI creates a new Ethereum protocol API for full nodes.
func NewPrivateQLightAPI[T crypto.PrivateKey, P crypto.PublicKey](peerUpdater RunningPeerAuthUpdater, rpcClient *rpc.Client) *PrivateQLightAPI[T,P] {
	return &PrivateQLightAPI[T,P]{peerUpdater: peerUpdater, rpcClient: rpcClient}
}

func (p *PrivateQLightAPI[T,P]) SetCurrentToken(token string) {
	p.tokenHolder.SetCurrentToken(token)
	p.peerUpdater.UpdateTokenForRunningQPeers(token)
	if p.rpcClient != nil {
		if len(token) > 0 {
			p.rpcClient.WithHTTPCredentials(p.tokenHolder.HttpCredentialsProvider)
		} else {
			p.rpcClient.WithHTTPCredentials(nil)
		}
	}
}

func (p *PrivateQLightAPI[T,P]) GetCurrentToken() string {
	return p.tokenHolder.CurrentToken()
}

func (p *PrivateQLightAPI[T,P]) ReloadPlugin() error {
	return p.tokenHolder.ReloadPlugin()
}
