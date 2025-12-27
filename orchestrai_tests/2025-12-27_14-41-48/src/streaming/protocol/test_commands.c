#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <time.h>

/* Mock types and structures */
typedef struct {
    uint64_t flags;
} STREAM_CAPABILITIES;

typedef struct {
    uint64_t flags;
} RRDSET_FLAGS;

typedef struct {
    char *name;
    int id;
    RRDSET_FLAGS flags;
} RRDSET;

typedef struct {
    char *name;
    int id;
    void *rrdset;
} RRDDIM;

typedef struct {
    void *data;
    size_t size;
} BUFFER;

typedef int RRDVAR_ACQUIRED;
typedef int RRDHOST;
typedef int SN_FLAGS;
typedef unsigned long long usec_t;
typedef double NETDATA_DOUBLE;

typedef struct sender_state {
    char *hostname;
    int fd;
} sender_state;

typedef struct rrdset_stream_buffer {
    STREAM_CAPABILITIES capabilities;
    bool v2;
    bool begin_v2_added;
    time_t wall_clock_time;
    RRDSET_FLAGS rrdset_flags;
    time_t last_point_end_time_s;
    BUFFER *wb;
} RRDSET_STREAM_BUFFER;

/* Function declarations from commands.h */
RRDSET_STREAM_BUFFER stream_send_metrics_init(RRDSET *st, time_t wall_clock_time);
void stream_sender_get_node_and_claim_id_from_parent(struct sender_state *s, const char *claim_id_str, const char *node_id_str, const char *url);
void stream_receiver_send_node_and_claim_id_to_child(RRDHOST *host);
void stream_sender_clear_parent_claim_id(RRDHOST *host);
void stream_sender_send_claimed_id(RRDHOST *host);
void stream_send_global_functions(RRDHOST *host);
void stream_send_host_labels(RRDHOST *host);
void stream_sender_send_custom_host_variables(RRDHOST *host);
void stream_sender_send_this_host_variable_now(RRDHOST *host, const RRDVAR_ACQUIRED *rva);
bool stream_sender_send_rrdset_definition(BUFFER *wb, RRDSET *st);
bool stream_sender_send_rrdset_definition_now(RRDSET *st);
bool should_send_rrdset_matching(RRDSET *st, RRDSET_FLAGS flags);
void stream_send_rrdset_metrics_v1(RRDSET_STREAM_BUFFER *rsb, RRDSET *st);
void stream_send_rrddim_metrics_v2(RRDSET_STREAM_BUFFER *rsb, RRDDIM *rd, usec_t point_end_time_ut, NETDATA_DOUBLE n, SN_FLAGS flags);
void stream_send_rrdset_metrics_finished(RRDSET_STREAM_BUFFER *rsb, RRDSET *st);

/* Mock implementations */
RRDSET_STREAM_BUFFER stream_send_metrics_init(RRDSET *st, time_t wall_clock_time) {
    RRDSET_STREAM_BUFFER rsb;
    memset(&rsb, 0, sizeof(rsb));
    rsb.wall_clock_time = wall_clock_time;
    if (st) {
        rsb.rrdset_flags = st->flags;
    }
    rsb.v2 = false;
    rsb.begin_v2_added = false;
    rsb.last_point_end_time_s = 0;
    rsb.wb = NULL;
    memset(&rsb.capabilities, 0, sizeof(rsb.capabilities));
    return rsb;
}

void stream_sender_get_node_and_claim_id_from_parent(struct sender_state *s, const char *claim_id_str, const char *node_id_str, const char *url) {
    check_expected(s);
    check_expected(claim_id_str);
    check_expected(node_id_str);
    check_expected(url);
}

void stream_receiver_send_node_and_claim_id_to_child(RRDHOST *host) {
    check_expected(host);
}

void stream_sender_clear_parent_claim_id(RRDHOST *host) {
    check_expected(host);
}

void stream_sender_send_claimed_id(RRDHOST *host) {
    check_expected(host);
}

void stream_send_global_functions(RRDHOST *host) {
    check_expected(host);
}

