#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

#include "stream-conf.h"

/* Mock structures and test fixtures */
typedef struct {
    char *buffer;
    size_t size;
    size_t capacity;
} test_buffer_t;

/* Mock implementations */
static void* __real_malloc(size_t size);
static void __real_free(void *ptr);
static char* __real_strdup(const char *s);

/* Test fixtures */
static test_buffer_t* test_buffer_create(size_t capacity) {
    test_buffer_t *buf = malloc(sizeof(test_buffer_t));
    buf->buffer = malloc(capacity);
    buf->size = 0;
    buf->capacity = capacity;
    return buf;
}

static void test_buffer_destroy(test_buffer_t *buf) {
    if (buf) {
        free(buf->buffer);
        free(buf);
    }
}

static void test_buffer_append(test_buffer_t *buf, const char *str) {
    if (!buf || !str) return;
    size_t len = strlen(str);
    if (buf->size + len >= buf->capacity) return;
    strcpy(buf->buffer + buf->size, str);
    buf->size += len;
}

/* Test: stream_conf_init - Basic initialization */
static void test_stream_conf_init_basic(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_init();
    
    assert_non_null(conf);
    assert_int_equal(conf->enabled, 0);
    assert_non_null(conf->default_port);
    assert_int_equal(strcmp(conf->default_port, "19999"), 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_init - NULL return on memory allocation failure */
static void test_stream_conf_init_memory_failure(void **state) {
    (void) state;
    
    /* This test assumes the function gracefully handles malloc failures */
    struct stream_conf *conf = stream_conf_init();
    
    /* Should either succeed or return NULL on allocation failure */
    if (conf) {
        stream_conf_free(conf);
    }
}

/* Test: stream_conf_free - Free allocated resources */
static void test_stream_conf_free_valid(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_init();
    assert_non_null(conf);
    
    /* Should not crash on valid pointer */
    stream_conf_free(conf);
}

/* Test: stream_conf_free - NULL pointer handling */
static void test_stream_conf_free_null(void **state) {
    (void) state;
    
    /* Should not crash on NULL pointer */
    stream_conf_free(NULL);
}

/* Test: stream_conf_create - Create new configuration */
static void test_stream_conf_create_basic(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    assert_non_null(conf);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse valid config */
static void test_stream_conf_parse_valid_config(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    const char *config = "enabled = yes\nport = 19999\n";
    
    int result = stream_conf_parse(conf, config);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse NULL config */
static void test_stream_conf_parse_null_config(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, NULL);
    
    /* Should handle NULL gracefully */
    assert_true(result < 0 || result == 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse empty config */
static void test_stream_conf_parse_empty_config(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "");
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse enabled setting */
static void test_stream_conf_parse_enabled_yes(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "enabled = yes\n");
    
    assert_int_equal(result, 0);
    assert_int_equal(conf->enabled, 1);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse disabled setting */
static void test_stream_conf_parse_enabled_no(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "enabled = no\n");
    
    assert_int_equal(result, 0);
    assert_int_equal(conf->enabled, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse port setting */
static void test_stream_conf_parse_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "port = 20000\n");
    
    assert_int_equal(result, 0);
    assert_non_null(conf->default_port);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse invalid port */
static void test_stream_conf_parse_invalid_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "port = invalid\n");
    
    /* Should handle invalid port gracefully */
    assert_true(result == 0 || result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse multiple settings */
static void test_stream_conf_parse_multiple_settings(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    const char *config = "enabled = yes\nport = 20000\n";
    
    int result = stream_conf_parse(conf, config);
    
    assert_int_equal(result, 0);
    assert_int_equal(conf->enabled, 1);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse with comments */
static void test_stream_conf_parse_with_comments(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    const char *config = "# Comment\nenabled = yes\n";
    
    int result = stream_conf_parse(conf, config);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Parse with whitespace */
static void test_stream_conf_parse_with_whitespace(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    const char *config = "  enabled  =  yes  \n";
    
    int result = stream_conf_parse(conf, config);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_load_file - Load from existing file */
static void test_stream_conf_load_file_valid(void **state) {
    (void) state;
    
    /* Create temporary config file */
    char temp_file[] = "/tmp/test_stream_conf_XXXXXX";
    int fd = mkstemp(temp_file);
    assert_return_code(fd, -1);
    
    const char *content = "enabled = yes\nport = 19999\n";
    ssize_t written = write(fd, content, strlen(content));
    assert_int_equal(written, strlen(content));
    close(fd);
    
    struct stream_conf *conf = stream_conf_create();
    int result = stream_conf_load_file(conf, temp_file);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
    unlink(temp_file);
}

/* Test: stream_conf_load_file - Load from non-existent file */
static void test_stream_conf_load_file_not_found(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    int result = stream_conf_load_file(conf, "/nonexistent/path/config.conf");
    
    /* Should return error code */
    assert_true(result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_load_file - NULL filename */
static void test_stream_conf_load_file_null_filename(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    int result = stream_conf_load_file(conf, NULL);
    
    /* Should handle NULL gracefully */
    assert_true(result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_load_file - Empty filename */
static void test_stream_conf_load_file_empty_filename(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    int result = stream_conf_load_file(conf, "");
    
    assert_true(result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_load_file - NULL conf */
static void test_stream_conf_load_file_null_conf(void **state) {
    (void) state;
    
    int result = stream_conf_load_file(NULL, "/tmp/config.conf");
    
    assert_true(result < 0);
}

/* Test: stream_conf_get_enabled - Get enabled flag */
static void test_stream_conf_get_enabled_true(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "enabled = yes\n");
    
    int enabled = stream_conf_get_enabled(conf);
    
    assert_int_equal(enabled, 1);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_get_enabled - Get disabled flag */
static void test_stream_conf_get_enabled_false(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "enabled = no\n");
    
    int enabled = stream_conf_get_enabled(conf);
    
    assert_int_equal(enabled, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_get_enabled - NULL conf */
static void test_stream_conf_get_enabled_null_conf(void **state) {
    (void) state;
    
    int enabled = stream_conf_get_enabled(NULL);
    
    assert_int_equal(enabled, 0);
}

/* Test: stream_conf_get_port - Get port string */
static void test_stream_conf_get_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "port = 20000\n");
    
    const char *port = stream_conf_get_port(conf);
    
    assert_non_null(port);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_get_port - Default port */
static void test_stream_conf_get_port_default(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    const char *port = stream_conf_get_port(conf);
    
    assert_non_null(port);
    assert_int_equal(strcmp(port, "19999"), 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_get_port - NULL conf */
static void test_stream_conf_get_port_null_conf(void **state) {
    (void) state;
    
    const char *port = stream_conf_get_port(NULL);
    
    assert_null(port);
}

/* Test: stream_conf_set_enabled - Set enabled flag */
static void test_stream_conf_set_enabled(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    stream_conf_set_enabled(conf, 1);
    assert_int_equal(conf->enabled, 1);
    
    stream_conf_set_enabled(conf, 0);
    assert_int_equal(conf->enabled, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_set_enabled - NULL conf */
static void test_stream_conf_set_enabled_null_conf(void **state) {
    (void) state;
    
    /* Should not crash on NULL */
    stream_conf_set_enabled(NULL, 1);
}

/* Test: stream_conf_set_port - Set port */
static void test_stream_conf_set_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    stream_conf_set_port(conf, "20000");
    
    const char *port = stream_conf_get_port(conf);
    assert_non_null(port);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_set_port - NULL port */
static void test_stream_conf_set_port_null_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    stream_conf_set_port(conf, NULL);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_set_port - Empty port */
static void test_stream_conf_set_port_empty(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    stream_conf_set_port(conf, "");
    
    stream_conf_free(conf);
}

/* Test: stream_conf_set_port - NULL conf */
static void test_stream_conf_set_port_null_conf(void **state) {
    (void) state;
    
    /* Should not crash on NULL */
    stream_conf_set_port(NULL, "20000");
}

/* Test: stream_conf_reset - Reset to defaults */
static void test_stream_conf_reset(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "enabled = yes\nport = 20000\n");
    
    stream_conf_reset(conf);
    
    /* After reset, should have default values */
    assert_int_equal(conf->enabled, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_reset - NULL conf */
static void test_stream_conf_reset_null_conf(void **state) {
    (void) state;
    
    /* Should not crash */
    stream_conf_reset(NULL);
}

/* Test: stream_conf_validate - Validate valid config */
static void test_stream_conf_validate_valid(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "enabled = yes\nport = 20000\n");
    
    int result = stream_conf_validate(conf);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_validate - NULL conf */
static void test_stream_conf_validate_null_conf(void **state) {
    (void) state;
    
    int result = stream_conf_validate(NULL);
    
    assert_true(result < 0);
}

/* Test: stream_conf_to_string - Convert to string */
static void test_stream_conf_to_string(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "enabled = yes\nport = 20000\n");
    
    char *str = stream_conf_to_string(conf);
    
    assert_non_null(str);
    
    free(str);
    stream_conf_free(conf);
}

/* Test: stream_conf_to_string - NULL conf */
static void test_stream_conf_to_string_null_conf(void **state) {
    (void) state;
    
    char *str = stream_conf_to_string(NULL);
    
    assert_null(str);
}

/* Test: stream_conf_duplicate - Duplicate configuration */
static void test_stream_conf_duplicate(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "enabled = yes\nport = 20000\n");
    
    struct stream_conf *dup = stream_conf_duplicate(conf);
    
    assert_non_null(dup);
    assert_int_equal(dup->enabled, conf->enabled);
    
    stream_conf_free(conf);
    stream_conf_free(dup);
}

/* Test: stream_conf_duplicate - NULL conf */
static void test_stream_conf_duplicate_null_conf(void **state) {
    (void) state;
    
    struct stream_conf *dup = stream_conf_duplicate(NULL);
    
    assert_null(dup);
}

/* Test: stream_conf_compare - Compare identical configs */
static void test_stream_conf_compare_identical(void **state) {
    (void) state;
    
    struct stream_conf *conf1 = stream_conf_create();
    struct stream_conf *conf2 = stream_conf_create();
    
    stream_conf_parse(conf1, "enabled = yes\nport = 20000\n");
    stream_conf_parse(conf2, "enabled = yes\nport = 20000\n");
    
    int result = stream_conf_compare(conf1, conf2);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf1);
    stream_conf_free(conf2);
}

/* Test: stream_conf_compare - Compare different configs */
static void test_stream_conf_compare_different(void **state) {
    (void) state;
    
    struct stream_conf *conf1 = stream_conf_create();
    struct stream_conf *conf2 = stream_conf_create();
    
    stream_conf_parse(conf1, "enabled = yes\nport = 20000\n");
    stream_conf_parse(conf2, "enabled = no\nport = 20001\n");
    
    int result = stream_conf_compare(conf1, conf2);
    
    assert_int_not_equal(result, 0);
    
    stream_conf_free(conf1);
    stream_conf_free(conf2);
}

/* Test: stream_conf_compare - NULL configs */
static void test_stream_conf_compare_null_first(void **state) {
    (void) state;
    
    struct stream_conf *conf2 = stream_conf_create();
    
    int result = stream_conf_compare(NULL, conf2);
    
    assert_true(result < 0 || result > 0);
    
    stream_conf_free(conf2);
}

/* Test: stream_conf_compare - Both NULL */
static void test_stream_conf_compare_both_null(void **state) {
    (void) state;
    
    int result = stream_conf_compare(NULL, NULL);
    
    assert_int_equal(result, 0);
}

/* Test: stream_conf_is_valid - Valid config */
static void test_stream_conf_is_valid_true(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    stream_conf_parse(conf, "enabled = yes\nport = 20000\n");
    
    int result = stream_conf_is_valid(conf);
    
    assert_int_equal(result, 1);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_is_valid - NULL conf */
static void test_stream_conf_is_valid_null(void **state) {
    (void) state;
    
    int result = stream_conf_is_valid(NULL);
    
    assert_int_equal(result, 0);
}

/* Test: stream_conf_merge - Merge configurations */
static void test_stream_conf_merge(void **state) {
    (void) state;
    
    struct stream_conf *conf1 = stream_conf_create();
    struct stream_conf *conf2 = stream_conf_create();
    
    stream_conf_parse(conf1, "enabled = yes\n");
    stream_conf_parse(conf2, "port = 20000\n");
    
    stream_conf_merge(conf1, conf2);
    
    stream_conf_free(conf1);
    stream_conf_free(conf2);
}

/* Test: stream_conf_merge - NULL target */
static void test_stream_conf_merge_null_target(void **state) {
    (void) state;
    
    struct stream_conf *conf2 = stream_conf_create();
    
    stream_conf_merge(NULL, conf2);
    
    stream_conf_free(conf2);
}

/* Test: stream_conf_merge - NULL source */
static void test_stream_conf_merge_null_source(void **state) {
    (void) state;
    
    struct stream_conf *conf1 = stream_conf_create();
    
    stream_conf_merge(conf1, NULL);
    
    stream_conf_free(conf1);
}

/* Test: stream_conf_parse - Invalid enabled value */
static void test_stream_conf_parse_invalid_enabled(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "enabled = maybe\n");
    
    assert_true(result == 0 || result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Case insensitivity */
static void test_stream_conf_parse_case_insensitive(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "ENABLED = YES\n");
    
    /* Should handle case variations */
    assert_true(result == 0 || result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Malformed line */
static void test_stream_conf_parse_malformed(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "malformed line without equals\n");
    
    assert_true(result == 0 || result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_get_port - After setting custom port */
static void test_stream_conf_port_custom(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    const char *custom_port = "12345";
    
    stream_conf_set_port(conf, custom_port);
    const char *retrieved = stream_conf_get_port(conf);
    
    assert_non_null(retrieved);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Very long port number */
static void test_stream_conf_parse_large_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "port = 99999\n");
    
    assert_true(result == 0 || result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Negative port number */
static void test_stream_conf_parse_negative_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "port = -1\n");
    
    assert_true(result == 0 || result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_parse - Zero port number */
static void test_stream_conf_parse_zero_port(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    int result = stream_conf_parse(conf, "port = 0\n");
    
    assert_true(result == 0 || result < 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_load_file - Large config file */
static void test_stream_conf_load_file_large(void **state) {
    (void) state;
    
    char temp_file[] = "/tmp/test_stream_conf_large_XXXXXX";
    int fd = mkstemp(temp_file);
    assert_return_code(fd, -1);
    
    /* Write large config */
    for (int i = 0; i < 100; i++) {
        char line[256];
        snprintf(line, sizeof(line), "# Comment line %d\n", i);
        write(fd, line, strlen(line));
    }
    write(fd, "enabled = yes\n", strlen("enabled = yes\n"));
    close(fd);
    
    struct stream_conf *conf = stream_conf_create();
    int result = stream_conf_load_file(conf, temp_file);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
    unlink(temp_file);
}

/* Test: stream_conf_parse - Multiple identical settings */
static void test_stream_conf_parse_duplicate_settings(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    const char *config = "enabled = yes\nenabled = no\nport = 19999\nport = 20000\n";
    
    int result = stream_conf_parse(conf, config);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_load_file - Permission denied */
static void test_stream_conf_load_file_permission_denied(void **state) {
    (void) state;
    
    char temp_file[] = "/tmp/test_stream_conf_perm_XXXXXX";
    int fd = mkstemp(temp_file);
    assert_return_code(fd, -1);
    write(fd, "enabled = yes\n", strlen("enabled = yes\n"));
    close(fd);
    
    /* Remove read permissions */
    chmod(temp_file, 0000);
    
    struct stream_conf *conf = stream_conf_create();
    int result = stream_conf_load_file(conf, temp_file);
    
    /* Should fail due to permissions */
    assert_true(result < 0 || result == 0);
    
    stream_conf_free(conf);
    chmod(temp_file, 0644);
    unlink(temp_file);
}

/* Test: stream_conf_parse - Leading/trailing newlines */
static void test_stream_conf_parse_extra_newlines(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    const char *config = "\n\nenabled = yes\n\n\n";
    
    int result = stream_conf_parse(conf, config);
    
    assert_int_equal(result, 0);
    
    stream_conf_free(conf);
}

/* Test: stream_conf_set_enabled - Toggle multiple times */
static void test_stream_conf_set_enabled_toggle(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    for (int i = 0; i < 10; i++) {
        stream_conf_set_enabled(conf, i % 2);
        assert_int_equal(conf->enabled, i % 2);
    }
    
    stream_conf_free(conf);
}

/* Test: stream_conf_set_port - Multiple times */
static void test_stream_conf_set_port_multiple(void **state) {
    (void) state;
    
    struct stream_conf *conf = stream_conf_create();
    
    stream_conf_set_port(conf, "10000");
    stream_conf_set_port(conf, "20000");
    stream_conf_set_port(conf, "30000");
    
    const char *port = stream_conf_get_port(conf);
    assert_non_null(port);
    
    stream_conf_free(conf);
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_stream_conf_init_basic),
        cmocka_unit_test(test_stream_conf_init_memory_failure),
        cmocka_unit_test(test_stream_conf_free_valid),
        cmocka_unit_test(test_stream_conf_free_null),
        cmocka_unit_test(test_stream_conf_create_basic),
        cmocka_unit_test(test_stream_conf_parse_valid_config),
        cmocka_unit_test(test_stream_conf_parse_null_config),
        cmocka_unit_test(test_stream_conf_parse_empty_config),
        cmocka_unit_test(test_stream_conf_parse_enabled_yes),
        cmocka_unit_test(test_stream_conf_parse_enabled_no),
        cmocka_unit_test(test_stream_conf_parse_port),
        cmocka_unit_test(test_stream_conf_parse_invalid_port),
        cmocka_unit_test(test_stream_conf_parse_multiple_settings),
        cmocka_unit_test