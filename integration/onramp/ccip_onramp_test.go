//go:build testnet_integration

package ccip_test

import (
	"context"
	"crypto/ed25519"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"

	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

// helpers
// ContractAddresses holds the addresses of all deployed CCIP contracts
type ContractAddresses struct {
	// SUI contract addresses
	LinkTokenPackageID    string
	LinkTokenCoinMetadata string
	LinkTokenTreasuryCap  string

	// CCIP
	CCIPPackageID string
	CCIPStateRef  string

	// Clock object
	ClockObject string

	// Token Pools
	LinkLockReleaseTokenPool         string
	LinkLockReleaseTokenPoolState    string
	LinkLockReleaseTokenPoolOwnerCap string

	// ETH Mint Burn Token Pools
	ETHMintBurnTokenPool         string
	ETHMintBurnTokenPoolState    string
	ETHMintBurnTokenPoolOwnerCap string

	// CCIP Onramp
	CCIPOnrampPackageID string
	CCIPOnrampState     string
	CCIPOnrampOwnerCap  string

	// Coin Objects
	LinkCoinObjects []string
	ETHCoinObjects  []string
}

// setupContractAddresses sets up the addresses of all deployed CCIP contracts
func SetupContractAddresses() ContractAddresses {
	return ContractAddresses{
		// SUI contracts
		LinkTokenPackageID:    "0xe3c005c4195ec60a3468ce01238df650e4fedbd36e517bf75b9d2ee90cce8a8b",
		LinkTokenCoinMetadata: "0x2b7aee90f1ce4d6a34bed21d954fcdab04fdf391dd3a012b65641a0a8a2c5f7a",
		LinkTokenTreasuryCap:  "0x0bcba9548545dd8b5563580b3466102ce1267e0ec8a80c26f259b820ff366c02",

		// ETH Token

		// Clock object
		ClockObject: "0x6",

		CCIPPackageID: "0x1245ccd9b14d187b00f12f37906271f18c334fb9fc1d83aa1261acda571e8746",

		// CCIP contracts
		CCIPStateRef: "0x9e390b3af3d8047a54c69f587915a2705c6c5988a70744e36c241ef592d03ae5",

		// Token Pools
		LinkLockReleaseTokenPool:         "0xc94c375d9f6f279837e3efe0edb176868b0c82dc85100652f0b27bd4d1333eae",
		LinkLockReleaseTokenPoolState:    "0x1ce09f2b2b236f96cadda63be3bd327faef9359af3ce2a838a19af617cafe136",
		LinkLockReleaseTokenPoolOwnerCap: "0xe016e6773d518b8aef268c0a16ee78eda4a9da98e1f9004132358fcb7a15358c",

		// Mint Burn Token Pools
		ETHMintBurnTokenPool:         "0xd630d8ff05bab63de7052ac2fdc83add042a85b6101cc2d2dacd9e93b3641b31",
		ETHMintBurnTokenPoolState:    "0xd5d4c474825d3b01cd67da90012d6ddb4727d8893748f88b9536904c4564ef85",
		ETHMintBurnTokenPoolOwnerCap: "0x9d4a840ccb14ecbc2fdff6e302493d4ca2af884c07d996faa8eb63f5f0d902dc",

		// CCIP Onramp
		CCIPOnrampPackageID: "0xc35e20e484c080e6751463e8338dc549ad6b8ef8e730622f22a5a7d793acd544",
		CCIPOnrampState:     "0x4c4e4b73dce27d6e97eff438a55ae9167d56f89bf807015ef71a17ab95b09791",
		CCIPOnrampOwnerCap:  "0xd3cdaae719b15e5281df754e8e28ff88ff0c685f45a93c47abbe3333a7e64ddc",

		// Coin Objects (available LINK and ETH)
		LinkCoinObjects: []string{
			"0xafc6373c8d6878fa165bf878bc9b28eddaa1fde2d9e9abbfc00f2233454c5096",
			"0x21e87d51607bd019038a67ce72e51069568d541352c50b88f7e58605c34c1f92",
		},
		ETHCoinObjects: []string{
			"0x00c1b6dca46a7bbd260f931b5e6dc9b6a48fa4057ca784dc667d032768adf01e",
			"0xb567558579df9108412a0df9ce82997cec98f55e0e32ad902eb70843f737c5e3",
			"0x6b5588a48b94f8e99b8b5ded7df3c09c9aa6045d45266073722261c922643c20",
		},
	}
}

// TestCCIPOnrampSend tests the CCIP onramp send functionality
func TestCCIPOnrampSend(t *testing.T) {
	lggr := logger.Test(t)

	// Setup addresses for the test
	addresses := SetupContractAddresses()

	// Create keystore and get account
	keystoreInstance, err := keystore.NewSuiKeystore(lggr, "")
	require.NoError(t, err)

	accountAddress := testutils.GetAccountAndKeyFromSui(t, lggr)
	lggr.Infow("Using account", "address", accountAddress)

	// Get private key for signing
	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes := []byte(publicKey)

	_, txManager, _ := testutils.SetupClients(t, testutils.TestnetUrl, keystoreInstance)

	t.Run("CCIP Send two tokens to lock release token pool", func(t *testing.T) {
		//t.Skip("Skipping test")

		// Set up arguments for the PTB
		ptbArgs := createCCIPSendPForTwoTokensTBArgs(addresses)
		txID := "ccip_send_test_two_tokens"

		chainWriterConfig := configureChainWriterForMultipleTokens(addresses, publicKeyBytes)
		chainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, chainWriterConfig, false)

		c := context.Background()
		ctx, cancel := context.WithCancel(c)
		defer cancel()

		err = chainWriter.Start(ctx)
		require.NoError(t, err)

		lggr.Infow("ptbArgs", ptbArgs)
		lggr.Infow("addresses", addresses)
		lggr.Infow("publicKeyBytes", publicKeyBytes)
		lggr.Infow("chainWriter", chainWriter)
		lggr.Infow("chainWriterConfig", chainWriterConfig)

		lggr.Infow("Submitting transaction",
			"txID", txID,
			"accountAddress", accountAddress,
			"ptbArgs", ptbArgs,
			"chainWriterConfig", chainWriterConfig)

		err = chainWriter.SubmitTransaction(ctx,
			chainwriter.PTBChainWriterModuleName,
			"ccip_send",
			&ptbArgs,
			txID,
			accountAddress,
			&commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			nil,
		)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
			if statusErr != nil {
				return false
			}

			return status == commonTypes.Finalized
		}, 5*time.Second, 1*time.Second, "Transaction final state not reached")

		chainWriter.Close()
	})
}

