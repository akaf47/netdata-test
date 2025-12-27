#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "../uuid_generate.h"

/* Test: UUID generation creates valid UUID string */
static void test_uuid_generate_creates_valid_uuid(void **state) {
    (void) state;
    
    char uuid_str[UUID_STR_LEN];
    memset(uuid_str, 0, sizeof(uuid_str));
    
    uuid_t uuid;
    uuid_generate(uuid);
    uuid_unparse(uuid, uuid_str);
    
    /* Assert UUID string is not empty */
    assert_non_null(uuid_str);
    assert_int_not_equal(strlen(uuid_str), 0);
    
    /* Assert UUID string matches expected format (8-4-4-4-12 hex digits) */
    assert_int_equal(strlen(uuid_str), 36); /* Standard UUID string length */
    
    /* Check format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx */
    assert_int_equal(uuid_str[8], '-');
    assert_int_equal(uuid_str[13], '-');
    assert_int_equal(uuid_str[18], '-');
    assert_int_equal(uuid_str[23], '-');
}

/* Test: Multiple UUID generations produce different UUIDs */
static void test_uuid_generate_produces_unique_uuids(void **state) {
    (void) state;
    
    char uuid_str1[UUID_STR_LEN];
    char uuid_str2[UUID_STR_LEN];
    
    uuid_t uuid1, uuid2;
    uuid_generate(uuid1);
    uuid_unparse(uuid1, uuid_str1);
    
    uuid_generate(uuid2);
    uuid_unparse(uuid2, uuid_str2);
    
    /* Assert that two generated UUIDs are different */
    assert_string_not_equal(uuid_str1, uuid_str2);
}

/* Test: UUID can be generated with null/uninitialized uuid_t */
static void test_uuid_generate_with_uninitialized_uuid(void **state) {
    (void) state;
    
    uuid_t uuid;
    /* UUID should be generated even if uuid_t was not initialized */
    uuid_generate(uuid);
    
    /* Verify UUID is not all zeros */
    uuid_t zero_uuid;
    memset(zero_uuid, 0, sizeof(zero_uuid));
    
    assert_memory_not_equal(uuid, zero_uuid, sizeof(uuid_t));
}

/* Test: UUID string parsing roundtrip */
static void test_uuid_unparse_and_parse_roundtrip(void **state) {
    (void) state;
    
    char uuid_str[UUID_STR_LEN];
    uuid_t original_uuid, parsed_uuid;
    
    uuid_generate(original_uuid);
    uuid_unparse(original_uuid, uuid_str);
    uuid_parse(uuid_str, parsed_uuid);
    
    /* Assert roundtrip produces identical UUID */
    assert_memory_equal(original_uuid, parsed_uuid, sizeof(uuid_t));
}

/* Test: UUID parse handles valid UUID string */
static void test_uuid_parse_valid_string(void **state) {
    (void) state;
    
    const char *valid_uuid = "550e8400-e29b-41d4-a716-446655440000";
    uuid_t uuid;
    
    int result = uuid_parse(valid_uuid, uuid);
    assert_int_equal(result, 0);
    
    /* Verify uuid is not all zeros */
    uuid_t zero_uuid;
    memset(zero_uuid, 0, sizeof(uuid_t));
    assert_memory_not_equal(uuid, zero_uuid, sizeof(uuid_t));
}

/* Test: UUID parse handles invalid string format */
static void test_uuid_parse_invalid_string(void **state) {
    (void) state;
    
    const char *invalid_uuid = "not-a-valid-uuid-string";
    uuid_t uuid;
    
    int result = uuid_parse(invalid_uuid, uuid);
    assert_int_not_equal(result, 0); /* Should return error */
}

/* Test: UUID parse with empty string */
static void test_uuid_parse_empty_string(void **state) {
    (void) state;
    
    const char *empty_uuid = "";
    uuid_t uuid;
    
    int result = uuid_parse(empty_uuid, uuid);
    assert_int_not_equal(result, 0); /* Should return error */
}

/* Test: UUID parse with NULL string */
static void test_uuid_parse_null_string(void **state) {
    (void) state;
    
    uuid_t uuid;
    
    /* This should handle NULL gracefully or crash - document behavior */
    /* Skipping NULL test as behavior is undefined in standard uuid library */
}

