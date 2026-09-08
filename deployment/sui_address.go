package deployment

import (
	"strings"

	suicodec "github.com/smartcontractkit/chainlink-sui/codec"
)

// Sui addresses are 32-byte hex values; case and leading zeroes are formatting only.

// SuiAddressesEqual compares valid Sui addresses canonically; other strings compare exactly.
func SuiAddressesEqual(a, b string) bool {
	if a == "" || b == "" {
		return a == b
	}
	// Trim only for address parsing; opaque values use the original strings below.
	canonicalA, errA := suicodec.ToSuiAddress(strings.TrimSpace(a))
	canonicalB, errB := suicodec.ToSuiAddress(strings.TrimSpace(b))
	if errA == nil && errB == nil {
		return canonicalA == canonicalB
	}
	return a == b
}

// QualifierContainsSuiAddress detects an address-derived qualifier. It supports padded and
// unpadded 0x-prefixed forms; bare hex is intentionally not matched because writers used 0x.
func QualifierContainsSuiAddress(qualifier, address string) bool {
	if qualifier == "" || address == "" {
		return false
	}
	q := strings.ToLower(qualifier)
	canonical, ok := canonicalSuiAddress(address)
	if !ok {
		return strings.Contains(q, strings.ToLower(address))
	}
	if strings.Contains(q, canonical) {
		return true
	}
	unpadded := "0x" + strings.TrimLeft(canonical[2:], "0")
	return unpadded != "0x" && strings.Contains(q, unpadded)
}

// canonicalSuiAddress returns the canonical form of a valid Sui address.
func canonicalSuiAddress(address string) (string, bool) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", false
	}
	canonical, err := suicodec.ToSuiAddress(address)
	return canonical, err == nil
}
