// Copyright 2015 The go-ethereum Authors
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

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/pavelkrolevets/MIR-pro/cmd/utils"
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/mps"
	"github.com/pavelkrolevets/MIR-pro/core/rawdb"
	"github.com/pavelkrolevets/MIR-pro/core/state"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/log"
	"github.com/pavelkrolevets/MIR-pro/metrics"
	"github.com/pavelkrolevets/MIR-pro/private"
	"github.com/pavelkrolevets/MIR-pro/private/engine/notinuse"
	"github.com/pavelkrolevets/MIR-pro/rlp"
	"gopkg.in/urfave/cli.v1"
)

var (
	initCommand = cli.Command{
		Action:    utils.MigrateFlags(initGenesisType),
		Name:      "init",
		Usage:     "Bootstrap and initialize a new genesis block",
		ArgsUsage: "<genesisPath>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CryptoSwitchFlag,
			utils.CryptoGostCurveFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}
	dumpGenesisCommand = cli.Command{
		Action:    utils.MigrateFlags(dumpGenesis[nist.PublicKey]),
		Name:      "dumpgenesis",
		Usage:     "Dumps genesis block JSON configuration to stdout",
		ArgsUsage: "",
		Flags: []cli.Flag{
			utils.MainnetMirFlag,
			utils.SoyuzFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The dumpgenesis command dumps the genesis block configuration in JSON format to stdout.`,
	}
	importCommand = cli.Command{
		Action:    utils.MigrateFlags(importChain[nist.PrivateKey, nist.PublicKey]),
		Name:      "import",
		Usage:     "Import a blockchain file",
		ArgsUsage: "<filename> (<filename 2> ... <filename N>) ",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
			utils.GCModeFlag,
			utils.SnapshotFlag,
			utils.CacheDatabaseFlag,
			utils.CacheGCFlag,
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
			utils.TxLookupLimitFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The import command imports blocks from an RLP-encoded form. The form can be one file
with several RLP-encoded blocks, or several files can be used.

If only one file is used, import error will result in failure. If several files are used,
processing will proceed even if an individual RLP-file import failure occurs.`,
	}
	mpsdbUpgradeCommand = cli.Command{
		Action:    utils.MigrateFlags(mpsdbUpgrade[nist.PrivateKey, nist.PublicKey]),
		Name:      "mpsdbupgrade",
		Usage:     "Upgrade a standalone DB to an MPS DB",
		ArgsUsage: "",
		Flags: []cli.Flag{
			utils.DataDirFlag,
		},
		Description: `
Checks if the chain config isMPS parameter value.
If false, it upgrades the DB to be MPS enabled (builds the trie of private states) and if successful sets isMPS to true.
If true, exits displaying an error message that the DB is already MPS.`,
		Category: "BLOCKCHAIN COMMANDS",
	}
	exportCommand = cli.Command{
		Action:    utils.MigrateFlags(exportChain[nist.PrivateKey, nist.PublicKey]),
		Name:      "export",
		Usage:     "Export blockchain into file",
		ArgsUsage: "<filename> [<blockNumFirst> <blockNumLast>]",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
Requires a first argument of the file to write to.
Optional second and third arguments control the first and
last block to write. In this mode, the file will be appended
if already existing. If the file ends with .gz, the output will
be gzipped.`,
	}
	importPreimagesCommand = cli.Command{
		Action:    utils.MigrateFlags(importPreimages[nist.PrivateKey, nist.PublicKey]),
		Name:      "import-preimages",
		Usage:     "Import the preimage database from an RLP stream",
		ArgsUsage: "<datafile>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
	The import-preimages command imports hash preimages from an RLP encoded stream.`,
	}
	exportPreimagesCommand = cli.Command{
		Action:    utils.MigrateFlags(exportPreimages[nist.PrivateKey, nist.PublicKey]),
		Name:      "export-preimages",
		Usage:     "Export the preimage database into an RLP stream",
		ArgsUsage: "<dumpfile>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The export-preimages command export hash preimages to an RLP encoded stream`,
	}
	dumpCommand = cli.Command{
		Action:    utils.MigrateFlags(dump[nist.PrivateKey, nist.PublicKey]),
		Name:      "dump",
		Usage:     "Dump a specific block from storage",
		ArgsUsage: "[<blockHash> | <blockNum>]...",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
			utils.IterativeOutputFlag,
			utils.ExcludeCodeFlag,
			utils.ExcludeStorageFlag,
			utils.IncludeIncompletesFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The arguments are interpreted as block numbers or hashes.
Use "ethereum dump 0" to dump the genesis block.`,
	}
)

// In the regular Genesis / ChainConfig struct, due to the way go deserializes
// json, IsQuorum defaults to false (when not specified). Here we specify it as
// a pointer so we can make the distinction and default unspecified to true.
func getIsQuorum(file io.Reader) bool {
	altGenesis := new(struct {
		Config *struct {
			IsQuorum *bool `json:"isQuorum"`
		} `json:"config"`
	})

	if err := json.NewDecoder(file).Decode(altGenesis); err != nil {
		utils.Fatalf("invalid genesis file: %v", err)
	}

	// unspecified defaults to true
	// Mir defaults to false
	return altGenesis.Config.IsQuorum != nil
}

func initGenesisType(ctx *cli.Context) error {
	// Mir - set crypto before the start of services 
	if ctx.GlobalString(utils.CryptoSwitchFlag.Name) != "" {
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
			return initGenesis[gost3410.PrivateKey,gost3410.PublicKey](ctx)
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "gost_csp" {
			crypto.CryptoAlg = crypto.GOST_CSP
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "nist" {
			fmt.Println(`
			╔═╗┬─┐┬ ┬┌─┐┌┬┐┌─┐  ╔╗╔╦╔═╗╔╦╗
			║  ├┬┘└┬┘├─┘ │ │ │  ║║║║╚═╗ ║ 
			╚═╝┴└─ ┴ ┴   ┴ └─┘  ╝╚╝╩╚═╝ ╩ `)
			return initGenesis[nist.PrivateKey,nist.PublicKey](ctx)
		} else if ctx.GlobalString(utils.CryptoSwitchFlag.Name) == "pqc" {
			crypto.CryptoAlg = crypto.PQC
		} else {
			fmt.Errorf("wrong crypto flag")
		}
	}
	return nil
}
// initGenesis will initialise the given JSON format genesis file and writes it as
// the zero'd block (i.e. genesis) or will fail hard if it can't succeed.
func initGenesis[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	// Make sure we have a valid genesis JSON
	genesisPath := ctx.Args().First()
	if len(genesisPath) == 0 {
		utils.Fatalf("Must supply path to genesis JSON file")
	}
	file, err := os.Open(genesisPath)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer file.Close()

	genesis := new(core.Genesis[P])
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		utils.Fatalf("invalid genesis file: %v", err)
	}

	// Quorum
	file.Seek(0, 0)
	genesis.Config.IsQuorum = getIsQuorum(file)
	fmt.Printf("Is Quorum %t \n ", genesis.Config.IsQuorum)
	// check the data given as a part of newMaxConfigData to ensure that
	// its in expected order
	err = genesis.Config.CheckMaxCodeConfigData()
	if err != nil {
		utils.Fatalf("maxCodeSize data invalid: %v", err)
	}
	if genesis.Config.IsQuorum {
		err = genesis.Config.CheckTransitionsData()
		if err != nil {
			utils.Fatalf("transitions data invalid: %v", err)
		}
		if genesis.Config.QBFT != nil && genesis.Config.QBFT.ValidatorContractAddress != (common.Address{}) {
			qbftExtra := new(types.QBFTExtra)
			err := rlp.DecodeBytes(genesis.ExtraData[:], qbftExtra)
			if err != nil || len(qbftExtra.Validators) > 0 {
				utils.Fatalf("invalid genesis file: cant combine extraData validators and config.qbft.validatorcontractaddress at the same time")
			}
		}
		if genesis.Config.IBFT != nil && genesis.Config.IBFT.ValidatorContractAddress != (common.Address{}) {
			istanbulExtra := new(types.IstanbulExtra)
			err := rlp.DecodeBytes(genesis.ExtraData[:], istanbulExtra)
			if err != nil || len(istanbulExtra.Validators) > 0 {
				utils.Fatalf("invalid genesis file: cant combine extraData validators and config.ibft.validatorcontractaddress at the same time")
			}
		}
	}
	// End Quorum

	// Open and initialise both full and light databases
	stack, _ := makeConfigNode[T,P](ctx)
	defer stack.Close()

	for _, name := range []string{"chaindata", "lightchaindata"} {
		chaindb, err := stack.OpenDatabase(name, 0, 0, "", false)
		if err != nil {
			utils.Fatalf("Failed to open database: %v", err)
		}
		_, hash, err := core.SetupGenesisBlock(chaindb, genesis)
		if err != nil {
			utils.Fatalf("Failed to write genesis block: %v", err)
		}
		chaindb.Close()
		log.Info("Successfully wrote genesis state", "database", name, "hash", hash)
	}
	return nil
}

func dumpGenesis[P crypto.PublicKey](ctx *cli.Context) error {
	// TODO(rjl493456442) support loading from the custom datadir
	genesis := utils.MakeGenesis[P](ctx)
	if genesis == nil {
		genesis = core.DefaultGenesisBlock[P]()
	}
	if err := json.NewEncoder(os.Stdout).Encode(genesis); err != nil {
		utils.Fatalf("could not encode genesis")
	}
	return nil
}

func importChain[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}
	// Start metrics export if enabled
	utils.SetupMetrics(ctx)
	// Start system runtime metrics collection
	go metrics.CollectProcessMetrics(3 * time.Second)

	stack, _ := makeConfigNode[T,P](ctx)
	defer stack.Close()

	chain, db := utils.MakeChain(ctx, stack, true)
	defer db.Close()

	// Start periodically gathering memory profiles
	var peakMemAlloc, peakMemSys uint64
	go func() {
		stats := new(runtime.MemStats)
		for {
			runtime.ReadMemStats(stats)
			if atomic.LoadUint64(&peakMemAlloc) < stats.Alloc {
				atomic.StoreUint64(&peakMemAlloc, stats.Alloc)
			}
			if atomic.LoadUint64(&peakMemSys) < stats.Sys {
				atomic.StoreUint64(&peakMemSys, stats.Sys)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	// Import the chain
	start := time.Now()

	var importErr error

	if len(ctx.Args()) == 1 {
		if err := utils.ImportChain[T,P](chain, ctx.Args().First()); err != nil {
			importErr = err
			log.Error("Import error", "err", err)
		}
	} else {
		for _, arg := range ctx.Args() {
			if err := utils.ImportChain[T,P](chain, arg); err != nil {
				importErr = err
				log.Error("Import error", "file", arg, "err", err)
			}
		}
	}
	chain.Stop()
	fmt.Printf("Import done in %v.\n\n", time.Since(start))

	// Output pre-compaction stats mostly to see the import trashing
	showLeveldbStats(db)

	// Print the memory statistics used by the importing
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)

	fmt.Printf("Object memory: %.3f MB current, %.3f MB peak\n", float64(mem.Alloc)/1024/1024, float64(atomic.LoadUint64(&peakMemAlloc))/1024/1024)
	fmt.Printf("System memory: %.3f MB current, %.3f MB peak\n", float64(mem.Sys)/1024/1024, float64(atomic.LoadUint64(&peakMemSys))/1024/1024)
	fmt.Printf("Allocations:   %.3f million\n", float64(mem.Mallocs)/1000000)
	fmt.Printf("GC pause:      %v\n\n", time.Duration(mem.PauseTotalNs))

	if ctx.GlobalBool(utils.NoCompactionFlag.Name) {
		return nil
	}

	// Compact the entire database to more accurately measure disk io and print the stats
	start = time.Now()
	fmt.Println("Compacting entire database...")
	if err := db.Compact(nil, nil); err != nil {
		utils.Fatalf("Compaction failed: %v", err)
	}
	fmt.Printf("Compaction done in %v.\n\n", time.Since(start))

	showLeveldbStats(db)
	return importErr
}

func mpsdbUpgrade[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	stack, _ := makeConfigNode[T,P](ctx)
	defer stack.Close()

	// initialise the tx manager with the dummy tx mgr
	private.Ptm = &notinuse.DBUpgradePrivateTransactionManager{}

	chain, db := utils.MakeChain(ctx, stack, true)

	if chain.Config().IsMPS {
		utils.Fatalf("The database is already upgraded to support multiple private states.")
	}

	currentBlockNumber := chain.CurrentBlock().Number().Int64()
	fmt.Printf("Current block number %v\n", currentBlockNumber)

	return mps.UpgradeDB[P](db, chain)
}

func exportChain[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	stack, _ := makeConfigNode[T,P](ctx)
	defer stack.Close()

	chain, _ := utils.MakeChain(ctx, stack, true)
	start := time.Now()

	var err error
	fp := ctx.Args().First()
	if len(ctx.Args()) < 3 {
		err = utils.ExportChain(chain, fp)
	} else {
		// This can be improved to allow for numbers larger than 9223372036854775807
		first, ferr := strconv.ParseInt(ctx.Args().Get(1), 10, 64)
		last, lerr := strconv.ParseInt(ctx.Args().Get(2), 10, 64)
		if ferr != nil || lerr != nil {
			utils.Fatalf("Export error in parsing parameters: block number not an integer\n")
		}
		if first < 0 || last < 0 {
			utils.Fatalf("Export error: block number must be greater than 0\n")
		}
		if head := chain.CurrentFastBlock(); uint64(last) > head.NumberU64() {
			utils.Fatalf("Export error: block number %d larger than head block %d\n", uint64(last), head.NumberU64())
		}
		err = utils.ExportAppendChain(chain, fp, uint64(first), uint64(last))
	}

	if err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}
	fmt.Printf("Export done in %v\n", time.Since(start))
	return nil
}

// importPreimages imports preimage data from the specified file.
func importPreimages[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	stack, _ := makeConfigNode[T,P](ctx)
	defer stack.Close()

	db := utils.MakeChainDatabase(ctx, stack, false)
	start := time.Now()

	if err := utils.ImportPreimages[P](db, ctx.Args().First()); err != nil {
		utils.Fatalf("Import error: %v\n", err)
	}
	fmt.Printf("Import done in %v\n", time.Since(start))
	return nil
}

// exportPreimages dumps the preimage data to specified json file in streaming way.
func exportPreimages[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	stack, _ := makeConfigNode[T,P](ctx)
	defer stack.Close()

	db := utils.MakeChainDatabase(ctx, stack, true)
	start := time.Now()

	if err := utils.ExportPreimages(db, ctx.Args().First()); err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}
	fmt.Printf("Export done in %v\n", time.Since(start))
	return nil
}

func dump[T crypto.PrivateKey, P crypto.PublicKey](ctx *cli.Context) error {
	stack, _ := makeConfigNode[T,P](ctx)
	defer stack.Close()

	db := utils.MakeChainDatabase(ctx, stack, false)
	for _, arg := range ctx.Args() {
		var header *types.Header[P]
		if hashish(arg) {
			hash := common.HexToHash(arg)
			number := rawdb.ReadHeaderNumber(db, hash)
			if number != nil {
				header = rawdb.ReadHeader[P](db, hash, *number)
			}
		} else {
			number, _ := strconv.Atoi(arg)
			hash := rawdb.ReadCanonicalHash(db, uint64(number))
			if hash != (common.Hash{}) {
				header = rawdb.ReadHeader[P](db, hash, uint64(number))
			}
		}
		if header == nil {
			fmt.Println("{}")
			utils.Fatalf("block not found")
		} else {
			state, err := state.New[P](header.Root, state.NewDatabase[P](db), nil)
			if err != nil {
				utils.Fatalf("could not create new state: %v", err)
			}
			excludeCode := ctx.Bool(utils.ExcludeCodeFlag.Name)
			excludeStorage := ctx.Bool(utils.ExcludeStorageFlag.Name)
			includeMissing := ctx.Bool(utils.IncludeIncompletesFlag.Name)
			if ctx.Bool(utils.IterativeOutputFlag.Name) {
				state.IterativeDump(excludeCode, excludeStorage, !includeMissing, json.NewEncoder(os.Stdout))
			} else {
				if includeMissing {
					fmt.Printf("If you want to include accounts with missing preimages, you need iterative output, since" +
						" otherwise the accounts will overwrite each other in the resulting mapping.")
				}
				fmt.Printf("%v %s\n", includeMissing, state.Dump(excludeCode, excludeStorage, false))
			}
		}
	}
	return nil
}

// hashish returns true for strings that look like hashes.
func hashish(x string) bool {
	_, err := strconv.Atoi(x)
	return err != nil
}
