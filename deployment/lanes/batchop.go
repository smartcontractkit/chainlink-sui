package lanes

import (
	"fmt"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func appendMCMSBatchOpFromCall(
	out *sequences.OnChainOutput,
	chainSelector uint64,
	call sui_ops.TransactionCall,
	deps sui_ops.OpTxDeps,
) error {
	if deps.Signer != nil {
		return nil
	}
	if call.PackageID == "" {
		return nil
	}

	tx, err := utils.TransactionCallToMCMSTransaction(call)
	if err != nil {
		return fmt.Errorf("create MCMS transaction from op call: %w", err)
	}

	out.BatchOps = append(out.BatchOps, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(chainSelector),
		Transactions:  []mcmstypes.Transaction{tx},
	})
	return nil
}
