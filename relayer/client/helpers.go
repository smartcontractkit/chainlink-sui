package client

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/blake2b"

	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/transaction"
)

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
		// get object's details
		objectDetails, err := c.ReadObjectId(ctx, arg.(string))
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
		if objectDetails.Owner.GetKind() == suirpcv2.Owner_SHARED && objectDetails.Owner.GetVersion() != 0 {
			objectArg = transaction.ObjectArg{
				SharedObject: &transaction.SharedObjectRef{
					ObjectId:             *objectIdBytes,
					InitialSharedVersion: objectDetails.Owner.GetVersion(),
					Mutable:              mutable,
				},
			}
		} else if objectDetails.Owner.GetKind() == suirpcv2.Owner_ADDRESS && objectDetails.Owner.GetAddress() != "" {
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
		} else {
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
		switch v := arg.(type) {
		case []string:
			// Already []string, convert directly
			return convertAddresses(tx, v)

		case []interface{}:
			// JSON-decoded slice -> could be []byte or hex or base64 strings
			addresses := make([]models.SuiAddressBytes, len(v))
			for i, raw := range v {
				s, ok := raw.(string)
				if !ok {
					return nil, fmt.Errorf("vector<address> element is not string: %T", raw)
				}

				s = strings.TrimPrefix(s, "0x")
				b, err := hex.DecodeString(s)
				if err != nil {
					// fallback: base64 decode (Go JSON encodes []byte → base64 string)
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

func (c *PTBClient) convertBlockvisionResponse(resp *models.SuiTransactionBlockResponse) SuiTransactionBlockResponse {
	result := SuiTransactionBlockResponse{
		TxDigest: resp.Digest,
		Timestamp: func() uint64 {
			if resp.TimestampMs == "" {
				return 0
			}
			ts, err := strconv.ParseUint(resp.TimestampMs, 10, 64)
			if err != nil {
				c.log.Errorw("failed to parse timestamp", "error", err, "timestamp", resp.TimestampMs)
				return 0
			}

			return ts
		}(),
		Height: func() uint64 {
			if resp.Checkpoint == "" {
				return 0
			}
			h, err := strconv.ParseUint(resp.Checkpoint, 10, 64)
			if err != nil {
				c.log.Errorw("failed to parse height", "error", err, "height", resp.Checkpoint)
				return 0
			}

			return h
		}(),
		Status: SuiExecutionStatus{
			Status: resp.Effects.Status.Status,
			Error:  resp.Effects.Status.Error,
		},
		ObjectChanges: resp.ObjectChanges,
		Events:        resp.Events,
		Effects:       resp.Effects,
	}

	// Note: Full conversion of effects, events, and object changes would require
	// detailed mapping between blockvision and internal models
	// For now, keeping the basic structure

	return result
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
		return transaction.TypeTag{}, fmt.Errorf("type string cannot be empty")
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
