// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package simple

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

// SimpleABI is the input ABI used to generate the binding from.
const SimpleABI = "[{\"inputs\":[],\"name\":\"getValue\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"v\",\"type\":\"bytes16\"}],\"name\":\"setValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

var SimpleParsedABI, _ = abi.JSON[gost3410.PublicKey](strings.NewReader(SimpleABI))

// SimpleBin is the compiled bytecode used for deploying new contracts.
var SimpleBin = "0x608060405234801561001057600080fd5b50610206806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806307ea02e1146100465780637d7ee50e14610064578063c7ff210314610082575b600080fd5b61004e61009e565b60405161005b919061012c565b60405180910390f35b61006c6100b4565b604051610079919061012c565b60405180910390f35b61009c60048036038101906100979190610178565b6100c5565b005b60008060009054906101000a900460801b905090565b60008054906101000a900460801b81565b806000806101000a8154816fffffffffffffffffffffffffffffffff021916908360801c021790555050565b60007fffffffffffffffffffffffffffffffff0000000000000000000000000000000082169050919050565b610126816100f1565b82525050565b6000602082019050610141600083018461011d565b92915050565b600080fd5b610155816100f1565b811461016057600080fd5b50565b6000813590506101728161014c565b92915050565b60006020828403121561018e5761018d610147565b5b600061019c84828501610163565b9150509291505056fea2646970667358221220a94d387a337f7be406d713b9d25400b18463d3def57c02777bfb3af704806ec864736f6c63782d302e382e31342d646576656c6f702e323032342e31312e32362b636f6d6d69742e38306434396633372e6d6f64005e"

// DeploySimple deploys a new Ethereum contract, binding an instance of Simple to it.
func DeploySimple[P crypto.PublicKey](auth *bind.TransactOpts[P], backend bind.ContractBackend[P]) (common.Address, *types.Transaction[P], *Simple[P], error) {
	parsed, err := abi.JSON[P](strings.NewReader(SimpleABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SimpleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Simple[P]{SimpleCaller: SimpleCaller[P]{contract: contract}, SimpleTransactor: SimpleTransactor[P]{contract: contract}, SimpleFilterer: SimpleFilterer[P]{contract: contract}}, nil
}

// Simple is an auto generated Go binding around an Ethereum contract.
type Simple[P crypto.PublicKey] struct {
	SimpleCaller[P]     // Read-only binding to the contract
	SimpleTransactor[P] // Write-only binding to the contract
	SimpleFilterer[P]   // Log filterer for contract events
}

// SimpleCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleCaller[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// SimpleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleTransactor[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// SimpleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleFilterer[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// SimpleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleSession[P crypto.PublicKey] struct {
	Contract     *Simple[P]           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts[P] // Transaction auth options to use throughout this session
}

// SimpleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleCallerSession[P crypto.PublicKey] struct {
	Contract *SimpleCaller[P] // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// SimpleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleTransactorSession[P crypto.PublicKey] struct {
	Contract     *SimpleTransactor[P] // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts[P] // Transaction auth options to use throughout this session
}

// SimpleRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleRaw[P crypto.PublicKey] struct {
	Contract *Simple[P] // Generic contract binding to access the raw methods on
}

// SimpleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleCallerRaw[P crypto.PublicKey] struct {
	Contract *SimpleCaller[P] // Generic read-only contract binding to access the raw methods on
}

// SimpleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleTransactorRaw[P crypto.PublicKey] struct {
	Contract *SimpleTransactor[P] // Generic write-only contract binding to access the raw methods on
}

