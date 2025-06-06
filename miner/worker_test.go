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

package miner

import (
	"encoding/base64"
	"math/big"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/clique"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/private/engine"
	"github.com/stretchr/testify/assert"
)

const (
	// testCode is the testing contract binary code which will initialises some
	// variables in constructor
	testCode = "0x60806040527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0060005534801561003457600080fd5b5060fc806100436000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630c4dae8814603757806398a213cf146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506084565b005b60005481565b806000819055507fe9e44f9f7da8c559de847a3232b57364adc0354f15a2cd8dc636d54396f9587a6000546040518082815260200191505060405180910390a15056fea265627a7a723058208ae31d9424f2d0bc2a3da1a5dd659db2d71ec322a17db8f87e19e209e3a1ff4a64736f6c634300050a0032"

	// testGas is the gas required for contract deployment.
	testGas = 144109
)

var (
	// Test chain configurations
	testTxPoolConfig  core.TxPoolConfig
	ethashChainConfig *params.ChainConfig
	cliqueChainConfig *params.ChainConfig

	// Test accounts
	testBankKey, _  = crypto.GenerateKey[nist.PrivateKey]()
	testBankAddress = crypto.PubkeyToAddress[nist.PublicKey](*testBankKey.Public())
	testBankFunds   = big.NewInt(1000000000000000000)

	testUserKey, _  = crypto.GenerateKey[nist.PrivateKey]()
	testUserAddress = crypto.PubkeyToAddress[nist.PublicKey](*testUserKey.Public())

	// Test transactions
	pendingTxs []*types.Transaction[nist.PublicKey]
	newTxs     []*types.Transaction[nist.PublicKey]

	testConfig = &Config{
		Recommit: time.Second,
		GasFloor: params.GenesisGasLimit,
		GasCeil:  params.GenesisGasLimit,
	}
)

func init() {
	testTxPoolConfig = core.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""
	ethashChainConfig = params.TestChainConfig
	cliqueChainConfig = params.TestChainConfig
	cliqueChainConfig.Clique = &params.CliqueConfig{
		Period: 10,
		Epoch:  30000,
	}

	signer := types.LatestSigner[nist.PublicKey](params.TestChainConfig)
	tx1 := types.MustSignNewTx(testBankKey, signer, &types.AccessListTx{
		ChainID: params.TestChainConfig.ChainID,
		Nonce:   0,
		To:      &testUserAddress,
		Value:   big.NewInt(1000),
		Gas:     params.TxGas,
	})
	pendingTxs = append(pendingTxs, tx1)

	tx2 := types.MustSignNewTx(testBankKey, signer, &types.LegacyTx{
		Nonce: 1,
		To:    &testUserAddress,
		Value: big.NewInt(1000),
		Gas:   params.TxGas,
	})
	newTxs = append(newTxs, tx2)

	rand.Seed(time.Now().UnixNano())
}

// testWorkerBackend implements worker.Backend interfaces and wraps all information needed during the testing.
type testWorkerBackend [T crypto.PrivateKey,P crypto.PublicKey] struct {
	db         ethdb.Database
	txPool     *core.TxPool[P]
	chain      *core.BlockChain[P]
	testTxFeed event.Feed
	genesis    *core.Genesis[P]
	uncleBlock *types.Block[P]
}

