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