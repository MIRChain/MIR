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

package miner

import (
	"bytes"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/misc"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

const (
	// resultQueueSize is the size of channel listening to sealing result.
	resultQueueSize = 10

	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096

	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10

	// chainSideChanSize is the size of channel listening to ChainSideEvent.
	chainSideChanSize = 10

	// resubmitAdjustChanSize is the size of resubmitting interval adjustment channel.
	resubmitAdjustChanSize = 10

	// miningLogAtDepth is the number of confirmations before logging successful mining.
	miningLogAtDepth = 7

	// minRecommitInterval is the minimal time interval to recreate the mining block with
	// any newly arrived transactions.
	minRecommitInterval = 1 * time.Second

	// maxRecommitInterval is the maximum time interval to recreate the mining block with
	// any newly arrived transactions.
	maxRecommitInterval = 15 * time.Second

	// intervalAdjustRatio is the impact a single interval adjustment has on sealing work
	// resubmitting interval.
	intervalAdjustRatio = 0.1

	// intervalAdjustBias is applied during the new resubmit interval calculation in favor of
	// increasing upper limit or decreasing lower limit so that the limit can be reachable.
	intervalAdjustBias = 200 * 1000.0 * 1000.0

	// staleThreshold is the maximum depth of the acceptable stale block.
	staleThreshold = 7
)

// environment is the worker's current environment and holds all of the current state information.
type environment [P crypto.PublicKey] struct {
	signer types.Signer[P]

	state     *state.StateDB[P] // apply state changes here
	ancestors mapset.Set     // ancestor set (used for checking uncle parent validity)
	family    mapset.Set     // family set (used for checking uncle invalidity)
	uncles    mapset.Set     // uncle set
	tcount    int            // tx count in cycle
	gasPool   *core.GasPool  // available gas used to pack transactions

	header   *types.Header[P]
	txs      []*types.Transaction[P]
	receipts []*types.Receipt[P]

	// Quorum
	privateReceipts  []*types.Receipt[P]
	privateStateRepo mps.PrivateStateRepository[P]
	// End Quorum
}

// task contains all information for consensus engine sealing and result submitting.
type task [P crypto.PublicKey] struct {
	receipts  []*types.Receipt[P]
	state     *state.StateDB[P]
	block     *types.Block[P]
	createdAt time.Time

	// Quorum
	privateReceipts  []*types.Receipt[P]
	privateStateRepo mps.PrivateStateRepository[P]
	// End Quorum
}

const (
	commitInterruptNone int32 = iota
	commitInterruptNewHead
	commitInterruptResubmit
)

// newWorkReq represents a request for new sealing work submitting with relative interrupt notifier.
type newWorkReq struct {
	interrupt *int32
	noempty   bool
	timestamp int64
}

// intervalAdjust represents a resubmitting interval adjustment.
type intervalAdjust struct {
	ratio float64
	inc   bool
}

// worker is the main object which takes care of submitting new work to consensus engine
// and gathering the sealing result.
type worker [T crypto.PrivateKey, P crypto.PublicKey] struct {
	config      *Config
	chainConfig *params.ChainConfig
	engine      consensus.Engine[P]
	eth         Backend[T,P]
	chain       *core.BlockChain[P]

	// Feeds
	pendingLogsFeed event.Feed

	// Subscriptions
	mux          *event.TypeMux
	txsCh        chan core.NewTxsEvent[P]
	txsSub       event.Subscription
	chainHeadCh  chan core.ChainHeadEvent[P]
	chainHeadSub event.Subscription
	chainSideCh  chan core.ChainSideEvent[P]
	chainSideSub event.Subscription

	// Channels
	newWorkCh          chan *newWorkReq
	taskCh             chan *task[P]
	resultCh           chan *types.Block[P]
	startCh            chan struct{}
	exitCh             chan struct{}
	resubmitIntervalCh chan time.Duration
	resubmitAdjustCh   chan *intervalAdjust

	current      *environment[P]                 // An environment for current running cycle.
	localUncles  map[common.Hash]*types.Block[P] // A set of side blocks generated locally as the possible uncle blocks.
	remoteUncles map[common.Hash]*types.Block[P] // A set of side blocks as the possible uncle blocks.
	unconfirmed  *unconfirmedBlocks[P]           // A set of locally mined blocks pending canonicalness confirmations.

	mu       sync.RWMutex // The lock used to protect the coinbase and extra fields
	coinbase common.Address
	extra    []byte

	pendingMu    sync.RWMutex
	pendingTasks map[common.Hash]*task[P]

	snapshotMu    sync.RWMutex // The lock used to protect the block snapshot and state snapshot
	snapshotBlock *types.Block[P]
	snapshotState *state.StateDB[P]

	// atomic status counters
	running int32 // The indicator whether the consensus engine is running or not.
	newTxs  int32 // New arrival transaction count since last sealing work submitting.

	// noempty is the flag used to control whether the feature of pre-seal empty
	// block is enabled. The default value is false(pre-seal is enabled by default).
	// But in some special scenario the consensus engine will seal blocks instantaneously,
	// in this case this feature will add all empty blocks into canonical chain
	// non-stop and no real transaction will be included.
	noempty uint32

	// External functions
	isLocalBlock func(block *types.Block[P]) bool // Function used to determine whether the specified block is mined by local miner.

	// Test hooks
	newTaskHook  func(*task[P])                        // Method to call upon receiving a new sealing task.
	skipSealHook func(*task[P]) bool                   // Method to decide whether skipping the sealing.
	fullTaskHook func()                             // Method to call before pushing the full sealing task.
	resubmitHook func(time.Duration, time.Duration) // Method to call upon updating resubmitting interval.
}

func newWorker[T crypto.PrivateKey, P crypto.PublicKey](config *Config, chainConfig *params.ChainConfig, engine consensus.Engine[P], eth Backend[T,P], mux *event.TypeMux, isLocalBlock func(*types.Block[P]) bool, init bool) *worker[T,P] {
	worker := &worker[T,P]{
		config:             config,
		chainConfig:        chainConfig,
		engine:             engine,
		eth:                eth,
		mux:                mux,
		chain:              eth.BlockChain(),
		isLocalBlock:       isLocalBlock,
		localUncles:        make(map[common.Hash]*types.Block[P]),
		remoteUncles:       make(map[common.Hash]*types.Block[P]),
		unconfirmed:        newUnconfirmedBlocks[P](eth.BlockChain(), miningLogAtDepth),
		pendingTasks:       make(map[common.Hash]*task[P]),
		txsCh:              make(chan core.NewTxsEvent[P], txChanSize),
		chainHeadCh:        make(chan core.ChainHeadEvent[P], chainHeadChanSize),
		chainSideCh:        make(chan core.ChainSideEvent[P], chainSideChanSize),
		newWorkCh:          make(chan *newWorkReq),
		taskCh:             make(chan *task[P]),
		resultCh:           make(chan *types.Block[P], resultQueueSize),
		exitCh:             make(chan struct{}),
		startCh:            make(chan struct{}, 1),
		resubmitIntervalCh: make(chan time.Duration),
		resubmitAdjustCh:   make(chan *intervalAdjust, resubmitAdjustChanSize),
	}
	if _, ok := engine.(consensus.Istanbul[P]); ok || !chainConfig.IsQuorum || chainConfig.Clique != nil || chainConfig.Ethash != nil {
		// Subscribe NewTxsEvent for tx pool
		worker.txsSub = eth.TxPool().SubscribeNewTxsEvent(worker.txsCh)
		// Subscribe events for blockchain
		worker.chainHeadSub = eth.BlockChain().SubscribeChainHeadEvent(worker.chainHeadCh)
		worker.chainSideSub = eth.BlockChain().SubscribeChainSideEvent(worker.chainSideCh)
		// Sanitize recommit interval if the user-specified one is too short.
		recommit := worker.config.Recommit
		if recommit < minRecommitInterval {
			log.Warn("Sanitizing miner recommit interval", "provided", recommit, "updated", minRecommitInterval)
			recommit = minRecommitInterval
		}

		go worker.mainLoop()
		go worker.newWorkLoop(recommit)
		go worker.resultLoop()
		go worker.taskLoop()

		// Submit first work to initialize pending state.
		if init {
			worker.startCh <- struct{}{}
		}
	}
	return worker
}

// setEtherbase sets the etherbase used to initialize the block coinbase field.
func (w *worker[T,P]) setEtherbase(addr common.Address) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.coinbase = addr
}

