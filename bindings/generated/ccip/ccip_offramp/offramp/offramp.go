// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_offramp

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

const FunctionInfo = `[{"package":"ccip_offramp","module":"offramp","name":"accept_ownership","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"accept_ownership_from_object","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"from","type":"sui::object::UID"}]},{"package":"ccip_offramp","module":"offramp","name":"add_package_id","parameters":[{"name":"state","type":"OffRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"package_id","type":"address"}]},{"package":"ccip_offramp","module":"offramp","name":"apply_source_chain_config_updates","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"source_chains_selector","type":"vector<u64>"},{"name":"source_chains_is_enabled","type":"vector<bool>"},{"name":"source_chains_is_rmn_verification_disabled","type":"vector<bool>"},{"name":"source_chains_on_ramp","type":"vector<vector<u8>>"}]},{"package":"ccip_offramp","module":"offramp","name":"calculate_message_hash","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"message_id","type":"vector<u8>"},{"name":"source_chain_selector","type":"u64"},{"name":"dest_chain_selector","type":"u64"},{"name":"sequence_number","type":"u64"},{"name":"nonce","type":"u64"},{"name":"sender","type":"vector<u8>"},{"name":"receiver","type":"address"},{"name":"on_ramp","type":"vector<u8>"},{"name":"data","type":"vector<u8>"},{"name":"gas_limit","type":"u256"},{"name":"token_receiver","type":"address"},{"name":"source_pool_addresses","type":"vector<vector<u8>>"},{"name":"dest_token_addresses","type":"vector<address>"},{"name":"dest_gas_amounts","type":"vector<u32>"},{"name":"extra_datas","type":"vector<vector<u8>>"},{"name":"amounts","type":"vector<u256>"}]},{"package":"ccip_offramp","module":"offramp","name":"calculate_metadata_hash","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"source_chain_selector","type":"u64"},{"name":"dest_chain_selector","type":"u64"},{"name":"on_ramp","type":"vector<u8>"}]},{"package":"ccip_offramp","module":"offramp","name":"commit","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"clock","type":"clock::Clock"},{"name":"report_context","type":"vector<vector<u8>>"},{"name":"report","type":"vector<u8>"},{"name":"signatures","type":"vector<vector<u8>>"}]},{"package":"ccip_offramp","module":"offramp","name":"config_signers","parameters":[{"name":"state","type":"OCRConfig"}]},{"package":"ccip_offramp","module":"offramp","name":"config_transmitters","parameters":[{"name":"state","type":"OCRConfig"}]},{"package":"ccip_offramp","module":"offramp","name":"execute_ownership_transfer","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"OffRampState"},{"name":"to","type":"address"}]},{"package":"ccip_offramp","module":"offramp","name":"execute_ownership_transfer_to_mcms","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"OffRampState"},{"name":"registry","type":"Registry"},{"name":"to","type":"address"}]},{"package":"ccip_offramp","module":"offramp","name":"finish_execute","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"receiver_params","type":"osh::ReceiverParams"}]},{"package":"ccip_offramp","module":"offramp","name":"get_all_source_chain_configs","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"get_ccip_package_id","parameters":null},{"package":"ccip_offramp","module":"offramp","name":"get_dynamic_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"get_dynamic_config_fields","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"cfg","type":"DynamicConfig"}]},{"package":"ccip_offramp","module":"offramp","name":"get_execution_state","parameters":[{"name":"state","type":"OffRampState"},{"name":"source_chain_selector","type":"u64"},{"name":"sequence_number","type":"u64"}]},{"package":"ccip_offramp","module":"offramp","name":"get_merkle_root","parameters":[{"name":"state","type":"OffRampState"},{"name":"root","type":"vector<u8>"}]},{"package":"ccip_offramp","module":"offramp","name":"get_ocr3_base","parameters":[{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"get_source_chain_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"source_chain_selector","type":"u64"}]},{"package":"ccip_offramp","module":"offramp","name":"get_source_chain_config_fields","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"source_chain_config","type":"SourceChainConfig"}]},{"package":"ccip_offramp","module":"offramp","name":"get_static_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"get_static_config_fields","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"cfg","type":"StaticConfig"}]},{"package":"ccip_offramp","module":"offramp","name":"has_pending_transfer","parameters":[{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"init_execute","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"clock","type":"clock::Clock"},{"name":"report_context","type":"vector<vector<u8>>"},{"name":"report","type":"vector<u8>"}]},{"package":"ccip_offramp","module":"offramp","name":"initialize","parameters":[{"name":"state","type":"OffRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"fee_quoter_cap","type":"FeeQuoterCap"},{"name":"dest_transfer_cap","type":"osh::DestTransferCap"},{"name":"chain_selector","type":"u64"},{"name":"permissionless_execution_threshold_seconds","type":"u32"},{"name":"source_chains_selectors","type":"vector<u64>"},{"name":"source_chains_is_enabled","type":"vector<bool>"},{"name":"source_chains_is_rmn_verification_disabled","type":"vector<bool>"},{"name":"source_chains_on_ramp","type":"vector<vector<u8>>"}]},{"package":"ccip_offramp","module":"offramp","name":"manually_init_execute","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"clock","type":"clock::Clock"},{"name":"report_bytes","type":"vector<u8>"}]},{"package":"ccip_offramp","module":"offramp","name":"owner","parameters":[{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"pending_transfer_accepted","parameters":[{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"pending_transfer_from","parameters":[{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"pending_transfer_to","parameters":[{"name":"state","type":"OffRampState"}]},{"package":"ccip_offramp","module":"offramp","name":"remove_package_id","parameters":[{"name":"state","type":"OffRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"package_id","type":"address"}]},{"package":"ccip_offramp","module":"offramp","name":"set_dynamic_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"permissionless_execution_threshold_seconds","type":"u32"}]},{"package":"ccip_offramp","module":"offramp","name":"set_ocr3_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"config_digest","type":"vector<u8>"},{"name":"ocr_plugin_type","type":"u8"},{"name":"big_f","type":"u8"},{"name":"is_signature_verification_enabled","type":"bool"},{"name":"signers","type":"vector<vector<u8>>"},{"name":"transmitters","type":"vector<address>"}]},{"package":"ccip_offramp","module":"offramp","name":"transfer_ownership","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OffRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"new_owner","type":"address"}]},{"package":"ccip_offramp","module":"offramp","name":"type_and_version","parameters":null}]`

type IOfframp interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error)
	AddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error)
	RemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error)
	GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*models.SuiTransactionBlockResponse, error)
	FinishExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, receiverParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*models.SuiTransactionBlockResponse, error)
	GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*models.SuiTransactionBlockResponse, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, tokenReceiver string, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*models.SuiTransactionBlockResponse, error)
	SetOcr3Config(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*models.SuiTransactionBlockResponse, error)
	ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	Commit(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (*models.SuiTransactionBlockResponse, error)
	GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainConfig SourceChainConfig) (*models.SuiTransactionBlockResponse, error)
	GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error)
	SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, permissionlessExecutionThresholdSeconds uint32) (*models.SuiTransactionBlockResponse, error)
	ApplySourceChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsAddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsSetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsApplySourceChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsSetOcr3Config(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOfframpDevInspect
	Encoder() OfframpEncoder
	Bound() bind.IBoundContract
}

type IOfframpDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bind.Object, error)
	InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (bind.Object, error)
	ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (bind.Object, error)
	GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (byte, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, tokenReceiver string, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) ([]byte, error)
	ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([][]byte, error)
	ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]string, error)
	GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (uint64, error)
	GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, sourceChainSelector uint64) (SourceChainConfig, error)
	GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainConfig SourceChainConfig) ([]any, error)
	GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) ([]any, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (StaticConfig, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg StaticConfig) ([]any, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (DynamicConfig, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg DynamicConfig) ([]any, error)
	GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type OfframpEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(state bind.Object, ownerCap bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	AddPackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error)
	AddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	RemovePackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error)
	RemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOcr3Base(state bind.Object) (*bind.EncodedCall, error)
	GetOcr3BaseWithArgs(args ...any) (*bind.EncodedCall, error)
	InitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*bind.EncodedCall, error)
	InitExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	FinishExecute(ref bind.Object, state bind.Object, receiverParams bind.Object) (*bind.EncodedCall, error)
	FinishExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	ManuallyInitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*bind.EncodedCall, error)
	ManuallyInitExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	GetExecutionState(state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error)
	GetExecutionStateWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMetadataHash(ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*bind.EncodedCall, error)
	CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMessageHash(ref bind.Object, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, tokenReceiver string, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*bind.EncodedCall, error)
	CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error)
	SetOcr3Config(ref bind.Object, state bind.Object, ownerCap bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*bind.EncodedCall, error)
	SetOcr3ConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	ConfigSigners(state bind.Object) (*bind.EncodedCall, error)
	ConfigSignersWithArgs(args ...any) (*bind.EncodedCall, error)
	ConfigTransmitters(state bind.Object) (*bind.EncodedCall, error)
	ConfigTransmittersWithArgs(args ...any) (*bind.EncodedCall, error)
	Commit(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*bind.EncodedCall, error)
	CommitWithArgs(args ...any) (*bind.EncodedCall, error)
	GetMerkleRoot(state bind.Object, root []byte) (*bind.EncodedCall, error)
	GetMerkleRootWithArgs(args ...any) (*bind.EncodedCall, error)
	GetSourceChainConfig(ref bind.Object, state bind.Object, sourceChainSelector uint64) (*bind.EncodedCall, error)
	GetSourceChainConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetSourceChainConfigFields(ref bind.Object, sourceChainConfig SourceChainConfig) (*bind.EncodedCall, error)
	GetSourceChainConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllSourceChainConfigs(ref bind.Object, state bind.Object) (*bind.EncodedCall, error)
	GetAllSourceChainConfigsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfig(ref bind.Object, state bind.Object) (*bind.EncodedCall, error)
	GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfigFields(ref bind.Object, cfg StaticConfig) (*bind.EncodedCall, error)
	GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfig(ref bind.Object, state bind.Object) (*bind.EncodedCall, error)
	GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfigFields(ref bind.Object, cfg DynamicConfig) (*bind.EncodedCall, error)
	GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	SetDynamicConfig(ref bind.Object, state bind.Object, ownerCap bind.Object, permissionlessExecutionThresholdSeconds uint32) (*bind.EncodedCall, error)
	SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplySourceChainConfigUpdates(ref bind.Object, state bind.Object, ownerCap bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error)
	ApplySourceChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCcipPackageId() (*bind.EncodedCall, error)
	GetCcipPackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	Owner(state bind.Object) (*bind.EncodedCall, error)
	OwnerWithArgs(args ...any) (*bind.EncodedCall, error)
	HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error)
	HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error)
	TransferOwnership(ref bind.Object, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(ref bind.Object, state bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipFromObject(ref bind.Object, state bind.Object, from string) (*bind.EncodedCall, error)
	AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsAcceptOwnership(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ref bind.Object, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(ref bind.Object, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(ref bind.Object, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsAddPackageId(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRemovePackageId(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsSetDynamicConfig(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsSetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsApplySourceChainConfigUpdates(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsApplySourceChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsSetOcr3Config(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsSetOcr3ConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsTransferOwnership(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransfer(ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsAddAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAddAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRemoveAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRemoveAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error)
}

type OfframpContract struct {
	*bind.BoundContract
	offrampEncoder
	devInspect *OfframpDevInspect
}

type OfframpDevInspect struct {
	contract *OfframpContract
}

var _ IOfframp = (*OfframpContract)(nil)
var _ IOfframpDevInspect = (*OfframpDevInspect)(nil)

func NewOfframp(packageID string, client sui.ISuiAPI) (IOfframp, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_offramp", "offramp", client)
	if err != nil {
		return nil, err
	}

	c := &OfframpContract{
		BoundContract:  contract,
		offrampEncoder: offrampEncoder{BoundContract: contract},
	}
	c.devInspect = &OfframpDevInspect{contract: c}
	return c, nil
}

func (c *OfframpContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *OfframpContract) Encoder() OfframpEncoder {
	return c.offrampEncoder
}

func (c *OfframpContract) DevInspect() IOfframpDevInspect {
	return c.devInspect
}

type OffRampState struct {
	Id                                      string       `move:"sui::object::UID"`
	PackageIds                              []string     `move:"vector<address>"`
	Ocr3BaseState                           bind.Object  `move:"OCR3BaseState"`
	ChainSelector                           uint64       `move:"u64"`
	PermissionlessExecutionThresholdSeconds uint32       `move:"u32"`
	SourceChainConfigs                      bind.Object  `move:"VecMap<u64, SourceChainConfig>"`
	ExecutionStates                         bind.Object  `move:"Table<u64, Table<u64, u8>>"`
	Roots                                   bind.Object  `move:"Table<vector<u8>, u64>"`
	LatestPriceSequenceNumber               uint64       `move:"u64"`
	FeeQuoterCap                            *bind.Object `move:"0x1::option::Option<FeeQuoterCap>"`
	DestTransferCap                         *bind.Object `move:"0x1::option::Option<osh::DestTransferCap>"`
	OwnableState                            bind.Object  `move:"OwnableState"`
}

type OffRampObject struct {
	Id string `move:"sui::object::UID"`
}

type OffRampStatePointer struct {
	Id              string `move:"sui::object::UID"`
	OffRampObjectId string `move:"address"`
}

type SourceChainConfig struct {
	Router                    string `move:"address"`
	IsEnabled                 bool   `move:"bool"`
	MinSeqNr                  uint64 `move:"u64"`
	IsRmnVerificationDisabled bool   `move:"bool"`
	OnRamp                    []byte `move:"vector<u8>"`
}

type RampMessageHeader struct {
	MessageId           []byte `move:"vector<u8>"`
	SourceChainSelector uint64 `move:"u64"`
	DestChainSelector   uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
	Nonce               uint64 `move:"u64"`
}

type Any2SuiRampMessage struct {
	Header        RampMessageHeader      `move:"RampMessageHeader"`
	Sender        []byte                 `move:"vector<u8>"`
	Data          []byte                 `move:"vector<u8>"`
	Receiver      string                 `move:"address"`
	GasLimit      *big.Int               `move:"u256"`
	TokenReceiver string                 `move:"address"`
	TokenAmounts  []Any2SuiTokenTransfer `move:"vector<Any2SuiTokenTransfer>"`
}

type Any2SuiTokenTransfer struct {
	SourcePoolAddress []byte   `move:"vector<u8>"`
	DestTokenAddress  string   `move:"address"`
	DestGasAmount     uint32   `move:"u32"`
	ExtraData         []byte   `move:"vector<u8>"`
	Amount            *big.Int `move:"u256"`
}

type ExecutionReport struct {
	SourceChainSelector uint64             `move:"u64"`
	Message             Any2SuiRampMessage `move:"Any2SuiRampMessage"`
	OffchainTokenData   [][]byte           `move:"vector<vector<u8>>"`
	Proofs              [][]byte           `move:"vector<vector<u8>>"`
}

type CommitReport struct {
	PriceUpdates         PriceUpdates `move:"PriceUpdates"`
	BlessedMerkleRoots   []MerkleRoot `move:"vector<MerkleRoot>"`
	UnblessedMerkleRoots []MerkleRoot `move:"vector<MerkleRoot>"`
	RmnSignatures        [][]byte     `move:"vector<vector<u8>>"`
}

type PriceUpdates struct {
	TokenPriceUpdates []TokenPriceUpdate `move:"vector<TokenPriceUpdate>"`
	GasPriceUpdates   []GasPriceUpdate   `move:"vector<GasPriceUpdate>"`
}

type TokenPriceUpdate struct {
	SourceToken string   `move:"address"`
	UsdPerToken *big.Int `move:"u256"`
}

type GasPriceUpdate struct {
	DestChainSelector uint64   `move:"u64"`
	UsdPerUnitGas     *big.Int `move:"u256"`
}

type MerkleRoot struct {
	SourceChainSelector uint64 `move:"u64"`
	OnRampAddress       []byte `move:"vector<u8>"`
	MinSeqNr            uint64 `move:"u64"`
	MaxSeqNr            uint64 `move:"u64"`
	MerkleRoot          []byte `move:"vector<u8>"`
}

type StaticConfig struct {
	ChainSelector      uint64 `move:"u64"`
	RmnRemote          string `move:"address"`
	TokenAdminRegistry string `move:"address"`
	NonceManager       string `move:"address"`
}

type DynamicConfig struct {
	FeeQuoter                               string `move:"address"`
	PermissionlessExecutionThresholdSeconds uint32 `move:"u32"`
}

type StaticConfigSet struct {
	ChainSelector uint64 `move:"u64"`
}

type DynamicConfigSet struct {
	DynamicConfig DynamicConfig `move:"DynamicConfig"`
}

type SourceChainConfigSet struct {
	SourceChainSelector uint64            `move:"u64"`
	SourceChainConfig   SourceChainConfig `move:"SourceChainConfig"`
}

type SkippedAlreadyExecuted struct {
	SourceChainSelector uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
}

type ExecutionStateChanged struct {
	SourceChainSelector uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
	MessageId           []byte `move:"vector<u8>"`
	MessageHash         []byte `move:"vector<u8>"`
	State               byte   `move:"u8"`
}

type CommitReportAccepted struct {
	BlessedMerkleRoots   []MerkleRoot `move:"vector<MerkleRoot>"`
	UnblessedMerkleRoots []MerkleRoot `move:"vector<MerkleRoot>"`
	PriceUpdates         PriceUpdates `move:"PriceUpdates"`
}

type SkippedReportExecution struct {
	SourceChainSelector uint64 `move:"u64"`
}

type OFFRAMP struct {
}

type McmsCallback struct {
}

type McmsAcceptOwnershipProof struct {
}

type bcsOffRampState struct {
	Id                                      string
	PackageIds                              [][32]byte
	Ocr3BaseState                           bind.Object
	ChainSelector                           uint64
	PermissionlessExecutionThresholdSeconds uint32
	SourceChainConfigs                      bind.Object
	ExecutionStates                         bind.Object
	Roots                                   bind.Object
	LatestPriceSequenceNumber               uint64
	FeeQuoterCap                            *bind.Object
	DestTransferCap                         *bind.Object
	OwnableState                            bind.Object
}

func convertOffRampStateFromBCS(bcs bcsOffRampState) (OffRampState, error) {

	return OffRampState{
		Id: bcs.Id,
		PackageIds: func() []string {
			addrs := make([]string, len(bcs.PackageIds))
			for i, addr := range bcs.PackageIds {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		Ocr3BaseState:                           bcs.Ocr3BaseState,
		ChainSelector:                           bcs.ChainSelector,
		PermissionlessExecutionThresholdSeconds: bcs.PermissionlessExecutionThresholdSeconds,
		SourceChainConfigs:                      bcs.SourceChainConfigs,
		ExecutionStates:                         bcs.ExecutionStates,
		Roots:                                   bcs.Roots,
		LatestPriceSequenceNumber:               bcs.LatestPriceSequenceNumber,
		FeeQuoterCap:                            bcs.FeeQuoterCap,
		DestTransferCap:                         bcs.DestTransferCap,
		OwnableState:                            bcs.OwnableState,
	}, nil
}

type bcsOffRampStatePointer struct {
	Id              string
	OffRampObjectId [32]byte
}

func convertOffRampStatePointerFromBCS(bcs bcsOffRampStatePointer) (OffRampStatePointer, error) {

	return OffRampStatePointer{
		Id:              bcs.Id,
		OffRampObjectId: fmt.Sprintf("0x%x", bcs.OffRampObjectId),
	}, nil
}

type bcsSourceChainConfig struct {
	Router                    [32]byte
	IsEnabled                 bool
	MinSeqNr                  uint64
	IsRmnVerificationDisabled bool
	OnRamp                    []byte
}

func convertSourceChainConfigFromBCS(bcs bcsSourceChainConfig) (SourceChainConfig, error) {

	return SourceChainConfig{
		Router:                    fmt.Sprintf("0x%x", bcs.Router),
		IsEnabled:                 bcs.IsEnabled,
		MinSeqNr:                  bcs.MinSeqNr,
		IsRmnVerificationDisabled: bcs.IsRmnVerificationDisabled,
		OnRamp:                    bcs.OnRamp,
	}, nil
}

type bcsAny2SuiRampMessage struct {
	Header        RampMessageHeader
	Sender        []byte
	Data          []byte
	Receiver      [32]byte
	GasLimit      [32]byte
	TokenReceiver [32]byte
	TokenAmounts  []Any2SuiTokenTransfer
}

func convertAny2SuiRampMessageFromBCS(bcs bcsAny2SuiRampMessage) (Any2SuiRampMessage, error) {
	GasLimitField, err := bind.DecodeU256Value(bcs.GasLimit)
	if err != nil {
		return Any2SuiRampMessage{}, fmt.Errorf("failed to decode u256 field GasLimit: %w", err)
	}

	return Any2SuiRampMessage{
		Header:        bcs.Header,
		Sender:        bcs.Sender,
		Data:          bcs.Data,
		Receiver:      fmt.Sprintf("0x%x", bcs.Receiver),
		GasLimit:      GasLimitField,
		TokenReceiver: fmt.Sprintf("0x%x", bcs.TokenReceiver),
		TokenAmounts:  bcs.TokenAmounts,
	}, nil
}

type bcsAny2SuiTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  [32]byte
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            [32]byte
}

func convertAny2SuiTokenTransferFromBCS(bcs bcsAny2SuiTokenTransfer) (Any2SuiTokenTransfer, error) {
	AmountField, err := bind.DecodeU256Value(bcs.Amount)
	if err != nil {
		return Any2SuiTokenTransfer{}, fmt.Errorf("failed to decode u256 field Amount: %w", err)
	}

	return Any2SuiTokenTransfer{
		SourcePoolAddress: bcs.SourcePoolAddress,
		DestTokenAddress:  fmt.Sprintf("0x%x", bcs.DestTokenAddress),
		DestGasAmount:     bcs.DestGasAmount,
		ExtraData:         bcs.ExtraData,
		Amount:            AmountField,
	}, nil
}

type bcsExecutionReport struct {
	SourceChainSelector uint64
	Message             bcsAny2SuiRampMessage
	OffchainTokenData   [][]byte
	Proofs              [][]byte
}

func convertExecutionReportFromBCS(bcs bcsExecutionReport) (ExecutionReport, error) {
	MessageField, err := convertAny2SuiRampMessageFromBCS(bcs.Message)
	if err != nil {
		return ExecutionReport{}, fmt.Errorf("failed to convert nested struct Message: %w", err)
	}

	return ExecutionReport{
		SourceChainSelector: bcs.SourceChainSelector,
		Message:             MessageField,
		OffchainTokenData:   bcs.OffchainTokenData,
		Proofs:              bcs.Proofs,
	}, nil
}

type bcsTokenPriceUpdate struct {
	SourceToken [32]byte
	UsdPerToken [32]byte
}

func convertTokenPriceUpdateFromBCS(bcs bcsTokenPriceUpdate) (TokenPriceUpdate, error) {
	UsdPerTokenField, err := bind.DecodeU256Value(bcs.UsdPerToken)
	if err != nil {
		return TokenPriceUpdate{}, fmt.Errorf("failed to decode u256 field UsdPerToken: %w", err)
	}

	return TokenPriceUpdate{
		SourceToken: fmt.Sprintf("0x%x", bcs.SourceToken),
		UsdPerToken: UsdPerTokenField,
	}, nil
}

type bcsGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     [32]byte
}

func convertGasPriceUpdateFromBCS(bcs bcsGasPriceUpdate) (GasPriceUpdate, error) {
	UsdPerUnitGasField, err := bind.DecodeU256Value(bcs.UsdPerUnitGas)
	if err != nil {
		return GasPriceUpdate{}, fmt.Errorf("failed to decode u256 field UsdPerUnitGas: %w", err)
	}

	return GasPriceUpdate{
		DestChainSelector: bcs.DestChainSelector,
		UsdPerUnitGas:     UsdPerUnitGasField,
	}, nil
}

type bcsStaticConfig struct {
	ChainSelector      uint64
	RmnRemote          [32]byte
	TokenAdminRegistry [32]byte
	NonceManager       [32]byte
}

func convertStaticConfigFromBCS(bcs bcsStaticConfig) (StaticConfig, error) {

	return StaticConfig{
		ChainSelector:      bcs.ChainSelector,
		RmnRemote:          fmt.Sprintf("0x%x", bcs.RmnRemote),
		TokenAdminRegistry: fmt.Sprintf("0x%x", bcs.TokenAdminRegistry),
		NonceManager:       fmt.Sprintf("0x%x", bcs.NonceManager),
	}, nil
}

type bcsDynamicConfig struct {
	FeeQuoter                               [32]byte
	PermissionlessExecutionThresholdSeconds uint32
}

func convertDynamicConfigFromBCS(bcs bcsDynamicConfig) (DynamicConfig, error) {

	return DynamicConfig{
		FeeQuoter:                               fmt.Sprintf("0x%x", bcs.FeeQuoter),
		PermissionlessExecutionThresholdSeconds: bcs.PermissionlessExecutionThresholdSeconds,
	}, nil
}

type bcsDynamicConfigSet struct {
	DynamicConfig bcsDynamicConfig
}

func convertDynamicConfigSetFromBCS(bcs bcsDynamicConfigSet) (DynamicConfigSet, error) {
	DynamicConfigField, err := convertDynamicConfigFromBCS(bcs.DynamicConfig)
	if err != nil {
		return DynamicConfigSet{}, fmt.Errorf("failed to convert nested struct DynamicConfig: %w", err)
	}

	return DynamicConfigSet{
		DynamicConfig: DynamicConfigField,
	}, nil
}

type bcsSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceChainConfig   bcsSourceChainConfig
}

func convertSourceChainConfigSetFromBCS(bcs bcsSourceChainConfigSet) (SourceChainConfigSet, error) {
	SourceChainConfigField, err := convertSourceChainConfigFromBCS(bcs.SourceChainConfig)
	if err != nil {
		return SourceChainConfigSet{}, fmt.Errorf("failed to convert nested struct SourceChainConfig: %w", err)
	}

	return SourceChainConfigSet{
		SourceChainSelector: bcs.SourceChainSelector,
		SourceChainConfig:   SourceChainConfigField,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip_offramp::offramp::OffRampState", func(data []byte) (interface{}, error) {
		var temp bcsOffRampState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOffRampStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OffRampState
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::OffRampState>", func(data []byte) (interface{}, error) {
		var temps []bcsOffRampState
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OffRampState, len(temps))
		for i, temp := range temps {
			result, err := convertOffRampStateFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::OffRampObject", func(data []byte) (interface{}, error) {
		var result OffRampObject
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OffRampObject
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::OffRampObject>", func(data []byte) (interface{}, error) {
		var results []OffRampObject
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::OffRampStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsOffRampStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOffRampStatePointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OffRampStatePointer
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::OffRampStatePointer>", func(data []byte) (interface{}, error) {
		var temps []bcsOffRampStatePointer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OffRampStatePointer, len(temps))
		for i, temp := range temps {
			result, err := convertOffRampStatePointerFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SourceChainConfig", func(data []byte) (interface{}, error) {
		var temp bcsSourceChainConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSourceChainConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for SourceChainConfig
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::SourceChainConfig>", func(data []byte) (interface{}, error) {
		var temps []bcsSourceChainConfig
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]SourceChainConfig, len(temps))
		for i, temp := range temps {
			result, err := convertSourceChainConfigFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::RampMessageHeader", func(data []byte) (interface{}, error) {
		var result RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for RampMessageHeader
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::RampMessageHeader>", func(data []byte) (interface{}, error) {
		var results []RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::Any2SuiRampMessage", func(data []byte) (interface{}, error) {
		var temp bcsAny2SuiRampMessage
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAny2SuiRampMessageFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for Any2SuiRampMessage
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::Any2SuiRampMessage>", func(data []byte) (interface{}, error) {
		var temps []bcsAny2SuiRampMessage
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]Any2SuiRampMessage, len(temps))
		for i, temp := range temps {
			result, err := convertAny2SuiRampMessageFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::Any2SuiTokenTransfer", func(data []byte) (interface{}, error) {
		var temp bcsAny2SuiTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAny2SuiTokenTransferFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for Any2SuiTokenTransfer
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::Any2SuiTokenTransfer>", func(data []byte) (interface{}, error) {
		var temps []bcsAny2SuiTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]Any2SuiTokenTransfer, len(temps))
		for i, temp := range temps {
			result, err := convertAny2SuiTokenTransferFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::ExecutionReport", func(data []byte) (interface{}, error) {
		var temp bcsExecutionReport
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertExecutionReportFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ExecutionReport
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::ExecutionReport>", func(data []byte) (interface{}, error) {
		var temps []bcsExecutionReport
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]ExecutionReport, len(temps))
		for i, temp := range temps {
			result, err := convertExecutionReportFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::CommitReport", func(data []byte) (interface{}, error) {
		var result CommitReport
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for CommitReport
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::CommitReport>", func(data []byte) (interface{}, error) {
		var results []CommitReport
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::PriceUpdates", func(data []byte) (interface{}, error) {
		var result PriceUpdates
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for PriceUpdates
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::PriceUpdates>", func(data []byte) (interface{}, error) {
		var results []PriceUpdates
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::TokenPriceUpdate", func(data []byte) (interface{}, error) {
		var temp bcsTokenPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenPriceUpdateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TokenPriceUpdate
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::TokenPriceUpdate>", func(data []byte) (interface{}, error) {
		var temps []bcsTokenPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]TokenPriceUpdate, len(temps))
		for i, temp := range temps {
			result, err := convertTokenPriceUpdateFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::GasPriceUpdate", func(data []byte) (interface{}, error) {
		var temp bcsGasPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertGasPriceUpdateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for GasPriceUpdate
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::GasPriceUpdate>", func(data []byte) (interface{}, error) {
		var temps []bcsGasPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]GasPriceUpdate, len(temps))
		for i, temp := range temps {
			result, err := convertGasPriceUpdateFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::MerkleRoot", func(data []byte) (interface{}, error) {
		var result MerkleRoot
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for MerkleRoot
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::MerkleRoot>", func(data []byte) (interface{}, error) {
		var results []MerkleRoot
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::StaticConfig", func(data []byte) (interface{}, error) {
		var temp bcsStaticConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertStaticConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for StaticConfig
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::StaticConfig>", func(data []byte) (interface{}, error) {
		var temps []bcsStaticConfig
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]StaticConfig, len(temps))
		for i, temp := range temps {
			result, err := convertStaticConfigFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::DynamicConfig", func(data []byte) (interface{}, error) {
		var temp bcsDynamicConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDynamicConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DynamicConfig
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::DynamicConfig>", func(data []byte) (interface{}, error) {
		var temps []bcsDynamicConfig
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]DynamicConfig, len(temps))
		for i, temp := range temps {
			result, err := convertDynamicConfigFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::StaticConfigSet", func(data []byte) (interface{}, error) {
		var result StaticConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for StaticConfigSet
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::StaticConfigSet>", func(data []byte) (interface{}, error) {
		var results []StaticConfigSet
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::DynamicConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsDynamicConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDynamicConfigSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DynamicConfigSet
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::DynamicConfigSet>", func(data []byte) (interface{}, error) {
		var temps []bcsDynamicConfigSet
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]DynamicConfigSet, len(temps))
		for i, temp := range temps {
			result, err := convertDynamicConfigSetFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SourceChainConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsSourceChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSourceChainConfigSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for SourceChainConfigSet
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::SourceChainConfigSet>", func(data []byte) (interface{}, error) {
		var temps []bcsSourceChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]SourceChainConfigSet, len(temps))
		for i, temp := range temps {
			result, err := convertSourceChainConfigSetFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SkippedAlreadyExecuted", func(data []byte) (interface{}, error) {
		var result SkippedAlreadyExecuted
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for SkippedAlreadyExecuted
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::SkippedAlreadyExecuted>", func(data []byte) (interface{}, error) {
		var results []SkippedAlreadyExecuted
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::ExecutionStateChanged", func(data []byte) (interface{}, error) {
		var result ExecutionStateChanged
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ExecutionStateChanged
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::ExecutionStateChanged>", func(data []byte) (interface{}, error) {
		var results []ExecutionStateChanged
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::CommitReportAccepted", func(data []byte) (interface{}, error) {
		var result CommitReportAccepted
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for CommitReportAccepted
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::CommitReportAccepted>", func(data []byte) (interface{}, error) {
		var results []CommitReportAccepted
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SkippedReportExecution", func(data []byte) (interface{}, error) {
		var result SkippedReportExecution
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for SkippedReportExecution
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::SkippedReportExecution>", func(data []byte) (interface{}, error) {
		var results []SkippedReportExecution
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::OFFRAMP", func(data []byte) (interface{}, error) {
		var result OFFRAMP
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OFFRAMP
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::OFFRAMP>", func(data []byte) (interface{}, error) {
		var results []OFFRAMP
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsCallback
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::McmsCallback>", func(data []byte) (interface{}, error) {
		var results []McmsCallback
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::McmsAcceptOwnershipProof", func(data []byte) (interface{}, error) {
		var result McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsAcceptOwnershipProof
	bind.RegisterStructDecoder("vector<ccip_offramp::offramp::McmsAcceptOwnershipProof>", func(data []byte) (interface{}, error) {
		var results []McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *OfframpContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *OfframpContract) Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.Initialize(state, ownerCap, feeQuoterCap, destTransferCap, chainSelector, permissionlessExecutionThresholdSeconds, sourceChainsSelectors, sourceChainsIsEnabled, sourceChainsIsRmnVerificationDisabled, sourceChainsOnRamp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddPackageId executes the add_package_id Move function.
func (c *OfframpContract) AddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.AddPackageId(state, ownerCap, packageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemovePackageId executes the remove_package_id Move function.
func (c *OfframpContract) RemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.RemovePackageId(state, ownerCap, packageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOcr3Base executes the get_ocr3_base Move function.
func (c *OfframpContract) GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetOcr3Base(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// InitExecute executes the init_execute Move function.
func (c *OfframpContract) InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.InitExecute(ref, state, clock, reportContext, report)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// FinishExecute executes the finish_execute Move function.
func (c *OfframpContract) FinishExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, receiverParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.FinishExecute(ref, state, receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ManuallyInitExecute executes the manually_init_execute Move function.
func (c *OfframpContract) ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ManuallyInitExecute(ref, state, clock, reportBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetExecutionState executes the get_execution_state Move function.
func (c *OfframpContract) GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetExecutionState(state, sourceChainSelector, sequenceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMetadataHash executes the calculate_metadata_hash Move function.
func (c *OfframpContract) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.CalculateMetadataHash(ref, sourceChainSelector, destChainSelector, onRamp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMessageHash executes the calculate_message_hash Move function.
func (c *OfframpContract) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, tokenReceiver string, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.CalculateMessageHash(ref, messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, onRamp, data, gasLimit, tokenReceiver, sourcePoolAddresses, destTokenAddresses, destGasAmounts, extraDatas, amounts)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetOcr3Config executes the set_ocr3_config Move function.
func (c *OfframpContract) SetOcr3Config(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.SetOcr3Config(ref, state, ownerCap, configDigest, ocrPluginType, bigF, isSignatureVerificationEnabled, signers, transmitters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigSigners executes the config_signers Move function.
func (c *OfframpContract) ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ConfigSigners(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigTransmitters executes the config_transmitters Move function.
func (c *OfframpContract) ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ConfigTransmitters(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Commit executes the commit Move function.
func (c *OfframpContract) Commit(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.Commit(ref, state, clock, reportContext, report, signatures)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetMerkleRoot executes the get_merkle_root Move function.
func (c *OfframpContract) GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetMerkleRoot(state, root)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSourceChainConfig executes the get_source_chain_config Move function.
func (c *OfframpContract) GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetSourceChainConfig(ref, state, sourceChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSourceChainConfigFields executes the get_source_chain_config_fields Move function.
func (c *OfframpContract) GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainConfig SourceChainConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetSourceChainConfigFields(ref, sourceChainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllSourceChainConfigs executes the get_all_source_chain_configs Move function.
func (c *OfframpContract) GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetAllSourceChainConfigs(ref, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function.
func (c *OfframpContract) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetStaticConfig(ref, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfigFields executes the get_static_config_fields Move function.
func (c *OfframpContract) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetStaticConfigFields(ref, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfig executes the get_dynamic_config Move function.
func (c *OfframpContract) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetDynamicConfig(ref, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfigFields executes the get_dynamic_config_fields Move function.
func (c *OfframpContract) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetDynamicConfigFields(ref, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetDynamicConfig executes the set_dynamic_config Move function.
func (c *OfframpContract) SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, permissionlessExecutionThresholdSeconds uint32) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.SetDynamicConfig(ref, state, ownerCap, permissionlessExecutionThresholdSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplySourceChainConfigUpdates executes the apply_source_chain_config_updates Move function.
func (c *OfframpContract) ApplySourceChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ApplySourceChainConfigUpdates(ref, state, ownerCap, sourceChainsSelector, sourceChainsIsEnabled, sourceChainsIsRmnVerificationDisabled, sourceChainsOnRamp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCcipPackageId executes the get_ccip_package_id Move function.
func (c *OfframpContract) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetCcipPackageId()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *OfframpContract) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.Owner(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *OfframpContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.HasPendingTransfer(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *OfframpContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *OfframpContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *OfframpContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *OfframpContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.TransferOwnership(ref, state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *OfframpContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.AcceptOwnership(ref, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *OfframpContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.AcceptOwnershipFromObject(ref, state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAcceptOwnership executes the mcms_accept_ownership Move function.
func (c *OfframpContract) McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsAcceptOwnership(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *OfframpContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ExecuteOwnershipTransfer(ref, ownerCap, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *OfframpContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ExecuteOwnershipTransferToMcms(ref, ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *OfframpContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsRegisterUpgradeCap(ref, upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddPackageId executes the mcms_add_package_id Move function.
func (c *OfframpContract) McmsAddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsAddPackageId(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemovePackageId executes the mcms_remove_package_id Move function.
func (c *OfframpContract) McmsRemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsRemovePackageId(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetDynamicConfig executes the mcms_set_dynamic_config Move function.
func (c *OfframpContract) McmsSetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsSetDynamicConfig(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplySourceChainConfigUpdates executes the mcms_apply_source_chain_config_updates Move function.
func (c *OfframpContract) McmsApplySourceChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsApplySourceChainConfigUpdates(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetOcr3Config executes the mcms_set_ocr3_config Move function.
func (c *OfframpContract) McmsSetOcr3Config(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsSetOcr3Config(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsTransferOwnership executes the mcms_transfer_ownership Move function.
func (c *OfframpContract) McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsTransferOwnership(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsExecuteOwnershipTransfer executes the mcms_execute_ownership_transfer Move function.
func (c *OfframpContract) McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsExecuteOwnershipTransfer(ref, state, registry, deployerState, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddAllowedModules executes the mcms_add_allowed_modules Move function.
func (c *OfframpContract) McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsAddAllowedModules(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemoveAllowedModules executes the mcms_remove_allowed_modules Move function.
func (c *OfframpContract) McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsRemoveAllowedModules(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *OfframpDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.offrampEncoder.TypeAndVersion()
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

// GetOcr3Base executes the get_ocr3_base Move function using DevInspect to get return values.
//
// Returns: &OCR3BaseState
func (d *OfframpDevInspect) GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bind.Object, error) {
	encoded, err := d.contract.offrampEncoder.GetOcr3Base(state)
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

// InitExecute executes the init_execute Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *OfframpDevInspect) InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (bind.Object, error) {
	encoded, err := d.contract.offrampEncoder.InitExecute(ref, state, clock, reportContext, report)
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

// ManuallyInitExecute executes the manually_init_execute Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *OfframpDevInspect) ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (bind.Object, error) {
	encoded, err := d.contract.offrampEncoder.ManuallyInitExecute(ref, state, clock, reportBytes)
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

// GetExecutionState executes the get_execution_state Move function using DevInspect to get return values.
//
// Returns: u8
func (d *OfframpDevInspect) GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (byte, error) {
	encoded, err := d.contract.offrampEncoder.GetExecutionState(state, sourceChainSelector, sequenceNumber)
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

// CalculateMetadataHash executes the calculate_metadata_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OfframpDevInspect) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error) {
	encoded, err := d.contract.offrampEncoder.CalculateMetadataHash(ref, sourceChainSelector, destChainSelector, onRamp)
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

// CalculateMessageHash executes the calculate_message_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OfframpDevInspect) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, tokenReceiver string, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) ([]byte, error) {
	encoded, err := d.contract.offrampEncoder.CalculateMessageHash(ref, messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, onRamp, data, gasLimit, tokenReceiver, sourcePoolAddresses, destTokenAddresses, destGasAmounts, extraDatas, amounts)
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

// ConfigSigners executes the config_signers Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *OfframpDevInspect) ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([][]byte, error) {
	encoded, err := d.contract.offrampEncoder.ConfigSigners(state)
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
	result, ok := results[0].([][]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected [][]byte, got %T", results[0])
	}
	return result, nil
}

// ConfigTransmitters executes the config_transmitters Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *OfframpDevInspect) ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]string, error) {
	encoded, err := d.contract.offrampEncoder.ConfigTransmitters(state)
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
	result, ok := results[0].([]string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []string, got %T", results[0])
	}
	return result, nil
}

// GetMerkleRoot executes the get_merkle_root Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OfframpDevInspect) GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (uint64, error) {
	encoded, err := d.contract.offrampEncoder.GetMerkleRoot(state, root)
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

// GetSourceChainConfig executes the get_source_chain_config Move function using DevInspect to get return values.
//
// Returns: SourceChainConfig
func (d *OfframpDevInspect) GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, sourceChainSelector uint64) (SourceChainConfig, error) {
	encoded, err := d.contract.offrampEncoder.GetSourceChainConfig(ref, state, sourceChainSelector)
	if err != nil {
		return SourceChainConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return SourceChainConfig{}, err
	}
	if len(results) == 0 {
		return SourceChainConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(SourceChainConfig)
	if !ok {
		return SourceChainConfig{}, fmt.Errorf("unexpected return type: expected SourceChainConfig, got %T", results[0])
	}
	return result, nil
}

// GetSourceChainConfigFields executes the get_source_chain_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: bool
//	[2]: u64
//	[3]: bool
//	[4]: vector<u8>
func (d *OfframpDevInspect) GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainConfig SourceChainConfig) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetSourceChainConfigFields(ref, sourceChainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetAllSourceChainConfigs executes the get_all_source_chain_configs Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<u64>
//	[1]: vector<SourceChainConfig>
func (d *OfframpDevInspect) GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetAllSourceChainConfigs(ref, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function using DevInspect to get return values.
//
// Returns: StaticConfig
func (d *OfframpDevInspect) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (StaticConfig, error) {
	encoded, err := d.contract.offrampEncoder.GetStaticConfig(ref, state)
	if err != nil {
		return StaticConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return StaticConfig{}, err
	}
	if len(results) == 0 {
		return StaticConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(StaticConfig)
	if !ok {
		return StaticConfig{}, fmt.Errorf("unexpected return type: expected StaticConfig, got %T", results[0])
	}
	return result, nil
}

// GetStaticConfigFields executes the get_static_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u64
//	[1]: address
//	[2]: address
//	[3]: address
func (d *OfframpDevInspect) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg StaticConfig) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetStaticConfigFields(ref, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetDynamicConfig executes the get_dynamic_config Move function using DevInspect to get return values.
//
// Returns: DynamicConfig
func (d *OfframpDevInspect) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (DynamicConfig, error) {
	encoded, err := d.contract.offrampEncoder.GetDynamicConfig(ref, state)
	if err != nil {
		return DynamicConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return DynamicConfig{}, err
	}
	if len(results) == 0 {
		return DynamicConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(DynamicConfig)
	if !ok {
		return DynamicConfig{}, fmt.Errorf("unexpected return type: expected DynamicConfig, got %T", results[0])
	}
	return result, nil
}

// GetDynamicConfigFields executes the get_dynamic_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: u32
func (d *OfframpDevInspect) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, ref bind.Object, cfg DynamicConfig) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetDynamicConfigFields(ref, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetCcipPackageId executes the get_ccip_package_id Move function using DevInspect to get return values.
//
// Returns: address
func (d *OfframpDevInspect) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.offrampEncoder.GetCcipPackageId()
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

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *OfframpDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.offrampEncoder.Owner(state)
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

// HasPendingTransfer executes the has_pending_transfer Move function using DevInspect to get return values.
//
// Returns: bool
func (d *OfframpDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error) {
	encoded, err := d.contract.offrampEncoder.HasPendingTransfer(state)
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

// PendingTransferFrom executes the pending_transfer_from Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *OfframpDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.offrampEncoder.PendingTransferFrom(state)
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
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferTo executes the pending_transfer_to Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *OfframpDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.offrampEncoder.PendingTransferTo(state)
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
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<bool>
func (d *OfframpDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error) {
	encoded, err := d.contract.offrampEncoder.PendingTransferAccepted(state)
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
	result, ok := results[0].(*bool)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *bool, got %T", results[0])
	}
	return result, nil
}

type offrampEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c offrampEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// Initialize encodes a call to the initialize Move function.
func (c offrampEncoder) Initialize(state bind.Object, ownerCap bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"FeeQuoterCap",
		"osh::DestTransferCap",
		"u64",
		"u32",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}, []any{
		state,
		ownerCap,
		feeQuoterCap,
		destTransferCap,
		chainSelector,
		permissionlessExecutionThresholdSeconds,
		sourceChainsSelectors,
		sourceChainsIsEnabled,
		sourceChainsIsRmnVerificationDisabled,
		sourceChainsOnRamp,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"FeeQuoterCap",
		"osh::DestTransferCap",
		"u64",
		"u32",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AddPackageId encodes a call to the add_package_id Move function.
func (c offrampEncoder) AddPackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_package_id", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}, []any{
		state,
		ownerCap,
		packageId,
	}, nil)
}

// AddPackageIdWithArgs encodes a call to the add_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) AddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_package_id", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RemovePackageId encodes a call to the remove_package_id Move function.
func (c offrampEncoder) RemovePackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_package_id", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}, []any{
		state,
		ownerCap,
		packageId,
	}, nil)
}

// RemovePackageIdWithArgs encodes a call to the remove_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) RemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_package_id", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetOcr3Base encodes a call to the get_ocr3_base Move function.
func (c offrampEncoder) GetOcr3Base(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ocr3_base", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"&OCR3BaseState",
	})
}

// GetOcr3BaseWithArgs encodes a call to the get_ocr3_base Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetOcr3BaseWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ocr3_base", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&OCR3BaseState",
	})
}

// InitExecute encodes a call to the init_execute Move function.
func (c offrampEncoder) InitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("init_execute", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
	}, []any{
		ref,
		state,
		clock,
		reportContext,
		report,
	}, []string{
		"osh::ReceiverParams",
	})
}

// InitExecuteWithArgs encodes a call to the init_execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) InitExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("init_execute", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}

// FinishExecute encodes a call to the finish_execute Move function.
func (c offrampEncoder) FinishExecute(ref bind.Object, state bind.Object, receiverParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("finish_execute", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"osh::ReceiverParams",
	}, []any{
		ref,
		state,
		receiverParams,
	}, nil)
}

// FinishExecuteWithArgs encodes a call to the finish_execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) FinishExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"osh::ReceiverParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("finish_execute", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ManuallyInitExecute encodes a call to the manually_init_execute Move function.
func (c offrampEncoder) ManuallyInitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("manually_init_execute", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<u8>",
	}, []any{
		ref,
		state,
		clock,
		reportBytes,
	}, []string{
		"osh::ReceiverParams",
	})
}

// ManuallyInitExecuteWithArgs encodes a call to the manually_init_execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ManuallyInitExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("manually_init_execute", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}

// GetExecutionState encodes a call to the get_execution_state Move function.
func (c offrampEncoder) GetExecutionState(state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_execution_state", typeArgsList, typeParamsList, []string{
		"&OffRampState",
		"u64",
		"u64",
	}, []any{
		state,
		sourceChainSelector,
		sequenceNumber,
	}, []string{
		"u8",
	})
}

// GetExecutionStateWithArgs encodes a call to the get_execution_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetExecutionStateWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_execution_state", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// CalculateMetadataHash encodes a call to the calculate_metadata_hash Move function.
func (c offrampEncoder) CalculateMetadataHash(ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_metadata_hash", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"u64",
		"vector<u8>",
	}, []any{
		ref,
		sourceChainSelector,
		destChainSelector,
		onRamp,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMetadataHashWithArgs encodes a call to the calculate_metadata_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_metadata_hash", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// CalculateMessageHash encodes a call to the calculate_message_hash Move function.
func (c offrampEncoder) CalculateMessageHash(ref bind.Object, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, tokenReceiver string, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_message_hash", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"vector<u8>",
		"address",
		"vector<u8>",
		"vector<u8>",
		"u256",
		"address",
		"vector<vector<u8>>",
		"vector<address>",
		"vector<u32>",
		"vector<vector<u8>>",
		"vector<u256>",
	}, []any{
		ref,
		messageId,
		sourceChainSelector,
		destChainSelector,
		sequenceNumber,
		nonce,
		sender,
		receiver,
		onRamp,
		data,
		gasLimit,
		tokenReceiver,
		sourcePoolAddresses,
		destTokenAddresses,
		destGasAmounts,
		extraDatas,
		amounts,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMessageHashWithArgs encodes a call to the calculate_message_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"vector<u8>",
		"address",
		"vector<u8>",
		"vector<u8>",
		"u256",
		"address",
		"vector<vector<u8>>",
		"vector<address>",
		"vector<u32>",
		"vector<vector<u8>>",
		"vector<u256>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_message_hash", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// SetOcr3Config encodes a call to the set_ocr3_config Move function.
func (c offrampEncoder) SetOcr3Config(ref bind.Object, state bind.Object, ownerCap bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_ocr3_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u8>",
		"u8",
		"u8",
		"bool",
		"vector<vector<u8>>",
		"vector<address>",
	}, []any{
		ref,
		state,
		ownerCap,
		configDigest,
		ocrPluginType,
		bigF,
		isSignatureVerificationEnabled,
		signers,
		transmitters,
	}, nil)
}

// SetOcr3ConfigWithArgs encodes a call to the set_ocr3_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) SetOcr3ConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u8>",
		"u8",
		"u8",
		"bool",
		"vector<vector<u8>>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_ocr3_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ConfigSigners encodes a call to the config_signers Move function.
func (c offrampEncoder) ConfigSigners(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_signers", typeArgsList, typeParamsList, []string{
		"&OCRConfig",
	}, []any{
		state,
	}, []string{
		"vector<vector<u8>>",
	})
}

// ConfigSignersWithArgs encodes a call to the config_signers Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ConfigSignersWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OCRConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_signers", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
	})
}

// ConfigTransmitters encodes a call to the config_transmitters Move function.
func (c offrampEncoder) ConfigTransmitters(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_transmitters", typeArgsList, typeParamsList, []string{
		"&OCRConfig",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// ConfigTransmittersWithArgs encodes a call to the config_transmitters Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ConfigTransmittersWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OCRConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_transmitters", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// Commit encodes a call to the commit Move function.
func (c offrampEncoder) Commit(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("commit", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<vector<u8>>",
	}, []any{
		ref,
		state,
		clock,
		reportContext,
		report,
		signatures,
	}, nil)
}

// CommitWithArgs encodes a call to the commit Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) CommitWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("commit", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetMerkleRoot encodes a call to the get_merkle_root Move function.
func (c offrampEncoder) GetMerkleRoot(state bind.Object, root []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_merkle_root", typeArgsList, typeParamsList, []string{
		"&OffRampState",
		"vector<u8>",
	}, []any{
		state,
		root,
	}, []string{
		"u64",
	})
}

// GetMerkleRootWithArgs encodes a call to the get_merkle_root Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetMerkleRootWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_merkle_root", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetSourceChainConfig encodes a call to the get_source_chain_config Move function.
func (c offrampEncoder) GetSourceChainConfig(ref bind.Object, state bind.Object, sourceChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&OffRampState",
		"u64",
	}, []any{
		ref,
		state,
		sourceChainSelector,
	}, []string{
		"ccip_offramp::offramp::SourceChainConfig",
	})
}

// GetSourceChainConfigWithArgs encodes a call to the get_source_chain_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetSourceChainConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&OffRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_offramp::offramp::SourceChainConfig",
	})
}

// GetSourceChainConfigFields encodes a call to the get_source_chain_config_fields Move function.
func (c offrampEncoder) GetSourceChainConfigFields(ref bind.Object, sourceChainConfig SourceChainConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config_fields", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"ccip_offramp::offramp::SourceChainConfig",
	}, []any{
		ref,
		sourceChainConfig,
	}, []string{
		"address",
		"bool",
		"u64",
		"bool",
		"vector<u8>",
	})
}

// GetSourceChainConfigFieldsWithArgs encodes a call to the get_source_chain_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetSourceChainConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"ccip_offramp::offramp::SourceChainConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"bool",
		"u64",
		"bool",
		"vector<u8>",
	})
}

// GetAllSourceChainConfigs encodes a call to the get_all_source_chain_configs Move function.
func (c offrampEncoder) GetAllSourceChainConfigs(ref bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_source_chain_configs", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&OffRampState",
	}, []any{
		ref,
		state,
	}, []string{
		"vector<u64>",
		"vector<ccip_offramp::offramp::SourceChainConfig>",
	})
}

// GetAllSourceChainConfigsWithArgs encodes a call to the get_all_source_chain_configs Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetAllSourceChainConfigsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_source_chain_configs", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
		"vector<ccip_offramp::offramp::SourceChainConfig>",
	})
}

// GetStaticConfig encodes a call to the get_static_config Move function.
func (c offrampEncoder) GetStaticConfig(ref bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&OffRampState",
	}, []any{
		ref,
		state,
	}, []string{
		"ccip_offramp::offramp::StaticConfig",
	})
}

// GetStaticConfigWithArgs encodes a call to the get_static_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_offramp::offramp::StaticConfig",
	})
}

