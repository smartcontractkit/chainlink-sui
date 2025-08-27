// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_rmn_remote

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

type IRmnRemote interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetArm(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, localChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	Verify(ctx context.Context, opts *bind.CallOpts, ref bind.Object, offRampStateAddress string, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error)
	SetConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*models.SuiTransactionBlockResponse, error)
	GetVersionedConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetLocalChainSelector(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetReportDigestHeader(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Curse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error)
	CurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error)
	Uncurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error)
	UncurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetCursedSubjects(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	IsCursedGlobal(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	IsCursed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subject []byte) (*models.SuiTransactionBlockResponse, error)
	IsCursedU128(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subjectValue *big.Int) (*models.SuiTransactionBlockResponse, error)
	McmsSetConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsCurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsUncurse(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsUncurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IRmnRemoteDevInspect
	Encoder() RmnRemoteEncoder
}

type IRmnRemoteDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetArm(ctx context.Context, opts *bind.CallOpts) (string, error)
	Verify(ctx context.Context, opts *bind.CallOpts, ref bind.Object, offRampStateAddress string, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) (bool, error)
	GetVersionedConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([]any, error)
	GetLocalChainSelector(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (uint64, error)
	GetReportDigestHeader(ctx context.Context, opts *bind.CallOpts) ([]byte, error)
	GetCursedSubjects(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([][]byte, error)
	IsCursedGlobal(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (bool, error)
	IsCursed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subject []byte) (bool, error)
	IsCursedU128(ctx context.Context, opts *bind.CallOpts, ref bind.Object, subjectValue *big.Int) (bool, error)
}

type RmnRemoteEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	GetArm() (*bind.EncodedCall, error)
	GetArmWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(ref bind.Object, ownerCap bind.Object, localChainSelector uint64) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	Verify(ref bind.Object, offRampStateAddress string, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) (*bind.EncodedCall, error)
	VerifyWithArgs(args ...any) (*bind.EncodedCall, error)
	SetConfig(ref bind.Object, param bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*bind.EncodedCall, error)
	SetConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetVersionedConfig(ref bind.Object) (*bind.EncodedCall, error)
	GetVersionedConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetLocalChainSelector(ref bind.Object) (*bind.EncodedCall, error)
	GetLocalChainSelectorWithArgs(args ...any) (*bind.EncodedCall, error)
	GetReportDigestHeader() (*bind.EncodedCall, error)
	GetReportDigestHeaderWithArgs(args ...any) (*bind.EncodedCall, error)
	Curse(ref bind.Object, ownerCap bind.Object, subject []byte) (*bind.EncodedCall, error)
	CurseWithArgs(args ...any) (*bind.EncodedCall, error)
	CurseMultiple(ref bind.Object, param bind.Object, subjects [][]byte) (*bind.EncodedCall, error)
	CurseMultipleWithArgs(args ...any) (*bind.EncodedCall, error)
	Uncurse(ref bind.Object, ownerCap bind.Object, subject []byte) (*bind.EncodedCall, error)
	UncurseWithArgs(args ...any) (*bind.EncodedCall, error)
	UncurseMultiple(ref bind.Object, param bind.Object, subjects [][]byte) (*bind.EncodedCall, error)
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

func NewRmnRemote(packageID string, client sui.ISuiAPI) (*RmnRemoteContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "rmn_remote", client)
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

type McmsCallback struct {
}

type bcsReport struct {
	DestChainSelector           uint64
	RmnRemoteContractAddress    [32]byte
	OffRampAddress              [32]byte
	RmnHomeContractConfigDigest []byte
	MerkleRoots                 []MerkleRoot
}

func convertReportFromBCS(bcs bcsReport) (Report, error) {

	return Report{
		DestChainSelector:           bcs.DestChainSelector,
		RmnRemoteContractAddress:    fmt.Sprintf("0x%x", bcs.RmnRemoteContractAddress),
		OffRampAddress:              fmt.Sprintf("0x%x", bcs.OffRampAddress),
		RmnHomeContractConfigDigest: bcs.RmnHomeContractConfigDigest,
		MerkleRoots:                 bcs.MerkleRoots,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip::rmn_remote::RMNRemoteState", func(data []byte) (interface{}, error) {
		var result RMNRemoteState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::Config", func(data []byte) (interface{}, error) {
		var result Config
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::Signer", func(data []byte) (interface{}, error) {
		var result Signer
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::Report", func(data []byte) (interface{}, error) {
		var temp bcsReport
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertReportFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::MerkleRoot", func(data []byte) (interface{}, error) {
		var result MerkleRoot
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::ConfigSet", func(data []byte) (interface{}, error) {
		var result ConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::Cursed", func(data []byte) (interface{}, error) {
		var result Cursed
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::Uncursed", func(data []byte) (interface{}, error) {
		var result Uncursed
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::rmn_remote::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *RmnRemoteContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetArm executes the get_arm Move function.
func (c *RmnRemoteContract) GetArm(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.GetArm()
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

// Verify executes the verify Move function.
func (c *RmnRemoteContract) Verify(ctx context.Context, opts *bind.CallOpts, ref bind.Object, offRampStateAddress string, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.Verify(ref, offRampStateAddress, merkleRootSourceChainSelectors, merkleRootOnRampAddresses, merkleRootMinSeqNrs, merkleRootMaxSeqNrs, merkleRootValues, signatures)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetConfig executes the set_config Move function.
func (c *RmnRemoteContract) SetConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.SetConfig(ref, param, rmnHomeContractConfigDigest, signerOnchainPublicKeys, nodeIndexes, fSign)
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
func (c *RmnRemoteContract) CurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.CurseMultiple(ref, param, subjects)
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
func (c *RmnRemoteContract) UncurseMultiple(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, subjects [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.UncurseMultiple(ref, param, subjects)
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
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// GetArm executes the get_arm Move function using DevInspect to get return values.
//
// Returns: address
func (d *RmnRemoteDevInspect) GetArm(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.rmnRemoteEncoder.GetArm()
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

// Verify executes the verify Move function using DevInspect to get return values.
//
// Returns: bool
func (d *RmnRemoteDevInspect) Verify(ctx context.Context, opts *bind.CallOpts, ref bind.Object, offRampStateAddress string, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) (bool, error) {
	encoded, err := d.contract.rmnRemoteEncoder.Verify(ref, offRampStateAddress, merkleRootSourceChainSelectors, merkleRootOnRampAddresses, merkleRootMinSeqNrs, merkleRootMaxSeqNrs, merkleRootValues, signatures)
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
	return d.contract.Call(ctx, opts, encoded)
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
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
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
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
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
	result, ok := results[0].([][]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected [][]byte, got %T", results[0])
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
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
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
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
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
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
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

// GetArm encodes a call to the get_arm Move function.
func (c rmnRemoteEncoder) GetArm() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_arm", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"address",
	})
}

// GetArmWithArgs encodes a call to the get_arm Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) GetArmWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_arm", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
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

// Verify encodes a call to the verify Move function.
func (c rmnRemoteEncoder) Verify(ref bind.Object, offRampStateAddress string, merkleRootSourceChainSelectors []uint64, merkleRootOnRampAddresses [][]byte, merkleRootMinSeqNrs []uint64, merkleRootMaxSeqNrs []uint64, merkleRootValues [][]byte, signatures [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<u64>",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
	}, []any{
		ref,
		offRampStateAddress,
		merkleRootSourceChainSelectors,
		merkleRootOnRampAddresses,
		merkleRootMinSeqNrs,
		merkleRootMaxSeqNrs,
		merkleRootValues,
		signatures,
	}, []string{
		"bool",
	})
}

// VerifyWithArgs encodes a call to the verify Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) VerifyWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<u64>",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// SetConfig encodes a call to the set_config Move function.
func (c rmnRemoteEncoder) SetConfig(ref bind.Object, param bind.Object, rmnHomeContractConfigDigest []byte, signerOnchainPublicKeys [][]byte, nodeIndexes []uint64, fSign uint64) (*bind.EncodedCall, error) {
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
		param,
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
func (c rmnRemoteEncoder) CurseMultiple(ref bind.Object, param bind.Object, subjects [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("curse_multiple", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<vector<u8>>",
	}, []any{
		ref,
		param,
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
func (c rmnRemoteEncoder) UncurseMultiple(ref bind.Object, param bind.Object, subjects [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("uncurse_multiple", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<vector<u8>>",
	}, []any{
		ref,
		param,
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
