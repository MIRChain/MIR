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
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

func TestState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	// Long tests:
	st.slow(`^stAttackTest/ContractCreationSpam`)
	st.slow(`^stBadOpcode/badOpcodes`)
	st.slow(`^stPreCompiledContracts/modexp`)
	st.slow(`^stQuadraticComplexityTest/`)
	st.slow(`^stStaticCall/static_Call50000`)
	st.slow(`^stStaticCall/static_Return50000`)
	st.slow(`^stSystemOperationsTest/CallRecursiveBomb`)
	st.slow(`^stTransactionTest/Opcodes_TransactionInit`)

	// Very time consuming
	st.skipLoad(`^stTimeConsuming/`)

	// Uses 1GB RAM per tested fork
	st.skipLoad(`^stStaticCall/static_Call1MB`)

	// Broken tests:
	// Expected failures:
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Byzantium/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/Constantinople/3`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/0`, "bug in test")
	//st.fails(`^stRevertTest/RevertPrecompiledTouch(_storage)?\.json/ConstantinopleFix/3`, "bug in test")

	// For Istanbul, older tests were moved into LegacyTests
	for _, dir := range []string{
		stateTestDir,
		legacyStateTestDir,
	} {
		st.walk(t, dir, func(t *testing.T, name string, test *StateTest[nist.PrivateKey, nist.PublicKey]) {
			for _, subtest := range test.Subtests() {
				subtest := subtest
				key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
				name := name + "/" + key

				t.Run(key+"/trie", func(t *testing.T) {
					withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config[nist.PublicKey]) error {
						_, _, err := test.Run(subtest, vmconfig, false)
						return st.checkFailure(t, name+"/trie", err)
					})
				})
				t.Run(key+"/snap", func(t *testing.T) {
					withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config[nist.PublicKey]) error {
						snaps, statedb, err := test.Run(subtest, vmconfig, true)
						if _, err := snaps.Journal(statedb.IntermediateRoot(false)); err != nil {
							return err
						}
						return st.checkFailure(t, name+"/snap", err)
					})
				})
			}
		})
	}
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withTrace[P crypto.PublicKey](t *testing.T, gasLimit uint64, test func(vm.Config[P]) error) {
	// Use config from command line arguments.
	config := vm.Config[P]{EVMInterpreter: *testEVM, EWASMInterpreter: *testEWASM}
	err := test(config)
	if err == nil {
		return
	}

	// Test failed, re-run with tracing enabled.
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	tracer := vm.NewJSONLogger[P](&vm.LogConfig{DisableMemory: true}, w)
	config.Debug, config.Tracer = true, tracer
	err2 := test(config)
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	w.Flush()
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
	//t.Logf("EVM output: 0x%x", tracer.Output())
	//t.Logf("EVM error: %v", tracer.Error())
}
