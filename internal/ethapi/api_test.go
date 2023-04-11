package ethapi

import (
	"context"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jpmorganchase/quorum-security-plugin-sdk-go/proto"
	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/bloombits"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/multitenancy"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/private/engine"
	"github.com/pavelkrolevets/MIR-pro/private/engine/notinuse"
	"github.com/pavelkrolevets/MIR-pro/rpc"
	"github.com/pavelkrolevets/MIR-pro/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	arbitraryCtx          = context.Background()
	arbitraryPrivateFrom  = "arbitrary private from"
	arbitraryPrivateFor   = []string{"arbitrary party 1", "arbitrary party 2"}
	arbitraryMandatoryFor = []string{"arbitrary party 2"}
	privateTxArgs         = &PrivateTxArgs[nist.PrivateKey, nist.PublicKey]{
		PrivateFrom: arbitraryPrivateFrom,
		PrivateFor:  arbitraryPrivateFor,
	}
	arbitraryFrom         = common.BytesToAddress([]byte("arbitrary address"))
	arbitraryTo           = common.BytesToAddress([]byte("arbitrary address to"))
	arbitraryGas          = uint64(200000)
	arbitraryZeroGasPrice = big.NewInt(0)
	arbitraryZeroValue    = big.NewInt(0)
	arbitraryEmptyData    = new([]byte)
	arbitraryAccessList   = types.AccessList{}
	callTxArgs            = CallArgs{
		From:       &arbitraryFrom,
		To:         &arbitraryTo,
		Gas:        (*hexutil.Uint64)(&arbitraryGas),
		GasPrice:   (*hexutil.Big)(arbitraryZeroGasPrice),
		Value:      (*hexutil.Big)(arbitraryZeroValue),
		Data:       (*hexutil.Bytes)(arbitraryEmptyData),
		AccessList: &arbitraryAccessList,
	}

	arbitrarySimpleStorageContractEncryptedPayloadHash       = common.BytesToEncryptedPayloadHash([]byte("arbitrary payload hash"))
	arbitraryMandatoryRecipientsContractEncryptedPayloadHash = common.BytesToEncryptedPayloadHash([]byte("arbitrary payload hash of tx with mr"))

	simpleStorageContractCreationTx = types.NewContractCreation[nist.PublicKey](
		0,
		big.NewInt(0),
		hexutil.MustDecodeUint64("0x47b760"),
		big.NewInt(0),
		hexutil.MustDecode("0x6060604052341561000f57600080fd5b604051602080610149833981016040528080519060200190919050505b806000819055505b505b610104806100456000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632a1afcd914605157806360fe47b11460775780636d4ce63c146097575b600080fd5b3415605b57600080fd5b606160bd565b6040518082815260200191505060405180910390f35b3415608157600080fd5b6095600480803590602001909190505060c3565b005b341560a157600080fd5b60a760ce565b6040518082815260200191505060405180910390f35b60005481565b806000819055505b50565b6000805490505b905600a165627a7a72305820d5851baab720bba574474de3d09dbeaabc674a15f4dd93b974908476542c23f00029"))

	rawSimpleStorageContractCreationTx = types.NewContractCreation[nist.PublicKey](
		0,
		big.NewInt(0),
		hexutil.MustDecodeUint64("0x47b760"),
		big.NewInt(0),
		arbitrarySimpleStorageContractEncryptedPayloadHash.Bytes())

	arbitrarySimpleStorageContractAddress                    common.Address
	arbitraryStandardPrivateSimpleStorageContractAddress     common.Address
	arbitraryMandatoryRecipientsSimpleStorageContractAddress common.Address

	simpleStorageContractMessageCallTx                   *types.Transaction[nist.PublicKey]
	standardPrivateSimpleStorageContractMessageCallTx    *types.Transaction[nist.PublicKey]
	rawStandardPrivateSimpleStorageContractMessageCallTx *types.Transaction[nist.PublicKey]

	arbitraryCurrentBlockNumber = big.NewInt(1)

	publicStateDB  *state.StateDB[nist.PublicKey]
	privateStateDB *state.StateDB[nist.PublicKey]

	workdir string
)

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func setup() {
	log.Root().SetHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(true)))
	var err error

	memdb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase[nist.PublicKey](memdb)

	publicStateDB, err = state.New[nist.PublicKey](common.Hash{}, db, nil)
	if err != nil {
		panic(err)
	}
	privateStateDB, err = state.New[nist.PublicKey](common.Hash{}, db, nil)
	if err != nil {
		panic(err)
	}

	private.Ptm = &StubPrivateTransactionManager{}

	key, _ := crypto.GenerateKey[nist.PrivateKey]()
	from := crypto.PubkeyToAddress[nist.PublicKey](*key.Public())

	arbitrarySimpleStorageContractAddress = crypto.CreateAddress[nist.PublicKey](from, 0)

	simpleStorageContractMessageCallTx = types.NewTransaction[nist.PublicKey](
		0,
		arbitrarySimpleStorageContractAddress,
		big.NewInt(0),
		hexutil.MustDecodeUint64("0x47b760"),
		big.NewInt(0),
		hexutil.MustDecode("0x60fe47b1000000000000000000000000000000000000000000000000000000000000000d"))

	arbitraryStandardPrivateSimpleStorageContractAddress = crypto.CreateAddress[nist.PublicKey](from, 1)

	standardPrivateSimpleStorageContractMessageCallTx = types.NewTransaction[nist.PublicKey](
		0,
		arbitraryStandardPrivateSimpleStorageContractAddress,
		big.NewInt(0),
		hexutil.MustDecodeUint64("0x47b760"),
		big.NewInt(0),
		hexutil.MustDecode("0x60fe47b1000000000000000000000000000000000000000000000000000000000000000e"))

	rawStandardPrivateSimpleStorageContractMessageCallTx = types.NewTransaction[nist.PublicKey](
		0,
		arbitraryStandardPrivateSimpleStorageContractAddress,
		big.NewInt(0),
		hexutil.MustDecodeUint64("0x47b760"),
		big.NewInt(0),
		arbitrarySimpleStorageContractEncryptedPayloadHash.Bytes())

	workdir, err = ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
}

