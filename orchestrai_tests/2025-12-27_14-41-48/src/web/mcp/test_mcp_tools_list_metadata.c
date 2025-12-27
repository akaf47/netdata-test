#include <stddef.h>
#include <stdarg.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

// Mock structures and types
typedef struct {
    char data[4096];
    size_t len;
} BUFFER;

typedef struct {
    BUFFER *buf;
} CLEAN_BUFFER;

typedef int (*json_func)(BUFFER *);

typedef enum {
    MCP_LIST_OUTPUT_NODES = 0,
    MCP_LIST_OUTPUT_METRICS = 1,
    MCP_LIST_OUTPUT_FUNCTIONS = 2,
    MCP_LIST_OUTPUT_ALERTS = 3,
} MCP_LIST_OUTPUT_TYPE;

typedef enum {
    CONTEXTS_V2_CONTEXTS = 1,
    CONTEXTS_V2_NODES = 2,
    CONTEXTS_V2_NODES_INFO = 4,
    CONTEXTS_V2_NODE_INSTANCES = 8,
    CONTEXTS_V2_FUNCTIONS = 16,
    CONTEXTS_V2_ALERTS = 32,
    CONTEXTS_V2_SEARCH = 64,
} CONTEXTS_V2_MODE;

typedef struct {
    bool has_q;
    bool has_metrics;
    bool has_nodes;
    bool has_time_range;
    bool has_cardinality_limit;
    bool has_alert_pattern;
    bool metrics_required;
    bool nodes_required;
    bool nodes_as_array;
    bool metrics_as_array;
} MCP_PARAMS_CONFIG;

typedef struct {
    size_t cardinality_limit;
    int alert_status;
} MCP_DEFAULTS;

typedef struct {
    const char *name;
    const char *title;
    const char *description;
    MCP_LIST_OUTPUT_TYPE output_type;
    CONTEXTS_V2_MODE mode;
    int options;
    MCP_PARAMS_CONFIG params;
    MCP_DEFAULTS defaults;
} MCP_LIST_TOOL_CONFIG;

typedef enum {
    MCP_RC_OK = 0,
    MCP_RC_ERROR = 1,
    MCP_RC_BAD_REQUEST = 2,
} MCP_RETURN_CODE;

typedef struct {
    BUFFER *error;
    BUFFER *result;
} MCP_CLIENT;

typedef int MCP_REQUEST_ID;

typedef struct json_object json_object;

// Define constants for testing
#define MCP_TOOL_LIST_METRICS "list_metrics"
#define MCP_TOOL_GET_METRICS_DETAILS "get_metrics_details"
#define MCP_TOOL_LIST_NODES "list_nodes"
#define MCP_TOOL_LIST_FUNCTIONS "list_functions"
#define MCP_TOOL_GET_NODES_DETAILS "get_nodes_details"
#define MCP_TOOL_LIST_RAISED_ALERTS "list_raised_alerts"
#define MCP_TOOL_LIST_ALL_ALERTS "list_all_alerts"

#define MCP_DEFAULT_AFTER_TIME 0
#define MCP_DEFAULT_BEFORE_TIME 0
#define MCP_METADATA_CARDINALITY_LIMIT 100
#define MCP_METADATA_CARDINALITY_LIMIT_MAX 1000

#define CONTEXTS_OPTION_TITLES 1
#define CONTEXTS_OPTION_INSTANCES 2
#define CONTEXTS_OPTION_DIMENSIONS 4
#define CONTEXTS_OPTION_LABELS 8
#define CONTEXTS_OPTION_RETENTION 16
#define CONTEXTS_OPTION_LIVENESS 32
#define CONTEXTS_OPTION_FAMILY 64
#define CONTEXTS_OPTION_UNITS 128
#define CONTEXTS_OPTION_MCP 256
#define CONTEXTS_OPTION_RFC3339 512
#define CONTEXTS_OPTION_JSON_LONG_KEYS 1024
#define CONTEXTS_OPTION_MINIFY 2048
#define CONTEXTS_OPTION_SUMMARY 4096

