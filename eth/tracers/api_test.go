// Copyright 2021 The go-ethereum Authors
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

package tracers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

var (
	errStateNotFound       = errors.New("state not found")
	errBlockNotFound       = errors.New("block not found")
	errTransactionNotFound = errors.New("transaction not found")
)

type testBackend [P crypto.PublicKey] struct {
	chainConfig *params.ChainConfig
	engine      consensus.Engine[P]
	chaindb     ethdb.Database
	chain       *core.BlockChain[P]
}

var _ Backend[nist.PublicKey] = &testBackend[nist.PublicKey]{}

func newTestBackend[P crypto.PublicKey](t *testing.T, n int, gspec *core.Genesis[P], generator func(i int, b *core.BlockGen[P])) *testBackend[P] {
	backend := &testBackend[P]{
		chainConfig: params.TestChainConfig,
		engine:       ethash.NewFaker[P](),
		chaindb:     rawdb.NewMemoryDatabase(),
	}
	// Generate blocks for testing
	gspec.Config = backend.chainConfig
	var (
		gendb   = rawdb.NewMemoryDatabase()
		genesis = gspec.MustCommit(gendb)
	)
	blocks, _ := core.GenerateChain[P](backend.chainConfig, genesis, backend.engine, gendb, n, generator)

	// Import the canonical chain
	gspec.MustCommit(backend.chaindb)
	cacheConfig := &core.CacheConfig{
		TrieCleanLimit:    256,
		TrieDirtyLimit:    256,
		TrieTimeLimit:     5 * time.Minute,
		SnapshotLimit:     0,
		TrieDirtyDisabled: true, // Archive mode
	}
	chain, err := core.NewBlockChain[P](backend.chaindb, cacheConfig, backend.chainConfig, backend.engine, vm.Config[P]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	backend.chain = chain
	return backend
}

func (b *testBackend[P]) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header[P], error) {
	return b.chain.GetHeaderByHash(hash), nil
}

func (b *testBackend[P]) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header[P], error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentHeader(), nil
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}

func (b *testBackend[P]) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block[P], error) {
	return b.chain.GetBlockByHash(hash), nil
}

func (b *testBackend[P]) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block[P], error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock(), nil
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}

func (b *testBackend[P]) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction[P], common.Hash, uint64, uint64, error) {
	tx, hash, blockNumber, index := rawdb.ReadTransaction[P](b.chaindb, txHash)
	if tx == nil {
		return nil, common.Hash{}, 0, 0, errTransactionNotFound
	}
	return tx, hash, blockNumber, index, nil
}

func (b *testBackend[P]) RPCGasCap() uint64 {
	return 25000000
}

func (b *testBackend[P]) ChainConfig() *params.ChainConfig {
	return b.chainConfig
}

func (b *testBackend[P]) Engine() consensus.Engine[P] {
	return b.engine
}

func (b *testBackend[P]) ChainDb() ethdb.Database {
	return b.chaindb
}

func (b *testBackend[P]) StateAtBlock(ctx context.Context, block *types.Block[P], reexec uint64, base *state.StateDB[P], checkLive bool) (*state.StateDB[P], mps.PrivateStateRepository[P], error) {
	statedb, privateStateRepo, err := b.chain.StateAt(block.Root())
	if err != nil {
		return nil, nil, errStateNotFound
	}
	return statedb, privateStateRepo, nil
}

func (b *testBackend[P]) StateAtTransaction(ctx context.Context, block *types.Block[P], txIndex int, reexec uint64) (core.Message, vm.BlockContext, *state.StateDB[P], *state.StateDB[P], mps.PrivateStateRepository[P], error) {
	parent := b.chain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.BlockContext{}, nil, nil, nil, errBlockNotFound
	}
	statedb, privateStateRepo, err := b.chain.StateAt(parent.Root())
	if err != nil {
		return nil, vm.BlockContext{}, nil, nil, nil, errStateNotFound
	}
	psm, err := b.chain.PrivateStateManager().ResolveForUserContext(ctx)
	if err != nil {
		return nil, vm.BlockContext{}, nil, nil, nil, err
	}
	privateState, err := privateStateRepo.StatePSI(psm.ID)
	if err != nil {
		return nil, vm.BlockContext{}, nil, nil, nil, errStateNotFound
	}
	if txIndex == 0 && len(block.Transactions()) == 0 {
		return nil, vm.BlockContext{}, statedb, privateState, privateStateRepo, nil
	}
	// Recompute transactions up to the target index.
	signer := types.MakeSigner[P](b.chainConfig, block.Number())
	for idx, tx := range block.Transactions() {
		msg, _ := tx.AsMessage(signer)
		txContext := core.NewEVMTxContext(msg)
		context := core.NewEVMBlockContext[P](block.Header(), b.chain, nil)
		if idx == txIndex {
			return msg, context, statedb, privateState, privateStateRepo, nil
		}
		vmenv := vm.NewEVM(context, txContext, statedb, privateState, b.chainConfig, vm.Config[nist.PublicKey]{})
		if _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(tx.Gas())); err != nil {
			return nil, vm.BlockContext{}, nil, nil, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}
		statedb.Finalise(vmenv.ChainConfig().IsEIP158(block.Number()))
	}
	return nil, vm.BlockContext{}, nil, nil, nil, fmt.Errorf("transaction index %d out of range for block %#x", txIndex, block.Hash())
}

