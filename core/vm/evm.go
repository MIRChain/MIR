// Copyright 2014 The go-ethereum Authors
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

package vm

import (
	"errors"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/holiman/uint256"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

// note: Quorum, States, and Value Transfer
//
// In Quorum there is a tricky issue in one specific case when there is call from private state to public state:
// * The state db is selected based on the callee (public)
// * With every call there is an associated value transfer -- in our case this is 0
// * Thus, there is an implicit transfer of 0 value from the caller to callee on the public state
// * However in our scenario the caller is private
// * Thus, the transfer creates a ghost of the private account on the public state with no value, code, or storage
//
// The solution is to skip this transfer of 0 value under Quorum

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
// var emptyCodeHash = crypto.Keccak256Hash(nil)

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, common.Address, common.Address, *big.Int)
	// GetHashFunc returns the n'th block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

func (evm *EVM[P]) precompile(addr common.Address) (PrecompiledContract, bool) {
	var precompiles map[common.Address]PrecompiledContract
	switch {
	case evm.chainRules.IsBerlin:
		precompiles = PrecompiledContractsBerlin
	case evm.chainRules.IsIstanbul:
		precompiles = PrecompiledContractsIstanbul
	case evm.chainRules.IsByzantium:
		precompiles = PrecompiledContractsByzantium
	default:
		precompiles = PrecompiledContractsHomestead
	}
	p, ok := precompiles[addr]
	return p, ok
}

// Quorum
func (evm *EVM[P]) quorumPrecompile(addr common.Address) (QuorumPrecompiledContract[P], bool) {
	var quorumPrecompiles map[common.Address]QuorumPrecompiledContract[P]
	switch {
	case evm.chainRules.IsPrivacyPrecompile:
		quorumPrecompiles = map[common.Address]QuorumPrecompiledContract[P]{
			common.QuorumPrivacyPrecompileContractAddress(): &privacyMarker[P]{},
		}
	}
	p, ok := quorumPrecompiles[addr]
	return p, ok
}

// End Quorum

// run runs the given contract and takes care of running precompiles with a fallback to the byte code interpreter.
func run[P crypto.PublicKey](evm *EVM[P], contract *Contract[P], input []byte, readOnly bool) ([]byte, error) {
	// Quorum
	if contract.CodeAddr != nil {
		// Using CodeAddr is favour over contract.Address()
		// During DelegateCall() CodeAddr is the address of the delegated account
		address := *contract.CodeAddr
		if _, ok := evm.affectedContracts[address]; !ok {
			evm.affectedContracts[address] = MessageCall
		}
	}
	// End Quorum
	for _, interpreter := range evm.interpreters {
		if interpreter.CanRun(contract.Code) {
			if evm.interpreter != interpreter {
				// Ensure that the interpreter pointer is set back
				// to its current value upon return.
				defer func(i Interpreter[P]) {
					evm.interpreter = i
				}(evm.interpreter)
				evm.interpreter = interpreter
			}
			return interpreter.Run(contract, input, readOnly)
		}
	}
	return nil, errors.New("no compatible interpreter")
}

// BlockContext provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type BlockContext struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
}

// TxContext provides the EVM with information about a transaction.
// All fields can change between transactions.
type TxContext struct {
	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE
}

// Quorum
type PublicState StateDB
type PrivateState StateDB

// End Quorum

// EVM is the Ethereum Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The EVM should never be reused and is not thread safe.
type EVM [P crypto.PublicKey] struct {
	// Context provides auxiliary blockchain related information
	Context BlockContext
	TxContext
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	vmConfig Config[P]
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	interpreters []Interpreter[P]
	interpreter  Interpreter[P]
	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64

	// Quorum additions:
	publicState       PublicState
	privateState      PrivateState
	states            [1027]*state.StateDB[P] // TODO(joel) we should be able to get away with 1024 or maybe 1025
	currentStateDepth uint

	// This flag has different semantics from the `Interpreter:readOnly` flag (though they interact and could maybe
	// be simplified). This is set by Quorum when it's inside a Private State -> Public State read.
	quorumReadOnly bool
	readOnlyDepth  uint

	// Quorum: these are for privacy enhancements and multitenancy
	affectedContracts map[common.Address]AffectedReason // affected contract account address -> type
	currentTx         *types.Transaction[P]                // transaction currently being applied on this EVM

	// Quorum: these are for privacy marker transactions
	InnerApply          func(innerTx *types.Transaction[P]) error //Quorum
	InnerPrivateReceipt *types.Receipt[P]                         //Quorum
}

