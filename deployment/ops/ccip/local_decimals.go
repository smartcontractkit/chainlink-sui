package ccipops

import (
	"context"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type coinMetadataClient interface {
	GetCoinMetadata(ctx context.Context, coinType string) (models.CoinMetadataResponse, error)
}

// ResolveLocalDecimals returns on-chain decimals for a token type and optionally
// validates a caller-supplied value.
func ResolveLocalDecimals(
	ctx context.Context,
	client coinMetadataClient,
	tokenType string,
	supplied *byte,
) (byte, error) {
	if tokenType == "" {
		return 0, fmt.Errorf("token type is required to resolve local decimals")
	}

	coinType := normalizeCoinType(tokenType)
	metadata, err := client.GetCoinMetadata(ctx, coinType)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch coin metadata for type %s: %w", coinType, err)
	}
	if metadata.Decimals < 0 || metadata.Decimals > 255 {
		return 0, fmt.Errorf("invalid decimals %d for coin type %s", metadata.Decimals, coinType)
	}

	chainDecimals := byte(metadata.Decimals)
	if supplied != nil && *supplied != chainDecimals {
		return 0, fmt.Errorf(
			"local decimals mismatch for %s: supplied %d, on-chain metadata %d",
			coinType,
			*supplied,
			chainDecimals,
		)
	}

	return chainDecimals, nil
}

func normalizeCoinType(tokenType string) string {
	if strings.HasPrefix(tokenType, "0x") {
		return tokenType
	}
	return "0x" + tokenType
}

type configuredToken struct {
	CoinMetadataAddress string
	TokenType           string
}

func listConfiguredTokens(
	ctx context.Context,
	deps sui_ops.OpTxDeps,
	ccipPackageID string,
	ccipRefObjectID string,
) ([]configuredToken, error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(ccipPackageID, deps.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	opts := callOptsWithSigner(deps)
	ccipRef := bind.Object{Id: ccipRefObjectID}
	startKey := "0x0"
	const pageSize uint64 = 1000

	tokens := make([]configuredToken, 0)
	for {
		raw, err := contract.DevInspect().GetAllConfiguredTokens(ctx, opts, ccipRef, startKey, pageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to list configured tokens: %w", err)
		}
		if len(raw) < 3 {
			return nil, fmt.Errorf("unexpected GetAllConfiguredTokens response length %d", len(raw))
		}

		addresses, ok := raw[0].([]string)
		if !ok {
			return nil, fmt.Errorf("unexpected token address list type %T", raw[0])
		}

		nextKey, ok := raw[1].(string)
		if !ok {
			return nil, fmt.Errorf("unexpected next key type %T", raw[1])
		}

		hasMore, ok := raw[2].(bool)
		if !ok {
			return nil, fmt.Errorf("unexpected has_more type %T", raw[2])
		}

		for _, coinMetadataAddress := range addresses {
			config, err := contract.DevInspect().GetTokenConfigStruct(ctx, opts, ccipRef, coinMetadataAddress)
			if err != nil {
				return nil, fmt.Errorf("failed to get token config for %s: %w", coinMetadataAddress, err)
			}
			tokens = append(tokens, configuredToken{
				CoinMetadataAddress: coinMetadataAddress,
				TokenType:           config.TokenType,
			})
		}

		if !hasMore {
			break
		}
		startKey = nextKey
	}

	return tokens, nil
}

func callOptsWithSigner(deps sui_ops.OpTxDeps) *bind.CallOpts {
	opts := deps.GetCallOpts()
	if deps.Signer != nil {
		opts.Signer = deps.Signer
	}
	return opts
}
