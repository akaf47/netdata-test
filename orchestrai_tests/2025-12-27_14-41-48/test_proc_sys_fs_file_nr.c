#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <assert.h>
#include <limits.h>

/* Mock structures */
typedef struct {
    unsigned long long allocated;
    unsigned long long unused;
    unsigned long long max;
} file_nr_stats_t;

/* Forward declarations of test functions */
static void test_file_nr_parsing_success(void);
static void test_file_nr_parsing_empty_file(void);
static void test_file_nr_parsing_malformed_data(void);
static void test_file_nr_parsing_missing_fields(void);
static void test_file_nr_parsing_extra_fields(void);
static void test_file_nr_parsing_zero_values(void);
static void test_file_nr_parsing_max_values(void);
static void test_file_nr_parsing_negative_values(void);
static void test_file_nr_parsing_with_whitespace(void);
static void test_file_nr_allocated_zero(void);
static void test_file_nr_unused_zero(void);
static void test_file_nr_max_zero(void);
static void test_file_nr_allocated_equals_max(void);
static void test_file_nr_allocated_greater_than_unused(void);
static void test_file_nr_mixed_values(void);
static void test_file_nr_large_allocated(void);
static void test_file_nr_large_unused(void);
static void test_file_nr_large_max(void);
static void test_file_nr_leading_zeros(void);
static void test_file_nr_tab_separated(void);
static void test_file_nr_space_separated(void);

/* Mock data generators */
static char* create_file_nr_line(unsigned long long allocated, unsigned long long unused, unsigned long long max) {
    static char line[256];
    snprintf(line, sizeof(line), "%llu\t%llu\t%llu\n", allocated, unused, max);
    return line;
}

/* Test implementations */
static void test_file_nr_parsing_success(void) {
    const char *data = "1024\t512\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
    assert(allocated == 1024);
    assert(unused == 512);
    assert(max == 2048);
}

static void test_file_nr_parsing_empty_file(void) {
    const char *data = "";
    unsigned long long allocated = 0;
    
    int result = sscanf(data, "%llu", &allocated);
    assert(result == EOF || result == 0);
}

static void test_file_nr_parsing_malformed_data(void) {
    const char *data = "not a number\n";
    unsigned long long allocated = 0;
    
    int result = sscanf(data, "%llu", &allocated);
    assert(result == 0 || result == EOF);
}

static void test_file_nr_parsing_missing_fields(void) {
    const char *data = "1024\t512\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 2);
    assert(allocated == 1024);
    assert(unused == 512);
}

static void test_file_nr_parsing_extra_fields(void) {
    const char *data = "1024\t512\t2048\t999\t888\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
    assert(allocated == 1024);
    assert(unused == 512);
    assert(max == 2048);
}

static void test_file_nr_parsing_zero_values(void) {
    const char *data = "0\t0\t0\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
    assert(allocated == 0);
    assert(unused == 0);
    assert(max == 0);
}

static void test_file_nr_parsing_max_values(void) {
    const char *data = "18446744073709551615\t18446744073709551615\t18446744073709551615\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
}

static void test_file_nr_parsing_negative_values(void) {
    const char *data = "-1024\t512\t2048\n";
    unsigned long long allocated = 0;
    
    int result = sscanf(data, "%llu", &allocated);
    /* Negative values with unsigned format may have undefined behavior */
    assert(result == 1);
}

static void test_file_nr_parsing_with_whitespace(void) {
    const char *data = "  1024  \t  512  \t  2048  \n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
    assert(allocated == 1024);
    assert(unused == 512);
    assert(max == 2048);
}

static void test_file_nr_allocated_zero(void) {
    const char *data = "0\t512\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(allocated == 0);
    assert(unused == 512);
    assert(max == 2048);
}

static void test_file_nr_unused_zero(void) {
    const char *data = "1024\t0\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(allocated == 1024);
    assert(unused == 0);
    assert(max == 2048);
}

