#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <math.h>
#include <time.h>

/* Mock types and structures */
typedef struct {
    char *data;
    size_t size;
    size_t pos;
} BUFFER;

typedef struct {
    int update_every;
    unsigned int options;
    char *name;
} instance_config;

typedef struct {
    void *seb;
    void *smh;
} tier;

typedef struct {
    int update_every;
    char rrd_update_every;
    char *hostname;
} RRDHOST;

typedef struct {
    RRDHOST *rrdhost;
    char *id;
    int update_every;
} RRDSET;

typedef struct {
    RRDSET *rrdset;
    char *id;
    tier tiers[1];
} RRDDIM;

typedef struct {
    int sum;
    int count;
    int flags;
} STORAGE_POINT;

typedef struct {
    void *state;
} storage_engine_query_handle;

typedef struct {
    unsigned int buffered_metrics;
    unsigned int buffered_bytes;
    unsigned int data_lost_events;
    unsigned int lost_metrics;
    unsigned int lost_bytes;
} stats;

typedef struct {
    BUFFER *buffer;
    BUFFER *header;
    int buffered_metrics;
    int buffered_bytes;
    int used;
} simple_connector_buffer;

typedef struct {
    simple_connector_buffer *last_buffer;
    simple_connector_buffer *first_buffer;
    int total_buffered_metrics;
} simple_connector_data;

typedef struct {
    instance_config config;
    stats stats;
    void *connector_specific_data;
    void *buffer;
    int scheduled;
    int disabled;
    int skip_host;
    int skip_chart;
    time_t after;
    time_t before;
    int data_is_ready;
    int (*start_batch_formatting)(void *);
    int (*end_batch_formatting)(void *);
    int (*start_host_formatting)(void *, RRDHOST *);
    int (*end_host_formatting)(void *, RRDHOST *);
    int (*start_chart_formatting)(void *, RRDSET *);
    int (*end_chart_formatting)(void *, RRDSET *);
    int (*metric_formatting)(void *, RRDDIM *);
    int (*variables_formatting)(void *, RRDHOST *);
    void (*prepare_header)(void *);
    void *cond_var;
    void *mutex;
} instance;

typedef struct {
    time_t now;
    instance *instance_root;
} engine;

typedef double NETDATA_DOUBLE;

/* Global variables */
RRDHOST *localhost;

/* Function declarations (prototypes from process_data.c) */
size_t exporting_name_copy(char *dst, const char *src, size_t max_len);
int mark_scheduled_instances(struct engine *engine);
NETDATA_DOUBLE exporting_calculate_value_from_stored_data(
    struct instance *instance,
    RRDDIM *rd,
    time_t *last_timestamp);
void start_batch_formatting(struct engine *engine);
void start_host_formatting(struct engine *engine, RRDHOST *host);
void start_chart_formatting(struct engine *engine, RRDSET *st);
void metric_formatting(struct engine *engine, RRDDIM *rd);
void end_chart_formatting(struct engine *engine, RRDSET *st);
void variables_formatting(struct engine *engine, RRDHOST *host);
void end_host_formatting(struct engine *engine, RRDHOST *host);
void end_batch_formatting(struct engine *engine);
void prepare_buffers(struct engine *engine);
int flush_host_labels(struct instance *instance, RRDHOST *host);
int simple_connector_end_batch(struct instance *instance);

/* Mock helper functions */
static int mock_isalnum(unsigned char c) {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9');
}

static void mock_netdata_log_error(const char *fmt, ...) {
    (void)fmt;
}

static void mock_netdata_log_debug(int level, const char *fmt, ...) {
    (void)level;
    (void)fmt;
}

static void mock_netdata_log_info(const char *fmt, ...) {
    (void)fmt;
}

/* ==================== Tests for exporting_name_copy ==================== */

void test_exporting_name_copy_alphanumeric(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "test123", 256);
    
    assert_int_equal(result, 7);
    assert_string_equal(dst, "test123");
}

void test_exporting_name_copy_with_special_chars(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "test@#$123", 256);
    
    assert_int_equal(result, 10);
    assert_string_equal(dst, "test_____123");
}

