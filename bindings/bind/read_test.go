package bind

import (
	"encoding/json"
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
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

	var data suiclient.SuiObjectData
	if err := json.Unmarshal([]byte(mockedContent), &data); err != nil {
		t.Fatalf("failed to unmarshal mocked content: %v", err)
	}

	var target int
	err := GetCustomValueFromObjectData(data, &target)
	require.NoError(t, err)
	require.Equal(t, 10, target)
}
