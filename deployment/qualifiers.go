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
	holder = strings.TrimSpace(holder)
	if holder == "" {
		return TokenQualifier(symbol) + "-"
	}
	if canonical, ok := canonicalSuiAddress(holder); ok {
		holder = canonical
	} else {
		// Preserve legacy handling for non-address values.
		holder = strings.ToLower(holder)
		if !strings.HasPrefix(holder, "0x") {
			holder = "0x" + holder
		}
	}
	return TokenQualifier(symbol) + "-" + holder
}