// setExtra sets the content used to initialize the block extra field.
func (w *worker[T,P]) setExtra(extra []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.extra = extra
}

// setRecommitInterval updates the interval for miner sealing work recommitting.
func (w *worker[T,P]) setRecommitInterval(interval time.Duration) {
	w.resubmitIntervalCh <- interval
}

// disablePreseal disables pre-sealing mining feature
func (w *worker[T,P]) disablePreseal() {
	atomic.StoreUint32(&w.noempty, 1)
}

// enablePreseal enables pre-sealing mining feature
func (w *worker[T,P]) enablePreseal() {
	atomic.StoreUint32(&w.noempty, 0)
}

// pending returns the pending state and corresponding block.
func (w *worker[T,P]) pending(psi types.PrivateStateIdentifier) (*types.Block[P], *state.StateDB[P], *state.StateDB[P]) {
	// return a snapshot to avoid contention on currentMu mutex
	w.snapshotMu.RLock()
	defer w.snapshotMu.RUnlock()
	if w.snapshotState == nil {
		return nil, nil, nil
	}
	privateState, err := w.current.privateStateRepo.StatePSI(psi)
	if err != nil {
		log.Error("Unable to retrieve private state", "psi", psi, "err", err)
		return nil, nil, nil
	}
	return w.snapshotBlock, w.snapshotState.Copy(), privateState.Copy()
}

// pendingBlock returns pending block.
func (w *worker[T,P]) pendingBlock() *types.Block[P] {
	// return a snapshot to avoid contention on currentMu mutex
	w.snapshotMu.RLock()
	defer w.snapshotMu.RUnlock()
	return w.snapshotBlock
}

// start sets the running status as 1 and triggers new work submitting.
func (w *worker[T,P]) start() {
	atomic.StoreInt32(&w.running, 1)
	if istanbul, ok := w.engine.(consensus.Istanbul[P]); ok {
		istanbul.Start(w.chain, w.chain.CurrentBlock, rawdb.HasBadBlock[P])
	} else {
		log.Warn("no start of engine", "engine", w.engine)
	}
	w.startCh <- struct{}{}
}

// stop sets the running status as 0.
func (w *worker[T,P]) stop() {
	if istanbul, ok := w.engine.(consensus.Istanbul[P]); ok {
		istanbul.Stop()
	}
	atomic.StoreInt32(&w.running, 0)
}

