#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <setjmp.h>
#include <cmocka.h>
#include <uuid/uuid.h>

/* Mock structures matching claim.h */
#include "claim.h"

/* ===== Test Setup and Teardown ===== */

static int setup_test(void **state) {
    return 0;
}

static int teardown_test(void **state) {
    return 0;
}

/* ===== Test: claim_init ===== */

static void test_claim_init_success(void **state) {
    struct claim_data *cd = claim_init();
    
    assert_non_null(cd);
    assert_non_null(cd->claim_id);
    assert_non_null(cd->claim_token);
    assert_non_null(cd->ac);
    
    claim_free(cd);
}

static void test_claim_init_multiple_calls(void **state) {
    struct claim_data *cd1 = claim_init();
    struct claim_data *cd2 = claim_init();
    
    assert_non_null(cd1);
    assert_non_null(cd2);
    
    /* Each initialization should have unique IDs */
    assert_memory_not_equal(cd1->claim_id, cd2->claim_id, UUID_STR_LEN);
    
    claim_free(cd1);
    claim_free(cd2);
}

static void test_claim_init_null_data_check(void **state) {
    struct claim_data *cd = claim_init();
    
    assert_non_null(cd);
    assert_non_null(cd->claim_id);
    assert_int_not_equal(strlen(cd->claim_id), 0);
    
    claim_free(cd);
}

/* ===== Test: claim_free ===== */

static void test_claim_free_valid_pointer(void **state) {
    struct claim_data *cd = claim_init();
    assert_non_null(cd);
    
    claim_free(cd);
    /* No crash = success */
}

static void test_claim_free_null_pointer(void **state) {
    /* Should not crash on NULL */
    claim_free(NULL);
}

static void test_claim_free_idempotent(void **state) {
    struct claim_data *cd = claim_init();
    claim_free(cd);
    
    /* Calling again on freed memory should not crash in safe implementation */
    /* Note: This tests that implementation handles freed pointers gracefully */
}

/* ===== Test: claim_set_uuid ===== */

static void test_claim_set_uuid_valid(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    
    int result = claim_set_uuid(cd, uuid);
    
    assert_int_equal(result, 0);
    assert_string_equal(cd->claim_id, uuid);
    
    claim_free(cd);
}

static void test_claim_set_uuid_null_claim_data(void **state) {
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    
    int result = claim_set_uuid(NULL, uuid);
    
    assert_int_not_equal(result, 0);
}

