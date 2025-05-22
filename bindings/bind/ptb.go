package bind

import (
	"context"
	"fmt"
	"reflect"

	"golang.org/x/exp/slices"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
)

const defaultGasBudget = 200000000

/*
BuildPTBFromArgs creates a Programmable Transaction Builder (PTB) for a Sui Move call.
It converts Go values to Sui-compatible arguments and constructs a transaction that calls
the specified Move function.

Parameters:
  - ctx: The context for the operation
  - client: The Sui client used for fetching object information.
  - packageId: The address of the Move package to call
  - module: The name of the Move module
  - function: The name of the function to call
  - isObjectCreatedTransferred: Whether the function creates an object that should be transferred
  - recipient: The address that should receive any created objects (used only if isObjectCreatedTransferred is true)
  - args: The arguments to pass to the Move function
*/
func BuildPTBFromArgs(
	ctx context.Context,
	client suiclient.ClientImpl,
	packageId *sui.Address,
	module string,
	function string,
	isObjectCreatedTransferred bool,
	recipient string,
	args ...any) (*suiptb.ProgrammableTransactionBuilder, error) {
	ptb := suiptb.NewTransactionDataTransactionBuilder()

	ptbArgs := make([]suiptb.Argument, 0, len(args))
	for _, arg := range args {
		ptbArg, err := ToPTBArg(ctx, ptb, client, arg, true)
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument %v to Sui type: %w", arg, err)
		}
		ptbArgs = append(ptbArgs, ptbArg)
	}

	obj := ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       packageId,
			Module:        module,
			Function:      function,
			TypeArguments: []sui.TypeTag{},
			Arguments:     ptbArgs,
		}},
	)

	// Add the instruction to transfer the object if the function creates one
	if isObjectCreatedTransferred {
		recAddress, err := ToSuiAddress(recipient)
		if err != nil {
			return nil, fmt.Errorf("failed to convert signer address: %w", err)
		}
		recArg, err := ptb.Pure(recAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to encode recipient address: %w", err)
		}
		ptb.Command(suiptb.Command{
			TransferObjects: &suiptb.ProgrammableTransferObjects{
				Objects: []suiptb.Argument{obj},
				Address: recArg,
			}})
	}

	return ptb, nil
}

func ToPTBArg(
	ctx context.Context,
	ptb *suiptb.ProgrammableTransactionBuilder,
	client suiclient.ClientImpl,
	arg any,
	isMutable bool,
) (suiptb.Argument, error) {
	// check if the argument has already been included in the PTB args list
	// only if it's possible to BCS marshal the value (cases where the value is "pure" / not object)
	marshalledBytes, err := bcs.Marshal(arg)
	if err == nil {
		for idx, key := range ptb.Inputs.InsertOrderList {
			if key.Pure != nil && slices.Equal(*key.Pure, marshalledBytes) {
				//nolint:gosec
				foundIdx := uint16(idx)
				return suiptb.Argument{
					Input: &foundIdx,
				}, nil
			}
		}
	}

	switch v := arg.(type) {
	// ────────────────────── STRING ──────────────────────
	case string:
		if IsSuiAddress(v) {
			// attempt to treat it as an object first
			obj, err := ReadObject(ctx, v, client)
			if err == nil && obj.Error == nil && obj.Data != nil {
				return ptb.Obj(ToObjectArg(obj.Data, isMutable))
			}
			// otherwise treat as raw address
			addr, err := ToSuiAddress(v)
			if err != nil {
				return suiptb.Argument{}, fmt.Errorf("bad address %s: %w", v, err)
			}

			return ptb.Pure(addr)
		}

		return ptb.Pure(v)

	// ────────────────────── []byte  ─────────────────────
	case []byte:
		// empty vec<u8>  → 0-length prefix
		if len(v) == 0 {
			return ptb.Pure([]uint8{}) // sui-go encodes as single 0x00
		}
		// deep-copy into []uint8 (necessary because Pure keeps ref)
		out := make([]uint8, len(v))
		copy(out, v)

		return ptb.Pure(out)

	// ───────────────────── [][]byte  ────────────────────
	case [][]byte:
		if len(v) == 0 {
			return ptb.Pure([][]uint8{}) // vec<vec<u8>>{}
		}
		// deep-copy each element so duplicate contents don’t share pointer
		vv := make([][]uint8, len(v))
		for i, src := range v {
			dst := make([]uint8, len(src))
			copy(dst, src)
			vv[i] = dst
		}

		return ptb.Pure(vv)

	// ───────────────────── generic slice ────────────────
	default:
		// This could be recursive, but only to support the very rare case of <vector<vector<address>>
		rv := reflect.ValueOf(arg)
		if rv.Kind() == reflect.Slice {
			var vec []any
			for i := range rv.Len() {
				el := rv.Index(i).String()
				if IsSuiAddress(el) {
					address, err := ToSuiAddress(el)
					if err != nil {
						return suiptb.Argument{}, fmt.Errorf("failed to convert address %s to Sui type: %w", arg, err)
					}
					vec = append(vec, address)
				} else {
					el := rv.Index(i).Interface()
					vec = append(vec, el)
				}
			}

			return ptb.Pure(vec)
		}

		return ptb.Pure(arg)
	}
}

