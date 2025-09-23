package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/mockcontracts/deploy"
	"github.com/smartcontractkit/chainlink-sui/mockcontracts/events"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"go.uber.org/zap"
)

const pidFile = "sui.pid"

var DEFAULT_GAS_BUDGET = uint64(1000000000)

func main() {
	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Create logger
	lggr, err := logger.NewWith(func(cfg *zap.Config) {
		cfg.Level.SetLevel(zap.InfoLevel)
	})
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "setup":
		runSetup(lggr)
	case "post-publish":
		runPostPublish(lggr, os.Args[2:])
	case "emit-events":
		packageId := os.Args[2]
		runEmitEvents(lggr, packageId)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
🚀 Sui Mock Contracts CLI

USAGE:
    go run mockcontracts/main.go <COMMAND> [OPTIONS]

COMMANDS:
    setup                    Setup local Sui chain and fund the active account
    post-publish <file>      Parse the output of the publish command
    emit-events [options]    Emit events (scaffolding)

EXAMPLES:
    # Setup environment
    go run mockcontracts/main.go setup

    # Parse deployment output
    go run mockcontracts/main.go post-publish deployment_output.json

    # Emit events
    go run mockcontracts/main.go emit-events
`)
}

func runSetup(lggr logger.Logger) {
	lggr.Infow("Starting setup command")

	if _, err := os.Stat(pidFile); err == nil {
		// File exists, do nothing
		lggr.Infow("sui.pid file exists, not starting Sui node")
	} else if os.IsNotExist(err) {
		// File does not exist, start Sui node
		cmd, err := testutils.StartSuiNode(testutils.CLI)
		if err != nil {
			lggr.Errorw("Failed to start Sui node", "error", err)
			return
		}
		pid := cmd.Process.Pid
		lggr.Infow("Sui node PID", "pid", pid)
		// Write PID to file
		if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
			lggr.Errorw("Failed to write sui.pid file", "error", err)
			return
		}
	} else {
		// Some other error accessing the file
		lggr.Errorw("Error checking sui.pid file", "error", err)
		return
	}

	activeAddress, err := deploy.GetActiveAddress(lggr)
	if err != nil {
		lggr.Errorw("Failed to get active address", "error", err)
		return
	}
	lggr.Infow("Active address", "activeAddress", activeAddress)

	for range 3 {
		fundErr := testutils.FundWithFaucet(lggr, testutils.SuiLocalnet, activeAddress)
		if fundErr != nil {
			lggr.Errorw("Failed to fund account", "error", fundErr)
			return
		}
	}

	lggr.Infow("Setup completed successfully", "activeAddress", activeAddress)
}

func runPostPublish(lggr logger.Logger, args []string) {
	lggr.Infow("Starting post-publish command")

	if len(args) == 0 {
		lggr.Errorw("Missing output file argument")
		fmt.Println("Usage: go run mockcontracts/main.go post-publish <output-file>")
		os.Exit(1)
	}

	outputFile := args[0]
	lggr.Infow("Parsing deployment output", "file", outputFile)

	// Read the deployment output file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		lggr.Errorw("Failed to read output file", "error", err, "file", outputFile)
		return
	}

	// Parse JSON
	var deploymentOutput map[string]interface{}
	if err := json.Unmarshal(data, &deploymentOutput); err != nil {
		lggr.Errorw("Failed to parse JSON", "error", err)
		return
	}

	// Extract package ID from objectChanges
	if objectChanges, ok := deploymentOutput["objectChanges"].([]interface{}); ok {
		for _, change := range objectChanges {
			if changeMap, ok := change.(map[string]interface{}); ok {
				if changeType, ok := changeMap["type"].(string); ok && changeType == "published" {
					if packageId, ok := changeMap["packageId"].(string); ok {
						lggr.Infow("Found published package", "packageId", packageId)

						// Write package ID to a file for easy access
						if err := os.WriteFile("package_id.txt", []byte(packageId), 0644); err != nil {
							lggr.Errorw("Failed to write package ID file", "error", err)
						} else {
							lggr.Infow("Package ID saved to package_id.txt")
						}
						return
					}
				}
			}
		}
	}

	lggr.Errorw("No published package found in output")
}

func runEmitEvents(lggr logger.Logger, packageId string) {
	lggr.Infow("Starting emit-events command")

	client := sui.NewSuiClient(testutils.LocalUrl)
	ctx := context.Background()

	signer, accountAddress, _, err := events.GenerateAccount(lggr)

	fundWithFaucetErr := testutils.FundWithFaucet(lggr, testutils.SuiLocalnet, accountAddress)
	if fundWithFaucetErr != nil {
		lggr.Errorw("Failed to fund account", "error", fundWithFaucetErr)
		return
	}

	callOpts := bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
		GasBudget:        &DEFAULT_GAS_BUDGET,
	}

	ptb := transaction.NewTransaction()

	offrampBoundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"offramp",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create offramp bound contract", "error", err)
		return
	}

	call, err := offrampBoundContract.EncodeCallArgsWithGenerics(
		"emit_static_config_set_event",
		[]string{},
		[]string{},
		[]string{"u64"},
		[]any{1},
		nil,
	)
	if err != nil {
		lggr.Errorw("Failed to encode call", "error", err)
		return
	}

	_, err = offrampBoundContract.AppendPTB(ctx, &callOpts, ptb, call)
	if err != nil {
		lggr.Errorw("Failed to append PTB", "error", err)
		return
	}

	tx, err := bind.ExecutePTB(ctx, &callOpts, client, ptb)
	if err != nil {
		lggr.Errorw("Failed to execute transaction", "error", err)
		return
	}

	lggr.Infow("Event emission completed (scaffolding)", "tx", tx)

}
