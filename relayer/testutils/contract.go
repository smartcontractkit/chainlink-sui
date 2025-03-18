package testutils

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"
)

// Response represents the minimal fields expected in the JSON output.
type Response struct {
	ObjectChanges []any `json:"objectChanges"`
}

func BuildSetup(t *testing.T, packagePath string) string {
	t.Helper()
	lgr := logger.Test(t)
	// Get the current working directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// Navigate to the project root (assuming we're in relayer/testutils)
	projectRoot := filepath.Dir(filepath.Dir(cwd))
	contractPath := filepath.Join(projectRoot, packagePath)

	lgr.Debugw("Building contract setup", "path", contractPath)

	return contractPath
}

func cleanJSONOutput(output string) string {
	idx := strings.Index(output, "{")
	if idx == -1 {
		// No JSON object found, return original output.
		return output
	}

	return output[idx:]
}

func BuildContract(t *testing.T, contractPath string) {
	t.Helper()

	lgr := logger.Test(t)

	lgr.Infow("Building contract", "path", contractPath)

	cmd := exec.Command("sui", "move", "build", "--path",
		contractPath,
	)
	lgr.Debugw("Executing build command", "command", cmd.String())

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build contract: %s", string(output))
}

// PublishContract publishes a Move contract to the Sui network and extracts its package ID.
//
// The function constructs and executes a "sui client publish" command using the provided
// contractPath and gasBudget (if specified). It cleans the command output to remove any unwanted
// header text, unmarshals the resulting JSON, and iterates over the "objectChanges" array to find
// an entry of type "published". Once found, it returns the associated packageId along with the full
// cleaned JSON output.
//
// Parameters:
//
//	t            - A testing.T instance for error reporting.
//	contractPath - A string representing the filesystem path to the Move contract.
//	gasBudget    - A pointer to an int that specifies the gas budget for the publish transaction.
//	               If nil, a default value is used.
//
// Returns:
//
//	packageId    - The package ID extracted from the JSON output, typically for a published contract.
//	output       - The cleaned JSON output from the publish command.
//	error        - An error if the publish operation fails or if a valid package ID is not found.
func PublishContract(t *testing.T, contractPath string, gasBudget *int) (string, string, error) {
	t.Helper()
	lgr := logger.Test(t)

	lgr.Infow("Publishing contract", "path", contractPath)

	gasBudgetArg := "200000000"
	if gasBudget != nil {
		gasBudgetArg = string(rune(*gasBudget))
	}

	publishCmd := exec.Command("sui", "client", "publish",
		"--gas-budget", gasBudgetArg,
		"--json",
		contractPath,
	)

	publishOutput, err := publishCmd.CombinedOutput()
	require.NoError(t, err, "Failed to publish contract: %s", string(publishOutput))

	cleanedOutput := cleanJSONOutput(string(publishOutput))

	// Unmarshal the JSON into a map.
	var result map[string]any
	if err := json.Unmarshal([]byte(cleanedOutput), &result); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	changes, ok := result["objectChanges"].([]any)
	if !ok {
		return "", "", errors.New("objectChanges key not found or not a slice")
	}

	var packageId string
	for _, change := range changes {
		m, ok := change.(map[string]any)
		if !ok {
			continue
		}
		if m["type"] == "published" {
			if p, ok := m["packageId"].(string); ok {
				packageId = p
				break
			}
		}
	}

	if packageId == "" {
		return "", "", errors.New("package ID not found")
	}

	return packageId, cleanedOutput, nil
}

// ExtractObjectId parses the JSON output from a Sui publish command and extracts an object identifier
// associated with a specific Move struct name. It expects the JSON to contain an "objectChanges" array,
// which may include various types of changes such as "published" and "created". When a "published"
// entry is present, the function extracts the "packageId", whereas for other types it might extract
// the "objectId" if that’s what’s required.
//
// Parameters:
//
//	t              - A testing.T instance for error reporting.
//	publishOutput  - A string containing the raw JSON output from the Sui publish command.
//	moveStructName - The name of the Move struct to search for (e.g. "TodoList").
//
// Returns:
//
//	A string representing the extracted object identifier (for instance, the packageId for a published object)
//	and an error if the JSON cannot be parsed or no matching object is found.
//
// Example JSON configuration elements that this function processes:
//
//	{
//	     "type": "published",
//	     "packageId": "0x36a176c9b2d99b89e90804870af1584ff244da9723308491b9222f831141c2a6",
//	     "version": "1",
//	     "digest": "DWh8Sy2dbojnGbArjYPgQGdy829Yo3u7G4bvH9UrtJGm",
//	     "modules": [
//	        "cw_tests"
//	     ]
//	},
//
//	{
//	     "type": "created",
//	     "sender": "0x57a33a2fbf908667686407c7dad19590de369054d3d9ce9545af9d80392406a6",
//	     "owner": {
//	          "Shared": {
//	               "initial_shared_version": 3
//	          }
//	     },
//	     "objectType": "0x36a176c9b2d99b89e90804870af1584ff244da9723308491b9222f831141c2a6::cw_tests::TodoList",
//	     "objectId": "0xd525c34d6bc0d4306f16fb5b929be894a333df597019415bc2143e94bc0bc09f",
//	     "version": "3",
//	     "digest": "H8n9xbztGVrBjvBzqq8fnHvADvLRvK8cLeMxDb5SYo8V"
//	}
func ExtractObjectId(t *testing.T, publishOutput string, moveStructName string) (string, error) {
	t.Helper()

	var result map[string]any
	if err := json.Unmarshal([]byte(publishOutput), &result); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	changesAny, ok := result["objectChanges"].([]any)
	if !ok {
		return "", errors.New("objectChanges key not found or not a slice")
	}

	for _, change := range changesAny {
		m, ok := change.(map[string]any)
		if !ok {
			continue
		}

		// Check for a "created" change that contains the target moveStructName.
		if typ, _ok := m["type"].(string); !_ok || typ != "created" {
			continue
		}
		objectType, ok := m["objectType"].(string)
		if !ok || !strings.Contains(objectType, moveStructName) {
			continue
		}
		if objectId, ok := m["objectId"].(string); ok {
			return objectId, nil
		}
	}

	return "", errors.New("object ID not found")
}
