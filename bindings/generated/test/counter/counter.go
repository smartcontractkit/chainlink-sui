// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_counter

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

type ICounter interface {
	Initialize(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Increment(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	Decrement(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	Create(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	IncrementByOne(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	IncrementByOneNoContext(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	IncrementByTwo(ctx context.Context, opts *bind.CallOpts, admin bind.Object, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	IncrementByTwoNoContext(ctx context.Context, opts *bind.CallOpts, admin bind.Object, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	IncrementBy(ctx context.Context, opts *bind.CallOpts, counter bind.Object, by uint64) (*models.SuiTransactionBlockResponse, error)
	IncrementMult(ctx context.Context, opts *bind.CallOpts, counter bind.Object, a uint64, b uint64) (*models.SuiTransactionBlockResponse, error)
	IncrementByBytesLength(ctx context.Context, opts *bind.CallOpts, counter bind.Object, bytes []byte) (*models.SuiTransactionBlockResponse, error)
	GetCount(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCountUsingPointer(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCountNoEntry(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCoinValue(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coin bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetAddressList(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetSimpleResult(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetResultStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetMultiNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetTupleStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetOcrConfig(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetVectorOfU8(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetVectorOfAddresses(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetVectorOfVectorsOfU8(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ICounterDevInspect
	Encoder() CounterEncoder
}

type ICounterDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	Create(ctx context.Context, opts *bind.CallOpts) (bind.Object, error)
	IncrementByOne(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error)
	IncrementByOneNoContext(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error)
	GetCount(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error)
	GetCountUsingPointer(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error)
	GetCountNoEntry(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error)
	GetCoinValue(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coin bind.Object) (uint64, error)
	GetAddressList(ctx context.Context, opts *bind.CallOpts) (AddressList, error)
	GetSimpleResult(ctx context.Context, opts *bind.CallOpts) (SimpleResult, error)
	GetResultStruct(ctx context.Context, opts *bind.CallOpts) (ComplexResult, error)
	GetNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (NestedStruct, error)
	GetMultiNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (MultiNestedStruct, error)
	GetTupleStruct(ctx context.Context, opts *bind.CallOpts) ([]any, error)
	GetOcrConfig(ctx context.Context, opts *bind.CallOpts) (OCRConfig, error)
	GetVectorOfU8(ctx context.Context, opts *bind.CallOpts) ([]byte, error)
	GetVectorOfAddresses(ctx context.Context, opts *bind.CallOpts) ([]string, error)
	GetVectorOfVectorsOfU8(ctx context.Context, opts *bind.CallOpts) ([][]byte, error)
}

type CounterEncoder interface {
	Initialize() (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Increment(counter bind.Object) (*bind.EncodedCall, error)
	IncrementWithArgs(args ...any) (*bind.EncodedCall, error)
	Decrement(counter bind.Object) (*bind.EncodedCall, error)
	DecrementWithArgs(args ...any) (*bind.EncodedCall, error)
	Create() (*bind.EncodedCall, error)
	CreateWithArgs(args ...any) (*bind.EncodedCall, error)
	IncrementByOne(counter bind.Object) (*bind.EncodedCall, error)
	IncrementByOneWithArgs(args ...any) (*bind.EncodedCall, error)
	IncrementByOneNoContext(counter bind.Object) (*bind.EncodedCall, error)
	IncrementByOneNoContextWithArgs(args ...any) (*bind.EncodedCall, error)
	IncrementByTwo(admin bind.Object, counter bind.Object) (*bind.EncodedCall, error)
	IncrementByTwoWithArgs(args ...any) (*bind.EncodedCall, error)
	IncrementByTwoNoContext(admin bind.Object, counter bind.Object) (*bind.EncodedCall, error)
	IncrementByTwoNoContextWithArgs(args ...any) (*bind.EncodedCall, error)
	IncrementBy(counter bind.Object, by uint64) (*bind.EncodedCall, error)
	IncrementByWithArgs(args ...any) (*bind.EncodedCall, error)
	IncrementMult(counter bind.Object, a uint64, b uint64) (*bind.EncodedCall, error)
	IncrementMultWithArgs(args ...any) (*bind.EncodedCall, error)
	IncrementByBytesLength(counter bind.Object, bytes []byte) (*bind.EncodedCall, error)
	IncrementByBytesLengthWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCount(counter bind.Object) (*bind.EncodedCall, error)
	GetCountWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCountUsingPointer(counter bind.Object) (*bind.EncodedCall, error)
	GetCountUsingPointerWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCountNoEntry(counter bind.Object) (*bind.EncodedCall, error)
	GetCountNoEntryWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCoinValue(typeArgs []string, coin bind.Object) (*bind.EncodedCall, error)
	GetCoinValueWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetAddressList() (*bind.EncodedCall, error)
	GetAddressListWithArgs(args ...any) (*bind.EncodedCall, error)
	GetSimpleResult() (*bind.EncodedCall, error)
	GetSimpleResultWithArgs(args ...any) (*bind.EncodedCall, error)
	GetResultStruct() (*bind.EncodedCall, error)
	GetResultStructWithArgs(args ...any) (*bind.EncodedCall, error)
	GetNestedResultStruct() (*bind.EncodedCall, error)
	GetNestedResultStructWithArgs(args ...any) (*bind.EncodedCall, error)
	GetMultiNestedResultStruct() (*bind.EncodedCall, error)
	GetMultiNestedResultStructWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTupleStruct() (*bind.EncodedCall, error)
	GetTupleStructWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOcrConfig() (*bind.EncodedCall, error)
	GetOcrConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetVectorOfU8() (*bind.EncodedCall, error)
	GetVectorOfU8WithArgs(args ...any) (*bind.EncodedCall, error)
	GetVectorOfAddresses() (*bind.EncodedCall, error)
	GetVectorOfAddressesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetVectorOfVectorsOfU8() (*bind.EncodedCall, error)
	GetVectorOfVectorsOfU8WithArgs(args ...any) (*bind.EncodedCall, error)
}

type CounterContract struct {
	*bind.BoundContract
	counterEncoder
	devInspect *CounterDevInspect
}

type CounterDevInspect struct {
	contract *CounterContract
}

var _ ICounter = (*CounterContract)(nil)
var _ ICounterDevInspect = (*CounterDevInspect)(nil)

func NewCounter(packageID string, client sui.ISuiAPI) (*CounterContract, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "counter", client)
	if err != nil {
		return nil, err
	}

	c := &CounterContract{
		BoundContract:  contract,
		counterEncoder: counterEncoder{BoundContract: contract},
	}
	c.devInspect = &CounterDevInspect{contract: c}
	return c, nil
}

func (c *CounterContract) Encoder() CounterEncoder {
	return c.counterEncoder
}

func (c *CounterContract) DevInspect() ICounterDevInspect {
	return c.devInspect
}

type COUNTER struct {
}

type CounterIncremented struct {
	CounterId bind.Object `move:"ID"`
	NewValue  uint64      `move:"u64"`
}

type CounterDecremented struct {
	EventType string      `move:"0x1::string::String"`
	CounterId bind.Object `move:"ID"`
	NewValue  uint64      `move:"u64"`
}

type AdminCap struct {
	Id string `move:"sui::object::UID"`
}

type Counter struct {
	Id    string `move:"sui::object::UID"`
	Value uint64 `move:"u64"`
}

type CounterPointer struct {
	Id         string `move:"sui::object::UID"`
	CounterId  string `move:"address"`
	AdminCapId string `move:"address"`
}

type AddressList struct {
	Addresses []string `move:"vector<address>"`
	Count     uint64   `move:"u64"`
}

type SimpleResult struct {
	Value uint64 `move:"u64"`
}

type ComplexResult struct {
	Count     uint64 `move:"u64"`
	Addr      string `move:"address"`
	IsComplex bool   `move:"bool"`
	Bytes     []byte `move:"vector<u8>"`
}

type NestedStruct struct {
	IsNested           bool          `move:"bool"`
	DoubleCount        uint64        `move:"u64"`
	NestedStruct       ComplexResult `move:"ComplexResult"`
	NestedSimpleStruct SimpleResult  `move:"SimpleResult"`
}

type MultiNestedStruct struct {
	IsMultiNested      bool         `move:"bool"`
	DoubleCount        uint64       `move:"u64"`
	NestedStruct       NestedStruct `move:"NestedStruct"`
	NestedSimpleStruct SimpleResult `move:"SimpleResult"`
}

type ConfigInfo struct {
	ConfigDigest                   []byte `move:"vector<u8>"`
	BigF                           byte   `move:"u8"`
	N                              byte   `move:"u8"`
	IsSignatureVerificationEnabled bool   `move:"bool"`
}

type OCRConfig struct {
	ConfigInfo   ConfigInfo `move:"ConfigInfo"`
	Signers      [][]byte   `move:"vector<vector<u8>>"`
	Transmitters []string   `move:"vector<address>"`
}

type bcsCounterPointer struct {
	Id         string
	CounterId  [32]byte
	AdminCapId [32]byte
}

func convertCounterPointerFromBCS(bcs bcsCounterPointer) CounterPointer {
	return CounterPointer{
		Id:         bcs.Id,
		CounterId:  fmt.Sprintf("0x%x", bcs.CounterId),
		AdminCapId: fmt.Sprintf("0x%x", bcs.AdminCapId),
	}
}

type bcsAddressList struct {
	Addresses [][32]byte
	Count     uint64
}

func convertAddressListFromBCS(bcs bcsAddressList) AddressList {
	return AddressList{
		Addresses: func() []string {
			addrs := make([]string, len(bcs.Addresses))
			for i, addr := range bcs.Addresses {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		Count: bcs.Count,
	}
}

type bcsComplexResult struct {
	Count     uint64
	Addr      [32]byte
	IsComplex bool
	Bytes     []byte
}

func convertComplexResultFromBCS(bcs bcsComplexResult) ComplexResult {
	return ComplexResult{
		Count:     bcs.Count,
		Addr:      fmt.Sprintf("0x%x", bcs.Addr),
		IsComplex: bcs.IsComplex,
		Bytes:     bcs.Bytes,
	}
}

type bcsNestedStruct struct {
	IsNested           bool
	DoubleCount        uint64
	NestedStruct       bcsComplexResult
	NestedSimpleStruct SimpleResult
}

func convertNestedStructFromBCS(bcs bcsNestedStruct) NestedStruct {
	return NestedStruct{
		IsNested:           bcs.IsNested,
		DoubleCount:        bcs.DoubleCount,
		NestedStruct:       convertComplexResultFromBCS(bcs.NestedStruct),
		NestedSimpleStruct: bcs.NestedSimpleStruct,
	}
}

type bcsMultiNestedStruct struct {
	IsMultiNested      bool
	DoubleCount        uint64
	NestedStruct       bcsNestedStruct
	NestedSimpleStruct SimpleResult
}

func convertMultiNestedStructFromBCS(bcs bcsMultiNestedStruct) MultiNestedStruct {
	return MultiNestedStruct{
		IsMultiNested:      bcs.IsMultiNested,
		DoubleCount:        bcs.DoubleCount,
		NestedStruct:       convertNestedStructFromBCS(bcs.NestedStruct),
		NestedSimpleStruct: bcs.NestedSimpleStruct,
	}
}

type bcsOCRConfig struct {
	ConfigInfo   ConfigInfo
	Signers      [][]byte
	Transmitters [][32]byte
}

func convertOCRConfigFromBCS(bcs bcsOCRConfig) OCRConfig {
	return OCRConfig{
		ConfigInfo: bcs.ConfigInfo,
		Signers:    bcs.Signers,
		Transmitters: func() []string {
			addrs := make([]string, len(bcs.Transmitters))
			for i, addr := range bcs.Transmitters {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}
}

func init() {
	bind.RegisterStructDecoder("test::counter::COUNTER", func(data []byte) (interface{}, error) {
		var result COUNTER
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::CounterIncremented", func(data []byte) (interface{}, error) {
		var result CounterIncremented
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::CounterDecremented", func(data []byte) (interface{}, error) {
		var result CounterDecremented
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::AdminCap", func(data []byte) (interface{}, error) {
		var result AdminCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::Counter", func(data []byte) (interface{}, error) {
		var result Counter
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::CounterPointer", func(data []byte) (interface{}, error) {
		var temp bcsCounterPointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertCounterPointerFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::AddressList", func(data []byte) (interface{}, error) {
		var temp bcsAddressList
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertAddressListFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::SimpleResult", func(data []byte) (interface{}, error) {
		var result SimpleResult
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::ComplexResult", func(data []byte) (interface{}, error) {
		var temp bcsComplexResult
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertComplexResultFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::NestedStruct", func(data []byte) (interface{}, error) {
		var temp bcsNestedStruct
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertNestedStructFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::MultiNestedStruct", func(data []byte) (interface{}, error) {
		var temp bcsMultiNestedStruct
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertMultiNestedStructFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::ConfigInfo", func(data []byte) (interface{}, error) {
		var result ConfigInfo
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::counter::OCRConfig", func(data []byte) (interface{}, error) {
		var temp bcsOCRConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertOCRConfigFromBCS(temp)
		return result, nil
	})
}

// Initialize executes the initialize Move function.
func (c *CounterContract) Initialize(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function.
func (c *CounterContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Increment executes the increment Move function.
func (c *CounterContract) Increment(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.Increment(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Decrement executes the decrement Move function.
func (c *CounterContract) Decrement(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.Decrement(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Create executes the create Move function.
func (c *CounterContract) Create(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.Create()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementByOne executes the increment_by_one Move function.
func (c *CounterContract) IncrementByOne(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.IncrementByOne(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementByOneNoContext executes the increment_by_one_no_context Move function.
func (c *CounterContract) IncrementByOneNoContext(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.IncrementByOneNoContext(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementByTwo executes the increment_by_two Move function.
func (c *CounterContract) IncrementByTwo(ctx context.Context, opts *bind.CallOpts, admin bind.Object, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.IncrementByTwo(admin, counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementByTwoNoContext executes the increment_by_two_no_context Move function.
func (c *CounterContract) IncrementByTwoNoContext(ctx context.Context, opts *bind.CallOpts, admin bind.Object, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.IncrementByTwoNoContext(admin, counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementBy executes the increment_by Move function.
func (c *CounterContract) IncrementBy(ctx context.Context, opts *bind.CallOpts, counter bind.Object, by uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.IncrementBy(counter, by)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementMult executes the increment_mult Move function.
func (c *CounterContract) IncrementMult(ctx context.Context, opts *bind.CallOpts, counter bind.Object, a uint64, b uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.IncrementMult(counter, a, b)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementByBytesLength executes the increment_by_bytes_length Move function.
func (c *CounterContract) IncrementByBytesLength(ctx context.Context, opts *bind.CallOpts, counter bind.Object, bytes []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.IncrementByBytesLength(counter, bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCount executes the get_count Move function.
func (c *CounterContract) GetCount(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetCount(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCountUsingPointer executes the get_count_using_pointer Move function.
func (c *CounterContract) GetCountUsingPointer(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetCountUsingPointer(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCountNoEntry executes the get_count_no_entry Move function.
func (c *CounterContract) GetCountNoEntry(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetCountNoEntry(counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCoinValue executes the get_coin_value Move function.
func (c *CounterContract) GetCoinValue(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coin bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetCoinValue(typeArgs, coin)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAddressList executes the get_address_list Move function.
func (c *CounterContract) GetAddressList(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetAddressList()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSimpleResult executes the get_simple_result Move function.
func (c *CounterContract) GetSimpleResult(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetSimpleResult()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetResultStruct executes the get_result_struct Move function.
func (c *CounterContract) GetResultStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetResultStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetNestedResultStruct executes the get_nested_result_struct Move function.
func (c *CounterContract) GetNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetNestedResultStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetMultiNestedResultStruct executes the get_multi_nested_result_struct Move function.
func (c *CounterContract) GetMultiNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetMultiNestedResultStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTupleStruct executes the get_tuple_struct Move function.
func (c *CounterContract) GetTupleStruct(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetTupleStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOcrConfig executes the get_ocr_config Move function.
func (c *CounterContract) GetOcrConfig(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetOcrConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetVectorOfU8 executes the get_vector_of_u8 Move function.
func (c *CounterContract) GetVectorOfU8(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetVectorOfU8()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetVectorOfAddresses executes the get_vector_of_addresses Move function.
func (c *CounterContract) GetVectorOfAddresses(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetVectorOfAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetVectorOfVectorsOfU8 executes the get_vector_of_vectors_of_u8 Move function.
func (c *CounterContract) GetVectorOfVectorsOfU8(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.counterEncoder.GetVectorOfVectorsOfU8()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *CounterDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.counterEncoder.TypeAndVersion()
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

// Create executes the create Move function using DevInspect to get return values.
//
// Returns: Counter
func (d *CounterDevInspect) Create(ctx context.Context, opts *bind.CallOpts) (bind.Object, error) {
	encoded, err := d.contract.counterEncoder.Create()
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
	result, ok := results[0].(bind.Object)
	if !ok {
		return bind.Object{}, fmt.Errorf("unexpected return type: expected bind.Object, got %T", results[0])
	}
	return result, nil
}

// IncrementByOne executes the increment_by_one Move function using DevInspect to get return values.
//
// Returns: u64
func (d *CounterDevInspect) IncrementByOne(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error) {
	encoded, err := d.contract.counterEncoder.IncrementByOne(counter)
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

// IncrementByOneNoContext executes the increment_by_one_no_context Move function using DevInspect to get return values.
//
// Returns: u64
func (d *CounterDevInspect) IncrementByOneNoContext(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error) {
	encoded, err := d.contract.counterEncoder.IncrementByOneNoContext(counter)
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

// GetCount executes the get_count Move function using DevInspect to get return values.
//
// Returns: u64
func (d *CounterDevInspect) GetCount(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error) {
	encoded, err := d.contract.counterEncoder.GetCount(counter)
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

// GetCountUsingPointer executes the get_count_using_pointer Move function using DevInspect to get return values.
//
// Returns: u64
func (d *CounterDevInspect) GetCountUsingPointer(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error) {
	encoded, err := d.contract.counterEncoder.GetCountUsingPointer(counter)
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

// GetCountNoEntry executes the get_count_no_entry Move function using DevInspect to get return values.
//
// Returns: u64
func (d *CounterDevInspect) GetCountNoEntry(ctx context.Context, opts *bind.CallOpts, counter bind.Object) (uint64, error) {
	encoded, err := d.contract.counterEncoder.GetCountNoEntry(counter)
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

// GetCoinValue executes the get_coin_value Move function using DevInspect to get return values.
//
// Returns: u64
func (d *CounterDevInspect) GetCoinValue(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coin bind.Object) (uint64, error) {
	encoded, err := d.contract.counterEncoder.GetCoinValue(typeArgs, coin)
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

// GetAddressList executes the get_address_list Move function using DevInspect to get return values.
//
// Returns: AddressList
func (d *CounterDevInspect) GetAddressList(ctx context.Context, opts *bind.CallOpts) (AddressList, error) {
	encoded, err := d.contract.counterEncoder.GetAddressList()
	if err != nil {
		return AddressList{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return AddressList{}, err
	}
	if len(results) == 0 {
		return AddressList{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(AddressList)
	if !ok {
		return AddressList{}, fmt.Errorf("unexpected return type: expected AddressList, got %T", results[0])
	}
	return result, nil
}

// GetSimpleResult executes the get_simple_result Move function using DevInspect to get return values.
//
// Returns: SimpleResult
func (d *CounterDevInspect) GetSimpleResult(ctx context.Context, opts *bind.CallOpts) (SimpleResult, error) {
	encoded, err := d.contract.counterEncoder.GetSimpleResult()
	if err != nil {
		return SimpleResult{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return SimpleResult{}, err
	}
	if len(results) == 0 {
		return SimpleResult{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(SimpleResult)
	if !ok {
		return SimpleResult{}, fmt.Errorf("unexpected return type: expected SimpleResult, got %T", results[0])
	}
	return result, nil
}

// GetResultStruct executes the get_result_struct Move function using DevInspect to get return values.
//
// Returns: ComplexResult
func (d *CounterDevInspect) GetResultStruct(ctx context.Context, opts *bind.CallOpts) (ComplexResult, error) {
	encoded, err := d.contract.counterEncoder.GetResultStruct()
	if err != nil {
		return ComplexResult{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return ComplexResult{}, err
	}
	if len(results) == 0 {
		return ComplexResult{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(ComplexResult)
	if !ok {
		return ComplexResult{}, fmt.Errorf("unexpected return type: expected ComplexResult, got %T", results[0])
	}
	return result, nil
}

// GetNestedResultStruct executes the get_nested_result_struct Move function using DevInspect to get return values.
//
// Returns: NestedStruct
func (d *CounterDevInspect) GetNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (NestedStruct, error) {
	encoded, err := d.contract.counterEncoder.GetNestedResultStruct()
	if err != nil {
		return NestedStruct{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return NestedStruct{}, err
	}
	if len(results) == 0 {
		return NestedStruct{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(NestedStruct)
	if !ok {
		return NestedStruct{}, fmt.Errorf("unexpected return type: expected NestedStruct, got %T", results[0])
	}
	return result, nil
}

// GetMultiNestedResultStruct executes the get_multi_nested_result_struct Move function using DevInspect to get return values.
//
// Returns: MultiNestedStruct
func (d *CounterDevInspect) GetMultiNestedResultStruct(ctx context.Context, opts *bind.CallOpts) (MultiNestedStruct, error) {
	encoded, err := d.contract.counterEncoder.GetMultiNestedResultStruct()
	if err != nil {
		return MultiNestedStruct{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return MultiNestedStruct{}, err
	}
	if len(results) == 0 {
		return MultiNestedStruct{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(MultiNestedStruct)
	if !ok {
		return MultiNestedStruct{}, fmt.Errorf("unexpected return type: expected MultiNestedStruct, got %T", results[0])
	}
	return result, nil
}

// GetTupleStruct executes the get_tuple_struct Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u64
//	[1]: address
//	[2]: bool
//	[3]: MultiNestedStruct
func (d *CounterDevInspect) GetTupleStruct(ctx context.Context, opts *bind.CallOpts) ([]any, error) {
	encoded, err := d.contract.counterEncoder.GetTupleStruct()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetOcrConfig executes the get_ocr_config Move function using DevInspect to get return values.
//
// Returns: OCRConfig
func (d *CounterDevInspect) GetOcrConfig(ctx context.Context, opts *bind.CallOpts) (OCRConfig, error) {
	encoded, err := d.contract.counterEncoder.GetOcrConfig()
	if err != nil {
		return OCRConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return OCRConfig{}, err
	}
	if len(results) == 0 {
		return OCRConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(OCRConfig)
	if !ok {
		return OCRConfig{}, fmt.Errorf("unexpected return type: expected OCRConfig, got %T", results[0])
	}
	return result, nil
}

// GetVectorOfU8 executes the get_vector_of_u8 Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *CounterDevInspect) GetVectorOfU8(ctx context.Context, opts *bind.CallOpts) ([]byte, error) {
	encoded, err := d.contract.counterEncoder.GetVectorOfU8()
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

// GetVectorOfAddresses executes the get_vector_of_addresses Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *CounterDevInspect) GetVectorOfAddresses(ctx context.Context, opts *bind.CallOpts) ([]string, error) {
	encoded, err := d.contract.counterEncoder.GetVectorOfAddresses()
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

// GetVectorOfVectorsOfU8 executes the get_vector_of_vectors_of_u8 Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *CounterDevInspect) GetVectorOfVectorsOfU8(ctx context.Context, opts *bind.CallOpts) ([][]byte, error) {
	encoded, err := d.contract.counterEncoder.GetVectorOfVectorsOfU8()
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

type counterEncoder struct {
	*bind.BoundContract
}

// Initialize encodes a call to the initialize Move function.
func (c counterEncoder) Initialize() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c counterEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// Increment encodes a call to the increment Move function.
func (c counterEncoder) Increment(counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment", typeArgsList, typeParamsList, []string{
		"&mut Counter",
	}, []any{
		counter,
	}, nil)
}

// IncrementWithArgs encodes a call to the increment Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Decrement encodes a call to the decrement Move function.
func (c counterEncoder) Decrement(counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("decrement", typeArgsList, typeParamsList, []string{
		"&mut Counter",
	}, []any{
		counter,
	}, nil)
}

// DecrementWithArgs encodes a call to the decrement Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) DecrementWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("decrement", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Create encodes a call to the create Move function.
func (c counterEncoder) Create() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"test::counter::Counter",
	})
}

// CreateWithArgs encodes a call to the create Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) CreateWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::counter::Counter",
	})
}

// IncrementByOne encodes a call to the increment_by_one Move function.
func (c counterEncoder) IncrementByOne(counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_one", typeArgsList, typeParamsList, []string{
		"&mut Counter",
	}, []any{
		counter,
	}, []string{
		"u64",
	})
}

// IncrementByOneWithArgs encodes a call to the increment_by_one Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementByOneWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_one", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// IncrementByOneNoContext encodes a call to the increment_by_one_no_context Move function.
func (c counterEncoder) IncrementByOneNoContext(counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_one_no_context", typeArgsList, typeParamsList, []string{
		"&mut Counter",
	}, []any{
		counter,
	}, []string{
		"u64",
	})
}

// IncrementByOneNoContextWithArgs encodes a call to the increment_by_one_no_context Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementByOneNoContextWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_one_no_context", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// IncrementByTwo encodes a call to the increment_by_two Move function.
func (c counterEncoder) IncrementByTwo(admin bind.Object, counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_two", typeArgsList, typeParamsList, []string{
		"&AdminCap",
		"&mut Counter",
	}, []any{
		admin,
		counter,
	}, nil)
}

// IncrementByTwoWithArgs encodes a call to the increment_by_two Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementByTwoWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&AdminCap",
		"&mut Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_two", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IncrementByTwoNoContext encodes a call to the increment_by_two_no_context Move function.
func (c counterEncoder) IncrementByTwoNoContext(admin bind.Object, counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_two_no_context", typeArgsList, typeParamsList, []string{
		"&AdminCap",
		"&mut Counter",
	}, []any{
		admin,
		counter,
	}, nil)
}

// IncrementByTwoNoContextWithArgs encodes a call to the increment_by_two_no_context Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementByTwoNoContextWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&AdminCap",
		"&mut Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_two_no_context", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IncrementBy encodes a call to the increment_by Move function.
func (c counterEncoder) IncrementBy(counter bind.Object, by uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by", typeArgsList, typeParamsList, []string{
		"&mut Counter",
		"u64",
	}, []any{
		counter,
		by,
	}, nil)
}

// IncrementByWithArgs encodes a call to the increment_by Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementByWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Counter",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IncrementMult encodes a call to the increment_mult Move function.
func (c counterEncoder) IncrementMult(counter bind.Object, a uint64, b uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_mult", typeArgsList, typeParamsList, []string{
		"&mut Counter",
		"u64",
		"u64",
	}, []any{
		counter,
		a,
		b,
	}, nil)
}

// IncrementMultWithArgs encodes a call to the increment_mult Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementMultWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Counter",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_mult", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IncrementByBytesLength encodes a call to the increment_by_bytes_length Move function.
func (c counterEncoder) IncrementByBytesLength(counter bind.Object, bytes []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_bytes_length", typeArgsList, typeParamsList, []string{
		"&mut Counter",
		"vector<u8>",
	}, []any{
		counter,
		bytes,
	}, nil)
}

// IncrementByBytesLengthWithArgs encodes a call to the increment_by_bytes_length Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) IncrementByBytesLengthWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Counter",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("increment_by_bytes_length", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCount encodes a call to the get_count Move function.
func (c counterEncoder) GetCount(counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_count", typeArgsList, typeParamsList, []string{
		"&Counter",
	}, []any{
		counter,
	}, []string{
		"u64",
	})
}

// GetCountWithArgs encodes a call to the get_count Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetCountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_count", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetCountUsingPointer encodes a call to the get_count_using_pointer Move function.
func (c counterEncoder) GetCountUsingPointer(counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_count_using_pointer", typeArgsList, typeParamsList, []string{
		"&Counter",
	}, []any{
		counter,
	}, []string{
		"u64",
	})
}

// GetCountUsingPointerWithArgs encodes a call to the get_count_using_pointer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetCountUsingPointerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_count_using_pointer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetCountNoEntry encodes a call to the get_count_no_entry Move function.
func (c counterEncoder) GetCountNoEntry(counter bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_count_no_entry", typeArgsList, typeParamsList, []string{
		"&Counter",
	}, []any{
		counter,
	}, []string{
		"u64",
	})
}

// GetCountNoEntryWithArgs encodes a call to the get_count_no_entry Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetCountNoEntryWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Counter",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_count_no_entry", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetCoinValue encodes a call to the get_coin_value Move function.
func (c counterEncoder) GetCoinValue(typeArgs []string, coin bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_coin_value", typeArgsList, typeParamsList, []string{
		"&Coin<T>",
	}, []any{
		coin,
	}, []string{
		"u64",
	})
}

// GetCoinValueWithArgs encodes a call to the get_coin_value Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetCoinValueWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Coin<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_coin_value", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetAddressList encodes a call to the get_address_list Move function.
func (c counterEncoder) GetAddressList() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_address_list", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"test::counter::AddressList",
	})
}

// GetAddressListWithArgs encodes a call to the get_address_list Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetAddressListWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_address_list", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::counter::AddressList",
	})
}

// GetSimpleResult encodes a call to the get_simple_result Move function.
func (c counterEncoder) GetSimpleResult() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_simple_result", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"test::counter::SimpleResult",
	})
}

// GetSimpleResultWithArgs encodes a call to the get_simple_result Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetSimpleResultWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_simple_result", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::counter::SimpleResult",
	})
}

// GetResultStruct encodes a call to the get_result_struct Move function.
func (c counterEncoder) GetResultStruct() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_result_struct", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"test::counter::ComplexResult",
	})
}

// GetResultStructWithArgs encodes a call to the get_result_struct Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetResultStructWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_result_struct", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::counter::ComplexResult",
	})
}

// GetNestedResultStruct encodes a call to the get_nested_result_struct Move function.
func (c counterEncoder) GetNestedResultStruct() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_nested_result_struct", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"test::counter::NestedStruct",
	})
}

// GetNestedResultStructWithArgs encodes a call to the get_nested_result_struct Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetNestedResultStructWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_nested_result_struct", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::counter::NestedStruct",
	})
}

// GetMultiNestedResultStruct encodes a call to the get_multi_nested_result_struct Move function.
func (c counterEncoder) GetMultiNestedResultStruct() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_multi_nested_result_struct", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"test::counter::MultiNestedStruct",
	})
}

// GetMultiNestedResultStructWithArgs encodes a call to the get_multi_nested_result_struct Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetMultiNestedResultStructWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_multi_nested_result_struct", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::counter::MultiNestedStruct",
	})
}

// GetTupleStruct encodes a call to the get_tuple_struct Move function.
func (c counterEncoder) GetTupleStruct() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_tuple_struct", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"u64",
		"address",
		"bool",
		"test::counter::MultiNestedStruct",
	})
}

// GetTupleStructWithArgs encodes a call to the get_tuple_struct Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetTupleStructWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_tuple_struct", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
		"address",
		"bool",
		"test::counter::MultiNestedStruct",
	})
}

// GetOcrConfig encodes a call to the get_ocr_config Move function.
func (c counterEncoder) GetOcrConfig() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ocr_config", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"test::counter::OCRConfig",
	})
}

// GetOcrConfigWithArgs encodes a call to the get_ocr_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetOcrConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ocr_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::counter::OCRConfig",
	})
}

// GetVectorOfU8 encodes a call to the get_vector_of_u8 Move function.
func (c counterEncoder) GetVectorOfU8() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_vector_of_u8", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"vector<u8>",
	})
}

// GetVectorOfU8WithArgs encodes a call to the get_vector_of_u8 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetVectorOfU8WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_vector_of_u8", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// GetVectorOfAddresses encodes a call to the get_vector_of_addresses Move function.
func (c counterEncoder) GetVectorOfAddresses() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_vector_of_addresses", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"vector<address>",
	})
}

// GetVectorOfAddressesWithArgs encodes a call to the get_vector_of_addresses Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetVectorOfAddressesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_vector_of_addresses", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// GetVectorOfVectorsOfU8 encodes a call to the get_vector_of_vectors_of_u8 Move function.
func (c counterEncoder) GetVectorOfVectorsOfU8() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_vector_of_vectors_of_u8", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"vector<vector<u8>>",
	})
}

// GetVectorOfVectorsOfU8WithArgs encodes a call to the get_vector_of_vectors_of_u8 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c counterEncoder) GetVectorOfVectorsOfU8WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_vector_of_vectors_of_u8", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
	})
}
