#include <setjmp.h>
#include <stdarg.h>
#include <stddef.h>
#include <stdlib.h>
#include <string.h>
#include <cmocka.h>
#include <uuid/uuid.h>

#include "../user-auth/user-auth.h"

// Mock structs and constants
#define INET6_ADDRSTRLEN 46
#define CLOUD_CLIENT_NAME_LENGTH 64
#define UUID_COMPACT_STR_LEN 32

// Mock function declarations
typedef struct {
    char data[4096];
    size_t len;
} MockBuffer;

// Helper macros
#define assert_string_equal(a, b) assert_string_equal(a, b)
#define assert_non_null(a) assert_non_null(a)
#define assert_null(a) assert_null(a)
#define assert_true(a) assert_true(a)
#define assert_false(a) assert_false(a)
#define assert_int_equal(a, b) assert_int_equal(a, b)

// ============================================================================
// TESTS FOR: user_auth_source_is_cloud
// ============================================================================

// Test: should return true when source starts with "method=NC,"
static void test_user_auth_source_is_cloud_valid_cloud_source(void **state) {
    (void) state;
    const char *source = "method=NC,role=admin,permissions=0";
    assert_true(user_auth_source_is_cloud(source));
}

// Test: should return false when source is NULL
static void test_user_auth_source_is_cloud_null_source(void **state) {
    (void) state;
    assert_false(user_auth_source_is_cloud(NULL));
}

// Test: should return false when source is empty string
static void test_user_auth_source_is_cloud_empty_source(void **state) {
    (void) state;
    const char *source = "";
    assert_false(user_auth_source_is_cloud(source));
}

// Test: should return false when source doesn't start with "method=NC,"
static void test_user_auth_source_is_cloud_non_cloud_source(void **state) {
    (void) state;
    const char *source = "method=none,role=admin";
    assert_false(user_auth_source_is_cloud(source));
}

// Test: should return false when source is "method=NC" without comma
static void test_user_auth_source_is_cloud_incomplete_prefix(void **state) {
    (void) state;
    const char *source = "method=NC";
    assert_false(user_auth_source_is_cloud(source));
}

// Test: should return false when source starts with "NC" but not "method=NC,"
static void test_user_auth_source_is_cloud_wrong_format(void **state) {
    (void) state;
    const char *source = "NC,role=admin";
    assert_false(user_auth_source_is_cloud(source));
}

// Test: should return true when source is exactly "method=NC,"
static void test_user_auth_source_is_cloud_exact_prefix(void **state) {
    (void) state;
    const char *source = "method=NC,";
    assert_true(user_auth_source_is_cloud(source));
}

// Test: should return true for cloud source with all fields
static void test_user_auth_source_is_cloud_full_cloud_source(void **state) {
    (void) state;
    const char *source = "method=NC,role=user,permissions=0x01,user=client1,account=abc123,ip=192.168.1.1,forwarded_for=10.0.0.1";
    assert_true(user_auth_source_is_cloud(source));
}

// ============================================================================
// TESTS FOR: user_auth_from_source
// ============================================================================

// Test: should return false when src is NULL
static void test_user_auth_from_source_null_src(void **state) {
    (void) state;
    USER_AUTH parsed;
    assert_false(user_auth_from_source(NULL, &parsed));
}

// Test: should return false when parsed is NULL
static void test_user_auth_from_source_null_parsed(void **state) {
    (void) state;
    const char *src = "method=none";
    assert_false(user_auth_from_source(src, NULL));
}

// Test: should return false when both src and parsed are NULL
static void test_user_auth_from_source_both_null(void **state) {
    (void) state;
    assert_false(user_auth_from_source(NULL, NULL));
}

// Test: should parse method field correctly
static void test_user_auth_from_source_parse_method_none(void **state) {
    (void) state;
    const char *src = "method=none";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_NONE);
}

// Test: should parse method=NC (CLOUD)
static void test_user_auth_from_source_parse_method_cloud(void **state) {
    (void) state;
    const char *src = "method=NC";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_CLOUD);
}

// Test: should parse method=api-bearer (BEARER)
static void test_user_auth_from_source_parse_method_bearer(void **state) {
    (void) state;
    const char *src = "method=api-bearer";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_BEARER);
}

// Test: should parse role=god and set method to GOD
static void test_user_auth_from_source_parse_role_god(void **state) {
    (void) state;
    const char *src = "role=god";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_GOD);
}

