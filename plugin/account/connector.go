package account

import (
	"context"

	"github.com/MIRChain/MIR/crypto"
	iplugin "github.com/MIRChain/MIR/internal/plugin"
	"github.com/hashicorp/go-plugin"
	"github.com/jpmorganchase/quorum-account-plugin-sdk-go/proto"
	"google.golang.org/grpc"
)

const ConnectorName = "account"

type PluginConnector[T crypto.PrivateKey, P crypto.PublicKey] struct {
	plugin.Plugin
}

func (*PluginConnector[T, P]) GRPCServer(_ *plugin.GRPCBroker, _ *grpc.Server) error {
	return iplugin.ErrNotSupported
}

func (*PluginConnector[T, P]) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, cc *grpc.ClientConn) (interface{}, error) {
	return &service[T, P]{
		client: proto.NewAccountServiceClient(cc),
	}, nil
}
