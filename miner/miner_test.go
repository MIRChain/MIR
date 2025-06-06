// Copyright 2020 The go-ethereum Authors
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

// Package miner implements Ethereum block creation and mining.
package miner

import (
	"testing"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/clique"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/ethdb/memorydb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

type mockBackend [P crypto.PublicKey] struct {
	bc     *core.BlockChain[P]
	txPool *core.TxPool[P]
	db     ethdb.Database
}

var _ Backend[nist.PrivateKey,nist.PublicKey] = &mockBackend[nist.PublicKey]{} // check implementation

func NewMockBackend[P crypto.PublicKey] (bc *core.BlockChain[P], txPool *core.TxPool[P]) *mockBackend[P] {
	return &mockBackend[P]{
		bc:     bc,
		txPool: txPool,
	}
}

func (m *mockBackend[P]) BlockChain() *core.BlockChain[P] {
	return m.bc
}

func (m *mockBackend[P]) TxPool() *core.TxPool[P] {
	return m.txPool
}

func (m *mockBackend[P]) ChainDb() ethdb.Database {
	return m.db
}

type testBlockChain [P crypto.PublicKey] struct {
	statedb *state.StateDB[P]

	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain[P]) CurrentBlock() *types.Block[P] {
	return types.NewBlock[P](&types.Header[P]{
		GasLimit: bc.gasLimit,
	}, nil, nil, nil, trie.NewStackTrie[P](nil))
}

func (bc *testBlockChain[P]) GetBlock(hash common.Hash, number uint64) *types.Block[P] {
	return bc.CurrentBlock()
}

func (bc *testBlockChain[P]) StateAt(common.Hash) (*state.StateDB[P], mps.PrivateStateRepository[P], error) {
	return bc.statedb, nil, nil
}

func (bc *testBlockChain[P]) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent[P]) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

func TestMiner(t *testing.T) {
	miner, mux := createMiner[nist.PrivateKey,nist.PublicKey](t)
	miner.Start(common.HexToAddress("0x12345"))
	waitForMiningState(t, miner, true)
	// Start the downloader
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)
	// Stop the downloader and wait for the update loop to run
	mux.Post(downloader.DoneEvent[nist.PublicKey]{})
	waitForMiningState(t, miner, true)

	// Subsequent downloader events after a successful DoneEvent should not cause the
	// miner to start or stop. This prevents a security vulnerability
	// that would allow entities to present fake high blocks that would
	// stop mining operations by causing a downloader sync
	// until it was discovered they were invalid, whereon mining would resume.
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, true)

	mux.Post(downloader.FailedEvent{})
	waitForMiningState(t, miner, true)
}

// TestMinerDownloaderFirstFails tests that mining is only
// permitted to run indefinitely once the downloader sees a DoneEvent (success).
// An initial FailedEvent should allow mining to stop on a subsequent
// downloader StartEvent.
func TestMinerDownloaderFirstFails(t *testing.T) {
	miner, mux := createMiner[nist.PrivateKey,nist.PublicKey](t)
	miner.Start(common.HexToAddress("0x12345"))
	waitForMiningState(t, miner, true)
	// Start the downloader
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)

	// Stop the downloader and wait for the update loop to run
	mux.Post(downloader.FailedEvent{})
	waitForMiningState(t, miner, true)

	// Since the downloader hasn't yet emitted a successful DoneEvent,
	// we expect the miner to stop on next StartEvent.
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)

	// Downloader finally succeeds.
	mux.Post(downloader.DoneEvent[nist.PublicKey]{})
	waitForMiningState(t, miner, true)

	// Downloader starts again.
	// Since it has achieved a DoneEvent once, we expect miner
	// state to be unchanged.
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, true)

	mux.Post(downloader.FailedEvent{})
	waitForMiningState(t, miner, true)
}

func TestMinerStartStopAfterDownloaderEvents(t *testing.T) {
	miner, mux := createMiner[nist.PrivateKey,nist.PublicKey](t)

	miner.Start(common.HexToAddress("0x12345"))
	waitForMiningState(t, miner, true)
	// Start the downloader
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)

	// Downloader finally succeeds.
	mux.Post(downloader.DoneEvent[nist.PublicKey]{})
	waitForMiningState(t, miner, true)

	miner.Stop()
	waitForMiningState(t, miner, false)

	miner.Start(common.HexToAddress("0x678910"))
	waitForMiningState(t, miner, true)

	miner.Stop()
	waitForMiningState(t, miner, false)
}

