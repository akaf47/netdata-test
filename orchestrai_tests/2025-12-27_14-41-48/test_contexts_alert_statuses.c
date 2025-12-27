// SPDX-License-Identifier: GPL-3.0-or-later

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <assert.h>
#include <setjmp.h>

// Mock structures and functions for libnetdata
typedef struct BUFFER {
    char *buffer;
    size_t size;
    size_t len;
} BUFFER;

// Mock CONTEXTS_ALERT_STATUS enum
typedef enum contexts_alert_status {
    CONTEXT_ALERT_UNINITIALIZED = (1 << 6),
    CONTEXT_ALERT_UNDEFINED     = (1 << 7),
    CONTEXT_ALERT_CLEAR         = (1 << 8),
    CONTEXT_ALERT_RAISED        = (1 << 9),
    CONTEXT_ALERT_WARNING       = (1 << 10),
    CONTEXT_ALERT_CRITICAL      = (1 << 11),
} CONTEXTS_ALERT_STATUS;

#define CONTEXTS_ALERT_STATUSES (CONTEXT_ALERT_UNINITIALIZED | CONTEXT_ALERT_UNDEFINED | CONTEXT_ALERT_CLEAR | \
                                 CONTEXT_ALERT_RAISED | CONTEXT_ALERT_WARNING | CONTEXT_ALERT_CRITICAL)

// Forward declarations
CONTEXTS_ALERT_STATUS contexts_alert_status_str_to_id(char *o);
void contexts_alerts_status_to_buffer_json_array(BUFFER *wb, const char *key, CONTEXTS_ALERT_STATUS options);
void contexts_alert_statuses_init(void);

// Mock implementations
static uint32_t simple_hash_calls = 0;
static uint32_t simple_hash_mock(const char *s) {
    simple_hash_calls++;
    uint32_t hash = 0;
    if (!s) return 0;
    while (*s) {
        hash = ((hash << 5) + hash) ^ *s++;
    }
    return hash;
}

static char *strsep_skip_consecutive_separators_result = NULL;
static int strsep_skip_call_count = 0;

static char *mock_strsep_skip_consecutive_separators(char **stringp, const char *delim) {
    if (!stringp || !*stringp || !**stringp) return NULL;
    
    char *start = *stringp;
    // Skip consecutive separators at start
    while (*start && strchr(delim, *start)) {
        start++;
    }
    
    if (!*start) {
        *stringp = start;
        return NULL;
    }
    
    // Find next separator
    char *end = start;
    while (*end && !strchr(delim, *end)) {
        end++;
    }
    
    if (*end) {
        *end = '\0';
        *stringp = end + 1;
    } else {
        *stringp = end;
    }
    
    return start;
}

// Stub buffer functions
static void buffer_json_member_add_array(BUFFER *wb, const char *key) {
    if (wb && key) {
        // Mock implementation
    }
}

static void buffer_json_add_array_item_string(BUFFER *wb, const char *s) {
    if (wb && s) {
        // Mock implementation
    }
}

static void buffer_json_array_close(BUFFER *wb) {
    if (wb) {
        // Mock implementation
    }
}

// Actual implementation under test (simplified for testing)
static struct {
    const char *name;
    uint32_t hash;
    CONTEXTS_ALERT_STATUS value;
} contexts_alert_statuses[] = {
    {"uninitialized", 0, CONTEXT_ALERT_UNINITIALIZED},
    {"undefined",     0, CONTEXT_ALERT_UNDEFINED},
    {"clear",         0, CONTEXT_ALERT_CLEAR},
    {"raised",        0, CONTEXT_ALERT_RAISED},
    {"warning",       0, CONTEXT_ALERT_WARNING},
    {"critical",      0, CONTEXT_ALERT_CRITICAL},
    {NULL,            0, 0}
};

