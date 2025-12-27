package filecheck

import (
	"testing"
	"time"
)

type mockCacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// TestNewCache tests cache initialization
func TestNewCache(t *testing.T) {
	cache := newCache()
	if cache == nil {
		t.Fatal("expected cache to be initialized, got nil")
	}
}

// TestCacheSet tests setting values in cache
func TestCacheSet(t *testing.T) {
	cache := newCache()
	key := "test_key"
	value := "test_value"
	ttl := 1 * time.Minute

	cache.Set(key, value, ttl)

	if entry, exists := cache.cache[key]; !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	} else {
		if entry.value != value {
			t.Errorf("expected value %v, got %v", value, entry.value)
		}
	}
}

// TestCacheSetWithZeroTTL tests setting cache with zero TTL
func TestCacheSetWithZeroTTL(t *testing.T) {
	cache := newCache()
	key := "test_key"
	value := "test_value"

	cache.Set(key, value, 0)

	if entry, exists := cache.cache[key]; !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	} else {
		if entry.expiresAt != (time.Time{}) {
			t.Errorf("expected zero expiration time, got %v", entry.expiresAt)
		}
	}
}

// TestCacheSetWithNegativeTTL tests setting cache with negative TTL
func TestCacheSetWithNegativeTTL(t *testing.T) {
	cache := newCache()
	key := "test_key"
	value := "test_value"

	cache.Set(key, value, -1*time.Minute)

	if _, exists := cache.cache[key]; !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	}
}

// TestCacheSetWithLongTTL tests setting cache with very long TTL
func TestCacheSetWithLongTTL(t *testing.T) {
	cache := newCache()
	key := "test_key"
	value := "test_value"
	ttl := 1000000 * time.Hour

	cache.Set(key, value, ttl)

	if _, exists := cache.cache[key]; !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	}
}

// TestCacheGet tests retrieving values from cache
func TestCacheGet(t *testing.T) {
	cache := newCache()
	key := "test_key"
	expectedValue := "test_value"
	ttl := 1 * time.Minute

	cache.Set(key, expectedValue, ttl)
	value, exists := cache.Get(key)

	if !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	}
	if value != expectedValue {
		t.Errorf("expected value %v, got %v", expectedValue, value)
	}
}

// TestCacheGetNonExistentKey tests getting non-existent key
func TestCacheGetNonExistentKey(t *testing.T) {
	cache := newCache()

	value, exists := cache.Get("non_existent_key")

	if exists {
		t.Errorf("expected key to not exist, got value %v", value)
	}
	if value != nil {
		t.Errorf("expected nil value for non-existent key, got %v", value)
	}
}

// TestCacheGetExpiredEntry tests getting expired cache entry
func TestCacheGetExpiredEntry(t *testing.T) {
	cache := newCache()
	key := "expired_key"
	value := "test_value"

	// Set entry with expiration in the past
	cache.cache[key] = &cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(-1 * time.Minute),
	}

	retrieved, exists := cache.Get(key)

	if exists {
		t.Errorf("expected expired entry to not exist, got value %v", retrieved)
	}
	if retrieved != nil {
		t.Errorf("expected nil for expired entry, got %v", retrieved)
	}
}

// TestCacheDelete tests deleting cache entries
func TestCacheDelete(t *testing.T) {
	cache := newCache()
	key := "test_key"
	value := "test_value"
	ttl := 1 * time.Minute

	cache.Set(key, value, ttl)
	cache.Delete(key)

	if _, exists := cache.cache[key]; exists {
		t.Fatalf("expected key %s to be deleted", key)
	}
}

// TestCacheDeleteNonExistentKey tests deleting non-existent key
func TestCacheDeleteNonExistentKey(t *testing.T) {
	cache := newCache()

	// Should not panic
	cache.Delete("non_existent_key")
}

// TestCacheDeleteMultipleKeys tests deleting multiple keys
func TestCacheDeleteMultipleKeys(t *testing.T) {
	cache := newCache()
	keys := []string{"key1", "key2", "key3"}

	for _, key := range keys {
		cache.Set(key, "value", 1*time.Minute)
	}

	for _, key := range keys {
		cache.Delete(key)
	}

	for _, key := range keys {
		if _, exists := cache.cache[key]; exists {
			t.Fatalf("expected key %s to be deleted", key)
		}
	}
}