func configureChainWriterForMultipleTokens(addresses ContractAddresses, publicKeyBytes []byte) chainwriter.ChainWriterConfig {
	return chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			chainwriter.PTBChainWriterModuleName: {
				Name:     chainwriter.PTBChainWriterModuleName,
				ModuleID: "0x123",
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"ccip_send": {
						Name:      "ccip_send",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							// First command: create token params
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: strPtr(addresses.CCIPPackageID),
								ModuleId:  strPtr("dynamic_dispatcher"),
								Function:  strPtr("create_token_params"),
								Params:    []codec.SuiFunctionParam{},
							},
							// Second command: lock tokens in the token pool
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: strPtr(addresses.LinkLockReleaseTokenPool),
								ModuleId:  strPtr("lock_release_token_pool"),
								Function:  strPtr("lock_or_burn"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "ref",
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
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:      "c",
										Type:      "object_id",
										Required:  true,
										IsGeneric: true,
									},
									{
										Name:     "remote_chain_selector",
										Type:     "u64",
										Required: true,
									},
									{
										Name:     "token_params",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
										},
									},
								},
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: strPtr(addresses.CCIPOnrampPackageID),
								ModuleId:  strPtr("onramp"),
								Function:  strPtr("ccip_send"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "ref",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "onramp_state",
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
										Name:     "dest_chain_selector",
										Type:     "u64",
										Required: true,
									},
									{
										Name:     "receiver",
										Type:     "vector<u8>",
										Required: true,
									},
									{
										Name:     "data",
										Type:     "vector<u8>",
										Required: true,
									},
									{
										Name:     "token_params",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 1,
										},
									},
									{
										Name:      "fee_token_metadata",
										Type:      "object_id",
										Required:  true,
										IsMutable: testutils.BoolPointer(false),
									},
									{
										Name:      "fee_token",
										Type:      "object_id",
										Required:  true,
										IsGeneric: true,
									},
									{
										Name:     "extra_args",
										Type:     "vector<u8>",
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// createCCIPSendPTBArgs creates PTBArgMapping for a CCIP send operation
func createCCIPSendPForTwoTokensTBArgs(addresses ContractAddresses) chainwriter.Arguments {
	// Define a destination chain selector (e.g., Ethereum Sepolia)
	destChainSelector := uint64(2)
	linkTokenTypeTag := "0xe3c005c4195ec60a3468ce01238df650e4fedbd36e517bf75b9d2ee90cce8a8b::link_token::LINK_TOKEN"

	return chainwriter.Arguments{
		Args: map[string]any{
			"ref":                   addresses.CCIPStateRef,
			"clock":                 addresses.ClockObject,
			"remote_chain_selector": destChainSelector,
			"dest_chain_selector":   destChainSelector,
			"state":                 addresses.LinkLockReleaseTokenPoolState,
			"c":                     addresses.LinkCoinObjects[0],
			"onramp_state":          addresses.CCIPOnrampState,
			"receiver":              []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"data":                  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"fee_token_metadata":    addresses.LinkTokenCoinMetadata,
			"fee_token":             addresses.LinkCoinObjects[1],
			"extra_args":            []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		ArgTypes: map[string]string{
			"c":         linkTokenTypeTag,
			"fee_token": linkTokenTypeTag,
		},
	}
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}
