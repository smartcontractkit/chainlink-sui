package config

import common "github.com/smartcontractkit/chainlink-common/pkg/types/sui"

//go:fix inline
var (
	PTBChainWriterModuleName = common.PTBChainWriterModuleName
	CCIPExecute              = common.CCIPExecute
	CCIPCommit               = common.CCIPCommit
)

//go:fix inline
type ChainWriterConfig = common.ChainWriterConfig

//go:fix inline
type ChainWriterModule = common.ChainWriterModule

//go:fix inline
type ChainWriterPTBCommand = common.ChainWriterPTBCommand

//go:fix inline
type PrerequisiteObject = common.PrerequisiteObject

//go:fix inline
type ChainWriterFunction = common.ChainWriterFunction

//go:fix inline
type Arguments = common.Arguments

//go:fix inline
type ChainWriterSignal = common.ChainWriterSignal
