package privatecache

import (
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/core/state"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/ethdb"
	"github.com/MIRChain/MIR/log"
	"github.com/MIRChain/MIR/trie"
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
