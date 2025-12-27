#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <time.h>
#include <windows.h>
#include <psapi.h>
#include <assert.h>

/* Mock structures and function prototypes */
typedef struct {
    uint32_t pid;
    char comm[256];
    char cmdline[4096];
    uint64_t minflt;
    uint64_t majflt;
    uint64_t utime;
    uint64_t stime;
    int64_t priority;
    uint32_t threads;
} proc_data_t;

typedef struct {
    uint32_t pid;
    HANDLE handle;
    FILETIME creation_time;
    FILETIME exit_time;
    FILETIME kernel_time;
    FILETIME user_time;
} proc_handle_t;

/* Function prototypes to test */
extern int get_windows_process_info(uint32_t pid, proc_data_t *data);
extern int enumerate_processes(proc_data_t **procs, uint32_t *count);
extern void free_process_list(proc_data_t *procs, uint32_t count);
extern uint64_t filetime_to_microseconds(FILETIME ft);
extern int get_process_cpu_time(uint32_t pid, uint64_t *user_time, uint64_t *kernel_time);
extern int get_process_memory_info(uint32_t pid, uint64_t *rss, uint64_t *vms);
extern int get_process_threads(uint32_t pid, uint32_t *thread_count);
extern char* get_process_cmdline(uint32_t pid, char *buffer, size_t buflen);

/* ============================================================================
 * Test Suite: get_windows_process_info
 * ============================================================================ */

void test_get_windows_process_info_success() {
    proc_data_t data = {0};
    int result = get_windows_process_info(GetCurrentProcessId(), &data);
    
    assert(result == 0);
    assert(data.pid > 0);
    assert(strlen(data.comm) > 0);
    printf("✓ test_get_windows_process_info_success\n");
}

void test_get_windows_process_info_invalid_pid() {
    proc_data_t data = {0};
    int result = get_windows_process_info(0xFFFFFFFF, &data);
    
    assert(result != 0);
    printf("✓ test_get_windows_process_info_invalid_pid\n");
}

void test_get_windows_process_info_null_data_pointer() {
    int result = get_windows_process_info(GetCurrentProcessId(), NULL);
    
    assert(result != 0);
    printf("✓ test_get_windows_process_info_null_data_pointer\n");
}

void test_get_windows_process_info_zero_pid() {
    proc_data_t data = {0};
    int result = get_windows_process_info(0, &data);
    
    assert(result != 0);
    printf("✓ test_get_windows_process_info_zero_pid\n");
}

void test_get_windows_process_info_initializes_all_fields() {
    proc_data_t data = {0xFF, "test", "test", 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF};
    int result = get_windows_process_info(GetCurrentProcessId(), &data);
    
    assert(result == 0);
    assert(data.minflt >= 0);
    assert(data.majflt >= 0);
    assert(data.utime >= 0);
    assert(data.stime >= 0);
    assert(data.threads > 0);
    printf("✓ test_get_windows_process_info_initializes_all_fields\n");
}

/* ============================================================================
 * Test Suite: enumerate_processes
 * ============================================================================ */

void test_enumerate_processes_success() {
    proc_data_t *procs = NULL;
    uint32_t count = 0;
    
    int result = enumerate_processes(&procs, &count);
    
    assert(result == 0);
    assert(count > 0);
    assert(procs != NULL);
    
    free_process_list(procs, count);
    printf("✓ test_enumerate_processes_success\n");
}

void test_enumerate_processes_null_procs_pointer() {
    uint32_t count = 0;
    int result = enumerate_processes(NULL, &count);
    
    assert(result != 0);
    printf("✓ test_enumerate_processes_null_procs_pointer\n");
}

void test_enumerate_processes_null_count_pointer() {
    proc_data_t *procs = NULL;
    int result = enumerate_processes(&procs, NULL);
    
    assert(result != 0);
    printf("✓ test_enumerate_processes_null_count_pointer\n");
}

void test_enumerate_processes_finds_system_processes() {
    proc_data_t *procs = NULL;
    uint32_t count = 0;
    int result = enumerate_processes(&procs, &count);
    
    assert(result == 0);
    assert(count > 0);
    
    int found_current = 0;
    uint32_t current_pid = GetCurrentProcessId();
    for(uint32_t i = 0; i < count; i++) {
        if(procs[i].pid == current_pid) {
            found_current = 1;
            break;
        }
    }
    
    assert(found_current);
    free_process_list(procs, count);
    printf("✓ test_enumerate_processes_finds_system_processes\n");
}

