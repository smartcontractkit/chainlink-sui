// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_rmn_remote

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

var (
	_ = big.NewInt
)

const FunctionInfo = `[{"package":"ccip","module":"rmn_remote","name":"create_curser_cap","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"}]},{"package":"ccip","module":"rmn_remote","name":"create_curser_cap_and_transfer","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"recipient","type":"address"}]},{"package":"ccip","module":"rmn_remote","name":"curse","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"subject","type":"vector<u8>"}]},{"package":"ccip","module":"rmn_remote","name":"curse_multiple","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"subjects","type":"vector<vector<u8>>"}]},{"package":"ccip","module":"rmn_remote","name":"curse_multiple_with_curser_cap","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"curser_cap","type":"CurserCap"},{"name":"subjects","type":"vector<vector<u8>>"}]},{"package":"ccip","module":"rmn_remote","name":"curse_with_curser_cap","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"curser_cap","type":"CurserCap"},{"name":"subject","type":"vector<u8>"}]},{"package":"ccip","module":"rmn_remote","name":"deregister_curser_cap_ids","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"cap_ids","type":"vector<address>"}]},{"package":"ccip","module":"rmn_remote","name":"get_allowed_curser_cap_ids","parameters":[{"name":"ref","type":"CCIPObjectRef"}]},{"package":"ccip","module":"rmn_remote","name":"get_cursed_subjects","parameters":[{"name":"ref","type":"CCIPObjectRef"}]},{"package":"ccip","module":"rmn_remote","name":"get_local_chain_selector","parameters":[{"name":"ref","type":"CCIPObjectRef"}]},{"package":"ccip","module":"rmn_remote","name":"get_report_digest_header","parameters":null},{"package":"ccip","module":"rmn_remote","name":"get_versioned_config","parameters":[{"name":"ref","type":"CCIPObjectRef"}]},{"package":"ccip","module":"rmn_remote","name":"initialize","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"local_chain_selector","type":"u64"}]},{"package":"ccip","module":"rmn_remote","name":"initialize_allowed_curser_caps","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"initial_cap_ids","type":"vector<address>"}]},{"package":"ccip","module":"rmn_remote","name":"is_cursed","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"subject","type":"vector<u8>"}]},{"package":"ccip","module":"rmn_remote","name":"is_cursed_global","parameters":[{"name":"ref","type":"CCIPObjectRef"}]},{"package":"ccip","module":"rmn_remote","name":"is_cursed_u128","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"subject_value","type":"u128"}]},{"package":"ccip","module":"rmn_remote","name":"is_curser_cap_allowed","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"cap_id","type":"address"}]},{"package":"ccip","module":"rmn_remote","name":"mint_and_register_curser_cap","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"fast_registry","type":"FastRegistry"}]},{"package":"ccip","module":"rmn_remote","name":"register_curser_cap","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"fast_registry","type":"FastRegistry"},{"name":"curser_cap","type":"CurserCap"}]},{"package":"ccip","module":"rmn_remote","name":"register_curser_cap_ids","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"cap_ids","type":"vector<address>"}]},{"package":"ccip","module":"rmn_remote","name":"set_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"rmn_home_contract_config_digest","type":"vector<u8>"},{"name":"signer_onchain_public_keys","type":"vector<vector<u8>>"},{"name":"node_indexes","type":"vector<u64>"},{"name":"f_sign","type":"u64"}]},{"package":"ccip","module":"rmn_remote","name":"type_and_version","parameters":null},{"package":"ccip","module":"rmn_remote","name":"uncurse","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"subject","type":"vector<u8>"}]},{"package":"ccip","module":"rmn_remote","name":"uncurse_multiple","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"subjects","type":"vector<vector<u8>>"}]}]`

