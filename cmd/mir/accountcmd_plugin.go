package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/pluggable"
	"github.com/pavelkrolevets/MIR-pro/cmd/utils"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/plugin"
	"gopkg.in/urfave/cli.v1"
)

var (
	quorumAccountPluginCommands = cli.Command{
		Name:  "plugin",
		Usage: "Manage 'account' plugin accounts",
		Description: `
		geth account plugin
	
	Quorum supports alternate account management methods through the use of 'account' plugins.
	
	See docs.goquorum.com for more info. 
	`,
		Subcommands: []cli.Command{
			{
				Name:   "list",
				Usage:  "Print summary of existing 'account' plugin accounts",
				Action: utils.MigrateFlags(listPluginAccountsCLIAction[nist.PrivateKey, nist.PublicKey]),
				Flags: []cli.Flag{
					utils.PluginSettingsFlag, // flag is used implicitly by makeConfigNode()
					utils.PluginLocalVerifyFlag,
					utils.PluginPublicKeyFlag,
					utils.PluginSkipVerifyFlag,
				},
				Description: `
	geth account plugin list			
Print a short summary of all accounts for the given plugin settings`,
			},
			{
				Name:   "new",
				Usage:  "Create a new account using an 'account' plugin",
				Action: utils.MigrateFlags(createPluginAccountCLIAction[nist.PrivateKey, nist.PublicKey]),
				Flags: []cli.Flag{
					utils.PluginSettingsFlag,
					utils.PluginLocalVerifyFlag,
					utils.PluginPublicKeyFlag,
					utils.PluginSkipVerifyFlag,
					utils.AccountPluginNewAccountConfigFlag,
				},
				Description: fmt.Sprintf(`
	geth account plugin new

Creates a new account using an 'account' plugin and prints the address.

--%v and --%v flags are required.

Each 'account' plugin will have different requirements for the value of --%v.  
For more info see the documentation for the particular 'account' plugin being used.
`, utils.PluginSettingsFlag.Name, utils.AccountPluginNewAccountConfigFlag.Name, utils.AccountPluginNewAccountConfigFlag.Name),
			},
			{
				Name:   "import",
				Usage:  "Import a private key into a new account using an 'account' plugin",
				Action: utils.MigrateFlags(importPluginAccountCLIAction[nist.PrivateKey, nist.PublicKey]),
				Flags: []cli.Flag{
					utils.PluginSettingsFlag,
					utils.PluginLocalVerifyFlag,
					utils.PluginPublicKeyFlag,
					utils.PluginSkipVerifyFlag,
					utils.AccountPluginNewAccountConfigFlag,
				},
				ArgsUsage: "<keyFile>",
				Description: `
    geth account plugin import <keyfile>

Imports an unencrypted private key from <keyfile> and creates a new account using an 'account' plugin.
Prints the address.

The keyfile must contain an unencrypted private key in hexadecimal format.

--%v and --%v flags are required.				
				
Note:
Before using this import mechanism to transfer accounts that are already 'account' plugin-managed between nodes, consult 
the documentation for the particular 'account' plugin being used as it may support alternate methods for transferring.
`,
			},
		},
	}

	// supportedPlugins is the list of supported plugins for the account subcommand
	supportedPlugins = []plugin.PluginInterfaceName{plugin.AccountPluginInterfaceName}

	invalidPluginFlagsErr = fmt.Errorf("--%v and --%v flags must be set", utils.PluginSettingsFlag.Name, utils.AccountPluginNewAccountConfigFlag.Name)

	// makeConfigNodeDelegate is a wrapper for the makeConfigNode function.
	// It can be replaced with a stub for testing.
	makeConfigNodeDelegate configNodeMaker[nist.PrivateKey,nist.PublicKey] = standardConfigNodeMaker[nist.PrivateKey,nist.PublicKey]{}
)

func listPluginAccountsCLIAction[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	accts, err := listPluginAccounts[T,P](ctx)
	if err != nil {
		utils.Fatalf("%v", err)
	}

	var index int
	for _, acct := range accts {
		fmt.Printf("Account #%d: {%x} %s\n", index, acct.Address, &acct.URL)
		index++
	}

	return nil
}

