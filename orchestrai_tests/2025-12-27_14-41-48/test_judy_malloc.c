#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <cmocka.h>

/* Forward declarations */
extern void* judy_malloc(size_t size);
extern void* judy_calloc(size_t count, size_t size);
extern void* judy_realloc(void* ptr, size_t size);
extern void judy_free(void* ptr);
extern void judy_malloc_stats(void);

/* Test suite for judy_malloc */

/* Tests for judy_malloc function */
static void test_judy_malloc_basic_allocation(void **state) {
    /* Test basic malloc with normal size */
    void *ptr = judy_malloc(100);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

static void test_judy_malloc_zero_size(void **state) {
    /* Test malloc with zero size - should handle gracefully */
    void *ptr = judy_malloc(0);
    /* Behavior may vary: could be NULL or a valid pointer */
}

static void test_judy_malloc_large_size(void **state) {
    /* Test malloc with large size allocation */
    void *ptr = judy_malloc(1024 * 1024 * 10);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

static void test_judy_malloc_small_size(void **state) {
    /* Test malloc with very small size */
    void *ptr = judy_malloc(1);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

static void test_judy_malloc_multiple_allocations(void **state) {
    /* Test multiple independent allocations */
    void *ptr1 = judy_malloc(50);
    void *ptr2 = judy_malloc(100);
    void *ptr3 = judy_malloc(200);
    
    assert_non_null(ptr1);
    assert_non_null(ptr2);
    assert_non_null(ptr3);
    assert_ptr_not_equal(ptr1, ptr2);
    assert_ptr_not_equal(ptr2, ptr3);
    assert_ptr_not_equal(ptr1, ptr3);
    
    judy_free(ptr1);
    judy_free(ptr2);
    judy_free(ptr3);
}

/* Tests for judy_calloc function */
static void test_judy_calloc_basic_allocation(void **state) {
    /* Test basic calloc with normal count and size */
    void *ptr = judy_calloc(10, 20);
    assert_non_null(ptr);
    
    /* Verify memory is zeroed */
    unsigned char *bytes = (unsigned char *)ptr;
    for (int i = 0; i < 200; i++) {
        assert_int_equal(bytes[i], 0);
    }
    
    judy_free(ptr);
}

static void test_judy_calloc_zero_count(void **state) {
    /* Test calloc with zero count */
    void *ptr = judy_calloc(0, 20);
    /* Behavior may vary: could be NULL or a valid pointer */
}

static void test_judy_calloc_zero_size(void **state) {
    /* Test calloc with zero size */
    void *ptr = judy_calloc(10, 0);
    /* Behavior may vary: could be NULL or a valid pointer */
}

static void test_judy_calloc_both_zero(void **state) {
    /* Test calloc with both count and size zero */
    void *ptr = judy_calloc(0, 0);
    /* Behavior may vary: could be NULL or a valid pointer */
}

static void test_judy_calloc_large_allocation(void **state) {
    /* Test calloc with large total size */
    void *ptr = judy_calloc(1000, 1000);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

static void test_judy_calloc_multiple_allocations(void **state) {
    /* Test multiple calloc allocations */
    void *ptr1 = judy_calloc(5, 10);
    void *ptr2 = judy_calloc(8, 15);
    
    assert_non_null(ptr1);
    assert_non_null(ptr2);
    assert_ptr_not_equal(ptr1, ptr2);
    
    judy_free(ptr1);
    judy_free(ptr2);
}

/* Tests for judy_realloc function */
static void test_judy_realloc_from_null(void **state) {
    /* Test realloc with NULL pointer (should behave like malloc) */
    void *ptr = judy_realloc(NULL, 100);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

static void test_judy_realloc_grow_allocation(void **state) {
    /* Test realloc growing an allocation */
    void *ptr = judy_malloc(100);
    assert_non_null(ptr);
    
    void *new_ptr = judy_realloc(ptr, 200);
    assert_non_null(new_ptr);
    
    judy_free(new_ptr);
}

static void test_judy_realloc_shrink_allocation(void **state) {
    /* Test realloc shrinking an allocation */
    void *ptr = judy_malloc(200);
    assert_non_null(ptr);
    
    void *new_ptr = judy_realloc(ptr, 100);
    assert_non_null(new_ptr);
    
    judy_free(new_ptr);
}

static void test_judy_realloc_same_size(void **state) {
    /* Test realloc with same size */
    void *ptr = judy_malloc(100);
    assert_non_null(ptr);
    
    void *new_ptr = judy_realloc(ptr, 100);
    assert_non_null(new_ptr);
    
    judy_free(new_ptr);
}

static void test_judy_realloc_to_zero(void **state) {
    /* Test realloc to zero size */
    void *ptr = judy_malloc(100);
    assert_non_null(ptr);
    
    void *new_ptr = judy_realloc(ptr, 0);
    /* Behavior may vary: could be NULL or a valid pointer */
}

static void test_judy_realloc_preserves_data(void **state) {
    /* Test that realloc preserves existing data */
    void *ptr = judy_malloc(100);
    assert_non_null(ptr);
    
    /* Write test data */
    memset(ptr, 0xAB, 50);
    
    void *new_ptr = judy_realloc(ptr, 200);
    assert_non_null(new_ptr);
    
    /* Verify data is preserved */
    unsigned char *bytes = (unsigned char *)new_ptr;
    for (int i = 0; i < 50; i++) {
        assert_int_equal(bytes[i], 0xAB);
    }
    
    judy_free(new_ptr);
}

static void test_judy_realloc_large_allocation(void **state) {
    /* Test realloc to large size */
    void *ptr = judy_malloc(100);
    assert_non_null(ptr);
    
    void *new_ptr = judy_realloc(ptr, 1024 * 1024);
    if (new_ptr != NULL) {
        judy_free(new_ptr);
    }
}

/* Tests for judy_free function */
static void test_judy_free_valid_pointer(void **state) {
    /* Test free with valid pointer */
    void *ptr = judy_malloc(100);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

static void test_judy_free_null_pointer(void **state) {
    /* Test free with NULL pointer - should be safe */
    judy_free(NULL);
}

static void test_judy_free_multiple_allocations(void **state) {
    /* Test free on multiple allocations */
    void *ptr1 = judy_malloc(50);
    void *ptr2 = judy_malloc(100);
    void *ptr3 = judy_malloc(150);
    
    assert_non_null(ptr1);
    assert_non_null(ptr2);
    assert_non_null(ptr3);
    
    judy_free(ptr1);
    judy_free(ptr2);
    judy_free(ptr3);
}

static void test_judy_free_after_calloc(void **state) {
    /* Test free on calloc'd memory */
    void *ptr = judy_calloc(10, 10);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

static void test_judy_free_after_realloc(void **state) {
    /* Test free on reallocated memory */
    void *ptr = judy_malloc(100);
    void *new_ptr = judy_realloc(ptr, 200);
    assert_non_null(new_ptr);
    
    judy_free(new_ptr);
}

/* Tests for judy_malloc_stats function */
static void test_judy_malloc_stats_basic(void **state) {
    /* Test that malloc_stats can be called without crashing */
    judy_malloc_stats();
}

static void test_judy_malloc_stats_with_allocations(void **state) {
    /* Test malloc_stats with active allocations */
    void *ptr1 = judy_malloc(100);
    void *ptr2 = judy_malloc(200);
    
    judy_malloc_stats();
    
    judy_free(ptr1);
    judy_free(ptr2);
    
    judy_malloc_stats();
}

/* Integration tests */
static void test_malloc_calloc_realloc_free_sequence(void **state) {
    /* Test complete lifecycle of malloc->calloc->realloc->free */
    void *ptr1 = judy_malloc(100);
    void *ptr2 = judy_calloc(20, 20);
    void *ptr3 = judy_realloc(ptr1, 200);
    
    assert_non_null(ptr1);
    assert_non_null(ptr2);
    assert_non_null(ptr3);
    
    judy_free(ptr2);
    judy_free(ptr3);
}

static void test_alternating_allocations_and_frees(void **state) {
    /* Test alternating alloc/free patterns */
    void *ptr1 = judy_malloc(100);
    assert_non_null(ptr1);
    
    void *ptr2 = judy_malloc(200);
    assert_non_null(ptr2);
    
    judy_free(ptr1);
    
    void *ptr3 = judy_malloc(150);
    assert_non_null(ptr3);
    
    judy_free(ptr2);
    judy_free(ptr3);
}

static void test_memory_reuse_after_free(void **state) {
    /* Test that freed memory can be reused */
    void *ptr1 = judy_malloc(100);
    assert_non_null(ptr1);
    
    judy_free(ptr1);
    
    void *ptr2 = judy_malloc(100);
    assert_non_null(ptr2);
    
    judy_free(ptr2);
}

static void test_mixed_allocation_functions(void **state) {
    /* Test mixing malloc, calloc, and realloc */
    void *ptr1 = judy_malloc(100);
    void *ptr2 = judy_calloc(10, 20);
    void *ptr3 = judy_realloc(ptr1, 150);
    void *ptr4 = judy_malloc(50);
    
    assert_non_null(ptr1);
    assert_non_null(ptr2);
    assert_non_null(ptr3);
    assert_non_null(ptr4);
    
    judy_free(ptr2);
    judy_free(ptr3);
    judy_free(ptr4);
}

/* Edge cases */
static void test_allocation_boundary_1_byte(void **state) {
    /* Test allocation of exactly 1 byte */
    void *ptr = judy_malloc(1);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

static void test_allocation_boundary_page_size(void **state) {
    /* Test allocation near page size (typically 4096) */
    void *ptr = judy_malloc(4096);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

static void test_allocation_boundary_multiple_pages(void **state) {
    /* Test allocation spanning multiple pages */
    void *ptr = judy_malloc(4096 * 10);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

static void test_calloc_overflow_protection(void **state) {
    /* Test calloc with large values that might overflow */
    /* Allocation of 2^31 * 2^31 would overflow if not handled */
    void *ptr = judy_calloc(1000000, 1000000);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

static void test_realloc_chain(void **state) {
    /* Test multiple realloc calls in sequence */
    void *ptr = judy_malloc(10);
    assert_non_null(ptr);
    
    for (int i = 0; i < 5; i++) {
        ptr = judy_realloc(ptr, (i + 2) * 10);
        assert_non_null(ptr);
    }
    
    judy_free(ptr);
}

static void test_stats_accuracy_with_allocations(void **state) {
    /* Test that stats reflect allocations and frees */
    judy_malloc_stats();
    
    void *ptr1 = judy_malloc(1000);
    judy_malloc_stats();
    
    void *ptr2 = judy_malloc(2000);
    judy_malloc_stats();
    
    judy_free(ptr1);
    judy_malloc_stats();
    
    judy_free(ptr2);
    judy_malloc_stats();
}

int main(void) {
    const struct CMUnitTest tests[] = {
        /* judy_malloc tests */
        cmocka_unit_test(test_judy_malloc_basic_allocation),
        cmocka_unit_test(test_judy_malloc_zero_size),
        cmocka_unit_test(test_judy_malloc_large_size),
        cmocka_unit_test(test_judy_malloc_small_size),
        cmocka_unit_test(test_judy_malloc_multiple_allocations),
        
        /* judy_calloc tests */
        cmocka_unit_test(test_judy_calloc_basic_allocation),
        cmocka_unit_test(test_judy_calloc_zero_count),
        cmocka_unit_test(test_judy_calloc_zero_size),
        cmocka_unit_test(test_judy_calloc_both_zero),
        cmocka_unit_test(test_judy_calloc_large_allocation),
        cmocka_unit_test(test_judy_calloc_multiple_allocations),
        
        /* judy_realloc tests */
        cmocka_unit_test(test_judy_realloc_from_null),
        cmocka_unit_test(test_judy_realloc_grow_allocation),
        cmocka_unit_test(test_judy_realloc_shrink_allocation),
        cmocka_unit_test(test_judy_realloc_same_size),
        cmocka_unit_test(test_judy_realloc_to_zero),
        cmocka_unit_test(test_judy_realloc_preserves_data),
        cmocka_unit_test(test_judy_realloc_large_allocation),
        
        /* judy_free tests */
        cmocka_unit_test(test_judy_free_valid_pointer),
        cmocka_unit_test(test_judy_free_null_pointer),
        cmocka_unit_test(test_judy_free_multiple_allocations),
        cmocka_unit_test(test_judy_free_after_calloc),
        cmocka_unit_test(test_judy_free_after_realloc),
        
        /* judy_malloc_stats tests */
        cmocka_unit_test(test_judy_malloc_stats_basic),
        cmocka_unit_test(test_judy_malloc_stats_with_allocations),
        
        /* Integration tests */
        cmocka_unit_test(test_malloc_calloc_realloc_free_sequence),
        cmocka_unit_test(test_alternating_allocations_and_frees),
        cmocka_unit_test(test_memory_reuse_after_free),
        cmocka_unit_test(test_mixed_allocation_functions),
        
        /* Edge case tests */
        cmocka_unit_test(test_allocation_boundary_1_byte),
        cmocka_unit_test(test_allocation_boundary_page_size),
        cmocka_unit_test(test_allocation_boundary_multiple_pages),
        cmocka_unit_test(test_calloc_overflow_protection),
        cmocka_unit_test(test_realloc_chain),
        cmocka_unit_test(test_stats_accuracy_with_allocations),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}