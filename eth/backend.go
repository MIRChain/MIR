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

// Package eth implements the Ethereum protocol.
package eth

import (
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/accounts"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/hexutil"
	"github.com/pavelkrolevets/MIR-pro/consensus"
	"github.com/pavelkrolevets/MIR-pro/consensus/clique"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/bloombits"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state/pruner"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/core/vm"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/eth/downloader"
	"github.com/pavelkrolevets/MIR-pro/eth/ethconfig"
	"github.com/pavelkrolevets/MIR-pro/eth/filters"
	"github.com/pavelkrolevets/MIR-pro/eth/gasprice"
	"github.com/pavelkrolevets/MIR-pro/eth/protocols/eth"
	qlightproto "github.com/pavelkrolevets/MIR-pro/eth/protocols/qlight"
	"github.com/pavelkrolevets/MIR-pro/eth/protocols/snap"
	"github.com/pavelkrolevets/MIR-pro/ethdb"
	"github.com/pavelkrolevets/MIR-pro/event"
	"github.com/pavelkrolevets/MIR-pro/internal/ethapi"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/miner"
	"github.com/pavelkrolevets/MIR-pro/node"
	"github.com/pavelkrolevets/MIR-pro/p2p"
	"github.com/pavelkrolevets/MIR-pro/p2p/dnsdisc"
	"github.com/pavelkrolevets/MIR-pro/p2p/enode"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/plugin"
	"github.com/pavelkrolevets/MIR-pro/plugin/security"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/qlight"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"github.com/pavelkrolevets/MIR-pro/rpc"
)

// Config contains the configuration options of the ETH protocol.
// Deprecated: use ethconfig.Config instead.
// type Config[P crypto.PublicKey] ethconfig.Config[P]

// Ethereum implements the Ethereum full node service.
type Ethereum [T crypto.PrivateKey, P crypto.PublicKey] struct {
	config *ethconfig.Config[P]

	// Handlers
	txPool             *core.TxPool[P]
	blockchain         *core.BlockChain[P]
	handler            *handler[T,P]
	ethDialCandidates  enode.Iterator[P]
	snapDialCandidates enode.Iterator[P]

	// DB interfaces
	chainDb ethdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine[P]
	accountManager *accounts.Manager[P]

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *core.ChainIndexer[P]          // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

	APIBackend *EthAPIBackend[T,P]

	miner     *miner.Miner[T,P]
	gasPrice  *big.Int
	etherbase common.Address

	networkID     uint64
	netRPCService *ethapi.PublicNetAPI[T,P]

	p2pServer *p2p.Server[T,P]

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)

	// Quorum - consensus as eth-service (e.g. raft)
	consensusServicePendingLogsFeed *event.Feed
	qlightServerHandler             *handler[T,P]
	qlightP2pServer                 *p2p.Server[T,P]
	qlightTokenHolder               *qlight.TokenHolder[T,P]
}