// isRunning returns an indicator whether worker is running or not.
func (w *worker[T,P]) isRunning() bool {
	return atomic.LoadInt32(&w.running) == 1
}

// close terminates all background threads maintained by the worker.
// Note the worker does not support being closed multiple times.
func (w *worker[T,P]) close() {
	if w.current != nil && w.current.state != nil {
		w.current.state.StopPrefetcher()
	}
	atomic.StoreInt32(&w.running, 0)
	close(w.exitCh)
}

// recalcRecommit recalculates the resubmitting interval upon feedback.
func recalcRecommit(minRecommit, prev time.Duration, target float64, inc bool) time.Duration {
	var (
		prevF = float64(prev.Nanoseconds())
		next  float64
	)
	if inc {
		next = prevF*(1-intervalAdjustRatio) + intervalAdjustRatio*(target+intervalAdjustBias)
		max := float64(maxRecommitInterval.Nanoseconds())
		if next > max {
			next = max
		}
	} else {
		next = prevF*(1-intervalAdjustRatio) + intervalAdjustRatio*(target-intervalAdjustBias)
		min := float64(minRecommit.Nanoseconds())
		if next < min {
			next = min
		}
	}
	return time.Duration(int64(next))
}

// newWorkLoop is a standalone goroutine to submit new mining work upon received events.
func (w *worker[T,P]) newWorkLoop(recommit time.Duration) {
	var (
		interrupt   *int32
		minRecommit = recommit // minimal resubmit interval specified by user.
		timestamp   int64      // timestamp for each round of mining.
	)

	timer := time.NewTimer(0)
	defer timer.Stop()
	<-timer.C // discard the initial tick

	// commit aborts in-flight transaction execution with given signal and resubmits a new one.
	commit := func(noempty bool, s int32) {
		if interrupt != nil {
			atomic.StoreInt32(interrupt, s)
		}
		interrupt = new(int32)
		select {
		case w.newWorkCh <- &newWorkReq{interrupt: interrupt, noempty: noempty, timestamp: timestamp}:
		case <-w.exitCh:
			return
		}
		timer.Reset(recommit)
		atomic.StoreInt32(&w.newTxs, 0)
	}
	// clearPending cleans the stale pending tasks.
	clearPending := func(number uint64) {
		w.pendingMu.Lock()
		for h, t := range w.pendingTasks {
			if t.block.NumberU64()+staleThreshold <= number {
				delete(w.pendingTasks, h)
			}
		}
		w.pendingMu.Unlock()
	}

	for {
		select {
		case <-w.startCh:
			clearPending(w.chain.CurrentBlock().NumberU64())
			timestamp = time.Now().Unix()
			commit(false, commitInterruptNewHead)

		case head := <-w.chainHeadCh:
			if h, ok := w.engine.(consensus.Handler[P]); ok {
				h.NewChainHead()
			}
			clearPending(head.Block.NumberU64())
			timestamp = time.Now().Unix()
			commit(false, commitInterruptNewHead)

		case <-timer.C:
			// If mining is running resubmit a new work cycle periodically to pull in
			// higher priced transactions. Disable this overhead for pending blocks.
			if w.isRunning() && (w.chainConfig.Clique == nil || w.chainConfig.Clique.Period > 0) {
				// Short circuit if no new transaction arrives.
				if atomic.LoadInt32(&w.newTxs) == 0 {
					timer.Reset(recommit)
					continue
				}
				commit(true, commitInterruptResubmit)
			}

		case interval := <-w.resubmitIntervalCh:
			// Adjust resubmit interval explicitly by user.
			if interval < minRecommitInterval {
				log.Warn("Sanitizing miner recommit interval", "provided", interval, "updated", minRecommitInterval)
				interval = minRecommitInterval
			}
			log.Info("Miner recommit interval update", "from", minRecommit, "to", interval)
			minRecommit, recommit = interval, interval

			if w.resubmitHook != nil {
				w.resubmitHook(minRecommit, recommit)
			}

		case adjust := <-w.resubmitAdjustCh:
			// Adjust resubmit interval by feedback.
			if adjust.inc {
				before := recommit
				target := float64(recommit.Nanoseconds()) / adjust.ratio
				recommit = recalcRecommit(minRecommit, recommit, target, true)
				log.Trace("Increase miner recommit interval", "from", before, "to", recommit)
			} else {
				before := recommit
				recommit = recalcRecommit(minRecommit, recommit, float64(minRecommit.Nanoseconds()), false)
				log.Trace("Decrease miner recommit interval", "from", before, "to", recommit)
			}

			if w.resubmitHook != nil {
				w.resubmitHook(minRecommit, recommit)
			}

		case <-w.exitCh:
			return
		}
	}
}

