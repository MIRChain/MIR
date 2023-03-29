package testutils

import (
	"bytes"

	"github.com/pavelkrolevets/MIR-pro/common"
	istanbulcommon "github.com/pavelkrolevets/MIR-pro/consensus/istanbul/common"
	"github.com/pavelkrolevets/MIR-pro/core"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/gost3410"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/pavelkrolevets/MIR-pro/rlp"
)

func Genesis[P crypto.PublicKey](validators []common.Address, isQBFT bool) *core.Genesis[P] {
	// generate genesis block
	genesis := core.DefaultGenesisBlock[P]()
	genesis.Config = params.TestChainConfig
	// force enable Istanbul engine
	genesis.Config.Istanbul = &params.IstanbulConfig{}
	genesis.Config.Ethash = nil
	genesis.Difficulty = istanbulcommon.DefaultDifficulty
	genesis.Nonce = istanbulcommon.EmptyBlockNonce.Uint64()
	genesis.Mixhash = types.IstanbulDigest

	if isQBFT {
		appendValidators(genesis, validators)
	} else {
		appendValidatorsIstanbulExtra(genesis, validators)
	}

	return genesis
}

func GenesisAndKeys[T crypto.PrivateKey, P crypto.PublicKey](n int, isQBFT bool) (*core.Genesis[P], []T) {
	// Setup validators
	var nodeKeys = make([]T, n)
	var addrs = make([]common.Address, n)
	for i := 0; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey[T]()
		var pub P
		switch t:=any(&nodeKeys[i]).(type) {
		case *nist.PrivateKey:
			p:=any(&pub).(*nist.PublicKey)
			*p = *t.Public()
		case *gost3410.PrivateKey:
			p:=any(&pub).(*gost3410.PublicKey)
			*p = *t.Public()
		}
		addrs[i] = crypto.PubkeyToAddress(pub)
	}

	// generate genesis block
	genesis := Genesis[P](addrs, isQBFT)

	return genesis, nodeKeys
}

func appendValidatorsIstanbulExtra[P crypto.PublicKey](genesis *core.Genesis[P], addrs []common.Address) {
	if len(genesis.ExtraData) < types.IstanbulExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.IstanbulExtraVanity]

	ist := &types.IstanbulExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = append(genesis.ExtraData, istPayload...)
}

func appendValidators[P crypto.PublicKey](genesis *core.Genesis[P], addrs []common.Address) {
	vanity := append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(genesis.ExtraData))...)
	ist := &types.QBFTExtra{
		VanityData:    vanity,
		Validators:    addrs,
		Vote:          nil,
		CommittedSeal: [][]byte{},
		Round:         0,
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = istPayload
}
