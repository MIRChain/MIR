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

package core

import (
	"fmt"
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/misc"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/params"
)

// BlockGen creates blocks for testing.
// See GenerateChain for a detailed explanation.
type BlockGen [P crypto.PublicKey] struct {
	i       int
	parent  *types.Block[P]
	chain   []*types.Block[P]
	header  *types.Header[P]
	statedb *state.StateDB[P]

	gasPool  *GasPool
	txs      []*types.Transaction[P]
	receipts []*types.Receipt[P]
	uncles   []*types.Header[P]

	config *params.ChainConfig
	engine consensus.Engine[P]

	privateStatedb *state.StateDB[P] // Quorum
}

// SetCoinbase sets the coinbase of the generated block.
// It can be called at most once.
func (b *BlockGen[P]) SetCoinbase(addr common.Address) {
	if b.gasPool != nil {
		if len(b.txs) > 0 {
			panic("coinbase must be set before adding transactions")
		}
		panic("coinbase can only be set once")
	}
	b.header.Coinbase = addr
	b.gasPool = new(GasPool).AddGas(b.header.GasLimit)
}

// SetExtra sets the extra data field of the generated block.
func (b *BlockGen[P]) SetExtra(data []byte) {
	b.header.Extra = data
}

// SetNonce sets the nonce field of the generated block.
func (b *BlockGen[P]) SetNonce(nonce types.BlockNonce) {
	b.header.Nonce = nonce
}

// SetDifficulty sets the difficulty field of the generated block. This method is
// useful for Clique tests where the difficulty does not depend on time. For the
// ethash tests, please use OffsetTime, which implicitly recalculates the diff.
func (b *BlockGen[P]) SetDifficulty(diff *big.Int) {
	b.header.Difficulty = diff
}

// AddTx adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// AddTx panics if the transaction cannot be executed. In addition to
// the protocol-imposed limitations (gas limit, etc.), there are some
// further limitations on the content of transactions that can be
// added. Notably, contract code relying on the BLOCKHASH instruction
// will panic during execution.
func (b *BlockGen[P]) AddTx(tx *types.Transaction[P]) {
	b.AddTxWithChain(&BlockChain[P]{quorumConfig: &QuorumChainConfig{}}, tx)
}

// AddTxWithChain adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// AddTxWithChain panics if the transaction cannot be executed. In addition to
// the protocol-imposed limitations (gas limit, etc.), there are some
// further limitations on the content of transactions that can be
// added. If contract code relies on the BLOCKHASH instruction,
// the block in chain will be returned.
func (b *BlockGen[P]) AddTxWithChain(bc *BlockChain[P], tx *types.Transaction[P]) {
	if b.gasPool == nil {
		b.SetCoinbase(common.Address{})
	}
	b.statedb.Prepare(tx.Hash(), common.Hash{}, len(b.txs))
	// Quorum
	privateDb := b.privateStatedb
	if privateDb == nil {
		privateDb = b.statedb
	} else {
		b.privateStatedb.Prepare(tx.Hash(), common.Hash{}, len(b.txs))
	}
	// End Quorum

	receipt, _, err := ApplyTransaction[P](b.config, bc, &b.header.Coinbase, b.gasPool, b.statedb, privateDb, b.header, tx, &b.header.GasUsed, vm.Config[P]{}, false, nil, false)
	if err != nil {
		panic(err)
	}
	b.txs = append(b.txs, tx)
	b.receipts = append(b.receipts, receipt)
}

// GetBalance returns the balance of the given address at the generated block.
func (b *BlockGen[P]) GetBalance(addr common.Address) *big.Int {
	return b.statedb.GetBalance(addr)
}

// AddUncheckedTx forcefully adds a transaction to the block without any
// validation.
//
// AddUncheckedTx will cause consensus failures when used during real
// chain processing. This is best used in conjunction with raw block insertion.
func (b *BlockGen[P]) AddUncheckedTx(tx *types.Transaction[P]) {
	b.txs = append(b.txs, tx)
}

