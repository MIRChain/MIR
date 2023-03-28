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

package keystore

import (
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/csp"
)


type KeyCsp struct {
	Id uuid.UUID // Version 4 "random" for unique id not derived from key data
	// to simplify lookups we also store the address
	Address common.Address
	// Cert ID
	SubjectKeyId string
	// Pin to unlock the cert
	Pin 		string
}

type plainKeyJSONCsp struct {
	Address    string `json:"address"`
	Id         string `json:"id"`
	Version    int    `json:"version"`
}

type cipherparamsJSONCsp struct {
	IV string `json:"iv"`
}

func newKeyCsp(subjectKeyId string) (*KeyCsp, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}
	store, err := csp.SystemStore("My")
	if err != nil {
		return nil, err
	}
	defer store.Close()
	crt, err := store.GetBySubjectId(subjectKeyId)
	if err != nil {
		return nil, err
	}
	key := &KeyCsp{
		Id:         id,
		Address:    crypto.PubkeyToAddressCsp(crt.Info().PublicKeyBytes()),
		SubjectKeyId: subjectKeyId,
	}
	return key, nil
}

func storeNewKeyCsp[T crypto.PrivateKey, P crypto.PublicKey](ks keyStore[T], rand io.Reader, subjectKeyId string) (*KeyCsp, accounts.Account, error) {
	key, err := newKeyCsp(subjectKeyId)
	if err != nil {
		return nil, accounts.Account{}, err
	}
	a := accounts.Account{
		Address: key.Address,
		URL:     accounts.URL{Scheme: KeyStoreScheme, Path: ks.JoinPath(keyFileName(key.Address))},
	}
	if err := ks.StoreKeyCsp(a.URL.Path, key); err != nil {
		return nil, a, err
	}
	return key, a, err
}