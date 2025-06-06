// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

// PublicEthereumAPI provides an API to access Ethereum full node-related
// information.
type PublicEthereumAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	e *Ethereum[T,P]
}

// NewPublicEthereumAPI creates a new Ethereum protocol API for full nodes.
func NewPublicEthereumAPI[T crypto.PrivateKey, P crypto.PublicKey](e *Ethereum[T,P]) *PublicEthereumAPI[T,P] {
	return &PublicEthereumAPI[T,P]{e}
}

// Etherbase is the address that mining rewards will be send to
func (api *PublicEthereumAPI[T,P]) Etherbase() (common.Address, error) {
	return api.e.Etherbase()
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (api *PublicEthereumAPI[T,P]) Coinbase() (common.Address, error) {
	return api.Etherbase()
}

// Hashrate returns the POW hashrate
func (api *PublicEthereumAPI[T,P]) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(api.e.Miner().Hashrate())
}

// PublicMinerAPI provides an API to control the miner.
// It offers only methods that operate on data that pose no security risk when it is publicly accessible.
type PublicMinerAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	e *Ethereum[T,P]
}

// NewPublicMinerAPI create a new PublicMinerAPI instance.
func NewPublicMinerAPI [T crypto.PrivateKey, P crypto.PublicKey](e *Ethereum[T,P]) *PublicMinerAPI[T,P] {
	return &PublicMinerAPI[T,P]{e}
}

// Mining returns an indication if this node is currently mining.
func (api *PublicMinerAPI[T,P]) Mining() bool {
	return api.e.IsMining()
}

// PrivateMinerAPI provides private RPC methods to control the miner.
// These methods can be abused by external users and must be considered insecure for use by untrusted users.
type PrivateMinerAPI  [T crypto.PrivateKey, P crypto.PublicKey] struct {
	e *Ethereum[T,P]
}

// NewPrivateMinerAPI create a new RPC service which controls the miner of this node.
func NewPrivateMinerAPI [T crypto.PrivateKey, P crypto.PublicKey] (e *Ethereum[T,P]) *PrivateMinerAPI[T,P] {
	return &PrivateMinerAPI[T,P]{e: e}
}

// Start starts the miner with the given number of threads. If threads is nil,
// the number of workers started is equal to the number of logical CPUs that are
// usable by this process. If mining is already running, this method adjust the
// number of threads allowed to use and updates the minimum price required by the
// transaction pool.
func (api *PrivateMinerAPI[T,P]) Start(threads *int) error {
	if threads == nil {
		return api.e.StartMining(runtime.NumCPU())
	}
	return api.e.StartMining(*threads)
}

// Stop terminates the miner, both at the consensus engine level as well as at
// the block creation level.
func (api *PrivateMinerAPI[T,P]) Stop() {
	api.e.StopMining()
}

// SetExtra sets the extra data string that is included when this miner mines a block.
func (api *PrivateMinerAPI[T,P]) SetExtra(extra string) (bool, error) {
	if err := api.e.Miner().SetExtra([]byte(extra)); err != nil {
		return false, err
	}
	return true, nil
}

// SetGasPrice sets the minimum accepted gas price for the miner.
func (api *PrivateMinerAPI[T,P]) SetGasPrice(gasPrice hexutil.Big) bool {
	api.e.lock.Lock()
	api.e.gasPrice = (*big.Int)(&gasPrice)
	api.e.lock.Unlock()

	api.e.txPool.SetGasPrice((*big.Int)(&gasPrice))
	return true
}

// SetEtherbase sets the etherbase of the miner
func (api *PrivateMinerAPI[T,P]) SetEtherbase(etherbase common.Address) bool {
	// Quorum: Set return value, so user can be notified if it is disallowed.
	return api.e.SetEtherbase(etherbase)
}

// SetRecommitInterval updates the interval for miner sealing work recommitting.
func (api *PrivateMinerAPI[T,P]) SetRecommitInterval(interval int) {
	api.e.Miner().SetRecommitInterval(time.Duration(interval) * time.Millisecond)
}

