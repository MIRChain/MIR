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

package backends

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/MIRChain/MIR"
	"github.com/MIRChain/MIR/accounts/abi"
	"github.com/MIRChain/MIR/accounts/abi/bind"
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/common/hexutil"
	"github.com/MIRChain/MIR/common/math"
	"github.com/MIRChain/MIR/consensus/ethash"
	"github.com/MIRChain/MIR/core"
	"github.com/MIRChain/MIR/core/bloombits"
	"github.com/MIRChain/MIR/core/mps"
	"github.com/MIRChain/MIR/core/rawdb"
	"github.com/MIRChain/MIR/core/state"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/core/vm"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/eth"
	"github.com/MIRChain/MIR/eth/filters"
	"github.com/MIRChain/MIR/ethdb"
	"github.com/MIRChain/MIR/event"
	"github.com/MIRChain/MIR/log"
	"github.com/MIRChain/MIR/params"
	"github.com/MIRChain/MIR/rpc"
)

// This nil assignment ensures at compile time that SimulatedBackend implements bind.ContractBackend.
// var _ bind.ContractBackend = (*SimulatedBackend[P])(nil)

var (
	errBlockNumberUnsupported  = errors.New("simulatedBackend cannot access blocks other than the latest block")
	errBlockDoesNotExist       = errors.New("block does not exist in blockchain")
	errTransactionDoesNotExist = errors.New("transaction does not exist")
)

// SimulatedBackend implements bind.ContractBackend, simulating a blockchain in
// the background. Its main purpose is to allow for easy testing of contract bindings.
// Simulated backend implements the following interfaces:
// ChainReader, ChainStateReader, ContractBackend, ContractCaller, ContractFilterer, ContractTransactor,
// DeployBackend, GasEstimator, GasPricer, LogFilterer, PendingContractCaller, TransactionReader, and TransactionSender
type SimulatedBackend[P crypto.PublicKey] struct {
	database   ethdb.Database      // In memory database to store our testing data
	blockchain *core.BlockChain[P] // Ethereum blockchain to handle the consensus

	mu           sync.Mutex
	pendingBlock *types.Block[P]   // Currently pending block that will be imported on request
	pendingState *state.StateDB[P] // Currently pending state that will be the active on request

	events *filters.EventSystem[P] // Event system for filtering log events live

	config *params.ChainConfig
}

// NewSimulatedBackendWithDatabase creates a new binding backend based on the given database
// and uses a simulated blockchain for testing purposes.
// A simulated backend always uses chainID 1337.
func NewSimulatedBackendWithDatabase[P crypto.PublicKey](database ethdb.Database, alloc core.GenesisAlloc, gasLimit uint64) *SimulatedBackend[P] {
	genesis := core.Genesis[P]{Config: params.AllEthashProtocolChanges, GasLimit: gasLimit, Alloc: alloc}
	genesis.MustCommit(database)
	blockchain, _ := core.NewBlockChain[P](database, nil, genesis.Config, ethash.NewFaker[P](), vm.Config[P]{}, nil, nil, nil)

	backend := &SimulatedBackend[P]{
		database:   database,
		blockchain: blockchain,
		config:     genesis.Config,
		events:     filters.NewEventSystem[P](&filterBackend[P]{database, blockchain}, false),
	}
	backend.rollback()
	return backend
}

// Quorum
//
// Create a simulated backend based on existing Ethereum service
func NewSimulatedBackendFrom[T crypto.PrivateKey, P crypto.PublicKey](ethereum *eth.Ethereum[T, P]) *SimulatedBackend[P] {
	backend := &SimulatedBackend[P]{
		database:   ethereum.ChainDb(),
		blockchain: ethereum.BlockChain(),
		config:     ethereum.BlockChain().Config(),
		events:     filters.NewEventSystem[P](&filterBackend[P]{ethereum.ChainDb(), ethereum.BlockChain()}, false),
	}
	backend.rollback()
	return backend
}

