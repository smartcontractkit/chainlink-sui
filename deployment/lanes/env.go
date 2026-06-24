package lanes

import (
	"fmt"
	"sync"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// connectChainsEnv carries the active CLDF environment while ConnectChains runs.
// LaneAdapter address getters only receive datastore.DataStore, not Environment, so
// Sui resolves package IDs from addresses.json via LoadOnchainStatesui using this scope.
//
// CLD must invoke ConnectChains inside WithConnectChainsEnvironment. Other chain
// families are unaffected.
var connectChainsEnv struct {
	mu     sync.RWMutex
	env    cldf.Environment
	active bool
}

// WithConnectChainsEnvironment runs fn while SuiAdapter address getters can read
// ExistingAddresses from the given environment.
func WithConnectChainsEnvironment(e cldf.Environment, fn func() error) error {
	connectChainsEnv.mu.Lock()
	connectChainsEnv.env = e
	connectChainsEnv.active = true
	connectChainsEnv.mu.Unlock()

	defer func() {
		connectChainsEnv.mu.Lock()
		connectChainsEnv.active = false
		connectChainsEnv.env = cldf.Environment{}
		connectChainsEnv.mu.Unlock()
	}()

	return fn()
}

func connectChainsEnvironment() (cldf.Environment, error) {
	connectChainsEnv.mu.RLock()
	defer connectChainsEnv.mu.RUnlock()

	if !connectChainsEnv.active {
		return cldf.Environment{}, fmt.Errorf(
			"Sui lane address getters require ConnectChains to run inside lanes.WithConnectChainsEnvironment",
		)
	}
	return connectChainsEnv.env, nil
}
