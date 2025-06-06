// Copyright 2016 The go-ethereum Authors
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

package core

import (
	"container/heap"
	"math"
	"math/big"
	"sort"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// nonceHeap is a heap.Interface implementation over 64bit unsigned integers for
// retrieving sorted transactions from the possibly gapped future queue.
type nonceHeap []uint64

func (h nonceHeap) Len() int           { return len(h) }
func (h nonceHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h nonceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nonceHeap) Push(x interface{}) {
	*h = append(*h, x.(uint64))
}

func (h *nonceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// txSortedMap is a nonce->transaction hash map with a heap based index to allow
// iterating over the contents in a nonce-incrementing way.
type txSortedMap [P crypto.PublicKey] struct {
	items map[uint64]*types.Transaction[P] // Hash map storing the transaction data
	index *nonceHeap                    // Heap of nonces of all the stored transactions (non-strict mode)
	cache types.Transactions[P]          // Cache of the transactions already sorted
}

// newTxSortedMap creates a new nonce-sorted transaction map.
func newTxSortedMap[P crypto.PublicKey]() *txSortedMap[P] {
	return &txSortedMap[P]{
		items: make(map[uint64]*types.Transaction[P]),
		index: new(nonceHeap),
	}
}

// Get retrieves the current transactions associated with the given nonce.
func (m *txSortedMap[P]) Get(nonce uint64) *types.Transaction[P] {
	return m.items[nonce]
}

// Put inserts a new transaction into the map, also updating the map's nonce
// index. If a transaction already exists with the same nonce, it's overwritten.
func (m *txSortedMap[P]) Put(tx *types.Transaction[P]) {
	nonce := tx.Nonce()
	if m.items[nonce] == nil {
		heap.Push(m.index, nonce)
	}
	m.items[nonce], m.cache = tx, nil
}

// Forward removes all transactions from the map with a nonce lower than the
// provided threshold. Every removed transaction is returned for any post-removal
// maintenance.
func (m *txSortedMap[P]) Forward(threshold uint64) types.Transactions[P] {
	var removed types.Transactions[P]

	// Pop off heap items until the threshold is reached
	for m.index.Len() > 0 && (*m.index)[0] < threshold {
		nonce := heap.Pop(m.index).(uint64)
		removed = append(removed, m.items[nonce])
		delete(m.items, nonce)
	}
	// If we had a cached order, shift the front
	if m.cache != nil {
		m.cache = m.cache[len(removed):]
	}
	return removed
}

// Filter iterates over the list of transactions and removes all of them for which
// the specified function evaluates to true.
// Filter, as opposed to 'filter', re-initialises the heap after the operation is done.
// If you want to do several consecutive filterings, it's therefore better to first
// do a .filter(func1) followed by .Filter(func2) or reheap()
func (m *txSortedMap[P]) Filter(filter func(*types.Transaction[P]) bool) types.Transactions[P] {
	removed := m.filter(filter)
	// If transactions were removed, the heap and cache are ruined
	if len(removed) > 0 {
		m.reheap()
	}
	return removed
}

func (m *txSortedMap[P]) reheap() {
	*m.index = make([]uint64, 0, len(m.items))
	for nonce := range m.items {
		*m.index = append(*m.index, nonce)
	}
	heap.Init(m.index)
	m.cache = nil
}

// filter is identical to Filter, but **does not** regenerate the heap. This method
// should only be used if followed immediately by a call to Filter or reheap()
func (m *txSortedMap[P]) filter(filter func(*types.Transaction[P]) bool) types.Transactions[P] {
	var removed types.Transactions[P]

	// Collect all the transactions to filter out
	for nonce, tx := range m.items {
		if filter(tx) {
			removed = append(removed, tx)
			delete(m.items, nonce)
		}
	}
	if len(removed) > 0 {
		m.cache = nil
	}
	return removed
}

// Cap places a hard limit on the number of items, returning all transactions
// exceeding that limit.
func (m *txSortedMap[P]) Cap(threshold int) types.Transactions[P] {
	// Short circuit if the number of items is under the limit
	if len(m.items) <= threshold {
		return nil
	}
	// Otherwise gather and drop the highest nonce'd transactions
	var drops types.Transactions[P]

	sort.Sort(*m.index)
	for size := len(m.items); size > threshold; size-- {
		drops = append(drops, m.items[(*m.index)[size-1]])
		delete(m.items, (*m.index)[size-1])
	}
	*m.index = (*m.index)[:threshold]
	heap.Init(m.index)

	// If we had a cache, shift the back
	if m.cache != nil {
		m.cache = m.cache[:len(m.cache)-len(drops)]
	}
	return drops
}

// Remove deletes a transaction from the maintained map, returning whether the
// transaction was found.
func (m *txSortedMap[P]) Remove(nonce uint64) bool {
	// Short circuit if no transaction is present
	_, ok := m.items[nonce]
	if !ok {
		return false
	}
	// Otherwise delete the transaction and fix the heap index
	for i := 0; i < m.index.Len(); i++ {
		if (*m.index)[i] == nonce {
			heap.Remove(m.index, i)
			break
		}
	}
	delete(m.items, nonce)
	m.cache = nil

	return true
}

// Ready retrieves a sequentially increasing list of transactions starting at the
// provided nonce that is ready for processing. The returned transactions will be
// removed from the list.
//
// Note, all transactions with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (m *txSortedMap[P]) Ready(start uint64) types.Transactions[P]{
	// Short circuit if no transactions are available
	if m.index.Len() == 0 || (*m.index)[0] > start {
		return nil
	}
	// Otherwise start accumulating incremental transactions
	var ready types.Transactions[P]
	for next := (*m.index)[0]; m.index.Len() > 0 && (*m.index)[0] == next; next++ {
		ready = append(ready, m.items[next])
		delete(m.items, next)
		heap.Pop(m.index)
	}
	m.cache = nil

	return ready
}

// Len returns the length of the transaction map.
func (m *txSortedMap[P]) Len() int {
	return len(m.items)
}

func (m *txSortedMap[P]) flatten() types.Transactions[P]{
	// If the sorting was not cached yet, create and cache it
	if m.cache == nil {
		m.cache = make(types.Transactions[P], 0, len(m.items))
		for _, tx := range m.items {
			m.cache = append(m.cache, tx)
		}
		sort.Sort(types.TxByNonce[P](m.cache))
	}
	return m.cache
}

// Flatten creates a nonce-sorted slice of transactions based on the loosely
// sorted internal representation. The result of the sorting is cached in case
// it's requested again before any modifications are made to the contents.
func (m *txSortedMap[P]) Flatten() types.Transactions[P]{
	// Copy the cache to prevent accidental modifications
	cache := m.flatten()
	txs := make(types.Transactions[P], len(cache))
	copy(txs, cache)
	return txs
}

// LastElement returns the last element of a flattened list, thus, the
// transaction with the highest nonce
func (m *txSortedMap[P]) LastElement() *types.Transaction[P] {
	cache := m.flatten()
	return cache[len(cache)-1]
}

// txList is a "list" of transactions belonging to an account, sorted by account
// nonce. The same type can be used both for storing contiguous transactions for
// the executable/pending queue; and for storing gapped transactions for the non-
// executable/future queue, with minor behavioral changes.
type txList [P crypto.PublicKey] struct {
	strict bool         // Whether nonces are strictly continuous or not
	txs    *txSortedMap[P] // Heap indexed sorted hash map of the transactions

	costcap *big.Int // Price of the highest costing transaction (reset only if exceeds balance)
	gascap  uint64   // Gas limit of the highest spending transaction (reset only if exceeds block limit)
}

// newTxList create a new transaction list for maintaining nonce-indexable fast,
// gapped, sortable transaction lists.
func newTxList [P crypto.PublicKey](strict bool) *txList[P] {
	return &txList[P]{
		strict:  strict,
		txs:     newTxSortedMap[P](),
		costcap: new(big.Int),
	}
}

// Overlaps returns whether the transaction specified has the same nonce as one
// already contained within the list.
func (l *txList[P]) Overlaps(tx *types.Transaction[P]) bool {
	return l.txs.Get(tx.Nonce()) != nil
}

// Add tries to insert a new transaction into the list, returning whether the
// transaction was accepted, and if yes, any previous transaction it replaced.
//
// If the new transaction is accepted into the list, the lists' cost and gas
// thresholds are also potentially updated.
func (l *txList[P]) Add(tx *types.Transaction[P], priceBump uint64) (bool, *types.Transaction[P]) {
	// If there's an older better transaction, abort
	old := l.txs.Get(tx.Nonce())
	if old != nil {
		// threshold = oldGP * (100 + priceBump) / 100
		a := big.NewInt(100 + int64(priceBump))
		a = a.Mul(a, old.GasPrice())
		b := big.NewInt(100)
		threshold := a.Div(a, b)
		// Have to ensure that the new gas price is higher than the old gas
		// price as well as checking the percentage threshold to ensure that
		// this is accurate for low (Wei-level) gas price replacements
		if old.GasPriceCmp(tx) >= 0 || tx.GasPriceIntCmp(threshold) < 0 {
			return false, nil
		}
	}
	// Otherwise overwrite the old transaction with the current one
	l.txs.Put(tx)
	if cost := tx.Cost(); l.costcap.Cmp(cost) < 0 {
		l.costcap = cost
	}
	if gas := tx.Gas(); l.gascap < gas {
		l.gascap = gas
	}
	return true, old
}

// Forward removes all transactions from the list with a nonce lower than the
// provided threshold. Every removed transaction is returned for any post-removal
// maintenance.
func (l *txList[P]) Forward(threshold uint64) types.Transactions[P]{
	return l.txs.Forward(threshold)
}

// Filter removes all transactions from the list with a cost or gas limit higher
// than the provided thresholds. Every removed transaction is returned for any
// post-removal maintenance. Strict-mode invalidated transactions are also
// returned.
//
// This method uses the cached costcap and gascap to quickly decide if there's even
// a point in calculating all the costs or if the balance covers all. If the threshold
// is lower than the costgas cap, the caps will be reset to a new high after removing
// the newly invalidated transactions.
func (l *txList[P]) Filter(costLimit *big.Int, gasLimit uint64) (types.Transactions[P], types.Transactions[P]) {
	// If all transactions are below the threshold, short circuit
	if l.costcap.Cmp(costLimit) <= 0 && l.gascap <= gasLimit {
		return nil, nil
	}
	l.costcap = new(big.Int).Set(costLimit) // Lower the caps to the thresholds
	l.gascap = gasLimit

	// Filter out all the transactions above the account's funds
	removed := l.txs.Filter(func(tx *types.Transaction[P]) bool {
		return tx.Gas() > gasLimit || tx.Cost().Cmp(costLimit) > 0
	})

	if len(removed) == 0 {
		return nil, nil
	}
	var invalids types.Transactions[P]
	// If the list was strict, filter anything above the lowest nonce
	if l.strict {
		lowest := uint64(math.MaxUint64)
		for _, tx := range removed {
			if nonce := tx.Nonce(); lowest > nonce {
				lowest = nonce
			}
		}
		invalids = l.txs.filter(func(tx *types.Transaction[P]) bool { return tx.Nonce() > lowest })
	}
	l.txs.reheap()
	return removed, invalids
}

// Cap places a hard limit on the number of items, returning all transactions
// exceeding that limit.
func (l *txList[P]) Cap(threshold int) types.Transactions[P]{
	return l.txs.Cap(threshold)
}

// Remove deletes a transaction from the maintained list, returning whether the
// transaction was found, and also returning any transaction invalidated due to
// the deletion (strict mode only).
func (l *txList[P]) Remove(tx *types.Transaction[P]) (bool, types.Transactions[P]) {
	// Remove the transaction from the set
	nonce := tx.Nonce()
	if removed := l.txs.Remove(nonce); !removed {
		return false, nil
	}
	// In strict mode, filter out non-executable transactions
	if l.strict {
		return true, l.txs.Filter(func(tx *types.Transaction[P]) bool { return tx.Nonce() > nonce })
	}
	return true, nil
}

// Ready retrieves a sequentially increasing list of transactions starting at the
// provided nonce that is ready for processing. The returned transactions will be
// removed from the list.
//
// Note, all transactions with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (l *txList[P]) Ready(start uint64) types.Transactions[P]{
	return l.txs.Ready(start)
}

