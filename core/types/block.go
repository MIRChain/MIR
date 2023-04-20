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

// Package types contains data types related to Ethereum consensus.
package types

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/rlp"
)

var (
	EmptyRootHash  = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

func EmptyUncleHash[P crypto.PublicKey]() common.Hash {return rlpHash[P]([]*Header[P](nil))}

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

//go:generate gencodec -type Header -field-override headerMarshaling -out gen_header_json.go

// Header represents a block header in the Ethereum blockchain.
type Header [P crypto.PublicKey] struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom[P]       `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`
}

// field type overrides for gencodec
type headerMarshaling struct {
	Difficulty *hexutil.Big
	Number     *hexutil.Big
	GasLimit   hexutil.Uint64
	GasUsed    hexutil.Uint64
	Time       hexutil.Uint64
	Extra      hexutil.Bytes
	Hash       common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header[P]) Hash() common.Hash {
	// If the mix digest is equivalent to the predefined Istanbul digest, use Istanbul
	// specific hash calculation.
	if h.MixDigest == IstanbulDigest {
		// Seal is reserved in extra-data. To prove block is signed by the proposer.
		if istanbulHeader := FilteredHeader(h); istanbulHeader != nil {
			return rlpHash[P](istanbulHeader)
		}
	}
	return rlpHash[P](h)
}

// var headerSize = common.StorageSize(reflect.TypeOf(Header{}).Size())

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header[P]) Size() common.StorageSize {
	return common.StorageSize(reflect.TypeOf(Header[P]{}).Size()) + common.StorageSize(len(h.Extra)+(h.Difficulty.BitLen()+h.Number.BitLen())/8)
}

// SanityCheck checks a few basic things -- these checks are way beyond what
// any 'sane' production values should hold, and can mainly be used to prevent
// that the unbounded fields are stuffed with junk data to add processing
// overhead
func (h *Header[P]) SanityCheck() error {
	if h.Number != nil && !h.Number.IsUint64() {
		return fmt.Errorf("too large block number: bitlen %d", h.Number.BitLen())
	}
	if h.Difficulty != nil {
		if diffLen := h.Difficulty.BitLen(); diffLen > 80 {
			return fmt.Errorf("too large block difficulty: bitlen %d", diffLen)
		}
	}
	if eLen := len(h.Extra); eLen > 100*1024 {
		return fmt.Errorf("too large block extradata: size %d", eLen)
	}
	return nil
}

// QBFTHashWithRoundNumber gets the hash of the Header with Only commit seal set to its null value
func (h *Header[P]) QBFTHashWithRoundNumber(round uint32) common.Hash {
	return rlpHash[P](QBFTFilteredHeaderWithRound(h, round))
}

// EmptyBody returns true if there is no additional 'body' to complete the header
// that is: no transactions and no uncles.
func (h *Header[P]) EmptyBody() bool {
	return h.TxHash == EmptyRootHash && h.UncleHash ==  rlpHash[P]([]*Header[P](nil))
}

// EmptyReceipts returns true if there are no receipts for this header/block.
func (h *Header[P]) EmptyReceipts() bool {
	return h.ReceiptHash == EmptyRootHash
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions and uncles) together.
type Body [P crypto.PublicKey]struct {
	Transactions []*Transaction[P]
	Uncles       []*Header[P]
}