// TestCacheFlush tests clearing all cache entries
func TestCacheFlush(t *testing.T) {
	cache := newCache()

	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 1*time.Minute)
	cache.Set("key3", "value3", 1*time.Minute)

	cache.Flush()

	if len(cache.cache) != 0 {
		t.Errorf("expected cache to be empty after flush, got %d entries", len(cache.cache))
	}
}

// TestCacheFlushEmptyCache tests flushing already empty cache
func TestCacheFlushEmptyCache(t *testing.T) {
	cache := newCache()

	// Should not panic
	cache.Flush()

	if len(cache.cache) != 0 {
		t.Errorf("expected cache to remain empty, got %d entries", len(cache.cache))
	}
}

// TestCacheClear tests the Clear method if it exists
func TestCacheClear(t *testing.T) {
	cache := newCache()

	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 1*time.Minute)

	cache.Clear()

	if len(cache.cache) != 0 {
		t.Errorf("expected cache to be cleared, got %d entries", len(cache.cache))
	}
}

// TestCacheSetWithNilValue tests setting cache with nil value
func TestCacheSetWithNilValue(t *testing.T) {
	cache := newCache()
	key := "nil_key"

	cache.Set(key, nil, 1*time.Minute)

	value, exists := cache.Get(key)

	if !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	}
	if value != nil {
		t.Errorf("expected nil value, got %v", value)
	}
}

// TestCacheSetWithComplexValue tests setting cache with complex struct
func TestCacheSetWithComplexValue(t *testing.T) {
	cache := newCache()
	key := "complex_key"
	complexValue := map[string]interface{}{
		"name": "test",
		"age":  30,
		"tags": []string{"a", "b", "c"},
	}

	cache.Set(key, complexValue, 1*time.Minute)

	value, exists := cache.Get(key)

	if !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	}
	if value != complexValue {
		t.Errorf("expected value %v, got %v", complexValue, value)
	}
}

// TestCacheSetOverwrite tests overwriting existing cache entry
func TestCacheSetOverwrite(t *testing.T) {
	cache := newCache()
	key := "test_key"
	oldValue := "old_value"
	newValue := "new_value"
	ttl := 1 * time.Minute

	cache.Set(key, oldValue, ttl)
	cache.Set(key, newValue, ttl)

	value, exists := cache.Get(key)

	if !exists {
		t.Fatalf("expected key %s to exist in cache", key)
	}
	if value != newValue {
		t.Errorf("expected value %v, got %v", newValue, value)
	}
}

// TestCacheSetOverwriteWithDifferentTTL tests overwriting with different TTL
func TestCacheSetOverwriteWithDifferentTTL(t *testing.T) {
	cache := newCache()
	key := "test_key"
	value := "test_value"

	cache.Set(key, value, 1*time.Minute)
	cache.Set(key, value, 10*time.Minute)

	if entry, exists := cache.cache[key]; !exists {
		t.Fatalf("expected key %s to exist", key)
	} else {
		oldExpiry := entry.expiresAt
		cache.Set(key, value, 20*time.Minute)
		if cache.cache[key].expiresAt.Equal(oldExpiry) {
			t.Errorf("expected TTL to be updated")
		}
	}
}

// TestCacheMultipleConcurrentOperations tests multiple operations
func TestCacheMultipleConcurrentOperations(t *testing.T) {
	cache := newCache()

	// Set multiple entries
	for i := 0; i < 100; i++ {
		key := "key_" + string(rune(i))
		cache.Set(key, i, 1*time.Minute)
	}

	// Verify all entries exist
	for i := 0; i < 100; i++ {
		key := "key_" + string(rune(i))
		if _, exists := cache.Get(key); !exists {
			t.Errorf("expected key %s to exist", key)
		}
	}

	// Delete half of them
	for i := 0; i < 50; i++ {
		key := "key_" + string(rune(i))
		cache.Delete(key)
	}

	// Verify deletions
	for i := 0; i < 50; i++ {
		key := "key_" + string(rune(i))
		if _, exists := cache.Get(key); exists {
			t.Errorf("expected key %s to be deleted", key)
		}
	}

	// Verify remaining entries
	for i := 50; i < 100; i++ {
		key := "key_" + string(rune(i))
		if _, exists := cache.Get(key); !exists {
			t.Errorf("expected key %s to exist", key)
		}
	}
}

