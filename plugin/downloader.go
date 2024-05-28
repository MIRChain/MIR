package plugin

import (
	"fmt"
	"path"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto"
	"github.com/MIRChain/MIR/log"
)

// get plugin zip file from local or remote
type Downloader[T crypto.PrivateKey, P crypto.PublicKey] struct {
	pm *PluginManager[T, P]
}

func NewDownloader[T crypto.PrivateKey, P crypto.PublicKey](pm *PluginManager[T, P]) *Downloader[T, P] {
	return &Downloader[T, P]{
		pm: pm,
	}
}

func (d *Downloader[T, P]) Download(definition *PluginDefinition) (string, error) {
	// check if plugin is already in the local
	pluginFile := path.Join(d.pm.pluginBaseDir, definition.DistFileName())
	exist := common.FileExist(pluginFile)
	log.Debug("checking plugin zip file", "path", pluginFile, "exist", exist)
	if exist {
		return pluginFile, nil
	}
	if err := d.pm.centralClient.PluginDistribution(definition, pluginFile); err != nil {
		return "", fmt.Errorf("can't download from Plugin Central due to: %s. Please download the plugin manually and copy it to %s", err, d.pm.pluginBaseDir)
	}
	return pluginFile, nil
}
