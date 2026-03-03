//nolint:all
package testutils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type ObjectChange struct {
	Type            string   `json:"type"` // "published", "created", etc.
	Sender          string   `json:"sender,omitempty"`
	Owner           Owner    `json:"owner,omitempty"`
	ObjectType      string   `json:"objectType,omitempty"`
	ObjectTypeSnake string   `json:"object_type,omitempty"`
	ObjectID        string   `json:"objectId,omitempty"`
	ObjectIDSnake   string   `json:"object_id,omitempty"`
	Version         string   `json:"version,omitempty"`
	PreviousVersion string   `json:"previousVersion,omitempty"`
	Digest          string   `json:"digest,omitempty"`
	PackageID       string   `json:"packageId,omitempty"` // Only in type == "published"
	OutputState     string   `json:"outputState,omitempty"`
	OutputStateAlt  string   `json:"output_state,omitempty"`
	IDOperation     string   `json:"idOperation,omitempty"`
	IDOperationAlt  string   `json:"id_operation,omitempty"`
	Modules         []string `json:"modules,omitempty"` // Only in type == "published"
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
	ObjectChanges       []ObjectChange `json:"objectChanges"`
	ObjectChangesSnake  []ObjectChange `json:"object_changes"`
	ChangedObjects      []ObjectChange `json:"changed_objects"`
	ChangedObjectsCamel []ObjectChange `json:"changedObjects"`
}

func normalizeObjectChanges(changes []ObjectChange) []ObjectChange {
	for i := range changes {
		if changes[i].ObjectID == "" {
			changes[i].ObjectID = changes[i].ObjectIDSnake
		}
		if changes[i].ObjectType == "" {
			changes[i].ObjectType = changes[i].ObjectTypeSnake
		}
		if changes[i].OutputState == "" {
			changes[i].OutputState = changes[i].OutputStateAlt
		}
		if changes[i].IDOperation == "" {
			changes[i].IDOperation = changes[i].IDOperationAlt
		}

		// Newer Sui JSON uses object-level metadata rather than "type":"published/created".
		if changes[i].Type == "" {
			if changes[i].ObjectType == "package" || changes[i].OutputState == "OUTPUT_OBJECT_STATE_PACKAGE_WRITE" {
				changes[i].Type = "published"
				if changes[i].PackageID == "" {
					changes[i].PackageID = changes[i].ObjectID
				}
			} else if strings.EqualFold(changes[i].IDOperation, "CREATED") {
				changes[i].Type = "created"
			}
		}
	}
	return changes
}

func BuildSetup(t *testing.T, packagePath string) string {
	t.Helper()
	lgr := logger.Test(t)

	// Get the file path of the current source file
	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get current file path")
	// Get the directory containing the current file (which should be the testutils package)
	currentDir := filepath.Dir(currentFile)

	// Navigate to the project root (assuming we're in relayer/testutils)
	projectRoot := filepath.Dir(filepath.Dir(currentDir))
	contractPath := filepath.Join(projectRoot, packagePath)

	lgr.Debugw("Building contract setup", "path", contractPath)

	return contractPath
}

