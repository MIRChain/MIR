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

package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/common/hexutil"
	"github.com/MIRChain/MIR/common/math"
	"github.com/MIRChain/MIR/crypto/csp"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/nist"
)

var (
	testmsg         = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	testsig         = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey      = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc     = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
	testpubkeycgost = hexutil.MustDecode("0x45a3d0a1d7e6cf63214c40ee903446179a2714dc11782f6e8a30f74e3fedc0bd06c1c816b1199f5e322929d3d1d534e4be7f44635c9e03d73fff58bb733112c8")
)

func TestEcrecover(t *testing.T) {
	pubkey, err := Ecrecover[nist.PublicKey](testmsg, testsig)
	if err != nil {
		t.Fatalf("recover error: %s", err)
	}
	if !bytes.Equal(pubkey, testpubkey) {
		t.Errorf("pubkey mismatch: want: %x have: %x", testpubkey, pubkey)
	}
}

func TestVerifySignature(t *testing.T) {
	sig := testsig[:len(testsig)-1] // remove recovery id
	if !VerifySignature[nist.PublicKey](testpubkey, testmsg, sig) {
		t.Errorf("can't verify signature with uncompressed key")
	}
	if !VerifySignature[nist.PublicKey](testpubkeyc, testmsg, sig) {
		t.Errorf("can't verify signature with compressed key")
	}

	if VerifySignature[nist.PublicKey](nil, testmsg, sig) {
		t.Errorf("signature valid with no key")
	}
	if VerifySignature[nist.PublicKey](testpubkey, nil, sig) {
		t.Errorf("signature valid with no message")
	}
	if VerifySignature[nist.PublicKey](testpubkey, testmsg, nil) {
		t.Errorf("nil signature valid")
	}
	if VerifySignature[nist.PublicKey](testpubkey, testmsg, append(common.CopyBytes(sig), 1, 2, 3)) {
		t.Errorf("signature valid with extra bytes at the end")
	}
	if VerifySignature[nist.PublicKey](testpubkey, testmsg, sig[:len(sig)-2]) {
		t.Errorf("signature valid even though it's incomplete")
	}
	wrongkey := common.CopyBytes(testpubkey)
	wrongkey[10]++
	if VerifySignature[nist.PublicKey](wrongkey, testmsg, sig) {
		t.Errorf("signature valid with with wrong public key")
	}
}

// This test checks that VerifySignature rejects malleable signatures with s > N/2.
func TestVerifySignatureMalleable(t *testing.T) {
	sig := hexutil.MustDecode("0x638a54215d80a6713c8d523a6adc4e6e73652d859103a36b700851cb0e61b66b8ebfc1a610c57d732ec6e0a8f06a9a7a28df5051ece514702ff9cdff0b11f454")
	key := hexutil.MustDecode("0x03ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd3138")
	msg := hexutil.MustDecode("0xd301ce462d3e639518f482c7f03821fec1e602018630ce621e1e7851c12343a6")
	if VerifySignature[nist.PublicKey](key, msg, sig) {
		t.Error("VerifySignature returned true for malleable signature")
	}
}

