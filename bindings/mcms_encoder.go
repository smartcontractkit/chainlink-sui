package mcmsencoder

import (
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	module_usdc_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/usdc_token_pool"
	module_managed_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/managed_token/managed_token"
)

type CCIPEntrypointArgEncoder struct {
	registryObjID string
	client        sui.ISuiAPI
}

func NewCCIPEntrypointArgEncoder() *CCIPEntrypointArgEncoder {
	return &CCIPEntrypointArgEncoder{}
}

// package-level constructor variables so tests can inject fakes
var (
	newFeeQuoter            = module_fee_quoter.NewFeeQuoter
	newOfframp              = module_offramp.NewOfframp
	newOnramp               = module_onramp.NewOnramp
	newRouter               = module_router.NewRouter
	newBurnMintTokenPool    = module_burn_mint_token_pool.NewBurnMintTokenPool
	newLockReleaseTokenPool = module_lock_release_token_pool.NewLockReleaseTokenPool
	newManagedTokenPool     = module_managed_token_pool.NewManagedTokenPool
	newUsdcTokenPool        = module_usdc_token_pool.NewUsdcTokenPool
	newManagedToken         = module_managed_token.NewManagedToken
)

// MCMS SDK will call this to encode the entrypoint call
// Data is the raw BCS encoded bytes of the final function call
func (e *CCIPEntrypointArgEncoder) EncodeEntryPointArg(executingCallbackParams *transaction.Argument, target, module, function, stateObjID string, data []byte) (*bind.EncodedCall, error) {
	clock := bind.Object{Id: "0x6"}
	stateObj := bind.Object{Id: stateObjID}
	registryObj := bind.Object{Id: e.registryObjID}

	encodeWithCCIPObjectRefAndState := func() (*bind.EncodedCall, error) {
		// TODO: These should come from decoded data
		ccipRef := bind.Object{Id: "0x123"}
		offramp, err := module_offramp.NewOfframp(target, e.client)
		if err != nil {
			return nil, err
		}
		// The function signature is the same for all ccip entrypoints that require the ccip object ref, so we can use any of them to encode
		entrypointCall, err := offramp.Encoder().McmsSetDynamicConfigWithArgs(ccipRef, stateObj, registryObj, executingCallbackParams)
		if err != nil {
			return nil, err
		}
		// Override the module info with the actual target
		entrypointCall.Module.ModuleName = module
		// mcms entrypoint like functions are the target function prefixed with `mcms_`
		entrypointCall.Function = fmt.Sprintf("mcms_%s", strings.TrimPrefix(function, "mcms_"))
		return entrypointCall, nil
	}

	encodeDefaultWithTypeArgsAndClock := func() (*bind.EncodedCall, error) {
		burnMintTokenPool, err := module_burn_mint_token_pool.NewBurnMintTokenPool(target, e.client)
		if err != nil {
			return nil, err
		}
		// TODO: Find correct type args
		typeArgs := []string{"0x1::sui::SUI"}
		entrypointCall, err := burnMintTokenPool.Encoder().McmsSetChainRateLimiterConfigsWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams, clock)
		if err != nil {
			return nil, err
		}
		// Override the module info with the actual target
		entrypointCall.Module.ModuleName = module
		// mcms entrypoint like functions are the target function prefixed with `mcms_`
		entrypointCall.Function = fmt.Sprintf("mcms_%s", strings.TrimPrefix(function, "mcms_"))
		return entrypointCall, nil
	}

	encodeDefaultWithTypeArgs := func() (*bind.EncodedCall, error) {
		burnMintTokenPool, err := module_burn_mint_token_pool.NewBurnMintTokenPool(target, e.client)
		if err != nil {
			return nil, err
		}
		// TODO: Find correct type args
		typeArgs := []string{"0x1::sui::SUI"}
		entrypointCall, err := burnMintTokenPool.Encoder().McmsExecuteOwnershipTransferWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		if err != nil {
			return nil, err
		}
		// Override the module info with the actual target
		entrypointCall.Module.ModuleName = module
		// mcms entrypoint like functions are the target function prefixed with `mcms_`
		entrypointCall.Function = fmt.Sprintf("mcms_%s", strings.TrimPrefix(function, "mcms_"))
		return entrypointCall, nil
	}

	switch module {
	// FEE QUOTER
	case "fee_quoter":
		feeQuoter, err := module_fee_quoter.NewFeeQuoter(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "update_prices_with_owner_cap":
			return feeQuoter.Encoder().McmsUpdatePricesWithOwnerCapWithArgs(stateObj, registryObj, clock, executingCallbackParams)
		}

	// OFFRAMP
	case "offramp":
		switch function {
		case "accept_ownership":
		case "set_dynamic_config":
		case "apply_source_chain_config_updates":
		case "set_ocr3_config":
		case "transfer_ownership":
		case "execute_ownership_transfer":
			return encodeWithCCIPObjectRefAndState()
		}

	// ONRAMP
	case "onramp":
		onramp, err := module_onramp.NewOnramp(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
		case "set_dynamic_config":
		case "apply_dest_chain_config_updates":
		case "apply_allowlist_updates":
		case "transfer_ownership":
		case "execute_ownership_transfer":
			return encodeWithCCIPObjectRefAndState()
		case "initialize":
			// TODO: These should come from decoded data
			nonceManagerCap := bind.Object{Id: "0x456"}
			sourceTransferCap := bind.Object{Id: "0x789"}
			return onramp.Encoder().McmsInitializeWithArgs(stateObj, registryObj, nonceManagerCap, sourceTransferCap, executingCallbackParams)
		case "withdraw_fee_tokens":
			// TODO: These should come from decoded data
			coinMetadata := bind.Object{Id: "0xabc"}
			ccipRef := bind.Object{Id: "0x123"}
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			return onramp.Encoder().McmsWithdrawFeeTokensWithArgs(typeArgs, ccipRef, stateObj, registryObj, coinMetadata, executingCallbackParams)
		}

	// ROUTER
	case "router":
		router, err := module_router.NewRouter(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			return router.Encoder().McmsAcceptOwnershipWithArgs(stateObj, executingCallbackParams)
		}

	// BURN MINT TOKEN POOL
	case "burn_mint_token_pool":
		burnMintTokenPool, err := module_burn_mint_token_pool.NewBurnMintTokenPool(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			return burnMintTokenPool.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, executingCallbackParams)
		case "set_allowlist_enabled":
		case "apply_allowlist_updates":
		case "apply_chain_updates":
		case "add_remote_pool":
		case "remove_remote_pool":
		case "transfer_ownership":
		case "execute_ownership_transfer":
			return encodeDefaultWithTypeArgs()
		case "set_chain_rate_limiter_configs":
		case "set_chain_rate_limiter_config":
			return encodeDefaultWithTypeArgs()
		}

	// LOCK RELEASE TOKEN POOL
	case "lock_release_token_pool":
		lockReleaseTokenPool, err := module_lock_release_token_pool.NewLockReleaseTokenPool(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			return lockReleaseTokenPool.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, executingCallbackParams)
		case "set_rebalancer":
		case "set_allowlist_enabled":
		case "apply_allowlist_updates":
		case "apply_chain_updates":
		case "add_remote_pool":
		case "remove_remote_pool":
		case "transfer_ownership":
		case "execute_ownership_transfer":
			return encodeDefaultWithTypeArgs()
		case "set_chain_rate_limiter_configs":
		case "set_chain_rate_limiter_config":
			return encodeDefaultWithTypeArgsAndClock()

		}

	// MANAGED TOKEN POOL
	case "managed_token_pool":
		managedTokenPool, err := module_managed_token_pool.NewManagedTokenPool(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			return managedTokenPool.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, executingCallbackParams)
		case "set_allowlist_enabled":
		case "apply_allowlist_updates":
		case "apply_chain_updates":
		case "add_remote_pool":
		case "remove_remote_pool":
		case "transfer_ownership":
		case "execute_ownership_transfer":
			return encodeDefaultWithTypeArgs()
		case "set_chain_rate_limiter_configs":
		case "set_chain_rate_limiter_config":
			return encodeDefaultWithTypeArgsAndClock()
		}

	// USDC TOKEN POOL
	case "usdc_token_pool":
		usdcTokenPool, err := module_usdc_token_pool.NewUsdcTokenPool(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			return usdcTokenPool.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, executingCallbackParams)
		case "set_allowlist_enabled":
		case "apply_allowlist_updates":
		case "apply_chain_updates":
		case "add_remote_pool":
		case "remove_remote_pool":
		case "transfer_ownership":
		case "execute_ownership_transfer":
			return encodeDefaultWithTypeArgs()
		case "set_chain_rate_limiter_configs":
		case "set_chain_rate_limiter_config":
			return encodeDefaultWithTypeArgsAndClock()
		}

	// MANAGED TOKEN
	case "managed_token":
		managedToken, err := module_managed_token.NewManagedToken(target, e.client)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			return managedToken.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, executingCallbackParams)
		case "configure_new_minter":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			// TODO: This doesn't exist
			denyList := bind.Object{Id: "0xdef"}
			return managedToken.Encoder().McmsConfigureNewMinterWithArgs(typeArgs, stateObj, registryObj, denyList, executingCallbackParams)
		case "increment_mint_allowance":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			denyList := bind.Object{Id: "0xdef"}
			return managedToken.Encoder().McmsIncrementMintAllowanceWithArgs(typeArgs, stateObj, registryObj, denyList, executingCallbackParams)
		case "set_unlimited_mint_allowances":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			denyList := bind.Object{Id: "0xdef"}
			return managedToken.Encoder().McmsSetUnlimitedMintAllowancesWithArgs(typeArgs, stateObj, registryObj, denyList, executingCallbackParams)
		case "blocklist":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			denyList := bind.Object{Id: "0xdef"}
			return managedToken.Encoder().McmsBlocklistWithArgs(typeArgs, stateObj, registryObj, denyList, executingCallbackParams)
		case "unblocklist":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			denyList := bind.Object{Id: "0xdef"}
			return managedToken.Encoder().McmsUnblocklistWithArgs(typeArgs, stateObj, registryObj, denyList, executingCallbackParams)
		case "pause":
			// TODO: Find correct type args
			typeArgs := []string{"0x1::sui::SUI"}
			denyList := bind.Object{Id: "0xdef"}
			return managedToken.Encoder().McmsPauseWithArgs(typeArgs, stateObj, registryObj, denyList, executingCallbackParams)
		}

		// END OF MODULE SWITCH
	}

	// FALLBACK CASE: Use Fee Quoter as it has the most common function signatures
	// Fallback to fee quoter for any unhandled module/function
	// This works because most mcms functions have the same signature
	// state: &State, registry: &Registry, executing_callback_params: &ExecutingCallbackParams
	// If a function has a different signature, it should be handled explicitly above
	feeQuoter, err := module_fee_quoter.NewFeeQuoter(target, e.client)
	if err != nil {
		return nil, err
	}

	entryPointCall, err := feeQuoter.Encoder().McmsApplyFeeTokenUpdatesWithArgs(
		stateObj,
		registryObj,
		executingCallbackParams,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mcms_entrypoint call: %w", err)
	}
	// Override the module info with the actual target
	entryPointCall.Module.ModuleName = module
	// mcms entrypoint like functions are the target function prefixed with `mcms_`
	entryPointCall.Function = fmt.Sprintf("mcms_%s", strings.TrimPrefix(function, "mcms_"))

	return entryPointCall, nil
}
