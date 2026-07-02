//go:build integration

package loop

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockContractReader is a mock implementation of types.ContractReader for testing.
// It simulates the LOOP boundary behavior where params and return values are []byte.
type mockContractReader struct {
	mock.Mock
}

func (m *mockContractReader) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockContractReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	args := m.Called(ctx, readIdentifier, confidenceLevel, params, returnVal)
	
	// Simulate the LOOP boundary: copy the mock return value to the returnVal pointer
	if ret := args.Get(0); ret != nil {
		// returnVal is expected to be *[]byte
		if rv, ok := returnVal.(*[]byte); ok {
			*rv = ret.([]byte)
		}
	}
	
	return args.Error(1)
}

func (m *mockContractReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(types.BatchGetLatestValuesResult), args.Error(1)
}

func (m *mockContractReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	args := m.Called(ctx, contract, filter, limitAndSort, sequenceDataType)
	
	// Simulate LOOP boundary: sequenceDataType is *[]byte, we return mock data
	if ret := args.Get(0); ret != nil {
		return ret.([]types.Sequence), nil
	}
	
	return nil, args.Error(1)
}

func (m *mockContractReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	args := m.Called(ctx, bindings)
	return args.Error(0)
}

func (m *mockContractReader) Unbind(ctx context.Context, bindings []types.BoundContract) error {
	args := m.Called(ctx, bindings)
	return args.Error(0)
}

func (m *mockContractReader) Start(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *mockContractReader) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockContractReader) Ready() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockContractReader) HealthReport() map[string]error {
	args := m.Called()
	return args.Get(0).(map[string]error)
}

// TestLoopChainReader_GetLatestValue_VerifiesLOOPBoundary tests that the LoopChainReader
// properly serializes params to []byte and deserializes results using DecodeSuiJsonValue.
func TestLoopChainReader_GetLatestValue_VerifiesLOOPBoundary(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	
	mockReader := new(mockContractReader)
	loopReader := NewLoopChainReader(log, mockReader)
	
	ctx := context.Background()
	readIdentifier := "test-contract-echo_u64"
	params := map[string]any{"val": uint64(42)}
	
	// Expected JSON params after serialization
	expectedParams, err := json.Marshal(params)
	require.NoError(t, err)
	
	// Mock return value (simulating what comes back from the LOOP boundary)
	mockReturn := json.RawMessage("42")
	
	// Set up expectation: GetLatestValue receives []byte params
	mockReader.On("GetLatestValue", 
		ctx, 
		readIdentifier, 
		primitives.Finalized, 
		mock.AnythingOfType("*[]uint8"), // params as *[]byte
		mock.AnythingOfType("*[]uint8"), // returnVal as *[]byte
	).Run(func(args mock.Arguments) {
		// Verify params were serialized correctly
		paramBytes := args.Get(3).(*[]byte)
		require.JSONEq(t, string(expectedParams), string(*paramBytes))
		
		// Simulate setting the return value (as would happen across LOOP boundary)
		returnVal := args.Get(4).(*[]byte)
		*returnVal = []byte(mockReturn)
	}).Return(nil)
	
	var result uint64
	err = loopReader.GetLatestValue(ctx, readIdentifier, primitives.Finalized, params, &result)
	
	require.NoError(t, err)
	require.Equal(t, uint64(42), result)
	mockReader.AssertExpectations(t)
}

