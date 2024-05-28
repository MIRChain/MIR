// Copyright 2016 The go-ethereum Authors
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

package types

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/csp"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/crypto/nist"
	"github.com/MIRChain/MIR/params"
)

var ErrInvalidChainId = errors.New("invalid chain id for signer")

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache[P crypto.PublicKey] struct {
	signer Signer[P]
	from   common.Address
}

// MakeSigner returns a Signer based on the given chain config and block number.
func MakeSigner[P crypto.PublicKey](config *params.ChainConfig, blockNumber *big.Int) Signer[P] {
	var signer Signer[P]
	switch {
	case config.IsBerlin(blockNumber):
		signer = NewEIP2930Signer[P](config.ChainID)
	case config.IsEIP155(blockNumber):
		signer = NewEIP155Signer[P](config.ChainID)
	case config.IsHomestead(blockNumber):
		signer = HomesteadSigner[P]{}
	default:
		signer = FrontierSigner[P]{}
	}
	return signer
}

// LatestSigner returns the 'most permissive' Signer available for the given chain
// configuration. Specifically, this enables support of EIP-155 replay protection and
// EIP-2930 access list transactions when their respective forks are scheduled to occur at
// any block number in the chain config.
//
// Use this in transaction-handling code where the current block number is unknown. If you
// have the current block number available, use MakeSigner instead.
func LatestSigner[P crypto.PublicKey](config *params.ChainConfig) Signer[P] {
	if config.ChainID != nil {
		if config.BerlinBlock != nil || config.YoloV3Block != nil {
			return NewEIP2930Signer[P](config.ChainID)
		}
		if config.EIP155Block != nil {
			return NewEIP155Signer[P](config.ChainID)
		}
	}
	return HomesteadSigner[P]{}
}

// LatestSignerForChainID returns the 'most permissive' Signer available. Specifically,
// this enables support for EIP-155 replay protection and all implemented EIP-2718
// transaction types if chainID is non-nil.
//
// Use this in transaction-handling code where the current block number and fork
// configuration are unknown. If you have a ChainConfig, use LatestSigner instead.
// If you have a ChainConfig and know the current block number, use MakeSigner instead.
func LatestSignerForChainID[P crypto.PublicKey](chainID *big.Int) Signer[P] {
	if chainID == nil {
		return HomesteadSigner[P]{}
	}
	return NewEIP2930Signer[P](chainID)
}

// SignTx signs the transaction using the given signer and private key.
func SignTx[T crypto.PrivateKey, P crypto.PublicKey](tx *Transaction[P], s Signer[P], prv T) (*Transaction[P], error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

// SignNewTx creates a transaction and signs it.
func SignNewTx[T crypto.PrivateKey, P crypto.PublicKey](prv T, s Signer[P], txdata TxData) (*Transaction[P], error) {
	tx := NewTx[P](txdata)
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithSignature(s, sig)
}

// MustSignNewTx creates a transaction and signs it.
// This panics if the transaction cannot be signed.
func MustSignNewTx[T crypto.PrivateKey, P crypto.PublicKey](prv T, s Signer[P], txdata TxData) *Transaction[P] {
	tx, err := SignNewTx(prv, s, txdata)
	if err != nil {
		panic(err)
	}
	return tx
}

// Sender returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// Sender may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Sender[P crypto.PublicKey](signer Signer[P], tx *Transaction[P]) (common.Address, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache[P])
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.from.Store(sigCache[P]{signer: signer, from: addr})
	return addr, nil
}

// Signer encapsulates transaction signature handling. The name of this type is slightly
// misleading because Signers don't actually sign, they're just for validating and
// processing of signatures.
//
// Note that this interface is not a stable API and may change at any time to accommodate
// new protocol rules.
type Signer[P crypto.PublicKey] interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction[P]) (common.Address, error)

	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction[P], sig []byte) (r, s, v *big.Int, err error)
	ChainID() *big.Int

	// Hash returns 'signature hash', i.e. the transaction hash that is signed by the
	// private key. This hash does not uniquely identify the transaction.
	Hash(tx *Transaction[P]) common.Hash

	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer[P]) bool
}

type eip2930Signer[P crypto.PublicKey] struct{ EIP155Signer[P] }

// NewEIP2930Signer returns a signer that accepts EIP-2930 access list transactions,
// EIP-155 replay protected transactions, and legacy Homestead transactions.
func NewEIP2930Signer[P crypto.PublicKey](chainId *big.Int) Signer[P] {
	return eip2930Signer[P]{NewEIP155Signer[P](chainId)}
}

