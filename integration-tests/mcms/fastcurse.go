//go:build integration

package mcms

import (
	"encoding/json"
	"testing"
	"time"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/adapters"
	"github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	rmn_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type CCIPCurseMCMSTestSuite struct {
	MCMSTestSuite

	fastBypasserConfig *RoleConfig
	curserCapObjectID  string
}

func (s *CCIPCurseMCMSTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()
	s.deployFastMCMS()
}

func (s *CCIPCurseMCMSTestSuite) TestCurseMCMSTest() {
	s.T().Run("Direct curse/uncurse via sequences", func(_ *testing.T) {
		s.testDirectCurseUncurse()
	})
	s.T().Run("MCMS proposal curse/uncurse (slow OwnerCap)", func(_ *testing.T) {
		s.testMCMSCurseProposal()
	})
	s.T().Run("Fast MCMS bootstrap and curse via CurserCap", func(_ *testing.T) {
		s.testFastMCMSCurseViaCurserCap()
	})
	s.T().Run("Slow MCMS uncurse after fast CurserCap curse", func(_ *testing.T) {
		s.testSlowMCMSUncurseAfterFastCurse()
	})
}

func (s *CCIPCurseMCMSTestSuite) deployFastMCMS() {
	// Fast MCMS package + objects were published in MCMSTestSuite.SetupSuite (before CCIP).
	// Configure that same instance so CurserCap bootstrap uses the fast_mcms types CCIP was compiled against.
	bundle := s.NewOpBundleWithRegistry()

	_, err := cld_ops.ExecuteSequence(bundle, mcmsops.ConfigureMCMSSequence, s.deps, mcmsops.ConfigureMCMSSeqInput{
		ChainSelector:               uint64(s.chainSelector),
		PackageId:                   s.fastMcmsPackageID,
		McmsAccountOwnerCapObjectId: s.fastOwnerCapObj,
		McmsAccountStateObjectId:    s.fastAccountObj,
		McmsMultisigStateObjectId:   s.fastMcmsObj,
		Bypasser:                    s.bypasserConfig.Config,
		Proposer:                    s.proposerConfig.Config,
	})
	s.Require().NoError(err, "configuring fast MCMS contract")

	s.fastBypasserConfig = s.bypasserConfig
	s.Require().NotEqual(s.mcmsObj, s.fastMcmsObj, "fast and slow MCMS must be distinct instances")

	_, err = cld_ops.ExecuteOperation(bundle, mcmsops.MCMSTransferOwnershipOp, s.deps, mcmsops.MCMSTransferOwnershipInput{
		McmsPackageID:   s.fastMcmsPackageID,
		OwnerCap:        s.fastOwnerCapObj,
		AccountObjectID: s.fastAccountObj,
	})
	s.Require().NoError(err, "initiating fast MCMS ownership transfer to self")

	acceptReport, err := cld_ops.ExecuteSequence(bundle, mcmsops.AcceptMCMSOwnershipSequence, s.deps, mcmsops.AcceptMCMSOwnershipSeqInput{
		ChainSelector:             uint64(s.chainSelector),
		PackageId:                 s.fastMcmsPackageID,
		McmsMultisigStateObjectId: s.fastMcmsObj,
		TimelockObjectId:          s.fastTimelockObj,
		McmsAccountStateObjectId:  s.fastAccountObj,
		McmsRegistryObjectId:      s.fastRegistryObj,
		McmsDeployerStateObjectId: s.fastDeployerStateObj,
		TimelockConfig: utils.TimelockConfig{
			MCMSAction: types.TimelockActionSchedule,
		},
	})
	s.Require().NoError(err, "generating fast MCMS accept-ownership proposal")

	s.executeFastProposalE2e(&acceptReport.Output, s.proposerConfig, 0)

	_, err = cld_ops.ExecuteOperation(bundle, mcmsops.MCMSExecuteTransferOwnershipOp, s.deps, mcmsops.MCMSExecuteTransferOwnershipInput{
		McmsPackageID:         s.fastMcmsPackageID,
		OwnerCap:              s.fastOwnerCapObj,
		AccountObjectID:       s.fastAccountObj,
		RegistryObjectID:      s.fastRegistryObj,
		DeployerStateObjectID: s.fastDeployerStateObj,
	})
	s.Require().NoError(err, "fast MCMS ownership transfer to self")
	s.T().Logf("✅ Configured fast MCMS registry %s", s.fastRegistryObj)
}

