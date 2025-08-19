// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

type IMcms interface {
	SetRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, clock bind.Object, role byte, root []byte, validUntil uint64, chainId *big.Int, multisigAddr string, preOpCount uint64, postOpCount uint64, overridePreviousRoot bool, metadataProof [][]byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error)
	Execute(ctx context.Context, opts *bind.CallOpts, state bind.Object, clock bind.Object, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) (*models.SuiTransactionBlockResponse, error)
	DispatchTimelockScheduleBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error)
	DispatchTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error)
	DispatchTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error)
	DispatchTimelockCancel(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error)
	DispatchTimelockUpdateMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error)
	DispatchTimelockBlockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error)
	DispatchTimelockUnblockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error)
	ExecuteDispatchToAccount(ctx context.Context, opts *bind.CallOpts, registry bind.Object, accountState bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteDispatchToDeployer(ctx context.Context, opts *bind.CallOpts, registry bind.Object, deployerState bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteTimelockScheduleBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteTimelockCancel(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteTimelockUpdateMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteTimelockBlockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteTimelockUnblockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteSetConfig(ctx context.Context, opts *bind.CallOpts, registry bind.Object, state bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetConfig(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, role byte, chainId *big.Int, signerAddresses [][]byte, signerGroups []byte, groupQuorums []byte, groupParents []byte, clearRoot bool) (*models.SuiTransactionBlockResponse, error)
	VerifyMerkleProof(ctx context.Context, opts *bind.CallOpts, proof [][]byte, root []byte, leaf []byte) (*models.SuiTransactionBlockResponse, error)
	ComputeEthMessageHash(ctx context.Context, opts *bind.CallOpts, root []byte, validUntil uint64) (*models.SuiTransactionBlockResponse, error)
	HashOpLeaf(ctx context.Context, opts *bind.CallOpts, domainSeparator []byte, op Op) (*models.SuiTransactionBlockResponse, error)
	SeenSignedHashes(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error)
	ExpiringRootAndOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error)
	RootMetadata(ctx context.Context, opts *bind.CallOpts, multisig Multisig) (*models.SuiTransactionBlockResponse, error)
	GetRootMetadata(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error)
	GetOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error)
	GetRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error)
	GetConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error)
	NumGroups(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	MaxNumSigners(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	BypasserRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	CancellerRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	ProposerRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	TimelockRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	IsValidRole(ctx context.Context, opts *bind.CallOpts, role byte) (*models.SuiTransactionBlockResponse, error)
	ZeroHash(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Role(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error)
	ChainId(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error)
	RootMetadataMultisig(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error)
	PreOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error)
	PostOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error)
	OverridePreviousRoot(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error)
	ConfigSigners(ctx context.Context, opts *bind.CallOpts, config Config) (*models.SuiTransactionBlockResponse, error)
	ConfigGroupQuorums(ctx context.Context, opts *bind.CallOpts, config Config) (*models.SuiTransactionBlockResponse, error)
	ConfigGroupParents(ctx context.Context, opts *bind.CallOpts, config Config) (*models.SuiTransactionBlockResponse, error)
	TimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) (*models.SuiTransactionBlockResponse, error)
	TimelockGetBlockedFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, index uint64) (*models.SuiTransactionBlockResponse, error)
	TimelockIsOperation(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error)
	TimelockIsOperationPending(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error)
	TimelockIsOperationReady(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error)
	TimelockIsOperationDone(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error)
	TimelockGetTimestamp(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error)
	TimelockMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (*models.SuiTransactionBlockResponse, error)
	TimelockGetBlockedFunctions(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (*models.SuiTransactionBlockResponse, error)
	TimelockGetBlockedFunctionsCount(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (*models.SuiTransactionBlockResponse, error)
	CreateCalls(ctx context.Context, opts *bind.CallOpts, targets []string, moduleNames []string, functionNames []string, datas [][]byte) (*models.SuiTransactionBlockResponse, error)
	HashOperationBatch(ctx context.Context, opts *bind.CallOpts, calls []Call, predecessor []byte, salt []byte) (*models.SuiTransactionBlockResponse, error)
	SignerView(ctx context.Context, opts *bind.CallOpts, signer Signer) (*models.SuiTransactionBlockResponse, error)
	FunctionName(ctx context.Context, opts *bind.CallOpts, function Function) (*models.SuiTransactionBlockResponse, error)
	ModuleName(ctx context.Context, opts *bind.CallOpts, function Function) (*models.SuiTransactionBlockResponse, error)
	Target(ctx context.Context, opts *bind.CallOpts, function Function) (*models.SuiTransactionBlockResponse, error)
	Data(ctx context.Context, opts *bind.CallOpts, call Call) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMcmsDevInspect
	Encoder() McmsEncoder
}

type IMcmsDevInspect interface {
	Execute(ctx context.Context, opts *bind.CallOpts, state bind.Object, clock bind.Object, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) (TimelockCallbackParams, error)
	DispatchTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) ([]bind.Object, error)
	DispatchTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelockCallbackParams TimelockCallbackParams) ([]bind.Object, error)
	ExecuteDispatchToDeployer(ctx context.Context, opts *bind.CallOpts, registry bind.Object, deployerState bind.Object, executingCallbackParams bind.Object) (bind.Object, error)
	ExecuteTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) ([]bind.Object, error)
	ExecuteTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, executingCallbackParams bind.Object) ([]bind.Object, error)
	VerifyMerkleProof(ctx context.Context, opts *bind.CallOpts, proof [][]byte, root []byte, leaf []byte) (bool, error)
	ComputeEthMessageHash(ctx context.Context, opts *bind.CallOpts, root []byte, validUntil uint64) ([]byte, error)
	HashOpLeaf(ctx context.Context, opts *bind.CallOpts, domainSeparator []byte, op Op) ([]byte, error)
	SeenSignedHashes(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (bind.Object, error)
	ExpiringRootAndOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) ([]any, error)
	RootMetadata(ctx context.Context, opts *bind.CallOpts, multisig Multisig) (RootMetadata, error)
	GetRootMetadata(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (RootMetadata, error)
	GetOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (uint64, error)
	GetRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) ([]any, error)
	GetConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (Config, error)
	NumGroups(ctx context.Context, opts *bind.CallOpts) (uint64, error)
	MaxNumSigners(ctx context.Context, opts *bind.CallOpts) (uint64, error)
	BypasserRole(ctx context.Context, opts *bind.CallOpts) (byte, error)
	CancellerRole(ctx context.Context, opts *bind.CallOpts) (byte, error)
	ProposerRole(ctx context.Context, opts *bind.CallOpts) (byte, error)
	TimelockRole(ctx context.Context, opts *bind.CallOpts) (byte, error)
	IsValidRole(ctx context.Context, opts *bind.CallOpts, role byte) (bool, error)
	ZeroHash(ctx context.Context, opts *bind.CallOpts) ([]byte, error)
	Role(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (byte, error)
	ChainId(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*big.Int, error)
	RootMetadataMultisig(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (string, error)
	PreOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (uint64, error)
	PostOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (uint64, error)
	OverridePreviousRoot(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (bool, error)
	ConfigSigners(ctx context.Context, opts *bind.CallOpts, config Config) ([]Signer, error)
	ConfigGroupQuorums(ctx context.Context, opts *bind.CallOpts, config Config) ([]byte, error)
	ConfigGroupParents(ctx context.Context, opts *bind.CallOpts, config Config) ([]byte, error)
	TimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) ([]bind.Object, error)
	TimelockGetBlockedFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, index uint64) (Function, error)
	TimelockIsOperation(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (bool, error)
	TimelockIsOperationPending(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (bool, error)
	TimelockIsOperationReady(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, id []byte) (bool, error)
	TimelockIsOperationDone(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (bool, error)
	TimelockGetTimestamp(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (uint64, error)
	TimelockMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (uint64, error)
	TimelockGetBlockedFunctions(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) ([]Function, error)
	TimelockGetBlockedFunctionsCount(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (uint64, error)
	CreateCalls(ctx context.Context, opts *bind.CallOpts, targets []string, moduleNames []string, functionNames []string, datas [][]byte) ([]Call, error)
	HashOperationBatch(ctx context.Context, opts *bind.CallOpts, calls []Call, predecessor []byte, salt []byte) ([]byte, error)
	SignerView(ctx context.Context, opts *bind.CallOpts, signer Signer) ([]any, error)
	FunctionName(ctx context.Context, opts *bind.CallOpts, function Function) (string, error)
	ModuleName(ctx context.Context, opts *bind.CallOpts, function Function) (string, error)
	Target(ctx context.Context, opts *bind.CallOpts, function Function) (string, error)
	Data(ctx context.Context, opts *bind.CallOpts, call Call) ([]byte, error)
}

type McmsEncoder interface {
	SetRoot(state bind.Object, clock bind.Object, role byte, root []byte, validUntil uint64, chainId *big.Int, multisigAddr string, preOpCount uint64, postOpCount uint64, overridePreviousRoot bool, metadataProof [][]byte, signatures [][]byte) (*bind.EncodedCall, error)
	SetRootWithArgs(args ...any) (*bind.EncodedCall, error)
	Execute(state bind.Object, clock bind.Object, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) (*bind.EncodedCall, error)
	ExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	DispatchTimelockScheduleBatch(timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error)
	DispatchTimelockScheduleBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	DispatchTimelockExecuteBatch(timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error)
	DispatchTimelockExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	DispatchTimelockBypasserExecuteBatch(timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error)
	DispatchTimelockBypasserExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	DispatchTimelockCancel(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error)
	DispatchTimelockCancelWithArgs(args ...any) (*bind.EncodedCall, error)
	DispatchTimelockUpdateMinDelay(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error)
	DispatchTimelockUpdateMinDelayWithArgs(args ...any) (*bind.EncodedCall, error)
	DispatchTimelockBlockFunction(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error)
	DispatchTimelockBlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	DispatchTimelockUnblockFunction(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error)
	DispatchTimelockUnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteDispatchToAccount(registry bind.Object, accountState bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteDispatchToAccountWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteDispatchToDeployer(registry bind.Object, deployerState bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteDispatchToDeployerWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteTimelockScheduleBatch(timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteTimelockScheduleBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteTimelockExecuteBatch(timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteTimelockExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteTimelockBypasserExecuteBatch(executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteTimelockBypasserExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteTimelockCancel(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteTimelockCancelWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteTimelockUpdateMinDelay(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteTimelockUpdateMinDelayWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteTimelockBlockFunction(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteTimelockBlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteTimelockUnblockFunction(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteTimelockUnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteSetConfig(registry bind.Object, state bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error)
	ExecuteSetConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	SetConfig(param bind.Object, state bind.Object, role byte, chainId *big.Int, signerAddresses [][]byte, signerGroups []byte, groupQuorums []byte, groupParents []byte, clearRoot bool) (*bind.EncodedCall, error)
	SetConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	VerifyMerkleProof(proof [][]byte, root []byte, leaf []byte) (*bind.EncodedCall, error)
	VerifyMerkleProofWithArgs(args ...any) (*bind.EncodedCall, error)
	ComputeEthMessageHash(root []byte, validUntil uint64) (*bind.EncodedCall, error)
	ComputeEthMessageHashWithArgs(args ...any) (*bind.EncodedCall, error)
	HashOpLeaf(domainSeparator []byte, op Op) (*bind.EncodedCall, error)
	HashOpLeafWithArgs(args ...any) (*bind.EncodedCall, error)
	SeenSignedHashes(state bind.Object, role byte) (*bind.EncodedCall, error)
	SeenSignedHashesWithArgs(args ...any) (*bind.EncodedCall, error)
	ExpiringRootAndOpCount(state bind.Object, role byte) (*bind.EncodedCall, error)
	ExpiringRootAndOpCountWithArgs(args ...any) (*bind.EncodedCall, error)
	RootMetadata(multisig Multisig) (*bind.EncodedCall, error)
	RootMetadataWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRootMetadata(state bind.Object, role byte) (*bind.EncodedCall, error)
	GetRootMetadataWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOpCount(state bind.Object, role byte) (*bind.EncodedCall, error)
	GetOpCountWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRoot(state bind.Object, role byte) (*bind.EncodedCall, error)
	GetRootWithArgs(args ...any) (*bind.EncodedCall, error)
	GetConfig(state bind.Object, role byte) (*bind.EncodedCall, error)
	GetConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	NumGroups() (*bind.EncodedCall, error)
	NumGroupsWithArgs(args ...any) (*bind.EncodedCall, error)
	MaxNumSigners() (*bind.EncodedCall, error)
	MaxNumSignersWithArgs(args ...any) (*bind.EncodedCall, error)
	BypasserRole() (*bind.EncodedCall, error)
	BypasserRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	CancellerRole() (*bind.EncodedCall, error)
	CancellerRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	ProposerRole() (*bind.EncodedCall, error)
	ProposerRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockRole() (*bind.EncodedCall, error)
	TimelockRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	IsValidRole(role byte) (*bind.EncodedCall, error)
	IsValidRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	ZeroHash() (*bind.EncodedCall, error)
	ZeroHashWithArgs(args ...any) (*bind.EncodedCall, error)
	Role(rootMetadata RootMetadata) (*bind.EncodedCall, error)
	RoleWithArgs(args ...any) (*bind.EncodedCall, error)
	ChainId(rootMetadata RootMetadata) (*bind.EncodedCall, error)
	ChainIdWithArgs(args ...any) (*bind.EncodedCall, error)
	RootMetadataMultisig(rootMetadata RootMetadata) (*bind.EncodedCall, error)
	RootMetadataMultisigWithArgs(args ...any) (*bind.EncodedCall, error)
	PreOpCount(rootMetadata RootMetadata) (*bind.EncodedCall, error)
	PreOpCountWithArgs(args ...any) (*bind.EncodedCall, error)
	PostOpCount(rootMetadata RootMetadata) (*bind.EncodedCall, error)
	PostOpCountWithArgs(args ...any) (*bind.EncodedCall, error)
	OverridePreviousRoot(rootMetadata RootMetadata) (*bind.EncodedCall, error)
	OverridePreviousRootWithArgs(args ...any) (*bind.EncodedCall, error)
	ConfigSigners(config Config) (*bind.EncodedCall, error)
	ConfigSignersWithArgs(args ...any) (*bind.EncodedCall, error)
	ConfigGroupQuorums(config Config) (*bind.EncodedCall, error)
	ConfigGroupQuorumsWithArgs(args ...any) (*bind.EncodedCall, error)
	ConfigGroupParents(config Config) (*bind.EncodedCall, error)
	ConfigGroupParentsWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockExecuteBatch(timelock bind.Object, clock bind.Object, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) (*bind.EncodedCall, error)
	TimelockExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockGetBlockedFunction(timelock bind.Object, index uint64) (*bind.EncodedCall, error)
	TimelockGetBlockedFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockIsOperation(timelock bind.Object, id []byte) (*bind.EncodedCall, error)
	TimelockIsOperationWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockIsOperationPending(timelock bind.Object, id []byte) (*bind.EncodedCall, error)
	TimelockIsOperationPendingWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockIsOperationReady(timelock bind.Object, clock bind.Object, id []byte) (*bind.EncodedCall, error)
	TimelockIsOperationReadyWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockIsOperationDone(timelock bind.Object, id []byte) (*bind.EncodedCall, error)
	TimelockIsOperationDoneWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockGetTimestamp(timelock bind.Object, id []byte) (*bind.EncodedCall, error)
	TimelockGetTimestampWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockMinDelay(timelock bind.Object) (*bind.EncodedCall, error)
	TimelockMinDelayWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockGetBlockedFunctions(timelock bind.Object) (*bind.EncodedCall, error)
	TimelockGetBlockedFunctionsWithArgs(args ...any) (*bind.EncodedCall, error)
	TimelockGetBlockedFunctionsCount(timelock bind.Object) (*bind.EncodedCall, error)
	TimelockGetBlockedFunctionsCountWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateCalls(targets []string, moduleNames []string, functionNames []string, datas [][]byte) (*bind.EncodedCall, error)
	CreateCallsWithArgs(args ...any) (*bind.EncodedCall, error)
	HashOperationBatch(calls []Call, predecessor []byte, salt []byte) (*bind.EncodedCall, error)
	HashOperationBatchWithArgs(args ...any) (*bind.EncodedCall, error)
	SignerView(signer Signer) (*bind.EncodedCall, error)
	SignerViewWithArgs(args ...any) (*bind.EncodedCall, error)
	FunctionName(function Function) (*bind.EncodedCall, error)
	FunctionNameWithArgs(args ...any) (*bind.EncodedCall, error)
	ModuleName(function Function) (*bind.EncodedCall, error)
	ModuleNameWithArgs(args ...any) (*bind.EncodedCall, error)
	Target(function Function) (*bind.EncodedCall, error)
	TargetWithArgs(args ...any) (*bind.EncodedCall, error)
	Data(call Call) (*bind.EncodedCall, error)
	DataWithArgs(args ...any) (*bind.EncodedCall, error)
}

type McmsContract struct {
	*bind.BoundContract
	mcmsEncoder
	devInspect *McmsDevInspect
}

type McmsDevInspect struct {
	contract *McmsContract
}

var _ IMcms = (*McmsContract)(nil)
var _ IMcmsDevInspect = (*McmsDevInspect)(nil)

func NewMcms(packageID string, client sui.ISuiAPI) (*McmsContract, error) {
	contract, err := bind.NewBoundContract(packageID, "mcms", "mcms", client)
	if err != nil {
		return nil, err
	}

	c := &McmsContract{
		BoundContract: contract,
		mcmsEncoder:   mcmsEncoder{BoundContract: contract},
	}
	c.devInspect = &McmsDevInspect{contract: c}
	return c, nil
}

func (c *McmsContract) Encoder() McmsEncoder {
	return c.mcmsEncoder
}

func (c *McmsContract) DevInspect() IMcmsDevInspect {
	return c.devInspect
}

type MultisigState struct {
	Id        string   `move:"sui::object::UID"`
	Bypasser  Multisig `move:"Multisig"`
	Canceller Multisig `move:"Multisig"`
	Proposer  Multisig `move:"Multisig"`
}

type Multisig struct {
	Role                   byte                   `move:"u8"`
	Signers                bind.Object            `move:"VecMap<vector<u8>, Signer>"`
	Config                 Config                 `move:"Config"`
	SeenSignedHashes       bind.Object            `move:"VecMap<vector<u8>, bool>"`
	ExpiringRootAndOpCount ExpiringRootAndOpCount `move:"ExpiringRootAndOpCount"`
	RootMetadata           RootMetadata           `move:"RootMetadata"`
}

type Signer struct {
	Addr  []byte `move:"vector<u8>"`
	Index byte   `move:"u8"`
	Group byte   `move:"u8"`
}

type Config struct {
	Signers      []Signer `move:"vector<Signer>"`
	GroupQuorums []byte   `move:"vector<u8>"`
	GroupParents []byte   `move:"vector<u8>"`
}

type ExpiringRootAndOpCount struct {
	Root       []byte `move:"vector<u8>"`
	ValidUntil uint64 `move:"u64"`
	OpCount    uint64 `move:"u64"`
}

type Op struct {
	Role         byte     `move:"u8"`
	ChainId      *big.Int `move:"u256"`
	Multisig     string   `move:"address"`
	Nonce        uint64   `move:"u64"`
	To           string   `move:"address"`
	ModuleName   string   `move:"0x1::string::String"`
	FunctionName string   `move:"0x1::string::String"`
	Data         []byte   `move:"vector<u8>"`
}

type RootMetadata struct {
	Role                 byte     `move:"u8"`
	ChainId              *big.Int `move:"u256"`
	Multisig             string   `move:"address"`
	PreOpCount           uint64   `move:"u64"`
	PostOpCount          uint64   `move:"u64"`
	OverridePreviousRoot bool     `move:"bool"`
}

type TimelockCallbackParams struct {
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Data         []byte `move:"vector<u8>"`
	Role         byte   `move:"u8"`
}

type MultisigStateInitialized struct {
	Bypasser  byte `move:"u8"`
	Canceller byte `move:"u8"`
	Proposer  byte `move:"u8"`
}

type ConfigSet struct {
	Role          byte   `move:"u8"`
	Config        Config `move:"Config"`
	IsRootCleared bool   `move:"bool"`
}

type NewRoot struct {
	Role       byte         `move:"u8"`
	Root       []byte       `move:"vector<u8>"`
	ValidUntil uint64       `move:"u64"`
	Metadata   RootMetadata `move:"RootMetadata"`
}

type OpExecuted struct {
	Role         byte     `move:"u8"`
	ChainId      *big.Int `move:"u256"`
	Multisig     string   `move:"address"`
	Nonce        uint64   `move:"u64"`
	To           string   `move:"address"`
	ModuleName   string   `move:"0x1::string::String"`
	FunctionName string   `move:"0x1::string::String"`
	Data         []byte   `move:"vector<u8>"`
}

type MCMS struct {
}

type McmsCallback struct {
}

type Timelock struct {
	Id               string      `move:"sui::object::UID"`
	MinDelay         uint64      `move:"u64"`
	Timestamps       bind.Object `move:"Table<vector<u8>, u64>"`
	BlockedFunctions bind.Object `move:"VecSet<Function>"`
}

type Call struct {
	Function Function `move:"Function"`
	Data     []byte   `move:"vector<u8>"`
}

type Function struct {
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
}

type TimelockInitialized struct {
	MinDelay uint64 `move:"u64"`
}

type BypasserCallInitiated struct {
	Index        uint64 `move:"u64"`
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Data         []byte `move:"vector<u8>"`
}

type Cancelled struct {
	Id []byte `move:"vector<u8>"`
}

type CallScheduled struct {
	Id           []byte `move:"vector<u8>"`
	Index        uint64 `move:"u64"`
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Data         []byte `move:"vector<u8>"`
	Predecessor  []byte `move:"vector<u8>"`
	Salt         []byte `move:"vector<u8>"`
	Delay        uint64 `move:"u64"`
}

type CallInitiated struct {
	Id           []byte `move:"vector<u8>"`
	Index        uint64 `move:"u64"`
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Data         []byte `move:"vector<u8>"`
}

type UpdateMinDelay struct {
	OldMinDelay uint64 `move:"u64"`
	NewMinDelay uint64 `move:"u64"`
}

type FunctionBlocked struct {
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
}

type FunctionUnblocked struct {
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
}

type bcsMultisigState struct {
	Id        string
	Bypasser  bcsMultisig
	Canceller bcsMultisig
	Proposer  bcsMultisig
}

func convertMultisigStateFromBCS(bcs bcsMultisigState) (MultisigState, error) {
	BypasserField, err := convertMultisigFromBCS(bcs.Bypasser)
	if err != nil {
		return MultisigState{}, fmt.Errorf("failed to convert nested struct Bypasser: %w", err)
	}
	CancellerField, err := convertMultisigFromBCS(bcs.Canceller)
	if err != nil {
		return MultisigState{}, fmt.Errorf("failed to convert nested struct Canceller: %w", err)
	}
	ProposerField, err := convertMultisigFromBCS(bcs.Proposer)
	if err != nil {
		return MultisigState{}, fmt.Errorf("failed to convert nested struct Proposer: %w", err)
	}

	return MultisigState{
		Id:        bcs.Id,
		Bypasser:  BypasserField,
		Canceller: CancellerField,
		Proposer:  ProposerField,
	}, nil
}

type bcsMultisig struct {
	Role                   byte
	Signers                bind.Object
	Config                 Config
	SeenSignedHashes       bind.Object
	ExpiringRootAndOpCount ExpiringRootAndOpCount
	RootMetadata           bcsRootMetadata
}

func convertMultisigFromBCS(bcs bcsMultisig) (Multisig, error) {
	RootMetadataField, err := convertRootMetadataFromBCS(bcs.RootMetadata)
	if err != nil {
		return Multisig{}, fmt.Errorf("failed to convert nested struct RootMetadata: %w", err)
	}

	return Multisig{
		Role:                   bcs.Role,
		Signers:                bcs.Signers,
		Config:                 bcs.Config,
		SeenSignedHashes:       bcs.SeenSignedHashes,
		ExpiringRootAndOpCount: bcs.ExpiringRootAndOpCount,
		RootMetadata:           RootMetadataField,
	}, nil
}

type bcsOp struct {
	Role         byte
	ChainId      [32]byte
	Multisig     [32]byte
	Nonce        uint64
	To           [32]byte
	ModuleName   string
	FunctionName string
	Data         []byte
}

func convertOpFromBCS(bcs bcsOp) (Op, error) {
	ChainIdField, err := bind.DecodeU256Value(bcs.ChainId)
	if err != nil {
		return Op{}, fmt.Errorf("failed to decode u256 field ChainId: %w", err)
	}

	return Op{
		Role:         bcs.Role,
		ChainId:      ChainIdField,
		Multisig:     fmt.Sprintf("0x%x", bcs.Multisig),
		Nonce:        bcs.Nonce,
		To:           fmt.Sprintf("0x%x", bcs.To),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
		Data:         bcs.Data,
	}, nil
}

type bcsRootMetadata struct {
	Role                 byte
	ChainId              [32]byte
	Multisig             [32]byte
	PreOpCount           uint64
	PostOpCount          uint64
	OverridePreviousRoot bool
}

func convertRootMetadataFromBCS(bcs bcsRootMetadata) (RootMetadata, error) {
	ChainIdField, err := bind.DecodeU256Value(bcs.ChainId)
	if err != nil {
		return RootMetadata{}, fmt.Errorf("failed to decode u256 field ChainId: %w", err)
	}

	return RootMetadata{
		Role:                 bcs.Role,
		ChainId:              ChainIdField,
		Multisig:             fmt.Sprintf("0x%x", bcs.Multisig),
		PreOpCount:           bcs.PreOpCount,
		PostOpCount:          bcs.PostOpCount,
		OverridePreviousRoot: bcs.OverridePreviousRoot,
	}, nil
}

type bcsNewRoot struct {
	Role       byte
	Root       []byte
	ValidUntil uint64
	Metadata   bcsRootMetadata
}

func convertNewRootFromBCS(bcs bcsNewRoot) (NewRoot, error) {
	MetadataField, err := convertRootMetadataFromBCS(bcs.Metadata)
	if err != nil {
		return NewRoot{}, fmt.Errorf("failed to convert nested struct Metadata: %w", err)
	}

	return NewRoot{
		Role:       bcs.Role,
		Root:       bcs.Root,
		ValidUntil: bcs.ValidUntil,
		Metadata:   MetadataField,
	}, nil
}

type bcsOpExecuted struct {
	Role         byte
	ChainId      [32]byte
	Multisig     [32]byte
	Nonce        uint64
	To           [32]byte
	ModuleName   string
	FunctionName string
	Data         []byte
}

func convertOpExecutedFromBCS(bcs bcsOpExecuted) (OpExecuted, error) {
	ChainIdField, err := bind.DecodeU256Value(bcs.ChainId)
	if err != nil {
		return OpExecuted{}, fmt.Errorf("failed to decode u256 field ChainId: %w", err)
	}

	return OpExecuted{
		Role:         bcs.Role,
		ChainId:      ChainIdField,
		Multisig:     fmt.Sprintf("0x%x", bcs.Multisig),
		Nonce:        bcs.Nonce,
		To:           fmt.Sprintf("0x%x", bcs.To),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
		Data:         bcs.Data,
	}, nil
}

type bcsCall struct {
	Function bcsFunction
	Data     []byte
}

func convertCallFromBCS(bcs bcsCall) (Call, error) {
	FunctionField, err := convertFunctionFromBCS(bcs.Function)
	if err != nil {
		return Call{}, fmt.Errorf("failed to convert nested struct Function: %w", err)
	}

	return Call{
		Function: FunctionField,
		Data:     bcs.Data,
	}, nil
}

type bcsFunction struct {
	Target       [32]byte
	ModuleName   string
	FunctionName string
}

func convertFunctionFromBCS(bcs bcsFunction) (Function, error) {

	return Function{
		Target:       fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
	}, nil
}

type bcsBypasserCallInitiated struct {
	Index        uint64
	Target       [32]byte
	ModuleName   string
	FunctionName string
	Data         []byte
}

func convertBypasserCallInitiatedFromBCS(bcs bcsBypasserCallInitiated) (BypasserCallInitiated, error) {

	return BypasserCallInitiated{
		Index:        bcs.Index,
		Target:       fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
		Data:         bcs.Data,
	}, nil
}

type bcsCallScheduled struct {
	Id           []byte
	Index        uint64
	Target       [32]byte
	ModuleName   string
	FunctionName string
	Data         []byte
	Predecessor  []byte
	Salt         []byte
	Delay        uint64
}

func convertCallScheduledFromBCS(bcs bcsCallScheduled) (CallScheduled, error) {

	return CallScheduled{
		Id:           bcs.Id,
		Index:        bcs.Index,
		Target:       fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
		Data:         bcs.Data,
		Predecessor:  bcs.Predecessor,
		Salt:         bcs.Salt,
		Delay:        bcs.Delay,
	}, nil
}

type bcsCallInitiated struct {
	Id           []byte
	Index        uint64
	Target       [32]byte
	ModuleName   string
	FunctionName string
	Data         []byte
}

func convertCallInitiatedFromBCS(bcs bcsCallInitiated) (CallInitiated, error) {

	return CallInitiated{
		Id:           bcs.Id,
		Index:        bcs.Index,
		Target:       fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
		Data:         bcs.Data,
	}, nil
}

type bcsFunctionBlocked struct {
	Target       [32]byte
	ModuleName   string
	FunctionName string
}

func convertFunctionBlockedFromBCS(bcs bcsFunctionBlocked) (FunctionBlocked, error) {

	return FunctionBlocked{
		Target:       fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
	}, nil
}

type bcsFunctionUnblocked struct {
	Target       [32]byte
	ModuleName   string
	FunctionName string
}

func convertFunctionUnblockedFromBCS(bcs bcsFunctionUnblocked) (FunctionUnblocked, error) {

	return FunctionUnblocked{
		Target:       fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("mcms::mcms::MultisigState", func(data []byte) (interface{}, error) {
		var temp bcsMultisigState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertMultisigStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Multisig", func(data []byte) (interface{}, error) {
		var temp bcsMultisig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertMultisigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Signer", func(data []byte) (interface{}, error) {
		var result Signer
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Config", func(data []byte) (interface{}, error) {
		var result Config
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::ExpiringRootAndOpCount", func(data []byte) (interface{}, error) {
		var result ExpiringRootAndOpCount
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Op", func(data []byte) (interface{}, error) {
		var temp bcsOp
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOpFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::RootMetadata", func(data []byte) (interface{}, error) {
		var temp bcsRootMetadata
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertRootMetadataFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::TimelockCallbackParams", func(data []byte) (interface{}, error) {
		var result TimelockCallbackParams
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::MultisigStateInitialized", func(data []byte) (interface{}, error) {
		var result MultisigStateInitialized
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::ConfigSet", func(data []byte) (interface{}, error) {
		var result ConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::NewRoot", func(data []byte) (interface{}, error) {
		var temp bcsNewRoot
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertNewRootFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::OpExecuted", func(data []byte) (interface{}, error) {
		var temp bcsOpExecuted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOpExecutedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::MCMS", func(data []byte) (interface{}, error) {
		var result MCMS
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Timelock", func(data []byte) (interface{}, error) {
		var result Timelock
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Call", func(data []byte) (interface{}, error) {
		var temp bcsCall
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCallFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Function", func(data []byte) (interface{}, error) {
		var temp bcsFunction
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFunctionFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::TimelockInitialized", func(data []byte) (interface{}, error) {
		var result TimelockInitialized
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::BypasserCallInitiated", func(data []byte) (interface{}, error) {
		var temp bcsBypasserCallInitiated
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertBypasserCallInitiatedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::Cancelled", func(data []byte) (interface{}, error) {
		var result Cancelled
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::CallScheduled", func(data []byte) (interface{}, error) {
		var temp bcsCallScheduled
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCallScheduledFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::CallInitiated", func(data []byte) (interface{}, error) {
		var temp bcsCallInitiated
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCallInitiatedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::UpdateMinDelay", func(data []byte) (interface{}, error) {
		var result UpdateMinDelay
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::FunctionBlocked", func(data []byte) (interface{}, error) {
		var temp bcsFunctionBlocked
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFunctionBlockedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms::FunctionUnblocked", func(data []byte) (interface{}, error) {
		var temp bcsFunctionUnblocked
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFunctionUnblockedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// SetRoot executes the set_root Move function.
func (c *McmsContract) SetRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, clock bind.Object, role byte, root []byte, validUntil uint64, chainId *big.Int, multisigAddr string, preOpCount uint64, postOpCount uint64, overridePreviousRoot bool, metadataProof [][]byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.SetRoot(state, clock, role, root, validUntil, chainId, multisigAddr, preOpCount, postOpCount, overridePreviousRoot, metadataProof, signatures)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Execute executes the execute Move function.
func (c *McmsContract) Execute(ctx context.Context, opts *bind.CallOpts, state bind.Object, clock bind.Object, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.Execute(state, clock, role, chainId, multisigAddr, nonce, to, moduleName, functionName, data, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DispatchTimelockScheduleBatch executes the dispatch_timelock_schedule_batch Move function.
func (c *McmsContract) DispatchTimelockScheduleBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.DispatchTimelockScheduleBatch(timelock, clock, timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DispatchTimelockExecuteBatch executes the dispatch_timelock_execute_batch Move function.
func (c *McmsContract) DispatchTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.DispatchTimelockExecuteBatch(timelock, clock, timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DispatchTimelockBypasserExecuteBatch executes the dispatch_timelock_bypasser_execute_batch Move function.
func (c *McmsContract) DispatchTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.DispatchTimelockBypasserExecuteBatch(timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DispatchTimelockCancel executes the dispatch_timelock_cancel Move function.
func (c *McmsContract) DispatchTimelockCancel(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.DispatchTimelockCancel(timelock, timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DispatchTimelockUpdateMinDelay executes the dispatch_timelock_update_min_delay Move function.
func (c *McmsContract) DispatchTimelockUpdateMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.DispatchTimelockUpdateMinDelay(timelock, timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DispatchTimelockBlockFunction executes the dispatch_timelock_block_function Move function.
func (c *McmsContract) DispatchTimelockBlockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.DispatchTimelockBlockFunction(timelock, timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DispatchTimelockUnblockFunction executes the dispatch_timelock_unblock_function Move function.
func (c *McmsContract) DispatchTimelockUnblockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.DispatchTimelockUnblockFunction(timelock, timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteDispatchToAccount executes the execute_dispatch_to_account Move function.
func (c *McmsContract) ExecuteDispatchToAccount(ctx context.Context, opts *bind.CallOpts, registry bind.Object, accountState bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteDispatchToAccount(registry, accountState, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteDispatchToDeployer executes the execute_dispatch_to_deployer Move function.
func (c *McmsContract) ExecuteDispatchToDeployer(ctx context.Context, opts *bind.CallOpts, registry bind.Object, deployerState bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteDispatchToDeployer(registry, deployerState, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteTimelockScheduleBatch executes the execute_timelock_schedule_batch Move function.
func (c *McmsContract) ExecuteTimelockScheduleBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteTimelockScheduleBatch(timelock, clock, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteTimelockExecuteBatch executes the execute_timelock_execute_batch Move function.
func (c *McmsContract) ExecuteTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteTimelockExecuteBatch(timelock, clock, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteTimelockBypasserExecuteBatch executes the execute_timelock_bypasser_execute_batch Move function.
func (c *McmsContract) ExecuteTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteTimelockBypasserExecuteBatch(executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteTimelockCancel executes the execute_timelock_cancel Move function.
func (c *McmsContract) ExecuteTimelockCancel(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteTimelockCancel(timelock, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteTimelockUpdateMinDelay executes the execute_timelock_update_min_delay Move function.
func (c *McmsContract) ExecuteTimelockUpdateMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteTimelockUpdateMinDelay(timelock, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteTimelockBlockFunction executes the execute_timelock_block_function Move function.
func (c *McmsContract) ExecuteTimelockBlockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteTimelockBlockFunction(timelock, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteTimelockUnblockFunction executes the execute_timelock_unblock_function Move function.
func (c *McmsContract) ExecuteTimelockUnblockFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteTimelockUnblockFunction(timelock, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteSetConfig executes the execute_set_config Move function.
func (c *McmsContract) ExecuteSetConfig(ctx context.Context, opts *bind.CallOpts, registry bind.Object, state bind.Object, executingCallbackParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExecuteSetConfig(registry, state, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetConfig executes the set_config Move function.
func (c *McmsContract) SetConfig(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, role byte, chainId *big.Int, signerAddresses [][]byte, signerGroups []byte, groupQuorums []byte, groupParents []byte, clearRoot bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.SetConfig(param, state, role, chainId, signerAddresses, signerGroups, groupQuorums, groupParents, clearRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// VerifyMerkleProof executes the verify_merkle_proof Move function.
func (c *McmsContract) VerifyMerkleProof(ctx context.Context, opts *bind.CallOpts, proof [][]byte, root []byte, leaf []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.VerifyMerkleProof(proof, root, leaf)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ComputeEthMessageHash executes the compute_eth_message_hash Move function.
func (c *McmsContract) ComputeEthMessageHash(ctx context.Context, opts *bind.CallOpts, root []byte, validUntil uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ComputeEthMessageHash(root, validUntil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HashOpLeaf executes the hash_op_leaf Move function.
func (c *McmsContract) HashOpLeaf(ctx context.Context, opts *bind.CallOpts, domainSeparator []byte, op Op) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.HashOpLeaf(domainSeparator, op)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SeenSignedHashes executes the seen_signed_hashes Move function.
func (c *McmsContract) SeenSignedHashes(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.SeenSignedHashes(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExpiringRootAndOpCount executes the expiring_root_and_op_count Move function.
func (c *McmsContract) ExpiringRootAndOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ExpiringRootAndOpCount(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RootMetadata executes the root_metadata Move function.
func (c *McmsContract) RootMetadata(ctx context.Context, opts *bind.CallOpts, multisig Multisig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.RootMetadata(multisig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRootMetadata executes the get_root_metadata Move function.
func (c *McmsContract) GetRootMetadata(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.GetRootMetadata(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOpCount executes the get_op_count Move function.
func (c *McmsContract) GetOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.GetOpCount(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRoot executes the get_root Move function.
func (c *McmsContract) GetRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.GetRoot(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetConfig executes the get_config Move function.
func (c *McmsContract) GetConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.GetConfig(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// NumGroups executes the num_groups Move function.
func (c *McmsContract) NumGroups(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.NumGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// MaxNumSigners executes the max_num_signers Move function.
func (c *McmsContract) MaxNumSigners(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.MaxNumSigners()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// BypasserRole executes the bypasser_role Move function.
func (c *McmsContract) BypasserRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.BypasserRole()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CancellerRole executes the canceller_role Move function.
func (c *McmsContract) CancellerRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.CancellerRole()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ProposerRole executes the proposer_role Move function.
func (c *McmsContract) ProposerRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ProposerRole()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockRole executes the timelock_role Move function.
func (c *McmsContract) TimelockRole(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockRole()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsValidRole executes the is_valid_role Move function.
func (c *McmsContract) IsValidRole(ctx context.Context, opts *bind.CallOpts, role byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.IsValidRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ZeroHash executes the zero_hash Move function.
func (c *McmsContract) ZeroHash(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ZeroHash()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Role executes the role Move function.
func (c *McmsContract) Role(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.Role(rootMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ChainId executes the chain_id Move function.
func (c *McmsContract) ChainId(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ChainId(rootMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RootMetadataMultisig executes the root_metadata_multisig Move function.
func (c *McmsContract) RootMetadataMultisig(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.RootMetadataMultisig(rootMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PreOpCount executes the pre_op_count Move function.
func (c *McmsContract) PreOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.PreOpCount(rootMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PostOpCount executes the post_op_count Move function.
func (c *McmsContract) PostOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.PostOpCount(rootMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// OverridePreviousRoot executes the override_previous_root Move function.
func (c *McmsContract) OverridePreviousRoot(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.OverridePreviousRoot(rootMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigSigners executes the config_signers Move function.
func (c *McmsContract) ConfigSigners(ctx context.Context, opts *bind.CallOpts, config Config) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ConfigSigners(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigGroupQuorums executes the config_group_quorums Move function.
func (c *McmsContract) ConfigGroupQuorums(ctx context.Context, opts *bind.CallOpts, config Config) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ConfigGroupQuorums(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigGroupParents executes the config_group_parents Move function.
func (c *McmsContract) ConfigGroupParents(ctx context.Context, opts *bind.CallOpts, config Config) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ConfigGroupParents(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockExecuteBatch executes the timelock_execute_batch Move function.
func (c *McmsContract) TimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockExecuteBatch(timelock, clock, targets, moduleNames, functionNames, datas, predecessor, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockGetBlockedFunction executes the timelock_get_blocked_function Move function.
func (c *McmsContract) TimelockGetBlockedFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, index uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockGetBlockedFunction(timelock, index)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockIsOperation executes the timelock_is_operation Move function.
func (c *McmsContract) TimelockIsOperation(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockIsOperation(timelock, id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockIsOperationPending executes the timelock_is_operation_pending Move function.
func (c *McmsContract) TimelockIsOperationPending(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockIsOperationPending(timelock, id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockIsOperationReady executes the timelock_is_operation_ready Move function.
func (c *McmsContract) TimelockIsOperationReady(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockIsOperationReady(timelock, clock, id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockIsOperationDone executes the timelock_is_operation_done Move function.
func (c *McmsContract) TimelockIsOperationDone(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockIsOperationDone(timelock, id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockGetTimestamp executes the timelock_get_timestamp Move function.
func (c *McmsContract) TimelockGetTimestamp(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockGetTimestamp(timelock, id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockMinDelay executes the timelock_min_delay Move function.
func (c *McmsContract) TimelockMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockMinDelay(timelock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockGetBlockedFunctions executes the timelock_get_blocked_functions Move function.
func (c *McmsContract) TimelockGetBlockedFunctions(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockGetBlockedFunctions(timelock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TimelockGetBlockedFunctionsCount executes the timelock_get_blocked_functions_count Move function.
func (c *McmsContract) TimelockGetBlockedFunctionsCount(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.TimelockGetBlockedFunctionsCount(timelock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateCalls executes the create_calls Move function.
func (c *McmsContract) CreateCalls(ctx context.Context, opts *bind.CallOpts, targets []string, moduleNames []string, functionNames []string, datas [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.CreateCalls(targets, moduleNames, functionNames, datas)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HashOperationBatch executes the hash_operation_batch Move function.
func (c *McmsContract) HashOperationBatch(ctx context.Context, opts *bind.CallOpts, calls []Call, predecessor []byte, salt []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.HashOperationBatch(calls, predecessor, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SignerView executes the signer_view Move function.
func (c *McmsContract) SignerView(ctx context.Context, opts *bind.CallOpts, signer Signer) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.SignerView(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// FunctionName executes the function_name Move function.
func (c *McmsContract) FunctionName(ctx context.Context, opts *bind.CallOpts, function Function) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.FunctionName(function)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ModuleName executes the module_name Move function.
func (c *McmsContract) ModuleName(ctx context.Context, opts *bind.CallOpts, function Function) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.ModuleName(function)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Target executes the target Move function.
func (c *McmsContract) Target(ctx context.Context, opts *bind.CallOpts, function Function) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.Target(function)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Data executes the data Move function.
func (c *McmsContract) Data(ctx context.Context, opts *bind.CallOpts, call Call) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsEncoder.Data(call)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Execute executes the execute Move function using DevInspect to get return values.
//
// Returns: TimelockCallbackParams
func (d *McmsDevInspect) Execute(ctx context.Context, opts *bind.CallOpts, state bind.Object, clock bind.Object, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) (TimelockCallbackParams, error) {
	encoded, err := d.contract.mcmsEncoder.Execute(state, clock, role, chainId, multisigAddr, nonce, to, moduleName, functionName, data, proof)
	if err != nil {
		return TimelockCallbackParams{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TimelockCallbackParams{}, err
	}
	if len(results) == 0 {
		return TimelockCallbackParams{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TimelockCallbackParams)
	if !ok {
		return TimelockCallbackParams{}, fmt.Errorf("unexpected return type: expected TimelockCallbackParams, got %T", results[0])
	}
	return result, nil
}

// DispatchTimelockExecuteBatch executes the dispatch_timelock_execute_batch Move function using DevInspect to get return values.
//
// Returns: vector<ExecutingCallbackParams>
func (d *McmsDevInspect) DispatchTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) ([]bind.Object, error) {
	encoded, err := d.contract.mcmsEncoder.DispatchTimelockExecuteBatch(timelock, clock, timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]bind.Object)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []bind.Object, got %T", results[0])
	}
	return result, nil
}

// DispatchTimelockBypasserExecuteBatch executes the dispatch_timelock_bypasser_execute_batch Move function using DevInspect to get return values.
//
// Returns: vector<ExecutingCallbackParams>
func (d *McmsDevInspect) DispatchTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelockCallbackParams TimelockCallbackParams) ([]bind.Object, error) {
	encoded, err := d.contract.mcmsEncoder.DispatchTimelockBypasserExecuteBatch(timelockCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]bind.Object)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []bind.Object, got %T", results[0])
	}
	return result, nil
}

// ExecuteDispatchToDeployer executes the execute_dispatch_to_deployer Move function using DevInspect to get return values.
//
// Returns: UpgradeTicket
func (d *McmsDevInspect) ExecuteDispatchToDeployer(ctx context.Context, opts *bind.CallOpts, registry bind.Object, deployerState bind.Object, executingCallbackParams bind.Object) (bind.Object, error) {
	encoded, err := d.contract.mcmsEncoder.ExecuteDispatchToDeployer(registry, deployerState, executingCallbackParams)
	if err != nil {
		return bind.Object{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return bind.Object{}, err
	}
	if len(results) == 0 {
		return bind.Object{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bind.Object)
	if !ok {
		return bind.Object{}, fmt.Errorf("unexpected return type: expected bind.Object, got %T", results[0])
	}
	return result, nil
}

// ExecuteTimelockExecuteBatch executes the execute_timelock_execute_batch Move function using DevInspect to get return values.
//
// Returns: vector<ExecutingCallbackParams>
func (d *McmsDevInspect) ExecuteTimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) ([]bind.Object, error) {
	encoded, err := d.contract.mcmsEncoder.ExecuteTimelockExecuteBatch(timelock, clock, executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]bind.Object)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []bind.Object, got %T", results[0])
	}
	return result, nil
}

// ExecuteTimelockBypasserExecuteBatch executes the execute_timelock_bypasser_execute_batch Move function using DevInspect to get return values.
//
// Returns: vector<ExecutingCallbackParams>
func (d *McmsDevInspect) ExecuteTimelockBypasserExecuteBatch(ctx context.Context, opts *bind.CallOpts, executingCallbackParams bind.Object) ([]bind.Object, error) {
	encoded, err := d.contract.mcmsEncoder.ExecuteTimelockBypasserExecuteBatch(executingCallbackParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]bind.Object)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []bind.Object, got %T", results[0])
	}
	return result, nil
}

// VerifyMerkleProof executes the verify_merkle_proof Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDevInspect) VerifyMerkleProof(ctx context.Context, opts *bind.CallOpts, proof [][]byte, root []byte, leaf []byte) (bool, error) {
	encoded, err := d.contract.mcmsEncoder.VerifyMerkleProof(proof, root, leaf)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// ComputeEthMessageHash executes the compute_eth_message_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsDevInspect) ComputeEthMessageHash(ctx context.Context, opts *bind.CallOpts, root []byte, validUntil uint64) ([]byte, error) {
	encoded, err := d.contract.mcmsEncoder.ComputeEthMessageHash(root, validUntil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// HashOpLeaf executes the hash_op_leaf Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsDevInspect) HashOpLeaf(ctx context.Context, opts *bind.CallOpts, domainSeparator []byte, op Op) ([]byte, error) {
	encoded, err := d.contract.mcmsEncoder.HashOpLeaf(domainSeparator, op)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// SeenSignedHashes executes the seen_signed_hashes Move function using DevInspect to get return values.
//
// Returns: VecMap<vector<u8>, bool>
func (d *McmsDevInspect) SeenSignedHashes(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (bind.Object, error) {
	encoded, err := d.contract.mcmsEncoder.SeenSignedHashes(state, role)
	if err != nil {
		return bind.Object{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return bind.Object{}, err
	}
	if len(results) == 0 {
		return bind.Object{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bind.Object)
	if !ok {
		return bind.Object{}, fmt.Errorf("unexpected return type: expected bind.Object, got %T", results[0])
	}
	return result, nil
}

// ExpiringRootAndOpCount executes the expiring_root_and_op_count Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<u8>
//	[1]: u64
//	[2]: u64
func (d *McmsDevInspect) ExpiringRootAndOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) ([]any, error) {
	encoded, err := d.contract.mcmsEncoder.ExpiringRootAndOpCount(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// RootMetadata executes the root_metadata Move function using DevInspect to get return values.
//
// Returns: RootMetadata
func (d *McmsDevInspect) RootMetadata(ctx context.Context, opts *bind.CallOpts, multisig Multisig) (RootMetadata, error) {
	encoded, err := d.contract.mcmsEncoder.RootMetadata(multisig)
	if err != nil {
		return RootMetadata{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return RootMetadata{}, err
	}
	if len(results) == 0 {
		return RootMetadata{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(RootMetadata)
	if !ok {
		return RootMetadata{}, fmt.Errorf("unexpected return type: expected RootMetadata, got %T", results[0])
	}
	return result, nil
}

// GetRootMetadata executes the get_root_metadata Move function using DevInspect to get return values.
//
// Returns: RootMetadata
func (d *McmsDevInspect) GetRootMetadata(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (RootMetadata, error) {
	encoded, err := d.contract.mcmsEncoder.GetRootMetadata(state, role)
	if err != nil {
		return RootMetadata{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return RootMetadata{}, err
	}
	if len(results) == 0 {
		return RootMetadata{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(RootMetadata)
	if !ok {
		return RootMetadata{}, fmt.Errorf("unexpected return type: expected RootMetadata, got %T", results[0])
	}
	return result, nil
}

// GetOpCount executes the get_op_count Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) GetOpCount(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.GetOpCount(state, role)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// GetRoot executes the get_root Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<u8>
//	[1]: u64
func (d *McmsDevInspect) GetRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) ([]any, error) {
	encoded, err := d.contract.mcmsEncoder.GetRoot(state, role)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetConfig executes the get_config Move function using DevInspect to get return values.
//
// Returns: Config
func (d *McmsDevInspect) GetConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, role byte) (Config, error) {
	encoded, err := d.contract.mcmsEncoder.GetConfig(state, role)
	if err != nil {
		return Config{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return Config{}, err
	}
	if len(results) == 0 {
		return Config{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(Config)
	if !ok {
		return Config{}, fmt.Errorf("unexpected return type: expected Config, got %T", results[0])
	}
	return result, nil
}

// NumGroups executes the num_groups Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) NumGroups(ctx context.Context, opts *bind.CallOpts) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.NumGroups()
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// MaxNumSigners executes the max_num_signers Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) MaxNumSigners(ctx context.Context, opts *bind.CallOpts) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.MaxNumSigners()
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// BypasserRole executes the bypasser_role Move function using DevInspect to get return values.
//
// Returns: u8
func (d *McmsDevInspect) BypasserRole(ctx context.Context, opts *bind.CallOpts) (byte, error) {
	encoded, err := d.contract.mcmsEncoder.BypasserRole()
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// CancellerRole executes the canceller_role Move function using DevInspect to get return values.
//
// Returns: u8
func (d *McmsDevInspect) CancellerRole(ctx context.Context, opts *bind.CallOpts) (byte, error) {
	encoded, err := d.contract.mcmsEncoder.CancellerRole()
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// ProposerRole executes the proposer_role Move function using DevInspect to get return values.
//
// Returns: u8
func (d *McmsDevInspect) ProposerRole(ctx context.Context, opts *bind.CallOpts) (byte, error) {
	encoded, err := d.contract.mcmsEncoder.ProposerRole()
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// TimelockRole executes the timelock_role Move function using DevInspect to get return values.
//
// Returns: u8
func (d *McmsDevInspect) TimelockRole(ctx context.Context, opts *bind.CallOpts) (byte, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockRole()
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// IsValidRole executes the is_valid_role Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDevInspect) IsValidRole(ctx context.Context, opts *bind.CallOpts, role byte) (bool, error) {
	encoded, err := d.contract.mcmsEncoder.IsValidRole(role)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// ZeroHash executes the zero_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsDevInspect) ZeroHash(ctx context.Context, opts *bind.CallOpts) ([]byte, error) {
	encoded, err := d.contract.mcmsEncoder.ZeroHash()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// Role executes the role Move function using DevInspect to get return values.
//
// Returns: u8
func (d *McmsDevInspect) Role(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (byte, error) {
	encoded, err := d.contract.mcmsEncoder.Role(rootMetadata)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// ChainId executes the chain_id Move function using DevInspect to get return values.
//
// Returns: u256
func (d *McmsDevInspect) ChainId(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (*big.Int, error) {
	encoded, err := d.contract.mcmsEncoder.ChainId(rootMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *big.Int, got %T", results[0])
	}
	return result, nil
}

// RootMetadataMultisig executes the root_metadata_multisig Move function using DevInspect to get return values.
//
// Returns: address
func (d *McmsDevInspect) RootMetadataMultisig(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (string, error) {
	encoded, err := d.contract.mcmsEncoder.RootMetadataMultisig(rootMetadata)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// PreOpCount executes the pre_op_count Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) PreOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.PreOpCount(rootMetadata)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// PostOpCount executes the post_op_count Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) PostOpCount(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.PostOpCount(rootMetadata)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// OverridePreviousRoot executes the override_previous_root Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDevInspect) OverridePreviousRoot(ctx context.Context, opts *bind.CallOpts, rootMetadata RootMetadata) (bool, error) {
	encoded, err := d.contract.mcmsEncoder.OverridePreviousRoot(rootMetadata)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// ConfigSigners executes the config_signers Move function using DevInspect to get return values.
//
// Returns: vector<Signer>
func (d *McmsDevInspect) ConfigSigners(ctx context.Context, opts *bind.CallOpts, config Config) ([]Signer, error) {
	encoded, err := d.contract.mcmsEncoder.ConfigSigners(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]Signer)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []Signer, got %T", results[0])
	}
	return result, nil
}

// ConfigGroupQuorums executes the config_group_quorums Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsDevInspect) ConfigGroupQuorums(ctx context.Context, opts *bind.CallOpts, config Config) ([]byte, error) {
	encoded, err := d.contract.mcmsEncoder.ConfigGroupQuorums(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// ConfigGroupParents executes the config_group_parents Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsDevInspect) ConfigGroupParents(ctx context.Context, opts *bind.CallOpts, config Config) ([]byte, error) {
	encoded, err := d.contract.mcmsEncoder.ConfigGroupParents(config)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// TimelockExecuteBatch executes the timelock_execute_batch Move function using DevInspect to get return values.
//
// Returns: vector<ExecutingCallbackParams>
func (d *McmsDevInspect) TimelockExecuteBatch(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) ([]bind.Object, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockExecuteBatch(timelock, clock, targets, moduleNames, functionNames, datas, predecessor, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]bind.Object)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []bind.Object, got %T", results[0])
	}
	return result, nil
}

// TimelockGetBlockedFunction executes the timelock_get_blocked_function Move function using DevInspect to get return values.
//
// Returns: Function
func (d *McmsDevInspect) TimelockGetBlockedFunction(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, index uint64) (Function, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockGetBlockedFunction(timelock, index)
	if err != nil {
		return Function{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return Function{}, err
	}
	if len(results) == 0 {
		return Function{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(Function)
	if !ok {
		return Function{}, fmt.Errorf("unexpected return type: expected Function, got %T", results[0])
	}
	return result, nil
}

// TimelockIsOperation executes the timelock_is_operation Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDevInspect) TimelockIsOperation(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (bool, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockIsOperation(timelock, id)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// TimelockIsOperationPending executes the timelock_is_operation_pending Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDevInspect) TimelockIsOperationPending(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (bool, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockIsOperationPending(timelock, id)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// TimelockIsOperationReady executes the timelock_is_operation_ready Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDevInspect) TimelockIsOperationReady(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, clock bind.Object, id []byte) (bool, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockIsOperationReady(timelock, clock, id)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// TimelockIsOperationDone executes the timelock_is_operation_done Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDevInspect) TimelockIsOperationDone(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (bool, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockIsOperationDone(timelock, id)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// TimelockGetTimestamp executes the timelock_get_timestamp Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) TimelockGetTimestamp(ctx context.Context, opts *bind.CallOpts, timelock bind.Object, id []byte) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockGetTimestamp(timelock, id)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// TimelockMinDelay executes the timelock_min_delay Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) TimelockMinDelay(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockMinDelay(timelock)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// TimelockGetBlockedFunctions executes the timelock_get_blocked_functions Move function using DevInspect to get return values.
//
// Returns: vector<Function>
func (d *McmsDevInspect) TimelockGetBlockedFunctions(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) ([]Function, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockGetBlockedFunctions(timelock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]Function)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []Function, got %T", results[0])
	}
	return result, nil
}

// TimelockGetBlockedFunctionsCount executes the timelock_get_blocked_functions_count Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsDevInspect) TimelockGetBlockedFunctionsCount(ctx context.Context, opts *bind.CallOpts, timelock bind.Object) (uint64, error) {
	encoded, err := d.contract.mcmsEncoder.TimelockGetBlockedFunctionsCount(timelock)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// CreateCalls executes the create_calls Move function using DevInspect to get return values.
//
// Returns: vector<Call>
func (d *McmsDevInspect) CreateCalls(ctx context.Context, opts *bind.CallOpts, targets []string, moduleNames []string, functionNames []string, datas [][]byte) ([]Call, error) {
	encoded, err := d.contract.mcmsEncoder.CreateCalls(targets, moduleNames, functionNames, datas)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]Call)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []Call, got %T", results[0])
	}
	return result, nil
}

// HashOperationBatch executes the hash_operation_batch Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsDevInspect) HashOperationBatch(ctx context.Context, opts *bind.CallOpts, calls []Call, predecessor []byte, salt []byte) ([]byte, error) {
	encoded, err := d.contract.mcmsEncoder.HashOperationBatch(calls, predecessor, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// SignerView executes the signer_view Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<u8>
//	[1]: u8
//	[2]: u8
func (d *McmsDevInspect) SignerView(ctx context.Context, opts *bind.CallOpts, signer Signer) ([]any, error) {
	encoded, err := d.contract.mcmsEncoder.SignerView(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// FunctionName executes the function_name Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *McmsDevInspect) FunctionName(ctx context.Context, opts *bind.CallOpts, function Function) (string, error) {
	encoded, err := d.contract.mcmsEncoder.FunctionName(function)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// ModuleName executes the module_name Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *McmsDevInspect) ModuleName(ctx context.Context, opts *bind.CallOpts, function Function) (string, error) {
	encoded, err := d.contract.mcmsEncoder.ModuleName(function)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// Target executes the target Move function using DevInspect to get return values.
//
// Returns: address
func (d *McmsDevInspect) Target(ctx context.Context, opts *bind.CallOpts, function Function) (string, error) {
	encoded, err := d.contract.mcmsEncoder.Target(function)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// Data executes the data Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsDevInspect) Data(ctx context.Context, opts *bind.CallOpts, call Call) ([]byte, error) {
	encoded, err := d.contract.mcmsEncoder.Data(call)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

type mcmsEncoder struct {
	*bind.BoundContract
}

// SetRoot encodes a call to the set_root Move function.
func (c mcmsEncoder) SetRoot(state bind.Object, clock bind.Object, role byte, root []byte, validUntil uint64, chainId *big.Int, multisigAddr string, preOpCount uint64, postOpCount uint64, overridePreviousRoot bool, metadataProof [][]byte, signatures [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_root", typeArgsList, typeParamsList, []string{
		"&mut MultisigState",
		"&Clock",
		"u8",
		"vector<u8>",
		"u64",
		"u256",
		"address",
		"u64",
		"u64",
		"bool",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
	}, []any{
		state,
		clock,
		role,
		root,
		validUntil,
		chainId,
		multisigAddr,
		preOpCount,
		postOpCount,
		overridePreviousRoot,
		metadataProof,
		signatures,
	}, nil)
}

// SetRootWithArgs encodes a call to the set_root Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) SetRootWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut MultisigState",
		"&Clock",
		"u8",
		"vector<u8>",
		"u64",
		"u256",
		"address",
		"u64",
		"u64",
		"bool",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_root", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Execute encodes a call to the execute Move function.
func (c mcmsEncoder) Execute(state bind.Object, clock bind.Object, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute", typeArgsList, typeParamsList, []string{
		"&mut MultisigState",
		"&Clock",
		"u8",
		"u256",
		"address",
		"u64",
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
		"vector<vector<u8>>",
	}, []any{
		state,
		clock,
		role,
		chainId,
		multisigAddr,
		nonce,
		to,
		moduleName,
		functionName,
		data,
		proof,
	}, []string{
		"mcms::mcms::TimelockCallbackParams",
	})
}

// ExecuteWithArgs encodes a call to the execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut MultisigState",
		"&Clock",
		"u8",
		"u256",
		"address",
		"u64",
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute", typeArgsList, typeParamsList, expectedParams, args, []string{
		"mcms::mcms::TimelockCallbackParams",
	})
}

// DispatchTimelockScheduleBatch encodes a call to the dispatch_timelock_schedule_batch Move function.
func (c mcmsEncoder) DispatchTimelockScheduleBatch(timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_schedule_batch", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"&Clock",
		"mcms::mcms::TimelockCallbackParams",
	}, []any{
		timelock,
		clock,
		timelockCallbackParams,
	}, nil)
}

// DispatchTimelockScheduleBatchWithArgs encodes a call to the dispatch_timelock_schedule_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DispatchTimelockScheduleBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"&Clock",
		"mcms::mcms::TimelockCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_schedule_batch", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DispatchTimelockExecuteBatch encodes a call to the dispatch_timelock_execute_batch Move function.
func (c mcmsEncoder) DispatchTimelockExecuteBatch(timelock bind.Object, clock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_execute_batch", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"&Clock",
		"mcms::mcms::TimelockCallbackParams",
	}, []any{
		timelock,
		clock,
		timelockCallbackParams,
	}, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// DispatchTimelockExecuteBatchWithArgs encodes a call to the dispatch_timelock_execute_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DispatchTimelockExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"&Clock",
		"mcms::mcms::TimelockCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_execute_batch", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// DispatchTimelockBypasserExecuteBatch encodes a call to the dispatch_timelock_bypasser_execute_batch Move function.
func (c mcmsEncoder) DispatchTimelockBypasserExecuteBatch(timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_bypasser_execute_batch", typeArgsList, typeParamsList, []string{
		"mcms::mcms::TimelockCallbackParams",
	}, []any{
		timelockCallbackParams,
	}, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// DispatchTimelockBypasserExecuteBatchWithArgs encodes a call to the dispatch_timelock_bypasser_execute_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DispatchTimelockBypasserExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms::TimelockCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_bypasser_execute_batch", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// DispatchTimelockCancel encodes a call to the dispatch_timelock_cancel Move function.
func (c mcmsEncoder) DispatchTimelockCancel(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_cancel", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}, []any{
		timelock,
		timelockCallbackParams,
	}, nil)
}

// DispatchTimelockCancelWithArgs encodes a call to the dispatch_timelock_cancel Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DispatchTimelockCancelWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_cancel", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DispatchTimelockUpdateMinDelay encodes a call to the dispatch_timelock_update_min_delay Move function.
func (c mcmsEncoder) DispatchTimelockUpdateMinDelay(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_update_min_delay", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}, []any{
		timelock,
		timelockCallbackParams,
	}, nil)
}

// DispatchTimelockUpdateMinDelayWithArgs encodes a call to the dispatch_timelock_update_min_delay Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DispatchTimelockUpdateMinDelayWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_update_min_delay", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DispatchTimelockBlockFunction encodes a call to the dispatch_timelock_block_function Move function.
func (c mcmsEncoder) DispatchTimelockBlockFunction(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_block_function", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}, []any{
		timelock,
		timelockCallbackParams,
	}, nil)
}

// DispatchTimelockBlockFunctionWithArgs encodes a call to the dispatch_timelock_block_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DispatchTimelockBlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_block_function", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DispatchTimelockUnblockFunction encodes a call to the dispatch_timelock_unblock_function Move function.
func (c mcmsEncoder) DispatchTimelockUnblockFunction(timelock bind.Object, timelockCallbackParams TimelockCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_unblock_function", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}, []any{
		timelock,
		timelockCallbackParams,
	}, nil)
}

// DispatchTimelockUnblockFunctionWithArgs encodes a call to the dispatch_timelock_unblock_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DispatchTimelockUnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"mcms::mcms::TimelockCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("dispatch_timelock_unblock_function", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteDispatchToAccount encodes a call to the execute_dispatch_to_account Move function.
func (c mcmsEncoder) ExecuteDispatchToAccount(registry bind.Object, accountState bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_dispatch_to_account", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"&mut AccountState",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		accountState,
		executingCallbackParams,
	}, nil)
}

// ExecuteDispatchToAccountWithArgs encodes a call to the execute_dispatch_to_account Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteDispatchToAccountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"&mut AccountState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_dispatch_to_account", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteDispatchToDeployer encodes a call to the execute_dispatch_to_deployer Move function.
func (c mcmsEncoder) ExecuteDispatchToDeployer(registry bind.Object, deployerState bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_dispatch_to_deployer", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		deployerState,
		executingCallbackParams,
	}, []string{
		"UpgradeTicket",
	})
}

// ExecuteDispatchToDeployerWithArgs encodes a call to the execute_dispatch_to_deployer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteDispatchToDeployerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_dispatch_to_deployer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"UpgradeTicket",
	})
}

// ExecuteTimelockScheduleBatch encodes a call to the execute_timelock_schedule_batch Move function.
func (c mcmsEncoder) ExecuteTimelockScheduleBatch(timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_schedule_batch", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"&Clock",
		"ExecutingCallbackParams",
	}, []any{
		timelock,
		clock,
		executingCallbackParams,
	}, nil)
}

// ExecuteTimelockScheduleBatchWithArgs encodes a call to the execute_timelock_schedule_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteTimelockScheduleBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"&Clock",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_schedule_batch", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteTimelockExecuteBatch encodes a call to the execute_timelock_execute_batch Move function.
func (c mcmsEncoder) ExecuteTimelockExecuteBatch(timelock bind.Object, clock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_execute_batch", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"&Clock",
		"ExecutingCallbackParams",
	}, []any{
		timelock,
		clock,
		executingCallbackParams,
	}, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// ExecuteTimelockExecuteBatchWithArgs encodes a call to the execute_timelock_execute_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteTimelockExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"&Clock",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_execute_batch", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// ExecuteTimelockBypasserExecuteBatch encodes a call to the execute_timelock_bypasser_execute_batch Move function.
func (c mcmsEncoder) ExecuteTimelockBypasserExecuteBatch(executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_bypasser_execute_batch", typeArgsList, typeParamsList, []string{
		"ExecutingCallbackParams",
	}, []any{
		executingCallbackParams,
	}, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// ExecuteTimelockBypasserExecuteBatchWithArgs encodes a call to the execute_timelock_bypasser_execute_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteTimelockBypasserExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_bypasser_execute_batch", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// ExecuteTimelockCancel encodes a call to the execute_timelock_cancel Move function.
func (c mcmsEncoder) ExecuteTimelockCancel(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_cancel", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}, []any{
		timelock,
		executingCallbackParams,
	}, nil)
}

// ExecuteTimelockCancelWithArgs encodes a call to the execute_timelock_cancel Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteTimelockCancelWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_cancel", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteTimelockUpdateMinDelay encodes a call to the execute_timelock_update_min_delay Move function.
func (c mcmsEncoder) ExecuteTimelockUpdateMinDelay(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_update_min_delay", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}, []any{
		timelock,
		executingCallbackParams,
	}, nil)
}

// ExecuteTimelockUpdateMinDelayWithArgs encodes a call to the execute_timelock_update_min_delay Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteTimelockUpdateMinDelayWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_update_min_delay", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteTimelockBlockFunction encodes a call to the execute_timelock_block_function Move function.
func (c mcmsEncoder) ExecuteTimelockBlockFunction(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_block_function", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}, []any{
		timelock,
		executingCallbackParams,
	}, nil)
}

// ExecuteTimelockBlockFunctionWithArgs encodes a call to the execute_timelock_block_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteTimelockBlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_block_function", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteTimelockUnblockFunction encodes a call to the execute_timelock_unblock_function Move function.
func (c mcmsEncoder) ExecuteTimelockUnblockFunction(timelock bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_unblock_function", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}, []any{
		timelock,
		executingCallbackParams,
	}, nil)
}

// ExecuteTimelockUnblockFunctionWithArgs encodes a call to the execute_timelock_unblock_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteTimelockUnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_timelock_unblock_function", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteSetConfig encodes a call to the execute_set_config Move function.
func (c mcmsEncoder) ExecuteSetConfig(registry bind.Object, state bind.Object, executingCallbackParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_set_config", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"&mut MultisigState",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		state,
		executingCallbackParams,
	}, nil)
}

// ExecuteSetConfigWithArgs encodes a call to the execute_set_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExecuteSetConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"&mut MultisigState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_set_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetConfig encodes a call to the set_config Move function.
func (c mcmsEncoder) SetConfig(param bind.Object, state bind.Object, role byte, chainId *big.Int, signerAddresses [][]byte, signerGroups []byte, groupQuorums []byte, groupParents []byte, clearRoot bool) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_config", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut MultisigState",
		"u8",
		"u256",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
		"bool",
	}, []any{
		param,
		state,
		role,
		chainId,
		signerAddresses,
		signerGroups,
		groupQuorums,
		groupParents,
		clearRoot,
	}, nil)
}

// SetConfigWithArgs encodes a call to the set_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) SetConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut MultisigState",
		"u8",
		"u256",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
		"bool",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// VerifyMerkleProof encodes a call to the verify_merkle_proof Move function.
func (c mcmsEncoder) VerifyMerkleProof(proof [][]byte, root []byte, leaf []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify_merkle_proof", typeArgsList, typeParamsList, []string{
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<u8>",
	}, []any{
		proof,
		root,
		leaf,
	}, []string{
		"bool",
	})
}

// VerifyMerkleProofWithArgs encodes a call to the verify_merkle_proof Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) VerifyMerkleProofWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify_merkle_proof", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// ComputeEthMessageHash encodes a call to the compute_eth_message_hash Move function.
func (c mcmsEncoder) ComputeEthMessageHash(root []byte, validUntil uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("compute_eth_message_hash", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"u64",
	}, []any{
		root,
		validUntil,
	}, []string{
		"vector<u8>",
	})
}

// ComputeEthMessageHashWithArgs encodes a call to the compute_eth_message_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ComputeEthMessageHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("compute_eth_message_hash", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// HashOpLeaf encodes a call to the hash_op_leaf Move function.
func (c mcmsEncoder) HashOpLeaf(domainSeparator []byte, op Op) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("hash_op_leaf", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"mcms::mcms::Op",
	}, []any{
		domainSeparator,
		op,
	}, []string{
		"vector<u8>",
	})
}

// HashOpLeafWithArgs encodes a call to the hash_op_leaf Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) HashOpLeafWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"mcms::mcms::Op",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("hash_op_leaf", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// SeenSignedHashes encodes a call to the seen_signed_hashes Move function.
func (c mcmsEncoder) SeenSignedHashes(state bind.Object, role byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("seen_signed_hashes", typeArgsList, typeParamsList, []string{
		"&MultisigState",
		"u8",
	}, []any{
		state,
		role,
	}, []string{
		"VecMap<vector<u8>, bool>",
	})
}

// SeenSignedHashesWithArgs encodes a call to the seen_signed_hashes Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) SeenSignedHashesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&MultisigState",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("seen_signed_hashes", typeArgsList, typeParamsList, expectedParams, args, []string{
		"VecMap<vector<u8>, bool>",
	})
}

// ExpiringRootAndOpCount encodes a call to the expiring_root_and_op_count Move function.
func (c mcmsEncoder) ExpiringRootAndOpCount(state bind.Object, role byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("expiring_root_and_op_count", typeArgsList, typeParamsList, []string{
		"&MultisigState",
		"u8",
	}, []any{
		state,
		role,
	}, []string{
		"vector<u8>",
		"u64",
		"u64",
	})
}

// ExpiringRootAndOpCountWithArgs encodes a call to the expiring_root_and_op_count Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ExpiringRootAndOpCountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&MultisigState",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("expiring_root_and_op_count", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
		"u64",
		"u64",
	})
}

// RootMetadata encodes a call to the root_metadata Move function.
func (c mcmsEncoder) RootMetadata(multisig Multisig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("root_metadata", typeArgsList, typeParamsList, []string{
		"&Multisig",
	}, []any{
		multisig,
	}, []string{
		"mcms::mcms::RootMetadata",
	})
}

// RootMetadataWithArgs encodes a call to the root_metadata Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) RootMetadataWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Multisig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("root_metadata", typeArgsList, typeParamsList, expectedParams, args, []string{
		"mcms::mcms::RootMetadata",
	})
}

// GetRootMetadata encodes a call to the get_root_metadata Move function.
func (c mcmsEncoder) GetRootMetadata(state bind.Object, role byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_root_metadata", typeArgsList, typeParamsList, []string{
		"&MultisigState",
		"u8",
	}, []any{
		state,
		role,
	}, []string{
		"mcms::mcms::RootMetadata",
	})
}

// GetRootMetadataWithArgs encodes a call to the get_root_metadata Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) GetRootMetadataWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&MultisigState",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_root_metadata", typeArgsList, typeParamsList, expectedParams, args, []string{
		"mcms::mcms::RootMetadata",
	})
}

// GetOpCount encodes a call to the get_op_count Move function.
func (c mcmsEncoder) GetOpCount(state bind.Object, role byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_op_count", typeArgsList, typeParamsList, []string{
		"&MultisigState",
		"u8",
	}, []any{
		state,
		role,
	}, []string{
		"u64",
	})
}

// GetOpCountWithArgs encodes a call to the get_op_count Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) GetOpCountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&MultisigState",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_op_count", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetRoot encodes a call to the get_root Move function.
func (c mcmsEncoder) GetRoot(state bind.Object, role byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_root", typeArgsList, typeParamsList, []string{
		"&MultisigState",
		"u8",
	}, []any{
		state,
		role,
	}, []string{
		"vector<u8>",
		"u64",
	})
}

// GetRootWithArgs encodes a call to the get_root Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) GetRootWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&MultisigState",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_root", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
		"u64",
	})
}

// GetConfig encodes a call to the get_config Move function.
func (c mcmsEncoder) GetConfig(state bind.Object, role byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_config", typeArgsList, typeParamsList, []string{
		"&MultisigState",
		"u8",
	}, []any{
		state,
		role,
	}, []string{
		"mcms::mcms::Config",
	})
}

// GetConfigWithArgs encodes a call to the get_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) GetConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&MultisigState",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"mcms::mcms::Config",
	})
}

// NumGroups encodes a call to the num_groups Move function.
func (c mcmsEncoder) NumGroups() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("num_groups", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"u64",
	})
}

// NumGroupsWithArgs encodes a call to the num_groups Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) NumGroupsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("num_groups", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// MaxNumSigners encodes a call to the max_num_signers Move function.
func (c mcmsEncoder) MaxNumSigners() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("max_num_signers", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"u64",
	})
}

// MaxNumSignersWithArgs encodes a call to the max_num_signers Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) MaxNumSignersWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("max_num_signers", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// BypasserRole encodes a call to the bypasser_role Move function.
func (c mcmsEncoder) BypasserRole() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("bypasser_role", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"u8",
	})
}

// BypasserRoleWithArgs encodes a call to the bypasser_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) BypasserRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("bypasser_role", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// CancellerRole encodes a call to the canceller_role Move function.
func (c mcmsEncoder) CancellerRole() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("canceller_role", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"u8",
	})
}

// CancellerRoleWithArgs encodes a call to the canceller_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) CancellerRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("canceller_role", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// ProposerRole encodes a call to the proposer_role Move function.
func (c mcmsEncoder) ProposerRole() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("proposer_role", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"u8",
	})
}

// ProposerRoleWithArgs encodes a call to the proposer_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ProposerRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("proposer_role", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// TimelockRole encodes a call to the timelock_role Move function.
func (c mcmsEncoder) TimelockRole() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_role", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"u8",
	})
}

// TimelockRoleWithArgs encodes a call to the timelock_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_role", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// IsValidRole encodes a call to the is_valid_role Move function.
func (c mcmsEncoder) IsValidRole(role byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_valid_role", typeArgsList, typeParamsList, []string{
		"u8",
	}, []any{
		role,
	}, []string{
		"bool",
	})
}

// IsValidRoleWithArgs encodes a call to the is_valid_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) IsValidRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_valid_role", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// ZeroHash encodes a call to the zero_hash Move function.
func (c mcmsEncoder) ZeroHash() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("zero_hash", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"vector<u8>",
	})
}

// ZeroHashWithArgs encodes a call to the zero_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ZeroHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("zero_hash", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// Role encodes a call to the role Move function.
func (c mcmsEncoder) Role(rootMetadata RootMetadata) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("role", typeArgsList, typeParamsList, []string{
		"&RootMetadata",
	}, []any{
		rootMetadata,
	}, []string{
		"u8",
	})
}

// RoleWithArgs encodes a call to the role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) RoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RootMetadata",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("role", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// ChainId encodes a call to the chain_id Move function.
func (c mcmsEncoder) ChainId(rootMetadata RootMetadata) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("chain_id", typeArgsList, typeParamsList, []string{
		"&RootMetadata",
	}, []any{
		rootMetadata,
	}, []string{
		"u256",
	})
}

// ChainIdWithArgs encodes a call to the chain_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ChainIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RootMetadata",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("chain_id", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u256",
	})
}

// RootMetadataMultisig encodes a call to the root_metadata_multisig Move function.
func (c mcmsEncoder) RootMetadataMultisig(rootMetadata RootMetadata) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("root_metadata_multisig", typeArgsList, typeParamsList, []string{
		"&RootMetadata",
	}, []any{
		rootMetadata,
	}, []string{
		"address",
	})
}

// RootMetadataMultisigWithArgs encodes a call to the root_metadata_multisig Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) RootMetadataMultisigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RootMetadata",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("root_metadata_multisig", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// PreOpCount encodes a call to the pre_op_count Move function.
func (c mcmsEncoder) PreOpCount(rootMetadata RootMetadata) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pre_op_count", typeArgsList, typeParamsList, []string{
		"&RootMetadata",
	}, []any{
		rootMetadata,
	}, []string{
		"u64",
	})
}

// PreOpCountWithArgs encodes a call to the pre_op_count Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) PreOpCountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RootMetadata",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pre_op_count", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// PostOpCount encodes a call to the post_op_count Move function.
func (c mcmsEncoder) PostOpCount(rootMetadata RootMetadata) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("post_op_count", typeArgsList, typeParamsList, []string{
		"&RootMetadata",
	}, []any{
		rootMetadata,
	}, []string{
		"u64",
	})
}

// PostOpCountWithArgs encodes a call to the post_op_count Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) PostOpCountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RootMetadata",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("post_op_count", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// OverridePreviousRoot encodes a call to the override_previous_root Move function.
func (c mcmsEncoder) OverridePreviousRoot(rootMetadata RootMetadata) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("override_previous_root", typeArgsList, typeParamsList, []string{
		"&RootMetadata",
	}, []any{
		rootMetadata,
	}, []string{
		"bool",
	})
}

// OverridePreviousRootWithArgs encodes a call to the override_previous_root Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) OverridePreviousRootWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RootMetadata",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("override_previous_root", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// ConfigSigners encodes a call to the config_signers Move function.
func (c mcmsEncoder) ConfigSigners(config Config) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_signers", typeArgsList, typeParamsList, []string{
		"&Config",
	}, []any{
		config,
	}, []string{
		"vector<mcms::mcms::Signer>",
	})
}

// ConfigSignersWithArgs encodes a call to the config_signers Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ConfigSignersWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Config",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_signers", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<mcms::mcms::Signer>",
	})
}

// ConfigGroupQuorums encodes a call to the config_group_quorums Move function.
func (c mcmsEncoder) ConfigGroupQuorums(config Config) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_group_quorums", typeArgsList, typeParamsList, []string{
		"&Config",
	}, []any{
		config,
	}, []string{
		"vector<u8>",
	})
}

// ConfigGroupQuorumsWithArgs encodes a call to the config_group_quorums Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ConfigGroupQuorumsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Config",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_group_quorums", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// ConfigGroupParents encodes a call to the config_group_parents Move function.
func (c mcmsEncoder) ConfigGroupParents(config Config) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_group_parents", typeArgsList, typeParamsList, []string{
		"&Config",
	}, []any{
		config,
	}, []string{
		"vector<u8>",
	})
}

// ConfigGroupParentsWithArgs encodes a call to the config_group_parents Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ConfigGroupParentsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Config",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_group_parents", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// TimelockExecuteBatch encodes a call to the timelock_execute_batch Move function.
func (c mcmsEncoder) TimelockExecuteBatch(timelock bind.Object, clock bind.Object, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_execute_batch", typeArgsList, typeParamsList, []string{
		"&mut Timelock",
		"&Clock",
		"vector<address>",
		"vector<0x1::string::String>",
		"vector<0x1::string::String>",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<u8>",
	}, []any{
		timelock,
		clock,
		targets,
		moduleNames,
		functionNames,
		datas,
		predecessor,
		salt,
	}, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// TimelockExecuteBatchWithArgs encodes a call to the timelock_execute_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockExecuteBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Timelock",
		"&Clock",
		"vector<address>",
		"vector<0x1::string::String>",
		"vector<0x1::string::String>",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_execute_batch", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ExecutingCallbackParams>",
	})
}

