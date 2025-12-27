package nginxplus

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test cache initialization and basic operations
func TestNewCache(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "should create new cache instance",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				assert.NotNil(t, cache)
			},
		},
		{
			name: "should initialize empty cache",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				assert.Equal(t, 0, cache.Size())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestCacheSet(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    interface{}
		testFunc func(*testing.T, *Cache)
	}{
		{
			name:  "should set value in cache",
			key:   "key1",
			value: "value1",
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Set("key1", "value1")
				val, exists := cache.Get("key1")
				assert.True(t, exists)
				assert.Equal(t, "value1", val)
			},
		},
		{
			name:  "should overwrite existing value",
			key:   "key1",
			value: "value2",
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Set("key1", "value1")
				cache.Set("key1", "value2")
				val, exists := cache.Get("key1")
				assert.True(t, exists)
				assert.Equal(t, "value2", val)
			},
		},
		{
			name:  "should handle empty key",
			key:   "",
			value: "value",
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Set("", "value")
				val, exists := cache.Get("")
				assert.True(t, exists)
				assert.Equal(t, "value", val)
			},
		},
		{
			name:  "should handle nil value",
			key:   "key1",
			value: nil,
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Set("key1", nil)
				val, exists := cache.Get("key1")
				assert.True(t, exists)
				assert.Nil(t, val)
			},
		},
		{
			name:  "should store different types",
			key:   "int_key",
			value: 42,
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Set("int_key", 42)
				val, exists := cache.Get("int_key")
				assert.True(t, exists)
				assert.Equal(t, 42, val)
			},
		},
		{
			name:  "should store complex types",
			key:   "map_key",
			value: map[string]interface{}{"nested": "value"},
			testFunc: func(t *testing.T, cache *Cache) {
				testMap := map[string]interface{}{"nested": "value"}
				cache.Set("map_key", testMap)
				val, exists := cache.Get("map_key")
				assert.True(t, exists)
				assert.Equal(t, testMap, val)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache()
			tt.testFunc(t, cache)
		})
	}
}

func TestCacheGet(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Cache)
		key      string
		expected interface{}
		exists   bool
	}{
		{
			name: "should retrieve existing value",
			setup: func(c *Cache) {
				c.Set("key1", "value1")
			},
			key:      "key1",
			expected: "value1",
			exists:   true,
		},
		{
			name: "should return false for non-existent key",
			setup: func(c *Cache) {
				// Empty cache
			},
			key:      "nonexistent",
			expected: nil,
			exists:   false,
		},
		{
			name: "should handle empty key lookup",
			setup: func(c *Cache) {
				c.Set("", "empty_key_value")
			},
			key:      "",
			expected: "empty_key_value",
			exists:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache()
			tt.setup(cache)
			val, exists := cache.Get(tt.key)
			assert.Equal(t, tt.exists, exists)
			if exists {
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

func TestCacheDelete(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Cache)
		key      string
		testFunc func(*testing.T, *Cache)
	}{
		{
			name: "should delete existing key",
			setup: func(c *Cache) {
				c.Set("key1", "value1")
			},
			key: "key1",
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Delete("key1")
				_, exists := cache.Get("key1")
				assert.False(t, exists)
			},
		},
		{
			name: "should handle deletion of non-existent key",
			setup: func(c *Cache) {
				// Empty cache
			},
			key: "nonexistent",
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Delete("nonexistent")
				_, exists := cache.Get("nonexistent")
				assert.False(t, exists)
			},
		},
		{
			name: "should delete with empty key",
			setup: func(c *Cache) {
				c.Set("", "value")
			},
			key: "",
			testFunc: func(t *testing.T, cache *Cache) {
				cache.Delete("")
				_, exists := cache.Get("")
				assert.False(t, exists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache()
			tt.setup(cache)
			tt.testFunc(t, cache)
		})
	}
}

func TestCacheClear(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "should clear all entries",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				cache.Set("key1", "value1")
				cache.Set("key2", "value2")
				cache.Set("key3", "value3")
				cache.Clear()
				assert.Equal(t, 0, cache.Size())
			},
		},
		{
			name: "should handle clearing empty cache",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				cache.Clear()
				assert.Equal(t, 0, cache.Size())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestCacheSize(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "should return zero for empty cache",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				assert.Equal(t, 0, cache.Size())
			},
		},
		{
			name: "should return correct size after adding items",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				cache.Set("key1", "value1")
				cache.Set("key2", "value2")
				cache.Set("key3", "value3")
				assert.Equal(t, 3, cache.Size())
			},
		},
		{
			name: "should update size after deletion",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				cache.Set("key1", "value1")
				cache.Set("key2", "value2")
				cache.Delete("key1")
				assert.Equal(t, 1, cache.Size())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestCacheExists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Cache)
		key      string
		expected bool
	}{
		{
			name: "should return true for existing key",
			setup: func(c *Cache) {
				c.Set("key1", "value1")
			},
			key:      "key1",
			expected: true,
		},
		{
			name: "should return false for non-existent key",
			setup: func(c *Cache) {
				// Empty cache
			},
			key:      "nonexistent",
			expected: false,
		},
		{
			name: "should return false for deleted key",
			setup: func(c *Cache) {
				c.Set("key1", "value1")
				c.Delete("key1")
			},
			key:      "key1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache()
			tt.setup(cache)
			exists := cache.Exists(tt.key)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestCacheConcurrency(t *testing.T) {
	t.Run("should handle concurrent operations safely", func(t *testing.T) {
		cache := NewCache()
		done := make(chan bool)

		// Concurrent writes
		for i := 0; i < 10; i++ {
			go func(id int) {
				for j := 0; j < 100; j++ {
					cache.Set("key", id*100+j)
				}
				done <- true
			}(i)
		}

		// Concurrent reads
		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					cache.Get("key")
				}
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 20; i++ {
			<-done
		}

		// Cache should have at least one value
		assert.Greater(t, cache.Size(), 0)
	})
}

func TestCacheMultipleOperations(t *testing.T) {
	t.Run("should handle complex operation sequences", func(t *testing.T) {
		cache := NewCache()

		// Add multiple entries
		for i := 0; i < 5; i++ {
			cache.Set("key"+string(rune(i)), "value"+string(rune(i)))
		}
		assert.Equal(t, 5, cache.Size())

		// Delete some entries
		cache.Delete("key0")
		cache.Delete("key2")
		assert.Equal(t, 3, cache.Size())

		// Verify remaining entries
		for i := 0; i < 5; i++ {
			_, exists := cache.Get("key" + string(rune(i)))
			if i == 0 || i == 2 {
				assert.False(t, exists)
			} else {
				assert.True(t, exists)
			}
		}

		// Clear and verify
		cache.Clear()
		assert.Equal(t, 0, cache.Size())
	})
}

func TestCacheEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "should handle very long keys",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				longKey := ""
				for i := 0; i < 10000; i++ {
					longKey += "a"
				}
				cache.Set(longKey, "value")
				val, exists := cache.Get(longKey)
				assert.True(t, exists)
				assert.Equal(t, "value", val)
			},
		},
		{
			name: "should handle special characters in keys",
			testFunc: func(t *testing.T) {
				cache := NewCache()
				specialKeys := []string{
					"key!@#$%^&*()",
					"key\nwith\nnewlines",
					"key\twith\ttabs",
					"key with spaces",
					"key/with/slashes",
				}
				for _, key := range specialKeys {
					cache.Set(key, "value")
					val, exists := cache.Get(key)
					assert.True(t, exists)
					assert.Equal(t, "value", val)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}