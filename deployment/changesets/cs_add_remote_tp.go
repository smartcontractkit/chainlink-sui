package changesets

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
)

type AddRemoteTPConfig struct {
	SuiChainSelector uint64
	TokenPoolTypes   []string

	PoolPackageId          string
	TokenpoolStateObjectId string
	TokenPoolOwnerCapId    string
	CoinObjectTypeArg      string
	RemoteChainSelectors   []uint64
	RemotePoolAddressToAdd []string
}

var _ cldf.ChangeSetV2[AddRemoteTPConfig] = AddRemoteTP{}

// DeployAptosChain deploys Sui chain packages and modules
type AddRemoteTP struct{}

// Apply implements deployment.ChangeSetV2.
func (d AddRemoteTP) Apply(e cldf.Environment, config AddRemoteTPConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.SuiChainSelector]

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	// Todo: validate that len of TokenpoolTypes == RemoteChainSelectors == RemotePoolAddressToAdd
	for _, tokenPoolType := range config.TokenPoolTypes {
		if tokenPoolType == "lnr" {
			for i, chainSelector := range config.RemoteChainSelectors {
				_, err := cld_ops.ExecuteOperation(
					e.OperationsBundle,
					lockreleasetokenpoolops.LockReleaseTokenPoolAddRemotePoolOp,
					deps,
					lockreleasetokenpoolops.LockReleaseTokenPoolAddRemotePoolInput{
						LockReleaseTokenPoolPackageId: config.PoolPackageId,
						CoinObjectTypeArg:             config.CoinObjectTypeArg,
						StateObjectId:                 config.TokenpoolStateObjectId,
						OwnerCap:                      config.TokenPoolOwnerCapId,
						RemoteChainSelector:           chainSelector,
						RemotePoolAddress:             config.RemotePoolAddressToAdd[i], // one address at a time
					},
				)
				if err != nil {
					return cldf.ChangesetOutput{}, err
				}
			}
		}

		if tokenPoolType == "bnm" {
			for i, chainSelector := range config.RemoteChainSelectors {
				_, err := cld_ops.ExecuteOperation(
					e.OperationsBundle,
					burnminttokenpoolops.BurnMintTokenPoolAddRemotePoolOp,
					deps,
					burnminttokenpoolops.BurnMintTokenPoolAddRemotePoolInput{
						BurnMintTokenPoolPackageId: config.PoolPackageId,
						CoinObjectTypeArg:          config.CoinObjectTypeArg,
						StateObjectId:              config.TokenpoolStateObjectId,
						OwnerCap:                   config.TokenPoolOwnerCapId,
						RemoteChainSelector:        chainSelector,
						RemotePoolAddress:          config.RemotePoolAddressToAdd[i], // one address at a time
					},
				)
				if err != nil {
					return cldf.ChangesetOutput{}, err
				}
			}
		}

		if tokenPoolType == "managed" {
			for i, chainSelector := range config.RemoteChainSelectors {

				_, err := cld_ops.ExecuteOperation(
					e.OperationsBundle,
					managedtokenpoolops.ManagedTokenPoolAddRemotePoolOp,
					deps,
					managedtokenpoolops.ManagedTokenPoolAddRemotePoolInput{
						ManagedTokenPoolPackageId: config.PoolPackageId,
						CoinObjectTypeArg:         config.CoinObjectTypeArg,
						StateObjectId:             config.TokenpoolStateObjectId,
						OwnerCap:                  config.TokenPoolOwnerCapId,
						RemoteChainSelector:       chainSelector,
						RemotePoolAddress:         config.RemotePoolAddressToAdd[i], // one address at a time
					},
				)
				if err != nil {
					return cldf.ChangesetOutput{}, err
				}
			}
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d AddRemoteTP) VerifyPreconditions(e cldf.Environment, config AddRemoteTPConfig) error {
	return nil
}
