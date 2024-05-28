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

package vm

import (
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/params"
)

type (
	executionFunc[P crypto.PublicKey] func(pc *uint64, interpreter *EVMInterpreter[P], callContext *ScopeContext[P]) ([]byte, error)
	gasFunc[P crypto.PublicKey]       func(*EVM[P], *Contract[P], *Stack, *Memory, uint64) (uint64, error) // last parameter is the requested memory size as a uint64
	// memorySizeFunc returns the required size, and whether the operation overflowed a uint64
	memorySizeFunc func(*Stack) (size uint64, overflow bool)
)

type operation[P crypto.PublicKey] struct {
	// execute is the operation function
	execute     executionFunc[P]
	constantGas uint64
	dynamicGas  gasFunc[P]
	// minStack tells how many stack items are required
	minStack int
	// maxStack specifies the max length the stack can have for this operation
	// to not overflow the stack.
	maxStack int

	// memorySize returns the memory size required for the operation
	memorySize memorySizeFunc

	halts   bool // indicates whether the operation should halt further execution
	jumps   bool // indicates whether the program counter should not increment
	writes  bool // determines whether this a state modifying operation
	reverts bool // determines whether the operation reverts state (implicitly halts)
	returns bool // determines whether the operations sets the return data content
}

// var (
// 	frontierInstructionSet         = newFrontierInstructionSet[nist.PublicKey]()
// 	homesteadInstructionSet        = newHomesteadInstructionSet[nist.PublicKey]()
// 	tangerineWhistleInstructionSet = newTangerineWhistleInstructionSet[nist.PublicKey]()
// 	spuriousDragonInstructionSet   = newSpuriousDragonInstructionSet[nist.PublicKey]()
// 	byzantiumInstructionSet        = newByzantiumInstructionSet[nist.PublicKey]()
// 	constantinopleInstructionSet   = newConstantinopleInstructionSet[nist.PublicKey]()
// 	istanbulInstructionSet         = newIstanbulInstructionSet[nist.PublicKey]()
// 	berlinInstructionSet           = newBerlinInstructionSet[nist.PublicKey]()
// )

// JumpTable contains the EVM opcodes supported at a given fork.
type JumpTable[P crypto.PublicKey] [256]*operation[P]

// newBerlinInstructionSet returns the frontier, homestead, byzantium,
// contantinople, istanbul, petersburg and berlin instructions.
func newBerlinInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	instructionSet := newIstanbulInstructionSet[P]()
	enable2929(&instructionSet) // Access lists for trie accesses https://eips.ethereum.org/EIPS/eip-2929
	return instructionSet
}

// newIstanbulInstructionSet returns the frontier, homestead, byzantium,
// contantinople, istanbul and petersburg instructions.
func newIstanbulInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	instructionSet := newConstantinopleInstructionSet[P]()

	enable1344(&instructionSet) // ChainID opcode - https://eips.ethereum.org/EIPS/eip-1344
	enable1884(&instructionSet) // Reprice reader opcodes - https://eips.ethereum.org/EIPS/eip-1884
	enable2200(&instructionSet) // Net metered SSTORE - https://eips.ethereum.org/EIPS/eip-2200

	return instructionSet
}

// newConstantinopleInstructionSet returns the frontier, homestead,
// byzantium and contantinople instructions.
func newConstantinopleInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	instructionSet := newByzantiumInstructionSet[P]()
	instructionSet[SHL] = &operation[P]{
		execute:     opSHL[P],
		constantGas: GasFastestStep,
		minStack:    minStack(2, 1),
		maxStack:    maxStack(2, 1),
	}
	instructionSet[SHR] = &operation[P]{
		execute:     opSHR[P],
		constantGas: GasFastestStep,
		minStack:    minStack(2, 1),
		maxStack:    maxStack(2, 1),
	}
	instructionSet[SAR] = &operation[P]{
		execute:     opSAR[P],
		constantGas: GasFastestStep,
		minStack:    minStack(2, 1),
		maxStack:    maxStack(2, 1),
	}
	instructionSet[EXTCODEHASH] = &operation[P]{
		execute:     opExtCodeHash[P],
		constantGas: params.ExtcodeHashGasConstantinople,
		minStack:    minStack(1, 1),
		maxStack:    maxStack(1, 1),
	}
	instructionSet[CREATE2] = &operation[P]{
		execute:     opCreate2[P],
		constantGas: params.Create2Gas,
		dynamicGas:  gasCreate2[P],
		minStack:    minStack(4, 1),
		maxStack:    maxStack(4, 1),
		memorySize:  memoryCreate2,
		writes:      true,
		returns:     true,
	}
	return instructionSet
}

