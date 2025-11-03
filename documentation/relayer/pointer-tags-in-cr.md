# PointerTag Usage in ChainReader

## Overview

`PointerTag` is a feature in the Sui ChainReader module that enables automatic resolution of object IDs from owned objects. Instead of manually providing object IDs as parameters, the ChainReader can dynamically fetch them from the blockchain using pointer tags that reference specific objects owned by a contract.

## How it Works

When a function parameter has a `PointerTag` configured, the ChainReader:

1. **Validates the pointer tag** and extracts module, pointer name, and derivation key
2. **Looks up the parent field name** from the global `common.PointerConfigs` registry based on the pointer name
3. **Fetches owned objects** from the specified package on the Sui blockchain from the same contract being called
4. **Matches objects** based on the module and object type
5. **Extracts the parent object ID** from the field (as defined in the registry)
6. **Derives the child object ID** using the derivation key and automatically populates the parameter

## Configuration Format

```go
type PointerTag struct {
    Module        string // e.g., "state_object", "offramp", "counter"
    PointerName   string // e.g., "CCIPObjectRefPointer", "OffRampStatePointer"
    FieldName     string // OPTIONAL: Ignored in favor of registry lookup
    DerivationKey string // e.g., "CCIPObjectRef", "OffRampState", "Counter"
    PackageID     string // OPTIONAL: Override for cross-package pointers
}

type SuiFunctionParam struct {
    Type       string      // Parameter type (e.g., "object_id")
    Name       string      // Parameter name
    PointerTag *PointerTag // Optional pointer tag for automatic object resolution
    Required   bool        // Whether parameter is required
    // ... other fields
}
```

### PointerTag Fields

- `Module` - The Sui module name containing the object type (required)
- `PointerName` - The object/struct type to search for, typically ends with "Pointer" (required)
- `FieldName` - **[OPTIONAL/IGNORED]** This field is not used by the implementation. The parent field name is automatically looked up from the global `common.PointerConfigs` registry based on the `PointerName`
- `DerivationKey` - The key used to derive the child object ID from the parent object ID (required)
- `PackageID` - **[OPTIONAL]** Override the package ID for cross-package pointer dependencies. If empty, the calling contract's package ID is used

**Note**: With the introduction of derived objects in Sui, pointer objects now store a parent object ID, and child object IDs are deterministically derived using derivation keys. The parent field name is looked up from `common.PointerConfigs`, and `DerivationKey` specifies which child object to derive.

## Pointer Registry

All pointer types must be registered in `relayer/common/pointer_config.go` in the `PointerConfigs` map before they can be used. This registry is the single source of truth for pointer configurations and defines:

- **Module** - The Sui module containing the pointer object
- **Pointer** - The pointer object type name (e.g., "OffRampStatePointer", "CounterPointer")
- **ParentFieldName** - The field in the pointer object containing the parent object ID

### Adding New Pointer Types

When adding support for a new pointer type, you must update `common.PointerConfigs`:

```go
// In relayer/common/pointer_config.go
var PointerConfigs = map[string][]PointerConfig{
    "mycontract": {
        {
            Module:          "mymodule",
            Pointer:         "MyPointer",
            ParentFieldName: "my_parent_object_id",
        },
    },
}
```

The registry key should be the contract/module name (case-insensitive). Once registered, the pointer type can be used in PointerTags without specifying the field name - it will be automatically looked up.

## Usage Example

```go
pointerTag := &codec.PointerTag{
    Module:        "counter",
    PointerName:   "CounterPointer",
    DerivationKey: "Counter",
    // FieldName is optional and will be looked up from common.PointerConfigs
    // which maps "CounterPointer" -> "counter_object_id"
}

// Function configuration
"get_count_using_pointer": {
    Name:          "get_count_using_pointer",
    SignerAddress: accountAddress,
    Params: []codec.SuiFunctionParam{
        {
            Type:       "object_id",
            Name:       "counter_id",
            PointerTag: pointerTag,
            Required:   true,
        },
    },
}
```

**Breaking down the pointer tag components:**
- `Module: "counter"` - Module name
- `PointerName: "CounterPointer"` - Pointer object type
- `DerivationKey: "Counter"` - Derivation key to derive the Counter child object from the parent
- Field name is automatically looked up from `common.PointerConfigs` ("counter_object_id" for CounterPointer)

> __IMPORTANT__: the pointer object MUST be owned by the contract.

### Calling the Function

When using PointerTag, no explicit parameters are needed:

```go
err = chainReader.GetLatestValue(
    context.Background(),
    strings.Join([]string{packageId, "counter", "get_count_using_pointer"}, "-"),
    primitives.Finalized,
    map[string]any{}, // Empty - parameter populated automatically
    &retUint64,
)
```

## Implementation Details

The ChainReader's `prepareArguments` function:

1. **Pre-loads parent object IDs** during `Bind()` for known pointer types (OffRamp, OnRamp, CCIP, Counter)
2. **Identifies pointer parameters** by checking for non-nil `PointerTag`
3. **Validates pointer tag** using `PointerTag.Validate()` method
4. **Builds pointer queries** by grouping tags by module/object type
5. **Retrieves parent object IDs** from cache (or fetches on-demand if not cached)
6. **Derives child object IDs** using `DeriveObjectIDWithVectorU8Key(parentID, derivationKey)`
7. **Populates argument map** with derived object IDs


## Benefits

- **Automatic object resolution** - No need to manually track object IDs
- **Dynamic parameter population** - Objects are resolved at call time
- **Simplified API calls** - Reduces the complexity of function invocations
- **Type safety** - Automatic conversion to appropriate object types (`bind.Object{Id: value}` for `object_id` type)
- **Performance optimized** - Parent object IDs are pre-loaded at binding time, reducing RPC calls
- **Deterministic derivation** - Child object IDs are derived offchain without additional RPC calls

## Limitations

- Only works with objects owned by the contract package
- Requires objects to exist and be accessible via `ReadOwnedObjects`
- Pointer types must be pre-registered in `common.PointerConfigs` before use
- The PointerTag's Module, PointerName, and DerivationKey must be correctly specified
- Field names in the registry must match the actual on-chain object structure exactly