// AffectedReason defines a type of operation that was applied to a contract.
type AffectedReason byte

const (
	_        AffectedReason = iota
	Creation AffectedReason = iota
	MessageCall
)

// NewEVM returns a new EVM. The returned EVM is not thread safe and should
// only ever be used *once*.
func NewEVM[P crypto.PublicKey](blockCtx BlockContext, txCtx TxContext, statedb, privateState StateDB, chainConfig *params.ChainConfig, vmConfig Config[P]) *EVM[P] {
	evm := &EVM[P]{
		Context:      blockCtx,
		TxContext:    txCtx,
		StateDB:      statedb,
		vmConfig:     vmConfig,
		chainConfig:  chainConfig,
		chainRules:   chainConfig.Rules(blockCtx.BlockNumber),
		interpreters: make([]Interpreter[P], 0, 1),

		publicState:  statedb,
		privateState: privateState,

		affectedContracts: make(map[common.Address]AffectedReason),
	}

	if chainConfig.IsEWASM(blockCtx.BlockNumber) {
		// to be implemented by EVM-C and Wagon PRs.
		// if vmConfig.EWASMInterpreter != "" {
		//  extIntOpts := strings.Split(vmConfig.EWASMInterpreter, ":")
		//  path := extIntOpts[0]
		//  options := []string{}
		//  if len(extIntOpts) > 1 {
		//    options = extIntOpts[1..]
		//  }
		//  evm.interpreters = append(evm.interpreters, NewEVMVCInterpreter(evm, vmConfig, options))
		// } else {
		// 	evm.interpreters = append(evm.interpreters, NewEWASMInterpreter(evm, vmConfig))
		// }
		panic("No supported ewasm interpreter yet.")
	}

	evm.Push(privateState)

	// vmConfig.EVMInterpreter will be used by EVM-C, it won't be checked here
	// as we always want to have the built-in EVM as the failover option.
	evm.interpreters = append(evm.interpreters, NewEVMInterpreter(evm, vmConfig))
	evm.interpreter = evm.interpreters[0]

	return evm
}

// Reset resets the EVM with a new transaction context.Reset
// This is not threadsafe and should only be done very cautiously.
func (evm *EVM[P]) Reset(txCtx TxContext, statedb StateDB, privateStateDB StateDB) {
	evm.TxContext = txCtx
	evm.StateDB = statedb
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM[P]) Cancel() {
	atomic.StoreInt32(&evm.abort, 1)
}

// Cancelled returns true if Cancel has been called
func (evm *EVM[P]) Cancelled() bool {
	return atomic.LoadInt32(&evm.abort) == 1
}

