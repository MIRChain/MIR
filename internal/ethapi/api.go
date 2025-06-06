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

package ethapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/abi"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/accounts/pluggable"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/consensus/clique"
	"github.com/pavelkrolevets/MIR-pro/consensus/ethash"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/multitenancy"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/private/engine"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

type TransactionType uint8

const (
	FillTransaction TransactionType = iota + 1
	RawTransaction
	NormalTransaction
)

// PublicEthereumAPI provides an API to access Ethereum related information.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicEthereumAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	b Backend[T,P]
}

// NewPublicEthereumAPI creates a new Ethereum protocol API.
func NewPublicEthereumAPI [T crypto.PrivateKey, P crypto.PublicKey](b Backend[T,P]) *PublicEthereumAPI[T,P] {
	return &PublicEthereumAPI[T,P]{b}
}

// GasPrice returns a suggestion for a gas price.
func (s *PublicEthereumAPI[T,P]) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	price, err := s.b.SuggestPrice(ctx)
	return (*hexutil.Big)(price), err
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (s *PublicEthereumAPI[T,P]) Syncing() (interface{}, error) {
	progress := s.b.Downloader().Progress()

	// Return not syncing if the synchronisation already completed
	if progress.CurrentBlock >= progress.HighestBlock {
		return false, nil
	}
	// Otherwise gather the block sync stats
	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(progress.StartingBlock),
		"currentBlock":  hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":  hexutil.Uint64(progress.HighestBlock),
		"pulledStates":  hexutil.Uint64(progress.PulledStates),
		"knownStates":   hexutil.Uint64(progress.KnownStates),
	}, nil
}

func (s *PublicEthereumAPI[T,P]) GetPrivacyPrecompileAddress() common.Address {
	return common.QuorumPrivacyPrecompileContractAddress()
}

// PublicTxPoolAPI offers and API for the transaction pool. It only operates on data that is non confidential.
type PublicTxPoolAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	b Backend[T,P]
}

// NewPublicTxPoolAPI creates a new tx pool service that gives information about the transaction pool.
func NewPublicTxPoolAPI [T crypto.PrivateKey, P crypto.PublicKey] (b Backend[T,P]) *PublicTxPoolAPI[T,P] {
	return &PublicTxPoolAPI[T,P]{b}
}

// Content returns the transactions contained within the transaction pool.
func (s *PublicTxPoolAPI[T,P]) Content() map[string]map[string]map[string]*RPCTransaction {
	content := map[string]map[string]map[string]*RPCTransaction{
		"pending": make(map[string]map[string]*RPCTransaction),
		"queued":  make(map[string]map[string]*RPCTransaction),
	}
	pending, queue := s.b.TxPoolContent()

	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx)
		}
		content["pending"][account.Hex()] = dump
	}
	// Flatten the queued transactions
	for account, txs := range queue {
		dump := make(map[string]*RPCTransaction)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = newRPCPendingTransaction(tx)
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}

// Status returns the number of pending and queued transaction in the pool.
func (s *PublicTxPoolAPI[T,P]) Status() map[string]hexutil.Uint {
	pending, queue := s.b.Stats()
	return map[string]hexutil.Uint{
		"pending": hexutil.Uint(pending),
		"queued":  hexutil.Uint(queue),
	}
}

// Inspect retrieves the content of the transaction pool and flattens it into an
// easily inspectable list.
func (s *PublicTxPoolAPI[T,P]) Inspect() map[string]map[string]map[string]string {
	content := map[string]map[string]map[string]string{
		"pending": make(map[string]map[string]string),
		"queued":  make(map[string]map[string]string),
	}
	pending, queue := s.b.TxPoolContent()

	// Define a formatter to flatten a transaction into a string
	var format = func(tx *types.Transaction[P]) string {
		if to := tx.To(); to != nil {
			return fmt.Sprintf("%s: %v wei + %v gas × %v wei", tx.To().Hex(), tx.Value(), tx.Gas(), tx.GasPrice())
		}
		return fmt.Sprintf("contract creation: %v wei + %v gas × %v wei", tx.Value(), tx.Gas(), tx.GasPrice())
	}
	// Flatten the pending transactions
	for account, txs := range pending {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["pending"][account.Hex()] = dump
	}
	// Flatten the queued transactions
	for account, txs := range queue {
		dump := make(map[string]string)
		for _, tx := range txs {
			dump[fmt.Sprintf("%d", tx.Nonce())] = format(tx)
		}
		content["queued"][account.Hex()] = dump
	}
	return content
}

// PublicAccountAPI provides an API to access accounts managed by this node.
// It offers only methods that can retrieve accounts.
type PublicAccountAPI [P crypto.PublicKey] struct {
	am *accounts.Manager[P]
}

// NewPublicAccountAPI creates a new PublicAccountAPI.
func NewPublicAccountAPI[P crypto.PublicKey](am *accounts.Manager[P]) *PublicAccountAPI[P] {
	return &PublicAccountAPI[P]{am: am}
}

// Accounts returns the collection of accounts this node manages
func (s *PublicAccountAPI[P]) Accounts() []common.Address {
	return s.am.Accounts()
}

// PrivateAccountAPI provides an API to access accounts managed by this node.
// It offers methods to create, (un)lock en list accounts. Some methods accept
// passwords and are therefore considered private by default.
type PrivateAccountAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	am        *accounts.Manager[P]
	nonceLock *AddrLocker
	b         Backend[T,P]
}

// NewPrivateAccountAPI create a new PrivateAccountAPI.
func NewPrivateAccountAPI[T crypto.PrivateKey, P crypto.PublicKey](b Backend[T,P], nonceLock *AddrLocker) *PrivateAccountAPI[T,P] {
	return &PrivateAccountAPI[T,P]{
		am:        b.AccountManager(),
		nonceLock: nonceLock,
		b:         b,
	}
}

// listAccounts will return a list of addresses for accounts this node manages.
func (s *PrivateAccountAPI[T,P]) ListAccounts() []common.Address {
	return s.am.Accounts()
}

// rawWallet is a JSON representation of an accounts.Wallet interface, with its
// data contents extracted into plain fields.
type rawWallet struct {
	URL      string             `json:"url"`
	Status   string             `json:"status"`
	Failure  string             `json:"failure,omitempty"`
	Accounts []accounts.Account `json:"accounts,omitempty"`
}

// ListWallets will return a list of wallets this node manages.
func (s *PrivateAccountAPI[T,P]) ListWallets() []rawWallet {
	wallets := make([]rawWallet, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		status, failure := wallet.Status()

		raw := rawWallet{
			URL:      wallet.URL().String(),
			Status:   status,
			Accounts: wallet.Accounts(),
		}
		if failure != nil {
			raw.Failure = failure.Error()
		}
		wallets = append(wallets, raw)
	}
	return wallets
}

// OpenWallet initiates a hardware wallet opening procedure, establishing a USB
// connection and attempting to authenticate via the provided passphrase. Note,
// the method may return an extra challenge requiring a second open (e.g. the
// Trezor PIN matrix challenge).
func (s *PrivateAccountAPI[T,P]) OpenWallet(url string, passphrase *string) error {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return err
	}
	pass := ""
	if passphrase != nil {
		pass = *passphrase
	}
	return wallet.Open(pass)
}

// DeriveAccount requests a HD wallet to derive a new account, optionally pinning
// it for later reuse.
func (s *PrivateAccountAPI[T,P]) DeriveAccount(url string, path string, pin *bool) (accounts.Account, error) {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return accounts.Account{}, err
	}
	derivPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return accounts.Account{}, err
	}
	if pin == nil {
		pin = new(bool)
	}
	return wallet.Derive(derivPath, *pin)
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PrivateAccountAPI[T,P]) NewAccount(password string) (common.Address, error) {
	ks, err := fetchKeystore[T,P](s.am)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := ks.NewAccount(password)
	if err == nil {
		log.Info("Your new key was generated", "address", acc.Address)
		log.Warn("Please backup your key file!", "path", acc.URL.Path)
		log.Warn("Please remember your password!")
		return acc.Address, nil
	}
	return common.Address{}, err
}

// fetchKeystore retrieves the encrypted keystore from the account manager.
func fetchKeystore[T crypto.PrivateKey, P crypto.PublicKey](am *accounts.Manager[P]) (*keystore.KeyStore[T,P], error) {
	if ks := am.Backends(reflect.TypeOf(&keystore.KeyStore[T,P]{})); len(ks) > 0 {
		return ks[0].(*keystore.KeyStore[T,P]), nil
	}
	return nil, errors.New("local keystore not used")
}

// ImportRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.
func (s *PrivateAccountAPI[T,P]) ImportRawKey(privkey string, password string) (common.Address, error) {
	key, err := crypto.HexToECDSA[T](privkey)
	if err != nil {
		return common.Address{}, err
	}
	ks, err := fetchKeystore[T,P](s.am)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := ks.ImportECDSA(key, password)
	return acc.Address, err
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
func (s *PrivateAccountAPI[T,P]) UnlockAccount(ctx context.Context, addr common.Address, password string, duration *uint64) (bool, error) {
	// When the API is exposed by external RPC(http, ws etc), unless the user
	// explicitly specifies to allow the insecure account unlocking, otherwise
	// it is disabled.
	if s.b.ExtRPCEnabled() && !s.b.AccountManager().Config().InsecureUnlockAllowed {
		return false, errors.New("account unlock with HTTP access is forbidden")
	}

	const max = uint64(time.Duration(math.MaxInt64) / time.Second)
	var d time.Duration
	if duration == nil {
		d = 300 * time.Second
	} else if *duration > max {
		return false, errors.New("unlock duration too large")
	} else {
		d = time.Duration(*duration) * time.Second
	}
	err := s.unlockAccount(addr, password, d)
	if err != nil {
		log.Warn("Failed account unlock attempt", "address", addr, "err", err)
	}
	return err == nil, err
}

func (s *PrivateAccountAPI[T,P]) unlockAccount(addr common.Address, password string, duration time.Duration) error {
	acct := accounts.Account{Address: addr}

	backend, err := s.am.Backend(acct)
	if err != nil {
		return err
	}

	switch b := backend.(type) {
	case *pluggable.Backend[P]:
		return b.TimedUnlock(acct, password, duration)
	case *keystore.KeyStore[T,P]:
		return b.TimedUnlock(acct, password, duration)
	default:
		return errors.New("unlock only supported for keystore or plugin wallets")
	}
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (s *PrivateAccountAPI[T,P]) LockAccount(addr common.Address) bool {
	if err := s.lockAccount(addr); err != nil {
		log.Warn("Failed account lock attempt", "address", addr, "err", err)
		return false
	}

	return true
}

func (s *PrivateAccountAPI[T,P]) lockAccount(addr common.Address) error {
	acct := accounts.Account{Address: addr}

	backend, err := s.am.Backend(acct)
	if err != nil {
		return err
	}

	switch b := backend.(type) {
	case *pluggable.Backend[P]:
		return b.Lock(acct)
	case *keystore.KeyStore[T,P]:
		return b.Lock(addr)
	default:
		return errors.New("lock only supported for keystore or plugin wallets")
	}
}

// signTransaction sets defaults and signs the given transaction
// NOTE: the caller needs to ensure that the nonceLock is held, if applicable,
// and release it after the transaction has been submitted to the tx pool
func (s *PrivateAccountAPI[T,P]) signTransaction(ctx context.Context, args *SendTxArgs[T,P], passwd string) (*types.Transaction[P], error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: args.From}
	wallet, err := s.am.Find(account)
	if err != nil {
		return nil, err
	}
	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and sign with the wallet
	tx := args.toTransaction()

	// Quorum
	if args.IsPrivate() {
		tx.SetPrivate()
	}
	var chainID *big.Int
	if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) && !tx.IsPrivate() {
		chainID = config.ChainID
	}
	// /Quorum

	return wallet.SignTxWithPassphrase(account, passwd, tx, chainID)
}

// SendTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given passwd isn't
// able to decrypt the key it fails.
func (s *PrivateAccountAPI[T,P]) SendTransaction(ctx context.Context, args SendTxArgs[T,P], passwd string) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}

	// Quorum
	_, replaceDataWithHash, data, err := checkAndHandlePrivateTransaction(ctx, s.b, args.toTransaction(), &args.PrivateTxArgs, args.From, NormalTransaction)
	if err != nil {
		return common.Hash{}, err
	}
	if replaceDataWithHash {
		// replace the original payload with encrypted payload hash
		args.Data = data.BytesTypeRef()
	}
	// /Quorum

	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction send attempt", "from", args.From, "to", args.To, "value", args.Value.ToInt(), "err", err)
		return common.Hash{}, err
	}

	// Quorum
	if signed.IsPrivate() && s.b.IsPrivacyMarkerTransactionCreationEnabled() {
		// Look up the wallet containing the requested signer
		account := accounts.Account{Address: args.From}
		wallet, err := s.am.Find(account)
		if err != nil {
			return common.Hash{}, err
		}

		pmt, err := createPrivacyMarkerTransaction(s.b, signed, &args.PrivateTxArgs)
		if err != nil {
			log.Warn("Failed to create privacy marker transaction for private transaction", "from", args.From, "to", args.To, "value", args.Value.ToInt(), "err", err)
			return common.Hash{}, err
		}

		var pmtChainID *big.Int // PMT is public so will have different chainID used in signing compared to the internal tx
		if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) {
			pmtChainID = config.ChainID
		}

		signed, err = wallet.SignTxWithPassphrase(account, passwd, pmt, pmtChainID)
		if err != nil {
			log.Warn("Failed to sign privacy marker transaction for private transaction", "from", args.From, "to", args.To, "value", args.Value.ToInt(), "err", err)
			return common.Hash{}, err
		}
	}
	// /Quorum

	return SubmitTransaction(ctx, s.b, signed, args.PrivateFrom, false)
}

// SignTransaction will create a transaction from the given arguments and
// tries to sign it with the key associated with args.From. If the given passwd isn't
// able to decrypt the key it fails. The transaction is returned in RLP-form, not broadcast
// to other nodes
func (s *PrivateAccountAPI[T,P]) SignTransaction(ctx context.Context, args SendTxArgs[T,P], passwd string) (*SignTransactionResult[P], error) {
	// No need to obtain the noncelock mutex, since we won't be sending this
	// tx into the transaction pool, but right back to the user
	if args.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}
	if args.GasPrice == nil {
		return nil, fmt.Errorf("gasPrice not specified")
	}
	if args.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}
	// Before actually sign the transaction, ensure the transaction fee is reasonable.
	if err := checkTxFee(args.GasPrice.ToInt(), uint64(*args.Gas), s.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	signed, err := s.signTransaction(ctx, &args, passwd)
	if err != nil {
		log.Warn("Failed transaction sign attempt", "from", args.From, "to", args.To, "value", args.Value.ToInt(), "err", err)
		return nil, err
	}
	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult[P]{data, signed}, nil
}

// Sign calculates an Ethereum ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/pavelkrolevets/MIR-pro/wiki/Management-APIs#personal_sign
func (s *PrivateAccountAPI[T,P]) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Assemble sign the data with the wallet
	signature, err := wallet.SignTextWithPassphrase(account, passwd, data)
	if err != nil {
		log.Warn("Failed data sign attempt", "address", addr, "err", err)
		return nil, err
	}
	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// EcRecover returns the address for the account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be 27 or 28 for legacy reasons.
//
// https://github.com/pavelkrolevets/MIR-pro/wiki/Management-APIs#personal_ecRecover
func (s *PrivateAccountAPI[T,P]) EcRecover(ctx context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.SigToPub[P](accounts.TextHash(data), sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(rpk), nil
}

// SignAndSendTransaction was renamed to SendTransaction. This method is deprecated
// and will be removed in the future. It primary goal is to give clients time to update.
func (s *PrivateAccountAPI[T,P]) SignAndSendTransaction(ctx context.Context, args SendTxArgs[T,P], passwd string) (common.Hash, error) {
	return s.SendTransaction(ctx, args, passwd)
}

// InitializeWallet initializes a new wallet at the provided URL, by generating and returning a new private key.
// func (s *PrivateAccountAPI[T,P]) InitializeWallet(ctx context.Context, url string) (string, error) {
// 	wallet, err := s.am.Wallet(url)
// 	if err != nil {
// 		return "", err
// 	}

// 	entropy, err := bip39.NewEntropy(256)
// 	if err != nil {
// 		return "", err
// 	}

// 	mnemonic, err := bip39.NewMnemonic(entropy)
// 	if err != nil {
// 		return "", err
// 	}

// 	seed := bip39.NewSeed(mnemonic, "")

// 	switch wallet := wallet.(type) {
// 	case *scwallet.Wallet:
// 		return mnemonic, wallet.Initialize(seed)
// 	default:
// 		return "", fmt.Errorf("specified wallet does not support initialization")
// 	}
// }

// Unpair deletes a pairing between wallet and geth.
// func (s *PrivateAccountAPI[T,P]) Unpair(ctx context.Context, url string, pin string) error {
// 	wallet, err := s.am.Wallet(url)
// 	if err != nil {
// 		return err
// 	}

// 	switch wallet := wallet.(type) {
// 	case *scwallet.Wallet:
// 		return wallet.Unpair([]byte(pin))
// 	default:
// 		return fmt.Errorf("specified wallet does not support pairing")
// 	}
// }

// PublicBlockChainAPI provides an API to access the Ethereum blockchain.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicBlockChainAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	b Backend[T,P]
}

// NewPublicBlockChainAPI creates a new Ethereum blockchain API.
func NewPublicBlockChainAPI[T crypto.PrivateKey, P crypto.PublicKey](b Backend[T,P]) *PublicBlockChainAPI[T,P] {
	return &PublicBlockChainAPI[T,P]{b}
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (api *PublicBlockChainAPI[T,P]) ChainId() (*hexutil.Big, error) {
	// if current block is at or past the EIP-155 replay-protection fork block, return chainID from config
	if config := api.b.ChainConfig(); config.IsEIP155(api.b.CurrentBlock().Number()) {
		return (*hexutil.Big)(config.ChainID), nil
	}
	return nil, fmt.Errorf("chain not synced beyond EIP-155 replay-protection fork block")
}

// GetPSI - retunrs the PSI that was resolved based on the client request
func (s *PublicBlockChainAPI[T,P]) GetPSI(ctx context.Context) (string, error) {
	psm, err := s.b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return "", err
	}
	return psm.ID.String(), nil
}

// BlockNumber returns the block number of the chain head.
func (s *PublicBlockChainAPI[T,P]) BlockNumber() hexutil.Uint64 {
	header, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber) // latest header should always be available
	return hexutil.Uint64(header.Number.Uint64())
}

// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (s *PublicBlockChainAPI[T,P]) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	return (*hexutil.Big)(state.GetBalance(address)), state.Error()
}

// Result structs for GetProof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}
type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// GetProof returns the Merkle-proof for a given account and optionally some storage keys.
func (s *PublicBlockChainAPI[T,P]) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}

	storageTrie := state.StorageTrie(address)
	storageHash := types.EmptyRootHash
	codeHash := state.GetCodeHash(address)
	storageProof := make([]StorageResult, len(storageKeys))

	// if we have a storageTrie, (which means the account exists), we can update the storagehash
	if storageTrie != nil {
		storageHash = storageTrie.Hash()
	} else {
		// no storageTrie means the account does not exist, so the codeHash is the hash of an empty bytearray.
		codeHash = crypto.Keccak256Hash[P](nil)
	}

	// create the proof for the storageKeys
	for i, key := range storageKeys {
		if storageTrie != nil {
			proof, storageError := state.GetStorageProof(address, common.HexToHash(key))
			if storageError != nil {
				return nil, storageError
			}
			storageProof[i] = StorageResult{key, (*hexutil.Big)(state.GetState(address, common.HexToHash(key)).Big()), toHexSlice(proof)}
		} else {
			storageProof[i] = StorageResult{key, &hexutil.Big{}, []string{}}
		}
	}

	// create the accountProof
	accountProof, proofErr := state.GetProof(address)
	if proofErr != nil {
		return nil, proofErr
	}

	return &AccountResult{
		Address:      address,
		AccountProof: toHexSlice(accountProof),
		Balance:      (*hexutil.Big)(state.GetBalance(address)),
		CodeHash:     codeHash,
		Nonce:        hexutil.Uint64(state.GetNonce(address)),
		StorageHash:  storageHash,
		StorageProof: storageProof,
	}, state.Error()
}

// GetHeaderByNumber returns the requested canonical block header.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
func (s *PublicBlockChainAPI[T,P]) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	header, err := s.b.HeaderByNumber(ctx, number)
	if header != nil && err == nil {
		response := s.rpcMarshalHeader(ctx, header)
		if number == rpc.PendingBlockNumber {
			// Pending header need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetHeaderByHash returns the requested header by hash.
func (s *PublicBlockChainAPI[T,P]) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	header, _ := s.b.HeaderByHash(ctx, hash)
	if header != nil {
		return s.rpcMarshalHeader(ctx, header)
	}
	return nil
}

// GetBlockByNumber returns the requested canonical block.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
// * When fullTx is true all transactions in the block are returned, otherwise
//   only the transaction hash is returned.
func (s *PublicBlockChainAPI[T,P]) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, number)
	if block != nil && err == nil {
		response, err := s.rpcMarshalBlock(ctx, block, true, fullTx)
		if err == nil && number == rpc.PendingBlockNumber {
			// Pending blocks need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI[T,P]) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, hash)
	if block != nil {
		return s.rpcMarshalBlock(ctx, block, true, fullTx)
	}
	return nil, err
}

// GetUncleByBlockNumberAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI[T,P]) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", blockNr, "hash", block.Hash(), "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader[P](uncles[index])
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}

// GetUncleByBlockHashAndIndex returns the uncle block for the given block hash and index. When fullTx is true
// all transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (s *PublicBlockChainAPI[T,P]) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, blockHash)
	if block != nil {
		uncles := block.Uncles()
		if index >= hexutil.Uint(len(uncles)) {
			log.Debug("Requested uncle not found", "number", block.Number(), "hash", blockHash, "index", index)
			return nil, nil
		}
		block = types.NewBlockWithHeader[P](uncles[index])
		return s.rpcMarshalBlock(ctx, block, false, false)
	}
	return nil, err
}

// GetUncleCountByBlockNumber returns number of uncles in the block for the given block number
func (s *PublicBlockChainAPI[T,P]) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetUncleCountByBlockHash returns number of uncles in the block for the given block hash
func (s *PublicBlockChainAPI[T,P]) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Uncles()))
		return &n
	}
	return nil
}

// GetCode returns the code stored at the given address in the state for the given block number.
func (s *PublicBlockChainAPI[T,P]) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers are also allowed.
func (s *PublicBlockChainAPI[T,P]) GetStorageAt(ctx context.Context, address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	res := state.GetState(address, common.HexToHash(key))
	return res[:], state.Error()
}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From       *common.Address   `json:"from"`
	To         *common.Address   `json:"to"`
	Gas        *hexutil.Uint64   `json:"gas"`
	GasPrice   *hexutil.Big      `json:"gasPrice"`
	Value      *hexutil.Big      `json:"value"`
	Data       *hexutil.Bytes    `json:"data"`
	AccessList *types.AccessList `json:"accessList"`
}

// ToMessage converts CallArgs to the Message type used by the core evm
func (args *CallArgs) ToMessage(globalGasCap uint64) types.Message {
	// Set sender address or use zero address if none specified.
	var addr common.Address
	if args.From != nil {
		addr = *args.From
	}

	// Set default gas & gas price if none were set
	gas := globalGasCap
	if gas == 0 {
		gas = uint64(math.MaxUint64 / 2)
	}
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != 0 && globalGasCap < gas {
		log.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap
	}
	gasPrice := new(big.Int)
	if args.GasPrice != nil {
		gasPrice = args.GasPrice.ToInt()
	}
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}
	var data []byte
	if args.Data != nil {
		data = *args.Data
	}
	var accessList types.AccessList
	if args.AccessList != nil {
		accessList = *args.AccessList
	}

	msg := types.NewMessage(addr, args.To, 0, value, gas, gasPrice, data, accessList, false)
	return msg
}

// OverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// StateOverride is the collection of overridden accounts.
type StateOverride [P crypto.PublicKey] map[common.Address]OverrideAccount

// Apply overrides the fields of specified accounts into the given state.
func (diff *StateOverride[P]) Apply(state *state.StateDB[P]) error {
	if diff == nil {
		return nil
	}
	for addr, account := range *diff {
		// Override account nonce.
		if account.Nonce != nil {
			state.SetNonce(addr, uint64(*account.Nonce))
		}
		// Override account(contract) code.
		if account.Code != nil {
			state.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil {
			state.SetBalance(addr, (*big.Int)(*account.Balance))
		}
		if account.State != nil && account.StateDiff != nil {
			return fmt.Errorf("account %s has both 'state' and 'stateDiff'", addr.Hex())
		}
		// Replace entire state if caller requires.
		if account.State != nil {
			state.SetStorage(addr, *account.State)
		}
		// Apply state diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				state.SetState(addr, key, value)
			}
		}
	}
	return nil
}

// Quorum - Multitenancy
// Before returning the result, we need to inspect the EVM and
// perform verification check
func DoCall[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], args CallArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride[P], vmCfg vm.Config[P], timeout time.Duration, globalGasCap uint64) (*core.ExecutionResult, error) {
	defer func(start time.Time) { log.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	/*if err := overrides.Apply(state.(eth.EthAPIState)); err != nil {
		return nil, err
	}*/
	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	// Get a new instance of the EVM.
	msg := args.ToMessage(globalGasCap)
	evm, vmError, err := b.GetEVM(ctx, msg, state, header, nil)
	if err != nil {
		return nil, err
	}
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel()
	}()

	// Execute the message.
	gp := new(core.GasPool).AddGas(math.MaxUint64)
	result, applyErr := core.ApplyMessage(evm, msg, gp)
	if err := vmError(); err != nil {
		return nil, err
	}

	// If the timer caused an abort, return an appropriate error message
	if evm.Cancelled() {
		return nil, fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}
	if applyErr != nil {
		return result, fmt.Errorf("err: %w (supplied gas %d)", applyErr, msg.Gas())
	}
	return result, nil
}

func newRevertError[P crypto.PublicKey](result *core.ExecutionResult) *revertError {
	reason, errUnpack := abi.UnpackRevert[P](result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// revertError is an API error that encompassas an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}

// Call executes the given transaction on the state for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the state/blockchain and is
// useful to execute and retrieve values.
// Quorum
// - replaced the default 5s time out with the value passed in vm.calltimeout
// - multi tenancy verification
func (s *PublicBlockChainAPI[T,P]) Call(ctx context.Context, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride[P]) (hexutil.Bytes, error) {
	var accounts map[common.Address]OverrideAccount
	if overrides != nil {
		accounts = *overrides
	}
	stateOverride := StateOverride[P](accounts)
	result, err := DoCall(ctx, s.b, args, blockNrOrHash, &stateOverride, vm.Config[P]{}, s.b.CallTimeOut(), s.b.RPCGasCap())
	if err != nil {
		return nil, err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(result.Revert()) > 0 {
		return nil, newRevertError[P](result)
	}
	return result.Return(), result.Err
}

func DoEstimateGas[T crypto.PrivateKey, P crypto.PublicKey] (ctx context.Context, b Backend[T,P], args CallArgs, blockNrOrHash rpc.BlockNumberOrHash, gasCap uint64) (hexutil.Uint64, error) {
	// Binary search the gas requirement, as it may be higher than the amount used
	var (
		lo  uint64 = params.TxGas - 1
		hi  uint64
		cap uint64
	)
	// Use zero address if sender unspecified.
	if args.From == nil {
		args.From = new(common.Address)
	}
	// Determine the highest gas limit can be used during the estimation.
	if args.Gas != nil && uint64(*args.Gas) >= params.TxGas {
		hi = uint64(*args.Gas)
	} else {
		// Retrieve the block to act as the gas ceiling
		block, err := b.BlockByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return 0, err
		}
		if block == nil {
			return 0, errors.New("block not found")
		}
		hi = block.GasLimit()
	}
	// Recap the highest gas limit with account's available balance.
	if args.GasPrice != nil && args.GasPrice.ToInt().BitLen() != 0 {
		state, _, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
		if err != nil {
			return 0, err
		}
		balance := state.GetBalance(*args.From) // from can't be nil
		available := new(big.Int).Set(balance)
		if args.Value != nil {
			if args.Value.ToInt().Cmp(available) >= 0 {
				return 0, errors.New("insufficient funds for transfer")
			}
			available.Sub(available, args.Value.ToInt())
		}
		allowance := new(big.Int).Div(available, args.GasPrice.ToInt())

		// If the allowance is larger than maximum uint64, skip checking
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := args.Value
			if transfer == nil {
				transfer = new(hexutil.Big)
			}
			log.Warn("Gas estimation capped by limited funds", "original", hi, "balance", balance,
				"sent", transfer.ToInt(), "gasprice", args.GasPrice.ToInt(), "fundable", allowance)
			hi = allowance.Uint64()
		}
	}
	// Recap the highest gas allowance with specified gascap.
	if gasCap != 0 && hi > gasCap {
		log.Warn("Caller gas above allowance, capping", "requested", hi, "cap", gasCap)
		hi = gasCap
	}
	cap = hi

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) (bool, *core.ExecutionResult, error) {
		args.Gas = (*hexutil.Uint64)(&gas)

		result, err := DoCall(ctx, b, args, blockNrOrHash, nil, vm.Config[P]{}, 0, gasCap)
		if err != nil {
			if errors.Is(err, core.ErrIntrinsicGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			return true, nil, err // Bail out
		}
		return result.Failed(), result, nil
	}
	// Execute the binary search and hone in on an executable gas limit
	for lo+1 < hi {
		mid := (hi + lo) / 2
		failed, _, err := executable(mid)

		// If the error is not nil(consensus error), it means the provided message
		// call or transaction will never be accepted no matter how much gas it is
		// assigned. Return the error directly, don't struggle any more.
		if err != nil {
			return 0, err
		}
		if failed {
			lo = mid
		} else {
			hi = mid
		}
	}
	// Reject the transaction as invalid if it still fails at the highest allowance
	if hi == cap {
		failed, result, err := executable(hi)
		if err != nil {
			return 0, err
		}
		if failed {
			if result != nil && result.Err != vm.ErrOutOfGas {
				if len(result.Revert()) > 0 {
					return 0, newRevertError[P](result)
				}
				return 0, result.Err
			}
			// Otherwise, the specified gas cap is too low
			return 0, fmt.Errorf("gas required exceeds allowance (%d)", cap)
		}
	}

	//QUORUM

	//We don't know if this is going to be a private or public transaction
	//It is possible to have a data field that has a lower intrinsic value than the PTM hash
	//so this checks that if we were to place a PTM hash (with all non-zero values) here then the transaction would
	//still run
	//This makes the return value a potential over-estimate of gas, rather than the exact cost to run right now

	//if the transaction has a value then it cannot be private, so we can skip this check
	if args.Value != nil && args.Value.ToInt().Cmp(big.NewInt(0)) == 0 {
		currentBlockHeight := b.CurrentHeader().Number
		homestead := b.ChainConfig().IsHomestead(currentBlockHeight)
		istanbul := b.ChainConfig().IsIstanbul(currentBlockHeight)

		var data []byte
		if args.Data == nil {
			data = nil
		} else {
			data = []byte(*args.Data)
		}
		var accessList types.AccessList
		if args.AccessList != nil {
			accessList = *args.AccessList
		}
		intrinsicGasPublic, err := core.IntrinsicGas(data, accessList, args.To == nil, homestead, istanbul)
		if err != nil {
			return 0, err
		}
		intrinsicGasPrivate, err := core.IntrinsicGas(common.Hex2Bytes(common.MaxPrivateIntrinsicDataHex), accessList, args.To == nil, homestead, istanbul)
		if err != nil {
			return 0, err
		}

		if intrinsicGasPrivate > intrinsicGasPublic {
			if math.MaxUint64-hi < intrinsicGasPrivate-intrinsicGasPublic {
				return 0, fmt.Errorf("private intrinsic gas addition exceeds allowance")
			}
			return hexutil.Uint64(hi + (intrinsicGasPrivate - intrinsicGasPublic)), nil
		}

	}

	//END QUORUM

	return hexutil.Uint64(hi), nil
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (s *PublicBlockChainAPI[T,P]) EstimateGas(ctx context.Context, args CallArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	return DoEstimateGas(ctx, s.b, args, bNrOrHash, s.b.RPCGasCap())
}

// ExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
// execution status, the amount of gas used and the return value
type ExecutionResult struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// StructLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type StructLogRes struct {
	Pc      uint64             `json:"pc"`
	Op      string             `json:"op"`
	Gas     uint64             `json:"gas"`
	GasCost uint64             `json:"gasCost"`
	Depth   int                `json:"depth"`
	Error   error              `json:"error,omitempty"`
	Stack   *[]string          `json:"stack,omitempty"`
	Memory  *[]string          `json:"memory,omitempty"`
	Storage *map[string]string `json:"storage,omitempty"`
}

// FormatLogs formats EVM returned structured logs for json output
func FormatLogs(logs []vm.StructLog) []StructLogRes {
	formatted := make([]StructLogRes, len(logs))
	for index, trace := range logs {
		formatted[index] = StructLogRes{
			Pc:      trace.Pc,
			Op:      trace.Op.String(),
			Gas:     trace.Gas,
			GasCost: trace.GasCost,
			Depth:   trace.Depth,
			Error:   trace.Err,
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = fmt.Sprintf("%x", math.PaddedBigBytes(stackValue, 32))
			}
			formatted[index].Stack = &stack
		}
		if trace.Memory != nil {
			memory := make([]string, 0, (len(trace.Memory)+31)/32)
			for i := 0; i+32 <= len(trace.Memory); i += 32 {
				memory = append(memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
			}
			formatted[index].Memory = &memory
		}
		if trace.Storage != nil {
			storage := make(map[string]string)
			for i, storageValue := range trace.Storage {
				storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
			}
			formatted[index].Storage = &storage
		}
	}
	return formatted
}

// RPCMarshalHeader converts the given header to the RPC output .
func RPCMarshalHeader[P crypto.PublicKey](head *types.Header[P]) map[string]interface{} {
	return map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             head.Hash(),
		"parentHash":       head.ParentHash,
		"nonce":            head.Nonce,
		"mixHash":          head.MixDigest,
		"sha3Uncles":       head.UncleHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"miner":            head.Coinbase,
		"difficulty":       (*hexutil.Big)(head.Difficulty),
		"extraData":        hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(head.Size()),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}
}

// RPCMarshalBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func RPCMarshalBlock[P crypto.PublicKey](block *types.Block[P], inclTx bool, fullTx bool) (map[string]interface{}, error) {
	fields := RPCMarshalHeader(block.Header())
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(tx *types.Transaction[P]) (interface{}, error) {
			return tx.Hash(), nil
		}
		if fullTx {
			formatTx = func(tx *types.Transaction[P]) (interface{}, error) {
				return newRPCTransactionFromBlockHash(block, tx.Hash()), nil
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range txs {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}
	uncles := block.Uncles()
	uncleHashes := make([]common.Hash, len(uncles))
	for i, uncle := range uncles {
		uncleHashes[i] = uncle.Hash()
	}
	fields["uncles"] = uncleHashes

	return fields, nil
}

// rpcMarshalHeader uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicBlockchainAPI`.
func (s *PublicBlockChainAPI[T,P]) rpcMarshalHeader(ctx context.Context, header *types.Header[P]) map[string]interface{} {
	fields := RPCMarshalHeader(header)
	fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(ctx, header.Hash()))
	return fields
}