func newTestWorkerBackend[T crypto.PrivateKey,P crypto.PublicKey](t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine[P], db ethdb.Database, n int) *testWorkerBackend[T,P] {
	var gspec = core.Genesis[P]{
		Config: chainConfig,
		Alloc:  core.GenesisAlloc{testBankAddress: {Balance: testBankFunds}},
	}

	switch e := engine.(type) {
	case *clique.Clique[P]:
		gspec.ExtraData = make([]byte, 32+common.AddressLength+crypto.SignatureLength)
		copy(gspec.ExtraData[32:32+common.AddressLength], testBankAddress.Bytes())
		e.Authorize(testBankAddress, func(account accounts.Account, s string, data []byte) ([]byte, error) {
			return crypto.Sign(crypto.Keccak256[P](data), testBankKey)
		})
	case *ethash.Ethash[P]:
	default:
		t.Fatalf("unexpected consensus engine type: %T", engine)
	}
	genesis := gspec.MustCommit(db)

	chain, _ := core.NewBlockChain[P](db, &core.CacheConfig{TrieDirtyDisabled: true}, gspec.Config, engine, vm.Config[P]{}, nil, nil, nil)
	txpool := core.NewTxPool[P](testTxPoolConfig, chainConfig, chain)

	// Generate a small n-block chain and an uncle block for it
	if n > 0 {
		blocks, _ := core.GenerateChain[P](chainConfig, genesis, engine, db, n, func(i int, gen *core.BlockGen[P]) {
			gen.SetCoinbase(testBankAddress)
		})
		if _, err := chain.InsertChain(blocks); err != nil {
			t.Fatalf("failed to insert origin chain: %v", err)
		}
	}
	parent := genesis
	if n > 0 {
		parent = chain.GetBlockByHash(chain.CurrentBlock().ParentHash())
	}
	blocks, _ := core.GenerateChain[P](chainConfig, parent, engine, db, 1, func(i int, gen *core.BlockGen[P]) {
		gen.SetCoinbase(testUserAddress)
	})

	return &testWorkerBackend[T,P]{
		db:         db,
		chain:      chain,
		txPool:     txpool,
		genesis:    &gspec,
		uncleBlock: blocks[0],
	}
}

func (b *testWorkerBackend[T,P]) ChainDb() ethdb.Database      { return b.db }
func (b *testWorkerBackend[T,P]) BlockChain() *core.BlockChain[P] { return b.chain }
func (b *testWorkerBackend[T,P]) TxPool() *core.TxPool[P]         { return b.txPool }

func (b *testWorkerBackend[T,P]) newRandomUncle() *types.Block[P] {
	var parent *types.Block[P]
	cur := b.chain.CurrentBlock()
	if cur.NumberU64() == 0 {
		parent = b.chain.Genesis()
	} else {
		parent = b.chain.GetBlockByHash(b.chain.CurrentBlock().ParentHash())
	}
	blocks, _ := core.GenerateChain[P](b.chain.Config(), parent, b.chain.Engine(), b.db, 1, func(i int, gen *core.BlockGen[P]) {
		var addr = make([]byte, common.AddressLength)
		rand.Read(addr)
		gen.SetCoinbase(common.BytesToAddress(addr))
	})
	return blocks[0]
}

func (b *testWorkerBackend[T,P]) newRandomTx(creation bool, private bool) *types.Transaction[P] {
	var signer types.Signer[P]
	signer = types.HomesteadSigner[P]{}
	if private {
		signer = types.QuorumPrivateTxSigner[P]{}
	}
	var tx *types.Transaction[P]
	var key T
	switch t:=any(&testBankKey).(type){
	case *nist.PrivateKey:
		tt:=any(&key).(*nist.PrivateKey)
		*tt=*t
	case *gost3410.PrivateKey:
		tt:=any(&key).(*gost3410.PrivateKey)
		*tt=*t
	}
	if creation {
		tx, _ = types.SignTx[T,P](types.NewContractCreation[P](b.txPool.Nonce(testBankAddress), big.NewInt(0), testGas, nil, common.FromHex(testCode)), signer, key)
	} else {
		tx, _ = types.SignTx[T,P](types.NewTransaction[P](b.txPool.Nonce(testBankAddress), testUserAddress, big.NewInt(1000), params.TxGas, nil, nil), signer, key)
	}
	return tx
}

func newTestWorker(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine[nist.PublicKey], db ethdb.Database, blocks int) (*worker[nist.PrivateKey,nist.PublicKey], *testWorkerBackend[nist.PrivateKey,nist.PublicKey]) {
	backend := newTestWorkerBackend[nist.PrivateKey,nist.PublicKey](t, chainConfig, engine, db, blocks)
	backend.txPool.AddLocals(pendingTxs)
	w := newWorker[nist.PrivateKey,nist.PublicKey](testConfig, chainConfig, engine, backend, new(event.TypeMux), nil, false)
	w.setEtherbase(testBankAddress)
	return w, backend
}

