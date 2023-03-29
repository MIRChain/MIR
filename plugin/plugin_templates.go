package plugin

import (
	"context"

	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/plugin/account"
	"github.com/pavelkrolevets/MIR-pro/plugin/helloworld"
	"github.com/pavelkrolevets/MIR-pro/plugin/qlight"
	"github.com/pavelkrolevets/MIR-pro/plugin/security"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// a template that returns the hello world plugin instance
type HelloWorldPluginTemplate [T crypto.PrivateKey, P crypto.PublicKey] struct {
	*basePlugin[T,P]
}

func (p *HelloWorldPluginTemplate[T,P]) Get() (helloworld.PluginHelloWorld, error) {
	return &helloworld.ReloadablePluginHelloWorld{
		DeferFunc: func() (helloworld.PluginHelloWorld, error) {
			raw, err := p.dispense(helloworld.ConnectorName)
			if err != nil {
				return nil, err
			}
			return raw.(helloworld.PluginHelloWorld), nil
		},
	}, nil
}

type SecurityPluginTemplate [T crypto.PrivateKey, P crypto.PublicKey] struct {
	*basePlugin[T,P]
}

// TLSConfigurationSource returns an implementation of security.TLSConfigurationSource which could be nil
// in case the plugin doesn't implement the corresponding service. In order to verify that, it attempts
// to make a call and inspect the error.
func (sp *SecurityPluginTemplate[T,P]) TLSConfigurationSource() (security.TLSConfigurationSource, error) {
	raw, err := sp.dispense(security.TLSConfigurationConnectorName)
	if err != nil {
		return nil, err
	}
	tlsConfigurationSource := raw.(security.TLSConfigurationSource)
	// try to invoke the method to test if the plugin actually implements the service
	_, err = tlsConfigurationSource.Get(context.Background())
	rpcStatus, ok := status.FromError(err)
	if ok && rpcStatus.Code() == codes.Unimplemented {
		log.Info("Security: Plugin doesn't implement TLSConfigurationSource service", "err", err)
		return nil, nil
	}
	return tlsConfigurationSource, nil
}

// AuthenticationManager returns an implementation of security.AuthenticationManager which could be
// a deferred implemenation or a disabled implementation.
//
// The deferred implementation delegates to the actual implemenation (which is the plugin client).
//
// The disabled implementation allows no authentication verification.
func (sp *SecurityPluginTemplate[T,P]) AuthenticationManager() (security.AuthenticationManager, error) {
	deferFunc := func() (security.AuthenticationManager, error) {
		raw, err := sp.dispense(security.AuthenticationConnectorName)
		if err != nil {
			return nil, err
		}
		return raw.(security.AuthenticationManager), nil
	}
	if am, err := deferFunc(); err != nil {
		return nil, err
	} else {
		// try to invoke the method to test if the plugin actually implements the service
		_, err = am.Authenticate(context.Background(), "")
		rpcStatus, ok := status.FromError(err)
		if ok && rpcStatus.Code() == codes.Unimplemented {
			log.Info("Security: Plugin doesn't implement AuthenticationManager service", "err", err)
			return security.NewDisabledAuthenticationManager(), nil
		}
	}
	return security.NewDeferredAuthenticationManager(deferFunc), nil
}

type ReloadableAccountServiceFactory [T crypto.PrivateKey, P crypto.PublicKey] struct {
	*basePlugin[T,P]
}

func (f *ReloadableAccountServiceFactory[T,P]) Create() (account.Service, error) {
	am := &account.ReloadableService{
		DispenseFunc: func() (account.Service, error) {
			raw, err := f.dispense(account.ConnectorName)
			if err != nil {
				return nil, err
			}
			return raw.(account.Service), nil
		},
	}

	return am, nil
}

type QLightTokenManagerPluginTemplate [T crypto.PrivateKey, P crypto.PublicKey] struct {
	*basePlugin[T,P]
}

func (p *QLightTokenManagerPluginTemplate[T,P]) Get() (qlight.PluginTokenManager, error) {
	return &qlight.ReloadablePluginTokenManager{
		DeferFunc: func() (qlight.PluginTokenManager, error) {
			raw, err := p.dispense(qlight.ConnectorName)
			if err != nil {
				return nil, err
			}
			return raw.(qlight.PluginTokenManager), nil
		},
	}, nil
}

func (p *QLightTokenManagerPluginTemplate[T,P]) ManagedPlugin() managedPlugin {
	return p
}

type QLightTokenManagerPluginTemplateInterface interface {
	Get() (qlight.PluginTokenManager, error)
	Start() (err error)
	Stop() (err error)
	ManagedPlugin() managedPlugin
}

//go:generate mockgen -source=plugin_templates.go -destination plugin_templates_mockery.go -package plugin
// var _ QLightTokenManagerPluginTemplateInterface = &QLightTokenManagerPluginTemplate{}
var _ QLightTokenManagerPluginTemplateInterface = &MockQLightTokenManagerPluginTemplateInterface{}