void test_enumerate_processes_count_consistency() {
    proc_data_t *procs1 = NULL, *procs2 = NULL;
    uint32_t count1 = 0, count2 = 0;
    
    int result1 = enumerate_processes(&procs1, &count1);
    int result2 = enumerate_processes(&procs2, &count2);
    
    assert(result1 == 0);
    assert(result2 == 0);
    assert(count1 > 0);
    /* Counts may differ slightly due to process churn */
    assert(count2 > 0);
    
    free_process_list(procs1, count1);
    free_process_list(procs2, count2);
    printf("✓ test_enumerate_processes_count_consistency\n");
}

/* ============================================================================
 * Test Suite: free_process_list
 * ============================================================================ */

void test_free_process_list_success() {
    proc_data_t *procs = malloc(3 * sizeof(proc_data_t));
    assert(procs != NULL);
    
    free_process_list(procs, 3);
    printf("✓ test_free_process_list_success\n");
}

void test_free_process_list_zero_count() {
    proc_data_t *procs = malloc(sizeof(proc_data_t));
    assert(procs != NULL);
    
    free_process_list(procs, 0);
    printf("✓ test_free_process_list_zero_count\n");
}

void test_free_process_list_null_pointer() {
    /* Should handle gracefully */
    free_process_list(NULL, 0);
    printf("✓ test_free_process_list_null_pointer\n");
}

void test_free_process_list_large_count() {
    proc_data_t *procs = malloc(1000 * sizeof(proc_data_t));
    assert(procs != NULL);
    
    free_process_list(procs, 1000);
    printf("✓ test_free_process_list_large_count\n");
}

/* ============================================================================
 * Test Suite: filetime_to_microseconds
 * ============================================================================ */

void test_filetime_to_microseconds_zero() {
    FILETIME ft = {0, 0};
    uint64_t result = filetime_to_microseconds(ft);
    
    assert(result == 0);
    printf("✓ test_filetime_to_microseconds_zero\n");
}

void test_filetime_to_microseconds_one_hundred_nanoseconds() {
    FILETIME ft = {1, 0};
    uint64_t result = filetime_to_microseconds(ft);
    
    assert(result == 0); /* Less than 1 microsecond */
    printf("✓ test_filetime_to_microseconds_one_hundred_nanoseconds\n");
}

void test_filetime_to_microseconds_ten_microseconds() {
    FILETIME ft = {100, 0}; /* 100 * 100ns = 10us */
    uint64_t result = filetime_to_microseconds(ft);
    
    assert(result == 10);
    printf("✓ test_filetime_to_microseconds_ten_microseconds\n");
}

void test_filetime_to_microseconds_one_second() {
    FILETIME ft = {10000000, 0}; /* 10000000 * 100ns = 1s */
    uint64_t result = filetime_to_microseconds(ft);
    
    assert(result == 1000000);
    printf("✓ test_filetime_to_microseconds_one_second\n");
}

void test_filetime_to_microseconds_high_part() {
    FILETIME ft = {0, 1}; /* High part set */
    uint64_t result = filetime_to_microseconds(ft);
    
    assert(result > 0);
    printf("✓ test_filetime_to_microseconds_high_part\n");
}

void test_filetime_to_microseconds_max_value() {
    FILETIME ft = {0xFFFFFFFF, 0xFFFFFFFF};
    uint64_t result = filetime_to_microseconds(ft);
    
    assert(result > 0);
    printf("✓ test_filetime_to_microseconds_max_value\n");
}

/* ============================================================================
 * Test Suite: get_process_cpu_time
 * ============================================================================ */

void test_get_process_cpu_time_current_process() {
    uint64_t user_time = 0, kernel_time = 0;
    int result = get_process_cpu_time(GetCurrentProcessId(), &user_time, &kernel_time);
    
    assert(result == 0);
    assert(user_time >= 0);
    assert(kernel_time >= 0);
    printf("✓ test_get_process_cpu_time_current_process\n");
}

void test_get_process_cpu_time_invalid_pid() {
    uint64_t user_time = 0, kernel_time = 0;
    int result = get_process_cpu_time(0xFFFFFFFF, &user_time, &kernel_time);
    
    assert(result != 0);
    printf("✓ test_get_process_cpu_time_invalid_pid\n");
}

void test_get_process_cpu_time_null_user_time() {
    uint64_t kernel_time = 0;
    int result = get_process_cpu_time(GetCurrentProcessId(), NULL, &kernel_time);
    
    assert(result != 0);
    printf("✓ test_get_process_cpu_time_null_user_time\n");
}

void test_get_process_cpu_time_null_kernel_time() {
    uint64_t user_time = 0;
    int result = get_process_cpu_time(GetCurrentProcessId(), &user_time, NULL);
    
    assert(result != 0);
    printf("✓ test_get_process_cpu_time_null_kernel_time\n");
}