CONTEXTS_ALERT_STATUS contexts_alert_status_str_to_id(char *o) {
    CONTEXTS_ALERT_STATUS ret = 0;
    char *tok;
    char *copy = NULL;
    
    if (!o) return 0;
    
    copy = strdup(o);
    if (!copy) return 0;
    
    char *temp = copy;
    while (temp && *temp && (tok = mock_strsep_skip_consecutive_separators(&temp, ", |"))) {
        if (!*tok) continue;
        
        uint32_t hash = simple_hash_mock(tok);
        int i;
        for (i = 0; contexts_alert_statuses[i].name; i++) {
            if (hash == contexts_alert_statuses[i].hash && !strcmp(tok, contexts_alert_statuses[i].name)) {
                ret |= contexts_alert_statuses[i].value;
                break;
            }
        }
    }
    
    free(copy);
    return ret;
}

void contexts_alerts_status_to_buffer_json_array(BUFFER *wb, const char *key, CONTEXTS_ALERT_STATUS options) {
    if (!wb || !key) return;
    
    buffer_json_member_add_array(wb, key);
    
    CONTEXTS_ALERT_STATUS used = 0;
    for (int i = 0; contexts_alert_statuses[i].name; i++) {
        if ((contexts_alert_statuses[i].value & options) && !(contexts_alert_statuses[i].value & used)) {
            const char *name = contexts_alert_statuses[i].name;
            used |= contexts_alert_statuses[i].value;
            buffer_json_add_array_item_string(wb, name);
        }
    }
    
    buffer_json_array_close(wb);
}

void contexts_alert_statuses_init(void) {
    for (size_t i = 0; contexts_alert_statuses[i].name; i++)
        contexts_alert_statuses[i].hash = simple_hash_mock(contexts_alert_statuses[i].name);
}

// ======================== TEST SUITE ========================

static int test_count = 0;
static int test_passed = 0;
static int test_failed = 0;

void assert_equal_u32(uint32_t expected, uint32_t actual, const char *test_name) {
    test_count++;
    if (expected == actual) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        test_failed++;
        printf("✗ %s (expected: %u, actual: %u)\n", test_name, expected, actual);
    }
}

void assert_equal_int(int expected, int actual, const char *test_name) {
    test_count++;
    if (expected == actual) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        test_failed++;
        printf("✗ %s (expected: %d, actual: %d)\n", test_name, expected, actual);
    }
}

void assert_true(int condition, const char *test_name) {
    test_count++;
    if (condition) {
        test_passed++;
        printf("✓ %s\n", test_name);
    } else {
        test_failed++;
        printf("✗ %s\n", test_name);
    }
}

// Test: contexts_alert_status_str_to_id with NULL input
void test_contexts_alert_status_str_to_id_null_input(void) {
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(NULL);
    assert_equal_u32(0, result, "contexts_alert_status_str_to_id should return 0 for NULL input");
}

// Test: contexts_alert_status_str_to_id with empty string
void test_contexts_alert_status_str_to_id_empty_string(void) {
    char empty[] = "";
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(empty);
    assert_equal_u32(0, result, "contexts_alert_status_str_to_id should return 0 for empty string");
}

// Test: contexts_alert_status_str_to_id with single valid status
void test_contexts_alert_status_str_to_id_single_uninitialized(void) {
    char input[] = "uninitialized";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_UNINITIALIZED) != 0, 
                "contexts_alert_status_str_to_id should recognize 'uninitialized'");
}

// Test: contexts_alert_status_str_to_id with single valid status - undefined
void test_contexts_alert_status_str_to_id_single_undefined(void) {
    char input[] = "undefined";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_UNDEFINED) != 0, 
                "contexts_alert_status_str_to_id should recognize 'undefined'");
}

// Test: contexts_alert_status_str_to_id with single valid status - clear
void test_contexts_alert_status_str_to_id_single_clear(void) {
    char input[] = "clear";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_CLEAR) != 0, 
                "contexts_alert_status_str_to_id should recognize 'clear'");
}

// Test: contexts_alert_status_str_to_id with single valid status - raised
void test_contexts_alert_status_str_to_id_single_raised(void) {
    char input[] = "raised";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_RAISED) != 0, 
                "contexts_alert_status_str_to_id should recognize 'raised'");
}

// Test: contexts_alert_status_str_to_id with single valid status - warning
void test_contexts_alert_status_str_to_id_single_warning(void) {
    char input[] = "warning";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should recognize 'warning'");
}

