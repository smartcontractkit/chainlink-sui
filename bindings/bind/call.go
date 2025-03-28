package bind

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	sui_signer "github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	sui_p "github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// To allow PTB, each method returns a .Build() method (don't require client or signer) that would return the tx payload and .Execute() method (require signer and client) that would build and send the tx to the network
type IMethod interface {
	Build(opts TxOpts, signerAddress string) (*suiptb.ProgrammableTransactionBuilder, error)
	Execute(ctx context.Context, opts TxOpts, signer sui_signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error)
}

var _ IMethod = (*Method)(nil)

type BuildFunc func(opts TxOpts, signerAddress string) (*suiptb.ProgrammableTransactionBuilder, error)
type ExecuteFunc func(ctx context.Context, opts TxOpts, signer sui_signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error)

type Method struct {
	buildFunc BuildFunc
	execFunc  ExecuteFunc
}

func NewMethod(buildFunc BuildFunc, execFunc ExecuteFunc) *Method {
	return &Method{
		buildFunc: buildFunc,
		execFunc:  execFunc,
	}
}

func (m *Method) Build(opts TxOpts, signerAddress string) (*suiptb.ProgrammableTransactionBuilder, error) {
	return m.buildFunc(opts, signerAddress)
}

func (m *Method) Execute(ctx context.Context, opts TxOpts, signer sui_signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error) {
	return m.execFunc(ctx, opts, signer, client)
}

func MakeExecute(buildFn BuildFunc) ExecuteFunc {
	return func(ctx context.Context, opts TxOpts, signer sui_signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error) {
		ptb, err := buildFn(opts, string(signer.PubKey))
		if err != nil {
			return nil, err
		}

		txBytes, err := FinishTransactionFromBuilder(ctx, ptb, opts, signer.Address, client)
		if err != nil {
			return nil, err
		}

		return SignAndSendTx(ctx, signer, client, txBytes)
	}
}

func BuildCallTransaction(opts TxOpts, packageID, module, function string, args []any) (*suiptb.ProgrammableTransactionBuilder, error) {
	ptb := suiptb.NewTransactionDataTransactionBuilder()

	pkgObjectId, err := ToSuiAddress(packageID)
	if err != nil {
		return nil, err
	}

	callArgs := make([]suiptb.Argument, 0, len(args))
	for _, arg := range args {
		callArgs = append(callArgs, ptb.MustPure(arg))
	}

	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       pkgObjectId,
			Module:        module,
			Function:      function,
			TypeArguments: []sui_p.TypeTag{},
			Arguments:     callArgs,
		}},
	)

	return ptb, nil
}

// Encoding
func Encode(paramTypes []string, paramValues []any) (encodedArgs []any, err error) {
	args, err := serializeArgs(paramTypes, paramValues)
	if err != nil {
		return nil, err
	}

	return args, nil
}

// Improve this function to return typed values ([][]byte?)
func serializeArgs(paramTypes []string, paramValues []any) ([]any, error) {
	if len(paramValues) != len(paramTypes) {
		return nil, fmt.Errorf("paramTypes and paramValues must have the same length")
	}

	functionValues := make([]any, len(paramValues))
	for i, v := range paramValues {
		value, err := codec.EncodeToSuiValue(paramTypes[i], v)
		if err != nil {
			msg := fmt.Errorf("failed to encode value: %v (%w)", value, err)
			return nil, msg
		}

		functionValues[i] = value
	}

	return functionValues, nil
}
