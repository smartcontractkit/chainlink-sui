//go:build integration

package txm

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// TODO: uncomment the following struct when
// type Counter struct {
// 	Value string `json:"value"`
// }

// setupClients initializes the Sui and relayer clients.
func setupClients(t *testing.T, rpcURL string, _keystore keystore.Keystore, accountAddress string) (sui.ISuiAPI, *client.Client, *SuiTxm) {
	t.Helper()
	suiClient := sui.NewSuiClient(rpcURL)

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

	return suiClient, relayerClient, txManager
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
	packageId, _, err := testutils.PublishContract(t, "test", contractPath, accountAddress, nil)
	require.NoError(t, err)

	initializeOutput := testutils.CallContractFromCLI(t, packageId, accountAddress, "counter", "initialize", nil)
	require.NoError(t, err)

	counterObjectId, err := testutils.QueryCreatedObjectID(initializeOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// TODO: add suiClient in another PR
	_, _, txManager := setupClients(t, testutils.LocalUrl, _keystore, accountAddress)

	// Step 2: Define multiple test scenarios
	testScenarios := []struct {
		name          string
		txID          string
		txMeta        *commontypes.TxMeta
		sender        string
		function      string
		typeArgs      []string
		args          []any
		expectErr     bool
		expectedValue string
	}{
		{
			name:          "Basic enqueue test",
			txID:          "integration-test-txID-1",
			txMeta:        &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:        accountAddress,
			function:      fmt.Sprintf("%s::counter::increment", packageId),
			typeArgs:      []string{"address"},
			args:          []any{counterObjectId},
			expectErr:     false,
			expectedValue: "1",
		},
		{
			name:          "Another enqueue test",
			txID:          "integration-test-txID-2",
			txMeta:        &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:        accountAddress,
			function:      fmt.Sprintf("%s::counter::increment", packageId),
			typeArgs:      []string{"address"},
			args:          []any{counterObjectId},
			expectErr:     false,
			expectedValue: "2",
		},
	}

	ctx := context.Background()
	err = txManager.Start(ctx)
	require.NoError(t, err, "Failed to start transaction manager")

	// Step 3: Execute each test scenario
	//nolint:paralleltest
	for _, tc := range testScenarios {
		t.Run(tc.name, func(t *testing.T) {
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

					return status == commontypes.Unconfirmed
				}, 60*time.Second, 1*time.Second, "Transaction should eventually reach expected status")

				transaction, err := txManager.transactionRepository.GetTransaction(tc.txID)
				require.NoError(t, err, "Failed to get transaction from repository")
				assert.Equal(t, StateSubmitted, transaction.State, "Transaction state should be Unconfirmed")
				assert.Equal(t, 1, transaction.Attempt, "Transaction attempts should be 1")

				// TODO: this will be moved to a separate test once we implmement the confirmer routine
				// objectDetails := fetchObjectDetails(t, suiClient, counterObjectId)
				// counter := extractStruct[Counter](t, objectDetails.Data.Content.Fields)
				// assert.Contains(t, counter.Value, tc.expectedValue, "Counter value does not match")
			}
		})
	}
	txManager.Close()
}