func (s eip2930Signer[P]) ChainID() *big.Int {
	return s.chainId
}

func (s eip2930Signer[P]) Equal(s2 Signer[P]) bool {
	x, ok := s2.(eip2930Signer[P])
	return ok && x.chainId.Cmp(s.chainId) == 0
}

func (s eip2930Signer[P]) Sender(tx *Transaction[P]) (common.Address, error) {
	// Quorum
	if tx.IsPrivate() {
		return QuorumPrivateTxSigner[P]{}.Sender(tx)
	}
	// End Quorum
	V, R, S := tx.RawSignatureValues()
	switch tx.Type() {
	case LegacyTxType:
		if !tx.Protected() {
			return HomesteadSigner[P]{}.Sender(tx)
		}
		V = new(big.Int).Sub(V, s.chainIdMul)
		V.Sub(V, big8)
	case AccessListTxType:
		// ACL txs are defined to use 0 and 1 as their recovery id, add
		// 27 to become equivalent to unprotected Homestead signatures.
		V = new(big.Int).Add(V, big.NewInt(27))
	default:
		return common.Address{}, ErrTxTypeNotSupported
	}
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	return recoverPlain[P](s.Hash(tx), R, S, V, true)
}

func (s eip2930Signer[P]) SignatureValues(tx *Transaction[P], sig []byte) (R, S, V *big.Int, err error) {
	switch txdata := tx.inner.(type) {
	case *LegacyTx:
		return s.EIP155Signer.SignatureValues(tx, sig)
	case *AccessListTx:
		// Check that chain ID of tx matches the signer. We also accept ID zero here,
		// because it indicates that the chain ID was not specified in the tx.
		if txdata.ChainID.Sign() != 0 && txdata.ChainID.Cmp(s.chainId) != 0 {
			return nil, nil, nil, ErrInvalidChainId
		}
		R, S, _ = decodeSignature[P](sig)
		V = big.NewInt(int64(sig[64]))
	default:
		return nil, nil, nil, ErrTxTypeNotSupported
	}
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s eip2930Signer[P]) Hash(tx *Transaction[P]) common.Hash {
	switch tx.Type() {
	case LegacyTxType:
		return rlpHash[P]([]interface{}{
			tx.Nonce(),
			tx.GasPrice(),
			tx.Gas(),
			tx.To(),
			tx.Value(),
			tx.Data(),
			s.chainId, uint(0), uint(0),
		})
	case AccessListTxType:
		return prefixedRlpHash[P](
			tx.Type(),
			[]interface{}{
				s.chainId,
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				tx.AccessList(),
			})
	default:
		// This _should_ not happen, but in case someone sends in a bad
		// json struct via RPC, it's probably more prudent to return an
		// empty hash instead of killing the node with a panic
		//panic("Unsupported transaction type: %d", tx.typ)
		return common.Hash{}
	}
}

// EIP155Signer implements Signer using the EIP-155 rules. This accepts transactions which
// are replay-protected as well as unprotected homestead transactions.
type EIP155Signer[P crypto.PublicKey] struct {
	chainId, chainIdMul *big.Int
}

