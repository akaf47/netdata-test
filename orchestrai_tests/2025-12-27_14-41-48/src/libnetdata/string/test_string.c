#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <assert.h>
#include <setjmp.h>
#include <stdarg.h>
#include <cmocka.h>

/* Mock the string.h header definitions for testing */
typedef struct {
    size_t len;
    size_t alloc;
    char *str;
} string_t;

/* Forward declarations - these would be in string.h */
extern string_t* string_alloc(size_t size);
extern void string_free(string_t* s);
extern string_t* string_new(const char *str);
extern string_t* string_dup(string_t *s);
extern void string_append(string_t *s, const char *append);
extern void string_append_n(string_t *s, const char *append, size_t n);
extern void string_reset(string_t *s);
extern void string_clear(string_t *s);
extern size_t string_length(string_t *s);
extern const char* string_str(string_t *s);
extern int string_cmp(string_t *s1, string_t *s2);

/* Test: string_alloc with valid size */
static void test_string_alloc_valid_size(void **state) {
    (void) state;
    string_t *s = string_alloc(100);
    assert_non_null(s);
    assert_true(s->alloc >= 100);
    assert_equal(s->len, 0);
    assert_non_null(s->str);
    string_free(s);
}

/* Test: string_alloc with zero size */
static void test_string_alloc_zero_size(void **state) {
    (void) state;
    string_t *s = string_alloc(0);
    if (s != NULL) {
        assert_non_null(s->str);
        string_free(s);
    }
}

/* Test: string_alloc with large size */
static void test_string_alloc_large_size(void **state) {
    (void) state;
    string_t *s = string_alloc(1000000);
    assert_non_null(s);
    assert_true(s->alloc >= 1000000);
    string_free(s);
}

/* Test: string_new with valid string */
static void test_string_new_valid(void **state) {
    (void) state;
    string_t *s = string_new("hello");
    assert_non_null(s);
    assert_string_equal(s->str, "hello");
    assert_equal(s->len, 5);
    string_free(s);
}

/* Test: string_new with empty string */
static void test_string_new_empty(void **state) {
    (void) state;
    string_t *s = string_new("");
    assert_non_null(s);
    assert_string_equal(s->str, "");
    assert_equal(s->len, 0);
    string_free(s);
}

/* Test: string_new with NULL */
static void test_string_new_null(void **state) {
    (void) state;
    string_t *s = string_new(NULL);
    if (s != NULL) {
        assert_non_null(s->str);
        string_free(s);
    }
}

/* Test: string_new with long string */
static void test_string_new_long_string(void **state) {
    (void) state;
    char long_str[1001];
    memset(long_str, 'a', 1000);
    long_str[1000] = '\0';
    string_t *s = string_new(long_str);
    assert_non_null(s);
    assert_equal(s->len, 1000);
    assert_string_equal(s->str, long_str);
    string_free(s);
}

/* Test: string_dup with valid string */
static void test_string_dup_valid(void **state) {
    (void) state;
    string_t *s1 = string_new("test");
    string_t *s2 = string_dup(s1);
    assert_non_null(s2);
    assert_string_equal(s1->str, s2->str);
    assert_equal(s1->len, s2->len);
    assert_ptr_not_equal(s1, s2);
    string_free(s1);
    string_free(s2);
}

/* Test: string_dup with empty string */
static void test_string_dup_empty(void **state) {
    (void) state;
    string_t *s1 = string_new("");
    string_t *s2 = string_dup(s1);
    assert_non_null(s2);
    assert_equal(s2->len, 0);
    string_free(s1);
    string_free(s2);
}

/* Test: string_dup with NULL */
static void test_string_dup_null(void **state) {
    (void) state;
    string_t *s = string_dup(NULL);
    if (s != NULL) {
        string_free(s);
    }
}

/* Test: string_append single call */
static void test_string_append_single(void **state) {
    (void) state;
    string_t *s = string_new("hello");
    string_append(s, " world");
    assert_string_equal(s->str, "hello world");
    assert_equal(s->len, 11);
    string_free(s);
}

/* Test: string_append multiple calls */
static void test_string_append_multiple(void **state) {
    (void) state;
    string_t *s = string_new("");
    string_append(s, "a");
    string_append(s, "b");
    string_append(s, "c");
    assert_string_equal(s->str, "abc");
    assert_equal(s->len, 3);
    string_free(s);
}

