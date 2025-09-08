//go:build integration

package ccipops

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/ops/link"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"

	"github.com/stretchr/testify/require"
)

func TestUpgradeRegistryOperations(t *testing.T) {
	t.Parallel()

	signer, client := testenv.SetupEnvironment(t)

	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(500_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)

	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	// Deploy MCMS
	mcmsReport, err := cld_ops.ExecuteOperation(bundle, mcmsops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Contract")

	configDigestHex := "e3b1c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	configDigest, err := hex.DecodeString(configDigestHex)
	require.NoError(t, err, "failed to decode config digest")

	publicKey1Hex := "8a1b2c3d4e5f60718293a4b5c6d7e8f901234567"
	publicKey1, err := hex.DecodeString(publicKey1Hex)
	require.NoError(t, err, "failed to decode public key 1")

	publicKey2Hex := "7b8c9dab0c1d2e3f405162738495a6b7c8d9e0f1"
	publicKey2, err := hex.DecodeString(publicKey2Hex)
	require.NoError(t, err, "failed to decode public key 2")

	publicKey3Hex := "1234567890abcdef1234567890abcdef12345678"
	publicKey3, err := hex.DecodeString(publicKey3Hex)
	require.NoError(t, err, "failed to decode public key 3")

	publicKey4Hex := "90abcdef1234567890abcdef1234567890abcdef"
	publicKey4, err := hex.DecodeString(publicKey4Hex)
	require.NoError(t, err, "failed to decode public key 4")

	signerAddress, err := signer.GetAddress()
	require.NoError(t, err, "failed to get signer address")

	// Deploy and initialize CCIP with upgrade registry
	report, err := cld_ops.ExecuteSequence(bundle, DeployAndInitCCIPSequence, deps, DeployAndInitCCIPSeqInput{
		LinkTokenCoinMetadataObjectId: linkReport.Output.Objects.CoinMetadataObjectId,
		LocalChainSelector:            1,
		DestChainSelector:             2,
		DeployCCIPInput: DeployCCIPInput{
			McmsPackageId: mcmsReport.Output.PackageId,
			McmsOwner:     signerAddress,
		},
		MaxFeeJuelsPerMsg:            "100000000",
		TokenPriceStalenessThreshold: 60,
		// Fee Quoter configuration
		AddMinFeeUsdCents:    []uint32{3000},
		AddMaxFeeUsdCents:    []uint32{30000},
		AddDeciBps:           []uint16{1000},
		AddDestGasOverhead:   []uint32{1000000},
		AddDestBytesOverhead: []uint32{1000},
		AddIsEnabled:         []bool{true},
		RemoveTokens:         []string{},
		// Fee Quoter destination chain configuration
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           2,
		MaxDataBytes:                      2000,
		MaxPerMsgGasLimit:                 5000000,
		DestGasOverhead:                   1000000,
		DestGasPerPayloadByteBase:         byte(2),
		DestGasPerPayloadByteHigh:         byte(5),
		DestGasPerPayloadByteThreshold:    uint16(10),
		DestDataAvailabilityOverheadGas:   300000,
		DestGasPerDataAvailabilityByte:    4,
		DestDataAvailabilityMultiplierBps: 1,
		ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c},
		EnforceOutOfOrder:                 false,
		DefaultTokenFeeUsdCents:           3,
		DefaultTokenDestGasOverhead:       100000,
		DefaultTxGasLimit:                 500000,
		GasMultiplierWeiPerEth:            100,
		GasPriceStalenessThreshold:        1000000000,
		NetworkFeeUsdCents:                10,
		// Premium multiplier updates
		PremiumMultiplierWeiPerEth: []uint64{10},

		RmnHomeContractConfigDigest: configDigest,
		SignerOnchainPublicKeys:     [][]byte{publicKey1, publicKey2, publicKey3, publicKey4},
		NodeIndexes:                 []uint64{0, 1, 2, 3},
		FSign:                       uint64(1),
	})
	require.NoError(t, err, "failed to execute CCIP deploy sequence")
	require.NotEmpty(t, report.Output.CCIPPackageId, "CCIP package ID should not be empty")
	require.NotEmpty(t, report.Output.Objects.UpgradeRegistryObjectId, "UpgradeRegistry object ID should not be empty")

	t.Run("Test Function Restrictions", func(t *testing.T) {
		// Test updating function restrictions
		_, err := cld_ops.ExecuteOperation(bundle, UpdateFunctionRestrictionsOp, deps, UpdateFunctionRestrictionsInput{
			CCIPPackageId:    report.Output.CCIPPackageId,
			StateObjectId:    report.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId: report.Output.Objects.OwnerCapObjectId,
			ModuleName:       "test_module",
			FunctionName:     "test_function",
			BlockedVersions:  []uint64{1, 2, 3},
		})
		require.NoError(t, err, "failed to update function restrictions")

		// Test getting function restrictions
		getFuncRestrictionsReport, err := cld_ops.ExecuteOperation(bundle, GetFunctionRestrictionsOp, deps, GetFunctionRestrictionsInput{
			CCIPPackageId: report.Output.CCIPPackageId,
			StateObjectId: report.Output.Objects.CCIPObjectRefObjectId,
			ModuleName:    "test_module",
			FunctionName:  "test_function",
		})
		require.NoError(t, err, "failed to get function restrictions")
		require.Equal(t, []uint64{1, 2, 3}, getFuncRestrictionsReport.Output.Objects.BlockedVersions, "blocked versions should match")

		// Test checking if function is allowed
		isFuncAllowedReport, err := cld_ops.ExecuteOperation(bundle, IsFunctionAllowedOp, deps, IsFunctionAllowedInput{
			CCIPPackageId:   report.Output.CCIPPackageId,
			StateObjectId:   report.Output.Objects.CCIPObjectRefObjectId,
			ModuleName:      "test_module",
			FunctionName:    "test_function",
			ContractVersion: 1,
		})
		require.NoError(t, err, "failed to check if function is allowed")
		require.False(t, isFuncAllowedReport.Output.Objects.IsAllowed, "function version 1 should be blocked")

		// Test with allowed version
		isFuncAllowedReport2, err := cld_ops.ExecuteOperation(bundle, IsFunctionAllowedOp, deps, IsFunctionAllowedInput{
			CCIPPackageId:   report.Output.CCIPPackageId,
			StateObjectId:   report.Output.Objects.CCIPObjectRefObjectId,
			ModuleName:      "test_module",
			FunctionName:    "test_function",
			ContractVersion: 4,
		})
		require.NoError(t, err, "failed to check if function is allowed")
		require.True(t, isFuncAllowedReport2.Output.Objects.IsAllowed, "function version 4 should be allowed")
	})

	t.Run("Test Module Restrictions", func(t *testing.T) {
		// Test updating module restrictions
		_, err := cld_ops.ExecuteOperation(bundle, UpdateModuleRestrictionsOp, deps, UpdateModuleRestrictionsInput{
			CCIPPackageId:    report.Output.CCIPPackageId,
			StateObjectId:    report.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId: report.Output.Objects.OwnerCapObjectId,
			ModuleName:       "test_module_2",
			BlockedVersions:  []uint64{5, 6},
		})
		require.NoError(t, err, "failed to update module restrictions")

		// Test getting module restrictions
		getModuleRestrictionsReport, err := cld_ops.ExecuteOperation(bundle, GetModuleRestrictionsOp, deps, GetModuleRestrictionsInput{
			CCIPPackageId: report.Output.CCIPPackageId,
			StateObjectId: report.Output.Objects.CCIPObjectRefObjectId,
			ModuleName:    "test_module_2",
		})
		require.NoError(t, err, "failed to get module restrictions")
		require.Equal(t, []uint64{5, 6}, getModuleRestrictionsReport.Output.Objects.BlockedVersions, "blocked versions should match")

		// Test checking if module is allowed
		isModuleAllowedReport, err := cld_ops.ExecuteOperation(bundle, IsModuleAllowedOp, deps, IsModuleAllowedInput{
			CCIPPackageId:   report.Output.CCIPPackageId,
			StateObjectId:   report.Output.Objects.CCIPObjectRefObjectId,
			ModuleName:      "test_module_2",
			ContractVersion: 5,
		})
		require.NoError(t, err, "failed to check if module is allowed")
		require.False(t, isModuleAllowedReport.Output.Objects.IsAllowed, "module version 5 should be blocked")

		// Test with allowed version
		isModuleAllowedReport2, err := cld_ops.ExecuteOperation(bundle, IsModuleAllowedOp, deps, IsModuleAllowedInput{
			CCIPPackageId:   report.Output.CCIPPackageId,
			StateObjectId:   report.Output.Objects.CCIPObjectRefObjectId,
			ModuleName:      "test_module_2",
			ContractVersion: 7,
		})
		require.NoError(t, err, "failed to check if module is allowed")
		require.True(t, isModuleAllowedReport2.Output.Objects.IsAllowed, "module version 7 should be allowed")
	})
}
