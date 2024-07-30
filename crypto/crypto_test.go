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
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/common/hexutil"
	"github.com/MIRChain/MIR/crypto/csp"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/gost3411"
	"github.com/MIRChain/MIR/crypto/nist"
)

var testAddrHex = "970e8128ab834e8eac17ab8e3812f010678cf791"
var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"

// These tests are sanity checks.
// They should ensure that we don't e.g. use Sha3-224 instead of Sha3-256
// and that the sha3 library uses keccak-f permutation.
func TestKeccak256Hash(t *testing.T) {
	msg := []byte("abc")
	exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := Keccak256Hash[nist.PublicKey](in); return h[:] }, msg, exp)
}

func Test3411Strebog(t *testing.T) {
	msg := []byte{
		0xd1, 0xe5, 0x20, 0xe2, 0xe5, 0xf2, 0xf0, 0xe8,
		0x2c, 0x20, 0xd1, 0xf2, 0xf0, 0xe8, 0xe1, 0xee,
		0xe6, 0xe8, 0x20, 0xe2, 0xed, 0xf3, 0xf6, 0xe8,
		0x2c, 0x20, 0xe2, 0xe5, 0xfe, 0xf2, 0xfa, 0x20,
		0xf1, 0x20, 0xec, 0xee, 0xf0, 0xff, 0x20, 0xf1,
		0xf2, 0xf0, 0xe5, 0xeb, 0xe0, 0xec, 0xe8, 0x20,
		0xed, 0xe0, 0x20, 0xf5, 0xf0, 0xe0, 0xe1, 0xf0,
		0xfb, 0xff, 0x20, 0xef, 0xeb, 0xfa, 0xea, 0xfb,
		0x20, 0xc8, 0xe3, 0xee, 0xf0, 0xe5, 0xe2, 0xfb,
	}
	exp := []byte{
		0x9d, 0xd2, 0xfe, 0x4e, 0x90, 0x40, 0x9e, 0x5d,
		0xa8, 0x7f, 0x53, 0x97, 0x6d, 0x74, 0x05, 0xb0,
		0xc0, 0xca, 0xc6, 0x28, 0xfc, 0x66, 0x9a, 0x74,
		0x1d, 0x50, 0x06, 0x3c, 0x55, 0x7e, 0x8f, 0x50,
	}
	checkhash(t, "Streebog-256-array", func(in []byte) []byte { h := Keccak256Hash[gost3410.PublicKey](in); return h[:] }, msg, exp)
}

func TestKeccak256Hasher(t *testing.T) {
	msg := []byte("abc")
	exp, _ := hex.DecodeString("4e03657aea45a94fc7d47ba826c8d667c0d1e6e33a64a036ec44f58fa12d6c45")
	hasher := NewKeccakState[nist.PublicKey]()
	checkhash(t, "Sha3-256-array", func(in []byte) []byte { h := HashData[nist.PublicKey](hasher, in); return h[:] }, msg, exp)
}

func Test3411StreebogHasher(t *testing.T) {
	msg := []byte{
		0xd1, 0xe5, 0x20, 0xe2, 0xe5, 0xf2, 0xf0, 0xe8,
		0x2c, 0x20, 0xd1, 0xf2, 0xf0, 0xe8, 0xe1, 0xee,
		0xe6, 0xe8, 0x20, 0xe2, 0xed, 0xf3, 0xf6, 0xe8,
		0x2c, 0x20, 0xe2, 0xe5, 0xfe, 0xf2, 0xfa, 0x20,
		0xf1, 0x20, 0xec, 0xee, 0xf0, 0xff, 0x20, 0xf1,
		0xf2, 0xf0, 0xe5, 0xeb, 0xe0, 0xec, 0xe8, 0x20,
		0xed, 0xe0, 0x20, 0xf5, 0xf0, 0xe0, 0xe1, 0xf0,
		0xfb, 0xff, 0x20, 0xef, 0xeb, 0xfa, 0xea, 0xfb,
		0x20, 0xc8, 0xe3, 0xee, 0xf0, 0xe5, 0xe2, 0xfb,
	}
	exp := []byte{
		0x9d, 0xd2, 0xfe, 0x4e, 0x90, 0x40, 0x9e, 0x5d,
		0xa8, 0x7f, 0x53, 0x97, 0x6d, 0x74, 0x05, 0xb0,
		0xc0, 0xca, 0xc6, 0x28, 0xfc, 0x66, 0x9a, 0x74,
		0x1d, 0x50, 0x06, 0x3c, 0x55, 0x7e, 0x8f, 0x50,
	}
	hasher := NewKeccakState[gost3410.PublicKey]()
	checkhash(t, "Streebog-256-array", func(in []byte) []byte { h := HashData[gost3410.PublicKey](hasher, in); return h[:] }, msg, exp)
}

