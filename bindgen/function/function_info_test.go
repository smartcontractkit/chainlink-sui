package function

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractFunctionInfo(t *testing.T) {
	packageNames, err := getPackagesWithFunctionInfo("../../bindings/generated")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	require.Greater(t, len(packageNames), 0)
}
