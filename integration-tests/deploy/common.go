//go:build integration

package deploy

import (
	"context"

	"github.com/block-vision/sui-go-sdk/sui"
	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/suite"

	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
)

type DeployTestSuite struct {
	suite.Suite
	lggr   logger.Logger
	signer bindutils.SuiSigner
	client sui.ISuiAPI
	env    cldf.Environment

	// Cached deployment addresses
	linkTokenPackageID     string
	linkTokenMetadataID    string
	linkTokenTreasuryCapID string
	ccipPackageID          string
	ccipObjectRef          string
	mcmsPackageID          string
	deployerAddr           string
}

func (s *DeployTestSuite) SetupSuite() {
	s.signer, s.client = testenv.SetupEnvironment(s.T())
	s.lggr = logger.Test(s.T())

	// Setup operation registry
	ops := make([]*cld_ops.Operation[any, any, any], len(opregistry.AllOperations))
	for i := range opregistry.AllOperations {
		ops[i] = &opregistry.AllOperations[i]
	}
	registry := cld_ops.NewOperationRegistry(ops...)

	bundle := cld_ops.NewBundle(
		func() context.Context { return s.T().Context() },
		s.lggr,
		cld_ops.NewMemoryReporter(),
		cld_ops.WithOperationRegistry(registry),
	)

	s.env = cldf.Environment{
		Name:              "test",
		Logger:            s.lggr,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: chain.NewBlockChains(
			map[uint64]chain.BlockChain{
				cselectors.SUI_LOCALNET.Selector: cldfsui.Chain{
					ChainMetadata: cldfsui.ChainMetadata{
						Selector: cselectors.SUI_LOCALNET.Selector,
					},
					Client: s.client,
					Signer: s.signer,
				},
			}),
		OperationsBundle: bundle,
	}
}
