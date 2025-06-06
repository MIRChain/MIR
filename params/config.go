// Copyright 2016 The go-ethereum Authors
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

package params

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/common/math"
	"github.com/pavelkrolevets/MIR-pro/log"
	"golang.org/x/crypto/sha3"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetMirGenesisHash = common.HexToHash("0x00ff17a7de0e3833ca49125f2f266e09c5ba1927e107ab34374e19e97d356ec4")
	MainnetEthGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	SoyuzGenesisHash      = common.HexToHash("0x8e09a1617e6d5f98f5f784e0077d4e054842c2f9ae99d4570e1e4470b1e6d190")
	RopstenGenesisHash    = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	RinkebyGenesisHash    = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
)

// TrustedCheckpoints associates each known checkpoint with the genesis hash of
// the chain it belongs to.
var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetEthGenesisHash: MainnetEthTrustedCheckpoint,
	RinkebyGenesisHash:    RinkebyTrustedCheckpoint,
}

// CheckpointOracles associates each known checkpoint oracles with the genesis hash of
// the chain it belongs to.
var CheckpointOracles = map[common.Hash]*CheckpointOracleConfig{
	MainnetEthGenesisHash: MainnetEthCheckpointOracle,
	RinkebyGenesisHash:    RinkebyCheckpointOracle,
}

