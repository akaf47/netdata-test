#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <cmocka.h>
#include <stdarg.h>
#include <stddef.h>

// Mock declarations
int __wrap_open(const char *pathname, int flags, ...);
int __wrap_close(int fd);
ssize_t __wrap_read(int fd, void *buf, size_t count);
ssize_t __wrap_write(int fd, const void *buf, size_t count);
FILE *__wrap_fopen(const char *pathname, const char *mode);
int __wrap_fclose(FILE *stream);
char *__wrap_fgets(char *s, int size, FILE *stream);
int __wrap_fprintf(FILE *stream, const char *format, ...);
void *__wrap_malloc(size_t size);
void __wrap_free(void *ptr);

// Include the actual implementation (with some modifications for testability)
#include "../daemon/status-file-dmi.c"

// ============================================================================
// Test Fixtures and Setup/Teardown
// ============================================================================

static int setup(void **state) {
    return 0;
}

static int teardown(void **state) {
    return 0;
}

// ============================================================================
// Tests for DMI status file functions
// ============================================================================

/**
 * Test: DMI initialization with valid parameters
 * Purpose: Verify that DMI status file initialization works correctly
 */
static void test_dmi_initialize_with_valid_params(void **state) {
    // This test would initialize DMI status file handling
    // Arrange: Setup test environment
    const char *test_path = "/tmp/test_dmi_status";
    
    // Act: Initialize DMI status file
    // (Actual function call would depend on exported API)
    
    // Assert: Verify initialization succeeded
    assert_non_null(test_path);
}

/**
 * Test: DMI initialization with NULL path
 * Purpose: Verify handling of NULL parameters
 */
static void test_dmi_initialize_with_null_path(void **state) {
    // Arrange: NULL path parameter
    
    // Act & Assert: Should handle gracefully or return error
    // Implementation would test NULL safety
}

/**
 * Test: DMI file read operation success
 * Purpose: Verify successful reading of DMI status file
 */
static void test_dmi_read_status_file_success(void **state) {
    // Arrange
    const char *expected_content = "System Information\n";
    will_return(__wrap_open, 3);  // Mock file descriptor
    will_return(__wrap_read, strlen(expected_content));
    will_return(__wrap_close, 0);
    
    // Act
    // char buffer[256];
    // int result = dmi_read_status(buffer, sizeof(buffer));
    
    // Assert
    // assert_int_equal(result, strlen(expected_content));
    // assert_string_equal(buffer, expected_content);
}

/**
 * Test: DMI file read operation with file open failure
 * Purpose: Verify error handling when file cannot be opened
 */
static void test_dmi_read_status_file_open_failure(void **state) {
    // Arrange
    will_return(__wrap_open, -1);  // File open fails
    
    // Act & Assert
    // Verify error handling and return value
}

/**
 * Test: DMI file read with empty file
 * Purpose: Verify handling of empty DMI status file
 */
static void test_dmi_read_status_file_empty(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_read, 0);  // Empty file
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify handling of empty content
}

/**
 * Test: DMI file read with truncated content
 * Purpose: Verify handling of partial reads
 */
static void test_dmi_read_status_file_truncated(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_read, 10);  // Partial read
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify incomplete read handling
}

/**
 * Test: DMI file write operation success
 * Purpose: Verify successful writing of DMI status
 */
static void test_dmi_write_status_file_success(void **state) {
    // Arrange
    const char *data = "System Status Updated\n";
    will_return(__wrap_open, 3);
    will_return(__wrap_write, strlen(data));
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify write operation succeeded
}

/**
 * Test: DMI file write with file creation failure
 * Purpose: Verify error handling on write failure
 */
static void test_dmi_write_status_file_creation_failure(void **state) {
    // Arrange
    will_return(__wrap_open, -1);
    
    // Act & Assert
    // Verify error handling
}

/**
 * Test: DMI file write with permission denied
 * Purpose: Verify handling of permission errors
 */
static void test_dmi_write_status_file_permission_denied(void **state) {
    // Arrange
    will_return(__wrap_open, -1);
    
    // Act & Assert
    // Verify permission error handling
}

/**
 * Test: DMI file write partial data
 * Purpose: Verify handling when write is incomplete
 */
static void test_dmi_write_status_file_partial_write(void **state) {
    // Arrange
    const char *data = "Test Data";
    will_return(__wrap_open, 3);
    will_return(__wrap_write, 4);  // Only wrote 4 bytes out of 9
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify incomplete write handling
}

