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

	"github.com/pavelkrolevets/MIR-pro/common"
)


func (ks keyStorePlain[T]) GetKeyCsp(addr common.Address, filename string) (*KeyCsp, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	key := new(KeyCsp)
	if err := json.NewDecoder(fd).Decode(key); err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	return key, nil
}

func (ks keyStorePlain[T]) StoreKeyCsp(filename string, key *KeyCsp) error {
	content, err := json.Marshal(key)
	if err != nil {
		return err
	}
	return writeKeyFile(filename, content)
}
