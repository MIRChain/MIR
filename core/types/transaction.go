// Copyright 2014 The go-ethereum Authors
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
	"bytes"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/private/engine"
	"github.com/pavelkrolevets/MIR-pro/rlp"
)

var (
	ErrInvalidSig           = errors.New("invalid transaction v, r, s values")
	ErrUnexpectedProtection = errors.New("transaction type does not supported EIP-155 protected signatures")
	ErrInvalidTxType        = errors.New("transaction type not valid in this context")
	ErrTxTypeNotSupported   = errors.New("transaction type not supported")
	errEmptyTypedTx         = errors.New("empty typed transaction bytes")
)

// Transaction types.
const (
	LegacyTxType = iota
	AccessListTxType
)

// Transaction is an Ethereum transaction.
type Transaction [P crypto.PublicKey] struct {
	inner TxData    // Consensus contents of a transaction
	time  time.Time // Time first seen locally (spam avoidance)

	// caches
	hash atomic.Value
	size atomic.Value
	from atomic.Value

	privacyMetadata *PrivacyMetadata
}

// NewTx creates a new transaction.
func NewTx[P crypto.PublicKey](inner TxData) *Transaction[P] {
	tx := new(Transaction[P])
	tx.setDecoded(inner.copy(), 0)
	return tx
}

// TxData is the underlying data of a transaction.
//
// This is implemented by LegacyTx and AccessListTx.
type TxData interface {
	txType() byte // returns the type ID
	copy() TxData // creates a deep copy and initializes all fields

	chainID() *big.Int
	accessList() AccessList
	data() []byte
	gas() uint64
	gasPrice() *big.Int
	value() *big.Int
	nonce() uint64
	to() *common.Address

	rawSignatureValues() (v, r, s *big.Int)
	setSignatureValues(chainID, v, r, s *big.Int)
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction[P]) EncodeRLP(w io.Writer) error {
	if tx.Type() == LegacyTxType {
		return rlp.Encode(w, tx.inner)
	}
	// It's an EIP-2718 typed TX envelope.
	buf := encodeBufferPool.Get().(*bytes.Buffer)
	defer encodeBufferPool.Put(buf)
	buf.Reset()
	if err := tx.encodeTyped(buf); err != nil {
		return err
	}
	return rlp.Encode(w, buf.Bytes())
}

