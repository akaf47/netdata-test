#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include <cmocka.h>
#include <time.h>

// Mock structures and forward declarations
typedef struct notification_recipient {
    char *name;
    char *endpoint;
    int enabled;
} notification_recipient_t;

typedef struct health_notification {
    char *id;
    char *name;
    char *message;
    int severity;
    time_t timestamp;
    notification_recipient_t **recipients;
    int recipient_count;
} health_notification_t;

// Forward declarations for functions to test
extern health_notification_t* health_notification_create(const char *id, const char *name);
extern int health_notification_add_recipient(health_notification_t *notification, const char *name, const char *endpoint);
extern int health_notification_send(health_notification_t *notification);
extern void health_notification_free(health_notification_t *notification);
extern int health_notification_set_message(health_notification_t *notification, const char *message);
extern int health_notification_set_severity(health_notification_t *notification, int severity);

// Test: health_notification_create with valid inputs
static void test_health_notification_create_valid(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test Notification");
    
    assert_non_null(notif);
    assert_string_equal(notif->id, "test-id");
    assert_string_equal(notif->name, "Test Notification");
    assert_int_equal(notif->severity, 0);
    assert_non_null(notif->timestamp);
    
    health_notification_free(notif);
}

// Test: health_notification_create with NULL id
static void test_health_notification_create_null_id(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create(NULL, "Test Notification");
    
    assert_null(notif);
}

// Test: health_notification_create with NULL name
static void test_health_notification_create_null_name(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", NULL);
    
    assert_null(notif);
}

// Test: health_notification_create with empty id
static void test_health_notification_create_empty_id(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("", "Test Notification");
    
    assert_null(notif);
}

// Test: health_notification_create with empty name
static void test_health_notification_create_empty_name(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "");
    
    assert_null(notif);
}

// Test: health_notification_add_recipient with valid inputs
static void test_health_notification_add_recipient_valid(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_add_recipient(notif, "admin", "admin@example.com");
    
    assert_int_equal(result, 0);
    assert_int_equal(notif->recipient_count, 1);
    
    health_notification_free(notif);
}

// Test: health_notification_add_recipient with NULL notification
static void test_health_notification_add_recipient_null_notification(void **state)
{
    (void) state;
    
    int result = health_notification_add_recipient(NULL, "admin", "admin@example.com");
    
    assert_int_equal(result, -1);
}

