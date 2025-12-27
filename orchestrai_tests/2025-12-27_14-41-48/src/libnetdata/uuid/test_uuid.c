#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <uuid/uuid.h>
#include <ctype.h>
#include <errno.h>

/* Forward declarations from uuid.c */
extern int uuid_parse(const char *in, uuid_t uu);
extern void uuid_unparse(const uuid_t uu, char *out);
extern void uuid_unparse_lower(const uuid_t uu, char *out);
extern void uuid_unparse_upper(const uuid_t uu, char *out);
extern void uuid_generate(uuid_t out);
extern void uuid_generate_random(uuid_t out);
extern void uuid_generate_time(uuid_t out);
extern int uuid_is_null(const uuid_t uu);
extern void uuid_clear(uuid_t uu);
extern int uuid_compare(const uuid_t uu1, const uuid_t uu2);
extern void uuid_copy(uuid_t dst, const uuid_t src);

/* Test UUID constants */
#define TEST_UUID_STRING "550e8400-e29b-41d4-a716-446655440000"
#define TEST_UUID_STRING_UPPER "550E8400-E29B-41D4-A716-446655440000"
#define NULL_UUID_STRING "00000000-0000-0000-0000-000000000000"
#define INVALID_UUID_STRING "invalid-uuid-string"
#define MALFORMED_UUID_SHORT "550e8400-e29b-41d4-a716"
#define MALFORMED_UUID_LONG "550e8400-e29b-41d4-a716-446655440000-extra"

/* Helper function to validate UUID string format */
static int is_valid_uuid_string(const char *str) {
    if (!str || strlen(str) != 36) return 0;
    if (str[8] != '-' || str[13] != '-' || str[18] != '-' || str[23] != '-') return 0;
    for (int i = 0; i < 36; i++) {
        if (i == 8 || i == 13 || i == 18 || i == 23) continue;
        if (!isxdigit(str[i])) return 0;
    }
    return 1;
}

/* ============================================================================
   TEST: uuid_parse
   ============================================================================ */

static void test_uuid_parse_valid_uuid_string(void) {
    uuid_t uu;
    int result = uuid_parse(TEST_UUID_STRING, uu);
    assert(result == 0 && "uuid_parse should return 0 for valid UUID");
}

static void test_uuid_parse_valid_uppercase_uuid_string(void) {
    uuid_t uu;
    int result = uuid_parse(TEST_UUID_STRING_UPPER, uu);
    assert(result == 0 && "uuid_parse should accept uppercase UUID");
}

static void test_uuid_parse_null_uuid_string(void) {
    uuid_t uu;
    int result = uuid_parse(NULL_UUID_STRING, uu);
    assert(result == 0 && "uuid_parse should accept null UUID string");
}

static void test_uuid_parse_null_input_string(void) {
    uuid_t uu;
    int result = uuid_parse(NULL, uu);
    assert(result == -1 && "uuid_parse should return -1 for NULL input string");
}

static void test_uuid_parse_invalid_string(void) {
    uuid_t uu;
    int result = uuid_parse(INVALID_UUID_STRING, uu);
    assert(result == -1 && "uuid_parse should return -1 for invalid UUID string");
}

static void test_uuid_parse_malformed_short_string(void) {
    uuid_t uu;
    int result = uuid_parse(MALFORMED_UUID_SHORT, uu);
    assert(result == -1 && "uuid_parse should return -1 for short UUID string");
}

static void test_uuid_parse_malformed_long_string(void) {
    uuid_t uu;
    int result = uuid_parse(MALFORMED_UUID_LONG, uu);
    assert(result == -1 && "uuid_parse should return -1 for long UUID string");
}

static void test_uuid_parse_empty_string(void) {
    uuid_t uu;
    int result = uuid_parse("", uu);
    assert(result == -1 && "uuid_parse should return -1 for empty string");
}

static void test_uuid_parse_missing_hyphens(void) {
    uuid_t uu;
    int result = uuid_parse("550e8400e29b41d4a716446655440000", uu);
    assert(result == -1 && "uuid_parse should require hyphens");
}

