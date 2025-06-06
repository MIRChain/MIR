// Copyright 2019 The go-ethereum Authors
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

// Package graphql provides a GraphQL interface to Ethereum node data.
package graphql

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/pavelkrolevets/MIR-pro"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/eth/filters"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

var (
	errBlockInvariant = errors.New("block objects must be instantiated with at least one of num or hash")
)

type Long int64

// ImplementsGraphQLType returns true if Long implements the provided GraphQL type.
func (b Long) ImplementsGraphQLType(name string) bool { return name == "Long" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (b *Long) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		// uncomment to support hex values
		//if strings.HasPrefix(input, "0x") {
		//	// apply leniency and support hex representations of longs.
		//	value, err := hexutil.DecodeUint64(input)
		//	*b = Long(value)
		//	return err
		//} else {
		value, err := strconv.ParseInt(input, 10, 64)
		*b = Long(value)
		return err
		//}
	case int32:
		*b = Long(input)
	case int64:
		*b = Long(input)
	default:
		err = fmt.Errorf("unexpected type %T for Long", input)
	}
	return err
}

// Account represents an Ethereum account at a particular block.
type Account [T crypto.PrivateKey, P crypto.PublicKey] struct {
	backend       ethapi.Backend[T,P]
	address       common.Address
	blockNrOrHash rpc.BlockNumberOrHash
}

// getState fetches the StateDB object for an account.
func (a *Account[T,P]) getState(ctx context.Context) (vm.MinimalApiState, error) {
	stat, _, err := a.backend.StateAndHeaderByNumberOrHash(ctx, a.blockNrOrHash)
	return stat, err
}

func (a *Account[T,P]) Address(ctx context.Context) (common.Address, error) {
	return a.address, nil
}

func (a *Account[T,P]) Balance(ctx context.Context) (hexutil.Big, error) {
	state, err := a.getState(ctx)
	if err != nil {
		return hexutil.Big{}, err
	}
	return hexutil.Big(*state.GetBalance(a.address)), nil
}

func (a *Account[T,P]) TransactionCount(ctx context.Context) (hexutil.Uint64, error) {
	state, err := a.getState(ctx)
	if err != nil {
		return 0, err
	}
	return hexutil.Uint64(state.GetNonce(a.address)), nil
}

func (a *Account[T,P]) Code(ctx context.Context) (hexutil.Bytes, error) {
	state, err := a.getState(ctx)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	return state.GetCode(a.address), nil
}

func (a *Account[T,P]) Storage(ctx context.Context, args struct{ Slot common.Hash }) (common.Hash, error) {
	state, err := a.getState(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return state.GetState(a.address, args.Slot), nil
}

// Log represents an individual log message. All arguments are mandatory.
type Log [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	backend     ethapi.Backend[T,P]
	transaction *Transaction[T,P]
	log         *types.Log
}

func (l *Log[T,P]) Transaction(ctx context.Context) *Transaction[T,P] {
	return l.transaction
}

func (l *Log[T,P]) Account(ctx context.Context, args BlockNumberArgs) *Account[T,P] {
	return &Account[T,P]{
		backend:       l.backend,
		address:       l.log.Address,
		blockNrOrHash: args.NumberOrLatest(),
	}
}

func (l *Log[T,P]) Index(ctx context.Context) int32 {
	return int32(l.log.Index)
}

func (l *Log[T,P]) Topics(ctx context.Context) []common.Hash {
	return l.log.Topics
}

func (l *Log[T,P]) Data(ctx context.Context) hexutil.Bytes {
	return l.log.Data
}

// AccessTuple represents EIP-2930
type AccessTuple struct {
	address     common.Address
	storageKeys *[]common.Hash
}

func (at *AccessTuple) Address(ctx context.Context) common.Address {
	return at.address
}

func (at *AccessTuple) StorageKeys(ctx context.Context) *[]common.Hash {
	return at.storageKeys
}

// Transaction represents an Ethereum transaction.
// backend and hash are mandatory; all others will be fetched when required.
type Transaction [T crypto.PrivateKey, P crypto.PublicKey] struct {
	backend       ethapi.Backend[T,P]
	hash          common.Hash
	tx            *types.Transaction[P]
	block         *Block[T,P]
	index         uint64
	receiptGetter receiptGetter[T,P]
}

// resolve returns the internal transaction object, fetching it if needed.
func (t *Transaction[T,P]) resolve(ctx context.Context) (*types.Transaction[P], error) {
	if t.tx == nil {
		tx, blockHash, _, index := rawdb.ReadTransaction[P](t.backend.ChainDb(), t.hash)
		if tx != nil {
			t.tx = tx
			blockNrOrHash := rpc.BlockNumberOrHashWithHash(blockHash, false)
			t.block = &Block[T,P]{
				backend:      t.backend,
				numberOrHash: &blockNrOrHash,
			}
			t.index = index
		} else {
			t.tx = t.backend.GetPoolTransaction(t.hash)
		}
	}
	return t.tx, nil
}

func (t *Transaction[T,P]) Hash(ctx context.Context) common.Hash {
	return t.hash
}

func (t *Transaction[T,P]) InputData(ctx context.Context) (hexutil.Bytes, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return hexutil.Bytes{}, err
	}
	return tx.Data(), nil
}

