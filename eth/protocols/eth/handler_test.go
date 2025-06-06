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
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/trie"
)

var (
	// testKey is a private key to use for funding a tester account.
	testKey, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

	// testAddr is the Ethereum address of the tester account.
	testAddr = crypto.PubkeyToAddress[nist.PublicKey](*testKey.Public())
)

// testBackend is a mock implementation of the live Ethereum message handler. Its
// purpose is to allow testing the request/reply workflows and wire serialization
// in the `eth` protocol without actually doing any data processing.
type testBackend [T crypto.PrivateKey,P crypto.PublicKey] struct {
	db     ethdb.Database
	chain  *core.BlockChain[P]
	txpool *core.TxPool[P]
}

// newTestBackend creates an empty chain and wraps it into a mock backend.
func newTestBackend[T crypto.PrivateKey,P crypto.PublicKey](blocks int) *testBackend[T,P] {
	return newTestBackendWithGenerator[T,P](blocks, nil)
}

// newTestBackend creates a chain with a number of explicitly defined blocks and
// wraps it into a mock backend.
func newTestBackendWithGenerator[T crypto.PrivateKey,P crypto.PublicKey](blocks int, generator func(int, *core.BlockGen[P])) *testBackend[T,P] {
	// Create a database pre-initialize with a genesis block
	db := rawdb.NewMemoryDatabase()
	(&core.Genesis[P]{
		Config: params.TestChainConfig,
		Alloc:  core.GenesisAlloc{testAddr: {Balance: big.NewInt(1000000)}},
	}).MustCommit(db)

	chain, _ := core.NewBlockChain[P](db, nil, params.TestChainConfig,  ethash.NewFaker[P](), vm.Config[P]{}, nil, nil, nil)

	bs, _ := core.GenerateChain[P](params.TestChainConfig, chain.Genesis(),  ethash.NewFaker[P](), db, blocks, generator)
	if _, err := chain.InsertChain(bs); err != nil {
		panic(err)
	}
	txconfig := core.DefaultTxPoolConfig
	txconfig.Journal = "" // Don't litter the disk with test journals

	return &testBackend[T,P]{
		db:     db,
		chain:  chain,
		txpool: core.NewTxPool[P](txconfig, params.TestChainConfig, chain),
	}
}

// close tears down the transaction pool and chain behind the mock backend.
func (b *testBackend[T,P]) close() {
	b.txpool.Stop()
	b.chain.Stop()
}

func (b *testBackend[T,P]) Chain() *core.BlockChain[P]     { return b.chain }
func (b *testBackend[T,P]) StateBloom() *trie.SyncBloom { return nil }
func (b *testBackend[T,P]) TxPool() TxPool[P]             { return b.txpool }

func (b *testBackend[T,P]) RunPeer(peer *Peer[T,P], handler Handler[T,P]) error {
	// Normally the backend would do peer mainentance and handshakes. All that
	// is omitted and we will just give control back to the handler.
	return handler(peer)
}
func (b *testBackend[T,P]) PeerInfo(enode.ID) interface{} { panic("not implemented") }

func (b *testBackend[T,P]) AcceptTxs() bool {
	panic("data processing tests should be done in the handler package")
}
func (b *testBackend[T,P]) Handle(*Peer[T,P], Packet) error {
	panic("data processing tests should be done in the handler package")
}

// Tests that block headers can be retrieved from a remote chain based on user queries.
func TestGetBlockHeaders65(t *testing.T) { testGetBlockHeaders(t, ETH65) }
func TestGetBlockHeaders66(t *testing.T) { testGetBlockHeaders(t, ETH66) }