func teardown() {
	log.Root().SetHandler(log.DiscardHandler())
	os.RemoveAll(workdir)
}

func TestDoEstimateGas_whenNoValueTx_Pre_Istanbul(t *testing.T) {
	assert := assert.New(t)

	estimation, err := DoEstimateGas[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{CurrentHeadNumber: big.NewInt(10)}, callTxArgs, rpc.BlockNumberOrHashWithNumber(10), math.MaxInt64)

	assert.NoError(err, "gas estimation")
	assert.Equal(hexutil.Uint64(25352), estimation, "estimation for a public or private tx")
}

func TestDoEstimateGas_whenNoValueTx_Istanbul(t *testing.T) {
	assert := assert.New(t)

	estimation, err := DoEstimateGas[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{IstanbulBlock: big.NewInt(0), CurrentHeadNumber: big.NewInt(10)}, callTxArgs, rpc.BlockNumberOrHashWithNumber(10), math.MaxInt64)

	assert.NoError(err, "gas estimation")
	assert.Equal(hexutil.Uint64(22024), estimation, "estimation for a public or private tx")
}

func TestSimulateExecution_whenStandardPrivateCreation(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractCreationTx, privateTxArgs)

	assert.NoError(err, "simulate execution")
	assert.Empty(affectedCACreationTxHashes, "creation tx should not have any affected contract creation tx hashes")
	assert.Equal(common.Hash{}, merkleRoot, "no private state validation")
}

func TestSimulateExecution_whenPartyProtectionCreation(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection

	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractCreationTx, privateTxArgs)

	assert.NoError(err, "simulation execution")
	assert.Empty(affectedCACreationTxHashes, "creation tx should not have any affected contract creation tx hashes")
	assert.Equal(common.Hash{}, merkleRoot, "no private state validation")
}

func TestSimulateExecution_whenCreationWithStateValidation(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStateValidation

	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractCreationTx, privateTxArgs)

	assert.NoError(err, "simulate execution")
	assert.Empty(affectedCACreationTxHashes, "creation tx should not have any affected contract creation tx hashes")
	assert.NotEqual(common.Hash{}, merkleRoot, "private state validation")
}

func TestSimulateExecution_whenStandardPrivateMessageCall(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	privateStateDB.SetCode(arbitraryStandardPrivateSimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000002"))
	privateStateDB.SetState(arbitraryStandardPrivateSimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, standardPrivateSimpleStorageContractMessageCallTx, privateTxArgs)

	log.Debug("state", "state", privateStateDB.GetState(arbitraryStandardPrivateSimpleStorageContractAddress, common.Hash{0}))

	assert.NoError(err, "simulate execution")
	assert.Empty(affectedCACreationTxHashes, "standard private contract should not have any affected contract creation tx hashes")
	assert.Equal(common.Hash{}, merkleRoot, "no private state validation")
}

func TestSimulateExecution_StandardPrivateMessageCallSucceedsWheContractNotAvailableLocally(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, standardPrivateSimpleStorageContractMessageCallTx, privateTxArgs)

	log.Debug("state", "state", privateStateDB.GetState(arbitraryStandardPrivateSimpleStorageContractAddress, common.Hash{0}))

	assert.NoError(err, "simulate execution")
	assert.Empty(affectedCACreationTxHashes, "standard private contract should not have any affected contract creation tx hashes")
	assert.Equal(common.Hash{}, merkleRoot, "no private state validation")
}

func TestSimulateExecution_whenPartyProtectionMessageCall(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection

	privateStateDB.SetCode(arbitrarySimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000001"))
	privateStateDB.SetPrivacyMetadata(arbitrarySimpleStorageContractAddress, &state.PrivacyMetadata{
		PrivacyFlag:    privateTxArgs.PrivacyFlag,
		CreationTxHash: arbitrarySimpleStorageContractEncryptedPayloadHash,
	})

	privateStateDB.SetState(arbitrarySimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractMessageCallTx, privateTxArgs)

	expectedCACreationTxHashes := []common.EncryptedPayloadHash{arbitrarySimpleStorageContractEncryptedPayloadHash}

	log.Debug("state", "state", privateStateDB.GetState(arbitrarySimpleStorageContractAddress, common.Hash{0}))

	assert.NoError(err, "simulate execution")
	assert.NotEmpty(affectedCACreationTxHashes, "affected contract accounts' creation transacton hashes")
	assert.Equal(common.Hash{}, merkleRoot, "no private state validation")
	assert.True(len(affectedCACreationTxHashes) == len(expectedCACreationTxHashes))
}

func TestSimulateExecution_whenPartyProtectionMessageCallAndPrivacyEnhancementsDisabled(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection

	params.QuorumTestChainConfig.PrivacyEnhancementsBlock = nil
	defer func() { params.QuorumTestChainConfig.PrivacyEnhancementsBlock = big.NewInt(0) }()

	stbBackend := &StubBackend[nist.PrivateKey,nist.PublicKey]{}
	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, stbBackend, arbitraryFrom, simpleStorageContractMessageCallTx, privateTxArgs)

	// the simulation returns early without executing the transaction
	assert.False(stbBackend.getEVMCalled, "simulation is ended early - before getEVM is called")
	assert.NoError(err, "simulate execution")
	assert.Empty(affectedCACreationTxHashes, "affected contract accounts' creation transacton hashes")
	assert.Equal(common.Hash{}, merkleRoot, "no private state validation")
}

