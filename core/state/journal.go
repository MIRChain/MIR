// Copyright 2016 The go-ethereum Authors
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

package state

import (
	"math/big"
	"sync"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto"
)

// journalEntry is a modification entry in the state change journal that can be
// reverted on demand.
type journalEntry[P crypto.PublicKey] interface {
	// revert undoes the changes introduced by this journal entry.
	revert(*StateDB[P])

	// dirtied returns the Ethereum address modified by this journal entry.
	dirtied() *common.Address
}

// journal contains the list of state modifications applied since the last state
// commit. These are tracked to be able to be reverted in case of an execution
// exception or revertal request.
type journal[P crypto.PublicKey] struct {
	entries []journalEntry[P]      // Current changes tracked by the journal
	dirties map[common.Address]int // Dirty accounts and the number of changes
	mutex   sync.Mutex
}

// newJournal create a new initialized journal.
func newJournal[P crypto.PublicKey]() *journal[P] {
	return &journal[P]{
		dirties: make(map[common.Address]int),
	}
}

// append inserts a new modification entry to the end of the change journal.
func (j *journal[P]) append(entry journalEntry[P]) {
	defer j.mutex.Unlock()
	j.mutex.Lock()
	j.entries = append(j.entries, entry)
	if addr := entry.dirtied(); addr != nil {
		j.dirties[*addr]++
	}
}

// revert undoes a batch of journalled modifications along with any reverted
// dirty handling too.
func (j *journal[P]) revert(statedb *StateDB[P], snapshot int) {
	defer j.mutex.Unlock()
	j.mutex.Lock()
	for i := len(j.entries) - 1; i >= snapshot; i-- {
		// Undo the changes made by the operation
		j.entries[i].revert(statedb)

		// Drop any dirty tracking induced by the change
		if addr := j.entries[i].dirtied(); addr != nil {
			if j.dirties[*addr]--; j.dirties[*addr] == 0 {
				delete(j.dirties, *addr)
			}
		}
	}
	j.entries = j.entries[:snapshot]
}

// dirty explicitly sets an address to dirty, even if the change entries would
// otherwise suggest it as clean. This method is an ugly hack to handle the RIPEMD
// precompile consensus exception.
func (j *journal[P]) dirty(addr common.Address) {
	defer j.mutex.Unlock()
	j.mutex.Lock()
	j.dirties[addr]++
}

// length returns the current number of entries in the journal.
func (j *journal[P]) length() int {
	return len(j.entries)
}

type (
	// Changes to the account trie.
	createObjectChange[P crypto.PublicKey] struct {
		account *common.Address
	}
	resetObjectChange[P crypto.PublicKey] struct {
		prev         *stateObject[P]
		prevdestruct bool
	}
	suicideChange[P crypto.PublicKey] struct {
		account     *common.Address
		prev        bool // whether account had already suicided
		prevbalance *big.Int
	}

	// Changes to individual accounts.
	balanceChange[P crypto.PublicKey] struct {
		account *common.Address
		prev    *big.Int
	}
	nonceChange[P crypto.PublicKey] struct {
		account *common.Address
		prev    uint64
	}
	storageChange[P crypto.PublicKey] struct {
		account       *common.Address
		key, prevalue common.Hash
	}
	codeChange[P crypto.PublicKey] struct {
		account            *common.Address
		prevcode, prevhash []byte
	}
	// Quorum - changes to AccountExtraData
	accountExtraDataChange[P crypto.PublicKey] struct {
		account *common.Address
		prev    *AccountExtraData
	}
	// Changes to other state values.
	refundChange[P crypto.PublicKey] struct {
		prev uint64
	}
	addLogChange[P crypto.PublicKey] struct {
		txhash common.Hash
	}
	addPreimageChange[P crypto.PublicKey] struct {
		hash common.Hash
	}
	touchChange[P crypto.PublicKey] struct {
		account *common.Address
	}
	// Changes to the access list
	accessListAddAccountChange[P crypto.PublicKey] struct {
		address *common.Address
	}
	accessListAddSlotChange[P crypto.PublicKey] struct {
		address *common.Address
		slot    *common.Hash
	}
)

