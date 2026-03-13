//go:build integration

package mcms

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/deployment/adapters"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type CCIPCurseMCMSTestSuite struct {
	MCMSTestSuite
}

func (s *CCIPCurseMCMSTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()
}

func (s *CCIPCurseMCMSTestSuite) TestCurseMCMSTest() {
	s.T().Run("Direct curse/uncurse via sequences", func(_ *testing.T) {
		s.testDirectCurseUncurse()
	})
	s.T().Run("MCMS proposal curse/uncurse", func(_ *testing.T) {
		s.testMCMSCurseProposal()
	})
}

// newAdapter creates a CurseAdapter pre-filled from suite deployment state,
// bypassing Initialize() since we already have all object IDs in memory.
func (s *CCIPCurseMCMSTestSuite) newAdapter() *adapters.CurseAdapter {
	return &adapters.CurseAdapter{
		CCIPAddress:          s.ccipPackageId,
		CCIPObjectRef:        s.ccipObjects.CCIPObjectRefObjectId,
		CCIPOwnerCapObjectID: s.ccipObjects.OwnerCapObjectId,
		RouterAddress:        s.ccipRouterPackageId,
		RouterStateObjectID:  s.ccipRouterObjects.RouterStateObjectId,
	}
}

// buildSuiChains creates a BlockChains instance backed by the suite's client.
// When includeSigner is true the chain carries the deployer signer, enabling
// direct on-chain execution inside sequences. Pass false for MCMS proposal
// mode where transactions are only encoded and returned as batch operations.
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

// buildEnv returns a minimal cldf.Environment whose BlockChains is backed by the suite's client.
// The signer is always included because the Sui DevInspect layer requires one even for read-only calls.
func (s *CCIPCurseMCMSTestSuite) buildEnv() cldf.Environment {
	return cldf.Environment{
		BlockChains: s.buildSuiChains(true),
	}
}

// assertIsCursed uses the adapter's IsSubjectCursedOnChain to check and assert the
// expected curse state, keeping assertions consistent with the adapter under test.
func (s *CCIPCurseMCMSTestSuite) assertIsCursed(a *adapters.CurseAdapter, env cldf.Environment, subject fastcurse.Subject, expected bool) {
	s.T().Helper()
	cursed, err := a.IsSubjectCursedOnChain(env, uint64(s.chainSelector), subject)
	s.Require().NoError(err, "IsSubjectCursedOnChain")
	s.Require().Equal(expected, cursed, "unexpected curse state for subject %x", subject)
}

// testDirectCurseUncurse verifies that the adapter's Curse() and Uncurse() sequences
// can execute on-chain when the deployer signer holds the RMN Remote OwnerCap.
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
	s.T().Logf("✅ Subject %x cursed directly on Sui chain", subject)

	s.assertIsCursed(a, env, subject, true)

	_, err = cld_ops.ExecuteSequence(bundle, a.Uncurse(), chains, curseInput)
	s.Require().NoError(err, "executing Uncurse() sequence directly")
	s.T().Logf("✅ Subject %x uncursed directly on Sui chain", subject)

	s.assertIsCursed(a, env, subject, false)
}

// testMCMSCurseProposal verifies that the adapter's Curse() and Uncurse() sequences
// can be driven through an MCMS bypasser proposal after CCIP ownership (including
// the RMN Remote OwnerCap) has been transferred to the MCMS contract.
func (s *CCIPCurseMCMSTestSuite) testMCMSCurseProposal() {
	s.RunOwnershipCCIPTransfer()

	a := s.newAdapter()
	env := s.buildEnv()
	// No signer on chains: sequences encode transactions for MCMS instead of executing directly.
	chains := s.buildSuiChains(false)
	bundle := s.NewOpBundle()

	subject := fastcurse.GlobalCurseSubject()
	curseInput := fastcurse.CurseInput{
		ChainSelector: uint64(s.chainSelector),
		Subjects:      []fastcurse.Subject{subject},
	}

	// --- Curse ---
	curseReport, err := cld_ops.ExecuteSequence(bundle, a.Curse(), chains, curseInput)
	s.Require().NoError(err, "building curse MCMS batch operations")
	s.Require().Len(curseReport.Output.BatchOps, 1, "expected exactly one batch op from Curse()")

	curseProposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.mcmsPackageID,
		MCMSStateObjID:     s.mcmsObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		TimelockObjID:      s.timelockObj,
		DeployerStateObjID: s.deployerStateObj,
		Description:        "Integration test: curse global subject via MCMS bypasser",
		BatchOp:            curseReport.Output.BatchOps[0],
		TimelockConfig: utils.TimelockConfig{
			MCMSAction: types.TimelockActionBypass,
		},
	})
	s.Require().NoError(err, "generating curse MCMS proposal")

	s.ExecuteProposalE2e(curseProposal, s.bypasserConfig, 0)
	s.T().Logf("✅ Global curse executed via MCMS bypasser proposal")

	s.assertIsCursed(a, env, subject, true)

	// --- Uncurse ---
	uncurseReport, err := cld_ops.ExecuteSequence(bundle, a.Uncurse(), chains, curseInput)
	s.Require().NoError(err, "building uncurse MCMS batch operations")
	s.Require().Len(uncurseReport.Output.BatchOps, 1, "expected exactly one batch op from Uncurse()")

	uncurseProposal, err := utils.GenerateProposal(s.T().Context(), utils.GenerateProposalInput{
		ChainSelector:      uint64(s.chainSelector),
		Client:             s.client,
		MCMSPackageID:      s.mcmsPackageID,
		MCMSStateObjID:     s.mcmsObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		TimelockObjID:      s.timelockObj,
		DeployerStateObjID: s.deployerStateObj,
		Description:        "Integration test: uncurse global subject via MCMS bypasser",
		BatchOp:            uncurseReport.Output.BatchOps[0],
		TimelockConfig: utils.TimelockConfig{
			MCMSAction: types.TimelockActionBypass,
		},
	})
	s.Require().NoError(err, "generating uncurse MCMS proposal")

	s.ExecuteProposalE2e(uncurseProposal, s.bypasserConfig, 0)
	s.T().Logf("✅ Global uncurse executed via MCMS bypasser proposal")

	s.assertIsCursed(a, env, subject, false)
}
