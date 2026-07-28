package offramp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// TestGetNormalizedModule_ReconstructionMatchesGRPCDecoder proves that the gRPC-backed
// GetNormalizedModule reconstruction (client.FunctionDescriptorToNormalizedFunctionMap) produces a
// JSON-RPC-shaped function map that DecodeParameters decodes to the same argument types as the
// production gRPC decoder DecodeParametersFromFunctionDescriptor. This is the migration safety
// check: the reconstructed map must be a faithful stand-in for the deprecated
// sui_getNormalizedMoveModule response shape.
func TestGetNormalizedModule_ReconstructionMatchesGRPCDecoder(t *testing.T) {
	lggr := logger.Test(t)

	tests := []struct {
		name       string
		parameters []*suirpcv2.OpenSignature
		// wantErrSubstr is non-empty when both decoders must error with this substring; empty when
		// both must succeed and produce equal paramTypes.
		wantErrSubstr string
	}{
		{
			name:       "immutable vector<u8>",
			parameters: []*suirpcv2.OpenSignature{openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), vectorBody(primitiveBody(suirpcv2.OpenSignatureBody_U8)))},
		},
		{
			name:       "nested vector<vector<u8>>",
			parameters: []*suirpcv2.OpenSignature{openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), vectorBody(vectorBody(primitiveBody(suirpcv2.OpenSignatureBody_U8))))},
		},
		{
			name:       "bare u64",
			parameters: []*suirpcv2.OpenSignature{openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), primitiveBody(suirpcv2.OpenSignatureBody_U64))},
		},
		{
			name:       "immutable datatype reference",
			parameters: []*suirpcv2.OpenSignature{openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), datatypeBody("0xccip::state_object::CCIPObjectRef"))},
		},
		{
			name:       "owned datatype treated as reference",
			parameters: []*suirpcv2.OpenSignature{openSig(nil, datatypeBody("0xccip::client::Any2SuiMessage"))},
		},
		{
			name:       "mutable datatype reference",
			parameters: []*suirpcv2.OpenSignature{openSig(AnyPointer(suirpcv2.OpenSignature_MUTABLE), datatypeBody("0xreceiver::dummy_receiver::CCIPReceiverState"))},
		},
		{
			name:          "vector of type parameter is unsupported",
			parameters:    []*suirpcv2.OpenSignature{openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), vectorBody(typeParameterBody(0)))},
			wantErrSubstr: "TypeParameter",
		},
		{
			name:          "top-level type parameter is unsupported",
			parameters:    []*suirpcv2.OpenSignature{openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), typeParameterBody(0))},
			wantErrSubstr: "TypeParameter",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fd := &suirpcv2.FunctionDescriptor{Parameters: tc.parameters}
			funcMap := client.FunctionDescriptorToNormalizedFunctionMap(fd)

			grpcParamTypes, grpcErr := DecodeParametersFromFunctionDescriptor(lggr, fd)
			mapParamTypes, mapErr := DecodeParameters(lggr, funcMap, "parameters")

			if tc.wantErrSubstr != "" {
				require.Error(t, grpcErr, "gRPC decoder should error")
				require.Error(t, mapErr, "reconstruction decoder should error")
				assert.Contains(t, grpcErr.Error(), tc.wantErrSubstr)
				assert.Contains(t, mapErr.Error(), tc.wantErrSubstr)
				return
			}

			require.NoError(t, grpcErr, "gRPC decoder should succeed")
			require.NoError(t, mapErr, "reconstruction decoder should succeed")
			assert.Equal(t, grpcParamTypes, mapParamTypes, "reconstructed map must decode to the same param types as the gRPC descriptor")
		})
	}
}
