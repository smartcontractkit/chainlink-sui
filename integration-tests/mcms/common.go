//go:build integration

package mcms

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	mcmsencoder "github.com/smartcontractkit/chainlink-sui/bindings"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	module_user "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_user"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"

	cselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

type RoleConfig struct {
	Role   suisdk.TimelockRole
	Count  int
	Quorum uint8
	Keys   []*ecdsa.PrivateKey
	Config *types.Config
}

func CreateConfig(role suisdk.TimelockRole, count int, quorum uint8) *RoleConfig {
	signers := make([]common.Address, count)
	signerKeys := make([]*ecdsa.PrivateKey, count)

	for i := range signers {
		signerKeys[i], _ = crypto.GenerateKey()
		signers[i] = crypto.PubkeyToAddress(signerKeys[i].PublicKey)
	}
	slices.SortFunc(signers[:], func(a, b common.Address) int {
		return a.Cmp(b)
	})

	return &RoleConfig{
		Role:   role,
		Count:  count,
		Quorum: quorum,
		Keys:   signerKeys,
		Config: &types.Config{
			Quorum:  quorum,
			Signers: signers[:],
		},
	}
}

type MCMSTestSuite struct {
	suite.Suite

	client sui.ISuiAPI
	signer bindutils.SuiSigner

	chainSelector types.ChainSelector

	// MCMS
	mcmsPackageID    string
	mcmsOwnerAddress string
	mcmsObj          string
	timelockObj      string
	registryObj      string
	deployerStateObj string
	accountObj       string
	ownerCapObj      string

	bypasserConfig *RoleConfig
	proposerConfig *RoleConfig

	// Ops
	deps   sui_ops.OpTxDeps
	bundle cld_ops.Bundle

	// CCIP
	// LINK
	linkPackageId string
	linkObjects   linkops.DeployLinkObjects

	// CCIP
	ccipPackageId       string
	latestCcipPackageId string
	ccipObjects         ccipops.DeployCCIPSeqObjects

	// Router
	ccipRouterPackageId       string
	latestCcipRouterPackageId string
	ccipRouterObjects         routerops.DeployCCIPRouterObjects

	// Onramp
	ccipOnrampPackageId       string
	latestCcipOnrampPackageId string
	ccipOnrampObjects         onrampops.DeployCCIPOnRampSeqObjects

	// offramp
	ccipOfframpPackageId       string
	latestCcipOfframpPackageId string
	ccipOfframpObjects         offrampops.DeployCCIPOffRampSeqObjects
}

