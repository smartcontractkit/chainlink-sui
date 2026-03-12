package deployment

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

func TestLoadOnchainStatesui(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		env     cldf.Environment
		want    map[uint64]CCIPChainState
		wantErr bool
	}{
		{
			name: "load sui state successfully",
			env: cldf.Environment{
				Name:              "test",
				ExistingAddresses: loadTestAddressBook(t),
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						9762610643973837292: sui.Chain{},
					}),
			},
			want:    map[uint64]CCIPChainState{9762610643973837292: getExpectedSuiChainState()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := LoadOnchainStatesui(tt.env)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("LoadOnchainStatesui() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("LoadOnchainStatesui() succeeded unexpectedly")
			}

			// Use cmp.Diff with SortSlices to handle unordered arrays
			if diff := cmp.Diff(tt.want[9762610643973837292], got[9762610643973837292],
				cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
				t.Errorf("LoadOnchainStatesui() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func loadTestAddressBook(t *testing.T) cldf.AddressBook {
	filePath := "testdata/addresses.json"
	b, err := os.ReadFile(filePath)
	require.NoError(t, err, "failed to read test address book file")

	addrsByChain := make(map[uint64]map[string]cldf.TypeAndVersion)
	err = json.Unmarshal(b, &addrsByChain)
	require.NoError(t, err, "failed to unmarshal JSON")

	return cldf.NewMemoryAddressBookFromMap(addrsByChain)
}

func getExpectedSuiChainState() CCIPChainState {
	return CCIPChainState{
		MCMSPackageID:               "0x56c76b6e9b53b071cec8d50f2c18031f52423cb0f31bbb573fe202364deb82e3",
		MCMSStateObjectID:           "0xf59eadb8d872202e6ebfdc714800c1a6457ac5fd7e1f9a593370a6f98f46dc30",
		MCMSRegistryObjectID:        "0xd545d92abb3e7d5fe488fdfc0653169f8dbdcbce294a7c80e96a05fd890d5249",
		MCMSAccountStateObjectID:    "0xe4ff650466b35c592f255637c87242ee2b03ca3ba736eacfbb886487403a6512",
		MCMSAccountOwnerCapObjectID: "0xea8150e8d66c6aaad16a649b0302a9baeadf3fb6d8f4da2380a7332226e13f52",
		MCMSTimelockObjectID:        "0xaba16fb3b71f52b2d3cb668dd5a8bd600eb87345ea7ae7ed819dab800c2bc0f3",
		CCIPAddress:                 "0xece742a763bddf1e36629fa06b605497e413241afd14f05e558e80eef4f64e95",
		CCIPObjectRef:               "0xbeace36c3c1e1f37c5806c4954140d15cbc7b0002ef7ccb490de26e82f5ec4ca",
		CCIPOwnerCapObjectId:        "0x874447a3ac6ae37bf545c6d71b0fa4a3d0af56a3f339faf4a4bab16ca4956ce7",
		CCIPUpgradeCapObjectId:      "0x99aef4fefc921681d57823e479fc424f43457ef9ebdf6cf1cbd91379c27ef19c",
		CCIPRouterAddress:           "0xed4613bd35004954c07150c3e9b10230f5e23e3058bc2ca0e3e676cb43eb4dc1",
		CCIPRouterStateObjectID:     "0xbb2486d233b0d358f82fb8c4c5c75881e65069ea8ebe5ab692a636c9e0eff7cd",
		FeeQuoterCapId:              "0xceac727ef0d9a8494323478ece8b883877d9178a573e1678565e5195eac878e8",
		OnRampAddress:               "0xf87c6010be571a304f0d860857204bc66f037842156f0f6c9d80be265fd83752",
		OnRampStateObjectId:         "0x75ec1e10b4302f7c69476eb196c88a0aa43a4d509bbbe5cc1feb213e4b6dd58b",
		OnRampOwnerCapObjectId:      "0x8101795ff02d4935a05fb519e1b21b83855a970639cf0c28fa7a51f7d2e689ae",
		OnRampUpgradeCapId:          "0x9495408a99b47022e94ef655c107cb0bff96638c8763e05d2239f661f48d47df",
		OffRampAddress:              "0x9438693fb18f5660aff9277240a2282be44dc01cdd7eed4e1d8de0591ad52c03",
		OffRampOwnerCapId:           "0x8f7eb5b7879449519b39db447b94c87b6bb88f0a1deff40b3e0ffef3d1058f69",
		OffRampUpgradeCapId:         "0x012fd673ed8cda310acd840635bfed201bcd76a446515eb4c569fa2a4155d75b",
		OffRampStateObjectId:        "0x4fdacea0d627df26a6f34ac62952ed8d3c32ea70aacb90a7f2134b39e36a79cd",
		LinkTokenAddress:            "0x59d6fe83c19eb26733acb19cd32522277319a7f2accd360af7c04285c015375d",
		LinkTokenCoinMetadataId:     "0x8afb916ec72b91d28f519539659ebb1200b1824ff1f8d4c8f433acbb03017f2f",
		LinkTokenTreasuryCapId:      "0x60f35af35748a3f2f4324c2646cf9187d111ed17cb204043c1cf85c0f83dde0f",
		ManagedTokens: map[string]ManagedTokenState{
			"LINK": {
				PackageID:        "0xa1ec7fc00a6f40db9693ad1415d0c193ad3906494428cf252621037bd7117e29",
				StateObjectId:    "0xd93ceb06feeab8b0b5333a7f76f2e9d0d79e71063b8a71bce654954f7bfd631d",
				OwnerCapObjectId: "0x0eac90f883ea0902156b340c41e0d178af3f012fe63c33f9663b7be245bee4fe",
				MinterCapObjectIds: []string{
					"0x125afa59deccfec819aaa67bf14f049982d2a1ca87c49c8614da5ea2dc438f78",
					"0xe25afa59deccfec819aaa67bf14f049982d2a1ca87c49c8614da5ea2dc438f72",
				},
				PublisherObjectId: "0x52f33e4724128431084e207406304383c05660042000016240438147268f184e",
			},
			"CCIP BnM": {
				PackageID:        "0xa1ec7fc00a6f40db9693ad1415d0c193ad3906494428cf252621037bd7117e2a",
				StateObjectId:    "0xd93ceb06feeab8b0b5333a7f76f2e9d0d79e71063b8a71bce654954f7bfd631e",
				OwnerCapObjectId: "0x0eac90f883ea0902156b340c41e0d178af3f012fe63c33f9663b7be245bee4ff",
				MinterCapObjectIds: []string{
					"0x125afa59deccfec819aaa67bf14f049982d2a1ca87c49c8614da5ea2dc438f71",
					"0xe25afa59deccfec819aaa67bf14f049982d2a1ca87c49c8614da5ea2dc438f7b",
				},
				PublisherObjectId: "0x52f33e4724128431084e207406304383c05660042000016240438147268f184a",
			},
		},
		ManagedTokenFaucets: map[string]ManagedTokenFaucetState{
			"LINK": {
				PackageID:          "0x52f33e4724128431084e207406304383c05660042000016240438147268f1851",
				StateObjectId:      "0x52f33e4724128431084e207406304383c05660042000016240438147268f184f",
				UpgradeCapObjectId: "0x52f33e4724128431084e207406304383c05660042000016240438147268f1850",
			},
		},
		LnRTokenPools: map[string]CCIPPoolState{
			"CCIP-LnR": {
				PackageID:        "0x90f10215a219a1f2e30f746e68c1e4f2b39f593d83a1980c7cf4e4738591a7e6",
				StateObjectId:    "0x4e6c205b2f8c5159795f536035c573338fbb345db473829a95fceb0f82890fd4",
				OwnerCapObjectId: "0x9aeef2d381d70fe075970961c15dcb4812e36a5520549f4cd559150589afdd98",
				RebalancerCapIds: []string{"0x9aeef2d381d70fe075970961c15dcb4812e36a5520549f4cd559150589afdd97", "0x9aeef2d381d70fe075970961c15dcb4812e36a5520549f4cd559150589afdd99"},
			},
		},
		BnMTokenPools: map[string]CCIPPoolState{
			"CCIP-BnM": {
				PackageID:        "0x70f10215a219a1f2e30f746e68c1e4f2b39f593d83a1980c7cf4e4738591a7e5",
				StateObjectId:    "0x2d5b1f5a1e7b4049694e425f24b46e227eaf234ca362719394ebd0ae71790fc2",
				OwnerCapObjectId: "0x78dde1c270c6edf64e586850b04cba3701d259441f4390e3bc4490407a9ecc86",
			},
			"CCIP-BnM2": {
				PackageID:        "0x81023260b32ab203041857f07dcf0349ef59ec45039b21ad8d051846ab8f6",
				StateObjectId:    "0x3e6c205b2f8c5159795f536035c573338fbb345db473829a95fceb0f82890fd3",
				OwnerCapObjectId: "0x89eef2d381d70fe075970961c15dcb4812e36a5520549f4cd559150589afdd97",
			},
		},
		ManagedTokenPools: map[string]CCIPPoolState{
			"LINK": {
				PackageID:        "0xa10215a219a1f2e30f746e68c1e4f2b39f593d83a1980c7cf4e4738591a7e7",
				StateObjectId:    "0x5f6c205b2f8c5159795f536035c573338fbb345db473829a95fceb0f82890fd5",
				OwnerCapObjectId: "0xabeef2d381d70fe075970961c15dcb4812e36a5520549f4cd559150589afdd99",
			},
		},
	}
}
