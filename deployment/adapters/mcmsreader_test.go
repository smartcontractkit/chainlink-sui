package adapters_test

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	ccipmcms "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/adapters"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

const (
	testMcmsReaderCCIPPackageID = "0x1111111111111111111111111111111111111111111111111111111111111111"
	testMcmsReaderCCIPObjectRef = "0x2222222222222222222222222222222222222222222222222222222222222222"

	testFastMcmsPackageID = "0x1212121212121212121212121212121212121212121212121212121212121212"
	testFastMcmsState     = "0x2323232323232323232323232323232323232323232323232323232323232323"
	testFastMcmsRegistry  = "0x3434343434343434343434343434343434343434343434343434343434343434"
	testFastMcmsAccount   = "0x4545454545454545454545454545454545454545454545454545454545454545"
	testFastMcmsOwnerCap  = "0x5656565656565656565656565656565656565656565656565656565656565656"
	testFastMcmsTimelock  = "0x6767676767676767676767676767676767676767676767676767676767676767"
	testFastMcmsDeployer  = "0x7878787878787878787878787878787878787878787878787878787878787878"
)

func TestMCMSReader_RMNMCMSQualifier_SelectsFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, ab.Save(selector, testMcmsReaderCCIPPackageID, cldf.NewTypeAndVersion(deployment.SuiCCIPType, deployment.Version1_0_0)))
	require.NoError(t, ab.Save(selector, testMcmsReaderCCIPObjectRef, cldf.NewTypeAndVersion(deployment.SuiCCIPObjectRefType, deployment.Version1_0_0)))
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: testFastMcmsPackageID,
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   testFastMcmsState,
			McmsRegistryObjectId:        testFastMcmsRegistry,
			McmsAccountStateObjectId:    testFastMcmsAccount,
			McmsAccountOwnerCapObjectId: testFastMcmsOwnerCap,
			TimelockObjectId:            testFastMcmsTimelock,
			McmsDeployerStateObjectId:   testFastMcmsDeployer,
		},
	}, deployment.MCMSInstanceFastCurse))

	env := cldf.Environment{
		ExistingAddresses: ab,
		BlockChains: cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
			selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
		}),
	}

	reader := &adapters.MCMSReader{}
	input := ccipmcms.Input{
		Qualifier:      "RMNMCMS",
		TimelockAction: mcmstypes.TimelockActionBypass,
	}

	timelockRef, err := reader.GetTimelockRef(env, selector, input)
	require.NoError(t, err)
	require.Equal(t, testFastMcmsTimelock, timelockRef.Address)
	require.Equal(t, selector, timelockRef.ChainSelector)

	mcmsRef, err := reader.GetMCMSRef(env, selector, input)
	require.NoError(t, err)
	require.Equal(t, testFastMcmsState, mcmsRef.Address)
	require.Equal(t, selector, mcmsRef.ChainSelector)
}