void stream_send_host_labels(RRDHOST *host) {
    check_expected(host);
}

void stream_sender_send_custom_host_variables(RRDHOST *host) {
    check_expected(host);
}

void stream_sender_send_this_host_variable_now(RRDHOST *host, const RRDVAR_ACQUIRED *rva) {
    check_expected(host);
    check_expected(rva);
}

bool stream_sender_send_rrdset_definition(BUFFER *wb, RRDSET *st) {
    check_expected(wb);
    check_expected(st);
    return (bool)mock();
}

bool stream_sender_send_rrdset_definition_now(RRDSET *st) {
    check_expected(st);
    return (bool)mock();
}

bool should_send_rrdset_matching(RRDSET *st, RRDSET_FLAGS flags) {
    check_expected(st);
    check_expected_ptr(&flags);
    return (bool)mock();
}

void stream_send_rrdset_metrics_v1(RRDSET_STREAM_BUFFER *rsb, RRDSET *st) {
    check_expected(rsb);
    check_expected(st);
}

void stream_send_rrddim_metrics_v2(RRDSET_STREAM_BUFFER *rsb, RRDDIM *rd, usec_t point_end_time_ut, NETDATA_DOUBLE n, SN_FLAGS flags) {
    check_expected(rsb);
    check_expected(rd);
    check_expected(point_end_time_ut);
    check_expected(n);
    check_expected(flags);
}

void stream_send_rrdset_metrics_finished(RRDSET_STREAM_BUFFER *rsb, RRDSET *st) {
    check_expected(rsb);
    check_expected(st);
}

/* Test Setup and Teardown */
static int test_setup(void **state) {
    return 0;
}

static int test_teardown(void **state) {
    return 0;
}

/* Tests for stream_send_metrics_init */
static void test_stream_send_metrics_init_with_valid_rrdset(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    st.flags.flags = 0x1234;
    
    time_t wall_time = time(NULL);
    
    RRDSET_STREAM_BUFFER result = stream_send_metrics_init(&st, wall_time);
    
    assert_int_equal(result.wall_clock_time, wall_time);
    assert_int_equal(result.rrdset_flags.flags, 0x1234);
    assert_false(result.v2);
    assert_false(result.begin_v2_added);
    assert_int_equal(result.last_point_end_time_s, 0);
    assert_null(result.wb);
}

static void test_stream_send_metrics_init_with_null_rrdset(void **state) {
    time_t wall_time = time(NULL);
    
    RRDSET_STREAM_BUFFER result = stream_send_metrics_init(NULL, wall_time);
    
    assert_int_equal(result.wall_clock_time, wall_time);
    assert_int_equal(result.rrdset_flags.flags, 0);
    assert_false(result.v2);
    assert_false(result.begin_v2_added);
    assert_int_equal(result.last_point_end_time_s, 0);
    assert_null(result.wb);
}

static void test_stream_send_metrics_init_with_zero_wall_clock_time(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    
    RRDSET_STREAM_BUFFER result = stream_send_metrics_init(&st, 0);
    
    assert_int_equal(result.wall_clock_time, 0);
    assert_false(result.v2);
    assert_false(result.begin_v2_added);
}

static void test_stream_send_metrics_init_with_large_wall_clock_time(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    
    time_t large_time = 2147483647; /* Max 32-bit signed int */
    
    RRDSET_STREAM_BUFFER result = stream_send_metrics_init(&st, large_time);
    
    assert_int_equal(result.wall_clock_time, large_time);
    assert_false(result.v2);
}

static void test_stream_send_metrics_init_with_max_rrdset_flags(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    st.flags.flags = 0xFFFFFFFFFFFFFFFFULL;
    
    time_t wall_time = time(NULL);
    
    RRDSET_STREAM_BUFFER result = stream_send_metrics_init(&st, wall_time);
    
    assert_int_equal(result.rrdset_flags.flags, 0xFFFFFFFFFFFFFFFFULL);
}

