package core

import (
	"context"
	"encoding/base64"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/privatecache"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/private/engine"
	"github.com/pavelkrolevets/MIR-pro/rpc"
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
	testKey, _  = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress[nist.PublicKey](*testKey.Public())
)

func buildTestChain(n int, config *params.ChainConfig) ([]*types.Block[nist.PublicKey], map[common.Hash]*types.Block[nist.PublicKey], *BlockChain[nist.PublicKey]) {
	testdb := rawdb.NewMemoryDatabase()
	genesis := GenesisBlockForTesting[nist.PublicKey](testdb, testAddress, big.NewInt(1000000000))
	blocks, _ := GenerateChain[nist.PublicKey](config, genesis,  ethash.NewFaker[nist.PublicKey](), testdb, n, func(i int, block *BlockGen[nist.PublicKey]) {
		block.SetCoinbase(common.Address{0})

		signer := types.QuorumPrivateTxSigner[nist.PublicKey]{}
		tx, err := types.SignTx[nist.PrivateKey, nist.PublicKey](types.NewContractCreation[nist.PublicKey](block.TxNonce(testAddress), big.NewInt(0), testGas, nil, common.FromHex(testCode)), signer, testKey)
		if err != nil {
			panic(err)
		}
		block.AddTx(tx)
	})

	hashes := make([]common.Hash, n+1)
	hashes[len(hashes)-1] = genesis.Hash()
	blockm := make(map[common.Hash]*types.Block[nist.PublicKey], n+1)
	blockm[genesis.Hash()] = genesis
	for i, b := range blocks {
		hashes[len(hashes)-i-2] = b.Hash()
		blockm[b.Hash()] = b
	}

	blockchain, _ := NewBlockChain[nist.PublicKey](testdb, nil, config,  ethash.NewFaker[nist.PublicKey](), vm.Config[nist.PublicKey]{}, nil, nil, nil)
	return blocks, blockm, blockchain
}

func TestMultiplePSMRStateCreated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockptm := private.NewMockPrivateTransactionManager(mockCtrl)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()
	private.Ptm = mockptm

	mockpsm := mps.NewMockPrivateStateManager[nist.PublicKey](mockCtrl)

	mockptm.EXPECT().Receive(gomock.Not(common.EncryptedPayloadHash{})).Return("", []string{"psi1", "psi2"}, common.FromHex(testCode), nil, nil).AnyTimes()
	mockptm.EXPECT().Receive(common.EncryptedPayloadHash{}).Return("", []string{}, common.EncryptedPayloadHash{}.Bytes(), nil, nil).AnyTimes()
	mockptm.EXPECT().HasFeature(engine.MultiplePrivateStates).Return(true)
	mockptm.EXPECT().Groups().Return(PrivacyGroups, nil).AnyTimes()

	mockpsm.EXPECT().ResolveForManagedParty("psi1").Return(&PSI1PSM, nil).AnyTimes()
	mockpsm.EXPECT().ResolveForManagedParty("psi2").Return(&PSI2PSM, nil).AnyTimes()
	mockpsm.EXPECT().PSIs().Return([]types.PrivateStateIdentifier{PSI1PSM.ID, PSI2PSM.ID, types.DefaultPrivateStateIdentifier, types.ToPrivateStateIdentifier("other")}).AnyTimes()

	blocks, blockmap, blockchain := buildTestChain(2, params.QuorumMPSTestChainConfig)
	cache := state.NewDatabase[nist.PublicKey](blockchain.db)
	privateCacheProvider := privatecache.NewPrivateCacheProvider[nist.PublicKey](blockchain.db, nil, cache, false)
	blockchain.privateStateManager = mockpsm

	for _, block := range blocks {
		parent := blockmap[block.ParentHash()]
		statedb, _ := state.New[nist.PublicKey](parent.Root(), blockchain.StateCache(), nil)
		mockpsm.EXPECT().StateRepository(gomock.Any()).Return(mps.NewMultiplePrivateStateRepository[nist.PublicKey](blockchain.db, cache, common.Hash{}, privateCacheProvider)).AnyTimes()

		privateStateRepo, err := blockchain.PrivateStateManager().StateRepository(parent.Root())
		assert.NoError(t, err)

		publicReceipts, privateReceipts, _, _, _ := blockchain.Processor().Process(block, statedb, privateStateRepo, vm.Config[nist.PublicKey]{})

		//managed states tests
		for _, privateReceipt := range privateReceipts {
			expectedContractAddress := privateReceipt.ContractAddress

			emptyState, _ := privateStateRepo.DefaultState()
			assert.True(t, emptyState.Exist(expectedContractAddress))
			assert.Equal(t, emptyState.GetCodeSize(expectedContractAddress), 0)
			ps1, _ := privateStateRepo.StatePSI(types.PrivateStateIdentifier("psi1"))
			assert.True(t, ps1.Exist(expectedContractAddress))
			assert.NotEqual(t, ps1.GetCodeSize(expectedContractAddress), 0)
			ps2, _ := privateStateRepo.StatePSI(types.PrivateStateIdentifier("psi2"))
			assert.True(t, ps2.Exist(expectedContractAddress))
			assert.NotEqual(t, ps2.GetCodeSize(expectedContractAddress), 0)

		}
		//CommitAndWrite to db
		privateStateRepo.CommitAndWrite(false, block)

		//managed states test
		for _, privateReceipt := range privateReceipts {
			expectedContractAddress := privateReceipt.ContractAddress
			latestBlockRoot := block.Root()
			_, privDb, _ := blockchain.StateAtPSI(latestBlockRoot, types.ToPrivateStateIdentifier("empty"))
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.Equal(t, privDb.GetCodeSize(expectedContractAddress), 0)
			//contract exists on both psi states
			_, privDb, _ = blockchain.StateAtPSI(latestBlockRoot, types.PrivateStateIdentifier("psi1"))
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.NotEqual(t, privDb.GetCodeSize(expectedContractAddress), 0)
			_, privDb, _ = blockchain.StateAtPSI(latestBlockRoot, types.PrivateStateIdentifier("psi2"))
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.NotEqual(t, privDb.GetCodeSize(expectedContractAddress), 0)
			//contract should exist on default private state but no contract code
			_, privDb, _ = blockchain.StateAtPSI(latestBlockRoot, types.DefaultPrivateStateIdentifier)
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.Equal(t, privDb.GetCodeSize(expectedContractAddress), 0)
			//contract should exist on random state but no contract code
			_, privDb, _ = blockchain.StateAtPSI(latestBlockRoot, types.ToPrivateStateIdentifier("other"))
			assert.True(t, privDb.Exist(expectedContractAddress))
			assert.Equal(t, privDb.GetCodeSize(expectedContractAddress), 0)
		}

		//mergeReceipts test
		for _, pubReceipt := range publicReceipts {
			assert.Equal(t, 0, len(pubReceipt.PSReceipts))
		}
		for _, privReceipt := range privateReceipts {
			assert.Equal(t, 2, len(privReceipt.PSReceipts))
			assert.NotEqual(t, nil, privReceipt.PSReceipts["psi1"])
			assert.NotEqual(t, nil, privReceipt.PSReceipts["psi2"])
		}

		allReceipts := privateStateRepo.MergeReceipts(publicReceipts, privateReceipts)
		for _, receipt := range allReceipts {
			assert.Equal(t, 3, len(receipt.PSReceipts))
			assert.NotEqual(t, nil, receipt.PSReceipts["empty"])
			assert.NotEqual(t, nil, receipt.PSReceipts["psi1"])
			assert.NotEqual(t, nil, receipt.PSReceipts["psi2"])
		}
	}
}

