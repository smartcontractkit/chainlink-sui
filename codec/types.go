package codec

import (
	"errors"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
)

var AccountZero = make([]byte, 32)

type PTBCommandDependency struct {
	CommandIndex uint16
	ResultIndex  *uint16
}

// PointerTag defines the structured format for pointer tags used in chain reader.
// Pointer tags specify how to derive object IDs from pointer objects stored on-chain.
type PointerTag struct {
	// Module name containing the pointer object (e.g. "state_object", "offramp", "counter")
	Module string `json:"module"`
	// PointerName is the object type to search for (e.g. "CCIPObjectRefPointer", "OffRampStatePointer")
	PointerName string `json:"pointerName"`
	// FieldName is OPTIONAL and NOT USED by the implementation. The parent field name is automatically
	// looked up from the global common.PointerConfigs registry based on the PointerName.
	// This field exists for backward compatibility or future implementations to override static code but is currently ignored.
	FieldName string `json:"fieldName,omitempty"`
	// DerivationKey is the key used to derive the child object ID from the parent object ID (e.g. "CCIPObjectRef", "CCIP_OWNABLE")
	DerivationKey string `json:"derivationKey"`
	// PackageID is the package ID for the Pointer object if it differs from the calling contract's package ID
	// This is used for cross-package pointer dependencies (e.g. offramp package depending on CCIP package CCIPObjectRef)
	// If empty, the calling contract's package ID is used
	PackageID string `json:"packageId,omitempty"`
}

func (p PointerTag) Validate() error {
	if p.Module == "" {
		return errors.New("PointerTag.Module is required")
	}
	if p.PointerName == "" {
		return errors.New("PointerTag.Pointer is required")
	}
	// FieldName is optional - it's looked up from common.PointerConfigs
	if p.DerivationKey == "" {
		return errors.New("PointerTag.DerivationKey is required")
	}
	return nil
}

// SuiFunctionParam defines a parameter for a Sui function call
type SuiFunctionParam struct {
	// Name of the parameter
	Name string
	// PointerTag (optional) specify how to derive object IDs from pointer objects stored on-chain.
	PointerTag *PointerTag
	// Type of the parameter (e.g., "u64", "String", "vector<u8>", "ptb_dependency")
	Type string
	// IsMutable specifies if the object is mutable or not (optional - defaults to true)
	IsMutable *bool
	// IsGeneric specifies if the parameter is a generic argument
	GenericType *string
	// Whether the parameter is required
	Required bool
	// Default value to use if not provided
	DefaultValue any
	// Result from a previous PTB Command (optional). It is used for expressive construction of PTB commands
	PTBDependency *PTBCommandDependency
	// GenericDependency maps to internal helpers for fetching an unknown generic type required by the parameter
	GenericDependency *string
}

type SuiPTBCommandType string

const (
	SuiPTBCommandMoveCall SuiPTBCommandType = "move_call"
	SuiPTBCommandPublish  SuiPTBCommandType = "publish"
	SuiPTBCommandTransfer SuiPTBCommandType = "transfer"
)

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