func listPluginAccounts[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) ([]accounts.Account, error) {
	if !ctx.IsSet(utils.PluginSettingsFlag.Name) {
		return []accounts.Account{}, fmt.Errorf("--%v required", utils.PluginSettingsFlag.Name)
	}

	p, err := setupAccountPluginForCLI[T,P](ctx)
	if err != nil {
		return []accounts.Account{}, err
	}
	defer func() {
		if err := p.teardown(); err != nil {
			log.Error("error tearing down account plugin", "err", err)
		}
	}()

	return p.accounts(), nil
}

func createPluginAccountCLIAction[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	account, err := createPluginAccount[T,P](ctx)
	if err != nil {
		utils.Fatalf("unable to create plugin-backed account: %v", err)
	}
	writePluginAccountToStdOut(account)
	return nil
}

func createPluginAccount[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) (accounts.Account, error) {
	if !ctx.IsSet(utils.PluginSettingsFlag.Name) || !ctx.IsSet(utils.AccountPluginNewAccountConfigFlag.Name) {
		return accounts.Account{}, invalidPluginFlagsErr
	}

	newAcctCfg, err := getNewAccountConfigFromCLI(ctx)
	if err != nil {
		return accounts.Account{}, err
	}

	p, err := setupAccountPluginForCLI[T,P](ctx)
	if err != nil {
		return accounts.Account{}, err
	}
	defer func() {
		if err := p.teardown(); err != nil {
			log.Error("error tearing down account plugin", "err", err)
		}
	}()

	return p.NewAccount(newAcctCfg)
}

func importPluginAccountCLIAction[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	account, err := importPluginAccount[T,P](ctx)
	if err != nil {
		utils.Fatalf("unable to import key and create plugin-backed account: %v", err)
	}
	writePluginAccountToStdOut(account)
	return nil
}

func importPluginAccount[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) (accounts.Account, error) {
	keyfile := ctx.Args().First()
	if len(keyfile) == 0 {
		return accounts.Account{}, errors.New("keyfile must be given as argument")
	}
	key, err := crypto.LoadECDSA[T](keyfile)
	if err != nil {
		return accounts.Account{}, fmt.Errorf("Failed to load the private key: %v", err)
	}
	keyBytes := crypto.FromECDSA(key)
	keyHex := hex.EncodeToString(keyBytes)

	if !ctx.IsSet(utils.PluginSettingsFlag.Name) || !ctx.IsSet(utils.AccountPluginNewAccountConfigFlag.Name) {
		return accounts.Account{}, invalidPluginFlagsErr
	}

	newAcctCfg, err := getNewAccountConfigFromCLI(ctx)
	if err != nil {
		return accounts.Account{}, err
	}

	p, err := setupAccountPluginForCLI[T,P](ctx)
	if err != nil {
		return accounts.Account{}, err
	}
	defer func() {
		if err := p.teardown(); err != nil {
			log.Error("error tearing down account plugin", "err", err)
		}
	}()

	return p.ImportRawKey(keyHex, newAcctCfg)
}

func getNewAccountConfigFromCLI(ctx *cli.Context) (map[string]interface{}, error) {
	data := ctx.String(utils.AccountPluginNewAccountConfigFlag.Name)
	conf, err := plugin.ReadMultiFormatConfig(data)
	if err != nil {
		return nil, fmt.Errorf("invalid account creation config provided: %v", err)
	}
	// plugin backend expects config to be json map
	confMap := new(map[string]interface{})
	if err := json.Unmarshal(conf, confMap); err != nil {
		return nil, fmt.Errorf("invalid account creation config provided: %v", err)
	}
	return *confMap, nil
}

type accountPlugin [T crypto.PrivateKey, P crypto.PublicKey] struct {
	pluggable.AccountCreator
	am *accounts.Manager[P]
	pm *plugin.PluginManager[T,P]
}

func (c *accountPlugin[T,P]) teardown() error {
	return c.pm.Stop()
}