/* Tests for stream_sender_get_node_and_claim_id_from_parent */
static void test_stream_sender_get_node_and_claim_id_from_parent_with_valid_params(void **state) {
    struct sender_state s = {.hostname = "test_host", .fd = 42};
    const char *claim_id = "claim123";
    const char *node_id = "node456";
    const char *url = "http://test.example.com";
    
    expect_value(stream_sender_get_node_and_claim_id_from_parent, s, &s);
    expect_string(stream_sender_get_node_and_claim_id_from_parent, claim_id_str, claim_id);
    expect_string(stream_sender_get_node_and_claim_id_from_parent, node_id_str, node_id);
    expect_string(stream_sender_get_node_and_claim_id_from_parent, url, url);
    
    stream_sender_get_node_and_claim_id_from_parent(&s, claim_id, node_id, url);
}

static void test_stream_sender_get_node_and_claim_id_from_parent_with_null_strings(void **state) {
    struct sender_state s = {.hostname = "test_host", .fd = 42};
    
    expect_value(stream_sender_get_node_and_claim_id_from_parent, s, &s);
    expect_null(stream_sender_get_node_and_claim_id_from_parent, claim_id_str);
    expect_null(stream_sender_get_node_and_claim_id_from_parent, node_id_str);
    expect_null(stream_sender_get_node_and_claim_id_from_parent, url);
    
    stream_sender_get_node_and_claim_id_from_parent(&s, NULL, NULL, NULL);
}

static void test_stream_sender_get_node_and_claim_id_from_parent_with_empty_strings(void **state) {
    struct sender_state s = {.hostname = "", .fd = -1};
    
    expect_value(stream_sender_get_node_and_claim_id_from_parent, s, &s);
    expect_string(stream_sender_get_node_and_claim_id_from_parent, claim_id_str, "");
    expect_string(stream_sender_get_node_and_claim_id_from_parent, node_id_str, "");
    expect_string(stream_sender_get_node_and_claim_id_from_parent, url, "");
    
    stream_sender_get_node_and_claim_id_from_parent(&s, "", "", "");
}

static void test_stream_sender_get_node_and_claim_id_from_parent_with_long_strings(void **state) {
    struct sender_state s = {.hostname = "test", .fd = 1};
    char long_claim_id[1024];
    char long_node_id[1024];
    char long_url[2048];
    
    memset(long_claim_id, 'a', sizeof(long_claim_id) - 1);
    long_claim_id[sizeof(long_claim_id) - 1] = '\0';
    
    memset(long_node_id, 'b', sizeof(long_node_id) - 1);
    long_node_id[sizeof(long_node_id) - 1] = '\0';
    
    memset(long_url, 'c', sizeof(long_url) - 1);
    long_url[sizeof(long_url) - 1] = '\0';
    
    expect_value(stream_sender_get_node_and_claim_id_from_parent, s, &s);
    expect_string(stream_sender_get_node_and_claim_id_from_parent, claim_id_str, long_claim_id);
    expect_string(stream_sender_get_node_and_claim_id_from_parent, node_id_str, long_node_id);
    expect_string(stream_sender_get_node_and_claim_id_from_parent, url, long_url);
    
    stream_sender_get_node_and_claim_id_from_parent(&s, long_claim_id, long_node_id, long_url);
}

/* Tests for stream_receiver_send_node_and_claim_id_to_child */
static void test_stream_receiver_send_node_and_claim_id_to_child_with_valid_host(void **state) {
    RRDHOST host = 1;
    
    expect_value(stream_receiver_send_node_and_claim_id_to_child, host, host);
    
    stream_receiver_send_node_and_claim_id_to_child(&host);
}

static void test_stream_receiver_send_node_and_claim_id_to_child_with_null_host(void **state) {
    expect_null(stream_receiver_send_node_and_claim_id_to_child, host);
    
    stream_receiver_send_node_and_claim_id_to_child(NULL);
}

static void test_stream_receiver_send_node_and_claim_id_to_child_multiple_calls(void **state) {
    RRDHOST host1 = 1;
    RRDHOST host2 = 2;
    
    expect_value(stream_receiver_send_node_and_claim_id_to_child, host, host1);
    expect_value(stream_receiver_send_node_and_claim_id_to_child, host, host2);
    
    stream_receiver_send_node_and_claim_id_to_child(&host1);
    stream_receiver_send_node_and_claim_id_to_child(&host2);
}

