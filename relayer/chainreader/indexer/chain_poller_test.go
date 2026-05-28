package indexer

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// mockSuiPTBClient is a mock implementation of client.SuiPTBClient
type mockSuiPTBClient struct {
	mock.Mock
}

func (m *mockSuiPTBClient) MoveCall(ctx context.Context, req client.MoveCallRequest) (client.TxnMetaData, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(client.TxnMetaData), args.Error(1)
}

func (m *mockSuiPTBClient) SendTransaction(ctx context.Context, execRequest *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error) {
	args := m.Called(ctx, execRequest)
	return args.Get(0).(*suirpcv2.ExecuteTransactionResponse), args.Error(1)
}

func (m *mockSuiPTBClient) ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor []byte) ([]*suirpcv2.Object, error) {
	args := m.Called(ctx, ownerAddress, cursor)
	return args.Get(0).([]*suirpcv2.Object), args.Error(1)
}

func (m *mockSuiPTBClient) ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor []byte) ([]*suirpcv2.Object, error) {
	args := m.Called(ctx, ownerAddress, structType, cursor)
	return args.Get(0).([]*suirpcv2.Object), args.Error(1)
}

func (m *mockSuiPTBClient) ReadObjectId(ctx context.Context, objectId string) (*suirpcv2.Object, error) {
	args := m.Called(ctx, objectId)
	return args.Get(0).(*suirpcv2.Object), args.Error(1)
}

func (m *mockSuiPTBClient) ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	mArgs := m.Called(ctx, packageId, module, function, args, argTypes, typeArgs)
	return mArgs.Get(0).([]any), mArgs.Error(1)
}

func (m *mockSuiPTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte) (*suirpcv2.ExecuteTransactionResponse, error) {
	args := m.Called(ctx, txBytesRaw, signerPublicKey)
	return args.Get(0).(*suirpcv2.ExecuteTransactionResponse), args.Error(1)
}

func (m *mockSuiPTBClient) QueryEvents(ctx context.Context, filter client.EventFilterByMoveEventModule, limit *uint, cursor *client.EventId, sortOptions *client.QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	args := m.Called(ctx, filter, limit, cursor, sortOptions)
	return args.Get(0).(*models.PaginatedEventsResponse), args.Error(1)
}

func (m *mockSuiPTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *suirpcv2.Checkpoint, limit *uint64) ([]*suirpcv2.ExecutedTransaction, error) {
	args := m.Called(ctx, fromAddress, cursor, limit)
	return args.Get(0).([]*suirpcv2.ExecutedTransaction), args.Error(1)
}

func (m *mockSuiPTBClient) GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error) {
	args := m.Called(ctx, digest)
	return args.Get(0).(client.TransactionResult), args.Error(1)
}

func (m *mockSuiPTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]*suirpcv2.Object, error) {
	args := m.Called(ctx, address)
	return args.Get(0).([]*suirpcv2.Object), args.Error(1)
}

func (m *mockSuiPTBClient) QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]*suirpcv2.Object, error) {
	args := m.Called(ctx, address, coinType)
	return args.Get(0).([]*suirpcv2.Object), args.Error(1)
}

func (m *mockSuiPTBClient) EstimateGas(ctx context.Context, tx *transaction.Transaction) (uint64, error) {
	args := m.Called(ctx, tx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockSuiPTBClient) GetReferenceGasPrice(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *mockSuiPTBClient) FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType client.TransactionRequestType) (*suirpcv2.ExecuteTransactionResponse, error) {
	args := m.Called(ctx, txnSigner, tx, requestType)
	return args.Get(0).(*suirpcv2.ExecuteTransactionResponse), args.Error(1)
}

func (m *mockSuiPTBClient) BlockByDigest(ctx context.Context, txDigest string) (*suirpcv2.Checkpoint, error) {
	args := m.Called(ctx, txDigest)
	return args.Get(0).(*suirpcv2.Checkpoint), args.Error(1)
}

func (m *mockSuiPTBClient) GetBlockById(ctx context.Context, checkpointDigest string) (*suirpcv2.Checkpoint, error) {
	args := m.Called(ctx, checkpointDigest)
	return args.Get(0).(*suirpcv2.Checkpoint), args.Error(1)
}

func (m *mockSuiPTBClient) GetLatestEpoch(ctx context.Context) (*suirpcv2.Epoch, error) {
	args := m.Called(ctx)
	return args.Get(0).(*suirpcv2.Epoch), args.Error(1)
}