void test_exporting_name_copy_with_dots(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "test.chart.metric", 256);
    
    assert_int_equal(result, 17);
    assert_string_equal(dst, "test.chart.metric");
}

void test_exporting_name_copy_with_spaces(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "test name", 256);
    
    assert_int_equal(result, 9);
    assert_string_equal(dst, "test_name");
}

void test_exporting_name_copy_with_hyphens(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "test-chart", 256);
    
    assert_int_equal(result, 10);
    assert_string_equal(dst, "test_chart");
}

void test_exporting_name_copy_empty_string(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "", 256);
    
    assert_int_equal(result, 0);
    assert_string_equal(dst, "");
}

void test_exporting_name_copy_max_len_exceeded(void **state) {
    (void)state;
    char dst[10];
    size_t result;
    
    result = exporting_name_copy(dst, "verylongname", 5);
    
    assert_int_equal(result, 5);
    assert_true(dst[5] == '\0');
}

void test_exporting_name_copy_single_char(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "a", 256);
    
    assert_int_equal(result, 1);
    assert_string_equal(dst, "a");
}

void test_exporting_name_copy_single_special_char(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, "@", 256);
    
    assert_int_equal(result, 1);
    assert_string_equal(dst, "_");
}

void test_exporting_name_copy_single_dot(void **state) {
    (void)state;
    char dst[256];
    size_t result;
    
    result = exporting_name_copy(dst, ".", 256);
    
    assert_int_equal(result, 1);
    assert_string_equal(dst, ".");
}

/* ==================== Tests for mark_scheduled_instances ==================== */

void test_mark_scheduled_instances_no_instances(void **state) {
    (void)state;
    engine test_engine = {0};
    test_engine.instance_root = NULL;
    test_engine.now = 1000;
    localhost = calloc(1, sizeof(RRDHOST));
    localhost->rrd_update_every = 10;
    
    int result = mark_scheduled_instances(&test_engine);
    
    assert_int_equal(result, 0);
    free(localhost);
}

void test_mark_scheduled_instances_single_scheduled(void **state) {
    (void)state;
    localhost = calloc(1, sizeof(RRDHOST));
    localhost->rrd_update_every = 10;
    
    instance inst = {0};
    inst.disabled = 0;
    inst.config.update_every = 100;
    inst.next = NULL;
    
    engine test_engine = {0};
    test_engine.instance_root = &inst;
    test_engine.now = 950; /* 950 % 100 = 50, >= (100 - 10) = 90? No */
    
    int result = mark_scheduled_instances(&test_engine);
    
    free(localhost);
}

void test_mark_scheduled_instances_disabled(void **state) {
    (void)state;
    localhost = calloc(1, sizeof(RRDHOST));
    localhost->rrd_update_every = 10;
    
    instance inst = {0};
    inst.disabled = 1;
    inst.config.update_every = 100;
    inst.next = NULL;
    
    engine test_engine = {0};
    test_engine.instance_root = &inst;
    test_engine.now = 1000;
    
    int result = mark_scheduled_instances(&test_engine);
    
    assert_int_equal(result, 0);
    assert_int_equal(inst.scheduled, 0);
    free(localhost);
}

void test_mark_scheduled_instances_multiple(void **state) {
    (void)state;
    localhost = calloc(1, sizeof(RRDHOST));
    localhost->rrd_update_every = 10;
    
    instance inst1 = {0};
    inst1.disabled = 0;
    inst1.config.update_every = 100;
    
    instance inst2 = {0};
    inst2.disabled = 0;
    inst2.config.update_every = 50;
    inst2.next = NULL;
    
    inst1.next = &inst2;
    
    engine test_engine = {0};
    test_engine.instance_root = &inst1;
    test_engine.now = 1000;
    
    int result = mark_scheduled_instances(&test_engine);
    
    free(localhost);
}

/* ==================== Tests for exporting_calculate_value_from_stored_data ==================== */

void test_exporting_calculate_value_empty_data(void **state) {
    (void)state;
    instance inst = {0};
    inst.after = 100;
    inst.before = 200;
    inst.config.options = 0;
    
    RRDHOST host = {0};
    host.rrd_update_every = 10;
    
    RRDSET st = {0};
    st.rrdhost = &host;
    st.update_every = 10;
    
    RRDDIM rd = {0};
    rd.rrdset = &st;
    
    time_t last_timestamp = 0;
    
    NETDATA_DOUBLE result = exporting_calculate_value_from_stored_data(&inst, &rd, &last_timestamp);
    
    assert_true(isnan(result));
}

