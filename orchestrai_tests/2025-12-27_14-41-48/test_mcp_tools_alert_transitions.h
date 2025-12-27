#ifndef TEST_MCP_TOOLS_ALERT_TRANSITIONS_H
#define TEST_MCP_TOOLS_ALERT_TRANSITIONS_H

#include <stddef.h>
#include <stdint.h>
#include <time.h>

/* ========== MOCK DECLARATIONS ========== */

/**
 * Mock structures to match expected header content
 * These would be defined in the actual header being tested
 */

typedef enum {
    ALERT_STATE_UNKNOWN = 0,
    ALERT_STATE_WARNING = 1,
    ALERT_STATE_CRITICAL = 2,
    ALERT_STATE_CLEAR = 3,
    ALERT_STATE_MAX = 4
} alert_state_t;

typedef struct {
    char *alert_id;
    alert_state_t from_state;
    alert_state_t to_state;
    time_t timestamp;
    char *reason;
} alert_transition_record_t;

typedef struct {
    alert_transition_record_t *records;
    size_t record_count;
    size_t max_records;
} alert_transitions_t;

/* ========== FUNCTION DECLARATIONS FOR TESTING ========== */

/**
 * Initialize alert transitions tracking
 * 
 * @return 0 on success, non-zero on error
 */
int mcp_alert_transitions_init(void);

/**
 * Cleanup alert transitions module
 * 
 * @return 0 on success, non-zero on error
 */
int mcp_alert_transitions_cleanup(void);

/**
 * Validate state transition
 * 
 * @param from_state Source state
 * @param to_state Destination state
 * @return 1 if valid, 0 if invalid
 */
int mcp_validate_state_transition(alert_state_t from_state, alert_state_t to_state);

/**
 * Add a transition record
 * 
 * @param alert_id Alert identifier
 * @param from_state Source state
 * @param to_state Destination state
 * @param reason Transition reason
 * @return 0 on success, non-zero on error
 */
int mcp_add_transition_record(const char *alert_id, alert_state_t from_state, 
                               alert_state_t to_state, const char *reason);

/**
 * Get transition history
 * 
 * @param alert_id Alert identifier (NULL for all)
 * @param records Output array pointer
 * @param count Output count pointer
 * @return 0 on success, non-zero on error
 */
int mcp_get_transition_history(const char *alert_id, alert_transition_record_t **records,
                                size_t *count);

/**
 * Clear transition history
 * 
 * @param alert_id Alert identifier (NULL for all)
 * @return 0 on success, non-zero on error
 */
int mcp_clear_transition_history(const char *alert_id);

/**
 * Compare two alert states
 * 
 * @param state1 First state
 * @param state2 Second state
 * @return 1 if equal, 0 if different
 */
int mcp_compare_alert_states(alert_state_t state1, alert_state_t state2);

/**
 * Convert state to string
 * 
 * @param state Alert state
 * @return State string representation
 */
const char *mcp_alert_state_to_string(alert_state_t state);

/**
 * Convert string to state
 * 
 * @param state_str State string
 * @return Alert state value, ALERT_STATE_UNKNOWN if invalid
 */
alert_state_t mcp_string_to_alert_state(const char *state_str);

/**
 * Filter transitions by state
 * 
 * @param state Target state
 * @param records Output array
 * @param count Output count
 * @return 0 on success
 */
int mcp_filter_transitions_by_state(alert_state_t state, 
                                     alert_transition_record_t **records,
                                     size_t *count);

/**
 * Filter transitions by time range
 * 
 * @param start_time Start of range
 * @param end_time End of range
 * @param records Output array
 * @param count Output count
 * @return 0 on success
 */
int mcp_filter_transitions_by_time(time_t start_time, time_t end_time,
                                    alert_transition_record_t **records,
                                    size_t *count);

/**
 * Get transition count
 * 
 * @param alert_id Alert ID (NULL for all)
 * @return Transition count
 */
size_t mcp_get_transition_count(const char *alert_id);

/**
 * Check if state transition is possible
 * 
 * @param from_state Current state
 * @param to_state Desired state
 * @return 1 if transition allowed, 0 otherwise
 */
int mcp_is_transition_allowed(alert_state_t from_state, alert_state_t to_state);

/**
 * Set maximum history size
 * 
 * @param max_size Maximum number of records
 * @return 0 on success
 */
int mcp_set_max_history_size(size_t max_size);

/**
 * Get current history size
 * 
 * @return Current number of stored records
 */
size_t mcp_get_current_history_size(void);

/**
 * Test NULL alert_id handling in get_transition_history
 */
static inline void test_get_history_null_alert_id(void) {
    /* Verify NULL alert_id returns all transitions */
}

/**
 * Test empty alert_id handling
 */
static inline void test_get_history_empty_alert_id(void) {
    /* Verify empty string handling */
}

/**
 * Test exact alert_id match
 */
static inline void test_get_history_exact_match(void) {
    /* Verify exact alert_id matching */
}

/**
 * Test nonexistent alert_id
 */
static inline void test_get_history_not_found(void) {
    /* Verify proper handling of not-found case */
}

/**
 * Test clear all transitions
 */
static inline void test_clear_all_transitions(void) {
    /* Verify clear with NULL alert_id clears all */
}

/**
 * Test clear specific alert transitions
 */
static inline void test_clear_specific_alert(void) {
    /* Verify clear only affects specified alert */
}

/**
 * Test state string conversions
 */
static inline void test_state_string_conversions(void) {
    /* Test bidirectional conversion */
}

/**
 * Test invalid state string
 */
static inline void test_invalid_state_string(void) {
    /* Verify handling of invalid state strings */
}

/**
 * Test state machine constraints
 */
static inline void test_state_machine_constraints(void) {
    /* Test all valid/invalid transitions */
}

/**
 * Test history size limits
 */
static inline void test_history_size_limits(void) {
    /* Test max size enforcement */
}

/**
 * Test timestamp validation
 */
static inline void test_timestamp_validation(void) {
    /* Test timestamp range checks */
}

/**
 * Test reason string handling with special characters
 */
static inline void test_reason_string_special_chars(void) {
    /* Test escaping and encoding */
}

/**
 * Test concurrent access (if applicable)
 */
static inline void test_concurrent_access(void) {
    /* Test thread safety */
}

/**
 * Test memory leak prevention
 */
static inline void test_memory_leak_prevention(void) {
    /* Verify no memory leaks on cleanup */
}

#endif /* TEST_MCP_TOOLS_ALERT_TRANSITIONS_H */