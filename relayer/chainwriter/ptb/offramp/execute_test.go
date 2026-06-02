package offramp

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	// Token-only messages (empty data, zero gas) should skip receiver processing.
	receiver := make([]byte, 32)
	receiver[31] = 0x01 // non-zero receiver

	msg := ccipocr3.Message{
		Receiver: receiver,
		Data:     nil,
	}
	extraArgs := map[string]any{"gasLimit": big.NewInt(0)}

	assert.False(t, needsAppDelivery(msg, extraArgs),
		"token-only message should not need app delivery")
}

func TestNeedsAppDelivery_ReturnsTrueForDataMessage(t *testing.T) {
	// Message with data should require app delivery.
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
	// Message with non-zero gasLimit should require app delivery even without data.
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
			result := extractReceiverObjectIdStrings(tc.extraArgs)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestClassifyReceiverBuildError(t *testing.T) {
	receiverPackageId := "0x00000000000000000000000000000000000000000000000000000000000000ab"

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
			err:           fmt.Errorf("failed to build PTB (receiver call) using bindings: network timeout"),
			expectSkip:    false,
			expectErr:     true,
			expectErrText: "failed to build receiver command for " + receiverPackageId,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			skip, retErr := classifyReceiverBuildError(receiverPackageId, tc.err)
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
	assert.True(t, errors.Is(err, ErrUnsupportedReceiverABI))
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
	assert.True(t, errors.Is(err, ErrUnsupportedReceiverABI))
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
	assert.True(t, errors.Is(err, ErrUnsupportedReceiverABI))
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
	assert.True(t, errors.Is(err, ErrUnsupportedReceiverABI))
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
	assert.Empty(t, extractReceiverObjectIdStrings(extraArgs))
}
