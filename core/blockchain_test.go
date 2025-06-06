// Copyright 2014 The go-ethereum Authors
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

package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"github.com/stretchr/testify/assert"
)

// So we can deterministically seed different blockchains
var (
	canonicalSeed = 1
	forkSeed      = 2
)

// newCanonical creates a chain database, and injects a deterministic canonical
// chain. Depending on the full flag, if creates either a full block chain or a
// header only chain.
func newCanonical[P crypto.PublicKey](engine consensus.Engine[P], n int, full bool) (ethdb.Database, *BlockChain[P], error) {
	var (
		db      = rawdb.NewMemoryDatabase()
		genesis = new(Genesis[P]).MustCommit(db)
	)

	// Initialize a fresh chain with only a genesis block
	blockchain, _ := NewBlockChain[P](db, nil, params.AllEthashProtocolChanges, engine, vm.Config[P]{}, nil, nil, nil)
	// Create and inject the requested chain
	if n == 0 {
		return db, blockchain, nil
	}
	if full {
		// Full block-chain requested
		blocks := makeBlockChain[P](genesis, n, engine, db, canonicalSeed)
		_, err := blockchain.InsertChain(blocks)
		return db, blockchain, err
	}
	// Header-only chain requested
	headers := makeHeaderChain[P](genesis.Header(), n, engine, db, canonicalSeed)
	_, err := blockchain.InsertHeaderChain(headers, 1)
	return db, blockchain, err
}

// Test fork of length N starting from block i
func testFork[P crypto.PublicKey](t *testing.T, blockchain *BlockChain[P], i, n int, full bool, comparator func(td1, td2 *big.Int)) {
	// Copy old chain up to #i into a new db
	db, blockchain2, err := newCanonical[P]( ethash.NewFaker[P](), i, full)
	if err != nil {
		t.Fatal("could not make new canonical in testFork", err)
	}
	defer blockchain2.Stop()

	// Assert the chains have the same header/block at #i
	var hash1, hash2 common.Hash
	if full {
		hash1 = blockchain.GetBlockByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetBlockByNumber(uint64(i)).Hash()
	} else {
		hash1 = blockchain.GetHeaderByNumber(uint64(i)).Hash()
		hash2 = blockchain2.GetHeaderByNumber(uint64(i)).Hash()
	}
	if hash1 != hash2 {
		t.Errorf("chain content mismatch at %d: have hash %v, want hash %v", i, hash2, hash1)
	}
	// Extend the newly created chain
	var (
		blockChainB  []*types.Block[P]
		headerChainB []*types.Header[P]
	)
	if full {
		blockChainB = makeBlockChain[P](blockchain2.CurrentBlock(), n,  ethash.NewFaker[P](), db, forkSeed)
		if _, err := blockchain2.InsertChain(blockChainB); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	} else {
		headerChainB = makeHeaderChain[P](blockchain2.CurrentHeader(), n,  ethash.NewFaker[P](), db, forkSeed)
		if _, err := blockchain2.InsertHeaderChain(headerChainB, 1); err != nil {
			t.Fatalf("failed to insert forking chain: %v", err)
		}
	}
	// Sanity check that the forked chain can be imported into the original
	var tdPre, tdPost *big.Int

	if full {
		tdPre = blockchain.GetTdByHash(blockchain.CurrentBlock().Hash())
		if err := testBlockChainImport(blockChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked block chain: %v", err)
		}
		tdPost = blockchain.GetTdByHash(blockChainB[len(blockChainB)-1].Hash())
	} else {
		tdPre = blockchain.GetTdByHash(blockchain.CurrentHeader().Hash())
		if err := testHeaderChainImport(headerChainB, blockchain); err != nil {
			t.Fatalf("failed to import forked header chain: %v", err)
		}
		tdPost = blockchain.GetTdByHash(headerChainB[len(headerChainB)-1].Hash())
	}
	// Compare the total difficulties of the chains
	comparator(tdPre, tdPost)
}

// testBlockChainImport tries to process a chain of blocks, writing them into
// the database if successful.
func testBlockChainImport[P crypto.PublicKey](chain types.Blocks[P], blockchain *BlockChain[P]) error {
	for _, block := range chain {
		// Try and process the block
		err := blockchain.engine.VerifyHeader(blockchain, block.Header(), true)
		if err == nil {
			err = blockchain.validator.ValidateBody(block)
		}
		if err != nil {
			if err == ErrKnownBlock {
				continue
			}
			return err
		}
		statedb, err := state.New[P](blockchain.GetBlockByHash(block.ParentHash()).Root(), blockchain.stateCache, nil)
		if err != nil {
			return err
		}
		privateStateRepo, repoErr := blockchain.PrivateStateManager().StateRepository(block.ParentHash())
		if repoErr != nil {
			return repoErr
		}
		receipts, _, _, usedGas, err := blockchain.processor.Process(block, statedb, privateStateRepo, vm.Config[P]{})
		if err != nil {
			blockchain.reportBlock(block, receipts, err)
			return err
		}
		err = blockchain.validator.ValidateState(block, statedb, receipts, usedGas)
		if err != nil {
			blockchain.reportBlock(block, receipts, err)
			return err
		}
		blockchain.chainmu.Lock()
		rawdb.WriteTd(blockchain.db, block.Hash(), block.NumberU64(), new(big.Int).Add(block.Difficulty(), blockchain.GetTdByHash(block.ParentHash())))
		rawdb.WriteBlock(blockchain.db, block)
		statedb.Commit(false)
		blockchain.chainmu.Unlock()
	}
	return nil
}

// testHeaderChainImport tries to process a chain of header, writing them into
// the database if successful.
func testHeaderChainImport[P crypto.PublicKey](chain []*types.Header[P], blockchain *BlockChain[P]) error {
	for _, header := range chain {
		// Try and validate the header
		if err := blockchain.engine.VerifyHeader(blockchain, header, false); err != nil {
			return err
		}
		// Manually insert the header into the database, but don't reorganise (allows subsequent testing)
		blockchain.chainmu.Lock()
		rawdb.WriteTd(blockchain.db, header.Hash(), header.Number.Uint64(), new(big.Int).Add(header.Difficulty, blockchain.GetTdByHash(header.ParentHash)))
		rawdb.WriteHeader(blockchain.db, header)
		blockchain.chainmu.Unlock()
	}
	return nil
}

func TestLastBlock(t *testing.T) {
	_, blockchain, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), 0, true)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	blocks := makeBlockChain[nist.PublicKey](blockchain.CurrentBlock(), 1, ethash.NewFullFaker[nist.PublicKey](), blockchain.db, 0)
	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatalf("Failed to insert block: %v", err)
	}
	if blocks[len(blocks)-1].Hash() != rawdb.ReadHeadBlockHash(blockchain.db) {
		t.Fatalf("Write/Get HeadBlockHash failed")
	}
}

// Tests that given a starting canonical chain of a given size, it can be extended
// with various length chains.
func TestExtendCanonicalHeaders(t *testing.T) { testExtendCanonical(t, false) }
func TestExtendCanonicalBlocks(t *testing.T)  { testExtendCanonical(t, true) }

func testExtendCanonical(t *testing.T, full bool) {
	length := 5

	// Make first chain starting from genesis
	_, processor, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Start fork from current height
	testFork(t, processor, length, 1, full, better)
	testFork(t, processor, length, 2, full, better)
	testFork(t, processor, length, 5, full, better)
	testFork(t, processor, length, 10, full, better)
}

// Tests that given a starting canonical chain of a given size, creating shorter
// forks do not take canonical ownership.
func TestShorterForkHeaders(t *testing.T) { testShorterFork(t, false) }
func TestShorterForkBlocks(t *testing.T)  { testShorterFork(t, true) }

func testShorterFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	worse := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) >= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected less than %v", td2, td1)
		}
	}
	// Sum of numbers must be less than `length` for this to be a shorter fork
	testFork(t, processor, 0, 3, full, worse)
	testFork(t, processor, 0, 7, full, worse)
	testFork(t, processor, 1, 1, full, worse)
	testFork(t, processor, 1, 7, full, worse)
	testFork(t, processor, 5, 3, full, worse)
	testFork(t, processor, 5, 4, full, worse)
}

// Tests that given a starting canonical chain of a given size, creating longer
// forks do take canonical ownership.
func TestLongerForkHeaders(t *testing.T) { testLongerFork(t, false) }
func TestLongerForkBlocks(t *testing.T)  { testLongerFork(t, true) }

func testLongerFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	better := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) <= 0 {
			t.Errorf("total difficulty mismatch: have %v, expected more than %v", td2, td1)
		}
	}
	// Sum of numbers must be greater than `length` for this to be a longer fork
	testFork(t, processor, 0, 11, full, better)
	testFork(t, processor, 0, 15, full, better)
	testFork(t, processor, 1, 10, full, better)
	testFork(t, processor, 1, 12, full, better)
	testFork(t, processor, 5, 6, full, better)
	testFork(t, processor, 5, 8, full, better)
}

// Tests that given a starting canonical chain of a given size, creating equal
// forks do take canonical ownership.
func TestEqualForkHeaders(t *testing.T) { testEqualFork(t, false) }
func TestEqualForkBlocks(t *testing.T)  { testEqualFork(t, true) }

func testEqualFork(t *testing.T, full bool) {
	length := 10

	// Make first chain starting from genesis
	_, processor, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), length, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer processor.Stop()

	// Define the difficulty comparator
	equal := func(td1, td2 *big.Int) {
		if td2.Cmp(td1) != 0 {
			t.Errorf("total difficulty mismatch: have %v, want %v", td2, td1)
		}
	}
	// Sum of numbers must be equal to `length` for this to be an equal fork
	testFork(t, processor, 0, 10, full, equal)
	testFork(t, processor, 1, 9, full, equal)
	testFork(t, processor, 2, 8, full, equal)
	testFork(t, processor, 5, 5, full, equal)
	testFork(t, processor, 6, 4, full, equal)
	testFork(t, processor, 9, 1, full, equal)
}

// Tests that chains missing links do not get accepted by the processor.
func TestBrokenHeaderChain(t *testing.T) { testBrokenChain(t, false) }
func TestBrokenBlockChain(t *testing.T)  { testBrokenChain(t, true) }

func testBrokenChain(t *testing.T, full bool) {
	// Make chain starting from genesis
	db, blockchain, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), 10, full)
	if err != nil {
		t.Fatalf("failed to make new canonical chain: %v", err)
	}
	defer blockchain.Stop()

	// Create a forked chain, and try to insert with a missing link
	if full {
		chain := makeBlockChain[nist.PublicKey](blockchain.CurrentBlock(), 5,  ethash.NewFaker[nist.PublicKey](), db, forkSeed)[1:]
		if err := testBlockChainImport(chain, blockchain); err == nil {
			t.Errorf("broken block chain not reported")
		}
	} else {
		chain := makeHeaderChain[nist.PublicKey](blockchain.CurrentHeader(), 5,  ethash.NewFaker[nist.PublicKey](), db, forkSeed)[1:]
		if err := testHeaderChainImport(chain, blockchain); err == nil {
			t.Errorf("broken header chain not reported")
		}
	}
}

// Tests that reorganising a long difficult chain after a short easy one
// overwrites the canonical numbers and links in the database.
func TestReorgLongHeaders(t *testing.T) { testReorgLong[nist.PublicKey](t, false) }
func TestReorgLongBlocks(t *testing.T)  { testReorgLong[nist.PublicKey](t, true) }

func testReorgLong[P crypto.PublicKey](t *testing.T, full bool) {
	testReorg[P](t, []int64{0, 0, -9}, []int64{0, 0, 0, -9}, 393280, full)
}

// Tests that reorganising a short difficult chain after a long easy one
// overwrites the canonical numbers and links in the database.
func TestReorgShortHeaders(t *testing.T) { testReorgShort[nist.PublicKey](t, false) }
func TestReorgShortBlocks(t *testing.T)  { testReorgShort[nist.PublicKey](t, true) }

func testReorgShort[P crypto.PublicKey](t *testing.T, full bool) {
	// Create a long easy chain vs. a short heavy one. Due to difficulty adjustment
	// we need a fairly long chain of blocks with different difficulties for a short
	// one to become heavyer than a long one. The 96 is an empirical value.
	easy := make([]int64, 96)
	for i := 0; i < len(easy); i++ {
		easy[i] = 60
	}
	diff := make([]int64, len(easy)-1)
	for i := 0; i < len(diff); i++ {
		diff[i] = -9
	}
	testReorg[P](t, easy, diff, 12615120, full)
}

