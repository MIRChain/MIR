// Copyright (c) 2013 Kyle Isom <kyle@tyrfingr.is>
// Copyright (c) 2012 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package ecies

import (
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/subtle"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"math/big"

	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/nist"
)

var (
	ErrImport                     = fmt.Errorf("ecies: failed to import key")
	ErrInvalidCurve               = fmt.Errorf("ecies: invalid elliptic curve")
	ErrInvalidPublicKey           = fmt.Errorf("ecies: invalid public key")
	ErrSharedKeyIsPointAtInfinity = fmt.Errorf("ecies: shared key is point at infinity")
	ErrSharedKeyTooBig            = fmt.Errorf("ecies: shared key params are too big")
)

// PublicKey is a representation of an elliptic curve public key.
type PublicKey[P crypto.PublicKey] struct {
	X *big.Int
	Y *big.Int
	elliptic.Curve
	Params *ECIESParams
}

// Export an ECIES public key as an ECDSA public key.
func (e *PublicKey[P]) ExportECDSA() P {
	var pub P
	switch p := any(&pub).(type) {
	case *nist.PublicKey:
		*p = nist.PublicKey{&ecdsa.PublicKey{Curve: e.Curve, X: e.X, Y: e.Y}}
	case *gost3410.PublicKey:
		*p = gost3410.PublicKey{C: e.Curve, X: e.X, Y: e.Y}
	}
	return pub
}

// Import an ECDSA public key as an ECIES public key.
func ImportECDSAPublic[P crypto.PublicKey](pub P) *PublicKey[P] {
	switch p := any(&pub).(type) {
	case *nist.PublicKey:
		return &PublicKey[P]{
			X:      p.X,
			Y:      p.Y,
			Curve:  p.Curve,
			Params: ParamsFromCurve(p.Curve),
		}
	case *gost3410.PublicKey:
		return &PublicKey[P]{
			X:      p.X,
			Y:      p.Y,
			Curve:  p.C,
			Params: ECIES_AES128_SHA256,
		}
	default:
		panic("cant infer pub key")
	}
}

// PrivateKey is a representation of an elliptic curve private key.
type PrivateKey[T crypto.PrivateKey, P crypto.PublicKey] struct {
	PublicKey[P]
	D *big.Int
}

// Export an ECIES private key as an ECDSA private key.
func (prv *PrivateKey[T, P]) ExportECDSA() T {
	pub := &prv.PublicKey
	pubECDSA := pub.ExportECDSA()
	var privGeneric T
	switch t := any(&privGeneric).(type) {
	case *nist.PrivateKey:
		p := any(&pubECDSA).(*nist.PublicKey)
		*t = nist.PrivateKey{&ecdsa.PrivateKey{PublicKey: *p.PublicKey, D: prv.D}}
	case *gost3410.PrivateKey:
		p := any(&pubECDSA).(*gost3410.PublicKey)
		*t = gost3410.PrivateKey{PublicKey: *p, C: gost3410.GostCurve, Key: prv.D}
	}
	return privGeneric
}

// Import an ECDSA private key as an ECIES private key.
func ImportECDSA[T crypto.PrivateKey, P crypto.PublicKey](prv T) *PrivateKey[T, P] {
	switch t := any(&prv).(type) {
	case *nist.PrivateKey:
		var pub P
		p := any(&pub).(*nist.PublicKey)
		*p = *t.Public()
		pubKey := ImportECDSAPublic[P](pub)
		return &PrivateKey[T, P]{*pubKey, t.D}
	case *gost3410.PrivateKey:
		var pub P
		p := any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
		pubKey := ImportECDSAPublic[P](pub)
		return &PrivateKey[T, P]{*pubKey, t.Key}
	default:
		panic("cant infer priv key for ecies")
	}
}