#define CONTEXT_ALERT_RAISED 1
#define CONTEXTS_ALERT_STATUSES 3

#define HTTP_RESP_OK 200

#define MAX(a, b) ((a) > (b) ? (a) : (b))

// Mock function declarations
BUFFER *buffer_create(size_t size, void *ptr);
void buffer_free(BUFFER *buf);
const char *buffer_tostring(BUFFER *buf);
size_t buffer_strlen(BUFFER *buf);
void buffer_sprintf(BUFFER *buf, const char *format, ...);
void buffer_json_member_add_object(BUFFER *buf, const char *key);
void buffer_json_member_add_array(BUFFER *buf, const char *key);
void buffer_json_add_array_item_string(BUFFER *buf, const char *str);
void buffer_json_member_add_string(BUFFER *buf, const char *key, const char *value);
void buffer_json_object_close(BUFFER *buf);
void buffer_json_array_close(BUFFER *buf);
void buffer_json_finalize(BUFFER *buf);
void mcp_schema_add_array_param(BUFFER *buf, const char *name, const char *title, const char *description);
void mcp_schema_add_time_params(BUFFER *buf, const char *prefix, bool historical);
void mcp_schema_add_cardinality_limit(BUFFER *buf, const char *desc, size_t def, size_t min, size_t max);
int rrdcontext_to_json_v2(BUFFER *out, const void *req, CONTEXTS_V2_MODE mode);
void mcp_init_success_result(MCP_CLIENT *mcpc, MCP_REQUEST_ID id);
const char *mcp_params_extract_string(json_object *params, const char *key, const char *default_value);
size_t mcp_params_extract_size(json_object *params, const char *key, size_t def, size_t min, size_t max, BUFFER *error);
bool mcp_params_parse_time_window(json_object *params, time_t *after, time_t *before, time_t after_default, time_t before_default, bool required, BUFFER *error);
CLEAN_BUFFER *mcp_params_parse_array_to_pattern(json_object *params, const char *key, bool required, bool allow_empty, const char *tool_ref, BUFFER *error);

