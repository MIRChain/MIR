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

package backend

import (
	"math/big"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul"
	istanbulcommon "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/common"
	ibftcore "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/ibft/core"
	ibftengine "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/ibft/engine"
	qbftcore "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/qbft/core"
	qbftengine "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/qbft/engine"
	qbfttypes "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/qbft/types"
	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul/validator"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/log"
)

const (
	// fetcherID is the ID indicates the block is from Istanbul engine
	fetcherID = "istanbul"
)

// New creates an Ethereum backend for Istanbul core engine.
func New[T crypto.PrivateKey, P crypto.PublicKey](config *istanbul.Config, privateKey T, db ethdb.Database) *Backend[T,P] {
	// Allocate the snapshot caches and create the engine
	recents, _ := lru.NewARC(inmemorySnapshots)
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)
	var pub P
	switch t:=any(&privateKey).(type) {
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
	}
	sb := &Backend[T,P]{
		config:           config,
		istanbulEventMux: new(event.TypeMux),
		privateKey:       privateKey,
		address:          crypto.PubkeyToAddress(pub),
		logger:           log.New(),
		db:               db,
		commitCh:         make(chan *types.Block[P], 1),
		recents:          recents,
		candidates:       make(map[common.Address]bool),
		coreStarted:      false,
		recentMessages:   recentMessages,
		knownMessages:    knownMessages,
	}

	sb.qbftEngine = qbftengine.NewEngine[P](sb.config, sb.address, sb.Sign)
	sb.ibftEngine = ibftengine.NewEngine[P](sb.config, sb.address, sb.Sign)

	return sb
}

// ----------------------------------------------------------------------------

type Backend [T crypto.PrivateKey, P crypto.PublicKey] struct {
	config *istanbul.Config

	privateKey T
	address    common.Address

	core istanbul.Core

	ibftEngine *ibftengine.Engine[P]
	qbftEngine *qbftengine.Engine[P]

	istanbulEventMux *event.TypeMux

	logger log.Logger

	db ethdb.Database

	chain        consensus.ChainHeaderReader[P]
	currentBlock func() *types.Block[P]
	hasBadBlock  func(db ethdb.Reader, hash common.Hash) bool

	// the channels for istanbul engine notifications
	commitCh          chan *types.Block[P]
	proposedBlockHash common.Hash
	sealMu            sync.Mutex
	coreStarted       bool
	coreMu            sync.RWMutex

	// Current list of candidates we are pushing
	candidates map[common.Address]bool
	// Protects the signer fields
	candidatesLock sync.RWMutex
	// Snapshots for recent block to speed up reorgs
	recents *lru.ARCCache

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster[P]

	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages

	qbftConsensusEnabled bool // qbft consensus
}

func (sb *Backend[T,P]) Engine() istanbul.Engine[P] {
	return sb.EngineForBlockNumber(nil)
}

func (sb *Backend[T,P]) EngineForBlockNumber(blockNumber *big.Int) istanbul.Engine[P] {
	switch {
	case blockNumber != nil && sb.IsQBFTConsensusAt(blockNumber):
		return sb.qbftEngine
	case blockNumber == nil && sb.IsQBFTConsensus():
		return sb.qbftEngine
	default:
		return sb.ibftEngine
	}
}

// zekun: HACK
func (sb *Backend[T,P]) CalcDifficulty(chain consensus.ChainHeaderReader[P], time uint64, parent *types.Header[P]) *big.Int {
	return sb.EngineForBlockNumber(parent.Number).CalcDifficulty(chain, time, parent)
}

// Address implements istanbul.Backend.Address
func (sb *Backend[T,P]) Address() common.Address {
	return sb.Engine().Address()
}

// Validators implements istanbul.Backend.Validators
func (sb *Backend[T,P]) Validators(proposal istanbul.Proposal) istanbul.ValidatorSet {
	return sb.getValidators(proposal.Number().Uint64(), proposal.Hash())
}

// Broadcast implements istanbul.Backend.Broadcast
func (sb *Backend[T,P]) Broadcast(valSet istanbul.ValidatorSet, code uint64, payload []byte) error {
	// send to others
	sb.Gossip(valSet, code, payload)
	// send to self
	msg := istanbul.MessageEvent{
		Code:    code,
		Payload: payload,
	}
	go sb.istanbulEventMux.Post(msg)
	return nil
}

