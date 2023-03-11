package csp

//#include "common.h"
import "C"

import (
	"fmt"
	"math/big"
	"unsafe"

	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
)
type PublicKey struct {
	Curve *gost3410.Curve
	Ds    int
	X     *big.Int
	Y     *big.Int
}

// Get pub key values from cert
func (c Cert) Public() *PublicKey {
	ci := CertInfo{c.pCert.pCertInfo}
	pb := ci.pCertInfo.SubjectPublicKeyInfo.PublicKey
	pubKeyBytes := C.GoBytes(unsafe.Pointer(pb.pbData), C.int(pb.cbData))[2:66]
	curveBitLen := gost3410.CurveIdGostR34102001CryptoProAParamSet().P.BitLen()
	curveByteLen := curveBitLen/8
	reverse(pubKeyBytes)
	return &PublicKey{
		gost3410.CurveIdGostR34102001CryptoProAParamSet(),
		curveByteLen,
		new(big.Int).SetBytes(pubKeyBytes[curveByteLen : 2*curveByteLen]),
		new(big.Int).SetBytes(pubKeyBytes[:curveByteLen]),
	}
}

func NewPublicKey(raw []byte) (*PublicKey, error) {
	curve := gost3410.CurveIdGostR34102001CryptoProAParamSet()
	curveBitLen := curve.P.BitLen()
	curveByteLen := curveBitLen/8
	if len(raw) != 2*curveByteLen {
		return nil, fmt.Errorf("raw pub key bytes not equal curve params %d", curveByteLen)
	}
	key := make([]byte, 2*curveByteLen)
	copy(key, raw)
	reverse(key)
	return &PublicKey{
		curve,
		curveByteLen,
		new(big.Int).SetBytes(key[curveByteLen : 2*curveByteLen]),
		new(big.Int).SetBytes(key[:curveByteLen]),
	}, nil
}

func (p *PublicKey) Raw() []byte {
	curve := gost3410.CurveIdGostR34102001CryptoProAParamSet()
	curveBitLen := curve.P.BitLen()
	curveByteLen := curveBitLen/8
	key := make([]byte, 2*curveByteLen)
	copy(key[:curveByteLen], p.Y.Bytes())
	copy(key[curveByteLen:2*curveByteLen], p.X.Bytes())
	reverse(key)
	return key
}

func (prv PublicKey) GetX() *big.Int {
	return prv.X
}
func (prv PublicKey) GetY() *big.Int {
	return prv.Y
}