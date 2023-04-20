package core

import (
	"context"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/privatecache"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

type DefaultPrivateStateManager [P crypto.PublicKey]struct {
	// Low level persistent database to store final content in
	db                   ethdb.Database
	repoCache            state.Database
	privateCacheProvider privatecache.Provider
}

func newDefaultPrivateStateManager[P crypto.PublicKey](db ethdb.Database, privateCacheProvider privatecache.Provider) *DefaultPrivateStateManager[P] {
	return &DefaultPrivateStateManager[P]{
		db:                   db,
		repoCache:            privateCacheProvider.GetCacheWithConfig(),
		privateCacheProvider: privateCacheProvider,
	}
}

func (d *DefaultPrivateStateManager[P]) StateRepository(blockHash common.Hash) (mps.PrivateStateRepository[P], error) {
	return mps.NewDefaultPrivateStateRepository[P](d.db, d.repoCache, d.privateCacheProvider, blockHash)
}

func (d *DefaultPrivateStateManager[P]) ResolveForManagedParty(_ string) (*mps.PrivateStateMetadata, error) {
	return mps.DefaultPrivateStateMetadata, nil
}

func (d *DefaultPrivateStateManager[P]) ResolveForUserContext(ctx context.Context) (*mps.PrivateStateMetadata, error) {
	psi, ok := rpc.PrivateStateIdentifierFromContext(ctx)
	if !ok {
		psi = types.DefaultPrivateStateIdentifier
	}
	return &mps.PrivateStateMetadata{ID: psi, Type: mps.Resident}, nil
}

func (d *DefaultPrivateStateManager[P]) PSIs() []types.PrivateStateIdentifier {
	return []types.PrivateStateIdentifier{
		types.DefaultPrivateStateIdentifier,
	}
}

func (d *DefaultPrivateStateManager[P]) NotIncludeAny(_ *mps.PrivateStateMetadata, _ ...string) bool {
	// with default implementation, all managedParties are members of the psm
	return false
}

func (d *DefaultPrivateStateManager[P]) CheckAt(root common.Hash) error {
	_, err := state.New[P](rawdb.GetPrivateStateRoot(d.db, root), d.repoCache, nil)
	return err
}

func (d *DefaultPrivateStateManager[P]) TrieDB() *trie.Database {
	return d.repoCache.TrieDB()
}
