#ifndef TEST_POLL_EVENTS_H
#define TEST_POLL_EVENTS_H

#include <sys/poll.h>

/*
 * Test Suite for poll-events.h header file
 * 
 * This test file provides comprehensive coverage of:
 * - All macro definitions
 * - All type definitions
 * - All function declarations
 * - All constants
 * - Edge cases and boundary conditions
 */

/* Test macro definitions */
void test_macro_definitions(void);

/* Test type definitions */
void test_poll_event_type_definition(void);

/* Test function declarations - Initialization */
void test_declare_poll_events_init(void);
void test_declare_poll_events_cleanup(void);
void test_declare_poll_events_destroy(void);

/* Test function declarations - Add/Remove */
void test_declare_poll_events_add(void);
void test_declare_poll_events_remove(void);

/* Test function declarations - Poll */
void test_declare_poll_events_poll(void);
void test_declare_poll_events_poll_timeout(void);
void test_declare_poll_events_poll_blocking(void);

/* Test function declarations - Get */
void test_declare_poll_events_get_fd(void);
void test_declare_poll_events_get_data(void);
void test_declare_poll_events_get_revents(void);
void test_declare_poll_events_get_events(void);

/* Test function declarations - Utilities */
void test_declare_poll_events_count(void);
void test_declare_poll_events_clear(void);
void test_declare_poll_events_update(void);
void test_declare_poll_events_exists(void);

/* Test header include guards */
void test_header_include_guards(void);

/* Test for circular dependencies */
void test_no_circular_dependencies(void);

/* Test constant definitions */
void test_constant_definitions(void);

/* Test size definitions */
void test_size_definitions(void);

/* Test return code definitions */
void test_return_code_definitions(void);

/* Integration test - header completeness */
void test_header_completeness(void);

#endif /* TEST_POLL_EVENTS_H */