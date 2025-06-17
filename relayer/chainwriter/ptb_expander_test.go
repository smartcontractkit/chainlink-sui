//go:build unit

package chainwriter_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/mocks"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
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

func GetTestChainWriterConfig() chainwriter.ChainWriterConfig {
	return chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			chainwriter.PTBChainWriterModuleName: {
				Functions: map[string]*chainwriter.ChainWriterFunction{
					chainwriter.CCIPExecuteReportFunctionName: {
						AddressMappings: addressMappings,
					},
				},
			},
		},
	}
}

func GenerateBSCEncodedTokenPoolInfo() ([]any, string) {
	tokenPoolInfo := chainwriter.GetPoolInfosResult{
		TokenPoolPackageIds: []chainwriter.SuiAddress{
			{0x1},
		},
		TokenPoolStateAddresses: []chainwriter.SuiAddress{
			{0x2},
		},
		TokenPoolModules: []chainwriter.SuiAddress{
			{0x3},
		},
		TokenTypes: []string{"0x66::link::LINK", "0x7::usdc::USDC"},
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

	cwConfig := GetTestChainWriterConfig()

	// Note: We're testing the structure, not the exact type compatibility
	// In real usage, a proper SuiPTBClient would be passed
	expander := chainwriter.NewSuiPTBExpander(lggr, nil, cwConfig)

	require.NotNil(t, expander)
	assert.Equal(t, addressMappings["ccipObjectRef"], expander.AddressMappings["ccipObjectRef"])
	assert.Equal(t, addressMappings["offRampState"], expander.AddressMappings["offRampState"])
	assert.Equal(t, addressMappings["clockObject"], expander.AddressMappings["clockObject"])
}

// setupMockExpectations sets up the mock expectations for ReadFunction based on test scenarios
func setupMockExpectations(mockClient *mocks.MockSuiPTBClient, testCases []struct {
	name                 string
	tokenAmounts         []ccipocr3.RampTokenAmount
	expected             ExpectedPoolInfo
	expectedReceiverInfo ExpectedReceiverInfo
	expectedLen          int
	expectError          bool
}) {
	// Set up expectations for each test case individually
	for _, tc := range testCases {
		if tc.expectError {
			// For error cases, you might want to return an error
			mockClient.EXPECT().
				ReadFunction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					"token_admin_registry",
					"get_pool_infos",
					gomock.Any(),
					[]string{"object_id", "vector<address>"},
				).
				Return(nil, fmt.Errorf("mock error for test: %s", tc.name)).
				Times(1)
		} else {
			// Generate response based on the specific test case
			response := generateMockResponseForTestCase(tc)

			mockClient.EXPECT().
				ReadFunction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					"token_admin_registry",
					"get_pool_infos",
					gomock.Any(),
					[]string{"object_id", "vector<address>"},
				).
				Return(response, nil).
				Times(1)
		}
	}
}