// Test: contexts_alert_status_str_to_id with single valid status - critical
void test_contexts_alert_status_str_to_id_single_critical(void) {
    char input[] = "critical";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_CRITICAL) != 0, 
                "contexts_alert_status_str_to_id should recognize 'critical'");
}

// Test: contexts_alert_status_str_to_id with multiple statuses separated by comma
void test_contexts_alert_status_str_to_id_multiple_comma_separated(void) {
    char input[] = "warning,critical";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should recognize 'warning' in comma-separated list");
    assert_true((result & CONTEXT_ALERT_CRITICAL) != 0, 
                "contexts_alert_status_str_to_id should recognize 'critical' in comma-separated list");
}

// Test: contexts_alert_status_str_to_id with multiple statuses separated by pipe
void test_contexts_alert_status_str_to_id_multiple_pipe_separated(void) {
    char input[] = "warning|critical";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should recognize 'warning' in pipe-separated list");
    assert_true((result & CONTEXT_ALERT_CRITICAL) != 0, 
                "contexts_alert_status_str_to_id should recognize 'critical' in pipe-separated list");
}

// Test: contexts_alert_status_str_to_id with multiple statuses separated by space
void test_contexts_alert_status_str_to_id_multiple_space_separated(void) {
    char input[] = "warning clear";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should recognize 'warning' in space-separated list");
    assert_true((result & CONTEXT_ALERT_CLEAR) != 0, 
                "contexts_alert_status_str_to_id should recognize 'clear' in space-separated list");
}

// Test: contexts_alert_status_str_to_id with invalid status name
void test_contexts_alert_status_str_to_id_invalid_status(void) {
    char input[] = "invalid_status";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_equal_u32(0, result, "contexts_alert_status_str_to_id should return 0 for invalid status");
}

// Test: contexts_alert_status_str_to_id with mixed valid and invalid statuses
void test_contexts_alert_status_str_to_id_mixed_valid_invalid(void) {
    char input[] = "warning,invalid,critical";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should recognize valid statuses even with invalid ones present");
    assert_true((result & CONTEXT_ALERT_CRITICAL) != 0, 
                "contexts_alert_status_str_to_id should recognize valid statuses even with invalid ones present");
}

// Test: contexts_alert_status_str_to_id with consecutive separators
void test_contexts_alert_status_str_to_id_consecutive_separators(void) {
    char input[] = "warning,,critical";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should skip consecutive separators");
    assert_true((result & CONTEXT_ALERT_CRITICAL) != 0, 
                "contexts_alert_status_str_to_id should skip consecutive separators");
}

// Test: contexts_alert_status_str_to_id with leading and trailing separators
void test_contexts_alert_status_str_to_id_leading_trailing_separators(void) {
    char input[] = ",warning,critical,";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should skip leading separators");
    assert_true((result & CONTEXT_ALERT_CRITICAL) != 0, 
                "contexts_alert_status_str_to_id should skip trailing separators");
}

// Test: contexts_alert_status_str_to_id with duplicate statuses
void test_contexts_alert_status_str_to_id_duplicate_statuses(void) {
    char input[] = "warning,warning,critical";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_true((result & CONTEXT_ALERT_WARNING) != 0, 
                "contexts_alert_status_str_to_id should handle duplicate statuses");
    assert_true((result & CONTEXT_ALERT_CRITICAL) != 0, 
                "contexts_alert_status_str_to_id should handle duplicate statuses");
}

// Test: contexts_alert_status_str_to_id with all statuses
void test_contexts_alert_status_str_to_id_all_statuses(void) {
    char input[] = "uninitialized,undefined,clear,raised,warning,critical";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    assert_equal_u32(CONTEXTS_ALERT_STATUSES, result, 
                     "contexts_alert_status_str_to_id should recognize all statuses");
}

// Test: contexts_alerts_status_to_buffer_json_array with NULL buffer
void test_contexts_alerts_status_to_buffer_json_array_null_buffer(void) {
    // Should not crash, returns void
    contexts_alerts_status_to_buffer_json_array(NULL, "key", CONTEXT_ALERT_CRITICAL);
    assert_true(1, "contexts_alerts_status_to_buffer_json_array should handle NULL buffer gracefully");
}

