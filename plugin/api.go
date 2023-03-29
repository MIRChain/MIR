package plugin

import "github.com/pavelkrolevets/MIR-pro/crypto"

type PluginManagerAPI [T crypto.PrivateKey, P crypto.PublicKey] struct {
	pm *PluginManager[T,P]
}

func NewPluginManagerAPI[T crypto.PrivateKey, P crypto.PublicKey](pm *PluginManager[T,P]) *PluginManagerAPI[T,P] {
	return &PluginManagerAPI[T,P]{
		pm: pm,
	}
}

func (pmapi *PluginManagerAPI[T,P]) ReloadPlugin(name PluginInterfaceName) (bool, error) {
	return pmapi.pm.Reload(name)
}
