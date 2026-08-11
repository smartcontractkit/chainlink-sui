package testutils

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// TokenPoolType represents the different types of token pools available
type TokenPoolType string

const (
	TokenPoolTypeLockRelease TokenPoolType = "lock_release_token_pool"
	TokenPoolTypeBurnMint    TokenPoolType = "burn_mint_token_pool"
	TokenPoolTypeManaged     TokenPoolType = "managed_token_pool"
	TokenPoolTypeBase        TokenPoolType = "token_pool"
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

func getCreateTokenTransferParams(ccipOnrampPackageID string) sui.ChainWriterPTBCommand {
	return sui.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: new(ccipOnrampPackageID),
		ModuleId:  new("onramp_state_helper"),
		Function:  new("create_token_transfer_params"),
		Params: []sui.SuiFunctionParam{
			{
				Name:      "token_receiver",
				Type:      "vector<u8>",
				Required:  true,
				IsMutable: new(false),
			},
		},
	}
}

func getCCIPSendCommand(ccipOnrampPackageID string, feeTokenType string) sui.ChainWriterPTBCommand {
	return sui.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: new(ccipOnrampPackageID),
		ModuleId:  new("onramp"),
		Function:  new("ccip_send"),
		Params: []sui.SuiFunctionParam{
			{
				Name:      "ccip_object_ref_mutable",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(true),
			},
			{
				Name:      "onramp_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(true),
			},
			{
				Name:      "clock",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(false),
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
				PTBDependency: &sui.PTBCommandDependency{
					CommandIndex: 0,
				},
			},
			{
				Name:      "fee_token_metadata",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(false),
			},
			{
				Name:        "fee_token",
				Type:        "object_id",
				Required:    true,
				GenericType: new(feeTokenType),
				IsMutable:   new(true),
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
func getLRLockOrBurnCommand(tokenPoolPackageID string, tokenType string) sui.ChainWriterPTBCommand {
	return sui.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: new(tokenPoolPackageID),
		ModuleId:  new("lock_release_token_pool"),
		Function:  new("lock_or_burn"),
		Params: []sui.SuiFunctionParam{
			{
				Name:      "ccip_object_ref",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(false),
			},
			{
				Name:     "token_transfer_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &sui.PTBCommandDependency{
					CommandIndex: 0,
				},
			},
			{
				Name:        "c_link",
				Type:        "object_id",
				Required:    true,
				GenericType: new(tokenType),
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
				IsMutable: new(false),
			},
			{
				Name:      "link_lock_release_token_pool_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(true),
			},
		},
	}
}

// getBMLockOrBurnCommand returns a ChainWriterPTBCommand for the lock_or_burn function of the burn_mint_token_pool module
func getBMLockOrBurnCommand(tokenPoolPackageID string, ethTokenType string) sui.ChainWriterPTBCommand {
	return sui.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: new(tokenPoolPackageID),
		ModuleId:  new("burn_mint_token_pool"),
		Function:  new("lock_or_burn"),
		Params: []sui.SuiFunctionParam{
			{
				Name:      "ccip_object_ref",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(false),
			},
			{
				Name:     "token_transfer_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &sui.PTBCommandDependency{
					CommandIndex: 0,
				},
			},
			{
				Name:        "c_eth",
				Type:        "object_id",
				Required:    true,
				GenericType: new(ethTokenType),
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
				IsMutable: new(false),
			},
			{
				Name:      "eth_burn_mint_token_pool_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(true),
			},
		},
	}
}

