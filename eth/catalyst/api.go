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

// Package catalyst implements the temporary eth1/eth2 RPC integration.
package catalyst

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/eth"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/node"
	chainParams "github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// Register adds catalyst APIs to the node.
func Register[T crypto.PrivateKey, P crypto.PublicKey](stack *node.Node[T,P], backend *eth.Ethereum[T,P]) error {
	chainconfig := backend.BlockChain().Config()
	if chainconfig.CatalystBlock == nil {
		return errors.New("catalystBlock is not set in genesis config")
	} else if chainconfig.CatalystBlock.Sign() != 0 {
		return errors.New("catalystBlock of genesis config must be zero")
	}

	log.Warn("Catalyst mode enabled")
	stack.RegisterAPIs([]rpc.API{
		{
			Namespace: "consensus",
			Version:   "1.0",
			Service:   newConsensusAPI(backend),
			Public:    true,
		},
	})
	return nil
}

type consensusAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	eth *eth.Ethereum[T,P]
}

func newConsensusAPI[T crypto.PrivateKey, P crypto.PublicKey](eth *eth.Ethereum[T,P]) *consensusAPI[T,P] {
	return &consensusAPI[T,P]{eth: eth}
}

// blockExecutionEnv gathers all the data required to execute
// a block, either when assembling it or when inserting it.
type blockExecutionEnv [T crypto.PrivateKey, P crypto.PublicKey] struct {
	chain   *core.BlockChain[P]
	state   *state.StateDB
	tcount  int
	gasPool *core.GasPool

	header   *types.Header
	txs      []*types.Transaction[P]
	receipts []*types.Receipt[P]

	// Quorum
	privateStateRepo  mps.PrivateStateRepository
	privateState      *state.StateDB
	forceNonParty     bool
	isInnerPrivateTxn bool
	privateReceipts   []*types.Receipt[P]
}

func (env *blockExecutionEnv[T,P]) commitTransaction(tx *types.Transaction[P], coinbase common.Address) error {
	vmconfig := *env.chain.GetVMConfig()
	receipt, privateReceipt, err := core.ApplyTransaction[P](env.chain.Config(), env.chain, &coinbase, env.gasPool, env.state, env.privateState, env.header, tx, &env.header.GasUsed, vmconfig, env.forceNonParty, env.privateStateRepo, env.isInnerPrivateTxn)
	if err != nil {
		return err
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)
	env.privateReceipts = append(env.privateReceipts, privateReceipt)
	return nil
}

func (api *consensusAPI[T,P]) makeEnv(parent *types.Block[P], header *types.Header) (*blockExecutionEnv[T,P], error) {
	state, mpsr, err := api.eth.BlockChain().StateAt(parent.Root())
	if err != nil {
		return nil, err
	}
	privateState, err := mpsr.DefaultState() // TODO merge add PSI?
	if err != nil {
		return nil, err
	}
	env := &blockExecutionEnv[T,P]{
		chain:   api.eth.BlockChain(),
		state:   state,
		header:  header,
		gasPool: new(core.GasPool).AddGas(header.GasLimit),
		// Quorum
		privateState: privateState,
	}
	return env, nil
}

