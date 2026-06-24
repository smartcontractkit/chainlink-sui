package client

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/blake2b"

	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/transaction"
)

func (c *PTBClient) transformMoveCallArgFromSignature(
	ctx context.Context,
	tx *transaction.Transaction,
	arg any,
	body *suirpcv2.OpenSignatureBody,
	mutable bool,
) (*transaction.Argument, error) {
	if body == nil {
		return nil, errors.New("missing parameter type signature")
	}

	c.log.Debugw("transformMoveCallArgFromSignature", "type", body.GetType(), "arg", arg)

	switch body.GetType() {
	case suirpcv2.OpenSignatureBody_VECTOR:
		instantiations := body.GetTypeParameterInstantiation()
		if len(instantiations) != 1 {
			return nil, fmt.Errorf("unsupported vector type arity: %d", len(instantiations))
		}

		inner := instantiations[0]
		switch inner.GetType() {
		case suirpcv2.OpenSignatureBody_U8:
			pureArg := tx.Pure(arg)
			return &pureArg, nil
		case suirpcv2.OpenSignatureBody_ADDRESS:
			return c.transformVectorAddressArg(tx, arg)
		case suirpcv2.OpenSignatureBody_VECTOR:
			nested := inner.GetTypeParameterInstantiation()
			if len(nested) == 1 && nested[0].GetType() == suirpcv2.OpenSignatureBody_U8 {
				pureArg := tx.Pure(arg)
				return &pureArg, nil
			}
			return nil, errors.New("unsupported nested vector type")
		default:
			return nil, fmt.Errorf("unsupported vector element type: %v", inner.GetType())
		}
	case suirpcv2.OpenSignatureBody_U8,
		suirpcv2.OpenSignatureBody_U16,
		suirpcv2.OpenSignatureBody_U32,
		suirpcv2.OpenSignatureBody_U64,
		suirpcv2.OpenSignatureBody_U128,
		suirpcv2.OpenSignatureBody_U256,
		suirpcv2.OpenSignatureBody_BOOL,
		suirpcv2.OpenSignatureBody_ADDRESS:
		pureArg := tx.Pure(arg)
		return &pureArg, nil
	default:
		argTypeString := suirpcv2.OpenSignatureBody_Type_name[int32(body.GetType())]
		return c.TransformTransactionArg(ctx, tx, arg, argTypeString, mutable)
	}
}

func (c *PTBClient) transformVectorAddressArg(tx *transaction.Transaction, arg any) (*transaction.Argument, error) {
	switch v := arg.(type) {
	case []string:
		return convertAddresses(tx, v)
	case []interface{}:
		addresses := make([]models.SuiAddressBytes, len(v))
		for i, raw := range v {
			s, ok := raw.(string)
			if !ok {
				return nil, fmt.Errorf("vector<address> element is not string: %T", raw)
			}

			s = strings.TrimPrefix(s, "0x")
			b, err := hex.DecodeString(s)
			if err != nil {
				b, err = base64.StdEncoding.DecodeString(s)
				if err != nil {
					return nil, fmt.Errorf("failed to decode address %q: %w", s, err)
				}
			}

			if len(b) != 32 {
				return nil, fmt.Errorf("address at index %d has wrong length %d, want 32", i, len(b))
			}

			var addr models.SuiAddressBytes
			copy(addr[:], b)
			addresses[i] = addr
		}

		pureArg := tx.Pure(addresses)
		return &pureArg, nil
	default:
		return nil, fmt.Errorf("expected []string or []interface{} for vector<address>, got %T", arg)
	}
}