// TODO: refactor so suites are per product
func (s *MCMSTestSuite) SetupSuite() {
	signer, client := testenv.SetupEnvironment(s.T())
	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(10_000_000_000) // needs to be high for publishing and big proposals
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
				Signer:           signer,
			}
		},
	}

	// Convert slice of values to slice of pointers
	registry := cld_ops.NewOperationRegistry(
		opregistry.AllOperations...,
	)

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		s.T().Context,
		logger.Test(s.T()),
		reporter,
		cld_ops.WithOperationRegistry(registry),
	)

	bypasserCount := 2
	bypasserQuorum := 2
	bypasserConfig := CreateConfig(suisdk.TimelockRoleBypasser, bypasserCount, uint8(bypasserQuorum))
	proposerCount := 3
	proposerQuorum := 2
	proposerConfig := CreateConfig(suisdk.TimelockRoleProposer, proposerCount, uint8(proposerQuorum))

	deployInput := mcmsops.DeployMCMSSeqInput{
		ChainSelector: cselectors.SUI_TESTNET.Selector,
		Bypasser:      bypasserConfig.Config,
		Proposer:      proposerConfig.Config,
	}

	mcmsDeploymentReport, err := cld_ops.ExecuteSequence(bundle, mcmsops.DeployMCMSSequence, deps, deployInput)
	require.NoError(s.T(), err, "deploying MCMS contract")

	s.mcmsPackageID = mcmsDeploymentReport.Output.PackageId
	s.mcmsObj = mcmsDeploymentReport.Output.Objects.McmsMultisigStateObjectId
	s.timelockObj = mcmsDeploymentReport.Output.Objects.TimelockObjectId
	s.registryObj = mcmsDeploymentReport.Output.Objects.McmsRegistryObjectId
	s.deployerStateObj = mcmsDeploymentReport.Output.Objects.McmsDeployerStateObjectId
	s.accountObj = mcmsDeploymentReport.Output.Objects.McmsAccountStateObjectId
	s.ownerCapObj = mcmsDeploymentReport.Output.Objects.McmsAccountOwnerCapObjectId

	s.mcmsOwnerAddress, err = signer.GetAddress()
	require.NoError(s.T(), err, "getting MCMS owner address")

	s.bypasserConfig = bypasserConfig
	s.proposerConfig = proposerConfig

	s.client = client
	s.signer = signer

	s.chainSelector = types.ChainSelector(cselectors.SUI_TESTNET.Selector)

	s.deps = deps
	s.bundle = bundle

	// Accept MCMS ownership to itself
	acceptProposal := mcmsDeploymentReport.Output.AcceptOwnershipProposal
	// Execute the proposal to accept ownership
	s.ExecuteProposalE2e(&acceptProposal, s.proposerConfig, 0)

	rep, err := cld_ops.ExecuteOperation(s.bundle, mcmsops.MCMSExecuteTransferOwnershipOp, s.deps, mcmsops.MCMSExecuteTransferOwnershipInput{
		McmsPackageID:         s.mcmsPackageID,
		OwnerCap:              s.ownerCapObj,
		AccountObjectID:       s.accountObj,
		RegistryObjectID:      s.registryObj,
		DeployerStateObjectID: s.deployerStateObj,
	})
	s.Require().NoError(err, "executing ownership transfer to self")
	s.T().Logf("✅ Transferred ownership of MCMS to itself in tx: %s", rep.Output.Digest)

	s.SetupCCIP()
}

func (s *MCMSTestSuite) NewOpBundle() cld_ops.Bundle {
	reporter := cld_ops.NewMemoryReporter()
	return cld_ops.NewBundle(
		s.T().Context,
		logger.Test(s.T()),
		reporter,
	)
}