// TestLoopChainReader_GetLatestValue_DecodesComplexTypes tests decoding of complex
// return types using the Sui JSON decoder.
func TestLoopChainReader_GetLatestValue_DecodesComplexTypes(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	
	mockReader := new(mockContractReader)
	loopReader := NewLoopChainReader(log, mockReader)
	
	ctx := context.Background()
	
	tests := []struct {
		name           string
		mockResponse   string
		expectedResult any
		targetType     any
	}{
		{
			name:           "uint64",
			mockResponse:   `100`,
			expectedResult: uint64(100),
			targetType:     new(uint64),
		},
		{
			name:           "string",
			mockResponse:   `"hello"`,
			expectedResult: "hello",
			targetType:     new(string),
		},
		{
			name:           "struct",
			mockResponse:   `{"first": 10, "second": 20}`,
			expectedResult: map[string]any{"first": float64(10), "second": float64(20)},
			targetType:     new(map[string]any),
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReader.ExpectedCalls = nil
			
			mockReader.On("GetLatestValue",
				ctx,
				"test-read",
				primitives.Finalized,
				mock.AnythingOfType("*[]uint8"),
				mock.AnythingOfType("*[]uint8"),
			).Run(func(args mock.Arguments) {
				returnVal := args.Get(4).(*[]byte)
				*returnVal = []byte(tt.mockResponse)
			}).Return(nil)
			
			switch v := tt.targetType.(type) {
			case *uint64:
				err := loopReader.GetLatestValue(ctx, "test-read", primitives.Finalized, map[string]any{}, v)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, *v)
			case *string:
				err := loopReader.GetLatestValue(ctx, "test-read", primitives.Finalized, map[string]any{}, v)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, *v)
			case *map[string]any:
				err := loopReader.GetLatestValue(ctx, "test-read", primitives.Finalized, map[string]any{}, v)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, *v)
			}
			
			mockReader.AssertExpectations(t)
		})
	}
}

// TestLoopChainReader_QueryKey_VerifiesExpressionSerialization tests that QueryKey
// properly serializes expressions for the LOOP boundary.
func TestLoopChainReader_QueryKey_VerifiesExpressionSerialization(t *testing.T) {
	t.Parallel()
	t.Skip("TODO: implement when QueryKey expression serialization is better understood")
	
	log := logger.Test(t)
	mockReader := new(mockContractReader)
	loopReader := NewLoopChainReader(log, mockReader)
	
	ctx := context.Background()
	contract := types.BoundContract{Name: "test", Address: "0x123"}
	filter := query.KeyFilter{
		Key: "test_event",
		Expressions: []query.Expression{
			// Add expression here
		},
	}
	limitAndSort := query.LimitAndSort{
		Limit: query.CountLimit(10),
	}
	
	// Mock event data as would come back from LOOP
	mockEventData := []byte(`{"value": 123}`)
	mockReader.On("QueryKey",
		ctx,
		contract,
		mock.Anything, // filter with serialized expressions
		limitAndSort,
		mock.AnythingOfType("*[]uint8"),
	).Return([]types.Sequence{
		{
			Data: &mockEventData,
		},
	}, nil)
	
	type TestEvent struct {
		Value uint64 `json:"value"`
	}
	
	var event TestEvent
	sequences, err := loopReader.QueryKey(ctx, contract, filter, limitAndSort, &event)
	
	require.NoError(t, err)
	require.Len(t, sequences, 1)
	require.Equal(t, uint64(123), event.Value)
}

// TestLoopChainReader_ProxiesCalls tests that all non-transforming calls are proxied
// directly to the underlying reader.
func TestLoopChainReader_ProxiesCalls(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	mockReader := new(mockContractReader)
	loopReader := NewLoopChainReader(log, mockReader)
	
	ctx := context.Background()
	
	t.Run("Name", func(t *testing.T) {
		mockReader.On("Name").Return("test-reader")
		name := loopReader.Name()
		require.Equal(t, "test-reader", name)
		mockReader.AssertExpectations(t)
	})
	
	t.Run("Ready", func(t *testing.T) {
		mockReader.On("Ready").Return(nil)
		err := loopReader.Ready()
		require.NoError(t, err)
		mockReader.AssertExpectations(t)
	})
	
	t.Run("HealthReport", func(t *testing.T) {
		report := map[string]error{"test": nil}
		mockReader.On("HealthReport").Return(report)
		result := loopReader.HealthReport()
		require.Equal(t, report, result)
		mockReader.AssertExpectations(t)
	})
	
	t.Run("Start", func(t *testing.T) {
		mockReader.On("Start", ctx).Return(nil)
		err := loopReader.Start(ctx)
		require.NoError(t, err)
		mockReader.AssertExpectations(t)
	})
	
	t.Run("Close", func(t *testing.T) {
		mockReader.On("Close").Return(nil)
		err := loopReader.Close()
		require.NoError(t, err)
		mockReader.AssertExpectations(t)
	})
	
	t.Run("Bind", func(t *testing.T) {
		bindings := []types.BoundContract{{Name: "test", Address: "0x123"}}
		mockReader.On("Bind", ctx, bindings).Return(nil)
		err := loopReader.Bind(ctx, bindings)
		require.NoError(t, err)
		mockReader.AssertExpectations(t)
	})
	
	t.Run("Unbind", func(t *testing.T) {
		bindings := []types.BoundContract{{Name: "test", Address: "0x123"}}
		mockReader.On("Unbind", ctx, bindings).Return(nil)
		err := loopReader.Unbind(ctx, bindings)
		require.NoError(t, err)
		mockReader.AssertExpectations(t)
	})
}

