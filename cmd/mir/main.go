// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// geth is the official command-line client for Ethereum.
package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/accounts/keystore"
	"github.com/pavelkrolevets/MIR-pro/accounts/pluggable"
	"github.com/pavelkrolevets/MIR-pro/cmd/utils"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/console/prompt"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/eth"

	// "github.com/pavelkrolevets/MIR-pro/eth"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/ethclient"
	"github.com/pavelkrolevets/MIR-pro/internal/debug"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/internal/flags"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/metrics"
	"github.com/pavelkrolevets/MIR-pro/node"

	// "github.com/pavelkrolevets/MIR-pro/permission"
	"github.com/pavelkrolevets/MIR-pro/plugin"
	"gopkg.in/urfave/cli.v1"
)

const (
	clientIdentifier = "mir" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""
	// The app that holds all commands and flags.
	app = flags.NewApp(gitCommit, gitDate, "the mir-ethereum command line interface")
	// flags that configure the node
	nodeFlags = []cli.Flag{
		utils.IdentityFlag,
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.BootnodesFlag,
		utils.DataDirFlag,
		utils.AncientFlag,
		utils.MinFreeDiskSpaceFlag,
		utils.KeyStoreDirFlag,
		utils.ExternalSignerFlag,
		utils.NoUSBFlag,
		utils.USBFlag,
		utils.SmartCardDaemonPathFlag,
		utils.OverrideBerlinFlag,
		utils.EthashCacheDirFlag,
		utils.EthashCachesInMemoryFlag,
		utils.EthashCachesOnDiskFlag,
		utils.EthashCachesLockMmapFlag,
		utils.EthashDatasetDirFlag,
		utils.EthashDatasetsInMemoryFlag,
		utils.EthashDatasetsOnDiskFlag,
		utils.EthashDatasetsLockMmapFlag,
		utils.TxPoolLocalsFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
		utils.SyncModeFlag,
		utils.ExitWhenSyncedFlag,
		utils.GCModeFlag,
		utils.SnapshotFlag,
		utils.TxLookupLimitFlag,
		utils.LightServeFlag,
		utils.LightIngressFlag,
		utils.LightEgressFlag,
		utils.LightMaxPeersFlag,
		utils.LightNoPruneFlag,
		utils.LightKDFFlag,
		utils.UltraLightServersFlag,
		utils.UltraLightFractionFlag,
		utils.UltraLightOnlyAnnounceFlag,
		utils.LightNoSyncServeFlag,
		utils.AuthorizationListFlag,
		utils.BloomFilterSizeFlag,
		utils.CacheFlag,
		utils.CacheDatabaseFlag,
		utils.CacheTrieFlag,
		utils.CacheTrieJournalFlag,
		utils.CacheTrieRejournalFlag,
		utils.CacheGCFlag,
		utils.CacheSnapshotFlag,
		utils.CacheNoPrefetchFlag,
		utils.CachePreimagesFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.MiningEnabledFlag,
		utils.MinerThreadsFlag,
		utils.MinerNotifyFlag,
		utils.MinerGasTargetFlag,
		utils.MinerGasLimitFlag,
		utils.MinerGasPriceFlag,
		utils.MinerEtherbaseFlag,
		utils.MinerExtraDataFlag,
		utils.MinerRecommitIntervalFlag,
		utils.MinerNoVerfiyFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.DNSDiscoveryFlag,
		utils.MainnetMirFlag,
		utils.MainnetEthFlag,
		utils.DeveloperFlag,
		utils.DeveloperPeriodFlag,
		utils.SoyuzFlag,
		utils.RinkebyFlag,
		utils.VMEnableDebugFlag,
		utils.NetworkIdFlag,
		utils.EthStatsURLFlag,
		utils.FakePoWFlag,
		utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.GpoPercentileFlag,
		utils.GpoMaxGasPriceFlag,
		utils.EWASMInterpreterFlag,
		utils.EVMInterpreterFlag,
		utils.MinerNotifyFullFlag,
		configFileFlag,
		utils.CatalystFlag,
		// Quorum
		utils.QuorumImmutabilityThreshold,
		utils.EnableNodePermissionFlag,
		utils.RaftLogDirFlag,
		utils.RaftModeFlag,
		utils.RaftBlockTimeFlag,
		utils.RaftJoinExistingFlag,
		utils.RaftPortFlag,
		utils.RaftDNSEnabledFlag,
		utils.EmitCheckpointsFlag,
		utils.IstanbulRequestTimeoutFlag,
		utils.IstanbulBlockPeriodFlag,
		utils.PluginSettingsFlag,
		utils.PluginSkipVerifyFlag,
		utils.PluginLocalVerifyFlag,
		utils.PluginPublicKeyFlag,
		utils.AllowedFutureBlockTimeFlag,
		utils.EVMCallTimeOutFlag,
		utils.MultitenancyFlag,
		utils.RevertReasonFlag,
		utils.QuorumEnablePrivateTrieCache,
		utils.QuorumEnablePrivacyMarker,
		utils.QuorumPTMUnixSocketFlag,
		utils.QuorumPTMUrlFlag,
		utils.QuorumPTMTimeoutFlag,
		utils.QuorumPTMDialTimeoutFlag,
		utils.QuorumPTMHttpIdleTimeoutFlag,
		utils.QuorumPTMHttpWriteBufferSizeFlag,
		utils.QuorumPTMHttpReadBufferSizeFlag,
		utils.QuorumPTMTlsModeFlag,
		utils.QuorumPTMTlsRootCaFlag,
		utils.QuorumPTMTlsClientCertFlag,
		utils.QuorumPTMTlsClientKeyFlag,
		utils.QuorumPTMTlsInsecureSkipVerify,
		utils.QuorumLightServerFlag,
		utils.QuorumLightServerP2PListenPortFlag,
		utils.QuorumLightServerP2PMaxPeersFlag,
		utils.QuorumLightServerP2PNetrestrictFlag,
		utils.QuorumLightServerP2PPermissioningFlag,
		utils.QuorumLightServerP2PPermissioningPrefixFlag,
		utils.QuorumLightClientFlag,
		utils.QuorumLightClientPSIFlag,
		utils.QuorumLightClientTokenEnabledFlag,
		utils.QuorumLightClientTokenValueFlag,
		utils.QuorumLightClientTokenManagementFlag,
		utils.QuorumLightClientRPCTLSFlag,
		utils.QuorumLightClientRPCTLSInsecureSkipVerifyFlag,
		utils.QuorumLightClientRPCTLSCACertFlag,
		utils.QuorumLightClientRPCTLSCertFlag,
		utils.QuorumLightClientRPCTLSKeyFlag,
		utils.QuorumLightClientServerNodeFlag,
		utils.QuorumLightClientServerNodeRPCFlag,
		utils.QuorumLightTLSFlag,
		utils.QuorumLightTLSCertFlag,
		utils.QuorumLightTLSKeyFlag,
		utils.QuorumLightTLSCACertsFlag,
		utils.QuorumLightTLSClientAuthFlag,
		utils.QuorumLightTLSCipherSuitesFlag,
		// End-Quorum
		
		//Mir
		utils.SignerCertFlag,
		utils.CryptoSwitchFlag,
		utils.CryptoGostCurveFlag,
	}

	rpcFlags = []cli.Flag{
		utils.HTTPEnabledFlag,
		utils.HTTPListenAddrFlag,
		utils.HTTPPortFlag,
		utils.HTTPCORSDomainFlag,
		utils.HTTPVirtualHostsFlag,
		utils.LegacyRPCEnabledFlag,
		utils.LegacyRPCListenAddrFlag,
		utils.LegacyRPCPortFlag,
		utils.LegacyRPCCORSDomainFlag,
		utils.LegacyRPCVirtualHostsFlag,
		utils.LegacyRPCApiFlag,
		utils.GraphQLEnabledFlag,
		utils.GraphQLCORSDomainFlag,
		utils.GraphQLVirtualHostsFlag,
		utils.HTTPApiFlag,
		utils.HTTPPathPrefixFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.WSPathPrefixFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
		utils.InsecureUnlockAllowedFlag,
		utils.RPCGlobalGasCapFlag,
		utils.RPCGlobalTxFeeCapFlag,
		utils.AllowUnprotectedTxs,
	}

	metricsFlags = []cli.Flag{
		utils.MetricsEnabledFlag,
		utils.MetricsEnabledExpensiveFlag,
		utils.MetricsHTTPFlag,
		utils.MetricsPortFlag,
		utils.MetricsEnableInfluxDBFlag,
		utils.MetricsInfluxDBEndpointFlag,
		utils.MetricsInfluxDBDatabaseFlag,
		utils.MetricsInfluxDBUsernameFlag,
		utils.MetricsInfluxDBPasswordFlag,
		utils.MetricsInfluxDBTagsFlag,
	}
)

