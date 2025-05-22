//go:build integration

package mcms_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type TestSetupOutputs struct {
	mcmsPackageId         string
	mcmsTestPackageId     string
	mcmsPublishOutput     testutils.TxnMetaWithObjectChanges
	mcmsTestPublishOutput testutils.TxnMetaWithObjectChanges
}

const (
	MINUTE_IN_MS = 60_000
)

// ------------------------------------------
//
//	Setup and Helpers
//
// ------------------------------------------
// setupTestEnvironment sets up the test environment with a local Sui node and deploys the counter contract
func setupTestEnvironment(t *testing.T) (
	log logger.Logger,
	accountAddress string,
	ptbClient *client.PTBClient,
	outputs TestSetupOutputs,
	publicKeyBytes []byte,
) {
	t.Helper()

	log = logger.Test(t)
	accountAddress = testutils.GetAccountAndKeyFromSui(t, log)

	// setup keystore instance
	keystoreInstance, keystoreErr := keystore.NewSuiKeystore(log, "")
	require.NoError(t, keystoreErr)

	// get the private key from keystore
	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes = []byte(publicKey)

	// Start local Sui node
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	// Ensure the process is killed when the test completes
	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	log.Debugw("Started Sui node")

	// Fund the account
	err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	// Create client
	ptbClient, err = client.NewPTBClient(log, testutils.LocalUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, err)

	// Build and publish contracts
	gasBudget := int(2000000000)
	mcmsContractPath := testutils.BuildSetup(t, "contracts/mcms/mcms")
	mcmsPackageId, mcmsPublishOutput, err := testutils.PublishContract(t, "ChainlinkManyChainMultisig", mcmsContractPath, "", &gasBudget)
	require.NoError(t, err)
	log.Debugw("Published MCMS Contract", "packageId", mcmsPackageId)

	mcmsTestContractPath := testutils.BuildSetup(t, "contracts/mcms/mcms_test")
	testutils.PatchContractDevAddressTOML(t, mcmsTestContractPath, "mcms", mcmsPackageId)
	mcmsTestPackageId, mcmsTestPublishOutput, err := testutils.PublishContract(
		t,
		"TestModule",
		mcmsTestContractPath,
		fmt.Sprint("mcms=", mcmsPackageId),
		&gasBudget,
	)
	require.NoError(t, err)
	log.Debugw("Published MCMS Test Contract", "packageId", mcmsTestPackageId)

	// Debug print results
	testutils.PrettyPrintDebug(log, mcmsTestPublishOutput, "mcms_test_publish")
	testutils.PrettyPrintDebug(log, mcmsPublishOutput, "mcms_publish")

	outputs = TestSetupOutputs{
		mcmsPackageId:         mcmsPackageId,
		mcmsTestPackageId:     mcmsTestPackageId,
		mcmsPublishOutput:     mcmsPublishOutput,
		mcmsTestPublishOutput: mcmsTestPublishOutput,
	}

	return log, accountAddress, ptbClient, outputs, publicKeyBytes
}