// mainLoop is a standalone goroutine to regenerate the sealing task based on the received event.
func (w *worker[T,P]) mainLoop() {
	defer w.txsSub.Unsubscribe()
	defer w.chainHeadSub.Unsubscribe()
	defer w.chainSideSub.Unsubscribe()

	for {
		select {
		case req := <-w.newWorkCh:
			w.commitNewWork(req.interrupt, req.noempty, req.timestamp)

		case ev := <-w.chainSideCh:
			// Short circuit for duplicate side blocks
			if _, exist := w.localUncles[ev.Block.Hash()]; exist {
				continue
			}
			if _, exist := w.remoteUncles[ev.Block.Hash()]; exist {
				continue
			}
			// Add side block to possible uncle block set depending on the author.
			if w.isLocalBlock != nil && w.isLocalBlock(ev.Block) {
				w.localUncles[ev.Block.Hash()] = ev.Block
			} else {
				w.remoteUncles[ev.Block.Hash()] = ev.Block
			}
			// If our mining block contains less than 2 uncle blocks,
			// add the new uncle block if valid and regenerate a mining block.
			if w.isRunning() && w.current != nil && w.current.uncles.Cardinality() < 2 {
				start := time.Now()
				if err := w.commitUncle(w.current, ev.Block.Header()); err == nil {
					var uncles []*types.Header[P]
					w.current.uncles.Each(func(item interface{}) bool {
						hash, ok := item.(common.Hash)
						if !ok {
							return false
						}
						uncle, exist := w.localUncles[hash]
						if !exist {
							uncle, exist = w.remoteUncles[hash]
						}
						if !exist {
							return false
						}
						uncles = append(uncles, uncle.Header())
						return false
					})
					w.commit(uncles, nil, true, start)
				}
			}

		case ev := <-w.txsCh:
			// Apply transactions to the pending state if we're not mining.
			//
			// Note all transactions received may not be continuous with transactions
			// already included in the current mining block. These transactions will
			// be automatically eliminated.
			if !w.isRunning() && w.current != nil {
				// If block is already full, abort
				if gp := w.current.gasPool; gp != nil && gp.Gas() < params.TxGas {
					continue
				}
				w.mu.RLock()
				coinbase := w.coinbase
				w.mu.RUnlock()

				txs := make(map[common.Address]types.Transactions[P])
				for _, tx := range ev.Txs {
					acc, _ := types.Sender(w.current.signer, tx)
					txs[acc] = append(txs[acc], tx)
				}
				txset := types.NewTransactionsByPriceAndNonce(w.current.signer, txs)
				tcount := w.current.tcount
				w.commitTransactions(txset, coinbase, nil)
				// Only update the snapshot if any new transactons were added
				// to the pending block
				if tcount != w.current.tcount {
					w.updateSnapshot()
				}
			} else {
				// Special case, if the consensus engine is 0 period clique(dev mode),
				// submit mining work here since all empty submission will be rejected
				// by clique. Of course the advance sealing(empty submission) is disabled.
				if w.chainConfig.Clique != nil && w.chainConfig.Clique.Period == 0 {
					w.commitNewWork(nil, true, time.Now().Unix())
				}
			}
			atomic.AddInt32(&w.newTxs, int32(len(ev.Txs)))

		// System stopped
		case <-w.exitCh:
			return
		case <-w.txsSub.Err():
			return
		case <-w.chainHeadSub.Err():
			return
		case <-w.chainSideSub.Err():
			return
		}
	}
}

// taskLoop is a standalone goroutine to fetch sealing task from the generator and
// push them to consensus engine.
func (w *worker[T,P]) taskLoop() {
	var (
		stopCh chan struct{}
		prev   common.Hash
	)

	// interrupt aborts the in-flight sealing task.
	interrupt := func() {
		if stopCh != nil {
			close(stopCh)
			stopCh = nil
		}
	}
	for {
		select {
		case task := <-w.taskCh:
			if w.newTaskHook != nil {
				w.newTaskHook(task)
			}
			// Reject duplicate sealing work due to resubmitting.
			sealHash := w.engine.SealHash(task.block.Header())
			if sealHash == prev {
				continue
			}
			// Interrupt previous sealing operation
			interrupt()
			stopCh, prev = make(chan struct{}), sealHash

			if w.skipSealHook != nil && w.skipSealHook(task) {
				continue
			}
			w.pendingMu.Lock()
			w.pendingTasks[sealHash] = task
			w.pendingMu.Unlock()

			if err := w.engine.Seal(w.chain, task.block, w.resultCh, stopCh); err != nil {
				log.Warn("Block sealing failed", "err", err)
			}
		case <-w.exitCh:
			interrupt()
			return
		}
	}
}