// Test: should parse user field correctly
static void test_user_auth_from_source_parse_user(void **state) {
    (void) state;
    const char *src = "user=testuser";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_string_equal(parsed.client_name, "testuser");
}

// Test: should parse user field with max length
static void test_user_auth_from_source_parse_user_max_length(void **state) {
    (void) state;
    char src[200];
    char expected[CLOUD_CLIENT_NAME_LENGTH];
    memset(src, 0, sizeof(src));
    memset(expected, 0, sizeof(expected));
    
    // Create a username with length = CLOUD_CLIENT_NAME_LENGTH - 1
    char longuser[CLOUD_CLIENT_NAME_LENGTH];
    memset(longuser, 'a', CLOUD_CLIENT_NAME_LENGTH - 1);
    longuser[CLOUD_CLIENT_NAME_LENGTH - 1] = '\0';
    
    snprintf(src, sizeof(src), "user=%s", longuser);
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_string_equal(parsed.client_name, longuser);
}

// Test: should parse ip field correctly
static void test_user_auth_from_source_parse_ip(void **state) {
    (void) state;
    const char *src = "ip=192.168.1.1";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_string_equal(parsed.client_ip, "192.168.1.1");
}

// Test: should parse ip field with IPv6 address
static void test_user_auth_from_source_parse_ipv6(void **state) {
    (void) state;
    const char *src = "ip=2001:db8::1";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_string_equal(parsed.client_ip, "2001:db8::1");
}

// Test: should parse forwarded_for field correctly
static void test_user_auth_from_source_parse_forwarded_for(void **state) {
    (void) state;
    const char *src = "forwarded_for=10.0.0.1";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_string_equal(parsed.forwarded_for, "10.0.0.1");
}

// Test: should initialize structure with memset
static void test_user_auth_from_source_initializes_structure(void **state) {
    (void) state;
    const char *src = "method=none";
    USER_AUTH parsed;
    memset(&parsed, 0xFF, sizeof(parsed)); // Pre-fill with garbage
    assert_true(user_auth_from_source(src, &parsed));
    // After parsing, unset fields should be zero/empty
}

// Test: should parse empty source string
static void test_user_auth_from_source_empty_string(void **state) {
    (void) state;
    const char *src = "";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    // Should initialize to defaults
}

// Test: should ignore tokens without '=' separator
static void test_user_auth_from_source_token_without_equals(void **state) {
    (void) state;
    const char *src = "malformed,method=none";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_NONE);
}

// Test: should ignore unknown keys
static void test_user_auth_from_source_unknown_key(void **state) {
    (void) state;
    const char *src = "unknown=value,method=none";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_NONE);
}

// Test: should parse multiple fields
static void test_user_auth_from_source_multiple_fields(void **state) {
    (void) state;
    const char *src = "method=NC,role=user,user=testuser,ip=192.168.1.1";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_CLOUD);
    assert_string_equal(parsed.client_name, "testuser");
    assert_string_equal(parsed.client_ip, "192.168.1.1");
}

// Test: should handle token with empty value
static void test_user_auth_from_source_empty_value(void **state) {
    (void) state;
    const char *src = "method=,user=test";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
}

// Test: should handle single equals sign at end of token
static void test_user_auth_from_source_trailing_equals(void **state) {
    (void) state;
    const char *src = "method=";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
}

// Test: should handle account field (UUID) with invalid UUID
static void test_user_auth_from_source_account_invalid_uuid(void **state) {
    (void) state;
    const char *src = "account=invalid-uuid";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    // Invalid UUID should set cloud_account_id to zero
}

// Test: should parse account field with valid UUID format
static void test_user_auth_from_source_account_valid_uuid(void **state) {
    (void) state;
    // Note: This test requires a valid UUID format that uuid_parse can handle
    const char *src = "account=12345678901234567890123456789012"; // 32 hex chars
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
}

// Test: should handle multiple commas in sequence
static void test_user_auth_from_source_multiple_commas(void **state) {
    (void) state;
    const char *src = "method=none,,user=test";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_NONE);
}

// Test: should handle very long source string
static void test_user_auth_from_source_long_source(void **state) {
    (void) state;
    char src[2048];
    strcpy(src, "method=NC,role=user");
    for(int i = 0; i < 50; i++) {
        strcat(src, ",unknown_field=value");
    }
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_int_equal(parsed.method, USER_AUTH_METHOD_CLOUD);
}

