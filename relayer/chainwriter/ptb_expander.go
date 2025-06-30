package chainwriter

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

const (
	DEFAULT_NR_OFFRAMP_PTB_COMMANDS  = 2
	OFFRAMP_TOKEN_POOL_FUNCTION_NAME = "release_or_mint"
	SUI_PATH_COMPONENTS_COUNT        = 3
)

type SuiAddress [32]byte

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

// PTBExpander defines the interface for expanding PTB (Programmable Transaction Block) commands
// for OffRamp execution in the Sui blockchain.
//
// This interface provides methods to generate PTB commands and arguments needed for CCIP
// (Cross-Chain Interoperability Protocol) message execution on Sui. It handles token pool
// operations, receiver contract calls, and the complete OffRamp execution flow.
//
// The PTB expansion process involves:
// 1. Looking up token pool addresses for token transfers
// 2. Generating PTB commands for token pool operations
// 3. Creating receiver call commands when messages have receivers
// 4. Orchestrating the complete OffRamp execution PTB
//
// Implementations of this interface should handle the complex task of translating high-level
// CCIP execution requests into low-level Sui Move function calls within a PTB structure.
type PTBExpander interface {
	// GetTokenPoolByTokenAddress retrieves token pool information for given token amounts.
	//
	// This method queries the token admin registry to find the corresponding token pools
	// for each token address in the provided token amounts. It returns detailed information
	// about each token pool including package ID, module ID, function name, and state address.
	//
	// The lookup process involves:
	// 1. Querying the token admin registry with token addresses
	// 2. Resolving token pool state addresses from state pointer objects
	// 3. Gathering package and module information for PTB command generation
	//
	// Parameters:
	//   - lggr: Logger instance for debugging and tracing operations
	//   - tokenAmounts: Slice of RampTokenAmount containing token addresses and transfer amounts
	//
	// Returns:
	//   - []TokenPool: Slice of TokenPool structs containing comprehensive pool information
	//   - error: Error if token pool lookup fails, token is not registered, or network issues occur
	//
	// Usage:
	//   This method is typically called during PTB expansion to resolve token addresses
	//   to their corresponding pool contracts before generating PTB commands.
	//
	// Example:
	//   tokenPools, err := expander.GetTokenPoolByTokenAddress(logger, tokenAmounts)
	//   if err != nil {
	//       return fmt.Errorf("failed to get token pools: %w", err)
	//   }
	GetTokenPoolByTokenAddress(
		lggr logger.Logger,
		tokenAmounts []ccipocr3.RampTokenAmount,
		signerPublicKey []byte,
	) ([]TokenPool, error)

	// GetOffRampPTB creates the complete end-to-end OffRamp Execute PTB command sequence and updates execution arguments.
	//
	// This method orchestrates the entire PTB creation process for executing CCIP messages
	// on the destination chain. It processes execution reports, generates token pool commands,
	// creates receiver call commands, assembles the final PTB with proper dependencies, and
	// generates the corresponding argument map with resolved values.
	//
	// The generated PTB typically includes:
	// 1. init_execute command to start the execution process and validate the report
	// 2. Token pool release_or_mint commands for each token transfer in the messages
	// 3. Receiver call commands for messages that specify receiver contracts
	// 4. finish_execute command to complete the execution and emit events
	//
	// The method handles both command generation and argument resolution:
	// - PTB Command Generation: Creates the sequence of Move function calls
	// - Argument Updates: Resolves object IDs, amounts, addresses, and dependencies
	// - Report validation and parsing
	// - Token pool address resolution
	// - PTB command dependency management
	// - Error handling and rollback scenarios
	//
	// Parameters:
	//   - lggr: Logger instance for debugging and tracing the entire execution flow
	//   - args: SuiOffRampExecCallArgs containing the execution report, context, and metadata
	//   - ptbConfigs: ChainWriterFunction configuration containing PTB command templates and settings
	//
	// Returns:
	//   - ptbCommands: Slice of ChainWriterPTBCommand representing the complete executable PTB
	//   - updatedArgs: Map of arguments with resolved values for PTB execution (object IDs, amounts, addresses, etc.)
	//   - err: Error if PTB generation or argument resolution fails at any step
	//
	// Usage:
	//   This is the main entry point for PTB expansion, called by the chain writer
	//   when executing CCIP messages on Sui. The returned PTB commands and updated arguments
	//   can be directly submitted to the Sui network for execution.
	//
	// Example:
	//   ptbCmds, updatedArgs, err := expander.GetOffRampPTB(logger, execArgs, config)
	//   if err != nil {
	//       return fmt.Errorf("failed to generate OffRamp PTB: %w", err)
	//   }
	//   // Submit ptbCmds with updatedArgs to Sui network
	GetOffRampPTB(
		lggr logger.Logger,
		args SuiOffRampExecCallArgs,
		ptbConfigs *ChainWriterFunction,
		signerPublicKey []byte,
	) (ptbCommands []ChainWriterPTBCommand, updatedArgs any, err error)
}

