package view

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
)

type RouterView struct {
	ContractMetaData

	IsTestRouter bool              `json:"isTestRouter"`
	OnRamps      map[uint64]string `json:"onRamps"`  // Map of DestinationChainSelector to OnRampAddress
	OffRamps     map[uint64]string `json:"offRamps"` // Map of DestinationChainSelector to OffRampAddress
}

// GenerateRouterView generates a router view for a given router by querying the on-chain state
func GenerateRouterView(
	ctx context.Context,
	chain sui.Chain,
	routerPackageID string,
	routerStateObjectID string,
) (RouterView, error) {
	if routerPackageID == "" || routerStateObjectID == "" {
		return RouterView{}, fmt.Errorf("routerPackageID and routerStateObjectID cannot be empty")
	}

	// Create router contract binding
	routerContract, err := module_router.NewRouter(routerPackageID, chain.Client)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to create router contract binding: %w", err)
	}

	routerStateObj := bind.Object{Id: routerStateObjectID}
	callOpts := &bind.CallOpts{Signer: chain.Signer}

	// Get owner
	owner, err := routerContract.DevInspect().Owner(ctx, callOpts, routerStateObj)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get owner: %w", err)
	}

	// Get type and version
	typeAndVersion, err := routerContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	destChainSelectors, err := routerContract.DevInspect().GetDestChains(ctx, callOpts, routerStateObj)
	if err != nil {
		return RouterView{}, fmt.Errorf("failed to get destination chain selectors: %w", err)
	}

	onRamps := make(map[uint64]string)
	for _, selector := range destChainSelectors {
		onRamp, err := routerContract.DevInspect().GetOnRamp(ctx, callOpts, routerStateObj, selector)
		if err != nil {
			return RouterView{}, fmt.Errorf("failed to get onRamp for selector %d: %w", selector, err)
		}
		onRamps[selector] = onRamp
	}

	return RouterView{
		ContractMetaData: ContractMetaData{
			Address:        routerPackageID,
			Owner:          owner,
			TypeAndVersion: typeAndVersion,
			StateObjectID:  routerStateObjectID,
		},
		IsTestRouter: false, // TODO: Determine from contract state or config
		OnRamps:      onRamps,
	}, nil
}
