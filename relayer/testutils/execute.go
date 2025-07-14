package testutils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
)

// ExecutionReport represents the execution report structure from offramp.move
// Matches the Move struct: public struct ExecutionReport has drop
type ExecutionReport struct {
	SourceChainSelector uint64             `json:"source_chain_selector"`
	Message             Any2SuiRampMessage `json:"message"`
	OffchainTokenData   [][]byte           `json:"offchain_token_data"`
	Proofs              [][]byte           `json:"proofs"`
}

// RampMessageHeader represents the message header structure from offramp.move
// Matches the Move struct: public struct RampMessageHeader has drop
type RampMessageHeader struct {
	MessageId           []byte `json:"message_id"`
	SourceChainSelector uint64 `json:"source_chain_selector"`
	DestChainSelector   uint64 `json:"dest_chain_selector"`
	SequenceNumber      uint64 `json:"sequence_number"`
	Nonce               uint64 `json:"nonce"`
}

// Any2SuiRampMessage represents the ramp message structure from offramp.move
// Matches the Move struct: public struct Any2SuiRampMessage has drop
type Any2SuiRampMessage struct {
	Header       RampMessageHeader      `json:"header"`
	Sender       []byte                 `json:"sender"`
	Data         []byte                 `json:"data"`
	Receiver     []byte                 `json:"receiver"`  // address in Move becomes []byte for 32-byte Sui address
	GasLimit     *big.Int               `json:"gas_limit"` // u256 in Move becomes *big.Int in Go
	TokenAmounts []Any2SuiTokenTransfer `json:"token_amounts"`
}

// Any2SuiTokenTransfer represents a token transfer structure from offramp.move
// Matches the Move struct: public struct Any2SuiTokenTransfer has drop
type Any2SuiTokenTransfer struct {
	SourcePoolAddress []byte   `json:"source_pool_address"`
	DestTokenAddress  []byte   `json:"dest_token_address"` // address in Move becomes []byte for 32-byte Sui address
	DestGasAmount     uint32   `json:"dest_gas_amount"`
	ExtraData         []byte   `json:"extra_data"`
	Amount            *big.Int `json:"amount"` // u256 in Move becomes *big.Int in Go
}

// Helper function to create a new ExecutionReport with proper initialization
func NewExecutionReport(sourceChainSelector uint64) *ExecutionReport {
	return &ExecutionReport{
		SourceChainSelector: sourceChainSelector,
		Message:             Any2SuiRampMessage{},
		OffchainTokenData:   make([][]byte, 0),
		Proofs:              make([][]byte, 0),
	}
}

// Helper function to create a RampMessageHeader
func NewRampMessageHeader(
	messageId []byte,
	sourceChainSelector uint64,
	destChainSelector uint64,
	sequenceNumber uint64,
	nonce uint64,
) RampMessageHeader {
	return RampMessageHeader{
		MessageId:           messageId,
		SourceChainSelector: sourceChainSelector,
		DestChainSelector:   destChainSelector,
		SequenceNumber:      sequenceNumber,
		Nonce:               nonce,
	}
}

// Helper function to create an Any2SuiRampMessage
func NewAny2SuiRampMessage(
	header RampMessageHeader,
	sender []byte,
	data []byte,
	receiver []byte,
	gasLimit *big.Int,
	tokenAmounts []Any2SuiTokenTransfer,
) Any2SuiRampMessage {
	if gasLimit == nil {
		gasLimit = big.NewInt(0)
	}
	if tokenAmounts == nil {
		tokenAmounts = make([]Any2SuiTokenTransfer, 0)
	}

	return Any2SuiRampMessage{
		Header:       header,
		Sender:       sender,
		Data:         data,
		Receiver:     receiver,
		GasLimit:     gasLimit,
		TokenAmounts: tokenAmounts,
	}
}

// Helper function to create an Any2SuiTokenTransfer
func NewAny2SuiTokenTransfer(
	sourcePoolAddress []byte,
	destTokenAddress []byte,
	destGasAmount uint32,
	extraData []byte,
	amount *big.Int,
) Any2SuiTokenTransfer {
	if amount == nil {
		amount = big.NewInt(0)
	}
	if extraData == nil {
		extraData = make([]byte, 0)
	}

	return Any2SuiTokenTransfer{
		SourcePoolAddress: sourcePoolAddress,
		DestTokenAddress:  destTokenAddress,
		DestGasAmount:     destGasAmount,
		ExtraData:         extraData,
		Amount:            amount,
	}
}

