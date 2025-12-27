package tomcat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStatusResponse tests the StatusResponse struct and its methods
func TestStatusResponse_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expected  *StatusResponse
		wantError bool
	}{
		{
			name: "valid response with all fields",
			data: []byte(`{
				"jvm": {"memory": {"max": 1000, "committed": 900, "used": 800}},
				"connector": [{"name": "http-nio-8080", "processing": 5, "requestCount": 100, "bytesReceived": 1000, "bytesSent": 2000}],
				"threadInfo": {"currentThreadCount": 10, "peakThreadCount": 20},
				"requestInfo": {"maxTime": 5000, "processingTime": 10000, "requestCount": 50, "errorCount": 2, "bytesReceived": 5000, "bytesSent": 10000}
			}`),
			expected: &StatusResponse{
				JVM: &JVM{
					Memory: &Memory{
						Max:       1000,
						Committed: 900,
						Used:      800,
					},
				},
				Connector: []Connector{
					{
						Name:         "http-nio-8080",
						Processing:   5,
						RequestCount: 100,
						BytesReceived: 1000,
						BytesSent:    2000,
					},
				},
				ThreadInfo: &ThreadInfo{
					CurrentThreadCount: 10,
					PeakThreadCount:    20,
				},
				RequestInfo: &RequestInfo{
					MaxTime:        5000,
					ProcessingTime: 10000,
					RequestCount:   50,
					ErrorCount:     2,
					BytesReceived:  5000,
					BytesSent:      10000,
				},
			},
			wantError: false,
		},
		{
			name:      "empty response",
			data:      []byte(`{}`),
			expected:  &StatusResponse{},
			wantError: false,
		},
		{
			name:      "null jvm",
			data:      []byte(`{"jvm": null}`),
			expected:  &StatusResponse{JVM: nil},
			wantError: false,
		},
		{
			name:      "null connector",
			data:      []byte(`{"connector": null}`),
			expected:  &StatusResponse{Connector: nil},
			wantError: false,
		},
		{
			name:      "empty connector array",
			data:      []byte(`{"connector": []}`),
			expected:  &StatusResponse{Connector: []Connector{}},
			wantError: false,
		},
		{
			name:      "null threadInfo",
			data:      []byte(`{"threadInfo": null}`),
			expected:  &StatusResponse{ThreadInfo: nil},
			wantError: false,
		},
		{
			name:      "null requestInfo",
			data:      []byte(`{"requestInfo": null}`),
			expected:  &StatusResponse{RequestInfo: nil},
			wantError: false,
		},
		{
			name:      "invalid json",
			data:      []byte(`{invalid}`),
			wantError: true,
		},
		{
			name:      "empty bytes",
			data:      []byte(``),
			wantError: true,
		},
		{
			name: "connector with all fields",
			data: []byte(`{
				"connector": [
					{"name": "http-nio-8080", "processing": 5, "requestCount": 100, "bytesReceived": 1000, "bytesSent": 2000},
					{"name": "http-nio-8443", "processing": 10, "requestCount": 200, "bytesReceived": 2000, "bytesSent": 4000}
				]
			}`),
			expected: &StatusResponse{
				Connector: []Connector{
					{
						Name:         "http-nio-8080",
						Processing:   5,
						RequestCount: 100,
						BytesReceived: 1000,
						BytesSent:    2000,
					},
					{
						Name:         "http-nio-8443",
						Processing:   10,
						RequestCount: 200,
						BytesReceived: 2000,
						BytesSent:    4000,
					},
				},
			},
			wantError: false,
		},
		{
			name: "jvm with partial memory info",
			data: []byte(`{
				"jvm": {"memory": {"max": 2000}}
			}`),
			expected: &StatusResponse{
				JVM: &JVM{
					Memory: &Memory{
						Max: 2000,
					},
				},
			},
			wantError: false,
		},
		{
			name: "threadInfo with zero values",
			data: []byte(`{
				"threadInfo": {"currentThreadCount": 0, "peakThreadCount": 0}
			}`),
			expected: &StatusResponse{
				ThreadInfo: &ThreadInfo{
					CurrentThreadCount: 0,
					PeakThreadCount:    0,
				},
			},
			wantError: false,
		},
		{
			name: "requestInfo with zero values",
			data: []byte(`{
				"requestInfo": {"maxTime": 0, "processingTime": 0, "requestCount": 0, "errorCount": 0, "bytesReceived": 0, "bytesSent": 0}
			}`),
			expected: &StatusResponse{
				RequestInfo: &RequestInfo{
					MaxTime:        0,
					ProcessingTime: 0,
					RequestCount:   0,
					ErrorCount:     0,
					BytesReceived:  0,
					BytesSent:      0,
				},
			},
			wantError: false,
		},
		{
			name: "large numeric values",
			data: []byte(`{
				"jvm": {"memory": {"max": 9223372036854775807, "committed": 9223372036854775806, "used": 9223372036854775805}},
				"requestInfo": {"maxTime": 999999999999, "processingTime": 999999999999}
			}`),
			expected: &StatusResponse{
				JVM: &JVM{
					Memory: &Memory{
						Max:       9223372036854775807,
						Committed: 9223372036854775806,
						Used:      9223372036854775805,
					},
				},
				RequestInfo: &RequestInfo{
					MaxTime:        999999999999,
					ProcessingTime: 999999999999,
				},
			},
			wantError: false,
		},
		{
			name: "negative numeric values",
			data: []byte(`{
				"jvm": {"memory": {"max": -1000, "committed": -500, "used": -100}},
				"threadInfo": {"currentThreadCount": -1, "peakThreadCount": -1}
			}`),
			expected: &StatusResponse{
				JVM: &JVM{
					Memory: &Memory{
						Max:       -1000,
						Committed: -500,
						Used:      -100,
					},
				},
				ThreadInfo: &ThreadInfo{
					CurrentThreadCount: -1,
					PeakThreadCount:    -1,
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sr StatusResponse
			err := json.Unmarshal(tt.data, &sr)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, &sr)
		})
	}
}

