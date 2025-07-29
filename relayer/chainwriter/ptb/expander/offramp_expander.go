package expander

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

var (
	DEFAULT_NR_OFFRAMP_PTB_COMMANDS  = 2
	OFFRAMP_TOKEN_POOL_FUNCTION_NAME = "release_or_mint"
	SUI_PATH_COMPONENTS_COUNT        = 3
)

type SuiAddress [32]byte

const SuiAddressLength = 32

// UnmarshalJSON implements the json.Unmarshaler interface for SuiAddress
func (s *SuiAddress) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Try to decode as base64
	decoded, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return fmt.Errorf("failed to decode base64 string for SuiAddress: %w", err)
	}

	// Ensure we have exactly 32 bytes
	if len(decoded) != SuiAddressLength {
		return fmt.Errorf("SuiAddress must be exactly 32 bytes, got %d", len(decoded))
	}

	// Copy the decoded bytes to the array
	copy(s[:], decoded)

	return nil
}

type TokenPool struct {
	CoinMetadata          string
	TokenType             string // e.g. "sui:0x66::link_module::LINK"
	PackageId             string
	ModuleId              string
	Function              string
	TokenPoolStateAddress string
	Index                 int
}

type SuiOffRampExecCallArgs struct {
	ReportContext [2][32]byte                `mapstructure:"ReportContext"`
	Report        []byte                     `mapstructure:"Report"`
	Info          ccipocr3.ExecuteReportInfo `mapstructure:"Info"`
}

// OffRampPTBArgs represents arguments for OffRamp PTB expansion operations
type OffRampPTBArgs struct {
	ExecArgs   SuiOffRampExecCallArgs
	PTBConfigs *config.ChainWriterFunction
}

// OffRampPTBResult represents the result of OffRamp PTB expansion operations
type OffRampPTBResult struct {
	PTBCommands []config.ChainWriterPTBCommand
	UpdatedArgs map[string]any
	TypeArgs    map[string]string
}

// SuiPTBExpander is a concrete implementation of the PTBExpander interface for Sui blockchain
// It implements the generic interface specifically for OffRamp operations
type SuiPTBExpander struct {
	lggr            logger.Logger
	ptbClient       client.SuiPTBClient
	AddressMappings map[string]string
}

// NewSuiPTBExpander creates a new instance of SuiPTBExpander
func NewSuiPTBExpander(
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	addressMappings map[string]string,
) *SuiPTBExpander {
	return &SuiPTBExpander{
		lggr:            lggr,
		ptbClient:       ptbClient,
		AddressMappings: addressMappings,
	}
}

type GetPoolInfosResult struct {
	TokenPoolPackageIds     []SuiAddress `json:"token_pool_package_ids"`
	TokenPoolStateAddresses []SuiAddress `json:"token_pool_state_addresses"`
	TokenPoolModules        []string     `json:"token_pool_modules"`
	TokenTypes              []string     `json:"token_types"`
}

