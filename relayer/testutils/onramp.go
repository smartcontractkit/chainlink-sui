package testutils

import (
	"fmt"

	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func strPtr(s string) *string {
	return &s
}

// TokenPoolType represents the different types of token pools available
type TokenPoolType string

const (
	TokenPoolTypeLockRelease TokenPoolType = "lock_release_token_pool"
	TokenPoolTypeBurnMint    TokenPoolType = "burn_mint_token_pool"
	TokenPoolTypeManaged     TokenPoolType = "managed_token_pool"
	TokenPoolTypeBase        TokenPoolType = "token_pool"
	TokenPoolTypeUSDC        TokenPoolType = "usdc_token_pool"
)

// String returns the string representation of TokenPoolType
func (t TokenPoolType) String() string {
	return string(t)
}

type TokenToolDetails struct {
	TokenPoolPackageId string
	TokenPoolType      TokenPoolType
}

func getCCIPSendCommand(ccipOnrampPackageId string, previousCommandIndex int) cwConfig.ChainWriterPTBCommand {
	return cwConfig.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: strPtr(ccipOnrampPackageId),
		ModuleId:  strPtr("onramp"),
		Function:  strPtr("ccip_send"),
		Params: []codec.SuiFunctionParam{
			{
				Name:      "ccip_object_ref_mutable",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(true),
			},
			{
				Name:      "onramp_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(true),
			},
			{
				Name:      "clock",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:     "data",
				Type:     "vector<u8>",
				Required: true,
			},
			{
				Name:     "token_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: uint16(previousCommandIndex),
				},
			},
			{
				Name:      "fee_token_metadata",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:      "fee_token",
				Type:      "object_id",
				Required:  true,
				IsGeneric: true,
				IsMutable: BoolPointer(true),
			},
			{
				Name:     "extra_args",
				Type:     "vector<u8>",
				Required: true,
			},
		},
	}
}

// getLRLockOrBurnCommand returns a ChainWriterPTBCommand for the lock_or_burn function of the lock_release_token_pool module
func getLRLockOrBurnCommand(tokenPoolPackageId string, previousCommandIndex int) cwConfig.ChainWriterPTBCommand {
	return cwConfig.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: strPtr(tokenPoolPackageId),
		ModuleId:  strPtr("lock_release_token_pool"),
		Function:  strPtr("lock_or_burn"),
		Params: []codec.SuiFunctionParam{
			{
				Name:      "ccip_object_ref",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:      "c_link",
				Type:      "object_id",
				Required:  true,
				IsGeneric: true,
			},
			{
				Name:     "token_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: uint16(previousCommandIndex),
				},
			},
			{
				Name:      "clock",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:      "link_lock_release_token_pool_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(true),
			},
		},
	}
}

// getBMLockOrBurnCommand returns a ChainWriterPTBCommand for the lock_or_burn function of the burn_mint_token_pool module
func getBMLockOrBurnCommand(tokenPoolPackageId string, previousCommandIndex int) cwConfig.ChainWriterPTBCommand {
	return cwConfig.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: strPtr(tokenPoolPackageId),
		ModuleId:  strPtr("burn_mint_token_pool"),
		Function:  strPtr("lock_or_burn"),
		Params: []codec.SuiFunctionParam{
			{
				Name:      "ccip_object_ref",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:      "c_eth",
				Type:      "object_id",
				Required:  true,
				IsGeneric: true,
			},
			{
				Name:     "token_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: uint16(previousCommandIndex),
				},
			},
			{
				Name:      "clock",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:      "eth_burn_mint_token_pool_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(true),
			},
		},
	}
}

// getManagedLockOrBurnCommand returns a ChainWriterPTBCommand for the lock_or_burn function of the managed_token_pool module
func getManagedLockOrBurnCommand(tokenPoolPackageId string, previousCommandIndex int) cwConfig.ChainWriterPTBCommand {
	return cwConfig.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: strPtr(tokenPoolPackageId),
		ModuleId:  strPtr("managed_token_pool"),
		Function:  strPtr("lock_or_burn"),
		Params: []codec.SuiFunctionParam{
			{
				Name:      "ccip_object_ref",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:      "c_managed_eth",
				Type:      "object_id",
				Required:  true,
				IsGeneric: true,
			},
			{
				Name:     "token_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: uint16(previousCommandIndex),
				},
			},
			{
				Name:      "clock",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:      "deny_list",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(true),
			},
			{
				Name:      "eth_managed_token_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(true),
			},
			{
				Name:      "eth_managed_token_pool_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(true),
			},
		},
	}
}