// getManagedLockOrBurnCommand returns a ChainWriterPTBCommand for the lock_or_burn function of the managed_token_pool module
func getManagedLockOrBurnCommand(tokenPoolPackageID string, ethTokenType string) sui.ChainWriterPTBCommand {
	return sui.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: new(tokenPoolPackageID),
		ModuleId:  new("managed_token_pool"),
		Function:  new("lock_or_burn"),
		Params: []sui.SuiFunctionParam{
			{
				Name:      "ccip_object_ref",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(false),
			},
			{
				Name:     "token_transfer_params",
				Type:     "ptb_dependency",
				Required: true,
				PTBDependency: &sui.PTBCommandDependency{
					CommandIndex: 0,
				},
			},
			{
				Name:        "c_managed_eth",
				Type:        "object_id",
				Required:    true,
				GenericType: new(ethTokenType),
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
				IsMutable: new(false),
			},
			{
				Name:      "deny_list",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(true),
			},
			{
				Name:      "eth_managed_token_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(true),
			},
			{
				Name:      "eth_managed_token_pool_state",
				Type:      "object_id",
				Required:  true,
				IsMutable: new(true),
			},
		},
	}
}

// ConfigureOnRampChainWriter creates a single ChainWriterConfig that contains
// two PTB configurations: one for message passing and one for token transfers with messaging
func ConfigureOnRampChainWriter(
	log logger.Logger,
	ccipPackageId string,
	ccipOnrampPackageId string,
	tokenPools []TokenToolDetails,
	publicKeyBytes []byte,
	feeTokenType string,
	linkTokenType string,
	ethTokenType string,
) (sui.ChainWriterConfig, error) {
	functions := make(map[string]*sui.ChainWriterFunction)

	// Build PTB for message passing only
	messagePassingCommands := []sui.ChainWriterPTBCommand{
		getCreateTokenTransferParams(ccipPackageId),
		getCCIPSendCommand(ccipOnrampPackageId, feeTokenType),
	}

	functions["message_passing"] = &sui.ChainWriterFunction{
		Name:        "message_passing",
		PublicKey:   publicKeyBytes,
		Params:      []sui.SuiFunctionParam{},
		PTBCommands: messagePassingCommands,
	}

	// Build PTB for token transfers with messaging (if token pools are provided)
	if len(tokenPools) > 0 {
		tokenTransferCommands := []sui.ChainWriterPTBCommand{
			getCreateTokenTransferParams(ccipPackageId),
		}

		for _, tokenPool := range tokenPools {
			switch tokenPool.TokenPoolType {
			case TokenPoolTypeLockRelease:
				lockOrBurnCommand := getLRLockOrBurnCommand(tokenPool.TokenPoolPackageId, linkTokenType)
				tokenTransferCommands = append(tokenTransferCommands, lockOrBurnCommand)
			case TokenPoolTypeBurnMint:
				burnMintCommand := getBMLockOrBurnCommand(tokenPool.TokenPoolPackageId, ethTokenType)
				log.Debugw("burnMintCommand", "burnMintCommand", burnMintCommand)
				tokenTransferCommands = append(tokenTransferCommands, burnMintCommand)
			case TokenPoolTypeManaged:
				managedCommand := getManagedLockOrBurnCommand(tokenPool.TokenPoolPackageId, ethTokenType)
				tokenTransferCommands = append(tokenTransferCommands, managedCommand)
			default:
				return sui.ChainWriterConfig{}, fmt.Errorf("unknown token pool type: %s", tokenPool.TokenPoolType)
			}
		}

		ccipSendCommand := getCCIPSendCommand(ccipOnrampPackageId, feeTokenType)
		tokenTransferCommands = append(tokenTransferCommands, ccipSendCommand)

		functions["token_transfer_with_messaging"] = &sui.ChainWriterFunction{
			Name:        "token_transfer_with_messaging",
			PublicKey:   publicKeyBytes,
			Params:      []sui.SuiFunctionParam{},
			PTBCommands: tokenTransferCommands,
		}
	}

	return sui.ChainWriterConfig{
		Modules: map[string]*sui.ChainWriterModule{
			cwConfig.PTBChainWriterModuleName: {
				Name:      cwConfig.PTBChainWriterModuleName,
				ModuleID:  "0x123",
				Functions: functions,
			},
		},
	}, nil
}