// resultLoop is a standalone goroutine to handle sealing result submitting
// and flush relative data to the database.
func (w *worker[T,P]) resultLoop() {
	for {
		select {
		case block := <-w.resultCh:
			// Short circuit when receiving empty result.
			if block == nil {
				continue
			}
			// Short circuit when receiving duplicate result caused by resubmitting.
			if w.chain.HasBlock(block.Hash(), block.NumberU64()) {
				continue
			}
			var (
				sealhash = w.engine.SealHash(block.Header())
				hash     = block.Hash()
			)
			w.pendingMu.RLock()
			task, exist := w.pendingTasks[sealhash]
			w.pendingMu.RUnlock()
			if !exist {
				log.Error("Block found but no relative pending task", "number", block.Number(), "sealhash", sealhash, "hash", hash)
				continue
			}
			// Different block could share same sealhash, deep copy here to prevent write-write conflict.
			var (
				pubReceipts = make([]*types.Receipt[P], len(task.receipts))
				prvReceipts = make([]*types.Receipt[P], len(task.privateReceipts))
				logs        []*types.Log
			)
			offset := len(task.receipts)
			for i, receipt := range task.receipts {
				// add block location fields
				receipt.BlockHash = hash
				receipt.BlockNumber = block.Number()
				receipt.TransactionIndex = uint(i)

				pubReceipts[i] = new(types.Receipt[P])
				*pubReceipts[i] = *receipt
				// Update the block hash in all logs since it is now available and not when the
				// receipt/log of individual transactions were created.
				for _, log := range receipt.Logs {
					log.BlockHash = hash
				}
				logs = append(logs, receipt.Logs...)

				// If this was a public privacy marker transaction then there will be an associated private receipt to handle.
				if receipt.PSReceipts != nil {
					tx := block.Transaction(receipt.TxHash)
					if tx.IsPrivacyMarker() {
						for _, markerReceipt := range receipt.PSReceipts {
							markerReceipt.BlockHash = hash
							markerReceipt.BlockNumber = block.Number()
							markerReceipt.TransactionIndex = uint(i)

							// Update the block hash in all logs since it is now available and not when the
							// receipt/log of individual transactions were created.
							for _, log := range markerReceipt.Logs {
								log.BlockHash = hash
							}
							// Note that we don't append logs here, else will get duplicates.
						}
					}
				}
			}

			for i, receipt := range task.privateReceipts {
				// add block location fields
				receipt.BlockHash = hash
				receipt.BlockNumber = block.Number()
				receipt.TransactionIndex = uint(i + offset)

				prvReceipts[i] = new(types.Receipt[P])
				*prvReceipts[i] = *receipt
				// Update the block hash in all logs since it is now available and not when the
				// receipt/log of individual transactions were created.
				for _, log := range receipt.Logs {
					log.BlockHash = hash
				}
				logs = append(logs, receipt.Logs...)

				for _, psReceipt := range receipt.PSReceipts {
					// if block location fields are already populated then this is a privacy marker receipt and it is
					// already handled - skip processing it again
					if psReceipt.BlockNumber != nil {
						continue
					}
					// add block location fields
					psReceipt.BlockHash = hash
					psReceipt.BlockNumber = block.Number()
					psReceipt.TransactionIndex = uint(i + offset)

					// Update the block hash in all logs since it is now available and not when the
					// receipt/log of individual transactions were created.
					for _, log := range psReceipt.Logs {
						log.BlockHash = hash
					}
					logs = append(logs, psReceipt.Logs...)
				}
			}

			allReceipts := task.privateStateRepo.MergeReceipts(pubReceipts, prvReceipts)

			// Commit block and state to database.
			_, err := w.chain.WriteBlockWithState(block, allReceipts, logs, task.state, task.privateStateRepo, true)
			if err != nil {
				log.Error("Failed writing block to chain", "err", err)
				continue
			}
			if err := rawdb.WritePrivateBlockBloom(w.eth.ChainDb(), block.NumberU64(), prvReceipts); err != nil {
				log.Error("Failed writing private block bloom", "err", err)
				continue
			}
			log.Info("Successfully sealed new block", "number", block.Number(), "sealhash", sealhash, "hash", hash,
				"elapsed", common.PrettyDuration(time.Since(task.createdAt)))

			// Broadcast the block and announce chain insertion event
			w.mux.Post(core.NewMinedBlockEvent[P]{Block: block})

			// Insert the block into the set of pending ones to resultLoop for confirmations
			w.unconfirmed.Insert(block.NumberU64(), block.Hash())

		case <-w.exitCh:
			return
		}
	}
}

// makeCurrent creates a new environment for the current cycle.
func (w *worker[T,P]) makeCurrent(parent *types.Block[P], header *types.Header[P]) error {
	// Retrieve the parent state to execute on top and start a prefetcher for
	// the miner to speed block sealing up a bit
	publicState, privateStateRepo, err := w.chain.StateAt(parent.Root())
	if err != nil {
		return err
	}
	publicState.StartPrefetcher("miner")

	env := &environment[P]{
		signer:    types.MakeSigner[P](w.chainConfig, header.Number),
		state:     publicState,
		ancestors: mapset.NewSet(),
		family:    mapset.NewSet(),
		uncles:    mapset.NewSet(),
		header:    header,
		// Quorum
		privateStateRepo: privateStateRepo,
	}
	// when 08 is processed ancestors contain 07 (quick block)
	for _, ancestor := range w.chain.GetBlocksFromHash(parent.Hash(), 7) {
		for _, uncle := range ancestor.Uncles() {
			env.family.Add(uncle.Hash())
		}
		env.family.Add(ancestor.Hash())
		env.ancestors.Add(ancestor.Hash())
	}
	// Keep track of transactions which return errors so they can be removed
	env.tcount = 0

	// Swap out the old work with the new one, terminating any leftover prefetcher
	// processes in the mean time and starting a new one.
	if w.current != nil && w.current.state != nil {
		w.current.state.StopPrefetcher()
	}
	w.current = env
	return nil
}

// commitUncle adds the given block to uncle block set, returns error if failed to add.
func (w *worker[T,P]) commitUncle(env *environment[P], uncle *types.Header[P]) error {
	hash := uncle.Hash()
	if env.uncles.Contains(hash) {
		return errors.New("uncle not unique")
	}
	if env.header.ParentHash == uncle.ParentHash {
		return errors.New("uncle is sibling")
	}
	if !env.ancestors.Contains(uncle.ParentHash) {
		return errors.New("uncle's parent unknown")
	}
	if env.family.Contains(hash) {
		return errors.New("uncle already included")
	}
	env.uncles.Add(uncle.Hash())
	return nil
}

