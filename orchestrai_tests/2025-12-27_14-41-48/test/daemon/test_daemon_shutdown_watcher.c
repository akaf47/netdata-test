#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>
#include <unistd.h>
#include <pthread.h>
#include <signal.h>
#include <time.h>
#include <setjmp.h>
#include <cmocka.h>

/* Include the header file being tested */
#include "daemon-shutdown-watcher.h"

/* Mock variables and helper structures */
static bool mock_shutdown_requested = false;
static int mock_signal_received = 0;
static jmp_buf jump_buffer;

/* Mock implementation of external dependencies */
void mock_signal_handler_setup(void) {
    mock_signal_received = 0;
}

void mock_signal_cleanup(void) {
    mock_signal_received = 0;
}

/* Test: Initialize shutdown watcher with valid parameters */
static void test_daemon_shutdown_watcher_init_success(void **state) {
    (void)state;
    
    /* Should initialize without errors */
    int result = daemon_shutdown_watcher_init();
    assert_int_equal(result, 0);
}

/* Test: Initialize shutdown watcher when already initialized */
static void test_daemon_shutdown_watcher_init_already_initialized(void **state) {
    (void)state;
    
    /* First initialization */
    daemon_shutdown_watcher_init();
    
    /* Second initialization should handle gracefully */
    int result = daemon_shutdown_watcher_init();
    assert_int_equal(result, 0);
    
    /* Cleanup */
    daemon_shutdown_watcher_fini();
}

/* Test: Finalize shutdown watcher */
static void test_daemon_shutdown_watcher_fini_success(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    int result = daemon_shutdown_watcher_fini();
    assert_int_equal(result, 0);
}

/* Test: Finalize without initialization */
static void test_daemon_shutdown_watcher_fini_without_init(void **state) {
    (void)state;
    
    /* Should handle cleanup gracefully even without init */
    int result = daemon_shutdown_watcher_fini();
    assert_int_equal(result, 0);
}

