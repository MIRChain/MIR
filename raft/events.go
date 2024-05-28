package raft

import (
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
)

type InvalidRaftOrdering[P crypto.PublicKey] struct {
	// Current head of the chain
	headBlock *types.Block[P]

	// New block that should point to the head, but doesn't
	invalidBlock *types.Block[P]
}
