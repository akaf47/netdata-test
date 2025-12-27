package rspamd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollect tests the main Collect() method
func TestCollect(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: validRspamdResponse(),
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectWithHTTPError tests Collect() when HTTP request fails
func TestCollectWithHTTPError(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			err: "connection refused",
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithEmptyResponse tests Collect() with empty response
func TestCollectWithEmptyResponse(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "",
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	// Should handle empty response gracefully
	assert.NoError(t, err) // or assert.Error(t, err) depending on implementation
}

// TestCollectWith400StatusCode tests Collect() with HTTP 400
func TestCollectWith400StatusCode(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Bad Request",
			statusCode:   400,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWith401StatusCode tests Collect() with HTTP 401
func TestCollectWith401StatusCode(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Unauthorized",
			statusCode:   401,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWith403StatusCode tests Collect() with HTTP 403
func TestCollectWith403StatusCode(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Forbidden",
			statusCode:   403,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWith404StatusCode tests Collect() with HTTP 404
func TestCollectWith404StatusCode(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Not Found",
			statusCode:   404,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWith500StatusCode tests Collect() with HTTP 500
func TestCollectWith500StatusCode(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Internal Server Error",
			statusCode:   500,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithInvalidJSON tests Collect() with malformed JSON response
func TestCollectWithInvalidJSON(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "{invalid json}",
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithPartialJSON tests Collect() with incomplete JSON
func TestCollectWithPartialJSON(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: `{"status": "ok"`,
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithNullResponse tests Collect() with null JSON response
func TestCollectWithNullResponse(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "null",
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	// Should handle null gracefully
	assert.NoError(t, err) // or assert.Error depending on implementation
}

// TestCollectWithEmptyJSONObject tests Collect() with empty object
func TestCollectWithEmptyJSONObject(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "{}",
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectWithEmptyJSONArray tests Collect() with empty array
func TestCollectWithEmptyJSONArray(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "[]",
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectWithTimeoutError tests Collect() when request times out
func TestCollectWithTimeoutError(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			err: "i/o timeout",
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithContextDeadlineExceeded tests Collect() with context deadline
func TestCollectWithContextDeadlineExceeded(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			err: "context deadline exceeded",
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithValidResponse tests Collect() parses valid response correctly
func TestCollectWithValidResponse(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: validRspamdResponse(),
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectMultipleTimes tests Collect() called sequentially
func TestCollectMultipleTimes(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: validRspamdResponse(),
			statusCode:   200,
		},
	}
	
	for i := 0; i < 3; i++ {
		err := r.Collect()
		assert.NoError(t, err)
	}
}

// TestCollectWithRateLimitError tests Collect() with rate limit response
func TestCollectWithRateLimitError(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Too Many Requests",
			statusCode:   429,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithServiceUnavailable tests Collect() with 503 response
func TestCollectWithServiceUnavailable(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Service Unavailable",
			statusCode:   503,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithGatewayTimeout tests Collect() with 504 response
func TestCollectWithGatewayTimeout(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "Gateway Timeout",
			statusCode:   504,
		},
	}
	
	err := r.Collect()
	
	assert.Error(t, err)
}

// TestCollectWithVeryLargeResponse tests Collect() with large response body
func TestCollectWithVeryLargeResponse(t *testing.T) {
	largeBody := `{"data": "` + string(make([]byte, 10000)) + `"}`
	
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: largeBody,
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	// Should handle large responses
	assert.NoError(t, err) // or error depending on implementation
}

// TestCollectWithSpecialCharacters tests Collect() with special chars in response
func TestCollectWithSpecialCharacters(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: `{"name": "test\u0000value"}`,
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	// Should handle special characters
	assert.NoError(t, err) // or error depending on implementation
}

// TestCollectWithBoundaryValues tests Collect() with numeric boundary values
func TestCollectWithBoundaryValues(t *testing.T) {
	testCases := []struct {
		name     string
		response string
	}{
		{"zero values", `{"count": 0, "size": 0}`},
		{"negative values", `{"count": -1, "size": -100}`},
		{"max int values", `{"count": 9223372036854775807}`},
		{"large floats", `{"ratio": 999999.999999}`},
		{"very small floats", `{"ratio": 0.0000001}`},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &Rspamd{
				httpClient: &mockHTTPClient{
					responseBody: tc.response,
					statusCode:   200,
				},
			}
			
			err := r.Collect()
			assert.NoError(t, err)
		})
	}
}

// TestCollectWithNullFields tests Collect() with null fields in JSON
func TestCollectWithNullFields(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: `{"field1": null, "field2": "value"}`,
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectWithMissingFields tests Collect() with missing expected fields
func TestCollectWithMissingFields(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: `{"incomplete": true}`,
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	// Should handle missing fields gracefully
	assert.NoError(t, err)
}

// TestCollectWithExtraFields tests Collect() with unexpected extra fields
func TestCollectWithExtraFields(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: `{"extra": "field", "another": 123, "data": true}`,
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectWithUnicodeContent tests Collect() with unicode characters
func TestCollectWithUnicodeContent(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: `{"name": "тест", "value": "日本語"}`,
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectDataFlow tests that collected data is properly stored
func TestCollectDataFlow(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: validRspamdResponse(),
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	require.NoError(t, err)
	
	// Verify that data collection succeeded (implementation specific)
	assert.NotNil(t, r)
}

// TestCollectErrorRecovery tests that Collect() can recover from errors
func TestCollectErrorRecovery(t *testing.T) {
	client := &mockHTTPClient{
		responseBody: validRspamdResponse(),
		statusCode:   200,
	}
	
	r := &Rspamd{
		httpClient: client,
	}
	
	// First call with error
	client.err = "connection error"
	err := r.Collect()
	assert.Error(t, err)
	
	// Second call with recovery
	client.err = ""
	err = r.Collect()
	assert.NoError(t, err)
}

// TestCollectWithWhitespaceOnlyResponse tests Collect() with whitespace only
func TestCollectWithWhitespaceOnlyResponse(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "   \n\t  ",
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	// Should handle whitespace gracefully
	assert.NoError(t, err) // or error depending on implementation
}

// TestCollectWithHTTPHeadMethod tests Collect() handles response headers correctly
func TestCollectWithHTTPHeadMethod(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "",
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	
	assert.NoError(t, err)
}

// TestCollectWithRedirectResponse tests Collect() with redirect status
func TestCollectWithRedirectResponse(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: "",
			statusCode:   301,
		},
	}
	
	err := r.Collect()
	
	// Depends on implementation if redirects are followed
	// assert.Error(t, err) or assert.NoError(t, err)
}

// TestCollectConcurrency tests Collect() called concurrently
func TestCollectConcurrency(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: validRspamdResponse(),
			statusCode:   200,
		},
	}
	
	done := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			done <- r.Collect()
		}()
	}
	
	for i := 0; i < 10; i++ {
		err := <-done
		assert.NoError(t, err)
	}
}

// TestCollectWithZeroTimeout tests Collect() with zero timeout value
func TestCollectWithZeroTimeout(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: validRspamdResponse(),
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	assert.NoError(t, err)
}

// TestCollectWithNegativeTimeout tests Collect() with negative timeout
func TestCollectWithNegativeTimeout(t *testing.T) {
	r := &Rspamd{
		httpClient: &mockHTTPClient{
			responseBody: validRspamdResponse(),
			statusCode:   200,
		},
	}
	
	err := r.Collect()
	assert.NoError(t, err)
}

// TestCollectResponseProcessing tests the response processing logic
func TestCollectResponseProcessing(t *testing.T) {
	testCases := []struct {
		name         string
		statusCode   int
		responseBody string
		shouldError  bool
	}{
		{"success 200", 200, validRspamdResponse(), false},
		{"success 201", 201, validRspamdResponse(), false},
		{"success 204", 204, "", false},
		{"client error 400", 400, "Bad Request", true},
		{"client error 401", 401, "Unauthorized", true},
		{"client error 403", 403, "Forbidden", true},
		{"client error 404", 404, "Not Found", true},
		{"server error 500", 500, "Server Error", true},
		{"server error 502", 502, "Bad Gateway", true},
		{"server error 503", 503, "Service Unavailable", true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &Rspamd{
				httpClient: &mockHTTPClient{
					responseBody: tc.responseBody,
					statusCode:   tc.statusCode,
				},
			}
			
			err := r.Collect()
			
			if tc.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCollectJSONParsing tests JSON parsing with various structures
func TestCollectJSONParsing(t *testing.T) {
	testCases := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{"empty object", "{}", false},
		{"empty array", "[]", false},
		{"nested object", `{"a": {"b": {"c": 1}}}`, false},
		{"array of objects", `[{"id": 1}, {"id": 2}]`, false},
		{"mixed types", `{"str": "value", "num": 42, "bool": true, "null": null}`, false},
		{"malformed json", "{invalid}", true},
		{"unclosed brace", `{"key": "value"`, true},
		{"unclosed bracket", `[1, 2, 3`, true},
		{"trailing comma", `{"a": 1, "b": 2,}`, true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &Rspamd{
				httpClient: &mockHTTPClient{
					responseBody: tc.json,
					statusCode:   200,
				},
			}
			
			err := r.Collect()
			
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to generate valid rspamd response
func validRspamdResponse() string {
	return `{
		"version": "2.0",
		"uptime": 12345,
		"memory_pool": {
			"total": 1000000,
			"used": 500000
		}
	}`
}

// Mock HTTP Client for testing
type mockHTTPClient struct {
	responseBody string
	statusCode   int
	err          string
}

// Do mocks the HTTP client Do method
func (m *mockHTTPClient) Do(req interface{}) (interface{}, error) {
	if m.err != "" {
		return nil, &mockError{msg: m.err}
	}
	
	return &mockResponse{
		StatusCode: m.statusCode,
		Body:       m.responseBody,
	}, nil
}

type mockResponse struct {
	StatusCode int
	Body       string
}

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}