// New creates a new Ethereum object (including the
// initialisation of the common Ethereum object)
func New[T crypto.PrivateKey, P crypto.PublicKey](stack *node.Node[T,P], config *ethconfig.Config[P]) (*Ethereum[T,P], error) {
	// Ensure configuration values are compatible and sane
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run eth.Ethereum in light sync mode, use les.LightEthereum")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	if config.Miner.GasPrice == nil || config.Miner.GasPrice.Cmp(common.Big0) <= 0 {
		log.Warn("Sanitizing invalid miner gas price", "provided", config.Miner.GasPrice, "updated", ethconfig.Defaults[P]().Miner.GasPrice)
		config.Miner.GasPrice = new(big.Int).Set(ethconfig.Defaults[P]().Miner.GasPrice)
	}
	if config.NoPruning && config.TrieDirtyCache > 0 {
		if config.SnapshotCache > 0 {
			config.TrieCleanCache += config.TrieDirtyCache * 3 / 5
			config.SnapshotCache += config.TrieDirtyCache * 2 / 5
		} else {
			config.TrieCleanCache += config.TrieDirtyCache
		}
		config.TrieDirtyCache = 0
	}
	log.Info("Allocated trie memory caches", "clean", common.StorageSize(config.TrieCleanCache)*1024*1024, "dirty", common.StorageSize(config.TrieDirtyCache)*1024*1024)

	// Transfer mining-related config to the ethash config.
	ethashConfig := config.Ethash
	ethashConfig.NotifyFull = config.Miner.NotifyFull

	// Assemble the Ethereum object
	chainDb, err := stack.OpenDatabaseWithFreezer("chaindata", config.DatabaseCache, config.DatabaseHandles, config.DatabaseFreezer, "eth/db/chaindata/", false)
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlockWithOverride(chainDb, config.Genesis, config.OverrideBerlin)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	if err := pruner.RecoverPruning[P](stack.ResolvePath(""), chainDb, stack.ResolvePath(config.TrieCleanCacheJournal)); err != nil {
		log.Error("Failed to recover state", "error", err)
	}
	if chainConfig.IsQuorum {
		// changes to manipulate the chain id for migration from 2.0.2 and below version to 2.0.3
		// version of Quorum  - this is applicable for v2.0.3 onwards
		if (chainConfig.ChainID != nil && chainConfig.ChainID.Int64() == 1) || config.NetworkId == 1 {
			return nil, errors.New("Cannot have chain id or network id as 1.")
		}

		if config.QuorumChainConfig.PrivacyMarkerEnabled() && chainConfig.PrivacyPrecompileBlock == nil {
			return nil, errors.New("Privacy marker transactions require privacyPrecompileBlock to be set in genesis.json")
		}
		if chainConfig.Istanbul != nil && (chainConfig.IBFT != nil || chainConfig.QBFT != nil) {
			return nil, errors.New("the attributes config.Istanbul and config.[IBFT|QBFT] are mutually exclusive on the genesis file")
		}
		if chainConfig.IBFT != nil && chainConfig.QBFT != nil {
			return nil, errors.New("the attributes config.IBFT and config.QBFT are mutually exclusive on the genesis file")
		}
	}

	if !rawdb.GetIsQuorumEIP155Activated(chainDb) && chainConfig.ChainID != nil {
		//Upon starting the node, write the flag to disallow changing ChainID/EIP155 block after HF
		rawdb.WriteQuorumEIP155Activation(chainDb)
	}

	eth := &Ethereum[T,P]{
		config:            config,
		chainDb:           chainDb,
		eventMux:          stack.EventMux(),
		accountManager:    stack.AccountManager(),
		engine:            ethconfig.CreateConsensusEngine(stack, chainConfig, config, config.Miner.Notify, config.Miner.Noverify, chainDb),
		closeBloomHandler: make(chan struct{}),
		networkID:         config.NetworkId,
		gasPrice:          config.Miner.GasPrice,
		etherbase:         config.Miner.Etherbase,
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		bloomIndexer:      core.NewBloomIndexer[P](chainDb, params.BloomBitsBlocks, params.BloomConfirms),
		p2pServer:         stack.Server(),
		// Quorum
		qlightP2pServer:                 stack.QServer(),
		consensusServicePendingLogsFeed: new(event.Feed),
	}

	// Quorum: Set protocol Name/Version
	// keep `var protocolName = "eth"` as is, and only update the quorum consensus specific protocol
	// This is used to enable the eth service to return multiple devp2p subprotocols.
	// Previously, for istanbul/64 istnbul/99 and clique (v2.6) `protocolName` would be overridden and
	// set to the consensus subprotocol name instead of "eth", meaning the node would no longer
	// communicate over the "eth" subprotocol, e.g. "eth" or "istanbul/99" but not eth" and "istanbul/99".
	// With this change, support is added so that the "eth" subprotocol remains and optionally a consensus subprotocol
	// can be added allowing the node to communicate over "eth" and an optional consensus subprotocol, e.g. "eth" and "istanbul/100"
	if chainConfig.IsQuorum {
		quorumProtocol := eth.engine.Protocol()
		// set the quorum specific consensus devp2p subprotocol, eth subprotocol remains set to protocolName as in upstream geth.
		quorumConsensusProtocolName = quorumProtocol.Name
		quorumConsensusProtocolVersions = quorumProtocol.Versions
		quorumConsensusProtocolLengths = quorumProtocol.Lengths
	}

	// force to set the istanbul etherbase to node key address
	if chainConfig.Istanbul != nil || chainConfig.IBFT != nil || chainConfig.QBFT != nil {
		key:= stack.GetNodeKey()
		var pub P
		switch t:=any(&key).(type) {
		case *nist.PrivateKey:
			p:=any(&pub).(*nist.PublicKey)
			*p = *t.Public()
		case *gost3410.PrivateKey:
			p:=any(&pub).(*gost3410.PublicKey)
			*p = *t.Public()
		}
		eth.etherbase = crypto.PubkeyToAddress(pub)
	}
	bcVersion := rawdb.ReadDatabaseVersion(chainDb)
	var dbVer = "<nil>"
	if bcVersion != nil {
		dbVer = fmt.Sprintf("%d", *bcVersion)
	}
	log.Info("Initialising Ethereum protocol", "network", config.NetworkId, "dbversion", dbVer)
	if chainConfig.IsQuorum {
		log.Info("Initialising Quorum consensus protocol", "name", quorumConsensusProtocolName, "versions", quorumConsensusProtocolVersions, "network", config.NetworkId, "dbversion", dbVer)
	}

	if !config.SkipBcVersionCheck {
		if bcVersion != nil && *bcVersion > core.BlockChainVersion {
			return nil, fmt.Errorf("database version is v%d, Geth %s only supports v%d", *bcVersion, params.VersionWithMeta, core.BlockChainVersion)
		} else if bcVersion == nil || *bcVersion < core.BlockChainVersion {
			if bcVersion != nil { // only print warning on upgrade, not on init
				log.Warn("Upgrade blockchain database version", "from", dbVer, "to", core.BlockChainVersion)
			}
			rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
		}
	}
	var (
		vmConfig = vm.Config[P]{
			EnablePreimageRecording: config.EnablePreimageRecording,
			EWASMInterpreter:        config.EWASMInterpreter,
			EVMInterpreter:          config.EVMInterpreter,
		}
		cacheConfig = &core.CacheConfig{
			TrieCleanLimit:      config.TrieCleanCache,
			TrieCleanJournal:    stack.ResolvePath(config.TrieCleanCacheJournal),
			TrieCleanRejournal:  config.TrieCleanCacheRejournal,
			TrieCleanNoPrefetch: config.NoPrefetch,
			TrieDirtyLimit:      config.TrieDirtyCache,
			TrieDirtyDisabled:   config.NoPruning,
			TrieTimeLimit:       config.TrieTimeout,
			SnapshotLimit:       config.SnapshotCache,
			Preimages:           config.Preimages,
		}
	)
	newBlockChainFunc := core.NewBlockChain[P]
	if config.QuorumChainConfig.MultiTenantEnabled() {
		newBlockChainFunc = core.NewMultitenantBlockChain[P]
	}
	eth.blockchain, err = newBlockChainFunc(chainDb, cacheConfig, chainConfig, eth.engine, vmConfig, eth.shouldPreserve, &config.TxLookupLimit, &config.QuorumChainConfig)
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			log.Error("panic occurred", "err", p)
			err := eth.Stop()
			if err != nil {
				log.Error("error while closing", "err", err)
			}
			os.Exit(1)
		}
	}()
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		eth.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	eth.bloomIndexer.Start(eth.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = stack.ResolvePath(config.TxPool.Journal)
	}
	eth.txPool = core.NewTxPool[P](config.TxPool, chainConfig, eth.blockchain)

	// Permit the downloader to use the trie cache allowance during fast sync
	cacheLimit := cacheConfig.TrieCleanLimit + cacheConfig.TrieDirtyLimit + cacheConfig.SnapshotLimit
	checkpoint := config.Checkpoint
	if checkpoint == nil {
		checkpoint = params.TrustedCheckpoints[genesisHash]
	}

	if eth.config.QuorumLightClient.Enabled() {
		clientCache, err := qlight.NewClientCache(chainDb)
		if eth.config.QuorumLightClient.TokenEnabled {
			switch eth.config.QuorumLightClient.TokenManagement {
			case "client-security-plugin":
				log.Info("Starting qlight client with auth token enabled without external API and token from argument, plugin has to be provided")
				eth.qlightTokenHolder, err = qlight.NewTokenHolder(config.QuorumLightClient.PSI, stack.PluginManager())
				if err != nil {
					return nil, fmt.Errorf("new token holder: %w", err)
				}
				eth.qlightTokenHolder.SetCurrentToken(eth.config.QuorumLightClient.TokenValue)
			case "none":
				log.Warn("Starting qlight client with auth token enabled but without a token management strategy. This is for development purposes only.")
				eth.qlightTokenHolder, err = qlight.NewTokenHolder[T,P](config.QuorumLightClient.PSI, nil)
				if err != nil {
					return nil, fmt.Errorf("new token holder: %w", err)
				}
				eth.qlightTokenHolder.SetCurrentToken(eth.config.QuorumLightClient.TokenValue)
			case "external":
				log.Info("Starting qlight client with auth token enabled and `external` token management strategy.")
				eth.qlightTokenHolder, err = qlight.NewTokenHolder[T,P](config.QuorumLightClient.PSI, nil)
				if err != nil {
					return nil, fmt.Errorf("new token holder: %w", err)
				}
			default:
				return nil, fmt.Errorf("Invalid value %s for `qlight.client.token.management`", eth.config.QuorumLightClient.TokenManagement)
			}
		}
		if err != nil {
			return nil, err
		}
		if eth.handler, err = newQLightClientHandler[T,P](&handlerConfig[T,P]{
			Database:           chainDb,
			Chain:              eth.blockchain,
			TxPool:             eth.txPool,
			Network:            config.NetworkId,
			Sync:               config.SyncMode,
			BloomCache:         uint64(cacheLimit),
			EventMux:           eth.eventMux,
			Checkpoint:         checkpoint,
			AuthorizationList:  config.AuthorizationList,
			RaftMode:           config.RaftMode,
			Engine:             eth.engine,
			psi:                config.QuorumLightClient.PSI,
			privateClientCache: clientCache,
			tokenHolder:        eth.qlightTokenHolder,
		}); err != nil {
			return nil, err
		}
		if eth.qlightTokenHolder != nil {
			eth.qlightTokenHolder.SetPeerUpdater(eth.handler.peers)
		}
		eth.blockchain.SetPrivateStateRootHashValidator(clientCache)
	} else {
		if eth.handler, err = newHandler[T,P](&handlerConfig[T,P]{
			Database:          chainDb,
			Chain:             eth.blockchain,
			TxPool:            eth.txPool,
			Network:           config.NetworkId,
			Sync:              config.SyncMode,
			BloomCache:        uint64(cacheLimit),
			EventMux:          eth.eventMux,
			Checkpoint:        checkpoint,
			AuthorizationList: config.AuthorizationList,
			RaftMode:          config.RaftMode,
			Engine:            eth.engine,
			tokenHolder:       eth.qlightTokenHolder,
		}); err != nil {
			return nil, err
		}
		if eth.config.QuorumLightServer {
			authManProvider := func() security.AuthenticationManager {
				_, authManager, _ := stack.GetSecuritySupports()
				return authManager
			}
			if eth.qlightServerHandler, err = newQLightServerHandler[T,P](&handlerConfig[T,P]{
				Database:                 chainDb,
				Chain:                    eth.blockchain,
				TxPool:                   eth.txPool,
				Network:                  config.NetworkId,
				Sync:                     config.SyncMode,
				BloomCache:               uint64(cacheLimit),
				EventMux:                 eth.eventMux,
				Checkpoint:               checkpoint,
				AuthorizationList:        config.AuthorizationList,
				RaftMode:                 config.RaftMode,
				Engine:                   eth.engine,
				authProvider:             qlight.NewAuthProvider(eth.blockchain.PrivateStateManager(), authManProvider),
				privateBlockDataResolver: qlight.NewPrivateBlockDataResolver(eth.blockchain.PrivateStateManager(), private.Ptm),
			}); err != nil {
				return nil, err
			}
		}
	}

	eth.miner = miner.New[T,P](eth, &config.Miner, chainConfig, eth.EventMux(), eth.engine, eth.isLocalBlock)
	eth.miner.SetExtra(makeExtraData(config.Miner.ExtraData, eth.blockchain.Config().IsQuorum))
	key := stack.GetNodeKey()
	var pub P
	switch t:=any(&key).(type) {
	case *nist.PrivateKey:
		p:=any(&pub).(*nist.PublicKey)
		*p = *t.Public()
	case *gost3410.PrivateKey:
		p:=any(&pub).(*gost3410.PublicKey)
		*p = *t.Public()
	}
	hexNodeId := fmt.Sprintf("%x", crypto.FromECDSAPub(pub)[1:]) // Quorum
	if eth.config.QuorumLightClient.Enabled() {
		var (
			proxyClient *rpc.Client
			err         error
		)
		// setup rpc client TLS context
		if eth.config.QuorumLightClient.RPCTLS {
			tlsConfig, err := qlight.NewTLSConfig(&qlight.TLSConfig{
				InsecureSkipVerify: eth.config.QuorumLightClient.RPCTLSInsecureSkipVerify,
				CACertFileName:     eth.config.QuorumLightClient.RPCTLSCACert,
				CertFileName:       eth.config.QuorumLightClient.RPCTLSCert,
				KeyFileName:        eth.config.QuorumLightClient.RPCTLSKey,
			})
			if err != nil {
				return nil, err
			}
			customHttpClient := &http.Client{
				Transport: http.DefaultTransport,
			}
			customHttpClient.Transport.(*http.Transport).TLSClientConfig = tlsConfig
			proxyClient, err = rpc.DialHTTPWithClient(eth.config.QuorumLightClient.ServerNodeRPC, customHttpClient)
			if err != nil {
				return nil, err
			}
		} else {
			proxyClient, err = rpc.Dial(eth.config.QuorumLightClient.ServerNodeRPC)
			if err != nil {
				return nil, err
			}
		}

		if eth.config.QuorumLightClient.TokenEnabled {
			proxyClient = proxyClient.WithHTTPCredentials(eth.qlightTokenHolder.HttpCredentialsProvider)
		}

		if len(eth.config.QuorumLightClient.PSI) > 0 {
			proxyClient = proxyClient.WithPSI(types.PrivateStateIdentifier(eth.config.QuorumLightClient.PSI))
		}
		// TODO qlight - need to find a better way to inject the rpc client into the tx manager
		rpcClientSetter, ok := private.Ptm.(private.HasRPCClient)
		if ok {
			rpcClientSetter.SetRPCClient(proxyClient)
		}
		eth.APIBackend = &EthAPIBackend[T,P]{stack.Config().ExtRPCEnabled(), stack.Config().AllowUnprotectedTxs, eth, nil, hexNodeId, config.EVMCallTimeOut, proxyClient}
	} else {
		eth.APIBackend = &EthAPIBackend[T,P]{stack.Config().ExtRPCEnabled(), stack.Config().AllowUnprotectedTxs, eth, nil, hexNodeId, config.EVMCallTimeOut, nil}
	}

	if eth.APIBackend.allowUnprotectedTxs {
		log.Info("Unprotected transactions allowed")
	}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.Miner.GasPrice
	}
	eth.APIBackend.gpo = gasprice.NewOracle[P](eth.APIBackend, gpoParams)

	// Setup DNS discovery iterators.
	dnsclient := dnsdisc.NewClient[T](dnsdisc.Config[P]{})
	eth.ethDialCandidates, err = dnsclient.NewIterator(eth.config.EthDiscoveryURLs...)
	if err != nil {
		return nil, err
	}
	eth.snapDialCandidates, err = dnsclient.NewIterator(eth.config.SnapDiscoveryURLs...)
	if err != nil {
		return nil, err
	}
	// Start the RPC service
	eth.netRPCService = ethapi.NewPublicNetAPI(eth.p2pServer, config.NetworkId)

	// Register the backend on the node
	stack.RegisterAPIs(eth.APIs())
	if eth.config.QuorumLightClient.Enabled() && eth.config.QuorumLightClient.TokenEnabled &&
		eth.config.QuorumLightClient.TokenManagement == "external" {
		stack.RegisterAPIs(eth.QLightClientAPIs())
	}
	stack.RegisterProtocols(eth.Protocols())
	if eth.config.QuorumLightServer {
		stack.RegisterQProtocols(eth.QProtocols())
	}
	stack.RegisterLifecycle(eth)
	// Check for unclean shutdown
	if uncleanShutdowns, discards, err := rawdb.PushUncleanShutdownMarker(chainDb); err != nil {
		log.Error("Could not update unclean-shutdown-marker list", "error", err)
	} else {
		if discards > 0 {
			log.Warn("Old unclean shutdowns found", "count", discards)
		}
		for _, tstamp := range uncleanShutdowns {
			t := time.Unix(int64(tstamp), 0)
			log.Warn("Unclean shutdown detected", "booted", t,
				"age", common.PrettyAge(t))
		}
	}
	return eth, nil
}