func newUint64(val uint64) *uint64 { return &val }

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetMirChainConfig = &ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(10),
		ConstantinopleBlock: big.NewInt(10),
		PetersburgBlock:     big.NewInt(10),
		IstanbulBlock:       big.NewInt(10),
		MuirGlacierBlock:    big.NewInt(10),
		BerlinBlock:         big.NewInt(10),
		Ethash:              new(EthashConfig),
	}

	// SoyuzChainConfig contains the chain parameters to run a node on the Soyuz test network.
	SoyuzChainConfig = &ChainConfig{
		ChainID:             big.NewInt(3),
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(10),
		PetersburgBlock:     big.NewInt(10),
		IstanbulBlock:       big.NewInt(10),
		MuirGlacierBlock:    big.NewInt(10),
		BerlinBlock:         big.NewInt(10),
		Ethash:              new(EthashConfig),
	}

	MainnetTerminalTotalDifficulty, _ = new(big.Int).SetString("58_750_000_000_000_000_000_000", 0)

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetEthChainConfig = &ChainConfig{
		ChainID:                       big.NewInt(1),
		HomesteadBlock:                big.NewInt(1_150_000),
		DAOForkBlock:                  big.NewInt(1_920_000),
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(2_463_000),
		EIP150Hash:                    common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:                   big.NewInt(2_675_000),
		EIP158Block:                   big.NewInt(2_675_000),
		ByzantiumBlock:                big.NewInt(4_370_000),
		ConstantinopleBlock:           big.NewInt(7_280_000),
		PetersburgBlock:               big.NewInt(7_280_000),
		IstanbulBlock:                 big.NewInt(9_069_000),
		MuirGlacierBlock:              big.NewInt(9_200_000),
		BerlinBlock:                   big.NewInt(12_244_000),
		LondonBlock:                   big.NewInt(12_965_000),
		ArrowGlacierBlock:             big.NewInt(13_773_000),
		GrayGlacierBlock:              big.NewInt(15_050_000),
		TerminalTotalDifficulty:       MainnetTerminalTotalDifficulty, // 58_750_000_000_000_000_000_000
		TerminalTotalDifficultyPassed: true,
		ShanghaiTime:                  newUint64(1681338455),
		Ethash:                        new(EthashConfig),
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetEthTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 371,
		SectionHead:  common.HexToHash("0x50fd3cec5376ede90ef9129772022690cd1467f22c18abb7faa11e793c51e9c9"),
		CHTRoot:      common.HexToHash("0xb57b4b22a77b5930847b1ca9f62daa11eae6578948cb7b18997f2c0fe5757025"),
		BloomRoot:    common.HexToHash("0xa338f8a868a194fa90327d0f5877f656a9f3640c618d2a01a01f2e76ef9ef954"),
	}

	// MainnetCheckpointOracle contains a set of configs for the main network oracle.
	MainnetEthCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0x9a9070028361F7AAbeB3f2F2Dc07F82C4a98A02a"),
		Signers: []common.Address{
			common.HexToAddress("0x1b2C260efc720BE89101890E4Db589b44E950527"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// RinkebyChainConfig contains the chain parameters to run a node on the Rinkeby test network.
	RinkebyChainConfig = &ChainConfig{
		ChainID:             big.NewInt(4),
		HomesteadBlock:      big.NewInt(1),
		DAOForkBlock:        nil,
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(2),
		EIP150Hash:          common.HexToHash("0x9b095b36c15eaf13044373aef8ee0bd3a382a5abb92e402afa44b8249c3a90e9"),
		EIP155Block:         big.NewInt(3),
		EIP158Block:         big.NewInt(3),
		ByzantiumBlock:      big.NewInt(1_035_301),
		ConstantinopleBlock: big.NewInt(3_660_663),
		PetersburgBlock:     big.NewInt(4_321_234),
		IstanbulBlock:       big.NewInt(5_435_345),
		MuirGlacierBlock:    nil,
		BerlinBlock:         big.NewInt(8_290_928),
		Clique: &CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}

	// RinkebyTrustedCheckpoint contains the light client trusted checkpoint for the Rinkeby test network.
	RinkebyTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 254,
		SectionHead:  common.HexToHash("0x0cba01dd71baa22ac8fa0b105bc908e94f9ecfbc79b4eb97427fe07b5851dd10"),
		CHTRoot:      common.HexToHash("0x5673d8fc49c9c7d8729068640e4b392d46952a5a38798973bac1cf1d0d27ad7d"),
		BloomRoot:    common.HexToHash("0x70e01232b66df9a7778ae3291c9217afb9a2d9f799f32d7b912bd37e7bce83a8"),
	}

	// RinkebyCheckpointOracle contains a set of configs for the Rinkeby test network oracle.
	RinkebyCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0xebe8eFA441B9302A0d7eaECc277c09d20D684540"),
		Signers: []common.Address{
			common.HexToAddress("0xd9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f3"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
		},
		Threshold: 2,
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	// AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, new(EthashConfig), nil, nil, nil, nil, nil, false, 32, 35, big.NewInt(0), big.NewInt(0), nil, nil, false, nil, nil}

	AllEthashProtocolChanges = &ChainConfig{
		ChainID:                       big.NewInt(1337),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP150Hash:                    common.Hash{},
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		TerminalTotalDifficulty:       nil,
		TerminalTotalDifficultyPassed: false,
		Ethash:                        new(EthashConfig),
		Clique:                        nil,
	}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	// AllCliqueProtocolChanges = &ChainConfig{big.NewInt(10), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, &CliqueConfig{Period: 0, Epoch: 30000}, nil, nil, nil, nil, false, 32, 32, big.NewInt(0), big.NewInt(0), nil, nil, false, nil, nil}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	AllCliqueProtocolChanges = &ChainConfig{
		ChainID:                       big.NewInt(1337),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP150Hash:                    common.Hash{},
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             nil,
		GrayGlacierBlock:              nil,
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		TerminalTotalDifficulty:       nil,
		TerminalTotalDifficultyPassed: false,
		Ethash:                        nil,
		Clique:                        &CliqueConfig{Period: 0, Epoch: 30000},
		Istanbul:                      nil,
		IBFT:                          nil,
		QBFT:                          nil,
		Transitions:                   nil,
		IsQuorum:                      false,
		TransactionSizeLimit:          32,
		MaxCodeSize:                   32,
		QIP714Block:                   big.NewInt(0),
		MaxCodeSizeChangeBlock:        big.NewInt(0),
		MaxCodeSizeConfig:             nil,
		PrivacyEnhancementsBlock:      nil,
		IsMPS:                         false,
		PrivacyPrecompileBlock:        nil,
		EnableGasPriceBlock:           nil,
	}

	// Quorum chainID should 10
	TestChainConfig = &ChainConfig{
		ChainID:                  big.NewInt(10),
		HomesteadBlock:           big.NewInt(0),
		DAOForkBlock:             nil,
		DAOForkSupport:           false,
		EIP150Block:              big.NewInt(0),
		EIP150Hash:               common.Hash{},
		EIP155Block:              big.NewInt(0),
		EIP158Block:              big.NewInt(0),
		ByzantiumBlock:           big.NewInt(0),
		ConstantinopleBlock:      big.NewInt(0),
		PetersburgBlock:          big.NewInt(0),
		IstanbulBlock:            big.NewInt(0),
		MuirGlacierBlock:         big.NewInt(0),
		BerlinBlock:              big.NewInt(0),
		YoloV3Block:              nil,
		EWASMBlock:               nil,
		CatalystBlock:            nil,
		Ethash:                   new(EthashConfig),
		Clique:                   nil,
		Istanbul:                 nil,
		IBFT:                     nil,
		QBFT:                     nil,
		Transitions:              nil,
		IsQuorum:                 false,
		TransactionSizeLimit:     32,
		MaxCodeSize:              32,
		QIP714Block:              big.NewInt(0),
		MaxCodeSizeChangeBlock:   big.NewInt(0),
		MaxCodeSizeConfig:        nil,
		PrivacyEnhancementsBlock: nil,
		IsMPS:                    false,
		PrivacyPrecompileBlock:   nil,
		EnableGasPriceBlock:      nil,
	}
	TestRules = TestChainConfig.Rules(new(big.Int))

	// QuorumTestChainConfig    = &ChainConfig{big.NewInt(10), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil, new(EthashConfig), nil, nil, nil, nil, nil, true, 64, 32, big.NewInt(0), big.NewInt(0), nil, big.NewInt(0), false, nil, nil}
	QuorumTestChainConfig = &ChainConfig{
		ChainID:                  big.NewInt(10),
		HomesteadBlock:           big.NewInt(0),
		DAOForkBlock:             nil,
		DAOForkSupport:           false,
		EIP150Block:              big.NewInt(0),
		EIP150Hash:               common.Hash{},
		EIP155Block:              big.NewInt(0),
		EIP158Block:              big.NewInt(0),
		ByzantiumBlock:           big.NewInt(0),
		ConstantinopleBlock:      big.NewInt(0),
		PetersburgBlock:          big.NewInt(0),
		IstanbulBlock:            big.NewInt(0),
		MuirGlacierBlock:         big.NewInt(0),
		BerlinBlock:              big.NewInt(0),
		YoloV3Block:              nil,
		EWASMBlock:               nil,
		CatalystBlock:            nil,
		Ethash:                   new(EthashConfig),
		Clique:                   nil,
		Istanbul:                 nil,
		IBFT:                     nil,
		QBFT:                     nil,
		Transitions:              nil,
		IsQuorum:                 false,
		TransactionSizeLimit:     64,
		MaxCodeSize:              32,
		QIP714Block:              big.NewInt(0),
		MaxCodeSizeChangeBlock:   big.NewInt(0),
		MaxCodeSizeConfig:        nil,
		PrivacyEnhancementsBlock: nil,
		IsMPS:                    false,
		PrivacyPrecompileBlock:   nil,
		EnableGasPriceBlock:      nil,
	}

	QuorumMPSTestChainConfig = &ChainConfig{
		ChainID:                  big.NewInt(10),
		HomesteadBlock:           big.NewInt(0),
		DAOForkBlock:             nil,
		DAOForkSupport:           false,
		EIP150Block:              big.NewInt(0),
		EIP150Hash:               common.Hash{},
		EIP155Block:              big.NewInt(0),
		EIP158Block:              big.NewInt(0),
		ByzantiumBlock:           big.NewInt(0),
		ConstantinopleBlock:      big.NewInt(0),
		PetersburgBlock:          big.NewInt(0),
		IstanbulBlock:            big.NewInt(0),
		MuirGlacierBlock:         big.NewInt(0),
		BerlinBlock:              big.NewInt(0),
		YoloV3Block:              nil,
		EWASMBlock:               nil,
		CatalystBlock:            nil,
		Ethash:                   new(EthashConfig),
		Clique:                   nil,
		Istanbul:                 nil,
		IBFT:                     nil,
		QBFT:                     nil,
		Transitions:              nil,
		IsQuorum:                 false,
		TransactionSizeLimit:     64,
		MaxCodeSize:              32,
		QIP714Block:              big.NewInt(0),
		MaxCodeSizeChangeBlock:   big.NewInt(0),
		MaxCodeSizeConfig:        nil,
		PrivacyEnhancementsBlock: nil,
		IsMPS:                    true,
		PrivacyPrecompileBlock:   nil,
		EnableGasPriceBlock:      nil,
	}
)

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  common.Hash `json:"sectionHead"`
	CHTRoot      common.Hash `json:"chtRoot"`
	BloomRoot    common.Hash `json:"bloomRoot"`
}

