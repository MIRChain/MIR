package raft

import (
	"crypto"
	"crypto/ecdsa"
	"sync"
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto/csp"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/eth"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

type RaftService [T crypto.PrivateKey]  struct {
	blockchain     *core.BlockChain
	chainDb        ethdb.Database // Block chain database
	txMu           sync.Mutex
	txPool         *core.TxPool
	accountManager *accounts.Manager
	downloader     *downloader.Downloader

	raftProtocolManager *ProtocolManager[T,P][T]
	startPeers          []*enode.Node

	// we need an event mux to instantiate the blockchain
	eventMux         *event.TypeMux
	minter           *minter[T]
	nodeKey          T
	calcGasLimitFunc func(block *types.Block) uint64

	pendingLogsFeed *event.Feed
}

func New [T ecdsa.PrivateKey | gost3410.PrivateKey | csp.Cert ] (stack *node.Node, chainConfig *params.ChainConfig, raftId, raftPort uint16, joinExisting bool, blockTime time.Duration, e *eth.Ethereum, startPeers []*enode.Node, raftLogDir string, useDns bool) (*RaftService[T], error) {
	service := &RaftService[T]{
		eventMux:         stack.EventMux(),
		chainDb:          e.ChainDb(),
		blockchain:       e.BlockChain(),
		txPool:           e.TxPool(),
		accountManager:   e.AccountManager(),
		downloader:       e.Downloader(),
		startPeers:       startPeers,
		nodeKey:		  stack.GetNodeKey(),
		calcGasLimitFunc: e.CalcGasLimit,
		pendingLogsFeed:  e.ConsensusServicePendingLogsFeed(),
	}

	etherbase, _ := e.Etherbase() //Quorum
	service.minter = newMinter(chainConfig, service, blockTime, etherbase)

	var err error
	if service.raftProtocolManager, err = NewProtocolManager(raftId, raftPort, service.blockchain, service.eventMux, startPeers, joinExisting, raftLogDir, service.minter, service.downloader, useDns, stack.Server()); err != nil {
		return nil, err
	}

	stack.RegisterAPIs(service.apis())
	stack.RegisterLifecycle(service)

	return service, nil
}

// Utility methods

func (service *RaftService[T]) apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "raft",
			Version:   "1.0",
			Service:   NewPublicRaftAPI(service),
			Public:    true,
		},
	}
}

// Backend interface methods:

func (service *RaftService[T]) AccountManager() *accounts.Manager { return service.accountManager }
func (service *RaftService[T]) BlockChain() *core.BlockChain      { return service.blockchain }
func (service *RaftService[T]) ChainDb() ethdb.Database           { return service.chainDb }
func (service *RaftService[T]) DappDb() ethdb.Database            { return nil }
func (service *RaftService[T]) EventMux() *event.TypeMux          { return service.eventMux }
func (service *RaftService[T]) TxPool() *core.TxPool              { return service.txPool }

// node.Lifecycle interface methods:

// Start implements node.Service, starting the background data propagation thread
// of the protocol.
func (service *RaftService[T]) Start() error {
	service.raftProtocolManager.Start()
	return nil
}

// Stop implements node.Service, stopping the background data propagation thread
// of the protocol.
func (service *RaftService[T]) Stop() error {
	service.blockchain.Stop()
	service.raftProtocolManager.Stop()
	service.minter.stop()
	service.eventMux.Stop()

	// handles gracefully if freezedb process is already stopped
	service.chainDb.Close()

	log.Info("Raft stopped")
	return nil
}
