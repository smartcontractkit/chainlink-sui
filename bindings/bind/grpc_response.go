package bind

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	v2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const defaultGasPrice uint64 = 10_000

func mapExecuteResponseToModels(resp *suirpcv2.ExecuteTransactionResponse) (*models.SuiTransactionBlockResponse, error) {
	if resp == nil || resp.Transaction == nil {
		return nil, errors.New("empty execute transaction response")
	}

	tx := resp.Transaction
	result := &models.SuiTransactionBlockResponse{
		Digest: tx.GetDigest(),
	}

	if tx.Effects != nil && tx.Effects.Status != nil {
		status := "success"
		if !tx.Effects.Status.GetSuccess() {
			status = "failure"
		}
		result.Effects.Status.Status = status
		if errMsg := tx.Effects.Status.GetError(); errMsg != nil {
			result.Effects.Status.Error = errMsg.GetDescription()
		}
	}

	if tx.Effects != nil {
		result.ObjectChanges = mapChangedObjectsToModels(tx.Effects.GetChangedObjects())
	}

	return result, nil
}

func mapChangedObjectsToModels(changed []*suirpcv2.ChangedObject) []models.ObjectChange {
	if len(changed) == 0 {
		return nil
	}

	out := make([]models.ObjectChange, 0, len(changed))
	for _, obj := range changed {
		if obj == nil {
			continue
		}

		change := models.ObjectChange{
			ObjectId:   obj.GetObjectId(),
			ObjectType: obj.GetObjectType(),
			Digest:     obj.GetOutputDigest(),
			Version:    strconv.FormatUint(obj.GetOutputVersion(), 10),
		}
		if obj.GetInputVersion() != 0 {
			change.PreviousVersion = strconv.FormatUint(obj.GetInputVersion(), 10)
		}

		if obj.GetObjectType() == "package" ||
			obj.GetOutputState() == suirpcv2.ChangedObject_OUTPUT_OBJECT_STATE_PACKAGE_WRITE {
			change.Type = "published"
			change.PackageId = obj.GetObjectId()
		} else {
			switch obj.GetIdOperation() {
			case suirpcv2.ChangedObject_CREATED:
				change.Type = "created"
			case suirpcv2.ChangedObject_DELETED:
				change.Type = "deleted"
			default:
				change.Type = "mutated"
			}
		}

		if owner := obj.GetOutputOwner(); owner != nil {
			change.Owner = mapOwnerToInterface(owner)
		}

		out = append(out, change)
	}

	return out
}

func mapOwnerToInterface(owner *suirpcv2.Owner) any {
	if owner == nil {
		return nil
	}

	switch owner.GetKind() {
	case suirpcv2.Owner_ADDRESS:
		return map[string]string{"AddressOwner": owner.GetAddress()}
	case suirpcv2.Owner_OBJECT:
		return map[string]string{"ObjectOwner": owner.GetAddress()}
	case suirpcv2.Owner_SHARED:
		return map[string]any{
			"Shared": map[string]uint64{
				"initial_shared_version": owner.GetVersion(),
			},
		}
	case suirpcv2.Owner_IMMUTABLE:
		return "Immutable"
	default:
		return nil
	}
}

func mapGrpcObjectToResolved(obj *suirpcv2.Object) (*resolvedObject, error) {
	if obj == nil {
		return nil, errors.New("object is nil")
	}

	resolved := &resolvedObject{
		ObjectId: obj.GetObjectId(),
		Version:  obj.GetVersion(),
		Digest:   obj.GetDigest(),
	}

	if owner := obj.GetOwner(); owner != nil {
		switch owner.GetKind() {
		case suirpcv2.Owner_ADDRESS:
			resolved.Owner.AddressOwner = owner.GetAddress()
		case suirpcv2.Owner_OBJECT:
			resolved.Owner.ObjectOwner = owner.GetAddress()
		case suirpcv2.Owner_SHARED:
			if owner.GetVersion() != 0 {
				v := owner.GetVersion()
				resolved.InitialSharedVersion = &v
				resolved.Owner.Shared.InitialSharedVersion = v
			}
		case suirpcv2.Owner_IMMUTABLE:
			resolved.Owner = models.ObjectOwner{}
		}
	}

	return resolved, nil
}

func mapGrpcCoinToObjectRef(coin *suirpcv2.Object) *models.SuiObjectRef {
	if coin == nil {
		return nil
	}

	return &models.SuiObjectRef{
		ObjectId: coin.GetObjectId(),
		Version:  coin.GetVersion(),
		Digest:   coin.GetDigest(),
	}
}

func decodeSignatureStrings(signatures []string) ([][]byte, error) {
	out := make([][]byte, 0, len(signatures))
	for i, sig := range signatures {
		raw, err := base64.StdEncoding.DecodeString(sig)
		if err != nil {
			return nil, fmt.Errorf("failed to decode signature %d: %w", i, err)
		}
		out = append(out, raw)
	}

	return out, nil
}

func transactionChangedObjectsReadMaskPaths() []string {
	return []string{
		"digest",
		"effects.status",
		"effects.changed_objects",
		"effects.changed_objects.object_id",
		"effects.changed_objects.object_type",
		"effects.changed_objects.id_operation",
		"effects.changed_objects.output_state",
		"effects.changed_objects.output_owner",
		"effects.changed_objects.output_version",
		"effects.changed_objects.output_digest",
		"effects.changed_objects.input_version",
	}
}

func buildExecuteRequest(bcsBytes []byte, signatures []string) (*v2.ExecuteTransactionRequest, error) {
	sigBytes, err := decodeSignatureStrings(signatures)
	if err != nil {
		return nil, err
	}

	userSigs := make([]*v2.UserSignature, len(sigBytes))
	for i, sig := range sigBytes {
		userSigs[i] = &v2.UserSignature{Bcs: &v2.Bcs{Value: sig}}
	}

	return &v2.ExecuteTransactionRequest{
		Transaction: &v2.Transaction{Bcs: &v2.Bcs{Value: bcsBytes}},
		Signatures:  userSigs,
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: transactionChangedObjectsReadMaskPaths(),
		},
	}, nil
}
