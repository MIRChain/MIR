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

package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jpmorganchase/quorum-security-plugin-sdk-go/proto"
	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/bloombits"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/eth/gasprice"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/miner"
	"github.com/pavelkrolevets/MIR-pro/params"
	pcore "github.com/pavelkrolevets/MIR-pro/permission/core"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

// EthAPIBackend implements ethapi.Backend for full nodes
type EthAPIBackend [T crypto.PrivateKey, P crypto.PublicKey] struct {
	extRPCEnabled       bool
	allowUnprotectedTxs bool
	eth                 *Ethereum[T,P]
	gpo                 *gasprice.Oracle

	// Quorum
	//
	// hex node id from node public key
	hexNodeId string

	// timeout value for call
	evmCallTimeOut time.Duration
	// Quorum
	proxyClient *rpc.Client
}

// var _ ethapi.Backend = &EthAPIBackend{}
// var _ tracers.Backend = &EthAPIBackend{}

func (b *EthAPIBackend[T,P]) ProxyEnabled() bool {
	return b.eth.config.QuorumLightClient.Enabled()
}

func (b *EthAPIBackend[T,P]) ProxyClient() *rpc.Client {
	return b.proxyClient
}

// ChainConfig returns the active chain configuration.
func (b *EthAPIBackend[T,P]) ChainConfig() *params.ChainConfig {
	return b.eth.blockchain.Config()
}

// PSMR returns the private state metadata resolver.
func (b *EthAPIBackend[T,P]) PSMR() mps.PrivateStateMetadataResolver {
	return b.eth.blockchain.PrivateStateManager()
}

func (b *EthAPIBackend[T,P]) CurrentBlock() *types.Block[P] {
	return b.eth.blockchain.CurrentBlock()
}

func (b *EthAPIBackend[T,P]) SetHead(number uint64) {
	b.eth.handler.downloader.Cancel()
	b.eth.blockchain.SetHead(number)
}

func (b *EthAPIBackend[T,P]) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if number == rpc.PendingBlockNumber {
		block := b.eth.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if number == rpc.LatestBlockNumber {
		return b.eth.blockchain.CurrentBlock().Header(), nil
	}
	return b.eth.blockchain.GetHeaderByNumber(uint64(number)), nil
}

func (b *EthAPIBackend[T,P]) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.HeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.eth.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.eth.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		return header, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *EthAPIBackend[T,P]) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.eth.blockchain.GetHeaderByHash(hash), nil
}

func (b *EthAPIBackend[T,P]) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block[P], error) {
	// Pending block is only known by the miner
	if number == rpc.PendingBlockNumber {
		if b.eth.handler.raftMode {
			// Use latest instead.
			return b.eth.blockchain.CurrentBlock(), nil
		}
		block := b.eth.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if number == rpc.LatestBlockNumber {
		return b.eth.blockchain.CurrentBlock(), nil
	}
	return b.eth.blockchain.GetBlockByNumber(uint64(number)), nil
}

func (b *EthAPIBackend[T,P]) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block[P], error) {
	return b.eth.blockchain.GetBlockByHash(hash), nil
}

func (b *EthAPIBackend[T,P]) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block[P], error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.BlockByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header := b.eth.blockchain.GetHeaderByHash(hash)
		if header == nil {
			return nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.eth.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, errors.New("hash is not currently canonical")
		}
		block := b.eth.blockchain.GetBlock(hash, header.Number.Uint64())
		if block == nil {
			return nil, errors.New("header found, but block body is missing")
		}
		return block, nil
	}
	return nil, errors.New("invalid arguments; neither block nor hash specified")
}

