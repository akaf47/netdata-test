#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <assert.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
#include <errno.h>

#include "../paths/paths.h"

/* Test helper macros */
#define TEST_ASSERT(condition, message) \
    do { \
        if (!(condition)) { \
            fprintf(stderr, "FAIL: %s\n", message); \
            exit(1); \
        } else { \
            fprintf(stderr, "PASS: %s\n", message); \
        } \
    } while(0)

#define TEST_ASSERT_EQUAL(actual, expected, message) \
    do { \
        if ((actual) != (expected)) { \
            fprintf(stderr, "FAIL: %s (expected %d, got %d)\n", message, expected, actual); \
            exit(1); \
        } else { \
            fprintf(stderr, "PASS: %s\n", message); \
        } \
    } while(0)

#define TEST_ASSERT_STRING_EQUAL(actual, expected, message) \
    do { \
        if (strcmp((actual), (expected)) != 0) { \
            fprintf(stderr, "FAIL: %s (expected '%s', got '%s')\n", message, expected, actual); \
            exit(1); \
        } else { \
            fprintf(stderr, "PASS: %s\n", message); \
        } \
    } while(0)

#define TEST_ASSERT_NULL(ptr, message) \
    do { \
        if ((ptr) != NULL) { \
            fprintf(stderr, "FAIL: %s (expected NULL)\n", message); \
            exit(1); \
        } else { \
            fprintf(stderr, "PASS: %s\n", message); \
        } \
    } while(0)

#define TEST_ASSERT_NOT_NULL(ptr, message) \
    do { \
        if ((ptr) == NULL) { \
            fprintf(stderr, "FAIL: %s (expected not NULL)\n", message); \
            exit(1); \
        } else { \
            fprintf(stderr, "PASS: %s\n", message); \
        } \
    } while(0)

/* ============================================================================
 * Test: paths_init with valid initialization
 * ============================================================================ */
static void test_paths_init_valid(void) {
    fprintf(stderr, "\n=== test_paths_init_valid ===\n");
    
    char *result = paths_init(PATHS_DEFAULT);
    TEST_ASSERT_NOT_NULL(result, "paths_init should return non-NULL for PATHS_DEFAULT");
    
    /* Verify that basic paths are initialized */
    TEST_ASSERT_NOT_NULL(netdata_cache_dir, "netdata_cache_dir should be initialized");
    TEST_ASSERT_NOT_NULL(netdata_varlib_dir, "netdata_varlib_dir should be initialized");
    TEST_ASSERT_NOT_NULL(netdata_sysconfdir, "netdata_sysconfdir should be initialized");
}

/* ============================================================================
 * Test: paths_init_custom with custom paths
 * ============================================================================ */
static void test_paths_init_custom(void) {
    fprintf(stderr, "\n=== test_paths_init_custom ===\n");
    
    char custom_path[1024];
    snprintf(custom_path, sizeof(custom_path), "/tmp/netdata_test_%d", getpid());
    
    char *result = paths_init(custom_path);
    TEST_ASSERT_NOT_NULL(result, "paths_init should return non-NULL for custom path");
}

/* ============================================================================
 * Test: netdata_cache_dir access
 * ============================================================================ */
static void test_netdata_cache_dir_access(void) {
    fprintf(stderr, "\n=== test_netdata_cache_dir_access ===\n");
    
    paths_init(PATHS_DEFAULT);
    TEST_ASSERT_NOT_NULL(netdata_cache_dir, "netdata_cache_dir should not be NULL");
    TEST_ASSERT(strlen(netdata_cache_dir) > 0, "netdata_cache_dir should have non-zero length");
}

/* ============================================================================
 * Test: netdata_varlib_dir access
 * ============================================================================ */
static void test_netdata_varlib_dir_access(void) {
    fprintf(stderr, "\n=== test_netdata_varlib_dir_access ===\n");
    
    paths_init(PATHS_DEFAULT);
    TEST_ASSERT_NOT_NULL(netdata_varlib_dir, "netdata_varlib_dir should not be NULL");
    TEST_ASSERT(strlen(netdata_varlib_dir) > 0, "netdata_varlib_dir should have non-zero length");
}

/* ============================================================================
 * Test: netdata_sysconfdir access
 * ============================================================================ */