// rpcMarshalBlock uses the generalized output filler, then adds the total difficulty field, which requires
// a `PublicBlockchainAPI`.
func (s *PublicBlockChainAPI[T,P]) rpcMarshalBlock(ctx context.Context, b *types.Block[P], inclTx bool, fullTx bool) (map[string]interface{}, error) {
	fields, err := RPCMarshalBlock(b, inclTx, fullTx)
	if err != nil {
		return nil, err
	}
	if inclTx {
		fields["totalDifficulty"] = (*hexutil.Big)(s.b.GetTd(ctx, b.Hash()))
	}
	return fields, err
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash      `json:"blockHash"`
	BlockNumber      *hexutil.Big      `json:"blockNumber"`
	From             common.Address    `json:"from"`
	Gas              hexutil.Uint64    `json:"gas"`
	GasPrice         *hexutil.Big      `json:"gasPrice"`
	Hash             common.Hash       `json:"hash"`
	Input            hexutil.Bytes     `json:"input"`
	Nonce            hexutil.Uint64    `json:"nonce"`
	To               *common.Address   `json:"to"`
	TransactionIndex *hexutil.Uint64   `json:"transactionIndex"`
	Value            *hexutil.Big      `json:"value"`
	Type             hexutil.Uint64    `json:"type"`
	Accesses         *types.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big      `json:"chainId,omitempty"`
	V                *hexutil.Big      `json:"v"`
	R                *hexutil.Big      `json:"r"`
	S                *hexutil.Big      `json:"s"`
}

// newRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func newRPCTransaction[P crypto.PublicKey](tx *types.Transaction[P], blockHash common.Hash, blockNumber uint64, index uint64) *RPCTransaction {
	// Determine the signer. For replay-protected transactions, use the most permissive
	// signer, because we assume that signers are backwards-compatible with old
	// transactions. For non-protected transactions, the homestead signer signer is used
	// because the return value of ChainId is zero for those transactions.
	var signer types.Signer[P]
	if tx.Protected() && !tx.IsPrivate() {
		signer = types.LatestSignerForChainID[P](tx.ChainId())
	} else {
		signer = types.HomesteadSigner[P]{}
	}

	from, _ := types.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	result := &RPCTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	if tx.Type() == types.AccessListTxType {
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	}
	return result
}

// newRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func newRPCPendingTransaction[P crypto.PublicKey](tx *types.Transaction[P]) *RPCTransaction {
	return newRPCTransaction(tx, common.Hash{}, 0, 0)
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockIndex[P crypto.PublicKey](b *types.Block[P], index uint64) *RPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return newRPCTransaction(txs[index], b.Hash(), b.NumberU64(), index)
}

// newRPCRawTransactionFromBlockIndex returns the bytes of a transaction given a block and a transaction index.
func newRPCRawTransactionFromBlockIndex[P crypto.PublicKey](b *types.Block[P], index uint64) hexutil.Bytes {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	blob, _ := txs[index].MarshalBinary()
	return blob
}

// newRPCTransactionFromBlockHash returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockHash[P crypto.PublicKey](b *types.Block[P], hash common.Hash) *RPCTransaction {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newRPCTransactionFromBlockIndex(b, uint64(idx))
		}
	}
	return nil
}

// accessListResult returns an optional accesslist
// Its the result of the `debug_createAccessList` RPC call.
// It contains an error if the transaction itself failed.
type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

// CreateAccessList creates a EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.
func (s *PublicBlockChainAPI[T,P]) CreateAccessList(ctx context.Context, args SendTxArgs[T,P], blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	bNrOrHash := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	acl, gasUsed, vmerr, err := AccessList(ctx, s.b, bNrOrHash, args)
	if err != nil {
		return nil, err
	}
	result := &accessListResult{Accesslist: &acl, GasUsed: hexutil.Uint64(gasUsed)}
	if vmerr != nil {
		result.Error = vmerr.Error()
	}
	return result, nil
}

// AccessList creates an access list for the given transaction.
// If the accesslist creation fails an error is returned.
// If the transaction itself fails, an vmErr is returned.
func AccessList[T crypto.PrivateKey, P crypto.PublicKey] (ctx context.Context, b Backend[T,P], blockNrOrHash rpc.BlockNumberOrHash, args SendTxArgs[T,P]) (acl types.AccessList, gasUsed uint64, vmErr error, err error) {
	// Retrieve the execution context
	db, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if db == nil || err != nil {
		return nil, 0, nil, err
	}
	// If the gas amount is not set, extract this as it will depend on access
	// lists and we'll need to reestimate every time
	nogas := args.Gas == nil

	// Ensure any missing fields are filled, extract the recipient and input data
	if err := args.setDefaults(ctx, b); err != nil {
		return nil, 0, nil, err
	}
	var to common.Address
	if args.To != nil {
		to = *args.To
	} else {
		to = crypto.CreateAddress[P](args.From, uint64(*args.Nonce))
	}
	var input []byte
	if args.Input != nil {
		input = *args.Input
	} else if args.Data != nil {
		input = *args.Data
	}
	// Retrieve the precompiles since they don't need to be added to the access list
	precompiles := vm.ActivePrecompiles(b.ChainConfig().Rules(header.Number))

	// Create an initial tracer
	prevTracer := vm.NewAccessListTracer[P](nil, args.From, to, precompiles)
	if args.AccessList != nil {
		prevTracer = vm.NewAccessListTracer[P](*args.AccessList, args.From, to, precompiles)
	}
	for {
		// Retrieve the current access list to expand
		accessList := prevTracer.AccessList()
		log.Trace("Creating access list", "input", accessList)

		// If no gas amount was specified, each unique access list needs it's own
		// gas calculation. This is quite expensive, but we need to be accurate
		// and it's convered by the sender only anyway.
		if nogas {
			args.Gas = nil
			if err := args.setDefaults(ctx, b); err != nil {
				return nil, 0, nil, err // shouldn't happen, just in case
			}
		}
		// Copy the original db so we don't modify it
		// statedb := db.Copy()
		statedb := db.(*state.StateDB[P]).Copy()
		msg := types.NewMessage(args.From, args.To, uint64(*args.Nonce), args.Value.ToInt(), uint64(*args.Gas), args.GasPrice.ToInt(), input, accessList, false)

		// Apply the transaction with the access list tracer
		tracer := vm.NewAccessListTracer[P](accessList, args.From, to, precompiles)
		config := vm.Config[P]{Tracer: tracer, Debug: true}
		vmenv, _, err := b.GetEVM(ctx, msg, statedb, header, &config)
		if err != nil {
			return nil, 0, nil, err
		}
		res, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(msg.Gas()))
		if err != nil {
			return nil, 0, nil, fmt.Errorf("failed to apply transaction: %v err: %v", args.toTransaction().Hash(), err)
		}
		if tracer.Equal(prevTracer) {
			return accessList, res.UsedGas, res.Err, nil
		}
		prevTracer = tracer
	}
}

// PublicTransactionPoolAPI exposes methods for the RPC interface
type PublicTransactionPoolAPI [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	b         Backend[T,P]
	nonceLock *AddrLocker
	signer    types.Signer[P]
}

// NewPublicTransactionPoolAPI creates a new RPC service with methods specific for the transaction pool.
func NewPublicTransactionPoolAPI[T crypto.PrivateKey, P crypto.PublicKey] (b Backend[T,P], nonceLock *AddrLocker) *PublicTransactionPoolAPI[T,P] {
	// The signer used by the API should always be the 'latest' known one because we expect
	// signers to be backwards-compatible with old transactions.
	signer := types.LatestSigner[P](b.ChainConfig())
	return &PublicTransactionPoolAPI[T,P]{b, nonceLock, signer}
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (s *PublicTransactionPoolAPI[T,P]) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (s *PublicTransactionPoolAPI[T,P]) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(len(block.Transactions()))
		return &n
	}
	return nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI[T,P]) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI[T,P]) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *RPCTransaction {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (s *PublicTransactionPoolAPI[T,P]) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByNumber(ctx, blockNr); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (s *PublicTransactionPoolAPI[T,P]) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	if block, _ := s.b.BlockByHash(ctx, blockHash); block != nil {
		return newRPCRawTransactionFromBlockIndex(block, uint64(index))
	}
	return nil
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number
func (s *PublicTransactionPoolAPI[T,P]) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	// Ask transaction pool for the nonce which includes pending transactions
	if blockNr, ok := blockNrOrHash.Number(); ok && blockNr == rpc.PendingBlockNumber {
		nonce, err := s.b.GetPoolNonce(ctx, address)
		if err != nil {
			return nil, err
		}
		return (*hexutil.Uint64)(&nonce), nil
	}
	// Resolve block number and use its state to ask for the nonce
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	nonce := state.GetNonce(address)
	return (*hexutil.Uint64)(&nonce), state.Error()
}

// Quorum

type PrivacyMetadataWithMandatoryRecipients struct {
	*state.PrivacyMetadata
	MandatoryRecipients []string `json:"mandatoryFor,omitempty"`
}

func (s *PublicTransactionPoolAPI[T,P]) GetContractPrivacyMetadata(ctx context.Context, address common.Address) (*PrivacyMetadataWithMandatoryRecipients, error) {
	state, _, err := s.b.StateAndHeaderByNumber(ctx, rpc.LatestBlockNumber)
	if state == nil || err != nil {
		return nil, err
	}
	var mandatoryRecipients []string

	privacyMetadata, err := state.GetPrivacyMetadata(address)
	if privacyMetadata == nil || err != nil {
		return nil, err
	}

	if privacyMetadata.PrivacyFlag == engine.PrivacyFlagMandatoryRecipients {
		mandatoryRecipients, err = private.Ptm.GetMandatory(privacyMetadata.CreationTxHash)
		if len(mandatoryRecipients) == 0 || err != nil {
			return nil, err
		}
	}

	return &PrivacyMetadataWithMandatoryRecipients{privacyMetadata, mandatoryRecipients}, nil
}

// End Quorum

// GetTransactionByHash returns the transaction for the given hash
func (s *PublicTransactionPoolAPI[T,P]) GetTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	// Try to return an already finalized transaction
	tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return newRPCTransaction(tx, blockHash, blockNumber, index), nil
	}
	// No finalized transaction, try to retrieve it from the pool
	if tx := s.b.GetPoolTransaction(hash); tx != nil {
		return newRPCPendingTransaction(tx), nil
	}

	// Transaction unknown, return as such
	return nil, nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (s *PublicTransactionPoolAPI[T,P]) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	// Retrieve a finalized transaction, or a pooled otherwise
	tx, _, _, _, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		if tx = s.b.GetPoolTransaction(hash); tx == nil {
			// Transaction not found anywhere, abort
			return nil, nil
		}
	}
	// Serialize to RLP and return
	return tx.MarshalBinary()
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (s *PublicTransactionPoolAPI[T,P]) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	tx, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, nil
	}
	receipts, err := s.b.GetReceipts(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	if len(receipts) <= int(index) {
		return nil, nil
	}
	receipt := receipts[index]

	// Quorum: note that upstream code has been refactored into this method
	return getTransactionReceiptCommonCode(tx, blockHash, blockNumber, hash, index, receipt)
}

// Quorum
// Common code extracted from GetTransactionReceipt() to enable reuse
func getTransactionReceiptCommonCode[P crypto.PublicKey](tx *types.Transaction[P], blockHash common.Hash, blockNumber uint64, hash common.Hash, index uint64, receipt *types.Receipt[P]) (map[string]interface{}, error) {
	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   hash,
		"transactionIndex":  hexutil.Uint64(index),
		"from":              tx.From(),
		"to":                tx.To(),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(tx.Type()),
		// Quorum
		"isPrivacyMarkerTransaction": tx.IsPrivacyMarker(),
	}

	// Quorum
	if len(receipt.RevertReason) > 0 {
		fields["revertReason"] = hexutil.Encode(receipt.RevertReason)
	}
	// End Quorum

	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	}
	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields, nil
}

// Quorum
// GetPrivateTransactionByHash accepts the hash for a privacy marker transaction,
// but returns the associated private transaction
func (s *PublicTransactionPoolAPI[T,P]) GetPrivateTransactionByHash(ctx context.Context, hash common.Hash) (*RPCTransaction, error) {
	if !private.IsQuorumPrivacyEnabled() {
		return nil, fmt.Errorf("PrivateTransactionManager is not enabled")
	}
	psm, err := s.b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}

	// first need the privacy marker transaction
	pmt, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}

	// now retrieve the private transaction
	if pmt != nil {
		tx, managedParties, _, err := private.FetchPrivateTransaction[P](pmt.Data())
		if err != nil {
			return nil, err
		}
		if tx != nil && !s.b.PSMR().NotIncludeAny(psm, managedParties...) {
			return newRPCTransaction(tx, blockHash, blockNumber, index), nil
		}
	}

	// Transaction unknown or not a participant in the private transaction, return as such
	return nil, nil
}

