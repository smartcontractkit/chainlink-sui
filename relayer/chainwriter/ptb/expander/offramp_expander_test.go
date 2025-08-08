// / DO NOT EDIT - this will be removed
//go:build unit

package expander_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/expander"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/mocks"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

const (
	//nolint:gosec // G101: This is a Sui struct type identifier, not credentials
	tokenAdminRegistryStructType = "0x6::token_admin_registry::GetPoolInfosResult"
)

var addressMappings = map[string]string{
	// CCIP
	"ccipPackageId": "0x1245ccd9b14d187b00f12f37906271f18c334fb9fc1d83aa1261acda571e8746",
	"ccipObjectRef": "0x123",

	// OffRamp
	"offRampPackageId": "0x123",
	"offRampState":     "0x456",

	// SUI
	"clockObject": "0x6",
}

// ExpectedPoolInfo represents the expected values for a test scenario
type ExpectedPoolInfo struct {
	TokenPoolPackageIds     []string
	TokenPoolStateAddresses []string
	TokenPoolModules        []string
	TokenTypes              []string
}

type ExpectedReceiverInfo struct {
	PackageIds []string
}

func GenerateBSCEncodedTokenPoolInfo() ([]any, string) {
	tokenPoolInfo := expander.GetPoolInfosResult{
		TokenPoolPackageIds: []expander.SuiAddress{
			{0x1},
		},
		TokenPoolStateAddresses: []expander.SuiAddress{
			{0x2},
		},
		TokenPoolModules: []string{"0x3"},
		TokenTypes:       []string{"0x66::link::LINK", "0x7::usdc::USDC"},
	}

	bcsBytes, err := bcs.Marshal(tokenPoolInfo)
	if err != nil {
		panic(err)
	}

	// Convert []byte to []any to match Sui's response format
	bcsAsAny := make([]any, len(bcsBytes))
	for i, b := range bcsBytes {
		bcsAsAny[i] = float64(b) // Sui responses typically use float64 for numbers
	}

	structType := tokenAdminRegistryStructType

	return bcsAsAny, structType
}

//nolint:paralleltest // This test cannot run in parallel due to shared mock expectations
func TestNewSuiPTBExpander(t *testing.T) {
	lggr := logger.Test(t)

	// Note: We're testing the structure, not the exact type compatibility
	// In real usage, a proper SuiPTBClient would be passed
	ptbExpander := expander.NewSuiPTBExpander(lggr, nil, addressMappings)

	require.NotNil(t, ptbExpander)
	assert.Equal(t, addressMappings["ccipObjectRef"], ptbExpander.AddressMappings["ccipObjectRef"])
	assert.Equal(t, addressMappings["offRampState"], ptbExpander.AddressMappings["offRampState"])
	assert.Equal(t, addressMappings["clockObject"], ptbExpander.AddressMappings["clockObject"])
}

func TestGeneratePTBCommandsForTokenPools(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	tests := []struct {
		name        string
		tokenPools  []expander.TokenPool
		expectedLen int
	}{
		{
			name: "single token pool",
			tokenPools: []expander.TokenPool{
				{
					CoinMetadata:          "0xtoken1",
					TokenType:             "0x66::link_module::LINK",
					PackageId:             "0x123",
					ModuleId:              "token_admin_registry",
					Function:              "release_or_mint",
					TokenPoolStateAddress: "0x456",
					Index:                 0,
				},
			},
			expectedLen: 1,
		},
		{
			name: "multiple token pools",
			tokenPools: []expander.TokenPool{
				{
					CoinMetadata:          "0xtoken1",
					TokenType:             "0x66::link_module::LINK",
					PackageId:             "0x123",
					ModuleId:              "token_admin_registry",
					Function:              "release_or_mint",
					TokenPoolStateAddress: "0x456",
					Index:                 0,
				},
				{
					CoinMetadata:          "0xtoken2",
					TokenType:             "0x77::usdc_module::USDC",
					PackageId:             "0x789",
					ModuleId:              "token_admin_registry",
					Function:              "release_or_mint",
					TokenPoolStateAddress: "0xabc",
					Index:                 1,
				},
			},
			expectedLen: 2,
		},
		{
			name:        "empty token pools",
			tokenPools:  []expander.TokenPool{},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ptbCommands, err := expander.GeneratePTBCommandsForTokenPools(lggr, tt.tokenPools)

			require.NoError(t, err)
			assert.Len(t, ptbCommands, tt.expectedLen)

			// Verify PTB command structure
			for i, cmd := range ptbCommands {
				assert.NotNil(t, cmd.PackageId)
				assert.NotNil(t, cmd.ModuleId)
				assert.NotNil(t, cmd.Function)
				assert.Equal(t, tt.tokenPools[i].PackageId, *cmd.PackageId)
				assert.Equal(t, tt.tokenPools[i].ModuleId, *cmd.ModuleId)
				assert.Equal(t, tt.tokenPools[i].Function, *cmd.Function)

				// Verify parameters
				assert.Len(t, cmd.Params, 5) // Should have 5 parameters

				// Check specific parameter names
				paramNames := make([]string, len(cmd.Params))
				for j, param := range cmd.Params {
					paramNames[j] = param.Name
				}
				assert.Contains(t, paramNames, "ccip_object_ref")
				assert.Contains(t, paramNames, "clock")
				assert.Contains(t, paramNames, fmt.Sprintf("pool_%d", i+1))
				assert.Contains(t, paramNames, "receiver_params")
				assert.Contains(t, paramNames, fmt.Sprintf("index_%d", i+1))
			}
		})
	}
}