func makeExtraData(extra []byte, isQuorum bool) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"geth",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.GetMaximumExtraDataSize(isQuorum) {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.GetMaximumExtraDataSize(isQuorum))
		extra = nil
	}
	return extra
}

func (s *Ethereum[T,P]) QLightClientAPIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "qlight",
			Version:   "1.0",
			Service:   qlight.NewPrivateQLightAPI[T,P](s.handler.peers, s.APIBackend.proxyClient),
			Public:    false,
		},
	}
}

// APIs return the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *Ethereum[T,P]) APIs() []rpc.API {
	apis := ethapi.GetAPIs[T,P](s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append all the local APIs and return
	apis = append(apis, []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.handler.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI[P](s.APIBackend, false, 5*time.Minute),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
	return apis
}

func (s *Ethereum[T,P]) ResetWithGenesisBlock(gb *types.Block[P]) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *Ethereum[T,P]) Etherbase() (eb common.Address, err error) {
	s.lock.RLock()
	etherbase := s.etherbase
	s.lock.RUnlock()

	if etherbase != (common.Address{}) {
		return etherbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			etherbase := accounts[0].Address

			s.lock.Lock()
			s.etherbase = etherbase
			s.lock.Unlock()

			log.Info("Etherbase automatically configured", "address", etherbase)
			return etherbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("etherbase must be explicitly specified")
}

// isLocalBlock checks whether the specified block is mined
// by local miner accounts.
//
// We regard two types of accounts as local miner account: etherbase
// and accounts specified via `txpool.locals` flag.
func (s *Ethereum[T,P]) isLocalBlock(block *types.Block[P]) bool {
	author, err := s.engine.Author(block.Header())
	if err != nil {
		log.Warn("Failed to retrieve block author", "number", block.NumberU64(), "hash", block.Hash(), "err", err)
		return false
	}
	// Check whether the given address is etherbase.
	s.lock.RLock()
	etherbase := s.etherbase
	s.lock.RUnlock()
	if author == etherbase {
		return true
	}
	// Check whether the given address is specified by `txpool.local`
	// CLI flag.
	for _, account := range s.config.TxPool.Locals {
		if account == author {
			return true
		}
	}
	return false
}

// shouldPreserve checks whether we should preserve the given block
// during the chain reorg depending on whether the author of block
// is a local account.
func (s *Ethereum[T,P]) shouldPreserve(block *types.Block[P]) bool {
	// The reason we need to disable the self-reorg preserving for clique
	// is it can be probable to introduce a deadlock.
	//
	// e.g. If there are 7 available signers
	//
	// r1   A
	// r2     B
	// r3       C
	// r4         D
	// r5   A      [X] F G
	// r6    [X]
	//
	// In the round5, the inturn signer E is offline, so the worst case
	// is A, F and G sign the block of round5 and reject the block of opponents
	// and in the round6, the last available signer B is offline, the whole
	// network is stuck.
	if _, ok := s.engine.(*clique.Clique[P]); ok {
		return false
	}
	return s.isLocalBlock(block)
}

// SetEtherbase sets the mining reward address.
// Quorum: method now has a return value
func (s *Ethereum[T,P]) SetEtherbase(etherbase common.Address) bool {
	//Quorum
	consensusAlgo := s.handler.getConsensusAlgorithm()
	if consensusAlgo == "istanbul" || consensusAlgo == "clique" || consensusAlgo == "raft" {
		log.Error("Cannot set etherbase with selected consensus mechanism")
		return false
	}
	//End-Quorum

	s.lock.Lock()
	s.etherbase = etherbase
	s.lock.Unlock()

	s.miner.SetEtherbase(etherbase)
	return true
}

// StartMining starts the miner with the given number of CPU threads. If mining
// is already running, this method adjust the number of threads allowed to use
// and updates the minimum price required by the transaction pool.
func (s *Ethereum[T,P]) StartMining(threads int) error {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		log.Info("Updated mining threads", "threads", threads)
		if threads == 0 {
			threads = -1 // Disable the miner from within
		}
		th.SetThreads(threads)
	}
	// If the miner was not running, initialize it
	if !s.IsMining() {
		// Propagate the initial price point to the transaction pool
		s.lock.RLock()
		price := s.gasPrice
		s.lock.RUnlock()
		s.txPool.SetGasPrice(price)

		// Configure the local mining address
		eb, err := s.Etherbase()
		if err != nil {
			log.Error("Cannot start mining without etherbase", "err", err)
			return fmt.Errorf("etherbase missing: %v", err)
		}
		if clique, ok := s.engine.(*clique.Clique[P]); ok {
			wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
			if wallet == nil || err != nil {
				log.Error("Etherbase account unavailable locally", "err", err)
				return fmt.Errorf("signer missing: %v", err)
			}
			clique.Authorize(eb, wallet.SignData)
		}
		// If mining is started, we can disable the transaction rejection mechanism
		// introduced to speed sync times.
		atomic.StoreUint32(&s.handler.acceptTxs, 1)

		go s.miner.Start(eb)
	}
	return nil
}