func FinishTransactionFromBuilder(ctx context.Context, ptb *suiptb.ProgrammableTransactionBuilder, opts TxOpts, signer string, client suiclient.ClientImpl) ([]byte, error) {
	txData, err := finishTransactionFromBuilder(ctx, ptb, opts, signer, client)
	if err != nil {
		return nil, err
	}

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction data: %w", err)
	}

	return txBytes, nil
}

func FinishDevInspectTransactionFromBuilder(ctx context.Context, ptb *suiptb.ProgrammableTransactionBuilder, opts TxOpts, signer string, client suiclient.ClientImpl) ([]byte, error) {
	txData, err := finishTransactionFromBuilder(ctx, ptb, opts, signer, client)
	if err != nil {
		return nil, err
	}

	txBytes, err := bcs.Marshal(txData.V1.Kind)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction kind data: %w", err)
	}

	return txBytes, nil
}

func finishTransactionFromBuilder(ctx context.Context, ptb *suiptb.ProgrammableTransactionBuilder, opts TxOpts, signer string, client suiclient.ClientImpl) (*suiptb.TransactionData, error) {
	pt := ptb.Finish()

	address, err := ToSuiAddress(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer address")
	}

	var coinData *sui.ObjectRef
	if opts.GasObject != "" {
		coinData, err = ToSuiObjectRef(ctx, client, opts.GasObject, signer)
	} else {
		coinData, err = FetchDefaultGasCoinRef(ctx, client, signer)
	}
	if err != nil {
		return nil, err
	}

	gasBudget := uint64(defaultGasBudget)
	if opts.GasBudget != nil {
		gasBudget = *opts.GasBudget
	}
	gasPrice := suiclient.DefaultGasPrice
	if opts.GasPrice != nil {
		gasPrice = *opts.GasPrice
	}
	txData := suiptb.NewTransactionData(
		address,
		pt,
		[]*sui.ObjectRef{coinData},
		gasBudget,
		gasPrice,
	)

	return &txData, nil
}

func ToObjectArg(object *suiclient.SuiObjectData, isMutable bool) suiptb.ObjectArg {
	if object != nil && object.Owner != nil && object.Owner.ObjectOwnerInternal != nil &&
		object.Owner.ObjectOwnerInternal.Shared != nil && object.Owner.ObjectOwnerInternal.Shared.InitialSharedVersion != nil {
		return suiptb.ObjectArg{
			SharedObject: &suiptb.SharedObjectArg{
				Id:                   object.ObjectId,
				Mutable:              isMutable,
				InitialSharedVersion: *object.Owner.Shared.InitialSharedVersion,
			},
		}
	}

	// TODO: Could there be a receiving object option?
	return suiptb.ObjectArg{
		ImmOrOwnedObject: &sui.ObjectRef{
			ObjectId: object.ObjectId,
			Version:  object.Version.Uint64(),
			Digest:   object.Digest,
		},
	}
}
