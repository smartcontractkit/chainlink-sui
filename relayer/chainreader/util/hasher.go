package util

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// Hasher implementation adapted from Aptos hasher for Sui
// With the following modifications:
// - Uses Sui-specific types: Any2SuiRampMessage, Any2SuiTokenTransfer, etc.
// - Uses Sui-specific constants: leafDomainSeparator and any2SuiMessageHash
// - Adapted for Sui address format (32-byte addresses vs 20-byte in Aptos)

var (
	// const LEAF_DOMAIN_SEPARATOR: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000000";
	leafDomainSeparator = [32]byte{}

	// see hash::keccak256(b"Any2SuiMessageHashV1") in calculate_metadata_hash
	any2SuiMessageHash = keccak256Fixed([]byte("Any2SuiMessageHashV1"))
)

type MessageHasherV1 struct {
	lggr logger.Logger
}

type any2SuiTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  [32]byte
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

func NewMessageHasherV1(lggr logger.Logger) *MessageHasherV1 {
	return &MessageHasherV1{
		lggr: lggr,
	}
}

func (h *MessageHasherV1) Hash(ctx context.Context, report *codec.ExecutionReport, onRampAddress []byte) ([32]byte, error) {
	rampTokenAmounts := make([]any2SuiTokenTransfer, len(report.Message.TokenAmounts))
	for i, rta := range report.Message.TokenAmounts {
		// Convert Sui address to 32-byte array
		var destTokenAddress [32]byte

		// Handle the conversion from models.SuiAddress to [32]byte
		destTokenBytes, err := hex.DecodeString("0x" + string(rta.DestTokenAddress))
		if err != nil {
			return [32]byte{}, fmt.Errorf("failed to decode dest token address: %w", err)
		}
		if len(destTokenBytes) != 32 {
			return [32]byte{}, fmt.Errorf("invalid dest token address length: expected 32, got %d", len(destTokenBytes))
		}
		copy(destTokenAddress[:], destTokenBytes)

		rampTokenAmounts[i] = any2SuiTokenTransfer{
			SourcePoolAddress: rta.SourcePoolAddress,
			DestTokenAddress:  destTokenAddress,
			DestGasAmount:     rta.DestGasAmount,
			ExtraData:         rta.ExtraData,
			Amount:            rta.Amount,
		}
	}

	metaDataHash, err := computeMetadataHash(
		report.SourceChainSelector,
		report.Message.Header.DestChainSelector,
		onRampAddress,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("compute metadata hash: %w", err)
	}

	var messageID [32]byte
	copy(messageID[:], report.Message.Header.MessageID)

	// Convert Sui address to 32-byte array
	var receiverAddress [32]byte
	receiverBytes, err := hex.DecodeString("0x" + string(report.Message.Receiver))
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode receiver address: %w", err)
	}
	if len(receiverBytes) != 32 {
		return [32]byte{}, fmt.Errorf("invalid receiver address length: expected 32, got %d", len(receiverBytes))
	}
	copy(receiverAddress[:], receiverBytes)

	msgHash, err := computeMessageDataHash(
		metaDataHash,
		messageID,
		receiverAddress,
		report.Message.Header.SequenceNumber,
		report.Message.GasLimit,
		report.Message.Header.Nonce,
		report.Message.Sender,
		report.Message.Data,
		rampTokenAmounts,
	)
	if err != nil {
		return [32]byte{}, fmt.Errorf("compute message hash: %w", err)
	}

	return msgHash, nil
}

