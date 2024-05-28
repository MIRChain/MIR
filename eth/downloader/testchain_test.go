// Copyright 2018 The go-ethereum Authors
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

package downloader

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/consensus/ethash"
	"github.com/MIRChain/MIR/core"
	"github.com/MIRChain/MIR/core/rawdb"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/nist"
	"github.com/MIRChain/MIR/params"
)

// Test chain parameters.
var (
	testKey, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress[nist.PublicKey](*testKey.Public())
	testDB      = rawdb.NewMemoryDatabase()
	testGenesis = core.GenesisBlockForTesting[nist.PublicKey](testDB, testAddress, big.NewInt(1000000000))
)

// The common prefix of all test chains:
var testChainBase = newTestChain[nist.PrivateKey, nist.PublicKey](blockCacheMaxItems+200, testGenesis)

// Different forks on top of the base chain:
var testChainForkLightA, testChainForkLightB, testChainForkHeavy *testChain[nist.PrivateKey, nist.PublicKey]

func init() {
	var forkLen = int(fullMaxForkAncestry + 50)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() { testChainForkLightA = testChainBase.makeFork(forkLen, false, 1); wg.Done() }()
	go func() { testChainForkLightB = testChainBase.makeFork(forkLen, false, 2); wg.Done() }()
	go func() { testChainForkHeavy = testChainBase.makeFork(forkLen, true, 3); wg.Done() }()
	wg.Wait()
}

type testChain[T crypto.PrivateKey, P crypto.PublicKey] struct {
	genesis  *types.Block[P]
	chain    []common.Hash
	headerm  map[common.Hash]*types.Header[P]
	blockm   map[common.Hash]*types.Block[P]
	receiptm map[common.Hash][]*types.Receipt[P]
	tdm      map[common.Hash]*big.Int
}

// newTestChain creates a blockchain of the given length.
func newTestChain[T crypto.PrivateKey, P crypto.PublicKey](length int, genesis *types.Block[P]) *testChain[T, P] {
	tc := new(testChain[T, P]).copy(length)
	tc.genesis = genesis
	tc.chain = append(tc.chain, genesis.Hash())
	tc.headerm[tc.genesis.Hash()] = tc.genesis.Header()
	tc.tdm[tc.genesis.Hash()] = tc.genesis.Difficulty()
	tc.blockm[tc.genesis.Hash()] = tc.genesis
	tc.generate(length-1, 0, genesis, false)
	return tc
}

// makeFork creates a fork on top of the test chain.
func (tc *testChain[T, P]) makeFork(length int, heavy bool, seed byte) *testChain[T, P] {
	fork := tc.copy(tc.len() + length)
	fork.generate(length, seed, tc.headBlock(), heavy)
	return fork
}

// shorten creates a copy of the chain with the given length. It panics if the
// length is longer than the number of available blocks.
func (tc *testChain[T, P]) shorten(length int) *testChain[T, P] {
	if length > tc.len() {
		panic(fmt.Errorf("can't shorten test chain to %d blocks, it's only %d blocks long", length, tc.len()))
	}
	return tc.copy(length)
}

func (tc *testChain[T, P]) copy(newlen int) *testChain[T, P] {
	cpy := &testChain[T, P]{
		genesis:  tc.genesis,
		headerm:  make(map[common.Hash]*types.Header[P], newlen),
		blockm:   make(map[common.Hash]*types.Block[P], newlen),
		receiptm: make(map[common.Hash][]*types.Receipt[P], newlen),
		tdm:      make(map[common.Hash]*big.Int, newlen),
	}
	for i := 0; i < len(tc.chain) && i < newlen; i++ {
		hash := tc.chain[i]
		cpy.chain = append(cpy.chain, tc.chain[i])
		cpy.tdm[hash] = tc.tdm[hash]
		cpy.blockm[hash] = tc.blockm[hash]
		cpy.headerm[hash] = tc.headerm[hash]
		cpy.receiptm[hash] = tc.receiptm[hash]
	}
	return cpy
}

