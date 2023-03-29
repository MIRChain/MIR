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
	"sync"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul"
	qbfttypes "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/qbft/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// newRoundState creates a new roundState instance with the given view and validatorSet
func newRoundState[P crypto.PublicKey](view *istanbul.View, validatorSet istanbul.ValidatorSet, preprepare *qbfttypes.Preprepare[P], preparedRound *big.Int, preparedBlock istanbul.Proposal, pendingRequest *Request, hasBadProposal func(hash common.Hash) bool) *roundState[P] {
	return &roundState[P]{
		round:      view.Round,
		sequence:   view.Sequence,
		Preprepare: preprepare,
		//Prepares:       newMessageSet(validatorSet),
		//Commits:        newMessageSet(validatorSet),
		QBFTPrepares:   newQBFTMsgSet(validatorSet),
		QBFTCommits:    newQBFTMsgSet(validatorSet),
		preparedRound:  preparedRound,
		preparedBlock:  preparedBlock,
		mu:             new(sync.RWMutex),
		pendingRequest: pendingRequest,
		hasBadProposal: hasBadProposal,
		preprepareSent: big.NewInt(0),
	}
}

// roundState stores the consensus state
type roundState [P crypto.PublicKey] struct {
	round      *big.Int
	sequence   *big.Int
	Preprepare *qbfttypes.Preprepare[P]

	QBFTPrepares *qbftMsgSet
	QBFTCommits  *qbftMsgSet

	pendingRequest *Request
	preparedRound  *big.Int
	preparedBlock  istanbul.Proposal

	mu             *sync.RWMutex
	hasBadProposal func(hash common.Hash) bool

	// Keep track of preprepare sent messages
	preprepareSent *big.Int
}

func (s *roundState[P]) Subject() *Subject {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare == nil {
		return nil
	}

	return &Subject{
		View: &istanbul.View{
			Round:    new(big.Int).Set(s.round),
			Sequence: new(big.Int).Set(s.sequence),
		},
		Digest: s.Preprepare.Proposal.Hash(),
	}
}

func (s *roundState[P]) SetPreprepare(preprepare *qbfttypes.Preprepare[P]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Preprepare = preprepare
}

func (s *roundState[P]) Proposal() istanbul.Proposal {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.Preprepare != nil {
		return s.Preprepare.Proposal
	}

	return nil
}

func (s *roundState[P]) SetRound(r *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.round = new(big.Int).Set(r)
}

func (s *roundState[P]) Round() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.round
}

func (s *roundState[P]) SetSequence(seq *big.Int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence = seq
}

func (s *roundState[P]) Sequence() *big.Int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sequence
}