func init() {
	// Initialize the CLI app and start Mirc
	app.Action = mir
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2013-2023 The go-ethereum and Mir Authors"
	app.Commands = []cli.Command{
		// See chaincmd.go:
		initCommand,
		mpsdbUpgradeCommand,
		importCommand,
		exportCommand,
		importPreimagesCommand,
		exportPreimagesCommand,
		removedbCommand,
		dumpCommand,
		dumpGenesisCommand,
		// See accountcmd.go:
		accountCommand,
		walletCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		// See misccmd.go:
		makecacheCommand,
		makedagCommand,
		versionCommand,
		versionCheckCommand,
		licenseCommand,
		// See config.go
		dumpConfigCommand,
		// see dbcmd.go
		dbCommand,
		// See cmd/utils/flags_legacy.go
		utils.ShowDeprecated,
		// See snapshot.go
		snapshotCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	app.Flags = append(app.Flags, metricsFlags...)

	app.Before = func(ctx *cli.Context) error {
		return debug.Setup(ctx)
	}
	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		prompt.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// prepare manipulates memory cache allowance and setups metric system.
// This function should be called before launching devp2p stack.
func prepare(ctx *cli.Context) {
	// If we're running a known preset, log it for convenience.
	switch {
	case ctx.GlobalIsSet(utils.SoyuzFlag.Name):
		log.Info("Starting Mir on Soyuz 3 testnet...")

	case ctx.GlobalIsSet(utils.DeveloperFlag.Name):
		log.Info("Starting Mir in ephemeral dev mode...")

	case !ctx.GlobalIsSet(utils.NetworkIdFlag.Name):
		log.Info("Starting Mir on mainnet...")

	case ctx.GlobalIsSet(utils.RinkebyFlag.Name):
		log.Info("Starting Mir on Rinkeby testnet...")
	}
	// If we're a full node on mainnet without --cache specified, bump default cache allowance
	if ctx.GlobalString(utils.SyncModeFlag.Name) != "light" && !ctx.GlobalIsSet(utils.CacheFlag.Name) && !ctx.GlobalIsSet(utils.NetworkIdFlag.Name) {
		// Make sure we're not on any supported preconfigured testnet either
		if !ctx.GlobalIsSet(utils.SoyuzFlag.Name) && !ctx.GlobalIsSet(utils.DeveloperFlag.Name) {
			// Nope, we're really on mainnet. Bump that cache up!
			log.Info("Bumping default cache on mainnet", "provided", ctx.GlobalInt(utils.CacheFlag.Name), "updated", 4096)
			ctx.GlobalSet(utils.CacheFlag.Name, strconv.Itoa(4096))
		}
	}
	// If we're running a light client on any network, drop the cache to some meaningfully low amount
	if ctx.GlobalString(utils.SyncModeFlag.Name) == "light" && !ctx.GlobalIsSet(utils.CacheFlag.Name) {
		log.Info("Dropping default light client cache", "provided", ctx.GlobalInt(utils.CacheFlag.Name), "updated", 128)
		ctx.GlobalSet(utils.CacheFlag.Name, strconv.Itoa(128))
	}

	// Start metrics export if enabled
	utils.SetupMetrics(ctx)

	// Start system runtime metrics collection
	go metrics.CollectProcessMetrics(3 * time.Second)
}

// mir is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func mir(ctx *cli.Context) error {
	fmt.Println(`
	███╗   ███╗██╗██████╗      ██████╗██╗  ██╗ █████╗ ██╗███╗   ██╗
	████╗ ████║██║██╔══██╗    ██╔════╝██║  ██║██╔══██╗██║████╗  ██║
	██╔████╔██║██║██████╔╝    ██║     ███████║███████║██║██╔██╗ ██║
	██║╚██╔╝██║██║██╔══██╗    ██║     ██╔══██║██╔══██║██║██║╚██╗██║
	██║ ╚═╝ ██║██║██║  ██║    ╚██████╗██║  ██║██║  ██║██║██║ ╚████║
	╚═╝     ╚═╝╚═╝╚═╝  ╚═╝     ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝
																   `)
	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	// Mir - set crypto before the start of services
	switch {
	case ctx.GlobalIsSet(utils.CryptoSwitchFlag.Name):
		if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔═╗╔═╗╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║ ╦║ ║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╚═╝╚═╝╚═╝ ╩ `)
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetA" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetA()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetB" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetB()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-tc26-gost-3410-12-256-paramSetC" {
				gost3410.GostCurve = gost3410.CurveIdtc26gost341012256paramSetC()
			}
			if ctx.GlobalString(utils.CryptoGostCurveFlag.Name) == "id-GostR3410-2001-CryptoPro-A-ParamSet" {
				gost3410.GostCurve = gost3410.CurveIdGostR34102001CryptoProAParamSet()
			}
			stack, backend := makeFullNode[gost3410.PrivateKey,gost3410.PublicKey](ctx)
			defer stack.Close()		
			startNode(ctx, stack, backend)
			stack.Wait()
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			crypto.CryptoAlg = crypto.GOST_CSP
			// Mir - check if signer cert is loaded
			// if stack.Config().SignerCert.Bytes() == nil {
			// 	return fmt.Errorf("signer cert cant be nil")
			// }
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔╗╔╦╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║║║║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╝╚╝╩╚═╝ ╩ `)
			stack, backend := makeFullNode[nist.PrivateKey,nist.PublicKey](ctx)
			defer stack.Close()
			startNode(ctx, stack, backend)
			stack.Wait()
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			crypto.CryptoAlg = crypto.PQC
		} else {
			return fmt.Errorf("wrong crypto flag")
		}
	default:
		return fmt.Errorf("crypto flag not set")
	}
	return nil
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context, stack *node.Node[T,P], backend ethapi.Backend[T,P]) {
	log.DoEmitCheckpoints = ctx.GlobalBool(utils.EmitCheckpointsFlag.Name)
	debug.Memsize.Add("node", stack)

	// raft mode does not support --exitwhensynced
	if ctx.GlobalBool(utils.ExitWhenSyncedFlag.Name) && ctx.GlobalBool(utils.RaftModeFlag.Name) {
		utils.Fatalf("raft consensus does not support --exitwhensynced")
	}

	// Start up the node itself
	utils.StartNode(ctx, stack)

	// Now that the plugin manager has been started we register the account plugin with the corresponding account backend.  All other account management is disabled when using External Signer
	if !ctx.IsSet(utils.ExternalSignerFlag.Name) && stack.PluginManager().IsEnabled(plugin.AccountPluginInterfaceName) {
		b := stack.AccountManager().Backends(reflect.TypeOf(&pluggable.Backend[P]{}))[0].(*pluggable.Backend[P])
		if err := stack.PluginManager().AddAccountPluginToBackend(b); err != nil {
			log.Error("failed to setup account plugin", "err", err)
		}
	}

	// Unlock any account specifically requested
	unlockAccounts(ctx, stack)

	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent[P], 16)
	stack.AccountManager().Subscribe(events)

	// Create a client to interact with local geth node.
	rpcClient, err := stack.Attach()
	if err != nil {
		utils.Fatalf("Failed to attach to self: %v", err)
	}
	ethClient := ethclient.NewClient[P](rpcClient)

	// Quorum
	if ctx.GlobalBool(utils.MultitenancyFlag.Name) && !stack.PluginManager().IsEnabled(plugin.SecurityPluginInterfaceName) {
		utils.Fatalf("multitenancy requires RPC Security Plugin to be configured")
	}
	// End Quorum

	go func() {
		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				var derivationPaths []accounts.DerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPaths = append(derivationPaths, accounts.LegacyLedgerBaseDerivationPath)
				}
				derivationPaths = append(derivationPaths, accounts.DefaultBaseDerivationPath)

				event.Wallet.SelfDerive(derivationPaths, ethClient)

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	// Spawn a standalone goroutine for status synchronization monitoring,
	// close the node when synchronization is complete if user required.
	if ctx.GlobalBool(utils.ExitWhenSyncedFlag.Name) {
		go func() {
			sub := stack.EventMux().Subscribe(downloader.DoneEvent[P]{})
			defer sub.Unsubscribe()
			for {
				event := <-sub.Chan()
				if event == nil {
					continue
				}
				done, ok := event.Data.(downloader.DoneEvent[P])
				if !ok {
					continue
				}
				if timestamp := time.Unix(int64(done.Latest.Time), 0); time.Since(timestamp) < 10*time.Minute {
					log.Info("Synchronisation completed", "latestnum", done.Latest.Number, "latesthash", done.Latest.Hash(),
						"age", common.PrettyAge(timestamp))
					stack.Close()
				}
			}
		}()
	}

	// Quorum
	//
	// checking if permissions is enabled and staring the permissions service
	// if stack.Config().EnableNodePermission {
	// 	stack.Server().SetIsNodePermissioned(permission.IsNodePermissioned)
	// 	if stack.IsPermissionEnabled() {
	// 		var permissionService *permission.PermissionCtrl
	// 		if err := stack.Lifecycle(&permissionService); err != nil {
	// 			utils.Fatalf("Permission service not runnning: %v", err)
	// 		}
	// 		if err := permissionService.AfterStart(); err != nil {
	// 			utils.Fatalf("Permission service post construct failure: %v", err)
	// 		}
	// 	}
	// }

	// Start auxiliary services if enabled
	if ctx.GlobalBool(utils.MiningEnabledFlag.Name) || ctx.GlobalBool(utils.DeveloperFlag.Name) {
		// Mining only makes sense if a full Ethereum node is running
		if ctx.GlobalString(utils.SyncModeFlag.Name) == "light" {
			utils.Fatalf("Light clients do not support mining")
		}
		ethBackend, ok := backend.(*eth.EthAPIBackend[T,P])
		if !ok {
			utils.Fatalf("Ethereum service not running: %v", err)
		}
		// Set the gas price to the limits from the CLI and start mining
		gasprice := utils.GlobalBig(ctx, utils.MinerGasPriceFlag.Name)
		ethBackend.TxPool().SetGasPrice(gasprice)
		// start mining
		threads := ctx.GlobalInt(utils.MinerThreadsFlag.Name)
		if err := ethBackend.StartMining(threads); err != nil {
			utils.Fatalf("Failed to start mining: %v", err)
		}
	}

	// checks quorum features that depend on the ethereum service
	quorumValidateEthService(stack, ctx.GlobalBool(utils.RaftModeFlag.Name))
}

// unlockAccounts unlocks any account specifically requested.
func unlockAccounts[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context, stack *node.Node[T,P]) {
	var unlocks []string
	inputs := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for _, input := range inputs {
		if trimmed := strings.TrimSpace(input); trimmed != "" {
			unlocks = append(unlocks, trimmed)
		}
	}
	// Short circuit if there is no account to unlock.
	if len(unlocks) == 0 {
		return
	}
	// If insecure account unlocking is not allowed if node's APIs are exposed to external.
	// Print warning log to user and skip unlocking.
	if !stack.Config().InsecureUnlockAllowed && stack.Config().ExtRPCEnabled() {
		utils.Fatalf("Account unlock with HTTP access is forbidden!")
	}
	ks := stack.AccountManager().Backends(reflect.TypeOf(&keystore.KeyStore[T,P]{}))[0].(*keystore.KeyStore[T,P])
	passwords := utils.MakePasswordList(ctx)
	for i, account := range unlocks {
		unlockAccount(ks, account, i, passwords)
	}
}
