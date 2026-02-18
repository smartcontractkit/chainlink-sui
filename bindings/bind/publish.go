package bind

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

type PackageID = string

type PublishRequest struct {
	CompiledModules []string `json:"compiled_modules"`
	Dependencies    []string `json:"dependencies"`
}

func PublishPackage(
	ctx context.Context,
	opts *CallOpts,
	client sui.ISuiAPI,
	req PublishRequest,
) (PackageID, *models.SuiTransactionBlockResponse, error) {
	if opts == nil {
		return "", nil, errors.New("opts cannot be nil")
	}
	if opts.Signer == nil {
		return "", nil, errors.New("opts.Signer cannot be nil")
	}

	var modules = make([][]byte, 0, len(req.CompiledModules))
	for _, encodedModule := range req.CompiledModules {
		decodedModule, err := bindutils.DecodeBase64(encodedModule)
		if err != nil {
			return "", nil, fmt.Errorf("failed to decode module: %w", err)
		}
		modules = append(modules, decodedModule)
	}

	deps := make([]models.SuiAddress, len(req.Dependencies))
	for i, dep := range req.Dependencies {
		addr := models.SuiAddress(dep)
		deps[i] = addr
	}

	signerAddressStr, err := opts.Signer.GetAddress()
	if err != nil {
		return "", nil, err
	}
	signerAddress, err := bindutils.ConvertAddressToString(signerAddressStr)
	if err != nil {
		return "", nil, fmt.Errorf("invalid signer address %v: %w", signerAddressStr, err)
	}

	gasBudgetValueDefault := uint64(1_000_000_000)
	if opts.GasBudget == nil {
		opts.GasBudget = &gasBudgetValueDefault // 500M MIST default for publish
	}

	ptb := transaction.NewTransaction()
	arg := ptb.Publish(modules, deps)
	// The program object is transferred to the signer once deployed
	recArg := ptb.Pure(signerAddress)
	ptb.TransferObjects([]transaction.Argument{arg}, recArg)

	tx, err := ExecutePTB(ctx, opts, client, ptb)
	if err != nil {
		return "", nil, fmt.Errorf("failed to execute publish transaction: %w", err)
	}

	pkgId, err := FindPackageIdFromPublishTx(*tx)
	if err != nil {
		return "", nil, err
	}

	return pkgId, tx, nil
}

func FindPackageIdFromPublishTx(tx models.SuiTransactionBlockResponse) (string, error) {
	if len(tx.ObjectChanges) == 0 {
		return "", errors.New("no object changes in transaction")
	}

	for _, change := range tx.ObjectChanges {
		if change.Type == "published" && change.PackageId != "" {
			return change.PackageId, nil
		}
	}

	return "", errors.New("package ID not found in transaction")
}

func FindObjectIdFromPublishTx(tx models.SuiTransactionBlockResponse, module, object string) (string, error) {
	if tx.ObjectChanges == nil {
		return "", errors.New("no object changes in transaction")
	}

	for _, change := range tx.ObjectChanges {
		if change.Type == "created" && change.ObjectType != "" {
			objectType := change.ObjectType

			// first, strip the generics since it'll contain '::' substrings
			if genericStart := strings.Index(objectType, "<"); genericStart != -1 {
				objectType = objectType[:genericStart]
			}

			parts := strings.Split(objectType, "::")

			const minPartsCount = 3
			if len(parts) >= minPartsCount {
				lastPart := parts[len(parts)-1]

				objectName := lastPart

				// Build module name from middle parts, incase there's more than 1
				// TODO: is this possible eg for a nested object?
				middleParts := parts[1 : len(parts)-1]
				moduleName := strings.Join(middleParts, "::")

				if objectName == object && moduleName == module {
					return change.ObjectId, nil
				}
			}
		}
	}

	return "", fmt.Errorf("object ID (module %s, object %s) not found in transaction", module, object)
}

// FindCoinObjectIdFromTx finds a coin object ID from a transaction response by looking for created objects of type Coin<T>
func FindCoinObjectIdFromTx(tx models.SuiTransactionBlockResponse, coinType string) (string, error) {
	expectedType := fmt.Sprintf("0x2::coin::Coin<%s>", coinType)

	for _, change := range tx.ObjectChanges {
		if change.Type == "created" && change.ObjectType == expectedType {
			return change.ObjectId, nil
		}
	}

	return "", fmt.Errorf("coin object of type %s not found in transaction", coinType)
}