func TestGenerateBlockAndImportEthash(t *testing.T) {
	testGenerateBlockAndImport(t, false)
}

func TestGenerateBlockAndImportClique(t *testing.T) {
	testGenerateBlockAndImport(t, true)
}

func testGenerateBlockAndImport(t *testing.T, isClique bool) {
	var (
		engine      consensus.Engine[nist.PublicKey]
		chainConfig *params.ChainConfig
		db          = rawdb.NewMemoryDatabase()
	)
	if isClique {
		chainConfig = params.AllCliqueProtocolChanges
		chainConfig.Clique = &params.CliqueConfig{Period: 1, Epoch: 30000}
		engine = clique.New[nist.PublicKey](chainConfig.Clique, db)
	} else {
		chainConfig = params.AllEthashProtocolChanges
		engine =  ethash.NewFaker[nist.PublicKey]()
	}

	w, b := newTestWorker(t, chainConfig, engine, db, 0)
	defer w.close()

	// This test chain imports the mined blocks.
	db2 := rawdb.NewMemoryDatabase()
	b.genesis.MustCommit(db2)
	chain, _ := core.NewBlockChain[nist.PublicKey](db2, nil, b.chain.Config(), engine, vm.Config[nist.PublicKey]{}, nil, nil, nil)
	defer chain.Stop()

	// Ignore empty commit here for less noise.
	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		return len(task.receipts) == 0
	}

	// Wait for mined blocks.
	sub := w.mux.Subscribe(core.NewMinedBlockEvent[nist.PublicKey]{})
	defer sub.Unsubscribe()

	// Start mining!
	w.start()

	for i := 0; i < 5; i++ {
		b.txPool.AddLocal(b.newRandomTx(true, false))
		b.txPool.AddLocal(b.newRandomTx(false, false))
		w.postSideBlock(core.ChainSideEvent[nist.PublicKey]{Block: b.newRandomUncle()})
		w.postSideBlock(core.ChainSideEvent[nist.PublicKey]{Block: b.newRandomUncle()})

		select {
		case ev := <-sub.Chan():
			block := ev.Data.(core.NewMinedBlockEvent[nist.PublicKey]).Block
			if _, err := chain.InsertChain([]*types.Block[nist.PublicKey]{block}); err != nil {
				t.Fatalf("failed to insert new mined block %d: %v", block.NumberU64(), err)
			}
		case <-time.After(3 * time.Second): // Worker needs 1s to include new changes.
			t.Fatalf("timeout")
		}
	}
}

func TestEmptyWorkEthash(t *testing.T) {
	testEmptyWork(t, ethashChainConfig,  ethash.NewFaker[nist.PublicKey]())
}
func TestEmptyWorkClique(t *testing.T) {
	testEmptyWork(t, cliqueChainConfig, clique.New[nist.PublicKey](cliqueChainConfig.Clique, rawdb.NewMemoryDatabase()))
}

func testEmptyWork(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine[nist.PublicKey]) {
	defer engine.Close()

	w, _ := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	defer w.close()

	var (
		taskIndex int
		taskCh    = make(chan struct{}, 2)
	)
	checkEqual := func(t *testing.T, task *task[nist.PublicKey], index int) {
		// The first empty work without any txs included
		receiptLen, balance := 0, big.NewInt(0)
		if index == 1 {
			// The second full work with 1 tx included
			receiptLen, balance = 1, big.NewInt(1000)
		}
		if len(task.receipts) != receiptLen {
			t.Fatalf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
		}
		if task.state.GetBalance(testUserAddress).Cmp(balance) != 0 {
			t.Fatalf("account balance mismatch: have %d, want %d", task.state.GetBalance(testUserAddress), balance)
		}
	}
	w.newTaskHook = func(task *task[nist.PublicKey]) {
		if task.block.NumberU64() == 1 {
			checkEqual(t, task, taskIndex)
			taskIndex += 1
			taskCh <- struct{}{}
		}
	}
	w.skipSealHook = func(task *task[nist.PublicKey]) bool { return true }
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	w.start() // Start mining!
	for i := 0; i < 2; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(3 * time.Second).C:
			t.Error("new task timeout")
		}
	}
}

