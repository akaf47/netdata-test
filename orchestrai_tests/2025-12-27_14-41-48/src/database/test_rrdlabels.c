#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <assert.h>

#include "rrdlabels.h"

/* Mock structures and test helpers */
typedef struct {
    char *key;
    char *value;
} label_pair;

/* Test fixtures */
static RRDLABELS *test_labels = NULL;

/* Setup and teardown functions */
static int setup(void **state) {
    test_labels = rrdlabels_create();
    assert_non_null(test_labels);
    return 0;
}

static int teardown(void **state) {
    if (test_labels) {
        rrdlabels_free(test_labels);
        test_labels = NULL;
    }
    return 0;
}

/* Tests for rrdlabels_create */
static void test_rrdlabels_create_returns_non_null(void **state) {
    RRDLABELS *labels = rrdlabels_create();
    assert_non_null(labels);
    rrdlabels_free(labels);
}

static void test_rrdlabels_create_multiple_instances(void **state) {
    RRDLABELS *labels1 = rrdlabels_create();
    RRDLABELS *labels2 = rrdlabels_create();
    assert_non_null(labels1);
    assert_non_null(labels2);
    assert_ptr_not_equal(labels1, labels2);
    rrdlabels_free(labels1);
    rrdlabels_free(labels2);
}

/* Tests for rrdlabels_free */
static void test_rrdlabels_free_null_pointer(void **state) {
    /* Should not crash or cause issues */
    rrdlabels_free(NULL);
    assert_true(1);
}

static void test_rrdlabels_free_valid_pointer(void **state) {
    RRDLABELS *labels = rrdlabels_create();
    assert_non_null(labels);
    rrdlabels_free(labels);
    assert_true(1);
}

/* Tests for rrdlabels_add */
static void test_rrdlabels_add_simple_label(void **state) {
    int result = rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 1);
}

static void test_rrdlabels_add_multiple_labels(void **state) {
    int result1 = rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result2 = rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    assert_int_equal(result1, 1);
    assert_int_equal(result2, 1);
}

static void test_rrdlabels_add_update_existing_label(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_add(test_labels, "key1", "value_updated", RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 1);
}

static void test_rrdlabels_add_empty_key(void **state) {
    int result = rrdlabels_add(test_labels, "", "value1", RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 0);
}

static void test_rrdlabels_add_null_key(void **state) {
    int result = rrdlabels_add(test_labels, NULL, "value1", RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 0);
}

static void test_rrdlabels_add_null_value(void **state) {
    int result = rrdlabels_add(test_labels, "key1", NULL, RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 0);
}

static void test_rrdlabels_add_empty_value(void **state) {
    int result = rrdlabels_add(test_labels, "key1", "", RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 1);
}

static void test_rrdlabels_add_null_labels(void **state) {
    int result = rrdlabels_add(NULL, "key1", "value1", RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 0);
}

static void test_rrdlabels_add_different_sources(void **state) {
    int result1 = rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result2 = rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_CONFIG);
    int result3 = rrdlabels_add(test_labels, "key3", "value3", RRDLABEL_SRC_K8S);
    int result4 = rrdlabels_add(test_labels, "key4", "value4", RRDLABEL_SRC_PROMETHEUS);
    assert_int_equal(result1, 1);
    assert_int_equal(result2, 1);
    assert_int_equal(result3, 1);
    assert_int_equal(result4, 1);
}

static void test_rrdlabels_add_special_characters_in_value(void **state) {
    int result = rrdlabels_add(test_labels, "key1", "value with spaces!@#$%", RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 1);
}

static void test_rrdlabels_add_long_strings(void **state) {
    char long_key[256];
    char long_value[1024];
    memset(long_key, 'a', 255);
    long_key[255] = '\0';
    memset(long_value, 'b', 1023);
    long_value[1023] = '\0';
    int result = rrdlabels_add(test_labels, long_key, long_value, RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 1);
}

/* Tests for rrdlabels_get */
static void test_rrdlabels_get_existing_label(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    const char *value = rrdlabels_get(test_labels, "key1");
    assert_non_null(value);
    assert_string_equal(value, "value1");
}

static void test_rrdlabels_get_non_existing_label(void **state) {
    const char *value = rrdlabels_get(test_labels, "nonexistent");
    assert_null(value);
}

static void test_rrdlabels_get_null_labels(void **state) {
    const char *value = rrdlabels_get(NULL, "key1");
    assert_null(value);
}

static void test_rrdlabels_get_null_key(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    const char *value = rrdlabels_get(test_labels, NULL);
    assert_null(value);
}

