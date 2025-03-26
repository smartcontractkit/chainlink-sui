//go:build integration

package txm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// TodoList represents the structure of stored data.
type Counter struct {
	Value string `json:"value"`
}

// setupClients initializes the Sui and relayer clients.
func setupClients(t *testing.T, rpcURL string, _keystore keystore.Keystore, accountAddress string) (sui.ISuiAPI, *client.Client, *SuiTxm) {
	t.Helper()
	suiClient := sui.NewSuiClient(rpcURL)

	logg, err := logger.New()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	relayerClient, err := client.NewClient(logg, suiClient, nil, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to create relayer client: %v", err)
	}

	// Get the private key from the keystore using the account address
	signerInstance, err := _keystore.GetSignerFromAddress(accountAddress)
	require.NoError(t, err)

	txManager, err := NewSuiTxm(logg, relayerClient, _keystore, true, signerInstance)
	if err != nil {
		t.Fatalf("Failed to create SuiTxm: %v", err)
	}

	return suiClient, relayerClient, txManager
}

// fetchObjectDetails retrieves an object from the Sui network.
func fetchObjectDetails(t *testing.T, suiClient sui.ISuiAPI, objectID string) *models.SuiObjectResponse {
	t.Helper()
	objectDetails, err := suiClient.SuiGetObject(context.Background(), models.SuiGetObjectRequest{
		ObjectId: objectID,
		Options: models.SuiObjectDataOptions{
			ShowContent:             true,
			ShowDisplay:             true,
			ShowType:                true,
			ShowBcs:                 true,
			ShowOwner:               true,
			ShowPreviousTransaction: true,
			ShowStorageRebate:       true,
		},
	})
	if err != nil {
		t.Fatalf("Failed to get object details: %v", err)
	}

	return &objectDetails
}

// extractStruct parses object details into a struct
func extractStruct[T any](t *testing.T, payload any) *T {
	t.Helper()
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	var obj T
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	return &obj
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
	packageId, deploymentOutput, err := testutils.PublishContract(t, "cw_tests", contractPath, accountAddress, nil)
	require.NoError(t, err)

	listObjectId, err := testutils.ExtractObjectId(t, deploymentOutput, "Counter")
	require.NoError(t, err)

	rpcURL := testutils.LocalUrl
	suiClient, _, txManager := setupClients(t, rpcURL, _keystore, accountAddress)

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
			args:          []any{listObjectId},
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
			args:          []any{listObjectId},
			expectErr:     false,
			expectedValue: "2",
		},
	}

	// Step 3: Execute each test scenario
	//nolint:paralleltest
	for _, tc := range testScenarios {
		t.Run(tc.name, func(t *testing.T) {
			err := txManager.Enqueue(context.Background(), tc.txID, tc.txMeta,
				tc.sender, tc.function, nil, tc.typeArgs, tc.args, false)

			if tc.expectErr {
				assert.Error(t, err, "Expected an error but Enqueue succeeded")
			} else {
				// Step 4: Validate results
				objectDetails := fetchObjectDetails(t, suiClient, listObjectId)
				counter := extractStruct[Counter](t, objectDetails.Data.Content.Fields)
				assert.Contains(t, counter.Value, tc.expectedValue, "Counter value does not match")
			}
		})
	}
}
