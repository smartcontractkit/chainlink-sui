//go:build integration

package mcms

import (
	"crypto/ecdsa"
	"slices"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	cselectors "github.com/smartcontractkit/chain-selectors"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

type RoleConfig struct {
	Count  int
	Quorum uint8
	Keys   []*ecdsa.PrivateKey
	Config *types.Config
}

func CreateConfig(count int, quorum uint8) *RoleConfig {
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
			b := uint64(500_000_000) // needs to be high for publishing
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		s.T().Context,
		logger.Test(s.T()),
		reporter,
	)

	bypasserCount := 2
	bypasserQuorum := 2
	bypasserConfig := CreateConfig(bypasserCount, uint8(bypasserQuorum))
	proposerCount := 3
	proposerQuorum := 2
	proposerConfig := CreateConfig(proposerCount, uint8(proposerQuorum))

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
}
