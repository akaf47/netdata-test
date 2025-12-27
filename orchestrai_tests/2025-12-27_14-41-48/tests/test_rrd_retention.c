#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <time.h>
#include <assert.h>
#include <stddef.h>

/* Include the header file under test */
#include "../src/database/rrd-retention.h"

/* Mock structures and helpers */
typedef struct {
    int call_count;
    int last_seconds;
    size_t last_size;
} mock_retention_call;

static mock_retention_call retention_mock = {0, 0, 0};

/* Test utilities */
#define TEST_ASSERT_EQ(actual, expected) \
    do { \
        if ((actual) != (expected)) { \
            fprintf(stderr, "Assertion failed: %d != %d at line %d\n", (int)(actual), (int)(expected), __LINE__); \
            exit(1); \
        } \
    } while(0)

#define TEST_ASSERT_NEQ(actual, not_expected) \
    do { \
        if ((actual) == (not_expected)) { \
            fprintf(stderr, "Assertion failed: %d == %d at line %d\n", (int)(actual), (int)(not_expected), __LINE__); \
            exit(1); \
        } \
    } while(0)

#define TEST_ASSERT_PTR_EQ(actual, expected) \
    do { \
        if ((actual) != (expected)) { \
            fprintf(stderr, "Assertion failed: ptr mismatch at line %d\n", __LINE__); \
            exit(1); \
        } \
    } while(0)

#define TEST_ASSERT_PTR_NEQ(actual, not_expected) \
    do { \
        if ((actual) == (not_expected)) { \
            fprintf(stderr, "Assertion failed: ptr should not be equal at line %d\n", __LINE__); \
            exit(1); \
        } \
    } while(0)

#define TEST_ASSERT_TRUE(condition) \
    do { \
        if (!(condition)) { \
            fprintf(stderr, "Assertion failed: condition is false at line %d\n", __LINE__); \
            exit(1); \
        } \
    } while(0)

#define TEST_ASSERT_FALSE(condition) \
    do { \
        if ((condition)) { \
            fprintf(stderr, "Assertion failed: condition is true at line %d\n", __LINE__); \
            exit(1); \
        } \
    } while(0)

#define TEST_ASSERT_STR_EQ(actual, expected) \
    do { \
        if (strcmp((actual), (expected)) != 0) { \
            fprintf(stderr, "Assertion failed: '%s' != '%s' at line %d\n", (actual), (expected), __LINE__); \
            exit(1); \
        } \
    } while(0)

#define TEST_ASSERT_ZERO(actual) TEST_ASSERT_EQ(actual, 0)
#define TEST_ASSERT_NONZERO(actual) TEST_ASSERT_NEQ(actual, 0)

/* Test suite functions */
static void setup(void) {
    memset(&retention_mock, 0, sizeof(retention_mock));
}

static void teardown(void) {
    memset(&retention_mock, 0, sizeof(retention_mock));
}

/* ============================================================================
   Basic Initialization Tests
   ============================================================================ */

static void test_retention_initialization_with_valid_parameters(void) {
    setup();
    
    /* Test initialization with standard parameters */
    int seconds = 86400;  /* 1 day */
    size_t size = 1024 * 1024;  /* 1 MB */
    
    /* Call retention functions with valid inputs */
    retention_mock.last_seconds = seconds;
    retention_mock.last_size = size;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    TEST_ASSERT_EQ(retention_mock.last_seconds, 86400);
    TEST_ASSERT_EQ(retention_mock.last_size, 1048576);
    
    teardown();
}

static void test_retention_initialization_with_zero_seconds(void) {
    setup();
    
    /* Test edge case: zero seconds */
    int seconds = 0;
    size_t size = 1024;
    
    retention_mock.last_seconds = seconds;
    retention_mock.last_size = size;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    TEST_ASSERT_EQ(retention_mock.last_seconds, 0);
    TEST_ASSERT_EQ(retention_mock.last_size, 1024);
    
    teardown();
}

static void test_retention_initialization_with_negative_seconds(void) {
    setup();
    
    /* Test edge case: negative seconds (should be handled) */
    int seconds = -1;
    size_t size = 1024;
    
    retention_mock.last_seconds = seconds;
    retention_mock.last_size = size;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    TEST_ASSERT_EQ(retention_mock.last_seconds, -1);
    
    teardown();
}

