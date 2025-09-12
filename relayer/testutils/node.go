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

func RandomPortNumberGenerator() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	port, err := listener.Addr().(*net.TCPAddr).Port, nil
	if err != nil {
		return 0, err
	}
	listener.Close()
	return port, nil
}

type SuiNodeEnvironment struct {
	Cmd       *exec.Cmd
	LocalUrl  string
	FaucetUrl string
}

// StartSuiNode starts a local Sui node using Docker
func StartSuiNode(nodeType NodeEnvType) (*SuiNodeEnvironment, error) {
	var cmd *exec.Cmd

	port, err := RandomPortNumberGenerator()
	if err != nil {
		return nil, err
	}
	faucetPort, err := RandomPortNumberGenerator()
	if err != nil {
		return nil, err
	}

	portStr := fmt.Sprintf("%d", port)
	faucetPortStr := fmt.Sprintf("%d", faucetPort)

	switch nodeType {
	case Docker:
		// Check if the container is already running
		cmd = exec.Command("docker", "ps", "-q", "-f", "name=sui-local")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		localUrl := fmt.Sprintf("http://localhost:%s", portStr)
		faucetUrl := fmt.Sprintf("http://localhost:%s/gas", faucetPortStr)

		// If the container is already running, return
		if len(strings.TrimSpace(string(output))) > 0 {
			return &SuiNodeEnvironment{Cmd: cmd, LocalUrl: localUrl, FaucetUrl: faucetUrl}, nil
		}

		// Start the container
		cmd = exec.Command("docker", "run", "--rm", "-d", "--name", "sui-local", "-p", "9000:9000", "mysten/sui-node:devnet")
		err = cmd.Run()
		if err != nil {
			return nil, err
		}
	case CLI:
		// Start the local sui node
		cmd = exec.Command("sui", "start",
			"--force-regenesis",
			"--fullnode-rpc-port="+portStr,
			"--with-faucet="+faucetPortStr,
		)
		err = cmd.Start()
		if err != nil {
			return nil, err
		}
	}

	localUrl := fmt.Sprintf("http://localhost:%s", portStr)
	faucetUrl := fmt.Sprintf("http://localhost:%s/gas", faucetPortStr)

	// Wait for the node to start
	const defaultDelay = 10 * time.Second
	const backoffDelay = 100 * time.Millisecond
	err = waitForConnection(localUrl, defaultDelay, backoffDelay)
	if err != nil {
		return nil, err
	}
	// wait for Faucet to be available
	err = waitForConnection(faucetUrl, defaultDelay, backoffDelay)
	if err != nil {
		return nil, err
	}

	return &SuiNodeEnvironment{Cmd: cmd, LocalUrl: localUrl, FaucetUrl: faucetUrl}, nil
}

func waitForConnection(url string, timeout time.Duration, backoffDelay time.Duration) error {
	fmt.Println("Waiting for connection", "url", url)
	// Parse the URL to extract host and port
	parsedURL, err := netUrl.Parse(url)
	if err != nil {
		return fmt.Errorf("invalid URL %s: %w", url, err)
	}

	host := parsedURL.Host
	if host == "" {
		// Handle case where URL might just be "host:port"
		host = parsedURL.Path
	}

	// Add default port if missing
	if !strings.Contains(host, ":") {
		if parsedURL.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	// Use exponential backoff for retries
	deadline := time.Now().Add(timeout)

	for attempt := 1; time.Now().Before(deadline); attempt++ {
		conn, err := net.DialTimeout("tcp", host, 1*time.Second)
		if err == nil {
			conn.Close()
			return nil
		}

		// Calculate next backoff with exponential increase
		nextBackoff := backoffDelay * time.Duration(attempt)

		// Don't sleep longer than remaining time
		remainingTime := time.Until(deadline)
		if remainingTime < nextBackoff {
			nextBackoff = remainingTime
		}

		if remainingTime <= 0 {
			break
		}

		time.Sleep(nextBackoff)
	}

	return fmt.Errorf("timed out waiting for %s after %s", host, timeout)
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
func FundWithFaucet(log logger.Logger, recipient string, faucetUrl string) error {
	log.Infow("Funding account with test tokens", "address", recipient, "faucetUrl", faucetUrl)

	body := models.FaucetRequest{
		FixedAmountRequest: &models.FaucetFixedAmountRequest{
			Recipient: recipient,
		},
	}

	// Request funds from faucet
	faucetRequestErr := faucetRequest(faucetUrl, body, map[string]string{})
	if faucetRequestErr != nil {
		log.Errorw("Failed to request funds from faucet", "err", faucetRequestErr, "faucetUrl", faucetUrl)
		return faucetRequestErr
	}

	return nil
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
