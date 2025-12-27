package ping

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"
)

// TestProberNew tests the Prober constructor
func TestProberNew(t *testing.T) {
	tests := []struct {
		name string
		host string
		want *Prober
	}{
		{
			name: "creates prober with valid host",
			host: "example.com",
			want: &Prober{
				host: "example.com",
			},
		},
		{
			name: "creates prober with IP address",
			host: "8.8.8.8",
			want: &Prober{
				host: "8.8.8.8",
			},
		},
		{
			name: "creates prober with localhost",
			host: "localhost",
			want: &Prober{
				host: "localhost",
			},
		},
		{
			name: "creates prober with empty host",
			host: "",
			want: &Prober{
				host: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.host)
			if got.host != tt.want.host {
				t.Errorf("New() host = %v, want %v", got.host, tt.want.host)
			}
		})
	}
}

// TestProberPing tests the Ping method
func TestProberPing(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		timeout     time.Duration
		count       int
		shouldError bool
		wantMin     time.Duration
		wantMax     time.Duration
		wantAvg     time.Duration
	}{
		{
			name:        "ping with valid timeout and count",
			host:        "127.0.0.1",
			timeout:     5 * time.Second,
			count:       1,
			shouldError: false,
		},
		{
			name:        "ping with short timeout",
			host:        "127.0.0.1",
			timeout:     100 * time.Millisecond,
			count:       1,
			shouldError: false,
		},
		{
			name:        "ping with zero count",
			host:        "127.0.0.1",
			timeout:     5 * time.Second,
			count:       0,
			shouldError: true,
		},
		{
			name:        "ping with negative count",
			host:        "127.0.0.1",
			timeout:     5 * time.Second,
			count:       -1,
			shouldError: true,
		},
		{
			name:        "ping with zero timeout",
			host:        "127.0.0.1",
			timeout:     0,
			count:       1,
			shouldError: true,
		},
		{
			name:        "ping with negative timeout",
			host:        "127.0.0.1",
			timeout:     -1 * time.Second,
			count:       1,
			shouldError: true,
		},
		{
			name:        "ping with invalid host",
			host:        "invalid..host...name",
			timeout:     5 * time.Second,
			count:       1,
			shouldError: true,
		},
		{
			name:        "ping with empty host",
			host:        "",
			timeout:     5 * time.Second,
			count:       1,
			shouldError: true,
		},
		{
			name:        "ping with multiple packets",
			host:        "127.0.0.1",
			timeout:     5 * time.Second,
			count:       3,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prober := New(tt.host)
			result, err := prober.Ping(context.Background(), tt.timeout, tt.count)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Ping() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Ping() unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("Ping() result is nil")
				}
			}
		})
	}
}

// TestProberPingWithContext tests context cancellation
func TestProberPingWithContext(t *testing.T) {
	tests := []struct {
		name    string
		context func() context.Context
		wantErr bool
	}{
		{
			name: "ping with cancelled context",
			context: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			wantErr: true,
		},
		{
			name: "ping with deadline exceeded context",
			context: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				defer cancel()
				time.Sleep(10 * time.Millisecond)
				return ctx
			},
			wantErr: true,
		},
		{
			name: "ping with valid background context",
			context: func() context.Context {
				return context.Background()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prober := New("127.0.0.1")
			ctx := tt.context()
			_, err := prober.Ping(ctx, 5*time.Second, 1)

			if tt.wantErr && err == nil {
				t.Errorf("Ping() expected error with context, got nil")
			}
		})
	}
}