// PrivateAdminAPI is the collection of Ethereum full node-related APIs
// exposed over the private admin endpoint.
type PrivateAdminAPI  [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	eth *Ethereum[T,P]
}

// NewPrivateAdminAPI creates a new API definition for the full node private
// admin methods of the Ethereum service.
func NewPrivateAdminAPI[T crypto.PrivateKey, P crypto.PublicKey](eth *Ethereum[T,P]) *PrivateAdminAPI[T,P] {
	return &PrivateAdminAPI[T,P]{eth: eth}
}

// ExportChain exports the current blockchain into a local file,
// or a range of blocks if first and last are non-nil
func (api *PrivateAdminAPI[T,P]) ExportChain(file string, first *uint64, last *uint64) (bool, error) {
	if first == nil && last != nil {
		return false, errors.New("last cannot be specified without first")
	}
	if first != nil && last == nil {
		head := api.eth.BlockChain().CurrentHeader().Number.Uint64()
		last = &head
	}
	if _, err := os.Stat(file); err == nil {
		// File already exists. Allowing overwrite could be a DoS vecotor,
		// since the 'file' may point to arbitrary paths on the drive
		return false, errors.New("location would overwrite an existing file")
	}
	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if first != nil {
		if err := api.eth.BlockChain().ExportN(writer, *first, *last); err != nil {
			return false, err
		}
	} else if err := api.eth.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}

func hasAllBlocks[P crypto.PublicKey](chain *core.BlockChain[P], bs []*types.Block[P]) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
			return false
		}
	}

	return true
}

// ImportChain imports a blockchain from a local file.
func (api *PrivateAdminAPI[T,P]) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}

	// Run actual the import in pre-configured batches
	stream := rlp.NewStream(reader, 0)

	blocks, index := make([]*types.Block[P], 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(types.Block[P])
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.eth.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.eth.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}

// PublicDebugAPI is the collection of Ethereum full node APIs exposed
// over the public debugging endpoint.
type PublicDebugAPI [T crypto.PrivateKey, P crypto.PublicKey]struct {
	eth *Ethereum[T,P]
}

// NewPublicDebugAPI creates a new API definition for the full node-
// related public debug methods of the Ethereum service.
func NewPublicDebugAPI[T crypto.PrivateKey, P crypto.PublicKey](eth *Ethereum[T,P]) *PublicDebugAPI[T,P] {
	return &PublicDebugAPI[T,P]{eth: eth}
}

// DumpBlock retrieves the entire state of the database at a given block.
// Quorum adds an additional parameter to support private state dump
func (api *PublicDebugAPI[T,P]) DumpBlock(ctx context.Context, blockNr rpc.BlockNumber, typ *string) (state.Dump, error) {
	publicState, privateState, err := api.getStateDbsFromBlockNumber(ctx, blockNr)
	if err != nil {
		return state.Dump{}, err
	}

	if typ != nil && *typ == "private" {
		return privateState.RawDump(false, false, true), nil
	}
	return publicState.RawDump(false, false, true), nil
}

func (api *PublicDebugAPI[T,P]) PrivateStateRoot(ctx context.Context, blockNr rpc.BlockNumber) (common.Hash, error) {
	_, privateState, err := api.getStateDbsFromBlockNumber(ctx, blockNr)
	if err != nil {
		return common.Hash{}, err
	}
	return privateState.IntermediateRoot(true), nil
}

func (api *PublicDebugAPI[T,P]) DefaultStateRoot(ctx context.Context, blockNr rpc.BlockNumber) (common.Hash, error) {
	psm, err := api.eth.blockchain.PrivateStateManager().StateRepository(api.eth.blockchain.CurrentBlock().Hash())
	if err != nil {
		return common.Hash{}, err
	}
	defaultPSM := psm.DefaultStateMetadata()

	var privateState *state.StateDB[P]
	if blockNr == rpc.PendingBlockNumber {
		// If we're dumping the pending state, we need to request
		// both the pending block as well as the pending state from
		// the miner and operate on those
		_, _, privateState = api.eth.miner.Pending(defaultPSM.ID)
		// this is a copy of the private state so it is OK to do IntermediateRoot
		return privateState.IntermediateRoot(true), nil
	}
	var block *types.Block[P]
	if blockNr == rpc.LatestBlockNumber {
		block = api.eth.blockchain.CurrentBlock()
	} else {
		block = api.eth.blockchain.GetBlockByNumber(uint64(blockNr))
	}
	if block == nil {
		return common.Hash{}, fmt.Errorf("block #%d not found", blockNr)
	}
	_, privateState, err = api.eth.BlockChain().StateAtPSI(block.Root(), defaultPSM.ID)
	if err != nil {
		return common.Hash{}, err
	}
	return privateState.IntermediateRoot(true), nil
}

