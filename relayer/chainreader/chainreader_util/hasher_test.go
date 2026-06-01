package chainreaderutil

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestMessageHasherV1_MetadataHash(t *testing.T) {
	t.Run("should match expected metadata hash", func(t *testing.T) {
		sourceChainSelector := uint64(123456789)
		destChainSelector := uint64(987654321)
		onRamp := []byte("source-onramp-address")

		metadataHash, err := computeMetadataHash(sourceChainSelector, destChainSelector, onRamp)
		require.NoError(t, err)

		expectedMetadataHash := "b62ec658417caa5bcc6ff1d8c45f8b1cb52e1b0ed71603a04b250b107ed836d9"
		actualMetadataHash := hex.EncodeToString(metadataHash[:])

		assert.Equal(t, expectedMetadataHash, actualMetadataHash)
	})

	t.Run("should produce different hash when source chain selector changes", func(t *testing.T) {
		sourceChainSelector := uint64(123456789)
		destChainSelector := uint64(987654321)
		onRamp := []byte("source-onramp-address")

		metadataHash, err := computeMetadataHash(sourceChainSelector, destChainSelector, onRamp)
		require.NoError(t, err)

		metadataHashDifferentSource, err := computeMetadataHash(sourceChainSelector+1, destChainSelector, onRamp)
		require.NoError(t, err)

		assert.NotEqual(t, metadataHash, metadataHashDifferentSource)

		expectedMetadataHashDifferentSource := "89da72ab93f7bd546d60b58a1e1b5f628fd456fe163614ff1e31a2413ca1b55a"
		actualMetadataHashDifferentSource := hex.EncodeToString(metadataHashDifferentSource[:])

		assert.Equal(t, expectedMetadataHashDifferentSource, actualMetadataHashDifferentSource)
	})

	t.Run("should produce different hash when destination chain selector changes", func(t *testing.T) {
		sourceChainSelector := uint64(123456789)
		destChainSelector := uint64(987654321)
		onRamp := []byte("source-onramp-address")

		metadataHash, err := computeMetadataHash(sourceChainSelector, destChainSelector, onRamp)
		require.NoError(t, err)

		metadataHashDifferentDest, err := computeMetadataHash(sourceChainSelector, destChainSelector+1, onRamp)
		require.NoError(t, err)

		assert.NotEqual(t, metadataHash, metadataHashDifferentDest)
	})

	t.Run("should produce different hash when on_ramp changes", func(t *testing.T) {
		sourceChainSelector := uint64(123456789)
		destChainSelector := uint64(987654321)
		onRamp := []byte("source-onramp-address")

		metadataHash, err := computeMetadataHash(sourceChainSelector, destChainSelector, onRamp)
		require.NoError(t, err)

		differentOnRamp := []byte("different-onramp-address")
		metadataHashDifferentOnRamp, err := computeMetadataHash(sourceChainSelector, destChainSelector, differentOnRamp)
		require.NoError(t, err)

		assert.NotEqual(t, metadataHash, metadataHashDifferentOnRamp)
	})
}