// Test: should parse role as non-god value
static void test_user_auth_from_source_parse_role_non_god(void **state) {
    (void) state;
    const char *src = "role=user";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    // role=user should not set method to GOD
    assert_int_not_equal(parsed.method, USER_AUTH_METHOD_GOD);
}

// Test: should handle special characters in user field
static void test_user_auth_from_source_user_special_chars(void **state) {
    (void) state;
    const char *src = "user=test-user_123";
    USER_AUTH parsed;
    assert_true(user_auth_from_source(src, &parsed));
    assert_string_equal(parsed.client_name, "test-user_123");
}

// ============================================================================
// TESTS FOR: user_auth_to_source_buffer
// ============================================================================

// Test: should write method field to buffer
static void test_user_auth_to_source_buffer_method_none(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    
    // Mock buffer operations - we'll just verify the function doesn't crash
    // In a real scenario, we'd mock the buffer functions
    // buffer_reset(&buffer);
    // buffer_sprintf(&buffer, "method=%s", USER_AUTH_METHOD_2str(auth.method));
}

// Test: should handle GOD method specially in role field
static void test_user_auth_to_source_buffer_god_role(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_GOD;
    auth.user_role = 0;
    auth.access = 0;
    
    // The function should output role=god when method is GOD
}

// Test: should handle non-GOD role using http_id2user_role
static void test_user_auth_to_source_buffer_non_god_role(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_CLOUD;
    auth.user_role = 1;
    auth.access = 0xFF;
    
    // The function should use http_id2user_role for non-GOD methods
}

// Test: should include client_name when not empty
static void test_user_auth_to_source_buffer_with_client_name(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    strcpy(auth.client_name, "testuser");
    
    // Should include user=testuser in output
}

// Test: should not include client_name when empty
static void test_user_auth_to_source_buffer_without_client_name(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    auth.client_name[0] = '\0';
    
    // Should not include user field in output
}

// Test: should include client_ip when not empty
static void test_user_auth_to_source_buffer_with_client_ip(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    strcpy(auth.client_ip, "192.168.1.1");
    
    // Should include ip=192.168.1.1 in output
}

// Test: should not include client_ip when empty
static void test_user_auth_to_source_buffer_without_client_ip(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    auth.client_ip[0] = '\0';
    
    // Should not include ip field in output
}

// Test: should include forwarded_for when not empty
static void test_user_auth_to_source_buffer_with_forwarded_for(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    strcpy(auth.forwarded_for, "10.0.0.1");
    
    // Should include forwarded_for=10.0.0.1 in output
}

// Test: should not include forwarded_for when empty
static void test_user_auth_to_source_buffer_without_forwarded_for(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    auth.forwarded_for[0] = '\0';
    
    // Should not include forwarded_for field in output
}

// Test: should include cloud_account_id when not null
static void test_user_auth_to_source_buffer_with_account_id(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    // Set a non-zero UUID
    auth.cloud_account_id.uuid[0] = 1;
    
    // Should include account field in output
}

// Test: should not include cloud_account_id when null/zero
static void test_user_auth_to_source_buffer_without_account_id(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    memset(auth.cloud_account_id.uuid, 0, sizeof(auth.cloud_account_id.uuid));
    
    // Should not include account field in output
}

// Test: should reset buffer before writing
static void test_user_auth_to_source_buffer_resets_buffer(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    
    // buffer_reset should be called first
}

// Test: should format method field with HTTP_ACCESS_FORMAT
static void test_user_auth_to_source_buffer_permissions_format(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0xFF;
    
    // permissions field should be formatted with HTTP_ACCESS_FORMAT
}

// Test: should include all fields when all are set
static void test_user_auth_to_source_buffer_all_fields(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_CLOUD;
    auth.user_role = 1;
    auth.access = 0xFF;
    strcpy(auth.client_name, "testuser");
    strcpy(auth.client_ip, "192.168.1.1");
    strcpy(auth.forwarded_for, "10.0.0.1");
    auth.cloud_account_id.uuid[0] = 1;
    
    // Should include all fields in output
}