func (s *MCMSTestSuite) SetupCCIP() {
	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(s.bundle, linkops.DeployLINKOp, s.deps, cld_ops.EmptyInput{})
	require.NoError(s.T(), err, "failed to deploy LINK token")
	s.linkPackageId = linkReport.Output.PackageId
	s.linkObjects = linkReport.Output.Objects

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
	ccipReport, err := cld_ops.ExecuteSequence(s.bundle, ccipops.DeployAndInitCCIPSequence, s.deps, ccipops.DeployAndInitCCIPSeqInput{
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
	require.NotEmpty(s.T(), ccipReport.Output.CCIPPackageId, "CCIP package ID should not be empty")

	s.linkObjects = linkReport.Output.Objects
	s.ccipPackageId = ccipReport.Output.CCIPPackageId
	s.latestCcipPackageId = ccipReport.Output.CCIPPackageId
	s.ccipObjects = ccipReport.Output.Objects

	// Deploy Router
	routerReport, err := cld_ops.ExecuteOperation(s.bundle, routerops.DeployCCIPRouterOp, s.deps, routerops.DeployCCIPRouterInput{
		McmsPackageId: s.mcmsPackageID,
		McmsOwner:     s.mcmsOwnerAddress,
	})
	require.NoError(s.T(), err, "failed to execute CCIP deploy sequence")

	s.ccipRouterPackageId = routerReport.Output.PackageId
	s.latestCcipRouterPackageId = routerReport.Output.PackageId
	s.ccipRouterObjects = routerReport.Output.Objects

	// Deploy Onramp
	ccipOnRampSeqInput := deployment.DefaultOnRampSeqConfig
	ccipOnRampSeqInput.DeployCCIPOnRampInput.CCIPPackageId = ccipReport.Output.CCIPPackageId
	ccipOnRampSeqInput.DeployCCIPOnRampInput.MCMSPackageId = s.mcmsPackageID
	ccipOnRampSeqInput.DeployCCIPOnRampInput.MCMSOwnerPackageId = s.mcmsOwnerAddress
	ccipOnRampSeqInput.OnRampInitializeInput.NonceManagerCapId = ccipReport.Output.Objects.NonceManagerCapObjectId
	ccipOnRampSeqInput.OnRampInitializeInput.SourceTransferCapId = ccipReport.Output.Objects.SourceTransferCapObjectId
	ccipOnRampSeqInput.OnRampInitializeInput.ChainSelector = uint64(s.chainSelector)
	ccipOnRampSeqInput.OnRampInitializeInput.FeeAggregator = s.mcmsOwnerAddress
	ccipOnRampSeqInput.OnRampInitializeInput.AllowListAdmin = s.mcmsOwnerAddress
	ccipOnRampSeqInput.OnRampInitializeInput.DestChainSelectors = []uint64{cselectors.ETHEREUM_MAINNET.Selector}
	ccipOnRampSeqInput.OnRampInitializeInput.DestChainRouters = []string{routerReport.Output.PackageId}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.DestChainSelector = []uint64{cselectors.ETHEREUM_MAINNET.Selector}
	ccipOnRampSeqInput.ApplyAllowListUpdatesInput.DestChainSelector = []uint64{cselectors.ETHEREUM_MAINNET.Selector}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.DestChainRouters = []string{routerReport.Output.PackageId}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.CCIPObjectRefId = ccipReport.Output.Objects.CCIPObjectRefObjectId

	ccipOnRampSeqReport, err := operations.ExecuteSequence(s.bundle, onrampops.DeployAndInitCCIPOnRampSequence, s.deps, ccipOnRampSeqInput)
	require.NoError(s.T(), err, "failed to execute CCIP OnRamp deploy sequence")

	s.ccipOnrampPackageId = ccipOnRampSeqReport.Output.CCIPOnRampPackageId
	s.latestCcipOnrampPackageId = ccipOnRampSeqReport.Output.CCIPOnRampPackageId
	s.ccipOnrampObjects = ccipOnRampSeqReport.Output.Objects

	// Deploy offramp
	ccipOffRampSeqInput := deployment.DefaultOffRampSeqConfig
	// note: this is a regression, can't acess other chains state very cleanly
	onRampBytes := [][]byte{
		{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	// Inject dynamic values for deployment
	ccipOffRampSeqInput.CCIPObjectRefId = ccipReport.Output.Objects.CCIPObjectRefObjectId
	ccipOffRampSeqInput.DeployCCIPOffRampInput.CCIPPackageId = ccipReport.Output.CCIPPackageId
	ccipOffRampSeqInput.DeployCCIPOffRampInput.MCMSPackageId = s.mcmsPackageID

	ccipOffRampSeqInput.InitializeOffRampInput.DestTransferCapId = ccipReport.Output.Objects.DestTransferCapObjectId
	ccipOffRampSeqInput.InitializeOffRampInput.FeeQuoterCapId = ccipReport.Output.Objects.FeeQuoterCapObjectId
	ccipOffRampSeqInput.InitializeOffRampInput.ChainSelector = uint64(s.chainSelector)
	ccipOffRampSeqInput.InitializeOffRampInput.SourceChainSelectors = []uint64{
		cselectors.ETHEREUM_MAINNET.Selector,
	}
	ccipOffRampSeqInput.InitializeOffRampInput.SourceChainsOnRamp = onRampBytes

	ccipOffRampSeqReport, err := operations.ExecuteSequence(s.bundle, offrampops.DeployAndInitCCIPOffRampSequence, s.deps, ccipOffRampSeqInput)
	require.NoError(s.T(), err, "failed to execute CCIP OffRamp deploy sequence")

	s.ccipOfframpPackageId = ccipOffRampSeqReport.Output.CCIPOffRampPackageId
	s.latestCcipOfframpPackageId = ccipOffRampSeqReport.Output.CCIPOffRampPackageId
	s.ccipOfframpObjects = ccipOffRampSeqReport.Output.Objects
}

func (s *MCMSTestSuite) SignProposal(proposal *mcms.Proposal, roleConfig *RoleConfig) {
	inspector, err := suisdk.NewInspector(s.client, s.signer, s.mcmsPackageID, roleConfig.Role)
	s.Require().NoError(err, "creating inspector for op count query")

	inspectorsMap := map[types.ChainSelector]sdk.Inspector{
		s.chainSelector: inspector,
	}

	signable, err := mcms.NewSignable(proposal, inspectorsMap)
	s.Require().NoError(err, "creating signable proposal")

	for i := 0; i < len(roleConfig.Keys) && i < roleConfig.Count; i++ {
		_, err = signable.SignAndAppend(mcms.NewPrivateKeySigner(roleConfig.Keys[i]))
		s.Require().NoError(err, "signing proposal")
	}

	// Need to query inspector with MCMS state object ID
	quorumMet, err := signable.ValidateSignatures(s.T().Context())
	s.Require().NoError(err, "Error validating signatures")
	s.Require().True(quorumMet, "Quorum not met")
}

func (s *MCMSTestSuite) ConvertProposal(timelockProposal *mcms.TimelockProposal) *mcms.Proposal {
	// Convert the Timelock Proposal into a MCMS Proposal
	timelockConverter, err := suisdk.NewTimelockConverter()
	s.Require().NoError(err)

	convertersMap := map[types.ChainSelector]sdk.TimelockConverter{
		s.chainSelector: timelockConverter,
	}
	proposal, _, err := timelockProposal.Convert(s.T().Context(), convertersMap)
	s.Require().NoError(err)
	return &proposal
}

func (s *MCMSTestSuite) SetRoot(proposal *mcms.Proposal, roleConfig *RoleConfig) {
	encoders, err := proposal.GetEncoders()
	s.Require().NoError(err)
	suiEncoder := encoders[s.chainSelector].(*suisdk.Encoder)
	executor, err := suisdk.NewExecutor(s.client, s.signer, suiEncoder, mcmsencoder.NewCCIPEntrypointArgEncoder(s.registryObj, s.deployerStateObj), s.mcmsPackageID, roleConfig.Role, s.mcmsObj, s.accountObj, s.registryObj, s.timelockObj)
	s.Require().NoError(err, "creating executor for Sui mcms contract")

	executors := map[types.ChainSelector]sdk.Executor{
		s.chainSelector: executor,
	}
	executable, err := mcms.NewExecutable(proposal, executors)
	s.Require().NoError(err, "Error creating executable")

	_, err = executable.SetRoot(s.T().Context(), s.chainSelector)
	s.Require().NoError(err)

}

func (s *MCMSTestSuite) Execute(timelockProposal *mcms.TimelockProposal, proposal *mcms.Proposal, proposalDelay time.Duration, roleConfig *RoleConfig) []types.TransactionResult {
	encoders, err := proposal.GetEncoders()
	s.Require().NoError(err)
	suiEncoder := encoders[s.chainSelector].(*suisdk.Encoder)
	executor, err := suisdk.NewExecutor(s.client, s.signer, suiEncoder, mcmsencoder.NewCCIPEntrypointArgEncoder(s.registryObj, s.deployerStateObj), s.mcmsPackageID, roleConfig.Role, s.mcmsObj, s.accountObj, s.registryObj, s.timelockObj)
	s.Require().NoError(err, "creating executor for Sui mcms contract")

	executors := map[types.ChainSelector]sdk.Executor{
		s.chainSelector: executor,
	}
	executable, err := mcms.NewExecutable(proposal, executors)
	s.Require().NoError(err, "Error creating executable")

	var responses []types.TransactionResult
	for i := range proposal.Operations {
		res, execErr := executable.Execute(s.T().Context(), i)
		s.Require().NoError(execErr)
		responses = append(responses, res)
	}
	if roleConfig.Role == suisdk.TimelockRoleProposer {
		// If proposer, some time needs to pass before the proposal can be executed sleep for delay5s
		time.Sleep(proposalDelay)

		timelockExecutor, tErr := suisdk.NewTimelockExecutor(
			s.client,
			s.signer,
			mcmsencoder.NewCCIPEntrypointArgEncoder(s.registryObj, s.deployerStateObj),
			s.mcmsPackageID,
			s.registryObj,
			s.accountObj,
		)

		s.Require().NoError(tErr, "creating timelock executor for Sui mcms contract")
		timelockExecutors := map[types.ChainSelector]sdk.TimelockExecutor{
			s.chainSelector: timelockExecutor,
		}
		timelockExecutable, execErr := mcms.NewTimelockExecutable(s.T().Context(), timelockProposal, timelockExecutors)
		s.Require().NoError(execErr)

		res, terr := timelockExecutable.Execute(s.T().Context(), 0, mcms.WithCallProxy(s.timelockObj))
		s.Require().NoError(terr)
		responses = append(responses, res)
	}

	return responses
}

func (s *MCMSTestSuite) ExecuteProposalE2e(timelockProposal *mcms.TimelockProposal, roleConfig *RoleConfig, proposalDelay time.Duration) []types.TransactionResult {
	proposal := s.ConvertProposal(timelockProposal)
	s.SignProposal(proposal, roleConfig)
	s.SetRoot(proposal, roleConfig)
	responses := s.Execute(timelockProposal, proposal, proposalDelay, roleConfig)
	s.T().Logf("✅ Executed MCMS proposal: %s", timelockProposal.Description)
	return responses
}

// Reused in other tests
func (s *MCMSTestSuite) RunOwnershipCCIPTransfer() {
	// 1a. Transfer OwnerCap of CCIP to MCMS (this should be done in the initial deployment sequence)
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
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	s.T().Logf("✅ Transferred ownership of CCIP to MCMS in tx: %s", tx.Digest)

	// 1b. Transfer ownership of the CCIP Router to MCMS
	ccipRouterContract, err := module_router.NewRouter(s.ccipRouterPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip router contract")

	tx, err = ccipRouterContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipRouterObjects.RouterStateObjectId},
		bind.Object{Id: s.ccipRouterObjects.OwnerCapObjectId},
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP Router to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	s.T().Logf("✅ Transferred ownership of CCIP Router to MCMS in tx: %s", tx.Digest)

	// 1c. Transfer ownership of the CCIP OnRamp to MCMS
	ccipOnRampContract, err := module_onramp.NewOnramp(s.ccipOnrampPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip onramp contract")

	tx, err = ccipOnRampContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId},
		bind.Object{Id: s.ccipOnrampObjects.StateObjectId},
		bind.Object{Id: s.ccipOnrampObjects.OwnerCapObjectId},
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP OnRamp to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	s.T().Logf("✅ Transferred ownership of CCIP OnRamp to MCMS in tx: %s", tx.Digest)

	// 1d. Transfer ownership of the CCIP OffRamp to MCMS
	ccipOffRampContract, err := module_offramp.NewOfframp(s.ccipOfframpPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip offramp contract")

	tx, err = ccipOffRampContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId},
		bind.Object{Id: s.ccipOfframpObjects.StateObjectId},
		bind.Object{Id: s.ccipOfframpObjects.OwnerCapId},
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP OffRamp to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	// 2. Proposal execution with acceptance from MCMS (through bypasser)
	input := ownershipops.AcceptCCIPOwnershipInput{
		// MCMS related
		MCMSPackageId:          s.mcmsPackageID,
		MCMSStateObjId:         s.mcmsObj,
		MCMSTimelockObjId:      s.timelockObj,
		MCMSAccountObjId:       s.accountObj,
		MCMSRegistryObjId:      s.registryObj,
		MCMSDeployerStateObjId: s.deployerStateObj,

		CCIPPackageId: s.ccipPackageId,
		CCIPObjectRef: s.ccipObjects.CCIPObjectRefObjectId,

		RouterPackageId:     s.ccipRouterPackageId,
		RouterStateObjectId: s.ccipRouterObjects.RouterStateObjectId,

		OnRampPackageId:     s.ccipOnrampPackageId,
		OnRampStateObjectId: s.ccipOnrampObjects.StateObjectId,

		OffRampPackageId:     s.ccipOfframpPackageId,
		OffRampStateObjectId: s.ccipOfframpObjects.StateObjectId,

		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},
		ChainSelector: uint64(s.chainSelector),
	}
	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(s.bundle, ownershipops.AcceptCCIPOwnershipSeq, s.deps, input)
	s.Require().NoError(err, "executing ownership acceptance proposal sequence")

	timelockProposal := acceptOwnershipProposalReport.Output

	// 3. Execute transfer ownership from original owner
	// 3.1. Execute the proposal
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// 3.2. Finish the ownership transfer with the original owner signer
	executeTransferInput := ownershipops.ExecuteOwnershipTransferToMcmsSeqInput{
		StateObject: &ccipops.ExecuteOwnershipTransferToMcmsStateObjectInput{
			CCIPPackageId:         s.ccipPackageId,
			CCIPObjectRefObjectId: s.ccipObjects.CCIPObjectRefObjectId,
			OwnerCapObjectId:      s.ccipObjects.OwnerCapObjectId,
			RegistryObjectId:      s.registryObj,
			To:                    s.mcmsPackageID,
		},
		OnRamp: &onrampops.ExecuteOwnershipTransferToMcmsOnRampInput{
			OnRampPackageId:     s.ccipOnrampPackageId,
			OnRampRefObjectId:   s.ccipObjects.CCIPObjectRefObjectId,
			OwnerCapObjectId:    s.ccipOnrampObjects.OwnerCapObjectId,
			OnRampStateObjectId: s.ccipOnrampObjects.StateObjectId,
			RegistryObjectId:    s.registryObj,
			To:                  s.mcmsPackageID,
		},
		OffRamp: &offrampops.ExecuteOwnershipTransferToMcmsOffRampInput{
			OffRampPackageId:     s.ccipOfframpPackageId,
			OffRampRefObjectId:   s.ccipObjects.CCIPObjectRefObjectId,
			OwnerCapObjectId:     s.ccipOfframpObjects.OwnerCapId,
			OffRampStateObjectId: s.ccipOfframpObjects.StateObjectId,
			RegistryObjectId:     s.registryObj,
			To:                   s.mcmsPackageID,
		},
		Router: &routerops.ExecuteOwnershipTransferToMcmsRouterInput{
			RouterPackageId:     s.ccipRouterPackageId,
			OwnerCapObjectId:    s.ccipRouterObjects.OwnerCapObjectId,
			RouterStateObjectId: s.ccipRouterObjects.RouterStateObjectId,
			RegistryObjectId:    s.registryObj,
			To:                  s.mcmsPackageID,
		},
	}
	_, err = cld_ops.ExecuteSequence(s.bundle, ownershipops.ExecuteOwnershipTransferToMcmsSequence, s.deps, executeTransferInput)
	s.Require().NoError(err, "executing final ownership transfer to MCMS")

	// 4. Verify the new owner is MCMS
	newOwner, err := ccipContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId})
	s.Require().NoError(err, "getting new owner of CCIP state object")
	s.Require().Equal(s.mcmsPackageID, newOwner, "new owner of CCIP should be MCMS")

	newOwnerOnRamp, err := ccipOnRampContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipOnrampObjects.StateObjectId})
	s.Require().NoError(err, "getting new owner of OnRamp state object")
	s.Require().Equal(s.mcmsPackageID, newOwnerOnRamp, "new owner of OnRamp should be MCMS")

	newOwnerOffRamp, err := ccipOffRampContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipOfframpObjects.StateObjectId})
	s.Require().NoError(err, "getting new owner of OffRamp state object")
	s.Require().Equal(s.mcmsPackageID, newOwnerOffRamp, "new owner of OffRamp should be MCMS")

	newOwnerRouter, err := ccipRouterContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipRouterObjects.RouterStateObjectId})
	s.Require().NoError(err, "getting new owner of Router state object")
	s.Require().Equal(s.mcmsPackageID, newOwnerRouter, "new owner of Router should be MCMS")
}

