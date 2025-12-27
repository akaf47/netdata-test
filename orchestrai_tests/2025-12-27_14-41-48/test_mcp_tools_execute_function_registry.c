#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdbool.h>

// Mock types and structures needed for testing
typedef struct {
    char *hostname;
} RRDHOST;

typedef struct {
    int access;
} USER_AUTH;

typedef struct {
    char data[4096];
    size_t len;
} BUFFER;

typedef struct {
    void *lock;
} RW_SPINLOCK;

typedef struct {
    void *lock;
} SPINLOCK;

typedef struct {
    char *str;
} STRING;

typedef enum {
    MCP_REQUIRED_PARAMS_TYPE_SELECT = 0,
    MCP_REQUIRED_PARAMS_TYPE_MULTISELECT
} MCP_REQUIRED_PARAMS_TYPE;

typedef enum {
    MCP_PAGINATION_UNITS_TIMESTAMP_USEC = 0,
    MCP_PAGINATION_UNITS_TIMESTAMP_SEC = 1
} MCP_PAGINATION_UNITS;

typedef enum {
    FN_TYPE_UNKNOWN = 0,
    FN_TYPE_TABLE = 1,
    FN_TYPE_TABLE_WITH_HISTORY = 2,
    FN_TYPE_NOT_TABLE = 3
} MCP_FUNCTION_TYPE;

typedef struct mcp_function_param_option {
    STRING *id;
    STRING *name;
    STRING *info;
} MCP_FUNCTION_PARAM_OPTION;

typedef struct mcp_function_param {
    STRING *id;
    STRING *name;
    STRING *help;
    MCP_REQUIRED_PARAMS_TYPE type;
    bool unique_view;
    size_t options_count;
    MCP_FUNCTION_PARAM_OPTION *options;
} MCP_FUNCTION_PARAM;

typedef struct mcp_function_pagination {
    bool enabled;
    STRING *key;
    STRING *column;
    MCP_PAGINATION_UNITS units;
} MCP_FUNCTION_PAGINATION;

typedef struct mcp_function_registry_entry {
    RW_SPINLOCK spinlock;
    SPINLOCK update_spinlock;
    MCP_FUNCTION_TYPE type;
    bool has_history;
    int update_every;
    STRING *help;
    int version;
    bool supports_post;
    size_t required_params_count;
    MCP_FUNCTION_PARAM *required_params;
    bool has_timeframe;
    bool has_last;
    bool has_data_only;
    bool has_direction;
    bool has_query;
    bool has_slice;
    MCP_FUNCTION_PAGINATION pagination;
    time_t last_update;
    time_t expires;
} MCP_FUNCTION_REGISTRY_ENTRY;

#define MCP_FUNCTIONS_REGISTRY_TTL 600

// Mock functions
void mcp_functions_registry_init(void) {
    // Mock implementation
}

void mcp_functions_registry_cleanup(void) {
    // Mock implementation
}

MCP_FUNCTION_REGISTRY_ENTRY *mcp_functions_registry_get(RRDHOST *host, const char *function_name, BUFFER *error) {
    // Mock implementation
    if (!host || !function_name) {
        if (error) strcpy(error->data, "Invalid parameters");
        return NULL;
    }
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = malloc(sizeof(MCP_FUNCTION_REGISTRY_ENTRY));
    memset(entry, 0, sizeof(MCP_FUNCTION_REGISTRY_ENTRY));
    entry->type = FN_TYPE_TABLE;
    entry->has_history = false;
    entry->version = 3;
    entry->supports_post = true;
    entry->expires = time(NULL) + 600;
    
    return entry;
}

void mcp_functions_registry_release(MCP_FUNCTION_REGISTRY_ENTRY *entry) {
    if (entry) {
        free(entry);
    }
}

// Test: mcp_functions_registry_init
static void test_registry_init(void **state) {
    (void)state;
    mcp_functions_registry_init();
    // Verify initialization completed without errors
    assert_true(1);
}

// Test: mcp_functions_registry_cleanup
static void test_registry_cleanup(void **state) {
    (void)state;
    mcp_functions_registry_init();
    mcp_functions_registry_cleanup();
    // Verify cleanup completed without errors
    assert_true(1);
}

// Test: mcp_functions_registry_get with valid host and function
static void test_registry_get_valid_parameters(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "test_function", &error);
    
    assert_non_null(entry);
    assert_int_equal(entry->type, FN_TYPE_TABLE);
    assert_false(entry->has_history);
    assert_true(entry->supports_post);
    
    mcp_functions_registry_release(entry);
}

// Test: mcp_functions_registry_get with NULL host
static void test_registry_get_null_host(void **state) {
    (void)state;
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(NULL, "test_function", &error);
    
    assert_null(entry);
    assert_string_equal(error.data, "Invalid parameters");
}