// TimelockGetBlockedFunction encodes a call to the timelock_get_blocked_function Move function.
func (c mcmsEncoder) TimelockGetBlockedFunction(timelock bind.Object, index uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_blocked_function", typeArgsList, typeParamsList, []string{
		"&Timelock",
		"u64",
	}, []any{
		timelock,
		index,
	}, []string{
		"mcms::mcms::Function",
	})
}

// TimelockGetBlockedFunctionWithArgs encodes a call to the timelock_get_blocked_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockGetBlockedFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_blocked_function", typeArgsList, typeParamsList, expectedParams, args, []string{
		"mcms::mcms::Function",
	})
}

// TimelockIsOperation encodes a call to the timelock_is_operation Move function.
func (c mcmsEncoder) TimelockIsOperation(timelock bind.Object, id []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation", typeArgsList, typeParamsList, []string{
		"&Timelock",
		"vector<u8>",
	}, []any{
		timelock,
		id,
	}, []string{
		"bool",
	})
}

// TimelockIsOperationWithArgs encodes a call to the timelock_is_operation Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockIsOperationWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// TimelockIsOperationPending encodes a call to the timelock_is_operation_pending Move function.
func (c mcmsEncoder) TimelockIsOperationPending(timelock bind.Object, id []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation_pending", typeArgsList, typeParamsList, []string{
		"&Timelock",
		"vector<u8>",
	}, []any{
		timelock,
		id,
	}, []string{
		"bool",
	})
}

