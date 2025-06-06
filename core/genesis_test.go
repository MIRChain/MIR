// Copyright 2017 The go-ethereum Authors
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
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3411"

	//"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	//"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/params"
)

func TestDefaultGenesisBlock(t *testing.T) {
	block := DefaultGenesisBlock[gost3410.PublicKey]().ToBlock(nil)
	if block.Hash() != params.MainnetMirGenesisHash {
		t.Errorf("wrong mainnet genesis hash, got %v, want %v", block.Hash(), params.MainnetMirGenesisHash)
	}
	block = DefaultSoyuzGenesisBlock[gost3410.PublicKey]().ToBlock(nil)
	if block.Hash() != params.SoyuzGenesisHash {
		t.Errorf("wrong ropsten genesis hash, got %v, want %v", block.Hash(), params.SoyuzGenesisHash)
	}
}

func TestSetupGenesis(t *testing.T) {
	// Quorum: customized test cases for quorum
	var (
		customg = Genesis[gost3410.PublicKey]{
			Config: &params.ChainConfig{HomesteadBlock: big.NewInt(3), IsQuorum: true},
			Alloc: GenesisAlloc{
				{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
			},
		}
		oldcustomg = customg
	)
	oldcustomg.Config = &params.ChainConfig{HomesteadBlock: big.NewInt(2)}
	tests := []struct {
		name       string
		fn         func(ethdb.Database) (*params.ChainConfig, common.Hash, error)
		wantConfig *params.ChainConfig
		wantHash   common.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock[gost3410.PublicKey](db, new(Genesis[gost3410.PublicKey]))
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: params.AllEthashProtocolChanges,
		},
		{
			name: "no block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock[gost3410.PublicKey](db, nil)
			},
			wantHash:   params.MainnetEthGenesisHash,
			wantConfig: params.MainnetEthChainConfig,
		},
		{
			name: "mainnet block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock[gost3410.PublicKey]().MustCommit(db)
				return SetupGenesisBlock[gost3410.PublicKey](db, nil)
			},
			wantHash:   params.MainnetEthGenesisHash,
			wantConfig: params.MainnetEthChainConfig,
		},
		{
			name: "genesis with incorrect SizeLimit",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				customg.Config.TransactionSizeLimit = 100000
				customg.Config.MaxCodeSize = 32
				return SetupGenesisBlock[gost3410.PublicKey](db, &customg)
			},
			wantErr:    errors.New("Genesis transaction size limit must be between 32 and 128"),
			wantConfig: customg.Config,
		},
		{
			name: "genesis with incorrect max code size ",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				customg.Config.TransactionSizeLimit = 64
				customg.Config.MaxCodeSize = 100000
				return SetupGenesisBlock[gost3410.PublicKey](db, &customg)
			},
			wantErr:    errors.New("Genesis max code size must be between 24 and 128"),
			wantConfig: customg.Config,
		},
	}

	for _, test := range tests {
		db := rawdb.NewMemoryDatabase()
		config, hash, err := test.fn(db)
		// Check the return values.
		if !reflect.DeepEqual(err, test.wantErr) {
			spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
			t.Errorf("%s: returned error %#v, want %#v", test.name, spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
		}
		if !reflect.DeepEqual(config, test.wantConfig) {
			t.Errorf("%s:\nreturned %v\nwant     %v", test.name, config, test.wantConfig)
		}
		if hash != test.wantHash {
			t.Errorf("%s: returned hash %s, want %s", test.name, hash.Hex(), test.wantHash.Hex())
		} else if err == nil {
			// Check database content.
			stored := rawdb.ReadBlock[gost3410.PublicKey](db, test.wantHash, 0)
			if stored.Hash() != test.wantHash {
				t.Errorf("%s: block in DB has hash %s, want %s", test.name, stored.Hash(), test.wantHash)
			}
		}
	}
}

// TestGenesisHashes checks the congruity of default genesis data to corresponding hardcoded genesis hash values.
func TestGenesisHashes(t *testing.T) {
	h := gost3411.New256()
	h.Write([]byte("мирумир"))
	t.Log(hex.EncodeToString(h.Sum(nil)))
	h.Reset()
	h.Write([]byte("Союз1"))
	t.Log(hex.EncodeToString(h.Sum(nil)))

	cases := []struct {
		genesis *Genesis[gost3410.PublicKey]
		hash    common.Hash
	}{
		{
			genesis: DefaultGenesisBlock[gost3410.PublicKey](),
			hash:    params.MainnetMirGenesisHash,
		},
		{
			genesis: DefaultSoyuzGenesisBlock[gost3410.PublicKey](),
			hash:    params.SoyuzGenesisHash,
		},
	}
	for i, c := range cases {
		b := c.genesis.MustCommit(rawdb.NewMemoryDatabase())
		if got := b.Hash(); got != c.hash {
			t.Errorf("case: %d, want: %s, got: %s", i, c.hash.Hex(), got.Hex())
		}
	}
}
