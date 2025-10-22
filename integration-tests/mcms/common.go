//go:build integration

package mcms

import (
	"crypto/ecdsa"
	"slices"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	mcmsencoder "github.com/smartcontractkit/chainlink-sui/bindings"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	cselectors "github.com/smartcontractkit/chain-selectors"
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
	accountObj       string
	ownerCapObj      string

	bypasserConfig *RoleConfig
	proposerConfig *RoleConfig

	// Ops
	deps   sui_ops.OpTxDeps
	bundle cld_ops.Bundle
}

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
	ops := make([]*cld_ops.Operation[any, any, any], len(opregistry.AllOperations))
	for i := range opregistry.AllOperations {
		ops[i] = &opregistry.AllOperations[i]
	}
	registry := cld_ops.NewOperationRegistry(
		ops...,
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
		McmsPackageID:    s.mcmsPackageID,
		OwnerCap:         s.ownerCapObj,
		AccountObjectID:  s.accountObj,
		RegistryObjectID: s.registryObj,
	})
	s.Require().NoError(err, "executing ownership transfer to self")
	s.T().Logf("✅ Transferred ownership of MCMS to itself in tx: %s", rep.Output.Digest)
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
	executor, err := suisdk.NewExecutor(s.client, s.signer, suiEncoder, mcmsencoder.NewCCIPEntrypointArgEncoder(s.registryObj), s.mcmsPackageID, roleConfig.Role, s.mcmsObj, s.accountObj, s.registryObj, s.timelockObj)
	s.Require().NoError(err, "creating executor for Sui mcms contract")

	executors := map[types.ChainSelector]sdk.Executor{
		s.chainSelector: executor,
	}
	executable, err := mcms.NewExecutable(proposal, executors)
	s.Require().NoError(err, "Error creating executable")

	_, err = executable.SetRoot(s.T().Context(), s.chainSelector)
	s.Require().NoError(err)

}

func (s *MCMSTestSuite) Execute(timelockProposal *mcms.TimelockProposal, proposal *mcms.Proposal, proposalDelay time.Duration, roleConfig *RoleConfig) {
	encoders, err := proposal.GetEncoders()
	s.Require().NoError(err)
	suiEncoder := encoders[s.chainSelector].(*suisdk.Encoder)
	executor, err := suisdk.NewExecutor(s.client, s.signer, suiEncoder, mcmsencoder.NewCCIPEntrypointArgEncoder(s.registryObj), s.mcmsPackageID, roleConfig.Role, s.mcmsObj, s.accountObj, s.registryObj, s.timelockObj)
	s.Require().NoError(err, "creating executor for Sui mcms contract")

	executors := map[types.ChainSelector]sdk.Executor{
		s.chainSelector: executor,
	}
	executable, err := mcms.NewExecutable(proposal, executors)
	s.Require().NoError(err, "Error creating executable")

	for i := range proposal.Operations {
		_, execErr := executable.Execute(s.T().Context(), i)
		s.Require().NoError(execErr)
	}
	if roleConfig.Role == suisdk.TimelockRoleProposer {
		// If proposer, some time needs to pass before the proposal can be executed sleep for delay5s
		time.Sleep(proposalDelay)

		timelockExecutor, tErr := suisdk.NewTimelockExecutor(
			s.client,
			s.signer,
			mcmsencoder.NewCCIPEntrypointArgEncoder(s.registryObj),
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

		_, terr := timelockExecutable.Execute(s.T().Context(), 0, mcms.WithCallProxy(s.timelockObj))
		s.Require().NoError(terr)
	}
}

func (s *MCMSTestSuite) ExecuteProposalE2e(timelockProposal *mcms.TimelockProposal, roleConfig *RoleConfig, proposalDelay time.Duration) {
	proposal := s.ConvertProposal(timelockProposal)
	s.SignProposal(proposal, roleConfig)
	s.SetRoot(proposal, roleConfig)
	s.Execute(timelockProposal, proposal, proposalDelay, roleConfig)
}