// TimelockIsOperationPendingWithArgs encodes a call to the timelock_is_operation_pending Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockIsOperationPendingWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation_pending", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// TimelockIsOperationReady encodes a call to the timelock_is_operation_ready Move function.
func (c mcmsEncoder) TimelockIsOperationReady(timelock bind.Object, clock bind.Object, id []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation_ready", typeArgsList, typeParamsList, []string{
		"&Timelock",
		"&Clock",
		"vector<u8>",
	}, []any{
		timelock,
		clock,
		id,
	}, []string{
		"bool",
	})
}

// TimelockIsOperationReadyWithArgs encodes a call to the timelock_is_operation_ready Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockIsOperationReadyWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
		"&Clock",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation_ready", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// TimelockIsOperationDone encodes a call to the timelock_is_operation_done Move function.
func (c mcmsEncoder) TimelockIsOperationDone(timelock bind.Object, id []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation_done", typeArgsList, typeParamsList, []string{
		"&Timelock",
		"vector<u8>",
	}, []any{
		timelock,
		id,
	}, []string{
		"bool",
	})
}

// TimelockIsOperationDoneWithArgs encodes a call to the timelock_is_operation_done Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockIsOperationDoneWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_is_operation_done", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// TimelockGetTimestamp encodes a call to the timelock_get_timestamp Move function.
func (c mcmsEncoder) TimelockGetTimestamp(timelock bind.Object, id []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_timestamp", typeArgsList, typeParamsList, []string{
		"&Timelock",
		"vector<u8>",
	}, []any{
		timelock,
		id,
	}, []string{
		"u64",
	})
}

