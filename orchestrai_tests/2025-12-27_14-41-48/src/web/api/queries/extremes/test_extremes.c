#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <setjmp.h>
#include <cmocka.h>

/* Mock structures and forward declarations */
struct query_result {
    double min;
    double max;
    int count;
};

/* Assuming the extremes.h defines the following */
typedef struct {
    double min;
    double max;
    int initialized;
} extremes_t;

/* Function prototypes (from extremes.c) */
extern void extremes_create(extremes_t *e);
extern void extremes_reset(extremes_t *e);
extern void extremes_add(extremes_t *e, double value);
extern double extremes_min(extremes_t *e);
extern double extremes_max(extremes_t *e);
extern int extremes_count(extremes_t *e);

/* Test: extremes_create - initialization */
static void test_extremes_create_initializes_structure(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    assert_int_equal(extremes.initialized, 1);
    assert_true(isnan(extremes.min) || extremes.min == 0);
    assert_true(isnan(extremes.max) || extremes.max == 0);
}

/* Test: extremes_reset - resets state */
static void test_extremes_reset_clears_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    extremes_add(&extremes, 5.0);
    
    extremes_reset(&extremes);
    
    assert_true(isnan(extremes.min) || extremes.min == 0);
    assert_true(isnan(extremes.max) || extremes.max == 0);
}

/* Test: extremes_add - single positive value */
static void test_extremes_add_single_positive_value(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 5.0);
    
    assert_double_equal(extremes_min(&extremes), 5.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
}

/* Test: extremes_add - single negative value */
static void test_extremes_add_single_negative_value(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, -5.0);
    
    assert_double_equal(extremes_min(&extremes), -5.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), -5.0, 0.0001);
}

/* Test: extremes_add - zero value */
static void test_extremes_add_zero_value(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 0.0);
    
    assert_double_equal(extremes_min(&extremes), 0.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 0.0, 0.0001);
}

/* Test: extremes_add - multiple values ascending */
static void test_extremes_add_multiple_values_ascending(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 1.0);
    extremes_add(&extremes, 2.0);
    extremes_add(&extremes, 3.0);
    extremes_add(&extremes, 4.0);
    extremes_add(&extremes, 5.0);
    
    assert_double_equal(extremes_min(&extremes), 1.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
    assert_int_equal(extremes_count(&extremes), 5);
}

/* Test: extremes_add - multiple values descending */
static void test_extremes_add_multiple_values_descending(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 5.0);
    extremes_add(&extremes, 4.0);
    extremes_add(&extremes, 3.0);
    extremes_add(&extremes, 2.0);
    extremes_add(&extremes, 1.0);
    
    assert_double_equal(extremes_min(&extremes), 1.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
}

/* Test: extremes_add - multiple values unordered */
static void test_extremes_add_multiple_values_unordered(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 3.0);
    extremes_add(&extremes, 1.0);
    extremes_add(&extremes, 5.0);
    extremes_add(&extremes, 2.0);
    extremes_add(&extremes, 4.0);
    
    assert_double_equal(extremes_min(&extremes), 1.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
}

/* Test: extremes_add - mixed positive and negative */
static void test_extremes_add_mixed_positive_negative(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, -5.0);
    extremes_add(&extremes, 3.0);
    extremes_add(&extremes, -2.0);
    extremes_add(&extremes, 7.0);
    
    assert_double_equal(extremes_min(&extremes), -5.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 7.0, 0.0001);
}

/* Test: extremes_add - duplicate values */
static void test_extremes_add_duplicate_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 5.0);
    extremes_add(&extremes, 5.0);
    extremes_add(&extremes, 5.0);
    
    assert_double_equal(extremes_min(&extremes), 5.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
    assert_int_equal(extremes_count(&extremes), 3);
}

/* Test: extremes_add - very small values (near zero) */
static void test_extremes_add_very_small_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 0.0001);
    extremes_add(&extremes, 0.0002);
    extremes_add(&extremes, 0.00005);
    
    assert_double_equal(extremes_min(&extremes), 0.00005, 0.000001);
    assert_double_equal(extremes_max(&extremes), 0.0002, 0.000001);
}

/* Test: extremes_add - very large values */
static void test_extremes_add_very_large_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 1e10);
    extremes_add(&extremes, 1e15);
    extremes_add(&extremes, 1e12);
    
    assert_double_equal(extremes_min(&extremes), 1e10, 1e9);
    assert_double_equal(extremes_max(&extremes), 1e15, 1e14);
}

/* Test: extremes_add - negative large values */
static void test_extremes_add_negative_large_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, -1e10);
    extremes_add(&extremes, -1e15);
    extremes_add(&extremes, -1e12);
    
    assert_double_equal(extremes_min(&extremes), -1e15, 1e14);
    assert_double_equal(extremes_max(&extremes), -1e10, 1e9);
}

/* Test: extremes_add - NaN handling if supported */
static void test_extremes_add_nan_value(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 5.0);
    extremes_add(&extremes, NAN);
    extremes_add(&extremes, 3.0);
    
    assert_double_equal(extremes_min(&extremes), 3.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
}

/* Test: extremes_add - infinity values */
static void test_extremes_add_infinity_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, INFINITY);
    extremes_add(&extremes, 5.0);
    
    assert_double_equal(extremes_max(&extremes), INFINITY, 0.0001);
    assert_double_equal(extremes_min(&extremes), 5.0, 0.0001);
}