// NewSimulatedBackend creates a new binding backend using a simulated blockchain
// for testing purposes.
// A simulated backend always uses chainID 1337.
func NewSimulatedBackend[P crypto.PublicKey](alloc core.GenesisAlloc, gasLimit uint64) *SimulatedBackend[P] {
	return NewSimulatedBackendWithDatabase[P](rawdb.NewMemoryDatabase(), alloc, gasLimit)
}

// Close terminates the underlying blockchain's update loop.
func (b *SimulatedBackend[P]) Close() error {
	b.blockchain.Stop()
	return nil
}

// Commit imports all the pending transactions as a single block and starts a
// fresh new state.
func (b *SimulatedBackend[P]) Commit() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, err := b.blockchain.InsertChain([]*types.Block[P]{b.pendingBlock}); err != nil {
		panic(err) // This cannot happen unless the simulator is wrong, fail in that case
	}
	b.rollback()
}

// Rollback aborts all pending transactions, reverting to the last committed state.
func (b *SimulatedBackend[P]) Rollback() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.rollback()
}

func (b *SimulatedBackend[P]) rollback() {
	blocks, _ := core.GenerateChain[P](b.config, b.blockchain.CurrentBlock(), ethash.NewFaker[P](), b.database, 1, func(int, *core.BlockGen[P]) {})

	b.pendingBlock = blocks[0]
	b.pendingState, _ = state.New[P](b.pendingBlock.Root(), b.blockchain.StateCache(), nil)
}

// stateByBlockNumber retrieves a state by a given blocknumber.
func (b *SimulatedBackend[P]) stateByBlockNumber(ctx context.Context, blockNumber *big.Int) (*state.StateDB[P], error) {
	if blockNumber == nil || blockNumber.Cmp(b.blockchain.CurrentBlock().Number()) == 0 {
		statedb, _, err := b.blockchain.State()
		return statedb, err
	}
	block, err := b.blockByNumberNoLock(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	statedb, _, err := b.blockchain.StateAt(block.Root())
	return statedb, err
}

// CodeAt returns the code associated with a certain account in the blockchain.
func (b *SimulatedBackend[P]) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	stateDB, _, _ = b.blockchain.State()
	return stateDB.GetCode(contract), nil
}

// BalanceAt returns the wei balance of a certain account in the blockchain.
func (b *SimulatedBackend[P]) BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	stateDB, _, _ = b.blockchain.State()
	return stateDB.GetBalance(contract), nil
}

// NonceAt returns the nonce of a certain account in the blockchain.
func (b *SimulatedBackend[P]) NonceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return 0, err
	}
	stateDB, _, _ = b.blockchain.State()
	return stateDB.GetNonce(contract), nil
}

// StorageAt returns the value of key in the storage of an account in the blockchain.
func (b *SimulatedBackend[P]) StorageAt(ctx context.Context, contract common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateDB, err := b.stateByBlockNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	stateDB, _, _ = b.blockchain.State()
	val := stateDB.GetState(contract, key)
	return val[:], nil
}

// TransactionReceipt returns the receipt of a transaction.
func (b *SimulatedBackend[P]) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt[P], error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	receipt, _, _, _ := rawdb.ReadReceipt[P](b.database, txHash, b.config)
	return receipt, nil
}

// TransactionByHash checks the pool of pending transactions in addition to the
// blockchain. The isPending return value indicates whether the transaction has been
// mined yet. Note that the transaction may not be part of the canonical chain even if
// it's not pending.
func (b *SimulatedBackend[P]) TransactionByHash(ctx context.Context, txHash common.Hash) (*types.Transaction[P], bool, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	tx := b.pendingBlock.Transaction(txHash)
	if tx != nil {
		return tx, true, nil
	}
	tx, _, _, _ = rawdb.ReadTransaction[P](b.database, txHash)
	if tx != nil {
		return tx, false, nil
	}
	return nil, false, ethereum.NotFound
}