func TestSimulateExecution_whenStateValidationMessageCall(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStateValidation

	privateStateDB.SetCode(arbitrarySimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000001"))
	privateStateDB.SetPrivacyMetadata(arbitrarySimpleStorageContractAddress, &state.PrivacyMetadata{
		PrivacyFlag:    privateTxArgs.PrivacyFlag,
		CreationTxHash: arbitrarySimpleStorageContractEncryptedPayloadHash,
	})

	privateStateDB.SetState(arbitrarySimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	affectedCACreationTxHashes, merkleRoot, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractMessageCallTx, privateTxArgs)

	expectedCACreationTxHashes := []common.EncryptedPayloadHash{arbitrarySimpleStorageContractEncryptedPayloadHash}

	log.Debug("state", "state", privateStateDB.GetState(arbitrarySimpleStorageContractAddress, common.Hash{0}))

	assert.NoError(err, "simulate execution")
	assert.NotEmpty(affectedCACreationTxHashes, "affected contract accounts' creation transacton hashes")
	assert.NotEqual(common.Hash{}, merkleRoot, "private state validation")
	assert.True(len(affectedCACreationTxHashes) == len(expectedCACreationTxHashes))
}

//mix and match flags
func TestSimulateExecution_PrivacyFlagPartyProtectionCallingStandardPrivateContract_Error(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection

	privateStateDB.SetCode(arbitraryStandardPrivateSimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000002"))
	privateStateDB.SetState(arbitraryStandardPrivateSimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	_, _, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, standardPrivateSimpleStorageContractMessageCallTx, privateTxArgs)

	log.Debug("state", "state", privateStateDB.GetState(arbitraryStandardPrivateSimpleStorageContractAddress, common.Hash{0}))

	assert.Error(err, "simulate execution")
}

func TestSimulateExecution_StandardPrivateFlagCallingPartyProtectionContract_Error(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	privateStateDB.SetCode(arbitrarySimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000001"))
	privateStateDB.SetPrivacyMetadata(arbitrarySimpleStorageContractAddress, &state.PrivacyMetadata{
		PrivacyFlag:    engine.PrivacyFlagPartyProtection,
		CreationTxHash: arbitrarySimpleStorageContractEncryptedPayloadHash,
	})

	privateStateDB.SetState(arbitrarySimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	_, _, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractMessageCallTx, privateTxArgs)

	assert.Error(err, "simulate execution")
}

func TestSimulateExecution_StandardPrivateFlagCallingStateValidationContract_Error(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	privateStateDB.SetCode(arbitrarySimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000001"))
	privateStateDB.SetPrivacyMetadata(arbitrarySimpleStorageContractAddress, &state.PrivacyMetadata{
		PrivacyFlag:    engine.PrivacyFlagStateValidation,
		CreationTxHash: arbitrarySimpleStorageContractEncryptedPayloadHash,
	})

	privateStateDB.SetState(arbitrarySimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	_, _, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractMessageCallTx, privateTxArgs)

	log.Debug("state", "state", privateStateDB.GetState(arbitrarySimpleStorageContractAddress, common.Hash{0}))

	assert.Error(err, "simulate execution")
}

func TestSimulateExecution_PartyProtectionFlagCallingStateValidationContract_Error(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection

	privateStateDB.SetCode(arbitrarySimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000001"))
	privateStateDB.SetPrivacyMetadata(arbitrarySimpleStorageContractAddress, &state.PrivacyMetadata{
		PrivacyFlag:    engine.PrivacyFlagStateValidation,
		CreationTxHash: arbitrarySimpleStorageContractEncryptedPayloadHash,
	})

	privateStateDB.SetState(arbitrarySimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	_, _, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractMessageCallTx, privateTxArgs)

	log.Debug("state", "state", privateStateDB.GetState(arbitrarySimpleStorageContractAddress, common.Hash{0}))

	assert.Error(err, "simulate execution")
}

func TestSimulateExecution_StateValidationFlagCallingPartyProtectionContract_Error(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStateValidation

	privateStateDB.SetCode(arbitrarySimpleStorageContractAddress, hexutil.MustDecode("0x608060405234801561001057600080fd5b506040516020806101618339810180604052602081101561003057600080fd5b81019080805190602001909291905050508060008190555050610109806100586000396000f3fe6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146099575b600080fd5b348015605957600080fd5b50608360048036036020811015606e57600080fd5b810190808035906020019092919050505060c1565b6040518082815260200191505060405180910390f35b34801560a457600080fd5b5060ab60d4565b6040518082815260200191505060405180910390f35b6000816000819055506000549050919050565b6000805490509056fea165627a7a723058203624ca2e3479d3fa5a12d97cf3dae0d9a6de3a3b8a53c8605b9cd398d9766b9f00290000000000000000000000000000000000000000000000000000000000000001"))
	privateStateDB.SetPrivacyMetadata(arbitrarySimpleStorageContractAddress, &state.PrivacyMetadata{
		PrivacyFlag:    engine.PrivacyFlagPartyProtection,
		CreationTxHash: arbitrarySimpleStorageContractEncryptedPayloadHash,
	})

	privateStateDB.SetState(arbitrarySimpleStorageContractAddress, common.Hash{0}, common.Hash{100})
	privateStateDB.Commit(true)

	_, _, err := simulateExecutionForPE[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, arbitraryFrom, simpleStorageContractMessageCallTx, privateTxArgs)

	//expectedCACreationTxHashes := []common.EncryptedPayloadHash{arbitrarySimpleStorageContractEncryptedPayloadHash}

	log.Debug("state", "state", privateStateDB.GetState(arbitrarySimpleStorageContractAddress, common.Hash{0}))

	assert.Error(err, "simulate execution")
}

func TestHandlePrivateTransaction_whenInvalidFlag(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = 4

	_, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.Error(err, "invalid privacyFlag")
}

func TestHandlePrivateTransaction_whenPrivateFromDoesNotMatchPrivateState(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockpsm := mps.NewMockPrivateStateManager[nist.PublicKey](mockCtrl)
	mockpsm.EXPECT().ResolveForUserContext(gomock.Any()).Return(mps.NewPrivateStateMetadata("PS1", "PS1", "", mps.Resident, []string{"some address"}), nil).AnyTimes()

	_, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &MPSStubBackend[nist.PrivateKey,nist.PublicKey]{psmr: mockpsm}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.Error(err, "The PrivateFrom (arbitrary private from) address does not match the specified private state (PS1) ")
}

func TestHandlePrivateTransaction_whenPrivateFromMatchesPrivateState(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockpsm := mps.NewMockPrivateStateManager[nist.PublicKey](mockCtrl)
	mockpsm.EXPECT().ResolveForUserContext(gomock.Any()).Return(mps.NewPrivateStateMetadata("PS1", "PS1", "", mps.Resident, []string{"some address"}), nil).AnyTimes()

	// empty data field means that checkAndHandlePrivateTransaction exits without doing handlePrivateTransaction
	emptyTx := types.NewContractCreation[nist.PublicKey](
		0,
		big.NewInt(0),
		hexutil.MustDecodeUint64("0x47b760"),
		big.NewInt(0),
		nil)

	mpsTxArgs := &PrivateTxArgs[nist.PrivateKey,nist.PublicKey]{
		PrivateFrom: "some address",
		PrivateFor:  []string{"arbitrary party 1", "arbitrary party 2"},
	}
	_, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &MPSStubBackend[nist.PrivateKey,nist.PublicKey]{psmr: mockpsm}, emptyTx, mpsTxArgs, arbitraryFrom, NormalTransaction)

	assert.Nil(err)
}

func TestHandlePrivateTransaction_withPartyProtectionTxAndPrivacyEnhancementsIsDisabled(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = 1
	params.QuorumTestChainConfig.PrivacyEnhancementsBlock = nil
	defer func() { params.QuorumTestChainConfig.PrivacyEnhancementsBlock = big.NewInt(0) }()

	_, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.Error(err, "PrivacyEnhancements are disabled. Can only accept transactions with PrivacyFlag=0(StandardPrivate).")
}

func TestHandlePrivateTransaction_whenStandardPrivateCreation(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	isPrivate, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	if err != nil {
		t.Fatalf("%s", err)
	}

	assert.True(isPrivate, "must be a private transaction")
}

func TestHandlePrivateTransaction_whenStandardPrivateCallingContractThatIsNotAvailableLocally(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	isPrivate, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, standardPrivateSimpleStorageContractMessageCallTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.NoError(err, "no error expected")

	assert.True(isPrivate, "must be a private transaction")
}

func TestHandlePrivateTransaction_whenPartyProtectionCallingContractThatIsNotAvailableLocally(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection

	isPrivate, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, standardPrivateSimpleStorageContractMessageCallTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.Error(err, "handle invalid message call")

	assert.True(isPrivate, "must be a private transaction")
}

func TestHandlePrivateTransaction_whenPartyProtectionCallingStandardPrivate(t *testing.T) {
	assert := assert.New(t)
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection

	isPrivate, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, standardPrivateSimpleStorageContractMessageCallTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.Error(err, "handle invalid message call")

	assert.True(isPrivate, "must be a private transaction")
}

func TestHandlePrivateTransaction_whenRawStandardPrivateCreation(t *testing.T) {
	assert := assert.New(t)
	private.Ptm = &StubPrivateTransactionManager{creation: true}
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	isPrivate, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, rawSimpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, RawTransaction)

	assert.NoError(err, "raw standard private creation succeeded")
	assert.True(isPrivate, "must be a private transaction")
}

func TestHandlePrivateTransaction_whenRawStandardPrivateMessageCall(t *testing.T) {
	assert := assert.New(t)
	private.Ptm = &StubPrivateTransactionManager{creation: false}
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate

	_, err := handlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, rawStandardPrivateSimpleStorageContractMessageCallTx, privateTxArgs, arbitraryFrom, RawTransaction)

	assert.NoError(err, "raw standard private msg call succeeded")

}

func TestHandlePrivateTransaction_whenMandatoryRecipients(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTM := private.NewMockPrivateTransactionManager(mockCtrl)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
		privateTxArgs.MandatoryRecipients = nil
	}()
	private.Ptm = mockTM
	privateTxArgs.MandatoryRecipients = arbitraryMandatoryFor
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagMandatoryRecipients

	var capturedMetadata engine.ExtraMetadata

	mockTM.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(arg1 interface{}, arg2 string, arg3 interface{}, arg4 *engine.ExtraMetadata) {
			capturedMetadata = *arg4
		}).Times(1)

	_, err := handlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.NoError(err)
	assert.Equal(engine.PrivacyFlagMandatoryRecipients, capturedMetadata.PrivacyFlag)
	assert.Equal(arbitraryMandatoryFor, capturedMetadata.MandatoryRecipients)

}

