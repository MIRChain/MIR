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

package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type keyStorePlain [T crypto.PrivateKey] struct {
	keysDirPath string
}

func (ks keyStorePlain[T]) GetKey(addr common.Address, filename, auth string) (*Key[T], error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	key := new(Key[T])
	if err := json.NewDecoder(fd).Decode(key); err != nil {
		return nil, err
	}
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have address %x, want %x", key.Address, addr)
	}
	return key, nil
}

func (ks keyStorePlain[T]) StoreKey(filename string, key *Key[T], auth string) error {
	content, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return writeKeyFile(filename, content)
}

func (ks keyStorePlain[T]) JoinPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(ks.keysDirPath, filename)
}