// HashEqual returns an indicator comparing the itself hash with given one.
func (c *TrustedCheckpoint) HashEqual(hash common.Hash) bool {
	if c.Empty() {
		return hash == common.Hash{}
	}
	return c.Hash() == hash
}

// Hash returns the hash of checkpoint's four key fields(index, sectionHead, chtRoot and bloomTrieRoot).
func (c *TrustedCheckpoint) Hash() common.Hash {
	var sectionIndex [8]byte
	binary.BigEndian.PutUint64(sectionIndex[:], c.SectionIndex)

	w := sha3.NewLegacyKeccak256()
	w.Write(sectionIndex[:])
	w.Write(c.SectionHead[:])
	w.Write(c.CHTRoot[:])
	w.Write(c.BloomRoot[:])

	var h common.Hash
	w.Sum(h[:0])
	return h
}

// Empty returns an indicator whether the checkpoint is regarded as empty.
func (c *TrustedCheckpoint) Empty() bool {
	return c.SectionHead == (common.Hash{}) || c.CHTRoot == (common.Hash{}) || c.BloomRoot == (common.Hash{})
}

// CheckpointOracleConfig represents a set of checkpoint contract(which acts as an oracle)
// config which used for light client checkpoint syncing.
type CheckpointOracleConfig struct {
	Address   common.Address   `json:"address"`
	Signers   []common.Address `json:"signers"`
	Threshold uint64           `json:"threshold"`
}