// GetStaticConfigFields encodes a call to the get_static_config_fields Move function.
func (c offrampEncoder) GetStaticConfigFields(ref bind.Object, cfg StaticConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"ccip_offramp::offramp::StaticConfig",
	}, []any{
		ref,
		cfg,
	}, []string{
		"u64",
		"address",
		"address",
		"address",
	})
}

// GetStaticConfigFieldsWithArgs encodes a call to the get_static_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"ccip_offramp::offramp::StaticConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
		"address",
		"address",
		"address",
	})
}

// GetDynamicConfig encodes a call to the get_dynamic_config Move function.
func (c offrampEncoder) GetDynamicConfig(ref bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&OffRampState",
	}, []any{
		ref,
		state,
	}, []string{
		"ccip_offramp::offramp::DynamicConfig",
	})
}

// GetDynamicConfigWithArgs encodes a call to the get_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_offramp::offramp::DynamicConfig",
	})
}

// GetDynamicConfigFields encodes a call to the get_dynamic_config_fields Move function.
func (c offrampEncoder) GetDynamicConfigFields(ref bind.Object, cfg DynamicConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"ccip_offramp::offramp::DynamicConfig",
	}, []any{
		ref,
		cfg,
	}, []string{
		"address",
		"u32",
	})
}