// SetupAddressMappings initializes and populates all required address mappings for PTB expansion operations.
//
// This function performs discovery and resolution of critical CCIP infrastructure addresses by:
// 1. Using the provided OffRamp package ID to query and discover the CCIP package ID
// 2. Reading owned objects to locate the OffRamp state pointer and extract the state address
// 3. Reading CCIP package objects to find the CCIP object reference and owner capability addresses
// 4. Assembling a complete address mapping required for subsequent PTB operations
//
// The returned address mappings include:
//   - ccipPackageId: The main CCIP package identifier discovered from the OffRamp
//   - ccipObjectRef: Reference to the main CCIP state object for operations
//   - ccipOwnerCap: Owner capability object for privileged operations
//   - clockObject: Sui system clock object (fixed at 0x6)
//   - offRampPackageId: The provided OffRamp package identifier
//   - offRampState: The OffRamp state object address discovered from state pointer
//
// Parameters:
//   - ctx: Context for the operation, used for request lifecycle management
//   - lggr: Logger instance for debugging and operational visibility
//   - ptbClient: Sui PTB client for reading blockchain state and objects
//   - offRampPackageId: The OffRamp package identifier to start discovery from
//   - publicKey: Public key bytes for generating signer address for read operations
//
// Returns:
//   - map[string]string: Complete address mappings required for PTB expansion
//   - error: Error if any discovery step fails, objects are missing, or network issues occur
//
// Usage:
//
//	This function should be called once during PTB expander initialization to discover
//	and cache all necessary addresses. The returned mappings are then used throughout
//	the PTB expansion process for referencing the correct on-chain objects.
//
// Example:
//
//	addressMappings, err := SetupAddressMappings(ctx, logger, ptbClient, offRampPkgId, pubKey)
//	if err != nil {
//	    return fmt.Errorf("failed to setup address mappings: %w", err)
//	}
//	expander := NewSuiPTBExpander(logger, ptbClient, addressMappings)
func SetupAddressMappings(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	offRampPackageId string,
	publicKey []byte,
) (map[string]string, error) {
	// address mappings for the expander
	addressMappings := map[string]string{
		"ccipPackageId":    "",
		"ccipObjectRef":    "",
		"ccipOwnerCap":     "",
		"clockObject":      "0x6",
		"offRampPackageId": offRampPackageId,
		"offRampState":     "",
	}

	// Use the `toAddress` (offramp package ID) from the config overrides to get the offramp pointer object
	signerAddress, err := client.GetAddressFromPublicKey(publicKey)
	if err != nil {
		lggr.Errorw("Error getting signer address", "error", err)
		return nil, err
	}
	getCCIPPackageIdResponse, err := ptbClient.ReadFunction(ctx, signerAddress, addressMappings["offRampPackageId"], "offramp", "get_ccip_package_id", []any{}, []string{})
	if err != nil {
		lggr.Errorw("Error reading ccip package id", "error", err)
		return nil, err
	}
	lggr.Debugw("getCCIPPackageIdResponse", "getCCIPPackageIdResponse", getCCIPPackageIdResponse)
	// Parse the response to get the returned address as a hex string
	var addressBytes []byte

	// Handle both byte slice and base64 string responses
	switch v := getCCIPPackageIdResponse[0].(type) {
	case []byte:
		// Response is already raw bytes ([]byte and []uint8 are the same type)
		addressBytes = v
	case string:
		// Response is base64-encoded string, decode it
		var decodeErr error
		addressBytes, decodeErr = base64.StdEncoding.DecodeString(v)
		if decodeErr != nil {
			lggr.Errorw("Error decoding base64 ccip package id", "error", decodeErr)
			return nil, decodeErr
		}
	default:
		lggr.Errorw("Unexpected type for ccip package id response", "type", fmt.Sprintf("%T", getCCIPPackageIdResponse[0]))
		return nil, fmt.Errorf("unexpected type for ccip package id response, got %T", getCCIPPackageIdResponse[0])
	}
	// Convert bytes to hex string with "0x" prefix
	ccipPackageId := "0x" + hex.EncodeToString(addressBytes)
	addressMappings["ccipPackageId"] = ccipPackageId

	lggr.Debugw("ccipPackageId", "ccipPackageId", addressMappings["ccipPackageId"])
	lggr.Debugw("offRampPackageId", "offrampPackageId", addressMappings["offRampPackageId"])

	// get the offramp state object
	offrampOwnedObjects, err := ptbClient.ReadOwnedObjects(ctx, addressMappings["offRampPackageId"], nil)
	if err != nil {
		lggr.Errorw("Error reading offramp state object", "error", err)
		return nil, err
	}
	for _, ccipOwnedObject := range offrampOwnedObjects {
		if ccipOwnedObject.Data.Type != "" && strings.Contains(ccipOwnedObject.Data.Type, "offramp::OffRampStatePointer") {
			lggr.Debugw("Found offramp state object pointer", "fields", ccipOwnedObject.Data.Content.Fields)
			// parse the object into a map
			parsedObject := ccipOwnedObject.Data.Content.Fields
			lggr.Debugw("offRampStatePointer Parsed", "offRampStatePointer", parsedObject)
			addressMappings["offRampState"] = parsedObject["off_ramp_state_id"].(string)

			break
		}
	}
	if addressMappings["offRampState"] == "" {
		lggr.Errorw("Address mappings are not populated", "addressMappings", addressMappings)
		return nil, fmt.Errorf("address mappings are missing required fields for expander (offRampState)")
	}

	// Get the object pointer present in the CCIP package ID
	ccipOwnedObjects, err := ptbClient.ReadOwnedObjects(ctx, addressMappings["ccipPackageId"], nil)
	if err != nil {
		lggr.Errorw("Error reading ccip object ref", "error", err)
		return nil, err
	}
	for _, ccipOwnedObject := range ccipOwnedObjects {
		if ccipOwnedObject.Data.Type != "" && strings.Contains(ccipOwnedObject.Data.Type, "state_object::CCIPObjectRefPointer") {
			// parse the object into a map
			parsedObject := ccipOwnedObject.Data.Content.Fields
			if err != nil {
				lggr.Errorw("Error parsing ccip object ref", "error", err)
				return nil, err
			}
			lggr.Debugw("ccipObjectRefPointer", "ccipObjectRefPointer", parsedObject)
			addressMappings["ccipObjectRef"] = parsedObject["object_ref_id"].(string)
			addressMappings["ccipOwnerCap"] = parsedObject["owner_cap_id"].(string)

			break
		}
	}
	// check that address mappings are populated
	if addressMappings["ccipObjectRef"] == "" || addressMappings["ccipOwnerCap"] == "" {
		lggr.Errorw("Address mappings are not populated", "addressMappings", addressMappings)
		return nil, fmt.Errorf("address mappings are missing required fields for expander (ccipObjectRef, ccipOwnerCap)")
	}

	lggr.Debugw("Address mappings for expander", "addressMappings", addressMappings)

	return addressMappings, nil
}

