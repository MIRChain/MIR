package nist

import (
	"crypto/ecdsa"
	"math/big"
)

type PublicKey struct {
	*ecdsa.PublicKey
}

func (p PublicKey) GetX() *big.Int {
	return p.X
}

func (p PublicKey) GetY() *big.Int {
	return p.Y
}