package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type BackfillLocalDecimalsSeqInput struct {
	CCIPPackageId string // effective/latest CCIP package: used for reads (verify) and direct execution
	// OriginalCCIPPackageId is the original/defining CCIP package, used as the MCMS on-chain identity
	// (proposal target) for the proposal-routed backfill ops. Falls back to CCIPPackageId when empty.
	OriginalCCIPPackageId string
	StateObjectId         string
	OwnerCapObjectId      string
	VerifyOnly            bool
	// ProposalOnly encodes MCMS backfill ops without executing them. Verify still uses deps.Signer.
	ProposalOnly bool
}

type BackfillLocalDecimalsSeqOutput struct {
	Reports    []cld_ops.Report[any, any]
	McmsDefs   []cld_ops.Definition
	McmsInputs []any
}

var BackfillLocalDecimalsSequence = cld_ops.NewSequence(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "backfill_local_decimals"),
	semver.MustParse("0.1.0"),
	"Verifies and backfills local token decimals, producing MCMS proposal inputs when ProposalOnly is set",
	func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BackfillLocalDecimalsSeqInput) (BackfillLocalDecimalsSeqOutput, error) {
		reports := make([]cld_ops.Report[any, any], 0)

		verifyReport, err := cld_ops.ExecuteOperation(
			b,
			TokenAdminRegistryVerifyLocalDecimalsOp,
			deps,
			VerifyLocalDecimalsInput{
				CCIPPackageId: input.CCIPPackageId,
				StateObjectId: input.StateObjectId,
			},
		)
		if err != nil {
			return BackfillLocalDecimalsSeqOutput{}, fmt.Errorf("failed to verify local decimals: %w", err)
		}
		reports = append(reports, verifyReport.ToGenericReport())

		mismatches := verifyReport.Output.Objects.Mismatches
		if len(mismatches) > 0 && input.VerifyOnly {
			return BackfillLocalDecimalsSeqOutput{Reports: reports}, fmt.Errorf(
				"local decimals verification failed for %d token(s)",
				len(mismatches),
			)
		}
		if input.VerifyOnly || len(mismatches) == 0 {
			return BackfillLocalDecimalsSeqOutput{Reports: reports}, nil
		}

		if input.ProposalOnly || deps.Signer == nil {
			// MCMS-routed: identity must be the original package; PTB dispatches against the latest.
			originalPkgId := input.OriginalCCIPPackageId
			if originalPkgId == "" {
				originalPkgId = input.CCIPPackageId
			}
			encodeDeps := deps
			if input.ProposalOnly {
				encodeDeps.Signer = nil
			}
			defs := make([]cld_ops.Definition, 0, len(mismatches))
			inputs := make([]any, 0, len(mismatches))
			for _, mismatch := range mismatches {
				report, err := cld_ops.ExecuteOperation(
					b,
					TokenAdminRegistryBackfillLocalDecimalsOp,
					encodeDeps,
					BackfillLocalDecimalsInput{
						CCIPPackageId:       originalPkgId,       // original = MCMS on-chain identity
						LatestPackageId:     input.CCIPPackageId, // upgraded = PTB dispatch binary
						StateObjectId:       input.StateObjectId,
						OwnerCapObjectId:    input.OwnerCapObjectId,
						CoinMetadataAddress: mismatch.CoinMetadataAddress,
						TokenType:           mismatch.TokenType,
					},
				)
				if err != nil {
					return BackfillLocalDecimalsSeqOutput{Reports: reports}, fmt.Errorf(
						"failed to encode backfill local decimals for %s: %w",
						mismatch.CoinMetadataAddress,
						err,
					)
				}
				defs = append(defs, report.Def)
				inputs = append(inputs, report.Input)
				reports = append(reports, report.ToGenericReport())
			}
			return BackfillLocalDecimalsSeqOutput{
				Reports:    reports,
				McmsDefs:   defs,
				McmsInputs: inputs,
			}, nil
		}

		backfillReport, err := cld_ops.ExecuteOperation(
			b,
			TokenAdminRegistryBackfillAllLocalDecimalsOp,
			deps,
			BackfillAllLocalDecimalsInput{
				CCIPPackageId:    input.CCIPPackageId,
				StateObjectId:    input.StateObjectId,
				OwnerCapObjectId: input.OwnerCapObjectId,
			},
		)
		if err != nil {
			return BackfillLocalDecimalsSeqOutput{Reports: reports}, fmt.Errorf("failed to backfill local decimals: %w", err)
		}
		reports = append(reports, backfillReport.ToGenericReport())

		verifyAfterReport, err := cld_ops.ExecuteOperation(
			b,
			TokenAdminRegistryVerifyLocalDecimalsOp,
			deps,
			VerifyLocalDecimalsInput{
				CCIPPackageId: input.CCIPPackageId,
				StateObjectId: input.StateObjectId,
			},
		)
		if err != nil {
			return BackfillLocalDecimalsSeqOutput{Reports: reports}, fmt.Errorf("failed to re-verify local decimals: %w", err)
		}
		reports = append(reports, verifyAfterReport.ToGenericReport())

		if len(verifyAfterReport.Output.Objects.Mismatches) > 0 {
			return BackfillLocalDecimalsSeqOutput{Reports: reports}, fmt.Errorf(
				"local decimals still mismatched for %d token(s) after backfill",
				len(verifyAfterReport.Output.Objects.Mismatches),
			)
		}

		return BackfillLocalDecimalsSeqOutput{Reports: reports}, nil
	},
)