func (s *CCIPCurseMCMSTestSuite) newAdapter() *adapters.CurseAdapter {
	return &adapters.CurseAdapter{
		CCIPAddress:          s.ccipPackageId,
		CCIPObjectRef:        s.ccipObjects.CCIPObjectRefObjectId,
		CCIPOwnerCapObjectID: s.ccipObjects.OwnerCapObjectId,
		CurserCapObjectID:    s.curserCapObjectID,
		RouterAddress:        s.ccipRouterPackageId,
		RouterStateObjectID:  s.ccipRouterObjects.RouterStateObjectId,
	}
}

func (s *CCIPCurseMCMSTestSuite) buildSuiChains(includeSigner bool) cldf_chain.BlockChains {
	c := cldfsui.Chain{
		ChainMetadata: cldfsui.ChainMetadata{Selector: uint64(s.chainSelector)},
		Client:        s.client,
	}
	if includeSigner {
		c.Signer = s.signer
	}
	return cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		uint64(s.chainSelector): c,
	})
}

func (s *CCIPCurseMCMSTestSuite) buildEnv() cldf.Environment {
	return cldf.Environment{
		BlockChains: s.buildSuiChains(true),
	}
}

func (s *CCIPCurseMCMSTestSuite) assertIsCursed(a *adapters.CurseAdapter, env cldf.Environment, subject fastcurse.Subject, expected bool) {
	s.T().Helper()
	cursed, err := a.IsSubjectCursedOnChain(env, uint64(s.chainSelector), subject)
	s.Require().NoError(err, "IsSubjectCursedOnChain")
	s.Require().Equal(expected, cursed, "unexpected curse state for subject %x", subject)
}

func (s *CCIPCurseMCMSTestSuite) testDirectCurseUncurse() {
	a := s.newAdapter()
	env := s.buildEnv()
	chains := s.buildSuiChains(true)
	bundle := s.NewOpBundle()

	subject := a.SelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector)
	curseInput := fastcurse.CurseInput{
		ChainSelector: uint64(s.chainSelector),
		Subjects:      []fastcurse.Subject{subject},
	}

	s.assertIsCursed(a, env, subject, false)

	_, err := cld_ops.ExecuteSequence(bundle, a.Curse(), chains, curseInput)
	s.Require().NoError(err, "executing Curse() sequence directly")
	s.assertIsCursed(a, env, subject, true)

	_, err = cld_ops.ExecuteSequence(bundle, a.Uncurse(), chains, curseInput)
	s.Require().NoError(err, "executing Uncurse() sequence directly")
	s.assertIsCursed(a, env, subject, false)
}

func (s *CCIPCurseMCMSTestSuite) testMCMSCurseProposal() {
	s.RunOwnershipCCIPTransfer()

	a := s.newAdapter()
	env := s.buildEnv()
	chains := s.buildSuiChains(false)
	bundle := s.NewOpBundle()

	subject := fastcurse.GlobalCurseSubject()
	curseInput := fastcurse.CurseInput{
		ChainSelector: uint64(s.chainSelector),
		Subjects:      []fastcurse.Subject{subject},
	}

	curseReport, err := cld_ops.ExecuteSequence(bundle, a.Curse(), chains, curseInput)
	s.Require().NoError(err, "building curse MCMS batch operations")

	curseProposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.mcmsPackageID,
		MCMSStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		Description:        "Integration test: curse global subject via slow MCMS bypasser",
		BatchOp:            curseReport.Output.BatchOps[0],
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	s.Require().NoError(err)

	s.ExecuteProposalE2e(curseProposal, s.bypasserConfig, 0)
	s.assertIsCursed(a, env, subject, true)

	uncurseReport, err := cld_ops.ExecuteSequence(bundle, a.Uncurse(), chains, curseInput)
	s.Require().NoError(err)

	uncurseProposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.mcmsPackageID,
		MCMSStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		Description:        "Integration test: uncurse global subject via slow MCMS bypasser",
		BatchOp:            uncurseReport.Output.BatchOps[0],
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	s.Require().NoError(err)

	s.ExecuteProposalE2e(uncurseProposal, s.bypasserConfig, 0)
	s.assertIsCursed(a, env, subject, false)
}

func (s *CCIPCurseMCMSTestSuite) ensureCCIPOwnedBySlowMCMS() {
	ccipContract, err := module_state_object.NewStateObject(s.ccipPackageId, s.client)
	s.Require().NoError(err)
	owner, err := ccipContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId})
	s.Require().NoError(err)
	if owner != s.mcmsPackageID {
		s.RunOwnershipCCIPTransfer()
	}
}

