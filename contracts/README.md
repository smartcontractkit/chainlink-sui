# Sui Move Contracts

This repository contains Sui Move contracts converted from Aptos Move contracts.

## Contracts

### Counter

A simple counter contract that demonstrates:
- Creating and sharing a Sui object
- Incrementing a counter by 1
- Incrementing a counter by a multiple of two values

### Echo

An echo contract that demonstrates:
- Creating and emitting different types of events in Sui
- Various view functions that return different data types

## Building and Testing

To build the contracts:

```bash
sui move build -d
```

> Note: The `-d` flag specifies that a development address (from `Move.toml`) should be used.

To test the contracts:

```bash
sui move test
```

## Publishing
do the following within `test` directory:
1. review this [page](https://docs.sui.io/references/cli/client) to create an environment, ideally devnet for testing,
and create an account. You can also use this [cheatsheet](https://docs.sui.io/references/cli/cheatsheet) as a reference. 
2. run `sui client faucet` to get some devnet SUI tokens. Verify the balance with `sui client gas`.
3. run `sui client objects` and you should see exactly 1 coin object, which is the tokens your just received.
4. run `sui client publish --gas-budget 50000000`
You should see lots of output. But look for the published objects:
```
│ Published Objects:                                                                                    │
│  ┌──                                                                                                  │
│  │ PackageID: 0x476dee4c8e4f8147b42b9034ec45ff40a837e61b211f944f6b1ad7d8d51aec9d                      │
│  │ Version: 1                                                                                         │
│  │ Digest: BRvoq9qJKygTG1La9P51RmVmRZmUht6MXKPX1ToJcsQJ                                               │
│  │ Modules: counter, echo                                                                             │
│  └──    
```
and created objects:
```
├───────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Created Objects:                                                                                      │
│  ┌──                                                                                                  │
│  │ ObjectID: 0x53326cd28679d3607d87a91c29a41f10ee34763ee12421b159e08183f6df8c72                       │
│  │ Sender: 0xf05692ca14ad4ee63c71fab652cd0544f1fa3eca048f36e1f75421147158d6ec                         │
│  │ Owner: Account Address ( 0xf05692ca14ad4ee63c71fab652cd0544f1fa3eca048f36e1f75421147158d6ec )      │
│  │ ObjectType: 0x476dee4c8e4f8147b42b9034ec45ff40a837e61b211f944f6b1ad7d8d51aec9d::counter::AdminCap  │
│  │ Version: 151                                                                                       │
│  │ Digest: 8Yhg5Gp9Rv6S2SjBs5gVEWkVdg7PMAzxbn2Fpc7HGzY                                                │
│  └──                                                                                                  │
│  ┌──                                                                                                  │
│  │ ObjectID: 0x608ea95bbcc36167efa841b5f65a092519bdf39ca7a5aac49bf3df081a6eb9d6                       │
│  │ Sender: 0xf05692ca14ad4ee63c71fab652cd0544f1fa3eca048f36e1f75421147158d6ec                         │
│  │ Owner: Account Address ( 0xf05692ca14ad4ee63c71fab652cd0544f1fa3eca048f36e1f75421147158d6ec )      │
│  │ ObjectType: 0x2::package::UpgradeCap                                                               │
│  │ Version: 151                                                                                       │
│  │ Digest: 4GFg7wbb7Gn2hwhNsrzzmdeBoCH9G1shMDFzRQ1777eg                                               │
│  └──                                                                                                  │
│  ┌──                                                                                                  │
│  │ ObjectID: 0xc0c6faa0ef6b6991e3ba1574ed5c59333bb31e01f402492e898aac4623c32460                       │
│  │ Sender: 0xf05692ca14ad4ee63c71fab652cd0544f1fa3eca048f36e1f75421147158d6ec                         │
│  │ Owner: Shared( 151 )                                                                               │
│  │ ObjectType: 0x476dee4c8e4f8147b42b9034ec45ff40a837e61b211f944f6b1ad7d8d51aec9d::counter::Counter   │
│  │ Version: 151                                                                                       │
│  │ Digest: Er4L9DfhTKyd3cgp3U3MUtBPqnhkh6uVPcoZHkRk5PJc                                               │
│  └──    
```
5. run `sui client objects` again to see your new objects: `AdminCap` and `UpgradeCap`. Note: the `Counter` object is shared.
6. run `sui client call --package package_id --module counter --function increment --args counter_object_id` to increment the counter by 1.
7. run `sui client call --package package_id --module counter --function get_count --args counter_object_id --dev-inspect` to get the count, which should be 1.

```
Execution Result
  Return values
    Sui TypeTag: SuiTypeTag("u64")
    Bytes: [7, 0, 0, 0, 0, 0, 0, 0]
```

8. run `sui client call  --package 0x476dee4c8e4f8147b42b9034ec45ff40a837e61b211f944f6b1ad7d8d51aec9d --module counter --function increment_by_two --args 0x53326cd28679d3607d87a91c29a41f10ee34763ee12421b159e08183f6df8c72 0xc0c6faa0ef6b6991e3ba1574ed5c59333bb31e01f402492e898aac4623c32460` to pass in the admin cap object to increase the counter by 2.

**Note: a move call will always cost gas no matter it changes on-chain states or not. However, by adding --dev-inspect flag, we don't need to pay gas for a "view" function**

## Key Differences between Aptos and Sui Move

The converted contracts showcase several key differences between Aptos and Sui Move:

1. **Object Model**: 
   - Sui uses an object-centric model with `UID` as the first field in key structs
   - Objects are explicitly created and shared/transferred

2. **Events**:
   - In Sui, events are emitted directly with `event::emit()`
   - No need for event handles like in Aptos

3. **Transaction Context**:
   - Sui functions use `TxContext` instead of `&signer`
   - Object ownership is managed through this context

4. **Module Initialization**:
   - Sui uses `init` function instead of Aptos's `init_module`

5. **Global Storage**:
   - Sui doesn't use Aptos's `move_to`, `borrow_global`, etc.
   - Objects are passed directly to functions

## License

This project is available under the same license as the original contracts. 