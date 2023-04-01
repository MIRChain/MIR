// Copyright 2015 The go-ethereum Authors
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

package filters

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

func makeReceipt[P crypto.PublicKey](addr common.Address) *types.Receipt[P] {
	receipt := types.NewReceipt[P](nil, false, 0)
	receipt.Logs = []*types.Log{
		{Address: addr},
	}
	receipt.Bloom = types.CreateBloom(types.Receipts[P]{receipt})
	return receipt
}

func BenchmarkFilters(b *testing.B) {
	dir, err := ioutil.TempDir("", "filtertest")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var (
		db, _   = rawdb.NewLevelDBDatabase(dir, 0, 0, "", false)
		backend = &testBackend[nist.PublicKey]{db: db}
		key1, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr1   = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())
		addr2   = common.BytesToAddress([]byte("jeff"))
		addr3   = common.BytesToAddress([]byte("ethereum"))
		addr4   = common.BytesToAddress([]byte("random addresses please"))
	)
	defer db.Close()

	genesis := core.GenesisBlockForTesting[nist.PublicKey](db, addr1, big.NewInt(1000000))
	chain, receipts := core.GenerateChain[nist.PublicKey](params.TestChainConfig, genesis,  ethash.NewFaker[nist.PublicKey](), db, 100010, func(i int, gen *core.BlockGen[nist.PublicKey]) {
		switch i {
		case 2403:
			receipt := makeReceipt[nist.PublicKey](addr1)
			gen.AddUncheckedReceipt(receipt)
		case 1034:
			receipt := makeReceipt[nist.PublicKey](addr2)
			gen.AddUncheckedReceipt(receipt)
		case 34:
			receipt := makeReceipt[nist.PublicKey](addr3)
			gen.AddUncheckedReceipt(receipt)
		case 99999:
			receipt := makeReceipt[nist.PublicKey](addr4)
			gen.AddUncheckedReceipt(receipt)

		}
	})
	for i, block := range chain {
		rawdb.WriteBlock(db, block)
		rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
		rawdb.WriteHeadBlockHash(db, block.Hash())
		rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), receipts[i])
	}
	b.ResetTimer()

	filter := NewRangeFilter[nist.PublicKey](backend, 0, -1, []common.Address{addr1, addr2, addr3, addr4}, nil, "")

	for i := 0; i < b.N; i++ {
		logs, _ := filter.Logs(context.Background())
		if len(logs) != 4 {
			b.Fatal("expected 4 logs, got", len(logs))
		}
	}
}

