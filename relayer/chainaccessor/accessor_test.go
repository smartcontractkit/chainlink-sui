package chainaccessor

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

func TestSuiAddressRoundTrip(t *testing.T) {
	t.Parallel()

	const addr = "0x2"
	b, err := suiAddressToBytes(addr)
	require.NoError(t, err)
	require.Len(t, b, 32)
	assert.Equal(t, byte(0x02), b[31])

	// FromBytes produces the zero-padded canonical form.
	got := suiAddressFromBytes(b)
	assert.Len(t, got, 2+64)
	roundTripped, err := suiAddressToBytes(got)
	require.NoError(t, err)
	assert.Equal(t, b, roundTripped)
}

func TestAsUint64(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   any
		want uint64
	}{
		{"string", "12345", 12345},
		{"float64", float64(42), 42},
		{"uint64", uint64(7), 7},
		{"big string", "18446744073709551615", ^uint64(0)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got, err := asUint64(c.in)
			require.NoError(t, err)
			assert.Equal(t, c.want, got)
		})
	}

	_, err := asUint64(struct{}{})
	assert.Error(t, err)
}

func TestAsBytes(t *testing.T) {
	t.Parallel()

	hexBytes, err := asBytes("0x0a0b0c")
	require.NoError(t, err)
	assert.Equal(t, []byte{0x0a, 0x0b, 0x0c}, hexBytes)

	sliceBytes, err := asBytes([]any{float64(1), float64(2), float64(255)})
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 255}, sliceBytes)

	empty, err := asBytes("")
	require.NoError(t, err)
	assert.Nil(t, empty)
}

func TestToBytes32(t *testing.T) {
	t.Parallel()

	short := toBytes32([]byte{0xaa, 0xbb})
	assert.Equal(t, byte(0xaa), short[30])
	assert.Equal(t, byte(0xbb), short[31])

	long := make([]byte, 40)
	long[39] = 0xff
	got := toBytes32(long)
	assert.Equal(t, byte(0xff), got[31])
}

func TestBindingsCache(t *testing.T) {
	t.Parallel()

	c := newBindingsCache()

	_, err := c.getPackageAddress(ContractNameOnRamp)
	assert.ErrorIs(t, err, ErrNotBound)

	c.setPackageAddress(ContractNameOnRamp, "0xabc")
	addr, err := c.getPackageAddress(ContractNameOnRamp)
	require.NoError(t, err)
	assert.Equal(t, "0xabc", addr)

	_, err = c.getStateObjectID("onramp")
	assert.ErrorIs(t, err, ErrNotBound)
	c.setStateObjectID("onramp", "0xstate")
	state, err := c.getStateObjectID("onramp")
	require.NoError(t, err)
	assert.Equal(t, "0xstate", state)

	_, err = c.getCCIPObjectRefID()
	assert.ErrorIs(t, err, ErrNotBound)
	c.setCCIPObjectRefID("0xref")
	ref, err := c.getCCIPObjectRefID()
	require.NoError(t, err)
	assert.Equal(t, "0xref", ref)
}

func TestDecodeCCIPMessageSent(t *testing.T) {
	t.Parallel()

	data := map[string]any{
		"dest_chain_selector": "20",
		"sequence_number":     "5",
		"message": map[string]any{
			"header": map[string]any{
				"message_id":            "0x" + "ab",
				"source_chain_selector": "10",
				"dest_chain_selector":   "20",
				"sequence_number":       "5",
				"nonce":                 "3",
			},
			"sender":           "0x01",
			"receiver":         "0x02",
			"data":             "0x0304",
			"fee_token_amount": "100",
			"fee_value_juels":  "200",
		},
	}

	msg, ok, err := decodeCCIPMessageSent(data)
	require.NoError(t, err)
	require.True(t, ok)

	assert.Equal(t, ccipocr3.ChainSelector(10), msg.Header.SourceChainSelector)
	assert.Equal(t, ccipocr3.ChainSelector(20), msg.Header.DestChainSelector)
	assert.Equal(t, ccipocr3.SeqNum(5), msg.Header.SequenceNumber)
	assert.Equal(t, uint64(3), msg.Header.Nonce)
	assert.Equal(t, byte(0xab), msg.Header.MessageID[31])
	assert.Equal(t, []byte{0x03, 0x04}, []byte(msg.Data))
	assert.Equal(t, big.NewInt(100), msg.FeeTokenAmount.Int)
	assert.Equal(t, big.NewInt(200), msg.FeeValueJuels.Int)
}

func TestDecodeExecutionStateChanged(t *testing.T) {
	t.Parallel()

	src, seq, err := decodeExecutionStateChanged(map[string]any{
		"source_chain_selector": "10",
		"sequence_number":       "42",
		"state":                 float64(2),
	})
	require.NoError(t, err)
	assert.Equal(t, ccipocr3.ChainSelector(10), src)
	assert.Equal(t, ccipocr3.SeqNum(42), seq)
}

func TestDecodeTimestampedPrice(t *testing.T) {
	t.Parallel()

	price, err := decodeTimestampedPrice([]any{map[string]any{
		"value":     "123456789",
		"timestamp": "1700000000",
	}})
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(123456789), price.Value)
	assert.Equal(t, uint32(1700000000), price.Timestamp)
}

func TestDecodeFeeQuoterDestChainConfig(t *testing.T) {
	t.Parallel()

	m := map[string]any{
		"is_enabled":                            true,
		"enforce_out_of_order":                  false,
		"max_number_of_tokens_per_msg":          float64(5),
		"max_data_bytes":                        float64(30000),
		"max_per_msg_gas_limit":                 float64(2000000),
		"dest_gas_overhead":                     float64(100000),
		"dest_gas_per_payload_byte_base":        float64(16),
		"dest_gas_per_payload_byte_high":        float64(40),
		"dest_gas_per_payload_byte_threshold":   float64(3000),
		"dest_data_availability_overhead_gas":   float64(700),
		"dest_gas_per_data_availability_byte":   float64(16),
		"dest_data_availability_multiplier_bps": float64(1),
		"default_token_fee_usd_cents":           float64(50),
		"default_token_dest_gas_overhead":       float64(90000),
		"default_tx_gas_limit":                  float64(200000),
		"gas_multiplier_wei_per_eth":            "1100000000000000000",
		"network_fee_usd_cents":                 float64(10),
		"gas_price_staleness_threshold":         float64(90000),
	}

	cfg, err := decodeFeeQuoterDestChainConfig(m)
	require.NoError(t, err)
	assert.True(t, cfg.IsEnabled)
	assert.False(t, cfg.EnforceOutOfOrder)
	assert.Equal(t, uint16(5), cfg.MaxNumberOfTokensPerMsg)
	assert.Equal(t, uint32(2000000), cfg.MaxPerMsgGasLimit)
	assert.Equal(t, uint64(1100000000000000000), cfg.GasMultiplierWeiPerEth)
}