func (t *Transaction[T,P]) Gas(ctx context.Context) (hexutil.Uint64, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return 0, err
	}
	return hexutil.Uint64(tx.Gas()), nil
}

func (t *Transaction[T,P]) GasPrice(ctx context.Context) (hexutil.Big, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return hexutil.Big{}, err
	}
	return hexutil.Big(*tx.GasPrice()), nil
}

func (t *Transaction[T,P]) Value(ctx context.Context) (hexutil.Big, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return hexutil.Big{}, err
	}
	return hexutil.Big(*tx.Value()), nil
}

func (t *Transaction[T,P]) Nonce(ctx context.Context) (hexutil.Uint64, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return 0, err
	}
	return hexutil.Uint64(tx.Nonce()), nil
}

func (t *Transaction[T,P]) To(ctx context.Context, args BlockNumberArgs) (*Account[T,P], error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return nil, err
	}
	to := tx.To()
	if to == nil {
		return nil, nil
	}
	return &Account[T,P]{
		backend:       t.backend,
		address:       *to,
		blockNrOrHash: args.NumberOrLatest(),
	}, nil
}

func (t *Transaction[T,P]) From(ctx context.Context, args BlockNumberArgs) (*Account[T,P], error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return nil, err
	}
	signer := types.LatestSigner[P](t.backend.ChainConfig())
	from, _ := types.Sender(signer, tx)
	return &Account[T,P]{
		backend:       t.backend,
		address:       from,
		blockNrOrHash: args.NumberOrLatest(),
	}, nil
}

func (t *Transaction[T,P]) Block(ctx context.Context) (*Block[T,P], error) {
	if _, err := t.resolve(ctx); err != nil {
		return nil, err
	}
	return t.block, nil
}

func (t *Transaction[T,P]) Index(ctx context.Context) (*int32, error) {
	if _, err := t.resolve(ctx); err != nil {
		return nil, err
	}
	if t.block == nil {
		return nil, nil
	}
	index := int32(t.index)
	return &index, nil
}

// (Quorum) receiptGetter allows Transaction to have different behaviours for getting transaction receipts
// (e.g. getting standard receipts or privacy precompile receipts from the db)
type receiptGetter [T crypto.PrivateKey, P crypto.PublicKey] interface {
	get(ctx context.Context) (*types.Receipt[P], error)
}

// (Quorum) transactionReceiptGetter implements receiptGetter and provides the standard behaviour for getting transaction
// receipts from the db
type transactionReceiptGetter [T crypto.PrivateKey, P crypto.PublicKey] struct {
	tx *Transaction[T,P]
}

func (g *transactionReceiptGetter[T,P]) get(ctx context.Context) (*types.Receipt[P], error) {
	if _, err := g.tx.resolve(ctx); err != nil {
		return nil, err
	}
	if g.tx.block == nil {
		return nil, nil
	}
	receipts, err := g.tx.block.resolveReceipts(ctx)
	if err != nil {
		return nil, err
	}
	return receipts[g.tx.index], nil
}

// (Quorum) privateTransactionReceiptGetter implements receiptGetter and gets privacy precompile transaction receipts
// from the the db
type privateTransactionReceiptGetter [T crypto.PrivateKey, P crypto.PublicKey] struct {
	pmt *Transaction[T,P]
}

func (g *privateTransactionReceiptGetter[T,P]) get(ctx context.Context) (*types.Receipt[P], error) {
	if _, err := g.pmt.resolve(ctx); err != nil {
		return nil, err
	}
	if g.pmt.block == nil {
		return nil, nil
	}
	receipts, err := g.pmt.block.resolveReceipts(ctx)
	if err != nil {
		return nil, err
	}
	receipt := receipts[g.pmt.index]

	psm, err := g.pmt.backend.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}

	privateReceipt := receipt.PSReceipts[psm.ID]
	if privateReceipt == nil {
		return nil, errors.New("could not find receipt for private transaction")
	}

	return privateReceipt, nil
}

// getReceipt returns the receipt associated with this transaction, if any.
func (t *Transaction[T,P]) getReceipt(ctx context.Context) (*types.Receipt[P], error) {
	// default to standard receipt getter if one is not set
	if t.receiptGetter == nil {
		t.receiptGetter = &transactionReceiptGetter[T,P]{tx: t}
	}
	return t.receiptGetter.get(ctx)
}

