// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_onramp

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_common "github.com/smartcontractkit/chainlink-sui/bindings/common"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
)

type IOnramp interface {
	TypeAndVersion() bind.IMethod
	IsChainSupported(state string, destChainSelector uint64) bind.IMethod
	GetExpectedNextSequenceNumber(state string, destChainSelector uint64) bind.IMethod
	SetDynamicConfig(state string, param string, feeAggregator string, allowlistAdmin string) bind.IMethod
	ApplyDestChainConfigUpdates(state string, param string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) bind.IMethod
	GetDestChainConfig(state string, destChainSelector uint64) bind.IMethod
	GetAllowedSendersList(state string, destChainSelector uint64) bind.IMethod
	ApplyAllowlistUpdates(state string, param string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod
	ApplyAllowlistUpdatesByAdmin(state string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod
	GetOutboundNonce(ref module_common.CCIPObjectRef, destChainSelector uint64, sender string) bind.IMethod
	GetStaticConfig(state string) bind.IMethod
	GetStaticConfigFields(cfg StaticConfig) bind.IMethod
	GetDynamicConfig(state string) bind.IMethod
	GetDynamicConfigFields(cfg DynamicConfig) bind.IMethod
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

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type OnRampState struct {
	Id             string `move:"sui::object::UID"`
	ChainSelector  uint64 `move:"u64"`
	FeeAggregator  string `move:"address"`
	AllowlistAdmin string `move:"address"`
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

// Functions

func (c *OnrampContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "type_and_version", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) IsChainSupported(state string, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "is_chain_supported", false, "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "is_chain_supported", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetExpectedNextSequenceNumber(state string, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_expected_next_sequence_number", false, "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_expected_next_sequence_number", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) SetDynamicConfig(state string, param string, feeAggregator string, allowlistAdmin string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "set_dynamic_config", false, "", state, param, feeAggregator, allowlistAdmin)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "set_dynamic_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) ApplyDestChainConfigUpdates(state string, param string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "apply_dest_chain_config_updates", false, "", state, param, destChainSelectors, destChainEnabled, destChainAllowlistEnabled)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "apply_dest_chain_config_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetDestChainConfig(state string, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_dest_chain_config", false, "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_dest_chain_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetAllowedSendersList(state string, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_allowed_senders_list", false, "", state, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_allowed_senders_list", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) ApplyAllowlistUpdates(state string, param string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "apply_allowlist_updates", false, "", state, param, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "apply_allowlist_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) ApplyAllowlistUpdatesByAdmin(state string, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "apply_allowlist_updates_by_admin", false, "", state, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_outbound_nonce", false, "", ref, destChainSelector, sender)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_outbound_nonce", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetStaticConfig(state string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_static_config", false, "", state)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_static_config_fields", false, "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_static_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OnrampContract) GetDynamicConfig(state string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_dynamic_config", false, "", state)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "onramp", "get_dynamic_config_fields", false, "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "onramp", "get_dynamic_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