// BlockByHash retrieves a block based on the block hash.
func (b *SimulatedBackend[P]) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block[P], error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if hash == b.pendingBlock.Hash() {
		return b.pendingBlock, nil
	}

	block := b.blockchain.GetBlockByHash(hash)
	if block != nil {
		return block, nil
	}

	return nil, errBlockDoesNotExist
}

// BlockByNumber retrieves a block from the database by number, caching it
// (associated with its hash) if found.
func (b *SimulatedBackend[P]) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block[P], error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.blockByNumberNoLock(ctx, number)
}

// blockByNumberNoLock retrieves a block from the database by number, caching it
// (associated with its hash) if found without Lock.
func (b *SimulatedBackend[P]) blockByNumberNoLock(ctx context.Context, number *big.Int) (*types.Block[P], error) {
	if number == nil || number.Cmp(b.pendingBlock.Number()) == 0 {
		return b.blockchain.CurrentBlock(), nil
	}

	block := b.blockchain.GetBlockByNumber(uint64(number.Int64()))
	if block == nil {
		return nil, errBlockDoesNotExist
	}

	return block, nil
}

// HeaderByHash returns a block header from the current canonical chain.
func (b *SimulatedBackend[P]) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header[P], error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if hash == b.pendingBlock.Hash() {
		return b.pendingBlock.Header(), nil
	}

	header := b.blockchain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errBlockDoesNotExist
	}

	return header, nil
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (b *SimulatedBackend[P]) HeaderByNumber(ctx context.Context, block *big.Int) (*types.Header[P], error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if block == nil || block.Cmp(b.pendingBlock.Number()) == 0 {
		return b.blockchain.CurrentHeader(), nil
	}

	return b.blockchain.GetHeaderByNumber(uint64(block.Int64())), nil
}

// TransactionCount returns the number of transactions in a given block.
func (b *SimulatedBackend[P]) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if blockHash == b.pendingBlock.Hash() {
		return uint(b.pendingBlock.Transactions().Len()), nil
	}

	block := b.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return uint(0), errBlockDoesNotExist
	}

	return uint(block.Transactions().Len()), nil
}

// TransactionInBlock returns the transaction for a specific block at a specific index.
func (b *SimulatedBackend[P]) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction[P], error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if blockHash == b.pendingBlock.Hash() {
		transactions := b.pendingBlock.Transactions()
		if uint(len(transactions)) < index+1 {
			return nil, errTransactionDoesNotExist
		}

		return transactions[index], nil
	}

	block := b.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return nil, errBlockDoesNotExist
	}

	transactions := block.Transactions()
	if uint(len(transactions)) < index+1 {
		return nil, errTransactionDoesNotExist
	}

	return transactions[index], nil
}

// PendingCodeAt returns the code associated with an account in the pending state.
func (b *SimulatedBackend[P]) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.pendingState.GetCode(contract), nil
}

func newRevertError[P crypto.PublicKey](result *core.ExecutionResult) *revertError {
	reason, errUnpack := abi.UnpackRevert[P](result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// revertError is an API error that encompasses an EVM revert with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revert.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

// CallContract executes a contract call.
func (b *SimulatedBackend[P]) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if blockNumber != nil && blockNumber.Cmp(b.blockchain.CurrentBlock().Number()) != 0 {
		return nil, errBlockNumberUnsupported
	}
	stateDB, _, err := b.blockchain.State()
	if err != nil {
		return nil, err
	}
	res, err := b.callContract(ctx, call, b.blockchain.CurrentBlock(), stateDB, stateDB)
	if err != nil {
		return nil, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(res.Revert()) > 0 {
		return nil, newRevertError[P](res)
	}
	return res.Return(), res.Err
}

// PendingCallContract executes a contract call on the pending state.
func (b *SimulatedBackend[P]) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	defer b.pendingState.RevertToSnapshot(b.pendingState.Snapshot())

	res, err := b.callContract(ctx, call, b.pendingBlock, b.pendingState, b.pendingState)
	if err != nil {
		return nil, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(res.Revert()) > 0 {
		return nil, newRevertError[P](res)
	}
	return res.Return(), res.Err
}

// PendingNonceAt implements PendingStateReader.PendingNonceAt, retrieving
// the nonce currently pending for the account.
func (b *SimulatedBackend[P]) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.pendingState.GetOrNewStateObject(account).Nonce(), nil
}

// SuggestGasPrice implements ContractTransactor.SuggestGasPrice. Since the simulated
// chain doesn't have miners, we just return a gas price of 1 for any call.
func (b *SimulatedBackend[P]) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}

