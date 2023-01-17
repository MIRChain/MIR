package rawdb

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
)

// HasBadBlock returns whether the block with the hash is a bad block. dep: Istanbul
func HasBadBlock(db ethdb.Reader, hash common.Hash) bool {
	return ReadBadBlock(db, hash) != nil
}
