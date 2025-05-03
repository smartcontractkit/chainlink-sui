//go:build integration

package chainwriter_test

import (
	"context"
	"crypto/ed25519"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/test-go/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
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
	privateKey, err := testState.KeystoreGateway.GetPrivateKeyByAddress(testState.AccountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes := []byte(publicKey)
	countContract := testState.Contracts[0]
	packageId := countContract.ModuleID
	objectID := countContract.Objects[0].ObjectID
	// ChainWriter configuration
	chainWriterConfig := chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			"counter": {
				Name:     "counter",
				ModuleID: testState.Contracts[0].ModuleID,
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"increment": {
						Name:      "increment",
						PublicKey: publicKeyBytes,
						Params: []codec.SuiFunctionParam{
							{
								Name:         "counter",
								Type:         "address",
								Required:     true,
								DefaultValue: nil,
							},
						},
					},
				},
			},
			chainwriter.PTBChainWriterModuleName: {
				Name:     chainwriter.PTBChainWriterModuleName,
				ModuleID: "0x2",
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"ptb_call": {
						Name:      "ptb_call",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
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
								Order: 1,
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
	defer chainWriter.Close()
	require.NoError(t, err)

	// Test scenarios
	testScenarios := []struct {
		name           string
		txID           string
		txMeta         *commonTypes.TxMeta
		sender         string
		contractName   string
		functionName   string
		args           map[string]any
		expectError    error
		expectedResult string
		status         commonTypes.TransactionStatus
		numberAttemps  int
	}{
		{
			name:           "Test ChainWriter with valid parameters",
			txID:           "test-txID",
			txMeta:         &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:         testState.AccountAddress,
			contractName:   "counter",
			functionName:   "increment",
			args:           map[string]any{"counter": objectID},
			expectError:    nil,
			expectedResult: "1",
			status:         commonTypes.Finalized,
			numberAttemps:  1,
		},
		{
			name:           "Test ChainWriter with PTB",
			txID:           "test-ptb-txID",
			txMeta:         &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:         testState.AccountAddress,
			contractName:   chainwriter.PTBChainWriterModuleName,
			functionName:   "ptb_call",
			args:           map[string]any{"counter": objectID},
			expectError:    nil,
			expectedResult: "2",
			status:         commonTypes.Finalized,
			numberAttemps:  1,
		},
		{
			name:           "Test ChainWriter with invalid function call",
			txID:           "test-txID",
			txMeta:         &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:         testState.AccountAddress,
			contractName:   "counter",
			functionName:   "nonexistent_function",
			args:           map[string]any{"counter": objectID},
			expectError:    commonTypes.ErrNotFound,
			expectedResult: "",
			status:         commonTypes.Failed,
			numberAttemps:  1,
		},
		{
			name:           "Test ChainWriter with invalid contract",
			txID:           "test-txID",
			txMeta:         &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:         testState.AccountAddress,
			contractName:   "nonexistent_contract",
			functionName:   "increment",
			args:           map[string]any{"counter": objectID},
			expectError:    commonTypes.ErrNotFound,
			expectedResult: "",
			status:         commonTypes.Failed,
			numberAttemps:  1,
		},
		{

			name:         "Test ChainWriter with invalid arguments",
			txID:         "wrong-args-txID",
			txMeta:       &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:       testState.AccountAddress,
			contractName: "counter",
			functionName: "increment",
			args: map[string]any{
				"counter":     objectID,
				"invalid_arg": "invalid_value",
				"extra_arg":   "extra_value",
				"extra_arg2":  123,
			},
			expectError:    errors.New("argument count mismatch"),
			expectedResult: "",
			status:         commonTypes.Failed,
			numberAttemps:  1,
		},
		{

			name:         "Test ChainWriter with the same trasaction ID",
			txID:         "test-txID",
			txMeta:       &commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:       testState.AccountAddress,
			contractName: "counter",
			functionName: "increment",
			args: map[string]any{
				"counter": objectID,
			},
			expectError:    errors.New("transaction already exists"),
			expectedResult: "",
			status:         commonTypes.Failed,
			numberAttemps:  1,
		},
	}

	//nolint:paralleltest
	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Submit the transaction
			err = chainWriter.SubmitTransaction(
				ctx, scenario.contractName, scenario.functionName,
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
				}, 60*time.Second, 1*time.Second, "Transaction final state not reached")

				objectDetails, err := testState.SuiGateway.ReadObjectId(ctx, objectID)
				require.NoError(t, err)
				counter := testutils.ExtractStruct[Counter](t, objectDetails)
				assert.Contains(t, counter.Value, scenario.expectedResult, "Counter value does not match")
				tx, err := testState.TxStore.GetTransaction(scenario.txID)
				require.NoError(t, err, "Failed to get transaction from repository")
				assert.Equal(t, scenario.numberAttemps, tx.Attempt, "Transaction attempts do not match")
			}
		})
	}
}