// StopMining terminates the miner, both at the consensus engine level as well as
// at the block creation level.
func (s *Ethereum[T,P]) StopMining() {
	// Update the thread count within the consensus engine
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := s.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	// Stop the block creating itself
	s.miner.Stop()
}

func (s *Ethereum[T,P]) IsMining() bool      { return s.miner.Mining() }
func (s *Ethereum[T,P]) Miner() *miner.Miner[T,P] { return s.miner }

func (s *Ethereum[T,P]) AccountManager() *accounts.Manager[P]  { return s.accountManager }
func (s *Ethereum[T,P]) BlockChain() *core.BlockChain[P]      { return s.blockchain }
func (s *Ethereum[T,P]) TxPool() *core.TxPool[P]               { return s.txPool }
func (s *Ethereum[T,P]) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Ethereum[T,P]) Engine() consensus.Engine[P]           { return s.engine }
func (s *Ethereum[T,P]) ChainDb() ethdb.Database            { return s.chainDb }
func (s *Ethereum[T,P]) IsListening() bool                  { return true } // Always listening
func (s *Ethereum[T,P]) Downloader() *downloader.Downloader[T,P] { return s.handler.downloader }
func (s *Ethereum[T,P]) Synced() bool                       { return atomic.LoadUint32(&s.handler.acceptTxs) == 1 }
func (s *Ethereum[T,P]) ArchiveMode() bool                  { return s.config.NoPruning }
func (s *Ethereum[T,P]) BloomIndexer() *core.ChainIndexer[P]   { return s.bloomIndexer }

