// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package v5wire

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"hash"
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/csp"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"golang.org/x/crypto/hkdf"
)

const (
	// Encryption/authentication parameters.
	aesKeySize   = 16
	gcmNonceSize = 12
)

// Nonce represents a nonce used for AES/GCM.
type Nonce [gcmNonceSize]byte

// EncodePubkey encodes a public key.
func EncodePubkey[P crypto.PublicKey](key P) []byte {
	// switch key.Curve {
	// case crypto.S256():
	// 	return 
	// default:
	// 	panic("unsupported curve " + key.Curve.Params().Name + " in EncodePubkey")
	// }
	res := crypto.CompressPubkey[P](key)
	return res
}

// DecodePubkey decodes a public key in compressed format.
func DecodePubkey[P crypto.PublicKey](e []byte) (P, error) {
	// switch curve {
	// case crypto.S256():
	// 	if len(e) != 33 {
	// 		return nil, errors.New("wrong size public key data")
	// 	}
	// 	return 
	// default:
	// 	
	// }
	pub, err := crypto.DecompressPubkey[P](e)
	if err != nil {
		return crypto.ZeroPublicKey[P](), fmt.Errorf("error %s in DecodePubkey", err)
	}
	return pub, nil
}

// idNonceHash computes the ID signature hash used in the handshake.
func idNonceHash(h hash.Hash, challenge, ephkey []byte, destID enode.ID) []byte {
	h.Reset()
	h.Write([]byte("discovery v5 identity proof"))
	h.Write(challenge)
	h.Write(ephkey)
	h.Write(destID[:])
	return h.Sum(nil)
}

// makeIDSignature creates the ID nonce signature.
func makeIDSignature[T crypto.PrivateKey, P crypto.PublicKey](hash hash.Hash, key T, challenge, ephkey []byte, destID enode.ID) ([]byte, error) {
	input := idNonceHash(hash, challenge, ephkey, destID)
	// switch key.Curve {
	// case crypto.S256():
	// 	idsig, err := crypto.Sign(input, key)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return idsig[:len(idsig)-1], nil // remove recovery ID
	// default:
	// 	return nil, fmt.Errorf("unsupported curve %s", key.Curve.Params().Name)
	// }
	idsig, err := crypto.Sign(input, key)
	if err != nil {
		return nil, err
	}
	return idsig[:len(idsig)-1], nil // remove recovery ID
}

// s256raw is an unparsed secp256k1 public key ENR entry.
type s256raw []byte

func (s256raw) ENRKey() string { return "secp256k1" }

// verifyIDSignature checks that signature over idnonce was made by the given node.
func verifyIDSignature[P crypto.PublicKey](hash hash.Hash, sig []byte, n *enode.Node[P], challenge, ephkey []byte, destID enode.ID) error {
	switch idscheme := n.Record().IdentityScheme(); idscheme {
	case "v4":
		var pubkey s256raw
		if n.Load(&pubkey) != nil {
			return errors.New("no secp256k1 public key in record")
		}
		input := idNonceHash(hash, challenge, ephkey, destID)
		if !crypto.VerifySignature[P](pubkey, input, sig) {
			return errInvalidNonceSig
		}
		return nil
	default:
		return fmt.Errorf("can't verify ID nonce signature against scheme %q", idscheme)
	}
}

type hashFn func() hash.Hash

// deriveKeys creates the session keys.
func deriveKeys[T crypto.PrivateKey, P crypto.PublicKey](hash hashFn, priv T, pub P, n1, n2 enode.ID, challenge []byte) *session {
	const text = "discovery v5 key agreement"
	var info = make([]byte, 0, len(text)+len(n1)+len(n2))
	info = append(info, text...)
	info = append(info, n1[:]...)
	info = append(info, n2[:]...)

	eph := ecdh[T,P](priv, pub)
	if eph == nil {
		return nil
	}
	kdf := hkdf.New(hash, eph, challenge, info)
	sec := session{writeKey: make([]byte, aesKeySize), readKey: make([]byte, aesKeySize)}
	kdf.Read(sec.writeKey)
	kdf.Read(sec.readKey)
	for i := range eph {
		eph[i] = 0
	}
	return &sec
}

// ecdh creates a shared secret.
func ecdh[T crypto.PrivateKey, P crypto.PublicKey](privkey T, pubkey P) []byte {
	var secX, secY *big.Int 
	switch p:=any(&privkey).(type) {
	case *nist.PrivateKey:
		pub := any(&pubkey).(*nist.PublicKey) 
		secX, secY = pub.ScalarMult(pub.X, pub.Y, p.D.Bytes())
		if secX == nil {
			return nil
		}
	case *gost3410.PrivateKey:
		pub := any(&pubkey).(*gost3410.PublicKey)
		secX, secY = pub.C.ScalarMult(pub.X, pub.Y, p.Key.Bytes())
		sec := make([]byte, 33)
		sec[0] = 0x02 | byte(secY.Bit(0))
		math.ReadBits(secX, sec[1:])
		return sec
	case *csp.Cert:
		// TODO	
	default:
		panic("cant infer priv key")
	}
	sec := make([]byte, 33)
	sec[0] = 0x02 | byte(secY.Bit(0))
	math.ReadBits(secX, sec[1:])
	return sec
}

// encryptGCM encrypts pt using AES-GCM with the given key and nonce. The ciphertext is
// appended to dest, which must not overlap with plaintext. The resulting ciphertext is 16
// bytes longer than plaintext because it contains an authentication tag.
func encryptGCM(dest, key, nonce, plaintext, authData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Errorf("can't create block cipher: %v", err))
	}
	aesgcm, err := cipher.NewGCMWithNonceSize(block, gcmNonceSize)
	if err != nil {
		panic(fmt.Errorf("can't create GCM: %v", err))
	}
	return aesgcm.Seal(dest, nonce, plaintext, authData), nil
}

// decryptGCM decrypts ct using AES-GCM with the given key and nonce.
func decryptGCM(key, nonce, ct, authData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("can't create block cipher: %v", err)
	}
	if len(nonce) != gcmNonceSize {
		return nil, fmt.Errorf("invalid GCM nonce size: %d", len(nonce))
	}
	aesgcm, err := cipher.NewGCMWithNonceSize(block, gcmNonceSize)
	if err != nil {
		return nil, fmt.Errorf("can't create GCM: %v", err)
	}
	pt := make([]byte, 0, len(ct))
	return aesgcm.Open(pt, nonce, ct, authData)
}
