package testutils

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"
	"os/exec"
	"strings"
	"testing"
)

func BuildContract(t *testing.T, contractPath string) {
	logger := logger.Test(t)
	t.Helper()

	logger.Infow("Building contract", "path", contractPath)

	cmd := exec.Command("sui", "move", "build", "--path", contractPath, "--dev")
	logger.Debugw("Executing build command", "command", cmd.String())

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build contract: %s", string(output))
}

func PublishContract(t *testing.T, contractPath string, gasBudget *int) string {
	logger := logger.Test(t)
	t.Helper()

	logger.Infow("Publishing contract", "path", contractPath)

	gasBudgetArg := "200000000"
	if gasBudget != nil {
		gasBudgetArg = string(rune(*gasBudget))
	}

	publishCmd := exec.Command("sui", "client", "publish",
		"--gas-budget", gasBudgetArg,
		"--dev", contractPath)

	publishOutput, err := publishCmd.CombinedOutput()
	require.NoError(t, err, "Failed to publish contract: %s", string(publishOutput))

	lines := strings.Split(string(publishOutput), "\n")
	var packageId string
	for _, line := range lines {
		if strings.Contains(line, "PackageID:") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "0x") {
					packageId = part
					break
				}
			}
		}
		if packageId != "" {
			break
		}
	}
	require.NotEmpty(t, packageId, "Failed to extract packageId from publish output")
	logger.Debugw("Published contract", "packageID", packageId)

	return packageId
}
