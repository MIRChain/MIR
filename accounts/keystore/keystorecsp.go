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

// Package keystore implements encrypted storage of secp256k1 private keys.
//
// Keys are stored as encrypted JSON files according to the Web3 Secret Storage specification.
// See https://github.com/ethereum/wiki/wiki/Web3-Secret-Storage-Definition for more information.
package keystore

import (
	crand "crypto/rand"
	"math/big"
	"os"
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/csp"
	"github.com/pavelkrolevets/MIR-pro/log"
)

func (ks *KeyStore[T,P]) SignHashCsp(a accounts.Account, hash []byte) ([]byte, error) {
	// Look up the key to sign with and abort if it cannot be found
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	unlockedKey, found := ks.unlocked[a.Address]
	if !found {
		return nil, ErrLocked
	}
	store, err := csp.SystemStore("My")
	if err != nil {
		return nil, err
	}
	defer store.Close()
	crt, err := store.GetBySubjectId(unlockedKey.KeyCsp.SubjectKeyId)
	if err != nil {
		return nil, err
	}
	defer crt.Close()
	// Sign the hash using plain ECDSA operations
	return crypto.Sign(hash, crt)
}

func (ks *KeyStore[T,P]) DeleteCsp(a accounts.Account, passphrase string) error {
	// Decrypting the key isn't really necessary, but we do
	// it anyway to check the password and zero out the key
	// immediately afterwards.
	a, key, err := ks.getDecryptedKeyCsp(a)
	if key != nil {
		key.SubjectKeyId = ""
	}
	if err != nil {
		return err
	}
	// The order is crucial here. The key is dropped from the
	// cache after the file is gone so that a reload happening in
	// between won't insert it into the cache again.
	err = os.Remove(a.URL.Path)
	if err == nil {
		ks.cache.delete(a)
		ks.refreshWallets()
	}
	return err
}

func (ks *KeyStore[T,P]) SignTxCsp(a accounts.Account, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	// Look up the key to sign with and abort if it cannot be found
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	unlockedKey, found := ks.unlocked[a.Address]
	if !found {
		return nil, ErrLocked
	}

	store, err := csp.SystemStore("My")
	if err != nil {
		return nil, err
	}
	defer store.Close()
	crt, err := store.GetBySubjectId(unlockedKey.KeyCsp.SubjectKeyId)
	if err != nil {
		return nil, err
	}
	defer crt.Close()
	var priv T
	t := any(&priv).(*csp.Cert)
	*t = crt
	// start quorum specific
	if tx.IsPrivate() {
		log.Info("Private transaction signing with QuorumPrivateTxSigner")
		return types.SignTx[T,P](tx, types.QuorumPrivateTxSigner[P]{}, priv)
	} // End quorum specific

	// Depending on the presence of the chain ID, sign with 2718 or homestead
	signer := types.LatestSignerForChainID[P](chainID)
	return types.SignTx[T,P](tx, signer, priv)
}

func (ks *KeyStore[T,P]) SignHashWithPassphraseCsp(a accounts.Account, subjectKeyId string, pin string, hash []byte) (signature []byte, err error) {
	_, key, err := ks.getDecryptedKeyCsp(a)
	if err != nil {
		return nil, err
	}
	store, err := csp.SystemStore("My")
	if err != nil {
		return nil, err
	}
	defer store.Close()
	crt, err := store.GetBySubjectId(key.SubjectKeyId)
	if err != nil {
		return nil, err
	}
	defer crt.Close()
	return crypto.Sign(hash, &crt)
}

func (ks *KeyStore[T,P]) SignTxWithPassphraseCsp(a accounts.Account, subjectKeyId string, pin string, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	_, key, err := ks.getDecryptedKeyCsp(a)
	if err != nil {
		return nil, err
	}
	store, err := csp.SystemStore("My")
	if err != nil {
		return nil, err
	}
	defer store.Close()
	crt, err := store.GetBySubjectId(key.SubjectKeyId)
	if err != nil {
		return nil, err
	}
	defer crt.Close()
	var priv T
	t := any(&priv).(*csp.Cert)
	*t = crt
	if tx.IsPrivate() {
		return types.SignTx[T,P](tx, types.QuorumPrivateTxSigner[P]{}, priv)
	}
	// Depending on the presence of the chain ID, sign with or without replay protection.
	signer := types.LatestSignerForChainID[P](chainID)
	return types.SignTx[T,P](tx, signer, priv)
}

func (ks *KeyStore[T,P]) getDecryptedKeyCsp(a accounts.Account) (accounts.Account, *KeyCsp, error) {
	a, err := ks.Find(a)
	if err != nil {
		return a, nil, err
	}
	key, err := ks.storage.GetKeyCsp(a.Address, a.URL.Path)
	return a, key, err
}

func (ks *KeyStore[T,P]) NewAccountCsp(subjectKeyId string) (accounts.Account, error) {
	_, account, err := storeNewKeyCsp[T,P](ks.storage, crand.Reader, subjectKeyId)
	if err != nil {
		return accounts.Account{}, err
	}
	// Add the account to the cache immediately rather
	// than waiting for file system notifications to pick it up.
	ks.cache.add(account)
	ks.refreshWallets()
	return account, nil
}

// Unlock unlocks the given account indefinitely.
func (ks *KeyStore[T,P]) UnlockCsp(a accounts.Account, pin string) error {
	return ks.TimedUnlockCsp(a, pin, 0)
}

// Lock removes the private key with the given address from memory.
func (ks *KeyStore[T,P]) LockCsp(addr common.Address) error {
	ks.mu.Lock()
	if unl, found := ks.unlocked[addr]; found {
		ks.mu.Unlock()
		ks.expire(addr, unl, time.Duration(0)*time.Nanosecond)
	} else {
		ks.mu.Unlock()
	}
	return nil
}

func (ks *KeyStore[T,P]) TimedUnlockCsp(a accounts.Account, pin string, timeout time.Duration) error {
	a, key, err := ks.getDecryptedKeyCsp(a)
	if err != nil {
		return err
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()
	u, found := ks.unlocked[a.Address]
	if found {
		if u.abort == nil {
			// The address was unlocked indefinitely, so unlocking
			// it with a timeout would be confusing.
			zeroKeyCsp(u.SubjectKeyId)
			return nil
		}
		// Terminate the expire goroutine and replace it below.
		close(u.abort)
	}
	if timeout > 0 {
		u = &unlocked[T]{KeyCsp: key, abort: make(chan struct{})}
		go ks.expire(a.Address, u, timeout)
	} else {
		u = &unlocked[T]{KeyCsp: key}
	}
	ks.unlocked[a.Address] = u
	return nil
}
func zeroKeyCsp(k string) {
	k = ""
}