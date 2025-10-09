//go:build integration

package mcms

import (
	"encoding/hex"
	"testing"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	"github.com/stretchr/testify/require"
)

type CCIPMCMSTestSuite struct {
	MCMSTestSuite

	ccipPackageId string
	ccipObjects   ccipops.DeployCCIPSeqObjects
}

func (s *CCIPMCMSTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()

	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(s.bundle, linkops.DeployLINKOp, s.deps, cld_ops.EmptyInput{})
	require.NoError(s.T(), err, "failed to deploy LINK token")

	configDigestHex := "e3b1c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	configDigest, err := hex.DecodeString(configDigestHex)
	require.NoError(s.T(), err, "failed to decode config digest")

	publicKey1Hex := "8a1b2c3d4e5f60718293a4b5c6d7e8f901234567"
	publicKey1, err := hex.DecodeString(publicKey1Hex)
	require.NoError(s.T(), err, "failed to decode public key 1")

	publicKey2Hex := "7b8c9dab0c1d2e3f405162738495a6b7c8d9e0f1"
	publicKey2, err := hex.DecodeString(publicKey2Hex)
	require.NoError(s.T(), err, "failed to decode public key 2")

	publicKey3Hex := "1234567890abcdef1234567890abcdef12345678"
	publicKey3, err := hex.DecodeString(publicKey3Hex)
	require.NoError(s.T(), err, "failed to decode public key 3")

	publicKey4Hex := "90abcdef1234567890abcdef1234567890abcdef"
	publicKey4, err := hex.DecodeString(publicKey4Hex)
	require.NoError(s.T(), err, "failed to decode public key 4")

	// Use the same seq as in production deployment
	report, err := cld_ops.ExecuteSequence(s.bundle, ccipops.DeployAndInitCCIPSequence, s.deps, ccipops.DeployAndInitCCIPSeqInput{
		LinkTokenCoinMetadataObjectId: linkReport.Output.Objects.CoinMetadataObjectId,
		LocalChainSelector:            1,
		DestChainSelector:             2,
		DeployCCIPInput: ccipops.DeployCCIPInput{
			McmsPackageId: s.mcmsPackageID,
			McmsOwner:     s.mcmsOwnerAddress,
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
	require.NoError(s.T(), err, "failed to execute CCIP deploy sequence")
	require.NotEmpty(s.T(), report.Output.CCIPPackageId, "CCIP package ID should not be empty")

	s.ccipPackageId = report.Output.CCIPPackageId
	s.ccipObjects = report.Output.Objects
}

func (s *CCIPMCMSTestSuite) Test_CCIP_MCMS() {
	s.T().Run("Transfer Ownership of CCIP to MCMS", func(t *testing.T) {
		RunTestCCIPOwnershipTransfer(s)
	})

	s.T().Run("Execute config proposal against CCIP from MCMS", func(t *testing.T) {
		RunTestCCIPProposal(s)
	})
}

// TODO: For prod env, the initial deployment sequence should start the ownership transfer flow of every deployed contract
func RunTestCCIPOwnershipTransfer(s *CCIPMCMSTestSuite) {
	// 1. Transfer OwnerCap of CCIP to MCMS
	ccipContract, err := module_state_object.NewStateObject(s.ccipPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip state object contract")

	tx, err := ccipContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId},
		bind.Object{Id: s.ccipObjects.OwnerCapObjectId},
		// TODO: not sure which address should be
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	s.T().Logf("✅ Transferred ownership of CCIP to MCMS in tx: %s", tx.Digest)

	// 2. Proposal execution with acceptance from MCMS (through proposer)
	// acceptOwnershipProposal, err := mcmsops.


	// 3. Execute transfer ownership from original owner
}

func RunTestCCIPProposal(s *CCIPMCMSTestSuite) {

}
