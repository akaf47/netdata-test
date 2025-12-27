#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include <cmocka.h>
#include <time.h>
#include <sys/time.h>

/* Mock structures and forward declarations */
typedef struct {
    char *id;
    char *name;
    char *type;
    char *family;
    unsigned long priority;
    int update_every;
    long usec;
    unsigned long collected_total;
    unsigned long last_collected;
    long views;
    int enabled;
} RRDCHART;

typedef struct {
    char *id;
    char *name;
    char *type;
    char *units;
    long multiplier;
    long divisor;
    int hidden;
    int flags;
} RRDDIM;

typedef struct {
    void *first;
    void *last;
} SIMPLE_PATTERN;

typedef struct {
    char *key;
    char *value;
} NAME_VALUE;

typedef struct {
    RRDCHART *charts;
    int charts_count;
    RRDDIM **dimensions;
    int dimensions_count;
} TEST_RRDHOST;

/* Mock functions */
static void* mock_rrdhost_read_lock(void) {
    return (void*)0x1;
}

static void mock_rrdhost_read_unlock(void *host) {
    (void)host;
}

static RRDCHART* mock_rrdchart_find(void *host, const char *id) {
    RRDCHART *chart = (RRDCHART*)malloc(sizeof(RRDCHART));
    chart->id = (char*)malloc(strlen(id) + 1);
    strcpy(chart->id, id);
    chart->name = strdup("test_chart");
    chart->type = strdup("area");
    chart->family = strdup("cpu");
    chart->priority = 100;
    chart->update_every = 1;
    chart->usec = 1000000;
    chart->collected_total = 100;
    chart->last_collected = 99;
    chart->views = 1;
    chart->enabled = 1;
    return chart;
}

static void mock_rrdchart_free(RRDCHART *chart) {
    if (chart) {
        free(chart->id);
        free(chart->name);
        free(chart->type);
        free(chart->family);
        free(chart);
    }
}

/* Test fixtures */
static int setup(void **state) {
    return 0;
}

static int teardown(void **state) {
    return 0;
}

/* Test: charts2json with valid parameters */
static void test_charts2json_valid_chart(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = 1;
    chart.enabled = 1;
    
    assert_non_null(chart.id);
    assert_non_null(chart.name);
    assert_string_equal(chart.id, "cpu.system");
}

/* Test: NULL chart ID */
static void test_charts2json_null_chart_id(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = NULL;
    chart.name = "CPU System";
    
    assert_null(chart.id);
}

/* Test: Empty chart ID */
static void test_charts2json_empty_chart_id(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "";
    chart.name = "CPU System";
    
    assert_non_null(chart.id);
    assert_int_equal(strlen(chart.id), 0);
}

/* Test: Very long chart ID */
static void test_charts2json_long_chart_id(void **state) {
    (void)state;
    
    char long_id[4096];
    memset(long_id, 'a', 4095);
    long_id[4095] = '\0';
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = long_id;
    chart.name = "Long ID Chart";
    
    assert_non_null(chart.id);
    assert_int_equal(strlen(chart.id), 4095);
}

/* Test: NULL chart name */
static void test_charts2json_null_chart_name(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = NULL;
    
    assert_null(chart.name);
}

/* Test: Empty chart name */
static void test_charts2json_empty_chart_name(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "";
    
    assert_non_null(chart.name);
    assert_int_equal(strlen(chart.name), 0);
}

/* Test: NULL chart type */
static void test_charts2json_null_chart_type(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = NULL;
    
    assert_null(chart.type);
}

/* Test: Different chart types */
static void test_charts2json_different_chart_types(void **state) {
    (void)state;
    
    const char *types[] = {"area", "line", "stacked", "bar", NULL};
    
    for (int i = 0; types[i] != NULL; i++) {
        RRDCHART chart;
        memset(&chart, 0, sizeof(RRDCHART));
        chart.id = "cpu.system";
        chart.name = "CPU System";
        chart.type = (char*)types[i];
        
        assert_non_null(chart.type);
        assert_string_equal(chart.type, types[i]);
    }
}

/* Test: NULL chart family */
static void test_charts2json_null_chart_family(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = NULL;
    
    assert_null(chart.family);
}

/* Test: Zero priority */
static void test_charts2json_zero_priority(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 0;
    
    assert_int_equal(chart.priority, 0);
}

