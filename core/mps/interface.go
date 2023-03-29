//go:generate mockgen -source interface.go -destination=mock_interface.go -package=mps

package mps

import (
	"context"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// PrivateStateManager interface separates
type PrivateStateManager [P crypto.PublicKey] interface {
	PrivateStateMetadataResolver
	// StateRepository returns repository corresponding to a block hash
	StateRepository(blockHash common.Hash) (PrivateStateRepository[P], error)
	// CheckAt verifies if there's a state being managed at a block hash
	CheckAt(blockHash common.Hash) error
	// TrieDB returns the trie database
	TrieDB() *trie.Database
}

type PrivateStateMetadataResolver interface {
	ResolveForManagedParty(managedParty string) (*PrivateStateMetadata, error)
	ResolveForUserContext(ctx context.Context) (*PrivateStateMetadata, error)
	// PSIs returns list of types.PrivateStateIdentifier being managed
	PSIs() []types.PrivateStateIdentifier
	// NotIncludeAny returns true if NONE of the managedParties is a member
	// of the given psm, otherwise returns false
	NotIncludeAny(psm *PrivateStateMetadata, managedParties ...string) bool
}

// PrivateStateRepository abstracts how we handle private state(s) including
// retrieving from and peristing private states to the underlying database
type PrivateStateRepository [P crypto.PublicKey] interface {
	PrivateStateRoot(psi types.PrivateStateIdentifier) (common.Hash, error)
	StatePSI(psi types.PrivateStateIdentifier) (*state.StateDB, error)
	CommitAndWrite(isEIP158 bool, block *types.Block[P]) error
	Commit(isEIP158 bool, block *types.Block[P]) error
	Copy() PrivateStateRepository[P]
	Reset() error
	DefaultState() (*state.StateDB, error)
	DefaultStateMetadata() *PrivateStateMetadata
	IsMPS() bool
	MergeReceipts(pub, priv types.Receipts[P]) types.Receipts[P]
}
