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
	"github.com/smartcontractkit/mcms/sdk"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	mcmsencoder "github.com/smartcontractkit/chainlink-sui/bindings"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
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
	s.T().Run("Explicit initialize allowlist via slow MCMS", func(_ *testing.T) {
		s.testExplicitInitializeAllowlistViaSlowMCMS()
	})
	s.T().Run("Fast MCMS bootstrap and curse via CurserCap", func(_ *testing.T) {
		s.testFastMCMSCurseViaCurserCap()
	})
	s.T().Run("CurserCap allowlist state after bootstrap", func(_ *testing.T) {
		s.testCurserCapAllowlistAfterBootstrap()
	})
	s.T().Run("Slow MCMS uncurse after fast CurserCap curse", func(_ *testing.T) {
		s.testSlowMCMSUncurseAfterFastCurse()
	})
	s.T().Run("Deregister CurserCap blocks fast curse", func(_ *testing.T) {
		s.testDeregisterCurserCapBlocksFastCurse()
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

func (s *CCIPCurseMCMSTestSuite) ccipObjectRef() bind.Object {
	return bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}
}

func (s *CCIPCurseMCMSTestSuite) rmnRemoteContract() module_rmn_remote.IRmnRemote {
	contract, err := module_rmn_remote.NewRmnRemote(s.ccipPackageId, s.client)
	s.Require().NoError(err)
	return contract
}

func (s *CCIPCurseMCMSTestSuite) assertAllowlistState(capID string, capAllowed bool) {
	s.T().Helper()

	contract := s.rmnRemoteContract()
	ctx := s.T().Context()
	opts := s.deps.GetCallOpts()
	ref := s.ccipObjectRef()

	ids, err := contract.DevInspect().GetAllowedCurserCapIds(ctx, opts, ref)
	s.Require().NoError(err)
	if capID == "" {
		s.Require().Empty(ids)
		return
	}

	if capAllowed {
		s.Require().Contains(ids, capID)
	} else {
		s.Require().NotContains(ids, capID)
	}

	allowed, err := contract.DevInspect().IsCurserCapAllowed(ctx, opts, ref, capID)
	s.Require().NoError(err)
	s.Require().Equal(capAllowed, allowed)
}

func (s *CCIPCurseMCMSTestSuite) executeSlowMCMSFromGenericReport(genericReport cld_ops.Report[any, any]) {
	bundle := s.NewOpBundleWithRegistry()
	deps := s.deps
	deps.Signer = nil

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
	s.Require().NoError(err, "generating slow MCMS proposal")
	s.ExecuteProposalE2e(&seqResult.Output, s.bypasserConfig, 0)
}

func (s *CCIPCurseMCMSTestSuite) slowCurseUncurseSeqInput(subjects []fastcurse.Subject) rmn_ops.CurseUncurseSeqInput {
	return rmn_ops.CurseUncurseSeqInput{
		CCIPAddress:          s.ccipPackageId,
		CCIPObjectRef:        s.ccipObjects.CCIPObjectRefObjectId,
		CCIPOwnerCapObjectID: s.ccipObjects.OwnerCapObjectId,
		ChainSelector:        uint64(s.chainSelector),
		Subjects:             subjects,
	}
}

func (s *CCIPCurseMCMSTestSuite) testDirectCurseUncurse() {
	a := s.newAdapter()
	env := s.buildEnv()
	chains := s.buildSuiChains(true)
	bundle := s.NewOpBundleWithRegistry()

	subject := a.SelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector)
	seqInput := s.slowCurseUncurseSeqInput([]fastcurse.Subject{subject})

	s.assertIsCursed(a, env, subject, false)

	_, err := cld_ops.ExecuteSequence(bundle, rmn_ops.CurseSequence, chains, seqInput)
	s.Require().NoError(err, "executing slow OwnerCap curse sequence directly")
	s.assertIsCursed(a, env, subject, true)

	_, err = cld_ops.ExecuteSequence(bundle, rmn_ops.UncurseSequence, chains, seqInput)
	s.Require().NoError(err, "executing slow OwnerCap uncurse sequence directly")
	s.assertIsCursed(a, env, subject, false)
}

