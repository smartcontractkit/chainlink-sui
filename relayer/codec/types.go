package codec

import (
	common "github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	"github.com/smartcontractkit/chainlink-sui/codec"
)

//go:fix inline
var AccountZero = codec.AccountZero

//go:fix inline
var NormalizeSuiAddress = codec.NormalizeSuiAddress

//go:fix inline
type PTBCommandDependency = common.PTBCommandDependency

//go:fix inline
type PointerTag = common.PointerTag

//go:fix inline
type SuiFunctionParam = common.SuiFunctionParam

//go:fix inline
type SuiPTBCommandType = common.SuiPTBCommandType

const (
	SuiPTBCommandMoveCall = common.SuiPTBCommandMoveCall
	SuiPTBCommandPublish  = common.SuiPTBCommandPublish
	SuiPTBCommandTransfer = common.SuiPTBCommandTransfer
)

//go:fix inline
type ConfigSet = codec.ConfigSet

//go:fix inline
type SourceChainConfigSet = codec.SourceChainConfigSet

//go:fix inline
type SourceChainConfig = codec.SourceChainConfig

//go:fix inline
type ExecutionReport = codec.ExecutionReport

//go:fix inline
type RampMessageHeader = codec.RampMessageHeader

//go:fix inline
type Any2SuiTokenTransfer = codec.Any2SuiTokenTransfer

//go:fix inline
type Any2SuiRampMessage = codec.Any2SuiRampMessage

//go:fix inline
type ExecutionStateChanged = codec.ExecutionStateChanged
