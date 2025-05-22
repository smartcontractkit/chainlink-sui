package testutils

import (
	"math/big"
	"os/exec"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

const (
	maxConcurrentRequests     = 5
	defaultTransactionTimeout = 10 * time.Second
	defaultNumberRetries      = 5
	defaultGasLimit           = 10000000
)

type TestState struct {
	AccountAddress  string
	SuiGateway      *client.PTBClient
	KeystoreGateway *keystore.SuiKeystore
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
func SetupClients(t *testing.T, rpcURL string, _keystore loop.Keystore) (*client.PTBClient, *txm.SuiTxm, *txm.InMemoryStore) {
	t.Helper()

	logg, err := logger.New()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	relayerClient, err := client.NewPTBClient(logg, rpcURL, nil, defaultTransactionTimeout, _keystore, maxConcurrentRequests, "WaitForEffectsCert")
	if err != nil {
		t.Fatalf("Failed to create relayer client: %v", err)
	}

	t.Log("relayerClient", relayerClient)

	store := txm.NewTxmStoreImpl()
	conf := txm.DefaultConfigSet

	retryManager := txm.NewDefaultRetryManager(defaultNumberRetries)
	gasLimit := big.NewInt(defaultGasLimit)
	gasManager := txm.NewSuiGasManager(logg, relayerClient, *gasLimit, 0)

	txManager, err := txm.NewSuiTxm(logg, relayerClient, _keystore, conf, store, retryManager, gasManager)
	if err != nil {
		t.Fatalf("Failed to create SuiTxm: %v", err)
	}

	return relayerClient, txManager, store
}

// BootstrapTestEnvironment sets up a complete test environment for a Sui relayer integration test.
// It performs the following steps:
//  1. Starts a local Sui node of the specified NodeEnvType.
//  2. Initializes a new keystore and retrieves the account address and associated signer.
//  3. Funds the account using the Sui faucet.
//  4. Iterates over the provided contracts metadata:
//     - Builds and compiles each contract.
//     - Publishes the contract to the Sui network.
//     - Queries and records the contract object IDs.
//  5. Initializes the Sui and relayer clients as well as the transaction manager (TXM).
//  6. Aggregates all the resources into a TestState struct, which is subsequently used for testing.
//
// Parameters:
//   - t: testing.T used for logging and error handling.
//   - nodeType: the type of node to start (NodeEnvType).
//   - contractsMetadata: a slice of Contracts containing metadata for the contracts to be deployed.
//
// Returns:
//   - *TestState: a pointer to the fully bootstrapped test environment containing:
//   - AccountAddress (string)
//   - SuiGateway (*client.PTBClient)
//   - KeystoreGateway (keystore.Keystore)
//   - TxManager (*txm.SuiTxm)
//   - TxStore (*txm.InMemoryStore)
//   - Signer (signer.SuiSigner)
//   - Contracts ([]Contracts)
//   - Cmd (exec.Cmd) — the running process for the Sui node.
//
// Example usage:
//
//	func TestMyFunction(t *testing.T) {
//	    contractsMeta := []testutils.Contracts{
//	        {
//	            Path: "contracts/mycontract/",
//	            Name: "mycontract",
//	            Objects: []testutils.ContractObject{
//	                {
//	                    ObjectID:    "0x123",
//	                    PackageName: "mycontract",
//	                    StructName:  "MyStruct",
//	                },
//	            },
//	        },
//	    }
//
//	    // Bootstraps the environment with a local Sui node and provided contracts metadata.
//	    testState := testutils.BootstrapTestEnvironment(t, testutils.CLI, contractsMeta)
//
//	    // Use testState to access the account address, TXM, client, etc.
//	    // Example: reading the account address
//	    t.Logf("Account Address: %s", testState.AccountAddress)
//
//	    // Continue with further tests using testState...
//	}
func BootstrapTestEnvironment(t *testing.T, nodeType NodeEnvType, contractsMetadata []Contracts) *TestState {
	t.Helper()
	_logger := logger.Test(t)
	_logger.Debugw("Starting Sui node")

	cmd, err := StartSuiNode(nodeType)
	require.NoError(t, err)

	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	_keystore, err := keystore.NewSuiKeystore(_logger, "")
	require.NoError(t, err)
	accountAddress := GetAccountAndKeyFromSui(t, _logger)

	err = FundWithFaucet(_logger, constant.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	contractsState := []Contracts{}

	for _, contract := range contractsMetadata {
		_logger.Infow("Building contract", contract)
		contractPath := BuildSetup(t, contract.Path)
		BuildContract(t, contractPath)
		packageId, publishOutput, err := PublishContract(t, contract.Name, contractPath, accountAddress, nil)
		_logger.Infow("Publish contract", "packageId", packageId, "Contract Name", contract.Name)
		require.NoError(t, err)
		c := &Contracts{
			Path:     contractPath,
			Name:     contract.Name,
			ModuleID: packageId,
			Objects:  []ContractObject{},
		}

		for _, obj := range contract.Objects {
			objectId, err := QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, obj.PackageName, obj.StructName)
			_logger.Debugw("Query object ID", "objectId", objectId, "PackageName", obj.PackageName, "StructName", obj.StructName)
			require.NoError(t, err)

			c.Objects = append(
				c.Objects, ContractObject{
					ObjectID:    objectId,
					PackageName: obj.PackageName,
					StructName:  obj.StructName,
				})
		}
		contractsState = append(contractsState, *c)
		_logger.Debugw("Contract state", contractsState)
	}

	suiClient, txManager, txStore := SetupClients(t, LocalUrl, _keystore)

	return &TestState{
		AccountAddress:  accountAddress,
		SuiGateway:      suiClient,
		KeystoreGateway: &_keystore,
		TxManager:       txManager,
		TxStore:         txStore,
		Contracts:       contractsState,
		Cmd:             *cmd,
	}
}