type MaxCodeConfigStruct struct {
	Block *big.Int `json:"block,omitempty"`
	Size  uint64   `json:"size,omitempty"`
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)
	BerlinBlock         *big.Int `json:"berlinBlock,omitempty"`         // Berlin switch block (nil = no fork, 0 = already on berlin)
	LondonBlock         *big.Int `json:"londonBlock,omitempty"`         // London switch block (nil = no fork, 0 = already on london)
	ArrowGlacierBlock   *big.Int `json:"arrowGlacierBlock,omitempty"`   // Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated)
	GrayGlacierBlock    *big.Int `json:"grayGlacierBlock,omitempty"`    // Eip-5133 (bomb delay) switch block (nil = no fork, 0 = already activated)
	MergeNetsplitBlock  *big.Int `json:"mergeNetsplitBlock,omitempty"`  // Virtual fork after The Merge to use as a network splitter

	// Fork scheduling was switched from blocks to timestamps here

	ShanghaiTime *uint64 `json:"shanghaiTime,omitempty"` // Shanghai switch time (nil = no fork, 0 = already on shanghai)
	CancunTime   *uint64 `json:"cancunTime,omitempty"`   // Cancun switch time (nil = no fork, 0 = already on cancun)
	PragueTime   *uint64 `json:"pragueTime,omitempty"`   // Prague switch time (nil = no fork, 0 = already on prague)

	YoloV3Block   *big.Int `json:"yoloV3Block,omitempty"`   // YOLO v3: Gas repricings TODO @holiman add EIP references
	EWASMBlock    *big.Int `json:"ewasmBlock,omitempty"`    // EWASM switch block (nil = no fork, 0 = already activated)
	CatalystBlock *big.Int `json:"catalystBlock,omitempty"` // Catalyst switch block (nil = no fork, 0 = already on catalyst)

	// Various consensus engines
	Ethash   *EthashConfig   `json:"ethash,omitempty"`
	Clique   *CliqueConfig   `json:"clique,omitempty"`
	Istanbul *IstanbulConfig `json:"istanbul,omitempty"` // Quorum
	IBFT     *IBFTConfig     `json:"ibft,omitempty"`     // Quorum
	QBFT     *QBFTConfig     `json:"qbft,omitempty"`     // Quorum

	// TerminalTotalDifficulty is the amount of total difficulty reached by
	// the network that triggers the consensus upgrade.
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

	// TerminalTotalDifficultyPassed is a flag specifying that the network already
	// passed the terminal total difficulty. Its purpose is to disable legacy sync
	// even without having seen the TTD locally (safer long term).
	TerminalTotalDifficultyPassed bool `json:"terminalTotalDifficultyPassed,omitempty"`

	// Start of Quorum specific configs

	Transitions          []Transition `json:"transitions,omitempty"` // Quorum - transition config based on the block number
	IsQuorum             bool         `json:"isQuorum"`              // Quorum flag
	TransactionSizeLimit uint64       `json:"txnSizeLimit"`          // Quorum - transaction size limit
	MaxCodeSize          uint64       `json:"maxCodeSize"`           // Quorum -  maximum CodeSize of contract

	// QIP714Block implements the permissions related changes
	QIP714Block            *big.Int `json:"qip714Block,omitempty"`
	MaxCodeSizeChangeBlock *big.Int `json:"maxCodeSizeChangeBlock,omitempty"`
	// to track multiple changes to maxCodeSize
	MaxCodeSizeConfig        []MaxCodeConfigStruct `json:"maxCodeSizeConfig,omitempty"`
	PrivacyEnhancementsBlock *big.Int              `json:"privacyEnhancementsBlock,omitempty"`
	IsMPS                    bool                  `json:"isMPS"`                            // multiple private states flag
	PrivacyPrecompileBlock   *big.Int              `json:"privacyPrecompileBlock,omitempty"` // Switch block to enable privacy precompiled contract to process privacy marker transactions
	EnableGasPriceBlock      *big.Int              `json:"enableGasPriceBlock,omitempty"`    // Switch block to enable usage of gas price

	// End of Quorum specific configs
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c *EthashConfig) String() string {
	return "ethash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period                 uint64 `json:"period"`                 // Number of seconds between blocks to enforce
	Epoch                  uint64 `json:"epoch"`                  // Epoch length to reset votes and checkpoint
	AllowedFutureBlockTime uint64 `json:"allowedFutureBlockTime"` // Max time (in seconds) from current time allowed for blocks, before they're considered future blocks
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// IstanbulConfig is the consensus engine configs for Istanbul based sealing.
type IstanbulConfig struct {
	Epoch          uint64   `json:"epoch"`                    // Epoch length to reset votes and checkpoint
	ProposerPolicy uint64   `json:"policy"`                   // The policy for proposer selection
	Ceil2Nby3Block *big.Int `json:"ceil2Nby3Block,omitempty"` // Number of confirmations required to move from one state to next [2F + 1 to Ceil(2N/3)]
	TestQBFTBlock  *big.Int `json:"testQBFTBlock,omitempty"`  // Fork block at which block confirmations are done using qbft consensus instead of ibft
}

// String implements the stringer interface, returning the consensus engine details.
func (c *IstanbulConfig) String() string {
	return "istanbul"
}

type BFTConfig struct {
	EpochLength              uint64         `json:"epochlength"`                       // Number of blocks that should pass before pending validator votes are reset
	BlockPeriodSeconds       uint64         `json:"blockperiodseconds"`                // Minimum time between two consecutive IBFT or QBFT blocks’ timestamps in seconds
	EmptyBlockPeriodSeconds  *uint64        `json:"emptyblockperiodseconds,omitempty"` // Minimum time between two consecutive IBFT or QBFT a block and empty block’ timestamps in seconds
	RequestTimeoutSeconds    uint64         `json:"requesttimeoutseconds"`             // Minimum request timeout for each IBFT or QBFT round in milliseconds
	ProposerPolicy           uint64         `json:"policy"`                            // The policy for proposer selection
	Ceil2Nby3Block           *big.Int       `json:"ceil2Nby3Block,omitempty"`          // Number of confirmations required to move from one state to next [2F + 1 to Ceil(2N/3)]
	ValidatorContractAddress common.Address `json:"validatorcontractaddress"`          // Smart contract address for list of validators
}

type IBFTConfig struct {
	*BFTConfig
}

func (c IBFTConfig) String() string {
	return "istanbul"
}

type QBFTConfig struct {
	*BFTConfig
	BlockReward            *math.HexOrDecimal256 `json:"blockReward,omitempty"`            // Reward from start, works only on QBFT consensus protocol
	BeneficiaryMode        *string               `json:"beneficiaryMode,omitempty"`        // Mode for setting the beneficiary, either: list, besu, validators (beneficiary list is the list of validators)
	BeneficiaryList        []common.Address      `json:"beneficiaryList,omitempty"`        // List of wallet addresses that have benefit at every new block (list mode)
	MiningBeneficiary      *common.Address       `json:"miningBeneficiary,omitempty"`      // Wallet address that benefits at every new block (besu mode)
	ValidatorSelectionMode *string               `json:"validatorselectionmode,omitempty"` // Select model for validators
	Validators             []common.Address      `json:"validators"`                       // Validators list
}

func (c QBFTConfig) String() string {
	return QBFT
}

const (
	IBFT = "ibft"
	QBFT = "qbft"

	ContractMode    = "contract"
	BlockHeaderMode = "blockheader"
)

type Transition struct {
	Block                        *big.Int              `json:"block"`
	Algorithm                    string                `json:"algorithm,omitempty"`
	EpochLength                  uint64                `json:"epochlength,omitempty"`                  // Number of blocks that should pass before pending validator votes are reset
	BlockPeriodSeconds           uint64                `json:"blockperiodseconds,omitempty"`           // Minimum time between two consecutive IBFT or QBFT blocks’ timestamps in seconds
	EmptyBlockPeriodSeconds      *uint64               `json:"emptyblockperiodseconds,omitempty"`      // Minimum time between two consecutive IBFT or QBFT a block and empty block’ timestamps in seconds
	RequestTimeoutSeconds        uint64                `json:"requesttimeoutseconds,omitempty"`        // Minimum request timeout for each IBFT or QBFT round in milliseconds
	ContractSizeLimit            uint64                `json:"contractsizelimit,omitempty"`            // Maximum smart contract code size
	ValidatorContractAddress     common.Address        `json:"validatorcontractaddress"`               // Smart contract address for list of validators
	Validators                   []common.Address      `json:"validators"`                             // List of validators
	ValidatorSelectionMode       string                `json:"validatorselectionmode,omitempty"`       // Validator selection mode to switch to
	EnhancedPermissioningEnabled *bool                 `json:"enhancedPermissioningEnabled,omitempty"` // aka QIP714Block
	PrivacyEnhancementsEnabled   *bool                 `json:"privacyEnhancementsEnabled,omitempty"`   // privacy enhancements (mandatory party, private state validation)
	PrivacyPrecompileEnabled     *bool                 `json:"privacyPrecompileEnabled,omitempty"`     // enable marker transactions support
	GasPriceEnabled              *bool                 `json:"gasPriceEnabled,omitempty"`              // enable gas price
	MinerGasLimit                uint64                `json:"miner.gaslimit,omitempty"`               // Gas Limit
	TwoFPlusOneEnabled           *bool                 `json:"2FPlus1Enabled,omitempty"`               // Ceil(2N/3) is the default you need to explicitly use 2F + 1
	TransactionSizeLimit         uint64                `json:"transactionSizeLimit,omitempty"`         // Modify TransactionSizeLimit
	BlockReward                  *math.HexOrDecimal256 `json:"blockReward,omitempty"`                  // validation rewards
	BeneficiaryMode              *string               `json:"beneficiaryMode,omitempty"`              // Mode for setting the beneficiary, either: list, besu, validators (beneficiary list is the list of validators)
	BeneficiaryList              []common.Address      `json:"beneficiaryList,omitempty"`              // List of wallet addresses that have benefit at every new block (list mode)
	MiningBeneficiary            *common.Address       `json:"miningBeneficiary,omitempty"`            // Wallet address that benefits at every new block (besu mode)
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Ethash != nil:
		engine = c.Ethash
	case c.Clique != nil:
		engine = c.Clique
	case c.Istanbul != nil:
		engine = c.Istanbul
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{ChainID: %v Homestead: %v DAO: %v DAOSupport: %v EIP150: %v EIP155: %v EIP158: %v Byzantium: %v IsQuorum: %v Constantinople: %v TransactionSizeLimit: %v MaxCodeSize: %v Petersburg: %v Istanbul: %v, Muir Glacier: %v, Berlin: %v  Catalyst: %v YOLO v3: %v PrivacyEnhancements: %v PrivacyPrecompile: %v EnableGasPriceBlock: %v Engine: %v}",
		c.ChainID,
		c.HomesteadBlock,
		c.DAOForkBlock,
		c.DAOForkSupport,
		c.EIP150Block,
		c.EIP155Block,
		c.EIP158Block,
		c.ByzantiumBlock,
		c.IsQuorum,
		c.ConstantinopleBlock,
		c.TransactionSizeLimit,
		c.MaxCodeSize,
		c.PetersburgBlock,
		c.IstanbulBlock,
		c.MuirGlacierBlock,
		c.BerlinBlock,
		c.CatalystBlock,
		c.YoloV3Block,
		c.PrivacyEnhancementsBlock, //Quorum
		c.PrivacyPrecompileBlock,   //Quorum
		c.EnableGasPriceBlock,      //Quorum
		engine,
	)
}