// Generate an elliptic curve public / private keypair. If params is nil,
// the recommended default parameters for the key will be chosen.
func GenerateKey[T crypto.PrivateKey, P crypto.PublicKey](rand io.Reader, curve elliptic.Curve, params *ECIESParams) (*PrivateKey[T, P], error) {
	var privKey T
	switch any(&privKey).(type) {
	case *nist.PrivateKey:
		pb, x, y, err := elliptic.GenerateKey(curve, rand)
		if err != nil {
			return nil, err
		}
		prv := new(PrivateKey[T, P])
		prv.PublicKey.X = x
		prv.PublicKey.Y = y
		prv.PublicKey.Curve = curve
		prv.D = new(big.Int).SetBytes(pb)
		if params == nil {
			params = ParamsFromCurve(curve)
		}
		prv.PublicKey.Params = params
		return prv, nil
	case *gost3410.PrivateKey:
		key, err := gost3410.GenPrivateKey(gost3410.GostCurve, rand)
		if err != nil {
			return nil, err
		}
		prv := new(PrivateKey[T, P])
		prv.PublicKey.X = key.X
		prv.PublicKey.Y = key.Y
		prv.PublicKey.Curve = key.C
		prv.D = key.Key
		if params == nil {
			params = ECIES_AES128_SHA256
		}
		prv.PublicKey.Params = params
		return prv, nil
	default:
		panic("cant infer priv key for ecies")
	}
}

// MaxSharedKeyLength returns the maximum length of the shared key the
// public key can produce.
func MaxSharedKeyLength[P crypto.PublicKey](pub *PublicKey[P]) int {
	return (pub.Curve.Params().BitSize + 7) / 8
}

