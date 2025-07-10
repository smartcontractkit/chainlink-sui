package bind

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

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
	var modules = make([][]byte, 0, len(req.CompiledModules))
	for _, encodedModule := range req.CompiledModules {
		decodedModule, err := bindutils.DecodeBase64(encodedModule)
		if err != nil {
			return "", nil, fmt.Errorf("failed to decode module: %w", err)
		}
		modules = append(modules, decodedModule)
	}

	// dependencies are already strings
	deps := req.Dependencies

	signerAddressStr, err := opts.Signer.GetAddress()
	if err != nil {
		return "", nil, err
	}
	signerAddress, err := bindutils.ConvertAddressToString(signerAddressStr)
	if err != nil {
		return "", nil, fmt.Errorf("invalid signer address %v: %w", signerAddressStr, err)
	}

	// TODO: we are using unsafe_publish here because module fields are []models.SuiAddressBytes?
	// https://github.com/block-vision/sui-go-sdk/blob/b382e3a4ec2e9233461cdecbd45e6c031166234a/transaction/transaction_data.go#L232
	// https://github.com/MystenLabs/sui/blob/6d8ceed4d727a6c9ce9f5879b3cd1b2e8605affa/crates/sui-graphql-rpc/src/types/transaction_block_kind/programmable.rs#L134

	// Convert modules to base64
	moduleStrs := make([]string, len(modules))
	for i, module := range modules {
		moduleB64 := base64.StdEncoding.EncodeToString(module)
		moduleStrs[i] = moduleB64
	}

	// Set gas budget
	gasBudget := "200000000" // 200M MIST default for publish
	if opts.GasBudget != nil {
		gasBudget = fmt.Sprintf("%d", *opts.GasBudget)
	}

	// Convert gas object to pointer if provided
	var gasObj *string
	if opts.GasObject != "" {
		gasObj = &opts.GasObject
	}

	// Use the client's Publish method instead of PTB
	publishResp, err := client.Publish(ctx, models.PublishRequest{
		Sender:          signerAddress,
		CompiledModules: moduleStrs,
		Dependencies:    deps,
		GasBudget:       gasBudget,
		Gas:             gasObj,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to create publish transaction: %w", err)
	}

	// the SDK returns base64 encoded bytes, we need to decode them
	txBytesDecoded, err := base64.StdEncoding.DecodeString(publishResp.TxBytes)
	if err != nil {
		return "", nil, fmt.Errorf("failed to decode transaction bytes: %w", err)
	}

	tx, err := SignAndSendTx(ctx, opts.Signer, client, txBytesDecoded, opts.WaitForExecution)
	if err != nil {
		msg := fmt.Errorf("failed to execute tx when publishing: %w", err)
		return "", nil, msg
	}

	if tx.Effects.Status.Status == "failure" {
		return "", nil, fmt.Errorf("transaction failed: %v", tx.Effects.Status.Error)
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
