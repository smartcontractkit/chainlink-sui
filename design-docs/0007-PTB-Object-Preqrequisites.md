# PTB Prerequisite Objects

## Overview

PrerequisiteObjects is an extension to the ChainWriter configuration that allows the PTB (Programmable Transaction Block) constructor to automatically fetch owned objects and populate PTB arguments before construction. This eliminates the need for clients to manually query and provide object IDs or object data.i

The owner can be any address in the network such as a Package or an EOA.

## Architecture Components

### 1. Configuration Structure (`config.go`)

The `PrerequisiteObject` struct is defined in the ChainWriter configuration:

```go
type PrerequisiteObject struct {
    OwnerId  *string  // Owner address of the object to fetch
    Name     string   // Identifier name for the prerequisite
    Tag      string   // Object type tag to match against
    SetKeys  bool     // Whether to extract object keys or just use object ID
}
```

### 2. Integration in ChainWriter (`chainwriter.go`)

PrerequisiteObjects are configured at the function level within `ChainWriterFunction`:

```go
type ChainWriterFunction struct {
    Name                string
    PublicKey          []byte
    PrerequisiteObjects []PrerequisiteObject  // Added field
    Params             []codec.SuiFunctionParam
    PTBCommands        []ChainWriterPTBCommand
}
```

### 3. Processing Logic (`ptb_constructor.go`)

The main processing happens in the `PTBConstructor.BuildPTBCommands()` method (in `relayer/chainwriter/ptb_constructor.go`):

```go
	// Attempt to fill args with pre-requisite object data
	err := p.FetchPrereqObjects(ctx, txnConfig.PrerequisiteObjects, &args, overrideToAddress)
	if err != nil {
		return nil, err
	}
```

## Two Use Cases

### Case A: Object ID Prerequisites (`SetKeys: false`)

**Purpose**: When you need the object ID itself as a PTB argument.

**Configuration Example**:
```go
PrerequisiteObjects: []chainwriter.PrerequisiteObject{
    {
        OwnerId: &accountAddress,
        Name:    "admin_cap_id", 
        Tag:     "counter::AdminCap",
        SetKeys: false,  // Only store the object ID
    },
}
```

**Behavior**: 
- Searches for objects owned by `accountAddress` (or any address set in the `OwnerId` field) matching type `counter::AdminCap` (or whatever is specified in `Tag`)
- Adds the object ID to the args map: `args["admin_cap_id"] = "0x1234..."`
- The key used in the args would be the `Name` field's value

### Case B: Object Contents Prerequisites (`SetKeys: true`)

**Purpose**: When you need the object's internal data fields as PTB arguments.

**Configuration Example**:
```go
PrerequisiteObjects: []chainwriter.PrerequisiteObject{
    {
        OwnerId: &accountAddress,
        Name:    "counter_id",  // Name less important here
        Tag:     "counter::CounterPointer", 
        SetKeys: true,  // Extract all object fields
    },
}
```

**Behavior**:
- Searches for objects owned by `accountAddress` matching type `counter::CounterPointer`
- The object found contains multiple keys, for example `field1`, `field2`, etc..
- Parses the object's JSON content and adds each field: `args["field1"] = value1, args["field2"] = value2, etc.`

## Implementation Details

### Object Fetching Process (`FetchPrereqObjects`) in `relayer/chainwriter/ptb_constructor.go`

```go
func (p *PTBConstructor) FetchPrereqObjects(ctx context.Context, prereqObjects []PrerequisiteObject, args *map[string]any, ownerFallback *string) error {
	for _, prereq := range prereqObjects {
		// set the owner fallback if the ownerId is not provided
		if prereq.OwnerId == nil {
			if ownerFallback == nil {
				return fmt.Errorf("ownerId or ownerFallback required for pre-requisite object %s", prereq.Name)
			}

			prereq.OwnerId = ownerFallback
		}
		// fetch owned objects
		ownedObjects, err := p.client.ReadOwnedObjects(ctx, *prereq.OwnerId, nil)
		if err != nil {
			return err
		}

		// check each returned object
		for _, ownedObject := range ownedObjects {
			// object tag matches
			if ownedObject.Data.Type != nil && strings.Contains(*ownedObject.Data.Type, prereq.Tag) {
				p.log.Debugw("Found pre-requisite object", "Object", ownedObject.Data, "Prereq", prereq)
				// object must be parsed and its keys added to the args map
				if prereq.SetKeys {
					// parse the object into a map
					parsedObject := map[string]any{}
					err := json.Unmarshal(ownedObject.Data.Content.Data.MoveObject.Fields, &parsedObject)
					if err != nil {
						return err
					}

					// add each key and value to the args map
					for key, value := range parsedObject {
						(*args)[key] = value
					}
				} else {
					// add the object id to the args map
					(*args)[prereq.Name] = ownedObject.Data.ObjectId.String()
				}
			}
		}
	}

	return nil
}
```

### Owner Address Fallback

The system supports an owner address fallback mechanism:
- If `OwnerId` is `nil` in the config, it uses the `toAddress` from the `SendTransaction` method in ChainWriter
- This allows dynamic owner resolution at runtime rather than needing it to be configured at compile time
- See `relayer/chainwriter/chainwriter_test.go` for a test case

## Test Examples

### Test Case A: Object ID Prerequisites

See `relayer/chainwriter/ptb_constructor_test.go`

```go
//nolint:paralleltest
	t.Run("Should fill a valid prerequisite object ID in CW config", func(t *testing.T) {
		// we only pass the counter ID as the other object ID (admin cap) is populated by the pre-requisites
		args := map[string]any{
			"counter_id": counterObjectId,
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "get_count_with_object_id_prereq", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, publicKeyBytes, ptb)
		prettyPrintDebug(log, ptbResult)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
	})
```

### Test Case B: Object Contents Prerequisites

See `relayer/chainwriter/ptb_constructor_test.go`

```go
//nolint:paralleltest
	t.Run("Should fill a valid prerequisite object keys in CW config", func(t *testing.T) {
		// pass no args as it should be populated by the pre-requisites
		args := map[string]any{}

		ptb, err := constructor.BuildPTBCommands(ctx, "counter", "get_count_with_object_keys_prereq", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, publicKeyBytes, ptb)
		prettyPrintDebug(log, ptbResult)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
	})
```

## Execution Flow

1. **Configuration**: PrerequisiteObjects are defined in the ChainWriter config
2. **Transaction Submission**: Client calls `SubmitTransaction()` with minimal args
3. **PTB Building**: `BuildPTBCommands()` is called, which:
   - First calls `FetchPrereqObjects()` to populate missing arguments
   - Then processes PTB commands with the enriched args map
4. **Object Fetching**: `FetchPrereqObjects()` queries owned objects and populates args
5. **PTB Construction**: Normal PTB building proceeds with the enriched arguments

## Benefits

1. **Simplified Client Interface**: Clients don't need to manually query object IDs
2. **Dynamic Object Resolution**: Objects are resolved at execution time
3. **Flexible Data Extraction**: Support for both object IDs and object contents
4. **Owner Address Flexibility**: Support for fallback owner addresses

This system reduces amount of code written in the ChainWriter handler and relies on the correct configuration to ensure that prerequisites are fetched prior to PTB construction.
