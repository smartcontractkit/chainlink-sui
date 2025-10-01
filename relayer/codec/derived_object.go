package codec

import (
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/chainlink-sui/bindings/utils"

	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/transaction"
	"golang.org/x/crypto/blake2b"
)

const (
	HashingIntentScopeChildObjectID = 0xf0

	SuiFrameworkAddress = "0x2"
)

// DeriveDynamicFieldID computes a deterministic ObjectID for a dynamic field
// given its parent address, key type tag, and serialized key bytes.
//
// This mirrors the Sui Rust implementation from:
// sui-types/src/dynamic_field.rs:derive_dynamic_field_id()
//
// Algorithm:
//
//	hash = Blake2b256(
//	    0xf0 +                          // HashingIntentScope::ChildObjectId
//	    parent_address +                // 32 bytes
//	    len(key_bytes) as little-endian + // 8 bytes
//	    key_bytes +                     // BCS-serialized key
//	    bcs(key_type_tag)              // BCS-serialized TypeTag
//	)
//	result = hash[0:32]  // First 32 bytes = ObjectID
func DeriveDynamicFieldID(parentAddress string, keyTypeTag *transaction.TypeTag, keyBytes []byte) (string, error) {
	normalizedParent, err := utils.ConvertAddressToString(parentAddress)
	if err != nil {
		return "", fmt.Errorf("invalid parent address: %w", err)
	}

	parentBytes, err := utils.ConvertStringToAddressBytes(normalizedParent)
	if err != nil {
		return "", fmt.Errorf("failed to convert parent address to bytes: %w", err)
	}

	keyTypeTagBytes, err := mystenbcs.Marshal(keyTypeTag)
	if err != nil {
		return "", fmt.Errorf("failed to BCS serialize key type tag: %w", err)
	}

	hasher, err := blake2b.New256(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create Blake2b hasher: %w", err)
	}

	// Hash: intent_scope || parent || len(key) || key || key_type_tag
	hasher.Write([]byte{HashingIntentScopeChildObjectID})
	hasher.Write(parentBytes[:])

	keyLenBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyLenBytes, uint64(len(keyBytes)))
	hasher.Write(keyLenBytes)

	hasher.Write(keyBytes)
	hasher.Write(keyTypeTagBytes)

	hash := hasher.Sum(nil)

	objectID, err := utils.ConvertBytesToAddress(hash)
	if err != nil {
		return "", fmt.Errorf("failed to convert hash to address: %w", err)
	}

	return objectID, nil
}

func DeriveObjectID(parentAddress string, keyTypeTag *transaction.TypeTag, keyBytes []byte) (string, error) {
	suiFrameworkBytes, err := utils.ConvertStringToAddressBytes(SuiFrameworkAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert sui framework address to bytes: %w", err)
	}

	// Wrap the key type in DerivedObjectKey<K>
	wrapperTypeTag := &transaction.TypeTag{
		Struct: &transaction.StructTag{
			Address:    *suiFrameworkBytes,
			Module:     "derived_object",
			Name:       "DerivedObjectKey",
			TypeParams: []*transaction.TypeTag{keyTypeTag},
		},
	}

	return DeriveDynamicFieldID(parentAddress, wrapperTypeTag, keyBytes)
}

// DeriveDerivedObjectID computes the deterministic ObjectID created with sui::derived_object::claim().
func DeriveDerivedObjectID(parentObjectId string, keyPackageId string, keyModule string, keyStructName string, keyValue []byte) (string, error) {
	normalizedPackageID, err := utils.ConvertAddressToString(keyPackageId)
	if err != nil {
		return "", fmt.Errorf("invalid key package ID: %w", err)
	}

	packageBytes, err := utils.ConvertStringToAddressBytes(normalizedPackageID)
	if err != nil {
		return "", fmt.Errorf("failed to convert package ID to bytes: %w", err)
	}

	keyTypeTag := &transaction.TypeTag{
		Struct: &transaction.StructTag{
			Address:    *packageBytes,
			Module:     keyModule,
			Name:       keyStructName,
			TypeParams: []*transaction.TypeTag{},
		},
	}

	return DeriveObjectID(parentObjectId, keyTypeTag, keyValue)
}