// generate creates a chain of n blocks starting at and including parent.
// the returned hash chain is ordered head->parent. In addition, every 22th block
// contains a transaction and every 5th an uncle to allow testing correct block
// reassembly.
func (tc *testChain[T, P]) generate(n int, seed byte, parent *types.Block[P], heavy bool) {
	// start := time.Now()
	// defer func() { fmt.Printf("test chain generated in %v\n", time.Since(start)) }()

	blocks, receipts := core.GenerateChain[P](params.TestChainConfig, parent, ethash.NewFaker[P](), testDB, n, func(i int, block *core.BlockGen[P]) {
		block.SetCoinbase(common.Address{seed})
		// If a heavy chain is requested, delay blocks to raise difficulty
		if heavy {
			block.OffsetTime(-1)
		}
		// Include transactions to the miner to make blocks more interesting.
		if parent == tc.genesis && i%22 == 0 {
			signer := types.MakeSigner[P](params.TestChainConfig, block.Number())
			var key T
			switch t := any(&testKey).(type) {
			case *nist.PrivateKey:
				tt := any(&key).(*nist.PrivateKey)
				*tt = *t
			case *gost3410.PrivateKey:
				tt := any(&key).(*gost3410.PrivateKey)
				*tt = *t
			}
			tx, err := types.SignTx[T, P](types.NewTransaction[P](block.TxNonce(testAddress), common.Address{seed}, big.NewInt(1000), params.TxGas, nil, nil), signer, key)
			if err != nil {
				panic(err)
			}
			block.AddTx(tx)
		}
		// if the block number is a multiple of 5, add a bonus uncle to the block
		if i > 0 && i%5 == 0 {
			block.AddUncle(&types.Header[P]{
				ParentHash: block.PrevBlock(i - 1).Hash(),
				Number:     big.NewInt(block.Number().Int64() - 1),
			})
		}
	})

	// Convert the block-chain into a hash-chain and header/block maps
	td := new(big.Int).Set(tc.td(parent.Hash()))
	for i, b := range blocks {
		td := td.Add(td, b.Difficulty())
		hash := b.Hash()
		tc.chain = append(tc.chain, hash)
		tc.blockm[hash] = b
		tc.headerm[hash] = b.Header()
		tc.receiptm[hash] = receipts[i]
		tc.tdm[hash] = new(big.Int).Set(td)
	}
}

// len returns the total number of blocks in the chain.
func (tc *testChain[T, P]) len() int {
	return len(tc.chain)
}

// headBlock returns the head of the chain.
func (tc *testChain[T, P]) headBlock() *types.Block[P] {
	return tc.blockm[tc.chain[len(tc.chain)-1]]
}

// td returns the total difficulty of the given block.
func (tc *testChain[T, P]) td(hash common.Hash) *big.Int {
	return tc.tdm[hash]
}

// headersByHash returns headers in order from the given hash.
func (tc *testChain[T, P]) headersByHash(origin common.Hash, amount int, skip int, reverse bool) []*types.Header[P] {
	num, _ := tc.hashToNumber(origin)
	return tc.headersByNumber(num, amount, skip, reverse)
}

// headersByNumber returns headers from the given number.
func (tc *testChain[T, P]) headersByNumber(origin uint64, amount int, skip int, reverse bool) []*types.Header[P] {
	result := make([]*types.Header[P], 0, amount)

	if !reverse {
		for num := origin; num < uint64(len(tc.chain)) && len(result) < amount; num += uint64(skip) + 1 {
			if header, ok := tc.headerm[tc.chain[int(num)]]; ok {
				result = append(result, header)
			}
		}
	} else {
		for num := int64(origin); num >= 0 && len(result) < amount; num -= int64(skip) + 1 {
			if header, ok := tc.headerm[tc.chain[int(num)]]; ok {
				result = append(result, header)
			}
		}
	}
	return result
}

// receipts returns the receipts of the given block hashes.
func (tc *testChain[T, P]) receipts(hashes []common.Hash) [][]*types.Receipt[P] {
	results := make([][]*types.Receipt[P], 0, len(hashes))
	for _, hash := range hashes {
		if receipt, ok := tc.receiptm[hash]; ok {
			results = append(results, receipt)
		}
	}
	return results
}

// bodies returns the block bodies of the given block hashes.
func (tc *testChain[T, P]) bodies(hashes []common.Hash) ([][]*types.Transaction[P], [][]*types.Header[P]) {
	transactions := make([][]*types.Transaction[P], 0, len(hashes))
	uncles := make([][]*types.Header[P], 0, len(hashes))
	for _, hash := range hashes {
		if block, ok := tc.blockm[hash]; ok {
			transactions = append(transactions, block.Transactions())
			uncles = append(uncles, block.Uncles())
		}
	}
	return transactions, uncles
}

func (tc *testChain[T, P]) hashToNumber(target common.Hash) (uint64, bool) {
	for num, hash := range tc.chain {
		if hash == target {
			return uint64(num), true
		}
	}
	return 0, false
}
