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

package core

import (
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common"
)

func (c *core[P]) handleFinalCommitted() error {
	logger := c.logger.New("state", c.state)
	logger.Trace("Received a final committed proposal")

	// startNewRound() needs to be called asynchronously when the transition to qbft happens
	// This is required so that the stop() on core can successfully unsubscribe from events
	nextSeq := new(big.Int).Add(c.currentView().Sequence, big.NewInt(1))
	if c.backend.IsQBFTConsensusAt(nextSeq) {
		go c.startNewRound(common.Big0)
	} else {
		c.startNewRound(common.Big0)
	}
	return nil
}
