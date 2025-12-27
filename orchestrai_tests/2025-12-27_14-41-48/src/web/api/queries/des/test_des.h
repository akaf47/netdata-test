#ifndef TEST_DES_H
#define TEST_DES_H

#include <stddef.h>
#include <stdint.h>
#include <math.h>

/**
 * DES (Double Exponential Smoothing) Test Suite Header
 * 
 * This header provides test utilities and mock declarations for testing
 * the DES query implementation in netdata.
 */

/* Test assertion macros */
#define TEST_ASSERT_EQUAL_DOUBLE(actual, expected, tolerance, msg) \
    assert_equal_double(actual, expected, tolerance, msg)

#define TEST_ASSERT_EQUAL_INT(actual, expected, msg) \
    assert_equal_int(actual, expected, msg)

#define TEST_ASSERT_NULL(ptr, msg) \
    assert_equal_pointer(ptr, NULL, msg)

#define TEST_ASSERT_NOT_NULL(ptr, msg) \
    assert_not_null(ptr, msg)

#define TEST_ASSERT_TRUE(cond, msg) \
    assert_true(cond, msg)

#define TEST_ASSERT_FALSE(cond, msg) \
    assert_true(!(cond), msg)

#define TEST_ASSERT_NAN(val, msg) \
    assert_true(isnan(val), msg)

#define TEST_ASSERT_INF(val, msg) \
    assert_true(isinf(val), msg)

#define TEST_ASSERT_FINITE(val, msg) \
    assert_true(isfinite(val), msg)

/* Mock DES state structure */
typedef struct {
    double alpha;              /* Exponential smoothing parameter */
    double beta;               /* Trend smoothing parameter */
    double s1;                 /* Level component */
    double s2;                 /* Trend component */
    int initialized;           /* Initialization flag */
    size_t count;              /* Sample count */
} DES_STATE;

/* Forward declarations of test functions */
void assert_equal_double(double actual, double expected, double tolerance, const char *msg);
void assert_equal_int(int actual, int expected, const char *msg);
void assert_equal_pointer(void *actual, void *expected, const char *msg);
void assert_not_null(void *ptr, const char *msg);
void assert_true(int condition, const char *msg);

/* DES Operation Tests */
void test_des_init_null_state(void);
void test_des_init_valid_state(void);
void test_des_init_boundary_alpha_zero(void);
void test_des_init_boundary_alpha_one(void);
void test_des_init_boundary_beta_zero(void);
void test_des_init_boundary_beta_one(void);
void test_des_init_negative_alpha(void);
void test_des_init_negative_beta(void);
void test_des_init_large_alpha(void);
void test_des_init_large_beta(void);

void test_des_apply_null_state(void);
void test_des_apply_nan_value(void);
void test_des_apply_first_value(void);
void test_des_apply_second_value_increase(void);
void test_des_apply_second_value_decrease(void);
void test_des_apply_zero_value(void);
void test_des_apply_negative_value(void);
void test_des_apply_very_large_value(void);
void test_des_apply_very_small_value(void);
void test_des_apply_infinite_value(void);
void test_des_apply_negative_infinite_value(void);
void test_des_apply_multiple_values_sequence(void);
void test_des_apply_alpha_zero_behavior(void);
void test_des_apply_alpha_one_behavior(void);
void test_des_apply_beta_zero_behavior(void);
void test_des_apply_beta_one_behavior(void);
void test_des_apply_fractional_alpha(void);
void test_des_apply_identical_values(void);
void test_des_apply_oscillating_values(void);

/* State Management Tests */
void test_des_state_persistence_across_calls(void);
void test_des_apply_count_increment(void);
void test_des_apply_count_overflow_handling(void);
void test_des_multiple_independent_states(void);
void test_des_reinitialization(void);

/* Integration Tests */
void test_des_full_workflow(void);
void test_des_trend_detection(void);
void test_des_trend_detection_downtrend(void);
void test_des_no_division_by_zero(void);
void test_des_memory_safety_bounds(void);

#endif /* TEST_DES_H */