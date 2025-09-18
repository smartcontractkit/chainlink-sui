//go:build integration

package ccipops

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"

	"github.com/stretchr/testify/require"
)

func TestFeeQuoterOperations(t *testing.T) {
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
	require.NoError(t, err, "failed to deploy MCMS")

	// Deploy CCIP
	deployReport, err := cld_ops.ExecuteOperation(bundle, DeployCCIPOp, deps, DeployCCIPInput{
		McmsPackageId: mcmsReport.Output.PackageId,
		McmsOwner:     "0x1234567890abcdef1234567890abcdef12345678", // dummy address
	})
	require.NoError(t, err, "failed to deploy CCIP")

	// Initialize Fee Quoter
	initFQReport, err := cld_ops.ExecuteOperation(
		bundle,
		FeeQuoterInitializeOp,
		deps,
		InitFeeQuoterInput{
			CCIPPackageId:                 deployReport.Output.PackageId,
			StateObjectId:                 deployReport.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId:              deployReport.Output.Objects.OwnerCapObjectId,
			MaxFeeJuelsPerMsg:             "100000000",
			LinkTokenCoinMetadataObjectId: linkReport.Output.Objects.CoinMetadataObjectId,
			TokenPriceStalenessThreshold:  60,
			FeeTokens:                     []string{linkReport.Output.Objects.CoinMetadataObjectId},
		},
	)
	require.NoError(t, err, "failed to initialize fee quoter")
	require.NotEmpty(t, initFQReport.Output.Objects.FeeQuoterCapObjectId, "fee quoter cap should be created during initialization")
	require.NotEmpty(t, initFQReport.Output.Objects.FeeQuoterStateObjectId, "fee quoter state should be created during initialization")

	// Test ApplyFeeTokenUpdates operation
	t.Run("ApplyFeeTokenUpdates", func(t *testing.T) {
		// Add a new fee token (using LINK as an example)
		applyFeeTokenUpdatesReport, err := cld_ops.ExecuteOperation(
			bundle,
			FeeQuoterApplyFeeTokenUpdatesOp,
			deps,
			FeeQuoterApplyFeeTokenUpdatesInput{
				CCIPPackageId:     deployReport.Output.PackageId,
				StateObjectId:     deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:  deployReport.Output.Objects.OwnerCapObjectId,
				FeeTokensToRemove: []string{},                                               // Remove none
				FeeTokensToAdd:    []string{linkReport.Output.Objects.CoinMetadataObjectId}, // Add LINK token
			},
		)
		require.NoError(t, err, "failed to apply fee token updates")
		require.NotEmpty(t, applyFeeTokenUpdatesReport.Digest, "apply fee token updates transaction should have a digest")
	})

	// Test ApplyTokenTransferFeeConfigUpdates operation
	t.Run("ApplyTokenTransferFeeConfigUpdates", func(t *testing.T) {
		applyTokenTransferFeeConfigReport, err := cld_ops.ExecuteOperation(
			bundle,
			FeeQuoterApplyTokenTransferFeeConfigUpdatesOp,
			deps,
			FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
				CCIPPackageId:        deployReport.Output.PackageId,
				StateObjectId:        deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:     deployReport.Output.Objects.OwnerCapObjectId,
				DestChainSelector:    1, // Test destination chain
				AddTokens:            []string{linkReport.Output.Objects.CoinMetadataObjectId},
				AddMinFeeUsdCents:    []uint32{3000},  // $0.30 minimum fee
				AddMaxFeeUsdCents:    []uint32{30000}, // $3.00 maximum fee
				AddDeciBps:           []uint16{1000},  // 0.1% fee
				AddDestGasOverhead:   []uint32{1000000},
				AddDestBytesOverhead: []uint32{1000},
				AddIsEnabled:         []bool{true},
				RemoveTokens:         []string{},
			},
		)
		require.NoError(t, err, "failed to apply token transfer fee config updates")
		require.NotEmpty(t, applyTokenTransferFeeConfigReport.Digest, "apply token transfer fee config updates transaction should have a digest")
	})

	// Test ApplyDestChainConfigUpdates operation
	t.Run("ApplyDestChainConfigUpdates", func(t *testing.T) {
		applyDestChainConfigReport, err := cld_ops.ExecuteOperation(
			bundle,
			FeeQuoterApplyDestChainConfigUpdatesOp,
			deps,
			FeeQuoterApplyDestChainConfigUpdatesInput{
				CCIPPackageId:                     deployReport.Output.PackageId,
				StateObjectId:                     deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:                  deployReport.Output.Objects.OwnerCapObjectId,
				DestChainSelector:                 1,
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
				ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c}, // EVM chain family
				EnforceOutOfOrder:                 false,
				DefaultTokenFeeUsdCents:           3,
				DefaultTokenDestGasOverhead:       100000,
				DefaultTxGasLimit:                 500000,
				GasMultiplierWeiPerEth:            100,
				GasPriceStalenessThreshold:        300,
				NetworkFeeUsdCents:                5,
			},
		)
		require.NoError(t, err, "failed to apply dest chain config updates")
		require.NotEmpty(t, applyDestChainConfigReport.Digest, "apply dest chain config updates transaction should have a digest")
	})

	// Test ApplyPremiumMultiplierWeiPerEthUpdates operation
	t.Run("ApplyPremiumMultiplierWeiPerEthUpdates", func(t *testing.T) {
		applyPremiumMultiplierReport, err := cld_ops.ExecuteOperation(
			bundle,
			FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp,
			deps,
			FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
				CCIPPackageId:              deployReport.Output.PackageId,
				StateObjectId:              deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:           deployReport.Output.Objects.OwnerCapObjectId,
				Tokens:                     []string{linkReport.Output.Objects.CoinMetadataObjectId},
				PremiumMultiplierWeiPerEth: []uint64{1100000000000000000}, // 10% premium (1.1e18)
			},
		)
		require.NoError(t, err, "failed to apply premium multiplier updates")
		require.NotEmpty(t, applyPremiumMultiplierReport.Digest, "apply premium multiplier updates transaction should have a digest")
	})

	// Test UpdatePrices operation (using fee quoter cap)
	t.Run("UpdatePrices", func(t *testing.T) {
		// Create some test price data
		sourceTokens := []string{linkReport.Output.Objects.CoinMetadataObjectId}
		sourceUsdPerToken := []*big.Int{big.NewInt(1000000000000000000)} // 1 USD in wei
		gasDestChainSelectors := []uint64{1}
		gasUsdPerUnitGas := []*big.Int{big.NewInt(2000000000000000)} // 0.002 USD per unit gas

		updatePricesReport, err := cld_ops.ExecuteOperation(
			bundle,
			FeeQuoterUpdateTokenPricesOp,
			deps,
			FeeQuoterUpdateTokenPricesInput{
				CCIPPackageId:         deployReport.Output.PackageId,
				CCIPObjectRef:         deployReport.Output.Objects.CCIPObjectRefObjectId,
				FeeQuoterCapId:        initFQReport.Output.Objects.FeeQuoterCapObjectId,
				SourceTokens:          sourceTokens,
				SourceUsdPerToken:     sourceUsdPerToken,
				GasDestChainSelectors: gasDestChainSelectors,
				GasUsdPerUnitGas:      gasUsdPerUnitGas,
			},
		)
		require.NoError(t, err, "failed to update prices")
		require.NotEmpty(t, updatePricesReport.Digest, "update prices transaction should have a digest")
	})

	// Test NewFeeQuoterCap operation
	t.Run("NewFeeQuoterCap", func(t *testing.T) {
		newCapReport, err := cld_ops.ExecuteOperation(
			bundle,
			FeeQuoterNewFeeQuoterCapOp,
			deps,
			NewFeeQuoterCapInput{
				CCIPPackageId:    deployReport.Output.PackageId,
				OwnerCapObjectId: deployReport.Output.Objects.OwnerCapObjectId,
			},
		)
		require.NoError(t, err, "failed to create new fee quoter cap")
		require.NotEmpty(t, newCapReport.Output.Objects.FeeQuoterCapObjectId, "fee quoter cap object ID should not be empty")

		// Test UpdatePricesWithOwnerCap operation using the new cap
		t.Run("UpdatePricesWithOwnerCap", func(t *testing.T) {
			// Create some test price data
			sourceTokens := []string{linkReport.Output.Objects.CoinMetadataObjectId}
			sourceUsdPerToken := []*big.Int{big.NewInt(1000000000000000000)} // 1 USD in wei
			gasDestChainSelectors := []uint64{1}
			gasUsdPerUnitGas := []*big.Int{big.NewInt(2000000000000000)} // 0.002 USD per unit gas

			updatePricesWithOwnerCapReport, err := cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterUpdatePricesWithOwnerCapOp,
				deps,
				FeeQuoterUpdatePricesWithOwnerCapInput{
					CCIPPackageId:         deployReport.Output.PackageId,
					CCIPObjectRef:         deployReport.Output.Objects.CCIPObjectRefObjectId,
					OwnerCapObjectId:      deployReport.Output.Objects.OwnerCapObjectId,
					SourceTokens:          sourceTokens,
					SourceUsdPerToken:     sourceUsdPerToken,
					GasDestChainSelectors: gasDestChainSelectors,
					GasUsdPerUnitGas:      gasUsdPerUnitGas,
				},
			)
			require.NoError(t, err, "failed to update prices with owner cap")
			require.NotEmpty(t, updatePricesWithOwnerCapReport.Digest, "update prices with owner cap transaction should have a digest")
		})

		// Test DestroyFeeQuoterCap operation
		t.Run("DestroyFeeQuoterCap", func(t *testing.T) {
			destroyCapReport, err := cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterDestroyFeeQuoterCapOp,
				deps,
				DestroyFeeQuoterCapInput{
					CCIPPackageId:        deployReport.Output.PackageId,
					OwnerCapObjectId:     deployReport.Output.Objects.OwnerCapObjectId,
					FeeQuoterCapObjectId: newCapReport.Output.Objects.FeeQuoterCapObjectId,
				},
			)
			require.NoError(t, err, "failed to destroy fee quoter cap")
			require.NotEmpty(t, destroyCapReport.Digest, "destroy transaction should have a digest")
		})
	})

	// Test error cases
	t.Run("ErrorCases", func(t *testing.T) {
		// Test with invalid package ID
		t.Run("InvalidPackageID", func(t *testing.T) {
			_, err := cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterApplyFeeTokenUpdatesOp,
				deps,
				FeeQuoterApplyFeeTokenUpdatesInput{
					CCIPPackageId:     "invalid_package_id",
					StateObjectId:     deployReport.Output.Objects.CCIPObjectRefObjectId,
					OwnerCapObjectId:  deployReport.Output.Objects.OwnerCapObjectId,
					FeeTokensToRemove: []string{},
					FeeTokensToAdd:    []string{},
				},
			)
			require.Error(t, err, "should fail with invalid package ID")
		})

		// Test with invalid object IDs
		t.Run("InvalidObjectIDs", func(t *testing.T) {
			_, err := cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterApplyFeeTokenUpdatesOp,
				deps,
				FeeQuoterApplyFeeTokenUpdatesInput{
					CCIPPackageId:     deployReport.Output.PackageId,
					StateObjectId:     "invalid_state_id",
					OwnerCapObjectId:  "invalid_owner_cap_id",
					FeeTokensToRemove: []string{},
					FeeTokensToAdd:    []string{},
				},
			)
			require.Error(t, err, "should fail with invalid object IDs")
		})
	})

	// Test complex scenarios
	t.Run("ComplexScenarios", func(t *testing.T) {
		// Test multiple fee token updates
		t.Run("MultipleFeeTokenUpdates", func(t *testing.T) {
			// First add a token
			_, err := cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterApplyFeeTokenUpdatesOp,
				deps,
				FeeQuoterApplyFeeTokenUpdatesInput{
					CCIPPackageId:     deployReport.Output.PackageId,
					StateObjectId:     deployReport.Output.Objects.CCIPObjectRefObjectId,
					OwnerCapObjectId:  deployReport.Output.Objects.OwnerCapObjectId,
					FeeTokensToRemove: []string{},
					FeeTokensToAdd:    []string{linkReport.Output.Objects.CoinMetadataObjectId},
				},
			)
			require.NoError(t, err, "failed to add fee token")

			// Then remove it
			_, err = cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterApplyFeeTokenUpdatesOp,
				deps,
				FeeQuoterApplyFeeTokenUpdatesInput{
					CCIPPackageId:     deployReport.Output.PackageId,
					StateObjectId:     deployReport.Output.Objects.CCIPObjectRefObjectId,
					OwnerCapObjectId:  deployReport.Output.Objects.OwnerCapObjectId,
					FeeTokensToRemove: []string{linkReport.Output.Objects.CoinMetadataObjectId},
					FeeTokensToAdd:    []string{},
				},
			)
			require.NoError(t, err, "failed to remove fee token")
		})

		// Test multiple token transfer fee config updates
		t.Run("MultipleTokenTransferFeeConfigUpdates", func(t *testing.T) {
			// Add multiple token configs
			_, err := cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterApplyTokenTransferFeeConfigUpdatesOp,
				deps,
				FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
					CCIPPackageId:        deployReport.Output.PackageId,
					StateObjectId:        deployReport.Output.Objects.CCIPObjectRefObjectId,
					OwnerCapObjectId:     deployReport.Output.Objects.OwnerCapObjectId,
					DestChainSelector:    1,
					AddTokens:            []string{linkReport.Output.Objects.CoinMetadataObjectId},
					AddMinFeeUsdCents:    []uint32{1000},  // $0.10
					AddMaxFeeUsdCents:    []uint32{10000}, // $1.00
					AddDeciBps:           []uint16{500},   // 0.05%
					AddDestGasOverhead:   []uint32{500000},
					AddDestBytesOverhead: []uint32{500},
					AddIsEnabled:         []bool{true},
					RemoveTokens:         []string{},
				},
			)
			require.NoError(t, err, "failed to add token transfer fee config")

			// Remove the config
			_, err = cld_ops.ExecuteOperation(
				bundle,
				FeeQuoterApplyTokenTransferFeeConfigUpdatesOp,
				deps,
				FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
					CCIPPackageId:        deployReport.Output.PackageId,
					StateObjectId:        deployReport.Output.Objects.CCIPObjectRefObjectId,
					OwnerCapObjectId:     deployReport.Output.Objects.OwnerCapObjectId,
					DestChainSelector:    1,
					AddTokens:            []string{},
					AddMinFeeUsdCents:    []uint32{},
					AddMaxFeeUsdCents:    []uint32{},
					AddDeciBps:           []uint16{},
					AddDestGasOverhead:   []uint32{},
					AddDestBytesOverhead: []uint32{},
					AddIsEnabled:         []bool{},
					RemoveTokens:         []string{linkReport.Output.Objects.CoinMetadataObjectId},
				},
			)
			require.NoError(t, err, "failed to remove token transfer fee config")
		})
	})
}
