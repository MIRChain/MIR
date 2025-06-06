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

package tests

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/state/snapshot"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"golang.org/x/crypto/sha3"
)

// StateTest checks transaction processing without block context.
// See https://github.com/ethereum/EIPs/issues/176 for the test format specification.
type StateTest [T crypto.PrivateKey, P crypto.PublicKey] struct {
	json stJSON[T,P]
}

// StateSubtest selects a specific configuration of a General State Test.
type StateSubtest struct {
	Fork  string
	Index int
}

func (t *StateTest[T,P]) UnmarshalJSON(in []byte) error {
	return json.Unmarshal(in, &t.json)
}

type stJSON [T crypto.PrivateKey, P crypto.PublicKey] struct {
	Env  stEnv                    `json:"env"`
	Pre  core.GenesisAlloc        `json:"pre"`
	Tx   stTransaction[T,P]            `json:"transaction"`
	Out  hexutil.Bytes            `json:"out"`
	Post map[string][]stPostState `json:"post"`
}

type stPostState struct {
	Root    common.UnprefixedHash `json:"hash"`
	Logs    common.UnprefixedHash `json:"logs"`
	Indexes struct {
		Data  int `json:"data"`
		Gas   int `json:"gas"`
		Value int `json:"value"`
	}
}

//go:generate gencodec -type stEnv -field-override stEnvMarshaling -out gen_stenv.go

type stEnv struct {
	Coinbase   common.Address `json:"currentCoinbase"   gencodec:"required"`
	Difficulty *big.Int       `json:"currentDifficulty" gencodec:"required"`
	GasLimit   uint64         `json:"currentGasLimit"   gencodec:"required"`
	Number     uint64         `json:"currentNumber"     gencodec:"required"`
	Timestamp  uint64         `json:"currentTimestamp"  gencodec:"required"`
}

type stEnvMarshaling struct {
	Coinbase   common.UnprefixedAddress
	Difficulty *math.HexOrDecimal256
	GasLimit   math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
}

//go:generate gencodec -type stTransaction -field-override stTransactionMarshaling -out gen_sttransaction.go

type stTransaction [T crypto.PrivateKey, P crypto.PublicKey]struct {
	GasPrice    *big.Int            `json:"gasPrice"`
	Nonce       uint64              `json:"nonce"`
	To          string              `json:"to"`
	Data        []string            `json:"data"`
	AccessLists []*types.AccessList `json:"accessLists,omitempty"`
	GasLimit    []uint64            `json:"gasLimit"`
	Value       []string            `json:"value"`
	PrivateKey  []byte              `json:"secretKey"`
}

type stTransactionMarshaling struct {
	GasPrice   *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	GasLimit   []math.HexOrDecimal64
	PrivateKey hexutil.Bytes
}

// GetChainConfig takes a fork definition and returns a chain config.
// The fork definition can be
// - a plain forkname, e.g. `Byzantium`,
// - a fork basename, and a list of EIPs to enable; e.g. `Byzantium+1884+1283`.
func GetChainConfig[P crypto.PublicKey](forkString string) (baseConfig *params.ChainConfig, eips []int, err error) {
	var (
		splitForks            = strings.Split(forkString, "+")
		ok                    bool
		baseName, eipsStrings = splitForks[0], splitForks[1:]
	)
	if baseConfig, ok = Forks[baseName]; !ok {
		return nil, nil, UnsupportedForkError{baseName}
	}
	for _, eip := range eipsStrings {
		if eipNum, err := strconv.Atoi(eip); err != nil {
			return nil, nil, fmt.Errorf("syntax error, invalid eip number %v", eipNum)
		} else {
			if !vm.ValidEip[P](eipNum) {
				return nil, nil, fmt.Errorf("syntax error, invalid eip number %v", eipNum)
			}
			eips = append(eips, eipNum)
		}
	}
	return baseConfig, eips, nil
}

// Subtests returns all valid subtests of the test.
func (t *StateTest[T,P]) Subtests() []StateSubtest {
	var sub []StateSubtest
	for fork, pss := range t.json.Post {
		for i := range pss {
			sub = append(sub, StateSubtest{fork, i})
		}
	}
	return sub
}

// Run executes a specific subtest and verifies the post-state and logs
func (t *StateTest[T,P]) Run(subtest StateSubtest, vmconfig vm.Config[P], snapshotter bool) (*snapshot.Tree[P], *state.StateDB[P], error) {
	snaps, statedb, root, err := t.RunNoVerify(subtest, vmconfig, snapshotter)
	if err != nil {
		return snaps, statedb, err
	}
	post := t.json.Post[subtest.Fork][subtest.Index]
	// N.B: We need to do this in a two-step process, because the first Commit takes care
	// of suicides, and we need to touch the coinbase _after_ it has potentially suicided.
	if root != common.Hash(post.Root) {
		return snaps, statedb, fmt.Errorf("post state root mismatch: got %x, want %x", root, post.Root)
	}
	if logs := rlpHash(statedb.Logs()); logs != common.Hash(post.Logs) {
		return snaps, statedb, fmt.Errorf("post state logs hash mismatch: got %x, want %x", logs, post.Logs)
	}
	return snaps, statedb, nil
}

