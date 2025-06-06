// Copyright 2018 The go-ethereum Authors
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

package rawdb

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"sort"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// ReadCanonicalHash retrieves the hash assigned to a canonical block number.
func ReadCanonicalHash(db ethdb.Reader, number uint64) common.Hash {
	data, _ := db.Ancient(freezerHashTable, number)
	if len(data) == 0 {
		data, _ = db.Get(headerHashKey(number))
		// In the background freezer is moving data from leveldb to flatten files.
		// So during the first check for ancient db, the data is not yet in there,
		// but when we reach into leveldb, the data was already moved. That would
		// result in a not found error.
		if len(data) == 0 {
			data, _ = db.Ancient(freezerHashTable, number)
		}
	}
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteCanonicalHash stores the hash assigned to a canonical block number.
func WriteCanonicalHash(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	if err := db.Put(headerHashKey(number), hash.Bytes()); err != nil {
		log.Crit("Failed to store number to hash mapping", "err", err)
	}
}

// DeleteCanonicalHash removes the number to hash canonical mapping.
func DeleteCanonicalHash(db ethdb.KeyValueWriter, number uint64) {
	if err := db.Delete(headerHashKey(number)); err != nil {
		log.Crit("Failed to delete number to hash mapping", "err", err)
	}
}

// ReadAllHashes retrieves all the hashes assigned to blocks at a certain heights,
// both canonical and reorged forks included.
func ReadAllHashes(db ethdb.Iteratee, number uint64) []common.Hash {
	prefix := headerKeyPrefix(number)

	hashes := make([]common.Hash, 0, 1)
	it := db.NewIterator(prefix, nil)
	defer it.Release()

	for it.Next() {
		if key := it.Key(); len(key) == len(prefix)+32 {
			hashes = append(hashes, common.BytesToHash(key[len(key)-32:]))
		}
	}
	return hashes
}

// ReadAllCanonicalHashes retrieves all canonical number and hash mappings at the
// certain chain range. If the accumulated entries reaches the given threshold,
// abort the iteration and return the semi-finish result.
func ReadAllCanonicalHashes(db ethdb.Iteratee, from uint64, to uint64, limit int) ([]uint64, []common.Hash) {
	// Short circuit if the limit is 0.
	if limit == 0 {
		return nil, nil
	}
	var (
		numbers []uint64
		hashes  []common.Hash
	)
	// Construct the key prefix of start point.
	start, end := headerHashKey(from), headerHashKey(to)
	it := db.NewIterator(nil, start)
	defer it.Release()

	for it.Next() {
		if bytes.Compare(it.Key(), end) >= 0 {
			break
		}
		if key := it.Key(); len(key) == len(headerPrefix)+8+1 && bytes.Equal(key[len(key)-1:], headerHashSuffix) {
			numbers = append(numbers, binary.BigEndian.Uint64(key[len(headerPrefix):len(headerPrefix)+8]))
			hashes = append(hashes, common.BytesToHash(it.Value()))
			// If the accumulated entries reaches the limit threshold, return.
			if len(numbers) >= limit {
				break
			}
		}
	}
	return numbers, hashes
}

// ReadHeaderNumber returns the header number assigned to a hash.
func ReadHeaderNumber(db ethdb.KeyValueReader, hash common.Hash) *uint64 {
	data, _ := db.Get(headerNumberKey(hash))
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// WriteHeaderNumber stores the hash->number mapping.
func WriteHeaderNumber(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	key := headerNumberKey(hash)
	enc := encodeBlockNumber(number)
	if err := db.Put(key, enc); err != nil {
		log.Crit("Failed to store hash to number mapping", "err", err)
	}
}

// DeleteHeaderNumber removes hash->number mapping.
func DeleteHeaderNumber(db ethdb.KeyValueWriter, hash common.Hash) {
	if err := db.Delete(headerNumberKey(hash)); err != nil {
		log.Crit("Failed to delete hash to number mapping", "err", err)
	}
}

// ReadHeadHeaderHash retrieves the hash of the current canonical head header.
func ReadHeadHeaderHash(db ethdb.KeyValueReader) common.Hash {
	data, _ := db.Get(headHeaderKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadHeaderHash stores the hash of the current canonical head header.
func WriteHeadHeaderHash(db ethdb.KeyValueWriter, hash common.Hash) {
	if err := db.Put(headHeaderKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last header's hash", "err", err)
	}
}

// ReadHeadBlockHash retrieves the hash of the current canonical head block.
func ReadHeadBlockHash(db ethdb.KeyValueReader) common.Hash {
	data, _ := db.Get(headBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadBlockHash stores the head block's hash.
func WriteHeadBlockHash(db ethdb.KeyValueWriter, hash common.Hash) {
	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last block's hash", "err", err)
	}
}

// ReadHeadFastBlockHash retrieves the hash of the current fast-sync head block.
func ReadHeadFastBlockHash(db ethdb.KeyValueReader) common.Hash {
	data, _ := db.Get(headFastBlockKey)
	if len(data) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(data)
}

// WriteHeadFastBlockHash stores the hash of the current fast-sync head block.
func WriteHeadFastBlockHash(db ethdb.KeyValueWriter, hash common.Hash) {
	if err := db.Put(headFastBlockKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last fast block's hash", "err", err)
	}
}

// ReadLastPivotNumber retrieves the number of the last pivot block. If the node
// full synced, the last pivot will always be nil.
func ReadLastPivotNumber(db ethdb.KeyValueReader) *uint64 {
	data, _ := db.Get(lastPivotKey)
	if len(data) == 0 {
		return nil
	}
	var pivot uint64
	if err := rlp.DecodeBytes(data, &pivot); err != nil {
		log.Error("Invalid pivot block number in database", "err", err)
		return nil
	}
	return &pivot
}

// WriteLastPivotNumber stores the number of the last pivot block.
func WriteLastPivotNumber(db ethdb.KeyValueWriter, pivot uint64) {
	enc, err := rlp.EncodeToBytes(pivot)
	if err != nil {
		log.Crit("Failed to encode pivot block number", "err", err)
	}
	if err := db.Put(lastPivotKey, enc); err != nil {
		log.Crit("Failed to store pivot block number", "err", err)
	}
}

// ReadFastTrieProgress retrieves the number of tries nodes fast synced to allow
// reporting correct numbers across restarts.
func ReadFastTrieProgress(db ethdb.KeyValueReader) uint64 {
	data, _ := db.Get(fastTrieProgressKey)
	if len(data) == 0 {
		return 0
	}
	return new(big.Int).SetBytes(data).Uint64()
}

// WriteFastTrieProgress stores the fast sync trie process counter to support
// retrieving it across restarts.
func WriteFastTrieProgress(db ethdb.KeyValueWriter, count uint64) {
	if err := db.Put(fastTrieProgressKey, new(big.Int).SetUint64(count).Bytes()); err != nil {
		log.Crit("Failed to store fast sync trie progress", "err", err)
	}
}

// ReadTxIndexTail retrieves the number of oldest indexed block
// whose transaction indices has been indexed. If the corresponding entry
// is non-existent in database it means the indexing has been finished.
func ReadTxIndexTail(db ethdb.KeyValueReader) *uint64 {
	data, _ := db.Get(txIndexTailKey)
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// WriteTxIndexTail stores the number of oldest indexed block
// into database.
func WriteTxIndexTail(db ethdb.KeyValueWriter, number uint64) {
	if err := db.Put(txIndexTailKey, encodeBlockNumber(number)); err != nil {
		log.Crit("Failed to store the transaction index tail", "err", err)
	}
}

// ReadFastTxLookupLimit retrieves the tx lookup limit used in fast sync.
func ReadFastTxLookupLimit(db ethdb.KeyValueReader) *uint64 {
	data, _ := db.Get(fastTxLookupLimitKey)
	if len(data) != 8 {
		return nil
	}
	number := binary.BigEndian.Uint64(data)
	return &number
}

// WriteFastTxLookupLimit stores the txlookup limit used in fast sync into database.
func WriteFastTxLookupLimit(db ethdb.KeyValueWriter, number uint64) {
	if err := db.Put(fastTxLookupLimitKey, encodeBlockNumber(number)); err != nil {
		log.Crit("Failed to store transaction lookup limit for fast sync", "err", err)
	}
}

// Quorum
// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
func readHeaderRLP[P crypto.PublicKey](db ethdb.Reader, hash common.Hash, number uint64) (rlp.RawValue, *types.Header[P]) {
	// First try to look up the data in ancient database. Extra hash
	// comparison is necessary since ancient database only maintains
	// the canonical data.
	data, _ := db.Ancient(freezerHeaderTable, number)

	// Quorum: parse header to make sure we compare using the right hash (IBFT hash is based on a filtered header)
	if len(data) > 0 {
		header := decodeHeaderRLP[P](data)
		if header.Hash() == hash {
			return data, header
		}
	}
	// End Quorum

	// Then try to look up the data in leveldb.
	data, _ = db.Get(headerKey(number, hash))
	if len(data) > 0 {
		return data, decodeHeaderRLP[P](data) // Quorum: return decodeHeaderRLP(data)
	}
	// In the background freezer is moving data from leveldb to flatten files.
	// So during the first check for ancient db, the data is not yet in there,
	// but when we reach into leveldb, the data was already moved. That would
	// result in a not found error.
	data, _ = db.Ancient(freezerHeaderTable, number)

	// Quorum: parse header to make sure we compare using the right hash (IBFT hash is based on a filtered header)
	if len(data) > 0 {
		header := decodeHeaderRLP[P](data)
		if header.Hash() == hash {
			return data, header
		}
	}
	// End Quorum

	return nil, nil // Can't find the data anywhere.
}

// Quorum
func decodeHeaderRLP[P crypto.PublicKey](data rlp.RawValue) *types.Header[P] {
	header := new(types.Header[P])
	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
		log.Error("Invalid block header RLP", "err", err)
		return nil
	}
	return header
}

// ReadHeaderRLP retrieves a block header in its raw RLP database encoding.
func ReadHeaderRLP[P crypto.PublicKey](db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	// Quorum: original code implemented inside `readHeaderRLP(...)` with some modifications from Quorum
	data, _ := readHeaderRLP[P](db, hash, number)
	return data
}

// HasHeader verifies the existence of a block header corresponding to the hash.
func HasHeader(db ethdb.Reader, hash common.Hash, number uint64) bool {
	if has, err := db.Ancient(freezerHashTable, number); err == nil && common.BytesToHash(has) == hash {
		return true
	}
	if has, err := db.Has(headerKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadHeader retrieves the block header corresponding to the hash.
func ReadHeader[P crypto.PublicKey](db ethdb.Reader, hash common.Hash, number uint64) *types.Header[P] {
	data, header := readHeaderRLP[P](db, hash, number)
	if data == nil {
		log.Trace("header data not found in ancient or level db", "hash", hash)
		return nil
	}
	if header == nil {
		log.Error("Invalid block header RLP", "hash", hash)
		return nil
	}
	return header
}

// WriteHeader stores a block header into the database and also stores the hash-
// to-number mapping.
func WriteHeader[P crypto.PublicKey](db ethdb.KeyValueWriter, header *types.Header[P]) {
	var (
		hash   = header.Hash()
		number = header.Number.Uint64()
	)
	// Write the hash -> number mapping
	WriteHeaderNumber(db, hash, number)

	// Write the encoded header
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		log.Crit("Failed to RLP encode header", "err", err)
	}
	key := headerKey(number, hash)
	if err := db.Put(key, data); err != nil {
		log.Crit("Failed to store header", "err", err)
	}
}

// DeleteHeader removes all block header data associated with a hash.
func DeleteHeader(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	deleteHeaderWithoutNumber(db, hash, number)
	if err := db.Delete(headerNumberKey(hash)); err != nil {
		log.Crit("Failed to delete hash to number mapping", "err", err)
	}
}

// deleteHeaderWithoutNumber removes only the block header but does not remove
// the hash to number mapping.
func deleteHeaderWithoutNumber(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	if err := db.Delete(headerKey(number, hash)); err != nil {
		log.Crit("Failed to delete header", "err", err)
	}
}

// ReadBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
func ReadBodyRLP(db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	// First try to look up the data in ancient database. Extra hash
	// comparison is necessary since ancient database only maintains
	// the canonical data.
	data, _ := db.Ancient(freezerBodiesTable, number)
	if len(data) > 0 {
		h, _ := db.Ancient(freezerHashTable, number)
		if common.BytesToHash(h) == hash {
			return data
		}
	}
	// Then try to look up the data in leveldb.
	data, _ = db.Get(blockBodyKey(number, hash))
	if len(data) > 0 {
		return data
	}
	// In the background freezer is moving data from leveldb to flatten files.
	// So during the first check for ancient db, the data is not yet in there,
	// but when we reach into leveldb, the data was already moved. That would
	// result in a not found error.
	data, _ = db.Ancient(freezerBodiesTable, number)
	if len(data) > 0 {
		h, _ := db.Ancient(freezerHashTable, number)
		if common.BytesToHash(h) == hash {
			return data
		}
	}
	return nil // Can't find the data anywhere.
}

// ReadCanonicalBodyRLP retrieves the block body (transactions and uncles) for the canonical
// block at number, in RLP encoding.
func ReadCanonicalBodyRLP(db ethdb.Reader, number uint64) rlp.RawValue {
	// If it's an ancient one, we don't need the canonical hash
	data, _ := db.Ancient(freezerBodiesTable, number)
	if len(data) == 0 {
		// Need to get the hash
		data, _ = db.Get(blockBodyKey(number, ReadCanonicalHash(db, number)))
		// In the background freezer is moving data from leveldb to flatten files.
		// So during the first check for ancient db, the data is not yet in there,
		// but when we reach into leveldb, the data was already moved. That would
		// result in a not found error.
		if len(data) == 0 {
			data, _ = db.Ancient(freezerBodiesTable, number)
		}
	}
	return data
}

// WriteBodyRLP stores an RLP encoded block body into the database.
func WriteBodyRLP(db ethdb.KeyValueWriter, hash common.Hash, number uint64, rlp rlp.RawValue) {
	if err := db.Put(blockBodyKey(number, hash), rlp); err != nil {
		log.Crit("Failed to store block body", "err", err)
	}
}

// HasBody verifies the existence of a block body corresponding to the hash.
func HasBody(db ethdb.Reader, hash common.Hash, number uint64) bool {
	if has, err := db.Ancient(freezerHashTable, number); err == nil && common.BytesToHash(has) == hash {
		return true
	}
	if has, err := db.Has(blockBodyKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadBody retrieves the block body corresponding to the hash.
func ReadBody[P crypto.PublicKey](db ethdb.Reader, hash common.Hash, number uint64) *types.Body[P] {
	data := ReadBodyRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	body := new(types.Body[P])
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		log.Error("Invalid block body RLP", "hash", hash, "err", err)
		return nil
	}
	return body
}

// WriteBody stores a block body into the database.
func WriteBody[P crypto.PublicKey](db ethdb.KeyValueWriter, hash common.Hash, number uint64, body *types.Body[P]) {
	data, err := rlp.EncodeToBytes(body)
	if err != nil {
		log.Crit("Failed to RLP encode body", "err", err)
	}
	WriteBodyRLP(db, hash, number, data)
}

// DeleteBody removes all block body data associated with a hash.
func DeleteBody(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	if err := db.Delete(blockBodyKey(number, hash)); err != nil {
		log.Crit("Failed to delete block body", "err", err)
	}
}

// ReadTdRLP retrieves a block's total difficulty corresponding to the hash in RLP encoding.
func ReadTdRLP(db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	// First try to look up the data in ancient database. Extra hash
	// comparison is necessary since ancient database only maintains
	// the canonical data.
	data, _ := db.Ancient(freezerDifficultyTable, number)
	if len(data) > 0 {
		h, _ := db.Ancient(freezerHashTable, number)
		if common.BytesToHash(h) == hash {
			return data
		}
	}
	// Then try to look up the data in leveldb.
	data, _ = db.Get(headerTDKey(number, hash))
	if len(data) > 0 {
		return data
	}
	// In the background freezer is moving data from leveldb to flatten files.
	// So during the first check for ancient db, the data is not yet in there,
	// but when we reach into leveldb, the data was already moved. That would
	// result in a not found error.
	data, _ = db.Ancient(freezerDifficultyTable, number)
	if len(data) > 0 {
		h, _ := db.Ancient(freezerHashTable, number)
		if common.BytesToHash(h) == hash {
			return data
		}
	}
	return nil // Can't find the data anywhere.
}

// ReadTd retrieves a block's total difficulty corresponding to the hash.
func ReadTd(db ethdb.Reader, hash common.Hash, number uint64) *big.Int {
	data := ReadTdRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	td := new(big.Int)
	if err := rlp.Decode(bytes.NewReader(data), td); err != nil {
		log.Error("Invalid block total difficulty RLP", "hash", hash, "err", err)
		return nil
	}
	return td
}

// WriteTd stores the total difficulty of a block into the database.
func WriteTd(db ethdb.KeyValueWriter, hash common.Hash, number uint64, td *big.Int) {
	data, err := rlp.EncodeToBytes(td)
	if err != nil {
		log.Crit("Failed to RLP encode block total difficulty", "err", err)
	}
	if err := db.Put(headerTDKey(number, hash), data); err != nil {
		log.Crit("Failed to store block total difficulty", "err", err)
	}
}

// DeleteTd removes all block total difficulty data associated with a hash.
func DeleteTd(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	if err := db.Delete(headerTDKey(number, hash)); err != nil {
		log.Crit("Failed to delete block total difficulty", "err", err)
	}
}

// HasReceipts verifies the existence of all the transaction receipts belonging
// to a block.
func HasReceipts(db ethdb.Reader, hash common.Hash, number uint64) bool {
	if has, err := db.Ancient(freezerHashTable, number); err == nil && common.BytesToHash(has) == hash {
		return true
	}
	if has, err := db.Has(blockReceiptsKey(number, hash)); !has || err != nil {
		return false
	}
	return true
}

// ReadReceiptsRLP retrieves all the transaction receipts belonging to a block in RLP encoding.
func ReadReceiptsRLP(db ethdb.Reader, hash common.Hash, number uint64) rlp.RawValue {
	// First try to look up the data in ancient database. Extra hash
	// comparison is necessary since ancient database only maintains
	// the canonical data.
	data, _ := db.Ancient(freezerReceiptTable, number)
	if len(data) > 0 {
		h, _ := db.Ancient(freezerHashTable, number)
		if common.BytesToHash(h) == hash {
			return data
		}
	}
	// Then try to look up the data in leveldb.
	data, _ = db.Get(blockReceiptsKey(number, hash))
	if len(data) > 0 {
		return data
	}
	// In the background freezer is moving data from leveldb to flatten files.
	// So during the first check for ancient db, the data is not yet in there,
	// but when we reach into leveldb, the data was already moved. That would
	// result in a not found error.
	data, _ = db.Ancient(freezerReceiptTable, number)
	if len(data) > 0 {
		h, _ := db.Ancient(freezerHashTable, number)
		if common.BytesToHash(h) == hash {
			return data
		}
	}
	return nil // Can't find the data anywhere.
}

// ReadRawReceipts retrieves all the transaction receipts belonging to a block.
// The receipt metadata fields are not guaranteed to be populated, so they
// should not be used. Use ReadReceipts instead if the metadata is needed.
func ReadRawReceipts[P crypto.PublicKey](db ethdb.Reader, hash common.Hash, number uint64) types.Receipts[P] {
	// Retrieve the flattened receipt slice
	data := ReadReceiptsRLP(db, hash, number)
	if len(data) == 0 {
		return nil
	}
	// split the data into the standard receipt rlp list and the quorum extraData bytes
	_, extraData, err := rlp.SplitList(data)
	if err != nil {
		log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
		return nil
	}
	// reslice data to remove extraData and get the receipt rlp list as the result from rlp.SplitList does not include the list header bytes
	vanillaDataWithListHeader := data[0 : len(data)-len(extraData)]

	// Convert the receipts from their storage form to their internal representation
	storageReceipts := []*types.ReceiptForStorage[P]{}
	if err := rlp.DecodeBytes(vanillaDataWithListHeader, &storageReceipts); err != nil {
		log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
		return nil
	}
	receipts := make(types.Receipts[P], len(storageReceipts))
	for i, storageReceipt := range storageReceipts {
		receipts[i] = (*types.Receipt[P])(storageReceipt)
	}
	if len(extraData) > 0 {
		quorumExtraDataReceipts := []*types.QuorumReceiptExtraData[P]{}
		if err := rlp.DecodeBytes(extraData, &quorumExtraDataReceipts); err != nil {
			log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
			return nil
		}
		for i, quorumExtraDataReceipt := range quorumExtraDataReceipts {
			if quorumExtraDataReceipt != nil {
				receipts[i].FillReceiptExtraDataFromStorage(quorumExtraDataReceipt)
			}
		}
	}
	return receipts
}

// ReadReceipts retrieves all the transaction receipts belonging to a block, including
// its correspoinding metadata fields. If it is unable to populate these metadata
// fields then nil is returned.
//
// The current implementation populates these metadata fields by reading the receipts'
// corresponding block body, so if the block body is not found it will return nil even
// if the receipt itself is stored.
func ReadReceipts[P crypto.PublicKey](db ethdb.Reader, hash common.Hash, number uint64, config *params.ChainConfig) types.Receipts[P] {
	// We're deriving many fields from the block body, retrieve beside the receipt
	receipts := ReadRawReceipts[P](db, hash, number)
	if receipts == nil {
		return nil
	}
	body := ReadBody[P](db, hash, number)
	if body == nil {
		log.Error("Missing body but have receipt", "hash", hash, "number", number)
		return nil
	}
	if err := receipts.DeriveFields(config, hash, number, body.Transactions); err != nil {
		log.Error("Failed to derive block receipts fields", "hash", hash, "number", number, "err", err)
		return nil
	}
	return receipts
}

// WriteReceipts stores all the transaction receipts belonging to a block.
func WriteReceipts[P crypto.PublicKey](db ethdb.KeyValueWriter, hash common.Hash, number uint64, receipts types.Receipts[P]) {
	// Convert the receipts into their storage form and serialize them
	storageReceipts := make([]*types.ReceiptForStorage[P], len(receipts))
	quorumReceiptsExtraData := make([]*types.QuorumReceiptExtraData[P], len(receipts))
	extraDataEmpty := true
	for i, receipt := range receipts {
		storageReceipts[i] = (*types.ReceiptForStorage[P])(receipt)
		quorumReceiptsExtraData[i] = &receipt.QuorumReceiptExtraData
		if !receipt.QuorumReceiptExtraData.IsEmpty() {
			extraDataEmpty = false
		}
	}
	bytes, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		log.Crit("Failed to encode block receipts", "err", err)
	}
	if !extraDataEmpty {
		bytesExtraData, err := rlp.EncodeToBytes(quorumReceiptsExtraData)
		if err != nil {
			log.Crit("Failed to encode block receipts", "err", err)
		}
		// the vanilla receipts and the extra data receipts are concatenated and stored as a single value
		bytes = append(bytes, bytesExtraData...)
	}
	// Store the flattened receipt slice
	if err := db.Put(blockReceiptsKey(number, hash), bytes); err != nil {
		log.Crit("Failed to store block receipts", "err", err)
	}
}

// DeleteReceipts removes all receipt data associated with a block hash.
func DeleteReceipts(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	if err := db.Delete(blockReceiptsKey(number, hash)); err != nil {
		log.Crit("Failed to delete block receipts", "err", err)
	}
}

// ReadBlock retrieves an entire block corresponding to the hash, assembling it
// back from the stored header and body. If either the header or body could not
// be retrieved nil is returned.
//
// Note, due to concurrent download of header and block body the header and thus
// canonical hash can be stored in the database but the body data not (yet).
func ReadBlock[P crypto.PublicKey](db ethdb.Reader, hash common.Hash, number uint64) *types.Block[P] {
	header := ReadHeader[P](db, hash, number)
	if header == nil {
		return nil
	}
	body := ReadBody[P](db, hash, number)
	if body == nil {
		return nil
	}
	return types.NewBlockWithHeader[P](header).WithBody(body.Transactions, body.Uncles)
}

// WriteBlock serializes a block into the database, header and body separately.
func WriteBlock[P crypto.PublicKey](db ethdb.KeyValueWriter, block *types.Block[P]) {
	WriteBody(db, block.Hash(), block.NumberU64(), block.Body())
	WriteHeader(db, block.Header())
}

// WriteAncientBlock writes entire block data into ancient store and returns the total written size.
func WriteAncientBlock[P crypto.PublicKey](db ethdb.AncientWriter, block *types.Block[P], receipts types.Receipts[P], td *big.Int) int {
	// Encode all block components to RLP format.
	headerBlob, err := rlp.EncodeToBytes(block.Header())
	if err != nil {
		log.Crit("Failed to RLP encode block header", "err", err)
	}
	bodyBlob, err := rlp.EncodeToBytes(block.Body())
	if err != nil {
		log.Crit("Failed to RLP encode body", "err", err)
	}
	storageReceipts := make([]*types.ReceiptForStorage[P], len(receipts))
	for i, receipt := range receipts {
		storageReceipts[i] = (*types.ReceiptForStorage[P])(receipt)
	}
	receiptBlob, err := rlp.EncodeToBytes(storageReceipts)
	if err != nil {
		log.Crit("Failed to RLP encode block receipts", "err", err)
	}
	tdBlob, err := rlp.EncodeToBytes(td)
	if err != nil {
		log.Crit("Failed to RLP encode block total difficulty", "err", err)
	}
	// Write all blob to flatten files.
	err = db.AppendAncient(block.NumberU64(), block.Hash().Bytes(), headerBlob, bodyBlob, receiptBlob, tdBlob)
	if err != nil {
		log.Crit("Failed to write block data to ancient store", "err", err)
	}
	return len(headerBlob) + len(bodyBlob) + len(receiptBlob) + len(tdBlob) + common.HashLength
}

// DeleteBlock removes all block data associated with a hash.
func DeleteBlock(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	DeleteReceipts(db, hash, number)
	DeleteHeader(db, hash, number)
	DeleteBody(db, hash, number)
	DeleteTd(db, hash, number)
}

// DeleteBlockWithoutNumber removes all block data associated with a hash, except
// the hash to number mapping.
func DeleteBlockWithoutNumber(db ethdb.KeyValueWriter, hash common.Hash, number uint64) {
	DeleteReceipts(db, hash, number)
	deleteHeaderWithoutNumber(db, hash, number)
	DeleteBody(db, hash, number)
	DeleteTd(db, hash, number)
}

const badBlockToKeep = 10

type badBlock [P crypto.PublicKey] struct {
	Header *types.Header[P]
	Body   *types.Body[P]
}

// badBlockList implements the sort interface to allow sorting a list of
// bad blocks by their number in the reverse order.
type badBlockList [P crypto.PublicKey] []*badBlock[P]

func (s badBlockList[P]) Len() int { return len(s) }
func (s badBlockList[P]) Less(i, j int) bool {
	return s[i].Header.Number.Uint64() < s[j].Header.Number.Uint64()
}
func (s badBlockList[P]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ReadBadBlock retrieves the bad block with the corresponding block hash.
func ReadBadBlock[P crypto.PublicKey](db ethdb.Reader, hash common.Hash) *types.Block[P] {
	blob, err := db.Get(badBlockKey)
	if err != nil {
		return nil
	}
	var badBlocks badBlockList[P]
	if err := rlp.DecodeBytes(blob, &badBlocks); err != nil {
		return nil
	}
	for _, bad := range badBlocks {
		if bad.Header.Hash() == hash {
			return types.NewBlockWithHeader[P](bad.Header).WithBody(bad.Body.Transactions, bad.Body.Uncles)
		}
	}
	return nil
}

// ReadAllBadBlocks retrieves all the bad blocks in the database.
// All returned blocks are sorted in reverse order by number.
func ReadAllBadBlocks[P crypto.PublicKey](db ethdb.Reader) []*types.Block[P] {
	blob, err := db.Get(badBlockKey)
	if err != nil {
		return nil
	}
	var badBlocks badBlockList[P]
	if err := rlp.DecodeBytes(blob, &badBlocks); err != nil {
		return nil
	}
	var blocks []*types.Block[P]
	for _, bad := range badBlocks {
		blocks = append(blocks, types.NewBlockWithHeader[P](bad.Header).WithBody(bad.Body.Transactions, bad.Body.Uncles))
	}
	return blocks
}

// WriteBadBlock serializes the bad block into the database. If the cumulated
// bad blocks exceeds the limitation, the oldest will be dropped.
func WriteBadBlock[P crypto.PublicKey](db ethdb.KeyValueStore, block *types.Block[P]) {
	blob, err := db.Get(badBlockKey)
	if err != nil {
		log.Warn("Failed to load old bad blocks", "error", err)
	}
	var badBlocks badBlockList[P]
	if len(blob) > 0 {
		if err := rlp.DecodeBytes(blob, &badBlocks); err != nil {
			log.Crit("Failed to decode old bad blocks", "error", err)
		}
	}
	for _, b := range badBlocks {
		if b.Header.Number.Uint64() == block.NumberU64() && b.Header.Hash() == block.Hash() {
			log.Info("Skip duplicated bad block", "number", block.NumberU64(), "hash", block.Hash())
			return
		}
	}
	badBlocks = append(badBlocks, &badBlock[P]{
		Header: block.Header(),
		Body:   block.Body(),
	})
	sort.Sort(sort.Reverse(badBlocks))
	if len(badBlocks) > badBlockToKeep {
		badBlocks = badBlocks[:badBlockToKeep]
	}
	data, err := rlp.EncodeToBytes(badBlocks)
	if err != nil {
		log.Crit("Failed to encode bad blocks", "err", err)
	}
	if err := db.Put(badBlockKey, data); err != nil {
		log.Crit("Failed to write bad blocks", "err", err)
	}
}

// DeleteBadBlocks deletes all the bad blocks from the database
func DeleteBadBlocks(db ethdb.KeyValueWriter) {
	if err := db.Delete(badBlockKey); err != nil {
		log.Crit("Failed to delete bad blocks", "err", err)
	}
}

// FindCommonAncestor returns the last common ancestor of two block headers
func FindCommonAncestor[P crypto.PublicKey](db ethdb.Reader, a, b *types.Header[P]) *types.Header[P] {
	for bn := b.Number.Uint64(); a.Number.Uint64() > bn; {
		a = ReadHeader[P](db, a.ParentHash, a.Number.Uint64()-1)
		if a == nil {
			return nil
		}
	}
	for an := a.Number.Uint64(); an < b.Number.Uint64(); {
		b = ReadHeader[P](db, b.ParentHash, b.Number.Uint64()-1)
		if b == nil {
			return nil
		}
	}
	for a.Hash() != b.Hash() {
		a = ReadHeader[P](db, a.ParentHash, a.Number.Uint64()-1)
		if a == nil {
			return nil
		}
		b = ReadHeader[P](db, b.ParentHash, b.Number.Uint64()-1)
		if b == nil {
			return nil
		}
	}
	return a
}

// ReadHeadHeader returns the current canonical head header.
func ReadHeadHeader[P crypto.PublicKey](db ethdb.Reader) *types.Header[P] {
	headHeaderHash := ReadHeadHeaderHash(db)
	if headHeaderHash == (common.Hash{}) {
		return nil
	}
	headHeaderNumber := ReadHeaderNumber(db, headHeaderHash)
	if headHeaderNumber == nil {
		return nil
	}
	return ReadHeader[P](db, headHeaderHash, *headHeaderNumber)
}

// ReadHeadBlock returns the current canonical head block.
func ReadHeadBlock[P crypto.PublicKey](db ethdb.Reader) *types.Block[P] {
	headBlockHash := ReadHeadBlockHash(db)
	if headBlockHash == (common.Hash{}) {
		return nil
	}
	headBlockNumber := ReadHeaderNumber(db, headBlockHash)
	if headBlockNumber == nil {
		return nil
	}
	return ReadBlock[P](db, headBlockHash, *headBlockNumber)
}
