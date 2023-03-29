package core

import (
	"context"
	"errors"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/plugin/account"
)

// <Quorum>

type approvalCreatorService [T crypto.PrivateKey, P crypto.PublicKey] struct {
	creator account.CreatorService
	ui      UIClientAPI[T,P]
}

// NewApprovalCreatorService adds a wrapper to the provided creator service which requires UI approval before executing the service's methods
func NewApprovalCreatorService[T crypto.PrivateKey, P crypto.PublicKey](creator account.CreatorService, ui UIClientAPI[T,P]) account.CreatorService {
	return &approvalCreatorService[T,P]{
		creator: creator,
		ui:      ui,
	}
}

func (s *approvalCreatorService[T,P]) NewAccount(ctx context.Context, newAccountConfig interface{}) (accounts.Account, error) {
	if resp, err := s.ui.ApproveNewAccount(&NewAccountRequest{MetadataFromContext(ctx)}); err != nil {
		return accounts.Account{}, err
	} else if !resp.Approved {
		return accounts.Account{}, ErrRequestDenied
	}

	return s.creator.NewAccount(ctx, newAccountConfig)
}

// ImportRawKey is unsupported in the clef external API for parity with the available keystore account functionality
func (s *approvalCreatorService[T,P]) ImportRawKey(_ context.Context, _ string, _ interface{}) (accounts.Account, error) {
	return accounts.Account{}, errors.New("not supported")
}