func (ch createObjectChange[P]) revert(s *StateDB[P]) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	delete(s.stateObjects, *ch.account)
	delete(s.stateObjectsDirty, *ch.account)
}

func (ch createObjectChange[P]) dirtied() *common.Address {
	return ch.account
}

func (ch resetObjectChange[P]) revert(s *StateDB[P]) {
	s.setStateObject(ch.prev)
	if !ch.prevdestruct && s.snap != nil {
		delete(s.snapDestructs, ch.prev.addrHash)
	}
}

func (ch resetObjectChange[P]) dirtied() *common.Address {
	return nil
}

func (ch suicideChange[P]) revert(s *StateDB[P]) {
	obj := s.getStateObject(*ch.account)
	if obj != nil {
		obj.suicided = ch.prev
		obj.setBalance(ch.prevbalance)
	}
}

func (ch suicideChange[P]) dirtied() *common.Address {
	return ch.account
}

var ripemd = common.HexToAddress("0000000000000000000000000000000000000003")

func (ch touchChange[P]) revert(s *StateDB[P]) {
}

func (ch touchChange[P]) dirtied() *common.Address {
	return ch.account
}

func (ch balanceChange[P]) revert(s *StateDB[P]) {
	s.getStateObject(*ch.account).setBalance(ch.prev)
}

func (ch balanceChange[P]) dirtied() *common.Address {
	return ch.account
}

func (ch nonceChange[P]) revert(s *StateDB[P]) {
	s.getStateObject(*ch.account).setNonce(ch.prev)
}

func (ch nonceChange[P]) dirtied() *common.Address {
	return ch.account
}

func (ch codeChange[P]) revert(s *StateDB[P]) {
	s.getStateObject(*ch.account).setCode(common.BytesToHash(ch.prevhash), ch.prevcode)
}

func (ch codeChange[P]) dirtied() *common.Address {
	return ch.account
}

// Quorum
func (ch accountExtraDataChange[P]) revert(s *StateDB[P]) {
	s.getStateObject(*ch.account).setAccountExtraData(ch.prev)
}

func (ch accountExtraDataChange[P]) dirtied() *common.Address {
	return ch.account
}

// End Quorum - Privacy Enhancements

func (ch storageChange[P]) revert(s *StateDB[P]) {
	s.getStateObject(*ch.account).setState(ch.key, ch.prevalue)
}

func (ch storageChange[P]) dirtied() *common.Address {
	return ch.account
}

func (ch refundChange[P]) revert(s *StateDB[P]) {
	s.refund = ch.prev
}

func (ch refundChange[P]) dirtied() *common.Address {
	return nil
}

func (ch addLogChange[P]) revert(s *StateDB[P]) {
	logs := s.logs[ch.txhash]
	if len(logs) == 1 {
		delete(s.logs, ch.txhash)
	} else {
		s.logs[ch.txhash] = logs[:len(logs)-1]
	}
	s.logSize--
}

func (ch addLogChange[P]) dirtied() *common.Address {
	return nil
}

func (ch addPreimageChange[P]) revert(s *StateDB[P]) {
	delete(s.preimages, ch.hash)
}

func (ch addPreimageChange[P]) dirtied() *common.Address {
	return nil
}

func (ch accessListAddAccountChange[P]) revert(s *StateDB[P]) {
	/*
		One important invariant here, is that whenever a (addr, slot) is added, if the
		addr is not already present, the add causes two journal entries:
		- one for the address,
		- one for the (address,slot)
		Therefore, when unrolling the change, we can always blindly delete the
		(addr) at this point, since no storage adds can remain when come upon
		a single (addr) change.
	*/
	s.accessList.DeleteAddress(*ch.address)
}

func (ch accessListAddAccountChange[P]) dirtied() *common.Address {
	return nil
}

func (ch accessListAddSlotChange[P]) revert(s *StateDB[P]) {
	s.accessList.DeleteSlot(*ch.address, *ch.slot)
}

func (ch accessListAddSlotChange[P]) dirtied() *common.Address {
	return nil
}
