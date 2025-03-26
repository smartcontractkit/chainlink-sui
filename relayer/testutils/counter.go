package testutils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
)

// HasCounterObject checks if a counter object exists for the given packageId.
func HasCounterObject(t *testing.T, c sui.Client, packageId string) bool {
	t.Helper()
	ctx := context.Background()

	// Query for objects owned by the package address
	resp, err := c.SuiXGetOwnedObjects(ctx, models.SuiXGetOwnedObjectsRequest{
		Address: packageId,
		Query: models.SuiObjectResponseQuery{
			Filter: map[string]any{
				"MatchType": map[string]any{
					"TypeName": packageId + "::counter::Counter",
				},
			},
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
			},
		},
	})
	require.NoError(t, err)

	return len(resp.Data) > 0
}

// ReadCounterValue gets the current value of a counter object.
func ReadCounterValue(t *testing.T, c sui.Client, counterId string) int {
	t.Helper()
	ctx := context.Background()

	resp, err := c.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: counterId,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	})
	require.NoError(t, err)

	// Access the fields directly from the SuiMoveObject embedded in SuiParsedData
	value, ok := resp.Data.Content.Fields["value"].(float64)
	require.True(t, ok, "Failed to parse counter value")

	return int(value)
}

// FindCounterID locates the counter object ID for a given package ID.
func FindCounterID(t *testing.T, c sui.ISuiAPI, packageId string) string {
	t.Helper()
	ctx := context.Background()

	// Query for objects owned by the package address
	resp, err := c.SuiXGetOwnedObjects(ctx, models.SuiXGetOwnedObjectsRequest{
		Address: packageId,
		Query: models.SuiObjectResponseQuery{
			Filter: map[string]any{
				"StructType": packageId + "::counter::Counter",
			},
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
			},
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.Data)

	return resp.Data[0].Data.ObjectId
}

// Deploys and initializes the counter contract
// Returns the package ID
// TODO: Depending on the CLI to publish and call contracts causes non-deterministic test failures
//   - Consider using the Sui SDK to publish contracts
func DeployCounterContract(t *testing.T) (string, string, error) {
	t.Helper()

	log := logger.Test(t)
	// Compile and publish the counter contract
	packagePath := "/contracts/test/"

	// Get the current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Navigate to the project root (assuming we're in relayer/testutils)
	projectRoot := filepath.Dir(filepath.Dir(cwd))
	contractPath := filepath.Join(projectRoot, packagePath)

	log.Infow("Deploying counter contract", "path", contractPath)

	// Build the contract
	cmd := exec.Command("sui", "move", "build", "--path", contractPath, "--dev")
	log.Debugw("Executing counter build command", "command", cmd.String())

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build contract: %s", string(output))

	// Publish the contract
	publishCmd := exec.Command("sui", "client", "publish",
		"--gas-budget", "200000000",
		"--dev", contractPath)

	publishOutput, err := publishCmd.CombinedOutput()
	require.NoError(t, err, "Failed to publish contract: %s", string(publishOutput))

	// Extract the package ID from the output
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

	log.Debugw("Published counter contract", "packageID", packageId)

	// Initialize the counter
	// Call the init function to create the counter
	// Initialize the counter using the Sui CLI
	initCmd := exec.Command("sui", "client", "call",
		"--package", packageId,
		"--module", "counter",
		"--function", "initialize",
		"--gas-budget", "20000000")

	log.Debugw("Executing counter init command", "command", initCmd.String())

	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("failed to initialize counter: %w, output: %s", err, string(initOutput))
	}

	// Extract the ObjectID from the initialization output
	lines = strings.Split(string(initOutput), "\n")
	var objectID string
	for _, line := range lines {
		if strings.Contains(line, "ObjectID:") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "0x") {
					objectID = part
					break
				}
			}
		}
		if objectID != "" {
			break
		}
	}
	require.NotEmpty(t, objectID, "Failed to extract ObjectID from initialization output")

	log.Debugw("Initialized counter contract", "objectID", objectID)

	// Wait for the transaction to be confirmed
	// TODO: Implement a more robust way to wait for the transaction to be confirmed
	const DefaultTxConfirmDelay = 2 * time.Second
	time.Sleep(DefaultTxConfirmDelay)

	return packageId, objectID, nil
}
