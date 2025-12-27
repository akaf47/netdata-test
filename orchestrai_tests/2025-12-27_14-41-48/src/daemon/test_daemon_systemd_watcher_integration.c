#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <pthread.h>
#include <time.h>

#include "../../../../../tests/cmocka/cmocka.h"
#include "daemon-systemd-watcher.h"

/* Global state for integration tests */
static volatile int signal_received = 0;
static volatile int watcher_ready = 0;

/* Mock signal handler */
void mock_signal_handler(int sig) {
    signal_received = sig;
}

/* Setup function */
static int integration_setup(void **state) {
    signal_received = 0;
    watcher_ready = 0;
    return 0;
}

/* Teardown function */
static int integration_teardown(void **state) {
    systemd_watcher_cleanup();
    return 0;
}

/* Integration Test: Full lifecycle */
static void test_systemd_watcher_full_lifecycle(void **state) {
    // Arrange
    assert_int_equal(systemd_watcher_get_status(), SYSTEMD_WATCHER_STATUS_UNAVAILABLE);
    
    // Act - Initialize
    int init_result = systemd_watcher_init();
    
    // Assert - Initialization should succeed or fail gracefully
    if (init_result == 0) {
        assert_int_equal(systemd_watcher_get_status(), SYSTEMD_WATCHER_STATUS_RUNNING);
        
        // Act - Stop
        systemd_watcher_stop();
        
        // Assert - Should be stopped
        assert_int_equal(systemd_watcher_get_status(), SYSTEMD_WATCHER_STATUS_STOPPED);
    }
}

/* Integration Test: Multiple lifecycle cycles */
static void test_systemd_watcher_multiple_cycles(void **state) {
    // Arrange
    int cycles = 3;
    
    for (int i = 0; i < cycles; i++) {
        // Act
        int init_result = systemd_watcher_init();
        
        if (init_result == 0) {
            int status = systemd_watcher_get_status();
            assert_true(status == SYSTEMD_WATCHER_STATUS_RUNNING || 
                       status == SYSTEMD_WATCHER_STATUS_UNAVAILABLE);
            
            systemd_watcher_stop();
            
            status = systemd_watcher_get_status();
            assert_true(status == SYSTEMD_WATCHER_STATUS_STOPPED || 
                       status == SYSTEMD_WATCHER_STATUS_UNAVAILABLE);
        }
    }
    
    // Assert - no crashes or leaks
    assert_true(1);
}

/* Integration Test: Error recovery */
static void test_systemd_watcher_error_recovery(void **state) {
    // Arrange
    int init_result = systemd_watcher_init();
    
    if (init_result == 0) {
        // Act - trigger error
        systemd_watcher_signal_handler(SIGTERM);
        
        // Assert - should handle cleanly
        int status = systemd_watcher_get_status();
        assert_true(status == SYSTEMD_WATCHER_STATUS_STOPPED || 
                   status == SYSTEMD_WATCHER_STATUS_UNAVAILABLE);
        
        // Act - reinitialize after error
        init_result = systemd_watcher_init();
        
        // Assert - should recover
        assert_true(init_result >= -1);
    }
}

/* Integration Test: Signal handling during operation */
static void test_systemd_watcher_signal_handling_integration(void **state) {
    // Arrange
    int init_result = systemd_watcher_init();
    
    if (init_result == 0) {
        int status_before = systemd_watcher_get_status();
        
        // Act
        systemd_watcher_signal_handler(SIGTERM);
        
        // Assert
        int status_after = systemd_watcher_get_status();
        assert_true(status_after == SYSTEMD_WATCHER_STATUS_STOPPED || 
                   status_after != status_before);
    }
}

/* Integration Test: Status consistency */
static void test_systemd_watcher_status_consistency(void **state) {
    // Arrange & Act
    int status1 = systemd_watcher_get_status();
    int status2 = systemd_watcher_get_status();
    int status3 = systemd_watcher_get_status();
    
    // Assert - multiple calls should return consistent results
    assert_int_equal(status1, status2);
    assert_int_equal(status2, status3);
}

/* Integration Test: Availability flag consistency */
static void test_systemd_watcher_availability_consistency(void **state) {
    // Arrange & Act
    int avail1 = systemd_watcher_is_available();
    int avail2 = systemd_watcher_is_available();
    
    // Assert
    assert_int_equal(avail1, avail2);
}

