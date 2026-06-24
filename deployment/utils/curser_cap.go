package utils

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/block-vision/sui-go-sdk/utils"
	"golang.org/x/crypto/blake2b"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

type transactionChangedObjectsClient interface {
	GetTransactionChangedObjects(ctx context.Context, digest string) ([]*suirpcv2.ChangedObject, error)
}

// ResolveCurserCapObjectID returns the CurserCap registered in the fast MCMS Registry.
// mint_and_register_curser_cap stores the cap in the registry bag, so it may not appear as a
// top-level created object in transaction effects.
func ResolveCurserCapObjectID(
	ctx context.Context,
	suiClient client.SuiPTBClient,
	txDigest string,
	fastRegistryObjectID string,
	ccipPackageID string,
) (string, error) {
	if id, err := FindCurserCapObjectIDFromTx(ctx, suiClient, txDigest); err == nil {
		return id, nil
	}
	return FindCurserCapInFastRegistry(ctx, suiClient, fastRegistryObjectID, ccipPackageID)
}

// FindCurserCapObjectIDFromTx extracts the minted CurserCap object ID from transaction effects.
func FindCurserCapObjectIDFromTx(ctx context.Context, suiClient client.SuiPTBClient, txDigest string) (string, error) {
	fetcher, ok := suiClient.(transactionChangedObjectsClient)
	if !ok {
		return "", fmt.Errorf("sui client does not support fetching transaction object changes")
	}
	changed, err := fetcher.GetTransactionChangedObjects(ctx, txDigest)
	if err != nil {
		return "", fmt.Errorf("fetch transaction %s: %w", txDigest, err)
	}
	if id, err := findCurserCapInChangedObjects(changed); err == nil {
		return id, nil
	}

	// Fallback for legacy object-type layouts in unit tests.
	curserCapID, err := bind.FindObjectIDFromChangedObjects(changed, "rmn_remote", "CurserCap")
	if err != nil {
		return "", fmt.Errorf("find CurserCap in transaction %s: %w", txDigest, err)
	}
	return curserCapID, nil
}

// FindCurserCapInFastRegistry reads the CurserCap object ID from the fast MCMS Registry bag.
func FindCurserCapInFastRegistry(
	ctx context.Context,
	suiClient client.SuiPTBClient,
	fastRegistryObjectID string,
	ccipPackageID string,
) (string, error) {
	registryObj, err := suiClient.ReadObjectId(ctx, fastRegistryObjectID)
	if err != nil {
		return "", fmt.Errorf("read fast registry %s: %w", fastRegistryObjectID, err)
	}

	bagID, err := extractBagIDFromRegistry(registryObj)
	if err != nil {
		return "", err
	}

	normalizedCCIP, err := bindutils.ConvertAddressToString(ccipPackageID)
	if err != nil {
		return "", fmt.Errorf("normalize ccip package id: %w", err)
	}

	keyCandidates := []string{
		strings.TrimPrefix(normalizedCCIP, "0x"),
		normalizedCCIP,
	}
	seen := make(map[string]struct{}, len(keyCandidates))
	for _, key := range keyCandidates {
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		capID, lookupErr := lookupCurserCapInBag(ctx, suiClient, bagID, key)
		if lookupErr == nil {
			return capID, nil
		}
	}

	return "", fmt.Errorf("CurserCap not found in fast registry %s for ccip package %s", fastRegistryObjectID, normalizedCCIP)
}

func findCurserCapInChangedObjects(changed []*suirpcv2.ChangedObject) (string, error) {
	for _, obj := range changed {
		if obj == nil || obj.GetObjectId() == "" {
			continue
		}
		objectType := obj.GetObjectType()
		if idx := strings.Index(objectType, "<"); idx != -1 {
			objectType = objectType[:idx]
		}
		if strings.HasSuffix(objectType, "::CurserCap") {
			return obj.GetObjectId(), nil
		}
	}
	return "", fmt.Errorf("no changed object with type suffix ::CurserCap")
}

func extractBagIDFromRegistry(registryObj *suirpcv2.Object) (string, error) {
	if registryObj == nil || registryObj.GetJson() == nil {
		return "", fmt.Errorf("registry object has no json content")
	}
	root, ok := registryObj.GetJson().AsInterface().(map[string]any)
	if !ok {
		return "", fmt.Errorf("registry json is not an object")
	}
	packageCaps, ok := root["package_caps"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("registry missing package_caps field")
	}
	bagID, ok := findNestedObjectID(packageCaps, "id")
	if !ok {
		return "", fmt.Errorf("registry package_caps bag id not found in json")
	}
	return bagID, nil
}