func TestMPSReset(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockptm := private.NewMockPrivateTransactionManager(mockCtrl)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()
	private.Ptm = mockptm

	mockpsm := mps.NewMockPrivateStateManager[nist.PublicKey](mockCtrl)

	mockptm.EXPECT().Receive(gomock.Not(common.EncryptedPayloadHash{})).Return("", []string{"psi1", "psi2"}, common.FromHex(testCode), nil, nil).AnyTimes()
	mockptm.EXPECT().Receive(common.EncryptedPayloadHash{}).Return("", []string{}, common.EncryptedPayloadHash{}.Bytes(), nil, nil).AnyTimes()
	mockptm.EXPECT().HasFeature(engine.MultiplePrivateStates).Return(true)
	mockptm.EXPECT().Groups().Return(PrivacyGroups, nil).AnyTimes()

	mockpsm.EXPECT().ResolveForManagedParty("psi1").Return(&PSI1PSM, nil).AnyTimes()
	mockpsm.EXPECT().ResolveForManagedParty("psi2").Return(&PSI2PSM, nil).AnyTimes()
	mockpsm.EXPECT().PSIs().Return([]types.PrivateStateIdentifier{PSI1PSM.ID, PSI2PSM.ID}).AnyTimes()

	blocks, blockmap, blockchain := buildTestChain(2, params.QuorumMPSTestChainConfig)
	blockchain.privateStateManager = mockpsm
	cache := state.NewDatabase[nist.PublicKey](blockchain.db)
	privateCacheProvider := privatecache.NewPrivateCacheProvider[nist.PublicKey](blockchain.db, nil, cache, false)

	for _, block := range blocks {
		parent := blockmap[block.ParentHash()]
		statedb, _ := state.New[nist.PublicKey](parent.Root(), blockchain.StateCache(), nil)
		mockpsm.EXPECT().StateRepository(gomock.Any()).Return(mps.NewMultiplePrivateStateRepository[nist.PublicKey](blockchain.db, cache, common.Hash{}, privateCacheProvider)).AnyTimes()

		privateStateRepo, err := blockchain.PrivateStateManager().StateRepository(parent.Root())
		assert.NoError(t, err)

		_, privateReceipts, _, _, _ := blockchain.Processor().Process(block, statedb, privateStateRepo, vm.Config[nist.PublicKey]{})

		for _, privateReceipt := range privateReceipts {
			expectedContractAddress := privateReceipt.ContractAddress

			emptyState, _ := privateStateRepo.DefaultState()
			assert.True(t, emptyState.Exist(expectedContractAddress))
			assert.Equal(t, emptyState.GetCodeSize(expectedContractAddress), 0)
			ps1, _ := privateStateRepo.StatePSI(types.PrivateStateIdentifier("psi1"))
			assert.True(t, ps1.Exist(expectedContractAddress))
			assert.NotEqual(t, ps1.GetCodeSize(expectedContractAddress), 0)
			ps2, _ := privateStateRepo.StatePSI(types.PrivateStateIdentifier("psi2"))
			assert.True(t, ps2.Exist(expectedContractAddress))
			assert.NotEqual(t, ps2.GetCodeSize(expectedContractAddress), 0)

			privateStateRepo.Reset()

			emptyState, _ = privateStateRepo.DefaultState()
			assert.False(t, emptyState.Exist(expectedContractAddress))
			assert.Equal(t, emptyState.GetCodeSize(expectedContractAddress), 0)
			ps1, _ = privateStateRepo.StatePSI(types.PrivateStateIdentifier("psi1"))
			assert.False(t, ps1.Exist(expectedContractAddress))
			assert.Equal(t, ps1.GetCodeSize(expectedContractAddress), 0)
			ps2, _ = privateStateRepo.StatePSI(types.PrivateStateIdentifier("psi2"))
			assert.False(t, ps2.Exist(expectedContractAddress))
			assert.Equal(t, ps2.GetCodeSize(expectedContractAddress), 0)
		}
	}
}

