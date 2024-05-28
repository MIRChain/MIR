package ecies

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/MIRChain/MIR/crypto/gost3410"
)

// Validate the ECDH component.
func TestSharedKeyGost(t *testing.T) {
	prv1, err := GenerateKey[gost3410.PrivateKey, gost3410.PublicKey](rand.Reader, DefaultCurve, nil)
	if err != nil {
		t.Fatal(err)
	}
	skLen := MaxSharedKeyLength(&prv1.PublicKey) / 2

	prv2, err := GenerateKey[gost3410.PrivateKey, gost3410.PublicKey](rand.Reader, DefaultCurve, nil)
	if err != nil {
		t.Fatal(err)
	}

	sk1, err := prv1.GenerateShared(&prv2.PublicKey, skLen, skLen)
	if err != nil {
		t.Fatal(err)
	}

	sk2, err := prv2.GenerateShared(&prv1.PublicKey, skLen, skLen)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(sk1, sk2) {
		t.Fatal(ErrBadSharedKeys)
	}
}

// Verify that an encrypted message can be successfully decrypted.
func TestEncryptDecryptGost(t *testing.T) {
	prv1, err := GenerateKey[gost3410.PrivateKey, gost3410.PublicKey](rand.Reader, DefaultCurve, nil)
	if err != nil {
		t.Fatal(err)
	}

	prv2, err := GenerateKey[gost3410.PrivateKey, gost3410.PublicKey](rand.Reader, DefaultCurve, nil)
	if err != nil {
		t.Fatal(err)
	}

	message := []byte("Hello, world.")
	ct, err := Encrypt[gost3410.PrivateKey](rand.Reader, &prv2.PublicKey, message, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	pt, err := prv2.Decrypt(ct, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pt, message) {
		t.Fatal("ecies: plaintext doesn't match message")
	}

	_, err = prv1.Decrypt(ct, nil, nil)
	if err == nil {
		t.Fatal("ecies: encryption should not have succeeded")
	}
}

func TestDecryptShared2Gost(t *testing.T) {
	prv, err := GenerateKey[gost3410.PrivateKey, gost3410.PublicKey](rand.Reader, gost3410.GostCurve, nil)
	if err != nil {
		t.Fatal(err)
	}
	message := []byte("Hello, world.")
	shared2 := []byte("shared data 2")
	ct, err := Encrypt[gost3410.PrivateKey](rand.Reader, &prv.PublicKey, message, nil, shared2)
	if err != nil {
		t.Fatal(err)
	}

	// Check that decrypting with correct shared data works.
	pt, err := prv.Decrypt(ct, nil, shared2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt, message) {
		t.Fatal("ecies: plaintext doesn't match message")
	}

	// Decrypting without shared data or incorrect shared data fails.
	if _, err = prv.Decrypt(ct, nil, nil); err == nil {
		t.Fatal("ecies: decrypting without shared data didn't fail")
	}
	if _, err = prv.Decrypt(ct, nil, []byte("garbage")); err == nil {
		t.Fatal("ecies: decrypting with incorrect shared data didn't fail")
	}
}

func TestBasicKeyValidationGost(t *testing.T) {
	badBytes := []byte{0, 1, 5, 6, 7, 8, 9}

	prv, err := GenerateKey[gost3410.PrivateKey, gost3410.PublicKey](rand.Reader, gost3410.GostCurve, nil)
	if err != nil {
		t.Fatal(err)
	}

	message := []byte("Hello, world.")
	ct, err := Encrypt[gost3410.PrivateKey](rand.Reader, &prv.PublicKey, message, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, b := range badBytes {
		ct[0] = b
		_, err := prv.Decrypt(ct, nil, nil)
		if err != ErrInvalidPublicKey {
			t.Fatal("ecies: validated an invalid key")
		}
	}
}