func TestStreamUncleBlock(t *testing.T) {
	ethash :=  ethash.NewFaker[nist.PublicKey]()
	defer ethash.Close()

	w, b := newTestWorker(t, ethashChainConfig, ethash, rawdb.NewMemoryDatabase(), 1)
	defer w.close()

	var taskCh = make(chan struct{})

	taskIndex := 0
	w.newTaskHook = func(task *task[nist.PublicKey]) {
		if task.block.NumberU64() == 2 {
			// The first task is an empty task, the second
			// one has 1 pending tx, the third one has 1 tx
			// and 1 uncle.
			if taskIndex == 2 {
				have := task.block.Header().UncleHash
				want := types.CalcUncleHash([]*types.Header[nist.PublicKey]{b.uncleBlock.Header()})
				if have != want {
					t.Errorf("uncle hash mismatch: have %s, want %s", have.Hex(), want.Hex())
				}
			}
			taskCh <- struct{}{}
			taskIndex += 1
		}
	}
	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	w.start()

	for i := 0; i < 2; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(time.Second).C:
			t.Error("new task timeout")
		}
	}

	w.postSideBlock(core.ChainSideEvent[nist.PublicKey]{Block: b.uncleBlock})

	select {
	case <-taskCh:
	case <-time.NewTimer(time.Second).C:
		t.Error("new task timeout")
	}
}

func TestRegenerateMiningBlockEthash(t *testing.T) {
	testRegenerateMiningBlock(t, ethashChainConfig,  ethash.NewFaker[nist.PublicKey]())
}

func TestRegenerateMiningBlockClique(t *testing.T) {
	testRegenerateMiningBlock(t, cliqueChainConfig, clique.New[nist.PublicKey](cliqueChainConfig.Clique, rawdb.NewMemoryDatabase()))
}

func testRegenerateMiningBlock(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine[nist.PublicKey]) {
	defer engine.Close()

	w, b := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	defer w.close()

	var taskCh = make(chan struct{})

	taskIndex := 0
	w.newTaskHook = func(task *task[nist.PublicKey]) {
		if task.block.NumberU64() == 1 {
			// The first task is an empty task, the second
			// one has 1 pending tx, the third one has 2 txs
			if taskIndex == 2 {
				receiptLen, balance := 2, big.NewInt(2000)
				if len(task.receipts) != receiptLen {
					t.Errorf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
				}
				if task.state.GetBalance(testUserAddress).Cmp(balance) != 0 {
					t.Errorf("account balance mismatch: have %d, want %d", task.state.GetBalance(testUserAddress), balance)
				}
			}
			taskCh <- struct{}{}
			taskIndex += 1
		}
	}
	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}

	w.start()
	// Ignore the first two works
	for i := 0; i < 2; i += 1 {
		select {
		case <-taskCh:
		case <-time.NewTimer(time.Second).C:
			t.Error("new task timeout")
		}
	}
	b.txPool.AddLocals(newTxs)
	time.Sleep(time.Second)

	select {
	case <-taskCh:
	case <-time.NewTimer(time.Second).C:
		t.Error("new task timeout")
	}
}

func TestAdjustIntervalEthash(t *testing.T) {
	testAdjustInterval(t, ethashChainConfig,  ethash.NewFaker[nist.PublicKey]())
}

func TestAdjustIntervalClique(t *testing.T) {
	testAdjustInterval(t, cliqueChainConfig, clique.New[nist.PublicKey](cliqueChainConfig.Clique, rawdb.NewMemoryDatabase()))
}