static void test_rrdlabels_get_empty_key(void **state) {
    const char *value = rrdlabels_get(test_labels, "");
    assert_null(value);
}

static void test_rrdlabels_get_after_update(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key1", "value_updated", RRDLABEL_SRC_AUTO);
    const char *value = rrdlabels_get(test_labels, "key1");
    assert_non_null(value);
    assert_string_equal(value, "value_updated");
}

/* Tests for rrdlabels_del */
static void test_rrdlabels_del_existing_label(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_del(test_labels, "key1");
    assert_int_equal(result, 1);
    const char *value = rrdlabels_get(test_labels, "key1");
    assert_null(value);
}

static void test_rrdlabels_del_non_existing_label(void **state) {
    int result = rrdlabels_del(test_labels, "nonexistent");
    assert_int_equal(result, 0);
}

static void test_rrdlabels_del_null_labels(void **state) {
    int result = rrdlabels_del(NULL, "key1");
    assert_int_equal(result, 0);
}

static void test_rrdlabels_del_null_key(void **state) {
    int result = rrdlabels_del(test_labels, NULL);
    assert_int_equal(result, 0);
}

static void test_rrdlabels_del_empty_key(void **state) {
    int result = rrdlabels_del(test_labels, "");
    assert_int_equal(result, 0);
}

static void test_rrdlabels_del_after_delete(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_del(test_labels, "key1");
    int result = rrdlabels_del(test_labels, "key1");
    assert_int_equal(result, 0);
}

/* Tests for rrdlabels_count */
static void test_rrdlabels_count_empty(void **state) {
    int count = rrdlabels_count(test_labels);
    assert_int_equal(count, 0);
}

static void test_rrdlabels_count_after_add(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int count = rrdlabels_count(test_labels);
    assert_int_equal(count, 1);
}

static void test_rrdlabels_count_multiple_labels(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key3", "value3", RRDLABEL_SRC_AUTO);
    int count = rrdlabels_count(test_labels);
    assert_int_equal(count, 3);
}

static void test_rrdlabels_count_after_delete(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    rrdlabels_del(test_labels, "key1");
    int count = rrdlabels_count(test_labels);
    assert_int_equal(count, 1);
}

static void test_rrdlabels_count_null_labels(void **state) {
    int count = rrdlabels_count(NULL);
    assert_int_equal(count, 0);
}

/* Tests for rrdlabels_copy */
static void test_rrdlabels_copy_non_null(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    RRDLABELS *copy = rrdlabels_copy(test_labels);
    assert_non_null(copy);
    const char *value = rrdlabels_get(copy, "key1");
    assert_non_null(value);
    assert_string_equal(value, "value1");
    rrdlabels_free(copy);
}

static void test_rrdlabels_copy_empty(void **state) {
    RRDLABELS *copy = rrdlabels_copy(test_labels);
    assert_non_null(copy);
    int count = rrdlabels_count(copy);
    assert_int_equal(count, 0);
    rrdlabels_free(copy);
}

static void test_rrdlabels_copy_multiple_labels(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key3", "value3", RRDLABEL_SRC_AUTO);
    RRDLABELS *copy = rrdlabels_copy(test_labels);
    assert_non_null(copy);
    assert_int_equal(rrdlabels_count(copy), 3);
    assert_string_equal(rrdlabels_get(copy, "key1"), "value1");
    assert_string_equal(rrdlabels_get(copy, "key2"), "value2");
    assert_string_equal(rrdlabels_get(copy, "key3"), "value3");
    rrdlabels_free(copy);
}

static void test_rrdlabels_copy_null_labels(void **state) {
    RRDLABELS *copy = rrdlabels_copy(NULL);
    assert_null(copy);
}

static void test_rrdlabels_copy_independence(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    RRDLABELS *copy = rrdlabels_copy(test_labels);
    
    /* Modify original */
    rrdlabels_add(test_labels, "key1", "modified", RRDLABEL_SRC_AUTO);
    
    /* Check copy is unchanged */
    const char *value = rrdlabels_get(copy, "key1");
    assert_string_equal(value, "value1");
    rrdlabels_free(copy);
}

/* Tests for rrdlabels_clear */
static void test_rrdlabels_clear_non_empty(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    rrdlabels_clear(test_labels);
    int count = rrdlabels_count(test_labels);
    assert_int_equal(count, 0);
}

static void test_rrdlabels_clear_empty(void **state) {
    rrdlabels_clear(test_labels);
    int count = rrdlabels_count(test_labels);
    assert_int_equal(count, 0);
}

static void test_rrdlabels_clear_null_labels(void **state) {
    /* Should handle gracefully */
    rrdlabels_clear(NULL);
    assert_true(1);
}