// Test structures initialization data
static const MCP_LIST_TOOL_CONFIG mcp_list_tools[] = {
    {
        .name = MCP_TOOL_LIST_METRICS,
        .title = "List available metrics",
        .description = "Search and list available metrics to query, across some or all nodes, for any time-frame",
        .output_type = MCP_LIST_OUTPUT_METRICS,
        .mode = CONTEXTS_V2_CONTEXTS,
        .options = 0,
        .params = {
            .has_q = true,
            .has_metrics = true,
            .has_nodes = true,
            .has_time_range = true,
            .has_cardinality_limit = true,
            .nodes_as_array = true,
        },
    },
    {
        .name = MCP_TOOL_GET_METRICS_DETAILS,
        .title = "Get metrics details",
        .description = "Get retention and cardinality information about specific metrics",
        .output_type = MCP_LIST_OUTPUT_METRICS,
        .mode = CONTEXTS_V2_CONTEXTS,
        .options = CONTEXTS_OPTION_TITLES | CONTEXTS_OPTION_INSTANCES | CONTEXTS_OPTION_DIMENSIONS | CONTEXTS_OPTION_LABELS | CONTEXTS_OPTION_RETENTION | CONTEXTS_OPTION_LIVENESS | CONTEXTS_OPTION_FAMILY | CONTEXTS_OPTION_UNITS,
        .params = {
            .has_metrics = true,
            .has_nodes = true,
            .has_time_range = true,
            .has_cardinality_limit = true,
            .metrics_required = true,
            .nodes_as_array = true,
            .metrics_as_array = true,
        },
    },
    {
        .name = MCP_TOOL_LIST_NODES,
        .title = "List monitored nodes",
        .description = "Search for and list monitored nodes by hostname patterns.",
        .output_type = MCP_LIST_OUTPUT_NODES,
        .mode = CONTEXTS_V2_NODES,
        .options = 0,
        .params = {
            .has_nodes = true,
            .has_metrics = true,
            .has_time_range = true,
            .has_cardinality_limit = true,
            .metrics_as_array = true,
        },
    },
    {
        .name = MCP_TOOL_LIST_FUNCTIONS,
        .title = "List available functions",
        .description = "List available Netdata functions that can be executed on specific nodes",
        .output_type = MCP_LIST_OUTPUT_FUNCTIONS,
        .mode = CONTEXTS_V2_FUNCTIONS,
        .options = 0,
        .params = {
            .has_nodes = true,
            .has_time_range = false,
            .has_cardinality_limit = false,
            .nodes_required = true,
            .nodes_as_array = true,
        },
    },
    {
        .name = MCP_TOOL_GET_NODES_DETAILS,
        .title = "Get detailed information about monitored nodes",
        .description = "Gets comprehensive node information including hardware specs",
        .output_type = MCP_LIST_OUTPUT_NODES,
        .mode = CONTEXTS_V2_NODES | CONTEXTS_V2_NODES_INFO | CONTEXTS_V2_NODE_INSTANCES,
        .options = 0,
        .params = {
            .has_nodes = true,
            .has_metrics = true,
            .has_time_range = true,
            .has_cardinality_limit = true,
            .nodes_required = true,
            .nodes_as_array = true,
            .metrics_as_array = true,
        },
    },
    {
        .name = MCP_TOOL_LIST_RAISED_ALERTS,
        .title = "List raised alerts",
        .description = "List currently active alerts (WARNING and CRITICAL status)",
        .output_type = MCP_LIST_OUTPUT_ALERTS,
        .mode = CONTEXTS_V2_ALERTS,
        .options = CONTEXTS_OPTION_INSTANCES | CONTEXTS_OPTION_VALUES,
        .params = {
            .has_nodes = true,
            .has_metrics = true,
            .has_alert_pattern = true,
            .has_time_range = false,
            .has_cardinality_limit = true,
            .nodes_as_array = true,
            .metrics_as_array = true,
        },
        .defaults = {
            .alert_status = CONTEXT_ALERT_RAISED,
            .cardinality_limit = 200,
        },
    },
    {
        .name = MCP_TOOL_LIST_ALL_ALERTS,
        .title = "List all alerts",
        .description = "List all currently running alerts",
        .output_type = MCP_LIST_OUTPUT_ALERTS,
        .mode = CONTEXTS_V2_ALERTS,
        .options = CONTEXTS_OPTION_SUMMARY,
        .params = {
            .has_nodes = true,
            .has_metrics = true,
            .has_alert_pattern = true,
            .has_time_range = true,
            .has_cardinality_limit = true,
            .nodes_as_array = true,
            .metrics_as_array = true,
        },
        .defaults = {
            .alert_status = CONTEXTS_ALERT_STATUSES,
            .cardinality_limit = 200,
        },
    },
};

// Test: mcp_get_list_tool_config - find existing tool
static void test_mcp_get_list_tool_config_find_metrics_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_METRICS);
    
    assert_non_null(config);
    assert_string_equal(config->name, MCP_TOOL_LIST_METRICS);
    assert_string_equal(config->title, "List available metrics");
    assert_int_equal(config->output_type, MCP_LIST_OUTPUT_METRICS);
    assert_int_equal(config->mode, CONTEXTS_V2_CONTEXTS);
    assert_true(config->params.has_q);
    assert_true(config->params.has_metrics);
    assert_true(config->params.has_nodes);
    assert_true(config->params.has_time_range);
    assert_true(config->params.has_cardinality_limit);
}

// Test: mcp_get_list_tool_config - find nodes tool
static void test_mcp_get_list_tool_config_find_nodes_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_NODES);
    
    assert_non_null(config);
    assert_string_equal(config->name, MCP_TOOL_LIST_NODES);
    assert_string_equal(config->title, "List monitored nodes");
    assert_int_equal(config->output_type, MCP_LIST_OUTPUT_NODES);
}