func testReorg[P crypto.PublicKey](t *testing.T, first, second []int64, td int64, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical[P]( ethash.NewFaker[P](), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	// Insert an easy and a difficult chain afterwards
	easyBlocks, _ := GenerateChain[P](params.TestChainConfig, blockchain.CurrentBlock(),  ethash.NewFaker[P](), db, len(first), func(i int, b *BlockGen[P]) {
		b.OffsetTime(first[i])
	})
	diffBlocks, _ := GenerateChain[P](params.TestChainConfig, blockchain.CurrentBlock(),  ethash.NewFaker[P](), db, len(second), func(i int, b *BlockGen[P]) {
		b.OffsetTime(second[i])
	})
	if full {
		if _, err := blockchain.InsertChain(easyBlocks); err != nil {
			t.Fatalf("failed to insert easy chain: %v", err)
		}
		if _, err := blockchain.InsertChain(diffBlocks); err != nil {
			t.Fatalf("failed to insert difficult chain: %v", err)
		}
	} else {
		easyHeaders := make([]*types.Header[P], len(easyBlocks))
		for i, block := range easyBlocks {
			easyHeaders[i] = block.Header()
		}
		diffHeaders := make([]*types.Header[P], len(diffBlocks))
		for i, block := range diffBlocks {
			diffHeaders[i] = block.Header()
		}
		if _, err := blockchain.InsertHeaderChain(easyHeaders, 1); err != nil {
			t.Fatalf("failed to insert easy chain: %v", err)
		}
		if _, err := blockchain.InsertHeaderChain(diffHeaders, 1); err != nil {
			t.Fatalf("failed to insert difficult chain: %v", err)
		}
	}
	// Check that the chain is valid number and link wise
	if full {
		prev := blockchain.CurrentBlock()
		for block := blockchain.GetBlockByNumber(blockchain.CurrentBlock().NumberU64() - 1); block.NumberU64() != 0; prev, block = block, blockchain.GetBlockByNumber(block.NumberU64()-1) {
			if prev.ParentHash() != block.Hash() {
				t.Errorf("parent block hash mismatch: have %x, want %x", prev.ParentHash(), block.Hash())
			}
		}
	} else {
		prev := blockchain.CurrentHeader()
		for header := blockchain.GetHeaderByNumber(blockchain.CurrentHeader().Number.Uint64() - 1); header.Number.Uint64() != 0; prev, header = header, blockchain.GetHeaderByNumber(header.Number.Uint64()-1) {
			if prev.ParentHash != header.Hash() {
				t.Errorf("parent header hash mismatch: have %x, want %x", prev.ParentHash, header.Hash())
			}
		}
	}
	// Make sure the chain total difficulty is the correct one
	want := new(big.Int).Add(blockchain.genesisBlock.Difficulty(), big.NewInt(td))
	if full {
		if have := blockchain.GetTdByHash(blockchain.CurrentBlock().Hash()); have.Cmp(want) != 0 {
			t.Errorf("total difficulty mismatch: have %v, want %v", have, want)
		}
	} else {
		if have := blockchain.GetTdByHash(blockchain.CurrentHeader().Hash()); have.Cmp(want) != 0 {
			t.Errorf("total difficulty mismatch: have %v, want %v", have, want)
		}
	}
}

// Tests that the insertion functions detect banned hashes.
func TestBadHeaderHashes(t *testing.T) { testBadHashes(t, false) }
func TestBadBlockHashes(t *testing.T)  { testBadHashes(t, true) }

func testBadHashes(t *testing.T, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	// Create a chain, ban a hash and try to import
	if full {
		blocks := makeBlockChain[nist.PublicKey](blockchain.CurrentBlock(), 3,  ethash.NewFaker[nist.PublicKey](), db, 10)

		BadHashes[blocks[2].Header().Hash()] = true
		defer func() { delete(BadHashes, blocks[2].Header().Hash()) }()

		_, err = blockchain.InsertChain(blocks)
	} else {
		headers := makeHeaderChain[nist.PublicKey](blockchain.CurrentHeader(), 3,  ethash.NewFaker[nist.PublicKey](), db, 10)

		BadHashes[headers[2].Hash()] = true
		defer func() { delete(BadHashes, headers[2].Hash()) }()

		_, err = blockchain.InsertHeaderChain(headers, 1)
	}
	if !errors.Is(err, ErrBlacklistedHash) {
		t.Errorf("error mismatch: have: %v, want: %v", err, ErrBlacklistedHash)
	}
}

// Tests that bad hashes are detected on boot, and the chain rolled back to a
// good state prior to the bad hash.
func TestReorgBadHeaderHashes(t *testing.T) { testReorgBadHashes(t, false) }
func TestReorgBadBlockHashes(t *testing.T)  { testReorgBadHashes(t, true) }

func testReorgBadHashes(t *testing.T, full bool) {
	// Create a pristine chain and database
	db, blockchain, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), 0, full)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	// Create a chain, import and ban afterwards
	headers := makeHeaderChain[nist.PublicKey](blockchain.CurrentHeader(), 4,  ethash.NewFaker[nist.PublicKey](), db, 10)
	blocks := makeBlockChain[nist.PublicKey](blockchain.CurrentBlock(), 4,  ethash.NewFaker[nist.PublicKey](), db, 10)

	if full {
		if _, err = blockchain.InsertChain(blocks); err != nil {
			t.Errorf("failed to import blocks: %v", err)
		}
		if blockchain.CurrentBlock().Hash() != blocks[3].Hash() {
			t.Errorf("last block hash mismatch: have: %x, want %x", blockchain.CurrentBlock().Hash(), blocks[3].Header().Hash())
		}
		BadHashes[blocks[3].Header().Hash()] = true
		defer func() { delete(BadHashes, blocks[3].Header().Hash()) }()
	} else {
		if _, err = blockchain.InsertHeaderChain(headers, 1); err != nil {
			t.Errorf("failed to import headers: %v", err)
		}
		if blockchain.CurrentHeader().Hash() != headers[3].Hash() {
			t.Errorf("last header hash mismatch: have: %x, want %x", blockchain.CurrentHeader().Hash(), headers[3].Hash())
		}
		BadHashes[headers[3].Hash()] = true
		defer func() { delete(BadHashes, headers[3].Hash()) }()
	}
	blockchain.Stop()

	// Create a new BlockChain and check that it rolled back the state.
	ncm, err := NewBlockChain[nist.PublicKey](blockchain.db, nil, blockchain.chainConfig,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create new chain manager: %v", err)
	}
	if full {
		if ncm.CurrentBlock().Hash() != blocks[2].Header().Hash() {
			t.Errorf("last block hash mismatch: have: %x, want %x", ncm.CurrentBlock().Hash(), blocks[2].Header().Hash())
		}
		if blocks[2].Header().GasLimit != ncm.GasLimit() {
			t.Errorf("last  block gasLimit mismatch: have: %d, want %d", ncm.GasLimit(), blocks[2].Header().GasLimit)
		}
	} else {
		if ncm.CurrentHeader().Hash() != headers[2].Hash() {
			t.Errorf("last header hash mismatch: have: %x, want %x", ncm.CurrentHeader().Hash(), headers[2].Hash())
		}
	}
	ncm.Stop()
}

// Tests chain insertions in the face of one entity containing an invalid nonce.
func TestHeadersInsertNonceError(t *testing.T) { testInsertNonceError(t, false) }
func TestBlocksInsertNonceError(t *testing.T)  { testInsertNonceError(t, true) }

func testInsertNonceError(t *testing.T, full bool) {
	for i := 1; i < 25 && !t.Failed(); i++ {
		// Create a pristine chain and database
		db, blockchain, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), 0, full)
		if err != nil {
			t.Fatalf("failed to create pristine chain: %v", err)
		}
		defer blockchain.Stop()

		// Create and insert a chain with a failing nonce
		var (
			failAt  int
			failRes int
			failNum uint64
		)
		if full {
			blocks := makeBlockChain[nist.PublicKey](blockchain.CurrentBlock(), i,  ethash.NewFaker[nist.PublicKey](), db, 0)

			failAt = rand.Int() % len(blocks)
			failNum = blocks[failAt].NumberU64()

			blockchain.engine = ethash.NewFakeFailer[nist.PublicKey](failNum)
			failRes, err = blockchain.InsertChain(blocks)
		} else {
			headers := makeHeaderChain[nist.PublicKey](blockchain.CurrentHeader(), i,  ethash.NewFaker[nist.PublicKey](), db, 0)

			failAt = rand.Int() % len(headers)
			failNum = headers[failAt].Number.Uint64()

			blockchain.engine = ethash.NewFakeFailer[nist.PublicKey](failNum)
			blockchain.hc.engine = blockchain.engine
			failRes, err = blockchain.InsertHeaderChain(headers, 1)
		}
		// Check that the returned error indicates the failure
		if failRes != failAt {
			t.Errorf("test %d: failure (%v) index mismatch: have %d, want %d", i, err, failRes, failAt)
		}
		// Check that all blocks after the failing block have been inserted
		for j := 0; j < i-failAt; j++ {
			if full {
				if block := blockchain.GetBlockByNumber(failNum + uint64(j)); block != nil {
					t.Errorf("test %d: invalid block in chain: %v", i, block)
				}
			} else {
				if header := blockchain.GetHeaderByNumber(failNum + uint64(j)); header != nil {
					t.Errorf("test %d: invalid header in chain: %v", i, header)
				}
			}
		}
	}
}