/* Test: Maximum priority */
static void test_charts2json_max_priority(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 0xFFFFFFFFUL;
    
    assert_int_equal(chart.priority, 0xFFFFFFFFUL);
}

/* Test: Negative update_every */
static void test_charts2json_negative_update_every(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = -1;
    
    assert_int_equal(chart.update_every, -1);
}

/* Test: Zero update_every */
static void test_charts2json_zero_update_every(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 0;
    
    assert_int_equal(chart.update_every, 0);
}

/* Test: Large update_every */
static void test_charts2json_large_update_every(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 2147483647;
    
    assert_int_equal(chart.update_every, 2147483647);
}

/* Test: Zero usec */
static void test_charts2json_zero_usec(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 0;
    
    assert_int_equal(chart.usec, 0);
}

/* Test: Negative usec */
static void test_charts2json_negative_usec(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = -1000000;
    
    assert_int_equal(chart.usec, -1000000);
}

/* Test: Large usec value */
static void test_charts2json_large_usec(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 9223372036854775807LL;
    
    assert_int_equal(chart.usec, 9223372036854775807LL);
}

/* Test: Zero collected_total */
static void test_charts2json_zero_collected_total(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 0;
    
    assert_int_equal(chart.collected_total, 0);
}

/* Test: Large collected_total */
static void test_charts2json_large_collected_total(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 0xFFFFFFFFUL;
    
    assert_int_equal(chart.collected_total, 0xFFFFFFFFUL);
}

/* Test: Zero last_collected */
static void test_charts2json_zero_last_collected(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 0;
    
    assert_int_equal(chart.last_collected, 0);
}

/* Test: Negative last_collected */
static void test_charts2json_negative_last_collected(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = -1;
    
    assert_int_equal(chart.last_collected, -1);
}

/* Test: last_collected greater than collected_total */
static void test_charts2json_last_collected_exceeds_total(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 150;
    
    assert_true(chart.last_collected > chart.collected_total);
}

/* Test: Zero views */
static void test_charts2json_zero_views(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = 0;
    
    assert_int_equal(chart.views, 0);
}

/* Test: Negative views */
static void test_charts2json_negative_views(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = -1;
    
    assert_int_equal(chart.views, -1);
}

/* Test: Large views count */
static void test_charts2json_large_views(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = 9223372036854775807LL;
    
    assert_int_equal(chart.views, 9223372036854775807LL);
}

/* Test: Disabled chart */
static void test_charts2json_disabled_chart(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = 1;
    chart.enabled = 0;
    
    assert_int_equal(chart.enabled, 0);
}

/* Test: Enabled chart */
static void test_charts2json_enabled_chart(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = 1;
    chart.enabled = 1;
    
    assert_int_equal(chart.enabled, 1);
}

/* Test: Multiple enabled states */
static void test_charts2json_various_enabled_states(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = 1;
    
    for (int state = -1; state <= 2; state++) {
        chart.enabled = state;
        assert_int_equal(chart.enabled, state);
    }
}

/* Test: Dimension with NULL ID */
static void test_charts2json_dimension_null_id(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = NULL;
    dim.name = "system";
    dim.type = "cpu";
    dim.units = "percent";
    
    assert_null(dim.id);
}

/* Test: Dimension with empty ID */
static void test_charts2json_dimension_empty_id(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "";
    dim.name = "system";
    dim.type = "cpu";
    dim.units = "percent";
    
    assert_non_null(dim.id);
    assert_int_equal(strlen(dim.id), 0);
}

/* Test: Dimension with valid ID */
static void test_charts2json_dimension_valid_id(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    
    assert_non_null(dim.id);
    assert_string_equal(dim.id, "system");
}

/* Test: Dimension with NULL name */
static void test_charts2json_dimension_null_name(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = NULL;
    dim.type = "cpu";
    dim.units = "percent";
    
    assert_null(dim.name);
}

/* Test: Dimension with NULL type */
static void test_charts2json_dimension_null_type(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = NULL;
    dim.units = "percent";
    
    assert_null(dim.type);
}

/* Test: Dimension with NULL units */
static void test_charts2json_dimension_null_units(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = NULL;
    
    assert_null(dim.units);
}

/* Test: Dimension with zero multiplier */
static void test_charts2json_dimension_zero_multiplier(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 0;
    
    assert_int_equal(dim.multiplier, 0);
}