func TestHandlePrivateTransaction_whenRawPrivateWithMandatoryRecipients(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTM := private.NewMockPrivateTransactionManager(mockCtrl)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
		privateTxArgs.MandatoryRecipients = nil
	}()
	private.Ptm = mockTM
	privateTxArgs.MandatoryRecipients = arbitraryMandatoryFor

	privateTxArgs.PrivacyFlag = engine.PrivacyFlagMandatoryRecipients

	var capturedMetadata engine.ExtraMetadata

	mockTM.EXPECT().ReceiveRaw(gomock.Any()).Times(1)

	mockTM.EXPECT().SendSignedTx(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(arg1 interface{}, arg2 []string, arg3 *engine.ExtraMetadata) {
			capturedMetadata = *arg3
		}).Times(1)

	_, err := handlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, RawTransaction)

	assert.NoError(err)
	assert.Equal(engine.PrivacyFlagMandatoryRecipients, capturedMetadata.PrivacyFlag)
	assert.Equal(arbitraryMandatoryFor, capturedMetadata.MandatoryRecipients)

}

func TestHandlePrivateTransaction_whenMandatoryRecipientsDataInvalid(t *testing.T) {
	assert := assert.New(t)

	privateTxArgs.PrivacyFlag = engine.PrivacyFlagMandatoryRecipients

	_, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.Error(err, "missing mandatory recipients data. if no mandatory recipients required consider using PrivacyFlag=1(PartyProtection)")

}

func TestHandlePrivateTransaction_whenNoMandatoryRecipientsData(t *testing.T) {
	assert := assert.New(t)

	privateTxArgs.PrivacyFlag = engine.PrivacyFlagPartyProtection
	defer func() {
		privateTxArgs.MandatoryRecipients = nil
	}()
	privateTxArgs.MandatoryRecipients = arbitraryMandatoryFor

	_, _, _, err := checkAndHandlePrivateTransaction[nist.PrivateKey,nist.PublicKey](arbitraryCtx, &StubBackend[nist.PrivateKey,nist.PublicKey]{}, simpleStorageContractCreationTx, privateTxArgs, arbitraryFrom, NormalTransaction)

	assert.Error(err, "privacy metadata invalid. mandatory recipients are only applicable for PrivacyFlag=2(MandatoryRecipients)")

}

