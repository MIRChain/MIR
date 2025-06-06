// Copyright 2019 The go-ethereum Authors
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
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/metrics"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/prometheus/tsdb/fileutil"
)

var (
	// errReadOnly is returned if the freezer is opened in read only mode. All the
	// mutations are disallowed.
	errReadOnly = errors.New("read only")

	// errUnknownTable is returned if the user attempts to read from a table that is
	// not tracked by the freezer.
	errUnknownTable = errors.New("unknown table")

	// errOutOrderInsertion is returned if the user attempts to inject out-of-order
	// binary blobs into the freezer.
	errOutOrderInsertion = errors.New("the append operation is out-order")

	// errSymlinkDatadir is returned if the ancient directory specified by user
	// is a symbolic link.
	errSymlinkDatadir = errors.New("symbolic link datadir is not supported")
)

const (
	// freezerRecheckInterval is the frequency to check the key-value database for
	// chain progression that might permit new blocks to be frozen into immutable
	// storage.
	freezerRecheckInterval = time.Minute

	// freezerBatchLimit is the maximum number of blocks to freeze in one batch
	// before doing an fsync and deleting it from the key-value store.
	freezerBatchLimit = 30000
)

// freezer is an memory mapped append-only database to store immutable chain data
// into flat files:
//
// - The append only nature ensures that disk writes are minimized.
// - The memory mapping ensures we can max out system memory for caching without
//   reserving it for go-ethereum. This would also reduce the memory requirements
//   of Geth, and thus also GC overhead.
type freezer [P crypto.PublicKey] struct {
	// WARNING: The `frozen` field is accessed atomically. On 32 bit platforms, only
	// 64-bit aligned fields can be atomic. The struct is guaranteed to be so aligned,
	// so take advantage of that (https://golang.org/pkg/sync/atomic/#pkg-note-BUG).
	frozen    uint64 // Number of blocks already frozen
	threshold uint64 // Number of recent blocks not to freeze (params.FullImmutabilityThreshold apart from tests)

	readonly     bool
	tables       map[string]*freezerTable // Data tables for storing everything
	instanceLock fileutil.Releaser        // File-system lock to prevent double opens

	trigger chan chan struct{} // Manual blocking freeze trigger, test determinism

	quit      chan struct{}
	closeOnce sync.Once
}

// newFreezer creates a chain freezer that moves ancient chain data into
// append-only flat file containers.
func newFreezer[P crypto.PublicKey](datadir string, namespace string, readonly bool) (*freezer[P], error) {
	// Create the initial freezer object
	var (
		readMeter  = metrics.NewRegisteredMeter(namespace+"ancient/read", nil)
		writeMeter = metrics.NewRegisteredMeter(namespace+"ancient/write", nil)
		sizeGauge  = metrics.NewRegisteredGauge(namespace+"ancient/size", nil)
	)
	// Ensure the datadir is not a symbolic link if it exists.
	if info, err := os.Lstat(datadir); !os.IsNotExist(err) {
		if info.Mode()&os.ModeSymlink != 0 {
			log.Warn("Symbolic link ancient database is not supported", "path", datadir)
			return nil, errSymlinkDatadir
		}
	}
	// Leveldb uses LOCK as the filelock filename. To prevent the
	// name collision, we use FLOCK as the lock name.
	lock, _, err := fileutil.Flock(filepath.Join(datadir, "FLOCK"))
	if err != nil {
		return nil, err
	}
	// Open all the supported data tables
	freezer := &freezer[P]{
		readonly:     readonly,
		threshold:    params.FullImmutabilityThreshold,
		tables:       make(map[string]*freezerTable),
		instanceLock: lock,
		trigger:      make(chan chan struct{}),
		quit:         make(chan struct{}),
	}
	for name, disableSnappy := range FreezerNoSnappy {
		table, err := newTable(datadir, name, readMeter, writeMeter, sizeGauge, disableSnappy)
		if err != nil {
			for _, table := range freezer.tables {
				table.Close()
			}
			lock.Release()
			return nil, err
		}
		freezer.tables[name] = table
	}
	if err := freezer.repair(); err != nil {
		for _, table := range freezer.tables {
			table.Close()
		}
		lock.Release()
		return nil, err
	}
	log.Info("Opened ancient database", "database", datadir, "readonly", readonly)
	return freezer, nil
}

