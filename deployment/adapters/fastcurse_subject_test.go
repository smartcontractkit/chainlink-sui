package adapters

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/stretchr/testify/require"
)

func TestSubjectInCursedSubjects(t *testing.T) {
	t.Parallel()

	laneSubject := fastcurse.GenericSelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector)
	globalSubject := fastcurse.GlobalCurseSubject()
	otherLaneSubject := fastcurse.GenericSelectorToSubject(cselectors.SUI_TESTNET.Selector)

	tests := []struct {
		name           string
		cursedSubjects [][]byte
		subject        fastcurse.Subject
		want           bool
	}{
		{
			name:           "lane subject explicitly cursed",
			cursedSubjects: [][]byte{laneSubject[:]},
			subject:        laneSubject,
			want:           true,
		},
		{
			name:           "lane subject not in set",
			cursedSubjects: [][]byte{otherLaneSubject[:]},
			subject:        laneSubject,
			want:           false,
		},
		{
			name:           "lane subject not cursed when only global is in set",
			cursedSubjects: [][]byte{globalSubject[:]},
			subject:        laneSubject,
			want:           false,
		},
		{
			name:           "global subject in set",
			cursedSubjects: [][]byte{globalSubject[:]},
			subject:        globalSubject,
			want:           true,
		},
		{
			name:           "empty set",
			cursedSubjects: nil,
			subject:        laneSubject,
			want:           false,
		},
		{
			name:           "ignores wrong length entries",
			cursedSubjects: [][]byte{laneSubject[:8], laneSubject[:]},
			subject:        laneSubject,
			want:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, subjectInCursedSubjects(tt.cursedSubjects, tt.subject))
		})
	}
}

func TestSubjectInCursedSubjects_GlobalCurseDoesNotImplyLaneCursed(t *testing.T) {
	t.Parallel()

	globalSubject := fastcurse.GlobalCurseSubject()
	globalOnly := [][]byte{globalSubject[:]}
	laneSubject := fastcurse.GenericSelectorToSubject(cselectors.SUI_TESTNET.Selector)

	require.False(t, subjectInCursedSubjects(globalOnly, laneSubject))
	require.True(t, subjectInCursedSubjects(globalOnly, globalSubject))
}
