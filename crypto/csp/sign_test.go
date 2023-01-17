package csp

import (
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
	skid := "7e2747ab453e7e5c3c12feebb10d253cc772b7f5"
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

func TestSignDigest(t *testing.T) {
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
	var res bool
	if res, err = VerifySignatureRaw(digest, sig, crt.Info().PublicKeyBytes()[2:66]); err != nil {
		t.Errorf("%+v", err)
	}
	assert.Equal(t, true, res)
	if res, err = VerifyDigestField(digest, sig, crt.Info().PublicKeyBytes()[2:66]); err != nil {
		t.Errorf("%+v", err)
	}
	assert.Equal(t, true, res)
	hash.Reset()
}