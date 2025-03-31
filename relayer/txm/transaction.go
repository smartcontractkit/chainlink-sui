package txm

import commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

type SuiTx struct {
	transactionID string
	sender        string
	metadata      *commontypes.TxMeta
	timestamp     uint64
	payload       []byte
	attempt       int
	state         string
}