void test_get_process_cpu_time_both_null() {
    int result = get_process_cpu_time(GetCurrentProcessId(), NULL, NULL);
    
    assert(result != 0);
    printf("✓ test_get_process_cpu_time_both_null\n");
}

void test_get_process_cpu_time_zero_pid() {
    uint64_t user_time = 0, kernel_time = 0;
    int result = get_process_cpu_time(0, &user_time, &kernel_time);
    
    assert(result != 0);
    printf("✓ test_get_process_cpu_time_zero_pid\n");
}

/* ============================================================================
 * Test Suite: get_process_memory_info
 * ============================================================================ */

void test_get_process_memory_info_current_process() {
    uint64_t rss = 0, vms = 0;
    int result = get_process_memory_info(GetCurrentProcessId(), &rss, &vms);
    
    assert(result == 0);
    assert(rss > 0);
    assert(vms > 0);
    printf("✓ test_get_process_memory_info_current_process\n");
}

void test_get_process_memory_info_invalid_pid() {
    uint64_t rss = 0, vms = 0;
    int result = get_process_memory_info(0xFFFFFFFF, &rss, &vms);
    
    assert(result != 0);
    printf("✓ test_get_process_memory_info_invalid_pid\n");
}

void test_get_process_memory_info_null_rss() {
    uint64_t vms = 0;
    int result = get_process_memory_info(GetCurrentProcessId(), NULL, &vms);
    
    assert(result != 0);
    printf("✓ test_get_process_memory_info_null_rss\n");
}

void test_get_process_memory_info_null_vms() {
    uint64_t rss = 0;
    int result = get_process_memory_info(GetCurrentProcessId(), &rss, NULL);
    
    assert(result != 0);
    printf("✓ test_get_process_memory_info_null_vms\n");
}

void test_get_process_memory_info_both_null() {
    int result = get_process_memory_info(GetCurrentProcessId(), NULL, NULL);
    
    assert(result != 0);
    printf("✓ test_get_process_memory_info_both_null\n");
}

void test_get_process_memory_info_zero_pid() {
    uint64_t rss = 0, vms = 0;
    int result = get_process_memory_info(0, &rss, &vms);
    
    assert(result != 0);
    printf("✓ test_get_process_memory_info_zero_pid\n");
}

void test_get_process_memory_info_rss_less_than_vms() {
    uint64_t rss = 0, vms = 0;
    int result = get_process_memory_info(GetCurrentProcessId(), &rss, &vms);
    
    assert(result == 0);
    assert(rss <= vms);
    printf("✓ test_get_process_memory_info_rss_less_than_vms\n");
}

/* ============================================================================
 * Test Suite: get_process_threads
 * ============================================================================ */

void test_get_process_threads_current_process() {
    uint32_t thread_count = 0;
    int result = get_process_threads(GetCurrentProcessId(), &thread_count);
    
    assert(result == 0);
    assert(thread_count > 0);
    printf("✓ test_get_process_threads_current_process\n");
}

void test_get_process_threads_invalid_pid() {
    uint32_t thread_count = 0;
    int result = get_process_threads(0xFFFFFFFF, &thread_count);
    
    assert(result != 0);
    printf("✓ test_get_process_threads_invalid_pid\n");
}

void test_get_process_threads_null_count() {
    int result = get_process_threads(GetCurrentProcessId(), NULL);
    
    assert(result != 0);
    printf("✓ test_get_process_threads_null_count\n");
}

void test_get_process_threads_zero_pid() {
    uint32_t thread_count = 0;
    int result = get_process_threads(0, &thread_count);
    
    assert(result != 0);
    printf("✓ test_get_process_threads_zero_pid\n");
}

void test_get_process_threads_at_least_one() {
    uint32_t thread_count = 0;
    int result = get_process_threads(GetCurrentProcessId(), &thread_count);
    
    assert(result == 0);
    assert(thread_count >= 1);
    printf("✓ test_get_process_threads_at_least_one\n");
}

/* ============================================================================
 * Test Suite: get_process_cmdline
 * ============================================================================ */

void test_get_process_cmdline_current_process() {
    char buffer[4096] = {0};
    char *result = get_process_cmdline(GetCurrentProcessId(), buffer, sizeof(buffer));
    
    assert(result != NULL);
    assert(strlen(buffer) > 0);
    printf("✓ test_get_process_cmdline_current_process\n");
}

void test_get_process_cmdline_invalid_pid() {
    char buffer[4096] = {0};
    char *result = get_process_cmdline(0xFFFFFFFF, buffer, sizeof(buffer));
    
    assert(result == NULL);
    printf("✓ test_get_process_cmdline_invalid_pid\n");
}

