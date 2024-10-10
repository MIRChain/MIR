package vm

import (
	"testing"

	"github.com/MIRChain/MIR/common"
	"github.com/MIRChain/MIR/crypto/gost3410"
	"github.com/MIRChain/MIR/params"
	"github.com/stretchr/testify/require"
)

func TestActivePrecompilesGost(t *testing.T) {
	tests := []struct {
		name string
		evm  *EVM[gost3410.PublicKey]
		want []common.Address
	}{
		{
			name: "istanbul-plus-quorum-privacy",
			evm: &EVM[gost3410.PublicKey]{
				chainRules: params.Rules{
					IsIstanbul:          true,
					IsPrivacyPrecompile: true,
				},
			},
			want: []common.Address{
				common.BytesToAddress([]byte{1}),
				common.BytesToAddress([]byte{2}),
				common.BytesToAddress([]byte{3}),
				common.BytesToAddress([]byte{4}),
				common.BytesToAddress([]byte{5}),
				common.BytesToAddress([]byte{6}),
				common.BytesToAddress([]byte{7}),
				common.BytesToAddress([]byte{8}),
				common.BytesToAddress([]byte{9}),
				common.QuorumPrivacyPrecompileContractAddress(),
			},
		},
		{
			name: "homestead-plus-quorum-privacy",
			evm: &EVM[gost3410.PublicKey]{
				chainRules: params.Rules{
					IsHomestead:         true,
					IsPrivacyPrecompile: true,
				},
			},
			want: []common.Address{
				common.BytesToAddress([]byte{1}),
				common.BytesToAddress([]byte{2}),
				common.BytesToAddress([]byte{3}),
				common.BytesToAddress([]byte{4}),
				common.QuorumPrivacyPrecompileContractAddress(),
			},
		},
		{
			name: "istanbul",
			evm: &EVM[gost3410.PublicKey]{
				chainRules: params.Rules{
					IsIstanbul:          true,
					IsPrivacyPrecompile: false,
				},
			},
			want: []common.Address{
				common.BytesToAddress([]byte{1}),
				common.BytesToAddress([]byte{2}),
				common.BytesToAddress([]byte{3}),
				common.BytesToAddress([]byte{4}),
				common.BytesToAddress([]byte{5}),
				common.BytesToAddress([]byte{6}),
				common.BytesToAddress([]byte{7}),
				common.BytesToAddress([]byte{8}),
				common.BytesToAddress([]byte{9}),
			},
		},
		{
			name: "homestead",
			evm: &EVM[gost3410.PublicKey]{
				chainRules: params.Rules{
					IsHomestead:         true,
					IsPrivacyPrecompile: false,
				},
			},
			want: []common.Address{
				common.BytesToAddress([]byte{1}),
				common.BytesToAddress([]byte{2}),
				common.BytesToAddress([]byte{3}),
				common.BytesToAddress([]byte{4}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ActivePrecompiles(tt.evm.chainRules)
			require.ElementsMatchf(t, tt.want, got, "want: %v, got: %v", tt.want, got)
		})
	}
}
