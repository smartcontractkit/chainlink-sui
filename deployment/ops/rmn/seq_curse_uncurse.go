package rmn

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// CurseUncurseSeqInput holds the context required by CurseSequence and UncurseSequence.
// It uses fastcurse.Subject (a [16]byte type alias) to bridge between the generic fastcurse
// interface and the RMN Remote contract's [][]byte representation.
type CurseUncurseSeqInput struct {
	CCIPAddress          string
	CCIPObjectRef        string
	CCIPOwnerCapObjectID string
	ChainSelector        uint64
	Subjects             []fastcurse.Subject
}

func executeCurseUncurse(
	b cldf_ops.Bundle,
	chains cldf_chain.BlockChains,
	in CurseUncurseSeqInput,
	op *cldf_ops.Operation[CurseUncurseChainInput, sui_ops.OpTxResult[NoObjects], sui_ops.OpTxDeps],
	opName string,
) (sequences.OnChainOutput, error) {
	chain, ok := chains.SuiChains()[in.ChainSelector]
	if !ok {
		return sequences.OnChainOutput{}, fmt.Errorf("Sui chain with selector %d not found in environment", in.ChainSelector)
	}

	subjectBytes := make([][]byte, len(in.Subjects))
	for i, subject := range in.Subjects {
		s := subject
		subjectBytes[i] = s[:]
	}

	deps := sui_ops.OpTxDeps{
		Client: chain.Client,
		Signer: chain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			gasBudget := uint64(400_000_000)
			return &bind.CallOpts{WaitForExecution: true, GasBudget: &gasBudget}
		},
		SuiRPC: chain.URL,
	}

	opInput := CurseUncurseChainInput{
		CCIPPackageId:    in.CCIPAddress,
		StateObjectId:    in.CCIPObjectRef,
		OwnerCapObjectId: in.CCIPOwnerCapObjectID,
		Subjects:         subjectBytes,
	}

	report, err := cldf_ops.ExecuteOperation(b, op, deps, opInput)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to execute %s operation on Sui chain %d: %w", opName, in.ChainSelector, err)
	}

	call := report.Output.Call
	tx, err := suisdk.NewTransactionWithStateObj(
		call.Module, call.Function, call.PackageID,
		call.Data, call.Module, []string{},
		call.StateObjID, call.TypeArgs,
	)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to create MCMS transaction: %w", err)
	}

	return sequences.OnChainOutput{
		BatchOps: []mcmstypes.BatchOperation{{
			ChainSelector: mcmstypes.ChainSelector(in.ChainSelector),
			Transactions:  []mcmstypes.Transaction{tx},
		}},
	}, nil
}

// CurseSequence curses the given subjects via RMN Remote on the specified Sui chain.
var CurseSequence = cldf_ops.NewSequence(
	"sui-curse-sequence",
	semver.MustParse("1.0.0"),
	"Curse sequence for Sui",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in CurseUncurseSeqInput) (sequences.OnChainOutput, error) {
		return executeCurseUncurse(b, chains, in, CurseChainOp, "curse")
	},
)

// UncurseSequence lifts the curse on given subjects via RMN Remote on the specified Sui chain.
var UncurseSequence = cldf_ops.NewSequence(
	"sui-uncurse-sequence",
	semver.MustParse("1.0.0"),
	"Uncurse sequence for Sui",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in CurseUncurseSeqInput) (sequences.OnChainOutput, error) {
		return executeCurseUncurse(b, chains, in, UncurseChainOp, "uncurse")
	},
)
