//go:build integration

package testenv

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

// Constants for magic numbers and repeated values
const (
	DefaultNodeStartTimeout = 30 * time.Second
	NodeReadyPollInterval   = 100 * time.Millisecond
	FaucetMaxRetries        = 10
	FaucetRetryDelay        = time.Second
	HTTPClientTimeout       = 30 * time.Second
	SeedSize                = ed25519.SeedSize
)

type TestEnvironment struct {
	nodeCmd    *exec.Cmd
	mu         sync.Mutex
	rpcPort    int
	faucetPort int
	logger     logger.Logger
}

var (
	instance *TestEnvironment
	errSetup error
	refCount int
	refMu    sync.Mutex
)

func SetupEnvironment(t *testing.T) (utils.SuiSigner, sui.ISuiAPI) {
	t.Helper()

	log := logger.Test(t)

	err := Setup(log)
	require.NoError(t, err, "Failed to setup test environment")

	t.Cleanup(func() {
		Cleanup()
	})

	return CreateTestAccount(t)
}

// Setup initializes the test environment. This is the low-level function that
// starts a Sui node with random ports to avoid conflicts between test packages.
//
// For most tests, use SetupEnvironment() instead, which handles cleanup automatically.
//
// If you need more control, you can use Setup/Cleanup/CreateTestAccount directly:
//
//	func TestMain(m *testing.M) {
//	    if err := testenv.Setup(); err != nil {
//	        log.Fatal(err)
//	    }
//	    code := m.Run()
//	    testenv.Cleanup()
//	    os.Exit(code)
//	}
func Setup(log logger.Logger) error {
	refMu.Lock()
	defer refMu.Unlock()

	refCount++

	if refCount == 1 {
		instance = &TestEnvironment{
			logger: log,
		}
		errSetup = instance.initialize()
	}

	return errSetup
}

func Cleanup() {
	refMu.Lock()
	defer refMu.Unlock()

	refCount--

	// only cleanup when all tests are done
	if refCount > 0 {
		return
	}

	if instance == nil {
		return
	}

	instance.cleanup()
	instance = nil
	errSetup = nil
}

// CreateTestAccount creates a new test account with funding from the faucet.
// This requires the test environment to be set up first (via Setup() or SetupEnvironment()).
func CreateTestAccount(t *testing.T) (utils.SuiSigner, sui.ISuiAPI) {
	t.Helper()

	refMu.Lock()
	defer refMu.Unlock()

	// ensure the environment is setup
	require.NoError(t, errSetup, "test setup failed")
	require.NotNil(t, instance, "test environment not initialized")

	signer, err := createAccount(t)
	require.NoError(t, err, "Failed to create test account")

	return signer, createClient()
}

func (te *TestEnvironment) initialize() error {
	te.mu.Lock()
	defer te.mu.Unlock()

	ports, err := randomPorts(2)
	if err != nil {
		return fmt.Errorf("failed to generate ports: %w", err)
	}
	te.rpcPort = ports[0]
	te.faucetPort = ports[1]

	te.logger.Infof("TestEnv: Starting Sui node with RPC port %d and faucet port %d\n", te.rpcPort, te.faucetPort)

	cmd := exec.Command("sui", "start", //nolint:gosec
		"--force-regenesis",
		"--fullnode-rpc-port", fmt.Sprintf("%d", te.rpcPort),
		fmt.Sprintf("--with-faucet=%d", te.faucetPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Sui node: %w", err)
	}
	te.nodeCmd = cmd

	// wait for node to be ready
	ctx, cancel := context.WithTimeout(context.Background(), DefaultNodeStartTimeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			te.cleanup()
			return fmt.Errorf("timeout waiting for Sui node to be ready")
		default:
			client := sui.NewSuiClient(fmt.Sprintf("http://localhost:%d", te.rpcPort))
			_, err := client.SuiGetChainIdentifier(context.Background())
			if err == nil {
				return nil
			}
			time.Sleep(NodeReadyPollInterval)
		}
	}
}

func (te *TestEnvironment) cleanup() {
	te.mu.Lock()
	defer te.mu.Unlock()

	if te.nodeCmd != nil && te.nodeCmd.Process != nil {
		te.logger.Info("TestEnv: Cleaning up Sui node")
		if err := te.nodeCmd.Process.Kill(); err != nil {
			te.logger.Info("Failed to kill Sui node process:", err)
		}
		_ = te.nodeCmd.Wait()
	}

	te.nodeCmd = nil
}

func createClient() sui.ISuiAPI {
	return sui.NewSuiClient(fmt.Sprintf("http://localhost:%d", instance.rpcPort))
}

func createAccount(t *testing.T) (utils.SuiSigner, error) {
	t.Helper()

	privateKey := ed25519.NewKeyFromSeed(randomSeed())
	signer := utils.NewTestPrivateKeySigner(privateKey)

	signerAddress, err := signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get signer address: %w", err)
	}

	faucetURL := fmt.Sprintf("http://localhost:%d", instance.faucetPort)
	log := logger.Test(t)
	err = fundWithFaucet(log, faucetURL, signerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to fund account: %w", err)
	}

	return signer, nil
}

func randomPorts(count int) ([]int, error) {
	listeners := make([]net.Listener, count)
	ports := make([]int, count)

	// open all listeners before closing to avoid the same port being returned
	for i := 0; i < count; i++ { //nolint:intrange
		listener, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			for j := 0; j < i; j++ { //nolint:intrange
				listeners[j].Close()
			}

			return nil, fmt.Errorf("failed to listen on random port %d: %w", i, err)
		}
		listeners[i] = listener
		addr := listener.Addr().(*net.TCPAddr)
		ports[i] = addr.Port
	}

	for _, listener := range listeners {
		listener.Close()
	}

	return ports, nil
}

func randomSeed() []byte {
	seed := make([]byte, SeedSize)
	_, err := rand.Read(seed)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random seed: %+v", err))
	}

	return seed
}

func fundWithFaucet(log logger.Logger, faucetURL string, address string) error {
	log.Infof("Funding account %s using faucet at %s", address, faucetURL)

	var lastErr error
	client := &http.Client{Timeout: HTTPClientTimeout}

	for range FaucetMaxRetries {
		requestURL := fmt.Sprintf("%s/gas", faucetURL)
		jsonBody := fmt.Sprintf(`{"FixedAmountRequest": {"recipient": "%s"}}`, address)

		req, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(jsonBody))
		if err != nil {
			lastErr = err
			time.Sleep(FaucetRetryDelay)

			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(FaucetRetryDelay)

			continue
		}

		closeResponse := func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Warnf("Failed to close response body: %v", closeErr)
			}
		}
		closeResponse()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			log.Infof("Successfully funded account %s", address)
			return nil
		}

		lastErr = fmt.Errorf("faucet returned status %d", resp.StatusCode)
		time.Sleep(FaucetRetryDelay)
	}

	return fmt.Errorf("failed to fund account after retries: %w", lastErr)
}
