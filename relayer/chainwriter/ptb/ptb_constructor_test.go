//go:build integration

package ptb_test

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb"
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

func fakeExecutePTB(ctx context.Context, tx *transaction.Transaction) (string, error) {
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

	ctx := context.Background()

	// Test data
	packageID := "0x2742f32b2f375f9054a571f9e50ea6fedb91a181379db1869c27bcc6c8cfb955"
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
	writerConfig := config.ChainWriterConfig{} // Empty config, not needed for this test
	constructor := ptb.NewPTBConstructor(writerConfig, mockClient, log)

	builder := transaction.NewTransaction()

	// Test cases
	t.Run("Valid move call command", func(t *testing.T) {
		t.Parallel()

		cmd := config.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: &packageID,
			ModuleId:  &moduleID,
			Function:  &functionName,
			Params:    []codec.SuiFunctionParam{},
		}

		args := config.Arguments{Args: map[string]any{}}
		cachedArgs := map[string]transaction.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.NoError(t, err)
	})

	t.Run("Missing package ID", func(t *testing.T) {
		t.Parallel()

		cmd := config.ChainWriterPTBCommand{
			Type:     codec.SuiPTBCommandMoveCall,
			ModuleId: &moduleID,
			Function: &functionName,
			Params:   []codec.SuiFunctionParam{},
		}

		args := config.Arguments{Args: map[string]any{}}
		cachedArgs := map[string]transaction.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.Error(t, err)
		require.Contains(t, err.Error(), "missing required parameter 'PackageId'")
	})

	t.Run("Missing module ID", func(t *testing.T) {
		t.Parallel()

		cmd := config.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: &packageID,
			Function:  &functionName,
			Params:    []codec.SuiFunctionParam{},
		}

		args := config.Arguments{Args: map[string]any{}}
		cachedArgs := map[string]transaction.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.Error(t, err)
	})

	t.Run("Missing function name", func(t *testing.T) {
		t.Parallel()

		cmd := config.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: &packageID,
			ModuleId:  &moduleID,
			Params:    []codec.SuiFunctionParam{},
		}

		args := config.Arguments{Args: map[string]any{}}
		cachedArgs := map[string]transaction.Argument{}

		_, err := constructor.ProcessMoveCall(ctx, builder, cmd, &args, &cachedArgs)
		require.Error(t, err)
	})
}

