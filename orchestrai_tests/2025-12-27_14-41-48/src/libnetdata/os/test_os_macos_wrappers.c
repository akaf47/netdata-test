// SPDX-License-Identifier: GPL-3.0-or-later
// Comprehensive tests for os-macos-wrappers.h and os-macos-wrappers.c

#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <errno.h>
#include <sys/sysctl.h>
#include <cmocka.h>
#include <setjmp.h>

// Mock setup for testing
typedef struct {
    jmp_buf jump_buffer;
    int errno_value;
    int sysctlbyname_return;
    size_t sysctlbyname_nlen;
    bool should_fail;
} test_context_t;

// Forward declarations
int getsysctl_by_name(const char *name, void *ptr, size_t len);

// Mock functions
int __wrap_sysctlbyname(const char *name, void *oldp, size_t *oldlenp, 
                       const void *newp, size_t newlen) {
    check_expected_ptr(name);
    check_expected_ptr(oldp);
    check_expected_ptr(oldlenp);
    
    if (oldlenp != NULL) {
        *oldlenp = mock_type(size_t);
    }
    
    return mock_type(int);
}

void __wrap_netdata_log_error(const char *fmt, ...) {
    // Mock logging - do nothing in tests
}

// Test fixtures
static int setup_test(void **state) {
    test_context_t *ctx = malloc(sizeof(test_context_t));
    if (ctx == NULL) return -1;
    memset(ctx, 0, sizeof(test_context_t));
    *state = ctx;
    return 0;
}

static int teardown_test(void **state) {
    test_context_t *ctx = (test_context_t *)*state;
    if (ctx) free(ctx);
    return 0;
}

// Test: getsysctl_by_name - successful call with matching length
static void test_getsysctl_by_name_success_int32(void **state) {
    int32_t test_value = 0;
    const char *sysctl_name = "hw.ncpu";
    size_t expected_len = sizeof(int32_t);
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);  // nlen returned
    will_return(__wrap_sysctlbyname, 0);              // success
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 0);
}

// Test: getsysctl_by_name - successful call with uint32_t
static void test_getsysctl_by_name_success_uint32(void **state) {
    uint32_t test_value = 0;
    const char *sysctl_name = "hw.memsize";
    size_t expected_len = sizeof(uint32_t);
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);
    will_return(__wrap_sysctlbyname, 0);
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 0);
}

// Test: getsysctl_by_name - successful call with uint64_t
static void test_getsysctl_by_name_success_uint64(void **state) {
    uint64_t test_value = 0;
    const char *sysctl_name = "hw.memsize";
    size_t expected_len = sizeof(uint64_t);
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);
    will_return(__wrap_sysctlbyname, 0);
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 0);
}

// Test: getsysctl_by_name - sysctlbyname returns -1 (error)
static void test_getsysctl_by_name_sysctlbyname_error(void **state) {
    int32_t test_value = 0;
    const char *sysctl_name = "hw.nonexistent";
    size_t expected_len = sizeof(int32_t);
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);
    will_return(__wrap_sysctlbyname, -1);  // Error return
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 1);  // Should return 1 on error
}

// Test: getsysctl_by_name - length mismatch (too small)
static void test_getsysctl_by_name_length_mismatch_small(void **state) {
    int32_t test_value = 0;
    const char *sysctl_name = "hw.ncpu";
    size_t expected_len = sizeof(int32_t);
    size_t returned_len = expected_len - 1;  // Mismatch: too small
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, returned_len);
    will_return(__wrap_sysctlbyname, 0);  // sysctlbyname succeeds
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 1);  // Should return 1 on length mismatch
}

// Test: getsysctl_by_name - length mismatch (too large)
static void test_getsysctl_by_name_length_mismatch_large(void **state) {
    int32_t test_value = 0;
    const char *sysctl_name = "hw.ncpu";
    size_t expected_len = sizeof(int32_t);
    size_t returned_len = expected_len + 1;  // Mismatch: too large
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, returned_len);
    will_return(__wrap_sysctlbyname, 0);  // sysctlbyname succeeds
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 1);  // Should return 1 on length mismatch
}

// Test: getsysctl_by_name - with empty string name (edge case)
static void test_getsysctl_by_name_empty_name(void **state) {
    int32_t test_value = 0;
    const char *sysctl_name = "";
    size_t expected_len = sizeof(int32_t);
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);
    will_return(__wrap_sysctlbyname, -1);  // Empty name should fail
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 1);
}

