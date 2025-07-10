# Sui contract bindings

This package contains bindings for all Sui contracts. It is used to publish packages and interact with deployed contracts.

## Generate contract bindings

Inside the nix shell run the `bindings:generate` task. It will execute the [`generate_bindings`](../scripts/generate_bindings.sh) script.

```
(nix:nix-shell-env)$ task bindings:generate

task: [bindings:generate] sh ./scripts/generate_bindings.sh
Generating bindings for Move Sui contracts...

	##############################################################
	Generating Go bindings for: ./contracts/test/sources/counter.move
	##############################################################

  (...)

```

You'll see that some structs and functions aren't parseable. This should be fine as long as we don't use those specific bindings

The contract bindings will be stored under [./generated](./generated/). Unless we need a binding that `bindgen` can't generate automatically, we shouldn't need to manually change these files.

## Package bindings

Package bindings are responsible for publishing packages and can be used as entry point to get contract bindings. They can't be auto-generated, as publishing needs very custom changes per package.

Package bindings live in [./packages](./packages/). Each Move package should have a single package binding.

## Using bindings

### Execution Example

This is an example that publishes a package and interacts with one of its contracts. When invoking the contract interface functions directly, this creates a PTB with a single MoveCall.

```go
import (
  "github.com/block-vision/sui-go-sdk/sui"
  "github.com/smartcontractkit/chainlink-sui/bindings/bind"
  "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

func PublishAndIncrementCounter(client *sui.Client, signer utils.SuiSigner) {
  ctx := context.Background()

  opts := &bind.CallOpts{
    Signer: signer,
    WaitForExecution: true,
  }

  // Deploys the Test package using the Package binding
  testPackage, tx, err := PublishTest(ctx, opts, client)
  
  counter := testPackage.Counter()
  initTx, err := counter.Initialize(ctx, opts)
  
  // Find the created counter object
  var counterObjectId string
  var initialSharedVersion *uint64
  for _, change := range initTx.ObjectChanges {
    if change.Type == "created" && strings.Contains(change.ObjectType, "::counter::Counter") {
      counterObjectId = change.ObjectId
      // Get initial shared version if it's a shared object
      if change.Owner != nil {
        share := change.GetObjectOwnerShare()
        if share.InitialSharedVersion != nil {
          version := uint64(*share.InitialSharedVersion)
          initialSharedVersion = &version
        }
      }
    }
  }

  // Create object reference for the counter
  // Note that InitialSharedVersion could be omitted, in which case, it would automatically
  // be resolved via sui_getObject
  counterObject := bind.Object{
    Id:                   counterObjectId,
    InitialSharedVersion: initialSharedVersion,
  }

  // Increment the counter
  incrementTx, err := counter.Increment(ctx, opts, counterObject)
  ...
}
```

### Multi-command PTB Example

```go
import (
  "github.com/block-vision/sui-go-sdk/transaction"
)

func ExecuteWithPTB(client *sui.Client, signer utils.SuiSigner, counter ICounter) {
  // Create a PTB
  ptb := transaction.NewTransaction()
  
  // Get the function call encoder
  ptbEncoder := counter.Encoder()
  
  // Add multiple calls to the PTB
  moveCall1, err := ptbEncoder.Increment(ptb, counterObject)
  moveCall2, err := ptbEncoder.IncrementBy(ptb, counterObject, 5)
  
  // Execute the PTB and wait for the results
  opts := &bind.CallOpts{
    Signer: signer,
    WaitForExecution: true,
  }
  tx, err := bind.ExecutePTB(ctx, opts, client, ptb)
  ...
}
```
