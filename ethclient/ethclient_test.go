// Copyright 2016 The go-ethereum Authors
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

package ethclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	ethereum "github.com/pavelkrolevets/MIR-pro"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/eth"
	"github.com/pavelkrolevets/MIR-pro/eth/ethconfig"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/stretchr/testify/assert"
)

// Verify that Client implements the ethereum interfaces.
var (
	_ = ethereum.ChainReader[nist.PublicKey](&Client[nist.PublicKey]{})
	_ = ethereum.TransactionReader[nist.PublicKey](&Client[nist.PublicKey]{})
	_ = ethereum.ChainStateReader(&Client[nist.PublicKey]{})
	_ = ethereum.ChainSyncReader(&Client[nist.PublicKey]{})
	_ = ethereum.ContractCaller(&Client[nist.PublicKey]{})
	_ = ethereum.GasEstimator(&Client[nist.PublicKey]{})
	_ = ethereum.GasPricer(&Client[nist.PublicKey]{})
	_ = ethereum.LogFilterer(&Client[nist.PublicKey]{})
	_ = ethereum.PendingStateReader(&Client[nist.PublicKey]{})
	// _ = ethereum.PendingStateEventer(&Client{})
	_ = ethereum.PendingContractCaller(&Client[nist.PublicKey]{})
)

