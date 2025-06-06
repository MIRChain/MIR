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

package ethclient

import (
	"errors"
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// senderFromServer is a types.Signer that remembers the sender address returned by the RPC
// server. It is stored in the transaction's sender address cache to avoid an additional
// request in TransactionSender.
type senderFromServer [P crypto.PublicKey]  struct {
	addr      common.Address
	blockhash common.Hash
}

var errNotCached = errors.New("sender not cached")

func setSenderFromServer[P crypto.PublicKey](tx *types.Transaction[P], addr common.Address, block common.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	types.Sender[P](&senderFromServer[P]{addr, block}, tx)
}

func (s *senderFromServer[P]) Equal(other types.Signer[P]) bool {
	os, ok := other.(*senderFromServer[P])
	return ok && os.blockhash == s.blockhash
}

func (s *senderFromServer[P]) Sender(tx *types.Transaction[P]) (common.Address, error) {
	if s.blockhash == (common.Hash{}) {
		return common.Address{}, errNotCached
	}
	return s.addr, nil
}

func (s *senderFromServer[P]) ChainID() *big.Int {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer[P]) Hash(tx *types.Transaction[P]) common.Hash {
	panic("can't sign with senderFromServer")
}
func (s *senderFromServer[P]) SignatureValues(tx *types.Transaction[P], sig []byte) (R, S, V *big.Int, err error) {
	panic("can't sign with senderFromServer")
}
