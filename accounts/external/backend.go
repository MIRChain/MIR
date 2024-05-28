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

package external

import (
	"fmt"
	"math/big"
	"sync"

	ethereum "github.com/MIRChain/MIR"
	"github.com/MIRChain/MIR/accounts"
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/common/hexutil"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/event"
	"github.com/MIRChain/MIR/log"
	"github.com/MIRChain/MIR/rpc"
	"github.com/MIRChain/MIR/signer/core"
)

type ExternalBackend[P crypto.PublicKey] struct {
	signers []accounts.Wallet[P]
}

func (eb *ExternalBackend[P]) Wallets() []accounts.Wallet[P] {
	return eb.signers
}

func NewExternalBackend[P crypto.PublicKey](endpoint string) (*ExternalBackend[P], error) {
	signer, err := NewExternalSigner[P](endpoint)
	if err != nil {
		return nil, err
	}
	return &ExternalBackend[P]{
		signers: []accounts.Wallet[P]{signer},
	}, nil
}

func (eb *ExternalBackend[P]) Subscribe(sink chan<- accounts.WalletEvent[P]) event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		<-quit
		return nil
	})
}

// ExternalSigner provides an API to interact with an external signer (clef)
// It proxies request to the external signer while forwarding relevant
// request headers
type ExternalSigner[P crypto.PublicKey] struct {
	client   *rpc.Client
	endpoint string
	status   string
	cacheMu  sync.RWMutex
	cache    []accounts.Account
}

func NewExternalSigner[P crypto.PublicKey](endpoint string) (*ExternalSigner[P], error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	extsigner := &ExternalSigner[P]{
		client:   client,
		endpoint: endpoint,
	}
	// Check if reachable
	version, err := extsigner.pingVersion()
	if err != nil {
		return nil, err
	}
	extsigner.status = fmt.Sprintf("ok [version=%v]", version)
	return extsigner, nil
}

func (api *ExternalSigner[P]) URL() accounts.URL {
	return accounts.URL{
		Scheme: "extapi",
		Path:   api.endpoint,
	}
}

func (api *ExternalSigner[P]) Status() (string, error) {
	return api.status, nil
}

func (api *ExternalSigner[P]) Open(passphrase string) error {
	return fmt.Errorf("operation not supported on external signers")
}

func (api *ExternalSigner[P]) Close() error {
	return fmt.Errorf("operation not supported on external signers")
}

func (api *ExternalSigner[P]) Accounts() []accounts.Account {
	var accnts []accounts.Account
	res, err := api.listAccounts()
	if err != nil {
		log.Error("account listing failed", "error", err)
		return accnts
	}
	for _, addr := range res {
		accnts = append(accnts, accounts.Account{
			URL: accounts.URL{
				Scheme: "extapi",
				Path:   api.endpoint,
			},
			Address: addr,
		})
	}
	api.cacheMu.Lock()
	api.cache = accnts
	api.cacheMu.Unlock()
	return accnts
}

func (api *ExternalSigner[P]) Contains(account accounts.Account) bool {
	api.cacheMu.RLock()
	defer api.cacheMu.RUnlock()
	if api.cache == nil {
		// If we haven't already fetched the accounts, it's time to do so now
		api.cacheMu.RUnlock()
		api.Accounts()
		api.cacheMu.RLock()
	}
	for _, a := range api.cache {
		if a.Address == account.Address && (account.URL == (accounts.URL{}) || account.URL == api.URL()) {
			return true
		}
	}
	return false
}

func (api *ExternalSigner[P]) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, fmt.Errorf("operation not supported on external signers")
}

func (api *ExternalSigner[P]) SelfDerive(bases []accounts.DerivationPath, chain ethereum.ChainStateReader) {
	log.Error("operation SelfDerive not supported on external signers")
}

func (api *ExternalSigner[P]) signHash(account accounts.Account, hash []byte) ([]byte, error) {
	return []byte{}, fmt.Errorf("operation not supported on external signers")
}