// Close terminates the chain freezer, unmapping all the data files.
func (f *freezer[P]) Close() error {
	var errs []error
	f.closeOnce.Do(func() {
		close(f.quit)
		for _, table := range f.tables {
			if err := table.Close(); err != nil {
				errs = append(errs, err)
			}
		}
		if err := f.instanceLock.Release(); err != nil {
			errs = append(errs, err)
		}
	})
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}
	return nil
}

// HasAncient returns an indicator whether the specified ancient data exists
// in the freezer.
func (f *freezer[P]) HasAncient(kind string, number uint64) (bool, error) {
	if table := f.tables[kind]; table != nil {
		return table.has(number), nil
	}
	return false, nil
}

// Ancient retrieves an ancient binary blob from the append-only immutable files.
func (f *freezer[P]) Ancient(kind string, number uint64) ([]byte, error) {
	if table := f.tables[kind]; table != nil {
		return table.Retrieve(number)
	}
	return nil, errUnknownTable
}

// Ancients returns the length of the frozen items.
func (f *freezer[P]) Ancients() (uint64, error) {
	return atomic.LoadUint64(&f.frozen), nil
}

// AncientSize returns the ancient size of the specified category.
func (f *freezer[P]) AncientSize(kind string) (uint64, error) {
	if table := f.tables[kind]; table != nil {
		return table.size()
	}
	return 0, errUnknownTable
}

// AppendAncient injects all binary blobs belong to block at the end of the
// append-only immutable table files.
//
// Notably, this function is lock free but kind of thread-safe. All out-of-order
// injection will be rejected. But if two injections with same number happen at
// the same time, we can get into the trouble.
func (f *freezer[P]) AppendAncient(number uint64, hash, header, body, receipts, td []byte) (err error) {
	if f.readonly {
		return errReadOnly
	}
	// Ensure the binary blobs we are appending is continuous with freezer.
	if atomic.LoadUint64(&f.frozen) != number {
		return errOutOrderInsertion
	}
	// Rollback all inserted data if any insertion below failed to ensure
	// the tables won't out of sync.
	defer func() {
		if err != nil {
			rerr := f.repair()
			if rerr != nil {
				log.Crit("Failed to repair freezer", "err", rerr)
			}
			log.Info("Append ancient failed", "number", number, "err", err)
		}
	}()
	// Inject all the components into the relevant data tables
	if err := f.tables[freezerHashTable].Append(f.frozen, hash[:]); err != nil {
		log.Error("Failed to append ancient hash", "number", f.frozen, "hash", hash, "err", err)
		return err
	}
	if err := f.tables[freezerHeaderTable].Append(f.frozen, header); err != nil {
		log.Error("Failed to append ancient header", "number", f.frozen, "hash", hash, "err", err)
		return err
	}
	if err := f.tables[freezerBodiesTable].Append(f.frozen, body); err != nil {
		log.Error("Failed to append ancient body", "number", f.frozen, "hash", hash, "err", err)
		return err
	}
	if err := f.tables[freezerReceiptTable].Append(f.frozen, receipts); err != nil {
		log.Error("Failed to append ancient receipts", "number", f.frozen, "hash", hash, "err", err)
		return err
	}
	if err := f.tables[freezerDifficultyTable].Append(f.frozen, td); err != nil {
		log.Error("Failed to append ancient difficulty", "number", f.frozen, "hash", hash, "err", err)
		return err
	}
	atomic.AddUint64(&f.frozen, 1) // Only modify atomically
	return nil
}

// TruncateAncients discards any recent data above the provided threshold number.
func (f *freezer[P]) TruncateAncients(items uint64) error {
	if f.readonly {
		return errReadOnly
	}
	if atomic.LoadUint64(&f.frozen) <= items {
		return nil
	}
	for _, table := range f.tables {
		if err := table.truncate(items); err != nil {
			return err
		}
	}
	atomic.StoreUint64(&f.frozen, items)
	return nil
}

