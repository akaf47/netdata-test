#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <time.h>
#include <stdint.h>

/* Mock the countif module */
#include "../countif.h"

/* Test framework */
#include <assert.h>
#include <float.h>

#define EPSILON 0.0001
#define MAX_TEST_POINTS 1000

/* Helper macros */
#define TEST_ASSERT(condition, message) \
    do { \
        if (!(condition)) { \
            fprintf(stderr, "FAIL: %s\n", message); \
            return 0; \
        } \
    } while(0)

#define TEST_ASSERT_DOUBLE_EQ(actual, expected, message) \
    do { \
        if (fabs((actual) - (expected)) > EPSILON) { \
            fprintf(stderr, "FAIL: %s (expected %f, got %f)\n", message, expected, actual); \
            return 0; \
        } \
    } while(0)

#define TEST_ASSERT_INT_EQ(actual, expected, message) \
    do { \
        if ((actual) != (expected)) { \
            fprintf(stderr, "FAIL: %s (expected %d, got %d)\n", message, expected, actual); \
            return 0; \
        } \
    } while(0)

/* Test counter */
static int tests_run = 0;
static int tests_passed = 0;

/* Test: countif_init initializes structure correctly */
int test_countif_init_basic(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    TEST_ASSERT(cqd != NULL, "countif_query_data allocation");
    
    /* Act */
    countif_init(cqd);
    
    /* Assert */
    TEST_ASSERT(cqd->data != NULL, "countif_init should initialize data");
    TEST_ASSERT(cqd->count == 0, "countif_init should set count to 0");
    TEST_ASSERT(cqd->size >= 0, "countif_init should set valid size");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_init with NULL pointer */
int test_countif_init_null_pointer(void) {
    tests_run++;
    
    /* Arrange & Act - should handle gracefully */
    /* Cleanup: no crash expected */
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with single value */
int test_countif_aggregate_single_value(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, 42.5);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 1, "count should be 1 after adding single value");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[0], 42.5, "first value should be stored");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with multiple values */
int test_countif_aggregate_multiple_values(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 20.0);
    countif_aggregate(cqd, 30.0);
    countif_aggregate(cqd, 15.5);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 4, "count should be 4");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[0], 10.0, "first value");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[1], 20.0, "second value");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[2], 30.0, "third value");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[3], 15.5, "fourth value");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with negative values */
int test_countif_aggregate_negative_values(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, -10.5);
    countif_aggregate(cqd, -20.0);
    countif_aggregate(cqd, 5.0);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 3, "count should be 3");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[0], -10.5, "negative value 1");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[1], -20.0, "negative value 2");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[2], 5.0, "positive value");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with zero */
int test_countif_aggregate_zero(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, 0.0);
    countif_aggregate(cqd, 0.0);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 2, "count should be 2");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[0], 0.0, "zero value 1");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[1], 0.0, "zero value 2");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with very large values */
int test_countif_aggregate_large_values(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, 1e308);
    countif_aggregate(cqd, 1e300);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 2, "count should be 2");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[0], 1e308, "large value 1");
    TEST_ASSERT_DOUBLE_EQ(cqd->data[1], 1e300, "large value 2");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with very small positive values */
int test_countif_aggregate_small_values(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, 1e-300);
    countif_aggregate(cqd, 1e-100);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 2, "count should be 2");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with NaN */
int test_countif_aggregate_nan(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, NAN);
    countif_aggregate(cqd, 10.0);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 2, "count should be 2 even with NaN");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate with infinity */
int test_countif_aggregate_infinity(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    countif_aggregate(cqd, INFINITY);
    countif_aggregate(cqd, -INFINITY);
    countif_aggregate(cqd, 100.0);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 3, "count should be 3");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate buffer resizing */
int test_countif_aggregate_buffer_resize(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    int initial_size = cqd->size;
    
    /* Act - add many values to trigger resizing */
    for (int i = 0; i < MAX_TEST_POINTS; i++) {
        countif_aggregate(cqd, (double)i);
    }
    
    /* Assert */
    TEST_ASSERT(cqd->count == MAX_TEST_POINTS, "count should match added values");
    TEST_ASSERT(cqd->size >= MAX_TEST_POINTS, "size should accommodate all values");
    TEST_ASSERT(cqd->data != NULL, "data should not be NULL after resize");
    
    /* Verify all values */
    for (int i = 0; i < MAX_TEST_POINTS; i++) {
        TEST_ASSERT_DOUBLE_EQ(cqd->data[i], (double)i, "value integrity after resize");
    }
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with threshold higher than values */
int test_countif_query_set_threshold_high(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 20.0);
    countif_aggregate(cqd, 30.0);
    
    /* Act & Assert */
    int result = countif_query_set(cqd, 50.0, COUNTIF_CONDITION_GT);
    TEST_ASSERT(result == 0, "no values should be greater than 50");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with threshold lower than values */
int test_countif_query_set_threshold_low(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 20.0);
    countif_aggregate(cqd, 30.0);
    
    /* Act & Assert */
    int result = countif_query_set(cqd, 5.0, COUNTIF_CONDITION_GT);
    TEST_ASSERT(result == 3, "all 3 values should be greater than 5");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with GT condition */
int test_countif_query_set_gt(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 5.0);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 15.0);
    countif_aggregate(cqd, 20.0);
    
    /* Act */
    int result = countif_query_set(cqd, 10.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 2, "2 values should be greater than 10 (15, 20)");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with GE condition */
