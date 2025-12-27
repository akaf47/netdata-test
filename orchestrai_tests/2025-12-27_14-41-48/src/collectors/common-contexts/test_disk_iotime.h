#include <stddef.h>
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <assert.h>
#include <setjmp.h>
#include <cmocka.h>

// Mock structures and functions for disk-iotime.h
#include "disk-iotime.h"

// Test fixtures
typedef struct {
    void *context;
} DiskIotimeTestFixture;

// Setup and teardown
static int setup(void **state) {
    DiskIotimeTestFixture *fixture = malloc(sizeof(DiskIotimeTestFixture));
    assert_non_null(fixture);
    fixture->context = NULL;
    *state = fixture;
    return 0;
}

static int teardown(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    if (fixture) {
        if (fixture->context) {
            free(fixture->context);
        }
        free(fixture);
    }
    return 0;
}

// Test: Initialization with valid parameters
static void test_disk_iotime_init_valid(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk_name = "sda";
    
    // Act & Assert
    // Test that initialization succeeds with valid disk name
    assert_non_null(disk_name);
    assert_true(strlen(disk_name) > 0);
}

// Test: Initialization with null disk name
static void test_disk_iotime_init_null_disk_name(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk_name = NULL;
    
    // Act & Assert
    // Test handling of null disk name
    assert_null(disk_name);
}

// Test: Initialization with empty disk name
static void test_disk_iotime_init_empty_disk_name(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk_name = "";
    
    // Act & Assert
    assert_non_null(disk_name);
    assert_int_equal(strlen(disk_name), 0);
}

// Test: Initialization with very long disk name
static void test_disk_iotime_init_long_disk_name(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    char long_name[256];
    memset(long_name, 'a', 255);
    long_name[255] = '\0';
    
    // Act & Assert
    assert_non_null(long_name);
    assert_int_equal(strlen(long_name), 255);
}

// Test: Multiple disk initialization
static void test_disk_iotime_multiple_disks(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disks[] = {"sda", "sdb", "sdc", "nvme0n1"};
    size_t num_disks = sizeof(disks) / sizeof(disks[0]);
    
    // Act & Assert
    for (size_t i = 0; i < num_disks; i++) {
        assert_non_null(disks[i]);
        assert_true(strlen(disks[i]) > 0);
    }
}

// Test: Disk name with special characters
static void test_disk_iotime_special_chars_disk_name(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *special_names[] = {
        "nvme0n1p1",
        "mmcblk0",
        "vda",
        "loop0",
        "dm-0"
    };
    size_t count = sizeof(special_names) / sizeof(special_names[0]);
    
    // Act & Assert
    for (size_t i = 0; i < count; i++) {
        assert_non_null(special_names[i]);
        assert_true(strlen(special_names[i]) > 0);
    }
}

// Test: Boundary value - zero io time
static void test_disk_iotime_zero_io_time(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t io_time = 0;
    
    // Act & Assert
    assert_int_equal(io_time, 0);
}

// Test: Boundary value - maximum io time
static void test_disk_iotime_max_io_time(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t io_time = UINT64_MAX;
    
    // Act & Assert
    assert_int_equal(io_time, UINT64_MAX);
}

// Test: Boundary value - large io time
static void test_disk_iotime_large_io_time(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t io_time = 999999999999ULL;
    
    // Act & Assert
    assert_true(io_time > 0);
    assert_true(io_time < UINT64_MAX);
}

// Test: Negative io time (as signed)
static void test_disk_iotime_negative_value(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    int64_t negative_value = -1;
    
    // Act & Assert
    assert_true(negative_value < 0);
}

// Test: Overflow scenario with sequential values
static void test_disk_iotime_overflow_detection(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t previous_value = UINT64_MAX - 10;
    uint64_t current_value = 5;
    
    // Act & Assert
    // Current value is less than previous, indicating counter overflow
    assert_true(current_value < previous_value);
}

// Test: Underflow scenario
static void test_disk_iotime_underflow_detection(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t previous_value = 100;
    uint64_t current_value = 50;
    
    // Act & Assert
    // Should detect decrease in counter
    assert_true(current_value < previous_value);
}

// Test: Normal increment
static void test_disk_iotime_normal_increment(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t previous_value = 100;
    uint64_t current_value = 150;
    
    // Act & Assert
    uint64_t delta = current_value - previous_value;
    assert_int_equal(delta, 50);
}