// Quorum
// adds quorum specific protocols to the Protocols() function which in the associated upstream geth version returns
// only one subprotocol, "eth", and the supported versions of the "eth" protocol.
// Quorum uses the eth service to run configurable consensus protocols, e.g. istanbul. Thru release v20.10.0
// the "eth" subprotocol would be replaced with a modified subprotocol, e.g. "istanbul/99" which would contain all the "eth"
// messages + the istanbul message and be communicated over the consensus specific subprotocol ("istanbul"), and
// not over "eth".
// Now the eth service supports multiple protocols, e.g. "eth" and an optional consensus
// protocol, e.g. "istanbul/100".
// /Quorum

// Protocols returns all the currently configured
// network protocols to start.
func (s *Ethereum[T,P]) Protocols() []p2p.Protocol[T,P] {
	if s.config.QuorumLightClient.Enabled() {
		protos := qlightproto.MakeProtocolsClient[T,P]((*qlightClientHandler[T,P])(s.handler), s.networkID, s.ethDialCandidates)
		return protos
	}

	protos := eth.MakeProtocols[T,P]((*ethHandler[T,P])(s.handler), s.networkID, s.ethDialCandidates)
	if s.config.SnapshotCache > 0 {
		protos = append(protos, snap.MakeProtocols[T,P]((*snapHandler[T,P])(s.handler), s.snapDialCandidates)...)
	}

	// /Quorum
	// add additional quorum consensus protocol if set and if not set to "eth", e.g. istanbul
	if quorumConsensusProtocolName != "" && quorumConsensusProtocolName != eth.ProtocolName {
		quorumProtos := s.quorumConsensusProtocols((*ethHandler[T,P])(s.handler), s.networkID, s.ethDialCandidates)
		protos = append(protos, quorumProtos...)
	}
	// /end Quorum

	return protos
}