// generateMockResponseForTestCase creates a mock response tailored to the specific test case
func generateMockResponseForTestCase(tc struct {
	name                 string
	tokenAmounts         []ccipocr3.RampTokenAmount
	expected             ExpectedPoolInfo
	expectedReceiverInfo ExpectedReceiverInfo
	expectedLen          int
	expectError          bool
}) *suiclient.ExecutionResultType {
	// Create token pool info based on the test case expectations
	tokenPoolInfo := chainwriter.GetPoolInfosResult{}

	// Convert expected hex strings back to SuiAddress for BCS encoding
	for i := range tc.expected.TokenPoolPackageIds {
		if i < len(tc.expected.TokenPoolPackageIds) {
			// Decode hex string to bytes and convert to SuiAddress
			packageIdBytes, _ := hex.DecodeString(tc.expected.TokenPoolPackageIds[i])
			var packageIdAddr chainwriter.SuiAddress
			copy(packageIdAddr[:], packageIdBytes)
			tokenPoolInfo.TokenPoolPackageIds = append(tokenPoolInfo.TokenPoolPackageIds, packageIdAddr)
		}
	}

	for i := range tc.expected.TokenPoolStateAddresses {
		if i < len(tc.expected.TokenPoolStateAddresses) {
			stateAddrBytes, _ := hex.DecodeString(tc.expected.TokenPoolStateAddresses[i])
			var stateAddr chainwriter.SuiAddress
			copy(stateAddr[:], stateAddrBytes)
			tokenPoolInfo.TokenPoolStateAddresses = append(tokenPoolInfo.TokenPoolStateAddresses, stateAddr)
		}
	}

	for i := range tc.expected.TokenPoolModules {
		if i < len(tc.expected.TokenPoolModules) {
			moduleBytes, _ := hex.DecodeString(tc.expected.TokenPoolModules[i])
			var moduleAddr chainwriter.SuiAddress
			copy(moduleAddr[:], moduleBytes)
			tokenPoolInfo.TokenPoolModules = append(tokenPoolInfo.TokenPoolModules, moduleAddr)
		}
	}

	// Add token types
	tokenPoolInfo.TokenTypes = tc.expected.TokenTypes

	// BCS encode the response
	bcsBytes, err := bcs.Marshal(tokenPoolInfo)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal token pool info for test %s: %v", tc.name, err))
	}

	// Convert to Sui response format
	bcsAsAny := make([]any, len(bcsBytes))
	for i, b := range bcsBytes {
		bcsAsAny[i] = float64(b)
	}

	structType := tokenAdminRegistryStructType
	results := []any{bcsAsAny, structType}

	return &suiclient.ExecutionResultType{
		ReturnValues: []suiclient.ReturnValueType{results},
	}
}