func TestMessageHasherV1_MessageHash(t *testing.T) {
	messageIDHex := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	sourceChainSelector := uint64(123456789)
	destChainSelector := uint64(987654321)
	sequenceNumber := uint64(42)
	nonce := uint64(0)
	senderHex := "8765432109fedcba8765432109fedcba87654321"
	receiverHex := "0000000000000000000000000000000000000000000000000000000000001234"
	onRamp := []byte("source-onramp-address")
	data := []byte("sample message data")
	gasLimit := big.NewInt(500000)

	t.Run("should match expected message hash with no tokens", func(t *testing.T) {
		tokenReceiverHex := "0000000000000000000000000000000000000000000000000000000000000000"

		messageID := hexTo32Bytes(t, messageIDHex)
		receiver := hexTo32Bytes(t, receiverHex)
		sender, err := hex.DecodeString(senderHex)
		require.NoError(t, err)
		tokenReceiver := hexTo32Bytes(t, tokenReceiverHex)

		metadataHash, err := computeMetadataHash(sourceChainSelector, destChainSelector, onRamp)
		require.NoError(t, err)

		messageHash, err := computeMessageDataHash(
			metadataHash,
			messageID,
			receiver,
			sequenceNumber,
			gasLimit,
			tokenReceiver,
			nonce,
			sender,
			data,
			[]any2SuiTokenTransfer{},
		)
		require.NoError(t, err)

		expectedHashNoTokens := "9f9be87e216efa0b1571131d9295e3802c5c9a3d6e369d230c72520a2e854a9e"
		actualHashNoTokens := hex.EncodeToString(messageHash[:])

		assert.Equal(t, expectedHashNoTokens, actualHashNoTokens)
	})

	t.Run("should match expected message hash with tokens", func(t *testing.T) {
		tokenReceiverHex := "0000000000000000000000000000000000000000000000000000000000005678"

		messageID := hexTo32Bytes(t, messageIDHex)
		receiver := hexTo32Bytes(t, receiverHex)
		sender, err := hex.DecodeString(senderHex)
		require.NoError(t, err)
		tokenReceiver := hexTo32Bytes(t, tokenReceiverHex)

		sourcePoolAddress1, err := hex.DecodeString("abcdef1234567890abcdef1234567890abcdef12")
		require.NoError(t, err)
		destTokenAddress1 := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
		extraData1, err := hex.DecodeString("00112233")
		require.NoError(t, err)

		sourcePoolAddress2, err := hex.DecodeString("123456789abcdef123456789abcdef123456789a")
		require.NoError(t, err)
		destTokenAddress2 := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000009abc")
		extraData2, err := hex.DecodeString("ffeeddcc")
		require.NoError(t, err)

		tokenAmounts := []any2SuiTokenTransfer{
			{
				SourcePoolAddress: sourcePoolAddress1,
				DestTokenAddress:  destTokenAddress1,
				DestGasAmount:     10000,
				ExtraData:         extraData1,
				Amount:            big.NewInt(1000000),
			},
			{
				SourcePoolAddress: sourcePoolAddress2,
				DestTokenAddress:  destTokenAddress2,
				DestGasAmount:     20000,
				ExtraData:         extraData2,
				Amount:            big.NewInt(5000000),
			},
		}

		metadataHash, err := computeMetadataHash(sourceChainSelector, destChainSelector, onRamp)
		require.NoError(t, err)

		messageHash, err := computeMessageDataHash(
			metadataHash,
			messageID,
			receiver,
			sequenceNumber,
			gasLimit,
			tokenReceiver,
			nonce,
			sender,
			data,
			tokenAmounts,
		)
		require.NoError(t, err)

		expectedHashWithTokens := "d183d22cb0b713da1b6b42d9c35cc9e1268257ff703c6579d6aa68fdfb1ff4b2"
		actualHashWithTokens := hex.EncodeToString(messageHash[:])

		assert.Equal(t, expectedHashWithTokens, actualHashWithTokens)
	})

	t.Run("hashes should be different when tokens are included", func(t *testing.T) {
		tokenReceiverNoTokens := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000000000")
		tokenReceiverWithTokens := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")

		messageID := hexTo32Bytes(t, messageIDHex)
		receiver := hexTo32Bytes(t, receiverHex)
		sender, err := hex.DecodeString(senderHex)
		require.NoError(t, err)

		metadataHash, err := computeMetadataHash(sourceChainSelector, destChainSelector, onRamp)
		require.NoError(t, err)

		hashNoTokens, err := computeMessageDataHash(
			metadataHash,
			messageID,
			receiver,
			sequenceNumber,
			gasLimit,
			tokenReceiverNoTokens,
			nonce,
			sender,
			data,
			[]any2SuiTokenTransfer{},
		)
		require.NoError(t, err)

		sourcePoolAddress1, err := hex.DecodeString("abcdef1234567890abcdef1234567890abcdef12")
		require.NoError(t, err)
		destTokenAddress1 := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
		extraData1, err := hex.DecodeString("00112233")
		require.NoError(t, err)

		tokenAmounts := []any2SuiTokenTransfer{
			{
				SourcePoolAddress: sourcePoolAddress1,
				DestTokenAddress:  destTokenAddress1,
				DestGasAmount:     10000,
				ExtraData:         extraData1,
				Amount:            big.NewInt(1000000),
			},
		}

		hashWithTokens, err := computeMessageDataHash(
			metadataHash,
			messageID,
			receiver,
			sequenceNumber,
			gasLimit,
			tokenReceiverWithTokens,
			nonce,
			sender,
			data,
			tokenAmounts,
		)
		require.NoError(t, err)

		assert.NotEqual(t, hashNoTokens, hashWithTokens)
	})
}

