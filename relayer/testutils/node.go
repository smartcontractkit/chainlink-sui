package testutils

import (
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"os/exec"
	"strings"
	"time"
)

// NodeEnvType represents the type of Sui node environment to run the localnet with
type NodeEnvType int

const (
	// Docker represents running a Sui node in a Docker container
	Docker NodeEnvType = iota
	// CLI represents running a Sui node via the Sui CLI
	CLI
)

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
		cmd := exec.Command("sui", "start", "--with-faucet")
		err := cmd.Start()
		if err != nil {
			return err
		}
	}

	// Wait for the node to start
	time.Sleep(5 * time.Second)
	return nil
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
func FundWithFaucet(logger logger.Logger, network string, recipient string) error {
	// In a real implementation, this would call the Sui faucet API
	// For simplicity in testing, we'll just log that we're "funding" the account
	logger.Infow("Funding account with test tokens", "address", recipient)

	faucetHost, err := sui.GetFaucetHost(network)
	if err != nil {
		logger.Errorw("GetFaucetHost err:", err)
		return err
	}

	logger.Infow("Faucet Host found", "host", faucetHost)

	header := map[string]string{}
	err = sui.RequestSuiFromFaucet(faucetHost, recipient, header)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("Request DevNet Sui From Faucet success")
	return nil
}