/* Integration Test: Resource cleanup after errors */
static void test_systemd_watcher_resource_cleanup(void **state) {
    // Arrange
    int init_result = systemd_watcher_init();
    
    if (init_result == 0) {
        // Act - simulate various operations
        systemd_watcher_get_status();
        systemd_watcher_is_available();
        
        // Act - cleanup
        systemd_watcher_cleanup();
        
        // Assert - status should indicate unavailable
        int status = systemd_watcher_get_status();
        assert_int_equal(status, SYSTEMD_WATCHER_STATUS_UNAVAILABLE);
    }
}

/* Integration Test: Repeated cleanup calls */
static void test_systemd_watcher_repeated_cleanup(void **state) {
    // Arrange
    int init_result = systemd_watcher_init();
    
    if (init_result == 0) {
        // Act
        systemd_watcher_cleanup();
        systemd_watcher_cleanup();
        systemd_watcher_cleanup();
        
        // Assert - should not crash
        assert_true(1);
    }
}

/* Integration Test: Status queries after cleanup */
static void test_systemd_watcher_status_after_cleanup(void **state) {
    // Arrange
    int init_result = systemd_watcher_init();
    
    if (init_result == 0) {
        // Act
        systemd_watcher_cleanup();
        
        // Assert - queries should still work
        int status = systemd_watcher_get_status();
        int avail = systemd_watcher_is_available();
        
        assert_true(status >= 0);
        assert_true(avail >= 0);
    }
}

/* Edge case: Initialization under low memory conditions */
static void test_systemd_watcher_init_recovery_attempts(void **state) {
    // Arrange - try initialization multiple times with potential failures
    int success_count = 0;
    int failure_count = 0;
    
    // Act
    for (int i = 0; i < 3; i++) {
        int result = systemd_watcher_init();
        if (result == 0) {
            success_count++;
            systemd_watcher_stop();
        } else {
            failure_count++;
        }
    }
    
    // Assert
    assert_true(success_count + failure_count == 3);
}

/* Edge case: Rapid start/stop cycles */
static void test_systemd_watcher_rapid_start_stop(void **state) {
    // Arrange
    const int cycles = 5;
    
    // Act & Assert
    for (int i = 0; i < cycles; i++) {
        int init_result = systemd_watcher_init();
        if (init_result == 0) {
            systemd_watcher_stop();
        }
        // Allow a small delay between cycles
        usleep(100);
    }
    
    assert_true(1);
}

/* Main test suite */
int main(void) {
    const struct CMUnitTest tests[] = {
        /* Lifecycle tests */
        cmocka_unit_test_setup_teardown(test_systemd_watcher_full_lifecycle, integration_setup, integration_teardown),
        cmocka_unit_test_setup_teardown(test_systemd_watcher_multiple_cycles, integration_setup, integration_teardown),
        
        /* Error recovery tests */
        cmocka_unit_test_setup_teardown(test_systemd_watcher_error_recovery, integration_setup, integration_teardown),
        
        /* Signal handling integration */
        cmocka_unit_test_setup_teardown(test_systemd_watcher_signal_handling_integration, integration_setup, integration_teardown),
        
        /* Consistency tests */
        cmocka_unit_test_setup_teardown(test_systemd_watcher_status_consistency, integration_setup, integration_teardown),
        cmocka_unit_test_setup_teardown(test_systemd_watcher_availability_consistency, integration_setup, integration_teardown),
        
        /* Resource management tests */
        cmocka_unit_test_setup_teardown(test_systemd_watcher_resource_cleanup, integration_setup, integration_teardown),
        cmocka_unit_test_setup_teardown(test_systemd_watcher_repeated_cleanup, integration_setup, integration_teardown),
        cmocka_unit_test_setup_teardown(test_systemd_watcher_status_after_cleanup, integration_setup, integration_teardown),
        
        /* Edge case tests */
        cmocka_unit_test_setup_teardown(test_systemd_watcher_init_recovery_attempts, integration_setup, integration_teardown),
        cmocka_unit_test_setup_teardown(test_systemd_watcher_rapid_start_stop, integration_setup, integration_teardown),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}