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

package istanbul

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"golang.org/x/crypto/sha3"
)

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress[P crypto.PublicKey](data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256[P](data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub[P](hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(pubkey), nil
}

// GetSignatureAddressNoHashing gets the signer address from the signature without first hashing the data
func GetSignatureAddressNoHashing[P crypto.PublicKey](data []byte, sig []byte) (common.Address, error) {
	pubkey, err := crypto.SigToPub[P](data, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(pubkey), nil
}

func CheckValidatorSignature[P crypto.PublicKey](valSet ValidatorSet, data []byte, sig []byte) (common.Address, error) {
	// 1. Get signature address
	signer, err := GetSignatureAddress[P](data, sig)
	if err != nil {
		log.Error("Failed to get signer address", "err", err)
		return common.Address{}, err
	}

	// 2. Check validator
	if _, val := valSet.GetByAddress(signer); val != nil {
		return val.Address(), nil
	}

	return common.Address{}, ErrUnauthorizedAddress
}
