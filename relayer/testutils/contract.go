//nolint:all
package testutils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

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

	// Collect contract dirs that need Published.toml cleanup: the package itself
	// and any local dependencies declared in Move.toml (e.g. test_secondary).
	dirsToClean := []string{contractPath}
	if moveToml, err := os.ReadFile(filepath.Join(contractPath, "Move.toml")); err == nil {
		for _, line := range strings.Split(string(moveToml), "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "local") {
				continue
			}
			for _, q := range []string{`"`, `'`} {
				prefix := "local = " + q
				if idx := strings.Index(line, prefix); idx != -1 {
					relPath := line[idx+len(prefix):]
					if end := strings.Index(relPath, q); end != -1 {
						p := relPath[:end]
						if strings.HasPrefix(p, ".") || strings.HasPrefix(p, "/") {
							dirsToClean = append(dirsToClean, filepath.Join(contractPath, p))
						}
					}
				}
			}
		}
	}

	// Remove Published.toml and ephemeral pubfiles for the MAIN package only.
	// Dependency dirs (e.g. test_secondary) keep their Published.toml so that
	// --with-unpublished-dependencies can detect packages that were already
	// published by a prior PublishContract call in the same test. Without this,
	// the dependency gets re-published as a separate on-chain package, causing
	// TypeMismatch errors when objects from the original publish are used.
	os.Remove(filepath.Join(contractPath, "Published.toml"))
	if pubGlob, err := filepath.Glob(filepath.Join(contractPath, "Pub.*.toml")); err == nil {
		for _, f := range pubGlob {
			os.Remove(f)
		}
	}

	// Snapshot Move.toml and Move.lock for each dir so we can restore them
	// after the test. This prevents patchMoveTomlEnvironment and the Sui CLI
	// from permanently dirtying the source tree.
	type fileSnapshot struct {
		path string
		data []byte
	}
	var snapshots []fileSnapshot
	for _, dir := range dirsToClean {
		for _, name := range []string{"Move.toml", "Move.lock"} {
			p := filepath.Join(dir, name)
			if data, err := os.ReadFile(p); err == nil {
				snapshots = append(snapshots, fileSnapshot{p, data})
			}
		}
	}

	// Register cleanup to restore source tree files after the test.
	t.Cleanup(func() {
		for _, dir := range dirsToClean {
			os.WriteFile(filepath.Join(dir, "Published.toml"), []byte{}, 0644) //nolint:errcheck
			pubGlob, _ := filepath.Glob(filepath.Join(dir, "Pub.*.toml"))
			for _, f := range pubGlob {
				os.Remove(f)
			}
		}
		for _, s := range snapshots {
			os.WriteFile(s.path, s.data, 0644) //nolint:errcheck
		}
	})

	// Create a CLI environment with a unique name derived from the current
	// node's chain ID, then switch to it. This avoids stale chain ID caching
	// that occurs when "sui client new-env" is called with an alias that
	// already exists from a previous node (with a different chain ID).
	chainID, err := GetChainIdentifier(LocalURL)
	require.NoError(t, err, "failed to get chain identifier before publish")

	envName := ensureCLIEnvForChainID(chainID)

	// Patch Move.toml [environments] with a matching entry so the CLI's
	// chain ID check passes.
	for _, dir := range dirsToClean {
		patchMoveTomlEnvironment(filepath.Join(dir, "Move.toml"), envName, chainID)
	}

	logFile, _ := os.OpenFile("/Users/felix/dev/chainlink/.cursor/debug-7c7360.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if logFile != nil {
		fmt.Fprintf(logFile, "{\"sessionId\":\"7c7360\",\"hypothesisId\":\"A\",\"location\":\"contract.go:PublishContract\",\"message\":\"pre-publish state\",\"data\":{\"chainID\":%q,\"envName\":%q,\"contractPath\":%q,\"dirsToClean\":%q},\"timestamp\":%d}\n", chainID, envName, contractPath, dirsToClean, time.Now().UnixMilli())
		logFile.Close()
	}

	publishCmd := exec.Command("sui", "client", "publish",
		"--gas-budget", gasBudgetArg,
		"--json",
		"--silence-warnings",
		"--with-unpublished-dependencies",
		contractPath,
	)

	publishOutput, err := publishCmd.CombinedOutput()

	logFile2, _ := os.OpenFile("/Users/felix/dev/chainlink/.cursor/debug-7c7360.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if logFile2 != nil {
		outSnippet := string(publishOutput)
		if len(outSnippet) > 500 {
			outSnippet = outSnippet[:500]
		}
		fmt.Fprintf(logFile2, "{\"sessionId\":\"7c7360\",\"hypothesisId\":\"A\",\"location\":\"contract.go:PublishContract:post\",\"message\":\"publish result\",\"data\":{\"err\":%q,\"outputSnippet\":%q},\"timestamp\":%d}\n", fmt.Sprint(err), outSnippet, time.Now().UnixMilli())
		logFile2.Close()
	}

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
}