// Tests that fast importing a block chain produces the same chain data as the
// classical full block processing.
func TestFastVsFullChains(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis[nist.PublicKey]{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{address: {Balance: funds}},
		}
		genesis = gspec.MustCommit(gendb)
		signer  = types.LatestSigner[nist.PublicKey](gspec.Config)
	)
	blocks, receipts := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), gendb, 1024, func(i int, block *BlockGen[nist.PublicKey]) {
		block.SetCoinbase(common.Address{0x00})

		// If the block number is multiple of 3, send a few bonus transactions to the miner
		if i%3 == 2 {
			for j := 0; j < i%4+1; j++ {
				tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), params.TxGas, nil, nil), signer, key)
				if err != nil {
					panic(err)
				}
				block.AddTx(tx)
			}
		}
		// If the block number is a multiple of 5, add a few bonus uncles to the block
		if i%5 == 5 {
			block.AddUncle(&types.Header[nist.PublicKey]{ParentHash: block.PrevBlock(i - 1).Hash(), Number: big.NewInt(int64(i - 1))})
		}
	})
	// Import the chain as an archive node for the comparison baseline
	archiveDb := rawdb.NewMemoryDatabase()
	gspec.MustCommit(archiveDb)
	archive, _ := NewBlockChain[nist.PublicKey](archiveDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer archive.Stop()

	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	// Fast import the chain as a non-archive node to test
	fastDb := rawdb.NewMemoryDatabase()
	gspec.MustCommit(fastDb)
	fast, _ := NewBlockChain[nist.PublicKey](fastDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer fast.Stop()

	headers := make([]*types.Header[nist.PublicKey], len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts, 0); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	// Freezer style fast import the chain.
	frdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp freezer dir: %v", err)
	}
	defer os.Remove(frdir)
	ancientDb, err := rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), frdir, "", false)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	gspec.MustCommit(ancientDb)
	ancient, _ := NewBlockChain[nist.PublicKey](ancientDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer ancient.Stop()

	if n, err := ancient.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(len(blocks)/2)); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	// Iterate over all chain data components, and cross reference
	for i := 0; i < len(blocks); i++ {
		num, hash := blocks[i].NumberU64(), blocks[i].Hash()

		if ftd, atd := fast.GetTdByHash(hash), archive.GetTdByHash(hash); ftd.Cmp(atd) != 0 {
			t.Errorf("block #%d [%x]: td mismatch: fastdb %v, archivedb %v", num, hash, ftd, atd)
		}
		if antd, artd := ancient.GetTdByHash(hash), archive.GetTdByHash(hash); antd.Cmp(artd) != 0 {
			t.Errorf("block #%d [%x]: td mismatch: ancientdb %v, archivedb %v", num, hash, antd, artd)
		}
		if fheader, aheader := fast.GetHeaderByHash(hash), archive.GetHeaderByHash(hash); fheader.Hash() != aheader.Hash() {
			t.Errorf("block #%d [%x]: header mismatch: fastdb %v, archivedb %v", num, hash, fheader, aheader)
		}
		if anheader, arheader := ancient.GetHeaderByHash(hash), archive.GetHeaderByHash(hash); anheader.Hash() != arheader.Hash() {
			t.Errorf("block #%d [%x]: header mismatch: ancientdb %v, archivedb %v", num, hash, anheader, arheader)
		}
		if fblock, arblock, anblock := fast.GetBlockByHash(hash), archive.GetBlockByHash(hash), ancient.GetBlockByHash(hash); fblock.Hash() != arblock.Hash() || anblock.Hash() != arblock.Hash() {
			t.Errorf("block #%d [%x]: block mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, fblock, anblock, arblock)
		} else if types.DeriveSha(fblock.Transactions(), trie.NewStackTrie[nist.PublicKey](nil)) != types.DeriveSha(arblock.Transactions(), trie.NewStackTrie[nist.PublicKey](nil)) || types.DeriveSha(anblock.Transactions(), trie.NewStackTrie[nist.PublicKey](nil)) != types.DeriveSha(arblock.Transactions(), trie.NewStackTrie[nist.PublicKey](nil)) {
			t.Errorf("block #%d [%x]: transactions mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, fblock.Transactions(), anblock.Transactions(), arblock.Transactions())
		} else if types.CalcUncleHash(fblock.Uncles()) != types.CalcUncleHash(arblock.Uncles()) || types.CalcUncleHash(anblock.Uncles()) != types.CalcUncleHash(arblock.Uncles()) {
			t.Errorf("block #%d [%x]: uncles mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, fblock.Uncles(), anblock, arblock.Uncles())
		}
		if freceipts, anreceipts, areceipts := rawdb.ReadReceipts[nist.PublicKey](fastDb, hash, *rawdb.ReadHeaderNumber(fastDb, hash), fast.Config()), rawdb.ReadReceipts[nist.PublicKey](ancientDb, hash, *rawdb.ReadHeaderNumber(ancientDb, hash), fast.Config()), rawdb.ReadReceipts[nist.PublicKey](archiveDb, hash, *rawdb.ReadHeaderNumber(archiveDb, hash), fast.Config()); types.DeriveSha(freceipts, trie.NewStackTrie[nist.PublicKey](nil)) != types.DeriveSha(areceipts, trie.NewStackTrie[nist.PublicKey](nil)) {
			t.Errorf("block #%d [%x]: receipts mismatch: fastdb %v, ancientdb %v, archivedb %v", num, hash, freceipts, anreceipts, areceipts)
		}
	}
	// Check that the canonical chains are the same between the databases
	for i := 0; i < len(blocks)+1; i++ {
		if fhash, ahash := rawdb.ReadCanonicalHash(fastDb, uint64(i)), rawdb.ReadCanonicalHash(archiveDb, uint64(i)); fhash != ahash {
			t.Errorf("block #%d: canonical hash mismatch: fastdb %v, archivedb %v", i, fhash, ahash)
		}
		if anhash, arhash := rawdb.ReadCanonicalHash(ancientDb, uint64(i)), rawdb.ReadCanonicalHash(archiveDb, uint64(i)); anhash != arhash {
			t.Errorf("block #%d: canonical hash mismatch: ancientdb %v, archivedb %v", i, anhash, arhash)
		}
	}
}

// Tests that various import methods move the chain head pointers to the correct
// positions.
func TestLightVsFastVsFullChainHeads(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{address: {Balance: funds}}}
		genesis = gspec.MustCommit(gendb)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), gendb, int(height), nil)

	// makeDb creates a db instance for testing.
	makeDb := func() (ethdb.Database, func()) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Fatalf("failed to create temp freezer dir: %v", err)
		}
		defer os.Remove(dir)
		db, err := rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), dir, "", false)
		if err != nil {
			t.Fatalf("failed to create temp freezer db: %v", err)
		}
		gspec.MustCommit(db)
		return db, func() { os.RemoveAll(dir) }
	}
	// Configure a subchain to roll back
	remove := blocks[height/2].NumberU64()

	// Create a small assertion method to check the three heads
	assert := func(t *testing.T, kind string, chain *BlockChain[nist.PublicKey], header uint64, fast uint64, block uint64) {
		t.Helper()

		if num := chain.CurrentBlock().NumberU64(); num != block {
			t.Errorf("%s head block mismatch: have #%v, want #%v", kind, num, block)
		}
		if num := chain.CurrentFastBlock().NumberU64(); num != fast {
			t.Errorf("%s head fast-block mismatch: have #%v, want #%v", kind, num, fast)
		}
		if num := chain.CurrentHeader().Number.Uint64(); num != header {
			t.Errorf("%s head header mismatch: have #%v, want #%v", kind, num, header)
		}
	}
	// Import the chain as an archive node and ensure all pointers are updated
	archiveDb, delfn := makeDb()
	defer delfn()

	archiveCaching := *defaultCacheConfig
	archiveCaching.TrieDirtyDisabled = true

	archive, _ := NewBlockChain[nist.PublicKey](archiveDb, &archiveCaching, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if n, err := archive.InsertChain(blocks); err != nil {
		t.Fatalf("failed to process block %d: %v", n, err)
	}
	defer archive.Stop()

	assert(t, "archive", archive, height, height, height)
	archive.SetHead(remove - 1)
	assert(t, "archive", archive, height/2, height/2, height/2)

	// Import the chain as a non-archive node and ensure all pointers are updated
	fastDb, delfn := makeDb()
	defer delfn()
	fast, _ := NewBlockChain[nist.PublicKey](fastDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer fast.Stop()

	headers := make([]*types.Header[nist.PublicKey], len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := fast.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := fast.InsertReceiptChain(blocks, receipts, 0); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	assert(t, "fast", fast, height, height, 0)
	fast.SetHead(remove - 1)
	assert(t, "fast", fast, height/2, height/2, 0)

	// Import the chain as a ancient-first node and ensure all pointers are updated
	ancientDb, delfn := makeDb()
	defer delfn()
	ancient, _ := NewBlockChain[nist.PublicKey](ancientDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer ancient.Stop()

	if n, err := ancient.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(3*len(blocks)/4)); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	assert(t, "ancient", ancient, height, height, 0)
	ancient.SetHead(remove - 1)
	assert(t, "ancient", ancient, 0, 0, 0)

	if frozen, err := ancientDb.Ancients(); err != nil || frozen != 1 {
		t.Fatalf("failed to truncate ancient store, want %v, have %v", 1, frozen)
	}
	// Import the chain as a light node and ensure all pointers are updated
	lightDb, delfn := makeDb()
	defer delfn()
	light, _ := NewBlockChain[nist.PublicKey](lightDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if n, err := light.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	defer light.Stop()

	assert(t, "light", light, height, 0, 0)
	light.SetHead(remove - 1)
	assert(t, "light", light, height/2, 0, 0)
}

// Tests that chain reorganisations handle transaction removals and reinsertions.
func TestChainTxReorgs(t *testing.T) {
	var (
		key1, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA[nist.PrivateKey]("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key3, _ = crypto.HexToECDSA[nist.PrivateKey]("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addr1   = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())
		addr2   = crypto.PubkeyToAddress[nist.PublicKey](*key2.Public())
		addr3   = crypto.PubkeyToAddress[nist.PublicKey](*key3.Public())
		db      = rawdb.NewMemoryDatabase()
		gspec   = &Genesis[nist.PublicKey]{
			Config:   params.TestChainConfig,
			GasLimit: 700000000,
			Alloc: GenesisAlloc{
				addr1: {Balance: big.NewInt(1000000)},
				addr2: {Balance: big.NewInt(1000000)},
				addr3: {Balance: big.NewInt(1000000)},
			},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSigner[nist.PublicKey](gspec.Config)
	)

	// Create two transactions shared between the chains:
	//  - postponed: transaction included at a later block in the forked chain
	//  - swapped: transaction included at the same block number in the forked chain
	postponed, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](0, addr1, big.NewInt(1000), params.TxGas, nil, nil), signer, key1)
	swapped, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](1, addr1, big.NewInt(1000), params.TxGas, nil, nil), signer, key1)

	// Create two transactions that will be dropped by the forked chain:
	//  - pastDrop: transaction dropped retroactively from a past block
	//  - freshDrop: transaction dropped exactly at the block where the reorg is detected
	var pastDrop, freshDrop *types.Transaction[nist.PublicKey]

	// Create three transactions that will be added in the forked chain:
	//  - pastAdd:   transaction added before the reorganization is detected
	//  - freshAdd:  transaction added at the exact block the reorg is detected
	//  - futureAdd: transaction added after the reorg has already finished
	var pastAdd, freshAdd, futureAdd *types.Transaction[nist.PublicKey]

	chain, _ := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), db, 3, func(i int, gen *BlockGen[nist.PublicKey]) {
		switch i {
		case 0:
			pastDrop, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](gen.TxNonce(addr2), addr2, big.NewInt(1000), params.TxGas, nil, nil), signer, key2)

			gen.AddTx(pastDrop)  // This transaction will be dropped in the fork from below the split point
			gen.AddTx(postponed) // This transaction will be postponed till block #3 in the fork

		case 2:
			freshDrop, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](gen.TxNonce(addr2), addr2, big.NewInt(1000), params.TxGas, nil, nil), signer, key2)

			gen.AddTx(freshDrop) // This transaction will be dropped in the fork from exactly at the split point
			gen.AddTx(swapped)   // This transaction will be swapped out at the exact height

			gen.OffsetTime(9) // Lower the block difficulty to simulate a weaker chain
		}
	})
	// Import the chain. This runs all block validation rules.
	blockchain, _ := NewBlockChain[nist.PublicKey](db, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if i, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert original chain[%d]: %v", i, err)
	}
	defer blockchain.Stop()

	// overwrite the old chain
	chain, _ = GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), db, 5, func(i int, gen *BlockGen[nist.PublicKey]) {
		switch i {
		case 0:
			pastAdd, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](gen.TxNonce(addr3), addr3, big.NewInt(1000), params.TxGas, nil, nil), signer, key3)
			gen.AddTx(pastAdd) // This transaction needs to be injected during reorg

		case 2:
			gen.AddTx(postponed) // This transaction was postponed from block #1 in the original chain
			gen.AddTx(swapped)   // This transaction was swapped from the exact current spot in the original chain

			freshAdd, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](gen.TxNonce(addr3), addr3, big.NewInt(1000), params.TxGas, nil, nil), signer, key3)
			gen.AddTx(freshAdd) // This transaction will be added exactly at reorg time

		case 3:
			futureAdd, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](gen.TxNonce(addr3), addr3, big.NewInt(1000), params.TxGas, nil, nil), signer, key3)
			gen.AddTx(futureAdd) // This transaction will be added after a full reorg
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}

	// removed tx
	for i, tx := range (types.Transactions[nist.PublicKey]{pastDrop, freshDrop}) {
		if txn, _, _, _ := rawdb.ReadTransaction[nist.PublicKey](db, tx.Hash()); txn != nil {
			t.Errorf("drop %d: tx %v found while shouldn't have been", i, txn)
		}
		if rcpt, _, _, _ := rawdb.ReadReceipt[nist.PublicKey](db, tx.Hash(), blockchain.Config()); rcpt != nil {
			t.Errorf("drop %d: receipt %v found while shouldn't have been", i, rcpt)
		}
	}
	// added tx
	for i, tx := range (types.Transactions[nist.PublicKey]{pastAdd, freshAdd, futureAdd}) {
		if txn, _, _, _ := rawdb.ReadTransaction[nist.PublicKey](db, tx.Hash()); txn == nil {
			t.Errorf("add %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := rawdb.ReadReceipt[nist.PublicKey](db, tx.Hash(), blockchain.Config()); rcpt == nil {
			t.Errorf("add %d: expected receipt to be found", i)
		}
	}
	// shared tx
	for i, tx := range (types.Transactions[nist.PublicKey]{postponed, swapped}) {
		if txn, _, _, _ := rawdb.ReadTransaction[nist.PublicKey](db, tx.Hash()); txn == nil {
			t.Errorf("share %d: expected tx to be found", i)
		}
		if rcpt, _, _, _ := rawdb.ReadReceipt[nist.PublicKey](db, tx.Hash(), blockchain.Config()); rcpt == nil {
			t.Errorf("share %d: expected receipt to be found", i)
		}
	}
}

func TestLogReorgs(t *testing.T) {
	var (
		key1, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())
		db      = rawdb.NewMemoryDatabase()
		// this code generates a log
		code    = common.Hex2Bytes("60606040525b7f24ec1d3ff24c2f6ff210738839dbc339cd45a5294d85c79361016243157aae7b60405180905060405180910390a15b600a8060416000396000f360606040526008565b00")
		gspec   = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{addr1: {Balance: big.NewInt(10000000000000)}}}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSigner[nist.PublicKey](gspec.Config)
	)

	blockchain, _ := NewBlockChain[nist.PublicKey](db, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer blockchain.Stop()

	rmLogsCh := make(chan RemovedLogsEvent[nist.PublicKey])
	blockchain.SubscribeRemovedLogsEvent(rmLogsCh)
	chain, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis,  ethash.NewFaker[nist.PublicKey](), db, 2, func(i int, gen *BlockGen[nist.PublicKey]) {
		if i == 1 {
			tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewContractCreation[nist.PublicKey](gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), code), signer, key1)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	chain, _ = GenerateChain[nist.PublicKey](params.TestChainConfig, genesis,  ethash.NewFaker[nist.PublicKey](), db, 3, func(i int, gen *BlockGen[nist.PublicKey]) {})
	done := make(chan struct{})
	go func() {
		ev := <-rmLogsCh
		if len(ev.Logs) == 0 {
			t.Error("expected logs")
		}
		close(done)
	}()
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}
	timeout := time.NewTimer(1 * time.Second)
	defer timeout.Stop()
	select {
	case <-done:
	case <-timeout.C:
		t.Fatal("Timeout. There is no RemovedLogsEvent has been sent.")
	}
}

// This EVM code generates a log when the contract is created.
var logCode = common.Hex2Bytes("60606040525b7f24ec1d3ff24c2f6ff210738839dbc339cd45a5294d85c79361016243157aae7b60405180905060405180910390a15b600a8060416000396000f360606040526008565b00")

// This test checks that log events and RemovedLogsEvent are sent
// when the chain reorganizes.
func TestLogRebirth(t *testing.T) {
	var (
		key1, _       = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1         = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())
		db            = rawdb.NewMemoryDatabase()
		gspec         = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{addr1: {Balance: big.NewInt(10000000000000)}}}
		genesis       = gspec.MustCommit(db)
		signer        = types.LatestSigner[nist.PublicKey](gspec.Config)
		engine        =  ethash.NewFaker[nist.PublicKey]()
		blockchain, _ = NewBlockChain[nist.PublicKey](db, nil, gspec.Config, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	)

	defer blockchain.Stop()

	// The event channels.
	newLogCh := make(chan []*types.Log, 10)
	rmLogsCh := make(chan RemovedLogsEvent[nist.PublicKey], 10)
	blockchain.SubscribeLogsEvent(newLogCh)
	blockchain.SubscribeRemovedLogsEvent(rmLogsCh)

	// This chain contains a single log.
	chain, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 2, func(i int, gen *BlockGen[nist.PublicKey]) {
		if i == 1 {
			tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewContractCreation[nist.PublicKey](gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), logCode), signer, key1)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}
	checkLogEvents(t, newLogCh, rmLogsCh, 1, 0)

	// Generate long reorg chain containing another log. Inserting the
	// chain removes one log and adds one.
	forkChain, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 2, func(i int, gen *BlockGen[nist.PublicKey]) {
		if i == 1 {
			tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewContractCreation[nist.PublicKey](gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), logCode), signer, key1)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
			gen.OffsetTime(-9) // higher block difficulty
		}
	})
	if _, err := blockchain.InsertChain(forkChain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}
	checkLogEvents(t, newLogCh, rmLogsCh, 1, 1)

	// This chain segment is rooted in the original chain, but doesn't contain any logs.
	// When inserting it, the canonical chain switches away from forkChain and re-emits
	// the log event for the old chain, as well as a RemovedLogsEvent for forkChain.
	newBlocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, chain[len(chain)-1], engine, db, 1, func(i int, gen *BlockGen[nist.PublicKey]) {})
	if _, err := blockchain.InsertChain(newBlocks); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}
	checkLogEvents(t, newLogCh, rmLogsCh, 1, 1)
}

// This test is a variation of TestLogRebirth. It verifies that log events are emitted
// when a side chain containing log events overtakes the canonical chain.
func TestSideLogRebirth(t *testing.T) {
	var (
		key1, _       = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1         = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())
		db            = rawdb.NewMemoryDatabase()
		gspec         = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{addr1: {Balance: big.NewInt(10000000000000)}}}
		genesis       = gspec.MustCommit(db)
		signer        = types.LatestSigner[nist.PublicKey](gspec.Config)
		blockchain, _ = NewBlockChain[nist.PublicKey](db, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	)

	defer blockchain.Stop()

	newLogCh := make(chan []*types.Log, 10)
	rmLogsCh := make(chan RemovedLogsEvent[nist.PublicKey], 10)
	blockchain.SubscribeLogsEvent(newLogCh)
	blockchain.SubscribeRemovedLogsEvent(rmLogsCh)

	chain, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis,  ethash.NewFaker[nist.PublicKey](), db, 2, func(i int, gen *BlockGen[nist.PublicKey]) {
		if i == 1 {
			gen.OffsetTime(-9) // higher block difficulty

		}
	})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}
	checkLogEvents(t, newLogCh, rmLogsCh, 0, 0)

	// Generate side chain with lower difficulty
	sideChain, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis,  ethash.NewFaker[nist.PublicKey](), db, 2, func(i int, gen *BlockGen[nist.PublicKey]) {
		if i == 1 {
			tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewContractCreation[nist.PublicKey](gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), logCode), signer, key1)
			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			gen.AddTx(tx)
		}
	})
	if _, err := blockchain.InsertChain(sideChain); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}
	checkLogEvents(t, newLogCh, rmLogsCh, 0, 0)

	// Generate a new block based on side chain.
	newBlocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, sideChain[len(sideChain)-1],  ethash.NewFaker[nist.PublicKey](), db, 1, func(i int, gen *BlockGen[nist.PublicKey]) {})
	if _, err := blockchain.InsertChain(newBlocks); err != nil {
		t.Fatalf("failed to insert forked chain: %v", err)
	}
	checkLogEvents(t, newLogCh, rmLogsCh, 1, 0)
}