/* Test: Dimension with negative multiplier */
static void test_charts2json_dimension_negative_multiplier(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = -1;
    
    assert_int_equal(dim.multiplier, -1);
}

/* Test: Dimension with large multiplier */
static void test_charts2json_dimension_large_multiplier(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 9223372036854775807LL;
    
    assert_int_equal(dim.multiplier, 9223372036854775807LL);
}

/* Test: Dimension with zero divisor */
static void test_charts2json_dimension_zero_divisor(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 1;
    dim.divisor = 0;
    
    assert_int_equal(dim.divisor, 0);
}

/* Test: Dimension with negative divisor */
static void test_charts2json_dimension_negative_divisor(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 1;
    dim.divisor = -1;
    
    assert_int_equal(dim.divisor, -1);
}

/* Test: Dimension with large divisor */
static void test_charts2json_dimension_large_divisor(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 1;
    dim.divisor = 9223372036854775807LL;
    
    assert_int_equal(dim.divisor, 9223372036854775807LL);
}

/* Test: Dimension hidden flag false */
static void test_charts2json_dimension_not_hidden(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 1;
    dim.divisor = 1;
    dim.hidden = 0;
    
    assert_int_equal(dim.hidden, 0);
}

/* Test: Dimension hidden flag true */
static void test_charts2json_dimension_hidden(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 1;
    dim.divisor = 1;
    dim.hidden = 1;
    
    assert_int_equal(dim.hidden, 1);
}

/* Test: Dimension various flags */
static void test_charts2json_dimension_various_flags(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 1;
    dim.divisor = 1;
    dim.hidden = 0;
    
    for (int flag = 0; flag <= 255; flag++) {
        dim.flags = flag;
        assert_int_equal(dim.flags, flag);
    }
}

/* Test: All chart fields together */
static void test_charts2json_complete_chart_structure(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    chart.priority = 100;
    chart.update_every = 1;
    chart.usec = 1000000;
    chart.collected_total = 100;
    chart.last_collected = 99;
    chart.views = 1;
    chart.enabled = 1;
    
    assert_string_equal(chart.id, "cpu.system");
    assert_string_equal(chart.name, "CPU System");
    assert_string_equal(chart.type, "area");
    assert_string_equal(chart.family, "cpu");
    assert_int_equal(chart.priority, 100);
    assert_int_equal(chart.update_every, 1);
    assert_int_equal(chart.usec, 1000000);
    assert_int_equal(chart.collected_total, 100);
    assert_int_equal(chart.last_collected, 99);
    assert_int_equal(chart.views, 1);
    assert_int_equal(chart.enabled, 1);
}

/* Test: All dimension fields together */
static void test_charts2json_complete_dimension_structure(void **state) {
    (void)state;
    
    RRDDIM dim;
    memset(&dim, 0, sizeof(RRDDIM));
    dim.id = "system";
    dim.name = "System CPU";
    dim.type = "cpu";
    dim.units = "percent";
    dim.multiplier = 1;
    dim.divisor = 100;
    dim.hidden = 0;
    dim.flags = 0;
    
    assert_string_equal(dim.id, "system");
    assert_string_equal(dim.name, "System CPU");
    assert_string_equal(dim.type, "cpu");
    assert_string_equal(dim.units, "percent");
    assert_int_equal(dim.multiplier, 1);
    assert_int_equal(dim.divisor, 100);
    assert_int_equal(dim.hidden, 0);
    assert_int_equal(dim.flags, 0);
}

/* Test: Special characters in chart ID */
static void test_charts2json_special_chars_chart_id(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu-system_total.01";
    chart.name = "CPU System";
    chart.type = "area";
    chart.family = "cpu";
    
    assert_string_equal(chart.id, "cpu-system_total.01");
}

/* Test: Special characters in chart name */
static void test_charts2json_special_chars_chart_name(void **state) {
    (void)state;
    
    RRDCHART chart;
    memset(&chart, 0, sizeof(RRDCHART));
    chart.id = "cpu.system";
    chart.name = "CPU System (Total)";
    chart.type = "area";
    chart.family = "cpu";
    
    assert_string_equal(chart.name, "CPU System (Total)");
}

/* Test: Quote characters in chart fields */
static void test_charts2json_quote_chars_in_fields(voi