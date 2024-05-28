package pluggable

import (
	"context"
	"math/big"
	"sync"
	"time"

	ethereum "github.com/MIRChain/MIR"
	"github.com/MIRChain/MIR/accounts"
	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	plugin "github.com/MIRChain/MIR/plugin/account"
)

type wallet[P crypto.PublicKey] struct {
	url           accounts.URL
	mu            sync.Mutex
	pluginService plugin.Service
}

func (w *wallet[P]) setPluginService(s plugin.Service) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.pluginService = s

	return nil
}

func (w *wallet[P]) URL() accounts.URL {
	return w.url
}

func (w *wallet[P]) Status() (string, error) {
	return w.pluginService.Status(context.Background())
}

func (w *wallet[P]) Open(passphrase string) error {
	return w.pluginService.Open(context.Background(), passphrase)
}

func (w *wallet[P]) Close() error {
	return w.pluginService.Close(context.Background())
}

func (w *wallet[P]) Accounts() []accounts.Account {
	if w.pluginService == nil {
		return []accounts.Account{}
	}
	return w.pluginService.Accounts(context.Background())
}

func (w *wallet[P]) Contains(account accounts.Account) bool {
	return w.pluginService.Contains(context.Background(), account)
}

func (w *wallet[P]) Derive(_ accounts.DerivationPath, _ bool) (accounts.Account, error) {
	return accounts.Account{}, accounts.ErrNotSupported
}

func (w *wallet[P]) SelfDerive(_ []accounts.DerivationPath, _ ethereum.ChainStateReader) {}

func (w *wallet[P]) SignData(account accounts.Account, _ string, data []byte) ([]byte, error) {
	return w.pluginService.Sign(context.Background(), account, crypto.Keccak256[P](data))
}

func (w *wallet[P]) SignDataWithPassphrase(account accounts.Account, passphrase, _ string, data []byte) ([]byte, error) {
	return w.pluginService.UnlockAndSign(context.Background(), account, crypto.Keccak256[P](data), passphrase)
}

func (w *wallet[P]) SignText(account accounts.Account, text []byte) ([]byte, error) {
	return w.pluginService.Sign(context.Background(), account, accounts.TextHash(text))
}

func (w *wallet[P]) SignTextWithPassphrase(account accounts.Account, passphrase string, text []byte) ([]byte, error) {
	return w.pluginService.UnlockAndSign(context.Background(), account, accounts.TextHash(text), passphrase)
}

func (w *wallet[P]) SignTx(account accounts.Account, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	toSign, signer := prepareTxForSign(tx, chainID)

	sig, err := w.pluginService.Sign(context.Background(), account, toSign.Bytes())
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signer, sig)
}

func (w *wallet[P]) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction[P], chainID *big.Int) (*types.Transaction[P], error) {
	toSign, signer := prepareTxForSign(tx, chainID)

	sig, err := w.pluginService.UnlockAndSign(context.Background(), account, toSign.Bytes(), passphrase)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(signer, sig)
}

func (w *wallet[P]) timedUnlock(account accounts.Account, password string, duration time.Duration) error {
	return w.pluginService.TimedUnlock(context.Background(), account, password, duration)
}

func (w *wallet[P]) lock(account accounts.Account) error {
	return w.pluginService.Lock(context.Background(), account)
}

func (w *wallet[P]) newAccount(newAccountConfig interface{}) (accounts.Account, error) {
	return w.pluginService.NewAccount(context.Background(), newAccountConfig)
}

func (w *wallet[P]) importRawKey(rawKey string, newAccountConfig interface{}) (accounts.Account, error) {
	return w.pluginService.ImportRawKey(context.Background(), rawKey, newAccountConfig)
}

// prepareTxForSign determines which Signer to use for the given tx and chainID, and returns the Signer's hash of the tx and the Signer itself
func prepareTxForSign[P crypto.PublicKey](tx *types.Transaction[P], chainID *big.Int) (common.Hash, types.Signer[P]) {
	var s types.Signer[P]

	if tx.IsPrivate() {
		s = types.QuorumPrivateTxSigner[P]{}
	} else {
		s = types.LatestSignerForChainID[P](chainID)
	}

	return s.Hash(tx), s
}
