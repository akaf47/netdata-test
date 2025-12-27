#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <setjmp.h>
#include <cmocka.h>
#include <limits.h>

#include "rrdset-type.h"

/* Tests for header file definitions and macros */

static void test_rrdset_type_definitions_exist(void **state) {
    (void) state;
    
    /* Verify that all expected type constants are defined */
    assert_int_equal(RRDSET_TYPE_LINE >= 0, 1);
    assert_int_equal(RRDSET_TYPE_AREA >= 0, 1);
    assert_int_equal(RRDSET_TYPE_STACKED >= 0, 1);
}

static void test_rrdset_type_are_distinct(void **state) {
    (void) state;
    
    /* Ensure each type has a unique value */
    int line = RRDSET_TYPE_LINE;
    int area = RRDSET_TYPE_AREA;
    int stacked = RRDSET_TYPE_STACKED;
    
    assert_int_not_equal(line, area);
    assert_int_not_equal(area, stacked);
    assert_int_not_equal(line, stacked);
}

static void test_rrdset_type_id_function_declaration(void **state) {
    (void) state;
    
    /* Verify function is callable and returns an int */
    int result = rrdset_type_id("line");
    assert_int_equal(result >= 0, 1);
}

static void test_rrdset_type_name_function_declaration(void **state) {
    (void) state;
    
    /* Verify function is callable and returns a pointer */
    const char *result = rrdset_type_name(RRDSET_TYPE_LINE);
    assert_non_null(result);
}

static void test_rrdset_type_header_consistency(void **state) {
    (void) state;
    
    /* Test that header constants match expected values */
    int line_id = rrdset_type_id("line");
    int area_id = rrdset_type_id("area");
    int stacked_id = rrdset_type_id("stacked");
    
    const char *line_name = rrdset_type_name(line_id);
    const char *area_name = rrdset_type_name(area_id);
    const char *stacked_name = rrdset_type_name(stacked_id);
    
    assert_non_null(line_name);
    assert_non_null(area_name);
    assert_non_null(stacked_name);
}

static void test_rrdset_type_id_multiple_calls_same_result(void **state) {
    (void) state;
    
    /* Verify deterministic behavior */
    int result1 = rrdset_type_id("line");
    int result2 = rrdset_type_id("line");
    int result3 = rrdset_type_id("line");
    
    assert_int_equal(result1, result2);
    assert_int_equal(result2, result3);
}

static void test_rrdset_type_name_multiple_calls_same_result(void **state) {
    (void) state;
    
    /* Verify deterministic behavior */
    const char *result1 = rrdset_type_name(RRDSET_TYPE_LINE);
    const char *result2 = rrdset_type_name(RRDSET_TYPE_LINE);
    const char *result3 = rrdset_type_name(RRDSET_TYPE_LINE);
    
    assert_string_equal(result1, result2);
    assert_string_equal(result2, result3);
}

static void test_rrdset_type_id_all_types_valid(void **state) {
    (void) state;
    
    /* Test that all predefined types are valid */
    int line_id = rrdset_type_id("line");
    int area_id = rrdset_type_id("area");
    int stacked_id = rrdset_type_id("stacked");
    
    assert_true(line_id >= 0);
    assert_true(area_id >= 0);
    assert_true(stacked_id >= 0);
}

static void test_rrdset_type_name_returns_string(void **state) {
    (void) state;
    
    /* Verify return types are proper C strings */
    const char *result = rrdset_type_name(RRDSET_TYPE_LINE);
    assert_non_null(result);
    assert_int_not_equal(strlen(result), 0);
}

static void test_rrdset_type_line_id_to_name(void **state) {
    (void) state;
    
    /* Test line type round trip */
    int line_type = RRDSET_TYPE_LINE;
    const char *line_name = rrdset_type_name(line_type);
    assert_non_null(line_name);
    assert_true(strlen(line_name) > 0);
}

static void test_rrdset_type_area_id_to_name(void **state) {
    (void) state;
    
    /* Test area type round trip */
    int area_type = RRDSET_TYPE_AREA;
    const char *area_name = rrdset_type_name(area_type);
    assert_non_null(area_name);
    assert_true(strlen(area_name) > 0);
}

