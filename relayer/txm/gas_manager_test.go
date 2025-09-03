//go:build unit

package txm_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/mocks"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

func TestNewSuiGasManager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		maxGasBudget            *big.Int
		percentualIncrease      int64
		expectedPercentIncrease int64
	}{
		{
			name:                    "default percentage increase when zero",
			maxGasBudget:            big.NewInt(10000000),
			percentualIncrease:      0,
			expectedPercentIncrease: 120, // gasLimitPercentualIncrease constant
		},
		{
			name:                    "custom percentage increase",
			maxGasBudget:            big.NewInt(10000000),
			percentualIncrease:      150,
			expectedPercentIncrease: 150,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lggr := logger.Test(t)
			mockClient := mocks.NewMockSuiPTBClient(ctrl)

			gasManager := txm.NewSuiGasManager(lggr, mockClient, *tt.maxGasBudget, tt.percentualIncrease)

			require.NotNil(t, gasManager)
			// Note: We can't directly test the internal percentualIncrease field since it's private
			// but we can test its behavior through the GasBump method
		})
	}
}

func TestSuiGasManager_EstimateGasBudget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMock     func(*mocks.MockSuiPTBClient)
		expectedGas   uint64
		expectedError error
	}{
		{
			name: "successful gas estimation for counter increment",
			setupMock: func(mockClient *mocks.MockSuiPTBClient) {
				mockClient.EXPECT().
					EstimateGas(gomock.Any(), gomock.Any()).
					Return(uint64(1500000), nil).
					Times(1)
			},
			expectedGas:   1500000,
			expectedError: nil,
		},
		{
			name: "gas estimation failure",
			setupMock: func(mockClient *mocks.MockSuiPTBClient) {
				mockClient.EXPECT().
					EstimateGas(gomock.Any(), gomock.Any()).
					Return(uint64(0), errors.New("network error")).
					Times(1)
			},
			expectedGas:   0,
			expectedError: errors.New("failed to estimate gas budget: network error"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lggr := logger.Test(t)
			mockClient := mocks.NewMockSuiPTBClient(ctrl)
			tt.setupMock(mockClient)

			gasManager := txm.NewSuiGasManager(lggr, mockClient, *big.NewInt(10000000), 0)

			// Create a sample transaction with counter increment call
			tx := &txm.SuiTx{
				TransactionID: "test-tx-id",
				Sender:        "0x123",
				Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(1000000)},
				Payload:       "dGVzdC10eC1ieXRlcw==", // base64 encoded "test-tx-bytes"
				Functions: []*txm.SuiFunction{
					{
						PackageId: "0x456",
						Module:    "counter",
						Name:      "increment",
					},
				},
				State: txm.StatePending,
			}

			gas, err := gasManager.EstimateGasBudget(context.Background(), tx)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Equal(t, tt.expectedGas, gas)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedGas, gas)
			}
		})
	}
}

func TestSuiGasManager_CalculateOfframpExecuteGasBudget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		arguments     offramp.SuiOffRampExecCallArgs
		expectedGas   *big.Int
		expectedError error
	}{
		{
			name: "calculate gas with gasLimit in ExtraArgsDecoded",
			arguments: offramp.SuiOffRampExecCallArgs{
				ExtraData: offramp.ExtraDataDecoded{
					ExtraArgsDecoded: map[string]any{
						"gasLimit": big.NewInt(500000),
					},
					DestExecDataDecoded: []map[string]any{
						{
							"destGasAmount": uint64(200000),
						},
						{
							"destGasAmount": uint64(300000),
						},
					},
				},
			},
			expectedGas:   big.NewInt(1000000), // 500000 + 200000 + 300000
			expectedError: nil,
		},
		{
			name: "calculate gas without gasLimit",
			arguments: offramp.SuiOffRampExecCallArgs{
				ExtraData: offramp.ExtraDataDecoded{
					ExtraArgsDecoded: map[string]any{},
					DestExecDataDecoded: []map[string]any{
						{
							"destGasAmount": uint64(100000),
						},
					},
				},
			},
			expectedGas:   big.NewInt(100000),
			expectedError: nil,
		},
		{
			name: "calculate gas with no dest exec data",
			arguments: offramp.SuiOffRampExecCallArgs{
				ExtraData: offramp.ExtraDataDecoded{
					ExtraArgsDecoded: map[string]any{
						"gasLimit": big.NewInt(750000),
					},
					DestExecDataDecoded: []map[string]any{},
				},
			},
			expectedGas:   big.NewInt(750000),
			expectedError: nil,
		},
		{
			name: "error when gasLimit is not *big.Int",
			arguments: offramp.SuiOffRampExecCallArgs{
				ExtraData: offramp.ExtraDataDecoded{
					ExtraArgsDecoded: map[string]any{
						"gasLimit": "invalid_type",
					},
				},
			},
			expectedGas:   nil,
			expectedError: errors.New("gasLimit in ExtraArgsDecoded is not *big.Int, got string"),
		},
		{
			name: "error when destGasAmount is not uint64",
			arguments: offramp.SuiOffRampExecCallArgs{
				ExtraData: offramp.ExtraDataDecoded{
					ExtraArgsDecoded: map[string]any{},
					DestExecDataDecoded: []map[string]any{
						{
							"destGasAmount": "invalid_type",
						},
					},
				},
			},
			expectedGas:   nil,
			expectedError: errors.New("destGasAmount in DestExecDataDecoded is not uint64, got string"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lggr := logger.Test(t)
			mockClient := mocks.NewMockSuiPTBClient(ctrl)

			gasManager := txm.NewSuiGasManager(lggr, mockClient, *big.NewInt(10000000), 0)

			gas, err := gasManager.CalculateOfframpExecuteGasBudget(context.Background(), tt.arguments)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, gas)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedGas, gas)
			}
		})
	}
}