func testAdjustInterval(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine[nist.PublicKey]) {
	defer engine.Close()

	w, _ := newTestWorker(t, chainConfig, engine, rawdb.NewMemoryDatabase(), 0)
	defer w.close()

	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	var (
		progress = make(chan struct{}, 10)
		result   = make([]float64, 0, 10)
		index    = 0
		start    uint32
	)
	w.resubmitHook = func(minInterval time.Duration, recommitInterval time.Duration) {
		// Short circuit if interval checking hasn't started.
		if atomic.LoadUint32(&start) == 0 {
			return
		}
		var wantMinInterval, wantRecommitInterval time.Duration

		switch index {
		case 0:
			wantMinInterval, wantRecommitInterval = 3*time.Second, 3*time.Second
		case 1:
			origin := float64(3 * time.Second.Nanoseconds())
			estimate := origin*(1-intervalAdjustRatio) + intervalAdjustRatio*(origin/0.8+intervalAdjustBias)
			wantMinInterval, wantRecommitInterval = 3*time.Second, time.Duration(estimate)*time.Nanosecond
		case 2:
			estimate := result[index-1]
			min := float64(3 * time.Second.Nanoseconds())
			estimate = estimate*(1-intervalAdjustRatio) + intervalAdjustRatio*(min-intervalAdjustBias)
			wantMinInterval, wantRecommitInterval = 3*time.Second, time.Duration(estimate)*time.Nanosecond
		case 3:
			wantMinInterval, wantRecommitInterval = time.Second, time.Second
		}

		// Check interval
		if minInterval != wantMinInterval {
			t.Errorf("resubmit min interval mismatch: have %v, want %v ", minInterval, wantMinInterval)
		}
		if recommitInterval != wantRecommitInterval {
			t.Errorf("resubmit interval mismatch: have %v, want %v", recommitInterval, wantRecommitInterval)
		}
		result = append(result, float64(recommitInterval.Nanoseconds()))
		index += 1
		progress <- struct{}{}
	}
	w.start()

	time.Sleep(time.Second) // Ensure two tasks have been summitted due to start opt
	atomic.StoreUint32(&start, 1)

	w.setRecommitInterval(3 * time.Second)
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}

	w.resubmitAdjustCh <- &intervalAdjust{inc: true, ratio: 0.8}
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}

	w.resubmitAdjustCh <- &intervalAdjust{inc: false}
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}

	w.setRecommitInterval(500 * time.Millisecond)
	select {
	case <-progress:
	case <-time.NewTimer(time.Second).C:
		t.Error("interval reset timeout")
	}
}

var PSI1PSM = mps.PrivateStateMetadata{
	ID:          "psi1",
	Name:        "psi1",
	Description: "private state 1",
	Type:        mps.Resident,
	Addresses:   nil,
}

var PSI2PSM = mps.PrivateStateMetadata{
	ID:          "psi2",
	Name:        "psi2",
	Description: "private state 2",
	Type:        mps.Resident,
	Addresses:   nil,
}

