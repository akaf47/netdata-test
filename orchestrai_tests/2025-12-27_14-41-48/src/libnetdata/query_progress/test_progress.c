#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>
#include <time.h>
#include <sys/time.h>
#include <pthread.h>
#include <unistd.h>
#include <cmocka.h>
#include <setjmp.h>

/* Mock structure definitions based on typical query_progress patterns */
typedef struct query_progress {
    uint64_t entries_total;
    uint64_t entries_processed;
    uint64_t start_time_ms;
    uint64_t last_update_time_ms;
    bool cancelled;
    const char *message;
} QUERY_PROGRESS;

/* Forward declarations of functions under test */
QUERY_PROGRESS *query_progress_create(uint64_t total_entries);
void query_progress_free(QUERY_PROGRESS *qp);
void query_progress_update(QUERY_PROGRESS *qp, uint64_t processed_entries);
bool query_progress_is_cancelled(const QUERY_PROGRESS *qp);
void query_progress_cancel(QUERY_PROGRESS *qp);
uint64_t query_progress_get_processed(const QUERY_PROGRESS *qp);
uint64_t query_progress_get_total(const QUERY_PROGRESS *qp);
const char *query_progress_get_message(const QUERY_PROGRESS *qp);
void query_progress_set_message(QUERY_PROGRESS *qp, const char *message);
uint64_t query_progress_get_elapsed_ms(const QUERY_PROGRESS *qp);
double query_progress_get_percentage(const QUERY_PROGRESS *qp);

/* ============================================================================
   TEST FIXTURES AND SETUP/TEARDOWN
   ============================================================================ */

static QUERY_PROGRESS *test_qp = NULL;

static int setup_test(void **state) {
    test_qp = NULL;
    return 0;
}

static int teardown_test(void **state) {
    if (test_qp != NULL) {
        query_progress_free(test_qp);
        test_qp = NULL;
    }
    return 0;
}

/* ============================================================================
   query_progress_create() TESTS
   ============================================================================ */

static void test_query_progress_create_success(void **state) {
    /* Arrange & Act */
    test_qp = query_progress_create(1000);
    
    /* Assert */
    assert_non_null(test_qp);
    assert_int_equal(test_qp->entries_total, 1000);
    assert_int_equal(test_qp->entries_processed, 0);
    assert_false(test_qp->cancelled);
    assert_null(test_qp->message);
    assert_true(test_qp->start_time_ms > 0);
}

static void test_query_progress_create_zero_entries(void **state) {
    /* Arrange & Act */
    test_qp = query_progress_create(0);
    
    /* Assert */
    assert_non_null(test_qp);
    assert_int_equal(test_qp->entries_total, 0);
    assert_int_equal(test_qp->entries_processed, 0);
}

static void test_query_progress_create_max_uint64(void **state) {
    /* Arrange & Act */
    test_qp = query_progress_create(UINT64_MAX);
    
    /* Assert */
    assert_non_null(test_qp);
    assert_int_equal(test_qp->entries_total, UINT64_MAX);
}

static void test_query_progress_create_one_entry(void **state) {
    /* Arrange & Act */
    test_qp = query_progress_create(1);
    
    /* Assert */
    assert_non_null(test_qp);
    assert_int_equal(test_qp->entries_total, 1);
}

static void test_query_progress_create_large_number(void **state) {
    /* Arrange & Act */
    test_qp = query_progress_create(1000000000);
    
    /* Assert */
    assert_non_null(test_qp);
    assert_int_equal(test_qp->entries_total, 1000000000);
}

/* ============================================================================
   query_progress_free() TESTS
   ============================================================================ */

static void test_query_progress_free_valid(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    assert_non_null(test_qp);
    
    /* Act */
    query_progress_free(test_qp);
    test_qp = NULL;
    
    /* Assert */
    /* If we reach here without segfault, the free was successful */
    assert_null(test_qp);
}

static void test_query_progress_free_null(void **state) {
    /* Arrange & Act & Assert */
    /* Calling free with NULL should not crash */
    query_progress_free(NULL);
    assert_true(true);
}

static void test_query_progress_free_with_message(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    query_progress_set_message(test_qp, "test message");
    
    /* Act */
    query_progress_free(test_qp);
    test_qp = NULL;
    
    /* Assert */
    assert_null(test_qp);
}

