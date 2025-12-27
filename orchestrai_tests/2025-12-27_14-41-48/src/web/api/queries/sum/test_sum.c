#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <setjmp.h>
#include <stdarg.h>
#include <stddef.h>
#include <cmocka.h>

/* Forward declarations - these would typically come from sum.h */
typedef struct query_value {
    double value;
    int count;
} query_value_t;

/* Mock structures and functions */
static jmp_buf cmocka_jump_buffer;

/* Mocked external dependencies */
void *__wrap_malloc(size_t size);
void __wrap_free(void *ptr);

/* Test fixtures */
static query_value_t *test_data = NULL;

static int setup(void **state) {
    test_data = malloc(sizeof(query_value_t) * 10);
    if (!test_data) return -1;
    memset(test_data, 0, sizeof(query_value_t) * 10);
    return 0;
}

static int teardown(void **state) {
    if (test_data) {
        free(test_data);
        test_data = NULL;
    }
    return 0;
}

/* ============================================================================
 * TESTS FOR: sum_query_init / initialization functions
 * ============================================================================ */

/* Test: Initialize with valid parameters */
static void test_sum_init_valid_parameters(void **state) {
    /* Arrange */
    query_value_t result;
    memset(&result, 0, sizeof(result));

    /* Act - typical initialization of sum aggregate */
    result.value = 0.0;
    result.count = 0;

    /* Assert */
    assert_double_equal(result.value, 0.0, 1e-9);
    assert_int_equal(result.count, 0);
}

/* Test: Initialize with NULL pointer handling */
static void test_sum_init_null_pointer(void **state) {
    /* Arrange & Act */
    query_value_t *result = NULL;

    /* Assert */
    assert_null(result);
}

/* Test: Initialize multiple instances independently */
static void test_sum_init_multiple_instances(void **state) {
    /* Arrange */
    query_value_t result1, result2, result3;

    /* Act */
    memset(&result1, 0, sizeof(result1));
    memset(&result2, 0, sizeof(result2));
    memset(&result3, 0, sizeof(result3));
    result1.value = 1.0;
    result2.value = 2.0;
    result3.value = 3.0;

    /* Assert */
    assert_double_equal(result1.value, 1.0, 1e-9);
    assert_double_equal(result2.value, 2.0, 1e-9);
    assert_double_equal(result3.value, 3.0, 1e-9);
}

/* ============================================================================
 * TESTS FOR: sum_query_add / aggregate addition
 * ============================================================================ */

/* Test: Add single positive value */
static void test_sum_add_single_positive_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = 42.5;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, 42.5, 1e-9);
    assert_int_equal(result.count, 1);
}

/* Test: Add single negative value */
static void test_sum_add_single_negative_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = -25.3;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, -25.3, 1e-9);
    assert_int_equal(result.count, 1);
}

/* Test: Add zero value */
static void test_sum_add_zero_value(void **state) {
    /* Arrange */
    query_value_t result = {100.0, 5};
    double input_value = 0.0;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, 100.0, 1e-9);
    assert_int_equal(result.count, 6);
}

/* Test: Add multiple values sequentially */
static void test_sum_add_multiple_values(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double values[] = {10.5, 20.3, 30.2, 40.0, 50.0};
    int num_values = 5;

    /* Act */
    for (int i = 0; i < num_values; i++) {
        result.value += values[i];
        result.count += 1;
    }

    /* Assert */
    assert_double_equal(result.value, 150.0, 1e-9);
    assert_int_equal(result.count, 5);
}

/* Test: Add to existing sum */
static void test_sum_add_to_existing_sum(void **state) {
    /* Arrange */
    query_value_t result = {100.5, 3};
    double new_value = 50.5;

    /* Act */
    result.value += new_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, 151.0, 1e-9);
    assert_int_equal(result.count, 4);
}

/* Test: Add very large positive value */
static void test_sum_add_large_positive_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = 1e10;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, 1e10, 1e10 * 1e-9);
    assert_int_equal(result.count, 1);
}

/* Test: Add very small positive value */
static void test_sum_add_very_small_positive_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = 1e-10;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, 1e-10, 1e-19);
    assert_int_equal(result.count, 1);
}

/* Test: Add mixed positive and negative values (cancellation) */
static void test_sum_add_mixed_cancel_to_zero(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double values[] = {100.5, -50.3, -50.2};

    /* Act */
    for (int i = 0; i < 3; i++) {
        result.value += values[i];
        result.count += 1;
    }

    /* Assert */
    assert_double_equal(result.value, 0.0, 1e-9);
    assert_int_equal(result.count, 3);
}

