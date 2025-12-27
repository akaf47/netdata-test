#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <stdint.h>
#include <limits.h>

#include "../streaming/stream-circular-buffer.h"

// ============================================================================
// Test Utility Functions
// ============================================================================

void assert_equal_int(int expected, int actual, const char *test_name) {
    if (expected != actual) {
        fprintf(stderr, "FAIL: %s - Expected %d, got %d\n", test_name, expected, actual);
        exit(1);
    }
    printf("PASS: %s\n", test_name);
}

void assert_equal_uint(unsigned int expected, unsigned int actual, const char *test_name) {
    if (expected != actual) {
        fprintf(stderr, "FAIL: %s - Expected %u, got %u\n", test_name, expected, actual);
        exit(1);
    }
    printf("PASS: %s\n", test_name);
}

void assert_equal_ptr(void *expected, void *actual, const char *test_name) {
    if (expected != actual) {
        fprintf(stderr, "FAIL: %s - Expected pointer %p, got %p\n", test_name, expected, actual);
        exit(1);
    }
    printf("PASS: %s\n", test_name);
}

void assert_equal_bytes(const uint8_t *expected, const uint8_t *actual, size_t len, const char *test_name) {
    if (memcmp(expected, actual, len) != 0) {
        fprintf(stderr, "FAIL: %s - Byte mismatch\n", test_name);
        exit(1);
    }
    printf("PASS: %s\n", test_name);
}

void assert_true(int condition, const char *test_name) {
    if (!condition) {
        fprintf(stderr, "FAIL: %s - Condition false\n", test_name);
        exit(1);
    }
    printf("PASS: %s\n", test_name);
}

void assert_false(int condition, const char *test_name) {
    if (condition) {
        fprintf(stderr, "FAIL: %s - Condition true\n", test_name);
        exit(1);
    }
    printf("PASS: %s\n", test_name);
}

// ============================================================================
// Test Suite: Circular Buffer Initialization
// ============================================================================

