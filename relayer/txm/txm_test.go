//go:build integration

package txm

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type Counter struct {
	Value string `json:"value"`
}

// setupClients initializes the Sui and relayer clients.
func setupClients(t *testing.T, rpcURL string, _keystore keystore.Keystore, accountAddress string) (*client.Client, *SuiTxm, signer.SuiSigner) {
	t.Helper()

	logg, err := logger.New()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Get the private key from the keystore using the account address
	signerInstance, err := _keystore.GetSignerFromAddress(accountAddress)
	require.NoError(t, err)

	relayerClient, err := client.NewClient(logg, rpcURL, nil, 10*time.Second, &signerInstance)
	if err != nil {
		t.Fatalf("Failed to create relayer client: %v", err)
	}

	store := NewTxmStoreImpl()
	conf := DefaultConfigSet

	txManager, err := NewSuiTxm(logg, relayerClient, _keystore, conf, signerInstance, store)
	if err != nil {
		t.Fatalf("Failed to create SuiTxm: %v", err)
	}

	return relayerClient, txManager, signerInstance
}

//nolint:paralleltest
func TestEnqueueIntegration(t *testing.T) {
	// Step 1: Setup

	_logger := logger.Test(t)
	_logger.Debugw("Starting Sui node")

	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	_keystore, err := keystore.NewSuiKeystore(_logger, "", keystore.PrivateKeySigner)
	require.NoError(t, err)
	accountAddress := testutils.GetAccountAndKeyFromSui(t, _logger)

	err = testutils.FundWithFaucet(_logger, constant.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	contractPath := testutils.BuildSetup(t, "contracts/test/")
	testutils.BuildContract(t, contractPath)
	packageId, publishOutput, err := testutils.PublishContract(t, "cw_tests", contractPath, accountAddress, nil)
	require.NoError(t, err)

	counterObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	suiClient, txManager, signerInstance := setupClients(t, testutils.LocalUrl, _keystore, accountAddress)

	// Step 2: Define multiple test scenarios
	testScenarios := []struct {
		name            string
		txID            string
		txMeta          *commontypes.TxMeta
		sender          string
		function        string
		typeArgs        []string
		args            []any
		expectErr       bool
		expectedValue   string
		finalState      commontypes.TransactionStatus
		storeFinalState TransactionState
		numberAttemps   int
		drainAccount    bool
	}{
		{
			name:            "Valid enqueue test",
			txID:            "integration-test-txID-1",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::increment", packageId),
			typeArgs:        []string{"address"},
			args:            []any{counterObjectId},
			expectErr:       false,
			expectedValue:   "1",
			finalState:      commontypes.Finalized,
			storeFinalState: StateFinalized,
			numberAttemps:   1,
			drainAccount:    false,
		},
		{
			name:            "Another valid enqueue test",
			txID:            "integration-test-txID-2",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::increment", packageId),
			typeArgs:        []string{"address"},
			args:            []any{counterObjectId},
			expectErr:       false,
			expectedValue:   "2",
			finalState:      commontypes.Finalized,
			storeFinalState: StateFinalized,
			numberAttemps:   1,
			drainAccount:    false,
		},
		{
			name:            "Invalid enqueue test (wrong function)",
			txID:            "wrong-function-test-txID",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::i-do-not-exist", packageId),
			typeArgs:        []string{"address"},
			args:            []any{counterObjectId},
			expectErr:       true,
			expectedValue:   "",
			finalState:      commontypes.Failed,
			storeFinalState: StateFailed,
			numberAttemps:   1,
			drainAccount:    false,
		},
		{
			name:            "Invalid enqueue test (no gas in wallet)",
			txID:            "low-gas-test-txID",
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::increment", packageId),
			typeArgs:        []string{"address"},
			args:            []any{counterObjectId},
			expectErr:       true,
			expectedValue:   "",
			finalState:      commontypes.Failed,
			storeFinalState: StateFailed,
			numberAttemps:   1,
			drainAccount:    true,
		},
	}

	ctx := context.Background()
	err = txManager.Start(ctx)
	require.NoError(t, err, "Failed to start transaction manager")

	// Step 3: Execute each test scenario
	//nolint:paralleltest
	for _, tc := range testScenarios {
		t.Run(tc.name, func(t *testing.T) {
			if tc.drainAccount {
				txManager.lggr.Infow("Draining account coins from account address", accountAddress)
				coins, err := suiClient.GetCoinsByAddress(ctx, accountAddress)
				burnAddress := "0x000000000000000000000000000000000000dead"
				require.NoError(t, err, "Failed to get coin objects")
				err = testutils.DrainAccountCoins(ctx, txManager.lggr, &signerInstance, suiClient, coins, burnAddress)
				require.NoError(t, err, "Failed to drain account coins")

				// Wait a moment for transactions to be confirmed
				time.Sleep(2 * time.Second)

				coins, err = suiClient.GetCoinsByAddress(ctx, accountAddress)
				require.NoError(t, err, "Failed to get coin objects")
				assert.Empty(t, coins, "Expected no coins left in the account")
			}

			tx, err := txManager.Enqueue(ctx, tc.txID, tc.txMeta,
				tc.sender, tc.function, nil, tc.typeArgs, tc.args, false)

			if tc.expectErr {
				assert.Error(t, err, "Expected an error but Enqueue succeeded")
			} else {
				require.Eventually(t, func() bool {
					status, statusErr := txManager.GetTransactionStatus(ctx, (*tx).TransactionID)
					if statusErr != nil {
						return false
					}

					return status == tc.finalState
				}, 60*time.Second, 1*time.Second, "Transaction final state not reached")

				tx2, err := txManager.transactionRepository.GetTransaction((*tx).TransactionID)
				require.NoError(t, err, "Failed to get transaction from repository")
				assert.NotNil(t, tx2.Digest, "Transaction digest should not be nil")

				transaction, err := txManager.transactionRepository.GetTransaction(tc.txID)
				require.NoError(t, err, "Failed to get transaction from repository")
				assert.Equal(t, tc.storeFinalState, transaction.State, "Transaction state should be Finalized")
				assert.Equal(t, tc.numberAttemps, transaction.Attempt, "Transaction attempts should be 1")

				objectDetails, err := suiClient.ReadObjectId(ctx, counterObjectId)
				require.NoError(t, err, "Failed to get object details")
				counter := testutils.ExtractStruct[Counter](t, objectDetails)
				assert.Contains(t, counter.Value, tc.expectedValue, "Counter value does not match")
			}
		})
	}
	txManager.Close()
}