func (s *MCMSTestSuite) VerifyVersion(packageId string, expectedVersion string) {
	// we can use mcmsuser contract always since the interface is the same
	userContract, err := module_user.NewMcmsUser(packageId, s.client)
	s.Require().NoError(err, "creating upgraded contract")
	// Call type_and_version function
	version, err := userContract.DevInspect().TypeAndVersion(s.T().Context(), s.deps.GetCallOpts())
	require.NoError(s.T(), err, "getting type and version")

	s.T().Logf("✅ New package version: %s", version)
	require.Equal(s.T(), expectedVersion, version, "version should match expected")
}

func (s *MCMSTestSuite) GetUpgradedAddress(result *models.SuiTransactionBlockResponse, mcmsPackageID string) (string, error) {
	s.T().Helper()

	if result == nil || result.Events == nil {
		return "", errors.New("result is nil or events are nil")
	}

	for _, event := range result.Events {
		if isUpgradeEvent(event, mcmsPackageID) {
			return processUpgradeEvent(s.T(), event)
		}
	}

	return "", errors.New("upgrade receipt committed event not found")
}

// isUpgradeEvent checks if the event is an upgrade receipt committed event
func isUpgradeEvent(event models.SuiEventResponse, mcmsPackageID string) bool {
	return event.PackageId == mcmsPackageID &&
		event.TransactionModule == "mcms_deployer" &&
		strings.Contains(event.Type, "UpgradeReceiptCommitted")
}

