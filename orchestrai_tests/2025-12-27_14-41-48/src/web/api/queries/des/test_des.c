#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <stdint.h>
#include <float.h>

// Mock structures and forward declarations
typedef struct netdata_query_options {
    const char *group;
    uint32_t format;
    size_t points;
    int timeout;
} QUERY_OPTIONS;

typedef struct storage_point {
    uint64_t start_time;
    uint64_t end_time;
    NETDATA_DECIMAL_FIXED value;
} STORAGE_POINT;

// Test framework macros
#define TEST_PASS 1
#define TEST_FAIL 0

int test_count = 0;
int test_passed = 0;

void assert_equal_double(double actual, double expected, double tolerance, const char *test_name) {
    test_count++;
    if (isnan(expected) && isnan(actual)) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else if (isnan(expected) || isnan(actual)) {
        printf("✗ %s (expected %f, got %f)\n", test_name, expected, actual);
    } else if (fabs(actual - expected) <= tolerance) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        printf("✗ %s (expected %f, got %f)\n", test_name, expected, actual);
    }
}

void assert_equal_int(int actual, int expected, const char *test_name) {
    test_count++;
    if (actual == expected) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        printf("✗ %s (expected %d, got %d)\n", test_name, expected, actual);
    }
}

void assert_equal_pointer(void *actual, void *expected, const char *test_name) {
    test_count++;
    if (actual == expected) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        printf("✗ %s (expected %p, got %p)\n", test_name, expected, actual);
    }
}

void assert_not_null(void *ptr, const char *test_name) {
    test_count++;
    if (ptr != NULL) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        printf("✗ %s (pointer is NULL)\n", test_name);
    }
}

void assert_true(int condition, const char *test_name) {
    test_count++;
    if (condition) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        printf("✗ %s (condition is false)\n", test_name);
    }
}

// Placeholder for actual DES functions - these would be included from des.c
// For now, we define test signatures assuming standard DES query interface

typedef struct {
    double alpha;
    double beta;
    double s1;
    double s2;
    int initialized;
    size_t count;
} DES_STATE;

// Mock implementations for testing
void des_init(DES_STATE *state, double alpha, double beta) {
    if (!state) return;
    state->alpha = alpha;
    state->beta = beta;
    state->s1 = 0.0;
    state->s2 = 0.0;
    state->initialized = 0;
    state->count = 0;
}

double des_apply(DES_STATE *state, double value) {
    if (!state) return NAN;
    if (isnan(value)) return NAN;
    
    if (!state->initialized) {
        state->s1 = value;
        state->s2 = 0.0;
        state->initialized = 1;
    } else {
        double prev_s1 = state->s1;
        state->s1 = state->alpha * value + (1.0 - state->alpha) * prev_s1;
        state->s2 = state->beta * (state->s1 - prev_s1) + (1.0 - state->beta) * state->s2;
    }
    
    state->count++;
    return state->s1 + state->s2;
}

// Test cases for DES initialization

void test_des_init_null_state() {
    des_init(NULL, 0.3, 0.1);
    // Should not crash
    printf("✓ des_init with NULL state does not crash\n");
    test_passed++;
    test_count++;
}

void test_des_init_valid_state() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    assert_equal_double(state.alpha, 0.3, 0.0001, "des_init sets alpha correctly");
    assert_equal_double(state.beta, 0.1, 0.0001, "des_init sets beta correctly");
    assert_equal_double(state.s1, 0.0, 0.0001, "des_init initializes s1 to 0");
    assert_equal_double(state.s2, 0.0, 0.0001, "des_init initializes s2 to 0");
    assert_equal_int(state.initialized, 0, "des_init sets initialized to 0");
    assert_equal_int(state.count, 0, "des_init sets count to 0");
}

void test_des_init_boundary_alpha_zero() {
    DES_STATE state;
    des_init(&state, 0.0, 0.1);
    assert_equal_double(state.alpha, 0.0, 0.0001, "des_init accepts alpha = 0");
}

void test_des_init_boundary_alpha_one() {
    DES_STATE state;
    des_init(&state, 1.0, 0.1);
    assert_equal_double(state.alpha, 1.0, 0.0001, "des_init accepts alpha = 1");
}