//nolint:paralleltest // This test cannot run in parallel due to shared mock expectations
func TestSuiPTBExpander_GetTokenPoolByTokenAddress(t *testing.T) {
	lggr := logger.Test(t)

	cwConfig := GetTestChainWriterConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSuiPTBClient := mocks.NewMockSuiPTBClient(ctrl)

	tests := []struct {
		name                 string
		tokenAmounts         []ccipocr3.RampTokenAmount
		expected             ExpectedPoolInfo
		expectedReceiverInfo ExpectedReceiverInfo
		expectedLen          int
		expectError          bool
	}{
		{
			name: "single token amount",
			tokenAmounts: []ccipocr3.RampTokenAmount{
				{
					DestTokenAddress: []byte("0x66::link::LINK"),
				},
			},
			expected: ExpectedPoolInfo{
				TokenPoolPackageIds: []string{
					hex.EncodeToString([]byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolStateAddresses: []string{
					hex.EncodeToString([]byte{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolModules: []string{
					hex.EncodeToString([]byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenTypes: []string{"0x66::link::LINK"},
			},
			expectedReceiverInfo: ExpectedReceiverInfo{},
			expectedLen:          1,
			expectError:          false,
		},
		{
			name: "multiple token amounts",
			tokenAmounts: []ccipocr3.RampTokenAmount{
				{
					DestTokenAddress: []byte("0x66::link::LINK"),
				},
				{
					DestTokenAddress: []byte("0x7::usdc::USDC"),
				},
			},
			expected: ExpectedPoolInfo{
				TokenPoolPackageIds: []string{
					hex.EncodeToString([]byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolStateAddresses: []string{
					hex.EncodeToString([]byte{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolModules: []string{
					hex.EncodeToString([]byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenTypes: []string{"0x66::link::LINK", "0x7::usdc::USDC"},
			},
			expectedReceiverInfo: ExpectedReceiverInfo{},
			expectedLen:          2,
			expectError:          false,
		},
		{
			name:         "empty token amounts",
			tokenAmounts: []ccipocr3.RampTokenAmount{},
			expected: ExpectedPoolInfo{
				TokenPoolPackageIds:     []string{},
				TokenPoolStateAddresses: []string{},
				TokenPoolModules:        []string{},
				TokenTypes:              []string{},
			},
			expectedReceiverInfo: ExpectedReceiverInfo{},
			expectedLen:          0,
			expectError:          false,
		},
		{
			name: "error scenario - network failure",
			tokenAmounts: []ccipocr3.RampTokenAmount{
				{
					DestTokenAddress: []byte("0x66::link::LINK"),
				},
			},
			expected:             ExpectedPoolInfo{}, // Empty since we expect error
			expectedReceiverInfo: ExpectedReceiverInfo{},
			expectedLen:          0,
			expectError:          true,
		},
	}

	// Set up mock expectations based on test cases
	setupMockExpectations(mockSuiPTBClient, tests)

	expander := chainwriter.NewSuiPTBExpander(lggr, mockSuiPTBClient, cwConfig)

	_, publicKey, _, err := testutils.GenerateAccountKeyPair(t, lggr)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenPools, err := expander.GetTokenPoolByTokenAddress(lggr, tt.tokenAmounts, publicKey)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, tokenPools)
			} else {
				require.NoError(t, err)
				assert.Len(t, tokenPools, tt.expectedLen)

				lggr.Debugw("tokenPools", "tokenPools", tokenPools)

				// Verify token pool structure for each result
				for i, tokenPool := range tokenPools {
					assert.Equal(t, string(tt.tokenAmounts[i].DestTokenAddress), tokenPool.CoinMetadata)
					assert.Equal(t, tt.expected.TokenTypes[i], tokenPool.TokenType)
					assert.Equal(t, tt.expected.TokenPoolPackageIds[i], tokenPool.PackageId)
					assert.Equal(t, tt.expected.TokenPoolModules[i], tokenPool.ModuleId)
					assert.Equal(t, chainwriter.OFFRAMP_TOKEN_POOL_FUNCTION_NAME, tokenPool.Function)
					assert.Equal(t, tt.expected.TokenPoolStateAddresses[i], tokenPool.TokenPoolStateAddress)
					assert.Equal(t, i, tokenPool.Index)
				}
			}
		})
	}
}

//nolint:paralleltest // This test cannot run in parallel due to shared mock expectations
func TestSuiPTBExpander_GetOffRampPTB(t *testing.T) {
	lggr := logger.Test(t)

	cwConfig := GetTestChainWriterConfig()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockSuiPTBClient := mocks.NewMockSuiPTBClient(ctrl)

	executePTBcommands := []chainwriter.ChainWriterPTBCommand{
		{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: chainwriter.AnyPointer(addressMappings["offRampPackageId"]),
			ModuleId:  chainwriter.AnyPointer("ccip_offramp"),
			Function:  chainwriter.AnyPointer("init_execute"),
			Params: []codec.SuiFunctionParam{
				{
					Name:     "ref",
					Type:     "object_id",
					Required: true,
				},
				{
					Name:     "state",
					Type:     "object_id",
					Required: true,
				},
				{
					Name:     "clock",
					Type:     "object_id",
					Required: true,
				},
				{
					Name:     "report_context",
					Type:     "vector<vector<u8>>",
					Required: true,
				},
				{
					Name:     "report",
					Type:     "vector<u8>",
					Required: true,
				},
			},
		},
		{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: chainwriter.AnyPointer(addressMappings["offRampPackageId"]),
			ModuleId:  chainwriter.AnyPointer("ccip_offramp"),
			Function:  chainwriter.AnyPointer("finish_execute"),
			Params: []codec.SuiFunctionParam{
				{
					Name:     "receiver_params",
					Type:     "ptb_dependency",
					Required: true,
					PTBDependency: &codec.PTBCommandDependency{
						CommandIndex: uint16(0),
					},
				},
			},
		},
	}

	lggr.Debugw("executePTBcommands", "commands", executePTBcommands)

	tests := []struct {
		name                 string
		args                 chainwriter.SuiOffRampExecCallArgs
		ptbConfigs           *chainwriter.ChainWriterFunction
		expectedPoolInfo     ExpectedPoolInfo
		expectedReceiverInfo ExpectedReceiverInfo
		expectError          bool
		errorMessage         string
		expectedPTBCommands  int
	}{
		{
			name: "valid configuration with single token",
			args: chainwriter.SuiOffRampExecCallArgs{
				ReportContext: [2][32]byte{},
				Report:        []byte("test report"),
				Info: ccipocr3.ExecuteReportInfo{
					AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{
						{
							Messages: []ccipocr3.Message{
								{
									TokenAmounts: []ccipocr3.RampTokenAmount{
										{
											DestTokenAddress: []byte("0x66::link::LINK"),
										},
									},
									Data:     []byte("test data"),
									Receiver: []byte("0x66::ccip_dummy_receiver::CCIPDummyReceiver"),
								},
							},
						},
					},
				},
			},
			ptbConfigs: &chainwriter.ChainWriterFunction{
				PTBCommands: executePTBcommands,
			},
			expectedPoolInfo: ExpectedPoolInfo{
				TokenPoolPackageIds: []string{
					hex.EncodeToString([]byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolStateAddresses: []string{
					hex.EncodeToString([]byte{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolModules: []string{
					hex.EncodeToString([]byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenTypes: []string{"0x66::link::LINK"},
			},
			expectedReceiverInfo: ExpectedReceiverInfo{
				PackageIds: []string{"0x66"},
			},
			expectError:         false,
			expectedPTBCommands: 4, // init_execute + 1 token pool + receiver call + finish_execute
		},
		{
			name: "invalid configuration with wrong number of PTB commands",
			args: chainwriter.SuiOffRampExecCallArgs{
				ReportContext: [2][32]byte{},
				Report:        []byte("test report"),
				Info: ccipocr3.ExecuteReportInfo{
					AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{},
				},
			},
			ptbConfigs: &chainwriter.ChainWriterFunction{
				PTBCommands: executePTBcommands[:1],
			},
			expectedPoolInfo:     ExpectedPoolInfo{}, // Empty since we expect error
			expectedReceiverInfo: ExpectedReceiverInfo{},
			expectError:          true,
			errorMessage:         "expected 2 PTB commands, got 1",
			expectedPTBCommands:  0,
		},
		{
			name: "multiple messages with multiple token amounts and no receiver",
			args: chainwriter.SuiOffRampExecCallArgs{
				ReportContext: [2][32]byte{},
				Report:        []byte("test report"),
				Info: ccipocr3.ExecuteReportInfo{
					AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{
						{
							Messages: []ccipocr3.Message{
								{
									TokenAmounts: []ccipocr3.RampTokenAmount{
										{
											DestTokenAddress: []byte("0x66::link::LINK"),
										},
										{
											DestTokenAddress: []byte("0x7::usdc::USDC"),
										},
									},
									Receiver: []byte("0x66::ccip_dummy_receiver::CCIPDummyReceiver"),
								},
								{
									TokenAmounts: []ccipocr3.RampTokenAmount{
										{
											DestTokenAddress: []byte("0x8::eth::ETH"),
										},
									},
								},
							},
						},
					},
				},
			},
			ptbConfigs: &chainwriter.ChainWriterFunction{
				PTBCommands: executePTBcommands,
			},
			expectedPoolInfo: ExpectedPoolInfo{
				TokenPoolPackageIds: []string{
					hex.EncodeToString([]byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolStateAddresses: []string{
					hex.EncodeToString([]byte{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolModules: []string{
					hex.EncodeToString([]byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenTypes: []string{"0x66::link::LINK", "0x7::usdc::USDC", "0x8::eth::ETH"},
			},
			expectedReceiverInfo: ExpectedReceiverInfo{}, // empty since we have no data field
			expectError:          false,
			expectedPTBCommands:  5, // init_execute + 3 token pools + finish_execute
		},
		{
			name: "multiple messages with multiple token amounts with receiver",
			args: chainwriter.SuiOffRampExecCallArgs{
				ReportContext: [2][32]byte{},
				Report:        []byte("test report"),
				Info: ccipocr3.ExecuteReportInfo{
					AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{
						{
							Messages: []ccipocr3.Message{
								{
									TokenAmounts: []ccipocr3.RampTokenAmount{
										{
											DestTokenAddress: []byte("0x66::link::LINK"),
										},
										{
											DestTokenAddress: []byte("0x7::usdc::USDC"),
										},
									},
									Receiver: []byte("0x66::ccip_dummy_receiver::CCIPDummyReceiver"),
									Data:     []byte("test data"),
								},
								{
									TokenAmounts: []ccipocr3.RampTokenAmount{
										{
											DestTokenAddress: []byte("0x8::eth::ETH"),
										},
									},
								},
							},
						},
					},
				},
			},
			ptbConfigs: &chainwriter.ChainWriterFunction{
				PTBCommands: executePTBcommands,
			},
			expectedPoolInfo: ExpectedPoolInfo{
				TokenPoolPackageIds: []string{
					hex.EncodeToString([]byte{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolStateAddresses: []string{
					hex.EncodeToString([]byte{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenPoolModules: []string{
					hex.EncodeToString([]byte{0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
					hex.EncodeToString([]byte{0x9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}),
				},
				TokenTypes: []string{"0x66::link::LINK", "0x7::usdc::USDC", "0x8::eth::ETH"},
			},
			expectedReceiverInfo: ExpectedReceiverInfo{
				PackageIds: []string{"0x66"},
			},
			expectError:         false,
			expectedPTBCommands: 6, // init_execute + 3 token pools + receiver call + finish_execute
		},
		{
			name: "empty messages - no tokens",
			args: chainwriter.SuiOffRampExecCallArgs{
				ReportContext: [2][32]byte{},
				Report:        []byte("test report"),
				Info: ccipocr3.ExecuteReportInfo{
					AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{
						{
							Messages: []ccipocr3.Message{},
						},
					},
				},
			},
			ptbConfigs: &chainwriter.ChainWriterFunction{
				PTBCommands: executePTBcommands,
			},
			expectedPoolInfo: ExpectedPoolInfo{
				TokenPoolPackageIds:     []string{},
				TokenPoolStateAddresses: []string{},
				TokenPoolModules:        []string{},
				TokenTypes:              []string{},
			},
			expectedReceiverInfo: ExpectedReceiverInfo{}, // empty since we have no data field
			expectError:          false,
			expectedPTBCommands:  2, // init_execute + finish_execute (no token pools)
		},
	}

	// Set up mock expectations for each test case
	for _, tc := range tests {
		if !tc.expectError {
			// Generate response based on test case expectations
			response := generateMockResponseForOffRampTest(lggr, tc)

			mockSuiPTBClient.EXPECT().
				ReadFunction(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					"token_admin_registry",
					"get_pool_infos",
					gomock.Any(),
					[]string{"object_id", "vector<address>"},
				).
				Return(response, nil).
				Times(1)

			// Add mock expectations for is_registered_receiver calls
			for _, report := range tc.args.Info.AbstractReports {
				for _, message := range report.Messages {
					if len(message.Receiver) > 0 && len(message.Data) > 0 {
						receiverParts := strings.Split(string(message.Receiver), "::")
						if len(receiverParts) != chainwriter.SUI_PATH_COMPONENTS_COUNT {
							continue
						}

						receiverAddress := fmt.Sprintf("%s::%s::%s", receiverParts[0], receiverParts[1], receiverParts[2])
						results := []any{true, "bool"} // Assume all receivers are registered for these tests

						expectedResult := &suiclient.ExecutionResultType{
							ReturnValues: []suiclient.ReturnValueType{results},
						}

						mockSuiPTBClient.EXPECT().
							ReadFunction(
								gomock.Any(),
								gomock.Any(),
								gomock.Any(),
								"ccip",
								"is_registered_receiver",
								[]any{
									addressMappings["ccipObjectRef"],
									receiverAddress,
								},
								[]string{"object_id", "address"},
							).
							Return(expectedResult, nil).
							Times(1)
					}
				}
			}
		}
	}

	expander := chainwriter.NewSuiPTBExpander(lggr, mockSuiPTBClient, cwConfig)

	_, publicKey, _, err := testutils.GenerateAccountKeyPair(t, lggr)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptbCommands, updatedArgs, err := expander.GetOffRampPTB(lggr, tt.args, tt.ptbConfigs, publicKey)

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				assert.Nil(t, ptbCommands)
				assert.Nil(t, updatedArgs)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, ptbCommands)
				assert.NotNil(t, updatedArgs)

				// Verify PTB command count
				assert.Len(t, ptbCommands, tt.expectedPTBCommands)

				// Verify that arguments contain expected keys
				args, ok := updatedArgs.(map[string]any)
				assert.True(t, ok)
				assert.Equal(t, addressMappings["ccipObjectRef"], args["ccip_state_ref"])
				assert.Equal(t, addressMappings["clockObject"], args["clock_ref"])
				assert.Equal(t, 0, args["remote_chain_selector"])

				// Calculate expected number of token amounts across all messages
				expectedTokenAmounts := 0
				for _, report := range tt.args.Info.AbstractReports {
					for _, message := range report.Messages {
						expectedTokenAmounts += len(message.TokenAmounts)
					}
				}

				lggr.Debugw("PTB commands", "ptbCommands", ptbCommands)

				// Verify command structure and ordering
				expectedTotalCommands := 2 + expectedTokenAmounts + len(tt.expectedReceiverInfo.PackageIds) // init_execute + token_pools + finish_execute
				lggr.Debugw("expectedTotalCommands", "expectedTotalCommands", expectedTotalCommands)
				assert.Len(t, ptbCommands, expectedTotalCommands, "Total PTB commands should include init, token pools, receiver calls and finish")

				if expectedTotalCommands > 0 {
					// Verify first command is init_execute
					assert.Equal(t, "init_execute", *ptbCommands[0].Function)
					assert.Equal(t, "ccip_offramp", *ptbCommands[0].ModuleId)

					// Verify the last command is finish_execute
					assert.Equal(t, "finish_execute", *ptbCommands[len(ptbCommands)-1].Function)
					assert.Equal(t, "ccip_offramp", *ptbCommands[len(ptbCommands)-1].ModuleId)
				}

				// Verify token pool commands are in the middle (if any)
				for i := 1; i <= expectedTokenAmounts; i++ {
					if i < len(ptbCommands) {
						tokenPoolCmd := ptbCommands[i]
						assert.Equal(t, "release_or_mint", *tokenPoolCmd.Function)
						assert.Equal(t, tt.expectedPoolInfo.TokenPoolModules[i-1], *tokenPoolCmd.ModuleId)

						// Verify PTB dependency structure in token pool commands
						receiverParamFound := false
						for _, param := range tokenPoolCmd.Params {
							if param.Name == "receiver_params" && param.PTBDependency != nil {
								receiverParamFound = true
								// Token pool commands should reference command at index i (1-based indexing adjusted for 0-based)
								//nolint:gosec // G115: PTB commands are typically small in number, overflow extremely unlikely
								assert.Equal(t, uint16(i-1), param.PTBDependency.CommandIndex,
									"Token pool command %d should reference command index %d", i, i-1)
							}
						}
						assert.True(t, receiverParamFound, "Token pool command should have receiver_params with PTB dependency")
					}
				}

				// Verify the last command's PTB dependency points to the last token pool command
				if expectedTokenAmounts > 0 && len(ptbCommands) > 1 {
					lastCmd := ptbCommands[len(ptbCommands)-1]
					lastCmdDependencyFound := false
					for _, param := range lastCmd.Params {
						if param.PTBDependency != nil {
							lastCmdDependencyFound = true
							// Should reference the last token pool command (which is at index expectedTokenAmounts)
							//nolint:gosec // G115: PTB commands are typically small in number, overflow extremely unlikely
							expectedIndex := uint16(expectedTokenAmounts) + uint16(len(tt.expectedReceiverInfo.PackageIds))
							assert.Equal(t, expectedIndex, param.PTBDependency.CommandIndex,
								"Last command should reference the last token pool command at index %d", expectedIndex)
						}
					}
					assert.True(t, lastCmdDependencyFound, "Last command should have a PTB dependency")
				}
			}
		})
	}
}

// generateMockResponseForOffRampTest creates a mock response tailored to the specific OffRamp test case
func generateMockResponseForOffRampTest(lggr logger.Logger, tc struct {
	name                 string
	args                 chainwriter.SuiOffRampExecCallArgs
	ptbConfigs           *chainwriter.ChainWriterFunction
	expectedPoolInfo     ExpectedPoolInfo
	expectedReceiverInfo ExpectedReceiverInfo
	expectError          bool
	errorMessage         string
	expectedPTBCommands  int
}) *suiclient.ExecutionResultType {
	// Create token pool info based on the test case expectations
	tokenPoolInfo := chainwriter.GetPoolInfosResult{}

	lggr.Debugw("Running mock response for test", "test", tc.name)

	// Convert expected hex strings back to SuiAddress for BCS encoding
	for i := range tc.expectedPoolInfo.TokenPoolPackageIds {
		if i < len(tc.expectedPoolInfo.TokenPoolPackageIds) {
			// Decode hex string to bytes and convert to SuiAddress
			packageIdBytes, _ := hex.DecodeString(tc.expectedPoolInfo.TokenPoolPackageIds[i])
			var packageIdAddr chainwriter.SuiAddress
			copy(packageIdAddr[:], packageIdBytes)
			tokenPoolInfo.TokenPoolPackageIds = append(tokenPoolInfo.TokenPoolPackageIds, packageIdAddr)
		}
	}

	for i := range tc.expectedPoolInfo.TokenPoolStateAddresses {
		if i < len(tc.expectedPoolInfo.TokenPoolStateAddresses) {
			stateAddrBytes, _ := hex.DecodeString(tc.expectedPoolInfo.TokenPoolStateAddresses[i])
			var stateAddr chainwriter.SuiAddress
			copy(stateAddr[:], stateAddrBytes)
			tokenPoolInfo.TokenPoolStateAddresses = append(tokenPoolInfo.TokenPoolStateAddresses, stateAddr)
		}
	}

	for i := range tc.expectedPoolInfo.TokenPoolModules {
		if i < len(tc.expectedPoolInfo.TokenPoolModules) {
			moduleBytes, _ := hex.DecodeString(tc.expectedPoolInfo.TokenPoolModules[i])
			var moduleAddr chainwriter.SuiAddress
			copy(moduleAddr[:], moduleBytes)
			tokenPoolInfo.TokenPoolModules = append(tokenPoolInfo.TokenPoolModules, moduleAddr)
		}
	}

	// Add token types
	tokenPoolInfo.TokenTypes = tc.expectedPoolInfo.TokenTypes

	// BCS encode the response
	bcsBytes, err := bcs.Marshal(tokenPoolInfo)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal token pool info for test %s: %v", tc.name, err))
	}

	// Convert to Sui response format
	bcsAsAny := make([]any, len(bcsBytes))
	for i, b := range bcsBytes {
		bcsAsAny[i] = float64(b)
	}

	structType := tokenAdminRegistryStructType
	results := []any{bcsAsAny, structType}

	return &suiclient.ExecutionResultType{
		ReturnValues: []suiclient.ReturnValueType{results},
	}
}

func TestGeneratePTBCommandsForTokenPools(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	tests := []struct {
		name        string
		tokenPools  []chainwriter.TokenPool
		expectedLen int
	}{
		{
			name: "single token pool",
			tokenPools: []chainwriter.TokenPool{
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
			tokenPools: []chainwriter.TokenPool{
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
			tokenPools:  []chainwriter.TokenPool{},
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ptbCommands, err := chainwriter.GeneratePTBCommandsForTokenPools(lggr, tt.tokenPools)

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
				assert.Len(t, cmd.Params, 6) // Should have 6 parameters

				// Check specific parameter names
				paramNames := make([]string, len(cmd.Params))
				for j, param := range cmd.Params {
					paramNames[j] = param.Name
				}
				assert.Contains(t, paramNames, fmt.Sprintf("ref_%d", i+1))
				assert.Contains(t, paramNames, "clock")
				assert.Contains(t, paramNames, "remote_chain_selector")
				assert.Contains(t, paramNames, "receiver_params")
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
		tokenPools []chainwriter.TokenPool
	}{
		{
			name: "single token pool",
			tokenPools: []chainwriter.TokenPool{
				{
					PackageId: "0x123",
					Index:     0,
				},
			},
		},
		{
			name: "multiple token pools",
			tokenPools: []chainwriter.TokenPool{
				{
					PackageId: "0x123",
					Index:     0,
				},
				{
					PackageId: "0x456",
					Index:     1,
				},
			},
		},
		{
			name:       "empty token pools",
			tokenPools: []chainwriter.TokenPool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			args, err := chainwriter.GenerateArgumentsForTokenPools(ccipStateRef, clockRef, lggr, tt.tokenPools)

			require.NoError(t, err)
			assert.NotNil(t, args)

			// Verify common arguments
			assert.Equal(t, ccipStateRef, args["ccip_state_ref"])
			assert.Equal(t, clockRef, args["clock_ref"])
			assert.Equal(t, 0, args["remote_chain_selector"])

			// Verify token pool specific arguments
			for i, tokenPool := range tt.tokenPools {
				poolKey := fmt.Sprintf("pool_%d", i+1)
				indexKey := fmt.Sprintf("index_%d", i+1)

				assert.Equal(t, tokenPool.PackageId, args[poolKey])
				assert.Equal(t, tokenPool.Index, args[indexKey])
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
			args, err := chainwriter.GenerateReceiverCallArguments(lggr, tt.messages, tt.previousCommandIndex, ccipObjectRef)

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

	cwConfig := GetTestChainWriterConfig()
	expander := chainwriter.NewSuiPTBExpander(lggr, mockSuiPTBClient, cwConfig)

	_, publicKey, _, err := testutils.GenerateAccountKeyPair(t, lggr)
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
					if len(receiverParts) != chainwriter.SUI_PATH_COMPONENTS_COUNT {
						continue // Skip invalid format cases
					}

					if tt.errorMessage == "mock error" {
						mockSuiPTBClient.EXPECT().
							ReadFunction(
								gomock.Any(),
								gomock.Any(),
								gomock.Any(),
								"ccip",
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

					expectedResult := &suiclient.ExecutionResultType{
						ReturnValues: []suiclient.ReturnValueType{results},
					}

					mockSuiPTBClient.EXPECT().
						ReadFunction(
							gomock.Any(),
							gomock.Any(),
							gomock.Any(),
							"ccip",
							"is_registered_receiver",
							gomock.Any(),
							[]string{"object_id", "address"},
						).
						Return(expectedResult, nil).
						Times(1)
				}
			}

			filteredMessages, err := expander.FilterRegisteredReceivers(lggr, tt.messages, publicKey)

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
						assert.Len(t, receiverParts, chainwriter.SUI_PATH_COMPONENTS_COUNT)
					}
				}
			}
		})
	}
}
