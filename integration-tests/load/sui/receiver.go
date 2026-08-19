package sui

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// ParseEVMReceiver32 parses a hex EVM receiver and returns a 32-byte left-padded value.
// Accepts either 20-byte (EVM address) or 32-byte already padded input.
func ParseEVMReceiver32(receiver string) ([]byte, error) {
	raw := strings.TrimPrefix(strings.TrimSpace(receiver), "0x")
	if raw == "" {
		return nil, fmt.Errorf("receiver is empty")
	}

	decoded, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("receiver is not valid hex: %w", err)
	}

	switch len(decoded) {
	case 20:
		padded := make([]byte, 32)
		copy(padded[12:], decoded)
		if isAllZero(decoded) {
			return nil, fmt.Errorf("receiver EVM address cannot be zero")
		}
		return padded, nil
	case 32:
		if isAllZero(decoded[12:]) {
			return nil, fmt.Errorf("receiver EVM address cannot be zero")
		}
		return decoded, nil
	default:
		return nil, fmt.Errorf("receiver must decode to 20 or 32 bytes, got %d", len(decoded))
	}
}

func isAllZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

// ResolveReceiverState resolves the CCIPReceiverState shared object ID from a receiver package.
//
// packageID is the Sui receiver package ID (hex-encoded, with or without 0x prefix).
// The object is found by querying the package's owned objects for the CCIPReceiverState pointer
// and returning its parent object ID (the actual shared CCIPReceiverState object).
func ResolveReceiverState(ctx context.Context, ptbClient *client.PTBClient, packageID string) (string, error) {
	cleanPkg := strings.TrimPrefix(strings.TrimSpace(packageID), "0x")
	if cleanPkg == "" {
		return "", fmt.Errorf("receiver package ID is empty")
	}

	if ptbClient == nil {
		return "", fmt.Errorf("ptbClient is required to resolve receiver state")
	}

	parentObjectID, err := ptbClient.GetParentObjectID(ctx, packageID, "ccip_dummy_receiver", "CCIPReceiverState")
	if err != nil {
		return "", fmt.Errorf("failed to resolve receiver state from package %s: %w", packageID, err)
	}

	return parentObjectID, nil
}

// SuiObjectIdToBytes32 converts a Sui object ID hex string to a 32-byte array, left-padded with zeros.
func SuiObjectIdToBytes32(objectID string) ([32]byte, error) {
	var result [32]byte
	clean := strings.TrimPrefix(strings.TrimSpace(objectID), "0x")
	if clean == "" {
		return result, fmt.Errorf("object ID is empty")
	}
	decoded := common.FromHex("0x" + clean)
	if len(decoded) == 0 || len(decoded) > 32 {
		return result, fmt.Errorf("object ID must decode to 1-32 bytes, got %d", len(decoded))
	}
	copy(result[32-len(decoded):], decoded)
	return result, nil
}
