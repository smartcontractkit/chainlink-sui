package codec

import "github.com/smartcontractkit/chainlink-sui/codec"

//go:fix inline
var AccountZero = codec.AccountZero

//go:fix inline
type PTBCommandDependency = codec.PTBCommandDependency

//go:fix inline
type PointerTag = codec.PointerTag

//go:fix inline
type SuiFunctionParam = codec.SuiFunctionParam

//go:fix inline
type SuiPTBCommandType = codec.SuiPTBCommandType

const (
	SuiPTBCommandMoveCall = codec.SuiPTBCommandMoveCall
	SuiPTBCommandPublish  = codec.SuiPTBCommandPublish
	SuiPTBCommandTransfer = codec.SuiPTBCommandTransfer
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
