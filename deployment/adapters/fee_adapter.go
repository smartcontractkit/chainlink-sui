package adapters

import (
	"fmt"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	fees "github.com/smartcontractkit/chainlink-ccip/deployment/fees"
	lanes "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	datastore_utils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	suilanes "github.com/smartcontractkit/chainlink-sui/deployment/lanes"
)

var _ fees.FeeAdapter = &SuiFeeAdapter{}

// Label keys stashed on the FeeQuoter ref by GetFeeContractRef so the downstream
// SetTokenTransferFee / GetOnchainTokenTransferFeeConfig can recover the CCIP state object
// ref, owner cap, and upgraded package id without a datastore of their own. Format: "<key>:<id>".
const (
	suiFeeCCIPObjectRefLabel = "ccip-object-ref"
	suiFeeCCIPOwnerCapLabel  = "ccip-owner-cap"
	suiFeeLatestCCIPPkgLabel = "ccip-latest-package-id"
)

// SuiFeeAdapter implements fees.FeeAdapter for Sui CCIP 1.6.0.
//
// On Sui the FeeQuoter is a module within the CCIP package (not a separate contract), so the
// FeeQuoter "address" is the CCIP package id, its state lives in the shared CCIPObjectRef, and
// it is authorized by the CCIPOwnerCap. Token transfer fee configs are keyed on-chain by the
// coin metadata object address. GetFeeContractRef returns the CCIP package ref versioned at
// 1.6.0 (the adapter version; Sui contract refs are stored at 1.0.0) and stashes the
// CCIPObjectRef / CCIPOwnerCap / latest package id as labels so the FQ ops and DevInspect reads
// have everything they need.
//
// ApplyDestChainConfigUpdates and GetOnchainDestChainConfig return nil or an error until wired.
type SuiFeeAdapter struct{}

// GetFeeContractRef returns the FeeQuoter address ref for the source chain. On Sui the
// FeeQuoter is the CCIP package, so this resolves the CCIP package ref from the datastore
// (not via an OnRamp dynamic-config read like EVM/Solana), enriches it with the CCIPObjectRef,
// CCIPOwnerCap, and optional latest CCIP package id as labels, and versions it at 1.6.0 so the
// generic fee flow selects this adapter. The onRamp ref is validated for presence.
func (a *SuiFeeAdapter) GetFeeContractRef(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, ds datastore.DataStore, onRamp datastore.AddressRef, src uint64, dst uint64) (datastore.AddressRef, error) {
	if onRamp.Address == "" {
		return datastore.AddressRef{}, fmt.Errorf("onRamp ref has empty address for src %d dst %d", src, dst)
	}
	fqRef, err := datastore_utils.FindAndFormatRef(ds, datastore.AddressRef{
		ChainSelector: src,
		Type:          datastore.ContractType(suideploy.SuiCCIPType),
	}, src, datastore_utils.FullRef)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("failed to find Sui CCIP package (FeeQuoter) ref on chain %d: %w", src, err)
	}
	ccipObjRef := firstRefAddress(findRefsByType(ds, src, datastore.ContractType(suideploy.SuiCCIPObjectRefType)))
	if ccipObjRef == "" {
		return datastore.AddressRef{}, fmt.Errorf("CCIP object ref not found on chain %d", src)
	}
	ownerCap := firstRefAddress(findRefsByType(ds, src, datastore.ContractType(suideploy.SuiCCIPOwnerCapObjectIDType)))
	if ownerCap == "" {
		return datastore.AddressRef{}, fmt.Errorf("CCIP owner cap not found on chain %d", src)
	}
	latestPkg := firstRefAddress(findRefsByType(ds, src, datastore.ContractType(suideploy.SuiLatestCCIPPackageIDType)))

	labels := append(fqRef.Labels.List(),
		suiFeeCCIPObjectRefLabel+":"+ccipObjRef,
		suiFeeCCIPOwnerCapLabel+":"+ownerCap,
	)
	if latestPkg != "" {
		labels = append(labels, suiFeeLatestCCIPPkgLabel+":"+latestPkg)
	}
	fqRef.Labels = datastore.NewLabelSet(labels...)
	fqRef.Version = semver.MustParse("1.6.0")
	return fqRef, nil
}

// GetDefaultTokenTransferFeeConfig returns the chain-agnostic default token transfer fee
// configuration, matching the EVM and Solana fee adapters.
func (a *SuiFeeAdapter) GetDefaultTokenTransferFeeConfig(src uint64, dst uint64) fees.TokenTransferFeeArgs {
	return fees.GetDefaultChainAgnosticTokenTransferFeeConfig(src, dst)
}

// GetDefaultDestChainConfig returns the Sui FeeQuoter destination chain config defaults.
// It delegates to the Sui lane adapter so the defaults have a single source of truth.
func (a *SuiFeeAdapter) GetDefaultDestChainConfig(_, _ uint64) lanes.FeeQuoterDestChainConfig {
	return (&suilanes.SuiAdapter{}).GetFeeQuoterDestChainConfig()
}