// TestJVM_UnmarshalJSON tests JVM struct unmarshaling
func TestJVM_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expected  *JVM
		wantError bool
	}{
		{
			name: "valid jvm with memory",
			data: []byte(`{"memory": {"max": 1000, "committed": 900, "used": 800}}`),
			expected: &JVM{
				Memory: &Memory{
					Max:       1000,
					Committed: 900,
					Used:      800,
				},
			},
			wantError: false,
		},
		{
			name:      "jvm with null memory",
			data:      []byte(`{"memory": null}`),
			expected:  &JVM{Memory: nil},
			wantError: false,
		},
		{
			name:      "empty jvm",
			data:      []byte(`{}`),
			expected:  &JVM{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jvm JVM
			err := json.Unmarshal(tt.data, &jvm)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, &jvm)
		})
	}
}

// TestMemory_UnmarshalJSON tests Memory struct unmarshaling
func TestMemory_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expected  *Memory
		wantError bool
	}{
		{
			name: "valid memory",
			data: []byte(`{"max": 2000, "committed": 1800, "used": 1500}`),
			expected: &Memory{
				Max:       2000,
				Committed: 1800,
				Used:      1500,
			},
			wantError: false,
		},
		{
			name:      "memory with zero values",
			data:      []byte(`{"max": 0, "committed": 0, "used": 0}`),
			expected:  &Memory{Max: 0, Committed: 0, Used: 0},
			wantError: false,
		},
		{
			name:      "empty memory",
			data:      []byte(`{}`),
			expected:  &Memory{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mem Memory
			err := json.Unmarshal(tt.data, &mem)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, &mem)
		})
	}
}

