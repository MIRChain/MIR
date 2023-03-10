package nist

import "crypto/ecdsa"

type PrivateKey struct {
	*ecdsa.PrivateKey
}

func (p *PrivateKey)Public() *PublicKey {
	return &PublicKey{&p.PublicKey}
}
