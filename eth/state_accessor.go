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

package eth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

// stateAtBlock retrieves the state database associated with a certain block.
// If no state is locally available for the given block, a number of blocks
// are attempted to be reexecuted to generate the desired state. The optional
// base layer statedb can be passed then it's regarded as the statedb of the
// parent block.
func (eth *Ethereum[T,P]) stateAtBlock(block *types.Block[P], reexec uint64, base *state.StateDB[P], checkLive bool) (statedb *state.StateDB[P], privateStateDB mps.PrivateStateRepository[P], err error) {
	var (
		current  *types.Block[P]
		database state.Database
		report   = true
		origin   = block.NumberU64()
	)
	// Check the live database first if we have the state fully available, use that.
	if checkLive {
		statedb, privateStateDB, err = eth.blockchain.StateAt(block.Root())
		if err == nil {
			return
		}
	}
	if base != nil {
		// The optional base statedb is given, mark the start point as parent block
		statedb, database, report = base, base.Database(), false
		current = eth.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	} else {
		// Otherwise try to reexec blocks until we find a state or reach our limit
		current = block

		// Create an ephemeral trie.Database for isolating the live one. Otherwise
		// the internal junks created by tracing will be persisted into the disk.
		database = state.NewDatabaseWithConfig[P](eth.chainDb, &trie.Config{Cache: 16})

		// If we didn't check the dirty database, do check the clean one, otherwise
		// we would rewind past a persisted block (specific corner case is chain
		// tracing from the genesis).
		if !checkLive {
			statedb, err = state.New[P](current.Root(), database, nil)
			if err == nil {
				// Quorum
				_, privateStateDB, err = eth.blockchain.StateAt(current.Root())
				if err == nil {
					return statedb, privateStateDB, nil
				}
				// End Quorum
			}
		}
		// Database does not have the state for the given block, try to regenerate
		for i := uint64(0); i < reexec; i++ {
			if current.NumberU64() == 0 {
				return nil, nil, errors.New("genesis state is missing")
			}
			parent := eth.blockchain.GetBlock(current.ParentHash(), current.NumberU64()-1)
			if parent == nil {
				return nil, nil, fmt.Errorf("missing block %v %d", current.ParentHash(), current.NumberU64()-1)
			}
			current = parent

			statedb, err = state.New[P](current.Root(), database, nil)
			if err == nil {
				// Quorum
				_, privateStateDB, err = eth.blockchain.StateAt(current.Root())
				if err == nil {
					break
				}
				// End Quorum
			}
		}
		if err != nil {
			switch err.(type) {
			case *trie.MissingNodeError:
				return nil, nil, fmt.Errorf("required historical state unavailable (reexec=%d)", reexec)
			default:
				return nil, nil, err
			}
		}
	}
	// State was available at historical point, regenerate
	var (
		start  = time.Now()
		logged time.Time
		parent common.Hash
	)
	for current.NumberU64() < origin {
		// Print progress logs if long enough time elapsed
		if time.Since(logged) > 8*time.Second && report {
			log.Info("Regenerating historical state", "block", current.NumberU64()+1, "target", origin, "remaining", origin-current.NumberU64()-1, "elapsed", time.Since(start))
			logged = time.Now()
		}
		// Retrieve the next block to regenerate and process it
		next := current.NumberU64() + 1
		if current = eth.blockchain.GetBlockByNumber(next); current == nil {
			return nil, nil, fmt.Errorf("block #%d not found", next)
		}
		_, _, _, _, err = eth.blockchain.Processor().Process(current, statedb, privateStateDB, vm.Config[P]{})
		if err != nil {
			return nil, nil, fmt.Errorf("processing block %d failed: %v", current.NumberU64(), err)
		}
		var root common.Hash
		// Finalize the state so any modifications are written to the trie
		root, err = statedb.Commit(eth.blockchain.Config().IsEIP158(current.Number()))
		if err != nil {
			return nil, nil, err
		}
		statedb, err = state.New[P](root, database, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("state reset after block %d failed: %v", current.NumberU64(), err)
		}

		// Quorum
		err = privateStateDB.Commit(eth.blockchain.Config().IsEIP158(block.Number()), block)
		if err != nil {
			return nil, nil, err
		}
		if err := privateStateDB.Reset(); err != nil {
			return nil, nil, fmt.Errorf("private state reset after block %d failed: %v", block.NumberU64(), err)
		}
		// End Quorum

		database.TrieDB().Reference(root, common.Hash{})
		if parent != (common.Hash{}) {
			database.TrieDB().Dereference(parent)
		}
		parent = root
	}
	if report {
		nodes, imgs := database.TrieDB().Size()
		log.Info("Historical state regenerated", "block", current.NumberU64(), "elapsed", time.Since(start), "nodes", nodes, "preimages", imgs)
	}
	return statedb, privateStateDB, nil
}

