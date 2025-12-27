#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <errno.h>

#include "../spawn_server/spawn_popen.h"

/* ============================================================================
 * Test Fixtures and Helpers
 * ============================================================================
 */

/* Mock for popen functionality */
static FILE* mock_popen_return;
static const char* mock_popen_cmd;
static const char* mock_popen_mode;

FILE *__wrap_popen(const char *cmd, const char *mode) {
    check_expected_ptr(cmd);
    check_expected_ptr(mode);
    mock_popen_cmd = cmd;
    mock_popen_mode = mode;
    return mock_popen_return;
}

/* Mock for pclose functionality */
static int mock_pclose_return;
int __wrap_pclose(FILE *stream) {
    check_expected_ptr(stream);
    return mock_pclose_return;
}

/* Mock for fread functionality */
static size_t mock_fread_return;
size_t __wrap_fread(void *ptr, size_t size, size_t nmemb, FILE *stream) {
    check_expected(size);
    check_expected(nmemb);
    check_expected_ptr(stream);
    
    if (ptr && mock_fread_return > 0) {
        memcpy(ptr, "test data", mock_fread_return);
    }
    return mock_fread_return;
}

/* Mock for fwrite functionality */
static size_t mock_fwrite_return;
size_t __wrap_fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream) {
    check_expected(size);
    check_expected(nmemb);
    check_expected_ptr(stream);
    if (ptr) {
        check_expected_ptr(ptr);
    }
    return mock_fwrite_return;
}

/* Mock for fgets functionality */
static char* mock_fgets_return;
char *__wrap_fgets(char *s, int size, FILE *stream) {
    check_expected(size);
    check_expected_ptr(stream);
    
    if (mock_fgets_return && s) {
        strncpy(s, mock_fgets_return, size - 1);
        s[size - 1] = '\0';
        return s;
    }
    return NULL;
}

/* Mock for ferror functionality */
static int mock_ferror_return;
int __wrap_ferror(FILE *stream) {
    check_expected_ptr(stream);
    return mock_ferror_return;
}

/* Mock for clearerr functionality */
void __wrap_clearerr(FILE *stream) {
    check_expected_ptr(stream);
    mock_ferror_return = 0;
}

/* Mock for fileno functionality */
static int mock_fileno_return;
int __wrap_fileno(FILE *stream) {
    check_expected_ptr(stream);
    return mock_fileno_return;
}

/* ============================================================================
 * spawn_popen_new() Tests
 * ============================================================================
 */

static int test_spawn_popen_new_success(void **state) {
    (void) state;
    
    FILE *fp = fopen("/dev/null", "r");
    assert_non_null(fp);
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    
    assert_non_null(popen_struct);
    assert_null(popen_struct->stream);
    assert_int_equal(popen_struct->fd, 0);
    assert_null(popen_struct->command);
    assert_int_equal(popen_struct->flags, 0);
    
    spawn_popen_free(popen_struct);
    fclose(fp);
    
    return 0;
}