func TestGenerateArgumentsForTokenPools(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	ccipStateRef := "0x123"
	clockRef := "0x456"

	tests := []struct {
		name       string
		tokenPools []expander.TokenPool
	}{
		{
			name: "single token pool",
			tokenPools: []expander.TokenPool{
				{
					TokenPoolStateAddress: "0x123",
					Index:                 0,
				},
			},
		},
		{
			name: "multiple token pools",
			tokenPools: []expander.TokenPool{
				{
					TokenPoolStateAddress: "0x123",
					Index:                 0,
				},
				{
					TokenPoolStateAddress: "0x456",
					Index:                 1,
				},
			},
		},
		{
			name:       "empty token pools",
			tokenPools: []expander.TokenPool{},
		},
	}

	for _, tt := range tests {
		//nolint:gosec // G115:
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			args, _, err := expander.GenerateArgumentsForTokenPools(ccipStateRef, clockRef, lggr, tt.tokenPools)

			require.NoError(t, err)
			assert.NotNil(t, args)

			// Verify common arguments
			assert.Equal(t, ccipStateRef, args["ccip_object_ref"])
			assert.Equal(t, clockRef, args["clock"])

			// Verify token pool specific arguments
			for i, tokenPool := range tt.tokenPools {
				poolKey := fmt.Sprintf("pool_%d", i+1)
				indexKey := fmt.Sprintf("index_%d", i+1)

				assert.Equal(t, tokenPool.TokenPoolStateAddress, args[poolKey])
				assert.Equal(t, uint64(tokenPool.Index), args[indexKey])
			}
		})
	}
}

func TestGenerateReceiverCallArguments(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	ccipObjectRef := "0x123"

	tests := []struct {
		name                 string
		messages             []ccipocr3.Message
		previousCommandIndex uint16
		expectedArgs         map[string]any
		expectError          bool
	}{
		{
			name: "single message with receiver and data",
			messages: []ccipocr3.Message{
				{
					Receiver: []byte("0xpackage::module::function"),
					Data:     []byte("test data"),
				},
			},
			previousCommandIndex: 1,
			expectedArgs: map[string]any{
				"ccip_object_ref": ccipObjectRef,
				"package_id_2":    "0xpackage",
			},
			expectError: false,
		},
		{
			name: "multiple messages with receivers and data",
			messages: []ccipocr3.Message{
				{
					Receiver: []byte("0xpackage1::module1::function1"),
					Data:     []byte("test data 1"),
				},
				{
					Receiver: []byte("0xpackage2::module2::function2"),
					Data:     []byte("test data 2"),
				},
			},
			previousCommandIndex: 1,
			expectedArgs: map[string]any{
				"ccip_object_ref": ccipObjectRef,
				"package_id_2":    "0xpackage1",
				"package_id_3":    "0xpackage2",
			},
			expectError: false,
		},
		{
			name: "message without receiver",
			messages: []ccipocr3.Message{
				{
					Data: []byte("test data"),
				},
			},
			previousCommandIndex: 1,
			expectedArgs: map[string]any{
				"ccip_object_ref": ccipObjectRef,
			},
			expectError: false,
		},
		{
			name: "message without data",
			messages: []ccipocr3.Message{
				{
					Receiver: []byte("0xpackage::module::function"),
				},
			},
			previousCommandIndex: 1,
			expectedArgs: map[string]any{
				"ccip_object_ref": ccipObjectRef,
			},
			expectError: false,
		},
		{
			name:                 "empty messages",
			messages:             []ccipocr3.Message{},
			previousCommandIndex: 1,
			expectedArgs: map[string]any{
				"ccip_object_ref": ccipObjectRef,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			args, err := expander.GenerateReceiverCallArguments(lggr, tt.messages, tt.previousCommandIndex, ccipObjectRef)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, args)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, args)

				// Verify common arguments
				assert.Equal(t, ccipObjectRef, args["ccip_object_ref"])

				// Verify receiver-specific arguments
				commandIndex := tt.previousCommandIndex + 1
				for _, message := range tt.messages {
					if len(message.Receiver) > 0 && len(message.Data) > 0 {
						packageKey := fmt.Sprintf("package_id_%d", commandIndex)
						receiverParts := strings.Split(string(message.Receiver), "::")
						assert.Equal(t, receiverParts[0], args[packageKey])
						commandIndex++
					}
				}

				// Verify total number of arguments
				expectedArgCount := 1 // ccip_object_ref
				for _, message := range tt.messages {
					if len(message.Receiver) > 0 && len(message.Data) > 0 {
						expectedArgCount++
					}
				}
				assert.Len(t, args, expectedArgCount)
			}
		})
	}
}