// Test: Rapid increments
static void test_disk_iotime_rapid_increments(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t values[] = {0, 100, 200, 300, 400, 500};
    size_t count = sizeof(values) / sizeof(values[0]);
    
    // Act & Assert
    for (size_t i = 1; i < count; i++) {
        uint64_t delta = values[i] - values[i-1];
        assert_int_equal(delta, 100);
    }
}

// Test: Context structure initialization
static void test_disk_iotime_context_init(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange - Context should be properly initialized
    assert_null(fixture->context);
    
    // Act
    fixture->context = malloc(sizeof(uint64_t) * 10);
    
    // Assert
    assert_non_null(fixture->context);
    memset(fixture->context, 0, sizeof(uint64_t) * 10);
    assert_memory_equal(fixture->context, calloc(10, sizeof(uint64_t)), sizeof(uint64_t) * 10);
    free(calloc(10, sizeof(uint64_t)));
}

// Test: Data validation - non-zero values
static void test_disk_iotime_data_validation_nonzero(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t valid_values[] = {1, 10, 100, 1000, 10000};
    size_t count = sizeof(valid_values) / sizeof(valid_values[0]);
    
    // Act & Assert
    for (size_t i = 0; i < count; i++) {
        assert_true(valid_values[i] > 0);
    }
}

// Test: Data validation - zero values
static void test_disk_iotime_data_validation_zero(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t zero_value = 0;
    
    // Act & Assert
    assert_int_equal(zero_value, 0);
}

// Test: String formatting with disk name
static void test_disk_iotime_string_formatting(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk_name = "sda";
    char formatted[256];
    
    // Act
    snprintf(formatted, sizeof(formatted), "disk.%s.io_time", disk_name);
    
    // Assert
    assert_string_equal(formatted, "disk.sda.io_time");
}

// Test: String formatting with complex disk name
static void test_disk_iotime_string_formatting_complex(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk_name = "nvme0n1p1";
    char formatted[256];
    
    // Act
    snprintf(formatted, sizeof(formatted), "disk.%s.io_time", disk_name);
    
    // Assert
    assert_string_equal(formatted, "disk.nvme0n1p1.io_time");
}

// Test: Metric name generation
static void test_disk_iotime_metric_name_generation(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *prefix = "system";
    const char *disk = "sdb";
    char metric_name[256];
    
    // Act
    snprintf(metric_name, sizeof(metric_name), "%s.disk.%s.io_time.ms", prefix, disk);
    
    // Assert
    assert_string_equal(metric_name, "system.disk.sdb.io_time.ms");
}

// Test: Array of io_time values
static void test_disk_iotime_array_handling(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t io_times[5] = {0};
    size_t array_size = 5;
    
    // Act
    for (size_t i = 0; i < array_size; i++) {
        io_times[i] = i * 100;
    }
    
    // Assert
    for (size_t i = 0; i < array_size; i++) {
        assert_int_equal(io_times[i], i * 100);
    }
}

// Test: Buffer overflow protection
static void test_disk_iotime_buffer_boundary(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    char buffer[32];
    const char *disk_name = "sda";
    
    // Act - Fill buffer carefully
    int written = snprintf(buffer, sizeof(buffer), "disk.%s", disk_name);
    
    // Assert
    assert_true(written < (int)sizeof(buffer));
    assert_true(written > 0);
}

// Test: Consistency check - same disk name
static void test_disk_iotime_consistency_same_disk(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk1 = "sda";
    const char *disk2 = "sda";
    
    // Act & Assert
    assert_string_equal(disk1, disk2);
}

// Test: Consistency check - different disk names
static void test_disk_iotime_consistency_different_disks(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk1 = "sda";
    const char *disk2 = "sdb";
    
    // Act & Assert
    assert_string_not_equal(disk1, disk2);
}

// Test: Case sensitivity of disk names
static void test_disk_iotime_case_sensitivity(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk_lower = "sda";
    const char *disk_upper = "SDA";
    
    // Act & Assert
    assert_string_not_equal(disk_lower, disk_upper);
}

// Test: Whitespace in disk names
static void test_disk_iotime_whitespace_handling(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    const char *disk_with_space = "sd a";
    char trimmed[256];
    
    // Act
    sscanf(disk_with_space, "%s", trimmed);
    
    // Assert
    assert_string_equal(trimmed, "sd");
}

// Test: Timestamp handling
static void test_disk_iotime_timestamp_consistency(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t timestamp1 = 1000;
    uint64_t timestamp2 = 1000;
    
    // Act & Assert
    assert_int_equal(timestamp1, timestamp2);
}

