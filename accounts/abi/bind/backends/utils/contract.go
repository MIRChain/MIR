// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/MIRChain/MIR"
	"github.com/MIRChain/MIR/accounts/abi"
	"github.com/MIRChain/MIR/accounts/abi/bind"
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = common.Big1
	_ = event.NewSubscription
)

// GasEstimationABI is the input ABI used to generate the binding from.
const GasEstimationABI = "[{\"inputs\":[],\"name\":\"Assert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OOG\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PureRevert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"Revert\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"Valid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var GasEstimationParsedABI, _ = abi.JSON[gost3410.PublicKey](strings.NewReader(GasEstimationABI))

// GasEstimationFuncSigs maps the 4-byte function signature to its string representation.
var GasEstimationFuncSigs = map[string]string{
	"b9b046f9": "Assert()",
	"50f6fe34": "OOG()",
	"aa8b1d30": "PureRevert()",
	"d8b98391": "Revert()",
	"e09fface": "Valid()",
}

// GasEstimationBin is the compiled bytecode used for deploying new contracts.
var GasEstimationBin = "0x608060405234801561001057600080fd5b5061014c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806350f6fe341461005c578063aa8b1d3014610066578063b9b046f91461006e578063d8b9839114610076578063e09fface14610064575b600080fd5b61006461007e565b005b610064600080fd5b610064610093565b61006461009d565b60005b8061008b816100d9565b915050610081565b61009b610100565b565b60405162461bcd60e51b815260206004820152600d60248201526c3932bb32b93a103932b0b9b7b760991b604482015260640160405180910390fd5b6000600182016100f957634e487b7160e01b600052601160045260246000fd5b5060010190565b634e487b7160e01b600052600160045260246000fdfea26469706673582212208dcba1fae05ef5cad0febf6f88e275d8b0e913d4906b5bd9606ba63c6936abb164736f6c63430008120033"

// DeployGasEstimation deploys a new Ethereum contract, binding an instance of GasEstimation to it.
func DeployGasEstimation[P crypto.PublicKey](auth *bind.TransactOpts[P], backend bind.ContractBackend[P]) (common.Address, *types.Transaction[P], *GasEstimation[P], error) {
	parsed, err := abi.JSON[P](strings.NewReader(GasEstimationABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(GasEstimationBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &GasEstimation[P]{GasEstimationCaller: GasEstimationCaller[P]{contract: contract}, GasEstimationTransactor: GasEstimationTransactor[P]{contract: contract}, GasEstimationFilterer: GasEstimationFilterer[P]{contract: contract}}, nil
}

// GasEstimation is an auto generated Go binding around an Ethereum contract.
type GasEstimation[P crypto.PublicKey] struct {
	GasEstimationCaller[P]     // Read-only binding to the contract
	GasEstimationTransactor[P] // Write-only binding to the contract
	GasEstimationFilterer[P]   // Log filterer for contract events
}

// GasEstimationCaller is an auto generated read-only Go binding around an Ethereum contract.
type GasEstimationCaller[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// GasEstimationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GasEstimationTransactor[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// GasEstimationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GasEstimationFilterer[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// GasEstimationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GasEstimationSession[P crypto.PublicKey] struct {
	Contract     *GasEstimation[P]    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts[P] // Transaction auth options to use throughout this session
}

// GasEstimationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GasEstimationCallerSession[P crypto.PublicKey] struct {
	Contract *GasEstimationCaller[P] // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// GasEstimationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GasEstimationTransactorSession[P crypto.PublicKey] struct {
	Contract     *GasEstimationTransactor[P] // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts[P]        // Transaction auth options to use throughout this session
}

// GasEstimationRaw is an auto generated low-level Go binding around an Ethereum contract.
type GasEstimationRaw[P crypto.PublicKey] struct {
	Contract *GasEstimation[P] // Generic contract binding to access the raw methods on
}

// GasEstimationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GasEstimationCallerRaw[P crypto.PublicKey] struct {
	Contract *GasEstimationCaller[P] // Generic read-only contract binding to access the raw methods on
}

// GasEstimationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GasEstimationTransactorRaw[P crypto.PublicKey] struct {
	Contract *GasEstimationTransactor[P] // Generic write-only contract binding to access the raw methods on
}

