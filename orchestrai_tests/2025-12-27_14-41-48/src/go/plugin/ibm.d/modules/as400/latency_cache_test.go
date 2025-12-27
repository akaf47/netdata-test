package as400

import (
	"sync"
	"testing"
	"time"
)

// TestLatencyCacheInitialization tests cache initialization
func TestLatencyCacheInitialization(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should create new cache with default capacity",
			test: func(t *testing.T) {
				cache := NewLatencyCache()
				if cache == nil {
					t.Fatal("expected cache to be initialized, got nil")
				}
			},
		},
		{
			name: "should create cache with specified capacity",
			test: func(t *testing.T) {
				capacity := 1000
				cache := NewLatencyCacheWithCapacity(capacity)
				if cache == nil {
					t.Fatal("expected cache to be initialized, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestLatencyCacheSet tests adding items to cache
func TestLatencyCacheSet(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     float64
		expectErr bool
	}{
		{
			name:      "should set value with valid key and value",
			key:       "query_1",
			value:     123.45,
			expectErr: false,
		},
		{
			name:      "should set value with zero value",
			key:       "query_2",
			value:     0.0,
			expectErr: false,
		},
		{
			name:      "should set value with negative value",
			key:       "query_3",
			value:     -1.5,
			expectErr: false,
		},
		{
			name:      "should set value with empty key",
			key:       "",
			value:     100.0,
			expectErr: false,
		},
		{
			name:      "should set value with large value",
			key:       "query_large",
			value:     999999999.99,
			expectErr: false,
		},
		{
			name:      "should overwrite existing key",
			key:       "query_1",
			value:     999.99,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			err := cache.Set(tt.key, tt.value)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error=%v, got error=%v", tt.expectErr, err)
			}
		})
	}
}

// TestLatencyCacheGet tests retrieving items from cache
func TestLatencyCacheGet(t *testing.T) {
	tests := []struct {
		name      string
		setupKey  string
		setupVal  float64
		getKey    string
		expectVal float64
		expectOk  bool
	}{
		{
			name:      "should get existing value",
			setupKey:  "query_1",
			setupVal:  123.45,
			getKey:    "query_1",
			expectVal: 123.45,
			expectOk:  true,
		},
		{
			name:      "should return false for non-existent key",
			setupKey:  "query_1",
			setupVal:  123.45,
			getKey:    "query_999",
			expectVal: 0,
			expectOk:  false,
		},
		{
			name:      "should get zero value",
			setupKey:  "query_zero",
			setupVal:  0.0,
			getKey:    "query_zero",
			expectVal: 0.0,
			expectOk:  true,
		},
		{
			name:      "should get negative value",
			setupKey:  "query_neg",
			setupVal:  -100.5,
			getKey:    "query_neg",
			expectVal: -100.5,
			expectOk:  true,
		},
		{
			name:      "should handle empty key",
			setupKey:  "",
			setupVal:  50.0,
			getKey:    "",
			expectVal: 50.0,
			expectOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			cache.Set(tt.setupKey, tt.setupVal)
			val, ok := cache.Get(tt.getKey)
			if ok != tt.expectOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectOk, ok)
			}
			if ok && val != tt.expectVal {
				t.Errorf("expected value=%v, got value=%v", tt.expectVal, val)
			}
		})
	}
}

// TestLatencyCacheDelete tests removing items from cache
func TestLatencyCacheDelete(t *testing.T) {
	tests := []struct {
		name     string
		setupKey string
		deleteKey string
		expectOk bool
	}{
		{
			name:     "should delete existing key",
			setupKey: "query_1",
			deleteKey: "query_1",
			expectOk: true,
		},
		{
			name:     "should return false when deleting non-existent key",
			setupKey: "query_1",
			deleteKey: "query_999",
			expectOk: false,
		},
		{
			name:     "should return false when deleting from empty cache",
			setupKey: "",
			deleteKey: "query_1",
			expectOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			if tt.setupKey != "" {
				cache.Set(tt.setupKey, 100.0)
			}
			ok := cache.Delete(tt.deleteKey)
			if ok != tt.expectOk {
				t.Errorf("expected ok=%v, got ok=%v", tt.expectOk, ok)
			}
			// Verify deletion
			_, exists := cache.Get(tt.deleteKey)
			if exists {
				t.Error("expected key to be deleted, but it still exists")
			}
		})
	}
}

