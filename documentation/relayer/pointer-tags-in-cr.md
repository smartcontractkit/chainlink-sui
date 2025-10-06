# PointerTag Usage in ChainReader

## Overview

`PointerTag` is a feature in the Sui ChainReader module that enables automatic resolution of object IDs from owned objects. Instead of manually providing object IDs as parameters, the ChainReader can dynamically fetch them from the blockchain using pointer tags that reference specific objects owned by a contract.

## How it Works

When a function parameter has a `PointerTag` configured, the ChainReader:

1. **Validates the pointer tag** and extracts module, pointer name, field name, and derivation key
2. **Fetches owned objects** from the specified package on the Sui blockchain from the same contract being called
3. **Matches objects** based on the module and object type
4. **Extracts the parent object ID** from the field
5. **Derives the child object ID** using the derivation key and automatically populates the parameter

## Configuration Format

```go
type PointerTag struct {
    Module        string // e.g., "state_object", "offramp", "counter"
    PointerName   string // e.g., "CCIPObjectRefPointer", "OffRampStatePointer"
    FieldName     string // e.g., "ccip_object_id", "off_ramp_object_id"
    DerivationKey string // e.g., "CCIPObjectRef", "OffRampState", "Counter"
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

- `Module` - The Sui module name containing the object type
- `PointerName` - The object/struct type to search for (typically ends with "Pointer" in the contract)
- `FieldName` - The field within the pointer object containing the parent object ID
- `DerivationKey` - The key used to derive the child object ID from the parent object ID

**Note**: With the introduction of derived objects in Sui, pointer objects now store a parent object ID, and child object IDs are deterministically derived using derivation keys. The `FieldName` refers to the parent object ID field in the pointer, and `DerivationKey` specifies which child object to derive.

## Usage Example

```go
pointerTag := &codec.PointerTag{
    Module:        "counter",
    PointerName:   "CounterPointer",
    FieldName:     "counter_object_id",
    DerivationKey: "Counter",
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
- `FieldName: "counter_object_id"` - Field in CounterPointer containing the parent object ID
- `DerivationKey: "Counter"` - Derivation key to derive the Counter child object from the parent

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
- All PointerTag fields (Module, PointerName, FieldName, DerivationKey) must be correctly specified
- Field names must match exactly between the tag and the actual object structure