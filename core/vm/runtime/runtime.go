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

package runtime

import (
	"math"
	"math/big"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/params"
)

// Config is a basic type specifying certain configuration flags for running
// the EVM.
type Config [P crypto.PublicKey] struct {
	ChainConfig *params.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   vm.Config[P]

	State     *state.StateDB[P]
	GetHashFn func(n uint64) common.Hash
}

// sets defaults on the config
func setDefaults[P crypto.PublicKey](cfg *Config[P]) {
	if cfg.ChainConfig == nil {
		cfg.ChainConfig = &params.ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      new(big.Int),
			DAOForkBlock:        new(big.Int),
			DAOForkSupport:      false,
			EIP150Block:         new(big.Int),
			EIP150Hash:          common.Hash{},
			EIP155Block:         new(big.Int),
			EIP158Block:         new(big.Int),
			ByzantiumBlock:      new(big.Int),
			ConstantinopleBlock: new(big.Int),
			PetersburgBlock:     new(big.Int),
			IstanbulBlock:       new(big.Int),
			MuirGlacierBlock:    new(big.Int),
			BerlinBlock:         new(big.Int),
			YoloV3Block:         nil,
		}
	}

	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = math.MaxUint64
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256[P]([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
}

// Execute executes the code using the input as call data during the execution.
// It returns the EVM's return value, the new state and an error if it failed.
//
// Execute sets up an in-memory, temporary, environment for the execution of
// the given code. It makes sure that it's restored to its original state afterwards.
func Execute[P crypto.PublicKey](code, input []byte, cfg *Config[P]) ([]byte, *state.StateDB[P], error) {
	if cfg == nil {
		cfg = new(Config[P])
	}
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New[P](common.Hash{}, state.NewDatabase[P](rawdb.NewMemoryDatabase()), nil)
	}
	var (
		address = common.BytesToAddress([]byte("contract"))
		vmenv   = NewEnv[P](cfg)
		sender  = vm.AccountRef(cfg.Origin)
	)
	if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber); rules.IsBerlin {
		cfg.State.PrepareAccessList(cfg.Origin, &address, vm.ActivePrecompiles(rules), nil)
	}
	cfg.State.CreateAccount(address)
	// set the receiver's (the executing contract) code for execution.
	cfg.State.SetCode(address, code)
	// Call the code with the given configuration.
	ret, _, err := vmenv.Call(
		sender,
		common.BytesToAddress([]byte("contract")),
		input,
		cfg.GasLimit,
		cfg.Value,
	)

	return ret, cfg.State, err
}

// Create executes the code using the EVM create method
func Create[P crypto.PublicKey](input []byte, cfg *Config[P]) ([]byte, common.Address, uint64, error) {
	if cfg == nil {
		cfg = new(Config[P])
	}
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New[P](common.Hash{}, state.NewDatabase[P](rawdb.NewMemoryDatabase()), nil)
	}
	var (
		vmenv  = NewEnv[P](cfg)
		sender = vm.AccountRef(cfg.Origin)
	)
	if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber); rules.IsBerlin {
		cfg.State.PrepareAccessList(cfg.Origin, nil, vm.ActivePrecompiles(rules), nil)
	}
	// Call the code with the given configuration.
	code, address, leftOverGas, err := vmenv.Create(
		sender,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	return code, address, leftOverGas, err
}

// Call executes the code given by the contract's address. It will return the
// EVM's return value or an error if it failed.
//
// Call, unlike Execute, requires a config and also requires the State field to
// be set.
func Call[P crypto.PublicKey](address common.Address, input []byte, cfg *Config[P]) ([]byte, uint64, error) {
	setDefaults(cfg)

	vmenv := NewEnv[P](cfg)

	sender := cfg.State.GetOrNewStateObject(cfg.Origin)
	statedb := cfg.State

	if rules := cfg.ChainConfig.Rules(vmenv.Context.BlockNumber); rules.IsBerlin {
		statedb.PrepareAccessList(cfg.Origin, &address, vm.ActivePrecompiles(rules), nil)
	}
	// Call the code with the given configuration.
	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	return ret, leftOverGas, err
}