// SuiPTBExpander is a concrete implementation of the PTBExpander interface for Sui blockchain
type SuiPTBExpander struct {
	lggr            logger.Logger
	ptbClient       client.SuiPTBClient
	AddressMappings map[string]string
}

// NewSuiPTBExpander creates a new instance of SuiPTBExpander
func NewSuiPTBExpander(lggr logger.Logger, ptbClient client.SuiPTBClient, chainWriterConfig ChainWriterConfig) *SuiPTBExpander {
	addressMappings := chainWriterConfig.Modules[PTBChainWriterModuleName].Functions[CCIPExecuteReportFunctionName].AddressMappings
	return &SuiPTBExpander{
		lggr:            lggr,
		ptbClient:       ptbClient,
		AddressMappings: addressMappings,
	}
}

type GetPoolInfosResult struct {
	TokenPoolPackageIds     []SuiAddress `json:"token_pool_package_ids"`
	TokenPoolStateAddresses []SuiAddress `json:"token_pool_state_addresses"`
	TokenPoolModules        []SuiAddress `json:"token_pool_modules"`
	TokenTypes              []string     `json:"token_types"`
}

// GetTokenPoolByTokenAddress gets token pool addresses for given token addresses
func (s *SuiPTBExpander) GetTokenPoolByTokenAddress(
	lggr logger.Logger,
	tokenAmounts []ccipocr3.RampTokenAmount,
	signerPublicKey []byte,
) ([]TokenPool, error) {
	coinMetadataAddresses := make([][]byte, len(tokenAmounts))
	for i, tokenAmount := range tokenAmounts {
		coinMetadataAddresses[i] = []byte(tokenAmount.DestTokenAddress)
	}

	signerAddress, err := client.GetAddressFromPublicKey(signerPublicKey)
	if err != nil {
		return nil, err
	}

	poolInfoResponseValues, err := s.ptbClient.ReadFunction(
		context.Background(),
		signerAddress,
		s.AddressMappings["ccipPackageId"],
		"token_admin_registry",
		"get_pool_infos",
		[]any{
			s.AddressMappings["ccipObjectRef"],
			coinMetadataAddresses,
		},
		[]string{
			"object_id",
			"vector<address>",
		},
	)

	if err != nil {
		lggr.Errorw("Error getting pool infos", "error", err)
		return nil, err
	}

	var tokenPoolInfo GetPoolInfosResult
	lggr.Debugw("tokenPoolInfo", "tokenPoolInfo", poolInfoResponseValues[0])
	err = codec.DecodeSuiJsonValue(poolInfoResponseValues[0], &tokenPoolInfo)
	if err != nil {
		return nil, err
	}

	lggr.Debugw("tokenPoolInfo Decoded", "tokenPoolInfo", tokenPoolInfo)

	tokenPools := make([]TokenPool, len(tokenAmounts))
	for i, tokenAmount := range tokenAmounts {
		lggr.Debugw("Getting pool address for token",
			"tokenAddress", tokenAmount.DestTokenAddress,
			"poolIndex", i)

		tokenPools[i] = TokenPool{
			CoinMetadata:          string(tokenAmount.DestTokenAddress),
			TokenType:             tokenPoolInfo.TokenTypes[i],
			PackageId:             hex.EncodeToString(tokenPoolInfo.TokenPoolPackageIds[i][:]),
			ModuleId:              hex.EncodeToString(tokenPoolInfo.TokenPoolModules[i][:]),
			Function:              OFFRAMP_TOKEN_POOL_FUNCTION_NAME,
			TokenPoolStateAddress: hex.EncodeToString(tokenPoolInfo.TokenPoolStateAddresses[i][:]),
			Index:                 i,
		}
	}

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

func (s *SuiPTBExpander) GetOffRampPTB(
	lggr logger.Logger,
	args SuiOffRampExecCallArgs,
	ptbConfigs *ChainWriterFunction,
	signerPublicKey []byte,
) (ptbCommands []ChainWriterPTBCommand, updatedArgs any, err error) {
	// update the args with the new values
	// update the hot potato reference index of the last token pool command.

	// We will have 2 base PTB commands from config:
	// 1. init_execute
	// 2. finish_execute

	// We will insert token pool commands between them
	if len(ptbConfigs.PTBCommands) != DEFAULT_NR_OFFRAMP_PTB_COMMANDS {
		return nil, nil, fmt.Errorf("expected %d PTB commands, got %d", DEFAULT_NR_OFFRAMP_PTB_COMMANDS, len(ptbConfigs.PTBCommands))
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

	tokenPoolStateAddresses, err := s.GetTokenPoolByTokenAddress(lggr, tokenAmounts, signerPublicKey)
	if err != nil {
		return nil, nil, err
	}

	generatedTokenPoolCommands, err := GeneratePTBCommandsForTokenPools(lggr, tokenPoolStateAddresses)
	if err != nil {
		return nil, nil, err
	}

	// TODO: filter  out messages that have a receiver that is not registered

	// Generate receiver call commands
	//nolint:gosec // G115:
	receiverCommands, err := GenerateReceiverCallCommands(lggr, messages, uint16(len(generatedTokenPoolCommands)))
	if err != nil {
		return nil, nil, err
	}

	// Construct the final PTB commands by inserting generated commands between config commands
	//nolint:gosec // G115:
	finalPTBCommands := make([]ChainWriterPTBCommand, 0, len(ptbConfigs.PTBCommands)+len(generatedTokenPoolCommands)+len(receiverCommands))

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
	tokenPoolArgs, err := GenerateArgumentsForTokenPools(s.AddressMappings["ccipObjectRef"], s.AddressMappings["clockObject"], lggr, tokenPoolStateAddresses)
	if err != nil {
		return nil, nil, err
	}

	filteredMessages, err := s.FilterRegisteredReceivers(lggr, messages, signerPublicKey)
	if err != nil {
		return nil, nil, err
	}

	// Generate receiver call arguments
	//nolint:gosec // G115:
	receiverArgs, err := GenerateReceiverCallArguments(lggr, filteredMessages, uint16(len(generatedTokenPoolCommands)), s.AddressMappings["ccipObjectRef"])
	if err != nil {
		return nil, nil, err
	}

	// Merge token pool and receiver arguments
	ptbArguments := make(map[string]any)
	for k, v := range tokenPoolArgs {
		ptbArguments[k] = v
	}
	for k, v := range receiverArgs {
		ptbArguments[k] = v
	}

	return finalPTBCommands, ptbArguments, nil
}

func GenerateReceiverCallCommands(
	lggr logger.Logger,
	messages []ccipocr3.Message,
	previousCommandIndex uint16,
) ([]ChainWriterPTBCommand, error) {
	var receiverCommands []ChainWriterPTBCommand
	receiverIndex := previousCommandIndex + 1
	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			// Parse the receiver string into packageID:moduleID:functionName format
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}

			receiverCommands = append(receiverCommands, ChainWriterPTBCommand{
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
			receiverModuleId := receiverParts[1]
			receiverFunctionName := receiverParts[2]

			receiverAddress := fmt.Sprintf("%s::%s::%s", receiverPackageId, receiverModuleId, receiverFunctionName)

			signerAddress, err := client.GetAddressFromPublicKey(signerPublicKey)
			if err != nil {
				return nil, err
			}

			lggr.Debugw("Getting receiver config", "receiverAddress", receiverAddress)

			isRegisteredResponseValues, err := s.ptbClient.ReadFunction(
				context.Background(),
				signerAddress,
				s.AddressMappings["ccipPackageId"],
				"ccip",
				"is_registered_receiver",
				[]any{
					s.AddressMappings["ccipObjectRef"],
					receiverAddress,
				},
				[]string{
					"object_id",
					"address",
				},
			)
			isRegistered := isRegisteredResponseValues[0].(bool)

			if err != nil {
				lggr.Errorw("Error getting pool infos", "error", err)
				return nil, err
			}
			lggr.Debugw("isRegistered", "isRegistered", isRegistered)

			if isRegistered {
				registeredReceivers = append(registeredReceivers, message)
			}
		}
	}

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
) ([]ChainWriterPTBCommand, error) {
	ptbCommands := make([]ChainWriterPTBCommand, len(tokenPools))
	for i, tokenPool := range tokenPools {
		cmdIndex := i + 1
		previousCommandIndex := i
		// We need to increment the index by 1 because the first command (`init_execute`) is already added
		lggr.Infow("Generating PTB command from token pool", "tokenPool", tokenPool, "with index", i)
		ptbCommands[i] = ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: AnyPointer(tokenPool.PackageId),
			ModuleId:  AnyPointer(tokenPool.ModuleId),
			Function:  AnyPointer(tokenPool.Function),
			Params: []codec.SuiFunctionParam{
				{
					Name:     fmt.Sprintf("ref_%d", cmdIndex),
					Type:     "object_id",
					Required: true,
				},
				{
					Name:      "clock",
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(false),
				},
				{
					Name:      fmt.Sprintf("pool_%d", cmdIndex),
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(true),
				},
				{
					Name:     "remote_chain_selector",
					Type:     "u64",
					Required: true,
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
			},
		}
	}

	return ptbCommands, nil
}

// GenerateArgumentsFromTokenAmounts generates PTB arguments for token addresses
func GenerateArgumentsForTokenPools(
	ccipStateRef string,
	clockRef string,
	lggr logger.Logger,
	tokenPools []TokenPool,
) (map[string]any, error) {
	arguments := make(map[string]any)

	arguments["ccip_state_ref"] = ccipStateRef
	arguments["clock_ref"] = clockRef
	arguments["remote_chain_selector"] = 0

	for i, tokenPool := range tokenPools {
		cmdIndex := i + 1
		arguments[fmt.Sprintf("pool_%d", cmdIndex)] = tokenPool.PackageId
		arguments[fmt.Sprintf("index_%d", cmdIndex)] = tokenPool.Index
	}

	return arguments, nil
}

func AnyPointer[T any](v T) *T {
	return &v
}