// Quorum
// GetPrivateTransactionReceipt accepts the hash for a privacy marker transaction,
// but returns the receipt of the associated private transaction
func (s *PublicTransactionPoolAPI[T,P]) GetPrivateTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	// first need the privacy marker transaction
	pmt, blockHash, blockNumber, index, err := s.b.GetTransaction(ctx, hash)
	if err != nil {
		return nil, err
	}
	if pmt == nil {
		// Transaction unknown, return as such
		return nil, errors.New("privacy marker transaction not found")
	}

	// now retrieve the private transaction
	tx, _, _, err := private.FetchPrivateTransaction[P](pmt.Data())
	if err != nil {
		return nil, err
	}
	// Transaction not found, or not a participant in the private transaction, return as such
	if tx == nil {
		return nil, errors.New("private transaction not found for this participant")
	}

	// get receipt for the privacy marker transaction
	receipts, err := s.b.GetReceipts(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	if len(receipts) <= int(index) {
		return nil, errors.New("could not find receipt for private transaction")
	}
	pmtReceipt := receipts[index]

	// now extract the receipt for the private transaction
	psm, err := s.b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}
	receipt := pmtReceipt.PSReceipts[psm.ID]
	if receipt == nil {
		return nil, errors.New("could not find receipt for private transaction")
	}

	return getTransactionReceiptCommonCode(tx, blockHash, blockNumber, hash, index, receipt)
}

// Quorum: if signing a private TX, set with tx.SetPrivate() before calling this method.
// sign is a helper function that signs a transaction with the private key of the given address.
func (s *PublicTransactionPoolAPI[T,P]) sign(addr common.Address, tx *types.Transaction[P]) (*types.Transaction[P], error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}

	// Quorum
	var chainID *big.Int
	if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) && !tx.IsPrivate() {
		chainID = config.ChainID
	}
	// /Quorum

	// Request the wallet to sign the transaction
	return wallet.SignTx(account, tx, chainID)
}

// SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
// Quorum: introducing additional arguments encapsulated in PrivateTxArgs struct
//		   to support private transactions processing.
type SendTxArgs [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	PrivateTxArgs[T,P] // Quorum

	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// For non-legacy transactions
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

func (s SendTxArgs[T,P]) IsPrivate() bool {
	return s.PrivateFor != nil
}

// SendRawTxArgs represents the arguments to submit a new signed private transaction into the transaction pool.
type SendRawTxArgs [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	PrivateTxArgs[T,P]
}

// Additional arguments used in private transactions
type PrivateTxArgs [T crypto.PrivateKey, P crypto.PublicKey]  struct {
	// PrivateFrom is the public key of the sending party.
	// The public key must be available in the Private Transaction Manager (i.e.: Tessera) which is paired with this geth node.
	// Empty value means the Private Transaction Manager will use the first public key
	// in its list of available keys which it maintains.
	PrivateFrom string `json:"privateFrom"`
	// PrivateFor is the list of public keys which are available in the Private Transaction Managers in the network.
	// The transaction payload is only visible to those party to the transaction.
	PrivateFor          []string               `json:"privateFor"`
	PrivateTxType       string                 `json:"restriction"`
	PrivacyFlag         engine.PrivacyFlagType `json:"privacyFlag"`
	MandatoryRecipients []string               `json:"mandatoryFor"`
}

func (args *PrivateTxArgs[T,P]) SetDefaultPrivateFrom(ctx context.Context, b Backend[T,P]) error {
	if args.PrivateFor != nil && len(args.PrivateFrom) == 0 && b.ChainConfig().IsMPS {
		psm, err := b.PSMR().ResolveForUserContext(ctx)
		if err != nil {
			return err
		}
		args.PrivateFrom = psm.Addresses[0]
	}
	return nil
}

func (args *PrivateTxArgs[T,P]) SetRawTransactionPrivateFrom(ctx context.Context, b Backend[T,P], tx *types.Transaction[P]) error {
	if args.PrivateFor != nil && b.ChainConfig().IsMPS {
		hash := common.BytesToEncryptedPayloadHash(tx.Data())
		_, retrievedPrivateFrom, _, err := private.Ptm.ReceiveRaw(hash)
		if err != nil {
			return err
		}
		if len(args.PrivateFrom) == 0 {
			args.PrivateFrom = retrievedPrivateFrom
		}
		if args.PrivateFrom != retrievedPrivateFrom {
			return fmt.Errorf("The PrivateFrom address retrieved from the privacy manager does not match private PrivateFrom (%s) specified in transaction arguments.", args.PrivateFrom)
		}
		psm, err := b.PSMR().ResolveForUserContext(ctx)
		if err != nil {
			return err
		}
		if psm.NotIncludeAny(args.PrivateFrom) {
			return fmt.Errorf("The PrivateFrom address does not match the specified private state (%s)", psm.ID)
		}
	}
	return nil
}

// setDefaults fills in default values for unspecified tx fields.
func (args *SendTxArgs[T,P]) setDefaults(ctx context.Context, b Backend[T,P]) error {
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Value == nil {
		args.Value = new(hexutil.Big)
	}
	if args.Nonce == nil {
		nonce, err := b.GetPoolNonce(ctx, args.From)
		if err != nil {
			return err
		}
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
	}
	if args.To == nil {
		// Contract creation
		var input []byte
		if args.Data != nil {
			input = *args.Data
		} else if args.Input != nil {
			input = *args.Input
		}
		if len(input) == 0 {
			return errors.New(`contract creation without any data provided`)
		}
	}
	// Estimate the gas usage if necessary.
	if args.Gas == nil {
		// For backwards-compatibility reason, we try both input and data
		// but input is preferred.
		input := args.Input
		if input == nil {
			input = args.Data
		}
		callArgs := CallArgs{
			From:       &args.From, // From shouldn't be nil
			To:         args.To,
			GasPrice:   args.GasPrice,
			Value:      args.Value,
			Data:       input,
			AccessList: args.AccessList,
		}
		pendingBlockNr := rpc.BlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
		estimated, err := DoEstimateGas(ctx, b, callArgs, pendingBlockNr, b.RPCGasCap())
		if err != nil {
			return err
		}
		args.Gas = &estimated
		log.Trace("Estimate gas usage automatically", "gas", args.Gas)
	}
	if args.ChainID == nil {
		id := (*hexutil.Big)(b.ChainConfig().ChainID)
		args.ChainID = id
	}

	// Quorum
	if args.PrivateTxType == "" {
		args.PrivateTxType = "restricted"
	}
	return args.SetDefaultPrivateFrom(ctx, b)
	// End Quorum
}

// toTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *SendTxArgs[T,P]) toTransaction() *types.Transaction[P] {
	var input []byte
	if args.Input != nil {
		input = *args.Input
	} else if args.Data != nil {
		input = *args.Data
	}
	var data types.TxData
	if args.AccessList == nil {
		data = &types.LegacyTx{
			To:       args.To,
			Nonce:    uint64(*args.Nonce),
			Gas:      uint64(*args.Gas),
			GasPrice: (*big.Int)(args.GasPrice),
			Value:    (*big.Int)(args.Value),
			Data:     input,
		}
	} else {
		data = &types.AccessListTx{
			To:         args.To,
			ChainID:    (*big.Int)(args.ChainID),
			Nonce:      uint64(*args.Nonce),
			Gas:        uint64(*args.Gas),
			GasPrice:   (*big.Int)(args.GasPrice),
			Value:      (*big.Int)(args.Value),
			Data:       input,
			AccessList: *args.AccessList,
		}
	}
	return types.NewTx[P](data)
}

// SubmitTransaction is a helper function that submits tx to txPool and logs a message.
func SubmitTransaction[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], tx *types.Transaction[P], privateFrom string, isRaw bool) (common.Hash, error) {
	// If the transaction fee cap is already specified, ensure the
	// fee of the given transaction is _reasonable_.
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), b.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}
	if !b.UnprotectedAllowed() && !tx.Protected() {
		// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
		return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
	}
	// Print a log with full tx details for manual investigations and interventions
	// Quorum
	var signer types.Signer[P]
	if tx.IsPrivate() {
		signer = types.QuorumPrivateTxSigner[P]{}
	} else {
		signer = types.MakeSigner[P](b.ChainConfig(), b.CurrentBlock().Number())
	}
	from, err := types.Sender(signer, tx)
	if err != nil {
		return common.Hash{}, err
	}
	// Quorum
	// Need to do authorization check for Ethereum Account being used in signing.
	// We only care about private transactions (or the private transaction relating to a privacy marker)
	if token, ok := b.SupportsMultitenancy(ctx); ok {
		tx := tx
		// If we are sending a Privacy Marker Transaction, then get the private txn details
		if tx.IsPrivacyMarker() {
			tx, _, _, err = private.FetchPrivateTransaction[P](tx.Data())
			if err != nil {
				return common.Hash{}, err
			}
		}
		innerFrom, err := types.Sender(signer, tx)
		if err != nil {
			return common.Hash{}, err
		}

		if tx.IsPrivate() {
			psm, err := b.PSMR().ResolveForUserContext(ctx)
			if err != nil {
				return common.Hash{}, err
			}
			eoaSecAttr := (&multitenancy.PrivateStateSecurityAttribute{}).WithPSI(psm.ID).WithSelfEOAIf(isRaw, innerFrom)
			psm, err = b.PSMR().ResolveForManagedParty(privateFrom)
			if err != nil {
				return common.Hash{}, err
			}
			privateFromSecAttr := (&multitenancy.PrivateStateSecurityAttribute{}).WithPSI(psm.ID).WithSelfEOAIf(isRaw, innerFrom)
			if isAuthorized, _ := multitenancy.IsAuthorized(token, eoaSecAttr, privateFromSecAttr); !isAuthorized {
				return common.Hash{}, multitenancy.ErrNotAuthorized
			}
		}
	}
	if !b.UnprotectedAllowed() && !tx.Protected() {
		// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
		return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
	}
	if err := b.SendTx(ctx, tx); err != nil {
		return common.Hash{}, err
	}

	if tx.To() == nil {
		addr := crypto.CreateAddress[P](from, tx.Nonce())
		log.Info("Submitted contract creation", "hash", tx.Hash().Hex(), "from", from, "nonce", tx.Nonce(), "contract", addr.Hex(), "value", tx.Value())
		log.EmitCheckpoint(log.TxCreated, "tx", tx.Hash().Hex(), "to", addr.Hex())
	} else {
		log.Info("Submitted transaction", "hash", tx.Hash().Hex(), "from", from, "nonce", tx.Nonce(), "recipient", tx.To(), "value", tx.Value())
		log.EmitCheckpoint(log.TxCreated, "tx", tx.Hash().Hex(), "to", tx.To().Hex())
	}
	return tx.Hash(), nil
}

