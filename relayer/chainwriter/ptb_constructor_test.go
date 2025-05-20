//go:build integration

package chainwriter_test

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// ------------------------------------------
//
//	Setup and Helpers
//
// ------------------------------------------
// setupTestEnvironment sets up the test environment with a local Sui node and deploys the counter contract
func setupTestEnvironment(t *testing.T) (
	log logger.Logger,
	accountAddress string,
	relayerClient *client.PTBClient,
	keystoreInstance keystore.SuiKeystore,
	packageId string,
	counterObjectId string,
) {
	t.Helper()

	log = logger.Test(t)
	accountAddress = testutils.GetAccountAndKeyFromSui(t, log)

	// Start local Sui node
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	// Ensure the process is killed when the test completes
	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	log.Debugw("Started Sui node")

	// Fund the account
	err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	// Set up keystore and signer
	keystoreInstance, err = keystore.NewSuiKeystore(log, "")
	require.NoError(t, err)

	relayerClient, err = client.NewPTBClient(log, testutils.LocalUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, err)

	// Build and publish contract
	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, publishOutput, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	// Get the counter object ID
	counterObjectId, err = testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	log.Debugw("Counter object created", "counterObjectId", counterObjectId)

	return log, accountAddress, relayerClient, keystoreInstance, packageId, counterObjectId
}

func stringPointer(s string) *string {
	return &s
}

func fakeExecutePTB(ctx context.Context, ptb *suiptb.ProgrammableTransactionBuilder) (string, error) {
	return "0x1234567890abcdef", nil
}

func prettyPrintDebug(log logger.Logger, data any) {
	resultJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Errorw("Failed to marshal data to JSON", "error", err)
	} else {
		log.Debugf("PTB Result:\n%s", string(resultJSON))
	}
}

// ------------------------------------------
//
// # Tests without actual contract interaction
//
// ------------------------------------------
func TestPTBConstructor_ProcessMoveCall(t *testing.T) {
	t.Parallel()
	t.Skip("Skipping test until PTB mock for execution is implemented")

	ctx := context.Background()

	// Test data
	packageID := "0x1234567890abcdef"
	moduleID := "test_module"
	functionName := "test_function"

	// Setup mock client
	mockClient := &testutils.FakeSuiPTBClient{
		Status: client.TransactionResult{
			Status: "failure",
			Error:  "ErrGasBudgetTooHigh",
		},
	}

	log := logger.Test(t)
	config := chainwriter.ChainWriterConfig{} // Empty config, not needed for this test
	constructor := chainwriter.NewPTBConstructor(config, mockClient, log)

	builder := suiptb.NewTransactionDataTransactionBuilder()

	// Test cases
	t.Run("Valid move call command", func(t *testing.T) {
		t.Parallel()

		cmd := chainwriter.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: &packageID,
			ModuleId:  &moduleID,
			Function:  &functionName,
			Params:    []codec.SuiFunctionParam{},
		}

		args := map[string]any{}
		cachedArgs := map[string]suiptb.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.NoError(t, err)
	})

	t.Run("Missing package ID", func(t *testing.T) {
		t.Parallel()

		cmd := chainwriter.ChainWriterPTBCommand{
			Type:     codec.SuiPTBCommandMoveCall,
			ModuleId: &moduleID,
			Function: &functionName,
			Params:   []codec.SuiFunctionParam{},
		}

		args := map[string]any{}
		cachedArgs := map[string]suiptb.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing required parameter 'PackageId'")
	})

	t.Run("Missing module ID", func(t *testing.T) {
		t.Parallel()

		cmd := chainwriter.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: &packageID,
			Function:  &functionName,
			Params:    []codec.SuiFunctionParam{},
		}

		args := map[string]any{}
		cachedArgs := map[string]suiptb.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.Error(t, err)
	})

	t.Run("Missing function name", func(t *testing.T) {
		t.Parallel()

		cmd := chainwriter.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: &packageID,
			ModuleId:  &moduleID,
			Params:    []codec.SuiFunctionParam{},
		}

		args := map[string]any{}
		cachedArgs := map[string]suiptb.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.Error(t, err)
	})

	t.Run("Invalid package ID format", func(t *testing.T) {
		t.Parallel()

		invalidPackageID := "invalid-hex"
		cmd := chainwriter.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: &invalidPackageID,
			ModuleId:  &moduleID,
			Function:  &functionName,
			Params:    []codec.SuiFunctionParam{},
		}

		args := map[string]any{}
		cachedArgs := map[string]suiptb.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.Error(t, err)
	})
}

