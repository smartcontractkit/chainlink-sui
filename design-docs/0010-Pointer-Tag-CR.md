# PointerTag Usage in ChainReader

## Overview

`PointerTag` is a feature in the Sui ChainReader module that enables automatic resolution of object IDs from owned objects. Instead of manually providing object IDs as parameters, the ChainReader can dynamically fetch them from the blockchain using pointer tags that reference specific objects owned by a contract.

## How it Works

When a function parameter has a `PointerTag` configured, the ChainReader:

1. **Parses the pointer tag** format: `_::module::ObjectName::fieldName`
2. **Fetches owned objects** from the specified package on the Sui blockchain from the same contract being called
3. **Matches objects** based on the module and object type
4. **Extracts the field value** and automatically populates the parameter

## Configuration Format

```go
type SuiFunctionParam struct {
    Type       string   // Parameter type (e.g., "object_id")
    Name       string   // Parameter name 
    PointerTag *string  // Optional pointer tag for automatic resolution
    Required   bool     // Whether parameter is required
    // ... other fields
}
```

### PointerTag Syntax

The pointer tag follows this format:
```
_::module::ObjectType::fieldName
```

- `_` - Placeholder for package ID (automatically replaced)
- `module` - The Sui module name containing the object type
- `ObjectType` - The object/struct type to search for, this usually ends with the word "Pointer" in the contract
- `fieldName` - The field within the object to extract as the parameter value

## Usage Example

From the test configuration:

```go
pointerTag := "_::counter::CounterPointer::counter_id"

// Function configuration
"get_count_using_pointer": {
    Name:          "get_count_using_pointer", 
    SignerAddress: accountAddress,
    Params: []codec.SuiFunctionParam{
        {
            Type:       "object_id",
            Name:       "counter_id", 
            PointerTag: &pointerTag,
            Required:   true,
        },
    },
}
```

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

1. **Identifies pointer parameters** by checking for non-nil `PointerTag`
2. **Validates pointer tag format** (must have exactly 4 parts separated by `::`)
3. **Builds pointer queries** by grouping tags by module/object type
4. **Fetches owned objects** using `client.ReadOwnedObjects()`
5. **Matches and extracts values** from object fields
6. **Populates argument map** with resolved values


## Benefits

- **Automatic object resolution** - No need to manually track object IDs
- **Dynamic parameter population** - Objects are resolved at call time
- **Simplified API calls** - Reduces the complexity of function invocations
- **Type safety** - Automatic conversion to appropriate object types (`bind.Object{Id: value}` for `object_id` type)

## Limitations

- Only works with objects owned by the contract package
- Requires objects to exist and be accessible via `ReadOwnedObjects`
- Pointer tag format must be strictly followed
- Field names must match exactly between the tag and the actual object structure