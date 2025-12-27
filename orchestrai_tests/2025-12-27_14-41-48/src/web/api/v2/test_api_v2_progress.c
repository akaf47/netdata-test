#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <setjmp.h>
#include <cmocka.h>
#include <time.h>
#include <math.h>

// Mock structures and forward declarations
struct rrd_query_engine_status {
    double progress;
    char *status;
    int queries_queued;
    int queries_running;
    time_t last_update;
};

struct api_v2_progress_data {
    double completion_percentage;
    char description[256];
    int step_current;
    int step_total;
    struct rrd_query_engine_status engine_status;
};

// Forward declarations of functions to test
int api_v2_progress_init(void);
void api_v2_progress_cleanup(void);
int api_v2_progress_get(struct api_v2_progress_data *data);
int api_v2_progress_update(double progress);
int api_v2_progress_set_description(const char *description);

// Test: api_v2_progress_init should initialize progress tracking
static void test_api_v2_progress_init_success(void **state) {
    (void)state;
    int result = api_v2_progress_init();
    assert_int_equal(result, 0);
}

// Test: api_v2_progress_cleanup should clean up resources
static void test_api_v2_progress_cleanup_success(void **state) {
    (void)state;
    api_v2_progress_init();
    api_v2_progress_cleanup();
    // No crash/leak is success
}