func TestFilters(t *testing.T) {
	dir, err := ioutil.TempDir("", "filtertest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var (
		db, _   = rawdb.NewLevelDBDatabase(dir, 0, 0, "", false)
		backend = &testBackend[nist.PublicKey]{db: db}
		key1, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr    = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())

		hash1 = common.BytesToHash([]byte("topic1"))
		hash2 = common.BytesToHash([]byte("topic2"))
		hash3 = common.BytesToHash([]byte("topic3"))
		hash4 = common.BytesToHash([]byte("topic4"))
		hash5 = common.BytesToHash([]byte("privateTopic"))
	)
	defer db.Close()

	genesis := core.GenesisBlockForTesting[nist.PublicKey](db, addr, big.NewInt(1000000))
	chain, receipts := core.GenerateChain[nist.PublicKey](params.TestChainConfig, genesis,  ethash.NewFaker[nist.PublicKey](), db, 1000, func(i int, gen *core.BlockGen[nist.PublicKey]) {
		switch i {
		case 1:
			receipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash1},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction[nist.PublicKey](1, common.HexToAddress("0x1"), big.NewInt(1), 1, big.NewInt(1), nil))
		case 2:
			receipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash2},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction[nist.PublicKey](2, common.HexToAddress("0x2"), big.NewInt(2), 2, big.NewInt(2), nil))

		case 998:
			receipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash3},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction[nist.PublicKey](998, common.HexToAddress("0x998"), big.NewInt(998), 998, big.NewInt(998), nil))
			// Add pseudo Quorum private transaction
			privateReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			privateReceipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash5},
				},
			}
			if err := rawdb.WritePrivateBlockBloom(db, 999, []*types.Receipt[nist.PublicKey]{privateReceipt}); err != nil {
				t.Fatal(err)
			}
			gen.AddUncheckedReceipt(privateReceipt)
			gen.AddUncheckedTx(types.NewTransaction[nist.PublicKey](998, common.HexToAddress("0x998"), big.NewInt(998), 998, big.NewInt(998), nil))
		case 999:
			receipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			receipt.Logs = []*types.Log{
				{
					Address: addr,
					Topics:  []common.Hash{hash4},
				},
			}
			gen.AddUncheckedReceipt(receipt)
			gen.AddUncheckedTx(types.NewTransaction[nist.PublicKey](999, common.HexToAddress("0x999"), big.NewInt(999), 999, big.NewInt(999), nil))
		}
	})
	for i, block := range chain {
		rawdb.WriteBlock(db, block)
		rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
		rawdb.WriteHeadBlockHash(db, block.Hash())
		rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), receipts[i])
	}

	filter := NewRangeFilter[nist.PublicKey](backend, 0, -1, []common.Address{addr}, [][]common.Hash{{hash1, hash2, hash3, hash4}}, "")

	logs, _ := filter.Logs(context.Background())
	if len(logs) != 4 {
		t.Error("expected 4 log, got", len(logs))
	}

	filter = NewRangeFilter[nist.PublicKey](backend, 900, 999, []common.Address{addr}, [][]common.Hash{{hash3}}, "")
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = NewRangeFilter[nist.PublicKey](backend, 990, -1, []common.Address{addr}, [][]common.Hash{{hash3}}, "")
	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash3 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash3, logs[0].Topics[0])
	}

	filter = NewRangeFilter[nist.PublicKey](backend, 1, 10, nil, [][]common.Hash{{hash1, hash2}}, "")

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 2 {
		t.Error("expected 2 log, got", len(logs))
	}

	failHash := common.BytesToHash([]byte("fail"))
	filter = NewRangeFilter[nist.PublicKey](backend, 0, -1, nil, [][]common.Hash{{failHash}}, "")

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}

	failAddr := common.BytesToAddress([]byte("failmenow"))
	filter = NewRangeFilter[nist.PublicKey](backend, 0, -1, []common.Address{failAddr}, nil, "")

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}

	filter = NewRangeFilter[nist.PublicKey](backend, 0, -1, nil, [][]common.Hash{{failHash}, {hash1}}, "")

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 0 {
		t.Error("expected 0 log, got", len(logs))
	}

	// Quorum

	// Test individual private log with NewBlockFilter (query filter with block hash)
	filter = NewBlockFilter[nist.PublicKey](backend, chain[998].Hash(), nil, [][]common.Hash{{hash5}}, "")

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 1 {
		t.Error("expected 1 log, got", len(logs))
	}
	if len(logs) > 0 && logs[0].Topics[0] != hash5 {
		t.Errorf("expected log[0].Topics[0] to be %x, got %x", hash5, logs[0].Topics[0])
	}

	// Test a mix of public and private logs with NewBlockFilter (query filter with block hash)
	filter = NewBlockFilter[nist.PublicKey](backend, chain[998].Hash(), nil, [][]common.Hash{{hash3, hash5}}, "")

	logs, _ = filter.Logs(context.Background())
	if len(logs) != 2 {
		t.Error("expected 2 log, got", len(logs))
	}

}

