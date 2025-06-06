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
	"encoding/json"
	"io"
	"math/big"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type JSONLogger [P crypto.PublicKey] struct {
	encoder *json.Encoder
	cfg     *LogConfig
}

// var _ Tracer = &JSONLogger{}

// NewJSONLogger creates a new EVM tracer that prints execution steps as JSON objects
// into the provided stream.
func NewJSONLogger[P crypto.PublicKey](cfg *LogConfig, writer io.Writer) *JSONLogger[P] {
	l := &JSONLogger[P]{json.NewEncoder(writer), cfg}
	if l.cfg == nil {
		l.cfg = &LogConfig{}
	}
	return l
}

func (l *JSONLogger[P]) CaptureStart(env *EVM[P], from, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
}

func (l *JSONLogger[P]) CaptureFault(*EVM[P], uint64, OpCode, uint64, uint64, *ScopeContext[P], int, error) {}

// CaptureState outputs state information on the logger.
func (l *JSONLogger[P]) CaptureState(env *EVM[P], pc uint64, op OpCode, gas, cost uint64, scope *ScopeContext[P], rData []byte, depth int, err error) {
	memory := scope.Memory
	stack := scope.Stack

	log := StructLog{
		Pc:            pc,
		Op:            op,
		Gas:           gas,
		GasCost:       cost,
		MemorySize:    memory.Len(),
		Storage:       nil,
		Depth:         depth,
		RefundCounter: env.StateDB.GetRefund(),
		Err:           err,
	}
	if !l.cfg.DisableMemory {
		log.Memory = memory.Data()
	}
	if !l.cfg.DisableStack {
		//TODO(@holiman) improve this
		logstack := make([]*big.Int, len(stack.Data()))
		for i, item := range stack.Data() {
			logstack[i] = item.ToBig()
		}
		log.Stack = logstack
	}
	if !l.cfg.DisableReturnData {
		log.ReturnData = rData
	}
	l.encoder.Encode(log)
}

// CaptureEnd is triggered at end of execution.
func (l *JSONLogger[P]) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	type endLog struct {
		Output  string              `json:"output"`
		GasUsed math.HexOrDecimal64 `json:"gasUsed"`
		Time    time.Duration       `json:"time"`
		Err     string              `json:"error,omitempty"`
	}
	if err != nil {
		l.encoder.Encode(endLog{common.Bytes2Hex(output), math.HexOrDecimal64(gasUsed), t, err.Error()})
	}
	l.encoder.Encode(endLog{common.Bytes2Hex(output), math.HexOrDecimal64(gasUsed), t, ""})
}
