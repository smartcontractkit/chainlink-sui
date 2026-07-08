package mcmsops

import (
	"context"
	"math/big"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	cselectors "github.com/smartcontractkit/chain-selectors"
	mocksui "github.com/smartcontractkit/mcms/sdk/sui/mocks/sui"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

// setConfigMockClient wraps the MCMS testify mock client so it also satisfies
// the GetMoveModuleFunction shape the Sui bindings occasionally consult.
// Mirrors proposalGenerateMockClient in op_proposal_generate_test.go.
type setConfigMockClient struct {
	*mocksui.SuiPTBClient
}

func (c *setConfigMockClient) GetMoveModuleFunction(
	_ context.Context, _, _, _ string,
) (*suirpcv2.FunctionDescriptor, error) {
	return nil, nil
}

// mockSharedSuiObjectForSetConfig produces a minimal shared-object descriptor
// the ObjectResolver can consume during DevInspect BCS assembly.
func mockSharedSuiObjectForSetConfig(objectID string) *suirpcv2.Object {
	digest := "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT"
	version := uint64(1)
	sharedVersion := uint64(1)
	ownerKind := suirpcv2.Owner_SHARED
	return &suirpcv2.Object{
		ObjectId: &objectID,
		Version:  &version,
		Digest:   &digest,
		Owner: &suirpcv2.Owner{
			Kind:    &ownerKind,
			Version: &sharedVersion,
		},
	}
}

// setupSetConfigMockClient builds a mock Sui client whose SimulatePTB return
// value determines the min_delay the defensive check observes.
func setupSetConfigMockClient(t *testing.T, minDelaySeconds uint64) *setConfigMockClient {
	t.Helper()
	mockClient := mocksui.NewSuiPTBClient(t)
	mockClient.On("ReadObjectId", mock.Anything, mock.Anything).
		Return(
			func(_ context.Context, objectID string) *suirpcv2.Object {
				return mockSharedSuiObjectForSetConfig(objectID)
			},
			func(_ context.Context, _ string) error {
				return nil
			},
		).
		Maybe()
	mockClient.On("GetReferenceGasPrice", mock.Anything).
		Return(big.NewInt(1000), nil).
		Maybe()
	mockClient.On("SimulatePTB", mock.Anything, mock.Anything).
		Return([]any{minDelaySeconds}, nil).
		Maybe()
	return &setConfigMockClient{mockClient}
}

func newSetConfigTestBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()
	registry := cld_ops.NewOperationRegistry(SetConfigMCMSOp.AsUntyped())
	return cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
		cld_ops.WithOperationRegistry(registry),
	)
}

// minimalValidConfig produces a small valid MCMS Config so the handler can
// proceed past DevInspect on the accept/warn paths. Two signers with quorum
// one mirrors the shape used in seq_deploy_test.go.
func minimalValidConfig() types.Config {
	return types.Config{
		Quorum: 1,
		Signers: []common.Address{
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
		},
	}
}

