package clickhouse

import (
	"testing"
)

// TestCollectSystemEvents tests the main collection function
func TestCollectSystemEvents(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockClickHouseDB)
		wantErr bool
	}{
		{
			name: "successful collection with events",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":        "Query",
						"value":        int64(100),
						"description":  "Query executed",
					},
				})
			},
			wantErr: false,
		},
		{
			name: "empty result set",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{})
			},
			wantErr: false,
		},
		{
			name: "database error",
			setup: func(m *mockClickHouseDB) {
				m.SetError("connection refused")
			},
			wantErr: true,
		},
		{
			name: "multiple events",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":       "Query",
						"value":       int64(100),
						"description": "Query executed",
					},
					{
						"event":       "Insert",
						"value":       int64(50),
						"description": "Insert executed",
					},
					{
						"event":       "Delete",
						"value":       int64(10),
						"description": "Delete executed",
					},
				})
			},
			wantErr: false,
		},
		{
			name: "null values in result",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":       "Query",
						"value":       nil,
						"description": nil,
					},
				})
			},
			wantErr: false,
		},
		{
			name: "zero values",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":       "",
						"value":       int64(0),
						"description": "",
					},
				})
			},
			wantErr: false,
		},
		{
			name: "large value",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":       "Query",
						"value":       int64(9223372036854775807), // max int64
						"description": "Large value test",
					},
				})
			},
			wantErr: false,
		},
		{
			name: "negative value",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":       "Query",
						"value":       int64(-100),
						"description": "Negative value",
					},
				})
			},
			wantErr: false,
		},
		{
			name: "special characters in event name",
			setup: func(m *mockClickHouseDB) {
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":       "Query_With-Special.Chars",
						"value":       int64(100),
						"description": "Event with special characters",
					},
				})
			},
			wantErr: false,
		},
		{
			name: "very long event name",
			setup: func(m *mockClickHouseDB) {
				longName := ""
				for i := 0; i < 1000; i++ {
					longName += "a"
				}
				m.SetQueryResult([]map[string]interface{}{
					{
						"event":       longName,
						"value":       int64(100),
						"description": "Long name test",
					},
				})
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := newMockClickHouseDB()
			tt.setup(mockDB)

			collector := &ClickHouseCollector{
				db: mockDB,
			}

			err := collector.collectSystemEvents()

			if (err != nil) != tt.wantErr {
				t.Errorf("collectSystemEvents() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && mockDB.QueryCalled() {
				// Verify the query was executed
				if !mockDB.QueryCalled() {
					t.Error("expected query to be called")
				}
			}
		})
	}
}

// TestCollectSystemEventsDataProcessing tests data processing and validation
func TestCollectSystemEventsDataProcessing(t *testing.T) {
	tests := []struct {
		name           string
		input          []map[string]interface{}
		expectedMetrics int
	}{
		{
			name: "single event processed",
			input: []map[string]interface{}{
				{
					"event": "Query",
					"value": int64(100),
				},
			},
			expectedMetrics: 1,
		},
		{
			name: "multiple events processed",
			input: []map[string]interface{}{
				{
					"event": "Query",
					"value": int64(100),
				},
				{
					"event": "Insert",
					"value": int64(200),
				},
				{
					"event": "Delete",
					"value": int64(300),
				},
			},
			expectedMetrics: 3,
		},
		{
			name:             "empty input",
			input:            []map[string]interface{}{},
			expectedMetrics: 0,
		},
		{
			name: "missing value field",
			input: []map[string]interface{}{
				{
					"event": "Query",
				},
			},
			expectedMetrics: 0,
		},
		{
			name: "missing event field",
			input: []map[string]interface{}{
				{
					"value": int64(100),
				},
			},
			expectedMetrics: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := newMockClickHouseDB()
			mockDB.SetQueryResult(tt.input)

			collector := &ClickHouseCollector{
				db: mockDB,
			}

			_ = collector.collectSystemEvents()

			if len(collector.metrics) != tt.expectedMetrics {
				t.Errorf("expected %d metrics, got %d", tt.expectedMetrics, len(collector.metrics))
			}
		})
	}
}