// ------------------------------------------------
//
//	MCMS Test Contract
//
// ------------------------------------------------
//
//nolint:paralleltest
func TestMCMS(t *testing.T) {
	// Set up the test environment
	log, _, ptbClient, outputs, pubKeyBytes := setupTestEnvironment(t)

	// Extract relevant IDs
	ownerCapObjId, err := testutils.QueryCreatedObjectID(outputs.mcmsPublishOutput.ObjectChanges, outputs.mcmsPackageId, "mcms_account", "OwnerCap")
	require.NoError(t, err)
	accountStateObjId, err := testutils.QueryCreatedObjectID(outputs.mcmsPublishOutput.ObjectChanges, outputs.mcmsPackageId, "mcms_account", "AccountState")
	require.NoError(t, err)
	multisigStateObjId, err := testutils.QueryCreatedObjectID(outputs.mcmsPublishOutput.ObjectChanges, outputs.mcmsPackageId, "mcms", "MultisigState")
	require.NoError(t, err)
	timelockObjId, err := testutils.QueryCreatedObjectID(outputs.mcmsPublishOutput.ObjectChanges, outputs.mcmsPackageId, "mcms", "Timelock")
	require.NoError(t, err)
	registryObjId, err := testutils.QueryCreatedObjectID(outputs.mcmsPublishOutput.ObjectChanges, outputs.mcmsPackageId, "mcms_registry", "Registry")
	require.NoError(t, err)

	log.Debugw("OwnerCap object created", "ownerCapObjectId", ownerCapObjId)
	log.Debugw("AccountState object created", "accountStateObjectId", accountStateObjId)
	log.Debugw("MultisigState object created", "multisigStateObjectId", multisigStateObjId)
	log.Debugw("Timelock object created", "timelockObjectId", timelockObjId)
	log.Debugw("Registry object created", "registryObjectId", registryObjId)

	// Some PTB functions that will be re-used in the chainwriter config below
	setConfigOnlyPTBFunc := chainwriter.ChainWriterPTBCommand{
		// set_config
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: &outputs.mcmsPackageId,
		ModuleId:  testutils.StringPointer("mcms"),
		Function:  testutils.StringPointer("set_config"),
		Params: []codec.SuiFunctionParam{
			{
				Name:     "owner_cap_id",
				Type:     "object_id",
				Required: true,
			},
			{
				Name:     "multisig_state_id",
				Type:     "object_id",
				Required: true,
			},
			{
				Name:     "role",
				Type:     "number",
				Required: true,
			},
			{
				Name:     "chain_id",
				Type:     "number",
				Required: true,
			},
			{
				Name:     "signer_addresses",
				Type:     "vector<vector<u8>>",
				Required: true,
			},
			{
				Name:     "signer_groups",
				Type:     "vector<u8>",
				Required: true,
			},
			{
				Name:     "group_quorums",
				Type:     "vector<u8>",
				Required: true,
			},
			{
				Name:     "group_parents",
				Type:     "vector<u8>",
				Required: true,
			},
			{
				Name:     "clear_root",
				Type:     "bool",
				Required: true,
			},
		},
	}

	setRootPTBFunc := chainwriter.ChainWriterPTBCommand{
		// set_root
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: &outputs.mcmsPackageId,
		ModuleId:  testutils.StringPointer("mcms"),
		Function:  testutils.StringPointer("set_root"),
		Params: []codec.SuiFunctionParam{
			{
				Name:     "multisig_state_id",
				Type:     "object_id",
				Required: true,
			},
			{
				Name:         "clock",
				Type:         "object_id",
				Required:     true,
				DefaultValue: "0x06",
				IsMutable:    testutils.BoolPointer(false),
			},
			{
				Name:     "role",
				Type:     "u8",
				Required: true,
			},
			{
				Name:     "root",
				Type:     "vector<u8>",
				Required: true,
			},
			{
				Name:     "valid_until",
				Type:     "u64",
				Required: true,
			},
			{
				Name:     "chain_id",
				Type:     "u256",
				Required: true,
			},
			{
				Name:     "multisig_addr",
				Type:     "object_id",
				Required: true,
			},
			{
				Name:     "pre_op_count",
				Type:     "u64",
				Required: true,
			},
			{
				Name:     "post_op_count",
				Type:     "u64",
				Required: true,
			},
			{
				Name:     "override_previous_root",
				Type:     "bool",
				Required: true,
			},
			{
				Name:     "metadata_proof",
				Type:     "vector<vector<u8>>",
				Required: true,
			},
			{
				Name:     "signatures",
				Type:     "vector<vector<u8>>",
				Required: true,
			},
		},
	}

	executePTBFunc := chainwriter.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: &outputs.mcmsPackageId,
		ModuleId:  testutils.StringPointer("mcms"),
		Function:  testutils.StringPointer("execute"),
		Params: []codec.SuiFunctionParam{
			{
				Name:     "multisig_state_id",
				Type:     "object_id",
				Required: true,
			},
			{
				Name:         "clock",
				Type:         "object_id",
				Required:     true,
				DefaultValue: "0x06",
				IsMutable:    testutils.BoolPointer(false),
			},
			{
				Name:     "role",
				Type:     "u8",
				Required: true,
			},
			{
				Name:     "chain_id",
				Type:     "u256",
				Required: true,
			},
			{
				Name:     "multisig_addr",
				Type:     "vector<u8>",
				Required: true,
			},
			{
				Name:     "nonce",
				Type:     "u64",
				Required: true,
			},
			{
				Name:     "to",
				Type:     "address",
				Required: true,
			},
			{
				Name:     "module_name",
				Type:     "string",
				Required: true,
			},
			{
				Name:     "function_name",
				Type:     "string",
				Required: true,
			},
			{
				Name:     "data",
				Type:     "vector<u8>",
				Required: true,
			},
			{
				Name:     "proof",
				Type:     "vector<vector<u8>>",
				Required: true,
			},
		},
	}

	timelockScheduleBatchPTBFunc := chainwriter.ChainWriterPTBCommand{
		Type:      codec.SuiPTBCommandMoveCall,
		PackageId: &outputs.mcmsPackageId,
		ModuleId:  testutils.StringPointer("mcms"),
		Function:  testutils.StringPointer("dispatch_timelock_schedule_batch"),
		Params: []codec.SuiFunctionParam{
			{
				Name:     "timelock",
				Type:     "object_id",
				Required: true,
			},
			{
				Name:      "clock",
				Type:      "object_id",
				Required:  true,
				IsMutable: testutils.BoolPointer(false),
			},
			{
				Name: "timelock_callback_params",
				Type: "hot_potato",
				PTBDependency: &codec.PTBCommandDependency{
					CommandIndex: 2,
					ResultIndex:  nil,
				},
				Required: true,
			},
		},
	}

	// Create PTB Constructor config
	config := chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			"mcms_ptb_test": {
				Name:     "mcms_ptb_test",
				ModuleID: "0x123",
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"set_config": {
						Name: "set_config",
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							setConfigOnlyPTBFunc,
						},
					},
					"set_config_and_root": {
						Name: "set_config_and_root",
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							setConfigOnlyPTBFunc,
							setRootPTBFunc,
						},
					},
					"timelock_execute": {
						Name: "timelock_execute",
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							// set_config
							setConfigOnlyPTBFunc,
							// set_root
							setRootPTBFunc,
							// execute
							executePTBFunc,
							// timelock schedule batch
							timelockScheduleBatchPTBFunc,
						},
					},
				},
			},
		},
	}

	constructor := chainwriter.NewPTBConstructor(config, ptbClient, log)
	require.NotNil(t, constructor)

	ctx := context.Background()

	//nolint:paralleltest
	t.Run("MCMS set_config with invalid role", func(t *testing.T) {
		// Create quorums array for config
		quorums := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			quorums[i] = 0
		}
		quorums[0] = 2

		// Create parents array for config
		parents := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			parents[i] = 0
		}

		// Create signer addresses
		signerAddresses := [][]byte{
			{0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			{0x00, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			{0x00, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
		}

		args := chainwriter.Arguments{
			Args: map[string]any{
				"owner_cap_id":      ownerCapObjId,
				"multisig_state_id": multisigStateObjId,
				"role":              uint8(100),
				"chain_id":          uint256.Int{77},
				"signer_addresses":  signerAddresses,
				"signer_groups":     []uint8{0, 0, 0},
				"group_quorums":     quorums,
				"group_parents":     parents,
				"clear_root":        false,
			},
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "mcms_ptb_test", "set_config", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, pubKeyBytes, ptb)
		testutils.PrettyPrintDebug(log, ptbResult, "ptb_result")
		require.NoError(t, err)
		require.Equal(t, "failure", ptbResult.Status.Status)
	})

	//nolint:paralleltest
	t.Run("MCMS set_config with valid role", func(t *testing.T) {
		// Create quorums array for config
		quorums := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			quorums[i] = 0
		}
		quorums[0] = 2

		// Create parents array for config
		parents := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			parents[i] = 0
		}

		// Create signer addresses
		signerAddresses := [][]byte{
			{0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			{0x00, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			{0x00, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
		}

		args := chainwriter.Arguments{
			Args: map[string]any{
				"owner_cap_id":      ownerCapObjId,
				"multisig_state_id": multisigStateObjId,
				"role":              uint8(1),
				"chain_id":          uint256.Int{77},
				"signer_addresses":  signerAddresses,
				"signer_groups":     []uint8{0, 0, 0},
				"group_quorums":     quorums,
				"group_parents":     parents,
				"clear_root":        false,
			},
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "mcms_ptb_test", "set_config", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, pubKeyBytes, ptb)
		testutils.PrettyPrintDebug(log, ptbResult, "ptb_result")
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
	})

	//nolint:paralleltest
	t.Run("MCMS set config and set root with invalid root and signatures", func(t *testing.T) {
		acc_1_priv_key, _ := testutils.GenerateFromHexSeed(testutils.ACCOUNT_1_SEED)
		acc_2_priv_key, _ := testutils.GenerateFromHexSeed(testutils.ACCOUNT_2_SEED)
		acc_3_priv_key, _ := testutils.GenerateFromHexSeed(testutils.ACCOUNT_3_SEED)

		// Get addresses from private keys
		acc_1_addr := crypto.PubkeyToAddress(acc_1_priv_key.PublicKey).Bytes()
		acc_2_addr := crypto.PubkeyToAddress(acc_2_priv_key.PublicKey).Bytes()
		acc_3_addr := crypto.PubkeyToAddress(acc_3_priv_key.PublicKey).Bytes()

		// Create a slice of addresses
		addresses := [][]byte{acc_1_addr, acc_2_addr, acc_3_addr}

		// Sort addresses lexicographically
		sort.Slice(addresses, func(i, j int) bool {
			return bytes.Compare(addresses[i], addresses[j]) < 0
		})

		// Create quorums array for config
		quorums := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			quorums[i] = 0
		}
		quorums[0] = 2

		// Create parents array for config
		parents := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			parents[i] = 0
		}

		args := chainwriter.Arguments{
			Args: map[string]any{
				"owner_cap_id":           ownerCapObjId,
				"multisig_state_id":      multisigStateObjId,
				"role":                   uint8(0),
				"chain_id":               uint256.Int{77},
				"signer_addresses":       addresses,
				"signer_groups":          []uint8{0, 0, 0},
				"group_quorums":          quorums,
				"group_parents":          parents,
				"clear_root":             false,
				"root":                   [][]byte{},
				"multisig_addr":          []byte(outputs.mcmsPackageId),
				"pre_op_count":           uint64(0),
				"post_op_count":          uint64(1),
				"override_previous_root": false,
				"metadata_proof":         [][]byte{{0x01}},
				"signatures":             [][]byte{{0x05}, {0x06}, {0x07}},
				"valid_until":            uint64(999),
				"clock":                  "0x06",
			},
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "mcms_ptb_test", "set_config_and_root", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, pubKeyBytes, ptb)
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "failure", ptbResult.Status.Status)
	})

	//nolint:paralleltest
	t.Run("MCMS timelock execution happy path", func(t *testing.T) {
		proposerRole := uint8(2)
		signers := []testutils.ECDSAKeyPair{}

		for range 3 {
			acc_priv_key, err := testutils.GenerateKeyPair()
			require.NoError(t, err)

			signers = append(signers, testutils.ECDSAKeyPair{
				PrivateKey: *acc_priv_key,
				PublicKey:  acc_priv_key.PublicKey,
				Address:    crypto.PubkeyToAddress(acc_priv_key.PublicKey).Bytes(),
			})
		}

		// Sort addresses lexicographically
		sort.Slice(signers, func(i, j int) bool {
			return bytes.Compare(signers[i].Address, signers[j].Address) < 0
		})

		// Parse the package ID string into an address
		packageAddr, err := sui.AddressFromHex(outputs.mcmsPackageId)
		require.NoError(t, err)

		// Convert to bytes in the correct format
		packageBytes := packageAddr.Bytes()

		// override package address bytes to zeros because the address that is used in [dev-addresses]
		// in the mcms's TOML file is 0x00
		for i := range 32 {
			packageBytes[i] = 0
		}

		// Construct a timelock operation and serialize its data to be included in the main operation's
		// `data` entry
		scheduleBatchOps := []testutils.TimelockOperation{
			{
				Target:       packageBytes,
				ModuleName:   "mcms",
				FunctionName: "timelock_schedule_batch",
				Data:         []byte{},
			},
		}
		scheduleBatchPredecessor := []byte{}
		salt := []byte{}
		delay := uint64(0)
		serializedBatchData, err := testutils.SerializeScheduleBatchParams(scheduleBatchOps, scheduleBatchPredecessor, salt, delay)
		require.NoError(t, err)

		// Create a valid root for the test using the merkle tree utilities
		// Define the operation for the merkle tree
		op := testutils.Op{
			Role:         proposerRole,
			ChainID:      big.NewInt(77),
			MultiSig:     packageBytes,
			Nonce:        0,
			To:           packageBytes, // Target address for the operation
			ModuleName:   "mcms",
			FunctionName: "timelock_schedule_batch",
			Data:         serializedBatchData, // Operation data
		}

		// Define the root metadata
		rootMetadata := testutils.RootMetadata{
			Role:                 proposerRole,
			ChainID:              big.NewInt(77),
			MultiSig:             packageBytes,
			PreOpCount:           0,
			PostOpCount:          1,
			OverridePreviousRoot: false,
		}

		// Generate the merkle tree with the operation and metadata
		merkleTree, err := testutils.GenerateMerkleTree([]testutils.Op{op}, rootMetadata)
		require.NoError(t, err)

		// Get the root and proof for the metadata leaf (index 0)
		root := merkleTree.GetRoot()
		rootBytes := make([]byte, 32)
		copy(rootBytes, root[:])
		metadataProof := merkleTree.GetProof(0)

		// Convert proof to the format expected by the contract
		metadataProofBytes := make([][]byte, len(metadataProof))
		for i, p := range metadataProof {
			metadataProofBytes[i] = p[:]
		}

		// Calculate the hash that needs to be signed
		//nolint:all
		validUntil := uint64(time.Now().UnixMilli() + MINUTE_IN_MS)
		signedHash := testutils.CalculateSignedHash(root, validUntil)

		log.Debugw("signedHash", "value", signedHash, "length", len(signedHash))

		// Generate signatures
		signatures := testutils.GenerateSignatures(t,
			[]ecdsa.PrivateKey{signers[0].PrivateKey, signers[1].PrivateKey, signers[2].PrivateKey},
			signedHash)

		for i, signature := range signatures {
			_ = testutils.VerifySignatureRecovery(t, signature, signedHash[:], signers[i].Address)
			log.Debugf("Address [%d]: %d", i, signers[i].Address)
		}

		// Prepare signer addresses array
		signerAddresses := make([][]uint8, 3)
		for i := range 3 {
			signerAddresses[i] = signers[i].Address
		}

		// Create quorums array
		quorums := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			quorums[i] = 0
		}
		quorums[0] = 2

		// Create parents array
		parents := make([]uint8, 32)
		for i := 1; i < 32; i++ {
			parents[i] = 0
		}

		// Get the proof for the next operation
		executeProof := merkleTree.GetProof(1)
		executeProofBytes := make([][]byte, len(executeProof))
		for i, p := range executeProof {
			executeProofBytes[i] = p[:]
		}

		log.Debugw("rootBytes", "value", rootBytes, "length", len(rootBytes))

		args := chainwriter.Arguments{
			Args: map[string]any{
				// Args for set_config
				"owner_cap_id":      ownerCapObjId,
				"multisig_state_id": multisigStateObjId,
				"role":              proposerRole,
				"chain_id":          uint256.Int{77},
				"signer_addresses":  signerAddresses,
				"signer_groups":     []uint8{0, 0, 0},
				"group_quorums":     quorums,
				"group_parents":     parents,
				"clear_root":        false,

				// Additional args for set_root operation
				"root":                   rootBytes,
				"multisig_addr":          packageBytes,
				"pre_op_count":           uint64(0),
				"post_op_count":          uint64(1),
				"override_previous_root": false,
				"metadata_proof":         metadataProofBytes,
				"signatures":             signatures,
				"valid_until":            validUntil,
				"clock":                  "0x06",

				// Additional args for execute operation
				"nonce":         uint64(0),
				"to":            packageBytes,
				"module_name":   "mcms",
				"function_name": "timelock_schedule_batch",
				"data":          serializedBatchData,
				"proof":         executeProofBytes,

				// Additional args for timelock_execute_batch
				"timelock": timelockObjId,
			},
		}

		ptb, err := constructor.BuildPTBCommands(ctx, "mcms_ptb_test", "timelock_execute", args, nil)
		require.NoError(t, err)
		require.NotNil(t, ptb)

		// Execute the PTB command
		ptbResult, err := ptbClient.FinishPTBAndSend(ctx, pubKeyBytes, ptb)
		testutils.PrettyPrintDebug(log, ptbResult, "ptb_result")
		require.NoError(t, err)
		require.NotEmpty(t, ptbResult)
		require.Equal(t, "success", ptbResult.Status.Status)
	})
}
