package lanes

// LatestPackageIDsConfig holds optional upgraded Sui package IDs for MCMS PTB routing.
// Package IDs in the address book remain the MCMS registry identity; latest IDs route
// execution to upgraded bytecode.
//
// Populated by CLD via WithSuiLatestPackageIDs from durable pipeline YAML resolver output.
type LatestPackageIDsConfig struct {
	OffRamp string
	CCIP    string
	OnRamp  string
	Router  string
}

// resolveLatestPackageIDs returns upgraded package IDs for MCMS execution on chainSelector.
// Empty when no resolver scope is active or the selector has no overrides (pre-upgrade is OK).
func resolveLatestPackageIDs(chainSelector uint64) LatestPackageIDsConfig {
	return currentLatestPackageIDs(chainSelector)
}