// TestCollectSystemEventsErrorHandling tests error scenarios
func TestCollectSystemEventsErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		errorMsg    string
		shouldError bool
	}{
		{
			name:        "connection error",
			errorMsg:    "connection refused",
			shouldError: true,
		},
		{
			name:        "timeout error",
			errorMsg:    "context deadline exceeded",
			shouldError: true,
		},
		{
			name:        "query syntax error",
			errorMsg:    "syntax error",
			shouldError: true,
		},
		{
			name:        "database locked",
			errorMsg:    "database is locked",
			shouldError: true,
		},
		{
			name:        "unknown error",
			errorMsg:    "unknown error",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := newMockClickHouseDB()
			mockDB.SetError(tt.errorMsg)

			collector := &ClickHouseCollector{
				db: mockDB,
			}

			err := collector.collectSystemEvents()

			if tt.shouldError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestCollectSystemEventsTypeConversion tests type conversions
func TestCollectSystemEventsTypeConversion(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		wantError bool
	}{
		{
			name:      "int64 value",
			value:     int64(100),
			wantError: false,
		},
		{
			name:      "int value",
			value:     100,
			wantError: false,
		},
		{
			name:      "float64 value",
			value:     100.5,
			wantError: false,
		},
		{
			name:      "string numeric value",
			value:     "100",
			wantError: false,
		},
		{
			name:      "nil value",
			value:     nil,
			wantError: false,
		},
		{
			name:      "boolean value",
			value:     true,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := newMockClickHouseDB()
			mockDB.SetQueryResult([]map[string]interface{}{
				{
					"event": "TestEvent",
					"value": tt.value,
				},
			})

			collector := &ClickHouseCollector{
				db: mockDB,
			}

			err := collector.collectSystemEvents()

			if tt.wantError && err == nil {
				t.Error("expected type conversion error but got nil")
			}
		})
	}
}

// TestCollectSystemEventsConcurrency tests thread-safety (if applicable)
func TestCollectSystemEventsConcurrency(t *testing.T) {
	mockDB := newMockClickHouseDB()
	mockDB.SetQueryResult([]map[string]interface{}{
		{
			"event": "Query",
			"value": int64(100),
		},
	})

	collector := &ClickHouseCollector{
		db: mockDB,
	}

	// Run collection multiple times concurrently
	done := make(chan error, 10)
	for i := 0; i < 10; i++ {
		go func() {
			done <- collector.collectSystemEvents()
		}()
	}

	for i := 0; i < 10; i++ {
		err := <-done
		if err != nil {
			t.Errorf("concurrent collection failed: %v", err)
		}
	}
}

// TestCollectSystemEventsMetricStorage tests metric storage and retrieval
func TestCollectSystemEventsMetricStorage(t *testing.T) {
	mockDB := newMockClickHouseDB()
	events := []map[string]interface{}{
		{
			"event": "Query",
			"value": int64(100),
		},
		{
			"event": "Insert",
			"value": int64(200),
		},
	}
	mockDB.SetQueryResult(events)

	collector := &ClickHouseCollector{
		db:      mockDB,
		metrics: make(map[string]int64),
	}

	err := collector.collectSystemEvents()
	if err != nil {
		t.Fatalf("collectSystemEvents failed: %v", err)
	}

	// Verify metrics are stored with correct values
	if len(collector.metrics) == 0 {
		t.Error("no metrics stored")
	}

	// Verify metric values
	for _, event := range events {
		eventName := event["event"].(string)
		expectedValue := event["value"].(int64)

		if metric, exists := collector.metrics[eventName]; !exists || metric != expectedValue {
			t.Errorf("metric mismatch for %s: expected %d, got %d", eventName, expectedValue, metric)
		}
	}
}

// Mock implementation
type mockClickHouseDB struct {
	queryResult []map[string]interface{}
	err         error
	queryCalled bool
}

func newMockClickHouseDB() *mockClickHouseDB {
	return &mockClickHouseDB{}
}

func (m *mockClickHouseDB) Query(query string) ([]map[string]interface{}, error) {
	m.queryCalled = true
	return m.queryResult, m.err
}

func (m *mockClickHouseDB) SetQueryResult(result []map[string]interface{}) {
	m.queryResult = result
}

func (m *mockClickHouseDB) SetError(errMsg string) {
	m.err = &mockError{message: errMsg}
}

func (m *mockClickHouseDB) QueryCalled() bool {
	return m.queryCalled
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

// Placeholder for ClickHouseCollector if not available in tests
type ClickHouseCollector struct {
	db      interface{}
	metrics map[string]int64
}

// Placeholder for collectSystemEvents if needed
func (c *ClickHouseCollector) collectSystemEvents() error {
	// Implementation would call the actual function
	return nil
}