func TestMPSFilters(t *testing.T) {
	dir, err := ioutil.TempDir("", "filtermpstest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var (
		db, _   = rawdb.NewLevelDBDatabase(dir, 0, 0, "", false)
		backend = &testBackend[nist.PublicKey]{db: db}
		key1, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		addr    = crypto.PubkeyToAddress[nist.PublicKey](*key1.Public())
		hash1   = common.BytesToHash([]byte("topic1"))
	)
	defer db.Close()

	noPSILog := []*types.Log{
		{
			Address: addr,
			Topics:  []common.Hash{hash1},
		},
	}

	psi2PSILog := []*types.Log{
		{
			Address: addr,
			Topics:  []common.Hash{hash1},
			PSI:     types.PrivateStateIdentifier("psi2"),
		},
	}
	psi1PSILog := []*types.Log{
		{
			Address: addr,
			Topics:  []common.Hash{hash1},
			PSI:     types.PrivateStateIdentifier("psi1"),
		},
	}

	genesis := core.GenesisBlockForTesting[nist.PublicKey](db, addr, big.NewInt(1000000))
	chain, receipts := core.GenerateChain[nist.PublicKey](params.QuorumMPSTestChainConfig, genesis,  ethash.NewFaker[nist.PublicKey](), db, 1000, func(i int, gen *core.BlockGen[nist.PublicKey]) {
		switch i {
		case 1:
			//log on private transaction
			//has "private" psi receipt
			tx := types.NewTransaction[nist.PublicKey](1, common.HexToAddress("0x1"), big.NewInt(1), 1, big.NewInt(1), nil)
			tx.SetPrivate()
			privateReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			privateReceipt.Logs = noPSILog
			privateReceipt.PSReceipts = make(map[types.PrivateStateIdentifier]*types.Receipt[nist.PublicKey])
			psiReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			psiReceipt.Logs = psi2PSILog
			privateReceipt.PSReceipts[types.PrivateStateIdentifier("psi2")] = psiReceipt
			if err := rawdb.WritePrivateBlockBloom(db, 2, []*types.Receipt[nist.PublicKey]{privateReceipt}); err != nil {
				t.Fatal(err)
			}
			gen.AddUncheckedReceipt(privateReceipt)
			gen.AddUncheckedTx(tx)
		case 2:
			//no log on private transaction
			//has "psi1" receipt
			tx := types.NewTransaction[nist.PublicKey](2, common.HexToAddress("0x2"), big.NewInt(2), 2, big.NewInt(2), nil)
			tx.SetPrivate()
			privateReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			privateReceipt.PSReceipts = make(map[types.PrivateStateIdentifier]*types.Receipt[nist.PublicKey])
			psiReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			psiReceipt.Logs = psi1PSILog
			privateReceipt.PSReceipts[types.PrivateStateIdentifier("psi1")] = psiReceipt
			if err := rawdb.WritePrivateBlockBloom(db, 3, []*types.Receipt[nist.PublicKey]{privateReceipt}); err != nil {
				t.Fatal(err)
			}
			gen.AddUncheckedReceipt(privateReceipt)
			gen.AddUncheckedTx(tx)
		case 998:
			//no log on private transaction
			//has "psi2" psi receipt
			tx := types.NewTransaction[nist.PublicKey](998, common.HexToAddress("0x998"), big.NewInt(998), 998, big.NewInt(998), nil)
			tx.SetPrivate()
			privateReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			privateReceipt.PSReceipts = make(map[types.PrivateStateIdentifier]*types.Receipt[nist.PublicKey])
			psiReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			psiReceipt.Logs = psi2PSILog
			privateReceipt.PSReceipts[types.PrivateStateIdentifier("psi2")] = psiReceipt
			if err := rawdb.WritePrivateBlockBloom(db, 999, []*types.Receipt[nist.PublicKey]{privateReceipt}); err != nil {
				t.Fatal(err)
			}
			gen.AddUncheckedReceipt(privateReceipt)
			gen.AddUncheckedTx(tx)
		case 999:
			//log on private transaction
			//no psi receipt
			tx := types.NewTransaction[nist.PublicKey](999, common.HexToAddress("0x999"), big.NewInt(999), 999, big.NewInt(999), nil)
			tx.SetPrivate()
			privateReceipt := types.NewReceipt[nist.PublicKey](nil, false, 0)
			privateReceipt.Logs = noPSILog
			if err := rawdb.WritePrivateBlockBloom(db, 1000, []*types.Receipt[nist.PublicKey]{privateReceipt}); err != nil {
				t.Fatal(err)
			}
			gen.AddUncheckedReceipt(privateReceipt)
			gen.AddUncheckedTx(tx)
		}
	})
	for i, block := range chain {
		rawdb.WriteBlock(db, block)
		rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
		rawdb.WriteHeadBlockHash(db, block.Hash())
		rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), receipts[i])
	}

	//no psi filter: but only gets top level private receipt logs(no psi)
	filter := NewRangeFilter[nist.PublicKey](backend, 0, -1, []common.Address{addr}, [][]common.Hash{{hash1}}, "")
	logs, _ := filter.Logs(context.Background())
	if len(logs) != 2 {
		t.Error("expected 2 logs, got", len(logs))
	}

	//test filtering "psi2" logs: gets psi2 logs and top level private receipt logs(no psi)
	filter = NewRangeFilter[nist.PublicKey](backend, 0, -1, []common.Address{addr}, [][]common.Hash{{hash1}}, types.ToPrivateStateIdentifier("psi2"))
	ctx := rpc.WithPrivateStateIdentifier(context.Background(), types.ToPrivateStateIdentifier("psi2"))
	logs, _ = filter.Logs(ctx)
	if len(logs) != 3 {
		t.Error("expected 3 logs, got", len(logs))
	}

	//test filtering "psi1" logs: gets psi1 logs and top level private receipt logs(no psi)
	filter = NewRangeFilter[nist.PublicKey](backend, 0, -1, []common.Address{addr}, [][]common.Hash{{hash1}}, types.PrivateStateIdentifier("psi1"))
	ctx = rpc.WithPrivateStateIdentifier(context.Background(), types.ToPrivateStateIdentifier("psi1"))
	logs, _ = filter.Logs(ctx)
	if len(logs) != 3 {
		t.Error("expected 3 logs, got", len(logs))
	}
}
