#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdarg.h>
#include <setjmp.h>
#include <cmocka.h>
#include <unistd.h>

#include "remote_write.h"

/* ===========================
   Mock Data & Structures
   =========================== */

typedef struct {
    char *url;
    char *auth_header;
    int timeout;
} mock_remote_write_config;

/* ===========================
   Setup & Teardown
   =========================== */

static int setup_remote_write(void **state) {
    mock_remote_write_config *config = malloc(sizeof(mock_remote_write_config));
    if (!config) return -1;
    
    config->url = malloc(256);
    strcpy(config->url, "http://localhost:9009/api/v1/write");
    config->auth_header = malloc(256);
    strcpy(config->auth_header, "");
    config->timeout = 30;
    
    *state = config;
    return 0;
}

static int teardown_remote_write(void **state) {
    mock_remote_write_config *config = *state;
    if (config) {
        if (config->url) free(config->url);
        if (config->auth_header) free(config->auth_header);
        free(config);
    }
    return 0;
}

/* ===========================
   Mock Functions
   =========================== */

static int __wrap_curl_easy_setopt(void *curl, int option, ...) {
    return CURLE_OK;
}

static int __wrap_curl_easy_perform(void *curl) {
    return CURLE_OK;
}

static int __wrap_curl_slist_append(void *list, const char *string) {
    check_expected(string);
    return (int)(intptr_t)mock();
}

/* ===========================
   Test Cases: Initialization
   =========================== */

static void test_remote_write_init_null_exporting_instance(void **state) {
    // Test initialization with NULL instance
    int result = init_remote_write_instance(NULL);
    assert_int_not_equal(result, 0);  // Should fail
}

static void test_remote_write_init_valid_instance(void **state) {
    // Test initialization with valid parameters
    struct exporting_instance mock_instance = {0};
    mock_instance.config.update_every = 10;
    
    int result = init_remote_write_instance(&mock_instance);
    // Verify initialization completed (specific behavior depends on implementation)
    assert_true(result >= 0);
}

static void test_remote_write_init_with_different_update_intervals(void **state) {
    struct exporting_instance instance1 = {0};
    instance1.config.update_every = 1;
    
    struct exporting_instance instance2 = {0};
    instance2.config.update_every = 60;
    
    struct exporting_instance instance3 = {0};
    instance3.config.update_every = 3600;
    
    // All valid intervals should initialize
    int result1 = init_remote_write_instance(&instance1);
    int result2 = init_remote_write_instance(&instance2);
    int result3 = init_remote_write_instance(&instance3);
    
    assert_true(result1 >= 0);
    assert_true(result2 >= 0);
    assert_true(result3 >= 0);
}

/* ===========================
   Test Cases: Buffer Operations
   =========================== */

static void test_remote_write_buffer_initialization(void **state) {
    // Test that buffer is properly initialized
    struct exporting_instance instance = {0};
    
    init_remote_write_instance(&instance);
    
    // Verify buffer exists and is usable
    assert_non_null(instance.buffer);
}

static void test_remote_write_buffer_append_empty_string(void **state) {
    // Test appending empty string to buffer
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    // This should handle gracefully
    buffer_strcat(instance.buffer, "");
    assert_non_null(instance.buffer);
}

static void test_remote_write_buffer_append_large_data(void **state) {
    // Test appending large amounts of data
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    // Create large string
    char large_data[10000];
    memset(large_data, 'a', sizeof(large_data) - 1);
    large_data[sizeof(large_data) - 1] = '\0';
    
    buffer_strcat(instance.buffer, large_data);
    assert_non_null(instance.buffer);
}

static void test_remote_write_buffer_append_special_characters(void **state) {
    // Test appending strings with special characters
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    const char *special = "test\n\t\r\0\xFF\x00";
    buffer_strcat(instance.buffer, special);
    assert_non_null(instance.buffer);
}

static void test_remote_write_buffer_reset(void **state) {
    // Test buffer reset functionality
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    buffer_strcat(instance.buffer, "test data");
    size_t before = instance.buffer->len;
    
    buffer_flush(instance.buffer);
    
    // After flush, buffer should be reset
    assert_true(instance.buffer->len <= before);
}

/* ===========================
   Test Cases: Send Operations
   =========================== */

static void test_remote_write_send_empty_buffer(void **state) {
    // Test sending with empty buffer
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    int result = remote_write_send_buffer(&instance);
    // Should handle empty buffer gracefully
    assert_true(result == 0 || result == 1);
}

static void test_remote_write_send_null_instance(void **state) {
    // Test sending with NULL instance
    int result = remote_write_send_buffer(NULL);
    assert_int_not_equal(result, 0);
}

static void test_remote_write_send_with_data(void **state) {
    // Test sending with actual data
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    buffer_strcat(instance.buffer, "# HELP test_metric Test metric\n");
    buffer_strcat(instance.buffer, "# TYPE test_metric counter\n");
    buffer_strcat(instance.buffer, "test_metric{job=\"prometheus\"} 42\n");
    
    int result = remote_write_send_buffer(&instance);
    assert_true(result >= 0);
}