// Number returns the block number of the block being generated.
func (b *BlockGen[P]) Number() *big.Int {
	return new(big.Int).Set(b.header.Number)
}

// AddUncheckedReceipt forcefully adds a receipts to the block without a
// backing transaction.
//
// AddUncheckedReceipt will cause consensus failures when used during real
// chain processing. This is best used in conjunction with raw block insertion.
func (b *BlockGen[P]) AddUncheckedReceipt(receipt *types.Receipt[P]) {
	b.receipts = append(b.receipts, receipt)
}

// TxNonce returns the next valid transaction nonce for the
// account at addr. It panics if the account does not exist.
func (b *BlockGen[P]) TxNonce(addr common.Address) uint64 {
	if !b.statedb.Exist(addr) {
		panic("account does not exist")
	}
	return b.statedb.GetNonce(addr)
}

// AddUncle adds an uncle header to the generated block.
func (b *BlockGen[P]) AddUncle(h *types.Header[P]) {
	b.uncles = append(b.uncles, h)
}

// PrevBlock returns a previously generated block by number. It panics if
// num is greater or equal to the number of the block being generated.
// For index -1, PrevBlock returns the parent block given to GenerateChain.
func (b *BlockGen[P]) PrevBlock(index int) *types.Block[P] {
	if index >= b.i {
		panic(fmt.Errorf("block index %d out of range (%d,%d)", index, -1, b.i))
	}
	if index == -1 {
		return b.parent
	}
	return b.chain[index]
}

// OffsetTime modifies the time instance of a block, implicitly changing its
// associated difficulty. It's useful to test scenarios where forking is not
// tied to chain length directly.
func (b *BlockGen[P]) OffsetTime(seconds int64) {
	b.header.Time += uint64(seconds)
	if b.header.Time <= b.parent.Header().Time {
		panic("block time out of range")
	}
	chainreader := &fakeChainReader[P]{config: b.config}
	b.header.Difficulty = b.engine.CalcDifficulty(chainreader, b.header.Time, b.parent.Header())
}

// GenerateChain creates a chain of n blocks. The first block's
// parent will be the provided parent. db is used to store
// intermediate states and should contain the parent's state trie.
//
// The generator function is called with a new block generator for
// every block. Any transactions and uncles added to the generator
// become part of the block. If gen is nil, the blocks will be empty
// and their coinbase will be the zero address.
//
// Blocks created by GenerateChain do not contain valid proof of work
// values. Inserting them into BlockChain requires use of FakePow or
// a similar non-validating proof of work implementation.
func GenerateChain[P crypto.PublicKey](config *params.ChainConfig, parent *types.Block[P], engine consensus.Engine[P], db ethdb.Database, n int, gen func(int, *BlockGen[P])) ([]*types.Block[P], []types.Receipts[P]) {
	if config == nil {
		config = params.TestChainConfig
	}
	blocks, receipts := make(types.Blocks[P], n), make([]types.Receipts[P], n)
	chainreader := &fakeChainReader[P]{config: config}
	// Quorum: add `privateStatedb` argument
	genblock := func(i int, parent *types.Block[P], statedb *state.StateDB[P], privateStatedb *state.StateDB[P]) (*types.Block[P], types.Receipts[P]) {
		b := &BlockGen[P]{i: i, chain: blocks, parent: parent, statedb: statedb, privateStatedb: privateStatedb, config: config, engine: engine}
		b.header = makeHeader[P](chainreader, parent, statedb, b.engine)

		// Mutate the state and block according to any hard-fork specs
		if daoBlock := config.DAOForkBlock; daoBlock != nil {
			limit := new(big.Int).Add(daoBlock, params.DAOForkExtraRange)
			if b.header.Number.Cmp(daoBlock) >= 0 && b.header.Number.Cmp(limit) < 0 {
				if config.DAOForkSupport {
					b.header.Extra = common.CopyBytes(params.DAOForkBlockExtra)
				}
			}
		}
		if config.DAOForkSupport && config.DAOForkBlock != nil && config.DAOForkBlock.Cmp(b.header.Number) == 0 {
			misc.ApplyDAOHardFork(statedb)
		}
		// Execute any user modifications to the block
		if gen != nil {
			gen(i, b)
		}
		if b.engine != nil {
			// Finalize and seal the block
			block, _ := b.engine.FinalizeAndAssemble(chainreader, b.header, statedb, b.txs, b.uncles, b.receipts)

			// Write state changes to db
			root, err := statedb.Commit(config.IsEIP158(b.header.Number))
			if err != nil {
				panic(fmt.Sprintf("state write error: %v", err))
			}
			if err := statedb.Database().TrieDB().Commit(root, false, nil); err != nil {
				panic(fmt.Sprintf("trie write error: %v", err))
			}
			return block, b.receipts
		}
		return nil, nil
	}
	for i := 0; i < n; i++ {
		statedb, err := state.New[P](parent.Root(), state.NewDatabase[P](db), nil)
		if err != nil {
			panic(err)
		}
		privateStatedb, err := state.New[P](parent.Root(), state.NewDatabase[P](db), nil) // Quorum
		if err != nil {
			panic(err)
		}
		// Quorum: add `privateStatedb` argument
		block, receipt := genblock(i, parent, statedb, privateStatedb)
		blocks[i] = block
		receipts[i] = receipt
		parent = block
	}
	return blocks, receipts
}

