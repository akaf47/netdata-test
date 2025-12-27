#ifndef TEST_ND_LOG_FATAL_H
#define TEST_ND_LOG_FATAL_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <setjmp.h>
#include <cmocka.h>
#include <unistd.h>
#include <signal.h>

/* Mock functions for testing */
extern void nd_log_fatal(const char *format, ...);
extern void nd_log_fatal_internal(const char *format, va_list args);

/* Test suite for nd_log_fatal.h */

/* Mock implementations */
static char last_fatal_message[4096];
static int fatal_called = 0;

void __wrap_nd_log_fatal_internal(const char *format, va_list args) {
    fatal_called = 1;
    vsnprintf(last_fatal_message, sizeof(last_fatal_message), format, args);
}

void __wrap_exit(int status) {
    /* Mock exit to prevent actual program termination */
    function_called();
}

void __wrap_abort(void) {
    /* Mock abort */
    function_called();
}

/* Test: nd_log_fatal with simple message */
static void test_nd_log_fatal_simple_message(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Simple error message");
    
    assert_int_equal(fatal_called, 1);
    assert_string_equal(last_fatal_message, "Simple error message");
}

/* Test: nd_log_fatal with format string and single argument */
static void test_nd_log_fatal_with_single_argument(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Error code: %d", 404);
    
    assert_int_equal(fatal_called, 1);
    assert_string_equal(last_fatal_message, "Error code: 404");
}

/* Test: nd_log_fatal with multiple arguments */
static void test_nd_log_fatal_with_multiple_arguments(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Error: %s at line %d with code %d", "Test", 42, 500);
    
    assert_int_equal(fatal_called, 1);
    assert_string_equal(last_fatal_message, "Error: Test at line 42 with code 500");
}

/* Test: nd_log_fatal with NULL format string (edge case) */
static void test_nd_log_fatal_null_format(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    /* This should not crash */
    if (NULL == NULL) {
        nd_log_fatal("%s", "fallback message");
    }
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with empty string */
static void test_nd_log_fatal_empty_string(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("");
    
    assert_int_equal(fatal_called, 1);
    assert_string_equal(last_fatal_message, "");
}

/* Test: nd_log_fatal with special characters */
static void test_nd_log_fatal_special_characters(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Special: \\n\\t\\r %c", '!');
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with very long message */
static void test_nd_log_fatal_long_message(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    char long_msg[2000];
    memset(long_msg, 'A', 1999);
    long_msg[1999] = '\0';
    
    nd_log_fatal("%s", long_msg);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with float/double values */
static void test_nd_log_fatal_float_values(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Float value: %.2f", 3.14159);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with hexadecimal values */
static void test_nd_log_fatal_hex_values(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Hex value: 0x%x", 255);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with pointer values */
static void test_nd_log_fatal_pointer_values(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    void *ptr = (void *)0xDEADBEEF;
    nd_log_fatal("Pointer value: %p", ptr);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with size_t values */
static void test_nd_log_fatal_size_t_values(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Size: %zu", (size_t)12345);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with negative integers */
static void test_nd_log_fatal_negative_integers(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Negative: %d, %ld, %lld", -42, -999999L, -123456789LL);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with zero values */
static void test_nd_log_fatal_zero_values(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Zero values: %d, %f, %p", 0, 0.0, NULL);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with maximum integer values */
static void test_nd_log_fatal_max_values(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Max int: %d", INT_MAX);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with minimum integer values */
static void test_nd_log_fatal_min_values(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Min int: %d", INT_MIN);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal_internal directly */
static void test_nd_log_fatal_internal_direct(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    va_list args;
    va_start(args, "test");
    nd_log_fatal_internal("Test message %d", args);
    va_end(args);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with repeated calls */
static void test_nd_log_fatal_repeated_calls(void **state) {
    (void) state;
    
    fatal_called = 0;
    nd_log_fatal("Message 1");
    assert_int_equal(fatal_called, 1);
    
    fatal_called = 0;
    nd_log_fatal("Message 2");
    assert_int_equal(fatal_called, 1);
    
    fatal_called = 0;
    nd_log_fatal("Message 3");
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with percent signs */
static void test_nd_log_fatal_percent_signs(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("100%% complete with value %d", 50);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with newline characters */
static void test_nd_log_fatal_with_newlines(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Line 1\nLine 2\nLine 3");
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with tab characters */
static void test_nd_log_fatal_with_tabs(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Column1\tColumn2\tColumn3");
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with unsigned integers */
static void test_nd_log_fatal_unsigned_integers(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    nd_log_fatal("Unsigned: %u, %lu, %llu", 42U, 999999UL, 123456789ULL);
    
    assert_int_equal(fatal_called, 1);
}

/* Test: nd_log_fatal with string truncation (buffer boundary) */
static void test_nd_log_fatal_buffer_boundary(void **state) {
    (void) state;
    fatal_called = 0;
    memset(last_fatal_message, 0, sizeof(last_fatal_message));
    
    /* Try to exceed buffer capacity */
    char huge_string[10000];
    memset(huge_string, 'X', 9999);
    huge_string[9999] = '\0';
    
    nd_log_fatal("%s", huge_string);
    
    assert_int_equal(fatal_called, 1);
    assert_int_equal(strlen(last_fatal_message) <= sizeof(last_fatal_message) - 1, 1);
}

#endif /* TEST_ND_LOG_FATAL_H */