package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	netUrl "net/url"
	"os"
	"os/exec"
	"strconv"
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

// StartSuiNode starts a local Sui node using Docker or the Sui CLI.
//
// For the CLI path it retries `sui start` until the fullnode RPC and faucet are
// both reachable. Retrying is necessary because Sui 1.75.x removed the post-launch
// grace period, so the faucet can race the fullnode loading genesis state and
// `sui start` exits with "No address found with sufficient coins" before the ports
// come up. Each attempt detects early process exit so a failed attempt is retried
// in ~1s instead of waiting the full per-attempt timeout.
func StartSuiNode(nodeType NodeEnvType) (*exec.Cmd, error) {
	switch nodeType {
	case Docker:
		return startDockerNode()
	case CLI:
		return startCliNodeWithRetry()
	default:
		return nil, fmt.Errorf("unknown node type: %v", nodeType)
	}
}

// startDockerNode starts a Sui node in a Docker container and waits for its RPC
// and faucet ports. The container persists across test runs, so no retry is used.
func startDockerNode() (*exec.Cmd, error) {
	cmd := exec.CommandContext(context.Background(), "docker", "ps", "-q", "-f", "name=sui-local")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(output))) > 0 {
		return cmd, nil
	}
	cmd = exec.CommandContext(context.Background(), "docker", "run", "--rm", "-d", "--name", "sui-local", "-p", "9000:9000", "mysten/sui-node:devnet")
	if err = cmd.Run(); err != nil {
		return nil, err
	}
	delay := startTimeoutFromEnv()
	const backoffDelay = 1000 * time.Millisecond
	if err := waitForConnection(LocalURL, delay, backoffDelay); err != nil {
		return nil, err
	}
	if err := waitForConnection(LocalFaucetURL, delay, backoffDelay); err != nil {
		return nil, err
	}
	return cmd, nil
}

// suiStartRetryDeadline bounds how long StartSuiNode retries `sui start` before
// giving up. Each attempt waits for the ports plus a stabilization window, so
// this allows several fresh starts if `sui start` keeps exiting from the
// validator health-check race.
const suiStartRetryDeadline = 3 * time.Minute

// startCliNodeWithRetry repeatedly starts `sui start --with-faucet --force-regenesis`
// until the node and faucet are reachable, or suiStartRetryDeadline elapses.
func startCliNodeWithRetry() (*exec.Cmd, error) {
	perAttempt := startTimeoutFromEnv()
	deadline := time.Now().Add(suiStartRetryDeadline)
	for attempt := 1; time.Now().Before(deadline); attempt++ {
		cmd, ok, ferr := startCliNodeOnce(perAttempt)
		if ferr != nil {
			return nil, ferr
		}
		if ok {
			return cmd, nil
		}
		fmt.Fprintf(os.Stderr, "[StartSuiNode] attempt %d did not come up; retrying\n", attempt)
		time.Sleep(200 * time.Millisecond)
	}
	return nil, fmt.Errorf("sui start did not become reachable within %s; see prior `sui start` output tails (likely the Sui 1.75.x faucet/fullnode startup race)", suiStartRetryDeadline)
}

// startCliNodeOnce starts `sui start` once, captures its output, and waits for the
// fullnode RPC and faucet to accept connections. It returns (cmd, true, nil) on
// success, (nil, false, nil) if the process exited or the ports were not reached
// in time (retryable), or (nil, false, err) on a fatal setup error.
func startCliNodeOnce(perAttempt time.Duration) (*exec.Cmd, bool, error) {
	logFile, err := os.CreateTemp("", "sui-start-*.log")
	if err != nil {
		return nil, false, fmt.Errorf("create sui start log: %w", err)
	}
	suiLogPath := logFile.Name()
	fmt.Fprintf(os.Stderr, "[StartSuiNode] capturing `sui start` output: %s\n", suiLogPath)

	cmd := exec.CommandContext(context.Background(), "sui", "start", "--with-faucet", "--force-regenesis")
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if startErr := cmd.Start(); startErr != nil {
		logFile.Close()
		return nil, false, fmt.Errorf("start sui: %w", startErr)
	}
	// The child holds its own dup'd fd; closing the parent copy avoids leaking it.
	logFile.Close()

	// Detect process exit so a failed attempt is retried immediately instead of
	// waiting the full per-attempt timeout for a process that has already died.
	exited := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(exited)
	}()

	ok := waitForConnectionOrExit(LocalURL, perAttempt, exited) &&
		waitForConnectionOrExit(LocalFaucetURL, perAttempt, exited)
	if !ok {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			<-exited
		}
		fmt.Fprintf(os.Stderr, "[StartSuiNode] node not reachable; sui start output tail:\n%s\n", tailLog(suiLogPath))
		return nil, false, nil
	}

	// Stabilize: the ports are up, but `sui start` can still exit ~9-12s after
	// launch if a validator health check fails, because the post-launch grace
	// period was removed in Sui 1.75.x. Wait out that window watching for process
	// exit; if it dies, retry with a fresh node instead of handing the test a
	// doomed process that will die mid-test.
	if exitedBeforeTimeout(exited, stabilizeWindowFromEnv()) {
		fmt.Fprintf(os.Stderr, "[StartSuiNode] node exited during stabilization window; retrying\n")
		return nil, false, nil
	}
	return cmd, true, nil
}