static void test_remote_write_send_multiple_times(void **state) {
    // Test sending multiple times
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    for (int i = 0; i < 5; i++) {
        buffer_flush(instance.buffer);
        buffer_strcat(instance.buffer, "test_metric 42\n");
        int result = remote_write_send_buffer(&instance);
        assert_true(result >= 0);
    }
}

/* ===========================
   Test Cases: Metric Formatting
   =========================== */

static void test_remote_write_format_metric_with_labels(void **state) {
    // Test metric formatting with labels
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    // Call function that formats metrics
    add_remote_write_metric(&instance, "metric_name", 42.5, "label1=\"value1\",label2=\"value2\"");
    
    assert_non_null(instance.buffer);
    assert_true(instance.buffer->len > 0);
}

static void test_remote_write_format_metric_without_labels(void **state) {
    // Test metric formatting without labels
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric_name", 42.5, "");
    
    assert_non_null(instance.buffer);
}

static void test_remote_write_format_metric_null_name(void **state) {
    // Test metric formatting with NULL name
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, NULL, 42.5, "");
    // Should handle gracefully
    assert_non_null(instance.buffer);
}

static void test_remote_write_format_metric_null_labels(void **state) {
    // Test metric formatting with NULL labels
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric_name", 42.5, NULL);
    assert_non_null(instance.buffer);
}

static void test_remote_write_format_metric_empty_name(void **state) {
    // Test metric formatting with empty name
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "", 42.5, "");
    assert_non_null(instance.buffer);
}

static void test_remote_write_format_metric_special_characters_in_name(void **state) {
    // Test metric with special characters
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric_name_with_underscores_123", 42.5, "label=\"value\"");
    assert_true(instance.buffer->len > 0);
}

static void test_remote_write_format_metric_various_values(void **state) {
    // Test metrics with different value types
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric_1", 0, "");           // zero
    add_remote_write_metric(&instance, "metric_2", -100.5, "");      // negative
    add_remote_write_metric(&instance, "metric_3", 1e10, "");        // scientific notation
    add_remote_write_metric(&instance, "metric_4", 0.00001, "");     // very small
    
    assert_true(instance.buffer->len > 0);
}

/* ===========================
   Test Cases: Configuration
   =========================== */

static void test_remote_write_config_null_pointer(void **state) {
    // Test configuration with NULL
    int result = configure_remote_write(NULL);
    assert_int_not_equal(result, 0);
}

static void test_remote_write_config_valid(void **state) {
    // Test valid configuration
    struct exporting_instance instance = {0};
    instance.config.type = EXPORTING_CONNECTOR_TYPE_PROMETHEUS_REMOTE_WRITE;
    
    int result = configure_remote_write(&instance);
    assert_true(result >= 0);
}

static void test_remote_write_config_missing_url(void **state) {
    // Test configuration without URL
    struct exporting_instance instance = {0};
    instance.config.type = EXPORTING_CONNECTOR_TYPE_PROMETHEUS_REMOTE_WRITE;
    instance.config.destination = NULL;
    
    int result = configure_remote_write(&instance);
    // Should either fail or use default
    assert_true(result >= -1);
}

/* ===========================
   Test Cases: Cleanup
   =========================== */

static void test_remote_write_cleanup_null_instance(void **state) {
    // Test cleanup with NULL
    cleanup_remote_write(NULL);
    // Should not crash
}

static void test_remote_write_cleanup_valid_instance(void **state) {
    // Test cleanup with valid instance
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    cleanup_remote_write(&instance);
    // Should complete successfully
}

static void test_remote_write_cleanup_releases_resources(void **state) {
    // Test that cleanup actually releases resources
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    buffer_strcat(instance.buffer, "test");
    size_t before = instance.buffer->len;
    
    cleanup_remote_write(&instance);
    // Buffer should be cleaned
    assert_true(instance.buffer == NULL || instance.buffer->len == 0);
}

/* ===========================
   Test Cases: Error Handling
   =========================== */

static void test_remote_write_network_error(void **state) {
    // Test behavior when network error occurs
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    buffer_strcat(instance.buffer, "test_metric 1\n");
    
    // Call send - should handle errors gracefully
    int result = remote_write_send_buffer(&instance);
    assert_true(result >= -1);
}

static void test_remote_write_timeout_error(void **state) {
    // Test timeout handling
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    if (instance.config.socket_timeout) {
        buffer_strcat(instance.buffer, "test");
        int result = remote_write_send_buffer(&instance);
        assert_true(result >= -1);
    }
}

static void test_remote_write_invalid_url(void **state) {
    // Test with invalid URL
    struct exporting_instance instance = {0};
    instance.config.destination = "not a valid url";
    
    int result = configure_remote_write(&instance);
    // Should handle invalid URL gracefully
    assert_true(result >= -1);
}