static void test_rrdlabels_clear_allows_readd(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_clear(test_labels);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    assert_int_equal(rrdlabels_count(test_labels), 1);
    assert_string_equal(rrdlabels_get(test_labels, "key2"), "value2");
}

/* Tests for rrdlabels_to_buffer */
static void test_rrdlabels_to_buffer_empty(void **state) {
    BUFFER *buf = buffer_create(4096, NULL);
    size_t result = rrdlabels_to_buffer(test_labels, buf, "prefix", "suffix", "newline");
    assert_true(result >= 0);
    buffer_free(buf);
}

static void test_rrdlabels_to_buffer_single_label(void **state) {
    BUFFER *buf = buffer_create(4096, NULL);
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    size_t result = rrdlabels_to_buffer(test_labels, buf, "", "=", "\n");
    assert_true(result > 0);
    buffer_free(buf);
}

static void test_rrdlabels_to_buffer_multiple_labels(void **state) {
    BUFFER *buf = buffer_create(4096, NULL);
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key3", "value3", RRDLABEL_SRC_AUTO);
    size_t result = rrdlabels_to_buffer(test_labels, buf, "[", "]", " ");
    assert_true(result > 0);
    buffer_free(buf);
}

static void test_rrdlabels_to_buffer_null_labels(void **state) {
    BUFFER *buf = buffer_create(4096, NULL);
    size_t result = rrdlabels_to_buffer(NULL, buf, "", "=", "\n");
    assert_int_equal(result, 0);
    buffer_free(buf);
}

static void test_rrdlabels_to_buffer_null_buffer(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    size_t result = rrdlabels_to_buffer(test_labels, NULL, "", "=", "\n");
    assert_int_equal(result, 0);
}

static void test_rrdlabels_to_buffer_null_separators(void **state) {
    BUFFER *buf = buffer_create(4096, NULL);
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    /* Should handle NULL separators gracefully */
    size_t result = rrdlabels_to_buffer(test_labels, buf, NULL, NULL, NULL);
    assert_true(result >= 0);
    buffer_free(buf);
}

/* Tests for rrdlabels_foreach */
typedef struct {
    int count;
    int call_count;
} foreach_context;

static int foreach_callback(RRDLABELS *labels, const char *key, const char *value, void *data) {
    foreach_context *ctx = (foreach_context *)data;
    ctx->count++;
    if (key && value) {
        ctx->call_count++;
    }
    return 0;
}

static void test_rrdlabels_foreach_empty(void **state) {
    foreach_context ctx = {0, 0};
    int result = rrdlabels_foreach(test_labels, foreach_callback, &ctx);
    assert_int_equal(ctx.count, 0);
}

static void test_rrdlabels_foreach_single_label(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    foreach_context ctx = {0, 0};
    int result = rrdlabels_foreach(test_labels, foreach_callback, &ctx);
    assert_int_equal(ctx.count, 1);
}

static void test_rrdlabels_foreach_multiple_labels(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key3", "value3", RRDLABEL_SRC_AUTO);
    foreach_context ctx = {0, 0};
    int result = rrdlabels_foreach(test_labels, foreach_callback, &ctx);
    assert_int_equal(ctx.count, 3);
}

static void test_rrdlabels_foreach_null_labels(void **state) {
    foreach_context ctx = {0, 0};
    int result = rrdlabels_foreach(NULL, foreach_callback, &ctx);
    assert_int_equal(ctx.count, 0);
}

static void test_rrdlabels_foreach_null_callback(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_foreach(test_labels, NULL, NULL);
    assert_int_equal(result, 0);
}

/* Tests for rrdlabels_callback_replace */
static void test_rrdlabels_callback_replace_all_sources(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_CONFIG);
    int result = rrdlabels_callback_replace(test_labels, foreach_callback, NULL, RRDLABEL_SRC_AUTO);
    assert_true(result >= 0);
}

static void test_rrdlabels_callback_replace_specific_source(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_CONFIG);
    int result = rrdlabels_callback_replace(test_labels, foreach_callback, NULL, RRDLABEL_SRC_CONFIG);
    assert_true(result >= 0);
}

static void test_rrdlabels_callback_replace_null_labels(void **state) {
    int result = rrdlabels_callback_replace(NULL, foreach_callback, NULL, RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 0);
}

static void test_rrdlabels_callback_replace_null_callback(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_callback_replace(test_labels, NULL, NULL, RRDLABEL_SRC_AUTO);
    assert_int_equal(result, 0);
}

/* Tests for rrdlabels_list */
static void test_rrdlabels_list_empty(void **state) {
    char *result = rrdlabels_list(test_labels);
    assert_null(result);
}

