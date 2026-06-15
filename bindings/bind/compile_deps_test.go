package bind

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/contracts"
)

func TestOrderPublishDependencies_IncludesMissingFastMCMS(t *testing.T) {
	mcms := "0xaaa"
	fastMcms := "0xbbb"
	std := "0x1"
	sui := "0x2"

	// Sui CLI may only return mcms + framework packages when fast_mcms is omitted.
	deps := orderPublishDependencies(
		[]string{mcms, std, sui},
		[]string{mcms, fastMcms},
	)

	require.Equal(t, []string{mcms, fastMcms, std, sui}, deps)
}

func TestOrderPublishDependencies_Deduplicates(t *testing.T) {
	mcms := "0xaaa"
	fastMcms := "0xbbb"

	deps := orderPublishDependencies(
		[]string{mcms, fastMcms, "0x1"},
		[]string{mcms, fastMcms},
	)

	require.Equal(t, []string{mcms, fastMcms, "0x1"}, deps)
}

func TestRequiredPublishDeps_OnRampIncludesCCIPTransitiveDeps(t *testing.T) {
	deps := orderPublishDependencies(
		[]string{"0xccc", "0x1", "0x2"},
		requiredPublishDeps(contracts.CCIPOnramp, map[string]string{
			"mcms":      "0xaaa",
			"fast_mcms": "0xbbb",
			"ccip":      "0xccc",
		}),
	)

	require.Equal(t, []string{"0xaaa", "0xbbb", "0xccc", "0x1", "0x2"}, deps)
}

func TestEnrichPublishDepsPreservingOrder_InsertsBeforeFramework(t *testing.T) {
	deps := enrichPublishDepsPreservingOrder([]string{"0xaaa", "0x1", "0x2"})
	require.Equal(t, []string{"0xaaa", "0x1", "0x2"}, deps)
}