func TestGetContractPrivacyMetadata(t *testing.T) {
	assert := assert.New(t)

	keystore, _, _ := createKeystore[nist.PrivateKey,nist.PublicKey](t)

	stbBackend := &StubBackend[nist.PrivateKey,nist.PublicKey]{}
	stbBackend.multitenancySupported = false
	stbBackend.isPrivacyMarkerTransactionCreationEnabled = false
	stbBackend.ks = keystore
	stbBackend.accountManager = accounts.NewManager[nist.PublicKey](&accounts.Config{InsecureUnlockAllowed: true}, stbBackend)
	stbBackend.poolNonce = 999

	public := NewPublicTransactionPoolAPI[nist.PrivateKey,nist.PublicKey](stbBackend, nil)

	privacyMetadata, _ := public.GetContractPrivacyMetadata(arbitraryCtx, arbitrarySimpleStorageContractAddress)

	assert.Equal(engine.PrivacyFlagPartyProtection, privacyMetadata.PrivacyFlag)
	assert.Equal(arbitrarySimpleStorageContractEncryptedPayloadHash, privacyMetadata.CreationTxHash)
	assert.Equal(0, len(privacyMetadata.MandatoryRecipients))
}

func TestGetContractPrivacyMetadataWhenMandatoryRecipients(t *testing.T) {
	assert := assert.New(t)

	keystore, _, _ := createKeystore[nist.PrivateKey,nist.PublicKey](t)

	stbBackend := &StubBackend[nist.PrivateKey,nist.PublicKey]{}
	stbBackend.multitenancySupported = false
	stbBackend.isPrivacyMarkerTransactionCreationEnabled = false
	stbBackend.ks = keystore
	stbBackend.accountManager = accounts.NewManager[nist.PublicKey](&accounts.Config{InsecureUnlockAllowed: true}, stbBackend)
	stbBackend.poolNonce = 999

	public := NewPublicTransactionPoolAPI[nist.PrivateKey,nist.PublicKey](stbBackend, nil)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTM := private.NewMockPrivateTransactionManager(mockCtrl)

	saved := private.Ptm
	defer func() {
		private.Ptm = saved
	}()
	private.Ptm = mockTM

	var capturedTxHash common.EncryptedPayloadHash

	mockTM.EXPECT().GetMandatory(gomock.Any()).
		DoAndReturn(func(arg1 common.EncryptedPayloadHash) ([]string, error) {
			capturedTxHash = arg1
			return arbitraryMandatoryFor, nil
		}).Times(1)

	privacyMetadata, _ := public.GetContractPrivacyMetadata(arbitraryCtx, arbitraryMandatoryRecipientsSimpleStorageContractAddress)

	assert.Equal(arbitraryMandatoryRecipientsContractEncryptedPayloadHash, capturedTxHash)

	assert.Equal(engine.PrivacyFlagMandatoryRecipients, privacyMetadata.PrivacyFlag)
	assert.Equal(arbitraryMandatoryRecipientsContractEncryptedPayloadHash, privacyMetadata.CreationTxHash)
	assert.Equal(arbitraryMandatoryFor, privacyMetadata.MandatoryRecipients)
}

func TestSubmitPrivateTransaction(t *testing.T) {
	assert := assert.New(t)

	keystore, fromAcct, toAcct := createKeystore[nist.PrivateKey,nist.PublicKey](t)

	stbBackend := &StubBackend[nist.PrivateKey,nist.PublicKey]{}
	stbBackend.multitenancySupported = false
	stbBackend.isPrivacyMarkerTransactionCreationEnabled = false
	stbBackend.ks = keystore
	stbBackend.accountManager = accounts.NewManager[nist.PublicKey](&accounts.Config{InsecureUnlockAllowed: true}, stbBackend)
	stbBackend.poolNonce = 999

	privateAccountAPI := NewPrivateAccountAPI[nist.PrivateKey,nist.PublicKey](stbBackend, nil)

	gas := hexutil.Uint64(999999)
	nonce := hexutil.Uint64(123)
	payload := hexutil.Bytes(([]byte("0x43d3e767000000000000000000000000000000000000000000000000000000000000000a"))[:])
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate
	txArgs := SendTxArgs[nist.PrivateKey,nist.PublicKey]{PrivateTxArgs: *privateTxArgs, From: fromAcct.Address, To: &toAcct.Address, Gas: &gas, Nonce: &nonce, Data: &payload}

	_, err := privateAccountAPI.SendTransaction(arbitraryCtx, txArgs, "")

	assert.NoError(err)
	assert.True(stbBackend.sendTxCalled, "transaction was not sent")
	assert.True(stbBackend.txThatWasSent.IsPrivate(), "must be a private transaction")
	assert.Equal(fromAcct.Address, stbBackend.txThatWasSent.From(), "incorrect 'From' address on transaction")
	assert.Equal(toAcct.Address, *stbBackend.txThatWasSent.To(), "incorrect 'To' address on transaction")
	assert.Equal(uint64(123), stbBackend.txThatWasSent.Nonce(), "incorrect nonce on transaction")
}