func TestToFilterArg(t *testing.T) {
	blockHashErr := fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
	addresses := []common.Address{
		common.HexToAddress("0xD36722ADeC3EdCB29c8e7b5a47f352D701393462"),
	}
	blockHash := common.HexToHash(
		"0xeb94bb7d78b73657a9d7a99792413f50c0a45c51fc62bdcb08a53f18e9a2b4eb",
	)

	for _, testCase := range []struct {
		name   string
		input  ethereum.FilterQuery
		output interface{}
		err    error
	}{
		{
			"without BlockHash",
			ethereum.FilterQuery{
				Addresses: addresses,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x1",
				"toBlock":   "0x2",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with nil fromBlock and nil toBlock",
			ethereum.FilterQuery{
				Addresses: addresses,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x0",
				"toBlock":   "latest",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with negative fromBlock and negative toBlock",
			ethereum.FilterQuery{
				Addresses: addresses,
				FromBlock: big.NewInt(-1),
				ToBlock:   big.NewInt(-1),
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "pending",
				"toBlock":   "pending",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"blockHash": blockHash,
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash and from block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				ToBlock:   big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and both from / to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := toFilterArg(testCase.input)
			if (testCase.err == nil) != (err == nil) {
				t.Fatalf("expected error %v but got %v", testCase.err, err)
			}
			if testCase.err != nil {
				if testCase.err.Error() != err.Error() {
					t.Fatalf("expected error %v but got %v", testCase.err, err)
				}
			} else if !reflect.DeepEqual(testCase.output, output) {
				t.Fatalf("expected filter arg %v but got %v", testCase.output, output)
			}
		})
	}
}

var (
	testKey, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr    = crypto.PubkeyToAddress[nist.PublicKey](*testKey.Public())
	testBalance = big.NewInt(2e10)
)

func newTestBackend(t *testing.T) (*node.Node[nist.PrivateKey,nist.PublicKey], []*types.Block[nist.PublicKey]) {
	// Generate test chain.
	genesis, blocks := generateTestChain()
	// Create node
	n, err := node.New(&node.Config[nist.PrivateKey,nist.PublicKey]{})
	if err != nil {
		t.Fatalf("can't create new node: %v", err)
	}
	// Create Ethereum Service
	config := &ethconfig.Config[nist.PublicKey]{Genesis: genesis}
	config.Ethash.PowMode = ethash.ModeFake
	ethservice, err := eth.New(n, config)
	if err != nil {
		t.Fatalf("can't create new ethereum service: %v", err)
	}
	// Import the test chain.
	if err := n.Start(); err != nil {
		t.Fatalf("can't start test node: %v", err)
	}
	if _, err := ethservice.BlockChain().InsertChain(blocks[1:]); err != nil {
		t.Fatalf("can't import test blocks: %v", err)
	}
	return n, blocks
}

func generateTestChain() (*core.Genesis[nist.PublicKey], []*types.Block[nist.PublicKey]) {
	db := rawdb.NewMemoryDatabase()
	config := params.AllEthashProtocolChanges
	genesis := &core.Genesis[nist.PublicKey]{
		Config:    config,
		Alloc:     core.GenesisAlloc{testAddr: {Balance: testBalance}},
		ExtraData: []byte("test genesis"),
		Timestamp: 9000,
	}
	generate := func(i int, g *core.BlockGen[nist.PublicKey]) {
		g.OffsetTime(5)
		g.SetExtra([]byte("test"))
	}
	gblock := genesis.ToBlock(db)
	engine :=  ethash.NewFaker[nist.PublicKey]()
	blocks, _ := core.GenerateChain[nist.PublicKey](config, gblock, engine, db, 1, generate)
	blocks = append([]*types.Block[nist.PublicKey]{gblock}, blocks...)
	return genesis, blocks
}

func TestEthClient(t *testing.T) {
	backend, chain := newTestBackend(t)
	client, _ := backend.Attach()
	defer backend.Close()
	defer client.Close()

	tests := map[string]struct {
		test func(t *testing.T)
	}{
		"TestHeader": {
			func(t *testing.T) { testHeader(t, chain, client) },
		},
		"TestBalanceAt": {
			func(t *testing.T) { testBalanceAt[nist.PublicKey](t, client) },
		},
		"TestTxInBlockInterrupted": {
			func(t *testing.T) { testTransactionInBlockInterrupted[nist.PublicKey](t, client) },
		},
		"TestChainID": {
			func(t *testing.T) { testChainID[nist.PublicKey](t, client) },
		},
		"TestGetBlock": {
			func(t *testing.T) { testGetBlock[nist.PublicKey](t, client) },
		},
		"TestStatusFunctions": {
			func(t *testing.T) { testStatusFunctions[nist.PublicKey](t, client) },
		},
		"TestCallContract": {
			func(t *testing.T) { testCallContract[nist.PublicKey](t, client) },
		},
		"TestAtFunctions": {
			func(t *testing.T) { testAtFunctions[nist.PublicKey](t, client) },
		},
	}

	t.Parallel()
	for name, tt := range tests {
		t.Run(name, tt.test)
	}
}

func testHeader[P crypto.PublicKey](t *testing.T, chain []*types.Block[P], client *rpc.Client) {
	tests := map[string]struct {
		block   *big.Int
		want    *types.Header[P]
		wantErr error
	}{
		"genesis": {
			block: big.NewInt(0),
			want:  chain[0].Header(),
		},
		"first_block": {
			block: big.NewInt(1),
			want:  chain[1].Header(),
		},
		"future_block": {
			block:   big.NewInt(1000000000),
			want:    nil,
			wantErr: ethereum.NotFound,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ec := NewClient[P](client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.HeaderByNumber(ctx, tt.block)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("HeaderByNumber(%v) error = %q, want %q", tt.block, err, tt.wantErr)
			}
			if got != nil && got.Number != nil && got.Number.Sign() == 0 {
				got.Number = big.NewInt(0) // hack to make DeepEqual work
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("HeaderByNumber(%v)\n   = %v\nwant %v", tt.block, got, tt.want)
			}
		})
	}
}

func testBalanceAt[P crypto.PublicKey](t *testing.T, client *rpc.Client) {
	tests := map[string]struct {
		account common.Address
		block   *big.Int
		want    *big.Int
		wantErr error
	}{
		"valid_account": {
			account: testAddr,
			block:   big.NewInt(1),
			want:    testBalance,
		},
		"non_existent_account": {
			account: common.Address{1},
			block:   big.NewInt(1),
			want:    big.NewInt(0),
		},
		"future_block": {
			account: testAddr,
			block:   big.NewInt(1000000000),
			want:    big.NewInt(0),
			wantErr: errors.New("header not found"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ec := NewClient[P](client)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := ec.BalanceAt(ctx, tt.account, tt.block)
			if tt.wantErr != nil && (err == nil || err.Error() != tt.wantErr.Error()) {
				t.Fatalf("BalanceAt(%x, %v) error = %q, want %q", tt.account, tt.block, err, tt.wantErr)
			}
			if got.Cmp(tt.want) != 0 {
				t.Fatalf("BalanceAt(%x, %v) = %v, want %v", tt.account, tt.block, got, tt.want)
			}
		})
	}
}

func testTransactionInBlockInterrupted[P crypto.PublicKey](t *testing.T, client *rpc.Client) {
	ec := NewClient[P](client)

	// Get current block by number
	block, err := ec.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Test tx in block interupted
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	tx, err := ec.TransactionInBlock(ctx, block.Hash(), 1)
	if tx != nil {
		t.Fatal("transaction should be nil")
	}
	if err == nil || err == ethereum.NotFound {
		t.Fatal("error should not be nil/notfound")
	}
	// Test tx in block not found
	if _, err := ec.TransactionInBlock(context.Background(), block.Hash(), 1); err != ethereum.NotFound {
		t.Fatal("error should be ethereum.NotFound")
	}
}

func testChainID[P crypto.PublicKey](t *testing.T, client *rpc.Client) {
	ec := NewClient[P](client)
	id, err := ec.ChainID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || id.Cmp(params.AllEthashProtocolChanges.ChainID) != 0 {
		t.Fatalf("ChainID returned wrong number: %+v", id)
	}
}

func testGetBlock[P crypto.PublicKey](t *testing.T, client *rpc.Client) {
	ec := NewClient[P](client)
	// Get current block number
	blockNumber, err := ec.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blockNumber != 1 {
		t.Fatalf("BlockNumber returned wrong number: %d", blockNumber)
	}
	// Get current block by number
	block, err := ec.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.NumberU64() != blockNumber {
		t.Fatalf("BlockByNumber returned wrong block: want %d got %d", blockNumber, block.NumberU64())
	}
	// Get current block by hash
	blockH, err := ec.BlockByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Hash() != blockH.Hash() {
		t.Fatalf("BlockByHash returned wrong block: want %v got %v", block.Hash().Hex(), blockH.Hash().Hex())
	}
	// Get header by number
	header, err := ec.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != header.Hash() {
		t.Fatalf("HeaderByNumber returned wrong header: want %v got %v", block.Header().Hash().Hex(), header.Hash().Hex())
	}
	// Get header by hash
	headerH, err := ec.HeaderByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != headerH.Hash() {
		t.Fatalf("HeaderByHash returned wrong header: want %v got %v", block.Header().Hash().Hex(), headerH.Hash().Hex())
	}
}

func testStatusFunctions[P crypto.PublicKey](t *testing.T, client *rpc.Client) {
	ec := NewClient[P](client)

	// Sync progress
	progress, err := ec.SyncProgress(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progress != nil {
		t.Fatalf("unexpected progress: %v", progress)
	}
	// NetworkID
	networkID, err := ec.NetworkID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if networkID.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("unexpected networkID: %v", networkID)
	}
	// SuggestGasPrice (should suggest 1 Gwei)
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gasPrice.Cmp(big.NewInt(1000000000)) != 0 {
		t.Fatalf("unexpected gas price: %v", gasPrice)
	}
}

func testCallContract[P crypto.PublicKey](t *testing.T, client *rpc.Client) {
	ec := NewClient[P](client)

	// EstimateGas
	msg := ethereum.CallMsg{
		From:     testAddr,
		To:       &common.Address{},
		Gas:      21000,
		GasPrice: big.NewInt(1),
		Value:    big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	// CallContract
	if _, err := ec.CallContract(context.Background(), msg, big.NewInt(1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PendingCallCOntract
	if _, err := ec.PendingCallContract(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testAtFunctions[P crypto.PublicKey](t *testing.T, client *rpc.Client) {
	ec := NewClient[P](client)
	// send a transaction for some interesting pending status
	sendTransaction(ec)
	time.Sleep(100 * time.Millisecond)
	// Check pending transaction count
	pending, err := ec.PendingTransactionCount(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pending != 1 {
		t.Fatalf("unexpected pending, wanted 1 got: %v", pending)
	}
	// Query balance
	balance, err := ec.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	penBalance, err := ec.PendingBalanceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Cmp(penBalance) == 0 {
		t.Fatalf("unexpected balance: %v %v", balance, penBalance)
	}
	// NonceAt
	nonce, err := ec.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	penNonce, err := ec.PendingNonceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if penNonce != nonce+1 {
		t.Fatalf("unexpected nonce: %v %v", nonce, penNonce)
	}
	// StorageAt
	storage, err := ec.StorageAt(context.Background(), testAddr, common.Hash{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	penStorage, err := ec.PendingStorageAt(context.Background(), testAddr, common.Hash{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(storage, penStorage) {
		t.Fatalf("unexpected storage: %v %v", storage, penStorage)
	}
	// CodeAt
	code, err := ec.CodeAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	penCode, err := ec.PendingCodeAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(code, penCode) {
		t.Fatalf("unexpected code: %v %v", code, penCode)
	}
}

func sendTransaction[P crypto.PublicKey](ec *Client[P]) error {
	// Retrieve chainID
	chainID, err := ec.ChainID(context.Background())
	if err != nil {
		return err
	}
	// Create transaction
	tx := types.NewTransaction[P](0, common.Address{1}, big.NewInt(1), 22000, big.NewInt(1), nil)
	signer := types.LatestSignerForChainID[P](chainID)
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), testKey)
	if err != nil {
		return err
	}
	signedTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return err
	}
	// Send transaction
	return ec.SendTransaction(context.Background(), signedTx, bind.PrivateTxArgs{PrivateFor: nil})
}

// Quorum

func TestClient_PreparePrivateTransaction_whenTypical(t *testing.T) {
	testObject := NewClient[nist.PublicKey](nil)

	_, err := testObject.PreparePrivateTransaction([]byte("arbitrary payload"), "arbitrary private from")

	assert.Error(t, err)
}

func TestClient_PreparePrivateTransaction_whenClientIsConfigured(t *testing.T) {
	expectedData := []byte("arbitrary payload")
	expectedDataEPH := common.BytesToEncryptedPayloadHash(expectedData)
	testObject := NewClient[nist.PublicKey](nil)
	testObject.pc = &privateTransactionManagerStubClient{expectedData}

	actualData, err := testObject.PreparePrivateTransaction([]byte("arbitrary payload"), "arbitrary private from")

	assert.NoError(t, err)
	assert.Equal(t, expectedDataEPH, actualData)
}

type privateTransactionManagerStubClient struct {
	expectedData []byte
}

func (s *privateTransactionManagerStubClient) StoreRaw(data []byte, from string) (common.EncryptedPayloadHash, error) {
	return common.BytesToEncryptedPayloadHash(data), nil
}