// TimelockGetTimestampWithArgs encodes a call to the timelock_get_timestamp Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockGetTimestampWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_timestamp", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// TimelockMinDelay encodes a call to the timelock_min_delay Move function.
func (c mcmsEncoder) TimelockMinDelay(timelock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_min_delay", typeArgsList, typeParamsList, []string{
		"&Timelock",
	}, []any{
		timelock,
	}, []string{
		"u64",
	})
}

// TimelockMinDelayWithArgs encodes a call to the timelock_min_delay Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockMinDelayWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_min_delay", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// TimelockGetBlockedFunctions encodes a call to the timelock_get_blocked_functions Move function.
func (c mcmsEncoder) TimelockGetBlockedFunctions(timelock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_blocked_functions", typeArgsList, typeParamsList, []string{
		"&Timelock",
	}, []any{
		timelock,
	}, []string{
		"vector<mcms::mcms::Function>",
	})
}

// TimelockGetBlockedFunctionsWithArgs encodes a call to the timelock_get_blocked_functions Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockGetBlockedFunctionsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_blocked_functions", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<mcms::mcms::Function>",
	})
}

// TimelockGetBlockedFunctionsCount encodes a call to the timelock_get_blocked_functions_count Move function.
func (c mcmsEncoder) TimelockGetBlockedFunctionsCount(timelock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_blocked_functions_count", typeArgsList, typeParamsList, []string{
		"&Timelock",
	}, []any{
		timelock,
	}, []string{
		"u64",
	})
}

