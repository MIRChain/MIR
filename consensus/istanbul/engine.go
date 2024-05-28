package istanbul

import (
	"math/big"
	"time"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/consensus"
	"github.com/MIRChain/MIR/core/state"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
)

type Engine[P crypto.PublicKey] interface {
	Address() common.Address
	Author(header *types.Header[P]) (common.Address, error)
	ExtractGenesisValidators(header *types.Header[P]) ([]common.Address, error)
	Signers(header *types.Header[P]) ([]common.Address, error)
	CommitHeader(header *types.Header[P], seals [][]byte, round *big.Int) error
	VerifyBlockProposal(chain consensus.ChainHeaderReader[P], block *types.Block[P], validators ValidatorSet) (time.Duration, error)
	VerifyHeader(chain consensus.ChainHeaderReader[P], header *types.Header[P], parents []*types.Header[P], validators ValidatorSet) error
	VerifyUncles(chain consensus.ChainReader[P], block *types.Block[P]) error
	VerifySeal(chain consensus.ChainHeaderReader[P], header *types.Header[P], validators ValidatorSet) error
	Prepare(chain consensus.ChainHeaderReader[P], header *types.Header[P], validators ValidatorSet) error
	Finalize(chain consensus.ChainHeaderReader[P], header *types.Header[P], state *state.StateDB[P], txs []*types.Transaction[P], uncles []*types.Header[P])
	FinalizeAndAssemble(chain consensus.ChainHeaderReader[P], header *types.Header[P], state *state.StateDB[P], txs []*types.Transaction[P], uncles []*types.Header[P], receipts []*types.Receipt[P]) (*types.Block[P], error)
	Seal(chain consensus.ChainHeaderReader[P], block *types.Block[P], validators ValidatorSet) (*types.Block[P], error)
	SealHash(header *types.Header[P]) common.Hash
	CalcDifficulty(chain consensus.ChainHeaderReader[P], time uint64, parent *types.Header[P]) *big.Int
	WriteVote(header *types.Header[P], candidate common.Address, authorize bool) error
	ReadVote(header *types.Header[P]) (candidate common.Address, authorize bool, err error)
}