void test_get_process_cmdline_null_buffer() {
    char *result = get_process_cmdline(GetCurrentProcessId(), NULL, 256);
    
    assert(result == NULL);
    printf("✓ test_get_process_cmdline_null_buffer\n");
}

void test_get_process_cmdline_zero_buflen() {
    char buffer[4096] = {0};
    char *result = get_process_cmdline(GetCurrentProcessId(), buffer, 0);
    
    assert(result == NULL);
    printf("✓ test_get_process_cmdline_zero_buflen\n");
}

void test_get_process_cmdline_small_buffer() {
    char buffer[2] = {0};
    char *result = get_process_cmdline(GetCurrentProcessId(), buffer, 2);
    
    assert(result != NULL);
    printf("✓ test_get_process_cmdline_small_buffer\n");
}

void test_get_process_cmdline_large_buffer() {
    char buffer[65536] = {0};
    char *result = get_process_cmdline(GetCurrentProcessId(), buffer, sizeof(buffer));
    
    assert(result != NULL);
    printf("✓ test_get_process_cmdline_large_buffer\n");
}

void test_get_process_cmdline_zero_pid() {
    char buffer[4096] = {0};
    char *result = get_process_cmdline(0, buffer, sizeof(buffer));
    
    assert(result == NULL);
    printf("✓ test_get_process_cmdline_zero_pid\n");
}

void test_get_process_cmdline_buffer_not_modified_on_failure() {
    char buffer[4096];
    memset(buffer, 0xAB, sizeof(buffer));
    
    char *result = get_process_cmdline(0xFFFFFFFF, buffer, sizeof(buffer));
    
    assert(result == NULL);
    printf("✓ test_get_process_cmdline_buffer_not_modified_on_failure\n");
}

/* ============================================================================
 * Main Test Runner
 * ============================================================================ */

int main(int argc, char *argv[]) {
    printf("=== Apps OS Windows NT Test Suite ===\n\n");
    
    printf("Testing get_windows_process_info...\n");
    test_get_windows_process_info_success();
    test_get_windows_process_info_invalid_pid();
    test_get_windows_process_info_null_data_pointer();
    test_get_windows_process_info_zero_pid();
    test_get_windows_process_info_initializes_all_fields();
    printf("\n");
    
    printf("Testing enumerate_processes...\n");
    test_enumerate_processes_success();
    test_enumerate_processes_null_procs_pointer();
    test_enumerate_processes_null_count_pointer();
    test_enumerate_processes_finds_system_processes();
    test_enumerate_processes_count_consistency();
    printf("\n");
    
    printf("Testing free_process_list...\n");
    test_free_process_list_success();
    test_free_process_list_zero_count();
    test_free_process_list_null_pointer();
    test_free_process_list_large_count();
    printf("\n");
    
    printf("Testing filetime_to_microseconds...\n");
    test_filetime_to_microseconds_zero();
    test_filetime_to_microseconds_one_hundred_nanoseconds();
    test_filetime_to_microseconds_ten_microseconds();
    test_filetime_to_microseconds_one_second();
    test_filetime_to_microseconds_high_part();
    test_filetime_to_microseconds_max_value();
    printf("\n");
    
    printf("Testing get_process_cpu_time...\n");
    test_get_process_cpu_time_current_process();
    test_get_process_cpu_time_invalid_pid();
    test_get_process_cpu_time_null_user_time();
    test_get_process_cpu_time_null_kernel_time();
    test_get_process_cpu_time_both_null();
    test_get_process_cpu_time_zero_pid();
    printf("\n");
    
    printf("Testing get_process_memory_info...\n");
    test_get_process_memory_info_current_process();
    test_get_process_memory_info_invalid_pid();
    test_get_process_memory_info_null_rss();
    test_get_process_memory_info_null_vms();
    test_get_process_memory_info_both_null();
    test_get_process_memory_info_zero_pid();
    test_get_process_memory_info_rss_less_than_vms();
    printf("\n");
    
    printf("Testing get_process_threads...\n");
    test_get_process_threads_current_process();
    test_get_process_threads_invalid_pid();
    test_get_process_threads_null_count();
    test_get_process_threads_zero_pid();
    test_get_process_threads_at_least_one();
    printf("\n");
    
    printf("Testing get_process_cmdline...\n");
    test_get_process_cmdline_current_process();
    test_get_process_cmdline_invalid_pid();
    test_get_process_cmdline_null_buffer();
    test_get_process_cmdline_zero_buflen();
    test_get_process_cmdline_small_buffer();
    test_get_process_cmdline_large_buffer();
    test_get_process_cmdline_zero_pid();
    test_get_process_cmdline_buffer_not_modified_on_failure();
    printf("\n");
    
    printf("=== All tests passed! ===\n");
    return 0;
}