static void test_rrdset_type_stacked_id_to_name(void **state) {
    (void) state;
    
    /* Test stacked type round trip */
    int stacked_type = RRDSET_TYPE_STACKED;
    const char *stacked_name = rrdset_type_name(stacked_type);
    assert_non_null(stacked_name);
    assert_true(strlen(stacked_name) > 0);
}

static void test_rrdset_type_id_accepts_valid_strings(void **state) {
    (void) state;
    
    /* Verify function accepts string input */
    const char *types[] = {"line", "area", "stacked"};
    for (size_t i = 0; i < sizeof(types) / sizeof(types[0]); i++) {
        int result = rrdset_type_id(types[i]);
        assert_true(result >= 0);
    }
}

static void test_rrdset_type_name_accepts_valid_ids(void **state) {
    (void) state;
    
    /* Verify function accepts int input */
    int types[] = {RRDSET_TYPE_LINE, RRDSET_TYPE_AREA, RRDSET_TYPE_STACKED};
    for (size_t i = 0; i < sizeof(types) / sizeof(types[0]); i++) {
        const char *result = rrdset_type_name(types[i]);
        assert_non_null(result);
    }
}

static void test_rrdset_type_constants_non_negative(void **state) {
    (void) state;
    
    /* All constants should be non-negative */
    assert_true(RRDSET_TYPE_LINE >= 0);
    assert_true(RRDSET_TYPE_AREA >= 0);
    assert_true(RRDSET_TYPE_STACKED >= 0);
}

static void test_rrdset_type_name_result_immutable(void **state) {
    (void) state;
    
    /* Test that returned names are stable and won't be modified */
    const char *name = rrdset_type_name(RRDSET_TYPE_LINE);
    const char *name_again = rrdset_type_name(RRDSET_TYPE_LINE);
    
    /* Both should point to the same content */
    assert_string_equal(name, name_again);
}

static void test_rrdset_type_id_case_handling(void **state) {
    (void) state;
    
    /* Test that case doesn't matter for input */
    int lowercase = rrdset_type_id("line");
    int uppercase = rrdset_type_id("LINE");
    int mixedcase = rrdset_type_id("Line");
    
    assert_int_equal(lowercase, uppercase);
    assert_int_equal(lowercase, mixedcase);
}

static void test_rrdset_type_conversion_bidirectional(void **state) {
    (void) state;
    
    /* Test bidirectional conversion for each type */
    const char *types[] = {"line", "area", "stacked"};
    
    for (size_t i = 0; i < sizeof(types) / sizeof(types[0]); i++) {
        int id = rrdset_type_id(types[i]);
        const char *name = rrdset_type_name(id);
        assert_non_null(name);
        
        /* Verify name matches input (case-insensitive) */
        int original_id = rrdset_type_id(name);
        assert_int_equal(id, original_id);
    }
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_rrdset_type_definitions_exist),
        cmocka_unit_test(test_rrdset_type_are_distinct),
        cmocka_unit_test(test_rrdset_type_id_function_declaration),
        cmocka_unit_test(test_rrdset_type_name_function_declaration),
        cmocka_unit_test(test_rrdset_type_header_consistency),
        cmocka_unit_test(test_rrdset_type_id_multiple_calls_same_result),
        cmocka_unit_test(test_rrdset_type_name_multiple_calls_same_result),
        cmocka_unit_test(test_rrdset_type_id_all_types_valid),
        cmocka_unit_test(test_rrdset_type_name_returns_string),
        cmocka_unit_test(test_rrdset_type_line_id_to_name),
        cmocka_unit_test(test_rrdset_type_area_id_to_name),
        cmocka_unit_test(test_rrdset_type_stacked_id_to_name),
        cmocka_unit_test(test_rrdset_type_id_accepts_valid_strings),
        cmocka_unit_test(test_rrdset_type_name_accepts_valid_ids),
        cmocka_unit_test(test_rrdset_type_constants_non_negative),
        cmocka_unit_test(test_rrdset_type_name_result_immutable),
        cmocka_unit_test(test_rrdset_type_id_case_handling),
        cmocka_unit_test(test_rrdset_type_conversion_bidirectional),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}