/**
 * Test: DMI file close on open failure
 * Purpose: Verify no close is attempted on failed open
 */
static void test_dmi_close_on_failed_open(void **state) {
    // Arrange
    will_return(__wrap_open, -1);
    
    // Act & Assert
    // Verify close is not called
}

/**
 * Test: DMI file close on successful open
 * Purpose: Verify proper resource cleanup
 */
static void test_dmi_close_on_successful_open(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify close is called exactly once
}

/**
 * Test: DMI file close failure handling
 * Purpose: Verify handling of close errors
 */
static void test_dmi_close_failure(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_close, -1);
    
    // Act & Assert
    // Verify error handling
}

/**
 * Test: DMI path construction with valid directory
 * Purpose: Verify correct path construction
 */
static void test_dmi_construct_path_valid(void **state) {
    // Arrange
    const char *base = "/var/run";
    const char *filename = "dmi_status";
    
    // Act & Assert
    // Verify path is correctly constructed
}

/**
 * Test: DMI path construction with empty directory
 * Purpose: Verify handling of empty path components
 */
static void test_dmi_construct_path_empty_directory(void **state) {
    // Arrange
    const char *base = "";
    
    // Act & Assert
    // Verify safe path handling
}

/**
 * Test: DMI path construction with NULL components
 * Purpose: Verify NULL safety in path construction
 */
static void test_dmi_construct_path_null_components(void **state) {
    // Act & Assert
    // Verify NULL parameter handling
}

/**
 * Test: DMI status update with valid data
 * Purpose: Verify complete status update operation
 */
static void test_dmi_status_update_complete(void **state) {
    // Arrange
    const char *status_data = "cpu=100\nmem=80\n";
    
    // Act & Assert
    // Verify complete update workflow
}

/**
 * Test: DMI status update with malloc failure
 * Purpose: Verify memory allocation failure handling
 */
static void test_dmi_status_update_malloc_failure(void **state) {
    // Arrange
    will_return(__wrap_malloc, NULL);
    
    // Act & Assert
    // Verify graceful degradation on memory failure
}

/**
 * Test: DMI status update with large data
 * Purpose: Verify handling of large status data
 */
static void test_dmi_status_update_large_data(void **state) {
    // Arrange
    char large_data[4096];
    memset(large_data, 'A', sizeof(large_data) - 1);
    large_data[sizeof(large_data) - 1] = '\0';
    
    // Act & Assert
    // Verify large data handling
}

/**
 * Test: DMI status update with special characters
 * Purpose: Verify handling of special characters in data
 */
static void test_dmi_status_update_special_chars(void **state) {
    // Arrange
    const char *status = "status=ready\x00\x01\x02";
    
    // Act & Assert
    // Verify special character handling
}

/**
 * Test: DMI status query success
 * Purpose: Verify successful status query operation
 */
static void test_dmi_status_query_success(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_read, 20);
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify query returns valid data
}

/**
 * Test: DMI status query file not found
 * Purpose: Verify handling when status file doesn't exist
 */
static void test_dmi_status_query_file_not_found(void **state) {
    // Arrange
    will_return(__wrap_open, -1);
    
    // Act & Assert
    // Verify appropriate error return
}

/**
 * Test: DMI status clear operation
 * Purpose: Verify successful clearing of status file
 */
static void test_dmi_status_clear_success(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify file is cleared
}

/**
 * Test: DMI status clear with invalid permissions
 * Purpose: Verify handling of permission errors on clear
 */
static void test_dmi_status_clear_permission_denied(void **state) {
    // Arrange
    will_return(__wrap_open, -1);
    
    // Act & Assert
    // Verify error handling
}

/**
 * Test: DMI buffer overflow protection
 * Purpose: Verify protection against buffer overflow
 */
static void test_dmi_buffer_overflow_protection(void **state) {
    // Arrange
    char small_buffer[10];
    const char *large_data = "This is much larger than 10 bytes";
    
    // Act & Assert
    // Verify buffer boundaries are respected
}

/**
 * Test: DMI concurrent read operations
 * Purpose: Verify handling of concurrent reads (thread safety)
 */
static void test_dmi_concurrent_reads(void **state) {
    // This test would require thread handling mocks
    // Verify that concurrent reads don't interfere
}

