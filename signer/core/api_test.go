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

package core_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/signer/core"
	"github.com/pavelkrolevets/MIR-pro/signer/fourbyte"
	"github.com/pavelkrolevets/MIR-pro/signer/storage"
)

//Used for testing
type headlessUi [T crypto.PrivateKey, P crypto.PublicKey] struct {
	approveCh chan string // to send approve/deny
	inputCh   chan string // to send password
}

func (ui *headlessUi[T,P]) OnInputRequired(info core.UserInputRequest) (core.UserInputResponse, error) {
	input := <-ui.inputCh
	return core.UserInputResponse{Text: input}, nil
}

func (ui *headlessUi[T,P]) OnSignerStartup(info core.StartupInfo)        {}
func (ui *headlessUi[T,P]) RegisterUIServer(api *core.UIServerAPI[T,P])       {}
func (ui *headlessUi[T,P]) OnApprovedTx(tx ethapi.SignTransactionResult[P]) {}

func (ui *headlessUi[T,P]) ApproveTx(request *core.SignTxRequest[P]) (core.SignTxResponse[P], error) {

	switch <-ui.approveCh {
	case "Y":
		return core.SignTxResponse[P]{request.Transaction, true}, nil
	case "M": // modify
		// The headless UI always modifies the transaction
		old := big.Int(request.Transaction.Value)
		newVal := big.NewInt(0).Add(&old, big.NewInt(1))
		request.Transaction.Value = hexutil.Big(*newVal)
		return core.SignTxResponse[P]{request.Transaction, true}, nil
	default:
		return core.SignTxResponse[P]{request.Transaction, false}, nil
	}
}

func (ui *headlessUi[T,P]) ApproveSignData(request *core.SignDataRequest) (core.SignDataResponse, error) {
	approved := (<-ui.approveCh == "Y")
	return core.SignDataResponse{approved}, nil
}

func (ui *headlessUi[T,P]) ApproveListing(request *core.ListRequest) (core.ListResponse, error) {
	approval := <-ui.approveCh
	//fmt.Printf("approval %s\n", approval)
	switch approval {
	case "A":
		return core.ListResponse{request.Accounts}, nil
	case "1":
		l := make([]accounts.Account, 1)
		l[0] = request.Accounts[1]
		return core.ListResponse{l}, nil
	default:
		return core.ListResponse{nil}, nil
	}
}

func (ui *headlessUi[T,P]) ApproveNewAccount(request *core.NewAccountRequest) (core.NewAccountResponse, error) {
	if <-ui.approveCh == "Y" {
		return core.NewAccountResponse{true}, nil
	}
	return core.NewAccountResponse{false}, nil
}

func (ui *headlessUi[T,P]) ShowError(message string) {
	//stdout is used by communication
	fmt.Fprintln(os.Stderr, message)
}

func (ui *headlessUi[T,P]) ShowInfo(message string) {
	//stdout is used by communication
	fmt.Fprintln(os.Stderr, message)
}

func tmpDirName(t *testing.T) string {
	d, err := ioutil.TempDir("", "eth-keystore-test")
	if err != nil {
		t.Fatal(err)
	}
	d, err = filepath.EvalSymlinks(d)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func setup[T crypto.PrivateKey, P crypto.PublicKey](t *testing.T) (*core.SignerAPI[T,P], *headlessUi[T,P]) {
	db, err := fourbyte.New[P]()
	if err != nil {
		t.Fatal(err.Error())
	}
	ui := &headlessUi[T,P]{make(chan string, 20), make(chan string, 20)}
	am := core.StartClefAccountManager[T,P](tmpDirName(t), true, true, nil, "")
	api := core.NewSignerAPI[T,P](am, 1337, true, ui, db, true, &storage.NoStorage{})
	return api, ui

}
func createAccount[T crypto.PrivateKey, P crypto.PublicKey](ui *headlessUi[T,P], api *core.SignerAPI[T,P], t *testing.T) {
	ui.approveCh <- "Y"
	ui.inputCh <- "a_long_password"
	_, err := api.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// Some time to allow changes to propagate
	time.Sleep(250 * time.Millisecond)
}

func failCreateAccountWithPassword[T crypto.PrivateKey, P crypto.PublicKey](ui *headlessUi[T,P], api *core.SignerAPI[T,P], password string, t *testing.T) {

	ui.approveCh <- "Y"
	// We will be asked three times to provide a suitable password
	ui.inputCh <- password
	ui.inputCh <- password
	ui.inputCh <- password

	addr, err := api.New(context.Background())
	if err == nil {
		t.Fatal("Should have returned an error")
	}
	if addr != (common.Address{}) {
		t.Fatal("Empty address should be returned")
	}
}

func failCreateAccount[T crypto.PrivateKey, P crypto.PublicKey](ui *headlessUi[T,P], api *core.SignerAPI[T,P], t *testing.T) {
	ui.approveCh <- "N"
	addr, err := api.New(context.Background())
	if err != core.ErrRequestDenied {
		t.Fatal(err)
	}
	if addr != (common.Address{}) {
		t.Fatal("Empty address should be returned")
	}
}

func list[T crypto.PrivateKey, P crypto.PublicKey](ui *headlessUi[T,P], api *core.SignerAPI[T,P], t *testing.T) ([]common.Address, error) {
	ui.approveCh <- "A"
	return api.List(context.Background())

}

func TestNewAcc(t *testing.T) {
	api, control := setup[nist.PrivateKey,nist.PublicKey](t)
	verifyNum := func(num int) {
		list, err := list(control, api, t)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if len(list) != num {
			t.Errorf("Expected %d accounts, got %d", num, len(list))
		}
	}
	// Testing create and create-deny
	createAccount(control, api, t)
	createAccount(control, api, t)
	failCreateAccount(control, api, t)
	failCreateAccount(control, api, t)
	createAccount(control, api, t)
	failCreateAccount(control, api, t)
	createAccount(control, api, t)
	failCreateAccount(control, api, t)
	verifyNum(4)

	// Fail to create this, due to bad password
	failCreateAccountWithPassword(control, api, "short", t)
	failCreateAccountWithPassword(control, api, "longerbutbad\rfoo", t)
	verifyNum(4)

	// Testing listing:
	// Listing one Account
	control.approveCh <- "1"
	list, err := api.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("List should only show one Account")
	}
	// Listing denied
	control.approveCh <- "Nope"
	list, err = api.List(context.Background())
	if len(list) != 0 {
		t.Fatalf("List should be empty")
	}
	if err != core.ErrRequestDenied {
		t.Fatal("Expected deny")
	}
}

