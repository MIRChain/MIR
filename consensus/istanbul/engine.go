package istanbul

import (
	"math/big"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type Engine [P crypto.PublicKey] interface {
	Address() common.Address
	Author(header *types.Header) (common.Address, error)
	ExtractGenesisValidators(header *types.Header) ([]common.Address, error)
	Signers(header *types.Header) ([]common.Address, error)
	CommitHeader(header *types.Header, seals [][]byte, round *big.Int) error
	VerifyBlockProposal(chain consensus.ChainHeaderReader, block *types.Block[P], validators ValidatorSet) (time.Duration, error)
	VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators ValidatorSet) error
	VerifyUncles(chain consensus.ChainReader[P], block *types.Block[P]) error
	VerifySeal(chain consensus.ChainHeaderReader, header *types.Header, validators ValidatorSet) error
	Prepare(chain consensus.ChainHeaderReader, header *types.Header, validators ValidatorSet) error
	Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB[P], txs []*types.Transaction[P], uncles []*types.Header)
	FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB[P], txs []*types.Transaction[P], uncles []*types.Header, receipts []*types.Receipt[P]) (*types.Block[P], error)
	Seal(chain consensus.ChainHeaderReader, block *types.Block[P], validators ValidatorSet) (*types.Block[P], error)
	SealHash(header *types.Header) common.Hash
	CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int
	WriteVote(header *types.Header, candidate common.Address, authorize bool) error
	ReadVote(header *types.Header) (candidate common.Address, authorize bool, err error)
}
