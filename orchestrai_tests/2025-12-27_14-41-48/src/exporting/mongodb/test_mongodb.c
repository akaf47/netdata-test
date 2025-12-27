#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

/* Mock structures and functions */
#include "../../../database/engine/rrdeng.h"

/* Forward declarations for mocked functions */
int mock_mongoc_init(void);
void mock_mongoc_cleanup(void);
mongoc_uri_t *mock_mongoc_uri_new(const char *uri_string);
mongoc_client_t *mock_mongoc_client_new_from_uri(const mongoc_uri_t *uri);
mongoc_database_t *mock_mongoc_client_get_database(mongoc_client_t *client, const char *name);
mongoc_collection_t *mock_mongoc_database_get_collection(mongoc_database_t *database, const char *name);
int mock_mongoc_collection_insert_one(mongoc_collection_t *collection, const bson_t *document, const mongoc_insert_flags_t flags, bson_t *reply, bson_error_t *error);
bson_t *mock_bson_new(void);
void mock_bson_destroy(bson_t *bson);
bool mock_bson_append_utf8(bson_t *bson, const char *key, int key_length, const char *v_utf8, int v_utf8_len);
bool mock_bson_append_double(bson_t *bson, const char *key, int key_length, double v_double);
bool mock_bson_append_int64(bson_t *bson, const char *key, int key_length, int64_t v_int64);
bool mock_bson_append_int32(bson_t *bson, const char *key, int key_length, int32_t v_int32);
bool mock_bson_append_bool(bson_t *bson, const char *key, int key_length, bool v_bool);
bool mock_bson_append_document(bson_t *bson, const char *key, int key_length, const bson_t *v_document);
char *mock_bson_as_json(const bson_t *bson, size_t *string_len);
void mock_bson_free(char *str);

/* Test helper structures */
typedef struct {
    int mongoc_init_called;
    int mongoc_cleanup_called;
    int bson_new_called;
    int bson_destroy_called;
    int insert_one_called;
    int insert_one_return_value;
    bson_error_t last_error;
} mongodb_test_state;

static mongodb_test_state test_state;

/* Setup and teardown */
static int setup(void **state) {
    memset(&test_state, 0, sizeof(mongodb_test_state));
    test_state.insert_one_return_value = 1;
    *state = &test_state;
    return 0;
}

static int teardown(void **state) {
    memset(&test_state, 0, sizeof(mongodb_test_state));
    return 0;
}

/* Tests for mongoc_init */
static void test_mongoc_init_success(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    will_return(__wrap_mongoc_init, 1);
    int result = mongoc_init();
    assert_int_equal(result, 1);
}

static void test_mongoc_init_multiple_calls(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    will_return(__wrap_mongoc_init, 1);
    will_return(__wrap_mongoc_init, 1);
    assert_int_equal(mongoc_init(), 1);
    assert_int_equal(mongoc_init(), 1);
}

/* Tests for mongoc_cleanup */
static void test_mongoc_cleanup_after_init(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    will_return(__wrap_mongoc_init, 1);
    will_return(__wrap_mongoc_cleanup, 0);
    mongoc_init();
    mongoc_cleanup();
    assert_int_equal(ts->mongoc_cleanup_called, 1);
}

static void test_mongoc_cleanup_without_init(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    will_return(__wrap_mongoc_cleanup, 0);
    mongoc_cleanup();
}

/* Tests for mongoc_uri_new */
static void test_mongoc_uri_new_valid_uri(void **state) {
    const char *uri = "mongodb://localhost:27017";
    mongoc_uri_t *expected_uri = (mongoc_uri_t *)0x12345678;
    
    will_return(__wrap_mongoc_uri_new, expected_uri);
    mongoc_uri_t *result = mongoc_uri_new(uri);
    assert_ptr_equal(result, expected_uri);
}

static void test_mongoc_uri_new_invalid_uri(void **state) {
    const char *uri = "invalid://uri";
    will_return(__wrap_mongoc_uri_new, NULL);
    mongoc_uri_t *result = mongoc_uri_new(uri);
    assert_null(result);
}