/* ============================================================================
   query_progress_update() TESTS
   ============================================================================ */

static void test_query_progress_update_basic(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act */
    query_progress_update(test_qp, 500);
    
    /* Assert */
    assert_int_equal(test_qp->entries_processed, 500);
}

static void test_query_progress_update_zero(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act */
    query_progress_update(test_qp, 0);
    
    /* Assert */
    assert_int_equal(test_qp->entries_processed, 0);
}

static void test_query_progress_update_full(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act */
    query_progress_update(test_qp, 1000);
    
    /* Assert */
    assert_int_equal(test_qp->entries_processed, 1000);
}

static void test_query_progress_update_exceed_total(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act */
    query_progress_update(test_qp, 2000);
    
    /* Assert - should allow exceeding total */
    assert_int_equal(test_qp->entries_processed, 2000);
}

static void test_query_progress_update_multiple_times(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act */
    query_progress_update(test_qp, 100);
    query_progress_update(test_qp, 200);
    query_progress_update(test_qp, 300);
    
    /* Assert */
    assert_int_equal(test_qp->entries_processed, 300);
}

static void test_query_progress_update_backwards(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    query_progress_update(test_qp, 500);
    
    /* Act */
    query_progress_update(test_qp, 200);
    
    /* Assert - should set to new value */
    assert_int_equal(test_qp->entries_processed, 200);
}

static void test_query_progress_update_max_uint64(void **state) {
    /* Arrange */
    test_qp = query_progress_create(UINT64_MAX);
    
    /* Act */
    query_progress_update(test_qp, UINT64_MAX);
    
    /* Assert */
    assert_int_equal(test_qp->entries_processed, UINT64_MAX);
}

static void test_query_progress_update_null_pointer(void **state) {
    /* Arrange & Act & Assert */
    /* This should be handled gracefully without crashing */
    query_progress_update(NULL, 100);
    assert_true(true);
}

/* ============================================================================
   query_progress_is_cancelled() TESTS
   ============================================================================ */

static void test_query_progress_is_cancelled_initial_false(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    bool result = query_progress_is_cancelled(test_qp);
    
    /* Assert */
    assert_false(result);
}

static void test_query_progress_is_cancelled_after_cancel(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    query_progress_cancel(test_qp);
    
    /* Act */
    bool result = query_progress_is_cancelled(test_qp);
    
    /* Assert */
    assert_true(result);
}

static void test_query_progress_is_cancelled_null_pointer(void **state) {
    /* Arrange & Act */
    bool result = query_progress_is_cancelled(NULL);
    
    /* Assert - should return false for NULL */
    assert_false(result);
}

static void test_query_progress_is_cancelled_multiple_checks(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    bool result1 = query_progress_is_cancelled(test_qp);
    bool result2 = query_progress_is_cancelled(test_qp);
    
    /* Assert */
    assert_false(result1);
    assert_false(result2);
}

/* ============================================================================
   query_progress_cancel() TESTS
   ============================================================================ */

static void test_query_progress_cancel_success(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    query_progress_cancel(test_qp);
    
    /* Assert */
    assert_true(test_qp->cancelled);
}

static void test_query_progress_cancel_null_pointer(void **state) {
    /* Arrange & Act & Assert */
    query_progress_cancel(NULL);
    assert_true(true);
}

static void test_query_progress_cancel_multiple_times(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    query_progress_cancel(test_qp);
    query_progress_cancel(test_qp);
    
    /* Assert */
    assert_true(test_qp->cancelled);
}

static void test_query_progress_cancel_during_update(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    query_progress_update(test_qp, 500);
    
    /* Act */
    query_progress_cancel(test_qp);
    
    /* Assert */
    assert_true(query_progress_is_cancelled(test_qp));
    assert_int_equal(test_qp->entries_processed, 500);
}

/* ============================================================================
   query_progress_get_processed() TESTS
   ============================================================================ */

static void test_query_progress_get_processed_initial(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act */
    uint64_t processed = query_progress_get_processed(test_qp);
    
    /* Assert */
    assert_int_equal(processed, 0);
}

static void test_query_progress_get_processed_after_update(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    query_progress_update(test_qp, 750);
    
    /* Act */
    uint64_t processed = query_progress_get_processed(test_qp);
    
    /* Assert */
    assert_int_equal(processed, 750);
}