func checkLogEvents(t *testing.T, logsCh <-chan []*types.Log, rmLogsCh <-chan RemovedLogsEvent[nist.PublicKey], wantNew, wantRemoved int) {
	t.Helper()

	if len(logsCh) != wantNew {
		t.Fatalf("wrong number of log events: got %d, want %d", len(logsCh), wantNew)
	}
	if len(rmLogsCh) != wantRemoved {
		t.Fatalf("wrong number of removed log events: got %d, want %d", len(rmLogsCh), wantRemoved)
	}
	// Drain events.
	for i := 0; i < len(logsCh); i++ {
		<-logsCh
	}
	for i := 0; i < len(rmLogsCh); i++ {
		<-rmLogsCh
	}
}

func TestReorgSideEvent(t *testing.T) {
	var (
		db      = rawdb.NewMemoryDatabase()
		key1, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())
		gspec   = &Genesis[nist.PublicKey]{
			Config: params.TestChainConfig,
			Alloc:  GenesisAlloc{addr1: {Balance: big.NewInt(10000000000000)}},
		}
		genesis = gspec.MustCommit(db)
		signer  = types.LatestSigner[nist.PublicKey](gspec.Config)
	)

	blockchain, _ := NewBlockChain[nist.PublicKey](db, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer blockchain.Stop()

	chain, _ := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), db, 3, func(i int, gen *BlockGen[nist.PublicKey]) {})
	if _, err := blockchain.InsertChain(chain); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	replacementBlocks, _ := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), db, 4, func(i int, gen *BlockGen[nist.PublicKey]) {
		tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewContractCreation[nist.PublicKey](gen.TxNonce(addr1), new(big.Int), 1000000, new(big.Int), nil), signer, key1)
		if i == 2 {
			gen.OffsetTime(-9)
		}
		if err != nil {
			t.Fatalf("failed to create tx: %v", err)
		}
		gen.AddTx(tx)
	})
	chainSideCh := make(chan ChainSideEvent[nist.PublicKey], 64)
	blockchain.SubscribeChainSideEvent(chainSideCh)
	if _, err := blockchain.InsertChain(replacementBlocks); err != nil {
		t.Fatalf("failed to insert chain: %v", err)
	}

	// first two block of the secondary chain are for a brief moment considered
	// side chains because up to that point the first one is considered the
	// heavier chain.
	expectedSideHashes := map[common.Hash]bool{
		replacementBlocks[0].Hash(): true,
		replacementBlocks[1].Hash(): true,
		chain[0].Hash():             true,
		chain[1].Hash():             true,
		chain[2].Hash():             true,
	}

	i := 0

	const timeoutDura = 10 * time.Second
	timeout := time.NewTimer(timeoutDura)
done:
	for {
		select {
		case ev := <-chainSideCh:
			block := ev.Block
			if _, ok := expectedSideHashes[block.Hash()]; !ok {
				t.Errorf("%d: didn't expect %x to be in side chain", i, block.Hash())
			}
			i++

			if i == len(expectedSideHashes) {
				timeout.Stop()

				break done
			}
			timeout.Reset(timeoutDura)

		case <-timeout.C:
			t.Fatal("Timeout. Possibly not all blocks were triggered for sideevent")
		}
	}

	// make sure no more events are fired
	select {
	case e := <-chainSideCh:
		t.Errorf("unexpected event fired: %v", e)
	case <-time.After(250 * time.Millisecond):
	}

}

// Tests if the canonical block can be fetched from the database during chain insertion.
func TestCanonicalBlockRetrieval(t *testing.T) {
	_, blockchain, err := newCanonical[nist.PublicKey]( ethash.NewFaker[nist.PublicKey](), 0, true)
	if err != nil {
		t.Fatalf("failed to create pristine chain: %v", err)
	}
	defer blockchain.Stop()

	chain, _ := GenerateChain[nist.PublicKey](blockchain.chainConfig, blockchain.genesisBlock,  ethash.NewFaker[nist.PublicKey](), blockchain.db, 10, func(i int, gen *BlockGen[nist.PublicKey]) {})

	var pend sync.WaitGroup
	pend.Add(len(chain))

	for i := range chain {
		go func(block *types.Block[nist.PublicKey]) {
			defer pend.Done()

			// try to retrieve a block by its canonical hash and see if the block data can be retrieved.
			for {
				ch := rawdb.ReadCanonicalHash(blockchain.db, block.NumberU64())
				if ch == (common.Hash{}) {
					continue // busy wait for canonical hash to be written
				}
				if ch != block.Hash() {
					t.Errorf("unknown canonical hash, want %s, got %s", block.Hash().Hex(), ch.Hex())
					return
				}
				fb := rawdb.ReadBlock[nist.PublicKey](blockchain.db, ch, block.NumberU64())
				if fb == nil {
					t.Errorf("unable to retrieve block %d for canonical hash: %s", block.NumberU64(), ch.Hex())
					return
				}
				if fb.Hash() != block.Hash() {
					t.Errorf("invalid block hash for block %d, want %s, got %s", block.NumberU64(), block.Hash().Hex(), fb.Hash().Hex())
					return
				}
				return
			}
		}(chain[i])

		if _, err := blockchain.InsertChain(types.Blocks[nist.PublicKey]{chain[i]}); err != nil {
			t.Fatalf("failed to insert block %d: %v", i, err)
		}
	}
	pend.Wait()
}

func TestEIP155Transition(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		db         = rawdb.NewMemoryDatabase()
		key, _     = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address    = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds      = big.NewInt(1000000000)
		deleteAddr = common.Address{1}
		gspec      = &Genesis[nist.PublicKey]{
			Config: &params.ChainConfig{ChainID: big.NewInt(10), EIP150Block: big.NewInt(0), EIP155Block: big.NewInt(2), HomesteadBlock: new(big.Int)},
			Alloc:  GenesisAlloc{address: {Balance: funds}, deleteAddr: {Balance: new(big.Int)}},
		}
		genesis = gspec.MustCommit(db)
	)

	blockchain, _ := NewBlockChain[nist.PublicKey](db, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer blockchain.Stop()

	blocks, _ := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), db, 4, func(i int, block *BlockGen[nist.PublicKey]) {
		var (
			tx      *types.Transaction[nist.PublicKey]
			err     error
			basicTx = func(signer types.Signer[nist.PublicKey]) (*types.Transaction[nist.PublicKey], error) {
				return types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), common.Address{}, new(big.Int), 21000, new(big.Int), nil), signer, key)
			}
		)
		switch i {
		case 0:
			tx, err = basicTx(types.HomesteadSigner[nist.PublicKey]{})
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		case 2:
			tx, err = basicTx(types.HomesteadSigner[nist.PublicKey]{})
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)

			tx, err = basicTx(types.LatestSigner[nist.PublicKey](gspec.Config))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		case 3:
			tx, err = basicTx(types.HomesteadSigner[nist.PublicKey]{})
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)

			tx, err = basicTx(types.LatestSigner[nist.PublicKey](gspec.Config))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		}
	})

	if _, err := blockchain.InsertChain(blocks); err != nil {
		t.Fatal(err)
	}
	block := blockchain.GetBlockByNumber(1)
	if block.Transactions()[0].Protected() {
		t.Error("Expected block[0].txs[0] to not be replay protected")
	}

	block = blockchain.GetBlockByNumber(3)
	if block.Transactions()[0].Protected() {
		t.Error("Expected block[3].txs[0] to not be replay protected")
	}
	if !block.Transactions()[1].Protected() {
		t.Error("Expected block[3].txs[1] to be replay protected")
	}
	if _, err := blockchain.InsertChain(blocks[4:]); err != nil {
		t.Fatal(err)
	}

	// generate an invalid chain id transaction
	config := &params.ChainConfig{ChainID: big.NewInt(2), EIP150Block: big.NewInt(0), EIP155Block: big.NewInt(2), HomesteadBlock: new(big.Int)}
	blocks, _ = GenerateChain[nist.PublicKey](config, blocks[len(blocks)-1],  ethash.NewFaker[nist.PublicKey](), db, 4, func(i int, block *BlockGen[nist.PublicKey]) {
		var (
			tx      *types.Transaction[nist.PublicKey]
			err     error
			basicTx = func(signer types.Signer[nist.PublicKey]) (*types.Transaction[nist.PublicKey], error) {
				return types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), common.Address{}, new(big.Int), 21000, new(big.Int), nil), signer, key)
			}
		)
		if i == 0 {
			tx, err = basicTx(types.LatestSigner[nist.PublicKey](config))
			if err != nil {
				t.Fatal(err)
			}
			block.AddTx(tx)
		}
	})
	_, err := blockchain.InsertChain(blocks)
	assert.ErrorIs(t, err, types.ErrInvalidChainId)
}

func TestEIP161AccountRemoval(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		db      = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		theAddr = common.Address{1}
		gspec   = &Genesis[nist.PublicKey]{
			Config: &params.ChainConfig{
				ChainID:        big.NewInt(10),
				HomesteadBlock: new(big.Int),
				EIP155Block:    new(big.Int),
				EIP150Block:    new(big.Int),
				EIP158Block:    big.NewInt(2),
			},
			Alloc: GenesisAlloc{address: {Balance: funds}},
		}
		genesis = gspec.MustCommit(db)
	)
	blockchain, _ := NewBlockChain[nist.PublicKey](db, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer blockchain.Stop()

	blocks, _ := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), db, 3, func(i int, block *BlockGen[nist.PublicKey]) {
		var (
			tx     *types.Transaction[nist.PublicKey]
			err    error
			signer = types.LatestSigner[nist.PublicKey](gspec.Config)
		)
		switch i {
		case 0:
			tx, err = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		case 1:
			tx, err = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		case 2:
			tx, err = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), theAddr, new(big.Int), 21000, new(big.Int), nil), signer, key)
		}
		if err != nil {
			t.Fatal(err)
		}
		block.AddTx(tx)
	})
	// account must exist pre eip 161
	if _, err := blockchain.InsertChain(types.Blocks[nist.PublicKey]{blocks[0]}); err != nil {
		t.Fatal(err)
	}
	if st, _, _ := blockchain.State(); !st.Exist(theAddr) {
		t.Error("expected account to exist")
	}

	// account needs to be deleted post eip 161
	if _, err := blockchain.InsertChain(types.Blocks[nist.PublicKey]{blocks[1]}); err != nil {
		t.Fatal(err)
	}
	if st, _, _ := blockchain.State(); st.Exist(theAddr) {
		t.Error("account should not exist")
	}

	// account mustn't be created post eip 161
	if _, err := blockchain.InsertChain(types.Blocks[nist.PublicKey]{blocks[2]}); err != nil {
		t.Fatal(err)
	}
	if st, _, _ := blockchain.State(); st.Exist(theAddr) {
		t.Error("account should not exist")
	}
}

// This is a regression test (i.e. as weird as it is, don't delete it ever), which
// tests that under weird reorg conditions the blockchain and its internal header-
// chain return the same latest block/header.
//
// https://github.com/pavelkrolevets/MIR-pro/pull/15941
func TestBlockchainHeaderchainReorgConsistency(t *testing.T) {
	// Generate a canonical chain to act as the main dataset
	engine :=  ethash.NewFaker[nist.PublicKey]()

	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)
	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 64, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{1}) })

	// Generate a bunch of fork blocks, each side forking from the canonical chain
	forks := make([]*types.Block[nist.PublicKey], len(blocks))
	for i := 0; i < len(forks); i++ {
		parent := genesis
		if i > 0 {
			parent = blocks[i-1]
		}
		fork, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, parent, engine, db, 1, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{2}) })
		forks[i] = fork[0]
	}
	// Import the canonical and fork chain side by side, verifying the current block
	// and current header consistency
	diskdb := rawdb.NewMemoryDatabase()
	new(Genesis[nist.PublicKey]).MustCommit(diskdb)

	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	for i := 0; i < len(blocks); i++ {
		if _, err := chain.InsertChain(blocks[i : i+1]); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", i, err)
		}
		if chain.CurrentBlock().Hash() != chain.CurrentHeader().Hash() {
			t.Errorf("block %d: current block/header mismatch: block #%d [%x..], header #%d [%x..]", i, chain.CurrentBlock().Number(), chain.CurrentBlock().Hash().Bytes()[:4], chain.CurrentHeader().Number, chain.CurrentHeader().Hash().Bytes()[:4])
		}
		if _, err := chain.InsertChain(forks[i : i+1]); err != nil {
			t.Fatalf(" fork %d: failed to insert into chain: %v", i, err)
		}
		if chain.CurrentBlock().Hash() != chain.CurrentHeader().Hash() {
			t.Errorf(" fork %d: current block/header mismatch: block #%d [%x..], header #%d [%x..]", i, chain.CurrentBlock().Number(), chain.CurrentBlock().Hash().Bytes()[:4], chain.CurrentHeader().Number, chain.CurrentHeader().Hash().Bytes()[:4])
		}
	}
}

// Tests that importing small side forks doesn't leave junk in the trie database
// cache (which would eventually cause memory issues).
func TestTrieForkGC(t *testing.T) {
	// Generate a canonical chain to act as the main dataset
	engine :=  ethash.NewFaker[nist.PublicKey]()

	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)
	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 2*TriesInMemory, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{1}) })

	// Generate a bunch of fork blocks, each side forking from the canonical chain
	forks := make([]*types.Block[nist.PublicKey], len(blocks))
	for i := 0; i < len(forks); i++ {
		parent := genesis
		if i > 0 {
			parent = blocks[i-1]
		}
		fork, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, parent, engine, db, 1, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{2}) })
		forks[i] = fork[0]
	}
	// Import the canonical and fork chain side by side, forcing the trie cache to cache both
	diskdb := rawdb.NewMemoryDatabase()
	new(Genesis[nist.PublicKey]).MustCommit(diskdb)

	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	for i := 0; i < len(blocks); i++ {
		if _, err := chain.InsertChain(blocks[i : i+1]); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", i, err)
		}
		if _, err := chain.InsertChain(forks[i : i+1]); err != nil {
			t.Fatalf("fork %d: failed to insert into chain: %v", i, err)
		}
	}
	// Dereference all the recent tries and ensure no past trie is left in
	for i := 0; i < TriesInMemory; i++ {
		chain.stateCache.TrieDB().Dereference(blocks[len(blocks)-1-i].Root())
		chain.stateCache.TrieDB().Dereference(forks[len(blocks)-1-i].Root())
	}
	if len(chain.stateCache.TrieDB().Nodes()) > 0 {
		t.Fatalf("stale tries still alive after garbase collection")
	}
}

