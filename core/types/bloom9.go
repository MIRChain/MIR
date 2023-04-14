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

package types

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
)

type bytesBacked interface {
	Bytes() []byte
}

const (
	// BloomByteLength represents the number of bytes used in a header log bloom.
	BloomByteLength = 256

	// BloomBitLength represents the number of bits used in a header log bloom.
	BloomBitLength = 8 * BloomByteLength
)

// Bloom represents a 2048 bit bloom filter.
type Bloom[P crypto.PublicKey] [BloomByteLength]byte

// BytesToBloom converts a byte slice to a bloom filter.
// It panics if b is not of suitable size.
func BytesToBloom[P crypto.PublicKey](b []byte) Bloom[P] {
	var bloom Bloom[P]
	bloom.SetBytes(b)
	return bloom
}

// SetBytes sets the content of b to the given bytes.
// It panics if d is not of suitable size.
func (b *Bloom[P]) SetBytes(d []byte) {
	if len(b) < len(d) {
		panic(fmt.Sprintf("bloom bytes too big %d %d", len(b), len(d)))
	}
	copy(b[BloomByteLength-len(d):], d)
}

// Add adds d to the filter. Future calls of Test(d) will return true.
func (b *Bloom[P]) Add(d []byte) {
	b.add(d, make([]byte, 6))
}

// add is internal version of Add, which takes a scratch buffer for reuse (needs to be at least 6 bytes)
func (b *Bloom[P]) add(d []byte, buf []byte) {
	i1, v1, i2, v2, i3, v3 := bloomValues[P](d, buf)
	b[i1] |= v1
	b[i2] |= v2
	b[i3] |= v3
}

// Quorum
// OrBloom executes an Or operation on the bloom
func (b *Bloom[P]) OrBloom(bl []byte) {
	bin := new(big.Int).SetBytes(b[:])
	input := new(big.Int).SetBytes(bl[:])
	bin.Or(bin, input)
	b.SetBytes(bin.Bytes())
}

// Big converts b to a big integer.
// Note: Converting a bloom filter to a big.Int and then calling GetBytes
// does not return the same bytes, since big.Int will trim leading zeroes
func (b Bloom[P]) Big() *big.Int {
	return new(big.Int).SetBytes(b[:])
}

// Bytes returns the backing byte slice of the bloom
func (b Bloom[P]) Bytes() []byte {
	return b[:]
}

// Test checks if the given topic is present in the bloom filter
func (b Bloom[P]) Test(topic []byte) bool {
	i1, v1, i2, v2, i3, v3 := bloomValues[P](topic, make([]byte, 6))
	return v1 == v1&b[i1] &&
		v2 == v2&b[i2] &&
		v3 == v3&b[i3]
}

// MarshalText encodes b as a hex string with 0x prefix.
func (b Bloom[P]) MarshalText() ([]byte, error) {
	return hexutil.Bytes(b[:]).MarshalText()
}

// UnmarshalText b as a hex string with 0x prefix.
func (b *Bloom[P]) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Bloom", input, b[:])
}

// CreateBloom creates a bloom filter out of the give Receipts (+Logs)
func CreateBloom[P crypto.PublicKey](receipts Receipts[P]) Bloom[P] {
	buf := make([]byte, 6)
	var bin Bloom[P]
	for _, receipt := range receipts {
		for _, log := range receipt.Logs {
			bin.add(log.Address.Bytes(), buf)
			for _, b := range log.Topics {
				bin.add(b[:], buf)
			}
		}
	}
	return bin
}

// LogsBloom returns the bloom bytes for the given logs
func LogsBloom[P crypto.PublicKey](logs []*Log) []byte {
	buf := make([]byte, 6)
	var bin Bloom[P]
	for _, log := range logs {
		bin.add(log.Address.Bytes(), buf)
		for _, b := range log.Topics {
			bin.add(b[:], buf)
		}
	}
	return bin[:]
}

// Bloom9 returns the bloom filter for the given data
func Bloom9[P crypto.PublicKey](data []byte) []byte {
	var b Bloom[P]
	b.SetBytes(data)
	return b.Bytes()
}

// bloomValues returns the bytes (index-value pairs) to set for the given data
func bloomValues[P crypto.PublicKey](data []byte, hashbuf []byte) (uint, byte, uint, byte, uint, byte) {
	var pub P
	switch any(&pub).(type){
	case *nist.PublicKey:
		sha := hasherPoolNist.Get().(crypto.KeccakState)
		sha.Reset()
		sha.Write(data)
		sha.Read(hashbuf)
		hasherPoolNist.Put(sha)
	case *gost3410.PublicKey:
		sha := hasherPoolGost.Get().(crypto.KeccakState)
		sha.Reset()
		sha.Write(data)
		sha.Read(hashbuf)
		hasherPoolGost.Put(sha)
	}
	// The actual bits to flip
	v1 := byte(1 << (hashbuf[1] & 0x7))
	v2 := byte(1 << (hashbuf[3] & 0x7))
	v3 := byte(1 << (hashbuf[5] & 0x7))
	// The indices for the bytes to OR in
	i1 := BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf)&0x7ff)>>3) - 1
	i2 := BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf[2:])&0x7ff)>>3) - 1
	i3 := BloomByteLength - uint((binary.BigEndian.Uint16(hashbuf[4:])&0x7ff)>>3) - 1

	return i1, v1, i2, v2, i3, v3
}

// BloomLookup is a convenience-method to check presence int he bloom filter
func BloomLookup[P crypto.PublicKey](bin Bloom[P], topic bytesBacked) bool {
	return bin.Test(topic.Bytes())
}