// Interpreter returns the current interpreter
func (evm *EVM[P]) Interpreter() Interpreter[P] {
	return evm.interpreter
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (evm *EVM[P]) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	evm.Push(getDualState(evm, addr))
	defer func() { evm.Pop() }()

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if value.Sign() != 0 && !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	snapshot := evm.StateDB.Snapshot()
	p, isPrecompile := evm.precompile(addr)
	qp, isQuorumPrecompile := evm.quorumPrecompile(addr) // Quorum

	if !evm.StateDB.Exist(addr) {
		if !isPrecompile && !isQuorumPrecompile && evm.chainRules.IsEIP158 && value.Sign() == 0 {
			// Calling a non existing account, don't do anything, but ping the tracer
			if evm.vmConfig.Debug && evm.depth == 0 {
				evm.vmConfig.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
				evm.vmConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, gas, nil
		}
		// If we are executing the quorum PMT precompile, then don't add it to state.
		// (When executing the PMT precompile, we are using private state - adding the account can cause differences in private state root when using with MPS.)
		if !isQuorumPrecompile {
			evm.StateDB.CreateAccount(addr)
		}
	}

	// Quorum
	if evm.ChainConfig().IsQuorum {
		// skip transfer if value == 0 (see note: Quorum, States, and Value Transfer)
		if value.Sign() != 0 {
			if evm.quorumReadOnly {
				return nil, gas, ErrReadOnlyValueTransfer
			}
			evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)
		}
		// End Quorum
	} else {
		evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)
	}

	// Capture the tracer start/end events in debug mode
	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
		defer func(startGas uint64, startTime time.Time) { // Lazy evaluation of the parameters
			evm.vmConfig.Tracer.CaptureEnd(ret, startGas-gas, time.Since(startTime), err)
		}(gas, time.Now())
	}

	if isQuorumPrecompile {
		ret, gas, err = RunQuorumPrecompiledContract(evm, qp, input, gas)
	} else if isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		code := evm.StateDB.GetCode(addr)
		addrCopy := addr
		// If the account has no code, we can abort here
		// The depth-check is already done, and precompiles handled above
		contract := NewContract[P](caller, AccountRef(addrCopy), value, gas)
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), code)
		ret, err = run(evm, contract, input, false)
		gas = contract.Gas
	}
	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
		// TODO: consider clearing up unused snapshots:
		//} else {
		//	evm.StateDB.DiscardSnapshot(snapshot)
	}
	return ret, gas, err
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (evm *EVM[P]) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	// Quorum
	evm.Push(getDualState(evm, addr))
	defer func() { evm.Pop() }()
	// End Quorum

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	// Note although it's noop to transfer X ether to caller itself. But
	// if caller doesn't have enough balance, it would be an error to allow
	// over-charging itself. So the check here is necessary.
	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	var snapshot = evm.StateDB.Snapshot()

	// It is allowed to call precompiles, even via delegatecall
	if qp, isQuorumPrecompile := evm.quorumPrecompile(addr); isQuorumPrecompile { // Quorum
		ret, gas, err = RunQuorumPrecompiledContract(evm, qp, input, gas)
	} else if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		addrCopy := addr
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		contract := NewContract[P](caller, AccountRef(caller.Address()), value, gas)
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		ret, err = run(evm, contract, input, false)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (evm *EVM[P]) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	// Quorum
	evm.Push(getDualState(evm, addr))
	defer func() { evm.Pop() }()
	// End Quorum

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	var snapshot = evm.StateDB.Snapshot()

	// It is allowed to call precompiles, even via delegatecall
	if qp, isQuorumPrecompile := evm.quorumPrecompile(addr); isQuorumPrecompile { // Quorum
		ret, gas, err = RunQuorumPrecompiledContract(evm, qp, input, gas)
	} else if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		addrCopy := addr
		// Initialise a new contract and make initialise the delegate values
		contract := NewContract[P](caller, AccountRef(caller.Address()), nil, gas).AsDelegate()
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		ret, err = run(evm, contract, input, false)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (evm *EVM[P]) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Quorum
	// use the right state (public or private)
	stateDb := getDualState(evm, addr)
	// End Quorum

	// We take a snapshot here. This is a bit counter-intuitive, and could probably be skipped.
	// However, even a staticcall is considered a 'touch'. On mainnet, static calls were introduced
	// after all empty accounts were deleted, so this is not required. However, if we omit this,
	// then certain tests start failing; stRevertTest/RevertPrecompiledTouchExactOOG.json.
	// We could change this, but for now it's left for legacy reasons
	var snapshot = stateDb.Snapshot()

	// We do an AddBalance of zero here, just in order to trigger a touch.
	// This doesn't matter on Mainnet, where all empties are gone at the time of Byzantium,
	// but is the correct thing to do and matters on other networks, in tests, and potential
	// future scenarios
	stateDb.AddBalance(addr, big0)

	if qp, isQuorumPrecompile := evm.quorumPrecompile(addr); isQuorumPrecompile { // Quorum
		ret, gas, err = RunQuorumPrecompiledContract(evm, qp, input, gas)
	} else if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// At this point, we use a copy of address. If we don't, the go compiler will
		// leak the 'contract' to the outer scope, and make allocation for 'contract'
		// even if the actual execution ends on RunPrecompiled above.
		addrCopy := addr
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		contract := NewContract[P](caller, AccountRef(addrCopy), new(big.Int), gas)
		contract.SetCallCode(&addrCopy, stateDb.GetCodeHash(addrCopy), stateDb.GetCode(addrCopy))
		// When an error was returned by the EVM or when setting the creation code
		// above we revert to the snapshot and consume any gas remaining. Additionally
		// when we're in Homestead this also counts for code storage gas errors.
		ret, err = run(evm, contract, input, true)
		gas = contract.Gas
	}
	if err != nil {
		stateDb.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

type codeAndHash [P crypto.PublicKey] struct {
	code []byte
	hash common.Hash
}

func (c *codeAndHash[P]) Hash() common.Hash {
	if c.hash == (common.Hash{}) {
		c.hash = crypto.Keccak256Hash[P](c.code)
	}
	return c.hash
}

// create creates a new contract using code as deployment code.
func (evm *EVM[P]) create(caller ContractRef, codeAndHash *codeAndHash[P], gas uint64, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if evm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gas, ErrDepth
	}
	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, common.Address{}, gas, ErrInsufficientBalance
	}

	// We add this to the access list _before_ taking a snapshot. Even if the creation fails,
	// the access-list change should not be rolled back
	if evm.chainRules.IsBerlin {
		evm.StateDB.AddAddressToAccessList(address)
	}

	// Quorum
	// Get the right state in case of a dual state environment. If a sender
	// is a transaction (depth == 0) use the public state to derive the address
	// and increment the nonce of the public state. If the sender is a contract
	// (depth > 0) use the private state to derive the nonce and increment the
	// nonce on the private state only.
	//
	// If the transaction went to a public contract the private and public state
	// are the same.
	var creatorStateDb StateDB
	if evm.depth > 0 {
		creatorStateDb = evm.privateState
	} else {
		creatorStateDb = evm.publicState
	}

	nonce := creatorStateDb.GetNonce(caller.Address())
	creatorStateDb.SetNonce(caller.Address(), nonce+1)

	// Ensure there's no existing contract already at the designated address
	contractHash := evm.StateDB.GetCodeHash(address)
	if evm.StateDB.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != crypto.Keccak256Hash[P](nil)) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(address)
	// Quorum
	evm.affectedContracts[address] = Creation
	// End Quorum
	if evm.chainRules.IsEIP158 {
		evm.StateDB.SetNonce(address, 1)
	}
	if evm.currentTx != nil && evm.currentTx.IsPrivate() && evm.currentTx.PrivacyMetadata() != nil {
		// for calls (reading contract state) or finding the affected contracts there is no transaction
		if evm.currentTx.PrivacyMetadata().PrivacyFlag.IsNotStandardPrivate() {
			pm := state.NewStatePrivacyMetadata(common.BytesToEncryptedPayloadHash(evm.currentTx.Data()), evm.currentTx.PrivacyMetadata().PrivacyFlag)
			evm.StateDB.SetPrivacyMetadata(address, pm)
			log.Trace("Set Privacy Metadata", "key", address, "privacyMetadata", pm)
		}
	}
	if evm.ChainConfig().IsQuorum {
		// skip transfer if value == 0 (see note: Quorum, States, and Value Transfer)
		if value.Sign() != 0 {
			if evm.quorumReadOnly {
				return nil, common.Address{}, gas, ErrReadOnlyValueTransfer
			}
			evm.Context.Transfer(evm.StateDB, caller.Address(), address, value)
		}
	} else {
		evm.Context.Transfer(evm.StateDB, caller.Address(), address, value)
	}

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract[P](caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, address, gas, nil
	}

	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(evm, caller.Address(), address, true, codeAndHash.code, gas, value)
	}
	start := time.Now()

	ret, err := run(evm, contract, nil, false)

	maxCodeSize := evm.ChainConfig().GetMaxCodeSize(evm.Context.BlockNumber)
	if maxCodeSize < params.MaxCodeSize {
		maxCodeSize = params.MaxCodeSize
	}

	// Check whether the max code size has been exceeded, assign err if the case.
	if err == nil && evm.chainRules.IsEIP158 && len(ret) > maxCodeSize {
		err = ErrMaxCodeSizeExceeded
	}

	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			evm.StateDB.SetCode(address, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil && (evm.chainRules.IsHomestead || err != ErrCodeStoreOutOfGas) {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}
	return ret, address, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM[P]) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	// Quorum
	// Get the right state in case of a dual state environment. If a sender
	// is a transaction (depth == 0) use the public state to derive the address
	// and increment the nonce of the public state. If the sender is a contract
	// (depth > 0) use the private state to derive the nonce and increment the
	// nonce on the private state only.
	//
	// If the transaction went to a public contract the private and public state
	// are the same.
	var creatorStateDb StateDB
	if evm.depth > 0 {
		creatorStateDb = evm.privateState
	} else {
		creatorStateDb = evm.publicState
	}

	// Ensure there's no existing contract already at the designated address
	nonce := creatorStateDb.GetNonce(caller.Address())
	contractAddr = crypto.CreateAddress[P](caller.Address(), nonce)
	return evm.create(caller, &codeAndHash[P]{code: code}, gas, value, contractAddr)
}