func TestDecompressPubkey(t *testing.T) {
	key, err := DecompressPubkey[nist.PublicKey](testpubkeyc)
	if err != nil {
		t.Fatal(err)
	}
	if uncompressed := FromECDSAPub(key); !bytes.Equal(uncompressed, testpubkey) {
		t.Errorf("wrong public key result: got %x, want %x", uncompressed, testpubkey)
	}
	if _, err := DecompressPubkey[nist.PublicKey](nil); err == nil {
		t.Errorf("no error for nil pubkey")
	}
	if _, err := DecompressPubkey[nist.PublicKey](testpubkeyc[:5]); err == nil {
		t.Errorf("no error for incomplete pubkey")
	}
	if _, err := DecompressPubkey[nist.PublicKey](append(common.CopyBytes(testpubkeyc), 1, 2, 3)); err == nil {
		t.Errorf("no error for pubkey with extra bytes at the end")
	}

	keyGost, err := DecompressPubkey[gost3410.PublicKey](testpubkeycgost)
	if err != nil {
		t.Fatal(err)
	}
	if uncompressed := FromECDSAPub(keyGost); !bytes.Equal(uncompressed, testpubkeycgost) {
		t.Errorf("wrong public key result: got %x, want %x", uncompressed, testpubkeycgost)
	}
	if _, err := DecompressPubkey[gost3410.PublicKey](nil); err == nil {
		t.Errorf("no error for nil pubkey")
	}
	if _, err := DecompressPubkey[gost3410.PublicKey](testpubkeycgost[:5]); err == nil {
		t.Errorf("no error for incomplete pubkey")
	}
	if _, err := DecompressPubkey[gost3410.PublicKey](append(common.CopyBytes(testpubkeycgost), 1, 2, 3)); err == nil {
		t.Errorf("no error for pubkey with extra bytes at the end")
	}

	store, err := csp.SystemStore("My")
	if err != nil {
		t.Errorf("Store error: %s", err)
	}
	defer store.Close()
	crt, err := store.GetBySubjectId("4ac93fc08bc0efd24180b0fa47f7309c257e8c85")
	if err != nil {
		t.Errorf("Get cert error: %s", err)
	}
	defer crt.Close()

	keyCsp, err := DecompressPubkey[csp.PublicKey](crt.Public().Raw())
	if err != nil {
		t.Fatal(err)
	}
	if uncompressed := FromECDSAPub(keyCsp); !bytes.Equal(uncompressed, crt.Public().Raw()) {
		t.Errorf("wrong public key result: got %x, want %x", uncompressed, crt.Public().Raw())
	}
	if _, err := DecompressPubkey[csp.PublicKey](nil); err == nil {
		t.Errorf("no error for nil pubkey")
	}
	if _, err := DecompressPubkey[csp.PublicKey](crt.Public().Raw()[:5]); err == nil {
		t.Errorf("no error for incomplete pubkey")
	}
	if _, err := DecompressPubkey[csp.PublicKey](append(common.CopyBytes(crt.Public().Raw()), 1, 2, 3)); err == nil {
		t.Errorf("no error for pubkey with extra bytes at the end")
	}
}

func TestCompressPubkey(t *testing.T) {
	key := &ecdsa.PublicKey{
		Curve: S256(),
		X:     math.MustParseBig256("0xe32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a"),
		Y:     math.MustParseBig256("0x0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652"),
	}
	compressed := CompressPubkey(nist.PublicKey{key})
	if !bytes.Equal(compressed, testpubkeyc) {
		t.Errorf("wrong public key result: got %x, want %x", compressed, testpubkeyc)
	}
}

func TestPubkeyRandom(t *testing.T) {
	const runs = 200

	for i := 0; i < runs; i++ {
		key, err := GenerateKey[nist.PrivateKey]()
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		pubkey2, err := DecompressPubkey[nist.PublicKey](CompressPubkey(nist.PublicKey{&key.PublicKey}))
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if !reflect.DeepEqual(key.PublicKey, *pubkey2.PublicKey) {
			t.Fatalf("iteration %d: keys not equal", i)
		}
	}
}

func BenchmarkEcrecoverSignature(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Ecrecover[nist.PublicKey](testmsg, testsig); err != nil {
			b.Fatal("ecrecover error", err)
		}
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	sig := testsig[:len(testsig)-1] // remove recovery id
	for i := 0; i < b.N; i++ {
		if !VerifySignature[nist.PublicKey](testpubkey, testmsg, sig) {
			b.Fatal("verify error")
		}
	}
}

func BenchmarkDecompressPubkey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := DecompressPubkey[nist.PublicKey](testpubkeyc); err != nil {
			b.Fatal(err)
		}
	}
}

