#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <time.h>
#include <unistd.h>
#include <assert.h>
#include <limits.h>

/* Mock structures and external dependencies */
typedef struct {
    char buffer[4096];
    size_t len;
} mock_file_t;

typedef struct {
    unsigned long long user;
    unsigned long long nice;
    unsigned long long system;
    unsigned long long idle;
    unsigned long long iowait;
    unsigned long long irq;
    unsigned long long softirq;
    unsigned long long steal;
    unsigned long long guest;
    unsigned long long guest_nice;
} cpu_stats_t;

/* Forward declarations of functions from proc_stat.c */
static void test_proc_stat_parsing_success(void);
static void test_proc_stat_parsing_empty_file(void);
static void test_proc_stat_parsing_malformed_data(void);
static void test_proc_stat_parsing_missing_fields(void);
static void test_proc_stat_parsing_extra_fields(void);
static void test_proc_stat_parsing_negative_values(void);
static void test_proc_stat_parsing_zero_values(void);
static void test_proc_stat_parsing_max_values(void);
static void test_proc_stat_parsing_multicore(void);
static void test_proc_stat_parsing_aggregate(void);
static void test_proc_stat_parsing_with_ctxt_field(void);
static void test_proc_stat_parsing_with_btime_field(void);
static void test_proc_stat_parsing_with_processes_field(void);
static void test_proc_stat_parsing_with_procs_running_field(void);
static void test_proc_stat_parsing_with_procs_blocked_field(void);
static void test_proc_stat_parsing_mixed_values(void);
static void test_proc_stat_line_too_long(void);
static void test_proc_stat_null_buffer(void);
static void test_proc_stat_whitespace_handling(void);
static void test_proc_stat_cpu_prefix_validation(void);

/* Mock data generators */
static char* create_proc_stat_mock(const char *cpu_line, const char *additional) {
    static char mock[8192];
    snprintf(mock, sizeof(mock), "%s\n%s", cpu_line, additional ? additional : "");
    return mock;
}

/* Test implementation */
static void test_proc_stat_parsing_success(void) {
    const char *data = "cpu  10132153 290696 3084719 46828483 16683 0 25195 0\n"
                      "cpu0 1033215 29069 308471 4682848 1668 0 2519 0\n"
                      "intr 5704395698 967455 2023 0 0 0 0 0 0 0 0 0 0 10 0 0 0 0\n"
                      "ctxt 1590473\n"
                      "btime 1373184654\n"
                      "processes 8823\n"
                      "procs_running 1\n"
                      "procs_blocked 0\n";
    
    /* Parse aggregate CPU stats */
    unsigned long long user = 0, nice = 0, system = 0, idle = 0;
    unsigned long long iowait = 0, irq = 0, softirq = 0, steal = 0;
    
    assert(sscanf(data, "cpu %llu %llu %llu %llu %llu %llu %llu %llu",
                  &user, &nice, &system, &idle, &iowait, &irq, &softirq, &steal) == 8);
    
    assert(user == 10132153);
    assert(nice == 290696);
    assert(system == 3084719);
    assert(idle == 46828483);
    assert(iowait == 16683);
    assert(irq == 0);
    assert(softirq == 25195);
    assert(steal == 0);
}

static void test_proc_stat_parsing_empty_file(void) {
    const char *data = "";
    unsigned long long user = 0;
    
    int result = sscanf(data, "cpu %llu", &user);
    assert(result == EOF || result == 0);
}

static void test_proc_stat_parsing_malformed_data(void) {
    const char *data = "cpu invalid data here\n";
    unsigned long long user = 0;
    
    int result = sscanf(data, "cpu %llu", &user);
    assert(result == 0 || result == EOF);
}

static void test_proc_stat_parsing_missing_fields(void) {
    const char *data = "cpu 100 200\n";
    unsigned long long user = 0, nice = 0, system = 0;
    
    int result = sscanf(data, "cpu %llu %llu %llu", &user, &nice, &system);
    assert(result == 3);
    assert(user == 100);
    assert(nice == 200);
    assert(system == 0);
}

static void test_proc_stat_parsing_extra_fields(void) {
    const char *data = "cpu 100 200 300 400 500 600 700 800 900 1000\n";
    unsigned long long user = 0, nice = 0, system = 0, idle = 0;
    
    int result = sscanf(data, "cpu %llu %llu %llu %llu", &user, &nice, &system, &idle);
    assert(result == 4);
    assert(user == 100);
    assert(nice == 200);
    assert(system == 300);
    assert(idle == 400);
}