func testGetBlockHeaders(t *testing.T, protocol uint) {
	t.Parallel()

	backend := newTestBackend[nist.PrivateKey,nist.PublicKey](maxHeadersServe + 15)
	defer backend.close()

	peer, _ := newTestPeer[nist.PrivateKey,nist.PublicKey]("peer", protocol, backend)
	defer peer.close()

	// Create a "random" unknown hash for testing
	var unknown common.Hash
	for i := range unknown {
		unknown[i] = byte(i)
	}
	// Create a batch of tests for various scenarios
	limit := uint64(maxHeadersServe)
	tests := []struct {
		query  *GetBlockHeadersPacket // The query to execute for header retrieval
		expect []common.Hash          // The hashes of the block whose headers are expected
	}{
		// A single random block should be retrievable by hash and number too
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: backend.chain.GetBlockByNumber(limit / 2).Hash()}, Amount: 1},
			[]common.Hash{backend.chain.GetBlockByNumber(limit / 2).Hash()},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Amount: 1},
			[]common.Hash{backend.chain.GetBlockByNumber(limit / 2).Hash()},
		},
		// Multiple headers should be retrievable in both directions
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 1).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 2).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 1).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 2).Hash(),
			},
		},
		// Multiple headers with skip lists should be retrievable
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 4).Hash(),
				backend.chain.GetBlockByNumber(limit/2 + 8).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(limit / 2).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 4).Hash(),
				backend.chain.GetBlockByNumber(limit/2 - 8).Hash(),
			},
		},
		// The chain endpoints should be retrievable
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 0}, Amount: 1},
			[]common.Hash{backend.chain.GetBlockByNumber(0).Hash()},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64()}, Amount: 1},
			[]common.Hash{backend.chain.CurrentBlock().Hash()},
		},
		// Ensure protocol limits are honored
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() - 1}, Amount: limit + 10, Reverse: true},
			backend.chain.GetBlockHashesFromHash(backend.chain.CurrentBlock().Hash(), limit),
		},
		// Check that requesting more than available is handled gracefully
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() - 4}, Skip: 3, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64() - 4).Hash(),
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64()).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 4}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(4).Hash(),
				backend.chain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check that requesting more than available is handled gracefully, even if mid skip
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() - 4}, Skip: 2, Amount: 3},
			[]common.Hash{
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64() - 4).Hash(),
				backend.chain.GetBlockByNumber(backend.chain.CurrentBlock().NumberU64() - 1).Hash(),
			},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 4}, Skip: 2, Amount: 3, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(4).Hash(),
				backend.chain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check a corner case where requesting more can iterate past the endpoints
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: 2}, Amount: 5, Reverse: true},
			[]common.Hash{
				backend.chain.GetBlockByNumber(2).Hash(),
				backend.chain.GetBlockByNumber(1).Hash(),
				backend.chain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back into the chain start
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: backend.chain.GetBlockByNumber(3).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64 - 1},
			[]common.Hash{
				backend.chain.GetBlockByNumber(3).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back to the same header
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: backend.chain.GetBlockByNumber(1).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64},
			[]common.Hash{
				backend.chain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check that non existing headers aren't returned
		{
			&GetBlockHeadersPacket{Origin: HashOrNumber{Hash: unknown}, Amount: 1},
			[]common.Hash{},
		}, {
			&GetBlockHeadersPacket{Origin: HashOrNumber{Number: backend.chain.CurrentBlock().NumberU64() + 1}, Amount: 1},
			[]common.Hash{},
		},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the headers to expect in the response
		var headers []*types.Header[nist.PublicKey]
		for _, hash := range tt.expect {
			headers = append(headers, backend.chain.GetBlockByHash(hash).Header())
		}
		// Send the hash request and verify the response
		if protocol <= ETH65 {
			p2p.Send(peer.app, GetBlockHeadersMsg, tt.query)
			if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, headers); err != nil {
				t.Errorf("test %d: headers mismatch: %v", i, err)
			}
		} else {
			p2p.Send(peer.app, GetBlockHeadersMsg, GetBlockHeadersPacket66{
				RequestId:             123,
				GetBlockHeadersPacket: tt.query,
			})
			if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, BlockHeadersPacket66[nist.PublicKey]{
				RequestId:          123,
				BlockHeadersPacket: headers,
			}); err != nil {
				t.Errorf("test %d: headers mismatch: %v", i, err)
			}
		}
		// If the test used number origins, repeat with hashes as the too
		if tt.query.Origin.Hash == (common.Hash{}) {
			if origin := backend.chain.GetBlockByNumber(tt.query.Origin.Number); origin != nil {
				tt.query.Origin.Hash, tt.query.Origin.Number = origin.Hash(), 0

				if protocol <= ETH65 {
					p2p.Send(peer.app, GetBlockHeadersMsg, tt.query)
					if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, headers); err != nil {
						t.Errorf("test %d: headers mismatch: %v", i, err)
					}
				} else {
					p2p.Send(peer.app, GetBlockHeadersMsg, GetBlockHeadersPacket66{
						RequestId:             456,
						GetBlockHeadersPacket: tt.query,
					})
					if err := p2p.ExpectMsg(peer.app, BlockHeadersMsg, BlockHeadersPacket66[nist.PublicKey]{
						RequestId:          456,
						BlockHeadersPacket: headers,
					}); err != nil {
						t.Errorf("test %d: headers mismatch: %v", i, err)
					}
				}
			}
		}
	}
}

// Tests that block contents can be retrieved from a remote chain based on their hashes.
func TestGetBlockBodies65(t *testing.T) { testGetBlockBodies(t, ETH65) }
func TestGetBlockBodies66(t *testing.T) { testGetBlockBodies(t, ETH66) }

