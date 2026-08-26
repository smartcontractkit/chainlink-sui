package client

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink-sui/relayer/client/suigrpcconn"
)

// chainIdentifierHexLen is the length of the short chain identifier that sui_getChainIdentifier
// returns and that Move.toml [environments] / Published.toml chain-id expect: the first 8 hex
// characters (4 bytes) of the genesis checkpoint digest.
const chainIdentifierHexLen = 8

// GetChainIdentifier returns the chain identifier of the node reachable at rpcURL, over gRPC.
//
// It calls the LedgerService GetServiceInfo RPC and reads its chain_id field, which (per the Sui
// proto) is the digest of the genesis checkpoint, base58-encoded. sui_getChainIdentifier and the
// Sui CLI instead surface the first 8 hex chars of that digest, which is the form Move.toml
// [environments] and Published.toml chain-id require. GetChainIdentifier decodes and truncates here
// so callers get the JSON-RPC-equivalent value without speaking JSON-RPC.
//
// rpcURL is a JSON-RPC-style URL (http(s)://host:port) or a bare gRPC target (host:port). The
// scheme is stripped when present to obtain the gRPC dial target, and TLS is selected by
// grpcTargetUsesTLS (local/non-443 plaintext, :443 TLS), matching
// how NewPTBClientFromConfig dials the same node. No auth token is injected, so this is suitable for
// local nodes; authenticated public endpoints should go through a fully-configured PTBClient instead.
func GetChainIdentifier(ctx context.Context, rpcURL string) (string, error) {
	target, useTLS, err := grpcTargetFromRPCURL(rpcURL)
	if err != nil {
		return "", err
	}

	conn := suigrpcconn.NewConnectionWithAuth(target, "", useTLS)
	defer func() { _ = conn.Close() }()

	if connectErr := conn.Connect(ctx); connectErr != nil {
		return "", fmt.Errorf("grpc connect to %s: %w", target, connectErr)
	}

	ledger, err := conn.LedgerService(ctx)
	if err != nil {
		return "", err
	}

	// GetServiceInfo is the same RPC that PTBClient.GetCheckpointAvailability uses; its chain_id is
	// the genesis checkpoint digest (base58). Convert it to the 8-hex-char identifier the CLI and
	// Move.toml use so this is a drop-in replacement for sui_getChainIdentifier.
	resp, err := ledger.GetServiceInfo(ctx, &suirpcv2.GetServiceInfoRequest{})
	if err != nil {
		return "", fmt.Errorf("GetServiceInfo on %s: %w", target, err)
	}

	id, err := chainIdentifierFromDigest(resp.GetChainId())
	if err != nil {
		return "", fmt.Errorf("deriving chain identifier from GetServiceInfo chain_id: %w", err)
	}
	return id, nil
}

// chainIdentifierFromDigest converts the genesis checkpoint digest (base58-encoded, as returned by
// GetServiceInfo.chain_id) into the 8-hex-char chain identifier that sui_getChainIdentifier returns.
func chainIdentifierFromDigest(b58Digest string) (string, error) {
	b58Digest = strings.TrimSpace(b58Digest)
	if b58Digest == "" {
		return "", errors.New("empty chain id")
	}
	digest, err := base58.Decode(b58Digest)
	if err != nil {
		return "", fmt.Errorf("base58 decode chain id: %w", err)
	}
	full := hex.EncodeToString(digest)
	if len(full) < chainIdentifierHexLen {
		return "", fmt.Errorf("decoded chain id too short (%d bytes) for %q", len(digest), b58Digest)
	}
	return full[:chainIdentifierHexLen], nil
}

// grpcTargetFromRPCURL converts an RPC endpoint into a gRPC dial target and reports whether TLS
// should be used. It accepts both a JSON-RPC URL (http(s)://host:port) and a bare gRPC target
// (host:port, as used by LocalGrpcURL and a scheme-less SUI_RPC_URL); the scheme is stripped when
// present. It reuses grpcTargetUsesTLS so TLS selection stays consistent with
// NewPTBClientFromConfig (e.g. local 127.0.0.1:9000 stays plaintext, public :443 uses TLS).
func grpcTargetFromRPCURL(rpcURL string) (target string, useTLS bool, err error) {
	rpcURL = strings.TrimSpace(rpcURL)
	if rpcURL == "" {
		return "", false, errors.New("empty rpc url")
	}
	if i := strings.Index(rpcURL, "://"); i >= 0 {
		target = rpcURL[i+3:]
	} else {
		target = rpcURL
	}
	if target == "" {
		return "", false, fmt.Errorf("rpc url has an empty host: %q", rpcURL)
	}
	return target, grpcTargetUsesTLS(target), nil
}