func TestSubmitPrivateTransactionWithPrivacyMarkerEnabled(t *testing.T) {
	assert := assert.New(t)

	keystore, fromAcct, toAcct := createKeystore[nist.PrivateKey,nist.PublicKey](t)

	params.QuorumTestChainConfig.PrivacyPrecompileBlock = big.NewInt(0)
	defer func() { params.QuorumTestChainConfig.PrivacyPrecompileBlock = nil }()

	stbBackend := &StubBackend[nist.PrivateKey,nist.PublicKey]{}
	stbBackend.multitenancySupported = false
	stbBackend.isPrivacyMarkerTransactionCreationEnabled = true
	stbBackend.ks = keystore
	stbBackend.accountManager = accounts.NewManager[nist.PublicKey](&accounts.Config{InsecureUnlockAllowed: true}, stbBackend)

	privateAccountAPI := NewPrivateAccountAPI[nist.PrivateKey,nist.PublicKey](stbBackend, nil)

	gas := hexutil.Uint64(999999)
	nonce := hexutil.Uint64(123)
	payload := hexutil.Bytes(([]byte("0x43d3e767000000000000000000000000000000000000000000000000000000000000000a"))[:])
	privateTxArgs.PrivacyFlag = engine.PrivacyFlagStandardPrivate
	txArgs := SendTxArgs[nist.PrivateKey,nist.PublicKey]{PrivateTxArgs: *privateTxArgs, From: fromAcct.Address, To: &toAcct.Address, Gas: &gas, Nonce: &nonce, Data: &payload}

	_, err := privateAccountAPI.SendTransaction(arbitraryCtx, txArgs, "")

	assert.NoError(err)
	assert.True(stbBackend.sendTxCalled, "transaction was not sent")
	assert.False(stbBackend.txThatWasSent.IsPrivate(), "transaction was private, instead of privacy marker transaction (public)")
	assert.Equal(fromAcct.Address, stbBackend.txThatWasSent.From(), "expected privacy marker transaction to have same 'from' address as internal private tx")
	assert.Equal(common.QuorumPrivacyPrecompileContractAddress(), *stbBackend.txThatWasSent.To(), "transaction 'To' address should be privacy marker precompile")
	assert.Equal(uint64(nonce), stbBackend.txThatWasSent.Nonce(), "incorrect nonce on transaction")
	assert.NotEqual(hexutil.Uint64(stbBackend.txThatWasSent.Gas()), gas, "privacy marker transaction should not have same gas value as internal private tx")
}

func TestSetRawTransactionPrivateFrom(t *testing.T) {
	somePTMAddr := "some-ptm-addr"
	psiID := types.PrivateStateIdentifier("myPSI")
	mpsPTMAddrs := []string{somePTMAddr}

	tests := []struct {
		name                  string
		receiveRawPrivateFrom string
		argsPrivateFrom       string
		wantPrivateFrom       string
	}{
		{
			name:                  "receiveRawPrivateFromIfNoArgPrivateFrom",
			receiveRawPrivateFrom: somePTMAddr,
			argsPrivateFrom:       "",
			wantPrivateFrom:       somePTMAddr,
		},
		{
			name:                  "argPrivateFromOnly",
			receiveRawPrivateFrom: "",
			argsPrivateFrom:       somePTMAddr,
			wantPrivateFrom:       somePTMAddr,
		},
		{
			name:                  "equalArgAndReceiveRawPrivateFrom",
			receiveRawPrivateFrom: somePTMAddr,
			argsPrivateFrom:       somePTMAddr,
			wantPrivateFrom:       somePTMAddr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			savedPTM := private.Ptm
			defer func() { private.Ptm = savedPTM }()

			mockPTM := private.NewMockPrivateTransactionManager(ctrl)
			mockPTM.EXPECT().ReceiveRaw(gomock.Any()).Return(nil, somePTMAddr, nil, nil).Times(1)
			private.Ptm = mockPTM

			psm := mps.NewPrivateStateMetadata(psiID, "", "", 0, mpsPTMAddrs)

			mockPSMR := mps.NewMockPrivateStateMetadataResolver(ctrl)
			mockPSMR.EXPECT().ResolveForUserContext(gomock.Any()).Return(psm, nil).Times(1)

			b := &MPSStubBackend[nist.PrivateKey,nist.PublicKey]{
				psmr: mockPSMR,
			}

			tx := types.NewTransaction[nist.PublicKey](0, common.Address{}, nil, 0, nil, []byte("ptm-hash"))

			args := &PrivateTxArgs[nist.PrivateKey,nist.PublicKey]{
				PrivateFor:  []string{"some-ptm-recipient"},
				PrivateFrom: tt.argsPrivateFrom,
			}

			err := args.SetRawTransactionPrivateFrom(context.Background(), b, tx)

			require.NoError(t, err)
			require.Equal(t, tt.wantPrivateFrom, args.PrivateFrom)
		})
	}
}

func TestSetRawTransactionPrivateFrom_DifferentArgPrivateFromAndReceiveRawPrivateFrom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	savedPTM := private.Ptm
	defer func() { private.Ptm = savedPTM }()

	receiveRawPrivateFrom := "some-ptm-addr"
	argsPrivateFrom := "other-ptm-addr"

	mockPTM := private.NewMockPrivateTransactionManager(ctrl)
	mockPTM.EXPECT().ReceiveRaw(gomock.Any()).Return(nil, receiveRawPrivateFrom, nil, nil).Times(1)
	private.Ptm = mockPTM

	b := &MPSStubBackend[nist.PrivateKey,nist.PublicKey]{}

	tx := types.NewTransaction[nist.PublicKey](0, common.Address{}, nil, 0, nil, []byte("ptm-hash"))

	args := &PrivateTxArgs[nist.PrivateKey,nist.PublicKey]{
		PrivateFor:  []string{"some-ptm-recipient"},
		PrivateFrom: argsPrivateFrom,
	}

	err := args.SetRawTransactionPrivateFrom(context.Background(), b, tx)

	require.EqualError(t, err, "The PrivateFrom address retrieved from the privacy manager does not match private PrivateFrom (other-ptm-addr) specified in transaction arguments.")
}

func TestSetRawTransactionPrivateFrom_ResolvePrivateFromIsNotMPSTenantAddr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	savedPTM := private.Ptm
	defer func() { private.Ptm = savedPTM }()

	receiveRawPrivateFrom := "some-ptm-addr"
	psiID := types.PrivateStateIdentifier("myPSI")

	mpsPTMAddrs := []string{"other-ptm-addr"}

	mockPTM := private.NewMockPrivateTransactionManager(ctrl)
	mockPTM.EXPECT().ReceiveRaw(gomock.Any()).Return(nil, receiveRawPrivateFrom, nil, nil).Times(1)
	private.Ptm = mockPTM

	psm := mps.NewPrivateStateMetadata(psiID, "", "", 0, mpsPTMAddrs)

	mockPSMR := mps.NewMockPrivateStateMetadataResolver(ctrl)
	mockPSMR.EXPECT().ResolveForUserContext(gomock.Any()).Return(psm, nil).Times(1)

	b := &MPSStubBackend[nist.PrivateKey,nist.PublicKey]{
		psmr: mockPSMR,
	}

	tx := types.NewTransaction[nist.PublicKey](0, common.Address{}, nil, 0, nil, []byte("ptm-hash"))

	args := &PrivateTxArgs[nist.PrivateKey,nist.PublicKey]{
		PrivateFor: []string{"some-ptm-recipient"},
	}

	err := args.SetRawTransactionPrivateFrom(context.Background(), b, tx)

	require.EqualError(t, err, "The PrivateFrom address does not match the specified private state (myPSI)")
}

