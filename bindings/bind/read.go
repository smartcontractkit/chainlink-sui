package bind

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

func ReadObject(ctx context.Context, objectId string, client sui.ISuiAPI) (*models.SuiObjectResponse, error) {
	// Normalize the object ID
	normalizedId, err := bindutils.ConvertAddressToString(objectId)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID %v: %w", objectId, err)
	}

	req := models.SuiGetObjectRequest{
		ObjectId: normalizedId,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
			ShowOwner:   true,
			ShowType:    true,
		},
	}

	object, err := client.SuiGetObject(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error getting object with id %s: %w", objectId, err)
	}

	// Return the object data
	return &object, nil
}

// Decodes "value" field into a user provided pointer of any type.
func GetCustomValueFromObjectData[T any](resp *models.SuiObjectResponse, target *T) error {
	if resp == nil || resp.Data == nil || resp.Data.Content == nil {
		return fmt.Errorf("object does not contain any content")
	}

	// The content should be a SuiMoveObject which embeds Fields
	// is already a map[string]interface{}
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