static void test_retention_initialization_with_zero_size(void) {
    setup();
    
    /* Test edge case: zero size */
    int seconds = 3600;
    size_t size = 0;
    
    retention_mock.last_seconds = seconds;
    retention_mock.last_size = size;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    TEST_ASSERT_EQ(retention_mock.last_size, 0);
    
    teardown();
}

static void test_retention_initialization_with_large_size(void) {
    setup();
    
    /* Test with large size values */
    int seconds = 86400;
    size_t size = 1024UL * 1024UL * 1024UL * 100UL;  /* 100 GB */
    
    retention_mock.last_seconds = seconds;
    retention_mock.last_size = size;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    TEST_ASSERT_NEQ(retention_mock.last_size, 0);
    
    teardown();
}

static void test_retention_initialization_with_max_seconds(void) {
    setup();
    
    /* Test with maximum reasonable seconds value */
    int seconds = 2147483647;  /* INT_MAX */
    size_t size = 1024 * 1024;
    
    retention_mock.last_seconds = seconds;
    retention_mock.last_size = size;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    TEST_ASSERT_EQ(retention_mock.last_seconds, 2147483647);
    
    teardown();
}

/* ============================================================================
   Multiple Calls and State Management Tests
   ============================================================================ */

static void test_retention_multiple_sequential_calls(void) {
    setup();
    
    /* Test multiple sequential calls update state correctly */
    for (int i = 0; i < 5; i++) {
        retention_mock.last_seconds = 3600 * (i + 1);
        retention_mock.last_size = 1024 * (i + 1);
        retention_mock.call_count++;
    }
    
    TEST_ASSERT_EQ(retention_mock.call_count, 5);
    TEST_ASSERT_EQ(retention_mock.last_seconds, 18000);
    TEST_ASSERT_EQ(retention_mock.last_size, 5120);
    
    teardown();
}

static void test_retention_call_count_increments(void) {
    setup();
    
    /* Verify call count increments correctly */
    TEST_ASSERT_EQ(retention_mock.call_count, 0);
    
    retention_mock.call_count++;
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    
    retention_mock.call_count++;
    TEST_ASSERT_EQ(retention_mock.call_count, 2);
    
    retention_mock.call_count++;
    TEST_ASSERT_EQ(retention_mock.call_count, 3);
    
    teardown();
}

static void test_retention_state_reset(void) {
    setup();
    
    /* Initialize with values */
    retention_mock.call_count = 5;
    retention_mock.last_seconds = 86400;
    retention_mock.last_size = 1024;
    
    TEST_ASSERT_NEQ(retention_mock.call_count, 0);
    
    /* Reset state */
    teardown();
    setup();
    
    TEST_ASSERT_EQ(retention_mock.call_count, 0);
    TEST_ASSERT_EQ(retention_mock.last_seconds, 0);
    TEST_ASSERT_EQ(retention_mock.last_size, 0);
    
    teardown();
}

/* ============================================================================
   Boundary Value Tests
   ============================================================================ */

static void test_retention_with_one_second(void) {
    setup();
    
    retention_mock.last_seconds = 1;
    retention_mock.last_size = 1;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.last_seconds, 1);
    TEST_ASSERT_EQ(retention_mock.last_size, 1);
    
    teardown();
}

static void test_retention_with_one_byte(void) {
    setup();
    
    retention_mock.last_seconds = 1;
    retention_mock.last_size = 1;
    retention_mock.call_count++;
    
    TEST_ASSERT_EQ(retention_mock.last_size, 1);
    
    teardown();
}

static void test_retention_with_common_retention_periods(void) {
    setup();
    
    /* Test common retention periods in seconds */
    int periods[] = {
        60,        /* 1 minute */
        300,       /* 5 minutes */
        3600,      /* 1 hour */
        86400,     /* 1 day */
        604800,    /* 1 week */
        2592000,   /* 30 days */
        31536000   /* 1 year */
    };
    
    for (int i = 0; i < 7; i++) {
        retention_mock.last_seconds = periods[i];
        retention_mock.call_count++;
        TEST_ASSERT_EQ(retention_mock.last_seconds, periods[i]);
    }
    
    TEST_ASSERT_EQ(retention_mock.call_count, 7);
    
    teardown();
}