// Len returns the length of the transaction list.
func (l *txList[P]) Len() int {
	return l.txs.Len()
}

// Empty returns whether the list of transactions is empty or not.
func (l *txList[P]) Empty() bool {
	return l.Len() == 0
}

// Flatten creates a nonce-sorted slice of transactions based on the loosely
// sorted internal representation. The result of the sorting is cached in case
// it's requested again before any modifications are made to the contents.
func (l *txList[P]) Flatten() types.Transactions[P]{
	return l.txs.Flatten()
}

// LastElement returns the last element of a flattened list, thus, the
// transaction with the highest nonce
func (l *txList[P]) LastElement() *types.Transaction[P] {
	return l.txs.LastElement()
}

// priceHeap is a heap.Interface implementation over transactions for retrieving
// price-sorted transactions to discard when the pool fills up.
type priceHeap[P crypto.PublicKey] []*types.Transaction[P]

func (h priceHeap[P]) Len() int      { return len(h) }
func (h priceHeap[P]) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h priceHeap[P]) Less(i, j int) bool {
	// Sort primarily by price, returning the cheaper one
	switch h[i].GasPriceCmp(h[j]) {
	case -1:
		return true
	case 1:
		return false
	}
	// If the prices match, stabilize via nonces (high nonce is worse)
	return h[i].Nonce() > h[j].Nonce()
}

