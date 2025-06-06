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

package state

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/metrics"
	"github.com/pavelkrolevets/MIR-pro/rlp"
)

// var emptyCodeHash = crypto.Keccak256(nil)

type Code []byte

func (c Code) String() string {
	return string(c) //strings.Join(Disassemble(c), " ")
}

type Storage map[common.Hash]common.Hash

func (s Storage) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

func (s Storage) Copy() Storage {
	cpy := make(Storage)
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}

// stateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call CommitTrie to write the modified storage trie into a database.
type stateObject [P crypto.PublicKey] struct {
	address  common.Address
	addrHash common.Hash // hash of ethereum address of the account
	data     Account
	db       *StateDB[P]

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	trie Trie // storage trie, which becomes non-nil on first access
	code Code // contract bytecode, which gets set when code is loaded

	// Quorum
	// contains extra data that is linked to the account
	accountExtraData *AccountExtraData
	// as there are many fields in accountExtraData which might be concurrently changed
	// this is to make sure we can keep track of changes individually.
	accountExtraDataMutex sync.Mutex

	originStorage  Storage // Storage cache of original entries to dedup rewrites, reset for every transaction
	pendingStorage Storage // Storage entries that need to be flushed to disk, at the end of an entire block
	dirtyStorage   Storage // Storage entries that have been modified in the current transaction execution
	fakeStorage    Storage // Fake storage which constructed by caller for debugging purpose.

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
	// Quorum
	// flag to track changes in AccountExtraData
	dirtyAccountExtraData bool

	mux sync.Mutex
}

// empty returns whether the account is considered empty.
func (s *stateObject[P]) empty() bool {
	return s.data.Nonce == 0 && s.data.Balance.Sign() == 0 && bytes.Equal(s.data.CodeHash, crypto.Keccak256[P](nil))
}

// Account is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
}

// newObject creates a state object.
func newObject[P crypto.PublicKey] (db *StateDB[P], address common.Address, data Account) *stateObject[P] {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}
	if data.CodeHash == nil {
		var emptyCodeHash = crypto.Keccak256[P](nil)
		data.CodeHash = emptyCodeHash
	}
	if data.Root == (common.Hash{}) {
		data.Root = emptyRoot
	}
	return &stateObject[P]{
		db:             db,
		address:        address,
		addrHash:       crypto.Keccak256Hash[P](address[:]),
		data:           data,
		originStorage:  make(Storage),
		pendingStorage: make(Storage),
		dirtyStorage:   make(Storage),
	}
}

// EncodeRLP implements rlp.Encoder.
func (s *stateObject[P]) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, s.data)
}

// setError remembers the first non-nil error it is called with.
func (s *stateObject[P]) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *stateObject[P]) markSuicided() {
	s.suicided = true
}

func (s *stateObject[P]) touch() {
	s.db.journal.append(touchChange[P]{
		account: &s.address,
	})
	if s.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		s.db.journal.dirty(s.address)
	}
}

func (s *stateObject[P]) getTrie(db Database) Trie {
	if s.trie == nil {
		// Try fetching from prefetcher first
		// We don't prefetch empty tries
		if s.data.Root != emptyRoot && s.db.prefetcher != nil {
			// When the miner is creating the pending state, there is no
			// prefetcher
			s.trie = s.db.prefetcher.trie(s.data.Root)
		}
		if s.trie == nil {
			var err error
			s.trie, err = db.OpenStorageTrie(s.addrHash, s.data.Root)
			if err != nil {
				s.trie, _ = db.OpenStorageTrie(s.addrHash, common.Hash{})
				s.setError(fmt.Errorf("can't create storage trie: %v", err))
			}
		}
	}
	return s.trie
}

func (so *stateObject[P]) storageRoot(db Database) common.Hash {
	return so.getTrie(db).Hash()
}