func testGetBlockBodies(t *testing.T, protocol uint) {
	t.Parallel()

	backend := newTestBackend[nist.PrivateKey, nist.PublicKey](maxBodiesServe + 15)
	defer backend.close()

	peer, _ := newTestPeer[nist.PrivateKey,nist.PublicKey]("peer", protocol, backend)
	defer peer.close()

	// Create a batch of tests for various scenarios
	limit := maxBodiesServe
	tests := []struct {
		random    int           // Number of blocks to fetch randomly from the chain
		explicit  []common.Hash // Explicitly requested blocks
		available []bool        // Availability of explicitly requested blocks
		expected  int           // Total number of existing blocks to expect
	}{
		{1, nil, nil, 1},             // A single random block should be retrievable
		{10, nil, nil, 10},           // Multiple random blocks should be retrievable
		{limit, nil, nil, limit},     // The maximum possible blocks should be retrievable
		{limit + 1, nil, nil, limit}, // No more than the possible block count should be returned
		{0, []common.Hash{backend.chain.Genesis().Hash()}, []bool{true}, 1},      // The genesis block should be retrievable
		{0, []common.Hash{backend.chain.CurrentBlock().Hash()}, []bool{true}, 1}, // The chains head block should be retrievable
		{0, []common.Hash{{}}, []bool{false}, 0},                                 // A non existent block should not be returned

		// Existing and non-existing blocks interleaved should not cause problems
		{0, []common.Hash{
			{},
			backend.chain.GetBlockByNumber(1).Hash(),
			{},
			backend.chain.GetBlockByNumber(10).Hash(),
			{},
			backend.chain.GetBlockByNumber(100).Hash(),
			{},
		}, []bool{false, true, false, true, false, true, false}, 3},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the hashes to request, and the response to expectva
		var (
			hashes []common.Hash
			bodies []*BlockBody[nist.PublicKey]
			seen   = make(map[int64]bool)
		)
		for j := 0; j < tt.random; j++ {
			for {
				num := rand.Int63n(int64(backend.chain.CurrentBlock().NumberU64()))
				if !seen[num] {
					seen[num] = true

					block := backend.chain.GetBlockByNumber(uint64(num))
					hashes = append(hashes, block.Hash())
					if len(bodies) < tt.expected {
						bodies = append(bodies, &BlockBody[nist.PublicKey]{Transactions: block.Transactions(), Uncles: block.Uncles()})
					}
					break
				}
			}
		}
		for j, hash := range tt.explicit {
			hashes = append(hashes, hash)
			if tt.available[j] && len(bodies) < tt.expected {
				block := backend.chain.GetBlockByHash(hash)
				bodies = append(bodies, &BlockBody[nist.PublicKey]{Transactions: block.Transactions(), Uncles: block.Uncles()})
			}
		}
		// Send the hash request and verify the response
		if protocol <= ETH65 {
			p2p.Send(peer.app, GetBlockBodiesMsg, hashes)
			if err := p2p.ExpectMsg(peer.app, BlockBodiesMsg, bodies); err != nil {
				t.Errorf("test %d: bodies mismatch: %v", i, err)
			}
		} else {
			p2p.Send(peer.app, GetBlockBodiesMsg, GetBlockBodiesPacket66{
				RequestId:            123,
				GetBlockBodiesPacket: hashes,
			})
			if err := p2p.ExpectMsg(peer.app, BlockBodiesMsg, BlockBodiesPacket66[nist.PublicKey]{
				RequestId:         123,
				BlockBodiesPacket: bodies,
			}); err != nil {
				t.Errorf("test %d: bodies mismatch: %v", i, err)
			}
		}
	}
}

// Tests that the state trie nodes can be retrieved based on hashes.
func TestGetNodeData65(t *testing.T) { testGetNodeData(t, ETH65) }
func TestGetNodeData66(t *testing.T) { testGetNodeData(t, ETH66) }