func TestPrivateStateMetadataResolver(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockptm := private.NewMockPrivateTransactionManager(mockCtrl)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()
	private.Ptm = mockptm

	mockptm.EXPECT().Receive(gomock.Not(common.EncryptedPayloadHash{})).Return("", []string{"AAA", "CCC"}, common.FromHex(testCode), nil, nil).AnyTimes()
	mockptm.EXPECT().Receive(common.EncryptedPayloadHash{}).Return("", []string{}, common.EncryptedPayloadHash{}.Bytes(), nil, nil).AnyTimes()
	mockptm.EXPECT().HasFeature(engine.MultiplePrivateStates).Return(true)
	mockptm.EXPECT().Groups().Return(PrivacyGroups, nil).AnyTimes()

	_, _, blockchain := buildTestChain(1, params.QuorumMPSTestChainConfig)

	mpsm := blockchain.privateStateManager

	psm1, _ := mpsm.ResolveForManagedParty("AAA")
	psm2, _ := mpsm.ResolveForManagedParty("CCC")
	_, err := mpsm.ResolveForManagedParty("TEST")
	assert.Equal(t, psm1, privacyGroupToPrivateStateMetadata(PG1))
	assert.Equal(t, psm2, privacyGroupToPrivateStateMetadata(PG2))
	assert.Error(t, err, "unable to find private state metadata for managed party TEST")

	ctx := rpc.WithPrivateStateIdentifier(context.Background(), types.ToPrivateStateIdentifier("RG1"))
	psm1, _ = mpsm.ResolveForUserContext(ctx)
	assert.Equal(t, psm1, privacyGroupToPrivateStateMetadata(PG1))
	ctx = rpc.WithPrivateStateIdentifier(context.Background(), types.ToPrivateStateIdentifier("OTHER"))
	_, err = mpsm.ResolveForUserContext(ctx)
	assert.Error(t, err, "unable to find private state for context psi OTHER")
	_, err = mpsm.ResolveForUserContext(context.Background())
	assert.Error(t, err, "unable to find private state for context psi private")

	assert.Contains(t, mpsm.PSIs(), types.PrivateStateIdentifier("RG1"))
	assert.Contains(t, mpsm.PSIs(), types.PrivateStateIdentifier("RG2"))
	assert.Contains(t, mpsm.PSIs(), types.PrivateStateIdentifier("LEGACY1"))
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

var PG1 = engine.PrivacyGroup{
	Type:           "RESIDENT",
	Name:           "RG1",
	PrivacyGroupId: "RG1",
	Description:    "Resident Group 1",
	From:           "",
	Members:        []string{"AAA", "BBB"},
}

var PG2 = engine.PrivacyGroup{
	Type:           "RESIDENT",
	Name:           "RG2",
	PrivacyGroupId: "RG2",
	Description:    "Resident Group 2",
	From:           "",
	Members:        []string{"CCC", "DDD"},
}

var PrivacyGroups = []engine.PrivacyGroup{
	{
		Type:           "RESIDENT",
		Name:           "RG1",
		PrivacyGroupId: base64.StdEncoding.EncodeToString([]byte("RG1")),
		Description:    "Resident Group 1",
		From:           "",
		Members:        []string{"AAA", "BBB"},
	},
	{
		Type:           "RESIDENT",
		Name:           "RG2",
		PrivacyGroupId: base64.StdEncoding.EncodeToString([]byte("RG2")),
		Description:    "Resident Group 2",
		From:           "",
		Members:        []string{"CCC", "DDD"},
	},
	{
		Type:           "LEGACY",
		Name:           "LEGACY1",
		PrivacyGroupId: "LEGACY1",
		Description:    "Legacy Group 1",
		From:           "",
		Members:        []string{"LEG1", "LEG2"},
	},
}