// EstimateGas executes the requested code against the currently pending block/state and
// returns the used amount of gas.
func (b *SimulatedBackend[P]) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Determine the lowest and highest possible gas limits to binary search in between
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	if call.Gas >= params.TxGas {
		hi = call.Gas
	} else {
		hi = b.pendingBlock.GasLimit()
	}
	// Recap the highest gas allowance with account's balance.
	if call.GasPrice != nil && call.GasPrice.BitLen() != 0 {
		balance := b.pendingState.GetBalance(call.From) // from can't be nil
		available := new(big.Int).Set(balance)
		if call.Value != nil {
			if call.Value.Cmp(available) >= 0 {
				return 0, errors.New("insufficient funds for transfer")
			}
			available.Sub(available, call.Value)
		}
		allowance := new(big.Int).Div(available, call.GasPrice)
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := call.Value
			if transfer == nil {
				transfer = new(big.Int)
			}
			log.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", transfer, "gasprice", call.GasPrice, "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) (bool, *core.ExecutionResult, error) {
		call.Gas = gas

		snapshot := b.pendingState.Snapshot()
		res, err := b.callContract(ctx, call, b.pendingBlock, b.pendingState, b.pendingState)
		b.pendingState.RevertToSnapshot(snapshot)

		if err != nil {
			if errors.Is(err, core.ErrIntrinsicGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			return true, nil, err // Bail out
		}
		return res.Failed(), res, nil
	}
	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := executable(mid)

		// If the error is not nil(consensus error), it means the provided message
		// call or transaction will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		failed, result, err := executable(hi)
		if err != nil {
			return 0, err
		}
		if failed {
			if result != nil && result.Err != vm.ErrOutOfGas {
				if len(result.Revert()) > 0 {
					return 0, newRevertError[P](result)
				}
				return 0, result.Err
			}
			// Otherwise, the specified gas cap is too low
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}
	return hi, nil
}

// callContract implements common code between normal and pending contract calls.
// state is modified during execution, make sure to copy it if necessary.
func (b *SimulatedBackend[P]) callContract(ctx context.Context, call ethereum.CallMsg, block *types.Block[P], stateDB *state.StateDB[P], privateState *state.StateDB[P]) (*core.ExecutionResult, error) {
	// Ensure message is initialized properly.
	if call.GasPrice == nil {
		call.GasPrice = big.NewInt(1)
	}
	if call.Gas == 0 {
		call.Gas = 50000000
	}
	if call.Value == nil {
		call.Value = new(big.Int)
	}
	// Set infinite balance to the fake caller account.
	from := stateDB.GetOrNewStateObject(call.From)
	from.SetBalance(math.MaxBig256)
	// Execute the call.
	msg := callMsg{call}

	txContext := core.NewEVMTxContext(msg)
	evmContext := core.NewEVMBlockContext[P](block.Header(), b.blockchain, nil)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmEnv := vm.NewEVM[P](evmContext, txContext, stateDB, privateState, b.config, vm.Config[P]{})
	gasPool := new(core.GasPool).AddGas(math.MaxUint64)

	return core.NewStateTransition(vmEnv, msg, gasPool).TransitionDb()
}

