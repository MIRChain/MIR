package core

import (
	"math/big"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/params"
)

// callHelper makes it easier to do proper calls and use the state transition object.
// It also manages the nonces of the caller and keeps private and public state, which
// can be freely modified outside of the helper.
type callHelper [T crypto.PrivateKey, P crypto.PublicKey] struct {
	db ethdb.Database

	nonces map[common.Address]uint64
	header types.Header[P]
	gp     *GasPool

	PrivateState, PublicState *state.StateDB[P]
}

// TxNonce returns the pending nonce
func (cg *callHelper[T,P]) TxNonce(addr common.Address) uint64 {
	return cg.nonces[addr]
}

// MakeCall makes does a call to the recipient using the given input. It can switch between private and public
// by setting the private boolean flag. It returns an error if the call failed.
func (cg *callHelper[T,P]) MakeCall(private bool, key T, to common.Address, input []byte) error {
	var pub P
	switch t:=any(&key).(type){
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
	}
	var (
		from = crypto.PubkeyToAddress[P](pub)
		err  error
	)

	// TODO(joel): these are just stubbed to the same values as in dual_state_test.go
	cg.header.Number = new(big.Int)
	cg.header.Time = uint64(43)
	cg.header.Difficulty = new(big.Int).SetUint64(1000488)
	cg.header.GasLimit = 4700000

	signer := types.MakeSigner[P](params.QuorumTestChainConfig, cg.header.Number)
	if private {
		signer = types.QuorumPrivateTxSigner[P]{}
	}

	tx, err := types.SignTx(types.NewTransaction[P](cg.TxNonce(from), to, new(big.Int), 1000000, new(big.Int), input), signer, key)

	if err != nil {
		return err
	}
	defer func() { cg.nonces[from]++ }()
	msg, err := tx.AsMessage(signer)
	if err != nil {
		return err
	}

	publicState, privateState := cg.PublicState, cg.PrivateState
	if !private {
		privateState = publicState
	}
	// TODO(joel): can we just pass nil instead of bc?
	bc, _ := NewBlockChain[P](cg.db, nil, params.QuorumTestChainConfig, ethash.NewFaker[P](), vm.Config[P]{}, nil, nil, nil)
	txContext := NewEVMTxContext(msg)
	evmContext := NewEVMBlockContext[P](&cg.header, bc, &from)
	vmenv := vm.NewEVM(evmContext, txContext, publicState, privateState, params.QuorumTestChainConfig, vm.Config[P]{})
	sender := vm.AccountRef(msg.From())
	vmenv.Call(sender, to, msg.Data(), 100000000, new(big.Int))
	return err
}

// MakeCallHelper returns a new callHelper
func MakeCallHelper[T crypto.PrivateKey, P crypto.PublicKey]() *callHelper[T,P] {
	memdb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase[P](memdb)

	publicState, err := state.New[P](common.Hash{}, db, nil)
	if err != nil {
		panic(err)
	}
	privateState, err := state.New[P](common.Hash{}, db, nil)
	if err != nil {
		panic(err)
	}
	cg := &callHelper[T,P]{
		db:           memdb,
		nonces:       make(map[common.Address]uint64),
		gp:           new(GasPool).AddGas(5000000),
		PublicState:  publicState,
		PrivateState: privateState,
	}
	return cg
}
