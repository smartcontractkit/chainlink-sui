//nolint:all
package testutils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"testing"

	"github.com/pelletier/go-toml/v2"

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
		"--dev",
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

	gasBudgetArg := "200000000"
	if gasBudget != nil {
		gasBudgetArg = strconv.Itoa(*gasBudget)
	}

	publishCmd := exec.Command("sui", "client", "publish",
		"--gas-budget", gasBudgetArg,
		"--json",
		"--silence-warnings",
		"--dev",
		"--with-unpublished-dependencies",
		contractPath,
	)

	publishOutput, err := publishCmd.CombinedOutput()
	require.NoError(t, err, "Failed to publish contract: %s", string(publishOutput))

	// This is a hack to skip the warnings from the CLI output by searching for "digest" with regex
	// and then extracting the JSON from there.
	idx, err := findDigestIndex(string(publishOutput))
	require.NoError(t, err)
	cleanedOutput := "{" + string(publishOutput)[idx:]

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

// PatchContractDevAddressTOML edits one entry under [dev-addresses].
// contractPath : folder that contains Move.toml
// name         : key to patch (e.g. "mcms")
// address      : new hex value (e.g. "0x0000")
func PatchContractDevAddressTOML(t *testing.T, contractPath, name, address string) {
	t.Helper()

	moveToml := filepath.Join(contractPath, "Move.toml")
	raw, err := os.ReadFile(moveToml)
	require.NoError(t, err, "read Move.toml")

	// Decode into a generic map[string]any
	var doc map[string]any
	err = toml.Unmarshal(raw, &doc)
	require.NoError(t, err, "parse TOML")

	// Ensure [dev-addresses] table exists
	devAddrs, ok := doc["dev-addresses"].(map[string]any)
	if !ok {
		devAddrs = make(map[string]any)
		doc["dev-addresses"] = devAddrs
	}

	// Set / overwrite the single entry
	devAddrs[name] = address

	// Re-encode with default indentation
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	err = enc.Encode(doc)
	require.NoError(t, err, "encode TOML")

	err = os.WriteFile(moveToml, buf.Bytes(), 0o644)
	require.NoError(t, err, "write Move.toml")
}
