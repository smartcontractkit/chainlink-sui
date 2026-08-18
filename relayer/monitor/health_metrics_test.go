package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aptosTypes "github.com/smartcontractkit/chainlink-aptos/relayer/types"
)

func testChainInfo() aptosTypes.ChainInfo {
	return aptosTypes.ChainInfo{
		ChainFamilyName: "sui",
		ChainID:         "test-chain-id",
		NetworkName:     "testnet",
		NetworkNameFull: "Sui Testnet",
	}
}

func TestNewHealthMetrics(t *testing.T) {
	t.Parallel()

	hm, err := NewHealthMetrics(testChainInfo())
	require.NoError(t, err)
	require.NotNil(t, hm)

	// Verify all component metrics are initialized
	components := []string{
		ComponentSuiRelayer,
		ComponentTxM,
		ComponentChainPoller,
		ComponentEventsIndexer,
		ComponentTransactionsIndexer,
		ComponentChainReader,
		ComponentChainWriter,
	}

	for _, component := range components {
		// Check that initial health status is set to false (unhealthy)
		status, exists := hm.GetHealthStatus(component)
		assert.True(t, exists, "Component %s should have initial status", component)
		assert.False(t, status, "Initial status for %s should be unhealthy", component)

		// Check that initial last success time is set
		lastSuccess, exists := hm.GetLastSuccessTime(component)
		assert.True(t, exists, "Component %s should have initial last success time", component)
		assert.False(t, lastSuccess.IsZero(), "Last success time for %s should not be zero", component)
	}
}

func TestRecordHealth(t *testing.T) {
	t.Parallel()

	hm, err := NewHealthMetrics(testChainInfo())
	require.NoError(t, err)

	ctx := context.Background()
	component := ComponentTxM

	// Initial state should be unhealthy
	status, exists := hm.GetHealthStatus(component)
	assert.True(t, exists)
	assert.False(t, status)

	// Record healthy status
	hm.RecordHealth(ctx, component, true)
	status, exists = hm.GetHealthStatus(component)
	assert.True(t, exists)
	assert.True(t, status)

	// Record unhealthy status (flip-flop)
	hm.RecordHealth(ctx, component, false)
	status, exists = hm.GetHealthStatus(component)
	assert.True(t, exists)
	assert.False(t, status)

	// Record healthy again (another flip-flop)
	hm.RecordHealth(ctx, component, true)
	status, exists = hm.GetHealthStatus(component)
	assert.True(t, exists)
	assert.True(t, status)
}

func TestRecordLastSuccess(t *testing.T) {
	t.Parallel()

	hm, err := NewHealthMetrics(testChainInfo())
	require.NoError(t, err)

	ctx := context.Background()
	component := ComponentEventsIndexer

	// Get initial last success time
	initialTime, exists := hm.GetLastSuccessTime(component)
	require.True(t, exists)

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Record a new success
	hm.RecordLastSuccess(ctx, component)

	// Verify last success time was updated
	newTime, exists := hm.GetLastSuccessTime(component)
	require.True(t, exists)
	assert.True(t, newTime.After(initialTime), "Last success time should be updated")
}

func TestRecordHealthFromReport(t *testing.T) {
	t.Parallel()

	hm, err := NewHealthMetrics(testChainInfo())
	require.NoError(t, err)

	ctx := context.Background()

	// Simulate a health report with mixed statuses
	report := map[string]error{
		ComponentSuiRelayer:    nil,            // healthy
		ComponentTxM:           nil,            // healthy
		ComponentEventsIndexer: assert.AnError, // unhealthy
		ComponentChainReader:   nil,            // healthy
	}

	hm.RecordHealthFromReport(ctx, report)

	// Verify statuses were recorded correctly
	status, exists := hm.GetHealthStatus(ComponentSuiRelayer)
	assert.True(t, exists)
	assert.True(t, status)

	status, exists = hm.GetHealthStatus(ComponentTxM)
	assert.True(t, exists)
	assert.True(t, status)

	status, exists = hm.GetHealthStatus(ComponentEventsIndexer)
	assert.True(t, exists)
	assert.False(t, status)

	status, exists = hm.GetHealthStatus(ComponentChainReader)
	assert.True(t, exists)
	assert.True(t, status)
}

func TestUpdateProcessingLag(t *testing.T) {
	t.Parallel()

	hm, err := NewHealthMetrics(testChainInfo())
	require.NoError(t, err)

	ctx := context.Background()

	// Record a success for a component
	hm.RecordLastSuccess(ctx, ComponentTxM)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Update processing lag - this should not panic and should complete
	hm.UpdateProcessingLag(ctx)

	// Verify the component still has a valid last success time
	lastSuccess, exists := hm.GetLastSuccessTime(ComponentTxM)
	assert.True(t, exists)
	assert.False(t, lastSuccess.IsZero())
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	hm, err := NewHealthMetrics(testChainInfo())
	require.NoError(t, err)

	ctx := context.Background()
	done := make(chan struct{})

	// Start multiple goroutines that access health metrics concurrently
	for range 10 {
		go func() {
			for j := range 100 {
				hm.RecordHealth(ctx, ComponentTxM, j%2 == 0)
				hm.RecordLastSuccess(ctx, ComponentEventsIndexer)
				hm.UpdateProcessingLag(ctx)
				_, _ = hm.GetHealthStatus(ComponentTxM)
				_, _ = hm.GetLastSuccessTime(ComponentEventsIndexer)
			}
			done <- struct{}{}
		}()
	}

	// Wait for all goroutines to complete
	for range 10 {
		<-done
	}

	// If we got here without panics or deadlocks, the test passes
}

func TestFlipFlopDetection(t *testing.T) {
	t.Parallel()

	hm, err := NewHealthMetrics(testChainInfo())
	require.NoError(t, err)

	ctx := context.Background()
	component := ComponentSuiRelayer

	// Initial state is unhealthy (false)
	// Recording unhealthy again should NOT be a flip-flop
	hm.RecordHealth(ctx, component, false)
	status, _ := hm.GetHealthStatus(component)
	assert.False(t, status)

	// Recording healthy IS a flip-flop
	hm.RecordHealth(ctx, component, true)
	status, _ = hm.GetHealthStatus(component)
	assert.True(t, status)

	// Recording healthy again should NOT be a flip-flop
	hm.RecordHealth(ctx, component, true)
	status, _ = hm.GetHealthStatus(component)
	assert.True(t, status)

	// Recording unhealthy IS a flip-flop
	hm.RecordHealth(ctx, component, false)
	status, _ = hm.GetHealthStatus(component)
	assert.False(t, status)
}
