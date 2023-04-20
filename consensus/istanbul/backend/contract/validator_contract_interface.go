// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	ethereum "github.com/pavelkrolevets/MIR-pro"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	// _ = bind.Bind
	_ = common.Big1
	// _ = types.BloomLookup
	_ = event.NewSubscription
)

// ValidatorContractInterfaceABI is the input ABI used to generate the binding from.
const ValidatorContractInterfaceABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// var ValidatorContractInterfaceParsedABI, _ = abi.JSON(strings.NewReader(ValidatorContractInterfaceABI))

// ValidatorContractInterface is an auto generated Go binding around an Ethereum contract.
type ValidatorContractInterface [P crypto.PublicKey] struct {
	ValidatorContractInterfaceCaller[P]     // Read-only binding to the contract
	ValidatorContractInterfaceTransactor[P] // Write-only binding to the contract
	ValidatorContractInterfaceFilterer[P]   // Log filterer for contract events
}

// ValidatorContractInterfaceCaller[P] is an auto generated read-only Go binding around an Ethereum contract.
type ValidatorContractInterfaceCaller [P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// ValidatorContractInterfaceTransactor[P] is an auto generated write-only Go binding around an Ethereum contract.
type ValidatorContractInterfaceTransactor[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// ValidatorContractInterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatorContractInterfaceFilterer[P crypto.PublicKey] struct {
	contract *bind.BoundContract[P] // Generic contract wrapper for the low level calls
}

// ValidatorContractInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatorContractInterfaceSession [P crypto.PublicKey] struct {
	Contract     *ValidatorContractInterface[P] // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts[P]         // Transaction auth options to use throughout this session
}

// ValidatorContractInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatorContractInterfaceCallerSession [P crypto.PublicKey] struct {
	Contract *ValidatorContractInterfaceCaller[P]// Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// ValidatorContractInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatorContractInterfaceTransactorSession[P crypto.PublicKey]  struct {
	Contract     *ValidatorContractInterfaceTransactor[P] // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts[P]                   // Transaction auth options to use throughout this session
}

// ValidatorContractInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatorContractInterfaceRaw [P crypto.PublicKey] struct {
	Contract *ValidatorContractInterface[P] // Generic contract binding to access the raw methods on
}

// ValidatorContractInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatorContractInterfaceCallerRaw [P crypto.PublicKey] struct {
	Contract *ValidatorContractInterfaceCaller[P] // Generic read-only contract binding to access the raw methods on
}

// ValidatorContractInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatorContractInterfaceTransactorRaw [P crypto.PublicKey] struct {
	Contract *ValidatorContractInterfaceTransactor[P] // Generic write-only contract binding to access the raw methods on
}

// NewValidatorContractInterface creates a new instance of ValidatorContractInterface, bound to a specific deployed contract.
func NewValidatorContractInterface[P crypto.PublicKey](address common.Address, backend bind.ContractBackend[P]) (*ValidatorContractInterface[P], error) {
	contract, err := bindValidatorContractInterface[P](address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ValidatorContractInterface[P]{ValidatorContractInterfaceCaller: ValidatorContractInterfaceCaller[P]{contract: contract}, ValidatorContractInterfaceTransactor: ValidatorContractInterfaceTransactor[P]{contract: contract}, ValidatorContractInterfaceFilterer: ValidatorContractInterfaceFilterer[P]{contract: contract}}, nil
}

// NewValidatorContractInterfaceCaller creates a new read-only instance of ValidatorContractInterface, bound to a specific deployed contract.
func NewValidatorContractInterfaceCaller[P crypto.PublicKey](address common.Address, caller bind.ContractCaller) (*ValidatorContractInterfaceCaller[P], error) {
	contract, err := bindValidatorContractInterface[P](address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorContractInterfaceCaller[P]{contract: contract}, nil
}

// NewValidatorContractInterfaceTransactor creates a new write-only instance of ValidatorContractInterface, bound to a specific deployed contract.
func NewValidatorContractInterfaceTransactor[P crypto.PublicKey](address common.Address, transactor bind.ContractTransactor[P]) (*ValidatorContractInterfaceTransactor[P], error) {
	contract, err := bindValidatorContractInterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorContractInterfaceTransactor[P]{contract: contract}, nil
}

// NewValidatorContractInterfaceFilterer creates a new log filterer instance of ValidatorContractInterface, bound to a specific deployed contract.
func NewValidatorContractInterfaceFilterer[P crypto.PublicKey](address common.Address, filterer bind.ContractFilterer) (*ValidatorContractInterfaceFilterer[P], error) {
	contract, err := bindValidatorContractInterface[P](address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatorContractInterfaceFilterer[P]{contract: contract}, nil
}

// bindValidatorContractInterface binds a generic wrapper to an already deployed contract.
func bindValidatorContractInterface[P crypto.PublicKey](address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor[P], filterer bind.ContractFilterer) (*bind.BoundContract[P], error) {
	parsed, err := abi.JSON[P](strings.NewReader(ValidatorContractInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorContractInterface *ValidatorContractInterfaceRaw[P]) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidatorContractInterface.Contract.ValidatorContractInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorContractInterface *ValidatorContractInterfaceRaw[P]) Transfer(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _ValidatorContractInterface.Contract.ValidatorContractInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorContractInterface *ValidatorContractInterfaceRaw[P]) Transact(opts *bind.TransactOpts[P], method string, params ...interface{}) (*types.Transaction[P], error) {
	return _ValidatorContractInterface.Contract.ValidatorContractInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorContractInterface *ValidatorContractInterfaceCallerRaw[P]) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidatorContractInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorContractInterface *ValidatorContractInterfaceTransactorRaw[P]) Transfer(opts *bind.TransactOpts[P]) (*types.Transaction[P], error) {
	return _ValidatorContractInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorContractInterface *ValidatorContractInterfaceTransactorRaw[P]) Transact(opts *bind.TransactOpts[P], method string, params ...interface{}) (*types.Transaction[P], error) {
	return _ValidatorContractInterface.Contract.contract.Transact(opts, method, params...)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_ValidatorContractInterface *ValidatorContractInterfaceCaller[P]) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _ValidatorContractInterface.contract.Call(opts, &out, "getValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_ValidatorContractInterface *ValidatorContractInterfaceSession[P]) GetValidators() ([]common.Address, error) {
	return _ValidatorContractInterface.Contract.GetValidators(&_ValidatorContractInterface.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_ValidatorContractInterface *ValidatorContractInterfaceCallerSession[P]) GetValidators() ([]common.Address, error) {
	return _ValidatorContractInterface.Contract.GetValidators(&_ValidatorContractInterface.CallOpts)
}