// Tests that doing large reorgs works even if the state associated with the
// forking point is not available any more.
func TestLargeReorgTrieGC(t *testing.T) {
	// Generate the original common chain segment and the two competing forks
	engine :=  ethash.NewFaker[nist.PublicKey]()

	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)

	shared, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 64, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{1}) })
	original, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, shared[len(shared)-1], engine, db, 2*TriesInMemory, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{2}) })
	competitor, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, shared[len(shared)-1], engine, db, 2*TriesInMemory+1, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{3}) })

	// Import the shared chain and the original canonical one
	diskdb := rawdb.NewMemoryDatabase()
	new(Genesis[nist.PublicKey]).MustCommit(diskdb)

	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if _, err := chain.InsertChain(shared); err != nil {
		t.Fatalf("failed to insert shared chain: %v", err)
	}
	if _, err := chain.InsertChain(original); err != nil {
		t.Fatalf("failed to insert original chain: %v", err)
	}
	// Ensure that the state associated with the forking point is pruned away
	if node, _ := chain.stateCache.TrieDB().Node(shared[len(shared)-1].Root()); node != nil {
		t.Fatalf("common-but-old ancestor still cache")
	}
	// Import the competitor chain without exceeding the canonical's TD and ensure
	// we have not processed any of the blocks (protection against malicious blocks)
	if _, err := chain.InsertChain(competitor[:len(competitor)-2]); err != nil {
		t.Fatalf("failed to insert competitor chain: %v", err)
	}
	for i, block := range competitor[:len(competitor)-2] {
		if node, _ := chain.stateCache.TrieDB().Node(block.Root()); node != nil {
			t.Fatalf("competitor %d: low TD chain became processed", i)
		}
	}
	// Import the head of the competitor chain, triggering the reorg and ensure we
	// successfully reprocess all the stashed away blocks.
	if _, err := chain.InsertChain(competitor[len(competitor)-2:]); err != nil {
		t.Fatalf("failed to finalize competitor chain: %v", err)
	}
	for i, block := range competitor[:len(competitor)-TriesInMemory] {
		if node, _ := chain.stateCache.TrieDB().Node(block.Root()); node != nil {
			t.Fatalf("competitor %d: competing chain state missing", i)
		}
	}
}

func TestBlockchainRecovery(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{address: {Balance: funds}}}
		genesis = gspec.MustCommit(gendb)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), gendb, int(height), nil)

	// Import the chain as a ancient-first node and ensure all pointers are updated
	frdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp freezer dir: %v", err)
	}
	defer os.Remove(frdir)

	ancientDb, err := rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), frdir, "", false)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	gspec.MustCommit(ancientDb)
	ancient, _ := NewBlockChain[nist.PublicKey](ancientDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)

	headers := make([]*types.Header[nist.PublicKey], len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := ancient.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(3*len(blocks)/4)); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	rawdb.WriteLastPivotNumber(ancientDb, blocks[len(blocks)-1].NumberU64()) // Force fast sync behavior
	ancient.Stop()

	// Destroy head fast block manually
	midBlock := blocks[len(blocks)/2]
	rawdb.WriteHeadFastBlockHash(ancientDb, midBlock.Hash())

	// Reopen broken blockchain again
	ancient, _ = NewBlockChain[nist.PublicKey](ancientDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer ancient.Stop()
	if num := ancient.CurrentBlock().NumberU64(); num != 0 {
		t.Errorf("head block mismatch: have #%v, want #%v", num, 0)
	}
	if num := ancient.CurrentFastBlock().NumberU64(); num != midBlock.NumberU64() {
		t.Errorf("head fast-block mismatch: have #%v, want #%v", num, midBlock.NumberU64())
	}
	if num := ancient.CurrentHeader().Number.Uint64(); num != midBlock.NumberU64() {
		t.Errorf("head header mismatch: have #%v, want #%v", num, midBlock.NumberU64())
	}
}

func TestIncompleteAncientReceiptChainInsertion(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{address: {Balance: funds}}}
		genesis = gspec.MustCommit(gendb)
	)
	height := uint64(1024)
	blocks, receipts := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), gendb, int(height), nil)

	// Import the chain as a ancient-first node and ensure all pointers are updated
	frdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp freezer dir: %v", err)
	}
	defer os.Remove(frdir)
	ancientDb, err := rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), frdir, "", false)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	gspec.MustCommit(ancientDb)
	ancient, _ := NewBlockChain[nist.PublicKey](ancientDb, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer ancient.Stop()

	headers := make([]*types.Header[nist.PublicKey], len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := ancient.InsertHeaderChain(headers, 1); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	// Abort ancient receipt chain insertion deliberately
	ancient.terminateInsert = func(hash common.Hash, number uint64) bool {
		return number == blocks[len(blocks)/2].NumberU64()
	}
	previousFastBlock := ancient.CurrentFastBlock()
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(3*len(blocks)/4)); err == nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	if ancient.CurrentFastBlock().NumberU64() != previousFastBlock.NumberU64() {
		t.Fatalf("failed to rollback ancient data, want %d, have %d", previousFastBlock.NumberU64(), ancient.CurrentFastBlock().NumberU64())
	}
	if frozen, err := ancient.db.Ancients(); err != nil || frozen != 1 {
		t.Fatalf("failed to truncate ancient data")
	}
	ancient.terminateInsert = nil
	if n, err := ancient.InsertReceiptChain(blocks, receipts, uint64(3*len(blocks)/4)); err != nil {
		t.Fatalf("failed to insert receipt %d: %v", n, err)
	}
	if ancient.CurrentFastBlock().NumberU64() != blocks[len(blocks)-1].NumberU64() {
		t.Fatalf("failed to insert ancient recept chain after rollback")
	}
}

// Tests that importing a very large side fork, which is larger than the canon chain,
// but where the difficulty per block is kept low: this means that it will not
// overtake the 'canon' chain until after it's passed canon by about 200 blocks.
//
// Details at:
//  - https://github.com/pavelkrolevets/MIR-pro/issues/18977
//  - https://github.com/pavelkrolevets/MIR-pro/pull/18988
func TestLowDiffLongChain(t *testing.T) {
	// Generate a canonical chain to act as the main dataset
	engine :=  ethash.NewFaker[nist.PublicKey]()
	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)

	// We must use a pretty long chain to ensure that the fork doesn't overtake us
	// until after at least 128 blocks post tip
	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 6*TriesInMemory, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{1})
		b.OffsetTime(-9)
	})

	// Import the canonical chain
	diskdb := rawdb.NewMemoryDatabase()
	new(Genesis[nist.PublicKey]).MustCommit(diskdb)

	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	// Generate fork chain, starting from an early block
	parent := blocks[10]
	fork, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, parent, engine, db, 8*TriesInMemory, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{2})
	})

	// And now import the fork
	if i, err := chain.InsertChain(fork); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", i, err)
	}
	head := chain.CurrentBlock()
	if got := fork[len(fork)-1].Hash(); got != head.Hash() {
		t.Fatalf("head wrong, expected %x got %x", head.Hash(), got)
	}
	// Sanity check that all the canonical numbers are present
	header := chain.CurrentHeader()
	for number := head.NumberU64(); number > 0; number-- {
		if hash := chain.GetHeaderByNumber(number).Hash(); hash != header.Hash() {
			t.Fatalf("header %d: canonical hash mismatch: have %x, want %x", number, hash, header.Hash())
		}
		header = chain.GetHeader(header.ParentHash, number-1)
	}
}

// Tests that importing a sidechain (S), where
// - S is sidechain, containing blocks [Sn...Sm]
// - C is canon chain, containing blocks [G..Cn..Cm]
// - A common ancestor is placed at prune-point + blocksBetweenCommonAncestorAndPruneblock
// - The sidechain S is prepended with numCanonBlocksInSidechain blocks from the canon chain
func testSideImport(t *testing.T, numCanonBlocksInSidechain, blocksBetweenCommonAncestorAndPruneblock int) {

	// Generate a canonical chain to act as the main dataset
	engine :=  ethash.NewFaker[nist.PublicKey]()
	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)

	// Generate and import the canonical chain
	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 2*TriesInMemory, nil)
	diskdb := rawdb.NewMemoryDatabase()
	new(Genesis[nist.PublicKey]).MustCommit(diskdb)
	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	lastPrunedIndex := len(blocks) - TriesInMemory - 1
	lastPrunedBlock := blocks[lastPrunedIndex]
	firstNonPrunedBlock := blocks[len(blocks)-TriesInMemory]

	// Verify pruning of lastPrunedBlock
	if chain.HasBlockAndState(lastPrunedBlock.Hash(), lastPrunedBlock.NumberU64()) {
		t.Errorf("Block %d not pruned", lastPrunedBlock.NumberU64())
	}
	// Verify firstNonPrunedBlock is not pruned
	if !chain.HasBlockAndState(firstNonPrunedBlock.Hash(), firstNonPrunedBlock.NumberU64()) {
		t.Errorf("Block %d pruned", firstNonPrunedBlock.NumberU64())
	}
	// Generate the sidechain
	// First block should be a known block, block after should be a pruned block. So
	// canon(pruned), side, side...

	// Generate fork chain, make it longer than canon
	parentIndex := lastPrunedIndex + blocksBetweenCommonAncestorAndPruneblock
	parent := blocks[parentIndex]
	fork, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, parent, engine, db, 2*TriesInMemory, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{2})
	})
	// Prepend the parent(s)
	var sidechain []*types.Block[nist.PublicKey]
	for i := numCanonBlocksInSidechain; i > 0; i-- {
		sidechain = append(sidechain, blocks[parentIndex+1-i])
	}
	sidechain = append(sidechain, fork...)
	_, err = chain.InsertChain(sidechain)
	if err != nil {
		t.Errorf("Got error, %v", err)
	}
	head := chain.CurrentBlock()
	if got := fork[len(fork)-1].Hash(); got != head.Hash() {
		t.Fatalf("head wrong, expected %x got %x", head.Hash(), got)
	}
}

// Tests that importing a sidechain (S), where
// - S is sidechain, containing blocks [Sn...Sm]
// - C is canon chain, containing blocks [G..Cn..Cm]
// - The common ancestor Cc is pruned
// - The first block in S: Sn, is == Cn
// That is: the sidechain for import contains some blocks already present in canon chain.
// So the blocks are
// [ Cn, Cn+1, Cc, Sn+3 ... Sm]
//   ^    ^    ^  pruned
func TestPrunedImportSide(t *testing.T) {
	//glogger := log.NewGlogHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(false)))
	//glogger.Verbosity(3)
	//log.Root().SetHandler(log.Handler(glogger))
	testSideImport(t, 3, 3)
	testSideImport(t, 3, -3)
	testSideImport(t, 10, 0)
	testSideImport(t, 1, 10)
	testSideImport(t, 1, -10)
}

func TestInsertKnownHeaders(t *testing.T)      { testInsertKnownChainData(t, "headers") }
func TestInsertKnownReceiptChain(t *testing.T) { testInsertKnownChainData(t, "receipts") }
func TestInsertKnownBlocks(t *testing.T)       { testInsertKnownChainData(t, "blocks") }

func testInsertKnownChainData(t *testing.T, typ string) {
	engine :=  ethash.NewFaker[nist.PublicKey]()

	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)

	blocks, receipts := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 32, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{1}) })
	// A longer chain but total difficulty is lower.
	blocks2, receipts2 := GenerateChain[nist.PublicKey](params.TestChainConfig, blocks[len(blocks)-1], engine, db, 65, func(i int, b *BlockGen[nist.PublicKey]) { b.SetCoinbase(common.Address{1}) })
	// A shorter chain but total difficulty is higher.
	blocks3, receipts3 := GenerateChain[nist.PublicKey](params.TestChainConfig, blocks[len(blocks)-1], engine, db, 64, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{1})
		b.OffsetTime(-9) // A higher difficulty
	})
	// Import the shared chain and the original canonical one
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp freezer dir: %v", err)
	}
	defer os.Remove(dir)
	chaindb, err := rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), dir, "", false)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	new(Genesis[nist.PublicKey]).MustCommit(chaindb)
	defer os.RemoveAll(dir)

	chain, err := NewBlockChain[nist.PublicKey](chaindb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}

	var (
		inserter func(blocks []*types.Block[nist.PublicKey], receipts []types.Receipts[nist.PublicKey]) error
		asserter func(t *testing.T, block *types.Block[nist.PublicKey])
	)
	if typ == "headers" {
		inserter = func(blocks []*types.Block[nist.PublicKey], receipts []types.Receipts[nist.PublicKey]) error {
			headers := make([]*types.Header[nist.PublicKey], 0, len(blocks))
			for _, block := range blocks {
				headers = append(headers, block.Header())
			}
			_, err := chain.InsertHeaderChain(headers, 1)
			return err
		}
		asserter = func(t *testing.T, block *types.Block[nist.PublicKey]) {
			if chain.CurrentHeader().Hash() != block.Hash() {
				t.Fatalf("current head header mismatch, have %v, want %v", chain.CurrentHeader().Hash().Hex(), block.Hash().Hex())
			}
		}
	} else if typ == "receipts" {
		inserter = func(blocks []*types.Block[nist.PublicKey], receipts []types.Receipts[nist.PublicKey]) error {
			headers := make([]*types.Header[nist.PublicKey], 0, len(blocks))
			for _, block := range blocks {
				headers = append(headers, block.Header())
			}
			_, err := chain.InsertHeaderChain(headers, 1)
			if err != nil {
				return err
			}
			_, err = chain.InsertReceiptChain(blocks, receipts, 0)
			return err
		}
		asserter = func(t *testing.T, block *types.Block[nist.PublicKey]) {
			if chain.CurrentFastBlock().Hash() != block.Hash() {
				t.Fatalf("current head fast block mismatch, have %v, want %v", chain.CurrentFastBlock().Hash().Hex(), block.Hash().Hex())
			}
		}
	} else {
		inserter = func(blocks []*types.Block[nist.PublicKey], receipts []types.Receipts[nist.PublicKey]) error {
			_, err := chain.InsertChain(blocks)
			return err
		}
		asserter = func(t *testing.T, block *types.Block[nist.PublicKey]) {
			if chain.CurrentBlock().Hash() != block.Hash() {
				t.Fatalf("current head block mismatch, have %v, want %v", chain.CurrentBlock().Hash().Hex(), block.Hash().Hex())
			}
		}
	}

	if err := inserter(blocks, receipts); err != nil {
		t.Fatalf("failed to insert chain data: %v", err)
	}

	// Reimport the chain data again. All the imported
	// chain data are regarded "known" data.
	if err := inserter(blocks, receipts); err != nil {
		t.Fatalf("failed to insert chain data: %v", err)
	}
	asserter(t, blocks[len(blocks)-1])

	// Import a long canonical chain with some known data as prefix.
	rollback := blocks[len(blocks)/2].NumberU64()

	chain.SetHead(rollback - 1)
	if err := inserter(append(blocks, blocks2...), append(receipts, receipts2...)); err != nil {
		t.Fatalf("failed to insert chain data: %v", err)
	}
	asserter(t, blocks2[len(blocks2)-1])

	// Import a heavier shorter but higher total difficulty chain with some known data as prefix.
	if err := inserter(append(blocks, blocks3...), append(receipts, receipts3...)); err != nil {
		t.Fatalf("failed to insert chain data: %v", err)
	}
	asserter(t, blocks3[len(blocks3)-1])

	// Import a longer but lower total difficulty chain with some known data as prefix.
	if err := inserter(append(blocks, blocks2...), append(receipts, receipts2...)); err != nil {
		t.Fatalf("failed to insert chain data: %v", err)
	}
	// The head shouldn't change.
	asserter(t, blocks3[len(blocks3)-1])

	// Rollback the heavier chain and re-insert the longer chain again
	chain.SetHead(rollback - 1)
	if err := inserter(append(blocks, blocks2...), append(receipts, receipts2...)); err != nil {
		t.Fatalf("failed to insert chain data: %v", err)
	}
	asserter(t, blocks2[len(blocks2)-1])
}

