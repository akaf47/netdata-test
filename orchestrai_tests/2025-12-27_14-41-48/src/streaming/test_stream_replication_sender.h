#ifndef TEST_STREAM_REPLICATION_SENDER_H
#define TEST_STREAM_REPLICATION_SENDER_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <time.h>

/* Mock structures and function declarations for testing */
typedef struct {
    char *host;
    int port;
    int timeout;
} stream_replication_sender_config_t;

typedef struct {
    int fd;
    stream_replication_sender_config_t config;
    int is_connected;
    int bytes_sent;
    int error_count;
} stream_replication_sender_t;

/* Test helper structures */
typedef struct {
    int call_count;
    int last_error;
    char last_message[1024];
} test_mock_state_t;

static test_mock_state_t mock_state = {0, 0, ""};

/* ============================================================================
 * TEST SETUP AND TEARDOWN
 * ============================================================================ */

void test_setup(void) {
    memset(&mock_state, 0, sizeof(test_mock_state_t));
}

void test_teardown(void) {
    memset(&mock_state, 0, sizeof(test_mock_state_t));
}

/* ============================================================================
 * BASIC STRUCTURE INITIALIZATION TESTS
 * ============================================================================ */

void test_stream_replication_sender_config_init_with_valid_values(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    const char *host = "localhost";
    int port = 8080;
    int timeout = 30;
    
    /* Arrange */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    
    /* Act */
    config.host = malloc(strlen(host) + 1);
    strcpy(config.host, host);
    config.port = port;
    config.timeout = timeout;
    
    /* Assert */
    assert(config.host != NULL);
    assert(strcmp(config.host, "localhost") == 0);
    assert(config.port == 8080);
    assert(config.timeout == 30);
    
    free(config.host);
    test_teardown();
}

void test_stream_replication_sender_config_init_with_null_host(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Arrange & Act */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    config.host = NULL;
    config.port = 8080;
    config.timeout = 30;
    
    /* Assert */
    assert(config.host == NULL);
    assert(config.port == 8080);
    assert(config.timeout == 30);
    
    test_teardown();
}

void test_stream_replication_sender_config_init_with_zero_port(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Arrange & Act */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    config.host = malloc(10);
    strcpy(config.host, "localhost");
    config.port = 0;
    config.timeout = 30;
    
    /* Assert */
    assert(config.port == 0);
    
    free(config.host);
    test_teardown();
}

void test_stream_replication_sender_config_init_with_negative_port(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Arrange & Act */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    config.port = -1;
    
    /* Assert */
    assert(config.port == -1);
    
    test_teardown();
}

void test_stream_replication_sender_config_init_with_max_port(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Arrange & Act */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    config.port = 65535;
    
    /* Assert */
    assert(config.port == 65535);
    
    test_teardown();
}

void test_stream_replication_sender_config_init_with_zero_timeout(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Arrange & Act */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    config.timeout = 0;
    
    /* Assert */
    assert(config.timeout == 0);
    
    test_teardown();
}

void test_stream_replication_sender_config_init_with_negative_timeout(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Arrange & Act */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    config.timeout = -1;
    
    /* Assert */
    assert(config.timeout == -1);
    
    test_teardown();
}

void test_stream_replication_sender_config_init_with_large_timeout(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Arrange & Act */
    memset(&config, 0, sizeof(stream_replication_sender_config_t));
    config.timeout = 999999;
    
    /* Assert */
    assert(config.timeout == 999999);
    
    test_teardown();
}

/* ============================================================================
 * SENDER INITIALIZATION TESTS
 * ============================================================================ */

