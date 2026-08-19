package sui

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/bech32"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// NewSuiClient creates a Sui PTB client.
// grpcTarget should be host:port. If empty, it is derived from rpcURL.
func NewSuiClient(t *testing.T, rpcURL string, grpcTarget string, grpcToken string) (*client.PTBClient, error) {
	normalizedTarget, err := normalizeGrpcTarget(rpcURL, grpcTarget)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize gRPC target: %w", err)
	}

	if strings.TrimSpace(grpcToken) == "" {
		// Keep backwards compatibility with public nodes where token is optional
		// while still enabling the gRPC path in PTBClient config.
		grpcToken = "unused"
	}

	cfg := client.PTBClientConfig{
		GrpcTarget:         normalizedTarget,
		GrpcToken:          grpcToken,
		TransactionTimeout: 60 * time.Second,
		DefaultRequestType: client.WaitForLocalExecution,
	}

	lggr := logger.Test(t)

	ptbClient, err := client.NewPTBClient(lggr, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Sui PTB client: %w", err)
	}

	return ptbClient, nil
}

func normalizeGrpcTarget(rpcURL string, grpcTarget string) (string, error) {
	target := strings.TrimSpace(grpcTarget)
	if target != "" {
		return target, nil
	}

	rpcURL = strings.TrimSpace(rpcURL)
	if rpcURL == "" {
		return "", fmt.Errorf("rpcURL is empty and grpcTarget is not set")
	}

	if strings.Contains(rpcURL, "://") {
		u, err := url.Parse(rpcURL)
		if err != nil {
			return "", fmt.Errorf("invalid rpcURL %q: %w", rpcURL, err)
		}
		host := u.Hostname()
		if host == "" {
			return "", fmt.Errorf("invalid rpcURL %q: missing hostname", rpcURL)
		}
		port := u.Port()
		if port == "" {
			switch u.Scheme {
			case "https":
				port = "443"
			case "http":
				port = "80"
			default:
				return "", fmt.Errorf("unsupported rpcURL scheme %q", u.Scheme)
			}
		}
		return net.JoinHostPort(host, port), nil
	}

	if _, _, err := net.SplitHostPort(rpcURL); err == nil {
		return rpcURL, nil
	}

	if strings.Contains(rpcURL, "/") {
		return "", fmt.Errorf("invalid grpc target %q: expected host:port", rpcURL)
	}

	// Bare host without explicit port defaults to TLS endpoint.
	return net.JoinHostPort(rpcURL, "443"), nil
}

// NewSuiSigner creates a Sui signer from a bech32-encoded private key.
// The key must be in suiprivkey1... format (Ed25519).
// Returns the signer and the sender's Sui address.
func NewSuiSigner(bech32PrivKey string) (bindutils.SuiSigner, string, error) {
	seed, err := hexFromSuiBech32PrivKey(bech32PrivKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode Sui private key: %w", err)
	}

	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode seed hex: %w", err)
	}

	if len(seedBytes) != 32 {
		return nil, "", fmt.Errorf("invalid seed length: expected 32 bytes, got %d", len(seedBytes))
	}

	privateKey := ed25519.NewKeyFromSeed(seedBytes)
	s := signer.NewPrivateKeySigner(privateKey)

	address, err := s.GetAddress()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get Sui address: %w", err)
	}

	return s, address, nil
}

// hexFromSuiBech32PrivKey decodes a bech32-encoded Sui private key
// (suiprivkey1...) to a hex-encoded 32-byte seed.
func hexFromSuiBech32PrivKey(bech string) (string, error) {
	hrp, data5, err := bech32.Decode(bech)
	if err != nil {
		return "", fmt.Errorf("failed to decode bech32: %w", err)
	}
	if hrp != "suiprivkey" {
		return "", fmt.Errorf("unexpected HRP: %s (expected suiprivkey)", hrp)
	}

	dataBytes, err := bech32.ConvertBits(data5, 5, 8, false)
	if err != nil {
		return "", fmt.Errorf("failed to convert bits: %w", err)
	}

	if len(dataBytes) != 33 {
		return "", fmt.Errorf("decoded privkey wrong length: %d bytes (expected 33)", len(dataBytes))
	}

	// dataBytes[0] is the flag byte (0x00 for Ed25519)
	seed := dataBytes[1:]
	if len(seed) != 32 {
		return "", fmt.Errorf("unexpected seed length: %d (expected 32)", len(seed))
	}

	return hex.EncodeToString(seed), nil
}
