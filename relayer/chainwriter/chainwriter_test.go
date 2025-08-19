//go:build integration

package chainwriter_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/test-go/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type Counter struct {
	Value string `json:"value"`
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}

//nolint:paralleltest
func TestChainWriterSubmitTransaction(t *testing.T) {
	_logger := logger.Test(t)
	metadata := []testutils.Contracts{
		{
			Path:     "contracts/test/",
			Name:     "test",
			ModuleID: "0x1",
			Objects: []testutils.ContractObject{
				{
					ObjectID:    "0x1",
					PackageName: "counter",
					StructName:  "Counter",
				},
			},
		},
	}

	testState := testutils.BootstrapTestEnvironment(t, testutils.CLI, metadata)
	publicKeyBytes := testState.PublicKeyBytes

	countContract := testState.Contracts[0]
	packageId := countContract.ModuleID
	objectID := countContract.Objects[0].ObjectID
	// ChainWriter configuration
	chainWriterConfig := config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			"counter": {
				Name:     "counter",
				ModuleID: testState.Contracts[0].ModuleID,
				Functions: map[string]*config.ChainWriterFunction{
					"increment": {
						Name:      "increment",
						PublicKey: publicKeyBytes,
						Params: []codec.SuiFunctionParam{
							{
								Name:         "counter",
								Type:         "object_id",
								Required:     true,
								DefaultValue: nil,
							},
						},
					},
				},
			},
			config.PTBChainWriterModuleName: {
				Name:     config.PTBChainWriterModuleName,
				ModuleID: "0x2",
				Functions: map[string]*config.ChainWriterFunction{
					"ptb_call": {
						Name:      "ptb_call",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []config.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  strPtr("counter"),
								Function:  strPtr("increment"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter",
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
								ModuleId:  strPtr("counter"),
								Function:  strPtr("get_coin_value"),
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

	_logger.Infow("ChainWriterConfig", "config", chainWriterConfig)

	chainWriter, err := chainwriter.NewSuiChainWriter(_logger, testState.TxManager, chainWriterConfig, false)
	require.NoError(t, err)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	err = chainWriter.Start(ctx)
	require.NoError(t, err)
	err = testState.TxManager.Start(ctx)
	require.NoError(t, err)

	defer chainWriter.Close()
	defer testState.TxManager.Close()

	// Simple map style for builder pattern
	simpleArgs := map[string]any{
		"counter": objectID,
	}

	// Get coins to use - need at least 2 coins (one for function arg, one for gas)
	coins, err := testState.SuiGateway.GetCoinsByAddress(ctx, testState.AccountAddress)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins), 2, "Need at least 2 coins for this test")

	// Use the first coin as the test input
	testCoin := coins[1]

	// Common validation functions
	getCounterValue := func() (string, error) {
		objectDetails, callBackErr := testState.SuiGateway.ReadObjectId(ctx, objectID)
		if callBackErr != nil {
			return "", callBackErr
		}

		return objectDetails.Content.SuiMoveObject.Fields["value"].(string), nil
	}
	getCoinBalance := func() (string, error) {
		return testCoin.Balance, nil
	}

	getErrorValue := func() (string, error) {
		return "", nil
	}

	// Test scenarios
	testScenarios := []struct {
		name             string
		txID             string
		txMeta           *commonTypes.TxMeta
		sender           string
		contractName     string
		functionName     string
		args             config.Arguments
		expectError      error
		expectedResult   string
		status           commonTypes.TransactionStatus
		numberAttemps    int
		getExpectedValue func() (string, error) // Function to fetch the expected value for comparison
	}{
		{
			name:             "Test ChainWriter with valid parameters",
			txID:             "test-txID",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:           testState.AccountAddress,
			contractName:     "counter",
			functionName:     "increment",
			args:             config.Arguments{Args: map[string]any{"counter": objectID}},
			expectError:      nil,
			expectedResult:   "1",
			status:           commonTypes.Finalized,
			numberAttemps:    1,
			getExpectedValue: getCounterValue,
		},
		{
			name:             "Test ChainWriter with PTB using builder",
			txID:             "test-ptb-txID",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:           testState.AccountAddress,
			contractName:     config.PTBChainWriterModuleName,
			functionName:     "ptb_call",
			args:             config.Arguments{Args: simpleArgs},
			expectError:      nil,
			expectedResult:   "2",
			status:           commonTypes.Finalized,
			numberAttemps:    1,
			getExpectedValue: getCounterValue,
		},
		{
			name:             "Test ChainWriter with missing argument for PTB using builder",
			txID:             "test-ptb-txID-missing-arg",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:           testState.AccountAddress,
			contractName:     config.PTBChainWriterModuleName,
			functionName:     "ptb_call",
			args:             config.Arguments{Args: map[string]any{}},
			expectError:      errors.New("required parameter counter has no value"),
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:             "Test ChainWriter with PTB using simple map",
			txID:             "test-ptb-simple-map",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:           testState.AccountAddress,
			contractName:     config.PTBChainWriterModuleName,
			functionName:     "ptb_call",
			args:             config.Arguments{Args: simpleArgs},
			expectError:      nil,
			expectedResult:   "3",
			status:           commonTypes.Finalized,
			numberAttemps:    1,
			getExpectedValue: getCounterValue,
		},
		{
			name:             "Test ChainWriter with invalid function call",
			txID:             "test-txID-invalid-func",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:           testState.AccountAddress,
			contractName:     "counter",
			functionName:     "nonexistent_function",
			args:             config.Arguments{Args: map[string]any{"counter": objectID}},
			expectError:      commonTypes.ErrNotFound,
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:             "Test ChainWriter with invalid contract",
			txID:             "test-txID-invalid-contract",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:           testState.AccountAddress,
			contractName:     "nonexistent_contract",
			functionName:     "increment",
			args:             config.Arguments{Args: map[string]any{"counter": objectID}},
			expectError:      commonTypes.ErrNotFound,
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:         "Test ChainWriter with invalid arguments",
			txID:         "wrong-args-txID",
			txMeta:       &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:       testState.AccountAddress,
			contractName: "counter",
			functionName: "increment",
			args: config.Arguments{Args: map[string]any{
				"counter":     objectID,
				"invalid_arg": "invalid_value",
				"extra_arg":   "extra_value",
				"extra_arg2":  123,
			}},
			expectError:      errors.New("argument count mismatch"),
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:         "Test ChainWriter with the same transaction ID",
			txID:         "test-txID",
			txMeta:       &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:       testState.AccountAddress,
			contractName: "counter",
			functionName: "increment",
			args: config.Arguments{Args: map[string]any{
				"counter": objectID,
			}},
			expectError:      errors.New("transaction already exists"),
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:         "Test ChainWriter with generic function",
			txID:         "test-get-coin-value-txID",
			txMeta:       &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:       testState.AccountAddress,
			contractName: config.PTBChainWriterModuleName,
			functionName: "get_coin_value_ptb",
			args: config.Arguments{
				Args: map[string]any{
					"coin": testCoin.CoinObjectId,
				},
				ArgTypes: map[string]string{
					"coin": "0x2::sui::SUI",
				},
			},
			expectError:      nil,
			expectedResult:   testCoin.Balance,
			status:           commonTypes.Finalized,
			numberAttemps:    1,
			getExpectedValue: getCoinBalance,
		},
	}

	//nolint:paralleltest
	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Submit the transaction
			err = chainWriter.SubmitTransaction(ctx, scenario.contractName, scenario.functionName,
				scenario.args,
				scenario.txID, scenario.sender,
				scenario.txMeta, nil,
			)
			if scenario.expectError != nil {
				require.Error(t, err)
				require.Equal(t, scenario.expectError, err)
			} else {
				require.NoError(t, err)

				require.Eventually(t, func() bool {
					status, statusErr := chainWriter.GetTransactionStatus(ctx, scenario.txID)
					if statusErr != nil {
						return false
					}

					return status == scenario.status
				}, 10*time.Second, 1*time.Second, "Transaction final state not reached")

				actualValue, err := scenario.getExpectedValue()
				require.NoError(t, err)

				assert.Equal(t, scenario.expectedResult, actualValue, "Expected value does not match")

				tx, err := testState.TxStore.GetTransaction(scenario.txID)
				require.NoError(t, err, "Failed to get transaction from repository")
				assert.Equal(t, scenario.numberAttemps, tx.Attempt, "Transaction attempts do not match")
			}
		})
	}
}