func findNestedObjectID(node map[string]any, key string) (string, bool) {
	if raw, ok := node[key]; ok {
		switch typed := raw.(type) {
		case map[string]any:
			if id, ok := typed["id"].(string); ok && id != "" {
				return id, true
			}
		case string:
			if typed != "" {
				return typed, true
			}
		}
	}
	if fields, ok := node["fields"].(map[string]any); ok {
		if id, ok := findNestedObjectID(fields, key); ok {
			return id, true
		}
	}
	for _, child := range node {
		childMap, ok := child.(map[string]any)
		if !ok {
			continue
		}
		if id, ok := findNestedObjectID(childMap, key); ok {
			return id, true
		}
	}
	return "", false
}

func lookupCurserCapInBag(ctx context.Context, suiClient client.SuiPTBClient, bagID, packageAddressKey string) (string, error) {
	dfID, err := deriveAsciiStringDynamicFieldID(bagID, packageAddressKey)
	if err != nil {
		return "", err
	}

	dfObj, err := suiClient.ReadObjectId(ctx, dfID)
	if err != nil {
		return "", err
	}

	return extractCurserCapIDFromDynamicField(dfObj)
}

func extractCurserCapIDFromDynamicField(dfObj *suirpcv2.Object) (string, error) {
	if dfObj == nil || dfObj.GetJson() == nil {
		return "", fmt.Errorf("dynamic field object has no json content")
	}
	root, ok := dfObj.GetJson().AsInterface().(map[string]any)
	if !ok {
		return "", fmt.Errorf("dynamic field json is not an object")
	}
	if value, ok := root["value"].(map[string]any); ok {
		if capID, ok := findNestedObjectID(value, "id"); ok {
			return capID, nil
		}
	}
	if fields, ok := root["fields"].(map[string]any); ok {
		if value, ok := fields["value"].(map[string]any); ok {
			if capID, ok := findNestedObjectID(value, "id"); ok {
				return capID, nil
			}
		}
	}
	if capID, ok := findNestedObjectID(root, "id"); ok {
		return capID, nil
	}
	return "", fmt.Errorf("dynamic field does not contain CurserCap id")
}

func deriveAsciiStringDynamicFieldID(parentObjectID string, key string) (string, error) {
	keyBytes, err := mystenbcs.Marshal([]byte(key))
	if err != nil {
		return "", fmt.Errorf("bcs-serialize ascii key: %w", err)
	}

	typeTagBytes, err := bcsEncodeAsciiStringTypeTag()
	if err != nil {
		return "", err
	}

	return deriveDynamicFieldIDFromBytes(parentObjectID, keyBytes, typeTagBytes)
}

func bcsEncodeAsciiStringTypeTag() ([]byte, error) {
	stdAddr, err := transaction.ConvertSuiAddressStringToBytes(utils.NormalizeSuiAddress("0x1"))
	if err != nil {
		return nil, fmt.Errorf("convert 0x1 address: %w", err)
	}

	var out []byte
	out = append(out, 0x07) // TypeTag::Struct
	out = append(out, stdAddr[:]...)
	out = append(out, byte(len("ascii")))
	out = append(out, []byte("ascii")...)
	out = append(out, byte(len("String")))
	out = append(out, []byte("String")...)
	out = append(out, 0x00) // no type parameters
	return out, nil
}

func deriveDynamicFieldIDFromBytes(parentAddress string, bcsKeyBytes []byte, bcsKeyTypeTagBytes []byte) (string, error) {
	normalizedParent := utils.NormalizeSuiAddress(parentAddress)
	parentBytes, err := transaction.ConvertSuiAddressStringToBytes(normalizedParent)
	if err != nil {
		return "", fmt.Errorf("convert parent address to bytes: %w", err)
	}

	hasher, err := blake2b.New256(nil)
	if err != nil {
		return "", fmt.Errorf("create blake2b hasher: %w", err)
	}

	hasher.Write([]byte{0xf0})
	hasher.Write(parentBytes[:])

	keyLenBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyLenBytes, uint64(len(bcsKeyBytes)))
	hasher.Write(keyLenBytes)
	hasher.Write(bcsKeyBytes)
	hasher.Write(bcsKeyTypeTagBytes)

	hash := hasher.Sum(nil)
	addrBytes := transaction.ConvertSuiAddressBytesToString(models.SuiAddressBytes(hash))
	return string(addrBytes), nil
}