func (b *testBackend[P]) GetBlockchain() *core.BlockChain[P] {
	return b.chain
}

func TestTraceCall(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &core.Genesis[nist.PublicKey]{Alloc: core.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.Ether)},
		accounts[1].addr: {Balance: big.NewInt(params.Ether)},
		accounts[2].addr: {Balance: big.NewInt(params.Ether)},
	}}
	genBlocks := 10
	signer := types.HomesteadSigner[nist.PublicKey]{}
	api := NewAPI[nist.PublicKey](newTestBackend[nist.PublicKey](t, genBlocks, genesis, func(i int, b *core.BlockGen[nist.PublicKey]) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
		b.AddTx(tx)
	}))

	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		call        ethapi.CallArgs
		config      *TraceCallConfig[nist.PublicKey]
		expectErr   error
		expect      interface{}
	}{
		// Standard JSON trace upon the genesis, plain transfer.
		{
			blockNumber: rpc.BlockNumber(0),
			call: ethapi.CallArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect: &ethapi.ExecutionResult{
				Gas:         params.TxGas,
				Failed:      false,
				ReturnValue: "",
				StructLogs:  []ethapi.StructLogRes{},
			},
		},
		// Standard JSON trace upon the head, plain transfer.
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			call: ethapi.CallArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect: &ethapi.ExecutionResult{
				Gas:         params.TxGas,
				Failed:      false,
				ReturnValue: "",
				StructLogs:  []ethapi.StructLogRes{},
			},
		},
		// Standard JSON trace upon the non-existent block, error expects
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			call: ethapi.CallArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: fmt.Errorf("block #%d not found", genBlocks+1),
			expect:    nil,
		},
		// Standard JSON trace upon the latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.CallArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect: &ethapi.ExecutionResult{
				Gas:         params.TxGas,
				Failed:      false,
				ReturnValue: "",
				StructLogs:  []ethapi.StructLogRes{},
			},
		},
		// Standard JSON trace upon the pending block
		{
			blockNumber: rpc.PendingBlockNumber,
			call: ethapi.CallArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect: &ethapi.ExecutionResult{
				Gas:         params.TxGas,
				Failed:      false,
				ReturnValue: "",
				StructLogs:  []ethapi.StructLogRes{},
			},
		},
	}
	for _, testspec := range testSuite {
		tc := &TraceCallConfig[nist.PublicKey]{}
		if testspec.config != nil {
			tc.LogConfig = testspec.config.LogConfig
			tc.Tracer = testspec.config.Tracer
			tc.Timeout = testspec.config.Timeout
			tc.Reexec = testspec.config.Reexec
		}
		result, err := api.TraceCall(context.Background(), testspec.call, rpc.BlockNumberOrHash{BlockNumber: &testspec.blockNumber}, tc)
		if testspec.expectErr != nil {
			if err == nil {
				t.Errorf("Expect error %v, get nothing", testspec.expectErr)
				continue
			}
			if !reflect.DeepEqual(err, testspec.expectErr) {
				t.Errorf("Error mismatch, want %v, get %v", testspec.expectErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("Expect no error, get %v", err)
				continue
			}
			if !reflect.DeepEqual(result, testspec.expect) {
				t.Errorf("Result mismatch, want %v, get %v", testspec.expect, result)
			}
		}
	}
}