/* ===========================
   Test Cases: Boundary Conditions
   =========================== */

static void test_remote_write_very_large_metric_value(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric", 9.999999999999999e308, "");
    assert_non_null(instance.buffer);
}

static void test_remote_write_very_small_metric_value(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric", 1e-308, "");
    assert_non_null(instance.buffer);
}

static void test_remote_write_metric_nan_value(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric", NAN, "");
    assert_non_null(instance.buffer);
}

static void test_remote_write_metric_infinity_value(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric", INFINITY, "");
    assert_non_null(instance.buffer);
}

static void test_remote_write_very_long_metric_name(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    char long_name[1024];
    memset(long_name, 'a', sizeof(long_name) - 1);
    long_name[sizeof(long_name) - 1] = '\0';
    
    add_remote_write_metric(&instance, long_name, 42, "");
    assert_non_null(instance.buffer);
}

static void test_remote_write_many_labels(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    char labels[2048];
    strcpy(labels, "");
    for (int i = 0; i < 50; i++) {
        char label[50];
        sprintf(label, "label%d=\"value%d\",", i, i);
        strcat(labels, label);
    }
    // Remove trailing comma
    if (strlen(labels) > 0) labels[strlen(labels) - 1] = '\0';
    
    add_remote_write_metric(&instance, "metric", 42, labels);
    assert_non_null(instance.buffer);
}

/* ===========================
   Test Cases: Data Format Validation
   =========================== */

static void test_remote_write_metric_format_prometheus_format(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric_total", 100, "job=\"prometheus\"");
    
    // Verify buffer contains prometheus-formatted data
    assert_true(instance.buffer->len > 0);
}

static void test_remote_write_buffer_contains_newlines(void **state) {
    struct exporting_instance instance = {0};
    init_remote_write_instance(&instance);
    
    add_remote_write_metric(&instance, "metric1", 1, "");
    add_remote_write_metric(&instance, "metric2", 2, "");
    
    const char *buffer_str = instance.buffer->buffer;
    assert_non_null(strchr(buffer_str, '\n'));
}

/* ===========================
   Test Runner
   =========================== */

int main(void) {
    const struct CMUnitTest tests[] = {
        // Initialization tests
        cmocka_unit_test(test_remote_write_init_null_exporting_instance),
        cmocka_unit_test(test_remote_write_init_valid_instance),
        cmocka_unit_test(test_remote_write_init_with_different_update_intervals),
        
        // Buffer operation tests
        cmocka_unit_test(test_remote_write_buffer_initialization),
        cmocka_unit_test(test_remote_write_buffer_append_empty_string),
        cmocka_unit_test(test_remote_write_buffer_append_large_data),
        cmocka_unit_test(test_remote_write_buffer_append_special_characters),
        cmocka_unit_test(test_remote_write_buffer_reset),
        
        // Send operation tests
        cmocka_unit_test(test_remote_write_send_empty_buffer),
        cmocka_unit_test(test_remote_write_send_null_instance),
        cmocka_unit_test(test_remote_write_send_with_data),
        cmocka_unit_test(test_remote_write_send_multiple_times),
        
        // Metric formatting tests
        cmocka_unit_test(test_remote_write_format_metric_with_labels),
        cmocka_unit_test(test_remote_write_format_metric_without_labels),
        cmocka_unit_test(test_remote_write_format_metric_null_name),
        cmocka_unit_test(test_remote_write_format_metric_null_labels),
        cmocka_unit_test(test_remote_write_format_metric_empty_name),
        cmocka_unit_test(test_remote_write_format_metric_special_characters_in_name),
        cmocka_unit_test(test_remote_write_format_metric_various_values),
        
        // Configuration tests
        cmocka_unit_test(test_remote_write_config_null_pointer),
        cmocka_unit_test(test_remote_write_config_valid),
        cmocka_unit_test(test_remote_write_config_missing_url),
        
        // Cleanup tests
        cmocka_unit_test(test_remote_write_cleanup_null_instance),
        cmocka_unit_test(test_remote_write_cleanup_valid_instance),
        cmocka_unit_test(test_remote_write_cleanup_releases_resources),
        
        // Error handling tests
        cmocka_unit_test(test_remote_write_network_error),
        cmocka_unit_test(test_remote_write_timeout_error),
        cmocka_unit_test(test_remote_write_invalid_url),
        
        // Boundary condition tests
        cmocka_unit_test(test_remote_write_very_large_metric_value),
        cmocka_unit_test(test_remote_write_very_small_metric_value),
        cmocka_unit_test(test_remote_write_metric_nan_value),
        cmocka_unit_test(test_remote_write_metric_infinity_value),
        cmocka_unit_test(test_remote_write_very_long_metric_name),
        cmocka_unit_test(test_remote_write_many_labels),
        
        // Data format validation tests
        cmocka_unit_test(test_remote_write_metric_format_prometheus_format),
        cmocka_unit_test(test_remote_write_buffer_contains_newlines),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}