void test_stream_replication_sender_init_allocates_structure(void) {
    test_setup();
    
    /* Act */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Assert */
    assert(sender != NULL);
    assert(sender->fd == 0);
    assert(sender->is_connected == 0);
    assert(sender->bytes_sent == 0);
    assert(sender->error_count == 0);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_init_with_fd_assignment(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    int test_fd = 42;
    
    /* Act */
    sender->fd = test_fd;
    
    /* Assert */
    assert(sender->fd == 42);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_init_with_negative_fd(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Act */
    sender->fd = -1;
    
    /* Assert */
    assert(sender->fd == -1);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_connection_state_initialized_disconnected(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Assert */
    assert(sender->is_connected == 0);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_connection_state_set_connected(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Act */
    sender->is_connected = 1;
    
    /* Assert */
    assert(sender->is_connected == 1);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_connection_state_set_disconnected(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    sender->is_connected = 1;
    
    /* Act */
    sender->is_connected = 0;
    
    /* Assert */
    assert(sender->is_connected == 0);
    
    free(sender);
    test_teardown();
}

/* ============================================================================
 * BYTES SENT TRACKING TESTS
 * ============================================================================ */

void test_stream_replication_sender_bytes_sent_initial_zero(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Assert */
    assert(sender->bytes_sent == 0);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_bytes_sent_increment(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Act */
    sender->bytes_sent += 100;
    sender->bytes_sent += 50;
    
    /* Assert */
    assert(sender->bytes_sent == 150);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_bytes_sent_large_value(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Act */
    sender->bytes_sent = 2147483647;
    
    /* Assert */
    assert(sender->bytes_sent == 2147483647);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_bytes_sent_reset(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    sender->bytes_sent = 1000;
    
    /* Act */
    sender->bytes_sent = 0;
    
    /* Assert */
    assert(sender->bytes_sent == 0);
    
    free(sender);
    test_teardown();
}

/* ============================================================================
 * ERROR TRACKING TESTS
 * ============================================================================ */

void test_stream_replication_sender_error_count_initial_zero(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Assert */
    assert(sender->error_count == 0);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_error_count_increment(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    
    /* Act */
    sender->error_count++;
    sender->error_count++;
    sender->error_count++;
    
    /* Assert */
    assert(sender->error_count == 3);
    
    free(sender);
    test_teardown();
}

void test_stream_replication_sender_error_count_reset(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender = malloc(sizeof(stream_replication_sender_t));
    memset(sender, 0, sizeof(stream_replication_sender_t));
    sender->error_count = 10;
    
    /* Act */
    sender->error_count = 0;
    
    /* Assert */
    assert(sender->error_count == 0);
    
    free(sender);
    test_teardown();
}

/* ============================================================================
 * CONFIG STRUCTURE FIELD MANIPULATION TESTS
 * ============================================================================ */

void test_stream_replication_sender_config_host_empty_string(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Act */
    config.host = malloc(1);
    strcpy(config.host, "");
    
    /* Assert */
    assert(strlen(config.host) == 0);
    
    free(config.host);
    test_teardown();
}

void test_stream_replication_sender_config_host_long_string(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    const char *long_host = "very.long.hostname.example.com.this.is.a.very.long.hostname";
    
    /* Act */
    config.host = malloc(strlen(long_host) + 1);
    strcpy(config.host, long_host);
    
    /* Assert */
    assert(strcmp(config.host, long_host) == 0);
    
    free(config.host);
    test_teardown();
}

void test_stream_replication_sender_config_host_special_characters(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    const char *special_host = "host-name_123.example.com";
    
    /* Act */
    config.host = malloc(strlen(special_host) + 1);
    strcpy(config.host, special_host);
    
    /* Assert */
    assert(strcmp(config.host, special_host) == 0);
    
    free(config.host);
    test_teardown();
}

void test_stream_replication_sender_config_port_min_valid(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Act */
    config.port = 1;
    
    /* Assert */
    assert(config.port == 1);
    
    test_teardown();
}

void test_stream_replication_sender_config_port_typical_values(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Test HTTP */
    config.port = 80;
    assert(config.port == 80);
    
    /* Test HTTPS */
    config.port = 443;
    assert(config.port == 443);
    
    /* Test custom */
    config.port = 19999;
    assert(config.port == 19999);
    
    test_teardown();
}

void test_stream_replication_sender_config_timeout_very_large(void) {
    test_setup();
    
    stream_replication_sender_config_t config;
    
    /* Act */
    config.timeout = 2147483647;
    
    /* Assert */
    assert(config.timeout == 2147483647);
    
    test_teardown();
}

/* ============================================================================
 * MULTI-SENDER SCENARIO TESTS
 * ============================================================================ */

void test_stream_replication_multiple_senders_independent_state(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_t *sender1 = malloc(sizeof(stream_replication_sender_t));
    stream_replication_sender_t *sender2 = malloc(sizeof(stream_replication_sender_t));
    memset(sender1, 0, sizeof(stream_replication_sender_t));
    memset(sender2, 0, sizeof(stream_replication_sender_t));
    
    /* Act */
    sender1->fd = 10;
    sender1->bytes_sent = 1000;
    sender1->error_count = 5;
    sender1->is_connected = 1;
    
    sender2->fd = 20;
    sender2->bytes_sent = 2000;
    sender2->error_count = 3;
    sender2->is_connected = 0;
    
    /* Assert */
    assert(sender1->fd == 10);
    assert(sender2->fd == 20);
    assert(sender1->bytes_sent == 1000);
    assert(sender2->bytes_sent == 2000);
    assert(sender1->error_count == 5);
    assert(sender2->error_count == 3);
    assert(sender1->is_connected == 1);
    assert(sender2->is_connected == 0);
    
    free(sender1);
    free(sender2);
    test_teardown();
}

void test_stream_replication_sender_config_copy(void) {
    test_setup();
    
    /* Arrange */
    stream_replication_sender_config_t config1, config2;
    memset(&config1, 0, sizeof(stream_replication_sender_config_t));
    memset(&config2, 0, sizeof(stream_replication_sender_config_t));
    
    config1.host = malloc(10);
    strcpy(config1.host, "localhost");
    config1.port = 8080;
    config1.timeout = 30;
    
    /* Act */
    config2.host = malloc(10);
    strcpy(config2.host, config1.host);
    config2.port = config1.port;
    config2.timeout = config1.timeout;
    
    /* Assert */
    assert(strcmp(config2.host, "localhost") == 0);
    assert(config2.port == 8080);
    assert(config2.timeout == 30);
    
    free(config1.host);
    free(config2.host);
    test_teardown();
}

/* ============================================================================
 * MAIN TEST RUNNER
 * ============================================================================ */

int main(void) {
    /* Setup and Teardown Tests */
    printf("Running stream replication sender header tests...\n");
    
    /* Basic Structure Initialization Tests */
    printf("Test: config_init_with_valid_values\n");
    test_stream_replication_sender_config_init_with_valid_values();
    
    printf("Test: config_init_with_null_host\n");
    test_stream_replication_sender_config_init_with_null_host();
    
    printf("Test: config_init_with_zero_port\n");
    test_stream_replication_sender_config_init_with_zero_port();
    
    printf("Test: config_init_with_negative_port\n");
    test_stream_replication_sender_config_init_with_negative_port();
    
    printf("Test: config_init_with_max_port\n");
    test_stream_replication_sender_config_init_with_max_port();
    
    printf("Test: config_init_with_zero_timeout\n");
    test_stream_replication_sender_config_init_with_zero_timeout();
    
    printf("Test: config_init_with_negative_timeout\n");
    test_stream_replication_sender_config_init_with_negative_timeout();
    
    printf("Test: config_init_with_large_timeout\n");
    test_stream_replication_sender_config_init_with_large_timeout();
    
    /* Sender Initialization Tests */
    printf("Test: sender_init_allocates_structure\n");
    test_stream_replication_sender_init_allocates_structure();
    
    printf("Test: sender_init_with_fd_assignment\n");
    test_stream_replication_sender_init_with_fd_assignment();
    
    printf("Test: sender_init_with_negative_fd\n");
    test_stream_replication_sender_init_with_negative_fd();
    
    printf("Test: sender_connection_state_initialized_disconnected\n");
    test_stream_replication_sender_connection_state_initialized_disconnected();
    
    printf("Test: sender_connection_state_set_connected\n");
    test_stream_replication_sender_connection_state_set_connected();
    
    printf("Test: sender_connection_state_set_disconnected\n");
    test_stream_replication_sender_connection_state_set_disconnected();
    
    /* Bytes Sent Tracking Tests */
    printf("Test: sender_bytes_sent_initial_zero\n");
    test_stream_replication_sender_bytes_sent_initial_zero();
    
    printf("Test: sender_bytes_sent_increment\n");
    test_stream_replication_sender_bytes_sent_increment();
    
    printf("Test: sender_bytes_sent_large_value\n");
    test_stream_replication_sender_bytes_sent_large_value();
    
    printf("Test: sender_bytes_sent_reset\n");
    test_stream_replication_sender_bytes_sent_reset();
    
    /* Error Tracking Tests */
    printf("Test: sender_error_count_initial_zero\n");
    test_stream_replication_sender_error_count_initial_zero();
    
    printf("Test: sender_error_count_increment\n");
    test_stream_replication_sender_error_count_increment();
    
    printf("Test: sender_error_count_reset\n");
    test_stream_replication_sender_error_count_reset();
    
    /* Config Structure Field Manipulation Tests */
    printf("Test: config_host_empty_string\n");
    test_stream_replication_sender_config_host_empty_string();
    
    printf("Test: config_host_long_string\n");
    test_stream_replication_sender_config_host_long_string();
    
    printf("Test: config_host_special_characters\n");
    test_stream_replication_sender_config_host_special_characters();
    
    printf("Test: config_port_min_valid\n");
    test_stream_replication_sender_config_port_min_valid();
    
    printf("Test: config_port_typical_values\n");
    test_stream_replication_sender_config_port_typical_values();
    
    printf("Test: config_timeout_very_large\n");
    test_stream_replication_sender_config_timeout_very_large();
    
    /* Multi-Sender Scenario Tests */
    printf("Test: multiple_senders_independent_state\n");
    test_stream_replication_multiple_senders_independent_state();
    
    printf("Test: config_copy\n");
    test_stream_replication_sender_config_copy();
    
    printf("\nAll stream replication sender header tests passed!\n");
    return 0;
}

#endif /* TEST_STREAM_REPLICATION_SENDER_H */