func (h *priceHeap[P]) Push(x interface{}) {
	*h = append(*h, x.(*types.Transaction[P]))
}

func (h *priceHeap[P]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return x
}

// txPricedList is a price-sorted heap to allow operating on transactions pool
// contents in a price-incrementing way. It's built opon the all transactions
// in txpool but only interested in the remote part. It means only remote transactions
// will be considered for tracking, sorting, eviction, etc.
type txPricedList [P crypto.PublicKey] struct {
	all     *txLookup[P]  // Pointer to the map of all transactions
	remotes *priceHeap[P] // Heap of prices of all the stored **remote** transactions
	stales  int        // Number of stale price points to (re-heap trigger)
}

// newTxPricedList creates a new price-sorted transaction heap.
func newTxPricedList[P crypto.PublicKey](all *txLookup[P]) *txPricedList[P] {
	return &txPricedList[P]{
		all:     all,
		remotes: new(priceHeap[P]),
	}
}

// Put inserts a new transaction into the heap.
func (l *txPricedList[P]) Put(tx *types.Transaction[P], local bool) {
	if local {
		return
	}
	heap.Push(l.remotes, tx)
}

// Removed notifies the prices transaction list that an old transaction dropped
// from the pool. The list will just keep a counter of stale objects and update
// the heap if a large enough ratio of transactions go stale.
func (l *txPricedList[P]) Removed(count int) {
	// Bump the stale counter, but exit if still too low (< 25%)
	l.stales += count
	if l.stales <= len(*l.remotes)/4 {
		return
	}
	// Seems we've reached a critical number of stale transactions, reheap
	l.Reheap()
}