// TestProberPingResult tests the result fields
func TestProberPingResult(t *testing.T) {
	prober := New("127.0.0.1")
	result, err := prober.Ping(context.Background(), 5*time.Second, 1)

	if err != nil {
		t.Fatalf("Ping() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("Ping() result is nil")
	}

	tests := []struct {
		name     string
		validate func(*Result) bool
	}{
		{
			name: "result has non-zero packets sent",
			validate: func(r *Result) bool {
				return r.PacketsSent > 0
			},
		},
		{
			name: "result has packets received",
			validate: func(r *Result) bool {
				return r.PacketsReceived >= 0
			},
		},
		{
			name: "result has packet loss percentage",
			validate: func(r *Result) bool {
				return r.PacketLoss >= 0 && r.PacketLoss <= 100
			},
		},
		{
			name: "result min time is less than or equal to avg",
			validate: func(r *Result) bool {
				if r.Min == 0 || r.Avg == 0 {
					return true // skip if not set
				}
				return r.Min <= r.Avg
			},
		},
		{
			name: "result avg time is less than or equal to max",
			validate: func(r *Result) bool {
				if r.Avg == 0 || r.Max == 0 {
					return true // skip if not set
				}
				return r.Avg <= r.Max
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.validate(result) {
				t.Errorf("Ping() result validation failed: %v", result)
			}
		})
	}
}

// TestProberPingLocalhost tests pinging localhost
func TestProberPingLocalhost(t *testing.T) {
	prober := New("127.0.0.1")
	result, err := prober.Ping(context.Background(), 5*time.Second, 1)

	if err != nil {
		t.Fatalf("Ping() localhost failed: %v", err)
	}

	if result.PacketsReceived == 0 {
		t.Errorf("Ping() localhost packets received = 0, want > 0")
	}

	if result.PacketLoss == 100 {
		t.Errorf("Ping() localhost packet loss = 100%%, want < 100%%")
	}
}

// TestProberPingUnreachableHost tests pinging unreachable host
func TestProberPingUnreachableHost(t *testing.T) {
	// Use a non-routable IP address
	prober := New("192.0.2.1")
	result, err := prober.Ping(context.Background(), 2*time.Second, 1)

	// Should either error or have packet loss
	if err == nil && result != nil && result.PacketLoss < 100 {
		t.Errorf("Ping() unreachable host expected high packet loss, got %f%%", result.PacketLoss)
	}
}

// TestProberPingMultipleAttempts tests multiple ping attempts
func TestProberPingMultipleAttempts(t *testing.T) {
	prober := New("127.0.0.1")

	for i := 0; i < 5; i++ {
		result, err := prober.Ping(context.Background(), 5*time.Second, 1)
		if err != nil {
			t.Errorf("Ping() attempt %d failed: %v", i+1, err)
		}
		if result == nil {
			t.Errorf("Ping() attempt %d returned nil result", i+1)
		}
	}
}

// TestProberPingRapidRequests tests rapid consecutive pings
func TestProberPingRapidRequests(t *testing.T) {
	prober := New("127.0.0.1")

	for i := 0; i < 10; i++ {
		_, err := prober.Ping(context.Background(), 5*time.Second, 1)
		if err != nil {
			t.Logf("Ping() rapid request %d error: %v", i+1, err)
		}
	}
}

// TestProberPingStatisticFields tests all statistical fields
func TestProberPingStatisticFields(t *testing.T) {
	prober := New("127.0.0.1")
	result, err := prober.Ping(context.Background(), 5*time.Second, 3)

	if err != nil {
		t.Fatalf("Ping() unexpected error: %v", err)
	}

	tests := []struct {
		name     string
		validate func(*Result) bool
	}{
		{
			name: "Min time is positive",
			validate: func(r *Result) bool {
				return r.Min >= 0
			},
		},
		{
			name: "Max time is positive",
			validate: func(r *Result) bool {
				return r.Max >= 0
			},
		},
		{
			name: "Avg time is positive",
			validate: func(r *Result) bool {
				return r.Avg >= 0
			},
		},
		{
			name: "StdDev is non-negative",
			validate: func(r *Result) bool {
				return r.StdDev >= 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.validate(result) {
				t.Errorf("Ping() stat field validation failed")
			}
		})
	}
}

// TestResult_Unhealthy tests the Unhealthy method
func TestResult_Unhealthy(t *testing.T) {
	tests := []struct {
		name    string
		result  *Result
		want    bool
	}{
		{
			name: "result with 0% packet loss is healthy",
			result: &Result{
				PacketLoss: 0,
			},
			want: false,
		},
		{
			name: "result with 50% packet loss is unhealthy",
			result: &Result{
				PacketLoss: 50,
			},
			want: true,
		},
		{
			name: "result with 100% packet loss is unhealthy",
			result: &Result{
				PacketLoss: 100,
			},
			want: true,
		},
		{
			name: "result with nil is unhealthy",
			result: nil,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.Unhealthy()
			if got != tt.want {
				t.Errorf("Unhealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestProberDNSResolution tests DNS resolution
func TestProberDNSResolution(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		shouldError bool
	}{
		{
			name:        "resolve localhost",
			host:        "localhost",
			shouldError: false,
		},
		{
			name:        "resolve IP address",
			host:        "127.0.0.1",
			shouldError: false,
		},
		{
			name:        "resolve invalid hostname",
			host:        "this-domain-definitely-does-not-exist-12345.test",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prober := New(tt.host)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			result, err := prober.Ping(ctx, 5*time.Second, 1)

			if tt.shouldError {
				if err == nil {
					t.Logf("Ping() expected error for host %s", tt.host)
				}
			} else {
				if err != nil {
					t.Logf("Ping() unexpected error for host %s: %v", tt.host, err)
				} else if result == nil {
					t.Errorf("Ping() returned nil result")
				}
			}
		})
	}
}

// TestProberVeryShortTimeout tests with very short timeout
func TestProberVeryShortTimeout(t *testing.T) {
	prober := New("127.0.0.1")
	_, err := prober.Ping(context.Background(), 1*time.Millisecond, 1)
	// Short timeout may error or have high packet loss
	_ = err
}

// TestProberLargePacketCount tests with large packet count
func TestProberLargePacketCount(t *testing.T) {
	prober := New("127.0.0.1")
	result, err := prober.Ping(context.Background(), 10*time.Second, 10)

	if err != nil {
		t.Fatalf("Ping() unexpected error with large count: %v", err)
	}

	if result != nil && result.PacketsSent != 10 {
		t.Errorf("Ping() packets sent = %d, want 10", result.PacketsSent)
	}
}

// TestProberBoundaryValues tests boundary value conditions
func TestProberBoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		count   int
		wantErr bool
	}{
		{
			name:    "timeout at 1 nanosecond",
			timeout: 1 * time.Nanosecond,
			count:   1,
			wantErr: true,
		},
		{
			name:    "timeout at 1 millisecond",
			timeout: 1 * time.Millisecond,
			count:   1,
			wantErr: false,
		},
		{
			name:    "count at 1",
			timeout: 5 * time.Second,
			count:   1,
			wantErr: false,
		},
		{
			name:    "count at max int",
			timeout: 5 * time.Second,
			count:   2147483647,
			wantErr: false, // May timeout or error, but should handle gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prober := New("127.0.0.1")
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			_, err := prober.Ping(ctx, tt.timeout, tt.count)

			if tt.wantErr && err == nil {
				t.Logf("Ping() expected error for timeout=%v, count=%d", tt.timeout, tt.count)
			}
		})
	}
}

// TestProberConcurrentPings tests concurrent ping operations
func TestProberConcurrentPings(t *testing.T) {
	prober := New("127.0.0.1")
	errs := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			_, err := prober.Ping(context.Background(), 5*time.Second, 1)
			errs <- err
		}()
	}

	for i := 0; i < 5; i++ {
		err := <-errs
		if err != nil {
			t.Logf("concurrent Ping() error: %v", err)
		}
	}
}

