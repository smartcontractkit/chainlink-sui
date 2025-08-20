# Keystore

The Keystore is responsible for securely holding keys used to sign transactions. It is used by the Transaction Manager to sign transactions.

The keystore implementation is provided by Chainlink Core. It must implemented the interface `core.Keystore` from `chainlink-common`.

```go
type Keystore interface {
	Accounts(ctx context.Context) (accounts []string, err error)
	// Sign returns data signed by account.
	// nil data can be used as a no-op to check for account existence.
	Sign(ctx context.Context, account string, data []byte) (signed []byte, err error)
}
```

We have a test implementation of the Keystore in the `relayer/testutils/keystore.go` file.
