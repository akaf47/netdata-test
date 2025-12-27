#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>
#include <setjmp.h>
#include <cmocka.h>

#include "simple_pattern.h"

/* Mock functions and test fixtures */
static void test_simple_pattern_compile_null_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile(NULL);
    assert_null(result);
}

static void test_simple_pattern_compile_empty_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("");
    assert_null(result);
}

static void test_simple_pattern_compile_single_wildcard(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("*");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_literal_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_wildcard_prefix(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("*test");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_wildcard_suffix(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test*");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_wildcard_both(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("*test*");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_multiple_wildcards(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("*test*match*");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_question_mark(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test?");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_pipe_separator(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test|match");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_multiple_pipe_separators(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test|match|pattern");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_complex_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("*test*|match*|*pattern");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_matches_null_pattern(void **state) {
    (void)state;
    int result = simple_pattern_matches(NULL, "test");
    assert_int_equal(result, 0);
}

static void test_simple_pattern_matches_null_string(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test");
    int result = simple_pattern_matches(pattern, NULL);
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_empty_string(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test");
    int result = simple_pattern_matches(pattern, "");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_exact_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test");
    int result = simple_pattern_matches(pattern, "other");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_wildcard_all(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*");
    int result = simple_pattern_matches(pattern, "anything");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_wildcard_all_empty(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*");
    int result = simple_pattern_matches(pattern, "");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_prefix_wildcard(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test");
    int result = simple_pattern_matches(pattern, "mytest");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_prefix_wildcard_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test");
    int result = simple_pattern_matches(pattern, "testing");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_suffix_wildcard(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test*");
    int result = simple_pattern_matches(pattern, "testing");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_suffix_wildcard_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test*");
    int result = simple_pattern_matches(pattern, "mytest");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_both_wildcards(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test*");
    int result = simple_pattern_matches(pattern, "mytesting");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_both_wildcards_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test*");
    int result = simple_pattern_matches(pattern, "mystring");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_question_mark_single(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test?");
    int result = simple_pattern_matches(pattern, "tests");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_question_mark_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test?");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_question_mark_too_long(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test?");
    int result = simple_pattern_matches(pattern, "tests123");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_multiple_question_marks(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("t?s?");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_pipe_first_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test|match");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_pipe_second_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test|match");
    int result = simple_pattern_matches(pattern, "match");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_pipe_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test|match");
    int result = simple_pattern_matches(pattern, "other");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_pipe_with_wildcards_first(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test|match");
    int result = simple_pattern_matches(pattern, "mytest");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_pipe_with_wildcards_second(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test|*match");
    int result = simple_pattern_matches(pattern, "mymatch");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_case_sensitive(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("Test");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_long_string(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test");
    char long_string[1024];
    strcpy(long_string, "this is a very long string containing the word test in it");
    int result = simple_pattern_matches(pattern, long_string);
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_free_null(void **state) {
    (void)state;
    simple_pattern_free(NULL);
}

static void test_simple_pattern_free_valid(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test");
    simple_pattern_free(pattern);
}

static void test_simple_pattern_free_complex(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test*|match*|*pattern");
    simple_pattern_free(pattern);
}

static void test_simple_pattern_compile_special_chars(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test@pattern");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_compile_numbers(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test123");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_matches_numbers(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("123");
    int result = simple_pattern_matches(pattern, "123");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_wildcard_with_numbers(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*123*");
    int result = simple_pattern_matches(pattern, "abc123def");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_compile_dot_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *result = simple_pattern_compile("test.pattern");
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_matches_exact_after_compile(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("exact");
    int result = simple_pattern_matches(pattern, "exact");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_with_spaces(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test pattern");
    int result = simple_pattern_matches(pattern, "test pattern");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_wildcard_with_spaces(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test*");
    int result = simple_pattern_matches(pattern, "my test string");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_multiple_pipes_first(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("first|second|third");
    int result = simple_pattern_matches(pattern, "first");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_multiple_pipes_middle(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("first|second|third");
    int result = simple_pattern_matches(pattern, "second");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_multiple_pipes_last(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("first|second|third");
    int result = simple_pattern_matches(pattern, "third");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_boundary_single_char_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("a");
    int result = simple_pattern_matches(pattern, "a");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_boundary_single_char_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("a");
    int result = simple_pattern_matches(pattern, "b");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_pipe_with_empty_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test|");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_complex_mixed_wildcards_and_question(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test?*");
    int result = simple_pattern_matches(pattern, "mytestxstring");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_question_at_start(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("?test");
    int result = simple_pattern_matches(pattern, "atest");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_question_at_start_no_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("?test");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 0);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_compile_very_long_pattern(void **state) {
    (void)state;
    char pattern[512];
    strcpy(pattern, "verylongpattern");
    strcat(pattern, "|anotherlongpattern");
    strcat(pattern, "|thirdlongpattern");
    SIMPLE_PATTERN *result = simple_pattern_compile(pattern);
    assert_non_null(result);
    simple_pattern_free(result);
}

static void test_simple_pattern_matches_very_long_pattern(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("verylongpattern|anotherlongpattern|thirdlongpattern");
    int result = simple_pattern_matches(pattern, "anotherlongpattern");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_suffix_wildcard_exact_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("test*");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_prefix_wildcard_exact_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static void test_simple_pattern_matches_both_wildcards_exact_match(void **state) {
    (void)state;
    SIMPLE_PATTERN *pattern = simple_pattern_compile("*test*");
    int result = simple_pattern_matches(pattern, "test");
    assert_int_equal(result, 1);
    simple_pattern_free(pattern);
}

static const struct CMUnitTest tests[] = {
    cmocka_unit_test(test_simple_pattern_compile_null_pattern),
    cmocka_unit_test(test_simple_pattern_compile_empty_pattern),
    cmocka_unit_test(test_simple_pattern_compile_single_wildcard),
    cmocka_unit_test(test_simple_pattern_compile_literal_pattern),
    cmocka_unit_test(test_simple_pattern_compile_wildcard_prefix),
    cmocka_unit_test(test_simple_pattern_compile_wildcard_suffix),
    cmocka_unit_test(test_simple_pattern_compile_wildcard_both),
    cmocka_unit_test(test_simple_pattern_compile_multiple_wildcards),
    cmocka_unit_test(test_simple_pattern_compile_question_mark),
    cmocka_unit_test(test_simple_pattern_compile_pipe_separator),
    cmocka_unit_test(test_simple_pattern_compile_multiple_pipe_separators),
    cmocka_unit_test(test_simple_pattern_compile_complex_pattern),
    cmocka_unit_test(test_simple_pattern_matches_null_pattern),
    cmocka_unit_test(test_simple_pattern_matches_null_string),
    cmocka_unit_test(test_simple_pattern_matches_empty_string),
    cmocka_unit_test(test_simple_pattern_matches_exact_match),
    cmocka_unit_test(test_simple_pattern_matches_no_match),
    cmocka_unit_test(test_simple_pattern_matches_wildcard_all),
    cmocka_unit_test(test_simple_pattern_matches_wildcard_all_empty),
    cmocka_unit_test(test_simple_pattern_matches_prefix_wildcard),
    cmocka_unit_test(test_simple_pattern_matches_prefix_wildcard_no_match),
    cmocka_unit_test(test_simple_pattern_matches_suffix_wildcard),
    cmocka_unit_test(test_simple_pattern_matches_suffix_wildcard_no_match),
    cmocka_unit_test(test_simple_pattern_matches_both_wildcards),
    cmocka_unit_test(test_simple_pattern_matches_both_wildcards_no_match),
    cmocka_unit_test(test_simple_pattern_matches_question_mark_single),
    cmocka_unit_test(test_simple_pattern_matches_question_mark_no_match),
    cmocka_unit_test(test_simple_pattern_matches_question_mark_too_long),
    cmocka_unit_test(test_simple_pattern_matches_multiple_question_marks),
    cmocka_unit_test(test_simple_pattern_matches_pipe_first_match),
    cmocka_unit_test(test_simple_pattern_matches_pipe_second_match),
    cmocka_unit_test(test_simple_pattern_matches_pipe_no_match),
    cmocka_unit_test(test_simple_pattern_matches_pipe_with_wildcards_first),
    cmocka_unit_test(test_simple_pattern_matches_pipe_with_wildcards_second),
    cmocka_unit_test(test_simple_pattern_matches_case_sensitive),
    cmocka_unit_test(test_simple_pattern_matches_long_string),
    cmocka_unit_test(test_simple_pattern_free_null),
    cmocka_unit_test(test_simple_pattern_free_valid),
    cmocka_unit_test(test_simple_pattern_free_complex),
    cmocka_unit_test(test_simple_pattern_compile_special_chars),
    cmocka_unit_test(test_simple_pattern_compile_numbers),
    cmocka_unit_test(test_simple_pattern_matches_numbers),
    cmocka_unit_test(test_simple_pattern_matches_wildcard_with_numbers),
    cmocka_unit_test(test_simple_pattern_compile_dot_pattern),
    cmocka_unit_test(test_simple_pattern_matches_exact_after_compile),
    cmocka_unit_test(test_simple_pattern_matches_with_spaces),
    cmocka_unit_test(test_simple_pattern_matches_wildcard_with_spaces),
    cmocka_unit_test(test_simple_pattern_matches_multiple_pipes_first),
    cmocka_unit_test(test_simple_pattern_matches_multiple_pipes_middle),
    cmocka_unit_test(test_simple_pattern_matches_multiple_pipes_last),
    cmocka_unit_test(test_simple_pattern_matches_boundary_single_char_pattern),
    cmocka_unit_test(test_simple_pattern_matches_boundary_single_char_no_match),
    cmocka_unit_test(test_simple_pattern_matches_pipe_with_empty_pattern),
    cmocka_unit_test(test_simple_pattern_matches_complex_mixed_wildcards_and_question),
    cmocka_unit_test(test_simple_pattern_matches_question_at_start),
    cmocka_unit_test(test_simple_pattern_matches_question_at_start_no_match),
    cmocka_unit_test(test_simple_pattern_compile_very_long_pattern),
    cmocka_unit_test(test_simple_pattern_matches_very_long_pattern),
    cmocka_unit_test(test_simple_pattern_matches_suffix_wildcard_exact_match),
    cmocka_unit_test(test_simple_pattern_matches_prefix_wildcard_exact_match),
    cmocka_unit_test(test_simple_pattern_matches_both_wildcards_exact_match),
};

int main(void) {
    return cmocka_run_group_tests(tests, NULL, NULL);
}