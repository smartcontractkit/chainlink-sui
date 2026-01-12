package rmn

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type NoObjects struct{}

type CurseUncurseChainInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	Subjects         [][]byte
}

var CurseChainOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "curse_chain"),
	semver.MustParse("0.1.0"),
	"Curse a chain selector in the CCIP RMN Remote contract",
	curseChainHandler,
)

func curseChainHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input CurseUncurseChainInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	if len(input.Subjects) == 0 {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("at least one subject is required to curse")
	}

	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	encodedCall, err := contract.Encoder().CurseMultiple(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.Subjects,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode curse call: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to build transaction call for curse: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of curse_chain on RMN Remote as no signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute curse_chain on RMN Remote: %w", err)
	}

	b.Logger.Infow("Chains cursed on RMN Remote", "digest", tx.Digest, "count", len(input.Subjects))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var UncurseChainOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "uncurse_chain"),
	semver.MustParse("0.1.0"),
	"Uncurse a chain selector in the CCIP RMN Remote contract",
	uncurseChainHandler,
)

func uncurseChainHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input CurseUncurseChainInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	if len(input.Subjects) == 0 {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("at least one subject is required to uncurse")
	}

	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	encodedCall, err := contract.Encoder().UncurseMultiple(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.Subjects,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode uncurse call: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to build transaction call for uncurse: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of uncurse_chain on RMN Remote as no signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute uncurse_chain on RMN Remote: %w", err)
	}

	b.Logger.Infow("Chains uncursed on RMN Remote", "digest", tx.Digest, "count", len(input.Subjects))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}