// encodeTyped writes the canonical encoding of a typed transaction to w.
func (tx *Transaction[P]) encodeTyped(w *bytes.Buffer) error {
	w.WriteByte(tx.Type())
	return rlp.Encode(w, tx.inner)
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For EIP-2718 typed
// transactions, it returns the type and payload.
func (tx *Transaction[P]) MarshalBinary() ([]byte, error) {
	if tx.Type() == LegacyTxType {
		return rlp.EncodeToBytes(tx.inner)
	}
	var buf bytes.Buffer
	err := tx.encodeTyped(&buf)
	return buf.Bytes(), err
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction[P]) DecodeRLP(s *rlp.Stream) error {
	kind, size, err := s.Kind()
	switch {
	case err != nil:
		return err
	case kind == rlp.List:
		// It's a legacy transaction.
		var inner LegacyTx
		err := s.Decode(&inner)
		if err == nil {
			tx.setDecoded(&inner, int(rlp.ListSize(size)))
		}
		return err
	case kind == rlp.String:
		// It's an EIP-2718 typed TX envelope.
		var b []byte
		if b, err = s.Bytes(); err != nil {
			return err
		}
		inner, err := tx.decodeTyped(b)
		if err == nil {
			tx.setDecoded(inner, len(b))
		}
		return err
	default:
		return rlp.ErrExpectedList
	}
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP2718 typed transactions.
func (tx *Transaction[P]) UnmarshalBinary(b []byte) error {
	if len(b) > 0 && b[0] > 0x7f {
		// It's a legacy transaction.
		var data LegacyTx
		err := rlp.DecodeBytes(b, &data)
		if err != nil {
			return err
		}
		tx.setDecoded(&data, len(b))
		return nil
	}
	// It's an EIP2718 typed transaction envelope.
	inner, err := tx.decodeTyped(b)
	if err != nil {
		return err
	}
	tx.setDecoded(inner, len(b))
	return nil
}

// decodeTyped decodes a typed transaction from the canonical format.
func (tx *Transaction[P]) decodeTyped(b []byte) (TxData, error) {
	if len(b) == 0 {
		return nil, errEmptyTypedTx
	}
	switch b[0] {
	case AccessListTxType:
		var inner AccessListTx
		err := rlp.DecodeBytes(b[1:], &inner)
		return &inner, err
	default:
		return nil, ErrTxTypeNotSupported
	}
}

// setDecoded sets the inner transaction and size after decoding.
func (tx *Transaction[P]) setDecoded(inner TxData, size int) {
	tx.inner = inner
	tx.time = time.Now()
	if size > 0 {
		tx.size.Store(common.StorageSize(size))
	}
}

func sanityCheckSignature[P crypto.PublicKey](v *big.Int, r *big.Int, s *big.Int, maybeProtected bool) error {
	if isProtectedV(v) && !maybeProtected {
		return ErrUnexpectedProtection
	}

	var plainV byte
	if isProtectedV(v) {
		chainID := deriveChainId(v).Uint64()
		plainV = byte(v.Uint64() - 35 - 2*chainID)
	} else if maybeProtected {
		// Only EIP-155 signatures can be optionally protected. Since
		// we determined this v value is not protected, it must be a
		// raw 27 or 28.
		plainV = byte(v.Uint64() - 27)
	} else {
		// If the signature is not optionally protected, we assume it
		// must already be equal to the recovery id.
		plainV = byte(v.Uint64())
	}
	if !crypto.ValidateSignatureValues[P](plainV, r, s, false) {
		return ErrInvalidSig
	}

	return nil
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28 && v != 1 && v != 0
	}
	// anything not 27 or 28 is considered protected
	return true
}

// Protected says whether the transaction is replay-protected.
func (tx *Transaction[P]) Protected() bool {
	switch tx := tx.inner.(type) {
	case *LegacyTx:
		return tx.V != nil && isProtectedV(tx.V)
	default:
		return true
	}
}

// Type returns the transaction type.
func (tx *Transaction[P]) Type() uint8 {
	return tx.inner.txType()
}

// ChainId returns the EIP155 chain ID of the transaction. The return value will always be
// non-nil. For legacy transactions which are not replay-protected, the return value is
// zero.
func (tx *Transaction[P]) ChainId() *big.Int {
	return tx.inner.chainID()
}

// Data returns the input data of the transaction.
func (tx *Transaction[P]) Data() []byte { return tx.inner.data() }

// AccessList returns the access list of the transaction.
func (tx *Transaction[P]) AccessList() AccessList { return tx.inner.accessList() }

// Gas returns the gas limit of the transaction.
func (tx *Transaction[P]) Gas() uint64 { return tx.inner.gas() }

// GasPrice returns the gas price of the transaction.
func (tx *Transaction[P]) GasPrice() *big.Int { return new(big.Int).Set(tx.inner.gasPrice()) }

// Value returns the ether amount of the transaction.
func (tx *Transaction[P]) Value() *big.Int { return new(big.Int).Set(tx.inner.value()) }

// Nonce returns the sender account nonce of the transaction.
func (tx *Transaction[P]) Nonce() uint64 { return tx.inner.nonce() }

// To returns the recipient address of the transaction.
// For contract-creation transactions, To returns nil.
func (tx *Transaction[P]) To() *common.Address {
	// Copy the pointed-to address.
	ito := tx.inner.to()
	if ito == nil {
		return nil
	}
	cpy := *ito
	return &cpy
}

// Cost returns gas * gasPrice + value.
func (tx *Transaction[P]) Cost() *big.Int {
	total := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(tx.Gas()))
	total.Add(total, tx.Value())
	return total
}

// RawSignatureValues returns the V, R, S signature values of the transaction.
// The return values should not be modified by the caller.
func (tx *Transaction[P]) RawSignatureValues() (v, r, s *big.Int) {
	return tx.inner.rawSignatureValues()
}

// GasPriceCmp compares the gas prices of two transactions.
func (tx *Transaction[P]) GasPriceCmp(other *Transaction[P]) int {
	return tx.inner.gasPrice().Cmp(other.inner.gasPrice())
}

// GasPriceIntCmp compares the gas price of the transaction against the given price.
func (tx *Transaction[P]) GasPriceIntCmp(other *big.Int) int {
	return tx.inner.gasPrice().Cmp(other)
}

// Hash returns the transaction hash.
func (tx *Transaction[P]) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	var h common.Hash
	if tx.Type() == LegacyTxType {
		h = rlpHash[P](tx.inner)
	} else {
		h = prefixedRlpHash[P](tx.Type(), tx.inner)
	}
	tx.hash.Store(h)
	return h
}

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previously cached value.
func (tx *Transaction[P]) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.inner)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be in the [R || S || V] format where V is 0 or 1.
func (tx *Transaction[P]) WithSignature(signer Signer[P], sig []byte) (*Transaction[P], error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := tx.inner.copy()
	cpy.setSignatureValues(signer.ChainID(), v, r, s)
	return &Transaction[P]{inner: cpy, time: tx.time}, nil
}

