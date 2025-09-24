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
		runSetupCommand(lggr, os.Args[2:])
	case "post-publish":
		runPostPublishCommand(lggr, os.Args[2:])
	case "emit-events":
		runEmitEventsCommand(lggr, os.Args[2:])
	case "emit-single-event":
		runEmitSingleEventCommand(lggr, os.Args[2:])
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
    setup                     Setup local Sui chain and fund the active account
    post-publish [options]    Parse the output of the publish command
    emit-events [options]     Emit events from all contracts
    emit-single-event [opts]  Emit a single event from a specific contract

OPTIONS:
    Use -h or --help after any command to see command-specific options.

EXAMPLES:
    # Setup environment
    go run mockcontracts/main.go setup

    # Parse deployment output (flag-based)
    go run mockcontracts/main.go post-publish -file deployment_output.json
    # Or backwards-compatible positional argument
    go run mockcontracts/main.go post-publish deployment_output.json

    # Emit events from all contracts
    go run mockcontracts/main.go emit-events -package-id-file package_id.txt

    # Emit single event
    go run mockcontracts/main.go emit-single-event -package-id-file package_id.txt -function-name emit_static_config_set_event -contract-name offramp

    # Get help for specific commands
    go run mockcontracts/main.go setup -h
    go run mockcontracts/main.go emit-single-event -h`)
}

func runSetupCommand(lggr logger.Logger, args []string) {
	// Parse flags for setup command
	flagSet := flag.NewFlagSet("setup", flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s setup [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Setup local Sui chain and fund the active account\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flagSet.PrintDefaults()
	}

	if err := flagSet.Parse(args); err != nil {
		lggr.Errorw("Failed to parse setup flags", "error", err)
		os.Exit(1)
	}

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

func runPostPublishCommand(lggr logger.Logger, args []string) {
	// Parse flags for post-publish command
	var outputFile string
	flagSet := flag.NewFlagSet("post-publish", flag.ExitOnError)
	flagSet.StringVar(&outputFile, "file", "", "Deployment output file (required)")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s post-publish [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parse the output of the publish command\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flagSet.PrintDefaults()
	}

	if err := flagSet.Parse(args); err != nil {
		lggr.Errorw("Failed to parse post-publish flags", "error", err)
		os.Exit(1)
	}

	// Support positional argument for backwards compatibility
	if outputFile == "" && flagSet.NArg() > 0 {
		outputFile = flagSet.Arg(0)
	}

	if outputFile == "" {
		lggr.Errorw("Missing output file argument")
		flagSet.Usage()
		os.Exit(1)
	}

	lggr.Infow("Starting post-publish command")
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

func runEmitEventsCommand(lggr logger.Logger, args []string) {
	// Parse flags for emit-events command
	var packageIdFile string
	flagSet := flag.NewFlagSet("emit-events", flag.ExitOnError)
	flagSet.StringVar(&packageIdFile, "package-id-file", "package_id.txt", "File containing the package ID")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s emit-events [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Emit events from all contracts\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flagSet.PrintDefaults()
	}

	if err := flagSet.Parse(args); err != nil {
		lggr.Errorw("Failed to parse emit-events flags", "error", err)
		os.Exit(1)
	}

	packageId, err := os.ReadFile(packageIdFile)
	if err != nil {
		lggr.Errorw("Failed to read package ID file", "error", err)
		os.Exit(1)
	}

	lggr.Infow("Starting emit-events command")
	lggr.Infow("Package ID", "packageId", string(packageId))

	runEmitEvents(lggr, string(packageId))
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

func runEmitSingleEventCommand(lggr logger.Logger, args []string) {
	// Parse flags for emit-single-event command
	var packageIdFile string
	var functionName string
	var contractName string
	flagSet := flag.NewFlagSet("emit-single-event", flag.ExitOnError)
	flagSet.StringVar(&packageIdFile, "package-id-file", "package_id.txt", "File containing the package ID")
	flagSet.StringVar(&functionName, "function-name", "", "Name of the function to emit (required)")
	flagSet.StringVar(&contractName, "contract-name", "", "Name of the contract to emit the event from (required)")
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s emit-single-event [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Emit a single event from a specific contract\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flagSet.PrintDefaults()
	}

	if err := flagSet.Parse(args); err != nil {
		lggr.Errorw("Failed to parse emit-single-event flags", "error", err)
		os.Exit(1)
	}

	if functionName == "" || contractName == "" {
		lggr.Errorw("Missing required arguments")
		flagSet.Usage()
		os.Exit(1)
	}

	packageId, err := os.ReadFile(packageIdFile)
	if err != nil {
		lggr.Errorw("Failed to read package ID file", "error", err)
		os.Exit(1)
	}

	runEmitSingleEvent(lggr, string(packageId), functionName, contractName)
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
	default:
		lggr.Errorw("Unknown contract name", "contractName", contractName)
		fmt.Printf("Unknown contract: %s\n", contractName)
		fmt.Println("Available contracts: offramp")
		return
	}
}
