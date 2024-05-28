// Copyright 2015 The go-ethereum Authors
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

// Package ethapi implements the general Ethereum API functions.
package ethapi

import (
	"context"
	"math/big"
	"time"

	"github.com/MIRChain/MIR/accounts"
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/consensus"
	"github.com/MIRChain/MIR/core"
	"github.com/MIRChain/MIR/core/bloombits"
	"github.com/MIRChain/MIR/core/mps"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/core/vm"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/eth/downloader"
	"github.com/MIRChain/MIR/ethdb"
	"github.com/MIRChain/MIR/event"
	"github.com/MIRChain/MIR/params"
	"github.com/MIRChain/MIR/rpc"
	"github.com/jpmorganchase/quorum-security-plugin-sdk-go/proto"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend[T crypto.PrivateKey, P crypto.PublicKey] interface {
	// General Ethereum API
	Downloader() *downloader.Downloader[T, P]
	SuggestPrice(ctx context.Context) (*big.Int, error)
	ChainDb() ethdb.Database
	AccountManager() *accounts.Manager[P]
	ExtRPCEnabled() bool
	RPCGasCap() uint64        // global gas cap for eth_call over rpc: DoS protection
	RPCTxFeeCap() float64     // global tx fee cap for all transaction related APIs
	UnprotectedAllowed() bool // allows only for EIP155 transactions.

	// Blockchain API
	SetHead(number uint64)
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header[P], error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header[P], error)
	HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header[P], error)
	CurrentHeader() *types.Header[P]
	CurrentBlock() *types.Block[P]
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block[P], error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block[P], error)
	BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block[P], error)
	StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (vm.MinimalApiState, *types.Header[P], error)
	StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (vm.MinimalApiState, *types.Header[P], error)
	GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts[P], error)
	GetTd(ctx context.Context, hash common.Hash) *big.Int
	GetEVM(ctx context.Context, msg core.Message, state vm.MinimalApiState, header *types.Header[P], vmConfig *vm.Config[P]) (*vm.EVM[P], func() error, error)
	SubscribeChainEvent(ch chan<- core.ChainEvent[P]) event.Subscription
	SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent[P]) event.Subscription
	SubscribeChainSideEvent(ch chan<- core.ChainSideEvent[P]) event.Subscription

	// Transaction pool API
	SendTx(ctx context.Context, signedTx *types.Transaction[P]) error
	GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction[P], common.Hash, uint64, uint64, error)
	GetPoolTransactions() (types.Transactions[P], error)
	GetPoolTransaction(txHash common.Hash) *types.Transaction[P]
	GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error)
	Stats() (pending int, queued int)
	TxPoolContent() (map[common.Address]types.Transactions[P], map[common.Address]types.Transactions[P])
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent[P]) event.Subscription

	// Filter API
	BloomStatus() (uint64, uint64)
	GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error)
	ServiceFilter(ctx context.Context, session *bloombits.MatcherSession)
	SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription
	SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription
	SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent[P]) event.Subscription

	ChainConfig() *params.ChainConfig
	Engine() consensus.Engine[P]

	// Quorum
	CallTimeOut() time.Duration
	// AccountExtraDataStateGetterByNumber returns state getter at a given block height
	AccountExtraDataStateGetterByNumber(ctx context.Context, number rpc.BlockNumber) (vm.AccountExtraDataStateGetter, error)
	PSMR() mps.PrivateStateMetadataResolver
	SupportsMultitenancy(rpcCtx context.Context) (*proto.PreAuthenticatedAuthenticationToken, bool)
	// IsPrivacyMarkerTransactionCreationEnabled returns true if privacy marker transactions are enabled and should be created
	IsPrivacyMarkerTransactionCreationEnabled() bool
}

func GetAPIs[T crypto.PrivateKey, P crypto.PublicKey](apiBackend Backend[T, P]) []rpc.API {
	nonceLock := new(AddrLocker)
	return []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicTransactionPoolProxyAPI(apiBackend, nonceLock),
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(apiBackend),
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(apiBackend.AccountManager()),
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountProxyAPI[T, P](apiBackend, nonceLock),
			Public:    false,
		},
	}
}