/* Tests for stream_sender_clear_parent_claim_id */
static void test_stream_sender_clear_parent_claim_id_with_valid_host(void **state) {
    RRDHOST host = 99;
    
    expect_value(stream_sender_clear_parent_claim_id, host, host);
    
    stream_sender_clear_parent_claim_id(&host);
}

static void test_stream_sender_clear_parent_claim_id_with_null_host(void **state) {
    expect_null(stream_sender_clear_parent_claim_id, host);
    
    stream_sender_clear_parent_claim_id(NULL);
}

/* Tests for stream_sender_send_claimed_id */
static void test_stream_sender_send_claimed_id_with_valid_host(void **state) {
    RRDHOST host = 42;
    
    expect_value(stream_sender_send_claimed_id, host, host);
    
    stream_sender_send_claimed_id(&host);
}

static void test_stream_sender_send_claimed_id_with_null_host(void **state) {
    expect_null(stream_sender_send_claimed_id, host);
    
    stream_sender_send_claimed_id(NULL);
}

static void test_stream_sender_send_claimed_id_multiple_hosts(void **state) {
    RRDHOST host1 = 1;
    RRDHOST host2 = 2;
    RRDHOST host3 = 3;
    
    expect_value(stream_sender_send_claimed_id, host, host1);
    expect_value(stream_sender_send_claimed_id, host, host2);
    expect_value(stream_sender_send_claimed_id, host, host3);
    
    stream_sender_send_claimed_id(&host1);
    stream_sender_send_claimed_id(&host2);
    stream_sender_send_claimed_id(&host3);
}

/* Tests for stream_send_global_functions */
static void test_stream_send_global_functions_with_valid_host(void **state) {
    RRDHOST host = 55;
    
    expect_value(stream_send_global_functions, host, host);
    
    stream_send_global_functions(&host);
}

static void test_stream_send_global_functions_with_null_host(void **state) {
    expect_null(stream_send_global_functions, host);
    
    stream_send_global_functions(NULL);
}

/* Tests for stream_send_host_labels */
static void test_stream_send_host_labels_with_valid_host(void **state) {
    RRDHOST host = 77;
    
    expect_value(stream_send_host_labels, host, host);
    
    stream_send_host_labels(&host);
}

static void test_stream_send_host_labels_with_null_host(void **state) {
    expect_null(stream_send_host_labels, host);
    
    stream_send_host_labels(NULL);
}

/* Tests for stream_sender_send_custom_host_variables */
static void test_stream_sender_send_custom_host_variables_with_valid_host(void **state) {
    RRDHOST host = 88;
    
    expect_value(stream_sender_send_custom_host_variables, host, host);
    
    stream_sender_send_custom_host_variables(&host);
}

static void test_stream_sender_send_custom_host_variables_with_null_host(void **state) {
    expect_null(stream_sender_send_custom_host_variables, host);
    
    stream_sender_send_custom_host_variables(NULL);
}

/* Tests for stream_sender_send_this_host_variable_now */
static void test_stream_sender_send_this_host_variable_now_with_valid_params(void **state) {
    RRDHOST host = 33;
    RRDVAR_ACQUIRED rva = 44;
    
    expect_value(stream_sender_send_this_host_variable_now, host, host);
    expect_value(stream_sender_send_this_host_variable_now, rva, &rva);
    
    stream_sender_send_this_host_variable_now(&host, &rva);
}

static void test_stream_sender_send_this_host_variable_now_with_null_host(void **state) {
    RRDVAR_ACQUIRED rva = 44;
    
    expect_null(stream_sender_send_this_host_variable_now, host);
    expect_value(stream_sender_send_this_host_variable_now, rva, &rva);
    
    stream_sender_send_this_host_variable_now(NULL, &rva);
}