func (s *CCIPCurseMCMSTestSuite) testMCMSCurseProposal() {
	s.RunOwnershipCCIPTransfer()

	a := s.newAdapter()
	env := s.buildEnv()
	chains := s.buildSuiChains(false)
	bundle := s.NewOpBundleWithRegistry()

	subject := fastcurse.GlobalCurseSubject()
	seqInput := s.slowCurseUncurseSeqInput([]fastcurse.Subject{subject})

	curseReport, err := cld_ops.ExecuteSequence(bundle, rmn_ops.CurseSequence, chains, seqInput)
	s.Require().NoError(err, "building slow OwnerCap curse MCMS batch operations")

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

	uncurseReport, err := cld_ops.ExecuteSequence(bundle, rmn_ops.UncurseSequence, chains, seqInput)
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

	if s.curserCapObjectID == "" {
		capID, err := utils.FindCurserCapInFastRegistry(
			s.T().Context(),
			s.client,
			s.fastRegistryObj,
			s.ccipPackageId,
		)
		if err == nil {
			s.curserCapObjectID = capID
		}
	}
	if s.curserCapObjectID != "" {
		s.T().Logf("Reusing CurserCap %s in fast MCMS Registry", s.curserCapObjectID)
		return
	}

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

	slowAdapter := &adapters.CurseAdapter{
		CCIPAddress:          s.ccipPackageId,
		CCIPObjectRef:        s.ccipObjects.CCIPObjectRefObjectId,
		CCIPOwnerCapObjectID: s.ccipObjects.OwnerCapObjectId,
		RouterAddress:        s.ccipRouterPackageId,
		RouterStateObjectID:  s.ccipRouterObjects.RouterStateObjectId,
	}

	uncurseReport, err := cld_ops.ExecuteSequence(bundle, slowAdapter.Uncurse(), chains, curseInput)
	s.Require().NoError(err, "building slow MCMS uncurse batch operations for cleanup")

	uncurseProposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.mcmsPackageID,
		MCMSStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		Description:        "Integration test: slow MCMS uncurse cleanup after fast CurserCap curse",
		BatchOp:            uncurseReport.Output.BatchOps[0],
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	s.Require().NoError(err)

	s.ExecuteProposalE2e(uncurseProposal, s.bypasserConfig, 0)
	s.assertIsCursed(slowAdapter, env, subject, false)
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

func (s *CCIPCurseMCMSTestSuite) testExplicitInitializeAllowlistViaSlowMCMS() {
	s.ensureCCIPOwnedBySlowMCMS()
	s.assertAllowlistState("", false)

	bundle := s.NewOpBundleWithRegistry()
	deps := s.deps
	deps.Signer = nil

	report, err := cld_ops.ExecuteOperation(bundle, rmn_ops.McmsInitializeAllowedCurserCapsOp, deps, rmn_ops.McmsInitializeAllowedCurserCapsInput{
		CCIPPackageId:        s.ccipPackageId,
		StateObjectId:        s.ccipObjects.CCIPObjectRefObjectId,
		SlowOwnerCapObjectId: s.ccipObjects.OwnerCapObjectId,
		InitialCurserCapIds:  nil,
	})
	s.Require().NoError(err, "encoding initialize_allowed_curser_caps leaf")
	s.executeSlowMCMSFromGenericReport(report.ToGenericReport())

	s.assertAllowlistState("", false)
}

func (s *CCIPCurseMCMSTestSuite) testCurserCapAllowlistAfterBootstrap() {
	s.bootstrapCurserCap()

	contract := s.rmnRemoteContract()
	ctx := s.T().Context()
	opts := s.deps.GetCallOpts()
	ref := s.ccipObjectRef()

	allowed, err := contract.DevInspect().IsCurserCapAllowed(ctx, opts, ref, s.curserCapObjectID)
	s.Require().NoError(err)
	s.Require().True(allowed, "bootstrapped cap must be permitted to curse")

	ids, err := contract.DevInspect().GetAllowedCurserCapIds(ctx, opts, ref)
	s.Require().NoError(err)
	s.Require().Contains(ids, s.curserCapObjectID)
}

func (s *CCIPCurseMCMSTestSuite) testDeregisterCurserCapBlocksFastCurse() {
	s.bootstrapCurserCap()
	s.assertAllowlistState(s.curserCapObjectID, true)

	bundle := s.NewOpBundleWithRegistry()
	deps := s.deps
	deps.Signer = nil

	deregisterReport, err := cld_ops.ExecuteOperation(bundle, rmn_ops.McmsDeregisterCurserCapIdsOp, deps, rmn_ops.McmsDeregisterCurserCapIdsInput{
		CCIPPackageId:        s.ccipPackageId,
		StateObjectId:        s.ccipObjects.CCIPObjectRefObjectId,
		SlowOwnerCapObjectId: s.ccipObjects.OwnerCapObjectId,
		CurserCapObjectIds:   []string{s.curserCapObjectID},
	})
	s.Require().NoError(err, "encoding deregister_curser_cap_ids leaf")
	s.executeSlowMCMSFromGenericReport(deregisterReport.ToGenericReport())
	s.assertAllowlistState(s.curserCapObjectID, false)

	a := s.newAdapter()
	env := s.buildEnv()
	chains := s.buildSuiChains(false)
	subject := a.SelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector)
	curseInput := fastcurse.CurseInput{
		ChainSelector: uint64(s.chainSelector),
		Subjects:      []fastcurse.Subject{subject},
	}

	s.assertIsCursed(a, env, subject, false)

	curseReport, err := cld_ops.ExecuteSequence(s.NewOpBundle(), a.Curse(), chains, curseInput)
	s.Require().NoError(err, "building fast curse MCMS batch operations after deregistration")

	proposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.fastMcmsPackageID,
		MCMSStateObjID:     s.fastMcmsObj,
		TimelockObjID:      s.fastTimelockObj,
		AccountObjID:       s.fastAccountObj,
		RegistryObjID:      s.fastRegistryObj,
		DeployerStateObjID: s.fastDeployerStateObj,
		Description:        "Integration test: fast MCMS curse should fail after allowlist deregistration",
		BatchOp:            curseReport.Output.BatchOps[0],
		TimelockConfig:     utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	s.Require().NoError(err)

	s.executeFastProposalE2eExpectFailure(proposal, s.fastBypasserConfig)
	s.assertIsCursed(a, env, subject, false)
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

func (s *CCIPCurseMCMSTestSuite) executeFastProposalE2eExpectFailure(timelockProposal *mcms.TimelockProposal, roleConfig *RoleConfig) {
	saved := s.snapshotSlowMCMSEndpoints()
	s.mcmsPackageID = s.fastMcmsPackageID
	s.mcmsObj = s.fastMcmsObj
	s.timelockObj = s.fastTimelockObj
	s.registryObj = s.fastRegistryObj
	s.deployerStateObj = s.fastDeployerStateObj
	s.accountObj = s.fastAccountObj
	defer s.restoreSlowMCMSEndpoints(saved)

	proposal := s.ConvertProposal(timelockProposal)
	s.SignProposal(proposal, roleConfig)
	s.SetRoot(proposal, roleConfig)

	encoders, err := proposal.GetEncoders()
	s.Require().NoError(err)
	suiEncoder := encoders[s.chainSelector].(*suisdk.Encoder)
	executor, err := suisdk.NewExecutor(
		s.client,
		s.signer,
		suiEncoder,
		mcmsencoder.NewCCIPEntrypointArgEncoder(s.registryObj, s.deployerStateObj),
		s.mcmsPackageID,
		roleConfig.Role,
		s.mcmsObj,
		s.accountObj,
		s.registryObj,
		s.timelockObj,
	)
	s.Require().NoError(err, "creating executor for fast MCMS contract")

	executors := map[types.ChainSelector]sdk.Executor{
		s.chainSelector: executor,
	}
	executable, err := mcms.NewExecutable(proposal, executors)
	s.Require().NoError(err, "creating fast MCMS executable")

	var execErr error
	for i := range proposal.Operations {
		_, execErr = executable.Execute(s.T().Context(), i)
		if execErr != nil {
			break
		}
	}
	s.Require().Error(execErr, "expected fast MCMS curse proposal to fail after allowlist deregistration")
}
