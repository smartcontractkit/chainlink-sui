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
	ctx := context.Background()
	gasLimit := int64(10000000)
	_logger := logger.Test(t)
	suiClient, txManager, txStore, accountAddress, _, publicKeyBytes, packageId, objectId := testutils.SetupTestEnv(t, ctx, _logger, gasLimit)

	// ChainWriter configuration
	chainWriterConfig := config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			"counter": {
				Name:     "counter",
				ModuleID: "counter",
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
				},
			},
		},
	}

	_logger.Infow("ChainWriterConfig", "config", chainWriterConfig)

	chainWriter, err := chainwriter.NewSuiChainWriter(_logger, txManager, chainWriterConfig, false)
	require.NoError(t, err)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	err = chainWriter.Start(ctx)
	require.NoError(t, err)
	err = txManager.Start(ctx)
	require.NoError(t, err)

	defer chainWriter.Close()
	defer txManager.Close()

	// Simple map style for builder pattern
	simpleArgs := map[string]any{
		"counter": objectId,
	}

	// Get coins to use - need at least 2 coins (one for function arg, one for gas)
	coins, err := suiClient.GetCoinsByAddress(ctx, accountAddress)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins), 2, "Need at least 2 coins for this test")

	// Common validation functions
	getCounterValue := func() (string, error) {
		objectDetails, callBackErr := suiClient.ReadObjectId(ctx, objectId)
		if callBackErr != nil {
			return "", callBackErr
		}

		return objectDetails.Content.SuiMoveObject.Fields["value"].(string), nil
	}
	//getCoinBalance := func() (string, error) {
	//	return testCoin.Balance, nil
	//}

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
		args             map[string]any
		expectError      error
		expectedResult   string
		status           commonTypes.TransactionStatus
		numberAttemps    int
		getExpectedValue func() (string, error) // Function to fetch the expected value for comparison
	}{
		{
			name:             "Test ChainWriter with valid parameters",
			txID:             "test-txID",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:           accountAddress,
			contractName:     "counter",
			functionName:     "increment",
			args:             map[string]any{"counter": objectId},
			expectError:      nil,
			expectedResult:   "0",
			status:           commonTypes.Finalized,
			numberAttemps:    1,
			getExpectedValue: getCounterValue,
		},
		{
			name:             "Test ChainWriter with PTB using builder",
			txID:             "test-ptb-txID",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:           accountAddress,
			contractName:     config.PTBChainWriterModuleName,
			functionName:     "ptb_call",
			args:             simpleArgs,
			expectError:      nil,
			expectedResult:   "1",
			status:           commonTypes.Finalized,
			numberAttemps:    1,
			getExpectedValue: getCounterValue,
		},
		{
			name:             "Test ChainWriter with missing argument for PTB using builder",
			txID:             "test-ptb-txID-missing-arg",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:           accountAddress,
			contractName:     config.PTBChainWriterModuleName,
			functionName:     "ptb_call",
			args:             map[string]any{},
			expectError:      errors.New("required parameter counter has no value"),
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:             "Test ChainWriter with PTB using simple map",
			txID:             "test-ptb-simple-map",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:           accountAddress,
			contractName:     config.PTBChainWriterModuleName,
			functionName:     "ptb_call",
			args:             simpleArgs,
			expectError:      nil,
			expectedResult:   "2",
			status:           commonTypes.Finalized,
			numberAttemps:    1,
			getExpectedValue: getCounterValue,
		},
		{
			name:             "Test ChainWriter with invalid function call",
			txID:             "test-txID-invalid-func",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:           accountAddress,
			contractName:     "counter",
			functionName:     "nonexistent_function",
			args:             map[string]any{"counter": objectId},
			expectError:      commonTypes.ErrNotFound,
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:             "Test ChainWriter with invalid contract",
			txID:             "test-txID-invalid-contract",
			txMeta:           &commonTypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:           accountAddress,
			contractName:     "nonexistent_contract",
			functionName:     "increment",
			args:             map[string]any{"counter": objectId},
			expectError:      commonTypes.ErrNotFound,
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
		{
			name:         "Test ChainWriter with the same transaction ID",
			txID:         "test-txID",
			txMeta:       &commonTypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:       accountAddress,
			contractName: "counter",
			functionName: "increment",
			args: map[string]any{
				"counter": objectId,
			},
			expectError:      errors.New("transaction already exists"),
			expectedResult:   "",
			status:           commonTypes.Failed,
			numberAttemps:    1,
			getExpectedValue: getErrorValue,
		},
	}

	//nolint:paralleltest
	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Submit the transaction
			err = chainWriter.SubmitTransaction(ctx, scenario.contractName, scenario.functionName,
				scenario.args,
				scenario.txID, packageId,
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

				tx, err := txStore.GetTransaction(scenario.txID)
				require.NoError(t, err, "Failed to get transaction from repository")
				assert.Equal(t, scenario.numberAttemps, tx.Attempt, "Transaction attempts do not match")
			}
		})
	}
}