// NewGasEstimation creates a new instance of GasEstimation, bound to a specific deployed contract.
func NewGasEstimation[P crypto.PublicKey](address common.Address, backend bind.ContractBackend[P]) (*GasEstimation[P], error) {
	contract, err := bindGasEstimation[P](address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GasEstimation[P]{GasEstimationCaller: GasEstimationCaller[P]{contract: contract}, GasEstimationTransactor: GasEstimationTransactor[P]{contract: contract}, GasEstimationFilterer: GasEstimationFilterer[P]{contract: contract}}, nil
}

// NewGasEstimationCaller creates a new read-only instance of GasEstimation, bound to a specific deployed contract.
func NewGasEstimationCaller[P crypto.PublicKey](address common.Address, caller bind.ContractCaller) (*GasEstimationCaller[P], error) {
	contract, err := bindGasEstimation[P](address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GasEstimationCaller[P]{contract: contract}, nil
}

// NewGasEstimationTransactor creates a new write-only instance of GasEstimation, bound to a specific deployed contract.
func NewGasEstimationTransactor[P crypto.PublicKey](address common.Address, transactor bind.ContractTransactor[P]) (*GasEstimationTransactor[P], error) {
	contract, err := bindGasEstimation[P](address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GasEstimationTransactor[P]{contract: contract}, nil
}

// NewGasEstimationFilterer creates a new log filterer instance of GasEstimation, bound to a specific deployed contract.
func NewGasEstimationFilterer[P crypto.PublicKey](address common.Address, filterer bind.ContractFilterer) (*GasEstimationFilterer[P], error) {
	contract, err := bindGasEstimation[P](address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GasEstimationFilterer[P]{contract: contract}, nil
}

// bindGasEstimation binds a generic wrapper to an already deployed contract.
func bindGasEstimation[P crypto.PublicKey](address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor[P], filterer bind.ContractFilterer) (*bind.BoundContract[P], error) {
	parsed, err := abi.JSON[P](strings.NewReader(GasEstimationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasEstimation *GasEstimationRaw[P]) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasEstimation.Contract.GasEstimationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasEstimation *GasEstimationRaw[P]) Transfer(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _GasEstimation.Contract.GasEstimationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasEstimation *GasEstimationRaw[P]) Transact(opts *bind.TransactOpts[P], method string, params ...interface{}) (*types.Transaction[P], error) {
	return _GasEstimation.Contract.GasEstimationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GasEstimation *GasEstimationCallerRaw[P]) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GasEstimation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GasEstimation *GasEstimationTransactorRaw[P]) Transfer(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _GasEstimation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GasEstimation *GasEstimationTransactorRaw[P]) Transact(opts *bind.TransactOpts[P], method string, params ...interface{}) (*types.Transaction[P], error) {
	return _GasEstimation.Contract.contract.Transact(opts, method, params...)
}

// Assert is a paid mutator transaction binding the contract method 0x4d2a6de5.
//
// Solidity: function Assert() returns()
func (_GasEstimation *GasEstimationTransactor[P]) Assert(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _GasEstimation.contract.Transact(opts, "Assert")
}

// Assert is a paid mutator transaction binding the contract method 0x4d2a6de5.
//
// Solidity: function Assert() returns()
func (_GasEstimation *GasEstimationSession[P]) Assert() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.Assert(&_GasEstimation.TransactOpts)
}

// Assert is a paid mutator transaction binding the contract method 0x4d2a6de5.
//
// Solidity: function Assert() returns()
func (_GasEstimation *GasEstimationTransactorSession[P]) Assert() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.Assert(&_GasEstimation.TransactOpts)
}

// OOG is a paid mutator transaction binding the contract method 0x01433d3b.
//
// Solidity: function OOG() returns()
func (_GasEstimation *GasEstimationTransactor[P]) OOG(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _GasEstimation.contract.Transact(opts, "OOG")
}

// OOG is a paid mutator transaction binding the contract method 0x01433d3b.
//
// Solidity: function OOG() returns()
func (_GasEstimation *GasEstimationSession[P]) OOG() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.OOG(&_GasEstimation.TransactOpts)
}

// OOG is a paid mutator transaction binding the contract method 0x01433d3b.
//
// Solidity: function OOG() returns()
func (_GasEstimation *GasEstimationTransactorSession[P]) OOG() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.OOG(&_GasEstimation.TransactOpts)
}

// PureRevert is a paid mutator transaction binding the contract method 0xeb247ceb.
//
// Solidity: function PureRevert() returns()
func (_GasEstimation *GasEstimationTransactor[P]) PureRevert(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _GasEstimation.contract.Transact(opts, "PureRevert")
}

// PureRevert is a paid mutator transaction binding the contract method 0xeb247ceb.
//
// Solidity: function PureRevert() returns()
func (_GasEstimation *GasEstimationSession[P]) PureRevert() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.PureRevert(&_GasEstimation.TransactOpts)
}

// PureRevert is a paid mutator transaction binding the contract method 0xeb247ceb.
//
// Solidity: function PureRevert() returns()
func (_GasEstimation *GasEstimationTransactorSession[P]) PureRevert() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.PureRevert(&_GasEstimation.TransactOpts)
}

// Revert is a paid mutator transaction binding the contract method 0xc4bef2db.
//
// Solidity: function Revert() returns()
func (_GasEstimation *GasEstimationTransactor[P]) Revert(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _GasEstimation.contract.Transact(opts, "Revert")
}

// Revert is a paid mutator transaction binding the contract method 0xc4bef2db.
//
// Solidity: function Revert() returns()
func (_GasEstimation *GasEstimationSession[P]) Revert() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.Revert(&_GasEstimation.TransactOpts)
}

// Revert is a paid mutator transaction binding the contract method 0xc4bef2db.
//
// Solidity: function Revert() returns()
func (_GasEstimation *GasEstimationTransactorSession[P]) Revert() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.Revert(&_GasEstimation.TransactOpts)
}

// Valid is a paid mutator transaction binding the contract method 0x4a694494.
//
// Solidity: function Valid() returns()
func (_GasEstimation *GasEstimationTransactor[P]) Valid(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _GasEstimation.contract.Transact(opts, "Valid")
}

// Valid is a paid mutator transaction binding the contract method 0x4a694494.
//
// Solidity: function Valid() returns()
func (_GasEstimation *GasEstimationSession[P]) Valid() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.Valid(&_GasEstimation.TransactOpts)
}

// Valid is a paid mutator transaction binding the contract method 0x4a694494.
//
// Solidity: function Valid() returns()
func (_GasEstimation *GasEstimationTransactorSession[P]) Valid() (*types.Transaction[P], error) {
	return _GasEstimation.Contract.Valid(&_GasEstimation.TransactOpts)
}
