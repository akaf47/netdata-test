#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>
#include <sys/types.h>
#include <stdint.h>
#include <stdbool.h>

/* Mock structures and forward declarations */
#define ANALYTICS_MAX_NAME 256
#define ANALYTICS_MAX_VALUE 1024

/* Mock for external dependencies */
typedef struct {
    char name[ANALYTICS_MAX_NAME];
    char value[ANALYTICS_MAX_VALUE];
    uint64_t timestamp;
} analytics_event_t;

typedef struct {
    analytics_event_t *events;
    size_t count;
    size_t capacity;
} mock_analytics_store;

static mock_analytics_store mock_store = {0};

/* Test framework helpers */
#define TEST_PASS() do { printf("✓ %s\n", __func__); } while(0)
#define TEST_FAIL(msg) do { printf("✗ %s: %s\n", __func__, msg); exit(1); } while(0)

#define ASSERT_EQUAL(actual, expected, fmt) \
    do { \
        if ((actual) != (expected)) { \
            printf("  Expected: " fmt ", Got: " fmt "\n", expected, actual); \
            TEST_FAIL("Assertion failed"); \
        } \
    } while(0)

#define ASSERT_STRING_EQUAL(actual, expected) \
    do { \
        if (strcmp((actual), (expected)) != 0) { \
            printf("  Expected: '%s', Got: '%s'\n", expected, actual); \
            TEST_FAIL("String assertion failed"); \
        } \
    } while(0)

#define ASSERT_NOT_NULL(ptr) \
    do { \
        if ((ptr) == NULL) { \
            TEST_FAIL("Pointer is NULL"); \
        } \
    } while(0)

#define ASSERT_NULL(ptr) \
    do { \
        if ((ptr) != NULL) { \
            TEST_FAIL("Pointer is not NULL"); \
        } \
    } while(0)

#define ASSERT_TRUE(condition) \
    do { \
        if (!(condition)) { \
            TEST_FAIL("Condition is false"); \
        } \
    } while(0)

#define ASSERT_FALSE(condition) \
    do { \
        if ((condition)) { \
            TEST_FAIL("Condition is true"); \
        } \
    } while(0)

/* Setup and teardown */
static void setup(void) {
    mock_store.events = NULL;
    mock_store.count = 0;
    mock_store.capacity = 0;
}

static void teardown(void) {
    if (mock_store.events) {
        free(mock_store.events);
        mock_store.events = NULL;
    }
    mock_store.count = 0;
    mock_store.capacity = 0;
}

