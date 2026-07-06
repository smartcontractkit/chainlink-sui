package rmn

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	suideployutils "github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

// CurseUncurseSeqInput holds the context required by CurseSequence and UncurseSequence.
// It uses fastcurse.Subject (a [16]byte type alias) to bridge between the generic fastcurse
// interface and the RMN Remote contract's [][]byte representation.
type CurseUncurseSeqInput struct {
	CCIPAddress          string
	LatestCCIPPackageID  string
	CCIPObjectRef        string
	CCIPOwnerCapObjectID string
	ChainSelector        uint64
	Subjects             []fastcurse.Subject
}

// FastCurseSeqInput holds the context required by FastCurseCurseSequence.
type FastCurseSeqInput struct {
	CCIPAddress         string
	LatestCCIPPackageID string
	CCIPObjectRef       string
	CurserCapObjectID   string
	ChainSelector       uint64
	Subjects            []fastcurse.Subject
}

func subjectsToBytes(subjects []fastcurse.Subject) [][]byte {
	out := make([][]byte, len(subjects))
	for i, subject := range subjects {
		out[i] = append([]byte(nil), subject[:]...)
	}
	return out
}

func batchOpFromCall(chainSelector uint64, call sui_ops.TransactionCall) (sequences.OnChainOutput, error) {
	tx, err := suideployutils.TransactionCallToMCMSTransaction(call)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to create MCMS transaction: %w", err)
	}
	return sequences.OnChainOutput{
		BatchOps: []mcmstypes.BatchOperation{{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  []mcmstypes.Transaction{tx},
		}},
	}, nil
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
		CCIPPackageId:       in.CCIPAddress,
		LatestCCIPPackageId: in.LatestCCIPPackageID,
		StateObjectId:       in.CCIPObjectRef,
		OwnerCapObjectId:    in.CCIPOwnerCapObjectID,
		Subjects:            subjectsToBytes(in.Subjects),
	}

	report, err := cldf_ops.ExecuteOperation(b, op, deps, opInput)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to execute %s operation on Sui chain %d: %w", opName, in.ChainSelector, err)
	}

	return batchOpFromCall(in.ChainSelector, report.Output.Call)
}
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

func executeFastCurse(
	b cldf_ops.Bundle,
	chains cldf_chain.BlockChains,
	in FastCurseSeqInput,
) (sequences.OnChainOutput, error) {
	chain, ok := chains.SuiChains()[in.ChainSelector]
	if !ok {
		return sequences.OnChainOutput{}, fmt.Errorf("Sui chain with selector %d not found in environment", in.ChainSelector)
	}
	if in.CurserCapObjectID == "" {
		return sequences.OnChainOutput{}, fmt.Errorf("curser cap object id is required for fast curse sequence")
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

	report, err := cldf_ops.ExecuteOperation(b, CurseWithCurserCapOp, deps, CurseWithCurserCapInput{
		CCIPPackageId:       in.CCIPAddress,
		LatestCCIPPackageId: in.LatestCCIPPackageID,
		StateObjectId:       in.CCIPObjectRef,
		CurserCapObjectId:   in.CurserCapObjectID,
		Subjects:            subjectsToBytes(in.Subjects),
	})
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to execute fast curse on Sui chain %d: %w", in.ChainSelector, err)
	}

	return batchOpFromCall(in.ChainSelector, report.Output.Call)
}

// FastCurseCurseSequence curses subjects via CurserCap held in the fast MCMS Registry.
var FastCurseCurseSequence = cldf_ops.NewSequence(
	"sui-fast-curse-sequence",
	semver.MustParse("1.0.0"),
	"Fast curse sequence for Sui via CurserCap",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in FastCurseSeqInput) (sequences.OnChainOutput, error) {
		return executeFastCurse(b, chains, in)
	},
)
