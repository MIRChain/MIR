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
	"crypto"
	"errors"
	"fmt"
	"io"
	"math/big"
)

type PrivateKey struct {
	PublicKey
	C   *Curve
	Key *big.Int
}

func NewPrivateKey(c *Curve, raw []byte) (*PrivateKey, error) {
	pointSize := c.PointSize()
	if len(raw) != pointSize {
		return nil, fmt.Errorf("gogost/gost3410: len(key) != %d", pointSize)
	}
	key := make([]byte, pointSize)
	for i := 0; i < len(key); i++ {
		key[i] = raw[len(raw)-i-1]
	}
	k := bytes2big(key)
	if k.Cmp(zero) == 0 {
		return nil, errors.New("gogost/gost3410: zero private key")
	}
	x, y:= c.ScalarBaseMult(k.Mod(k, c.Q).Bytes())
	pk := PublicKey{c, x, y}
	return &PrivateKey{pk, c, k.Mod(k, c.Q)}, nil
}

func GenPrivateKey(c *Curve, rand io.Reader) (*PrivateKey, error) {
	raw := make([]byte, c.PointSize())
	if _, err := io.ReadFull(rand, raw); err != nil {
		return nil, err
	}
	return NewPrivateKey(c, raw)
}

func (prv *PrivateKey) Raw() []byte {
	raw := pad(prv.Key.Bytes(), prv.C.PointSize())
	// reverse(raw)
	return raw
}

func (prv *PrivateKey) SignDigest(digest []byte, rand io.Reader) ([]byte, error) {
	// 1. Select random nonce k in [1, N-1]
	// 2. Compute kG
	// 3. r = kG.x mod N (kG.x is the x coordinate of the point kG)
	//    Repeat from step 1 if r = 0
	// 4. e = H(m)
	// 5. s = (ke + dr) mod N (GOST 3410-2012)
	//    Repeat from step 1 if s = 0
	// 6. Return (r,s)

	e := bytes2big(digest)
	e.Mod(e, prv.C.Q)
	if e.Cmp(zero) == 0 {
		e = big.NewInt(1)
	}
	kRaw := make([]byte, prv.C.PointSize())
	var err error
	var k *big.Int
	var r *big.Int
	d := big.NewInt(0)
	s := big.NewInt(0)
	for iteration := uint32(0); ; iteration++ {
		if _, err = io.ReadFull(rand, kRaw); err != nil {
			return nil, err
		}
		k = bytes2big(kRaw)
		k.Mod(k, prv.C.Q)
		if k.Cmp(zero) == 0 {
			continue
		}
		r, _ = prv.C.ScalarMult(prv.C.X, prv.C.Y, k.Bytes())
		r.Mod(r, prv.C.Q)
		if r.Cmp(zero) == 0 {
			continue
		}
		d.Mul(prv.Key, r)
		k.Mul(k, e)
		s.Add(d, k)
		s.Mod(s, prv.C.Q)
		if s.Cmp(zero) == 0 {
			continue
		}
		pointSize := prv.C.PointSize()
		return append(
			pad(r.Bytes(), pointSize),
			pad(s.Bytes(), pointSize)...,
		), nil
	}
}

func (prv *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return prv.SignDigest(digest, rand)
}

func (prv *PrivateKey) Public() *PublicKey {
	return &prv.PublicKey
}