// GetState retrieves a value from the account storage trie.
func (s *stateObject[P]) GetState(db Database, key common.Hash) common.Hash {
	// If the fake storage is set, only lookup the state here(in the debugging mode)
	if s.fakeStorage != nil {
		return s.fakeStorage[key]
	}
	// If we have a dirty value for this state entry, return it
	value, dirty := s.dirtyStorage[key]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return s.GetCommittedState(db, key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
func (s *stateObject[P]) GetCommittedState(db Database, key common.Hash) common.Hash {
	// If the fake storage is set, only lookup the state here(in the debugging mode)
	if s.fakeStorage != nil {
		return s.fakeStorage[key]
	}
	// If we have a pending write or clean cached, return that
	if value, pending := s.pendingStorage[key]; pending {
		return value
	}
	if value, cached := s.originStorage[key]; cached {
		return value
	}
	// If no live objects are available, attempt to use snapshots
	var (
		enc   []byte
		err   error
		meter *time.Duration
	)
	readStart := time.Now()
	if metrics.EnabledExpensive {
		// If the snap is 'under construction', the first lookup may fail. If that
		// happens, we don't want to double-count the time elapsed. Thus this
		// dance with the metering.
		defer func() {
			if meter != nil {
				*meter += time.Since(readStart)
			}
		}()
	}
	if s.db.snap != nil {
		if metrics.EnabledExpensive {
			meter = &s.db.SnapshotStorageReads
		}
		// If the object was destructed in *this* block (and potentially resurrected),
		// the storage has been cleared out, and we should *not* consult the previous
		// snapshot about any storage values. The only possible alternatives are:
		//   1) resurrect happened, and new slot values were set -- those should
		//      have been handles via pendingStorage above.
		//   2) we don't have new values, and can deliver empty response back
		if _, destructed := s.db.snapDestructs[s.addrHash]; destructed {
			return common.Hash{}
		}
		enc, err = s.db.snap.Storage(s.addrHash, crypto.Keccak256Hash[P](key.Bytes()))
	}
	// If snapshot unavailable or reading from it failed, load from the database
	if s.db.snap == nil || err != nil {
		if meter != nil {
			// If we already spent time checking the snapshot, account for it
			// and reset the readStart
			*meter += time.Since(readStart)
			readStart = time.Now()
		}
		if metrics.EnabledExpensive {
			meter = &s.db.StorageReads
		}
		if enc, err = s.getTrie(db).TryGet(key.Bytes()); err != nil {
			s.setError(err)
			return common.Hash{}
		}
	}
	var value common.Hash
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			s.setError(err)
		}
		value.SetBytes(content)
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	s.originStorage[key] = value
	return value
}

// SetState updates a value in account storage.
func (s *stateObject[P]) SetState(db Database, key, value common.Hash) {
	// If the fake storage is set, put the temporary state update here.
	if s.fakeStorage != nil {
		s.fakeStorage[key] = value
		return
	}
	// If the new value is the same as old, don't set
	prev := s.GetState(db, key)
	if prev == value {
		return
	}
	// New value is different, update and journal the change
	s.db.journal.append(storageChange[P]{
		account:  &s.address,
		key:      key,
		prevalue: prev,
	})
	s.setState(key, value)
}

// SetStorage replaces the entire state storage with the given one.
//
// After this function is called, all original state will be ignored and state
// lookup only happens in the fake state storage.
//
// Note this function should only be used for debugging purpose.
func (s *stateObject[P]) SetStorage(storage map[common.Hash]common.Hash) {
	// Allocate fake storage if it's nil.
	if s.fakeStorage == nil {
		s.fakeStorage = make(Storage)
	}
	for key, value := range storage {
		s.fakeStorage[key] = value
	}
	// Don't bother journal since this function should only be used for
	// debugging and the `fake` storage won't be committed to database.
}

func (s *stateObject[P]) setState(key, value common.Hash) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.dirtyStorage[key] = value
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (s *stateObject[P]) finalise(prefetch bool) {
	s.mux.Lock()
	defer s.mux.Unlock()

	slotsToPrefetch := make([][]byte, 0, len(s.dirtyStorage))

	for key, value := range s.dirtyStorage {
		s.pendingStorage[key] = value
		if value != s.originStorage[key] {
			slotsToPrefetch = append(slotsToPrefetch, common.CopyBytes(key[:])) // Copy needed for closure
		}
	}
	if s.db.prefetcher != nil && prefetch && len(slotsToPrefetch) > 0 && s.data.Root != emptyRoot {
		s.db.prefetcher.prefetch(s.data.Root, slotsToPrefetch)
	}
	if len(s.dirtyStorage) > 0 {
		s.dirtyStorage = make(Storage)
	}
}

// updateTrie writes cached storage modifications into the object's storage trie.
// It will return nil if the trie has not been loaded and no changes have been made
func (s *stateObject[P]) updateTrie(db Database) Trie {
	// Make sure all dirty slots are finalized into the pending storage area
	s.finalise(false) // Don't prefetch any more, pull directly if need be
	if len(s.pendingStorage) == 0 {
		return s.trie
	}
	// Track the amount of time wasted on updating the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageUpdates += time.Since(start) }(time.Now())
	}
	// The snapshot storage map for the object
	var storage map[common.Hash][]byte
	// Insert all the pending updates into the trie
	tr := s.getTrie(db)
	hasher := s.db.hasher

	s.mux.Lock()
	defer s.mux.Unlock()
	usedStorage := make([][]byte, 0, len(s.pendingStorage))
	for key, value := range s.pendingStorage {
		// Skip noop changes, persist actual changes
		if value == s.originStorage[key] {
			continue
		}
		s.originStorage[key] = value

		var v []byte
		if (value == common.Hash{}) {
			s.setError(tr.TryDelete(key[:]))
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ = rlp.EncodeToBytes(common.TrimLeftZeroes(value[:]))
			s.setError(tr.TryUpdate(key[:], v))
		}
		// If state snapshotting is active, cache the data til commit
		if s.db.snap != nil {
			if storage == nil {
				// Retrieve the old storage map, if available, create a new one otherwise
				if storage = s.db.snapStorage[s.addrHash]; storage == nil {
					storage = make(map[common.Hash][]byte)
					s.db.snapStorage[s.addrHash] = storage
				}
			}
			storage[crypto.HashData[P](hasher, key[:])] = v // v will be nil if value is 0x00
		}
		usedStorage = append(usedStorage, common.CopyBytes(key[:])) // Copy needed for closure
	}
	if s.db.prefetcher != nil {
		s.db.prefetcher.used(s.data.Root, usedStorage)
	}
	if len(s.pendingStorage) > 0 {
		s.pendingStorage = make(Storage)
	}
	return tr
}