func (t *Transaction[T,P]) Status(ctx context.Context) (*Long, error) {
	receipt, err := t.getReceipt(ctx)
	if err != nil || receipt == nil {
		return nil, err
	}
	ret := Long(receipt.Status)
	return &ret, nil
}

func (t *Transaction[T,P]) GasUsed(ctx context.Context) (*Long, error) {
	receipt, err := t.getReceipt(ctx)
	if err != nil || receipt == nil {
		return nil, err
	}
	ret := Long(receipt.GasUsed)
	return &ret, nil
}

func (t *Transaction[T,P]) CumulativeGasUsed(ctx context.Context) (*Long, error) {
	receipt, err := t.getReceipt(ctx)
	if err != nil || receipt == nil {
		return nil, err
	}
	ret := Long(receipt.CumulativeGasUsed)
	return &ret, nil
}

func (t *Transaction[T,P]) CreatedContract(ctx context.Context, args BlockNumberArgs) (*Account[T,P], error) {
	receipt, err := t.getReceipt(ctx)
	if err != nil || receipt == nil || receipt.ContractAddress == (common.Address{}) {
		return nil, err
	}
	return &Account[T,P]{
		backend:       t.backend,
		address:       receipt.ContractAddress,
		blockNrOrHash: args.NumberOrLatest(),
	}, nil
}

func (t *Transaction[T,P]) Logs(ctx context.Context) (*[]*Log[T,P], error) {
	receipt, err := t.getReceipt(ctx)
	if err != nil || receipt == nil {
		return nil, err
	}
	ret := make([]*Log[T,P], 0, len(receipt.Logs))
	for _, log := range receipt.Logs {
		ret = append(ret, &Log[T,P]{
			backend:     t.backend,
			transaction: t,
			log:         log,
		})
	}
	return &ret, nil
}

func (t *Transaction[T,P]) Type(ctx context.Context) (*int32, error) {
	tx, err := t.resolve(ctx)
	if err != nil {
		return nil, err
	}
	txType := int32(tx.Type())
	return &txType, nil
}

func (t *Transaction[T,P]) AccessList(ctx context.Context) (*[]*AccessTuple, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return nil, err
	}
	accessList := tx.AccessList()
	ret := make([]*AccessTuple, 0, len(accessList))
	for _, al := range accessList {
		ret = append(ret, &AccessTuple{
			address:     al.Address,
			storageKeys: &al.StorageKeys,
		})
	}
	return &ret, nil
}

func (t *Transaction[T,P]) R(ctx context.Context) (hexutil.Big, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return hexutil.Big{}, err
	}
	_, r, _ := tx.RawSignatureValues()
	return hexutil.Big(*r), nil
}

func (t *Transaction[T,P]) S(ctx context.Context) (hexutil.Big, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return hexutil.Big{}, err
	}
	_, _, s := tx.RawSignatureValues()
	return hexutil.Big(*s), nil
}

func (t *Transaction[T,P]) V(ctx context.Context) (hexutil.Big, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return hexutil.Big{}, err
	}
	v, _, _ := tx.RawSignatureValues()
	return hexutil.Big(*v), nil
}

type BlockType int

// Block represents an Ethereum block.
// backend, and numberOrHash are mandatory. All other fields are lazily fetched
// when required.
type Block [T crypto.PrivateKey, P crypto.PublicKey] struct {
	backend      ethapi.Backend[T,P]
	numberOrHash *rpc.BlockNumberOrHash
	hash         common.Hash
	header       *types.Header[P]
	block        *types.Block[P]
	receipts     []*types.Receipt[P]
}

// resolve returns the internal Block object representing this block, fetching
// it if necessary.
func (b *Block[T,P]) resolve(ctx context.Context) (*types.Block[P], error) {
	if b.block != nil {
		return b.block, nil
	}
	if b.numberOrHash == nil {
		latest := rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
		b.numberOrHash = &latest
	}
	var err error
	b.block, err = b.backend.BlockByNumberOrHash(ctx, *b.numberOrHash)
	if b.block != nil && b.header == nil {
		b.header = b.block.Header()
		if hash, ok := b.numberOrHash.Hash(); ok {
			b.hash = hash
		}
	}
	return b.block, err
}

// resolveHeader returns the internal Header object for this block, fetching it
// if necessary. Call this function instead of `resolve` unless you need the
// additional data (transactions and uncles).
func (b *Block[T,P]) resolveHeader(ctx context.Context) (*types.Header[P], error) {
	if b.numberOrHash == nil && b.hash == (common.Hash{}) {
		return nil, errBlockInvariant
	}
	var err error
	if b.header == nil {
		if b.hash != (common.Hash{}) {
			b.header, err = b.backend.HeaderByHash(ctx, b.hash)
		} else {
			b.header, err = b.backend.HeaderByNumberOrHash(ctx, *b.numberOrHash)
		}
	}
	return b.header, err
}

