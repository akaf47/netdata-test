#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include <cmocka.h>
#include <assert.h>

// Mock structures and stubs for testing
#include "line_splitter.h"

// ============================================================================
// LINE_SPLITTER STRUCTURE TESTS
// ============================================================================

// Test: Initialize line splitter with valid parameters
static void test_line_splitter_create_valid(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    
    assert_non_null(splitter);
    assert_non_null(splitter->buffer);
    assert_int_equal(splitter->size, 0);
    assert_int_equal(splitter->capacity, LINE_SPLITTER_INITIAL_CAPACITY);
    
    line_splitter_destroy(splitter);
}

// Test: Create multiple line splitters
static void test_line_splitter_create_multiple(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter1 = line_splitter_create();
    LINE_SPLITTER *splitter2 = line_splitter_create();
    LINE_SPLITTER *splitter3 = line_splitter_create();
    
    assert_non_null(splitter1);
    assert_non_null(splitter2);
    assert_non_null(splitter3);
    assert_ptr_not_equal(splitter1, splitter2);
    assert_ptr_not_equal(splitter2, splitter3);
    
    line_splitter_destroy(splitter1);
    line_splitter_destroy(splitter2);
    line_splitter_destroy(splitter3);
}

// Test: Add single line to splitter
static void test_line_splitter_add_line_single(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *line = "test line";
    
    int result = line_splitter_add_line(splitter, line);
    
    assert_int_equal(result, 0);
    assert_int_equal(splitter->size, strlen(line));
    
    line_splitter_destroy(splitter);
}

// Test: Add multiple lines to splitter
static void test_line_splitter_add_line_multiple(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *line1 = "first line";
    const char *line2 = "second line";
    const char *line3 = "third line";
    
    line_splitter_add_line(splitter, line1);
    line_splitter_add_line(splitter, line2);
    line_splitter_add_line(splitter, line3);
    
    assert_int_equal(splitter->size, strlen(line1) + strlen(line2) + strlen(line3));
    
    line_splitter_destroy(splitter);
}

// Test: Add empty line
static void test_line_splitter_add_line_empty(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *line = "";
    
    int result = line_splitter_add_line(splitter, line);
    
    assert_int_equal(result, 0);
    assert_int_equal(splitter->size, 0);
    
    line_splitter_destroy(splitter);
}

// Test: Add NULL line (edge case)
static void test_line_splitter_add_line_null(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    
    int result = line_splitter_add_line(splitter, NULL);
    
    // Should return error code
    assert_int_not_equal(result, 0);
    
    line_splitter_destroy(splitter);
}

// Test: Buffer expansion when size exceeds capacity
static void test_line_splitter_buffer_expansion(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    int initial_capacity = splitter->capacity;
    
    // Create a line larger than initial capacity
    char *large_line = (char *)malloc(initial_capacity + 100);
    memset(large_line, 'a', initial_capacity + 99);
    large_line[initial_capacity + 99] = '\0';
    
    line_splitter_add_line(splitter, large_line);
    
    assert_int_greater_than(splitter->capacity, initial_capacity);
    
    free(large_line);
    line_splitter_destroy(splitter);
}