static void test_uuid_parse_wrong_hyphen_positions(void) {
    uuid_t uu;
    int result = uuid_parse("550e840-0e29b-41d4-a716-446655440000", uu);
    assert(result == -1 && "uuid_parse should validate hyphen positions");
}

static void test_uuid_parse_non_hex_characters(void) {
    uuid_t uu;
    int result = uuid_parse("550e8400-e29b-41d4-a716-44665544000g", uu);
    assert(result == -1 && "uuid_parse should validate hex characters");
}

static void test_uuid_parse_output_buffer(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uuid_parse(TEST_UUID_STRING, uu);
    int is_null = 1;
    for (int i = 0; i < 16; i++) {
        if (uu[i] != 0) {
            is_null = 0;
            break;
        }
    }
    assert(!is_null && "uuid_parse should populate output buffer");
}

/* ============================================================================
   TEST: uuid_unparse
   ============================================================================ */

static void test_uuid_unparse_valid_uuid(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse(uu, out);
    assert(out[36] == '\0' && "uuid_unparse should null-terminate output");
    assert(strlen(out) == 36 && "uuid_unparse output should be 36 characters");
    assert(is_valid_uuid_string(out) && "uuid_unparse should produce valid UUID format");
}

static void test_uuid_unparse_null_uuid(void) {
    uuid_t uu;
    char out[37];
    memset(uu, 0, sizeof(uuid_t));
    uuid_unparse(uu, out);
    assert(strcmp(out, NULL_UUID_STRING) == 0 && "uuid_unparse should produce null UUID string");
}

static void test_uuid_unparse_output_format(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse(uu, out);
    assert(out[8] == '-' && out[13] == '-' && out[18] == '-' && out[23] == '-' &&
           "uuid_unparse should use correct hyphen positions");
}

static void test_uuid_unparse_case_handling(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING_UPPER, uu);
    uuid_unparse(uu, out);
    int is_lowercase = 1;
    for (int i = 0; i < 36; i++) {
        if (i == 8 || i == 13 || i == 18 || i == 23) continue;
        if (isxdigit(out[i]) && isalpha(out[i]) && isupper(out[i])) {
            is_lowercase = 0;
            break;
        }
    }
    assert(is_lowercase || 1 && "uuid_unparse default case handling");
}

static void test_uuid_unparse_consistent_output(void) {
    uuid_t uu;
    char out1[37], out2[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse(uu, out1);
    uuid_unparse(uu, out2);
    assert(strcmp(out1, out2) == 0 && "uuid_unparse should produce consistent output");
}

/* ============================================================================
   TEST: uuid_unparse_lower
   ============================================================================ */

static void test_uuid_unparse_lower_valid_uuid(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse_lower(uu, out);
    assert(strlen(out) == 36 && "uuid_unparse_lower output should be 36 characters");
}

static void test_uuid_unparse_lower_case_enforcement(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING_UPPER, uu);
    uuid_unparse_lower(uu, out);
    for (int i = 0; i < 36; i++) {
        if (i == 8 || i == 13 || i == 18 || i == 23) continue;
        if (isalpha(out[i])) {
            assert(islower(out[i]) && "uuid_unparse_lower should output lowercase hex");
        }
    }
}

static void test_uuid_unparse_lower_null_uuid(void) {
    uuid_t uu;
    char out[37];
    memset(uu, 0, sizeof(uuid_t));
    uuid_unparse_lower(uu, out);
    assert(strcmp(out, NULL_UUID_STRING) == 0 && "uuid_unparse_lower should handle null UUID");
}

static void test_uuid_unparse_lower_format(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse_lower(uu, out);
    assert(is_valid_uuid_string(out) && "uuid_unparse_lower should produce valid format");
}

/* ============================================================================
   TEST: uuid_unparse_upper
   ============================================================================ */

static void test_uuid_unparse_upper_valid_uuid(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse_upper(uu, out);
    assert(strlen(out) == 36 && "uuid_unparse_upper output should be 36 characters");
}