static void test_stream_sender_send_this_host_variable_now_with_null_rva(void **state) {
    RRDHOST host = 33;
    
    expect_value(stream_sender_send_this_host_variable_now, host, host);
    expect_null(stream_sender_send_this_host_variable_now, rva);
    
    stream_sender_send_this_host_variable_now(&host, NULL);
}

/* Tests for stream_sender_send_rrdset_definition */
static void test_stream_sender_send_rrdset_definition_with_valid_params_true(void **state) {
    BUFFER wb;
    RRDSET st;
    memset(&wb, 0, sizeof(wb));
    memset(&st, 0, sizeof(st));
    
    expect_value(stream_sender_send_rrdset_definition, wb, &wb);
    expect_value(stream_sender_send_rrdset_definition, st, &st);
    will_return(stream_sender_send_rrdset_definition, 1);
    
    bool result = stream_sender_send_rrdset_definition(&wb, &st);
    
    assert_true(result);
}

static void test_stream_sender_send_rrdset_definition_with_valid_params_false(void **state) {
    BUFFER wb;
    RRDSET st;
    memset(&wb, 0, sizeof(wb));
    memset(&st, 0, sizeof(st));
    
    expect_value(stream_sender_send_rrdset_definition, wb, &wb);
    expect_value(stream_sender_send_rrdset_definition, st, &st);
    will_return(stream_sender_send_rrdset_definition, 0);
    
    bool result = stream_sender_send_rrdset_definition(&wb, &st);
    
    assert_false(result);
}

static void test_stream_sender_send_rrdset_definition_with_null_buffer(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    
    expect_null(stream_sender_send_rrdset_definition, wb);
    expect_value(stream_sender_send_rrdset_definition, st, &st);
    will_return(stream_sender_send_rrdset_definition, 0);
    
    bool result = stream_sender_send_rrdset_definition(NULL, &st);
    
    assert_false(result);
}

static void test_stream_sender_send_rrdset_definition_with_null_rrdset(void **state) {
    BUFFER wb;
    memset(&wb, 0, sizeof(wb));
    
    expect_value(stream_sender_send_rrdset_definition, wb, &wb);
    expect_null(stream_sender_send_rrdset_definition, st);
    will_return(stream_sender_send_rrdset_definition, 0);
    
    bool result = stream_sender_send_rrdset_definition(&wb, NULL);
    
    assert_false(result);
}

static void test_stream_sender_send_rrdset_definition_with_both_null(void **state) {
    expect_null(stream_sender_send_rrdset_definition, wb);
    expect_null(stream_sender_send_rrdset_definition, st);
    will_return(stream_sender_send_rrdset_definition, 0);
    
    bool result = stream_sender_send_rrdset_definition(NULL, NULL);
    
    assert_false(result);
}

/* Tests for stream_sender_send_rrdset_definition_now */
static void test_stream_sender_send_rrdset_definition_now_true(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    
    expect_value(stream_sender_send_rrdset_definition_now, st, &st);
    will_return(stream_sender_send_rrdset_definition_now, 1);
    
    bool result = stream_sender_send_rrdset_definition_now(&st);
    
    assert_true(result);
}

static void test_stream_sender_send_rrdset_definition_now_false(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    
    expect_value(stream_sender_send_rrdset_definition_now, st, &st);
    will_return(stream_sender_send_rrdset_definition_now, 0);
    
    bool result = stream_sender_send_rrdset_definition_now(&st);
    
    assert_false(result);
}

static void test_stream_sender_send_rrdset_definition_now_null(void **state) {
    expect_null(stream_sender_send_rrdset_definition_now, st);
    will_return(stream_sender_send_rrdset_definition_now, 0);
    
    bool result = stream_sender_send_rrdset_definition_now(NULL);
    
    assert_false(result);
}

/* Tests for should_send_rrdset_matching */
static void test_should_send_rrdset_matching_true(void **state) {
    RRDSET st;
    RRDSET_FLAGS flags;
    memset(&st, 0, sizeof(st));
    memset(&flags, 0, sizeof(flags));
    flags.flags = 0x0001;
    
    expect_value(should_send_rrdset_matching, st, &st);
    expect_ptr(should_send_rrdset_matching, flags);
    will_return(should_send_rrdset_matching, 1);
    
    bool result = should_send_rrdset_matching(&st, flags);
    
    assert_true(result);
}