// PrivateDebugAPI is the collection of Ethereum full node APIs exposed over
// the private debugging endpoint.
type PrivateDebugAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	eth *Ethereum[T,P]
}

// NewPrivateDebugAPI creates a new API definition for the full node-related
// private debug methods of the Ethereum service.
func NewPrivateDebugAPI[T crypto.PrivateKey, P crypto.PublicKey](eth *Ethereum[T,P]) *PrivateDebugAPI[T,P] {
	return &PrivateDebugAPI[T,P]{eth: eth}
}

// Preimage is a debug API function that returns the preimage for a sha3 hash, if known.
func (api *PrivateDebugAPI[T,P]) Preimage(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	if preimage := rawdb.ReadPreimage(api.eth.ChainDb(), hash); preimage != nil {
		return preimage, nil
	}
	return nil, errors.New("unknown preimage")
}

// BadBlockArgs represents the entries in the list returned when bad blocks are queried.
type BadBlockArgs struct {
	Hash  common.Hash            `json:"hash"`
	Block map[string]interface{} `json:"block"`
	RLP   string                 `json:"rlp"`
}

// GetBadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
// and returns them as a JSON list of block-hashes
func (api *PrivateDebugAPI[T,P]) GetBadBlocks(ctx context.Context) ([]*BadBlockArgs, error) {
	var (
		err     error
		blocks  = rawdb.ReadAllBadBlocks[P](api.eth.chainDb)
		results = make([]*BadBlockArgs, 0, len(blocks))
	)
	for _, block := range blocks {
		var (
			blockRlp  string
			blockJSON map[string]interface{}
		)
		if rlpBytes, err := rlp.EncodeToBytes(block); err != nil {
			blockRlp = err.Error() // Hacky, but hey, it works
		} else {
			blockRlp = fmt.Sprintf("0x%x", rlpBytes)
		}
		if blockJSON, err = ethapi.RPCMarshalBlock(block, true, true); err != nil {
			blockJSON = map[string]interface{}{"error": err.Error()}
		}
		results = append(results, &BadBlockArgs{
			Hash:  block.Hash(),
			RLP:   blockRlp,
			Block: blockJSON,
		})
	}
	return results, nil
}

// AccountRangeMaxResults is the maximum number of results to be returned per call
const AccountRangeMaxResults = 256

// AccountRange enumerates all accounts in the given block and start point in paging request
func (api *PublicDebugAPI[T,P]) AccountRange(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash, start []byte, maxResults int, nocode, nostorage, incompletes bool) (state.IteratorDump, error) {
	psm, err := api.eth.blockchain.PrivateStateManager().ResolveForUserContext(ctx)
	if err != nil {
		return state.IteratorDump{}, err
	}
	var stateDb *state.StateDB[P]

	if number, ok := blockNrOrHash.Number(); ok {
		if number == rpc.PendingBlockNumber {
			// If we're dumping the pending state, we need to request
			// both the pending block as well as the pending state from
			// the miner and operate on those
			_, stateDb, _ = api.eth.miner.Pending(psm.ID)
		} else {
			var block *types.Block[P]
			if number == rpc.LatestBlockNumber {
				block = api.eth.blockchain.CurrentBlock()
			} else {
				block = api.eth.blockchain.GetBlockByNumber(uint64(number))
			}
			if block == nil {
				return state.IteratorDump{}, fmt.Errorf("block #%d not found", number)
			}
			stateDb, _, err = api.eth.BlockChain().StateAtPSI(block.Root(), psm.ID)
			if err != nil {
				return state.IteratorDump{}, err
			}
		}
	} else if hash, ok := blockNrOrHash.Hash(); ok {
		block := api.eth.blockchain.GetBlockByHash(hash)
		if block == nil {
			return state.IteratorDump{}, fmt.Errorf("block %s not found", hash.Hex())
		}
		stateDb, _, err = api.eth.BlockChain().StateAtPSI(block.Root(), psm.ID)
		if err != nil {
			return state.IteratorDump{}, err
		}
	} else {
		return state.IteratorDump{}, errors.New("either block number or block hash must be specified")
	}

	if maxResults > AccountRangeMaxResults || maxResults <= 0 {
		maxResults = AccountRangeMaxResults
	}
	return stateDb.IteratorDump(nocode, nostorage, incompletes, start, maxResults), nil
}

