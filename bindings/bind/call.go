package bind

import (
	"context"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

type TxOpts struct {
	GasObject string
	GasBudget string
}

func SignAndSendTx(ctx context.Context, signer signer.Signer, client sui.ISuiAPI, unsignedTx models.TxnMetaData) (*models.SuiTransactionBlockResponse, error) {
	signedTx := unsignedTx.SignSerializedSigWith(signer.PriKey)
	blockReq := &models.SuiExecuteTransactionBlockRequest{
		TxBytes:   signedTx.TxBytes,
		Signature: []string{signedTx.Signature},
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		// RequestType:
	}

	tx, err := client.SuiExecuteTransactionBlock(ctx, *blockReq)
	if err != nil {
		// TODO: include more details about the function and arguments
		msg := fmt.Sprintf("tx failed calling move method: %v", err)
		return nil, errors.New(msg)
	}

	return &tx, nil
}

func BuildCallRequest(opts TxOpts, signerAddress, packageObjectId, module, function string, args []any) models.MoveCallRequest {
	return models.MoveCallRequest{
		Signer:          signerAddress,
		PackageObjectId: packageObjectId,
		Module:          module,
		Function:        function,
		TypeArguments:   []any{},
		Arguments:       args,
		Gas:             &opts.GasObject,
		GasBudget:       opts.GasBudget,
		ExecutionMode:   "WaitForCommit",
	}
}

// To allow PTB, each method returns a .BuildTx() method (don't require client or signer) that would return the tx payload and .Execute() method (require signer and client) that would build and send the tx to the network
type IMethod interface {
	Build(opts TxOpts, signerAddress string) (models.TxnMetaData, error)
	Execute(ctx context.Context, opts TxOpts, signer signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error)
}

var _ IMethod = (*Method)(nil)

type BuildFunc func(opts TxOpts, signerAddress string) (models.TxnMetaData, error)
type ExecuteFunc func(ctx context.Context, opts TxOpts, signer signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error)

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

func (m *Method) Build(opts TxOpts, signerAddress string) (models.TxnMetaData, error) {
	return m.buildFunc(opts, signerAddress)
}

func (m *Method) Execute(ctx context.Context, opts TxOpts, signer signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error) {
	return m.execFunc(ctx, opts, signer, client)
}

func MakeExecute(buildFn BuildFunc) ExecuteFunc {
	return func(ctx context.Context, opts TxOpts, signer signer.Signer, client sui.ISuiAPI) (*models.SuiTransactionBlockResponse, error) {
		unsignedTx, err := buildFn(opts, string(signer.PubKey))
		if err != nil {
			return nil, err
		}
		return SignAndSendTx(ctx, signer, client, unsignedTx)
	}
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
		// TODO: codec should live in a common place, we could have circular dependencies
		value, err := codec.EncodeToSuiValue(paramTypes[i], v)
		if err != nil {
			msg := fmt.Errorf("failed to encode value: %v (%v)", value, err)
			return nil, msg
		}

		functionValues[i] = value
	}

	return functionValues, nil
}
