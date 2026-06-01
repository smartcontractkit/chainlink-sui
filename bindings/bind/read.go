package bind

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/block-vision/sui-go-sdk/models"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

func ReadObject(ctx context.Context, objectId string, chainClient client.BindingsClient) (*models.SuiObjectResponse, error) {
	normalizedId, err := bindutils.ConvertAddressToString(objectId)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID %v: %w", objectId, err)
	}

	obj, err := chainClient.ReadObjectId(ctx, normalizedId)
	if err != nil {
		return nil, fmt.Errorf("error getting object with id %s: %w", objectId, err)
	}

	resp := &models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			ObjectId: obj.GetObjectId(),
			Version:  strconv.FormatUint(obj.GetVersion(), 10),
			Digest:   obj.GetDigest(),
		},
	}

	if obj.GetJson() != nil {
		fields := obj.GetJson().AsInterface()
		if fieldMap, ok := fields.(map[string]any); ok {
			resp.Data.Content = &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: fieldMap,
				},
			}
		}
	}

	return resp, nil
}

func GetCustomValueFromObjectData[T any](resp *models.SuiObjectResponse, target *T) error {
	if resp == nil || resp.Data == nil || resp.Data.Content == nil {
		return fmt.Errorf("object does not contain any content")
	}

	if resp.Data.Content.SuiMoveObject.Fields == nil {
		return fmt.Errorf("object content does not have fields")
	}
	moveObject := resp.Data.Content.SuiMoveObject.Fields

	valueField, exists := moveObject["value"]
	if !exists {
		return fmt.Errorf("object does not contain a 'value' field")
	}

	jsonBytes, err := json.Marshal(valueField)
	if err != nil {
		return fmt.Errorf("failed to marshal value field: %w", err)
	}

	if err := json.Unmarshal(jsonBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal object data value: %w", err)
	}

	return nil
}