// StorageRangeResult is the result of a debug_storageRangeAt API call.
type StorageRangeResult struct {
	Storage storageMap   `json:"storage"`
	NextKey *common.Hash `json:"nextKey"` // nil if Storage includes the last key in the trie.
}

type storageMap map[common.Hash]storageEntry

type storageEntry struct {
	Key   *common.Hash `json:"key"`
	Value common.Hash  `json:"value"`
}

// StorageRangeAt returns the storage at the given block height and transaction index.
func (api *PrivateDebugAPI[T,P]) StorageRangeAt(ctx context.Context, blockHash common.Hash, txIndex int, contractAddress common.Address, keyStart hexutil.Bytes, maxResult int) (StorageRangeResult, error) {
	// Retrieve the block
	block := api.eth.blockchain.GetBlockByHash(blockHash)
	if block == nil {
		return StorageRangeResult{}, fmt.Errorf("block %#x not found", blockHash)
	}
	_, _, statedb, _, _, err := api.eth.stateAtTransaction(ctx, block, txIndex, 0)
	if err != nil {
		return StorageRangeResult{}, err
	}
	st := statedb.StorageTrie(contractAddress)
	if st == nil {
		return StorageRangeResult{}, fmt.Errorf("account %x doesn't exist", contractAddress)
	}
	return storageRangeAt(st, keyStart, maxResult)
}

func storageRangeAt(st state.Trie, start []byte, maxResult int) (StorageRangeResult, error) {
	it := trie.NewIterator(st.NodeIterator(start))
	result := StorageRangeResult{Storage: storageMap{}}
	for i := 0; i < maxResult && it.Next(); i++ {
		_, content, _, err := rlp.Split(it.Value)
		if err != nil {
			return StorageRangeResult{}, err
		}
		e := storageEntry{Value: common.BytesToHash(content)}
		if preimage := st.GetKey(it.Key); preimage != nil {
			preimage := common.BytesToHash(preimage)
			e.Key = &preimage
		}
		result.Storage[common.BytesToHash(it.Key)] = e
	}
	// Add the 'next key' so clients can continue downloading.
	if it.Next() {
		next := common.BytesToHash(it.Key)
		result.NextKey = &next
	}
	return result, nil
}

// GetModifiedAccountsByNumber returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI[T,P]) GetModifiedAccountsByNumber(startNum uint64, endNum *uint64) ([]common.Address, error) {
	var startBlock, endBlock *types.Block[P]

	startBlock = api.eth.blockchain.GetBlockByNumber(startNum)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startNum)
	}

	if endNum == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByNumber(*endNum)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %d not found", *endNum)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