func createKeystore[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T) (*keystore.KeyStore[T,P], accounts.Account, accounts.Account) {
	assert := assert.New(t)

	keystore := keystore.NewKeyStore[T,P](filepath.Join(workdir, "keystore"), keystore.StandardScryptN, keystore.StandardScryptP)
	fromAcct, err := keystore.NewAccount("")
	assert.NoError(err)
	toAcct, err := keystore.NewAccount("")
	assert.NoError(err)

	return keystore, fromAcct, toAcct
}

type StubBackend [T crypto.PrivateKey,P crypto.PublicKey]struct {
	getEVMCalled                              bool
	sendTxCalled                              bool
	txThatWasSent                             *types.Transaction[P]
	mockAccountExtraDataStateGetter           *vm.MockAccountExtraDataStateGetter
	multitenancySupported                     bool
	isPrivacyMarkerTransactionCreationEnabled bool
	accountManager                            *accounts.Manager[P]
	ks                                        *keystore.KeyStore[T,P]
	poolNonce                                 uint64
	allowUnprotectedTxs                       bool

	IstanbulBlock     *big.Int
	CurrentHeadNumber *big.Int
}

var _ Backend[nist.PrivateKey, nist.PublicKey] = &StubBackend[nist.PrivateKey, nist.PublicKey]{}

func (sb *StubBackend[T,P]) UnprotectedAllowed() bool {
	return sb.allowUnprotectedTxs
}

func (sb *StubBackend[T,P]) CurrentHeader() *types.Header {
	return &types.Header{Number: sb.CurrentHeadNumber}
}

