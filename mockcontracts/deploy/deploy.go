package deploy

import (
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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func main() {
	fmt.Println("Hello, World!")
}

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
func PublishContract(lggr logger.Logger, packageName string, contractPath string, gasBudget *int) (string, TxnMetaWithObjectChanges, error) {

	lggr.Infow("Publishing contract", "name", packageName, "path", contractPath)

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
	if err != nil {
		lggr.Errorw("Failed to publish contract", "error", err, "output", string(publishOutput))
		return "", TxnMetaWithObjectChanges{}, err
	}

	// This is a hack to skip the warnings from the CLI output by searching for "digest" with regex
	// and then extracting the JSON from there.
	idx, err := findDigestIndex(string(publishOutput))
	if err != nil {
		return "", TxnMetaWithObjectChanges{}, err
	}
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

	return packageId, parsedPublishTxn, nil
}

func ParsePublishOutputFromFile(filename string) (string, TxnMetaWithObjectChanges, error) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return "", TxnMetaWithObjectChanges{}, err
	}

	// This is a hack to skip the warnings from the CLI output by searching for "digest" with regex
	// and then extracting the JSON from there.
	idx, err := findDigestIndex(string(jsonData))
	if err != nil {
		return "", TxnMetaWithObjectChanges{}, err
	}
	cleanedOutput := "{" + string(jsonData)[idx:]

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

	return packageId, parsedPublishTxn, nil
}

func findDigestIndex(input string) (int, error) {
	digestRegex := regexp.MustCompile(`"digest":\s*"[A-Za-z0-9]+"`)
	loc := digestRegex.FindStringIndex(input)
	if loc == nil {
		return -1, errors.New("digest not found")
	}

	return loc[0], nil
}

func BuildSetup(lggr logger.Logger, packagePath string) (string, error) {

	// Get the file path of the current source file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get current file path")
	}
	// Get the directory containing the current file (which should be the testutils package)
	currentDir := filepath.Dir(currentFile)

	// Navigate to the project root (assuming we're in relayer/testutils)
	projectRoot := filepath.Dir(filepath.Dir(currentDir))
	contractPath := filepath.Join(projectRoot, packagePath)

	lggr.Debugw("Building contract setup", "path", contractPath)

	return contractPath, nil
}

func BuildContract(lggr logger.Logger, contractPath string) error {

	lggr.Infow("Building contract", "path", contractPath)

	cmd := exec.Command("sui", "move", "build", "--path",
		contractPath,
		"--dev",
	)
	lggr.Debugw("Executing build command", "command", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build contract: %s", string(output))
	}

	return nil
}

// GetActiveAddress returns the active Sui client address
func GetActiveAddress(lggr logger.Logger) (string, error) {
	lggr.Debugw("Getting active Sui address")

	cmd := exec.Command("sui", "client", "active-address")
	output, err := cmd.CombinedOutput()
	if err != nil {
		lggr.Errorw("Failed to get active address", "error", err, "output", string(output))
		return "", fmt.Errorf("failed to get active address: %w", err)
	}

	// Trim whitespace and return the address
	address := string(output)
	address = regexp.MustCompile(`\s+`).ReplaceAllString(address, "")

	lggr.Infow("Retrieved active address", "address", address)
	return address, nil
}