// processUpgradeEvent processes an upgrade event and returns the new package address
func processUpgradeEvent(t *testing.T, event models.SuiEventResponse) (string, error) {
	t.Helper()

	if event.ParsedJson == nil {
		return "", errors.New("parsed json is nil")
	}

	oldAddr := event.ParsedJson["old_package_address"]
	newAddr := event.ParsedJson["new_package_address"]
	oldVer := event.ParsedJson["old_version"]
	newVer := event.ParsedJson["new_version"]

	newAddress, err := validateAddressChange(t, oldAddr, newAddr)
	if err != nil {
		return "", err
	}

	err = validateVersionIncrement(t, oldVer, newVer)
	if err != nil {
		return "", err
	}

	return newAddress, nil
}

// validateAddressChange validates that the package address changed correctly
func validateAddressChange(t *testing.T, oldAddr, newAddr any) (string, error) {
	t.Helper()

	if oldAddr == nil || newAddr == nil {
		return "", errors.New("package addresses are nil")
	}

	oldAddrStr := fmt.Sprintf("%v", oldAddr)
	newAddrStr := fmt.Sprintf("%v", newAddr)

	if oldAddrStr == newAddrStr {
		t.Errorf("ERROR: Package address did not change! Old: %v, New: %v", oldAddr, newAddr)
		return "", errors.New("package address did not change")
	}

	return newAddrStr, nil
}

// validateVersionIncrement validates that the version incremented correctly
func validateVersionIncrement(t *testing.T, oldVer, newVer any) error {
	t.Helper()

	if oldVer == nil || newVer == nil {
		return nil // Version validation is optional
	}

	oldVersion, oldParseOk := parseVersion(t, oldVer, "old")
	newVersion, newParseOk := parseVersion(t, newVer, "new")

	if !oldParseOk || !newParseOk {
		return nil // Skip validation if parsing failed
	}

	expectedVersion := oldVersion + 1
	if newVersion != expectedVersion {
		t.Errorf("ERROR: Version did not increment correctly! Old: %.0f, New: %.0f (expected %.0f)",
			oldVersion, newVersion, expectedVersion)

		return errors.New("version did not increment correctly")
	}

	return nil
}

// parseVersion parses a version value from interface{} to float64
func parseVersion(t *testing.T, version any, versionType string) (float64, bool) {
	t.Helper()

	switch v := version.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			return parsed, true
		}
		t.Logf("Warning: Could not parse %s version string '%s' as number", versionType, v)

		return 0, false
	default:
		t.Logf("Warning: Unsupported %s version type: %T", versionType, v)
		return 0, false
	}
}
