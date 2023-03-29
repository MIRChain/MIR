package qlight

import (
	"encoding/base64"
	"fmt"

	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/plugin/security"
	"github.com/pavelkrolevets/MIR-pro/private/engine"
	"github.com/pavelkrolevets/MIR-pro/private/engine/qlightptm"
	"github.com/pavelkrolevets/MIR-pro/rlp"
)

type PrivateStateRootHashValidator interface {
	ValidatePrivateStateRoot(blockHash common.Hash, blockPublicStateRoot common.Hash) error
}

type PrivateClientCache interface {
	PrivateStateRootHashValidator
	AddPrivateBlock(blockPrivateData BlockPrivateData) error
	CheckAndAddEmptyEntry(hash common.EncryptedPayloadHash)
}

type PrivateBlockDataResolver [P crypto.PublicKey]interface {
	PrepareBlockPrivateData(block *types.Block[P], psi string) (*BlockPrivateData, error)
}

type AuthManagerProvider func() security.AuthenticationManager

type AuthProvider interface {
	Initialize() error
	Authorize(token string, psi string) error
}

type CacheWithEmpty interface {
	Cache(privateTxData *qlightptm.CachablePrivateTransactionData) error
	CheckAndAddEmptyToCache(hash common.EncryptedPayloadHash)
}

type BlockPrivateData struct {
	BlockHash           common.Hash
	PSI                 types.PrivateStateIdentifier
	PrivateStateRoot    common.Hash
	PrivateTransactions []PrivateTransactionData
}

type QLightCacheKey struct {
	BlockHash common.Hash
	PSI       types.PrivateStateIdentifier
}

func (k *QLightCacheKey) String() string {
	bytes, err := rlp.EncodeToBytes(k)
	if err != nil {
		return err.Error()
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

type PrivateTransactionData struct {
	Hash     *common.EncryptedPayloadHash
	Payload  []byte
	Extra    *engine.ExtraMetadata
	IsSender bool
}

func (d *PrivateTransactionData) ToCachable() *qlightptm.CachablePrivateTransactionData {
	return &qlightptm.CachablePrivateTransactionData{
		Hash: *d.Hash,
		QuorumPrivateTxData: engine.QuorumPayloadExtra{
			Payload:       fmt.Sprintf("0x%x", d.Payload),
			ExtraMetaData: d.Extra,
			IsSender:      d.IsSender,
		},
	}
}