// UpdateRoot sets the trie root to the current root hash of
func (s *stateObject[P]) updateRoot(db Database) {
	// If nothing changed, don't bother with hashing anything
	if s.updateTrie(db) == nil {
		return
	}
	// Track the amount of time wasted on hashing the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageHashes += time.Since(start) }(time.Now())
	}
	s.data.Root = s.trie.Hash()
}

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (s *stateObject[P]) CommitTrie(db Database) error {
	// If nothing changed, don't bother with hashing anything
	if s.updateTrie(db) == nil {
		return nil
	}
	if s.dbErr != nil {
		return s.dbErr
	}
	// Track the amount of time wasted on committing the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageCommits += time.Since(start) }(time.Now())
	}
	root, err := s.trie.Commit(nil)
	if err == nil {
		s.data.Root = root
	}
	return err
}

// AddBalance adds amount to s's balance.
// It is used to add funds to the destination account of a transfer.
func (s *stateObject[P]) AddBalance(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if s.empty() {
			s.touch()
		}
		return
	}
	s.SetBalance(new(big.Int).Add(s.Balance(), amount))
}

// SubBalance removes amount from s's balance.
// It is used to remove funds from the origin account of a transfer.
func (s *stateObject[P]) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalance(new(big.Int).Sub(s.Balance(), amount))
}

func (s *stateObject[P]) SetBalance(amount *big.Int) {
	s.db.journal.append(balanceChange[P]{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Balance),
	})
	s.setBalance(amount)
}

