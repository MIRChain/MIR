package istanbulcommon

import (
	"github.com/pavelkrolevets/MIR-pro/common"
	"github.com/pavelkrolevets/MIR-pro/crypto"
	"github.com/pavelkrolevets/MIR-pro/crypto/nist"
)

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress[P crypto.PublicKey](data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256[nist.PublicKey](data)
	// 2. Recover public key
	pubkey, err := crypto.SigToPub[P](hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(pubkey), nil
}
