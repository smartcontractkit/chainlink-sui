package ccipops

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type BackfillAllLocalDecimalsInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
}

type BackfillAllLocalDecimalsOutput struct {
	BackfilledTokens []string
}

var backfillAllLocalDecimalsHandler = func(
	b cld_ops.Bundle,
	deps sui_ops.OpTxDeps,
	input BackfillAllLocalDecimalsInput,
) (sui_ops.OpTxResult[BackfillAllLocalDecimalsOutput], error) {
	tokens, err := listConfiguredTokens(b.GetContext(), deps, input.CCIPPackageId, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[BackfillAllLocalDecimalsOutput]{}, err
	}

	backfilled := make([]string, 0, len(tokens))
	for _, token := range tokens {
		_, err := backfillLocalDecimalsHandler(b, deps, BackfillLocalDecimalsInput{
			CCIPPackageId:       input.CCIPPackageId,
			StateObjectId:       input.StateObjectId,
			OwnerCapObjectId:    input.OwnerCapObjectId,
			CoinMetadataAddress: token.CoinMetadataAddress,
			TokenType:           token.TokenType,
		})
		if err != nil {
			return sui_ops.OpTxResult[BackfillAllLocalDecimalsOutput]{}, fmt.Errorf(
				"failed to backfill local decimals for %s: %w",
				token.CoinMetadataAddress,
				err,
			)
		}
		backfilled = append(backfilled, token.CoinMetadataAddress)
	}

	return sui_ops.OpTxResult[BackfillAllLocalDecimalsOutput]{
		PackageId: input.CCIPPackageId,
		Objects: BackfillAllLocalDecimalsOutput{
			BackfilledTokens: backfilled,
		},
	}, nil
}

var TokenAdminRegistryBackfillAllLocalDecimalsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "backfill_all_local_decimals"),
	semver.MustParse("0.1.0"),
	"Backfills local token decimals for every registered token using on-chain coin metadata",
	backfillAllLocalDecimalsHandler,
)

type VerifyLocalDecimalsInput struct {
	CCIPPackageId string
	StateObjectId string
}

type VerifyLocalDecimalsMismatch struct {
	CoinMetadataAddress string
	TokenType           string
	RegistryDecimals    *byte
	ChainDecimals       byte
}

type VerifyLocalDecimalsOutput struct {
	CheckedTokens int
	Mismatches    []VerifyLocalDecimalsMismatch
}

var verifyLocalDecimalsHandler = func(
	b cld_ops.Bundle,
	deps sui_ops.OpTxDeps,
	input VerifyLocalDecimalsInput,
) (sui_ops.OpTxResult[VerifyLocalDecimalsOutput], error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[VerifyLocalDecimalsOutput]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	tokens, err := listConfiguredTokens(b.GetContext(), deps, input.CCIPPackageId, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[VerifyLocalDecimalsOutput]{}, err
	}

	opts := deps.GetCallOpts()
	ccipRef := bind.Object{Id: input.StateObjectId}
	mismatches := make([]VerifyLocalDecimalsMismatch, 0)

	for _, token := range tokens {
		chainDecimals, err := ResolveLocalDecimals(b.GetContext(), deps.Client, token.TokenType, nil)
		if err != nil {
			return sui_ops.OpTxResult[VerifyLocalDecimalsOutput]{}, fmt.Errorf(
				"failed to resolve on-chain decimals for %s: %w",
				token.CoinMetadataAddress,
				err,
			)
		}

		registryDecimals, err := contract.DevInspect().GetLocalDecimalsForToken(
			b.GetContext(),
			opts,
			ccipRef,
			token.CoinMetadataAddress,
		)
		if err != nil {
			mismatches = append(mismatches, VerifyLocalDecimalsMismatch{
				CoinMetadataAddress: token.CoinMetadataAddress,
				TokenType:           token.TokenType,
				RegistryDecimals:    nil,
				ChainDecimals:       chainDecimals,
			})
			continue
		}

		if registryDecimals != chainDecimals {
			reg := registryDecimals
			mismatches = append(mismatches, VerifyLocalDecimalsMismatch{
				CoinMetadataAddress: token.CoinMetadataAddress,
				TokenType:           token.TokenType,
				RegistryDecimals:    &reg,
				ChainDecimals:       chainDecimals,
			})
		}
	}

	return sui_ops.OpTxResult[VerifyLocalDecimalsOutput]{
		PackageId: input.CCIPPackageId,
		Objects: VerifyLocalDecimalsOutput{
			CheckedTokens: len(tokens),
			Mismatches:    mismatches,
		},
	}, nil
}

var TokenAdminRegistryVerifyLocalDecimalsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "verify_local_decimals"),
	semver.MustParse("0.1.0"),
	"Verifies registry local decimals match on-chain coin metadata for all registered tokens",
	verifyLocalDecimalsHandler,
)

// ValidateRegisterPoolLocalDecimals checks MCMS/owner registration input against on-chain metadata.
func ValidateRegisterPoolLocalDecimals(
	ctx context.Context,
	client coinMetadataClient,
	tokenType string,
	localDecimals byte,
) error {
	_, err := ResolveLocalDecimals(ctx, client, tokenType, &localDecimals)
	return err
}
