package chainwriter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
)

func TestGetOWnedObject(t *testing.T) {

	client := sui.NewSuiClient("")

	ownedObjectsReq := models.SuiXGetOwnedObjectsRequest{
		Address: "0xdec6eef8180d0a62a1dab9182c28bb854cfe0cf013d24344c3e6cb3de6b10572",
		Query: models.SuiObjectResponseQuery{
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
				ShowType:    true,
				ShowOwner:   true,
			},
		},
		Limit: uint64(50),
	}

	response, err := client.SuiXGetOwnedObjects(context.Background(), ownedObjectsReq)
	require.NoError(t, err)

	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("RESP", string(data))

}
