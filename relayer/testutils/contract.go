package testutils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"
)

type ObjectChange struct {
	Type            string   `json:"type"` // "published", "created", etc.
	Sender          string   `json:"sender,omitempty"`
	Owner           Owner    `json:"owner,omitempty"`
	ObjectType      string   `json:"objectType,omitempty"`
	ObjectID        string   `json:"objectId,omitempty"`
	Version         string   `json:"version,omitempty"`
	PreviousVersion string   `json:"previousVersion,omitempty"`
	Digest          string   `json:"digest,omitempty"`
	PackageID       string   `json:"packageId,omitempty"` // Only in type == "published"
	Modules         []string `json:"modules,omitempty"`   // Only in type == "published"
}

type Owner struct {
	AddressOwner *string      `json:"AddressOwner,omitempty"`
	Shared       *SharedOwner `json:"Shared,omitempty"`
	Immutable    *string      `json:"Immutable,omitempty"`
}

type SharedOwner struct {
	InitialSharedVersion int `json:"initial_shared_version"`
}

type TxnMetaWithObjectChanges struct {
	ObjectChanges []ObjectChange `json:"objectChanges"`
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

func findDigestIndex(input string) (int, error) {
	digestRegex := regexp.MustCompile(`"digest":\s*"[A-Za-z0-9]+"`)
	loc := digestRegex.FindStringIndex(input)
	if loc == nil {
		return -1, errors.New("digest not found")
	}
	return loc[0], nil
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

// LoadCompiledModules given a path to an already built contract, this method will
// find all the files ending with `.mv`
func LoadCompiledModules(packageName string, contractPath string) ([]string, error) {
	var modules []string

	dir := filepath.Join(contractPath, "/build/", packageName, "bytecode_modules/")

	// check each item in the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// find `.mv` files
		if !info.IsDir() && filepath.Ext(path) == ".mv" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			encoded := base64.StdEncoding.EncodeToString(data)
			modules = append(modules, encoded)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return modules, nil
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
//	packageName  - A string representing the contract name (package name in Move.toml).
//	contractPath - A string representing the filesystem path to the Move contract.
//	gasBudget    - A pointer to an int that specifies the gas budget for the publish transaction.
//	               If nil, a default value is used.
//
// Returns:
//
//	packageId    - The package ID extracted from the JSON output, typically for a published contract.
//	output       - The cleaned JSON output from the publish command.
//	error        - An error if the publish operation fails or if a valid package ID is not found.
func PublishContract(t *testing.T, packageName string, contractPath string, accountAddress string, gasBudget *int) (string, TxnMetaWithObjectChanges, error) {
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

	// This is a hack to skip the warnings from the CLI output by searching for "digest" with regex
	// and then extracting the JSON from there.
	idx, err := findDigestIndex(string(publishOutput))
	require.NoError(t, err)
	cleanedOutput := "{" + string(publishOutput)[idx:]

	lgr.Debug(cleanedOutput)

	// Unmarshal the JSON into a map.
	var parsedPublishTxn TxnMetaWithObjectChanges
	if err := json.Unmarshal([]byte(cleanedOutput), &parsedPublishTxn); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	changes := parsedPublishTxn.ObjectChanges

	var packageId string
	for _, change := range changes {
		if change.Type == "published" {
			packageId = change.PackageID
			break
		}
	}
	require.NotEmpty(t, packageId, "Package ID not found")

	return packageId, parsedPublishTxn, nil
}

func CallContractFromCLI(t *testing.T, packageId string, accountAddress string, module string, function string, gasBudget *int) TxnMetaWithObjectChanges {
	t.Helper()

	gasBudgetArg := "200000000"
	if gasBudget != nil {
		gasBudgetArg = string(rune(*gasBudget))
	}

	initializeCmd := exec.Command("sui", "client", "call",
		"--package", packageId,
		"--module", module,
		"--function", function,
		"--gas-budget", gasBudgetArg,
		"--json",
	)

	initializeOutput, err := initializeCmd.CombinedOutput()
	require.NoError(t, err, "Failed to initialize contract: %s", string(initializeOutput))

	// Unmarshal the JSON into a map.
	var parsedInitializeTxn TxnMetaWithObjectChanges
	if err := json.Unmarshal([]byte(initializeOutput), &parsedInitializeTxn); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	return parsedInitializeTxn
}

// QueryCreatedObjectID queries the created object ID for a given package ID, module, and struct name.
func QueryCreatedObjectID(objectChanges []ObjectChange, packageID, module, structName string) (string, error) {
	expectedType := fmt.Sprintf("%s::%s::%s", packageID, module, structName)

	for _, change := range objectChanges {
		if change.Type == "created" && change.ObjectType == expectedType {
			return change.ObjectID, nil
		}
	}

	return "", fmt.Errorf("object of type %s not found", expectedType)
}

// ExtractObjectId parses the JSON output from a Sui publish command and extracts an object identifier
// associated with a specific Move struct name. It expects the JSON to contain an "objectChanges" array,
// which may include various types of changes such as "published" and "created". When a "published"
// entry is present, the function extracts the "packageId", whereas for other types it might extract
// the "objectId" if that's what's required.
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