// resolveReceipts returns the list of receipts for this block, fetching them
// if necessary.
func (b *Block[T,P]) resolveReceipts(ctx context.Context) ([]*types.Receipt[P], error) {
	if b.receipts == nil {
		hash := b.hash
		if hash == (common.Hash{}) {
			header, err := b.resolveHeader(ctx)
			if err != nil {
				return nil, err
			}
			hash = header.Hash()
		}
		receipts, err := b.backend.GetReceipts(ctx, hash)
		if err != nil {
			return nil, err
		}
		b.receipts = receipts
	}
	return b.receipts, nil
}

func (b *Block[T,P]) Number(ctx context.Context) (Long, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return 0, err
	}

	return Long(header.Number.Uint64()), nil
}

func (b *Block[T,P]) Hash(ctx context.Context) (common.Hash, error) {
	if b.hash == (common.Hash{}) {
		header, err := b.resolveHeader(ctx)
		if err != nil {
			return common.Hash{}, err
		}
		b.hash = header.Hash()
	}
	return b.hash, nil
}

func (b *Block[T,P]) GasLimit(ctx context.Context) (Long, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return 0, err
	}
	return Long(header.GasLimit), nil
}

func (b *Block[T,P]) GasUsed(ctx context.Context) (Long, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return 0, err
	}
	return Long(header.GasUsed), nil
}

func (b *Block[T,P]) Parent(ctx context.Context) (*Block[T,P], error) {
	// If the block header hasn't been fetched, and we'll need it, fetch it.
	if b.numberOrHash == nil && b.header == nil {
		if _, err := b.resolveHeader(ctx); err != nil {
			return nil, err
		}
	}
	if b.header != nil && b.header.Number.Uint64() > 0 {
		num := rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(b.header.Number.Uint64() - 1))
		return &Block[T,P]{
			backend:      b.backend,
			numberOrHash: &num,
			hash:         b.header.ParentHash,
		}, nil
	}
	return nil, nil
}

func (b *Block[T,P]) Difficulty(ctx context.Context) (hexutil.Big, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return hexutil.Big{}, err
	}
	return hexutil.Big(*header.Difficulty), nil
}

func (b *Block[T,P]) Timestamp(ctx context.Context) (hexutil.Uint64, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return 0, err
	}
	return hexutil.Uint64(header.Time), nil
}

func (b *Block[T,P]) Nonce(ctx context.Context) (hexutil.Bytes, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	return header.Nonce[:], nil
}

func (b *Block[T,P]) MixHash(ctx context.Context) (common.Hash, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return header.MixDigest, nil
}

func (b *Block[T,P]) TransactionsRoot(ctx context.Context) (common.Hash, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return header.TxHash, nil
}

func (b *Block[T,P]) StateRoot(ctx context.Context) (common.Hash, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return header.Root, nil
}

func (b *Block[T,P]) ReceiptsRoot(ctx context.Context) (common.Hash, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return header.ReceiptHash, nil
}

func (b *Block[T,P]) OmmerHash(ctx context.Context) (common.Hash, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	return header.UncleHash, nil
}

func (b *Block[T,P]) OmmerCount(ctx context.Context) (*int32, error) {
	block, err := b.resolve(ctx)
	if err != nil || block == nil {
		return nil, err
	}
	count := int32(len(block.Uncles()))
	return &count, err
}

func (b *Block[T,P]) Ommers(ctx context.Context) (*[]*Block[T,P], error) {
	block, err := b.resolve(ctx)
	if err != nil || block == nil {
		return nil, err
	}
	ret := make([]*Block[T,P], 0, len(block.Uncles()))
	for _, uncle := range block.Uncles() {
		blockNumberOrHash := rpc.BlockNumberOrHashWithHash(uncle.Hash(), false)
		ret = append(ret, &Block[T,P]{
			backend:      b.backend,
			numberOrHash: &blockNumberOrHash,
			header:       uncle,
		})
	}
	return &ret, nil
}

func (b *Block[T,P]) ExtraData(ctx context.Context) (hexutil.Bytes, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	return header.Extra, nil
}

func (b *Block[T,P]) LogsBloom(ctx context.Context) (hexutil.Bytes, error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return hexutil.Bytes{}, err
	}
	return header.Bloom.Bytes(), nil
}

