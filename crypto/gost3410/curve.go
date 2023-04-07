// GoGOST -- Pure Go GOST cryptographic functions library
// Copyright (C) 2015-2021 Sergey Matveev <stargrave@stargrave.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gost3410

import (
	"crypto/elliptic"
	"errors"
	"math/big"
)

var (
	zero    *big.Int = big.NewInt(0)
	bigInt1 *big.Int = big.NewInt(1)
	bigInt2 *big.Int = big.NewInt(2)
	bigInt3 *big.Int = big.NewInt(3)
	bigInt4 *big.Int = big.NewInt(4)
)

type Curve struct {
	Name string // Just simple identifier

	P *big.Int // Characteristic of the underlying prime field
	Q *big.Int // Elliptic curve subgroup order

	Co *big.Int // Cofactor

	// Equation coefficients of the elliptic curve in canonical form
	A *big.Int
	B *big.Int

	// Equation coefficients of the elliptic curve in twisted Edwards form
	E *big.Int
	D *big.Int

	// Basic point X and Y coordinates
	X *big.Int
	Y *big.Int

	// Cached s/t parameters for Edwards curve points conversion
	edS *big.Int
	edT *big.Int
}

func NewCurve(p, q, a, b, x, y, e, d, co *big.Int) (*Curve, error) {
	c := Curve{
		Name: "unknown",
		P:    p,
		Q:    q,
		A:    a,
		B:    b,
		X:    x,
		Y:    y,
	}
	r1 := big.NewInt(0)
	r2 := big.NewInt(0)
	r1.Mul(c.Y, c.Y)
	r1.Mod(r1, c.P)
	r2.Mul(c.X, c.X)
	r2.Add(r2, c.A)
	r2.Mul(r2, c.X)
	r2.Add(r2, c.B)
	r2.Mod(r2, c.P)
	c.pos(r2)
	if r1.Cmp(r2) != 0 {
		return nil, errors.New("gogost/gost3410: invalid curve parameters")
	}
	if e != nil && d != nil {
		c.E = e
		c.D = d
	}
	if co == nil {
		c.Co = bigInt1
	} else {
		c.Co = co
	}
	return &c, nil
}

func (c *Curve) PointSize() int {
	return PointSize(c.P)
}

func (c *Curve) pos(v *big.Int) {
	if v.Cmp(zero) < 0 {
		v.Add(v, c.P)
	}
}

func (c *Curve) Params() *elliptic.CurveParams {
	return &elliptic.CurveParams{
		P: c.P,
		N: c.Q,
		B: c.B,
		Gx: c.X, 
		Gy: c.Y,
		BitSize: c.P.BitLen(),
		Name: c.Name,
	}
}

func (c *Curve) AsElliptic() elliptic.Curve {
	return c
}

func (c *Curve) Add(x1, y1, x2, y2 *big.Int) (x, y *big.Int) {
	var t, tx, ty, X, Y big.Int
	if x1.Cmp(x2) == 0 && y1.Cmp(y2) == 0 {
		// double
		t.Mul(x1, x1)
		t.Mul(&t, bigInt3)
		t.Add(&t, c.A)
		tx.Mul(bigInt2, y2)
		tx.ModInverse(&tx, c.P)
		t.Mul(&t, &tx)
		t.Mod(&t, c.P)
	} else {
		tx.Sub(x2, x1)
		tx.Mod(&tx, c.P)
		c.pos(&tx)
		ty.Sub(y2, y1)
		ty.Mod(&ty, c.P)
		c.pos(&ty)
		t.ModInverse(&tx, c.P)
		t.Mul(&t, &ty)
		t.Mod(&t, c.P)
	}
	tx.Mul(&t, &t)
	tx.Sub(&tx, x1)
	tx.Sub(&tx, x2)
	tx.Mod(&tx, c.P)
	c.pos(&tx)
	ty.Sub(x1, &tx)
	ty.Mul(&ty, &t)
	ty.Sub(&ty, y1)
	ty.Mod(&ty, c.P)
	c.pos(&ty)
	X.Set(&tx)
	Y.Set(&ty)
	return &X, &Y
}

func(c *Curve) Double(x1, y1 *big.Int) (x, y *big.Int) {
	return c.Add(x1, y1, x1, y1)
}