// TestLatencyCacheClear tests clearing entire cache
func TestLatencyCacheClear(t *testing.T) {
	tests := []struct {
		name     string
		itemsToAdd int
	}{
		{
			name:     "should clear empty cache",
			itemsToAdd: 0,
		},
		{
			name:     "should clear cache with single item",
			itemsToAdd: 1,
		},
		{
			name:     "should clear cache with multiple items",
			itemsToAdd: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			for i := 0; i < tt.itemsToAdd; i++ {
				key := "query_" + string(rune(i))
				cache.Set(key, float64(i*10))
			}
			cache.Clear()
			
			// Verify all items are deleted
			for i := 0; i < tt.itemsToAdd; i++ {
				key := "query_" + string(rune(i))
				_, exists := cache.Get(key)
				if exists {
					t.Errorf("expected key %s to be cleared, but it still exists", key)
				}
			}
		})
	}
}

// TestLatencyCacheSize tests getting cache size
func TestLatencyCacheSize(t *testing.T) {
	tests := []struct {
		name     string
		itemsToAdd int
		expectedSize int
	}{
		{
			name:     "should return 0 for empty cache",
			itemsToAdd: 0,
			expectedSize: 0,
		},
		{
			name:     "should return correct size with single item",
			itemsToAdd: 1,
			expectedSize: 1,
		},
		{
			name:     "should return correct size with multiple items",
			itemsToAdd: 50,
			expectedSize: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			for i := 0; i < tt.itemsToAdd; i++ {
				key := "query_" + string(rune(i))
				cache.Set(key, float64(i))
			}
			size := cache.Size()
			if size != tt.expectedSize {
				t.Errorf("expected size=%d, got size=%d", tt.expectedSize, size)
			}
		})
	}
}

// TestLatencyCacheExists tests checking if key exists
func TestLatencyCacheExists(t *testing.T) {
	tests := []struct {
		name     string
		setupKey string
		checkKey string
		expected bool
	}{
		{
			name:     "should return true for existing key",
			setupKey: "query_1",
			checkKey: "query_1",
			expected: true,
		},
		{
			name:     "should return false for non-existent key",
			setupKey: "query_1",
			checkKey: "query_999",
			expected: false,
		},
		{
			name:     "should return false for empty cache",
			setupKey: "",
			checkKey: "query_1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			if tt.setupKey != "" {
				cache.Set(tt.setupKey, 100.0)
			}
			exists := cache.Exists(tt.checkKey)
			if exists != tt.expected {
				t.Errorf("expected exists=%v, got exists=%v", tt.expected, exists)
			}
		})
	}
}

// TestLatencyCacheGetAll tests retrieving all items
func TestLatencyCacheGetAll(t *testing.T) {
	tests := []struct {
		name     string
		itemsToAdd map[string]float64
	}{
		{
			name:     "should return empty map for empty cache",
			itemsToAdd: make(map[string]float64),
		},
		{
			name:     "should return single item",
			itemsToAdd: map[string]float64{"query_1": 100.0},
		},
		{
			name:     "should return multiple items",
			itemsToAdd: map[string]float64{
				"query_1": 100.0,
				"query_2": 200.0,
				"query_3": 300.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			for k, v := range tt.itemsToAdd {
				cache.Set(k, v)
			}
			all := cache.GetAll()
			if len(all) != len(tt.itemsToAdd) {
				t.Errorf("expected %d items, got %d", len(tt.itemsToAdd), len(all))
			}
			for k, v := range tt.itemsToAdd {
				if val, ok := all[k]; !ok || val != v {
					t.Errorf("expected key %s with value %v, got %v", k, v, val)
				}
			}
		})
	}
}