// SetTokenTransferFee sets token transfer fee config on the Sui FeeQuoter as an MCMS proposal.
// Per destination chain it maps the generic Settings into add/remove token updates and invokes
// FeeQuoterApplyTokenTransferFeeConfigUpdates in proposal mode, bridging the encoded call into
// OnChainOutput.BatchOps. Token keys are coin metadata object addresses.
func (a *SuiFeeAdapter) SetTokenTransferFee(_ datastore.DataStore, feeRef datastore.AddressRef) *cldf_ops.Sequence[fees.SetTokenTransferFeeSequenceInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"sui-fee-adapter:set-token-transfer-fee",
		semver.MustParse("1.6.0"),
		"Set token transfer fee config on the Sui FeeQuoter as an MCMS proposal",
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, input fees.SetTokenTransferFeeSequenceInput) (sequences.OnChainOutput, error) {
			chain, ok := chains.SuiChains()[input.Selector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("sui chain with selector %d not found", input.Selector)
			}
			ccipPkg := feeRef.Address
			if ccipPkg == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("fee ref has empty CCIP package address on chain %d", input.Selector)
			}
			ccipObjRef := feeRefLabelValue(feeRef, suiFeeCCIPObjectRefLabel)
			ownerCap := feeRefLabelValue(feeRef, suiFeeCCIPOwnerCapLabel)
			if ccipObjRef == "" || ownerCap == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("fee ref is missing ccip-object-ref/ccip-owner-cap labels on chain %d", input.Selector)
			}
			latestPkg := feeRefLabelValue(feeRef, suiFeeLatestCCIPPkgLabel)
			deps := suiDeps(chain)

			batchOps := make([]mcmstypes.BatchOperation, 0)
			for dst, dstCfg := range input.Settings {
				update, err := buildSuiTokenTransferFeeUpdate(dstCfg)
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("token transfer fee config for dst %d: %w", dst, err)
				}
				if len(update.AddTokens) == 0 && len(update.RemoveTokens) == 0 {
					b.Logger.Warnf("sui SetTokenTransferFee: no token updates for dst %d; skipping", dst)
					continue
				}
				r, err := cldf_ops.ExecuteOperation(b, ccipops.FeeQuoterApplyTokenTransferFeeConfigUpdatesOp, deps, ccipops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
					CCIPPackageId:        ccipPkg,
					LatestPackageId:      latestPkg,
					StateObjectId:        ccipObjRef,
					OwnerCapObjectId:     ownerCap,
					DestChainSelector:    dst,
					AddTokens:            update.AddTokens,
					AddMinFeeUsdCents:    update.AddMinFeeUsdCents,
					AddMaxFeeUsdCents:    update.AddMaxFeeUsdCents,
					AddDeciBps:           update.AddDeciBps,
					AddDestGasOverhead:   update.AddDestGasOverhead,
					AddDestBytesOverhead: update.AddDestBytesOverhead,
					AddIsEnabled:         update.AddIsEnabled,
					RemoveTokens:         update.RemoveTokens,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to set token transfer fee for dst %d: %w", dst, err)
				}
				out, err := batchOpFromCall(input.Selector, r.Output.Call)
				if err != nil {
					return sequences.OnChainOutput{}, err
				}
				batchOps = append(batchOps, out.BatchOps...)
			}
			return sequences.OnChainOutput{BatchOps: batchOps}, nil
		},
	)
}

// GetOnchainTokenTransferFeeConfig reads a token's transfer fee config from the Sui FeeQuoter
// via DevInspect. The token arg is the coin metadata object address. Returns a zero config
// (matching the Move behavior) when no custom config is set on-chain.
func (a *SuiFeeAdapter) GetOnchainTokenTransferFeeConfig(b cldf_ops.Bundle, chains cldf_chain.BlockChains, feeRef datastore.AddressRef, src uint64, dst uint64, token string) (fees.TokenTransferFeeArgs, error) {
	chain, ok := chains.SuiChains()[src]
	if !ok {
		return fees.TokenTransferFeeArgs{}, fmt.Errorf("sui chain with selector %d not found", src)
	}
	if token == "" {
		return fees.TokenTransferFeeArgs{}, fmt.Errorf("token (coin metadata address) is empty for src %d dst %d", src, dst)
	}
	ccipPkg := feeRef.Address
	if ccipPkg == "" {
		return fees.TokenTransferFeeArgs{}, fmt.Errorf("fee ref has empty CCIP package address for src %d dst %d", src, dst)
	}
	ccipObjRef := feeRefLabelValue(feeRef, suiFeeCCIPObjectRefLabel)
	if ccipObjRef == "" {
		return fees.TokenTransferFeeArgs{}, fmt.Errorf("fee ref is missing ccip-object-ref label for src %d dst %d", src, dst)
	}
	// DevInspect reads current state, so target the upgraded package head when present.
	readPkg := feeRefLabelValue(feeRef, suiFeeLatestCCIPPkgLabel)
	if readPkg == "" {
		readPkg = ccipPkg
	}
	contract, err := module_fee_quoter.NewFeeQuoter(readPkg, chain.Client)
	if err != nil {
		return fees.TokenTransferFeeArgs{}, fmt.Errorf("failed to instantiate FeeQuoter at %s on chain %d: %w", readPkg, src, err)
	}
	cfg, err := contract.DevInspect().GetTokenTransferFeeConfig(b.GetContext(), &bind.CallOpts{Signer: chain.Signer}, bind.Object{Id: ccipObjRef}, dst, token)
	if err != nil {
		return fees.TokenTransferFeeArgs{}, fmt.Errorf("failed to read token transfer fee config from FeeQuoter at %s for src %d dst %d token %s: %w", readPkg, src, dst, token, err)
	}
	b.Logger.Infof("Fetched on-chain token transfer fee config for src %d, dst %d, token %s: %+v", src, dst, token, cfg)
	return fqConfigToArgs(cfg), nil
}