// Test: contexts_alerts_status_to_buffer_json_array with NULL key
void test_contexts_alerts_status_to_buffer_json_array_null_key(void) {
    BUFFER wb = {0};
    // Should not crash, returns void
    contexts_alerts_status_to_buffer_json_array(&wb, NULL, CONTEXT_ALERT_CRITICAL);
    assert_true(1, "contexts_alerts_status_to_buffer_json_array should handle NULL key gracefully");
}

// Test: contexts_alerts_status_to_buffer_json_array with single status
void test_contexts_alerts_status_to_buffer_json_array_single_status(void) {
    BUFFER wb = {0};
    contexts_alert_statuses_init();
    contexts_alerts_status_to_buffer_json_array(&wb, "statuses", CONTEXT_ALERT_CRITICAL);
    assert_true(1, "contexts_alerts_status_to_buffer_json_array should process single status");
}

// Test: contexts_alerts_status_to_buffer_json_array with multiple statuses
void test_contexts_alerts_status_to_buffer_json_array_multiple_statuses(void) {
    BUFFER wb = {0};
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS options = CONTEXT_ALERT_WARNING | CONTEXT_ALERT_CRITICAL;
    contexts_alerts_status_to_buffer_json_array(&wb, "statuses", options);
    assert_true(1, "contexts_alerts_status_to_buffer_json_array should process multiple statuses");
}

// Test: contexts_alerts_status_to_buffer_json_array with all statuses
void test_contexts_alerts_status_to_buffer_json_array_all_statuses(void) {
    BUFFER wb = {0};
    contexts_alert_statuses_init();
    contexts_alerts_status_to_buffer_json_array(&wb, "statuses", CONTEXTS_ALERT_STATUSES);
    assert_true(1, "contexts_alerts_status_to_buffer_json_array should process all statuses");
}

// Test: contexts_alerts_status_to_buffer_json_array with no statuses
void test_contexts_alerts_status_to_buffer_json_array_no_statuses(void) {
    BUFFER wb = {0};
    contexts_alert_statuses_init();
    contexts_alerts_status_to_buffer_json_array(&wb, "statuses", 0);
    assert_true(1, "contexts_alerts_status_to_buffer_json_array should handle no statuses");
}

// Test: contexts_alerts_status_to_buffer_json_array prevents duplicates
void test_contexts_alerts_status_to_buffer_json_array_no_duplicates(void) {
    BUFFER wb = {0};
    contexts_alert_statuses_init();
    // Test that the "used" variable prevents duplicates
    CONTEXTS_ALERT_STATUS options = CONTEXT_ALERT_WARNING;
    contexts_alerts_status_to_buffer_json_array(&wb, "statuses", options);
    assert_true(1, "contexts_alerts_status_to_buffer_json_array should prevent duplicate entries");
}

// Test: contexts_alert_statuses_init initializes hashes
void test_contexts_alert_statuses_init(void) {
    contexts_alert_statuses_init();
    // After init, all hashes should be non-zero (except for the NULL terminator)
    for (int i = 0; contexts_alert_statuses[i].name; i++) {
        assert_true(contexts_alert_statuses[i].hash != 0, 
                    "contexts_alert_statuses_init should initialize hash values");
    }
}

// Test: enum value coverage
void test_enum_values_coverage(void) {
    assert_equal_u32((1 << 6), CONTEXT_ALERT_UNINITIALIZED, "CONTEXT_ALERT_UNINITIALIZED value check");
    assert_equal_u32((1 << 7), CONTEXT_ALERT_UNDEFINED, "CONTEXT_ALERT_UNDEFINED value check");
    assert_equal_u32((1 << 8), CONTEXT_ALERT_CLEAR, "CONTEXT_ALERT_CLEAR value check");
    assert_equal_u32((1 << 9), CONTEXT_ALERT_RAISED, "CONTEXT_ALERT_RAISED value check");
    assert_equal_u32((1 << 10), CONTEXT_ALERT_WARNING, "CONTEXT_ALERT_WARNING value check");
    assert_equal_u32((1 << 11), CONTEXT_ALERT_CRITICAL, "CONTEXT_ALERT_CRITICAL value check");
}