func extractJSONOutput(output string) (string, error) {
	start := strings.Index(output, "{")
	end := strings.LastIndex(output, "}")
	if start == -1 || end == -1 || end < start {
		return "", fmt.Errorf("json output not found")
	}
	return output[start : end+1], nil
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

	lgr.Infow("Publishing contract", "name", packageName, "path", contractPath)

	gasBudgetArg := "800000000"
	if gasBudget != nil {
		gasBudgetArg = strconv.Itoa(*gasBudget)
	}

	publishCmd := exec.Command("sui", "client", "publish",
		"--gas-budget", gasBudgetArg,
		"--json",
		"--silence-warnings",
		"--with-unpublished-dependencies",
		"-e", "local", // Explicitly use local environment from Move.toml
		contractPath,
	)

	publishOutput, err := publishCmd.CombinedOutput()
	require.NoError(t, err, "Failed to publish contract: %s", string(publishOutput))

	cleanedOutput, err := extractJSONOutput(string(publishOutput))
	require.NoError(t, err)

	// Unmarshal the JSON into a map.
	var parsedPublishTxn TxnMetaWithObjectChanges
	err = json.Unmarshal([]byte(cleanedOutput), &parsedPublishTxn)
	require.NoError(t, err, "Failed to parse publish output: %s", cleanedOutput)

	if len(parsedPublishTxn.ObjectChanges) == 0 && len(parsedPublishTxn.ObjectChangesSnake) > 0 {
		parsedPublishTxn.ObjectChanges = parsedPublishTxn.ObjectChangesSnake
	}
	if len(parsedPublishTxn.ObjectChanges) == 0 && len(parsedPublishTxn.ChangedObjects) > 0 {
		parsedPublishTxn.ObjectChanges = parsedPublishTxn.ChangedObjects
	}
	if len(parsedPublishTxn.ObjectChanges) == 0 && len(parsedPublishTxn.ChangedObjectsCamel) > 0 {
		parsedPublishTxn.ObjectChanges = parsedPublishTxn.ChangedObjectsCamel
	}
	parsedPublishTxn.ObjectChanges = normalizeObjectChanges(parsedPublishTxn.ObjectChanges)

	changes := parsedPublishTxn.ObjectChanges

	var packageId string
	for _, change := range changes {
		if change.Type == "published" && change.PackageID != "" {
			packageId = change.PackageID
			break
		}

		// Newer Sui output can represent package publish as objectType=package.
		if (change.ObjectType == "package" || change.OutputState == "OUTPUT_OBJECT_STATE_PACKAGE_WRITE") && change.ObjectID != "" {
			packageId = change.ObjectID
			break
		}
	}
	require.NotEmpty(t, packageId, "Package ID not found")

	return packageId, parsedPublishTxn, nil
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

// patchContractTOMLSection edits one entry under the specified TOML section.
// contractPath : folder that contains Move.toml
// section      : TOML section name (e.g., "addresses")
// name         : key to patch (e.g. "mcms", "test_secondary")
// address      : new hex value (e.g. "0x0000", "0x123...")
func patchContractTOMLSection(t *testing.T, contractPath, addresses, name, address string) {
	t.Helper()

	// Only resolve relative paths to absolute paths
	if !filepath.IsAbs(contractPath) {
		// Get the file path of the current source file
		_, currentFile, _, ok := runtime.Caller(0)
		require.True(t, ok, "Failed to get current file path")
		// Get the directory containing the current file (which should be the testutils package)
		currentDir := filepath.Dir(currentFile)

		// Navigate to the project root (assuming we're in relayer/testutils)
		projectRoot := filepath.Dir(filepath.Dir(currentDir))
		contractPath = filepath.Join(projectRoot, contractPath)
	}

	moveToml := filepath.Join(contractPath, "Move.toml")
	raw, err := os.ReadFile(moveToml)
	require.NoError(t, err, "read Move.toml")

	// Decode into a generic map[string]any
	var doc map[string]any
	err = toml.Unmarshal(raw, &doc)
	require.NoError(t, err, "parse TOML")

	if addresses == "addresses" {
		// Ensure the section [addresses] table exists
		addrs, ok := doc[addresses].(map[string]any)
		if !ok {
			addrs = make(map[string]any)
			doc[addresses] = addrs
		}

		// Set / overwrite the single entry
		addrs[name] = address

		// Re-encode with default indentation
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		enc.SetIndentTables(true)
		err = enc.Encode(doc)
		require.NoError(t, err, "encode TOML")

		err = os.WriteFile(moveToml, buf.Bytes(), 0o644)
		require.NoError(t, err, "write Move.toml")
	} else if addresses == "environments" {
		// Add entry under [environments]. If the section exists, only add/replace
		// the entry; if it doesn't, append a new section with blank lines above/below.
		envs, ok := doc[addresses].(map[string]any)
		if ok {
			envs[name] = address

			var buf bytes.Buffer
			enc := toml.NewEncoder(&buf)
			enc.SetIndentTables(true)
			err = enc.Encode(doc)
			require.NoError(t, err, "encode TOML")

			err = os.WriteFile(moveToml, buf.Bytes(), 0o644)
			require.NoError(t, err, "write Move.toml")
		} else {
			// Append with a leading and trailing empty line.
			if len(raw) == 0 || raw[len(raw)-1] != '\n' {
				raw = append(raw, '\n')
			}
			appendSection := fmt.Sprintf("\n[environments]\n%s = \"%s\"\n\n", name, address)
			err = os.WriteFile(moveToml, append(raw, []byte(appendSection)...), 0o644)
			require.NoError(t, err, "write Move.toml")
		}
	}

	// Log resulting TOML contents for debugging.
	finalToml, err := os.ReadFile(moveToml)
	require.NoError(t, err, "read patched Move.toml")
	t.Logf("Patched Move.toml (%s):\n%s\n", moveToml, string(finalToml))
}

// PatchContractAddressTOML edits one entry under [addresses].
func PatchContractAddressTOML(t *testing.T, contractPath, name, address string) {
	patchContractTOMLSection(t, contractPath, "addresses", name, address)
}

func PatchEnvironmentTOML(contractPath, environment, chainID string) {
	patchContractTOMLSectionNoTest(contractPath, "environments", environment, chainID)
}

func patchContractTOMLSectionNoTest(contractPath, addresses, name, address string) {
	// Only resolve relative paths to absolute paths
	if !filepath.IsAbs(contractPath) {
		// Get the file path of the current source file
		_, currentFile, _, _ := runtime.Caller(0)
		// require.True(t, ok, "Failed to get current file path")
		// Get the directory containing the current file (which should be the testutils package)
		currentDir := filepath.Dir(currentFile)

		// Navigate to the project root (assuming we're in relayer/testutils)
		projectRoot := filepath.Dir(filepath.Dir(currentDir))
		contractPath = filepath.Join(projectRoot, contractPath)
	}

	moveToml := filepath.Join(contractPath, "Move.toml")
	raw, _ := os.ReadFile(moveToml)
	// require.NoError(t, err, "read Move.toml")

	// Decode into a generic map[string]any
	var doc map[string]any
	_ = toml.Unmarshal(raw, &doc)
	// require.NoError(t, err, "parse TOML")

	if addresses == "addresses" {
		// Ensure the section [addresses] table exists
		addrs, ok := doc[addresses].(map[string]any)
		if !ok {
			addrs = make(map[string]any)
			doc[addresses] = addrs
		}

		// Set / overwrite the single entry
		addrs[name] = address

		// Re-encode with default indentation
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		enc.SetIndentTables(true)
		_ = enc.Encode(doc)
		// require.NoError(t, err, "encode TOML")

		_ = os.WriteFile(moveToml, buf.Bytes(), 0o644)
		// require.NoError(t, err, "write Move.toml")
	} else if addresses == "environments" {
		// Add entry under [environments]. If the section exists, only add/replace
		// the entry; if it doesn't, append a new section with blank lines above/below.
		envs, ok := doc[addresses].(map[string]any)
		if ok {
			envs[name] = address

			var buf bytes.Buffer
			enc := toml.NewEncoder(&buf)
			enc.SetIndentTables(true)
			_ = enc.Encode(doc)
			// require.NoError(t, err, "encode TOML")

			_ = os.WriteFile(moveToml, buf.Bytes(), 0o644)
			// require.NoError(t, err, "write Move.toml")
		} else {
			// Append with a leading and trailing empty line.
			if len(raw) == 0 || raw[len(raw)-1] != '\n' {
				raw = append(raw, '\n')
			}
			appendSection := fmt.Sprintf("\n[environments]\n%s = \"%s\"\n\n", name, address)
			_ = os.WriteFile(moveToml, append(raw, []byte(appendSection)...), 0o644)
			// require.NoError(t, err, "write Move.toml")
		}
	}

	// Log resulting TOML contents for debugging.
	finalToml, _ := os.ReadFile(moveToml)
	// require.NoError(t, err, "read patched Move.toml")
	log.Printf("Patched Move.toml (%s):\n%s\n", moveToml, string(finalToml))
}

// CleanupTestContracts removes the [published.local] entries from Published.toml files
// for all test contracts. This should be called at the start of tests AND registered
// with t.Cleanup to ensure a clean state for each test run
func CleanupTestContracts() {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return
	}
	currentDir := filepath.Dir(currentFile)
	projectRoot := filepath.Dir(filepath.Dir(currentDir))

	contractPaths := []string{
		filepath.Join(projectRoot, "contracts", "test"),
		filepath.Join(projectRoot, "contracts", "test_secondary"),
	}

	for _, path := range contractPaths {
		removeLocalPublishedEntry(path)
	}
}

func removeLocalPublishedEntry(contractPath string) {
	publishedToml := filepath.Join(contractPath, "Published.toml")
	content, err := os.ReadFile(publishedToml)
	if err != nil {
		return
	}

	// Parse TOML
	var doc map[string]interface{}
	if err := toml.Unmarshal(content, &doc); err != nil {
		return
	}

	// Check if there's a published section
	published, ok := doc["published"].(map[string]interface{})
	if !ok {
		return
	}

	// Remove the local entry if it exists
	if _, hasLocal := published["local"]; hasLocal {
		delete(published, "local")

		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		enc.SetIndentTables(true)
		if err := enc.Encode(doc); err != nil {
			return
		}
		_ = os.WriteFile(publishedToml, buf.Bytes(), 0o644)
	}
}