// Test: Split lines by newline character
static void test_line_splitter_split_newline(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "line1\nline2\nline3";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    assert_string_equal(lines[0].line, "line1");
    assert_string_equal(lines[1].line, "line2");
    assert_string_equal(lines[2].line, "line3");
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Split lines with different line endings
static void test_line_splitter_split_crlf(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "line1\r\nline2\r\nline3";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    assert_int_equal(lines[0].length, 5);
    assert_int_equal(lines[1].length, 5);
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Split single line without newline
static void test_line_splitter_split_single_line(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "single line";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    assert_string_equal(lines[0].line, input);
    assert_int_equal(lines[0].length, strlen(input));
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Split empty buffer
static void test_line_splitter_split_empty(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    // Should return empty array or NULL
    if (lines != NULL) {
        assert_int_equal(lines[0].length, 0);
    }
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Split with trailing newline
static void test_line_splitter_split_trailing_newline(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "line1\nline2\n";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    assert_string_equal(lines[0].line, "line1");
    assert_string_equal(lines[1].line, "line2");
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Split with only newlines
static void test_line_splitter_split_only_newlines(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "\n\n\n";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Clear buffer
static void test_line_splitter_clear(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    line_splitter_add_line(splitter, "test content");
    
    assert_int_greater_than(splitter->size, 0);
    
    line_splitter_clear(splitter);
    
    assert_int_equal(splitter->size, 0);
    
    line_splitter_destroy(splitter);
}

// Test: Clear empty buffer
static void test_line_splitter_clear_empty(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    
    line_splitter_clear(splitter);
    
    assert_int_equal(splitter->size, 0);
    assert_non_null(splitter->buffer);
    
    line_splitter_destroy(splitter);
}

// Test: Reset splitter
static void test_line_splitter_reset(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    line_splitter_add_line(splitter, "test");
    int capacity_before = splitter->capacity;
    
    line_splitter_reset(splitter);
    
    assert_int_equal(splitter->size, 0);
    assert_int_equal(splitter->capacity, LINE_SPLITTER_INITIAL_CAPACITY);
    
    line_splitter_destroy(splitter);
}

// Test: Destroy NULL pointer (edge case)
static void test_line_splitter_destroy_null(void **state) {
    (void)state;
    
    // Should not crash
    line_splitter_destroy(NULL);
}

// Test: Add very long line
static void test_line_splitter_add_very_long_line(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    char *long_line = (char *)malloc(10000);
    memset(long_line, 'x', 9999);
    long_line[9999] = '\0';
    
    int result = line_splitter_add_line(splitter, long_line);
    
    assert_int_equal(result, 0);
    assert_int_equal(splitter->size, 9999);
    
    free(long_line);
    line_splitter_destroy(splitter);
}

// Test: Add lines with special characters
static void test_line_splitter_add_special_characters(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *line = "line with\ttabs\tand\x00special";
    
    int result = line_splitter_add_line(splitter, line);
    
    assert_int_equal(result, 0);
    
    line_splitter_destroy(splitter);
}

// Test: Split with mixed line endings
static void test_line_splitter_split_mixed_endings(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "line1\nline2\r\nline3\rline4";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Line splitter with sequential operations
static void test_line_splitter_sequential_operations(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    
    // Add, split, clear, add again
    line_splitter_add_line(splitter, "first");
    line_splitter_split(splitter);
    line_splitter_clear(splitter);
    line_splitter_add_line(splitter, "second");
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    assert_string_equal(lines[0].line, "second");
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Line splitter get buffer content
static void test_line_splitter_get_buffer(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *content = "buffer content";
    
    line_splitter_add_line(splitter, content);
    const char *buffer = line_splitter_get_buffer(splitter);
    
    assert_non_null(buffer);
    assert_string_equal(buffer, content);
    
    line_splitter_destroy(splitter);
}

// Test: Line splitter get size
static void test_line_splitter_get_size(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    assert_int_equal(line_splitter_get_size(splitter), 0);
    
    line_splitter_add_line(splitter, "test");
    assert_int_equal(line_splitter_get_size(splitter), 4);
    
    line_splitter_destroy(splitter);
}

// Test: Line splitter get capacity
static void test_line_splitter_get_capacity(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    int capacity = line_splitter_get_capacity(splitter);
    
    assert_int_greater_than(capacity, 0);
    
    line_splitter_destroy(splitter);
}

// Test: Reuse splitter after destroy and recreate
static void test_line_splitter_reuse(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter1 = line_splitter_create();
    line_splitter_add_line(splitter1, "first");
    line_splitter_destroy(splitter1);
    
    LINE_SPLITTER *splitter2 = line_splitter_create();
    line_splitter_add_line(splitter2, "second");
    
    assert_int_equal(splitter2->size, 6);
    
    line_splitter_destroy(splitter2);
}

// Test: Split with whitespace only lines
static void test_line_splitter_split_whitespace_lines(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "line1\n   \nline3";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    assert_string_equal(lines[0].line, "line1");
    assert_string_equal(lines[1].line, "   ");
    assert_string_equal(lines[2].line, "line3");
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// Test: Consecutive newlines handling
static void test_line_splitter_consecutive_newlines(void **state) {
    (void)state;
    
    LINE_SPLITTER *splitter = line_splitter_create();
    const char *input = "line1\n\n\nline2";
    
    line_splitter_add_line(splitter, input);
    SPLIT_LINE *lines = line_splitter_split(splitter);
    
    assert_non_null(lines);
    
    line_splitter_free_lines(lines);
    line_splitter_destroy(splitter);
}

// ============================================================================
// TEST SUITE SETUP
// ============================================================================

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_line_splitter_create_valid),
        cmocka_unit_test(test_line_splitter_create_multiple),
        cmocka_unit_test(test_line_splitter_add_line_single),
        cmocka_unit_test(test_line_splitter_add_line_multiple),
        cmocka_unit_test(test_line_splitter_add_line_empty),
        cmocka_unit_test(test_line_splitter_add_line_null),
        cmocka_unit_test(test_line_splitter_buffer_expansion),
        cmocka_unit_test(test_line_splitter_split_newline),
        cmocka_unit_test(test_line_splitter_split_crlf),
        cmocka_unit_test(test_line_splitter_split_single_line),
        cmocka_unit_test(test_line_splitter_split_empty),
        cmocka_unit_test(test_line_splitter_split_trailing_newline),
        cmocka_unit_test(test_line_splitter_split_only_newlines),
        cmocka_unit_test(test_line_splitter_clear),
        cmocka_unit_test(test_line_splitter_clear_empty),
        cmocka_unit_test(test_line_splitter_reset),
        cmocka_unit_test(test_line_splitter_destroy_null),
        cmocka_unit_test(test_line_splitter_add_very_long_line),
        cmocka_unit_test(test_line_splitter_add_special_characters),
        cmocka_unit_test(test_line_splitter_split_mixed_endings),
        cmocka_unit_test(test_line_splitter_sequential_operations),
        cmocka_unit_test(test_line_splitter_get_buffer),
        cmocka_unit_test(test_line_splitter_get_size),
        cmocka_unit_test(test_line_splitter_get_capacity),
        cmocka_unit_test(test_line_splitter_reuse),
        cmocka_unit_test(test_line_splitter_split_whitespace_lines),
        cmocka_unit_test(test_line_splitter_consecutive_newlines),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}