// Sync flushes all data tables to disk.
func (f *freezer[P]) Sync() error {
	return f.SyncRetry(1, 1*time.Second)
}

// SyncRetry
// Quorum
// add retry to sync
func (f *freezer[P]) SyncRetry(retry uint8, delay time.Duration) error {
	var errs []error
	for _, table := range f.tables {
		if err := table.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	hasError := len(errs) > 0
	if hasError && retry < 5 {
		log.Info("sync", "retry", retry, "errors", errs)
		time.Sleep(delay)
		return f.SyncRetry(retry+1, delay*2)
	} else if hasError {
		return fmt.Errorf("%v", errs)
	}
	return nil
}

// End Quorum

// freeze is a background thread that periodically checks the blockchain for any
// import progress and moves ancient data from the fast database into the freezer.
//
// This functionality is deliberately broken off from block importing to avoid
// incurring additional data shuffling delays on block propagation.
func (f *freezer[P]) freeze(db ethdb.KeyValueStore) {
	nfdb := &nofreezedb{KeyValueStore: db}

	var (
		backoff   bool
		triggered chan struct{} // Used in tests
	)
	for {
		select {
		case <-f.quit:
			log.Info("Freezer shutting down")
			return
		default:
		}
		if backoff {
			// If we were doing a manual trigger, notify it
			if triggered != nil {
				triggered <- struct{}{}
				triggered = nil
			}
			select {
			case <-time.NewTimer(freezerRecheckInterval).C:
				backoff = false
			case triggered = <-f.trigger:
				backoff = false
			case <-f.quit:
				return
			}
		}
		// Retrieve the freezing threshold.
		hash := ReadHeadBlockHash(nfdb)
		if hash == (common.Hash{}) {
			log.Debug("Current full block hash unavailable") // new chain, empty database
			backoff = true
			continue
		}
		number := ReadHeaderNumber(nfdb, hash)
		threshold := int(atomic.LoadUint64(&f.threshold))

		switch {
		case number == nil:
			log.Error("Current full block number unavailable", "hash", hash)
			backoff = true
			continue

		case *number < uint64(params.GetImmutabilityThresholdWithDefault(threshold)):
			log.Debug("Current full block not old enough", "number", *number, "hash", hash, "delay", params.GetImmutabilityThresholdWithDefault(threshold))
			backoff = true
			continue

		case *number-uint64(params.GetImmutabilityThresholdWithDefault(threshold)) <= f.frozen:
			log.Debug("Ancient blocks frozen already", "number", *number, "hash", hash, "frozen", f.frozen)
			backoff = true
			continue
		}
		head := ReadHeader[P](nfdb, hash, *number)
		if head == nil {
			log.Error("Current full block unavailable", "number", *number, "hash", hash)
			backoff = true
			continue
		}
		// Seems we have data ready to be frozen, process in usable batches
		limit := *number - uint64(params.GetImmutabilityThresholdWithDefault(threshold))
		if limit-f.frozen > freezerBatchLimit {
			limit = f.frozen + freezerBatchLimit
		}
		var (
			start    = time.Now()
			first    = f.frozen
			ancients = make([]common.Hash, 0, limit-f.frozen)
		)
		for f.frozen <= limit {
			// Retrieves all the components of the canonical block
			hash := ReadCanonicalHash(nfdb, f.frozen)
			if hash == (common.Hash{}) {
				log.Error("Canonical hash missing, can't freeze", "number", f.frozen)
				break
			}
			header := ReadHeaderRLP[P](nfdb, hash, f.frozen)
			if len(header) == 0 {
				log.Error("Block header missing, can't freeze", "number", f.frozen, "hash", hash)
				break
			}
			body := ReadBodyRLP(nfdb, hash, f.frozen)
			if len(body) == 0 {
				log.Error("Block body missing, can't freeze", "number", f.frozen, "hash", hash)
				break
			}
			receipts := ReadReceiptsRLP(nfdb, hash, f.frozen)
			if len(receipts) == 0 {
				log.Error("Block receipts missing, can't freeze", "number", f.frozen, "hash", hash)
				break
			}
			td := ReadTdRLP(nfdb, hash, f.frozen)
			if len(td) == 0 {
				log.Error("Total difficulty missing, can't freeze", "number", f.frozen, "hash", hash)
				break
			}
			log.Trace("Deep froze ancient block", "number", f.frozen, "hash", hash)
			// Inject all the components into the relevant data tables
			if err := f.AppendAncient(f.frozen, hash[:], header, body, receipts, td); err != nil {
				break
			}
			ancients = append(ancients, hash)
		}
		// Batch of blocks have been frozen, flush them before wiping from leveldb
		if err := f.Sync(); err != nil {
			log.Crit("Failed to flush frozen tables", "err", err)
		}
		// Wipe out all data from the active database
		batch := db.NewBatch()
		for i := 0; i < len(ancients); i++ {
			// Always keep the genesis block in active database
			if first+uint64(i) != 0 {
				DeleteBlockWithoutNumber(batch, ancients[i], first+uint64(i))
				DeleteCanonicalHash(batch, first+uint64(i))
			}
		}
		if err := batch.Write(); err != nil {
			log.Crit("Failed to delete frozen canonical blocks", "err", err)
		}
		batch.Reset()

		// Wipe out side chains also and track dangling side chians
		var dangling []common.Hash
		for number := first; number < f.frozen; number++ {
			// Always keep the genesis block in active database
			if number != 0 {
				dangling = ReadAllHashes(db, number)
				for _, hash := range dangling {
					log.Trace("Deleting side chain", "number", number, "hash", hash)
					DeleteBlock(batch, hash, number)
				}
			}
		}
		if err := batch.Write(); err != nil {
			log.Crit("Failed to delete frozen side blocks", "err", err)
		}
		batch.Reset()

		// Step into the future and delete and dangling side chains
		if f.frozen > 0 {
			tip := f.frozen
			for len(dangling) > 0 {
				drop := make(map[common.Hash]struct{})
				for _, hash := range dangling {
					log.Debug("Dangling parent from freezer", "number", tip-1, "hash", hash)
					drop[hash] = struct{}{}
				}
				children := ReadAllHashes(db, tip)
				for i := 0; i < len(children); i++ {
					// Dig up the child and ensure it's dangling
					child := ReadHeader[P](nfdb, children[i], tip)
					if child == nil {
						log.Error("Missing dangling header", "number", tip, "hash", children[i])
						continue
					}
					if _, ok := drop[child.ParentHash]; !ok {
						children = append(children[:i], children[i+1:]...)
						i--
						continue
					}
					// Delete all block data associated with the child
					log.Debug("Deleting dangling block", "number", tip, "hash", children[i], "parent", child.ParentHash)
					DeleteBlock(batch, children[i], tip)
				}
				dangling = children
				tip++
			}
			if err := batch.Write(); err != nil {
				log.Crit("Failed to delete dangling side blocks", "err", err)
			}
		}
		// Log something friendly for the user
		context := []interface{}{
			"blocks", f.frozen - first, "elapsed", common.PrettyDuration(time.Since(start)), "number", f.frozen - 1,
		}
		if n := len(ancients); n > 0 {
			context = append(context, []interface{}{"hash", ancients[n-1]}...)
		}
		log.Info("Deep froze chain segment", context...)

		// Avoid database thrashing with tiny writes
		if f.frozen-first < freezerBatchLimit {
			backoff = true
		}
	}
}

// repair truncates all data tables to the same length.
func (f *freezer[P]) repair() error {
	min := uint64(math.MaxUint64)
	for _, table := range f.tables {
		items := atomic.LoadUint64(&table.items)
		if min > items {
			min = items
		}
	}
	for _, table := range f.tables {
		if err := table.truncate(min); err != nil {
			return err
		}
	}
	atomic.StoreUint64(&f.frozen, min)
	return nil
}
