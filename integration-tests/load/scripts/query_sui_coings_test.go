package scripts

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
)

// TestQuerySuiCoins queries all SUI coins for the configured signer.
//
// Usage: go test -v -run TestQuerySuiCoins ./scripts/
func TestQuerySuiCoins(t *testing.T) {
	cwd, _ := os.Getwd()
	if err := os.Chdir(".."); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(cwd) })

	ctx := context.Background()

	// Only need env name + private key — no run config needed.
	envName := "testnet"
	_, _, _, suiGrpcToken, err := config.LoadEnvConfig(envName)
	if err != nil {
		t.Fatalf("LoadEnvConfig: %v", err)
	}

	networks, err := config.LoadNetworks(envName)
	if err != nil {
		t.Fatalf("LoadNetworks: %v", err)
	}

	// Find the Sui network by its chain selector.
	suiNetwork, err := config.FindNetworkBySelector(networks, 9762610643973837292)
	if err != nil {
		t.Fatalf("FindNetworkBySelector: %v", err)
	}
	if len(suiNetwork.RPCs) == 0 {
		t.Fatal("no RPCs for Sui network")
	}

	rpc := suiNetwork.RPCs[0]
	grpcToken := suiGrpcToken
	if grpcToken == "" {
		grpcToken = rpc.GrpcToken
	}
	ptbClient, err := sui.NewSuiClient(t, rpc.HTTPURL, rpc.GrpcTarget, grpcToken)
	if err != nil {
		t.Fatalf("NewSuiClient: %v", err)
	}

	senderAddress := "0x3f6d6a9e3f7707485bf51c02a6bc6cb6e17dffe7f3e160b3c5520d55d1de8398"
	coins, err := ptbClient.QueryCoinsByAddress(ctx, senderAddress, "0x2::coin::Coin")
	if err != nil {
		t.Fatalf("query SUI coins for sender: %v", err)
	}
	if len(coins) == 0 {
		t.Fatalf("no SUI coins found for sender %s", senderAddress)
	}
	for _, c := range coins {
		fmt.Printf("Coin Type: %s, Coin ID: %s, Balance: %d\n", c.GetObjectType(), c.GetObjectId(), c.GetBalance())
	}
}