func (s *Ethereum[T,P]) QProtocols() []p2p.Protocol[T,P] {
	protos := qlightproto.MakeProtocolsServer[T,P]((*qlightServerHandler[T,P])(s.qlightServerHandler), s.networkID, s.ethDialCandidates)
	return protos
}

// Start implements node.Lifecycle, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *Ethereum[T,P]) Start() error {
	eth.StartENRUpdater(s.blockchain, s.p2pServer.LocalNode())

	// Start the bloom bits servicing goroutines
	s.startBloomHandlers(params.BloomBitsBlocks)

	// Figure out a max peers count based on the server limits
	maxPeers := s.p2pServer.MaxPeers
	if s.config.LightServ > 0 {
		if s.config.LightPeers >= s.p2pServer.MaxPeers {
			return fmt.Errorf("invalid peer config: light peer count (%d) >= total peer count (%d)", s.config.LightPeers, s.p2pServer.MaxPeers)
		}
		maxPeers -= s.config.LightPeers
	}
	// Start the networking layer and the light server if requested
	if s.config.QuorumLightClient.Enabled() {
		s.handler.StartQLightClient()
	} else {
		s.handler.Start(maxPeers)
		if s.qlightServerHandler != nil {
			s.qlightServerHandler.StartQLightServer(s.qlightP2pServer.MaxPeers)
		}
	}

	return nil
}

