package adapters

import (
	"fmt"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	fees "github.com/smartcontractkit/chainlink-ccip/deployment/fees"
	lanes "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	datastore_utils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"

	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	suilanes "github.com/smartcontractkit/chainlink-sui/deployment/lanes"
)

var _ fees.FeeAdapter = &SuiFeeAdapter{}

// SuiFeeAdapter implements fees.FeeAdapter for Sui CCIP 1.6.0.
//
// On Sui the FeeQuoter is a module within the CCIP package (not a separate contract),
// so the FeeQuoter "address" is the CCIP package id and its state lives in the shared
// CCIPObjectRef, authorized by the CCIPOwnerCap. The fee adapter therefore resolves the
// FeeQuoter as the CCIP package ref. Token transfer fee configs are keyed on-chain by the
// coin metadata object address.
//
// Methods are being filled in incrementally; SetTokenTransferFee,
// GetOnchainTokenTransferFeeConfig, ApplyDestChainConfigUpdates, and GetOnchainDestChainConfig
// return nil or an error until they are wired up.
type SuiFeeAdapter struct{}

// GetFeeContractRef returns the FeeQuoter address ref for the source chain. On Sui the
// FeeQuoter is the CCIP package, so this resolves the CCIP package ref from the datastore
// rather than reading it from the OnRamp's dynamic config (as EVM/Solana do). The onRamp ref
// is validated for presence but is not used to derive the FeeQuoter.
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

// ================================================================
// === Stubs: not yet implemented                                ===
// ================================================================

func (a *SuiFeeAdapter) SetTokenTransferFee(_ datastore.DataStore, _ datastore.AddressRef) *cldf_ops.Sequence[fees.SetTokenTransferFeeSequenceInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

func (a *SuiFeeAdapter) GetOnchainTokenTransferFeeConfig(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, _ datastore.AddressRef, _ uint64, _ uint64, _ string) (fees.TokenTransferFeeArgs, error) {
	return fees.TokenTransferFeeArgs{}, fmt.Errorf("GetOnchainTokenTransferFeeConfig is not implemented on SuiFeeAdapter yet")
}

func (a *SuiFeeAdapter) ApplyDestChainConfigUpdates(_ datastore.DataStore, _ datastore.AddressRef) *cldf_ops.Sequence[fees.ApplyDestChainConfigSequenceInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

func (a *SuiFeeAdapter) GetOnchainDestChainConfig(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, _ datastore.AddressRef, _ uint64, _ uint64) (lanes.FeeQuoterDestChainConfig, error) {
	return lanes.FeeQuoterDestChainConfig{}, fmt.Errorf("GetOnchainDestChainConfig is not implemented on SuiFeeAdapter yet")
}
