#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <time.h>
#include <unistd.h>

/* Mock structures and function declarations */
typedef struct {
    const char *name;
    void (*function)(void);
    int enabled;
    long iterations;
    long total_time_ns;
} benchmark_entry;

/* Forward declarations for tested functions */
void benchmark_registry_init(void);
void benchmark_registry_cleanup(void);
int benchmark_registry_add(const char *name, void (*func)(void));
int benchmark_registry_remove(const char *name);
int benchmark_registry_execute(void);
int benchmark_registry_execute_single(const char *name);
benchmark_entry* benchmark_registry_get(const char *name);
int benchmark_registry_count(void);
void benchmark_registry_print_results(void);

/* Test helper functions */
static void dummy_benchmark_function(void) {
    volatile int i = 0;
    for (i = 0; i < 1000; i++);
}

static void dummy_benchmark_function_2(void) {
    volatile int i = 0;
    for (i = 0; i < 500; i++);
}

/* Test Suite: Registry Initialization and Cleanup */
static void test_benchmark_registry_init_success(void) {
    benchmark_registry_init();
    assert(benchmark_registry_count() == 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_init_success passed\n");
}

static void test_benchmark_registry_cleanup_after_init(void) {
    benchmark_registry_init();
    benchmark_registry_add("test1", dummy_benchmark_function);
    benchmark_registry_cleanup();
    assert(benchmark_registry_count() == 0);
    printf("✓ test_benchmark_registry_cleanup_after_init passed\n");
}

/* Test Suite: Adding Benchmarks */
static void test_benchmark_registry_add_single_benchmark(void) {
    benchmark_registry_init();
    int result = benchmark_registry_add("test_func", dummy_benchmark_function);
    assert(result == 0);
    assert(benchmark_registry_count() == 1);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_single_benchmark passed\n");
}

static void test_benchmark_registry_add_multiple_benchmarks(void) {
    benchmark_registry_init();
    int r1 = benchmark_registry_add("test1", dummy_benchmark_function);
    int r2 = benchmark_registry_add("test2", dummy_benchmark_function_2);
    int r3 = benchmark_registry_add("test3", dummy_benchmark_function);
    assert(r1 == 0);
    assert(r2 == 0);
    assert(r3 == 0);
    assert(benchmark_registry_count() == 3);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_multiple_benchmarks passed\n");
}

static void test_benchmark_registry_add_null_name(void) {
    benchmark_registry_init();
    int result = benchmark_registry_add(NULL, dummy_benchmark_function);
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_null_name passed\n");
}

static void test_benchmark_registry_add_null_function(void) {
    benchmark_registry_init();
    int result = benchmark_registry_add("test", NULL);
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_null_function passed\n");
}

static void test_benchmark_registry_add_empty_name(void) {
    benchmark_registry_init();
    int result = benchmark_registry_add("", dummy_benchmark_function);
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_empty_name passed\n");
}

static void test_benchmark_registry_add_duplicate_name(void) {
    benchmark_registry_init();
    benchmark_registry_add("duplicate", dummy_benchmark_function);
    int result = benchmark_registry_add("duplicate", dummy_benchmark_function_2);
    assert(result != 0);
    assert(benchmark_registry_count() == 1);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_duplicate_name passed\n");
}

static void test_benchmark_registry_add_long_name(void) {
    benchmark_registry_init();
    char long_name[512];
    memset(long_name, 'a', 511);
    long_name[511] = '\0';
    int result = benchmark_registry_add(long_name, dummy_benchmark_function);
    assert(result == 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_long_name passed\n");
}

/* Test Suite: Removing Benchmarks */
static void test_benchmark_registry_remove_existing(void) {
    benchmark_registry_init();
    benchmark_registry_add("to_remove", dummy_benchmark_function);
    int result = benchmark_registry_remove("to_remove");
    assert(result == 0);
    assert(benchmark_registry_count() == 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_remove_existing passed\n");
}

static void test_benchmark_registry_remove_nonexistent(void) {
    benchmark_registry_init();
    int result = benchmark_registry_remove("does_not_exist");
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_remove_nonexistent passed\n");
}

static void test_benchmark_registry_remove_null_name(void) {
    benchmark_registry_init();
    int result = benchmark_registry_remove(NULL);
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_remove_null_name passed\n");
}

static void test_benchmark_registry_remove_empty_name(void) {
    benchmark_registry_init();
    int result = benchmark_registry_remove("");
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_remove_empty_name passed\n");
}

static void test_benchmark_registry_remove_from_multiple(void) {
    benchmark_registry_init();
    benchmark_registry_add("keep1", dummy_benchmark_function);
    benchmark_registry_add("remove_me", dummy_benchmark_function_2);
    benchmark_registry_add("keep2", dummy_benchmark_function);
    int result = benchmark_registry_remove("remove_me");
    assert(result == 0);
    assert(benchmark_registry_count() == 2);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_remove_from_multiple passed\n");
}

/* Test Suite: Getting Benchmarks */
static void test_benchmark_registry_get_existing(void) {
    benchmark_registry_init();
    benchmark_registry_add("test_get", dummy_benchmark_function);
    benchmark_entry *entry = benchmark_registry_get("test_get");
    assert(entry != NULL);
    assert(strcmp(entry->name, "test_get") == 0);
    assert(entry->function == dummy_benchmark_function);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_get_existing passed\n");
}

static void test_benchmark_registry_get_nonexistent(void) {
    benchmark_registry_init();
    benchmark_entry *entry = benchmark_registry_get("nonexistent");
    assert(entry == NULL);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_get_nonexistent passed\n");
}

static void test_benchmark_registry_get_null_name(void) {
    benchmark_registry_init();
    benchmark_entry *entry = benchmark_registry_get(NULL);
    assert(entry == NULL);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_get_null_name passed\n");
}

static void test_benchmark_registry_get_empty_name(void) {
    benchmark_registry_init();
    benchmark_entry *entry = benchmark_registry_get("");
    assert(entry == NULL);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_get_empty_name passed\n");
}

static void test_benchmark_registry_get_case_sensitive(void) {
    benchmark_registry_init();
    benchmark_registry_add("TestFunc", dummy_benchmark_function);
    benchmark_entry *entry1 = benchmark_registry_get("TestFunc");
    benchmark_entry *entry2 = benchmark_registry_get("testfunc");
    assert(entry1 != NULL);
    assert(entry2 == NULL);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_get_case_sensitive passed\n");
}

/* Test Suite: Executing Benchmarks */
static void test_benchmark_registry_execute_all(void) {
    benchmark_registry_init();
    benchmark_registry_add("exec1", dummy_benchmark_function);
    benchmark_registry_add("exec2", dummy_benchmark_function_2);
    int result = benchmark_registry_execute();
    assert(result == 0 || result == 2);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_all passed\n");
}

static void test_benchmark_registry_execute_empty(void) {
    benchmark_registry_init();
    int result = benchmark_registry_execute();
    assert(result == 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_empty passed\n");
}

static void test_benchmark_registry_execute_single_existing(void) {
    benchmark_registry_init();
    benchmark_registry_add("single", dummy_benchmark_function);
    int result = benchmark_registry_execute_single("single");
    assert(result == 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_single_existing passed\n");
}

static void test_benchmark_registry_execute_single_nonexistent(void) {
    benchmark_registry_init();
    int result = benchmark_registry_execute_single("nonexistent");
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_single_nonexistent passed\n");
}

static void test_benchmark_registry_execute_single_null_name(void) {
    benchmark_registry_init();
    int result = benchmark_registry_execute_single(NULL);
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_single_null_name passed\n");
}

static void test_benchmark_registry_execute_single_empty_name(void) {
    benchmark_registry_init();
    int result = benchmark_registry_execute_single("");
    assert(result != 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_single_empty_name passed\n");
}

static void test_benchmark_registry_execute_updates_timing(void) {
    benchmark_registry_init();
    benchmark_registry_add("timing_test", dummy_benchmark_function);
    benchmark_entry *entry = benchmark_registry_get("timing_test");
    long time_before = entry->total_time_ns;
    benchmark_registry_execute_single("timing_test");
    long time_after = entry->total_time_ns;
    assert(time_after >= time_before);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_updates_timing passed\n");
}

/* Test Suite: Registry Count */
static void test_benchmark_registry_count_empty(void) {
    benchmark_registry_init();
    assert(benchmark_registry_count() == 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_count_empty passed\n");
}

static void test_benchmark_registry_count_multiple(void) {
    benchmark_registry_init();
    for (int i = 0; i < 5; i++) {
        char name[32];
        snprintf(name, sizeof(name), "bench_%d", i);
        benchmark_registry_add(name, dummy_benchmark_function);
    }
    assert(benchmark_registry_count() == 5);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_count_multiple passed\n");
}

/* Test Suite: Print Results */
static void test_benchmark_registry_print_results_empty(void) {
    benchmark_registry_init();
    benchmark_registry_print_results();
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_print_results_empty passed\n");
}

static void test_benchmark_registry_print_results_with_data(void) {
    benchmark_registry_init();
    benchmark_registry_add("print_test", dummy_benchmark_function);
    benchmark_registry_execute_single("print_test");
    benchmark_registry_print_results();
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_print_results_with_data passed\n");
}

/* Test Suite: Complex Scenarios */
static void test_benchmark_registry_add_remove_add_cycle(void) {
    benchmark_registry_init();
    benchmark_registry_add("cycle", dummy_benchmark_function);
    assert(benchmark_registry_count() == 1);
    benchmark_registry_remove("cycle");
    assert(benchmark_registry_count() == 0);
    benchmark_registry_add("cycle", dummy_benchmark_function_2);
    assert(benchmark_registry_count() == 1);
    benchmark_entry *entry = benchmark_registry_get("cycle");
    assert(entry->function == dummy_benchmark_function_2);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_add_remove_add_cycle passed\n");
}

static void test_benchmark_registry_stress_many_additions(void) {
    benchmark_registry_init();
    int count = 100;
    for (int i = 0; i < count; i++) {
        char name[32];
        snprintf(name, sizeof(name), "stress_%d", i);
        int result = benchmark_registry_add(name, dummy_benchmark_function);
        assert(result == 0);
    }
    assert(benchmark_registry_count() == count);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_stress_many_additions passed\n");
}

static void test_benchmark_registry_stress_many_removals(void) {
    benchmark_registry_init();
    int count = 50;
    for (int i = 0; i < count; i++) {
        char name[32];
        snprintf(name, sizeof(name), "remove_%d", i);
        benchmark_registry_add(name, dummy_benchmark_function);
    }
    for (int i = 0; i < count; i++) {
        char name[32];
        snprintf(name, sizeof(name), "remove_%d", i);
        int result = benchmark_registry_remove(name);
        assert(result == 0);
    }
    assert(benchmark_registry_count() == 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_stress_many_removals passed\n");
}

static void test_benchmark_registry_execute_all_after_removals(void) {
    benchmark_registry_init();
    benchmark_registry_add("exec_a", dummy_benchmark_function);
    benchmark_registry_add("exec_b", dummy_benchmark_function_2);
    benchmark_registry_add("exec_c", dummy_benchmark_function);
    benchmark_registry_remove("exec_b");
    int result = benchmark_registry_execute();
    assert(result == 0 || result == 2);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_execute_all_after_removals passed\n");
}

/* Test Suite: Edge Cases for Timing */
static void test_benchmark_registry_timing_precision(void) {
    benchmark_registry_init();
    benchmark_registry_add("precision", dummy_benchmark_function);
    for (int i = 0; i < 3; i++) {
        benchmark_registry_execute_single("precision");
    }
    benchmark_entry *entry = benchmark_registry_get("precision");
    assert(entry->total_time_ns > 0);
    benchmark_registry_cleanup();
    printf("✓ test_benchmark_registry_timing_precision passed\n");
}

/* Main test runner */
int main(void) {
    printf("\n========== BENCHMARK REGISTRY TEST SUITE ==========\n\n");
    
    /* Initialization tests */
    test_benchmark_registry_init_success();
    test_benchmark_registry_cleanup_after_init();
    
    /* Addition tests */
    test_benchmark_registry_add_single_benchmark();
    test_benchmark_registry_add_multiple_benchmarks();
    test_benchmark_registry_add_null_name();
    test_benchmark_registry_add_null_function();
    test_benchmark_registry_add_empty_name();
    test_benchmark_registry_add_duplicate_name();
    test_benchmark_registry_add_long_name();
    
    /* Removal tests */
    test_benchmark_registry_remove_existing();
    test_benchmark_registry_remove_nonexistent();
    test_benchmark_registry_remove_null_name();
    test_benchmark_registry_remove_empty_name();
    test_benchmark_registry_remove_from_multiple();
    
    /* Get tests */
    test_benchmark_registry_get_existing();
    test_benchmark_registry_get_nonexistent();
    test_benchmark_registry_get_null_name();
    test_benchmark_registry_get_empty_name();
    test_benchmark_registry_get_case_sensitive();
    
    /* Execution tests */
    test_benchmark_registry_execute_all();
    test_benchmark_registry_execute_empty();
    test_benchmark_registry_execute_single_existing();
    test_benchmark_registry_execute_single_nonexistent();
    test_benchmark_registry_execute_single_null_name();
    test_benchmark_registry_execute_single_empty_name();
    test_benchmark_registry_execute_updates_timing();
    
    /* Count tests */
    test_benchmark_registry_count_empty();
    test_benchmark_registry_count_multiple();
    
    /* Print results tests */
    test_benchmark_registry_print_results_empty();
    test_benchmark_registry_print_results_with_data();
    
    /* Complex scenario tests */
    test_benchmark_registry_add_remove_add_cycle();
    test_benchmark_registry_stress_many_additions();
    test_benchmark_registry_stress_many_removals();
    test_benchmark_registry_execute_all_after_removals();
    
    /* Edge case tests */
    test_benchmark_registry_timing_precision();
    
    printf("\n========== ALL TESTS PASSED ==========\n\n");
    return 0;
}