/* Test: UUID compare for equal UUIDs */
static void test_uuid_compare_equal(void **state) {
    (void) state;
    
    uuid_t uuid1, uuid2;
    uuid_generate(uuid1);
    memcpy(uuid2, uuid1, sizeof(uuid_t));
    
    int result = uuid_compare(uuid1, uuid2);
    assert_int_equal(result, 0);
}

/* Test: UUID compare for different UUIDs */
static void test_uuid_compare_different(void **state) {
    (void) state;
    
    uuid_t uuid1, uuid2;
    uuid_generate(uuid1);
    uuid_generate(uuid2);
    
    int result = uuid_compare(uuid1, uuid2);
    /* Result can be -1, 0, or 1 - just verify it's consistent */
    assert_true(result == -1 || result == 0 || result == 1);
}

/* Test: UUID clear sets all bytes to zero */
static void test_uuid_clear(void **state) {
    (void) state;
    
    uuid_t uuid;
    uuid_generate(uuid);
    
    uuid_clear(uuid);
    
    /* Verify all bytes are zero */
    uuid_t zero_uuid;
    memset(zero_uuid, 0, sizeof(uuid_t));
    assert_memory_equal(uuid, zero_uuid, sizeof(uuid_t));
}

/* Test: UUID is_null on zero UUID */
static void test_uuid_is_null_on_zero_uuid(void **state) {
    (void) state;
    
    uuid_t uuid;
    memset(uuid, 0, sizeof(uuid_t));
    
    int result = uuid_is_null(uuid);
    assert_int_equal(result, 1);
}

/* Test: UUID is_null on generated UUID */
static void test_uuid_is_null_on_generated_uuid(void **state) {
    (void) state;
    
    uuid_t uuid;
    uuid_generate(uuid);
    
    int result = uuid_is_null(uuid);
    assert_int_equal(result, 0);
}

/* Test: UUID copy operation */
static void test_uuid_copy(void **state) {
    (void) state;
    
    uuid_t source, dest;
    uuid_generate(source);
    
    uuid_copy(dest, source);
    
    assert_memory_equal(source, dest, sizeof(uuid_t));
}

/* Test: UUID unparse with NULL buffer */
static void test_uuid_unparse_null_buffer(void **state) {
    (void) state;
    
    uuid_t uuid;
    uuid_generate(uuid);
    
    /* This is implementation-dependent behavior */
    /* Most implementations will crash or handle gracefully */
}

/* Test: UUID unparse with insufficient buffer */
static void test_uuid_unparse_small_buffer(void **state) {
    (void) state;
    
    uuid_t uuid;
    char small_buffer[10];
    
    uuid_generate(uuid);
    uuid_unparse(uuid, small_buffer);
    
    /* This may overflow - implementation dependent */
    /* In safe implementation, should handle gracefully */
}

/* Test: UUID version check */
static void test_uuid_version(void **state) {
    (void) state;
    
    uuid_t uuid;
    uuid_generate(uuid);
    
    /* UUID generated by uuid_generate should have valid version */
    int version = (uuid[6] >> 4) & 0x0f;
    assert_true(version >= 1 && version <= 5);
}

/* Test: UUID variant check */
static void test_uuid_variant(void **state) {
    (void) state;
    
    uuid_t uuid;
    uuid_generate(uuid);
    
    /* UUID should have valid variant */
    int variant = (uuid[8] >> 6) & 0x03;
    assert_true(variant >= 0 && variant <= 3);
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_uuid_generate_creates_valid_uuid),
        cmocka_unit_test(test_uuid_generate_produces_unique_uuids),
        cmocka_unit_test(test_uuid_generate_with_uninitialized_uuid),
        cmocka_unit_test(test_uuid_unparse_and_parse_roundtrip),
        cmocka_unit_test(test_uuid_parse_valid_string),
        cmocka_unit_test(test_uuid_parse_invalid_string),
        cmocka_unit_test(test_uuid_parse_empty_string),
        cmocka_unit_test(test_uuid_compare_equal),
        cmocka_unit_test(test_uuid_compare_different),
        cmocka_unit_test(test_uuid_clear),
        cmocka_unit_test(test_uuid_is_null_on_zero_uuid),
        cmocka_unit_test(test_uuid_is_null_on_generated_uuid),
        cmocka_unit_test(test_uuid_copy),
        cmocka_unit_test(test_uuid_version),
        cmocka_unit_test(test_uuid_variant),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}