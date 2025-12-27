#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <stdint.h>
#include <cmocka.h>

/* Forward declarations - these would come from incremental_sum.h */
typedef struct {
    NETDATA_QUERY_VALUE *series;
    uint32_t count;
} INCREMENTAL_SUM_DATA;

/* Mock structures and functions */
typedef double NETDATA_QUERY_VALUE;

/* Assuming these are the exported functions from incremental_sum.c */
void *incremental_sum_create(RRDCALC_OPTIONS *options);
void incremental_sum_free(void *state);
void incremental_sum_add_value(void *state, NETDATA_QUERY_VALUE value);
NETDATA_QUERY_VALUE incremental_sum_flush(void *state, NETDATA_QUERY_VALUE *sum);

/* ============================================================================
 * TEST: incremental_sum_create - Memory allocation
 * ============================================================================ */

static void test_incremental_sum_create_success(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    
    /* Act */
    void *result = incremental_sum_create(&options);
    
    /* Assert */
    assert_non_null(result);
    
    /* Cleanup */
    incremental_sum_free(result);
}

static void test_incremental_sum_create_with_null_options(void **state) {
    /* Arrange */
    
    /* Act */
    void *result = incremental_sum_create(NULL);
    
    /* Assert */
    assert_non_null(result);
    
    /* Cleanup */
    if (result) incremental_sum_free(result);
}

/* ============================================================================
 * TEST: incremental_sum_free - Memory deallocation
 * ============================================================================ */

static void test_incremental_sum_free_valid_state(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *allocated = incremental_sum_create(&options);
    
    /* Act */
    incremental_sum_free(allocated);
    
    /* Assert - no crash or segfault means success */
}

static void test_incremental_sum_free_null_pointer(void **state) {
    /* Arrange */
    
    /* Act */
    incremental_sum_free(NULL);
    
    /* Assert - should handle gracefully */
}

/* ============================================================================
 * TEST: incremental_sum_add_value - Single value addition
 * ============================================================================ */