func TestMessageHasherV2_DifferentObjectIds_DifferentHash(t *testing.T) {
	messageID := hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	receiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001234")
	sender, err := hex.DecodeString("8765432109fedcba8765432109fedcba87654321")
	require.NoError(t, err)
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000000000")
	data := []byte("sample message data")
	gasLimit := big.NewInt(500000)
	sequenceNumber := uint64(42)
	nonce := uint64(0)

	metadataHash, err := computeMetadataHash(uint64(123456789), uint64(987654321), []byte("source-onramp-address"))
	require.NoError(t, err)

	objectIdA := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001111")
	objectIdB := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000002222")

	t.Run("hash with object IDs [A] differs from hash with object IDs [B]", func(t *testing.T) {
		hashA, err := computeMessageDataHashV2(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{}, [][32]byte{objectIdA},
		)
		require.NoError(t, err)

		hashB, err := computeMessageDataHashV2(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{}, [][32]byte{objectIdB},
		)
		require.NoError(t, err)

		assert.NotEqual(t, hashA, hashB)
	})

	t.Run("hash with empty object IDs differs from hash with [A]", func(t *testing.T) {
		hashEmpty, err := computeMessageDataHashV2(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{}, [][32]byte{},
		)
		require.NoError(t, err)

		hashWithA, err := computeMessageDataHashV2(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{}, [][32]byte{objectIdA},
		)
		require.NoError(t, err)

		assert.NotEqual(t, hashEmpty, hashWithA)
	})

	t.Run("hash with [A,B] differs from [B,A] (order matters)", func(t *testing.T) {
		hashAB, err := computeMessageDataHashV2(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{}, [][32]byte{objectIdA, objectIdB},
		)
		require.NoError(t, err)

		hashBA, err := computeMessageDataHashV2(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{}, [][32]byte{objectIdB, objectIdA},
		)
		require.NoError(t, err)

		assert.NotEqual(t, hashAB, hashBA)
	})

	t.Run("V2 hash with empty object IDs differs from V1 hash (different structure)", func(t *testing.T) {
		hashV1, err := computeMessageDataHash(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{},
		)
		require.NoError(t, err)

		hashV2Empty, err := computeMessageDataHashV2(
			metadataHash, messageID, receiver, sequenceNumber, gasLimit, tokenReceiver, nonce,
			sender, data, []any2SuiTokenTransfer{}, [][32]byte{},
		)
		require.NoError(t, err)

		assert.NotEqual(t, hashV1, hashV2Empty,
			"V2 hash must differ from V1 even with empty objectIds because V2 includes the objectIdsHash term")
	})
}

func TestMessageHasherV2_Deterministic(t *testing.T) {
	messageID := hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	receiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000001234")
	sender, err := hex.DecodeString("8765432109fedcba8765432109fedcba87654321")
	require.NoError(t, err)
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
	data := []byte("test payload")
	gasLimit := big.NewInt(200000)

	metadataHash, err := computeMetadataHash(uint64(1000), uint64(2000), []byte("onramp"))
	require.NoError(t, err)

	objectId := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000aabbcc")

	hash1, err := computeMessageDataHashV2(
		metadataHash, messageID, receiver, uint64(1), gasLimit, tokenReceiver, uint64(0),
		sender, data, []any2SuiTokenTransfer{}, [][32]byte{objectId},
	)
	require.NoError(t, err)

	hash2, err := computeMessageDataHashV2(
		metadataHash, messageID, receiver, uint64(1), gasLimit, tokenReceiver, uint64(0),
		sender, data, []any2SuiTokenTransfer{}, [][32]byte{objectId},
	)
	require.NoError(t, err)

	assert.Equal(t, hash1, hash2, "same inputs must produce same hash")

	// Cross-language parity: this expected value matches the Move test
	// test_calculate_message_hash_v2_parity in offramp_test.move
	expectedHashHex := "1463b1b58f28f74dd73d4447da139d065051ddbb292549847a8c315d19148fc1"
	actualHashHex := hex.EncodeToString(hash1[:])
	assert.Equal(t, expectedHashHex, actualHashHex,
		"Go V2 hash must match Move calculate_message_hash_v2 for identical inputs")
}

