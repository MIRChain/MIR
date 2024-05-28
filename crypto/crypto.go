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

package crypto

import (
	"bufio"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/common/math"
	"github.com/MIRChain/MIR/crypto/csp"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/gost3411"
	"github.com/MIRChain/MIR/crypto/nist"
	"github.com/MIRChain/MIR/crypto/secp256k1"
	"github.com/MIRChain/MIR/params"
	"github.com/MIRChain/MIR/rlp"
	"golang.org/x/crypto/sha3"
)

// SignatureLength indicates the byte length required to carry a signature with recovery id.
const SignatureLength = 64 + 1 // 64 bytes ECDSA signature + 1 byte recovery id

// RecoveryIDOffset points to the byte offset within the signature that contains the recovery id.
const RecoveryIDOffset = 64

// DigestLength sets the signature digest exact length
const DigestLength = 32

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

var errInvalidPubkey = errors.New("invalid secp256k1 public key")

type CryptoType int

const (
	NIST CryptoType = iota
	GOST
	GOST_CSP
	PQC
)

var CryptoAlg CryptoType = NIST

type PrivateKey interface {
	nist.PrivateKey | gost3410.PrivateKey | csp.Cert | *nist.PrivateKey | *gost3410.PrivateKey | *csp.Cert
}

type PublicKey interface {
	nist.PublicKey | gost3410.PublicKey | csp.PublicKey | *nist.PublicKey | *gost3410.PublicKey | *csp.PublicKey
	GetX() *big.Int
	GetY() *big.Int
}

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// NewKeccakState creates a new KeccakState
// Mir: for code reusage we use same KeccakState interface for Streebog hash
func NewKeccakState[P PublicKey]() KeccakState {
	var pub P
	switch any(&pub).(type) {
	case *nist.PublicKey:
		return sha3.NewLegacyKeccak256().(KeccakState)
	case *gost3410.PublicKey:
		return gost3411.New256().(KeccakState)
	case *csp.PublicKey:
		return csp.New256().(KeccakState)
	default:
		panic("cant infer crypto type for hashing alg")
	}
}

