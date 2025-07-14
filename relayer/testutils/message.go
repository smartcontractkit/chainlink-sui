package testutils

import (
	"fmt"
	"math/big"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"golang.org/x/crypto/sha3"
)

var (
	// const LEAF_DOMAIN_SEPARATOR: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000000";
	leafDomainSeparator = [32]byte{}

	// see aptos_hash::keccak256(b"Any2SuiMessageHashV1") in calculate_metadata_hash
	any2SuiMessageHash = Keccak256Fixed([]byte("Any2SuiMessageHashV1"))
)

func Keccak256Fixed(in []byte) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	hash.Write(in)
	var h [32]byte
	copy(h[:], hash.Sum(nil))

	return h
}

// BCS encoding helper functions to match Move contract expectations
func encodeUint256(value *big.Int) []byte {
	s := &bcs.Serializer{}
	s.U256(*value)

	return s.ToBytes()
}

func encodeUint32(value uint32) []byte {
	s := &bcs.Serializer{}
	s.U32(value)

	return s.ToBytes()
}

func encodeBytes(value []byte) []byte {
	s := &bcs.Serializer{}
	s.WriteBytes(value)

	return s.ToBytes()
}

// This is the equivalent of ccip_offramp::calculate_message_hash.
// This is similar to the EVM version, except for 32-byte addresses and no dynamic offsets.
func ComputeMessageDataHash(
	metadataHash [32]byte,
	messageID [32]byte,
	receiver []byte,
	sequenceNumber uint64,
	gasLimit *big.Int,
	nonce uint64,
	sender []byte,
	data []byte,
	tokenAmounts []ccipocr3.RampTokenAmount,
	destGasAmount uint32,
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

	// Manually encode tokens to match the Move implementation, because abi.Pack has different behavior
	// for dynamic types.
	var tokenHashData []byte
	tokenHashData = append(tokenHashData, encodeUint256(big.NewInt(int64(len(tokenAmounts))))...)
	for _, token := range tokenAmounts {
		tokenHashData = append(tokenHashData, encodeBytes(token.SourcePoolAddress)...)
		tokenHashData = append(tokenHashData, token.DestTokenAddress[:]...)
		// Note: DestGasAmount doesn't exist in ccipocr3.RampTokenAmount, using a default value
		tokenHashData = append(tokenHashData, encodeUint32(destGasAmount)...)
		tokenHashData = append(tokenHashData, encodeBytes(token.ExtraData)...)
		tokenHashData = append(tokenHashData, encodeUint256(token.Amount.Int)...)
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

// This is the equivalent of ccip_offramp::calculate_metadata_hash.
// This is similar to the EVM version, except for the separator, 32-byte addresses, and no dynamic offsets.
// See https://github.com/smartcontractkit/chainlink-aptos/blob/d2cf1852ffdbf80fa55b0c834ebef7f44a46d843/contracts/ccip/ccip_offramp/sources/offramp.move#L1044
func ComputeMetadataHash(
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
