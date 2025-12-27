package powerdns_recursor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
		fmt.Fprintf(w, "answers=950\n")
		fmt.Fprintf(w, "servfail-answers=10\n")
		fmt.Fprintf(w, "nxdomain-answers=20\n")
		fmt.Fprintf(w, "nodata-answers=15\n")
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

func TestCollect_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_InvalidURL(t *testing.T) {
	c := &Collector{
		URL: "http://invalid-host-that-does-not-exist-12345:8001",
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_MalformedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=not-a-number\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_PartialMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
		fmt.Fprintf(w, "invalid-line\n")
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

func TestCollect_ZeroValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=0\n")
		fmt.Fprintf(w, "answers=0\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.NotNil(t, m)
	assert.Equal(t, uint64(0), m["questions"])
	assert.Equal(t, uint64(0), m["answers"])
}

func TestCollect_LargeValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=9223372036854775807\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.NotNil(t, m)
	assert.Equal(t, uint64(9223372036854775807), m["questions"])
}

func TestCollect_MultipleMetricsSameKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000\n")
		fmt.Fprintf(w, "questions=2000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.NotNil(t, m)
	assert.Equal(t, uint64(2000), m["questions"])
}

func TestCollect_WhitespaceHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "  questions=1000  \n")
		fmt.Fprintf(w, "\tanswers=950\t\n")
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

func TestCollect_NoValueAfterEquals(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_NegativeValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=-100\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_FloatingPointValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions=1000.5\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_SpecialCharactersInKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions-special=1000\n")
		fmt.Fprintf(w, "answers_count=950\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.NotNil(t, m)
	assert.Equal(t, uint64(1000), m["questions-special"])
	assert.Equal(t, uint64(950), m["answers_count"])
}

func TestCollect_ConnectionTimeout(t *testing.T) {
	c := &Collector{
		URL: "http://127.0.0.1:1",
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_AllMetricsTypes(t *testing.T) {
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
}

func TestCollect_EmptyKeyName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "=1000\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_OnlyEqualsSign(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "=\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_NoEqualsSign(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "questions-without-value\n")
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_StatusBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_StatusUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_StatusNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}

func TestCollect_StatusServiceUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	c := &Collector{
		URL:    server.URL,
		client: &http.Client{},
	}

	m := c.Collect()
	require.Nil(t, m)
}