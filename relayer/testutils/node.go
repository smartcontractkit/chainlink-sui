package testutils

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// NodeEnvType represents the type of Sui node environment to run the localnet with
type NodeEnvType int

const (
	// Docker represents running a Sui node in a Docker container
	Docker NodeEnvType = iota
	// CLI represents running a Sui node via the Sui CLI
	CLI
)

const SuiLocalEndpoint = "127.0.0.1:9000"

// StartSuiNode starts a local Sui node using Docker
func StartSuiNode(nodeType NodeEnvType) error {
	switch nodeType {
	case Docker:
		// Check if the container is already running
		cmd := exec.Command("docker", "ps", "-q", "-f", "name=sui-local")
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		// If the container is already running, return
		if len(strings.TrimSpace(string(output))) > 0 {
			return nil
		}

		// Start the container
		cmd = exec.Command("docker", "run", "--rm", "-d", "--name", "sui-local", "-p", "9000:9000", "mysten/sui-node:devnet")
		err = cmd.Run()
		if err != nil {
			return err
		}
	case CLI:
		// Start the local sui node
		cmd := exec.Command("sui", "start", "--with-faucet", "--force-regenesis")
		err := cmd.Start()
		if err != nil {
			return err
		}
	}

	// Wait for the node to start
	const defaultDelay = 10 * time.Second
	err := waitForConnection(SuiLocalEndpoint, defaultDelay)
	if err != nil {
		return err
	}

	// TODO: This creates a default local address, e.g. for publishing contracts. We should not need this. Any interaction with the chain should be through a controlled wallet
	cmd := exec.Command("sui", "client", "new-address", "ed25519")
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = waitForConnection(constant.FaucetLocalnetEndpoint, defaultDelay)
	if err != nil {
		return err
	}

	return nil
}

func waitForConnection(url string, timeout time.Duration) error {
	// Remove "http://" prefix if present
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", url, 1*time.Second)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for %s", url)
}

// FundWithFaucet Funds a Sui account with test tokens using the Sui faucet API.
// NOTE: The Sui faucet must be already running.
//
// It logs the funding details and attempts to request tokens from the faucet.
// Parameters:
// - logger: A logger instance used to log the funding process.
// - network: The network from which the faucet tokens are requested. Use "sui/constant" (e.g., "constant.SuiLocalnet").
// - recipient: The recipient's address to fund.
// Returns an error if the faucet request fails or if there is an issue determining the faucet host.
func FundWithFaucet(log logger.Logger, network string, recipient string) error {
	// In a real implementation, this would call the Sui faucet API
	// For simplicity in testing, we'll just log that we're "funding" the account
	log.Infow("Funding account with test tokens", "address", recipient)

	faucetHost, err := sui.GetFaucetHost(network)
	if err != nil {
		log.Errorw("GetFaucetHost err:", err)
		return err
	}

	log.Infow("Faucet Host found", "host", faucetHost)

	header := map[string]string{}
	err = sui.RequestSuiFromFaucet(faucetHost, recipient, header)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Request DevNet Sui From Faucet success")

	return nil
}
