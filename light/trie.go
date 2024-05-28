// Copyright 2015 The go-ethereum Authors
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

package light

import (
	"context"
	"errors"
	"fmt"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/core/rawdb"
	"github.com/MIRChain/MIR/core/state"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/ethdb"
	"github.com/MIRChain/MIR/trie"
)

// var (
// 	sha3Nil = crypto.Keccak256Hash(nil)
// )

func NewState[P crypto.PublicKey](ctx context.Context, head *types.Header[P], odr OdrBackend[P]) *state.StateDB[P] {
	state, _ := state.New[P](head.Root, NewStateDatabase(ctx, head, odr), nil)
	return state
}

func NewStateDatabase[P crypto.PublicKey](ctx context.Context, head *types.Header[P], odr OdrBackend[P]) state.Database {
	return &odrDatabase[P]{ctx, StateTrieID(head), odr}
}

type odrDatabase[P crypto.PublicKey] struct {
	ctx     context.Context
	id      *TrieID
	backend OdrBackend[P]
}

func (db *odrDatabase[P]) OpenTrie(root common.Hash) (state.Trie, error) {
	return &odrTrie[P]{db: db, id: db.id}, nil
}

func (db *odrDatabase[P]) OpenStorageTrie(addrHash, root common.Hash) (state.Trie, error) {
	return &odrTrie[P]{db: db, id: StorageTrieID(db.id, addrHash, root)}, nil
}

func (db *odrDatabase[P]) CopyTrie(t state.Trie) state.Trie {
	switch t := t.(type) {
	case *odrTrie[P]:
		cpy := &odrTrie[P]{db: t.db, id: t.id}
		if t.trie != nil {
			cpytrie := *t.trie
			cpy.trie = &cpytrie
		}
		return cpy
	default:
		panic(fmt.Errorf("unknown trie type %T", t))
	}
}

func (db *odrDatabase[P]) ContractCode(addrHash, codeHash common.Hash) ([]byte, error) {
	if codeHash == crypto.Keccak256Hash[P](nil) {
		return nil, nil
	}
	code := rawdb.ReadCode(db.backend.Database(), codeHash)
	if len(code) != 0 {
		return code, nil
	}
	id := *db.id
	id.AccKey = addrHash[:]
	req := &CodeRequest{Id: &id, Hash: codeHash}
	err := db.backend.Retrieve(db.ctx, req)
	return req.Data, err
}

func (db *odrDatabase[P]) ContractCodeSize(addrHash, codeHash common.Hash) (int, error) {
	code, err := db.ContractCode(addrHash, codeHash)
	return len(code), err
}

func (db *odrDatabase[P]) TrieDB() *trie.Database {
	return nil
}

type stubAccountExtraDataLinker struct {
}

func newAccountExtraDataLinkerStub() rawdb.AccountExtraDataLinker {
	return &stubAccountExtraDataLinker{}
}

func (pml *stubAccountExtraDataLinker) GetAccountExtraDataRoot(_ common.Hash) common.Hash {
	return common.Hash{}
}

func (pml *stubAccountExtraDataLinker) Link(_, _ common.Hash) error {
	return nil
}

func (db *odrDatabase[P]) AccountExtraDataLinker() rawdb.AccountExtraDataLinker {
	return newAccountExtraDataLinkerStub()
}

type odrTrie[P crypto.PublicKey] struct {
	db   *odrDatabase[P]
	id   *TrieID
	trie *trie.Trie[P]
}

func (t *odrTrie[P]) TryGet(key []byte) ([]byte, error) {
	key = crypto.Keccak256[P](key)
	var res []byte
	err := t.do(key, func() (err error) {
		res, err = t.trie.TryGet(key)
		return err
	})
	return res, err
}

func (t *odrTrie[P]) TryUpdate(key, value []byte) error {
	key = crypto.Keccak256[P](key)
	return t.do(key, func() error {
		return t.trie.TryUpdate(key, value)
	})
}