func (b *Block[T,P]) TotalDifficulty(ctx context.Context) (hexutil.Big, error) {
	h := b.hash
	if h == (common.Hash{}) {
		header, err := b.resolveHeader(ctx)
		if err != nil {
			return hexutil.Big{}, err
		}
		h = header.Hash()
	}
	return hexutil.Big(*b.backend.GetTd(ctx, h)), nil
}

// BlockNumberArgs encapsulates arguments to accessors that specify a block number.
type BlockNumberArgs struct {
	// TODO: Ideally we could use input unions to allow the query to specify the
	// block parameter by hash, block number, or tag but input unions aren't part of the
	// standard GraphQL schema SDL yet, see: https://github.com/graphql/graphql-spec/issues/488
	Block *hexutil.Uint64
}

// NumberOr returns the provided block number argument, or the "current" block number or hash if none
// was provided.
func (a BlockNumberArgs) NumberOr(current rpc.BlockNumberOrHash) rpc.BlockNumberOrHash {
	if a.Block != nil {
		blockNr := rpc.BlockNumber(*a.Block)
		return rpc.BlockNumberOrHashWithNumber(blockNr)
	}
	return current
}

// NumberOrLatest returns the provided block number argument, or the "latest" block number if none
// was provided.
func (a BlockNumberArgs) NumberOrLatest() rpc.BlockNumberOrHash {
	return a.NumberOr(rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
}

func (b *Block[T,P]) Miner(ctx context.Context, args BlockNumberArgs) (*Account[T,P], error) {
	header, err := b.resolveHeader(ctx)
	if err != nil {
		return nil, err
	}
	return &Account[T,P]{
		backend:       b.backend,
		address:       header.Coinbase,
		blockNrOrHash: args.NumberOrLatest(),
	}, nil
}

func (b *Block[T,P]) TransactionCount(ctx context.Context) (*int32, error) {
	block, err := b.resolve(ctx)
	if err != nil || block == nil {
		return nil, err
	}
	count := int32(len(block.Transactions()))
	return &count, err
}

func (b *Block[T,P]) Transactions(ctx context.Context) (*[]*Transaction[T,P], error) {
	block, err := b.resolve(ctx)
	if err != nil || block == nil {
		return nil, err
	}
	ret := make([]*Transaction[T,P], 0, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		ret = append(ret, &Transaction[T,P]{
			backend: b.backend,
			hash:    tx.Hash(),
			tx:      tx,
			block:   b,
			index:   uint64(i),
		})
	}
	return &ret, nil
}

func (b *Block[T,P]) TransactionAt(ctx context.Context, args struct{ Index int32 }) (*Transaction[T,P], error) {
	block, err := b.resolve(ctx)
	if err != nil || block == nil {
		return nil, err
	}
	txs := block.Transactions()
	if args.Index < 0 || int(args.Index) >= len(txs) {
		return nil, nil
	}
	tx := txs[args.Index]
	return &Transaction[T,P]{
		backend: b.backend,
		hash:    tx.Hash(),
		tx:      tx,
		block:   b,
		index:   uint64(args.Index),
	}, nil
}

func (b *Block[T,P]) OmmerAt(ctx context.Context, args struct{ Index int32 }) (*Block[T,P], error) {
	block, err := b.resolve(ctx)
	if err != nil || block == nil {
		return nil, err
	}
	uncles := block.Uncles()
	if args.Index < 0 || int(args.Index) >= len(uncles) {
		return nil, nil
	}
	uncle := uncles[args.Index]
	blockNumberOrHash := rpc.BlockNumberOrHashWithHash(uncle.Hash(), false)
	return &Block[T,P]{
		backend:      b.backend,
		numberOrHash: &blockNumberOrHash,
		header:       uncle,
	}, nil
}

// BlockFilterCriteria encapsulates criteria passed to a `logs` accessor inside
// a block.
type BlockFilterCriteria struct {
	Addresses *[]common.Address // restricts matches to events created by specific contracts

	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Examples:
	// {} or nil          matches any topic list
	// {{A}}              matches topic A in first position
	// {{}, {B}}          matches any topic in first position, B in second position
	// {{A}, {B}}         matches topic A in first position, B in second position
	// {{A, B}}, {C, D}}  matches topic (A OR B) in first position, (C OR D) in second position
	Topics *[][]common.Hash
}

// runFilter accepts a filter and executes it, returning all its results as
// `Log` objects.
func runFilter[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, be ethapi.Backend[T,P], filter *filters.Filter[P]) ([]*Log[T,P], error) {
	logs, err := filter.Logs(ctx)
	if err != nil || logs == nil {
		return nil, err
	}
	ret := make([]*Log[T,P], 0, len(logs))
	for _, log := range logs {
		ret = append(ret, &Log[T,P]{
			backend:     be,
			transaction: &Transaction[T,P]{backend: be, hash: log.TxHash},
			log:         log,
		})
	}
	return ret, nil
}

func (b *Block[T,P]) Logs(ctx context.Context, args struct{ Filter BlockFilterCriteria }) ([]*Log[T,P], error) {
	var addresses []common.Address
	if args.Filter.Addresses != nil {
		addresses = *args.Filter.Addresses
	}
	var topics [][]common.Hash
	if args.Filter.Topics != nil {
		topics = *args.Filter.Topics
	}
	hash := b.hash
	if hash == (common.Hash{}) {
		header, err := b.resolveHeader(ctx)
		if err != nil {
			return nil, err
		}
		hash = header.Hash()
	}
	// Construct the range filter
	psm, err := b.backend.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}
	filter := filters.NewBlockFilter[P](b.backend, hash, addresses, topics, psm.ID)

	// Run the filter and return all the logs
	return runFilter(ctx, b.backend, filter)
}