func computeMessageDataHash(
	metadataHash [32]byte,
	messageID [32]byte,
	receiver [32]byte,
	sequenceNumber uint64,
	gasLimit *big.Int,
	nonce uint64,
	sender []byte,
	data []byte,
	tokenAmounts []any2SuiTokenTransfer,
) ([32]byte, error) {
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create uint64 ABI type: %w", err)
	}

	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create uint256 ABI type: %w", err)
	}

	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create bytes32 ABI type: %w", err)
	}

	headerArgs := abi.Arguments{
		{Type: bytes32Type}, // messageID
		{Type: bytes32Type}, // receiver as bytes32
		{Type: uint64Type},  // sequenceNumber
		{Type: uint256Type}, // gasLimit
		{Type: uint64Type},  // nonce
	}
	headerEncoded, err := headerArgs.Pack(
		messageID,
		receiver,
		sequenceNumber,
		gasLimit,
		nonce,
	)
	if err != nil {
		return [32]byte{}, err
	}
	headerHash := crypto.Keccak256Hash(headerEncoded)

	senderHash := crypto.Keccak256Hash(sender)

	dataHash := crypto.Keccak256Hash(data)

	// Manually encode tokens to match the Move implementation
	var tokenHashData []byte
	tokenHashData = append(tokenHashData, encodeUint256(big.NewInt(int64(len(tokenAmounts))))...)
	for _, token := range tokenAmounts {
		tokenHashData = append(tokenHashData, encodeBytes(token.SourcePoolAddress)...)
		tokenHashData = append(tokenHashData, token.DestTokenAddress[:]...)
		tokenHashData = append(tokenHashData, encodeUint32(token.DestGasAmount)...)
		tokenHashData = append(tokenHashData, encodeBytes(token.ExtraData)...)
		tokenHashData = append(tokenHashData, encodeUint256(token.Amount)...)
	}
	tokenAmountsHash := crypto.Keccak256Hash(tokenHashData)

	finalArgs := abi.Arguments{
		{Type: bytes32Type}, // LEAF_DOMAIN_SEPARATOR
		{Type: bytes32Type}, // metadataHash
		{Type: bytes32Type}, // headerHash
		{Type: bytes32Type}, // senderHash
		{Type: bytes32Type}, // dataHash
		{Type: bytes32Type}, // tokenAmountsHash
	}

	finalEncoded, err := finalArgs.Pack(
		leafDomainSeparator,
		metadataHash,
		headerHash,
		senderHash,
		dataHash,
		tokenAmountsHash,
	)
	if err != nil {
		return [32]byte{}, err
	}

	return crypto.Keccak256Hash(finalEncoded), nil
}

func computeMetadataHash(
	sourceChainSelector uint64,
	destinationChainSelector uint64,
	onRamp []byte,
) ([32]byte, error) {
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create uint64 ABI type: %w", err)
	}

	bytes32Type, err := abi.NewType("bytes32", "", nil)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to create bytes32 ABI type: %w", err)
	}

	onRampHash := crypto.Keccak256Hash(onRamp)

	args := abi.Arguments{
		{Type: bytes32Type}, // ANY_2_SUI_MESSAGE_HASH
		{Type: uint64Type},  // sourceChainSelector
		{Type: uint64Type},  // destinationChainSelector (i_chainSelector)
		{Type: bytes32Type}, // onRamp
	}

	encoded, err := args.Pack(
		any2SuiMessageHash,
		sourceChainSelector,
		destinationChainSelector,
		onRampHash,
	)
	if err != nil {
		return [32]byte{}, err
	}

	metadataHash := crypto.Keccak256Hash(encoded)

	return metadataHash, nil
}

func encodeUint256(n *big.Int) []byte {
	return common.LeftPadBytes(n.Bytes(), 32)
}

func encodeUint32(n uint32) []byte {
	return common.LeftPadBytes(new(big.Int).SetUint64(uint64(n)).Bytes(), 32)
}

func encodeBytes(b []byte) []byte {
	encodedLength := common.LeftPadBytes(big.NewInt(int64(len(b))).Bytes(), 32)
	padLen := (32 - (len(b) % 32)) % 32
	result := make([]byte, 32+len(b)+padLen)
	copy(result[:32], encodedLength)
	copy(result[32:], b)

	return result
}

func keccak256Fixed(in []byte) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	hash.Write(in)
	var h [32]byte
	copy(h[:], hash.Sum(nil))

	return h
}