// SignData signs keccak256(data). The mimetype parameter describes the type of data being signed
func (api *ExternalSigner[P]) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	var res hexutil.Bytes
	var signAddress = common.NewMixedcaseAddress(account.Address)
	if err := api.client.Call(&res, "account_signData",
		mimeType,
		&signAddress, // Need to use the pointer here, because of how MarshalJSON is defined
		hexutil.Encode(data)); err != nil {
		return nil, err
	}
	// If V is on 27/28-form, convert to 0/1 for Clique
	if mimeType == accounts.MimetypeClique && (res[64] == 27 || res[64] == 28) {
		res[64] -= 27 // Transform V from 27/28 to 0/1 for Clique use
	}
	return res, nil
}

func (api *ExternalSigner[P]) SignText(account accounts.Account, text []byte) ([]byte, error) {
	var signature hexutil.Bytes
	var signAddress = common.NewMixedcaseAddress(account.Address)
	if err := api.client.Call(&signature, "account_signData",
		accounts.MimetypeTextPlain,
		&signAddress, // Need to use the pointer here, because of how MarshalJSON is defined
		hexutil.Encode(text)); err != nil {
		return nil, err
	}
	if signature[64] == 27 || signature[64] == 28 {
		// If clef is used as a backend, it may already have transformed
		// the signature to ethereum-type signature.
		signature[64] -= 27 // Transform V from Ethereum-legacy to 0/1
	}
	return signature, nil
}

// signTransactionResult represents the signinig result returned by clef.
type signTransactionResult[P crypto.PublicKey] struct {
	Raw hexutil.Bytes         `json:"raw"`
	Tx  *types.Transaction[P] `json:"tx"`
}

func (api *ExternalSigner[P]) SignTx(account accounts.Account, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	data := hexutil.Bytes(tx.Data())
	var to *common.MixedcaseAddress
	if tx.To() != nil {
		t := common.NewMixedcaseAddress(*tx.To())
		to = &t
	}
	args := &core.SendTxArgs[P]{
		Data:     &data,
		Nonce:    hexutil.Uint64(tx.Nonce()),
		Value:    hexutil.Big(*tx.Value()),
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: hexutil.Big(*tx.GasPrice()),
		To:       to,
		From:     common.NewMixedcaseAddress(account.Address),
		// Quorum
		IsPrivate: tx.IsPrivate(),
	}
	// We should request the default chain id that we're operating with
	// (the chain we're executing on)
	if chainID != nil {
		args.ChainID = (*hexutil.Big)(chainID)
	}
	// However, if the user asked for a particular chain id, then we should
	// use that instead.
	if tx.Type() != types.LegacyTxType && tx.ChainId() != nil {
		args.ChainID = (*hexutil.Big)(tx.ChainId())
	}
	if tx.Type() == types.AccessListTxType {
		accessList := tx.AccessList()
		args.AccessList = &accessList
	}
	var res signTransactionResult[P]
	if err := api.client.Call(&res, "account_signTransaction", args); err != nil {
		return nil, err
	}
	return res.Tx, nil
}

func (api *ExternalSigner[P]) SignTextWithPassphrase(account accounts.Account, passphrase string, text []byte) ([]byte, error) {
	return []byte{}, fmt.Errorf("password-operations not supported on external signers")
}

func (api *ExternalSigner[P]) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	return nil, fmt.Errorf("password-operations not supported on external signers")
}
func (api *ExternalSigner[P]) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	return nil, fmt.Errorf("password-operations not supported on external signers")
}

func (api *ExternalSigner[P]) listAccounts() ([]common.Address, error) {
	var res []common.Address
	if err := api.client.Call(&res, "account_list"); err != nil {
		return nil, err
	}
	return res, nil
}

func (api *ExternalSigner[P]) pingVersion() (string, error) {
	var v string
	if err := api.client.Call(&v, "account_version"); err != nil {
		return "", err
	}
	return v, nil
}