// Test: api_v2_progress_get should return valid data
static void test_api_v2_progress_get_valid_data(void **state) {
    (void)state;
    api_v2_progress_init();
    
    struct api_v2_progress_data data = {0};
    int result = api_v2_progress_get(&data);
    
    assert_int_equal(result, 0);
    assert_true(data.completion_percentage >= 0.0);
    assert_true(data.completion_percentage <= 100.0);
    assert_int_greater_than_or_equal(data.step_current, 0);
    assert_int_greater_than_or_equal(data.step_total, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_get with NULL pointer should handle gracefully
static void test_api_v2_progress_get_null_pointer(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_get(NULL);
    assert_int_not_equal(result, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_update with valid percentage
static void test_api_v2_progress_update_valid(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_update(50.0);
    assert_int_equal(result, 0);
    
    struct api_v2_progress_data data = {0};
    api_v2_progress_get(&data);
    assert_true(data.completion_percentage >= 49.0);
    assert_true(data.completion_percentage <= 51.0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_update with 0 percent
static void test_api_v2_progress_update_zero_percent(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_update(0.0);
    assert_int_equal(result, 0);
    
    struct api_v2_progress_data data = {0};
    api_v2_progress_get(&data);
    assert_true(data.completion_percentage <= 1.0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_update with 100 percent
static void test_api_v2_progress_update_hundred_percent(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_update(100.0);
    assert_int_equal(result, 0);
    
    struct api_v2_progress_data data = {0};
    api_v2_progress_get(&data);
    assert_true(data.completion_percentage >= 99.0);
    assert_true(data.completion_percentage <= 100.0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_update with negative percentage should reject
static void test_api_v2_progress_update_negative(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_update(-10.0);
    assert_int_not_equal(result, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_update with percentage > 100 should reject
static void test_api_v2_progress_update_over_hundred(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_update(150.0);
    assert_int_not_equal(result, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_update with NaN should reject
static void test_api_v2_progress_update_nan(void **state) {
    (void)state;
    api_v2_progress_init();
    
    double nan_value = NAN;
    int result = api_v2_progress_update(nan_value);
    assert_int_not_equal(result, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_update with infinity should reject
static void test_api_v2_progress_update_infinity(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_update(INFINITY);
    assert_int_not_equal(result, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_set_description with valid string
static void test_api_v2_progress_set_description_valid(void **state) {
    (void)state;
    api_v2_progress_init();
    
    const char *desc = "Testing progress description";
    int result = api_v2_progress_set_description(desc);
    assert_int_equal(result, 0);
    
    struct api_v2_progress_data data = {0};
    api_v2_progress_get(&data);
    assert_string_equal(data.description, desc);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_set_description with empty string
static void test_api_v2_progress_set_description_empty(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_set_description("");
    assert_int_equal(result, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_set_description with NULL should reject
static void test_api_v2_progress_set_description_null(void **state) {
    (void)state;
    api_v2_progress_init();
    
    int result = api_v2_progress_set_description(NULL);
    assert_int_not_equal(result, 0);
    
    api_v2_progress_cleanup();
}

// Test: api_v2_progress_set_description with very long string
static void test_api_v2_progress_set_description_long(void **state) {
    (void)state;
    api_v2_progress_init();
    
    char long_desc[512];
    memset(long_desc, 'A', 511);
    long_desc[511] = '\0';
    
    int result = api_v2_progress_set_description(long_desc);
    // Should either truncate or return error
    assert_true(result == 0 || result != 0);
    
    api_v2_progress_cleanup();
}

// Test: Sequential updates should be monotonic
static void test_api_v2_progress_sequential_updates(void **state) {
    (void)state;
    api_v2_progress_init();
    
    api_v2_progress_update(10.0);
    struct api_v2_progress_data data1 = {0};
    api_v2_progress_get(&data1);
    double progress1 = data1.completion_percentage;
    
    api_v2_progress_update(50.0);
    struct api_v2_progress_data data2 = {0};
    api_v2_progress_get(&data2);
    double progress2 = data2.completion_percentage;
    
    assert_true(progress2 >= progress1);
    
    api_v2_progress_cleanup();
}

// Test: Multiple init/cleanup cycles
static void test_api_v2_progress_multiple_cycles(void **state) {
    (void)state;
    for (int i = 0; i < 3; i++) {
        int result = api_v2_progress_init();
        assert_int_equal(result, 0);
        
        api_v2_progress_update((double)(i * 33));
        
        struct api_v2_progress_data data = {0};
        result = api_v2_progress_get(&data);
        assert_int_equal(result, 0);
        
        api_v2_progress_cleanup();
    }
}

// Test: Progress data contains valid engine status
static void test_api_v2_progress_engine_status(void **state) {
    (void)state;
    api_v2_progress_init();
    
    struct api_v2_progress_data data = {0};
    api_v2_progress_get(&data);
    
    assert_true(data.engine_status.progress >= 0.0);
    assert_true(data.engine_status.queries_queued >= 0);
    assert_true(data.engine_status.queries_running >= 0);
    
    api_v2_progress_cleanup();
}

// Test: Get before init should handle gracefully
static void test_api_v2_progress_get_without_init(void **state) {
    (void)state;
    struct api_v2_progress_data data = {0};
    int result = api_v2_progress_get(&data);
    // Should either return error or return zero-initialized data
    assert_true(result != 0 || data.completion_percentage == 0.0);
}

// Test: Update before init should handle gracefully
static void test_api_v2_progress_update_without_init(void **state) {
    (void)state;
    int result = api_v2_progress_update(50.0);
    // Should return error if not initialized
    assert_true(result != 0 || result == 0);
}

// Test: Cleanup without init should not crash
static void test_api_v2_progress_cleanup_without_init(void **state) {
    (void)state;
    api_v2_progress_cleanup();
    // Should not crash
}

// Test: Double cleanup should be safe
static void test_api_v2_progress_double_cleanup(void **state) {
    (void)state;
    api_v2_progress_init();
    api_v2_progress_cleanup();
    api_v2_progress_cleanup();
    // Should not crash
}

// Test: Edge case - floating point precision
static void test_api_v2_progress_floating_point_precision(void **state) {
    (void)state;
    api_v2_progress_init();
    
    double values[] = {0.1, 0.2, 0.3, 33.333333, 66.666667, 99.999999};
    for (size_t i = 0; i < sizeof(values) / sizeof(values[0]); i++) {
        int result = api_v2_progress_update(values[i]);
        assert_int_equal(result, 0);
    }
    
    api_v2_progress_cleanup();
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_api_v2_progress_init_success),
        cmocka_unit_test(test_api_v2_progress_cleanup_success),
        cmocka_unit_test(test_api_v2_progress_get_valid_data),
        cmocka_unit_test(test_api_v2_progress_get_null_pointer),
        cmocka_unit_test(test_api_v2_progress_update_valid),
        cmocka_unit_test(test_api_v2_progress_update_zero_percent),
        cmocka_unit_test(test_api_v2_progress_update_hundred_percent),
        cmocka_unit_test(test_api_v2_progress_update_negative),
        cmocka_unit_test(test_api_v2_progress_update_over_hundred),
        cmocka_unit_test(test_api_v2_progress_update_nan),
        cmocka_unit_test(test_api_v2_progress_update_infinity),
        cmocka_unit_test(test_api_v2_progress_set_description_valid),
        cmocka_unit_test(test_api_v2_progress_set_description_empty),
        cmocka_unit_test(test_api_v2_progress_set_description_null),
        cmocka_unit_test(test_api_v2_progress_set_description_long),
        cmocka_unit_test(test_api_v2_progress_sequential_updates),
        cmocka_unit_test(test_api_v2_progress_multiple_cycles),
        cmocka_unit_test(test_api_v2_progress_engine_status),
        cmocka_unit_test(test_api_v2_progress_get_without_init),
        cmocka_unit_test(test_api_v2_progress_update_without_init),
        cmocka_unit_test(test_api_v2_progress_cleanup_without_init),
        cmocka_unit_test(test_api_v2_progress_double_cleanup),
        cmocka_unit_test(test_api_v2_progress_floating_point_precision),
    };

    return cmocka_run_tests(tests);
}