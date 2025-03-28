package bind

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	sui_pattokan "github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
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
	signer signer.Signer,
	client sui.ISuiAPI,
	req PublishRequest,
) (PackageID, *models.SuiTransactionBlockResponse, error) {
	var modules [][]byte
	for _, encodedModule := range req.CompiledModules {
		decodedModule, err := codec.DecodeBase64(encodedModule)
		if err != nil {
			return "", nil, fmt.Errorf("failed to decode module: %v", err)
		}
		modules = append(modules, decodedModule)
	}

	var deps []*sui_pattokan.Address
	for _, dep := range req.Dependencies {
		suiAddressDep, err := ToSuiAddress(dep)
		if err != nil {
			return "", nil, fmt.Errorf("failed to convert dependency address: %v", err)
		}
		deps = append(deps, suiAddressDep)
	}

	signerAddress, err := ToSuiAddress(signer.Address)
	if err != nil {
		return "", nil, fmt.Errorf("failed to convert signer address: %v", err)
	}

	// Construct the Transaction using PTB
	ptb := suiptb.NewTransactionDataTransactionBuilder()
	arg := ptb.PublishUpgradeable(modules, deps)
	// The program object is transferred to the signer address once deployed
	recArg, err := ptb.Pure(signerAddress)
	ptb.Command(suiptb.Command{
		TransferObjects: &suiptb.ProgrammableTransferObjects{
			Objects: []suiptb.Argument{arg},
			Address: recArg,
		}})

	// Finish transaction details and encode it
	txBytes, err := FinishTransactionFromBuilder(ctx, ptb, opts, signer.Address, client)
	if err != nil {
		return "", nil, err
	}

	// Sign and send Transaction
	tx, err := SignAndSendTx(ctx, signer, client, txBytes)
	if err != nil {
		msg := fmt.Errorf("failed to execute tx when publishing: %v", err)
		return "", nil, msg
	}

	if tx.Effects.Status.Status == "failure" {
		return "", nil, fmt.Errorf("transaction failed: %v", tx.Effects.Status.Error)
	}

	// Find the object ID from the transaction
	pkgId, err := FindPackageIdFromPublishTx(*tx)
	if err != nil {
		return "", nil, err
	}

	return pkgId, tx, err
}

func FindPackageIdFromPublishTx(tx models.SuiTransactionBlockResponse) (string, error) {
	for _, change := range tx.ObjectChanges {
		if change.Type == "published" {
			return change.PackageId, nil
		}
	}

	return "", errors.New("package ID not found in transaction")
}

func FindObjectIdFromPublishTx(tx *models.SuiTransactionBlockResponse, module string) (string, error) {
	for _, change := range tx.ObjectChanges {
		if change.Type == "created" && strings.Contains(change.ObjectType, module) {
			return change.ObjectId, nil
		}
	}

	return "", errors.New("object ID not found in transaction")
}
