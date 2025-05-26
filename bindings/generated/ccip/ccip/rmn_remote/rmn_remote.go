// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_rmn_remote

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

type IRmnRemote interface {
	TypeAndVersion() bind.IMethod
	Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap, localChainSelector uint64) bind.IMethod
	Verify(ref module_common.CCIPObjectRef, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) bind.IMethod
	GetArm() bind.IMethod
	SetConfig(ref module_common.CCIPObjectRef, param module_common.OwnerCap, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) bind.IMethod
	GetVersionedConfig(ref module_common.CCIPObjectRef) bind.IMethod
	GetVersionedConfigFields(vc VersionedConfig) bind.IMethod
	GetLocalChainSelector(ref module_common.CCIPObjectRef) bind.IMethod
	GetReportDigestHeader() bind.IMethod
	Curse(ref module_common.CCIPObjectRef, ownerCap module_common.OwnerCap, subject []byte) bind.IMethod
	CurseMultiple(ref module_common.CCIPObjectRef, param module_common.OwnerCap, subjects [][]byte) bind.IMethod
	Uncurse(ref module_common.CCIPObjectRef, ownerCap module_common.OwnerCap, subject []byte) bind.IMethod
	UncurseMultiple(ref module_common.CCIPObjectRef, param module_common.OwnerCap, subjects [][]byte) bind.IMethod
	GetCursedSubjects(ref module_common.CCIPObjectRef) bind.IMethod
	IsCursedGlobal(ref module_common.CCIPObjectRef) bind.IMethod
	IsCursed(ref module_common.CCIPObjectRef, subject []byte) bind.IMethod
	IsCursedU128(ref module_common.CCIPObjectRef, subjectValue *big.Int) bind.IMethod
	GetActiveSigners(ref module_common.CCIPObjectRef) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type RmnRemoteContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IRmnRemote = (*RmnRemoteContract)(nil)

func NewRmnRemote(packageID string, client suiclient.ClientImpl) (*RmnRemoteContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &RmnRemoteContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *RmnRemoteContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type RMNRemoteState struct {
	Id                 string `move:"sui::object::UID"`
	LocalChainSelector uint64 `move:"u64"`
	Config             Config `move:"Config"`
	ConfigCount        uint32 `move:"u32"`
}

type Config struct {
	RmnHomeContractConfigDigest []byte   `move:"vector<u8>"`
	Signers                     []Signer `move:"vector<Signer>"`
	FSign                       uint64   `move:"u64"`
}

type Signer struct {
	OnchainPublicKey []byte `move:"vector<u8>"`
	NodeIndex        uint64 `move:"u64"`
}

type Report struct {
	DestChainSelector           uint64       `move:"u64"`
	RmnRemoteContractAddress    string       `move:"address"`
	OffRampAddress              string       `move:"address"`
	RmnHomeContractConfigDigest []byte       `move:"vector<u8>"`
	MerkleRoots                 []MerkleRoot `move:"vector<MerkleRoot>"`
}

type MerkleRoot struct {
	SourceChainSelector uint64 `move:"u64"`
	OnRampAddress       []byte `move:"vector<u8>"`
	MinSeqNr            uint64 `move:"u64"`
	MaxSeqNr            uint64 `move:"u64"`
	MerkleRoot          []byte `move:"vector<u8>"`
}

type VersionedConfig struct {
	Version uint32 `move:"u32"`
	Config  Config `move:"Config"`
}

type ConfigSet struct {
	Version uint32 `move:"u32"`
	Config  Config `move:"Config"`
}

type Cursed struct {
	Subjects [][]byte `move:"vector<vector<u8>>"`
}

type Uncursed struct {
	Subjects [][]byte `move:"vector<vector<u8>>"`
}

// Functions

func (c *RmnRemoteContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "type_and_version", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap, localChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "initialize", false, "", ref, param, localChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) Verify(ref module_common.CCIPObjectRef, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "verify", false, "", ref, merkleRootSourceChainSelectors, merkleRootOnRampAddresses, merkleRootMinSeqNrs, merkleRootMaxSeqNrs, merkleRootValues, signatures)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "verify", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) GetArm() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "get_arm", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "get_arm", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) SetConfig(ref module_common.CCIPObjectRef, param module_common.OwnerCap, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "set_config", false, "", ref, param, rmnHomeContractConfigDigest, signerOnchainPublicKeys, nodeIndexes, fSign)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "set_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) GetVersionedConfig(ref module_common.CCIPObjectRef) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "get_versioned_config", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "get_versioned_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) GetVersionedConfigFields(vc VersionedConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "get_versioned_config_fields", false, "", vc)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "get_versioned_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) GetLocalChainSelector(ref module_common.CCIPObjectRef) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "get_local_chain_selector", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "get_local_chain_selector", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) GetReportDigestHeader() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "get_report_digest_header", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "get_report_digest_header", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) Curse(ref module_common.CCIPObjectRef, ownerCap module_common.OwnerCap, subject []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "curse", false, "", ref, ownerCap, subject)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "curse", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) CurseMultiple(ref module_common.CCIPObjectRef, param module_common.OwnerCap, subjects [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "curse_multiple", false, "", ref, param, subjects)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "curse_multiple", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) Uncurse(ref module_common.CCIPObjectRef, ownerCap module_common.OwnerCap, subject []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "uncurse", false, "", ref, ownerCap, subject)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "uncurse", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) UncurseMultiple(ref module_common.CCIPObjectRef, param module_common.OwnerCap, subjects [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "uncurse_multiple", false, "", ref, param, subjects)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "uncurse_multiple", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) GetCursedSubjects(ref module_common.CCIPObjectRef) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "get_cursed_subjects", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "get_cursed_subjects", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) IsCursedGlobal(ref module_common.CCIPObjectRef) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "is_cursed_global", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "is_cursed_global", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) IsCursed(ref module_common.CCIPObjectRef, subject []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "is_cursed", false, "", ref, subject)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "is_cursed", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) IsCursedU128(ref module_common.CCIPObjectRef, subjectValue *big.Int) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "is_cursed_u128", false, "", ref, subjectValue)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "is_cursed_u128", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RmnRemoteContract) GetActiveSigners(ref module_common.CCIPObjectRef) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "rmn_remote", "get_active_signers", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "rmn_remote", "get_active_signers", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