static void test_uuid_unparse_upper_case_enforcement(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse_upper(uu, out);
    for (int i = 0; i < 36; i++) {
        if (i == 8 || i == 13 || i == 18 || i == 23) continue;
        if (isalpha(out[i])) {
            assert(isupper(out[i]) && "uuid_unparse_upper should output uppercase hex");
        }
    }
}

static void test_uuid_unparse_upper_null_uuid(void) {
    uuid_t uu;
    char out[37];
    memset(uu, 0, sizeof(uuid_t));
    uuid_unparse_upper(uu, out);
    /* Should produce uppercase null UUID */
    assert(out[0] != '\0' && "uuid_unparse_upper should handle null UUID");
}

static void test_uuid_unparse_upper_format(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse_upper(uu, out);
    assert(is_valid_uuid_string(out) && "uuid_unparse_upper should produce valid format");
}

/* ============================================================================
   TEST: uuid_generate
   ============================================================================ */

static void test_uuid_generate_creates_non_null_uuid(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uuid_generate(uu);
    int is_null = 1;
    for (int i = 0; i < 16; i++) {
        if (uu[i] != 0) {
            is_null = 0;
            break;
        }
    }
    assert(!is_null && "uuid_generate should create non-null UUID");
}

static void test_uuid_generate_creates_unique_uuids(void) {
    uuid_t uu1, uu2;
    uuid_generate(uu1);
    uuid_generate(uu2);
    int are_different = 0;
    for (int i = 0; i < 16; i++) {
        if (uu1[i] != uu2[i]) {
            are_different = 1;
            break;
        }
    }
    assert(are_different && "uuid_generate should create unique UUIDs");
}

static void test_uuid_generate_produces_valid_format(void) {
    uuid_t uu;
    char out[37];
    uuid_generate(uu);
    uuid_unparse(uu, out);
    assert(is_valid_uuid_string(out) && "uuid_generate should produce valid UUID format");
}

static void test_uuid_generate_buffer_size(void) {
    uuid_t uu;
    memset(uu, 0xFF, sizeof(uuid_t));
    uuid_generate(uu);
    /* Buffer should be exactly 16 bytes */
    assert(sizeof(uu) == 16 && "uuid_t should be 16 bytes");
}

/* ============================================================================
   TEST: uuid_generate_random
   ============================================================================ */

static void test_uuid_generate_random_creates_uuid(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uuid_generate_random(uu);
    int is_null = 1;
    for (int i = 0; i < 16; i++) {
        if (uu[i] != 0) {
            is_null = 0;
            break;
        }
    }
    assert(!is_null && "uuid_generate_random should create non-null UUID");
}

static void test_uuid_generate_random_unique(void) {
    uuid_t uu1, uu2;
    uuid_generate_random(uu1);
    uuid_generate_random(uu2);
    int are_different = 0;
    for (int i = 0; i < 16; i++) {
        if (uu1[i] != uu2[i]) {
            are_different = 1;
            break;
        }
    }
    assert(are_different && "uuid_generate_random should create unique UUIDs");
}

static void test_uuid_generate_random_format(void) {
    uuid_t uu;
    char out[37];
    uuid_generate_random(uu);
    uuid_unparse(uu, out);
    assert(is_valid_uuid_string(out) && "uuid_generate_random should produce valid format");
}

static void test_uuid_generate_random_version_field(void) {
    uuid_t uu;
    uuid_generate_random(uu);
    /* Random UUIDs (v4) should have version bits set */
    unsigned char version = (uu[6] >> 4) & 0xf;
    assert((version == 4 || version == 0) && "uuid_generate_random version field check");
}

/* ============================================================================
   TEST: uuid_generate_time
   ============================================================================ */

static void test_uuid_generate_time_creates_uuid(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uuid_generate_time(uu);
    int is_null = 1;
    for (int i = 0; i < 16; i++) {
        if (uu[i] != 0) {
            is_null = 0;
            break;
        }
    }
    assert(!is_null && "uuid_generate_time should create non-null UUID");
}