type IRmnRemote interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, localChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	SetConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*models.SuiTransactionBlockResponse, error)
	GetVersionedConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetLocalChainSelector(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetReportDigestHeader(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Curse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error)
	CurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error)
	CurseWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, curserCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error)
	CurseMultipleWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, curserCap bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error)
	CreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	CreateCurserCapAndTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error)
	MintAndRegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object) (*models.SuiTransactionBlockResponse, error)
	RegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object, curserCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	InitializeAllowedCurserCaps(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, initialCapIds []string) (*models.SuiTransactionBlockResponse, error)
	RegisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, capIds []string) (*models.SuiTransactionBlockResponse, error)
	DeregisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, capIds []string) (*models.SuiTransactionBlockResponse, error)
	IsCurserCapAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, capId string) (*models.SuiTransactionBlockResponse, error)
	GetAllowedCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	Uncurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error)
	UncurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetCursedSubjects(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	IsCursedGlobal(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	IsCursed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error)
	IsCursedU128(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subjectValue *big.Int) (*models.SuiTransactionBlockResponse, error)
	McmsSetConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsUncurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsUncurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCurseWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCurseMultipleWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCreateCurserCapAndTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object, curserCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsMintAndRegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsInitializeAllowedCurserCaps(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsDeregisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IRmnRemoteDevInspect
	Encoder() RmnRemoteEncoder
	Bound() bind.IBoundContract
}

type IRmnRemoteDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetVersionedConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([]any, error)
	GetLocalChainSelector(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (uint64, error)
	GetReportDigestHeader(ctx context.Context, opts *bind.CallOpts) ([]byte, error)
	CreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (bind.Object, error)
	IsCurserCapAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, capId string) (bool, error)
	GetAllowedCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([]string, error)
	GetCursedSubjects(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([][]byte, error)
	IsCursedGlobal(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (bool, error)
	IsCursed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subject []byte) (bool, error)
	IsCursedU128(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subjectValue *big.Int) (bool, error)
	McmsCreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (bind.Object, error)
}

type RmnRemoteEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(ref bind.Object, ownerCap bind.Object, localChainSelector uint64) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	SetConfig(ref bind.Object, ownerCap bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*bind.EncodedCall, error)
	SetConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetVersionedConfig(ref bind.Object) (*bind.EncodedCall, error)
	GetVersionedConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetLocalChainSelector(ref bind.Object) (*bind.EncodedCall, error)
	GetLocalChainSelectorWithArgs(args ...any) (*bind.EncodedCall, error)
	GetReportDigestHeader() (*bind.EncodedCall, error)
	GetReportDigestHeaderWithArgs(args ...any) (*bind.EncodedCall, error)
	Curse(ref bind.Object, ownerCap bind.Object, subject []byte) (*bind.EncodedCall, error)
	CurseWithArgs(args ...any) (*bind.EncodedCall, error)
	CurseMultiple(ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*bind.EncodedCall, error)
	CurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error)
	CurseWithCurserCap(ref bind.Object, curserCap bind.Object, subject []byte) (*bind.EncodedCall, error)
	CurseWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	CurseMultipleWithCurserCap(ref bind.Object, curserCap bind.Object, subjects [][]byte) (*bind.EncodedCall, error)
	CurseMultipleWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateCurserCap(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	CreateCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateCurserCapAndTransfer(ref bind.Object, ownerCap bind.Object, recipient string) (*bind.EncodedCall, error)
	CreateCurserCapAndTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	MintAndRegisterCurserCap(ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object) (*bind.EncodedCall, error)
	MintAndRegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterCurserCap(ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object, curserCap bind.Object) (*bind.EncodedCall, error)
	RegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	InitializeAllowedCurserCaps(ref bind.Object, ownerCap bind.Object, initialCapIds []string) (*bind.EncodedCall, error)
	InitializeAllowedCurserCapsWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterCurserCapIds(ref bind.Object, ownerCap bind.Object, capIds []string) (*bind.EncodedCall, error)
	RegisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error)
	DeregisterCurserCapIds(ref bind.Object, ownerCap bind.Object, capIds []string) (*bind.EncodedCall, error)
	DeregisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsCurserCapAllowed(ref bind.Object, capId string) (*bind.EncodedCall, error)
	IsCurserCapAllowedWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowedCurserCapIds(ref bind.Object) (*bind.EncodedCall, error)
	GetAllowedCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error)
	Uncurse(ref bind.Object, ownerCap bind.Object, subject []byte) (*bind.EncodedCall, error)
	UncurseWithArgs(args ...any) (*bind.EncodedCall, error)
	UncurseMultiple(ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*bind.EncodedCall, error)
	UncurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCursedSubjects(ref bind.Object) (*bind.EncodedCall, error)
	GetCursedSubjectsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsCursedGlobal(ref bind.Object) (*bind.EncodedCall, error)
	IsCursedGlobalWithArgs(args ...any) (*bind.EncodedCall, error)
	IsCursed(ref bind.Object, subject []byte) (*bind.EncodedCall, error)
	IsCursedWithArgs(args ...any) (*bind.EncodedCall, error)
	IsCursedU128(ref bind.Object, subjectValue *big.Int) (*bind.EncodedCall, error)
	IsCursedU128WithArgs(args ...any) (*bind.EncodedCall, error)
	McmsSetConfig(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsSetConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsCurse(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsCurseWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsCurseMultiple(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsCurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsUncurse(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsUncurseWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsUncurseMultiple(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsUncurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsCurseWithCurserCap(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsCurseWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsCurseMultipleWithCurserCap(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsCurseMultipleWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsCreateCurserCap(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsCreateCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsCreateCurserCapAndTransfer(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsCreateCurserCapAndTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterCurserCap(ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object, curserCap bind.Object) (*bind.EncodedCall, error)
	McmsRegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsMintAndRegisterCurserCap(ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsMintAndRegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsInitializeAllowedCurserCaps(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsInitializeAllowedCurserCapsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterCurserCapIds(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRegisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsDeregisterCurserCapIds(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsDeregisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error)
}

type RmnRemoteContract struct {
	*bind.BoundContract
	rmnRemoteEncoder
	devInspect *RmnRemoteDevInspect
}

type RmnRemoteDevInspect struct {
	contract *RmnRemoteContract
}

var _ IRmnRemote = (*RmnRemoteContract)(nil)
var _ IRmnRemoteDevInspect = (*RmnRemoteDevInspect)(nil)

func NewRmnRemote(packageID string, chainClient client.BindingsClient) (IRmnRemote, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "rmn_remote", chainClient)
	if err != nil {
		return nil, err
	}

	c := &RmnRemoteContract{
		BoundContract:    contract,
		rmnRemoteEncoder: rmnRemoteEncoder{BoundContract: contract},
	}
	c.devInspect = &RmnRemoteDevInspect{contract: c}
	return c, nil
}

func (c *RmnRemoteContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *RmnRemoteContract) Encoder() RmnRemoteEncoder {
	return c.rmnRemoteEncoder
}

func (c *RmnRemoteContract) DevInspect() IRmnRemoteDevInspect {
	return c.devInspect
}

type RMNRemoteState struct {
	Id                 string      `move:"sui::object::UID"`
	LocalChainSelector uint64      `move:"u64"`
	Config             Config      `move:"Config"`
	ConfigCount        uint32      `move:"u32"`
	Signers            bind.Object `move:"VecMap<vector<u8>, bool>"`
	CursedSubjects     bind.Object `move:"VecMap<vector<u8>, bool>"`
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

type CurserCap struct {
	Id string `move:"sui::object::UID"`
}

type AllowedCurserCaps struct {
	Id            string      `move:"sui::object::UID"`
	AllowedCapIds bind.Object `move:"VecMap<address, bool>"`
}

type CurserCapRegistered struct {
	CapId string `move:"address"`
}

type CurserCapDeregistered struct {
	CapId string `move:"address"`
}

// TypeAndVersion executes the type_and_version Move function.
func (c *RmnRemoteContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *RmnRemoteContract) Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, localChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.Initialize(ref, ownerCap, localChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetConfig executes the set_config Move function.
func (c *RmnRemoteContract) SetConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.SetConfig(ref, ownerCap, rmnHomeContractConfigDigest, signerOnchainPublicKeys, nodeIndexes, fSign)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetVersionedConfig executes the get_versioned_config Move function.
func (c *RmnRemoteContract) GetVersionedConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.GetVersionedConfig(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetLocalChainSelector executes the get_local_chain_selector Move function.
func (c *RmnRemoteContract) GetLocalChainSelector(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.GetLocalChainSelector(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetReportDigestHeader executes the get_report_digest_header Move function.
func (c *RmnRemoteContract) GetReportDigestHeader(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.GetReportDigestHeader()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Curse executes the curse Move function.
func (c *RmnRemoteContract) Curse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.Curse(ref, ownerCap, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CurseMultiple executes the curse_multiple Move function.
func (c *RmnRemoteContract) CurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.CurseMultiple(ref, ownerCap, subjects)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CurseWithCurserCap executes the curse_with_curser_cap Move function.
func (c *RmnRemoteContract) CurseWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, curserCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.CurseWithCurserCap(ref, curserCap, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CurseMultipleWithCurserCap executes the curse_multiple_with_curser_cap Move function.
func (c *RmnRemoteContract) CurseMultipleWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, curserCap bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.CurseMultipleWithCurserCap(ref, curserCap, subjects)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateCurserCap executes the create_curser_cap Move function.
func (c *RmnRemoteContract) CreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.CreateCurserCap(ref, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateCurserCapAndTransfer executes the create_curser_cap_and_transfer Move function.
func (c *RmnRemoteContract) CreateCurserCapAndTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.CreateCurserCapAndTransfer(ref, ownerCap, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// MintAndRegisterCurserCap executes the mint_and_register_curser_cap Move function.
func (c *RmnRemoteContract) MintAndRegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.MintAndRegisterCurserCap(ref, ownerCap, fastRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterCurserCap executes the register_curser_cap Move function.
func (c *RmnRemoteContract) RegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object, curserCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.RegisterCurserCap(ref, ownerCap, fastRegistry, curserCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// InitializeAllowedCurserCaps executes the initialize_allowed_curser_caps Move function.
func (c *RmnRemoteContract) InitializeAllowedCurserCaps(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, initialCapIds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.InitializeAllowedCurserCaps(ref, ownerCap, initialCapIds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterCurserCapIds executes the register_curser_cap_ids Move function.
func (c *RmnRemoteContract) RegisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, capIds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.RegisterCurserCapIds(ref, ownerCap, capIds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DeregisterCurserCapIds executes the deregister_curser_cap_ids Move function.
func (c *RmnRemoteContract) DeregisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, capIds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.DeregisterCurserCapIds(ref, ownerCap, capIds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsCurserCapAllowed executes the is_curser_cap_allowed Move function.
func (c *RmnRemoteContract) IsCurserCapAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, capId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.IsCurserCapAllowed(ref, capId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowedCurserCapIds executes the get_allowed_curser_cap_ids Move function.
func (c *RmnRemoteContract) GetAllowedCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.GetAllowedCurserCapIds(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Uncurse executes the uncurse Move function.
func (c *RmnRemoteContract) Uncurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.Uncurse(ref, ownerCap, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UncurseMultiple executes the uncurse_multiple Move function.
func (c *RmnRemoteContract) UncurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.UncurseMultiple(ref, ownerCap, subjects)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCursedSubjects executes the get_cursed_subjects Move function.
func (c *RmnRemoteContract) GetCursedSubjects(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.GetCursedSubjects(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsCursedGlobal executes the is_cursed_global Move function.
func (c *RmnRemoteContract) IsCursedGlobal(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.IsCursedGlobal(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsCursed executes the is_cursed Move function.
func (c *RmnRemoteContract) IsCursed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.IsCursed(ref, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsCursedU128 executes the is_cursed_u128 Move function.
func (c *RmnRemoteContract) IsCursedU128(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subjectValue *big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.IsCursedU128(ref, subjectValue)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetConfig executes the mcms_set_config Move function.
func (c *RmnRemoteContract) McmsSetConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsSetConfig(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsCurse executes the mcms_curse Move function.
func (c *RmnRemoteContract) McmsCurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsCurse(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsCurseMultiple executes the mcms_curse_multiple Move function.
func (c *RmnRemoteContract) McmsCurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsCurseMultiple(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsUncurse executes the mcms_uncurse Move function.
func (c *RmnRemoteContract) McmsUncurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsUncurse(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsUncurseMultiple executes the mcms_uncurse_multiple Move function.
func (c *RmnRemoteContract) McmsUncurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsUncurseMultiple(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsCurseWithCurserCap executes the mcms_curse_with_curser_cap Move function.
func (c *RmnRemoteContract) McmsCurseWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsCurseWithCurserCap(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsCurseMultipleWithCurserCap executes the mcms_curse_multiple_with_curser_cap Move function.
func (c *RmnRemoteContract) McmsCurseMultipleWithCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsCurseMultipleWithCurserCap(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsCreateCurserCap executes the mcms_create_curser_cap Move function.
func (c *RmnRemoteContract) McmsCreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsCreateCurserCap(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsCreateCurserCapAndTransfer executes the mcms_create_curser_cap_and_transfer Move function.
func (c *RmnRemoteContract) McmsCreateCurserCapAndTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsCreateCurserCapAndTransfer(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterCurserCap executes the mcms_register_curser_cap Move function.
func (c *RmnRemoteContract) McmsRegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object, curserCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsRegisterCurserCap(ref, slowRegistry, fastRegistry, params, curserCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsMintAndRegisterCurserCap executes the mcms_mint_and_register_curser_cap Move function.
func (c *RmnRemoteContract) McmsMintAndRegisterCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsMintAndRegisterCurserCap(ref, slowRegistry, fastRegistry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsInitializeAllowedCurserCaps executes the mcms_initialize_allowed_curser_caps Move function.
func (c *RmnRemoteContract) McmsInitializeAllowedCurserCaps(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsInitializeAllowedCurserCaps(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterCurserCapIds executes the mcms_register_curser_cap_ids Move function.
func (c *RmnRemoteContract) McmsRegisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsRegisterCurserCapIds(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsDeregisterCurserCapIds executes the mcms_deregister_curser_cap_ids Move function.
func (c *RmnRemoteContract) McmsDeregisterCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.McmsDeregisterCurserCapIds(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *RmnRemoteDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.rmnRemoteEncoder.TypeAndVersion()
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
	var result string
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return "", fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// GetVersionedConfig executes the get_versioned_config Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u32
//	[1]: Config
func (d *RmnRemoteDevInspect) GetVersionedConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([]any, error) {
	encoded, err := d.contract.rmnRemoteEncoder.GetVersionedConfig(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) != 2 {
		return nil, fmt.Errorf("expected 2 return values, got %d", len(results))
	}
	decoded := make([]any, 2)
	var ret0 uint32
	if err := bind.DecodeJSONReturn(results[0], &ret0); err != nil {
		return nil, fmt.Errorf("failed to decode return value 0: %w", err)
	}
	decoded[0] = ret0
	var ret1 Config
	if err := bind.DecodeJSONReturn(results[1], &ret1); err != nil {
		return nil, fmt.Errorf("failed to decode return value 1: %w", err)
	}
	decoded[1] = ret1
	return decoded, nil
}

// GetLocalChainSelector executes the get_local_chain_selector Move function using DevInspect to get return values.
//
// Returns: u64
func (d *RmnRemoteDevInspect) GetLocalChainSelector(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (uint64, error) {
	encoded, err := d.contract.rmnRemoteEncoder.GetLocalChainSelector(ref)
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
	var result uint64
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return 0, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// GetReportDigestHeader executes the get_report_digest_header Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *RmnRemoteDevInspect) GetReportDigestHeader(ctx context.Context, opts *bind.CallOpts) ([]byte, error) {
	encoded, err := d.contract.rmnRemoteEncoder.GetReportDigestHeader()
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
	var result []byte
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return nil, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// CreateCurserCap executes the create_curser_cap Move function using DevInspect to get return values.
//
// Returns: CurserCap
func (d *RmnRemoteDevInspect) CreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (bind.Object, error) {
	encoded, err := d.contract.rmnRemoteEncoder.CreateCurserCap(ref, ownerCap)
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
	var result bind.Object
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return bind.Object{}, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// IsCurserCapAllowed executes the is_curser_cap_allowed Move function using DevInspect to get return values.
//
// Returns: bool
func (d *RmnRemoteDevInspect) IsCurserCapAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, capId string) (bool, error) {
	encoded, err := d.contract.rmnRemoteEncoder.IsCurserCapAllowed(ref, capId)
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
	var result bool
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return false, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// GetAllowedCurserCapIds executes the get_allowed_curser_cap_ids Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *RmnRemoteDevInspect) GetAllowedCurserCapIds(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([]string, error) {
	encoded, err := d.contract.rmnRemoteEncoder.GetAllowedCurserCapIds(ref)
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
	var result []string
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return nil, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// GetCursedSubjects executes the get_cursed_subjects Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *RmnRemoteDevInspect) GetCursedSubjects(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([][]byte, error) {
	encoded, err := d.contract.rmnRemoteEncoder.GetCursedSubjects(ref)
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
	var result [][]byte
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return nil, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// IsCursedGlobal executes the is_cursed_global Move function using DevInspect to get return values.
//
// Returns: bool
func (d *RmnRemoteDevInspect) IsCursedGlobal(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (bool, error) {
	encoded, err := d.contract.rmnRemoteEncoder.IsCursedGlobal(ref)
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
	var result bool
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return false, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// IsCursed executes the is_cursed Move function using DevInspect to get return values.
//
// Returns: bool
func (d *RmnRemoteDevInspect) IsCursed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subject []byte) (bool, error) {
	encoded, err := d.contract.rmnRemoteEncoder.IsCursed(ref, subject)
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
	var result bool
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return false, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// IsCursedU128 executes the is_cursed_u128 Move function using DevInspect to get return values.
//
// Returns: bool
func (d *RmnRemoteDevInspect) IsCursedU128(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subjectValue *big.Int) (bool, error) {
	encoded, err := d.contract.rmnRemoteEncoder.IsCursedU128(ref, subjectValue)
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
	var result bool
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return false, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// McmsCreateCurserCap executes the mcms_create_curser_cap Move function using DevInspect to get return values.
//
// Returns: CurserCap
func (d *RmnRemoteDevInspect) McmsCreateCurserCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (bind.Object, error) {
	encoded, err := d.contract.rmnRemoteEncoder.McmsCreateCurserCap(ref, registry, params)
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
	var result bind.Object
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return bind.Object{}, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

type rmnRemoteEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c rmnRemoteEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c rmnRemoteEncoder) Initialize(ref bind.Object, ownerCap bind.Object, localChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u64",
	}, []any{
		ref,
		ownerCap,
		localChainSelector,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetConfig encodes a call to the set_config Move function.
func (c rmnRemoteEncoder) SetConfig(ref bind.Object, ownerCap bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_config", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<u8>",
		"vector<vector<u8>>",
		"vector<u64>",
		"u64",
	}, []any{
		ref,
		ownerCap,
		rmnHomeContractConfigDigest,
		signerOnchainPublicKeys,
		nodeIndexes,
		fSign,
	}, nil)
}

// SetConfigWithArgs encodes a call to the set_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) SetConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<u8>",
		"vector<vector<u8>>",
		"vector<u64>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetVersionedConfig encodes a call to the get_versioned_config Move function.
func (c rmnRemoteEncoder) GetVersionedConfig(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_versioned_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"u32",
		"ccip::rmn_remote::Config",
	})
}

// GetVersionedConfigWithArgs encodes a call to the get_versioned_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) GetVersionedConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_versioned_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u32",
		"ccip::rmn_remote::Config",
	})
}

// GetLocalChainSelector encodes a call to the get_local_chain_selector Move function.
func (c rmnRemoteEncoder) GetLocalChainSelector(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_local_chain_selector", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"u64",
	})
}

// GetLocalChainSelectorWithArgs encodes a call to the get_local_chain_selector Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) GetLocalChainSelectorWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_local_chain_selector", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetReportDigestHeader encodes a call to the get_report_digest_header Move function.
func (c rmnRemoteEncoder) GetReportDigestHeader() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_report_digest_header", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"vector<u8>",
	})
}

// GetReportDigestHeaderWithArgs encodes a call to the get_report_digest_header Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) GetReportDigestHeaderWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_report_digest_header", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// Curse encodes a call to the curse Move function.
func (c rmnRemoteEncoder) Curse(ref bind.Object, ownerCap bind.Object, subject []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<u8>",
	}, []any{
		ref,
		ownerCap,
		subject,
	}, nil)
}

// CurseWithArgs encodes a call to the curse Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) CurseWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CurseMultiple encodes a call to the curse_multiple Move function.
func (c rmnRemoteEncoder) CurseMultiple(ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse_multiple", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<vector<u8>>",
	}, []any{
		ref,
		ownerCap,
		subjects,
	}, nil)
}

// CurseMultipleWithArgs encodes a call to the curse_multiple Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) CurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse_multiple", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CurseWithCurserCap encodes a call to the curse_with_curser_cap Move function.
func (c rmnRemoteEncoder) CurseWithCurserCap(ref bind.Object, curserCap bind.Object, subject []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse_with_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&CurserCap",
		"vector<u8>",
	}, []any{
		ref,
		curserCap,
		subject,
	}, nil)
}

// CurseWithCurserCapWithArgs encodes a call to the curse_with_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) CurseWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&CurserCap",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse_with_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CurseMultipleWithCurserCap encodes a call to the curse_multiple_with_curser_cap Move function.
func (c rmnRemoteEncoder) CurseMultipleWithCurserCap(ref bind.Object, curserCap bind.Object, subjects [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse_multiple_with_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&CurserCap",
		"vector<vector<u8>>",
	}, []any{
		ref,
		curserCap,
		subjects,
	}, nil)
}

// CurseMultipleWithCurserCapWithArgs encodes a call to the curse_multiple_with_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) CurseMultipleWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&CurserCap",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse_multiple_with_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CreateCurserCap encodes a call to the create_curser_cap Move function.
func (c rmnRemoteEncoder) CreateCurserCap(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}, []any{
		ref,
		ownerCap,
	}, []string{
		"ccip::rmn_remote::CurserCap",
	})
}

// CreateCurserCapWithArgs encodes a call to the create_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) CreateCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_curser_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::rmn_remote::CurserCap",
	})
}

// CreateCurserCapAndTransfer encodes a call to the create_curser_cap_and_transfer Move function.
func (c rmnRemoteEncoder) CreateCurserCapAndTransfer(ref bind.Object, ownerCap bind.Object, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_curser_cap_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"address",
	}, []any{
		ref,
		ownerCap,
		recipient,
	}, nil)
}

// CreateCurserCapAndTransferWithArgs encodes a call to the create_curser_cap_and_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) CreateCurserCapAndTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_curser_cap_and_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// MintAndRegisterCurserCap encodes a call to the mint_and_register_curser_cap Move function.
func (c rmnRemoteEncoder) MintAndRegisterCurserCap(ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_register_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"&mut FastRegistry",
	}, []any{
		ref,
		ownerCap,
		fastRegistry,
	}, nil)
}

// MintAndRegisterCurserCapWithArgs encodes a call to the mint_and_register_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) MintAndRegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"&mut FastRegistry",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_register_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RegisterCurserCap encodes a call to the register_curser_cap Move function.
func (c rmnRemoteEncoder) RegisterCurserCap(ref bind.Object, ownerCap bind.Object, fastRegistry bind.Object, curserCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"&mut FastRegistry",
		"ccip::rmn_remote::CurserCap",
	}, []any{
		ref,
		ownerCap,
		fastRegistry,
		curserCap,
	}, nil)
}

// RegisterCurserCapWithArgs encodes a call to the register_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) RegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"&mut FastRegistry",
		"ccip::rmn_remote::CurserCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// InitializeAllowedCurserCaps encodes a call to the initialize_allowed_curser_caps Move function.
func (c rmnRemoteEncoder) InitializeAllowedCurserCaps(ref bind.Object, ownerCap bind.Object, initialCapIds []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize_allowed_curser_caps", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
	}, []any{
		ref,
		ownerCap,
		initialCapIds,
	}, nil)
}

// InitializeAllowedCurserCapsWithArgs encodes a call to the initialize_allowed_curser_caps Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) InitializeAllowedCurserCapsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize_allowed_curser_caps", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RegisterCurserCapIds encodes a call to the register_curser_cap_ids Move function.
func (c rmnRemoteEncoder) RegisterCurserCapIds(ref bind.Object, ownerCap bind.Object, capIds []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_curser_cap_ids", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
	}, []any{
		ref,
		ownerCap,
		capIds,
	}, nil)
}

// RegisterCurserCapIdsWithArgs encodes a call to the register_curser_cap_ids Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) RegisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_curser_cap_ids", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DeregisterCurserCapIds encodes a call to the deregister_curser_cap_ids Move function.
func (c rmnRemoteEncoder) DeregisterCurserCapIds(ref bind.Object, ownerCap bind.Object, capIds []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("deregister_curser_cap_ids", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
	}, []any{
		ref,
		ownerCap,
		capIds,
	}, nil)
}

// DeregisterCurserCapIdsWithArgs encodes a call to the deregister_curser_cap_ids Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) DeregisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("deregister_curser_cap_ids", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IsCurserCapAllowed encodes a call to the is_curser_cap_allowed Move function.
func (c rmnRemoteEncoder) IsCurserCapAllowed(ref bind.Object, capId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_curser_cap_allowed", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		capId,
	}, []string{
		"bool",
	})
}

// IsCurserCapAllowedWithArgs encodes a call to the is_curser_cap_allowed Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) IsCurserCapAllowedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_curser_cap_allowed", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetAllowedCurserCapIds encodes a call to the get_allowed_curser_cap_ids Move function.
func (c rmnRemoteEncoder) GetAllowedCurserCapIds(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_curser_cap_ids", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"vector<address>",
	})
}

// GetAllowedCurserCapIdsWithArgs encodes a call to the get_allowed_curser_cap_ids Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) GetAllowedCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_curser_cap_ids", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// Uncurse encodes a call to the uncurse Move function.
func (c rmnRemoteEncoder) Uncurse(ref bind.Object, ownerCap bind.Object, subject []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("uncurse", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<u8>",
	}, []any{
		ref,
		ownerCap,
		subject,
	}, nil)
}

// UncurseWithArgs encodes a call to the uncurse Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) UncurseWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("uncurse", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// UncurseMultiple encodes a call to the uncurse_multiple Move function.
func (c rmnRemoteEncoder) UncurseMultiple(ref bind.Object, ownerCap bind.Object, subjects [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("uncurse_multiple", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<vector<u8>>",
	}, []any{
		ref,
		ownerCap,
		subjects,
	}, nil)
}

// UncurseMultipleWithArgs encodes a call to the uncurse_multiple Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) UncurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("uncurse_multiple", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCursedSubjects encodes a call to the get_cursed_subjects Move function.
func (c rmnRemoteEncoder) GetCursedSubjects(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_cursed_subjects", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"vector<vector<u8>>",
	})
}

// GetCursedSubjectsWithArgs encodes a call to the get_cursed_subjects Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) GetCursedSubjectsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_cursed_subjects", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
	})
}

// IsCursedGlobal encodes a call to the is_cursed_global Move function.
func (c rmnRemoteEncoder) IsCursedGlobal(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_cursed_global", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"bool",
	})
}

// IsCursedGlobalWithArgs encodes a call to the is_cursed_global Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) IsCursedGlobalWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_cursed_global", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// IsCursed encodes a call to the is_cursed Move function.
func (c rmnRemoteEncoder) IsCursed(ref bind.Object, subject []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_cursed", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"vector<u8>",
	}, []any{
		ref,
		subject,
	}, []string{
		"bool",
	})
}

// IsCursedWithArgs encodes a call to the is_cursed Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) IsCursedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_cursed", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// IsCursedU128 encodes a call to the is_cursed_u128 Move function.
func (c rmnRemoteEncoder) IsCursedU128(ref bind.Object, subjectValue *big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_cursed_u128", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u128",
	}, []any{
		ref,
		subjectValue,
	}, []string{
		"bool",
	})
}

// IsCursedU128WithArgs encodes a call to the is_cursed_u128 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) IsCursedU128WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u128",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_cursed_u128", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// McmsSetConfig encodes a call to the mcms_set_config Move function.
func (c rmnRemoteEncoder) McmsSetConfig(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_config", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsSetConfigWithArgs encodes a call to the mcms_set_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsSetConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsCurse encodes a call to the mcms_curse Move function.
func (c rmnRemoteEncoder) McmsCurse(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsCurseWithArgs encodes a call to the mcms_curse Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsCurseWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsCurseMultiple encodes a call to the mcms_curse_multiple Move function.
func (c rmnRemoteEncoder) McmsCurseMultiple(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse_multiple", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsCurseMultipleWithArgs encodes a call to the mcms_curse_multiple Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsCurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse_multiple", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsUncurse encodes a call to the mcms_uncurse Move function.
func (c rmnRemoteEncoder) McmsUncurse(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_uncurse", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsUncurseWithArgs encodes a call to the mcms_uncurse Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsUncurseWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_uncurse", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsUncurseMultiple encodes a call to the mcms_uncurse_multiple Move function.
func (c rmnRemoteEncoder) McmsUncurseMultiple(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_uncurse_multiple", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsUncurseMultipleWithArgs encodes a call to the mcms_uncurse_multiple Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsUncurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_uncurse_multiple", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsCurseWithCurserCap encodes a call to the mcms_curse_with_curser_cap Move function.
func (c rmnRemoteEncoder) McmsCurseWithCurserCap(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse_with_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut FastRegistry",
		"FastExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsCurseWithCurserCapWithArgs encodes a call to the mcms_curse_with_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsCurseWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut FastRegistry",
		"FastExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse_with_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsCurseMultipleWithCurserCap encodes a call to the mcms_curse_multiple_with_curser_cap Move function.
func (c rmnRemoteEncoder) McmsCurseMultipleWithCurserCap(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse_multiple_with_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut FastRegistry",
		"FastExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsCurseMultipleWithCurserCapWithArgs encodes a call to the mcms_curse_multiple_with_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsCurseMultipleWithCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut FastRegistry",
		"FastExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_curse_multiple_with_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsCreateCurserCap encodes a call to the mcms_create_curser_cap Move function.
func (c rmnRemoteEncoder) McmsCreateCurserCap(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_create_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, []string{
		"ccip::rmn_remote::CurserCap",
	})
}

// McmsCreateCurserCapWithArgs encodes a call to the mcms_create_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsCreateCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_create_curser_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::rmn_remote::CurserCap",
	})
}

// McmsCreateCurserCapAndTransfer encodes a call to the mcms_create_curser_cap_and_transfer Move function.
func (c rmnRemoteEncoder) McmsCreateCurserCapAndTransfer(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_create_curser_cap_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsCreateCurserCapAndTransferWithArgs encodes a call to the mcms_create_curser_cap_and_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsCreateCurserCapAndTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_create_curser_cap_and_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRegisterCurserCap encodes a call to the mcms_register_curser_cap Move function.
func (c rmnRemoteEncoder) McmsRegisterCurserCap(ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object, curserCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"&mut FastRegistry",
		"ExecutingCallbackParams",
		"ccip::rmn_remote::CurserCap",
	}, []any{
		ref,
		slowRegistry,
		fastRegistry,
		params,
		curserCap,
	}, nil)
}

// McmsRegisterCurserCapWithArgs encodes a call to the mcms_register_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsRegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"&mut FastRegistry",
		"ExecutingCallbackParams",
		"ccip::rmn_remote::CurserCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsMintAndRegisterCurserCap encodes a call to the mcms_mint_and_register_curser_cap Move function.
func (c rmnRemoteEncoder) McmsMintAndRegisterCurserCap(ref bind.Object, slowRegistry bind.Object, fastRegistry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_mint_and_register_curser_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"&mut FastRegistry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		slowRegistry,
		fastRegistry,
		params,
	}, nil)
}

// McmsMintAndRegisterCurserCapWithArgs encodes a call to the mcms_mint_and_register_curser_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsMintAndRegisterCurserCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"&mut FastRegistry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_mint_and_register_curser_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsInitializeAllowedCurserCaps encodes a call to the mcms_initialize_allowed_curser_caps Move function.
func (c rmnRemoteEncoder) McmsInitializeAllowedCurserCaps(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_initialize_allowed_curser_caps", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsInitializeAllowedCurserCapsWithArgs encodes a call to the mcms_initialize_allowed_curser_caps Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsInitializeAllowedCurserCapsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_initialize_allowed_curser_caps", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRegisterCurserCapIds encodes a call to the mcms_register_curser_cap_ids Move function.
func (c rmnRemoteEncoder) McmsRegisterCurserCapIds(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_curser_cap_ids", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsRegisterCurserCapIdsWithArgs encodes a call to the mcms_register_curser_cap_ids Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsRegisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_curser_cap_ids", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsDeregisterCurserCapIds encodes a call to the mcms_deregister_curser_cap_ids Move function.
func (c rmnRemoteEncoder) McmsDeregisterCurserCapIds(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_deregister_curser_cap_ids", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsDeregisterCurserCapIdsWithArgs encodes a call to the mcms_deregister_curser_cap_ids Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) McmsDeregisterCurserCapIdsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_deregister_curser_cap_ids", typeArgsList, typeParamsList, expectedParams, args, nil)
}