// Transactions implements DerivableList for transactions.
type Transactions[P crypto.PublicKey] []*Transaction[P]

// Len returns the length of s.
func (s Transactions[P]) Len() int { return len(s) }

// EncodeIndex encodes the i'th transaction to w. Note that this does not check for errors
// because we assume that *Transaction will only ever contain valid txs that were either
// constructed by decoding or via public API in this package.
func (s Transactions[P]) EncodeIndex(i int, w *bytes.Buffer) {
	tx := s[i]
	if tx.Type() == LegacyTxType {
		rlp.Encode(w, tx.inner)
	} else {
		tx.encodeTyped(w)
	}
}

// TxDifference returns a new set which is the difference between a and b.
func TxDifference[P crypto.PublicKey](a, b Transactions[P]) Transactions[P] {
	keep := make(Transactions[P], 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce[P crypto.PublicKey] Transactions[P]

func (s TxByNonce[P]) Len() int           { return len(s) }
func (s TxByNonce[P]) Less(i, j int) bool { return s[i].Nonce() < s[j].Nonce() }
func (s TxByNonce[P]) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// TxByPriceAndTime implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type TxByPriceAndTime[P crypto.PublicKey] Transactions[P]

func (s TxByPriceAndTime[P]) Len() int { return len(s) }
func (s TxByPriceAndTime[P]) Less(i, j int) bool {
	// If the prices are equal, use the time the transaction was first seen for
	// deterministic sorting
	cmp := s[i].GasPrice().Cmp(s[j].GasPrice())
	if cmp == 0 {
		return s[i].time.Before(s[j].time)
	}
	return cmp > 0
}
func (s TxByPriceAndTime[P]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *TxByPriceAndTime[P]) Push(x interface{}) {
	*s = append(*s, x.(*Transaction[P]))
}

func (s *TxByPriceAndTime[P]) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// TransactionsByPriceAndNonce represents a set of transactions that can return
// transactions in a profit-maximizing sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByPriceAndNonce [P crypto.PublicKey] struct {
	txs    map[common.Address]Transactions[P] // Per account nonce-sorted list of transactions
	heads  TxByPriceAndTime[P]               // Next transaction for each unique account (price heap)
	signer Signer[P]                      // Signer for the set of transactions
}

// NewTransactionsByPriceAndNonce creates a transaction set that can retrieve
// price sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.
func NewTransactionsByPriceAndNonce[P crypto.PublicKey](signer Signer[P], txs map[common.Address]Transactions[P]) *TransactionsByPriceAndNonce[P] {
	// Initialize a price and received time based heap with the head transactions
	heads := make(TxByPriceAndTime[P], 0, len(txs))
	for from, accTxs := range txs {
		// Ensure the sender address is from the signer
		acc, err := Sender[P](signer, accTxs[0])
		if err != nil {
			log.Error("Failed to retrieve the sender address", "err", err)
		}
		if acc != from {
			delete(txs, from)
			continue
		}
		heads = append(heads, accTxs[0])
		txs[from] = accTxs[1:]
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	return &TransactionsByPriceAndNonce[P]{
		txs:    txs,
		heads:  heads,
		signer: signer,
	}
}

// Peek returns the next transaction by price.
func (t *TransactionsByPriceAndNonce[P]) Peek() *Transaction[P] {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0]
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByPriceAndNonce[P]) Shift() {
	acc, _ := Sender[P](t.signer, t.heads[0])
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		t.heads[0], t.txs[acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByPriceAndNonce[P]) Pop() {
	heap.Pop(&t.heads)
}

// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	to         *common.Address
	from       common.Address
	nonce      uint64
	amount     *big.Int
	gasLimit   uint64
	gasPrice   *big.Int
	data       []byte
	accessList AccessList
	checkNonce bool
	// Quorum
	isPrivate      bool
	isInnerPrivate bool
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, accessList AccessList, checkNonce bool) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		gasLimit:   gasLimit,
		gasPrice:   gasPrice,
		data:       data,
		accessList: accessList,
		checkNonce: checkNonce,
	}
}

