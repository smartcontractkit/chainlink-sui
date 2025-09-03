package testutils

import (
	"math/big"
	"os/exec"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
	"github.com/test-go/testify/require"
	"golang.org/x/net/context"
)

const (
	maxConcurrentRequests     = 5
	defaultTransactionTimeout = 10 * time.Second
	defaultNumberRetries      = 5
	defaultGasLimit           = 10000000
	waitTimeNextTest          = 2 * time.Second
)

type TestState struct {
	AccountAddress  string
	PublicKeyBytes  []byte
	SuiGateway      *client.PTBClient
	KeystoreGateway loop.Keystore
	TxManager       *txm.SuiTxm
	TxStore         *txm.InMemoryStore
	Contracts       []Contracts
	Cmd             exec.Cmd
}

type ContractObject struct {
	ObjectID    string
	PackageName string
	StructName  string
}

type Contracts struct {
	Path     string
	Name     string
	ModuleID string
	Objects  []ContractObject
}

// setupClients initializes the Sui and relayer clients.
func SetupClients(
	t *testing.T,
	rpcURL string,
	keystore loop.Keystore,
	logg logger.Logger,
	gasLimit int64,
) (*client.PTBClient, *txm.SuiTxm, *txm.InMemoryStore) {
	t.Helper()

	relayerClient, err := client.NewPTBClient(logg, rpcURL, nil, defaultTransactionTimeout, keystore, maxConcurrentRequests, "WaitForEffectsCert")
	if err != nil {
		t.Fatalf("Failed to create relayer client: %v", err)
	}

	t.Log("relayerClient", relayerClient)

	lggr := logger.Named(logg, "testutils")

	store := txm.NewTxmStoreImpl(lggr)
	conf := txm.DefaultConfigSet

	retryManager := txm.NewDefaultRetryManager(defaultNumberRetries)
	// Set max gas budget to be higher than provided gas limit to allow gas bumping
	maxGasBudget := big.NewInt(gasLimit * 2) // 2x the gas limit as max budget
	gasManager := txm.NewSuiGasManager(logg, relayerClient, *maxGasBudget, 0)

	txManager, err := txm.NewSuiTxm(logg, relayerClient, keystore, conf, store, retryManager, gasManager)
	if err != nil {
		t.Fatalf("Failed to create SuiTxm: %v", err)
	}

	return relayerClient, txManager, store
}

func SetupTestEnv(
	t *testing.T,
	ctx context.Context,
	lgr logger.Logger,
	gasLimit int64,
) (*client.PTBClient, *txm.SuiTxm, *txm.InMemoryStore, string, *TestKeystore, []byte, string, string) {
	cmd, err := StartSuiNode(CLI)
	require.NoError(t, err)

	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	// Used to wait for the tear down of one test before starting the next
	// since they both depend on the Sui node running on the same port
	time.Sleep(waitTimeNextTest)

	keystoreInstance := NewTestKeystore(t)
	accountAddress, publicKeyBytes := GetAccountAndKeyFromSui(keystoreInstance)

	faucetFundErr := FundWithFaucet(lgr, SuiLocalnet, accountAddress)
	require.NoError(t, faucetFundErr)

	contractPath := BuildSetup(t, "contracts/test")
	gasBudget := int(2000000000)
	packageId, tx, err := PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotNil(t, packageId)
	require.NotNil(t, tx)

	lgr.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := QueryCreatedObjectID(tx.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	suiClient, txManager, transactionRepository := SetupClients(t, LocalUrl, keystoreInstance, lgr, gasLimit)

	return suiClient, txManager, transactionRepository, accountAddress, keystoreInstance, publicKeyBytes, packageId, counterObjectId
}