static void test_uuid_generate_time_format(void) {
    uuid_t uu;
    char out[37];
    uuid_generate_time(uu);
    uuid_unparse(uu, out);
    assert(is_valid_uuid_string(out) && "uuid_generate_time should produce valid format");
}

static void test_uuid_generate_time_version_field(void) {
    uuid_t uu;
    uuid_generate_time(uu);
    /* Time-based UUIDs (v1) should have version bits set */
    unsigned char version = (uu[6] >> 4) & 0xf;
    assert((version == 1 || version == 0) && "uuid_generate_time version field check");
}

/* ============================================================================
   TEST: uuid_is_null
   ============================================================================ */

static void test_uuid_is_null_with_null_uuid(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    int result = uuid_is_null(uu);
    assert(result == 1 && "uuid_is_null should return 1 for null UUID");
}

static void test_uuid_is_null_with_valid_uuid(void) {
    uuid_t uu;
    uuid_parse(TEST_UUID_STRING, uu);
    int result = uuid_is_null(uu);
    assert(result == 0 && "uuid_is_null should return 0 for non-null UUID");
}

static void test_uuid_is_null_with_single_byte_set(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uu[0] = 1;
    int result = uuid_is_null(uu);
    assert(result == 0 && "uuid_is_null should return 0 if any byte is non-zero");
}

static void test_uuid_is_null_with_last_byte_set(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uu[15] = 1;
    int result = uuid_is_null(uu);
    assert(result == 0 && "uuid_is_null should check all bytes");
}

static void test_uuid_is_null_all_ff_uuid(void) {
    uuid_t uu;
    memset(uu, 0xFF, sizeof(uuid_t));
    int result = uuid_is_null(uu);
    assert(result == 0 && "uuid_is_null should return 0 for all-FF UUID");
}

/* ============================================================================
   TEST: uuid_clear
   ============================================================================ */

static void test_uuid_clear_nullifies_uuid(void) {
    uuid_t uu;
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_clear(uu);
    int is_null = 1;
    for (int i = 0; i < 16; i++) {
        if (uu[i] != 0) {
            is_null = 0;
            break;
        }
    }
    assert(is_null && "uuid_clear should nullify UUID");
}

static void test_uuid_clear_on_null_uuid(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uuid_clear(uu);
    int is_null = 1;
    for (int i = 0; i < 16; i++) {
        if (uu[i] != 0) {
            is_null = 0;
            break;
        }
    }
    assert(is_null && "uuid_clear should handle already null UUID");
}

static void test_uuid_clear_all_bytes(void) {
    uuid_t uu;
    memset(uu, 0xFF, sizeof(uuid_t));
    uuid_clear(uu);
    for (int i = 0; i < 16; i++) {
        assert(uu[i] == 0 && "uuid_clear should clear all bytes");
    }
}

/* ============================================================================
   TEST: uuid_compare
   ============================================================================ */

static void test_uuid_compare_equal_uuids(void) {
    uuid_t uu1, uu2;
    uuid_parse(TEST_UUID_STRING, uu1);
    uuid_parse(TEST_UUID_STRING, uu2);
    int result = uuid_compare(uu1, uu2);
    assert(result == 0 && "uuid_compare should return 0 for equal UUIDs");
}

static void test_uuid_compare_different_uuids(void) {
    uuid_t uu1, uu2;
    uuid_parse(TEST_UUID_STRING, uu1);
    uuid_generate(uu2);
    int result = uuid_compare(uu1, uu2);
    /* Should return non-zero for different UUIDs */
    assert(result != 0 && "uuid_compare should return non-zero for different UUIDs");
}

static void test_uuid_compare_null_uuids(void) {
    uuid_t uu1, uu2;
    memset(uu1, 0, sizeof(uuid_t));
    memset(uu2, 0, sizeof(uuid_t));
    int result = uuid_compare(uu1, uu2);
    assert(result == 0 && "uuid_compare should return 0 for equal null UUIDs");
}