func (b *EthAPIBackend[T,P]) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (vm.MinimalApiState, *types.Header, error) {
	psm, err := b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	// Pending state is only known by the miner
	if number == rpc.PendingBlockNumber {
		// Quorum
		if b.eth.handler.raftMode {
			// Use latest instead.
			header, err := b.HeaderByNumber(ctx, rpc.LatestBlockNumber)
			if header == nil || err != nil {
				return nil, nil, err
			}
			publicState, privateState, err := b.eth.BlockChain().StateAtPSI(header.Root, psm.ID)
			return EthAPIState{publicState, privateState}, header, err
		}
		block, publicState, privateState := b.eth.miner.Pending(psm.ID)
		if block == nil || publicState == nil || privateState == nil {
			return nil, nil, fmt.Errorf("Unable to retrieve the pending state from the miner.")
		}
		return EthAPIState{publicState, privateState}, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	if header == nil {
		return nil, nil, errors.New("header not found")
	}
	stateDb, privateState, err := b.eth.BlockChain().StateAtPSI(header.Root, psm.ID)
	return EthAPIState{stateDb, privateState}, header, err

}

func (b *EthAPIBackend[T,P]) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (vm.MinimalApiState, *types.Header, error) {
	if blockNr, ok := blockNrOrHash.Number(); ok {
		return b.StateAndHeaderByNumber(ctx, blockNr)
	}
	if hash, ok := blockNrOrHash.Hash(); ok {
		header, err := b.HeaderByHash(ctx, hash)
		if err != nil {
			return nil, nil, err
		}
		if header == nil {
			return nil, nil, errors.New("header for hash not found")
		}
		if blockNrOrHash.RequireCanonical && b.eth.blockchain.GetCanonicalHash(header.Number.Uint64()) != hash {
			return nil, nil, errors.New("hash is not currently canonical")
		}
		psm, err := b.PSMR().ResolveForUserContext(ctx)
		if err != nil {
			return nil, nil, err
		}
		stateDb, privateState, err := b.eth.BlockChain().StateAtPSI(header.Root, psm.ID)
		return EthAPIState{stateDb, privateState}, header, err

	}
	return nil, nil, errors.New("invalid arguments; neither block nor hash specified")
}

// Modified for Quorum:
// - If MPS is enabled then the list of receipts returned will contain all public receipts, plus the private receipts for this PSI.
// - if MPS is not enabled, then list will contain all public and private receipts
// Note that for a privacy marker transactions, the private receipts will remain under PSReceipts
func (b *EthAPIBackend[T,P]) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts[P], error) {
	receipts := b.eth.blockchain.GetReceiptsByHash(hash)
	psm, err := b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}

	psiReceipts := make([]*types.Receipt[P], len(receipts))
	for i := 0; i < len(receipts); i++ {
		psiReceipts[i] = receipts[i]
		if receipts[i].PSReceipts != nil {
			psReceipt, found := receipts[i].PSReceipts[psm.ID]
			// if PSReceipt found and this is not a privacy marker transaction receipt, then pull out the PSI receipt
			if found && receipts[i].TxHash == psReceipt.TxHash {
				psiReceipts[i] = psReceipt
			}
		}
	}

	return psiReceipts, nil
}

func (b *EthAPIBackend[T,P]) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	// Quorum
	// We should use the modified getReceipts to get the private receipts for PSI (MPS)
	receipts, err := b.GetReceipts(ctx, hash)
	if err != nil {
		return nil, err
	}
	// End Quorum
	if receipts == nil {
		return nil, nil
	}
	privateReceipts, err := b.eth.blockchain.GetPMTPrivateReceiptsByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	logs := make([][]*types.Log, len(receipts)+len(privateReceipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	for i, receipt := range privateReceipts {
		logs[len(receipts)+i] = receipt.Logs
	}
	return logs, nil
}

func (b *EthAPIBackend[T,P]) GetTd(ctx context.Context, hash common.Hash) *big.Int {
	return b.eth.blockchain.GetTdByHash(hash)
}

func (b *EthAPIBackend[T,P]) GetEVM(ctx context.Context, msg core.Message, state vm.MinimalApiState, header *types.Header, vmConfig *vm.Config[P]) (*vm.EVM[P], func() error, error) {
	statedb := state.(EthAPIState)
	vmError := func() error { return nil }
	if vmConfig == nil {
		vmConfig = b.eth.blockchain.GetVMConfig()
	}
	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext[P](header, b.eth.BlockChain(), nil)

	// Quorum
	// Set the private state to public state if contract address is not present in the private state
	to := common.Address{}
	if msg.To() != nil {
		to = *msg.To()
	}
	privateState := statedb.privateState
	if !privateState.Exist(to) {
		privateState = statedb.state
	}
	// End Quorum

	return vm.NewEVM(context, txContext, statedb.state, privateState, b.eth.blockchain.Config(), *vmConfig), vmError, nil
}

func (b *EthAPIBackend[T,P]) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent[P]) event.Subscription {
	return b.eth.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *EthAPIBackend[T,P]) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.eth.SubscribePendingLogs(ch) // Quorum
}

func (b *EthAPIBackend[T,P]) SubscribeChainEvent(ch chan<- core.ChainEvent[P]) event.Subscription {
	return b.eth.BlockChain().SubscribeChainEvent(ch)
}

func (b *EthAPIBackend[T,P]) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent[P]) event.Subscription {
	return b.eth.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *EthAPIBackend[T,P]) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent[P]) event.Subscription {
	return b.eth.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *EthAPIBackend[T,P]) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.eth.BlockChain().SubscribeLogsEvent(ch)
}

func (b *EthAPIBackend[T,P]) SendTx(ctx context.Context, signedTx *types.Transaction[P]) error {
	// validation for node need to happen here and cannot be done as a part of
	// validateTx in tx_pool.go as tx_pool validation will happen in every node
	if b.hexNodeId != "" && !pcore.ValidateNodeForTxn[P](b.hexNodeId, signedTx.From()) {
		return errors.New("cannot send transaction from this node")
	}
	return b.eth.txPool.AddLocal(signedTx)
}

