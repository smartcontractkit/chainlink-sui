package suierrors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractLockedObjectRef(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		wantOK      bool
		wantObjID   string
		wantVersion uint64
	}{
		{
			name:        "simple lock message",
			msg:         `Object (0xabc123, SequenceNumber(42), o#foo) already locked by a different transaction`,
			wantOK:      true,
			wantObjID:   "0xabc123",
			wantVersion: 42,
		},
		{
			name:        "full Sui equivocation log line (realistic)",
			msg:         `failed to execute transaction: {"code":-32002,"message":"Transaction is rejected as invalid by more than 1/3 of validators by stake (non-retriable). Non-retriable errors: [Object (0x23a4b83340069bd92db7ee2a22994d09f7ff1083af74a9151c9659a5a9662750, SequenceNumber(717214713), o#3R6R2XxDfWT6sXyzz4xX4r1mL6GUHNHCH64vkGwHbBJd) already locked by a different transaction: TransactionDigest(HEf6wjXGWSoemesir2LC7CGnb1c4cPrZruSjTXwpgAuJ) { k#80a11a4a.., k#8501c121.. } with 2768 stake; Object (0xdeadbeef, SequenceNumber(123), o#other) already locked by a different transaction: TransactionDigest(AAA) { k#bbbbbbbb.. } with 10 stake]. Retriable errors: [There are too many transactions pending in consensus { k#81888a6b.. } with 521 stake]."}`,
			wantOK:      true,
			wantObjID:   "0x23a4b83340069bd92db7ee2a22994d09f7ff1083af74a9151c9659a5a9662750",
			wantVersion: 717214713,
		},
		{
			name:        "no object pattern present",
			msg:         `failed to execute transaction: {"code":-32002,"message":"Some other error without Sui object lock pattern"}`,
			wantOK:      false,
			wantObjID:   "",
			wantVersion: 0,
		},
		{
			name: "malformed sequence number",
			msg:  `Object (0xabc123, SequenceNumber(foo), o#foo) already locked by a different transaction`,
			// regex does NOT match because (\d+) requires digits -> we get the "no match" path
			wantOK:      false,
			wantObjID:   "",
			wantVersion: 0,
		},
		{
			name:        "mixed-case hex address",
			msg:         `Object (0xABc123DEF456, SequenceNumber(99), o#bar) already locked by a different transaction`,
			wantOK:      true,
			wantObjID:   "0xABc123DEF456",
			wantVersion: 99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objID, ver, ok := ExtractLockedObjectRef(tt.msg)

			require.Equal(t, tt.wantOK, ok, "ok mismatch")
			require.Equal(t, tt.wantObjID, objID, "objectID mismatch")
			require.Equal(t, tt.wantVersion, ver, "version mismatch")
		})
	}
}
