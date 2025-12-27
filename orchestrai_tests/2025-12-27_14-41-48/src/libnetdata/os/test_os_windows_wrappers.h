#ifndef TEST_OS_WINDOWS_WRAPPERS_H
#define TEST_OS_WINDOWS_WRAPPERS_H

#ifdef _WIN32

#include <assert.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <windows.h>

/* Test Suite for os-windows-wrappers.h */

/* Mock structures and functions for testing */
typedef struct {
    int call_count;
    int return_value;
    void *return_ptr;
} mock_context;

/* Test: posix_memalign wrapper - success case */
void test_posix_memalign_success(void) {
    void *ptr = NULL;
    int result = posix_memalign(&ptr, 16, 1024);
    assert(result == 0);
    assert(ptr != NULL);
    assert(((uintptr_t)ptr % 16) == 0);
    free(ptr);
    printf("✓ test_posix_memalign_success\n");
}

/* Test: posix_memalign wrapper - invalid alignment */
void test_posix_memalign_invalid_alignment(void) {
    void *ptr = NULL;
    int result = posix_memalign(&ptr, 3, 1024);
    assert(result == EINVAL);
    printf("✓ test_posix_memalign_invalid_alignment\n");
}

/* Test: posix_memalign wrapper - zero size */
void test_posix_memalign_zero_size(void) {
    void *ptr = NULL;
    int result = posix_memalign(&ptr, 16, 0);
    assert(result == 0);
    if (ptr) free(ptr);
    printf("✓ test_posix_memalign_zero_size\n");
}

/* Test: posix_memalign wrapper - large size */
void test_posix_memalign_large_size(void) {
    void *ptr = NULL;
    int result = posix_memalign(&ptr, 16, 1024 * 1024);
    if (result == 0) {
        assert(ptr != NULL);
        assert(((uintptr_t)ptr % 16) == 0);
        free(ptr);
    }
    printf("✓ test_posix_memalign_large_size\n");
}

/* Test: posix_memalign wrapper - null pointer */
void test_posix_memalign_null_pointer(void) {
    int result = posix_memalign(NULL, 16, 1024);
    assert(result != 0);
    printf("✓ test_posix_memalign_null_pointer\n");
}

/* Test: clock_gettime wrapper - CLOCK_REALTIME */
void test_clock_gettime_realtime(void) {
    struct timespec ts;
    int result = clock_gettime(CLOCK_REALTIME, &ts);
    assert(result == 0);
    assert(ts.tv_sec > 0);
    assert(ts.tv_nsec >= 0 && ts.tv_nsec < 1000000000);
    printf("✓ test_clock_gettime_realtime\n");
}

/* Test: clock_gettime wrapper - CLOCK_MONOTONIC */
void test_clock_gettime_monotonic(void) {
    struct timespec ts;
    int result = clock_gettime(CLOCK_MONOTONIC, &ts);
    assert(result == 0);
    assert(ts.tv_sec >= 0);
    assert(ts.tv_nsec >= 0 && ts.tv_nsec < 1000000000);
    printf("✓ test_clock_gettime_monotonic\n");
}

/* Test: clock_gettime wrapper - CLOCK_PROCESS_CPUTIME_ID */
void test_clock_gettime_process_cputime(void) {
    struct timespec ts;
    int result = clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &ts);
    assert(result == 0);
    assert(ts.tv_sec >= 0);
    assert(ts.tv_nsec >= 0 && ts.tv_nsec < 1000000000);
    printf("✓ test_clock_gettime_process_cputime\n");
}

/* Test: clock_gettime wrapper - CLOCK_THREAD_CPUTIME_ID */
void test_clock_gettime_thread_cputime(void) {
    struct timespec ts;
    int result = clock_gettime(CLOCK_THREAD_CPUTIME_ID, &ts);
    assert(result == 0);
    assert(ts.tv_sec >= 0);
    assert(ts.tv_nsec >= 0 && ts.tv_nsec < 1000000000);
    printf("✓ test_clock_gettime_thread_cputime\n");
}

/* Test: clock_gettime wrapper - null timespec */
void test_clock_gettime_null_timespec(void) {
    int result = clock_gettime(CLOCK_REALTIME, NULL);
    assert(result != 0);
    printf("✓ test_clock_gettime_null_timespec\n");
}

/* Test: clock_gettime wrapper - invalid clock_id */
void test_clock_gettime_invalid_clock(void) {
    struct timespec ts;
    int result = clock_gettime(-1, &ts);
    assert(result != 0);
    printf("✓ test_clock_gettime_invalid_clock\n");
}