// Helper function to create an ExecutionReport from CCIP types
func GetExecutionReportFromCCIP(
	sourceChainSelector uint64,
	message cciptypes.Message,
	offchainTokenData [][]byte,
	proofs [][]byte,
	gasAmount uint32,
) ExecutionReport {
	// Convert CCIP message to Move message format
	// Convert MessageID from [32]byte to []byte
	messageIDBytes := make([]byte, DefaultByteSize)
	copy(messageIDBytes, message.Header.MessageID[:])
	header := NewRampMessageHeader(
		messageIDBytes,
		sourceChainSelector,
		uint64(message.Header.DestChainSelector),
		uint64(message.Header.SequenceNumber),
		message.Header.Nonce,
	)

	// Convert token amounts
	tokenAmounts := make([]Any2SuiTokenTransfer, len(message.TokenAmounts))
	for i, tokenAmount := range message.TokenAmounts {
		tokenAmounts[i] = NewAny2SuiTokenTransfer(
			tokenAmount.SourcePoolAddress,
			tokenAmount.DestTokenAddress,
			gasAmount,
			tokenAmount.ExtraData,
			tokenAmount.Amount.Int,
		)
	}

	// Create the ramp message
	rampMessage := NewAny2SuiRampMessage(
		header,
		message.Sender,
		message.Data,
		message.Receiver,
		big.NewInt(int64(gasAmount)),
		tokenAmounts,
	)

	return ExecutionReport{
		SourceChainSelector: sourceChainSelector,
		Message:             rampMessage,
		OffchainTokenData:   offchainTokenData,
		Proofs:              proofs,
	}
}

// SerializeExecutionReport serializes an ExecutionReport using BCS format to match the Move contract's expected deserialization format.
// The Move contract expects the following order:
// 1. SourceChainSelector (u64)
// 2. Message (Any2SuiRampMessage)
// 3. OffchainTokenData (vector<vector<u8>>)
// 4. Proofs (vector<vector<u8>>)
func SerializeExecutionReport(report ExecutionReport) ([]byte, error) {
	s := &bcs.Serializer{}

	// Serialize SourceChainSelector as u64
	s.U64(report.SourceChainSelector)
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize SourceChainSelector: %w", s.Error())
	}

	// Serialize Message (Any2SuiRampMessage)
	if err := serializeAny2SuiRampMessage(s, report.Message); err != nil {
		return nil, fmt.Errorf("failed to serialize Message: %w", err)
	}

	// Serialize OffchainTokenData as vector<vector<u8>>
	bcs.SerializeSequenceWithFunction(report.OffchainTokenData, s, func(s *bcs.Serializer, item []byte) {
		s.WriteBytes(item)
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize OffchainTokenData: %w", s.Error())
	}

	// Serialize Proofs as vector<vector<u8>>
	bcs.SerializeSequenceWithFunction(report.Proofs, s, func(s *bcs.Serializer, item []byte) {
		s.WriteBytes(item)
	})
	if s.Error() != nil {
		return nil, fmt.Errorf("failed to serialize Proofs: %w", s.Error())
	}

	return s.ToBytes(), nil
}

