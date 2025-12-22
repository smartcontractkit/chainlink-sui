//go:build integration

package txm_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	txm_export "github.com/smartcontractkit/chainlink-sui/relayer/txm"
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

func TestHandleLockCoinError(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)

	gasLimit := int64(200000000000)

	_, txm, store, _, _, _, _, _ :=
		testutils.SetupTestEnv(t, ctx, lggr, gasLimit)

	txID := "tx-locked-coins-1"
	initialTx := txm_export.SuiTx{
		TransactionID: txID,
		Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(1)},
	}

	err := store.AddTransaction(initialTx)
	require.NoError(t, err, "failed to seed transaction into store")

	lockedErrMsg := `failed to execute transaction: {"code":-32002,"message":"Transaction is rejected as invalid by more than 1/3 of validators by stake (non-retriable). 
					Non-retriable errors: [Object (0x23a4b83340069bd92db7ee2a22994d09f7ff1083af74a9151c9659a5a9662750, SequenceNumber(717214713), o#3R6R2XxDfWT6sXyzz4xX4r1mL6GUHNHCH64vkGwHbBJd) 
					already locked by a different transaction: TransactionDigest(HEf6wjXGWSoemesir2LC7CGnb1c4cPrZruSjTXwpgAuJ)]. Retriable errors: []"}`

	handled := txm_export.HandleLockCoinError(txm, initialTx, lockedErrMsg)
	require.True(t, handled, "expected locked-coin helper to handle the error")

	storedTx, err := store.GetTransaction(txID)
	require.NoError(t, err, "failed to read transaction from store after handling error")

	require.Equal(t, txm_export.StateFailed, storedTx.State, "transaction state should be Failed after lock-coin error")
	require.NotNil(t, storedTx.TxError, "TxError should be set on lock-coin error")
	require.Equal(t, suierrors.LockCoinErrors, storedTx.TxError.Category, "TxError category should be LockCoinErrors")

	snap := txm.SnapshotLockedCoins()
	require.NotEmpty(t, snap, "lockedCoins snapshot should not be empty after handling a locked-coin error")

	now := time.Now()
	yesterday := now.Add(-25 * time.Hour)

	txm.MarkLockedCoin("0x_old_coin", 1, yesterday)
	txm.MarkLockedCoin("0x_fresh_coin", 1, now)

	snap = txm.SnapshotLockedCoins()

	require.Len(t, snap, 2) // 2 because one 1 coins is from above `HandleLockCoinError`
	require.Contains(t, snap, txm_export.CoinKey("0x_fresh_coin", 1))
	require.NotContains(t, snap, txm_export.CoinKey("0x_old_coin", 1))
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}
