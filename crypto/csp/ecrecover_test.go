package csp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcrecover(t *testing.T) {
	store, err := SystemStore("My")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	skid := "7e2747ab453e7e5c3c12feebb10d253cc772b7f5"
	// crt, err := store.GetByThumb(signCertThumb)
	crt, err := store.GetBySubjectId(skid)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Key : %x", crt.Info().PublicKeyBytes()[2:66])
	defer crt.Close()
	testData := "Test string"
	hash, err := NewHash(HashOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer func(hash *Hash) {
		if err := hash.Close(); err != nil {
			t.Fatal(err)
		}
	}(hash)
	_, err = hash.Write([]byte(testData))
	if err != nil {
		t.Fatal(err)
	}
	digest := hash.Sum(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("hash digest: %x", digest)
	sig, err := Sign(digest, crt)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signature: %x", sig)
	pub, err := Ecrecover(digest, sig)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, pub[:64], crt.Info().PublicKeyBytes())
	hash.Reset()
}