package mps

import (
	"fmt"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/privatecache"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/log"
)

// DefaultPrivateStateRepository acts as the single private state in the original
// Quorum design.
type DefaultPrivateStateRepository [P crypto.PublicKey] struct {
	db ethdb.Database
	// cache of stateDB
	stateCache           state.Database
	privateCacheProvider privatecache.Provider
	// stateDB gives access to the underlying state
	stateDB *state.StateDB[P]
	root    common.Hash
}

func NewDefaultPrivateStateRepository[P crypto.PublicKey](db ethdb.Database, cache state.Database, privateCacheProvider privatecache.Provider, previousBlockHash common.Hash) (*DefaultPrivateStateRepository[P], error) {
	root := rawdb.GetPrivateStateRoot(db, previousBlockHash)

	statedb, err := state.New[P](root, cache, nil)
	if err != nil {
		return nil, err
	}

	return &DefaultPrivateStateRepository[P]{
		db:                   db,
		stateCache:           cache,
		privateCacheProvider: privateCacheProvider,
		stateDB:              statedb,
		root:                 root,
	}, nil
}

func (dpsr *DefaultPrivateStateRepository[P]) DefaultState() (*state.StateDB[P], error) {
	if dpsr == nil {
		return nil, fmt.Errorf("nil instance")
	}
	return dpsr.stateDB, nil
}

func (dpsr *DefaultPrivateStateRepository[P]) DefaultStateMetadata() *PrivateStateMetadata {
	return DefaultPrivateStateMetadata
}

func (dpsr *DefaultPrivateStateRepository[P]) IsMPS() bool {
	return false
}

func (dpsr *DefaultPrivateStateRepository[P]) PrivateStateRoot(psi types.PrivateStateIdentifier) (common.Hash, error) {
	return dpsr.root, nil
}

func (dpsr *DefaultPrivateStateRepository[P]) StatePSI(psi types.PrivateStateIdentifier) (*state.StateDB[P], error) {
	if psi != types.DefaultPrivateStateIdentifier {
		return nil, fmt.Errorf("only the 'private' psi is supported by the default private state manager")
	}
	return dpsr.stateDB, nil
}

func (dpsr *DefaultPrivateStateRepository[P]) Reset() error {
	return dpsr.stateDB.Reset(dpsr.root)
}

// CommitAndWrite commits the private state and writes to disk
func (dpsr *DefaultPrivateStateRepository[P]) CommitAndWrite(isEIP158 bool, block *types.Block[P]) error {
	privateRoot, err := dpsr.stateDB.Commit(isEIP158)
	if err != nil {
		return err
	}

	if err := rawdb.WritePrivateStateRoot(dpsr.db, block.Root(), privateRoot); err != nil {
		log.Error("Failed writing private state root", "err", err)
		return err
	}
	dpsr.privateCacheProvider.Commit(dpsr.stateCache, privateRoot)
	dpsr.privateCacheProvider.Reference(privateRoot, block.Root())
	return nil
}

// Commit commits the private state only
func (dpsr *DefaultPrivateStateRepository[P]) Commit(isEIP158 bool, block *types.Block[P]) error {
	var err error
	dpsr.root, err = dpsr.stateDB.Commit(isEIP158)
	return err
}

func (dpsr *DefaultPrivateStateRepository[P]) Copy() PrivateStateRepository[P] {
	return &DefaultPrivateStateRepository[P]{
		db:                   dpsr.db,
		stateCache:           dpsr.stateCache,
		privateCacheProvider: dpsr.privateCacheProvider,
		stateDB:              dpsr.stateDB.Copy(),
		root:                 dpsr.root,
	}
}

// Given a slice of public receipts and an overlapping (smaller) slice of
// private receipts, return a new slice where the default for each location is
// the public receipt but we take the private receipt in each place we have
// one.
func (dpsr *DefaultPrivateStateRepository[P]) MergeReceipts(pub, priv types.Receipts[P]) types.Receipts[P] {
	m := make(map[common.Hash]*types.Receipt[P])
	for _, receipt := range pub {
		m[receipt.TxHash] = receipt
	}
	for _, receipt := range priv {
		m[receipt.TxHash] = receipt
	}

	ret := make(types.Receipts[P], 0, len(pub))
	for _, pubReceipt := range pub {
		ret = append(ret, m[pubReceipt.TxHash])
	}

	return ret
}