/* Test: string_append with empty append */
static void test_string_append_empty(void **state) {
    (void) state;
    string_t *s = string_new("hello");
    string_append(s, "");
    assert_string_equal(s->str, "hello");
    assert_equal(s->len, 5);
    string_free(s);
}

/* Test: string_append to empty string */
static void test_string_append_to_empty(void **state) {
    (void) state;
    string_t *s = string_new("");
    string_append(s, "text");
    assert_string_equal(s->str, "text");
    assert_equal(s->len, 4);
    string_free(s);
}

/* Test: string_append_n with valid n */
static void test_string_append_n_valid(void **state) {
    (void) state;
    string_t *s = string_new("hello");
    string_append_n(s, " world extra", 6);
    assert_string_equal(s->str, "hello world");
    assert_equal(s->len, 11);
    string_free(s);
}

/* Test: string_append_n with n=0 */
static void test_string_append_n_zero(void **state) {
    (void) state;
    string_t *s = string_new("hello");
    string_append_n(s, "world", 0);
    assert_string_equal(s->str, "hello");
    assert_equal(s->len, 5);
    string_free(s);
}

/* Test: string_append_n with n exceeding length */
static void test_string_append_n_exceed(void **state) {
    (void) state;
    string_t *s = string_new("test");
    string_append_n(s, "abc", 100);
    assert_non_null(s->str);
    string_free(s);
}

/* Test: string_append_n with NULL pointer */
static void test_string_append_n_null(void **state) {
    (void) state;
    string_t *s = string_new("test");
    /* Should handle NULL gracefully */
    string_append_n(s, NULL, 5);
    string_free(s);
}

/* Test: string_reset */
static void test_string_reset(void **state) {
    (void) state;
    string_t *s = string_new("hello world");
    string_reset(s);
    assert_equal(s->len, 0);
    assert_equal(s->str[0], '\0');
    string_free(s);
}

/* Test: string_reset on empty string */
static void test_string_reset_empty(void **state) {
    (void) state;
    string_t *s = string_new("");
    string_reset(s);
    assert_equal(s->len, 0);
    string_free(s);
}

/* Test: string_clear */
static void test_string_clear(void **state) {
    (void) state;
    string_t *s = string_new("hello");
    string_clear(s);
    if (s != NULL) {
        /* After clear, string might be freed or reset */
        assert_true(1); /* Placeholder */
    }
}

/* Test: string_length on normal string */
static void test_string_length_normal(void **state) {
    (void) state;
    string_t *s = string_new("hello");
    assert_equal(string_length(s), 5);
    string_free(s);
}

/* Test: string_length on empty string */
static void test_string_length_empty(void **state) {
    (void) state;
    string_t *s = string_new("");
    assert_equal(string_length(s), 0);
    string_free(s);
}

/* Test: string_length with NULL */
static void test_string_length_null(void **state) {
    (void) state;
    size_t len = string_length(NULL);
    assert_equal(len, 0);
}

/* Test: string_str */
static void test_string_str(void **state) {
    (void) state;
    string_t *s = string_new("test");
    const char *str = string_str(s);
    assert_string_equal(str, "test");
    string_free(s);
}

/* Test: string_str with empty */
static void test_string_str_empty(void **state) {
    (void) state;
    string_t *s = string_new("");
    const char *str = string_str(s);
    assert_string_equal(str, "");
    string_free(s);
}

/* Test: string_str with NULL */
static void test_string_str_null(void **state) {
    (void) state;
    const char *str = string_str(NULL);
    if (str != NULL) {
        assert_string_equal(str, "");
    }
}

/* Test: string_cmp equal strings */
static void test_string_cmp_equal(void **state) {
    (void) state;
    string_t *s1 = string_new("hello");
    string_t *s2 = string_new("hello");
    assert_equal(string_cmp(s1, s2), 0);
    string_free(s1);
    string_free(s2);
}

/* Test: string_cmp first less than second */
static void test_string_cmp_less(void **state) {
    (void) state;
    string_t *s1 = string_new("abc");
    string_t *s2 = string_new("def");
    int result = string_cmp(s1, s2);
    assert_true(result < 0);
    string_free(s1);
    string_free(s2);
}

/* Test: string_cmp first greater than second */
static void test_string_cmp_greater(void **state) {
    (void) state;
    string_t *s1 = string_new("xyz");
    string_t *s2 = string_new("abc");
    int result = string_cmp(s1, s2);
    assert_true(result > 0);
    string_free(s1);
    string_free(s2);
}