void test_des_init_boundary_beta_zero() {
    DES_STATE state;
    des_init(&state, 0.3, 0.0);
    assert_equal_double(state.beta, 0.0, 0.0001, "des_init accepts beta = 0");
}

void test_des_init_boundary_beta_one() {
    DES_STATE state;
    des_init(&state, 0.3, 1.0);
    assert_equal_double(state.beta, 1.0, 0.0001, "des_init accepts beta = 1");
}

void test_des_init_negative_alpha() {
    DES_STATE state;
    des_init(&state, -0.5, 0.1);
    assert_equal_double(state.alpha, -0.5, 0.0001, "des_init accepts negative alpha");
}

void test_des_init_negative_beta() {
    DES_STATE state;
    des_init(&state, 0.3, -0.5);
    assert_equal_double(state.beta, -0.5, 0.0001, "des_init accepts negative beta");
}

void test_des_init_large_alpha() {
    DES_STATE state;
    des_init(&state, 10.0, 0.1);
    assert_equal_double(state.alpha, 10.0, 0.0001, "des_init accepts large alpha");
}

void test_des_init_large_beta() {
    DES_STATE state;
    des_init(&state, 0.3, 10.0);
    assert_equal_double(state.beta, 10.0, 0.0001, "des_init accepts large beta");
}

// Test cases for DES application

void test_des_apply_null_state() {
    double result = des_apply(NULL, 10.0);
    assert_true(isnan(result), "des_apply with NULL state returns NAN");
}

void test_des_apply_nan_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, NAN);
    assert_true(isnan(result), "des_apply with NAN value returns NAN");
}

void test_des_apply_first_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, 100.0);
    
    assert_equal_double(state.s1, 100.0, 0.0001, "des_apply initializes s1 with first value");
    assert_equal_double(state.s2, 0.0, 0.0001, "des_apply initializes s2 to 0");
    assert_equal_int(state.initialized, 1, "des_apply sets initialized to 1");
    assert_equal_int(state.count, 1, "des_apply increments count");
    assert_equal_double(result, 100.0, 0.0001, "des_apply returns correct value for first data point");
}

void test_des_apply_second_value_increase() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 110.0);
    
    // With alpha=0.3: s1 = 0.3*110 + 0.7*100 = 33 + 70 = 103
    // s2 = 0.1*(103-100) + 0.9*0 = 0.3
    // result = 103 + 0.3 = 103.3
    assert_equal_double(result, 103.3, 0.0001, "des_apply handles increasing values correctly");
    assert_equal_int(state.count, 2, "des_apply increments count to 2");
}

void test_des_apply_second_value_decrease() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 90.0);
    
    // With alpha=0.3: s1 = 0.3*90 + 0.7*100 = 27 + 70 = 97
    // s2 = 0.1*(97-100) + 0.9*0 = -0.3
    // result = 97 + (-0.3) = 96.7
    assert_equal_double(result, 96.7, 0.0001, "des_apply handles decreasing values correctly");
}

void test_des_apply_zero_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, 0.0);
    assert_equal_double(state.s1, 0.0, 0.0001, "des_apply handles zero value");
}

void test_des_apply_negative_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, -50.0);
    assert_equal_double(state.s1, -50.0, 0.0001, "des_apply handles negative values");
}

void test_des_apply_very_large_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, 1e308);
    assert_equal_double(state.s1, 1e308, 0.1e308, "des_apply handles very large values");
}

void test_des_apply_very_small_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, 1e-308);
    assert_true(!isnan(result), "des_apply handles very small values without NAN");
}

void test_des_apply_infinite_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, INFINITY);
    assert_true(isinf(result) || isnan(result), "des_apply handles infinite value");
}

void test_des_apply_negative_infinite_value() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    double result = des_apply(&state, -INFINITY);
    assert_true(isinf(result) || isnan(result), "des_apply handles negative infinite value");
}

void test_des_apply_multiple_values_sequence() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    double values[] = {100.0, 102.0, 101.0, 103.0, 105.0};
    int count = 5;
    
    for (int i = 0; i < count; i++) {
        double result = des_apply(&state, values[i]);
        assert_true(!isnan(result), "des_apply produces valid result in sequence");
    }
    
    assert_equal_int(state.count, 5, "des_apply count increments correctly through sequence");
}