static void test_mongoc_uri_new_null_uri(void **state) {
    will_return(__wrap_mongoc_uri_new, NULL);
    mongoc_uri_t *result = mongoc_uri_new(NULL);
    assert_null(result);
}

static void test_mongoc_uri_new_empty_uri(void **state) {
    will_return(__wrap_mongoc_uri_new, NULL);
    mongoc_uri_t *result = mongoc_uri_new("");
    assert_null(result);
}

static void test_mongoc_uri_new_with_authentication(void **state) {
    const char *uri = "mongodb://user:pass@localhost:27017/db";
    mongoc_uri_t *expected_uri = (mongoc_uri_t *)0x87654321;
    
    will_return(__wrap_mongoc_uri_new, expected_uri);
    mongoc_uri_t *result = mongoc_uri_new(uri);
    assert_ptr_equal(result, expected_uri);
}

static void test_mongoc_uri_new_with_replica_set(void **state) {
    const char *uri = "mongodb://localhost:27017,localhost:27018?replicaSet=rs0";
    mongoc_uri_t *expected_uri = (mongoc_uri_t *)0xABCDEF00;
    
    will_return(__wrap_mongoc_uri_new, expected_uri);
    mongoc_uri_t *result = mongoc_uri_new(uri);
    assert_ptr_equal(result, expected_uri);
}

/* Tests for mongoc_client_new_from_uri */
static void test_mongoc_client_new_from_uri_success(void **state) {
    mongoc_uri_t *uri = (mongoc_uri_t *)0x12345678;
    mongoc_client_t *expected_client = (mongoc_client_t *)0x87654321;
    
    will_return(__wrap_mongoc_client_new_from_uri, expected_client);
    mongoc_client_t *result = mongoc_client_new_from_uri(uri);
    assert_ptr_equal(result, expected_client);
}

static void test_mongoc_client_new_from_uri_null_uri(void **state) {
    will_return(__wrap_mongoc_client_new_from_uri, NULL);
    mongoc_client_t *result = mongoc_client_new_from_uri(NULL);
    assert_null(result);
}

static void test_mongoc_client_new_from_uri_connection_failure(void **state) {
    mongoc_uri_t *uri = (mongoc_uri_t *)0x12345678;
    will_return(__wrap_mongoc_client_new_from_uri, NULL);
    mongoc_client_t *result = mongoc_client_new_from_uri(uri);
    assert_null(result);
}

/* Tests for mongoc_client_get_database */
static void test_mongoc_client_get_database_success(void **state) {
    mongoc_client_t *client = (mongoc_client_t *)0x12345678;
    mongoc_database_t *expected_db = (mongoc_database_t *)0x87654321;
    const char *dbname = "netdata";
    
    will_return(__wrap_mongoc_client_get_database, expected_db);
    mongoc_database_t *result = mongoc_client_get_database(client, dbname);
    assert_ptr_equal(result, expected_db);
}

static void test_mongoc_client_get_database_null_client(void **state) {
    will_return(__wrap_mongoc_client_get_database, NULL);
    mongoc_database_t *result = mongoc_client_get_database(NULL, "netdata");
    assert_null(result);
}

static void test_mongoc_client_get_database_null_dbname(void **state) {
    mongoc_client_t *client = (mongoc_client_t *)0x12345678;
    will_return(__wrap_mongoc_client_get_database, NULL);
    mongoc_database_t *result = mongoc_client_get_database(client, NULL);
    assert_null(result);
}

static void test_mongoc_client_get_database_empty_dbname(void **state) {
    mongoc_client_t *client = (mongoc_client_t *)0x12345678;
    will_return(__wrap_mongoc_client_get_database, NULL);
    mongoc_database_t *result = mongoc_client_get_database(client, "");
    assert_null(result);
}

/* Tests for mongoc_database_get_collection */
static void test_mongoc_database_get_collection_success(void **state) {
    mongoc_database_t *database = (mongoc_database_t *)0x12345678;
    mongoc_collection_t *expected_collection = (mongoc_collection_t *)0x87654321;
    const char *collection_name = "metrics";
    
    will_return(__wrap_mongoc_database_get_collection, expected_collection);
    mongoc_collection_t *result = mongoc_database_get_collection(database, collection_name);
    assert_ptr_equal(result, expected_collection);
}

