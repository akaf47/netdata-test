package netdataapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid URL",
			url:     "http://localhost:19999",
			wantErr: false,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && client == nil {
				t.Errorf("NewClient() returned nil client with no error")
			}
		})
	}
}

func TestClient_GetInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/info" {
			t.Errorf("expected path /api/v1/info, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version":"1.0","hostname":"test"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := client.GetInfo(ctx)
	if err != nil {
		t.Errorf("GetInfo() error = %v", err)
		return
	}
	if info == nil {
		t.Errorf("GetInfo() returned nil info")
	}
}

func TestClient_GetInfo_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server error"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	info, err := client.GetInfo(ctx)
	if err == nil {
		t.Errorf("GetInfo() expected error, got nil")
	}
	if info != nil {
		t.Errorf("GetInfo() expected nil info, got %v", info)
	}
}

func TestClient_GetInfo_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	info, err := client.GetInfo(ctx)
	if err == nil {
		t.Errorf("GetInfo() expected context timeout error, got nil")
	}
	if info != nil {
		t.Errorf("GetInfo() expected nil info on timeout, got %v", info)
	}
}

func TestClient_GetChart(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/chart" {
			t.Errorf("expected path /api/v1/chart, got %s", r.URL.Path)
		}
		chart := r.URL.Query().Get("chart")
		if chart != "system.cpu" {
			t.Errorf("expected chart=system.cpu, got %s", chart)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"system.cpu","name":"cpu"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	chart, err := client.GetChart(ctx, "system.cpu")
	if err != nil {
		t.Errorf("GetChart() error = %v", err)
		return
	}
	if chart == nil {
		t.Errorf("GetChart() returned nil chart")
	}
}

func TestClient_GetChart_EmptyChartID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.GetChart(ctx, "")
	if err == nil {
		t.Errorf("GetChart() with empty chartID expected error, got nil")
	}
}

func TestClient_GetData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/data" {
			t.Errorf("expected path /api/v1/data, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"ok","data":[[1,2,3]]}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := client.GetData(ctx, &GetDataRequest{Chart: "system.cpu"})
	if err != nil {
		t.Errorf("GetData() error = %v", err)
		return
	}
	if data == nil {
		t.Errorf("GetData() returned nil data")
	}
}

func TestClient_GetData_NilRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.GetData(ctx, nil)
	if err == nil {
		t.Errorf("GetData() with nil request expected error, got nil")
	}
}

func TestClient_GetData_InvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.GetData(ctx, &GetDataRequest{Chart: "system.cpu"})
	if err == nil {
		t.Errorf("GetData() with invalid JSON expected error, got nil")
	}
}

func TestClient_QueryDatabase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/query" {
			t.Errorf("expected path /api/v2/query, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":{"labels":[],"data":[]}}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := client.QueryDatabase(ctx, "test_query")
	if err != nil {
		t.Errorf("QueryDatabase() error = %v", err)
		return
	}
	if result == nil {
		t.Errorf("QueryDatabase() returned nil result")
	}
}

func TestClient_QueryDatabase_EmptyQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.QueryDatabase(ctx, "")
	if err == nil {
		t.Errorf("QueryDatabase() with empty query expected error, got nil")
	}
}

func TestClient_Close(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestClient_CustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Authorization header, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	client.SetHeader("Authorization", "Bearer test-token")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.GetInfo(ctx)
	if err != nil {
		t.Errorf("GetInfo() with custom header error = %v", err)
	}
}

func TestClient_MultipleRequests(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"test"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 3; i++ {
		_, err = client.GetInfo(ctx)
		if err != nil {
			t.Errorf("GetInfo() iteration %d error = %v", i, err)
		}
	}

	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}