func (b *EthAPIBackend[T,P]) GetPoolTransactions() (types.Transactions[P], error) {
	pending, err := b.eth.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions[P]
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *EthAPIBackend[T,P]) GetPoolTransaction(hash common.Hash) *types.Transaction[P] {
	return b.eth.txPool.Get(hash)
}

func (b *EthAPIBackend[T,P]) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction[P], common.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index := rawdb.ReadTransaction[P](b.eth.ChainDb(), txHash)
	return tx, blockHash, blockNumber, index, nil
}

func (b *EthAPIBackend[T,P]) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.eth.txPool.Nonce(addr), nil
}

func (b *EthAPIBackend[T,P]) Stats() (pending int, queued int) {
	return b.eth.txPool.Stats()
}

func (b *EthAPIBackend[T,P]) TxPoolContent() (map[common.Address]types.Transactions[P], map[common.Address]types.Transactions[P]) {
	return b.eth.TxPool().Content()
}

func (b *EthAPIBackend[T,P]) TxPool() *core.TxPool[P] {
	return b.eth.TxPool()
}

func (b *EthAPIBackend[T,P]) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent[P]) event.Subscription {
	return b.eth.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *EthAPIBackend[T,P]) Downloader() *downloader.Downloader {
	return b.eth.Downloader()
}

func (b *EthAPIBackend[T,P]) SuggestPrice(ctx context.Context) (*big.Int, error) {
	if !b.ChainConfig().IsQuorum || b.ChainConfig().IsGasPriceEnabled(b.eth.blockchain.CurrentBlock().Number()) {
		return b.gpo.SuggestPrice(ctx)
	} else {
		return big.NewInt(0), nil
	}
}

func (b *EthAPIBackend[T,P]) ChainDb() ethdb.Database {
	return b.eth.ChainDb()
}

func (b *EthAPIBackend[T,P]) EventMux() *event.TypeMux {
	return b.eth.EventMux()
}

func (b *EthAPIBackend[T,P]) AccountManager() *accounts.Manager[P] {
	return b.eth.AccountManager()
}

func (b *EthAPIBackend[T,P]) ExtRPCEnabled() bool {
	return b.extRPCEnabled
}

func (b *EthAPIBackend[T,P]) UnprotectedAllowed() bool {
	return b.allowUnprotectedTxs
}

func (b *EthAPIBackend[T,P]) RPCGasCap() uint64 {
	return b.eth.config.RPCGasCap
}

func (b *EthAPIBackend[T,P]) RPCTxFeeCap() float64 {
	return b.eth.config.RPCTxFeeCap
}

func (b *EthAPIBackend[T,P]) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.eth.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *EthAPIBackend[T,P]) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.eth.bloomRequests)
	}
}

func (b *EthAPIBackend[T,P]) Engine() consensus.Engine[P] {
	return b.eth.engine
}

func (b *EthAPIBackend[T,P]) CurrentHeader() *types.Header {
	return b.eth.blockchain.CurrentHeader()
}

func (b *EthAPIBackend[T,P]) Miner() *miner.Miner {
	return b.eth.Miner()
}

func (b *EthAPIBackend[T,P]) StartMining(threads int) error {
	return b.eth.StartMining(threads)
}

func (b *EthAPIBackend[T,P]) StateAtBlock(ctx context.Context, block *types.Block[P], reexec uint64, base *state.StateDB, checkLive bool) (*state.StateDB, mps.PrivateStateRepository, error) {
	return b.eth.stateAtBlock(block, reexec, base, checkLive)
}

func (b *EthAPIBackend[T,P]) StateAtTransaction(ctx context.Context, block *types.Block[P], txIndex int, reexec uint64) (core.Message, vm.BlockContext, *state.StateDB, *state.StateDB, mps.PrivateStateRepository, error) {
	return b.eth.stateAtTransaction(ctx, block, txIndex, reexec)
}

// Quorum

func (b *EthAPIBackend[T,P]) CallTimeOut() time.Duration {
	return b.evmCallTimeOut
}

func (b *EthAPIBackend[T,P]) GetBlockchain() *core.BlockChain[P] {
	return b.eth.BlockChain()
}

// The validation of pre-requisite for multitenancy is done when EthService
// is being created. So it's safe to use the config value here.
func (b *EthAPIBackend[T,P]) SupportsMultitenancy(rpcCtx context.Context) (*proto.PreAuthenticatedAuthenticationToken, bool) {
	authToken := rpc.PreauthenticatedTokenFromContext(rpcCtx)
	if authToken != nil && b.eth.config.MultiTenantEnabled() {
		return authToken, true
	}
	return nil, false
}

