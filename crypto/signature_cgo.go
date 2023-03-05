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
	"github.com/pavelkrolevets/MIR-pro/crypto/csp"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/secp256k1"
)

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(digestHash, sig []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(digestHash, sig)
}

func EcrecoverGost(digestHash, sig []byte) ([]byte, error) {
	v := int((sig[64]) & ^byte(4))
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	X, Y, err := gost3410.RecoverCompact(*gost3410.CurveIdGostR34102001CryptoProAParamSet(), digestHash, r, s, v)
	pubKey := gost3410.PublicKey{
		C: gost3410.GostCurve,
		X: X,
		Y: Y,
	}
	if err != nil{
		return nil, err
	}
	return gost3410.Marshal(*gost3410.GostCurve, pubKey.X, pubKey.Y), nil
}

func EcrecoverCsp(digestHash, sig []byte) ([]byte, error) {
	hash, err := csp.NewHash(csp.HashOptions{HashAlg: csp.GOST_R3411_12_256})
	if err != nil {
		return nil, err
	}
	defer hash.Close()
	_, err = hash.Write(digestHash)
	if err != nil {
		return nil, err
	}
	digest := hash.Sum(nil)
	r, s, revHash := RevertCSP(digest, sig)
	pub := make([]byte, 0, 64)
	X, Y, err := gost3410.RecoverCompact(*gost3410.GostCurve, revHash, r, s, int(sig[64]))
	if err != nil{
		return nil, err
	}
	pub = append(pub, Y.Bytes()...)
	pub = append(pub, X.Bytes()...)
	reverse(pub)
	return pub, nil
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	x, y := elliptic.Unmarshal(S256(), s)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}
func SigToPubGost(hash, sig []byte) (*gost3410.PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}
	pubKey, err := gost3410.NewPublicKey(gost3410.GostCurve, s)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func SigToPubCsp(hash, sig []byte) (interface{}, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}
	store, err := csp.SystemStore("My")
	if err != nil {
		return nil, fmt.Errorf("store error: %s", err)
	}
	defer store.Close()

	return s, nil
}

// Sign calculates an ECDSA signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given digest cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(digestHash) != DigestLength {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", DigestLength, len(digestHash))
	}
	seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
	defer zeroBytes(seckey)
	return secp256k1.Sign(digestHash, seckey)
}

func SignGost(digestHash []byte, prv *gost3410.PrivateKey) (resSig []byte, err error) {
	if len(digestHash) != DigestLength {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", DigestLength, len(digestHash))
	}
	sig, err := prv.SignDigest(digestHash, rand.Reader)
	if err != nil {
		return nil, err
	}
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	for i := 0; i < (1+1)*2; i++ {
		X, Y, err := gost3410.RecoverCompact(*prv.C, digestHash, r, s, i)
		pub := prv.PublicKey()
		if err == nil && X.Cmp(pub.X) == 0 && Y.Cmp(pub.Y) == 0 {
			resSig = append(resSig, sig...)
			resSig = append(resSig, byte(i))
		}
	}
	return resSig, nil
}

func SignCsp(digestHash []byte, crt *csp.Cert) (resSig []byte, err error) {
	defer crt.Close()
	hash, err := csp.NewHash(csp.HashOptions{SignCert: *crt, HashAlg: csp.GOST_R3411_12_256})
	if err != nil {
		return nil, err
	}
	defer hash.Close()
	_, err = hash.Write(digestHash)
	if err != nil {
		return nil, err
	}
	digest := hash.Sum(nil)
	sig, err := hash.Sign()
	if err != nil {
		return nil, err
	}
	r, s, revHash := RevertCSP(digest, sig)
	pubKey := make([]byte, 64)
	copy(pubKey[:], crt.Info().PublicKeyBytes()[2:66])
	reverse(pubKey)
	pubKeyX := new(big.Int).SetBytes(pubKey[32 : 64])
	pubKeyY := new(big.Int).SetBytes(pubKey[ : 32])
	gostKey := gost3410.PublicKey{
		C: gost3410.CurveIdGostR34102001CryptoProAParamSet(),
		X: pubKeyX,
		Y: pubKeyY,
	}
	for i := 0; i < 4; i++ {
		X, Y, err := gost3410.RecoverCompact(*gost3410.CurveIdGostR34102001CryptoProAParamSet(), revHash, r, s, i)
		if err == nil && X.Cmp(gostKey.X) == 0 && Y.Cmp(gostKey.Y) == 0 {
			resSig = append(resSig, sig...)
			resSig = append(resSig, byte(i))
		}
	}
	return resSig, nil
}

// VerifySignature checks that the given public key created signature over digest.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature(pubkey, digestHash, signature []byte) bool {
	return secp256k1.VerifySignature(pubkey, digestHash, signature)
}

func VerifySignatureGost(pubkey, digestHash, signature []byte) bool {
		pub, err := gost3410.NewPublicKey(gost3410.GostCurve, pubkey)
		if err != nil {
			return false
		}
		ver, err := pub.VerifyDigest(digestHash, signature[:128])
		if err != nil {
			return false
		}
		return ver
}

func VerifySignatureScp(pubkey, digestHash, signature []byte) bool {
	res, err := csp.VerifySignatureRaw(digestHash, signature, pubkey)
	if err != nil {
		return false
	}
	return res
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	x, y := secp256k1.DecompressPubkey(pubkey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{X: x, Y: y, Curve: S256()}, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	return secp256k1.CompressPubkey(pubkey.X, pubkey.Y)
}

// S256 returns an instance of the secp256k1 curve.
func S256() elliptic.Curve {
	return secp256k1.S256()
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