// updateSnapshot updates pending snapshot block and state.
// Note this function assumes the current variable is thread safe.
func (w *worker[T,P]) updateSnapshot() {
	w.snapshotMu.Lock()
	defer w.snapshotMu.Unlock()

	var uncles []*types.Header[P]
	w.current.uncles.Each(func(item interface{}) bool {
		hash, ok := item.(common.Hash)
		if !ok {
			return false
		}
		uncle, exist := w.localUncles[hash]
		if !exist {
			uncle, exist = w.remoteUncles[hash]
		}
		if !exist {
			return false
		}
		uncles = append(uncles, uncle.Header())
		return false
	})

	w.snapshotBlock = types.NewBlock(
		w.current.header,
		w.current.txs,
		uncles,
		w.current.receipts,
		trie.NewStackTrie[P](nil),
	)
	w.snapshotState = w.current.state.Copy()
}

func (w *worker[T,P]) commitTransaction(tx *types.Transaction[P], coinbase common.Address) ([]*types.Log, error) {
	workerEnv := w.current
	publicStateDB := workerEnv.state
	snap := publicStateDB.Snapshot()
	privateStateRepo := workerEnv.privateStateRepo
	txnStart := time.Now()

	mpsReceipt, privateStateSnaphots, err := w.handleMPS(tx, coinbase, false)
	if err != nil {
		w.revertToPrivateStateSnapshots(privateStateSnaphots)
		return nil, err
	}
	privateStateDB, err := privateStateRepo.DefaultState()
	if err != nil {
		w.revertToPrivateStateSnapshots(privateStateSnaphots)
		return nil, err
	}
	privateStateDB.Prepare(tx.Hash(), common.Hash{}, workerEnv.tcount)
	publicStateDB.Prepare(tx.Hash(), common.Hash{}, workerEnv.tcount)
	privateStateSnaphots[privateStateRepo.DefaultStateMetadata().ID] = privateStateDB.Snapshot()
	receipt, privateReceipt, err := core.ApplyTransaction[P](w.chainConfig, w.chain, &coinbase, workerEnv.gasPool, publicStateDB, privateStateDB, workerEnv.header, tx, &workerEnv.header.GasUsed, *w.chain.GetVMConfig(), privateStateRepo.IsMPS(), privateStateRepo, false)
	if err != nil {
		publicStateDB.RevertToSnapshot(snap)
		w.revertToPrivateStateSnapshots(privateStateSnaphots)
		return nil, err
	}
	workerEnv.txs = append(workerEnv.txs, tx)
	workerEnv.receipts = append(workerEnv.receipts, receipt)
	log.EmitCheckpoint(log.TxCompleted, "tx", tx.Hash().Hex(), "time", time.Since(txnStart))

	logs := receipt.Logs

	// Quorum
	if privateReceipt != nil {
		newPrivateReceipt, privateLogs := core.HandlePrivateReceipt(receipt, privateReceipt, mpsReceipt, tx, privateStateDB, privateStateRepo, w.chain)
		workerEnv.privateReceipts = append(workerEnv.privateReceipts, newPrivateReceipt)
		logs = append(logs, privateLogs...)
	}
	// End Quorum

	return logs, nil
}

