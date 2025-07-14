package offramp

import (
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

const COMMIT_PTB_NAME = "ccip_commit"
const EXECUTE_PTB_NAME = "CCIPExecuteReport"

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func GenerateCommitPTB(
	lggr logger.Logger,
	offRampPackageId string,
	signerPublicKey []byte,
) *config.ChainWriterFunction {
	return &config.ChainWriterFunction{
		Name:      COMMIT_PTB_NAME,
		PublicKey: signerPublicKey,
		Params:    []codec.SuiFunctionParam{},
		PTBCommands: []config.ChainWriterPTBCommand{
			{
				Type:      codec.SuiPTBCommandMoveCall,
				PackageId: strPtr(offRampPackageId),
				ModuleId:  strPtr("offramp"),
				Function:  strPtr("commit"),
				Params: []codec.SuiFunctionParam{
					{
						Name:      "ccip_object_ref",
						Type:      "object_id",
						Required:  true,
						IsMutable: boolPtr(true),
					},
					{
						Name:      "state",
						Type:      "object_id",
						Required:  true,
						IsMutable: boolPtr(true),
					},
					{
						Name:      "clock",
						Type:      "object_id",
						Required:  true,
						IsMutable: boolPtr(false),
					},
					{
						Name:     "report_context",
						Type:     "vector<vector<u8>>",
						Required: true,
					},
					{
						Name:     "report",
						Type:     "vector<u8>",
						Required: true,
					},
					{
						Name:     "signatures",
						Type:     "vector<vector<u8>>",
						Required: true,
					},
				},
			},
		},
	}
}

func GenerateExecutePTB(
	lggr logger.Logger,
	offRampPackageId string,
	signerPublicKey []byte,
) *config.ChainWriterFunction {
	return &config.ChainWriterFunction{
		Name:      EXECUTE_PTB_NAME,
		PublicKey: signerPublicKey,
		Params:    []codec.SuiFunctionParam{},
		PTBCommands: []config.ChainWriterPTBCommand{
			{
				Type:      codec.SuiPTBCommandMoveCall,
				PackageId: strPtr(offRampPackageId),
				ModuleId:  strPtr("offramp"),
				Function:  strPtr("init_execute"),
				Params: []codec.SuiFunctionParam{
					{
						Name:      "ccip_object_ref",
						Type:      "object_id",
						Required:  true,
						IsMutable: boolPtr(false),
					},
					{
						Name:     "state",
						Type:     "object_id",
						Required: true,
					},
					{
						Name:      "clock",
						Type:      "object_id",
						Required:  true,
						IsMutable: boolPtr(false),
					},
					{
						Name:     "report_context",
						Type:     "vector<vector<u8>>",
						Required: true,
					},
					{
						Name:     "report",
						Type:     "vector<u8>",
						Required: true,
					},
				},
			},
			{
				Type:      codec.SuiPTBCommandMoveCall,
				PackageId: strPtr(offRampPackageId),
				ModuleId:  strPtr("offramp"),
				Function:  strPtr("finish_execute"),
				Params: []codec.SuiFunctionParam{
					{
						Name:      "state",
						Type:      "object_id",
						Required:  true,
						IsMutable: boolPtr(true),
					},
					{
						Name:     "receiver_params",
						Type:     "ptb_dependency",
						Required: true,
						PTBDependency: &codec.PTBCommandDependency{
							CommandIndex: uint16(0),
						},
					},
				},
			},
		},
	}
}
