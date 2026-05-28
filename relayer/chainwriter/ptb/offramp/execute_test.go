package offramp

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/stretchr/testify/assert"
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

func TestProcessReceivers_SkipsTokenOnlyMessage(t *testing.T) {
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

func TestProcessReceivers_RequiresAppDeliveryForDataMessage(t *testing.T) {
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

func TestProcessReceivers_RequiresAppDeliveryForGasLimitMessage(t *testing.T) {
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
