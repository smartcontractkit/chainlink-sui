package adapters

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	datastore_utils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/datastore"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
)

var _ tokensapi.TokenAdminRegistryReader = (*SuiTokenAdminRegistryReader)(nil)

// SuiTokenAdminRegistryReader reads the Sui TokenAdminRegistry for the generic token changesets.
// The Sui TAR is a shared object accessed through the CCIP package and the CCIP object ref, and is
// keyed by coin metadata object id. Sui pools self-register in the TAR during pool initialization,
// so a registered token maps to its deployed pool package id.
type SuiTokenAdminRegistryReader struct{}

// GetActivePool returns the pool package id registered for tokenRef in the Sui TokenAdminRegistry,
// as raw bytes. Returns empty bytes with no error when no pool is registered for the token.
func (r *SuiTokenAdminRegistryReader) GetActivePool(
	e deployment.Environment, chainSelector uint64, tokenRef datastore.AddressRef, _ ...datastore.AddressRef,
) ([]byte, error) {
	chain, ok := e.BlockChains.SuiChains()[chainSelector]
	if !ok {
		return nil, fmt.Errorf("sui chain with selector %d not found", chainSelector)
	}
	coinType := tokenRef.Address
	if coinType == "" {
		ref, err := datastore_utils.FindAndFormatRef(e.DataStore, tokenRef, chainSelector, datastore_utils.FullRef)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve token ref for chain %d: %w", chainSelector, err)
		}
		coinType = ref.Address
	}
	if coinType == "" {
		return nil, fmt.Errorf("token ref has no coin type address on chain %d", chainSelector)
	}

	ccipPkg := firstRefAddress(findRefsByType(e.DataStore, chainSelector, datastore.ContractType(suideploy.SuiCCIPType)))
	if ccipPkg == "" {
		return nil, fmt.Errorf("CCIP package not found on chain %d", chainSelector)
	}
	ccipObjRef := firstRefAddress(findRefsByType(e.DataStore, chainSelector, datastore.ContractType(suideploy.SuiCCIPObjectRefType)))
	if ccipObjRef == "" {
		return nil, fmt.Errorf("CCIP object ref not found on chain %d", chainSelector)
	}

	// The TAR is keyed by coin metadata object id, so resolve the coin type to its CoinMetadata id.
	// The coin metadata read uses a direct RPC, so the nil-signer proposal deps suffice.
	deps := suiDeps(chain)
	coinMetaReport, err := cldf_ops.ExecuteOperation(e.OperationsBundle, coin_ops.GetCoinSymbolOp, deps, coinType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coin metadata for type %s on chain %d: %w", coinType, chainSelector, err)
	}
	coinMetaID := coinMetaReport.Output.Id
	if coinMetaID == "" {
		return nil, fmt.Errorf("coin metadata object id not found for type %s on chain %d", coinType, chainSelector)
	}

	contract, err := module_token_admin_registry.NewTokenAdminRegistry(ccipPkg, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}
	// The binding's DevInspect path builds a PTB and simulates it, which requires a signer
	// as the PTB sender even though no transaction is signed or submitted. Use the execution
	// deps and thread the signer onto the call opts, matching the opts.Signer = deps.Signer
	// idiom used by the token-pool ops.
	ctx := e.OperationsBundle.GetContext()
	execDeps := suiDepsExec(chain)
	opts := execDeps.GetCallOpts()
	opts.Signer = execDeps.Signer
	tarRef := bind.Object{Id: ccipObjRef}

	registered, err := contract.DevInspect().IsPoolRegistered(ctx, opts, tarRef, coinMetaID)
	if err != nil {
		return nil, fmt.Errorf("failed to check pool registration for %s on chain %d: %w", coinMetaID, chainSelector, err)
	}
	if !registered {
		return nil, nil
	}
	pool, err := contract.DevInspect().GetPool(ctx, opts, tarRef, coinMetaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get registered pool for %s on chain %d: %w", coinMetaID, chainSelector, err)
	}
	if pool == "" {
		return nil, nil
	}
	return suideploy.StrToBytes(pool)
}

// GetTokenAdminRegistryRef returns the Sui CCIP object ref, which is the shared object the
// token_admin_registry module functions are invoked with.
func (r *SuiTokenAdminRegistryReader) GetTokenAdminRegistryRef(
	e deployment.Environment, chainSelector uint64,
) (datastore.AddressRef, error) {
	if _, ok := e.BlockChains.SuiChains()[chainSelector]; !ok {
		return datastore.AddressRef{}, fmt.Errorf("sui chain with selector %d not found", chainSelector)
	}
	refs := e.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(chainSelector),
		datastore.AddressRefByType(datastore.ContractType(suideploy.SuiCCIPObjectRefType)),
	)
	if len(refs) == 0 {
		return datastore.AddressRef{}, fmt.Errorf("token admin registry (CCIP object ref) not found in datastore for chain %d", chainSelector)
	}
	return refs[0], nil
}