func (b *Block[T,P]) Account(ctx context.Context, args struct {
	Address common.Address
}) (*Account[T,P], error) {
	if b.numberOrHash == nil {
		_, err := b.resolveHeader(ctx)
		if err != nil {
			return nil, err
		}
	}
	return &Account[T,P]{
		backend:       b.backend,
		address:       args.Address,
		blockNrOrHash: *b.numberOrHash,
	}, nil
}

// CallData encapsulates arguments to `call` or `estimateGas`.
// All arguments are optional.
type CallData struct {
	From     *common.Address // The Ethereum address the call is from.
	To       *common.Address // The Ethereum address the call is to.
	Gas      *hexutil.Uint64 // The amount of gas provided for the call.
	GasPrice *hexutil.Big    // The price of each unit of gas, in wei.
	Value    *hexutil.Big    // The value sent along with the call.
	Data     *hexutil.Bytes  // Any data sent with the call.
}

// CallResult encapsulates the result of an invocation of the `call` accessor.
type CallResult struct {
	data    hexutil.Bytes // The return data from the call
	gasUsed Long          // The amount of gas used
	status  Long          // The return status of the call - 0 for failure or 1 for success.
}

func (c *CallResult) Data() hexutil.Bytes {
	return c.data
}

func (c *CallResult) GasUsed() Long {
	return c.gasUsed
}

func (c *CallResult) Status() Long {
	return c.status
}

