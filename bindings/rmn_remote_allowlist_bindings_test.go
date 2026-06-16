package mcmsencoder

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
)

const (
	testAllowlistCCIPPackageID = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testAllowlistStateObjectID = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	testAllowlistCurserCapID   = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
)

func TestRmnRemoteAllowlistViewEncoders(t *testing.T) {
	t.Parallel()

	ref := bind.Object{Id: testAllowlistStateObjectID}

	contract, err := module_rmn_remote.NewRmnRemote(testAllowlistCCIPPackageID, nil)
	require.NoError(t, err)

	allowedCall, err := contract.Encoder().IsCurserCapAllowed(ref, testAllowlistCurserCapID)
	require.NoError(t, err)
	require.Equal(t, "rmn_remote", allowedCall.Module.ModuleName)
	require.Equal(t, "is_curser_cap_allowed", allowedCall.Function)

	idsCall, err := contract.Encoder().GetAllowedCurserCapIds(ref)
	require.NoError(t, err)
	require.Equal(t, "get_allowed_curser_cap_ids", idsCall.Function)
}

func TestRmnRemoteAllowlistAdminEncoders(t *testing.T) {
	t.Parallel()

	ref := bind.Object{Id: testAllowlistStateObjectID}
	ownerCap := bind.Object{Id: "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"}

	contract, err := module_rmn_remote.NewRmnRemote(testAllowlistCCIPPackageID, nil)
	require.NoError(t, err)

	initCall, err := contract.Encoder().InitializeAllowedCurserCaps(ref, ownerCap, []string{testAllowlistCurserCapID})
	require.NoError(t, err)
	require.Equal(t, "initialize_allowed_curser_caps", initCall.Function)

	registerCall, err := contract.Encoder().RegisterCurserCapIds(ref, ownerCap, []string{testAllowlistCurserCapID})
	require.NoError(t, err)
	require.Equal(t, "register_curser_cap_ids", registerCall.Function)

	deregisterCall, err := contract.Encoder().DeregisterCurserCapIds(ref, ownerCap, []string{testAllowlistCurserCapID})
	require.NoError(t, err)
	require.Equal(t, "deregister_curser_cap_ids", deregisterCall.Function)
}
