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

package core

import (
	"context"

	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

type StdIOUI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	client rpc.Client
}

func NewStdIOUI[T crypto.PrivateKey, P crypto.PublicKey]() *StdIOUI[T,P] {
	client, err := rpc.DialContext(context.Background(), "stdio://")
	if err != nil {
		log.Crit("Could not create stdio client", "err", err)
	}
	ui := &StdIOUI[T,P]{client: *client}
	return ui
}

func (ui *StdIOUI[T,P]) RegisterUIServer(api *UIServerAPI[T,P]) {
	ui.client.RegisterName("clef", api)
}

// dispatch sends a request over the stdio
func (ui *StdIOUI[T,P]) dispatch(serviceMethod string, args interface{}, reply interface{}) error {
	err := ui.client.Call(&reply, serviceMethod, args)
	if err != nil {
		log.Info("Error", "exc", err.Error())
	}
	return err
}

// notify sends a request over the stdio, and does not listen for a response
func (ui *StdIOUI[T,P]) notify(serviceMethod string, args interface{}) error {
	ctx := context.Background()
	err := ui.client.Notify(ctx, serviceMethod, args)
	if err != nil {
		log.Info("Error", "exc", err.Error())
	}
	return err
}

func (ui *StdIOUI[T,P]) ApproveTx(request *SignTxRequest[P]) (SignTxResponse[P], error) {
	var result SignTxResponse[P]
	err := ui.dispatch("ui_approveTx", request, &result)
	return result, err
}

func (ui *StdIOUI[T,P]) ApproveSignData(request *SignDataRequest) (SignDataResponse, error) {
	var result SignDataResponse
	err := ui.dispatch("ui_approveSignData", request, &result)
	return result, err
}

func (ui *StdIOUI[T,P]) ApproveListing(request *ListRequest) (ListResponse, error) {
	var result ListResponse
	err := ui.dispatch("ui_approveListing", request, &result)
	return result, err
}

func (ui *StdIOUI[T,P]) ApproveNewAccount(request *NewAccountRequest) (NewAccountResponse, error) {
	var result NewAccountResponse
	err := ui.dispatch("ui_approveNewAccount", request, &result)
	return result, err
}

func (ui *StdIOUI[T,P]) ShowError(message string) {
	err := ui.notify("ui_showError", &Message{message})
	if err != nil {
		log.Info("Error calling 'ui_showError'", "exc", err.Error(), "msg", message)
	}
}

func (ui *StdIOUI[T,P]) ShowInfo(message string) {
	err := ui.notify("ui_showInfo", Message{message})
	if err != nil {
		log.Info("Error calling 'ui_showInfo'", "exc", err.Error(), "msg", message)
	}
}
func (ui *StdIOUI[T,P]) OnApprovedTx(tx ethapi.SignTransactionResult[P]) {
	err := ui.notify("ui_onApprovedTx", tx)
	if err != nil {
		log.Info("Error calling 'ui_onApprovedTx'", "exc", err.Error(), "tx", tx)
	}
}

func (ui *StdIOUI[T,P]) OnSignerStartup(info StartupInfo) {
	err := ui.notify("ui_onSignerStartup", info)
	if err != nil {
		log.Info("Error calling 'ui_onSignerStartup'", "exc", err.Error(), "info", info)
	}
}
func (ui *StdIOUI[T,P]) OnInputRequired(info UserInputRequest) (UserInputResponse, error) {
	var result UserInputResponse
	err := ui.dispatch("ui_onInputRequired", info, &result)
	if err != nil {
		log.Info("Error calling 'ui_onInputRequired'", "exc", err.Error(), "info", info)
	}
	return result, err
}
