package bind

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

type ObjectResolver struct {
	client client.BindingsClient
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

func NewObjectResolver(chainClient client.BindingsClient) *ObjectResolver {
	return &ObjectResolver{
		client: chainClient,
		cache: &objectCache{
			cache: make(map[string]*resolvedObject),
		},
	}
}

func GetSharedObject(ctx context.Context, chainClient client.BindingsClient, objectId string) (*Object, error) {
	resolver := NewObjectResolver(chainClient)
	return resolver.GetSharedObject(ctx, objectId)
}

func (r *ObjectResolver) GetSharedObject(ctx context.Context, objectId string) (*Object, error) {
	normalizedId, err := bindutils.ConvertAddressToString(objectId)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID %s: %w", objectId, err)
	}

	resolved, err := r.resolveObject(ctx, normalizedId)
	if err != nil {
		return nil, err
	}

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
	if cached := r.cache.get(objectId); cached != nil {
		return cached, nil
	}

	obj, err := r.client.ReadObjectId(ctx, objectId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch object %s: %w", objectId, err)
	}

	resolved, err := mapGrpcObjectToResolved(obj)
	if err != nil {
		return nil, err
	}

	r.cache.set(objectId, resolved)

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