func (c *Curve) ScalarMult(xS, yS *big.Int, k []byte) (x, y *big.Int) {
	degree := new(big.Int).SetBytes(k)
	if degree.Cmp(zero) == 0 {
		return nil, nil
	}
	dg := big.NewInt(0).Sub(degree, bigInt1)
	tx := big.NewInt(0).Set(xS)
	ty := big.NewInt(0).Set(yS)
	cx := big.NewInt(0).Set(xS)
	cy := big.NewInt(0).Set(yS)
	for dg.Cmp(zero) != 0 {
		if dg.Bit(0) == 1 {
			tx, ty = c.Add(tx, ty, cx, cy)
		}
		dg.Rsh(dg, 1)
		cx, cy = c.Add(cx, cy, cx, cy)
	}
	return tx, ty
}

func (c *Curve) ScalarBaseMult(k []byte)(x, y *big.Int){
	return c.ScalarMult(c.X, c.Y, k)
}

func (our *Curve) Equal(their *Curve) bool {
	return our.P.Cmp(their.P) == 0 &&
		our.Q.Cmp(their.Q) == 0 &&
		our.A.Cmp(their.A) == 0 &&
		our.B.Cmp(their.B) == 0 &&
		our.X.Cmp(their.X) == 0 &&
		our.Y.Cmp(their.Y) == 0 &&
		((our.E == nil && their.E == nil) || our.E.Cmp(their.E) == 0) &&
		((our.D == nil && their.D == nil) || our.D.Cmp(their.D) == 0) &&
		our.Co.Cmp(their.Co) == 0
}

// DecompressPoint decompresses a point on the given curve given the X point and
// the solution to use.
func (curve *Curve) DecompressPoint(x *big.Int, ybit bool) (*big.Int, error) {
	// y = +-sqrt(x^3 + ax + b)
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)
	aX := new(big.Int).Mul(x, curve.A)
	x3.Add(x3, aX)
	x3.Add(x3, curve.B)
	y := x3.ModSqrt(x3, curve.P)
	if y == nil {
		return nil, errors.New("failed to decompress elliptic curve point from given X coordinate")
	}

	if ybit != (y.Bit(0) == 1) {
		y.Sub(curve.P, y)
	}
	if ybit != (y.Bit(0) == 1) {
		return nil, errors.New("ybit doesn't match oddness")
	}
	return y, nil
}

// Marshal converts a point on the curve into the uncompressed form specified in
// SEC 1, Version 2.0, Section 2.3.3. If the point is not on the curve (or is
// the conventional point at infinity), the behavior is undefined.
func Marshal(curve elliptic.Curve, x, y *big.Int) []byte {
	byteLen := (curve.Params().BitSize + 7) / 8

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point

	x.FillBytes(ret[1 : 1+byteLen])
	y.FillBytes(ret[1+byteLen : 1+2*byteLen])

	return ret
}

// Unmarshal converts a point, serialized by Marshal, into an x, y pair. It is
// an error if the point is not in uncompressed form, is not on the curve, or is
// the point at infinity. On error, x = nil.
func Unmarshal(curve elliptic.Curve, data []byte) (x, y *big.Int) {
	byteLen := (curve.Params().BitSize + 7) / 8
	if len(data) != 1+2*byteLen {
		return nil, nil
	}
	if data[0] != 4 { // uncompressed form
		return nil, nil
	}
	p := curve.Params().P
	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
	if x.Cmp(p) >= 0 || y.Cmp(p) >= 0 {
		return nil, nil
	}
	if !curve.IsOnCurve(x, y) {
		return nil, nil
	}
	return
}


func (curve *Curve) IsOnCurve(x, y *big.Int) bool {
	if x.Sign() < 0 || x.Cmp(curve.P) >= 0 ||
		y.Sign() < 0 || y.Cmp(curve.P) >= 0 {
		return false
	}
	// y = +-sqrt(x^3 + ax + b)
	x3 := new(big.Int).Mul(x, x)
	x3.Mul(x3, x)
	aX := new(big.Int).Mul(x, curve.A)
	x3.Add(x3, aX)
	x3.Add(x3, curve.B)
	res_y := x3.ModSqrt(x3, curve.P)
	if res_y.Cmp(y) != 0 {
		res_y.Sub(curve.P, res_y)
		if res_y.Cmp(y) != 0 {
			return false
		}
	}
	return true
}