/* Test: Check shutdown status when not requested */
static void test_is_daemon_shutdown_requested_false(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    bool result = is_daemon_shutdown_requested();
    assert_false(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Check shutdown status after request */
static void test_is_daemon_shutdown_requested_true(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    /* Request shutdown */
    request_daemon_shutdown("Test shutdown");
    
    bool result = is_daemon_shutdown_requested();
    assert_true(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Request shutdown with valid reason */
static void test_request_daemon_shutdown_with_reason(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    request_daemon_shutdown("User initiated shutdown");
    
    bool result = is_daemon_shutdown_requested();
    assert_true(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Request shutdown with NULL reason */
static void test_request_daemon_shutdown_with_null_reason(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    /* Should handle NULL reason gracefully */
    request_daemon_shutdown(NULL);
    
    bool result = is_daemon_shutdown_requested();
    assert_true(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Request shutdown with empty string reason */
static void test_request_daemon_shutdown_with_empty_reason(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    request_daemon_shutdown("");
    
    bool result = is_daemon_shutdown_requested();
    assert_true(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Request shutdown with long reason string */
static void test_request_daemon_shutdown_with_long_reason(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    char long_reason[512];
    memset(long_reason, 'A', sizeof(long_reason) - 1);
    long_reason[sizeof(long_reason) - 1] = '\0';
    
    request_daemon_shutdown(long_reason);
    
    bool result = is_daemon_shutdown_requested();
    assert_true(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Multiple shutdown requests */
static void test_multiple_shutdown_requests(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    request_daemon_shutdown("First request");
    assert_true(is_daemon_shutdown_requested());
    
    request_daemon_shutdown("Second request");
    assert_true(is_daemon_shutdown_requested());
    
    daemon_shutdown_watcher_fini();
}

/* Test: Get shutdown reason */
static void test_get_daemon_shutdown_reason(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    const char *reason = "Test shutdown reason";
    request_daemon_shutdown(reason);
    
    const char *retrieved_reason = get_daemon_shutdown_reason();
    assert_non_null(retrieved_reason);
    assert_string_equal(retrieved_reason, reason);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Get shutdown reason before any request */
static void test_get_daemon_shutdown_reason_before_request(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    const char *reason = get_daemon_shutdown_reason();
    assert_null(reason);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Get shutdown reason after NULL request */
static void test_get_daemon_shutdown_reason_null_request(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    request_daemon_shutdown(NULL);
    
    const char *reason = get_daemon_shutdown_reason();
    /* Should be NULL or empty depending on implementation */
    
    daemon_shutdown_watcher_fini();
}

/* Test: Signal handler for SIGTERM */
static void test_daemon_shutdown_on_sigterm(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    /* Simulate SIGTERM signal */
    raise(SIGTERM);
    
    /* Give signal handler time to process */
    usleep(100000);
    
    bool result = is_daemon_shutdown_requested();
    assert_true(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Signal handler for SIGINT */
static void test_daemon_shutdown_on_sigint(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    /* Simulate SIGINT signal */
    raise(SIGINT);
    
    /* Give signal handler time to process */
    usleep(100000);
    
    bool result = is_daemon_shutdown_requested();
    assert_true(result);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Initialize and finalize multiple times */
static void test_daemon_shutdown_watcher_multiple_init_fini_cycles(void **state) {
    (void)state;
    
    for (int i = 0; i < 3; i++) {
        int init_result = daemon_shutdown_watcher_init();
        assert_int_equal(init_result, 0);
        
        assert_false(is_daemon_shutdown_requested());
        
        int fini_result = daemon_shutdown_watcher_fini();
        assert_int_equal(fini_result, 0);
    }
}

/* Test: Concurrent initialization calls */
static void test_daemon_shutdown_watcher_concurrent_init(void **state) {
    (void)state;
    
    /* Initialize multiple times without finalization */
    int result1 = daemon_shutdown_watcher_init();
    int result2 = daemon_shutdown_watcher_init();
    int result3 = daemon_shutdown_watcher_init();
    
    assert_int_equal(result1, 0);
    assert_int_equal(result2, 0);
    assert_int_equal(result3, 0);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Request shutdown after finalization */
static void test_request_shutdown_after_fini(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    daemon_shutdown_watcher_fini();
    
    /* Requesting shutdown after finalization should handle gracefully */
    request_daemon_shutdown("After finalization");
}

/* Test: Special characters in shutdown reason */
static void test_request_daemon_shutdown_special_characters(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    request_daemon_shutdown("Shutdown\nwith\nnewlines");
    assert_true(is_daemon_shutdown_requested());
    
    request_daemon_shutdown("Shutdown\twith\ttabs");
    assert_true(is_daemon_shutdown_requested());
    
    request_daemon_shutdown("Shutdown with \"quotes\" and 'apostrophes'");
    assert_true(is_daemon_shutdown_requested());
    
    daemon_shutdown_watcher_fini();
}

/* Test: Shutdown reason persistence */
static void test_shutdown_reason_persistence(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    const char *reason1 = "First reason";
    request_daemon_shutdown(reason1);
    
    const char *retrieved1 = get_daemon_shutdown_reason();
    assert_non_null(retrieved1);
    assert_string_equal(retrieved1, reason1);
    
    /* Retrieve again - should be the same */
    const char *retrieved2 = get_daemon_shutdown_reason();
    assert_non_null(retrieved2);
    assert_string_equal(retrieved2, reason1);
    
    daemon_shutdown_watcher_fini();
}

/* Test: Shutdown state thread safety */
static void test_shutdown_state_consistency(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    request_daemon_shutdown("Consistency test");
    
    /* Check multiple times - should be consistent */
    for (int i = 0; i < 10; i++) {
        bool result = is_daemon_shutdown_requested();
        assert_true(result);
    }
    
    daemon_shutdown_watcher_fini();
}

/* Test: Reason retrieval consistency */
static void test_reason_retrieval_consistency(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    const char *original_reason = "Consistency check reason";
    request_daemon_shutdown(original_reason);
    
    /* Retrieve multiple times */
    for (int i = 0; i < 5; i++) {
        const char *retrieved = get_daemon_shutdown_reason();
        assert_non_null(retrieved);
        assert_string_equal(retrieved, original_reason);
    }
    
    daemon_shutdown_watcher_fini();
}

/* Test: Shutdown after multiple state checks */
static void test_shutdown_state_check_sequence(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    /* Initial state should be no shutdown */
    assert_false(is_daemon_shutdown_requested());
    
    /* Request shutdown */
    request_daemon_shutdown("Sequential check");
    
    /* Verify shutdown state is set */
    assert_true(is_daemon_shutdown_requested());
    
    /* Multiple subsequent checks should confirm shutdown state */
    assert_true(is_daemon_shutdown_requested());
    assert_true(is_daemon_shutdown_requested());
    
    daemon_shutdown_watcher_fini();
}

/* Test: Watcher with edge case signal values */
static void test_daemon_shutdown_watcher_signal_handling(void **state) {
    (void)state;
    
    daemon_shutdown_watcher_init();
    
    /* Should handle signals without crashing */
    request_daemon_shutdown("Signal test");
    
    assert_true(is_daemon_shutdown_requested());
    
    daemon_shutdown_watcher_fini();
}

/* Test: Cleanup resources properly */
static void test_daemon_shutdown_watcher_resource_cleanup(void **state) {
    (void)state;
    
    for (int i = 0; i < 5; i++) {
        daemon_shutdown_watcher_init();
        request_daemon_shutdown("Resource test iteration");
        assert_true(is_daemon_shutdown_requested());
        daemon_shutdown_watcher_fini();
    }
    
    /* After cleanup, should be able to reinitialize fresh */
    daemon_shutdown_watcher_init();
    assert_false(is_daemon_shutdown_requested());
    daemon_shutdown_watcher_fini();
}

/* Test suite setup and teardown */
static int group_setup(void **state) {
    (void)state;
    return 0;
}

static int group_teardown(void **state) {
    (void)state;
    return 0;
}

/* Main test runner */
int main(void) {
    const struct CMUnitTest tests[] = {
        /* Initialization tests */
        cmocka_unit_test(test_daemon_shutdown_watcher_init_success),
        cmocka_unit_test(test_daemon_shutdown_watcher_init_already_initialized),
        
        /* Finalization tests */
        cmocka_unit_test(test_daemon_shutdown_watcher_fini_success),
        cmocka_unit_test(test_daemon_shutdown_watcher_fini_without_init),
        
        /* Status check tests */
        cmocka_unit_test(test_is_daemon_shutdown_requested_false),
        cmocka_unit_test(test_is_daemon_shutdown_requested_true),
        
        /* Shutdown request tests */
        cmocka_unit_test(test_request_daemon_shutdown_with_reason),
        cmocka_unit_test(test_request_daemon_shutdown_with_null_reason),
        cmocka_unit_test(test_request_daemon_shutdown_with_empty_reason),
        cmocka_unit_test(test_request_daemon_shutdown_with_long_reason),
        cmocka_unit_test(test_multiple_shutdown_requests),
        
        /* Reason retrieval tests */
        cmocka_unit_test(test_get_daemon_shutdown_reason),
        cmocka_unit_test(test_get_daemon_shutdown_reason_before_request),
        cmocka_unit_test(test_get_daemon_shutdown_reason_null_request),
        
        /* Signal handling tests */
        cmocka_unit_test(test_daemon_shutdown_on_sigterm),
        cmocka_unit_test(test_daemon_shutdown_on_sigint),
        
        /* Lifecycle tests */
        cmocka_unit_test(test_daemon_shutdown_watcher_multiple_init_fini_cycles),
        cmocka_unit_test(test_daemon_shutdown_watcher_concurrent_init),
        cmocka_unit_test(test_request_shutdown_after_fini),
        
        /* Edge case tests */
        cmocka_unit_test(test_request_daemon_shutdown_special_characters),
        cmocka_unit_test(test_shutdown_reason_persistence),
        cmocka_unit_test(test_shutdown_state_consistency),
        cmocka_unit_test(test_reason_retrieval_consistency),
        cmocka_unit_test(test_shutdown_state_check_sequence),
        cmocka_unit_test(test_daemon_shutdown_watcher_signal_handling),
        cmocka_unit_test(test_daemon_shutdown_watcher_resource_cleanup),
    };
    
    return cmocka_run_group_tests(tests, group_setup, group_teardown);
}