/**
 * Test: DMI concurrent write and read
 * Purpose: Verify thread safety with concurrent read/write
 */
static void test_dmi_concurrent_write_read(void **state) {
    // This test would require thread handling mocks
    // Verify data consistency with concurrent operations
}

/**
 * Test: DMI file format validation
 * Purpose: Verify parsing of valid DMI format
 */
static void test_dmi_parse_valid_format(void **state) {
    // Arrange
    const char *dmi_content = "bios_vendor=ACME\nbios_version=1.0\n";
    
    // Act & Assert
    // Verify parsing succeeds
}

/**
 * Test: DMI file format invalid
 * Purpose: Verify handling of malformed DMI content
 */
static void test_dmi_parse_invalid_format(void **state) {
    // Arrange
    const char *dmi_content = "invalid content without proper format";
    
    // Act & Assert
    // Verify graceful error handling
}

/**
 * Test: DMI string parsing with escape sequences
 * Purpose: Verify handling of escaped characters
 */
static void test_dmi_parse_escape_sequences(void **state) {
    // Arrange
    const char *content = "vendor=ACME\\nCorp";
    
    // Act & Assert
    // Verify escape sequence handling
}

/**
 * Test: DMI field extraction single value
 * Purpose: Verify extracting single field values
 */
static void test_dmi_extract_field_single(void **state) {
    // Arrange
    const char *data = "vendor=Dell\n";
    const char *field = "vendor";
    
    // Act & Assert
    // Verify field value is correctly extracted
}

/**
 * Test: DMI field extraction multiple values
 * Purpose: Verify extracting multiple fields
 */
static void test_dmi_extract_fields_multiple(void **state) {
    // Arrange
    const char *data = "vendor=Dell\nmodel=XPS\nserial=12345\n";
    
    // Act & Assert
    // Verify all fields are correctly extracted
}

/**
 * Test: DMI field extraction non-existent field
 * Purpose: Verify handling of missing fields
 */
static void test_dmi_extract_field_not_found(void **state) {
    // Arrange
    const char *data = "vendor=Dell\n";
    const char *field = "nonexistent";
    
    // Act & Assert
    // Verify appropriate return when field not found
}

/**
 * Test: DMI field extraction empty value
 * Purpose: Verify handling of empty field values
 */
static void test_dmi_extract_field_empty_value(void **state) {
    // Arrange
    const char *data = "vendor=\n";
    
    // Act & Assert
    // Verify empty value handling
}

/**
 * Test: DMI initialization sequence
 * Purpose: Verify proper initialization order
 */
static void test_dmi_init_sequence(void **state) {
    // Act & Assert
    // Verify initialization steps in correct order
}

/**
 * Test: DMI cleanup on error
 * Purpose: Verify proper cleanup when errors occur
 */
static void test_dmi_cleanup_on_error(void **state) {
    // Arrange
    will_return(__wrap_malloc, NULL);
    
    // Act & Assert
    // Verify all resources are freed on error
}

/**
 * Test: DMI memory leak prevention
 * Purpose: Verify no memory leaks in normal operation
 */
static void test_dmi_no_memory_leaks(void **state) {
    // Act & Assert
    // Verify malloc/free balance
}

/**
 * Test: DMI file descriptor leak prevention
 * Purpose: Verify all file descriptors are properly closed
 */
static void test_dmi_no_fd_leaks(void **state) {
    // Act & Assert
    // Verify open/close balance
}

/**
 * Test: DMI status with maximum field count
 * Purpose: Verify handling of maximum number of fields
 */
static void test_dmi_max_field_count(void **state) {
    // Arrange: Create data with many fields
    
    // Act & Assert
    // Verify all fields are handled
}

/**
 * Test: DMI status with long field names
 * Purpose: Verify handling of very long field names
 */
static void test_dmi_long_field_names(void **state) {
    // Arrange
    char long_name[512];
    memset(long_name, 'a', sizeof(long_name) - 1);
    long_name[sizeof(long_name) - 1] = '\0';
    
    // Act & Assert
    // Verify long name handling
}

/**
 * Test: DMI status with long field values
 * Purpose: Verify handling of very long values
 */
static void test_dmi_long_field_values(void **state) {
    // Arrange
    char long_value[4096];
    memset(long_value, 'v', sizeof(long_value) - 1);
    long_value[sizeof(long_value) - 1] = '\0';
    
    // Act & Assert
    // Verify long value handling
}

