package core

import (
	"context"
	"fmt"

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

type MultiplePrivateStateManager [P crypto.PublicKey] struct {
	// Low level persistent database to store final content in
	db                     ethdb.Database
	privateStatesTrieCache state.Database
	privateCacheProvider   privatecache.Provider

	residentGroupByKey map[string]*mps.PrivateStateMetadata
	privacyGroupById   map[types.PrivateStateIdentifier]*mps.PrivateStateMetadata
}

func newMultiplePrivateStateManager [P crypto.PublicKey] (db ethdb.Database, privateCacheProvider privatecache.Provider, residentGroupByKey map[string]*mps.PrivateStateMetadata, privacyGroupById map[types.PrivateStateIdentifier]*mps.PrivateStateMetadata) (*MultiplePrivateStateManager[P], error) {
	return &MultiplePrivateStateManager[P]{
		db:                     db,
		privateStatesTrieCache: privateCacheProvider.GetCacheWithConfig(),
		privateCacheProvider:   privateCacheProvider,
		residentGroupByKey:     residentGroupByKey,
		privacyGroupById:       privacyGroupById,
	}, nil
}

func (m *MultiplePrivateStateManager[P]) StateRepository(blockHash common.Hash) (mps.PrivateStateRepository[P], error) {
	privateStatesTrieRoot := rawdb.GetPrivateStatesTrieRoot(m.db, blockHash)
	return mps.NewMultiplePrivateStateRepository[P](m.db, m.privateStatesTrieCache, privateStatesTrieRoot, m.privateCacheProvider)
}

func (m *MultiplePrivateStateManager[P]) ResolveForManagedParty(managedParty string) (*mps.PrivateStateMetadata, error) {
	psm, found := m.residentGroupByKey[managedParty]
	if !found {
		return nil, fmt.Errorf("unable to find private state metadata for managed party %s", managedParty)
	}
	return psm, nil
}

func (m *MultiplePrivateStateManager[P]) ResolveForUserContext(ctx context.Context) (*mps.PrivateStateMetadata, error) {
	psi, ok := rpc.PrivateStateIdentifierFromContext(ctx)
	if !ok {
		psi = types.DefaultPrivateStateIdentifier
	}
	psm, found := m.privacyGroupById[psi]
	if !found {
		return nil, fmt.Errorf("unable to find private state for context psi %s", psi)
	}
	return psm, nil
}

func (m *MultiplePrivateStateManager[P]) PSIs() []types.PrivateStateIdentifier {
	psis := make([]types.PrivateStateIdentifier, 0, len(m.privacyGroupById))
	for psi := range m.privacyGroupById {
		psis = append(psis, psi)
	}
	return psis
}

func (m *MultiplePrivateStateManager[P]) NotIncludeAny(psm *mps.PrivateStateMetadata, managedParties ...string) bool {
	return psm.NotIncludeAny(managedParties...)
}

func (m *MultiplePrivateStateManager[P]) CheckAt(root common.Hash) error {
	_, err := state.New(rawdb.GetPrivateStatesTrieRoot(m.db, root), m.privateStatesTrieCache, nil)
	return err
}

func (m *MultiplePrivateStateManager[P]) TrieDB() *trie.Database {
	return m.privateStatesTrieCache.TrieDB()
}
