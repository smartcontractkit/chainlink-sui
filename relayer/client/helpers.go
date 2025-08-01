package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
)

func (c *PTBClient) TransformTransactionArg(
	ctx context.Context,
	tx *transaction.Transaction,
	arg any,
	argType string,
	mutable bool,
) (*transaction.Argument, error) {
	switch argType {
	case "objectId", "object_id":
		objectIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(arg.(string)))
		if err != nil {
			return nil, err
		}
		// get object's details
		objectDetails, err := c.ReadObjectId(ctx, arg.(string))
		if err != nil {
			return nil, err
		}
		var objectOwner models.ObjectOwner
		var objectArg transaction.ObjectArg

		// handle truly immutable objects
		if ownerStr, ok := objectDetails.Owner.(string); ok && ownerStr == "Immutable" {
			var versionUint uint64
			var digestBytes *models.ObjectDigestBytes
			versionUint, err = strconv.ParseUint(objectDetails.Version, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version: %w", err)
			}
			digestBytes, err = transaction.ConvertObjectDigestStringToBytes(
				models.ObjectDigest(objectDetails.Digest),
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

		// convert the response map into ObjectOwner type
		ownerData, ok := objectDetails.Owner.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to convert owner to map")
		}
		ownerJSON, err := json.Marshal(ownerData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal owner data: %w", err)
		}
		err = json.Unmarshal(ownerJSON, &objectOwner)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal owner data: %w", err)
		}
		// construct the objectArg
		if objectOwner.Shared.InitialSharedVersion != 0 {
			objectArg = transaction.ObjectArg{
				SharedObject: &transaction.SharedObjectRef{
					ObjectId:             *objectIdBytes,
					InitialSharedVersion: objectOwner.Shared.InitialSharedVersion,
					Mutable:              mutable,
				},
			}
		} else if objectOwner.AddressOwner != "" {
			versionUint, err := strconv.ParseUint(objectDetails.Version, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse version: %w", err)
			}
			digestBytes, err := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest(objectDetails.Digest))
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
		} else {
			return nil, fmt.Errorf("unknown object owner: %v", objectOwner)
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
	case "vector<address>":
		// Handle vector of addresses
		if addresses, ok := arg.([]string); ok {
			// Convert each address string to proper Sui address bytes
			convertedAddresses := make([]models.SuiAddressBytes, len(addresses))
			for i, addr := range addresses {
				addressBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(addr))
				if err != nil {
					return nil, fmt.Errorf("failed to convert address %s to Sui address: %w", addr, err)
				}
				convertedAddresses[i] = *addressBytes
			}
			pureArg := tx.Pure(convertedAddresses)

			return &pureArg, nil
		}

		return nil, fmt.Errorf("expected []string for vector<address>, got %T", arg)
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

	coinObjectIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(coins[0].CoinObjectId))
	if err != nil {
		return models.SuiAddressBytes{}, 0, nil, err
	}
	versionUint, err := strconv.ParseUint(coins[0].Version, 10, 64)
	if err != nil {
		return models.SuiAddressBytes{}, 0, nil, fmt.Errorf("failed to parse version: %w", err)
	}
	digestBytes, err := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest(coins[0].Digest))
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

func (c *PTBClient) PayAllSui(ctx context.Context, toAddress string, coinObjectRefs []string, signer string) error {
	_, err := c.client.PayAllSui(ctx, models.PayAllSuiRequest{
		Recipient:   toAddress,
		GasBudget:   "10000000",
		Signer:      signer,
		SuiObjectId: coinObjectRefs,
	})
	if err != nil {
		return fmt.Errorf("failed to add PayAllSui command: %w", err)
	}

	return nil
}

// HashTxBytes is a helper method to hash (Blake2) the transaction bytes before signing
func (c *PTBClient) HashTxBytes(txBytes []byte) []byte {
	intentMessage := append([]byte{0, 0, 0}, txBytes...)
	digest := blake2b.Sum256(intentMessage)
	return digest[:]
}
