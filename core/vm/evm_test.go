package vm

import (
	"testing"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
	"github.com/pavelkrolevets/MIR-pro/params"
	"github.com/stretchr/testify/require"
)

func TestActivePrecompiles(t *testing.T) {
	tests := []struct {
		name string
		evm  *EVM[nist.PublicKey]
		want []common.Address
	}{
		{
			name: "istanbul-plus-quorum-privacy",
			evm: &EVM[nist.PublicKey]{
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
			evm: &EVM[nist.PublicKey]{
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
			evm: &EVM[nist.PublicKey]{
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
			evm: &EVM[nist.PublicKey]{
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