// newByzantiumInstructionSet returns the frontier, homestead and
// byzantium instructions.
func newByzantiumInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	instructionSet := newSpuriousDragonInstructionSet[P]()
	instructionSet[STATICCALL] = &operation[P]{
		execute:     opStaticCall[P],
		constantGas: params.CallGasEIP150,
		dynamicGas:  gasStaticCall[P],
		minStack:    minStack(6, 1),
		maxStack:    maxStack(6, 1),
		memorySize:  memoryStaticCall,
		returns:     true,
	}
	instructionSet[RETURNDATASIZE] = &operation[P]{
		execute:     opReturnDataSize[P],
		constantGas: GasQuickStep,
		minStack:    minStack(0, 1),
		maxStack:    maxStack(0, 1),
	}
	instructionSet[RETURNDATACOPY] = &operation[P]{
		execute:     opReturnDataCopy[P],
		constantGas: GasFastestStep,
		dynamicGas:  memoryCopierGas[P](2),
		minStack:    minStack(3, 0),
		maxStack:    maxStack(3, 0),
		memorySize:  memoryReturnDataCopy,
	}
	instructionSet[REVERT] = &operation[P]{
		execute:    opRevert[P],
		dynamicGas: pureMemoryGascost[P],
		minStack:   minStack(2, 0),
		maxStack:   maxStack(2, 0),
		memorySize: memoryRevert,
		reverts:    true,
		returns:    true,
	}
	return instructionSet
}

// EIP 158 a.k.a Spurious Dragon
func newSpuriousDragonInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	instructionSet := newTangerineWhistleInstructionSet[P]()
	instructionSet[EXP].dynamicGas = gasExpEIP158[P]
	return instructionSet

}

// EIP 150 a.k.a Tangerine Whistle
func newTangerineWhistleInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	instructionSet := newHomesteadInstructionSet[P]()
	instructionSet[BALANCE].constantGas = params.BalanceGasEIP150
	instructionSet[EXTCODESIZE].constantGas = params.ExtcodeSizeGasEIP150
	instructionSet[SLOAD].constantGas = params.SloadGasEIP150
	instructionSet[EXTCODECOPY].constantGas = params.ExtcodeCopyBaseEIP150
	instructionSet[CALL].constantGas = params.CallGasEIP150
	instructionSet[CALLCODE].constantGas = params.CallGasEIP150
	instructionSet[DELEGATECALL].constantGas = params.CallGasEIP150
	return instructionSet
}

// newHomesteadInstructionSet returns the frontier and homestead
// instructions that can be executed during the homestead phase.
func newHomesteadInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	instructionSet := newFrontierInstructionSet[P]()
	instructionSet[DELEGATECALL] = &operation[P]{
		execute:     opDelegateCall[P],
		dynamicGas:  gasDelegateCall[P],
		constantGas: params.CallGasFrontier,
		minStack:    minStack(6, 1),
		maxStack:    maxStack(6, 1),
		memorySize:  memoryDelegateCall,
		returns:     true,
	}
	return instructionSet
}

