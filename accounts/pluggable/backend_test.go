package pluggable

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/pluggable/internal/testutils/mock_plugin"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/stretchr/testify/require"
)

func TestBackend_Subscribe_NoOp(t *testing.T) {
	b := NewBackend[nist.PublicKey]()

	subscriber := make(chan accounts.WalletEvent[nist.PublicKey], 4)
	sub := b.Subscribe(subscriber)
	require.NotNil(t, sub)
	require.Len(t, subscriber, 0)

	sub.Unsubscribe()
	require.Len(t, subscriber, 0)
}

func TestBackend_Wallets_ReturnsCopy(t *testing.T) {
	wallets := []accounts.Wallet[nist.PublicKey]{
		&wallet[nist.PublicKey]{
			url: accounts.URL{
				Scheme: "http",
				Path:   "url1",
			},
		},
		&wallet[nist.PublicKey]{
			url: accounts.URL{
				Scheme: "http",
				Path:   "url2",
			},
		},
	}

	b := NewBackend[nist.PublicKey]()
	b.wallets = wallets

	got := b.Wallets()
	got[0] = &wallet[nist.PublicKey]{
		url: accounts.URL{
			Scheme: "http",
			Path:   "changedurl",
		},
	}
	require.Equal(t, "changedurl", got[0].URL().Path)

	unchanged := b.Wallets()
	require.Equal(t, "url1", unchanged[0].URL().Path)
}

func TestBackend_TimedUnlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_plugin.NewMockService(ctrl)
	mockClient.
		EXPECT().
		TimedUnlock(gomock.Any(), gomock.Eq(acct1), gomock.Eq("pwd"), gomock.Eq(time.Minute)).
		Return(nil)

	b := NewBackend[nist.PublicKey]()
	b.wallets[0].(*wallet[nist.PublicKey]).pluginService = mockClient

	err := b.TimedUnlock(acct1, "pwd", time.Minute)
	require.NoError(t, err)
}

func TestBackend_Lock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_plugin.NewMockService(ctrl)
	mockClient.
		EXPECT().
		Lock(gomock.Any(), gomock.Eq(acct1)).
		Return(nil)

	b := NewBackend[nist.PublicKey]()
	b.wallets[0].(*wallet[nist.PublicKey]).pluginService = mockClient

	err := b.Lock(acct1)
	require.NoError(t, err)
}

func TestBackend_NewAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newAccountConfig := struct{ config string }{config: "someconfig"}
	newAccount := accounts.Account{}

	mockClient := mock_plugin.NewMockService(ctrl)
	mockClient.
		EXPECT().
		NewAccount(gomock.Any(), gomock.Eq(newAccountConfig)).
		Return(newAccount, nil)

	b := NewBackend[nist.PublicKey]()
	b.wallets[0].(*wallet[nist.PublicKey]).pluginService = mockClient

	got, err := b.NewAccount(newAccountConfig)
	require.NoError(t, err)
	require.Equal(t, newAccount, got)
}

func TestBackend_ImportRawKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newAccountConfig := struct{ config string }{config: "someconfig"}
	newAccount := accounts.Account{}

	mockClient := mock_plugin.NewMockService(ctrl)
	mockClient.
		EXPECT().
		ImportRawKey(gomock.Any(), gomock.Eq("rawkey"), gomock.Eq(newAccountConfig)).
		Return(newAccount, nil)

	b := NewBackend[nist.PublicKey]()
	b.wallets[0].(*wallet[nist.PublicKey]).pluginService = mockClient

	got, err := b.ImportRawKey("rawkey", newAccountConfig)
	require.NoError(t, err)
	require.Equal(t, newAccount, got)
}