// patchMoveTomlEnvironment does a targeted text replacement of a chain ID value
// in Move.toml's [environments] section. Unlike PatchEnvironmentTOML, this avoids
// full TOML parse/re-encode which can silently corrupt the file format.
func patchMoveTomlEnvironment(moveTomlPath, envName, newChainID string) {
	content, err := os.ReadFile(moveTomlPath)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	inEnvSection := false
	envSectionIdx := -1
	patched := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "[environments]" {
			inEnvSection = true
			envSectionIdx = i
			continue
		}
		if inEnvSection && strings.HasPrefix(trimmed, "[") {
			inEnvSection = false
			break
		}

		if !inEnvSection {
			continue
		}

		// Match envName = 'value' or envName = "value" with optional indentation
		for _, q := range []string{`'`, `"`} {
			prefix := envName + " = " + q
			altPrefix := envName + "= " + q
			altPrefix2 := envName + " =" + q
			if strings.HasPrefix(trimmed, prefix) || strings.HasPrefix(trimmed, altPrefix) || strings.HasPrefix(trimmed, altPrefix2) {
				leading := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
				lines[i] = fmt.Sprintf("%s%s = %s%s%s", leading, envName, q, newChainID, q)
				patched = true
				break
			}
		}
		if patched {
			break
		}
	}

	if !patched && envSectionIdx >= 0 {
		// Entry doesn't exist in [environments]; insert it after the section header.
		newLine := fmt.Sprintf("  %s = '%s'", envName, newChainID)
		lines = append(lines[:envSectionIdx+1], append([]string{newLine}, lines[envSectionIdx+1:]...)...)
		patched = true
	}

	if !patched {
		return
	}

	os.WriteFile(moveTomlPath, []byte(strings.Join(lines, "\n")), 0644) //nolint:errcheck
}

// ensureCLIEnvForChainID creates a CLI environment alias unique to the given
// chain ID, then switches to it. Using a chain-ID-derived name avoids stale
// chain ID caching: "sui client new-env" silently fails when the alias already
// exists, so reusing a fixed "local" alias across node restarts (each with a
// different chain ID from --force-regenesis) leaves the CLI with an outdated
// chain ID that causes "Move.toml expects local to have chain ID …" errors.
func ensureCLIEnvForChainID(chainID string) string {
	envName := "local_" + chainID

	createCmd := exec.Command("sui", "client", "new-env", "--rpc", LocalURL, "--alias", envName)
	createCmd.CombinedOutput() //nolint:errcheck

	switchCmd := exec.Command("sui", "client", "switch", "--env", envName)
	switchCmd.CombinedOutput() //nolint:errcheck

	logFile, _ := os.OpenFile("/Users/felix/dev/chainlink/.cursor/debug-7c7360.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if logFile != nil {
		fmt.Fprintf(logFile, "{\"sessionId\":\"7c7360\",\"hypothesisId\":\"B\",\"location\":\"contract.go:ensureCLIEnvForChainID\",\"message\":\"CLI env created\",\"data\":{\"envName\":%q,\"chainID\":%q},\"timestamp\":%d}\n", envName, chainID, time.Now().UnixMilli())
		logFile.Close()
	}

	return envName
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
