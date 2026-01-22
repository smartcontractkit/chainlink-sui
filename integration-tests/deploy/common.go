//go:build integration

package deploy

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
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
	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/changesets"
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
	registry := cld_ops.NewOperationRegistry(opregistry.AllOperations...)

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

// findUnusedManagedTokenMinterCapID finds the mint cap ID that wasn't consumed by the faucet.
// The faucet consumes its mint cap during initialization, so we check which mint caps still
// exist on-chain. The one that exists is the unused one (from ConfigureDeployerAsMinter).
func (s *DeployTestSuite) findUnusedManagedTokenMinterCapID() (string, error) {
	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	if err != nil {
		return "", fmt.Errorf("failed to get addresses: %w", err)
	}

	ctx := s.T().Context()
	var unusedMintCapID string

	// Find all mint caps and check which ones still exist on-chain
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenMinterCapID {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				// Check if this object still exists on-chain (not consumed/deleted by faucet)
				resp, err := s.client.SuiGetObject(ctx, models.SuiGetObjectRequest{
					ObjectId: addr,
					Options: models.SuiObjectDataOptions{
						ShowOwner: true,
						ShowType:  true,
					},
				})
				// If the object exists (no error and no error in response), it's the unused one
				if err == nil && resp.Error == nil && resp.Data != nil {
					unusedMintCapID = addr
					s.T().Logf("Found unused managed token minter cap ID: %s", addr)
					break
				}
			}
		}
	}

	if unusedMintCapID == "" {
		return "", fmt.Errorf("no unused managed token minter cap ID found (all may have been consumed)")
	}

	return unusedMintCapID, nil
}