// newFrontierInstructionSet returns the frontier instructions
// that can be executed during the frontier phase.
func newFrontierInstructionSet[P crypto.PublicKey]() JumpTable[P] {
	return JumpTable[P]{
		STOP: {
			execute:     opStop[P],
			constantGas: 0,
			minStack:    minStack(0, 0),
			maxStack:    maxStack(0, 0),
			halts:       true,
		},
		ADD: {
			execute:     opAdd[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		MUL: {
			execute:     opMul[P],
			constantGas: GasFastStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		SUB: {
			execute:     opSub[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		DIV: {
			execute:     opDiv[P],
			constantGas: GasFastStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		SDIV: {
			execute:     opSdiv[P],
			constantGas: GasFastStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		MOD: {
			execute:     opMod[P],
			constantGas: GasFastStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		SMOD: {
			execute:     opSmod[P],
			constantGas: GasFastStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		ADDMOD: {
			execute:     opAddmod[P],
			constantGas: GasMidStep,
			minStack:    minStack(3, 1),
			maxStack:    maxStack(3, 1),
		},
		MULMOD: {
			execute:     opMulmod[P],
			constantGas: GasMidStep,
			minStack:    minStack(3, 1),
			maxStack:    maxStack(3, 1),
		},
		EXP: {
			execute:    opExp[P],
			dynamicGas: gasExpFrontier[P],
			minStack:   minStack(2, 1),
			maxStack:   maxStack(2, 1),
		},
		SIGNEXTEND: {
			execute:     opSignExtend[P],
			constantGas: GasFastStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		LT: {
			execute:     opLt[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		GT: {
			execute:     opGt[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		SLT: {
			execute:     opSlt[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		SGT: {
			execute:     opSgt[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		EQ: {
			execute:     opEq[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		ISZERO: {
			execute:     opIszero[P],
			constantGas: GasFastestStep,
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
		},
		AND: {
			execute:     opAnd[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		XOR: {
			execute:     opXor[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		OR: {
			execute:     opOr[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		NOT: {
			execute:     opNot[P],
			constantGas: GasFastestStep,
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
		},
		BYTE: {
			execute:     opByte[P],
			constantGas: GasFastestStep,
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
		},
		SHA3: {
			execute:     opSha3[P],
			constantGas: params.Sha3Gas,
			dynamicGas:  gasSha3[P],
			minStack:    minStack(2, 1),
			maxStack:    maxStack(2, 1),
			memorySize:  memorySha3,
		},
		ADDRESS: {
			execute:     opAddress[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		BALANCE: {
			execute:     opBalance[P],
			constantGas: params.BalanceGasFrontier,
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
		},
		ORIGIN: {
			execute:     opOrigin[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		CALLER: {
			execute:     opCaller[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		CALLVALUE: {
			execute:     opCallValue[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		CALLDATALOAD: {
			execute:     opCallDataLoad[P],
			constantGas: GasFastestStep,
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
		},
		CALLDATASIZE: {
			execute:     opCallDataSize[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		CALLDATACOPY: {
			execute:     opCallDataCopy[P],
			constantGas: GasFastestStep,
			dynamicGas:  memoryCopierGas[P](2),
			minStack:    minStack(3, 0),
			maxStack:    maxStack(3, 0),
			memorySize:  memoryCallDataCopy,
		},
		CODESIZE: {
			execute:     opCodeSize[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		CODECOPY: {
			execute:     opCodeCopy[P],
			constantGas: GasFastestStep,
			dynamicGas:  memoryCopierGas[P](2),
			minStack:    minStack(3, 0),
			maxStack:    maxStack(3, 0),
			memorySize:  memoryCodeCopy,
		},
		GASPRICE: {
			execute:     opGasprice[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		EXTCODESIZE: {
			execute:     opExtCodeSize[P],
			constantGas: params.ExtcodeSizeGasFrontier,
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
		},
		EXTCODECOPY: {
			execute:     opExtCodeCopy[P],
			constantGas: params.ExtcodeCopyBaseFrontier,
			dynamicGas:  memoryCopierGas[P](3),
			minStack:    minStack(4, 0),
			maxStack:    maxStack(4, 0),
			memorySize:  memoryExtCodeCopy,
		},
		BLOCKHASH: {
			execute:     opBlockhash[P],
			constantGas: GasExtStep,
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
		},
		COINBASE: {
			execute:     opCoinbase[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		TIMESTAMP: {
			execute:     opTimestamp[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		NUMBER: {
			execute:     opNumber[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		DIFFICULTY: {
			execute:     opDifficulty[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		GASLIMIT: {
			execute:     opGasLimit[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		POP: {
			execute:     opPop[P],
			constantGas: GasQuickStep,
			minStack:    minStack(1, 0),
			maxStack:    maxStack(1, 0),
		},
		MLOAD: {
			execute:     opMload[P],
			constantGas: GasFastestStep,
			dynamicGas:  pureMemoryGascost[P],
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
			memorySize:  memoryMLoad,
		},
		MSTORE: {
			execute:     opMstore[P],
			constantGas: GasFastestStep,
			dynamicGas:  pureMemoryGascost[P],
			minStack:    minStack(2, 0),
			maxStack:    maxStack(2, 0),
			memorySize:  memoryMStore,
		},
		MSTORE8: {
			execute:     opMstore8[P],
			constantGas: GasFastestStep,
			dynamicGas:  pureMemoryGascost[P],
			memorySize:  memoryMStore8,
			minStack:    minStack(2, 0),
			maxStack:    maxStack(2, 0),
		},
		SLOAD: {
			execute:     opSload[P],
			constantGas: params.SloadGasFrontier,
			minStack:    minStack(1, 1),
			maxStack:    maxStack(1, 1),
		},
		SSTORE: {
			execute:    opSstore[P],
			dynamicGas: gasSStore[P],
			minStack:   minStack(2, 0),
			maxStack:   maxStack(2, 0),
			writes:     true,
		},
		JUMP: {
			execute:     opJump[P],
			constantGas: GasMidStep,
			minStack:    minStack(1, 0),
			maxStack:    maxStack(1, 0),
			jumps:       true,
		},
		JUMPI: {
			execute:     opJumpi[P],
			constantGas: GasSlowStep,
			minStack:    minStack(2, 0),
			maxStack:    maxStack(2, 0),
			jumps:       true,
		},
		PC: {
			execute:     opPc[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		MSIZE: {
			execute:     opMsize[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		GAS: {
			execute:     opGas[P],
			constantGas: GasQuickStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		JUMPDEST: {
			execute:     opJumpdest[P],
			constantGas: params.JumpdestGas,
			minStack:    minStack(0, 0),
			maxStack:    maxStack(0, 0),
		},
		PUSH1: {
			execute:     opPush1[P],
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH2: {
			execute:     makePush[P](2, 2),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH3: {
			execute:     makePush[P](3, 3),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH4: {
			execute:     makePush[P](4, 4),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH5: {
			execute:     makePush[P](5, 5),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH6: {
			execute:     makePush[P](6, 6),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH7: {
			execute:     makePush[P](7, 7),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH8: {
			execute:     makePush[P](8, 8),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH9: {
			execute:     makePush[P](9, 9),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH10: {
			execute:     makePush[P](10, 10),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH11: {
			execute:     makePush[P](11, 11),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH12: {
			execute:     makePush[P](12, 12),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH13: {
			execute:     makePush[P](13, 13),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH14: {
			execute:     makePush[P](14, 14),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH15: {
			execute:     makePush[P](15, 15),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH16: {
			execute:     makePush[P](16, 16),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH17: {
			execute:     makePush[P](17, 17),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH18: {
			execute:     makePush[P](18, 18),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH19: {
			execute:     makePush[P](19, 19),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH20: {
			execute:     makePush[P](20, 20),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH21: {
			execute:     makePush[P](21, 21),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH22: {
			execute:     makePush[P](22, 22),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH23: {
			execute:     makePush[P](23, 23),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH24: {
			execute:     makePush[P](24, 24),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH25: {
			execute:     makePush[P](25, 25),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH26: {
			execute:     makePush[P](26, 26),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH27: {
			execute:     makePush[P](27, 27),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH28: {
			execute:     makePush[P](28, 28),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH29: {
			execute:     makePush[P](29, 29),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH30: {
			execute:     makePush[P](30, 30),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH31: {
			execute:     makePush[P](31, 31),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		PUSH32: {
			execute:     makePush[P](32, 32),
			constantGas: GasFastestStep,
			minStack:    minStack(0, 1),
			maxStack:    maxStack(0, 1),
		},
		DUP1: {
			execute:     makeDup[P](1),
			constantGas: GasFastestStep,
			minStack:    minDupStack(1),
			maxStack:    maxDupStack(1),
		},
		DUP2: {
			execute:     makeDup[P](2),
			constantGas: GasFastestStep,
			minStack:    minDupStack(2),
			maxStack:    maxDupStack(2),
		},
		DUP3: {
			execute:     makeDup[P](3),
			constantGas: GasFastestStep,
			minStack:    minDupStack(3),
			maxStack:    maxDupStack(3),
		},
		DUP4: {
			execute:     makeDup[P](4),
			constantGas: GasFastestStep,
			minStack:    minDupStack(4),
			maxStack:    maxDupStack(4),
		},
		DUP5: {
			execute:     makeDup[P](5),
			constantGas: GasFastestStep,
			minStack:    minDupStack(5),
			maxStack:    maxDupStack(5),
		},
		DUP6: {
			execute:     makeDup[P](6),
			constantGas: GasFastestStep,
			minStack:    minDupStack(6),
			maxStack:    maxDupStack(6),
		},
		DUP7: {
			execute:     makeDup[P](7),
			constantGas: GasFastestStep,
			minStack:    minDupStack(7),
			maxStack:    maxDupStack(7),
		},
		DUP8: {
			execute:     makeDup[P](8),
			constantGas: GasFastestStep,
			minStack:    minDupStack(8),
			maxStack:    maxDupStack(8),
		},
		DUP9: {
			execute:     makeDup[P](9),
			constantGas: GasFastestStep,
			minStack:    minDupStack(9),
			maxStack:    maxDupStack(9),
		},
		DUP10: {
			execute:     makeDup[P](10),
			constantGas: GasFastestStep,
			minStack:    minDupStack(10),
			maxStack:    maxDupStack(10),
		},
		DUP11: {
			execute:     makeDup[P](11),
			constantGas: GasFastestStep,
			minStack:    minDupStack(11),
			maxStack:    maxDupStack(11),
		},
		DUP12: {
			execute:     makeDup[P](12),
			constantGas: GasFastestStep,
			minStack:    minDupStack(12),
			maxStack:    maxDupStack(12),
		},
		DUP13: {
			execute:     makeDup[P](13),
			constantGas: GasFastestStep,
			minStack:    minDupStack(13),
			maxStack:    maxDupStack(13),
		},
		DUP14: {
			execute:     makeDup[P](14),
			constantGas: GasFastestStep,
			minStack:    minDupStack(14),
			maxStack:    maxDupStack(14),
		},
		DUP15: {
			execute:     makeDup[P](15),
			constantGas: GasFastestStep,
			minStack:    minDupStack(15),
			maxStack:    maxDupStack(15),
		},
		DUP16: {
			execute:     makeDup[P](16),
			constantGas: GasFastestStep,
			minStack:    minDupStack(16),
			maxStack:    maxDupStack(16),
		},
		SWAP1: {
			execute:     makeSwap[P](1),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(2),
			maxStack:    maxSwapStack(2),
		},
		SWAP2: {
			execute:     makeSwap[P](2),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(3),
			maxStack:    maxSwapStack(3),
		},
		SWAP3: {
			execute:     makeSwap[P](3),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(4),
			maxStack:    maxSwapStack(4),
		},
		SWAP4: {
			execute:     makeSwap[P](4),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(5),
			maxStack:    maxSwapStack(5),
		},
		SWAP5: {
			execute:     makeSwap[P](5),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(6),
			maxStack:    maxSwapStack(6),
		},
		SWAP6: {
			execute:     makeSwap[P](6),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(7),
			maxStack:    maxSwapStack(7),
		},
		SWAP7: {
			execute:     makeSwap[P](7),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(8),
			maxStack:    maxSwapStack(8),
		},
		SWAP8: {
			execute:     makeSwap[P](8),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(9),
			maxStack:    maxSwapStack(9),
		},
		SWAP9: {
			execute:     makeSwap[P](9),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(10),
			maxStack:    maxSwapStack(10),
		},
		SWAP10: {
			execute:     makeSwap[P](10),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(11),
			maxStack:    maxSwapStack(11),
		},
		SWAP11: {
			execute:     makeSwap[P](11),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(12),
			maxStack:    maxSwapStack(12),
		},
		SWAP12: {
			execute:     makeSwap[P](12),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(13),
			maxStack:    maxSwapStack(13),
		},
		SWAP13: {
			execute:     makeSwap[P](13),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(14),
			maxStack:    maxSwapStack(14),
		},
		SWAP14: {
			execute:     makeSwap[P](14),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(15),
			maxStack:    maxSwapStack(15),
		},
		SWAP15: {
			execute:     makeSwap[P](15),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(16),
			maxStack:    maxSwapStack(16),
		},
		SWAP16: {
			execute:     makeSwap[P](16),
			constantGas: GasFastestStep,
			minStack:    minSwapStack(17),
			maxStack:    maxSwapStack(17),
		},
		LOG0: {
			execute:    makeLog[P](0),
			dynamicGas: makeGasLog[P](0),
			minStack:   minStack(2, 0),
			maxStack:   maxStack(2, 0),
			memorySize: memoryLog,
			writes:     true,
		},
		LOG1: {
			execute:    makeLog[P](1),
			dynamicGas: makeGasLog[P](1),
			minStack:   minStack(3, 0),
			maxStack:   maxStack(3, 0),
			memorySize: memoryLog,
			writes:     true,
		},
		LOG2: {
			execute:    makeLog[P](2),
			dynamicGas: makeGasLog[P](2),
			minStack:   minStack(4, 0),
			maxStack:   maxStack(4, 0),
			memorySize: memoryLog,
			writes:     true,
		},
		LOG3: {
			execute:    makeLog[P](3),
			dynamicGas: makeGasLog[P](3),
			minStack:   minStack(5, 0),
			maxStack:   maxStack(5, 0),
			memorySize: memoryLog,
			writes:     true,
		},
		LOG4: {
			execute:    makeLog[P](4),
			dynamicGas: makeGasLog[P](4),
			minStack:   minStack(6, 0),
			maxStack:   maxStack(6, 0),
			memorySize: memoryLog,
			writes:     true,
		},
		CREATE: {
			execute:     opCreate[P],
			constantGas: params.CreateGas,
			dynamicGas:  pureMemoryGascost[P],
			minStack:    minStack(3, 1),
			maxStack:    maxStack(3, 1),
			memorySize:  memoryCreate,
			writes:      true,
			returns:     true,
		},
		CALL: {
			execute:     opCall[P],
			constantGas: params.CallGasFrontier,
			dynamicGas:  gasCall[P],
			minStack:    minStack(7, 1),
			maxStack:    maxStack(7, 1),
			memorySize:  memoryCall,
			returns:     true,
		},
		CALLCODE: {
			execute:     opCallCode[P],
			constantGas: params.CallGasFrontier,
			dynamicGas:  gasCallCode[P],
			minStack:    minStack(7, 1),
			maxStack:    maxStack(7, 1),
			memorySize:  memoryCall,
			returns:     true,
		},
		RETURN: {
			execute:    opReturn[P],
			dynamicGas: pureMemoryGascost[P],
			minStack:   minStack(2, 0),
			maxStack:   maxStack(2, 0),
			memorySize: memoryReturn,
			halts:      true,
		},
		SELFDESTRUCT: {
			execute:    opSuicide[P],
			dynamicGas: gasSelfdestruct[P],
			minStack:   minStack(1, 0),
			maxStack:   maxStack(1, 0),
			halts:      true,
			writes:     true,
		},
	}
}
