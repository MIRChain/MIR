package rawdb

import (
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/ethdb"
)

// HasBadBlock returns whether the block with the hash is a bad block. dep: Istanbul
func HasBadBlock[P crypto.PublicKey](db ethdb.Reader, hash common.Hash) bool {
	return ReadBadBlock[P](db, hash) != nil
}
