package privatecache

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type segregatedCacheProvider [P crypto.PublicKey] struct {
	db     ethdb.Database
	config *trie.Config
}

func (p *segregatedCacheProvider[P]) GetCache() state.Database {
	return state.NewDatabase[P](p.db)
}

func (p *segregatedCacheProvider[P]) GetCacheWithConfig() state.Database {
	return state.NewDatabaseWithConfig[P](p.db, p.config)
}

func (p *segregatedCacheProvider[P]) Commit(db state.Database, hash common.Hash) error {
	return db.TrieDB().Commit(hash, false, nil)
}
func (p *segregatedCacheProvider[P]) Reference(child, parent common.Hash) {
	// do nothing
}
