package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/block-vision/sui-go-sdk/sui"
	"go.uber.org/zap"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/mockcontracts/account"
	"github.com/smartcontractkit/chainlink-sui/mockcontracts/deploy"
	"github.com/smartcontractkit/chainlink-sui/mockcontracts/events"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
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
		lggr.Infow("Starting setup command")
		runSetup(lggr)
	case "post-publish":
		lggr.Infow("Starting post-publish command")
		runPostPublish(lggr, os.Args[2:])
	case "emit-events":
		lggr.Infow("Starting emit-events command")
		packageIdFile := os.Args[2]
		packageId, err := os.ReadFile(packageIdFile)
		if err != nil {
			lggr.Errorw("Failed to read package ID file", "error", err)
			os.Exit(1)
		}

		lggr.Infow("Package ID", "packageId", string(packageId))

		runEmitEvents(lggr, string(packageId))
	case "emit-single-event":
		lggr.Infow("Starting emit-single-event command")
		var packageIdFile string
		var functionName string
		var contractName string

		// Create a new FlagSet for this subcommand
		flagSet := flag.NewFlagSet("emit-single-event", flag.ExitOnError)
		flagSet.StringVar(&packageIdFile, "package-id-file", "package_id.txt", "File containing the package ID")
		flagSet.StringVar(&functionName, "function-name", "", "Name of the function to emit")
		flagSet.StringVar(&contractName, "contract-name", "", "Name of the contract to emit the event from")

		// Parse the remaining arguments (excluding the command name)
		flagSet.Parse(os.Args[2:])

		packageId, err := os.ReadFile(packageIdFile)
		if err != nil {
			lggr.Errorw("Failed to read package ID file", "error", err)
			os.Exit(1)
		}
		runEmitSingleEvent(lggr, string(packageId), functionName, contractName)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`🚀 Sui Mock Contracts CLI

USAGE:
    go run mockcontracts/main.go <COMMAND> [OPTIONS]

COMMANDS:
    setup                    Setup local Sui chain and fund the active account
    post-publish <file>      Parse the output of the publish command
    emit-events [options]    Emit events (scaffolding)
	emit-single-event [options] Emit a single event

EXAMPLES:
    # Setup environment
    go run mockcontracts/main.go setup

    # Parse deployment output
    go run mockcontracts/main.go post-publish deployment_output.json

    # Emit events
    go run mockcontracts/main.go emit-events

	# Emit single event
    go run mockcontracts/main.go emit-single-event -package-id-file package_id.txt -event-name emit_static_config_set_event -contract-name offramp
`)
}

func runSetup(lggr logger.Logger) {
	lggr.Infow("Starting setup command")

	// stop sui node if it is running
	deploy.StopSuiNode(lggr, pidFile)

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

	signer, accountAddress, _, err := account.GenerateAccount(lggr)

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

	offrampEmitter := events.NewOffRampEmitter(lggr, packageId, signer, callOpts, client, accountAddress)

	tx, err := offrampEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit offramp events", "error", err)
		return
	}

	lggr.Infow("Emitting offramp events", "tx", tx)

	tokenAdminRegistryEmitter := events.NewTokenAdminRegistryEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = tokenAdminRegistryEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit token admin registry events", "error", err)
		return
	}

	lggr.Infow("Emitting token admin registry events", "tx", tx)

	feeQuoterEmitter := events.NewFeeQuoterEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = feeQuoterEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit fee quoter events", "error", err)
		return
	}
	lggr.Infow("Emitting fee quoter events", "tx", tx)

	rmnRemoteEmitter := events.NewRMNRemoteEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = rmnRemoteEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit RMN remote events", "error", err)
		return
	}
	lggr.Infow("Emitting RMN remote events", "tx", tx)

	routerEmitter := events.NewRouterEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = routerEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit router events", "error", err)
		return
	}
	lggr.Infow("Emitting router events", "tx", tx)

	onrampEmitter := events.NewOnRampEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = onrampEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit onramp events", "error", err)
		return
	}
	lggr.Infow("Emitting onramp events", "tx", tx)

	managedTokenPoolEmitter := events.NewManagedTokenPoolEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = managedTokenPoolEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit managed token pool events", "error", err)
		return
	}
	lggr.Infow("Emitting managed token pool events", "tx", tx)

	burnMintTokenPoolEmitter := events.NewBurnMintTokenPoolEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = burnMintTokenPoolEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit burn mint token pool events", "error", err)
		return
	}
	lggr.Infow("Emitting burn mint token pool events", "tx", tx)

	lockReleaseTokenPoolEmitter := events.NewLockReleaseTokenPoolEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
	tx, err = lockReleaseTokenPoolEmitter.BatchEmitEvents(ctx, DEFAULT_GAS_BUDGET, nil)
	if err != nil {
		lggr.Errorw("Failed to emit lock release token pool events", "error", err)
		return
	}
	lggr.Infow("Emitting lock release token pool events", "tx", tx)
	lggr.Infow("Event emission completed (scaffolding)")
}

func runEmitSingleEvent(lggr logger.Logger, packageId string, functionName string, contractName string) {
	lggr.Infow("Starting emit-single-event command", "packageId", packageId, "functionName", functionName)

	client := sui.NewSuiClient(testutils.LocalUrl)
	ctx := context.Background()

	signer, accountAddress, _, err := account.GenerateAccount(lggr)
	if err != nil {
		lggr.Errorw("Failed to generate account", "error", err)
		return
	}

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

	switch contractName {
	case "offramp":
		offrampEmitter := events.NewOffRampEmitter(lggr, packageId, signer, callOpts, client, accountAddress)
		tx, err := offrampEmitter.EmitEvent(ctx, DEFAULT_GAS_BUDGET, functionName)
		if err != nil {
			lggr.Errorw("Failed to emit offramp event", "error", err)
			return
		}
		lggr.Infow("Executed PTB", "tx", tx)
		lggr.Infow("Event emission completed (scaffolding)")
		return
	}
}
