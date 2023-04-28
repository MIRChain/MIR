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

//go:build !nacl && !js && cgo
// +build !nacl,!js,cgo

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/crypto/secp256k1"
)

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover[P PublicKey](digestHash, sig []byte) ([]byte, error) {
	if len(digestHash) != 32 {
		return nil, fmt.Errorf("hash length should be 32 bytes, got %d", len(digestHash))
	}
	if len(sig) != 65 {
		return nil, fmt.Errorf("sig length should be 65 bytes, got %d", len(digestHash))
	}
	var pubKey P
	switch any(&pubKey).(type) {
	case *nist.PublicKey:
		return secp256k1.RecoverPubkey(digestHash, sig)
	case *gost3410.PublicKey:
		v := int((sig[64]) & ^byte(4))
		r := new(big.Int).SetBytes(sig[:32])
		s := new(big.Int).SetBytes(sig[32:64])
		revHash := make([]byte, 32)
		copy(revHash, digestHash)
		reverse(revHash)
		X, Y, err := gost3410.RecoverCompact(*gost3410.GostCurve, revHash, r, s, v)
		pubKey := gost3410.PublicKey{
			C: gost3410.GostCurve,
			X: X,
			Y: Y,
		}
		if err != nil{
			return nil, err
		}
		return gost3410.Marshal(gost3410.GostCurve, pubKey.X, pubKey.Y), nil
	default:
		return nil, fmt.Errorf("no crypto alg was set")
	}
}
// SigToPub returns the public key that created the given signature.
func SigToPub[P PublicKey](hash, sig []byte) (P, error) {
	var pubKey P
	switch p := any(&pubKey).(type) {
	case *nist.PublicKey:
		s, err := Ecrecover[P](hash, sig)
		if err != nil {
			return ZeroPublicKey[P](), err
		}

		x, y := elliptic.Unmarshal(S256(), s)
		*p = nist.PublicKey{&ecdsa.PublicKey{Curve: S256(), X: x, Y: y}}
	case *gost3410.PublicKey:
		s, err := Ecrecover[P](hash, sig)
		if err != nil {
			return ZeroPublicKey[P](), err
		}
		key, err := gost3410.NewPublicKey(gost3410.GostCurve,s)
		if err != nil {
			return ZeroPublicKey[P](), err
		}
		*p = *key
	default:
		return ZeroPublicKey[P](), fmt.Errorf("cant infer pub key type")
	}
	return pubKey, nil
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given digest cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign [T PrivateKey](digestHash []byte, key T) (sig []byte, err error) {
	switch prv:=any(&key).(type) {
	case *nist.PrivateKey:
		if len(digestHash) != DigestLength {
			return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", DigestLength, len(digestHash))
		}
		seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
		defer zeroBytes(seckey)
		return secp256k1.Sign(digestHash, seckey)
	case *gost3410.PrivateKey:
		var resSig []byte
		if len(digestHash) != DigestLength {
			return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", DigestLength, len(digestHash))
		}
		revHash := make([]byte, 32)
		copy(revHash, digestHash)
		reverse(revHash)
		sig, err := prv.SignDigest(revHash, rand.Reader)
		if err != nil {
			return nil, err
		}
		r := new(big.Int).SetBytes(sig[:32])
		s := new(big.Int).SetBytes(sig[32:64])
		for i := 0; i < (1+1)*2; i++ {
			X, Y, err := gost3410.RecoverCompact(*prv.C, revHash, r, s, i)
			if err == nil && X.Cmp(prv.PublicKey.X) == 0 && Y.Cmp(prv.PublicKey.Y) == 0 {
				resSig = append(resSig, sig...)
				resSig = append(resSig, byte(i))
			}
		}
		return resSig, nil
	default:
		return nil, fmt.Errorf("cant identify prv key type")
	}

}

// VerifySignature checks that the given public key created signature over digest.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature[P PublicKey](pubkey, digestHash, signature []byte) bool {
	var pubKey P
	switch any(&pubKey).(type) {
	case *nist.PublicKey:
		return secp256k1.VerifySignature(pubkey, digestHash, signature)
	case *gost3410.PublicKey:
		pub, err := gost3410.NewPublicKey(gost3410.GostCurve, pubkey[1:])
		if err != nil {
			return false
		}
		revHash := make([]byte, 32)
		copy(revHash, digestHash)
		reverse(revHash)
		ver, err := pub.VerifyDigest(revHash, signature[:64])
		if err != nil {
			return false
		}
		return ver
	default:
		panic("cant infer pub key type")
	}
}

// DecompressPubkey parses a public key in the 33-byte compressed format
// TODO GOST decompression
func DecompressPubkey[P PublicKey](pubkey []byte) (P, error) {
	var pubKey P
	switch p := any(&pubKey).(type) {
	case *nist.PublicKey:
		x, y := secp256k1.DecompressPubkey(pubkey)
		if x == nil {
			return ZeroPublicKey[P](), fmt.Errorf("invalid public key")
		}
		*p = nist.PublicKey{&ecdsa.PublicKey{X: x, Y: y, Curve: S256()}}
	case *gost3410.PublicKey:
		k, err := gost3410.NewPublicKey(gost3410.GostCurve, pubkey[1:])
		if err != nil {
			return ZeroPublicKey[P](), fmt.Errorf("invalid gost 3410 public key")
		}
		*p=*k
	default:
		return ZeroPublicKey[P](), fmt.Errorf("cant infer pub key type")
	}
	return pubKey, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
// TODO GOST compression
func CompressPubkey[P PublicKey](pubkey P) []byte {
	switch p := any(&pubkey).(type) {
	case *nist.PublicKey:
		return secp256k1.CompressPubkey(p.X, p.Y)
	case *gost3410.PublicKey:
		return gost3410.Marshal(gost3410.GostCurve, p.X, p.Y)
	default:
		return nil
	}

}

func RevertCSP(hash, signature []byte) (r, s *big.Int, revertHash []byte) {
	revertHash = make([]byte, 32)
	copy(revertHash, hash)
	reverse(revertHash)

	sig := make([]byte, 64)
	copy(sig, signature)
	reverse(sig)
	s = new(big.Int).SetBytes(sig[:32])
	r = new(big.Int).SetBytes(sig[32:64])

	return
}

func reverse(d []byte) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}