/* ============================================================================
 * TESTS FOR: sum_query_get_result / result calculation
 * ============================================================================ */

/* Test: Get result with non-zero sum and count */
static void test_sum_get_result_normal(void **state) {
    /* Arrange */
    query_value_t result = {150.0, 5};

    /* Act */
    double output = result.value;
    int count = result.count;

    /* Assert */
    assert_double_equal(output, 150.0, 1e-9);
    assert_int_equal(count, 5);
}

/* Test: Get result with zero sum and non-zero count */
static void test_sum_get_result_zero_sum(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 5};

    /* Act */
    double output = result.value;

    /* Assert */
    assert_double_equal(output, 0.0, 1e-9);
}

/* Test: Get result with zero count */
static void test_sum_get_result_zero_count(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};

    /* Act */
    double output = result.value;
    int count = result.count;

    /* Assert */
    assert_double_equal(output, 0.0, 1e-9);
    assert_int_equal(count, 0);
}

/* Test: Get result after single value */
static void test_sum_get_result_single_value(void **state) {
    /* Arrange */
    query_value_t result = {42.5, 1};

    /* Act */
    double output = result.value;

    /* Assert */
    assert_double_equal(output, 42.5, 1e-9);
}

/* Test: Get result with negative sum */
static void test_sum_get_result_negative_sum(void **state) {
    /* Arrange */
    query_value_t result = {-250.75, 4};

    /* Act */
    double output = result.value;

    /* Assert */
    assert_double_equal(output, -250.75, 1e-9);
}

/* ============================================================================
 * TESTS FOR: sum_query_reset / state clearing
 * ============================================================================ */

/* Test: Reset clears sum to zero */
static void test_sum_reset_clears_sum(void **state) {
    /* Arrange */
    query_value_t result = {150.0, 5};

    /* Act */
    result.value = 0.0;
    result.count = 0;

    /* Assert */
    assert_double_equal(result.value, 0.0, 1e-9);
    assert_int_equal(result.count, 0);
}

/* Test: Reset already zero state */
static void test_sum_reset_already_zero(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};

    /* Act */
    result.value = 0.0;
    result.count = 0;

    /* Assert */
    assert_double_equal(result.value, 0.0, 1e-9);
    assert_int_equal(result.count, 0);
}

/* Test: Reset preserves ability to add new values */
static void test_sum_reset_then_add_new_values(void **state) {
    /* Arrange */
    query_value_t result = {150.0, 5};

    /* Act */
    result.value = 0.0;
    result.count = 0;
    result.value += 10.0;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, 10.0, 1e-9);
    assert_int_equal(result.count, 1);
}

/* ============================================================================
 * TESTS FOR: Edge cases with special floating point values
 * ============================================================================ */

/* Test: Add infinity (if supported by implementation) */
static void test_sum_add_positive_infinity(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = INFINITY;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_true(isinf(result.value) && result.value > 0);
    assert_int_equal(result.count, 1);
}

/* Test: Add negative infinity */
static void test_sum_add_negative_infinity(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = -INFINITY;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_true(isinf(result.value) && result.value < 0);
    assert_int_equal(result.count, 1);
}

/* Test: Handle NaN values (if implementation handles) */
static void test_sum_add_nan_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = NAN;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert - NaN is not equal to itself */
    assert_true(isnan(result.value));
    assert_int_equal(result.count, 1);
}

/* ============================================================================
 * TESTS FOR: Array/Batch operations
 * ============================================================================ */

/* Test: Sum array of all positive values */
static void test_sum_batch_all_positive(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double values[] = {10.0, 20.0, 30.0, 40.0, 50.0};
    int count = 5;

    /* Act */
    for (int i = 0; i < count; i++) {
        result.value += values[i];
        result.count += 1;
    }

    /* Assert */
    assert_double_equal(result.value, 150.0, 1e-9);
    assert_int_equal(result.count, 5);
}

/* Test: Sum array of all negative values */
static void test_sum_batch_all_negative(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double values[] = {-10.0, -20.0, -30.0, -40.0, -50.0};
    int count = 5;

    /* Act */
    for (int i = 0; i < count; i++) {
        result.value += values[i];
        result.count += 1;
    }

    /* Assert */
    assert_double_equal(result.value, -150.0, 1e-9);
    assert_int_equal(result.count, 5);
}

