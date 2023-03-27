package rawdb

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
)

// HasBadBlock returns whether the block with the hash is a bad block. dep: Istanbul
func HasBadBlock[P crypto.PublicKey](db ethdb.Reader, hash common.Hash) bool {
	return ReadBadBlock[P](db, hash) != nil
}