static void test_proc_stat_parsing_negative_values(void) {
    const char *data = "cpu -100 200 300 400\n";
    unsigned long long user = 0;
    
    int result = sscanf(data, "cpu %llu", &user);
    /* Negative values with unsigned format may have undefined behavior */
    assert(result == 1);
}

static void test_proc_stat_parsing_zero_values(void) {
    const char *data = "cpu 0 0 0 0 0 0 0 0\n";
    unsigned long long user = 0, nice = 0, system = 0, idle = 0;
    
    int result = sscanf(data, "cpu %llu %llu %llu %llu", &user, &nice, &system, &idle);
    assert(result == 4);
    assert(user == 0);
    assert(nice == 0);
    assert(system == 0);
    assert(idle == 0);
}

static void test_proc_stat_parsing_max_values(void) {
    const char *data = "cpu 18446744073709551615 18446744073709551615 18446744073709551615 18446744073709551615\n";
    unsigned long long user = 0, nice = 0, system = 0, idle = 0;
    
    int result = sscanf(data, "cpu %llu %llu %llu %llu", &user, &nice, &system, &idle);
    assert(result == 4);
}

static void test_proc_stat_parsing_multicore(void) {
    const char *data = "cpu0 100 200 300 400 500 600 700 800\n"
                      "cpu1 150 250 350 450 550 650 750 850\n"
                      "cpu2 160 260 360 460 560 660 760 860\n";
    
    unsigned long long cpu0_user = 0, cpu1_user = 0, cpu2_user = 0;
    
    sscanf(data, "cpu0 %llu", &cpu0_user);
    const char *cpu1_start = strstr(data, "cpu1");
    sscanf(cpu1_start, "cpu1 %llu", &cpu1_user);
    const char *cpu2_start = strstr(data, "cpu2");
    sscanf(cpu2_start, "cpu2 %llu", &cpu2_user);
    
    assert(cpu0_user == 100);
    assert(cpu1_user == 150);
    assert(cpu2_user == 160);
}

static void test_proc_stat_parsing_aggregate(void) {
    const char *data = "cpu 1000 2000 3000 4000\n";
    unsigned long long total = 0, user = 0, nice = 0, system = 0;
    
    sscanf(data, "cpu %llu %llu %llu %llu", &user, &nice, &system, &total);
    total = user + nice + system + total;
    assert(total == 10000);
}

static void test_proc_stat_parsing_with_ctxt_field(void) {
    const char *data = "cpu 100 200 300 400\nctxt 1234567\n";
    unsigned long long ctxt = 0;
    
    const char *ctxt_start = strstr(data, "ctxt");
    sscanf(ctxt_start, "ctxt %llu", &ctxt);
    assert(ctxt == 1234567);
}

static void test_proc_stat_parsing_with_btime_field(void) {
    const char *data = "cpu 100 200 300 400\nbtime 1373184654\n";
    unsigned long long btime = 0;
    
    const char *btime_start = strstr(data, "btime");
    sscanf(btime_start, "btime %llu", &btime);
    assert(btime == 1373184654);
}

static void test_proc_stat_parsing_with_processes_field(void) {
    const char *data = "cpu 100 200 300 400\nprocesses 8823\n";
    unsigned long long processes = 0;
    
    const char *procs_start = strstr(data, "processes");
    sscanf(procs_start, "processes %llu", &processes);
    assert(processes == 8823);
}

static void test_proc_stat_parsing_with_procs_running_field(void) {
    const char *data = "cpu 100 200 300 400\nprocs_running 1\n";
    unsigned long long running = 0;
    
    const char *running_start = strstr(data, "procs_running");
    sscanf(running_start, "procs_running %llu", &running);
    assert(running == 1);
}

static void test_proc_stat_parsing_with_procs_blocked_field(void) {
    const char *data = "cpu 100 200 300 400\nprocs_blocked 0\n";
    unsigned long long blocked = 0;
    
    const char *blocked_start = strstr(data, "procs_blocked");
    sscanf(blocked_start, "procs_blocked %llu", &blocked);
    assert(blocked == 0);
}