void test_exporting_calculate_value_timeframe_normalization(void **state) {
    (void)state;
    instance inst = {0};
    inst.after = 107;
    inst.before = 217;
    inst.config.options = 0;
    
    RRDHOST host = {0};
    host.rrd_update_every = 10;
    
    RRDSET st = {0};
    st.rrdhost = &host;
    st.update_every = 10;
    
    RRDDIM rd = {0};
    rd.rrdset = &st;
    
    time_t last_timestamp = 0;
    
    NETDATA_DOUBLE result = exporting_calculate_value_from_stored_data(&inst, &rd, &last_timestamp);
    
    assert_true(isnan(result));
}

/* ==================== Tests for start_batch_formatting ==================== */

void test_start_batch_formatting_no_instances(void **state) {
    (void)state;
    engine test_engine = {0};
    test_engine.instance_root = NULL;
    
    start_batch_formatting(&test_engine);
}

void test_start_batch_formatting_not_scheduled(void **state) {
    (void)state;
    instance inst = {0};
    inst.scheduled = 0;
    inst.next = NULL;
    
    engine test_engine = {0};
    test_engine.instance_root = &inst;
    
    start_batch_formatting(&test_engine);
    
    assert_int_equal(inst.scheduled, 0);
}

void test_start_batch_formatting_success(void **state) {
    (void)state;
    instance inst = {0};
    inst.scheduled = 1;
    inst.next = NULL;
    inst.config.name = "test";
    inst.start_batch_formatting = NULL;
    
    engine test_engine = {0};
    test_engine.instance_root = &inst;
    
    start_batch_formatting(&test_engine);
}

/* ==================== Tests for flush_host_labels ==================== */

void test_flush_host_labels_no_buffer(void **state) {
    (void)state;
    instance inst = {0};
    inst.labels_buffer = NULL;
    
    RRDHOST host = {0};
    
    int result = flush_host_labels(&inst, &host);
    
    assert_int_equal(result, 0);
}

void test_flush_host_labels_with_buffer(void **state) {
    (void)state;
    instance inst = {0};
    inst.labels_buffer = calloc(1, sizeof(BUFFER));
    
    RRDHOST host = {0};
    
    int result = flush_host_labels(&inst, &host);
    
    assert_int_equal(result, 0);
    free(inst.labels_buffer);
}

/* ==================== Tests for simple_connector_end_batch ==================== */

void test_simple_connector_end_batch_create_buffer(void **state) {
    (void)state;
    simple_connector_buffer last_buf = {0};
    last_buf.buffer = NULL;
    last_buf.header = NULL;
    last_buf.next = &last_buf;
    last_buf.used = 0;
    
    simple_connector_data connector_data = {0};
    connector_data.last_buffer = &last_buf;
    connector_data.first_buffer = &last_buf;
    connector_data.total_buffered_metrics = 0;
    
    instance inst = {0};
    inst.connector_specific_data = &connector_data;
    inst.buffer = calloc(1, sizeof(BUFFER));
    inst.stats.buffered_metrics = 5;
    inst.stats.buffered_bytes = 0;
    inst.prepare_header = NULL;
    
    int result = simple_connector_end_batch(&inst);
    
    assert_int_equal(result, 0);
    assert_non_null(last_buf.buffer);
    
    free(inst.buffer);
}

