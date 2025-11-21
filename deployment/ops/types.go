package operations

import (
	"fmt"

	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type OpTxInput[I any] struct {
	Input     I
	NoExecute bool
}

type OpTxResult[O any] struct {
	Digest    string
	PackageId string
	Objects   O
	Call      TransactionCall
}

type TransactionCall struct {
	PackageID  string
	Module     string
	Function   string
	Data       []byte
	StateObjID string
	TypeArgs   []string
}

type OpTxDeps struct {
	Client sui.ISuiAPI
	Signer rel.SuiSigner
	// We could have some logic to modify the gas based on input
	GetCallOpts func() *bind.CallOpts
	SuiRPC      string
}

func NewSuiOperationName(pkg string, module string, action string) string {
	return fmt.Sprintf("sui-%s-%s-%s", pkg, module, action)
}

func ToTransactionCall(call *bind.EncodedCall, stateObjID string) (TransactionCall, error) {
	if call == nil {
		return TransactionCall{}, fmt.Errorf("nil call provided")
	}

	calldata, err := extractByteArgsFromEncodedCall(*call)
	if err != nil {
		return TransactionCall{}, fmt.Errorf("failed to extract byte args from encoded call: %w", err)
	}

	return TransactionCall{
		PackageID:  call.Module.PackageID,
		Module:     call.Module.ModuleName,
		Function:   call.Function,
		Data:       calldata,
		TypeArgs:   []string{},
		StateObjID: stateObjID,
	}, nil
}

func ToTransactionCallWithTypeArgs(call *bind.EncodedCall, stateObjID string, typeArgs []string) (TransactionCall, error) {
	if call == nil {
		return TransactionCall{}, fmt.Errorf("nil call provided")
	}

	calldata, err := extractByteArgsFromEncodedCall(*call)
	if err != nil {
		return TransactionCall{}, fmt.Errorf("failed to extract byte args from encoded call: %w", err)
	}

	return TransactionCall{
		PackageID:  call.Module.PackageID,
		Module:     call.Module.ModuleName,
		Function:   call.Function,
		Data:       calldata,
		TypeArgs:   typeArgs,
		StateObjID: stateObjID,
	}, nil
}

func extractByteArgsFromEncodedCall(encodedCall bind.EncodedCall) ([]byte, error) {
	var args []byte
	for _, callArg := range encodedCall.CallArgs {
		if callArg.CallArg.UnresolvedObject != nil {
			args = append(args, callArg.CallArg.UnresolvedObject.ObjectId[:]...)
		}
		if callArg.CallArg.Pure != nil {
			args = append(args, callArg.CallArg.Pure.Bytes...)
		}
		// we don't support resolved objects
		if callArg.CallArg.Object != nil {
			return nil, fmt.Errorf("resolved objects are not supported in transaction call encoding")
		}
	}

	return args, nil
}