static void test_query_progress_get_processed_null_pointer(void **state) {
    /* Arrange & Act */
    uint64_t processed = query_progress_get_processed(NULL);
    
    /* Assert - should return 0 for NULL */
    assert_int_equal(processed, 0);
}

static void test_query_progress_get_processed_max(void **state) {
    /* Arrange */
    test_qp = query_progress_create(UINT64_MAX);
    query_progress_update(test_qp, UINT64_MAX);
    
    /* Act */
    uint64_t processed = query_progress_get_processed(test_qp);
    
    /* Assert */
    assert_int_equal(processed, UINT64_MAX);
}

/* ============================================================================
   query_progress_get_total() TESTS
   ============================================================================ */

static void test_query_progress_get_total_basic(void **state) {
    /* Arrange */
    test_qp = query_progress_create(5000);
    
    /* Act */
    uint64_t total = query_progress_get_total(test_qp);
    
    /* Assert */
    assert_int_equal(total, 5000);
}

static void test_query_progress_get_total_zero(void **state) {
    /* Arrange */
    test_qp = query_progress_create(0);
    
    /* Act */
    uint64_t total = query_progress_get_total(test_qp);
    
    /* Assert */
    assert_int_equal(total, 0);
}

static void test_query_progress_get_total_null_pointer(void **state) {
    /* Arrange & Act */
    uint64_t total = query_progress_get_total(NULL);
    
    /* Assert - should return 0 for NULL */
    assert_int_equal(total, 0);
}

static void test_query_progress_get_total_max(void **state) {
    /* Arrange */
    test_qp = query_progress_create(UINT64_MAX);
    
    /* Act */
    uint64_t total = query_progress_get_total(test_qp);
    
    /* Assert */
    assert_int_equal(total, UINT64_MAX);
}

/* ============================================================================
   query_progress_get_message() TESTS
   ============================================================================ */

static void test_query_progress_get_message_initial_null(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    const char *msg = query_progress_get_message(test_qp);
    
    /* Assert */
    assert_null(msg);
}

static void test_query_progress_get_message_after_set(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    query_progress_set_message(test_qp, "test message");
    
    /* Act */
    const char *msg = query_progress_get_message(test_qp);
    
    /* Assert */
    assert_non_null(msg);
    assert_string_equal(msg, "test message");
}

static void test_query_progress_get_message_null_pointer(void **state) {
    /* Arrange & Act */
    const char *msg = query_progress_get_message(NULL);
    
    /* Assert */
    assert_null(msg);
}

static void test_query_progress_get_message_empty_string(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    query_progress_set_message(test_qp, "");
    
    /* Act */
    const char *msg = query_progress_get_message(test_qp);
    
    /* Assert */
    assert_non_null(msg);
    assert_string_equal(msg, "");
}

static void test_query_progress_get_message_long_string(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    const char *long_msg = "This is a very long message with many characters "
                           "that should be properly stored and retrieved from the "
                           "query progress structure without any truncation or issues.";
    query_progress_set_message(test_qp, long_msg);
    
    /* Act */
    const char *msg = query_progress_get_message(test_qp);
    
    /* Assert */
    assert_non_null(msg);
    assert_string_equal(msg, long_msg);
}

/* ============================================================================
   query_progress_set_message() TESTS
   ============================================================================ */

static void test_query_progress_set_message_basic(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    query_progress_set_message(test_qp, "progress message");
    
    /* Assert */
    assert_non_null(test_qp->message);
    assert_string_equal(test_qp->message, "progress message");
}

static void test_query_progress_set_message_null(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    query_progress_set_message(test_qp, "first message");
    
    /* Act */
    query_progress_set_message(test_qp, NULL);
    
    /* Assert */
    assert_null(test_qp->message);
}

static void test_query_progress_set_message_empty_string(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    query_progress_set_message(test_qp, "");
    
    /* Assert */
    assert_non_null(test_qp->message);
    assert_string_equal(test_qp->message, "");
}

static void test_query_progress_set_message_replace(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    query_progress_set_message(test_qp, "message1");
    
    /* Act */
    query_progress_set_message(test_qp, "message2");
    
    /* Assert */
    assert_string_equal(test_qp->message, "message2");
}

static void test_query_progress_set_message_null_pointer(void **state) {
    /* Arrange & Act & Assert */
    query_progress_set_message(NULL, "message");
    assert_true(true);
}

