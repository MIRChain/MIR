package extensionContracts

import (
	"github.com/pavelkrolevets/MIR-pro/core/state"
)

type AccountWithMetadata struct {
	State state.DumpAccount `json:"state"`
}
