package bind

import (
	"context"
	"fmt"

	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
)

const defaultGasBudget = 200000000

func BuildPTBFromArgs(ctx context.Context, client suiclient.ClientImpl, packageId *sui.Address, module string, function string, args ...any) (*suiptb.ProgrammableTransactionBuilder, error) {
	ptb := suiptb.NewTransactionDataTransactionBuilder()

	ptbArgs := make([]suiptb.Argument, 0, len(args))
	for _, arg := range args {
		var ptbArg suiptb.Argument
		var err error

		switch arg := arg.(type) {
		case string:
			if IsSuiAddress(arg) {
				// Fetch object information. Build argument based on it
				object, e := ReadObject(ctx, arg, client)
				if e != nil {
					return nil, e
				}
				ptbArg, err = ptb.Obj(toObjectArg(object))
			} else {
				ptbArg, err = ptb.Pure(arg)
			}
		// Rest of Sui primitive types work the same. Objects are the only difference
		default:
			ptbArg, err = ptb.Pure(arg)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to convert argument %v to Sui type: %w", arg, err)
		}
		ptbArgs = append(ptbArgs, ptbArg)
	}

	ptb.Command(suiptb.Command{
		MoveCall: &suiptb.ProgrammableMoveCall{
			Package:       packageId,
			Module:        module,
			Function:      function,
			TypeArguments: []sui.TypeTag{},
			Arguments:     ptbArgs,
		}},
	)

	return ptb, nil
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

func toObjectArg(object *suiclient.SuiObjectResponse) suiptb.ObjectArg {
	if object.Data.Owner.Shared != nil {
		return suiptb.ObjectArg{
			SharedObject: &suiptb.SharedObjectArg{
				Id:                   object.Data.ObjectId,
				Mutable:              true,
				InitialSharedVersion: *object.Data.Owner.Shared.InitialSharedVersion,
			},
		}
	}
	// TODO: Could there be a receiving object option?
	return suiptb.ObjectArg{
		ImmOrOwnedObject: &sui.ObjectRef{
			ObjectId: object.Data.ObjectId,
			Version:  object.Data.Version.Uint64(),
			Digest:   object.Data.Digest,
		},
	}
}