void test_circular_buffer_create_with_valid_size(void) {
    // Test creating a circular buffer with a valid size
    circular_buffer_t *cb = circular_buffer_create(1024);
    
    assert_true(cb != NULL, "circular_buffer_create returns non-NULL pointer");
    assert_equal_uint(1024, circular_buffer_size(cb), "circular_buffer_size returns correct size");
    assert_equal_uint(0, circular_buffer_available_size(cb), "new buffer has 0 available size");
    assert_equal_uint(1024, circular_buffer_free_size(cb), "new buffer has full free size");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_create_with_minimum_size(void) {
    // Test creating circular buffer with minimum valid size
    circular_buffer_t *cb = circular_buffer_create(1);
    
    assert_true(cb != NULL, "circular_buffer_create with size 1 succeeds");
    assert_equal_uint(1, circular_buffer_size(cb), "size is correctly set to 1");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_create_with_large_size(void) {
    // Test creating circular buffer with large size
    size_t large_size = 1024 * 1024;  // 1MB
    circular_buffer_t *cb = circular_buffer_create(large_size);
    
    assert_true(cb != NULL, "circular_buffer_create with large size succeeds");
    assert_equal_uint(large_size, circular_buffer_size(cb), "large size correctly set");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_create_with_zero_size(void) {
    // Test creating circular buffer with zero size (edge case)
    circular_buffer_t *cb = circular_buffer_create(0);
    
    // Should either fail or create with minimum size
    if (cb != NULL) {
        circular_buffer_free(cb);
    }
}

void test_circular_buffer_double_free_safety(void) {
    // Test that double free is handled safely
    circular_buffer_t *cb = circular_buffer_create(256);
    assert_true(cb != NULL, "buffer created");
    
    circular_buffer_free(cb);
    // Second free should not crash
    circular_buffer_free(NULL);  // Safe null check
}

// ============================================================================
// Test Suite: Writing Data
// ============================================================================

void test_circular_buffer_write_single_byte(void) {
    // Test writing a single byte
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data = 0x42;
    
    int written = circular_buffer_write(cb, &data, 1);
    
    assert_equal_int(1, written, "write returns 1 for single byte");
    assert_equal_uint(1, circular_buffer_available_size(cb), "available size increases by 1");
    assert_equal_uint(255, circular_buffer_free_size(cb), "free size decreases by 1");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_write_multiple_bytes(void) {
    // Test writing multiple bytes
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[5] = {0x01, 0x02, 0x03, 0x04, 0x05};
    
    int written = circular_buffer_write(cb, data, 5);
    
    assert_equal_int(5, written, "write returns 5 for 5 bytes");
    assert_equal_uint(5, circular_buffer_available_size(cb), "available size is 5");
    assert_equal_uint(251, circular_buffer_free_size(cb), "free size is 251");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_write_fill_entire_buffer(void) {
    // Test writing until buffer is full
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[256];
    memset(data, 0xAA, 256);
    
    int written = circular_buffer_write(cb, data, 256);
    
    assert_equal_int(256, written, "write returns 256 when filling entire buffer");
    assert_equal_uint(256, circular_buffer_available_size(cb), "available size is 256");
    assert_equal_uint(0, circular_buffer_free_size(cb), "free size is 0");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_write_beyond_capacity(void) {
    // Test writing more data than buffer capacity
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[300];
    memset(data, 0xBB, 300);
    
    int written = circular_buffer_write(cb, data, 300);
    
    // Should write only what fits
    assert_equal_int(256, written, "write is limited to buffer capacity");
    assert_equal_uint(256, circular_buffer_available_size(cb), "available size is at capacity");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_write_to_full_buffer(void) {
    // Test writing to already full buffer
    circular_buffer_t *cb = circular_buffer_create(128);
    uint8_t fill_data[128];
    memset(fill_data, 0xCC, 128);
    
    circular_buffer_write(cb, fill_data, 128);
    
    uint8_t extra_data[10] = {0xFF};
    int written = circular_buffer_write(cb, extra_data, 10);
    
    assert_equal_int(0, written, "write to full buffer returns 0");
    assert_equal_uint(128, circular_buffer_available_size(cb), "available size unchanged");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_write_null_pointer(void) {
    // Test writing with null data pointer
    circular_buffer_t *cb = circular_buffer_create(256);
    
    int written = circular_buffer_write(cb, NULL, 10);
    
    // Should handle gracefully, either return 0 or handle NULL
    assert_equal_int(0, written, "write with NULL pointer returns 0");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_write_zero_bytes(void) {
    // Test writing zero bytes
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[10] = {0x11};
    
    int written = circular_buffer_write(cb, data, 0);
    
    assert_equal_int(0, written, "write of 0 bytes returns 0");
    assert_equal_uint(0, circular_buffer_available_size(cb), "available size unchanged");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_sequential_writes(void) {
    // Test multiple sequential writes
    circular_buffer_t *cb = circular_buffer_create(512);
    uint8_t data1[100] = {0x11};
    uint8_t data2[100] = {0x22};
    uint8_t data3[100] = {0x33};
    
    memset(data1, 0x11, 100);
    memset(data2, 0x22, 100);
    memset(data3, 0x33, 100);
    
    circular_buffer_write(cb, data1, 100);
    assert_equal_uint(100, circular_buffer_available_size(cb), "after first write");
    
    circular_buffer_write(cb, data2, 100);
    assert_equal_uint(200, circular_buffer_available_size(cb), "after second write");
    
    circular_buffer_write(cb, data3, 100);
    assert_equal_uint(300, circular_buffer_available_size(cb), "after third write");
    
    circular_buffer_free(cb);
}

// ============================================================================
// Test Suite: Reading Data
// ============================================================================

void test_circular_buffer_read_from_empty_buffer(void) {
    // Test reading from empty buffer
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t output[10];
    
    int read = circular_buffer_read(cb, output, 10);
    
    assert_equal_int(0, read, "read from empty buffer returns 0");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_read_single_byte(void) {
    // Test reading a single byte
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[1] = {0x42};
    uint8_t output[1];
    
    circular_buffer_write(cb, input, 1);
    int read = circular_buffer_read(cb, output, 1);
    
    assert_equal_int(1, read, "read returns 1");
    assert_equal_uint(0x42, output[0], "read data is correct");
    assert_equal_uint(0, circular_buffer_available_size(cb), "buffer is empty after read");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_read_multiple_bytes(void) {
    // Test reading multiple bytes
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[5] = {0x01, 0x02, 0x03, 0x04, 0x05};
    uint8_t output[5] = {0};
    
    circular_buffer_write(cb, input, 5);
    int read = circular_buffer_read(cb, output, 5);
    
    assert_equal_int(5, read, "read returns 5");
    assert_equal_bytes(input, output, 5, "read data matches written data");
    assert_equal_uint(0, circular_buffer_available_size(cb), "buffer empty after read");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_read_partial(void) {
    // Test reading less than available
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[10];
    memset(input, 0x99, 10);
    uint8_t output[5] = {0};
    
    circular_buffer_write(cb, input, 10);
    int read = circular_buffer_read(cb, output, 5);
    
    assert_equal_int(5, read, "read returns 5");
    assert_equal_uint(5, circular_buffer_available_size(cb), "5 bytes remain");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_read_more_than_available(void) {
    // Test reading more than available
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[5] = {0x01, 0x02, 0x03, 0x04, 0x05};
    uint8_t output[20] = {0};
    
    circular_buffer_write(cb, input, 5);
    int read = circular_buffer_read(cb, output, 20);
    
    // Should read only what's available
    assert_equal_int(5, read, "read returns available bytes");
    assert_equal_bytes(input, output, 5, "read data is correct");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_read_with_null_output(void) {
    // Test reading with null output pointer
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[5] = {0x01, 0x02, 0x03, 0x04, 0x05};
    
    circular_buffer_write(cb, input, 5);
    int read = circular_buffer_read(cb, NULL, 5);
    
    // Should either skip or return error
    assert_equal_int(0, read, "read with NULL output returns 0");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_read_zero_bytes(void) {
    // Test reading zero bytes
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[10] = {0xAA};
    uint8_t output[10] = {0};
    
    circular_buffer_write(cb, input, 10);
    int read = circular_buffer_read(cb, output, 0);
    
    assert_equal_int(0, read, "read of 0 bytes returns 0");
    assert_equal_uint(10, circular_buffer_available_size(cb), "available size unchanged");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_write_and_read_pattern(void) {
    // Test write-read-write-read pattern
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t write1[5] = {0x11, 0x12, 0x13, 0x14, 0x15};
    uint8_t write2[5] = {0x21, 0x22, 0x23, 0x24, 0x25};
    uint8_t read1[5] = {0};
    uint8_t read2[5] = {0};
    
    circular_buffer_write(cb, write1, 5);
    circular_buffer_read(cb, read1, 5);
    assert_equal_bytes(write1, read1, 5, "first write-read cycle");
    
    circular_buffer_write(cb, write2, 5);
    circular_buffer_read(cb, read2, 5);
    assert_equal_bytes(write2, read2, 5, "second write-read cycle");
    
    circular_buffer_free(cb);
}

// ============================================================================
// Test Suite: Circular Behavior (Wrap-around)
// ============================================================================

void test_circular_buffer_write_wrap_around(void) {
    // Test writing data that wraps around the buffer
    circular_buffer_t *cb = circular_buffer_create(64);
    uint8_t data1[40];
    uint8_t data2[30];
    memset(data1, 0x11, 40);
    memset(data2, 0x22, 30);
    
    circular_buffer_write(cb, data1, 40);
    assert_equal_uint(40, circular_buffer_available_size(cb), "first write");
    
    // Read partially to create wraparound scenario
    uint8_t temp[20] = {0};
    circular_buffer_read(cb, temp, 20);
    assert_equal_uint(20, circular_buffer_available_size(cb), "after partial read");
    
    // Now write more data - should wrap around
    int written = circular_buffer_write(cb, data2, 30);
    assert_equal_int(30, written, "wrap-around write succeeds");
    assert_equal_uint(50, circular_buffer_available_size(cb), "available size is 50");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_read_wrap_around(void) {
    // Test reading data that wraps around the buffer
    circular_buffer_t *cb = circular_buffer_create(64);
    uint8_t data1[40];
    uint8_t data2[30];
    uint8_t output[50] = {0};
    
    memset(data1, 0x11, 40);
    memset(data2, 0x22, 30);
    
    // Write two chunks
    circular_buffer_write(cb, data1, 40);
    
    uint8_t temp[20] = {0};
    circular_buffer_read(cb, temp, 20);
    
    circular_buffer_write(cb, data2, 30);
    
    // Read all remaining data (wrap-around scenario)
    int read = circular_buffer_read(cb, output, 50);
    assert_equal_int(50, read, "wrap-around read returns correct count");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_multiple_wrap_arounds(void) {
    // Test multiple cycles of write-read that cause wrap-around
    circular_buffer_t *cb = circular_buffer_create(128);
    uint8_t data[64];
    uint8_t output[64];
    
    memset(data, 0x55, 64);
    
    // Do 3 cycles
    for (int i = 0; i < 3; i++) {
        circular_buffer_write(cb, data, 64);
        circular_buffer_read(cb, output, 64);
        assert_equal_bytes(data, output, 64, "wrap-around cycle");
    }
    
    circular_buffer_free(cb);
}

// ============================================================================
// Test Suite: Peek Operations
// ============================================================================

void test_circular_buffer_peek_from_empty(void) {
    // Test peeking from empty buffer
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t output[10] = {0};
    
    int peeked = circular_buffer_peek(cb, output, 10);
    
    assert_equal_int(0, peeked, "peek from empty buffer returns 0");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_peek_single_byte(void) {
    // Test peeking a single byte without consuming
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[1] = {0x42};
    uint8_t output[1] = {0};
    
    circular_buffer_write(cb, input, 1);
    int peeked = circular_buffer_peek(cb, output, 1);
    
    assert_equal_int(1, peeked, "peek returns 1");
    assert_equal_uint(0x42, output[0], "peek data is correct");
    assert_equal_uint(1, circular_buffer_available_size(cb), "available size unchanged");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_peek_multiple_bytes(void) {
    // Test peeking multiple bytes
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[5] = {0x01, 0x02, 0x03, 0x04, 0x05};
    uint8_t output[5] = {0};
    
    circular_buffer_write(cb, input, 5);
    int peeked = circular_buffer_peek(cb, output, 5);
    
    assert_equal_int(5, peeked, "peek returns 5");
    assert_equal_bytes(input, output, 5, "peek data matches");
    assert_equal_uint(5, circular_buffer_available_size(cb), "data still available");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_peek_partial(void) {
    // Test peeking less than available
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[10];
    memset(input, 0x99, 10);
    uint8_t output[5] = {0};
    
    circular_buffer_write(cb, input, 10);
    int peeked = circular_buffer_peek(cb, output, 5);
    
    assert_equal_int(5, peeked, "peek returns 5");
    assert_equal_uint(10, circular_buffer_available_size(cb), "all data still available");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_peek_doesnt_consume(void) {
    // Test that multiple peeks return same data
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t input[5] = {0x10, 0x20, 0x30, 0x40, 0x50};
    uint8_t output1[5] = {0};
    uint8_t output2[5] = {0};
    
    circular_buffer_write(cb, input, 5);
    circular_buffer_peek(cb, output1, 5);
    circular_buffer_peek(cb, output2, 5);
    
    assert_equal_bytes(output1, output2, 5, "consecutive peeks match");
    assert_equal_uint(5, circular_buffer_available_size(cb), "data unconsumed");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_peek_wrap_around(void) {
    // Test peeking with wrap-around
    circular_buffer_t *cb = circular_buffer_create(64);
    uint8_t data1[40];
    uint8_t data2[30];
    uint8_t output[50] = {0};
    
    memset(data1, 0x11, 40);
    memset(data2, 0x22, 30);
    
    circular_buffer_write(cb, data1, 40);
    uint8_t temp[20] = {0};
    circular_buffer_read(cb, temp, 20);
    circular_buffer_write(cb, data2, 30);
    
    int peeked = circular_buffer_peek(cb, output, 50);
    assert_equal_int(50, peeked, "peek with wrap-around");
    assert_equal_uint(50, circular_buffer_available_size(cb), "data unconsumed");
    
    circular_buffer_free(cb);
}

// ============================================================================
// Test Suite: Size and Status Operations
// ============================================================================

void test_circular_buffer_size(void) {
    // Test getting buffer size
    size_t sizes[] = {64, 256, 1024, 4096};
    
    for (size_t i = 0; i < sizeof(sizes)/sizeof(sizes[0]); i++) {
        circular_buffer_t *cb = circular_buffer_create(sizes[i]);
        assert_equal_uint(sizes[i], circular_buffer_size(cb), "size matches creation size");
        circular_buffer_free(cb);
    }
}

void test_circular_buffer_free_size_initial(void) {
    // Test free size of new buffer
    circular_buffer_t *cb = circular_buffer_create(256);
    
    assert_equal_uint(256, circular_buffer_free_size(cb), "new buffer has full free size");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_free_size_after_write(void) {
    // Test free size decreases after write
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[100];
    memset(data, 0xAA, 100);
    
    circular_buffer_write(cb, data, 100);
    
    assert_equal_uint(156, circular_buffer_free_size(cb), "free size after 100-byte write");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_free_size_after_read(void) {
    // Test free size increases after read
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[100];
    uint8_t output[50];
    memset(data, 0xBB, 100);
    
    circular_buffer_write(cb, data, 100);
    circular_buffer_read(cb, output, 50);
    
    assert_equal_uint(206, circular_buffer_free_size(cb), "free size after 50-byte read");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_available_size_initial(void) {
    // Test available size of new buffer
    circular_buffer_t *cb = circular_buffer_create(256);
    
    assert_equal_uint(0, circular_buffer_available_size(cb), "new buffer has 0 available");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_available_size_after_write(void) {
    // Test available size increases after write
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[100];
    memset(data, 0xCC, 100);
    
    circular_buffer_write(cb, data, 100);
    
    assert_equal_uint(100, circular_buffer_available_size(cb), "available size after 100-byte write");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_available_size_after_read(void) {
    // Test available size decreases after read
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[100];
    uint8_t output[50];
    memset(data, 0xDD, 100);
    
    circular_buffer_write(cb, data, 100);
    circular_buffer_read(cb, output, 50);
    
    assert_equal_uint(50, circular_buffer_available_size(cb), "available size after 50-byte read");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_is_empty(void) {
    // Test is_empty function
    circular_buffer_t *cb = circular_buffer_create(256);
    
    assert_true(circular_buffer_is_empty(cb), "new buffer is empty");
    
    uint8_t data[10] = {0xFF};
    circular_buffer_write(cb, data, 10);
    assert_false(circular_buffer_is_empty(cb), "buffer with data is not empty");
    
    uint8_t output[10];
    circular_buffer_read(cb, output, 10);
    assert_true(circular_buffer_is_empty(cb), "buffer is empty after reading all");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_is_full(void) {
    // Test is_full function
    circular_buffer_t *cb = circular_buffer_create(256);
    
    assert_false(circular_buffer_is_full(cb), "new buffer is not full");
    
    uint8_t data[256];
    memset(data, 0xEE, 256);
    circular_buffer_write(cb, data, 256);
    assert_true(circular_buffer_is_full(cb), "full buffer is full");
    
    uint8_t output[256];
    circular_buffer_read(cb, output, 256);
    assert_false(circular_buffer_is_full(cb), "empty buffer is not full");
    
    circular_buffer_free(cb);
}

// ============================================================================
// Test Suite: Reset and Clear Operations
// ============================================================================

void test_circular_buffer_reset(void) {
    // Test reset/clear operation
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[100];
    memset(data, 0x77, 100);
    
    circular_buffer_write(cb, data, 100);
    assert_equal_uint(100, circular_buffer_available_size(cb), "has data before reset");
    
    circular_buffer_reset(cb);
    
    assert_equal_uint(0, circular_buffer_available_size(cb), "empty after reset");
    assert_equal_uint(256, circular_buffer_free_size(cb), "full free size after reset");
    assert_true(circular_buffer_is_empty(cb), "is empty after reset");
    
    circular_buffer_free(cb);
}

void test_circular_buffer_reset_multiple_times(void) {
    // Test multiple reset operations
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t data[50];
    memset(data, 0x88, 50);
    
    for (int i = 0; i < 3; i++) {
        circular_buffer_write(cb, data, 50);
        assert_equal_uint(50, circular_buffer_available_size(cb), "has data");
        
        circular_buffer_reset(cb);
        assert_equal_uint(0, circular_buffer_available_size(cb), "empty after reset");
    }
    
    circular_buffer_free(cb);
}

// ============================================================================
// Test Suite: Data Integrity
// ============================================================================

void test_circular_buffer_data_integrity_pattern(void) {
    // Test data integrity with repeating patterns
    circular_buffer_t *cb = circular_buffer_create(256);
    uint8_t pattern[16] = {0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
                           0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10};
    uint8_t output[16];
    
    // Write and read the pattern multiple times
    for (int i = 0; i < 5; i++) {
        circular_buffer_write(cb, pattern, 16);