// HashData hashes the provided data using the KeccakState and returns a 32 byte hash
func HashData[P PublicKey](kh KeccakState, data []byte) (h common.Hash) {
	kh.Reset()
	kh.Write(data)
	kh.Read(h[:])
	return h
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256[P PublicKey](data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState[P]()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash[P PublicKey](data ...[]byte) (h common.Hash) {
	d := NewKeccakState[P]()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	return h
}

// Keccak512 calculates and returns the Keccak512 hash of the input data.
func Keccak512[P PublicKey](data ...[]byte) []byte {
	var pub P
	switch any(&pub).(type) {
	case *nist.PublicKey:
		d := sha3.NewLegacyKeccak512()
		for _, b := range data {
			d.Write(b)
		}
		return d.Sum(nil)
	case *gost3410.PublicKey:
		d := gost3411.New512()
		for _, b := range data {
			d.Write(b)
		}
		return d.Sum(nil)
	default:
		panic("cant infer crypto type for hashing alg")
	}
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress[P PublicKey](b common.Address, nonce uint64) common.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return common.BytesToAddress(Keccak256[P](data)[12:])
}

// CreateAddress2 creates an ethereum address given the address bytes, initial
// contract code hash and a salt.
func CreateAddress2[P PublicKey](b common.Address, salt [32]byte, inithash []byte) common.Address {
	return common.BytesToAddress(Keccak256[P]([]byte{0xff}, b.Bytes(), salt[:], inithash)[12:])
}

// ToECDSA creates a private key with the given D value.
func ToECDSA[T PrivateKey](d []byte) (T, error) {
	return toECDSA[T](d, true)
}

// ToECDSAUnsafe blindly converts a binary blob to a private key. It should almost
// never be used unless you are sure the input is valid and want to avoid hitting
// errors due to bad origin encoding (0 prefixes cut off).
func ToECDSAUnsafe[T PrivateKey](d []byte) T {
	priv, _ := toECDSA[T](d, false)
	return priv
}

// toECDSA creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSA[T PrivateKey](d []byte, strict bool) (T, error) {
	var prv T
	switch p := any(&prv).(type) {
	case *nist.PrivateKey:
		priv := new(ecdsa.PrivateKey)
		priv.PublicKey.Curve = S256()
		if strict && 8*len(d) != priv.Params().BitSize {
			return ZeroPrivateKey[T](), fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
		}
		priv.D = new(big.Int).SetBytes(d)

		// The priv.D must < N
		if priv.D.Cmp(secp256k1N) >= 0 {
			return ZeroPrivateKey[T](), fmt.Errorf("invalid private key, >=N")
		}
		// The priv.D must not be zero or negative.
		if priv.D.Sign() <= 0 {
			return ZeroPrivateKey[T](), fmt.Errorf("invalid private key, zero or negative")
		}

		priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
		if priv.PublicKey.X == nil {
			return ZeroPrivateKey[T](), errors.New("invalid private key")
		}
		*p = nist.PrivateKey{priv}
	case *gost3410.PrivateKey:
		priv := new(gost3410.PrivateKey)
		priv.C = gost3410.GostCurve
		priv.PublicKey.C = priv.C
		if strict && 8*len(d) != priv.C.Params().BitSize {
			return ZeroPrivateKey[T](), fmt.Errorf("invalid length, need %d bits", priv.C.Params().BitSize)
		}
		priv.Key = new(big.Int).SetBytes(d)

		// The priv.D must < N
		if priv.Key.Cmp(gost3410.GostCurve.Q) >= 0 {
			return ZeroPrivateKey[T](), fmt.Errorf("invalid private key, >=N")
		}
		// The priv.D must not be zero or negative.
		if priv.Key.Sign() <= 0 {
			return ZeroPrivateKey[T](), fmt.Errorf("invalid private key, zero or negative")
		}
		priv.PublicKey.X, priv.PublicKey.Y = priv.C.ScalarBaseMult(d)
		if priv.PublicKey.X == nil {
			return ZeroPrivateKey[T](), errors.New("invalid private key")
		}
		*p = *priv
	}
	return prv, nil
}

// FromECDSA exports a private key into a binary dump.
func FromECDSA[T PrivateKey](priv T) []byte {
	switch priv := any(&priv).(type) {
	case *nist.PrivateKey:
		if priv == nil {
			return nil
		}
		return math.PaddedBigBytes(priv.D, priv.Params().BitSize/8)
	case *gost3410.PrivateKey:
		if priv == nil {
			return nil
		}
		return math.PaddedBigBytes(priv.Key, priv.C.P.BitLen()/8)
	case *csp.Cert:
		return priv.Bytes()
	default:
		panic("cant infer priv key")
	}

}

// UnmarshalPubkey converts bytes to a secp256k1 public key.
func UnmarshalPubkey[P PublicKey](pub []byte) (P, error) {
	var pubKey P
	switch p := any(&pubKey).(type) {
	case *nist.PublicKey:
		x, y := elliptic.Unmarshal(S256(), pub)
		if x == nil {
			return ZeroPublicKey[P](), errInvalidPubkey
		}
		*p = nist.PublicKey{&ecdsa.PublicKey{Curve: S256(), X: x, Y: y}}
	case *gost3410.PublicKey:
		k, err := gost3410.NewPublicKey(gost3410.GostCurve, pub[1:])
		if err != nil {
			return ZeroPublicKey[P](), err
		}
		*p = *k
	case *csp.PublicKey:
		k, err := csp.NewPublicKey(pub)
		if err == nil {
			return ZeroPublicKey[P](), err
		}
		*p = *k
	}
	return pubKey, nil
}

func FromECDSAPub[P PublicKey](pub P) []byte {
	switch p := any(&pub).(type) {
	case *nist.PublicKey:
		if pub.GetX() == nil || pub.GetY() == nil {
			panic("nil nil")
		}
		return elliptic.Marshal(S256(), pub.GetX(), pub.GetY())
	case *gost3410.PublicKey:
		if pub.GetX() == nil || pub.GetY() == nil {
			panic("nil nil")
		}
		return gost3410.Marshal(gost3410.GostCurve, p.X, p.Y)
	case *csp.PublicKey:
		if pub.GetX() == nil || pub.GetY() == nil {
			panic("nil nil")
		}
		return csp.Marshal(*gost3410.CurveIdGostR34102001CryptoProAParamSet(), p.X, p.Y)
	default:
		panic("cant infer pubkey type")
	}
}

// HexToECDSA parses a secp256k1 private key.
func HexToECDSA[T PrivateKey](hexkey string) (T, error) {
	b, err := hex.DecodeString(hexkey)
	if byteErr, ok := err.(hex.InvalidByteError); ok {
		return ZeroPrivateKey[T](), fmt.Errorf("invalid hex character %q in private key", byte(byteErr))
	} else if err != nil {
		return ZeroPrivateKey[T](), errors.New("invalid hex data for private key")
	}
	return ToECDSA[T](b)
}

// LoadECDSA loads a secp256k1 private key from the given file.
func LoadECDSA[T PrivateKey](file string) (T, error) {
	fd, err := os.Open(file)
	if err != nil {
		return ZeroPrivateKey[T](), err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	buf := make([]byte, 64)
	n, err := readASCII(buf, r)
	if err != nil {
		return ZeroPrivateKey[T](), err
	} else if n != len(buf) {
		return ZeroPrivateKey[T](), fmt.Errorf("key file too short, want 64 hex characters")
	}
	if err := checkKeyFileEnd(r); err != nil {
		return ZeroPrivateKey[T](), err
	}

	return HexToECDSA[T](string(buf))
}

// readASCII reads into 'buf', stopping when the buffer is full or
// when a non-printable control character is encountered.
func readASCII(buf []byte, r *bufio.Reader) (n int, err error) {
	for ; n < len(buf); n++ {
		buf[n], err = r.ReadByte()
		switch {
		case err == io.EOF || buf[n] < '!':
			return n, nil
		case err != nil:
			return n, err
		}
	}
	return n, nil
}

// checkKeyFileEnd skips over additional newlines at the end of a key file.
func checkKeyFileEnd(r *bufio.Reader) error {
	for i := 0; ; i++ {
		b, err := r.ReadByte()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case b != '\n' && b != '\r':
			return fmt.Errorf("invalid character %q at end of key file", b)
		case i >= 2:
			return errors.New("key file too long, want 64 hex characters")
		}
	}
}

// SaveECDSA saves a secp256k1 private key to the given file with
// restrictive permissions. The key data is saved hex-encoded.
func SaveECDSA[T PrivateKey](file string, key T) error {
	k := hex.EncodeToString(FromECDSA(key))
	return ioutil.WriteFile(file, []byte(k), 0600)
}

func SaveECDSAGost(file string, key *gost3410.PrivateKey) error {
	k := hex.EncodeToString(key.Raw())
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// GenerateKey generates a new private key.
func GenerateKey[T PrivateKey]() (T, error) {
	var key T
	switch t := any(&key).(type) {
	case *nist.PrivateKey:
		new256k1, err := ecdsa.GenerateKey(S256(), rand.Reader)
		if err != nil {
			return ZeroPrivateKey[T](), fmt.Errorf("error generating secp256k1 key")
		}
		*t = nist.PrivateKey{new256k1}
	case *gost3410.PrivateKey:
		newGost3410, err := gost3410.GenPrivateKey(gost3410.GostCurve, rand.Reader)
		if err != nil {
			return ZeroPrivateKey[T](), fmt.Errorf("")
		}
		*t = *newGost3410
	case *csp.Cert:
		*t = *params.SignerCert
	}
	return key, nil
}

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues[P crypto.PublicKey](v byte, r, s *big.Int, homestead bool) bool {
	var pub P
	switch any(&pub).(type) {
	case *nist.PublicKey:
		if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
			return false
		}
		// reject upper range of s values (ECDSA malleability)
		// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
		if homestead && s.Cmp(secp256k1halfN) > 0 {
			return false
		}
		// Frontier: allow s to be in full N range
		return r.Cmp(secp256k1N) < 0 && s.Cmp(secp256k1N) < 0 && (v == 0 || v == 1 || v == 10 || v == 11)
	case *gost3410.PublicKey:
		if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
			return false
		}
		return r.Cmp(gost3410.GostCurve.Q) < 0 && s.Cmp(gost3410.GostCurve.Q) < 0 && (v == 0 || v == 1 || v == 10 || v == 11)
	case *csp.PublicKey:
		if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
			return false
		}
		return true
	default:
		return false
	}
}

func PubkeyToAddress[P PublicKey](pub P) common.Address {
	pubBytes := FromECDSAPub(pub)
	return common.BytesToAddress(Keccak256[P](pubBytes[1:])[12:])
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}

func ZeroPrivateKey[T PrivateKey]() T {
	var res T
	return res
}

func ZeroPublicKey[P PublicKey]() P {
	var res P
	return res
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
}
