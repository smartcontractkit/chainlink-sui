//go:build integration

package txm_test

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

const WAIT_TIME_NEXT_TEST = 2 * time.Second

type Counter struct {
	Value string `json:"value"`
}

//nolint:paralleltest
func TestEnqueueIntegration(t *testing.T) {
	t.Skip("Skipping unstable test")

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

	// Used to wait for the tear down of one test before starting the next
	// since they both depend on the Sui node running on the same port
	time.Sleep(WAIT_TIME_NEXT_TEST)

	_keystore, err := keystore.NewSuiKeystore(_logger, "")
	require.NoError(t, err)
	accountAddress := testutils.GetAccountAndKeyFromSui(t, _logger)

	privateKey, err := _keystore.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)

	publicKey := privateKey.Public().(ed25519.PublicKey)

	err = testutils.FundWithFaucet(_logger, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	contractPath := testutils.BuildSetup(t, "contracts/test/")
	testutils.BuildContract(t, contractPath)
	packageId, publishOutput, err := testutils.PublishContract(t, "cw_tests", contractPath, accountAddress, nil)
	require.NoError(t, err)

	counterObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	suiClient, txManager, transactionRepository := testutils.SetupClients(t, testutils.LocalUrl, _keystore)

	// Step 2: Define multiple test scenarios
	testScenarios := []struct {
		name            string
		txID            string
		signerPublicKey []byte
		txMeta          *commontypes.TxMeta
		sender          string
		function        string
		paramTypes      []string
		args            []any
		expectErr       bool
		expectedValue   string
		finalState      commontypes.TransactionStatus
		storeFinalState txm.TransactionState
		numberAttemps   int
		drainAccount    bool
	}{
		{
			name:            "Valid enqueue test",
			txID:            "integration-test-txID-1",
			signerPublicKey: publicKey,
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::increment", packageId),
			paramTypes:      []string{"object_id"},
			args:            []any{counterObjectId},
			expectErr:       false,
			expectedValue:   "1",
			finalState:      commontypes.Finalized,
			storeFinalState: txm.StateFinalized,
			numberAttemps:   1,
			drainAccount:    false,
		},
		{
			name:            "Another valid enqueue test",
			txID:            "integration-test-txID-2",
			signerPublicKey: publicKey,
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::increment", packageId),
			paramTypes:      []string{"object_id"},
			args:            []any{counterObjectId},
			expectErr:       false,
			expectedValue:   "2",
			finalState:      commontypes.Finalized,
			storeFinalState: txm.StateFinalized,
			numberAttemps:   1,
			drainAccount:    false,
		},
		{
			name:            "Invalid enqueue test (wrong function)",
			txID:            "wrong-function-test-txID",
			signerPublicKey: publicKey,
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::i-do-not-exist", packageId),
			paramTypes:      []string{"object_id"},
			args:            []any{counterObjectId},
			expectErr:       false,
			expectedValue:   "",
			finalState:      commontypes.Fatal,
			storeFinalState: txm.StateFailed,
			numberAttemps:   1,
			drainAccount:    false,
		},
		{
			name:            "Invalid enqueue test (no gas in wallet)",
			txID:            "low-gas-test-txID",
			signerPublicKey: publicKey,
			txMeta:          &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
			sender:          accountAddress,
			function:        fmt.Sprintf("%s::counter::increment", packageId),
			paramTypes:      []string{"object_id"},
			args:            []any{counterObjectId},
			expectErr:       true,
			expectedValue:   "",
			finalState:      commontypes.Failed,
			storeFinalState: txm.StateFailed,
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
				addr, err := client.GetAddressFromPublicKey(tc.signerPublicKey)

				require.NoError(t, err, "Failed to get address from public key")
				_logger.Infow("Draining account coins from account address", addr)
				coins, err := suiClient.GetCoinsByAddress(ctx, addr)
				require.NoError(t, err, "Failed to get coin objects")
				burnAddress := "0x000000000000000000000000000000000000dead"
				err = testutils.DrainAccountCoins(ctx, _logger, addr, _keystore, suiClient, coins, burnAddress)
				require.NoError(t, err, "Failed to drain account coins")

				// Wait a moment for transactions to be confirmed
				time.Sleep(2 * time.Second)

				coins, err = suiClient.GetCoinsByAddress(ctx, addr)
				require.NoError(t, err, "Failed to get coin objects")
				assert.Empty(t, coins, "Expected no coins left in the account")
			}

			tx, err := txManager.Enqueue(ctx, tc.txID, tc.txMeta,
				tc.signerPublicKey, tc.function, nil, tc.paramTypes, tc.args, false)

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

				tx2, err := transactionRepository.GetTransaction((*tx).TransactionID)
				require.NoError(t, err, "Failed to get transaction from repository")
				assert.NotNil(t, tx2.Digest, "Transaction digest should not be nil")

				transaction, err := transactionRepository.GetTransaction(tc.txID)
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

//nolint:paralleltest
func TestEnqueuePTBIntegration(t *testing.T) {
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

	// Used to wait for the tear down of one test before starting the next
	// since they both depend on the Sui node running on the same port
	time.Sleep(WAIT_TIME_NEXT_TEST)
	testState := testutils.BootstrapTestEnvironment(t, testutils.CLI, metadata)
	txManager := testState.TxManager

	numberFaucetCalls := 3

	for range numberFaucetCalls {
		err = testutils.FundWithFaucet(_logger, testutils.SuiLocalnet, testState.AccountAddress)
		require.NoError(t, err)
	}

	privateKey, err := testState.KeystoreGateway.GetPrivateKeyByAddress(testState.AccountAddress)
	require.NoError(t, err)

	publicKey := privateKey.Public().(ed25519.PublicKey)
	pubKeyBytes := []byte(publicKey)

	_logger.Infow("Test environment bootstrapped")

	countContract := testState.Contracts[0]
	packageId := countContract.ModuleID
	objectId := countContract.Objects[0].ObjectID

	chainWriterConfig := chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			"counter": {
				Name:     countContract.Name,
				ModuleID: packageId,
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"ptb_call": {
						Name:      "ptb_call",
						PublicKey: pubKeyBytes,
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
							},
						},
					},
				},
			},
		},
	}

	ptbConstructor := chainwriter.NewPTBConstructor(chainWriterConfig, testState.SuiGateway, _logger)

	gasLimit := int64(200000000000)

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
			sender:          testState.AccountAddress,
			signerPublicKey: pubKeyBytes,
			contractName:    chainwriter.PTBChainWriterModuleName,
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
			sender:          testState.AccountAddress,
			signerPublicKey: pubKeyBytes,
			contractName:    chainwriter.PTBChainWriterModuleName,
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
			sender:          testState.AccountAddress,
			signerPublicKey: pubKeyBytes,
			contractName:    chainwriter.PTBChainWriterModuleName,
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
			sender:          testState.AccountAddress,
			signerPublicKey: pubKeyBytes,
			contractName:    chainwriter.PTBChainWriterModuleName,
			functionName:    "ptb_call",
			args:            map[string]any{"counter": objectId},
			expectError:     nil,
			expectedResult:  "3",
			status:          commontypes.Finalized,
			numberAttemps:   1,
		},
	}

	ctx := context.Background()
	err = txManager.Start(ctx)
	require.NoError(t, err, "Failed to start transaction manager")

	// Step 3: Execute each test scenario
	//nolint:paralleltest
	for _, tc := range testScenarios {
		t.Run(tc.name, func(t *testing.T) {
			arg := chainwriter.Arguments{
				Args: tc.args.(map[string]any),
			}
			ptb, err := ptbConstructor.BuildPTBCommands(ctx, "counter", tc.functionName, arg, nil)
			if tc.expectError != nil {
				if err != nil {
					// Expected error occurred during PTB command building
					return
				}
			} else {
				require.NoError(t, err, "Failed to build PTB commands")
				tx, err := txManager.EnqueuePTB(ctx, tc.txID, tc.txMeta, tc.signerPublicKey, ptb, false)
				require.NoError(t, err, "Failed to enqueue PTB")

				require.Eventually(t, func() bool {
					status, statusErr := txManager.GetTransactionStatus(ctx, (*tx).TransactionID)
					if statusErr != nil {
						return false
					}

					return status == tc.status
				}, 5*time.Second, 1*time.Second, "Transaction final state not reached")
			}
		})
	}
	txManager.Close()
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}