func TestPrivatePSMRStateCreated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockptm := private.NewMockPrivateTransactionManager(mockCtrl)
	mockptm.EXPECT().HasFeature(engine.MultiplePrivateStates).Return(true)
	mockptm.EXPECT().Groups().Return([]engine.PrivacyGroup{
		{
			Type:           "RESIDENT",
			Name:           PSI1PSM.Name,
			PrivacyGroupId: base64.StdEncoding.EncodeToString([]byte(PSI1PSM.ID)),
			Description:    "Resident Group 1",
			From:           "",
			Members:        []string{"psi1"},
		},
		{
			Type:           "RESIDENT",
			Name:           PSI2PSM.Name,
			PrivacyGroupId: base64.StdEncoding.EncodeToString([]byte(PSI2PSM.ID)),
			Description:    "Resident Group 2",
			From:           "",
			Members:        []string{"psi2"},
		},
	}, nil)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()
	private.Ptm = mockptm

	db := rawdb.NewMemoryDatabase()
	chainConfig := params.AllCliqueProtocolChanges
	chainConfig.IsQuorum = true
	chainConfig.IsMPS = true
	defer func() { chainConfig.IsQuorum = false }()
	defer func() { chainConfig.IsMPS = false }()

	w, b := newTestWorker(t, chainConfig, clique.New[nist.PublicKey](chainConfig.Clique, db), db, 0)
	defer w.close()

	newBlock := make(chan *types.Block[nist.PublicKey])
	listenNewBlock := func() {
		sub := w.mux.Subscribe(core.NewMinedBlockEvent[nist.PublicKey]{})
		defer sub.Unsubscribe()

		for item := range sub.Chan() {
			newBlock <- item.Data.(core.NewMinedBlockEvent[nist.PublicKey]).Block
		}
	}
	w.start() // Start mining!

	// Ignore first 2 commits caused by start operation
	ignored := make(chan struct{}, 2)
	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		ignored <- struct{}{}
		return true
	}
	timer := time.NewTimer(3 * time.Second)
	for i := 0; i < 2; i++ {
		select {
		case <-ignored:
		case <-timer.C:
			t.Fatalf("timeout")
		}
	}
	timer.Stop()

	go listenNewBlock()

	// Ignore empty commit here for less noise
	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		return len(task.receipts) == 0
	}
	for i := 0; i < 5; i++ {
		randomPrivateTx := b.newRandomTx(true, true)
		mockptm.EXPECT().Receive(common.BytesToEncryptedPayloadHash(randomPrivateTx.Data())).Return("", []string{"psi1", "psi2"}, common.FromHex(testCode), nil, nil).AnyTimes()
		mockptm.EXPECT().Receive(common.EncryptedPayloadHash{}).Return("", []string{}, common.EncryptedPayloadHash{}.Bytes(), nil, nil).AnyTimes()
		expectedContractAddress := crypto.CreateAddress[nist.PublicKey](randomPrivateTx.From(), randomPrivateTx.Nonce())
		b.txPool.AddLocal(randomPrivateTx)
		select {
		case blk := <-newBlock:
			//check if the tx is present
			found := blk.Transaction(randomPrivateTx.Hash())
			if found == nil {
				continue
			}

			latestBlockRoot := b.BlockChain().CurrentBlock().Root()
			_, privDb, err := b.BlockChain().StateAtPSI(latestBlockRoot, types.ToPrivateStateIdentifier("empty"))
			assert.NoError(t, err)
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.Equal(t, privDb.GetCodeSize(expectedContractAddress), 0)
			//contract should exist on both psi states and not be empty
			_, privDb, _ = b.BlockChain().StateAtPSI(latestBlockRoot, types.ToPrivateStateIdentifier("psi1"))
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.False(t, privDb.Empty(expectedContractAddress))
			assert.NotEqual(t, privDb.GetCodeSize(expectedContractAddress), 0)
			_, privDb, _ = b.BlockChain().StateAtPSI(latestBlockRoot, types.ToPrivateStateIdentifier("psi2"))
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.False(t, privDb.Empty(expectedContractAddress))
			assert.NotEqual(t, privDb.GetCodeSize(expectedContractAddress), 0)
			//contract should exist on random state (delegated to emptystate) but no contract code
			_, privDb, _ = b.BlockChain().StateAtPSI(latestBlockRoot, types.ToPrivateStateIdentifier("other"))
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.Equal(t, privDb.GetCodeSize(expectedContractAddress), 0)
		case <-time.NewTimer(3 * time.Second).C: // Worker needs 1s to include new changes.
			t.Fatalf("timeout")
		}
	}

	logsChan := make(chan []*types.Log)
	sub := b.BlockChain().SubscribeLogsEvent(logsChan)
	defer sub.Unsubscribe()

	logsContractData := "6080604052348015600f57600080fd5b507f24ec1d3ff24c2f6ff210738839dbc339cd45a5294d85c79361016243157aae7b60405160405180910390a1603e8060496000396000f3fe6080604052600080fdfea265627a7a72315820937805cb4f2481262ad95d420ab93220f11ceaea518c7ccf119fc2c58f58050d64736f6c63430005110032"
	tx, _ := types.SignTx[nist.PrivateKey,nist.PublicKey](types.NewContractCreation[nist.PublicKey](b.txPool.Nonce(testBankAddress), big.NewInt(0), 470000, nil, common.FromHex(logsContractData)), types.QuorumPrivateTxSigner[nist.PublicKey]{}, testBankKey)

	mockptm.EXPECT().Receive(common.BytesToEncryptedPayloadHash(tx.Data())).Return("", []string{"psi1", "psi2"}, common.FromHex(logsContractData), nil, nil).AnyTimes()

	b.txPool.AddLocal(tx)

	select {
	case logs := <-logsChan:
		assert.Len(t, logs, 2)
		assert.Contains(t, []types.PrivateStateIdentifier{logs[0].PSI, logs[1].PSI}, types.PrivateStateIdentifier("psi1"))
		assert.Contains(t, []types.PrivateStateIdentifier{logs[0].PSI, logs[1].PSI}, types.PrivateStateIdentifier("psi2"))
	case <-time.NewTimer(3 * time.Second).C:
		t.Error("timeout")
	}
}