// TestLatencyCacheConcurrentAccess tests thread safety
func TestLatencyCacheConcurrentAccess(t *testing.T) {
	tests := []struct {
		name        string
		goroutines  int
		operationsPerGoroutine int
	}{
		{
			name:        "should handle concurrent writes from 10 goroutines",
			goroutines:  10,
			operationsPerGoroutine: 100,
		},
		{
			name:        "should handle concurrent writes from 100 goroutines",
			goroutines:  100,
			operationsPerGoroutine: 10,
		},
		{
			name:        "should handle concurrent reads and writes",
			goroutines:  50,
			operationsPerGoroutine: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCache()
			var wg sync.WaitGroup
			
			for i := 0; i < tt.goroutines; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for j := 0; j < tt.operationsPerGoroutine; j++ {
						key := "query_" + string(rune(id))
						cache.Set(key, float64(j))
						cache.Get(key)
						if j%10 == 0 {
							cache.Delete(key)
						}
					}
				}(i)
			}
			wg.Wait()
		})
	}
}

// TestLatencyCacheExpiration tests cache expiration if supported
func TestLatencyCacheExpiration(t *testing.T) {
	tests := []struct {
		name           string
		ttl            time.Duration
		waitDuration   time.Duration
		shouldExist    bool
	}{
		{
			name:           "should expire item after TTL",
			ttl:            100 * time.Millisecond,
			waitDuration:   200 * time.Millisecond,
			shouldExist:    false,
		},
		{
			name:           "should not expire item before TTL",
			ttl:            500 * time.Millisecond,
			waitDuration:   100 * time.Millisecond,
			shouldExist:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCacheWithTTL(tt.ttl)
			cache.Set("query_1", 100.0)
			time.Sleep(tt.waitDuration)
			_, ok := cache.Get("query_1")
			if ok != tt.shouldExist {
				t.Errorf("expected exists=%v, got exists=%v", tt.shouldExist, ok)
			}
		})
	}
}

// TestLatencyCacheEdgeCases tests boundary and edge cases
func TestLatencyCacheEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should handle very long key names",
			test: func(t *testing.T) {
				cache := NewLatencyCache()
				longKey := ""
				for i := 0; i < 10000; i++ {
					longKey += "a"
				}
				err := cache.Set(longKey, 123.45)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				val, ok := cache.Get(longKey)
				if !ok || val != 123.45 {
					t.Error("failed to retrieve value with long key")
				}
			},
		},
		{
			name: "should handle special characters in keys",
			test: func(t *testing.T) {
				cache := NewLatencyCache()
				specialKey := "query!@#$%^&*()"
				cache.Set(specialKey, 456.78)
				val, ok := cache.Get(specialKey)
				if !ok || val != 456.78 {
					t.Error("failed to handle special characters in key")
				}
			},
		},
		{
			name: "should handle MAX_FLOAT values",
			test: func(t *testing.T) {
				cache := NewLatencyCache()
				maxFloat := 1.7976931348623157e+308
				cache.Set("max", maxFloat)
				val, ok := cache.Get("max")
				if !ok || val != maxFloat {
					t.Error("failed to handle max float value")
				}
			},
		},
		{
			name: "should handle MIN_FLOAT values",
			test: func(t *testing.T) {
				cache := NewLatencyCache()
				minFloat := -1.7976931348623157e+308
				cache.Set("min", minFloat)
				val, ok := cache.Get("min")
				if !ok || val != minFloat {
					t.Error("failed to handle min float value")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestLatencyCacheCapacityManagement tests capacity and eviction
func TestLatencyCacheCapacityManagement(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		itemsToAdd int
	}{
		{
			name:     "should not exceed specified capacity",
			capacity: 100,
			itemsToAdd: 200,
		},
		{
			name:     "should handle capacity of 1",
			capacity: 1,
			itemsToAdd: 10,
		},
		{
			name:     "should handle large capacity",
			capacity: 10000,
			itemsToAdd: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLatencyCacheWithCapacity(tt.capacity)
			for i := 0; i < tt.itemsToAdd; i++ {
				key := "query_" + string(rune(i%1000))
				cache.Set(key, float64(i))
			}
			size := cache.Size()
			if size > tt.capacity {
				t.Errorf("cache size %d exceeded capacity %d", size, tt.capacity)
			}
		})
	}
}