// Gossip implements istanbul.Backend.Gossip
func (sb *Backend[T,P]) Gossip(valSet istanbul.ValidatorSet, code uint64, payload []byte) error {
	hash := istanbul.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	targets := make(map[common.Address]bool)
	for _, val := range valSet.List() {
		if val.Address() != sb.Address() {
			targets[val.Address()] = true
		}
	}
	if sb.broadcaster != nil && len(targets) > 0 {
		ps := sb.broadcaster.FindPeers(targets)
		for addr, p := range ps {
			ms, ok := sb.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)

			if sb.IsQBFTConsensus() {
				var outboundCode uint64 = istanbulMsg
				if _, ok := qbfttypes.MessageCodes()[code]; ok {
					outboundCode = code
				}
				go p.SendQBFTConsensus(outboundCode, payload)
			} else {
				go p.SendConsensus(istanbulMsg, payload)
			}
		}
	}
	return nil
}

// Commit implements istanbul.Backend.Commit
func (sb *Backend[T,P]) Commit(proposal istanbul.Proposal, seals [][]byte, round *big.Int) (err error) {
	// Check if the proposal is a valid block
	block, ok := proposal.(*types.Block[P])
	if !ok {
		sb.logger.Error("BFT: invalid block proposal", "proposal", proposal)
		return istanbulcommon.ErrInvalidProposal
	}

	// Commit header
	h := block.Header()
	err = sb.EngineForBlockNumber(h.Number).CommitHeader(h, seals, round)
	if err != nil {
		return
	}

	// Remove ValidatorSet added to ProposerPolicy registry, if not done, the registry keeps increasing size with each block height
	sb.config.ProposerPolicy.ClearRegistry()

	// update block's header
	block = block.WithSeal(h)

	sb.logger.Info("BFT: block proposal committed", "author", sb.Address(), "hash", proposal.Hash(), "number", proposal.Number().Uint64())

	// - if the proposed and committed blocks are the same, send the proposed hash
	//   to commit channel, which is being watched inside the engine.Seal() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	if sb.proposedBlockHash == block.Hash() {
		// feed block hash to Seal() and wait the Seal() result
		sb.commitCh <- block
		return nil
	}

	if sb.broadcaster != nil {
		sb.broadcaster.Enqueue(fetcherID, block)
	}

	return nil
}

// EventMux implements istanbul.Backend.EventMux
func (sb *Backend[T,P]) EventMux() *event.TypeMux {
	return sb.istanbulEventMux
}

// Verify implements istanbul.Backend.Verify
func (sb *Backend[T,P]) Verify(proposal istanbul.Proposal) (time.Duration, error) {
	// Check if the proposal is a valid block
	block, ok := proposal.(*types.Block[P])
	if !ok {
		sb.logger.Error("BFT: invalid block proposal", "proposal", proposal)
		return 0, istanbulcommon.ErrInvalidProposal
	}

	// check bad block
	if sb.HasBadProposal(block.Hash()) {
		sb.logger.Warn("BFT: bad block proposal", "proposal", proposal)
		return 0, core.ErrBlacklistedHash
	}

	header := block.Header()
	snap, err := sb.snapshot(sb.chain, header.Number.Uint64()-1, header.ParentHash, nil)
	if err != nil {
		return 0, err
	}

	return sb.EngineForBlockNumber(header.Number).VerifyBlockProposal(sb.chain, block, snap.ValSet)
}

// Sign implements istanbul.Backend.Sign
func (sb *Backend[T,P]) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256[P](data)
	return crypto.Sign(hashData, sb.privateKey)
}

// SignWithoutHashing implements istanbul.Backend.SignWithoutHashing and signs input data with the backend's private key without hashing the input data
func (sb *Backend[T,P]) SignWithoutHashing(data []byte) ([]byte, error) {
	return crypto.Sign(data, sb.privateKey)
}

// CheckSignature implements istanbul.Backend.CheckSignature
func (sb *Backend[T,P]) CheckSignature(data []byte, address common.Address, sig []byte) error {
	signer, err := istanbul.GetSignatureAddress[P](data, sig)
	if err != nil {
		return err
	}
	// Compare derived addresses
	if signer != address {
		return istanbulcommon.ErrInvalidSignature
	}

	return nil
}

// HasPropsal implements istanbul.Backend.HashBlock
func (sb *Backend[T,P]) HasPropsal(hash common.Hash, number *big.Int) bool {
	return sb.chain.GetHeader(hash, number.Uint64()) != nil
}

