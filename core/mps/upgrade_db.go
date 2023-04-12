package mps

import (
	"fmt"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/core/privatecache"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
)

// chainReader contains methods to access local blockchain
type chainReader [P crypto.PublicKey] interface {
	consensus.ChainReader[P]
	CurrentBlock() *types.Block[P]
	GetReceiptsByHash(hash common.Hash) types.Receipts[P]
}

// UpgradeDB performs the following database operations to enable MPS support
// 1. Construct and persist Empty Private State with empty accounts
// 2. Construct and persist trie of root hashes of existing private states
// 3. Update new mapping: block header root -> trie of private states root
// 4. Once upgrade is complete update the ChainConfig.isMPS to true
func UpgradeDB[P crypto.PublicKey](db ethdb.Database, chain chainReader[P]) error {
	currentBlockNumber := uint64(chain.CurrentBlock().Number().Int64())
	genesisHeader := chain.GetHeaderByNumber(0)

	privateStatesTrieRoot := rawdb.GetPrivateStatesTrieRoot(db, genesisHeader.Root)
	privateCacheProvider := privatecache.NewPrivateCacheProvider[P](db, nil, nil, false)
	mpsRepo, err := NewMultiplePrivateStateRepository[P](db, state.NewDatabase[P](db), privateStatesTrieRoot, privateCacheProvider)
	if err != nil {
		return err
	}
	emptyState, err := mpsRepo.DefaultState()
	if err != nil {
		return err
	}
	// pre-populate with dummy one as the state root is derived from block root hash
	privateState := &managedState[P]{}
	mpsRepo.managedStates[types.DefaultPrivateStateIdentifier] = privateState
	for idx := uint64(1); idx <= currentBlockNumber; idx++ {
		header := chain.GetHeaderByNumber(idx)
		// TODO consider periodic reports instead of logging about each block
		fmt.Printf("Processing block %v with hash %v\n", idx, header.Hash().Hex())
		block := chain.GetBlock(header.Hash(), header.Number.Uint64())
		// update Empty Private State
		receipts := chain.GetReceiptsByHash(header.Hash())
		receiptsUpdated := false
		for txIdx, tx := range block.Transactions() {
			if tx.IsPrivate() && tx.To() == nil {
				// this is a contract creation transaction
				receipt := receipts[txIdx]
				accountAddress := receipt.ContractAddress
				emptyState.CreateAccount(accountAddress)
				emptyState.SetNonce(accountAddress, 1)

				emptyReceipt := &types.Receipt[P]{
					PostState:         receipt.PostState,
					Status:            1,
					CumulativeGasUsed: receipt.CumulativeGasUsed,
					Bloom:             types.Bloom[P]{},
					Logs:              nil,
					TxHash:            receipt.TxHash,
					ContractAddress:   receipt.ContractAddress,
					GasUsed:           receipt.GasUsed,
					BlockHash:         receipt.BlockHash,
					BlockNumber:       receipt.BlockNumber,
					TransactionIndex:  receipt.TransactionIndex,
				}
				emptyReceipt.Bloom = types.CreateBloom(types.Receipts[P]{emptyReceipt})
				emptyReceipt.PSReceipts = map[types.PrivateStateIdentifier]*types.Receipt[P]{
					types.DefaultPrivateStateIdentifier: receipt,
					types.EmptyPrivateStateIdentifier:   emptyReceipt}
				receipts[txIdx] = emptyReceipt
				receiptsUpdated = true
			}
		}
		if receiptsUpdated {
			batch := db.NewBatch()
			rawdb.WriteReceipts(batch, block.Hash(), block.NumberU64(), receipts)
			err := batch.Write()
			if err != nil {
				return err
			}
		}
		// update trie of private state roots and new mapping with block root hash
		privateState.stateRootProviderFunc = func(_ bool) (common.Hash, error) {
			return rawdb.GetPrivateStateRoot(db, header.Root), nil
		}
		err = mpsRepo.CommitAndWrite(chain.Config().IsEIP158(block.Number()), block)
		if err != nil {
			return err
		}
	}
	// update isMPS in the chain config
	config := chain.Config()
	config.IsMPS = true
	rawdb.WriteChainConfig(db, rawdb.ReadCanonicalHash(db, 0), config)
	fmt.Printf("MPS DB upgrade finished successfully.\n")
	return nil
}