func Test3411Streebog512Hasher(t *testing.T) {
	msg := []byte{
		0xd1, 0xe5, 0x20, 0xe2, 0xe5, 0xf2, 0xf0, 0xe8,
		0x2c, 0x20, 0xd1, 0xf2, 0xf0, 0xe8, 0xe1, 0xee,
		0xe6, 0xe8, 0x20, 0xe2, 0xed, 0xf3, 0xf6, 0xe8,
		0x2c, 0x20, 0xe2, 0xe5, 0xfe, 0xf2, 0xfa, 0x20,
		0xf1, 0x20, 0xec, 0xee, 0xf0, 0xff, 0x20, 0xf1,
		0xf2, 0xf0, 0xe5, 0xeb, 0xe0, 0xec, 0xe8, 0x20,
		0xed, 0xe0, 0x20, 0xf5, 0xf0, 0xe0, 0xe1, 0xf0,
		0xfb, 0xff, 0x20, 0xef, 0xeb, 0xfa, 0xea, 0xfb,
		0x20, 0xc8, 0xe3, 0xee, 0xf0, 0xe5, 0xe2, 0xfb,
	}
	exp := []byte{
		0x1e, 0x88, 0xe6, 0x22, 0x26, 0xbf, 0xca, 0x6f,
		0x99, 0x94, 0xf1, 0xf2, 0xd5, 0x15, 0x69, 0xe0,
		0xda, 0xf8, 0x47, 0x5a, 0x3b, 0x0f, 0xe6, 0x1a,
		0x53, 0x00, 0xee, 0xe4, 0x6d, 0x96, 0x13, 0x76,
		0x03, 0x5f, 0xe8, 0x35, 0x49, 0xad, 0xa2, 0xb8,
		0x62, 0x0f, 0xcd, 0x7c, 0x49, 0x6c, 0xe5, 0xb3,
		0x3f, 0x0c, 0xb9, 0xdd, 0xdc, 0x2b, 0x64, 0x60,
		0x14, 0x3b, 0x03, 0xda, 0xba, 0xc9, 0xfb, 0x28,
	}
	checkhash(t, "Streebog-512-array", func(in []byte) []byte { h := Keccak512[gost3410.PublicKey](in); return h[:] }, msg, exp)
}

func TestToECDSAErrors(t *testing.T) {
	if _, err := HexToECDSA[nist.PrivateKey]("0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Fatal("HexToECDSA should've returned error")
	}
	if _, err := HexToECDSA[nist.PrivateKey]("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"); err == nil {
		t.Fatal("HexToECDSA should've returned error")
	}
}

func BenchmarkSha3(b *testing.B) {
	a := []byte("hello world")
	for i := 0; i < b.N; i++ {
		Keccak256[nist.PublicKey](a)
	}
}

func Benchmark3411Streebog(b *testing.B) {
	a := []byte("hello world")
	for i := 0; i < b.N; i++ {
		Keccak256[gost3410.PublicKey](a)
	}
}