// Test: getsysctl_by_name - with very long sysctl name
static void test_getsysctl_by_name_long_name(void **state) {
    int32_t test_value = 0;
    const char *sysctl_name = "very.long.sysctl.name.that.tests.boundary.conditions.in.macos";
    size_t expected_len = sizeof(int32_t);
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);
    will_return(__wrap_sysctlbyname, 0);
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 0);
}

// Test: getsysctl_by_name - with size 0 (edge case)
static void test_getsysctl_by_name_zero_size(void **state) {
    int32_t test_value = 0;
    const char *sysctl_name = "hw.ncpu";
    size_t expected_len = 0;
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);
    will_return(__wrap_sysctlbyname, 0);
    
    int result = getsysctl_by_name(sysctl_name, &test_value, expected_len);
    
    assert_int_equal(result, 0);
}

// Test: getsysctl_by_name - with large buffer size
static void test_getsysctl_by_name_large_size(void **state) {
    char test_buffer[4096];
    const char *sysctl_name = "kern.osversion";
    size_t expected_len = sizeof(test_buffer);
    
    expect_string(__wrap_sysctlbyname, name, sysctl_name);
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, expected_len);
    will_return(__wrap_sysctlbyname, 0);
    
    int result = getsysctl_by_name(sysctl_name, test_buffer, expected_len);
    
    assert_int_equal(result, 0);
}

// Test: GETSYSCTL_BY_NAME macro - successful execution
static void test_macro_getsysctl_by_name_success(void **state) {
    int32_t ncpu = 0;
    
    expect_string(__wrap_sysctlbyname, name, "hw.ncpu");
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, sizeof(int32_t));
    will_return(__wrap_sysctlbyname, 0);
    
    int result = GETSYSCTL_BY_NAME("hw.ncpu", ncpu);
    
    assert_int_equal(result, 0);
}

// Test: GETSYSCTL_BY_NAME macro - failure
static void test_macro_getsysctl_by_name_failure(void **state) {
    int32_t ncpu = 0;
    
    expect_string(__wrap_sysctlbyname, name, "hw.invalid");
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, sizeof(int32_t));
    will_return(__wrap_sysctlbyname, -1);
    
    int result = GETSYSCTL_BY_NAME("hw.invalid", ncpu);
    
    assert_int_equal(result, 1);
}

// Test: getsysctl_by_name - multiple sequential calls
static void test_getsysctl_by_name_sequential_calls(void **state) {
    int32_t ncpu = 0;
    uint64_t memsize = 0;
    size_t len1 = sizeof(int32_t);
    size_t len2 = sizeof(uint64_t);
    
    // First call
    expect_string(__wrap_sysctlbyname, name, "hw.ncpu");
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, len1);
    will_return(__wrap_sysctlbyname, 0);
    
    int result1 = getsysctl_by_name("hw.ncpu", &ncpu, len1);
    assert_int_equal(result1, 0);
    
    // Second call
    expect_string(__wrap_sysctlbyname, name, "hw.memsize");
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, len2);
    will_return(__wrap_sysctlbyname, 0);
    
    int result2 = getsysctl_by_name("hw.memsize", &memsize, len2);
    assert_int_equal(result2, 0);
}

// Test: getsysctl_by_name - alternating success and failure
static void test_getsysctl_by_name_alternating_success_failure(void **state) {
    int32_t val1 = 0;
    int32_t val2 = 0;
    size_t len = sizeof(int32_t);
    
    // Success
    expect_string(__wrap_sysctlbyname, name, "hw.ncpu");
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, len);
    will_return(__wrap_sysctlbyname, 0);
    
    int result1 = getsysctl_by_name("hw.ncpu", &val1, len);
    assert_int_equal(result1, 0);
    
    // Failure
    expect_string(__wrap_sysctlbyname, name, "hw.invalid");
    expect_any(__wrap_sysctlbyname, oldp);
    expect_any(__wrap_sysctlbyname, oldlenp);
    will_return(__wrap_sysctlbyname, len);
    will_return(__wrap_sysctlbyname, -1);
    
    int result2 = getsysctl_by_name("hw.invalid", &val2, len);
    assert_int_equal(result2, 1);
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_success_int32, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_success_uint32, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_success_uint64, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_sysctlbyname_error, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_length_mismatch_small, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_length_mismatch_large, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_empty_name, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_long_name, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_zero_size, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_large_size, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_macro_getsysctl_by_name_success, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_macro_getsysctl_by_name_failure, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_sequential_calls, setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_getsysctl_by_name_alternating_success_failure, setup_test, teardown_test),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}