// GetDynamicConfigFieldsWithArgs encodes a call to the get_dynamic_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"ccip_offramp::offramp::DynamicConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"u32",
	})
}

// SetDynamicConfig encodes a call to the set_dynamic_config Move function.
func (c offrampEncoder) SetDynamicConfig(ref bind.Object, state bind.Object, ownerCap bind.Object, permissionlessExecutionThresholdSeconds uint32) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"u32",
	}, []any{
		ref,
		state,
		ownerCap,
		permissionlessExecutionThresholdSeconds,
	}, nil)
}

// SetDynamicConfigWithArgs encodes a call to the set_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"u32",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplySourceChainConfigUpdates encodes a call to the apply_source_chain_config_updates Move function.
func (c offrampEncoder) ApplySourceChainConfigUpdates(ref bind.Object, state bind.Object, ownerCap bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_source_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}, []any{
		ref,
		state,
		ownerCap,
		sourceChainsSelector,
		sourceChainsIsEnabled,
		sourceChainsIsRmnVerificationDisabled,
		sourceChainsOnRamp,
	}, nil)
}

// ApplySourceChainConfigUpdatesWithArgs encodes a call to the apply_source_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ApplySourceChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_source_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCcipPackageId encodes a call to the get_ccip_package_id Move function.
func (c offrampEncoder) GetCcipPackageId() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ccip_package_id", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"address",
	})
}

