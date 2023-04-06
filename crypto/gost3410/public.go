// GoGOST -- Pure Go GOST cryptographic functions library
// Copyright (C) 2015-2021 Sergey Matveev <stargrave@stargrave.org>
// Copyright (C) 2022 Pavel Krolevets <pavelkrolevets@gmail.com>
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
	"crypto"
	"crypto/elliptic"
	"fmt"
	"math/big"
)

type PublicKey struct {
	C elliptic.Curve
	X *big.Int
	Y *big.Int
}
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
func NewPublicKey(c *Curve, raw []byte) (*PublicKey, error) {
	pointSize := c.PointSize()
	key := make([]byte, 2*pointSize)
	if len(raw) != len(key) {
		return &PublicKey{}, fmt.Errorf("gogost/gost3410: len(key) != %d", len(raw))
	}
	return &PublicKey{
		c,
		bytes2big(raw[:pointSize]),
		bytes2big(raw[pointSize : 2*pointSize]),
	}, nil
}

func (pub *PublicKey) Raw() []byte {
	pointSize := pub.C.Params().BitSize / 8
	raw := append(
		pad(pub.X.Bytes(), pointSize),
		pad(pub.Y.Bytes(), pointSize)...,
	)
	return raw
}

func (pub *PublicKey) VerifyDigest(digest, signature []byte) (bool, error) {
	pointSize := pub.C.Params().BitSize / 8
	if len(signature) != 2*pointSize {
		return false, fmt.Errorf("gogost/gost3410: len(signature) != %d", 2*pointSize)
	}
	r := bytes2big(signature[:pointSize])
	s := bytes2big(signature[pointSize:2*pointSize])
	if r.Cmp(zero) <= 0 ||
		r.Cmp(pub.C.Params().N) >= 0 ||
		s.Cmp(zero) <= 0 ||
		s.Cmp(pub.C.Params().N) >= 0 {
		return false, nil
	}
	e := bytes2big(digest)
	e.Mod(e, pub.C.Params().N)
	if e.Cmp(zero) == 0 {
		e = big.NewInt(1)
	}
	v := big.NewInt(0)
	v.ModInverse(e, pub.C.Params().N)
	z1 := big.NewInt(0)
	z2 := big.NewInt(0)
	z1.Mul(s, v)
	z1.Mod(z1, pub.C.Params().N)
	z2.Mul(r, v)
	z2.Mod(z2, pub.C.Params().N)
	z2.Sub(pub.C.Params().N, z2)
	p1x, p1y := pub.C.ScalarMult(pub.C.Params().Gx, pub.C.Params().Gy, z1.Bytes())
	q1x, q1y := pub.C.ScalarMult(pub.X, pub.Y, z2.Bytes())
	lm := big.NewInt(0)
	lm.Sub(q1x, p1x)
	if lm.Cmp(zero) < 0 {
		lm.Add(lm, pub.C.Params().P)
	}
	lm.ModInverse(lm, pub.C.Params().P)
	z1.Sub(q1y, p1y)
	lm.Mul(lm, z1)
	lm.Mod(lm, pub.C.Params().P)
	lm.Mul(lm, lm)
	lm.Mod(lm, pub.C.Params().P)
	lm.Sub(lm, p1x)
	lm.Sub(lm, q1x)
	lm.Mod(lm, pub.C.Params().P)
	if lm.Cmp(zero) < 0 {
		lm.Add(lm, pub.C.Params().P)
	}
	lm.Mod(lm, pub.C.Params().N)
	return lm.Cmp(r) == 0, nil
}

func (our *PublicKey) Equal(theirKey crypto.PublicKey) bool {
	their, ok := theirKey.(*PublicKey)
	if !ok {
		return false
	}
	return our.X.Cmp(their.X) == 0 && our.Y.Cmp(their.Y) == 0 && our.C.Params().Name == their.C.Params().Name
}

func RecoverCompact(curve Curve, digest []byte, r *big.Int, s *big.Int, i int) (*big.Int, *big.Int, error) {

	if r.Cmp(zero) <= 0 ||
		r.Cmp(curve.Q) >= 0 ||
		s.Cmp(zero) <= 0 ||
		s.Cmp(curve.Q) >= 0 {
		return nil, nil, fmt.Errorf("error at r")
	}
	z := bytes2big(digest)
	z.Mod(z, curve.Q)
	if z.Cmp(zero) == 0 {
		z = big.NewInt(1)
	}
	var w, u1, u2= new(big.Int), new(big.Int), new(big.Int)

	// Calculate a curve point R = ( x 1 , y 1 ) 
	//  R=(x 1,y 1) where x 1  is one of r r, r + n, r + 2n, etc. (provided x 1 is not too large for a field element) 
	// and y 1  is a value such that the curve equation is satisfied. Note that there may be several curve points satisfying these conditions, and each different R value results in a distinct recovered key.
	// Basically x = (n * i) + r

	Rx := new(big.Int).Mul(curve.Q, new(big.Int).SetInt64(int64(i/2)))
	Rx.Add(Rx, r)
	if Rx.Cmp(curve.P) != -1 {
		return nil, nil, fmt.Errorf("calculated Rx is larger than curve P")
	}
	Ry, err := curve.DecompressPoint(Rx, i%2 == 1)
	if err != nil {
		return nil, nil, fmt.Errorf("error at DecompressPoint")
	}
	w.ModInverse(r, curve.Q)
	
	u1.Mul(s, w)
	u1.Mod(u1, curve.Q)

	u2.Mul(z, w)
	u2.Neg(u2)
	u2.Mod(u2, curve.Q)

	u1Gx, u1Gy := curve.ScalarMult(curve.X, curve.Y, u1.Bytes())
	u2Rx, u2Ry := curve.ScalarMult(Rx, Ry, u2.Bytes())
	Qx, Qy := curve.Add(u1Gx, u1Gy, u2Rx, u2Ry)
	return Qx, Qy, nil

}

func (prv PublicKey) GetX() *big.Int {
	return prv.X
}

func (prv PublicKey) GetY() *big.Int {
	return prv.Y
}