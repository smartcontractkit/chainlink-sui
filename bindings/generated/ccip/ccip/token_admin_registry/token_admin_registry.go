// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_token_admin_registry

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

type ITokenAdminRegistry interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetPools(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) (*models.SuiTransactionBlockResponse, error)
	GetPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	GetTokenConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	GetTokenConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) (*models.SuiTransactionBlockResponse, error)
	GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, tokenConfig TokenConfig) (*models.SuiTransactionBlockResponse, error)
	GetAllConfiguredTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object, startKey string, maxCount uint64) (*models.SuiTransactionBlockResponse, error)
	RegisterPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, tokenPoolPackageId string, tokenPoolModule string, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	RegisterPoolByAdmin(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, tokenType bind.Object, initialAdministrator string, tokenPoolTypeProof bind.Object, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error)
	UnregisterPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	SetPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, lockOrBurnParams []string, releaseOrMintParams []string, param bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, newAdmin string) (*models.SuiTransactionBlockResponse, error)
	AcceptAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	IsAdministrator(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, administrator string) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ITokenAdminRegistryDevInspect
	Encoder() TokenAdminRegistryEncoder
}

type ITokenAdminRegistryDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetPools(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) ([]string, error)
	GetPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (string, error)
	GetTokenConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (TokenConfig, error)
	GetTokenConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) ([]TokenConfig, error)
	GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, tokenConfig TokenConfig) ([]any, error)
	GetAllConfiguredTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object, startKey string, maxCount uint64) ([]any, error)
	IsAdministrator(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, administrator string) (bool, error)
}

type TokenAdminRegistryEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	GetPools(ref bind.Object, coinMetadataAddresses []string) (*bind.EncodedCall, error)
	GetPoolsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetPool(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	GetPoolWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenConfig(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	GetTokenConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenConfigs(ref bind.Object, coinMetadataAddresses []string) (*bind.EncodedCall, error)
	GetTokenConfigsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenConfigData(tokenConfig TokenConfig) (*bind.EncodedCall, error)
	GetTokenConfigDataWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllConfiguredTokens(ref bind.Object, startKey string, maxCount uint64) (*bind.EncodedCall, error)
	GetAllConfiguredTokensWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterPool(typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, tokenPoolPackageId string, tokenPoolModule string, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, proof bind.Object) (*bind.EncodedCall, error)
	RegisterPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	RegisterPoolByAdmin(ref bind.Object, param bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, tokenType bind.Object, initialAdministrator string, tokenPoolTypeProof bind.Object, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error)
	RegisterPoolByAdminWithArgs(args ...any) (*bind.EncodedCall, error)
	UnregisterPool(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	UnregisterPoolWithArgs(args ...any) (*bind.EncodedCall, error)
	SetPool(typeArgs []string, ref bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, lockOrBurnParams []string, releaseOrMintParams []string, param bind.Object) (*bind.EncodedCall, error)
	SetPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	TransferAdminRole(ref bind.Object, coinMetadataAddress string, newAdmin string) (*bind.EncodedCall, error)
	TransferAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptAdminRole(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	AcceptAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	IsAdministrator(ref bind.Object, coinMetadataAddress string, administrator string) (*bind.EncodedCall, error)
	IsAdministratorWithArgs(args ...any) (*bind.EncodedCall, error)
}

type TokenAdminRegistryContract struct {
	*bind.BoundContract
	tokenAdminRegistryEncoder
	devInspect *TokenAdminRegistryDevInspect
}

type TokenAdminRegistryDevInspect struct {
	contract *TokenAdminRegistryContract
}

var _ ITokenAdminRegistry = (*TokenAdminRegistryContract)(nil)
var _ ITokenAdminRegistryDevInspect = (*TokenAdminRegistryDevInspect)(nil)

func NewTokenAdminRegistry(packageID string, client sui.ISuiAPI) (*TokenAdminRegistryContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "token_admin_registry", client)
	if err != nil {
		return nil, err
	}

	c := &TokenAdminRegistryContract{
		BoundContract:             contract,
		tokenAdminRegistryEncoder: tokenAdminRegistryEncoder{BoundContract: contract},
	}
	c.devInspect = &TokenAdminRegistryDevInspect{contract: c}
	return c, nil
}

func (c *TokenAdminRegistryContract) Encoder() TokenAdminRegistryEncoder {
	return c.tokenAdminRegistryEncoder
}

func (c *TokenAdminRegistryContract) DevInspect() ITokenAdminRegistryDevInspect {
	return c.devInspect
}

type TokenAdminRegistryState struct {
	Id           string      `move:"sui::object::UID"`
	TokenConfigs bind.Object `move:"LinkedTable<address, TokenConfig>"`
}

type TokenConfig struct {
	TokenPoolPackageId   string      `move:"address"`
	TokenPoolModule      string      `move:"0x1::string::String"`
	TokenType            string 	 `move:"ascii::String"`
	Administrator        string      `move:"address"`
	PendingAdministrator string      `move:"address"`
	TokenPoolTypeProof   string `move:"ascii::String"`
	LockOrBurnParams     []string    `move:"vector<address>"`
	ReleaseOrMintParams  []string    `move:"vector<address>"`
}

type PoolSet struct {
	CoinMetadataAddress   string      `move:"address"`
	PreviousPoolPackageId string      `move:"address"`
	NewPoolPackageId      string      `move:"address"`
	TokenPoolTypeProof    bind.Object `move:"ascii::String"`
	LockOrBurnParams      []string    `move:"vector<address>"`
	ReleaseOrMintParams   []string    `move:"vector<address>"`
}

type PoolRegistered struct {
	CoinMetadataAddress string      `move:"address"`
	TokenPoolPackageId  string      `move:"address"`
	Administrator       string      `move:"address"`
	TokenPoolTypeProof  bind.Object `move:"ascii::String"`
}

type PoolUnregistered struct {
	CoinMetadataAddress string `move:"address"`
	PreviousPoolAddress string `move:"address"`
}

type AdministratorTransferRequested struct {
	CoinMetadataAddress string `move:"address"`
	CurrentAdmin        string `move:"address"`
	NewAdmin            string `move:"address"`
}

type AdministratorTransferred struct {
	CoinMetadataAddress string `move:"address"`
	NewAdmin            string `move:"address"`
}

type bcsTokenConfig struct {
	TokenPoolPackageId   [32]byte
	TokenPoolModule      string
	TokenType            string
	Administrator        [32]byte
	PendingAdministrator [32]byte
	TokenPoolTypeProof   string
	LockOrBurnParams     [][32]byte
	ReleaseOrMintParams  [][32]byte
}

func convertTokenConfigFromBCS(bcs bcsTokenConfig) TokenConfig {
	return TokenConfig{
		TokenPoolPackageId:   fmt.Sprintf("0x%x", bcs.TokenPoolPackageId),
		TokenPoolModule:      bcs.TokenPoolModule,
		TokenType:            bcs.TokenType,
		Administrator:        fmt.Sprintf("0x%x", bcs.Administrator),
		PendingAdministrator: fmt.Sprintf("0x%x", bcs.PendingAdministrator),
		TokenPoolTypeProof:   bcs.TokenPoolTypeProof,
		LockOrBurnParams: func() []string {
			addrs := make([]string, len(bcs.LockOrBurnParams))
			for i, addr := range bcs.LockOrBurnParams {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		ReleaseOrMintParams: func() []string {
			addrs := make([]string, len(bcs.ReleaseOrMintParams))
			for i, addr := range bcs.ReleaseOrMintParams {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}
}

type bcsPoolSet struct {
	CoinMetadataAddress   [32]byte
	PreviousPoolPackageId [32]byte
	NewPoolPackageId      [32]byte
	TokenPoolTypeProof    bind.Object
	LockOrBurnParams      [][32]byte
	ReleaseOrMintParams   [][32]byte
}

func convertPoolSetFromBCS(bcs bcsPoolSet) PoolSet {
	return PoolSet{
		CoinMetadataAddress:   fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		PreviousPoolPackageId: fmt.Sprintf("0x%x", bcs.PreviousPoolPackageId),
		NewPoolPackageId:      fmt.Sprintf("0x%x", bcs.NewPoolPackageId),
		TokenPoolTypeProof:    bcs.TokenPoolTypeProof,
		LockOrBurnParams: func() []string {
			addrs := make([]string, len(bcs.LockOrBurnParams))
			for i, addr := range bcs.LockOrBurnParams {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		ReleaseOrMintParams: func() []string {
			addrs := make([]string, len(bcs.ReleaseOrMintParams))
			for i, addr := range bcs.ReleaseOrMintParams {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}
}

type bcsPoolRegistered struct {
	CoinMetadataAddress [32]byte
	TokenPoolPackageId  [32]byte
	Administrator       [32]byte
	TokenPoolTypeProof  bind.Object
}

func convertPoolRegisteredFromBCS(bcs bcsPoolRegistered) PoolRegistered {
	return PoolRegistered{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		TokenPoolPackageId:  fmt.Sprintf("0x%x", bcs.TokenPoolPackageId),
		Administrator:       fmt.Sprintf("0x%x", bcs.Administrator),
		TokenPoolTypeProof:  bcs.TokenPoolTypeProof,
	}
}

type bcsPoolUnregistered struct {
	CoinMetadataAddress [32]byte
	PreviousPoolAddress [32]byte
}

func convertPoolUnregisteredFromBCS(bcs bcsPoolUnregistered) PoolUnregistered {
	return PoolUnregistered{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		PreviousPoolAddress: fmt.Sprintf("0x%x", bcs.PreviousPoolAddress),
	}
}

type bcsAdministratorTransferRequested struct {
	CoinMetadataAddress [32]byte
	CurrentAdmin        [32]byte
	NewAdmin            [32]byte
}

func convertAdministratorTransferRequestedFromBCS(bcs bcsAdministratorTransferRequested) AdministratorTransferRequested {
	return AdministratorTransferRequested{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		CurrentAdmin:        fmt.Sprintf("0x%x", bcs.CurrentAdmin),
		NewAdmin:            fmt.Sprintf("0x%x", bcs.NewAdmin),
	}
}

type bcsAdministratorTransferred struct {
	CoinMetadataAddress [32]byte
	NewAdmin            [32]byte
}

func convertAdministratorTransferredFromBCS(bcs bcsAdministratorTransferred) AdministratorTransferred {
	return AdministratorTransferred{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		NewAdmin:            fmt.Sprintf("0x%x", bcs.NewAdmin),
	}
}

func init() {
	bind.RegisterStructDecoder("ccip::token_admin_registry::TokenAdminRegistryState", func(data []byte) (interface{}, error) {
		var result TokenAdminRegistryState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::TokenConfig", func(data []byte) (interface{}, error) {
		var temp bcsTokenConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertTokenConfigFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::PoolSet", func(data []byte) (interface{}, error) {
		var temp bcsPoolSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertPoolSetFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::PoolRegistered", func(data []byte) (interface{}, error) {
		var temp bcsPoolRegistered
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertPoolRegisteredFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::PoolUnregistered", func(data []byte) (interface{}, error) {
		var temp bcsPoolUnregistered
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertPoolUnregisteredFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::AdministratorTransferRequested", func(data []byte) (interface{}, error) {
		var temp bcsAdministratorTransferRequested
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertAdministratorTransferRequestedFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::AdministratorTransferred", func(data []byte) (interface{}, error) {
		var temp bcsAdministratorTransferred
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertAdministratorTransferredFromBCS(temp)
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *TokenAdminRegistryContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *TokenAdminRegistryContract) Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.Initialize(ref, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetPools executes the get_pools Move function.
func (c *TokenAdminRegistryContract) GetPools(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetPools(ref, coinMetadataAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetPool executes the get_pool Move function.
func (c *TokenAdminRegistryContract) GetPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetPool(ref, coinMetadataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenConfig executes the get_token_config Move function.
func (c *TokenAdminRegistryContract) GetTokenConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetTokenConfig(ref, coinMetadataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenConfigs executes the get_token_configs Move function.
func (c *TokenAdminRegistryContract) GetTokenConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetTokenConfigs(ref, coinMetadataAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenConfigData executes the get_token_config_data Move function.
func (c *TokenAdminRegistryContract) GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, tokenConfig TokenConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetTokenConfigData(tokenConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllConfiguredTokens executes the get_all_configured_tokens Move function.
func (c *TokenAdminRegistryContract) GetAllConfiguredTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object, startKey string, maxCount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetAllConfiguredTokens(ref, startKey, maxCount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterPool executes the register_pool Move function.
func (c *TokenAdminRegistryContract) RegisterPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, tokenPoolPackageId string, tokenPoolModule string, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.RegisterPool(typeArgs, ref, param, coinMetadata, tokenPoolPackageId, tokenPoolModule, initialAdministrator, lockOrBurnParams, releaseOrMintParams, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterPoolByAdmin executes the register_pool_by_admin Move function.
func (c *TokenAdminRegistryContract) RegisterPoolByAdmin(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, tokenType bind.Object, initialAdministrator string, tokenPoolTypeProof bind.Object, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.RegisterPoolByAdmin(ref, param, coinMetadataAddress, tokenPoolPackageId, tokenPoolModule, tokenType, initialAdministrator, tokenPoolTypeProof, lockOrBurnParams, releaseOrMintParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UnregisterPool executes the unregister_pool Move function.
func (c *TokenAdminRegistryContract) UnregisterPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.UnregisterPool(ref, coinMetadataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetPool executes the set_pool Move function.
func (c *TokenAdminRegistryContract) SetPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, lockOrBurnParams []string, releaseOrMintParams []string, param bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.SetPool(typeArgs, ref, coinMetadataAddress, tokenPoolPackageId, tokenPoolModule, lockOrBurnParams, releaseOrMintParams, param)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferAdminRole executes the transfer_admin_role Move function.
func (c *TokenAdminRegistryContract) TransferAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, newAdmin string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.TransferAdminRole(ref, coinMetadataAddress, newAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptAdminRole executes the accept_admin_role Move function.
func (c *TokenAdminRegistryContract) AcceptAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.AcceptAdminRole(ref, coinMetadataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsAdministrator executes the is_administrator Move function.
func (c *TokenAdminRegistryContract) IsAdministrator(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, administrator string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.IsAdministrator(ref, coinMetadataAddress, administrator)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *TokenAdminRegistryDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.TypeAndVersion()
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

// GetPools executes the get_pools Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *TokenAdminRegistryDevInspect) GetPools(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) ([]string, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetPools(ref, coinMetadataAddresses)
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

// GetPool executes the get_pool Move function using DevInspect to get return values.
//
// Returns: address
func (d *TokenAdminRegistryDevInspect) GetPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (string, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetPool(ref, coinMetadataAddress)
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

// GetTokenConfig executes the get_token_config Move function using DevInspect to get return values.
//
// Returns: TokenConfig
func (d *TokenAdminRegistryDevInspect) GetTokenConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (TokenConfig, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetTokenConfig(ref, coinMetadataAddress)
	if err != nil {
		return TokenConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TokenConfig{}, err
	}
	if len(results) == 0 {
		return TokenConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TokenConfig)
	if !ok {
		return TokenConfig{}, fmt.Errorf("unexpected return type: expected TokenConfig, got %T", results[0])
	}
	return result, nil
}

// GetTokenConfigs executes the get_token_configs Move function using DevInspect to get return values.
//
// Returns: vector<TokenConfig>
func (d *TokenAdminRegistryDevInspect) GetTokenConfigs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) ([]TokenConfig, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetTokenConfigs(ref, coinMetadataAddresses)
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
	result, ok := results[0].([]TokenConfig)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []TokenConfig, got %T", results[0])
	}
	return result, nil
}

// GetTokenConfigData executes the get_token_config_data Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: 0x1::string::String
//	[2]: ascii::String
//	[3]: address
//	[4]: address
//	[5]: ascii::String
//	[6]: vector<address>
//	[7]: vector<address>
func (d *TokenAdminRegistryDevInspect) GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, tokenConfig TokenConfig) ([]any, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetTokenConfigData(tokenConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetAllConfiguredTokens executes the get_all_configured_tokens Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<address>
//	[1]: address
//	[2]: bool
func (d *TokenAdminRegistryDevInspect) GetAllConfiguredTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object, startKey string, maxCount uint64) ([]any, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetAllConfiguredTokens(ref, startKey, maxCount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// IsAdministrator executes the is_administrator Move function using DevInspect to get return values.
//
// Returns: bool
func (d *TokenAdminRegistryDevInspect) IsAdministrator(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, administrator string) (bool, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.IsAdministrator(ref, coinMetadataAddress, administrator)
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

type tokenAdminRegistryEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c tokenAdminRegistryEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c tokenAdminRegistryEncoder) Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}, []any{
		ref,
		ownerCap,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetPools encodes a call to the get_pools Move function.
func (c tokenAdminRegistryEncoder) GetPools(ref bind.Object, coinMetadataAddresses []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_pools", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"vector<address>",
	}, []any{
		ref,
		coinMetadataAddresses,
	}, []string{
		"vector<address>",
	})
}

// GetPoolsWithArgs encodes a call to the get_pools Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetPoolsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_pools", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// GetPool encodes a call to the get_pool Move function.
func (c tokenAdminRegistryEncoder) GetPool(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_pool", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
	}, []string{
		"address",
	})
}

// GetPoolWithArgs encodes a call to the get_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetPoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_pool", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetTokenConfig encodes a call to the get_token_config Move function.
func (c tokenAdminRegistryEncoder) GetTokenConfig(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
	}, []string{
		"ccip::token_admin_registry::TokenConfig",
	})
}

// GetTokenConfigWithArgs encodes a call to the get_token_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetTokenConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::token_admin_registry::TokenConfig",
	})
}

// GetTokenConfigs encodes a call to the get_token_configs Move function.
func (c tokenAdminRegistryEncoder) GetTokenConfigs(ref bind.Object, coinMetadataAddresses []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_configs", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"vector<address>",
	}, []any{
		ref,
		coinMetadataAddresses,
	}, []string{
		"vector<ccip::token_admin_registry::TokenConfig>",
	})
}

// GetTokenConfigsWithArgs encodes a call to the get_token_configs Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetTokenConfigsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_configs", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ccip::token_admin_registry::TokenConfig>",
	})
}

// GetTokenConfigData encodes a call to the get_token_config_data Move function.
func (c tokenAdminRegistryEncoder) GetTokenConfigData(tokenConfig TokenConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_config_data", typeArgsList, typeParamsList, []string{
		"ccip::token_admin_registry::TokenConfig",
	}, []any{
		tokenConfig,
	}, []string{
		"address",
		"0x1::string::String",
		"ascii::String",
		"address",
		"address",
		"ascii::String",
		"vector<address>",
		"vector<address>",
	})
}

// GetTokenConfigDataWithArgs encodes a call to the get_token_config_data Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetTokenConfigDataWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip::token_admin_registry::TokenConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_config_data", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"0x1::string::String",
		"ascii::String",
		"address",
		"address",
		"ascii::String",
		"vector<address>",
		"vector<address>",
	})
}

// GetAllConfiguredTokens encodes a call to the get_all_configured_tokens Move function.
func (c tokenAdminRegistryEncoder) GetAllConfiguredTokens(ref bind.Object, startKey string, maxCount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_configured_tokens", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"u64",
	}, []any{
		ref,
		startKey,
		maxCount,
	}, []string{
		"vector<address>",
		"address",
		"bool",
	})
}

// GetAllConfiguredTokensWithArgs encodes a call to the get_all_configured_tokens Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetAllConfiguredTokensWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_configured_tokens", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
		"address",
		"bool",
	})
}

// RegisterPool encodes a call to the register_pool Move function.
func (c tokenAdminRegistryEncoder) RegisterPool(typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, tokenPoolPackageId string, tokenPoolModule string, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, proof bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("register_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&TreasuryCap<T>",
		"&CoinMetadata<T>",
		"address",
		"0x1::string::String",
		"address",
		"vector<address>",
		"vector<address>",
		"TypeProof",
	}, []any{
		ref,
		param,
		coinMetadata,
		tokenPoolPackageId,
		tokenPoolModule,
		initialAdministrator,
		lockOrBurnParams,
		releaseOrMintParams,
		proof,
	}, nil)
}

// RegisterPoolWithArgs encodes a call to the register_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) RegisterPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&TreasuryCap<T>",
		"&CoinMetadata<T>",
		"address",
		"0x1::string::String",
		"address",
		"vector<address>",
		"vector<address>",
		"TypeProof",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("register_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RegisterPoolByAdmin encodes a call to the register_pool_by_admin Move function.
func (c tokenAdminRegistryEncoder) RegisterPoolByAdmin(ref bind.Object, param bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, tokenType bind.Object, initialAdministrator string, tokenPoolTypeProof bind.Object, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_pool_by_admin", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"state_object::CCIPAdminProof",
		"address",
		"address",
		"0x1::string::String",
		"ascii::String",
		"address",
		"ascii::String",
		"vector<address>",
		"vector<address>",
	}, []any{
		ref,
		param,
		coinMetadataAddress,
		tokenPoolPackageId,
		tokenPoolModule,
		tokenType,
		initialAdministrator,
		tokenPoolTypeProof,
		lockOrBurnParams,
		releaseOrMintParams,
	}, nil)
}

// RegisterPoolByAdminWithArgs encodes a call to the register_pool_by_admin Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) RegisterPoolByAdminWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"state_object::CCIPAdminProof",
		"address",
		"address",
		"0x1::string::String",
		"ascii::String",
		"address",
		"ascii::String",
		"vector<address>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_pool_by_admin", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// UnregisterPool encodes a call to the unregister_pool Move function.
func (c tokenAdminRegistryEncoder) UnregisterPool(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unregister_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
	}, nil)
}

