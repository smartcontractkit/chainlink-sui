package indexer

import (
	"testing"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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

func gateTestMoveCallCommand(pkg, module, function string, numInputArgs int) *suirpcv2.Command {
	call := &suirpcv2.MoveCall{
		Package:  &pkg,
		Module:   &module,
		Function: &function,
	}
	for i := range numInputArgs {
		argKind := suirpcv2.Argument_INPUT
		inputIdx := uint32(i) //nolint:gosec
		call.Arguments = append(call.Arguments, &suirpcv2.Argument{Kind: &argKind, Input: &inputIdx})
	}

	return &suirpcv2.Command{Command: &suirpcv2.Command_MoveCall{MoveCall: call}}
}

func gateTestFailedTx(commands []*suirpcv2.Command, inputs []*suirpcv2.Input, abortCommand uint64, abortFunctionName *string) *suirpcv2.ExecutedTransaction {
	errKind := suirpcv2.ExecutionError_MOVE_ABORT
	abortCode := uint64(1)
	abortPkg := "0591c76156cdbf86aa4bcd42398c369c9d20151a23d9c5803f2b21e41f93ec61"
	abortModule := "offramp"
	abortFunction := uint32(5)
	abortInstruction := uint32(0)
	digest := "test-digest"
	success := false

	return &suirpcv2.ExecutedTransaction{
		Digest: &digest,
		Transaction: &suirpcv2.Transaction{
			Kind: &suirpcv2.TransactionKind{
				Data: &suirpcv2.TransactionKind_ProgrammableTransaction{
					ProgrammableTransaction: &suirpcv2.ProgrammableTransaction{
						Inputs:   inputs,
						Commands: commands,
					},
				},
			},
		},
		Effects: &suirpcv2.TransactionEffects{
			Status: &suirpcv2.ExecutionStatus{
				Success: &success,
				Error: &suirpcv2.ExecutionError{
					Kind:    &errKind,
					Command: &abortCommand,
					ErrorDetails: &suirpcv2.ExecutionError_Abort{
						Abort: &suirpcv2.MoveAbort{
							AbortCode: &abortCode,
							Location: &suirpcv2.MoveLocation{
								Package:      &abortPkg,
								Module:       &abortModule,
								Function:     &abortFunction,
								Instruction:  &abortInstruction,
								FunctionName: abortFunctionName,
							},
						},
					},
				},
			},
		},
	}
}

func TestProcessFailedTransaction_Gate(t *testing.T) {
	t.Parallel()

	const (
		offrampPkg = "0xofframp"
		latestPkg  = "0xlatestofframp"
		otherPkg   = "0xother"
	)

	strPtr := func(s string) *string { return &s }

	initExecuteCmd := gateTestMoveCallCommand(offrampPkg, "offramp", "init_execute", 5)
	finishExecuteCmd := gateTestMoveCallCommand(offrampPkg, "offramp", "finish_execute", 0)

	// Five INPUT args for init_execute; arg position 4 holds the report as a Pure
	// input, filled with bytes that fail BCS deserialization.
	garbageInputs := []*suirpcv2.Input{{}, {}, {}, {}, {Pure: []byte{0xff, 0x01, 0x02}}}

	tests := []struct {
		name            string
		commands        []*suirpcv2.Command
		inputs          []*suirpcv2.Input
		abortCommand    uint64
		abortFunction   *string
		wantErrContains string // empty means expect (nil, nil) with no error
	}{
		{
			name:          "abort in command before init_execute is skipped",
			commands:      []*suirpcv2.Command{gateTestMoveCallCommand(otherPkg, "attacker", "fail_fast", 0), initExecuteCmd},
			abortCommand:  0,
			abortFunction: strPtr("fail_fast"),
		},
		{
			name:          "abort inside init_execute callee at same command index is skipped",
			commands:      []*suirpcv2.Command{initExecuteCmd, finishExecuteCmd},
			abortCommand:  0,
			abortFunction: strPtr("transmit"),
		},
		{
			name:          "abort with missing function name fails closed",
			commands:      []*suirpcv2.Command{initExecuteCmd, finishExecuteCmd},
			abortCommand:  1,
			abortFunction: nil,
		},
		{
			name:          "abort named init_execute is skipped",
			commands:      []*suirpcv2.Command{initExecuteCmd, finishExecuteCmd},
			abortCommand:  0,
			abortFunction: strPtr("init_execute"),
		},
		{
			name:          "abort after init_execute passes the gate",
			commands:      []*suirpcv2.Command{initExecuteCmd, finishExecuteCmd},
			inputs:        garbageInputs,
			abortCommand:  1,
			abortFunction: strPtr("finish_execute"),
			// Proceeding past the gate reaches report deserialization, which fails
			// on the garbage report bytes - proving the gate did not skip.
			wantErrContains: "failed to deserialize execution report",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tIndexer := &TransactionsIndexer{
				logger:                  logger.Test(t),
				executionEventModuleKey: "offramp",
				executeFunction:         "init_execute",
			}

			tx := gateTestFailedTx(tc.commands, tc.inputs, tc.abortCommand, tc.abortFunction)
			record, err := tIndexer.processFailedTransaction(
				t.Context(), tx, 0, "handle", offrampPkg, latestPkg, CheckpointMeta{})

			require.Nil(t, record)
			if tc.wantErrContains == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantErrContains)
			}
		})
	}
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
