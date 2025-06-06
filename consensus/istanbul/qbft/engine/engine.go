package qbftengine

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul/backend/contract"
	istanbulcommon "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul/validator"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// var (
// 	nilUncleHash = types.CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
// )

type SignerFn func(data []byte) ([]byte, error)

type Engine [P crypto.PublicKey] struct {
	cfg *istanbul.Config

	signer common.Address // Ethereum address of the signing key
	sign   SignerFn       // Signer function to authorize hashes with
}

func NewEngine[P crypto.PublicKey](cfg *istanbul.Config, signer common.Address, sign SignerFn) *Engine[P] {
	return &Engine[P]{
		cfg:    cfg,
		signer: signer,
		sign:   sign,
	}
}

func (e *Engine[P]) Author(header *types.Header[P]) (common.Address, error) {
	return header.Coinbase, nil
}

func (e *Engine[P]) CommitHeader(header *types.Header[P], seals [][]byte, round *big.Int) error {
	return ApplyHeaderQBFTExtra(
		header,
		writeCommittedSeals(seals),
		writeRoundNumber(round),
	)
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
func writeCommittedSeals(committedSeals [][]byte) ApplyQBFTExtra {
	return func(qbftExtra *types.QBFTExtra) error {
		if len(committedSeals) == 0 {
			return istanbulcommon.ErrInvalidCommittedSeals
		}

		for _, seal := range committedSeals {
			if len(seal) != types.IstanbulExtraSeal {
				return istanbulcommon.ErrInvalidCommittedSeals
			}
		}

		qbftExtra.CommittedSeal = make([][]byte, len(committedSeals))
		copy(qbftExtra.CommittedSeal, committedSeals)

		return nil
	}
}

// writeRoundNumber writes the extra-data field of a block header with given round.
func writeRoundNumber(round *big.Int) ApplyQBFTExtra {
	return func(qbftExtra *types.QBFTExtra) error {
		qbftExtra.Round = uint32(round.Uint64())
		return nil
	}
}

func (e *Engine[P]) VerifyBlockProposal(chain consensus.ChainHeaderReader[P], block *types.Block[P], validators istanbul.ValidatorSet) (time.Duration, error) {
	// check block body
	txnHash := types.DeriveSha(block.Transactions(), new(trie.Trie[P]))
	if txnHash != block.Header().TxHash {
		return 0, istanbulcommon.ErrMismatchTxhashes
	}

	uncleHash := types.CalcUncleHash(block.Uncles())
	if uncleHash != types.CalcUncleHash[P](nil) {
		return 0, istanbulcommon.ErrInvalidUncleHash
	}

	// verify the header of proposed block
	err := e.VerifyHeader(chain, block.Header(), nil, validators)
	if err == nil || err == istanbulcommon.ErrEmptyCommittedSeals {
		// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
		return 0, nil
	} else if err == consensus.ErrFutureBlock {
		return time.Until(time.Unix(int64(block.Header().Time), 0)), consensus.ErrFutureBlock
	}

	parentHeader := chain.GetHeaderByHash(block.ParentHash())
	config := e.cfg.GetConfig(parentHeader.Number)

	if config.EmptyBlockPeriod > config.BlockPeriod && len(block.Transactions()) == 0 {
		if block.Header().Time < parentHeader.Time+config.EmptyBlockPeriod {
			return 0, fmt.Errorf("empty block verification fail")
		}
	}

	return 0, err
}

func (e *Engine[P]) VerifyHeader(chain consensus.ChainHeaderReader[P], header *types.Header[P], parents []*types.Header[P], validators istanbul.ValidatorSet) error {
	return e.verifyHeader(chain, header, parents, validators)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (e *Engine[P]) verifyHeader(chain consensus.ChainHeaderReader[P], header *types.Header[P], parents []*types.Header[P], validators istanbul.ValidatorSet) error {
	if header.Number == nil {
		return istanbulcommon.ErrUnknownBlock
	}

	// Don't waste time checking blocks from the future (adjusting for allowed threshold)
	adjustedTimeNow := time.Now().Add(time.Duration(e.cfg.AllowedFutureBlockTime) * time.Second).Unix()
	if header.Time > uint64(adjustedTimeNow) {
		return consensus.ErrFutureBlock
	}

	if _, err := types.ExtractQBFTExtra(header); err != nil {
		return istanbulcommon.ErrInvalidExtraDataFormat
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.IstanbulDigest {
		return istanbulcommon.ErrInvalidMixDigest
	}

	// Ensure that the block doesn't contain any uncles which are meaningless in Istanbul
	if header.UncleHash != types.CalcUncleHash[P](nil) {
		return istanbulcommon.ErrInvalidUncleHash
	}

	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Difficulty == nil || header.Difficulty.Cmp(istanbulcommon.DefaultDifficulty) != 0 {
		return istanbulcommon.ErrInvalidDifficulty
	}

	return e.verifyCascadingFields(chain, header, validators, parents)
}

func (e *Engine[P]) VerifyHeaders(chain consensus.ChainHeaderReader[P], headers []*types.Header[P], seals []bool, validators istanbul.ValidatorSet) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		errored := false
		for i, header := range headers {
			var err error
			if errored {
				err = consensus.ErrUnknownAncestor
			} else {
				err = e.verifyHeader(chain, header, headers[:i], validators)
			}

			if err != nil {
				errored = true
			}

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (e *Engine[P]) verifyCascadingFields(chain consensus.ChainHeaderReader[P], header *types.Header[P], validators istanbul.ValidatorSet, parents []*types.Header[P]) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// Check parent
	var parent *types.Header[P]
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}

	// Ensure that the block's parent has right number and hash
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}

	// Ensure that the block's timestamp isn't too close to it's parent
	if parent.Time+e.cfg.GetConfig(parent.Number).BlockPeriod > header.Time {
		return istanbulcommon.ErrInvalidTimestamp
	}

	// Verify signer
	if err := e.verifySigner(chain, header, parents, validators); err != nil {
		return err
	}

	return e.verifyCommittedSeals(chain, header, parents, validators)
}

func (e *Engine[P]) verifySigner(chain consensus.ChainHeaderReader[P], header *types.Header[P], parents []*types.Header[P], validators istanbul.ValidatorSet) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return istanbulcommon.ErrUnknownBlock
	}

	// Resolve the authorization key and check against signers
	signer, err := e.Author(header)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := validators.GetByAddress(signer); v == nil {
		return istanbulcommon.ErrUnauthorized
	}

	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (e *Engine[P]) verifyCommittedSeals(chain consensus.ChainHeaderReader[P], header *types.Header[P], parents []*types.Header[P], validators istanbul.ValidatorSet) error {
	number := header.Number.Uint64()

	if number == 0 {
		// We don't need to verify committed seals in the genesis block
		return nil
	}

	extra, err := types.ExtractQBFTExtra(header)
	if err != nil {
		return err
	}
	committedSeal := extra.CommittedSeal

	// The length of Committed seals should be larger than 0
	if len(committedSeal) == 0 {
		return istanbulcommon.ErrEmptyCommittedSeals
	}

	validatorsCpy := validators.Copy()

	// Check whether the committed seals are generated by validators
	validSeal := 0
	committers, err := e.Signers(header)
	if err != nil {
		return err
	}

	for _, addr := range committers {
		if validatorsCpy.RemoveValidator(addr) {
			validSeal++
			continue
		}
		return istanbulcommon.ErrInvalidCommittedSeals
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= validators.F() {
		return istanbulcommon.ErrInvalidCommittedSeals
	}

	return nil
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (e *Engine[P]) VerifyUncles(chain consensus.ChainReader[P], block *types.Block[P]) error {
	if len(block.Uncles()) > 0 {
		return istanbulcommon.ErrInvalidUncleHash
	}
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (e *Engine[P]) VerifySeal(chain consensus.ChainHeaderReader[P], header *types.Header[P], validators istanbul.ValidatorSet) error {
	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return istanbulcommon.ErrUnknownBlock
	}

	// ensure that the difficulty equals to istanbulcommon.DefaultDifficulty
	if header.Difficulty.Cmp(istanbulcommon.DefaultDifficulty) != 0 {
		return istanbulcommon.ErrInvalidDifficulty
	}

	return e.verifySigner(chain, header, nil, validators)
}

func (e *Engine[P]) Prepare(chain consensus.ChainHeaderReader[P], header *types.Header[P], validators istanbul.ValidatorSet) error {
	header.Coinbase = common.Address{}
	header.Nonce = istanbulcommon.EmptyBlockNonce
	header.MixDigest = types.IstanbulDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// use the same difficulty for all blocks
	header.Difficulty = istanbulcommon.DefaultDifficulty

	// set header's timestamp
	header.Time = parent.Time + e.cfg.GetConfig(header.Number).BlockPeriod
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}

	currentBlockNumber := big.NewInt(0).SetUint64(number - 1)
	validatorContract := e.cfg.GetValidatorContractAddress(currentBlockNumber)
	if validatorContract != (common.Address{}) && e.cfg.GetValidatorSelectionMode(currentBlockNumber) == params.ContractMode {
		return ApplyHeaderQBFTExtra(
			header,
			WriteValidators([]common.Address{}),
		)
	} else {
		for _, transition := range e.cfg.Transitions {
			if transition.Block.Cmp(currentBlockNumber) == 0 && len(transition.Validators) > 0 {
				toRemove := make([]istanbul.Validator, 0, validators.Size())
				l := validators.List()
				for i := range l {
					toRemove = append(toRemove, l[i])
				}
				for i := range toRemove {
					validators.RemoveValidator(toRemove[i].Address())
				}
				for i := range transition.Validators {
					validators.AddValidator(transition.Validators[i])
				}
				break
			}
		}
		validatorsList := validator.SortedAddresses(validators.List())
		// add validators in snapshot to extraData's validators section
		return ApplyHeaderQBFTExtra(
			header,
			WriteValidators(validatorsList),
		)
	}
}

func WriteValidators(validators []common.Address) ApplyQBFTExtra {
	return func(qbftExtra *types.QBFTExtra) error {
		qbftExtra.Validators = validators
		return nil
	}
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (e *Engine[P]) Finalize(chain consensus.ChainHeaderReader[P], header *types.Header[P], state *state.StateDB[P], txs []*types.Transaction[P], uncles []*types.Header[P]) {
	// Accumulate any block and uncle rewards and commit the final state root
	e.accumulateRewards(chain, state, header, uncles, e.cfg.GetConfig(header.Number))
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash[P](nil)
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (e *Engine[P]) FinalizeAndAssemble(chain consensus.ChainHeaderReader[P], header *types.Header[P], state *state.StateDB[P], txs []*types.Transaction[P], uncles []*types.Header[P], receipts []*types.Receipt[P]) (*types.Block[P], error) {
	e.Finalize(chain, header, state, txs, uncles)
	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, new(trie.Trie[P])), nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (e *Engine[P]) Seal(chain consensus.ChainHeaderReader[P], block *types.Block[P], validators istanbul.ValidatorSet) (*types.Block[P], error) {
	if _, v := validators.GetByAddress(e.signer); v == nil {
		return block, istanbulcommon.ErrUnauthorized
	}

	header := block.Header()
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return block, consensus.ErrUnknownAncestor
	}

	// Set Coinbase
	header.Coinbase = e.signer

	return block.WithSeal(header), nil
}

func (e *Engine[P]) SealHash(header *types.Header[P]) common.Hash {
	header.Coinbase = e.signer
	return sigHash(header)
}

func (e *Engine[P]) CalcDifficulty(chain consensus.ChainHeaderReader[P], time uint64, parent *types.Header[P]) *big.Int {
	return new(big.Int)
}

func (e *Engine[P]) ExtractGenesisValidators(header *types.Header[P]) ([]common.Address, error) {
	extra, err := types.ExtractQBFTExtra(header)
	if err != nil {
		return nil, err
	}

	return extra.Validators, nil
}

func (e *Engine[P]) Signers(header *types.Header[P]) ([]common.Address, error) {
	extra, err := types.ExtractQBFTExtra(header)
	if err != nil {
		return []common.Address{}, err
	}
	committedSeal := extra.CommittedSeal
	proposalSeal := PrepareCommittedSeal(header, extra.Round)

	var addrs []common.Address
	// 1. Get committed seals from current header
	for _, seal := range committedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := istanbul.GetSignatureAddressNoHashing[P](proposalSeal, seal)
		if err != nil {
			return nil, istanbulcommon.ErrInvalidSignature
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func (e *Engine[P]) Address() common.Address {
	return e.signer
}

// FIXME: Need to update this for Istanbul
// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash[P crypto.PublicKey](header *types.Header[P]) (hash common.Hash) {
	hasher := crypto.NewKeccakState[P]()

	rlp.Encode(hasher, types.QBFTFilteredHeader(header))
	hasher.Sum(hash[:0])
	return hash
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal[P crypto.PublicKey](header *types.Header[P], round uint32) []byte {
	h := types.CopyHeader(header)
	return h.QBFTHashWithRoundNumber(round).Bytes()
}

func (e *Engine[P]) WriteVote(header *types.Header[P], candidate common.Address, authorize bool) error {
	return ApplyHeaderQBFTExtra(
		header,
		WriteVote(candidate, authorize),
	)
}

func WriteVote(candidate common.Address, authorize bool) ApplyQBFTExtra {
	return func(qbftExtra *types.QBFTExtra) error {
		voteType := types.QBFTDropVote
		if authorize {
			voteType = types.QBFTAuthVote
		}

		vote := &types.ValidatorVote{RecipientAddress: candidate, VoteType: voteType}
		qbftExtra.Vote = vote
		return nil
	}
}

func (e *Engine[P]) ReadVote(header *types.Header[P]) (candidate common.Address, authorize bool, err error) {
	qbftExtra, err := getExtra(header)
	if err != nil {
		return common.Address{}, false, err
	}

	var vote *types.ValidatorVote
	if qbftExtra.Vote == nil {
		vote = &types.ValidatorVote{RecipientAddress: common.Address{}, VoteType: types.QBFTDropVote}
	} else {
		vote = qbftExtra.Vote
	}

	// Tally up the new vote from the validator
	switch {
	case vote.VoteType == types.QBFTAuthVote:
		authorize = true
	case vote.VoteType == types.QBFTDropVote:
		authorize = false
	default:
		return common.Address{}, false, istanbulcommon.ErrInvalidVote
	}

	return vote.RecipientAddress, authorize, nil
}

func getExtra[P crypto.PublicKey](header *types.Header[P]) (*types.QBFTExtra, error) {
	if len(header.Extra) < types.IstanbulExtraVanity {
		// In this scenario, the header extradata only contains client specific information, hence create a new qbftExtra and set vanity
		vanity := append(header.Extra, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(header.Extra))...)
		return &types.QBFTExtra{
			VanityData:    vanity,
			Validators:    []common.Address{},
			CommittedSeal: [][]byte{},
			Round:         0,
			Vote:          nil,
		}, nil
	}

	// This is the case when Extra has already been set
	return types.ExtractQBFTExtra(header)
}

func setExtra[P crypto.PublicKey](h *types.Header[P], qbftExtra *types.QBFTExtra) error {
	payload, err := rlp.EncodeToBytes(qbftExtra)
	if err != nil {
		return err
	}

	h.Extra = payload
	return nil
}

func (e *Engine[P]) validatorsList(genesis *types.Header[P], config istanbul.Config) ([]common.Address, error) {
	var validators []common.Address
	if config.ValidatorContract != (common.Address{}) && config.GetValidatorSelectionMode(big.NewInt(0)) == params.ContractMode {
		log.Info("Initialising snap with contract validators", "address", config.ValidatorContract, "client", config.Client)

		validatorContractCaller, err := contract.NewValidatorContractInterfaceCaller[P](config.ValidatorContract, config.Client)

		if err != nil {
			return nil, fmt.Errorf("invalid smart contract in genesis alloc: %w", err)
		}

		opts := bind.CallOpts{
			Pending:     false,
			BlockNumber: big.NewInt(0),
		}
		validators, err = validatorContractCaller.GetValidators(&opts)
		if err != nil {
			log.Error("QBFT: invalid smart contract in genesis alloc", "err", err)
			return nil, err
		}
	} else {
		// Get the validators from genesis to create a snapshot
		var err error
		validators, err = e.ExtractGenesisValidators(genesis)
		if err != nil {
			log.Error("BFT: invalid genesis block", "err", err)
			return nil, err
		}
	}
	return validators, nil
}

// AccumulateRewards credits the beneficiary of the given block with a reward.
func (e *Engine[P]) accumulateRewards(chain consensus.ChainHeaderReader[P], state *state.StateDB[P], header *types.Header[P], uncles []*types.Header[P], cfg istanbul.Config) {
	blockReward := cfg.BlockReward
	if blockReward == nil {
		return // no reward
	}
	// Accumulate the rewards for a beneficiary
	reward := big.Int(*blockReward)
	if cfg.BeneficiaryMode == nil || *cfg.BeneficiaryMode == "" {
		if cfg.MiningBeneficiary != nil {
			state.AddBalance(*cfg.MiningBeneficiary, &reward) // implicit besu compatible mode
		}
		return
	}
	switch strings.ToLower(*cfg.BeneficiaryMode) {
	case "fixed":
		if cfg.MiningBeneficiary != nil {
			state.AddBalance(*cfg.MiningBeneficiary, &reward)
		}
	case "list":
		for _, b := range cfg.BeneficiaryList {
			state.AddBalance(b, &reward)
		}
	case "validators":
		genesis := chain.GetHeaderByNumber(0)
		if err := e.VerifyHeader(chain, genesis, nil, nil); err != nil {
			log.Error("BFT: invalid genesis block", "err", err)
			return
		}
		list, err := e.validatorsList(genesis, cfg)
		if err != nil {
			log.Error("get validators list", "err", err)
			return
		}
		for _, b := range list {
			state.AddBalance(b, &reward)
		}
	}
}