// UnregisterPoolWithArgs encodes a call to the unregister_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) UnregisterPoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unregister_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetPool encodes a call to the set_pool Move function.
func (c tokenAdminRegistryEncoder) SetPool(typeArgs []string, ref bind.Object, coinMetadataAddress string, tokenPoolPackageId string, tokenPoolModule string, lockOrBurnParams []string, releaseOrMintParams []string, param bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("set_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"address",
		"address",
		"0x1::string::String",
		"vector<address>",
		"vector<address>",
		"TypeProof",
	}, []any{
		ref,
		coinMetadataAddress,
		tokenPoolPackageId,
		tokenPoolModule,
		lockOrBurnParams,
		releaseOrMintParams,
		param,
	}, nil)
}

// SetPoolWithArgs encodes a call to the set_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) SetPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"address",
		"address",
		"0x1::string::String",
		"vector<address>",
		"vector<address>",
		"TypeProof",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("set_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// TransferAdminRole encodes a call to the transfer_admin_role Move function.
func (c tokenAdminRegistryEncoder) TransferAdminRole(ref bind.Object, coinMetadataAddress string, newAdmin string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_admin_role", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"address",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
		newAdmin,
	}, nil)
}

// TransferAdminRoleWithArgs encodes a call to the transfer_admin_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) TransferAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_admin_role", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptAdminRole encodes a call to the accept_admin_role Move function.
func (c tokenAdminRegistryEncoder) AcceptAdminRole(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_admin_role", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
	}, nil)
}

// AcceptAdminRoleWithArgs encodes a call to the accept_admin_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) AcceptAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_admin_role", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IsAdministrator encodes a call to the is_administrator Move function.
func (c tokenAdminRegistryEncoder) IsAdministrator(ref bind.Object, coinMetadataAddress string, administrator string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_administrator", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
		administrator,
	}, []string{
		"bool",
	})
}

// IsAdministratorWithArgs encodes a call to the is_administrator Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) IsAdministratorWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_administrator", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}
