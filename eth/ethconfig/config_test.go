package ethconfig

import (
	"math/big"
	"testing"

	"github.com/pavelkrolevets/MIR-pro/consensus/istanbul"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/stretchr/testify/assert"
)

func TestSetBFT(t *testing.T) {
	config := istanbul.DefaultConfig
	bftConfig := &params.BFTConfig{
		EpochLength:           10000,
		Ceil2Nby3Block:        big.NewInt(10),
		RequestTimeoutSeconds: 100,
	}
	setBFTConfig(config, bftConfig)
	assert.Equal(t, config.Ceil2Nby3Block, bftConfig.Ceil2Nby3Block)
	assert.Equal(t, config.Epoch, bftConfig.EpochLength)
	assert.Equal(t, config.RequestTimeout, bftConfig.RequestTimeoutSeconds*1000)
	assert.Equal(t, config.BlockPeriod, istanbul.DefaultConfig.BlockPeriod)
	assert.Equal(t, config.EmptyBlockPeriod, istanbul.DefaultConfig.EmptyBlockPeriod)
	assert.Equal(t, config.ProposerPolicy, istanbul.DefaultConfig.ProposerPolicy)

	bftConfig = &params.BFTConfig{
		EpochLength:           10000,
		Ceil2Nby3Block:        big.NewInt(10),
		RequestTimeoutSeconds: 100,
		BlockPeriodSeconds:    5,
	}
	setBFTConfig(config, bftConfig)
	assert.Equal(t, config.Ceil2Nby3Block, bftConfig.Ceil2Nby3Block)
	assert.Equal(t, config.Epoch, bftConfig.EpochLength)
	assert.Equal(t, config.RequestTimeout, bftConfig.RequestTimeoutSeconds*1000)
	assert.Equal(t, config.BlockPeriod, uint64(5))
	assert.Equal(t, config.EmptyBlockPeriod, uint64(0))
	assert.Equal(t, config.ProposerPolicy, istanbul.DefaultConfig.ProposerPolicy)

	bftConfig = &params.BFTConfig{
		EpochLength:           10000,
		Ceil2Nby3Block:        big.NewInt(10),
		RequestTimeoutSeconds: 100,
	}
	setBFTConfig(config, bftConfig)
	assert.Equal(t, config.Ceil2Nby3Block, bftConfig.Ceil2Nby3Block)
	assert.Equal(t, config.Epoch, bftConfig.EpochLength)
	assert.Equal(t, config.RequestTimeout, bftConfig.RequestTimeoutSeconds*1000)
	assert.Equal(t, config.BlockPeriod, istanbul.DefaultConfig.BlockPeriod)
	assert.Equal(t, config.EmptyBlockPeriod, uint64(0))
	assert.Equal(t, config.ProposerPolicy, istanbul.DefaultConfig.ProposerPolicy)
}
