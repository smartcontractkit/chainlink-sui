//go:build integration

package offramp_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/mcms"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// TestBrokenReceiverABI_NoPanic verifies that the relayer's DecodeParameters handles
// a registered receiver with generic type parameters (TypeParameter in normalized ABI)
// without panicking. This is the core fix for report 71024.
func TestBrokenReceiverABI_NoPanic(t *testing.T) {
	lggr := logger.Test(t)
	gasBudget := int64(1_000_000_000)

	// Start dedicated Sui node
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			if perr := cmd.Process.Kill(); perr != nil {
				t.Logf("Failed to kill Sui node process: %v", perr)
			}
		}
	})

	time.Sleep(3 * time.Second)

	// Setup signer and fund account
	keystoreInstance, accountAddress, publicKeyBytes := testutils.SetupTestSigner(t, context.Background(), lggr, gasBudget)
	lggr.Infow("Using account", "address", accountAddress)

	for range 3 {
		fundErr := testutils.FundWithFaucet(lggr, "localnet", accountAddress)
		require.NoError(t, fundErr)
	}

	// Setup client
	ptbClient, _, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr, gasBudget)
	suiClient := ptbClient.GetClient()

	signer := keystoreInstance.GetSuiSigner(context.Background(), fmt.Sprintf("%064x", publicKeyBytes))
	privateKeySigner := rel.NewPrivateKeySigner(signer.PriKey)

	gasBudgetU64 := uint64(gasBudget)
	opts := &bind.CallOpts{
		Signer:           privateKeySigner,
		WaitForExecution: true,
		GasBudget:        &gasBudgetU64,
	}

	mcmsPkg, _, err := mcms.PublishMCMS(context.Background(), opts, suiClient, testutils.LocalUrl)
	require.NoError(t, err, "failed to publish MCMS")
	mcmsPackageId := mcmsPkg.Address()

	ccipPkg, _, err := ccip.PublishCCIP(context.Background(), opts, suiClient, mcmsPackageId, "0x2", testutils.LocalUrl)
	require.NoError(t, err, "failed to publish CCIP")
	ccipPackageId := ccipPkg.Address()

	signerAddr, err := privateKeySigner.GetAddress()
	require.NoError(t, err)

	brokenReceiverArtifact, err := bind.CompilePackage(contracts.CCIPBrokenReceiver, map[string]string{
		"ccip":                 ccipPackageId,
		"ccip_broken_receiver": "0x0",
		"mcms":                 mcmsPackageId,
		"mcms_owner":           "0x2",
		"signer":               signerAddr,
	}, false, testutils.LocalUrl)
	require.NoError(t, err, "failed to compile broken receiver")

	brokenReceiverPackageId, _, err := bind.PublishPackage(context.Background(), opts, suiClient, bind.PublishRequest{
		CompiledModules: brokenReceiverArtifact.Modules,
		Dependencies:    brokenReceiverArtifact.Dependencies,
	})
	require.NoError(t, err, "failed to publish broken receiver")
	lggr.Infow("Published broken receiver", "packageId", brokenReceiverPackageId)

	// Fetch the normalized module — this is the data the relayer would get from chain
	normalizedModule, err := suiClient.SuiGetNormalizedMoveModule(context.Background(), models.GetNormalizedMoveModuleRequest{
		Package:    brokenReceiverPackageId,
		ModuleName: "broken_receiver",
	})
	require.NoError(t, err, "failed to get normalized module")

	functionSig, ok := normalizedModule.ExposedFunctions["ccip_receive"]
	require.True(t, ok, "ccip_receive not found in normalized module")

	funcMap, ok := functionSig.(map[string]any)
	require.True(t, ok, "function signature is not a map")

	params, ok := funcMap["parameters"].([]any)
	require.True(t, ok, "parameters field is not an array")
	require.Greater(t, len(params), 0)

	lggr.Infow("Broken receiver normalized ABI", "parameters", params)

	// The critical assertion: DecodeParameters must not panic and must return an error.
	// Before the fix this panicked with:
	//   interface conversion: interface {} is float64, not map[string]interface {}
	assert.NotPanics(t, func() {
		result, decodeErr := offramp.DecodeParameters(lggr, funcMap, "parameters")
		assert.Error(t, decodeErr, "DecodeParameters should return error for TypeParameter ABI")
		assert.Nil(t, result)
		assert.Contains(t, decodeErr.Error(), "TypeParameter")
	})
}