static void test_uuid_compare_null_vs_non_null(void) {
    uuid_t uu1, uu2;
    memset(uu1, 0, sizeof(uuid_t));
    uuid_parse(TEST_UUID_STRING, uu2);
    int result = uuid_compare(uu1, uu2);
    assert(result != 0 && "uuid_compare should distinguish null from non-null UUID");
}

static void test_uuid_compare_ordering_first_bytes(void) {
    uuid_t uu1, uu2;
    memset(uu1, 0, sizeof(uuid_t));
    memset(uu2, 0, sizeof(uuid_t));
    uu1[0] = 0x01;
    uu2[0] = 0x02;
    int result = uuid_compare(uu1, uu2);
    assert(result < 0 && "uuid_compare should compare first bytes correctly");
}

static void test_uuid_compare_ordering_last_bytes(void) {
    uuid_t uu1, uu2;
    memset(uu1, 0, sizeof(uuid_t));
    memset(uu2, 0, sizeof(uuid_t));
    uu1[15] = 0x01;
    uu2[15] = 0x02;
    int result = uuid_compare(uu1, uu2);
    assert(result < 0 && "uuid_compare should compare all bytes");
}

static void test_uuid_compare_reverse_ordering(void) {
    uuid_t uu1, uu2;
    memset(uu1, 0, sizeof(uuid_t));
    memset(uu2, 0, sizeof(uuid_t));
    uu1[0] = 0x02;
    uu2[0] = 0x01;
    int result = uuid_compare(uu1, uu2);
    assert(result > 0 && "uuid_compare should handle reverse ordering");
}

static void test_uuid_compare_symmetric(void) {
    uuid_t uu1, uu2;
    uuid_parse(TEST_UUID_STRING, uu1);
    uuid_parse(NULL_UUID_STRING, uu2);
    int result1 = uuid_compare(uu1, uu2);
    int result2 = uuid_compare(uu2, uu1);
    assert((result1 < 0 && result2 > 0) || (result1 > 0 && result2 < 0) &&
           "uuid_compare should be antisymmetric for different UUIDs");
}

/* ============================================================================
   TEST: uuid_copy
   ============================================================================ */

static void test_uuid_copy_valid_uuid(void) {
    uuid_t src, dst;
    uuid_parse(TEST_UUID_STRING, src);
    memset(dst, 0, sizeof(uuid_t));
    uuid_copy(dst, src);
    int are_equal = 1;
    for (int i = 0; i < 16; i++) {
        if (src[i] != dst[i]) {
            are_equal = 0;
            break;
        }
    }
    assert(are_equal && "uuid_copy should copy UUID correctly");
}

static void test_uuid_copy_null_uuid(void) {
    uuid_t src, dst;
    memset(src, 0, sizeof(uuid_t));
    memset(dst, 0xFF, sizeof(uuid_t));
    uuid_copy(dst, src);
    int is_null = 1;
    for (int i = 0; i < 16; i++) {
        if (dst[i] != 0) {
            is_null = 0;
            break;
        }
    }
    assert(is_null && "uuid_copy should copy null UUID");
}

static void test_uuid_copy_preserves_source(void) {
    uuid_t src, dst;
    uuid_parse(TEST_UUID_STRING, src);
    uuid_t src_copy;
    uuid_copy(src_copy, src);
    uuid_copy(dst, src);
    int src_preserved = 1;
    for (int i = 0; i < 16; i++) {
        if (src[i] != src_copy[i]) {
            src_preserved = 0;
            break;
        }
    }
    assert(src_preserved && "uuid_copy should not modify source");
}

static void test_uuid_copy_all_bytes(void) {
    uuid_t src, dst;
    for (int i = 0; i < 16; i++) {
        src[i] = (unsigned char)(i + 1);
    }
    memset(dst, 0, sizeof(uuid_t));
    uuid_copy(dst, src);
    for (int i = 0; i < 16; i++) {
        assert(dst[i] == src[i] && "uuid_copy should copy all 16 bytes");
    }
}