func TestUnmarshalPubkey(t *testing.T) {
	key, err := UnmarshalPubkey[nist.PublicKey](nil)
	if err != errInvalidPubkey || key.PublicKey != nil {
		t.Fatalf("expected error, got %v, %v", err, key)
	}
	key, err = UnmarshalPubkey[nist.PublicKey]([]byte{1, 2, 3})
	if err != errInvalidPubkey || key.PublicKey != nil {
		t.Fatalf("expected error, got %v, %v", err, key)
	}

	var (
		enc, _ = hex.DecodeString("04760c4460e5336ac9bbd87952a3c7ec4363fc0a97bd31c86430806e287b437fd1b01abc6e1db640cf3106b520344af1d58b00b57823db3e1407cbc433e1b6d04d")
		dec    = nist.PublicKey{&ecdsa.PublicKey{
			Curve: S256(),
			X:     hexutil.MustDecodeBig("0x760c4460e5336ac9bbd87952a3c7ec4363fc0a97bd31c86430806e287b437fd1"),
			Y:     hexutil.MustDecodeBig("0xb01abc6e1db640cf3106b520344af1d58b00b57823db3e1407cbc433e1b6d04d"),
		}}
	)
	key, err = UnmarshalPubkey[nist.PublicKey](enc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(key, dec) {
		t.Fatal("wrong result")
	}
}

func TestUnmarshalPubkeyGost(t *testing.T) {
	gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
	var (
		enc, _ = hex.DecodeString("04e4f910baf0152b2bac365a5ac2323323dd48a46db08c30c8b0cd140154dbc4218496a69e003c2b5eb14ec23b7e0a83e9212d33500b7764f74a06ad92be36e775")
		dec    = gost3410.PublicKey{
			C: gost3410.GostCurve,
			X: hexutil.MustDecodeBig("0xe4f910baf0152b2bac365a5ac2323323dd48a46db08c30c8b0cd140154dbc421"),
			Y: hexutil.MustDecodeBig("0x8496a69e003c2b5eb14ec23b7e0a83e9212d33500b7764f74a06ad92be36e775"),
		}
	)
	key, err := UnmarshalPubkey[gost3410.PublicKey](enc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(key, dec) {
		t.Fatalf("wrong result: want %d, have %d", dec.X, key.X)
	}
}

func TestSign(t *testing.T) {
	CryptoAlg = NIST
	key, _ := HexToECDSA[nist.PrivateKey](testPrivHex)
	addr := common.HexToAddress(testAddrHex)

	msg := Keccak256[nist.PublicKey]([]byte("foo"))
	sig, err := Sign(msg, key)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}
	recoveredPub, err := Ecrecover[nist.PublicKey](msg, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}
	pubKey, _ := UnmarshalPubkey[nist.PublicKey](recoveredPub)
	recoveredAddr := PubkeyToAddress(pubKey)
	if addr != recoveredAddr {
		t.Errorf("Address mismatch: want: %x have: %x", addr, recoveredAddr)
	}

	// should be equal to SigToPub
	recoveredPub2, err := SigToPub[nist.PublicKey](msg, sig)
	if err != nil {
		t.Errorf("ECRecover error: %s", err)
	}
	recoveredAddr2 := PubkeyToAddress(recoveredPub2)
	if addr != recoveredAddr2 {
		t.Errorf("Address mismatch: want: %x have: %x", addr, recoveredAddr2)
	}

	CryptoAlg = GOST
	gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
	gostKey, _ := gost3410.GenPrivateKey(gost3410.CurveIdGostR34102001CryptoProAParamSet(), rand.Reader)
	gostMsg := gost3411.New(32)
	gostMsg.Write(([]byte("foo")))
	digest := make([]byte, 32)
	gostMsg.Read(digest)
	gostSig, err := Sign(digest, *gostKey)
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}
	revHash := make([]byte, 32)
	copy(revHash, digest)
	reverse(revHash)
	ver, err := gostKey.Public().VerifyDigest(revHash, gostSig[:64])
	if err != nil {
		t.Errorf("Sign error: %s", err)
	}
	assert.Equal(t, true, ver)
	ver = VerifySignature[gost3410.PublicKey](FromECDSAPub(*gostKey.Public()), digest, gostSig)
	assert.Equal(t, true, ver)
	r := new(big.Int).SetBytes(gostSig[:32])
	s := new(big.Int).SetBytes(gostSig[32:64])
	X, Y, err := gost3410.RecoverCompact(*gost3410.GostCurve, revHash, r, s, int((gostSig[64]) & ^byte(4)))
	if err == nil && X.Cmp(gostKey.Public().X) == 0 && Y.Cmp(gostKey.Public().Y) == 0 {
		t.Log("Recovered X ", X.String())
		t.Log("Recovered Y ", Y.String())
	}

	recoveredGostPub, err := Ecrecover[gost3410.PublicKey](digest, gostSig)
	if err != nil {
		t.Fatalf("ECRecover error: %s", err)
	}
	if !bytes.Equal(gostKey.Public().X.Bytes(), recoveredGostPub[1:33]) {
		t.Fatalf("Address mismatch: want: %x have: %x", gostKey.Public().X.Bytes(), recoveredGostPub[1:33])
	}
	if !bytes.Equal(gostKey.Public().Y.Bytes(), recoveredGostPub[33:65]) {
		t.Fatalf("Address mismatch: want: %x have: %x", gostKey.Public().Y.Bytes(), recoveredGostPub[33:65])
	}

	if !bytes.Equal(gostKey.Public().Raw(), recoveredGostPub[1:]) {
		t.Fatalf("Address mismatch: want: %x have: %x", gostKey.Public().Raw(), recoveredGostPub[1:])
	}

	CryptoAlg = GOST_CSP
	store, err := csp.SystemStore("My")
	if err != nil {
		t.Fatalf("Store error: %s", err)
	}
	defer store.Close()
	// Cert should be without set pin
	crt, err := store.GetBySubjectId("4ac93fc08bc0efd24180b0fa47f7309c257e8c85")
	if err != nil {
		t.Fatalf("Get cert error: %s", err)
	}
	t.Logf("Cert pub key: %x", crt.Info().PublicKeyBytes()[:66])
	defer crt.Close()
	hash, err := csp.NewHash(csp.HashOptions{SignCert: crt, HashAlg: csp.GOST_R3411_12_256})
	if err != nil {
		t.Fatal(err)
	}
	_, err = hash.Write([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	digest = hash.Sum(nil)
	hash.Reset()
	hash.Close()
	t.Logf("hash digest: %x", digest)
	sig, err = Sign(digest, crt)
	if err != nil {
		t.Fatalf("Sign error: %s", err)
	}
	t.Log("Sig csp", len(sig))
	recoveredGostPub, err = Ecrecover[csp.PublicKey](digest, sig)
	if err != nil {
		t.Fatalf("ECRecover error: %s", err)
	}
	if !bytes.Equal(crt.Info().PublicKeyBytes()[1:66], recoveredGostPub) {
		t.Fatalf("Address mismatch: want: %x have: %x", crt.Info().PublicKeyBytes()[1:66], recoveredGostPub)
	}
}

func TestHashCSP(t *testing.T) {
	hash, err := csp.NewHash(csp.HashOptions{HashAlg: csp.GOST_R3411_12_256})
	if err != nil {
		t.Fatal(err)
	}
	_, err = hash.Write([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	digest := hash.Sum(nil)
	hash.Reset()
	hash.Close()
	t.Logf("hash digest: %x", digest)
}
func TestHashReadCSP(t *testing.T) {
	hasher := NewKeccakState[csp.PublicKey]()
	_, err := hasher.Write([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	digestNew := make([]byte, 32)
	hasher.Read(digestNew)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Read hash digest: %x", digestNew)
}
func TestSignCSPRecoverGOST(t *testing.T) {
	store, err := csp.SystemStore("My")
	if err != nil {
		t.Fatalf("Store error: %s", err)
	}
	defer store.Close()
	// Cert should be without set pin
	crt, err := store.GetBySubjectId("4ac93fc08bc0efd24180b0fa47f7309c257e8c85")
	if err != nil {
		t.Fatalf("Get cert error: %s", err)
	}
	defer crt.Close()
	hash, err := csp.NewHash(csp.HashOptions{SignCert: crt, HashAlg: csp.GOST_R3411_12_256})
	if err != nil {
		t.Fatal(err)
	}
	_, err = hash.Write([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	digest := hash.Sum(nil)
	hash.Reset()
	hash.Close()
	t.Logf("hash digest: %x", digest)
	sig, err := Sign(digest, crt)
	if err != nil {
		t.Fatalf("Sign error: %s", err)
	}
	t.Log("Sig csp", len(sig))
	recoveredGostPub, err := Ecrecover[csp.PublicKey](digest, sig)
	if err != nil {
		t.Fatalf("ECRecover error: %s", err)
	}
	if !bytes.Equal(crt.Info().PublicKeyBytes()[1:66], recoveredGostPub) {
		t.Fatalf("Address mismatch: want: %x have: %x", crt.Info().PublicKeyBytes()[1:66], recoveredGostPub)
	}
	// Trying to Ecrecover from CSP sig to pure GOST pub key
	// 1. Reverse and convert CSP sig to 65 bytes [r,s]+1
	revPub := make([]byte, 65)
	revPub[0] = 4
	cspPub := crt.Info().PublicKeyBytes()[2:66]
	reverse(cspPub)
	copy(revPub[1:33], cspPub[32:64])
	copy(revPub[33:65], cspPub[:32])
	if crt.Public().X.Cmp(new(big.Int).SetBytes(revPub[1:33])) != 0 {
		t.Fatalf("Pub key mismatch: want: %d have: %d", crt.Public().X, new(big.Int).SetBytes(revPub[1:33]))
	}
	revSig := make([]byte, 65)
	r, s, _ := RevertCSP(digest, sig)
	copy(revSig[:32], r.Bytes())
	copy(revSig[32:64], s.Bytes())
	revSig[64] = sig[64]
	// Verify csp sig using pure gost crypto
	ver := VerifySignature[gost3410.PublicKey](revPub, digest, revSig)
	assert.Equal(t, true, ver)

	recoveredGostPub, err = Ecrecover[gost3410.PublicKey](digest, revSig)
	if err != nil {
		t.Fatalf("ECRecover error: %s", err)
	}
	resPub := recoveredGostPub[1:]
	pub := make([]byte, 64)
	copy(pub[:32], resPub[32:64])
	copy(pub[32:64], resPub[:32])
	reverse(pub)
	if !bytes.Equal(crt.Info().PublicKeyBytes()[2:66], pub) {
		t.Fatalf("Pub key: want: %x have: %x", crt.Info().PublicKeyBytes()[2:66], pub)
	}
}

// func TestInvalidSign(t *testing.T) {
// 	if _, err := Sign(make([]byte, 1), nil); err == nil {
// 		t.Errorf("expected sign with hash 1 byte to error")
// 	}
// 	if _, err := Sign(make([]byte, 33), nil); err == nil {
// 		t.Errorf("expected sign with hash 33 byte to error")
// 	}
// }

func TestNewContractAddress(t *testing.T) {
	key, _ := HexToECDSA[nist.PrivateKey](testPrivHex)
	addr := common.HexToAddress(testAddrHex)
	genAddr := PubkeyToAddress(*key.Public())
	// sanity check before using addr to create contract address
	checkAddr(t, genAddr, addr)

	caddr0 := CreateAddress[nist.PublicKey](addr, 0)
	caddr1 := CreateAddress[nist.PublicKey](addr, 1)
	caddr2 := CreateAddress[nist.PublicKey](addr, 2)
	checkAddr(t, common.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d"), caddr0)
	checkAddr(t, common.HexToAddress("8bda78331c916a08481428e4b07c96d3e916d165"), caddr1)
	checkAddr(t, common.HexToAddress("c9ddedf451bc62ce88bf9292afb13df35b670699"), caddr2)
}

func TestLoadECDSA(t *testing.T) {
	tests := []struct {
		input string
		err   string
	}{
		// good
		{input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		{input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n"},
		{input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n\r"},
		{input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\r\n"},
		{input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n\n"},
		{input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n\r"},
		// bad
		{
			input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde",
			err:   "key file too short, want 64 hex characters",
		},
		{
			input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde\n",
			err:   "key file too short, want 64 hex characters",
		},
		{
			input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeX",
			err:   "invalid hex character 'X' in private key",
		},
		{
			input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdefX",
			err:   "invalid character 'X' at end of key file",
		},
		{
			input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n\n\n",
			err:   "key file too long, want 64 hex characters",
		},
	}

	for _, test := range tests {
		f, err := ioutil.TempFile("", "loadecdsa_test.*.txt")
		if err != nil {
			t.Fatal(err)
		}
		filename := f.Name()
		f.WriteString(test.input)
		f.Close()

		_, err = LoadECDSA[*nist.PrivateKey](filename)
		switch {
		case err != nil && test.err == "":
			t.Fatalf("unexpected error for input %q:\n  %v", test.input, err)
		case err != nil && err.Error() != test.err:
			t.Fatalf("wrong error for input %q:\n  %v", test.input, err)
		case err == nil && test.err != "":
			t.Fatalf("LoadECDSA did not return error for input %q", test.input)
		}
	}
}

func TestSaveECDSA(t *testing.T) {
	f, err := ioutil.TempFile("", "saveecdsa_test.*.txt")
	if err != nil {
		t.Fatal(err)
	}
	file := f.Name()
	f.Close()
	defer os.Remove(file)

	key, _ := HexToECDSA[nist.PrivateKey](testPrivHex)
	if err := SaveECDSA(file, key); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadECDSA[nist.PrivateKey](file)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(key, loaded) {
		t.Fatal("loaded key not equal to saved key")
	}
}

func TestValidateSignatureValues(t *testing.T) {
	check := func(expected bool, v byte, r, s *big.Int) {
		if ValidateSignatureValues[nist.PublicKey](v, r, s, false) != expected {
			t.Errorf("mismatch for v: %d r: %d s: %d want: %v", v, r, s, expected)
		}
	}
	minusOne := big.NewInt(-1)
	one := common.Big1
	zero := common.Big0
	secp256k1nMinus1 := new(big.Int).Sub(secp256k1N, common.Big1)

	// correct v,r,s
	check(true, 0, one, one)
	check(true, 1, one, one)
	// incorrect v, correct r,s,
	check(false, 2, one, one)
	check(false, 3, one, one)

	// incorrect v, combinations of incorrect/correct r,s at lower limit
	check(false, 2, zero, zero)
	check(false, 2, zero, one)
	check(false, 2, one, zero)
	check(false, 2, one, one)

	// correct v for any combination of incorrect r,s
	check(false, 0, zero, zero)
	check(false, 0, zero, one)
	check(false, 0, one, zero)

	check(false, 1, zero, zero)
	check(false, 1, zero, one)
	check(false, 1, one, zero)

	// correct sig with max r,s
	check(true, 0, secp256k1nMinus1, secp256k1nMinus1)
	// correct v, combinations of incorrect r,s at upper limit
	check(false, 0, secp256k1N, secp256k1nMinus1)
	check(false, 0, secp256k1nMinus1, secp256k1N)
	check(false, 0, secp256k1N, secp256k1N)

	// current callers ensures r,s cannot be negative, but let's test for that too
	// as crypto package could be used stand-alone
	check(false, 0, minusOne, one)
	check(false, 0, one, minusOne)
}

func checkhash(t *testing.T, name string, f func([]byte) []byte, msg, exp []byte) {
	sum := f(msg)
	if !bytes.Equal(exp, sum) {
		t.Fatalf("hash %s mismatch: want: %x have: %x", name, exp, sum)
	}
}

func checkAddr(t *testing.T, addr0, addr1 common.Address) {
	if addr0 != addr1 {
		t.Fatalf("address mismatch: want: %x have: %x", addr0, addr1)
	}
}

// test to help Python team with integration of libsecp256k1
// skip but keep it after they are done
func TestPythonIntegration(t *testing.T) {
	kh := "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	k0, _ := HexToECDSA[nist.PrivateKey](kh)

	msg0 := Keccak256[nist.PublicKey]([]byte("foo"))
	sig0, _ := Sign(msg0, k0)

	msg1 := common.FromHex("00000000000000000000000000000000")
	sig1, _ := Sign(msg0, k0)

	t.Logf("msg: %x, privkey: %s sig: %x\n", msg0, kh, sig0)
	t.Logf("msg: %x, privkey: %s sig: %x\n", msg1, kh, sig1)
}
