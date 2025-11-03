// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_onramp

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

const FunctionInfo = `[{"package":"ccip_onramp","module":"onramp","name":"accept_ownership","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"accept_ownership_from_object","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"from","type":"sui::object::UID"}]},{"package":"ccip_onramp","module":"onramp","name":"add_package_id","parameters":[{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"package_id","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"apply_allowlist_updates","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"dest_chain_selectors","type":"vector<u64>"},{"name":"dest_chain_allowlist_enabled","type":"vector<bool>"},{"name":"dest_chain_add_allowed_senders","type":"vector<vector<address>>"},{"name":"dest_chain_remove_allowed_senders","type":"vector<vector<address>>"}]},{"package":"ccip_onramp","module":"onramp","name":"apply_allowlist_updates_by_admin","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"dest_chain_selectors","type":"vector<u64>"},{"name":"dest_chain_allowlist_enabled","type":"vector<bool>"},{"name":"dest_chain_add_allowed_senders","type":"vector<vector<address>>"},{"name":"dest_chain_remove_allowed_senders","type":"vector<vector<address>>"}]},{"package":"ccip_onramp","module":"onramp","name":"apply_dest_chain_config_updates","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"dest_chain_selectors","type":"vector<u64>"},{"name":"dest_chain_allowlist_enabled","type":"vector<bool>"},{"name":"dest_chain_routers","type":"vector<address>"}]},{"package":"ccip_onramp","module":"onramp","name":"calculate_message_hash","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"on_ramp_address","type":"address"},{"name":"message_id","type":"vector<u8>"},{"name":"source_chain_selector","type":"u64"},{"name":"dest_chain_selector","type":"u64"},{"name":"sequence_number","type":"u64"},{"name":"nonce","type":"u64"},{"name":"sender","type":"address"},{"name":"receiver","type":"vector<u8>"},{"name":"data","type":"vector<u8>"},{"name":"fee_token","type":"address"},{"name":"fee_token_amount","type":"u64"},{"name":"source_pool_addresses","type":"vector<address>"},{"name":"dest_token_addresses","type":"vector<vector<u8>>"},{"name":"extra_datas","type":"vector<vector<u8>>"},{"name":"amounts","type":"vector<u64>"},{"name":"dest_exec_datas","type":"vector<vector<u8>>"},{"name":"extra_args","type":"vector<u8>"}]},{"package":"ccip_onramp","module":"onramp","name":"calculate_metadata_hash","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"source_chain_selector","type":"u64"},{"name":"dest_chain_selector","type":"u64"},{"name":"on_ramp_address","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"ccip_send","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"clock","type":"Clock"},{"name":"dest_chain_selector","type":"u64"},{"name":"receiver","type":"vector<u8>"},{"name":"data","type":"vector<u8>"},{"name":"token_params","type":"TokenTransferParams"},{"name":"fee_token_metadata","type":"CoinMetadata<T>"},{"name":"fee_token","type":"Coin<T>"},{"name":"extra_args","type":"vector<u8>"}]},{"package":"ccip_onramp","module":"onramp","name":"execute_ownership_transfer","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"OnRampState"},{"name":"to","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"execute_ownership_transfer_to_mcms","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"OnRampState"},{"name":"registry","type":"Registry"},{"name":"to","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"get_allowed_senders_list","parameters":[{"name":"state","type":"OnRampState"},{"name":"dest_chain_selector","type":"u64"}]},{"package":"ccip_onramp","module":"onramp","name":"get_ccip_package_id","parameters":null},{"package":"ccip_onramp","module":"onramp","name":"get_dest_chain_config","parameters":[{"name":"state","type":"OnRampState"},{"name":"dest_chain_selector","type":"u64"}]},{"package":"ccip_onramp","module":"onramp","name":"get_dynamic_config","parameters":[{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"get_dynamic_config_fields","parameters":[{"name":"cfg","type":"DynamicConfig"}]},{"package":"ccip_onramp","module":"onramp","name":"get_expected_next_sequence_number","parameters":[{"name":"state","type":"OnRampState"},{"name":"dest_chain_selector","type":"u64"}]},{"package":"ccip_onramp","module":"onramp","name":"get_fee","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"clock","type":"Clock"},{"name":"dest_chain_selector","type":"u64"},{"name":"receiver","type":"vector<u8>"},{"name":"data","type":"vector<u8>"},{"name":"token_addresses","type":"vector<address>"},{"name":"token_amounts","type":"vector<u64>"},{"name":"fee_token","type":"CoinMetadata<T>"},{"name":"extra_args","type":"vector<u8>"}]},{"package":"ccip_onramp","module":"onramp","name":"get_outbound_nonce","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"dest_chain_selector","type":"u64"},{"name":"sender","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"get_static_config","parameters":[{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"get_static_config_fields","parameters":[{"name":"cfg","type":"StaticConfig"}]},{"package":"ccip_onramp","module":"onramp","name":"has_pending_transfer","parameters":[{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"initialize","parameters":[{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"nonce_manager_cap","type":"NonceManagerCap"},{"name":"source_transfer_cap","type":"osh::SourceTransferCap"},{"name":"chain_selector","type":"u64"},{"name":"fee_aggregator","type":"address"},{"name":"allowlist_admin","type":"address"},{"name":"dest_chain_selectors","type":"vector<u64>"},{"name":"dest_chain_allowlist_enabled","type":"vector<bool>"},{"name":"dest_chain_routers","type":"vector<address>"}]},{"package":"ccip_onramp","module":"onramp","name":"is_chain_supported","parameters":[{"name":"state","type":"OnRampState"},{"name":"dest_chain_selector","type":"u64"}]},{"package":"ccip_onramp","module":"onramp","name":"owner","parameters":[{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"pending_transfer_accepted","parameters":[{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"pending_transfer_from","parameters":[{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"pending_transfer_to","parameters":[{"name":"state","type":"OnRampState"}]},{"package":"ccip_onramp","module":"onramp","name":"remove_package_id","parameters":[{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"package_id","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"set_dynamic_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"fee_aggregator","type":"address"},{"name":"allowlist_admin","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"transfer_ownership","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"new_owner","type":"address"}]},{"package":"ccip_onramp","module":"onramp","name":"type_and_version","parameters":null},{"package":"ccip_onramp","module":"onramp","name":"withdraw_fee_tokens","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"OnRampState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"fee_token_metadata","type":"CoinMetadata<T>"}]}]`

