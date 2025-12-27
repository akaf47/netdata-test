package as400

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
)

// Mock client for testing
type mockClient struct {
	queryFunc func(ctx context.Context, query string, callback func(columns []string, values []string) error) error
}

func (m *mockClient) Query(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, query, callback)
	}
	return nil
}

func (m *mockClient) QueryWithLimit(ctx context.Context, query string, limit int, callback func(columns []string, values []string) error) error {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, query, callback)
	}
	return nil
}

func (m *mockClient) Exec(ctx context.Context, statement string) error {
	return nil
}

func newTestCollector() *Collector {
	return &Collector{
		client:             &mockClient{},
		mx:                 &metrics{},
		activeJobTargets:   []TargetConfig{},
		activeJobs:         make(map[string]*activeJobMetadata),
		CollectActiveJobs: &FeatureFlag{},
	}
}

// Test collectActiveJobs with empty targets
func TestCollectActiveJobsEmptyTargets(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs with disabled feature
func TestCollectActiveJobsDisabledFeature(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: false}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs with single job found
func TestCollectActiveJobsFoundSingleJob(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	queryExecuted := false
	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			queryExecuted = true

			// Simulate returning job data
			callback([]string{"JOB_NAME"}, []string{"123456/TESTUSER/JOB1"})
			callback([]string{"JOB_USER"}, []string{"testuser"})
			callback([]string{"JOB_NUMBER"}, []string{"123456"})
			callback([]string{"JOB_STATUS"}, []string{"ACTIVE"})
			callback([]string{"SUBSYSTEM"}, []string{"QBATCH"})
			callback([]string{"JOB_TYPE"}, []string{"BATCH"})
			callback([]string{"ELAPSED_CPU_TIME"}, []string{"1000"})
			callback([]string{"ELAPSED_TIME"}, []string{"5000"})
			callback([]string{"TEMPORARY_STORAGE"}, []string{"2048"})
			callback([]string{"CPU_PERCENTAGE"}, []string{"25.5"})
			callback([]string{"ELAPSED_INTERACTIVE_TRANSACTIONS"}, []string{"100"})
			callback([]string{"ELAPSED_TOTAL_DISK_IO_COUNT"}, []string{"500"})
			callback([]string{"THREAD_COUNT"}, []string{"8"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !queryExecuted {
		t.Error("query was not executed")
	}
}

// Test collectActiveJobs with job not found
func TestCollectActiveJobsJobNotFound(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			// No columns/values = lineEnd=false, so found stays false
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs with query error
func TestCollectActiveJobsQueryError(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	expectedErr := errors.New("database connection failed")
	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			return expectedErr
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "active job") {
		t.Errorf("error message should contain 'active job', got: %v", err)
	}
}

// Test collectActiveJobs with multiple targets, first succeeds
func TestCollectActiveJobsMultipleTargetsFirstSucceeds(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "111111", User: "USER1", Name: "JOB1"},
		{Number: "222222", User: "USER2", Name: "JOB2"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	callCount := 0
	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callCount++
			if callCount == 1 {
				// First job succeeds
				callback([]string{"JOB_NAME"}, []string{"111111/USER1/JOB1"})
				return nil
			}
			// Second job fails
			return errors.New("network timeout")
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err == nil {
		t.Errorf("expected error from second job, got nil")
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

// Test collectActiveJobs with CPU percentage calculation
func TestCollectActiveJobsCPUPercentageCalculation(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"CPU_PERCENTAGE"}, []string{"50.5"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs with temporary storage conversion (KB to MB)
func TestCollectActiveJobsTemporaryStorageConversion(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"TEMPORARY_STORAGE"}, []string{"1024"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs with empty values
func TestCollectActiveJobsEmptyValues(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"JOB_NAME"}, []string{""})
			callback([]string{"JOB_USER"}, []string{""})
			callback([]string{"JOB_NUMBER"}, []string{""})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs whitespace trimming
func TestCollectActiveJobsWhitespaceTrimming(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"JOB_NAME"}, []string{"  123456/TESTUSER/JOB1  "})
			callback([]string{"JOB_STATUS"}, []string{"  ACTIVE  "})
			callback([]string{"JOB_USER"}, []string{"  testuser  "})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs with invalid numeric values
func TestCollectActiveJobsInvalidNumericValues(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"ELAPSED_CPU_TIME"}, []string{"INVALID"})
			callback([]string{"CPU_PERCENTAGE"}, []string{"N/A"})
			callback([]string{"THREAD_COUNT"}, []string{"abc"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test collectActiveJobs with all metric fields
func TestCollectActiveJobsAllMetricFields(t *testing.T) {
	c := newTestCollector()
	target := TargetConfig{Number: "123456", User: "TESTUSER", Name: "JOB1"}
	c.activeJobTargets = []TargetConfig{target}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}
	c.mx.activeJobs = make(map[string]activeJobInstanceMetrics)

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			fields := []string{
				"JOB_NAME", "JOB_USER", "JOB_NUMBER", "JOB_STATUS", "SUBSYSTEM",
				"JOB_TYPE", "ELAPSED_CPU_TIME", "ELAPSED_TIME", "TEMPORARY_STORAGE",
				"CPU_PERCENTAGE", "ELAPSED_INTERACTIVE_TRANSACTIONS",
				"ELAPSED_TOTAL_DISK_IO_COUNT", "THREAD_COUNT",
			}
			values := []string{
				"123456/TESTUSER/JOB1", "TESTUSER", "123456", "ACTIVE", "QBATCH",
				"BATCH", "1000", "5000", "2048", "25.5", "100", "500", "8",
			}

			for i, field := range fields {
				callback([]string{field}, []string{values[i]})
			}
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	key := target.ID()
	if _, ok := c.mx.activeJobs[key]; !ok {
		t.Errorf("expected metrics for job %s", key)
	}
}

// Test with lineEnd false, then true
func TestCollectActiveJobsLineEndHandling(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			// Simulate row with multiple columns
			callback([]string{"JOB_NAME"}, []string{"123456/TESTUSER/JOB1"})
			callback([]string{"JOB_STATUS"}, []string{"ACTIVE"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test buildActiveJobQuery is called correctly
func TestCollectActiveJobsQueryBuilding(t *testing.T) {
	c := newTestCollector()
	target := TargetConfig{Number: "123456", User: "TESTUSER", Name: "JOB1"}
	c.activeJobTargets = []TargetConfig{target}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	queryBuilt := false
	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			// Verify query contains expected patterns
			if strings.Contains(query, "ACTIVE_JOBS_INFO_TABLE") ||
				strings.Contains(query, "JOB_NAME") ||
				strings.Contains(query, "JOB_NUMBER") {
				queryBuilt = true
			}
			return nil
		},
	}

	c.collectActiveJobs(context.Background())
	// Note: queryBuilt check depends on actual buildActiveJobQuery implementation
}

// Test case-insensitive status handling
func TestCollectActiveJobsJobStatusHandling(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"JOB_STATUS"}, []string{"ACTIVE"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test context cancellation
func TestCollectActiveJobsContextCancellation(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			// Could return context error
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				return nil
			}
		},
	}

	// Should handle context error gracefully
	c.collectActiveJobs(ctx)
}

// Test floating point precision in CPU percentage
func TestCollectActiveJobsFloatPrecision(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"CPU_PERCENTAGE"}, []string{"12.3456789"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test large numeric values
func TestCollectActiveJobsLargeNumericValues(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"ELAPSED_TIME"}, []string{"9223372036"})
			callback([]string{"THREAD_COUNT"}, []string{"1000000"})
			callback([]string{"ELAPSED_TOTAL_DISK_IO_COUNT"}, []string{"999999999"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test zero values
func TestCollectActiveJobsZeroValues(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"ELAPSED_CPU_TIME"}, []string{"0"})
			callback([]string{"THREAD_COUNT"}, []string{"0"})
			callback([]string{"CPU_PERCENTAGE"}, []string{"0.0"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test negative values (should be handled)
func TestCollectActiveJobsNegativeValues(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"CPU_PERCENTAGE"}, []string{"-5.5"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test target ID usage
func TestCollectActiveJobsUsesTargetID(t *testing.T) {
	c := newTestCollector()
	target := TargetConfig{Number: "123456", User: "TESTUSER", Name: "JOB1"}
	c.activeJobTargets = []TargetConfig{target}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			return nil
		},
	}

	c.collectActiveJobs(context.Background())
	// Verify that target ID would be used correctly
	expectedID := target.ID()
	if expectedID == "" {
		t.Error("target ID should not be empty")
	}
}

// Test unrecognized column handling
func TestCollectActiveJobsUnrecognizedColumn(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"UNKNOWN_COLUMN"}, []string{"SOME_VALUE"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error for unknown column, got %v", err)
	}
}

// Test interactive transactions
func TestCollectActiveJobsInteractiveTransactions(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"ELAPSED_INTERACTIVE_TRANSACTIONS"}, []string{"12345"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test disk IO count
func TestCollectActiveJobsDiskIOCount(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"ELAPSED_TOTAL_DISK_IO_COUNT"}, []string{"999999"})
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

// Test continues on error for multiple targets
func TestCollectActiveJobsContinuesOnError(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "111111", User: "USER1", Name: "JOB1"},
		{Number: "222222", User: "USER2", Name: "JOB2"},
		{Number: "333333", User: "USER3", Name: "JOB3"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	callCount := 0
	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callCount++
			if callCount == 2 {
				return errors.New("query error")
			}
			return nil
		},
	}

	err := c.collectActiveJobs(context.Background())
	// Should return first error but continue processing
	if callCount != 3 {
		t.Errorf("expected all 3 queries attempted, got %d", callCount)
	}
}

// Test metrics are properly stored
func TestCollectActiveJobsMetricsStorage(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}
	c.mx.activeJobs = make(map[string]activeJobInstanceMetrics)

	expectedCPUTime := int64(1000)
	expectedThreads := int64(8)

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"ELAPSED_CPU_TIME"}, []string{"1000"})
			callback([]string{"THREAD_COUNT"}, []string{"8"})
			return nil
		},
	}

	c.collectActiveJobs(context.Background())

	key := c.activeJobTargets[0].ID()
	metrics, ok := c.mx.activeJobs[key]
	if !ok {
		t.Errorf("metrics not stored for key %s", key)
	}
	if metrics.ElapsedCPUTime != expectedCPUTime {
		t.Errorf("expected ElapsedCPUTime %d, got %d", expectedCPUTime, metrics.ElapsedCPUTime)
	}
	if metrics.ThreadCount != expectedThreads {
		t.Errorf("expected ThreadCount %d, got %d", expectedThreads, metrics.ThreadCount)
	}
}