// GetCcipPackageIdWithArgs encodes a call to the get_ccip_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetCcipPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ccip_package_id", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// Owner encodes a call to the owner Move function.
func (c offrampEncoder) Owner(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// HasPendingTransfer encodes a call to the has_pending_transfer Move function.
func (c offrampEncoder) HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// PendingTransferFrom encodes a call to the pending_transfer_from Move function.
func (c offrampEncoder) PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferTo encodes a call to the pending_transfer_to Move function.
func (c offrampEncoder) PendingTransferTo(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferAccepted encodes a call to the pending_transfer_accepted Move function.
func (c offrampEncoder) PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<bool>",
	})
}

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c offrampEncoder) TransferOwnership(ref bind.Object, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}, []any{
		ref,
		state,
		ownerCap,
		newOwner,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnership encodes a call to the accept_ownership Move function.
func (c offrampEncoder) AcceptOwnership(ref bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
	}, []any{
		ref,
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c offrampEncoder) AcceptOwnershipFromObject(ref bind.Object, state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut UID",
	}, []any{
		ref,
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAcceptOwnership encodes a call to the mcms_accept_ownership Move function.
func (c offrampEncoder) McmsAcceptOwnership(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsAcceptOwnershipWithArgs encodes a call to the mcms_accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c offrampEncoder) ExecuteOwnershipTransfer(ref bind.Object, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OffRampState",
		"address",
	}, []any{
		ref,
		ownerCap,
		state,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OffRampState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransferToMcms encodes a call to the execute_ownership_transfer_to_mcms Move function.
func (c offrampEncoder) ExecuteOwnershipTransferToMcms(ref bind.Object, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OffRampState",
		"&mut Registry",
		"address",
	}, []any{
		ref,
		ownerCap,
		state,
		registry,
		to,
	}, nil)
}

// ExecuteOwnershipTransferToMcmsWithArgs encodes a call to the execute_ownership_transfer_to_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OffRampState",
		"&mut Registry",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRegisterUpgradeCap encodes a call to the mcms_register_upgrade_cap Move function.
func (c offrampEncoder) McmsRegisterUpgradeCap(ref bind.Object, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}, []any{
		ref,
		upgradeCap,
		registry,
		state,
	}, nil)
}

