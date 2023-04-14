package privatecache

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type Provider interface {
	GetCache() state.Database
	GetCacheWithConfig() state.Database
	Commit(db state.Database, hash common.Hash) error
	Reference(child, parent common.Hash)
}

func NewPrivateCacheProvider[P crypto.PublicKey](db ethdb.Database, config *trie.Config, cache state.Database, privateCacheEnabled bool) Provider {
	if privateCacheEnabled {
		log.Info("Using UnifiedCacheProvider.")
		return &unifiedCacheProvider{
			cache: cache,
		}
	}
	log.Info("Using SegregatedCacheProvider.")
	return &segregatedCacheProvider[P]{db: db, config: config}
}
