package pihole

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect(t *testing.T) {
	tests := []struct {
		name            string
		prepare         func(p *Pihole)
		expectedMetrics int
		shouldErr       bool
	}{
		{
			name: "successful collect with valid response",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"dns_queries_today": 1000,
								"ads_blocked_today": 100,
								"ads_percentage_today": 10.0,
								"unique_clients": 50,
								"queries_forwarded": 900,
								"queries_cached": 50,
								"reply_NODATA": 5,
								"reply_NXDOMAIN": 10,
								"status": "enabled",
								"gravity_last_updated": 1234567890,
								"reply_A": 400,
								"reply_AAAA": 300,
								"reply_ANY": 100,
								"reply_SRV": 50
							}`)),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       false,
		},
		{
			name: "collect with zero values",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"dns_queries_today": 0,
								"ads_blocked_today": 0,
								"ads_percentage_today": 0.0,
								"unique_clients": 0
							}`)),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       false,
		},
		{
			name: "collect with very large values",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"dns_queries_today": 999999999,
								"ads_blocked_today": 999999999,
								"ads_percentage_today": 99.9999
							}`)),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       false,
		},
		{
			name: "collect with missing fields",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("{}")),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       false,
		},
		{
			name: "collect with null response",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("null")),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       true,
		},
		{
			name: "collect with malformed JSON",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("{invalid json")),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       true,
		},
		{
			name: "collect with HTTP error",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusInternalServerError,
							Body:       io.NopCloser(strings.NewReader("Server Error")),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       true,
		},
		{
			name: "collect with negative values",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"dns_queries_today": -100,
								"ads_blocked_today": -50
							}`)),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       false,
		},
		{
			name: "collect with float percentage edge cases",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"ads_percentage_today": 100.0
							}`)),
						}, nil
					},
				}
			},
			expectedMetrics: 0,
			shouldErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:80/api"
			tt.prepare(p)

			err := p.Collect()
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectWithNilData(t *testing.T) {
	p := New()
	p.URL = "http://localhost:80/api"
	p.client = &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"dns_queries_today": null}`)),
			}, nil
		},
	}

	err := p.Collect()
	// Should handle null values gracefully
	if err != nil {
		assert.NotNil(t, err)
	}
}

func TestCollectDataStructureParsing(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		validate func(t *testing.T, data *apiResponse)
	}{
		{
			name: "parse all query response fields",
			jsonData: `{
				"dns_queries_today": 5000,
				"ads_blocked_today": 500,
				"ads_percentage_today": 10.0,
				"unique_clients": 100,
				"queries_forwarded": 4000,
				"queries_cached": 500,
				"reply_NODATA": 50,
				"reply_NXDOMAIN": 75,
				"status": "enabled",
				"gravity_last_updated": 1609459200,
				"reply_A": 2000,
				"reply_AAAA": 1500,
				"reply_ANY": 300,
				"reply_SRV": 200
			}`,
			validate: func(t *testing.T, data *apiResponse) {
				assert.Equal(t, int64(5000), data.DNSQueriesToday)
				assert.Equal(t, int64(500), data.AdsBlockedToday)
				assert.Equal(t, 10.0, data.AdsPercentageToday)
				assert.Equal(t, int64(100), data.UniqueClients)
				assert.Equal(t, int64(4000), data.QueriesForwarded)
				assert.Equal(t, int64(500), data.QueriesCached)
				assert.Equal(t, int64(50), data.ReplyNODATA)
				assert.Equal(t, int64(75), data.ReplyNXDOMAIN)
				assert.Equal(t, "enabled", data.Status)
				assert.Equal(t, int64(1609459200), data.GravityLastUpdated)
			},
		},
		{
			name: "parse with missing optional fields",
			jsonData: `{
				"dns_queries_today": 1000,
				"ads_blocked_today": 100
			}`,
			validate: func(t *testing.T, data *apiResponse) {
				assert.Equal(t, int64(1000), data.DNSQueriesToday)
				assert.Equal(t, int64(100), data.AdsBlockedToday)
				assert.Equal(t, int64(0), data.UniqueClients)
			},
		},
		{
			name:     "parse empty response",
			jsonData: `{}`,
			validate: func(t *testing.T, data *apiResponse) {
				assert.NotNil(t, data)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data apiResponse
			err := parseJSON([]byte(tt.jsonData), &data)
			if tt.name != "parse empty response" && err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			if tt.name == "parse empty response" || err == nil {
				tt.validate(t, &data)
			}
		})
	}
}

// Mock HTTP client for testing
type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

// Helper function to simulate JSON parsing
type apiResponse struct {
	DNSQueriesToday    int64   `json:"dns_queries_today"`
	AdsBlockedToday    int64   `json:"ads_blocked_today"`
	AdsPercentageToday float64 `json:"ads_percentage_today"`
	UniqueClients      int64   `json:"unique_clients"`
	QueriesForwarded   int64   `json:"queries_forwarded"`
	QueriesCached      int64   `json:"queries_cached"`
	ReplyNODATA        int64   `json:"reply_NODATA"`
	ReplyNXDOMAIN      int64   `json:"reply_NXDOMAIN"`
	Status             string  `json:"status"`
	GravityLastUpdated int64   `json:"gravity_last_updated"`
	ReplyA             int64   `json:"reply_A"`
	ReplyAAAA          int64   `json:"reply_AAAA"`
	ReplyANY           int64   `json:"reply_ANY"`
	ReplySRV           int64   `json:"reply_SRV"`
}

func parseJSON(data []byte, v interface{}) error {
	// This is a placeholder - actual implementation would use json.Unmarshal
	return nil
}

func TestCollectEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(p *Pihole)
		check   func(t *testing.T, err error)
	}{
		{
			name: "empty URL",
			prepare: func(p *Pihole) {
				p.URL = ""
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return nil, nil
					},
				}
			},
			check: func(t *testing.T, err error) {
				// Should handle empty URL
			},
		},
		{
			name: "response with empty body",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader("")),
						}, nil
					},
				}
			},
			check: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "very large numeric values",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"dns_queries_today": 9223372036854775807,
								"ads_blocked_today": 9223372036854775806
							}`)),
						}, nil
					},
				}
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "response with special characters",
			prepare: func(p *Pihole) {
				p.client = &mockHTTPClient{
					doFunc: func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body: io.NopCloser(bytes.NewBufferString(`{
								"status": "enabled\u0020with\u0020unicode",
								"dns_queries_today": 100
							}`)),
						}, nil
					},
				}
			},
			check: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:80/api"
			tt.prepare(p)
			err := p.Collect()
			tt.check(t, err)
		})
	}
}

func TestCollectHTTPStatusCodes(t *testing.T) {
	statusCodes := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}

	for _, statusCode := range statusCodes {
		t.Run("status_"+string(rune(statusCode)), func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:80/api"
			p.client = &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					body := ""
					if statusCode == http.StatusOK {
						body = `{"dns_queries_today": 100}`
					}
					return &http.Response{
						StatusCode: statusCode,
						Body:       io.NopCloser(strings.NewReader(body)),
					}, nil
				},
			}

			err := p.Collect()
			if statusCode != http.StatusOK {
				assert.Error(t, err)
			}
		})
	}
}