// ------------------------------------------------
//
//	Tests prerequisite object filling
//
// ------------------------------------------------
//
//nolint:paralleltest
func TestPTBConstructor_PrereqObjectFill(t *testing.T) {
	ctx := context.Background()
	log, accountAddress, ptbClient, keystoreInstance, packageId, counterObjectId := setupTestEnvironment(t)
	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes := []byte(publicKey)

	txnSigner := signer.Signer{
		PriKey:  privateKey,
		PubKey:  publicKey,
		Address: accountAddress,
	}

	writerConfig := config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			"counter": {
				Name:     "counter",
				ModuleID: packageId,
				Functions: map[string]*config.ChainWriterFunction{
					"get_count_with_object_id_prereq": {
						Name:      "get_count_with_object_id_prereq",
						PublicKey: publicKeyBytes,
						PrerequisiteObjects: []config.PrerequisiteObject{
							{
								// we set the owner as the recently deployed counter contract
								OwnerId: &accountAddress,
								Name:    "admin_cap_id",
								Tag:     "counter::AdminCap",
								// we don't set the keys as we want to set the ID of the object in the PTB args
								SetKeys: false,
							},
						},
						PTBCommands: []config.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("increment_by_two_no_context"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "admin_cap_id",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "counter_id",
										Type:     "object_id",
										Required: true,
									},
								},
							},
						},
					},
					"get_count_with_object_keys_prereq": {
						Name:      "get_count_with_object_id_prereq",
						PublicKey: publicKeyBytes,
						PrerequisiteObjects: []config.PrerequisiteObject{
							{
								OwnerId: &accountAddress,
								// name doesn't matter here as we are setting the keys
								Name: "counter_id",
								Tag:  "counter::CounterPointer",
								// the keys of the returned object are set in the PTB args
								SetKeys: true,
							},
						},
						PTBCommands: []config.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("increment_by_two_no_context"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "admin_cap_id",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "counter_id",
										Type:     "object_id",
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	constructor := ptb.NewPTBConstructor(writerConfig, ptbClient, log)
	_ = transaction.NewTransaction()

	//nolint:paralleltest
	t.Run("Should fill a valid prerequisite object ID in CW config", func(t *testing.T) {
		// we only pass the counter ID as the other object ID (admin cap) is populated by the pre-requisites
		args := config.Arguments{Args: map[string]any{
			"counter_id": counterObjectId,
		}}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "get_count_with_object_id_prereq", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, &txnSigner, ptb, client.WaitForLocalExecution)
		prettyPrintDebug(log, ptbResult)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
	})

	//nolint:paralleltest
	t.Run("Should fill a valid prerequisite object keys in CW config", func(t *testing.T) {
		// pass no args as it should be populated by the pre-requisites
		args := config.Arguments{Args: map[string]any{}}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "get_count_with_object_keys_prereq", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, &txnSigner, ptb, client.WaitForLocalExecution)
		prettyPrintDebug(log, ptbResult)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
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

	txnSigner := signer.Signer{
		PriKey:  privateKey,
		PubKey:  publicKey,
		Address: accountAddress,
	}

	// Create PTB Constructor with config targeting the counter contract
	writerConfig := config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			"counter": {
				Name:     "counter",
				ModuleID: packageId,
				Functions: map[string]*config.ChainWriterFunction{
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
						PTBCommands: []config.ChainWriterPTBCommand{
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
							},
						},
					},
					"single_op_ptb": {
						Name:      "single_op_ptb",
						PublicKey: publicKeyBytes,
						PTBCommands: []config.ChainWriterPTBCommand{
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
							},
						},
					},
					"create_counter_manager": {
						Name:      "create_counter_manager",
						PublicKey: publicKeyBytes,
						PTBCommands: []config.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("create"),
								Params:    []codec.SuiFunctionParam{},
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
											ResultIndex:  testutils.Uint16Pointer(0),
										},
									},
								},
							},
						},
					},
					"manager_borrow_op_ptb": {
						Name:      "manager_borrow_op_ptb",
						PublicKey: publicKeyBytes,
						PTBCommands: []config.ChainWriterPTBCommand{
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
											ResultIndex:  testutils.Uint16Pointer(0),
										},
									},
								},
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
											ResultIndex:  testutils.Uint16Pointer(0),
										},
									},
								},
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
											ResultIndex:  testutils.Uint16Pointer(0),
										},
									},
									{
										Name:     "counter_borrow",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
											ResultIndex:  testutils.Uint16Pointer(1),
										},
									},
								},
							},
						},
					},
					"complex_operation": {
						Name:      "complex_operation",
						PublicKey: publicKeyBytes,
						PTBCommands: []config.ChainWriterPTBCommand{
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
							},
						},
					},
					"get_coin_value_ptb": {
						Name:      "get_coin_value_ptb",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []config.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  stringPointer("counter"),
								Function:  stringPointer("get_coin_value"),
								Params: []codec.SuiFunctionParam{
									{
										Name:      "coin",
										Type:      "object_id",
										Required:  true,
										IsGeneric: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	constructor := ptb.NewPTBConstructor(writerConfig, ptbClient, log)
	ctx := context.Background()

	// Test building and executing PTB commands
	//nolint:paralleltest
	t.Run("Single Operation PTB", func(t *testing.T) {
		args := config.Arguments{Args: map[string]any{
			"counter_id": counterObjectId,
		}}

		ptb, cError := constructor.BuildPTBCommands(ctx, "counter", "single_op_ptb", args, nil)
		require.NoError(t, cError)
		require.NotNil(t, ptb)

		// Execute the PTB command
		txHash, err := fakeExecutePTB(ctx, ptb)
		require.NoError(t, err)
		require.NotEmpty(t, txHash)
	})

	//nolint:paralleltest
	t.Run("Missing Module Error", func(t *testing.T) {
		args := config.Arguments{Args: map[string]any{
			"counter_id": counterObjectId,
		}}

		ptb, cError := constructor.BuildPTBCommands(ctx, "nonexistent_module", "get_count", args, nil)
		require.Error(t, cError)
		require.Nil(t, ptb)
	})

	//nolint:paralleltest
	t.Run("Missing Function Error", func(t *testing.T) {
		args := config.Arguments{Args: map[string]any{
			"counter_id": counterObjectId,
		}}

		ptb, cError := constructor.BuildPTBCommands(ctx, "counter", "nonexistent_function", args, nil)
		require.Error(t, cError)
		require.Nil(t, ptb)
	})

	//nolint:paralleltest
	t.Run("Missing Required Argument", func(t *testing.T) {
		args := config.Arguments{Args: map[string]any{}}

		ptb, cError := constructor.BuildPTBCommands(ctx, "incorrect_ptb", "get_count", args, nil)
		require.Error(t, cError)
		require.Nil(t, ptb)
	})

	//nolint:paralleltest
	t.Run("CounterManager Borrow Pattern", func(t *testing.T) {
		// Start by creating a Counter and its counter manager
		args := config.Arguments{Args: map[string]any{}}

		ptb, cError := constructor.BuildPTBCommands(ctx, "counter", "create_counter_manager", args, nil)
		require.NoError(t, cError)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, &txnSigner, ptb, client.WaitForLocalExecution)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
		prettyPrintDebug(log, ptbResult)

		// Borrow the Counter from the manager and pass it to increment then put it back
		var managerObjectId string
		// iterate through object changes
		for _, change := range ptbResult.ObjectChanges {
			if strings.Contains(change.ObjectType, "counter_manager") {
				managerObjectId = change.ObjectId
			}
		}

		args = config.Arguments{Args: map[string]any{
			"manager_object": managerObjectId,
		}}

		ptb, cError = constructor.BuildPTBCommands(ctx, "counter", "manager_borrow_op_ptb", args, nil)
		require.NoError(t, cError)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err = ptbClient.FinishPTBAndSend(ctx, &txnSigner, ptb, client.WaitForLocalExecution)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)

		// Expect 2 increment events
		incrementEventsCounter := 0
		for _, event := range ptbResult.Events {
			if strings.Contains(event.Type, "CounterIncremented") {
				incrementEventsCounter += 1
			}
		}
		require.Equal(t, 2, incrementEventsCounter)

		prettyPrintDebug(log, ptbResult)
	})

	//nolint:paralleltest
	t.Run("Complex Operation with Multiple Commands", func(t *testing.T) {
		args := config.Arguments{Args: map[string]any{
			"counter_id": counterObjectId,
		}}

		ptb, cError := constructor.BuildPTBCommands(ctx, "counter", "complex_operation", args, nil)
		require.NoError(t, cError)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, &txnSigner, ptb, client.WaitForLocalExecution)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		prettyPrintDebug(log, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
	})

	//nolint:paralleltest
	t.Run("PTB Constructor with Generic Type Tags - get_coin_value", func(t *testing.T) {
		// Fund the account
		cError := testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
		require.NoError(t, cError)

		// Get coins to use - need at least 2 coins (one for function arg, one for gas)
		coins, err := ptbClient.GetCoinsByAddress(ctx, txnSigner.Address)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(coins), 2, "Need at least 2 coins for this test")

		// Use the first coin as the test input
		testCoin := coins[1]
		log.Debugw("Using test coin with PTB constructor", "coinId", testCoin.CoinObjectId, "coinType", testCoin.CoinType)

		suiTypeTag := "0x2::sui::SUI"

		// Prepare arguments for the PTB constructor
		args := config.Arguments{
			Args: map[string]any{
				"coin": testCoin.CoinObjectId,
			},
			ArgTypes: map[string]string{
				"coin": suiTypeTag,
			},
		}

		// Use the constructor to build PTB commands for the generic function
		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "get_coin_value_ptb", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		log.Debugw("Executing generic function via PTB constructor", "ptb", ptb)

		// Execute the PTB command using the PTB client
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, &txnSigner, ptb, client.WaitForLocalExecution)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
		// Verify the function executed successfully
		log.Debugw("PTB Constructor generic function call successful", "coinValue", testCoin.Balance)
	})
}