// TestValidReceiverABI_DecodesSuccessfully verifies that the fix does not break
// legitimate receivers. The dummy receiver's ccip_receive has a concrete signature
// (no generics) and DecodeParameters must succeed on its normalized ABI.
func TestValidReceiverABI_DecodesSuccessfully(t *testing.T) {
	lggr := logger.Test(t)
	gasBudget := int64(1_000_000_000)

	// Start dedicated Sui node
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			if perr := cmd.Process.Kill(); perr != nil {
				t.Logf("Failed to kill Sui node process: %v", perr)
			}
		}
	})

	time.Sleep(3 * time.Second)

	keystoreInstance, accountAddress, publicKeyBytes := testutils.SetupTestSigner(t, context.Background(), lggr, gasBudget)
	lggr.Infow("Using account", "address", accountAddress)

	for range 3 {
		fundErr := testutils.FundWithFaucet(lggr, "localnet", accountAddress)
		require.NoError(t, fundErr)
	}

	ptbClient, _, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr, gasBudget)
	suiClient := ptbClient.GetClient()

	signer := keystoreInstance.GetSuiSigner(context.Background(), fmt.Sprintf("%064x", publicKeyBytes))
	privateKeySigner := rel.NewPrivateKeySigner(signer.PriKey)

	gasBudgetU64 := uint64(gasBudget)
	opts := &bind.CallOpts{
		Signer:           privateKeySigner,
		WaitForExecution: true,
		GasBudget:        &gasBudgetU64,
	}

	mcmsPkg, _, err := mcms.PublishMCMS(context.Background(), opts, suiClient, testutils.LocalUrl)
	require.NoError(t, err, "failed to publish MCMS")
	mcmsPackageId := mcmsPkg.Address()

	ccipPkg, _, err := ccip.PublishCCIP(context.Background(), opts, suiClient, mcmsPackageId, "0x2", testutils.LocalUrl)
	require.NoError(t, err, "failed to publish CCIP")
	ccipPackageId := ccipPkg.Address()

	signerAddr, err := privateKeySigner.GetAddress()
	require.NoError(t, err)

	dummyReceiverArtifact, err := bind.CompilePackage(contracts.CCIPDummyReceiver, map[string]string{
		"ccip":                ccipPackageId,
		"ccip_dummy_receiver": "0x0",
		"mcms":                mcmsPackageId,
		"mcms_owner":          "0x2",
		"signer":              signerAddr,
	}, false, testutils.LocalUrl)
	require.NoError(t, err, "failed to compile dummy receiver")

	dummyReceiverPackageId, _, err := bind.PublishPackage(context.Background(), opts, suiClient, bind.PublishRequest{
		CompiledModules: dummyReceiverArtifact.Modules,
		Dependencies:    dummyReceiverArtifact.Dependencies,
	})
	require.NoError(t, err, "failed to publish dummy receiver")
	lggr.Infow("Published dummy receiver", "packageId", dummyReceiverPackageId)

	// Fetch the normalized module
	normalizedModule, err := suiClient.SuiGetNormalizedMoveModule(context.Background(), models.GetNormalizedMoveModuleRequest{
		Package:    dummyReceiverPackageId,
		ModuleName: "dummy_receiver",
	})
	require.NoError(t, err)

	functionSig, ok := normalizedModule.ExposedFunctions["ccip_receive"]
	require.True(t, ok, "ccip_receive not found")

	funcMap, ok := functionSig.(map[string]any)
	require.True(t, ok)

	// DecodeParameters must succeed for valid receivers
	result, err := offramp.DecodeParameters(lggr, funcMap, "parameters")
	require.NoError(t, err, "DecodeParameters should succeed for dummy receiver ABI")
	require.NotNil(t, result)
	require.Greater(t, len(result), 0, "should decode at least one parameter type")

	lggr.Infow("Dummy receiver decoded parameter types", "paramTypes", result)
}