func TestStartWhileDownload(t *testing.T) {
	miner, mux := createMiner[nist.PrivateKey,nist.PublicKey](t)
	waitForMiningState(t, miner, false)
	miner.Start(common.HexToAddress("0x12345"))
	waitForMiningState(t, miner, true)
	// Stop the downloader and wait for the update loop to run
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)
	// Starting the miner after the downloader should not work
	miner.Start(common.HexToAddress("0x12345"))
	waitForMiningState(t, miner, false)
}

func TestStartStopMiner(t *testing.T) {
	miner, _ := createMiner[nist.PrivateKey,nist.PublicKey](t)
	waitForMiningState(t, miner, false)
	miner.Start(common.HexToAddress("0x12345"))
	waitForMiningState(t, miner, true)
	miner.Stop()
	waitForMiningState(t, miner, false)
}

func TestCloseMiner(t *testing.T) {
	miner, _ := createMiner[nist.PrivateKey,nist.PublicKey](t)
	waitForMiningState(t, miner, false)
	miner.Start(common.HexToAddress("0x12345"))
	waitForMiningState(t, miner, true)
	// Terminate the miner and wait for the update loop to run
	miner.Close()
	waitForMiningState(t, miner, false)
}

// TestMinerSetEtherbase checks that etherbase becomes set even if mining isn't
// possible at the moment
func TestMinerSetEtherbase(t *testing.T) {
	miner, mux := createMiner[nist.PrivateKey,nist.PublicKey](t)
	// Start with a 'bad' mining address
	miner.Start(common.HexToAddress("0xdead"))
	waitForMiningState(t, miner, true)
	// Start the downloader
	mux.Post(downloader.StartEvent{})
	waitForMiningState(t, miner, false)
	// Now user tries to configure proper mining address
	miner.Start(common.HexToAddress("0x1337"))
	// Stop the downloader and wait for the update loop to run
	mux.Post(downloader.DoneEvent[nist.PublicKey]{})

	waitForMiningState(t, miner, true)
	// The miner should now be using the good address
	if got, exp := miner.coinbase, common.HexToAddress("0x1337"); got != exp {
		t.Fatalf("Wrong coinbase, got %x expected %x", got, exp)
	}
}

// waitForMiningState waits until either
// * the desired mining state was reached
// * a timeout was reached which fails the test
func waitForMiningState[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T, m *Miner[T,P], mining bool) {
	t.Helper()

	var state bool
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		if state = m.Mining(); state == mining {
			return
		}
	}
	t.Fatalf("Mining() == %t, want %t", state, mining)
}

func createMiner[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T) (*Miner[T,P], *event.TypeMux) {
	// Create Ethash config
	config := Config{
		Etherbase: common.HexToAddress("123456789"),
	}
	// Create chainConfig
	memdb := memorydb.New()
	chainDB := rawdb.NewDatabase(memdb)
	genesis := core.DeveloperGenesisBlock[P](15, common.HexToAddress("12345"))
	chainConfig, _, err := core.SetupGenesisBlock(chainDB, genesis)
	if err != nil {
		t.Fatalf("can't create new chain config: %v", err)
	}
	// Create consensus engine
	engine := clique.New[P](chainConfig.Clique, chainDB)
	// Create Ethereum backend
	bc, err := core.NewBlockChain[P](chainDB, nil, chainConfig, engine, vm.Config[P]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("can't create new chain %v", err)
	}
	statedb, _ := state.New[P](common.Hash{}, state.NewDatabase[P](chainDB), nil)
	blockchain := &testBlockChain[P]{statedb, 10000000, new(event.Feed)}

	pool := core.NewTxPool[P](testTxPoolConfig, chainConfig, blockchain)
	backend := NewMockBackend[P](bc, pool)
	// Create event Mux
	mux := new(event.TypeMux)
	// Create Miner
	return New[T,P](backend, &config, chainConfig, mux, engine, nil), mux
}
