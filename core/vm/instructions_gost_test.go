// Copyright 2017 The go-ethereum Authors
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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/params"
	"github.com/holiman/uint256"
)

var twoOpMethodsGost map[string]executionFunc[gost3410.PublicKey]

func init() {

	// Params is a list of common edgecases that should be used for some common tests
	params := []string{
		"0000000000000000000000000000000000000000000000000000000000000000", // 0
		"0000000000000000000000000000000000000000000000000000000000000001", // +1
		"0000000000000000000000000000000000000000000000000000000000000005", // +5
		"7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", // + max -1
		"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", // + max
		"8000000000000000000000000000000000000000000000000000000000000000", // - max
		"8000000000000000000000000000000000000000000000000000000000000001", // - max+1
		"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffb", // - 5
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", // - 1
	}
	// Params are combined so each param is used on each 'side'
	commonParams = make([]*twoOperandParams, len(params)*len(params))
	for i, x := range params {
		for j, y := range params {
			commonParams[i*len(params)+j] = &twoOperandParams{x, y}
		}
	}
	twoOpMethodsGost = map[string]executionFunc[gost3410.PublicKey]{
		"add":     opAdd[gost3410.PublicKey],
		"sub":     opSub[gost3410.PublicKey],
		"mul":     opMul[gost3410.PublicKey],
		"div":     opDiv[gost3410.PublicKey],
		"sdiv":    opSdiv[gost3410.PublicKey],
		"mod":     opMod[gost3410.PublicKey],
		"smod":    opSmod[gost3410.PublicKey],
		"exp":     opExp[gost3410.PublicKey],
		"signext": opSignExtend[gost3410.PublicKey],
		"lt":      opLt[gost3410.PublicKey],
		"gt":      opGt[gost3410.PublicKey],
		"slt":     opSlt[gost3410.PublicKey],
		"sgt":     opSgt[gost3410.PublicKey],
		"eq":      opEq[gost3410.PublicKey],
		"and":     opAnd[gost3410.PublicKey],
		"or":      opOr[gost3410.PublicKey],
		"xor":     opXor[gost3410.PublicKey],
		"byte":    opByte[gost3410.PublicKey],
		"shl":     opSHL[gost3410.PublicKey],
		"shr":     opSHR[gost3410.PublicKey],
		"sar":     opSAR[gost3410.PublicKey],
	}
}

func TestByteOpGost(t *testing.T) {
	tests := []TwoOperandTestcase{
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", "00", "AB"},
		{"ABCDEF0908070605040302010000000000000000000000000000000000000000", "01", "CD"},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", "00", "00"},
		{"00CDEF090807060504030201ffffffffffffffffffffffffffffffffffffffff", "01", "CD"},
		{"0000000000000000000000000000000000000000000000000000000000102030", "1F", "30"},
		{"0000000000000000000000000000000000000000000000000000000000102030", "1E", "20"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "20", "00"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "FFFFFFFFFFFFFFFF", "00"},
	}
	testTwoOperandOp(t, tests, opByte[gost3410.PublicKey], "byte")
}

func TestSHLGost(t *testing.T) {
	// Testcases from https://github.com/ethereum/EIPs/blob/master/EIPS/eip-145.md#shl-shift-left
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000002"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "ff", "8000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "0101", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "8000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"},
	}
	testTwoOperandOp(t, tests, opSHL[gost3410.PublicKey], "shl")
}

func TestSHRGost(t *testing.T) {
	// Testcases from https://github.com/ethereum/EIPs/blob/master/EIPS/eip-145.md#shr-logical-shift-right
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "00", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "01", "4000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "ff", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0101", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
	}
	testTwoOperandOp(t, tests, opSHR[gost3410.PublicKey], "shr")
}