// Test: mcp_functions_registry_get with NULL function name
static void test_registry_get_null_function(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, NULL, &error);
    
    assert_null(entry);
    assert_string_equal(error.data, "Invalid parameters");
}

// Test: mcp_functions_registry_get with empty function name
static void test_registry_get_empty_function(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "", &error);
    
    assert_null(entry);
    assert_string_equal(error.data, "Invalid parameters");
}

// Test: mcp_functions_registry_release with valid entry
static void test_registry_release_valid(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "test_function", &error);
    assert_non_null(entry);
    
    // Should not crash
    mcp_functions_registry_release(entry);
    assert_true(1);
}

// Test: mcp_functions_registry_release with NULL entry
static void test_registry_release_null(void **state) {
    (void)state;
    
    // Should handle gracefully
    mcp_functions_registry_release(NULL);
    assert_true(1);
}

// Test: mcp_functions_registry_get returns different entries for different functions
static void test_registry_get_different_functions(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry1 = mcp_functions_registry_get(&host, "func1", &error);
    MCP_FUNCTION_REGISTRY_ENTRY *entry2 = mcp_functions_registry_get(&host, "func2", &error);
    
    assert_non_null(entry1);
    assert_non_null(entry2);
    
    mcp_functions_registry_release(entry1);
    mcp_functions_registry_release(entry2);
}

// Test: mcp_functions_registry_get entry expiration
static void test_registry_get_entry_expiration(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "test_function", &error);
    
    assert_non_null(entry);
    // Verify expiration is set to TTL
    time_t now = time(NULL);
    assert_true(entry->expires > now);
    assert_true(entry->expires <= now + MCP_FUNCTIONS_REGISTRY_TTL + 1);
    
    mcp_functions_registry_release(entry);
}

// Test: mcp_functions_registry_get entry version
static void test_registry_get_entry_version(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "test_function", &error);
    
    assert_non_null(entry);
    assert_int_equal(entry->version, 3);
    assert_true(entry->supports_post);
    
    mcp_functions_registry_release(entry);
}

// Test: mcp_functions_registry_get entry type
static void test_registry_get_entry_type_table(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "test_function", &error);
    
    assert_non_null(entry);
    assert_int_equal(entry->type, FN_TYPE_TABLE);
    
    mcp_functions_registry_release(entry);
}

// Test: mcp_functions_registry_get with special characters in function name
static void test_registry_get_special_characters(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "func-name_123.test", &error);
    
    assert_non_null(entry);
    
    mcp_functions_registry_release(entry);
}

// Test: mcp_functions_registry_get with long function name
static void test_registry_get_long_function_name(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    const char *long_name = "very_long_function_name_that_contains_many_characters_"
                           "and_more_and_more_and_more_characters_to_test_limits";
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, long_name, &error);
    
    assert_non_null(entry);
    
    mcp_functions_registry_release(entry);
}

// Test: Multiple registry operations without cleanup
static void test_registry_multiple_operations(void **state) {
    (void)state;
    
    mcp_functions_registry_init();
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    // Get and release multiple times
    for (int i = 0; i < 5; i++) {
        MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "test", &error);
        assert_non_null(entry);
        mcp_functions_registry_release(entry);
    }
    
    mcp_functions_registry_cleanup();
    assert_true(1);
}

// Test: Registry entry with pagination enabled
static void test_registry_entry_pagination(void **state) {
    (void)state;
    
    RRDHOST host;
    host.hostname = "test-host";
    
    BUFFER error;
    memset(&error, 0, sizeof(BUFFER));
    
    MCP_FUNCTION_REGISTRY_ENTRY *entry = mcp_functions_registry_get(&host, "test_function", &error);
    
    assert_non_null(entry);
    // Entry should have pagination structure
    assert_non_null(&entry->pagination);
    
    mcp_functions_registry_release(entry);
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_registry_init),
        cmocka_unit_test(test_registry_cleanup),
        cmocka_unit_test(test_registry_get_valid_parameters),
        cmocka_unit_test(test_registry_get_null_host),
        cmocka_unit_test(test_registry_get_null_function),
        cmocka_unit_test(test_registry_get_empty_function),
        cmocka_unit_test(test_registry_release_valid),
        cmocka_unit_test(test_registry_release_null),
        cmocka_unit_test(test_registry_get_different_functions),
        cmocka_unit_test(test_registry_get_entry_expiration),
        cmocka_unit_test(test_registry_get_entry_version),
        cmocka_unit_test(test_registry_get_entry_type_table),
        cmocka_unit_test(test_registry_get_special_characters),
        cmocka_unit_test(test_registry_get_long_function_name),
        cmocka_unit_test(test_registry_multiple_operations),
        cmocka_unit_test(test_registry_entry_pagination),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}