// AsMessage returns the transaction as a core.Message.
func (tx *Transaction[P]) AsMessage(s Signer[P]) (Message, error) {
	msg := Message{
		nonce:      tx.Nonce(),
		gasLimit:   tx.Gas(),
		gasPrice:   new(big.Int).Set(tx.GasPrice()),
		to:         tx.To(),
		amount:     tx.Value(),
		data:       tx.Data(),
		accessList: tx.AccessList(),
		checkNonce: true,
		// Quorum
		isPrivate: tx.IsPrivate(),
	}

	var err error
	msg.from, err = Sender[P](s, tx)
	return msg, err
}

func (m Message) From() common.Address   { return m.from }
func (m Message) To() *common.Address    { return m.to }
func (m Message) GasPrice() *big.Int     { return m.gasPrice }
func (m Message) Value() *big.Int        { return m.amount }
func (m Message) Gas() uint64            { return m.gasLimit }
func (m Message) Nonce() uint64          { return m.nonce }
func (m Message) Data() []byte           { return m.data }
func (m Message) AccessList() AccessList { return m.accessList }
func (m Message) CheckNonce() bool       { return m.checkNonce }

// Quorum

func NewTxPrivacyMetadata(privacyFlag engine.PrivacyFlagType) *PrivacyMetadata {
	return &PrivacyMetadata{
		PrivacyFlag: privacyFlag,
	}
}

func (tx *Transaction[P]) SetTxPrivacyMetadata(pm *PrivacyMetadata) {
	tx.privacyMetadata = pm
}

// PrivacyMetadata returns the privacy metadata of the transaction. (Quorum)
func (tx *Transaction[P]) PrivacyMetadata() *PrivacyMetadata {
	return tx.privacyMetadata
}

// From returns the sender address of the transaction. (Quorum)
func (tx *Transaction[P]) From() common.Address {
	if from, err := Sender[P](NewEIP2930Signer[P](tx.ChainId()), tx); err == nil {
		return from
	}
	return common.Address{}
}

// String returns the string representation of the transaction. (Quorum)
func (tx *Transaction[P]) String() string {
	var from, to string
	v, r, s := tx.RawSignatureValues()
	if v != nil {
		if f, err := Sender[P](NewEIP2930Signer[P](tx.ChainId()), tx); err != nil {
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if tx.To() == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", tx.To())
	}
	enc, _ := rlp.EncodeToBytes(&tx.inner)
	return fmt.Sprintf(`
	TX(%x)
	Contract: %v
	From:     %s
	To:       %s
	Nonce:    %v
	GasPrice: %#x
	GasLimit  %#x
	Value:    %#x
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x
	Hex:      %x
`,
		tx.Hash(),
		tx.To() == nil,
		from,
		to,
		tx.Nonce(),
		tx.Cost(),
		tx.Gas(),
		tx.Value(),
		tx.Data(),
		v,
		r,
		s,
		enc,
	)
}

func (m Message) IsPrivate() bool {
	return m.isPrivate
}

// Quorum
func (m Message) IsInnerPrivate() bool {
	return m.isInnerPrivate
}

// Quorum
func (m Message) WithInnerPrivateFlag(isInnerPrivateTxn bool) Message {
	m.isInnerPrivate = isInnerPrivateTxn
	return m
}

// overriding msg.data so that when tesseera.receive is invoked we get nothing back
func (m Message) WithEmptyPrivateData(b bool) Message {
	if b {
		m.data = common.EncryptedPayloadHash{}.Bytes()
	}
	return m
}

func (tx *Transaction[P]) IsPrivate() bool {
	v, _, _ := tx.RawSignatureValues()
	if v == nil {
		return false
	}
	return v.Uint64() == 37 || v.Uint64() == 38
}

/*
 * Indicates that a transaction is private, but doesn't necessarily set the correct v value, as it can be called on
 * an unsigned transaction.
 * pre homestead signer, all v values were v=27 or v=28, with EIP155Signer that change,
 * but SetPrivate() is also used on unsigned transactions to temporarily set the v value to indicate
 * the transaction is intended to be private, and so that the correct signer can be selected. The signer will correctly
 * set the valid v value (37 or 38): This helps minimize changes vs upstream go-ethereum code.
 */
func (tx *Transaction[P]) SetPrivate() {
	v, _, _ := tx.RawSignatureValues()
	if tx.IsPrivate() {
		return
	}
	if v.Int64() == 28 {
		v.SetUint64(38)
	} else {
		v.SetUint64(37)
	}
}

func (tx *Transaction[P]) IsPrivacyMarker() bool {
	return tx.To() != nil && *tx.To() == common.QuorumPrivacyPrecompileContractAddress()
}

// PrivacyMetadata encapsulates privacy information to be attached
// to a transaction being processed
type PrivacyMetadata struct {
	PrivacyFlag engine.PrivacyFlagType
}

// End Quorum
