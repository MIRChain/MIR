// Copyright 2017 The go-ethereum Authors
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

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// ChainHeaderReader defines a small collection of methods needed to access the local
// blockchain during header verification.
type ChainHeaderReader [P crypto.PublicKey] interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header[P]

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *types.Header[P]

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header[P]

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *types.Header[P]
}

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader [P crypto.PublicKey] interface {
	ChainHeaderReader[P]

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *types.Block[P]
}

// Engine is an algorithm agnostic consensus engine.
type Engine [P crypto.PublicKey] interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *types.Header[P]) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainHeaderReader[P], header *types.Header[P], seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainHeaderReader[P], headers []*types.Header[P], seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	VerifyUncles(chain ChainReader[P], block *types.Block[P]) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainHeaderReader[P], header *types.Header[P]) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// but does not assemble the block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainHeaderReader[P], header *types.Header[P], state *state.StateDB[P], txs []*types.Transaction[P],
		uncles []*types.Header[P])

	// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
	// rewards) and assembles the final block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	FinalizeAndAssemble(chain ChainHeaderReader[P], header *types.Header[P], state *state.StateDB[P], txs []*types.Transaction[P],
		uncles []*types.Header[P], receipts []*types.Receipt[P]) (*types.Block[P], error)

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	//
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	Seal(chain ChainHeaderReader[P], block *types.Block[P], results chan<- *types.Block[P], stop <-chan struct{}) error

	// SealHash returns the hash of a block prior to it being sealed.
	SealHash(header *types.Header[P]) common.Hash

	// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// that a new block should have.
	CalcDifficulty(chain ChainHeaderReader[P], time uint64, parent *types.Header[P]) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainHeaderReader[P]) []rpc.API

	// Protocol returns the protocol for this consensus
	Protocol() Protocol

	// Close terminates any background threads maintained by the consensus engine.
	Close() error
}

// Handler should be implemented is the consensus needs to handle and send peer's message
type Handler [P crypto.PublicKey] interface {
	// NewChainHead handles a new head block comes
	NewChainHead() error

	// HandleMsg handles a message from peer
	HandleMsg(address common.Address, data p2p.Msg) (bool, error)

	// SetBroadcaster sets the broadcaster to send message to peers
	SetBroadcaster(Broadcaster[P])
}

// PoW is a consensus engine based on proof-of-work.
type PoW [P crypto.PublicKey] interface {
	Engine[P]

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}

// Istanbul is a consensus engine to avoid byzantine failure
type Istanbul [P crypto.PublicKey] interface {
	Engine[P]

	// Start starts the engine
	Start(chain ChainHeaderReader[P], currentBlock func() *types.Block[P], hasBadBlock func(db ethdb.Reader, hash common.Hash) bool) error

	// Stop stops the engine
	Stop() error
}