// waitForConnectionOrExit polls the TCP endpoint until it accepts a connection,
// the sui process exits, or the timeout elapses. Returning early on process exit
// lets the caller retry immediately instead of waiting the full timeout for a
// process that has already died.
func waitForConnectionOrExit(url string, timeout time.Duration, exited <-chan struct{}) bool {
	parsedURL, err := netUrl.Parse(url)
	if err != nil || parsedURL.Host == "" {
		return false
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-exited:
			return false
		default:
		}
		d := net.Dialer{Timeout: 2 * time.Second}
		conn, derr := d.Dial("tcp", parsedURL.Host)
		if derr == nil {
			conn.Close()
			return true
		}
		select {
		case <-exited:
			return false
		case <-time.After(500 * time.Millisecond):
		}
	}
	return false
}

// exitedBeforeTimeout reports whether the sui process exited before the given
// duration elapsed. Used to hold through Sui's health-check exit window.
func exitedBeforeTimeout(exited <-chan struct{}, d time.Duration) bool {
	select {
	case <-exited:
		return true
	case <-time.After(d):
		return false
	}
}

// stabilizeWindowFromEnv returns how long StartSuiNode waits after the ports come
// up to confirm `sui start` does not exit from a validator health-check failure.
// Defaults to 15s, which covers Sui's ~9-12s health-check exit window with margin,
// and is overridable via SUI_NODE_STABILIZE_SECS.
func stabilizeWindowFromEnv() time.Duration {
	const defaultWindow = 15 * time.Second
	v := os.Getenv("SUI_NODE_STABILIZE_SECS")
	if v == "" {
		return defaultWindow
	}
	secs, err := strconv.Atoi(v)
	if err != nil || secs < 0 {
		return defaultWindow
	}
	return time.Duration(secs) * time.Second
}

// startTimeoutFromEnv returns the deadline used when waiting for the local Sui node
// and faucet to accept connections. It defaults to 30s and can be overridden with
// SUI_NODE_START_TIMEOUT_SECS, e.g. to give a heavier Sui version more time to come up.
func startTimeoutFromEnv() time.Duration {
	const defaultDelay = 30 * time.Second
	v := os.Getenv("SUI_NODE_START_TIMEOUT_SECS")
	if v == "" {
		return defaultDelay
	}
	secs, err := strconv.Atoi(v)
	if err != nil || secs <= 0 {
		return defaultDelay
	}
	return time.Duration(secs) * time.Second
}

// tailLog returns the trailing bytes of the captured `sui start` log so the real
// crash/panic reason is included in the error returned by StartSuiNode. Returns an
// empty string when no log was captured, e.g. the Docker path.
func tailLog(path string) string {
	if path == "" {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return fmt.Sprintf("(could not read sui start log %s: %v)", path, err)
	}
	defer f.Close()
	const tail = 8 * 1024
	if info, statErr := f.Stat(); statErr == nil && info.Size() > tail {
		if _, seekErr := f.Seek(-tail, io.SeekEnd); seekErr != nil {
			return ""
		}
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return ""
	}
	if len(data) == 0 {
		return "(sui start produced no output)"
	}
	return "--- sui start output (tail) ---\n" + string(data)
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
