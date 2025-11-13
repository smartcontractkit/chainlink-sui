package common

// PointerConfig defines the configuration for a pointer object type.
// Pointer objects store parent object IDs, and child objects are derived from them.
type PointerConfig struct {
	// Module is the Sui module containing the pointer object
	Module string
	// Pointer is the type name of the pointer object (e.g., "OffRampStatePointer")
	Pointer string
	// ParentFieldName is the field in the pointer object containing the parent object ID
	ParentFieldName string
}

// PointerConfigs maps contract/module names to their pointer configurations.
// This is the single source of truth for all pointer object configurations.
var PointerConfigs = map[string][]PointerConfig{
	"offramp": {
		{
			Module:          "offramp",
			Pointer:         "OffRampStatePointer",
			ParentFieldName: "off_ramp_object_id",
		},
	},
	"onramp": {
		{
			Module:          "onramp",
			Pointer:         "OnRampStatePointer",
			ParentFieldName: "on_ramp_object_id",
		},
	},
	"ccip": {
		{
			Module:          "state_object",
			Pointer:         "CCIPObjectRefPointer",
			ParentFieldName: "ccip_object_id",
		},
	},
	"state_object": {
		{
			Module:          "state_object",
			Pointer:         "CCIPObjectRefPointer",
			ParentFieldName: "ccip_object_id",
		},
	},
	"router": {
		{
			Module:          "router",
			Pointer:         "RouterStatePointer",
			ParentFieldName: "router_object_id",
		},
	},
	"burn_mint_token_pool": {
		{
			Module:          "burn_mint_token_pool",
			Pointer:         "BurnMintTokenPoolStatePointer",
			ParentFieldName: "burn_mint_token_pool_object_id",
		},
	},
	"managed_token_pool": {
		{
			Module:          "managed_token_pool",
			Pointer:         "ManagedTokenPoolStatePointer",
			ParentFieldName: "managed_token_pool_object_id",
		},
	},
	"usdc_token_pool": {
		{
			Module:          "usdc_token_pool",
			Pointer:         "USDCTokenPoolStatePointer",
			ParentFieldName: "usdc_token_pool_object_id",
		},
	},
	"lock_release_token_pool": {
		{
			Module:          "lock_release_token_pool",
			Pointer:         "LockReleaseTokenPoolStatePointer",
			ParentFieldName: "lock_release_token_pool_object_id",
		},
	},
	"counter": { // Test contract
		{
			Module:          "counter",
			Pointer:         "CounterPointer",
			ParentFieldName: "counter_object_id",
		},
	},
}

func GetParentFieldName(pointerName string) string {
	for _, configs := range PointerConfigs {
		for _, config := range configs {
			if config.Pointer == pointerName {
				return config.ParentFieldName
			}
		}
	}
	return ""
}

func GetPointerConfigsByContract(contractName string) []PointerConfig {
	return PointerConfigs[NormalizeName(contractName)]
}
