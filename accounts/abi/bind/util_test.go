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

package bind_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi/bind/backends"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
)

var testKey, _ = crypto.HexToECDSA[nist.PrivateKey]("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")

var waitDeployedTests = map[string]struct {
	code        string
	gas         uint64
	wantAddress common.Address
	wantErr     error
}{
	"successful deploy": {
		code:        `6060604052600a8060106000396000f360606040526008565b00`,
		gas:         3000000,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
	"empty code": {
		code:        ``,
		gas:         300000,
		wantErr:     bind.ErrNoCodeAfterDeploy,
		wantAddress: common.HexToAddress("0x3a220f351252089d385b29beca14e27f204c296a"),
	},
}

func TestWaitDeployed(t *testing.T) {
	for name, test := range waitDeployedTests {
		backend := backends.NewSimulatedBackend[nist.PublicKey](
			core.GenesisAlloc{
				crypto.PubkeyToAddress[nist.PublicKey](*testKey.Public()): {Balance: big.NewInt(10000000000)},
			},
			10000000,
		)
		defer backend.Close()

		// Create the transaction.
		tx := types.NewContractCreation[nist.PublicKey](0, big.NewInt(0), test.gas, big.NewInt(1), common.FromHex(test.code))
		tx, _ = types.SignTx[nist.PrivateKey,nist.PublicKey](tx, types.HomesteadSigner[nist.PublicKey]{}, testKey)

		// Wait for it to get mined in the background.
		var (
			err     error
			address common.Address
			mined   = make(chan struct{})
			ctx     = context.Background()
		)
		go func() {
			address, err = bind.WaitDeployed[nist.PublicKey](ctx, backend, tx)
			close(mined)
		}()

		// Send and mine the transaction.
		backend.SendTransaction(ctx, tx, bind.PrivateTxArgs{})
		backend.Commit()

		select {
		case <-mined:
			if err != test.wantErr {
				t.Errorf("test %q: error mismatch: want %q, got %q", name, test.wantErr, err)
			}
			if address != test.wantAddress {
				t.Errorf("test %q: unexpected contract address %s", name, address.Hex())
			}
		case <-time.After(2 * time.Second):
			t.Errorf("test %q: timeout", name)
		}
	}
}

func TestWaitDeployedCornerCases(t *testing.T) {
	backend := backends.NewSimulatedBackend[nist.PublicKey](
		core.GenesisAlloc{
			crypto.PubkeyToAddress[nist.PublicKey](*testKey.Public()): {Balance: big.NewInt(10000000000)},
		},
		10000000,
	)
	defer backend.Close()

	// Create a transaction to an account.
	code := "6060604052600a8060106000396000f360606040526008565b00"
	tx := types.NewTransaction[nist.PublicKey](0, common.HexToAddress("0x01"), big.NewInt(0), 3000000, big.NewInt(1), common.FromHex(code))
	tx, _ = types.SignTx[nist.PrivateKey,nist.PublicKey](tx, types.HomesteadSigner[nist.PublicKey]{}, testKey)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	backend.SendTransaction(ctx, tx, bind.PrivateTxArgs{})
	backend.Commit()
	notContentCreation := errors.New("tx is not contract creation")
	if _, err := bind.WaitDeployed[nist.PublicKey](ctx, backend, tx); err.Error() != notContentCreation.Error() {
		t.Errorf("error missmatch: want %q, got %q, ", notContentCreation, err)
	}

	// Create a transaction that is not mined.
	tx = types.NewContractCreation[nist.PublicKey](1, big.NewInt(0), 3000000, big.NewInt(1), common.FromHex(code))
	tx, _ = types.SignTx[nist.PrivateKey,nist.PublicKey](tx, types.HomesteadSigner[nist.PublicKey]{}, testKey)

	go func() {
		contextCanceled := errors.New("context canceled")
		if _, err := bind.WaitDeployed[nist.PublicKey](ctx, backend, tx); err.Error() != contextCanceled.Error() {
			t.Errorf("error missmatch: want %q, got %q, ", contextCanceled, err)
		}
	}()

	backend.SendTransaction(ctx, tx, bind.PrivateTxArgs{})
	cancel()
}