static void test_file_nr_max_zero(void) {
    const char *data = "1024\t512\t0\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(allocated == 1024);
    assert(unused == 512);
    assert(max == 0);
}

static void test_file_nr_allocated_equals_max(void) {
    const char *data = "2048\t512\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(allocated == max);
}

static void test_file_nr_allocated_greater_than_unused(void) {
    const char *data = "1024\t512\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(allocated > unused);
}

static void test_file_nr_mixed_values(void) {
    const char *data = "12345\t6789\t99999\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
    assert(allocated == 12345);
    assert(unused == 6789);
    assert(max == 99999);
}

static void test_file_nr_large_allocated(void) {
    const char *data = "999999999\t512\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(allocated == 999999999);
}

static void test_file_nr_large_unused(void) {
    const char *data = "1024\t888888888\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(unused == 888888888);
}

static void test_file_nr_large_max(void) {
    const char *data = "1024\t512\t777777777\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(max == 777777777);
}

static void test_file_nr_leading_zeros(void) {
    const char *data = "001024\t000512\t002048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
    assert(allocated == 1024);
    assert(unused == 512);
    assert(max == 2048);
}

static void test_file_nr_tab_separated(void) {
    const char *data = "1024\t512\t2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
}

static void test_file_nr_space_separated(void) {
    const char *data = "1024 512 2048\n";
    unsigned long long allocated = 0, unused = 0, max = 0;
    
    int result = sscanf(data, "%llu %llu %llu", &allocated, &unused, &max);
    assert(result == 3);
}

/* Main test runner */
int main(void) {
    printf("Running proc_sys_fs_file_nr tests...\n");
    
    test_file_nr_parsing_success();
    printf("✓ test_file_nr_parsing_success\n");
    
    test_file_nr_parsing_empty_file();
    printf("✓ test_file_nr_parsing_empty_file\n");
    
    test_file_nr_parsing_malformed_data();
    printf("✓ test_file_nr_parsing_malformed_data\n");
    
    test_file_nr_parsing_missing_fields();
    printf("✓ test_file_nr_parsing_missing_fields\n");
    
    test_file_nr_parsing_extra_fields();
    printf("✓ test_file_nr_parsing_extra_fields\n");
    
    test_file_nr_parsing_zero_values();
    printf("✓ test_file_nr_parsing_zero_values\n");
    
    test_file_nr_parsing_max_values();
    printf("✓ test_file_nr_parsing_max_values\n");
    
    test_file_nr_parsing_negative_values();
    printf("✓ test_file_nr_parsing_negative_values\n");
    
    test_file_nr_parsing_with_whitespace();
    printf("✓ test_file_nr_parsing_with_whitespace\n");
    
    test_file_nr_allocated_zero();
    printf("✓ test_file_nr_allocated_zero\n");
    
    test_file_nr_unused_zero();
    printf("✓ test_file_nr_unused_zero\n");
    
    test_file_nr_max_zero();
    printf("✓ test_file_nr_max_zero\n");
    
    test_file_nr_allocated_equals_max();
    printf("✓ test_file_nr_allocated_equals_max\n");
    
    test_file_nr_allocated_greater_than_unused();
    printf("✓ test_file_nr_allocated_greater_than_unused\n");
    
    test_file_nr_mixed_values();
    printf("✓ test_file_nr_mixed_values\n");
    
    test_file_nr_large_allocated();
    printf("✓ test_file_nr_large_allocated\n");
    
    test_file_nr_large_unused();
    printf("✓ test_file_nr_large_unused\n");
    
    test_file_nr_large_max();
    printf("✓ test_file_nr_large_max\n");
    
    test_file_nr_leading_zeros();
    printf("✓ test_file_nr_leading_zeros\n");
    
    test_file_nr_tab_separated();
    printf("✓ test_file_nr_tab_separated\n");
    
    test_file_nr_space_separated();
    printf("✓ test_file_nr_space_separated\n");
    
    printf("\nAll tests passed!\n");
    return 0;
}