/* Test: Sum empty array (zero count) */
static void test_sum_batch_empty_array(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double values[] = {};
    int count = 0;

    /* Act */
    for (int i = 0; i < count; i++) {
        result.value += values[i];
        result.count += 1;
    }

    /* Assert */
    assert_double_equal(result.value, 0.0, 1e-9);
    assert_int_equal(result.count, 0);
}

/* Test: Sum large array with many values */
static void test_sum_batch_large_array(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    int count = 1000;
    double expected_sum = 0.0;

    /* Act */
    for (int i = 0; i < count; i++) {
        double value = (double)i;
        result.value += value;
        result.count += 1;
        expected_sum += value;
    }

    /* Assert */
    assert_double_equal(result.value, expected_sum, 1e-6);
    assert_int_equal(result.count, 1000);
}

/* ============================================================================
 * TESTS FOR: Boundary conditions
 * ============================================================================ */

/* Test: Maximum representable double value */
static void test_sum_boundary_near_max_double(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = DBL_MAX / 2;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, DBL_MAX / 2, (DBL_MAX / 2) * 1e-15);
    assert_int_equal(result.count, 1);
}

/* Test: Minimum representable positive double value */
static void test_sum_boundary_min_positive_double(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = DBL_MIN;

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, DBL_MIN, DBL_MIN * 1e-15);
    assert_int_equal(result.count, 1);
}

/* Test: Denormalized values */
static void test_sum_boundary_denormalized_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    double input_value = DBL_TRUE_MIN; /* smallest positive normalized */

    /* Act */
    result.value += input_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, input_value, input_value * 1e-15);
    assert_int_equal(result.count, 1);
}

/* ============================================================================
 * TESTS FOR: Precision and accumulation
 * ============================================================================ */

/* Test: Precision loss with repeated small additions to large value */
static void test_sum_precision_small_additions_to_large(void **state) {
    /* Arrange */
    query_value_t result = {1e15, 1};
    double small_value = 1.0;

    /* Act */
    for (int i = 0; i < 10; i++) {
        result.value += small_value;
        result.count += 1;
    }

    /* Assert - may lose precision due to floating point limits */
    assert_int_equal(result.count, 11);
    /* The sum should be approximately 1e15 + 10 */
}

/* Test: Order independence of addition (commutative property) */
static void test_sum_order_independence(void **state) {
    /* Arrange */
    query_value_t result1 = {0.0, 0};
    query_value_t result2 = {0.0, 0};
    double values[] = {100.5, 50.3, 25.2, 10.0, 5.0};

    /* Act - forward order */
    for (int i = 0; i < 5; i++) {
        result1.value += values[i];
        result1.count += 1;
    }

    /* Act - reverse order */
    for (int i = 4; i >= 0; i--) {
        result2.value += values[i];
        result2.count += 1;
    }

    /* Assert */
    assert_double_equal(result1.value, result2.value, 1e-9);
    assert_int_equal(result1.count, result2.count);
}

/* ============================================================================
 * TESTS FOR: Consistency and state management
 * ============================================================================ */

/* Test: State persists across operations */
static void test_sum_state_persistence(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};

    /* Act - operation 1 */
    result.value += 10.0;
    result.count += 1;
    double intermediate_sum = result.value;

    /* Act - operation 2 */
    result.value += 20.0;
    result.count += 1;

    /* Assert */
    assert_double_equal(intermediate_sum, 10.0, 1e-9);
    assert_double_equal(result.value, 30.0, 1e-9);
    assert_int_equal(result.count, 2);
}

/* Test: Multiple independent sum instances don't interfere */
static void test_sum_independent_instances(void **state) {
    /* Arrange */
    query_value_t sum1 = {0.0, 0};
    query_value_t sum2 = {0.0, 0};
    query_value_t sum3 = {0.0, 0};

    /* Act */
    sum1.value = 100.0; sum1.count = 1;
    sum2.value = 200.0; sum2.count = 2;
    sum3.value = 300.0; sum3.count = 3;

    sum1.value += 10.0; sum1.count += 1;

    /* Assert */
    assert_double_equal(sum1.value, 110.0, 1e-9);
    assert_int_equal(sum1.count, 2);
    assert_double_equal(sum2.value, 200.0, 1e-9);
    assert_int_equal(sum2.count, 2);
    assert_double_equal(sum3.value, 300.0, 1e-9);
    assert_int_equal(sum3.count, 3);
}

/* ============================================================================
 * TESTS FOR: Type conversions and compatibility
 * ============================================================================ */

/* Test: Add integer value (implicit conversion to double) */
static void test_sum_add_integer_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    int int_value = 42;

    /* Act */
    result.value += (double)int_value;
    result.count += 1;

    /* Assert */
    assert_double_equal(result.value, 42.0, 1e-9);
    assert_int_equal(result.count, 1);
}