func (c *PTBClient) TransformTransactionArg(
	ctx context.Context,
	tx *transaction.Transaction,
	arg any,
	argType string,
	mutable bool,
) (*transaction.Argument, error) {
	c.log.Debugw("TransformTransactionArg", "argType", argType, "arg", arg)

	switch argType {
	case "objectId", "object_id", "DATATYPE":
		objectIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(arg.(string)))
		if err != nil {
			return nil, err
		}
		// get object's reference metadata (owner/version/digest only) - intentionally avoids fetching the
		// full (and growing) object contents/json, which is unnecessary for building an ObjectArg and was
		// the dominant source of read latency under load.
		objectDetails, err := c.readObjectMetadataInternal(ctx, arg.(string))
		if err != nil || objectDetails == nil {
			return nil, fmt.Errorf("failed to read object %s: %w", arg.(string), err)
		}
		var objectArg transaction.ObjectArg

		// handle truly immutable objects
		if objectDetails.Owner != nil && objectDetails.Owner.GetKind() == suirpcv2.Owner_IMMUTABLE {
			var versionUint uint64
			var digestBytes *models.ObjectDigestBytes
			versionUint = objectDetails.GetVersion()
			if err != nil {
				return nil, fmt.Errorf("failed to parse version: %w", err)
			}
			digestBytes, err = transaction.ConvertObjectDigestStringToBytes(
				models.ObjectDigest(objectDetails.GetDigest()),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to convert object digest: %w", err)
			}

			objectArg = transaction.ObjectArg{
				ImmOrOwnedObject: &transaction.SuiObjectRef{
					ObjectId: *objectIdBytes,
					Version:  versionUint,
					Digest:   *digestBytes,
				},
			}
			callArg := tx.Object(transaction.CallArg{Object: &objectArg})

			return &callArg, nil
		}

		// construct the objectArg
		switch objectDetails.Owner.GetKind() {
		case suirpcv2.Owner_SHARED:
			if objectDetails.Owner.GetVersion() == 0 {
				return nil, fmt.Errorf("shared object %s missing initial shared version", arg.(string))
			}
			objectArg = transaction.ObjectArg{
				SharedObject: &transaction.SharedObjectRef{
					ObjectId:             *objectIdBytes,
					InitialSharedVersion: objectDetails.Owner.GetVersion(),
					Mutable:              mutable,
				},
			}
		case suirpcv2.Owner_ADDRESS:
			if objectDetails.Owner.GetAddress() == "" {
				return nil, fmt.Errorf("address-owned object %s missing owner address", arg.(string))
			}
			digestBytes, err := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest(objectDetails.GetDigest()))
			if err != nil {
				return nil, fmt.Errorf("failed to convert object digest: %w", err)
			}
			objectArg = transaction.ObjectArg{
				ImmOrOwnedObject: &transaction.SuiObjectRef{
					ObjectId: *objectIdBytes,
					Version:  objectDetails.GetVersion(),
					Digest:   *digestBytes,
				},
			}
		default:
			return nil, fmt.Errorf("unknown object owner: %v", objectDetails.Owner.GetAddress())
		}

		// construct the arg
		transactionObjectArg := tx.Object(
			transaction.CallArg{
				Object: &objectArg,
			},
		)

		return &transactionObjectArg, nil
	case "string":
		// hex encode the string
		if str, ok := arg.(string); ok {
			hexStr := hex.EncodeToString([]byte(str))
			pureArg := tx.Pure(hexStr)

			return &pureArg, nil
		}
	case "vector<address>", "VECTOR":
		return c.transformVectorAddressArg(tx, arg)

	default:
		pureArg := tx.Pure(arg)
		return &pureArg, nil
	}

	return nil, fmt.Errorf("unknown argument type: %s", argType)
}

func (c *PTBClient) GetTransactionPaymentCoinForAddress(ctx context.Context, payer string) (models.SuiAddressBytes, uint64, models.ObjectDigestBytes, error) {
	coins, err := c.GetCoinsByAddress(ctx, payer)
	if err != nil {
		return models.SuiAddressBytes{}, 0, nil, err
	}
	if len(coins) == 0 {
		return models.SuiAddressBytes{}, 0, nil, fmt.Errorf("no coins available for gas payment")
	}

	//nolint:revive // var-naming: matches Sui object ID field naming
	coinObjectIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(coins[0].GetObjectId()))
	if err != nil {
		return models.SuiAddressBytes{}, 0, nil, err
	}
	versionUint := coins[0].GetVersion()

	digestBytes, err := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest(coins[0].GetDigest()))
	if err != nil {
		return models.SuiAddressBytes{}, 0, nil, fmt.Errorf("failed to convert object digest: %w", err)
	}

	return *coinObjectIdBytes, versionUint, *digestBytes, nil
}

// HashTxBytes is a helper method to hash (Blake2) the transaction bytes before signing
func (c *PTBClient) HashTxBytes(txBytes []byte) []byte {
	intentMessage := append([]byte{0, 0, 0}, txBytes...)
	digest := blake2b.Sum256(intentMessage)
	return digest[:]
}

// convertAddresses converts string address to sui address bytes
func convertAddresses(tx *transaction.Transaction, addresses []string) (*transaction.Argument, error) {
	converted := make([]models.SuiAddressBytes, len(addresses))
	for i, addr := range addresses {
		addressBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(addr))
		if err != nil {
			return nil, fmt.Errorf("failed to convert address %s to Sui address: %w", addr, err)
		}
		converted[i] = *addressBytes
	}
	pureArg := tx.Pure(converted)
	return &pureArg, nil
}

// Add helper method to create type tags
func (c *PTBClient) CreateTypeTag(typeStr string) (transaction.TypeTag, error) {
	if typeStr == "" {
		return transaction.TypeTag{}, errors.New("type string cannot be empty")
	}

	// Handle struct types (package::module::name)
	if strings.Contains(typeStr, "::") {
		parts := strings.Split(typeStr, "::")
		if len(parts) != 3 {
			return transaction.TypeTag{}, fmt.Errorf("invalid struct type format %q, expected package::module::name", typeStr)
		}

		packageID, module, name := parts[0], parts[1], parts[2]

		// Convert package ID to address bytes
		packageAddr := models.SuiAddress(packageID)
		addressBytes, err := transaction.ConvertSuiAddressStringToBytes(packageAddr)
		if err != nil {
			return transaction.TypeTag{}, fmt.Errorf("failed to convert package address %q: %w", packageID, err)
		}

		return transaction.TypeTag{
			Struct: &transaction.StructTag{
				Address:    *addressBytes,
				Module:     module,
				Name:       name,
				TypeParams: []*transaction.TypeTag{},
			},
		}, nil
	}

	// TODO: Handle primitive types if needed
	return transaction.TypeTag{}, fmt.Errorf("unsupported type format: %s", typeStr)
}