// NewSimple creates a new instance of Simple, bound to a specific deployed contract.
func NewSimple[P crypto.PublicKey](address common.Address, backend bind.ContractBackend[P]) (*Simple[P], error) {
	contract, err := bindSimple[P](address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Simple[P]{SimpleCaller: SimpleCaller[P]{contract: contract}, SimpleTransactor: SimpleTransactor[P]{contract: contract}, SimpleFilterer: SimpleFilterer[P]{contract: contract}}, nil
}

// NewSimpleCaller creates a new read-only instance of Simple, bound to a specific deployed contract.
func NewSimpleCaller[P crypto.PublicKey](address common.Address, caller bind.ContractCaller) (*SimpleCaller[P], error) {
	contract, err := bindSimple[P](address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleCaller[P]{contract: contract}, nil
}

// NewSimpleTransactor creates a new write-only instance of Simple, bound to a specific deployed contract.
func NewSimpleTransactor[P crypto.PublicKey](address common.Address, transactor bind.ContractTransactor[P]) (*SimpleTransactor[P], error) {
	contract, err := bindSimple[P](address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleTransactor[P]{contract: contract}, nil
}

// NewSimpleFilterer creates a new log filterer instance of Simple, bound to a specific deployed contract.
func NewSimpleFilterer[P crypto.PublicKey](address common.Address, filterer bind.ContractFilterer) (*SimpleFilterer[P], error) {
	contract, err := bindSimple[P](address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleFilterer[P]{contract: contract}, nil
}

// bindSimple binds a generic wrapper to an already deployed contract.
func bindSimple[P crypto.PublicKey](address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor[P], filterer bind.ContractFilterer) (*bind.BoundContract[P], error) {
	parsed, err := abi.JSON[P](strings.NewReader(SimpleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Simple *SimpleRaw[P]) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Simple.Contract.SimpleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Simple *SimpleRaw[P]) Transfer(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _Simple.Contract.SimpleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Simple *SimpleRaw[P]) Transact(opts *bind.TransactOpts[P], method string, params ...interface{}) (*types.Transaction[P], error) {
	return _Simple.Contract.SimpleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Simple *SimpleCallerRaw[P]) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Simple.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Simple *SimpleTransactorRaw[P]) Transfer(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _Simple.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Simple *SimpleTransactorRaw[P]) Transact(opts *bind.TransactOpts[P], method string, params ...interface{}) (*types.Transaction[P], error) {
	return _Simple.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x07ea02e1.
//
// Solidity: function getValue() view returns(bytes16)
func (_Simple *SimpleCaller[P]) GetValue(opts *bind.CallOpts) ([16]byte, error) {
	var out []interface{}
	err := _Simple.contract.Call(opts, &out, "getValue")

	if err != nil {
		return *new([16]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([16]byte)).(*[16]byte)

	return out0, err

}

// GetValue is a free data retrieval call binding the contract method 0x07ea02e1.
//
// Solidity: function getValue() view returns(bytes16)
func (_Simple *SimpleSession[P]) GetValue() ([16]byte, error) {
	return _Simple.Contract.GetValue(&_Simple.CallOpts)
}

// GetValue is a free data retrieval call binding the contract method 0x07ea02e1.
//
// Solidity: function getValue() view returns(bytes16)
func (_Simple *SimpleCallerSession[P]) GetValue() ([16]byte, error) {
	return _Simple.Contract.GetValue(&_Simple.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x7d7ee50e.
//
// Solidity: function value() view returns(bytes16)
func (_Simple *SimpleCaller[P]) Value(opts *bind.CallOpts) ([16]byte, error) {
	var out []interface{}
	err := _Simple.contract.Call(opts, &out, "value")

	if err != nil {
		return *new([16]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([16]byte)).(*[16]byte)

	return out0, err

}

// Value is a free data retrieval call binding the contract method 0x7d7ee50e.
//
// Solidity: function value() view returns(bytes16)
func (_Simple *SimpleSession[P]) Value() ([16]byte, error) {
	return _Simple.Contract.Value(&_Simple.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x7d7ee50e.
//
// Solidity: function value() view returns(bytes16)
func (_Simple *SimpleCallerSession[P]) Value() ([16]byte, error) {
	return _Simple.Contract.Value(&_Simple.CallOpts)
}

// SetValue is a paid mutator transaction binding the contract method 0xc7ff2103.
//
// Solidity: function setValue(bytes16 v) returns()
func (_Simple *SimpleTransactor[P]) SetValue(opts *bind.TransactOpts[P], v [16]byte) (*types.Transaction[P], error) {
	return _Simple.contract.Transact(opts, "setValue", v)
}

// SetValue is a paid mutator transaction binding the contract method 0xc7ff2103.
//
// Solidity: function setValue(bytes16 v) returns()
func (_Simple *SimpleSession[P]) SetValue(v [16]byte) (*types.Transaction[P], error) {
	return _Simple.Contract.SetValue(&_Simple.TransactOpts, v)
}

// SetValue is a paid mutator transaction binding the contract method 0xc7ff2103.
//
// Solidity: function setValue(bytes16 v) returns()
func (_Simple *SimpleTransactorSession[P]) SetValue(v [16]byte) (*types.Transaction[P], error) {
	return _Simple.Contract.SetValue(&_Simple.TransactOpts, v)
}
