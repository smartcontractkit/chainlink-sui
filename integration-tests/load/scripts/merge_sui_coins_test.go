package scripts

import (
	"context"
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
)

// TestMergeSuiCoins merges all SUI coins for the configured signer into one.
//
// Usage: go test -v -run TestMergeSuiCoins ./scripts/
func TestMergeSuiCoins(t *testing.T) {
	cwd, _ := os.Getwd()
	if err := os.Chdir(".."); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(cwd) })

	ctx := context.Background()

	// Only need env name + private key — no run config needed.
	envName := "testnet"
	suiPrivKey, _, _, err := config.LoadEnvConfig(envName)
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

	ptbClient, err := sui.NewSuiClient(t, suiNetwork.RPCs[0].HTTPURL)
	if err != nil {
		t.Fatalf("NewSuiClient: %v", err)
	}

	signer, senderAddress, err := sui.NewSuiSigner(suiPrivKey)
	if err != nil {
		t.Fatalf("NewSuiSigner: %v", err)
	}

	t.Logf("Signer address: %s", senderAddress)

	mergedID, err := sui.MergeAllSuiCoins(ctx, ptbClient, signer, senderAddress)
	if err != nil {
		t.Fatalf("MergeAllSuiCoins: %v", err)
	}

	t.Logf("All SUI coins merged into: %s", mergedID)
}