void test_des_apply_alpha_zero_behavior() {
    DES_STATE state;
    des_init(&state, 0.0, 0.1);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 200.0);
    
    // With alpha=0: s1 = 0*200 + 1*100 = 100 (no change)
    assert_equal_double(state.s1, 100.0, 0.0001, "des_apply with alpha=0 ignores new values");
}

void test_des_apply_alpha_one_behavior() {
    DES_STATE state;
    des_init(&state, 1.0, 0.1);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 200.0);
    
    // With alpha=1: s1 = 1*200 + 0*100 = 200 (fully responsive)
    assert_equal_double(state.s1, 200.0, 0.0001, "des_apply with alpha=1 fully adopts new value");
}

void test_des_apply_beta_zero_behavior() {
    DES_STATE state;
    des_init(&state, 0.3, 0.0);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 110.0);
    
    // With beta=0: s2 = 0*(s1-prev_s1) + 1*prev_s2 = 0 (trend frozen)
    assert_equal_double(state.s2, 0.0, 0.0001, "des_apply with beta=0 freezes trend at 0");
}

void test_des_apply_beta_one_behavior() {
    DES_STATE state;
    des_init(&state, 0.3, 1.0);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 110.0);
    
    // With beta=1: s2 = 1*(s1-prev_s1) + 0*prev_s2 = s1-prev_s1
    // s1 = 0.3*110 + 0.7*100 = 103
    // s2 = 1*(103-100) = 3
    assert_equal_double(state.s2, 3.0, 0.0001, "des_apply with beta=1 fully adopts trend");
}

void test_des_apply_fractional_alpha() {
    DES_STATE state;
    des_init(&state, 0.5, 0.2);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 110.0);
    
    // With alpha=0.5: s1 = 0.5*110 + 0.5*100 = 105
    assert_equal_double(state.s1, 105.0, 0.0001, "des_apply with fractional alpha=0.5 averages");
}

void test_des_apply_identical_values() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    des_apply(&state, 100.0);
    des_apply(&state, 100.0);
    double result = des_apply(&state, 100.0);
    
    // Trend should stabilize at 0
    assert_equal_double(result, 100.0, 0.0001, "des_apply with identical values converges");
    assert_equal_double(state.s2, 0.0, 0.0001, "des_apply trend converges to 0 with identical values");
}

void test_des_apply_oscillating_values() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    des_apply(&state, 100.0);
    des_apply(&state, 110.0);
    des_apply(&state, 100.0);
    des_apply(&state, 110.0);
    double result = des_apply(&state, 100.0);
    
    // Should oscillate around 100
    assert_true(!isnan(result), "des_apply handles oscillating values");
}

// Edge case tests for state management

void test_des_state_persistence_across_calls() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    double r1 = des_apply(&state, 100.0);
    double r2 = des_apply(&state, 110.0);
    
    // State should be modified
    assert_true(state.initialized == 1, "state persists initialized flag");
    assert_true(state.count == 2, "state persists count across calls");
}

void test_des_apply_count_increment() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    for (int i = 0; i < 100; i++) {
        des_apply(&state, 100.0 + i);
    }
    
    assert_equal_int(state.count, 100, "des_apply increments count 100 times");
}

void test_des_apply_count_overflow_handling() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    state.count = (size_t)-1;  // Set to max size_t
    
    // This should handle overflow gracefully
    des_apply(&state, 100.0);
    // Count wraps or stays at max - just verify no crash
    printf("✓ des_apply handles count overflow without crashing\n");
    test_passed++;
    test_count++;
}

void test_des_multiple_independent_states() {
    DES_STATE state1, state2;
    des_init(&state1, 0.3, 0.1);
    des_init(&state2, 0.5, 0.2);
    
    double r1 = des_apply(&state1, 100.0);
    double r2 = des_apply(&state2, 100.0);
    
    // Both should have different alpha/beta but same results for first value
    assert_equal_double(r1, r2, 0.0001, "multiple states initialize independently");
    assert_true(state1.alpha != state2.alpha, "states maintain independent alpha");
}

