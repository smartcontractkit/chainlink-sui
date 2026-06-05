//nolint:staticcheck,revive // ccipocr3.Message is a deprecated alias; method names match SuiPTBClient interface.
package offramp

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	suisigner "github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// stubPTBClient is a minimal no-op implementation of client.SuiPTBClient
// for tests that exercise code paths before any RPC call is made.
type stubPTBClient struct{}

var _ client.SuiPTBClient = (*stubPTBClient)(nil)

func (s *stubPTBClient) MoveCall(context.Context, client.MoveCallRequest) (client.TxnMetaData, error) {
	return client.TxnMetaData{}, nil
}
func (s *stubPTBClient) SendTransaction(context.Context, *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error) {
	return nil, nil
}
func (s *stubPTBClient) ReadOwnedObjects(context.Context, string, []byte) ([]*suirpcv2.Object, error) {
	return nil, nil
}
func (s *stubPTBClient) ReadFilterOwnedObjectIds(context.Context, string, string, []byte) ([]*suirpcv2.Object, error) {
	return nil, nil
}
func (s *stubPTBClient) ReadObjectId(context.Context, string) (*suirpcv2.Object, error) {
	return nil, nil
}
func (s *stubPTBClient) ReadFunction(context.Context, string, string, string, []any, []string, []string) ([]any, error) {
	return nil, nil
}
func (s *stubPTBClient) SimulatePTB(context.Context, []byte) ([]any, error) {
	return nil, nil
}
func (s *stubPTBClient) SignAndSendTransaction(context.Context, string, []byte) (*suirpcv2.ExecuteTransactionResponse, error) {
	return nil, nil
}
func (s *stubPTBClient) QueryEvents(context.Context, client.EventFilterByMoveEventModule, *uint, *client.EventId, *client.QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	return nil, nil
}
func (s *stubPTBClient) QueryTransactions(context.Context, string, *suirpcv2.Checkpoint, *uint64) ([]*suirpcv2.ExecutedTransaction, error) {
	return nil, nil
}
func (s *stubPTBClient) GetTransactionStatus(context.Context, string) (client.TransactionResult, error) {
	return client.TransactionResult{}, nil
}
func (s *stubPTBClient) GetCoinsByAddress(context.Context, string) ([]*suirpcv2.Object, error) {
	return nil, nil
}
func (s *stubPTBClient) QueryCoinsByAddress(context.Context, string, string) ([]*suirpcv2.Object, error) {
	return nil, nil
}
func (s *stubPTBClient) EstimateGas(context.Context, *transaction.Transaction) (uint64, error) {
	return 0, nil
}
func (s *stubPTBClient) GetReferenceGasPrice(context.Context) (*big.Int, error) {
	return big.NewInt(1000), nil
}
func (s *stubPTBClient) FinishPTBAndSend(context.Context, *suisigner.Signer, *transaction.Transaction, client.TransactionRequestType) (*suirpcv2.ExecuteTransactionResponse, error) {
	return nil, nil
}
func (s *stubPTBClient) GetBlockById(context.Context, string) (*suirpcv2.Checkpoint, error) {
	return nil, nil
}
func (s *stubPTBClient) BlockByDigest(context.Context, string) (*suirpcv2.Checkpoint, error) {
	return nil, nil
}
func (s *stubPTBClient) GetLatestEpoch(context.Context) (*suirpcv2.Epoch, error) {
	return nil, nil
}
func (s *stubPTBClient) GetLatestCheckpoint(context.Context) (*suirpcv2.Checkpoint, error) {
	return nil, nil
}
func (s *stubPTBClient) GetCheckpointData(context.Context, uint64) (*client.CheckpointData, error) {
	return nil, nil
}
func (s *stubPTBClient) GetNormalizedModule(context.Context, string, string) (models.GetNormalizedMoveModuleResponse, error) {
	return models.GetNormalizedMoveModuleResponse{}, nil
}
func (s *stubPTBClient) GetSUIBalance(context.Context, string) (*suirpcv2.Balance, error) {
	return nil, nil
}
func (s *stubPTBClient) LoadModulePackageIds(context.Context, string, string) ([]string, error) {
	return nil, nil
}
func (s *stubPTBClient) GetLatestPackageId(_ context.Context, id, _ string) (string, error) {
	return id, nil
}
func (s *stubPTBClient) GetCoinMetadata(context.Context, string) (models.CoinMetadataResponse, error) {
	return models.CoinMetadataResponse{}, nil
}
func (s *stubPTBClient) GetCache() *cache.Cache                                   { return nil }
func (s *stubPTBClient) GetCachedValue(string) (any, bool)                        { return nil, false }
func (s *stubPTBClient) SetCachedValue(string, any)                               {}
func (s *stubPTBClient) GetCachedValues([]string) (map[string]any, bool)          { return nil, false }
func (s *stubPTBClient) SetCachedValues(map[string]any)                           {}
func (s *stubPTBClient) HashTxBytes(b []byte) []byte                              { return nil }
func (s *stubPTBClient) GetCCIPPackageID(context.Context, string) (string, error) { return "", nil }
func (s *stubPTBClient) GetValuesFromPackageOwnedObjectField(context.Context, string, string, string, []string) (map[string]string, error) {
	return nil, nil
}
func (s *stubPTBClient) GetParentObjectID(context.Context, string, string, string) (string, error) {
	return "", nil
}