func TestOverridenTraceCall(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &core.Genesis[nist.PublicKey]{Alloc: core.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.Ether)},
		accounts[1].addr: {Balance: big.NewInt(params.Ether)},
		accounts[2].addr: {Balance: big.NewInt(params.Ether)},
	}}
	genBlocks := 10
	signer := types.HomesteadSigner[nist.PublicKey]{}
	api := NewAPI[nist.PublicKey](newTestBackend[nist.PublicKey](t, genBlocks, genesis, func(i int, b *core.BlockGen[nist.PublicKey]) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
		b.AddTx(tx)
	}))
	randomAccounts, tracer := newAccounts(3), "callTracer"

	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		call        ethapi.CallArgs
		config      *TraceCallConfig[nist.PublicKey]
		expectErr   error
		expect      *callTrace
	}{
		// Succcessful call with state overriding
		{
			blockNumber: rpc.PendingBlockNumber,
			call: ethapi.CallArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config: &TraceCallConfig[nist.PublicKey]{
				Tracer: &tracer,
				StateOverrides: &ethapi.StateOverride[nist.PublicKey]{
					randomAccounts[0].addr: ethapi.OverrideAccount{Balance: newRPCBalance(new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether)))},
				},
			},
			expectErr: nil,
			expect: &callTrace{
				Type:    "CALL",
				From:    randomAccounts[0].addr,
				To:      randomAccounts[1].addr,
				Gas:     newRPCUint64(24979000),
				GasUsed: newRPCUint64(0),
				Value:   (*hexutil.Big)(big.NewInt(1000)),
			},
		},
		// Invalid call without state overriding
		{
			blockNumber: rpc.PendingBlockNumber,
			call: ethapi.CallArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config: &TraceCallConfig[nist.PublicKey]{
				Tracer: &tracer,
			},
			expectErr: core.ErrInsufficientFundsForTransfer,
			expect:    nil,
		},
		// Successful simple contract call
		//
		// // SPDX-License-Identifier: GPL-3.0
		//
		//  pragma solidity >=0.7.0 <0.8.0;
		//
		//  /**
		//   * @title Storage
		//   * @dev Store & retrieve value in a variable
		//   */
		//  contract Storage {
		//      uint256 public number;
		//      constructor() {
		//          number = block.number;
		//      }
		//  }
		{
			blockNumber: rpc.PendingBlockNumber,
			call: ethapi.CallArgs{
				From: &randomAccounts[0].addr,
				To:   &randomAccounts[2].addr,
				Data: newRPCBytes(common.Hex2Bytes("8381f58a")), // call number()
			},
			config: &TraceCallConfig[nist.PublicKey]{
				Tracer: &tracer,
				StateOverrides: &ethapi.StateOverride[nist.PublicKey]{
					randomAccounts[2].addr: ethapi.OverrideAccount{
						Code:      newRPCBytes(common.Hex2Bytes("6080604052348015600f57600080fd5b506004361060285760003560e01c80638381f58a14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000548156fea2646970667358221220eab35ffa6ab2adfe380772a48b8ba78e82a1b820a18fcb6f59aa4efb20a5f60064736f6c63430007040033")),
						StateDiff: newStates([]common.Hash{{}}, []common.Hash{common.BigToHash(big.NewInt(123))}),
					},
				},
			},
			expectErr: nil,
			expect: &callTrace{
				Type:    "CALL",
				From:    randomAccounts[0].addr,
				To:      randomAccounts[2].addr,
				Input:   hexutil.Bytes(common.Hex2Bytes("8381f58a")),
				Output:  hexutil.Bytes(common.BigToHash(big.NewInt(123)).Bytes()),
				Gas:     newRPCUint64(24978936),
				GasUsed: newRPCUint64(2283),
				Value:   (*hexutil.Big)(big.NewInt(0)),
			},
		},
	}
	for _, testspec := range testSuite {
		result, err := api.TraceCall(context.Background(), testspec.call, rpc.BlockNumberOrHash{BlockNumber: &testspec.blockNumber}, testspec.config)
		if testspec.expectErr != nil {
			if err == nil {
				t.Errorf("Expect error %v, get nothing", testspec.expectErr)
				continue
			}
			if !errors.Is(err, testspec.expectErr) {
				t.Errorf("Error mismatch, want %v, get %v", testspec.expectErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("Expect no error, get %v", err)
				continue
			}
			ret := new(callTrace)
			if err := json.Unmarshal(result.(json.RawMessage), ret); err != nil {
				t.Fatalf("failed to unmarshal trace result: %v", err)
			}
			if !jsonEqual(ret, testspec.expect) {
				// uncomment this for easier debugging
				//have, _ := json.MarshalIndent(ret, "", " ")
				//want, _ := json.MarshalIndent(testspec.expect, "", " ")
				//t.Fatalf("trace mismatch: \nhave %+v\nwant %+v", string(have), string(want))
				t.Fatalf("trace mismatch: \nhave %+v\nwant %+v", ret, testspec.expect)
			}
		}
	}
}

