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
	ZeroAddress              string        = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

// String returns the string representation of TokenPoolType
func (t TokenPoolType) String() string {
	return string(t)
}

type TokenToolDetails struct {
	TokenPoolPackageId string
	TokenPoolType      TokenPoolType
}

func getCreateTokenTransferParams(ccipOnrampPackageId string) cwConfig.ChainWriterPTBCommand {
	return cwConfig.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: strPtr(ccipOnrampPackageId),
		ModuleId:  strPtr("onramp_state_helper"),
		Function:  strPtr("create_token_transfer_params"),
		Params: []codec.SuiFunctionParam{
			{
				Name:      "token_receiver",
				Type:      "vector<u8>",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
		},
	}
}

func getCCIPSendCommand(ccipOnrampPackageId string, feeTokenType string) cwConfig.ChainWriterPTBCommand {
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
				Name:     "destination_chain_selector",
				Type:     "u64",
				Required: true,
			},
			{
				Name:     "receiver",
				Type:     "vector<u8>",
				Required: true,
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
					CommandIndex: 0,
				},
			},
			{
				Name:      "fee_token_metadata",
				Type:      "object_id",
				Required:  true,
				IsMutable: BoolPointer(false),
			},
			{
				Name:        "fee_token",
				Type:        "object_id",
				Required:    true,
				GenericType: strPtr(feeTokenType),
				IsMutable:   BoolPointer(true),
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
func getLRLockOrBurnCommand(tokenPoolPackageId string, tokenType string) cwConfig.ChainWriterPTBCommand {
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
				Name:     "token_transfer_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: 0,
				},
			},
			{
				Name:        "c_link",
				Type:        "object_id",
				Required:    true,
				GenericType: strPtr(tokenType),
			},
			{
				Name:     "destination_chain_selector",
				Type:     "u64",
				Required: true,
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
func getBMLockOrBurnCommand(tokenPoolPackageId string, ethTokenType string) cwConfig.ChainWriterPTBCommand {
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
				Name:     "token_transfer_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: 0,
				},
			},
			{
				Name:        "c_eth",
				Type:        "object_id",
				Required:    true,
				GenericType: strPtr(ethTokenType),
			},
			{
				Name:     "destination_chain_selector",
				Type:     "u64",
				Required: true,
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
func getManagedLockOrBurnCommand(tokenPoolPackageId string, ethTokenType string) cwConfig.ChainWriterPTBCommand {
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
				Name:     "token_transfer_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: 0,
				},
			},
			{
				Name:        "c_managed_eth",
				Type:        "object_id",
				Required:    true,
				GenericType: strPtr(ethTokenType),
			},
			{
				Name:     "destination_chain_selector",
				Type:     "u64",
				Required: true,
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

// ConfigureOnRampChainWriter creates a single ChainWriterConfig that contains
// two PTB configurations: one for message passing and one for token transfers with messaging
func ConfigureOnRampChainWriter(
	ccipPackageId string,
	ccipOnrampPackageId string,
	tokenPools []TokenToolDetails,
	publicKeyBytes []byte,
	feeTokenType string,
	linkTokenType string,
	ethTokenType string,
) (cwConfig.ChainWriterConfig, error) {
	functions := make(map[string]*cwConfig.ChainWriterFunction)

	// Build PTB for message passing only
	messagePassingCommands := []cwConfig.ChainWriterPTBCommand{
		getCreateTokenTransferParams(ccipPackageId),
		getCCIPSendCommand(ccipOnrampPackageId, feeTokenType),
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
			getCreateTokenTransferParams(ccipPackageId),
		}

		for _, tokenPool := range tokenPools {
			switch tokenPool.TokenPoolType {
			case TokenPoolTypeLockRelease:
				lockOrBurnCommand := getLRLockOrBurnCommand(tokenPool.TokenPoolPackageId, linkTokenType)
				tokenTransferCommands = append(tokenTransferCommands, lockOrBurnCommand)
			case TokenPoolTypeBurnMint:
				burnMintCommand := getBMLockOrBurnCommand(tokenPool.TokenPoolPackageId, ethTokenType)
				tokenTransferCommands = append(tokenTransferCommands, burnMintCommand)
			case TokenPoolTypeManaged:
				managedCommand := getManagedLockOrBurnCommand(tokenPool.TokenPoolPackageId, ethTokenType)
				tokenTransferCommands = append(tokenTransferCommands, managedCommand)
			case TokenPoolTypeUSDC:
				// TODO: Add USDC token pool command when available
				return cwConfig.ChainWriterConfig{}, fmt.Errorf("usdc_token_pool not yet implemented")
			default:
				return cwConfig.ChainWriterConfig{}, fmt.Errorf("unknown token pool type: %s", tokenPool.TokenPoolType)
			}
		}

		ccipSendCommand := getCCIPSendCommand(ccipOnrampPackageId, feeTokenType)
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