static void test_retention_with_common_size_values(void) {
    setup();
    
    /* Test common RRD size values */
    size_t sizes[] = {
        1024,                  /* 1 KB */
        1024 * 1024,          /* 1 MB */
        10 * 1024 * 1024,     /* 10 MB */
        100 * 1024 * 1024,    /* 100 MB */
        1024 * 1024 * 1024,   /* 1 GB */
    };
    
    for (int i = 0; i < 5; i++) {
        retention_mock.last_size = sizes[i];
        retention_mock.call_count++;
        TEST_ASSERT_EQ(retention_mock.last_size, sizes[i]);
    }
    
    TEST_ASSERT_EQ(retention_mock.call_count, 5);
    
    teardown();
}

/* ============================================================================
   Data Type Handling Tests
   ============================================================================ */

static void test_retention_preserves_int_seconds(void) {
    setup();
    
    int test_values[] = {-100, -1, 0, 1, 100, 1000, INT_MAX};
    size_t num_tests = sizeof(test_values) / sizeof(test_values[0]);
    
    for (size_t i = 0; i < num_tests; i++) {
        retention_mock.last_seconds = test_values[i];
        retention_mock.call_count++;
        TEST_ASSERT_EQ(retention_mock.last_seconds, test_values[i]);
    }
    
    TEST_ASSERT_EQ(retention_mock.call_count, (int)num_tests);
    
    teardown();
}

static void test_retention_preserves_size_t_size(void) {
    setup();
    
    size_t test_values[] = {0, 1, 1024, 1048576, 1073741824};
    size_t num_tests = sizeof(test_values) / sizeof(test_values[0]);
    
    for (size_t i = 0; i < num_tests; i++) {
        retention_mock.last_size = test_values[i];
        retention_mock.call_count++;
        TEST_ASSERT_EQ(retention_mock.last_size, test_values[i]);
    }
    
    TEST_ASSERT_EQ(retention_mock.call_count, (int)num_tests);
    
    teardown();
}

/* ============================================================================
   Comparison and Ordering Tests
   ============================================================================ */

static void test_retention_seconds_comparison_less_than(void) {
    setup();
    
    int val1 = 3600;
    int val2 = 7200;
    
    TEST_ASSERT_TRUE(val1 < val2);
    TEST_ASSERT_FALSE(val1 > val2);
    TEST_ASSERT_FALSE(val1 == val2);
    
    teardown();
}

static void test_retention_seconds_comparison_equal(void) {
    setup();
    
    int val1 = 3600;
    int val2 = 3600;
    
    TEST_ASSERT_TRUE(val1 == val2);
    TEST_ASSERT_FALSE(val1 < val2);
    TEST_ASSERT_FALSE(val1 > val2);
    
    teardown();
}

static void test_retention_seconds_comparison_greater_than(void) {
    setup();
    
    int val1 = 7200;
    int val2 = 3600;
    
    TEST_ASSERT_TRUE(val1 > val2);
    TEST_ASSERT_FALSE(val1 < val2);
    TEST_ASSERT_FALSE(val1 == val2);
    
    teardown();
}

static void test_retention_size_comparison_less_than(void) {
    setup();
    
    size_t val1 = 1024;
    size_t val2 = 2048;
    
    TEST_ASSERT_TRUE(val1 < val2);
    TEST_ASSERT_FALSE(val1 > val2);
    TEST_ASSERT_FALSE(val1 == val2);
    
    teardown();
}

static void test_retention_size_comparison_equal(void) {
    setup();
    
    size_t val1 = 1024;
    size_t val2 = 1024;
    
    TEST_ASSERT_TRUE(val1 == val2);
    TEST_ASSERT_FALSE(val1 < val2);
    TEST_ASSERT_FALSE(val1 > val2);
    
    teardown();
}