// runSimulation runs a simulation of the given transaction.
// It returns the EVM instance upon completion
func runSimulation[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], from common.Address, tx *types.Transaction[P]) (*vm.EVM[P], []byte, error) {
	defer func(start time.Time) {
		log.Debug("Simulated Execution EVM call finished", "runtime", time.Since(start))
	}(time.Now())

	// Set sender address or use a default if none specified
	addr := from
	if addr == (common.Address{}) {
		if wallets := b.AccountManager().Wallets(); len(wallets) > 0 {
			if accountList := wallets[0].Accounts(); len(accountList) > 0 {
				addr = accountList[0].Address
			}
		}
	}

	// Create new call message
	msg := types.NewMessage(addr, tx.To(), tx.Nonce(), tx.Value(), tx.Gas(), tx.GasPrice(), tx.Data(), tx.AccessList(), false)

	// Setup context with timeout as gas un-metered
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer func() { cancel() }()

	// Get a new instance of the EVM.
	blockNumber := b.CurrentBlock().Number().Uint64()
	stateAtBlock, header, err := b.StateAndHeaderByNumber(ctx, rpc.BlockNumber(blockNumber))
	if stateAtBlock == nil || err != nil {
		return nil, nil, err
	}
	evm, _, err := b.GetEVM(ctx, msg, stateAtBlock, header, &vm.Config[P]{})
	if err != nil {
		return nil, nil, err
	}

	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel()
	}()

	var contractAddr common.Address

	// even the creation of a contract (init code) can invoke other contracts
	if tx.To() != nil {
		// removed contract availability checks as they are performed in checkAndHandlePrivateTransaction
		data, _, err := evm.Call(vm.AccountRef(addr), *tx.To(), tx.Data(), tx.Gas(), tx.Value())
		return evm, data, err
	} else {
		_, contractAddr, _, err = evm.Create(vm.AccountRef(addr), tx.Data(), tx.Gas(), tx.Value())
		//make sure that nonce is same in simulation as in actual block processing
		//simulation blockNumber will be behind block processing blockNumber by at least 1
		//only guaranteed to work for default config where EIP158=1
		if evm.ChainConfig().IsEIP158(big.NewInt(evm.Context.BlockNumber.Int64() + 1)) {
			evm.StateDB.SetNonce(contractAddr, 1)
		}
	}
	return evm, nil, err
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (s *PublicTransactionPoolAPI[T,P]) SendTransaction(ctx context.Context, args SendTxArgs[T,P]) (common.Hash, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: args.From}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return common.Hash{}, err
	}

	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}

	_, replaceDataWithHash, data, err := checkAndHandlePrivateTransaction(ctx, s.b, args.toTransaction(), &args.PrivateTxArgs, args.From, NormalTransaction)
	if err != nil {
		return common.Hash{}, err
	}
	if replaceDataWithHash {
		// replace the original payload with encrypted payload hash
		args.Data = data.BytesTypeRef()
	}
	// /Quorum

	// Assemble the transaction and sign with the wallet
	tx := args.toTransaction()

	// Quorum
	if args.IsPrivate() {
		tx.SetPrivate()
	}

	var chainID *big.Int
	if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) && !tx.IsPrivate() {
		chainID = config.ChainID
	}
	// /Quorum

	signed, err := wallet.SignTx(account, tx, chainID)
	if err != nil {
		return common.Hash{}, err
	}
	// Quorum
	if signed.IsPrivate() && s.b.IsPrivacyMarkerTransactionCreationEnabled() {
		pmt, err := createPrivacyMarkerTransaction(s.b, signed, &args.PrivateTxArgs)
		if err != nil {
			log.Warn("Failed to create privacy marker transaction for private transaction", "from", args.From, "to", args.To, "value", args.Value.ToInt(), "err", err)
			return common.Hash{}, err
		}

		var pmtChainID *big.Int // PMT is public so will have different chainID used in signing compared to the internal tx
		if config := s.b.ChainConfig(); config.IsEIP155(s.b.CurrentBlock().Number()) {
			pmtChainID = config.ChainID
		}

		signed, err = wallet.SignTx(account, pmt, pmtChainID)
		if err != nil {
			log.Warn("Failed to sign privacy marker transaction for private transaction", "from", args.From, "to", args.To, "value", args.Value.ToInt(), "err", err)
			return common.Hash{}, err
		}
	}

	return SubmitTransaction(ctx, s.b, signed, args.PrivateFrom, false)
}

// FillTransaction fills the defaults (nonce, gas, gasPrice) on a given unsigned transaction,
// and returns it to the caller for further processing (signing + broadcast)
func (s *PublicTransactionPoolAPI[T,P]) FillTransaction(ctx context.Context, args SendTxArgs[T,P]) (*SignTransactionResult[P], error) {
	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and obtain rlp
	// Quorum
	isPrivate, replaceDataWithHash, hash, err := checkAndHandlePrivateTransaction(ctx, s.b, args.toTransaction(), &args.PrivateTxArgs, args.From, FillTransaction)
	if err != nil {
		return nil, err
	}
	if replaceDataWithHash {
		// replace the original payload with encrypted payload hash
		args.Data = hash.BytesTypeRef()
	}
	// /Quorum

	tx := args.toTransaction()

	// Quorum
	if isPrivate {
		tx.SetPrivate()
	}
	// /Quorum
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult[P]{data, tx}, nil
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTransactionPoolAPI[T,P]) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	tx := new(types.Transaction[P])
	if err := tx.UnmarshalBinary(input); err != nil {
		return common.Hash{}, err
	}
	return SubmitTransaction(ctx, s.b, tx, "", true)
}

// Quorum
//
// SendRawPrivateTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (s *PublicTransactionPoolAPI[T,P]) SendRawPrivateTransaction(ctx context.Context, encodedTx hexutil.Bytes, args SendRawTxArgs[T,P]) (common.Hash, error) {

	tx := new(types.Transaction[P])
	if err := tx.UnmarshalBinary(encodedTx); err != nil {
		return common.Hash{}, err
	}

	// Quorum
	if err := args.SetRawTransactionPrivateFrom(ctx, s.b, tx); err != nil {
		return common.Hash{}, err
	}
	isPrivate, _, _, err := checkAndHandlePrivateTransaction(ctx, s.b, tx, &args.PrivateTxArgs, common.Address{}, RawTransaction)
	if err != nil {
		return common.Hash{}, err
	}
	if !isPrivate {
		return common.Hash{}, fmt.Errorf("transaction is not private")
	}

	return SubmitTransaction(ctx, s.b, tx, args.PrivateFrom, true)
}

// DistributePrivateTransaction will perform the simulation checks and send the private transactions data to the other
// private participants
// It then submits the entire private transaction to the attached PTM and sends it to other private participants,
// return the PTM generated hash, intended to be used in the Input field of a Privacy Marker Transaction
func (s *PublicTransactionPoolAPI[T,P]) DistributePrivateTransaction(ctx context.Context, encodedTx hexutil.Bytes, args SendRawTxArgs[T,P]) (string, error) {
	log.Info("distributing raw private tx")

	tx := new(types.Transaction[P])
	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		return "", err
	}

	log.Debug("deserialised raw private tx", "hash", tx.Hash())

	// Quorum
	if err := args.SetRawTransactionPrivateFrom(ctx, s.b, tx); err != nil {
		return "", err
	}
	isPrivate, _, _, err := checkAndHandlePrivateTransaction(ctx, s.b, tx, &args.PrivateTxArgs, common.Address{}, RawTransaction)
	if err != nil {
		return "", err
	}
	if !isPrivate {
		return "", fmt.Errorf("transaction is not private")
	}

	serialisedTx, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}

	_, _, txnHash, err := private.Ptm.Send(serialisedTx, args.PrivateFrom, args.PrivateFor, &engine.ExtraMetadata{})
	if err != nil {
		return "", err
	}
	log.Debug("private transaction sent to PTM", "generated ptm-hash", txnHash)
	return txnHash.Hex(), nil
}

// /Quorum

// Sign calculates an ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (s *PublicTransactionPoolAPI[T,P]) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Sign the requested hash with the wallet
	signature, err := wallet.SignText(account, data)
	if err == nil {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	return signature, err
}

// SignTransactionResult represents a RLP encoded signed transaction.
type SignTransactionResult[P crypto.PublicKey]struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction[P] `json:"tx"`
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (s *PublicTransactionPoolAPI[T,P]) SignTransaction(ctx context.Context, args SendTxArgs[T,P]) (*SignTransactionResult[P], error) {
	if args.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}
	if args.GasPrice == nil {
		return nil, fmt.Errorf("gasPrice not specified")
	}
	if args.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}
	// Quorum
	// setDefaults calls DoEstimateGas in ethereum1.9.0, private transaction is not supported for that feature
	// set gas to constant if nil
	if args.IsPrivate() && args.Gas == nil {
		gas := (hexutil.Uint64)(90000)
		args.Gas = &gas
	}
	// End Quorum
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Before actually sign the transaction, ensure the transaction fee is reasonable.
	if err := checkTxFee(args.GasPrice.ToInt(), uint64(*args.Gas), s.b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	// Quorum
	toSign := args.toTransaction()
	if args.IsPrivate() {
		toSign.SetPrivate()
	}
	// End Quorum

	tx, err := s.sign(args.From, toSign)
	if err != nil {
		return nil, err
	}
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult[P]{data, tx}, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (s *PublicTransactionPoolAPI[T,P]) PendingTransactions() ([]*RPCTransaction, error) {
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	accounts := make(map[common.Address]struct{})
	for _, wallet := range s.b.AccountManager().Wallets() {
		for _, account := range wallet.Accounts() {
			accounts[account.Address] = struct{}{}
		}
	}
	transactions := make([]*RPCTransaction, 0, len(pending))
	for _, tx := range pending {
		from, _ := types.Sender(s.signer, tx)
		if _, exists := accounts[from]; exists {
			transactions = append(transactions, newRPCPendingTransaction(tx))
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (s *PublicTransactionPoolAPI[T,P]) Resend(ctx context.Context, sendArgs SendTxArgs[T,P], gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	if sendArgs.Nonce == nil {
		return common.Hash{}, fmt.Errorf("missing transaction nonce in transaction spec")
	}
	// setDefaults calls DoEstimateGas in ethereum1.9.0, private transaction is not supported for that feature
	// set gas to constant if nil
	if sendArgs.IsPrivate() && sendArgs.Gas == nil {
		gas := (hexutil.Uint64)(90000)
		sendArgs.Gas = &gas
	}
	if err := sendArgs.setDefaults(ctx, s.b); err != nil {
		return common.Hash{}, err
	}
	matchTx := sendArgs.toTransaction()

	// Before replacing the old transaction, ensure the _new_ transaction fee is reasonable.
	var price = matchTx.GasPrice()
	if gasPrice != nil {
		price = gasPrice.ToInt()
	}
	var gas = matchTx.Gas()
	if gasLimit != nil {
		gas = uint64(*gasLimit)
	}
	if err := checkTxFee(price, gas, s.b.RPCTxFeeCap()); err != nil {
		return common.Hash{}, err
	}
	// Iterate the pending list for replacement
	pending, err := s.b.GetPoolTransactions()
	if err != nil {
		return common.Hash{}, err
	}
	for _, p := range pending {
		wantSigHash := s.signer.Hash(matchTx)
		pFrom, err := types.Sender(s.signer, p)
		if err == nil && pFrom == sendArgs.From && s.signer.Hash(p) == wantSigHash {
			// Match. Re-sign and send the transaction.
			if gasPrice != nil && (*big.Int)(gasPrice).Sign() != 0 {
				sendArgs.GasPrice = gasPrice
			}
			if gasLimit != nil && *gasLimit != 0 {
				sendArgs.Gas = gasLimit
			}
			newTx := sendArgs.toTransaction()
			// set v param to 37 to indicate private tx before submitting to the signer.
			if sendArgs.IsPrivate() {
				newTx.SetPrivate()
			}
			signedTx, err := s.sign(sendArgs.From, newTx)
			if err != nil {
				return common.Hash{}, err
			}
			if err = s.b.SendTx(ctx, signedTx); err != nil {
				return common.Hash{}, err
			}
			return signedTx.Hash(), nil
		}
	}
	return common.Hash{}, fmt.Errorf("transaction %#x not found", matchTx.Hash())
}

// PublicDebugAPI is the collection of Ethereum APIs exposed over the public
// debugging endpoint.
type PublicDebugAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	b Backend[T,P]
}

// NewPublicDebugAPI creates a new API definition for the public debug methods
// of the Ethereum service.
func NewPublicDebugAPI[T crypto.PrivateKey, P crypto.PublicKey](b Backend[T,P]) *PublicDebugAPI[T,P] {
	return &PublicDebugAPI[T,P]{b: b}
}

// GetBlockRlp retrieves the RLP encoded for of a single block.
func (api *PublicDebugAPI[T,P]) GetBlockRlp(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	encoded, err := rlp.EncodeToBytes(block)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", encoded), nil
}

// TestSignCliqueBlock fetches the given block number, and attempts to sign it as a clique header with the
// given address, returning the address of the recovered signature
//
// This is a temporary method to debug the externalsigner integration,
// TODO: Remove this method when the integration is mature
func (api *PublicDebugAPI[T,P]) TestSignCliqueBlock(ctx context.Context, address common.Address, number uint64) (common.Address, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return common.Address{}, fmt.Errorf("block #%d not found", number)
	}
	header := block.Header()
	header.Extra = make([]byte, 32+65)
	encoded := clique.CliqueRLP(header)

	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: address}
	wallet, err := api.b.AccountManager().Find(account)
	if err != nil {
		return common.Address{}, err
	}

	signature, err := wallet.SignData(account, accounts.MimetypeClique, encoded)
	if err != nil {
		return common.Address{}, err
	}
	sealHash := clique.SealHash(header).Bytes()
	log.Info("test signing of clique block",
		"Sealhash", fmt.Sprintf("%x", sealHash),
		"signature", fmt.Sprintf("%x", signature))
	pubkey, err := crypto.Ecrecover[P](sealHash, signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256[P](pubkey[1:])[12:])

	return signer, nil
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *PublicDebugAPI[T,P]) PrintBlock(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return spew.Sdump(block), nil
}

// SeedHash retrieves the seed hash of a block.
func (api *PublicDebugAPI[T,P]) SeedHash(ctx context.Context, number uint64) (string, error) {
	block, _ := api.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if block == nil {
		return "", fmt.Errorf("block #%d not found", number)
	}
	return fmt.Sprintf("0x%x", ethash.SeedHash(number)), nil
}

// PrivateDebugAPI is the collection of Ethereum APIs exposed over the private
// debugging endpoint.
type PrivateDebugAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	b Backend[T,P]
}