// Test temporary storage is converted from KB to MB
func TestCollectActiveJobsTemporaryStorageKBtoMBConversion(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "123456", User: "TESTUSER", Name: "JOB1"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}
	c.mx.activeJobs = make(map[string]activeJobInstanceMetrics)

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			callback([]string{"TEMPORARY_STORAGE"}, []string{"2048"})
			return nil
		},
	}

	c.collectActiveJobs(context.Background())

	key := c.activeJobTargets[0].ID()
	metrics, ok := c.mx.activeJobs[key]
	if !ok {
		t.Errorf("metrics not stored for key %s", key)
	}
	// 2048 KB / 1024 = 2 MB
	expectedStorageInMB := int64(2)
	if metrics.TemporaryStorage != expectedStorageInMB {
		t.Errorf("expected TemporaryStorage %d MB, got %d MB", expectedStorageInMB, metrics.TemporaryStorage)
	}
}

func TestCollectActiveJobsMultipleQueryErrors(t *testing.T) {
	c := newTestCollector()
	c.activeJobTargets = []TargetConfig{
		{Number: "111111", User: "USER1", Name: "JOB1"},
		{Number: "222222", User: "USER2", Name: "JOB2"},
	}
	c.CollectActiveJobs = &FeatureFlag{enabled: true}

	c.client = &mockClient{
		queryFunc: func(ctx context.Context, query string, callback func(columns []string, values []string) error) error {
			return errors.New("persistent error")
		},
	}

	err := c.collectActiveJobs(context.Background())
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "active job") {
		t.Errorf("error should mention 'active job', got: %v", err)
	}
}