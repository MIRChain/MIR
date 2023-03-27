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
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// NewTxsEvent is posted when a batch of transactions enter the transaction pool.
type NewTxsEvent [P crypto.PublicKey]struct{ Txs []*types.Transaction[P] }

// PendingStateEvent is posted pre mining and notifies of pending state changes.
type PendingStateEvent struct{}

// NewMinedBlockEvent is posted when a block has been imported.
type NewMinedBlockEvent [P crypto.PublicKey] struct{ Block *types.Block[P] }

// RemovedLogsEvent is posted when a reorg happens
type RemovedLogsEvent [P crypto.PublicKey] struct{ Logs []*types.Log }

type ChainEvent [P crypto.PublicKey] struct {
	Block *types.Block[P]
	Hash  common.Hash
	Logs  []*types.Log
}

type ChainSideEvent [P crypto.PublicKey] struct {
	Block *types.Block[P]
}

type ChainHeadEvent [P crypto.PublicKey] struct{ Block *types.Block[P] }