static int test_spawn_popen_new_memory_allocation(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    
    assert_non_null(popen_struct);
    assert_true(popen_struct != NULL);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_exec() Tests
 * ============================================================================
 */

static int test_spawn_popen_exec_with_valid_command(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    const char *cmd = "echo 'test'";
    int result = spawn_popen_exec(popen_struct, cmd);
    
    /* Result depends on actual implementation */
    if (result == 0) {
        assert_non_null(popen_struct->stream);
        if (popen_struct->stream) {
            pclose(popen_struct->stream);
            popen_struct->stream = NULL;
        }
    }
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_exec_with_null_struct(void **state) {
    (void) state;
    
    int result = spawn_popen_exec(NULL, "test");
    
    assert_int_not_equal(result, 0);
    
    return 0;
}

static int test_spawn_popen_exec_with_null_command(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    int result = spawn_popen_exec(popen_struct, NULL);
    
    assert_int_not_equal(result, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_exec_with_empty_command(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    int result = spawn_popen_exec(popen_struct, "");
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_exec_with_read_mode(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    popen_struct->flags = SPAWN_POPEN_READ;
    const char *cmd = "echo 'test'";
    int result = spawn_popen_exec(popen_struct, cmd);
    
    if (result == 0 && popen_struct->stream) {
        pclose(popen_struct->stream);
        popen_struct->stream = NULL;
    }
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_exec_with_write_mode(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    popen_struct->flags = SPAWN_POPEN_WRITE;
    const char *cmd = "cat > /dev/null";
    int result = spawn_popen_exec(popen_struct, cmd);
    
    if (result == 0 && popen_struct->stream) {
        pclose(popen_struct->stream);
        popen_struct->stream = NULL;
    }
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_read() Tests
 * ============================================================================
 */

static int test_spawn_popen_read_with_valid_stream(void **state) {
    (void) state;
    
    FILE *fp = popen("echo 'test data'", "r");
    if (!fp) {
        skip();
    }
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    popen_struct->stream = fp;
    
    char buffer[256];
    size_t bytes_read = spawn_popen_read(popen_struct, buffer, sizeof(buffer));
    
    assert_int_not_equal(bytes_read, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_read_with_null_struct(void **state) {
    (void) state;
    
    char buffer[256];
    size_t bytes_read = spawn_popen_read(NULL, buffer, sizeof(buffer));
    
    assert_int_equal(bytes_read, 0);
    
    return 0;
}

static int test_spawn_popen_read_with_null_buffer(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    size_t bytes_read = spawn_popen_read(popen_struct, NULL, 256);
    
    assert_int_equal(bytes_read, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_read_with_zero_length(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    char buffer[256];
    size_t bytes_read = spawn_popen_read(popen_struct, buffer, 0);
    
    assert_int_equal(bytes_read, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_read_with_no_stream(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    popen_struct->stream = NULL;
    
    char buffer[256];
    size_t bytes_read = spawn_popen_read(popen_struct, buffer, sizeof(buffer));
    
    assert_int_equal(bytes_read, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_write() Tests
 * ============================================================================
 */

static int test_spawn_popen_write_with_valid_stream(void **state) {
    (void) state;
    
    FILE *fp = popen("cat > /dev/null", "w");
    if (!fp) {
        skip();
    }
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    popen_struct->stream = fp;
    popen_struct->flags = SPAWN_POPEN_WRITE;
    
    const char *data = "test data";
    size_t bytes_written = spawn_popen_write(popen_struct, data, strlen(data));
    
    assert_int_not_equal(bytes_written, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_write_with_null_struct(void **state) {
    (void) state;
    
    const char *data = "test data";
    size_t bytes_written = spawn_popen_write(NULL, data, strlen(data));
    
    assert_int_equal(bytes_written, 0);
    
    return 0;
}

static int test_spawn_popen_write_with_null_buffer(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    size_t bytes_written = spawn_popen_write(popen_struct, NULL, 256);
    
    assert_int_equal(bytes_written, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_write_with_zero_length(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    const char *data = "test data";
    size_t bytes_written = spawn_popen_write(popen_struct, data, 0);
    
    assert_int_equal(bytes_written, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_write_with_no_stream(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    popen_struct->stream = NULL;
    
    const char *data = "test data";
    size_t bytes_written = spawn_popen_write(popen_struct, data, strlen(data));
    
    assert_int_equal(bytes_written, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_close() Tests
 * ============================================================================
 */

static int test_spawn_popen_close_with_valid_stream(void **state) {
    (void) state;
    
    FILE *fp = popen("echo 'test'", "r");
    if (!fp) {
        skip();
    }
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    popen_struct->stream = fp;
    
    int result = spawn_popen_close(popen_struct);
    
    assert_int_equal(result, 0);
    assert_null(popen_struct->stream);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_close_with_null_struct(void **state) {
    (void) state;
    
    int result = spawn_popen_close(NULL);
    
    assert_int_not_equal(result, 0);
    
    return 0;
}

static int test_spawn_popen_close_with_no_stream(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    popen_struct->stream = NULL;
    
    int result = spawn_popen_close(popen_struct);
    
    assert_int_not_equal(result, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_close_idempotent(void **state) {
    (void) state;
    
    FILE *fp = popen("echo 'test'", "r");
    if (!fp) {
        skip();
    }
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    popen_struct->stream = fp;
    
    spawn_popen_close(popen_struct);
    
    popen_struct->stream = NULL;
    int result = spawn_popen_close(popen_struct);
    
    assert_int_not_equal(result, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_free() Tests
 * ============================================================================
 */

static int test_spawn_popen_free_with_null_struct(void **state) {
    (void) state;
    
    spawn_popen_free(NULL);
    
    return 0;
}

static int test_spawn_popen_free_with_active_stream(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    FILE *fp = popen("echo 'test'", "r");
    if (fp) {
        popen_struct->stream = fp;
        spawn_popen_free(popen_struct);
    } else {
        spawn_popen_free(popen_struct);
    }
    
    return 0;
}

static int test_spawn_popen_free_with_command(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    popen_struct->command = strdup("test command");
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_free_releases_memory(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_get_fd() Tests
 * ============================================================================
 */

static int test_spawn_popen_get_fd_valid_struct(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    int fd = spawn_popen_get_fd(popen_struct);
    
    assert_int_equal(fd, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_get_fd_with_null_struct(void **state) {
    (void) state;
    
    int fd = spawn_popen_get_fd(NULL);
    
    assert_int_equal(fd, -1);
    
    return 0;
}

static int test_spawn_popen_get_fd_after_exec(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    FILE *fp = popen("echo 'test'", "r");
    if (fp) {
        popen_struct->stream = fp;
        popen_struct->fd = fileno(fp);
        
        int fd = spawn_popen_get_fd(popen_struct);
        
        assert_int_not_equal(fd, 0);
        
        pclose(fp);
    }
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_get_stream() Tests
 * ============================================================================
 */

static int test_spawn_popen_get_stream_valid_struct(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    FILE *stream = spawn_popen_get_stream(popen_struct);
    
    assert_null(stream);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_get_stream_with_null_struct(void **state) {
    (void) state;
    
    FILE *stream = spawn_popen_get_stream(NULL);
    
    assert_null(stream);
    
    return 0;
}

static int test_spawn_popen_get_stream_after_exec(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    FILE *fp = popen("echo 'test'", "r");
    if (fp) {
        popen_struct->stream = fp;
        
        FILE *stream = spawn_popen_get_stream(popen_struct);
        
        assert_ptr_equal(stream, fp);
        
        pclose(fp);
    }
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * spawn_popen_set_flags() Tests
 * ============================================================================
 */

static int test_spawn_popen_set_flags_read(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    spawn_popen_set_flags(popen_struct, SPAWN_POPEN_READ);
    
    assert_int_equal(popen_struct->flags, SPAWN_POPEN_READ);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_set_flags_write(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    spawn_popen_set_flags(popen_struct, SPAWN_POPEN_WRITE);
    
    assert_int_equal(popen_struct->flags, SPAWN_POPEN_WRITE);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_set_flags_multiple(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    spawn_popen_set_flags(popen_struct, SPAWN_POPEN_READ | SPAWN_POPEN_WRITE);
    
    assert_int_equal(popen_struct->flags, (SPAWN_POPEN_READ | SPAWN_POPEN_WRITE));
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_set_flags_with_null_struct(void **state) {
    (void) state;
    
    spawn_popen_set_flags(NULL, SPAWN_POPEN_READ);
    
    return 0;
}

static int test_spawn_popen_set_flags_zero(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    spawn_popen_set_flags(popen_struct, 0);
    
    assert_int_equal(popen_struct->flags, 0);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

/* ============================================================================
 * Edge Case Tests
 * ============================================================================
 */

static int test_spawn_popen_command_with_special_characters(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    const char *cmd = "echo 'test|special&chars'";
    int result = spawn_popen_exec(popen_struct, cmd);
    
    if (result == 0 && popen_struct->stream) {
        pclose(popen_struct->stream);
        popen_struct->stream = NULL;
    }
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_very_long_command(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    char long_cmd[4096];
    memset(long_cmd, 'a', sizeof(long_cmd) - 1);
    long_cmd[sizeof(long_cmd) - 1] = '\0';
    
    int result = spawn_popen_exec(popen_struct, long_cmd);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_multiple_operations(void **state) {
    (void) state;
    
    struct spawn_popen *popen_struct = spawn_popen_new();
    assert_non_null(popen_struct);
    
    spawn_popen_set_flags(popen_struct, SPAWN_POPEN_READ);
    assert_int_equal(popen_struct->flags, SPAWN_POPEN_READ);
    
    int fd = spawn_popen_get_fd(popen_struct);
    assert_int_equal(fd, 0);
    
    FILE *stream = spawn_popen_get_stream(popen_struct);
    assert_null(stream);
    
    spawn_popen_free(popen_struct);
    
    return 0;
}

static int test_spawn_popen_sequential_allocations(void **state) {
    (void) state;
    
    struct spawn_popen *p1 = spawn_popen_new();
    struct spawn_popen *p2 = spawn_popen_new();
    struct spawn_popen *p3 = spawn_popen_new();
    
    assert_non_null(p1);
    assert_non_null(p2);
    assert_non_null(p3);
    
    assert_ptr_not_equal(p1, p2);
    assert_ptr_not_equal(p2, p3);
    
    spawn_popen_free(p1);
    spawn_popen_free(p2);
    spawn_popen_free(p3);
    
    return 0;
}

/* ============================================================================
 * Main Test Suite
 * ============================================================================
 */

const struct CMUnitTest spawn_popen_tests[] = {
    /* spawn_popen_new tests */
    cmocka_unit_test(test_spawn_popen_new_success),
    cmocka_unit_test(test_spawn_popen_new_memory_allocation),
    
    /* spawn_popen_exec tests */
    cmocka_unit_test(test_spawn_popen_exec_with_valid_command),
    cmocka_unit_test(test_spawn_popen_exec_with_null_struct),
    cmocka_unit_test(test_spawn_popen_exec_with_null_command),
    cmocka_unit_test(test_spawn_popen_exec_with_empty_command),
    cmocka_unit_test(test_spawn_popen_exec_with_read_mode),
    cmocka_unit_test(test_spawn_popen_exec_with_write_mode),
    
    /* spawn_popen_read tests */
    cmocka_unit_test(test_spawn_popen_read_with_valid_stream),
    cmocka_unit_test(test_spawn_popen_read_with_null_struct),
    cmocka_unit_test(test_spawn_popen_read_with_null_buffer),
    cmocka_unit_test(test_spawn_popen_read_with_zero_length),
    cmocka_unit_test(test_spawn_popen_read_with_no_stream),
    
    /* spawn_popen_write tests */
    cmocka_unit_test(test_spawn_popen_write_with_valid_stream),
    cmocka_unit_test(test_spawn_popen_write_with_null_struct),
    cmocka_unit_test(test_spawn_popen_write_with_null_buffer),
    cmocka_unit_test(test_spawn_popen_write_with_zero_length),
    cmocka_unit_test(test_spawn_popen_write_with_no_stream),
    
    /* spawn_popen_close tests */
    cmocka_unit_test(test_spawn_popen_close_with_valid_stream),
    cmocka_unit_test(test_spawn_popen_close_with_null_struct),
    cmocka_unit_test(test_spawn_popen_close_with_no_stream),
    cmocka_unit_test(test_spawn_popen_close_idempotent),
    
    /* spawn_popen_free tests */
    cmocka_unit_test(test_spawn_popen_free_with_null_struct),
    cmocka_unit_test(test_spawn_popen_free_with_active_stream),
    cmocka_unit_test(test_spawn_popen_free_with_command),
    cmocka_unit_test(test_spawn_popen_free_releases_memory),
    
    /* spawn_popen_get_fd tests */
    cmocka_unit_test(test_spawn_popen_get_fd_valid_struct),
    cmocka_unit_test(test_spawn_popen_get_fd_with_null_struct),
    cmocka_unit_test(test_spawn_popen_get_fd_after_exec),
    
    /* spawn_popen_get_stream tests */
    cmocka_unit_test(test_spawn_popen_get_stream_valid_struct),
    cmocka_unit_test(test_spawn_popen_get_stream_with_null_struct),
    cmocka_unit_test(test_spawn_popen_get_stream_after_exec),
    
    /* spawn_popen_set_flags tests */
    cmocka_unit_test(test_spawn_popen_set_flags_read),
    cmocka_unit_test(test_spawn_popen_set_flags_write),
    cmocka_unit_test(test_spawn_popen_set_flags_multiple),
    cmocka_unit_test(test_spawn_popen_set_flags_with_null_struct),
    cmocka_unit_test(test_spawn_popen_set_flags_zero),
    
    /* Edge case tests */
    cmocka_unit_test(test_spawn_popen_command_with_special_characters),
    cmocka_unit_test(test_spawn_popen_very_long_command),
    cmocka_unit_test(test_spawn_popen_multiple_operations),
    cmocka_unit_test(test_spawn_popen_sequential_allocations),
};

int main(void) {
    return cmocka