// McmsRegisterUpgradeCapWithArgs encodes a call to the mcms_register_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAddPackageId encodes a call to the mcms_add_package_id Move function.
func (c offrampEncoder) McmsAddPackageId(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_add_package_id", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsAddPackageIdWithArgs encodes a call to the mcms_add_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsAddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_add_package_id", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRemovePackageId encodes a call to the mcms_remove_package_id Move function.
func (c offrampEncoder) McmsRemovePackageId(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_remove_package_id", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsRemovePackageIdWithArgs encodes a call to the mcms_remove_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsRemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_remove_package_id", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsSetDynamicConfig encodes a call to the mcms_set_dynamic_config Move function.
func (c offrampEncoder) McmsSetDynamicConfig(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_dynamic_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsSetDynamicConfigWithArgs encodes a call to the mcms_set_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsSetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsApplySourceChainConfigUpdates encodes a call to the mcms_apply_source_chain_config_updates Move function.
func (c offrampEncoder) McmsApplySourceChainConfigUpdates(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_source_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsApplySourceChainConfigUpdatesWithArgs encodes a call to the mcms_apply_source_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsApplySourceChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_source_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsSetOcr3Config encodes a call to the mcms_set_ocr3_config Move function.
func (c offrampEncoder) McmsSetOcr3Config(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_ocr3_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsSetOcr3ConfigWithArgs encodes a call to the mcms_set_ocr3_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsSetOcr3ConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_ocr3_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsTransferOwnership encodes a call to the mcms_transfer_ownership Move function.
func (c offrampEncoder) McmsTransferOwnership(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsTransferOwnershipWithArgs encodes a call to the mcms_transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsExecuteOwnershipTransfer encodes a call to the mcms_execute_ownership_transfer Move function.
func (c offrampEncoder) McmsExecuteOwnershipTransfer(ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		deployerState,
		params,
	}, nil)
}

// McmsExecuteOwnershipTransferWithArgs encodes a call to the mcms_execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAddAllowedModules encodes a call to the mcms_add_allowed_modules Move function.
func (c offrampEncoder) McmsAddAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_add_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		params,
	}, nil)
}

// McmsAddAllowedModulesWithArgs encodes a call to the mcms_add_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsAddAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_add_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRemoveAllowedModules encodes a call to the mcms_remove_allowed_modules Move function.
func (c offrampEncoder) McmsRemoveAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_remove_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		params,
	}, nil)
}

// McmsRemoveAllowedModulesWithArgs encodes a call to the mcms_remove_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsRemoveAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_remove_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}