func (t *odrTrie[P]) TryDelete(key []byte) error {
	key = crypto.Keccak256[P](key)
	return t.do(key, func() error {
		return t.trie.TryDelete(key)
	})
}

func (t *odrTrie[P]) Commit(onleaf trie.LeafCallback) (common.Hash, error) {
	if t.trie == nil {
		return t.id.Root, nil
	}
	return t.trie.Commit(onleaf)
}

func (t *odrTrie[P]) Hash() common.Hash {
	if t.trie == nil {
		return t.id.Root
	}
	return t.trie.Hash()
}

func (t *odrTrie[P]) NodeIterator(startkey []byte) trie.NodeIterator {
	return newNodeIterator(t, startkey)
}

func (t *odrTrie[P]) GetKey(sha []byte) []byte {
	return nil
}

func (t *odrTrie[P]) Prove(key []byte, fromLevel uint, proofDb ethdb.KeyValueWriter) error {
	return errors.New("not implemented, needs client/server interface split")
}

// do tries and retries to execute a function until it returns with no error or
// an error type other than MissingNodeError
func (t *odrTrie[P]) do(key []byte, fn func() error) error {
	for {
		var err error
		if t.trie == nil {
			t.trie, err = trie.New[P](t.id.Root, trie.NewDatabase(t.db.backend.Database()))
		}
		if err == nil {
			err = fn()
		}
		if _, ok := err.(*trie.MissingNodeError); !ok {
			return err
		}
		r := &TrieRequest[P]{Id: t.id, Key: key}
		if err := t.db.backend.Retrieve(t.db.ctx, r); err != nil {
			return err
		}
	}
}

type nodeIterator[P crypto.PublicKey] struct {
	trie.NodeIterator
	t   *odrTrie[P]
	err error
}

func newNodeIterator[P crypto.PublicKey](t *odrTrie[P], startkey []byte) trie.NodeIterator {
	it := &nodeIterator[P]{t: t}
	// Open the actual non-ODR trie if that hasn't happened yet.
	if t.trie == nil {
		it.do(func() error {
			t, err := trie.New[P](t.id.Root, trie.NewDatabase(t.db.backend.Database()))
			if err == nil {
				it.t.trie = t
			}
			return err
		})
	}
	it.do(func() error {
		it.NodeIterator = it.t.trie.NodeIterator(startkey)
		return it.NodeIterator.Error()
	})
	return it
}

func (it *nodeIterator[P]) Next(descend bool) bool {
	var ok bool
	it.do(func() error {
		ok = it.NodeIterator.Next(descend)
		return it.NodeIterator.Error()
	})
	return ok
}

// do runs fn and attempts to fill in missing nodes by retrieving.
func (it *nodeIterator[P]) do(fn func() error) {
	var lasthash common.Hash
	for {
		it.err = fn()
		missing, ok := it.err.(*trie.MissingNodeError)
		if !ok {
			return
		}
		if missing.NodeHash == lasthash {
			it.err = fmt.Errorf("retrieve loop for trie node %x", missing.NodeHash)
			return
		}
		lasthash = missing.NodeHash
		r := &TrieRequest[P]{Id: it.t.id, Key: nibblesToKey(missing.Path)}
		if it.err = it.t.db.backend.Retrieve(it.t.db.ctx, r); it.err != nil {
			return
		}
	}
}

func (it *nodeIterator[P]) Error() error {
	if it.err != nil {
		return it.err
	}
	return it.NodeIterator.Error()
}

func nibblesToKey(nib []byte) []byte {
	if len(nib) > 0 && nib[len(nib)-1] == 0x10 {
		nib = nib[:len(nib)-1] // drop terminator
	}
	if len(nib)&1 == 1 {
		nib = append(nib, 0) // make even
	}
	key := make([]byte, len(nib)/2)
	for bi, ni := 0, 0; ni < len(nib); bi, ni = bi+1, ni+2 {
		key[bi] = nib[ni]<<4 | nib[ni+1]
	}
	return key
}