// TestLoopChainReader_BatchGetLatestValues_VerifiesLOOPBoundary tests that
// BatchGetLatestValues properly serializes all requests and deserializes responses.
func TestLoopChainReader_BatchGetLatestValues_VerifiesLOOPBoundary(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	mockReader := new(mockContractReader)
	loopReader := NewLoopChainReader(log, mockReader)
	
	ctx := context.Background()
	
	request := types.BatchGetLatestValuesRequest{
		types.BoundContract{Name: "contract1", Address: "0x1"}: {
			{ReadName: "read1", Params: map[string]any{"key": "value1"}},
		},
	}
	
	mockResponse := types.BatchGetLatestValuesResult{
		types.BoundContract{Name: "contract1", Address: "0x1"}: {
			{ReadName: "read1"},
		},
	}
	
	// The mock will receive converted requests with []byte params
	mockReader.On("BatchGetLatestValues", ctx, mock.Anything).Run(func(args mock.Arguments) {
		// Verify the request was converted to use []byte params
		convertedRequest := args.Get(1).(types.BatchGetLatestValuesRequest)
		for _, reads := range convertedRequest {
			for _, read := range reads {
				// Params should be []byte after serialization
				_, ok := read.Params.([]byte)
				require.True(t, ok, "params should be []byte after LOOP serialization")
			}
		}
	}).Return(mockResponse, nil)
	
	result, err := loopReader.BatchGetLatestValues(ctx, request)
	
	require.NoError(t, err)
	require.NotNil(t, result)
	mockReader.AssertExpectations(t)
}

// TestLoopChainReader_GetLatestValue_PropagateErrors tests that errors from the
// underlying reader are properly propagated.
func TestLoopChainReader_GetLatestValue_PropagateErrors(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	mockReader := new(mockContractReader)
	loopReader := NewLoopChainReader(log, mockReader)
	
	ctx := context.Background()
	expectedErr := errors.New("underlying reader error")
	
	mockReader.On("GetLatestValue",
		ctx,
		"test-read",
		primitives.Finalized,
		mock.AnythingOfType("*[]uint8"),
		mock.AnythingOfType("*[]uint8"),
	).Return(expectedErr)
	
	var result uint64
	err := loopReader.GetLatestValue(ctx, "test-read", primitives.Finalized, map[string]any{}, &result)
	
	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	mockReader.AssertExpectations(t)
}

// TestLoopChainReader_DecodeErrors tests that decode errors are wrapped properly.
func TestLoopChainReader_DecodeErrors(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	mockReader := new(mockContractReader)
	loopReader := NewLoopChainReader(log, mockReader)
	
	ctx := context.Background()
	
	// Return invalid JSON that can't be decoded into the target type
	mockReader.On("GetLatestValue",
		ctx,
		"test-read",
		primitives.Finalized,
		mock.AnythingOfType("*[]uint8"),
		mock.AnythingOfType("*[]uint8"),
	).Run(func(args mock.Arguments) {
		returnVal := args.Get(4).(*[]byte)
		*returnVal = []byte(`invalid json`)
	}).Return(nil)
	
	var result uint64
	err := loopReader.GetLatestValue(ctx, "test-read", primitives.Finalized, map[string]any{}, &result)
	
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode")
	mockReader.AssertExpectations(t)
}
