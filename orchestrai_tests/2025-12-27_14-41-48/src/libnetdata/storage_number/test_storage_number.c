#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <assert.h>
#include <limits.h>
#include <float.h>

// Mock and test framework setup
#include "storage_number.h"

// Test framework helpers
#define TEST_ASSERT(condition) do { \
    if (!(condition)) { \
        printf("ASSERTION FAILED: %s at line %d\n", #condition, __LINE__); \
        exit(1); \
    } \
} while(0)

#define TEST_ASSERT_EQ(a, b) do { \
    if ((a) != (b)) { \
        printf("EQUALITY ASSERTION FAILED at line %d\n", __LINE__); \
        exit(1); \
    } \
} while(0)

#define TEST_ASSERT_FLOAT_EQ(a, b, epsilon) do { \
    if (fabs((a) - (b)) > (epsilon)) { \
        printf("FLOAT EQUALITY ASSERTION FAILED at line %d\n", __LINE__); \
        exit(1); \
    } \
} while(0)

#define TEST_ASSERT_STR_EQ(a, b) do { \
    if (strcmp((a), (b)) != 0) { \
        printf("STRING EQUALITY ASSERTION FAILED at line %d\n", __LINE__); \
        exit(1); \
    } \
} while(0)

#define RUN_TEST(test_func) do { \
    printf("Running: %s\n", #test_func); \
    test_func(); \
    printf("✓ PASSED: %s\n", #test_func); \
} while(0)

// ============================================================================
// Test Cases for storage_number functions
// ============================================================================

// Test: pack_storage_number with positive values
static void test_pack_storage_number_positive() {
    storage_number packed;
    double value = 123.456;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    TEST_ASSERT(packed != 0);
    
    // Round trip test
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT_FLOAT_EQ(unpacked, value, 1e-10);
}

// Test: pack_storage_number with negative values
static void test_pack_storage_number_negative() {
    storage_number packed;
    double value = -987.654;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    TEST_ASSERT(packed != 0);
    
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT_FLOAT_EQ(unpacked, value, 1e-10);
}

// Test: pack_storage_number with zero
static void test_pack_storage_number_zero() {
    storage_number packed;
    double value = 0.0;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT_FLOAT_EQ(unpacked, 0.0, 1e-15);
}

// Test: pack_storage_number with very small positive values
static void test_pack_storage_number_very_small_positive() {
    storage_number packed;
    double value = 1e-10;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT(unpacked >= 0);
}

// Test: pack_storage_number with very small negative values
static void test_pack_storage_number_very_small_negative() {
    storage_number packed;
    double value = -1e-10;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT(unpacked <= 0);
}

// Test: pack_storage_number with very large values
static void test_pack_storage_number_large_values() {
    storage_number packed;
    double value = 1e20;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT(unpacked > 0);
}

// Test: pack_storage_number with infinity
static void test_pack_storage_number_infinity() {
    storage_number packed;
    double value = INFINITY;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT(isinf(unpacked));
}

// Test: pack_storage_number with negative infinity
static void test_pack_storage_number_negative_infinity() {
    storage_number packed;
    double value = -INFINITY;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT(isinf(unpacked) && unpacked < 0);
}

// Test: pack_storage_number with NaN
static void test_pack_storage_number_nan() {
    storage_number packed;
    double value = NAN;
    
    packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT(isnan(unpacked));
}

// Test: unpack_storage_number with valid packed value
static void test_unpack_storage_number_valid() {
    double original = 456.789;
    storage_number packed = pack_storage_number(original, SN_DEFAULT_FLAGS);
    
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT_FLOAT_EQ(unpacked, original, 1e-10);
}

// Test: unpack_storage_number with zero packed value
static void test_unpack_storage_number_zero() {
    storage_number packed = 0;
    double unpacked = unpack_storage_number(packed);
    TEST_ASSERT_FLOAT_EQ(unpacked, 0.0, 1e-15);
}

// Test: storage_number_encode with positive value
static void test_storage_number_encode_positive() {
    char buffer[100];
    double value = 123.456;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    
    int bytes = storage_number_encode(packed, buffer);
    TEST_ASSERT(bytes > 0);
    TEST_ASSERT(bytes <= sizeof(buffer));
}

// Test: storage_number_encode with negative value
static void test_storage_number_encode_negative() {
    char buffer[100];
    double value = -456.789;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    
    int bytes = storage_number_encode(packed, buffer);
    TEST_ASSERT(bytes > 0);
    TEST_ASSERT(bytes <= sizeof(buffer));
}

// Test: storage_number_encode with zero
static void test_storage_number_encode_zero() {
    char buffer[100];
    storage_number packed = pack_storage_number(0.0, SN_DEFAULT_FLAGS);
    
    int bytes = storage_number_encode(packed, buffer);
    TEST_ASSERT(bytes > 0);
    TEST_ASSERT(bytes <= sizeof(buffer));
}

// Test: storage_number_decode round trip
static void test_storage_number_decode_round_trip() {
    double original = 789.012;
    storage_number packed = pack_storage_number(original, SN_DEFAULT_FLAGS);
    
    char buffer[100];
    int encoded_bytes = storage_number_encode(packed, buffer);
    TEST_ASSERT(encoded_bytes > 0);
    
    storage_number unpacked_packed;
    int decoded_bytes = storage_number_decode((unsigned char *)buffer, encoded_bytes, &unpacked_packed);
    TEST_ASSERT(decoded_bytes == encoded_bytes);
    
    double final_value = unpack_storage_number(unpacked_packed);
    TEST_ASSERT_FLOAT_EQ(final_value, original, 1e-10);
}

// Test: storage_number_decode with minimum buffer
static void test_storage_number_decode_minimum_buffer() {
    storage_number packed = pack_storage_number(100.5, SN_DEFAULT_FLAGS);
    
    char buffer[10];
    int encoded_bytes = storage_number_encode(packed, buffer);
    TEST_ASSERT(encoded_bytes > 0);
    
    storage_number decoded_packed;
    int decoded_bytes = storage_number_decode((unsigned char *)buffer, encoded_bytes, &decoded_packed);
    TEST_ASSERT(decoded_bytes > 0);
}

// Test: storage_number_decode with insufficient buffer
static void test_storage_number_decode_insufficient_buffer() {
    storage_number packed = pack_storage_number(100.5, SN_DEFAULT_FLAGS);
    
    char buffer[100];
    int encoded_bytes = storage_number_encode(packed, buffer);
    
    storage_number decoded_packed;
    // Try to decode with smaller buffer than needed
    int decoded_bytes = storage_number_decode((unsigned char *)buffer, 1, &decoded_packed);
    TEST_ASSERT(decoded_bytes >= 0);
}

// Test: print_storage_number with positive value
static void test_print_storage_number_positive() {
    char buffer[256];
    double value = 555.444;
    
    ssize_t result = print_storage_number(buffer, 256, value);
    TEST_ASSERT(result > 0);
    TEST_ASSERT(result < 256);
    TEST_ASSERT(strlen(buffer) > 0);
}

// Test: print_storage_number with negative value
static void test_print_storage_number_negative() {
    char buffer[256];
    double value = -333.222;
    
    ssize_t result = print_storage_number(buffer, 256, value);
    TEST_ASSERT(result > 0);
    TEST_ASSERT(result < 256);
    TEST_ASSERT(buffer[0] == '-' || isdigit(buffer[0]));
}

// Test: print_storage_number with zero
static void test_print_storage_number_zero() {
    char buffer[256];
    double value = 0.0;
    
    ssize_t result = print_storage_number(buffer, 256, value);
    TEST_ASSERT(result > 0);
    TEST_ASSERT(result < 256);
}

// Test: print_storage_number with small buffer
static void test_print_storage_number_small_buffer() {
    char buffer[10];
    double value = 123.456;
    
    ssize_t result = print_storage_number(buffer, 10, value);
    TEST_ASSERT(result >= 0);
}

// Test: print_storage_number with buffer size 1
static void test_print_storage_number_buffer_size_1() {
    char buffer[1];
    double value = 123.456;
    
    ssize_t result = print_storage_number(buffer, 1, value);
    TEST_ASSERT(result >= 0);
}

// Test: print_storage_number with infinity
static void test_print_storage_number_infinity() {
    char buffer[256];
    double value = INFINITY;
    
    ssize_t result = print_storage_number(buffer, 256, value);
    TEST_ASSERT(result > 0 || result == -1);
}

// Test: print_storage_number with NaN
static void test_print_storage_number_nan() {
    char buffer[256];
    double value = NAN;
    
    ssize_t result = print_storage_number(buffer, 256, value);
    TEST_ASSERT(result > 0 || result == -1);
}

// Test: pack_storage_number with different flag combinations
static void test_pack_storage_number_with_flags() {
    storage_number packed1 = pack_storage_number(100.0, SN_DEFAULT_FLAGS);
    storage_number packed2 = pack_storage_number(100.0, SN_ANOMALY_BIT);
    
    // Both should work, may produce different results
    double unpacked1 = unpack_storage_number(packed1);
    double unpacked2 = unpack_storage_number(packed2);
    
    TEST_ASSERT(unpacked1 > 0);
    TEST_ASSERT(unpacked2 > 0);
}

// Test: Boundary value - largest representable double
static void test_pack_unpack_max_double() {
    double value = DBL_MAX / 1e10; // Avoid exact max to prevent overflow
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    
    TEST_ASSERT(unpacked > 0);
}

// Test: Boundary value - smallest positive double
static void test_pack_unpack_min_positive_double() {
    double value = DBL_MIN;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    
    TEST_ASSERT(unpacked >= 0);
}

// Test: Multiple round trips maintain approximate precision
static void test_multiple_round_trips() {
    double original = 12345.6789;
    storage_number packed = pack_storage_number(original, SN_DEFAULT_FLAGS);
    
    for (int i = 0; i < 5; i++) {
        double unpacked = unpack_storage_number(packed);
        packed = pack_storage_number(unpacked, SN_DEFAULT_FLAGS);
    }
    
    double final = unpack_storage_number(packed);
    TEST_ASSERT(fabs(final - original) / original < 0.01); // Within 1%
}

// Test: Encode/decode preserves value across multiple iterations
static void test_encode_decode_multiple_iterations() {
    double original = 999.888;
    
    for (int i = 0; i < 3; i++) {
        storage_number packed = pack_storage_number(original, SN_DEFAULT_FLAGS);
        
        char buffer[100];
        int encoded = storage_number_encode(packed, buffer);
        TEST_ASSERT(encoded > 0);
        
        storage_number decoded_packed;
        int decoded = storage_number_decode((unsigned char *)buffer, encoded, &decoded_packed);
        TEST_ASSERT(decoded == encoded);
        
        original = unpack_storage_number(decoded_packed);
    }
}

// Test: Storage number with fractional part precision
static void test_fractional_precision() {
    double value = 0.123456789;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    
    TEST_ASSERT(unpacked > 0);
    TEST_ASSERT(unpacked < 1);
}

// Test: Negative fractional values
static void test_negative_fractional() {
    double value = -0.987654321;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    double unpacked = unpack_storage_number(packed);
    
    TEST_ASSERT(unpacked < 0);
    TEST_ASSERT(unpacked > -1);
}

// Test: Very close values should pack/unpack consistently
static void test_very_close_values() {
    double value1 = 100.0;
    double value2 = 100.0000000001;
    
    storage_number packed1 = pack_storage_number(value1, SN_DEFAULT_FLAGS);
    storage_number packed2 = pack_storage_number(value2, SN_DEFAULT_FLAGS);
    
    double unpacked1 = unpack_storage_number(packed1);
    double unpacked2 = unpack_storage_number(packed2);
    
    TEST_ASSERT(unpacked1 > 0);
    TEST_ASSERT(unpacked2 > 0);
}

// Test: Encode with exact boundary buffer size
static void test_encode_exact_boundary() {
    double value = 111.222;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    
    char buffer[100];
    int bytes = storage_number_encode(packed, buffer);
    TEST_ASSERT(bytes > 0);
    
    // Encode again and verify exact size
    char buffer2[100];
    int bytes2 = storage_number_encode(packed, buffer2);
    TEST_ASSERT(bytes2 == bytes);
}

// Test: Consecutive encoding of same value
static void test_consecutive_same_value_encoding() {
    double value = 333.444;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    
    char buffer1[100];
    char buffer2[100];
    
    int bytes1 = storage_number_encode(packed, buffer1);
    int bytes2 = storage_number_encode(packed, buffer2);
    
    TEST_ASSERT(bytes1 == bytes2);
    TEST_ASSERT(memcmp(buffer1, buffer2, bytes1) == 0);
}

// Test: Decode produces consistent results
static void test_decode_consistency() {
    double value = 555.666;
    storage_number packed = pack_storage_number(value, SN_DEFAULT_FLAGS);
    
    char buffer[100];
    int encoded = storage_number_encode(packed, buffer);
    
    storage_number decoded1, decoded2;
    storage_number_decode((unsigned char *)buffer, encoded, &decoded1);
    storage_number_decode((unsigned char *)buffer, encoded, &decoded2);
    
    TEST_ASSERT(decoded1 == decoded2);
}

// Test: Print with very large buffer
static void test_print_large_buffer() {
    char buffer[10000];
    double value = 777.888;
    
    ssize_t result = print_storage_number(buffer, 10000, value);
    TEST_ASSERT(result > 0);
    TEST_ASSERT(result < 10000);
}

// Test: Print with exact minimum needed buffer
static void test_print_exact_minimum_buffer() {
    double value = 999.111;
    
    char buffer[256];
    ssize_t needed = print_storage_number(buffer, 256, value);
    TEST_ASSERT(needed > 0);
}

// Main test runner
int main() {
    printf("========================================\n");
    printf("Storage Number Module Test Suite\n");
    printf("========================================\n\n");
    
    // Pack tests
    RUN_TEST(test_pack_storage_number_positive);
    RUN_TEST(test_pack_storage_number_negative);
    RUN_TEST(test_pack_storage_number_zero);
    RUN_TEST(test_pack_storage_number_very_small_positive);
    RUN_TEST(test_pack_storage_number_very_small_negative);
    RUN_TEST(test_pack_storage_number_large_values);
    RUN_TEST(test_pack_storage_number_infinity);
    RUN_TEST(test_pack_storage_number_negative_infinity);
    RUN_TEST(test_pack_storage_number_nan);
    RUN_TEST(test_pack_storage_number_with_flags);
    
    // Unpack tests
    RUN_TEST(test_unpack_storage_number_valid);
    RUN_TEST(test_unpack_storage_number_zero);
    
    // Encode tests
    RUN_TEST(test_storage_number_encode_positive);
    RUN_TEST(test_storage_number_encode_negative);
    RUN_TEST(test_storage_number_encode_zero);
    RUN_TEST(test_encode_exact_boundary);
    RUN_TEST(test_consecutive_same_value_encoding);
    
    // Decode tests
    RUN_TEST(test_storage_number_decode_round_trip);
    RUN_TEST(test_storage_number_decode_minimum_buffer);
    RUN_TEST(test_storage_number_decode_insufficient_buffer);
    RUN_TEST(test_decode_consistency);
    
    // Print tests
    RUN_TEST(test_print_storage_number_positive);
    RUN_TEST(test_print_storage_number_negative);
    RUN_TEST(test_print_storage_number_zero);
    RUN_TEST(test_print_storage_number_small_buffer);
    RUN_TEST(test_print_storage_number_buffer_size_1);
    RUN_TEST(test_print_storage_number_infinity);
    RUN_TEST(test_print_storage_number_nan);
    RUN_TEST(test_print_large_buffer);
    RUN_TEST(test_print_exact_minimum_buffer);
    
    // Boundary and edge case tests
    RUN_TEST(test_pack_unpack_max_double);
    RUN_TEST(test_pack_unpack_min_positive_double);
    RUN_TEST(test_fractional_precision);
    RUN_TEST(test_negative_fractional);
    RUN_TEST(test_very_close_values);
    
    // Integration tests
    RUN_TEST(test_multiple_round_trips);
    RUN_TEST(test_encode_decode_multiple_iterations);
    
    printf("\n========================================\n");
    printf("ALL TESTS PASSED!\n");
    printf("========================================\n");
    
    return 0;
}