// Quorum - validate code size and transaction size limit
func (c *ChainConfig) IsValid() error {

	if c.TransactionSizeLimit < 32 || c.TransactionSizeLimit > 128 {
		return errors.New("Genesis transaction size limit must be between 32 and 128")
	}

	if c.MaxCodeSize != 0 && (c.MaxCodeSize < 24 || c.MaxCodeSize > 128) {
		return errors.New("Genesis max code size must be between 24 and 128")
	}

	return nil
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *ChainConfig) IsHomestead(num *big.Int) bool {
	return isForked(c.HomesteadBlock, num)
}

// IsDAOFork returns whether num is either equal to the DAO fork block or greater.
func (c *ChainConfig) IsDAOFork(num *big.Int) bool {
	return isForked(c.DAOForkBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return isForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return isForked(c.EIP155Block, num)
}

// IsEIP158 returns whether num is either equal to the EIP158 fork block or greater.
func (c *ChainConfig) IsEIP158(num *big.Int) bool {
	return isForked(c.EIP158Block, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *ChainConfig) IsByzantium(num *big.Int) bool {
	return isForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *ChainConfig) IsConstantinople(num *big.Int) bool {
	return isForked(c.ConstantinopleBlock, num)
}

// IsMuirGlacier returns whether num is either equal to the Muir Glacier (EIP-2384) fork block or greater.
func (c *ChainConfig) IsMuirGlacier(num *big.Int) bool {
	return isForked(c.MuirGlacierBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return isForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && isForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return isForked(c.IstanbulBlock, num)
}

// IsBerlin returns whether num is either equal to the Berlin fork block or greater.
func (c *ChainConfig) IsBerlin(num *big.Int) bool {
	return isForked(c.BerlinBlock, num) || isForked(c.YoloV3Block, num)
}

// IsCatalyst returns whether num is either equal to the Merge fork block or greater.
func (c *ChainConfig) IsCatalyst(num *big.Int) bool {
	return isForked(c.CatalystBlock, num)
}

// IsEWASM returns whether num represents a block number after the EWASM fork
func (c *ChainConfig) IsEWASM(num *big.Int) bool {
	return isForked(c.EWASMBlock, num)
}

// Quorum
//
// IsQIP714 returns whether num represents a block number where permissions is enabled
func (c *ChainConfig) IsQIP714(num *big.Int) bool {
	enableEnhancedPermissioning := false
	c.GetTransitionValue(num, func(transition Transition) {
		if transition.EnhancedPermissioningEnabled != nil {
			enableEnhancedPermissioning = *transition.EnhancedPermissioningEnabled
		}
	})
	return isForked(c.QIP714Block, num) || enableEnhancedPermissioning
}

// Quorum
//
// GetMaxCodeSize returns maxCodeSize for the given block number
func (c *ChainConfig) GetMaxCodeSize(num *big.Int) int {
	maxCodeSize := MaxCodeSize

	if len(c.MaxCodeSizeConfig) > 0 {
		log.Warn("WARNING: The attribute config.maxCodeSizeConfig is deprecated and will be removed in the future, please use config.transitions.contractsizelimit on genesis file")
		for _, data := range c.MaxCodeSizeConfig {
			if data.Block.Cmp(num) > 0 {
				break
			}
			maxCodeSize = int(data.Size) * 1024
		}
	} else if c.MaxCodeSize > 0 {
		if c.MaxCodeSizeChangeBlock != nil && c.MaxCodeSizeChangeBlock.Cmp(big.NewInt(0)) >= 0 {
			if isForked(c.MaxCodeSizeChangeBlock, num) {
				maxCodeSize = int(c.MaxCodeSize) * 1024
			}
		} else {
			maxCodeSize = int(c.MaxCodeSize) * 1024
		}
	}

	c.GetTransitionValue(num, func(transition Transition) {
		if transition.ContractSizeLimit != 0 {
			maxCodeSize = int(transition.ContractSizeLimit) * 1024
		}
	})

	return maxCodeSize
}

// Quorum
// gets value at or after a transition
func (c *ChainConfig) GetTransitionValue(num *big.Int, callback func(transition Transition)) {
	if c != nil && num != nil && c.Transitions != nil {
		for i := 0; i < len(c.Transitions) && c.Transitions[i].Block.Cmp(num) <= 0; i++ {
			callback(c.Transitions[i])
		}
	}
}

// Quorum
//
// GetMinerMinGasLimit returns the miners minGasLimit for the given block number
func (c *ChainConfig) GetMinerMinGasLimit(num *big.Int, defaultValue uint64) uint64 {
	minGasLimit := defaultValue
	if c != nil && num != nil && len(c.Transitions) > 0 {
		for i := 0; i < len(c.Transitions) && c.Transitions[i].Block.Cmp(num) <= 0; i++ {
			if c.Transitions[i].MinerGasLimit != 0 {
				minGasLimit = c.Transitions[i].MinerGasLimit
			}
		}
	}
	return minGasLimit
}

// Quorum
//
// validates the maxCodeSizeConfig data passed in config
func (c *ChainConfig) CheckMaxCodeConfigData() error {
	if c.MaxCodeSize != 0 || (c.MaxCodeSizeChangeBlock != nil && c.MaxCodeSizeChangeBlock.Cmp(big.NewInt(0)) >= 0) {
		return errors.New("maxCodeSize & maxCodeSizeChangeBlock deprecated. Consider using maxCodeSizeConfig")
	}
	// validate max code size data
	// 1. Code size should not be less than 24 and greater than 128
	// 2. block entries are in ascending order
	prevBlock := big.NewInt(0)
	for _, data := range c.MaxCodeSizeConfig {
		if data.Size < 24 || data.Size > 128 {
			return errors.New("Genesis max code size must be between 24 and 128")
		}
		if data.Block == nil {
			return errors.New("Block number not given in maxCodeSizeConfig data")
		}
		if data.Block.Cmp(prevBlock) < 0 {
			return errors.New("invalid maxCodeSize detail, block order has to be ascending")
		}
		prevBlock = data.Block
	}

	return nil
}

func (c *ChainConfig) CheckTransitionsData() error {
	isQBFT := false
	if c.QBFT != nil {
		isQBFT = true
	}
	prevBlock := big.NewInt(0)
	for _, transition := range c.Transitions {
		if transition.Algorithm != "" && !strings.EqualFold(transition.Algorithm, IBFT) && !strings.EqualFold(transition.Algorithm, QBFT) {
			return ErrTransitionAlgorithm
		}
		if transition.ValidatorSelectionMode != "" && transition.ValidatorSelectionMode != ContractMode && transition.ValidatorSelectionMode != BlockHeaderMode {
			return ErrValidatorSelectionMode
		}
		if c.Istanbul != nil && c.Istanbul.TestQBFTBlock != nil && (strings.EqualFold(transition.Algorithm, IBFT) || strings.EqualFold(transition.Algorithm, QBFT)) {
			return ErrTestQBFTBlockAndTransitions
		}
		if len(c.MaxCodeSizeConfig) > 0 && transition.ContractSizeLimit != 0 {
			return ErrMaxCodeSizeConfigAndTransitions
		}
		if strings.EqualFold(transition.Algorithm, QBFT) {
			isQBFT = true
		}
		if transition.Block == nil {
			return ErrBlockNumberMissing
		}
		if transition.Block.Cmp(prevBlock) < 0 {
			return ErrBlockOrder
		}
		if isQBFT && strings.EqualFold(transition.Algorithm, IBFT) {
			return ErrTransition
		}
		if transition.ContractSizeLimit != 0 && (transition.ContractSizeLimit < 24 || transition.ContractSizeLimit > 128) {
			return ErrContractSizeLimit
		}
		if transition.ValidatorContractAddress != (common.Address{}) && transition.ValidatorSelectionMode != ContractMode {
			return ErrMissingValidatorSelectionMode
		}
		if transition.TransactionSizeLimit != 0 && transition.TransactionSizeLimit < 32 || transition.TransactionSizeLimit > 128 {
			return ErrTransactionSizeLimit
		}
		if transition.BeneficiaryMode != nil && *transition.BeneficiaryMode != "fixed" && *transition.BeneficiaryMode != "validators" && *transition.BeneficiaryMode != "" && *transition.BeneficiaryMode != "list" {
			return ErrBeneficiaryMode
		}
		prevBlock = transition.Block
	}
	return nil
}

// Quorum
//
// checks if changes to maxCodeSizeConfig proposed are compatible
// with already existing genesis data
func isMaxCodeSizeConfigCompatible(c1, c2 *ChainConfig, head *big.Int) (error, *big.Int, *big.Int) {
	if len(c1.MaxCodeSizeConfig) == 0 && len(c2.MaxCodeSizeConfig) == 0 {
		// maxCodeSizeConfig not used. return
		return nil, big.NewInt(0), big.NewInt(0)
	}

	// existing config had maxCodeSizeConfig and new one does not have the same return error
	if len(c1.MaxCodeSizeConfig) > 0 && len(c2.MaxCodeSizeConfig) == 0 {
		return fmt.Errorf("genesis file missing max code size information"), head, head
	}

	if len(c2.MaxCodeSizeConfig) > 0 && len(c1.MaxCodeSizeConfig) == 0 {
		return nil, big.NewInt(0), big.NewInt(0)
	}

	// check the number of records below current head in both configs
	// if they do not match throw an error
	c1RecsBelowHead := 0
	for _, data := range c1.MaxCodeSizeConfig {
		if data.Block.Cmp(head) <= 0 {
			c1RecsBelowHead++
		} else {
			break
		}
	}

	c2RecsBelowHead := 0
	for _, data := range c2.MaxCodeSizeConfig {
		if data.Block.Cmp(head) <= 0 {
			c2RecsBelowHead++
		} else {
			break
		}
	}

	// if the count of past records is not matching return error
	if c1RecsBelowHead != c2RecsBelowHead {
		return errors.New("maxCodeSizeConfig data incompatible. updating maxCodeSize for past"), head, head
	}

	// validate that each past record is matching exactly. if not return error
	for i := 0; i < c1RecsBelowHead; i++ {
		if c1.MaxCodeSizeConfig[i].Block.Cmp(c2.MaxCodeSizeConfig[i].Block) != 0 ||
			c1.MaxCodeSizeConfig[i].Size != c2.MaxCodeSizeConfig[i].Size {
			return errors.New("maxCodeSizeConfig data incompatible. maxCodeSize historical data does not match"), head, head
		}
	}

	return nil, big.NewInt(0), big.NewInt(0)
}

// Quorum
//
// checks if changes to transitions proposed are compatible
// with already existing genesis data
func isTransitionsConfigCompatible(c1, c2 *ChainConfig, head *big.Int) (error, *big.Int, *big.Int) {
	if len(c1.Transitions) == 0 && len(c2.Transitions) == 0 {
		// maxCodeSizeConfig not used. return
		return nil, big.NewInt(0), big.NewInt(0)
	}

	// existing config had Transitions and new one does not have the same return error
	if len(c1.Transitions) > 0 && len(c2.Transitions) == 0 {
		return fmt.Errorf("genesis file missing transitions information"), head, head
	}

	if len(c2.Transitions) > 0 && len(c1.Transitions) == 0 {
		return nil, big.NewInt(0), big.NewInt(0)
	}

	// check the number of records below current head in both configs
	// if they do not match throw an error
	c1RecsBelowHead := 0
	for _, data := range c1.Transitions {
		if data.Block.Cmp(head) <= 0 {
			c1RecsBelowHead++
		} else {
			break
		}
	}

	c2RecsBelowHead := 0
	for _, data := range c2.Transitions {
		if data.Block.Cmp(head) <= 0 {
			c2RecsBelowHead++
		} else {
			break
		}
	}

	// if the count of past records is not matching return error
	if c1RecsBelowHead != c2RecsBelowHead {
		return errors.New("transitions data incompatible. updating transitions for past"), head, head
	}

	// validate that each past record is matching exactly. if not return error
	for i := 0; i < c1RecsBelowHead; i++ {
		isSameBlock := c1.Transitions[i].Block.Cmp(c2.Transitions[i].Block) != 0
		if isSameBlock || c1.Transitions[i].Algorithm != c2.Transitions[i].Algorithm {
			return ErrTransitionIncompatible("Algorithm"), head, head
		}
		if isSameBlock || c1.Transitions[i].BlockPeriodSeconds != c2.Transitions[i].BlockPeriodSeconds {
			return ErrTransitionIncompatible("BlockPeriodSeconds"), head, head
		}
		if isSameBlock || c1.Transitions[i].RequestTimeoutSeconds != c2.Transitions[i].RequestTimeoutSeconds {
			return ErrTransitionIncompatible("RequestTimeoutSeconds"), head, head
		}
		if isSameBlock || c1.Transitions[i].EpochLength != c2.Transitions[i].EpochLength {
			return ErrTransitionIncompatible("EpochLength"), head, head
		}
		if isSameBlock || c1.Transitions[i].ContractSizeLimit != c2.Transitions[i].ContractSizeLimit {
			return ErrTransitionIncompatible("ContractSizeLimit"), head, head
		}
		if isSameBlock || c1.Transitions[i].ValidatorContractAddress != c2.Transitions[i].ValidatorContractAddress {
			return ErrTransitionIncompatible("ValidatorContractAddress"), head, head
		}
		if isSameBlock || c1.Transitions[i].ValidatorSelectionMode != c2.Transitions[i].ValidatorSelectionMode {
			return ErrTransitionIncompatible("ValidatorSelectionMode"), head, head
		}
		if isSameBlock || c1.Transitions[i].MinerGasLimit != c2.Transitions[i].MinerGasLimit {
			return ErrTransitionIncompatible("Miner GasLimit"), head, head
		}
		if isSameBlock || c1.Transitions[i].TwoFPlusOneEnabled != c2.Transitions[i].TwoFPlusOneEnabled {
			return ErrTransitionIncompatible("2FPlus1Enabled"), head, head
		}
		if isSameBlock || c1.Transitions[i].MinerGasLimit != c2.Transitions[i].MinerGasLimit {
			return ErrTransitionIncompatible("TransactionSizeLimit"), head, head
		}
	}

	return nil, big.NewInt(0), big.NewInt(0)
}

// Quorum
//
// IsPrivacyEnhancementsEnabled returns whether num represents a block number after the PrivacyEnhancementsEnabled fork
func (c *ChainConfig) IsPrivacyEnhancementsEnabled(num *big.Int) bool {
	isPrivacyEnhancementsEnabled := false
	c.GetTransitionValue(num, func(transition Transition) {
		if transition.PrivacyEnhancementsEnabled != nil {
			isPrivacyEnhancementsEnabled = *transition.PrivacyEnhancementsEnabled
		}
	})

	return isForked(c.PrivacyEnhancementsBlock, num) || isPrivacyEnhancementsEnabled
}

// Quorum
//
// Check whether num represents a block number after the PrivacyPrecompileBlock
func (c *ChainConfig) IsPrivacyPrecompileEnabled(num *big.Int) bool {
	isPrivacyPrecompileEnabled := false
	c.GetTransitionValue(num, func(transition Transition) {
		if transition.PrivacyPrecompileEnabled != nil {
			isPrivacyPrecompileEnabled = *transition.PrivacyPrecompileEnabled
		}
	})

	return isForked(c.PrivacyPrecompileBlock, num) || isPrivacyPrecompileEnabled
}

// Quorum
func (c *ChainConfig) GetTransactionSizeLimit(num *big.Int) uint64 {
	transactionSizeLimit := uint64(0)
	c.GetTransitionValue(num, func(transition Transition) {
		transactionSizeLimit = transition.TransactionSizeLimit
	})

	if transactionSizeLimit == 0 {
		transactionSizeLimit = c.TransactionSizeLimit
	}

	if transactionSizeLimit == 0 {
		transactionSizeLimit = 64
	}

	return transactionSizeLimit
}

// Quorum
//
// Check whether num represents a block number after the EnableGasPriceBlock
func (c *ChainConfig) IsGasPriceEnabled(num *big.Int) bool {
	isGasEnabled := false
	c.GetTransitionValue(num, func(transition Transition) {
		if transition.GasPriceEnabled != nil {
			isGasEnabled = *transition.GasPriceEnabled
		}
	})

	return isForked(c.EnableGasPriceBlock, num) || isGasEnabled
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64, isQuorumEIP155Activated bool) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// check if the maxCodesize data passed is compatible 1st
	// this is being handled separately as it can have breaks
	// at multiple block heights and cannot be handled with in
	// checkCompatible

	// compare the maxCodeSize data between the old and new config
	err, cBlock, newCfgBlock := isMaxCodeSizeConfigCompatible(c, newcfg, bhead)
	if err != nil {
		return newCompatError(err.Error(), cBlock, newCfgBlock)
	}
	// compare the transitions data between the old and new config
	err, cBlock, newCfgBlock = isTransitionsConfigCompatible(c, newcfg, bhead)
	if err != nil {
		return newCompatError(err.Error(), cBlock, newCfgBlock)
	}

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead, isQuorumEIP155Activated)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name     string
		block    *big.Int
		optional bool // if true, the fork may be nil and next fork is still allowed
	}
	var lastFork fork
	for _, cur := range []fork{
		{name: "homesteadBlock", block: c.HomesteadBlock},
		{name: "daoForkBlock", block: c.DAOForkBlock, optional: true},
		{name: "eip150Block", block: c.EIP150Block},
		{name: "eip155Block", block: c.EIP155Block},
		{name: "eip158Block", block: c.EIP158Block},
		{name: "byzantiumBlock", block: c.ByzantiumBlock},
		{name: "constantinopleBlock", block: c.ConstantinopleBlock},
		{name: "petersburgBlock", block: c.PetersburgBlock},
		{name: "istanbulBlock", block: c.IstanbulBlock},
		{name: "muirGlacierBlock", block: c.MuirGlacierBlock, optional: true},
		{name: "berlinBlock", block: c.BerlinBlock},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int, isQuorumEIP155Activated bool) *ConfigCompatError {
	if isForkIncompatible(c.HomesteadBlock, newcfg.HomesteadBlock, head) {
		return newCompatError("Homestead fork block", c.HomesteadBlock, newcfg.HomesteadBlock)
	}
	if isForkIncompatible(c.DAOForkBlock, newcfg.DAOForkBlock, head) {
		return newCompatError("DAO fork block", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if c.IsDAOFork(head) && c.DAOForkSupport != newcfg.DAOForkSupport {
		return newCompatError("DAO fork support flag", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if isForkIncompatible(c.EIP150Block, newcfg.EIP150Block, head) {
		return newCompatError("EIP150 fork block", c.EIP150Block, newcfg.EIP150Block)
	}
	if isQuorumEIP155Activated && c.ChainID != nil && isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, head) {
		return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isQuorumEIP155Activated && c.ChainID != nil && c.IsEIP155(head) && !configNumEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP155 chain ID", c.ChainID, newcfg.ChainID)
	}
	if isForkIncompatible(c.EIP158Block, newcfg.EIP158Block, head) {
		return newCompatError("EIP158 fork block", c.EIP158Block, newcfg.EIP158Block)
	}
	if c.IsEIP158(head) && !configNumEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP158 chain ID", c.EIP158Block, newcfg.EIP158Block)
	}
	if isForkIncompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, head) {
		return newCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	}
	if isForkIncompatible(c.ConstantinopleBlock, newcfg.ConstantinopleBlock, head) {
		return newCompatError("Constantinople fork block", c.ConstantinopleBlock, newcfg.ConstantinopleBlock)
	}
	if isForkIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, head) {
		// the only case where we allow Petersburg to be set in the past is if it is equal to Constantinople
		// mainly to satisfy fork ordering requirements which state that Petersburg fork be set if Constantinople fork is set
		if isForkIncompatible(c.ConstantinopleBlock, newcfg.PetersburgBlock, head) {
			return newCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
		}
	}
	if isForkIncompatible(c.IstanbulBlock, newcfg.IstanbulBlock, head) {
		return newCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if isForkIncompatible(c.MuirGlacierBlock, newcfg.MuirGlacierBlock, head) {
		return newCompatError("Muir Glacier fork block", c.MuirGlacierBlock, newcfg.MuirGlacierBlock)
	}
	if isForkIncompatible(c.BerlinBlock, newcfg.BerlinBlock, head) {
		return newCompatError("Berlin fork block", c.BerlinBlock, newcfg.BerlinBlock)
	}
	if isForkIncompatible(c.YoloV3Block, newcfg.YoloV3Block, head) {
		return newCompatError("YOLOv3 fork block", c.YoloV3Block, newcfg.YoloV3Block)
	}
	if isForkIncompatible(c.EWASMBlock, newcfg.EWASMBlock, head) {
		return newCompatError("ewasm fork block", c.EWASMBlock, newcfg.EWASMBlock)
	}
	if c.Istanbul != nil && newcfg.Istanbul != nil && isForkIncompatible(c.Istanbul.Ceil2Nby3Block, newcfg.Istanbul.Ceil2Nby3Block, head) {
		return newCompatError("Ceil 2N/3 fork block", c.Istanbul.Ceil2Nby3Block, newcfg.Istanbul.Ceil2Nby3Block)
	}
	if c.Istanbul != nil && newcfg.Istanbul != nil && isForkIncompatible(c.Istanbul.TestQBFTBlock, newcfg.Istanbul.TestQBFTBlock, head) {
		return newCompatError("Test QBFT fork block", c.Istanbul.TestQBFTBlock, newcfg.Istanbul.TestQBFTBlock)
	}
	if isForkIncompatible(c.QIP714Block, newcfg.QIP714Block, head) {
		return newCompatError("permissions fork block", c.QIP714Block, newcfg.QIP714Block)
	}
	if newcfg.MaxCodeSizeChangeBlock != nil && isForkIncompatible(c.MaxCodeSizeChangeBlock, newcfg.MaxCodeSizeChangeBlock, head) {
		return newCompatError("max code size change fork block", c.MaxCodeSizeChangeBlock, newcfg.MaxCodeSizeChangeBlock)
	}
	if isForkIncompatible(c.PrivacyEnhancementsBlock, newcfg.PrivacyEnhancementsBlock, head) {
		return newCompatError("Privacy Enhancements fork block", c.PrivacyEnhancementsBlock, newcfg.PrivacyEnhancementsBlock)
	}
	if isForkIncompatible(c.PrivacyPrecompileBlock, newcfg.PrivacyPrecompileBlock, head) {
		return newCompatError("Privacy Precompile fork block", c.PrivacyPrecompileBlock, newcfg.PrivacyPrecompileBlock)
	}
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool
	IsBerlin, IsCatalyst                                    bool
	// Quorum
	IsPrivacyEnhancementsEnabled bool
	IsPrivacyPrecompile          bool
	IsGasPriceEnabled            bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:          new(big.Int).Set(chainID),
		IsHomestead:      c.IsHomestead(num),
		IsEIP150:         c.IsEIP150(num),
		IsEIP155:         c.IsEIP155(num),
		IsEIP158:         c.IsEIP158(num),
		IsByzantium:      c.IsByzantium(num),
		IsConstantinople: c.IsConstantinople(num),
		IsPetersburg:     c.IsPetersburg(num),
		IsIstanbul:       c.IsIstanbul(num),
		IsBerlin:         c.IsBerlin(num),
		IsCatalyst:       c.IsCatalyst(num),
		// Quorum
		IsPrivacyEnhancementsEnabled: c.IsPrivacyEnhancementsEnabled(num),
		IsPrivacyPrecompile:          c.IsPrivacyPrecompileEnabled(num),
		IsGasPriceEnabled:            c.IsGasPriceEnabled(num),
	}
}
