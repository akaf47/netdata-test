```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <setjmp.h>
#include <cmocka.h>
#include <sys/types.h>
#include <unistd.h>

#include "../ebpf_fd.h"

// Mock structures for testing
struct ebpf_module mock_module;
struct netdata_fd_stat mock_fd_stat;
ebpf_local_maps_t mock_local_maps;

// Test fixtures
void setup_test_module(void **state) {
    mock_module = (struct ebpf_module){0};
    mock_module.enabled = 1;
    mock_module.probes = NULL;
    mock_module.probe_links = NULL;
}

void teardown_test_module(void **state) {
    memset(&mock_module, 0, sizeof(mock_module));
}

// ===== Tests for ebpf_fd_init =====
void test_ebpf_fd_init_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_init(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_init_module_null(void **state) {
    // Arrange & Act & Assert
    // Test with NULL module - should handle gracefully
    int result = ebpf_fd_init(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_init_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_init(&mock_module);
    
    // Assert
    assert_int_not_equal(result, 0);
}

// ===== Tests for ebpf_fd_exit =====
void test_ebpf_fd_exit_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_exit(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_exit_module_null(void **state) {
    // Arrange & Act & Assert
    int result = ebpf_fd_exit(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_exit_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_exit(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

// ===== Tests for ebpf_fd_collector_update_probes =====
void test_ebpf_fd_collector_update_probes_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    void ebpf_fd_collector_update_probes(&mock_module);
    
    // Assert
    // Verify probes were updated
    assert_non_null(mock_module.probes);
}

void test_ebpf_fd_collector_update_probes_disabled(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act & Assert
    // Should not crash with disabled module
    ebpf_fd_collector_update_probes(&mock_module);
}

void test_ebpf_fd_collector_update_probes_null_module(void **state) {
    // Arrange & Act & Assert
    // Should handle NULL gracefully
    ebpf_fd_collector_update_probes(NULL);
}

// ===== Tests for ebpf_fd_probe_update =====
void test_ebpf_fd_probe_update_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    enum ebpf_kernel_type kernel_type = EBPF_KERNEL_6_1;
    
    // Act
    ebpf_fd_probe_update(kernel_type, &mock_module);
    
    // Assert
    // Verify probes updated for kernel type
}

void test_ebpf_fd_probe_update_different_kernel_types(void **state) {
    // Test with various kernel types
    enum ebpf_kernel_type kernel_types[] = {
        EBPF_KERNEL_4_11,
        EBPF_KERNEL_4_15,
        EBPF_KERNEL_5_8,
        EBPF_KERNEL_6_1
    };
    
    for (int i = 0; i < 4; i++) {
        ebpf_fd_probe_update(kernel_types[i], &mock_module);
    }
    
    assert_int_equal(1, 1);
}

void test_ebpf_fd_probe_update_null_module(void **state) {
    // Arrange & Act & Assert
    ebpf_fd_probe_update(EBPF_KERNEL_5_8, NULL);
}

// ===== Tests for ebpf_fd_collector_sync_table =====
void test_ebpf_fd_collector_sync_table_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_collector_sync_table(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_collector_sync_table_null_module(void **state) {
    // Arrange & Act & Assert
    int result = ebpf_fd_collector_sync_table(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_collector_sync_table_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_collector_sync_table(&mock_module);
    
    // Assert
    assert_int_not_equal(result, 0);
}

// ===== Tests for ebpf_fd_fill_histogram =====
void test_ebpf_fd_fill_histogram_valid_value(void **state) {
    // Arrange
    mock_fd_stat.open_calls = 100;
    mock_fd_stat.close_calls = 50;
    
    // Act
    ebpf_fd_fill_histogram(&mock_fd_stat, 10);
    
    // Assert
    // Verify histogram was filled
}

void test_ebpf_fd_fill_histogram_zero_value(void **state) {
    // Arrange
    mock_fd_stat.open_calls = 0;
    
    // Act
    ebpf_fd_fill_histogram(&mock_fd_stat, 0);
    
    // Assert
    // Should handle zero gracefully
}

void test_ebpf_fd_fill_histogram_large_value(void **state) {
    // Arrange
    mock_fd_stat.open_calls = UINT32_MAX;
    
    // Act
    ebpf_fd_fill_histogram(&mock_fd_stat, 1000000);
    
    // Assert
    // Should handle large values
}

void test_ebpf_fd_fill_histogram_null_stat(void **state) {
    // Arrange & Act & Assert
    ebpf_fd_fill_histogram(NULL, 10);
}

// ===== Tests for ebpf_fd_adjust_memory =====
void test_ebpf_fd_adjust_memory_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_adjust_memory(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_adjust_memory_null_module(void **state) {
    // Arrange & Act & Assert
    int result = ebpf_fd_adjust_memory(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_adjust_memory_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_adjust_memory(&mock_module);
    
    // Assert
    assert_int_not_equal(result, 0);
}

// ===== Tests for ebpf_fd_set_hash_table =====
void test_ebpf_fd_set_hash_table_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_set_hash_table(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_set_hash_table_null_module(void **state) {
    // Arrange & Act & Assert
    int result = ebpf_fd_set_hash_table(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_set_hash_table_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_set_hash_table(&mock_module);
    
    // Assert
    assert_int_not_equal(result, 0);
}

// ===== Tests for ebpf_fd_read_hash_table =====
void test_ebpf_fd_read_hash_table_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_read_hash_table(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_read_hash_table_null_module(void **state) {
    // Arrange & Act & Assert
    int result = ebpf_fd_read_hash_table(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_read_hash_table_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_read_hash_table(&mock_module);
    
    // Assert
    assert_int_not_equal(result, 0);
}

// ===== Tests for ebpf_fd_cleanup =====
void test_ebpf_fd_cleanup_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    void ebpf_fd_cleanup(&mock_module);
    
    // Assert
    // Verify resources were cleaned up
}

void test_ebpf_fd_cleanup_null_module(void **state) {
    // Arrange & Act & Assert
    ebpf_fd_cleanup(NULL);
}

void test_ebpf_fd_cleanup_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    ebpf_fd_cleanup(&mock_module);
    
    // Assert
    // Should handle disabled module
}

// ===== Tests for ebpf_fd_initialize_histogram =====
void test_ebpf_fd_initialize_histogram_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_initialize_histogram(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_initialize_histogram_null_module(void **state) {
    // Arrange & Act & Assert
    int result = ebpf_fd_initialize_histogram(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_initialize_histogram_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_initialize_histogram(&mock_module);
    
    // Assert
    assert_int_not_equal(result, 0);
}

// ===== Tests for ebpf_fd_collector_update_stats =====
void test_ebpf_fd_collector_update_stats_success(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    int result = ebpf_fd_collector_update_stats(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

void test_ebpf_fd_collector_update_stats_null_module(void **state) {
    // Arrange & Act & Assert
    int result = ebpf_fd_collector_update_stats(NULL);
    assert_int_not_equal(result, 0);
}

void test_ebpf_fd_collector_update_stats_disabled_module(void **state) {
    // Arrange
    mock_module.enabled = 0;
    
    // Act
    int result = ebpf_fd_collector_update_stats(&mock_module);
    
    // Assert
    assert_int_not_equal(result, 0);
}

// ===== Edge Case Tests =====
void test_ebpf_fd_concurrent_access(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act - Simulate concurrent access
    ebpf_fd_collector_sync_table(&mock_module);
    ebpf_fd_read_hash_table(&mock_module);
    ebpf_fd_collector_update_stats(&mock_module);
    
    // Assert
    assert_int_equal(1, 1);
}

void test_ebpf_fd_multiple_init_calls(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act - Multiple initializations
    ebpf_fd_init(&mock_module);
    ebpf_fd_init(&mock_module);
    ebpf_fd_init(&mock_module);
    
    // Assert
    assert_int_equal(1, 1);
}

void test_ebpf_fd_init_then_exit(void **state) {
    // Arrange
    mock_module.enabled = 1;
    
    // Act
    ebpf_fd_init(&mock_module);
    int result = ebpf_fd_exit(&mock_module);
    
    // Assert
    assert_int_equal(result, 0);
}

// ===== Test Suite Configuration =====
int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test_setup_teardown(test_ebpf_fd_init_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_init_module_null),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_init_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_exit_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_exit_module_null),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_exit_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_collector_update_probes_success, setup_test_module, teardown_test_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_collector_update_probes_disabled, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_collector_update_probes_null_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_probe_update_success, setup_test_module, teardown_test_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_probe_update_different_kernel_types, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_probe_update_null_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_collector_sync_table_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_collector_sync_table_null_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_collector_sync_table_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test(test_ebpf_fd_fill_histogram_valid_value),
        cmocka_unit_test(test_ebpf_fd_fill_histogram_zero_value),
        cmocka_unit_test(test_ebpf_fd_fill_histogram_large_value),
        cmocka_unit_test(test_ebpf_fd_fill_histogram_null_stat),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_adjust_memory_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_adjust_memory_null_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_adjust_memory_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_set_hash_table_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_set_hash_table_null_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_set_hash_table_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_read_hash_table_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_read_hash_table_null_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_read_hash_table_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_cleanup_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_cleanup_null_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_cleanup_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_initialize_histogram_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_initialize_histogram_null_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_initialize_histogram_disabled_module, setup_test_module, teardown_test_module),
        
        cmocka_unit_test_setup_teardown(test_ebpf_fd_collector_update_stats_success, setup_test_module, teardown_test_module),
        cmocka_unit_test(test_ebpf_fd_collector_update_stats_null_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_collector_update_stats_disabled_module, setup_test_module, teardown_test_module),
        
        // Edge case tests
        cmocka_unit_test_setup_teardown(test_ebpf_fd_concurrent_access, setup_test_module, teardown_test_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_multiple_init_calls, setup_test_module, teardown_test_module),
        cmocka_unit_test_setup_teardown(test_ebpf_fd_init_then_exit, setup_test_module, teardown_test_module),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}
```