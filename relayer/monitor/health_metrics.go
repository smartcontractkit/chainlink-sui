package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	aptosTypes "github.com/smartcontractkit/chainlink-aptos/relayer/types"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"

	"github.com/smartcontractkit/chainlink-aptos/relayer/monitoring/metric/utils"
)

// Component names for health metrics
const (
	ComponentSuiRelayer          = "SuiRelayer"
	ComponentTxM                 = "TxM"
	ComponentChainPoller         = "ChainPoller"
	ComponentEventsIndexer       = "EventsIndexer"
	ComponentTransactionsIndexer = "TransactionsIndexer"
	ComponentChainReader         = "ChainReader"
	ComponentChainWriter         = "ChainWriter"
)

// HealthMetrics provides per-component health metrics that are pushed to Grafana.
// It tracks component status, flip-flop counts, and processing lag for monitoring
// and alerting purposes.
type HealthMetrics struct {
	chainInfo aptosTypes.ChainInfo

	// Metrics
	statusGauge               metric.Int64Gauge
	flipFlopCounter           metric.Int64Counter
	lastSuccessTimestampGauge metric.Int64Gauge
	processingLagGauge        metric.Float64Gauge

	// Internal state for tracking flip-flops
	mu               sync.RWMutex
	lastHealthStatus map[string]bool
	lastSuccessTime  map[string]time.Time
}

// NewHealthMetrics creates a new HealthMetrics instance and pre-registers all metrics.
// Metrics are registered at startup so they are always present for alerting.
func NewHealthMetrics(chainInfo aptosTypes.ChainInfo) (*HealthMetrics, error) {
	meter := beholder.GetMeter()

	statusGauge, err := meter.Int64Gauge(
		"sui_component_status",
		metric.WithDescription("Component health status: 1=healthy, 0=unhealthy"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create status gauge: %w", err)
	}

	flipFlopCounter, err := meter.Int64Counter(
		"sui_component_flip_flops_total",
		metric.WithDescription("Number of healthy/unhealthy state transitions"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create flip-flop counter: %w", err)
	}

	lastSuccessTimestampGauge, err := meter.Int64Gauge(
		"sui_component_last_success_timestamp",
		metric.WithDescription("Unix timestamp of last successful operation"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create last success timestamp gauge: %w", err)
	}

	processingLagGauge, err := meter.Float64Gauge(
		"sui_component_processing_lag_seconds",
		metric.WithDescription("Time since last successful operation in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create processing lag gauge: %w", err)
	}

	hm := &HealthMetrics{
		chainInfo:                 chainInfo,
		statusGauge:               statusGauge,
		flipFlopCounter:           flipFlopCounter,
		lastSuccessTimestampGauge: lastSuccessTimestampGauge,
		processingLagGauge:        processingLagGauge,
		lastHealthStatus:          make(map[string]bool),
		lastSuccessTime:           make(map[string]time.Time),
	}

	// Pre-register all component metrics with initial values
	components := []string{
		ComponentSuiRelayer,
		ComponentTxM,
		ComponentChainPoller,
		ComponentEventsIndexer,
		ComponentTransactionsIndexer,
		ComponentChainReader,
		ComponentChainWriter,
	}

	ctx := context.Background()
	now := time.Now()
	for _, component := range components {
		// Initialize with unknown/unhealthy state
		hm.lastHealthStatus[component] = false
		hm.lastSuccessTime[component] = now

		// Record initial metrics
		attrs := hm.getAttributes(component)
		hm.statusGauge.Record(ctx, 0, metric.WithAttributeSet(attrs))
		hm.lastSuccessTimestampGauge.Record(ctx, now.Unix(), metric.WithAttributeSet(attrs))
		hm.processingLagGauge.Record(ctx, 0, metric.WithAttributeSet(attrs))
	}

	return hm, nil
}

// RecordHealth records the health status of a component and tracks flip-flops.
func (hm *HealthMetrics) RecordHealth(ctx context.Context, component string, healthy bool) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	attrs := hm.getAttributes(component)

	// Check for flip-flop (state change)
	if lastStatus, exists := hm.lastHealthStatus[component]; exists && lastStatus != healthy {
		hm.flipFlopCounter.Add(ctx, 1, metric.WithAttributeSet(attrs))
	}
	hm.lastHealthStatus[component] = healthy

	// Record current status
	var statusValue int64
	if healthy {
		statusValue = 1
	}
	hm.statusGauge.Record(ctx, statusValue, metric.WithAttributeSet(attrs))
}

// RecordLastSuccess records the timestamp of a successful operation and updates
// the processing lag metric.
func (hm *HealthMetrics) RecordLastSuccess(ctx context.Context, component string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	now := time.Now()
	hm.lastSuccessTime[component] = now

	attrs := hm.getAttributes(component)
	hm.lastSuccessTimestampGauge.Record(ctx, now.Unix(), metric.WithAttributeSet(attrs))
	hm.processingLagGauge.Record(ctx, 0, metric.WithAttributeSet(attrs))
}

// UpdateProcessingLag updates the processing lag for all components based on their
// last success time. This should be called periodically.
func (hm *HealthMetrics) UpdateProcessingLag(ctx context.Context) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	now := time.Now()
	for component, lastSuccess := range hm.lastSuccessTime {
		lag := now.Sub(lastSuccess).Seconds()
		attrs := hm.getAttributes(component)
		hm.processingLagGauge.Record(ctx, lag, metric.WithAttributeSet(attrs))
	}
}

// RecordHealthFromReport processes a HealthReport map and records metrics for all components.
func (hm *HealthMetrics) RecordHealthFromReport(ctx context.Context, report map[string]error) {
	for component, err := range report {
		healthy := err == nil
		hm.RecordHealth(ctx, component, healthy)
	}
}

// getAttributes returns the OpenTelemetry attributes for a component.
func (hm *HealthMetrics) getAttributes(component string) attribute.Set {
	return attribute.NewSet(
		attribute.String("component", component),
		attribute.String("chain_family_name", utils.ValOrUnknown(hm.chainInfo.ChainFamilyName)),
		attribute.String("chain_id", utils.ValOrUnknown(hm.chainInfo.ChainID)),
		attribute.String("network_name", utils.ValOrUnknown(hm.chainInfo.NetworkName)),
		attribute.String("network_name_full", utils.ValOrUnknown(hm.chainInfo.NetworkNameFull)),
	)
}

// GetLastSuccessTime returns the last success time for a component.
func (hm *HealthMetrics) GetLastSuccessTime(component string) (time.Time, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	t, ok := hm.lastSuccessTime[component]
	return t, ok
}

// GetHealthStatus returns the last known health status for a component.
func (hm *HealthMetrics) GetHealthStatus(component string) (bool, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	status, ok := hm.lastHealthStatus[component]
	return status, ok
}
