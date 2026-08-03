package bind

import (
	"testing"

	"github.com/stretchr/testify/require"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
)

func TestMapChangedObjectsToModels_PublishedPackageWithCreatedIdOperation(t *testing.T) {
	objectID := "0xabc"
	changed := []*suirpcv2.ChangedObject{
		{
			ObjectId:    &objectID,
			ObjectType:  new("package"),
			IdOperation: suirpcv2.ChangedObject_CREATED.Enum(),
			OutputState: suirpcv2.ChangedObject_OUTPUT_OBJECT_STATE_PACKAGE_WRITE.Enum(),
		},
	}

	out := mapChangedObjectsToModels(changed)
	require.Len(t, out, 1)
	require.Equal(t, "published", out[0].Type)
	require.Equal(t, objectID, out[0].PackageId)
}