static void test_mongoc_database_get_collection_null_database(void **state) {
    will_return(__wrap_mongoc_database_get_collection, NULL);
    mongoc_collection_t *result = mongoc_database_get_collection(NULL, "metrics");
    assert_null(result);
}

static void test_mongoc_database_get_collection_null_name(void **state) {
    mongoc_database_t *database = (mongoc_database_t *)0x12345678;
    will_return(__wrap_mongoc_database_get_collection, NULL);
    mongoc_collection_t *result = mongoc_database_get_collection(database, NULL);
    assert_null(result);
}

static void test_mongoc_database_get_collection_empty_name(void **state) {
    mongoc_database_t *database = (mongoc_database_t *)0x12345678;
    will_return(__wrap_mongoc_database_get_collection, NULL);
    mongoc_collection_t *result = mongoc_database_get_collection(database, "");
    assert_null(result);
}

/* Tests for bson_new */
static void test_bson_new_success(void **state) {
    bson_t *expected_bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_new, expected_bson);
    bson_t *result = bson_new();
    assert_ptr_equal(result, expected_bson);
}

static void test_bson_new_null_return(void **state) {
    will_return(__wrap_bson_new, NULL);
    bson_t *result = bson_new();
    assert_null(result);
}

static void test_bson_new_multiple_calls(void **state) {
    bson_t *bson1 = (bson_t *)0x11111111;
    bson_t *bson2 = (bson_t *)0x22222222;
    
    will_return(__wrap_bson_new, bson1);
    will_return(__wrap_bson_new, bson2);
    
    bson_t *result1 = bson_new();
    bson_t *result2 = bson_new();
    
    assert_ptr_equal(result1, bson1);
    assert_ptr_equal(result2, bson2);
    assert_ptr_not_equal(result1, result2);
}

/* Tests for bson_destroy */
static void test_bson_destroy_valid_bson(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    bson_t *bson = (bson_t *)0x12345678;
    
    expect_value(__wrap_bson_destroy, bson, bson);
    bson_destroy(bson);
    assert_int_equal(ts->bson_destroy_called, 1);
}

static void test_bson_destroy_null_bson(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    expect_value(__wrap_bson_destroy, bson, NULL);
    bson_destroy(NULL);
}

/* Tests for bson_append_utf8 */
static void test_bson_append_utf8_success(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "hostname";
    const char *value = "server1.example.com";
    
    will_return(__wrap_bson_append_utf8, true);
    bool result = bson_append_utf8(bson, key, -1, value, -1);
    assert_true(result);
}

static void test_bson_append_utf8_failure(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "hostname";
    const char *value = "server1.example.com";
    
    will_return(__wrap_bson_append_utf8, false);
    bool result = bson_append_utf8(bson, key, -1, value, -1);
    assert_false(result);
}

static void test_bson_append_utf8_null_bson(void **state) {
    will_return(__wrap_bson_append_utf8, false);
    bool result = bson_append_utf8(NULL, "key", -1, "value", -1);
    assert_false(result);
}

static void test_bson_append_utf8_null_key(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    will_return(__wrap_bson_append_utf8, false);
    bool result = bson_append_utf8(bson, NULL, -1, "value", -1);
    assert_false(result);
}

static void test_bson_append_utf8_null_value(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    will_return(__wrap_bson_append_utf8, false);
    bool result = bson_append_utf8(bson, "key", -1, NULL, -1);
    assert_false(result);
}

static void test_bson_append_utf8_empty_key(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    will_return(__wrap_bson_append_utf8, false);
    bool result = bson_append_utf8(bson, "", -1, "value", -1);
    assert_false(result);
}

static void test_bson_append_utf8_empty_value(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "metric";
    
    will_return(__wrap_bson_append_utf8, true);
    bool result = bson_append_utf8(bson, key, -1, "", -1);
    assert_true(result);
}

