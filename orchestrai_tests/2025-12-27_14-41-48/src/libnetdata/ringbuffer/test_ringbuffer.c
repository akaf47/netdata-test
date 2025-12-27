#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "../ringbuffer/ringbuffer.h"

// Mock allocation tracking
static int mock_allocation_count = 0;
static int mock_allocation_fail_at = -1;

void* __wrap_malloc(size_t size) {
    if (mock_allocation_count == mock_allocation_fail_at) {
        return NULL;
    }
    mock_allocation_count++;
    return __real_malloc(size);
}

void __wrap_free(void* ptr) {
    __real_free(ptr);
}

// ============================================================================
// Test Fixtures and Setup
// ============================================================================

static void setup(void **state) {
    mock_allocation_count = 0;
    mock_allocation_fail_at = -1;
}

static void teardown(void **state) {
    mock_allocation_count = 0;
    mock_allocation_fail_at = -1;
}

// ============================================================================
// ringbuffer_create Tests
// ============================================================================

static void test_ringbuffer_create_success(void **state) {
    // Test: ringbuffer_create with valid size
    ringbuffer_t *rb = ringbuffer_create(100);
    
    assert_non_null(rb);
    assert_int_equal(rb->size, 100);
    assert_int_equal(rb->read_pos, 0);
    assert_int_equal(rb->write_pos, 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_create_zero_size(void **state) {
    // Test: ringbuffer_create with zero size - should create minimal buffer or fail gracefully
    ringbuffer_t *rb = ringbuffer_create(0);
    
    // Implementation dependent: either NULL or minimal buffer
    if (rb != NULL) {
        assert_int_equal(rb->size, 0);
        ringbuffer_free(rb);
    } else {
        assert_null(rb);
    }
}

static void test_ringbuffer_create_small_size(void **state) {
    // Test: ringbuffer_create with size 1
    ringbuffer_t *rb = ringbuffer_create(1);
    
    assert_non_null(rb);
    assert_int_equal(rb->size, 1);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_create_large_size(void **state) {
    // Test: ringbuffer_create with large size
    ringbuffer_t *rb = ringbuffer_create(1000000);
    
    assert_non_null(rb);
    assert_int_equal(rb->size, 1000000);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_create_allocation_failure(void **state) {
    // Test: ringbuffer_create when malloc fails
    mock_allocation_fail_at = 0;
    
    ringbuffer_t *rb = ringbuffer_create(100);
    
    assert_null(rb);
    
    mock_allocation_fail_at = -1;
}

// ============================================================================
// ringbuffer_free Tests
// ============================================================================

static void test_ringbuffer_free_valid(void **state) {
    // Test: ringbuffer_free with valid ringbuffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    ringbuffer_free(rb);
    // If we reach here, no segfault occurred
    assert_true(1);
}

static void test_ringbuffer_free_null(void **state) {
    // Test: ringbuffer_free with NULL pointer - should not crash
    ringbuffer_free(NULL);
    
    // If we reach here, no segfault occurred
    assert_true(1);
}

// ============================================================================
// ringbuffer_write Tests
// ============================================================================

static void test_ringbuffer_write_single_byte(void **state) {
    // Test: write single byte to empty ringbuffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t data = 0x42;
    int written = ringbuffer_write(rb, &data, 1);
    
    assert_int_equal(written, 1);
    assert_int_equal(rb->write_pos, 1);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_write_multiple_bytes(void **state) {
    // Test: write multiple bytes
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t data[] = {0x01, 0x02, 0x03, 0x04, 0x05};
    int written = ringbuffer_write(rb, data, 5);
    
    assert_int_equal(written, 5);
    assert_int_equal(rb->write_pos, 5);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_write_full_buffer(void **state) {
    // Test: write data that fills the buffer completely
    ringbuffer_t *rb = ringbuffer_create(10);
    assert_non_null(rb);
    
    uint8_t data[10];
    memset(data, 0xAA, 10);
    
    int written = ringbuffer_write(rb, data, 10);
    
    assert_int_equal(written, 10);
    assert_int_equal(rb->write_pos, 10);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_write_overflow(void **state) {
    // Test: write more data than buffer can hold
    ringbuffer_t *rb = ringbuffer_create(10);
    assert_non_null(rb);
    
    uint8_t data[15];
    memset(data, 0xBB, 15);
    
    int written = ringbuffer_write(rb, data, 15);
    
    // Should either write only what fits or wrap around
    assert_int_equal(written, 10);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_write_with_wraparound(void **state) {
    // Test: write that wraps around buffer boundary
    ringbuffer_t *rb = ringbuffer_create(10);
    assert_non_null(rb);
    
    // Write 7 bytes
    uint8_t data1[7];
    memset(data1, 0x11, 7);
    ringbuffer_write(rb, data1, 7);
    
    // Read 5 bytes to move read pointer
    uint8_t read_buf[5];
    ringbuffer_read(rb, read_buf, 5);
    
    // Write 8 bytes (should wrap)
    uint8_t data2[8];
    memset(data2, 0x22, 8);
    int written = ringbuffer_write(rb, data2, 8);
    
    assert_int_equal(written, 8);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_write_null_data(void **state) {
    // Test: write with NULL data pointer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    int written = ringbuffer_write(rb, NULL, 10);
    
    // Should handle gracefully - either return 0 or -1
    assert_true(written <= 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_write_to_null_ringbuffer(void **state) {
    // Test: write to NULL ringbuffer
    uint8_t data = 0x42;
    
    int written = ringbuffer_write(NULL, &data, 1);
    
    assert_true(written < 0 || written == 0);
}

static void test_ringbuffer_write_zero_length(void **state) {
    // Test: write zero bytes
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t data = 0x42;
    int written = ringbuffer_write(rb, &data, 0);
    
    assert_int_equal(written, 0);
    assert_int_equal(rb->write_pos, 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_write_sequential(void **state) {
    // Test: multiple sequential writes
    ringbuffer_t *rb = ringbuffer_create(20);
    assert_non_null(rb);
    
    uint8_t data1[5];
    memset(data1, 0x11, 5);
    int w1 = ringbuffer_write(rb, data1, 5);
    
    uint8_t data2[5];
    memset(data2, 0x22, 5);
    int w2 = ringbuffer_write(rb, data2, 5);
    
    uint8_t data3[5];
    memset(data3, 0x33, 5);
    int w3 = ringbuffer_write(rb, data3, 5);
    
    assert_int_equal(w1, 5);
    assert_int_equal(w2, 5);
    assert_int_equal(w3, 5);
    assert_int_equal(rb->write_pos, 15);
    
    ringbuffer_free(rb);
}

// ============================================================================
// ringbuffer_read Tests
// ============================================================================

static void test_ringbuffer_read_single_byte(void **state) {
    // Test: read single byte from ringbuffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data = 0x42;
    ringbuffer_write(rb, &write_data, 1);
    
    uint8_t read_data;
    int read = ringbuffer_read(rb, &read_data, 1);
    
    assert_int_equal(read, 1);
    assert_int_equal(read_data, 0x42);
    assert_int_equal(rb->read_pos, 1);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_multiple_bytes(void **state) {
    // Test: read multiple bytes
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data[] = {0x01, 0x02, 0x03, 0x04, 0x05};
    ringbuffer_write(rb, write_data, 5);
    
    uint8_t read_data[5];
    int read = ringbuffer_read(rb, read_data, 5);
    
    assert_int_equal(read, 5);
    assert_memory_equal(read_data, write_data, 5);
    assert_int_equal(rb->read_pos, 5);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_empty_buffer(void **state) {
    // Test: read from empty ringbuffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t read_data[10];
    int read = ringbuffer_read(rb, read_data, 10);
    
    // Should return 0 or error code
    assert_int_equal(read, 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_partial(void **state) {
    // Test: read less data than available
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data[] = {0x01, 0x02, 0x03, 0x04, 0x05};
    ringbuffer_write(rb, write_data, 5);
    
    uint8_t read_data[3];
    int read = ringbuffer_read(rb, read_data, 3);
    
    assert_int_equal(read, 3);
    assert_memory_equal(read_data, write_data, 3);
    assert_int_equal(rb->read_pos, 3);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_more_than_available(void **state) {
    // Test: request more data than available
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data[] = {0x01, 0x02, 0x03};
    ringbuffer_write(rb, write_data, 3);
    
    uint8_t read_data[10];
    int read = ringbuffer_read(rb, read_data, 10);
    
    // Should read only what's available
    assert_int_equal(read, 3);
    assert_memory_equal(read_data, write_data, 3);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_after_wraparound(void **state) {
    // Test: read data that was written across wraparound
    ringbuffer_t *rb = ringbuffer_create(10);
    assert_non_null(rb);
    
    // Write 8 bytes
    uint8_t data1[8];
    memset(data1, 0xAA, 8);
    ringbuffer_write(rb, data1, 8);
    
    // Read 6 bytes
    uint8_t read_buf[6];
    ringbuffer_read(rb, read_buf, 6);
    
    // Write 5 more bytes (wraps around)
    uint8_t data2[5];
    memset(data2, 0xBB, 5);
    ringbuffer_write(rb, data2, 5);
    
    // Read remaining 2 bytes from first write
    uint8_t read_buf2[2];
    int read = ringbuffer_read(rb, read_buf2, 2);
    
    assert_int_equal(read, 2);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_null_data(void **state) {
    // Test: read with NULL buffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data = 0x42;
    ringbuffer_write(rb, &write_data, 1);
    
    int read = ringbuffer_read(rb, NULL, 1);
    
    assert_true(read < 0 || read == 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_from_null_ringbuffer(void **state) {
    // Test: read from NULL ringbuffer
    uint8_t read_data;
    
    int read = ringbuffer_read(NULL, &read_data, 1);
    
    assert_true(read < 0 || read == 0);
}

static void test_ringbuffer_read_zero_length(void **state) {
    // Test: read zero bytes
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data = 0x42;
    ringbuffer_write(rb, &write_data, 1);
    
    uint8_t read_data;
    int read = ringbuffer_read(rb, &read_data, 0);
    
    assert_int_equal(read, 0);
    assert_int_equal(rb->read_pos, 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_read_sequential(void **state) {
    // Test: multiple sequential reads
    ringbuffer_t *rb = ringbuffer_create(20);
    assert_non_null(rb);
    
    uint8_t write_data[] = {0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A};
    ringbuffer_write(rb, write_data, 10);
    
    uint8_t read_buf1[3];
    int r1 = ringbuffer_read(rb, read_buf1, 3);
    
    uint8_t read_buf2[3];
    int r2 = ringbuffer_read(rb, read_buf2, 3);
    
    uint8_t read_buf3[4];
    int r3 = ringbuffer_read(rb, read_buf3, 4);
    
    assert_int_equal(r1, 3);
    assert_int_equal(r2, 3);
    assert_int_equal(r3, 4);
    assert_memory_equal(read_buf1, write_data, 3);
    assert_memory_equal(read_buf2, write_data + 3, 3);
    assert_memory_equal(read_buf3, write_data + 6, 4);
    
    ringbuffer_free(rb);
}

// ============================================================================
// ringbuffer_available Tests
// ============================================================================

static void test_ringbuffer_available_empty(void **state) {
    // Test: available on empty ringbuffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    int avail = ringbuffer_available(rb);
    
    assert_int_equal(avail, 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_available_partial(void **state) {
    // Test: available with partial data
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t data[30];
    memset(data, 0xCC, 30);
    ringbuffer_write(rb, data, 30);
    
    int avail = ringbuffer_available(rb);
    
    assert_int_equal(avail, 30);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_available_after_read(void **state) {
    // Test: available after reading some data
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t data[50];
    memset(data, 0xDD, 50);
    ringbuffer_write(rb, data, 50);
    
    uint8_t read_buf[20];
    ringbuffer_read(rb, read_buf, 20);
    
    int avail = ringbuffer_available(rb);
    
    assert_int_equal(avail, 30);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_available_full(void **state) {
    // Test: available when buffer is full
    ringbuffer_t *rb = ringbuffer_create(50);
    assert_non_null(rb);
    
    uint8_t data[50];
    memset(data, 0xEE, 50);
    ringbuffer_write(rb, data, 50);
    
    int avail = ringbuffer_available(rb);
    
    assert_int_equal(avail, 50);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_available_null(void **state) {
    // Test: available on NULL ringbuffer
    int avail = ringbuffer_available(NULL);
    
    assert_int_equal(avail, 0);
}

static void test_ringbuffer_available_after_wraparound(void **state) {
    // Test: available after data wraps around
    ringbuffer_t *rb = ringbuffer_create(10);
    assert_non_null(rb);
    
    // Write 8 bytes
    uint8_t data1[8];
    memset(data1, 0x11, 8);
    ringbuffer_write(rb, data1, 8);
    
    // Read 6 bytes
    uint8_t read_buf[6];
    ringbuffer_read(rb, read_buf, 6);
    
    // Write 5 more bytes
    uint8_t data2[5];
    memset(data2, 0x22, 5);
    ringbuffer_write(rb, data2, 5);
    
    int avail = ringbuffer_available(rb);
    
    assert_int_equal(avail, 7);
    
    ringbuffer_free(rb);
}

// ============================================================================
// ringbuffer_space Tests
// ============================================================================

static void test_ringbuffer_space_empty(void **state) {
    // Test: space in empty ringbuffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    int space = ringbuffer_space(rb);
    
    assert_int_equal(space, 100);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_space_partial(void **state) {
    // Test: space with partial data
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t data[30];
    memset(data, 0xFF, 30);
    ringbuffer_write(rb, data, 30);
    
    int space = ringbuffer_space(rb);
    
    assert_int_equal(space, 70);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_space_full(void **state) {
    // Test: space when buffer is full
    ringbuffer_t *rb = ringbuffer_create(50);
    assert_non_null(rb);
    
    uint8_t data[50];
    memset(data, 0x12, 50);
    ringbuffer_write(rb, data, 50);
    
    int space = ringbuffer_space(rb);
    
    assert_int_equal(space, 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_space_after_read(void **state) {
    // Test: space increases after reading
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t data[50];
    memset(data, 0x34, 50);
    ringbuffer_write(rb, data, 50);
    
    uint8_t read_buf[20];
    ringbuffer_read(rb, read_buf, 20);
    
    int space = ringbuffer_space(rb);
    
    assert_int_equal(space, 70);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_space_null(void **state) {
    // Test: space on NULL ringbuffer
    int space = ringbuffer_space(NULL);
    
    assert_int_equal(space, 0);
}

static void test_ringbuffer_space_after_wraparound(void **state) {
    // Test: space after data wraps around
    ringbuffer_t *rb = ringbuffer_create(10);
    assert_non_null(rb);
    
    // Write 8 bytes
    uint8_t data1[8];
    memset(data1, 0x11, 8);
    ringbuffer_write(rb, data1, 8);
    
    // Read 6 bytes
    uint8_t read_buf[6];
    ringbuffer_read(rb, read_buf, 6);
    
    // Write 5 more bytes
    uint8_t data2[5];
    memset(data2, 0x22, 5);
    ringbuffer_write(rb, data2, 5);
    
    int space = ringbuffer_space(rb);
    
    assert_int_equal(space, 3);
    
    ringbuffer_free(rb);
}

// ============================================================================
// ringbuffer_peek Tests (if available)
// ============================================================================

static void test_ringbuffer_peek_single_byte(void **state) {
    // Test: peek single byte without advancing read pointer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data = 0x42;
    ringbuffer_write(rb, &write_data, 1);
    
    uint8_t peek_data;
    int peeked = ringbuffer_peek(rb, &peek_data, 1);
    
    if (peeked == 1) {
        assert_int_equal(peek_data, 0x42);
        assert_int_equal(rb->read_pos, 0);  // Read pos shouldn't change
    }
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_peek_multiple_bytes(void **state) {
    // Test: peek multiple bytes
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t write_data[] = {0x01, 0x02, 0x03, 0x04, 0x05};
    ringbuffer_write(rb, write_data, 5);
    
    uint8_t peek_data[5];
    int peeked = ringbuffer_peek(rb, peek_data, 5);
    
    if (peeked == 5) {
        assert_memory_equal(peek_data, write_data, 5);
        assert_int_equal(rb->read_pos, 0);  // Read pos shouldn't change
    }
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_peek_empty(void **state) {
    // Test: peek on empty buffer
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t peek_data[10];
    int peeked = ringbuffer_peek(rb, peek_data, 10);
    
    assert_int_equal(peeked, 0);
    
    ringbuffer_free(rb);
}

// ============================================================================
// Integration Tests
// ============================================================================

static void test_ringbuffer_write_read_cycle(void **state) {
    // Test: complete write-read cycle
    ringbuffer_t *rb = ringbuffer_create(100);
    assert_non_null(rb);
    
    uint8_t original[] = {0x01, 0x02, 0x03, 0x04, 0x05};
    int written = ringbuffer_write(rb, original, 5);
    assert_int_equal(written, 5);
    
    int available = ringbuffer_available(rb);
    assert_int_equal(available, 5);
    
    uint8_t read_back[5];
    int read = ringbuffer_read(rb, read_back, 5);
    assert_int_equal(read, 5);
    
    assert_memory_equal(read_back, original, 5);
    
    int space = ringbuffer_space(rb);
    assert_int_equal(space, 100);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_multiple_write_read_cycles(void **state) {
    // Test: multiple write-read cycles
    ringbuffer_t *rb = ringbuffer_create(50);
    assert_non_null(rb);
    
    for (int i = 0; i < 3; i++) {
        uint8_t data[10];
        for (int j = 0; j < 10; j++) {
            data[j] = (i * 10) + j;
        }
        
        int written = ringbuffer_write(rb, data, 10);
        assert_int_equal(written, 10);
        
        uint8_t read_back[10];
        int read = ringbuffer_read(rb, read_back, 10);
        assert_int_equal(read, 10);
        
        assert_memory_equal(read_back, data, 10);
    }
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_mixed_partial_operations(void **state) {
    // Test: mixed partial read/write operations
    ringbuffer_t *rb = ringbuffer_create(20);
    assert_non_null(rb);
    
    // Write 15 bytes
    uint8_t data1[15];
    memset(data1, 0xAA, 15);
    ringbuffer_write(rb, data1, 15);
    assert_int_equal(ringbuffer_available(rb), 15);
    
    // Read 7 bytes
    uint8_t buf1[7];
    ringbuffer_read(rb, buf1, 7);
    assert_int_equal(ringbuffer_available(rb), 8);
    
    // Write 10 bytes (should wrap)
    uint8_t data2[10];
    memset(data2, 0xBB, 10);
    ringbuffer_write(rb, data2, 10);
    assert_int_equal(ringbuffer_available(rb), 18);
    
    // Read remaining 18 bytes
    uint8_t buf2[18];
    int read = ringbuffer_read(rb, buf2, 18);
    assert_int_equal(read, 18);
    assert_int_equal(ringbuffer_available(rb), 0);
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_stress_test(void **state) {
    // Test: stress test with many operations
    ringbuffer_t *rb = ringbuffer_create(256);
    assert_non_null(rb);
    
    for (int iteration = 0; iteration < 100; iteration++) {
        // Random writes
        uint8_t write_data[50];
        for (int i = 0; i < 50; i++) {
            write_data[i] = (iteration * 50 + i) % 256;
        }
        
        int space = ringbuffer_space(rb);
        int to_write = (space < 50) ? space : 50;
        
        if (to_write > 0) {
            int written = ringbuffer_write(rb, write_data, to_write);
            assert_int_equal(written, to_write);
        }
        
        // Random reads
        int available = ringbuffer_available(rb);
        int to_read = (available > 0) ? (1 + (available % 10)) : 0;
        
        if (to_read > 0) {
            uint8_t read_data[50];
            int read = ringbuffer_read(rb, read_data, to_read);
            assert_true(read >= 0 && read <= available);
        }
    }
    
    ringbuffer_free(rb);
}

static void test_ringbuffer_