// Create2 creates a new contract using code as deployment code.
//
// The different between Create2 with Create is Create2 uses sha3(0xff ++ msg.sender ++ salt ++ sha3(init_code))[12:]
// instead of the usual sender-and-nonce-hash as the address where the contract is initialized at.
func (evm *EVM[P]) Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int, salt *uint256.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	codeAndHash := &codeAndHash[P]{code: code}
	contractAddr = crypto.CreateAddress2[P](caller.Address(), salt.Bytes32(), codeAndHash.Hash().Bytes())
	return evm.create(caller, codeAndHash, gas, endowment, contractAddr)
}

// ChainConfig returns the environment's chain configuration
func (evm *EVM[P]) ChainConfig() *params.ChainConfig { return evm.chainConfig }

// Quorum functions for dual state
func getDualState[P crypto.PublicKey](evm *EVM[P], addr common.Address) StateDB {
	// priv: (a) -> (b)  (private)
	// pub:   a  -> [b]  (private -> public)
	// priv: (a) ->  b   (public)
	state := evm.StateDB

	if evm.PrivateState().Exist(addr) {
		state = evm.PrivateState()
	} else if evm.PublicState().Exist(addr) {
		state = evm.PublicState()
	}

	return state
}

func (evm *EVM[P]) PublicState() PublicState           { return evm.publicState }
func (evm *EVM[P]) PrivateState() PrivateState         { return evm.privateState }
func (evm *EVM[P]) SetCurrentTX(tx *types.Transaction[P]) { evm.currentTx = tx }
func (evm *EVM[P]) SetTxPrivacyMetadata(pm *types.PrivacyMetadata) {
	evm.currentTx.SetTxPrivacyMetadata(pm)
}
func (evm *EVM[P]) Push(statedb StateDB) {
	// Quorum : the read only depth to be set up only once for the entire
	// op code execution. This will be set first time transition from
	// private state to public state happens
	// statedb will be the state of the contract being called.
	// if a private contract is calling a public contract make it readonly.
	if !evm.quorumReadOnly && evm.privateState != statedb {
		evm.quorumReadOnly = true
		evm.readOnlyDepth = evm.currentStateDepth
	}

	if castedStateDb, ok := statedb.(*state.StateDB[P]); ok {
		evm.states[evm.currentStateDepth] = castedStateDb
		evm.currentStateDepth++
	}

	evm.StateDB = statedb
}
func (evm *EVM[P]) Pop() {
	evm.currentStateDepth--
	if evm.quorumReadOnly && evm.currentStateDepth == evm.readOnlyDepth {
		evm.quorumReadOnly = false
	}
	evm.StateDB = evm.states[evm.currentStateDepth-1]
}

