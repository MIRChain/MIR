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

package core

import (
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/mclock"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
)

// insertStats tracks and reports on block insertion.
type insertStats [P crypto.PublicKey] struct {
	queued, processed, ignored int
	usedGas                    uint64
	lastIndex                  int
	startTime                  mclock.AbsTime
}

// statsReportLimit is the time limit during import and export after which we
// always print out progress. This avoids the user wondering what's going on.
const statsReportLimit = 8 * time.Second

// report prints statistics if some number of blocks have been processed
// or more than a few seconds have passed since the last message.
func (st *insertStats[P]) report(chain []*types.Block[P], index int, dirty common.StorageSize) {
	// Fetch the timings for the batch
	var (
		now     = mclock.Now()
		elapsed = now.Sub(st.startTime)
	)
	// If we're at the last block of the batch or report period reached, log
	if index == len(chain)-1 || elapsed >= statsReportLimit {
		// Count the number of transactions in this segment
		var txs int
		for _, block := range chain[st.lastIndex : index+1] {
			txs += len(block.Transactions())
		}
		end := chain[index]

		// Assemble the log context and send it to the logger
		context := []interface{}{
			"blocks", st.processed, "txs", txs, "mgas", float64(st.usedGas) / 1000000,
			"elapsed", common.PrettyDuration(elapsed), "mgasps", float64(st.usedGas) * 1000 / float64(elapsed),
			"number", end.Number(), "hash", end.Hash(),
		}
		if timestamp := time.Unix(int64(end.Time()), 0); time.Since(timestamp) > time.Minute {
			context = append(context, []interface{}{"age", common.PrettyAge(timestamp)}...)
		}
		context = append(context, []interface{}{"dirty", dirty}...)

		if st.queued > 0 {
			context = append(context, []interface{}{"queued", st.queued}...)
		}
		if st.ignored > 0 {
			context = append(context, []interface{}{"ignored", st.ignored}...)
		}
		log.Info("Imported new chain segment", context...)

		// Bump the stats reported to the next section
		*st = insertStats[P]{startTime: now, lastIndex: index + 1}
	}
}

// insertIterator is a helper to assist during chain import.
type insertIterator [P crypto.PublicKey] struct {
	chain types.Blocks[P] // Chain of blocks being iterated over

	results <-chan error // Verification result sink from the consensus engine
	errors  []error      // Header verification errors for the blocks

	index     int       // Current offset of the iterator
	validator Validator[P] // Validator to run if verification succeeds
}

// newInsertIterator creates a new iterator based on the given blocks, which are
// assumed to be a contiguous chain.
func newInsertIterator[P crypto.PublicKey](chain types.Blocks[P], results <-chan error, validator Validator[P]) *insertIterator[P] {
	return &insertIterator[P]{
		chain:     chain,
		results:   results,
		errors:    make([]error, 0, len(chain)),
		index:     -1,
		validator: validator,
	}
}

// next returns the next block in the iterator, along with any potential validation
// error for that block. When the end is reached, it will return (nil, nil).
func (it *insertIterator[P]) next() (*types.Block[P], error) {
	// If we reached the end of the chain, abort
	if it.index+1 >= len(it.chain) {
		it.index = len(it.chain)
		return nil, nil
	}
	// Advance the iterator and wait for verification result if not yet done
	it.index++
	if len(it.errors) <= it.index {
		it.errors = append(it.errors, <-it.results)
	}
	if it.errors[it.index] != nil {
		return it.chain[it.index], it.errors[it.index]
	}
	// Block header valid, run body validation and return
	return it.chain[it.index], it.validator.ValidateBody(it.chain[it.index])
}

// peek returns the next block in the iterator, along with any potential validation
// error for that block, but does **not** advance the iterator.
//
// Both header and body validation errors (nil too) is cached into the iterator
// to avoid duplicating work on the following next() call.
func (it *insertIterator[P]) peek() (*types.Block[P], error) {
	// If we reached the end of the chain, abort
	if it.index+1 >= len(it.chain) {
		return nil, nil
	}
	// Wait for verification result if not yet done
	if len(it.errors) <= it.index+1 {
		it.errors = append(it.errors, <-it.results)
	}
	if it.errors[it.index+1] != nil {
		return it.chain[it.index+1], it.errors[it.index+1]
	}
	// Block header valid, ignore body validation since we don't have a parent anyway
	return it.chain[it.index+1], nil
}

// previous returns the previous header that was being processed, or nil.
func (it *insertIterator[P]) previous() *types.Header[P] {
	if it.index < 1 {
		return nil
	}
	return it.chain[it.index-1].Header()
}

// first returns the first block in the it.
func (it *insertIterator[P]) first() *types.Block[P] {
	return it.chain[0]
}

// remaining returns the number of remaining blocks.
func (it *insertIterator[P]) remaining() int {
	return len(it.chain) - it.index
}

// processed returns the number of processed blocks.
func (it *insertIterator[P]) processed() int {
	return it.index + 1
}
