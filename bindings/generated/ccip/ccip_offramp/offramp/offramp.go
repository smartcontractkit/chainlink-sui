// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_offramp

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

type IOfframp interface {
	TypeAndVersion() bind.IMethod
	Initialize(state bind.Object, param module_common.OwnerCap, feeQuoterCap module_common.FeeQuoterCap, destTransferCap module_common.DestTransferCap, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) bind.IMethod
	GetOcr3Base(state bind.Object) bind.IMethod
	InitExecute(ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) bind.IMethod
	FinishExecute(state bind.Object, receiverParams module_common.ReceiverParams) bind.IMethod
	ManuallyInitExecute(ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, reportBytes []byte) bind.IMethod
	GetExecutionState(state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) bind.IMethod
	SetOcr3Config(state bind.Object, param module_common.OwnerCap, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) bind.IMethod
	ConfigSigners(state module_common.OCRConfig) bind.IMethod
	ConfigTransmitters(state module_common.OCRConfig) bind.IMethod
	Commit(ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) bind.IMethod
	GetMerkleRoot(state bind.Object, root []byte) bind.IMethod
	GetSourceChainConfig(state bind.Object, sourceChainSelector uint64) bind.IMethod
	GetSourceChainConfigFields(sourceChainConfig SourceChainConfig) bind.IMethod
	GetAllSourceChainConfigs(state bind.Object) bind.IMethod
	GetStaticConfig(state bind.Object) bind.IMethod
	GetStaticConfigFields(cfg StaticConfig) bind.IMethod
	GetDynamicConfig(state bind.Object) bind.IMethod
	GetDynamicConfigFields(cfg DynamicConfig) bind.IMethod
	SetDynamicConfig(state bind.Object, param module_common.OwnerCap, permissionlessExecutionThresholdSeconds uint32) bind.IMethod
	ApplySourceChainConfigUpdates(state bind.Object, param module_common.OwnerCap, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) bind.IMethod
	GetCcipPackageId() bind.IMethod
	Owner(state bind.Object) bind.IMethod
	HasPendingTransfer(state bind.Object) bind.IMethod
	PendingTransferFrom(state bind.Object) bind.IMethod
	PendingTransferTo(state bind.Object) bind.IMethod
	PendingTransferAccepted(state bind.Object) bind.IMethod
	TransferOwnership(state bind.Object, ownerCap module_common.OwnerCap, newOwner string) bind.IMethod
	AcceptOwnership(state bind.Object) bind.IMethod
	AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod
	ExecuteOwnershipTransfer(ownerCap module_common.OwnerCap, ownableState module_common.OwnableState, to string) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type OfframpContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IOfframp = (*OfframpContract)(nil)

func NewOfframp(packageID string, client suiclient.ClientImpl) (*OfframpContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &OfframpContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *OfframpContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OffRampState struct {
	Id                                      string                         `move:"sui::object::UID"`
	ChainSelector                           uint64                         `move:"u64"`
	PermissionlessExecutionThresholdSeconds uint32                         `move:"u32"`
	LatestPriceSequenceNumber               uint64                         `move:"u64"`
	FeeQuoterCap                            *module_common.FeeQuoterCap    `move:"0x1::option::Option<FeeQuoterCap>"`
	DestTransferCap                         *module_common.DestTransferCap `move:"0x1::option::Option<ccip::common::DestTransferCap>"`
	OwnableState                            module_common.OwnableState     `move:"OwnableState"`
}

type OffRampStatePointer struct {
	Id             string `move:"sui::object::UID"`
	OffRampStateId string `move:"address"`
	OwnerCapId     string `move:"address"`
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
	Header       RampMessageHeader      `move:"RampMessageHeader"`
	Sender       []byte                 `move:"vector<u8>"`
	Data         []byte                 `move:"vector<u8>"`
	Receiver     string                 `move:"address"`
	GasLimit     uint256.Int            `move:"u256"`
	TokenAmounts []Any2SuiTokenTransfer `move:"vector<Any2SuiTokenTransfer>"`
}

type Any2SuiTokenTransfer struct {
	SourcePoolAddress []byte      `move:"vector<u8>"`
	DestTokenAddress  string      `move:"address"`
	DestGasAmount     uint32      `move:"u32"`
	ExtraData         []byte      `move:"vector<u8>"`
	Amount            uint256.Int `move:"u256"`
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
	SourceToken string      `move:"address"`
	UsdPerToken uint256.Int `move:"u256"`
}

type GasPriceUpdate struct {
	DestChainSelector uint64      `move:"u64"`
	UsdPerUnitGas     uint256.Int `move:"u256"`
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

// Functions

func (c *OfframpContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "type_and_version", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) Initialize(state bind.Object, param module_common.OwnerCap, feeQuoterCap module_common.FeeQuoterCap, destTransferCap module_common.DestTransferCap, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "initialize", false, "", "", state, param, feeQuoterCap, destTransferCap, chainSelector, permissionlessExecutionThresholdSeconds, sourceChainsSelectors, sourceChainsIsEnabled, sourceChainsIsRmnVerificationDisabled, sourceChainsOnRamp)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetOcr3Base(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_ocr3_base", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_ocr3_base", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) InitExecute(ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "init_execute", false, "", "", ref, state, clock, reportContext, report)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "init_execute", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) FinishExecute(state bind.Object, receiverParams module_common.ReceiverParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "finish_execute", false, "", "", state, receiverParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "finish_execute", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) ManuallyInitExecute(ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, reportBytes []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "manually_init_execute", false, "", "", ref, state, clock, reportBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "manually_init_execute", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetExecutionState(state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_execution_state", false, "", "", state, sourceChainSelector, sequenceNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_execution_state", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) SetOcr3Config(state bind.Object, param module_common.OwnerCap, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "set_ocr3_config", false, "", "", state, param, configDigest, ocrPluginType, bigF, isSignatureVerificationEnabled, signers, transmitters)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "set_ocr3_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) ConfigSigners(state module_common.OCRConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "config_signers", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "config_signers", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) ConfigTransmitters(state module_common.OCRConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "config_transmitters", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "config_transmitters", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) Commit(ref module_common.CCIPObjectRef, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "commit", false, "", "", ref, state, clock, reportContext, report, signatures)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "commit", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetMerkleRoot(state bind.Object, root []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_merkle_root", false, "", "", state, root)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_merkle_root", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetSourceChainConfig(state bind.Object, sourceChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_source_chain_config", false, "", "", state, sourceChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_source_chain_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetSourceChainConfigFields(sourceChainConfig SourceChainConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_source_chain_config_fields", false, "", "", sourceChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_source_chain_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetAllSourceChainConfigs(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_all_source_chain_configs", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_all_source_chain_configs", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetStaticConfig(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_static_config", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_static_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetStaticConfigFields(cfg StaticConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_static_config_fields", false, "", "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_static_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetDynamicConfig(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_dynamic_config", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_dynamic_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetDynamicConfigFields(cfg DynamicConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_dynamic_config_fields", false, "", "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_dynamic_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) SetDynamicConfig(state bind.Object, param module_common.OwnerCap, permissionlessExecutionThresholdSeconds uint32) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "set_dynamic_config", false, "", "", state, param, permissionlessExecutionThresholdSeconds)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "set_dynamic_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) ApplySourceChainConfigUpdates(state bind.Object, param module_common.OwnerCap, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "apply_source_chain_config_updates", false, "", "", state, param, sourceChainsSelector, sourceChainsIsEnabled, sourceChainsIsRmnVerificationDisabled, sourceChainsOnRamp)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "apply_source_chain_config_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetCcipPackageId() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_ccip_package_id", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_ccip_package_id", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) Owner(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "owner", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "owner", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) HasPendingTransfer(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "has_pending_transfer", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "has_pending_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) PendingTransferFrom(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "pending_transfer_from", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "pending_transfer_from", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) PendingTransferTo(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "pending_transfer_to", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "pending_transfer_to", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) PendingTransferAccepted(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "pending_transfer_accepted", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "pending_transfer_accepted", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) TransferOwnership(state bind.Object, ownerCap module_common.OwnerCap, newOwner string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "transfer_ownership", false, "", "", state, ownerCap, newOwner)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "transfer_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) AcceptOwnership(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "accept_ownership", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "accept_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "accept_ownership_from_object", false, "", "", state, from)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "accept_ownership_from_object", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) ExecuteOwnershipTransfer(ownerCap module_common.OwnerCap, ownableState module_common.OwnableState, to string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "execute_ownership_transfer", false, "", "", ownerCap, ownableState, to)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "execute_ownership_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