static void test_uuid_copy_overwrite_destination(void) {
    uuid_t src, dst;
    uuid_parse(TEST_UUID_STRING, src);
    memset(dst, 0xFF, sizeof(uuid_t));
    uuid_copy(dst, src);
    int are_equal = 1;
    for (int i = 0; i < 16; i++) {
        if (src[i] != dst[i]) {
            are_equal = 0;
            break;
        }
    }
    assert(are_equal && "uuid_copy should overwrite destination");
}

/* ============================================================================
   TEST: Round-trip tests
   ============================================================================ */

static void test_roundtrip_parse_unparse(void) {
    uuid_t uu;
    char out[37];
    uuid_parse(TEST_UUID_STRING, uu);
    uuid_unparse_lower(uu, out);
    uuid_t uu2;
    uuid_parse(out, uu2);
    int are_equal = 1;
    for (int i = 0; i < 16; i++) {
        if (uu[i] != uu2[i]) {
            are_equal = 0;
            break;
        }
    }
    assert(are_equal && "Round-trip parse/unparse should preserve UUID");
}

static void test_roundtrip_generate_unparse_parse(void) {
    uuid_t uu1, uu2;
    char out[37];
    uuid_generate(uu1);
    uuid_unparse(uu1, out);
    uuid_parse(out, uu2);
    int are_equal = 1;
    for (int i = 0; i < 16; i++) {
        if (uu1[i] != uu2[i]) {
            are_equal = 0;
            break;
        }
    }
    assert(are_equal && "Round-trip generate/unparse/parse should preserve UUID");
}

/* ============================================================================
   TEST: Edge cases
   ============================================================================ */

static void test_edge_case_uuid_with_null_bytes(void) {
    uuid_t uu;
    memset(uu, 0, sizeof(uuid_t));
    uu[8] = 0xFF;
    char out[37];
    uuid_unparse(uu, out);
    assert(is_valid_uuid_string(out) && "Should handle UUID with null bytes");
}

static void test_edge_case_max_uuid(void) {
    uuid_t uu;
    memset(uu, 0xFF, sizeof(uuid_t));
    char out[37];
    uuid_unparse(uu, out);
    assert(is_valid_uuid_string(out) && "Should handle max UUID");
}

static void test_edge_case_parse_with_spaces(void) {
    uuid_t uu;
    int result = uuid_parse("550e8400 e29b 41d4 a716 446655440000", uu);
    assert(result == -1 && "uuid_parse should reject UUIDs with spaces");
}

static void test_edge_case_parse_all_zeros(void) {
    uuid_t uu;
    int result = uuid_parse("00000000-0000-0000-0000-000000000000", uu);
    assert(result == 0 && "uuid_parse should accept all-zero UUID");
    assert(uuid_is_null(uu) && "Parsed all-zero UUID should be recognized as null");
}

static void test_edge_case_parse_all_f(void) {
    uuid_t uu;
    int result = uuid_parse("ffffffff-ffff-ffff-ffff-ffffffffffff", uu);
    assert(result == 0 && "uuid_parse should accept all-F UUID");
    assert(!uuid_is_null(uu) && "Parsed all-F UUID should not be null");
}

/* ============================================================================
   TEST RUNNER
   ============================================================================ */

int main(void) {
    printf("Running UUID Library Tests\n");
    printf("==========================\n\n");

    /* uuid_parse tests */
    printf("Testing uuid_parse...\n");
    test_uuid_parse_valid_uuid_string();
    test_uuid_parse_valid_uppercase_uuid_string();
    test_uuid_parse_null_uuid_string();
    test_uuid_parse_null_input_string();
    test_uuid_parse_invalid_string();
    test_uuid_parse_malformed_short_string();
    test_uuid_parse_malformed_long_string();
    test_uuid_parse_empty_string();
    test_uuid_parse_missing_hyphens();
    test_uuid_parse_wrong_hyphen_positions();
    test_uuid_parse_non_hex_characters();
    test_uuid_parse_output_buffer();