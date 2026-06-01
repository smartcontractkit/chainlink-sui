package codec_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestDeserializeExecutionReportV2_RoundtripWithTestutilsSerializer(t *testing.T) {
	objectID := leftPad32(t, "aabbcc")
	tokenReceiver := leftPad32(t, "5678")
	receiverPackage := []byte("0000000000000000000000000000000000000000000000000000000000001234")

	report := testutils.ExecutionReportV2{
		SourceChainSelector: 5009297550715157269,
		Message: testutils.Any2SuiRampMessageV2{
			Header: testutils.NewRampMessageHeader(
				leftPad32(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
				5009297550715157269,
				2000,
				42,
				0,
			),
			Sender:            []byte("8765432109fedcba87654321"),
			Data:              []byte("test payload"),
			Receiver:          receiverPackage,
			GasLimit:          big.NewInt(200000),
			TokenReceiver:     tokenReceiver,
			ReceiverObjectIds: [][]byte{objectID},
			TokenAmounts: []testutils.Any2SuiTokenTransfer{{
				SourcePoolAddress: []byte("source-pool"),
				DestTokenAddress:  leftPad32(t, "1234"),
				DestGasAmount:     10000,
				ExtraData:         []byte{0x01},
				Amount:            big.NewInt(1000),
			}},
		},
		OffchainTokenData: [][]byte{{0xab, 0xcd}},
		Proofs:            [][]byte{leftPad32(t, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")},
	}

	encoded, err := testutils.SerializeExecutionReportV2(report)
	require.NoError(t, err)

	decoded, err := codec.DeserializeExecutionReportV2(encoded)
	require.NoError(t, err)

	assert.Equal(t, report.SourceChainSelector, decoded.SourceChainSelector)
	assert.Equal(t, report.Message.Header.SequenceNumber, decoded.Message.Header.SequenceNumber)
	assert.Equal(t, report.Message.Sender, decoded.Message.Sender)
	assert.Equal(t, report.Message.Data, decoded.Message.Data)
	assert.Equal(t, models.SuiAddress("0000000000000000000000000000000000000000000000000000000000001234"), decoded.Message.Receiver)
	assert.Equal(t, 0, report.Message.GasLimit.Cmp(decoded.Message.GasLimit))
	assert.Equal(t, models.SuiAddressBytes(tokenReceiver), decoded.Message.TokenReceiver)
	require.Len(t, decoded.Message.ReceiverObjectIds, 1)
	assert.Equal(t, models.SuiAddressBytes(objectID), decoded.Message.ReceiverObjectIds[0])
	require.Len(t, decoded.Message.TokenAmounts, 1)
	assert.Equal(t, uint32(10000), decoded.Message.TokenAmounts[0].DestGasAmount)
	assert.Equal(t, report.OffchainTokenData, decoded.OffchainTokenData)
	assert.Equal(t, report.Proofs, decoded.Proofs)
}

func TestFormatReceiverObjectIDStrings(t *testing.T) {
	objectIDBytes := leftPad32(t, "aabbcc")
	var objectID models.SuiAddressBytes
	copy(objectID[:], objectIDBytes)

	assert.Equal(t,
		[]string{"0x0000000000000000000000000000000000000000000000000000000000aabbcc"},
		codec.FormatReceiverObjectIDStrings([]models.SuiAddressBytes{objectID}),
	)
}

func TestDeserializeExecutionReportV2_EmptyReceiverObjectIds(t *testing.T) {
	tokenReceiver := leftPad32(t, "5678")
	receiverPackage := []byte("0000000000000000000000000000000000000000000000000000000000001234")

	report := testutils.ExecutionReportV2{
		SourceChainSelector: 1000,
		Message: testutils.Any2SuiRampMessageV2{
			Header: testutils.NewRampMessageHeader(
				leftPad32(t, "1111111111111111111111111111111111111111111111111111111111111111"),
				1000,
				2000,
				1,
				0,
			),
			Sender:            []byte("sender"),
			Data:              []byte{},
			Receiver:          receiverPackage,
			GasLimit:          big.NewInt(0),
			TokenReceiver:     tokenReceiver,
			ReceiverObjectIds: [][]byte{},
			TokenAmounts:      []testutils.Any2SuiTokenTransfer{},
		},
		OffchainTokenData: [][]byte{},
		Proofs:            [][]byte{},
	}

	encoded, err := testutils.SerializeExecutionReportV2(report)
	require.NoError(t, err)

	decoded, err := codec.DeserializeExecutionReportV2(encoded)
	require.NoError(t, err)
	assert.Empty(t, decoded.Message.ReceiverObjectIds)
}

func TestDeserializeExecutionReportV2_LongerThanV1ForSameMessage(t *testing.T) {
	objectID := leftPad32(t, "aabbcc")
	tokenReceiver := leftPad32(t, "5678")
	receiverPackage := []byte("0000000000000000000000000000000000000000000000000000000000001234")
	messageID := leftPad32(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")

	header := testutils.NewRampMessageHeader(messageID, 1000, 2000, 1, 0)
	tokenAmounts := []testutils.Any2SuiTokenTransfer{{
		SourcePoolAddress: []byte("pool"),
		DestTokenAddress:  leftPad32(t, "1234"),
		DestGasAmount:     100,
		ExtraData:         []byte{0x01},
		Amount:            big.NewInt(50),
	}}

	v1Report := testutils.ExecutionReport{
		SourceChainSelector: 1000,
		Message: testutils.NewAny2SuiRampMessage(
			header,
			[]byte("sender"),
			[]byte("data"),
			receiverPackage,
			big.NewInt(500000),
			tokenReceiver,
			tokenAmounts,
		),
		OffchainTokenData: [][]byte{{0x01}},
		Proofs:            [][]byte{leftPad32(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")},
	}

	v2Report := testutils.ExecutionReportV2{
		SourceChainSelector: 1000,
		Message: testutils.Any2SuiRampMessageV2{
			Header:            header,
			Sender:            []byte("sender"),
			Data:              []byte("data"),
			Receiver:          receiverPackage,
			GasLimit:          big.NewInt(500000),
			TokenReceiver:     tokenReceiver,
			ReceiverObjectIds: [][]byte{objectID},
			TokenAmounts:      tokenAmounts,
		},
		OffchainTokenData: [][]byte{{0x01}},
		Proofs:            [][]byte{leftPad32(t, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")},
	}

	v1Bytes, err := testutils.SerializeExecutionReport(v1Report)
	require.NoError(t, err)

	v2Bytes, err := testutils.SerializeExecutionReportV2(v2Report)
	require.NoError(t, err)

	assert.Greater(t, len(v2Bytes), len(v1Bytes))

	_, err = codec.DeserializeExecutionReport(v2Bytes)
	require.Error(t, err)
}

func leftPad32(t *testing.T, hexSuffix string) []byte {
	t.Helper()
	b, err := hex.DecodeString(hexSuffix)
	require.NoError(t, err)
	require.LessOrEqual(t, len(b), 32)
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}
