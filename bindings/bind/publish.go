package bind

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type PackageID = string

type PublishRequest struct {
	CompiledModules []string `json:"compiled_modules"`
	Dependencies    []string `json:"dependencies"`
}

func PublishPackage(
	ctx context.Context,
	opts TxOpts,
	// TODO: Replace by a Signer common interface
	signer rel.SuiSigner,
	client suiclient.ClientImpl,
	req PublishRequest,
) (PackageID, *suiclient.SuiTransactionBlockResponse, error) {
	var modules = make([][]byte, 0, len(req.CompiledModules))
	for _, encodedModule := range req.CompiledModules {
		decodedModule, err := codec.DecodeBase64(encodedModule)
		if err != nil {
			return "", nil, fmt.Errorf("failed to decode module: %w", err)
		}
		modules = append(modules, decodedModule)
	}

	deps := make([]*sui.Address, 0, len(req.Dependencies))
	for _, dep := range req.Dependencies {
		suiAddressDep, err := ToSuiAddress(dep)
		if err != nil {
			return "", nil, fmt.Errorf("failed to convert dependency address: %w", err)
		}
		deps = append(deps, suiAddressDep)
	}

	_signerAddress, err := signer.GetAddress()
	if err != nil {
		return "", nil, err
	}
	signerAddress, err := ToSuiAddress(_signerAddress)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert signer address: %w", err)
	}

	// Construct the Transaction using PTB
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	arg := ptb.PublishUpgradeable(modules, deps)
	// The program object is transferred to the signer address once deployed
	recArg, err := ptb.Pure(signerAddress)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create tx argument: %w", err)
	}
	ptb.Command(suiptb.Command{
		TransferObjects: &suiptb.ProgrammableTransferObjects{
			Objects: []suiptb.Argument{arg},
			Address: recArg,
		}})

	// Finish transaction details and encode it
	txBytes, err := FinishTransactionFromBuilder(ctx, ptb, opts, _signerAddress, client)
	if err != nil {
		return "", nil, err
	}

	// Sign and send Transaction
	tx, err := SignAndSendTx(ctx, signer, client, txBytes)
	if err != nil {
		msg := fmt.Errorf("failed to execute tx when publishing: %w", err)
		return "", nil, msg
	}

	if tx.Effects.Data.V1.Status.Status == FailureResultType {
		return "", nil, fmt.Errorf("transaction failed: %v", tx.Effects.Data.V1.Status.Status)
	}

	// Find the object ID from the transaction
	pkgId, err := FindPackageIdFromPublishTx(*tx)
	if err != nil {
		return "", nil, err
	}

	return pkgId, tx, err
}

func FindPackageIdFromPublishTx(tx suiclient.SuiTransactionBlockResponse) (string, error) {
	for _, change := range tx.ObjectChanges {
		if change.Data.Published != nil {
			return change.Data.Published.PackageId.String(), nil
		}
	}

	return "", errors.New("package ID not found in transaction")
}

func FindObjectIdFromPublishTx(tx *suiclient.SuiTransactionBlockResponse, module, object string) (string, error) {
	for _, change := range tx.ObjectChanges {
		if change.Data.Created != nil && strings.Contains(change.Data.Created.ObjectType, fmt.Sprintf("%v::%v", module, object)) {
			return change.Data.Created.ObjectId.String(), nil
		}
	}

	return "", errors.New("object ID not found in transaction")
}
