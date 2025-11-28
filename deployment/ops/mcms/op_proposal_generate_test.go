package mcmsops

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	cselectors "github.com/smartcontractkit/chain-selectors"
	mocksui "github.com/smartcontractkit/mcms/sdk/sui/mocks/sui"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func newTestBundle(t *testing.T, registry *cld_ops.OperationRegistry) cld_ops.Bundle {
	t.Helper()
	reporter := cld_ops.NewMemoryReporter()
	return cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		reporter,
		cld_ops.WithOperationRegistry(registry),
	)
}

func TestMCMSDynamicProposalGenerateSeq(t *testing.T) {
	t.Parallel()

	// Create a registry with state object operations that support exporting the Call
	registry := cld_ops.NewOperationRegistry(
		ccipops.AddPackageIdStateObjectOp.AsUntyped(),
		ccipops.RemovePackageIdStateObjectOp.AsUntyped(),
		ccipops.TransferOwnershipStateObjectOp.AsUntyped(),
		ccipops.AcceptOwnershipStateObjectOp.AsUntyped(),
	)

	mockClient := mocksui.NewISuiAPI(t)
	// This response doesn't matter much
	mockClient.EXPECT().SuiGetObject(mock.Anything, mock.Anything).
		Return(models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				ObjectId: "0xf2facb344885659b11e707838ee131b407654f75f6589984af462c13de41ef84",
				Version:  "3",
				Digest:   "4TRR2ZC9r7UUDUeke2DUhHdRQkZWYjkygHrRSNVM4YmX",
				Owner:    nil,
			},
			Error: nil,
		}, nil)
	// This is the response from getOpCount
	mockClient.EXPECT().SuiDevInspectTransactionBlock(mock.Anything, mock.Anything).
		Return(models.SuiTransactionBlockResponse{
			Effects: models.SuiEffects{
				Status: models.ExecutionStatus{
					Status: "success",
					Error:  "",
				},
			},
			Results: json.RawMessage(`[{"returnValues":[[[1,0,0,0,0,0,0,0],"u64"]]}]`), // Returns 1
		}, nil)
	// Create mock dependencies
	deps := sui_ops.OpTxDeps{
		Client: mockClient,
		Signer: nil, // We don't need a real signer for this test since NoExecute=true
		GetCallOpts: func() *bind.CallOpts {
			return &bind.CallOpts{}
		},
	}

	// Test data
	testCCIPPackageId := "0x1234567890abcdef"
	testObjectRefId := "0xabcdef1234567890"
	testOwnerCapId := "0x9876543210fedcba"
	testPackageId := "0xdeadbeefcafebabe"
	testNewOwner := "0x1111111111111111"
	testTimelockObjID := "0x2222222222222222"
	testAccountObjID := "0x3333333333333333"
	testRegistryObjID := "0x4444444444444444"
	testDeployerStateObjID := "0x5555555555555555"
	testChainSelector := cselectors.SUI_TESTNET.Selector

	t.Run("Generate Proposal with Multiple Operations - Proposer Role", func(t *testing.T) {
		// Create operation definitions and inputs
		defs := []cld_ops.Definition{
			ccipops.AddPackageIdStateObjectOp.Def(),
			ccipops.TransferOwnershipStateObjectOp.Def(),
			ccipops.AcceptOwnershipStateObjectOp.Def(),
		}

		inputs := []any{
			ccipops.AddPackageIdStateObjectInput{
				CCIPPackageId:         testCCIPPackageId,
				CCIPObjectRefObjectId: testObjectRefId,
				OwnerCapObjectId:      testOwnerCapId,
				PackageId:             testPackageId,
			},
			ccipops.TransferOwnershipStateObjectInput{
				CCIPPackageId:         testCCIPPackageId,
				CCIPObjectRefObjectId: testObjectRefId,
				OwnerCapObjectId:      testOwnerCapId,
				To:                    testNewOwner,
			},
			ccipops.AcceptOwnershipStateObjectInput{
				CCIPPackageId:         testCCIPPackageId,
				CCIPObjectRefObjectId: testObjectRefId,
			},
		}

		proposalInput := ProposalGenerateInput{
			Defs:               defs,
			Inputs:             inputs,
			ChainSelector:      testChainSelector,
			MmcsPackageID:      testCCIPPackageId,
			McmsStateObjID:     testObjectRefId,
			TimelockObjID:      testTimelockObjID,
			AccountObjID:       testAccountObjID,
			RegistryObjID:      testRegistryObjID,
			DeployerStateObjID: testDeployerStateObjID,
			TimelockConfig: utils.TimelockConfig{
				MCMSAction:   types.TimelockActionSchedule,
				MinDelay:     time.Hour * 24,
				OverrideRoot: false,
			},
		}

		// Execute the operation
		bundle := newTestBundle(t, registry)
		result, err := cld_ops.ExecuteSequence(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		require.NoError(t, err, "should generate proposal successfully")

		// Verify the proposal structure
		proposal := result.Output
		assert.Equal(t, "v1", proposal.Version, "proposal version should be v1")
		assert.NotZero(t, proposal.ValidUntil, "proposal should have valid until timestamp")
		assert.Contains(t, proposal.Description, "Invokes the following set of operations", "description should contain operation description")
		assert.Contains(t, proposal.Description, ccipops.AddPackageIdStateObjectOp.Def().ID, "description should contain first operation")
		assert.Contains(t, proposal.Description, ccipops.TransferOwnershipStateObjectOp.Def().ID, "description should contain second operation")
		assert.Contains(t, proposal.Description, ccipops.AcceptOwnershipStateObjectOp.Def().ID, "description should contain third operation")

		// Verify timelock addresses
		require.Len(t, proposal.TimelockAddresses, 1, "should have one timelock address")
		assert.Equal(t, testTimelockObjID, proposal.TimelockAddresses[types.ChainSelector(testChainSelector)], "timelock address should match")

		// Verify chain metadata
		require.Len(t, proposal.ChainMetadata, 1, "should have one chain metadata")
		require.Equal(t, proposal.ChainMetadata[types.ChainSelector(testChainSelector)].StartingOpCount, uint64(1), "starting op count should be 1 as mocked")
		// Note: ChainMetadata structure verification simplified for test

		// Verify operations
		require.Len(t, proposal.Operations, 1, "should have one batch operation")
		batchOp := proposal.Operations[0]
		assert.Equal(t, testChainSelector, uint64(batchOp.ChainSelector), "batch operation chain selector should match")
		assert.Len(t, batchOp.Transactions, 3, "batch operation should contain 3 transactions")

		// Verify delay is set for proposer role
		assert.NotZero(t, proposal.Delay, "delay should be set for proposer role")
		// Note: Delay verification simplified for test

	})

	t.Run("Generate Proposal with Single Operation - Bypasser Role", func(t *testing.T) {
		// Create operation definitions and inputs for single operation
		defs := []cld_ops.Definition{
			ccipops.RemovePackageIdStateObjectOp.Def(),
		}

		inputs := []any{

			ccipops.RemovePackageIdStateObjectInput{
				CCIPPackageId:         testCCIPPackageId,
				CCIPObjectRefObjectId: testObjectRefId,
				OwnerCapObjectId:      testOwnerCapId,
				PackageId:             testPackageId,
			},
		}

		proposalInput := ProposalGenerateInput{
			Defs:               defs,
			Inputs:             inputs,
			MmcsPackageID:      testCCIPPackageId,
			McmsStateObjID:     testObjectRefId,
			TimelockObjID:      testTimelockObjID,
			AccountObjID:       testAccountObjID,
			RegistryObjID:      testRegistryObjID,
			DeployerStateObjID: testDeployerStateObjID,
			ChainSelector:      testChainSelector,
			TimelockConfig: utils.TimelockConfig{
				MCMSAction:   types.TimelockActionBypass,
				MinDelay:     0,
				OverrideRoot: false,
			},
		}

		// Execute the operation
		bundle := newTestBundle(t, registry)
		result, err := cld_ops.ExecuteSequence(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		require.NoError(t, err, "should generate proposal successfully")

		// Verify the proposal structure
		proposal := result.Output
		assert.Equal(t, "v1", proposal.Version, "proposal version should be v1")
		assert.Contains(t, proposal.Description, ccipops.RemovePackageIdStateObjectOp.Def().ID, "description should contain remove operation")

		// Verify no delay is set for bypasser role
		assert.Zero(t, proposal.Delay, "delay should not be set for bypasser role")

		// Verify single transaction
		require.Len(t, proposal.Operations, 1, "should have one batch operation")
		batchOp := proposal.Operations[0]
		assert.Len(t, batchOp.Transactions, 1, "batch operation should contain 1 transaction")
	})

	t.Run("Generate Proposal with Invalid Role", func(t *testing.T) {
		// Create a proposal with invalid role
		defs := []cld_ops.Definition{
			ccipops.AddPackageIdStateObjectOp.Def(),
		}

		inputs := []any{

			ccipops.AddPackageIdStateObjectInput{
				CCIPPackageId:         testCCIPPackageId,
				CCIPObjectRefObjectId: testObjectRefId,
				OwnerCapObjectId:      testOwnerCapId,
				PackageId:             testPackageId,
			},
		}

		proposalInput := ProposalGenerateInput{
			Defs:   defs,
			Inputs: inputs,

			MmcsPackageID:      testCCIPPackageId,
			McmsStateObjID:     testObjectRefId,
			TimelockObjID:      testTimelockObjID,
			AccountObjID:       testAccountObjID,
			RegistryObjID:      testRegistryObjID,
			DeployerStateObjID: testDeployerStateObjID,
			ChainSelector:      testChainSelector,
			TimelockConfig: utils.TimelockConfig{
				MCMSAction:   "bad_action",
				MinDelay:     time.Hour,
				OverrideRoot: false,
			},
		}

		// Execute the operation - should fail
		bundle := newTestBundle(t, registry)
		_, err := cld_ops.ExecuteSequence(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		require.Error(t, err, "should fail with invalid action")
		assert.Contains(t, err.Error(), "unsupported action", "error should mention `unsupported action`")
	})

	t.Run("Generate Proposal with Mismatched Definitions and Inputs", func(t *testing.T) {
		// Create mismatched definitions and inputs (more defs than inputs)
		defs := []cld_ops.Definition{
			ccipops.AddPackageIdStateObjectOp.Def(),
			ccipops.RemovePackageIdStateObjectOp.Def(),
		}

		inputs := []any{

			ccipops.AddPackageIdStateObjectInput{
				CCIPPackageId:         testCCIPPackageId,
				CCIPObjectRefObjectId: testObjectRefId,
				OwnerCapObjectId:      testOwnerCapId,
				PackageId:             testPackageId,
			},

			// Missing second input
		}

		proposalInput := ProposalGenerateInput{
			Defs:   defs,
			Inputs: inputs[:1], // Only one input for two definitions

			MmcsPackageID:      testCCIPPackageId,
			McmsStateObjID:     testObjectRefId,
			TimelockObjID:      testTimelockObjID,
			AccountObjID:       testAccountObjID,
			RegistryObjID:      testRegistryObjID,
			DeployerStateObjID: testDeployerStateObjID,
			ChainSelector:      testChainSelector,
			TimelockConfig: utils.TimelockConfig{
				MCMSAction:   types.TimelockActionSchedule,
				MinDelay:     time.Hour,
				OverrideRoot: false,
			},
		}

		// Execute the operation - should fail due to index out of bounds
		bundle := newTestBundle(t, registry)
		_, err := cld_ops.ExecuteSequence(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		require.Error(t, err, "should fail with mismatched definitions and inputs")
	})
}
