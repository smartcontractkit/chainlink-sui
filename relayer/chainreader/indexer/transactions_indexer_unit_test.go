package indexer

import (
	"testing"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"
)

func TestParseMoveAbortFromExecutionError(t *testing.T) {
	t.Parallel()

	kind := suirpcv2.ExecutionError_MOVE_ABORT
	command := uint64(1)
	functionName := "finish_execute"
	pkg := "0591c76156cdbf86aa4bcd42398c369c9d20151a23d9c5803f2b21e41f93ec61"
	module := "offramp"
	function := uint32(5)
	instruction := uint32(0)
	abortCode := uint64(1)

	execErr := &suirpcv2.ExecutionError{
		Kind:    &kind,
		Command: &command,
		ErrorDetails: &suirpcv2.ExecutionError_Abort{
			Abort: &suirpcv2.MoveAbort{
				AbortCode: &abortCode,
				Location: &suirpcv2.MoveLocation{
					Package:      &pkg,
					Module:       &module,
					Function:     &function,
					Instruction:  &instruction,
					FunctionName: &functionName,
				},
			},
		},
	}

	indexer := &TransactionsIndexer{}
	moveAbort, err := indexer.parseMoveAbortFromExecutionError(execErr)
	require.NoError(t, err)
	require.NotNil(t, moveAbort)
	require.Equal(t, uint64(1), moveAbort.CommandIndex)
	require.Equal(t, uint64(1), moveAbort.AbortCode)
	require.Equal(t, pkg, moveAbort.Location.Module.Address)
	require.Equal(t, module, moveAbort.Location.Module.Name)
	require.Equal(t, "finish_execute", *moveAbort.Location.FunctionName)
}

func TestParseMoveAbortFromExecutionError_RejectsMissingAbort(t *testing.T) {
	t.Parallel()

	description := "MoveAbort(MoveLocation { module: ModuleId { address: abc, name: Identifier(\"offramp\") }, function: 1, instruction: 0, function_name: Some(\"finish_execute\") }, 1) in command 1"
	execErr := &suirpcv2.ExecutionError{
		Description: &description,
	}

	indexer := &TransactionsIndexer{}
	_, err := indexer.parseMoveAbortFromExecutionError(execErr)
	require.Error(t, err)
}

func TestBatchContainsConfigSet(t *testing.T) {
	t.Parallel()

	pkg := "0xabc"
	txWithConfigSet := &suirpcv2.ExecutedTransaction{
		Events: &suirpcv2.TransactionEvents{
			Events: []*suirpcv2.Event{
				{EventType: strPtr("0xdef::offramp::ExecutionStateChanged")},
				{EventType: strPtr("0xabc::ocr3_base::ConfigSet")},
			},
		},
	}
	txWithoutConfigSet := &suirpcv2.ExecutedTransaction{
		Events: &suirpcv2.TransactionEvents{
			Events: []*suirpcv2.Event{
				{EventType: strPtr("0xabc::offramp::ExecutionStateChanged")},
			},
		},
	}
	txNoEvents := &suirpcv2.ExecutedTransaction{}

	require.True(t, batchContainsConfigSet(
		[]*suirpcv2.ExecutedTransaction{txWithoutConfigSet, txWithConfigSet},
		pkg, "ocr3_base", "ConfigSet",
	))
	require.False(t, batchContainsConfigSet(
		[]*suirpcv2.ExecutedTransaction{txWithoutConfigSet},
		pkg, "ocr3_base", "ConfigSet",
	))
	require.False(t, batchContainsConfigSet(
		[]*suirpcv2.ExecutedTransaction{txNoEvents},
		pkg, "ocr3_base", "ConfigSet",
	))
	require.False(t, batchContainsConfigSet(
		nil,
		pkg, "ocr3_base", "ConfigSet",
	))
}