func (s *stateObject[P]) setBalance(amount *big.Int) {
	s.data.Balance = amount
}

// Return the gas back to the origin. Used by the Virtual machine or Closures
func (s *stateObject[P]) ReturnGas(gas *big.Int) {}

func (s *stateObject[P]) deepCopy(db *StateDB[P]) *stateObject[P] {
	s.mux.Lock()
	defer s.mux.Unlock()

	stateObject := newObject[P](db, s.address, s.data)
	if s.trie != nil {
		stateObject.trie = db.db.CopyTrie(s.trie)
	}
	stateObject.code = s.code
	stateObject.dirtyStorage = s.dirtyStorage.Copy()
	stateObject.originStorage = s.originStorage.Copy()
	stateObject.pendingStorage = s.pendingStorage.Copy()
	stateObject.suicided = s.suicided
	stateObject.dirtyCode = s.dirtyCode
	stateObject.deleted = s.deleted
	// Quorum - copy AccountExtraData
	stateObject.accountExtraData = s.accountExtraData
	stateObject.dirtyAccountExtraData = s.dirtyAccountExtraData

	return stateObject
}

//
// Attribute accessors
//

// Returns the address of the contract/account
func (s *stateObject[P]) Address() common.Address {
	return s.address
}

// Code returns the contract code associated with this object, if any.
func (s *stateObject[P]) Code(db Database) []byte {
	if s.code != nil {
		return s.code
	}
	var emptyCodeHash = crypto.Keccak256[P](nil)
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(s.addrHash, common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code hash %x: %v", s.CodeHash(), err))
	}
	s.code = code
	return code
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (s *stateObject[P]) CodeSize(db Database) int {
	if s.code != nil {
		return len(s.code)
	}
	var emptyCodeHash = crypto.Keccak256[P](nil)
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return 0
	}
	size, err := db.ContractCodeSize(s.addrHash, common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.setError(fmt.Errorf("can't load code size %x: %v", s.CodeHash(), err))
	}
	return size
}

func (s *stateObject[P]) SetCode(codeHash common.Hash, code []byte) {
	prevcode := s.Code(s.db.db)
	s.db.journal.append(codeChange[P]{
		account:  &s.address,
		prevhash: s.CodeHash(),
		prevcode: prevcode,
	})
	s.setCode(codeHash, code)
}

func (s *stateObject[P]) setCode(codeHash common.Hash, code []byte) {
	s.code = code
	s.data.CodeHash = codeHash[:]
	s.dirtyCode = true
}

func (s *stateObject[P]) SetNonce(nonce uint64) {
	s.db.journal.append(nonceChange[P]{
		account: &s.address,
		prev:    s.data.Nonce,
	})
	s.setNonce(nonce)
}

func (s *stateObject[P]) setNonce(nonce uint64) {
	s.data.Nonce = nonce
}

// Quorum
// SetAccountExtraData modifies the AccountExtraData reference and journals it
func (s *stateObject[P]) SetAccountExtraData(extraData *AccountExtraData) {
	current, _ := s.AccountExtraData()
	s.db.journal.append(accountExtraDataChange[P]{
		account: &s.address,
		prev:    current,
	})
	s.setAccountExtraData(extraData)
}

// A new AccountExtraData will be created if not exists.
// This must be called after successfully acquiring accountExtraDataMutex lock
func (s *stateObject[P]) journalAccountExtraData() *AccountExtraData {
	current, _ := s.AccountExtraData()
	s.db.journal.append(accountExtraDataChange[P]{
		account: &s.address,
		prev:    current.copy(),
	})
	if current == nil {
		current = &AccountExtraData{}
	}
	return current
}

// Quorum
// SetStatePrivacyMetadata updates the PrivacyMetadata in AccountExtraData and journals it.
func (s *stateObject[P]) SetStatePrivacyMetadata(pm *PrivacyMetadata) {
	s.accountExtraDataMutex.Lock()
	defer s.accountExtraDataMutex.Unlock()

	newExtraData := s.journalAccountExtraData()
	newExtraData.PrivacyMetadata = pm
	s.setAccountExtraData(newExtraData)
}

