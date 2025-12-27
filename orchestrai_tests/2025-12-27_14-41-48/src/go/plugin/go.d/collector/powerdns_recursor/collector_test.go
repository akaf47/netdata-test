package powerdns_recursor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCollector(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
}

func TestCollectorInit_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL: server.URL,
	}

	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.client)
}

func TestCollectorInit_InvalidURL(t *testing.T) {
	c := &Collector{
		URL: "not-a-valid-url",
	}

	err := c.Init()
	assert.Error(t, err)
}

func TestCollectorInit_EmptyURL(t *testing.T) {
	c := &Collector{
		URL: "",
	}

	err := c.Init()
	assert.Error(t, err)
}

func TestCollectorInit_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := &Collector{
		URL:     server.URL,
		Timeout: 100 * time.Millisecond,
	}

	err := c.Init()
	assert.Error(t, err)
}

func TestCollectorInit_WithDefaultTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:     server.URL,
		Timeout: 5 * time.Second,
	}

	err := c.Init()
	assert.NoError(t, err)
}

func TestCollectorCheck_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	err := c.Check()
	assert.NoError(t, err)
}

func TestCollectorCheck_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	err := c.Check()
	assert.Error(t, err)
}

func TestCollectorCheck_ConnectionRefused(t *testing.T) {
	c := &Collector{
		URL: "http://127.0.0.1:1",
	}

	err := c.Check()
	assert.Error(t, err)
}

func TestCollectorCheck_NoClient(t *testing.T) {
	c := &Collector{
		URL: "http://localhost:8001",
	}

	err := c.Check()
	assert.Error(t, err)
}

func TestCollectorCollect_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
		fmt.Fprintf(w, "answers=950\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.NotNil(t, m)
	assert.Len(t, m, 2)
}

func TestCollectorCollect_Nil(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	assert.Nil(t, m)
}

func TestCollectorCollect_NoClient(t *testing.T) {
	c := &Collector{
		URL: "http://localhost:8001",
	}

	m := c.Collect()
	assert.Nil(t, m)
}

func TestCollectorName(t *testing.T) {
	c := &Collector{}
	name := c.Name()
	assert.NotEmpty(t, name)
	assert.Equal(t, "powerdns_recursor", name)
}

func TestCollectorDescription(t *testing.T) {
	c := &Collector{}
	desc := c.Description()
	assert.NotEmpty(t, desc)
}

func TestCollectorInit_URLValidation(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantError bool
	}{
		{"valid http", "http://localhost:8001", false},
		{"valid https", "https://localhost:8001", false},
		{"empty url", "", true},
		{"invalid scheme", "ftp://localhost:8001", true},
		{"malformed url", "ht!tp://localhost:8001", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "questions=1000\n")
			}))
			defer server.Close()

			if tt.url == "" {
				c := &Collector{
					URL: tt.url,
				}
				err := c.Init()
				assert.Error(t, err)
			}
		})
	}
}

func TestCollectorConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:     server.URL,
		Timeout: 10 * time.Second,
	}

	err := c.Init()
	assert.NoError(t, err)
	assert.Equal(t, server.URL, c.URL)
}

func TestCollectorMultipleInit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL: server.URL,
	}

	err1 := c.Init()
	assert.NoError(t, err1)

	err2 := c.Init()
	assert.NoError(t, err2)
}

func TestCollectorMultipleCollect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
		fmt.Fprintf(w, "answers=950\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m1 := c.Collect()
	require.NotNil(t, m1)

	m2 := c.Collect()
	require.NotNil(t, m2)

	assert.Equal(t, m1, m2)
}

func TestCollectorCheck_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	err := c.Check()
	assert.Error(t, err)
}

func TestCollectorInit_HTTPClientCreation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL: server.URL,
	}

	assert.Nil(t, c.client)
	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.client)
}

func TestCollectorCollect_WithMultipleMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "all-questions=5000\n")
		fmt.Fprintf(w, "answers0=100\n")
		fmt.Fprintf(w, "answers1=200\n")
		fmt.Fprintf(w, "answers2=300\n")
		fmt.Fprintf(w, "answers3=400\n")
		fmt.Fprintf(w, "answers4=500\n")
		fmt.Fprintf(w, "answers5=600\n")
		fmt.Fprintf(w, "nxdomain=50\n")
		fmt.Fprintf(w, "nodata=75\n")
		fmt.Fprintf(w, "servfail=25\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.NotNil(t, m)
	assert.Equal(t, 10, len(m))
	assert.Equal(t, uint64(5000), m["all-questions"])
	assert.Equal(t, uint64(100), m["answers0"])
}

func TestCollectorCollect_MalformedData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "invalid-data-format")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	assert.Nil(t, m)
}

func TestCollectorCheck_ValidMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
		fmt.Fprintf(w, "answers=950\n")
		fmt.Fprintf(w, "servfail-answers=10\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	err := c.Check()
	assert.NoError(t, err)
}

func TestCollectorInit_ClientTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:     server.URL,
		Timeout: 500 * time.Millisecond,
	}

	err := c.Init()
	assert.NoError(t, err)
}

func TestCollectorCollect_PartialFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
		fmt.Fprintf(w, "invalid=not-number\n")
		fmt.Fprintf(w, "answers=950\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.NotNil(t, m)
	assert.Equal(t, uint64(1000), m["questions"])
	assert.Equal(t, uint64(950), m["answers"])
}

func TestCollectorName_NotEmpty(t *testing.T) {
	c := &Collector{}
	name := c.Name()
	assert.NotEmpty(t, name)
}

func TestCollectorDescription_NotEmpty(t *testing.T) {
	c := &Collector{}
	desc := c.Description()
	assert.NotEmpty(t, desc)
}