// ================================================================
// === Stubs: not yet implemented                                ===
// ================================================================

func (a *SuiFeeAdapter) ApplyDestChainConfigUpdates(_ datastore.DataStore, _ datastore.AddressRef) *cldf_ops.Sequence[fees.ApplyDestChainConfigSequenceInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

func (a *SuiFeeAdapter) GetOnchainDestChainConfig(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, _ datastore.AddressRef, _ uint64, _ uint64) (lanes.FeeQuoterDestChainConfig, error) {
	return lanes.FeeQuoterDestChainConfig{}, fmt.Errorf("GetOnchainDestChainConfig is not implemented on SuiFeeAdapter yet")
}

// ================================================================
// === Helpers                                                   ===
// ================================================================

// feeRefLabelValue extracts the value following "<key>:" from the ref's labels, or empty.
func feeRefLabelValue(ref datastore.AddressRef, key string) string {
	prefix := key + ":"
	for _, l := range ref.Labels.List() {
		if strings.HasPrefix(l, prefix) {
			return strings.TrimPrefix(l, prefix)
		}
	}
	return ""
}

// suiTokenTransferFeeUpdate holds the add/remove token update slices for one destination chain.
type suiTokenTransferFeeUpdate struct {
	AddTokens            []string
	AddMinFeeUsdCents    []uint32
	AddMaxFeeUsdCents    []uint32
	AddDeciBps           []uint16
	AddDestGasOverhead   []uint32
	AddDestBytesOverhead []uint32
	AddIsEnabled         []bool
	RemoveTokens         []string
}

// buildSuiTokenTransferFeeUpdate maps the generic per-destination token fee config into the
// Sui FeeQuoter add/remove parallel slices. Enabled tokens are added with their config; nil or
// disabled tokens are removed. The Sui FeeQuoter takes removals as a separate list, so no
// default-fill is needed (unlike Solana).
func buildSuiTokenTransferFeeUpdate(dstCfg map[string]*fees.TokenTransferFeeArgs) (suiTokenTransferFeeUpdate, error) {
	var u suiTokenTransferFeeUpdate
	for token, feeCfg := range dstCfg {
		if token == "" {
			return suiTokenTransferFeeUpdate{}, fmt.Errorf("empty token address in token transfer fee config")
		}
		if feeCfg == nil || !feeCfg.IsEnabled {
			u.RemoveTokens = append(u.RemoveTokens, token)
			continue
		}
		u.AddTokens = append(u.AddTokens, token)
		u.AddMinFeeUsdCents = append(u.AddMinFeeUsdCents, feeCfg.MinFeeUSDCents)
		u.AddMaxFeeUsdCents = append(u.AddMaxFeeUsdCents, feeCfg.MaxFeeUSDCents)
		u.AddDeciBps = append(u.AddDeciBps, feeCfg.DeciBps)
		u.AddDestGasOverhead = append(u.AddDestGasOverhead, feeCfg.DestGasOverhead)
		u.AddDestBytesOverhead = append(u.AddDestBytesOverhead, feeCfg.DestBytesOverhead)
		u.AddIsEnabled = append(u.AddIsEnabled, true)
	}
	return u, nil
}

// fqConfigToArgs maps the Sui FeeQuoter TokenTransferFeeConfig onto the chain-agnostic args.
func fqConfigToArgs(cfg module_fee_quoter.TokenTransferFeeConfig) fees.TokenTransferFeeArgs {
	return fees.TokenTransferFeeArgs{
		MinFeeUSDCents:    cfg.MinFeeUsdCents,
		MaxFeeUSDCents:    cfg.MaxFeeUsdCents,
		DeciBps:           cfg.DeciBps,
		DestGasOverhead:   cfg.DestGasOverhead,
		DestBytesOverhead: cfg.DestBytesOverhead,
		IsEnabled:         cfg.IsEnabled,
	}
}