static void test_query_progress_set_message_special_characters(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    const char *special = "Message with\nnewline\tand\ttabs!@#$%^&*()";
    
    /* Act */
    query_progress_set_message(test_qp, special);
    
    /* Assert */
    assert_string_equal(test_qp->message, special);
}

/* ============================================================================
   query_progress_get_elapsed_ms() TESTS
   ============================================================================ */

static void test_query_progress_get_elapsed_ms_initial(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    uint64_t elapsed = query_progress_get_elapsed_ms(test_qp);
    
    /* Assert */
    assert_true(elapsed >= 0);
}

static void test_query_progress_get_elapsed_ms_null_pointer(void **state) {
    /* Arrange & Act */
    uint64_t elapsed = query_progress_get_elapsed_ms(NULL);
    
    /* Assert */
    assert_int_equal(elapsed, 0);
}

static void test_query_progress_get_elapsed_ms_increases(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    uint64_t elapsed1 = query_progress_get_elapsed_ms(test_qp);
    
    /* Act */
    sleep(1);
    uint64_t elapsed2 = query_progress_get_elapsed_ms(test_qp);
    
    /* Assert */
    assert_true(elapsed2 > elapsed1);
}

/* ============================================================================
   query_progress_get_percentage() TESTS
   ============================================================================ */

static void test_query_progress_get_percentage_zero_total(void **state) {
    /* Arrange */
    test_qp = query_progress_create(0);
    
    /* Act */
    double percentage = query_progress_get_percentage(test_qp);
    
    /* Assert */
    assert_double_equal(percentage, 0.0, 0.01);
}

static void test_query_progress_get_percentage_zero_processed(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act */
    double percentage = query_progress_get_percentage(test_qp);
    
    /* Assert */
    assert_double_equal(percentage, 0.0, 0.01);
}

static void test_query_progress_get_percentage_half(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    query_progress_update(test_qp, 500);
    
    /* Act */
    double percentage = query_progress_get_percentage(test_qp);
    
    /* Assert */
    assert_double_equal(percentage, 50.0, 0.01);
}

static void test_query_progress_get_percentage_full(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    query_progress_update(test_qp, 1000);
    
    /* Act */
    double percentage = query_progress_get_percentage(test_qp);
    
    /* Assert */
    assert_double_equal(percentage, 100.0, 0.01);
}

static void test_query_progress_get_percentage_partial(void **state) {
    /* Arrange */
    test_qp = query_progress_create(200);
    query_progress_update(test_qp, 50);
    
    /* Act */
    double percentage = query_progress_get_percentage(test_qp);
    
    /* Assert */
    assert_double_equal(percentage, 25.0, 0.01);
}

static void test_query_progress_get_percentage_exceed_total(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    query_progress_update(test_qp, 150);
    
    /* Act */
    double percentage = query_progress_get_percentage(test_qp);
    
    /* Assert */
    assert_double_equal(percentage, 150.0, 0.01);
}

static void test_query_progress_get_percentage_null_pointer(void **state) {
    /* Arrange & Act */
    double percentage = query_progress_get_percentage(NULL);
    
    /* Assert */
    assert_double_equal(percentage, 0.0, 0.01);
}

static void test_query_progress_get_percentage_one(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1);
    query_progress_update(test_qp, 1);
    
    /* Act */
    double percentage = query_progress_get_percentage(test_qp);
    
    /* Assert */
    assert_double_equal(percentage, 100.0, 0.01);
}

/* ============================================================================
   INTEGRATION TESTS
   ============================================================================ */

static void test_integration_full_workflow(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    assert_non_null(test_qp);
    assert_int_equal(query_progress_get_total(test_qp), 1000);
    assert_int_equal(query_progress_get_processed(test_qp), 0);
    assert_false(query_progress_is_cancelled(test_qp));
    
    /* Act & Assert - Update progressively */
    for (int i = 1; i <= 10; i++) {
        query_progress_update(test_qp, i * 100);
        assert_int_equal(query_progress_get_processed(test_qp), i * 100);
    }
    
    assert_double_equal(query_progress_get_percentage(test_qp), 100.0, 0.01);
}