// ECDH key agreement method used to establish secret keys for encryption.
func (prv *PrivateKey[T, P]) GenerateShared(pub *PublicKey[P], skLen, macLen int) (sk []byte, err error) {
	if prv.PublicKey.Curve != pub.Curve {
		return nil, ErrInvalidCurve
	}
	if skLen+macLen > MaxSharedKeyLength(pub) {
		return nil, ErrSharedKeyTooBig
	}

	x, _ := pub.Curve.ScalarMult(pub.X, pub.Y, prv.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}

var (
	ErrSharedTooLong  = fmt.Errorf("ecies: shared secret is too long")
	ErrInvalidMessage = fmt.Errorf("ecies: invalid message")
)

// NIST SP 800-56 Concatenation Key Derivation Function (see section 5.8.1).
func concatKDF(hash hash.Hash, z, s1 []byte, kdLen int) []byte {
	counterBytes := make([]byte, 4)
	k := make([]byte, 0, roundup(kdLen, hash.Size()))
	for counter := uint32(1); len(k) < kdLen; counter++ {
		binary.BigEndian.PutUint32(counterBytes, counter)
		hash.Reset()
		hash.Write(counterBytes)
		hash.Write(z)
		hash.Write(s1)
		k = hash.Sum(k)
	}
	return k[:kdLen]
}

// roundup rounds size up to the next multiple of blocksize.
func roundup(size, blocksize int) int {
	return size + blocksize - (size % blocksize)
}

// deriveKeys creates the encryption and MAC keys using concatKDF.
func deriveKeys(hash hash.Hash, z, s1 []byte, keyLen int) (Ke, Km []byte) {
	K := concatKDF(hash, z, s1, 2*keyLen)
	Ke = K[:keyLen]
	Km = K[keyLen:]
	hash.Reset()
	hash.Write(Km)
	Km = hash.Sum(Km[:0])
	return Ke, Km
}

// messageTag computes the MAC of a message (called the tag) as per
// SEC 1, 3.5.
func messageTag(hash func() hash.Hash, km, msg, shared []byte) []byte {
	mac := hmac.New(hash, km)
	mac.Write(msg)
	mac.Write(shared)
	tag := mac.Sum(nil)
	return tag
}

// Generate an initialisation vector for CTR mode.
func generateIV(params *ECIESParams, rand io.Reader) (iv []byte, err error) {
	iv = make([]byte, params.BlockSize)
	_, err = io.ReadFull(rand, iv)
	return
}

// symEncrypt carries out CTR encryption using the block cipher specified in the
func symEncrypt(rand io.Reader, params *ECIESParams, key, m []byte) (ct []byte, err error) {
	c, err := params.Cipher(key)
	if err != nil {
		return
	}

	iv, err := generateIV(params, rand)
	if err != nil {
		return
	}
	ctr := cipher.NewCTR(c, iv)

	ct = make([]byte, len(m)+params.BlockSize)
	copy(ct, iv)
	ctr.XORKeyStream(ct[params.BlockSize:], m)
	return
}

// symDecrypt carries out CTR decryption using the block cipher specified in
// the parameters
func symDecrypt(params *ECIESParams, key, ct []byte) (m []byte, err error) {
	c, err := params.Cipher(key)
	if err != nil {
		return
	}

	ctr := cipher.NewCTR(c, ct[:params.BlockSize])

	m = make([]byte, len(ct)-params.BlockSize)
	ctr.XORKeyStream(m, ct[params.BlockSize:])
	return
}

// Encrypt encrypts a message using ECIES as specified in SEC 1, 5.1.
//
// s1 and s2 contain shared information that is not part of the resulting
// ciphertext. s1 is fed into key derivation, s2 is fed into the MAC. If the
// shared information parameters aren't being used, they should be nil.
func Encrypt[T crypto.PrivateKey, P crypto.PublicKey](rand io.Reader, pub *PublicKey[P], m, s1, s2 []byte) (ct []byte, err error) {
	var params *ECIESParams
	var pubType P
	switch any(&pubType).(type) {
	case *nist.PublicKey:
		params, err = pubkeyParams(pub)
		if err != nil {
			return nil, err
		}
	case *gost3410.PublicKey:
		params = ECIES_AES128_SHA256
	}
	R, err := GenerateKey[T, P](rand, pub.Curve, params)
	if err != nil {
		return nil, err
	}

	z, err := R.GenerateShared(pub, params.KeyLen, params.KeyLen)
	if err != nil {
		return nil, err
	}

	hash := params.Hash()
	Ke, Km := deriveKeys(hash, z, s1, params.KeyLen)

	em, err := symEncrypt(rand, params, Ke, m)
	if err != nil || len(em) <= params.BlockSize {
		return nil, err
	}

	d := messageTag(params.Hash, Km, em, s2)

	var Rb []byte
	switch any(&pubType).(type) {
	case *nist.PublicKey:
		Rb = elliptic.Marshal(pub.Curve, R.PublicKey.X, R.PublicKey.Y)
	case *gost3410.PublicKey:
		Rb = gost3410.Marshal(pub.Curve, R.PublicKey.X, R.PublicKey.Y)
	}
	ct = make([]byte, len(Rb)+len(em)+len(d))
	copy(ct, Rb)
	copy(ct[len(Rb):], em)
	copy(ct[len(Rb)+len(em):], d)
	return ct, nil
}

// Decrypt decrypts an ECIES ciphertext.
func (prv *PrivateKey[T, P]) Decrypt(c, s1, s2 []byte) (m []byte, err error) {
	if len(c) == 0 {
		return nil, ErrInvalidMessage
	}
	var params *ECIESParams
	var pubType P
	switch any(&pubType).(type) {
	case *nist.PublicKey:
		params, err = pubkeyParams(&prv.PublicKey)
		if err != nil {
			return nil, err
		}
	case *gost3410.PublicKey:
		params = ECIES_AES128_SHA256
	}
	hash := params.Hash()

	var (
		rLen   int
		hLen   int = hash.Size()
		mStart int
		mEnd   int
	)

	switch c[0] {
	case 2, 3, 4:
		rLen = (prv.PublicKey.Curve.Params().BitSize + 7) / 4
		if len(c) < (rLen + hLen + 1) {
			return nil, ErrInvalidMessage
		}
	default:
		return nil, ErrInvalidPublicKey
	}

	mStart = rLen
	mEnd = len(c) - hLen

	R := new(PublicKey[P])
	R.Curve = prv.PublicKey.Curve
	switch any(&pubType).(type) {
	case *nist.PublicKey:
		R.X, R.Y = elliptic.Unmarshal(R.Curve, c[:rLen])
	case *gost3410.PublicKey:
		R.X, R.Y = gost3410.Unmarshal(R.Curve, c[:rLen])
	}
	if R.X == nil {
		return nil, ErrInvalidPublicKey
	}

	z, err := prv.GenerateShared(R, params.KeyLen, params.KeyLen)
	if err != nil {
		return nil, err
	}
	Ke, Km := deriveKeys(hash, z, s1, params.KeyLen)

	d := messageTag(params.Hash, Km, c[mStart:mEnd], s2)
	if subtle.ConstantTimeCompare(c[mEnd:], d) != 1 {
		return nil, ErrInvalidMessage
	}

	return symDecrypt(params, Ke, c[mStart:mEnd])
}
