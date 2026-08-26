package client

import (
	"encoding/hex"
	"testing"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChainIdentifierFromDigest verifies the genesis-digest -> 8-hex-char chain identifier
// conversion that makes gRPC GetServiceInfo.chain_id a drop-in for sui_getChainIdentifier.
// It is independent of any live node: a known digest is base58-encoded, then decoded back and
// truncated by chainIdentifierFromDigest, which must yield the first 8 hex chars of the digest.
func TestChainIdentifierFromDigest(t *testing.T) {
	// 32-byte genesis checkpoint digest whose first 4 bytes are 78 25 f4 0f.
	digest, err := hex.DecodeString("7825f40f" + "a1b2c3d4e5f60718293a4b5c6d7e8f90" + "00112233445566778899aabbccddeeff")
	require.NoError(t, err)
	require.Len(t, digest, 32)

	encoded := base58.Encode(digest)
	expected := hex.EncodeToString(digest)[:chainIdentifierHexLen] // "7825f40f"

	got, err := chainIdentifierFromDigest(encoded)
	require.NoError(t, err)
	assert.Equal(t, "7825f40f", got)
	assert.Equal(t, expected, got)

	// Round-trip stability: re-encoding the returned short id is not expected, but decoding the
	// full digest must always reproduce the same prefix regardless of base58 framing.
	for _, tc := range []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"valid", encoded, false},
		{"valid with surrounding whitespace", "  " + encoded + "  ", false},
		{"not base58", "this is not base58!!!", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := chainIdentifierFromDigest(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGrpcTargetFromRPCURL verifies RPC-endpoint -> gRPC dial target + TLS derivation without
// touching the network. It covers JSON-RPC URLs (with scheme) and bare gRPC targets (no scheme,
// as used by LocalGrpcURL / a scheme-less SUI_RPC_URL).
func TestGrpcTargetFromRPCURL(t *testing.T) {
	tests := []struct {
		name    string
		rpcURL  string
		target  string
		useTLS  bool
		wantErr bool
	}{
		{"http local", "http://127.0.0.1:9000", "127.0.0.1:9000", false, false},
		{"https local", "https://127.0.0.1:9000", "127.0.0.1:9000", false, false},
		{"bare local (LocalGrpcURL)", "127.0.0.1:9000", "127.0.0.1:9000", false, false},
		{"localhost local", "http://localhost:9000", "localhost:9000", false, false},
		{"https public 443", "https://sui.example.com:443", "sui.example.com:443", true, false},
		{"http public non-443", "http://sui.example.com:8080", "sui.example.com:8080", false, false},
		{"bare public 443", "sui.example.com:443", "sui.example.com:443", true, false},
		{"empty", "", "", false, true},
		{"scheme only", "http://", "", false, true},
		{"whitespace", "   ", "", false, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			target, useTLS, err := grpcTargetFromRPCURL(tc.rpcURL)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.target, target)
			assert.Equal(t, tc.useTLS, useTLS)
		})
	}
}