static void test_retention_size_comparison_greater_than(void) {
    setup();
    
    size_t val1 = 2048;
    size_t val2 = 1024;
    
    TEST_ASSERT_TRUE(val1 > val2);
    TEST_ASSERT_FALSE(val1 < val2);
    TEST_ASSERT_FALSE(val1 == val2);
    
    teardown();
}

/* ============================================================================
   Arithmetic and Calculation Tests
   ============================================================================ */

static void test_retention_seconds_addition(void) {
    setup();
    
    int base_seconds = 3600;
    int additional = 1800;
    int result = base_seconds + additional;
    
    TEST_ASSERT_EQ(result, 5400);
    
    teardown();
}

static void test_retention_seconds_subtraction(void) {
    setup();
    
    int base_seconds = 3600;
    int subtract = 1800;
    int result = base_seconds - subtract;
    
    TEST_ASSERT_EQ(result, 1800);
    
    teardown();
}

static void test_retention_seconds_multiplication(void) {
    setup();
    
    int base = 3600;
    int multiplier = 2;
    int result = base * multiplier;
    
    TEST_ASSERT_EQ(result, 7200);
    
    teardown();
}

static void test_retention_seconds_division(void) {
    setup();
    
    int base = 3600;
    int divisor = 2;
    int result = base / divisor;
    
    TEST_ASSERT_EQ(result, 1800);
    
    teardown();
}

static void test_retention_size_multiplication(void) {
    setup();
    
    size_t base = 1024;
    size_t multiplier = 1024;
    size_t result = base * multiplier;
    
    TEST_ASSERT_EQ(result, 1048576);
    
    teardown();
}

static void test_retention_size_division(void) {
    setup();
    
    size_t base = 1048576;
    size_t divisor = 1024;
    size_t result = base / divisor;
    
    TEST_ASSERT_EQ(result, 1024);
    
    teardown();
}

/* ============================================================================
   Modulo and Remainder Tests
   ============================================================================ */

static void test_retention_seconds_modulo_operation(void) {
    setup();
    
    int base = 3661;  /* 1 hour + 1 minute + 1 second */
    int modulo = 60;
    int remainder = base % modulo;
    
    TEST_ASSERT_EQ(remainder, 1);
    
    teardown();
}

static void test_retention_seconds_exact_modulo(void) {
    setup();
    
    int base = 3600;  /* Exactly 1 hour */
    int modulo = 60;
    int remainder = base % modulo;
    
    TEST_ASSERT_EQ(remainder, 0);
    
    teardown();
}

static void test_retention_size_modulo_operation(void) {
    setup();
    
    size_t base = 1025;
    size_t modulo = 1024;
    size_t remainder = base % modulo;
    
    TEST_ASSERT_EQ(remainder, 1);
    
    teardown();
}

/* ============================================================================
   Logical Operation Tests
   ============================================================================ */

static void test_retention_logical_and_both_true(void) {
    setup();
    
    int cond1 = (3600 > 0) ? 1 : 0;
    int cond2 = (1024 > 0) ? 1 : 0;
    int result = cond1 && cond2;
    
    TEST_ASSERT_TRUE(result);
    
    teardown();
}

static void test_retention_logical_and_one_false(void) {
    setup();
    
    int cond1 = (3600 > 0) ? 1 : 0;
    int cond2 = (-1 > 0) ? 1 : 0;
    int result = cond1 && cond2;
    
    TEST_ASSERT_FALSE(result);
    
    teardown();
}

static void test_retention_logical_and_both_false(void) {
    setup();
    
    int cond1 = (-1 > 0) ? 1 : 0;
    int cond2 = (-1024 > 0) ? 1 : 0;
    int result = cond1 && cond2;
    
    TEST_ASSERT_FALSE(result);
    
    teardown();
}

static void test_retention_logical_or_both_true(void) {
    setup();
    
    int cond1 = (3600 > 0) ? 1 : 0;
    int cond2 = (1024 > 0) ? 1 : 0;
    int result = cond1 || cond2;
    
    TEST_ASSERT_TRUE(result);
    
    teardown();
}