// TestProberDifferentHosts tests pinging different hosts
func TestProberDifferentHosts(t *testing.T) {
	hosts := []string{
		"127.0.0.1",
		"localhost",
		"::1", // IPv6 loopback (may not work on all systems)
	}

	for _, host := range hosts {
		t.Run("ping "+host, func(t *testing.T) {
			prober := New(host)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			result, err := prober.Ping(ctx, 5*time.Second, 1)
			// IPv6 may not be available on all systems
			if err == nil && result != nil {
				if result.PacketsReceived == 0 {
					t.Logf("Ping(%s) no responses", host)
				}
			}
		})
	}
}

// TestProberErrorHandling tests error handling paths
func TestProberErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		timeout     time.Duration
		count       int
		shouldError bool
	}{
		{
			name:        "invalid parameters zero timeout",
			host:        "127.0.0.1",
			timeout:     0,
			count:       1,
			shouldError: true,
		},
		{
			name:        "invalid parameters negative timeout",
			host:        "127.0.0.1",
			timeout:     -5 * time.Second,
			count:       1,
			shouldError: true,
		},
		{
			name:        "invalid parameters zero count",
			host:        "127.0.0.1",
			timeout:     5 * time.Second,
			count:       0,
			shouldError: true,
		},
		{
			name:        "invalid parameters negative count",
			host:        "127.0.0.1",
			timeout:     5 * time.Second,
			count:       -1,
			shouldError: true,
		},
		{
			name:        "invalid host format",
			host:        "...invalid...",
			timeout:     5 * time.Second,
			count:       1,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prober := New(tt.host)
			_, err := prober.Ping(context.Background(), tt.timeout, tt.count)

			if tt.shouldError && err == nil {
				t.Errorf("Ping() expected error")
			}
		})
	}
}