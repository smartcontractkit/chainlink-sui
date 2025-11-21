package view

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_nonce_manager "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/nonce_manager"
	module_receiver_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/receiver_registry"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
)

type CCIPView struct {
	ContractMetaData

	FeeQuoter          FeeQuoterView          `json:"feeQuoter"`
	RMNRemote          RMNRemoteView          `json:"rmnRemote"`
	TokenAdminRegistry TokenAdminRegistryView `json:"tokenAdminRegistry"`
	NonceManager       NonceManagerView       `json:"nonceManager"`
	ReceiverRegistry   ReceiverRegistryView   `json:"receiverRegistry"`
}

type FeeQuoterView struct {
	ContractMetaData

	FeeTokens               []string                            `json:"feeTokens"`
	StaticConfig            FeeQuoterStaticConfig               `json:"staticConfig"`
	DestinationChainConfigs map[uint64]FeeQuoterDestChainConfig `json:"destinationChainConfigs"`
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg            string `json:"maxFeeJuelsPerMsg"`
	LinkToken                    string `json:"linkToken"`
	TokenPriceStalenessThreshold uint64 `json:"tokenPriceStalenessThreshold"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool   `json:"isEnabled"`
	MaxNumberOfTokensPerMsg           uint16 `json:"maxNumberOfTokensPerMsg"`
	MaxDataBytes                      uint32 `json:"maxDataBytes"`
	MaxPerMsgGasLimit                 uint32 `json:"maxPerMsgGasLimit"`
	DestGasOverhead                   uint32 `json:"destGasOverhead"`
	DestGasPerPayloadByteBase         uint8  `json:"destGasPerPayloadByteBase"`
	DestGasPerPayloadByteHigh         uint8  `json:"destGasPerPayloadByteHigh"`
	DestGasPerPayloadByteThreshold    uint16 `json:"destGasPerPayloadByteThreshold"`
	DestDataAvailabilityOverheadGas   uint32 `json:"destDataAvailabilityOverheadGas"`
	DestGasPerDataAvailabilityByte    uint16 `json:"destGasPerDataAvailabilityByte"`
	DestDataAvailabilityMultiplierBps uint16 `json:"destDataAvailabilityMultiplierBps"`
	ChainFamilySelector               string `json:"chainFamilySelector"`
	EnforceOutOfOrder                 bool   `json:"enforceOutOfOrder"`
	DefaultTokenFeeUsdCents           uint16 `json:"defaultTokenFeeUsdCents"`
	DefaultTokenDestGasOverhead       uint32 `json:"defaultTokenDestGasOverhead"`
	DefaultTxGasLimit                 uint32 `json:"defaultTxGasLimit"`
	GasMultiplierWeiPerEth            uint64 `json:"gasMultiplierWeiPerEth"`
	GasPriceStalenessThreshold        uint32 `json:"gasPriceStalenessThreshold"`
	NetworkFeeUsdCents                uint32 `json:"networkFeeUsdCents"`
}

type RMNRemoteView struct {
	ContractMetaData
	IsCursed             bool                     `json:"isCursed"`
	Config               RMNRemoteVersionedConfig `json:"config"`
	CursedSubjectEntries []RMNRemoteCurseEntry    `json:"cursedSubjectEntries"`
}

type RMNRemoteVersionedConfig struct {
	Version uint32            `json:"version"`
	Signers []RMNRemoteSigner `json:"signers"`
	Fsign   uint64            `json:"fSign"`
}

type RMNRemoteSigner struct {
	OnchainPublicKey string `json:"onchain_public_key"` // Follow EVM snake_case
	NodeIndex        uint64 `json:"node_index"`
}

type RMNRemoteCurseEntry struct {
	Subject  string `json:"subject"`
	Selector uint64 `json:"selector"`
}

type TokenAdminRegistryView struct {
	ContractMetaData
	TokenConfigs map[string]TokenConfigView `json:"tokenConfigs"` // Token address => config
}

type TokenConfigView struct {
	TokenPoolPackageId  string   `json:"tokenPoolPackageId"`
	TokenPoolModule     string   `json:"tokenPoolModule"`
	TokenType           string   `json:"tokenType"`
	Administrator       string   `json:"administrator"`
	TokenPoolTypeProof  string   `json:"tokenPoolTypeProof"`
	LockOrBurnParams    []string `json:"lockOrBurnParams"`
	ReleaseOrMintParams []string `json:"releaseOrMintParams"`
}

type NonceManagerView struct {
	ContractMetaData
}

type ReceiverRegistryView struct {
	ContractMetaData
}

// GenerateCCIPView generates a CCIP view by querying the on-chain state
func GenerateCCIPView(
	ctx context.Context,
	chain sui.Chain,
	ccipPackageID string,
	ccipObjectRef string,
	routerPackageID string,
	routerStateObjectID string,
) (CCIPView, error) {
	if ccipPackageID == "" || ccipObjectRef == "" {
		return CCIPView{}, fmt.Errorf("ccipPackageID and ccipObjectRef cannot be empty")
	}
	ccipRefObj := bind.Object{Id: ccipObjectRef}
	callOpts := &bind.CallOpts{Signer: chain.Signer}

	// Create state object contract binding to get owner
	stateObjectContract, err := module_state_object.NewStateObject(ccipPackageID, chain.Client)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to create state object contract binding: %w", err)
	}

	// Get owner
	owner, err := stateObjectContract.DevInspect().Owner(ctx, callOpts, ccipRefObj)
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to get owner: %w", err)
	}

	// Generate FeeQuoter view
	feeQuoterView, err := generateFeeQuoterView(ctx, chain, ccipPackageID, ccipRefObj, callOpts, routerPackageID, routerStateObjectID)
	feeQuoterView.ContractMetaData.Owner = owner
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to generate fee quoter view: %w", err)
	}

	// Generate RMNRemote view
	rmnRemoteView, err := generateRMNRemoteView(ctx, chain, ccipPackageID, ccipRefObj, callOpts)
	rmnRemoteView.ContractMetaData.Owner = owner
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to generate rmn remote view: %w", err)
	}

	// Generate Token Admin Registry view
	tokenAdminRegistryView, err := generateTokenAdminRegistryView(ctx, chain, ccipPackageID, ccipRefObj, callOpts)
	tokenAdminRegistryView.ContractMetaData.Owner = owner
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to generate token admin registry view: %w", err)
	}

	// Generate Nonce Manager view
	nonceManagerView, err := generateNonceManagerView(ctx, chain, ccipPackageID, callOpts)
	nonceManagerView.ContractMetaData.Owner = owner
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to generate nonce manager view: %w", err)
	}

	// Generate Receiver Registry view
	receiverRegistryView, err := generateReceiverRegistryView(ctx, chain, ccipPackageID, callOpts)
	receiverRegistryView.ContractMetaData.Owner = owner
	if err != nil {
		return CCIPView{}, fmt.Errorf("failed to generate receiver registry view: %w", err)
	}

	return CCIPView{
		ContractMetaData: ContractMetaData{
			Address:       ccipPackageID,
			Owner:         owner,
			StateObjectID: ccipObjectRef,
		},
		FeeQuoter:          feeQuoterView,
		RMNRemote:          rmnRemoteView,
		TokenAdminRegistry: tokenAdminRegistryView,
		NonceManager:       nonceManagerView,
		ReceiverRegistry:   receiverRegistryView,
	}, nil
}

func generateFeeQuoterView(
	ctx context.Context,
	chain sui.Chain,
	ccipPackageID string,
	ccipRefObj bind.Object,
	callOpts *bind.CallOpts,
	routerPackageID string,
	routerStateObjectID string,
) (FeeQuoterView, error) {
	// Create fee quoter contract binding
	feeQuoterContract, err := module_fee_quoter.NewFeeQuoter(ccipPackageID, chain.Client)
	if err != nil {
		return FeeQuoterView{}, fmt.Errorf("failed to create fee quoter contract binding: %w", err)
	}

	// Get type and version
	typeAndVersion, err := feeQuoterContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return FeeQuoterView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	// Get fee tokens
	feeTokens, err := feeQuoterContract.DevInspect().GetFeeTokens(ctx, callOpts, ccipRefObj)
	if err != nil {
		return FeeQuoterView{}, fmt.Errorf("failed to get fee tokens: %w", err)
	}

	// Get static config
	staticConfig, err := feeQuoterContract.DevInspect().GetStaticConfig(ctx, callOpts, ccipRefObj)
	if err != nil {
		return FeeQuoterView{}, fmt.Errorf("failed to get static config: %w", err)
	}

	// TODO: Changesets are not configuring router, any configuration that requires GetDestChains will be empty
	// Get destination chains from router to query per-chain configs
	var destChainSelectors []uint64
	if routerPackageID != "" && routerStateObjectID != "" {
		routerContract, err := module_router.NewRouter(routerPackageID, chain.Client)
		if err != nil {
			return FeeQuoterView{}, fmt.Errorf("failed to create router contract binding: %w", err)
		}

		routerStateObj := bind.Object{Id: routerStateObjectID}
		destChainSelectors, err = routerContract.DevInspect().GetDestChains(ctx, callOpts, routerStateObj)
		if err != nil {
			return FeeQuoterView{}, fmt.Errorf("failed to get dest chains from router: %w", err)
		}
	}

	// Get destination chain configs for each destination chain
	destinationChainConfigs := make(map[uint64]FeeQuoterDestChainConfig)
	for _, destChainSelector := range destChainSelectors {
		destChainConfig, err := feeQuoterContract.DevInspect().GetDestChainConfig(ctx, callOpts, ccipRefObj, destChainSelector)
		if err != nil {
			// Chain might not be configured, skip it
			continue
		}

		destinationChainConfigs[destChainSelector] = FeeQuoterDestChainConfig{
			IsEnabled:                         destChainConfig.IsEnabled,
			MaxNumberOfTokensPerMsg:           destChainConfig.MaxNumberOfTokensPerMsg,
			MaxDataBytes:                      destChainConfig.MaxDataBytes,
			MaxPerMsgGasLimit:                 destChainConfig.MaxPerMsgGasLimit,
			DestGasOverhead:                   destChainConfig.DestGasOverhead,
			DestGasPerPayloadByteBase:         destChainConfig.DestGasPerPayloadByteBase,
			DestGasPerPayloadByteHigh:         destChainConfig.DestGasPerPayloadByteHigh,
			DestGasPerPayloadByteThreshold:    destChainConfig.DestGasPerPayloadByteThreshold,
			DestDataAvailabilityOverheadGas:   destChainConfig.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    destChainConfig.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: destChainConfig.DestDataAvailabilityMultiplierBps,
			ChainFamilySelector:               hex.EncodeToString(destChainConfig.ChainFamilySelector),
			EnforceOutOfOrder:                 destChainConfig.EnforceOutOfOrder,
			DefaultTokenFeeUsdCents:           destChainConfig.DefaultTokenFeeUsdCents,
			DefaultTokenDestGasOverhead:       destChainConfig.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:                 destChainConfig.DefaultTxGasLimit,
			GasMultiplierWeiPerEth:            destChainConfig.GasMultiplierWeiPerEth,
			GasPriceStalenessThreshold:        destChainConfig.GasPriceStalenessThreshold,
			NetworkFeeUsdCents:                destChainConfig.NetworkFeeUsdCents,
		}
	}

	return FeeQuoterView{
		ContractMetaData: ContractMetaData{
			Address:        ccipPackageID,
			Owner:          "",
			TypeAndVersion: typeAndVersion,
			StateObjectID:  ccipRefObj.Id,
		},
		FeeTokens: feeTokens,
		StaticConfig: FeeQuoterStaticConfig{
			MaxFeeJuelsPerMsg:            staticConfig.MaxFeeJuelsPerMsg.String(),
			LinkToken:                    staticConfig.LinkToken,
			TokenPriceStalenessThreshold: staticConfig.TokenPriceStalenessThreshold,
		},
		DestinationChainConfigs: destinationChainConfigs,
	}, nil
}

func generateRMNRemoteView(
	ctx context.Context,
	chain sui.Chain,
	ccipPackageID string,
	ccipRefObj bind.Object,
	callOpts *bind.CallOpts,
) (RMNRemoteView, error) {
	// Create RMN remote contract binding
	rmnRemoteContract, err := module_rmn_remote.NewRmnRemote(ccipPackageID, chain.Client)
	if err != nil {
		return RMNRemoteView{}, fmt.Errorf("failed to create rmn remote contract binding: %w", err)
	}

	// Get type and version
	typeAndVersion, err := rmnRemoteContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return RMNRemoteView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	// Get cursed status (global)
	isCursed, err := rmnRemoteContract.DevInspect().IsCursedGlobal(ctx, callOpts, ccipRefObj)
	if err != nil {
		return RMNRemoteView{}, fmt.Errorf("failed to get cursed status: %w", err)
	}

	// Get versioned config
	// Returns [0]: u32 (version), [1]: Config struct (with signers and f_sign)
	versionedConfigRaw, err := rmnRemoteContract.DevInspect().GetVersionedConfig(ctx, callOpts, ccipRefObj)
	if err != nil {
		return RMNRemoteView{}, fmt.Errorf("failed to get versioned config: %w", err)
	}

	var versionedConfig RMNRemoteVersionedConfig
	if len(versionedConfigRaw) >= 2 {
		version, ok := versionedConfigRaw[0].(uint32)
		if !ok {
			return RMNRemoteView{}, fmt.Errorf("unexpected type for version: got %T", versionedConfigRaw[0])
		}
		versionedConfig.Version = version

		// Parse the config struct (contains signers and f_sign)
		// Config is a struct with fields: signers (vector of Signer structs) and f_sign (u64)
		configData, ok := versionedConfigRaw[1].([]interface{})
		if ok && len(configData) >= 2 {
			// Parse signers
			if signersRaw, ok := configData[0].([]interface{}); ok {
				signers := make([]RMNRemoteSigner, 0, len(signersRaw))
				for _, signerRaw := range signersRaw {
					if signerData, ok := signerRaw.([]interface{}); ok && len(signerData) >= 2 {
						// Signer struct: [onchain_public_key (vector<u8>), node_index (u64)]
						pubKeyBytes, pkOk := signerData[0].([]byte)
						nodeIndex, niOk := signerData[1].(uint64)
						if pkOk && niOk {
							signers = append(signers, RMNRemoteSigner{
								OnchainPublicKey: "0x" + hex.EncodeToString(pubKeyBytes),
								NodeIndex:        nodeIndex,
							})
						}
					}
				}
				versionedConfig.Signers = signers
			}

			// Parse f_sign
			if fSign, ok := configData[1].(uint64); ok {
				versionedConfig.Fsign = fSign
			}
		}
	}

	// Get cursed subjects
	cursedSubjectsRaw, err := rmnRemoteContract.DevInspect().GetCursedSubjects(ctx, callOpts, ccipRefObj)
	if err != nil {
		return RMNRemoteView{}, fmt.Errorf("failed to get cursed subjects: %w", err)
	}

	cursedSubjectEntries := make([]RMNRemoteCurseEntry, 0, len(cursedSubjectsRaw))
	for _, subjectBytes := range cursedSubjectsRaw {
		// Try to parse as chain selector (u128 encoded as bytes)
		var selector uint64
		if len(subjectBytes) == 16 {
			// Convert bytes to u128, then to u64 for chain selector
			selectorBig := new(big.Int).SetBytes(subjectBytes)
			if selectorBig.IsUint64() {
				selector = selectorBig.Uint64()
			}
		}

		cursedSubjectEntries = append(cursedSubjectEntries, RMNRemoteCurseEntry{
			Subject:  "0x" + hex.EncodeToString(subjectBytes),
			Selector: selector,
		})
	}

	return RMNRemoteView{
		ContractMetaData: ContractMetaData{
			Address:        ccipPackageID,
			Owner:          "",
			TypeAndVersion: typeAndVersion,
			StateObjectID:  ccipRefObj.Id,
		},
		IsCursed:             isCursed,
		Config:               versionedConfig,
		CursedSubjectEntries: cursedSubjectEntries,
	}, nil
}

func generateTokenAdminRegistryView(
	ctx context.Context,
	chain sui.Chain,
	ccipPackageID string,
	ccipRefObj bind.Object,
	callOpts *bind.CallOpts,
) (TokenAdminRegistryView, error) {
	// Create token admin registry contract binding
	tokenAdminRegistryContract, err := module_token_admin_registry.NewTokenAdminRegistry(ccipPackageID, chain.Client)
	if err != nil {
		return TokenAdminRegistryView{}, fmt.Errorf("failed to create token admin registry contract binding: %w", err)
	}

	// Get type and version
	typeAndVersion, err := tokenAdminRegistryContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return TokenAdminRegistryView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	// Get all configured tokens
	// GetAllConfiguredTokens returns: [0]: vector<address> (tokens), [1]: address (next_key), [2]: bool (has_next)
	configuredTokensRaw, err := tokenAdminRegistryContract.DevInspect().GetAllConfiguredTokens(ctx, callOpts, ccipRefObj, "0x0", 1000)
	if err != nil {
		return TokenAdminRegistryView{}, fmt.Errorf("failed to get all configured tokens: %w", err)
	}

	tokenConfigs := make(map[string]TokenConfigView)
	if len(configuredTokensRaw) >= 1 {
		if tokens, ok := configuredTokensRaw[0].([]string); ok {
			// Query each token's config
			for _, tokenAddr := range tokens {
				tokenConfig, err := tokenAdminRegistryContract.DevInspect().GetTokenConfigStruct(ctx, callOpts, ccipRefObj, tokenAddr)
				if err != nil {
					// Token might not be configured properly, skip it
					continue
				}

				tokenType := tokenConfig.TokenType
				if !strings.HasPrefix(tokenType, "0x") {
					tokenType = "0x" + tokenType
				}

				tokenConfigs[tokenAddr] = TokenConfigView{
					TokenPoolPackageId:  tokenConfig.TokenPoolPackageId,
					TokenPoolModule:     tokenConfig.TokenPoolModule,
					TokenType:           tokenConfig.TokenType,
					Administrator:       tokenConfig.Administrator,
					TokenPoolTypeProof:  tokenConfig.TokenPoolTypeProof,
					LockOrBurnParams:    tokenConfig.LockOrBurnParams,
					ReleaseOrMintParams: tokenConfig.ReleaseOrMintParams,
				}
			}
		}
	}

	return TokenAdminRegistryView{
		ContractMetaData: ContractMetaData{
			Address:        ccipPackageID,
			Owner:          "",
			TypeAndVersion: typeAndVersion,
			StateObjectID:  ccipRefObj.Id,
		},
		TokenConfigs: tokenConfigs,
	}, nil
}

func generateNonceManagerView(
	ctx context.Context,
	chain sui.Chain,
	ccipPackageID string,
	callOpts *bind.CallOpts,
) (NonceManagerView, error) {
	// Create nonce manager contract binding
	nonceManagerContract, err := module_nonce_manager.NewNonceManager(ccipPackageID, chain.Client)
	if err != nil {
		return NonceManagerView{}, fmt.Errorf("failed to create nonce manager contract binding: %w", err)
	}

	// Get type and version
	typeAndVersion, err := nonceManagerContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return NonceManagerView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	return NonceManagerView{
		ContractMetaData: ContractMetaData{
			Address:        ccipPackageID,
			Owner:          "",
			TypeAndVersion: typeAndVersion,
		},
	}, nil
}

func generateReceiverRegistryView(
	ctx context.Context,
	chain sui.Chain,
	ccipPackageID string,
	callOpts *bind.CallOpts,
) (ReceiverRegistryView, error) {
	// Create receiver registry contract binding
	receiverRegistryContract, err := module_receiver_registry.NewReceiverRegistry(ccipPackageID, chain.Client)
	if err != nil {
		return ReceiverRegistryView{}, fmt.Errorf("failed to create receiver registry contract binding: %w", err)
	}

	// Get type and version
	typeAndVersion, err := receiverRegistryContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return ReceiverRegistryView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	return ReceiverRegistryView{
		ContractMetaData: ContractMetaData{
			Address:        ccipPackageID,
			Owner:          "",
			TypeAndVersion: typeAndVersion,
		},
	}, nil
}
