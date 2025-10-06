package bind

import (
	"encoding/binary"
	"fmt"

	"github.com/smartcontractkit/chainlink-sui/bindings/utils"

	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"golang.org/x/crypto/blake2b"
)

const (
	HashingIntentScopeChildObjectID = 0xf0
	SuiFrameworkAddress             = "0x2"
)

// deriveDynamicFieldIDFromBytes computes a deterministic ObjectID for a dynamic field
// given its parent address, serialized key bytes and serialized key type tag bytes.
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
func deriveDynamicFieldIDFromBytes(parentAddress string, bcsKeyBytes []byte, bcsKeyTypeTagBytes []byte) (string, error) {
	normalizedParent, err := utils.ConvertAddressToString(parentAddress)
	if err != nil {
		return "", fmt.Errorf("invalid parent address: %w", err)
	}

	parentBytes, err := utils.ConvertStringToAddressBytes(normalizedParent)
	if err != nil {
		return "", fmt.Errorf("failed to convert parent address to bytes: %w", err)
	}

	hasher, err := blake2b.New256(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create Blake2b hasher: %w", err)
	}

	// Hash: intent_scope || parent || len(key) || key || key_type_tag
	hasher.Write([]byte{HashingIntentScopeChildObjectID})
	hasher.Write(parentBytes[:])

	keyLenBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyLenBytes, uint64(len(bcsKeyBytes)))
	hasher.Write(keyLenBytes)

	hasher.Write(bcsKeyBytes)
	hasher.Write(bcsKeyTypeTagBytes)

	hash := hasher.Sum(nil)

	objectID, err := utils.ConvertBytesToAddress(hash)
	if err != nil {
		return "", fmt.Errorf("failed to convert hash to address: %w", err)
	}

	return objectID, nil
}

// DeriveObjectIDWithVectorU8Key constructs the BCS bytes for DerivedObjectKey<vector<u8>> TypeTag.
// keyBytes should be the raw vector<u8> value - this function will BCS-serialize it.
func DeriveObjectIDWithVectorU8Key(parentAddress string, keyBytes []byte) (string, error) {
	// BCS-serialize the key value (adds length prefix)
	bcsKeyBytes, err := mystenbcs.Marshal(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to BCS serialize key bytes: %w", err)
	}

	suiFrameworkBytes, err := utils.ConvertStringToAddressBytes(SuiFrameworkAddress)
	if err != nil {
		return "", fmt.Errorf("failed to convert sui framework address to bytes: %w", err)
	}

	// Manually construct BCS bytes for: TypeTag::Struct(DerivedObjectKey<vector<u8>>)
	// This avoids the SDK's BCS encoder bug with nested TypeTag enums.
	//
	// BCS format breakdown:
	//   0x07                        - TypeTag::Struct variant
	//   [32 bytes]                  - address (0x2)
	//   0x0e + "derived_object"     - module name with length prefix
	//   0x10 + "DerivedObjectKey"   - struct name with length prefix
	//   0x01                        - type params count: 1
	//   0x06 + 0x01                 - TypeTag::Vector(TypeTag::U8)

	var typeTagBytes []byte
	typeTagBytes = append(typeTagBytes, 0x07)                          // TypeTag::Struct
	typeTagBytes = append(typeTagBytes, suiFrameworkBytes[:]...)       // address
	typeTagBytes = append(typeTagBytes, 0x0e)                          // module length
	typeTagBytes = append(typeTagBytes, []byte("derived_object")...)   // module name
	typeTagBytes = append(typeTagBytes, 0x10)                          // struct name length
	typeTagBytes = append(typeTagBytes, []byte("DerivedObjectKey")...) // struct name
	typeTagBytes = append(typeTagBytes, 0x01)                          // type params count
	typeTagBytes = append(typeTagBytes, 0x06, 0x01)                    // vector<u8> TypeTag

	return deriveDynamicFieldIDFromBytes(parentAddress, bcsKeyBytes, typeTagBytes)
}
