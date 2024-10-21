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

package types

import (
	"math/big"
	"testing"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/gost3410"
)

func TestEIP155SigningGOST(t *testing.T) {
	key, _ := crypto.GenerateKey[gost3410.PrivateKey]()
	addr := crypto.PubkeyToAddress[gost3410.PublicKey](*key.Public())

	signer := NewEIP155Signer[gost3410.PublicKey](big.NewInt(18))
	tx, err := SignTx[gost3410.PrivateKey, gost3410.PublicKey](NewTransaction[gost3410.PublicKey](0, addr, new(big.Int), 0, new(big.Int), nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender[gost3410.PublicKey](signer, tx)
	if err != nil {
		t.Fatal(err)
	}
	if from != addr {
		t.Errorf("exected from and address to be equal. Got %x want %x", from, addr)
	}
}

func TestEIP155ChainIdGOST(t *testing.T) {
	key, _ := crypto.GenerateKey[gost3410.PrivateKey]()
	addr := crypto.PubkeyToAddress[gost3410.PublicKey](*key.Public())

	signer := NewEIP155Signer[gost3410.PublicKey](big.NewInt(18))
	tx, err := SignTx[gost3410.PrivateKey, gost3410.PublicKey](NewTransaction[gost3410.PublicKey](0, addr, new(big.Int), 0, new(big.Int), nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}
	if !tx.Protected() {
		t.Fatal("expected tx to be protected")
	}

	if tx.ChainId().Cmp(signer.chainId) != 0 {
		t.Error("expected chainId to be", signer.chainId, "got", tx.ChainId())
	}

	tx = NewTransaction[gost3410.PublicKey](0, addr, new(big.Int), 0, new(big.Int), nil)
	tx, err = SignTx[gost3410.PrivateKey, gost3410.PublicKey](tx, HomesteadSigner[gost3410.PublicKey]{}, key)
	if err != nil {
		t.Fatal(err)
	}

	if tx.Protected() {
		t.Error("didn't expect tx to be protected")
	}

	if tx.ChainId().Sign() != 0 {
		t.Error("expected chain id to be 0 got", tx.ChainId())
	}
}

func TestChainIdGOST(t *testing.T) {
	key, _ := crypto.GenerateKey[gost3410.PrivateKey]()

	tx := NewTransaction[gost3410.PublicKey](0, common.Address{}, new(big.Int), 0, new(big.Int), nil)

	var err error
	tx, err = SignTx[gost3410.PrivateKey, gost3410.PublicKey](tx, NewEIP155Signer[gost3410.PublicKey](big.NewInt(10)), key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Sender[gost3410.PublicKey](NewEIP155Signer[gost3410.PublicKey](big.NewInt(11)), tx)
	if err != ErrInvalidChainId {
		t.Error("expected error:", ErrInvalidChainId)
	}

	_, err = Sender[gost3410.PublicKey](NewEIP155Signer[gost3410.PublicKey](big.NewInt(10)), tx)
	if err != nil {
		t.Error("expected no error")
	}
}