func TestSARGost(t *testing.T) {
	// Testcases from https://github.com/ethereum/EIPs/blob/master/EIPS/eip-145.md#sar-arithmetic-shift-right
	tests := []TwoOperandTestcase{
		{"0000000000000000000000000000000000000000000000000000000000000001", "00", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"0000000000000000000000000000000000000000000000000000000000000001", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "01", "c000000000000000000000000000000000000000000000000000000000000000"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "ff", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0100", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"8000000000000000000000000000000000000000000000000000000000000000", "0101", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "00", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "01", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
		{"0000000000000000000000000000000000000000000000000000000000000000", "01", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"4000000000000000000000000000000000000000000000000000000000000000", "fe", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "f8", "000000000000000000000000000000000000000000000000000000000000007f"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "fe", "0000000000000000000000000000000000000000000000000000000000000001"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "ff", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", "0100", "0000000000000000000000000000000000000000000000000000000000000000"},
	}

	testTwoOperandOp(t, tests, opSAR[gost3410.PublicKey], "sar")
}

func TestAddModGost(t *testing.T) {
	var (
		env            = NewEVM(BlockContext{}, TxContext{}, nil, nil, params.TestChainConfig, Config[gost3410.PublicKey]{})
		stack          = newstack()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
		pc             = uint64(0)
	)
	tests := []struct {
		x        string
		y        string
		z        string
		expected string
	}{
		{"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe",
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe",
		},
	}
	// x + y = 0x1fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd
	// in 256 bit repr, fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd

	for i, test := range tests {
		x := new(uint256.Int).SetBytes(common.Hex2Bytes(test.x))
		y := new(uint256.Int).SetBytes(common.Hex2Bytes(test.y))
		z := new(uint256.Int).SetBytes(common.Hex2Bytes(test.z))
		expected := new(uint256.Int).SetBytes(common.Hex2Bytes(test.expected))
		stack.push(z)
		stack.push(y)
		stack.push(x)
		opAddmod[gost3410.PublicKey](&pc, evmInterpreter, &ScopeContext[gost3410.PublicKey]{nil, stack, nil})
		actual := stack.pop()
		if actual.Cmp(expected) != 0 {
			t.Errorf("Testcase %d, expected  %x, got %x", i, expected, actual)
		}
	}
}

// utility function to fill the json-file with testcases
// Enable this test to generate the 'testcases_xx.json' files
func TestWriteExpectedValuesGost(t *testing.T) {
	t.Skip("Enable this test to create json test cases.")

	for name, method := range twoOpMethodsGost {
		data, err := json.Marshal(getResult(commonParams, method))
		if err != nil {
			t.Fatal(err)
		}
		_ = ioutil.WriteFile(fmt.Sprintf("testdata/testcases_%v.json", name), data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestJsonTestcases runs through all the testcases defined as json-files
func TestJsonTestcasesGost(t *testing.T) {
	for name := range twoOpMethodsGost {
		data, err := ioutil.ReadFile(fmt.Sprintf("testdata/testcases_%v.json", name))
		if err != nil {
			t.Fatal("Failed to read file", err)
		}
		var testcases []TwoOperandTestcase
		json.Unmarshal(data, &testcases)
		testTwoOperandOp(t, testcases, twoOpMethodsGost[name], name)
	}
}

func BenchmarkOpAdd64Gost(b *testing.B) {
	x := "ffffffff"
	y := "fd37f3e2bba2c4f"

	opBenchmark(b, opAdd[gost3410.PublicKey], x, y)
}

func BenchmarkOpAdd128Gost(b *testing.B) {
	x := "ffffffffffffffff"
	y := "f5470b43c6549b016288e9a65629687"

	opBenchmark(b, opAdd[gost3410.PublicKey], x, y)
}

func BenchmarkOpAdd256Gost(b *testing.B) {
	x := "0802431afcbce1fc194c9eaa417b2fb67dc75a95db0bc7ec6b1c8af11df6a1da9"
	y := "a1f5aac137876480252e5dcac62c354ec0d42b76b0642b6181ed099849ea1d57"

	opBenchmark(b, opAdd[gost3410.PublicKey], x, y)
}

func BenchmarkOpSub64Gost(b *testing.B) {
	x := "51022b6317003a9d"
	y := "a20456c62e00753a"

	opBenchmark(b, opSub[gost3410.PublicKey], x, y)
}

func BenchmarkOpSub128Gost(b *testing.B) {
	x := "4dde30faaacdc14d00327aac314e915d"
	y := "9bbc61f5559b829a0064f558629d22ba"

	opBenchmark(b, opSub[gost3410.PublicKey], x, y)
}

func BenchmarkOpSub256Gost(b *testing.B) {
	x := "4bfcd8bb2ac462735b48a17580690283980aa2d679f091c64364594df113ea37"
	y := "97f9b1765588c4e6b69142eb00d20507301545acf3e1238c86c8b29be227d46e"

	opBenchmark(b, opSub[gost3410.PublicKey], x, y)
}

func BenchmarkOpMulGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opMul[gost3410.PublicKey], x, y)
}

func BenchmarkOpDiv256Gost(b *testing.B) {
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"
	opBenchmark(b, opDiv[gost3410.PublicKey], x, y)
}

func BenchmarkOpDiv128Gost(b *testing.B) {
	x := "fdedc7f10142ff97"
	y := "fbdfda0e2ce356173d1993d5f70a2b11"
	opBenchmark(b, opDiv[gost3410.PublicKey], x, y)
}

func BenchmarkOpDiv64Gost(b *testing.B) {
	x := "fcb34eb3"
	y := "f97180878e839129"
	opBenchmark(b, opDiv[gost3410.PublicKey], x, y)
}

func BenchmarkOpSdivGost(b *testing.B) {
	x := "ff3f9014f20db29ae04af2c2d265de17"
	y := "fe7fb0d1f59dfe9492ffbf73683fd1e870eec79504c60144cc7f5fc2bad1e611"

	opBenchmark(b, opSdiv[gost3410.PublicKey], x, y)
}

func BenchmarkOpModGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opMod[gost3410.PublicKey], x, y)
}

func BenchmarkOpSmodGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opSmod[gost3410.PublicKey], x, y)
}

func BenchmarkOpExpGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opExp[gost3410.PublicKey], x, y)
}

func BenchmarkOpSignExtendGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opSignExtend[gost3410.PublicKey], x, y)
}

func BenchmarkOpLtGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opLt[gost3410.PublicKey], x, y)
}

func BenchmarkOpGtGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opGt[gost3410.PublicKey], x, y)
}

func BenchmarkOpSltGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opSlt[gost3410.PublicKey], x, y)
}

func BenchmarkOpSgtGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opSgt[gost3410.PublicKey], x, y)
}

func BenchmarkOpEqGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opEq[gost3410.PublicKey], x, y)
}
func BenchmarkOpEq2Gost(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201fffffffe"
	opBenchmark(b, opEq[gost3410.PublicKey], x, y)
}
func BenchmarkOpAndGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opAnd[gost3410.PublicKey], x, y)
}

func BenchmarkOpOrGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opOr[gost3410.PublicKey], x, y)
}

func BenchmarkOpXorGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opXor[gost3410.PublicKey], x, y)
}

func BenchmarkOpByteGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup

	opBenchmark(b, opByte[gost3410.PublicKey], x, y)
}

func BenchmarkOpAddmodGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup
	z := alphabetSoup

	opBenchmark(b, opAddmod[gost3410.PublicKey], x, y, z)
}

func BenchmarkOpMulmodGost(b *testing.B) {
	x := alphabetSoup
	y := alphabetSoup
	z := alphabetSoup

	opBenchmark(b, opMulmod[gost3410.PublicKey], x, y, z)
}

func BenchmarkOpSHLGost(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSHL[gost3410.PublicKey], x, y)
}
func BenchmarkOpSHRGost(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSHR[gost3410.PublicKey], x, y)
}
func BenchmarkOpSARGost(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	y := "ff"

	opBenchmark(b, opSAR[gost3410.PublicKey], x, y)
}
func BenchmarkOpIsZeroGost(b *testing.B) {
	x := "FBCDEF090807060504030201ffffffffFBCDEF090807060504030201ffffffff"
	opBenchmark(b, opIszero[gost3410.PublicKey], x)
}

func TestOpMstoreGost(t *testing.T) {
	var (
		env            = NewEVM(BlockContext{}, TxContext{}, nil, nil, params.TestChainConfig, Config[gost3410.PublicKey]{})
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	mem.Resize(64)
	pc := uint64(0)
	v := "abcdef00000000000000abba000000000deaf000000c0de00100000000133700"
	stack.pushN(*new(uint256.Int).SetBytes(common.Hex2Bytes(v)), *new(uint256.Int))
	opMstore(&pc, evmInterpreter, &ScopeContext[gost3410.PublicKey]{mem, stack, nil})
	if got := common.Bytes2Hex(mem.GetCopy(0, 32)); got != v {
		t.Fatalf("Mstore fail, got %v, expected %v", got, v)
	}
	stack.pushN(*new(uint256.Int).SetUint64(0x1), *new(uint256.Int))
	opMstore(&pc, evmInterpreter, &ScopeContext[gost3410.PublicKey]{mem, stack, nil})
	if common.Bytes2Hex(mem.GetCopy(0, 32)) != "0000000000000000000000000000000000000000000000000000000000000001" {
		t.Fatalf("Mstore failed to overwrite previous value")
	}
}

func BenchmarkOpMstoreGost(bench *testing.B) {
	var (
		env            = NewEVM(BlockContext{}, TxContext{}, nil, nil, params.TestChainConfig, Config[gost3410.PublicKey]{})
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)

	env.interpreter = evmInterpreter
	mem.Resize(64)
	pc := uint64(0)
	memStart := new(uint256.Int)
	value := new(uint256.Int).SetUint64(0x1337)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		stack.pushN(*value, *memStart)
		opMstore(&pc, evmInterpreter, &ScopeContext[gost3410.PublicKey]{mem, stack, nil})
	}
}