// NewPrivateDebugAPI creates a new API definition for the private debug methods
// of the Ethereum service.
func NewPrivateDebugAPI[T crypto.PrivateKey, P crypto.PublicKey](b Backend[T,P]) *PrivateDebugAPI[T,P] {
	return &PrivateDebugAPI[T,P]{b: b}
}

// ChaindbProperty returns leveldb properties of the key-value database.
func (api *PrivateDebugAPI[T,P]) ChaindbProperty(property string) (string, error) {
	if property == "" {
		property = "leveldb.stats"
	} else if !strings.HasPrefix(property, "leveldb.") {
		property = "leveldb." + property
	}
	return api.b.ChainDb().Stat(property)
}

// ChaindbCompact flattens the entire key-value database into a single level,
// removing all unused slots and merging all keys.
func (api *PrivateDebugAPI[T,P]) ChaindbCompact() error {
	for b := byte(0); b < 255; b++ {
		log.Info("Compacting chain database", "range", fmt.Sprintf("0x%0.2X-0x%0.2X", b, b+1))
		if err := api.b.ChainDb().Compact([]byte{b}, []byte{b + 1}); err != nil {
			log.Error("Database compaction failed", "err", err)
			return err
		}
	}
	return nil
}

// SetHead rewinds the head of the blockchain to a previous block.
func (api *PrivateDebugAPI[T,P]) SetHead(number hexutil.Uint64) {
	api.b.SetHead(uint64(number))
}

// PublicNetAPI offers network related RPC methods
type PublicNetAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	net            *p2p.Server[T,P]
	networkVersion uint64
}

// NewPublicNetAPI creates a new net API instance.
func NewPublicNetAPI[T crypto.PrivateKey, P crypto.PublicKey](net *p2p.Server[T,P], networkVersion uint64) *PublicNetAPI[T,P] {
	return &PublicNetAPI[T,P]{net, networkVersion}
}

// Listening returns an indication if the node is listening for network connections.
func (s *PublicNetAPI[T,P]) Listening() bool {
	return true // always listening
}

// PeerCount returns the number of connected peers
func (s *PublicNetAPI[T,P]) PeerCount() hexutil.Uint {
	return hexutil.Uint(s.net.PeerCount())
}

// Version returns the current ethereum protocol version.
func (s *PublicNetAPI[T,P]) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func checkTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	feeEth := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))), new(big.Float).SetInt(big.NewInt(params.Ether)))
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f ether) exceeds the configured cap (%.2f ether)", feeFloat, cap)
	}
	return nil
}

// toHexSlice creates a slice of hex-strings based on []byte.
func toHexSlice(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = hexutil.Encode(b[i])
	}
	return r
}

// Quorum
// Please note: This is a temporary integration to improve performance in high-latency
// environments when sending many private transactions. It will be removed at a later
// date when account management is handled outside Ethereum.

type AsyncSendTxArgs[T crypto.PrivateKey, P crypto.PublicKey] struct {
	SendTxArgs[T,P]
	CallbackUrl string `json:"callbackUrl"`
}

type AsyncResultSuccess struct {
	Id     string      `json:"id,omitempty"`
	TxHash common.Hash `json:"txHash"`
}

type AsyncResultFailure struct {
	Id    string `json:"id,omitempty"`
	Error string `json:"error"`
}

type Async struct {
	sync.Mutex
	sem chan struct{}
}

func (s *PublicTransactionPoolAPI[T,P]) send(ctx context.Context, asyncArgs AsyncSendTxArgs[T,P]) {

	txHash, err := s.SendTransaction(ctx, asyncArgs.SendTxArgs)

	if asyncArgs.CallbackUrl != "" {

		//don't need to nil check this since id is required for every geth rpc call
		//even though this is stated in the specification as an "optional" parameter
		jsonId := ctx.Value("id").(*json.RawMessage)
		id := string(*jsonId)

		var resultResponse interface{}
		if err != nil {
			resultResponse = &AsyncResultFailure{Id: id, Error: err.Error()}
		} else {
			resultResponse = &AsyncResultSuccess{Id: id, TxHash: txHash}
		}

		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(resultResponse)
		if err != nil {
			log.Info("Error encoding callback JSON", "err", err.Error())
			return
		}
		_, err = http.Post(asyncArgs.CallbackUrl, "application/json", buf)
		if err != nil {
			log.Info("Error sending callback", "err", err.Error())
			return
		}
	}

}

func newAsync(n int) *Async {
	a := &Async{
		sem: make(chan struct{}, n),
	}
	return a
}

var async = newAsync(100)

// SendTransactionAsync creates a transaction for the given argument, signs it, and
// submits it to the transaction pool. This call returns immediately to allow sending
// many private transactions/bursts of transactions without waiting for the recipient
// parties to confirm receipt of the encrypted payloads. An optional callbackUrl may
// be specified--when a transaction is submitted to the transaction pool, it will be
// called with a POST request containing either {"error": "error message"} or
// {"txHash": "0x..."}.
//
// Please note: This is a temporary integration to improve performance in high-latency
// environments when sending many private transactions. It will be removed at a later
// date when account management is handled outside Ethereum.
func (s *PublicTransactionPoolAPI[T,P]) SendTransactionAsync(ctx context.Context, args AsyncSendTxArgs[T,P]) (common.Hash, error) {

	select {
	case async.sem <- struct{}{}:
		go func() {
			s.send(ctx, args)
			<-async.sem
		}()
		return common.Hash{}, nil
	default:
		return common.Hash{}, errors.New("too many concurrent requests")
	}
}

// GetQuorumPayload returns the contents of a private transaction
func (s *PublicBlockChainAPI[T,P]) GetQuorumPayload(ctx context.Context, digestHex string) (string, error) {
	if !private.IsQuorumPrivacyEnabled() {
		return "", fmt.Errorf("PrivateTransactionManager is not enabled")
	}
	psm, err := s.b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return "", err
	}
	if len(digestHex) < 3 {
		return "", fmt.Errorf("Invalid digest hex")
	}
	if digestHex[:2] == "0x" {
		digestHex = digestHex[2:]
	}
	b, err := hex.DecodeString(digestHex)
	if err != nil {
		return "", err
	}

	if len(b) != common.EncryptedPayloadHashLength {
		return "", fmt.Errorf("Expected a Quorum digest of length 64, but got %d", len(b))
	}
	_, managedParties, data, _, err := private.Ptm.Receive(common.BytesToEncryptedPayloadHash(b))
	if err != nil {
		return "", err
	}
	if s.b.PSMR().NotIncludeAny(psm, managedParties...) {
		return "0x", nil
	}
	return fmt.Sprintf("0x%x", data), nil
}

func (s *PublicBlockChainAPI[T,P]) GetQuorumPayloadExtra(ctx context.Context, digestHex string) (*engine.QuorumPayloadExtra, error) {
	if !private.IsQuorumPrivacyEnabled() {
		return nil, fmt.Errorf("PrivateTransactionManager is not enabled")
	}
	psm, err := s.b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}
	if len(digestHex) < 3 {
		return nil, fmt.Errorf("Invalid digest hex")
	}
	if digestHex[:2] == "0x" {
		digestHex = digestHex[2:]
	}
	b, err := hex.DecodeString(digestHex)
	if err != nil {
		return nil, err
	}

	if len(b) != common.EncryptedPayloadHashLength {
		return nil, fmt.Errorf("Expected a Quorum digest of length 64, but got %d", len(b))
	}
	_, managedParties, data, extraMetaData, err := private.Ptm.Receive(common.BytesToEncryptedPayloadHash(b))
	if err != nil {
		return nil, err
	}
	if s.b.PSMR().NotIncludeAny(psm, managedParties...) {
		return nil, nil
	}
	isSender := false
	if len(psm.Addresses) == 0 {
		isSender, _ = private.Ptm.IsSender(common.BytesToEncryptedPayloadHash(b))
	} else {
		isSender = !psm.NotIncludeAny(extraMetaData.Sender)
	}
	return &engine.QuorumPayloadExtra{
		Payload:       fmt.Sprintf("0x%x", data),
		ExtraMetaData: extraMetaData,
		IsSender:      isSender,
	}, nil
}

// DecryptQuorumPayload returns the decrypted version of the input transaction
func (s *PublicBlockChainAPI[T,P]) DecryptQuorumPayload(ctx context.Context, payloadHex string) (*engine.QuorumPayloadExtra, error) {
	if !private.IsQuorumPrivacyEnabled() {
		return nil, fmt.Errorf("PrivateTransactionManager is not enabled")
	}
	psm, err := s.b.PSMR().ResolveForUserContext(ctx)
	if err != nil {
		return nil, err
	}
	if len(payloadHex) < 3 {
		return nil, fmt.Errorf("Invalid payload hex")
	}
	if payloadHex[:2] == "0x" {
		payloadHex = payloadHex[2:]
	}
	b, err := hex.DecodeString(payloadHex)
	if err != nil {
		return nil, err
	}

	var payload common.DecryptRequest
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, err
	}
	// if we are MPS and the sender is not part of the resolved PSM - return empty
	if len(psm.Addresses) != 0 && psm.NotIncludeAny(base64.StdEncoding.EncodeToString(payload.SenderKey)) {
		return nil, nil
	}
	data, extraMetaData, err := private.Ptm.DecryptPayload(payload)
	if err != nil {
		return nil, err
	}

	return &engine.QuorumPayloadExtra{
		Payload:       fmt.Sprintf("0x%x", data),
		ExtraMetaData: extraMetaData,
		IsSender:      true,
	}, nil
}