// GetModifiedAccountsByHash returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI[T,P]) GetModifiedAccountsByHash(startHash common.Hash, endHash *common.Hash) ([]common.Address, error) {
	var startBlock, endBlock *types.Block[P]
	startBlock = api.eth.blockchain.GetBlockByHash(startHash)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startHash)
	}

	if endHash == nil {
		endBlock = startBlock
		startBlock = api.eth.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.eth.blockchain.GetBlockByHash(*endHash)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %x not found", *endHash)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

func (api *PrivateDebugAPI[T,P]) getModifiedAccounts(startBlock, endBlock *types.Block[P]) ([]common.Address, error) {
	if startBlock.Number().Uint64() >= endBlock.Number().Uint64() {
		return nil, fmt.Errorf("start block height (%d) must be less than end block height (%d)", startBlock.Number().Uint64(), endBlock.Number().Uint64())
	}
	triedb := api.eth.BlockChain().StateCache().TrieDB()

	oldTrie, err := trie.NewSecure[P](startBlock.Root(), triedb)
	if err != nil {
		return nil, err
	}
	newTrie, err := trie.NewSecure[P](endBlock.Root(), triedb)
	if err != nil {
		return nil, err
	}
	diff, _ := trie.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := trie.NewIterator(diff)

	var dirty []common.Address
	for iter.Next() {
		key := newTrie.GetKey(iter.Key)
		if key == nil {
			return nil, fmt.Errorf("no preimage found for hash %x", iter.Key)
		}
		dirty = append(dirty, common.BytesToAddress(key))
	}
	return dirty, nil
}

// Quorum

// StorageRoot returns the storage root of an account on the the given (optional) block height.
// If block number is not given the latest block is used.
func (s *PublicEthereumAPI[T,P]) StorageRoot(ctx context.Context, addr common.Address, blockNr *rpc.BlockNumber) (common.Hash, error) {
	var (
		pub, priv *state.StateDB[P]
		err       error
	)

	psm, err := s.e.blockchain.PrivateStateManager().ResolveForUserContext(ctx)
	if err != nil {
		return common.Hash{}, err
	}
	if blockNr == nil || blockNr.Int64() == rpc.LatestBlockNumber.Int64() {
		pub, priv, err = s.e.blockchain.StatePSI(psm.ID)
	} else {
		if ch := s.e.blockchain.GetHeaderByNumber(uint64(blockNr.Int64())); ch != nil {
			pub, priv, err = s.e.blockchain.StateAtPSI(ch.Root, psm.ID)
		} else {
			return common.Hash{}, fmt.Errorf("invalid block number")
		}
	}

	if err != nil {
		return common.Hash{}, err
	}

	if priv.Exist(addr) {
		return priv.GetStorageRoot(addr)
	}
	return pub.GetStorageRoot(addr)
}

// DumpAddress retrieves the state of an address at a given block.
// Quorum adds an additional parameter to support private state dump
func (api *PublicDebugAPI[T,P]) DumpAddress(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (state.DumpAccount, error) {
	publicState, privateState, err := api.getStateDbsFromBlockNumber(ctx, blockNr)
	if err != nil {
		return state.DumpAccount{}, err
	}

	if accountDump, ok := privateState.DumpAddress(address); ok {
		return accountDump, nil
	}
	if accountDump, ok := publicState.DumpAddress(address); ok {
		return accountDump, nil
	}
	return state.DumpAccount{}, errors.New("error retrieving state")
}

//Taken from DumpBlock, as it was reused in DumpAddress.
//Contains modifications from the original to return the private state db, as well as public.
func (api *PublicDebugAPI[T,P]) getStateDbsFromBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB[P], *state.StateDB[P], error) {
	psm, err := api.eth.blockchain.PrivateStateManager().ResolveForUserContext(ctx)
	if err != nil {
		return nil, nil, err
	}
	if blockNr == rpc.PendingBlockNumber {
		// If we're dumping the pending state, we need to request
		// both the pending block as well as the pending state from
		// the miner and operate on those
		_, publicState, privateState := api.eth.miner.Pending(psm.ID)
		return publicState, privateState, nil
	}

	var block *types.Block[P]
	if blockNr == rpc.LatestBlockNumber {
		block = api.eth.blockchain.CurrentBlock()
	} else {
		block = api.eth.blockchain.GetBlockByNumber(uint64(blockNr))
	}
	if block == nil {
		return nil, nil, fmt.Errorf("block #%d not found", blockNr)
	}
	return api.eth.BlockChain().StateAtPSI(block.Root(), psm.ID)
}