func BenchmarkOpSHA3Gost(bench *testing.B) {
	var (
		env            = NewEVM(BlockContext{}, TxContext{}, nil, nil, params.TestChainConfig, Config[gost3410.PublicKey]{})
		stack          = newstack()
		mem            = NewMemory()
		evmInterpreter = NewEVMInterpreter(env, env.vmConfig)
	)
	env.interpreter = evmInterpreter
	mem.Resize(32)
	pc := uint64(0)
	start := new(uint256.Int)

	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		stack.pushN(*uint256.NewInt(32), *start)
		opSha3(&pc, evmInterpreter, &ScopeContext[gost3410.PublicKey]{mem, stack, nil})
	}
}

func TestCreate2AddresesGost(t *testing.T) {
	type testcase struct {
		origin   string
		salt     string
		code     string
		expected string
	}

	for i, tt := range []testcase{
		{
			origin:   "0x0000000000000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "0x964Bd211B963BfFfb9Ac463bE86186d31CF62116",
		},
		{
			origin:   "0xdeadbeef00000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x00",
			expected: "0x200fFcf9ef5E62D186dBc4e977167B47E0168783",
		},
		{
			origin:   "0xdeadbeef00000000000000000000000000000000",
			salt:     "0xfeed000000000000000000000000000000000000",
			code:     "0x00",
			expected: "0x758AEafdC95f54E3201aCd880ed72ed8e912BCf8",
		},
		{
			origin:   "0x0000000000000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0xdeadbeef",
			expected: "0x3f0DF5eB7828e7a4628F5ffD37F866eB110059BC",
		},
		{
			origin:   "0x00000000000000000000000000000000deadbeef",
			salt:     "0xcafebabe",
			code:     "0xdeadbeef",
			expected: "0xc243C09c247CA5f57BFD1D1139E8ef35e0f82c80",
		},
		{
			origin:   "0x00000000000000000000000000000000deadbeef",
			salt:     "0xcafebabe",
			code:     "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			expected: "0x95580d2c06633d5aa4B6E95a21b98b6d873Bc3c3",
		},
		{
			origin:   "0x0000000000000000000000000000000000000000",
			salt:     "0x0000000000000000000000000000000000000000",
			code:     "0x",
			expected: "0xc2DD95E7Ac5EA43Feb70914060DE4F98A7b3c973",
		},
	} {

		origin := common.BytesToAddress(common.FromHex(tt.origin))
		salt := common.BytesToHash(common.FromHex(tt.salt))
		code := common.FromHex(tt.code)
		codeHash := crypto.Keccak256[gost3410.PublicKey](code)
		address := crypto.CreateAddress2[gost3410.PublicKey](origin, salt, codeHash)
		/*
			stack          := newstack()
			// salt, but we don't need that for this test
			stack.push(big.NewInt(int64(len(code)))) //size
			stack.push(big.NewInt(0)) // memstart
			stack.push(big.NewInt(0)) // value
			gas, _ := gasCreate2(params.GasTable{}, nil, nil, stack, nil, 0)
			fmt.Printf("Example %d\n* address `0x%x`\n* salt `0x%x`\n* init_code `0x%x`\n* gas (assuming no mem expansion): `%v`\n* result: `%s`\n\n", i,origin, salt, code, gas, address.String())
		*/
		expected := common.BytesToAddress(common.FromHex(tt.expected))
		if !bytes.Equal(expected.Bytes(), address.Bytes()) {
			t.Errorf("test %d: expected %s, got %s", i, expected.String(), address.String())
		}
	}
}