func mkTestTx[P crypto.PublicKey](from common.MixedcaseAddress) core.SendTxArgs[P] {
	to := common.NewMixedcaseAddress(common.HexToAddress("0x1337"))
	gas := hexutil.Uint64(21000)
	gasPrice := (hexutil.Big)(*big.NewInt(2000000000))
	value := (hexutil.Big)(*big.NewInt(1e18))
	nonce := (hexutil.Uint64)(0)
	data := hexutil.Bytes(common.Hex2Bytes("01020304050607080a"))
	tx := core.SendTxArgs[P]{
		From:     from,
		To:       &to,
		Gas:      gas,
		GasPrice: gasPrice,
		Value:    value,
		Data:     &data,
		Nonce:    nonce}
	return tx
}

func TestSignTx(t *testing.T) {
	var (
		list      []common.Address
		res, res2 *ethapi.SignTransactionResult[nist.PublicKey]
		err       error
	)

	api, control := setup[nist.PrivateKey,nist.PublicKey](t)
	createAccount(control, api, t)
	control.approveCh <- "A"
	list, err = api.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	a := common.NewMixedcaseAddress(list[0])

	methodSig := "test(uint)"
	tx := mkTestTx[nist.PublicKey](a)

	control.approveCh <- "Y"
	control.inputCh <- "wrongpassword"
	res, err = api.SignTransaction(context.Background(), tx, &methodSig)
	if res != nil {
		t.Errorf("Expected nil-response, got %v", res)
	}
	if err != keystore.ErrDecrypt {
		t.Errorf("Expected ErrLocked! %v", err)
	}
	control.approveCh <- "No way"
	res, err = api.SignTransaction(context.Background(), tx, &methodSig)
	if res != nil {
		t.Errorf("Expected nil-response, got %v", res)
	}
	if err != core.ErrRequestDenied {
		t.Errorf("Expected ErrRequestDenied! %v", err)
	}
	// Sign with correct password
	control.approveCh <- "Y"
	control.inputCh <- "a_long_password"
	res, err = api.SignTransaction(context.Background(), tx, &methodSig)

	if err != nil {
		t.Fatal(err)
	}
	parsedTx := &types.Transaction[nist.PublicKey]{}
	rlp.Decode(bytes.NewReader(res.Raw), parsedTx)

	//The tx should NOT be modified by the UI
	if parsedTx.Value().Cmp(tx.Value.ToInt()) != 0 {
		t.Errorf("Expected value to be unchanged, expected %v got %v", tx.Value, parsedTx.Value())
	}
	control.approveCh <- "Y"
	control.inputCh <- "a_long_password"

	res2, err = api.SignTransaction(context.Background(), tx, &methodSig)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(res.Raw, res2.Raw) {
		t.Error("Expected tx to be unmodified by UI")
	}

	//The tx is modified by the UI
	control.approveCh <- "M"
	control.inputCh <- "a_long_password"

	res2, err = api.SignTransaction(context.Background(), tx, &methodSig)
	if err != nil {
		t.Fatal(err)
	}
	parsedTx2 := &types.Transaction[nist.PublicKey]{}
	rlp.Decode(bytes.NewReader(res.Raw), parsedTx2)

	//The tx should be modified by the UI
	if parsedTx2.Value().Cmp(tx.Value.ToInt()) != 0 {
		t.Errorf("Expected value to be unchanged, got %v", parsedTx.Value())
	}
	if bytes.Equal(res.Raw, res2.Raw) {
		t.Error("Expected tx to be modified by UI")
	}

}
