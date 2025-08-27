package bind

import (
	"encoding/json"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

func TestGetCustomValueFromObjectData_Success(t *testing.T) {
	t.Parallel()
	mockedContent := `{
		"content": {
			"dataType": "moveObject",
			"fields": {
				"value": 10
			}
		}
	}`

	var data models.SuiObjectData
	if err := json.Unmarshal([]byte(mockedContent), &data); err != nil {
		t.Fatalf("failed to unmarshal mocked content: %v", err)
	}

	response := &models.SuiObjectResponse{
		Data: &data,
	}

	var target int
	err := GetCustomValueFromObjectData(response, &target)
	require.NoError(t, err)
	require.Equal(t, 10, target)
}
