// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_onramp

import (
	"context"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_common "github.com/smartcontractkit/chainlink-sui/bindings/common"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
	_ = uint256.NewInt
)

type IOnramp interface {
	TypeAndVersion() bind.IMethod
	Initialize(state bind.Object, param module_common.OwnerCap, nonceManagerCap module_common.NonceManagerCap, sourceTransferCap module_common.SourceTransferCap, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) bind.IMethod
	IsChainSupported(state bind.Object, destChainSelector uint64) bind.IMethod
	GetExpectedNextSequenceNumber(state bind.Object, destChainSelector uint64) bind.IMethod
	WithdrawFeeTokens(typeArgs string, state bind.Object, param module_common.OwnerCap, feeTokenMetadata bind.Object) bind.IMethod
	GetFee(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) bind.IMethod
	SetDynamicConfig(state bind.Object, param module_common.OwnerCap, feeAggregator string, allowlistAdmin string) bind.IMethod
	ApplyDestChainConfigUpdates(state bind.Object, param module_common.OwnerCap, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) bind.IMethod
	GetDestChainConfig(state bind.Object, destChainSelector uint64) bind.IMethod
	GetAllowedSendersList(state bind.Object, destChainSelector uint64) bind.IMethod
	ApplyAllowlistUpdates(state bind.Object, param module_common.OwnerCap, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod
	ApplyAllowlistUpdatesByAdmin(state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod
	GetOutboundNonce(ref module_common.CCIPObjectRef, destChainSelector uint64, sender string) bind.IMethod
	GetStaticConfig(state bind.Object) bind.IMethod
	GetStaticConfigFields(cfg StaticConfig) bind.IMethod
	GetDynamicConfig(state bind.Object) bind.IMethod
	GetDynamicConfigFields(cfg DynamicConfig) bind.IMethod
	CcipSend(typeArgs string, ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, receiver []byte, data []byte, tokenParams module_common.TokenParams, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) bind.IMethod
	GetCcipPackageId() bind.IMethod
	Owner(state bind.Object) bind.IMethod
	HasPendingTransfer(state bind.Object) bind.IMethod
	PendingTransferFrom(state bind.Object) bind.IMethod
	PendingTransferTo(state bind.Object) bind.IMethod
	PendingTransferAccepted(state bind.Object) bind.IMethod
	TransferOwnership(state bind.Object, ownerCap module_common.OwnerCap, newOwner string) bind.IMethod
	AcceptOwnership(state bind.Object) bind.IMethod
	AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type OnrampContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IOnramp = (*OnrampContract)(nil)

func NewOnramp(packageID string, client suiclient.ClientImpl) (*OnrampContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &OnrampContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *OnrampContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OnRampState struct {
	Id                string                           `move:"sui::object::UID"`
	ChainSelector     uint64                           `move:"u64"`
	FeeAggregator     string                           `move:"address"`
	AllowlistAdmin    string                           `move:"address"`
	NonceManagerCap   *module_common.NonceManagerCap   `move:"0x1::option::Option<NonceManagerCap>"`
	SourceTransferCap *module_common.SourceTransferCap `move:"0x1::option::Option<ccip::common::SourceTransferCap>"`
}

type OnRampStatePointer struct {
	Id            string `move:"sui::object::UID"`
	OnRampStateId string `move:"address"`
	OwnerCapId    string `move:"address"`
}

type DestChainConfig struct {
	IsEnabled        bool     `move:"bool"`
	SequenceNumber   uint64   `move:"u64"`
	AllowlistEnabled bool     `move:"bool"`
	AllowedSenders   []string `move:"vector<address>"`
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
	FeeValueJuels  uint256.Int            `move:"u256"`
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
	IsEnabled         bool   `move:"bool"`
	SequenceNumber    uint64 `move:"u64"`
	AllowlistEnabled  bool   `move:"bool"`
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

// Functions

func (c *OnrampContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "type_and_version", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) Initialize(state bind.Object, param module_common.OwnerCap, nonceManagerCap module_common.NonceManagerCap, sourceTransferCap module_common.SourceTransferCap, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "initialize", false, "", "", state, param, nonceManagerCap, sourceTransferCap, chainSelector, feeAggregator, allowlistAdmin, destChainSelectors, destChainEnabled, destChainAllowlistEnabled)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) IsChainSupported(state bind.Object, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "is_chain_supported", false, "", "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "is_chain_supported", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetExpectedNextSequenceNumber(state bind.Object, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_expected_next_sequence_number", false, "", "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_expected_next_sequence_number", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) WithdrawFeeTokens(typeArgs string, state bind.Object, param module_common.OwnerCap, feeTokenMetadata bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "withdraw_fee_tokens", false, "", typeArgs, state, param, feeTokenMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "withdraw_fee_tokens", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetFee(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_fee", false, "", typeArgs, ref, clock, destChainSelector, receiver, data, tokenAddresses, tokenAmounts, feeToken, extraArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_fee", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) SetDynamicConfig(state bind.Object, param module_common.OwnerCap, feeAggregator string, allowlistAdmin string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "set_dynamic_config", false, "", "", state, param, feeAggregator, allowlistAdmin)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "set_dynamic_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) ApplyDestChainConfigUpdates(state bind.Object, param module_common.OwnerCap, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "apply_dest_chain_config_updates", false, "", "", state, param, destChainSelectors, destChainEnabled, destChainAllowlistEnabled)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "apply_dest_chain_config_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetDestChainConfig(state bind.Object, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_dest_chain_config", false, "", "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_dest_chain_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetAllowedSendersList(state bind.Object, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_allowed_senders_list", false, "", "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_allowed_senders_list", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) ApplyAllowlistUpdates(state bind.Object, param module_common.OwnerCap, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "apply_allowlist_updates", false, "", "", state, param, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "apply_allowlist_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) ApplyAllowlistUpdatesByAdmin(state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "apply_allowlist_updates_by_admin", false, "", "", state, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "apply_allowlist_updates_by_admin", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetOutboundNonce(ref module_common.CCIPObjectRef, destChainSelector uint64, sender string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_outbound_nonce", false, "", "", ref, destChainSelector, sender)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_outbound_nonce", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetStaticConfig(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_static_config", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_static_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetStaticConfigFields(cfg StaticConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_static_config_fields", false, "", "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_static_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetDynamicConfig(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_dynamic_config", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_dynamic_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetDynamicConfigFields(cfg DynamicConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_dynamic_config_fields", false, "", "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_dynamic_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) CcipSend(typeArgs string, ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, receiver []byte, data []byte, tokenParams module_common.TokenParams, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "ccip_send", false, "", typeArgs, ref, state, clock, receiver, data, tokenParams, feeTokenMetadata, feeToken, extraArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "ccip_send", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetCcipPackageId() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_ccip_package_id", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_ccip_package_id", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) Owner(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "owner", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "owner", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) HasPendingTransfer(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "has_pending_transfer", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "has_pending_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) PendingTransferFrom(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "pending_transfer_from", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "pending_transfer_from", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) PendingTransferTo(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "pending_transfer_to", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "pending_transfer_to", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) PendingTransferAccepted(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "pending_transfer_accepted", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "pending_transfer_accepted", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) TransferOwnership(state bind.Object, ownerCap module_common.OwnerCap, newOwner string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "transfer_ownership", false, "", "", state, ownerCap, newOwner)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "transfer_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) AcceptOwnership(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "accept_ownership", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "accept_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "accept_ownership_from_object", false, "", "", state, from)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "accept_ownership_from_object", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