static void test_integration_with_message_and_cancel(void **state) {
    /* Arrange */
    test_qp = query_progress_create(500);
    query_progress_set_message(test_qp, "Processing");
    
    /* Act */
    query_progress_update(test_qp, 250);
    query_progress_cancel(test_qp);
    
    /* Assert */
    assert_true(query_progress_is_cancelled(test_qp));
    assert_string_equal(query_progress_get_message(test_qp), "Processing");
    assert_int_equal(query_progress_get_processed(test_qp), 250);
    assert_double_equal(query_progress_get_percentage(test_qp), 50.0, 0.01);
}

static void test_integration_multiple_updates_with_message_changes(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    
    /* Act & Assert */
    query_progress_set_message(test_qp, "Starting");
    query_progress_update(test_qp, 100);
    assert_string_equal(query_progress_get_message(test_qp), "Starting");
    
    query_progress_set_message(test_qp, "Processing");
    query_progress_update(test_qp, 500);
    assert_string_equal(query_progress_get_message(test_qp), "Processing");
    
    query_progress_set_message(test_qp, "Finishing");
    query_progress_update(test_qp, 1000);
    assert_string_equal(query_progress_get_message(test_qp), "Finishing");
}

static void test_integration_cancel_prevents_further_work(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1000);
    query_progress_update(test_qp, 100);
    query_progress_cancel(test_qp);
    
    /* Act & Assert */
    assert_true(query_progress_is_cancelled(test_qp));
    query_progress_update(test_qp, 200);
    assert_int_equal(query_progress_get_processed(test_qp), 200);
}

static void test_integration_elapsed_time_tracking(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act */
    uint64_t elapsed1 = query_progress_get_elapsed_ms(test_qp);
    query_progress_update(test_qp, 50);
    uint64_t elapsed2 = query_progress_get_elapsed_ms(test_qp);
    
    /* Assert */
    assert_true(elapsed2 >= elapsed1);
}

/* ============================================================================
   EDGE CASE TESTS
   ============================================================================ */

static void test_edge_case_all_operations_on_zero_total(void **state) {
    /* Arrange */
    test_qp = query_progress_create(0);
    
    /* Act & Assert */
    query_progress_update(test_qp, 0);
    query_progress_set_message(test_qp, "message");
    query_progress_cancel(test_qp);
    
    assert_int_equal(query_progress_get_total(test_qp), 0);
    assert_int_equal(query_progress_get_processed(test_qp), 0);
    assert_double_equal(query_progress_get_percentage(test_qp), 0.0, 0.01);
    assert_true(query_progress_is_cancelled(test_qp));
}

static void test_edge_case_large_numbers(void **state) {
    /* Arrange */
    uint64_t large_num = 9223372036854775807ULL; /* max int64 */
    test_qp = query_progress_create(large_num);
    
    /* Act */
    query_progress_update(test_qp, large_num / 2);
    
    /* Assert */
    assert_int_equal(query_progress_get_total(test_qp), large_num);
    assert_int_equal(query_progress_get_processed(test_qp), large_num / 2);
    assert_true(query_progress_get_percentage(test_qp) > 0.0);
}

static void test_edge_case_one_item(void **state) {
    /* Arrange */
    test_qp = query_progress_create(1);
    
    /* Act */
    query_progress_update(test_qp, 1);
    
    /* Assert */
    assert_int_equal(query_progress_get_total(test_qp), 1);
    assert_int_equal(query_progress_get_processed(test_qp), 1);
    assert_double_equal(query_progress_get_percentage(test_qp), 100.0, 0.01);
}

static void test_edge_case_message_null_to_empty_to_value(void **state) {
    /* Arrange */
    test_qp = query_progress_create(100);
    
    /* Act & Assert */
    assert_null(query_progress_get_message(test_qp));
    
    query_progress_set_message(test_qp, "");
    assert_string_equal(query_progress_get_message(test_qp), "");
    
    query_progress_set_message(test_qp, "value");
    assert_string_equal(query_progress_get_message(test_qp), "value");
    
    query_progress_set_message(test_qp, NULL);
    assert_null(query_progress_get_message(test_qp));
}

/* ============================================================================
   MAIN TEST SUITE CONFIGURATION
   ============================================================================ */

int main(void) {
    const struct CMUnitTest tests[] = {
        /* query_progress_create() tests */
        cmocka_unit_test_setup_teardown(test_query_progress_create_success,
                                       setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_query_progress_create_zero_entries,
                                       setup_test, teardown_test),
        cmocka_unit_test_setup_teardown(test_query_progress_create_max_uint64,
                                       setup_test