// Quorum
// for raw private transaction, privateTxArgs.privateFrom will be updated with value from Tessera when payload is retrieved
func checkAndHandlePrivateTransaction[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], tx *types.Transaction[P], privateTxArgs *PrivateTxArgs[T,P], from common.Address, txnType TransactionType) (isPrivate bool, replaceDataWithHash bool, hash common.EncryptedPayloadHash, err error) {
	replaceDataWithHash = false
	isPrivate = privateTxArgs != nil && privateTxArgs.PrivateFor != nil
	if !isPrivate {
		return
	}

	if err = privateTxArgs.PrivacyFlag.Validate(); err != nil {
		return
	}

	if !b.ChainConfig().IsPrivacyEnhancementsEnabled(b.CurrentBlock().Number()) && privateTxArgs.PrivacyFlag.IsNotStandardPrivate() {
		err = fmt.Errorf("PrivacyEnhancements are disabled. Can only accept transactions with PrivacyFlag=0(StandardPrivate).")
		return
	}

	if engine.PrivacyFlagMandatoryRecipients == privateTxArgs.PrivacyFlag && len(privateTxArgs.MandatoryRecipients) == 0 {
		err = fmt.Errorf("missing mandatory recipients data. if no mandatory recipients required consider using PrivacyFlag=1(PartyProtection)")
		return
	}

	if engine.PrivacyFlagMandatoryRecipients != privateTxArgs.PrivacyFlag && len(privateTxArgs.MandatoryRecipients) > 0 {
		err = fmt.Errorf("privacy metadata invalid. mandatory recipients are only applicable for PrivacyFlag=2(MandatoryRecipients)")
		return
	}

	// validate that PrivateFrom is one of the addresses of the private state resolved from the user context
	if b.ChainConfig().IsMPS {
		var psm *mps.PrivateStateMetadata
		psm, err = b.PSMR().ResolveForUserContext(ctx)
		if err != nil {
			return
		}
		if psm.NotIncludeAny(privateTxArgs.PrivateFrom) {
			err = fmt.Errorf("The PrivateFrom (%s) address does not match the specified private state (%s) ", privateTxArgs.PrivateFrom, psm.ID)
			return
		}
	}

	if len(tx.Data()) > 0 {
		// check private contract exists on the node initiating the transaction
		if tx.To() != nil && privateTxArgs.PrivacyFlag.IsNotStandardPrivate() {
			state, _, lerr := b.StateAndHeaderByNumber(ctx, rpc.BlockNumber(b.CurrentBlock().Number().Uint64()))
			if lerr != nil && state == nil {
				err = fmt.Errorf("state not found")
				return
			}
			if state.GetCode(*tx.To()) == nil {
				err = fmt.Errorf("contract not found. cannot transact")
				return
			}
		}

		replaceDataWithHash = true
		hash, err = handlePrivateTransaction(ctx, b, tx, privateTxArgs, from, txnType)
	}

	return
}

// Quorum
// If transaction is raw, the tx payload is indeed the hash of the encrypted payload.
// Then the sender key will set to privateTxArgs.privateFrom.
//
// For private transaction, run a simulated execution in order to
// 1. Find all affected private contract accounts then retrieve encrypted payload hashes of their creation txs
// 2. Calculate Merkle Root as the result of the simulated execution
// The above information along with private originating payload are sent to Transaction Manager
// to obtain hash of the encrypted private payload
func handlePrivateTransaction[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], tx *types.Transaction[P], privateTxArgs *PrivateTxArgs[T,P], from common.Address, txnType TransactionType) (hash common.EncryptedPayloadHash, err error) {
	defer func(start time.Time) {
		log.Debug("Handle Private Transaction finished", "took", time.Since(start))
	}(time.Now())

	data := tx.Data()

	log.Debug("sending private tx", "txnType", txnType, "data", common.FormatTerminalString(data), "privatefrom", privateTxArgs.PrivateFrom, "privatefor", privateTxArgs.PrivateFor, "privacyFlag", privateTxArgs.PrivacyFlag, "mandatoryfor", privateTxArgs.MandatoryRecipients)

	switch txnType {
	case FillTransaction:
		hash, err = private.Ptm.StoreRaw(data, privateTxArgs.PrivateFrom)
	case RawTransaction:
		hash, err = handleRawPrivateTransaction(ctx, b, tx, privateTxArgs, from)
	case NormalTransaction:
		hash, err = handleNormalPrivateTransaction(ctx, b, tx, data, privateTxArgs, from)
	}
	return
}

// Quorum
func handleRawPrivateTransaction[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], tx *types.Transaction[P], privateTxArgs *PrivateTxArgs[T,P], from common.Address) (hash common.EncryptedPayloadHash, err error) {
	data := tx.Data()
	hash = common.BytesToEncryptedPayloadHash(data)
	privatePayload, privateFrom, _, revErr := private.Ptm.ReceiveRaw(hash)
	if revErr != nil {
		return common.EncryptedPayloadHash{}, revErr
	}
	log.Trace("received raw payload", "hash", hash, "privatepayload", common.FormatTerminalString(privatePayload), "privateFrom", privateFrom)

	privateTxArgs.PrivateFrom = privateFrom
	var privateTx *types.Transaction[P]
	if tx.To() == nil {
		privateTx = types.NewContractCreation[P](tx.Nonce(), tx.Value(), tx.Gas(), tx.GasPrice(), privatePayload)
	} else {
		privateTx = types.NewTransaction[P](tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), tx.GasPrice(), privatePayload)
	}

	affectedCATxHashes, merkleRoot, err := simulateExecutionForPE(ctx, b, from, privateTx, privateTxArgs)
	log.Trace("after simulation", "affectedCATxHashes", affectedCATxHashes, "merkleRoot", merkleRoot, "privacyFlag", privateTxArgs.PrivacyFlag, "error", err)
	if err != nil {
		return
	}

	metadata := engine.ExtraMetadata{
		ACHashes:            affectedCATxHashes,
		ACMerkleRoot:        merkleRoot,
		PrivacyFlag:         privateTxArgs.PrivacyFlag,
		MandatoryRecipients: privateTxArgs.MandatoryRecipients,
	}
	_, _, data, err = private.Ptm.SendSignedTx(hash, privateTxArgs.PrivateFor, &metadata)
	if err != nil {
		return
	}

	log.Info("sent raw private signed tx",
		"data", common.FormatTerminalString(data),
		"hash", hash,
		"privatefrom", privateTxArgs.PrivateFrom,
		"privatefor", privateTxArgs.PrivateFor,
		"affectedCATxHashes", metadata.ACHashes,
		"merkleroot", metadata.ACHashes,
		"privacyflag", metadata.PrivacyFlag,
		"mandatoryrecipients", metadata.MandatoryRecipients)
	return
}

// Quorum
func handleNormalPrivateTransaction[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], tx *types.Transaction[P], data []byte, privateTxArgs *PrivateTxArgs[T,P], from common.Address) (hash common.EncryptedPayloadHash, err error) {
	affectedCATxHashes, merkleRoot, err := simulateExecutionForPE(ctx, b, from, tx, privateTxArgs)
	log.Trace("after simulation", "affectedCATxHashes", affectedCATxHashes, "merkleRoot", merkleRoot, "privacyFlag", privateTxArgs.PrivacyFlag, "error", err)
	if err != nil {
		return
	}

	metadata := engine.ExtraMetadata{
		ACHashes:            affectedCATxHashes,
		ACMerkleRoot:        merkleRoot,
		PrivacyFlag:         privateTxArgs.PrivacyFlag,
		MandatoryRecipients: privateTxArgs.MandatoryRecipients,
	}
	_, _, hash, err = private.Ptm.Send(data, privateTxArgs.PrivateFrom, privateTxArgs.PrivateFor, &metadata)
	if err != nil {
		return
	}

	log.Info("sent private signed tx",
		"data", common.FormatTerminalString(data),
		"hash", hash,
		"privatefrom", privateTxArgs.PrivateFrom,
		"privatefor", privateTxArgs.PrivateFor,
		"affectedCATxHashes", metadata.ACHashes,
		"merkleroot", metadata.ACHashes,
		"privacyflag", metadata.PrivacyFlag,
		"mandatoryrecipients", metadata.MandatoryRecipients)
	return
}

// (Quorum) createPrivacyMarkerTransaction creates a new privacy marker transaction (PMT) with the given signed privateTx.
// The private tx is sent only to the privateFor recipients. The resulting PMT's 'to' is the privacy precompile address and its 'data' is the
// privacy manager hash for the private tx.
func createPrivacyMarkerTransaction[T crypto.PrivateKey, P crypto.PublicKey](b Backend[T,P], privateTx *types.Transaction[P], privateTxArgs *PrivateTxArgs[T,P]) (*types.Transaction[P], error) {
	log.Trace("creating privacy marker transaction", "from", privateTx.From(), "to", privateTx.To())

	data := new(bytes.Buffer)
	err := json.NewEncoder(data).Encode(privateTx)
	if err != nil {
		return nil, err
	}

	_, _, ptmHash, err := private.Ptm.Send(data.Bytes(), privateTxArgs.PrivateFrom, privateTxArgs.PrivateFor, &engine.ExtraMetadata{})
	if err != nil {
		return nil, err
	}

	currentBlockHeight := b.CurrentHeader().Number
	istanbul := b.ChainConfig().IsIstanbul(currentBlockHeight)
	intrinsicGas, err := core.IntrinsicGas(ptmHash.Bytes(), privateTx.AccessList(), false, true, istanbul)
	if err != nil {
		return nil, err
	}

	pmt := types.NewTransaction[P](privateTx.Nonce(), common.QuorumPrivacyPrecompileContractAddress(), privateTx.Value(), intrinsicGas, privateTx.GasPrice(), ptmHash.Bytes())

	return pmt, nil
}

// Quorum
// simulateExecutionForPE simulates execution of a private transaction for enhanced privacy
//
// Returns hashes of encrypted payload of creation transactions for all affected contract accounts
// and the merkle root combining all affected contract accounts after the simulation
func simulateExecutionForPE[T crypto.PrivateKey, P crypto.PublicKey](ctx context.Context, b Backend[T,P], from common.Address, privateTx *types.Transaction[P], privateTxArgs *PrivateTxArgs[T,P]) (common.EncryptedPayloadHashes, common.Hash, error) {
	// skip simulation if privacy enhancements are disabled
	if !b.ChainConfig().IsPrivacyEnhancementsEnabled(b.CurrentBlock().Number()) {
		return nil, common.Hash{}, nil
	}

	evm, data, err := runSimulation(ctx, b, from, privateTx)
	if evm == nil {
		log.Debug("TX Simulation setup failed", "error", err)
		return nil, common.Hash{}, err
	}
	if err != nil {
		if privateTxArgs.PrivacyFlag.IsStandardPrivate() {
			log.Debug("An error occurred during StandardPrivate transaction simulation. "+
				"Continuing to simulation checks.", "error", err)
		} else {
			log.Trace("Simulated execution", "error", err)
			if len(data) > 0 && errors.Is(err, vm.ErrExecutionReverted) {
				reason, errUnpack := abi.UnpackRevert[P](data)

				reasonError := errors.New("execution reverted")
				if errUnpack == nil {
					reasonError = fmt.Errorf("execution reverted: %v", reason)
				}
				err = &revertError{
					error:  reasonError,
					reason: hexutil.Encode(data),
				}

			}
			return nil, common.Hash{}, err
		}
	}
	affectedContractsHashes := make(common.EncryptedPayloadHashes)
	var merkleRoot common.Hash
	addresses := evm.AffectedContracts()
	privacyFlag := privateTxArgs.PrivacyFlag
	log.Trace("after simulation run", "numberOfAffectedContracts", len(addresses), "privacyFlag", privacyFlag)
	for _, addr := range addresses {
		// GetPrivacyMetadata is invoked directly on the privateState (as the tx is private) and it returns:
		// 1. public contacts: privacyMetadata = nil, err = nil
		// 2. private contracts of type:
		// 2.1. StandardPrivate:     privacyMetadata = nil, err = "The provided contract does not have privacy metadata"
		// 2.2. PartyProtection/PSV: privacyMetadata = <data>, err = nil
		privacyMetadata, err := evm.StateDB.GetPrivacyMetadata(addr)
		log.Debug("Found affected contract", "address", addr.Hex(), "privacyMetadata", privacyMetadata)
		//privacyMetadata not found=non-party, or another db error
		if err != nil && privacyFlag.IsNotStandardPrivate() {
			return nil, common.Hash{}, errors.New("PrivacyMetadata not found: " + err.Error())
		}
		// when we run simulation, it's possible that affected contracts may contain public ones
		// public contract will not have any privacyMetadata attached
		// standard private will be nil
		if privacyMetadata == nil {
			continue
		}
		//if affecteds are not all the same return an error
		if privacyFlag != privacyMetadata.PrivacyFlag {
			return nil, common.Hash{}, errors.New("sent privacy flag doesn't match all affected contract flags")
		}

		affectedContractsHashes.Add(privacyMetadata.CreationTxHash)
	}
	//only calculate the merkle root if all contracts are psv
	if privacyFlag.Has(engine.PrivacyFlagStateValidation) {
		merkleRoot, err = evm.CalculateMerkleRoot()
		if err != nil {
			return nil, common.Hash{}, err
		}
	}
	log.Trace("post-execution run", "merkleRoot", merkleRoot, "affectedhashes", affectedContractsHashes)
	return affectedContractsHashes, merkleRoot, nil
}

//End-Quorum