// getLongAndShortChains returns two chains,
// A is longer, B is heavier
func getLongAndShortChains() (*BlockChain[nist.PublicKey], []*types.Block[nist.PublicKey], []*types.Block[nist.PublicKey], error) {
	// Generate a canonical chain to act as the main dataset
	engine :=  ethash.NewFaker[nist.PublicKey]()
	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)

	// Generate and import the canonical chain,
	// Offset the time, to keep the difficulty low
	longChain, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 80, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{1})
	})
	diskdb := rawdb.NewMemoryDatabase()
	new(Genesis[nist.PublicKey]).MustCommit(diskdb)

	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create tester chain: %v", err)
	}

	// Generate fork chain, make it shorter than canon, with common ancestor pretty early
	parentIndex := 3
	parent := longChain[parentIndex]
	heavyChain, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, parent, engine, db, 75, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{2})
		b.OffsetTime(-9)
	})
	// Verify that the test is sane
	var (
		longerTd  = new(big.Int)
		shorterTd = new(big.Int)
	)
	for index, b := range longChain {
		longerTd.Add(longerTd, b.Difficulty())
		if index <= parentIndex {
			shorterTd.Add(shorterTd, b.Difficulty())
		}
	}
	for _, b := range heavyChain {
		shorterTd.Add(shorterTd, b.Difficulty())
	}
	if shorterTd.Cmp(longerTd) <= 0 {
		return nil, nil, nil, fmt.Errorf("Test is moot, heavyChain td (%v) must be larger than canon td (%v)", shorterTd, longerTd)
	}
	longerNum := longChain[len(longChain)-1].NumberU64()
	shorterNum := heavyChain[len(heavyChain)-1].NumberU64()
	if shorterNum >= longerNum {
		return nil, nil, nil, fmt.Errorf("Test is moot, heavyChain num (%v) must be lower than canon num (%v)", shorterNum, longerNum)
	}
	return chain, longChain, heavyChain, nil
}

// TestReorgToShorterRemovesCanonMapping tests that if we
// 1. Have a chain [0 ... N .. X]
// 2. Reorg to shorter but heavier chain [0 ... N ... Y]
// 3. Then there should be no canon mapping for the block at height X
func TestReorgToShorterRemovesCanonMapping(t *testing.T) {
	chain, canonblocks, sideblocks, err := getLongAndShortChains()
	if err != nil {
		t.Fatal(err)
	}
	if n, err := chain.InsertChain(canonblocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	canonNum := chain.CurrentBlock().NumberU64()
	_, err = chain.InsertChain(sideblocks)
	if err != nil {
		t.Errorf("Got error, %v", err)
	}
	head := chain.CurrentBlock()
	if got := sideblocks[len(sideblocks)-1].Hash(); got != head.Hash() {
		t.Fatalf("head wrong, expected %x got %x", head.Hash(), got)
	}
	// We have now inserted a sidechain.
	if blockByNum := chain.GetBlockByNumber(canonNum); blockByNum != nil {
		t.Errorf("expected block to be gone: %v", blockByNum.NumberU64())
	}
	if headerByNum := chain.GetHeaderByNumber(canonNum); headerByNum != nil {
		t.Errorf("expected header to be gone: %v", headerByNum.Number.Uint64())
	}
}

// TestReorgToShorterRemovesCanonMappingHeaderChain is the same scenario
// as TestReorgToShorterRemovesCanonMapping, but applied on headerchain
// imports -- that is, for fast sync
func TestReorgToShorterRemovesCanonMappingHeaderChain(t *testing.T) {
	chain, canonblocks, sideblocks, err := getLongAndShortChains()
	if err != nil {
		t.Fatal(err)
	}
	// Convert into headers
	canonHeaders := make([]*types.Header[nist.PublicKey], len(canonblocks))
	for i, block := range canonblocks {
		canonHeaders[i] = block.Header()
	}
	if n, err := chain.InsertHeaderChain(canonHeaders, 0); err != nil {
		t.Fatalf("header %d: failed to insert into chain: %v", n, err)
	}
	canonNum := chain.CurrentHeader().Number.Uint64()
	sideHeaders := make([]*types.Header[nist.PublicKey], len(sideblocks))
	for i, block := range sideblocks {
		sideHeaders[i] = block.Header()
	}
	if n, err := chain.InsertHeaderChain(sideHeaders, 0); err != nil {
		t.Fatalf("header %d: failed to insert into chain: %v", n, err)
	}
	head := chain.CurrentHeader()
	if got := sideblocks[len(sideblocks)-1].Hash(); got != head.Hash() {
		t.Fatalf("head wrong, expected %x got %x", head.Hash(), got)
	}
	// We have now inserted a sidechain.
	if blockByNum := chain.GetBlockByNumber(canonNum); blockByNum != nil {
		t.Errorf("expected block to be gone: %v", blockByNum.NumberU64())
	}
	if headerByNum := chain.GetHeaderByNumber(canonNum); headerByNum != nil {
		t.Errorf("expected header to be gone: %v", headerByNum.Number.Uint64())
	}
}

func TestTransactionIndices(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{address: {Balance: funds}}}
		genesis = gspec.MustCommit(gendb)
		signer  = types.LatestSigner[nist.PublicKey](gspec.Config)
	)
	height := uint64(128)
	blocks, receipts := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), gendb, int(height), func(i int, block *BlockGen[nist.PublicKey]) {
		tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), params.TxGas, nil, nil), signer, key)
		if err != nil {
			panic(err)
		}
		block.AddTx(tx)
	})
	blocks2, _ := GenerateChain[nist.PublicKey](gspec.Config, blocks[len(blocks)-1],  ethash.NewFaker[nist.PublicKey](), gendb, 10, nil)

	check := func(tail *uint64, chain *BlockChain[nist.PublicKey]) {
		stored := rawdb.ReadTxIndexTail(chain.db)
		if tail == nil && stored != nil {
			t.Fatalf("Oldest indexded block mismatch, want nil, have %d", *stored)
		}
		if tail != nil && *stored != *tail {
			t.Fatalf("Oldest indexded block mismatch, want %d, have %d", *tail, *stored)
		}
		if tail != nil {
			for i := *tail; i <= chain.CurrentBlock().NumberU64(); i++ {
				block := rawdb.ReadBlock[nist.PublicKey](chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index == nil {
						t.Fatalf("Miss transaction indice, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
			for i := uint64(0); i < *tail; i++ {
				block := rawdb.ReadBlock[nist.PublicKey](chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index != nil {
						t.Fatalf("Transaction indice should be deleted, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
		}
	}
	frdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp freezer dir: %v", err)
	}
	defer os.Remove(frdir)
	ancientDb, err := rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), frdir, "", false)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	gspec.MustCommit(ancientDb)

	// Import all blocks into ancient db
	l := uint64(0)
	chain, err := NewBlockChain[nist.PublicKey](ancientDb, nil, params.TestChainConfig,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, &l, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	headers := make([]*types.Header[nist.PublicKey], len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := chain.InsertHeaderChain(headers, 0); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	if n, err := chain.InsertReceiptChain(blocks, receipts, 128); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	chain.Stop()
	ancientDb.Close()

	// Init block chain with external ancients, check all needed indices has been indexed.
	limit := []uint64{0, 32, 64, 128}
	for _, l := range limit {
		ancientDb, err = rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), frdir, "", false)
		if err != nil {
			t.Fatalf("failed to create temp freezer db: %v", err)
		}
		gspec.MustCommit(ancientDb)
		chain, err = NewBlockChain[nist.PublicKey](ancientDb, nil, params.TestChainConfig,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, &l, nil)
		if err != nil {
			t.Fatalf("failed to create tester chain: %v", err)
		}
		time.Sleep(50 * time.Millisecond) // Wait for indices initialisation
		var tail uint64
		if l != 0 {
			tail = uint64(128) - l + 1
		}
		check(&tail, chain)
		chain.Stop()
		ancientDb.Close()
	}

	// Reconstruct a block chain which only reserves HEAD-64 tx indices
	ancientDb, err = rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), frdir, "", false)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	gspec.MustCommit(ancientDb)

	limit = []uint64{0, 64 /* drop stale */, 32 /* shorten history */, 64 /* extend history */, 0 /* restore all */}
	tails := []uint64{0, 67 /* 130 - 64 + 1 */, 100 /* 131 - 32 + 1 */, 69 /* 132 - 64 + 1 */, 0}
	for i, l := range limit {
		chain, err = NewBlockChain[nist.PublicKey](ancientDb, nil, params.TestChainConfig,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, &l, nil)
		if err != nil {
			t.Fatalf("failed to create tester chain: %v", err)
		}
		chain.InsertChain(blocks2[i : i+1]) // Feed chain a higher block to trigger indices updater.
		time.Sleep(50 * time.Millisecond)   // Wait for indices initialisation
		check(&tails[i], chain)
		chain.Stop()
	}
}

func TestSkipStaleTxIndicesInFastSync(t *testing.T) {
	// Configure and generate a sample block chain
	var (
		gendb   = rawdb.NewMemoryDatabase()
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis[nist.PublicKey]{Config: params.TestChainConfig, Alloc: GenesisAlloc{address: {Balance: funds}}}
		genesis = gspec.MustCommit(gendb)
		signer  = types.LatestSigner[nist.PublicKey](gspec.Config)
	)
	height := uint64(128)
	blocks, receipts := GenerateChain[nist.PublicKey](gspec.Config, genesis,  ethash.NewFaker[nist.PublicKey](), gendb, int(height), func(i int, block *BlockGen[nist.PublicKey]) {
		tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(address), common.Address{0x00}, big.NewInt(1000), params.TxGas, nil, nil), signer, key)
		if err != nil {
			panic(err)
		}
		block.AddTx(tx)
	})

	check := func(tail *uint64, chain *BlockChain[nist.PublicKey]) {
		stored := rawdb.ReadTxIndexTail(chain.db)
		if tail == nil && stored != nil {
			t.Fatalf("Oldest indexded block mismatch, want nil, have %d", *stored)
		}
		if tail != nil && *stored != *tail {
			t.Fatalf("Oldest indexded block mismatch, want %d, have %d", *tail, *stored)
		}
		if tail != nil {
			for i := *tail; i <= chain.CurrentBlock().NumberU64(); i++ {
				block := rawdb.ReadBlock[nist.PublicKey](chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index == nil {
						t.Fatalf("Miss transaction indice, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
			for i := uint64(0); i < *tail; i++ {
				block := rawdb.ReadBlock[nist.PublicKey](chain.db, rawdb.ReadCanonicalHash(chain.db, i), i)
				if block.Transactions().Len() == 0 {
					continue
				}
				for _, tx := range block.Transactions() {
					if index := rawdb.ReadTxLookupEntry(chain.db, tx.Hash()); index != nil {
						t.Fatalf("Transaction indice should be deleted, number %d hash %s", i, tx.Hash().Hex())
					}
				}
			}
		}
	}

	frdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp freezer dir: %v", err)
	}
	defer os.Remove(frdir)
	ancientDb, err := rawdb.NewDatabaseWithFreezer[nist.PublicKey](rawdb.NewMemoryDatabase(), frdir, "", false)
	if err != nil {
		t.Fatalf("failed to create temp freezer db: %v", err)
	}
	gspec.MustCommit(ancientDb)

	// Import all blocks into ancient db, only HEAD-32 indices are kept.
	l := uint64(32)
	chain, err := NewBlockChain[nist.PublicKey](ancientDb, nil, params.TestChainConfig,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, &l, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	headers := make([]*types.Header[nist.PublicKey], len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	if n, err := chain.InsertHeaderChain(headers, 0); err != nil {
		t.Fatalf("failed to insert header %d: %v", n, err)
	}
	// The indices before ancient-N(32) should be ignored. After that all blocks should be indexed.
	if n, err := chain.InsertReceiptChain(blocks, receipts, 64); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	tail := uint64(32)
	check(&tail, chain)
}

// Benchmarks large blocks with value transfers to non-existing accounts
func benchmarkLargeNumberOfValueToNonexisting(b *testing.B, numTxs, numBlocks int, recipientFn func(uint64) common.Address, dataFn func(uint64) []byte) {
	var (
		signer          = types.HomesteadSigner[nist.PublicKey]{}
		testBankKey, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		testBankAddress = crypto.PubkeyToAddress[nist.PublicKey](*testBankKey.Public())
		bankFunds       = big.NewInt(100000000000000000)
		gspec           = Genesis[nist.PublicKey]{
			Config: params.TestChainConfig,
			Alloc: GenesisAlloc{
				testBankAddress: {Balance: bankFunds},
				common.HexToAddress("0xc0de"): {
					Code:    []byte{0x60, 0x01, 0x50},
					Balance: big.NewInt(0),
				}, // push 1, pop
			},
			GasLimit: 100e6, // 100 M
		}
	)
	// Generate the original common chain segment and the two competing forks
	engine :=  ethash.NewFaker[nist.PublicKey]()
	db := rawdb.NewMemoryDatabase()
	genesis := gspec.MustCommit(db)

	blockGenerator := func(i int, block *BlockGen[nist.PublicKey]) {
		block.SetCoinbase(common.Address{1})
		for txi := 0; txi < numTxs; txi++ {
			uniq := uint64(i*numTxs + txi)
			recipient := recipientFn(uniq)
			tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](uniq, recipient, big.NewInt(1), params.TxGas, big.NewInt(1), nil), signer, testBankKey)
			if err != nil {
				b.Error(err)
			}
			block.AddTx(tx)
		}
	}

	shared, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, numBlocks, blockGenerator)
	b.StopTimer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Import the shared chain and the original canonical one
		diskdb := rawdb.NewMemoryDatabase()
		gspec.MustCommit(diskdb)

		chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
		if err != nil {
			b.Fatalf("failed to create tester chain: %v", err)
		}
		b.StartTimer()
		if _, err := chain.InsertChain(shared); err != nil {
			b.Fatalf("failed to insert shared chain: %v", err)
		}
		b.StopTimer()
		if got := chain.CurrentBlock().Transactions().Len(); got != numTxs*numBlocks {
			b.Fatalf("Transactions were not included, expected %d, got %d", numTxs*numBlocks, got)

		}
	}
}

func BenchmarkBlockChain_1x1000ValueTransferToNonexisting(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)
	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(1337 + nonce))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}
	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}

