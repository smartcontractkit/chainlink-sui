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

const FunctionInfo = `[{"package":"ccip","module":"token_admin_registry","name":"accept_admin_role","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"get_all_configured_tokens","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"start_key","type":"address"},{"name":"max_count","type":"u64"}]},{"package":"ccip","module":"token_admin_registry","name":"get_pool","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"get_pool_local_token","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"token_pool_package_id","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"get_pools","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_addresses","type":"vector<address>"}]},{"package":"ccip","module":"token_admin_registry","name":"get_token_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"get_token_config_data","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"get_token_config_struct","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"initialize","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"}]},{"package":"ccip","module":"token_admin_registry","name":"is_administrator","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"},{"name":"administrator","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"is_pool_registered","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"register_pool","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"_","type":"TreasuryCap<T>"},{"name":"coin_metadata","type":"CoinMetadata<T>"},{"name":"initial_administrator","type":"address"},{"name":"lock_or_burn_params","type":"vector<address>"},{"name":"release_or_mint_params","type":"vector<address>"},{"name":"publisher_wrapper","type":"PublisherWrapper<TypeProof>"},{"name":"_proof","type":"TypeProof"}]},{"package":"ccip","module":"token_admin_registry","name":"register_pool_as_owner","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"},{"name":"package_address","type":"address"},{"name":"token_pool_module","type":"0x1::string::String"},{"name":"token_type","type":"ascii::String"},{"name":"initial_administrator","type":"address"},{"name":"token_pool_type_proof","type":"ascii::String"},{"name":"lock_or_burn_params","type":"vector<address>"},{"name":"release_or_mint_params","type":"vector<address>"}]},{"package":"ccip","module":"token_admin_registry","name":"transfer_admin_role","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"},{"name":"new_admin","type":"address"}]},{"package":"ccip","module":"token_admin_registry","name":"type_and_version","parameters":null},{"package":"ccip","module":"token_admin_registry","name":"unregister_pool","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"coin_metadata_address","type":"address"}]}]`

type ITokenAdminRegistry interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetPools(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) (*models.SuiTransactionBlockResponse, error)
	GetPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	GetTokenConfigStruct(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	GetPoolLocalToken(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokenPoolPackageId string) (*models.SuiTransactionBlockResponse, error)
	GetTokenConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	GetAllConfiguredTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object, startKey string, maxCount uint64) (*models.SuiTransactionBlockResponse, error)
	RegisterPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, publisherWrapper bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	RegisterPoolAsOwner(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ref bind.Object, coinMetadataAddress string, packageAddress string, tokenPoolModule string, tokenType string, initialAdministrator string, tokenPoolTypeProof string, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error)
	UnregisterPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	TransferAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, newAdmin string) (*models.SuiTransactionBlockResponse, error)
	AcceptAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	IsPoolRegistered(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error)
	IsAdministrator(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string, administrator string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsUnregisterPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsTransferAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsAcceptAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ITokenAdminRegistryDevInspect
	Encoder() TokenAdminRegistryEncoder
	Bound() bind.IBoundContract
}

type ITokenAdminRegistryDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetPools(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddresses []string) ([]string, error)
	GetPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (string, error)
	GetTokenConfigStruct(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (TokenConfig, error)
	GetPoolLocalToken(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokenPoolPackageId string) (string, error)
	GetTokenConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) ([]any, error)
	GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) ([]any, error)
	GetAllConfiguredTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object, startKey string, maxCount uint64) ([]any, error)
	IsPoolRegistered(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (bool, error)
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
	GetTokenConfigStruct(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	GetTokenConfigStructWithArgs(args ...any) (*bind.EncodedCall, error)
	GetPoolLocalToken(ref bind.Object, tokenPoolPackageId string) (*bind.EncodedCall, error)
	GetPoolLocalTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenConfig(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	GetTokenConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenConfigData(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	GetTokenConfigDataWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllConfiguredTokens(ref bind.Object, startKey string, maxCount uint64) (*bind.EncodedCall, error)
	GetAllConfiguredTokensWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterPool(typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, publisherWrapper bind.Object, proof bind.Object) (*bind.EncodedCall, error)
	RegisterPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	RegisterPoolAsOwner(ownerCap bind.Object, ref bind.Object, coinMetadataAddress string, packageAddress string, tokenPoolModule string, tokenType string, initialAdministrator string, tokenPoolTypeProof string, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error)
	RegisterPoolAsOwnerWithArgs(args ...any) (*bind.EncodedCall, error)
	UnregisterPool(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	UnregisterPoolWithArgs(args ...any) (*bind.EncodedCall, error)
	TransferAdminRole(ref bind.Object, coinMetadataAddress string, newAdmin string) (*bind.EncodedCall, error)
	TransferAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptAdminRole(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	AcceptAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	IsPoolRegistered(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error)
	IsPoolRegisteredWithArgs(args ...any) (*bind.EncodedCall, error)
	IsAdministrator(ref bind.Object, coinMetadataAddress string, administrator string) (*bind.EncodedCall, error)
	IsAdministratorWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterPool(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRegisterPoolWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsUnregisterPool(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsUnregisterPoolWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsTransferAdminRole(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsTransferAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsAcceptAdminRole(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAcceptAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewTokenAdminRegistry(packageID string, client sui.ISuiAPI) (ITokenAdminRegistry, error) {
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

func (c *TokenAdminRegistryContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *TokenAdminRegistryContract) Encoder() TokenAdminRegistryEncoder {
	return c.tokenAdminRegistryEncoder
}

func (c *TokenAdminRegistryContract) DevInspect() ITokenAdminRegistryDevInspect {
	return c.devInspect
}

type TokenAdminRegistryState struct {
	Id                               string      `move:"sui::object::UID"`
	TokenConfigs                     bind.Object `move:"LinkedTable<address, TokenConfig>"`
	TokenPoolPackageIdToCoinMetadata bind.Object `move:"LinkedTable<address, address>"`
}

type TokenConfig struct {
	TokenPoolPackageId   string   `move:"address"`
	TokenPoolModule      string   `move:"0x1::string::String"`
	TokenType            string   `move:"ascii::String"`
	Administrator        string   `move:"address"`
	PendingAdministrator string   `move:"address"`
	TokenPoolTypeProof   string   `move:"ascii::String"`
	LockOrBurnParams     []string `move:"vector<address>"`
	ReleaseOrMintParams  []string `move:"vector<address>"`
}

type PoolSet struct {
	CoinMetadataAddress   string   `move:"address"`
	PreviousPoolPackageId string   `move:"address"`
	NewPoolPackageId      string   `move:"address"`
	TokenPoolTypeProof    string   `move:"ascii::String"`
	LockOrBurnParams      []string `move:"vector<address>"`
	ReleaseOrMintParams   []string `move:"vector<address>"`
}

type PoolRegistered struct {
	CoinMetadataAddress string `move:"address"`
	TokenPoolPackageId  string `move:"address"`
	Administrator       string `move:"address"`
	TokenPoolTypeProof  string `move:"ascii::String"`
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

func convertTokenConfigFromBCS(bcs bcsTokenConfig) (TokenConfig, error) {

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
	}, nil
}

type bcsPoolSet struct {
	CoinMetadataAddress   [32]byte
	PreviousPoolPackageId [32]byte
	NewPoolPackageId      [32]byte
	TokenPoolTypeProof    string
	LockOrBurnParams      [][32]byte
	ReleaseOrMintParams   [][32]byte
}

func convertPoolSetFromBCS(bcs bcsPoolSet) (PoolSet, error) {

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
	}, nil
}

type bcsPoolRegistered struct {
	CoinMetadataAddress [32]byte
	TokenPoolPackageId  [32]byte
	Administrator       [32]byte
	TokenPoolTypeProof  string
}

func convertPoolRegisteredFromBCS(bcs bcsPoolRegistered) (PoolRegistered, error) {

	return PoolRegistered{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		TokenPoolPackageId:  fmt.Sprintf("0x%x", bcs.TokenPoolPackageId),
		Administrator:       fmt.Sprintf("0x%x", bcs.Administrator),
		TokenPoolTypeProof:  bcs.TokenPoolTypeProof,
	}, nil
}

type bcsPoolUnregistered struct {
	CoinMetadataAddress [32]byte
	PreviousPoolAddress [32]byte
}

func convertPoolUnregisteredFromBCS(bcs bcsPoolUnregistered) (PoolUnregistered, error) {

	return PoolUnregistered{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		PreviousPoolAddress: fmt.Sprintf("0x%x", bcs.PreviousPoolAddress),
	}, nil
}

type bcsAdministratorTransferRequested struct {
	CoinMetadataAddress [32]byte
	CurrentAdmin        [32]byte
	NewAdmin            [32]byte
}

func convertAdministratorTransferRequestedFromBCS(bcs bcsAdministratorTransferRequested) (AdministratorTransferRequested, error) {

	return AdministratorTransferRequested{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		CurrentAdmin:        fmt.Sprintf("0x%x", bcs.CurrentAdmin),
		NewAdmin:            fmt.Sprintf("0x%x", bcs.NewAdmin),
	}, nil
}

type bcsAdministratorTransferred struct {
	CoinMetadataAddress [32]byte
	NewAdmin            [32]byte
}

func convertAdministratorTransferredFromBCS(bcs bcsAdministratorTransferred) (AdministratorTransferred, error) {

	return AdministratorTransferred{
		CoinMetadataAddress: fmt.Sprintf("0x%x", bcs.CoinMetadataAddress),
		NewAdmin:            fmt.Sprintf("0x%x", bcs.NewAdmin),
	}, nil
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
	// Register vector decoder for TokenAdminRegistryState
	bind.RegisterStructDecoder("vector<ccip::token_admin_registry::TokenAdminRegistryState>", func(data []byte) (interface{}, error) {
		var results []TokenAdminRegistryState
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::TokenConfig", func(data []byte) (interface{}, error) {
		var temp bcsTokenConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TokenConfig
	bind.RegisterStructDecoder("vector<ccip::token_admin_registry::TokenConfig>", func(data []byte) (interface{}, error) {
		var temps []bcsTokenConfig
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]TokenConfig, len(temps))
		for i, temp := range temps {
			result, err := convertTokenConfigFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::PoolSet", func(data []byte) (interface{}, error) {
		var temp bcsPoolSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertPoolSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for PoolSet
	bind.RegisterStructDecoder("vector<ccip::token_admin_registry::PoolSet>", func(data []byte) (interface{}, error) {
		var temps []bcsPoolSet
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]PoolSet, len(temps))
		for i, temp := range temps {
			result, err := convertPoolSetFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::PoolRegistered", func(data []byte) (interface{}, error) {
		var temp bcsPoolRegistered
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertPoolRegisteredFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for PoolRegistered
	bind.RegisterStructDecoder("vector<ccip::token_admin_registry::PoolRegistered>", func(data []byte) (interface{}, error) {
		var temps []bcsPoolRegistered
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]PoolRegistered, len(temps))
		for i, temp := range temps {
			result, err := convertPoolRegisteredFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::PoolUnregistered", func(data []byte) (interface{}, error) {
		var temp bcsPoolUnregistered
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertPoolUnregisteredFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for PoolUnregistered
	bind.RegisterStructDecoder("vector<ccip::token_admin_registry::PoolUnregistered>", func(data []byte) (interface{}, error) {
		var temps []bcsPoolUnregistered
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]PoolUnregistered, len(temps))
		for i, temp := range temps {
			result, err := convertPoolUnregisteredFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::AdministratorTransferRequested", func(data []byte) (interface{}, error) {
		var temp bcsAdministratorTransferRequested
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAdministratorTransferRequestedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for AdministratorTransferRequested
	bind.RegisterStructDecoder("vector<ccip::token_admin_registry::AdministratorTransferRequested>", func(data []byte) (interface{}, error) {
		var temps []bcsAdministratorTransferRequested
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]AdministratorTransferRequested, len(temps))
		for i, temp := range temps {
			result, err := convertAdministratorTransferRequestedFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::token_admin_registry::AdministratorTransferred", func(data []byte) (interface{}, error) {
		var temp bcsAdministratorTransferred
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAdministratorTransferredFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for AdministratorTransferred
	bind.RegisterStructDecoder("vector<ccip::token_admin_registry::AdministratorTransferred>", func(data []byte) (interface{}, error) {
		var temps []bcsAdministratorTransferred
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]AdministratorTransferred, len(temps))
		for i, temp := range temps {
			result, err := convertAdministratorTransferredFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
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

// GetTokenConfigStruct executes the get_token_config_struct Move function.
func (c *TokenAdminRegistryContract) GetTokenConfigStruct(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetTokenConfigStruct(ref, coinMetadataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetPoolLocalToken executes the get_pool_local_token Move function.
func (c *TokenAdminRegistryContract) GetPoolLocalToken(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokenPoolPackageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetPoolLocalToken(ref, tokenPoolPackageId)
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

// GetTokenConfigData executes the get_token_config_data Move function.
func (c *TokenAdminRegistryContract) GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.GetTokenConfigData(ref, coinMetadataAddress)
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
func (c *TokenAdminRegistryContract) RegisterPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, publisherWrapper bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.RegisterPool(typeArgs, ref, param, coinMetadata, initialAdministrator, lockOrBurnParams, releaseOrMintParams, publisherWrapper, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterPoolAsOwner executes the register_pool_as_owner Move function.
func (c *TokenAdminRegistryContract) RegisterPoolAsOwner(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ref bind.Object, coinMetadataAddress string, packageAddress string, tokenPoolModule string, tokenType string, initialAdministrator string, tokenPoolTypeProof string, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.RegisterPoolAsOwner(ownerCap, ref, coinMetadataAddress, packageAddress, tokenPoolModule, tokenType, initialAdministrator, tokenPoolTypeProof, lockOrBurnParams, releaseOrMintParams)
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

// IsPoolRegistered executes the is_pool_registered Move function.
func (c *TokenAdminRegistryContract) IsPoolRegistered(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.IsPoolRegistered(ref, coinMetadataAddress)
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

// McmsRegisterPool executes the mcms_register_pool Move function.
func (c *TokenAdminRegistryContract) McmsRegisterPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.McmsRegisterPool(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsUnregisterPool executes the mcms_unregister_pool Move function.
func (c *TokenAdminRegistryContract) McmsUnregisterPool(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.McmsUnregisterPool(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsTransferAdminRole executes the mcms_transfer_admin_role Move function.
func (c *TokenAdminRegistryContract) McmsTransferAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.McmsTransferAdminRole(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAcceptAdminRole executes the mcms_accept_admin_role Move function.
func (c *TokenAdminRegistryContract) McmsAcceptAdminRole(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenAdminRegistryEncoder.McmsAcceptAdminRole(ref, registry, params)
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

// GetTokenConfigStruct executes the get_token_config_struct Move function using DevInspect to get return values.
//
// Returns: TokenConfig
func (d *TokenAdminRegistryDevInspect) GetTokenConfigStruct(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (TokenConfig, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetTokenConfigStruct(ref, coinMetadataAddress)
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

// GetPoolLocalToken executes the get_pool_local_token Move function using DevInspect to get return values.
//
// Returns: address
func (d *TokenAdminRegistryDevInspect) GetPoolLocalToken(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokenPoolPackageId string) (string, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetPoolLocalToken(ref, tokenPoolPackageId)
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
// Returns:
//
//	[0]: address
//	[1]: address
//	[2]: address
func (d *TokenAdminRegistryDevInspect) GetTokenConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) ([]any, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetTokenConfig(ref, coinMetadataAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
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
func (d *TokenAdminRegistryDevInspect) GetTokenConfigData(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) ([]any, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.GetTokenConfigData(ref, coinMetadataAddress)
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

// IsPoolRegistered executes the is_pool_registered Move function using DevInspect to get return values.
//
// Returns: bool
func (d *TokenAdminRegistryDevInspect) IsPoolRegistered(ctx context.Context, opts *bind.CallOpts, ref bind.Object, coinMetadataAddress string) (bool, error) {
	encoded, err := d.contract.tokenAdminRegistryEncoder.IsPoolRegistered(ref, coinMetadataAddress)
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

// GetTokenConfigStruct encodes a call to the get_token_config_struct Move function.
func (c tokenAdminRegistryEncoder) GetTokenConfigStruct(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_config_struct", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
	}, []string{
		"ccip::token_admin_registry::TokenConfig",
	})
}

// GetTokenConfigStructWithArgs encodes a call to the get_token_config_struct Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetTokenConfigStructWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_config_struct", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::token_admin_registry::TokenConfig",
	})
}

// GetPoolLocalToken encodes a call to the get_pool_local_token Move function.
func (c tokenAdminRegistryEncoder) GetPoolLocalToken(ref bind.Object, tokenPoolPackageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_pool_local_token", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		tokenPoolPackageId,
	}, []string{
		"address",
	})
}

// GetPoolLocalTokenWithArgs encodes a call to the get_pool_local_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) GetPoolLocalTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_pool_local_token", typeArgsList, typeParamsList, expectedParams, args, []string{
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
		"address",
		"address",
		"address",
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
		"address",
		"address",
		"address",
	})
}

// GetTokenConfigData encodes a call to the get_token_config_data Move function.
func (c tokenAdminRegistryEncoder) GetTokenConfigData(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_config_data", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
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
		"&CCIPObjectRef",
		"address",
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
func (c tokenAdminRegistryEncoder) RegisterPool(typeArgs []string, ref bind.Object, param bind.Object, coinMetadata bind.Object, initialAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string, publisherWrapper bind.Object, proof bind.Object) (*bind.EncodedCall, error) {
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
		"vector<address>",
		"vector<address>",
		"PublisherWrapper<TypeProof>",
		"TypeProof",
	}, []any{
		ref,
		param,
		coinMetadata,
		initialAdministrator,
		lockOrBurnParams,
		releaseOrMintParams,
		publisherWrapper,
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
		"vector<address>",
		"vector<address>",
		"PublisherWrapper<TypeProof>",
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

// RegisterPoolAsOwner encodes a call to the register_pool_as_owner Move function.
func (c tokenAdminRegistryEncoder) RegisterPoolAsOwner(ownerCap bind.Object, ref bind.Object, coinMetadataAddress string, packageAddress string, tokenPoolModule string, tokenType string, initialAdministrator string, tokenPoolTypeProof string, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_pool_as_owner", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut CCIPObjectRef",
		"address",
		"address",
		"0x1::string::String",
		"ascii::String",
		"address",
		"ascii::String",
		"vector<address>",
		"vector<address>",
	}, []any{
		ownerCap,
		ref,
		coinMetadataAddress,
		packageAddress,
		tokenPoolModule,
		tokenType,
		initialAdministrator,
		tokenPoolTypeProof,
		lockOrBurnParams,
		releaseOrMintParams,
	}, nil)
}

// RegisterPoolAsOwnerWithArgs encodes a call to the register_pool_as_owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) RegisterPoolAsOwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut CCIPObjectRef",
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
	return c.EncodeCallArgsWithGenerics("register_pool_as_owner", typeArgsList, typeParamsList, expectedParams, args, nil)
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

// IsPoolRegistered encodes a call to the is_pool_registered Move function.
func (c tokenAdminRegistryEncoder) IsPoolRegistered(ref bind.Object, coinMetadataAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_pool_registered", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		coinMetadataAddress,
	}, []string{
		"bool",
	})
}

// IsPoolRegisteredWithArgs encodes a call to the is_pool_registered Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) IsPoolRegisteredWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_pool_registered", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
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

// McmsRegisterPool encodes a call to the mcms_register_pool Move function.
func (c tokenAdminRegistryEncoder) McmsRegisterPool(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsRegisterPoolWithArgs encodes a call to the mcms_register_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) McmsRegisterPoolWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_register_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsUnregisterPool encodes a call to the mcms_unregister_pool Move function.
func (c tokenAdminRegistryEncoder) McmsUnregisterPool(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_unregister_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsUnregisterPoolWithArgs encodes a call to the mcms_unregister_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) McmsUnregisterPoolWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_unregister_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsTransferAdminRole encodes a call to the mcms_transfer_admin_role Move function.
func (c tokenAdminRegistryEncoder) McmsTransferAdminRole(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_admin_role", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsTransferAdminRoleWithArgs encodes a call to the mcms_transfer_admin_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) McmsTransferAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_transfer_admin_role", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAcceptAdminRole encodes a call to the mcms_accept_admin_role Move function.
func (c tokenAdminRegistryEncoder) McmsAcceptAdminRole(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_admin_role", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsAcceptAdminRoleWithArgs encodes a call to the mcms_accept_admin_role Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenAdminRegistryEncoder) McmsAcceptAdminRoleWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_accept_admin_role", typeArgsList, typeParamsList, expectedParams, args, nil)
}