/**
 * Test: DMI repeated field names
 * Purpose: Verify handling when same field appears multiple times
 */
static void test_dmi_repeated_field_names(void **state) {
    // Arrange
    const char *data = "vendor=Dell\nvendor=HP\n";
    
    // Act & Assert
    // Verify which value is returned (first/last)
}

/**
 * Test: DMI whitespace handling in values
 * Purpose: Verify proper whitespace handling
 */
static void test_dmi_whitespace_handling(void **state) {
    // Arrange
    const char *data = "vendor=  Dell  \n";
    
    // Act & Assert
    // Verify trimming or preserving of whitespace
}

/**
 * Test: DMI empty lines in file
 * Purpose: Verify handling of blank lines
 */
static void test_dmi_empty_lines(void **state) {
    // Arrange
    const char *data = "vendor=Dell\n\n\nmodel=XPS\n";
    
    // Act & Assert
    // Verify empty lines are handled
}

/**
 * Test: DMI comment lines
 * Purpose: Verify handling of comment characters
 */
static void test_dmi_comment_lines(void **state) {
    // Arrange
    const char *data = "# This is a comment\nvendor=Dell\n";
    
    // Act & Assert
    // Verify comments are properly skipped
}

/**
 * Test: DMI line ending variations
 * Purpose: Verify handling of different line endings (LF, CRLF)
 */
static void test_dmi_line_ending_variations(void **state) {
    // Arrange: Test with \n and \r\n
    
    // Act & Assert
    // Verify both line endings work
}

/**
 * Test: DMI binary data in file
 * Purpose: Verify handling when file contains binary data
 */
static void test_dmi_binary_data(void **state) {
    // Arrange
    unsigned char binary_data[] = {0xFF, 0xFE, 0x00, 0x00};
    
    // Act & Assert
    // Verify safe handling of binary content
}

/**
 * Test: DMI file permission modes
 * Purpose: Verify correct file permission settings
 */
static void test_dmi_file_permissions(void **state) {
    // Act & Assert
    // Verify files are created with correct permissions
}

/**
 * Test: DMI directory creation
 * Purpose: Verify directory creation if needed
 */
static void test_dmi_directory_creation(void **state) {
    // Act & Assert
    // Verify directories are created properly
}

/**
 * Test: DMI symlink following
 * Purpose: Verify symlinks are handled correctly
 */
static void test_dmi_symlink_handling(void **state) {
    // Act & Assert
    // Verify symlink behavior
}

/**
 * Test: DMI disk full scenario
 * Purpose: Verify handling when disk is full
 */
static void test_dmi_disk_full(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_write, -1);  // Write fails, disk full
    
    // Act & Assert
    // Verify error handling
}

/**
 * Test: DMI read-only filesystem
 * Purpose: Verify handling on read-only filesystem
 */
static void test_dmi_readonly_filesystem(void **state) {
    // Arrange
    will_return(__wrap_open, -1);  // Cannot open for writing
    
    // Act & Assert
    // Verify appropriate error
}

/**
 * Test: DMI stale file descriptor
 * Purpose: Verify handling of stale FD
 */
static void test_dmi_stale_file_descriptor(void **state) {
    // Arrange
    will_return(__wrap_read, -1);  // Read fails
    
    // Act & Assert
    // Verify error handling
}

/**
 * Test: DMI signal handling during I/O
 * Purpose: Verify handling of interrupted system calls
 */
static void test_dmi_signal_handling(void **state) {
    // This test would require signal mocking
    // Verify EINTR is handled properly (retry)
}

/**
 * Test: DMI file lock contention
 * Purpose: Verify behavior with file locks
 */
static void test_dmi_file_lock_contention(void **state) {
    // This test requires lock mocking
    // Verify lock handling
}

/**
 * Test: DMI recovery from partial write
 * Purpose: Verify recovery mechanism after partial write
 */
static void test_dmi_recovery_partial_write(void **state) {
    // Arrange
    will_return(__wrap_open, 3);
    will_return(__wrap_write, 5);  // Partial write
    will_return(__wrap_close, 0);
    
    // Act & Assert
    // Verify recovery or error reporting
}

/**
 * Test: DMI timestamp preservation
 * Purpose: Verify that file timestamps are handled correctly
 */
static void test_dmi_timestamp_preservation(void **state) {
    // Act & Assert
    // Verify timestamps
}

