package mcmsops

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCMSDynamicProposalGenerateSeq(t *testing.T) {
	t.Parallel()

	// Create a registry with state object operations that support exporting the Call
	registry := cld_ops.NewOperationRegistry(
		ccipops.AddPackageIdStateObjectOp.AsUntyped(),
		ccipops.RemovePackageIdStateObjectOp.AsUntyped(),
		ccipops.TransferOwnershipStateObjectOp.AsUntyped(),
		ccipops.AcceptOwnershipStateObjectOp.AsUntyped(),
	)

	// Create mock dependencies
	deps := sui_ops.OpTxDeps{
		Client: nil, // We don't need a real client for this test since NoExecute=true
		Signer: nil, // We don't need a real signer for this test since NoExecute=true
		GetCallOpts: func() *bind.CallOpts {
			return &bind.CallOpts{}
		},
	}

	// Create test bundle
	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
		cld_ops.WithOperationRegistry(registry),
	)
	// Test data
	testCCIPPackageId := "0x1234567890abcdef"
	testObjectRefId := "0xabcdef1234567890"
	testOwnerCapId := "0x9876543210fedcba"
	testPackageId := "0xdeadbeefcafebabe"
	testNewOwner := "0x1111111111111111"
	testTimelockObjID := "0x2222222222222222"
	testAccountObjID := "0x3333333333333333"
	testRegistryObjID := "0x4444444444444444"
	testChainSelector := uint64(123456789)

	t.Run("Generate Proposal with Multiple Operations - Proposer Role", func(t *testing.T) {
		// Create operation definitions and inputs
		defs := []cld_ops.Definition{
			ccipops.AddPackageIdStateObjectOp.Def(),
			ccipops.TransferOwnershipStateObjectOp.Def(),
			ccipops.AcceptOwnershipStateObjectOp.Def(),
		}

		inputs := []sui_ops.OpTxInput[any]{
			{
				Input: ccipops.AddPackageIdStateObjectInput{
					CCIPPackageId:         testCCIPPackageId,
					CCIPObjectRefObjectId: testObjectRefId,
					OwnerCapObjectId:      testOwnerCapId,
					PackageId:             testPackageId,
				},
				NoExecute: true,
			},
			{
				Input: ccipops.TransferOwnershipStateObjectInput{
					CCIPPackageId:         testCCIPPackageId,
					CCIPObjectRefObjectId: testObjectRefId,
					OwnerCapObjectId:      testOwnerCapId,
					To:                    testNewOwner,
				},
				NoExecute: true,
			},
			{
				Input: ccipops.AcceptOwnershipStateObjectInput{
					CCIPPackageId:         testCCIPPackageId,
					CCIPObjectRefObjectId: testObjectRefId,
				},
				NoExecute: true,
			},
		}

		proposalInput := Input{
			Defs:   defs,
			Inputs: inputs,

			MmcsPackageID:  testCCIPPackageId,
			McmsStateObjID: testObjectRefId,
			TimelockObjID:  testTimelockObjID,
			AccountObjID:   testAccountObjID,
			RegistryObjID:  testRegistryObjID,

			Role:  suisdk.TimelockRoleProposer,
			Delay: time.Hour * 24,

			ChainSelector: testChainSelector,
		}

		// Execute the operation
		result, err := cld_ops.ExecuteOperation(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
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
		assert.Equal(t, testTimelockObjID, proposal.TimelockAddresses[0], "timelock address should match")

		// Verify chain metadata
		require.Len(t, proposal.ChainMetadata, 1, "should have one chain metadata")
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

		inputs := []sui_ops.OpTxInput[any]{
			{
				Input: ccipops.RemovePackageIdStateObjectInput{
					CCIPPackageId:         testCCIPPackageId,
					CCIPObjectRefObjectId: testObjectRefId,
					OwnerCapObjectId:      testOwnerCapId,
					PackageId:             testPackageId,
				},
				NoExecute: true,
			},
		}

		proposalInput := Input{
			Defs:   defs,
			Inputs: inputs,

			MmcsPackageID:  testCCIPPackageId,
			McmsStateObjID: testObjectRefId,
			TimelockObjID:  testTimelockObjID,
			AccountObjID:   testAccountObjID,
			RegistryObjID:  testRegistryObjID,

			Role:  suisdk.TimelockRoleBypasser,
			Delay: 0, // No delay for bypasser

			ChainSelector: testChainSelector,
		}

		// Execute the operation
		result, err := cld_ops.ExecuteOperation(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		require.NoError(t, err, "should generate proposal successfully")

		// Verify the proposal structure
		proposal := result.Output
		assert.Equal(t, "v1", proposal.Version, "proposal version should be v1")
		assert.Contains(t, proposal.Description, "ccip.state_object.remove_package_id", "description should contain remove operation")

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

		inputs := []sui_ops.OpTxInput[any]{
			{
				Input: ccipops.AddPackageIdStateObjectInput{
					CCIPPackageId:         testCCIPPackageId,
					CCIPObjectRefObjectId: testObjectRefId,
					OwnerCapObjectId:      testOwnerCapId,
					PackageId:             testPackageId,
				},
				NoExecute: true,
			},
		}

		proposalInput := Input{
			Defs:   defs,
			Inputs: inputs,

			MmcsPackageID:  testCCIPPackageId,
			McmsStateObjID: testObjectRefId,
			TimelockObjID:  testTimelockObjID,
			AccountObjID:   testAccountObjID,
			RegistryObjID:  testRegistryObjID,

			Role:  suisdk.TimelockRole(100), // Invalid role
			Delay: time.Hour,

			ChainSelector: testChainSelector,
		}

		// Execute the operation - should fail
		_, err := cld_ops.ExecuteOperation(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		require.Error(t, err, "should fail with invalid role")
		assert.Contains(t, err.Error(), "unsupported role", "error should mention unsupported role")
	})

	t.Run("Generate Proposal with Mismatched Definitions and Inputs", func(t *testing.T) {
		// Create mismatched definitions and inputs (more defs than inputs)
		defs := []cld_ops.Definition{
			ccipops.AddPackageIdStateObjectOp.Def(),
			ccipops.RemovePackageIdStateObjectOp.Def(),
		}

		inputs := []sui_ops.OpTxInput[any]{
			{
				Input: ccipops.AddPackageIdStateObjectInput{
					CCIPPackageId:         testCCIPPackageId,
					CCIPObjectRefObjectId: testObjectRefId,
					OwnerCapObjectId:      testOwnerCapId,
					PackageId:             testPackageId,
				},
				NoExecute: true,
			},
			// Missing second input
		}

		proposalInput := Input{
			Defs:   defs,
			Inputs: inputs[:1], // Only one input for two definitions

			MmcsPackageID:  testCCIPPackageId,
			McmsStateObjID: testObjectRefId,
			TimelockObjID:  testTimelockObjID,
			AccountObjID:   testAccountObjID,
			RegistryObjID:  testRegistryObjID,

			Role:  suisdk.TimelockRoleProposer,
			Delay: time.Hour,

			ChainSelector: testChainSelector,
		}

		// Execute the operation - should fail due to index out of bounds
		_, err := cld_ops.ExecuteOperation(bundle, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		require.Error(t, err, "should fail with mismatched definitions and inputs")
	})
}