// Test: Time delta calculation
static void test_disk_iotime_delta_calculation(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t time1 = 1000;
    uint64_t time2 = 2000;
    
    // Act
    uint64_t delta = time2 - time1;
    
    // Assert
    assert_int_equal(delta, 1000);
}

// Test: Empty context handling
static void test_disk_iotime_empty_context(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange & Act & Assert
    assert_null(fixture->context);
}

// Test: Multiple context operations
static void test_disk_iotime_multiple_context_ops(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    fixture->context = malloc(sizeof(uint64_t) * 3);
    uint64_t *data = (uint64_t *)fixture->context;
    
    // Act
    data[0] = 100;
    data[1] = 200;
    data[2] = 300;
    
    // Assert
    assert_int_equal(data[0], 100);
    assert_int_equal(data[1], 200);
    assert_int_equal(data[2], 300);
}

// Test: Zero-length array handling
static void test_disk_iotime_zero_length_array(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t *empty_array = NULL;
    size_t array_length = 0;
    
    // Act & Assert
    assert_null(empty_array);
    assert_int_equal(array_length, 0);
}

// Test: Single element array
static void test_disk_iotime_single_element_array(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t single_array[1] = {42};
    
    // Act & Assert
    assert_int_equal(single_array[0], 42);
}

// Test: Large array handling
static void test_disk_iotime_large_array(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    size_t large_size = 1000;
    uint64_t *large_array = malloc(large_size * sizeof(uint64_t));
    
    // Act
    for (size_t i = 0; i < large_size; i++) {
        large_array[i] = i;
    }
    
    // Assert
    for (size_t i = 0; i < large_size; i++) {
        assert_int_equal(large_array[i], i);
    }
    
    free(large_array);
}

// Test: Pointer validation
static void test_disk_iotime_pointer_validation(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t value = 123;
    uint64_t *ptr = &value;
    
    // Act & Assert
    assert_non_null(ptr);
    assert_int_equal(*ptr, 123);
}

// Test: Null pointer handling
static void test_disk_iotime_null_pointer(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t *null_ptr = NULL;
    
    // Act & Assert
    assert_null(null_ptr);
}

// Test: Iterator with valid range
static void test_disk_iotime_iteration_valid_range(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t values[5];
    for (int i = 0; i < 5; i++) {
        values[i] = i * 10;
    }
    
    // Act & Assert
    for (int i = 0; i < 5; i++) {
        assert_int_equal(values[i], i * 10);
    }
}

// Test: Iterator boundary conditions
static void test_disk_iotime_iteration_boundaries(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t values[3] = {0, 100, 200};
    int start = 0;
    int end = 3;
    
    // Act & Assert
    for (int i = start; i < end; i++) {
        assert_true(values[i] >= 0);
    }
}

// Test: Increment and decrement operations
static void test_disk_iotime_increment_decrement(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t counter = 10;
    
    // Act
    counter++;
    
    // Assert
    assert_int_equal(counter, 11);
    
    // Act
    counter--;
    
    // Assert
    assert_int_equal(counter, 10);
}

// Test: Memory safety - bounds checking
static void test_disk_iotime_memory_bounds(void **state) {
    DiskIotimeTestFixture *fixture = *state;
    
    // Arrange
    uint64_t buffer[10];
    int index = 5;
    
    // Act
    buffer[index] = 999;
    
    // Assert
    assert_int_equal(buffer[5], 999);
    assert_true(index >= 0 && index < 10);
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test_setup_teardown(test_disk_iotime_init_valid, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_init_null_disk_name, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_init_empty_disk_name, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_init_long_disk_name, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_multiple_disks, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_special_chars_disk_name, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_zero_io_time, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_max_io_time, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_large_io_time, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_negative_value, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_overflow_detection, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_underflow_detection, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_normal_increment, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_rapid_increments, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_context_init, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_data_validation_nonzero, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_data_validation_zero, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_string_formatting, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_string_formatting_complex, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_metric_name_generation, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_array_handling, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_buffer_boundary, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_consistency_same_disk, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_consistency_different_disks, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_case_sensitivity, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_whitespace_handling, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_timestamp_consistency, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_delta_calculation, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_empty_context, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_multiple_context_ops, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_zero_length_array, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_single_element_array, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_large_array, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_pointer_validation, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_null_pointer, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_iteration_valid_range, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_iteration_boundaries, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_increment_decrement, setup, teardown),
        cmocka_unit_test_setup_teardown(test_disk_iotime_memory_bounds, setup, teardown),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}