// TestConnector_UnmarshalJSON tests Connector struct unmarshaling
func TestConnector_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expected  *Connector
		wantError bool
	}{
		{
			name: "valid connector",
			data: []byte(`{"name": "http-nio-8080", "processing": 5, "requestCount": 100, "bytesReceived": 1000, "bytesSent": 2000}`),
			expected: &Connector{
				Name:         "http-nio-8080",
				Processing:   5,
				RequestCount: 100,
				BytesReceived: 1000,
				BytesSent:    2000,
			},
			wantError: false,
		},
		{
			name: "connector with empty name",
			data: []byte(`{"name": "", "processing": 0, "requestCount": 0, "bytesReceived": 0, "bytesSent": 0}`),
			expected: &Connector{
				Name:         "",
				Processing:   0,
				RequestCount: 0,
				BytesReceived: 0,
				BytesSent:    0,
			},
			wantError: false,
		},
		{
			name:      "empty connector",
			data:      []byte(`{}`),
			expected:  &Connector{},
			wantError: false,
		},
		{
			name: "connector with negative values",
			data: []byte(`{"name": "http", "processing": -1, "requestCount": -100}`),
			expected: &Connector{
				Name:         "http",
				Processing:   -1,
				RequestCount: -100,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conn Connector
			err := json.Unmarshal(tt.data, &conn)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, &conn)
		})
	}
}

// TestThreadInfo_UnmarshalJSON tests ThreadInfo struct unmarshaling
func TestThreadInfo_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expected  *ThreadInfo
		wantError bool
	}{
		{
			name: "valid threadInfo",
			data: []byte(`{"currentThreadCount": 10, "peakThreadCount": 20}`),
			expected: &ThreadInfo{
				CurrentThreadCount: 10,
				PeakThreadCount:    20,
			},
			wantError: false,
		},
		{
			name: "threadInfo with zero values",
			data: []byte(`{"currentThreadCount": 0, "peakThreadCount": 0}`),
			expected: &ThreadInfo{
				CurrentThreadCount: 0,
				PeakThreadCount:    0,
			},
			wantError: false,
		},
		{
			name:      "empty threadInfo",
			data:      []byte(`{}`),
			expected:  &ThreadInfo{},
			wantError: false,
		},
		{
			name: "threadInfo with large values",
			data: []byte(`{"currentThreadCount": 999999, "peakThreadCount": 1000000}`),
			expected: &ThreadInfo{
				CurrentThreadCount: 999999,
				PeakThreadCount:    1000000,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ti ThreadInfo
			err := json.Unmarshal(tt.data, &ti)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, &ti)
		})
	}
}

// TestRequestInfo_UnmarshalJSON tests RequestInfo struct unmarshaling
func TestRequestInfo_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		expected  *RequestInfo
		wantError bool
	}{
		{
			name: "valid requestInfo",
			data: []byte(`{"maxTime": 5000, "processingTime": 10000, "requestCount": 50, "errorCount": 2, "bytesReceived": 5000, "bytesSent": 10000}`),
			expected: &RequestInfo{
				MaxTime:        5000,
				ProcessingTime: 10000,
				RequestCount:   50,
				ErrorCount:     2,
				BytesReceived:  5000,
				BytesSent:      10000,
			},
			wantError: false,
		},
		{
			name: "requestInfo with zero values",
			data: []byte(`{"maxTime": 0, "processingTime": 0, "requestCount": 0, "errorCount": 0, "bytesReceived": 0, "bytesSent": 0}`),
			expected: &RequestInfo{
				MaxTime:        0,
				ProcessingTime: 0,
				RequestCount:   0,
				ErrorCount:     0,
				BytesReceived:  0,
				BytesSent:      0,
			},
			wantError: false,
		},
		{
			name:      "empty requestInfo",
			data:      []byte(`{}`),
			expected:  &RequestInfo{},
			wantError: false,
		},
		{
			name: "requestInfo with large values",
			data: []byte(`{"maxTime": 9999999999, "processingTime": 9999999999, "requestCount": 9999999999, "errorCount": 999999, "bytesReceived": 9999999999, "bytesSent": 9999999999}`),
			expected: &RequestInfo{
				MaxTime:        9999999999,
				ProcessingTime: 9999999999,
				RequestCount:   9999999999,
				ErrorCount:     999999,
				BytesReceived:  9999999999,
				BytesSent:      9999999999,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ri RequestInfo
			err := json.Unmarshal(tt.data, &ri)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, &ri)
		})
	}
}

