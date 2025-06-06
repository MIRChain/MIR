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

package core

import (
	"math/big"
	"testing"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"golang.org/x/crypto/sha3"
)

// TestStateProcessorErrors tests the output from the 'core' errors
// as defined in core/error.go. These errors are generated when the
// blockchain imports bad blocks, meaning blocks which have valid headers but
// contain invalid transactions
func TestStateProcessorErrors(t *testing.T) {
	var (
		signer     = types.HomesteadSigner[nist.PublicKey]{}
		testKey, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		db         = rawdb.NewMemoryDatabase()
		gspec      = &Genesis[nist.PublicKey]{
			Config: params.TestChainConfig,
		}
		genesis       = gspec.MustCommit(db)
		blockchain, _ = NewBlockChain[nist.PublicKey](db, nil, gspec.Config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	)
	defer blockchain.Stop()
	var makeTx = func(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction[nist.PublicKey] {
		tx, _ := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewTransaction[nist.PublicKey](nonce, to, amount, gasLimit, gasPrice, data), signer, testKey)
		return tx
	}
	for i, tt := range []struct {
		txs  []*types.Transaction[nist.PublicKey]
		want string
	}{
		{
			txs: []*types.Transaction[nist.PublicKey]{
				makeTx(0, common.Address{}, big.NewInt(0), params.TxGas, nil, nil),
				makeTx(0, common.Address{}, big.NewInt(0), params.TxGas, nil, nil),
			},
			want: "could not apply tx 1 [0x36bfa6d14f1cd35a1be8cc2322982a595fabc0e799f09c1de3bad7bd5b1f7626]: nonce too low: address 0x71562b71999873DB5b286dF957af199Ec94617F7, tx: 0 state: 1",
		},
		{
			txs: []*types.Transaction[nist.PublicKey]{
				makeTx(100, common.Address{}, big.NewInt(0), params.TxGas, nil, nil),
			},
			want: "could not apply tx 0 [0x51cd272d41ef6011d8138e18bf4043797aca9b713c7d39a97563f9bbe6bdbe6f]: nonce too high: address 0x71562b71999873DB5b286dF957af199Ec94617F7, tx: 100 state: 0",
		},
		{
			txs: []*types.Transaction[nist.PublicKey]{
				makeTx(0, common.Address{}, big.NewInt(0), 2100000000, nil, nil),
			},
			want: "could not apply tx 0 [0xa6111e2753b0495c90a4e5b709db7fadc4c3c7dc83ca5e80c2c8aedc53e6fa2c]: gas limit reached",
		},
		{
			txs: []*types.Transaction[nist.PublicKey]{
				makeTx(0, common.Address{}, big.NewInt(0), 21001, nil, nil),
			},
			want: "invalid gas used (remote: 0 local: 21000)", // "could not apply tx 0 [0x54c58b530824b0bb84b7a98183f08913b5d74e1cebc368515ef3c65edf8eb56a]: gas limit reached",
		},
		{
			txs: []*types.Transaction[nist.PublicKey]{
				makeTx(0, common.Address{}, big.NewInt(1), params.TxGas, nil, nil),
			},
			want: "could not apply tx 0 [0x3094b17498940d92b13baccf356ce8bfd6f221e926abc903d642fa1466c5b50e]: insufficient funds for transfer: address 0x71562b71999873DB5b286dF957af199Ec94617F7",
		},
		{
			txs: []*types.Transaction[nist.PublicKey]{
				makeTx(0, common.Address{}, big.NewInt(0), params.TxGas, big.NewInt(0xffffff), nil),
			},
			want: "could not apply tx 0 [0xaa3f7d86802b1f364576d9071bf231e31d61b392d306831ac9cf706ff5371ce0]: insufficient funds for gas * price + value: address 0x71562b71999873DB5b286dF957af199Ec94617F7 have 0 want 352321515000",
		},
		{
			txs: []*types.Transaction[nist.PublicKey]{
				makeTx(0, common.Address{}, big.NewInt(0), params.TxGas, nil, nil),
				makeTx(1, common.Address{}, big.NewInt(0), params.TxGas, nil, nil),
				makeTx(2, common.Address{}, big.NewInt(0), params.TxGas, nil, nil),
				makeTx(3, common.Address{}, big.NewInt(0), params.TxGas-1000, big.NewInt(0), nil),
			},
			want: "could not apply tx 3 [0x836fab5882205362680e49b311a20646de03b630920f18ec6ee3b111a2cf6835]: intrinsic gas too low: have 20000, want 21000",
		},
		// The last 'core' error is ErrGasUintOverflow: "gas uint64 overflow", but in order to
		// trigger that one, we'd have to allocate a _huge_ chunk of data, such that the
		// multiplication len(data) +gas_per_byte overflows uint64. Not testable at the moment
	} {
		block := GenerateBadBlock[nist.PublicKey](genesis,  ethash.NewFaker[nist.PublicKey](), tt.txs)
		_, err := blockchain.InsertChain(types.Blocks[nist.PublicKey]{block})
		if err == nil {
			t.Fatal("block imported without errors")
		}
		if have, want := err.Error(), tt.want; have != want {
			t.Errorf("test %d:\nhave \"%v\"\nwant \"%v\"\n", i, have, want)
		}
	}
}

// GenerateBadBlock constructs a "block" which contains the transactions. The transactions are not expected to be
// valid, and no proper post-state can be made. But from the perspective of the blockchain, the block is sufficiently
// valid to be considered for import:
// - valid pow (fake), ancestry, difficulty, gaslimit etc
func GenerateBadBlock[P crypto.PublicKey](parent *types.Block[P], engine consensus.Engine[P], txs types.Transactions[P]) *types.Block[P] {
	header := &types.Header[P]{
		ParentHash: parent.Hash(),
		Coinbase:   parent.Coinbase(),
		Difficulty: engine.CalcDifficulty(&fakeChainReader[P]{params.TestChainConfig}, parent.Time()+10, &types.Header[P]{
			Number:     parent.Number(),
			Time:       parent.Time(),
			Difficulty: parent.Difficulty(),
			UncleHash:  parent.UncleHash(),
		}),
		GasLimit:  CalcGasLimit(parent, parent.GasLimit(), parent.GasLimit(), parent.GasLimit()),
		Number:    new(big.Int).Add(parent.Number(), common.Big1),
		Time:      parent.Time() + 10,
		UncleHash: types.EmptyUncleHash[P](),
	}
	var receipts []*types.Receipt[P]

	// The post-state result doesn't need to be correct (this is a bad block), but we do need something there
	// Preferably something unique. So let's use a combo of blocknum + txhash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(header.Number.Bytes())
	var cumulativeGas uint64
	for _, tx := range txs {
		txh := tx.Hash()
		hasher.Write(txh[:])
		receipt := types.NewReceipt[P](nil, false, cumulativeGas+tx.Gas())
		receipt.TxHash = tx.Hash()
		receipt.GasUsed = tx.Gas()
		receipts = append(receipts, receipt)
		cumulativeGas += tx.Gas()
	}
	header.Root = common.BytesToHash(hasher.Sum(nil))
	// Assemble and return the final block for sealing
	return types.NewBlock[P](header, txs, nil, receipts, trie.NewStackTrie[P](nil))
}