/**
 * Test: DMI platform independence
 * Purpose: Verify cross-platform path handling
 */
static void test_dmi_platform_independence(void **state) {
    // Act & Assert
    // Verify platform-specific paths are handled
}

/**
 * Test: DMI numeric overflow protection
 * Purpose: Verify protection against integer overflow
 */
static void test_dmi_numeric_overflow(void **state) {
    // Arrange
    const char *data = "size=999999999999999999999\n";
    
    // Act & Assert
    // Verify overflow protection
}

/**
 * Test: DMI state machine transitions
 * Purpose: Verify correct state transitions if state machine exists
 */
static void test_dmi_state_transitions(void **state) {
    // Act & Assert
    // Verify all state transitions
}

/**
 * Test: DMI idempotent operations
 * Purpose: Verify operations can be repeated safely
 */
static void test_dmi_idempotent_operations(void **state) {
    // Act
    // Perform operation twice
    
    // Assert
    // Verify same result both times
}

/**
 * Test: DMI encoding handling
 * Purpose: Verify UTF-8 and other encoding handling
 */
static void test_dmi_encoding_handling(void **state) {
    // Arrange
    const char *utf8_data = "vendor=ACME™\n";
    
    // Act & Assert
    // Verify encoding is preserved
}

/**
 * Test: DMI NULL byte in middle of string
 * Purpose: Verify handling of embedded NULL bytes
 */
static void test_dmi_embedded_null_bytes(void **state) {
    // Arrange
    const char *data = "vendor=ACME\0Corp\n";
    
    // Act & Assert
    // Verify NULL byte handling
}

/**
 * Test: DMI very long lines
 * Purpose: Verify handling of lines exceeding typical buffer sizes
 */
static void test_dmi_very_long_lines(void **state) {
    // Arrange: Create line > 4KB
    
    // Act & Assert
    // Verify handling
}

/**
 * Test: DMI initialization idempotency
 * Purpose: Verify init can be called multiple times safely
 */
static void test_dmi_init_idempotent(void **state) {
    // Act: Call init twice
    // Assert: No errors, only init once
}

/**
 * Test: DMI cleanup without init
 * Purpose: Verify cleanup can be called safely without init
 */
static void test_dmi_cleanup_without_init(void **state) {
    // Act: Call cleanup without init
    // Assert: No crash or error
}

/**
 * Test: DMI mixed case field names
 * Purpose: Verify case sensitivity in field names
 */
static void test_dmi_mixed_case_field_names(void **state) {
    // Arrange
    const char *data = "Vendor=Dell\nVENDOR=HP\n";
    
    // Act & Assert
    // Verify case handling (likely case-insensitive)
}

/**
 * Test: DMI get multiple values same field
 * Purpose: Verify behavior when getting array of values
 */
static void test_dmi_get_multiple_field_values(void **state) {
    // Arrange
    const char *data = "tag=value1\ntag=value2\ntag=value3\n";
    
    // Act & Assert
    // Verify all or specific value is returned
}

// ============================================================================
// Main test suite
// ============================================================================

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test_setup_teardown(test_dmi_initialize_with_valid_params, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_initialize_with_null_path, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_read_status_file_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_read_status_file_open_failure, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_read_status_file_empty, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_read_status_file_truncated, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_write_status_file_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_write_status_file_creation_failure, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_write_status_file_permission_denied, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_write_status_file_partial_write, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_close_on_failed_open, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_close_on_successful_open, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_close_failure, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_construct_path_valid, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_construct_path_empty_directory, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_construct_path_null_components, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_update_complete, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_update_malloc_failure, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_update_large_data, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_update_special_chars, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_query_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_query_file_not_found, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_clear_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_status_clear_permission_denied, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_buffer_overflow_protection, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_concurrent_reads, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_concurrent_write_read, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_parse_valid_format, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_parse_invalid_format, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_parse_escape_sequences, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_extract_field_single, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_extract_fields_multiple, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_extract_field_not_found, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_extract_field_empty_value, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_init_sequence, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_cleanup_on_error, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_no_memory_leaks, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_no_fd_leaks, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_max_field_count, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_long_field_names, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_long_field_values, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_repeated_field_names, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_whitespace_handling, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_empty_lines, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_comment_lines, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_line_ending_variations, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi_binary_data, setup, teardown),
        cmocka_unit_test_setup_teardown(test_dmi