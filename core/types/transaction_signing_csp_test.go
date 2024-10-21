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
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/csp"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/nist"
)

func TestEIP155SigningCSP(t *testing.T) {
	store, err := csp.SystemStore("My")
	if err != nil {
		t.Fatalf("Store error: %s", err)
	}
	defer store.Close()
	// Cert should be without set pin
	crt, err := store.GetBySubjectId("71732462bbc029d911e6d16a3ed00d9d1d772620")
	if err != nil {
		t.Fatalf("Get cert error: %s", err)
	}
	t.Logf("Cert pub key: %x", crt.Info().PublicKeyBytes()[2:66])
	defer crt.Close()
	addr := crypto.PubkeyToAddress[csp.PublicKey](*crt.Public())
	t.Logf("Cert address: %s", addr.Hex())
	signer := NewEIP155Signer[csp.PublicKey](big.NewInt(18))
	tx, err := SignTx[csp.Cert, csp.PublicKey](NewTransaction[csp.PublicKey](0, addr, new(big.Int), 0, new(big.Int), nil), signer, crt)
	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender[csp.PublicKey](signer, tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("From address: %s", from.Hex())
	if from != addr {
		t.Errorf("exected from and address to be equal. Got %x want %x", from, addr)
	}
}

func TestEIP155ChainIdCSP(t *testing.T) {
	store, err := csp.SystemStore("My")
	if err != nil {
		t.Fatalf("Store error: %s", err)
	}
	defer store.Close()
	// Cert should be without set pin
	crt, err := store.GetBySubjectId("71732462bbc029d911e6d16a3ed00d9d1d772620")
	if err != nil {
		t.Fatalf("Get cert error: %s", err)
	}
	t.Logf("Cert pub key: %x", crt.Info().PublicKeyBytes()[2:66])
	defer crt.Close()
	addr := crypto.PubkeyToAddress[csp.PublicKey](*crt.Public())

	signer := NewEIP155Signer[csp.PublicKey](big.NewInt(18))
	tx, err := SignTx[csp.Cert, csp.PublicKey](NewTransaction[csp.PublicKey](0, addr, new(big.Int), 0, new(big.Int), nil), signer, crt)
	if err != nil {
		t.Fatal(err)
	}
	if !tx.Protected() {
		t.Fatal("expected tx to be protected")
	}

	if tx.ChainId().Cmp(signer.chainId) != 0 {
		t.Error("expected chainId to be", signer.chainId, "got", tx.ChainId())
	}

	tx = NewTransaction[csp.PublicKey](0, addr, new(big.Int), 0, new(big.Int), nil)
	tx, err = SignTx[csp.Cert, csp.PublicKey](tx, HomesteadSigner[csp.PublicKey]{}, crt)
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

func TestChainIdCSP(t *testing.T) {
	key, _ := defaultTestKey()

	tx := NewTransaction[nist.PublicKey](0, common.Address{}, new(big.Int), 0, new(big.Int), nil)

	var err error
	tx, err = SignTx[nist.PrivateKey, nist.PublicKey](tx, NewEIP155Signer[nist.PublicKey](big.NewInt(10)), key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Sender[nist.PublicKey](NewEIP155Signer[nist.PublicKey](big.NewInt(11)), tx)
	if err != ErrInvalidChainId {
		t.Error("expected error:", ErrInvalidChainId)
	}

	_, err = Sender[nist.PublicKey](NewEIP155Signer[nist.PublicKey](big.NewInt(10)), tx)
	if err != nil {
		t.Error("expected no error")
	}
}

func TestSignTxCSP(t *testing.T) {
	store, err := csp.SystemStore("My")
	if err != nil {
		t.Fatalf("Store error: %s", err)
	}
	defer store.Close()
	// Cert should be without set pin
	crt, err := store.GetBySubjectId("71732462bbc029d911e6d16a3ed00d9d1d772620")
	if err != nil {
		t.Fatalf("Get cert error: %s", err)
	}
	defer crt.Close()

	crtPub := crt.Info().PublicKeyBytes()[2:66]
	reverse(crtPub)
	// Get address which will be used at pure GO GOST network
	value := big.NewInt(10000000)
	gasLimit := uint64(21000)
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	var data []byte
	tx := NewTransaction[gost3410.PublicKey](uint64(0), toAddress, value, gasLimit, big.NewInt(0), data)
	signerCSP := NewEIP155Signer[gost3410.PublicKey](big.NewInt(1515))
	txHash := signerCSP.Hash(tx)
	sig, err := crypto.Sign(txHash[:], crt)
	if err != nil {
		t.Fatal(err)
	}
	r, s, _ := crypto.RevertCSP(txHash[:], sig)
	resSig := make([]byte, 65)
	copy(resSig[:32], r.Bytes())
	copy(resSig[32:64], s.Bytes())
	resSig[64] = sig[64]
	// Get address which will be used at pure GO GOST network - recover to get the right value
	recoveredGostPub, err := crypto.Ecrecover[gost3410.PublicKey](txHash[:], resSig)
	if err != nil {
		t.Fatal(err)
	}
	var addrFrom common.Address
	copy(addrFrom[:], crypto.Keccak256[gost3410.PublicKey](recoveredGostPub[1:])[12:])
	pub := make([]byte, 64)
	copy(pub[:32], recoveredGostPub[33:65])
	copy(pub[32:64], recoveredGostPub[1:33])
	reverse(pub)
	// compare recovered and crt.pub
	if !bytes.Equal(crt.Info().PublicKeyBytes()[2:66], pub) {
		t.Fatal("Wrong recovered pub key")
	}
	signerGost := NewEIP155Signer[gost3410.PublicKey](big.NewInt(1515))
	signedTx, err := tx.WithSignature(signerGost, resSig)
	if err != nil {
		t.Fatal(err)
	}
	V, R, S := signedTx.RawSignatureValues()
	V = new(big.Int).Sub(V, signerGost.chainIdMul)
	V.Sub(V, big8)
	recoveredAddress, err := recoverPlain[gost3410.PublicKey](txHash, R, S, V, true)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(addrFrom, recoveredAddress) {
		t.Fatal("Wrong recovered pub key")
	}
}