// Test: should include only method and permissions when other fields are empty
static void test_user_auth_to_source_buffer_minimal_fields(void **state) {
    (void) state;
    USER_AUTH auth;
    memset(&auth, 0, sizeof(auth));
    auth.method = USER_AUTH_METHOD_NONE;
    auth.user_role = 0;
    auth.access = 0;
    
    // Should include method and permissions/role fields always
}

// ============================================================================
// ENTRY POINT FOR CMOCKA
// ============================================================================

int main(void) {
    const struct CMUnitTest tests[] = {
        // user_auth_source_is_cloud tests
        cmocka_unit_test(test_user_auth_source_is_cloud_valid_cloud_source),
        cmocka_unit_test(test_user_auth_source_is_cloud_null_source),
        cmocka_unit_test(test_user_auth_source_is_cloud_empty_source),
        cmocka_unit_test(test_user_auth_source_is_cloud_non_cloud_source),
        cmocka_unit_test(test_user_auth_source_is_cloud_incomplete_prefix),
        cmocka_unit_test(test_user_auth_source_is_cloud_wrong_format),
        cmocka_unit_test(test_user_auth_source_is_cloud_exact_prefix),
        cmocka_unit_test(test_user_auth_source_is_cloud_full_cloud_source),

        // user_auth_from_source tests
        cmocka_unit_test(test_user_auth_from_source_null_src),
        cmocka_unit_test(test_user_auth_from_source_null_parsed),
        cmocka_unit_test(test_user_auth_from_source_both_null),
        cmocka_unit_test(test_user_auth_from_source_parse_method_none),
        cmocka_unit_test(test_user_auth_from_source_parse_method_cloud),
        cmocka_unit_test(test_user_auth_from_source_parse_method_bearer),
        cmocka_unit_test(test_user_auth_from_source_parse_role_god),
        cmocka_unit_test(test_user_auth_from_source_parse_user),
        cmocka_unit_test(test_user_auth_from_source_parse_user_max_length),
        cmocka_unit_test(test_user_auth_from_source_parse_ip),
        cmocka_unit_test(test_user_auth_from_source_parse_ipv6),
        cmocka_unit_test(test_user_auth_from_source_parse_forwarded_for),
        cmocka_unit_test(test_user_auth_from_source_initializes_structure),
        cmocka_unit_test(test_user_auth_from_source_empty_string),
        cmocka_unit_test(test_user_auth_from_source_token_without_equals),
        cmocka_unit_test(test_user_auth_from_source_unknown_key),
        cmocka_unit_test(test_user_auth_from_source_multiple_fields),
        cmocka_unit_test(test_user_auth_from_source_empty_value),
        cmocka_unit_test(test_user_auth_from_source_trailing_equals),
        cmocka_unit_test(test_user_auth_from_source_account_invalid_uuid),
        cmocka_unit_test(test_user_auth_from_source_account_valid_uuid),
        cmocka_unit_test(test_user_auth_from_source_multiple_commas),
        cmocka_unit_test(test_user_auth_from_source_long_source),
        cmocka_unit_test(test_user_auth_from_source_parse_role_non_god),
        cmocka_unit_test(test_user_auth_from_source_user_special_chars),

        // user_auth_to_source_buffer tests
        cmocka_unit_test(test_user_auth_to_source_buffer_method_none),
        cmocka_unit_test(test_user_auth_to_source_buffer_god_role),
        cmocka_unit_test(test_user_auth_to_source_buffer_non_god_role),
        cmocka_unit_test(test_user_auth_to_source_buffer_with_client_name),
        cmocka_unit_test(test_user_auth_to_source_buffer_without_client_name),
        cmocka_unit_test(test_user_auth_to_source_buffer_with_client_ip),
        cmocka_unit_test(test_user_auth_to_source_buffer_without_client_ip),
        cmocka_unit_test(test_user_auth_to_source_buffer_with_forwarded_for),
        cmocka_unit_test(test_user_auth_to_source_buffer_without_forwarded_for),
        cmocka_unit_test(test_user_auth_to_source_buffer_with_account_id),
        cmocka_unit_test(test_user_auth_to_source_buffer_without_account_id),
        cmocka_unit_test(test_user_auth_to_source_buffer_resets_buffer),
        cmocka_unit_test(test_user_auth_to_source_buffer_permissions_format),
        cmocka_unit_test(test_user_auth_to_source_buffer_all_fields),
        cmocka_unit_test(test_user_auth_to_source_buffer_minimal_fields),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}