// TimelockGetBlockedFunctionsCountWithArgs encodes a call to the timelock_get_blocked_functions_count Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TimelockGetBlockedFunctionsCountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Timelock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("timelock_get_blocked_functions_count", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// CreateCalls encodes a call to the create_calls Move function.
func (c mcmsEncoder) CreateCalls(targets []string, moduleNames []string, functionNames []string, datas [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_calls", typeArgsList, typeParamsList, []string{
		"vector<address>",
		"vector<0x1::string::String>",
		"vector<0x1::string::String>",
		"vector<vector<u8>>",
	}, []any{
		targets,
		moduleNames,
		functionNames,
		datas,
	}, []string{
		"vector<mcms::mcms::Call>",
	})
}

// CreateCallsWithArgs encodes a call to the create_calls Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) CreateCallsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<address>",
		"vector<0x1::string::String>",
		"vector<0x1::string::String>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_calls", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<mcms::mcms::Call>",
	})
}

// HashOperationBatch encodes a call to the hash_operation_batch Move function.
func (c mcmsEncoder) HashOperationBatch(calls []Call, predecessor []byte, salt []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("hash_operation_batch", typeArgsList, typeParamsList, []string{
		"vector<mcms::mcms::Call>",
		"vector<u8>",
		"vector<u8>",
	}, []any{
		calls,
		predecessor,
		salt,
	}, []string{
		"vector<u8>",
	})
}

