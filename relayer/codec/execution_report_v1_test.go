package codec_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestCustomReportDeserializer(t *testing.T) {
	// Values preserved from the legacy pre-token_receiver hex fixture; re-encoded with
	// testutils so the blob matches the current V1 layout (includes token_receiver).
	messageID, err := hex.DecodeString("8869e580deb6dbc08e84fb41431d41d04f8849ed00be4a070dca7c34e2f78ecd")
	require.NoError(t, err)
	require.Len(t, messageID, 32)

	sender, err := hex.DecodeString("e30b40bfb1baeed9e4c62f145be85eb3d19ae932")
	require.NoError(t, err)

	receiver := []byte("4010af5717948371a0b649a59530f8e80e0e1247e015f05f1f3e09c715288dd0")

	destTokenAddress, err := hex.DecodeString("a1b6cf2e878987deb2624f9a122297abf6332d45b48c4df6fc3ea705f810980f")
	require.NoError(t, err)
	require.Len(t, destTokenAddress, 32)

	sourcePoolAddress, err := hex.DecodeString("bd10ffa3815c010d5cf7d38815a0eaabc959eb84")
	require.NoError(t, err)

	extraData, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000012")
	require.NoError(t, err)

	input := testutils.ExecutionReport{
		SourceChainSelector: 16015286601757825753,
		Message: testutils.NewAny2SuiRampMessage(
			testutils.NewRampMessageHeader(messageID, 16015286601757825753, 743186221051783445, 3, 0),
			sender,
			[]byte("I am a test ccip message"),
			receiver,
			big.NewInt(1000000),
			make([]byte, 32),
			[]testutils.Any2SuiTokenTransfer{
				testutils.NewAny2SuiTokenTransfer(
					sourcePoolAddress,
					destTokenAddress,
					100000,
					extraData,
					big.NewInt(10000000000000000),
				),
			},
		),
		OffchainTokenData: [][]byte{{}},
		Proofs:            [][]byte{},
	}

	encoded, err := testutils.SerializeExecutionReport(input)
	require.NoError(t, err)

	report, err := codec.DeserializeExecutionReport(encoded)
	require.NoError(t, err)

	t.Run("Verify integer values", func(t *testing.T) {
		require.Equal(t, uint64(16015286601757825753), report.SourceChainSelector)
		require.Equal(t, uint64(743186221051783445), report.Message.Header.DestChainSelector)
		require.Equal(t, uint64(3), report.Message.Header.SequenceNumber)
		require.Equal(t, uint64(0), report.Message.Header.Nonce)

		expectedGasLimit := big.NewInt(1000000)
		require.Equal(t, 0, expectedGasLimit.Cmp(report.Message.GasLimit))
	})

	t.Run("Verify token transfers", func(t *testing.T) {
		require.Len(t, report.Message.TokenAmounts, 1)
		tokenTransfer := report.Message.TokenAmounts[0]

		expectedAmount := big.NewInt(10000000000000000)
		require.Equal(t, 0, expectedAmount.Cmp(tokenTransfer.Amount))
		require.Equal(t, uint32(100000), tokenTransfer.DestGasAmount)
	})

	t.Run("Verify addresses", func(t *testing.T) {
		expectedReceiver := models.SuiAddress("4010af5717948371a0b649a59530f8e80e0e1247e015f05f1f3e09c715288dd0")
		require.Equal(t, expectedReceiver, report.Message.Receiver)

		expectedDestAddress := models.SuiAddress("a1b6cf2e878987deb2624f9a122297abf6332d45b48c4df6fc3ea705f810980f")
		require.Equal(t, expectedDestAddress, report.Message.TokenAmounts[0].DestTokenAddress)
	})

	t.Run("Verify message data", func(t *testing.T) {
		require.Equal(t, "I am a test ccip message", string(report.Message.Data))
	})
}

func TestCustomReportDeserializer_LegacyHexWithoutTokenReceiverRejected(t *testing.T) {
	legacyReportHex := "d91ad9c94fba41de8869e580deb6dbc08e84fb41431d41d04f8849ed00be4a070dca7c34e2f78ecdd91ad9c94fba41de15a9c133ee53500a0300000000000000000000000000000014e30b40bfb1baeed9e4c62f145be85eb3d19ae932184920616d206120746573742063636970206d6573736167654010af5717948371a0b649a59530f8e80e0e1247e015f05f1f3e09c715288dd040420f00000000000000000000000000000000000000000000000000000000000114bd10ffa3815c010d5cf7d38815a0eaabc959eb84a1b6cf2e878987deb2624f9a122297abf6332d45b48c4df6fc3ea705f810980fa08601002000000000000000000000000000000000000000000000000000000000000000120000c16ff2862300000000000000000000000000000000000000000000000000010000"
	data, err := hex.DecodeString(legacyReportHex)
	require.NoError(t, err)

	_, err = codec.DeserializeExecutionReport(data)
	require.Error(t, err)
}
