package raft

import (
	"sync"
	"time"

	"github.com/MIRChain/MIR/accounts"
	"github.com/MIRChain/MIR/core"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/eth"
	"github.com/MIRChain/MIR/eth/downloader"
	"github.com/MIRChain/MIR/ethdb"
	"github.com/MIRChain/MIR/event"
	"github.com/MIRChain/MIR/log"
	"github.com/MIRChain/MIR/node"
	"github.com/MIRChain/MIR/p2p/enode"
	"github.com/MIRChain/MIR/params"
	"github.com/MIRChain/MIR/rpc"
)

type RaftService[T crypto.PrivateKey, P crypto.PublicKey] struct {
	blockchain     *core.BlockChain[P]
	chainDb        ethdb.Database // Block chain database
	txMu           sync.Mutex
	txPool         *core.TxPool[P]
	accountManager *accounts.Manager[P]
	downloader     *downloader.Downloader[T, P]

	raftProtocolManager *ProtocolManager[T, P]
	startPeers          []*enode.Node[P]

	// we need an event mux to instantiate the blockchain
	eventMux         *event.TypeMux
	minter           *minter[T, P]
	nodeKey          T
	calcGasLimitFunc func(block *types.Block[P]) uint64

	pendingLogsFeed *event.Feed
}

func New[T crypto.PrivateKey, P crypto.PublicKey](stack *node.Node[T, P], chainConfig *params.ChainConfig, raftId, raftPort uint16, joinExisting bool, blockTime time.Duration, e *eth.Ethereum[T, P], startPeers []*enode.Node[P], raftLogDir string, useDns bool) (*RaftService[T, P], error) {
	service := &RaftService[T, P]{
		eventMux:         stack.EventMux(),
		chainDb:          e.ChainDb(),
		blockchain:       e.BlockChain(),
		txPool:           e.TxPool(),
		accountManager:   e.AccountManager(),
		downloader:       e.Downloader(),
		startPeers:       startPeers,
		nodeKey:          stack.GetNodeKey(),
		calcGasLimitFunc: e.CalcGasLimit,
		pendingLogsFeed:  e.ConsensusServicePendingLogsFeed(),
	}

	etherbase, _ := e.Etherbase() //Quorum
	service.minter = newMinter[T, P](chainConfig, service, blockTime, etherbase)

	var err error
	if service.raftProtocolManager, err = NewProtocolManager(raftId, raftPort, service.blockchain, service.eventMux, startPeers, joinExisting, raftLogDir, service.minter, service.downloader, useDns, stack.Server()); err != nil {
		return nil, err
	}

	stack.RegisterAPIs(service.apis())
	stack.RegisterLifecycle(service)

	return service, nil
}

// Utility methods

func (service *RaftService[T, P]) apis() []rpc.API {
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

func (service *RaftService[T, P]) AccountManager() *accounts.Manager[P] {
	return service.accountManager
}
func (service *RaftService[T, P]) BlockChain() *core.BlockChain[P] { return service.blockchain }
func (service *RaftService[T, P]) ChainDb() ethdb.Database         { return service.chainDb }
func (service *RaftService[T, P]) DappDb() ethdb.Database          { return nil }
func (service *RaftService[T, P]) EventMux() *event.TypeMux        { return service.eventMux }
func (service *RaftService[T, P]) TxPool() *core.TxPool[P]         { return service.txPool }

// node.Lifecycle interface methods:

// Start implements node.Service, starting the background data propagation thread
// of the protocol.
func (service *RaftService[T, P]) Start() error {
	service.raftProtocolManager.Start()
	return nil
}

// Stop implements node.Service, stopping the background data propagation thread
// of the protocol.
func (service *RaftService[T, P]) Stop() error {
	service.blockchain.Stop()
	service.raftProtocolManager.Stop()
	service.minter.stop()
	service.eventMux.Stop()

	// handles gracefully if freezedb process is already stopped
	service.chainDb.Close()

	log.Info("Raft stopped")
	return nil
}
