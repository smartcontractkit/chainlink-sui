package loop

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-sui/codec/loop"
)

//go:fix inline
const READ_COMPONENTS_COUNT = loop.READ_COMPONENTS_COUNT

//Deprecated: use loop.NewLoopChainReader
//
//go:fix inline
func NewLoopChainReader(log logger.Logger, reader types.ContractReader) types.ContractReader {
	return loop.NewLoopChainReader(log, reader)
}
