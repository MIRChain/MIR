package csp

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSignHash(t *testing.T) {
	// if signCertThumb == "2f2bf249a3c6856e86dc8f968bc5b983a76af7cc" {
	// 	t.Skip("certificate for sign test not provided")
	// }
	store, err := SystemStore("My")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	skid := "4ac93fc08bc0efd24180b0fa47f7309c257e8c85"
	// crt, err := store.GetByThumb(signCertThumb)
	crt, err := store.GetBySubjectId(skid)
	if err != nil {
		t.Fatal(err)
	}
	defer crt.Close()
	testData := "Test string"
	hash, err := NewHash(HashOptions{SignCert: crt})
	if err != nil {
		t.Fatal(err)
	}
	defer func(hash *Hash) {
		if err := hash.Close(); err != nil {
			t.Fatal(err)
		}
	}(hash)
	// fmt.Fprintf(hash, "%s", testData)
	_, err = hash.Write([]byte(testData))
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
	if len(sig) == 0 {
		t.Fatal("empty signature")
	}
	t.Logf("signature: %x", sig)
	// fmt.Fprintf(hash, "%s", testData)
	if err := hash.Verify(crt, sig); err != nil {
		t.Errorf("%+v", err)
	}
	hash.Reset()
	// fmt.Fprintf(hash, "%s", "wrong data")
	// var cryptErr Error
	// if err := hash.Verify(crt, sig); !errors.As(err, &cryptErr) {
	// 	t.Errorf("expected crypto Error, got %+v", err)
	// } else if cryptErr.Code != 0x80090006 {
	// 	t.Errorf("expected error 0x80090006, got %x", cryptErr.Code)
	// }
}

func TestHashSignVeritfy(t *testing.T) {
	store, err := SystemStore("My")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	skid := "4ac93fc08bc0efd24180b0fa47f7309c257e8c85"
	// crt, err := store.GetByThumb(signCertThumb)
	crt, err := store.GetBySubjectId(skid)
	if err != nil {
		t.Fatal(err)
	}
	defer crt.Close()
	testData := []byte([]byte("Foo"))
	hash, err := NewHash(HashOptions{SignCert: crt, HashAlg: GOST_R3411_12_256})
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
	sig, err := hash.Sign()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signature: %x", sig)
	err = hash.VerifyWithPub(sig, crt.Info().PublicKeyBytes())
	assert.NoError(t, err)
	hash.Reset()
}
func TestSignDigest(t *testing.T) {
	store, err := SystemStore("My")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	skid := "4ac93fc08bc0efd24180b0fa47f7309c257e8c85"
	// crt, err := store.GetByThumb(signCertThumb)
	crt, err := store.GetBySubjectId(skid)
	if err != nil {
		t.Fatal(err)
	}
	defer crt.Close()
	digest, err := hex.DecodeString("8ed192858e269a2fdb6f663fa70140930134c8f7f4b5590df7dbba8645fde84d") // hash "Foo"
	if err != nil {
		t.Fatal(err)
	}
	sig, err := Sign(crt, digest)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signature: %x", sig)
	res, err := VerifySignature(crt, digest, sig)
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("Sig isnt valid")
	}
	sig, err = hex.DecodeString("debcf9aeb9235ad8a50e09db99a06bb5a885f3a6b2e529b550843ba0978782ed4077cf47f21cab20238ae7d962178d4a3674e65b4d2360cfe0811755a1c52a7c")
	if err != nil {
		t.Fatal(err)
	}
	res, err = VerifySignatureRaw(digest, sig, crt.Info().PublicKeyBytes())
	if err != nil {
		t.Fatal(err)
	}
	if !res {
		t.Fatal("Sig isnt valid")
	}
	// Check wrong sig
	wrongSig, err := hex.DecodeString("1bf1b9c3bd6ca1871d0dff0a394c3e8331a6c91e83c2a4b86e80f52fb2fbcb7b56ad2b42b00844a0cdb78a6283c39c1e125c75924b7dc2893358b48a892288d6")
	if err != nil {
		t.Fatal(err)
	}
	res, err = VerifySignatureRaw(digest, wrongSig, crt.Info().PublicKeyBytes())
	if err == nil {
		t.Fatal(err)
	}
	if res {
		t.Fatal("Sig should be invalid")
	}
}