// GetProposer implements istanbul.Backend.GetProposer
func (sb *Backend[T,P]) GetProposer(number uint64) common.Address {
	if h := sb.chain.GetHeaderByNumber(number); h != nil {
		a, _ := sb.Author(h)
		return a
	}
	return common.Address{}
}

// ParentValidators implements istanbul.Backend.GetParentValidators
func (sb *Backend[T,P]) ParentValidators(proposal istanbul.Proposal) istanbul.ValidatorSet {
	if block, ok := proposal.(*types.Block[P]); ok {
		return sb.getValidators(block.Number().Uint64()-1, block.ParentHash())
	}
	return validator.NewSet(nil, sb.config.ProposerPolicy)
}

func (sb *Backend[T,P]) getValidators(number uint64, hash common.Hash) istanbul.ValidatorSet {
	snap, err := sb.snapshot(sb.chain, number, hash, nil)
	if err != nil {
		return validator.NewSet(nil, sb.config.ProposerPolicy)
	}
	return snap.ValSet
}

func (sb *Backend[T,P]) LastProposal() (istanbul.Proposal, common.Address) {
	block := sb.currentBlock()

	var proposer common.Address
	if block.Number().Cmp(common.Big0) > 0 {
		var err error
		proposer, err = sb.Author(block.Header())
		if err != nil {
			sb.logger.Error("BFT: last block proposal invalid", "err", err)
			return nil, common.Address{}
		}
	}

	// Return header only block here since we don't need block body
	return block, proposer
}

func (sb *Backend[T,P]) HasBadProposal(hash common.Hash) bool {
	if sb.hasBadBlock == nil {
		return false
	}
	return sb.hasBadBlock(sb.db, hash)
}

func (sb *Backend[T,P]) Close() error {
	return nil
}

// IsQBFTConsensus returns whether qbft consensus should be used
func (sb *Backend[T,P]) IsQBFTConsensus() bool {
	if sb.qbftConsensusEnabled {
		return true
	}
	if sb.chain != nil {
		qbftEnabled := sb.IsQBFTConsensusAt(sb.chain.CurrentHeader().Number)
		sb.qbftConsensusEnabled = qbftEnabled
		return qbftEnabled
	}
	return false
}

// IsQBFTConsensusForHeader checks if qbft consensus is enabled for the block height identified by the given header
func (sb *Backend[T,P]) IsQBFTConsensusAt(blockNumber *big.Int) bool {
	return sb.config.IsQBFTConsensusAt(blockNumber)
}

func (sb *Backend[T,P]) startIBFT() error {
	sb.logger.Info("BFT: activate IBFT")
	sb.logger.Trace("BFT: set ProposerPolicy sorter to ValidatorSortByStringFun")
	sb.config.ProposerPolicy.Use(istanbul.ValidatorSortByString())
	sb.qbftConsensusEnabled = false

	sb.core = ibftcore.New[P](sb, sb.config)
	if err := sb.core.Start(); err != nil {
		sb.logger.Error("BFT: failed to activate IBFT", "err", err)
		return err
	}

	return nil
}

func (sb *Backend[T,P]) startQBFT() error {
	sb.logger.Info("BFT: activate QBFT")
	sb.logger.Trace("BFT: set ProposerPolicy sorter to ValidatorSortByByteFunc")
	sb.config.ProposerPolicy.Use(istanbul.ValidatorSortByByte())
	sb.qbftConsensusEnabled = true

	sb.core = qbftcore.New[P](sb, sb.config)
	if err := sb.core.Start(); err != nil {
		sb.logger.Error("BFT: failed to activate QBFT", "err", err)
		return err
	}

	return nil
}

func (sb *Backend[T,P]) stop() error {
	core := sb.core
	sb.core = nil

	if core != nil {
		sb.logger.Info("BFT: deactivate")
		if err := core.Stop(); err != nil {
			sb.logger.Error("BFT: failed to deactivate", "err", err)
			return err
		}
	}

	sb.qbftConsensusEnabled = false

	return nil
}

// StartQBFTConsensus stops existing legacy ibft consensus and starts the new qbft consensus
func (sb *Backend[T,P]) StartQBFTConsensus() error {
	sb.logger.Info("BFT: switch from IBFT to QBFT")
	if err := sb.stop(); err != nil {
		return err
	}

	return sb.startQBFT()
}
