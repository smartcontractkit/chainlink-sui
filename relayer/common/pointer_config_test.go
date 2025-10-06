package common //nolint:revive // var-naming: package name is intentionally simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetParentFieldName(t *testing.T) {
	tests := []struct {
		name        string
		pointerName string
		want        string
	}{
		{
			name:        "OffRampStatePointer",
			pointerName: "OffRampStatePointer",
			want:        "off_ramp_object_id",
		},
		{
			name:        "OnRampStatePointer",
			pointerName: "OnRampStatePointer",
			want:        "on_ramp_object_id",
		},
		{
			name:        "CCIPObjectRefPointer",
			pointerName: "CCIPObjectRefPointer",
			want:        "ccip_object_id",
		},
		{
			name:        "RouterStatePointer",
			pointerName: "RouterStatePointer",
			want:        "router_object_id",
		},
		{
			name:        "CounterPointer",
			pointerName: "CounterPointer",
			want:        "counter_object_id",
		},
		{
			name:        "unknown pointer",
			pointerName: "UnknownPointer",
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetParentFieldName(tt.pointerName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPointerConfigsByContract(t *testing.T) {
	tests := []struct {
		name         string
		contractName string
		wantCount    int
		wantModule   string
	}{
		{
			name:         "OffRamp",
			contractName: "OffRamp",
			wantCount:    1,
			wantModule:   "offramp",
		},
		{
			name:         "OnRamp",
			contractName: "OnRamp",
			wantCount:    1,
			wantModule:   "onramp",
		},
		{
			name:         "ccip",
			contractName: "ccip",
			wantCount:    1,
			wantModule:   "state_object",
		},
		{
			name:         "Router",
			contractName: "Router",
			wantCount:    1,
			wantModule:   "router",
		},
		{
			name:         "counter",
			contractName: "counter",
			wantCount:    1,
			wantModule:   "counter",
		},
		{
			name:         "unknown contract",
			contractName: "Unknown",
			wantCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPointerConfigsByContract(tt.contractName)
			assert.Len(t, got, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, tt.wantModule, got[0].Module)
			}
		})
	}
}

func TestPointerConfig_AllConfigs(t *testing.T) {
	// Verify all known pointer configs are properly configured
	for contractName, configs := range PointerConfigs {
		t.Run(contractName, func(t *testing.T) {
			require.NotEmpty(t, configs, "Contract %s should have at least one pointer config", contractName)

			for _, config := range configs {
				assert.NotEmpty(t, config.Module, "Module should not be empty for %s", contractName)
				assert.NotEmpty(t, config.Pointer, "Pointer should not be empty for %s", contractName)
				assert.NotEmpty(t, config.ParentFieldName, "ParentFieldName should not be empty for %s", contractName)

				// Verify the parent field name can be found by pointer name
				foundField := GetParentFieldName(config.Pointer)
				assert.Equal(t, config.ParentFieldName, foundField,
					"GetParentFieldName should return the correct field for %s", config.Pointer)
			}
		})
	}
}