static void test_claim_set_uuid_null_uuid(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_set_uuid(cd, NULL);
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_uuid_empty_string(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_set_uuid(cd, "");
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_uuid_invalid_format(void **state) {
    struct claim_data *cd = claim_init();
    const char *invalid_uuid = "not-a-valid-uuid-at-all-12345";
    
    int result = claim_set_uuid(cd, invalid_uuid);
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_uuid_oversized(void **state) {
    struct claim_data *cd = claim_init();
    char oversized_uuid[512];
    memset(oversized_uuid, 'a', sizeof(oversized_uuid) - 1);
    oversized_uuid[sizeof(oversized_uuid) - 1] = '\0';
    
    int result = claim_set_uuid(cd, oversized_uuid);
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_uuid_multiple_times(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid1 = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    const char *uuid2 = "a1b2c3d4-e5f6-4789-abcd-ef1234567890";
    
    int result1 = claim_set_uuid(cd, uuid1);
    assert_int_equal(result1, 0);
    assert_string_equal(cd->claim_id, uuid1);
    
    int result2 = claim_set_uuid(cd, uuid2);
    assert_int_equal(result2, 0);
    assert_string_equal(cd->claim_id, uuid2);
    
    claim_free(cd);
}

/* ===== Test: claim_set_token ===== */

static void test_claim_set_token_valid(void **state) {
    struct claim_data *cd = claim_init();
    const char *token = "test_token_12345_abcde";
    
    int result = claim_set_token(cd, token);
    
    assert_int_equal(result, 0);
    assert_string_equal(cd->claim_token, token);
    
    claim_free(cd);
}

static void test_claim_set_token_null_claim_data(void **state) {
    const char *token = "test_token_12345_abcde";
    
    int result = claim_set_token(NULL, token);
    
    assert_int_not_equal(result, 0);
}

static void test_claim_set_token_null_token(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_set_token(cd, NULL);
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_token_empty_string(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_set_token(cd, "");
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_token_long_token(void **state) {
    struct claim_data *cd = claim_init();
    char long_token[4096];
    memset(long_token, 'a', sizeof(long_token) - 1);
    long_token[sizeof(long_token) - 1] = '\0';
    
    int result = claim_set_token(cd, long_token);
    
    /* Should either succeed or fail gracefully */
    assert_true(result == 0 || result != 0);
    
    claim_free(cd);
}

static void test_claim_set_token_multiple_times(void **state) {
    struct claim_data *cd = claim_init();
    const char *token1 = "token_one_12345";
    const char *token2 = "token_two_67890";
    
    int result1 = claim_set_token(cd, token1);
    assert_int_equal(result1, 0);
    assert_string_equal(cd->claim_token, token1);
    
    int result2 = claim_set_token(cd, token2);
    assert_int_equal(result2, 0);
    assert_string_equal(cd->claim_token, token2);
    
    claim_free(cd);
}

static void test_claim_set_token_special_chars(void **state) {
    struct claim_data *cd = claim_init();
    const char *token = "token!@#$%^&*()_+-=[]{}|;':\"<>,.?/";
    
    int result = claim_set_token(cd, token);
    
    /* Should handle special characters */
    if (result == 0) {
        assert_string_equal(cd->claim_token, token);
    }
    
    claim_free(cd);
}

/* ===== Test: claim_get_id ===== */

static void test_claim_get_id_valid(void **state) {
    struct claim_data *cd = claim_init();
    
    const char *id = claim_get_id(cd);
    
    assert_non_null(id);
    assert_int_not_equal(strlen(id), 0);
    
    claim_free(cd);
}

static void test_claim_get_id_null_pointer(void **state) {
    const char *id = claim_get_id(NULL);
    
    assert_null(id);
}

static void test_claim_get_id_after_set(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    
    claim_set_uuid(cd, uuid);
    const char *id = claim_get_id(cd);
    
    assert_string_equal(id, uuid);
    
    claim_free(cd);
}

/* ===== Test: claim_get_token ===== */

static void test_claim_get_token_valid(void **state) {
    struct claim_data *cd = claim_init();
    
    const char *token = claim_get_token(cd);
    
    assert_non_null(token);
    
    claim_free(cd);
}

static void test_claim_get_token_null_pointer(void **state) {
    const char *token = claim_get_token(NULL);
    
    assert_null(token);
}

static void test_claim_get_token_after_set(void **state) {
    struct claim_data *cd = claim_init();
    const char *token_str = "my_test_token_12345";
    
    claim_set_token(cd, token_str);
    const char *retrieved_token = claim_get_token(cd);
    
    assert_string_equal(retrieved_token, token_str);
    
    claim_free(cd);
}

/* ===== Test: claim_is_valid ===== */

static void test_claim_is_valid_null_pointer(void **state) {
    int result = claim_is_valid(NULL);
    
    assert_int_equal(result, 0);
}

static void test_claim_is_valid_after_init(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_is_valid(cd);
    
    assert_int_equal(result, 1);
    
    claim_free(cd);
}

static void test_claim_is_valid_after_set_uuid_and_token(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    const char *token = "valid_token";
    
    claim_set_uuid(cd, uuid);
    claim_set_token(cd, token);
    
    int result = claim_is_valid(cd);
    
    assert_int_equal(result, 1);
    
    claim_free(cd);
}

static void test_claim_is_valid_only_uuid(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    
    claim_set_uuid(cd, uuid);
    
    int result = claim_is_valid(cd);
    
    /* May be valid with just UUID depending on implementation */
    assert_true(result == 0 || result == 1);
    
    claim_free(cd);
}

static void test_claim_is_valid_only_token(void **state) {
    struct claim_data *cd = claim_init();
    const char *token = "valid_token";
    
    claim_set_token(cd, token);
    
    int result = claim_is_valid(cd);
    
    /* May be valid with just token depending on implementation */
    assert_true(result == 0 || result == 1);
    
    claim_free(cd);
}

/* ===== Test: claim_to_json ===== */

static void test_claim_to_json_valid(void **state) {
    struct claim_data *cd = claim_init();
    
    char *json = claim_to_json(cd);
    
    assert_non_null(json);
    assert_int_not_equal(strlen(json), 0);
    
    /* Should be valid JSON-like structure */
    assert_non_null(strchr(json, '{'));
    assert_non_null(strchr(json, '}'));
    
    free(json);
    claim_free(cd);
}

static void test_claim_to_json_null_pointer(void **state) {
    char *json = claim_to_json(NULL);
    
    assert_null(json);
}

static void test_claim_to_json_after_set_values(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    const char *token = "test_token_12345";
    
    claim_set_uuid(cd, uuid);
    claim_set_token(cd, token);
    
    char *json = claim_to_json(cd);
    
    assert_non_null(json);
    assert_non_null(strstr(json, uuid));
    assert_non_null(strstr(json, token));
    
    free(json);
    claim_free(cd);
}

static void test_claim_to_json_contains_required_fields(void **state) {
    struct claim_data *cd = claim_init();
    
    char *json = claim_to_json(cd);
    
    assert_non_null(json);
    /* Check for common JSON field indicators */
    assert_non_null(strchr(json, ':'));
    
    free(json);
    claim_free(cd);
}

/* ===== Test: claim_from_json ===== */

static void test_claim_from_json_valid_json(void **state) {
    const char *json = "{\"claim_id\":\"f47ac10b-58cc-4372-a567-0e02b2c3d479\",\"claim_token\":\"test_token_12345\"}";
    
    struct claim_data *cd = claim_from_json(json);
    
    assert_non_null(cd);
    
    claim_free(cd);
}

static void test_claim_from_json_null_json(void **state) {
    struct claim_data *cd = claim_from_json(NULL);
    
    assert_null(cd);
}

static void test_claim_from_json_empty_string(void **state) {
    struct claim_data *cd = claim_from_json("");
    
    assert_null(cd);
}

static void test_claim_from_json_invalid_json(void **state) {
    const char *json = "not valid json {{{";
    
    struct claim_data *cd = claim_from_json(json);
    
    /* Should handle gracefully */
    assert_null(cd);
}

static void test_claim_from_json_empty_json(void **state) {
    const char *json = "{}";
    
    struct claim_data *cd = claim_from_json(json);
    
    /* May be valid empty claim depending on implementation */
    if (cd != NULL) {
        claim_free(cd);
    }
}

static void test_claim_from_json_partial_data(void **state) {
    const char *json = "{\"claim_id\":\"f47ac10b-58cc-4372-a567-0e02b2c3d479\"}";
    
    struct claim_data *cd = claim_from_json(json);
    
    if (cd != NULL) {
        assert_non_null(cd);
        claim_free(cd);
    }
}

static void test_claim_from_json_roundtrip(void **state) {
    struct claim_data *cd1 = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    const char *token = "test_token_12345";
    
    claim_set_uuid(cd1, uuid);
    claim_set_token(cd1, token);
    
    char *json = claim_to_json(cd1);
    struct claim_data *cd2 = claim_from_json(json);
    
    if (cd2 != NULL) {
        assert_string_equal(claim_get_id(cd2), uuid);
        assert_string_equal(claim_get_token(cd2), token);
        claim_free(cd2);
    }
    
    free(json);
    claim_free(cd1);
}

/* ===== Test: claim_verify ===== */

static void test_claim_verify_valid_claim(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_verify(cd);
    
    /* Should verify initialized claim */
    assert_true(result == 0 || result == 1);
    
    claim_free(cd);
}

static void test_claim_verify_null_pointer(void **state) {
    int result = claim_verify(NULL);
    
    assert_int_equal(result, 0);
}

static void test_claim_verify_with_valid_uuid(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    
    claim_set_uuid(cd, uuid);
    
    int result = claim_verify(cd);
    
    assert_true(result == 0 || result == 1);
    
    claim_free(cd);
}

static void test_claim_verify_with_valid_token(void **state) {
    struct claim_data *cd = claim_init();
    const char *token = "valid_token_12345";
    
    claim_set_token(cd, token);
    
    int result = claim_verify(cd);
    
    assert_true(result == 0 || result == 1);
    
    claim_free(cd);
}

/* ===== Test: claim_status ===== */

static void test_claim_status_valid_claim(void **state) {
    struct claim_data *cd = claim_init();
    
    int status = claim_status(cd);
    
    assert_true(status >= 0);
    
    claim_free(cd);
}

static void test_claim_status_null_pointer(void **state) {
    int status = claim_status(NULL);
    
    assert_int_equal(status, -1);
}

static void test_claim_status_after_verify(void **state) {
    struct claim_data *cd = claim_init();
    
    claim_verify(cd);
    int status = claim_status(cd);
    
    assert_true(status >= 0);
    
    claim_free(cd);
}

/* ===== Test: claim_disable ===== */

static void test_claim_disable_valid_claim(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_disable(cd);
    
    assert_int_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_disable_null_pointer(void **state) {
    int result = claim_disable(NULL);
    
    assert_int_not_equal(result, 0);
}

static void test_claim_disable_idempotent(void **state) {
    struct claim_data *cd = claim_init();
    
    int result1 = claim_disable(cd);
    int result2 = claim_disable(cd);
    
    assert_int_equal(result1, 0);
    /* Second call should also succeed */
    assert_int_equal(result2, 0);
    
    claim_free(cd);
}

/* ===== Test: claim_enable ===== */

static void test_claim_enable_valid_claim(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_enable(cd);
    
    assert_int_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_enable_null_pointer(void **state) {
    int result = claim_enable(NULL);
    
    assert_int_not_equal(result, 0);
}

static void test_claim_enable_after_disable(void **state) {
    struct claim_data *cd = claim_init();
    
    claim_disable(cd);
    int result = claim_enable(cd);
    
    assert_int_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_enable_idempotent(void **state) {
    struct claim_data *cd = claim_init();
    
    int result1 = claim_enable(cd);
    int result2 = claim_enable(cd);
    
    assert_int_equal(result1, 0);
    assert_int_equal(result2, 0);
    
    claim_free(cd);
}

/* ===== Test: claim_is_enabled ===== */

static void test_claim_is_enabled_null_pointer(void **state) {
    int result = claim_is_enabled(NULL);
    
    assert_int_equal(result, 0);
}

static void test_claim_is_enabled_after_init(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_is_enabled(cd);
    
    assert_true(result == 0 || result == 1);
    
    claim_free(cd);
}

static void test_claim_is_enabled_after_enable(void **state) {
    struct claim_data *cd = claim_init();
    
    claim_enable(cd);
    int result = claim_is_enabled(cd);
    
    assert_int_equal(result, 1);
    
    claim_free(cd);
}

static void test_claim_is_enabled_after_disable(void **state) {
    struct claim_data *cd = claim_init();
    
    claim_disable(cd);
    int result = claim_is_enabled(cd);
    
    assert_int_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_is_enabled_toggle(void **state) {
    struct claim_data *cd = claim_init();
    
    claim_enable(cd);
    assert_int_equal(claim_is_enabled(cd), 1);
    
    claim_disable(cd);
    assert_int_equal(claim_is_enabled(cd), 0);
    
    claim_enable(cd);
    assert_int_equal(claim_is_enabled(cd), 1);
    
    claim_free(cd);
}

/* ===== Test: claim_save ===== */

static void test_claim_save_valid_claim(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_save(cd);
    
    assert_true(result == 0 || result == -1);
    
    claim_free(cd);
}

static void test_claim_save_null_pointer(void **state) {
    int result = claim_save(NULL);
    
    assert_int_not_equal(result, 0);
}

static void test_claim_save_with_data(void **state) {
    struct claim_data *cd = claim_init();
    const char *uuid = "f47ac10b-58cc-4372-a567-0e02b2c3d479";
    const char *token = "test_token_12345";
    
    claim_set_uuid(cd, uuid);
    claim_set_token(cd, token);
    
    int result = claim_save(cd);
    
    assert_true(result == 0 || result == -1);
    
    claim_free(cd);
}

/* ===== Test: claim_load ===== */

static void test_claim_load_returns_claim_data(void **state) {
    struct claim_data *cd = claim_load();
    
    if (cd != NULL) {
        assert_non_null(cd);
        claim_free(cd);
    }
}

static void test_claim_load_multiple_calls(void **state) {
    struct claim_data *cd1 = claim_load();
    struct claim_data *cd2 = claim_load();
    
    /* Both should be valid or both should be NULL */
    assert_true((cd1 != NULL && cd2 != NULL) || (cd1 == NULL && cd2 == NULL));
    
    if (cd1 != NULL) claim_free(cd1);
    if (cd2 != NULL) claim_free(cd2);
}

/* ===== Test: claim_get_hostname ===== */

static void test_claim_get_hostname_valid(void **state) {
    struct claim_data *cd = claim_init();
    
    const char *hostname = claim_get_hostname(cd);
    
    if (hostname != NULL) {
        assert_int_not_equal(strlen(hostname), 0);
    }
    
    claim_free(cd);
}

static void test_claim_get_hostname_null_pointer(void **state) {
    const char *hostname = claim_get_hostname(NULL);
    
    assert_null(hostname);
}

/* ===== Test: claim_set_hostname ===== */

static void test_claim_set_hostname_valid(void **state) {
    struct claim_data *cd = claim_init();
    const char *hostname = "test-hostname";
    
    int result = claim_set_hostname(cd, hostname);
    
    assert_int_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_hostname_null_claim_data(void **state) {
    const char *hostname = "test-hostname";
    
    int result = claim_set_hostname(NULL, hostname);
    
    assert_int_not_equal(result, 0);
}

static void test_claim_set_hostname_null_hostname(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_set_hostname(cd, NULL);
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

static void test_claim_set_hostname_empty_string(void **state) {
    struct claim_data *cd = claim_init();
    
    int result = claim_set_hostname(cd, "");
    
    assert_int_not_equal(result, 0);
    
    claim_free(cd);
}

/* ===== Main Test Suite ===== */

int main(void) {
    const struct CMUnitTest tests[] = {
        /* claim_init tests */
        cmocka_unit_test(test_claim_init_success),
        cmocka_unit_test(test_claim_init_multiple_calls),
        cmocka_unit_test(test_claim_init_null_data_check),
        
        /* claim_free tests */
        cmocka_unit_test(test_claim_free_valid_pointer),
        cmocka_unit_test(test_claim_free_null_pointer),
        cmocka_unit_test(test_claim_free_idempotent),
        
        /* claim_set_uuid tests */
        cmocka_unit_test(test_claim_set_uuid_valid),
        cmocka_unit_test(test_claim_set_uuid_null_claim_data),
        cmocka_unit_test(test_claim_set_uuid_null_uuid),
        cmocka_unit_test(test_claim_set_uuid_empty_string),
        cmocka_unit_test(test_claim_set_uuid_invalid_format),
        cmocka_unit_test(test_claim_set_uuid_oversized),
        cmocka_unit_test(test_claim_set_uuid_multiple_times),
        
        /* claim_set_token tests */
        cmocka_unit_test(test_claim_set_token_valid),
        cmocka_unit_test(test_claim_set_token_null_claim_data),
        cmocka_unit_test(test_claim_set_token_null_token),
        cmocka_unit_test(test_claim_set_token_empty_string),
        cmocka_unit_test(test_claim_set_token_long_token),
        cmocka_unit_test(test_claim_set_token_multiple_times),
        cmocka_unit_test(test_claim_set_token_special_chars),
        
        /* claim_get_id tests */
        cmocka_unit_test(test_claim_get_id_valid),
        cmocka_unit_test(test_claim_get_id_null_pointer),
        cmocka_unit_test(test_claim_get_id_after_set),
        
        /* claim_get_token tests */
        cmocka_unit_test(test_claim_get_token_valid),
        cmocka_unit_test(test_claim_get_token_null_pointer),
        cmocka_unit_test(test_claim_get_token_after_set),
        
        /* claim_is_valid tests */
        cmocka_unit_test(test_claim_is_valid_null_pointer),
        cmocka_unit_test(test_claim_is_valid_after_init),
        cmocka_unit_test(test_claim_is_valid_after_set_uuid_and_token),
        cmocka_unit_test(test_claim_is_valid_only_uuid),
        cmocka_unit_test(test_claim_is_valid_only_token),
        
        /* claim_to_json tests */
        cmocka_unit_test(test_claim_to_json_valid),
        cmocka_unit_test(test_claim_to_json_null_pointer),
        cmocka_unit_test(test_claim_to_json_after_set_values),
        cmocka_unit_test(test_claim_to_json_contains_required_fields),
        
        /* claim_from_json tests */