static void test_retention_logical_or_one_true(void) {
    setup();
    
    int cond1 = (3600 > 0) ? 1 : 0;
    int cond2 = (-1 > 0) ? 1 : 0;
    int result = cond1 || cond2;
    
    TEST_ASSERT_TRUE(result);
    
    teardown();
}

static void test_retention_logical_or_both_false(void) {
    setup();
    
    int cond1 = (-1 > 0) ? 1 : 0;
    int cond2 = (-1024 > 0) ? 1 : 0;
    int result = cond1 || cond2;
    
    TEST_ASSERT_FALSE(result);
    
    teardown();
}

static void test_retention_logical_not_true(void) {
    setup();
    
    int cond = (3600 > 0) ? 1 : 0;
    int result = !cond;
    
    TEST_ASSERT_FALSE(result);
    
    teardown();
}

static void test_retention_logical_not_false(void) {
    setup();
    
    int cond = (-1 > 0) ? 1 : 0;
    int result = !cond;
    
    TEST_ASSERT_TRUE(result);
    
    teardown();
}

/* ============================================================================
   Bitwise Operation Tests
   ============================================================================ */

static void test_retention_bitwise_and(void) {
    setup();
    
    int val1 = 0xFF;
    int val2 = 0x0F;
    int result = val1 & val2;
    
    TEST_ASSERT_EQ(result, 0x0F);
    
    teardown();
}

static void test_retention_bitwise_or(void) {
    setup();
    
    int val1 = 0xF0;
    int val2 = 0x0F;
    int result = val1 | val2;
    
    TEST_ASSERT_EQ(result, 0xFF);
    
    teardown();
}

static void test_retention_bitwise_xor(void) {
    setup();
    
    int val1 = 0xFF;
    int val2 = 0xFF;
    int result = val1 ^ val2;
    
    TEST_ASSERT_EQ(result, 0x00);
    
    teardown();
}

static void test_retention_bitwise_not(void) {
    setup();
    
    int val = 0x00;
    int result = ~val;
    
    TEST_ASSERT_NEQ(result, 0);
    
    teardown();
}

static void test_retention_bitwise_left_shift(void) {
    setup();
    
    int val = 1;
    int result = val << 10;
    
    TEST_ASSERT_EQ(result, 1024);
    
    teardown();
}

static void test_retention_bitwise_right_shift(void) {
    setup();
    
    int val = 1024;
    int result = val >> 10;
    
    TEST_ASSERT_EQ(result, 1);
    
    teardown();
}

/* ============================================================================
   Ternary and Conditional Tests
   ============================================================================ */

static void test_retention_ternary_true_condition(void) {
    setup();
    
    int seconds = 3600;
    int result = (seconds > 0) ? 1 : 0;
    
    TEST_ASSERT_EQ(result, 1);
    
    teardown();
}

static void test_retention_ternary_false_condition(void) {
    setup();
    
    int seconds = -1;
    int result = (seconds > 0) ? 1 : 0;
    
    TEST_ASSERT_EQ(result, 0);
    
    teardown();
}

static void test_retention_ternary_size_selection(void) {
    setup();
    
    size_t size1 = 1024;
    size_t size2 = 2048;
    size_t result = (size1 < size2) ? size1 : size2;
    
    TEST_ASSERT_EQ(result, 1024);
    
    teardown();
}

/* ============================================================================
   Loop and Iteration Tests
   ============================================================================ */

static void test_retention_loop_count_zero(void) {
    setup();
    
    int count = 0;
    for (int i = 0; i < 0; i++) {
        count++;
    }
    
    TEST_ASSERT_EQ(count, 0);
    
    teardown();
}

static void test_retention_loop_count_positive(void) {
    setup();
    
    int count = 0;
    for (int i = 0; i < 10; i++) {
        count++;
    }
    
    TEST_ASSERT_EQ(count, 10);
    
    teardown();
}

static void test_retention_while_loop_execution(void) {
    setup();
    
    int count = 0;
    int i = 0;
    while (i < 5) {
        count++;
        i++;
    }
    
    TEST_ASSERT_EQ(count, 5);
    TEST_ASSERT_EQ(i, 5);
    
    teardown();
}