func (w *worker[T,P]) commitTransactions(txs *types.TransactionsByPriceAndNonce[P], coinbase common.Address, interrupt *int32) bool {
	// Short circuit if current is nil
	if w.current == nil {
		return true
	}

	if w.current.gasPool == nil {
		w.current.gasPool = new(core.GasPool).AddGas(w.current.header.GasLimit)
	}

	var coalescedLogs []*types.Log

	loopStartTime := time.Now() // Quorum
	for {
		// In the following three cases, we will interrupt the execution of the transaction.
		// (1) new head block event arrival, the interrupt signal is 1
		// (2) worker start or restart, the interrupt signal is 1
		// (3) worker recreate the mining block with any newly arrived transactions, the interrupt signal is 2.
		// For the first two cases, the semi-finished work will be discarded.
		// For the third case, the semi-finished work will be submitted to the consensus engine.
		if interrupt != nil && atomic.LoadInt32(interrupt) != commitInterruptNone {
			log.Info("Aborting transaction processing due to 'commitInterruptNewHead',", "elapsed time", time.Since(loopStartTime)) // Quorum
			// Notify resubmit loop to increase resubmitting interval due to too frequent commits.
			if atomic.LoadInt32(interrupt) == commitInterruptResubmit {
				ratio := float64(w.current.header.GasLimit-w.current.gasPool.Gas()) / float64(w.current.header.GasLimit)
				if ratio < 0.1 {
					ratio = 0.1
				}
				w.resubmitAdjustCh <- &intervalAdjust{
					ratio: ratio,
					inc:   true,
				}
			}
			return atomic.LoadInt32(interrupt) == commitInterruptNewHead
		}
		// If we don't have enough gas for any further transactions then we're done
		if w.current.gasPool.Gas() < params.TxGas {
			log.Trace("Not enough gas for further transactions", "have", w.current.gasPool, "want", params.TxGas)
			break
		}
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := types.Sender(w.current.signer, tx)
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !w.chainConfig.IsEIP155(w.current.header.Number) && !tx.IsPrivate() {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", w.chainConfig.EIP155Block)

			txs.Pop()
			continue
		}
		// Start executing the transaction
		logs, err := w.commitTransaction(tx, coinbase)
		switch {
		case errors.Is(err, core.ErrGasLimitReached):
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Trace("Gas limit exceeded for current block", "sender", from)
			txs.Pop()

		case errors.Is(err, core.ErrNonceTooLow):
			// New head notification data race between the transaction pool and miner, shift
			log.Trace("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce())
			txs.Shift()

		case errors.Is(err, core.ErrNonceTooHigh):
			// Reorg notification data race between the transaction pool and miner, skip account =
			log.Trace("Skipping account with hight nonce", "sender", from, "nonce", tx.Nonce())
			txs.Pop()

		case errors.Is(err, nil):
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			w.current.tcount++
			txs.Shift()

		case errors.Is(err, core.ErrTxTypeNotSupported):
			// Pop the unsupported transaction without shifting in the next from the account
			log.Trace("Skipping unsupported transaction type", "sender", from, "type", tx.Type())
			txs.Pop()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Shift()
		}
	}

	if !w.isRunning() && len(coalescedLogs) > 0 {
		// We don't push the pendingLogsEvent while we are mining. The reason is that
		// when we are mining, the worker will regenerate a mining block every 3 seconds.
		// In order to avoid pushing the repeated pendingLog, we disable the pending log pushing.

		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		w.pendingLogsFeed.Send(cpy)
	}
	// Notify resubmit loop to decrease resubmitting interval if current interval is larger
	// than the user-specified one.
	if interrupt != nil {
		w.resubmitAdjustCh <- &intervalAdjust{inc: false}
	}
	return false
}

// commitNewWork generates several new sealing tasks based on the parent block.
func (w *worker[T,P]) commitNewWork(interrupt *int32, noempty bool, timestamp int64) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	tstart := time.Now()
	parent := w.chain.CurrentBlock()

	if parent.Time() >= uint64(timestamp) {
		timestamp = int64(parent.Time() + 1)
	}
	minGasLimit := w.chainConfig.GetMinerMinGasLimit(parent.Number(), params.DefaultMinGasLimit)
	num := parent.Number()
	header := &types.Header[P]{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   core.CalcGasLimit(parent, minGasLimit, w.config.GasFloor, w.config.GasCeil),
		Extra:      w.extra,
		Time:       uint64(timestamp),
	}
	// Only set the coinbase if our consensus engine is running (avoid spurious block rewards)
	if w.isRunning() {
		if w.coinbase == (common.Address{}) {
			log.Error("Refusing to mine without etherbase")
			return
		}
		header.Coinbase = w.coinbase
	}
	if err := w.engine.Prepare(w.chain, header); err != nil {
		log.Error("Failed to prepare header for mining", "err", err)
		return
	}
	// If we are care about TheDAO hard-fork check whether to override the extra-data or not
	if daoBlock := w.chainConfig.DAOForkBlock; daoBlock != nil {
		// Check whether the block is among the fork extra-override range
		limit := new(big.Int).Add(daoBlock, params.DAOForkExtraRange)
		if header.Number.Cmp(daoBlock) >= 0 && header.Number.Cmp(limit) < 0 {
			// Depending whether we support or oppose the fork, override differently
			if w.chainConfig.DAOForkSupport {
				header.Extra = common.CopyBytes(params.DAOForkBlockExtra)
			} else if bytes.Equal(header.Extra, params.DAOForkBlockExtra) {
				header.Extra = []byte{} // If miner opposes, don't let it use the reserved extra-data
			}
		}
	}
	// Could potentially happen if starting to mine in an odd state.
	err := w.makeCurrent(parent, header)
	if err != nil {
		log.Error("Failed to create mining context", "err", err)
		return
	}
	// Create the current work task and check any fork transitions needed
	env := w.current
	if w.chainConfig.DAOForkSupport && w.chainConfig.DAOForkBlock != nil && w.chainConfig.DAOForkBlock.Cmp(header.Number) == 0 {
		misc.ApplyDAOHardFork(env.state)
	}
	// Accumulate the uncles for the current block
	uncles := make([]*types.Header[P], 0, 2)
	commitUncles := func(blocks map[common.Hash]*types.Block[P]) {
		// Clean up stale uncle blocks first
		for hash, uncle := range blocks {
			if uncle.NumberU64()+staleThreshold <= header.Number.Uint64() {
				delete(blocks, hash)
			}
		}
		for hash, uncle := range blocks {
			if len(uncles) == 2 {
				break
			}
			if err := w.commitUncle(env, uncle.Header()); err != nil {
				log.Trace("Possible uncle rejected", "hash", hash, "reason", err)
			} else {
				log.Debug("Committing new uncle to block", "hash", hash)
				uncles = append(uncles, uncle.Header())
			}
		}
	}
	// Prefer to locally generated uncle
	commitUncles(w.localUncles)
	commitUncles(w.remoteUncles)

	// Create an empty block based on temporary copied state for
	// sealing in advance without waiting block execution finished.
	if !noempty && atomic.LoadUint32(&w.noempty) == 0 {
		w.commit(uncles, nil, false, tstart)
	}

	// Fill the block with all available pending transactions.
	pending, err := w.eth.TxPool().Pending()
	if err != nil {
		log.Error("Failed to fetch pending transactions", "err", err)
		return
	}
	// Short circuit if there is no available pending transactions.
	// But if we disable empty precommit already, ignore it. Since
	// empty block is necessary to keep the liveness of the network.
	if len(pending) == 0 && atomic.LoadUint32(&w.noempty) == 0 {
		w.updateSnapshot()
		return
	}
	// Split the pending transactions into locals and remotes
	localTxs, remoteTxs := make(map[common.Address]types.Transactions[P]), pending
	for _, account := range w.eth.TxPool().Locals() {
		if txs := remoteTxs[account]; len(txs) > 0 {
			delete(remoteTxs, account)
			localTxs[account] = txs
		}
	}
	if len(localTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(w.current.signer, localTxs)
		if w.commitTransactions(txs, w.coinbase, interrupt) {
			return
		}
	}
	if len(remoteTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(w.current.signer, remoteTxs)
		if w.commitTransactions(txs, w.coinbase, interrupt) {
			return
		}
	}
	w.commit(uncles, w.fullTaskHook, true, tstart)
}