// Stop implements node.Lifecycle, terminating all internal goroutines used by the
// Ethereum protocol.
func (s *Ethereum[T,P]) Stop() error {
	// Stop all the peer-related stuff first.
	if s.config.QuorumLightClient.Enabled() {
		s.handler.StopQLightClient()
	} else {
		if s.qlightServerHandler != nil {
			s.qlightServerHandler.StopQLightServer()
		}
		s.ethDialCandidates.Close()
		s.snapDialCandidates.Close()
		s.handler.Stop()
	}

	// Then stop everything else.
	s.bloomIndexer.Close()
	close(s.closeBloomHandler)
	s.txPool.Stop()
	s.miner.Stop()
	s.blockchain.Stop()
	s.engine.Close()
	rawdb.PopUncleanShutdownMarker(s.chainDb)
	s.chainDb.Close()
	s.eventMux.Stop()

	return nil
}

func (s *Ethereum[T,P]) CalcGasLimit(block *types.Block[P]) uint64 {
	minGasLimit := params.DefaultMinGasLimit
	if s != nil && s.config != nil && s.config.Genesis != nil {
		minGasLimit = s.config.Genesis.Config.GetMinerMinGasLimit(block.Number(), params.DefaultMinGasLimit)
	}
	return core.CalcGasLimit(block, minGasLimit, s.config.Miner.GasFloor, s.config.Miner.GasCeil)
}

// (Quorum)
// ConsensusServicePendingLogsFeed returns an event.Feed.  When the consensus protocol does not use eth.worker (e.g. raft), the event.Feed should be used to send logs from transactions included in the pending block
func (s *Ethereum[T,P]) ConsensusServicePendingLogsFeed() *event.Feed {
	return s.consensusServicePendingLogsFeed
}

// (Quorum)
// SubscribePendingLogs starts delivering logs from transactions included in the consensus engine's pending block to the given channel.
func (s *Ethereum[T,P]) SubscribePendingLogs(ch chan<- []*types.Log) event.Subscription {
	if s.config.RaftMode {
		return s.consensusServicePendingLogsFeed.Subscribe(ch)
	}
	return s.miner.SubscribePendingLogs(ch)
}

// NotifyRegisteredPluginService will ask to refresh the plugin service
// (Quorum)
func (s *Ethereum[T,P]) NotifyRegisteredPluginService(pluginManager *plugin.PluginManager[T,P]) error {
	if s.qlightTokenHolder == nil {
		return nil
	}
	switch s.config.QuorumLightClient.TokenManagement {
	case "client-security-plugin":
		return s.qlightTokenHolder.RefreshPlugin(pluginManager)
	}
	return nil
}
