package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"

	"github.com/smartcontractkit/chainlink-sui/relayer/client/suigrpcconn"
)

// GetChainIdentifier returns the chain identifier of the node reachable at rpcURL, over gRPC.
//
// It calls the LedgerService GetServiceInfo RPC and reads its chain_id field, which is the digest
// of the genesis checkpoint — the same clean chain identifier that sui_getChainIdentifier returns
// and that Move.toml [environments] expects. Unlike sui_getChainIdentifier (JSON-RPC), this uses the
// gRPC transport the relayer already speaks to the node, so callers no longer need a JSON-RPC client.
//
// rpcURL is a JSON-RPC-style URL (http(s)://host:port). The scheme is stripped to obtain the gRPC
// dial target and TLS is selected by grpcTargetUsesTLS (local/non-443 plaintext, :443 TLS), matching
// how NewPTBClientFromConfig dials the same node. No auth token is injected, so this is suitable for
// local nodes; authenticated public endpoints should go through a fully-configured PTBClient instead.
func GetChainIdentifier(ctx context.Context, rpcURL string) (string, error) {
	target, useTLS, err := grpcTargetFromRPCURL(rpcURL)
	if err != nil {
		return "", err
	}

	conn := suigrpcconn.NewConnectionWithAuth(target, "", useTLS)
	defer func() { _ = conn.Close() }()

	if err := conn.Connect(ctx); err != nil {
		return "", fmt.Errorf("grpc connect to %s: %w", target, err)
	}

	ledger, err := conn.LedgerService(ctx)
	if err != nil {
		return "", err
	}

	// GetServiceInfo is the same RPC that PTBClient.GetCheckpointAvailability uses; its chain_id is
	// the genesis-checkpoint digest, i.e. the chain identifier.
	resp, err := ledger.GetServiceInfo(ctx, &suirpcv2.GetServiceInfoRequest{})
	if err != nil {
		return "", fmt.Errorf("GetServiceInfo on %s: %w", target, err)
	}

	id := strings.TrimSpace(resp.GetChainId())
	if id == "" {
		return "", errors.New("GetServiceInfo returned an empty chain id")
	}
	return id, nil
}

// grpcTargetFromRPCURL converts a JSON-RPC URL (http(s)://host:port) into a gRPC dial target and
// reports whether TLS should be used. It reuses grpcTargetUsesTLS so TLS selection stays consistent
// with NewPTBClientFromConfig (e.g. local 127.0.0.1:9000 stays plaintext, public :443 uses TLS).
func grpcTargetFromRPCURL(rpcURL string) (target string, useTLS bool, err error) {
	switch {
	case strings.HasPrefix(rpcURL, "https://"):
		target = strings.TrimPrefix(rpcURL, "https://")
	case strings.HasPrefix(rpcURL, "http://"):
		target = strings.TrimPrefix(rpcURL, "http://")
	default:
		return "", false, fmt.Errorf("rpc url must include an http(s) scheme: %q", rpcURL)
	}
	if target == "" {
		return "", false, fmt.Errorf("rpc url has an empty host: %q", rpcURL)
	}
	return target, grpcTargetUsesTLS(target), nil
}
