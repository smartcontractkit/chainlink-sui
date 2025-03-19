package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

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
			Filter: map[string]interface{}{
				"MatchType": map[string]interface{}{
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
func FindCounterID(t *testing.T, c sui.Client, packageId string) string {
	t.Helper()
	ctx := context.Background()

	// Query for objects owned by the package address
	resp, err := c.SuiXGetOwnedObjects(ctx, models.SuiXGetOwnedObjectsRequest{
		Address: packageId,
		Query: models.SuiObjectResponseQuery{
			Filter: map[string]interface{}{
				"MatchType": map[string]interface{}{
					"TypeName": packageId + "::counter::Counter",
				},
			},
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
			},
		},
	})
	require.NoError(t, err)
	require.Greater(t, len(resp.Data), 0, "No counter objects found")

	return resp.Data[0].Data.ObjectId
}

// Deploys and initializes the counter contract
// Returns the package ID
func DeployCounterContract(t *testing.T, c sui.Client) (string, error) {
	t.Helper()

	// Compile and publish the counter contract
	packagePath := "/contracts/test/"

	// Get the current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Navigate to the project root (assuming we're in relayer/testutils)
	projectRoot := filepath.Dir(filepath.Dir(cwd))
	contractPath := filepath.Join(projectRoot, packagePath)

	// Build the contract
	cmd := exec.Command("sui", "move", "build", "--path", contractPath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build contract: %s", string(output))

	// Publish the contract
	publishCmd := exec.Command("sui", "client", "publish",
		"--path", contractPath,
		"--gas-budget", "200000000",
		"-d", "--json")

	publishOutput, err := publishCmd.CombinedOutput()
	require.NoError(t, err, "Failed to publish contract: %s", string(publishOutput))

	var publishData map[string]interface{}
	err = json.Unmarshal(publishOutput, &publishData)
	require.NoError(t, err, "Failed to parse publish output as JSON")

	// Extract the package ID from the JSON response
	packageId, ok := publishData["packageId"].(string)
	require.True(t, ok, "Failed to extract packageId from publish output")

	// Initialize the counter
	// Call the init function to create the counter
	// Initialize the counter using the Sui CLI
	initCmd := exec.Command("sui", "client", "call",
		"--package", packageId,
		"--module", "counter",
		"--function", "initialize",
		"--gas-budget", "20000000", "-d", "--json")

	initOutput, err := initCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to initialize counter: %w, output: %s", err, string(initOutput))
	}

	// Wait for the transaction to be confirmed
	time.Sleep(2 * time.Second)

	return packageId, nil
}