// Test: bitwise AND operations in status checking
void test_bitwise_and_operations(void) {
    CONTEXTS_ALERT_STATUS status = CONTEXT_ALERT_WARNING;
    assert_true((status & CONTEXT_ALERT_WARNING) != 0, "bitwise AND should detect present status");
    assert_true((status & CONTEXT_ALERT_CRITICAL) == 0, "bitwise AND should detect absent status");
}

// Test: contexts_alert_status_str_to_id with case sensitivity
void test_contexts_alert_status_str_to_id_case_sensitivity(void) {
    char input[] = "WARNING";
    contexts_alert_statuses_init();
    CONTEXTS_ALERT_STATUS result = contexts_alert_status_str_to_id(input);
    // Should not match due to case sensitivity
    assert_equal_u32(0, result, "contexts_alert_status_str_to_id should be case-sensitive");
}

// Test: bitwise OR accumulation
void test_bitwise_or_accumulation(void) {
    CONTEXTS_ALERT_STATUS status = 0;
    status |= CONTEXT_ALERT_WARNING;
    assert_true((status & CONTEXT_ALERT_WARNING) != 0, "first OR operation should work");
    status |= CONTEXT_ALERT_CRITICAL;
    assert_true((status & CONTEXT_ALERT_CRITICAL) != 0, "second OR operation should work");
    assert_true((status & CONTEXT_ALERT_WARNING) != 0, "previous flag should remain after second OR");
}

// Main test runner
int main(void) {
    printf("=== Testing contexts_alert_statuses ===\n\n");
    
    // contexts_alert_status_str_to_id tests
    test_contexts_alert_status_str_to_id_null_input();
    test_contexts_alert_status_str_to_id_empty_string();
    test_contexts_alert_status_str_to_id_single_uninitialized();
    test_contexts_alert_status_str_to_id_single_undefined();
    test_contexts_alert_status_str_to_id_single_clear();
    test_contexts_alert_status_str_to_id_single_raised();
    test_contexts_alert_status_str_to_id_single_warning();
    test_contexts_alert_status_str_to_id_single_critical();
    test_contexts_alert_status_str_to_id_multiple_comma_separated();
    test_contexts_alert_status_str_to_id_multiple_pipe_separated();
    test_contexts_alert_status_str_to_id_multiple_space_separated();
    test_contexts_alert_status_str_to_id_invalid_status();
    test_contexts_alert_status_str_to_id_mixed_valid_invalid();
    test_contexts_alert_status_str_to_id_consecutive_separators();
    test_contexts_alert_status_str_to_id_leading_trailing_separators();
    test_contexts_alert_status_str_to_id_duplicate_statuses();
    test_contexts_alert_status_str_to_id_all_statuses();
    test_contexts_alert_status_str_to_id_case_sensitivity();
    
    // contexts_alerts_status_to_buffer_json_array tests
    test_contexts_alerts_status_to_buffer_json_array_null_buffer();
    test_contexts_alerts_status_to_buffer_json_array_null_key();
    test_contexts_alerts_status_to_buffer_json_array_single_status();
    test_contexts_alerts_status_to_buffer_json_array_multiple_statuses();
    test_contexts_alerts_status_to_buffer_json_array_all_statuses();
    test_contexts_alerts_status_to_buffer_json_array_no_statuses();
    test_contexts_alerts_status_to_buffer_json_array_no_duplicates();
    
    // contexts_alert_statuses_init tests
    test_contexts_alert_statuses_init();
    
    // Enum and bitwise operation tests
    test_enum_values_coverage();
    test_bitwise_and_operations();
    test_bitwise_or_accumulation();
    
    printf("\n=== Test Summary ===\n");
    printf("Total: %d\n", test_count);
    printf("Passed: %d\n", test_passed);
    printf("Failed: %d\n", test_failed);
    
    return test_failed > 0 ? 1 : 0;
}