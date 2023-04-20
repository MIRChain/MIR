package qbftengine

import (
	"github.com/pavelkrolevets/MIR-pro/core/types"
	"github.com/pavelkrolevets/MIR-pro/crypto"
)

type ApplyQBFTExtra func(*types.QBFTExtra) error

func Combine(applies ...ApplyQBFTExtra) ApplyQBFTExtra {
	return func(extra *types.QBFTExtra) error {
		for _, apply := range applies {
			err := apply(extra)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func ApplyHeaderQBFTExtra[P crypto.PublicKey](header *types.Header[P], applies ...ApplyQBFTExtra) error {
	extra, err := getExtra(header)
	if err != nil {
		return err
	}

	err = Combine(applies...)(extra)
	if err != nil {
		return err
	}

	return setExtra(header, extra)
}
