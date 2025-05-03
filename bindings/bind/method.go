package bind

import (
	"context"
	"fmt"

	"github.com/pattonkan/sui-go/sui/suiptb"

	"github.com/pattonkan/sui-go/suiclient"

	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// To allow PTB, each method returns a .Build() method (don't require client or signer) that would return the tx payload and .Execute() method (require signer and client) that would build and send the tx to the network
type IMethod interface {
	// Build could use some client calls to fetch for object details
	Build(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error)
	// Executes the signed transaction. If successful, it will cause state changes in the chain
	Execute(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.SuiTransactionBlockResponse, error)
	// Return transaction execution effects including the gas cost summary, while the effects are not committed to the chain.
	Inspect(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.DevInspectTransactionBlockResponse, error)
}

var _ IMethod = (*Method)(nil)

type BuildFunc func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error)
type ExecuteFunc func(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.SuiTransactionBlockResponse, error)
type InspectFunc func(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.DevInspectTransactionBlockResponse, error)

type Method struct {
	buildFunc   BuildFunc
	execFunc    ExecuteFunc
	inspectFunc InspectFunc
}

func NewMethod(buildFunc BuildFunc, execFunc ExecuteFunc, inspectFunc InspectFunc) *Method {
	return &Method{
		buildFunc:   buildFunc,
		execFunc:    execFunc,
		inspectFunc: inspectFunc,
	}
}

func (m *Method) Build(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
	return m.buildFunc(ctx)
}

func (m *Method) Execute(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.SuiTransactionBlockResponse, error) {
	return m.execFunc(ctx, opts, signer, client)
}

func (m *Method) Inspect(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.DevInspectTransactionBlockResponse, error) {
	return m.inspectFunc(ctx, opts, signer, client)
}

func MakeExecute(buildFn BuildFunc) ExecuteFunc {
	return func(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.SuiTransactionBlockResponse, error) {
		ptb, err := buildFn(ctx)
		if err != nil {
			return nil, err
		}
		address, err := signer.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get address: %w", err)
		}

		txBytes, err := FinishTransactionFromBuilder(ctx, ptb, opts, address, client)
		if err != nil {
			return nil, err
		}

		receipt, err := SignAndSendTx(ctx, signer, client, txBytes)
		if err != nil {
			return nil, err
		}

		if receipt.Effects.Data.V1.Status.Status == FailureResultType {
			return nil, fmt.Errorf("transaction failed: %v", receipt.Effects.Data.V1.Status.Error)
		}

		return receipt, nil
	}
}

func MakeInspect(buildFn BuildFunc) InspectFunc {
	return func(ctx context.Context, opts TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*suiclient.DevInspectTransactionBlockResponse, error) {
		ptb, err := buildFn(ctx)
		if err != nil {
			return nil, err
		}
		address, err := signer.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get address: %w", err)
		}

		txBytes, err := FinishDevInspectTransactionFromBuilder(ctx, ptb, opts, address, client)
		if err != nil {
			return nil, err
		}

		receipt, err := DevInspectTx(ctx, address, client, txBytes)
		if err != nil {
			return nil, err
		}

		if receipt.Effects.Data.V1.Status.Status == FailureResultType {
			return nil, fmt.Errorf("transaction inspect failed: %v", receipt.Effects.Data.V1.Status.Error)
		}

		return receipt, nil
	}
}
