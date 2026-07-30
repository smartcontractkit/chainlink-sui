package lanes

import (
	"fmt"
	"maps"
	"sync"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// connectChainsEnv carries the active CLDF environment while ConnectChains runs.
// LaneAdapter address getters only receive datastore.DataStore, not Environment, so
// Sui resolves package IDs from LoadOnchainStatesui (datastore address refs first,
// legacy addresses.json fallback) using this scope.
//
// CLD must invoke ConnectChains inside WithConnectChainsEnvironment. Other chain
// families are unaffected.
var connectChainsEnv struct {
	mu     sync.RWMutex
	env    cldf.Environment
	active bool
}

// WithConnectChainsEnvironment runs fn while SuiAdapter address getters can read
// chain metadata from the given environment.
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

var latestPackageIDsEnv struct {
	mu         sync.RWMutex
	bySelector map[uint64]LatestPackageIDsConfig
	active     bool
}

// WithSuiLatestPackageIDs runs fn while lane sequences read optional upgraded package IDs
// keyed by Sui chain selector. CLD resolver populates this map from durable pipeline YAML.
func WithSuiLatestPackageIDs(bySelector map[uint64]LatestPackageIDsConfig, fn func() error) error {
	latestPackageIDsEnv.mu.Lock()
	latestPackageIDsEnv.bySelector = copyLatestPackageIDsBySelector(bySelector)
	latestPackageIDsEnv.active = true
	latestPackageIDsEnv.mu.Unlock()

	defer func() {
		latestPackageIDsEnv.mu.Lock()
		latestPackageIDsEnv.active = false
		latestPackageIDsEnv.bySelector = nil
		latestPackageIDsEnv.mu.Unlock()
	}()

	return fn()
}

// RunConnectChainsWithSuiScopes wraps ConnectChains with address-book and resolver latest-package-ID scopes.
func RunConnectChainsWithSuiScopes(
	e cldf.Environment,
	latestBySelector map[uint64]LatestPackageIDsConfig,
	fn func() error,
) error {
	return WithConnectChainsEnvironment(e, func() error {
		return WithSuiLatestPackageIDs(latestBySelector, fn)
	})
}

func currentLatestPackageIDs(chainSelector uint64) LatestPackageIDsConfig {
	latestPackageIDsEnv.mu.RLock()
	defer latestPackageIDsEnv.mu.RUnlock()
	if !latestPackageIDsEnv.active || latestPackageIDsEnv.bySelector == nil {
		return LatestPackageIDsConfig{}
	}
	return latestPackageIDsEnv.bySelector[chainSelector]
}

func copyLatestPackageIDsBySelector(in map[uint64]LatestPackageIDsConfig) map[uint64]LatestPackageIDsConfig {
	if len(in) == 0 {
		return nil
	}
	out := make(map[uint64]LatestPackageIDsConfig, len(in))
	maps.Copy(out, in)
	return out
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