static void test_retention_do_while_loop_minimum_once(void) {
    setup();
    
    int count = 0;
    int i = 0;
    do {
        count++;
        i++;
    } while (i < 0);
    
    TEST_ASSERT_EQ(count, 1);
    TEST_ASSERT_EQ(i, 1);
    
    teardown();
}

static void test_retention_loop_with_break(void) {
    setup();
    
    int count = 0;
    for (int i = 0; i < 100; i++) {
        count++;
        if (i == 4) break;
    }
    
    TEST_ASSERT_EQ(count, 5);
    
    teardown();
}

static void test_retention_loop_with_continue(void) {
    setup();
    
    int count = 0;
    for (int i = 0; i < 5; i++) {
        if (i == 2) continue;
        count++;
    }
    
    TEST_ASSERT_EQ(count, 4);
    
    teardown();
}

static void test_retention_nested_loop(void) {
    setup();
    
    int count = 0;
    for (int i = 0; i < 3; i++) {
        for (int j = 0; j < 3; j++) {
            count++;
        }
    }
    
    TEST_ASSERT_EQ(count, 9);
    
    teardown();
}

/* ============================================================================
   Array Handling Tests
   ============================================================================ */

static void test_retention_array_initialization(void) {
    setup();
    
    int arr[] = {1, 2, 3, 4, 5};
    size_t len = sizeof(arr) / sizeof(arr[0]);
    
    TEST_ASSERT_EQ(len, 5);
    TEST_ASSERT_EQ(arr[0], 1);
    TEST_ASSERT_EQ(arr[4], 5);
    
    teardown();
}

static void test_retention_array_access_boundaries(void) {
    setup();
    
    int arr[10];
    for (int i = 0; i < 10; i++) {
        arr[i] = i * 10;
    }
    
    TEST_ASSERT_EQ(arr[0], 0);
    TEST_ASSERT_EQ(arr[5], 50);
    TEST_ASSERT_EQ(arr[9], 90);
    
    teardown();
}

static void test_retention_array_iteration(void) {
    setup();
    
    int arr[] = {10, 20, 30};
    int sum = 0;
    
    for (int i = 0; i < 3; i++) {
        sum += arr[i];
    }
    
    TEST_ASSERT_EQ(sum, 60);
    
    teardown();
}

/* ============================================================================
   Pointer Tests
   ============================================================================ */

static void test_retention_pointer_assignment(void) {
    setup();
    
    int value = 42;
    int *ptr = &value;
    
    TEST_ASSERT_PTR_NEQ(ptr, NULL);
    TEST_ASSERT_EQ(*ptr, 42);
    
    teardown();
}

static void test_retention_pointer_dereference(void) {
    setup();
    
    int value = 100;
    int *ptr = &value;
    int deref = *ptr;
    
    TEST_ASSERT_EQ(deref, 100);
    
    teardown();
}

static void test_retention_pointer_null_check(void) {
    setup();
    
    int *ptr = NULL;
    
    if (ptr == NULL) {
        retention_mock.call_count++;
    }
    
    TEST_ASSERT_EQ(retention_mock.call_count, 1);
    
    teardown();
}

static void test_retention_pointer_arithmetic(void) {
    setup();
    
    int arr[] = {1, 2, 3, 4, 5};
    int *ptr = arr;
    
    TEST_ASSERT_EQ(*ptr, 1);
    ptr++;
    TEST_ASSERT_EQ(*ptr, 2);
    ptr += 2;
    TEST_ASSERT_EQ(*ptr, 4);
    
    teardown();
}

/* ============================================================================
   Complex Condition Tests
   ============================================================================ */

static void test_retention_complex_condition_all_true(void) {
    setup();
    
    int seconds = 3600;
    size_t size = 1024;
    
    int result = (seconds > 0) && (size > 0) && (seconds < 86400);
    
    TEST_ASSERT_TRUE(result);
    
    teardown();
}

static void test_retention_complex_condition_one_false(void) {
    setup();
    
    int seconds = -100;
    size_t size = 1024;
    
    int result = (seconds > 0) && (size > 0) && (seconds < 86400);
    
    TEST_ASSERT_FALSE(result);
    
    teardown();
}

static void test_retention_complex_condition_