static void test_bson_append_utf8_explicit_lengths(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "hostname";
    const char *value = "server1";
    
    will_return(__wrap_bson_append_utf8, true);
    bool result = bson_append_utf8(bson, key, 8, value, 7);
    assert_true(result);
}

/* Tests for bson_append_double */
static void test_bson_append_double_success(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "cpu_usage";
    double value = 45.67;
    
    will_return(__wrap_bson_append_double, true);
    bool result = bson_append_double(bson, key, -1, value);
    assert_true(result);
}

static void test_bson_append_double_zero(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_double, true);
    bool result = bson_append_double(bson, "value", -1, 0.0);
    assert_true(result);
}

static void test_bson_append_double_negative(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_double, true);
    bool result = bson_append_double(bson, "value", -1, -123.45);
    assert_true(result);
}

static void test_bson_append_double_very_large(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_double, true);
    bool result = bson_append_double(bson, "value", -1, 1.7976931348623157e+308);
    assert_true(result);
}

static void test_bson_append_double_very_small(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_double, true);
    bool result = bson_append_double(bson, "value", -1, 2.2250738585072014e-308);
    assert_true(result);
}

static void test_bson_append_double_failure(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_double, false);
    bool result = bson_append_double(bson, "value", -1, 45.67);
    assert_false(result);
}

static void test_bson_append_double_null_bson(void **state) {
    will_return(__wrap_bson_append_double, false);
    bool result = bson_append_double(NULL, "value", -1, 45.67);
    assert_false(result);
}

/* Tests for bson_append_int64 */
static void test_bson_append_int64_success(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "timestamp";
    int64_t value = 1609459200;
    
    will_return(__wrap_bson_append_int64, true);
    bool result = bson_append_int64(bson, key, -1, value);
    assert_true(result);
}

static void test_bson_append_int64_zero(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int64, true);
    bool result = bson_append_int64(bson, "value", -1, 0);
    assert_true(result);
}

static void test_bson_append_int64_negative(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int64, true);
    bool result = bson_append_int64(bson, "value", -1, -9223372036854775807LL);
    assert_true(result);
}

static void test_bson_append_int64_max(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int64, true);
    bool result = bson_append_int64(bson, "value", -1, 9223372036854775807LL);
    assert_true(result);
}

static void test_bson_append_int64_failure(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int64, false);
    bool result = bson_append_int64(bson, "value", -1, 12345);
    assert_false(result);
}

/* Tests for bson_append_int32 */
static void test_bson_append_int32_success(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "port";
    int32_t value = 27017;
    
    will_return(__wrap_bson_append_int32, true);
    bool result = bson_append_int32(bson, key, -1, value);
    assert_true(result);
}

static void test_bson_append_int32_zero(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int32, true);
    bool result = bson_append_int32(bson, "value", -1, 0);
    assert_true(result);
}

static void test_bson_append_int32_negative(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int32, true);
    bool result = bson_append_int32(bson, "value", -1, -2147483647);
    assert_true(result);
}

static void test_bson_append_int32_max(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int32, true);
    bool result = bson_append_int32(bson, "value", -1, 2147483647);
    assert_true(result);
}

static void test_bson_append_int32_failure(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_int32, false);
    bool result = bson_append_int32(bson, "value", -1, 1000);
    assert_false(result);
}

/* Tests for bson_append_bool */
static void test_bson_append_bool_true(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "enabled";
    
    will_return(__wrap_bson_append_bool, true);
    bool result = bson_append_bool(bson, key, -1, true);
    assert_true(result);
}

static void test_bson_append_bool_false(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *key = "enabled";
    
    will_return(__wrap_bson_append_bool, true);
    bool result = bson_append_bool(bson, key, -1, false);
    assert_true(result);
}

static void test_bson_append_bool_failure(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_bool, false);
    bool result = bson_append_bool(bson, "enabled", -1, true);
    assert_false(result);
}

/* Tests for bson_append_document */
static void test_bson_append_document_success(void **state) {
    bson_t *parent = (bson_t *)0x12345678;
    bson_t *child = (bson_t *)0x87654321;
    const char *key = "tags";
    
    will_return(__wrap_bson_append_document, true);
    bool result = bson_append_document(parent, key, -1, child);
    assert_true(result);
}