// Block represents an entire block in the Ethereum blockchain.
type Block [P crypto.PublicKey] struct {
	header       *Header[P]
	uncles       []*Header[P]
	transactions Transactions[P]

	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

func (b *Block[P]) String() string {
	return fmt.Sprintf("{Header: %v}", b.header)
}

// "external" block encoding. used for eth protocol, etc.
type extblock [P crypto.PublicKey] struct {
	Header *Header[P]
	Txs    []*Transaction[P]
	Uncles []*Header[P]
}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, UncleHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs, uncles
// and receipts.
func NewBlock[P crypto.PublicKey](header *Header[P], txs []*Transaction[P], uncles []*Header[P], receipts []*Receipt[P], hasher TrieHasher) *Block[P] {
	b := &Block[P]{header: CopyHeader(header), td: new(big.Int)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		b.header.TxHash = DeriveSha(Transactions[P](txs), hasher)
		b.transactions = make(Transactions[P], len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = EmptyRootHash
	} else {
		b.header.ReceiptHash = DeriveSha(Receipts[P](receipts), hasher)
		b.header.Bloom = CreateBloom(receipts)
	}

	if len(uncles) == 0 {
		b.header.UncleHash =  rlpHash[P]([]*Header[P](nil))
	} else {
		b.header.UncleHash = CalcUncleHash(uncles)
		b.uncles = make([]*Header[P], len(uncles))
		for i := range uncles {
			b.uncles[i] = CopyHeader(uncles[i])
		}
	}

	return b
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader[P crypto.PublicKey](header *Header[P]) *Block[P] {
	return &Block[P]{header: CopyHeader(header)}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader[P crypto.PublicKey](h *Header[P]) *Header[P] {
	cpy := *h
	if cpy.Difficulty = new(big.Int); h.Difficulty != nil {
		cpy.Difficulty.Set(h.Difficulty)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	return &cpy
}

// DecodeRLP decodes the Ethereum
func (b *Block[P]) DecodeRLP(s *rlp.Stream) error {
	var eb extblock[P]
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.uncles, b.transactions = eb.Header, eb.Uncles, eb.Txs
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block[P]) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock[P]{
		Header: b.header,
		Txs:    b.transactions,
		Uncles: b.uncles,
	})
}

// TODO: copies

func (b *Block[P]) Uncles() []*Header[P]          { return b.uncles }
func (b *Block[P]) Transactions() Transactions[P] { return b.transactions }

func (b *Block[P]) Transaction(hash common.Hash) *Transaction[P] {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

func (b *Block[P]) Number() *big.Int     { return new(big.Int).Set(b.header.Number) }
func (b *Block[P]) GasLimit() uint64     { return b.header.GasLimit }
func (b *Block[P]) GasUsed() uint64      { return b.header.GasUsed }
func (b *Block[P]) Difficulty() *big.Int { return new(big.Int).Set(b.header.Difficulty) }
func (b *Block[P]) Time() uint64         { return b.header.Time }

func (b *Block[P]) NumberU64() uint64        { return b.header.Number.Uint64() }
func (b *Block[P]) MixDigest() common.Hash   { return b.header.MixDigest }
func (b *Block[P]) Nonce() uint64            { return binary.BigEndian.Uint64(b.header.Nonce[:]) }
func (b *Block[P]) Bloom() Bloom[P]             { return b.header.Bloom }
func (b *Block[P]) Coinbase() common.Address { return b.header.Coinbase }
func (b *Block[P]) Root() common.Hash        { return b.header.Root }
func (b *Block[P]) ParentHash() common.Hash  { return b.header.ParentHash }
func (b *Block[P]) TxHash() common.Hash      { return b.header.TxHash }
func (b *Block[P]) ReceiptHash() common.Hash { return b.header.ReceiptHash }
func (b *Block[P]) UncleHash() common.Hash   { return b.header.UncleHash }
func (b *Block[P]) Extra() []byte            { return common.CopyBytes(b.header.Extra) }

func (b *Block[P]) Header() *Header[P] { return CopyHeader(b.header) }

// Body returns the non-header content of the block.
func (b *Block[P]) Body() *Body[P] { return &Body[P]{b.transactions, b.uncles} }

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block[P]) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// SanityCheck can be used to prevent that unbounded fields are
// stuffed with junk data to add processing overhead
func (b *Block[P]) SanityCheck() error {
	return b.header.SanityCheck()
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

func CalcUncleHash[P crypto.PublicKey](uncles []*Header[P]) common.Hash {
	if len(uncles) == 0 {
		return  rlpHash[P]([]*Header[P](nil))
	}
	return rlpHash[P](uncles)
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block[P]) WithSeal(header *Header[P]) *Block[P] {
	cpy := *header

	return &Block[P]{
		header:       &cpy,
		transactions: b.transactions,
		uncles:       b.uncles,
	}
}

// WithBody returns a new block with the given transaction and uncle contents.
func (b *Block[P]) WithBody(transactions []*Transaction[P], uncles []*Header[P]) *Block[P] {
	block := &Block[P]{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction[P], len(transactions)),
		uncles:       make([]*Header[P], len(uncles)),
	}
	copy(block.transactions, transactions)
	for i := range uncles {
		block.uncles[i] = CopyHeader(uncles[i])
	}
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block[P]) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

type Blocks[P crypto.PublicKey] []*Block[P]
