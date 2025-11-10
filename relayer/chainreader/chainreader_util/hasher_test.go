package chainreaderutil

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// Helper function to convert hex string to [32]byte array
func hexTo32Bytes(t *testing.T, hexStr string) [32]byte {
	bytes, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	require.Len(t, bytes, 32, "hex string must decode to exactly 32 bytes")

	var result [32]byte
	copy(result[:], bytes)
	return result
}