// SendTransaction updates the pending block to include the given transaction.
// It panics if the transaction is invalid.
func (b *SimulatedBackend[P]) SendTransaction(ctx context.Context, tx *types.Transaction[P], args bind.PrivateTxArgs) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check transaction validity.
	block := b.blockchain.CurrentBlock()
	signer := types.MakeSigner[P](b.blockchain.Config(), block.Number())
	sender, err := types.Sender(signer, tx)
	if err != nil {
		panic(fmt.Errorf("invalid transaction: %v", err))
	}
	nonce := b.pendingState.GetNonce(sender)
	if tx.Nonce() != nonce {
		panic(fmt.Errorf("invalid transaction nonce: got %d, want %d", tx.Nonce(), nonce))
	}

	// Include tx in chain.
	blocks, _ := core.GenerateChain[P](b.config, block, ethash.NewFaker[P](), b.database, 1, func(number int, block *core.BlockGen[P]) {
		for _, tx := range b.pendingBlock.Transactions() {
			block.AddTxWithChain(b.blockchain, tx)
		}
		block.AddTxWithChain(b.blockchain, tx)
	})
	stateDB, _, _ := b.blockchain.State()

	b.pendingBlock = blocks[0]
	b.pendingState, _ = state.New[P](b.pendingBlock.Root(), stateDB.Database(), nil)
	return nil
}

// PreparePrivateTransaction dummy implementation
func (b *SimulatedBackend[P]) PreparePrivateTransaction(data []byte, privateFrom string) (common.EncryptedPayloadHash, error) {
	return common.EncryptedPayloadHash{}, nil
}

func (b *SimulatedBackend[P]) DistributeTransaction(ctx context.Context, tx *types.Transaction[P], args bind.PrivateTxArgs) (string, error) {
	return tx.Hash().String(), nil
}

// FilterLogs executes a log filter operation, blocking during execution and
// returning all the results in one batch.
//
// TODO(karalabe): Deprecate when the subscription one can return past data too.
func (b *SimulatedBackend[P]) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	var filter *filters.Filter[P]
	if query.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = filters.NewBlockFilter[P](&filterBackend[P]{b.database, b.blockchain}, *query.BlockHash, query.Addresses, query.Topics, query.PSI)
	} else {
		// Initialize unset filter boundaries to run from genesis to chain head
		from := int64(0)
		if query.FromBlock != nil {
			from = query.FromBlock.Int64()
		}
		to := int64(-1)
		if query.ToBlock != nil {
			to = query.ToBlock.Int64()
		}
		// Construct the range filter
		filter = filters.NewRangeFilter[P](&filterBackend[P]{b.database, b.blockchain}, from, to, query.Addresses, query.Topics, query.PSI)
	}
	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]types.Log, len(logs))
	for i, nLog := range logs {
		res[i] = *nLog
	}
	return res, nil
}

