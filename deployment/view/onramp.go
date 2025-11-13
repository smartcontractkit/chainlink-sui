package view

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
)

type OnRampView struct {
	ContractMetaData

	StaticConfig          OnRampStaticConfig               `json:"staticConfig"`
	DynamicConfig         OnRampDynamicConfig              `json:"dynamicConfig"`
	DestChainSpecificData map[uint64]DestChainSpecificData `json:"destChainSpecificData"`
}

type OnRampStaticConfig struct {
	ChainSelector uint64 `json:"chainSelector"`
}

type OnRampDynamicConfig struct {
	FeeAggregator  string `json:"feeAggregator"`
	AllowlistAdmin string `json:"allowlistAdmin"`
}

type DestChainSpecificData struct {
	AllowedSendersList []string              `json:"allowedSendersList"`
	DestChainConfig    OnRampDestChainConfig `json:"destChainConfig"`
	ExpectedNextSeqNum uint64                `json:"expectedNextSeqNum"`
}

type OnRampDestChainConfig struct {
	SequenceNumber   uint64 `json:"sequenceNumber"`
	AllowlistEnabled bool   `json:"allowlistEnabled"`
	Router           string `json:"router"`
}

// GenerateOnRampView generates an onramp view for a given onramp by querying the on-chain state
func GenerateOnRampView(
	ctx context.Context,
	chain sui.Chain,
	onRampPackageID string,
	onRampStateObjectID string,
	routerPackageID string,
	routerStateObjectID string,
) (OnRampView, error) {
	if onRampPackageID == "" || onRampStateObjectID == "" {
		return OnRampView{}, fmt.Errorf("onRampPackageID and onRampStateObjectID cannot be empty")
	}

	// Create onramp contract binding
	onRampContract, err := module_onramp.NewOnramp(onRampPackageID, chain.Client)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to create onramp contract binding: %w", err)
	}

	onRampStateObj := bind.Object{Id: onRampStateObjectID}
	callOpts := &bind.CallOpts{Signer: chain.Signer}

	// Get owner
	owner, err := onRampContract.DevInspect().Owner(ctx, callOpts, onRampStateObj)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get owner: %w", err)
	}

	// Get type and version
	typeAndVersion, err := onRampContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	// Get static config
	staticConfig, err := onRampContract.DevInspect().GetStaticConfig(ctx, callOpts, onRampStateObj)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get static config: %w", err)
	}

	// Get dynamic config
	dynamicConfig, err := onRampContract.DevInspect().GetDynamicConfig(ctx, callOpts, onRampStateObj)
	if err != nil {
		return OnRampView{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	// TODO: Changesets are not configuring router, any configuration that requires GetDestChains will be empty
	// Query the router to get destination chains
	var destChainSelectors []uint64
	if routerPackageID != "" && routerStateObjectID != "" {
		routerContract, err := module_router.NewRouter(routerPackageID, chain.Client)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to create router contract binding: %w", err)
		}

		routerStateObj := bind.Object{Id: routerStateObjectID}
		destChainSelectors, err = routerContract.DevInspect().GetDestChains(ctx, callOpts, routerStateObj)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get dest chains from router: %w", err)
		}
	}

	// Get dest chain specific data for each known destination chain
	destChainSpecificData := make(map[uint64]DestChainSpecificData)
	for _, destChainSelector := range destChainSelectors {
		// Get dest chain config
		// GetDestChainConfig returns [0]: u64 (sequence_number), [1]: bool (allowlist_enabled), [2]: address (router)
		destChainConfigRaw, err := onRampContract.DevInspect().GetDestChainConfig(ctx, callOpts, onRampStateObj, destChainSelector)
		if err != nil {
			// Chain might not be configured, skip it
			continue
		}

		if len(destChainConfigRaw) < 3 {
			continue
		}

		sequenceNumber, ok := destChainConfigRaw[0].(uint64)
		if !ok {
			return OnRampView{}, fmt.Errorf("unexpected type for sequence number: got %T", destChainConfigRaw[0])
		}

		allowlistEnabled, ok := destChainConfigRaw[1].(bool)
		if !ok {
			return OnRampView{}, fmt.Errorf("unexpected type for allowlist enabled: got %T", destChainConfigRaw[1])
		}

		router, ok := destChainConfigRaw[2].(string)
		if !ok {
			return OnRampView{}, fmt.Errorf("unexpected type for router: got %T", destChainConfigRaw[2])
		}

		// Get allowed senders list
		// GetAllowedSendersList returns [0]: bool (allowlist_enabled), [1]: vector<address> (allowed_senders)
		allowedSendersRaw, err := onRampContract.DevInspect().GetAllowedSendersList(ctx, callOpts, onRampStateObj, destChainSelector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get allowed senders list for chain %d: %w", destChainSelector, err)
		}

		allowedSendersList := []string{}
		if len(allowedSendersRaw) >= 2 {
			if senders, ok := allowedSendersRaw[1].([]string); ok {
				allowedSendersList = senders
			}
		}

		// Get expected next sequence number
		expectedNextSeqNum, err := onRampContract.DevInspect().GetExpectedNextSequenceNumber(ctx, callOpts, onRampStateObj, destChainSelector)
		if err != nil {
			return OnRampView{}, fmt.Errorf("failed to get expected next sequence number for chain %d: %w", destChainSelector, err)
		}

		destChainSpecificData[destChainSelector] = DestChainSpecificData{
			AllowedSendersList: allowedSendersList,
			DestChainConfig: OnRampDestChainConfig{
				SequenceNumber:   sequenceNumber,
				AllowlistEnabled: allowlistEnabled,
				Router:           router,
			},
			ExpectedNextSeqNum: expectedNextSeqNum,
		}
	}

	return OnRampView{
		ContractMetaData: ContractMetaData{
			Address:        onRampPackageID,
			Owner:          owner,
			TypeAndVersion: typeAndVersion,
			StateObjectID:  onRampStateObjectID,
		},
		StaticConfig: OnRampStaticConfig{
			ChainSelector: staticConfig.ChainSelector,
		},
		DynamicConfig: OnRampDynamicConfig{
			FeeAggregator:  dynamicConfig.FeeAggregator,
			AllowlistAdmin: dynamicConfig.AllowlistAdmin,
		},
		DestChainSpecificData: destChainSpecificData,
	}, nil
}