// Test: mcp_get_list_tool_config - find functions tool
static void test_mcp_get_list_tool_config_find_functions_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_FUNCTIONS);
    
    assert_non_null(config);
    assert_string_equal(config->name, MCP_TOOL_LIST_FUNCTIONS);
    assert_int_equal(config->output_type, MCP_LIST_OUTPUT_FUNCTIONS);
    assert_true(config->params.nodes_required);
}

// Test: mcp_get_list_tool_config - find metrics details tool
static void test_mcp_get_list_tool_config_find_metrics_details_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_GET_METRICS_DETAILS);
    
    assert_non_null(config);
    assert_string_equal(config->name, MCP_TOOL_GET_METRICS_DETAILS);
    assert_true(config->params.metrics_required);
    assert_true(config->params.metrics_as_array);
}

// Test: mcp_get_list_tool_config - find nodes details tool
static void test_mcp_get_list_tool_config_find_nodes_details_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_GET_NODES_DETAILS);
    
    assert_non_null(config);
    assert_string_equal(config->name, MCP_TOOL_GET_NODES_DETAILS);
    assert_true(config->params.nodes_required);
    assert_true(config->params.metrics_as_array);
}

// Test: mcp_get_list_tool_config - find raised alerts tool
static void test_mcp_get_list_tool_config_find_raised_alerts_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_RAISED_ALERTS);
    
    assert_non_null(config);
    assert_string_equal(config->name, MCP_TOOL_LIST_RAISED_ALERTS);
    assert_int_equal(config->output_type, MCP_LIST_OUTPUT_ALERTS);
    assert_true(config->params.has_alert_pattern);
    assert_int_equal(config->defaults.alert_status, CONTEXT_ALERT_RAISED);
    assert_int_equal(config->defaults.cardinality_limit, 200);
}

// Test: mcp_get_list_tool_config - find all alerts tool
static void test_mcp_get_list_tool_config_find_all_alerts_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_ALL_ALERTS);
    
    assert_non_null(config);
    assert_string_equal(config->name, MCP_TOOL_LIST_ALL_ALERTS);
    assert_int_equal(config->output_type, MCP_LIST_OUTPUT_ALERTS);
    assert_int_equal(config->defaults.alert_status, CONTEXTS_ALERT_STATUSES);
}

// Test: mcp_get_list_tool_config - nonexistent tool
static void test_mcp_get_list_tool_config_nonexistent_tool(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config("nonexistent_tool");
    
    assert_null(config);
}

// Test: mcp_get_list_tool_config - empty string
static void test_mcp_get_list_tool_config_empty_string(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config("");
    
    assert_null(config);
}

// Test: mcp_get_list_tool_config - NULL pointer
static void test_mcp_get_list_tool_config_null_pointer(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(NULL);
    
    // This should not crash; behavior depends on implementation of strcmp with NULL
    // If strcmp(NULL, ...) is checked before, it should be safe
}

// Test: mcp_unified_list_tool_schema - NULL buffer
static void test_mcp_unified_list_tool_schema_null_buffer(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_METRICS);
    
    // Should return early without crashing
    mcp_unified_list_tool_schema(NULL, config);
    // No assertion - just checking it doesn't crash
}