func NewEIP155Signer[P crypto.PublicKey](chainId *big.Int) EIP155Signer[P] {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return EIP155Signer[P]{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

func (s EIP155Signer[P]) ChainID() *big.Int {
	return s.chainId
}

func (s EIP155Signer[P]) Equal(s2 Signer[P]) bool {
	eip155, ok := s2.(EIP155Signer[P])
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

func (s EIP155Signer[P]) Sender(tx *Transaction[P]) (common.Address, error) {
	if tx.IsPrivate() {
		return QuorumPrivateTxSigner[P]{}.Sender(tx)
	}
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	if !tx.Protected() {
		return HomesteadSigner[P]{}.Sender(tx)
	}
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	V, R, S := tx.RawSignatureValues()
	V = new(big.Int).Sub(V, s.chainIdMul)
	V.Sub(V, big8)
	return recoverPlain[P](s.Hash(tx), R, S, V, true)
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s EIP155Signer[P]) SignatureValues(tx *Transaction[P], sig []byte) (R, S, V *big.Int, err error) {
	if tx.Type() != LegacyTxType {
		return nil, nil, nil, ErrTxTypeNotSupported
	}
	R, S, V = decodeSignature[P](sig)
	if s.chainId.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, s.chainIdMul)
	}
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s EIP155Signer[P]) Hash(tx *Transaction[P]) common.Hash {
	return rlpHash[P]([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		s.chainId, uint(0), uint(0),
	})
}

// HomesteadTransaction implements TransactionInterface using the
// homestead rules.
type HomesteadSigner[P crypto.PublicKey] struct{ FrontierSigner[P] }

func (s HomesteadSigner[P]) ChainID() *big.Int {
	return nil
}

func (s HomesteadSigner[P]) Equal(s2 Signer[P]) bool {
	_, ok := s2.(HomesteadSigner[P])
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (hs HomesteadSigner[P]) SignatureValues(tx *Transaction[P], sig []byte) (r, s, v *big.Int, err error) {
	return hs.FrontierSigner.SignatureValues(tx, sig)
}

func (hs HomesteadSigner[P]) Sender(tx *Transaction[P]) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	v, r, s := tx.RawSignatureValues()
	return recoverPlain[P](hs.Hash(tx), r, s, v, true)
}

type FrontierSigner[P crypto.PublicKey] struct{}

func (s FrontierSigner[P]) ChainID() *big.Int {
	return nil
}

func (s FrontierSigner[P]) Equal(s2 Signer[P]) bool {
	_, ok := s2.(FrontierSigner[P])
	return ok
}

func (fs FrontierSigner[P]) Sender(tx *Transaction[P]) (common.Address, error) {
	if tx.Type() != LegacyTxType {
		return common.Address{}, ErrTxTypeNotSupported
	}
	v, r, s := tx.RawSignatureValues()
	return recoverPlain[P](fs.Hash(tx), r, s, v, false)
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs FrontierSigner[P]) SignatureValues(tx *Transaction[P], sig []byte) (r, s, v *big.Int, err error) {
	if tx.Type() != LegacyTxType {
		return nil, nil, nil, ErrTxTypeNotSupported
	}
	r, s, v = decodeSignature[P](sig)
	return r, s, v, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (fs FrontierSigner[P]) Hash(tx *Transaction[P]) common.Hash {
	return rlpHash[P]([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
	})
}

func decodeSignature[P crypto.PublicKey](sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	var pub P
	switch any(&pub).(type) {
	case *nist.PublicKey:
		r = new(big.Int).SetBytes(sig[:32])
		s = new(big.Int).SetBytes(sig[32:64])
		v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	case *gost3410.PublicKey:
		r = new(big.Int).SetBytes(sig[:32])
		s = new(big.Int).SetBytes(sig[32:64])
		v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	case *csp.PublicKey:
		resSig := sig[:64]
		reverse(resSig)
		r = new(big.Int).SetBytes(resSig[32:64])
		s = new(big.Int).SetBytes(resSig[:32])
		v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	}
	return r, s, v
}

func recoverPlain[P crypto.PublicKey](sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	if Vb.BitLen() > 8 {
		return common.Address{}, ErrInvalidSig
	}
	var offset uint64
	// private transaction has a v value of 37 or 38
	if isPrivate(Vb) {
		offset = 37
	} else {
		offset = 27
	}
	V := byte(Vb.Uint64() - offset)
	if !crypto.ValidateSignatureValues[P](V, R, S, homestead) {
		return common.Address{}, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	sig := make([]byte, crypto.SignatureLength)
	// MIR
	var p P
	switch any(&p).(type) {
	case *nist.PublicKey:
		r, s := R.Bytes(), S.Bytes()
		copy(sig[32-len(r):32], r)
		copy(sig[64-len(s):64], s)
		sig[64] = V
	case *gost3410.PublicKey:
		r, s := R.Bytes(), S.Bytes()
		copy(sig[32-len(r):32], r)
		copy(sig[64-len(s):64], s)
		sig[64] = V
	case *csp.PublicKey:
		resSig := make([]byte, 64)
		r, s := R.Bytes(), S.Bytes()
		copy(resSig[32-len(s):32], s)
		copy(resSig[64-len(r):64], r)
		reverse(resSig)
		copy(sig, resSig)
		sig[64] = V
	}
	// recover the public key from the signature
	pub, err := crypto.Ecrecover[P](sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || !(pub[0] == 4 || pub[0] == 64) {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256[P](pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		// TODO(joel): this given v = 37 / 38 this constrains us to chain id 1
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}

func reverse(d []byte) {
	for i, j := 0, len(d)-1; i < j; i, j = i+1, j-1 {
		d[i], d[j] = d[j], d[i]
	}
}