// ------------------------------------------------
//
//	Tests with contract interaction
//
// ------------------------------------------------
//
//nolint:paralleltest
func TestPTBConstructor_IntegrationWithCounter(t *testing.T) {
	// Set up the test environment
	log, accountAddress, ptbClient, keystoreInstance, packageId, counterObjectId := setupTestEnvironment(t)

	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes := []byte(publicKey)

	// Create PTB Constructor with config targeting the counter contract
	config := chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			"counter": {
				Name:     "counter",
				ModuleID: packageId,
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"get_count": {
						Name:      "get_count",
						PublicKey: publicKeyBytes,
						Params: []codec.SuiFunctionParam{
							{
								Name:     "counter_id",
								Type:     "object_id",
								Required: true,
							},
						},
					},
					"increment_counter": {
						Name:      "increment",
						PublicKey: publicKeyBytes,
						Params: []codec.SuiFunctionParam{
							{
								Name:     "counter_id",
								Type:     "object_id",
								Required: true,
							},
						},
					},
					"incorrect_ptb": {
						Name:      "incorrect_ptb",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("get_count"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_id",
										Type:     "object_id",
										Required: true,
									},
								},
								Order: 1,
							},
						},
					},
					"single_op_ptb": {
						Name:      "single_op_ptb",
						PublicKey: publicKeyBytes,
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("get_count"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_id",
										Type:     "object_id",
										Required: true,
									},
								},
								Order: 1,
							},
						},
					},
					"create_counter_manager": {
						Name:      "create_counter_manager",
						PublicKey: publicKeyBytes,
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("create"),
								Params:    []codec.SuiFunctionParam{},
								Order:     1,
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter_manager"),
								Function:  stringPointer("create"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_id",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
											ResultIndex:  0,
										},
									},
								},
								Order: 2,
							},
						},
					},
					"manager_borrow_op_ptb": {
						Name:      "manager_borrow_op_ptb",
						PublicKey: publicKeyBytes,
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter_manager"),
								Function:  stringPointer("borrow_counter"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "manager_object",
										Type:     "object_id",
										Required: true,
									},
								},
								Order: 1,
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("increment_by_one_no_context"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_object",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
											ResultIndex:  0,
										},
									},
								},
								Order: 2,
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("increment_by_one_no_context"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_object",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
											ResultIndex:  0,
										},
									},
								},
								Order: 3,
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter_manager"),
								Function:  stringPointer("return_counter"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "manager_object",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "counter_object",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
											ResultIndex:  0,
										},
									},
									{
										Name:     "counter_borrow",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
											ResultIndex:  1,
										},
									},
								},
								Order: 4,
							},
						},
					},
					"complex_operation": {
						Name:      "complex_operation",
						PublicKey: publicKeyBytes,
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("increment"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_id",
										Type:     "object_id",
										Required: true,
									},
								},
								Order: 1,
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("increment_by"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_id",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:         "increment_by",
										Type:         "u64",
										Required:     true,
										DefaultValue: uint64(10),
									},
								},
								Order: 2,
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("get_count"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter_id",
										Type:     "object_id",
										Required: true,
									},
								},
								Order: 3,
							},
						},
					},
				},
			},
		},
	}

	constructor := chainwriter.NewPTBConstructor(config, ptbClient, log)
	ctx := context.Background()

	// Test building and executing PTB commands
	//nolint:paralleltest
	t.Run("Single Operation PTB", func(t *testing.T) {
		args := map[string]any{
			"counter_id": counterObjectId,
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "single_op_ptb", args)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		txHash, err := fakeExecutePTB(ctx, ptb)
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
	})

	//nolint:paralleltest
	t.Run("Missing Module Error", func(t *testing.T) {
		args := map[string]any{
			"counter_id": counterObjectId,
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "nonexistent_module", "get_count", args)
		require.Error(t, err)
		require.Nil(t, ptb)
	})

	//nolint:paralleltest
	t.Run("Missing Function Error", func(t *testing.T) {
		args := map[string]any{
			"counter_id": counterObjectId,
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "nonexistent_function", args)
		require.Error(t, err)
		require.Nil(t, ptb)
	})

	//nolint:paralleltest
	t.Run("Missing Required Argument", func(t *testing.T) {
		args := map[string]any{}

		ptb, err := constructor.BuildPTBCommands(ctx, "incorrect_ptb", "get_count", args)
		require.Error(t, err)
		require.Nil(t, ptb)
	})

	//nolint:paralleltest
	t.Run("CounterManager Borrow Pattern", func(t *testing.T) {
		// Start by creating a Counter and its counter manager
		args := map[string]any{}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "create_counter_manager", args)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, publicKeyBytes, ptb)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
		prettyPrintDebug(log, ptbResult)

		// Borrow the Counter from the manager and pass it to increment then put it back
		var managerObjectId string
		// iterate through object changes
		for _, change := range ptbResult.ObjectChanges {
			if change.Data.Created != nil && strings.Contains(change.Data.Created.ObjectType, "counter_manager") {
				managerObjectId = change.Data.Created.ObjectId.String()
			}
		}

		args = map[string]any{
			"manager_object": managerObjectId,
		}

		ptb, err = constructor.BuildPTBCommands(ctx, "counter", "manager_borrow_op_ptb", args)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err = ptbClient.FinishPTBAndSend(ctx, publicKeyBytes, ptb)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)

		// Expect 2 increment events
		incrementEventsCounter := 0
		for _, event := range ptbResult.Events {
			if strings.Contains(event.Type.String(), "CounterIncremented") {
				incrementEventsCounter += 1
			}
		}
		require.Equal(t, 2, incrementEventsCounter)

		prettyPrintDebug(log, ptbResult)
	})

	//nolint:paralleltest
	t.Run("Complex Operation with Multiple Commands", func(t *testing.T) {
		args := map[string]any{
			"counter_id": counterObjectId,
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "complex_operation", args)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, publicKeyBytes, ptb)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		prettyPrintDebug(log, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
	})
}
