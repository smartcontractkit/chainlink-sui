package testutils

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
)

// HasCounterObject checks if a counter object exists for the given packageId.
func HasCounterObject(t *testing.T, c sui.ISuiAPI, packageId string) bool {
	t.Helper()
	ctx := context.Background()

	// Query for objects owned by the package address
	resp, err := c.SuiXGetOwnedObjects(ctx, models.SuiXGetOwnedObjectsRequest{
		Address: packageId,
		Query: models.SuiObjectResponseQuery{
			Filter: map[string]interface{}{
				"MatchType": map[string]interface{}{
					"TypeName": packageId + "::counter::Counter",
				},
			},
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
			},
		},
	})
	require.NoError(t, err)

	return len(resp.Data) > 0
}

// ReadCounterValue gets the current value of a counter object.
func ReadCounterValue(t *testing.T, c sui.ISuiAPI, counterId string) int {
	t.Helper()
	ctx := context.Background()

	resp, err := c.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: counterId,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	})
	require.NoError(t, err)

	// Access the fields directly from the SuiMoveObject embedded in SuiParsedData
	value, ok := resp.Data.Content.Fields["value"].(float64)
	require.True(t, ok, "Failed to parse counter value")

	return int(value)
}

// FindCounterID locates the counter object ID for a given package ID.
func FindCounterID(t *testing.T, c sui.ISuiAPI, packageId string) string {
	t.Helper()
	ctx := context.Background()

	// Query for objects owned by the package address
	resp, err := c.SuiXGetOwnedObjects(ctx, models.SuiXGetOwnedObjectsRequest{
		Address: packageId,
		Query: models.SuiObjectResponseQuery{
			Filter: map[string]interface{}{
				"MatchType": map[string]interface{}{
					"TypeName": packageId + "::counter::Counter",
				},
			},
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
			},
		},
	})
	require.NoError(t, err)
	require.Greater(t, len(resp.Data), 0, "No counter objects found")

	return resp.Data[0].Data.ObjectId
}

// parseStringToUint64 parses a string to uint64
func parseStringToUint64(str string) (uint64, error) {
	var result uint64
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		return 0, err
	}
	return result, nil
}
