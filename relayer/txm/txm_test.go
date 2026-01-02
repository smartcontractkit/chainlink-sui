//go:build integration

package txm_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type Counter struct {
	Value string `json:"value"`
}

//nolint:paralleltest

func TestEnqueuePTBIntegration(t *testing.T) {
	ctx := context.Background()
	_logger := logger.Test(t)
	_logger.Debugw("Starting Sui node")

	gasLimit := int64(200000000000)

	suiClient, txManager, _, accountAddress, _, publicKeyBytes, packageId, objectId := testutils.SetupTestEnv(t, ctx, _logger, gasLimit)

	chainWriterConfig := config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			"counter": {
				Name:     "Counter",
				ModuleID: packageId,
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

	ptbConstructor := ptb.NewPTBConstructor(chainWriterConfig, suiClient, _logger)

	// Step 2: Define multiple test scenarios
	testScenarios := []struct {
		name            string
		txID            string
		txMeta          *commontypes.TxMeta
		sender          string
		signerPublicKey []byte
		contractName    string
		functionName    string
		args            any
		expectError     error
		expectedResult  string
		status          commontypes.TransactionStatus
		numberAttemps   int
	}{
		{
			name:            "Test ChainWriter with valid parameters",
			txID:            "test-txID",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:          accountAddress,
			signerPublicKey: publicKeyBytes,
			contractName:    config.PTBChainWriterModuleName,
			functionName:    "ptb_call",
			args:            map[string]any{"counter": objectId},
			expectError:     nil,
			expectedResult:  "1",
			status:          commontypes.Finalized,
			numberAttemps:   1,
		},
		{
			name:            "Test ChainWriter with PTB",
			txID:            "test-ptb-txID",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:          accountAddress,
			signerPublicKey: publicKeyBytes,
			contractName:    config.PTBChainWriterModuleName,
			functionName:    "ptb_call",
			args:            map[string]any{"counter": objectId},
			expectError:     nil,
			expectedResult:  "2",
			status:          commontypes.Finalized,
			numberAttemps:   1,
		},
		{
			name:            "Test ChainWriter with missing argument for PTB",
			txID:            "test-ptb-txID-missing-arg",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:          accountAddress,
			signerPublicKey: publicKeyBytes,
			contractName:    config.PTBChainWriterModuleName,
			functionName:    "ptb_call",
			args:            map[string]any{}, // missing "counter"
			expectError:     errors.New("missing required parameter counter for command increment"),
			expectedResult:  "",
			status:          commontypes.Failed,
			numberAttemps:   1,
		},
		{
			name:            "Test ChainWriter with simple map args",
			txID:            "test-ptb-simple-map",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(gasLimit)},
			sender:          accountAddress,
			signerPublicKey: publicKeyBytes,
			contractName:    config.PTBChainWriterModuleName,
			functionName:    "ptb_call",
			args:            map[string]any{"counter": objectId},
			expectError:     nil,
			expectedResult:  "3",
			status:          commontypes.Finalized,
			numberAttemps:   3,
		},
		{
			name:            "Test ChainWriter with low gas budget requiring gas bump",
			txID:            "test-ptb-gas-management",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(1000000000)}, // Use small limit to trigger gas bumping
			sender:          accountAddress,
			signerPublicKey: publicKeyBytes,
			contractName:    config.PTBChainWriterModuleName,
			functionName:    "ptb_call",
			args:            map[string]any{"counter": objectId},
			expectError:     nil,
			expectedResult:  "4",
			status:          commontypes.Finalized,
			numberAttemps:   3, // Should succeed after gas bumps
		},
	}

	err := txManager.Start(ctx)
	require.NoError(t, err, "Failed to start transaction manager")

	functionConfig := chainWriterConfig.Modules["counter"].Functions["ptb_call"]

	// Step 3: Execute each test scenario
	//nolint:paralleltest
	for _, tc := range testScenarios {
		t.Run(tc.name, func(t *testing.T) {
			arg := config.Arguments{
				Args: tc.args.(map[string]any),
			}
			ptb, err := ptbConstructor.BuildPTBCommands(ctx, "counter", tc.functionName, arg, packageId, functionConfig)
			if tc.expectError != nil {
				require.Error(t, err, "Expected an error but BuildPTBCommands succeeded")
			} else {
				require.NoError(t, err, "Failed to build PTB commands")
				tx, err := txManager.EnqueuePTB(ctx, tc.txID, tc.txMeta, tc.signerPublicKey, ptb)
				require.NoError(t, err, "Failed to enqueue PTB")

				require.Eventually(t, func() bool {
					status, statusErr := txManager.GetTransactionStatus(ctx, (*tx).TransactionID)
					if statusErr != nil {
						return false
					}

					return status == tc.status
				}, 10*time.Second, 1*time.Second, "Transaction final state not reached")

			}
		})
	}
	txManager.Close()
}

func TestCoinReservations(t *testing.T) {
	ctx := context.Background()
	logger := logger.Test(t)
	logger.Debugw("Starting Sui node")

	gasLimit := int64(200000000000)

	suiClient, txManager, _, _, _, publicKeyBytes, packageId, objectID := testutils.SetupTestEnv(t, ctx, logger, gasLimit)

	chainWriterConfig := config.ChainWriterConfig{
		Modules: map[string]*config.ChainWriterModule{
			"counter": {
				Name:     "Counter",
				ModuleID: packageId,
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
								Function:  strPtr("increment_by"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "by",
										Type:     "u64",
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

	numTransactions := 10
	ptbConstructor := ptb.NewPTBConstructor(chainWriterConfig, suiClient, logger)
	numEnqueuedTransactions := 0

	// queue N transactions
	for i := range numTransactions {
		arg := config.Arguments{
			Args: map[string]any{
				"counter": objectID,
				"by":      uint64(i),
			},
		}
		ptb, err := ptbConstructor.BuildPTBCommands(ctx, "counter", "ptb_call", arg, packageId, chainWriterConfig.Modules["counter"].Functions["ptb_call"])
		require.NoError(t, err, "Failed to build PTB commands")

		txID := fmt.Sprintf("test-txID-%d", i)
		_, err = txManager.EnqueuePTB(ctx, txID, &commontypes.TxMeta{GasLimit: big.NewInt(gasLimit)}, publicKeyBytes, ptb)
		
		if err != nil {
			logger.Errorw("Failed to enqueue PTB", "error", err)
			continue
		}

		numEnqueuedTransactions++
	}

	// start the txm with multiple queued transactions
	err := txManager.Start(ctx)
	require.NoError(t, err, "Failed to start transaction manager")

	defer txManager.Close()

	successfulTransactions := map[string]bool{}

	require.Eventually(t, func() bool {
		successCount := 0

		for i := range numEnqueuedTransactions {
			txID := fmt.Sprintf("test-txID-%d", i)
			
			// Skip if the transaction has already been finalized
			if _, ok := successfulTransactions[txID]; ok {
				successCount++
				continue
			}
			
			status, statusErr := txManager.GetTransactionStatus(ctx, txID)
			if statusErr != nil {
				return false
			}

			if status == commontypes.Finalized {
				successCount++
				successfulTransactions[txID] = true
			}
		}

		return successCount == numEnqueuedTransactions
	}, 60*time.Second, 3*time.Second, "Transactions final state not reached")
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}
