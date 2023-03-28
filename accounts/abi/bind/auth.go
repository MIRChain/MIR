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

package bind

import (
	"errors"
	"io"
	"io/ioutil"
	"math/big"
	"reflect"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/external"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
)

// ErrNoChainID is returned whenever the user failed to specify a chain id.
var ErrNoChainID = errors.New("no chain id specified")

// ErrNotAuthorized is returned when an account is not properly unlocked.
var ErrNotAuthorized = errors.New("not authorized to sign this account")

// NewTransactor is a utility method to easily create a transaction signer from
// an encrypted json key stream and the associated passphrase.
//
// Deprecated: Use NewTransactorWithChainID instead.
func NewTransactor[T crypto.PrivateKey,P crypto.PublicKey](keyin io.Reader, passphrase string) (*TransactOpts[P], error) {
	log.Warn("WARNING: NewTransactor has been deprecated in favour of NewTransactorWithChainID")
	json, err := ioutil.ReadAll(keyin)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey[T,P](json, passphrase)
	if err != nil {
		return nil, err
	}
	return NewKeyedTransactor[T,P](key.PrivateKey), nil
}

// NewKeyStoreTransactor is a utility method to easily create a transaction signer from
// an decrypted key from a keystore.
//
// Deprecated: Use NewKeyStoreTransactorWithChainID instead.
func NewKeyStoreTransactor[T crypto.PrivateKey,P crypto.PublicKey](keystore *keystore.KeyStore[T,P], account accounts.Account) (*TransactOpts[P], error) {
	log.Warn("WARNING: NewKeyStoreTransactor has been deprecated in favour of NewTransactorWithChainID")
	var homesteadSigner types.Signer[P] = types.HomesteadSigner[P]{}
	return &TransactOpts[P]{
		From: account.Address,
		Signer: func(address common.Address, tx *types.Transaction[P]) (*types.Transaction[P], error) {
			if address != account.Address {
				return nil, ErrNotAuthorized
			}
			// Quorum
			signer := homesteadSigner
			if tx.IsPrivate() {
				signer = types.QuorumPrivateTxSigner[P]{}
			}
			// / Quorum
			signature, err := keystore.SignHash(account, signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}, nil
}

// NewKeyedTransactor is a utility method to easily create a transaction signer
// from a single private key.
//
// Deprecated: Use NewKeyedTransactorWithChainID instead.
func NewKeyedTransactor[T crypto.PrivateKey,P crypto.PublicKey](key T) *TransactOpts[P] {
	log.Warn("WARNING: NewKeyedTransactor has been deprecated in favour of NewKeyedTransactorWithChainID")
	var pub P
	switch t:=any(&key).(type) {
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
	}
	keyAddr := crypto.PubkeyToAddress(pub)
	var homesteadSigner types.Signer[P] = types.HomesteadSigner[P]{}
	return &TransactOpts[P]{
		From: keyAddr,
		Signer: func(address common.Address, tx *types.Transaction[P]) (*types.Transaction[P], error) {
			if address != keyAddr {
				return nil, ErrNotAuthorized
			}
			// Quorum
			signer := homesteadSigner
			if tx.IsPrivate() {
				signer = types.QuorumPrivateTxSigner[P]{}
			}
			// / Quorum
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

// NewTransactorWithChainID is a utility method to easily create a transaction signer from
// an encrypted json key stream and the associated passphrase.
func NewTransactorWithChainID[T crypto.PrivateKey,P crypto.PublicKey](keyin io.Reader, passphrase string, chainID *big.Int) (*TransactOpts[P], error) {
	json, err := ioutil.ReadAll(keyin)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey[T,P](json, passphrase)
	if err != nil {
		return nil, err
	}
	return NewKeyedTransactorWithChainID[T,P](key.PrivateKey, chainID)
}

// NewKeyStoreTransactorWithChainID is a utility method to easily create a transaction signer from
// an decrypted key from a keystore.
func NewKeyStoreTransactorWithChainID[T crypto.PrivateKey,P crypto.PublicKey](keystore *keystore.KeyStore[T,P], account accounts.Account, chainID *big.Int) (*TransactOpts[P], error) {
	if chainID == nil {
		return nil, ErrNoChainID
	}
	latestSigner := types.LatestSignerForChainID[P](chainID)
	log.Info("NewKeyStoreTransactorWithChainID", "latestSigner", reflect.TypeOf(latestSigner))
	return &TransactOpts[P]{
		From: account.Address,
		Signer: func(address common.Address, tx *types.Transaction[P]) (*types.Transaction[P], error) {
			if address != account.Address {
				return nil, ErrNotAuthorized
			}
			// Quorum
			signer := latestSigner
			if tx.IsPrivate() {
				signer = types.QuorumPrivateTxSigner[P]{}
			}
			// / Quorum
			signature, err := keystore.SignHash(account, signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}, nil
}

// NewKeyedTransactorWithChainID is a utility method to easily create a transaction signer
// from a single private key.
func NewKeyedTransactorWithChainID[T crypto.PrivateKey,P crypto.PublicKey](key T, chainID *big.Int) (*TransactOpts[P], error) {
	var pub P
	switch t:=any(&key).(type) {
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
	}
	keyAddr := crypto.PubkeyToAddress(pub)
	if chainID == nil {
		return nil, ErrNoChainID
	}
	latestSigner := types.LatestSignerForChainID[P](chainID)
	return &TransactOpts[P]{
		From: keyAddr,
		Signer: func(address common.Address, tx *types.Transaction[P]) (*types.Transaction[P], error) {
			if address != keyAddr {
				return nil, ErrNotAuthorized
			}
			// Quorum
			signer := latestSigner
			if tx.IsPrivate() {
				signer = types.QuorumPrivateTxSigner[P]{}
			}
			// / Quorum
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}, nil
}

// NewClefTransactor is a utility method to easily create a transaction signer
// with a clef backend.
func NewClefTransactor[T crypto.PrivateKey,P crypto.PublicKey](clef *external.ExternalSigner[P], account accounts.Account) *TransactOpts[P] {
	return &TransactOpts[P]{
		From: account.Address,
		Signer: func(address common.Address, transaction *types.Transaction[P]) (*types.Transaction[P], error) {
			if address != account.Address {
				return nil, ErrNotAuthorized
			}
			log.Info("Signing with NewClefTransactor")
			return clef.SignTx(account, transaction, transaction.ChainId()) // Clef enforces its own chain id
		},
	}
}

// Quorum
//
// NewWalletTransactor is a utility method to easily create a transaction signer
// from a wallet account
func NewWalletTransactor[P crypto.PublicKey](w accounts.Wallet[P], account accounts.Account, chainId *big.Int) *TransactOpts[P] {
	return &TransactOpts[P]{
		From: account.Address,
		Signer: func(address common.Address, transaction *types.Transaction[P]) (*types.Transaction[P], error) {
			if address != account.Address {
				return nil, ErrNotAuthorized
			}
			if transaction.ChainId() == nil {
				chainId = transaction.ChainId()
			}

			return w.SignTx(account, transaction, chainId)
		},
	}
}