/* Test: extremes_add - negative infinity */
static void test_extremes_add_negative_infinity(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, -INFINITY);
    extremes_add(&extremes, 5.0);
    
    assert_double_equal(extremes_min(&extremes), -INFINITY, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
}

/* Test: extremes_min - uninitialized structure */
static void test_extremes_min_uninitialized(void **state) {
    extremes_t extremes;
    memset(&extremes, 0, sizeof(extremes_t));
    
    double min = extremes_min(&extremes);
    assert_true(isnan(min) || min == 0);
}

/* Test: extremes_max - uninitialized structure */
static void test_extremes_max_uninitialized(void **state) {
    extremes_t extremes;
    memset(&extremes, 0, sizeof(extremes_t));
    
    double max = extremes_max(&extremes);
    assert_true(isnan(max) || max == 0);
}

/* Test: extremes_count - multiple values */
static void test_extremes_count_multiple_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    for (int i = 0; i < 100; i++) {
        extremes_add(&extremes, (double)i);
    }
    
    assert_int_equal(extremes_count(&extremes), 100);
}

/* Test: extremes_count - zero values */
static void test_extremes_count_zero_values(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    assert_int_equal(extremes_count(&extremes), 0);
}

/* Test: reset after adding values */
static void test_extremes_reset_after_multiple_adds(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 10.0);
    extremes_add(&extremes, 20.0);
    extremes_add(&extremes, 30.0);
    
    extremes_reset(&extremes);
    
    extremes_add(&extremes, 5.0);
    
    assert_double_equal(extremes_min(&extremes), 5.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), 5.0, 0.0001);
    assert_int_equal(extremes_count(&extremes), 1);
}

/* Test: alternating reset and add */
static void test_extremes_alternating_reset_add(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    for (int cycle = 0; cycle < 5; cycle++) {
        extremes_reset(&extremes);
        extremes_add(&extremes, 1.0 * (cycle + 1));
        
        assert_int_equal(extremes_count(&extremes), 1);
        assert_double_equal(extremes_min(&extremes), 1.0 * (cycle + 1), 0.0001);
    }
}

/* Test: large number of additions */
static void test_extremes_large_number_of_additions(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    int count = 10000;
    for (int i = 0; i < count; i++) {
        extremes_add(&extremes, (double)i);
    }
    
    assert_int_equal(extremes_count(&extremes), count);
    assert_double_equal(extremes_min(&extremes), 0.0, 0.0001);
    assert_double_equal(extremes_max(&extremes), (double)(count - 1), 1.0);
}

/* Test: floating point precision edge case */
static void test_extremes_floating_point_precision(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 0.1 + 0.2);
    extremes_add(&extremes, 0.3);
    
    assert_double_equal(extremes_min(&extremes), 0.3, 0.0001);
    assert_double_equal(extremes_max(&extremes), 0.3, 0.0001);
}

/* Test: single value then reset then single value */
static void test_extremes_single_reset_single(void **state) {
    extremes_t extremes;
    extremes_create(&extremes);
    
    extremes_add(&extremes, 10.0);
    assert_double_equal(extremes_min(&extremes), 10.0, 0.0001);
    assert_int_equal(extremes_count(&extremes), 1);
    
    extremes_reset(&extremes);
    assert_int_equal(extremes_count(&extremes), 0);
    
    extremes_add(&extremes, 20.0);
    assert_double_equal(extremes_min(&extremes), 20.0, 0.0001);
    assert_int_equal(extremes_count(&extremes), 1);
}

/* Main test runner */
int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_extremes_create_initializes_structure),
        cmocka_unit_test(test_extremes_reset_clears_values),
        cmocka_unit_test(test_extremes_add_single_positive_value),
        cmocka_unit_test(test_extremes_add_single_negative_value),
        cmocka_unit_test(test_extremes_add_zero_value),
        cmocka_unit_test(test_extremes_add_multiple_values_ascending),
        cmocka_unit_test(test_extremes_add_multiple_values_descending),
        cmocka_unit_test(test_extremes_add_multiple_values_unordered),
        cmocka_unit_test(test_extremes_add_mixed_positive_negative),
        cmocka_unit_test(test_extremes_add_duplicate_values),
        cmocka_unit_test(test_extremes_add_very_small_values),
        cmocka_unit_test(test_extremes_add_very_large_values),
        cmocka_unit_test(test_extremes_add_negative_large_values),
        cmocka_unit_test(test_extremes_add_nan_value),
        cmocka_unit_test(test_extremes_add_infinity_values),
        cmocka_unit_test(test_extremes_add_negative_infinity),
        cmocka_unit_test(test_extremes_min_uninitialized),
        cmocka_unit_test(test_extremes_max_uninitialized),
        cmocka_unit_test(test_extremes_count_multiple_values),
        cmocka_unit_test(test_extremes_count_zero_values),
        cmocka_unit_test(test_extremes_reset_after_multiple_adds),
        cmocka_unit_test(test_extremes_alternating_reset_add),
        cmocka_unit_test(test_extremes_large_number_of_additions),
        cmocka_unit_test(test_extremes_floating_point_precision),
        cmocka_unit_test(test_extremes_single_reset_single),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}