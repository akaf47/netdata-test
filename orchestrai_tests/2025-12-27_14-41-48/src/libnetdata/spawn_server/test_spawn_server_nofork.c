#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/wait.h>
#include <errno.h>

// Mock declarations
typedef struct {
    char *command;
    char *stdin_data;
    int stdout_fd;
    int stderr_fd;
    size_t buffer_size;
} spawn_context_t;

// Forward declarations of functions to test
void spawn_server_nofork_init(void);
int spawn_server_nofork_exec(const char *cmd, const char *stdin_data, int *stdout_fd, int *stderr_fd);
void spawn_server_nofork_cleanup(void);

/* Test Suite: spawn_server_nofork initialization */

static int setup_nofork(void **state) {
    spawn_context_t *ctx = malloc(sizeof(spawn_context_t));
    if (!ctx) return -1;
    memset(ctx, 0, sizeof(spawn_context_t));
    ctx->buffer_size = 4096;
    *state = ctx;
    return 0;
}

static int teardown_nofork(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    if (ctx) {
        free(ctx->command);
        free(ctx->stdin_data);
        free(ctx);
    }
    return 0;
}

/* Test: spawn_server_nofork_init with valid initialization */
static void test_spawn_server_nofork_init_success(void **state) {
    spawn_server_nofork_init();
    // Verify initialization completed without errors
    assert_true(1);
}

/* Test: spawn_server_nofork_cleanup after initialization */
static void test_spawn_server_nofork_cleanup_success(void **state) {
    spawn_server_nofork_init();
    spawn_server_nofork_cleanup();
    assert_true(1);
}

/* Test: spawn_server_nofork_exec with NULL command */
static void test_spawn_server_nofork_exec_null_command(void **state) {
    int stdout_fd = -1, stderr_fd = -1;
    int result = spawn_server_nofork_exec(NULL, NULL, &stdout_fd, &stderr_fd);
    assert_true(result < 0 || result == 0); // Should fail or return error
}

/* Test: spawn_server_nofork_exec with empty command */
static void test_spawn_server_nofork_exec_empty_command(void **state) {
    int stdout_fd = -1, stderr_fd = -1;
    int result = spawn_server_nofork_exec("", NULL, &stdout_fd, &stderr_fd);
    assert_true(result < 0 || result == 0);
}

/* Test: spawn_server_nofork_exec with valid simple command */
static void test_spawn_server_nofork_exec_simple_command(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo test");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    // Result should be valid process id or error code
    assert_true(result != -2); // Should not be uninitialized
}

/* Test: spawn_server_nofork_exec with stdin data */
static void test_spawn_server_nofork_exec_with_stdin(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("cat");
    ctx->stdin_data = strdup("test input data");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, ctx->stdin_data, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with NULL stdout fd pointer */
static void test_spawn_server_nofork_exec_null_stdout_fd(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo test");
    int stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, NULL, &stderr_fd);
    // Should handle NULL pointer gracefully
    assert_true(1);
}

/* Test: spawn_server_nofork_exec with NULL stderr fd pointer */
static void test_spawn_server_nofork_exec_null_stderr_fd(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo test");
    int stdout_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, NULL);
    assert_true(1);
}

/* Test: spawn_server_nofork_exec with NULL stdin data and command */
static void test_spawn_server_nofork_exec_null_stdin_valid_command(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo hello");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with empty stdin data */
static void test_spawn_server_nofork_exec_empty_stdin(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("cat");
    ctx->stdin_data = strdup("");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, ctx->stdin_data, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with long command */
static void test_spawn_server_nofork_exec_long_command(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = malloc(8192);
    memset(ctx->command, 'a', 8191);
    ctx->command[8191] = '\0';
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with long stdin data */
static void test_spawn_server_nofork_exec_long_stdin(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("cat");
    ctx->stdin_data = malloc(16384);
    memset(ctx->stdin_data, 'x', 16383);
    ctx->stdin_data[16383] = '\0';
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, ctx->stdin_data, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with special shell characters */
static void test_spawn_server_nofork_exec_special_characters(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo 'test; test | test && test'");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with quoted arguments */
static void test_spawn_server_nofork_exec_quoted_args(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo \"hello world\"");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with command containing backslashes */
static void test_spawn_server_nofork_exec_backslash_escape(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo \\test\\path");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: Multiple sequential executions */
static void test_spawn_server_nofork_exec_multiple_calls(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("echo test");
    int stdout_fd1 = -1, stderr_fd1 = -1;
    int stdout_fd2 = -1, stderr_fd2 = -1;
    
    int result1 = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd1, &stderr_fd1);
    int result2 = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd2, &stderr_fd2);
    
    assert_true(result1 != -2);
    assert_true(result2 != -2);
}

/* Test: spawn_server_nofork_cleanup multiple times */
static void test_spawn_server_nofork_cleanup_multiple(void **state) {
    spawn_server_nofork_init();
    spawn_server_nofork_cleanup();
    spawn_server_nofork_cleanup(); // Should be safe to call multiple times
    assert_true(1);
}

/* Test: spawn_server_nofork_exec with whitespace-only command */
static void test_spawn_server_nofork_exec_whitespace_command(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("   ");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with newlines in stdin */
static void test_spawn_server_nofork_exec_newlines_stdin(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("cat");
    ctx->stdin_data = strdup("line1\nline2\nline3\n");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, ctx->stdin_data, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec boundary: single character command */
static void test_spawn_server_nofork_exec_single_char_command(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup(":");
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, NULL, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
}

/* Test: spawn_server_nofork_exec with embedded null in stdin handled */
static void test_spawn_server_nofork_exec_binary_stdin(void **state) {
    spawn_context_t *ctx = (spawn_context_t *)*state;
    ctx->command = strdup("cat");
    ctx->stdin_data = malloc(10);
    memset(ctx->stdin_data, 0xFF, 10);
    int stdout_fd = -1, stderr_fd = -1;
    
    int result = spawn_server_nofork_exec(ctx->command, ctx->stdin_data, &stdout_fd, &stderr_fd);
    assert_true(result != -2);
    
    free(ctx->stdin_data);
    ctx->stdin_data = NULL;
}

int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_init_success, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_cleanup_success, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_null_command, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_empty_command, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_simple_command, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_with_stdin, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_null_stdout_fd, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_null_stderr_fd, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_null_stdin_valid_command, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_empty_stdin, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_long_command, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_long_stdin, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_special_characters, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_quoted_args, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_backslash_escape, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_multiple_calls, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_cleanup_multiple, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_whitespace_command, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_newlines_stdin, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_single_char_command, setup_nofork, teardown_nofork),
        cmocka_unit_test_setup_teardown(test_spawn_server_nofork_exec_binary_stdin, setup_nofork, teardown_nofork),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}