// TestCacheEmptyKeyHandling tests behavior with empty string key
func TestCacheEmptyKeyHandling(t *testing.T) {
	cache := newCache()
	emptyKey := ""
	value := "test_value"

	cache.Set(emptyKey, value, 1*time.Minute)

	retrieved, exists := cache.Get(emptyKey)

	if !exists {
		t.Fatalf("expected empty key to be stored")
	}
	if retrieved != value {
		t.Errorf("expected value %v, got %v", value, retrieved)
	}
}

// TestCacheSpecialCharacterKeys tests keys with special characters
func TestCacheSpecialCharacterKeys(t *testing.T) {
	cache := newCache()
	specialKeys := []string{
		"key:with:colons",
		"key/with/slashes",
		"key\\with\\backslashes",
		"key with spaces",
		"key\twith\ttabs",
		"key\nwith\nnewlines",
		"key-with-dashes",
		"key_with_underscores",
		"key.with.dots",
	}

	for _, key := range specialKeys {
		cache.Set(key, key, 1*time.Minute)
	}

	for _, key := range specialKeys {
		value, exists := cache.Get(key)
		if !exists {
			t.Errorf("expected key %s to exist", key)
		}
		if value != key {
			t.Errorf("expected value %s, got %v", key, value)
		}
	}
}

// TestCacheNumericKeys tests numeric string keys
func TestCacheNumericKeys(t *testing.T) {
	cache := newCache()
	numericKeys := []string{"0", "1", "-1", "999999", "3.14", "0x10"}

	for _, key := range numericKeys {
		cache.Set(key, key, 1*time.Minute)
	}

	for _, key := range numericKeys {
		value, exists := cache.Get(key)
		if !exists {
			t.Errorf("expected key %s to exist", key)
		}
		if value != key {
			t.Errorf("expected value %s, got %v", key, value)
		}
	}
}

// TestCacheUnicodeKeys tests unicode characters in keys
func TestCacheUnicodeKeys(t *testing.T) {
	cache := newCache()
	unicodeKeys := []string{"🔑key", "キー", "مفتاح", "klíč"}

	for _, key := range unicodeKeys {
		cache.Set(key, key, 1*time.Minute)
	}

	for _, key := range unicodeKeys {
		value, exists := cache.Get(key)
		if !exists {
			t.Errorf("expected key %s to exist", key)
		}
		if value != key {
			t.Errorf("expected value %s, got %v", key, value)
		}
	}
}

// TestCacheBoundaryConditions tests boundary conditions
func TestCacheBoundaryConditions(t *testing.T) {
	cache := newCache()

	// Test with very long key
	longKey := ""
	for i := 0; i < 10000; i++ {
		longKey += "a"
	}

	cache.Set(longKey, "value", 1*time.Minute)

	if _, exists := cache.Get(longKey); !exists {
		t.Errorf("expected long key to be stored")
	}
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

type cache struct {
	cache map[string]*cacheEntry
}

func newCache() *cache {
	return &cache{
		cache: make(map[string]*cacheEntry),
	}
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	expiresAt := time.Now().Add(ttl)
	c.cache[key] = &cacheEntry{
		value:     value,
		expiresAt: expiresAt,
	}
}

func (c *cache) Get(key string) (interface{}, bool) {
	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	if entry.expiresAt != (time.Time{}) && time.Now().After(entry.expiresAt) {
		delete(c.cache, key)
		return nil, false
	}

	return entry.value, true
}

func (c *cache) Delete(key string) {
	delete(c.cache, key)
}

func (c *cache) Flush() {
	c.cache = make(map[string]*cacheEntry)
}

func (c *cache) Clear() {
	for key := range c.cache {
		delete(c.cache, key)
	}
}