func makeHeader[P crypto.PublicKey](chain consensus.ChainReader[P], parent *types.Block[P], state *state.StateDB[P], engine consensus.Engine[P]) *types.Header[P] {
	var time uint64
	if parent.Time() == 0 {
		time = 10
	} else {
		time = parent.Time() + 10 // block time is fixed at 10 seconds
	}

	return &types.Header[P]{
		Root:       state.IntermediateRoot(chain.Config().IsEIP158(parent.Number())),
		ParentHash: parent.Hash(),
		Coinbase:   parent.Coinbase(),
		Difficulty: engine.CalcDifficulty(chain, time, &types.Header[P]{
			Number:     parent.Number(),
			Time:       time - 10,
			Difficulty: parent.Difficulty(),
			UncleHash:  parent.UncleHash(),
		}),
		GasLimit: CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit(), parent.GasLimit()),
		Number:   new(big.Int).Add(parent.Number(), common.Big1),
		Time:     time,
	}
}

// makeHeaderChain creates a deterministic chain of headers rooted at parent.
func makeHeaderChain[P crypto.PublicKey](parent *types.Header[P], n int, engine consensus.Engine[P], db ethdb.Database, seed int) []*types.Header[P] {
	blocks := makeBlockChain(types.NewBlockWithHeader[P](parent), n, engine, db, seed)
	headers := make([]*types.Header[P], len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	return headers
}

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain[P crypto.PublicKey](parent *types.Block[P], n int, engine consensus.Engine[P], db ethdb.Database, seed int) []*types.Block[P] {
	blocks, _ := GenerateChain(params.TestChainConfig, parent, engine, db, n, func(i int, b *BlockGen[P]) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return blocks
}

type fakeChainReader [P crypto.PublicKey] struct {
	config *params.ChainConfig
}

// Config returns the chain configuration.
func (cr *fakeChainReader[P]) Config() *params.ChainConfig {
	return cr.config
}

func (cr *fakeChainReader[P]) CurrentHeader() *types.Header[P]                            { return nil }
func (cr *fakeChainReader[P]) GetHeaderByNumber(number uint64) *types.Header[P]           { return nil }
func (cr *fakeChainReader[P]) GetHeaderByHash(hash common.Hash) *types.Header[P]          { return nil }
func (cr *fakeChainReader[P]) GetHeader(hash common.Hash, number uint64) *types.Header[P] { return nil }
func (cr *fakeChainReader[P]) GetBlock(hash common.Hash, number uint64) *types.Block[P]   { return nil }