func TestHashExecutionReportV2_MatchesMoveParity(t *testing.T) {
	messageID := hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	objectID := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000aabbcc")
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
	sender, err := hex.DecodeString("8765432109fedcba8765432109fedcba87654321")
	require.NoError(t, err)

	report := &codec.ExecutionReportV2{
		SourceChainSelector: 1000,
		Message: codec.Any2SuiRampMessageV2{
			Header: codec.RampMessageHeader{
				MessageID:           messageID[:],
				SourceChainSelector: 1000,
				DestChainSelector:   2000,
				SequenceNumber:      1,
				Nonce:               0,
			},
			Sender: sender,
			Data:   []byte("test payload"),
			Receiver: models.SuiAddress(
				"0000000000000000000000000000000000000000000000000000000000001234",
			),
			GasLimit:          big.NewInt(200000),
			TokenReceiver:     models.SuiAddressBytes(tokenReceiver),
			ReceiverObjectIds: []models.SuiAddressBytes{models.SuiAddressBytes(objectID)},
			TokenAmounts:      []codec.Any2SuiTokenTransfer{},
		},
	}

	hasher := NewMessageHasherV2(logger.Test(t))
	hash, err := hasher.HashExecutionReportV2(context.Background(), report, []byte("onramp"))
	require.NoError(t, err)

	expectedHashHex := "1463b1b58f28f74dd73d4447da139d065051ddbb292549847a8c315d19148fc1"
	assert.Equal(t, expectedHashHex, hex.EncodeToString(hash[:]))
}

func TestHashExecutionReportV2_AfterDeserializeRoundtrip(t *testing.T) {
	messageID := hexTo32Bytes(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	objectID := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000aabbcc")
	tokenReceiver := hexTo32Bytes(t, "0000000000000000000000000000000000000000000000000000000000005678")
	sender, err := hex.DecodeString("8765432109fedcba8765432109fedcba87654321")
	require.NoError(t, err)

	report := testutils.ExecutionReportV2{
		SourceChainSelector: 1000,
		Message: testutils.Any2SuiRampMessageV2{
			Header: testutils.NewRampMessageHeader(messageID[:], 1000, 2000, 1, 0),
			Sender: sender,
			Data:   []byte("test payload"),
			Receiver: []byte(
				"0000000000000000000000000000000000000000000000000000000000001234",
			),
			GasLimit:          big.NewInt(200000),
			TokenReceiver:     tokenReceiver[:],
			ReceiverObjectIds: [][]byte{objectID[:]},
			TokenAmounts:      []testutils.Any2SuiTokenTransfer{},
		},
	}

	encoded, err := testutils.SerializeExecutionReportV2(report)
	require.NoError(t, err)

	decoded, err := codec.DeserializeExecutionReportV2(encoded)
	require.NoError(t, err)

	hasher := NewMessageHasherV2(logger.Test(t))
	hashFromReport, err := hasher.HashExecutionReportV2(context.Background(), decoded, []byte("onramp"))
	require.NoError(t, err)

	directReport := &codec.ExecutionReportV2{
		SourceChainSelector: decoded.SourceChainSelector,
		Message:             decoded.Message,
	}
	hashDirect, err := hasher.HashExecutionReportV2(context.Background(), directReport, []byte("onramp"))
	require.NoError(t, err)

	assert.Equal(t, hashDirect, hashFromReport)

	expectedHashHex := "1463b1b58f28f74dd73d4447da139d065051ddbb292549847a8c315d19148fc1"
	assert.Equal(t, expectedHashHex, hex.EncodeToString(hashFromReport[:]))
}

// Helper function to convert hex string to [32]byte array
func hexTo32Bytes(t *testing.T, hexStr string) [32]byte {
	bytes, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	require.Len(t, bytes, 32, "hex string must decode to exactly 32 bytes")

	var result [32]byte
	copy(result[:], bytes)
	return result
}