static void test_bson_append_document_failure(void **state) {
    bson_t *parent = (bson_t *)0x12345678;
    bson_t *child = (bson_t *)0x87654321;
    
    will_return(__wrap_bson_append_document, false);
    bool result = bson_append_document(parent, "tags", -1, child);
    assert_false(result);
}

static void test_bson_append_document_null_parent(void **state) {
    bson_t *child = (bson_t *)0x87654321;
    
    will_return(__wrap_bson_append_document, false);
    bool result = bson_append_document(NULL, "tags", -1, child);
    assert_false(result);
}

static void test_bson_append_document_null_child(void **state) {
    bson_t *parent = (bson_t *)0x12345678;
    
    will_return(__wrap_bson_append_document, false);
    bool result = bson_append_document(parent, "tags", -1, NULL);
    assert_false(result);
}

/* Tests for bson_as_json */
static void test_bson_as_json_success(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *expected_json = "{\"key\": \"value\"}";
    size_t json_len = strlen(expected_json);
    
    will_return(__wrap_bson_as_json, (char *)expected_json);
    char *result = bson_as_json(bson, &json_len);
    assert_string_equal(result, expected_json);
    assert_int_equal(json_len, strlen(expected_json));
}

static void test_bson_as_json_empty_document(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *expected_json = "{}";
    size_t json_len = 2;
    
    will_return(__wrap_bson_as_json, (char *)expected_json);
    char *result = bson_as_json(bson, &json_len);
    assert_string_equal(result, expected_json);
    assert_int_equal(json_len, 2);
}

static void test_bson_as_json_null_bson(void **state) {
    size_t json_len = 0;
    
    will_return(__wrap_bson_as_json, NULL);
    char *result = bson_as_json(NULL, &json_len);
    assert_null(result);
}

static void test_bson_as_json_null_length(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *expected_json = "{\"key\": \"value\"}";
    
    will_return(__wrap_bson_as_json, (char *)expected_json);
    char *result = bson_as_json(bson, NULL);
    assert_string_equal(result, expected_json);
}

static void test_bson_as_json_complex_document(void **state) {
    bson_t *bson = (bson_t *)0x12345678;
    const char *expected_json = "{\"hostname\": \"server1\", \"cpu\": 45.5, \"timestamp\": 1609459200, \"enabled\": true}";
    size_t json_len = strlen(expected_json);
    
    will_return(__wrap_bson_as_json, (char *)expected_json);
    char *result = bson_as_json(bson, &json_len);
    assert_string_equal(result, expected_json);
}

/* Tests for bson_free */
static void test_bson_free_valid_string(void **state) {
    char *json_string = (char *)0x12345678;
    
    expect_value(__wrap_bson_free, str, json_string);
    bson_free(json_string);
}

static void test_bson_free_null_string(void **state) {
    expect_value(__wrap_bson_free, str, NULL);
    bson_free(NULL);
}

/* Tests for mongoc_collection_insert_one */
static void test_mongoc_collection_insert_one_success(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    mongoc_collection_t *collection = (mongoc_collection_t *)0x12345678;
    bson_t *document = (bson_t *)0x87654321;
    bson_t reply;
    bson_error_t error;
    
    memset(&reply, 0, sizeof(bson_t));
    memset(&error, 0, sizeof(bson_error_t));
    
    will_return(__wrap_mongoc_collection_insert_one, 1);
    int result = mongoc_collection_insert_one(collection, document, 0, &reply, &error);
    assert_int_equal(result, 1);
}

static void test_mongoc_collection_insert_one_failure(void **state) {
    mongodb_test_state *ts = (mongodb_test_state *)*state;
    mongoc_collection_t *collection = (mongoc_collection_t *)0x12345678;
    bson_t *document = (bson_t *)0x87654321;
    bson_t reply;
    bson_error_t error;
    
    memset(&reply, 0, sizeof(bson_t));
    memset(&error, 0, sizeof(bson_error_t));
    
    will_return(__wrap_mongoc_collection_insert_one, 0);
    int result = mongoc_collection_insert_one(collection, document,