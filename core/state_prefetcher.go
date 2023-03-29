// Copyright 2019 The go-ethereum Authors
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
	"sync"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/params"

	// Quorum
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/private"
)

// statePrefetcher is a basic Prefetcher, which blindly executes a block on top
// of an arbitrary state with the goal of prefetching potentially useful state
// data from disk before the main block processor start executing.
type statePrefetcher [P crypto.PublicKey] struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain[P]         // Canonical block chain
	engine consensus.Engine[P]    // Consensus engine used for block rewards

	pend sync.WaitGroup // Quorum: wait for MPS prefetching
}

// newStatePrefetcher initialises a new statePrefetcher.
func newStatePrefetcher[P crypto.PublicKey](config *params.ChainConfig, bc *BlockChain[P], engine consensus.Engine[P]) *statePrefetcher[P] {
	return &statePrefetcher[P]{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Prefetch processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb, but any changes are discarded. The
// only goal is to pre-cache transaction signatures and state trie nodes.
// Quorum: Add privateStateDb argument
func (p *statePrefetcher[P]) Prefetch(block *types.Block[P], statedb *state.StateDB, privateStateRepo mps.PrivateStateRepository[P], cfg vm.Config[P], interrupt *uint32) {
	var (
		header  = block.Header()
		gaspool = new(GasPool).AddGas(block.GasLimit())
	)
	// Iterate over and process the individual transactions
	byzantium := p.config.IsByzantium(block.Number())
	for i, tx := range block.Transactions() {
		// If block precaching was interrupted, abort
		if interrupt != nil && atomic.LoadUint32(interrupt) == 1 {
			return
		}

		// Quorum
		if tx.IsPrivate() && privateStateRepo.IsMPS() {
			p.prefetchMpsTransaction(block, tx, i, statedb.Copy(), privateStateRepo, cfg, interrupt)
		}
		privateStateDb, _ := privateStateRepo.DefaultState()
		privateStateDb.Prepare(tx.Hash(), block.Hash(), i)
		// End Quorum

		// Block precaching permitted to continue, execute the transaction
		statedb.Prepare(tx.Hash(), block.Hash(), i)

		innerApply := createInnerApply(block, tx, i, statedb, privateStateRepo, cfg, interrupt, p, privateStateDb)

		// Quorum: Add privateStateDb argument
		if err := precacheTransaction[P](p.config, p.bc, nil, gaspool, statedb, privateStateDb, header, tx, cfg, innerApply); err != nil {
			return // Ugh, something went horribly wrong, bail out
		}
		// If we're pre-byzantium, pre-load trie nodes for the intermediate root
		if !byzantium {
			statedb.IntermediateRoot(true)
		}
	}
	// If were post-byzantium, pre-load trie nodes for the final root hash
	if byzantium {
		statedb.IntermediateRoot(true)
	}
}

// precacheTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. The goal is not to execute
// the transaction successfully, rather to warm up touched data slots.
// Quorum: Add privateStateDb and isMPS arguments
func precacheTransaction[P crypto.PublicKey](config *params.ChainConfig, bc ChainContext[P], author *common.Address, gaspool *GasPool, statedb *state.StateDB, privateStateDb *state.StateDB, header *types.Header, tx *types.Transaction[P], cfg vm.Config[P], innerApply func(*types.Transaction[P]) error) error {
	// Convert the transaction into an executable message and pre-cache its sender
	msg, err := tx.AsMessage(types.MakeSigner[P](config, header.Number))
	if err != nil {
		return err
	}
	// Quorum
	// Create the EVM and execute the transaction
	context := NewEVMBlockContext(header, bc, author)
	txContext := NewEVMTxContext(msg)

	var evm *vm.EVM[P]
	// Quorum: Add privateStateDb argument
	if tx.IsPrivate() {
		evm = vm.NewEVM(context, txContext, statedb, privateStateDb, config, cfg)
	} else {
		evm = vm.NewEVM(context, txContext, statedb, statedb, config, cfg)
	}
	// End Quorum
	evm.SetCurrentTX(tx) // Quorum
	evm.InnerApply = innerApply
	// Add addresses to access list if applicable
	_, err = ApplyMessage(evm, msg, gaspool)
	return err
}

// Quorum

func (p *statePrefetcher[P]) prefetchMpsTransaction(block *types.Block[P], tx *types.Transaction[P], txIndex int, statedb *state.StateDB, privateStateRepo mps.PrivateStateRepository[P], cfg vm.Config[P], interrupt *uint32) {
	byzantium := p.config.IsByzantium(block.Number())
	// Block precaching permitted to continue, execute the transaction
	_, managedParties, _, _, err := private.Ptm.Receive(common.BytesToEncryptedPayloadHash(tx.Data()))
	if err != nil {
		return
	}
	for _, managedParty := range managedParties {
		if interrupt != nil && atomic.LoadUint32(interrupt) == 1 {
			return
		}
		psMetadata, err := p.bc.PrivateStateManager().ResolveForManagedParty(managedParty)
		if err != nil {
			continue
		}

		privateStateDb, err := privateStateRepo.StatePSI(psMetadata.ID)
		if err != nil {
			continue
		}
		p.pend.Add(1)

		innerApply := createInnerApply(block, tx, txIndex, statedb, privateStateRepo, cfg, interrupt, p, privateStateDb)

		go func(start time.Time, followup *types.Block[P], statedb *state.StateDB, privateStateDb *state.StateDB, tx *types.Transaction[P], gaspool *GasPool) {
			privateStateDb.Prepare(tx.Hash(), block.Hash(), txIndex)
			if err := precacheTransaction[P](p.config, p.bc, nil, gaspool, statedb, privateStateDb, followup.Header(), tx, cfg, innerApply); err != nil {
				return
			}
			// If we're pre-byzantium, pre-load trie nodes for the intermediate root
			if !byzantium {
				privateStateDb.IntermediateRoot(true)
			}
			p.pend.Done()
		}(time.Now(), block, statedb, privateStateDb, tx, new(GasPool).AddGas(tx.Gas())) // TODO ricardolyn: which gas: block or Tx?
	}
	p.pend.Wait()
}

func createInnerApply[P crypto.PublicKey](block *types.Block[P], tx *types.Transaction[P], txIndex int, statedb *state.StateDB, privateStateRepo mps.PrivateStateRepository[P], cfg vm.Config[P], interrupt *uint32, p *statePrefetcher[P], privateStateDb *state.StateDB) func(innerTx *types.Transaction[P]) error {
	return func(innerTx *types.Transaction[P]) error {
		if !tx.IsPrivacyMarker() {
			return nil
		} else if innerTx.IsPrivate() && privateStateRepo.IsMPS() {
			p.prefetchMpsTransaction(block, innerTx, txIndex, statedb.Copy(), privateStateRepo, cfg, interrupt)
			return nil
		} else {
			return precacheTransaction[P](p.config, p.bc, nil, new(GasPool).AddGas(innerTx.Gas()), statedb, privateStateDb, block.Header(), innerTx, cfg, nil)
		}
	}
}

// End Quorum