func (m *mockSuiPTBClient) GetLatestCheckpoint(ctx context.Context) (*suirpcv2.Checkpoint, error) {
	args := m.Called(ctx)
	return args.Get(0).(*suirpcv2.Checkpoint), args.Error(1)
}

func (m *mockSuiPTBClient) GetCheckpointData(ctx context.Context, checkpointSequenceNumber uint64) (*client.CheckpointData, error) {
	args := m.Called(ctx, checkpointSequenceNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.CheckpointData), args.Error(1)
}

func (m *mockSuiPTBClient) GetNormalizedModule(ctx context.Context, packageId string, moduleId string) (models.GetNormalizedMoveModuleResponse, error) {
	args := m.Called(ctx, packageId, moduleId)
	return args.Get(0).(models.GetNormalizedMoveModuleResponse), args.Error(1)
}

func (m *mockSuiPTBClient) GetSUIBalance(ctx context.Context, address string) (*suirpcv2.Balance, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(*suirpcv2.Balance), args.Error(1)
}

func (m *mockSuiPTBClient) LoadModulePackageIds(ctx context.Context, packageId string, module string) ([]string, error) {
	args := m.Called(ctx, packageId, module)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockSuiPTBClient) GetLatestPackageId(ctx context.Context, packageId string, module string) (string, error) {
	args := m.Called(ctx, packageId, module)
	return args.String(0), args.Error(1)
}

func (m *mockSuiPTBClient) GetClient() sui.ISuiAPI {
	args := m.Called()
	return args.Get(0).(sui.ISuiAPI)
}

func (m *mockSuiPTBClient) GetCache() *cache.Cache {
	args := m.Called()
	return args.Get(0).(*cache.Cache)
}

func (m *mockSuiPTBClient) GetCachedValue(key string) (any, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *mockSuiPTBClient) SetCachedValue(key string, value any) {
	m.Called(key, value)
}

func (m *mockSuiPTBClient) GetCachedValues(keys []string) (map[string]any, bool) {
	args := m.Called(keys)
	return args.Get(0).(map[string]any), args.Bool(1)
}

func (m *mockSuiPTBClient) SetCachedValues(keyValues map[string]any) {
	m.Called(keyValues)
}

func (m *mockSuiPTBClient) HashTxBytes(txBytes []byte) []byte {
	args := m.Called(txBytes)
	return args.Get(0).([]byte)
}

func (m *mockSuiPTBClient) GetCCIPPackageID(ctx context.Context, offRampPackageID string) (string, error) {
	args := m.Called(ctx, offRampPackageID)
	return args.String(0), args.Error(1)
}

func (m *mockSuiPTBClient) GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error) {
	args := m.Called(ctx, packageID, moduleID, objectName, fieldKeys)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *mockSuiPTBClient) GetParentObjectID(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error) {
	args := m.Called(ctx, packageID, moduleID, pointerObjectName)
	return args.String(0), args.Error(1)
}

func (m *mockSuiPTBClient) WithRateLimit(ctx context.Context, methodName string, f func(ctx context.Context) error) error {
	return f(ctx)
}

func TestNewChainPoller(t *testing.T) {
	log := logger.Test(t)
	mockClient := new(mockSuiPTBClient)

	cfg := config.ChainPollerConfig{
		PollingInterval:   100 * time.Millisecond,
		SyncTimeout:       10 * time.Second,
		ChannelBufferSize: 8,
	}

	selectorProvider := func() []*client.EventSelector {
		return []*client.EventSelector{
			{Package: "0x123", Module: "test", Event: "TestEvent"},
		}
	}

	poller := NewChainPoller(mockClient, log, cfg, selectorProvider)

	require.NotNil(t, poller)
	assert.NotNil(t, poller.EventsChannel())
	assert.NotNil(t, poller.TransactionsChannel())
}