void test_simple_connector_end_batch_ring_buffer_full(void **state) {
    (void)state;
    simple_connector_buffer buf1 = {0};
    buf1.buffer = calloc(1, sizeof(BUFFER));
    buf1.used = 1;
    
    simple_connector_buffer buf2 = {0};
    buf2.buffer = calloc(1, sizeof(BUFFER));
    buf2.used = 1;
    buf2.next = &buf1;
    buf2.buffered_metrics = 10;
    buf2.buffered_bytes = 100;
    
    buf1.next = &buf2;
    
    simple_connector_data connector_data = {0};
    connector_data.last_buffer = &buf2;
    connector_data.first_buffer = &buf1;
    connector_data.total_buffered_metrics = 0;
    
    instance inst = {0};
    inst.connector_specific_data = &connector_data;
    inst.buffer = calloc(1, sizeof(BUFFER));
    inst.stats.buffered_metrics = 5;
    inst.stats.buffered_bytes = 0;
    inst.stats.data_lost_events = 0;
    inst.stats.lost_metrics = 0;
    inst.stats.lost_bytes = 0;
    inst.prepare_header = NULL;
    
    int result = simple_connector_end_batch(&inst);
    
    assert_int_equal(result, 0);
    assert_int_equal(inst.stats.data_lost_events, 1);
    assert_int_equal(inst.stats.lost_metrics, 10);
    assert_int_equal(inst.stats.lost_bytes, 100);
    
    free(inst.buffer);
    free(buf1.buffer);
    free(buf2.buffer);
}

void test_simple_connector_end_batch_with_prepare_header(void **state) {
    (void)state;
    simple_connector_buffer last_buf = {0};
    last_buf.buffer = calloc(1, sizeof(BUFFER));
    last_buf.header = calloc(1, sizeof(BUFFER));
    last_buf.next = &last_buf;
    last_buf.used = 0;
    
    simple_connector_data connector_data = {0};
    connector_data.last_buffer = &last_buf;
    connector_data.first_buffer = &last_buf;
    connector_data.total_buffered_metrics = 0;
    
    int header_called = 0;
    void prepare_header_mock(void *instance) {
        header_called = 1;
    }
    
    instance inst = {0};
    inst.connector_specific_data = &connector_data;
    inst.buffer = calloc(1, sizeof(BUFFER));
    inst.stats.buffered_metrics = 3;
    inst.stats.buffered_bytes = 0;
    inst.prepare_header = NULL;
    
    int result = simple_connector_end_batch(&inst);
    
    assert_int_equal(result, 0);
    
    free(inst.buffer);
    free(last_buf.buffer);
    free(last_buf.header);
}

int main(void) {
    const struct CMUnitTest tests[] = {
        /* exporting_name_copy tests */
        cmocka_unit_test(test_exporting_name_copy_alphanumeric),
        cmocka_unit_test(test_exporting_name_copy_with_special_chars),
        cmocka_unit_test(test_exporting_name_copy_with_dots),
        cmocka_unit_test(test_exporting_name_copy_with_spaces),
        cmocka_unit_test(test_exporting_name_copy_with_hyphens),
        cmocka_unit_test(test_exporting_name_copy_empty_string),
        cmocka_unit_test(test_exporting_name_copy_max_len_exceeded),
        cmocka_unit_test(test_exporting_name_copy_single_char),
        cmocka_unit_test(test_exporting_name_copy_single_special_char),
        cmocka_unit_test(test_exporting_name_copy_single_dot),
        
        /* mark_scheduled_instances tests */
        cmocka_unit_test(test_mark_scheduled_instances_no_instances),
        cmocka_unit_test(test_mark_scheduled_instances_single_scheduled),
        cmocka_unit_test(test_mark_scheduled_instances_disabled),
        cmocka_unit_test(test_mark_scheduled_instances_multiple),
        
        /* exporting_calculate_value_from_stored_data tests */
        cmocka_unit_test(test_exporting_calculate_value_empty_data),
        cmocka_unit_test(test_exporting_calculate_value_timeframe_normalization),
        
        /* start_batch_formatting tests */
        cmocka_unit_test(test_start_batch_formatting_no_instances),
        cmocka_unit_test(test_start_batch_formatting_not_scheduled),
        cmocka_unit_test(test_start_batch_formatting_success),
        
        /* flush_host_labels tests */
        cmocka_unit_test(test_flush_host_labels_no_buffer),
        cmocka_unit_test(test_flush_host_labels_with_buffer),
        
        /* simple_connector_end_batch tests */
        cmocka_unit_test(test_simple_connector_end_batch_create_buffer),
        cmocka_unit_test(test_simple_connector_end_batch_ring_buffer_full),
        cmocka_unit_test(test_simple_connector_end_batch_with_prepare_header),
    };
    
    return cmocka_run_group_tests(tests, NULL, NULL);
}