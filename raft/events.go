package raft

import (
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type InvalidRaftOrdering [P crypto.PublicKey] struct {
	// Current head of the chain
	headBlock *types.Block[P]

	// New block that should point to the head, but doesn't
	invalidBlock *types.Block[P]
}