// stateAtTransaction returns the execution environment of a certain transaction.
func (eth *Ethereum[T,P]) stateAtTransaction(ctx context.Context, block *types.Block[P], txIndex int, reexec uint64) (core.Message, vm.BlockContext, *state.StateDB[P], *state.StateDB[P], mps.PrivateStateRepository[P], error) {
	// Short circuit if it's genesis block.
	if block.NumberU64() == 0 {
		return nil, vm.BlockContext{}, nil, nil, nil, errors.New("no transaction in genesis")
	}
	// Create the parent state database
	parent := eth.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.BlockContext{}, nil, nil, nil, fmt.Errorf("parent %#x not found", block.ParentHash())
	}
	// Quorum
	statedb, privateStateRepo, err := eth.stateAtBlock(parent, reexec, nil, true)
	if err != nil {
		return nil, vm.BlockContext{}, nil, nil, nil, err
	}
	psm, err := eth.blockchain.PrivateStateManager().ResolveForUserContext(ctx)
	if err != nil {
		return nil, vm.BlockContext{}, nil, nil, nil, err
	}
	privateStateDb, err := privateStateRepo.StatePSI(psm.ID)
	if err != nil {
		return nil, vm.BlockContext{}, nil, nil, nil, err
	}
	// End Quorum
	if txIndex == 0 && len(block.Transactions()) == 0 {
		return nil, vm.BlockContext{}, statedb, privateStateDb, privateStateRepo, nil
	}
	// Recompute transactions up to the target index.
	signer := types.MakeSigner[P](eth.blockchain.Config(), block.Number())
	for idx, tx := range block.Transactions() {
		// Quorum
		privateStateDbToUse := core.PrivateStateDBForTxn(eth.blockchain.Config().IsQuorum, tx, statedb, privateStateDb)
		// End Quorum
		// Assemble the transaction call message and return if the requested offset
		msg, _ := tx.AsMessage(signer)
		msg = eth.clearMessageDataIfNonParty(msg, psm) // Quorum
		txContext := core.NewEVMTxContext(msg)
		context := core.NewEVMBlockContext[P](block.Header(), eth.blockchain, nil)
		if idx == txIndex {
			return msg, context, statedb, privateStateDb, privateStateRepo, nil
		}
		// Not yet the searched for transaction, execute on top of the current state
		vmenv := vm.NewEVM(context, txContext, statedb, privateStateDbToUse, eth.blockchain.Config(), vm.Config[P]{})
		vmenv.SetCurrentTX(tx)
		vmenv.InnerApply = func(innerTx *types.Transaction[P]) error {
			return applyInnerTransaction(eth.blockchain, statedb, privateStateDbToUse, block.Header(), tx, vm.Config[P]{}, privateStateRepo.IsMPS(), privateStateRepo, vmenv, innerTx, idx)
		}
		if _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(tx.Gas())); err != nil {
			return nil, vm.BlockContext{}, nil, nil, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}
		// Ensure any modifications are committed to the state
		// Only delete empty objects if EIP158/161 (a.k.a Spurious Dragon) is in effect
		statedb.Finalise(vmenv.ChainConfig().IsEIP158(block.Number()))
	}
	return nil, vm.BlockContext{}, nil, nil, nil, fmt.Errorf("transaction index %d out of range for block %#x", txIndex, block.Hash())
}

// Quorum

func (eth *Ethereum[T,P]) GetBlockchain() *core.BlockChain[P] {
	return eth.BlockChain()
}

func applyInnerTransaction[P crypto.PublicKey](bc *core.BlockChain[P], stateDB *state.StateDB[P], privateStateDB *state.StateDB[P], header *types.Header[P], outerTx *types.Transaction[P], evmConf vm.Config[P], forceNonParty bool, privateStateRepo mps.PrivateStateRepository[P], vmenv *vm.EVM[P], innerTx *types.Transaction[P], txIndex int) error {
	var (
		author  *common.Address = nil // ApplyTransaction will determine the author from the header so we won't do it here
		gp      *core.GasPool   = new(core.GasPool).AddGas(outerTx.Gas())
		usedGas uint64          = 0
	)
	return core.ApplyInnerTransaction[P](bc, author, gp, stateDB, privateStateDB, header, outerTx, &usedGas, evmConf, forceNonParty, privateStateRepo, vmenv, innerTx, txIndex)
}

// clearMessageDataIfNonParty sets the message data to empty hash in case the private state is not party to the
// transaction. The effect is that when the private tx payload is resolved using the privacy manager the private part of
// the transaction is not retrieved and the transaction is being executed as if the node/private state is not party to
// the transaction.
func (eth *Ethereum[T,P]) clearMessageDataIfNonParty(msg types.Message, psm *mps.PrivateStateMetadata) types.Message {
	if msg.IsPrivate() {
		_, managedParties, _, _, _ := private.Ptm.Receive(common.BytesToEncryptedPayloadHash(msg.Data()))

		if eth.GetBlockchain().PrivateStateManager().NotIncludeAny(psm, managedParties...) {
			return msg.WithEmptyPrivateData(true)
		}
	}
	return msg
}