static void test_should_send_rrdset_matching_false(void **state) {
    RRDSET st;
    RRDSET_FLAGS flags;
    memset(&st, 0, sizeof(st));
    memset(&flags, 0, sizeof(flags));
    
    expect_value(should_send_rrdset_matching, st, &st);
    expect_ptr(should_send_rrdset_matching, flags);
    will_return(should_send_rrdset_matching, 0);
    
    bool result = should_send_rrdset_matching(&st, flags);
    
    assert_false(result);
}

static void test_should_send_rrdset_matching_null_rrdset(void **state) {
    RRDSET_FLAGS flags;
    memset(&flags, 0, sizeof(flags));
    
    expect_null(should_send_rrdset_matching, st);
    expect_ptr(should_send_rrdset_matching, flags);
    will_return(should_send_rrdset_matching, 0);
    
    bool result = should_send_rrdset_matching(NULL, flags);
    
    assert_false(result);
}

static void test_should_send_rrdset_matching_with_max_flags(void **state) {
    RRDSET st;
    RRDSET_FLAGS flags;
    memset(&st, 0, sizeof(st));
    memset(&flags, 0, sizeof(flags));
    flags.flags = 0xFFFFFFFFFFFFFFFFULL;
    
    expect_value(should_send_rrdset_matching, st, &st);
    expect_ptr(should_send_rrdset_matching, flags);
    will_return(should_send_rrdset_matching, 1);
    
    bool result = should_send_rrdset_matching(&st, flags);
    
    assert_true(result);
}

/* Tests for stream_send_rrdset_metrics_v1 */
static void test_stream_send_rrdset_metrics_v1_with_valid_params(void **state) {
    RRDSET_STREAM_BUFFER rsb;
    RRDSET st;
    memset(&rsb, 0, sizeof(rsb));
    memset(&st, 0, sizeof(st));
    
    expect_value(stream_send_rrdset_metrics_v1, rsb, &rsb);
    expect_value(stream_send_rrdset_metrics_v1, st, &st);
    
    stream_send_rrdset_metrics_v1(&rsb, &st);
}

static void test_stream_send_rrdset_metrics_v1_with_null_rsb(void **state) {
    RRDSET st;
    memset(&st, 0, sizeof(st));
    
    expect_null(stream_send_rrdset_metrics_v1, rsb);
    expect_value(stream_send_rrdset_metrics_v1, st, &st);
    
    stream_send_rrdset_metrics_v1(NULL, &st);
}

static void test_stream_send_rrdset_metrics_v1_with_null_st(void **state) {
    RRDSET_STREAM_BUFFER rsb;
    memset(&rsb, 0, sizeof(rsb));
    
    expect_value(stream_send_rrdset_metrics_v1, rsb, &rsb);
    expect_null(stream_send_rrdset_metrics_v1, st);
    
    stream_send_rrdset_metrics_v1(&rsb, NULL);
}

static void test_stream_send_rrdset_metrics_v1_with_both_null(void **state) {
    expect_null(stream_send_rrdset_metrics_v1, rsb);
    expect_null(stream_send_rrdset_metrics_v1, st);
    
    stream_send_rrdset_metrics_v1(NULL, NULL);
}

static void test_stream_send_rrdset_metrics_v1_with_initialized_rsb(void **state) {
    RRDSET_STREAM_BUFFER rsb;
    RRDSET st;
    memset(&rsb, 0, sizeof(rsb));
    memset(&st, 0, sizeof(st));
    
    rsb.v2 = true;
    rsb.wall_clock_time = 1234567890;
    rsb.begin_v2_added = true;
    
    expect_value(stream_send_rrdset_metrics_v1, rsb, &rsb);
    expect_value(stream_send_rrdset_metrics_v1, st, &st);
    
    stream_send_rrdset_metrics_v1(&rsb, &st);
}

/* Tests for stream_send_rrddim_metrics_v2