package deploy

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

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

func StopSuiNode(lggr logger.Logger, pidFile string) {
	lggr.Infow("Stopping Sui node if it is running", "pidFile", pidFile)
	pid, err := os.ReadFile(pidFile)
	if err != nil {
		lggr.Warnw("Failed to read pid file", "error", err)
		return
	}

	defer func() {
		os.Remove(pidFile)
		lggr.Infow("Removed pid file", "pidFile", pidFile)
	}()

	cmd := exec.Command("kill", "-9", string(pid))
	output, err := cmd.CombinedOutput()
	if err != nil {
		lggr.Warnw("Failed to stop Sui node", "error", err, "output", string(output))
	}
	lggr.Infow("Stopped Sui node", "output", string(output))
}

func ParsePublishOutputFromFile(lggr logger.Logger, filename string) (string, TxnMetaWithObjectChanges, error) {
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

	lggr.Infow("Parsed publish transaction", "packageId", packageId)

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