func (b *Block[T,P]) Call(ctx context.Context, args struct {
	Data ethapi.CallArgs
}) (*CallResult, error) {
	if b.numberOrHash == nil {
		_, err := b.resolve(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Quorum - replaced the default 5s time out with the value passed in vm.calltimeout
	result, err := ethapi.DoCall(ctx, b.backend, args.Data, *b.numberOrHash, nil, vm.Config[P]{}, b.backend.CallTimeOut(), b.backend.RPCGasCap())
	if err != nil {
		return nil, err
	}
	status := Long(1)
	if result.Failed() {
		status = 0
	}

	return &CallResult{
		data:    result.ReturnData,
		gasUsed: Long(result.UsedGas),
		status:  status,
	}, nil
}

func (b *Block[T,P]) EstimateGas(ctx context.Context, args struct {
	Data ethapi.CallArgs
}) (Long, error) {
	if b.numberOrHash == nil {
		_, err := b.resolveHeader(ctx)
		if err != nil {
			return 0, err
		}
	}
	gas, err := ethapi.DoEstimateGas(ctx, b.backend, args.Data, *b.numberOrHash, b.backend.RPCGasCap())
	return Long(gas), err
}

type Pending [T crypto.PrivateKey, P crypto.PublicKey] struct {
	backend ethapi.Backend[T,P]
}

func (p *Pending[T,P]) TransactionCount(ctx context.Context) (int32, error) {
	txs, err := p.backend.GetPoolTransactions()
	return int32(len(txs)), err
}

func (p *Pending[T,P]) Transactions(ctx context.Context) (*[]*Transaction[T,P], error) {
	txs, err := p.backend.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	ret := make([]*Transaction[T,P], 0, len(txs))
	for i, tx := range txs {
		ret = append(ret, &Transaction[T,P]{
			backend: p.backend,
			hash:    tx.Hash(),
			tx:      tx,
			index:   uint64(i),
		})
	}
	return &ret, nil
}

func (p *Pending[T,P]) Account(ctx context.Context, args struct {
	Address common.Address
}) *Account[T,P] {
	pendingBlockNr := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	return &Account[T,P]{
		backend:       p.backend,
		address:       args.Address,
		blockNrOrHash: pendingBlockNr,
	}
}

func (p *Pending[T,P]) Call(ctx context.Context, args struct {
	Data ethapi.CallArgs
}) (*CallResult, error) {
	pendingBlockNr := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)

	// Quorum - replaced the default 5s time out with the value passed in vm.calltimeout
	result, err := ethapi.DoCall(ctx, p.backend, args.Data, pendingBlockNr, nil, vm.Config[P]{}, p.backend.CallTimeOut(), p.backend.RPCGasCap())
	if err != nil {
		return nil, err
	}
	status := Long(1)
	if result.Failed() {
		status = 0
	}

	return &CallResult{
		data:    result.ReturnData,
		gasUsed: Long(result.UsedGas),
		status:  status,
	}, nil
}

func (p *Pending[T,P]) EstimateGas(ctx context.Context, args struct {
	Data ethapi.CallArgs
}) (Long, error) {
	pendingBlockNr := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	gas, err := ethapi.DoEstimateGas(ctx, p.backend, args.Data, pendingBlockNr, p.backend.RPCGasCap())
	return Long(gas), err
}

// Resolver is the top-level object in the GraphQL hierarchy.
type Resolver [T crypto.PrivateKey, P crypto.PublicKey] struct {
	backend ethapi.Backend[T,P]
}

func (r *Resolver[T,P]) Block(ctx context.Context, args struct {
	Number *Long
	Hash   *common.Hash
}) (*Block[T,P], error) {
	var block *Block[T,P]
	if args.Number != nil {
		if *args.Number < 0 {
			return nil, nil
		}
		number := rpc.BlockNumber(*args.Number)
		numberOrHash := rpc.BlockNumberOrHashWithNumber(number)
		block = &Block[T,P]{
			backend:      r.backend,
			numberOrHash: &numberOrHash,
		}
	} else if args.Hash != nil {
		numberOrHash := rpc.BlockNumberOrHashWithHash(*args.Hash, false)
		block = &Block[T,P]{
			backend:      r.backend,
			numberOrHash: &numberOrHash,
		}
	} else {
		numberOrHash := rpc.BlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
		block = &Block[T,P]{
			backend:      r.backend,
			numberOrHash: &numberOrHash,
		}
	}
	// Resolve the header, return nil if it doesn't exist.
	// Note we don't resolve block directly here since it will require an
	// additional network request for light client.
	h, err := block.resolveHeader(ctx)
	if err != nil {
		return nil, err
	} else if h == nil {
		return nil, nil
	}
	return block, nil
}

func (r *Resolver[T,P]) Blocks(ctx context.Context, args struct {
	From *Long
	To   *Long
}) ([]*Block[T,P], error) {
	from := rpc.BlockNumber(*args.From)

	var to rpc.BlockNumber
	if args.To != nil {
		to = rpc.BlockNumber(*args.To)
	} else {
		to = rpc.BlockNumber(r.backend.CurrentBlock().Number().Int64())
	}
	if to < from {
		return []*Block[T,P]{}, nil
	}
	ret := make([]*Block[T,P], 0, to-from+1)
	for i := from; i <= to; i++ {
		numberOrHash := rpc.BlockNumberOrHashWithNumber(i)
		ret = append(ret, &Block[T,P]{
			backend:      r.backend,
			numberOrHash: &numberOrHash,
		})
	}
	return ret, nil
}

func (r *Resolver[T,P]) Pending(ctx context.Context) *Pending[T,P] {
	return &Pending[T,P]{r.backend}
}

func (r *Resolver[T,P]) Transaction(ctx context.Context, args struct{ Hash common.Hash }) (*Transaction[T,P], error) {
	tx := &Transaction[T,P]{
		backend: r.backend,
		hash:    args.Hash,
	}
	// Resolve the transaction; if it doesn't exist, return nil.
	t, err := tx.resolve(ctx)
	if err != nil {
		return nil, err
	} else if t == nil {
		return nil, nil
	}
	return tx, nil
}

func (r *Resolver[T,P]) SendRawTransaction(ctx context.Context, args struct{ Data hexutil.Bytes }) (common.Hash, error) {
	tx := new(types.Transaction[P])
	if err := tx.UnmarshalBinary(args.Data); err != nil {
		return common.Hash{}, err
	}
	hash, err := ethapi.SubmitTransaction(ctx, r.backend, tx, "", true)
	return hash, err
}

// FilterCriteria encapsulates the arguments to `logs` on the root resolver object.
type FilterCriteria struct {
	FromBlock *hexutil.Uint64   // beginning of the queried range, nil means genesis block
	ToBlock   *hexutil.Uint64   // end of the range, nil means latest block
	Addresses *[]common.Address // restricts matches to events created by specific contracts

	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Examples:
	// {} or nil          matches any topic list
	// {{A}}              matches topic A in first position
	// {{}, {B}}          matches any topic in first position, B in second position
	// {{A}, {B}}         matches topic A in first position, B in second position
	// {{A, B}}, {C, D}}  matches topic (A OR B) in first position, (C OR D) in second position
	Topics *[][]common.Hash
}

func (r *Resolver[T,P]) Logs(ctx context.Context, args struct{ Filter FilterCriteria }) ([]*Log[T,P], error) {
	// Convert the RPC block numbers into internal representations
	begin := rpc.LatestBlockNumber.Int64()
	if args.Filter.FromBlock != nil {
		begin = int64(*args.Filter.FromBlock)
	}
	end := rpc.LatestBlockNumber.Int64()
	if args.Filter.ToBlock != nil {
		end = int64(*args.Filter.ToBlock)
	}
	var addresses []common.Address
	if args.Filter.Addresses != nil {
		addresses = *args.Filter.Addresses
	}
	var topics [][]common.Hash
	if args.Filter.Topics != nil {
		topics = *args.Filter.Topics
	}
	// Construct the range filter
	psm, err := r.backend.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}
	filter := filters.NewRangeFilter(filters.Backend[P](r.backend), begin, end, addresses, topics, psm.ID)
	return runFilter(ctx, r.backend, filter)
}

