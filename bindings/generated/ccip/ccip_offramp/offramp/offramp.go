// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_offramp

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

type IOfframp interface {
	TypeAndVersion() bind.IMethod
	GetOcr3Base(state string) bind.IMethod
	InitExecute(ref module_common.CCIPObjectRef, state string, clock string, reportContext [][]byte, report []byte) bind.IMethod
	ManuallyInitExecute(ref module_common.CCIPObjectRef, state string, clock string, reportBytes []byte) bind.IMethod
	GetExecutionState(state string, sourceChainSelector uint64, sequenceNumber uint64) bind.IMethod
	SetOcr3Config(state string, param string, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) bind.IMethod
	Commit(ref module_common.CCIPObjectRef, state string, clock string, reportContext [][]byte, report []byte, signatures [][]byte) bind.IMethod
	GetMerkleRoot(state string, root []byte) bind.IMethod
	GetSourceChainConfig(state string, sourceChainSelector uint64) bind.IMethod
	GetSourceChainConfigFields(sourceChainConfig SourceChainConfig) bind.IMethod
	GetAllSourceChainConfigs(state string) bind.IMethod
	GetStaticConfig(state string) bind.IMethod
	GetStaticConfigFields(cfg StaticConfig) bind.IMethod
	GetDynamicConfig(state string) bind.IMethod
	GetDynamicConfigFields(cfg DynamicConfig) bind.IMethod
	SetDynamicConfig(state string, param string, permissionlessExecutionThresholdSeconds uint32) bind.IMethod
	ApplySourceChainConfigUpdates(state string, param string, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) bind.IMethod
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

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type OffRampState struct {
	Id                                      string `move:"sui::object::UID"`
	ChainSelector                           uint64 `move:"u64"`
	PermissionlessExecutionThresholdSeconds uint32 `move:"u32"`
	LatestPriceSequenceNumber               uint64 `move:"u64"`
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
	GasLimit     *big.Int               `move:"u256"`
	TokenAmounts []Any2SuiTokenTransfer `move:"vector<Any2SuiTokenTransfer>"`
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

// Functions

func (c *OfframpContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "type_and_version", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetOcr3Base(state string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_ocr3_base", false, "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_ocr3_base", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) InitExecute(ref module_common.CCIPObjectRef, state string, clock string, reportContext [][]byte, report []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "init_execute", false, "", ref, state, clock, reportContext, report)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "init_execute", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) ManuallyInitExecute(ref module_common.CCIPObjectRef, state string, clock string, reportBytes []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "manually_init_execute", false, "", ref, state, clock, reportBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "manually_init_execute", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetExecutionState(state string, sourceChainSelector uint64, sequenceNumber uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_execution_state", false, "", state, sourceChainSelector, sequenceNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_execution_state", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) SetOcr3Config(state string, param string, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "set_ocr3_config", false, "", state, param, configDigest, ocrPluginType, bigF, isSignatureVerificationEnabled, signers, transmitters)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "set_ocr3_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) Commit(ref module_common.CCIPObjectRef, state string, clock string, reportContext [][]byte, report []byte, signatures [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "commit", false, "", ref, state, clock, reportContext, report, signatures)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "commit", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetMerkleRoot(state string, root []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_merkle_root", false, "", state, root)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_merkle_root", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetSourceChainConfig(state string, sourceChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_source_chain_config", false, "", state, sourceChainSelector)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_source_chain_config_fields", false, "", sourceChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_source_chain_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetAllSourceChainConfigs(state string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_all_source_chain_configs", false, "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_all_source_chain_configs", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetStaticConfig(state string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_static_config", false, "", state)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_static_config_fields", false, "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_static_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) GetDynamicConfig(state string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_dynamic_config", false, "", state)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "get_dynamic_config_fields", false, "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "get_dynamic_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) SetDynamicConfig(state string, param string, permissionlessExecutionThresholdSeconds uint32) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "set_dynamic_config", false, "", state, param, permissionlessExecutionThresholdSeconds)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "set_dynamic_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OfframpContract) ApplySourceChainConfigUpdates(state string, param string, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "offramp", "apply_source_chain_config_updates", false, "", state, param, sourceChainsSelector, sourceChainsIsEnabled, sourceChainsIsRmnVerificationDisabled, sourceChainsOnRamp)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "offramp", "apply_source_chain_config_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