func (sb *StubBackend[T,P]) Engine() consensus.Engine[P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SupportsMultitenancy(rpcCtx context.Context) (*proto.PreAuthenticatedAuthenticationToken, bool) {
	return nil, false
}

func (sb *StubBackend[T,P]) AccountExtraDataStateGetterByNumber(context.Context, rpc.BlockNumber) (vm.AccountExtraDataStateGetter, error) {
	return sb.mockAccountExtraDataStateGetter, nil
}

func (sb *StubBackend[T,P]) IsAuthorized(authToken *proto.PreAuthenticatedAuthenticationToken, attributes ...*multitenancy.PrivateStateSecurityAttribute) (bool, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetEVM(ctx context.Context, msg core.Message, state vm.MinimalApiState, header *types.Header, vmconfig *vm.Config[P]) (*vm.EVM[P], func() error, error) {
	sb.getEVMCalled = true
	vmCtx := core.NewEVMBlockContext[P](&types.Header{
		Coinbase:   arbitraryFrom,
		Number:     arbitraryCurrentBlockNumber,
		Time:       0,
		Difficulty: big.NewInt(0),
		GasLimit:   0,
	}, nil, &arbitraryFrom)
	txCtx := core.NewEVMTxContext(msg)
	vmError := func() error {
		return nil
	}
	config := params.QuorumTestChainConfig
	config.IstanbulBlock = sb.IstanbulBlock
	return vm.NewEVM[P](vmCtx, txCtx, publicStateDB, privateStateDB, config, vm.Config[P]{}), vmError, nil
}

func (sb *StubBackend[T,P]) CurrentBlock() *types.Block[P] {
	return types.NewBlock[P](&types.Header{
		Number: arbitraryCurrentBlockNumber,
	}, nil, nil, nil, new(trie.Trie[P]))
}

func (sb *StubBackend[T,P]) Downloader() *downloader.Downloader[T,P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ProtocolVersion() int {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (sb *StubBackend[T,P]) ChainDb() ethdb.Database {
	panic("implement me")
}

func (sb *StubBackend[T,P]) EventMux() *event.TypeMux {
	panic("implement me")
}

func (sb *StubBackend[T,P]) Wallets() []accounts.Wallet[P] {
	return sb.ks.Wallets()
}

func (sb *StubBackend[T,P]) Subscribe(sink chan<- accounts.WalletEvent[P]) event.Subscription {
	return nil
}

func (sb *StubBackend[T,P]) AccountManager() *accounts.Manager[P] {
	return sb.accountManager
}

func (sb *StubBackend[T,P]) ExtRPCEnabled() bool {
	panic("implement me")
}

func (sb *StubBackend[T,P]) CallTimeOut() time.Duration {
	panic("implement me")
}

func (sb *StubBackend[T,P]) RPCTxFeeCap() float64 {
	return 25000000
}

func (sb *StubBackend[T,P]) RPCGasCap() uint64 {
	return 25000000
}

func (sb *StubBackend[T,P]) SetHead(number uint64) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) HeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BlockByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*types.Block[P], error) {
	return sb.CurrentBlock(), nil
}

func (sb *StubBackend[T,P]) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (vm.MinimalApiState, *types.Header, error) {
	return &StubMinimalApiState{}, nil, nil
}

func (sb *StubBackend[T,P]) StateAndHeaderByNumberOrHash(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (vm.MinimalApiState, *types.Header, error) {
	return &StubMinimalApiState{}, nil, nil
}

func (sb *StubBackend[T,P]) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetTd(ctx context.Context, hash common.Hash) *big.Int {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeChainEvent(ch chan<- core.ChainEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SendTx(ctx context.Context, signedTx *types.Transaction[P]) error {
	sb.sendTxCalled = true
	sb.txThatWasSent = signedTx
	return nil
}

func (sb *StubBackend[T,P]) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction[P], common.Hash, uint64, uint64, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetPoolTransactions() (types.Transactions[P], error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetPoolTransaction(txHash common.Hash) *types.Transaction[P] {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return sb.poolNonce, nil
}

func (sb *StubBackend[T,P]) Stats() (pending int, queued int) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) TxPoolContent() (map[common.Address]types.Transactions[P], map[common.Address]types.Transactions[P]) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeNewTxsEvent(chan<- core.NewTxsEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) BloomStatus() (uint64, uint64) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent[P]) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) ChainConfig() *params.ChainConfig {
	return params.QuorumTestChainConfig
}

func (sb *StubBackend[T,P]) SubscribePendingLogsEvent(ch chan<- []*types.Log) event.Subscription {
	panic("implement me")
}

func (sb *StubBackend[T,P]) PSMR() mps.PrivateStateMetadataResolver {
	panic("implement me")
}

type MPSStubBackend [T crypto.PrivateKey, P crypto.PublicKey] struct {
	StubBackend[T,P]
	psmr mps.PrivateStateMetadataResolver
}

func (msb *MPSStubBackend[T,P]) ChainConfig() *params.ChainConfig {
	return params.QuorumMPSTestChainConfig
}

func (sb *MPSStubBackend[T,P]) PSMR() mps.PrivateStateMetadataResolver {
	return sb.psmr
}

func (sb *StubBackend[T,P]) IsPrivacyMarkerTransactionCreationEnabled() bool {
	return sb.isPrivacyMarkerTransactionCreationEnabled
}

type StubMinimalApiState struct {
}

func (StubMinimalApiState) GetBalance(addr common.Address) *big.Int {
	panic("implement me")
}

func (StubMinimalApiState) SetBalance(addr common.Address, balance *big.Int) {
	panic("implement me")
}

func (StubMinimalApiState) GetCode(addr common.Address) []byte {
	return nil
}

func (StubMinimalApiState) GetState(a common.Address, b common.Hash) common.Hash {
	panic("implement me")
}

func (StubMinimalApiState) GetNonce(addr common.Address) uint64 {
	panic("implement me")
}

func (StubMinimalApiState) SetNonce(addr common.Address, nonce uint64) {
	panic("implement me")
}

func (StubMinimalApiState) SetCode(common.Address, []byte) {
	panic("implement me")
}

func (StubMinimalApiState) GetPrivacyMetadata(addr common.Address) (*state.PrivacyMetadata, error) {
	if addr == arbitraryMandatoryRecipientsSimpleStorageContractAddress {
		return &state.PrivacyMetadata{
			CreationTxHash: arbitraryMandatoryRecipientsContractEncryptedPayloadHash,
			PrivacyFlag:    2,
		}, nil
	}
	return &state.PrivacyMetadata{
		CreationTxHash: arbitrarySimpleStorageContractEncryptedPayloadHash,
		PrivacyFlag:    1,
	}, nil
}

func (StubMinimalApiState) GetManagedParties(addr common.Address) ([]string, error) {
	panic("implement me")
}

func (StubMinimalApiState) GetRLPEncodedStateObject(addr common.Address) ([]byte, error) {
	panic("implement me")
}

func (StubMinimalApiState) GetProof(common.Address) ([][]byte, error) {
	panic("implement me")
}

func (StubMinimalApiState) GetStorageProof(common.Address, common.Hash) ([][]byte, error) {
	panic("implement me")
}

func (StubMinimalApiState) StorageTrie(addr common.Address) state.Trie {
	panic("implement me")
}

func (StubMinimalApiState) Error() error {
	panic("implement me")
}

func (StubMinimalApiState) GetCodeHash(common.Address) common.Hash {
	panic("implement me")
}

func (StubMinimalApiState) SetState(common.Address, common.Hash, common.Hash) {
	panic("implement me")
}

func (StubMinimalApiState) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	panic("implement me")
}

type StubPrivateTransactionManager struct {
	notinuse.PrivateTransactionManager
	creation bool
}

func (sptm *StubPrivateTransactionManager) Send(data []byte, from string, to []string, extra *engine.ExtraMetadata) (string, []string, common.EncryptedPayloadHash, error) {
	return "", nil, arbitrarySimpleStorageContractEncryptedPayloadHash, nil
}

func (sptm *StubPrivateTransactionManager) EncryptPayload(data []byte, from string, to []string, extra *engine.ExtraMetadata) ([]byte, error) {
	return nil, engine.ErrPrivateTxManagerNotSupported
}

func (sptm *StubPrivateTransactionManager) DecryptPayload(payload common.DecryptRequest) ([]byte, *engine.ExtraMetadata, error) {
	return nil, nil, engine.ErrPrivateTxManagerNotSupported
}

func (sptm *StubPrivateTransactionManager) StoreRaw(data []byte, from string) (common.EncryptedPayloadHash, error) {
	return arbitrarySimpleStorageContractEncryptedPayloadHash, nil
}

func (sptm *StubPrivateTransactionManager) SendSignedTx(data common.EncryptedPayloadHash, to []string, extra *engine.ExtraMetadata) (string, []string, []byte, error) {
	return "", nil, arbitrarySimpleStorageContractEncryptedPayloadHash.Bytes(), nil
}

func (sptm *StubPrivateTransactionManager) ReceiveRaw(data common.EncryptedPayloadHash) ([]byte, string, *engine.ExtraMetadata, error) {
	if sptm.creation {
		return hexutil.MustDecode("0x6060604052341561000f57600080fd5b604051602080610149833981016040528080519060200190919050505b806000819055505b505b610104806100456000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632a1afcd914605157806360fe47b11460775780636d4ce63c146097575b600080fd5b3415605b57600080fd5b606160bd565b6040518082815260200191505060405180910390f35b3415608157600080fd5b6095600480803590602001909190505060c3565b005b341560a157600080fd5b60a760ce565b6040518082815260200191505060405180910390f35b60005481565b806000819055505b50565b6000805490505b905600a165627a7a72305820d5851baab720bba574474de3d09dbeaabc674a15f4dd93b974908476542c23f00029"), "", nil, nil
	} else {
		return hexutil.MustDecode("0x60fe47b1000000000000000000000000000000000000000000000000000000000000000e"), "", nil, nil
	}
}

func (sptm *StubPrivateTransactionManager) HasFeature(f engine.PrivateTransactionManagerFeature) bool {
	return true
}