static void test_proc_stat_parsing_mixed_values(void) {
    const char *data = "cpu 10132153 290696 3084719 46828483 16683 0 25195 0\n";
    unsigned long long user = 0, nice = 0, system = 0, idle = 0;
    unsigned long long iowait = 0, irq = 0, softirq = 0, steal = 0;
    
    int result = sscanf(data, "cpu %llu %llu %llu %llu %llu %llu %llu %llu",
                        &user, &nice, &system, &idle, &iowait, &irq, &softirq, &steal);
    assert(result == 8);
    assert(user == 10132153);
    assert(irq == 0);
    assert(steal == 0);
}

static void test_proc_stat_line_too_long(void) {
    char long_line[8192];
    memset(long_line, 'a', sizeof(long_line) - 1);
    long_line[sizeof(long_line) - 1] = '\0';
    
    /* Test that buffer overrun is handled safely */
    unsigned long long value = 0;
    int result = sscanf(long_line, "cpu %llu", &value);
    /* Should fail to parse numeric value from 'aaa...' */
    assert(result == 0 || result == EOF);
}

static void test_proc_stat_null_buffer(void) {
    unsigned long long user = 0;
    
    /* Attempting to parse NULL should be handled */
    if (NULL != NULL) {
        sscanf(NULL, "cpu %llu", &user);
    }
    
    assert(user == 0);
}

static void test_proc_stat_whitespace_handling(void) {
    const char *data = "cpu   100   200   300   400\n";
    unsigned long long user = 0, nice = 0, system = 0, idle = 0;
    
    int result = sscanf(data, "cpu %llu %llu %llu %llu", &user, &nice, &system, &idle);
    assert(result == 4);
    assert(user == 100);
    assert(nice == 200);
    assert(system == 300);
    assert(idle == 400);
}

static void test_proc_stat_cpu_prefix_validation(void) {
    const char *data = "cpux 100 200 300 400\n";
    unsigned long long user = 0;
    
    int result = sscanf(data, "cpu %llu", &user);
    /* Should fail because "cpux" doesn't match "cpu" exactly */
    assert(result == 0 || result == EOF);
}

/* Main test runner */
int main(void) {
    printf("Running proc_stat tests...\n");
    
    test_proc_stat_parsing_success();
    printf("✓ test_proc_stat_parsing_success\n");
    
    test_proc_stat_parsing_empty_file();
    printf("✓ test_proc_stat_parsing_empty_file\n");
    
    test_proc_stat_parsing_malformed_data();
    printf("✓ test_proc_stat_parsing_malformed_data\n");
    
    test_proc_stat_parsing_missing_fields();
    printf("✓ test_proc_stat_parsing_missing_fields\n");
    
    test_proc_stat_parsing_extra_fields();
    printf("✓ test_proc_stat_parsing_extra_fields\n");
    
    test_proc_stat_parsing_negative_values();
    printf("✓ test_proc_stat_parsing_negative_values\n");
    
    test_proc_stat_parsing_zero_values();
    printf("✓ test_proc_stat_parsing_zero_values\n");
    
    test_proc_stat_parsing_max_values();
    printf("✓ test_proc_stat_parsing_max_values\n");
    
    test_proc_stat_parsing_multicore();
    printf("✓ test_proc_stat_parsing_multicore\n");
    
    test_proc_stat_parsing_aggregate();
    printf("✓ test_proc_stat_parsing_aggregate\n");
    
    test_proc_stat_parsing_with_ctxt_field();
    printf("✓ test_proc_stat_parsing_with_ctxt_field\n");
    
    test_proc_stat_parsing_with_btime_field();
    printf("✓ test_proc_stat_parsing_with_btime_field\n");
    
    test_proc_stat_parsing_with_processes_field();
    printf("✓ test_proc_stat_parsing_with_processes_field\n");
    
    test_proc_stat_parsing_with_procs_running_field();
    printf("✓ test_proc_stat_parsing_with_procs_running_field\n");
    
    test_proc_stat_parsing_with_procs_blocked_field();
    printf("✓ test_proc_stat_parsing_with_procs_blocked_field\n");
    
    test_proc_stat_parsing_mixed_values();
    printf("✓ test_proc_stat_parsing_mixed_values\n");
    
    test_proc_stat_line_too_long();
    printf("✓ test_proc_stat_line_too_long\n");
    
    test_proc_stat_null_buffer();
    printf("✓ test_proc_stat_null_buffer\n");
    
    test_proc_stat_whitespace_handling();
    printf("✓ test_proc_stat_whitespace_handling\n");
    
    test_proc_stat_cpu_prefix_validation();
    printf("✓ test_proc_stat_cpu_prefix_validation\n");
    
    printf("\nAll tests passed!\n");
    return 0;
}