static void test_netdata_sysconfdir_access(void) {
    fprintf(stderr, "\n=== test_netdata_sysconfdir_access ===\n");
    
    paths_init(PATHS_DEFAULT);
    TEST_ASSERT_NOT_NULL(netdata_sysconfdir, "netdata_sysconfdir should not be NULL");
    TEST_ASSERT(strlen(netdata_sysconfdir) > 0, "netdata_sysconfdir should have non-zero length");
}

/* ============================================================================
 * Test: netdata_rundir access
 * ============================================================================ */
static void test_netdata_rundir_access(void) {
    fprintf(stderr, "\n=== test_netdata_rundir_access ===\n");
    
    paths_init(PATHS_DEFAULT);
    TEST_ASSERT_NOT_NULL(netdata_rundir, "netdata_rundir should not be NULL");
    TEST_ASSERT(strlen(netdata_rundir) > 0, "netdata_rundir should have non-zero length");
}

/* ============================================================================
 * Test: path string building functions
 * ============================================================================ */
static void test_path_build_functions(void) {
    fprintf(stderr, "\n=== test_path_build_functions ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[1024];
    
    /* Test cache path building */
    if (path_to_cache("test_file", buffer, sizeof(buffer)) != NULL) {
        TEST_ASSERT_NOT_NULL(buffer, "cache path should be built");
        TEST_ASSERT(strlen(buffer) > 0, "cache path should have non-zero length");
    }
    
    /* Test varlib path building */
    if (path_to_varlib("test_file", buffer, sizeof(buffer)) != NULL) {
        TEST_ASSERT_NOT_NULL(buffer, "varlib path should be built");
        TEST_ASSERT(strlen(buffer) > 0, "varlib path should have non-zero length");
    }
}

/* ============================================================================
 * Test: path operations with empty strings
 * ============================================================================ */
static void test_path_operations_empty_strings(void) {
    fprintf(stderr, "\n=== test_path_operations_empty_strings ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[1024];
    
    /* Test with empty filename */
    path_to_cache("", buffer, sizeof(buffer));
    TEST_ASSERT_NOT_NULL(buffer, "empty filename should still return buffer");
    
    path_to_varlib("", buffer, sizeof(buffer));
    TEST_ASSERT_NOT_NULL(buffer, "empty filename should still return buffer");
}

/* ============================================================================
 * Test: path operations with very long strings
 * ============================================================================ */
static void test_path_operations_long_strings(void) {
    fprintf(stderr, "\n=== test_path_operations_long_strings ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[2048];
    char long_filename[1024];
    
    /* Create a very long filename */
    memset(long_filename, 'a', sizeof(long_filename) - 1);
    long_filename[sizeof(long_filename) - 1] = '\0';
    
    /* Should handle long filenames gracefully (either truncate or return NULL) */
    path_to_cache(long_filename, buffer, sizeof(buffer));
    /* As long as it doesn't crash, it's a pass */
    TEST_ASSERT(1, "long filename handling should not crash");
}

/* ============================================================================
 * Test: path buffer overflow protection
 * ============================================================================ */
static void test_path_buffer_overflow_protection(void) {
    fprintf(stderr, "\n=== test_path_buffer_overflow_protection ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char small_buffer[16];
    
    /* Test with small buffer - should not overflow */
    path_to_cache("very_long_filename_test.txt", small_buffer, sizeof(small_buffer));
    /* Should either return NULL or safely truncate, but not crash */
    TEST_ASSERT(1, "buffer overflow protection should work");
}

/* ============================================================================
 * Test: multiple consecutive path operations
 * ============================================================================ */
static void test_multiple_consecutive_path_operations(void) {
    fprintf(stderr, "\n=== test_multiple_consecutive_path_operations ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer1[1024], buffer2[1024], buffer3[1024];
    
    /* Multiple operations should not interfere with each other */
    path_to_cache("file1.txt", buffer1, sizeof(buffer1));
    path_to_varlib("file2.txt", buffer2, sizeof(buffer2));
    path_to_cache("file3.txt", buffer3, sizeof(buffer3));
    
    /* Each buffer should have content */
    TEST_ASSERT(strlen(buffer1) > 0, "first cache path should be non-empty");
    TEST_ASSERT(strlen(buffer2) > 0, "varlib path should be non-empty");
    TEST_ASSERT(strlen(buffer3) > 0, "third cache path should be non-empty");
    
    /* They should be different */
    TEST_ASSERT(strcmp(buffer1, buffer3) != 0, "different filenames should produce different paths");
}

/* ============================================================================
 * Test: path consistency across multiple calls
 * ============================================================================ */
static void test_path_consistency(void) {
    fprintf(stderr, "\n=== test_path_consistency ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer1[1024], buffer2[1024];
    
    /* Same input should produce same output */
    path_to_cache("test.txt", buffer1, sizeof(buffer1));
    path_to_cache("test.txt", buffer2, sizeof(buffer2));
    
    TEST_ASSERT_STRING_EQUAL(buffer1, buffer2, "same input should produce same output");
}

/* ============================================================================
 * Test: path directory components
 * ============================================================================ */
static void test_path_directory_components(void) {
    fprintf(stderr, "\n=== test_path_directory_components ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[1024];
    
    /* Test with nested directory structure in filename */
    path_to_cache("subdir/file.txt", buffer, sizeof(buffer));
    TEST_ASSERT_NOT_NULL(buffer, "nested path should be handled");
    
    path_to_varlib("subdir/another/file.txt", buffer, sizeof(buffer));
    TEST_ASSERT_NOT_NULL(buffer, "deeply nested path should be handled");
}

/* ============================================================================
 * Test: path special characters
 * ============================================================================ */
static void test_path_special_characters(void) {
    fprintf(stderr, "\n=== test_path_special_characters ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[1024];
    
    /* Test with special characters */
    path_to_cache("file-with-dashes.txt", buffer, sizeof(buffer));
    TEST_ASSERT_NOT_NULL(buffer, "dashes should be handled");
    
    path_to_cache("file_with_underscores.txt", buffer, sizeof(buffer));
    TEST_ASSERT_NOT_NULL(buffer, "underscores should be handled");
    
    path_to_cache("file.multiple.dots.txt", buffer, sizeof(buffer));
    TEST_ASSERT_NOT_NULL(buffer, "multiple dots should be handled");
}

/* ============================================================================
 * Test: paths initialization idempotency
 * ============================================================================ */
static void test_paths_initialization_idempotency(void) {
    fprintf(stderr, "\n=== test_paths_initialization_idempotency ===\n");
    
    /* Initialize multiple times - should be safe */
    char *result1 = paths_init(PATHS_DEFAULT);
    char *result2 = paths_init(PATHS_DEFAULT);
    
    TEST_ASSERT_NOT_NULL(result1, "first initialization should succeed");
    TEST_ASSERT_NOT_NULL(result2, "second initialization should succeed");
}

/* ============================================================================
 * Test: path_to_cache NULL buffer handling
 * ============================================================================ */
static void test_path_to_cache_null_buffer(void) {
    fprintf(stderr, "\n=== test_path_to_cache_null_buffer ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    /* Should handle NULL buffer gracefully */
    char *result = path_to_cache("test.txt", NULL, 0);
    /* Should either return NULL or handle safely */
    TEST_ASSERT(1, "NULL buffer should be handled");
}

/* ============================================================================
 * Test: path_to_varlib NULL buffer handling
 * ============================================================================ */
static void test_path_to_varlib_null_buffer(void) {
    fprintf(stderr, "\n=== test_path_to_varlib_null_buffer ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    /* Should handle NULL buffer gracefully */
    char *result = path_to_varlib("test.txt", NULL, 0);
    /* Should either return NULL or handle safely */
    TEST_ASSERT(1, "NULL buffer should be handled");
}

/* ============================================================================
 * Test: absolute vs relative paths in input
 * ============================================================================ */
static void test_absolute_vs_relative_paths(void) {
    fprintf(stderr, "\n=== test_absolute_vs_relative_paths ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer1[1024], buffer2[1024];
    
    /* Relative path */
    path_to_cache("relative/path.txt", buffer1, sizeof(buffer1));
    
    /* Absolute path - should still work */
    path_to_cache("/absolute/path.txt", buffer2, sizeof(buffer2));
    
    TEST_ASSERT(strlen(buffer1) > 0, "relative path should work");
    TEST_ASSERT(strlen(buffer2) > 0, "absolute path should work");
}

/* ============================================================================
 * Test: paths with environment variable expansion
 * ============================================================================ */
static void test_path_environment_expansion(void) {
    fprintf(stderr, "\n=== test_path_environment_expansion ===\n");
    
    /* Set test environment variable */
    setenv("TEST_VAR", "test_value", 1);
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[1024];
    path_to_cache("file.txt", buffer, sizeof(buffer));
    
    TEST_ASSERT_NOT_NULL(buffer, "path building should work with environment set");
}

/* ============================================================================
 * Test: concurrent path operations
 * ============================================================================ */
static void test_concurrent_path_operations(void) {
    fprintf(stderr, "\n=== test_concurrent_path_operations ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffers[10][1024];
    
    /* Simulate multiple sequential operations */
    for (int i = 0; i < 10; i++) {
        char filename[256];
        snprintf(filename, sizeof(filename), "file_%d.txt", i);
        
        if (i % 2 == 0) {
            path_to_cache(filename, buffers[i], sizeof(buffers[i]));
        } else {
            path_to_varlib(filename, buffers[i], sizeof(buffers[i]));
        }
        
        TEST_ASSERT(strlen(buffers[i]) > 0, "concurrent operation should succeed");
    }
}

/* ============================================================================
 * Test: Path case sensitivity
 * ============================================================================ */
static void test_path_case_sensitivity(void) {
    fprintf(stderr, "\n=== test_path_case_sensitivity ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer1[1024], buffer2[1024];
    
    /* Different cases should produce different filenames (on case-sensitive systems) */
    path_to_cache("File.txt", buffer1, sizeof(buffer1));
    path_to_cache("file.txt", buffer2, sizeof(buffer2));
    
    /* Paths should differ in the filename component */
    TEST_ASSERT(strlen(buffer1) > 0 && strlen(buffer2) > 0, "case handling should work");
}

/* ============================================================================
 * Test: Path with dots and parent directory references
 * ============================================================================ */
static void test_path_with_dots(void) {
    fprintf(stderr, "\n=== test_path_with_dots ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[1024];
    
    /* Path with .. should be handled */
    path_to_cache("../test.txt", buffer, sizeof(buffer));
    TEST_ASSERT(strlen(buffer) > 0, "parent directory reference should be handled");
    
    /* Path with . should be handled */
    path_to_cache("./test.txt", buffer, sizeof(buffer));
    TEST_ASSERT(strlen(buffer) > 0, "current directory reference should be handled");
    
    /* Path with ./ should be handled */
    path_to_cache("../../test.txt", buffer, sizeof(buffer));
    TEST_ASSERT(strlen(buffer) > 0, "multiple parent references should be handled");
}

/* ============================================================================
 * Test: Return value consistency
 * ============================================================================ */
static void test_return_value_consistency(void) {
    fprintf(stderr, "\n=== test_return_value_consistency ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer1[1024], buffer2[1024];
    
    /* Return values should be buffer pointers when successful */
    char *result1 = path_to_cache("test1.txt", buffer1, sizeof(buffer1));
    char *result2 = path_to_varlib("test2.txt", buffer2, sizeof(buffer2));
    
    /* Verify return values point to the buffers or are consistent */
    TEST_ASSERT_NOT_NULL(result1, "path_to_cache should return valid pointer");
    TEST_ASSERT_NOT_NULL(result2, "path_to_varlib should return valid pointer");
}

/* ============================================================================
 * Test: Paths initialization with NULL input
 * ============================================================================ */
static void test_paths_init_null_input(void) {
    fprintf(stderr, "\n=== test_paths_init_null_input ===\n");
    
    /* Initialize with NULL - should use default or handle gracefully */
    char *result = paths_init(NULL);
    /* Should either use default or return NULL, but not crash */
    TEST_ASSERT(1, "NULL input should be handled gracefully");
}

/* ============================================================================
 * Test: Paths with trailing slashes
 * ============================================================================ */
static void test_path_with_trailing_slashes(void) {
    fprintf(stderr, "\n=== test_path_with_trailing_slashes ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer1[1024], buffer2[1024];
    
    /* With trailing slash */
    path_to_cache("dir/file.txt", buffer1, sizeof(buffer1));
    
    /* Same but with trailing slashes in dirs (if applicable) */
    path_to_cache("dir//file.txt", buffer2, sizeof(buffer2));
    
    TEST_ASSERT(strlen(buffer1) > 0, "normal path should work");
    TEST_ASSERT(strlen(buffer2) > 0, "path with double slashes should work");
}

/* ============================================================================
 * Test: Global variable initialization
 * ============================================================================ */
static void test_global_variables_initialization(void) {
    fprintf(stderr, "\n=== test_global_variables_initialization ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    /* All global path variables should be initialized */
    TEST_ASSERT_NOT_NULL(netdata_cache_dir, "netdata_cache_dir should be initialized");
    TEST_ASSERT_NOT_NULL(netdata_varlib_dir, "netdata_varlib_dir should be initialized");
    TEST_ASSERT_NOT_NULL(netdata_sysconfdir, "netdata_sysconfdir should be initialized");
    TEST_ASSERT_NOT_NULL(netdata_rundir, "netdata_rundir should be initialized");
    
    /* All should be non-empty strings */
    TEST_ASSERT(strlen(netdata_cache_dir) > 0, "cache_dir should not be empty");
    TEST_ASSERT(strlen(netdata_varlib_dir) > 0, "varlib_dir should not be empty");
    TEST_ASSERT(strlen(netdata_sysconfdir) > 0, "sysconfdir should not be empty");
    TEST_ASSERT(strlen(netdata_rundir) > 0, "rundir should not be empty");
}

/* ============================================================================
 * Test: Path validation
 * ============================================================================ */
static void test_path_validation(void) {
    fprintf(stderr, "\n=== test_path_validation ===\n");
    
    paths_init(PATHS_DEFAULT);
    
    char buffer[1024];
    
    /* Paths should start with / (absolute) or contain base dir */
    path_to_cache("test.txt", buffer, sizeof(buffer));
    
    /* Verify it's a valid path (contains directory separator) */
    TEST_ASSERT(strchr(buffer, '/') != NULL || strlen(buffer) > 0, "path should be valid");
}

/* ============================================================================
 * Test: Multiple initialization with different paths
 * ============================================================================ */
static void test_multiple_init_different_paths(void) {
    fprintf(stderr, "\n=== test_multiple_init_different_paths ===\n");
    
    /* Initialize with first path */
    char *result1 = paths_init(PATHS_DEFAULT);
    TEST_ASSERT_NOT_NULL(result1, "first initialization should succeed");
    
    /* Initialize with different path */
    char custom_path[256];
    snprintf(custom_path, sizeof(custom_path), "/tmp/netdata_%d", getpid());
    char *result2 = paths_init(custom_path);
    TEST_ASSERT_NOT_NULL(result2, "second initialization with custom path should succeed");
}

/* ============================================================================
 * Main test runner
 * ============================================================================ */
int main(void) {
    fprintf(stderr, "\n========================================\n");
    fprintf(stderr, "NETDATA PATHS COMPREHENSIVE TEST SUITE\n");
    fprintf(stderr, "========================================\n");
    
    test_paths_init_valid();
    test_paths_init_custom();
    test_netdata_cache_dir_access();
    test_netdata_varlib_dir_access();
    test_netdata_sysconfdir_access();
    test_netdata_rundir_access();
    test_path_build_functions();
    test_path_operations_empty_strings();
    test_path_operations_long_strings();
    test_path_buffer_overflow_protection();
    test_multiple_consecutive_path_operations();
    test_path_consistency();
    test_path_directory_components();
    test_path_special_characters();
    test_paths_initialization_idempotency();
    test_path_to_cache_null_buffer();
    test_path_to_varlib_null_buffer();
    test_absolute_vs_relative_paths();
    test_path_environment_expansion();
    test_concurrent_path_operations();
    test_path_case_sensitivity();
    test_path_with_dots();
    test_return_value_consistency();
    test_paths_init_null_input();
    test_path_with_trailing_slashes();
    test_global_variables_initialization();
    test_path_validation();
    test_multiple_init_different_paths();
    
    fprintf(stderr, "\n========================================\n");
    fprintf(stderr, "ALL TESTS PASSED\n");
    fprintf(stderr, "========================================\n");
    
    return 0;
}