func getCreateTokenParamsCommand(ccipPackageId string) cwConfig.ChainWriterPTBCommand {
	return cwConfig.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: strPtr(ccipPackageId),
		ModuleId:  strPtr("dynamic_dispatcher"),
		Function:  strPtr("create_token_params"),
		Params: []codec.SuiFunctionParam{
			{
				Name:     "destination_chain_selector",
				Type:     "u64",
				Required: true,
			},
			{
				Name:     "receiver",
				Type:     "vector<u8>",
				Required: true,
			},
		},
	}
}

// ConfigureOnRampChainWriter creates a single ChainWriterConfig that contains
// two PTB configurations: one for message passing and one for token transfers with messaging
func ConfigureOnRampChainWriter(ccipPackageId string, ccipOnrampPackageId string, tokenPools []TokenToolDetails, publicKeyBytes []byte) (cwConfig.ChainWriterConfig, error) {
	functions := make(map[string]*cwConfig.ChainWriterFunction)

	// Build PTB for message passing only
	messagePassingCommands := []cwConfig.ChainWriterPTBCommand{
		getCreateTokenParamsCommand(ccipPackageId),
		getCCIPSendCommand(ccipOnrampPackageId, 0),
	}

	functions["message_passing"] = &cwConfig.ChainWriterFunction{
		Name:        "message_passing",
		PublicKey:   publicKeyBytes,
		Params:      []codec.SuiFunctionParam{},
		PTBCommands: messagePassingCommands,
	}

	// Build PTB for token transfers with messaging (if token pools are provided)
	if len(tokenPools) > 0 {
		tokenTransferCommands := []cwConfig.ChainWriterPTBCommand{
			getCreateTokenParamsCommand(ccipPackageId),
		}

		currentCommandIndex := 0
		for _, tokenPool := range tokenPools {
			switch tokenPool.TokenPoolType {
			case TokenPoolTypeLockRelease:
				lockOrBurnCommand := getLRLockOrBurnCommand(tokenPool.TokenPoolPackageId, currentCommandIndex)
				tokenTransferCommands = append(tokenTransferCommands, lockOrBurnCommand)
			case TokenPoolTypeBurnMint:
				burnMintCommand := getBMLockOrBurnCommand(tokenPool.TokenPoolPackageId, currentCommandIndex)
				tokenTransferCommands = append(tokenTransferCommands, burnMintCommand)
			case TokenPoolTypeManaged:
				managedCommand := getManagedLockOrBurnCommand(tokenPool.TokenPoolPackageId, currentCommandIndex)
				tokenTransferCommands = append(tokenTransferCommands, managedCommand)
			case TokenPoolTypeUSDC:
				// TODO: Add USDC token pool command when available
				return cwConfig.ChainWriterConfig{}, fmt.Errorf("usdc_token_pool not yet implemented")
			default:
				return cwConfig.ChainWriterConfig{}, fmt.Errorf("unknown token pool type: %s", tokenPool.TokenPoolType)
			}
			currentCommandIndex++
		}

		// Add the final CCIP send command
		// The ccip_send command should depend on the result of the lock_or_burn command
		cmdIndex := len(tokenTransferCommands) - 1
		ccipSendCommand := getCCIPSendCommand(ccipOnrampPackageId, cmdIndex)
		tokenTransferCommands = append(tokenTransferCommands, ccipSendCommand)

		functions["token_transfer_with_messaging"] = &cwConfig.ChainWriterFunction{
			Name:        "token_transfer_with_messaging",
			PublicKey:   publicKeyBytes,
			Params:      []codec.SuiFunctionParam{},
			PTBCommands: tokenTransferCommands,
		}
	}

	return cwConfig.ChainWriterConfig{
		Modules: map[string]*cwConfig.ChainWriterModule{
			cwConfig.PTBChainWriterModuleName: {
				Name:      cwConfig.PTBChainWriterModuleName,
				ModuleID:  "0x123",
				Functions: functions,
			},
		},
	}, nil
}
