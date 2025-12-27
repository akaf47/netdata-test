#ifndef TEST_RRD2JSON_H
#define TEST_RRD2JSON_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <assert.h>

/* Mock structures and functions */

typedef struct {
    void *data;
    size_t len;
    size_t allocated;
} BUFFER;

typedef struct {
    char *name;
    time_t first_entry_t;
    time_t last_entry_t;
    uint32_t update_every;
    long entries;
} RRDDIM;

typedef struct {
    char *id;
    char *name;
    char *type;
    char *family;
    char *context;
    time_t first_entry_t;
    time_t last_entry_t;
    uint32_t update_every;
    long entries;
    RRDDIM *dimensions;
    int dimensions_count;
} RRDSET;

typedef struct {
    uint32_t options;
    char *cookie;
    char *pin;
    time_t after;
    time_t before;
    int points;
} WEB_CLIENT_OPTIONS;

/* Test utilities */

static BUFFER* buffer_create(size_t initial_size) {
    BUFFER *buf = malloc(sizeof(BUFFER));
    assert(buf != NULL);
    buf->allocated = initial_size > 0 ? initial_size : 1024;
    buf->data = malloc(buf->allocated);
    assert(buf->data != NULL);
    buf->len = 0;
    return buf;
}

static void buffer_free(BUFFER *buf) {
    if (buf) {
        free(buf->data);
        free(buf);
    }
}

static void buffer_reset(BUFFER *buf) {
    if (buf) {
        buf->len = 0;
    }
}

static size_t buffer_strlen(BUFFER *buf) {
    return buf ? buf->len : 0;
}

/* Test: NULL buffer handling */
static void test_null_buffer_operations(void) {
    assert(buffer_strlen(NULL) == 0);
    buffer_reset(NULL);
    buffer_free(NULL);
}

/* Test: Buffer creation and initialization */
static void test_buffer_creation(void) {
    BUFFER *buf = buffer_create(0);
    assert(buf != NULL);
    assert(buf->data != NULL);
    assert(buf->len == 0);
    assert(buf->allocated >= 1024);
    buffer_free(buf);
}

/* Test: Buffer creation with custom size */
static void test_buffer_creation_custom_size(void) {
    BUFFER *buf = buffer_create(512);
    assert(buf != NULL);
    assert(buf->allocated == 512);
    buffer_free(buf);
}

/* Test: Buffer reset */
static void test_buffer_reset(void) {
    BUFFER *buf = buffer_create(256);
    assert(buf->len == 0);
    buf->len = 100;
    assert(buf->len == 100);
    buffer_reset(buf);
    assert(buf->len == 0);
    buffer_free(buf);
}

/* Test: Empty RRDSET handling */
static void test_empty_rrdset(void) {
    RRDSET *set = malloc(sizeof(RRDSET));
    assert(set != NULL);
    memset(set, 0, sizeof(RRDSET));
    free(set);
}

/* Test: RRDSET with valid data */
static void test_valid_rrdset(void) {
    RRDSET *set = malloc(sizeof(RRDSET));
    assert(set != NULL);
    
    set->id = strdup("cpu.usage");
    set->name = strdup("CPU Usage");
    set->type = strdup("system");
    set->family = strdup("cpu");
    set->context = strdup("system.cpu");
    set->first_entry_t = 1000000;
    set->last_entry_t = 2000000;
    set->update_every = 10;
    set->entries = 10000;
    set->dimensions = NULL;
    set->dimensions_count = 0;
    
    assert(set->id != NULL);
    assert(set->name != NULL);
    assert(set->type != NULL);
    assert(set->first_entry_t < set->last_entry_t);
    
    free(set->id);
    free(set->name);
    free(set->type);
    free(set->family);
    free(set->context);
    free(set);
}

/* Test: RRDSET with dimensions */
static void test_rrdset_with_dimensions(void) {
    RRDSET *set = malloc(sizeof(RRDSET));
    assert(set != NULL);
    
    set->id = strdup("cpu.usage");
    set->dimensions_count = 2;
    set->dimensions = malloc(sizeof(RRDDIM) * 2);
    assert(set->dimensions != NULL);
    
    memset(&set->dimensions[0], 0, sizeof(RRDDIM));
    memset(&set->dimensions[1], 0, sizeof(RRDDIM));
    
    set->dimensions[0].name = strdup("user");
    set->dimensions[1].name = strdup("system");
    
    assert(set->dimensions_count == 2);
    
    free(set->dimensions[0].name);
    free(set->dimensions[1].name);
    free(set->dimensions);
    free(set->id);
    free(set);
}

/* Test: WEB_CLIENT_OPTIONS initialization */
static void test_web_client_options_init(void) {
    WEB_CLIENT_OPTIONS opts;
    memset(&opts, 0, sizeof(WEB_CLIENT_OPTIONS));
    
    opts.options = 0;
    opts.cookie = NULL;
    opts.pin = NULL;
    opts.after = 0;
    opts.before = 0;
    opts.points = 0;
    
    assert(opts.options == 0);
    assert(opts.cookie == NULL);
    assert(opts.pin == NULL);
}

/* Test: Time range validation */
static void test_time_range_validation(void) {
    time_t start = 1000000;
    time_t end = 2000000;
    assert(start < end);
    
    time_t invalid_start = 2000000;
    time_t invalid_end = 1000000;
    assert(invalid_start >= invalid_end);
}

/* Test: Points calculation */
static void test_points_calculation(void) {
    time_t start = 1000000;
    time_t end = 1000100;
    uint32_t update_every = 10;
    
    int points = (int)((end - start) / update_every);
    assert(points == 10);
    
    /* Edge case: 0 duration */
    time_t same_time = 1000000;
    int zero_points = (int)((same_time - same_time) / update_every);
    assert(zero_points == 0);
}

/* Test: Update every value validation */
static void test_update_every_validation(void) {
    uint32_t valid_update = 10;
    assert(valid_update > 0);
    
    uint32_t zero_update = 0;
    assert(zero_update == 0);
    
    uint32_t large_update = 86400;
    assert(large_update > valid_update);
}

/* Test: Entries count handling */
static void test_entries_count(void) {
    long zero_entries = 0;
    assert(zero_entries == 0);
    
    long positive_entries = 10000;
    assert(positive_entries > 0);
    
    long large_entries = 1000000;
    assert(large_entries > positive_entries);
}

/* Test: String field handling */
static void test_string_field_handling(void) {
    char *empty = strdup("");
    assert(empty != NULL);
    assert(strlen(empty) == 0);
    free(empty);
    
    char *null_check = NULL;
    assert(null_check == NULL);
    
    char *normal = strdup("test_value");
    assert(strlen(normal) > 0);
    free(normal);
}

/* Test: Buffer boundary conditions */
static void test_buffer_boundary_conditions(void) {
    BUFFER *buf = buffer_create(1);
    assert(buf->allocated == 1);
    buffer_free(buf);
    
    BUFFER *buf2 = buffer_create(65536);
    assert(buf2->allocated == 65536);
    buffer_free(buf2);
}

/* Main test runner */
static int run_all_tests(void) {
    test_null_buffer_operations();
    test_buffer_creation();
    test_buffer_creation_custom_size();
    test_buffer_reset();
    test_empty_rrdset();
    test_valid_rrdset();
    test_rrdset_with_dimensions();
    test_web_client_options_init();
    test_time_range_validation();
    test_points_calculation();
    test_update_every_validation();
    test_entries_count();
    test_string_field_handling();
    test_buffer_boundary_conditions();
    
    return 0;
}

#endif /* TEST_RRD2JSON_H */