func (evm *EVM[P]) Depth() int { return evm.depth }

// We only need to revert the current state because when we call from private
// public state it's read only, there wouldn't be anything to reset.
// (A)->(B)->C->(B): A failure in (B) wouldn't need to reset C, as C was flagged
// read only.
func (evm *EVM[P]) RevertToSnapshot(snapshot int) {
	evm.StateDB.RevertToSnapshot(snapshot)
}

// Quorum
//
// Returns addresses of contracts which are newly created
func (evm *EVM[P]) CreatedContracts() []common.Address {
	addr := make([]common.Address, 0, len(evm.affectedContracts))
	for a, t := range evm.affectedContracts {
		if t == Creation {
			addr = append(addr, a)
		}
	}
	return addr[:]
}

// Quorum
//
// AffectedContracts returns all affected contracts that are the results of
// MessageCall transaction
func (evm *EVM[P]) AffectedContracts() []common.Address {
	addr := make([]common.Address, 0, len(evm.affectedContracts))
	for a, t := range evm.affectedContracts {
		if t == MessageCall {
			addr = append(addr, a)
		}
	}
	return addr[:]
}

// Quorum
//
// Return MerkleRoot of all affected contracts (due to both creation and message call)
func (evm *EVM[P]) CalculateMerkleRoot() (common.Hash, error) {
	combined := new(trie.Trie[P])
	for addr := range evm.affectedContracts {
		data, err := getDualState(evm, addr).GetRLPEncodedStateObject(addr)
		if err != nil {
			return common.Hash{}, err
		}
		if err := combined.TryUpdate(addr.Bytes(), data); err != nil {
			return common.Hash{}, err
		}
	}
	return combined.Hash(), nil
}
