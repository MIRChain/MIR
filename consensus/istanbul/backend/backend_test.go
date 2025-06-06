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

package backend

import (
	"bytes"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul"
	istanbulcommon "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul/validator"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
)

func TestSign(t *testing.T) {
	b := newBackend()
	defer b.Stop()
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	//Check signature recover
	hashData := crypto.Keccak256[nist.PublicKey](data)
	pubkey, _ := crypto.Ecrecover[nist.PublicKey](hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256[nist.PublicKey](pubkey[1:])[12:])
	if signer != getAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
	}
}

func TestCheckSignature(t *testing.T) {
	key, _ := generatePrivateKey()
	data := []byte("Here is a string....")
	hashData := crypto.Keccak256[nist.PublicKey](data)
	sig, _ := crypto.Sign[nist.PrivateKey](hashData, key)
	b := newBackend()
	defer b.Stop()
	a := getAddress()
	err := b.CheckSignature(data, a, sig)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	a = getInvalidAddress()
	err = b.CheckSignature(data, a, sig)
	if err != istanbulcommon.ErrInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, istanbulcommon.ErrInvalidSignature)
	}
}

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256[nist.PublicKey](data)
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		// CheckValidatorSignature should succeed
		addr, err := istanbul.CheckValidatorSignature[nist.PublicKey](vset, data, sig)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		validator := vset.GetByIndex(uint64(i))
		if addr != validator.Address() {
			t.Errorf("validator address mismatch: have %v, want %v", addr, validator.Address())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey[nist.PrivateKey]()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := istanbul.CheckValidatorSignature[nist.PublicKey](vset, data, sig)
	if err != istanbul.ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, istanbul.ErrUnauthorizedAddress)
	}
	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCommit(t *testing.T) {
	backend := newBackend()
	defer backend.Stop()

	commitCh := make(chan *types.Block[nist.PublicKey])
	// Case: it's a proposer, so the backend.commit will receive channel result from backend.Commit function
	testCases := []struct {
		expectedErr       error
		expectedSignature [][]byte
		expectedBlock     func() *types.Block[nist.PublicKey]
	}{
		{
			// normal case
			nil,
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-1)...)},
			func() *types.Block[nist.PublicKey] {
				chain, engine := NewBlockChain(1, big.NewInt(0))
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				return updateQBFTBlock(block, engine.Address())
			},
		},
		{
			// invalid signature
			istanbulcommon.ErrInvalidCommittedSeals,
			nil,
			func() *types.Block[nist.PublicKey] {
				chain, engine := NewBlockChain(1, big.NewInt(0))
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				return updateQBFTBlock(block, engine.Address())
			},
		},
	}

	for _, test := range testCases {
		expBlock := test.expectedBlock()
		go func() {
			result := <-backend.commitCh
			commitCh <- result
		}()

		backend.proposedBlockHash = expBlock.Hash()
		if err := backend.Commit(expBlock, test.expectedSignature, big.NewInt(0)); err != nil {
			if err != test.expectedErr {
				t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
			}
		}

		if test.expectedErr == nil {
			// to avoid race condition is occurred by goroutine
			select {
			case result := <-commitCh:
				if result.Hash() != expBlock.Hash() {
					t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timeout")
			}
		}
	}
}

func TestGetProposer(t *testing.T) {
	chain, engine := NewBlockChain(1, big.NewInt(0))
	defer engine.Stop()
	block := makeBlock(chain, engine, chain.Genesis())
	chain.InsertChain(types.Blocks[nist.PublicKey]{block})
	expected := engine.GetProposer(1)
	actual := engine.Address()
	if actual != expected {
		t.Errorf("proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
	}
}

// TestQBFTTransitionDeadlock test whether a deadlock occurs when testQBFTBlock is set to 1
// This was fixed as part of commit 2a8310663ecafc0233758ca7883676bf568e926e
func TestQBFTTransitionDeadlock(t *testing.T) {
	timeout := time.After(1 * time.Minute)
	done := make(chan bool)
	go func() {
		chain, engine := NewBlockChain(1, big.NewInt(1))
		defer engine.Stop()
		// Create an insert a new block into the chain.
		block := makeBlock(chain, engine, chain.Genesis())
		_, err := chain.InsertChain(types.Blocks[nist.PublicKey]{block})
		if err != nil {
			t.Errorf("Error inserting block: %v", err)
		}

		if err = engine.NewChainHead(); err != nil {
			t.Errorf("Error posting NewChainHead Event: %v", err)
		}

		if !engine.IsQBFTConsensus() {
			t.Errorf("IsQBFTConsensus() should return true after block insertion")
		}
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Deadlock occurred during IBFT to QBFT transition")
	case <-done:
	}
}

func TestIsQBFTConsensus(t *testing.T) {
	chain, engine := NewBlockChain(1, big.NewInt(2))
	defer engine.Stop()
	qbftConsensus := engine.IsQBFTConsensus()
	if qbftConsensus {
		t.Errorf("IsQBFTConsensus() should return false")
	}

	// Create an insert a new block into the chain.
	block := makeBlock(chain, engine, chain.Genesis())
	_, err := chain.InsertChain(types.Blocks[nist.PublicKey]{block})
	if err != nil {
		t.Errorf("Error inserting block: %v", err)
	}

	if err = engine.NewChainHead(); err != nil {
		t.Errorf("Error posting NewChainHead Event: %v", err)
	}

	secondBlock := makeBlock(chain, engine, block)
	_, err = chain.InsertChain(types.Blocks[nist.PublicKey]{secondBlock})
	if err != nil {
		t.Errorf("Error inserting block: %v", err)
	}

	qbftConsensus = engine.IsQBFTConsensus()
	if !qbftConsensus {
		t.Errorf("IsQBFTConsensus() should return true after block insertion")
	}
}

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d25ecdf2b0d7c0f7aae9aa624
 * Address: 0x70524d664ffe731100208a0154e556f9bb679ae6
 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (nist.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA[nist.PrivateKey](key)
}

func newTestValidatorSet(n int) (istanbul.ValidatorSet, []nist.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey[nist.PrivateKey]()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress[nist.PublicKey](*privateKey.Public())
	}
	vset := validator.NewSet(addrs, istanbul.NewRoundRobinProposerPolicy())
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

type Keys []nist.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress[nist.PublicKey](*slice[i].Public()).String(), crypto.PubkeyToAddress[nist.PublicKey](*slice[j].Public()).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func newBackend() (b *Backend[nist.PrivateKey,nist.PublicKey]) {
	_, b = NewBlockChain(1, big.NewInt(0))
	key, _ := generatePrivateKey()
	b.privateKey = key
	return
}
