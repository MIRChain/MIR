// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package simple

import (
	"math/big"
	"strings"

	ethereum "github.com/pavelkrolevets/MIR-pro"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SimpleABI is the input ABI used to generate the binding from.
const SimpleABI = "[{\"inputs\":[],\"name\":\"getValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"v\",\"type\":\"uint256\"}],\"name\":\"setValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"value\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

var SimpleParsedABI, _ = abi.JSON(strings.NewReader(SimpleABI))

// SimpleFuncSigs maps the 4-byte function signature to its string representation.
var SimpleFuncSigs = map[string]string{
	"20965255": "getValue()",
	"55241077": "setValue(uint256)",
	"3fa4f245": "value()",
}

// SimpleBin is the compiled bytecode used for deploying new contracts.
var SimpleBin = "0x608060405234801561001057600080fd5b5060bf8061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c8063209652551460415780633fa4f2451460575780635524107714605f575b600080fd5b6000545b60405190815260200160405180910390f35b604560005481565b606f606a3660046071565b600055565b005b600060208284031215608257600080fd5b503591905056fea2646970667358221220a614845e77b9c01d3ee20c8ed15a3bf3f6fbc3485b50904da48f58dcc97e256c64736f6c634300080e0033"

// DeploySimple deploys a new Ethereum contract, binding an instance of Simple to it.
func DeploySimple(auth *bind.TransactOpts[nist.PublicKey], backend bind.ContractBackend[nist.PublicKey]) (common.Address, *types.Transaction[nist.PublicKey], *Simple, error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SimpleBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Simple{SimpleCaller: SimpleCaller{contract: contract}, SimpleTransactor: SimpleTransactor{contract: contract}, SimpleFilterer: SimpleFilterer{contract: contract}}, nil
}

// Simple is an auto generated Go binding around an Ethereum contract.
type Simple struct {
	SimpleCaller     // Read-only binding to the contract
	SimpleTransactor // Write-only binding to the contract
	SimpleFilterer   // Log filterer for contract events
}

// SimpleCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleCaller struct {
	contract *bind.BoundContract[nist.PublicKey] // Generic contract wrapper for the low level calls
}

// SimpleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleTransactor struct {
	contract *bind.BoundContract[nist.PublicKey] // Generic contract wrapper for the low level calls
}

// SimpleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleFilterer struct {
	contract *bind.BoundContract[nist.PublicKey] // Generic contract wrapper for the low level calls
}

// SimpleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleSession struct {
	Contract     *Simple           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts[nist.PublicKey] // Transaction auth options to use throughout this session
}

// SimpleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleCallerSession struct {
	Contract *SimpleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SimpleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleTransactorSession struct {
	Contract     *SimpleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts[nist.PublicKey] // Transaction auth options to use throughout this session
}

// SimpleRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleRaw struct {
	Contract *Simple // Generic contract binding to access the raw methods on
}

// SimpleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleCallerRaw struct {
	Contract *SimpleCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleTransactorRaw struct {
	Contract *SimpleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimple creates a new instance of Simple, bound to a specific deployed contract.
func NewSimple(address common.Address, backend bind.ContractBackend[nist.PublicKey]) (*Simple, error) {
	contract, err := bindSimple(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Simple{SimpleCaller: SimpleCaller{contract: contract}, SimpleTransactor: SimpleTransactor{contract: contract}, SimpleFilterer: SimpleFilterer{contract: contract}}, nil
}

// NewSimpleCaller creates a new read-only instance of Simple, bound to a specific deployed contract.
func NewSimpleCaller(address common.Address, caller bind.ContractCaller) (*SimpleCaller, error) {
	contract, err := bindSimple(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleCaller{contract: contract}, nil
}

// NewSimpleTransactor creates a new write-only instance of Simple, bound to a specific deployed contract.
func NewSimpleTransactor(address common.Address, transactor bind.ContractTransactor[nist.PublicKey]) (*SimpleTransactor, error) {
	contract, err := bindSimple(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleTransactor{contract: contract}, nil
}

// NewSimpleFilterer creates a new log filterer instance of Simple, bound to a specific deployed contract.
func NewSimpleFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleFilterer, error) {
	contract, err := bindSimple(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleFilterer{contract: contract}, nil
}

// bindSimple binds a generic wrapper to an already deployed contract.
func bindSimple(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor[nist.PublicKey], filterer bind.ContractFilterer) (*bind.BoundContract[nist.PublicKey], error) {
	parsed, err := abi.JSON(strings.NewReader(SimpleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Simple *SimpleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Simple.Contract.SimpleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Simple *SimpleRaw) Transfer(opts *bind.TransactOpts[nist.PublicKey]) (*types.Transaction[nist.PublicKey], error) {
	return _Simple.Contract.SimpleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Simple *SimpleRaw) Transact(opts *bind.TransactOpts[nist.PublicKey], method string, params ...interface{}) (*types.Transaction[nist.PublicKey], error) {
	return _Simple.Contract.SimpleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Simple *SimpleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Simple.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Simple *SimpleTransactorRaw) Transfer(opts *bind.TransactOpts[nist.PublicKey]) (*types.Transaction[nist.PublicKey], error) {
	return _Simple.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Simple *SimpleTransactorRaw) Transact(opts *bind.TransactOpts[nist.PublicKey], method string, params ...interface{}) (*types.Transaction[nist.PublicKey], error) {
	return _Simple.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_Simple *SimpleCaller) GetValue(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Simple.contract.Call(opts, &out, "getValue")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_Simple *SimpleSession) GetValue() (*big.Int, error) {
	return _Simple.Contract.GetValue(&_Simple.CallOpts)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() view returns(uint256)
func (_Simple *SimpleCallerSession) GetValue() (*big.Int, error) {
	return _Simple.Contract.GetValue(&_Simple.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_Simple *SimpleCaller) Value(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Simple.contract.Call(opts, &out, "value")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_Simple *SimpleSession) Value() (*big.Int, error) {
	return _Simple.Contract.Value(&_Simple.CallOpts)
}

// Value is a free data retrieval call binding the contract method 0x3fa4f245.
//
// Solidity: function value() view returns(uint256)
func (_Simple *SimpleCallerSession) Value() (*big.Int, error) {
	return _Simple.Contract.Value(&_Simple.CallOpts)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 v) returns()
func (_Simple *SimpleTransactor) SetValue(opts *bind.TransactOpts[nist.PublicKey], v *big.Int) (*types.Transaction[nist.PublicKey], error) {
	return _Simple.contract.Transact(opts, "setValue", v)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 v) returns()
func (_Simple *SimpleSession) SetValue(v *big.Int) (*types.Transaction[nist.PublicKey], error) {
	return _Simple.Contract.SetValue(&_Simple.TransactOpts, v)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 v) returns()
func (_Simple *SimpleTransactorSession) SetValue(v *big.Int) (*types.Transaction[nist.PublicKey], error) {
	return _Simple.Contract.SetValue(&_Simple.TransactOpts, v)
}
