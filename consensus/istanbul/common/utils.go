package istanbulcommon

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256(data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}
