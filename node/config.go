// Copyright 2014 The go-ethereum Authors
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

package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/external"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/accounts/pluggable"

	// "github.com/pavelkrolevets/MIR-pro/accounts/scwallet"
	// "github.com/pavelkrolevets/MIR-pro/accounts/usbwallet"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/p2p/enr"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/plugin"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

const (
	datadirPrivateKey      = "nodekey"            // Path within the datadir to the node's private key
	datadirDefaultKeyStore = "keystore"           // Path within the datadir to the keystore
	datadirStaticNodes     = "static-nodes.json"  // Path within the datadir to the static node list
	datadirTrustedNodes    = "trusted-nodes.json" // Path within the datadir to the trusted node list
	datadirNodeDatabase    = "nodes"              // Path within the datadir to store the node infos
)

// Config represents a small collection of configuration values to fine tune the
// P2P network layer of a protocol stack. These values can be further extended by
// all registered services.
type Config [T crypto.PrivateKey, P crypto.PublicKey ] struct {
	// Name sets the instance name of the node. It must not contain the / character and is
	// used in the devp2p node identifier. The instance name of geth is "geth". If no
	// value is specified, the basename of the current executable is used.
	Name string `toml:"-"`

	// UserIdent, if set, is used as an additional component in the devp2p node identifier.
	UserIdent string `toml:",omitempty"`

	// Version should be set to the version number of the program. It is used
	// in the devp2p node identifier.
	Version string `toml:"-"`

	// DataDir is the file system folder the node should use for any data storage
	// requirements. The configured data directory will not be directly shared with
	// registered services, instead those can use utility methods to create/access
	// databases or flat files. This enables ephemeral nodes which can fully reside
	// in memory.
	DataDir string

	// RaftLogDir is the file system folder the node use for raft-state, raft-snap and
	// raft-wal folders.
	RaftLogDir string

	// Configuration of peer-to-peer networking.
	P2P p2p.Config[T, P]

	// Quorum
	QP2P *p2p.Config[T, P] `toml:",omitempty"`

	// KeyStoreDir is the file system folder that contains private keys. The directory can
	// be specified as a relative path, in which case it is resolved relative to the
	// current directory.
	//
	// If KeyStoreDir is empty, the default location is the "keystore" subdirectory of
	// DataDir. If DataDir is unspecified and KeyStoreDir is empty, an ephemeral directory
	// is created by New and destroyed when the node is stopped.
	KeyStoreDir string `toml:",omitempty"`

	// ExternalSigner specifies an external URI for a clef-type signer
	ExternalSigner string `toml:",omitempty"`

	// UseLightweightKDF lowers the memory and CPU requirements of the key store
	// scrypt KDF at the expense of security.
	UseLightweightKDF bool `toml:",omitempty"`

	// InsecureUnlockAllowed allows user to unlock accounts in unsafe http environment.
	InsecureUnlockAllowed bool `toml:",omitempty"`

	// NoUSB disables hardware wallet monitoring and connectivity.
	NoUSB bool `toml:",omitempty"`

	// USB enables hardware wallet monitoring and connectivity.
	USB bool `toml:",omitempty"`

	// SmartCardDaemonPath is the path to the smartcard daemon's socket
	SmartCardDaemonPath string `toml:",omitempty"`

	// IPCPath is the requested location to place the IPC endpoint. If the path is
	// a simple file name, it is placed inside the data directory (or on the root
	// pipe path on Windows), whereas if it's a resolvable path name (absolute or
	// relative), then that specific path is enforced. An empty path disables IPC.
	IPCPath string

	// HTTPHost is the host interface on which to start the HTTP RPC server. If this
	// field is empty, no HTTP API endpoint will be started.
	HTTPHost string

	// HTTPPort is the TCP port number on which to start the HTTP RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful
	// for ephemeral nodes).
	HTTPPort int `toml:",omitempty"`

	// HTTPCors is the Cross-Origin Resource Sharing header to send to requesting
	// clients. Please be aware that CORS is a browser enforced security, it's fully
	// useless for custom HTTP clients.
	HTTPCors []string `toml:",omitempty"`

	// HTTPVirtualHosts is the list of virtual hostnames which are allowed on incoming requests.
	// This is by default {'localhost'}. Using this prevents attacks like
	// DNS rebinding, which bypasses SOP by simply masquerading as being within the same
	// origin. These attacks do not utilize CORS, since they are not cross-domain.
	// By explicitly checking the Host-header, the server will not allow requests
	// made against the server with a malicious host domain.
	// Requests using ip address directly are not affected
	HTTPVirtualHosts []string `toml:",omitempty"`

	// HTTPModules is a list of API modules to expose via the HTTP RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	HTTPModules []string

	// HTTPTimeouts allows for customization of the timeout values used by the HTTP RPC
	// interface.
	HTTPTimeouts rpc.HTTPTimeouts

	// HTTPPathPrefix specifies a path prefix on which http-rpc is to be served.
	HTTPPathPrefix string `toml:",omitempty"`

	// WSHost is the host interface on which to start the websocket RPC server. If
	// this field is empty, no websocket API endpoint will be started.
	WSHost string

	// WSPort is the TCP port number on which to start the websocket RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful for
	// ephemeral nodes).
	WSPort int `toml:",omitempty"`

	// WSPathPrefix specifies a path prefix on which ws-rpc is to be served.
	WSPathPrefix string `toml:",omitempty"`

	// WSOrigins is the list of domain to accept websocket requests from. Please be
	// aware that the server can only act upon the HTTP request the client sends and
	// cannot verify the validity of the request header.
	WSOrigins []string `toml:",omitempty"`

	// WSModules is a list of API modules to expose via the websocket RPC interface.
	// If the module list is empty, all RPC API endpoints designated public will be
	// exposed.
	WSModules []string

	// WSExposeAll exposes all API modules via the WebSocket RPC interface rather
	// than just the public ones.
	//
	// *WARNING* Only set this if the node is running in a trusted network, exposing
	// private APIs to untrusted users is a major security risk.
	WSExposeAll bool `toml:",omitempty"`

	// GraphQLCors is the Cross-Origin Resource Sharing header to send to requesting
	// clients. Please be aware that CORS is a browser enforced security, it's fully
	// useless for custom HTTP clients.
	GraphQLCors []string `toml:",omitempty"`

	// GraphQLVirtualHosts is the list of virtual hostnames which are allowed on incoming requests.
	// This is by default {'localhost'}. Using this prevents attacks like
	// DNS rebinding, which bypasses SOP by simply masquerading as being within the same
	// origin. These attacks do not utilize CORS, since they are not cross-domain.
	// By explicitly checking the Host-header, the server will not allow requests
	// made against the server with a malicious host domain.
	// Requests using ip address directly are not affected
	GraphQLVirtualHosts []string `toml:",omitempty"`

	// Logger is a custom logger to use with the p2p.Server.
	Logger log.Logger `toml:",omitempty"`

	staticNodesWarning     bool
	trustedNodesWarning    bool
	oldGethResourceWarning bool

	// AllowUnprotectedTxs allows non EIP-155 protected transactions to be send over RPC.
	AllowUnprotectedTxs bool `toml:",omitempty"`

	// Quorum
	Plugins              *plugin.Settings `toml:",omitempty"`
	EnableNodePermission bool             `toml:",omitempty"` // comes from EnableNodePermissionFlag --permissioned.
	EnableMultitenancy   bool             `toml:",omitempty"` // comes from MultitenancyFlag flag
}

