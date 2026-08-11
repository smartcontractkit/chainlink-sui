package codec

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"unicode"

	"github.com/block-vision/sui-go-sdk/models"

	common "github.com/smartcontractkit/chainlink-common/pkg/types/sui"
)

var AccountZero = make([]byte, 32)

// Object represents a Sui object with its ID and optional initial shared version.
type Object struct {
	Id                   string
	InitialSharedVersion *uint64
}

// IsSuiAddress returns true if addr is a valid Sui address/ObjectID.
// It is an improvement over the sui-go-sdk's IsValidSuiAddress function.
func IsSuiAddress(addr string) bool {
	if !(strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X")) {
		return false
	}

	h := addr[2:]

	// 1..64 hex chars
	if len(h) == 0 || len(h) > 64 {
		return false
	}

	// hex only
	for _, r := range h {
		if !isHexRune(r) {
			return false
		}
	}

	// hex.DecodeString requires even length; allow odd by left-padding one '0'
	if len(h)%2 == 1 {
		h = "0" + h
	}

	b, err := hex.DecodeString(h)
	if err != nil {
		return false
	}

	// must be ≤32 bytes (32 bytes == 64 hex chars)
	return len(b) <= 32
}

// ToSuiAddress normalizes and validates a Sui address.
// It pads the address to 64 hex characters (32 bytes) with leading zeros.
func ToSuiAddress(address string) (string, error) {
	normalized := NormalizeSuiAddress(address)
	if !IsSuiAddress(string(normalized)) {
		return "", fmt.Errorf("invalid sui address: %s", address)
	}

	return string(normalized), nil
}

// NormalizeSuiAddress normalizes a Sui address to 64 hex characters.
func NormalizeSuiAddress(address string) string {
	if !strings.HasPrefix(address, "0x") && !strings.HasPrefix(address, "0X") {
		address = "0x" + address
	}

	// Remove 0x prefix for processing
	h := address[2:]

	// Pad with leading zeros to 64 characters (32 bytes)
	if len(h) < 64 {
		h = strings.Repeat("0", 64-len(h)) + h
	}

	// Convert to lowercase
	h = strings.ToLower(h)

	return "0x" + h
}

func isHexRune(r rune) bool {
	return unicode.IsDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

//go:fix inline
type PTBCommandDependency = common.PTBCommandDependency

//go:fix inline
type PointerTag = common.PointerTag

//go:fix inline
type SuiFunctionParam = common.SuiFunctionParam

//go:fix inline
type SuiPTBCommandType = common.SuiPTBCommandType

// OCRConfigSet event data
type ConfigSet struct {
	OcrPluginType byte
	ConfigDigest  []byte
	Signers       [][]byte
	// this is a list of addresses, we can treat them as strings
	Transmitters []string
	BigF         byte
}

// SourceChainConfigSet event data
type SourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceChainConfig   SourceChainConfig
}

// SourceChainConfig event data
type SourceChainConfig struct {
	Router                    string
	IsEnabled                 bool
	MinSeqNr                  uint64
	IsRMNVerificationDisabled bool
	OnRamp                    []byte
}

// ExecutionReport event data
type ExecutionReport struct {
	SourceChainSelector uint64
	Message             Any2SuiRampMessage
	OffchainTokenData   [][]byte
	Proofs              [][]byte
}

// RampMessageHeader event data
type RampMessageHeader struct {
	MessageID           []byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

// Any2SuiTokenTransfer event data
type Any2SuiTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  models.SuiAddress
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

// Any2SuiRampMessage event data
type Any2SuiRampMessage struct {
	Header        RampMessageHeader
	Sender        []byte
	Data          []byte
	Receiver      models.SuiAddress
	GasLimit      *big.Int
	TokenReceiver models.SuiAddressBytes
	TokenAmounts  []Any2SuiTokenTransfer
}

// ExecutionStateChanged event data
type ExecutionStateChanged struct {
	SourceChainSelector uint64 `json:"source_chain_selector"`
	SequenceNumber      uint64 `json:"sequence_number"`
	MessageId           []byte `json:"message_id"`
	MessageHash         []byte `json:"message_hash"`
	State               byte   `json:"state"`
}