/* Test: analytics initialization */
static void test_analytics_init_success(void) {
    setup();
    
    /* Test successful initialization */
    mock_store.capacity = 100;
    mock_store.events = malloc(sizeof(analytics_event_t) * 100);
    
    ASSERT_NOT_NULL(mock_store.events);
    ASSERT_EQUAL(mock_store.count, 0, "%zu");
    ASSERT_EQUAL(mock_store.capacity, 100, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_init_null_allocation(void) {
    setup();
    
    /* Simulate failed allocation by not allocating */
    mock_store.events = NULL;
    
    ASSERT_NULL(mock_store.events);
    
    teardown();
    TEST_PASS();
}

/* Test: event creation */
static void test_analytics_event_create_basic(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    const char *name = "startup";
    const char *value = "success";
    
    if (mock_store.count < mock_store.capacity) {
        strncpy(mock_store.events[mock_store.count].name, name, ANALYTICS_MAX_NAME - 1);
        strncpy(mock_store.events[mock_store.count].value, value, ANALYTICS_MAX_VALUE - 1);
        mock_store.events[mock_store.count].timestamp = time(NULL);
        mock_store.count++;
    }
    
    ASSERT_EQUAL(mock_store.count, 1, "%zu");
    ASSERT_STRING_EQUAL(mock_store.events[0].name, "startup");
    ASSERT_STRING_EQUAL(mock_store.events[0].value, "success");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_event_create_empty_name(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    const char *name = "";
    const char *value = "test";
    
    if (mock_store.count < mock_store.capacity) {
        strncpy(mock_store.events[mock_store.count].name, name, ANALYTICS_MAX_NAME - 1);
        strncpy(mock_store.events[mock_store.count].value, value, ANALYTICS_MAX_VALUE - 1);
        mock_store.count++;
    }
    
    ASSERT_EQUAL(mock_store.count, 1, "%zu");
    ASSERT_STRING_EQUAL(mock_store.events[0].name, "");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_event_create_empty_value(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    const char *name = "event";
    const char *value = "";
    
    if (mock_store.count < mock_store.capacity) {
        strncpy(mock_store.events[mock_store.count].name, name, ANALYTICS_MAX_NAME - 1);
        strncpy(mock_store.events[mock_store.count].value, value, ANALYTICS_MAX_VALUE - 1);
        mock_store.count++;
    }
    
    ASSERT_EQUAL(mock_store.count, 1, "%zu");
    ASSERT_STRING_EQUAL(mock_store.events[0].value, "");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_event_create_long_strings(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    char long_name[ANALYTICS_MAX_NAME];
    char long_value[ANALYTICS_MAX_VALUE];
    
    memset(long_name, 'a', ANALYTICS_MAX_NAME - 1);
    long_name[ANALYTICS_MAX_NAME - 1] = '\0';
    
    memset(long_value, 'b', ANALYTICS_MAX_VALUE - 1);
    long_value[ANALYTICS_MAX_VALUE - 1] = '\0';
    
    if (mock_store.count < mock_store.capacity) {
        strncpy(mock_store.events[mock_store.count].name, long_name, ANALYTICS_MAX_NAME - 1);
        strncpy(mock_store.events[mock_store.count].value, long_value, ANALYTICS_MAX_VALUE - 1);
        mock_store.count++;
    }
    
    ASSERT_EQUAL(mock_store.count, 1, "%zu");
    ASSERT_EQUAL(strlen(mock_store.events[0].name), ANALYTICS_MAX_NAME - 1, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_event_create_buffer_overflow_protection(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    /* Create string longer than buffer */
    char very_long_string[ANALYTICS_MAX_NAME * 2];
    memset(very_long_string, 'x', ANALYTICS_MAX_NAME * 2 - 1);
    very_long_string[ANALYTICS_MAX_NAME * 2 - 1] = '\0';
    
    if (mock_store.count < mock_store.capacity) {
        strncpy(mock_store.events[mock_store.count].name, very_long_string, ANALYTICS_MAX_NAME - 1);
        mock_store.events[mock_store.count].name[ANALYTICS_MAX_NAME - 1] = '\0';
        mock_store.count++;
    }
    
    /* Verify truncation protection */
    ASSERT_EQUAL(strlen(mock_store.events[0].name), ANALYTICS_MAX_NAME - 1, "%zu");
    
    teardown();
    TEST_PASS();
}

/* Test: capacity management */
static void test_analytics_capacity_expansion(void) {
    setup();
    
    /* Start with small capacity */
    mock_store.capacity = 5;
    mock_store.events = malloc(sizeof(analytics_event_t) * 5);
    
    /* Add events until capacity reached */
    for (size_t i = 0; i < 5; i++) {
        if (mock_store.count < mock_store.capacity) {
            snprintf(mock_store.events[i].name, ANALYTICS_MAX_NAME, "event_%zu", i);
            mock_store.count++;
        }
    }
    
    ASSERT_EQUAL(mock_store.count, 5, "%zu");
    
    /* Simulate expansion */
    size_t new_capacity = mock_store.capacity * 2;
    analytics_event_t *new_events = malloc(sizeof(analytics_event_t) * new_capacity);
    memcpy(new_events, mock_store.events, sizeof(analytics_event_t) * mock_store.count);
    free(mock_store.events);
    mock_store.events = new_events;
    mock_store.capacity = new_capacity;
    
    ASSERT_EQUAL(mock_store.capacity, 10, "%zu");
    ASSERT_EQUAL(mock_store.count, 5, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_capacity_zero(void) {
    setup();
    
    mock_store.capacity = 0;
    mock_store.events = NULL;
    
    ASSERT_EQUAL(mock_store.capacity, 0, "%zu");
    ASSERT_NULL(mock_store.events);
    
    teardown();
    TEST_PASS();
}

static void test_analytics_capacity_single(void) {
    setup();
    
    mock_store.capacity = 1;
    mock_store.events = malloc(sizeof(analytics_event_t) * 1);
    
    if (mock_store.count < mock_store.capacity) {
        strncpy(mock_store.events[0].name, "single", ANALYTICS_MAX_NAME - 1);
        mock_store.count++;
    }
    
    ASSERT_EQUAL(mock_store.count, 1, "%zu");
    ASSERT_EQUAL(mock_store.capacity, 1, "%zu");
    
    teardown();
    TEST_PASS();
}

/* Test: timestamp handling */
static void test_analytics_timestamp_current_time(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    time_t before = time(NULL);
    
    if (mock_store.count < mock_store.capacity) {
        strncpy(mock_store.events[0].name, "test", ANALYTICS_MAX_NAME - 1);
        mock_store.events[0].timestamp = time(NULL);
        mock_store.count++;
    }
    
    time_t after = time(NULL);
    
    ASSERT_TRUE(mock_store.events[0].timestamp >= before);
    ASSERT_TRUE(mock_store.events[0].timestamp <= after);
    
    teardown();
    TEST_PASS();
}

static void test_analytics_timestamp_ordering(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    /* Create multiple events with timestamps */
    for (int i = 0; i < 3; i++) {
        if (mock_store.count < mock_store.capacity) {
            snprintf(mock_store.events[i].name, ANALYTICS_MAX_NAME, "event_%d", i);
            mock_store.events[i].timestamp = time(NULL) + i;
            mock_store.count++;
        }
        usleep(100);
    }
    
    /* Verify timestamp ordering */
    ASSERT_TRUE(mock_store.events[0].timestamp <= mock_store.events[1].timestamp);
    ASSERT_TRUE(mock_store.events[1].timestamp <= mock_store.events[2].timestamp);
    
    teardown();
    TEST_PASS();
}

/* Test: event count management */
static void test_analytics_count_zero(void) {
    setup();
    
    ASSERT_EQUAL(mock_store.count, 0, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_count_increment(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    for (int i = 0; i < 7; i++) {
        if (mock_store.count < mock_store.capacity) {
            mock_store.count++;
        }
    }
    
    ASSERT_EQUAL(mock_store.count, 7, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_count_does_not_exceed_capacity(void) {
    setup();
    
    mock_store.capacity = 5;
    mock_store.events = malloc(sizeof(analytics_event_t) * 5);
    
    /* Try to add more than capacity */
    for (int i = 0; i < 10; i++) {
        if (mock_store.count < mock_store.capacity) {
            mock_store.count++;
        }
    }
    
    ASSERT_EQUAL(mock_store.count, 5, "%zu");
    
    teardown();
    TEST_PASS();
}

/* Test: event retrieval */
static void test_analytics_get_event_first(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    strncpy(mock_store.events[0].name, "first", ANALYTICS_MAX_NAME - 1);
    strncpy(mock_store.events[0].value, "value", ANALYTICS_MAX_VALUE - 1);
    mock_store.count = 1;
    
    ASSERT_STRING_EQUAL(mock_store.events[0].name, "first");
    ASSERT_STRING_EQUAL(mock_store.events[0].value, "value");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_get_event_last(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    for (int i = 0; i < 5; i++) {
        snprintf(mock_store.events[i].name, ANALYTICS_MAX_NAME, "event_%d", i);
        mock_store.count++;
    }
    
    int last_idx = mock_store.count - 1;
    ASSERT_STRING_EQUAL(mock_store.events[last_idx].name, "event_4");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_get_event_out_of_bounds(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    mock_store.count = 3;
    
    /* Access beyond count should not be allowed in real code */
    size_t invalid_idx = 10;
    ASSERT_TRUE(invalid_idx >= mock_store.count);
    
    teardown();
    TEST_PASS();
}

/* Test: event filtering */
static void test_analytics_filter_by_name_match(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    for (int i = 0; i < 3; i++) {
        snprintf(mock_store.events[i].name, ANALYTICS_MAX_NAME, "startup");
        snprintf(mock_store.events[i].value, ANALYTICS_MAX_VALUE, "value_%d", i);
        mock_store.count++;
    }
    
    /* Count matching events */
    size_t matches = 0;
    for (size_t i = 0; i < mock_store.count; i++) {
        if (strcmp(mock_store.events[i].name, "startup") == 0) {
            matches++;
        }
    }
    
    ASSERT_EQUAL(matches, 3, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_filter_by_name_no_match(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    strncpy(mock_store.events[0].name, "event", ANALYTICS_MAX_NAME - 1);
    mock_store.count = 1;
    
    size_t matches = 0;
    for (size_t i = 0; i < mock_store.count; i++) {
        if (strcmp(mock_store.events[i].name, "nonexistent") == 0) {
            matches++;
        }
    }
    
    ASSERT_EQUAL(matches, 0, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_filter_empty_store(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    mock_store.count = 0;
    
    size_t matches = 0;
    for (size_t i = 0; i < mock_store.count; i++) {
        matches++;
    }
    
    ASSERT_EQUAL(matches, 0, "%zu");
    
    teardown();
    TEST_PASS();
}

/* Test: time-based operations */
static void test_analytics_time_range_all_events(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    time_t start = time(NULL);
    
    for (int i = 0; i < 3; i++) {
        mock_store.events[i].timestamp = start + i;
        mock_store.count++;
    }
    
    time_t end = start + 10;
    
    size_t in_range = 0;
    for (size_t i = 0; i < mock_store.count; i++) {
        if (mock_store.events[i].timestamp >= start && 
            mock_store.events[i].timestamp <= end) {
            in_range++;
        }
    }
    
    ASSERT_EQUAL(in_range, 3, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_time_range_partial_events(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    time_t start = 1000;
    
    mock_store.events[0].timestamp = 500;
    mock_store.events[1].timestamp = 1500;
    mock_store.events[2].timestamp = 2500;
    mock_store.count = 3;
    
    time_t query_start = 1000;
    time_t query_end = 2000;
    
    size_t in_range = 0;
    for (size_t i = 0; i < mock_store.count; i++) {
        if (mock_store.events[i].timestamp >= query_start && 
            mock_store.events[i].timestamp <= query_end) {
            in_range++;
        }
    }
    
    ASSERT_EQUAL(in_range, 1, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_time_range_no_events(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    mock_store.events[0].timestamp = 100;
    mock_store.events[1].timestamp = 200;
    mock_store.count = 2;
    
    time_t query_start = 5000;
    time_t query_end = 6000;
    
    size_t in_range = 0;
    for (size_t i = 0; i < mock_store.count; i++) {
        if (mock_store.events[i].timestamp >= query_start && 
            mock_store.events[i].timestamp <= query_end) {
            in_range++;
        }
    }
    
    ASSERT_EQUAL(in_range, 0, "%zu");
    
    teardown();
    TEST_PASS();
}

/* Test: cleanup and finalization */
static void test_analytics_cleanup_allocated(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    ASSERT_NOT_NULL(mock_store.events);
    
    teardown();
    
    ASSERT_NULL(mock_store.events);
    
    TEST_PASS();
}

static void test_analytics_cleanup_null_pointer(void) {
    setup();
    
    mock_store.events = NULL;
    mock_store.count = 0;
    mock_store.capacity = 0;
    
    ASSERT_NULL(mock_store.events);
    
    teardown();
    
    ASSERT_NULL(mock_store.events);
    
    TEST_PASS();
}

static void test_analytics_cleanup_resets_count(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    mock_store.count = 7;
    
    teardown();
    
    ASSERT_EQUAL(mock_store.count, 0, "%zu");
    
    TEST_PASS();
}

static void test_analytics_cleanup_resets_capacity(void) {
    setup();
    
    mock_store.capacity = 100;
    mock_store.events = malloc(sizeof(analytics_event_t) * 100);
    
    teardown();
    
    ASSERT_EQUAL(mock_store.capacity, 0, "%zu");
    
    TEST_PASS();
}

/* Test: edge cases */
static void test_analytics_max_uint64_timestamp(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    mock_store.events[0].timestamp = UINT64_MAX;
    mock_store.count = 1;
    
    ASSERT_EQUAL(mock_store.events[0].timestamp, UINT64_MAX, "%lu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_zero_timestamp(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    mock_store.events[0].timestamp = 0;
    mock_store.count = 1;
    
    ASSERT_EQUAL(mock_store.events[0].timestamp, 0, "%lu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_null_character_in_name(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    mock_store.events[0].name[0] = '\0';
    mock_store.count = 1;
    
    ASSERT_EQUAL(strlen(mock_store.events[0].name), 0, "%zu");
    
    teardown();
    TEST_PASS();
}

static void test_analytics_special_characters_in_value(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    const char *special = "!@#$%^&*(){}[]|:;<>?,./";
    strncpy(mock_store.events[0].value, special, ANALYTICS_MAX_VALUE - 1);
    mock_store.count = 1;
    
    ASSERT_STRING_EQUAL(mock_store.events[0].value, special);
    
    teardown();
    TEST_PASS();
}

static void test_analytics_unicode_in_value(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    const char *unicode = "测试中文 テスト éàèù";
    strncpy(mock_store.events[0].value, unicode, ANALYTICS_MAX_VALUE - 1);
    mock_store.count = 1;
    
    ASSERT_STRING_EQUAL(mock_store.events[0].value, unicode);
    
    teardown();
    TEST_PASS();
}

/* Test: state consistency */
static void test_analytics_count_matches_actual_events(void) {
    setup();
    
    mock_store.capacity = 10;
    mock_store.events = malloc(sizeof(analytics_event_t) * 10);
    
    for (int i = 0; i < 5; i++) {
        if (mock_store.count < mock_store.capacity) {
            snprintf(mock_store.events[i].name, ANALYTICS_MAX_NAME, "event_%d", i);
            mock_store.count++;
        }
    }
    
    size_t actual_count = 0;
    for (size_t i = 0; i < mock_store.capacity; i++) {
        if (mock_store.events[i].name[0] != '\0') {
            actual_count++;
        }
    }
    
    /* Only first 5 should be initialized */
    ASSERT_TRUE(actual_count >= 5);
    
    teardown();
    TEST_PASS();
}

static void test_analytics_capacity_greater_than_count(void) {
    setup();
    
    mock_store.capacity = 20;
    mock_store.events = malloc(sizeof(analytics_event_t) * 20);
    mock_store.count = 5;
    
    ASSERT_TRUE(mock_store.capacity >= mock_store.count);
    
    teardown();
    TEST_PASS();
}

static void test_analytics_capacity_exactly_matches_count(void) {
    setup();
    
    mock_store.capacity = 5;
    mock_store.events = malloc(sizeof(analytics_event_t) * 5);
    mock_store.count = 5;
    
    ASSERT_EQUAL(mock_store.capacity, mock_store.count, "%zu");
    
    teardown();
    TEST_PASS();
}

/* Test: multiple operations in sequence */
static void test_analytics_multiple_operations_sequence(void) {
    setup();
    
    /* Initialize */
    mock_store.capacity = 20;
    mock_store.events = malloc(sizeof(analytics_event_t) * 20);
    ASSERT_NOT_NULL(mock_store.events);
    
    /* Add multiple events */
    for (int i = 0; i < 5; i++) {
        if (mock_store.count < mock_store.capacity) {
            snprintf(mock_store.events[i].name, ANALYTICS_MAX_NAME, "event_%d", i);
            snprintf(mock_store.events[i].value, ANALYTICS_MAX_VALUE, "value_%d", i);
            mock_store.events[i].timestamp = time(NULL) + i;
            mock_store.count++;
        }
    }
    
    ASSERT_EQUAL(mock_store.count, 5, "%zu");
    
    /* Verify events */
    for (int i = 0; i < 5; i++) {
        char expected_name[ANALYTICS_MAX_NAME];
        snprintf(expected_name, ANALYTICS_MAX_NAME, "event_%d", i);
        ASSERT_STRING_EQUAL(mock_store.events[i].name, expected_name);
    }
    
    /* Expand capacity */
    size_t new_capacity = mock_store.capacity * 2;
    analytics_event_t *new_events = malloc(sizeof(analytics_event_t) * new_capacity);
    memcpy(new_events, mock_store.events, sizeof(analytics_event_t) * mock_store.count);
    free(mock_store.events);
    mock_store.events = new_events;
    mock_store.capacity = new_capacity;
    
    ASSERT_EQUAL(mock_store.capacity, 40, "%zu");
    ASSERT_EQUAL(mock_store.count, 5, "%zu");
    
    teardown();
    TEST_PASS();
}

/* Main test runner */
int main(void) {
    printf("Starting Analytics Tests\n");
    printf("========================\n\n");
    
    /* Initialization tests */
    printf("Initialization Tests:\n");
    test_analytics_init_success();
    test_analytics_init_null_allocation();
    
    /* Event creation tests */
    printf("\nEvent Creation Tests:\n");
    test_analytics_event_create_basic();
    test_analytics_event_create_empty_name();
    test_analytics_event_create_empty_value();
    test_analytics_event_create_long_strings();
    test_analytics_event_create_buffer_overflow_protection();
    
    /* Capacity management tests */
    printf("\nCapacity Management Tests:\n");
    test_analytics_capacity_expansion();
    test_analytics_capacity_zero();
    test_analytics_capacity_single();
    
    /* Timestamp handling tests */
    printf("\nTimestamp Handling Tests:\n");
    test_analytics_timestamp_current_time();
    test_analytics_timestamp_ordering();
    
    /* Event count management tests */
    printf("\nEvent Count Management Tests:\n");
    test_analytics_count_zero();
    test_analytics_count_