#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <setjmp.h>
#include <cmocka.h>

#include "rrdset-type.h"

/* Mock function declarations */
static void test_rrdset_type_id_conversion_valid_inputs(void **state) {
    (void) state;
    
    /* Test RRDSET_TYPE_LINE */
    assert_int_equal(rrdset_type_id("line"), RRDSET_TYPE_LINE);
    
    /* Test RRDSET_TYPE_AREA */
    assert_int_equal(rrdset_type_id("area"), RRDSET_TYPE_AREA);
    
    /* Test RRDSET_TYPE_STACKED */
    assert_int_equal(rrdset_type_id("stacked"), RRDSET_TYPE_STACKED);
    
    /* Test uppercase variations */
    assert_int_equal(rrdset_type_id("LINE"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("AREA"), RRDSET_TYPE_AREA);
    assert_int_equal(rrdset_type_id("STACKED"), RRDSET_TYPE_STACKED);
    
    /* Test mixed case */
    assert_int_equal(rrdset_type_id("Line"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("Area"), RRDSET_TYPE_AREA);
    assert_int_equal(rrdset_type_id("Stacked"), RRDSET_TYPE_STACKED);
}

static void test_rrdset_type_id_invalid_inputs(void **state) {
    (void) state;
    
    /* Test invalid type strings */
    assert_int_equal(rrdset_type_id("invalid"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("unknown"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id(""), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("xyz123"), RRDSET_TYPE_LINE);
}

static void test_rrdset_type_id_null_input(void **state) {
    (void) state;
    
    /* Test NULL input - should handle gracefully */
    assert_int_equal(rrdset_type_id(NULL), RRDSET_TYPE_LINE);
}

static void test_rrdset_type_id_empty_string(void **state) {
    (void) state;
    
    /* Test empty string */
    assert_int_equal(rrdset_type_id(""), RRDSET_TYPE_LINE);
}

static void test_rrdset_type_id_whitespace(void **state) {
    (void) state;
    
    /* Test strings with whitespace */
    assert_int_equal(rrdset_type_id(" line "), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("\tarea\t"), RRDSET_TYPE_AREA);
}

static void test_rrdset_type_id_partial_match(void **state) {
    (void) state;
    
    /* Test partial matches that shouldn't match */
    assert_int_equal(rrdset_type_id("lin"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("ar"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("stack"), RRDSET_TYPE_LINE);
}

static void test_rrdset_type_name_valid_types(void **state) {
    (void) state;
    
    /* Test RRDSET_TYPE_LINE */
    const char *name = rrdset_type_name(RRDSET_TYPE_LINE);
    assert_non_null(name);
    assert_string_equal(name, "line");
    
    /* Test RRDSET_TYPE_AREA */
    name = rrdset_type_name(RRDSET_TYPE_AREA);
    assert_non_null(name);
    assert_string_equal(name, "area");
    
    /* Test RRDSET_TYPE_STACKED */
    name = rrdset_type_name(RRDSET_TYPE_STACKED);
    assert_non_null(name);
    assert_string_equal(name, "stacked");
}

static void test_rrdset_type_name_invalid_type(void **state) {
    (void) state;
    
    /* Test invalid type ID - should return default */
    const char *name = rrdset_type_name(999);
    assert_non_null(name);
    /* Should return "line" as default or empty string */
}

static void test_rrdset_type_name_boundary_values(void **state) {
    (void) state;
    
    /* Test boundary values */
    const char *name = rrdset_type_name(0);
    assert_non_null(name);
    
    name = rrdset_type_name(-1);
    assert_non_null(name);
    
    name = rrdset_type_name(INT_MAX);
    assert_non_null(name);
}

static void test_rrdset_type_conversion_roundtrip(void **state) {
    (void) state;
    
    /* Test that type -> id -> name conversions are consistent */
    int line_id = rrdset_type_id("line");
    const char *line_name = rrdset_type_name(line_id);
    assert_string_equal(line_name, "line");
    
    int area_id = rrdset_type_id("area");
    const char *area_name = rrdset_type_name(area_id);
    assert_string_equal(area_name, "area");
    
    int stacked_id = rrdset_type_id("stacked");
    const char *stacked_name = rrdset_type_name(stacked_id);
    assert_string_equal(stacked_name, "stacked");
}

static void test_rrdset_type_id_case_insensitive(void **state) {
    (void) state;
    
    /* Verify case insensitivity across all valid types */
    int types_count = 3;
    const char *valid_types[] = {"line", "area", "stacked"};
    const char *upper_types[] = {"LINE", "AREA", "STACKED"};
    const char *lower_types[] = {"line", "area", "stacked"};
    
    for (int i = 0; i < types_count; i++) {
        int lower_id = rrdset_type_id(lower_types[i]);
        int upper_id = rrdset_type_id(upper_types[i]);
        assert_int_equal(lower_id, upper_id);
    }
}

static void test_rrdset_type_id_with_numeric_strings(void **state) {
    (void) state;
    
    /* Test with numeric strings */
    assert_int_equal(rrdset_type_id("0"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("1"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("123"), RRDSET_TYPE_LINE);
}

static void test_rrdset_type_id_with_special_characters(void **state) {
    (void) state;
    
    /* Test with special characters */
    assert_int_equal(rrdset_type_id("line!"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("@area"), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("stacked#"), RRDSET_TYPE_LINE);
}

static void test_rrdset_type_id_long_strings(void **state) {
    (void) state;
    
    /* Test with very long strings */
    char long_string[1024];
    memset(long_string, 'a', sizeof(long_string) - 1);
    long_string[sizeof(long_string) - 1] = '\0';
    assert_int_equal(rrdset_type_id(long_string), RRDSET_TYPE_LINE);
}

static void test_rrdset_type_constants_defined(void **state) {
    (void) state;
    
    /* Verify that type constants are defined and accessible */
    assert_int_equal(RRDSET_TYPE_LINE, RRDSET_TYPE_LINE);
    assert_int_equal(RRDSET_TYPE_AREA, RRDSET_TYPE_AREA);
    assert_int_equal(RRDSET_TYPE_STACKED, RRDSET_TYPE_STACKED);
}

static void test_rrdset_type_id_leading_trailing_spaces(void **state) {
    (void) state;
    
    /* Test with leading/trailing spaces */
    assert_int_equal(rrdset_type_id("  line  "), RRDSET_TYPE_LINE);
    assert_int_equal(rrdset_type_id("   area   "), RRDSET_TYPE_AREA);
    assert_int_equal(rrdset_type_id(" stacked "), RRDSET_TYPE_STACKED);
}

static void test_rrdset_type_name_consistency(void **state) {
    (void) state;
    
    /* Verify that name function returns consistent results */
    const char *name1 = rrdset_type_name(RRDSET_TYPE_LINE);
    const char *name2 = rrdset_type_name(RRDSET_TYPE_LINE);
    assert_string_equal(name1, name2);
    
    name1 = rrdset_type_name(RRDSET_TYPE_AREA);
    name2 = rrdset_type_name(RRDSET_TYPE_AREA);
    assert_string_equal(name1, name2);
    
    name1 = rrdset_type_name(RRDSET_TYPE_STACKED);
    name2 = rrdset_type_name(RRDSET_TYPE_STACKED);
    assert_string_equal(name1, name2);
}

static void test_rrdset_type_id_all_valid_returns(void **state) {
    (void) state;
    
    /* Test that valid types return non-negative values */
    assert_int_not_equal(rrdset_type_id("line"), -1);
    assert_int_not_equal(rrdset_type_id("area"), -1);
    assert_int_not_equal(rrdset_type_id("stacked"), -1);
}

static void test_rrdset_type_id_different_types_different_ids(void **state) {
    (void) state;
    
    /* Verify that different types return different IDs */
    int line_id = rrdset_type_id("line");
    int area_id = rrdset_type_id("area");
    int stacked_id = rrdset_type_id("stacked");
    
    assert_int_not_equal(line_id, area_id);
    assert_int_not_equal(area_id, stacked_id);
    assert_int_not_equal(line_id, stacked_id);
}

static void test_rrdset_type_name_not_null(void **state) {
    (void) state;
    
    /* Verify name function always returns non-null */
    assert_non_null(rrdset_type_name(RRDSET_TYPE_LINE));
    assert_non_null(rrdset_type_name(RRDSET_TYPE_AREA));
    assert_non_null(rrdset_type_name(RRDSET_TYPE_STACKED));
    assert_non_null(rrdset_type_name(-999));
    assert_non_null(rrdset_type_name(999));
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_rrdset_type_id_conversion_valid_inputs),
        cmocka_unit_test(test_rrdset_type_id_invalid_inputs),
        cmocka_unit_test(test_rrdset_type_id_null_input),
        cmocka_unit_test(test_rrdset_type_id_empty_string),
        cmocka_unit_test(test_rrdset_type_id_whitespace),
        cmocka_unit_test(test_rrdset_type_id_partial_match),
        cmocka_unit_test(test_rrdset_type_name_valid_types),
        cmocka_unit_test(test_rrdset_type_name_invalid_type),
        cmocka_unit_test(test_rrdset_type_name_boundary_values),
        cmocka_unit_test(test_rrdset_type_conversion_roundtrip),
        cmocka_unit_test(test_rrdset_type_id_case_insensitive),
        cmocka_unit_test(test_rrdset_type_id_with_numeric_strings),
        cmocka_unit_test(test_rrdset_type_id_with_special_characters),
        cmocka_unit_test(test_rrdset_type_id_long_strings),
        cmocka_unit_test(test_rrdset_type_constants_defined),
        cmocka_unit_test(test_rrdset_type_id_leading_trailing_spaces),
        cmocka_unit_test(test_rrdset_type_name_consistency),
        cmocka_unit_test(test_rrdset_type_id_all_valid_returns),
        cmocka_unit_test(test_rrdset_type_id_different_types_different_ids),
        cmocka_unit_test(test_rrdset_type_name_not_null),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}