//nolint:paralleltest // This test cannot run in parallel due to shared mock expectations
func TestSuiPTBExpander_FilterRegisteredReceivers(t *testing.T) {
	lggr := logger.Test(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSuiPTBClient := mocks.NewMockSuiPTBClient(ctrl)

	ptbExpander := expander.NewSuiPTBExpander(lggr, mockSuiPTBClient, addressMappings)
	ctx := context.Background()
	_, publicKey, _, err := testutils.GenerateAccountKeyPair(t)
	require.NoError(t, err)

	tests := []struct {
		name          string
		messages      []ccipocr3.Message
		mockResponses []bool
		expectedCount int
		expectError   bool
		errorMessage  string
	}{
		{
			name: "single registered receiver",
			messages: []ccipocr3.Message{
				{
					Receiver: []byte("0x66::ccip_dummy_receiver::CCIPDummyReceiver"),
					Data:     []byte("test data"),
				},
			},
			mockResponses: []bool{true},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "multiple receivers with some registered",
			messages: []ccipocr3.Message{
				{
					Receiver: []byte("0x66::ccip_dummy_receiver::CCIPDummyReceiver"),
					Data:     []byte("test data 1"),
				},
				{
					Receiver: []byte("0x77::ccip_dummy_receiver::CCIPDummyReceiver"),
					Data:     []byte("test data 2"),
				},
			},
			mockResponses: []bool{true, false},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name: "no receivers",
			messages: []ccipocr3.Message{
				{
					Data: []byte("test data"),
				},
			},
			mockResponses: []bool{},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "invalid receiver format",
			messages: []ccipocr3.Message{
				{
					Receiver: []byte("invalid_format"),
					Data:     []byte("test data"),
				},
			},
			mockResponses: []bool{},
			expectedCount: 0,
			expectError:   true,
			errorMessage:  "invalid receiver format",
		},
		{
			name: "client error",
			messages: []ccipocr3.Message{
				{
					Receiver: []byte("0x66::ccip_dummy_receiver::CCIPDummyReceiver"),
					Data:     []byte("test data"),
				},
			},
			mockResponses: []bool{},
			expectedCount: 0,
			expectError:   true,
			errorMessage:  "mock error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			for i, message := range tt.messages {
				if len(message.Receiver) > 0 && len(message.Data) > 0 {
					receiverParts := strings.Split(string(message.Receiver), "::")
					if len(receiverParts) != expander.SUI_PATH_COMPONENTS_COUNT {
						continue // Skip invalid format cases
					}

					if tt.errorMessage == "mock error" {
						mockSuiPTBClient.EXPECT().
							ReadFunction(
								gomock.Any(),
								gomock.Any(),
								gomock.Any(),
								"receiver_registry",
								"is_registered_receiver",
								gomock.Any(),
								[]string{"object_id", "address"},
							).
							Return(nil, fmt.Errorf("mock error")).
							Times(1)

						break
					}

					response := tt.mockResponses[i]
					results := []any{response, "bool"}

					mockSuiPTBClient.EXPECT().
						ReadFunction(
							gomock.Any(),
							gomock.Any(),
							gomock.Any(),
							"receiver_registry",
							"is_registered_receiver",
							gomock.Any(),
							[]string{"object_id", "address"},
						).
						Return(results, nil).
						Times(1)
				}
			}

			filteredMessages, err := ptbExpander.FilterRegisteredReceivers(ctx, lggr, tt.messages, publicKey)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				assert.Nil(t, filteredMessages)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, filteredMessages)
				assert.Len(t, filteredMessages, tt.expectedCount)

				// Verify that only registered receivers are included
				for _, message := range filteredMessages {
					if len(message.Receiver) > 0 && len(message.Data) > 0 {
						receiverParts := strings.Split(string(message.Receiver), "::")
						assert.Len(t, receiverParts, expander.SUI_PATH_COMPONENTS_COUNT)
					}
				}
			}
		})
	}
}