// TestStatusResponseIntegration tests complex scenarios with multiple nested structures
func TestStatusResponseIntegration(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		wantError bool
		validate  func(t *testing.T, sr *StatusResponse)
	}{
		{
			name: "complete response with all components",
			data: []byte(`{
				"jvm": {"memory": {"max": 2000, "committed": 1800, "used": 1500}},
				"connector": [
					{"name": "http-nio-8080", "processing": 5, "requestCount": 100, "bytesReceived": 1000, "bytesSent": 2000},
					{"name": "http-nio-8443", "processing": 10, "requestCount": 200, "bytesReceived": 2000, "bytesSent": 4000}
				],
				"threadInfo": {"currentThreadCount": 25, "peakThreadCount": 50},
				"requestInfo": {"maxTime": 5000, "processingTime": 20000, "requestCount": 300, "errorCount": 5, "bytesReceived": 3000, "bytesSent": 6000}
			}`),
			wantError: false,
			validate: func(t *testing.T, sr *StatusResponse) {
				assert.NotNil(t, sr.JVM)
				assert.NotNil(t, sr.JVM.Memory)
				assert.Equal(t, int64(2000), sr.JVM.Memory.Max)
				assert.Len(t, sr.Connector, 2)
				assert.Equal(t, "http-nio-8080", sr.Connector[0].Name)
				assert.Equal(t, "http-nio-8443", sr.Connector[1].Name)
				assert.NotNil(t, sr.ThreadInfo)
				assert.Equal(t, int64(25), sr.ThreadInfo.CurrentThreadCount)
				assert.NotNil(t, sr.RequestInfo)
				assert.Equal(t, int64(5), sr.RequestInfo.ErrorCount)
			},
		},
		{
			name:      "response with only jvm",
			data:      []byte(`{"jvm": {"memory": {"max": 1000, "committed": 900, "used": 800}}}`),
			wantError: false,
			validate: func(t *testing.T, sr *StatusResponse) {
				assert.NotNil(t, sr.JVM)
				assert.Nil(t, sr.ThreadInfo)
				assert.Nil(t, sr.RequestInfo)
				assert.Len(t, sr.Connector, 0)
			},
		},
		{
			name:      "response with multiple connectors",
			data:      []byte(`{"connector": [{"name": "c1"}, {"name": "c2"}, {"name": "c3"}]}`),
			wantError: false,
			validate: func(t *testing.T, sr *StatusResponse) {
				assert.Len(t, sr.Connector, 3)
				assert.Nil(t, sr.JVM)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sr StatusResponse
			err := json.Unmarshal(tt.data, &sr)

			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, &sr)
			}
		})
	}
}

// TestEdgeCases tests edge cases and boundary conditions
func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		wantError bool
	}{
		{
			name:      "very deep nesting",
			data:      []byte(`{"jvm": {"memory": {"max": 1}}}`),
			wantError: false,
		},
		{
			name:      "unicode in connector name",
			data:      []byte(`{"connector": [{"name": "http-тест-8080"}]}`),
			wantError: false,
		},
		{
			name:      "special characters in name",
			data:      []byte(`{"connector": [{"name": "http://special-chars_8080"}]}`),
			wantError: false,
		},
		{
			name:      "scientific notation",
			data:      []byte(`{"jvm": {"memory": {"max": 1e6, "committed": 1e5, "used": 1e4}}}`),
			wantError: false,
		},
		{
			name:      "decimal numbers",
			data:      []byte(`{"jvm": {"memory": {"max": 1000.5}}}`),
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sr StatusResponse
			err := json.Unmarshal(tt.data, &sr)

			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}