/* Test: getppid wrapper */
void test_getppid(void) {
    pid_t ppid = getppid();
    assert(ppid > 0);
    printf("✓ test_getppid\n");
}

/* Test: gettid wrapper */
void test_gettid(void) {
    unsigned long tid = gettid();
    assert(tid > 0);
    printf("✓ test_gettid\n");
}

/* Test: strerror_r wrapper - success case */
void test_strerror_r_success(void) {
    char buf[256];
    int result = strerror_r(EINVAL, buf, sizeof(buf));
    assert(result == 0);
    assert(strlen(buf) > 0);
    printf("✓ test_strerror_r_success\n");
}

/* Test: strerror_r wrapper - small buffer */
void test_strerror_r_small_buffer(void) {
    char buf[2];
    int result = strerror_r(EINVAL, buf, sizeof(buf));
    assert(result != 0);
    printf("✓ test_strerror_r_small_buffer\n");
}

/* Test: strerror_r wrapper - null buffer */
void test_strerror_r_null_buffer(void) {
    int result = strerror_r(EINVAL, NULL, 256);
    assert(result != 0);
    printf("✓ test_strerror_r_null_buffer\n");
}

/* Test: strerror_r wrapper - invalid errno */
void test_strerror_r_invalid_errno(void) {
    char buf[256];
    int result = strerror_r(-9999, buf, sizeof(buf));
    assert(result == 0);
    assert(strlen(buf) > 0);
    printf("✓ test_strerror_r_invalid_errno\n");
}

/* Test: strerror_r wrapper - zero buffer size */
void test_strerror_r_zero_buffer_size(void) {
    char buf[256];
    int result = strerror_r(EINVAL, buf, 0);
    assert(result != 0);
    printf("✓ test_strerror_r_zero_buffer_size\n");
}

/* Test: mkdir wrapper - new directory */
void test_mkdir_new_directory(void) {
    const char *test_dir = "test_mkdir_dir";
    int result = mkdir(test_dir, 0755);
    assert(result == 0);
    /* Cleanup */
    RemoveDirectory(test_dir);
    printf("✓ test_mkdir_new_directory\n");
}

/* Test: mkdir wrapper - existing directory */
void test_mkdir_existing_directory(void) {
    const char *test_dir = "test_mkdir_dir";
    mkdir(test_dir, 0755);
    int result = mkdir(test_dir, 0755);
    assert(result != 0);
    /* Cleanup */
    RemoveDirectory(test_dir);
    printf("✓ test_mkdir_existing_directory\n");
}

/* Test: mkdir wrapper - null path */
void test_mkdir_null_path(void) {
    int result = mkdir(NULL, 0755);
    assert(result != 0);
    printf("✓ test_mkdir_null_path\n");
}

/* Test: mkdir wrapper - empty path */
void test_mkdir_empty_path(void) {
    int result = mkdir("", 0755);
    assert(result != 0);
    printf("✓ test_mkdir_empty_path\n");
}

/* Test: mkdir wrapper - invalid mode */
void test_mkdir_invalid_mode(void) {
    const char *test_dir = "test_mkdir_mode";
    int result = mkdir(test_dir, 0);
    if (result == 0) {
        RemoveDirectory(test_dir);
    }
    printf("✓ test_mkdir_invalid_mode\n");
}

/* Test runner */
void run_all_windows_wrapper_tests(void) {
    printf("\n=== Windows Wrapper Tests ===\n");
    test_posix_memalign_success();
    test_posix_memalign_invalid_alignment();
    test_posix_memalign_zero_size();
    test_posix_memalign_large_size();
    test_posix_memalign_null_pointer();
    test_clock_gettime_realtime();
    test_clock_gettime_monotonic();
    test_clock_gettime_process_cputime();
    test_clock_gettime_thread_cputime();
    test_clock_gettime_null_timespec();
    test_clock_gettime_invalid_clock();
    test_getppid();
    test_gettid();
    test_strerror_r_success();
    test_strerror_r_small_buffer();
    test_strerror_r_null_buffer();
    test_strerror_r_invalid_errno();
    test_strerror_r_zero_buffer_size();
    test_mkdir_new_directory();
    test_mkdir_existing_directory();
    test_mkdir_null_path();
    test_mkdir_empty_path();
    test_mkdir_invalid_mode();
    printf("\n=== All Windows Wrapper Tests Completed ===\n");
}

#endif /* _WIN32 */
#endif /* TEST_OS_WINDOWS_WRAPPERS_H */