// Cap finds all the transactions below the given price threshold, drops them
// from the priced list and returns them for further removal from the entire pool.
//
// Note: only remote transactions will be considered for eviction.
func (l *txPricedList[P]) Cap(threshold *big.Int) types.Transactions[P]{
	drop := make(types.Transactions[P], 0, 128) // Remote underpriced transactions to drop
	for len(*l.remotes) > 0 {
		// Discard stale transactions if found during cleanup
		cheapest := (*l.remotes)[0]
		if l.all.GetRemote(cheapest.Hash()) == nil { // Removed or migrated
			heap.Pop(l.remotes)
			l.stales--
			continue
		}
		// Stop the discards if we've reached the threshold
		if cheapest.GasPriceIntCmp(threshold) >= 0 {
			break
		}
		heap.Pop(l.remotes)
		drop = append(drop, cheapest)
	}
	return drop
}

// Underpriced checks whether a transaction is cheaper than (or as cheap as) the
// lowest priced (remote) transaction currently being tracked.
func (l *txPricedList[P]) Underpriced(tx *types.Transaction[P]) bool {
	// Discard stale price points if found at the heap start
	for len(*l.remotes) > 0 {
		head := []*types.Transaction[P](*l.remotes)[0]
		if l.all.GetRemote(head.Hash()) == nil { // Removed or migrated
			l.stales--
			heap.Pop(l.remotes)
			continue
		}
		break
	}
	// Check if the transaction is underpriced or not
	if len(*l.remotes) == 0 {
		return false // There is no remote transaction at all.
	}
	// If the remote transaction is even cheaper than the
	// cheapest one tracked locally, reject it.
	cheapest := []*types.Transaction[P](*l.remotes)[0]
	return cheapest.GasPriceCmp(tx) >= 0
}

// Discard finds a number of most underpriced transactions, removes them from the
// priced list and returns them for further removal from the entire pool.
//
// Note local transaction won't be considered for eviction.
func (l *txPricedList[P]) Discard(slots int, force bool) (types.Transactions[P], bool) {
	drop := make(types.Transactions[P], 0, slots) // Remote underpriced transactions to drop
	for len(*l.remotes) > 0 && slots > 0 {
		// Discard stale transactions if found during cleanup
		tx := heap.Pop(l.remotes).(*types.Transaction[P])
		if l.all.GetRemote(tx.Hash()) == nil { // Removed or migrated
			l.stales--
			continue
		}
		// Non stale transaction found, discard it
		drop = append(drop, tx)
		slots -= numSlots(tx)
	}
	// If we still can't make enough room for the new transaction
	if slots > 0 && !force {
		for _, tx := range drop {
			heap.Push(l.remotes, tx)
		}
		return nil, false
	}
	return drop, true
}

// Reheap forcibly rebuilds the heap based on the current remote transaction set.
func (l *txPricedList[P]) Reheap() {
	reheap := make(priceHeap[P], 0, l.all.RemoteCount())

	l.stales, l.remotes = 0, &reheap
	l.all.Range(func(hash common.Hash, tx *types.Transaction[P], local bool) bool {
		*l.remotes = append(*l.remotes, tx)
		return true
	}, false, true) // Only iterate remotes
	heap.Init(l.remotes)
}