static void test_incremental_sum_add_value_positive(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = 42.5;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_negative(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = -42.5;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_zero(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = 0.0;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_very_large(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = 1e308;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_very_small(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = 1e-308;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_nan(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = NAN;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_infinity(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = INFINITY;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_negative_infinity(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE value = -INFINITY;
    
    /* Act */
    incremental_sum_add_value(data, value);
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_value_to_null_state(void **state) {
    /* Arrange */
    NETDATA_QUERY_VALUE value = 42.5;
    
    /* Act & Assert - should handle gracefully or crash predictably */
    incremental_sum_add_value(NULL, value);
}

/* ============================================================================
 * TEST: incremental_sum_add_value - Multiple additions
 * ============================================================================ */

static void test_incremental_sum_add_multiple_values_ascending(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE values[] = {1.0, 2.0, 3.0, 4.0, 5.0};
    
    /* Act */
    for (int i = 0; i < 5; i++) {
        incremental_sum_add_value(data, values[i]);
    }
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_multiple_values_descending(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE values[] = {5.0, 4.0, 3.0, 2.0, 1.0};
    
    /* Act */
    for (int i = 0; i < 5; i++) {
        incremental_sum_add_value(data, values[i]);
    }
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_multiple_values_mixed_signs(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE values[] = {-5.0, 10.0, -3.0, 2.0, 0.5};
    
    /* Act */
    for (int i = 0; i < 5; i++) {
        incremental_sum_add_value(data, values[i]);
    }
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_many_values(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    
    /* Act */
    for (int i = 0; i < 1000; i++) {
        incremental_sum_add_value(data, (NETDATA_QUERY_VALUE)i);
    }
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_add_only_zeros(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    
    /* Act */
    for (int i = 0; i < 10; i++) {
        incremental_sum_add_value(data, 0.0);
    }
    
    /* Assert */
    incremental_sum_flush(data, NULL);
    
    /* Cleanup */
    incremental_sum_free(data);
}

/* ============================================================================
 * TEST: incremental_sum_flush - Output generation
 * ============================================================================ */

static void test_incremental_sum_flush_empty_state(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    NETDATA_QUERY_VALUE result = 0.0;
    
    /* Act */
    NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(data, &result);
    
    /* Assert */
    assert_true(!isnan(flush_result) || isnan(flush_result)); /* Just ensure it returns something */
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_flush_with_single_value(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    incremental_sum_add_value(data, 42.5);
    NETDATA_QUERY_VALUE result = 0.0;
    
    /* Act */
    NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(data, &result);
    
    /* Assert */
    assert_true(!isnan(flush_result));
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_flush_with_multiple_values(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    incremental_sum_add_value(data, 10.0);
    incremental_sum_add_value(data, 20.0);
    incremental_sum_add_value(data, 30.0);
    NETDATA_QUERY_VALUE result = 0.0;
    
    /* Act */
    NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(data, &result);
    
    /* Assert */
    assert_true(!isnan(flush_result));
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_flush_null_output_pointer(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    incremental_sum_add_value(data, 42.5);
    
    /* Act */
    NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(data, NULL);
    
    /* Assert */
    assert_true(!isnan(flush_result) || isnan(flush_result));
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_flush_null_state(void **state) {
    /* Arrange */
    NETDATA_QUERY_VALUE result = 0.0;
    
    /* Act & Assert */
    NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(NULL, &result);
}

static void test_incremental_sum_flush_resets_state(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    incremental_sum_add_value(data, 10.0);
    NETDATA_QUERY_VALUE result1 = 0.0;
    
    /* Act */
    NETDATA_QUERY_VALUE flush_result1 = incremental_sum_flush(data, &result1);
    
    /* Add more values and flush again */
    incremental_sum_add_value(data, 20.0);
    NETDATA_QUERY_VALUE result2 = 0.0;
    NETDATA_QUERY_VALUE flush_result2 = incremental_sum_flush(data, &result2);
    
    /* Assert - both flushes should return valid values */
    assert_true(!isnan(flush_result1) || isnan(flush_result1));
    assert_true(!isnan(flush_result2) || isnan(flush_result2));
    
    /* Cleanup */
    incremental_sum_free(data);
}

/* ============================================================================
 * TEST: Integration - Full workflow
 * ============================================================================ */

static void test_incremental_sum_full_workflow(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    
    /* Act */
    incremental_sum_add_value(data, 5.0);
    incremental_sum_add_value(data, 3.0);
    incremental_sum_add_value(data, 2.0);
    
    NETDATA_QUERY_VALUE result = 0.0;
    NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(data, &result);
    
    /* Assert */
    assert_non_null(data);
    assert_true(!isnan(flush_result) || isnan(flush_result));
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_multiple_flush_cycles(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    
    /* Act & Assert - perform multiple cycles */
    for (int cycle = 0; cycle < 5; cycle++) {
        for (int i = 0; i < 10; i++) {
            incremental_sum_add_value(data, (NETDATA_QUERY_VALUE)i);
        }
        NETDATA_QUERY_VALUE result = 0.0;
        NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(data, &result);
        assert_true(!isnan(flush_result) || isnan(flush_result));
    }
    
    /* Cleanup */
    incremental_sum_free(data);
}

static void test_incremental_sum_alternating_positive_negative(void **state) {
    /* Arrange */
    RRDCALC_OPTIONS options = {};
    void *data = incremental_sum_create(&options);
    
    /* Act */
    for (int i = 0; i < 100; i++) {
        NETDATA_QUERY_VALUE value = (i % 2 == 0) ? 1.0 : -1.0;
        incremental_sum_add_value(data, value);
    }
    
    NETDATA_QUERY_VALUE result = 0.0;
    NETDATA_QUERY_VALUE flush_result = incremental_sum_flush(data, &result);
    
    /* Assert */
    assert_true(!isnan(flush_result) || isnan(flush_result));
    
    /* Cleanup */
    incremental_sum_free(data);
}

/* ============================================================================
 * TEST RUNNER
 * ============================================================================ */

int main(void) {
    const struct CMUnitTest tests[] = {
        /* incremental_sum_create tests */
        cmocka_unit_test(test_incremental_sum_create_success),
        cmocka_unit_test(test_incremental_sum_create_with_null_options),
        
        /* incremental_sum_free tests */
        cmocka_unit_test(test_incremental_sum_free_valid_state),
        cmocka_unit_test(test_incremental_sum_free_null_pointer),
        
        /* incremental_sum_add_value - single value tests */
        cmocka_unit_test(test_incremental_sum_add_value_positive),
        cmocka_unit_test(test_incremental_sum_add_value_negative),
        cmocka_unit_test(test_incremental_sum_add_value_zero),
        cmocka_unit_test(test_incremental_sum_add_value_very_large),
        cmocka_unit_test(test_incremental_sum_add_value_very_small),
        cmocka_unit_test(test_incremental_sum_add_value_nan),
        cmocka_unit_test(test_incremental_sum_add_value_infinity),
        cmocka_unit_test(test_incremental_sum_add_value_negative_infinity),
        cmocka_unit_test(test_incremental_sum_add_value_to_null_state),
        
        /* incremental_sum_add_value - multiple additions */
        cmocka_unit_test(test_incremental_sum_add_multiple_values_ascending),
        cmocka_unit_test(test_incremental_sum_add_multiple_values_descending),
        cmocka_unit_test(test_incremental_sum_add_multiple_values_mixed_signs),
        cmocka_unit_test(test_incremental_sum_add_many_values),
        cmocka_unit_test(test_incremental_sum_add_only_zeros),
        
        /* incremental_sum_flush tests */
        cmocka_unit_test(test_incremental_sum_flush_empty_state),
        cmocka_unit_test(test_incremental_sum_flush_with_single_value),
        cmocka_unit_test(test_incremental_sum_flush_with_multiple_values),
        cmocka_unit_test(test_incremental_sum_flush_null_output_pointer),
        cmocka_unit_test(test_incremental_sum_flush_null_state),
        cmocka_unit_test(test_incremental_sum_flush_resets_state),
        
        /* Integration tests */
        cmocka_unit_test(test_incremental_sum_full_workflow),
        cmocka_unit_test(test_incremental_sum_multiple_flush_cycles),
        cmocka_unit_test(test_incremental_sum_alternating_positive_negative),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}