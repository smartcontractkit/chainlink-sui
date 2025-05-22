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

This is an example that publishes a package, and interacts with one of its contracts.

```go
import (
  "github.com/pattonkan/sui-go/suiclient"
  "github.com/smartcontractkit/chainlink-sui/bindings/bind"
  rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

func PublishAndIncrementCounter(client *suiclient.ClientImpl, signer rel.SuiSigner) {
  ctx := context.Background()

  // Deploys the Test package using the Package binding
  testPackage, tx, err := PublishTest(ctx, bind.TxOpts{}, signer, *client)

  // Get the object created in the `init` function
  counterObjectId, err := bind.FindObjectIdFromPublishTx(tx, "counter", "Counter")

  // Get the Counter contract binding
  counter := testPackage.Counter()

  // We construct the Increment method pasing the `increment` needed params
  increment := counter.Increment(counterObjectId)

  // We decide what to do with the method
  ptb, err := increment.Build(ctx) // we can extract the PTB
  tx, err := increment.Execute(ctx, bind.TxOpts{}, signer, *client) // or we can send the transaction
}
```