// IPCEndpoint resolves an IPC endpoint based on a configured value, taking into
// account the set data folders as well as the designated platform we're currently
// running on.
func (c *Config[T, P]) IPCEndpoint() string {
	// Short circuit if IPC has not been enabled
	if c.IPCPath == "" {
		return ""
	}
	// On windows we can only use plain top-level pipes
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(c.IPCPath, `\\.\pipe\`) {
			return c.IPCPath
		}
		return `\\.\pipe\` + c.IPCPath
	}
	// Resolve names into the data directory full paths otherwise
	if filepath.Base(c.IPCPath) == c.IPCPath {
		if c.DataDir == "" {
			return filepath.Join(os.TempDir(), c.IPCPath)
		}
		return filepath.Join(c.DataDir, c.IPCPath)
	}
	return c.IPCPath
}

// NodeDB returns the path to the discovery node database.
func (c *Config[T, P]) NodeDB() string {
	if c.DataDir == "" {
		return "" // ephemeral
	}
	return c.ResolvePath(datadirNodeDatabase)
}

func (c *Config[T, P]) QNodeDB() string {
	if c.DataDir == "" {
		return "" // ephemeral
	}
	return c.ResolvePath("qnodes")
}

// DefaultIPCEndpoint returns the IPC path used by default.
func DefaultIPCEndpoint[T crypto.PrivateKey, P crypto.PublicKey](clientIdentifier string) string {
	if clientIdentifier == "" {
		clientIdentifier = strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
		if clientIdentifier == "" {
			panic("empty executable name")
		}
	}
	config := &Config[T, P]{DataDir: DefaultDataDir(), IPCPath: clientIdentifier + ".ipc"}
	return config.IPCEndpoint()
}

// HTTPEndpoint resolves an HTTP endpoint based on the configured host interface
// and port parameters.
func (c *Config[T, P]) HTTPEndpoint() string {
	if c.HTTPHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.HTTPHost, c.HTTPPort)
}

// DefaultHTTPEndpoint returns the HTTP endpoint used by default.
func DefaultHTTPEndpoint[T crypto.PrivateKey, P crypto.PublicKey]() string {
	config := &Config[T, P]{HTTPHost: DefaultHTTPHost, HTTPPort: DefaultHTTPPort}
	return config.HTTPEndpoint()
}

// WSEndpoint resolves a websocket endpoint based on the configured host interface
// and port parameters.
func (c *Config[T,P]) WSEndpoint() string {
	if c.WSHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.WSHost, c.WSPort)
}

// DefaultWSEndpoint returns the websocket endpoint used by default.
func DefaultWSEndpoint[T crypto.PrivateKey, P crypto.PublicKey]() string {
	config := &Config[T,P]{WSHost: DefaultWSHost, WSPort: DefaultWSPort}
	return config.WSEndpoint()
}

// ExtRPCEnabled returns the indicator whether node enables the external
// RPC(http, ws or graphql).
func (c *Config[T,P]) ExtRPCEnabled() bool {
	return c.HTTPHost != "" || c.WSHost != ""
}

// NodeName returns the devp2p node identifier.
func (c *Config[T,P]) NodeName() string {
	name := c.name()
	// Backwards compatibility: previous versions used title-cased "Geth", keep that.
	if name == "mir" || name == "mir-testnet" {
		name = "Mir"
	}
	if c.UserIdent != "" {
		name += "/" + c.UserIdent
	}
	if c.Version != "" {
		name += "/v" + c.Version
	}
	name += "/" + runtime.GOOS + "-" + runtime.GOARCH
	name += "/" + runtime.Version()
	return name
}

func (c *Config[T,P]) name() string {
	if c.Name == "" {
		progname := strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
		if progname == "" {
			panic("empty executable name, set Config.Name")
		}
		return progname
	}
	return c.Name
}

// These resources are resolved differently for "geth" instances.
var isOldGethResource = map[string]bool{
	"chaindata":          true,
	"nodes":              true,
	"nodekey":            true,
	"static-nodes.json":  false, // no warning for these because they have their
	"trusted-nodes.json": false, // own separate warning.
}

// ResolvePath resolves path in the instance directory.
func (c *Config[T,P]) ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if c.DataDir == "" {
		return ""
	}
	// Backwards-compatibility: ensure that data directory files created
	// by geth 1.4 are used if they exist.
	if warn, isOld := isOldGethResource[path]; isOld {
		oldpath := ""
		if c.name() == "geth" {
			oldpath = filepath.Join(c.DataDir, path)
		}
		if oldpath != "" && common.FileExist(oldpath) {
			if warn {
				c.warnOnce(&c.oldGethResourceWarning, "Using deprecated resource file %s, please move this file to the 'geth' subdirectory of datadir.", oldpath)
			}
			return oldpath
		}
	}
	return filepath.Join(c.instanceDir(), path)
}

func (c *Config[T,P]) instanceDir() string {
	if c.DataDir == "" {
		return ""
	}
	return filepath.Join(c.DataDir, c.name())
}

// NodeKey retrieves the currently configured private key of the node, checking
// first any manually set key, falling back to the one found in the configured
// data folder. If no key can be found, a new one is generated.
func (c *Config[T,P]) NodeKey() T {
	// Use any specifically configured key.
	if !reflect.ValueOf(&c.P2P.PrivateKey).Elem().IsZero(){
		return c.P2P.PrivateKey
	}
	// Generate ephemeral key if no datadir is being used.
	if c.DataDir == "" {
		key, err := crypto.GenerateKey[T]()
		if err != nil {
			log.Crit(fmt.Sprintf("Failed to generate ephemeral node key: %v", err))
		}
		return key
	}
	keyfile := c.ResolvePath(datadirPrivateKey)
	if key, err := crypto.LoadECDSA[T](keyfile); err == nil {
		return key
	}
	// No persistent key found, generate and store a new one.
	key, err := crypto.GenerateKey[T]()
	if err != nil {
		log.Crit(fmt.Sprintf("Failed to generate node key: %v", err))
	}
	instanceDir := filepath.Join(c.DataDir, c.name())
	if err := os.MkdirAll(instanceDir, 0700); err != nil {
		log.Error(fmt.Sprintf("Failed to persist node key: %v", err))
		return key
	}
	keyfile = filepath.Join(instanceDir, datadirPrivateKey)
	if err := crypto.SaveECDSA(keyfile, key); err != nil {
		log.Error(fmt.Sprintf("Failed to persist node key: %v", err))
	}
	return key
}

// StaticNodes returns a list of node enode URLs configured as static nodes.
func (c *Config[T,P]) StaticNodes() []*enode.Node[P] {
	return c.parsePersistentNodes(&c.staticNodesWarning, c.ResolvePath(datadirStaticNodes))
}

// TrustedNodes returns a list of node enode URLs configured as trusted nodes.
func (c *Config[T,P]) TrustedNodes() []*enode.Node[P] {
	return c.parsePersistentNodes(&c.trustedNodesWarning, c.ResolvePath(datadirTrustedNodes))
}

// parsePersistentNodes parses a list of discovery node URLs loaded from a .json
// file from within the data directory.
func (c *Config[T,P]) parsePersistentNodes(w *bool, path string) []*enode.Node[P] {
	// Short circuit if no node config is present
	if c.DataDir == "" {
		return nil
	}
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	c.warnOnce(w, "Found deprecated node list file %s, please use the TOML config file instead.", path)

	// Load the nodes from the config file.
	var nodelist []string
	if err := common.LoadJSON(path, &nodelist); err != nil {
		log.Error(fmt.Sprintf("Can't load node list file: %v", err))
		return nil
	}
	// Interpret the list as a discovery node array
	var nodes []*enode.Node[P]
	for _, url := range nodelist {
		if url == "" {
			continue
		}
		node, err := enode.Parse[P](enr.SchemeMap{"v4": enode.V4ID[P]{}}, url)
		if err != nil {
			log.Error(fmt.Sprintf("Node URL %s: %v\n", url, err))
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// AccountConfig determines the settings for scrypt and keydirectory
func (c *Config[T,P]) AccountConfig() (int, int, string, error) {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if c.UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}

	var (
		keydir string
		err    error
	)
	switch {
	case filepath.IsAbs(c.KeyStoreDir):
		keydir = c.KeyStoreDir
	case c.DataDir != "":
		if c.KeyStoreDir == "" {
			keydir = filepath.Join(c.DataDir, datadirDefaultKeyStore)
		} else {
			keydir, err = filepath.Abs(c.KeyStoreDir)
		}
	case c.KeyStoreDir != "":
		keydir, err = filepath.Abs(c.KeyStoreDir)
	}
	return scryptN, scryptP, keydir, err
}

// Quorum
//
// Make sure plugin base dir exists
func (c *Config[T,P]) ResolvePluginBaseDir() error {
	if c.Plugins == nil {
		return nil
	}
	baseDir := c.Plugins.BaseDir.String()
	if baseDir == "" {
		baseDir = filepath.Join(c.DataDir, "plugins")
	}
	if !common.FileExist(baseDir) {
		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return err
		}
	}
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return err
	}
	c.Plugins.BaseDir = plugin.EnvironmentAwaredValue(absBaseDir)
	return nil
}

// check if smart-contract-based permissioning is enabled by reading `--permissioned` flag and checking permission config file
func (c *Config[T,P]) IsPermissionEnabled() bool {
	fullPath := filepath.Join(c.DataDir, params.PERMISSION_MODEL_CONFIG)
	if _, err := os.Stat(fullPath); err != nil {
		log.Warn(fmt.Sprintf("%s file is missing. Smart-contract-based permission service will be disabled", params.PERMISSION_MODEL_CONFIG), "error", err)
		return false
	}
	return true
}

func makeAccountManager [T crypto.PrivateKey, P crypto.PublicKey](conf *Config[T,P]) (*accounts.Manager[P], string, error) {
	scryptN, scryptP, keydir, err := conf.AccountConfig()
	var ephemeral string
	if keydir == "" {
		// There is no datadir.
		keydir, err = ioutil.TempDir("", "go-ethereum-keystore")
		ephemeral = keydir
	}

	if err != nil {
		return nil, "", err
	}
	if err := os.MkdirAll(keydir, 0700); err != nil {
		return nil, "", err
	}
	// Assemble the account manager and supported backends
	var backends []accounts.Backend[P]
	if len(conf.ExternalSigner) > 0 {
		log.Info("Using external signer", "url", conf.ExternalSigner)
		if extapi, err := external.NewExternalBackend[P](conf.ExternalSigner); err == nil {
			backends = append(backends, extapi)
		} else {
			return nil, "", fmt.Errorf("error connecting to external signer: %v", err)
		}
	}
	if len(backends) == 0 {
		// For now, we're using EITHER external signer OR local signers.
		// If/when we implement some form of lockfile for USB and keystore wallets,
		// we can have both, but it's very confusing for the user to see the same
		// accounts in both externally and locally, plus very racey.
		backends = append(backends, keystore.NewKeyStore[T,P](keydir, scryptN, scryptP))
		// if conf.USB {
		// 	// Start a USB hub for Ledger hardware wallets
		// 	if ledgerhub, err := usbwallet.NewLedgerHub(); err != nil {
		// 		log.Warn(fmt.Sprintf("Failed to start Ledger hub, disabling: %v", err))
		// 	} else {
		// 		backends = append(backends, ledgerhub)
		// 	}
		// 	// Start a USB hub for Trezor hardware wallets (HID version)
		// 	if trezorhub, err := usbwallet.NewTrezorHubWithHID(); err != nil {
		// 		log.Warn(fmt.Sprintf("Failed to start HID Trezor hub, disabling: %v", err))
		// 	} else {
		// 		backends = append(backends, trezorhub)
		// 	}
		// 	// Start a USB hub for Trezor hardware wallets (WebUSB version)
		// 	if trezorhub, err := usbwallet.NewTrezorHubWithWebUSB(); err != nil {
		// 		log.Warn(fmt.Sprintf("Failed to start WebUSB Trezor hub, disabling: %v", err))
		// 	} else {
		// 		backends = append(backends, trezorhub)
		// 	}
		// }
		// if len(conf.SmartCardDaemonPath) > 0 {
		// 	// Start a smart card hub
		// 	if schub, err := scwallet.NewHub(conf.SmartCardDaemonPath, scwallet.Scheme, keydir); err != nil {
		// 		log.Warn(fmt.Sprintf("Failed to start smart card hub, disabling: %v", err))
		// 	} else {
		// 		backends = append(backends, schub)
		// 	}
		// }
		if conf.Plugins != nil {
			if _, ok := conf.Plugins.Providers[plugin.AccountPluginInterfaceName]; ok {
				pluginBackend := pluggable.NewBackend[P]()
				backends = append(backends, pluginBackend)
			}
		}
	}

	return accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: conf.InsecureUnlockAllowed}, backends...), ephemeral, nil
}

var warnLock sync.Mutex

func (c *Config[T,P]) warnOnce(w *bool, format string, args ...interface{}) {
	warnLock.Lock()
	defer warnLock.Unlock()

	if *w {
		return
	}
	l := c.Logger
	if l == nil {
		l = log.Root()
	}
	l.Warn(fmt.Sprintf(format, args...))
	*w = true
}