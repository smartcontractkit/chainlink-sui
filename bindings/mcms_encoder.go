package mcmsencoder

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	module_upgrade_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/upgrade_registry"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	module_managed_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/managed_token/managed_token"
)

var SuiAddressLength = 32

type CCIPEntrypointArgEncoder struct {
	registryObjID      string
	deployerStateObjID string
}

func NewCCIPEntrypointArgEncoder(registryObjID string, deployerStateObjID string) *CCIPEntrypointArgEncoder {
	return &CCIPEntrypointArgEncoder{registryObjID: registryObjID, deployerStateObjID: deployerStateObjID}
}

func deserializeFirst32Bytes(data []byte) []byte {
	deserializer := bcs.NewDeserializer(data)
	return deserializer.ReadFixedBytes(SuiAddressLength)
}

func toHexString(data []byte) string {
	return fmt.Sprintf("0x%s", strings.ToLower(hex.EncodeToString(data)))
}

func overrideCall(call *bind.EncodedCall, module, function string) *bind.EncodedCall {
	call.Module.ModuleName = module
	call.Function = fmt.Sprintf("mcms_%s", strings.TrimPrefix(function, "mcms_"))

	return call
}

// MCMS SDK will call this to encode the entrypoint call
// Data is the raw BCS encoded bytes of the final function call
func (e *CCIPEntrypointArgEncoder) EncodeEntryPointArg(executingCallbackParams *transaction.Argument, target, module, function, stateObjID string, data []byte, typeArgs []string) (*bind.EncodedCall, error) {
	clock := bind.Object{Id: "0x6"}
	stateObj := bind.Object{Id: stateObjID}
	registryObj := bind.Object{Id: e.registryObjID}
	deployerStateObj := bind.Object{Id: e.deployerStateObjID}

	encodeWithCCIPObjectRefAndState := func() (*bind.EncodedCall, error) {
		// Deserialize the ccip object ref (always the first 32 bytes)
		ccipRefBytes := deserializeFirst32Bytes(data)
		ccipRef := bind.Object{Id: toHexString(ccipRefBytes)}

		offramp, err := module_offramp.NewOfframp(target, nil)
		if err != nil {
			return nil, err
		}
		// The function signature is the same for all ccip entrypoints that require the ccip object ref, so we can use any of them to encode
		entrypointCall, err := offramp.Encoder().McmsSetDynamicConfigWithArgs(ccipRef, stateObj, registryObj, executingCallbackParams)
		if err != nil {
			return nil, err
		}

		return overrideCall(entrypointCall, module, function), nil
	}

	encodeDefaultWithTypeArgsAndClock := func() (*bind.EncodedCall, error) {
		burnMintTokenPool, err := module_burn_mint_token_pool.NewBurnMintTokenPool(target, nil)
		if err != nil {
			return nil, err
		}

		entrypointCall, err := burnMintTokenPool.Encoder().McmsSetChainRateLimiterConfigsWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams, clock)
		if err != nil {
			return nil, err
		}

		return overrideCall(entrypointCall, module, function), nil
	}

	encodeDefaultWithTypeArgs := func() (*bind.EncodedCall, error) {
		burnMintTokenPool, err := module_burn_mint_token_pool.NewBurnMintTokenPool(target, nil)
		if err != nil {
			return nil, err
		}

		entrypointCall, err := burnMintTokenPool.Encoder().McmsSetAllowlistEnabledWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		if err != nil {
			return nil, err
		}

		return overrideCall(entrypointCall, module, function), nil
	}

	encodeExecuteOwnershipTransferWithTypeArgs := func() (*bind.EncodedCall, error) {
		burnMintTokenPool, err := module_burn_mint_token_pool.NewBurnMintTokenPool(target, nil)
		if err != nil {
			return nil, err
		}

		entrypointCall, err := burnMintTokenPool.Encoder().McmsExecuteOwnershipTransferWithArgs(typeArgs, stateObj, registryObj, deployerStateObj, executingCallbackParams)
		if err != nil {
			return nil, err
		}

		return overrideCall(entrypointCall, module, function), nil
	}

	switch module {
	// FEE QUOTER
	case "fee_quoter":
		feeQuoter, err := module_fee_quoter.NewFeeQuoter(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "update_prices_with_owner_cap":
			ccipRefBytes := deserializeFirst32Bytes(data)
			ccipRef := bind.Object{Id: toHexString(ccipRefBytes)}
			if ccipRef.Id != stateObj.Id {
				return nil, fmt.Errorf("ccip ref (%s) does not match state object (%s)", ccipRef.Id, stateObj.Id)
			}

			return feeQuoter.Encoder().McmsUpdatePricesWithOwnerCapWithArgs(ccipRef, registryObj, clock, executingCallbackParams)
		case "apply_fee_token_updates":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return feeQuoter.Encoder().McmsApplyFeeTokenUpdatesWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "apply_dest_chain_config_updates":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return feeQuoter.Encoder().McmsApplyDestChainConfigUpdatesWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "apply_token_transfer_fee_config_updates":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return feeQuoter.Encoder().McmsApplyTokenTransferFeeConfigUpdatesWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "apply_premium_multiplier_wei_per_eth_updates":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return feeQuoter.Encoder().McmsApplyPremiumMultiplierWeiPerEthUpdatesWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "new_fee_quoter_cap_and_transfer":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return feeQuoter.Encoder().McmsNewFeeQuoterCapAndTransferWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "destroy_fee_quoter_cap":
			deserializer := bcs.NewDeserializer(data)
			deserializer.ReadFixedBytes(SuiAddressLength)
			deserializer.ReadFixedBytes(SuiAddressLength)
			capBytes := deserializer.ReadFixedBytes(SuiAddressLength)
			ccipRef := bind.Object{Id: stateObjID}
			c := bind.Object{Id: toHexString(capBytes)}
			return feeQuoter.Encoder().McmsDestroyFeeQuoterCapWithArgs(ccipRef, registryObj, executingCallbackParams, c)
		default:
			return nil, fmt.Errorf("unsupported fee_quoter MCMS function: %q", function)
		}

	// STATE OBJECT
	case "state_object":
		moduleStateObj, err := module_state_object.NewStateObject(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			return moduleStateObj.Encoder().McmsAcceptOwnershipWithArgs(stateObj, registryObj, executingCallbackParams)
		case "transfer_ownership":
			// For state_object, the state object is the CCIP object ref
			return moduleStateObj.Encoder().McmsTransferOwnershipWithArgs(stateObj, registryObj, executingCallbackParams)
		case "execute_ownership_transfer":
			return moduleStateObj.Encoder().McmsExecuteOwnershipTransferWithArgs(stateObj, registryObj, deployerStateObj, executingCallbackParams)
		case "add_package_id":
			return moduleStateObj.Encoder().McmsAddPackageIdWithArgs(stateObj, registryObj, executingCallbackParams)
		case "remove_package_id":
			return moduleStateObj.Encoder().McmsRemovePackageIdWithArgs(stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return moduleStateObj.Encoder().McmsAddAllowedModulesWithArgs(registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return moduleStateObj.Encoder().McmsRemoveAllowedModulesWithArgs(registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported state_object MCMS function: %q", function)
		}

	// OFFRAMP
	case "offramp":
		offramp, err := module_offramp.NewOfframp(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "set_dynamic_config",
			"apply_source_chain_config_updates",
			"set_ocr3_config",
			"transfer_ownership":
			return encodeWithCCIPObjectRefAndState()
		case "execute_ownership_transfer":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return offramp.Encoder().McmsExecuteOwnershipTransferWithArgs(ccipRef, stateObj, registryObj, deployerStateObj, executingCallbackParams)
		case "accept_ownership":
			ccipObjectRef := bind.Object{Id: stateObjID} // For accept_ownership, the state object is the CCIP object ref
			stateObj := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}

			return offramp.Encoder().McmsAcceptOwnershipWithArgs(ccipObjectRef, stateObj, registryObj, executingCallbackParams)
		case "add_package_id":
			return offramp.Encoder().McmsAddPackageIdWithArgs(stateObj, registryObj, executingCallbackParams)
		case "remove_package_id":
			return offramp.Encoder().McmsRemovePackageIdWithArgs(stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return offramp.Encoder().McmsAddAllowedModulesWithArgs(registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return offramp.Encoder().McmsRemoveAllowedModulesWithArgs(registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported offramp MCMS function: %q", function)
		}

	// ONRAMP
	case "onramp":
		onramp, err := module_onramp.NewOnramp(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			ccipObjectRef := bind.Object{Id: stateObjID} // For accept_ownership, the state object is the CCIP object ref
			stateObj := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}

			return onramp.Encoder().McmsAcceptOwnershipWithArgs(ccipObjectRef, stateObj, registryObj, executingCallbackParams)

		case "set_dynamic_config",
			"apply_dest_chain_config_updates",
			"apply_allowlist_updates",
			"transfer_ownership":
			return encodeWithCCIPObjectRefAndState()
		case "execute_ownership_transfer":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return onramp.Encoder().McmsExecuteOwnershipTransferWithArgs(ccipRef, stateObj, registryObj, deployerStateObj, executingCallbackParams)
		case "add_package_id":
			return onramp.Encoder().McmsAddPackageIdWithArgs(stateObj, registryObj, executingCallbackParams)
		case "remove_package_id":
			return onramp.Encoder().McmsRemovePackageIdWithArgs(stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return onramp.Encoder().McmsAddAllowedModulesWithArgs(registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return onramp.Encoder().McmsRemoveAllowedModulesWithArgs(registryObj, executingCallbackParams)
		case "withdraw_fee_tokens":
			deserializer := bcs.NewDeserializer(data)
			ccipRefBytes := deserializer.ReadFixedBytes(SuiAddressLength)
			state := deserializer.ReadFixedBytes(SuiAddressLength)
			deserializer.ReadFixedBytes(SuiAddressLength) // skip owner cap, we don't need it
			feeTokenMetadata := deserializer.ReadFixedBytes(SuiAddressLength)

			coinMetadata := bind.Object{Id: toHexString(feeTokenMetadata)}
			ccipRef := bind.Object{Id: toHexString(ccipRefBytes)}

			if toHexString(state) != stateObj.Id {
				return nil, fmt.Errorf("state (%s) does not match state object (%s)", toHexString(state), stateObj.Id)
			}

			return onramp.Encoder().McmsWithdrawFeeTokensWithArgs(typeArgs, ccipRef, stateObj, registryObj, coinMetadata, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported onramp MCMS function: %q", function)
		}

	// ROUTER
	case "router":
		router, err := module_router.NewRouter(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			return router.Encoder().McmsAcceptOwnershipWithArgs(stateObj, registryObj, executingCallbackParams)
		case "transfer_ownership":
			return router.Encoder().McmsTransferOwnershipWithArgs(stateObj, registryObj, executingCallbackParams)
		case "execute_ownership_transfer":
			return router.Encoder().McmsExecuteOwnershipTransferWithArgs(stateObj, registryObj, deployerStateObj, executingCallbackParams)
		case "set_on_ramps":
			return router.Encoder().McmsSetOnRampsWithArgs(stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return router.Encoder().McmsAddAllowedModulesWithArgs(registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return router.Encoder().McmsRemoveAllowedModulesWithArgs(registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported router MCMS function: %q", function)
		}

	// BURN MINT TOKEN POOL
	case "burn_mint_token_pool":
		burnMintTokenPool, err := module_burn_mint_token_pool.NewBurnMintTokenPool(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			return burnMintTokenPool.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "set_allowlist_enabled",
			"apply_allowlist_updates",
			"apply_chain_updates",
			"add_remote_pool",
			"remove_remote_pool",
			"transfer_ownership":
			return encodeDefaultWithTypeArgs()
		case "execute_ownership_transfer":
			return encodeExecuteOwnershipTransferWithTypeArgs()
		case "set_chain_rate_limiter_configs",
			"set_chain_rate_limiter_config":
			return encodeDefaultWithTypeArgsAndClock()
		case "destroy_token_pool":
			deserializer := bcs.NewDeserializer(data)
			ccipRef := bind.Object{Id: toHexString(deserializer.ReadFixedBytes(SuiAddressLength))}
			state := deserializer.ReadFixedBytes(SuiAddressLength)
			if toHexString(state) != stateObj.Id {
				return nil, fmt.Errorf("state (%s) does not match state object (%s)", toHexString(state), stateObj.Id)
			}

			return burnMintTokenPool.Encoder().McmsDestroyTokenPoolWithArgs(typeArgs, ccipRef, stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return burnMintTokenPool.Encoder().McmsAddAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return burnMintTokenPool.Encoder().McmsRemoveAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported burn_mint_token_pool MCMS function: %q", function)
		}

	// LOCK RELEASE TOKEN POOL
	case "lock_release_token_pool":
		lockReleaseTokenPool, err := module_lock_release_token_pool.NewLockReleaseTokenPool(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":
			return lockReleaseTokenPool.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "set_rebalancer",
			"set_allowlist_enabled",
			"apply_allowlist_updates",
			"apply_chain_updates",
			"add_remote_pool",
			"remove_remote_pool",
			"transfer_ownership":
			return encodeDefaultWithTypeArgs()
		case "execute_ownership_transfer":
			return encodeExecuteOwnershipTransferWithTypeArgs()
		case "set_chain_rate_limiter_configs",
			"set_chain_rate_limiter_config":
			return encodeDefaultWithTypeArgsAndClock()
		case "destroy_rebalancer_cap":
			return lockReleaseTokenPool.Encoder().McmsDestroyRebalancerCapWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "withdraw_liquidity":
			return lockReleaseTokenPool.Encoder().McmsWithdrawLiquidityWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "provide_liquidity":
			deserializer := bcs.NewDeserializer(data)
			state := deserializer.ReadFixedBytes(SuiAddressLength)
			deserializer.ReadFixedBytes(SuiAddressLength) // skip rebalancer cap, we don't need it
			coin := bind.Object{Id: toHexString(deserializer.ReadFixedBytes(SuiAddressLength))}
			if toHexString(state) != stateObj.Id {
				return nil, fmt.Errorf("state (%s) does not match state object (%s)", toHexString(state), stateObj.Id)
			}

			return lockReleaseTokenPool.Encoder().McmsProvideLiquidityWithArgs(typeArgs, stateObj, registryObj, coin, executingCallbackParams)
		case "destroy_token_pool":
			deserializer := bcs.NewDeserializer(data)
			ccipRef := bind.Object{Id: toHexString(deserializer.ReadFixedBytes(SuiAddressLength))}
			state := deserializer.ReadFixedBytes(SuiAddressLength)
			if toHexString(state) != stateObj.Id {
				return nil, fmt.Errorf("state (%s) does not match state object (%s)", toHexString(state), stateObj.Id)
			}

			return lockReleaseTokenPool.Encoder().McmsDestroyTokenPoolWithArgs(typeArgs, ccipRef, stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return lockReleaseTokenPool.Encoder().McmsAddAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return lockReleaseTokenPool.Encoder().McmsRemoveAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported lock_release_token_pool MCMS function: %q", function)
		}

	// MANAGED TOKEN POOL
	case "managed_token_pool":
		managedTokenPool, err := module_managed_token_pool.NewManagedTokenPool(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":

			return managedTokenPool.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "set_allowlist_enabled",
			"apply_allowlist_updates",
			"apply_chain_updates",
			"add_remote_pool",
			"remove_remote_pool",
			"transfer_ownership":
			return encodeDefaultWithTypeArgs()
		case "execute_ownership_transfer":
			return encodeExecuteOwnershipTransferWithTypeArgs()
		case "set_chain_rate_limiter_configs",
			"set_chain_rate_limiter_config":
			return encodeDefaultWithTypeArgsAndClock()
		case "destroy_token_pool":
			deserializer := bcs.NewDeserializer(data)
			ccipRef := bind.Object{Id: toHexString(deserializer.ReadFixedBytes(SuiAddressLength))}
			state := deserializer.ReadFixedBytes(SuiAddressLength)
			if toHexString(state) != stateObj.Id {
				return nil, fmt.Errorf("state (%s) does not match state object (%s)", toHexString(state), stateObj.Id)
			}

			return managedTokenPool.Encoder().McmsDestroyTokenPoolWithArgs(typeArgs, ccipRef, stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return managedTokenPool.Encoder().McmsAddAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return managedTokenPool.Encoder().McmsRemoveAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported managed_token_pool MCMS function: %q", function)
		}

	// RMN REMOTE
	case "rmn_remote":
		rmnRemote, err := module_rmn_remote.NewRmnRemote(target, nil)
		if err != nil {
			return nil, err
		}
		ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
		switch function {
		case "set_config":
			return rmnRemote.Encoder().McmsSetConfigWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "curse", "curse_multiple":
			entrypointCall, err := rmnRemote.Encoder().McmsCurseMultipleWithArgs(ccipRef, registryObj, executingCallbackParams)
			if err != nil {
				return nil, err
			}
			return overrideCall(entrypointCall, module, function), nil
		case "uncurse", "uncurse_multiple":
			entrypointCall, err := rmnRemote.Encoder().McmsUncurseMultipleWithArgs(ccipRef, registryObj, executingCallbackParams)
			if err != nil {
				return nil, err
			}
			return overrideCall(entrypointCall, module, function), nil
		case "curse_with_curser_cap", "curse_multiple_with_curser_cap":
			entrypointCall, err := rmnRemote.Encoder().McmsCurseMultipleWithCurserCapWithArgs(ccipRef, registryObj, executingCallbackParams)
			if err != nil {
				return nil, err
			}
			return overrideCall(entrypointCall, module, function), nil
		case "create_curser_cap":
			entrypointCall, err := rmnRemote.Encoder().McmsCreateCurserCapWithArgs(ccipRef, registryObj, executingCallbackParams)
			if err != nil {
				return nil, err
			}
			return overrideCall(entrypointCall, module, function), nil
		case "create_curser_cap_and_transfer":
			entrypointCall, err := rmnRemote.Encoder().McmsCreateCurserCapAndTransferWithArgs(ccipRef, registryObj, executingCallbackParams)
			if err != nil {
				return nil, err
			}
			return overrideCall(entrypointCall, module, function), nil
		case "register_curser_cap":
			deserializer := bcs.NewDeserializer(data)
			deserializer.ReadFixedBytes(SuiAddressLength)
			deserializer.ReadFixedBytes(SuiAddressLength)
			fastRegBytes := deserializer.ReadFixedBytes(SuiAddressLength)
			curserCapBytes := deserializer.ReadFixedBytes(SuiAddressLength)
			fastReg := bind.Object{Id: toHexString(fastRegBytes)}
			curserCap := bind.Object{Id: toHexString(curserCapBytes)}
			entrypointCall, err := rmnRemote.Encoder().McmsRegisterCurserCapWithArgs(ccipRef, registryObj, fastReg, executingCallbackParams, curserCap)
			if err != nil {
				return nil, err
			}
			return overrideCall(entrypointCall, module, function), nil
		case "mint_and_register_curser_cap":
			deserializer := bcs.NewDeserializer(data)
			deserializer.ReadFixedBytes(SuiAddressLength)
			deserializer.ReadFixedBytes(SuiAddressLength)
			fastRegBytes := deserializer.ReadFixedBytes(SuiAddressLength)
			fastReg := bind.Object{Id: toHexString(fastRegBytes)}
			entrypointCall, err := rmnRemote.Encoder().McmsMintAndRegisterCurserCapWithArgs(ccipRef, registryObj, fastReg, executingCallbackParams)
			if err != nil {
				return nil, err
			}
			return overrideCall(entrypointCall, module, function), nil
		case "initialize_allowed_curser_caps":
			return rmnRemote.Encoder().McmsInitializeAllowedCurserCapsWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "register_curser_cap_ids":
			return rmnRemote.Encoder().McmsRegisterCurserCapIdsWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "deregister_curser_cap_ids":
			return rmnRemote.Encoder().McmsDeregisterCurserCapIdsWithArgs(ccipRef, registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported rmn_remote MCMS function: %q", function)
		}

	// TOKEN ADMIN REGISTRY
	case "token_admin_registry":
		tokenAdminRegistry, err := module_token_admin_registry.NewTokenAdminRegistry(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "initialize_local_decimals":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return tokenAdminRegistry.Encoder().McmsInitializeLocalDecimalsWithArgs(
				ccipRef,
				registryObj,
				executingCallbackParams,
			)
		case "backfill_local_decimals", "register_pool", "unregister_pool":
			deserializer := bcs.NewDeserializer(data)
			deserializer.ReadFixedBytes(SuiAddressLength)
			ccipRef := bind.Object{Id: toHexString(deserializer.ReadFixedBytes(SuiAddressLength))}
			switch function {
			case "backfill_local_decimals":
				return tokenAdminRegistry.Encoder().McmsBackfillLocalDecimalsWithArgs(
					ccipRef,
					registryObj,
					executingCallbackParams,
				)
			case "register_pool":
				return tokenAdminRegistry.Encoder().McmsRegisterPoolWithArgs(
					ccipRef,
					registryObj,
					executingCallbackParams,
				)
			case "unregister_pool":
				return tokenAdminRegistry.Encoder().McmsUnregisterPoolWithArgs(
					ccipRef,
					registryObj,
					executingCallbackParams,
				)
			}
		case "transfer_admin_role":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return tokenAdminRegistry.Encoder().McmsTransferAdminRoleWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "accept_admin_role":
			ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
			return tokenAdminRegistry.Encoder().McmsAcceptAdminRoleWithArgs(ccipRef, registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported token_admin_registry MCMS function: %q", function)
		}

	// UPGRADE REGISTRY
	case "upgrade_registry":
		upgradeRegistry, err := module_upgrade_registry.NewUpgradeRegistry(target, nil)
		if err != nil {
			return nil, err
		}
		ccipRef := bind.Object{Id: toHexString(deserializeFirst32Bytes(data))}
		switch function {
		case "block_version":
			return upgradeRegistry.Encoder().McmsBlockVersionWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "unblock_version":
			return upgradeRegistry.Encoder().McmsUnblockVersionWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "block_function":
			return upgradeRegistry.Encoder().McmsBlockFunctionWithArgs(ccipRef, registryObj, executingCallbackParams)
		case "unblock_function":
			return upgradeRegistry.Encoder().McmsUnblockFunctionWithArgs(ccipRef, registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported upgrade_registry MCMS function: %q", function)
		}

	// MANAGED TOKEN
	case "managed_token":
		managedToken, err := module_managed_token.NewManagedToken(target, nil)
		if err != nil {
			return nil, err
		}
		switch function {
		case "accept_ownership":

			return managedToken.Encoder().McmsAcceptOwnershipWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "configure_new_minter":

			return managedToken.Encoder().McmsConfigureNewMinterWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "increment_mint_allowance",
			"set_unlimited_mint_allowances",
			"blocklist",
			"unblocklist",
			"pause",
			"unpause":
			deserializer := bcs.NewDeserializer(data)
			state := deserializer.ReadFixedBytes(SuiAddressLength)
			deserializer.ReadFixedBytes(SuiAddressLength) // skip owner cap, we don't need it
			denyList := deserializer.ReadFixedBytes(SuiAddressLength)

			denyListObj := bind.Object{Id: toHexString(denyList)}

			if toHexString(state) != stateObj.Id {
				return nil, fmt.Errorf("state (%s) does not match state object (%s)", toHexString(state), stateObj.Id)
			}
			entrypointCall, err := managedToken.Encoder().McmsIncrementMintAllowanceWithArgs(typeArgs, stateObj, registryObj, denyListObj, executingCallbackParams)
			if err != nil {
				return nil, fmt.Errorf("failed to create mcms_entrypoint call: %w", err)
			}

			return overrideCall(entrypointCall, module, function), nil
		case "transfer_ownership":
			return managedToken.Encoder().McmsTransferOwnershipWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "execute_ownership_transfer":
			return managedToken.Encoder().McmsExecuteOwnershipTransferWithArgs(typeArgs, stateObj, registryObj, deployerStateObj, executingCallbackParams)
		case "destroy_managed_token":
			return managedToken.Encoder().McmsDestroyManagedTokenWithArgs(typeArgs, stateObj, registryObj, executingCallbackParams)
		case "add_allowed_modules":
			return managedToken.Encoder().McmsAddAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		case "remove_allowed_modules":
			return managedToken.Encoder().McmsRemoveAllowedModulesWithArgs(typeArgs, registryObj, executingCallbackParams)
		default:
			return nil, fmt.Errorf("unsupported managed_token MCMS function: %q", function)
		}
	}

	// Every supported module/function pair must be handled explicitly above: MCMS
	// callback signatures differ per function (extra objects, type arguments), so a
	// generic fallback would produce a plausible-looking but wrong call.
	return nil, fmt.Errorf("unsupported MCMS module: %q", module)
}