// HashOperationBatchWithArgs encodes a call to the hash_operation_batch Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) HashOperationBatchWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<mcms::mcms::Call>",
		"vector<u8>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("hash_operation_batch", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// SignerView encodes a call to the signer_view Move function.
func (c mcmsEncoder) SignerView(signer Signer) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("signer_view", typeArgsList, typeParamsList, []string{
		"&Signer",
	}, []any{
		signer,
	}, []string{
		"vector<u8>",
		"u8",
		"u8",
	})
}

// SignerViewWithArgs encodes a call to the signer_view Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) SignerViewWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Signer",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("signer_view", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
		"u8",
		"u8",
	})
}

// FunctionName encodes a call to the function_name Move function.
func (c mcmsEncoder) FunctionName(function Function) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_name", typeArgsList, typeParamsList, []string{
		"mcms::mcms::Function",
	}, []any{
		function,
	}, []string{
		"0x1::string::String",
	})
}

// FunctionNameWithArgs encodes a call to the function_name Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) FunctionNameWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms::Function",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_name", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// ModuleName encodes a call to the module_name Move function.
func (c mcmsEncoder) ModuleName(function Function) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("module_name", typeArgsList, typeParamsList, []string{
		"mcms::mcms::Function",
	}, []any{
		function,
	}, []string{
		"0x1::string::String",
	})
}

// ModuleNameWithArgs encodes a call to the module_name Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) ModuleNameWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms::Function",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("module_name", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// Target encodes a call to the target Move function.
func (c mcmsEncoder) Target(function Function) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("target", typeArgsList, typeParamsList, []string{
		"mcms::mcms::Function",
	}, []any{
		function,
	}, []string{
		"address",
	})
}

// TargetWithArgs encodes a call to the target Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) TargetWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms::Function",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("target", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// Data encodes a call to the data Move function.
func (c mcmsEncoder) Data(call Call) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("data", typeArgsList, typeParamsList, []string{
		"mcms::mcms::Call",
	}, []any{
		call,
	}, []string{
		"vector<u8>",
	})
}

// DataWithArgs encodes a call to the data Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsEncoder) DataWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms::Call",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("data", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}