// serializeAny2SuiRampMessage serializes an Any2SuiRampMessage struct
func serializeAny2SuiRampMessage(s *bcs.Serializer, message Any2SuiRampMessage) error {
	// Serialize Header (RampMessageHeader)
	if err := serializeRampMessageHeader(s, message.Header); err != nil {
		return fmt.Errorf("failed to serialize Header: %w", err)
	}

	// Serialize Sender as vector<u8>
	s.WriteBytes(message.Sender)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize Sender: %w", s.Error())
	}

	// Serialize Data as vector<u8>
	// Handle nil data by using empty byte slice
	dataBytes := message.Data
	if dataBytes == nil {
		dataBytes = make([]byte, 0)
	}
	s.WriteBytes(dataBytes)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize Data: %w", s.Error())
	}

	// Serialize Receiver as address (32 bytes)
	// Handle empty receiver by using zero address
	receiverBytes := make([]byte, DefaultByteSize)
	if len(message.Receiver) > 0 && string(message.Receiver) != "0x" {
		receiverStr := string(message.Receiver)

		// If receiver contains "::", extract just the address part (before first "::")
		if strings.Contains(receiverStr, "::") {
			parts := strings.Split(receiverStr, "::")
			if len(parts) >= 1 {
				receiverStr = parts[0] // Take only the address part
			}
		}

		// Remove "0x" prefix if present
		receiverStr = strings.TrimPrefix(receiverStr, "0x")

		// Decode hex string to bytes
		if decoded, err := hex.DecodeString(receiverStr); err == nil && len(decoded) <= DefaultByteSize {
			// Right-pad to 32 bytes (Sui addresses are 32 bytes)
			copy(receiverBytes[DefaultByteSize-len(decoded):], decoded)
		}
		// If decoding fails or address is invalid, receiverBytes remains all zeros
	}

	// If message.Receiver is empty or "0x", receiverBytes remains all zeros (zero address)
	s.FixedBytes(receiverBytes)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize Receiver: %w", s.Error())
	}

	// Serialize GasLimit as u256
	if message.GasLimit == nil {
		s.U256(*big.NewInt(0))
	} else {
		s.U256(*message.GasLimit)
	}
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize GasLimit: %w", s.Error())
	}

	// Serialize TokenAmounts as vector<Any2SuiTokenTransfer>
	bcs.SerializeSequenceWithFunction(message.TokenAmounts, s, func(s *bcs.Serializer, item Any2SuiTokenTransfer) {
		if err := serializeAny2SuiTokenTransfer(s, item); err != nil {
			s.SetError(err)
		}
	})
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize TokenAmounts: %w", s.Error())
	}

	return nil
}

// serializeRampMessageHeader serializes a RampMessageHeader struct
func serializeRampMessageHeader(s *bcs.Serializer, header RampMessageHeader) error {
	messageId := make([]byte, DefaultByteSize)
	if len(header.MessageId) > 0 {
		copy(messageId, header.MessageId)
	}
	s.FixedBytes(messageId)

	if s.Error() != nil {
		return fmt.Errorf("failed to serialize MessageId: %w", s.Error())
	}

	// Serialize SourceChainSelector as u64
	s.U64(header.SourceChainSelector)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize SourceChainSelector: %w", s.Error())
	}

	// Serialize DestChainSelector as u64
	s.U64(header.DestChainSelector)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize DestChainSelector: %w", s.Error())
	}

	// Serialize SequenceNumber as u64
	s.U64(header.SequenceNumber)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize SequenceNumber: %w", s.Error())
	}

	// Serialize Nonce as u64
	s.U64(header.Nonce)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize Nonce: %w", s.Error())
	}

	return nil
}

// serializeAny2SuiTokenTransfer serializes an Any2SuiTokenTransfer struct
func serializeAny2SuiTokenTransfer(s *bcs.Serializer, transfer Any2SuiTokenTransfer) error {
	// Serialize SourcePoolAddress as vector<u8>
	s.WriteBytes(transfer.SourcePoolAddress)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize SourcePoolAddress: %w", s.Error())
	}

	// Serialize DestTokenAddress as address (32 bytes)
	if len(transfer.DestTokenAddress) != DefaultByteSize {
		return fmt.Errorf("dest token address must be exactly 32 bytes, got %d", len(transfer.DestTokenAddress))
	}
	s.FixedBytes(transfer.DestTokenAddress)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize DestTokenAddress: %w", s.Error())
	}

	// Serialize DestGasAmount as u32
	s.U32(transfer.DestGasAmount)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize DestGasAmount: %w", s.Error())
	}

	// Serialize ExtraData as vector<u8>
	s.WriteBytes(transfer.ExtraData)
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize ExtraData: %w", s.Error())
	}

	// Serialize Amount as u256
	if transfer.Amount == nil {
		s.U256(*big.NewInt(0))
	} else {
		s.U256(*transfer.Amount)
	}
	if s.Error() != nil {
		return fmt.Errorf("failed to serialize Amount: %w", s.Error())
	}

	return nil
}