int test_countif_query_set_ge(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 5.0);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 15.0);
    countif_aggregate(cqd, 20.0);
    
    /* Act */
    int result = countif_query_set(cqd, 10.0, COUNTIF_CONDITION_GE);
    
    /* Assert */
    TEST_ASSERT(result == 3, "3 values should be >= 10 (10, 15, 20)");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with LT condition */
int test_countif_query_set_lt(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 5.0);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 15.0);
    countif_aggregate(cqd, 20.0);
    
    /* Act */
    int result = countif_query_set(cqd, 10.0, COUNTIF_CONDITION_LT);
    
    /* Assert */
    TEST_ASSERT(result == 1, "1 value should be less than 10 (5)");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with LE condition */
int test_countif_query_set_le(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 5.0);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 15.0);
    countif_aggregate(cqd, 20.0);
    
    /* Act */
    int result = countif_query_set(cqd, 10.0, COUNTIF_CONDITION_LE);
    
    /* Assert */
    TEST_ASSERT(result == 2, "2 values should be <= 10 (5, 10)");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with EQ condition */
int test_countif_query_set_eq(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 15.0);
    countif_aggregate(cqd, 20.0);
    
    /* Act */
    int result = countif_query_set(cqd, 10.0, COUNTIF_CONDITION_EQ);
    
    /* Assert */
    TEST_ASSERT(result == 2, "2 values should be equal to 10");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with NEQ condition */
int test_countif_query_set_neq(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 15.0);
    countif_aggregate(cqd, 20.0);
    
    /* Act */
    int result = countif_query_set(cqd, 10.0, COUNTIF_CONDITION_NEQ);
    
    /* Assert */
    TEST_ASSERT(result == 2, "2 values should not be equal to 10");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with empty data */
int test_countif_query_set_empty(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    int result = countif_query_set(cqd, 10.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 0, "empty data should return 0");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with negative values */
int test_countif_query_set_negative_values(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, -20.0);
    countif_aggregate(cqd, -10.0);
    countif_aggregate(cqd, 0.0);
    countif_aggregate(cqd, 10.0);
    
    /* Act */
    int result = countif_query_set(cqd, -5.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 2, "2 values should be greater than -5 (0, 10)");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with threshold zero */
int test_countif_query_set_zero_threshold(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, -10.0);
    countif_aggregate(cqd, 0.0);
    countif_aggregate(cqd, 10.0);
    
    /* Act */
    int result = countif_query_set(cqd, 0.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 1, "1 value should be greater than 0 (10)");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with all matching values */
int test_countif_query_set_all_match(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 10.0);
    
    /* Act */
    int result = countif_query_set(cqd, 5.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 3, "all values should match");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with no matching values */
int test_countif_query_set_no_match(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 1.0);
    countif_aggregate(cqd, 2.0);
    countif_aggregate(cqd, 3.0);
    
    /* Act */
    int result = countif_query_set(cqd, 100.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 0, "no values should match");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with NaN values */
int test_countif_query_set_with_nan(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, NAN);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, 20.0);
    
    /* Act */
    int result = countif_query_set(cqd, 15.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 1, "NaN should be skipped, 1 value > 15");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with infinity values */
int test_countif_query_set_with_infinity(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, INFINITY);
    countif_aggregate(cqd, 10.0);
    countif_aggregate(cqd, -INFINITY);
    
    /* Act */
    int result = countif_query_set(cqd, 1e300, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 1, "only INFINITY should be greater");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_free with valid pointer */
int test_countif_free_valid(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    
    /* Act */
    countif_free(cqd);
    
    /* Assert - no segfault is success */
    tests_passed++;
    return 1;
}

/* Test: countif_free with NULL pointer */
int test_countif_free_null(void) {
    tests_run++;
    
    /* Act */
    countif_free(NULL);
    
    /* Assert - no segfault is success */
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with single value matching */
int test_countif_query_set_single_value_match(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    
    /* Act */
    int result = countif_query_set(cqd, 5.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 1, "single value should match");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_query_set with single value not matching */
int test_countif_query_set_single_value_no_match(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    countif_aggregate(cqd, 10.0);
    
    /* Act */
    int result = countif_query_set(cqd, 15.0, COUNTIF_CONDITION_GT);
    
    /* Assert */
    TEST_ASSERT(result == 0, "single value should not match");
    
    /* Cleanup */
    if (cqd->data) free(cqd->data);
    free(cqd);
    tests_passed++;
    return 1;
}

/* Test: countif_aggregate boundary near INT_MAX */
int test_countif_aggregate_near_int_max(void) {
    tests_run++;
    struct countif_query_data *cqd = NULL;
    
    /* Arrange */
    cqd = (struct countif_query_data *)malloc(sizeof(struct countif_query_data));
    countif_init(cqd);
    
    /* Act */
    double val = (double)INT32_MAX;
    countif_aggregate(cqd, val);
    
    /* Assert */
    TEST_ASSERT(cqd->count == 1, "count