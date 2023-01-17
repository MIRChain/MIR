package csp

/*
#include "common.h"

*/
import "C"

import (
	"fmt"
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/crypto/csp/gost3410"
)


func Ecrecover(m, signature []byte) ([]byte, error) {
	curve := gost3410.CurveIdGostR34102001CryptoProAParamSet()
	pointSize := curve.PointSize()
	if len(signature) != 2*pointSize {
		return nil, fmt.Errorf("gogost/gost3410: len(signature) != %d", 2*pointSize)
	}
	r := bytes2big(signature[pointSize:])
	s := bytes2big(signature[:pointSize])

	if r.Cmp(big.NewInt(0)) <= 0 ||
		r.Cmp(curve.Q) >= 0 ||
		s.Cmp(big.NewInt(0)) <= 0 ||
		s.Cmp(curve.Q) >= 0 {
		return nil, fmt.Errorf("error at r")
	}
	z := bytes2big(m)
	z.Mod(z, curve.Q)
	if z.Cmp(big.NewInt(0)) == 0 {
		z = big.NewInt(1)
	}
	var Rx,Ry, w, u1, u2 = new(big.Int), new(big.Int), new(big.Int), new(big.Int), new(big.Int)
	
	x3 := new(big.Int).Mul(r, r)
	x3.Mul(x3, r)
	aX := new(big.Int).Mul(curve.A, r)
	x3.Add(x3, aX)
	x3.Add(x3, curve.B)
	x3.Mod(x3, curve.P)
	
	y0 := new(big.Int).ModSqrt(x3, curve.P)

	if y0.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf("error at computing y0")
	}

	y1 := new(big.Int).Sub(curve.P, y0)
	Rx.Set(r)
	Ry.Set(y0)

	w.ModInverse(r, curve.Q)
	
	u1.Mul(s, w)
	u1.Mod(u1, curve.Q)

	u2.Mul(z, w)
	u2.Neg(u2)
	u2.Mod(u2, curve.Q)

	u1Gx, u1Gy, err := curve.Exp(u1, curve.X, curve.Y)
	if err != nil {
		return nil, fmt.Errorf("error at computing u1Gx, u1Gy")
	}
	u2Rx, u2Ry, err := curve.Exp(u2, Rx, Ry)
	if err != nil {
		return nil, fmt.Errorf("error at computing u2Rx, u2Ry")
	}
	Qx, Qy := curve.Add(u1Gx, u1Gy, u2Rx, u2Ry)
	raw := append(
		pad(Qy.Bytes(), pointSize),
		pad(Qx.Bytes(), pointSize)...,
	)
	reverse(raw)
	recPub := append(
		[]byte{0x04, 0x40},
		raw...
	)
	if len(recPub) != 66 {
		return nil, fmt.Errorf("key len isnt 66")
	}
	ver, err := VerifySignatureRaw(m, signature, recPub)
	if ver {
		return append(raw, byte(0)), nil
	}	
	u2Rx_, u2Ry_, err :=  curve.Exp(Rx, y1, u2)
	if err != nil {
		return nil, fmt.Errorf("error at computing u2Rx_, u2Ry_")
	}
	Qx, Qy = curve.Add(u1Gx, u1Gy, u2Rx_, u2Ry_)
	_raw := append(
		pad(Qy.Bytes(), pointSize),
		pad(Qx.Bytes(), pointSize)...,
	)
	reverse(_raw)
	_recPub := append(
		[]byte{0x04, 0x40},
		_raw...
	)
	ver, err = VerifySignatureRaw(m, signature, _recPub)
	if ver {
		return append(raw, byte(1)), nil
	}	
	return nil, err
}

func pad(d []byte, size int) []byte {
	return append(make([]byte, size-len(d)), d...)
}

func reverse(d []byte) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}

func bytes2big(d []byte) *big.Int {
	return big.NewInt(0).SetBytes(d)
}