func (s *CCIPCurseMCMSTestSuite) bootstrapCurserCap() {
	s.ensureCCIPOwnedBySlowMCMS()

	bundle := s.NewOpBundleWithRegistry()
	deps := s.deps
	deps.Signer = nil

	report, err := cld_ops.ExecuteOperation(bundle, rmn_ops.McmsMintAndRegisterCurserCapOp, deps, rmn_ops.McmsMintAndRegisterCurserCapInput{
		CCIPPackageId:        s.ccipPackageId,
		StateObjectId:        s.ccipObjects.CCIPObjectRefObjectId,
		SlowOwnerCapObjectId: s.ccipObjects.OwnerCapObjectId,
		FastRegistryObjectId: s.fastRegistryObj,
	})
	s.Require().NoError(err, "encoding mint_and_register_curser_cap leaf")

	genericReport := report.ToGenericReport()
	mcmsConfig := mcmsops.ProposalGenerateInput{
		ChainSelector:      uint64(s.chainSelector),
		Defs:               []cld_ops.Definition{genericReport.Def},
		Inputs:             []any{genericReport.Input},
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	}
	seqResult, err := cld_ops.ExecuteSequence(bundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
	s.Require().NoError(err, "generating bootstrap MCMS proposal")

	results := s.ExecuteProposalE2e(&seqResult.Output, s.bypasserConfig, 0)
	s.Require().NotEmpty(results)
	execTxDigest, err := changesets.LastSuccessfulReportTxDigest(results)
	s.Require().NoError(err)

	curserCapID, err := utils.ResolveCurserCapObjectID(
		s.T().Context(),
		s.client,
		execTxDigest,
		s.fastRegistryObj,
		s.ccipPackageId,
	)
	s.Require().NoError(err, "finding CurserCap from executed bootstrap proposal")
	s.curserCapObjectID = curserCapID
	s.Require().NotEmpty(s.curserCapObjectID)
	s.T().Logf("Registered CurserCap %s in fast MCMS Registry", s.curserCapObjectID)

	ab := cldf.NewMemoryAddressBook()
	tv := cldf.NewTypeAndVersion(deployment.SuiCurserCapObjectIDType, deployment.Version1_0_0)
	s.Require().NoError(ab.Save(uint64(s.chainSelector), s.curserCapObjectID, tv))

	// Concern 2: off-chain validation accepts the issued cap id from state and rejects mismatches.
	envWithCap := cldf.Environment{
		ExistingAddresses: ab,
		BlockChains:       s.buildSuiChains(false),
	}
	cs := changesets.CurseUncurseChains{}
	s.Require().NoError(cs.VerifyPreconditions(envWithCap, changesets.CurseUncurseChainsConfig{
		SuiChainSelector: uint64(s.chainSelector),
		OperationType:    string(changesets.CurseOperationType),
		IsGlobalCurse:    true,
		IsFastCurse:      true,
	}))
	err = cs.VerifyPreconditions(envWithCap, changesets.CurseUncurseChainsConfig{
		SuiChainSelector:  uint64(s.chainSelector),
		OperationType:     string(changesets.CurseOperationType),
		IsGlobalCurse:     true,
		IsFastCurse:       true,
		CurserCapObjectId: "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "does not match registered CurserCap")
}

func (s *CCIPCurseMCMSTestSuite) testFastMCMSCurseViaCurserCap() {
	s.bootstrapCurserCap()

	a := s.newAdapter()
	s.Require().Equal(s.curserCapObjectID, a.CurserCapObjectID)
	env := s.buildEnv()
	chains := s.buildSuiChains(false)
	bundle := s.NewOpBundle()

	subject := a.SelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector)
	curseInput := fastcurse.CurseInput{
		ChainSelector: uint64(s.chainSelector),
		Subjects:      []fastcurse.Subject{subject},
	}

	s.assertIsCursed(a, env, subject, false)

	curseReport, err := cld_ops.ExecuteSequence(bundle, a.Curse(), chains, curseInput)
	s.Require().NoError(err, "building fast curse MCMS batch operations")

	tx := curseReport.Output.BatchOps[0].Transactions[0]
	s.Require().Equal(s.ccipPackageId, tx.To)
	var txFields suisdk.AdditionalFields
	s.Require().NoError(json.Unmarshal(tx.AdditionalFields, &txFields))
	s.Require().Equal("rmn_remote", txFields.ModuleName)
	s.Require().Equal("curse_multiple_with_curser_cap", txFields.Function)

	proposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.fastMcmsPackageID,
		MCMSStateObjID:     s.fastMcmsObj,
		TimelockObjID:      s.fastTimelockObj,
		AccountObjID:       s.fastAccountObj,
		RegistryObjID:      s.fastRegistryObj,
		DeployerStateObjID: s.fastDeployerStateObj,
		Description:        "Integration test: fast MCMS curse via CurserCap",
		BatchOp:            curseReport.Output.BatchOps[0],
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	s.Require().NoError(err)

	s.executeFastProposalE2e(proposal, s.fastBypasserConfig, 0)
	s.assertIsCursed(a, env, subject, true)
}

func (s *CCIPCurseMCMSTestSuite) testSlowMCMSUncurseAfterFastCurse() {
	s.bootstrapCurserCap()

	a := s.newAdapter()
	env := s.buildEnv()
	chains := s.buildSuiChains(false)
	bundle := s.NewOpBundle()

	subject := a.SelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector)
	curseInput := fastcurse.CurseInput{
		ChainSelector: uint64(s.chainSelector),
		Subjects:      []fastcurse.Subject{subject},
	}

	curseReport, err := cld_ops.ExecuteSequence(bundle, a.Curse(), chains, curseInput)
	s.Require().NoError(err)
	curseProposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.fastMcmsPackageID,
		MCMSStateObjID:     s.fastMcmsObj,
		TimelockObjID:      s.fastTimelockObj,
		AccountObjID:       s.fastAccountObj,
		RegistryObjID:      s.fastRegistryObj,
		DeployerStateObjID: s.fastDeployerStateObj,
		Description:        "Integration test: fast MCMS curse before slow uncurse",
		BatchOp:            curseReport.Output.BatchOps[0],
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	s.Require().NoError(err)
	s.executeFastProposalE2e(curseProposal, s.fastBypasserConfig, 0)
	s.assertIsCursed(a, env, subject, true)

	slowAdapter := &adapters.CurseAdapter{
		CCIPAddress:          s.ccipPackageId,
		CCIPObjectRef:        s.ccipObjects.CCIPObjectRefObjectId,
		CCIPOwnerCapObjectID: s.ccipObjects.OwnerCapObjectId,
		RouterAddress:        s.ccipRouterPackageId,
		RouterStateObjectID:  s.ccipRouterObjects.RouterStateObjectId,
	}

	uncurseReport, err := cld_ops.ExecuteSequence(bundle, slowAdapter.Uncurse(), chains, curseInput)
	s.Require().NoError(err)

	uncurseProposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.mcmsPackageID,
		MCMSStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		Description:        "Integration test: slow MCMS uncurse after fast CurserCap curse",
		BatchOp:            uncurseReport.Output.BatchOps[0],
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	s.Require().NoError(err)

	s.ExecuteProposalE2e(uncurseProposal, s.bypasserConfig, 0)
	s.assertIsCursed(slowAdapter, env, subject, false)
}

type mcmsEndpointSnapshot struct {
	packageID        string
	mcmsObj          string
	timelockObj      string
	registryObj      string
	deployerStateObj string
	accountObj       string
}

func (s *CCIPCurseMCMSTestSuite) snapshotSlowMCMSEndpoints() mcmsEndpointSnapshot {
	return mcmsEndpointSnapshot{
		packageID:        s.mcmsPackageID,
		mcmsObj:          s.mcmsObj,
		timelockObj:      s.timelockObj,
		registryObj:      s.registryObj,
		deployerStateObj: s.deployerStateObj,
		accountObj:       s.accountObj,
	}
}

func (s *CCIPCurseMCMSTestSuite) restoreSlowMCMSEndpoints(snap mcmsEndpointSnapshot) {
	s.mcmsPackageID = snap.packageID
	s.mcmsObj = snap.mcmsObj
	s.timelockObj = snap.timelockObj
	s.registryObj = snap.registryObj
	s.deployerStateObj = snap.deployerStateObj
	s.accountObj = snap.accountObj
}

func (s *CCIPCurseMCMSTestSuite) executeFastProposalE2e(timelockProposal *mcms.TimelockProposal, roleConfig *RoleConfig, proposalDelay time.Duration) {
	saved := s.snapshotSlowMCMSEndpoints()
	s.mcmsPackageID = s.fastMcmsPackageID
	s.mcmsObj = s.fastMcmsObj
	s.timelockObj = s.fastTimelockObj
	s.registryObj = s.fastRegistryObj
	s.deployerStateObj = s.fastDeployerStateObj
	s.accountObj = s.fastAccountObj
	defer s.restoreSlowMCMSEndpoints(saved)

	s.ExecuteProposalE2e(timelockProposal, roleConfig, proposalDelay)
}
