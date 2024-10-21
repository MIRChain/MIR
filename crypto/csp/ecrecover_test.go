package csp

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/stretchr/testify/assert"
)

func TestEcrecover(t *testing.T) {
	store, err := SystemStore("My")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	skid := "71732462bbc029d911e6d16a3ed00d9d1d772620"
	// crt, err := store.GetByThumb(signCertThumb)
	crt, err := store.GetBySubjectId(skid)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Pub Key : %x", crt.Info().PublicKeyBytes()[2:66])
	defer crt.Close()
	hash, err := NewHash(HashOptions{SignCert: crt, HashAlg: GOST_R3411_12_256})
	if err != nil {
		t.Fatal(err)
	}
	defer func(hash *Hash) {
		if err := hash.Close(); err != nil {
			t.Fatal(err)
		}
	}(hash)
	_, err = hash.Write([]byte(skid))
	if err != nil {
		t.Fatal(err)
	}
	digest := hash.Sum(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("hash digest: %x", digest)
	sig, err := hash.Sign()
	if err != nil {
		t.Fatal(err)
	}
	// sig, _ := hex.DecodeString("25e53018b290d115ae96ac6abb49c8e9fd88297499b7a68e600c78d1c20083285ee04229f089bea48703d90c6dd33ed1b4b9051cd5611e7ec9994c968bcb8412")
	t.Logf("signature: %s", hex.EncodeToString(sig))
	err = hash.VerifyWithPub(sig, crt.Info().PublicKeyBytes())
	assert.NoError(t, err)
	pubKey := make([]byte, 64)
	copy(pubKey[:], crt.Info().PublicKeyBytes()[2:66])
	reverse(pubKey)
	pubKeyX := new(big.Int).SetBytes(pubKey[32:64])
	pubKeyY := new(big.Int).SetBytes(pubKey[:32])
	t.Logf("pubKeyX : %s", pubKeyX.String())
	t.Logf("pubKeyY : %s", pubKeyY.String())
	gostKey := gost3410.PublicKey{
		C: gost3410.CurveIdGostR34102001CryptoProAParamSet(),
		X: pubKeyX,
		Y: pubKeyY,
	}
	r, s, revHash := revertCSP(digest, sig)
	recoveredPub := make([]byte, 0, 64)
	for i := 0; i < 4; i++ {
		X, Y, err := gost3410.RecoverCompact(*gost3410.CurveIdGostR34102001CryptoProAParamSet(), revHash, r, s, i)
		if err == nil && X.Cmp(gostKey.X) == 0 && Y.Cmp(gostKey.Y) == 0 {
			recoveredPub = append(recoveredPub, Y.Bytes()...)
			recoveredPub = append(recoveredPub, X.Bytes()...)
			reverse(recoveredPub)
		}
	}
	assert.Equal(t, recoveredPub, crt.Info().PublicKeyBytes()[2:66])
	hash.Reset()
}

func revertCSP(hash, signature []byte) (r, s *big.Int, revertHash []byte) {
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