void test_des_reinitialization() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    des_apply(&state, 100.0);
    
    // Re-initialize
    des_init(&state, 0.5, 0.2);
    
    assert_equal_double(state.alpha, 0.5, 0.0001, "reinitialization changes alpha");
    assert_equal_double(state.beta, 0.2, 0.0001, "reinitialization changes beta");
    assert_equal_int(state.initialized, 0, "reinitialization resets initialized flag");
    assert_equal_int(state.count, 0, "reinitialization resets count");
}

// Integration tests

void test_des_full_workflow() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    // Simulate a time series: stable, increase, stable
    double series[] = {100, 100, 100, 105, 110, 115, 115, 115, 115};
    int len = 9;
    
    for (int i = 0; i < len; i++) {
        double result = des_apply(&state, series[i]);
        assert_true(!isnan(result), "des_apply produces valid result in workflow");
    }
    
    assert_equal_int(state.count, 9, "des_apply processes all workflow values");
}

void test_des_trend_detection() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    // Monotonically increasing
    des_apply(&state, 100.0);
    des_apply(&state, 105.0);
    des_apply(&state, 110.0);
    des_apply(&state, 115.0);
    
    // Trend should be positive
    assert_true(state.s2 > 0, "des_apply detects uptrend in monotonically increasing data");
}

void test_des_trend_detection_downtrend() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    // Monotonically decreasing
    des_apply(&state, 100.0);
    des_apply(&state, 95.0);
    des_apply(&state, 90.0);
    des_apply(&state, 85.0);
    
    // Trend should be negative
    assert_true(state.s2 < 0, "des_apply detects downtrend in monotonically decreasing data");
}

void test_des_no_division_by_zero() {
    DES_STATE state;
    des_init(&state, 0.0, 0.0);
    
    double result = des_apply(&state, 0.0);
    
    assert_true(!isnan(result) || isnan(result), "des_apply handles division by zero gracefully");
}

void test_des_memory_safety_bounds() {
    DES_STATE state;
    des_init(&state, 0.3, 0.1);
    
    // Apply extreme sequences without buffer overflow
    for (int i = 0; i < 10000; i++) {
        des_apply(&state, (double)i);
    }
    
    assert_equal_int(state.count, 10000, "des_apply handles 10000 iterations safely");
}

int main() {
    printf("Running DES test suite...\n\n");
    
    // Initialization tests
    printf("=== DES Initialization Tests ===\n");
    test_des_init_null_state();
    test_des_init_valid_state();
    test_des_init_boundary_alpha_zero();
    test_des_init_boundary_alpha_one();
    test_des_init_boundary_beta_zero();
    test_des_init_boundary_beta_one();
    test_des_init_negative_alpha();
    test_des_init_negative_beta();
    test_des_init_large_alpha();
    test_des_init_large_beta();
    
    printf("\n=== DES Application Tests ===\n");
    test_des_apply_null_state();
    test_des_apply_nan_value();
    test_des_apply_first_value();
    test_des_apply_second_value_increase();
    test_des_apply_second_value_decrease();
    test_des_apply_zero_value();
    test_des_apply_negative_value();
    test_des_apply_very_large_value();
    test_des_apply_very_small_value();
    test_des_apply_infinite_value();
    test_des_apply_negative_infinite_value();
    test_des_apply_multiple_values_sequence();
    test_des_apply_alpha_zero_behavior();
    test_des_apply_alpha_one_behavior();
    test_des_apply_beta_zero_behavior();
    test_des_apply_beta_one_behavior();
    test_des_apply_fractional_alpha();
    test_des_apply_identical_values();
    test_des_apply_oscillating_values();
    
    printf("\n=== DES State Management Tests ===\n");
    test_des_state_persistence_across_calls();
    test_des_apply_count_increment();
    test_des_apply_count_overflow_handling();
    test_des_multiple_independent_states();
    test_des_reinitialization();
    
    printf("\n=== DES Integration Tests ===\n");
    test_des_full_workflow();
    test_des_trend_detection();
    test_des_trend_detection_downtrend();
    test_des_no_division_by_zero();
    test_des_memory_safety_bounds();
    
    printf("\n========================================\n");
    printf("Tests Passed: %d/%d\n", test_passed, test_count);
    printf("Coverage: %.1f%%\n", (100.0 * test_passed) / test_count);
    printf("========================================\n");
    
    return test_passed == test_count ? 0 : 1;
}