func TestSetConfigMcmsHandler_TimelockObjectID(t *testing.T) {
	t.Parallel()

	const (
		testMcmsPackageID = "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		testOwnerCap      = "0x2222222222222222222222222222222222222222222222222222222222222222"
		testMcmsObjectID  = "0x3333333333333333333333333333333333333333333333333333333333333333"
		testTimelockObjID = "0x4444444444444444444444444444444444444444444444444444444444444444"
	)
	chainSelector := cselectors.SUI_TESTNET.Selector
	maxSeconds := uint64(utils.MaxTimelockScheduleDelay / time.Second)

	buildDeps := func(client *setConfigMockClient) sui_ops.OpTxDeps {
		return sui_ops.OpTxDeps{
			Client: client,
			Signer: nil, // proposal-only mode; handler substitutes a DevInspectSigner
			GetCallOpts: func() *bind.CallOpts {
				return &bind.CallOpts{}
			},
		}
	}

	t.Run("skips defensive read when TimelockObjectID is empty", func(t *testing.T) {
		t.Parallel()
		// SimulatePTB return of ^uint64(0) would trip the cap if the check
		// ran. We assert it does NOT run by expecting a nil error from
		// the defensive block despite the malicious value.
		client := setupSetConfigMockClient(t, ^uint64(0))
		bundle := newSetConfigTestBundle(t)
		input := MCMSSetConfigInput{
			ChainSelector: chainSelector,
			McmsPackageID: testMcmsPackageID,
			OwnerCap:      testOwnerCap,
			McmsObjectID:  testMcmsObjectID,
			// TimelockObjectID left empty on purpose.
			Config:    minimalValidConfig(),
			ClearRoot: true,
		}
		_, err := cld_ops.ExecuteOperation(bundle, SetConfigMCMSOp, buildDeps(client), input)
		require.NoError(t, err)
	})

	t.Run("accepts a min_delay within the cap", func(t *testing.T) {
		t.Parallel()
		client := setupSetConfigMockClient(t, uint64(24*time.Hour/time.Second))
		bundle := newSetConfigTestBundle(t)
		input := MCMSSetConfigInput{
			ChainSelector:    chainSelector,
			McmsPackageID:    testMcmsPackageID,
			OwnerCap:         testOwnerCap,
			McmsObjectID:     testMcmsObjectID,
			TimelockObjectID: testTimelockObjID,
			Config:           minimalValidConfig(),
			ClearRoot:        true,
		}
		_, err := cld_ops.ExecuteOperation(bundle, SetConfigMCMSOp, buildDeps(client), input)
		require.NoError(t, err)
	})

	t.Run("warns but proceeds when min_delay is 0 (F5 bootstrap window)", func(t *testing.T) {
		t.Parallel()
		client := setupSetConfigMockClient(t, 0)
		bundle := newSetConfigTestBundle(t)
		input := MCMSSetConfigInput{
			ChainSelector:    chainSelector,
			McmsPackageID:    testMcmsPackageID,
			OwnerCap:         testOwnerCap,
			McmsObjectID:     testMcmsObjectID,
			TimelockObjectID: testTimelockObjID,
			Config:           minimalValidConfig(),
			ClearRoot:        true,
		}
		_, err := cld_ops.ExecuteOperation(bundle, SetConfigMCMSOp, buildDeps(client), input)
		require.NoError(t, err, "F5 must warn only, not error")
	})

	t.Run("rejects when current min_delay exceeds the cap (F8)", func(t *testing.T) {
		t.Parallel()
		client := setupSetConfigMockClient(t, maxSeconds+1)
		bundle := newSetConfigTestBundle(t)
		input := MCMSSetConfigInput{
			ChainSelector:    chainSelector,
			McmsPackageID:    testMcmsPackageID,
			OwnerCap:         testOwnerCap,
			McmsObjectID:     testMcmsObjectID,
			TimelockObjectID: testTimelockObjID,
			Config:           minimalValidConfig(),
			ClearRoot:        true,
		}
		_, err := cld_ops.ExecuteOperation(bundle, SetConfigMCMSOp, buildDeps(client), input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "F8 defense")
		assert.Contains(t, err.Error(), testTimelockObjID)
	})

	t.Run("rejects on u64 max (the F8 attack value)", func(t *testing.T) {
		t.Parallel()
		client := setupSetConfigMockClient(t, ^uint64(0))
		bundle := newSetConfigTestBundle(t)
		input := MCMSSetConfigInput{
			ChainSelector:    chainSelector,
			McmsPackageID:    testMcmsPackageID,
			OwnerCap:         testOwnerCap,
			McmsObjectID:     testMcmsObjectID,
			TimelockObjectID: testTimelockObjID,
			Config:           minimalValidConfig(),
			ClearRoot:        true,
		}
		_, err := cld_ops.ExecuteOperation(bundle, SetConfigMCMSOp, buildDeps(client), input)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "F8 defense")
	})
}
