package sui

import (
	"context"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	tokenadminregistry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// TokenPoolConfig holds the runtime-resolved token pool configuration.
type TokenPoolConfig struct {
	PoolPackageID       string
	PoolModule          string
	PoolKind            string
	PoolStateObjectID   string
	CoinType            string
	LatestPoolPackageID string
	LockOrBurnParams    []string
}

// ResolveTokenPoolConfig resolves the token pool configuration from the CCIP token admin registry.
//
// ccipPkgID is the original CCIP package ID from the address book.
// ccipObjectRefID is the shared CCIPObjectRef object ID.
// coinMetadataID is the token's CoinMetadata object ID.
func ResolveTokenPoolConfig(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	ccipPkgID string,
	ccipObjectRefID string,
	coinMetadataID string,
) (*TokenPoolConfig, error) {
	latestCcipPkg := ccipPkgID
	if pkg, err := resolveLatestPackageIDFromStateObject(ctx, ptbClient, ccipObjectRefID); err == nil && pkg != "" {
		latestCcipPkg = pkg
	}
	if fallback, err := ptbClient.GetLatestPackageId(ctx, latestCcipPkg, "ccip"); err == nil && fallback != "" {
		latestCcipPkg = fallback
	}

	registryContract, err := tokenadminregistry.NewTokenAdminRegistry(latestCcipPkg, ptbClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	tokenConfig, err := registryContract.DevInspect().GetTokenConfigStruct(
		ctx,
		&bind.CallOpts{Signer: signer},
		bind.Object{Id: ccipObjectRefID},
		coinMetadataID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get token config struct for coin %s: %w", coinMetadataID, err)
	}

	poolKind := derivePoolKind(tokenConfig.TokenPoolModule)
	poolStateObjectID := derivePoolStateObjectID(poolKind, tokenConfig.LockOrBurnParams)

	cfg := &TokenPoolConfig{
		PoolPackageID:     tokenConfig.TokenPoolPackageId,
		PoolModule:        tokenConfig.TokenPoolModule,
		PoolKind:          poolKind,
		PoolStateObjectID: poolStateObjectID,
		CoinType:          normalizeSuiTypeTag(tokenConfig.TokenType),
		LockOrBurnParams:  tokenConfig.LockOrBurnParams,
	}

	if pkg, err := resolveLatestPackageIDFromStateObject(ctx, ptbClient, cfg.PoolStateObjectID); err == nil && pkg != "" {
		cfg.LatestPoolPackageID = pkg
	} else {
		cfg.LatestPoolPackageID = cfg.PoolPackageID
	}
	if fallback, err := ptbClient.GetLatestPackageId(ctx, cfg.LatestPoolPackageID, cfg.PoolModule); err == nil && fallback != "" {
		cfg.LatestPoolPackageID = fallback
	}

	return cfg, nil
}

func normalizeSuiTypeTag(typeTag string) string {
	trimmed := strings.TrimSpace(typeTag)
	if trimmed == "" {
		return trimmed
	}
	if strings.HasPrefix(trimmed, "0x") || strings.HasPrefix(trimmed, "0X") {
		return trimmed
	}
	return "0x" + trimmed
}

func derivePoolKind(moduleName string) string {
	switch {
	case strings.Contains(moduleName, "managed_token_pool"):
		return "managed"
	case strings.Contains(moduleName, "lock_release_token_pool"):
		return "lock_release"
	case strings.Contains(moduleName, "burn_mint_token_pool"):
		return "burn_mint"
	default:
		return "unknown"
	}
}

func derivePoolStateObjectID(poolKind string, lockOrBurnParams []string) string {
	if len(lockOrBurnParams) == 0 {
		return ""
	}
	switch poolKind {
	case "managed":
		if len(lockOrBurnParams) > 3 {
			return lockOrBurnParams[3]
		}
	default:
		if len(lockOrBurnParams) > 1 {
			return lockOrBurnParams[1]
		}
	}
	return ""
}
