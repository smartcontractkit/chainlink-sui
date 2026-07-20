package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	netUrl "net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
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

// StartSuiNode starts a local Sui node using Docker
func StartSuiNode(nodeType NodeEnvType) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	switch nodeType {
	case Docker:
		// Check if the container is already running
		cmd = exec.Command("docker", "ps", "-q", "-f", "name=sui-local")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		// If the container is already running, return
		if len(strings.TrimSpace(string(output))) > 0 {
			return cmd, nil
		}

		// Start the container
		cmd = exec.Command("docker", "run", "--rm", "-d", "--name", "sui-local", "-p", "9000:9000", "mysten/sui-node:devnet")
		err = cmd.Run()
		if err != nil {
			return nil, err
		}
	case CLI:
		// Start the local sui node
		cmd = exec.Command("sui", "start", "--with-faucet", "--force-regenesis")
		err := cmd.Start()
		if err != nil {
			return nil, err
		}
	}

	// Wait for the node to start
	const defaultDelay = 10 * time.Second
	const backoffDelay = 100 * time.Millisecond
	err := waitForConnection(LocalURL, defaultDelay, backoffDelay)
	if err != nil {
		return nil, err
	}
	// wait for Faucet to be available
	err = waitForConnection(LocalFaucetURL, defaultDelay, backoffDelay)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func waitForConnection(url string, timeout time.Duration, backoffDelay time.Duration) error {
	parsedURL, err := netUrl.Parse(url)
	if err != nil {
		return fmt.Errorf("invalid URL %s: %w", url, err)
	}
	if parsedURL.Host == "" {
		return fmt.Errorf("invalid URL %s: missing host", url)
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		d := net.Dialer{
			Timeout: 5 * time.Second,
		}

		conn, err := d.Dial("tcp", parsedURL.Host)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(backoffDelay)
	}

	return fmt.Errorf("timed out waiting for %s after %s", parsedURL.Host, timeout)
}

func GetFaucetHost(network string) string {
	switch network {
	default:
		return LocalURL
	}
}

// FundWithFaucet Funds a Sui account with test tokens using the Sui faucet API.
// NOTE: The Sui faucet must be already running.
//
// It logs the funding details and attempts to request tokens from the faucet.
// Parameters:
// - logger: A logger instance used to log the funding process.
// - network: The network from which the faucet tokens are requested. Use "sui/constant" (e.g., "SuiLocalnet").
// - recipient: The recipient's address to fund.
// Returns an error if the faucet request fails or if there is an issue determining the faucet host.
func FundWithFaucet(log logger.Logger, network string, recipient string) error {
	log.Infow("Funding account with test tokens", "address", recipient)

	faucetHost, err := sui.GetFaucetHost(network)
	if err != nil {
		return err
	}

	body := models.FaucetRequest{
		FixedAmountRequest: &models.FaucetFixedAmountRequest{
			Recipient: recipient,
		},
	}

	const (
		maxAttempts    = 5
		initialBackoff = 500 * time.Millisecond
		maxBackoff     = 5 * time.Second
	)

	var lastErr error
	backoff := initialBackoff

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		faucetRequestErr := faucetRequest(faucetHost, body, map[string]string{})
		if faucetRequestErr == nil {
			return nil
		}
		lastErr = faucetRequestErr
		log.Warnw("Faucet request failed, will retry", "attempt", attempt, "err", faucetRequestErr)

		if attempt < maxAttempts {
			time.Sleep(backoff)
			backoff = min(backoff*2, maxBackoff)
		}
	}

	log.Errorw("Failed to request funds from faucet after retries", "err", lastErr)
	return lastErr
}

func faucetRequest(faucetUrl string, body any, headers map[string]string) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, faucetUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("request faucet failed, statusCode: %d, err: %+v", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