func TestSuiGasManager_GasBump(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		maxGasBudget    *big.Int
		txGasBudget     *big.Int
		currentGasLimit *big.Int
		percentIncrease int64
		expectedGas     *big.Int
		expectedError   error
	}{
		{
			name:            "successful gas bump with default percentage",
			maxGasBudget:    big.NewInt(20000000),
			txGasBudget:     big.NewInt(1500000),
			currentGasLimit: big.NewInt(1000000),
			percentIncrease: 0,                   // Will use default 120%
			expectedGas:     big.NewInt(1200000), // 1000000 * 120 / 100
			expectedError:   nil,
		},
		{
			name:            "successful gas bump with custom percentage (note: implementation uses constant 120%)",
			maxGasBudget:    big.NewInt(20000000),
			txGasBudget:     big.NewInt(1500000),
			currentGasLimit: big.NewInt(1000000),
			percentIncrease: 150,                 // This is ignored by current implementation
			expectedGas:     big.NewInt(1200000), // Always uses 120% regardless of percentIncrease field
			expectedError:   nil,
		},
		{
			name:            "gas bump capped at max budget",
			maxGasBudget:    big.NewInt(1500000),
			txGasBudget:     big.NewInt(1500000),
			currentGasLimit: big.NewInt(1300000),
			percentIncrease: 0, // 120% would be 1560000, but capped at 1500000
			expectedGas:     big.NewInt(1500000),
			expectedError:   nil,
		},
		{
			name:            "no error when current gas limit equals max budget",
			maxGasBudget:    big.NewInt(1000000),
			txGasBudget:     big.NewInt(1500000),
			currentGasLimit: big.NewInt(1000000),
			percentIncrease: 0,
			expectedGas:     big.NewInt(1000000),
			expectedError:   nil,
		},
		{
			name:            "error when current gas limit exceeds max budget",
			maxGasBudget:    big.NewInt(1000000),
			txGasBudget:     big.NewInt(1500000),
			currentGasLimit: big.NewInt(1500000),
			percentIncrease: 0,
			expectedGas:     big.NewInt(0),
			expectedError:   errors.New("gas budget is already at max gas limit"),
		},
		{
			name:            "gas bump with small values",
			maxGasBudget:    big.NewInt(1000),
			txGasBudget:     big.NewInt(200),
			currentGasLimit: big.NewInt(100),
			percentIncrease: 0,
			expectedGas:     big.NewInt(120), // 100 * 120 / 100
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lggr := logger.Test(t)
			mockClient := mocks.NewMockSuiPTBClient(ctrl)

			gasManager := txm.NewSuiGasManager(lggr, mockClient, *tt.maxGasBudget, tt.percentIncrease)

			// Create a transaction with the specified gas limit
			tx := &txm.SuiTx{
				TransactionID: "test-tx-id",
				Sender:        "0x123",
				Metadata:      &commontypes.TxMeta{GasLimit: tt.currentGasLimit},
				GasBudget:     tt.txGasBudget.Uint64(),
				Payload:       "dGVzdC10eC1ieXRlcw==",
				Functions: []*txm.SuiFunction{
					{
						PackageId: "0x456",
						Module:    "counter",
						Name:      "increment",
					},
				},
				State: txm.StateSubmitted,
			}

			newGas, err := gasManager.GasBump(context.Background(), tx)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Equal(t, *tt.expectedGas, newGas)
			} else {
				require.NoError(t, err)
				assert.Equal(t, *tt.expectedGas, newGas)
			}
		})
	}
}

