//go:build integration

package tests

import (
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
)

var DEFAULT_GAS_BUDGET uint64 = 1000000000

func FindCreatedObject(objectChanges []models.ObjectChange, objectPattern string) (string, *uint64, error) {
	var targetModule, targetObject string

	objectPattern = strings.TrimPrefix(objectPattern, "::")

	parts := strings.Split(objectPattern, "::")
	if len(parts) >= 2 {
		targetModule = parts[len(parts)-2]
		targetObject = parts[len(parts)-1]
	} else if len(parts) == 1 {
		targetObject = parts[0]
	}

	for _, change := range objectChanges {
		if change.Type == "created" && change.ObjectType != "" {
			objectType := change.ObjectType

			// strip generics since they contain '::' substrings
			baseType := objectType
			if genericStart := strings.Index(objectType, "<"); genericStart != -1 {
				baseType = objectType[:genericStart]
			}

			// split the type
			typeParts := strings.Split(baseType, "::")
			if len(typeParts) >= 3 {
				moduleName := typeParts[len(typeParts)-2]
				objectName := typeParts[len(typeParts)-1]

				moduleMatch := targetModule == "" || moduleName == targetModule
				objectMatch := objectName == targetObject

				if moduleMatch && objectMatch {
					objectId := change.ObjectId
					share := change.GetObjectOwnerShare()
					if share.InitialSharedVersion > 0 {
						version := share.InitialSharedVersion
						return objectId, &version, nil
					}

					return objectId, nil, nil
				}
			}
		}
	}

	return "", nil, fmt.Errorf("could not find object with pattern %s", objectPattern)
}