// RunNoVerify runs a specific subtest and returns the statedb and post-state root
func (t *StateTest[T,P]) RunNoVerify(subtest StateSubtest, vmconfig vm.Config[P], snapshotter bool) (*snapshot.Tree[P], *state.StateDB[P], common.Hash, error) {
	config, eips, err := GetChainConfig[P](subtest.Fork)
	if err != nil {
		return nil, nil, common.Hash{}, UnsupportedForkError{subtest.Fork}
	}
	vmconfig.ExtraEips = eips
	block := t.genesis(config).ToBlock(nil)
	snaps, statedb := MakePreState[P](rawdb.NewMemoryDatabase(), t.json.Pre, snapshotter)

	post := t.json.Post[subtest.Fork][subtest.Index]
	msg, err := t.json.Tx.toMessage(post)
	if err != nil {
		return nil, nil, common.Hash{}, err
	}

	// Prepare the EVM.
	txContext := core.NewEVMTxContext(msg)
	context := core.NewEVMBlockContext[P](block.Header(), nil, &t.json.Env.Coinbase)
	context.GetHash = vmTestBlockHash[P]
	evm := vm.NewEVM(context, txContext, statedb, statedb, config, vmconfig)

	// Execute the message.
	snapshot := statedb.Snapshot()
	gaspool := new(core.GasPool)
	gaspool.AddGas(block.GasLimit())
	if _, err := core.ApplyMessage(evm, msg, gaspool); err != nil {
		statedb.RevertToSnapshot(snapshot)
	}

	// Commit block
	statedb.Commit(config.IsEIP158(block.Number()))
	// Add 0-value mining reward. This only makes a difference in the cases
	// where
	// - the coinbase suicided, or
	// - there are only 'bad' transactions, which aren't executed. In those cases,
	//   the coinbase gets no txfee, so isn't created, and thus needs to be touched
	statedb.AddBalance(block.Coinbase(), new(big.Int))
	// And _now_ get the state root
	root := statedb.IntermediateRoot(config.IsEIP158(block.Number()))
	return snaps, statedb, root, nil
}

func (t *StateTest[T,P]) gasLimit(subtest StateSubtest) uint64 {
	return t.json.Tx.GasLimit[t.json.Post[subtest.Fork][subtest.Index].Indexes.Gas]
}

func MakePreState[P crypto.PublicKey](db ethdb.Database, accounts core.GenesisAlloc, snapshotter bool) (*snapshot.Tree[P], *state.StateDB[P]) {
	sdb := state.NewDatabase[P](db)
	statedb, _ := state.New[P](common.Hash{}, sdb, nil)
	for addr, a := range accounts {
		statedb.SetCode(addr, a.Code)
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
		for k, v := range a.Storage {
			statedb.SetState(addr, k, v)
		}
	}
	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false)

	var snaps *snapshot.Tree[P]
	if snapshotter {
		snaps, _ = snapshot.New[P](db, sdb.TrieDB(), 1, root, false, true, false)
	}
	statedb, _ = state.New(root, sdb, snaps)
	return snaps, statedb
}

func (t *StateTest[T,P]) genesis(config *params.ChainConfig) *core.Genesis[P] {
	return &core.Genesis[P]{
		Config:     config,
		Coinbase:   t.json.Env.Coinbase,
		Difficulty: t.json.Env.Difficulty,
		GasLimit:   t.json.Env.GasLimit,
		Number:     t.json.Env.Number,
		Timestamp:  t.json.Env.Timestamp,
		Alloc:      t.json.Pre,
	}
}

func (tx *stTransaction[T,P]) toMessage(ps stPostState) (core.Message, error) {
	// Derive sender from private key if present.
	var from common.Address
	if len(tx.PrivateKey) > 0 {
		key, err := crypto.ToECDSA[T](tx.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %v", err)
		}
		var pub P
		switch t:=any(&key).(type){
		case *nist.PrivateKey:
			p:=any(&pub).(*nist.PublicKey)
			*p = *t.Public()
		}
		from = crypto.PubkeyToAddress[P](pub)
	}
	// Parse recipient if present.
	var to *common.Address
	if tx.To != "" {
		to = new(common.Address)
		if err := to.UnmarshalText([]byte(tx.To)); err != nil {
			return nil, fmt.Errorf("invalid to address: %v", err)
		}
	}

	// Get values specific to this post state.
	if ps.Indexes.Data > len(tx.Data) {
		return nil, fmt.Errorf("tx data index %d out of bounds", ps.Indexes.Data)
	}
	if ps.Indexes.Value > len(tx.Value) {
		return nil, fmt.Errorf("tx value index %d out of bounds", ps.Indexes.Value)
	}
	if ps.Indexes.Gas > len(tx.GasLimit) {
		return nil, fmt.Errorf("tx gas limit index %d out of bounds", ps.Indexes.Gas)
	}
	dataHex := tx.Data[ps.Indexes.Data]
	valueHex := tx.Value[ps.Indexes.Value]
	gasLimit := tx.GasLimit[ps.Indexes.Gas]
	// Value, Data hex encoding is messy: https://github.com/ethereum/tests/issues/203
	value := new(big.Int)
	if valueHex != "0x" {
		v, ok := math.ParseBig256(valueHex)
		if !ok {
			return nil, fmt.Errorf("invalid tx value %q", valueHex)
		}
		value = v
	}
	data, err := hex.DecodeString(strings.TrimPrefix(dataHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid tx data %q", dataHex)
	}
	var accessList types.AccessList
	if tx.AccessLists != nil && tx.AccessLists[ps.Indexes.Data] != nil {
		accessList = *tx.AccessLists[ps.Indexes.Data]
	}
	msg := types.NewMessage(from, to, tx.Nonce, value, gasLimit, tx.GasPrice, data, accessList, true)
	return msg, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
