package privatecache

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

type segregatedCacheProvider struct {
	db     ethdb.Database
	config *trie.Config
}

func (p *segregatedCacheProvider) GetCache() state.Database {
	return state.NewDatabase(p.db)
}

func (p *segregatedCacheProvider) GetCacheWithConfig() state.Database {
	return state.NewDatabaseWithConfig(p.db, p.config)
}

func (p *segregatedCacheProvider) Commit(db state.Database, hash common.Hash) error {
	return db.TrieDB().Commit(hash, false, nil)
}
func (p *segregatedCacheProvider) Reference(child, parent common.Hash) {
	// do nothing
}