func testGetNodeData(t *testing.T, protocol uint) {
	t.Parallel()

	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA[nist.PrivateKey]("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA[nist.PrivateKey]("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress[nist.PublicKey](*acc1Key.Public())
	acc2Addr := crypto.PubkeyToAddress[nist.PublicKey](*acc2Key.Public())

	signer := types.HomesteadSigner[nist.PublicKey]{}
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen[nist.PublicKey]) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(testAddr), acc1Addr, big.NewInt(10000), params.TxGas, nil, nil), signer, testKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(testAddr), acc1Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, testKey)
			tx2, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			block.AddUncle(b2)
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
			block.AddUncle(b3)
		}
	}
	// Assemble the test environment
	backend := newTestBackendWithGenerator[nist.PrivateKey,nist.PublicKey](4, generator)
	defer backend.close()

	peer, _ := newTestPeer[nist.PrivateKey,nist.PublicKey]("peer", protocol, backend)
	defer peer.close()

	// Fetch for now the entire chain db
	var hashes []common.Hash

	it := backend.db.NewIterator(nil, nil)
	for it.Next() {
		if key := it.Key(); len(key) == common.HashLength {
			hashes = append(hashes, common.BytesToHash(key))
		}
	}
	it.Release()

	if protocol <= ETH65 {
		p2p.Send(peer.app, GetNodeDataMsg, hashes)
	} else {
		p2p.Send(peer.app, GetNodeDataMsg, GetNodeDataPacket66{
			RequestId:         123,
			GetNodeDataPacket: hashes,
		})
	}
	msg, err := peer.app.ReadMsg()
	if err != nil {
		t.Fatalf("failed to read node data response: %v", err)
	}
	if msg.Code != NodeDataMsg {
		t.Fatalf("response packet code mismatch: have %x, want %x", msg.Code, NodeDataMsg)
	}
	var data [][]byte
	if protocol <= ETH65 {
		if err := msg.Decode(&data); err != nil {
			t.Fatalf("failed to decode response node data: %v", err)
		}
	} else {
		var res NodeDataPacket66
		if err := msg.Decode(&res); err != nil {
			t.Fatalf("failed to decode response node data: %v", err)
		}
		data = res.NodeDataPacket
	}
	// Verify that all hashes correspond to the requested data, and reconstruct a state tree
	for i, want := range hashes {
		if hash := crypto.Keccak256Hash[nist.PublicKey](data[i]); hash != want {
			t.Errorf("data hash mismatch: have %x, want %x", hash, want)
		}
	}
	statedb := rawdb.NewMemoryDatabase()
	for i := 0; i < len(data); i++ {
		statedb.Put(hashes[i].Bytes(), data[i])
	}
	accounts := []common.Address{testAddr, acc1Addr, acc2Addr}
	for i := uint64(0); i <= backend.chain.CurrentBlock().NumberU64(); i++ {
		trie, _ := state.New[nist.PublicKey](backend.chain.GetBlockByNumber(i).Root(), state.NewDatabase[nist.PublicKey](statedb), nil)

		for j, acc := range accounts {
			state, _, _ := backend.chain.State()
			bw := state.GetBalance(acc)
			bh := trie.GetBalance(acc)

			if (bw != nil && bh == nil) || (bw == nil && bh != nil) {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
			if bw != nil && bh != nil && bw.Cmp(bw) != 0 {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
		}
	}
}

// Tests that the transaction receipts can be retrieved based on hashes.
func TestGetBlockReceipts65(t *testing.T) { testGetBlockReceipts(t, ETH65) }
func TestGetBlockReceipts66(t *testing.T) { testGetBlockReceipts(t, ETH66) }

func testGetBlockReceipts(t *testing.T, protocol uint) {
	t.Parallel()

	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA[nist.PrivateKey]("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA[nist.PrivateKey]("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress[nist.PublicKey](*acc1Key.Public())
	acc2Addr := crypto.PubkeyToAddress[nist.PublicKey](*acc2Key.Public())

	signer := types.HomesteadSigner[nist.PublicKey]{}
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen[nist.PublicKey]) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(testAddr), acc1Addr, big.NewInt(10000), params.TxGas, nil, nil), signer, testKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(testAddr), acc1Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, testKey)
			tx2, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewTransaction[nist.PublicKey](block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			block.AddUncle(b2)
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
			block.AddUncle(b3)
		}
	}
	// Assemble the test environment
	backend := newTestBackendWithGenerator[nist.PrivateKey,nist.PublicKey](4, generator)
	defer backend.close()

	peer, _ := newTestPeer[nist.PrivateKey,nist.PublicKey]("peer", protocol, backend)
	defer peer.close()

	// Collect the hashes to request, and the response to expect
	var (
		hashes   []common.Hash
		receipts [][]*types.Receipt[nist.PublicKey]
	)
	for i := uint64(0); i <= backend.chain.CurrentBlock().NumberU64(); i++ {
		block := backend.chain.GetBlockByNumber(i)

		hashes = append(hashes, block.Hash())
		receipts = append(receipts, backend.chain.GetReceiptsByHash(block.Hash()))
	}
	// Send the hash request and verify the response
	if protocol <= ETH65 {
		p2p.Send(peer.app, GetReceiptsMsg, hashes)
		if err := p2p.ExpectMsg(peer.app, ReceiptsMsg, receipts); err != nil {
			t.Errorf("receipts mismatch: %v", err)
		}
	} else {
		p2p.Send(peer.app, GetReceiptsMsg, GetReceiptsPacket66{
			RequestId:         123,
			GetReceiptsPacket: hashes,
		})
		if err := p2p.ExpectMsg(peer.app, ReceiptsMsg, ReceiptsPacket66[nist.PublicKey]{
			RequestId:      123,
			ReceiptsPacket: receipts,
		}); err != nil {
			t.Errorf("receipts mismatch: %v", err)
		}
	}
}