func TestTraceTransaction(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(2)
	genesis := &core.Genesis[nist.PublicKey]{Alloc: core.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.Ether)},
		accounts[1].addr: {Balance: big.NewInt(params.Ether)},
	}}
	target := common.Hash{}
	signer := types.HomesteadSigner[nist.PublicKey]{}
	api := NewAPI[nist.PublicKey](newTestBackend[nist.PublicKey](t, 1, genesis, func(i int, b *core.BlockGen[nist.PublicKey]) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
		b.AddTx(tx)
		target = tx.Hash()
	}))
	result, err := api.TraceTransaction(context.Background(), target, nil)
	if err != nil {
		t.Errorf("Failed to trace transaction %v", err)
	}
	if !reflect.DeepEqual(result, &ethapi.ExecutionResult{
		Gas:         params.TxGas,
		Failed:      false,
		ReturnValue: "",
		StructLogs:  []ethapi.StructLogRes{},
	}) {
		t.Error("Transaction tracing result is different")
	}
}

func TestTraceBlock(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &core.Genesis[nist.PublicKey]{Alloc: core.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.Ether)},
		accounts[1].addr: {Balance: big.NewInt(params.Ether)},
		accounts[2].addr: {Balance: big.NewInt(params.Ether)},
	}}
	genBlocks := 10
	signer := types.HomesteadSigner[nist.PublicKey]{}
	api := NewAPI[nist.PublicKey](newTestBackend[nist.PublicKey](t, genBlocks, genesis, func(i int, b *core.BlockGen[nist.PublicKey]) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, big.NewInt(0), nil), signer, accounts[0].key)
		b.AddTx(tx)
	}))

	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		config      *TraceConfig
		expect      interface{}
		expectErr   error
	}{
		// Trace genesis block, expect error
		{
			blockNumber: rpc.BlockNumber(0),
			config:      nil,
			expect:      nil,
			expectErr:   errors.New("genesis is not traceable"),
		},
		// Trace head block
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{
					Result: &ethapi.ExecutionResult{
						Gas:         params.TxGas,
						Failed:      false,
						ReturnValue: "",
						StructLogs:  []ethapi.StructLogRes{},
					},
				},
			},
		},
		// Trace non-existent block
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			config:      nil,
			expectErr:   fmt.Errorf("block #%d not found", genBlocks+1),
			expect:      nil,
		},
		// Trace latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{
					Result: &ethapi.ExecutionResult{
						Gas:         params.TxGas,
						Failed:      false,
						ReturnValue: "",
						StructLogs:  []ethapi.StructLogRes{},
					},
				},
			},
		},
		// Trace pending block
		{
			blockNumber: rpc.PendingBlockNumber,
			config:      nil,
			expectErr:   nil,
			expect: []*txTraceResult{
				{
					Result: &ethapi.ExecutionResult{
						Gas:         params.TxGas,
						Failed:      false,
						ReturnValue: "",
						StructLogs:  []ethapi.StructLogRes{},
					},
				},
			},
		},
	}
	for _, testspec := range testSuite {
		result, err := api.TraceBlockByNumber(context.Background(), testspec.blockNumber, testspec.config)
		if testspec.expectErr != nil {
			if err == nil {
				t.Errorf("Expect error %v, get nothing", testspec.expectErr)
				continue
			}
			if !reflect.DeepEqual(err, testspec.expectErr) {
				t.Errorf("Error mismatch, want %v, get %v", testspec.expectErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("Expect no error, get %v", err)
				continue
			}
			if !reflect.DeepEqual(result, testspec.expect) {
				t.Errorf("Result mismatch, want %v, get %v", testspec.expect, result)
			}
		}
	}
}

type Account struct {
	key  nist.PrivateKey
	addr common.Address
}

type Accounts []Account

func (a Accounts) Len() int           { return len(a) }
func (a Accounts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Accounts) Less(i, j int) bool { return bytes.Compare(a[i].addr.Bytes(), a[j].addr.Bytes()) < 0 }

func newAccounts(n int) (accounts Accounts) {
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey[nist.PrivateKey]()
		addr := crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		accounts = append(accounts, Account{key: key, addr: addr})
	}
	sort.Sort(accounts)
	return accounts
}

func newRPCBalance(balance *big.Int) **hexutil.Big {
	rpcBalance := (*hexutil.Big)(balance)
	return &rpcBalance
}

func newRPCUint64(number uint64) *hexutil.Uint64 {
	rpcUint64 := hexutil.Uint64(number)
	return &rpcUint64
}

func newRPCBytes(bytes []byte) *hexutil.Bytes {
	rpcBytes := hexutil.Bytes(bytes)
	return &rpcBytes
}

func newStates(keys []common.Hash, vals []common.Hash) *map[common.Hash]common.Hash {
	if len(keys) != len(vals) {
		panic("invalid input")
	}
	m := make(map[common.Hash]common.Hash)
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = vals[i]
	}
	return &m
}