func BenchmarkBlockChain_1x1000ValueTransferToExisting(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)
	b.StopTimer()
	b.ResetTimer()

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(1337))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}
	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}

func BenchmarkBlockChain_1x1000Executions(b *testing.B) {
	var (
		numTxs    = 1000
		numBlocks = 1
	)
	b.StopTimer()
	b.ResetTimer()

	recipientFn := func(nonce uint64) common.Address {
		return common.BigToAddress(big.NewInt(0).SetUint64(0xc0de))
	}
	dataFn := func(nonce uint64) []byte {
		return nil
	}
	benchmarkLargeNumberOfValueToNonexisting(b, numTxs, numBlocks, recipientFn, dataFn)
}

// Tests that importing a some old blocks, where all blocks are before the
// pruning point.
// This internally leads to a sidechain import, since the blocks trigger an
// ErrPrunedAncestor error.
// This may e.g. happen if
//   1. Downloader rollbacks a batch of inserted blocks and exits
//   2. Downloader starts to sync again
//   3. The blocks fetched are all known and canonical blocks
func TestSideImportPrunedBlocks(t *testing.T) {
	// Generate a canonical chain to act as the main dataset
	engine :=  ethash.NewFaker[nist.PublicKey]()
	db := rawdb.NewMemoryDatabase()
	genesis := new(Genesis[nist.PublicKey]).MustCommit(db)

	// Generate and import the canonical chain
	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 2*TriesInMemory, nil)
	diskdb := rawdb.NewMemoryDatabase()
	new(Genesis[nist.PublicKey]).MustCommit(diskdb)
	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}

	lastPrunedIndex := len(blocks) - TriesInMemory - 1
	lastPrunedBlock := blocks[lastPrunedIndex]

	// Verify pruning of lastPrunedBlock
	if chain.HasBlockAndState(lastPrunedBlock.Hash(), lastPrunedBlock.NumberU64()) {
		t.Errorf("Block %d not pruned", lastPrunedBlock.NumberU64())
	}
	firstNonPrunedBlock := blocks[len(blocks)-TriesInMemory]
	// Verify firstNonPrunedBlock is not pruned
	if !chain.HasBlockAndState(firstNonPrunedBlock.Hash(), firstNonPrunedBlock.NumberU64()) {
		t.Errorf("Block %d pruned", firstNonPrunedBlock.NumberU64())
	}
	// Now re-import some old blocks
	blockToReimport := blocks[5:8]
	_, err = chain.InsertChain(blockToReimport)
	if err != nil {
		t.Errorf("Got error, %v", err)
	}
}

// TestDeleteCreateRevert tests a weird state transition corner case that we hit
// while changing the internals of statedb. The workflow is that a contract is
// self destructed, then in a followup transaction (but same block) it's created
// again and the transaction reverted.
//
// The original statedb implementation flushed dirty objects to the tries after
// each transaction, so this works ok. The rework accumulated writes in memory
// first, but the journal wiped the entire state object on create-revert.
func TestDeleteCreateRevert(t *testing.T) {
	var (
		aa = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		bb = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		// Generate a canonical chain to act as the main dataset
		engine =  ethash.NewFaker[nist.PublicKey]()
		db     = rawdb.NewMemoryDatabase()

		// A sender who makes transactions, has some funds
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		gspec   = &Genesis[nist.PublicKey]{
			Config: params.TestChainConfig,
			Alloc: GenesisAlloc{
				address: {Balance: funds},
				// The address 0xAAAAA selfdestructs if called
				aa: {
					// Code needs to just selfdestruct
					Code:    []byte{byte(vm.PC), byte(vm.SELFDESTRUCT)},
					Nonce:   1,
					Balance: big.NewInt(0),
				},
				// The address 0xBBBB send 1 wei to 0xAAAA, then reverts
				bb: {
					Code: []byte{
						byte(vm.PC),          // [0]
						byte(vm.DUP1),        // [0,0]
						byte(vm.DUP1),        // [0,0,0]
						byte(vm.DUP1),        // [0,0,0,0]
						byte(vm.PUSH1), 0x01, // [0,0,0,0,1] (value)
						byte(vm.PUSH2), 0xaa, 0xaa, // [0,0,0,0,1, 0xaaaa]
						byte(vm.GAS),
						byte(vm.CALL),
						byte(vm.REVERT),
					},
					Balance: big.NewInt(1),
				},
			},
		}
		genesis = gspec.MustCommit(db)
	)

	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 1, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{1})
		// One transaction to AAAA
		tx, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](0, aa,
			big.NewInt(0), 50000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		b.AddTx(tx)
		// One transaction to BBBB
		tx, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](1, bb,
			big.NewInt(0), 100000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		b.AddTx(tx)
	})
	// Import the canonical chain
	diskdb := rawdb.NewMemoryDatabase()
	gspec.MustCommit(diskdb)

	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
}

// TestDeleteRecreateSlots tests a state-transition that contains both deletion
// and recreation of contract state.
// Contract A exists, has slots 1 and 2 set
// Tx 1: Selfdestruct A
// Tx 2: Re-create A, set slots 3 and 4
// Expected outcome is that _all_ slots are cleared from A, due to the selfdestruct,
// and then the new slots exist
func TestDeleteRecreateSlots(t *testing.T) {
	var (
		// Generate a canonical chain to act as the main dataset
		engine =  ethash.NewFaker[nist.PublicKey]()
		db     = rawdb.NewMemoryDatabase()
		// A sender who makes transactions, has some funds
		key, _    = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address   = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds     = big.NewInt(1000000000)
		bb        = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		aaStorage = make(map[common.Hash]common.Hash)          // Initial storage in AA
		aaCode    = []byte{byte(vm.PC), byte(vm.SELFDESTRUCT)} // Code for AA (simple selfdestruct)
	)
	// Populate two slots
	aaStorage[common.HexToHash("01")] = common.HexToHash("01")
	aaStorage[common.HexToHash("02")] = common.HexToHash("02")

	// The bb-code needs to CREATE2 the aa contract. It consists of
	// both initcode and deployment code
	// initcode:
	// 1. Set slots 3=3, 4=4,
	// 2. Return aaCode

	initCode := []byte{
		byte(vm.PUSH1), 0x3, // value
		byte(vm.PUSH1), 0x3, // location
		byte(vm.SSTORE),     // Set slot[3] = 1
		byte(vm.PUSH1), 0x4, // value
		byte(vm.PUSH1), 0x4, // location
		byte(vm.SSTORE), // Set slot[4] = 1
		// Slots are set, now return the code
		byte(vm.PUSH2), byte(vm.PC), byte(vm.SELFDESTRUCT), // Push code on stack
		byte(vm.PUSH1), 0x0, // memory start on stack
		byte(vm.MSTORE),
		// Code is now in memory.
		byte(vm.PUSH1), 0x2, // size
		byte(vm.PUSH1), byte(32 - 2), // offset
		byte(vm.RETURN),
	}
	if l := len(initCode); l > 32 {
		t.Fatalf("init code is too long for a pushx, need a more elaborate deployer")
	}
	bbCode := []byte{
		// Push initcode onto stack
		byte(vm.PUSH1) + byte(len(initCode)-1)}
	bbCode = append(bbCode, initCode...)
	bbCode = append(bbCode, []byte{
		byte(vm.PUSH1), 0x0, // memory start on stack
		byte(vm.MSTORE),
		byte(vm.PUSH1), 0x00, // salt
		byte(vm.PUSH1), byte(len(initCode)), // size
		byte(vm.PUSH1), byte(32 - len(initCode)), // offset
		byte(vm.PUSH1), 0x00, // endowment
		byte(vm.CREATE2),
	}...)

	initHash := crypto.Keccak256Hash[nist.PublicKey](initCode)
	aa := crypto.CreateAddress2[nist.PublicKey](bb, [32]byte{}, initHash[:])
	t.Logf("Destination address: %x\n", aa)

	gspec := &Genesis[nist.PublicKey]{
		Config: params.TestChainConfig,
		Alloc: GenesisAlloc{
			address: {Balance: funds},
			// The address 0xAAAAA selfdestructs if called
			aa: {
				// Code needs to just selfdestruct
				Code:    aaCode,
				Nonce:   1,
				Balance: big.NewInt(0),
				Storage: aaStorage,
			},
			// The contract BB recreates AA
			bb: {
				Code:    bbCode,
				Balance: big.NewInt(1),
			},
		},
	}
	genesis := gspec.MustCommit(db)

	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 1, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{1})
		// One transaction to AA, to kill it
		tx, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](0, aa,
			big.NewInt(0), 50000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		b.AddTx(tx)
		// One transaction to BB, to recreate AA
		tx, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](1, bb,
			big.NewInt(0), 100000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		b.AddTx(tx)
	})
	// Import the canonical chain
	diskdb := rawdb.NewMemoryDatabase()
	gspec.MustCommit(diskdb)
	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{
		Debug:  true,
		Tracer: vm.NewJSONLogger[nist.PublicKey](nil, os.Stdout),
	}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	statedb, _, _ := chain.State()

	// If all is correct, then slot 1 and 2 are zero
	if got, exp := statedb.GetState(aa, common.HexToHash("01")), (common.Hash{}); got != exp {
		t.Errorf("got %x exp %x", got, exp)
	}
	if got, exp := statedb.GetState(aa, common.HexToHash("02")), (common.Hash{}); got != exp {
		t.Errorf("got %x exp %x", got, exp)
	}
	// Also, 3 and 4 should be set
	if got, exp := statedb.GetState(aa, common.HexToHash("03")), common.HexToHash("03"); got != exp {
		t.Fatalf("got %x exp %x", got, exp)
	}
	if got, exp := statedb.GetState(aa, common.HexToHash("04")), common.HexToHash("04"); got != exp {
		t.Fatalf("got %x exp %x", got, exp)
	}
}

// TestDeleteRecreateAccount tests a state-transition that contains deletion of a
// contract with storage, and a recreate of the same contract via a
// regular value-transfer
// Expected outcome is that _all_ slots are cleared from A
func TestDeleteRecreateAccount(t *testing.T) {
	var (
		// Generate a canonical chain to act as the main dataset
		engine =  ethash.NewFaker[nist.PublicKey]()
		db     = rawdb.NewMemoryDatabase()
		// A sender who makes transactions, has some funds
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)

		aa        = common.HexToAddress("0x7217d81b76bdd8707601e959454e3d776aee5f43")
		aaStorage = make(map[common.Hash]common.Hash)          // Initial storage in AA
		aaCode    = []byte{byte(vm.PC), byte(vm.SELFDESTRUCT)} // Code for AA (simple selfdestruct)
	)
	// Populate two slots
	aaStorage[common.HexToHash("01")] = common.HexToHash("01")
	aaStorage[common.HexToHash("02")] = common.HexToHash("02")

	gspec := &Genesis[nist.PublicKey]{
		Config: params.TestChainConfig,
		Alloc: GenesisAlloc{
			address: {Balance: funds},
			// The address 0xAAAAA selfdestructs if called
			aa: {
				// Code needs to just selfdestruct
				Code:    aaCode,
				Nonce:   1,
				Balance: big.NewInt(0),
				Storage: aaStorage,
			},
		},
	}
	genesis := gspec.MustCommit(db)

	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 1, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{1})
		// One transaction to AA, to kill it
		tx, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](0, aa,
			big.NewInt(0), 50000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		b.AddTx(tx)
		// One transaction to AA, to recreate it (but without storage
		tx, _ = types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](1, aa,
			big.NewInt(1), 100000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		b.AddTx(tx)
	})
	// Import the canonical chain
	diskdb := rawdb.NewMemoryDatabase()
	gspec.MustCommit(diskdb)
	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{
		Debug:  true,
		Tracer: vm.NewJSONLogger[nist.PublicKey](nil, os.Stdout),
	}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}
	statedb, _, _ := chain.State()

	// If all is correct, then both slots are zero
	if got, exp := statedb.GetState(aa, common.HexToHash("01")), (common.Hash{}); got != exp {
		t.Errorf("got %x exp %x", got, exp)
	}
	if got, exp := statedb.GetState(aa, common.HexToHash("02")), (common.Hash{}); got != exp {
		t.Errorf("got %x exp %x", got, exp)
	}
}

