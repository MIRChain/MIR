// Copyright 2019 The go-ethereum Authors
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
	"fmt"
	"sort"

	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/params"
	"github.com/holiman/uint256"
)

// var activators = map[int]func(*JumpTable[nist.PublicKey]){
// 	2929: enable2929[nist.PublicKey],
// 	2200: enable2200[nist.PublicKey],
// 	1884: enable1884[nist.PublicKey],
// 	1344: enable1344[nist.PublicKey],
// }

// EnableEIP enables the given EIP on the config.
// This operation writes in-place, and callers need to ensure that the globally
// defined jump tables are not polluted.
func EnableEIP[P crypto.PublicKey](eipNum int, jt *JumpTable[P]) error {
	enablerFn, ok := map[int]func(*JumpTable[P]){
		2929: enable2929[P],
		2200: enable2200[P],
		1884: enable1884[P],
		1344: enable1344[P],
	}[eipNum]
	if !ok {
		return fmt.Errorf("undefined eip %d", eipNum)
	}
	enablerFn(jt)
	return nil
}

func ValidEip[P crypto.PublicKey](eipNum int) bool {
	_, ok := map[int]func(*JumpTable[P]){
		2929: enable2929[P],
		2200: enable2200[P],
		1884: enable1884[P],
		1344: enable1344[P],
	}[eipNum]
	return ok
}
func ActivateableEips[P crypto.PublicKey]() []string {
	var nums []string
	for k := range map[int]func(*JumpTable[P]){
		2929: enable2929[P],
		2200: enable2200[P],
		1884: enable1884[P],
		1344: enable1344[P],
	} {
		nums = append(nums, fmt.Sprintf("%d", k))
	}
	sort.Strings(nums)
	return nums
}

// enable1884 applies EIP-1884 to the given jump table:
// - Increase cost of BALANCE to 700
// - Increase cost of EXTCODEHASH to 700
// - Increase cost of SLOAD to 800
// - Define SELFBALANCE, with cost GasFastStep (5)
func enable1884[P crypto.PublicKey](jt *JumpTable[P]) {
	// Gas cost changes
	jt[SLOAD].constantGas = params.SloadGasEIP1884
	jt[BALANCE].constantGas = params.BalanceGasEIP1884
	jt[EXTCODEHASH].constantGas = params.ExtcodeHashGasEIP1884

	// New opcode
	jt[SELFBALANCE] = &operation[P]{
		execute:     opSelfBalance[P],
		constantGas: GasFastStep,
		minStack:    minStack(0, 1),
		maxStack:    maxStack(0, 1),
	}
}

func opSelfBalance[P crypto.PublicKey](pc *uint64, interpreter *EVMInterpreter[P], scope *ScopeContext[P]) ([]byte, error) {
	balance, _ := uint256.FromBig(interpreter.evm.StateDB.GetBalance(scope.Contract.Address()))
	scope.Stack.push(balance)
	return nil, nil
}

// enable1344 applies EIP-1344 (ChainID Opcode)
// - Adds an opcode that returns the current chain’s EIP-155 unique identifier
func enable1344[P crypto.PublicKey](jt *JumpTable[P]) {
	// New opcode
	jt[CHAINID] = &operation[P]{
		execute:     opChainID[P],
		constantGas: GasQuickStep,
		minStack:    minStack(0, 1),
		maxStack:    maxStack(0, 1),
	}
}

// opChainID implements CHAINID opcode
func opChainID[P crypto.PublicKey](pc *uint64, interpreter *EVMInterpreter[P], scope *ScopeContext[P]) ([]byte, error) {
	chainId, _ := uint256.FromBig(interpreter.evm.chainConfig.ChainID)
	scope.Stack.push(chainId)
	return nil, nil
}

// enable2200 applies EIP-2200 (Rebalance net-metered SSTORE)
func enable2200[P crypto.PublicKey](jt *JumpTable[P]) {
	jt[SLOAD].constantGas = params.SloadGasEIP2200
	jt[SSTORE].dynamicGas = gasSStoreEIP2200[P]
}

// enable2929 enables "EIP-2929: Gas cost increases for state access opcodes"
// https://eips.ethereum.org/EIPS/eip-2929
func enable2929[P crypto.PublicKey](jt *JumpTable[P]) {
	jt[SSTORE].dynamicGas = gasSStoreEIP2929[P]

	jt[SLOAD].constantGas = 0
	jt[SLOAD].dynamicGas = gasSLoadEIP2929[P]

	jt[EXTCODECOPY].constantGas = WarmStorageReadCostEIP2929
	jt[EXTCODECOPY].dynamicGas = gasExtCodeCopyEIP2929[P]

	jt[EXTCODESIZE].constantGas = WarmStorageReadCostEIP2929
	jt[EXTCODESIZE].dynamicGas = gasEip2929AccountCheck[P]

	jt[EXTCODEHASH].constantGas = WarmStorageReadCostEIP2929
	jt[EXTCODEHASH].dynamicGas = gasEip2929AccountCheck[P]

	jt[BALANCE].constantGas = WarmStorageReadCostEIP2929
	jt[BALANCE].dynamicGas = gasEip2929AccountCheck[P]

	jt[CALL].constantGas = WarmStorageReadCostEIP2929
	jt[CALL].dynamicGas = makeCallVariantGasCallEIP2929(gasCall[P])

	jt[CALLCODE].constantGas = WarmStorageReadCostEIP2929
	jt[CALLCODE].dynamicGas = makeCallVariantGasCallEIP2929(gasCallCode[P])

	jt[STATICCALL].constantGas = WarmStorageReadCostEIP2929
	jt[STATICCALL].dynamicGas = makeCallVariantGasCallEIP2929(gasStaticCall[P])

	jt[DELEGATECALL].constantGas = WarmStorageReadCostEIP2929
	jt[DELEGATECALL].dynamicGas = makeCallVariantGasCallEIP2929(gasDelegateCall[P])

	// This was previously part of the dynamic cost, but we're using it as a constantGas
	// factor here
	jt[SELFDESTRUCT].constantGas = params.SelfdestructGasEIP150
	jt[SELFDESTRUCT].dynamicGas = gasSelfdestructEIP2929[P]
}