// commit runs any post-transaction state modifications, assembles the final block
// and commits new work if consensus engine is running.
func (w *worker[T,P]) commit(uncles []*types.Header[P], interval func(), update bool, start time.Time) error {
	// Deep copy receipts here to avoid interaction between different tasks.
	receipts := copyReceipts(w.current.receipts)
	privateReceipts := copyReceipts(w.current.privateReceipts) // Quorum

	s := w.current.state.Copy()
	psrCopy := w.current.privateStateRepo.Copy()
	block, err := w.engine.FinalizeAndAssemble(w.chain, w.current.header, s, w.current.txs, uncles, receipts)
	if err != nil {
		return err
	}
	if w.isRunning() {
		if interval != nil {
			interval()
		}
		select {
		case w.taskCh <- &task[P]{receipts: receipts, privateReceipts: privateReceipts, state: s, privateStateRepo: psrCopy, block: block, createdAt: time.Now()}:
			w.unconfirmed.Shift(block.NumberU64() - 1)
			log.Info("Commit new mining work", "number", block.Number(), "sealhash", w.engine.SealHash(block.Header()),
				"uncles", len(uncles), "txs", w.current.tcount,
				"gas", block.GasUsed(), "fees", totalFees(block, receipts),
				"elapsed", common.PrettyDuration(time.Since(start)))

		case <-w.exitCh:
			log.Info("Worker has exited")
		}
	}
	if update {
		w.updateSnapshot()
	}
	return nil
}

// copyReceipts makes a deep copy of the given receipts.
func copyReceipts[P crypto.PublicKey](receipts []*types.Receipt[P]) []*types.Receipt[P] {
	result := make([]*types.Receipt[P], len(receipts))
	for i, l := range receipts {
		cpy := *l
		result[i] = &cpy
	}
	return result
}

// postSideBlock fires a side chain event, only use it for testing.
func (w *worker[T,P]) postSideBlock(event core.ChainSideEvent[P]) {
	select {
	case w.chainSideCh <- event:
	case <-w.exitCh:
	}
}

// totalFees computes total consumed fees in ETH. Block transactions and receipts have to have the same order.
func totalFees[P crypto.PublicKey](block *types.Block[P], receipts []*types.Receipt[P]) *big.Float {
	feesWei := new(big.Int)
	for i, tx := range block.Transactions() {
		feesWei.Add(feesWei, new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), tx.GasPrice()))
	}
	return new(big.Float).Quo(new(big.Float).SetInt(feesWei), new(big.Float).SetInt(big.NewInt(params.Ether)))
}

// Quorum
//
// revertToPrivateStateSnapshots attempts to revert all private states to the provided snapshots
func (w *worker[T,P]) revertToPrivateStateSnapshots(snapshots map[types.PrivateStateIdentifier]int) {
	for psi, snapshot := range snapshots {
		privateState, err := w.current.privateStateRepo.StatePSI(psi)
		if err == nil {
			privateState.RevertToSnapshot(snapshot)
		} else {
			log.Warn("unable to revert to snapshot", "psi", psi, "snapshot", snapshot, "err", err)
		}
	}
}

// Quorum
//
// handling MPS scenario for a private transaction. It also captures the snapshot
// before applying the transaction
//
// handleMPS returns the auxiliary receipt and the non-nil snapshots.
//
// Caller must check for error and reverts private states.
func (w *worker[T,P]) handleMPS(tx *types.Transaction[P], coinbase common.Address, applyOnPartyOnly bool) (mpsReceipt *types.Receipt[P], privateStateSnaphots map[types.PrivateStateIdentifier]int, err error) {
	workerEnv := w.current
	privateStateRepo := workerEnv.privateStateRepo
	// make sure we don't return NIL map
	privateStateSnaphots = make(map[types.PrivateStateIdentifier]int)
	if tx.IsPrivate() && privateStateRepo.IsMPS() {
		publicStateDBFactory := func() *state.StateDB[P] {
			db := workerEnv.state.Copy()
			db.Prepare(tx.Hash(), common.Hash{}, workerEnv.tcount)
			return db
		}
		privateStateDBFactory := func(psi types.PrivateStateIdentifier) (*state.StateDB[P], error) {
			db, err := privateStateRepo.StatePSI(psi)
			if err != nil {
				return nil, err
			}
			db.Prepare(tx.Hash(), common.Hash{}, workerEnv.tcount)
			privateStateSnaphots[psi] = db.Snapshot()
			return db, nil
		}
		mpsReceipt, err = core.ApplyTransactionOnMPS[P](w.chainConfig, w.chain, &coinbase, workerEnv.gasPool, publicStateDBFactory, privateStateDBFactory, workerEnv.header, tx, &workerEnv.header.GasUsed, *w.chain.GetVMConfig(), privateStateRepo, applyOnPartyOnly, false)
	}
	return
}
