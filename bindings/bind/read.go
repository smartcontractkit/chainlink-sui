package bind

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pattonkan/sui-go/suiclient"
)

func ReadObject(ctx context.Context, objectId string, client suiclient.ClientImpl) (*suiclient.SuiObjectResponse, error) {
	address, err := ToSuiAddress(objectId)
	if err != nil {
		return nil, err
	}

	object, err := client.GetObject(ctx, &suiclient.GetObjectRequest{
		ObjectId: address,
		Options: &suiclient.SuiObjectDataOptions{
			ShowContent: true,
			ShowOwner:   true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting object with id %s: %w", objectId, err)
	}

	// Return the object data
	return object, nil
}

// Decodes "value" field into a user provided pointer of any type.
func GetCustomValueFromObjectData[T any](data suiclient.SuiObjectData, target *T) error {
	if data.Content == nil || data.Content.Data.MoveObject.Fields == nil {
		return fmt.Errorf("object does not contain any fields")
	}

	type Data struct {
		Value T `json:"value"`
	}

	var valueData Data
	if err := json.Unmarshal(data.Content.Data.MoveObject.Fields, &valueData); err != nil {
		return fmt.Errorf("failed to unmarshal object data value: %w", err)
	}

	// Assign the value to the target
	*target = valueData.Value

	return nil
}
