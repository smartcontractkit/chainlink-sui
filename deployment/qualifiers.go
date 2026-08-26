package deployment

import "strings"

// ChainSingletonQualifier is used for types with one instance per chain.
const ChainSingletonQualifier = ""

// SupersededLabel marks a ref that no longer represents the active deployment.
const SupersededLabel = "superseded"

// TokenQualifier normalizes a token symbol for use as a datastore key.
func TokenQualifier(symbol string) string {
	return strings.Join(strings.Fields(symbol), "-")
}

// MinterCapQualifier qualifies a token minter capability by token and holder.
func MinterCapQualifier(symbol string, holder string) string {
	return TokenQualifier(symbol) + "-" + holder
}