type IOnramp interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*models.SuiTransactionBlockResponse, error)
	AddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error)
	RemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	WithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object, feeTokenMetadata bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error)
	SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, feeAggregator string, allowlistAdmin string) (*models.SuiTransactionBlockResponse, error)
	ApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*models.SuiTransactionBlockResponse, error)
	GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdatesByAdmin(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error)
	GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*models.SuiTransactionBlockResponse, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*models.SuiTransactionBlockResponse, error)
	CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error)
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
	McmsApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsWithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOnrampDevInspect
	Encoder() OnrampEncoder
	Bound() bind.IBoundContract
}

type IOnrampDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (bool, error)
	GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (uint64, error)
	GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (uint64, error)
	GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error)
	GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error)
	GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (StaticConfig, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (uint64, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (DynamicConfig, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) ([]any, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) ([]byte, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) ([]byte, error)
	CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) ([]byte, error)
	GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type OnrampEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(state bind.Object, ownerCap bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	AddPackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error)
	AddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	RemovePackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error)
	RemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	IsChainSupported(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	IsChainSupportedWithArgs(args ...any) (*bind.EncodedCall, error)
	GetExpectedNextSequenceNumber(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetExpectedNextSequenceNumberWithArgs(args ...any) (*bind.EncodedCall, error)
	WithdrawFeeTokens(typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object, feeTokenMetadata bind.Object) (*bind.EncodedCall, error)
	WithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetFee(typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error)
	GetFeeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetDynamicConfig(ref bind.Object, state bind.Object, ownerCap bind.Object, feeAggregator string, allowlistAdmin string) (*bind.EncodedCall, error)
	SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyDestChainConfigUpdates(ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*bind.EncodedCall, error)
	ApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestChainConfig(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetDestChainConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowedSendersList(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetAllowedSendersListWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyAllowlistUpdates(ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesByAdmin(ref bind.Object, state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesByAdminWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error)
	GetOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfig(state bind.Object) (*bind.EncodedCall, error)
	GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error)
	GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfig(state bind.Object) (*bind.EncodedCall, error)
	GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfigFields(cfg DynamicConfig) (*bind.EncodedCall, error)
	GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMessageHash(ref bind.Object, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*bind.EncodedCall, error)
	CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMetadataHash(ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*bind.EncodedCall, error)
	CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error)
	CcipSend(typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error)
	CcipSendWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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
	McmsApplyDestChainConfigUpdates(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsApplyAllowlistUpdates(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsTransferOwnership(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransfer(ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsWithdrawFeeTokens(typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsWithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsAddAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAddAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRemoveAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRemoveAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error)
}

type OnrampContract struct {
	*bind.BoundContract
	onrampEncoder
	devInspect *OnrampDevInspect
}

type OnrampDevInspect struct {
	contract *OnrampContract
}

var _ IOnramp = (*OnrampContract)(nil)
var _ IOnrampDevInspect = (*OnrampDevInspect)(nil)

func NewOnramp(packageID string, client sui.ISuiAPI) (IOnramp, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_onramp", "onramp", client)
	if err != nil {
		return nil, err
	}

	c := &OnrampContract{
		BoundContract: contract,
		onrampEncoder: onrampEncoder{BoundContract: contract},
	}
	c.devInspect = &OnrampDevInspect{contract: c}
	return c, nil
}

func (c *OnrampContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *OnrampContract) Encoder() OnrampEncoder {
	return c.onrampEncoder
}

func (c *OnrampContract) DevInspect() IOnrampDevInspect {
	return c.devInspect
}

type OnRampState struct {
	Id                string       `move:"sui::object::UID"`
	PackageIds        []string     `move:"vector<address>"`
	ChainSelector     uint64       `move:"u64"`
	FeeAggregator     string       `move:"address"`
	AllowlistAdmin    string       `move:"address"`
	DestChainConfigs  bind.Object  `move:"Table<u64, DestChainConfig>"`
	FeeTokens         bind.Object  `move:"Bag"`
	NonceManagerCap   *bind.Object `move:"0x1::option::Option<NonceManagerCap>"`
	SourceTransferCap *bind.Object `move:"0x1::option::Option<osh::SourceTransferCap>"`
	OwnableState      bind.Object  `move:"OwnableState"`
}

type OnRampObject struct {
	Id string `move:"sui::object::UID"`
}

type OnRampStatePointer struct {
	Id             string `move:"sui::object::UID"`
	OnRampObjectId string `move:"address"`
}

type DestChainConfig struct {
	SequenceNumber   uint64   `move:"u64"`
	AllowlistEnabled bool     `move:"bool"`
	AllowedSenders   []string `move:"vector<address>"`
	Router           string   `move:"address"`
}

type RampMessageHeader struct {
	MessageId           []byte `move:"vector<u8>"`
	SourceChainSelector uint64 `move:"u64"`
	DestChainSelector   uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
	Nonce               uint64 `move:"u64"`
}

type Sui2AnyRampMessage struct {
	Header         RampMessageHeader      `move:"RampMessageHeader"`
	Sender         string                 `move:"address"`
	Data           []byte                 `move:"vector<u8>"`
	Receiver       []byte                 `move:"vector<u8>"`
	ExtraArgs      []byte                 `move:"vector<u8>"`
	FeeToken       string                 `move:"address"`
	FeeTokenAmount uint64                 `move:"u64"`
	FeeValueJuels  *big.Int               `move:"u256"`
	TokenAmounts   []Sui2AnyTokenTransfer `move:"vector<Sui2AnyTokenTransfer>"`
}

type Sui2AnyTokenTransfer struct {
	SourcePoolAddress string `move:"address"`
	DestTokenAddress  []byte `move:"vector<u8>"`
	ExtraData         []byte `move:"vector<u8>"`
	Amount            uint64 `move:"u64"`
	DestExecData      []byte `move:"vector<u8>"`
}

type StaticConfig struct {
	ChainSelector uint64 `move:"u64"`
}

type DynamicConfig struct {
	FeeAggregator  string `move:"address"`
	AllowlistAdmin string `move:"address"`
}

type ConfigSet struct {
	StaticConfig  StaticConfig  `move:"StaticConfig"`
	DynamicConfig DynamicConfig `move:"DynamicConfig"`
}

type DestChainConfigSet struct {
	DestChainSelector uint64 `move:"u64"`
	SequenceNumber    uint64 `move:"u64"`
	AllowlistEnabled  bool   `move:"bool"`
	Router            string `move:"address"`
}

type CCIPMessageSent struct {
	DestChainSelector uint64             `move:"u64"`
	SequenceNumber    uint64             `move:"u64"`
	Message           Sui2AnyRampMessage `move:"Sui2AnyRampMessage"`
}

type AllowlistSendersAdded struct {
	DestChainSelector uint64   `move:"u64"`
	Senders           []string `move:"vector<address>"`
}

type AllowlistSendersRemoved struct {
	DestChainSelector uint64   `move:"u64"`
	Senders           []string `move:"vector<address>"`
}

type FeeTokenWithdrawn struct {
	FeeAggregator string `move:"address"`
	FeeToken      string `move:"address"`
	Amount        uint64 `move:"u64"`
}

type ONRAMP struct {
}

type McmsCallback struct {
}

type McmsAcceptOwnershipProof struct {
}

type bcsOnRampState struct {
	Id                string
	PackageIds        [][32]byte
	ChainSelector     uint64
	FeeAggregator     [32]byte
	AllowlistAdmin    [32]byte
	DestChainConfigs  bind.Object
	FeeTokens         bind.Object
	NonceManagerCap   *bind.Object
	SourceTransferCap *bind.Object
	OwnableState      bind.Object
}

func convertOnRampStateFromBCS(bcs bcsOnRampState) (OnRampState, error) {

	return OnRampState{
		Id: bcs.Id,
		PackageIds: func() []string {
			addrs := make([]string, len(bcs.PackageIds))
			for i, addr := range bcs.PackageIds {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		ChainSelector:     bcs.ChainSelector,
		FeeAggregator:     fmt.Sprintf("0x%x", bcs.FeeAggregator),
		AllowlistAdmin:    fmt.Sprintf("0x%x", bcs.AllowlistAdmin),
		DestChainConfigs:  bcs.DestChainConfigs,
		FeeTokens:         bcs.FeeTokens,
		NonceManagerCap:   bcs.NonceManagerCap,
		SourceTransferCap: bcs.SourceTransferCap,
		OwnableState:      bcs.OwnableState,
	}, nil
}

type bcsOnRampStatePointer struct {
	Id             string
	OnRampObjectId [32]byte
}

func convertOnRampStatePointerFromBCS(bcs bcsOnRampStatePointer) (OnRampStatePointer, error) {

	return OnRampStatePointer{
		Id:             bcs.Id,
		OnRampObjectId: fmt.Sprintf("0x%x", bcs.OnRampObjectId),
	}, nil
}

type bcsDestChainConfig struct {
	SequenceNumber   uint64
	AllowlistEnabled bool
	AllowedSenders   [][32]byte
	Router           [32]byte
}

func convertDestChainConfigFromBCS(bcs bcsDestChainConfig) (DestChainConfig, error) {

	return DestChainConfig{
		SequenceNumber:   bcs.SequenceNumber,
		AllowlistEnabled: bcs.AllowlistEnabled,
		AllowedSenders: func() []string {
			addrs := make([]string, len(bcs.AllowedSenders))
			for i, addr := range bcs.AllowedSenders {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		Router: fmt.Sprintf("0x%x", bcs.Router),
	}, nil
}

type bcsSui2AnyRampMessage struct {
	Header         RampMessageHeader
	Sender         [32]byte
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       [32]byte
	FeeTokenAmount uint64
	FeeValueJuels  [32]byte
	TokenAmounts   []Sui2AnyTokenTransfer
}

func convertSui2AnyRampMessageFromBCS(bcs bcsSui2AnyRampMessage) (Sui2AnyRampMessage, error) {
	FeeValueJuelsField, err := bind.DecodeU256Value(bcs.FeeValueJuels)
	if err != nil {
		return Sui2AnyRampMessage{}, fmt.Errorf("failed to decode u256 field FeeValueJuels: %w", err)
	}

	return Sui2AnyRampMessage{
		Header:         bcs.Header,
		Sender:         fmt.Sprintf("0x%x", bcs.Sender),
		Data:           bcs.Data,
		Receiver:       bcs.Receiver,
		ExtraArgs:      bcs.ExtraArgs,
		FeeToken:       fmt.Sprintf("0x%x", bcs.FeeToken),
		FeeTokenAmount: bcs.FeeTokenAmount,
		FeeValueJuels:  FeeValueJuelsField,
		TokenAmounts:   bcs.TokenAmounts,
	}, nil
}

type bcsSui2AnyTokenTransfer struct {
	SourcePoolAddress [32]byte
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            uint64
	DestExecData      []byte
}

func convertSui2AnyTokenTransferFromBCS(bcs bcsSui2AnyTokenTransfer) (Sui2AnyTokenTransfer, error) {

	return Sui2AnyTokenTransfer{
		SourcePoolAddress: fmt.Sprintf("0x%x", bcs.SourcePoolAddress),
		DestTokenAddress:  bcs.DestTokenAddress,
		ExtraData:         bcs.ExtraData,
		Amount:            bcs.Amount,
		DestExecData:      bcs.DestExecData,
	}, nil
}

type bcsDynamicConfig struct {
	FeeAggregator  [32]byte
	AllowlistAdmin [32]byte
}

func convertDynamicConfigFromBCS(bcs bcsDynamicConfig) (DynamicConfig, error) {

	return DynamicConfig{
		FeeAggregator:  fmt.Sprintf("0x%x", bcs.FeeAggregator),
		AllowlistAdmin: fmt.Sprintf("0x%x", bcs.AllowlistAdmin),
	}, nil
}

type bcsConfigSet struct {
	StaticConfig  StaticConfig
	DynamicConfig bcsDynamicConfig
}

func convertConfigSetFromBCS(bcs bcsConfigSet) (ConfigSet, error) {
	DynamicConfigField, err := convertDynamicConfigFromBCS(bcs.DynamicConfig)
	if err != nil {
		return ConfigSet{}, fmt.Errorf("failed to convert nested struct DynamicConfig: %w", err)
	}

	return ConfigSet{
		StaticConfig:  bcs.StaticConfig,
		DynamicConfig: DynamicConfigField,
	}, nil
}

type bcsDestChainConfigSet struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	AllowlistEnabled  bool
	Router            [32]byte
}

func convertDestChainConfigSetFromBCS(bcs bcsDestChainConfigSet) (DestChainConfigSet, error) {

	return DestChainConfigSet{
		DestChainSelector: bcs.DestChainSelector,
		SequenceNumber:    bcs.SequenceNumber,
		AllowlistEnabled:  bcs.AllowlistEnabled,
		Router:            fmt.Sprintf("0x%x", bcs.Router),
	}, nil
}

type bcsCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           bcsSui2AnyRampMessage
}

func convertCCIPMessageSentFromBCS(bcs bcsCCIPMessageSent) (CCIPMessageSent, error) {
	MessageField, err := convertSui2AnyRampMessageFromBCS(bcs.Message)
	if err != nil {
		return CCIPMessageSent{}, fmt.Errorf("failed to convert nested struct Message: %w", err)
	}

	return CCIPMessageSent{
		DestChainSelector: bcs.DestChainSelector,
		SequenceNumber:    bcs.SequenceNumber,
		Message:           MessageField,
	}, nil
}

type bcsAllowlistSendersAdded struct {
	DestChainSelector uint64
	Senders           [][32]byte
}

func convertAllowlistSendersAddedFromBCS(bcs bcsAllowlistSendersAdded) (AllowlistSendersAdded, error) {

	return AllowlistSendersAdded{
		DestChainSelector: bcs.DestChainSelector,
		Senders: func() []string {
			addrs := make([]string, len(bcs.Senders))
			for i, addr := range bcs.Senders {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}, nil
}

type bcsAllowlistSendersRemoved struct {
	DestChainSelector uint64
	Senders           [][32]byte
}

func convertAllowlistSendersRemovedFromBCS(bcs bcsAllowlistSendersRemoved) (AllowlistSendersRemoved, error) {

	return AllowlistSendersRemoved{
		DestChainSelector: bcs.DestChainSelector,
		Senders: func() []string {
			addrs := make([]string, len(bcs.Senders))
			for i, addr := range bcs.Senders {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}, nil
}

type bcsFeeTokenWithdrawn struct {
	FeeAggregator [32]byte
	FeeToken      [32]byte
	Amount        uint64
}

func convertFeeTokenWithdrawnFromBCS(bcs bcsFeeTokenWithdrawn) (FeeTokenWithdrawn, error) {

	return FeeTokenWithdrawn{
		FeeAggregator: fmt.Sprintf("0x%x", bcs.FeeAggregator),
		FeeToken:      fmt.Sprintf("0x%x", bcs.FeeToken),
		Amount:        bcs.Amount,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip_onramp::onramp::OnRampState", func(data []byte) (interface{}, error) {
		var temp bcsOnRampState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOnRampStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OnRampState
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::OnRampState>", func(data []byte) (interface{}, error) {
		var temps []bcsOnRampState
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OnRampState, len(temps))
		for i, temp := range temps {
			result, err := convertOnRampStateFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::OnRampObject", func(data []byte) (interface{}, error) {
		var result OnRampObject
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OnRampObject
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::OnRampObject>", func(data []byte) (interface{}, error) {
		var results []OnRampObject
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::OnRampStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsOnRampStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOnRampStatePointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OnRampStatePointer
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::OnRampStatePointer>", func(data []byte) (interface{}, error) {
		var temps []bcsOnRampStatePointer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OnRampStatePointer, len(temps))
		for i, temp := range temps {
			result, err := convertOnRampStatePointerFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::DestChainConfig", func(data []byte) (interface{}, error) {
		var temp bcsDestChainConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDestChainConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DestChainConfig
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::DestChainConfig>", func(data []byte) (interface{}, error) {
		var temps []bcsDestChainConfig
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]DestChainConfig, len(temps))
		for i, temp := range temps {
			result, err := convertDestChainConfigFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::RampMessageHeader", func(data []byte) (interface{}, error) {
		var result RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for RampMessageHeader
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::RampMessageHeader>", func(data []byte) (interface{}, error) {
		var results []RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::Sui2AnyRampMessage", func(data []byte) (interface{}, error) {
		var temp bcsSui2AnyRampMessage
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSui2AnyRampMessageFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for Sui2AnyRampMessage
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::Sui2AnyRampMessage>", func(data []byte) (interface{}, error) {
		var temps []bcsSui2AnyRampMessage
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]Sui2AnyRampMessage, len(temps))
		for i, temp := range temps {
			result, err := convertSui2AnyRampMessageFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::Sui2AnyTokenTransfer", func(data []byte) (interface{}, error) {
		var temp bcsSui2AnyTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSui2AnyTokenTransferFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for Sui2AnyTokenTransfer
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::Sui2AnyTokenTransfer>", func(data []byte) (interface{}, error) {
		var temps []bcsSui2AnyTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]Sui2AnyTokenTransfer, len(temps))
		for i, temp := range temps {
			result, err := convertSui2AnyTokenTransferFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::StaticConfig", func(data []byte) (interface{}, error) {
		var result StaticConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for StaticConfig
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::StaticConfig>", func(data []byte) (interface{}, error) {
		var results []StaticConfig
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::DynamicConfig", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::DynamicConfig>", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::ConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertConfigSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ConfigSet
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::ConfigSet>", func(data []byte) (interface{}, error) {
		var temps []bcsConfigSet
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]ConfigSet, len(temps))
		for i, temp := range temps {
			result, err := convertConfigSetFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::DestChainConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsDestChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDestChainConfigSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DestChainConfigSet
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::DestChainConfigSet>", func(data []byte) (interface{}, error) {
		var temps []bcsDestChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]DestChainConfigSet, len(temps))
		for i, temp := range temps {
			result, err := convertDestChainConfigSetFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::CCIPMessageSent", func(data []byte) (interface{}, error) {
		var temp bcsCCIPMessageSent
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCCIPMessageSentFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for CCIPMessageSent
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::CCIPMessageSent>", func(data []byte) (interface{}, error) {
		var temps []bcsCCIPMessageSent
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]CCIPMessageSent, len(temps))
		for i, temp := range temps {
			result, err := convertCCIPMessageSentFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::AllowlistSendersAdded", func(data []byte) (interface{}, error) {
		var temp bcsAllowlistSendersAdded
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAllowlistSendersAddedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for AllowlistSendersAdded
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::AllowlistSendersAdded>", func(data []byte) (interface{}, error) {
		var temps []bcsAllowlistSendersAdded
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]AllowlistSendersAdded, len(temps))
		for i, temp := range temps {
			result, err := convertAllowlistSendersAddedFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::AllowlistSendersRemoved", func(data []byte) (interface{}, error) {
		var temp bcsAllowlistSendersRemoved
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAllowlistSendersRemovedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for AllowlistSendersRemoved
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::AllowlistSendersRemoved>", func(data []byte) (interface{}, error) {
		var temps []bcsAllowlistSendersRemoved
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]AllowlistSendersRemoved, len(temps))
		for i, temp := range temps {
			result, err := convertAllowlistSendersRemovedFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::FeeTokenWithdrawn", func(data []byte) (interface{}, error) {
		var temp bcsFeeTokenWithdrawn
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFeeTokenWithdrawnFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for FeeTokenWithdrawn
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::FeeTokenWithdrawn>", func(data []byte) (interface{}, error) {
		var temps []bcsFeeTokenWithdrawn
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]FeeTokenWithdrawn, len(temps))
		for i, temp := range temps {
			result, err := convertFeeTokenWithdrawnFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::ONRAMP", func(data []byte) (interface{}, error) {
		var result ONRAMP
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ONRAMP
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::ONRAMP>", func(data []byte) (interface{}, error) {
		var results []ONRAMP
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsCallback
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::McmsCallback>", func(data []byte) (interface{}, error) {
		var results []McmsCallback
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::McmsAcceptOwnershipProof", func(data []byte) (interface{}, error) {
		var result McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsAcceptOwnershipProof
	bind.RegisterStructDecoder("vector<ccip_onramp::onramp::McmsAcceptOwnershipProof>", func(data []byte) (interface{}, error) {
		var results []McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *OnrampContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *OnrampContract) Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.Initialize(state, ownerCap, nonceManagerCap, sourceTransferCap, chainSelector, feeAggregator, allowlistAdmin, destChainSelectors, destChainAllowlistEnabled, destChainRouters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddPackageId executes the add_package_id Move function.
func (c *OnrampContract) AddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.AddPackageId(state, ownerCap, packageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemovePackageId executes the remove_package_id Move function.
func (c *OnrampContract) RemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.RemovePackageId(state, ownerCap, packageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsChainSupported executes the is_chain_supported Move function.
func (c *OnrampContract) IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.IsChainSupported(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetExpectedNextSequenceNumber executes the get_expected_next_sequence_number Move function.
func (c *OnrampContract) GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetExpectedNextSequenceNumber(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// WithdrawFeeTokens executes the withdraw_fee_tokens Move function.
func (c *OnrampContract) WithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object, feeTokenMetadata bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.WithdrawFeeTokens(typeArgs, ref, state, ownerCap, feeTokenMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFee executes the get_fee Move function.
func (c *OnrampContract) GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetFee(typeArgs, ref, clock, destChainSelector, receiver, data, tokenAddresses, tokenAmounts, feeToken, extraArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetDynamicConfig executes the set_dynamic_config Move function.
func (c *OnrampContract) SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, feeAggregator string, allowlistAdmin string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.SetDynamicConfig(ref, state, ownerCap, feeAggregator, allowlistAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyDestChainConfigUpdates executes the apply_dest_chain_config_updates Move function.
func (c *OnrampContract) ApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ApplyDestChainConfigUpdates(ref, state, ownerCap, destChainSelectors, destChainAllowlistEnabled, destChainRouters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestChainConfig executes the get_dest_chain_config Move function.
func (c *OnrampContract) GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetDestChainConfig(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowedSendersList executes the get_allowed_senders_list Move function.
func (c *OnrampContract) GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetAllowedSendersList(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *OnrampContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ApplyAllowlistUpdates(ref, state, ownerCap, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdatesByAdmin executes the apply_allowlist_updates_by_admin Move function.
func (c *OnrampContract) ApplyAllowlistUpdatesByAdmin(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ApplyAllowlistUpdatesByAdmin(ref, state, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOutboundNonce executes the get_outbound_nonce Move function.
func (c *OnrampContract) GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetOutboundNonce(ref, destChainSelector, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function.
func (c *OnrampContract) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetStaticConfig(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfigFields executes the get_static_config_fields Move function.
func (c *OnrampContract) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetStaticConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfig executes the get_dynamic_config Move function.
func (c *OnrampContract) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetDynamicConfig(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfigFields executes the get_dynamic_config_fields Move function.
func (c *OnrampContract) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetDynamicConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMessageHash executes the calculate_message_hash Move function.
func (c *OnrampContract) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.CalculateMessageHash(ref, onRampAddress, messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, data, feeToken, feeTokenAmount, sourcePoolAddresses, destTokenAddresses, extraDatas, amounts, destExecDatas, extraArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMetadataHash executes the calculate_metadata_hash Move function.
func (c *OnrampContract) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.CalculateMetadataHash(ref, sourceChainSelector, destChainSelector, onRampAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CcipSend executes the ccip_send Move function.
func (c *OnrampContract) CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.CcipSend(typeArgs, ref, state, clock, destChainSelector, receiver, data, tokenParams, feeTokenMetadata, feeToken, extraArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCcipPackageId executes the get_ccip_package_id Move function.
func (c *OnrampContract) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetCcipPackageId()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *OnrampContract) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.Owner(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *OnrampContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.HasPendingTransfer(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *OnrampContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *OnrampContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *OnrampContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *OnrampContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.TransferOwnership(ref, state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *OnrampContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.AcceptOwnership(ref, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *OnrampContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.AcceptOwnershipFromObject(ref, state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAcceptOwnership executes the mcms_accept_ownership Move function.
func (c *OnrampContract) McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsAcceptOwnership(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *OnrampContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ExecuteOwnershipTransfer(ref, ownerCap, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *OnrampContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ExecuteOwnershipTransferToMcms(ref, ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *OnrampContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsRegisterUpgradeCap(ref, upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddPackageId executes the mcms_add_package_id Move function.
func (c *OnrampContract) McmsAddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsAddPackageId(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemovePackageId executes the mcms_remove_package_id Move function.
func (c *OnrampContract) McmsRemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsRemovePackageId(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetDynamicConfig executes the mcms_set_dynamic_config Move function.
func (c *OnrampContract) McmsSetDynamicConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsSetDynamicConfig(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyDestChainConfigUpdates executes the mcms_apply_dest_chain_config_updates Move function.
func (c *OnrampContract) McmsApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsApplyDestChainConfigUpdates(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyAllowlistUpdates executes the mcms_apply_allowlist_updates Move function.
func (c *OnrampContract) McmsApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsApplyAllowlistUpdates(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsTransferOwnership executes the mcms_transfer_ownership Move function.
func (c *OnrampContract) McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsTransferOwnership(ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsExecuteOwnershipTransfer executes the mcms_execute_ownership_transfer Move function.
func (c *OnrampContract) McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsExecuteOwnershipTransfer(ref, state, registry, deployerState, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsWithdrawFeeTokens executes the mcms_withdraw_fee_tokens Move function.
func (c *OnrampContract) McmsWithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsWithdrawFeeTokens(typeArgs, ref, state, registry, feeTokenMetadata, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddAllowedModules executes the mcms_add_allowed_modules Move function.
func (c *OnrampContract) McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsAddAllowedModules(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemoveAllowedModules executes the mcms_remove_allowed_modules Move function.
func (c *OnrampContract) McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsRemoveAllowedModules(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *OnrampDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.onrampEncoder.TypeAndVersion()
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

// IsChainSupported executes the is_chain_supported Move function using DevInspect to get return values.
//
// Returns: bool
func (d *OnrampDevInspect) IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (bool, error) {
	encoded, err := d.contract.onrampEncoder.IsChainSupported(state, destChainSelector)
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

// GetExpectedNextSequenceNumber executes the get_expected_next_sequence_number Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OnrampDevInspect) GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetExpectedNextSequenceNumber(state, destChainSelector)
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

// GetFee executes the get_fee Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OnrampDevInspect) GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetFee(typeArgs, ref, clock, destChainSelector, receiver, data, tokenAddresses, tokenAmounts, feeToken, extraArgs)
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

// GetDestChainConfig executes the get_dest_chain_config Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u64
//	[1]: bool
//	[2]: address
func (d *OnrampDevInspect) GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error) {
	encoded, err := d.contract.onrampEncoder.GetDestChainConfig(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetAllowedSendersList executes the get_allowed_senders_list Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: bool
//	[1]: vector<address>
func (d *OnrampDevInspect) GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error) {
	encoded, err := d.contract.onrampEncoder.GetAllowedSendersList(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetOutboundNonce executes the get_outbound_nonce Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OnrampDevInspect) GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetOutboundNonce(ref, destChainSelector, sender)
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

// GetStaticConfig executes the get_static_config Move function using DevInspect to get return values.
//
// Returns: StaticConfig
func (d *OnrampDevInspect) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (StaticConfig, error) {
	encoded, err := d.contract.onrampEncoder.GetStaticConfig(state)
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
// Returns: u64
func (d *OnrampDevInspect) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetStaticConfigFields(cfg)
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

// GetDynamicConfig executes the get_dynamic_config Move function using DevInspect to get return values.
//
// Returns: DynamicConfig
func (d *OnrampDevInspect) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (DynamicConfig, error) {
	encoded, err := d.contract.onrampEncoder.GetDynamicConfig(state)
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
//	[1]: address
func (d *OnrampDevInspect) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) ([]any, error) {
	encoded, err := d.contract.onrampEncoder.GetDynamicConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// CalculateMessageHash executes the calculate_message_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OnrampDevInspect) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) ([]byte, error) {
	encoded, err := d.contract.onrampEncoder.CalculateMessageHash(ref, onRampAddress, messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, data, feeToken, feeTokenAmount, sourcePoolAddresses, destTokenAddresses, extraDatas, amounts, destExecDatas, extraArgs)
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

// CalculateMetadataHash executes the calculate_metadata_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OnrampDevInspect) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) ([]byte, error) {
	encoded, err := d.contract.onrampEncoder.CalculateMetadataHash(ref, sourceChainSelector, destChainSelector, onRampAddress)
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

// CcipSend executes the ccip_send Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OnrampDevInspect) CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) ([]byte, error) {
	encoded, err := d.contract.onrampEncoder.CcipSend(typeArgs, ref, state, clock, destChainSelector, receiver, data, tokenParams, feeTokenMetadata, feeToken, extraArgs)
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

// GetCcipPackageId executes the get_ccip_package_id Move function using DevInspect to get return values.
//
// Returns: address
func (d *OnrampDevInspect) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.onrampEncoder.GetCcipPackageId()
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
func (d *OnrampDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.onrampEncoder.Owner(state)
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
func (d *OnrampDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error) {
	encoded, err := d.contract.onrampEncoder.HasPendingTransfer(state)
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
func (d *OnrampDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.onrampEncoder.PendingTransferFrom(state)
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
func (d *OnrampDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.onrampEncoder.PendingTransferTo(state)
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
func (d *OnrampDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error) {
	encoded, err := d.contract.onrampEncoder.PendingTransferAccepted(state)
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

type onrampEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c onrampEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) Initialize(state bind.Object, ownerCap bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&OwnerCap",
		"NonceManagerCap",
		"osh::SourceTransferCap",
		"u64",
		"address",
		"address",
		"vector<u64>",
		"vector<bool>",
		"vector<address>",
	}, []any{
		state,
		ownerCap,
		nonceManagerCap,
		sourceTransferCap,
		chainSelector,
		feeAggregator,
		allowlistAdmin,
		destChainSelectors,
		destChainAllowlistEnabled,
		destChainRouters,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&OwnerCap",
		"NonceManagerCap",
		"osh::SourceTransferCap",
		"u64",
		"address",
		"address",
		"vector<u64>",
		"vector<bool>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AddPackageId encodes a call to the add_package_id Move function.
func (c onrampEncoder) AddPackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_package_id", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
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
func (c onrampEncoder) AddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
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
func (c onrampEncoder) RemovePackageId(state bind.Object, ownerCap bind.Object, packageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_package_id", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
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
func (c onrampEncoder) RemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
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

// IsChainSupported encodes a call to the is_chain_supported Move function.
func (c onrampEncoder) IsChainSupported(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_chain_supported", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"bool",
	})
}

// IsChainSupportedWithArgs encodes a call to the is_chain_supported Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) IsChainSupportedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_chain_supported", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetExpectedNextSequenceNumber encodes a call to the get_expected_next_sequence_number Move function.
func (c onrampEncoder) GetExpectedNextSequenceNumber(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_expected_next_sequence_number", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"u64",
	})
}

// GetExpectedNextSequenceNumberWithArgs encodes a call to the get_expected_next_sequence_number Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetExpectedNextSequenceNumberWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_expected_next_sequence_number", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// WithdrawFeeTokens encodes a call to the withdraw_fee_tokens Move function.
func (c onrampEncoder) WithdrawFeeTokens(typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object, feeTokenMetadata bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("withdraw_fee_tokens", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"&CoinMetadata<T>",
	}, []any{
		ref,
		state,
		ownerCap,
		feeTokenMetadata,
	}, nil)
}

// WithdrawFeeTokensWithArgs encodes a call to the withdraw_fee_tokens Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) WithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"&CoinMetadata<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("withdraw_fee_tokens", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetFee encodes a call to the get_fee Move function.
func (c onrampEncoder) GetFee(typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_fee", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"vector<address>",
		"vector<u64>",
		"&CoinMetadata<T>",
		"vector<u8>",
	}, []any{
		ref,
		clock,
		destChainSelector,
		receiver,
		data,
		tokenAddresses,
		tokenAmounts,
		feeToken,
		extraArgs,
	}, []string{
		"u64",
	})
}

// GetFeeWithArgs encodes a call to the get_fee Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetFeeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"vector<address>",
		"vector<u64>",
		"&CoinMetadata<T>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_fee", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// SetDynamicConfig encodes a call to the set_dynamic_config Move function.
func (c onrampEncoder) SetDynamicConfig(ref bind.Object, state bind.Object, ownerCap bind.Object, feeAggregator string, allowlistAdmin string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"address",
		"address",
	}, []any{
		ref,
		state,
		ownerCap,
		feeAggregator,
		allowlistAdmin,
	}, nil)
}

// SetDynamicConfigWithArgs encodes a call to the set_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplyDestChainConfigUpdates encodes a call to the apply_dest_chain_config_updates Move function.
func (c onrampEncoder) ApplyDestChainConfigUpdates(ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainRouters []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_dest_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<address>",
	}, []any{
		ref,
		state,
		ownerCap,
		destChainSelectors,
		destChainAllowlistEnabled,
		destChainRouters,
	}, nil)
}

// ApplyDestChainConfigUpdatesWithArgs encodes a call to the apply_dest_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) ApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_dest_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetDestChainConfig encodes a call to the get_dest_chain_config Move function.
func (c onrampEncoder) GetDestChainConfig(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"u64",
		"bool",
		"address",
	})
}

// GetDestChainConfigWithArgs encodes a call to the get_dest_chain_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetDestChainConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
		"bool",
		"address",
	})
}

// GetAllowedSendersList encodes a call to the get_allowed_senders_list Move function.
func (c onrampEncoder) GetAllowedSendersList(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_senders_list", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"bool",
		"vector<address>",
	})
}

// GetAllowedSendersListWithArgs encodes a call to the get_allowed_senders_list Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetAllowedSendersListWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_senders_list", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
		"vector<address>",
	})
}

// ApplyAllowlistUpdates encodes a call to the apply_allowlist_updates Move function.
func (c onrampEncoder) ApplyAllowlistUpdates(ref bind.Object, state bind.Object, ownerCap bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}, []any{
		ref,
		state,
		ownerCap,
		destChainSelectors,
		destChainAllowlistEnabled,
		destChainAddAllowedSenders,
		destChainRemoveAllowedSenders,
	}, nil)
}

// ApplyAllowlistUpdatesWithArgs encodes a call to the apply_allowlist_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplyAllowlistUpdatesByAdmin encodes a call to the apply_allowlist_updates_by_admin Move function.
func (c onrampEncoder) ApplyAllowlistUpdatesByAdmin(ref bind.Object, state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates_by_admin", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}, []any{
		ref,
		state,
		destChainSelectors,
		destChainAllowlistEnabled,
		destChainAddAllowedSenders,
		destChainRemoveAllowedSenders,
	}, nil)
}

// ApplyAllowlistUpdatesByAdminWithArgs encodes a call to the apply_allowlist_updates_by_admin Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) ApplyAllowlistUpdatesByAdminWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates_by_admin", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetOutboundNonce encodes a call to the get_outbound_nonce Move function.
func (c onrampEncoder) GetOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_outbound_nonce", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"address",
	}, []any{
		ref,
		destChainSelector,
		sender,
	}, []string{
		"u64",
	})
}

// GetOutboundNonceWithArgs encodes a call to the get_outbound_nonce Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_outbound_nonce", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetStaticConfig encodes a call to the get_static_config Move function.
func (c onrampEncoder) GetStaticConfig(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"ccip_onramp::onramp::StaticConfig",
	})
}

// GetStaticConfigWithArgs encodes a call to the get_static_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_onramp::onramp::StaticConfig",
	})
}

// GetStaticConfigFields encodes a call to the get_static_config_fields Move function.
func (c onrampEncoder) GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, []string{
		"ccip_onramp::onramp::StaticConfig",
	}, []any{
		cfg,
	}, []string{
		"u64",
	})
}

// GetStaticConfigFieldsWithArgs encodes a call to the get_static_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_onramp::onramp::StaticConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetDynamicConfig encodes a call to the get_dynamic_config Move function.
func (c onrampEncoder) GetDynamicConfig(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"ccip_onramp::onramp::DynamicConfig",
	})
}

// GetDynamicConfigWithArgs encodes a call to the get_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_onramp::onramp::DynamicConfig",
	})
}

// GetDynamicConfigFields encodes a call to the get_dynamic_config_fields Move function.
func (c onrampEncoder) GetDynamicConfigFields(cfg DynamicConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, []string{
		"ccip_onramp::onramp::DynamicConfig",
	}, []any{
		cfg,
	}, []string{
		"address",
		"address",
	})
}

// GetDynamicConfigFieldsWithArgs encodes a call to the get_dynamic_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_onramp::onramp::DynamicConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"address",
	})
}

// CalculateMessageHash encodes a call to the calculate_message_hash Move function.
func (c onrampEncoder) CalculateMessageHash(ref bind.Object, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_message_hash", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"address",
		"vector<u8>",
		"vector<u8>",
		"address",
		"u64",
		"vector<address>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<u8>",
	}, []any{
		ref,
		onRampAddress,
		messageId,
		sourceChainSelector,
		destChainSelector,
		sequenceNumber,
		nonce,
		sender,
		receiver,
		data,
		feeToken,
		feeTokenAmount,
		sourcePoolAddresses,
		destTokenAddresses,
		extraDatas,
		amounts,
		destExecDatas,
		extraArgs,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMessageHashWithArgs encodes a call to the calculate_message_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"address",
		"vector<u8>",
		"vector<u8>",
		"address",
		"u64",
		"vector<address>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<u8>",
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

// CalculateMetadataHash encodes a call to the calculate_metadata_hash Move function.
func (c onrampEncoder) CalculateMetadataHash(ref bind.Object, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_metadata_hash", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"u64",
		"address",
	}, []any{
		ref,
		sourceChainSelector,
		destChainSelector,
		onRampAddress,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMetadataHashWithArgs encodes a call to the calculate_metadata_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
		"u64",
		"address",
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

// CcipSend encodes a call to the ccip_send Move function.
func (c onrampEncoder) CcipSend(typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("ccip_send", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut OnRampState",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"TokenTransferParams",
		"&CoinMetadata<T>",
		"&mut Coin<T>",
		"vector<u8>",
	}, []any{
		ref,
		state,
		clock,
		destChainSelector,
		receiver,
		data,
		tokenParams,
		feeTokenMetadata,
		feeToken,
		extraArgs,
	}, []string{
		"vector<u8>",
	})
}

// CcipSendWithArgs encodes a call to the ccip_send Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) CcipSendWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut OnRampState",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"TokenTransferParams",
		"&CoinMetadata<T>",
		"&mut Coin<T>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("ccip_send", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// GetCcipPackageId encodes a call to the get_ccip_package_id Move function.
func (c onrampEncoder) GetCcipPackageId() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ccip_package_id", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"address",
	})
}

// GetCcipPackageIdWithArgs encodes a call to the get_ccip_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetCcipPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) Owner(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) PendingTransferTo(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) TransferOwnership(ref bind.Object, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) AcceptOwnership(ref bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
	}, []any{
		ref,
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c onrampEncoder) AcceptOwnershipFromObject(ref bind.Object, state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&mut UID",
	}, []any{
		ref,
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsAcceptOwnership(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) ExecuteOwnershipTransfer(ref bind.Object, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OnRampState",
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
func (c onrampEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OnRampState",
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
func (c onrampEncoder) ExecuteOwnershipTransferToMcms(ref bind.Object, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OnRampState",
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
func (c onrampEncoder) ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"OwnerCap",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsRegisterUpgradeCap(ref bind.Object, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) McmsAddPackageId(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_add_package_id", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
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
func (c onrampEncoder) McmsAddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
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
func (c onrampEncoder) McmsRemovePackageId(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_remove_package_id", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
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
func (c onrampEncoder) McmsRemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
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
func (c onrampEncoder) McmsSetDynamicConfig(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_dynamic_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsSetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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

// McmsApplyDestChainConfigUpdates encodes a call to the mcms_apply_dest_chain_config_updates Move function.
func (c onrampEncoder) McmsApplyDestChainConfigUpdates(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_dest_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsApplyDestChainConfigUpdatesWithArgs encodes a call to the mcms_apply_dest_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_dest_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsApplyAllowlistUpdates encodes a call to the mcms_apply_allowlist_updates Move function.
func (c onrampEncoder) McmsApplyAllowlistUpdates(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsApplyAllowlistUpdatesWithArgs encodes a call to the mcms_apply_allowlist_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_allowlist_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsTransferOwnership encodes a call to the mcms_transfer_ownership Move function.
func (c onrampEncoder) McmsTransferOwnership(ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsExecuteOwnershipTransfer(ref bind.Object, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
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

// McmsWithdrawFeeTokens encodes a call to the mcms_withdraw_fee_tokens Move function.
func (c onrampEncoder) McmsWithdrawFeeTokens(typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_withdraw_fee_tokens", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&mut Registry",
		"&CoinMetadata<T>",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		feeTokenMetadata,
		params,
	}, nil)
}

// McmsWithdrawFeeTokensWithArgs encodes a call to the mcms_withdraw_fee_tokens Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsWithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OnRampState",
		"&mut Registry",
		"&CoinMetadata<T>",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_withdraw_fee_tokens", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAddAllowedModules encodes a call to the mcms_add_allowed_modules Move function.
func (c onrampEncoder) McmsAddAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) McmsAddAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) McmsRemoveAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) McmsRemoveAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error) {
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
