package csp

/*
#include "common.h"

*/
import "C"

import (
	"fmt"
	"math/big"
	"unsafe"

	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
)

func (h *Hash) Sign() ([]byte, error) {
	var slen C.DWORD
	if C.CryptSignHash(h.hHash, h.dwKeySpec, nil, 0, nil, &slen) == 0 {
		return nil, getErr("Error calculating signature size")
	}
	if slen == 0 {
		return nil, nil
	}
	res := make([]byte, int(slen))
	if C.CryptSignHash(h.hHash, h.dwKeySpec, nil, 0, (*C.BYTE)(&res[0]), &slen) == 0 {
		return nil, getErr("Error calculating signature value")
	}
	return res, nil
}

func (h *Hash) Verify(signer Cert, sig []byte) error {
	var hPubKey C.HCRYPTKEY
	// Get the public key from the certificate
	if C.CryptImportPublicKeyInfo(h.hProv, C.MY_ENC_TYPE, &signer.pCert.pCertInfo.SubjectPublicKeyInfo, &hPubKey) == 0 {
		return getErr("Error getting certificate public key handle")
	}
	var ptr unsafe.Pointer
	if len(sig) > 0 {
		ptr = unsafe.Pointer(&sig[0])
	}
	if C.CryptVerifySignature(h.hHash, (*C.BYTE)(ptr), C.DWORD(len(sig)), hPubKey, nil, 0) == 0 {
		return getErr("Error verifying hash signature")
	}
	return nil
}

func Sign(digest []byte, crt Cert) ([]byte, error) {
	hash, err := NewHash(HashOptions{SignCert: crt})
	if err != nil {
		return nil, err
	}
	_, err = hash.Write(digest)
	if err != nil {
		return nil, err
	}
	sig, err := hash.Sign()
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func VerifySignature(crt Cert, digestHash, signature []byte) (bool, error) {
	hash, err := NewHash(HashOptions{SignCert: crt})
	if err != nil {
		return false, err
	}
	_, err = hash.Write(digestHash)
	if err != nil {
		return false, err
	}
	err = hash.Verify(crt, signature)
	if err != nil {
		return false, err
	}
	return true, nil
}

func VerifySignatureRaw(digestHash, signature, pub []byte) (bool, error) {
	hash, err := NewHash(HashOptions{})
	if err != nil {
		return false, err
	}
	_, err = hash.Write(digestHash)
	if err != nil {
		return false, err
	}
	err = hash.VerifyWithPub(signature, pub)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (h *Hash) VerifyWithPub(sig, pub []byte) error {
	blob := []byte{
        0x06, 0x20, 0x00, 0x00, 0x49, 0x2e, 0x00, 0x00, 0x4d, 0x41,
        0x47, 0x31, 0x00, 0x02, 0x00, 0x00, 0x30, 0x13, 0x06, 0x07,
        0x2a, 0x85, 0x03, 0x02, 0x02, 0x23, 0x01, 0x06, 0x08, 0x2a,
        0x85, 0x03, 0x07, 0x01, 0x01, 0x02, 0x02}
	blob = append(blob[:], pub[2:66]...)
	if len(blob) != 101 {
		return getErr("Wrogn pub key lengh")
	}
	var (
		res     Key
		decrKey C.HCRYPTKEY
		errMsg  = "Error importing key blob"
	)
	bufBytes := C.CBytes(blob)
	defer C.free(bufBytes)
	if C.CryptImportKey(h.hProv, (*C.BYTE)(bufBytes), C.DWORD(len(blob)), decrKey, 0, &res.hKey) == 0 {
		return getErr(errMsg)
	}
	var sigPtr unsafe.Pointer
	if len(sig) > 0 {
		sigPtr = unsafe.Pointer(&sig[0])
	}
	if C.CryptVerifySignature(h.hHash, (*C.BYTE)(sigPtr), C.DWORD(len(sig)), res.hKey, nil, 0) == 0 {
		return getErr("Error verifying hash signature")
	}
	return nil
}

func VerifyDigestField(digest, signature, pub []byte) (bool, error) {
	curve := gost3410.CurveIdGostR34102001CryptoProAParamSet()	
	pointSize := curve.PointSize()
	if len(signature) != 2*pointSize {
		return false, fmt.Errorf("gogost/gost3410: len(signature) != %d", 2*pointSize)
	}
	key := make([]byte, 2*pointSize)
	copy(key, pub)
	reverse(key)
	X := new(big.Int).SetBytes(key[pointSize : 2*pointSize])
	Y := new(big.Int).SetBytes(key[:pointSize])

	s := bytes2big(signature[:pointSize])
	r := bytes2big(signature[pointSize:])

	if r.Cmp(big.NewInt(0)) <= 0 ||
		r.Cmp(curve.Q) >= 0 ||
		s.Cmp(big.NewInt(0)) <= 0 ||
		s.Cmp(curve.Q) >= 0 {
		return false, nil
	}
	e := bytes2big(digest)
	e.Mod(e, curve.Q)
	if e.Cmp(big.NewInt(0)) == 0 {
		e = big.NewInt(1)
	}
	v := big.NewInt(0)
	v.ModInverse(e, curve.Q)
	z1 := big.NewInt(0)
	z2 := big.NewInt(0)
	z1.Mul(s, v)
	z1.Mod(z1, curve.Q)
	z2.Mul(r, v)
	z2.Mod(z2, curve.Q)
	z2.Sub(curve.Q, z2)
	p1x, p1y, err := curve.Exp(z1, curve.X, curve.Y)
	if err != nil {
		return false, err
	}
	q1x, q1y, err := curve.Exp(z2, X, Y)
	if err != nil {
		return false, err
	}
	lm := big.NewInt(0)
	lm.Sub(q1x, p1x)
	if lm.Cmp(big.NewInt(0)) < 0 {
		lm.Add(lm, curve.P)
	}
	lm.ModInverse(lm, curve.P)
	z1.Sub(q1y, p1y)
	lm.Mul(lm, z1)
	lm.Mod(lm, curve.P)
	lm.Mul(lm, lm)
	lm.Mod(lm, curve.P)
	lm.Sub(lm, p1x)
	lm.Sub(lm, q1x)
	lm.Mod(lm, curve.P)
	if lm.Cmp(big.NewInt(0)) < 0 {
		lm.Add(lm, curve.P)
	}
	lm.Mod(lm, curve.Q)
	return lm.Cmp(r) == 0, nil
}