// Test: health_notification_add_recipient with NULL name
static void test_health_notification_add_recipient_null_name(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_add_recipient(notif, NULL, "admin@example.com");
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_add_recipient with NULL endpoint
static void test_health_notification_add_recipient_null_endpoint(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_add_recipient(notif, "admin", NULL);
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_add_recipient with empty name
static void test_health_notification_add_recipient_empty_name(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_add_recipient(notif, "", "admin@example.com");
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_add_recipient with empty endpoint
static void test_health_notification_add_recipient_empty_endpoint(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_add_recipient(notif, "admin", "");
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_add_recipient multiple recipients
static void test_health_notification_add_recipient_multiple(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result1 = health_notification_add_recipient(notif, "admin1", "admin1@example.com");
    int result2 = health_notification_add_recipient(notif, "admin2", "admin2@example.com");
    int result3 = health_notification_add_recipient(notif, "admin3", "admin3@example.com");
    
    assert_int_equal(result1, 0);
    assert_int_equal(result2, 0);
    assert_int_equal(result3, 0);
    assert_int_equal(notif->recipient_count, 3);
    
    health_notification_free(notif);
}

// Test: health_notification_set_message with valid input
static void test_health_notification_set_message_valid(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_set_message(notif, "Test message");
    
    assert_int_equal(result, 0);
    assert_string_equal(notif->message, "Test message");
    
    health_notification_free(notif);
}

// Test: health_notification_set_message with NULL notification
static void test_health_notification_set_message_null_notification(void **state)
{
    (void) state;
    
    int result = health_notification_set_message(NULL, "Test message");
    
    assert_int_equal(result, -1);
}

// Test: health_notification_set_message with NULL message
static void test_health_notification_set_message_null_message(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_set_message(notif, NULL);
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_set_message with empty message
static void test_health_notification_set_message_empty(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_set_message(notif, "");
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_set_message with long message
static void test_health_notification_set_message_long(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    char long_message[4096];
    memset(long_message, 'a', 4095);
    long_message[4095] = '\0';
    
    int result = health_notification_set_message(notif, long_message);
    
    assert_int_equal(result, 0);
    assert_string_equal(notif->message, long_message);
    
    health_notification_free(notif);
}

// Test: health_notification_set_severity with valid severity
static void test_health_notification_set_severity_valid(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_set_severity(notif, 5);
    
    assert_int_equal(result, 0);
    assert_int_equal(notif->severity, 5);
    
    health_notification_free(notif);
}

// Test: health_notification_set_severity with NULL notification
static void test_health_notification_set_severity_null_notification(void **state)
{
    (void) state;
    
    int result = health_notification_set_severity(NULL, 5);
    
    assert_int_equal(result, -1);
}

// Test: health_notification_set_severity with zero
static void test_health_notification_set_severity_zero(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_set_severity(notif, 0);
    
    assert_int_equal(result, 0);
    assert_int_equal(notif->severity, 0);
    
    health_notification_free(notif);
}

// Test: health_notification_set_severity with negative value
static void test_health_notification_set_severity_negative(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_set_severity(notif, -1);
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_set_severity with very high value
static void test_health_notification_set_severity_high(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_set_severity(notif, 1000);
    
    assert_int_equal(result, 0);
    assert_int_equal(notif->severity, 1000);
    
    health_notification_free(notif);
}

// Test: health_notification_send with valid notification
static void test_health_notification_send_valid(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_add_recipient(notif, "admin", "admin@example.com");
    health_notification_set_message(notif, "Test message");
    
    int result = health_notification_send(notif);
    
    assert_int_equal(result, 0);
    
    health_notification_free(notif);
}

// Test: health_notification_send with NULL notification
static void test_health_notification_send_null_notification(void **state)
{
    (void) state;
    
    int result = health_notification_send(NULL);
    
    assert_int_equal(result, -1);
}

// Test: health_notification_send with no recipients
static void test_health_notification_send_no_recipients(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_set_message(notif, "Test message");
    
    int result = health_notification_send(notif);
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_send with no message
static void test_health_notification_send_no_message(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_add_recipient(notif, "admin", "admin@example.com");
    
    int result = health_notification_send(notif);
    
    assert_int_equal(result, -1);
    
    health_notification_free(notif);
}

// Test: health_notification_free with valid notification
static void test_health_notification_free_valid(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_add_recipient(notif, "admin", "admin@example.com");
    health_notification_set_message(notif, "Test message");
    
    health_notification_free(notif);
    // If we reach here, no crash occurred
    assert_true(1);
}

// Test: health_notification_free with NULL notification
static void test_health_notification_free_null(void **state)
{
    (void) state;
    
    // Should not crash
    health_notification_free(NULL);
    
    assert_true(1);
}

// Test: health_notification_create and verify timestamp is set
static void test_health_notification_timestamp_set(void **state)
{
    (void) state;
    
    time_t before = time(NULL);
    health_notification_t *notif = health_notification_create("test-id", "Test");
    time_t after = time(NULL);
    
    assert_non_null(notif);
    assert_true(notif->timestamp >= before);
    assert_true(notif->timestamp <= after);
    
    health_notification_free(notif);
}

// Test: Add recipient updates recipient list correctly
static void test_health_notification_recipient_list_integrity(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_add_recipient(notif, "user1", "user1@example.com");
    assert_non_null(notif->recipients);
    assert_string_equal(notif->recipients[0]->name, "user1");
    assert_string_equal(notif->recipients[0]->endpoint, "user1@example.com");
    
    health_notification_add_recipient(notif, "user2", "user2@example.com");
    assert_string_equal(notif->recipients[1]->name, "user2");
    assert_string_equal(notif->recipients[1]->endpoint, "user2@example.com");
    
    health_notification_free(notif);
}

// Test: health_notification_send with multiple recipients
static void test_health_notification_send_multiple_recipients(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_add_recipient(notif, "admin1", "admin1@example.com");
    health_notification_add_recipient(notif, "admin2", "admin2@example.com");
    health_notification_set_message(notif, "Test message");
    
    int result = health_notification_send(notif);
    
    assert_int_equal(result, 0);
    
    health_notification_free(notif);
}

// Test: health_notification with special characters in message
static void test_health_notification_special_chars_message(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    const char *special_msg = "Test!@#$%^&*()_+-=[]{}|;:',.<>?/\\";
    int result = health_notification_set_message(notif, special_msg);
    
    assert_int_equal(result, 0);
    assert_string_equal(notif->message, special_msg);
    
    health_notification_free(notif);
}

// Test: health_notification with special characters in endpoint
static void test_health_notification_special_chars_endpoint(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    int result = health_notification_add_recipient(notif, "user", "user+tag@example.co.uk");
    
    assert_int_equal(result, 0);
    
    health_notification_free(notif);
}

// Test: Multiple set_message calls (last one wins)
static void test_health_notification_multiple_message_calls(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_set_message(notif, "First message");
    health_notification_set_message(notif, "Second message");
    health_notification_set_message(notif, "Final message");
    
    assert_string_equal(notif->message, "Final message");
    
    health_notification_free(notif);
}

// Test: Multiple set_severity calls (last one wins)
static void test_health_notification_multiple_severity_calls(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    health_notification_set_severity(notif, 1);
    health_notification_set_severity(notif, 5);
    health_notification_set_severity(notif, 10);
    
    assert_int_equal(notif->severity, 10);
    
    health_notification_free(notif);
}

// Test: Very long id and name strings
static void test_health_notification_long_strings(void **state)
{
    (void) state;
    
    char long_id[1024];
    char long_name[1024];
    memset(long_id, 'a', 1023);
    long_id[1023] = '\0';
    memset(long_name, 'b', 1023);
    long_name[1023] = '\0';
    
    health_notification_t *notif = health_notification_create(long_id, long_name);
    
    assert_non_null(notif);
    assert_string_equal(notif->id, long_id);
    assert_string_equal(notif->name, long_name);
    
    health_notification_free(notif);
}

// Test: Recipient count initialization
static void test_health_notification_recipient_count_init(void **state)
{
    (void) state;
    
    health_notification_t *notif = health_notification_create("test-id", "Test");
    assert_non_null(notif);
    
    assert_int_equal(notif->recipient_count, 0);
    assert_null(notif->recipients);
    
    health_notification_free(notif);
}

int main(void)
{
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_health_notification_create_valid),
        cmocka_unit_test(test_health_notification_create_null_id),
        cmocka_unit_test(test_health_notification_create_null_name),
        cmocka_unit_test(test_health_notification_create_empty_id),
        cmocka_unit_test(test_health_notification_create_empty_name),
        cmocka_unit_test(test_health_notification_add_recipient_valid),
        cmocka_unit_test(test_health_notification_add_recipient_null_notification),
        cmocka_unit_test(test_health_notification_add_recipient_null_name),
        cmocka_unit_test(test_health_notification_add_recipient_null_endpoint),
        cmocka_unit_test(test_health_notification_add_recipient_empty_name),
        cmocka_unit_test(test_health_notification_add_recipient_empty_endpoint),
        cmocka_unit_test(test_health_notification_add_recipient_multiple),
        cmocka_unit_test(test_health_notification_set_message_valid),
        cmocka_unit_test(test_health_notification_set_message_null_notification),
        cmocka_unit_test(test_health_notification_set_message_null_message),
        cmocka_unit_test(test_health_notification_set_message_empty),
        cmocka_unit_test(test_health_notification_set_message_long),
        cmocka_unit_test(test_health_notification_set_severity_valid),
        cmocka_unit_test(test_health_notification_set_severity_null_notification),
        cmocka_unit_test(test_health_notification_set_severity_zero),
        cmocka_unit_test(test_health_notification_set_severity_negative),
        cmocka_unit_test(test_health_notification_set_severity_high),
        cmocka_unit_test(test_health_notification_send_valid),
        cmocka_unit_test(test_health_notification_send_null_notification),
        cmocka_unit_test(test_health_notification_send_no_recipients),
        cmocka_unit_test(test_health_notification_send_no_message),
        cmocka_unit_test(test_health_notification_free_valid),
        cmocka_unit_test(test_health_notification_free_null),
        cmocka_unit_test(test_health_notification_timestamp_set),
        cmocka_unit_test(test_health_notification_recipient_list_integrity),
        cmocka_unit_test(test_health_notification_send_multiple_recipients),
        cmocka_unit_test(test_health_notification_special_chars_message),
        cmocka_unit_test(test_health_notification_special_chars_endpoint),
        cmocka_unit_test(test_health_notification_multiple_message_calls),
        cmocka_unit_test(test_health_notification_multiple_severity_calls),
        cmocka_unit_test(test_health_notification_long_strings),
        cmocka_unit_test(test_health_notification_recipient_count_init),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}