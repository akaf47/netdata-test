#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <cmocka.h>

/* Test suite for judy-malloc.h header file definitions */

/* Tests for macro definitions and inline functions in header */

/* Test that judy_malloc macro/function is available */
static void test_judy_malloc_header_available(void **state) {
    /* Verify judy_malloc is accessible */
    void *ptr = judy_malloc(100);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

/* Test that judy_calloc macro/function is available */
static void test_judy_calloc_header_available(void **state) {
    /* Verify judy_calloc is accessible */
    void *ptr = judy_calloc(10, 10);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

/* Test that judy_realloc macro/function is available */
static void test_judy_realloc_header_available(void **state) {
    /* Verify judy_realloc is accessible */
    void *ptr = judy_malloc(100);
    if (ptr != NULL) {
        void *new_ptr = judy_realloc(ptr, 200);
        if (new_ptr != NULL) {
            judy_free(new_ptr);
        }
    }
}

/* Test that judy_free macro/function is available */
static void test_judy_free_header_available(void **state) {
    /* Verify judy_free is accessible */
    void *ptr = judy_malloc(100);
    if (ptr != NULL) {
        judy_free(ptr);
    }
}

/* Test that judy_malloc_stats macro/function is available */
static void test_judy_malloc_stats_header_available(void **state) {
    /* Verify judy_malloc_stats is accessible */
    judy_malloc_stats();
}

/* Test proper definition of functions */
static void test_judy_malloc_returns_void_pointer(void **state) {
    void *ptr = judy_malloc(50);
    assert_non_null(ptr);
    judy_free(ptr);
}

static void test_judy_calloc_returns_void_pointer(void **state) {
    void *ptr = judy_calloc(5, 10);
    assert_non_null(ptr);
    judy_free(ptr);
}

static void test_judy_realloc_returns_void_pointer(void **state) {
    void *ptr = judy_malloc(50);
    assert_non_null(ptr);
    
    void *new_ptr = judy_realloc(ptr, 100);
    assert_non_null(new_ptr);
    judy_free(new_ptr);
}

/* Test function parameter types */
static void test_judy_malloc_accepts_size_t(void **state) {
    size_t size = 100;
    void *ptr = judy_malloc(size);
    assert_non_null(ptr);
    judy_free(ptr);
}

static void test_judy_calloc_accepts_size_t_parameters(void **state) {
    size_t count = 10;
    size_t size = 20;
    void *ptr = judy_calloc(count, size);
    assert_non_null(ptr);
    judy_free(ptr);
}

static void test_judy_realloc_accepts_void_and_size_t(void **state) {
    void *ptr = judy_malloc(50);
    assert_non_null(ptr);
    
    size_t new_size = 100;
    void *new_ptr = judy_realloc(ptr, new_size);
    assert_non_null(new_ptr);
    judy_free(new_ptr);
}

static void test_judy_free_accepts_void_pointer(void **state) {
    void *ptr = judy_malloc(50);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

/* Test type conversions and casts in headers */
static void test_type_safety_malloc_return(void **state) {
    unsigned char *ptr = (unsigned char *)judy_malloc(10);
    assert_non_null(ptr);
    
    *ptr = 0xFF;
    assert_int_equal(*ptr, 0xFF);
    
    judy_free(ptr);
}

static void test_type_safety_calloc_return(void **state) {
    int *ptr = (int *)judy_calloc(5, sizeof(int));
    assert_non_null(ptr);
    
    ptr[0] = 42;
    assert_int_equal(ptr[0], 42);
    
    judy_free(ptr);
}

static void test_type_safety_realloc_return(void **state) {
    char *ptr = (char *)judy_malloc(10);
    assert_non_null(ptr);
    
    ptr[0] = 'A';
    assert_int_equal(ptr[0], 'A');
    
    char *new_ptr = (char *)judy_realloc(ptr, 20);
    assert_non_null(new_ptr);
    
    assert_int_equal(new_ptr[0], 'A');
    
    judy_free(new_ptr);
}

/* Test compound operations */
static void test_sequential_allocations_from_header(void **state) {
    void *ptr1 = judy_malloc(100);
    void *ptr2 = judy_calloc(10, 20);
    void *ptr3 = judy_realloc(ptr1, 200);
    
    assert_non_null(ptr1);
    assert_non_null(ptr2);
    assert_non_null(ptr3);
    assert_ptr_not_equal(ptr1, ptr2);
    
    judy_free(ptr2);
    judy_free(ptr3);
}

static void test_array_of_allocations(void **state) {
    const int NUM_ALLOCS = 5;
    void *ptrs[NUM_ALLOCS];
    
    for (int i = 0; i < NUM_ALLOCS; i++) {
        ptrs[i] = judy_malloc((i + 1) * 100);
        assert_non_null(ptrs[i]);
    }
    
    for (int i = 0; i < NUM_ALLOCS; i++) {
        judy_free(ptrs[i]);
    }
}

static void test_nested_allocation_and_realloc(void **state) {
    void *ptr = judy_malloc(50);
    assert_non_null(ptr);
    
    ptr = judy_realloc(ptr, 100);
    assert_non_null(ptr);
    
    ptr = judy_realloc(ptr, 200);
    assert_non_null(ptr);
    
    ptr = judy_realloc(ptr, 150);
    assert_non_null(ptr);
    
    judy_free(ptr);
}

static void test_allocation_with_zero_values(void **state) {
    void *ptr1 = judy_malloc(0);
    void *ptr2 = judy_calloc(0, 100);
    void *ptr3 = judy_calloc(100, 0);
    
    if (ptr1 != NULL) judy_free(ptr1);
    if (ptr2 != NULL) judy_free(ptr2);
    if (ptr3 != NULL) judy_free(ptr3);
}

static void test_allocation_sizes_across_ranges(void **state) {
    size_t sizes[] = {1, 16, 256, 4096, 65536, 1048576};
    const int num_sizes = sizeof(sizes) / sizeof(sizes[0]);
    
    for (int i = 0; i < num_sizes; i++) {
        void *ptr = judy_malloc(sizes[i]);
        if (ptr != NULL) {
            judy_free(ptr);
        }
    }
}

static void test_realloc_growth_pattern(void **state) {
    void *ptr = judy_malloc(10);
    assert_non_null(ptr);
    
    for (int i = 1; i <= 5; i++) {
        ptr = judy_realloc(ptr, i * 100);
        assert_non_null(ptr);
    }
    
    judy_free(ptr);
}

static void test_realloc_shrinkage_pattern(void **state) {
    void *ptr = judy_malloc(1000);
    assert_non_null(ptr);
    
    for (int i = 5; i >= 1; i--) {
        ptr = judy_realloc(ptr, i * 100);
        assert_non_null(ptr);
    }
    
    judy_free(ptr);
}

static void test_calloc_data_initialization(void **state) {
    void *ptr = judy_calloc(100, 1);
    assert_non_null(ptr);
    
    unsigned char *bytes = (unsigned char *)ptr;
    for (int i = 0; i < 100; i++) {
        assert_int_equal(bytes[i], 0);
    }
    
    judy_free(ptr);
}

static void test_free_and_reuse_pattern(void **state) {
    void *ptr1 = judy_malloc(100);
    assert_non_null(ptr1);
    
    void *addr1 = ptr1;
    judy_free(ptr1);
    
    void *ptr2 = judy_malloc(100);
    assert_non_null(ptr2);
    
    judy_free(ptr2);
}

static void test_malloc_stats_multiple_calls(void **state) {
    for (int i = 0; i < 10; i++) {
        judy_malloc_stats();
    }
}

static void test_malloc_stats_between_allocations(void **state) {
    judy_malloc_stats();
    
    void *ptr1 = judy_malloc(100);
    judy_malloc_stats();
    
    void *ptr2 = judy_calloc(10, 10);
    judy_malloc_stats();
    
    judy_free(ptr1);
    judy_malloc_stats();
    
    judy_free(ptr2);
    judy_malloc_stats();
}

int main(void) {
    const struct CMUnitTest tests[] = {
        /* Header availability tests */
        cmocka_unit_test(test_judy_malloc_header_available),
        cmocka_unit_test(test_judy_calloc_header_available),
        cmocka_unit_test(test_judy_realloc_header_available),
        cmocka_unit_test(test_judy_free_header_available),
        cmocka_unit_test(test_judy_malloc_stats_header_available),
        
        /* Return type tests */
        cmocka_unit_test(test_judy_malloc_returns_void_pointer),
        cmocka_unit_test(test_judy_calloc_returns_void_pointer),
        cmocka_unit_test(test_judy_realloc_returns_void_pointer),
        
        /* Parameter type tests */
        cmocka_unit_test(test_judy_malloc_accepts_size_t),
        cmocka_unit_test(test_judy_calloc_accepts_size_t_parameters),
        cmocka_unit_test(test_judy_realloc_accepts_void_and_size_t),
        cmocka_unit_test(test_judy_free_accepts_void_pointer),
        
        /* Type safety tests */
        cmocka_unit_test(test_type_safety_malloc_return),
        cmocka_unit_test(test_type_safety_calloc_return),
        cmocka_unit_test(test_type_safety_realloc_return),
        
        /* Compound operation tests */
        cmocka_unit_test(test_sequential_allocations_from_header),
        cmocka_unit_test(test_array_of_allocations),
        cmocka_unit_test(test_nested_allocation_and_realloc),
        
        /* Edge case tests */
        cmocka_unit_test(test_allocation_with_zero_values),
        cmocka_unit_test(test_allocation_sizes_across_ranges),
        cmocka_unit_test(test_realloc_growth_pattern),
        cmocka_unit_test(test_realloc_shrinkage_pattern),
        cmocka_unit_test(test_calloc_data_initialization),
        cmocka_unit_test(test_free_and_reuse_pattern),
        
        /* Stats tests */
        cmocka_unit_test(test_malloc_stats_multiple_calls),
        cmocka_unit_test(test_malloc_stats_between_allocations),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}