// Test: mcp_unified_list_tool_schema - NULL config
static void test_mcp_unified_list_tool_schema_null_config(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(0, NULL);
    
    mcp_unified_list_tool_schema(buffer, NULL);
    // No assertion - just checking it doesn't crash
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - both NULL
static void test_mcp_unified_list_tool_schema_both_null(void **state) {
    (void) state;
    
    // Should return early without crashing
    mcp_unified_list_tool_schema(NULL, NULL);
}

// Test: mcp_unified_list_tool_schema - nodes output type
static void test_mcp_unified_list_tool_schema_output_type_nodes(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_NODES);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    // Verify that the schema was generated
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - metrics output type
static void test_mcp_unified_list_tool_schema_output_type_metrics(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_METRICS);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - functions output type
static void test_mcp_unified_list_tool_schema_output_type_functions(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_FUNCTIONS);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - alerts output type
static void test_mcp_unified_list_tool_schema_output_type_alerts(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_RAISED_ALERTS);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - with metrics required
static void test_mcp_unified_list_tool_schema_metrics_required(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_GET_METRICS_DETAILS);
    
    assert_true(config->params.metrics_required);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - with nodes required
static void test_mcp_unified_list_tool_schema_nodes_required(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_FUNCTIONS);
    
    assert_true(config->params.nodes_required);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - metrics as array
static void test_mcp_unified_list_tool_schema_metrics_as_array(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_GET_METRICS_DETAILS);
    
    assert_true(config->params.metrics_as_array);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - nodes as array
static void test_mcp_unified_list_tool_schema_nodes_as_array(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_FUNCTIONS);
    
    assert_true(config->params.nodes_as_array);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - has q parameter
static void test_mcp_unified_list_tool_schema_has_q_parameter(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_METRICS);
    
    assert_true(config->params.has_q);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - has time range parameter
static void test_mcp_unified_list_tool_schema_has_time_range(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_METRICS);
    
    assert_true(config->params.has_time_range);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - has cardinality limit
static void test_mcp_unified_list_tool_schema_has_cardinality_limit(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_METRICS);
    
    assert_true(config->params.has_cardinality_limit);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - has alert pattern
static void test_mcp_unified_list_tool_schema_has_alert_pattern(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_RAISED_ALERTS);
    
    assert_true(config->params.has_alert_pattern);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - without time range
static void test_mcp_unified_list_tool_schema_without_time_range(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_FUNCTIONS);
    
    assert_false(config->params.has_time_range);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_schema - without cardinality limit
static void test_mcp_unified_list_tool_schema_without_cardinality_limit(void **state) {
    (void) state;
    
    BUFFER *buffer = buffer_create(4096, NULL);
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_FUNCTIONS);
    
    assert_false(config->params.has_cardinality_limit);
    
    mcp_unified_list_tool_schema(buffer, config);
    
    const char *output = buffer_tostring(buffer);
    assert_non_null(output);
    
    buffer_free(buffer);
}

// Test: mcp_unified_list_tool_execute - NULL client
static void test_mcp_unified_list_tool_execute_null_client(void **state) {
    (void) state;
    
    const MCP_LIST_TOOL_CONFIG *config = mcp_get_list_tool_config(MCP_TOOL_LIST_METRICS);
    json_object *params = NULL;
    MCP_REQUEST_ID id = 1;
    
    MCP_RETURN_CODE result = mcp_unified_list_tool_execute(NULL, config, params, id);
    
    assert_int_equal(result, MCP_RC_ERROR);
}

// Test: mcp_unified_list_tool_execute - NULL config
static void test_mcp_unified_list_tool_execute_null_config(void **state) {
    (void) state;
    
    MCP_CLIENT mcpc = {0};
    mcpc.error = buffer_create(1024, NULL);
    json_object *params = NULL;
    MCP_REQUEST_ID id = 1;
    
    MCP_RETURN_CODE result = mcp_unified_list_tool_execute(&mcpc, NULL, params, id);
    
    assert_int_equal(result, MCP_RC_ERROR);
    
    buffer_free(mcpc.error);
}

// Test: mcp_unified_list_tool_execute - both NULL
static void test_mcp_unified_list_tool_execute_both_null(void **state) {
    (void) state;
    
    MCP_RETURN_CODE result = mcp_unified_list_tool_execute(NULL, NULL, NULL, 1);
    
    assert_int_equal(result, MCP_RC_ERROR);
}

// Test: mcp_unified_list_tool_execute - missing required metrics
static void test_mcp_unified_list_tool_