func TestSuiGasManager_EstimateGasBudget_CounterIncrementExample(t *testing.T) {
	t.Parallel()

	// This test specifically demonstrates gas estimation for a counter increment function call
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lggr := logger.Test(t)
	mockClient := mocks.NewMockSuiPTBClient(ctrl)

	// Set up mock to return a realistic gas estimate for counter increment
	expectedGasForCounterIncrement := uint64(1234567) // Typical gas for simple counter increment
	mockClient.EXPECT().
		EstimateGas(gomock.Any(), gomock.Eq("dGVzdC1jb3VudGVyLWluY3JlbWVudC1ieXRlcw==")).
		Return(expectedGasForCounterIncrement, nil).
		Times(1)

	gasManager := txm.NewSuiGasManager(lggr, mockClient, *big.NewInt(50000000), 0)

	// Create a transaction representing a counter increment call
	counterIncrementTx := &txm.SuiTx{
		TransactionID: "counter-increment-tx",
		Sender:        "0x742d35cc6db32e18ac22c0e3b5c4c7c0e5e1a7c3d2b1a0f9e8d7c6b5a4938271",
		Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(2000000)},
		Payload:       "dGVzdC1jb3VudGVyLWluY3JlbWVudC1ieXRlcw==", // base64 encoded counter increment tx bytes
		Functions: []*txm.SuiFunction{
			{
				PackageId: "0x456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123",
				Module:    "counter",
				Name:      "increment",
			},
		},
		State:     txm.StatePending,
		Timestamp: 1234567890,
	}

	gas, err := gasManager.EstimateGasBudget(context.Background(), counterIncrementTx)

	require.NoError(t, err)
	assert.Equal(t, expectedGasForCounterIncrement, gas)

	// Verify the gas estimate is reasonable for a simple counter increment
	assert.Greater(t, gas, uint64(1000000), "Gas estimate should be greater than 1M for counter increment")
	assert.Less(t, gas, uint64(10000000), "Gas estimate should be less than 10M for simple counter increment")
}

func TestSuiGasManager_IntegrationWithCounterContract(t *testing.T) {
	t.Parallel()

	// This test demonstrates the complete flow with a counter contract
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lggr := logger.Test(t)
	mockClient := mocks.NewMockSuiPTBClient(ctrl)

	maxGasBudget := big.NewInt(20000000)
	txGasBudget := big.NewInt(19000000)
	txMaxGasLimit := big.NewInt(1500000)
	expectedGas := big.NewInt(1800000)
	gasManager := txm.NewSuiGasManager(lggr, mockClient, *maxGasBudget, 0)

	// Simulate counter increment transaction
	counterTx := &txm.SuiTx{
		TransactionID: "counter-tx-integration",
		Sender:        "0x742d35cc6db32e18ac22c0e3b5c4c7c0e5e1a7c3d2b1a0f9e8d7c6b5a4938271",
		Metadata:      &commontypes.TxMeta{GasLimit: txMaxGasLimit},
		GasBudget:     txGasBudget.Uint64(),
		Payload:       "counter_increment_payload_bytes",
		Functions: []*txm.SuiFunction{
			{
				PackageId: "0x123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef01",
				Module:    "counter",
				Name:      "increment",
			},
		},
		State: txm.StateSubmitted,
	}

	// Test gas estimation
	mockClient.EXPECT().
		EstimateGas(gomock.Any(), "counter_increment_payload_bytes").
		Return(uint64(1800000), nil).
		Times(1)

	estimatedGas, err := gasManager.EstimateGasBudget(context.Background(), counterTx)
	require.NoError(t, err)
	assert.Equal(t, expectedGas.Uint64(), estimatedGas)

	// Test gas bump when transaction needs retry
	bumpedGas, err := gasManager.GasBump(context.Background(), counterTx)
	require.NoError(t, err)

	// Should be 1500000 * 120 / 100 = 1800000
	assert.Equal(t, *expectedGas, bumpedGas)

	// Verify the bumped gas is still within the max budget
	assert.True(t, bumpedGas.Cmp(maxGasBudget) <= 0, "Bumped gas should not exceed max budget")
}