// TestDeleteRecreateSlotsAcrossManyBlocks tests multiple state-transition that contains both deletion
// and recreation of contract state.
// Contract A exists, has slots 1 and 2 set
// Tx 1: Selfdestruct A
// Tx 2: Re-create A, set slots 3 and 4
// Expected outcome is that _all_ slots are cleared from A, due to the selfdestruct,
// and then the new slots exist
func TestDeleteRecreateSlotsAcrossManyBlocks(t *testing.T) {
	var (
		// Generate a canonical chain to act as the main dataset
		engine =  ethash.NewFaker[nist.PublicKey]()
		db     = rawdb.NewMemoryDatabase()
		// A sender who makes transactions, has some funds
		key, _    = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address   = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds     = big.NewInt(1000000000)
		bb        = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		aaStorage = make(map[common.Hash]common.Hash)          // Initial storage in AA
		aaCode    = []byte{byte(vm.PC), byte(vm.SELFDESTRUCT)} // Code for AA (simple selfdestruct)
	)
	// Populate two slots
	aaStorage[common.HexToHash("01")] = common.HexToHash("01")
	aaStorage[common.HexToHash("02")] = common.HexToHash("02")

	// The bb-code needs to CREATE2 the aa contract. It consists of
	// both initcode and deployment code
	// initcode:
	// 1. Set slots 3=blocknum+1, 4=4,
	// 2. Return aaCode

	initCode := []byte{
		byte(vm.PUSH1), 0x1, //
		byte(vm.NUMBER),     // value = number + 1
		byte(vm.ADD),        //
		byte(vm.PUSH1), 0x3, // location
		byte(vm.SSTORE),     // Set slot[3] = number + 1
		byte(vm.PUSH1), 0x4, // value
		byte(vm.PUSH1), 0x4, // location
		byte(vm.SSTORE), // Set slot[4] = 4
		// Slots are set, now return the code
		byte(vm.PUSH2), byte(vm.PC), byte(vm.SELFDESTRUCT), // Push code on stack
		byte(vm.PUSH1), 0x0, // memory start on stack
		byte(vm.MSTORE),
		// Code is now in memory.
		byte(vm.PUSH1), 0x2, // size
		byte(vm.PUSH1), byte(32 - 2), // offset
		byte(vm.RETURN),
	}
	if l := len(initCode); l > 32 {
		t.Fatalf("init code is too long for a pushx, need a more elaborate deployer")
	}
	bbCode := []byte{
		// Push initcode onto stack
		byte(vm.PUSH1) + byte(len(initCode)-1)}
	bbCode = append(bbCode, initCode...)
	bbCode = append(bbCode, []byte{
		byte(vm.PUSH1), 0x0, // memory start on stack
		byte(vm.MSTORE),
		byte(vm.PUSH1), 0x00, // salt
		byte(vm.PUSH1), byte(len(initCode)), // size
		byte(vm.PUSH1), byte(32 - len(initCode)), // offset
		byte(vm.PUSH1), 0x00, // endowment
		byte(vm.CREATE2),
	}...)

	initHash := crypto.Keccak256Hash[nist.PublicKey](initCode)
	aa := crypto.CreateAddress2[nist.PublicKey](bb, [32]byte{}, initHash[:])
	t.Logf("Destination address: %x\n", aa)
	gspec := &Genesis[nist.PublicKey]{
		Config: params.TestChainConfig,
		Alloc: GenesisAlloc{
			address: {Balance: funds},
			// The address 0xAAAAA selfdestructs if called
			aa: {
				// Code needs to just selfdestruct
				Code:    aaCode,
				Nonce:   1,
				Balance: big.NewInt(0),
				Storage: aaStorage,
			},
			// The contract BB recreates AA
			bb: {
				Code:    bbCode,
				Balance: big.NewInt(1),
			},
		},
	}
	genesis := gspec.MustCommit(db)
	var nonce uint64

	type expectation struct {
		exist    bool
		blocknum int
		values   map[int]int
	}
	var current = &expectation{
		exist:    true, // exists in genesis
		blocknum: 0,
		values:   map[int]int{1: 1, 2: 2},
	}
	var expectations []*expectation
	var newDestruct = func(e *expectation) *types.Transaction[nist.PublicKey] {
		tx, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](nonce, aa,
			big.NewInt(0), 50000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		nonce++
		if e.exist {
			e.exist = false
			e.values = nil
		}
		t.Logf("block %d; adding destruct\n", e.blocknum)
		return tx
	}
	var newResurrect = func(e *expectation) *types.Transaction[nist.PublicKey] {
		tx, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](nonce, bb,
			big.NewInt(0), 100000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		nonce++
		if !e.exist {
			e.exist = true
			e.values = map[int]int{3: e.blocknum + 1, 4: 4}
		}
		t.Logf("block %d; adding resurrect\n", e.blocknum)
		return tx
	}

	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 150, func(i int, b *BlockGen[nist.PublicKey]) {
		var exp = new(expectation)
		exp.blocknum = i + 1
		exp.values = make(map[int]int)
		for k, v := range current.values {
			exp.values[k] = v
		}
		exp.exist = current.exist

		b.SetCoinbase(common.Address{1})
		if i%2 == 0 {
			b.AddTx(newDestruct(exp))
		}
		if i%3 == 0 {
			b.AddTx(newResurrect(exp))
		}
		if i%5 == 0 {
			b.AddTx(newDestruct(exp))
		}
		if i%7 == 0 {
			b.AddTx(newResurrect(exp))
		}
		expectations = append(expectations, exp)
		current = exp
	})
	// Import the canonical chain
	diskdb := rawdb.NewMemoryDatabase()
	gspec.MustCommit(diskdb)
	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{
		//Debug: true,
		//Tracer: vm.NewJSONLogger[nist.PublicKey](nil, os.Stdout),
	}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	var asHash = func(num int) common.Hash {
		return common.BytesToHash([]byte{byte(num)})
	}
	for i, block := range blocks {
		blockNum := i + 1
		if n, err := chain.InsertChain([]*types.Block[nist.PublicKey]{block}); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", n, err)
		}
		statedb, _, _ := chain.State()
		// If all is correct, then slot 1 and 2 are zero
		if got, exp := statedb.GetState(aa, common.HexToHash("01")), (common.Hash{}); got != exp {
			t.Errorf("block %d, got %x exp %x", blockNum, got, exp)
		}
		if got, exp := statedb.GetState(aa, common.HexToHash("02")), (common.Hash{}); got != exp {
			t.Errorf("block %d, got %x exp %x", blockNum, got, exp)
		}
		exp := expectations[i]
		if exp.exist {
			if !statedb.Exist(aa) {
				t.Fatalf("block %d, expected %v to exist, it did not", blockNum, aa)
			}
			for slot, val := range exp.values {
				if gotValue, expValue := statedb.GetState(aa, asHash(slot)), asHash(val); gotValue != expValue {
					t.Fatalf("block %d, slot %d, got %x exp %x", blockNum, slot, gotValue, expValue)
				}
			}
		} else {
			if statedb.Exist(aa) {
				t.Fatalf("block %d, expected %v to not exist, it did", blockNum, aa)
			}
		}
	}
}

// TestInitThenFailCreateContract tests a pretty notorious case that happened
// on mainnet over blocks 7338108, 7338110 and 7338115.
// - Block 7338108: address e771789f5cccac282f23bb7add5690e1f6ca467c is initiated
//   with 0.001 ether (thus created but no code)
// - Block 7338110: a CREATE2 is attempted. The CREATE2 would deploy code on
//   the same address e771789f5cccac282f23bb7add5690e1f6ca467c. However, the
//   deployment fails due to OOG during initcode execution
// - Block 7338115: another tx checks the balance of
//   e771789f5cccac282f23bb7add5690e1f6ca467c, and the snapshotter returned it as
//   zero.
//
// The problem being that the snapshotter maintains a destructset, and adds items
// to the destructset in case something is created "onto" an existing item.
// We need to either roll back the snapDestructs, or not place it into snapDestructs
// in the first place.
//
func TestInitThenFailCreateContract(t *testing.T) {
	var (
		// Generate a canonical chain to act as the main dataset
		engine =  ethash.NewFaker[nist.PublicKey]()
		db     = rawdb.NewMemoryDatabase()
		// A sender who makes transactions, has some funds
		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
		funds   = big.NewInt(1000000000)
		bb      = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
	)

	// The bb-code needs to CREATE2 the aa contract. It consists of
	// both initcode and deployment code
	// initcode:
	// 1. If blocknum < 1, error out (e.g invalid opcode)
	// 2. else, return a snippet of code
	initCode := []byte{
		byte(vm.PUSH1), 0x1, // y (2)
		byte(vm.NUMBER), // x (number)
		byte(vm.GT),     // x > y?
		byte(vm.PUSH1), byte(0x8),
		byte(vm.JUMPI), // jump to label if number > 2
		byte(0xFE),     // illegal opcode
		byte(vm.JUMPDEST),
		byte(vm.PUSH1), 0x2, // size
		byte(vm.PUSH1), 0x0, // offset
		byte(vm.RETURN), // return 2 bytes of zero-code
	}
	if l := len(initCode); l > 32 {
		t.Fatalf("init code is too long for a pushx, need a more elaborate deployer")
	}
	bbCode := []byte{
		// Push initcode onto stack
		byte(vm.PUSH1) + byte(len(initCode)-1)}
	bbCode = append(bbCode, initCode...)
	bbCode = append(bbCode, []byte{
		byte(vm.PUSH1), 0x0, // memory start on stack
		byte(vm.MSTORE),
		byte(vm.PUSH1), 0x00, // salt
		byte(vm.PUSH1), byte(len(initCode)), // size
		byte(vm.PUSH1), byte(32 - len(initCode)), // offset
		byte(vm.PUSH1), 0x00, // endowment
		byte(vm.CREATE2),
	}...)

	initHash := crypto.Keccak256Hash[nist.PublicKey](initCode)
	aa := crypto.CreateAddress2[nist.PublicKey](bb, [32]byte{}, initHash[:])
	t.Logf("Destination address: %x\n", aa)

	gspec := &Genesis[nist.PublicKey]{
		Config: params.TestChainConfig,
		Alloc: GenesisAlloc{
			address: {Balance: funds},
			// The address aa has some funds
			aa: {Balance: big.NewInt(100000)},
			// The contract BB tries to create code onto AA
			bb: {
				Code:    bbCode,
				Balance: big.NewInt(1),
			},
		},
	}
	genesis := gspec.MustCommit(db)
	nonce := uint64(0)
	blocks, _ := GenerateChain[nist.PublicKey](params.TestChainConfig, genesis, engine, db, 4, func(i int, b *BlockGen[nist.PublicKey]) {
		b.SetCoinbase(common.Address{1})
		// One transaction to BB
		tx, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](nonce, bb,
			big.NewInt(0), 100000, big.NewInt(1), nil), types.HomesteadSigner[nist.PublicKey]{}, key)
		b.AddTx(tx)
		nonce++
	})

	// Import the canonical chain
	diskdb := rawdb.NewMemoryDatabase()
	gspec.MustCommit(diskdb)
	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, params.TestChainConfig, engine, vm.Config[nist.PublicKey]{
		//Debug:  true,
		//Tracer: vm.NewJSONLogger[nist.PublicKey](nil, os.Stdout),
	}, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	statedb, _, _ := chain.State()
	if got, exp := statedb.GetBalance(aa), big.NewInt(100000); got.Cmp(exp) != 0 {
		t.Fatalf("Genesis err, got %v exp %v", got, exp)
	}
	// First block tries to create, but fails
	{
		block := blocks[0]
		if _, err := chain.InsertChain([]*types.Block[nist.PublicKey]{blocks[0]}); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", block.NumberU64(), err)
		}
		statedb, _, _ = chain.State()
		if got, exp := statedb.GetBalance(aa), big.NewInt(100000); got.Cmp(exp) != 0 {
			t.Fatalf("block %d: got %v exp %v", block.NumberU64(), got, exp)
		}
	}
	// Import the rest of the blocks
	for _, block := range blocks[1:] {
		if _, err := chain.InsertChain([]*types.Block[nist.PublicKey]{block}); err != nil {
			t.Fatalf("block %d: failed to insert into chain: %v", block.NumberU64(), err)
		}
	}
}

// TestEIP2718Transition tests that an EIP-2718 transaction will be accepted
// after the fork block has passed. This is verified by sending an EIP-2930
// access list transaction, which specifies a single slot access, and then
// checking that the gas usage of a hot SLOAD and a cold SLOAD are calculated
// correctly.
// func TestEIP2718Transition(t *testing.T) {
// 	var (
// 		aa = common.HexToAddress("0x000000000000000000000000000000000000aaaa")

// 		// Generate a canonical chain to act as the main dataset
// 		engine =  ethash.NewFaker[nist.PublicKey]()
// 		db     = rawdb.NewMemoryDatabase()

// 		// A sender who makes transactions, has some funds
// 		key, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
// 		address = crypto.PubkeyToAddress[nist.PublicKey](*key.Public())
// 		funds   = big.NewInt(1000000000)
// 		gspec   = &Genesis[nist.PublicKey]{
// 			Config: params.YoloV3ChainConfig,
// 			Alloc: GenesisAlloc{
// 				address: {Balance: funds},
// 				// The address 0xAAAA sloads 0x00 and 0x01
// 				aa: {
// 					Code: []byte{
// 						byte(vm.PC),
// 						byte(vm.PC),
// 						byte(vm.SLOAD),
// 						byte(vm.SLOAD),
// 					},
// 					Nonce:   0,
// 					Balance: big.NewInt(0),
// 				},
// 			},
// 		}
// 		genesis = gspec.MustCommit(db)
// 	)

// 	blocks, _ := GenerateChain[nist.PublicKey](gspec.Config, genesis, engine, db, 1, func(i int, b *BlockGen[nist.PublicKey]) {
// 		b.SetCoinbase(common.Address{1})

// 		// One transaction to 0xAAAA
// 		signer := types.LatestSigner[nist.PublicKey](gspec.Config)
// 		tx, _ := types.SignNewTx(key, signer, &types.AccessListTx{
// 			ChainID:  gspec.Config.ChainID,
// 			Nonce:    0,
// 			To:       &aa,
// 			Gas:      30000,
// 			GasPrice: big.NewInt(1),
// 			AccessList: types.AccessList{{
// 				Address:     aa,
// 				StorageKeys: []common.Hash{{0}},
// 			}},
// 		})
// 		b.AddTx(tx)
// 	})

// 	// Import the canonical chain
// 	diskdb := rawdb.NewMemoryDatabase()
// 	gspec.MustCommit(diskdb)

// 	chain, err := NewBlockChain[nist.PublicKey](diskdb, nil, gspec.Config, engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
// 	if err != nil {
// 		t.Fatalf("failed to create tester chain: %v", err)
// 	}
// 	if n, err := chain.InsertChain(blocks); err != nil {
// 		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
// 	}

// 	block := chain.GetBlockByNumber(1)

// 	// Expected gas is intrinsic + 2 * pc + hot load + cold load, since only one load is in the access list
// 	expected := params.TxGas + params.TxAccessListAddressGas + params.TxAccessListStorageKeyGas + vm.GasQuickStep*2 + vm.WarmStorageReadCostEIP2929 + vm.ColdSloadCostEIP2929
// 	if block.GasUsed() != expected {
// 		t.Fatalf("incorrect amount of gas spent: expected %d, got %d", expected, block.GasUsed())

// 	}
// }