func TestPrivateLegacyStateCreated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockptm := private.NewMockPrivateTransactionManager(mockCtrl)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()
	private.Ptm = mockptm

	db := rawdb.NewMemoryDatabase()
	chainConfig := params.AllCliqueProtocolChanges
	chainConfig.IsQuorum = true
	chainConfig.IsMPS = false
	defer func() { chainConfig.IsQuorum = false }()

	w, b := newTestWorker(t, chainConfig, clique.New[nist.PublicKey](chainConfig.Clique, db), db, 0)
	defer w.close()

	newBlock := make(chan *types.Block[nist.PublicKey])
	listenNewBlock := func() {
		sub := w.mux.Subscribe(core.NewMinedBlockEvent[nist.PublicKey]{})
		defer sub.Unsubscribe()

		for item := range sub.Chan() {
			newBlock <- item.Data.(core.NewMinedBlockEvent[nist.PublicKey]).Block
		}
	}
	w.start() // Start mining!

	// Ignore first 2 commits caused by start operation
	ignored := make(chan struct{}, 2)
	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		ignored <- struct{}{}
		return true
	}
	for i := 0; i < 2; i++ {
		<-ignored
	}

	go listenNewBlock()

	// Ignore empty commit here for less noise
	w.skipSealHook = func(task *task[nist.PublicKey]) bool {
		return len(task.receipts) == 0
	}
	for i := 0; i < 5; i++ {
		mockptm.EXPECT().Receive(gomock.Any()).Return("", []string{}, common.FromHex(testCode), nil, nil).Times(1)
		randomPrivateTx := b.newRandomTx(true, true)
		expectedContractAddress := crypto.CreateAddress[nist.PublicKey](randomPrivateTx.From(), randomPrivateTx.Nonce())
		b.txPool.AddLocal(randomPrivateTx)
		select {
		case blk := <-newBlock:
			//check if the tx is present
			found := blk.Transaction(randomPrivateTx.Hash())
			if found == nil {
				continue
			}

			latestBlockRoot := b.BlockChain().CurrentBlock().Root()
			//contract exists on default state
			_, privDb, _ := b.BlockChain().StateAtPSI(latestBlockRoot, types.DefaultPrivateStateIdentifier)
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.NotEqual(t, privDb.GetCodeSize(expectedContractAddress), 0)
			//only "private" state on legacy psm
			_, _, err := b.BlockChain().StateAtPSI(latestBlockRoot, types.ToPrivateStateIdentifier("other"))
			assert.Error(t, err, "Only the 'private' psi is supported by the default private state manager")
		case <-time.NewTimer(3 * time.Second).C: // Worker needs 1s to include new changes.
			t.Fatalf("timeout")
		}
	}
}