static void test_rrdlabels_list_single_label(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    char *result = rrdlabels_list(test_labels);
    assert_non_null(result);
    assert_string_contains(result, "key1");
    assert_string_contains(result, "value1");
    freez(result);
}

static void test_rrdlabels_list_multiple_labels(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(test_labels, "key2", "value2", RRDLABEL_SRC_AUTO);
    char *result = rrdlabels_list(test_labels);
    assert_non_null(result);
    freez(result);
}

static void test_rrdlabels_list_null_labels(void **state) {
    char *result = rrdlabels_list(NULL);
    assert_null(result);
}

/* Tests for rrdlabels_migrate */
static void test_rrdlabels_migrate_non_null(void **state) {
    RRDLABELS *source = rrdlabels_create();
    RRDLABELS *destination = rrdlabels_create();
    
    rrdlabels_add(source, "key1", "value1", RRDLABEL_SRC_AUTO);
    rrdlabels_add(source, "key2", "value2", RRDLABEL_SRC_CONFIG);
    
    int result = rrdlabels_migrate(source, destination);
    assert_true(result >= 0);
    
    rrdlabels_free(source);
    rrdlabels_free(destination);
}

static void test_rrdlabels_migrate_to_non_empty(void **state) {
    RRDLABELS *source = rrdlabels_create();
    RRDLABELS *destination = rrdlabels_create();
    
    rrdlabels_add(destination, "existing", "value", RRDLABEL_SRC_AUTO);
    rrdlabels_add(source, "key1", "value1", RRDLABEL_SRC_AUTO);
    
    int result = rrdlabels_migrate(source, destination);
    assert_true(result >= 0);
    
    rrdlabels_free(source);
    rrdlabels_free(destination);
}

static void test_rrdlabels_migrate_null_source(void **state) {
    RRDLABELS *destination = rrdlabels_create();
    int result = rrdlabels_migrate(NULL, destination);
    assert_int_equal(result, 0);
    rrdlabels_free(destination);
}

static void test_rrdlabels_migrate_null_destination(void **state) {
    RRDLABELS *source = rrdlabels_create();
    rrdlabels_add(source, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_migrate(source, NULL);
    assert_int_equal(result, 0);
    rrdlabels_free(source);
}

static void test_rrdlabels_migrate_both_null(void **state) {
    int result = rrdlabels_migrate(NULL, NULL);
    assert_int_equal(result, 0);
}

/* Tests for rrdlabels_match_simple_pattern */
static void test_rrdlabels_match_simple_pattern_match(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_match_simple_pattern(test_labels, "key1", "value1");
    assert_true(result);
}

static void test_rrdlabels_match_simple_pattern_no_match(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_match_simple_pattern(test_labels, "key1", "value2");
    assert_false(result);
}

static void test_rrdlabels_match_simple_pattern_nonexistent_key(void **state) {
    int result = rrdlabels_match_simple_pattern(test_labels, "nonexistent", "value1");
    assert_false(result);
}

static void test_rrdlabels_match_simple_pattern_null_labels(void **state) {
    int result = rrdlabels_match_simple_pattern(NULL, "key1", "value1");
    assert_false(result);
}

static void test_rrdlabels_match_simple_pattern_null_key(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_match_simple_pattern(test_labels, NULL, "value1");
    assert_false(result);
}

static void test_rrdlabels_match_simple_pattern_null_value(void **state) {
    rrdlabels_add(test_labels, "key1", "value1", RRDLABEL_SRC_AUTO);
    int result = rrdlabels_match_simple_pattern(test_labels, "key1", NULL);
    assert_false(result);
}

/* Tests for rrdlabels_value_is_pattern */
static void test_rrdlabels_value_is_pattern_true_glob(void **state) {
    int result = rrdlabels_value_is_pattern("*value*");
    assert_true(result);
}

static void test_rrdlabels_value_is_pattern_true_question(void **state) {
    int result = rrdlabels_value_is_pattern("value?");
    assert_true(result);
}

static void test_rrdlabels_value_is_pattern_false(void **state) {
    int result = rrdlabels_value_is_pattern("simplevalue");
    assert_false(result);
}

static void test_rrdlabels_value_is_pattern_null(void **state) {
    int result = rrdlabels_value_is_pattern(NULL);
    assert_false(result);
}

static void test_rrdlabels_value_is_pattern_empty(void **state) {
    int result = rrdlabels_value_is_pattern("");
    assert_false(result);
}

static void test_rrdlabels_value_is_pattern_only_special(void **state) {
    int result = rrdlabels_value_is_pattern("*");
    assert_true(