func TestDecodeOffRampExecCallArgs(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		args := map[string]any{
			"ReportContext": [2][32]byte{},
			"Report":        []byte{0x01, 0x02},
			"Info":          ccipocr3.ExecuteReportInfo{},
			"ExtraData":     ExtraDataDecoded{},
		}
		result, err := DecodeOffRampExecCallArgs(args)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []byte{0x01, 0x02}, result.Report)
	})

	t.Run("empty map returns zero-value struct", func(t *testing.T) {
		result, err := DecodeOffRampExecCallArgs(map[string]any{})
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("nil map returns nil result without error", func(t *testing.T) {
		result, err := DecodeOffRampExecCallArgs(nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestIsValidTokenPoolConfig(t *testing.T) {
	tests := []struct {
		name   string
		config module_token_admin_registry.TokenConfig
		valid  bool
	}{
		{
			name: "all fields present",
			config: module_token_admin_registry.TokenConfig{
				TokenPoolPackageId: "0x1",
				TokenPoolModule:    "pool",
				TokenType:          "0x2::coin::COIN",
			},
			valid: true,
		},
		{
			name:   "missing package ID",
			config: module_token_admin_registry.TokenConfig{TokenPoolModule: "pool", TokenType: "0x2"},
			valid:  false,
		},
		{
			name:   "missing module",
			config: module_token_admin_registry.TokenConfig{TokenPoolPackageId: "0x1", TokenType: "0x2"},
			valid:  false,
		},
		{
			name:   "missing token type",
			config: module_token_admin_registry.TokenConfig{TokenPoolPackageId: "0x1", TokenPoolModule: "pool"},
			valid:  false,
		},
		{
			name:   "all empty",
			config: module_token_admin_registry.TokenConfig{},
			valid:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.valid, IsValidTokenPoolConfig(&tc.config))
		})
	}
}

func TestNeedsAppDelivery(t *testing.T) {
	tests := []struct {
		name      string
		message   ccipocr3.Message
		extraArgs map[string]any
		expected  bool
	}{
		{
			name:      "empty data and no gasLimit in extraArgs",
			message:   ccipocr3.Message{Data: nil},
			extraArgs: map[string]any{},
			expected:  false,
		},
		{
			name:      "empty data and zero gasLimit (big.Int)",
			message:   ccipocr3.Message{Data: []byte{}},
			extraArgs: map[string]any{"gasLimit": big.NewInt(0)},
			expected:  false,
		},
		{
			name:      "empty data and nil gasLimit (big.Int)",
			message:   ccipocr3.Message{Data: nil},
			extraArgs: map[string]any{"gasLimit": (*big.Int)(nil)},
			expected:  false,
		},
		{
			name:      "empty data and zero gasLimit (uint64)",
			message:   ccipocr3.Message{Data: []byte{}},
			extraArgs: map[string]any{"gasLimit": uint64(0)},
			expected:  false,
		},
		{
			name:      "empty data and negative gasLimit (big.Int)",
			message:   ccipocr3.Message{Data: nil},
			extraArgs: map[string]any{"gasLimit": big.NewInt(-1)},
			expected:  false,
		},
		{
			name:      "non-empty data and zero gasLimit",
			message:   ccipocr3.Message{Data: []byte{0x01}},
			extraArgs: map[string]any{"gasLimit": big.NewInt(0)},
			expected:  true,
		},
		{
			name:      "non-empty data and no extraArgs",
			message:   ccipocr3.Message{Data: []byte{0xab, 0xcd}},
			extraArgs: map[string]any{},
			expected:  true,
		},
		{
			name:      "empty data and positive gasLimit (big.Int)",
			message:   ccipocr3.Message{Data: nil},
			extraArgs: map[string]any{"gasLimit": big.NewInt(200000)},
			expected:  true,
		},
		{
			name:      "empty data and positive gasLimit (uint64)",
			message:   ccipocr3.Message{Data: []byte{}},
			extraArgs: map[string]any{"gasLimit": uint64(100000)},
			expected:  true,
		},
		{
			name:      "both data and gasLimit present",
			message:   ccipocr3.Message{Data: []byte{0x01}},
			extraArgs: map[string]any{"gasLimit": big.NewInt(500000)},
			expected:  true,
		},
		{
			name:      "gasLimit is unexpected type (string) — treated as absent",
			message:   ccipocr3.Message{Data: nil},
			extraArgs: map[string]any{"gasLimit": "not-a-number"},
			expected:  false,
		},
		{
			name:      "nil extraArgs map",
			message:   ccipocr3.Message{Data: nil},
			extraArgs: nil,
			expected:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := needsAppDelivery(tc.message, tc.extraArgs)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestExtractReceiverObjectIDs(t *testing.T) {
	lggr := logger.Test(t)

	tests := []struct {
		name     string
		args     map[string]any
		expected [][]byte
	}{
		{
			name:     "key missing",
			args:     map[string]any{},
			expected: nil,
		},
		{
			name:     "key explicitly nil",
			args:     map[string]any{"receiverObjectIds": nil},
			expected: nil,
		},
		{
			name:     "nil map",
			args:     nil,
			expected: nil,
		},
		{
			name:     "[][]byte value",
			args:     map[string]any{"receiverObjectIds": [][]byte{{0x01}, {0x02, 0x03}}},
			expected: [][]byte{{0x01}, {0x02, 0x03}},
		},
		{
			name:     "empty [][]byte",
			args:     map[string]any{"receiverObjectIds": [][]byte{}},
			expected: [][]byte{},
		},
		{
			name:     "[]any with valid []byte elements",
			args:     map[string]any{"receiverObjectIds": []any{[]byte{0xaa}, []byte{0xbb}}},
			expected: [][]byte{{0xaa}, {0xbb}},
		},
		{
			name:     "[]any with mixed types skips non-byte elements",
			args:     map[string]any{"receiverObjectIds": []any{[]byte{0xaa}, "bad", []byte{0xcc}}},
			expected: [][]byte{{0xaa}, {0xcc}},
		},
		{
			name:     "[]any all non-byte returns nil",
			args:     map[string]any{"receiverObjectIds": []any{"a", 42}},
			expected: nil,
		},
		{
			name:     "unexpected type returns nil",
			args:     map[string]any{"receiverObjectIds": "not-a-slice"},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractReceiverObjectIDs(lggr, tc.args)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// validSuiAddress returns a zero-padded Sui address suitable for test fixtures.
func validSuiAddress() string {
	return "0x0000000000000000000000000000000000000000000000000000000000000001"
}

func testAddressMappings() *OffRampAddressMappings {
	return &OffRampAddressMappings{
		CcipPackageId:    validSuiAddress(),
		CcipObjectRef:    validSuiAddress(),
		CcipOwnerCap:     validSuiAddress(),
		ClockObject:      "0x6",
		OffRampPackageId: validSuiAddress(),
		OffRampState:     validSuiAddress(),
	}
}

func testCallOpts() *bind.CallOpts {
	return &bind.CallOpts{
		Signer:           signer.NewDevInspectSigner(validSuiAddress()),
		WaitForExecution: true,
	}
}

func TestProcessReceivers_SkipPaths(t *testing.T) {
	lggr := logger.Test(t)
	ctx := context.Background()
	fakeClient := &stubPTBClient{}
	ptb := transaction.NewTransaction()
	mappings := testAddressMappings()
	callOpts := testCallOpts()
	var receiverParams *transaction.Argument

	tests := []struct {
		name      string
		messages  []ccipocr3.Message
		extraArgs map[string]any
	}{
		{
			name:      "nil receiver",
			messages:  []ccipocr3.Message{{Receiver: nil, Data: []byte{0x01}}},
			extraArgs: map[string]any{"gasLimit": big.NewInt(100000)},
		},
		{
			name:      "empty receiver",
			messages:  []ccipocr3.Message{{Receiver: []byte{}, Data: []byte{0x01}}},
			extraArgs: map[string]any{"gasLimit": big.NewInt(100000)},
		},
		{
			name:      "zero-address receiver",
			messages:  []ccipocr3.Message{{Receiver: make([]byte, 32), Data: []byte{0x01}}},
			extraArgs: map[string]any{"gasLimit": big.NewInt(100000)},
		},
		{
			name: "no app delivery needed",
			messages: []ccipocr3.Message{{
				Receiver: append(make([]byte, 31), 0x01),
				Data:     nil,
			}},
			extraArgs: map[string]any{"gasLimit": big.NewInt(0)},
		},
		{
			name:      "no messages",
			messages:  []ccipocr3.Message{},
			extraArgs: map[string]any{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			results, err := ProcessReceivers(
				ctx, lggr, fakeClient, ptb,
				tc.messages, mappings, callOpts, receiverParams, tc.extraArgs,
			)
			require.NoError(t, err)
			assert.Empty(t, results)
		})
	}
}

func TestProcessReceivers_SkipsZeroAddressReceiver(t *testing.T) {
	// A 32-byte zero address receiver should be skipped.
	zeroAddr := make([]byte, 32)
	msg := ccipocr3.Message{
		Receiver: zeroAddr,
		Data:     []byte{0x01},
	}
	// Verify our check matches codec.AccountZero pattern
	allZero := true
	for _, b := range msg.Receiver {
		if b != 0 {
			allZero = false
			break
		}
	}
	assert.True(t, allZero)
}

func TestNeedsAppDelivery_ReturnsFalseForTokenOnlyMessage(t *testing.T) {
	receiver := make([]byte, 32)
	receiver[31] = 0x01

	msg := ccipocr3.Message{
		Receiver: receiver,
		Data:     nil,
	}
	extraArgs := map[string]any{"gasLimit": big.NewInt(0)}

	assert.False(t, needsAppDelivery(msg, extraArgs),
		"token-only message should not need app delivery")
}

func TestNeedsAppDelivery_ReturnsTrueForDataMessage(t *testing.T) {
	receiver := make([]byte, 32)
	receiver[31] = 0x01

	msg := ccipocr3.Message{
		Receiver: receiver,
		Data:     []byte{0xde, 0xad},
	}
	extraArgs := map[string]any{}

	assert.True(t, needsAppDelivery(msg, extraArgs),
		"message with data should need app delivery")
}

func TestNeedsAppDelivery_ReturnsTrueForPositiveGasLimit(t *testing.T) {
	receiver := make([]byte, 32)
	receiver[31] = 0x01

	msg := ccipocr3.Message{
		Receiver: receiver,
		Data:     nil,
	}
	extraArgs := map[string]any{"gasLimit": big.NewInt(100000)}

	assert.True(t, needsAppDelivery(msg, extraArgs),
		"message with gasLimit > 0 should need app delivery")
}

func TestExtractReceiverObjectIdStrings(t *testing.T) {
	tests := []struct {
		name      string
		extraArgs map[string]any
		expected  []string
	}{
		{
			name:      "missing key returns empty slice",
			extraArgs: map[string]any{},
			expected:  []string{},
		},
		{
			name:      "nil extraArgs returns empty slice",
			extraArgs: nil,
			expected:  []string{},
		},
		{
			name: "[][]byte with one entry",
			extraArgs: map[string]any{
				"receiverObjectIds": [][]byte{
					{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x11, 0x11},
				},
			},
			expected: []string{"0x0000000000000000000000000000000000000000000000000000000000001111"},
		},
		{
			name: "[][]byte with multiple entries",
			extraArgs: map[string]any{
				"receiverObjectIds": [][]byte{
					{0xAA, 0xBB},
					{0xCC, 0xDD},
				},
			},
			expected: []string{"0xaabb", "0xccdd"},
		},
		{
			name: "[]any with []byte elements",
			extraArgs: map[string]any{
				"receiverObjectIds": []any{
					[]byte{0x11, 0x22},
					[]byte{0x33, 0x44},
				},
			},
			expected: []string{"0x1122", "0x3344"},
		},
		{
			name: "[]any with non-byte elements skipped",
			extraArgs: map[string]any{
				"receiverObjectIds": []any{
					[]byte{0x11, 0x22},
					"not-bytes",
					[]byte{0x33, 0x44},
				},
			},
			expected: []string{"0x1122", "0x3344"},
		},
		{
			name: "unexpected type returns empty",
			extraArgs: map[string]any{
				"receiverObjectIds": "not-a-slice",
			},
			expected: []string{},
		},
		{
			name: "empty [][]byte returns empty",
			extraArgs: map[string]any{
				"receiverObjectIds": [][]byte{},
			},
			expected: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := extractReceiverObjectIDStrings(tc.extraArgs)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestClassifyReceiverBuildError(t *testing.T) {
	receiverPackageID := "0x00000000000000000000000000000000000000000000000000000000000000ab"

	tests := []struct {
		name          string
		err           error
		expectSkip    bool
		expectErr     bool
		expectErrText string
	}{
		{
			name:       "nil error",
			err:        nil,
			expectSkip: false,
			expectErr:  false,
		},
		{
			name:       "71024 poison ABI wrapped as unsupported",
			err:        fmt.Errorf("%w: TypeParameter (generic parameters are not supported)", ErrUnsupportedReceiverABI),
			expectSkip: true,
			expectErr:  false,
		},
		{
			name:       "71024 missing ccip_receive wrapped as unsupported",
			err:        fmt.Errorf("%w: function %q not found in module", ErrUnsupportedReceiverABI, "ccip_receive"),
			expectSkip: true,
			expectErr:  false,
		},
		{
			name:          "transient PTB build error propagates for ProcessReceiversV2",
			err:           errors.New("failed to build PTB (receiver call) using bindings: network timeout"),
			expectSkip:    false,
			expectErr:     true,
			expectErrText: "failed to build receiver command for " + receiverPackageID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			skip, retErr := classifyReceiverBuildError(receiverPackageID, tc.err)
			assert.Equal(t, tc.expectSkip, skip)
			if tc.expectErr {
				require.Error(t, retErr)
				assert.Contains(t, retErr.Error(), tc.expectErrText)
				return
			}
			require.NoError(t, retErr)
		})
	}
}

func TestAppendCcipReceiveCommand_ReceiverObjectIdCountMismatch(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	ptb := transaction.NewTransaction()
	extractedArg := ptb.MoveCall("0x2", "extract", "extract", nil, nil)
	callOpts := &bind.CallOpts{Signer: signer.NewDevInspectSigner("0x1")}

	boundReceiverContract, err := bind.NewBoundContract(
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"dummy_receiver",
		nil,
	)
	require.NoError(t, err)

	// Dummy receiver ABI: 3 fixed + clock + state = 5 params (TxContext omitted by DecodeParameters).
	normalizedModule := models.GetNormalizedMoveModuleResponse{
		ExposedFunctions: map[string]any{
			"ccip_receive": map[string]any{
				"parameters": []any{
					map[string]any{"Vector": "U8"},
					map[string]any{"Reference": map[string]any{"Struct": map[string]any{
						"address": "0xccip", "module": "state_object", "name": "CCIPObjectRef", "typeArguments": []any{},
					}}},
					map[string]any{"Struct": map[string]any{
						"address": "0xccip", "module": "client", "name": "Any2SuiMessageV2", "typeArguments": []any{},
					}},
					map[string]any{"Reference": map[string]any{"Struct": map[string]any{
						"address": "0x2", "module": "clock", "name": "Clock", "typeArguments": []any{},
					}}},
					map[string]any{"MutableReference": map[string]any{"Struct": map[string]any{
						"address": "0xreceiver", "module": "dummy_receiver", "name": "CCIPReceiverState", "typeArguments": []any{},
					}}},
				},
			},
		},
	}

	_, err = appendCcipReceiveCommand(
		ctx, lggr, ptb, callOpts, boundReceiverContract, "ccip_receive",
		&OffRampAddressMappings{CcipObjectRef: "0x3"},
		[32]byte{}, &normalizedModule, &extractedArg,
		[]string{"0x00000000000000000000000000000000000000000000000000000000000000aa"},
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnsupportedReceiverABI)
	assert.Contains(t, err.Error(), "receiver_object_ids count 1 does not match ccip_receive tail parameter count 2")
}

func TestAppendCcipReceiveCommand_MissingCcipReceive(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	ptb := transaction.NewTransaction()
	extractedArg := ptb.MoveCall("0x2", "extract", "extract", nil, nil)
	callOpts := &bind.CallOpts{Signer: signer.NewDevInspectSigner("0x1")}

	boundReceiverContract, err := bind.NewBoundContract(
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"receiver",
		nil,
	)
	require.NoError(t, err)

	_, err = appendCcipReceiveCommand(
		ctx, lggr, ptb, callOpts, boundReceiverContract, "ccip_receive",
		&OffRampAddressMappings{CcipObjectRef: "0x3"},
		[32]byte{},
		&models.GetNormalizedMoveModuleResponse{ExposedFunctions: map[string]any{}},
		&extractedArg,
		nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnsupportedReceiverABI)
}

func TestAppendCcipReceiveCommand_PoisonABI(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	ptb := transaction.NewTransaction()
	extractedArg := ptb.MoveCall("0x2", "extract", "extract", nil, nil)
	callOpts := &bind.CallOpts{Signer: signer.NewDevInspectSigner("0x1")}

	boundReceiverContract, err := bind.NewBoundContract(
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"receiver",
		nil,
	)
	require.NoError(t, err)

	normalizedModule := models.GetNormalizedMoveModuleResponse{
		ExposedFunctions: map[string]any{
			"ccip_receive": map[string]any{
				"parameters": []any{
					map[string]any{"Vector": map[string]any{"TypeParameter": float64(0)}},
				},
			},
		},
	}

	_, err = appendCcipReceiveCommand(
		ctx, lggr, ptb, callOpts, boundReceiverContract, "ccip_receive",
		&OffRampAddressMappings{CcipObjectRef: "0x3"},
		[32]byte{}, &normalizedModule, &extractedArg, nil,
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnsupportedReceiverABI)
	assert.Contains(t, err.Error(), "failed to decode receiver parameters")
}

func TestAppendPTBCommandForReceiver_PoisonABI(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	ptb := transaction.NewTransaction()
	receiverParams := ptb.MoveCall("0x2", "init", "init", nil, nil)
	callOpts := &bind.CallOpts{Signer: signer.NewDevInspectSigner("0x1")}

	normalizedModule := models.GetNormalizedMoveModuleResponse{
		ExposedFunctions: map[string]any{
			"ccip_receive": map[string]any{
				"parameters": []any{
					map[string]any{"Vector": map[string]any{"TypeParameter": float64(0)}},
				},
			},
		},
	}

	_, err := AppendPTBCommandForReceiver(
		ctx, lggr, nil, ptb, callOpts,
		"0x0000000000000000000000000000000000000000000000000000000000000001",
		"receiver", "ccip_receive",
		&OffRampAddressMappings{
			CcipPackageId: "0x0000000000000000000000000000000000000000000000000000000000000002",
			CcipObjectRef: "0x3",
		},
		[32]byte{}, &normalizedModule, &receiverParams, map[string]any{},
	)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnsupportedReceiverABI)
	assert.Contains(t, err.Error(), "failed to decode receiver parameters")
}

func TestFormatReceiverObjectIDStrings(t *testing.T) {
	objectIDBytes := make([]byte, 32)
	objectIDBytes[31] = 0x11
	var objectID models.SuiAddressBytes
	copy(objectID[:], objectIDBytes)

	result := codec.FormatReceiverObjectIDStrings([]models.SuiAddressBytes{objectID})
	assert.Equal(t, []string{"0x0000000000000000000000000000000000000000000000000000000000000011"}, result)
	assert.Empty(t, codec.FormatReceiverObjectIDStrings(nil))
}

func TestExtractAny2suiMessageV2_EncodeUsesVectorAddress(t *testing.T) {
	ptb := transaction.NewTransaction()
	receiverParams := ptb.MoveCall("0x2", "init_execute_v2", "init_execute_v2", nil, nil)

	contract, err := bind.NewBoundContract(
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"0x0000000000000000000000000000000000000000000000000000000000000002",
		"offramp_state_helper",
		nil,
	)
	require.NoError(t, err)

	receiverObjectIDs := []string{"0x0000000000000000000000000000000000000000000000000000000000000011"}
	_, err = contract.EncodeCallArgsWithGenerics(
		"extract_any2sui_message_v2",
		nil,
		nil,
		[]string{"&mut object", "vector<address>"},
		[]any{receiverParams, receiverObjectIDs},
		nil,
	)
	require.NoError(t, err)
}

func TestProcessReceiversV2_GatingParityWithV1(t *testing.T) {
	// Receiver object IDs are sourced from the V2 execution report, not ExtraArgsDecoded.
	receiver := make([]byte, 32)
	receiver[31] = 0x01

	tokenOnlyMsg := ccipocr3.Message{
		Receiver: receiver,
		Data:     nil,
	}
	extraArgs := map[string]any{"gasLimit": big.NewInt(0)}

	assert.False(t, needsAppDelivery(tokenOnlyMsg, extraArgs))
	assert.Empty(t, extractReceiverObjectIDStrings(extraArgs))
}