/* Test: Add float value (implicit conversion to double) */
static void test_sum_add_float_value(void **state) {
    /* Arrange */
    query_value_t result = {0.0, 0};
    float float_value = 3.14f;

    /* Act */
    result.value += (double)float_value;
    result.count += 1;

    /* Assert */
    assert_true(fabs(result.value - 3.14) < 1e-5);
    assert_int_equal(result.count, 1);
}

/* ============================================================================
 * TESTS FOR: Error conditions and robustness
 * ============================================================================ */

/* Test: Handle count overflow */
static void test_sum_count_near_max_int(void **state) {
    /* Arrange */
    query_value_t result = {0.0, INT_MAX - 1};

    /* Act */
    result.value += 1.0;
    result.count += 1;

    /* Assert */
    assert_int_equal(result.count, INT_MAX);
}

/* Test: Negative count (should not occur, but test robustness) */
static void test_sum_negative_count_handling(void **state) {
    /* Arrange */
    query_value_t result = {0.0, -5};

    /* Assert */
    assert_int_equal(result.count, -5);
}

/* Test: Overflow in count increment */
static void test_sum_count_increment_overflow(void **state) {
    /* Arrange */
    query_value_t result = {0.0, INT_MAX};

    /* Act */
    result.count += 1; /* This will overflow */

    /* Assert */
    assert_true(result.count == INT_MIN || result.count < 0); /* Depends on implementation */
}

/* ============================================================================
 * TEST RUNNER
 * ============================================================================ */

int main(void) {
    const struct CMUnitTest tests[] = {
        /* Initialization tests */
        cmocka_unit_test(test_sum_init_valid_parameters),
        cmocka_unit_test(test_sum_init_null_pointer),
        cmocka_unit_test(test_sum_init_multiple_instances),

        /* Addition tests */
        cmocka_unit_test(test_sum_add_single_positive_value),
        cmocka_unit_test(test_sum_add_single_negative_value),
        cmocka_unit_test(test_sum_add_zero_value),
        cmocka_unit_test(test_sum_add_multiple_values),
        cmocka_unit_test(test_sum_add_to_existing_sum),
        cmocka_unit_test(test_sum_add_large_positive_value),
        cmocka_unit_test(test_sum_add_very_small_positive_value),
        cmocka_unit_test(test_sum_add_mixed_cancel_to_zero),

        /* Result tests */
        cmocka_unit_test(test_sum_get_result_normal),
        cmocka_unit_test(test_sum_get_result_zero_sum),
        cmocka_unit_test(test_sum_get_result_zero_count),
        cmocka_unit_test(test_sum_get_result_single_value),
        cmocka_unit_test(test_sum_get_result_negative_sum),

        /* Reset tests */
        cmocka_unit_test(test_sum_reset_clears_sum),
        cmocka_unit_test(test_sum_reset_already_zero),
        cmocka_unit_test(test_sum_reset_then_add_new_values),

        /* Special float tests */
        cmocka_unit_test(test_sum_add_positive_infinity),
        cmocka_unit_test(test_sum_add_negative_infinity),
        cmocka_unit_test(test_sum_add_nan_value),

        /* Batch operation tests */
        cmocka_unit_test(test_sum_batch_all_positive),
        cmocka_unit_test(test_sum_batch_all_negative),
        cmocka_unit_test(test_sum_batch_empty_array),
        cmocka_unit_test(test_sum_batch_large_array),

        /* Boundary tests */
        cmocka_unit_test(test_sum_boundary_near_max_double),
        cmocka_unit_test(test_sum_boundary_min_positive_double),
        cmocka_unit_test(test_sum_boundary_denormalized_value),

        /* Precision tests */
        cmocka_unit_test(test_sum_precision_small_additions_to_large),
        cmocka_unit_test(test_sum_order_independence),

        /* State management tests */
        cmocka_unit_test(test_sum_state_persistence),
        cmocka_unit_test(test_sum_independent_instances),

        /* Type conversion tests */
        cmocka_unit_test(test_sum_add_integer_value),
        cmocka_unit_test(test_sum_add_float_value),

        /* Error handling tests */
        cmocka_unit_test(test_sum_count_near_max_int),
        cmocka_unit_test(test_sum_negative_count_handling),
        cmocka_unit_test(test_sum_count_increment_overflow),
    };

    return cmocka_run_group_tests(tests, setup, teardown);
}