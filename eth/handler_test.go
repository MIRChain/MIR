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
	"math/big"
	"sort"
	"sync"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/consensus/ethash"
	"github.com/MIRChain/MIR/core"
	"github.com/MIRChain/MIR/core/rawdb"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/core/vm"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/nist"
	"github.com/MIRChain/MIR/eth/downloader"
	"github.com/MIRChain/MIR/ethdb"
	"github.com/MIRChain/MIR/event"
	"github.com/MIRChain/MIR/params"
)

var (
	// testKey is a private key to use for funding a tester account.
	testKey, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

	// testAddr is the Ethereum address of the tester account.
	testAddr = crypto.PubkeyToAddress[nist.PublicKey](*testKey.Public())
)

// testTxPool is a mock transaction pool that blindly accepts all transactions.
// Its goal is to get around setting up a valid statedb for the balance and nonce
// checks.
type testTxPool[P crypto.PublicKey] struct {
	pool map[common.Hash]*types.Transaction[P] // Hash map of collected transactions

	txFeed event.Feed   // Notification feed to allow waiting for inclusion
	lock   sync.RWMutex // Protects the transaction pool
}

// newTestTxPool creates a mock transaction pool.
func newTestTxPool[P crypto.PublicKey]() *testTxPool[P] {
	return &testTxPool[P]{
		pool: make(map[common.Hash]*types.Transaction[P]),
	}
}

// Has returns an indicator whether txpool has a transaction
// cached with the given hash.
func (p *testTxPool[P]) Has(hash common.Hash) bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.pool[hash] != nil
}

// Get retrieves the transaction from local txpool with given
// tx hash.
func (p *testTxPool[P]) Get(hash common.Hash) *types.Transaction[P] {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.pool[hash]
}

// AddRemotes appends a batch of transactions to the pool, and notifies any
// listeners if the addition channel is non nil
func (p *testTxPool[P]) AddRemotes(txs []*types.Transaction[P]) []error {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, tx := range txs {
		p.pool[tx.Hash()] = tx
	}
	p.txFeed.Send(core.NewTxsEvent[P]{Txs: txs})
	return make([]error, len(txs))
}

// Pending returns all the transactions known to the pool
func (p *testTxPool[P]) Pending() (map[common.Address]types.Transactions[P], error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	batches := make(map[common.Address]types.Transactions[P])
	for _, tx := range p.pool {
		from, _ := types.Sender[P](types.HomesteadSigner[P]{}, tx)
		batches[from] = append(batches[from], tx)
	}
	for _, batch := range batches {
		sort.Sort(types.TxByNonce[P](batch))
	}
	return batches, nil
}

// SubscribeNewTxsEvent should return an event subscription of NewTxsEvent and
// send events to the given channel.
func (p *testTxPool[P]) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent[P]) event.Subscription {
	return p.txFeed.Subscribe(ch)
}

// testHandler is a live implementation of the Ethereum protocol handler, just
// preinitialized with some sane testing defaults and the transaction pool mocked
// out.
type testHandler[T crypto.PrivateKey, P crypto.PublicKey] struct {
	db      ethdb.Database
	chain   *core.BlockChain[P]
	txpool  *testTxPool[P]
	handler *handler[T, P]
}

// newTestHandler creates a new handler for testing purposes with no blocks.
func newTestHandler[T crypto.PrivateKey, P crypto.PublicKey]() *testHandler[T, P] {
	return newTestHandlerWithBlocks[T, P](0)
}

// newTestHandlerWithBlocks creates a new handler for testing purposes, with a
// given number of initial blocks.
func newTestHandlerWithBlocks[T crypto.PrivateKey, P crypto.PublicKey](blocks int) *testHandler[T, P] {
	// Create a database pre-initialize with a genesis block
	db := rawdb.NewMemoryDatabase()
	(&core.Genesis[P]{
		Config: params.TestChainConfig,
		Alloc:  core.GenesisAlloc{testAddr: {Balance: big.NewInt(1000000)}},
	}).MustCommit(db)

	chain, _ := core.NewBlockChain[P](db, nil, params.TestChainConfig, ethash.NewFaker[P](), vm.Config[P]{}, nil, nil, nil)

	bs, _ := core.GenerateChain[P](params.TestChainConfig, chain.Genesis(), ethash.NewFaker[P](), db, blocks, nil)
	if _, err := chain.InsertChain(bs); err != nil {
		panic(err)
	}
	txpool := newTestTxPool[P]()

	handler, _ := newHandler(&handlerConfig[T, P]{
		Database:   db,
		Chain:      chain,
		TxPool:     txpool,
		Network:    1,
		Sync:       downloader.FastSync,
		BloomCache: 1,
	})
	handler.Start(1000)

	return &testHandler[T, P]{
		db:      db,
		chain:   chain,
		txpool:  txpool,
		handler: handler,
	}
}

// close tears down the handler and all its internal constructs.
func (b *testHandler[T, P]) close() {
	b.handler.Stop()
	b.chain.Stop()
}