func (c *accountPlugin[T,P]) accounts() []accounts.Account {
	b := c.am.Backends(reflect.TypeOf(&pluggable.Backend[P]{}))
	if b == nil {
		return []accounts.Account{}
	}

	var accts []accounts.Account
	for _, wallet := range b[0].Wallets() {
		accts = append(accts, wallet.Accounts()...)
	}
	return accts
}

// startPluginManagerForAccountCLI is a helper func for use with the account plugin CLI.
// It creates and starts a new PluginManager with the provided CLI flags.
// The caller should call teardown on the returned accountPlugin to stop the plugin after use.
// The returned accountPlugin provides several methods necessary for the account plugin CLI, abstracting the underlying plugin/account types.
//
// This func should not be used for anything other than the account CLI.
// The account plugin, if present, is registered with the existing pluggable.Backend in the stack's AccountManager.
// This allows the AccountManager to use the account plugin even though the PluginManager is not registered with the stack.
// Instead of registering a plugin manager with the stack this is manually creating a plugin manager.
// This means that the plugin manager can be started without having to start the whole stack (P2P client, IPC interface, ...).
// The purpose of this is to help prevent issues/conflicts if an existing node is already running on this host.
//
func setupAccountPluginForCLI[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) (*accountPlugin[T,P], error) {
	stack, cfg := makeConfigNode[T,P](ctx)

	if cfg.Node.Plugins == nil {
		return nil, errors.New("no plugin config provided")
	}
	if err := cfg.Node.Plugins.CheckSettingsAreSupported(supportedPlugins); err != nil {
		return nil, err
	}
	if err := cfg.Node.ResolvePluginBaseDir(); err != nil {
		return nil, fmt.Errorf("unable to resolve plugin base dir due to %s", err)
	}

	pm, err := plugin.NewPluginManager[T,P](
		cfg.Node.UserIdent,
		cfg.Node.Plugins,
		ctx.Bool(utils.PluginSkipVerifyFlag.Name),
		ctx.Bool(utils.PluginLocalVerifyFlag.Name),
		ctx.String(utils.PluginPublicKeyFlag.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create plugin manager: %v", err)
	}
	if err := pm.Start(); err != nil {
		return nil, fmt.Errorf("unable to start plugin manager: %v", err)
	}

	b := stack.AccountManager().Backends(reflect.TypeOf(&pluggable.Backend[P]{}))[0].(*pluggable.Backend[P])
	if err := pm.AddAccountPluginToBackend(b); err != nil {
		return nil, fmt.Errorf("unable to load pluggable account backend: %v", err)
	}

	return &accountPlugin[T,P]{
		AccountCreator: b,
		am:             stack.AccountManager(),
		pm:             pm,
	}, nil
}

func writePluginAccountToStdOut(account accounts.Account) {
	fmt.Printf("\nYour new plugin-backed account was generated\n\n")
	fmt.Printf("Public address of the account:   %s\n", account.Address.Hex())
	fmt.Printf("Account URL: %s\n\n", account.URL.Path)
	fmt.Printf("- You can share your public address with anyone. Others need it to interact with you.\n")
	fmt.Printf("- You must NEVER share the secret key with anyone! The key controls access to your funds!\n")
	fmt.Printf("- Consider BACKING UP your account! The specifics of backing up will depend on the plugin backend being used.\n")
	fmt.Printf("- The plugin backend may require you to REMEMBER part/all of the new account config to retrieve the key in the future!\n  See the plugin specific documentation for more info.\n\n")
	fmt.Printf("- See the documentation for the plugin being used for more info.\n\n")
}

type configNodeMaker [T crypto.PrivateKey, P crypto.PublicKey] interface {
	makeConfigNode(ctx *cli.Context) (*node.Node[T,P], gethConfig[T,P])
}

// standardConfigNodeMaker is a wrapper around the makeConfigNode function to enable mocking in testing
type standardConfigNodeMaker [T crypto.PrivateKey, P crypto.PublicKey]struct{}

func (f standardConfigNodeMaker[T,P]) makeConfigNode(ctx *cli.Context) (*node.Node[T,P], gethConfig[T,P]) {
	return makeConfigNode[T,P](ctx)
}