// Expand performs OffRamp PTB expansion
func (s *SuiPTBExpander) Expand(
	ctx context.Context,
	lggr logger.Logger,
	args OffRampPTBArgs,
	signerPublicKey []byte,
) (OffRampPTBResult, error) {
	ptbCmds, updatedArgs, typeArgs, err := s.getOffRampPTB(ctx, lggr, args.ExecArgs, args.PTBConfigs, signerPublicKey)
	if err != nil {
		return OffRampPTBResult{}, err
	}

	return OffRampPTBResult{PTBCommands: ptbCmds, UpdatedArgs: updatedArgs, TypeArgs: typeArgs}, nil
}

// getTokenPoolByTokenAddress gets token pool addresses for given token addresses (internal method)
func (s *SuiPTBExpander) getTokenPoolByTokenAddress(
	ctx context.Context,
	lggr logger.Logger,
	tokenAmounts []ccipocr3.RampTokenAmount,
	signerPublicKey []byte,
) ([]TokenPool, error) {
	coinMetadataAddresses := make([]string, len(tokenAmounts))
	for i, tokenAmount := range tokenAmounts {
		address := tokenAmount.DestTokenAddress
		coinMetadataAddresses[i] = "0x" + hex.EncodeToString(address)
	}

	lggr.Debugw("getting token pool infos",
		"packageID", s.AddressMappings["ccipPackageId"],
		"ccipObjectRef", s.AddressMappings["ccipObjectRef"],
		"coinMetadataAddresses", coinMetadataAddresses)

	signerAddress, err := client.GetAddressFromPublicKey(signerPublicKey)
	if err != nil {
		return nil, err
	}

	poolInfos, err := s.ptbClient.ReadFunction(
		ctx,
		signerAddress,
		s.AddressMappings["ccipPackageId"],
		"token_admin_registry",
		"get_pool_infos",
		[]any{
			s.AddressMappings["ccipObjectRef"],
			coinMetadataAddresses,
		},
		[]string{"object_id", "vector<address>"},
	)
	if err != nil {
		lggr.Errorw("Error getting pool infos", "error", err)
		return nil, err
	}

	var tokenPoolInfo GetPoolInfosResult
	lggr.Debugw("tokenPoolInfo", "tokenPoolInfo", poolInfos[0])
	jsonBytes, err := json.Marshal(poolInfos[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &tokenPoolInfo)
	if err != nil {
		return nil, err
	}

	lggr.Debugw("Decoded tokenPoolInfo", "tokenPoolInfo", tokenPoolInfo)

	tokenPools := make([]TokenPool, len(tokenAmounts))
	for i, tokenAmount := range tokenAmounts {
		lggr.Debugw("\n\nGetting pool address for token",
			"tokenAddress", tokenAmount.DestTokenAddress,
			"poolIndex", i)

		packageId := hex.EncodeToString(tokenPoolInfo.TokenPoolPackageIds[i][:])
		if !strings.HasPrefix(packageId, "0x") {
			packageId = "0x" + packageId
		}

		tokenType := tokenPoolInfo.TokenTypes[i]
		if !strings.HasPrefix(tokenType, "0x") {
			tokenType = "0x" + tokenType
		}

		tokenPoolStateAddress := hex.EncodeToString(tokenPoolInfo.TokenPoolStateAddresses[i][:])
		if !strings.HasPrefix(tokenPoolStateAddress, "0x") {
			tokenPoolStateAddress = "0x" + tokenPoolStateAddress
		}

		tokenPools[i] = TokenPool{
			CoinMetadata:          "0x" + hex.EncodeToString(tokenAmount.DestTokenAddress),
			TokenType:             tokenType,
			PackageId:             packageId,
			ModuleId:              tokenPoolInfo.TokenPoolModules[i],
			Function:              OFFRAMP_TOKEN_POOL_FUNCTION_NAME,
			TokenPoolStateAddress: tokenPoolStateAddress,
			Index:                 i,
		}
	}

	lggr.Debugw("tokenPoolInfo Decoded", "tokenPoolInfo", tokenPools)

	return tokenPools, nil
}

// GetOffRampPTB creates end-to-end OffRamp Execute PTB command given report, ptb configs, and sui client
//
// This function orchestrates the creation of a complete Programmable Transaction Block (PTB) for OffRamp execution
// by combining base commands from configuration with dynamically generated token pool commands. The strategy involves:
//
// 1. Base PTB Commands:
//    - init_execute: Initializes the OffRamp execution
//    - finish_execute: Finalizes the execution
//
// 2. Dynamic Token Pool Commands:
//    - Generated for each token transfer in the messages field from the Execute Report
//    - Inserted between init_execute and finish_execute
//    - Handles token pool operations for each transfer
//
// 3. Message Processing:
//    - Extracts all messages and token amounts from the reports
//    - Looks up token pool information for each token
//    - Generates appropriate PTB commands for token operations
//
// The function ensures proper sequencing of commands and maintains the integrity of the OffRamp execution flow.

func (s *SuiPTBExpander) getOffRampPTB(
	ctx context.Context,
	lggr logger.Logger,
	args SuiOffRampExecCallArgs,
	ptbConfigs *config.ChainWriterFunction,
	signerPublicKey []byte,
) (ptbCommands []config.ChainWriterPTBCommand, updatedArgs map[string]any, typeArgs map[string]string, err error) {
	// update the args with the new values
	// update the hot potato reference index of the last token pool command.

	// We will have 2 base PTB commands from config:
	// 1. init_execute
	// 2. finish_execute

	// We will insert token pool commands between them
	if len(ptbConfigs.PTBCommands) != DEFAULT_NR_OFFRAMP_PTB_COMMANDS {
		return nil, nil, nil, fmt.Errorf("expected %d PTB commands, got %d", DEFAULT_NR_OFFRAMP_PTB_COMMANDS, len(ptbConfigs.PTBCommands))
	}

	tokenAmounts := make([]ccipocr3.RampTokenAmount, 0)
	messages := make([]ccipocr3.Message, 0)

	// save all messages in a single slice
	for _, report := range args.Info.AbstractReports {
		for _, message := range report.Messages {
			tokenAmounts = append(tokenAmounts, message.TokenAmounts...)
			messages = append(messages, message)
		}
	}

	tokenPoolStateAddresses, err := s.getTokenPoolByTokenAddress(ctx, lggr, tokenAmounts, signerPublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	generatedTokenPoolCommands, err := GeneratePTBCommandsForTokenPools(lggr, tokenPoolStateAddresses)
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: filter  out messages that have a receiver that is not registered

	// Generate receiver call commands
	//nolint:gosec // G115:
	receiverCommands, err := GenerateReceiverCallCommands(lggr, messages, uint16(len(generatedTokenPoolCommands)))
	if err != nil {
		return nil, nil, nil, err
	}

	// Construct the final PTB commands by inserting generated commands between config commands
	//nolint:gosec // G115:
	finalPTBCommands := make([]config.ChainWriterPTBCommand, 0, len(ptbConfigs.PTBCommands)+len(generatedTokenPoolCommands)+len(receiverCommands))

	// Add the first command from config (init_execute)
	finalPTBCommands = append(finalPTBCommands, ptbConfigs.PTBCommands[0])

	// Insert all generated token pool commands
	finalPTBCommands = append(finalPTBCommands, generatedTokenPoolCommands...)

	// Insert all generated receiver commands
	finalPTBCommands = append(finalPTBCommands, receiverCommands...)

	// Add the remaining commands from config (finish_execute)
	endCommand := ptbConfigs.PTBCommands[len(ptbConfigs.PTBCommands)-1]

	// Find and update the PTB dependency in the existing parameters
	for i := range endCommand.Params {
		if endCommand.Params[i].PTBDependency != nil {
			//nolint:gosec // G115: PTB commands are typically small in number, overflow extremely unlikely
			endCommand.Params[i].PTBDependency.CommandIndex = uint16(len(finalPTBCommands) - 1)
		}
	}

	finalPTBCommands = append(finalPTBCommands, endCommand)

	// Generate token pool arguments
	tokenPoolArgs, typeArgs, err := GenerateArgumentsForTokenPools(s.AddressMappings["ccipObjectRef"], s.AddressMappings["clockObject"], lggr, tokenPoolStateAddresses)
	if err != nil {
		return nil, nil, nil, err
	}

	filteredMessages, err := s.FilterRegisteredReceivers(ctx, lggr, messages, signerPublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Generate receiver call arguments
	//nolint:gosec // G115:
	receiverArgs, err := GenerateReceiverCallArguments(lggr, filteredMessages, uint16(len(generatedTokenPoolCommands)), s.AddressMappings["ccipObjectRef"])
	if err != nil {
		return nil, nil, nil, err
	}

	// Merge token pool and receiver arguments
	ptbArguments := make(map[string]any)
	for k, v := range tokenPoolArgs {
		ptbArguments[k] = v
	}
	for k, v := range receiverArgs {
		ptbArguments[k] = v
	}

	return finalPTBCommands, ptbArguments, typeArgs, nil
}

func GenerateReceiverCallCommands(
	lggr logger.Logger,
	messages []ccipocr3.Message,
	previousCommandIndex uint16,
) ([]config.ChainWriterPTBCommand, error) {
	var receiverCommands []config.ChainWriterPTBCommand
	receiverIndex := previousCommandIndex + 1
	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			// Parse the receiver string into packageID:moduleID:functionName format
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}

			receiverCommands = append(receiverCommands, config.ChainWriterPTBCommand{
				Type:      codec.SuiPTBCommandMoveCall,
				PackageId: AnyPointer(receiverParts[0]),
				ModuleId:  AnyPointer(receiverParts[1]),
				Function:  AnyPointer(receiverParts[2]),
				Params: []codec.SuiFunctionParam{
					{
						Name:     "ccip_object_ref",
						Type:     "object_id",
						Required: true,
					},
					{
						Name:     fmt.Sprintf("package_id_%d", receiverIndex),
						Type:     "address",
						Required: true,
					},
					{
						Name:     fmt.Sprintf("receiver_params_%d", receiverIndex),
						Type:     "ptb_dependency",
						Required: true,
						PTBDependency: &codec.PTBCommandDependency{
							// PTB commands are typically small in number, overflow extremely unlikely
							//nolint:gosec
							CommandIndex: receiverIndex - 1,
						},
					},
				},
			})
			receiverIndex++
		}
	}

	return receiverCommands, nil
}