func (b *EthAPIBackend[T,P]) AccountExtraDataStateGetterByNumber(ctx context.Context, number rpc.BlockNumber) (vm.AccountExtraDataStateGetter, error) {
	s, _, err := b.StateAndHeaderByNumber(ctx, number)
	return s, err
}

func (b *EthAPIBackend[T,P]) IsPrivacyMarkerTransactionCreationEnabled() bool {
	return b.eth.config.QuorumChainConfig.PrivacyMarkerEnabled() && b.ChainConfig().IsPrivacyPrecompileEnabled(b.eth.blockchain.CurrentBlock().Number())
}

// used by Quorum
type EthAPIState struct {
	state, privateState *state.StateDB
}

func (s EthAPIState) GetBalance(addr common.Address) *big.Int {
	if s.privateState.Exist(addr) {
		return s.privateState.GetBalance(addr)
	}
	return s.state.GetBalance(addr)
}

func (s EthAPIState) GetCode(addr common.Address) []byte {
	if s.privateState.Exist(addr) {
		return s.privateState.GetCode(addr)
	}
	return s.state.GetCode(addr)
}

func (s EthAPIState) SetNonce(addr common.Address, nonce uint64) {
	if s.privateState.Exist(addr) {
		s.privateState.SetNonce(addr, nonce)
	} else {
		s.state.SetNonce(addr, nonce)
	}
}

func (s EthAPIState) SetCode(addr common.Address, code []byte) {
	if s.privateState.Exist(addr) {
		s.privateState.SetCode(addr, code)
	} else {
		s.state.SetCode(addr, code)
	}
}

func (s EthAPIState) SetBalance(addr common.Address, balance *big.Int) {
	if s.privateState.Exist(addr) {
		s.privateState.SetBalance(addr, balance)
	} else {
		s.state.SetBalance(addr, balance)
	}
}

func (s EthAPIState) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	if s.privateState.Exist(addr) {
		s.privateState.SetStorage(addr, storage)
	} else {
		s.state.SetStorage(addr, storage)
	}
}

func (s EthAPIState) SetState(a common.Address, key common.Hash, value common.Hash) {
	if s.privateState.Exist(a) {
		s.privateState.SetState(a, key, value)
	} else {
		s.state.SetState(a, key, value)
	}
}

func (s EthAPIState) GetState(a common.Address, b common.Hash) common.Hash {
	if s.privateState.Exist(a) {
		return s.privateState.GetState(a, b)
	}
	return s.state.GetState(a, b)
}

func (s EthAPIState) GetNonce(addr common.Address) uint64 {
	if s.privateState.Exist(addr) {
		return s.privateState.GetNonce(addr)
	}
	return s.state.GetNonce(addr)
}

func (s EthAPIState) GetPrivacyMetadata(addr common.Address) (*state.PrivacyMetadata, error) {
	if s.privateState.Exist(addr) {
		return s.privateState.GetPrivacyMetadata(addr)
	}
	return nil, fmt.Errorf("%x: %w", addr, common.ErrNotPrivateContract)
}

func (s EthAPIState) GetManagedParties(addr common.Address) ([]string, error) {
	if s.privateState.Exist(addr) {
		return s.privateState.GetManagedParties(addr)
	}
	return nil, fmt.Errorf("%x: %w", addr, common.ErrNotPrivateContract)
}

func (s EthAPIState) GetRLPEncodedStateObject(addr common.Address) ([]byte, error) {
	getFunc := s.state.GetRLPEncodedStateObject
	if s.privateState.Exist(addr) {
		getFunc = s.privateState.GetRLPEncodedStateObject
	}
	return getFunc(addr)
}

func (s EthAPIState) GetProof(addr common.Address) ([][]byte, error) {
	if s.privateState.Exist(addr) {
		return s.privateState.GetProof(addr)
	}
	return s.state.GetProof(addr)
}

func (s EthAPIState) GetStorageProof(addr common.Address, h common.Hash) ([][]byte, error) {
	if s.privateState.Exist(addr) {
		return s.privateState.GetStorageProof(addr, h)
	}
	return s.state.GetStorageProof(addr, h)
}

func (s EthAPIState) StorageTrie(addr common.Address) state.Trie {
	if s.privateState.Exist(addr) {
		return s.privateState.StorageTrie(addr)
	}
	return s.state.StorageTrie(addr)
}

func (s EthAPIState) Error() error {
	if s.privateState.Error() != nil {
		return s.privateState.Error()
	}
	return s.state.Error()
}

func (s EthAPIState) GetCodeHash(addr common.Address) common.Hash {
	if s.privateState.Exist(addr) {
		return s.privateState.GetCodeHash(addr)
	}
	return s.state.GetCodeHash(addr)
}
