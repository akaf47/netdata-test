#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <windows.h>

/* Mock declarations for WinAPI functions */
ULONGLONG mock_GetTickCount64(void);
BOOL mock_GetSystemTimes(LPFILETIME lpIdleTime, LPFILETIME lpKernelTime, LPFILETIME lpUserTime);
BOOL mock_GetNumberOfProcessors(DWORD *lpNumberOfProcessors);

/* Assume the real function signature */
typedef struct {
    ULONGLONG idle_time;
    ULONGLONG kernel_time;
    ULONGLONG user_time;
    DWORD processor_count;
    int status;
} CPU_INFO;

/* Mock implementation of GetSystemCPU */
static CPU_INFO* GetSystemCPU(void);

/* ===== TEST FIXTURES ===== */
static int setup(void **state) {
    return 0;
}

static int teardown(void **state) {
    return 0;
}

/* ===== TEST CASES ===== */

/* Test: Successful CPU information retrieval */
static void test_GetSystemCPU_success(void **state) {
    FILETIME idle_time = {0x00000000, 0x00000001};
    FILETIME kernel_time = {0x00000000, 0x00000002};
    FILETIME user_time = {0x00000000, 0x00000003};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_equal(result->status, 0);
}

/* Test: GetSystemTimes fails - should return error */
static void test_GetSystemCPU_GetSystemTimes_fails(void **state) {
    will_return(mock_GetSystemTimes, FALSE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_not_equal(result->status, 0);
}

/* Test: Zero idle time - minimum boundary */
static void test_GetSystemCPU_zero_idle_time(void **state) {
    FILETIME idle_time = {0, 0};
    FILETIME kernel_time = {0, 1};
    FILETIME user_time = {0, 1};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_equal(result->idle_time, 0);
    assert_int_equal(result->status, 0);
}

/* Test: Maximum FILETIME values - boundary test */
static void test_GetSystemCPU_max_filetime_values(void **state) {
    FILETIME idle_time = {0xFFFFFFFF, 0xFFFFFFFF};
    FILETIME kernel_time = {0xFFFFFFFF, 0xFFFFFFFF};
    FILETIME user_time = {0xFFFFFFFF, 0xFFFFFFFF};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_equal(result->status, 0);
}

/* Test: Memory allocation failure */
static void test_GetSystemCPU_allocation_failure(void **state) {
    expect_value(__wrap_malloc, size, sizeof(CPU_INFO));
    will_return(__wrap_malloc, NULL);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_null(result);
}

/* Test: Null pointer handling */
static void test_GetSystemCPU_null_filetime_structures(void **state) {
    will_return(mock_GetSystemTimes, FALSE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_not_equal(result->status, 0);
}

/* Test: Multiple successive calls */
static void test_GetSystemCPU_multiple_calls(void **state) {
    FILETIME idle_time = {0, 1};
    FILETIME kernel_time = {0, 2};
    FILETIME user_time = {0, 3};
    
    for (int i = 0; i < 5; i++) {
        will_return(mock_GetSystemTimes, &idle_time);
        will_return(mock_GetSystemTimes, &kernel_time);
        will_return(mock_GetSystemTimes, &user_time);
        will_return(mock_GetSystemTimes, TRUE);
        
        CPU_INFO *result = GetSystemCPU();
        assert_non_null(result);
        assert_int_equal(result->status, 0);
        free(result);
    }
}

/* Test: Idle time greater than kernel time */
static void test_GetSystemCPU_idle_greater_than_kernel(void **state) {
    FILETIME idle_time = {0, 100};
    FILETIME kernel_time = {0, 50};
    FILETIME user_time = {0, 30};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_equal(result->status, 0);
}

/* Test: Kernel time equals user time */
static void test_GetSystemCPU_kernel_equals_user(void **state) {
    FILETIME idle_time = {0, 50};
    FILETIME kernel_time = {0, 25};
    FILETIME user_time = {0, 25};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_equal(result->status, 0);
}

/* Test: All times equal */
static void test_GetSystemCPU_all_times_equal(void **state) {
    FILETIME times = {0, 100};
    
    will_return(mock_GetSystemTimes, &times);
    will_return(mock_GetSystemTimes, &times);
    will_return(mock_GetSystemTimes, &times);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_equal(result->status, 0);
}

/* Test: FILETIME low/high parts boundary values */
static void test_GetSystemCPU_filetime_high_low_parts(void **state) {
    FILETIME idle_time = {1, 0};
    FILETIME kernel_time = {0, 1};
    FILETIME user_time = {0xFFFFFFFF, 0};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_int_equal(result->status, 0);
}

/* Test: Result structure is properly populated */
static void test_GetSystemCPU_result_structure_populated(void **state) {
    FILETIME idle_time = {0x11111111, 0x22222222};
    FILETIME kernel_time = {0x33333333, 0x44444444};
    FILETIME user_time = {0x55555555, 0x66666666};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    assert_non_zero(result->idle_time);
    assert_non_zero(result->kernel_time);
    assert_non_zero(result->user_time);
}

/* Test: Return value is non-null on success */
static void test_GetSystemCPU_returns_non_null(void **state) {
    FILETIME idle_time = {0, 1};
    FILETIME kernel_time = {0, 1};
    FILETIME user_time = {0, 1};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_non_null(result);
    free(result);
}

/* Test: Status field indicates success (0) on valid call */
static void test_GetSystemCPU_status_success(void **state) {
    FILETIME idle_time = {0, 1};
    FILETIME kernel_time = {0, 1};
    FILETIME user_time = {0, 1};
    
    will_return(mock_GetSystemTimes, &idle_time);
    will_return(mock_GetSystemTimes, &kernel_time);
    will_return(mock_GetSystemTimes, &user_time);
    will_return(mock_GetSystemTimes, TRUE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_int_equal(result->status, 0);
    free(result);
}

/* Test: Status field indicates failure on GetSystemTimes error */
static void test_GetSystemCPU_status_failure(void **state) {
    will_return(mock_GetSystemTimes, FALSE);
    
    CPU_INFO *result = GetSystemCPU();
    
    assert_true(result->status != 0);
    free(result);
}

/* ===== MAIN TEST RUNNER ===== */
int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_GetSystemTimes_fails, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_zero_idle_time, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_max_filetime_values, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_allocation_failure, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_null_filetime_structures, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_multiple_calls, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_idle_greater_than_kernel, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_kernel_equals_user, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_all_times_equal, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_filetime_high_low_parts, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_result_structure_populated, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_returns_non_null, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_status_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_GetSystemCPU_status_failure, setup, teardown),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}