func (s *SuiPTBExpander) FilterRegisteredReceivers(
	ctx context.Context,
	lggr logger.Logger,
	messages []ccipocr3.Message,
	signerPublicKey []byte,
) ([]ccipocr3.Message, error) {
	registeredReceivers := make([]ccipocr3.Message, 0)
	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}

			receiverPackageId := receiverParts[0]

			signerAddress, err := client.GetAddressFromPublicKey(signerPublicKey)
			if err != nil {
				return nil, err
			}

			lggr.Debugw("Getting receiver config", "receiverPackageId", receiverPackageId, "receiverFunctionName", receiverParts[2])

			result, err := s.ptbClient.ReadFunction(
				ctx,
				signerAddress,
				s.AddressMappings["ccipPackageId"],
				"receiver_registry",
				"is_registered_receiver",
				[]any{
					s.AddressMappings["ccipObjectRef"],
					receiverPackageId,
				},
				[]string{
					"object_id",
					"address",
				},
			)

			if err != nil {
				lggr.Errorw("Error getting pool infos", "error", err)
				return nil, err
			}

			var isRegistered bool
			lggr.Debugw("isRegistered", "isRegistered", result[0])
			err = codec.DecodeSuiJsonValue(result[0], &isRegistered)
			if err != nil {
				return nil, err
			}

			if isRegistered {
				registeredReceivers = append(registeredReceivers, message)
			}
		}
	}

	lggr.Debugw("registeredReceivers", "registeredReceivers", registeredReceivers)

	return registeredReceivers, nil
}