// AssembleBlock creates a new block, inserts it into the chain, and returns the "execution
// data" required for eth2 clients to process the new block.
func (api *consensusAPI[T,P]) AssembleBlock(params assembleBlockParams) (*executableData, error) {
	log.Info("Producing block", "parentHash", params.ParentHash)

	bc := api.eth.BlockChain()
	parent := bc.GetBlockByHash(params.ParentHash)
	if parent == nil {
		log.Warn("Cannot assemble block with parent hash to unknown block", "parentHash", params.ParentHash)
		return nil, fmt.Errorf("cannot assemble block with unknown parent %s", params.ParentHash)
	}

	pool := api.eth.TxPool()

	if parent.Time() >= params.Timestamp {
		return nil, fmt.Errorf("child timestamp lower than parent's: %d >= %d", parent.Time(), params.Timestamp)
	}
	if now := uint64(time.Now().Unix()); params.Timestamp > now+1 {
		wait := time.Duration(params.Timestamp-now) * time.Second
		log.Info("Producing block too far in the future", "wait", common.PrettyDuration(wait))
		time.Sleep(wait)
	}

	pending, err := pool.Pending()
	if err != nil {
		return nil, err
	}

	coinbase, err := api.eth.Etherbase()
	if err != nil {
		return nil, err
	}
	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		Coinbase:   coinbase,
		GasLimit:   parent.GasLimit(), // Keep the gas limit constant in this prototype
		Extra:      []byte{},
		Time:       params.Timestamp,
	}
	err = api.eth.Engine().Prepare(bc, header)
	if err != nil {
		return nil, err
	}

	env, err := api.makeEnv(parent, header)
	if err != nil {
		return nil, err
	}

	var (
		signer       = types.MakeSigner[P](bc.Config(), header.Number)
		txHeap       = types.NewTransactionsByPriceAndNonce(signer, pending)
		transactions []*types.Transaction[P]
	)
	for {
		if env.gasPool.Gas() < chainParams.TxGas {
			log.Trace("Not enough gas for further transactions", "have", env.gasPool, "want", chainParams.TxGas)
			break
		}
		tx := txHeap.Peek()
		if tx == nil {
			break
		}

		// The sender is only for logging purposes, and it doesn't really matter if it's correct.
		from, _ := types.Sender(signer, tx)

		// Execute the transaction
		env.state.Prepare(tx.Hash(), common.Hash{}, env.tcount)
		err = env.commitTransaction(tx, coinbase)
		switch err {
		case core.ErrGasLimitReached:
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Trace("Gas limit exceeded for current block", "sender", from)
			txHeap.Pop()

		case core.ErrNonceTooLow:
			// New head notification data race between the transaction pool and miner, shift
			log.Trace("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce())
			txHeap.Shift()

		case core.ErrNonceTooHigh:
			// Reorg notification data race between the transaction pool and miner, skip account =
			log.Trace("Skipping account with high nonce", "sender", from, "nonce", tx.Nonce())
			txHeap.Pop()

		case nil:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			env.tcount++
			txHeap.Shift()
			transactions = append(transactions, tx)

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txHeap.Shift()
		}
	}

	// Create the block.
	block, err := api.eth.Engine().FinalizeAndAssemble(bc, header, env.state, transactions, nil /* uncles */, env.receipts)
	if err != nil {
		return nil, err
	}
	return &executableData{
		BlockHash:    block.Hash(),
		ParentHash:   block.ParentHash(),
		Miner:        block.Coinbase(),
		StateRoot:    block.Root(),
		Number:       block.NumberU64(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		Timestamp:    block.Time(),
		ReceiptRoot:  block.ReceiptHash(),
		LogsBloom:    block.Bloom().Bytes(),
		Transactions: encodeTransactions(block.Transactions()),
	}, nil
}

func encodeTransactions[P crypto.PublicKey](txs []*types.Transaction[P]) [][]byte {
	var enc = make([][]byte, len(txs))
	for i, tx := range txs {
		enc[i], _ = tx.MarshalBinary()
	}
	return enc
}

func decodeTransactions[P crypto.PublicKey](enc [][]byte) ([]*types.Transaction[P], error) {
	var txs = make([]*types.Transaction[P], len(enc))
	for i, encTx := range enc {
		var tx types.Transaction[P]
		if err := tx.UnmarshalBinary(encTx); err != nil {
			return nil, fmt.Errorf("invalid transaction %d: %v", i, err)
		}
		txs[i] = &tx
	}
	return txs, nil
}

func insertBlockParamsToBlock[P crypto.PublicKey](params executableData) (*types.Block[P], error) {
	txs, err := decodeTransactions[P](params.Transactions)
	if err != nil {
		return nil, err
	}

	number := big.NewInt(0)
	number.SetUint64(params.Number)
	header := &types.Header{
		ParentHash:  params.ParentHash,
		UncleHash:   types.EmptyUncleHash,
		Coinbase:    params.Miner,
		Root:        params.StateRoot,
		TxHash:      types.DeriveSha(types.Transactions[P](txs), trie.NewStackTrie(nil)),
		ReceiptHash: params.ReceiptRoot,
		Bloom:       types.BytesToBloom(params.LogsBloom),
		Difficulty:  big.NewInt(1),
		Number:      number,
		GasLimit:    params.GasLimit,
		GasUsed:     params.GasUsed,
		Time:        params.Timestamp,
	}
	block := types.NewBlockWithHeader[P](header).WithBody(txs, nil /* uncles */)
	return block, nil
}

// NewBlock creates an Eth1 block, inserts it in the chain, and either returns true,
// or false + an error. This is a bit redundant for go, but simplifies things on the
// eth2 side.
func (api *consensusAPI[T,P]) NewBlock(params executableData) (*newBlockResponse, error) {
	parent := api.eth.BlockChain().GetBlockByHash(params.ParentHash)
	if parent == nil {
		return &newBlockResponse{false}, fmt.Errorf("could not find parent %x", params.ParentHash)
	}
	block, err := insertBlockParamsToBlock[P](params)
	if err != nil {
		return nil, err
	}

	_, err = api.eth.BlockChain().InsertChainWithoutSealVerification(block)
	return &newBlockResponse{err == nil}, err
}

// Used in tests to add a the list of transactions from a block to the tx pool.
func (api *consensusAPI[T,P]) addBlockTxs(block *types.Block[P]) error {
	for _, tx := range block.Transactions() {
		api.eth.TxPool().AddLocal(tx)
	}
	return nil
}

// FinalizeBlock is called to mark a block as synchronized, so
// that data that is no longer needed can be removed.
func (api *consensusAPI[T,P]) FinalizeBlock(blockHash common.Hash) (*genericResponse, error) {
	return &genericResponse{true}, nil
}

// SetHead is called to perform a force choice.
func (api *consensusAPI[T,P]) SetHead(newHead common.Hash) (*genericResponse, error) {
	return &genericResponse{true}, nil
}
