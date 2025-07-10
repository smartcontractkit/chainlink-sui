package bind

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

type ObjectResolver struct {
	client sui.ISuiAPI
	cache  *objectCache
}

type objectCache struct {
	mu    sync.RWMutex
	cache map[string]*resolvedObject
}

type resolvedObject struct {
	ObjectId             string
	Version              uint64
	Digest               string
	Owner                models.ObjectOwner
	InitialSharedVersion *uint64
}

func NewObjectResolver(client sui.ISuiAPI) *ObjectResolver {
	return &ObjectResolver{
		client: client,
		cache: &objectCache{
			cache: make(map[string]*resolvedObject),
		},
	}
}

func GetObject(ctx context.Context, client sui.ISuiAPI, objectId string) (*Object, error) {
	resolver := NewObjectResolver(client)
	return resolver.GetObject(ctx, objectId)
}

func (r *ObjectResolver) GetObject(ctx context.Context, objectId string) (*Object, error) {
	normalizedId, err := bindutils.ConvertAddressToString(objectId)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID %s: %w", objectId, err)
	}

	if cached := r.cache.get(normalizedId); cached != nil {
		return r.createObjectFromResolved(cached), nil
	}

	resolved, err := r.resolveObject(ctx, normalizedId)
	if err != nil {
		return nil, err
	}

	r.cache.set(normalizedId, resolved)

	return r.createObjectFromResolved(resolved), nil
}

func (r *ObjectResolver) ResolveCallArg(ctx context.Context, arg *transaction.CallArg, typeName string) (*transaction.CallArg, error) {
	if arg == nil {
		return nil, errors.New("nil CallArg")
	}

	if arg.UnresolvedPure != nil {
		return nil, errors.New("cannot handle UnresolvedPure")
	}

	if arg.UnresolvedObject != nil {
		objectId := fmt.Sprintf("0x%x", arg.UnresolvedObject.ObjectId)

		resolved, err := r.resolveObject(ctx, objectId)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve UnresolvedObject %s: %w", objectId, err)
		}

		isMutable := strings.HasPrefix(typeName, "&mut ")
		objectArg, err := r.createObjectArgWithMutability(resolved, isMutable)
		if err != nil {
			return nil, err
		}

		return &transaction.CallArg{
			Object: objectArg,
		}, nil
	}

	if arg.Pure == nil && arg.Object == nil {
		return nil, errors.New("invalid call arg, no Pure or Object field")
	}

	return arg, nil
}

func (r *ObjectResolver) resolveObject(ctx context.Context, objectId string) (*resolvedObject, error) {
	resp, err := r.client.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: objectId,
		Options: models.SuiObjectDataOptions{
			ShowOwner:               true,
			ShowType:                true,
			ShowPreviousTransaction: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch object %s: %w", objectId, err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("object error for %s: %v", objectId, resp.Error)
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("object %s not found", objectId)
	}

	version, err := parseVersionString(resp.Data.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version for object %s: %w", objectId, err)
	}

	resolved := &resolvedObject{
		ObjectId: resp.Data.ObjectId,
		Version:  version,
		Digest:   resp.Data.Digest,
	}

	if resp.Data.Owner != nil {
		// TODO: check this logic, use mapstructure if this is a map[string]any{}
		ownerBytes, err := json.Marshal(resp.Data.Owner)
		if err == nil {
			var owner models.ObjectOwner
			if err := json.Unmarshal(ownerBytes, &owner); err == nil {
				resolved.Owner = owner

				if owner.AddressOwner != "" {
					resolved.InitialSharedVersion = nil
				} else if owner.ObjectOwner != "" {
					resolved.InitialSharedVersion = nil
				} else if owner.Shared.InitialSharedVersion > 0 {
					v := owner.Shared.InitialSharedVersion
					resolved.InitialSharedVersion = &v
				}
			} else {
				var immutableOwner string
				if err := json.Unmarshal(ownerBytes, &immutableOwner); err == nil && immutableOwner == "Immutable" {
					resolved.Owner = models.ObjectOwner{}
				}
			}
		}
	}

	return resolved, nil
}

func (r *ObjectResolver) createObjectFromResolved(resolved *resolvedObject) *Object {
	return &Object{
		Id:                   resolved.ObjectId,
		InitialSharedVersion: resolved.InitialSharedVersion,
	}
}

func (r *ObjectResolver) createObjectArgWithMutability(resolved *resolvedObject, isMutable bool) (*transaction.ObjectArg, error) {
	objIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(resolved.ObjectId))
	if err != nil {
		return nil, fmt.Errorf("failed to convert object ID to bytes: %w", err)
	}

	digestBytes, err := bindutils.ConvertStringToDigestBytes(resolved.Digest)
	if err != nil {
		return nil, fmt.Errorf("failed to convert digest to bytes: %w", err)
	}

	if resolved.Owner.Shared.InitialSharedVersion > 0 {
		if resolved.InitialSharedVersion == nil {
			return nil, fmt.Errorf("shared object %s missing initial shared version", resolved.ObjectId)
		}

		return &transaction.ObjectArg{
			SharedObject: &transaction.SharedObjectRef{
				ObjectId:             *objIdBytes,
				InitialSharedVersion: *resolved.InitialSharedVersion,
				Mutable:              isMutable,
			},
		}, nil
	}

	if resolved.Owner.AddressOwner != "" {
		return &transaction.ObjectArg{
			ImmOrOwnedObject: &transaction.SuiObjectRef{
				ObjectId: *objIdBytes,
				Version:  resolved.Version,
				Digest:   *digestBytes,
			},
		}, nil
	}

	if resolved.Owner.ObjectOwner != "" {
		return &transaction.ObjectArg{
			ImmOrOwnedObject: &transaction.SuiObjectRef{
				ObjectId: *objIdBytes,
				Version:  resolved.Version,
				Digest:   *digestBytes,
			},
		}, nil
	}

	return &transaction.ObjectArg{
		ImmOrOwnedObject: &transaction.SuiObjectRef{
			ObjectId: *objIdBytes,
			Version:  resolved.Version,
			Digest:   *digestBytes,
		},
	}, nil
}

func (c *objectCache) get(objectId string) *resolvedObject {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cache[objectId]
}

func (c *objectCache) set(objectId string, resolved *resolvedObject) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[objectId] = resolved
}

func (r *ObjectResolver) ClearCache() {
	r.cache.mu.Lock()
	defer r.cache.mu.Unlock()
	r.cache.cache = make(map[string]*resolvedObject)
}