// Quorum
// SetStatePrivacyMetadata updates the PrivacyMetadata in AccountExtraData and journals it.
func (s *stateObject[P]) SetManagedParties(managedParties []string) {
	s.accountExtraDataMutex.Lock()
	defer s.accountExtraDataMutex.Unlock()

	newExtraData := s.journalAccountExtraData()
	newExtraData.ManagedParties = managedParties
	s.setAccountExtraData(newExtraData)
}

// Quorum
// setAccountExtraData modifies the AccountExtraData reference in this state object
func (s *stateObject[P]) setAccountExtraData(extraData *AccountExtraData) {
	s.accountExtraData = extraData
	s.dirtyAccountExtraData = true
}

func (s *stateObject[P]) CodeHash() []byte {
	return s.data.CodeHash
}

func (s *stateObject[P]) Balance() *big.Int {
	return s.data.Balance
}

func (s *stateObject[P]) Nonce() uint64 {
	return s.data.Nonce
}

// Quorum
// AccountExtraData returns the extra data in this state object.
// It will also update the reference by searching the accountExtraDataTrie.
//
// This method enforces on returning error and never returns (nil, nil).
func (s *stateObject[P]) AccountExtraData() (*AccountExtraData, error) {
	if s.accountExtraData != nil {
		return s.accountExtraData, nil
	}
	val, err := s.getCommittedAccountExtraData()
	if err != nil {
		return nil, err
	}
	s.accountExtraData = val
	return val, nil
}

// Quorum
// getCommittedAccountExtraData looks for an entry in accountExtraDataTrie.
//
// This method enforces on returning error and never returns (nil, nil).
func (s *stateObject[P]) getCommittedAccountExtraData() (*AccountExtraData, error) {
	val, err := s.db.accountExtraDataTrie.TryGet(s.address.Bytes())
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from the accountExtraDataTrie. Cause: %v", err)
	}
	if len(val) == 0 {
		return nil, fmt.Errorf("%s: %w", s.address.Hex(), common.ErrNoAccountExtraData)
	}
	var extraData AccountExtraData
	if err := rlp.DecodeBytes(val, &extraData); err != nil {
		return nil, fmt.Errorf("unable to decode to AccountExtraData. Cause: %v", err)
	}
	return &extraData, nil
}

// Quorum - Privacy Enhancements
// PrivacyMetadata returns the reference to PrivacyMetadata.
// It will returrn an error if no PrivacyMetadata is in the AccountExtraData.
func (s *stateObject[P]) PrivacyMetadata() (*PrivacyMetadata, error) {
	extraData, err := s.AccountExtraData()
	if err != nil {
		return nil, err
	}
	// extraData can't be nil. Refer to s.AccountExtraData()
	if extraData.PrivacyMetadata == nil {
		return nil, fmt.Errorf("no privacy metadata data for contract %s", s.address.Hex())
	}
	return extraData.PrivacyMetadata, nil
}

func (s *stateObject[P]) GetCommittedPrivacyMetadata() (*PrivacyMetadata, error) {
	extraData, err := s.getCommittedAccountExtraData()
	if err != nil {
		return nil, err
	}
	if extraData == nil || extraData.PrivacyMetadata == nil {
		return nil, fmt.Errorf("The provided contract does not have privacy metadata: %x", s.address)
	}
	return extraData.PrivacyMetadata, nil
}

// End Quorum - Privacy Enhancements

// ManagedParties will return empty if no account extra data found
func (s *stateObject[P]) ManagedParties() ([]string, error) {
	extraData, err := s.AccountExtraData()
	if errors.Is(err, common.ErrNoAccountExtraData) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	// extraData can't be nil. Refer to s.AccountExtraData()
	return extraData.ManagedParties, nil
}

// Never called, but must be present to allow stateObject to be used
// as a vm.Account interface that also satisfies the vm.ContractRef
// interface. Interfaces are awesome.
func (s *stateObject[P]) Value() *big.Int {
	panic("Value on stateObject should never be called")
}