func TestChainPoller_StartStop(t *testing.T) {
	log := logger.Test(t)
	mockClient := new(mockSuiPTBClient)

	// Mock GetLatestCheckpoint to return checkpoint 5
	mockClient.On("GetLatestCheckpoint", mock.Anything).Return(
		&suirpcv2.Checkpoint{SequenceNumber: ptr(uint64(5))},
		nil,
	)

	// Mock GetCheckpointData for checkpoints 5
	mockClient.On("GetCheckpointData", mock.Anything, uint64(5)).Return(
		&client.CheckpointData{
			Checkpoint: &suirpcv2.Checkpoint{
				SequenceNumber: ptr(uint64(5)),
				Digest:         ptr("checkpoint_digest_5"),
				Summary: &suirpcv2.CheckpointSummary{
					Timestamp: timestamppb.Now(),
				},
			},
			Transactions: []*suirpcv2.ExecutedTransaction{},
		},
		nil,
	).Maybe()

	cfg := config.ChainPollerConfig{
		PollingInterval:   1 * time.Second,
		SyncTimeout:       10 * time.Second,
		ChannelBufferSize: 8,
	}

	selectorProvider := func() []*client.EventSelector { return nil }

	poller := NewChainPoller(mockClient, log, cfg, selectorProvider)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := poller.Start(ctx)
	require.NoError(t, err)

	// Give it a moment to process
	time.Sleep(100 * time.Millisecond)

	err = poller.Close()
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func TestChainPoller_FilterEvents(t *testing.T) {
	log := logger.Test(t)
	mockClient := new(mockSuiPTBClient)

	cfg := config.ChainPollerConfig{
		PollingInterval:   1 * time.Second,
		SyncTimeout:       10 * time.Second,
		ChannelBufferSize: 8,
	}

	selectors := []*client.EventSelector{
		{Package: "0x123", Module: "test", Event: "MatchingEvent"},
	}

	selectorProvider := func() []*client.EventSelector {
		return selectors
	}

	poller := NewChainPoller(mockClient, log, cfg, selectorProvider)

	meta := CheckpointMeta{
		SequenceNumber: 1,
		Digest:         "test_digest",
		TimestampMs:    12345,
	}

	// Create transactions with events
	transactions := []*suirpcv2.ExecutedTransaction{
		{
			Digest: ptr("tx1"),
			Events: &suirpcv2.TransactionEvents{
				Events: []*suirpcv2.Event{
					{
						PackageId: ptr("0x123"),
						Module:    ptr("test"),
						EventType: ptr("0x123::test::MatchingEvent"),
					},
					{
						PackageId: ptr("0x456"),
						Module:    ptr("other"),
						EventType: ptr("0x456::other::OtherEvent"),
					},
				},
			},
		},
	}

	batch := poller.filterEvents(meta, transactions, selectors)

	require.Len(t, batch.Events, 1)
	assert.Equal(t, "tx1", batch.Events[0].TxDigest)
	assert.Equal(t, "0x123", batch.Events[0].Event.GetPackageId())
}

func TestChainPoller_ComputeStartSequence(t *testing.T) {
	log := logger.Test(t)

	tests := []struct {
		name                    string
		startCheckpointSequence *uint64
		backfillCheckpointCount *uint64
		latestCheckpoint        uint64
		expectedStart           uint64
	}{
		{
			name:             "no config - start from latest",
			latestCheckpoint: 100,
			expectedStart:    100,
		},
		{
			name:                    "explicit start checkpoint",
			startCheckpointSequence: ptr(uint64(50)),
			latestCheckpoint:        100,
			expectedStart:           50,
		},
		{
			name:                    "backfill count - latest - N",
			backfillCheckpointCount: ptr(uint64(10)),
			latestCheckpoint:        100,
			expectedStart:           90,
		},
		{
			name:                    "backfill count - clamps to 0",
			backfillCheckpointCount: ptr(uint64(200)),
			latestCheckpoint:        50,
			expectedStart:           0,
		},
		{
			name:                    "explicit start overrides backfill",
			startCheckpointSequence: ptr(uint64(75)),
			backfillCheckpointCount: ptr(uint64(10)),
			latestCheckpoint:        100,
			expectedStart:           75,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mockSuiPTBClient)
			mockClient.On("GetLatestCheckpoint", mock.Anything).Return(
				&suirpcv2.Checkpoint{SequenceNumber: ptr(tt.latestCheckpoint)},
				nil,
			).Once()

			cfg := config.ChainPollerConfig{
				PollingInterval:         1 * time.Second,
				SyncTimeout:             10 * time.Second,
				StartCheckpointSequence: tt.startCheckpointSequence,
				BackfillCheckpointCount: tt.backfillCheckpointCount,
				ChannelBufferSize:       8,
			}

			selectorProvider := func() []*client.EventSelector { return nil }
			poller := NewChainPoller(mockClient, log, cfg, selectorProvider)

			ctx := context.Background()
			startSeq, err := poller.computeStartSequence(ctx)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStart, startSeq)
			mockClient.AssertExpectations(t)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}