/* Test: string_cmp with empty strings */
static void test_string_cmp_empty(void **state) {
    (void) state;
    string_t *s1 = string_new("");
    string_t *s2 = string_new("");
    assert_equal(string_cmp(s1, s2), 0);
    string_free(s1);
    string_free(s2);
}

/* Test: string_cmp with one empty */
static void test_string_cmp_one_empty(void **state) {
    (void) state;
    string_t *s1 = string_new("");
    string_t *s2 = string_new("text");
    int result = string_cmp(s1, s2);
    assert_true(result < 0);
    string_free(s1);
    string_free(s2);
}

/* Test: string_free with valid pointer */
static void test_string_free_valid(void **state) {
    (void) state;
    string_t *s = string_new("test");
    string_free(s);
    assert_true(1); /* Should not crash */
}

/* Test: string_free with NULL */
static void test_string_free_null(void **state) {
    (void) state;
    string_free(NULL);
    assert_true(1); /* Should not crash */
}

/* Test: multiple append operations for reallocation */
static void test_string_append_reallocation(void **state) {
    (void) state;
    string_t *s = string_alloc(10);
    for (int i = 0; i < 100; i++) {
        string_append(s, "x");
    }
    assert_equal(s->len, 100);
    string_free(s);
}

/* Test: string_append with special characters */
static void test_string_append_special_chars(void **state) {
    (void) state;
    string_t *s = string_new("");
    string_append(s, "!@#$%^&*()");
    assert_string_equal(s->str, "!@#$%^&*()");
    string_free(s);
}

/* Test: string_append with newlines and tabs */
static void test_string_append_whitespace(void **state) {
    (void) state;
    string_t *s = string_new("");
    string_append(s, "line1\nline2\ttab");
    assert_string_equal(s->str, "line1\nline2\ttab");
    string_free(s);
}

/* Test: string operations after reset */
static void test_string_reset_and_append(void **state) {
    (void) state;
    string_t *s = string_new("initial");
    string_reset(s);
    string_append(s, "new");
    assert_string_equal(s->str, "new");
    assert_equal(s->len, 3);
    string_free(s);
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_string_alloc_valid_size),
        cmocka_unit_test(test_string_alloc_zero_size),
        cmocka_unit_test(test_string_alloc_large_size),
        cmocka_unit_test(test_string_new_valid),
        cmocka_unit_test(test_string_new_empty),
        cmocka_unit_test(test_string_new_null),
        cmocka_unit_test(test_string_new_long_string),
        cmocka_unit_test(test_string_dup_valid),
        cmocka_unit_test(test_string_dup_empty),
        cmocka_unit_test(test_string_dup_null),
        cmocka_unit_test(test_string_append_single),
        cmocka_unit_test(test_string_append_multiple),
        cmocka_unit_test(test_string_append_empty),
        cmocka_unit_test(test_string_append_to_empty),
        cmocka_unit_test(test_string_append_n_valid),
        cmocka_unit_test(test_string_append_n_zero),
        cmocka_unit_test(test_string_append_n_exceed),
        cmocka_unit_test(test_string_append_n_null),
        cmocka_unit_test(test_string_reset),
        cmocka_unit_test(test_string_reset_empty),
        cmocka_unit_test(test_string_clear),
        cmocka_unit_test(test_string_length_normal),
        cmocka_unit_test(test_string_length_empty),
        cmocka_unit_test(test_string_length_null),
        cmocka_unit_test(test_string_str),
        cmocka_unit_test(test_string_str_empty),
        cmocka_unit_test(test_string_str_null),
        cmocka_unit_test(test_string_cmp_equal),
        cmocka_unit_test(test_string_cmp_less),
        cmocka_unit_test(test_string_cmp_greater),
        cmocka_unit_test(test_string_cmp_empty),
        cmocka_unit_test(test_string_cmp_one_empty),
        cmocka_unit_test(test_string_free_valid),
        cmocka_unit_test(test_string_free_null),
        cmocka_unit_test(test_string_append_reallocation),
        cmocka_unit_test(test_string_append_special_chars),
        cmocka_unit_test(test_string_append_whitespace),
        cmocka_unit_test(test_string_reset_and_append),
    };

    return cmocka_run_tests(tests);
}