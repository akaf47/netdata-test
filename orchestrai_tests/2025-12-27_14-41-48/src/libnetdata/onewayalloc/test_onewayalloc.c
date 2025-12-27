#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <assert.h>
#include <limits.h>
#include "onewayalloc.h"

// Test counters and state
static int tests_run = 0;
static int tests_passed = 0;
static int tests_failed = 0;

#define TEST_START(name) do { \
    tests_run++; \
    printf("\n[TEST %d] %s... ", tests_run, name); \
    fflush(stdout); \
} while(0)

#define TEST_PASS() do { \
    tests_passed++; \
    printf("✓ PASS\n"); \
} while(0)

#define TEST_FAIL(msg) do { \
    tests_failed++; \
    printf("✗ FAIL: %s\n", msg); \
} while(0)

#define ASSERT_NOT_NULL(ptr, msg) do { \
    if ((ptr) == NULL) { \
        TEST_FAIL(msg); \
        return; \
    } \
} while(0)

#define ASSERT_NULL(ptr, msg) do { \
    if ((ptr) != NULL) { \
        TEST_FAIL(msg); \
        return; \
    } \
} while(0)

#define ASSERT_TRUE(cond, msg) do { \
    if (!(cond)) { \
        TEST_FAIL(msg); \
        return; \
    } \
} while(0)

#define ASSERT_EQUAL(actual, expected, msg) do { \
    if ((actual) != (expected)) { \
        TEST_FAIL(msg); \
        return; \
    } \
} while(0)

#define ASSERT_STR_EQUAL(actual, expected, msg) do { \
    if (strcmp((actual), (expected)) != 0) { \
        TEST_FAIL(msg); \
        return; \
    } \
} while(0)

