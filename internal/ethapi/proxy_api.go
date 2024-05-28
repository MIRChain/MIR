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
	"context"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/common/hexutil"
	"github.com/MIRChain/MIR/core/types"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/log"
	"github.com/MIRChain/MIR/rpc"
)

type ProxyAPISupport interface {
	ProxyEnabled() bool
	ProxyClient() *rpc.Client
}

// PublicTransactionPoolAPI exposes methods for the RPC interface
type PublicTransactionPoolProxyAPI[T crypto.PrivateKey, P crypto.PublicKey] struct {
	PublicTransactionPoolAPI[T, P]
	proxyClient *rpc.Client
}

// NewPublicTransactionPoolAPI creates a new RPC service with methods specific for the transaction pool.
func NewPublicTransactionPoolProxyAPI[T crypto.PrivateKey, P crypto.PublicKey](b Backend[T, P], nonceLock *AddrLocker) interface{} {
	apiSupport, ok := b.(ProxyAPISupport)
	if ok && apiSupport.ProxyEnabled() {
		signer := types.LatestSigner[P](b.ChainConfig())
		return &PublicTransactionPoolProxyAPI[T, P]{
			PublicTransactionPoolAPI[T, P]{b, nonceLock, signer},
			apiSupport.ProxyClient(),
		}
	}
	return NewPublicTransactionPoolAPI(b, nonceLock)
}

func (s *PublicTransactionPoolProxyAPI[T, P]) SendTransaction(ctx context.Context, args SendTxArgs[T, P]) (common.Hash, error) {
	log.Info("QLight - proxy enabled")
	var result common.Hash
	err := s.proxyClient.CallContext(ctx, &result, "eth_sendTransaction", args)
	return result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
	log.Info("QLight - proxy enabled")
	var result common.Hash
	err := s.proxyClient.CallContext(ctx, &result, "eth_sendRawTransaction", encodedTx)
	return result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) SendRawPrivateTransaction(ctx context.Context, encodedTx hexutil.Bytes, args SendRawTxArgs[T, P]) (common.Hash, error) {
	log.Info("QLight - proxy enabled")
	var result common.Hash
	err := s.proxyClient.CallContext(ctx, &result, "eth_sendRawPrivateTransaction", encodedTx, args)
	return result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) FillTransaction(ctx context.Context, args SendTxArgs[T, P]) (*SignTransactionResult[P], error) {
	log.Info("QLight - proxy enabled")
	var result SignTransactionResult[P]
	err := s.proxyClient.CallContext(ctx, &result, "eth_fillTransaction", args)
	return &result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) DistributePrivateTransaction(ctx context.Context, encodedTx hexutil.Bytes, args SendRawTxArgs[T, P]) (string, error) {
	log.Info("QLight - proxy enabled")
	var result string
	err := s.proxyClient.CallContext(ctx, &result, "eth_distributePrivateTransaction", encodedTx, args)
	return result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) Resend(ctx context.Context, sendArgs SendTxArgs[T, P], gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	log.Info("QLight - proxy enabled")
	var result common.Hash
	err := s.proxyClient.CallContext(ctx, &result, "eth_resend", sendArgs, gasPrice, gasLimit)
	return result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) SendTransactionAsync(ctx context.Context, args AsyncSendTxArgs[T, P]) (common.Hash, error) {
	log.Info("QLight - proxy enabled")
	var result common.Hash
	err := s.proxyClient.CallContext(ctx, &result, "eth_sendTransactionAsync", args)
	return result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	log.Info("QLight - proxy enabled")
	var result hexutil.Bytes
	err := s.proxyClient.Call(&result, "eth_sign", addr, data)
	return result, err
}

func (s *PublicTransactionPoolProxyAPI[T, P]) SignTransaction(ctx context.Context, args SendTxArgs[T, P]) (*SignTransactionResult[P], error) {
	log.Info("QLight - proxy enabled")
	var result SignTransactionResult[P]
	err := s.proxyClient.CallContext(ctx, &result, "eth_signTransaction", args)
	return &result, err
}

type PrivateAccountProxyAPI[T crypto.PrivateKey, P crypto.PublicKey] struct {
	PrivateAccountAPI[T, P]
	proxyClient *rpc.Client
}

func NewPrivateAccountProxyAPI[T crypto.PrivateKey, P crypto.PublicKey](b Backend[T, P], nonceLock *AddrLocker) interface{} {
	apiSupport, ok := b.(ProxyAPISupport)
	if ok && apiSupport.ProxyEnabled() {
		return &PrivateAccountProxyAPI[T, P]{
			PrivateAccountAPI[T, P]{
				am:        b.AccountManager(),
				nonceLock: nonceLock,
				b:         b,
			},
			apiSupport.ProxyClient(),
		}
	}
	return NewPrivateAccountAPI[T, P](b, nonceLock)
}

func (s *PrivateAccountProxyAPI[T, P]) SendTransaction(ctx context.Context, args SendTxArgs[T, P], passwd string) (common.Hash, error) {
	log.Info("QLight - proxy enabled")
	var result common.Hash
	err := s.proxyClient.CallContext(ctx, &result, "personal_sendTransaction", args, passwd)
	return result, err
}

func (s *PrivateAccountProxyAPI[T, P]) SignTransaction(ctx context.Context, args SendTxArgs[T, P], passwd string) (*SignTransactionResult[P], error) {
	log.Info("QLight - proxy enabled")
	var result SignTransactionResult[P]
	err := s.proxyClient.CallContext(ctx, &result, "personal_signTransaction", args, passwd)
	return &result, err
}

func (s *PrivateAccountProxyAPI[T, P]) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	log.Info("QLight - proxy enabled")
	var result hexutil.Bytes
	err := s.proxyClient.CallContext(ctx, &result, "personal_sign", data, addr, passwd)
	return result, err
}

func (s *PrivateAccountProxyAPI[T, P]) SignAndSendTransaction(ctx context.Context, args SendTxArgs[T, P], passwd string) (common.Hash, error) {
	return s.SendTransaction(ctx, args, passwd)
}