// SubscribeFilterLogs creates a background log filtering operation, returning a
// subscription immediately, which can be used to stream the found events.
func (b *SimulatedBackend[P]) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	// Subscribe to contract events
	sink := make(chan []*types.Log)

	sub, err := b.events.SubscribeLogs(query, sink)
	if err != nil {
		return nil, err
	}
	// Since we're getting logs in batches, we need to flatten them into a plain stream
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case logs := <-sink:
				for _, nlog := range logs {
					select {
					case ch <- *nlog:
					case err := <-sub.Err():
						return err
					case <-quit:
						return nil
					}
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// SubscribeNewHead returns an event subscription for a new header.
func (b *SimulatedBackend[P]) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header[P]) (ethereum.Subscription, error) {
	// subscribe to a new head
	sink := make(chan *types.Header[P])
	sub := b.events.SubscribeNewHeads(sink)

	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case head := <-sink:
				select {
				case ch <- head:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// AdjustTime adds a time shift to the simulated clock.
// It can only be called on empty blocks.
func (b *SimulatedBackend[P]) AdjustTime(adjustment time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.pendingBlock.Transactions()) != 0 {
		return errors.New("Could not adjust time on non-empty block")
	}

	blocks, _ := core.GenerateChain[P](b.config, b.blockchain.CurrentBlock(), ethash.NewFaker[P](), b.database, 1, func(number int, block *core.BlockGen[P]) {
		block.OffsetTime(int64(adjustment.Seconds()))
	})
	stateDB, _, _ := b.blockchain.State()

	b.pendingBlock = blocks[0]
	b.pendingState, _ = state.New[P](b.pendingBlock.Root(), stateDB.Database(), nil)

	return nil
}

// Blockchain returns the underlying blockchain.
func (b *SimulatedBackend[P]) Blockchain() *core.BlockChain[P] {
	return b.blockchain
}

// callMsg implements core.Message to allow passing it as a transaction simulator.
type callMsg struct {
	ethereum.CallMsg
}

func (m callMsg) From() common.Address         { return m.CallMsg.From }
func (m callMsg) Nonce() uint64                { return 0 }
func (m callMsg) CheckNonce() bool             { return false }
func (m callMsg) To() *common.Address          { return m.CallMsg.To }
func (m callMsg) GasPrice() *big.Int           { return m.CallMsg.GasPrice }
func (m callMsg) Gas() uint64                  { return m.CallMsg.Gas }
func (m callMsg) Value() *big.Int              { return m.CallMsg.Value }
func (m callMsg) Data() []byte                 { return m.CallMsg.Data }
func (m callMsg) AccessList() types.AccessList { return m.CallMsg.AccessList }

// filterBackend implements filters.Backend to support filtering for logs without
// taking bloom-bits acceleration structures into account.
type filterBackend[P crypto.PublicKey] struct {
	db ethdb.Database
	bc *core.BlockChain[P]
}

func (fb *filterBackend[P]) ChainDb() ethdb.Database  { return fb.db }
func (fb *filterBackend[P]) EventMux() *event.TypeMux { panic("not supported") }

func (fb *filterBackend[P]) PSMR() mps.PrivateStateMetadataResolver {
	return fb.bc.PrivateStateManager()
}

func (fb *filterBackend[P]) HeaderByNumber(ctx context.Context, block rpc.BlockNumber) (*types.Header[P], error) {
	if block == rpc.LatestBlockNumber {
		return fb.bc.CurrentHeader(), nil
	}
	return fb.bc.GetHeaderByNumber(uint64(block.Int64())), nil
}

func (fb *filterBackend[P]) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header[P], error) {
	return fb.bc.GetHeaderByHash(hash), nil
}

func (fb *filterBackend[P]) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts[P], error) {
	number := rawdb.ReadHeaderNumber(fb.db, hash)
	if number == nil {
		return nil, nil
	}
	return rawdb.ReadReceipts[P](fb.db, hash, *number, fb.bc.Config()), nil
}

func (fb *filterBackend[P]) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := rawdb.ReadHeaderNumber(fb.db, hash)
	if number == nil {
		return nil, nil
	}
	receipts := rawdb.ReadReceipts[P](fb.db, hash, *number, fb.bc.Config())
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (fb *filterBackend[P]) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent[P]) event.Subscription {
	return nullSubscription()
}

func (fb *filterBackend[P]) SubscribeChainEvent(ch chan<- core.ChainEvent[P]) event.Subscription {
	return fb.bc.SubscribeChainEvent(ch)
}

func (fb *filterBackend[P]) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent[P]) event.Subscription {
	return fb.bc.SubscribeRemovedLogsEvent(ch)
}

func (fb *filterBackend[P]) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return fb.bc.SubscribeLogsEvent(ch)
}

func (fb *filterBackend[P]) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return nullSubscription()
}

func (fb *filterBackend[P]) BloomStatus() (uint64, uint64) { return 4096, 0 }

func (fb *filterBackend[P]) ServiceFilter(ctx context.Context, ms *bloombits.MatcherSession) {
	panic("not supported")
}

func (fb *filterBackend[P]) AccountExtraDataStateGetterByNumber(context.Context, rpc.BlockNumber) (vm.AccountExtraDataStateGetter, error) {
	panic("not supported")
}

func nullSubscription() event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		<-quit
		return nil
	})
}
