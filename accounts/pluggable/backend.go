package pluggable

import (
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/event"
	plugin "github.com/pavelkrolevets/MIR-pro/plugin/account"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// var BackendType = reflect.TypeOf(&Backend{})

type Backend [P crypto.PublicKey]struct {
	wallets []accounts.Wallet[P]
}

func NewBackend[P crypto.PublicKey]() *Backend[P] {
	return &Backend[P]{
		wallets: []accounts.Wallet[P]{
			&wallet[P]{
				url: accounts.URL{
					Scheme: "plugin",
					Path:   "account",
				},
			},
		},
	}
}

func (b *Backend[P]) Wallets() []accounts.Wallet[P] {
	cpy := make([]accounts.Wallet[P], len(b.wallets))
	copy(cpy, b.wallets)
	return cpy
}

// Subscribe implements accounts.Backend, creating a new subscription that is a no-op and simply exits when the Unsubscribe is called
func (b *Backend[P]) Subscribe(_ chan<- accounts.WalletEvent[P]) event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		<-quit
		return nil
	})
}

func (b *Backend[P]) SetPluginService(s plugin.Service) error {
	return b.wallet().setPluginService(s)
}

func (b *Backend[P]) TimedUnlock(account accounts.Account, password string, duration time.Duration) error {
	return b.wallet().timedUnlock(account, password, duration)
}

func (b *Backend[P]) Lock(account accounts.Account) error {
	return b.wallet().lock(account)
}

// AccountCreator is the interface that wraps the plugin account creation methods.
// This interface is used to simplify the pluggable.Backend API available to the account plugin CLI and enables easier testing.
type AccountCreator interface {
	NewAccount(newAccountConfig interface{}) (accounts.Account, error)
	ImportRawKey(rawKey string, newAccountConfig interface{}) (accounts.Account, error)
}

func (b *Backend[P]) NewAccount(newAccountConfig interface{}) (accounts.Account, error) {
	return b.wallet().newAccount(newAccountConfig)
}

func (b *Backend[P]) ImportRawKey(rawKey string, newAccountConfig interface{}) (accounts.Account, error) {
	return b.wallet().importRawKey(rawKey, newAccountConfig)
}

func (b *Backend[P]) wallet() *wallet[P] {
	return b.wallets[0].(*wallet[P])
}
