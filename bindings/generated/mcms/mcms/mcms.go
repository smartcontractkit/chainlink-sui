// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms

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

type IMcms interface {
	SetRoot(state string, clock string, role byte, root []byte, validUntil uint64, chainId *big.Int, multisigAddr string, preOpCount uint64, postOpCount uint64, overridePreviousRoot bool, metadataProof [][]byte, signatures [][]byte) bind.IMethod
	Execute(state string, clock string, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) bind.IMethod
	DispatchTimelockScheduleBatch(timelock string, clock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod
	DispatchTimelockExecuteBatch(timelock string, clock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod
	DispatchTimelockBypasserExecuteBatch(timelockCallbackParams TimelockCallbackParams) bind.IMethod
	DispatchTimelockCancel(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod
	DispatchTimelockUpdateMinDelay(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod
	DispatchTimelockBlockFunction(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod
	DispatchTimelockUnblockFunction(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod
	SetConfig(param module_common.OwnerCap, state string, role byte, chainId *big.Int, signerAddresses [][]byte, signerGroups []byte, groupQuorums []byte, groupParents []byte, clearRoot bool) bind.IMethod
	VerifyMerkleProof(proof [][]byte, root []byte, leaf []byte) bind.IMethod
	ComputeEthMessageHash(root []byte, validUntil uint64) bind.IMethod
	HashOpLeaf(domainSeparator []byte, op Op) bind.IMethod
	SeenSignedHashes(state string, role byte) bind.IMethod
	ExpiringRootAndOpCount(state string, role byte) bind.IMethod
	RootMetadata(state string, role byte) bind.IMethod
	GetRootMetadata(state string, role byte) bind.IMethod
	GetOpCount(state string, role byte) bind.IMethod
	GetRoot(state string, role byte) bind.IMethod
	NumGroups() bind.IMethod
	MaxNumSigners() bind.IMethod
	BypasserRole() bind.IMethod
	CancellerRole() bind.IMethod
	ProposerRole() bind.IMethod
	TimelockRole() bind.IMethod
	IsValidRole(role byte) bind.IMethod
	ZeroHash() bind.IMethod
	TimelockExecuteBatch(timelock string, clock string, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) bind.IMethod
	TimelockGetBlockedFunction(timelock string, index uint64) bind.IMethod
	TimelockIsOperation(timelock string, id []byte) bind.IMethod
	TimelockIsOperationPending(timelock string, id []byte) bind.IMethod
	TimelockIsOperationReady(timelock string, clock string, id []byte) bind.IMethod
	TimelockIsOperationDone(timelock string, id []byte) bind.IMethod
	TimelockGetTimestamp(timelock string, id []byte) bind.IMethod
	TimelockMinDelay(timelock string) bind.IMethod
	TimelockGetBlockedFunctions(timelock string) bind.IMethod
	TimelockGetBlockedFunctionsCount(timelock string) bind.IMethod
	CreateCalls(targets []string, moduleNames []string, functionNames []string, datas [][]byte) bind.IMethod
	HashOperationBatch(calls []Call, predecessor []byte, salt []byte) bind.IMethod
	SignerView(signer Signer) bind.IMethod
	FunctionName(function Function) bind.IMethod
	ModuleName(function Function) bind.IMethod
	Target(function Function) bind.IMethod
	Data(call Call) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type McmsContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IMcms = (*McmsContract)(nil)

func NewMcms(packageID string, client suiclient.ClientImpl) (*McmsContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &McmsContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *McmsContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type MultisigState struct {
	Id        string   `move:"sui::object::UID"`
	Bypasser  Multisig `move:"Multisig"`
	Canceller Multisig `move:"Multisig"`
	Proposer  Multisig `move:"Multisig"`
}

type Multisig struct {
	Role                   byte                   `move:"u8"`
	Config                 Config                 `move:"Config"`
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

type Timelock struct {
	Id       string `move:"sui::object::UID"`
	MinDelay uint64 `move:"u64"`
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

// Functions

func (c *McmsContract) SetRoot(state string, clock string, role byte, root []byte, validUntil uint64, chainId *big.Int, multisigAddr string, preOpCount uint64, postOpCount uint64, overridePreviousRoot bool, metadataProof [][]byte, signatures [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "set_root", false, "", state, clock, role, root, validUntil, chainId, multisigAddr, preOpCount, postOpCount, overridePreviousRoot, metadataProof, signatures)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "set_root", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) Execute(state string, clock string, role byte, chainId *big.Int, multisigAddr string, nonce uint64, to string, moduleName string, functionName string, data []byte, proof [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "execute", false, "", state, clock, role, chainId, multisigAddr, nonce, to, moduleName, functionName, data, proof)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "execute", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) DispatchTimelockScheduleBatch(timelock string, clock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "dispatch_timelock_schedule_batch", false, "", timelock, clock, timelockCallbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "dispatch_timelock_schedule_batch", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) DispatchTimelockExecuteBatch(timelock string, clock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "dispatch_timelock_execute_batch", false, "", timelock, clock, timelockCallbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "dispatch_timelock_execute_batch", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) DispatchTimelockBypasserExecuteBatch(timelockCallbackParams TimelockCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "dispatch_timelock_bypasser_execute_batch", false, "", timelockCallbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "dispatch_timelock_bypasser_execute_batch", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) DispatchTimelockCancel(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "dispatch_timelock_cancel", false, "", timelock, timelockCallbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "dispatch_timelock_cancel", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) DispatchTimelockUpdateMinDelay(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "dispatch_timelock_update_min_delay", false, "", timelock, timelockCallbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "dispatch_timelock_update_min_delay", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) DispatchTimelockBlockFunction(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "dispatch_timelock_block_function", false, "", timelock, timelockCallbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "dispatch_timelock_block_function", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) DispatchTimelockUnblockFunction(timelock string, timelockCallbackParams TimelockCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "dispatch_timelock_unblock_function", false, "", timelock, timelockCallbackParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "dispatch_timelock_unblock_function", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) SetConfig(param module_common.OwnerCap, state string, role byte, chainId *big.Int, signerAddresses [][]byte, signerGroups []byte, groupQuorums []byte, groupParents []byte, clearRoot bool) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "set_config", false, "", param, state, role, chainId, signerAddresses, signerGroups, groupQuorums, groupParents, clearRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "set_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) VerifyMerkleProof(proof [][]byte, root []byte, leaf []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "verify_merkle_proof", false, "", proof, root, leaf)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "verify_merkle_proof", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) ComputeEthMessageHash(root []byte, validUntil uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "compute_eth_message_hash", false, "", root, validUntil)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "compute_eth_message_hash", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) HashOpLeaf(domainSeparator []byte, op Op) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "hash_op_leaf", false, "", domainSeparator, op)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "hash_op_leaf", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) SeenSignedHashes(state string, role byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "seen_signed_hashes", false, "", state, role)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "seen_signed_hashes", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) ExpiringRootAndOpCount(state string, role byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "expiring_root_and_op_count", false, "", state, role)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "expiring_root_and_op_count", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) RootMetadata(state string, role byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "root_metadata", false, "", state, role)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "root_metadata", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) GetRootMetadata(state string, role byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "get_root_metadata", false, "", state, role)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "get_root_metadata", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) GetOpCount(state string, role byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "get_op_count", false, "", state, role)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "get_op_count", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) GetRoot(state string, role byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "get_root", false, "", state, role)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "get_root", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) NumGroups() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "num_groups", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "num_groups", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) MaxNumSigners() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "max_num_signers", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "max_num_signers", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) BypasserRole() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "bypasser_role", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "bypasser_role", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) CancellerRole() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "canceller_role", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "canceller_role", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) ProposerRole() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "proposer_role", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "proposer_role", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockRole() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_role", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_role", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) IsValidRole(role byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "is_valid_role", false, "", role)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "is_valid_role", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) ZeroHash() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "zero_hash", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "zero_hash", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockExecuteBatch(timelock string, clock string, targets []string, moduleNames []string, functionNames []string, datas [][]byte, predecessor []byte, salt []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_execute_batch", false, "", timelock, clock, targets, moduleNames, functionNames, datas, predecessor, salt)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_execute_batch", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockGetBlockedFunction(timelock string, index uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_get_blocked_function", false, "", timelock, index)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_get_blocked_function", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockIsOperation(timelock string, id []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_is_operation", false, "", timelock, id)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_is_operation", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockIsOperationPending(timelock string, id []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_is_operation_pending", false, "", timelock, id)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_is_operation_pending", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockIsOperationReady(timelock string, clock string, id []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_is_operation_ready", false, "", timelock, clock, id)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_is_operation_ready", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockIsOperationDone(timelock string, id []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_is_operation_done", false, "", timelock, id)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_is_operation_done", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockGetTimestamp(timelock string, id []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_get_timestamp", false, "", timelock, id)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_get_timestamp", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockMinDelay(timelock string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_min_delay", false, "", timelock)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_min_delay", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockGetBlockedFunctions(timelock string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_get_blocked_functions", false, "", timelock)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_get_blocked_functions", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) TimelockGetBlockedFunctionsCount(timelock string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "timelock_get_blocked_functions_count", false, "", timelock)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "timelock_get_blocked_functions_count", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) CreateCalls(targets []string, moduleNames []string, functionNames []string, datas [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "create_calls", false, "", targets, moduleNames, functionNames, datas)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "create_calls", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) HashOperationBatch(calls []Call, predecessor []byte, salt []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "hash_operation_batch", false, "", calls, predecessor, salt)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "hash_operation_batch", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) SignerView(signer Signer) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "signer_view", false, "", signer)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "signer_view", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) FunctionName(function Function) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "function_name", false, "", function)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "function_name", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) ModuleName(function Function) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "module_name", false, "", function)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "module_name", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) Target(function Function) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "target", false, "", function)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "target", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsContract) Data(call Call) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "data", false, "", call)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "data", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