func (r *Resolver[T,P]) GasPrice(ctx context.Context) (hexutil.Big, error) {
	price, err := r.backend.SuggestPrice(ctx)
	return hexutil.Big(*price), err
}

func (r *Resolver[T,P]) ChainID(ctx context.Context) (hexutil.Big, error) {
	return hexutil.Big(*r.backend.ChainConfig().ChainID), nil
}

// SyncState represents the synchronisation status returned from the `syncing` accessor.
type SyncState struct {
	progress ethereum.SyncProgress
}

func (s *SyncState) StartingBlock() hexutil.Uint64 {
	return hexutil.Uint64(s.progress.StartingBlock)
}

func (s *SyncState) CurrentBlock() hexutil.Uint64 {
	return hexutil.Uint64(s.progress.CurrentBlock)
}

func (s *SyncState) HighestBlock() hexutil.Uint64 {
	return hexutil.Uint64(s.progress.HighestBlock)
}

func (s *SyncState) PulledStates() *hexutil.Uint64 {
	ret := hexutil.Uint64(s.progress.PulledStates)
	return &ret
}

func (s *SyncState) KnownStates() *hexutil.Uint64 {
	ret := hexutil.Uint64(s.progress.KnownStates)
	return &ret
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (r *Resolver[T,P]) Syncing() (*SyncState, error) {
	progress := r.backend.Downloader().Progress()

	// Return not syncing if the synchronisation already completed
	if progress.CurrentBlock >= progress.HighestBlock {
		return nil, nil
	}
	// Otherwise gather the block sync stats
	return &SyncState{progress}, nil
}

// Quorum

// PrivateTransaction returns the internal private transaction for privacy marker transactions
func (t *Transaction[T,P]) PrivateTransaction(ctx context.Context) (*Transaction[T,P], error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return nil, err
	}

	if !tx.IsPrivacyMarker() {
		// tx will not have a private tx so return early - no error to keep in line with other graphql behaviour (see PrivateInputData)
		return nil, nil
	}

	pvtTx, _, _, err := private.FetchPrivateTransaction[P](tx.Data())
	if err != nil {
		return nil, err
	}

	if pvtTx == nil {
		return nil, nil
	}

	return &Transaction[T,P]{
		backend:       t.backend,
		hash:          t.hash,
		tx:            pvtTx,
		block:         t.block,
		index:         t.index,
		receiptGetter: &privateTransactionReceiptGetter[T,P]{pmt: t},
	}, nil
}

func (t *Transaction[T,P]) IsPrivate(ctx context.Context) (*bool, error) {
	ret := false
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return &ret, err
	}
	ret = tx.IsPrivate()
	return &ret, nil
}

func (t *Transaction[T,P]) PrivateInputData(ctx context.Context) (*hexutil.Bytes, error) {
	tx, err := t.resolve(ctx)
	if err != nil || tx == nil {
		return &hexutil.Bytes{}, err
	}
	if tx.IsPrivate() {
		psm, err := t.backend.PSMR().ResolveForUserContext(ctx)
		if err != nil {
			return &hexutil.Bytes{}, err
		}
		_, managedParties, privateInputData, _, err := private.Ptm.Receive(common.BytesToEncryptedPayloadHash(tx.Data()))
		if err != nil || tx == nil {
			return &hexutil.Bytes{}, err
		}
		if t.backend.PSMR().NotIncludeAny(psm, managedParties...) {
			return &hexutil.Bytes{}, nil
		}
		ret := hexutil.Bytes(privateInputData)
		return &ret, nil
	}
	return &hexutil.Bytes{}, nil
}