// ============================================================================
// Test: onewayalloc_new with default size
// ============================================================================
static void test_onewayalloc_new_default_size(void) {
    TEST_START("onewayalloc_new with default size");
    
    ONEWAYALLOC *owa = onewayalloc_new(0);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc with size 0");
    
    // Verify basic structure
    ASSERT_NOT_NULL(owa, "Allocation returned NULL");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_new with custom size
// ============================================================================
static void test_onewayalloc_new_custom_size(void) {
    TEST_START("onewayalloc_new with custom size");
    
    size_t custom_size = 16 * 1024;
    ONEWAYALLOC *owa = onewayalloc_new(custom_size);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc with custom size");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_new with very small size
// ============================================================================
static void test_onewayalloc_new_small_size(void) {
    TEST_START("onewayalloc_new with very small size");
    
    ONEWAYALLOC *owa = onewayalloc_new(1);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc with size 1");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_new with large size
// ============================================================================
static void test_onewayalloc_new_large_size(void) {
    TEST_START("onewayalloc_new with large size");
    
    size_t large_size = 100 * 1024 * 1024;  // 100MB
    ONEWAYALLOC *owa = onewayalloc_new(large_size);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc with large size");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_malloc basic allocation
// ============================================================================
static void test_onewayalloc_malloc_basic(void) {
    TEST_START("onewayalloc_malloc basic allocation");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_malloc(owa, 100);
    ASSERT_NOT_NULL(ptr, "Failed to allocate 100 bytes");
    
    // Verify pointer is properly aligned
    ASSERT_TRUE((uintptr_t)ptr % sizeof(void*) == 0, "Pointer not properly aligned");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_malloc zero bytes
// ============================================================================
static void test_onewayalloc_malloc_zero_bytes(void) {
    TEST_START("onewayalloc_malloc zero bytes");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_malloc(owa, 0);
    // Zero allocation may return non-NULL or NULL depending on implementation
    // Both are acceptable in C
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_malloc single byte
// ============================================================================
static void test_onewayalloc_malloc_single_byte(void) {
    TEST_START("onewayalloc_malloc single byte");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_malloc(owa, 1);
    ASSERT_NOT_NULL(ptr, "Failed to allocate 1 byte");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_malloc multiple sequential allocations
// ============================================================================
static void test_onewayalloc_malloc_sequential(void) {
    TEST_START("onewayalloc_malloc multiple sequential allocations");
    
    ONEWAYALLOC *owa = onewayalloc_new(10240);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr1 = onewayalloc_malloc(owa, 256);
    ASSERT_NOT_NULL(ptr1, "First allocation failed");
    
    void *ptr2 = onewayalloc_malloc(owa, 256);
    ASSERT_NOT_NULL(ptr2, "Second allocation failed");
    
    void *ptr3 = onewayalloc_malloc(owa, 256);
    ASSERT_NOT_NULL(ptr3, "Third allocation failed");
    
    // Verify pointers are different
    ASSERT_TRUE(ptr1 != ptr2, "Pointers 1 and 2 should differ");
    ASSERT_TRUE(ptr2 != ptr3, "Pointers 2 and 3 should differ");
    ASSERT_TRUE(ptr1 != ptr3, "Pointers 1 and 3 should differ");
    
    // Verify sequential ordering
    ASSERT_TRUE((uintptr_t)ptr1 < (uintptr_t)ptr2, "First pointer should be less than second");
    ASSERT_TRUE((uintptr_t)ptr2 < (uintptr_t)ptr3, "Second pointer should be less than third");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_malloc write and read data
// ============================================================================
static void test_onewayalloc_malloc_write_read(void) {
    TEST_START("onewayalloc_malloc write and read data");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_malloc(owa, 100);
    ASSERT_NOT_NULL(ptr, "Failed to allocate");
    
    // Write test data
    char test_data[] = "Hello, World!";
    memcpy(ptr, test_data, strlen(test_data) + 1);
    
    // Verify data integrity
    ASSERT_STR_EQUAL((char*)ptr, test_data, "Data integrity check failed");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_malloc exhausts pool
// ============================================================================
static void test_onewayalloc_malloc_exhausts_pool(void) {
    TEST_START("onewayalloc_malloc exhausts pool");
    
    size_t pool_size = 1024;
    ONEWAYALLOC *owa = onewayalloc_new(pool_size);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    // Allocate close to pool size
    void *ptr1 = onewayalloc_malloc(owa, 512);
    ASSERT_NOT_NULL(ptr1, "First allocation failed");
    
    void *ptr2 = onewayalloc_malloc(owa, 512);
    // Second allocation should succeed or fail gracefully
    // Both behaviors are acceptable for one-way allocator
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_malloc with growing allocator (auto-expand)
// ============================================================================
static void test_onewayalloc_malloc_auto_expand(void) {
    TEST_START("onewayalloc_malloc auto-expand behavior");
    
    ONEWAYALLOC *owa = onewayalloc_new(100);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    // Allocate more than initial pool
    void *ptr1 = onewayalloc_malloc(owa, 200);
    ASSERT_NOT_NULL(ptr1, "Allocation should succeed even if > initial size");
    
    void *ptr2 = onewayalloc_malloc(owa, 200);
    ASSERT_NOT_NULL(ptr2, "Second large allocation should succeed");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_calloc basic allocation
// ============================================================================
static void test_onewayalloc_calloc_basic(void) {
    TEST_START("onewayalloc_calloc basic allocation");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_calloc(owa, 10, 10);
    ASSERT_NOT_NULL(ptr, "Failed to allocate with calloc");
    
    // Verify memory is zeroed
    char *byte_ptr = (char*)ptr;
    int all_zero = 1;
    for (int i = 0; i < 100; i++) {
        if (byte_ptr[i] != 0) {
            all_zero = 0;
            break;
        }
    }
    ASSERT_TRUE(all_zero, "Allocated memory should be zeroed");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_calloc with zero count
// ============================================================================
static void test_onewayalloc_calloc_zero_count(void) {
    TEST_START("onewayalloc_calloc with zero count");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_calloc(owa, 0, 10);
    // Zero allocation may return NULL or non-NULL
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_calloc with zero size
// ============================================================================
static void test_onewayalloc_calloc_zero_size(void) {
    TEST_START("onewayalloc_calloc with zero size");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_calloc(owa, 10, 0);
    // Zero allocation may return NULL or non-NULL
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_calloc both zero
// ============================================================================
static void test_onewayalloc_calloc_both_zero(void) {
    TEST_START("onewayalloc_calloc with both count and size zero");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_calloc(owa, 0, 0);
    // Zero allocation may return NULL or non-NULL
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_calloc large allocation
// ============================================================================
static void test_onewayalloc_calloc_large(void) {
    TEST_START("onewayalloc_calloc large allocation");
    
    ONEWAYALLOC *owa = onewayalloc_new(10240);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_calloc(owa, 100, 100);
    ASSERT_NOT_NULL(ptr, "Failed to allocate large block");
    
    // Verify zeroed
    uint32_t *int_ptr = (uint32_t*)ptr;
    int all_zero = 1;
    for (int i = 0; i < 2500; i++) {
        if (int_ptr[i] != 0) {
            all_zero = 0;
            break;
        }
    }
    ASSERT_TRUE(all_zero, "Large allocation should be zeroed");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_realloc basic
// ============================================================================
static void test_onewayalloc_realloc_basic(void) {
    TEST_START("onewayalloc_realloc basic");
    
    ONEWAYALLOC *owa = onewayalloc_new(2048);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr1 = onewayalloc_malloc(owa, 100);
    ASSERT_NOT_NULL(ptr1, "First allocation failed");
    
    // Fill with data
    memset(ptr1, 0xAA, 100);
    
    // Realloc to larger size
    void *ptr2 = onewayalloc_realloc(owa, ptr1, 100, 200);
    ASSERT_NOT_NULL(ptr2, "Realloc failed");
    
    // Verify original data preserved
    char *byte_ptr = (char*)ptr2;
    int data_preserved = 1;
    for (int i = 0; i < 100; i++) {
        if (byte_ptr[i] != 0xAA) {
            data_preserved = 0;
            break;
        }
    }
    ASSERT_TRUE(data_preserved, "Original data should be preserved after realloc");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_realloc zero old size
// ============================================================================
static void test_onewayalloc_realloc_zero_old_size(void) {
    TEST_START("onewayalloc_realloc with zero old size");
    
    ONEWAYALLOC *owa = onewayalloc_new(2048);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr1 = onewayalloc_malloc(owa, 100);
    ASSERT_NOT_NULL(ptr1, "First allocation failed");
    
    // Realloc with zero old size
    void *ptr2 = onewayalloc_realloc(owa, ptr1, 0, 200);
    ASSERT_NOT_NULL(ptr2, "Realloc with zero old size failed");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_realloc zero new size
// ============================================================================
static void test_onewayalloc_realloc_zero_new_size(void) {
    TEST_START("onewayalloc_realloc with zero new size");
    
    ONEWAYALLOC *owa = onewayalloc_new(2048);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr1 = onewayalloc_malloc(owa, 100);
    ASSERT_NOT_NULL(ptr1, "First allocation failed");
    
    // Realloc to zero size
    void *ptr2 = onewayalloc_realloc(owa, ptr1, 100, 0);
    // Result may be NULL or non-NULL depending on implementation
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_realloc smaller size
// ============================================================================
static void test_onewayalloc_realloc_smaller(void) {
    TEST_START("onewayalloc_realloc to smaller size");
    
    ONEWAYALLOC *owa = onewayalloc_new(2048);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr1 = onewayalloc_malloc(owa, 200);
    ASSERT_NOT_NULL(ptr1, "First allocation failed");
    
    memset(ptr1, 0xBB, 200);
    
    // Realloc to smaller size
    void *ptr2 = onewayalloc_realloc(owa, ptr1, 200, 100);
    ASSERT_NOT_NULL(ptr2, "Realloc to smaller size failed");
    
    // Verify data preserved (at least up to new size)
    char *byte_ptr = (char*)ptr2;
    int data_preserved = 1;
    for (int i = 0; i < 100; i++) {
        if (byte_ptr[i] != 0xBB) {
            data_preserved = 0;
            break;
        }
    }
    ASSERT_TRUE(data_preserved, "Data should be preserved in smaller realloc");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_realloc same size
// ============================================================================
static void test_onewayalloc_realloc_same_size(void) {
    TEST_START("onewayalloc_realloc same size");
    
    ONEWAYALLOC *owa = onewayalloc_new(2048);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr1 = onewayalloc_malloc(owa, 100);
    ASSERT_NOT_NULL(ptr1, "First allocation failed");
    
    memset(ptr1, 0xCC, 100);
    
    // Realloc to same size
    void *ptr2 = onewayalloc_realloc(owa, ptr1, 100, 100);
    ASSERT_NOT_NULL(ptr2, "Realloc same size failed");
    
    // Data may be preserved or not depending on implementation
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_strdupz basic
// ============================================================================
static void test_onewayalloc_strdupz_basic(void) {
    TEST_START("onewayalloc_strdupz basic");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    const char *original = "Hello, World!";
    char *dup = onewayalloc_strdupz(owa, original);
    ASSERT_NOT_NULL(dup, "strdupz failed");
    
    ASSERT_STR_EQUAL(dup, original, "Duplicated string should match original");
    ASSERT_TRUE(dup != original, "Should be a different pointer");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_strdupz empty string
// ============================================================================
static void test_onewayalloc_strdupz_empty(void) {
    TEST_START("onewayalloc_strdupz empty string");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    char *dup = onewayalloc_strdupz(owa, "");
    ASSERT_NOT_NULL(dup, "strdupz empty string failed");
    
    ASSERT_STR_EQUAL(dup, "", "Empty string should duplicate correctly");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_strdupz single char
// ============================================================================
static void test_onewayalloc_strdupz_single_char(void) {
    TEST_START("onewayalloc_strdupz single character");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    char *dup = onewayalloc_strdupz(owa, "A");
    ASSERT_NOT_NULL(dup, "strdupz single char failed");
    
    ASSERT_STR_EQUAL(dup, "A", "Single character should duplicate correctly");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_strdupz long string
// ============================================================================
static void test_onewayalloc_strdupz_long(void) {
    TEST_START("onewayalloc_strdupz long string");
    
    ONEWAYALLOC *owa = onewayalloc_new(10240);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    char long_str[5001];
    memset(long_str, 'x', 5000);
    long_str[5000] = '\0';
    
    char *dup = onewayalloc_strdupz(owa, long_str);
    ASSERT_NOT_NULL(dup, "strdupz long string failed");
    
    ASSERT_STR_EQUAL(dup, long_str, "Long string should duplicate correctly");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_strdupz multiple strings
// ============================================================================
static void test_onewayalloc_strdupz_multiple(void) {
    TEST_START("onewayalloc_strdupz multiple strings");
    
    ONEWAYALLOC *owa = onewayalloc_new(10240);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    char *dup1 = onewayalloc_strdupz(owa, "First");
    char *dup2 = onewayalloc_strdupz(owa, "Second");
    char *dup3 = onewayalloc_strdupz(owa, "Third");
    
    ASSERT_NOT_NULL(dup1, "First string dup failed");
    ASSERT_NOT_NULL(dup2, "Second string dup failed");
    ASSERT_NOT_NULL(dup3, "Third string dup failed");
    
    ASSERT_STR_EQUAL(dup1, "First", "First string should match");
    ASSERT_STR_EQUAL(dup2, "Second", "Second string should match");
    ASSERT_STR_EQUAL(dup3, "Third", "Third string should match");
    
    ASSERT_TRUE(dup1 != dup2, "Different strings should have different pointers");
    ASSERT_TRUE(dup2 != dup3, "Different strings should have different pointers");
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_destroy null pointer
// ============================================================================
static void test_onewayalloc_destroy_null(void) {
    TEST_START("onewayalloc_destroy with null pointer");
    
    // Should not crash
    onewayalloc_destroy(NULL);
    
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_destroy valid pointer
// ============================================================================
static void test_onewayalloc_destroy_valid(void) {
    TEST_START("onewayalloc_destroy with valid pointer");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    // Destroy should not crash and should free resources
    onewayalloc_destroy(owa);
    
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_destroy after allocations
// ============================================================================
static void test_onewayalloc_destroy_after_allocations(void) {
    TEST_START("onewayalloc_destroy after multiple allocations");
    
    ONEWAYALLOC *owa = onewayalloc_new(2048);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    // Make some allocations
    void *ptr1 = onewayalloc_malloc(owa, 100);
    void *ptr2 = onewayalloc_malloc(owa, 200);
    char *str = onewayalloc_strdupz(owa, "test string");
    
    // All should succeed
    ASSERT_NOT_NULL(ptr1, "First malloc failed");
    ASSERT_NOT_NULL(ptr2, "Second malloc failed");
    ASSERT_NOT_NULL(str, "strdupz failed");
    
    // Destroy should clean up all allocations
    onewayalloc_destroy(owa);
    
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_freez with valid buffer
// ============================================================================
static void test_onewayalloc_freez_valid(void) {
    TEST_START("onewayalloc_freez with valid buffer");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    void *ptr = onewayalloc_malloc(owa, 100);
    ASSERT_NOT_NULL(ptr, "malloc failed");
    
    // freez should work
    void *result = onewayalloc_freez(owa, ptr);
    // Result may be NULL or pointer depending on implementation
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc_freez with null pointer
// ============================================================================
static void test_onewayalloc_freez_null(void) {
    TEST_START("onewayalloc_freez with null pointer");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    // freez on NULL should handle gracefully
    void *result = onewayalloc_freez(owa, NULL);
    // Should not crash
    
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc stress test - many allocations
// ============================================================================
static void test_onewayalloc_stress_many_allocations(void) {
    TEST_START("onewayalloc stress test - many allocations");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024 * 1024);
    ASSERT_NOT_NULL(owa, "Failed to create onewayalloc");
    
    // Allocate many small blocks
    int count = 1000;
    void **ptrs = malloc(count * sizeof(void*));
    ASSERT_NOT_NULL(ptrs, "Failed to allocate pointer array");
    
    for (int i = 0; i < count; i++) {
        ptrs[i] = onewayalloc_malloc(owa, 16);
        // Some may fail if pool exhausted, that's okay
    }
    
    // Verify at least first few succeeded
    int success_count = 0;
    for (int i = 0; i < count; i++) {
        if (ptrs[i] != NULL) success_count++;
    }
    ASSERT_TRUE(success_count > 0, "At least some allocations should succeed");
    
    free(ptrs);
    onewayalloc_destroy(owa);
    TEST_PASS();
}

// ============================================================================
// Test: onewayalloc stress test - large and small mixed
// ============================================================================
static void test_onewayalloc_stress_mixed(void) {
    TEST_START("onewayalloc stress test - large and small mixed");
    
    ONEWAYALLOC *owa = onewayalloc_new(1024 * 1024);
    ASSERT_NOT_NULL(owa, "Failed to create on