func GenerateReceiverCallArguments(
	lggr logger.Logger,
	messages []ccipocr3.Message,
	previousCommandIndex uint16,
	ccipObjectRef string,
) (map[string]any, error) {
	arguments := make(map[string]any)

	arguments["ccip_object_ref"] = ccipObjectRef

	commandIndex := previousCommandIndex + 1

	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			lggr.Debugw("receiverParts", "receiverParts", message.Receiver)
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}
			arguments[fmt.Sprintf("package_id_%d", commandIndex)] = receiverParts[0]
			commandIndex++
		}
	}

	return arguments, nil
}

// Auxiliary functions

// GeneratePTBCommands generates PTB commands for token addresses
func GeneratePTBCommandsForTokenPools(
	lggr logger.Logger,
	tokenPools []TokenPool,
) ([]config.ChainWriterPTBCommand, error) {
	ptbCommands := make([]config.ChainWriterPTBCommand, len(tokenPools))
	for i, tokenPool := range tokenPools {
		cmdIndex := i + 1
		previousCommandIndex := i
		// We need to increment the index by 1 because the first command (`init_execute`) is already added
		lggr.Infow("Generating PTB command from token pool", "tokenPool", tokenPool, "with index", i)
		packageID := tokenPool.PackageId
		if !strings.HasPrefix(packageID, "0x") {
			packageID = "0x" + packageID
		}
		ptbCommands[i] = config.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: AnyPointer(packageID),
			ModuleId:  AnyPointer(tokenPool.ModuleId),
			Function:  AnyPointer(tokenPool.Function),
			Params: []codec.SuiFunctionParam{
				{
					Name:      "ccip_object_ref",
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(false),
				},
				{
					Name:     "receiver_params",
					Type:     "ptb_dependency",
					Required: true,
					PTBDependency: &codec.PTBCommandDependency{
						//nolint:gosec // G115: PTB commands are typically small in number, overflow extremely unlikely
						CommandIndex: uint16(previousCommandIndex),
					},
				},
				{
					Name:     fmt.Sprintf("index_%d", cmdIndex),
					Type:     "u64",
					Required: true,
				},
				{
					Name:      fmt.Sprintf("pool_%d", cmdIndex),
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(true),
					IsGeneric: true,
				},
				{
					Name:      "clock",
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(false),
				},
			},
		}
	}

	return ptbCommands, nil
}

// GenerateArgumentsFromTokenAmounts generates PTB arguments for token addresses
//
//nolint:gosec // G115:
func GenerateArgumentsForTokenPools(
	ccipObjectRef string,
	clockRef string,
	lggr logger.Logger,
	tokenPools []TokenPool,
) (map[string]any, map[string]string, error) {
	arguments := make(map[string]any)
	typeArgs := make(map[string]string)

	arguments["ccip_object_ref"] = ccipObjectRef
	arguments["clock"] = clockRef

	for i, tokenPool := range tokenPools {
		cmdIndex := i + 1
		arguments[fmt.Sprintf("pool_%d", cmdIndex)] = tokenPool.TokenPoolStateAddress
		arguments[fmt.Sprintf("index_%d", cmdIndex)] = uint64(tokenPool.Index)
		typeArgs[fmt.Sprintf("pool_%d", cmdIndex)] = tokenPool.TokenType
	}

	return arguments, typeArgs, nil
}

func AnyPointer[T any](v T) *T {
	return &v
}