func TestEcrecover_GOST(t *testing.T) {
	prvRaw := []byte{
		0x7A, 0x92, 0x9A, 0xDE, 0x78, 0x9B, 0xB9, 0xBE,
		0x10, 0xED, 0x35, 0x9D, 0xD3, 0x9A, 0x72, 0xC1,
		0x1B, 0x60, 0x96, 0x1F, 0x49, 0x39, 0x7E, 0xEE,
		0x1D, 0x19, 0xCE, 0x98, 0x91, 0xEC, 0x3B, 0x28,
	}
	dgst := []byte{
		0x2D, 0xFB, 0xC1, 0xB3, 0x72, 0xD8, 0x9A, 0x11,
		0x88, 0xC0, 0x9C, 0x52, 0xE0, 0xEE, 0xC6, 0x1F,
		0xCE, 0x52, 0x03, 0x2A, 0xB1, 0x02, 0x2E, 0x8E,
		0x67, 0xEC, 0xE6, 0x67, 0x2B, 0x04, 0x3E, 0xE5,
	}
	rnd := []byte{
		0x77, 0x10, 0x5C, 0x9B, 0x20, 0xBC, 0xD3, 0x12,
		0x28, 0x23, 0xC8, 0xCF, 0x6F, 0xCC, 0x7B, 0x95,
		0x6D, 0xE3, 0x38, 0x14, 0xE9, 0x5B, 0x7F, 0xE6,
		0x4F, 0xED, 0x92, 0x45, 0x94, 0xDC, 0xEA, 0xB3,
	}
	r := []byte{
		0x41, 0xAA, 0x28, 0xD2, 0xF1, 0xAB, 0x14, 0x82,
		0x80, 0xCD, 0x9E, 0xD5, 0x6F, 0xED, 0xA4, 0x19,
		0x74, 0x05, 0x35, 0x54, 0xA4, 0x27, 0x67, 0xB8,
		0x3A, 0xD0, 0x43, 0xFD, 0x39, 0xDC, 0x04, 0x93,
	}
	s := []byte{
		0x01, 0x45, 0x6C, 0x64, 0xBA, 0x46, 0x42, 0xA1,
		0x65, 0x3C, 0x23, 0x5A, 0x98, 0xA6, 0x02, 0x49,
		0xBC, 0xD6, 0xD3, 0xF7, 0x46, 0xB6, 0x31, 0xDF,
		0x92, 0x80, 0x14, 0xF6, 0xC5, 0xBF, 0x9C, 0x40,
	}
	c := gost3410.CurveIdGostR34102001TestParamSet()
	reverse(prvRaw)
	prv, err := gost3410.NewPrivateKey(c, prvRaw)
	if err != nil {
		t.Fatal(err)
	}
	sign, err := prv.SignDigest(dgst, bytes.NewBuffer(rnd))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(sign, append(r, s...)) {
		t.Fatalf("Signature mismatch: want: %x have: %x", sign, append(r, s...))
	}
	_r := new(big.Int).SetBytes(r)
	_s := new(big.Int).SetBytes(s)

	recovPubX, recovPubY, err := gost3410.RecoverCompact(*prv.C, dgst, _r, _s, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(prv.PublicKey.X.Bytes(), recovPubX.Bytes()) {
		t.Fatalf("Address mismatch: want: %x have: %x", prv.PublicKey.X.Bytes(), recovPubX.Bytes())
	}
	if !bytes.Equal(prv.PublicKey.Y.Bytes(), recovPubY.Bytes()) {
		t.Fatalf("Address mismatch: want: %x have: %x", prv.PublicKey.Y.Bytes(), recovPubY.Bytes())
	}
	testSig := make([]byte, 65)
	copy(testSig[:32], r)
	copy(testSig[32:64], s)
	testSig[64] = byte(0)
	gost3410.GostCurve = c
	revDgst := dgst
	reverse(revDgst)
	recoveredGostPub, err := Ecrecover[gost3410.PublicKey](revDgst, testSig[:])
	if err != nil {
		t.Fatalf("ECRecover error: %s", err)
	}
	if !bytes.Equal(prv.PublicKey.X.Bytes(), recoveredGostPub[1:33]) {
		t.Fatalf("Address mismatch: want: %x have: %x", prv.PublicKey.X.Bytes(), recoveredGostPub[1:33])
	}
	if !bytes.Equal(prv.PublicKey.Y.Bytes(), recoveredGostPub[33:65]) {
		t.Fatalf("Address mismatch: want: %x have: %x", prv.PublicKey.Y.Bytes(), recoveredGostPub[33:65])
	}
}
