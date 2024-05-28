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

package keystore

import (
	"math/big"

	"github.com/MIRChain/MIR/accounts"
	"github.com/MIRChain/MIR/core/types"
)

func (w *keystoreWallet[T, P]) SignTxCsp(account accounts.Account, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTx(account, tx, chainID)
}

func (w *keystoreWallet[T, P]) SignTxWithPassphraseCsp(account accounts.Account, subjectKeyId string, pin string, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTxWithPassphraseCsp(account, subjectKeyId, pin, tx, chainID)
}
