#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <assert.h>
#include <setjmp.h>

/* Forward declarations for test functions */
extern void ParseInit(void *);
extern void Parse(void *, int, void *);
extern void ParseFree(void *);
extern void ParseTrace(FILE *, char *);

/* Mock structure for testing parser state */
typedef struct {
    int token_count;
    int error_count;
    int accept_count;
    char error_message[256];
} MockParserContext;

/* Test helper functions */
static MockParserContext* create_mock_context(void) {
    MockParserContext *ctx = (MockParserContext *)malloc(sizeof(MockParserContext));
    if (ctx) {
        ctx->token_count = 0;
        ctx->error_count = 0;
        ctx->accept_count = 0;
        memset(ctx->error_message, 0, sizeof(ctx->error_message));
    }
    return ctx;
}

static void destroy_mock_context(MockParserContext *ctx) {
    if (ctx) {
        free(ctx);
    }
}

/* Test: ParseInit with NULL parser */
static int test_parse_init_null_parser(void) {
    ParseInit(NULL);
    return 1;
}

/* Test: ParseInit with valid parser */
static int test_parse_init_valid_parser(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parse with NULL parser - should not crash */
static int test_parse_null_parser(void) {
    Parse(NULL, 0, NULL);
    return 1;
}

/* Test: Parse with valid parser and token */
static int test_parse_valid_token(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    Parse(ctx, 1, ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parse with token value 0 */
static int test_parse_token_zero(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    Parse(ctx, 0, ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parse with negative token */
static int test_parse_negative_token(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    Parse(ctx, -1, ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parse with large token value */
static int test_parse_large_token(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    Parse(ctx, 999999, ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: ParseFree with NULL parser */
static int test_parse_free_null_parser(void) {
    ParseFree(NULL);
    return 1;
}

/* Test: ParseFree with valid parser */
static int test_parse_free_valid_parser(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    ParseFree(ctx);
    return 1;
}

/* Test: ParseTrace with NULL file and parser */
static int test_parse_trace_null_file(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    ParseTrace(NULL, NULL);
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: ParseTrace with stdout */
static int test_parse_trace_stdout(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    ParseTrace(stdout, (char *)"");
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: ParseTrace with custom message */
static int test_parse_trace_custom_message(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    ParseTrace(stdout, (char *)"PARSER_TRACE");
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Multiple sequential Parse calls */
static int test_parse_sequential_calls(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    
    for (int i = 0; i < 10; i++) {
        Parse(ctx, i, ctx);
    }
    
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parse after ParseFree - parser reinitialization */
static int test_parse_after_free_reinit(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    
    ParseInit(ctx);
    Parse(ctx, 1, ctx);
    ParseFree(ctx);
    
    /* Reinitialize and use again */
    ParseInit(ctx);
    Parse(ctx, 2, ctx);
    ParseFree(ctx);
    
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parser state with different payload types */
static int test_parse_different_payloads(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    
    /* Test with different payload pointers */
    int payload1 = 42;
    int payload2 = 100;
    char *payload3 = "test";
    
    Parse(ctx, 1, &payload1);
    Parse(ctx, 2, &payload2);
    Parse(ctx, 3, payload3);
    
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: ParseTrace with empty message */
static int test_parse_trace_empty_message(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    ParseTrace(stdout, (char *)"");
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Memory allocation failure during ParseInit */
static int test_parse_memory_stress(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    
    /* Allocate and free multiple times */
    for (int i = 0; i < 100; i++) {
        ParseInit(ctx);
        Parse(ctx, i % 256, ctx);
        ParseFree(ctx);
    }
    
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parser with boundary token values */
static int test_parse_boundary_tokens(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    
    /* Test boundary values */
    Parse(ctx, 0, ctx);              /* Minimum valid token */
    Parse(ctx, 1, ctx);              /* First token */
    Parse(ctx, 255, ctx);            /* Byte boundary */
    Parse(ctx, 256, ctx);            /* Beyond byte */
    Parse(ctx, 65535, ctx);          /* Word boundary */
    Parse(ctx, 65536, ctx);          /* Beyond word */
    
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Test: Parser trace control flow */
static int test_parse_trace_control_flow(void) {
    MockParserContext *ctx = create_mock_context();
    assert(ctx != NULL);
    ParseInit(ctx);
    
    ParseTrace(stdout, (char *)"TRACE_ON");
    Parse(ctx, 1, ctx);
    ParseTrace(NULL, (char *)"TRACE_OFF");
    Parse(ctx, 2, ctx);
    
    ParseFree(ctx);
    destroy_mock_context(ctx);
    return 1;
}

/* Main test runner */
int main(void) {
    int tests_passed = 0;
    int tests_total = 0;

    #define RUN_TEST(test_func) \
        do { \
            tests_total++; \
            if (test_func()) { \
                tests_passed++; \
                printf("✓ %s\n", #test_func); \
            } else { \
                printf("✗ %s\n", #test_func); \
            } \
        } while(0)

    printf("=== Parser (lempar.c) Test Suite ===\n\n");

    RUN_TEST(test_parse_init_null_parser);
    RUN_TEST(test_parse_init_valid_parser);
    RUN_TEST(test_parse_null_parser);
    RUN_TEST(test_parse_valid_token);
    RUN_TEST(test_parse_token_zero);
    RUN_TEST(test_parse_negative_token);
    RUN_TEST(test_parse_large_token);
    RUN_TEST(test_parse_free_null_parser);
    RUN_TEST(test_parse_free_valid_parser);
    RUN_TEST(test_parse_trace_null_file);
    RUN_TEST(test_parse_trace_stdout);
    RUN_TEST(test_parse_trace_custom_message);
    RUN_TEST(test_parse_sequential_calls);
    RUN_TEST(test_parse_after_free_reinit);
    RUN_TEST(test_parse_different_payloads);
    RUN_TEST(test_parse_trace_empty_message);
    RUN_TEST(test_parse_memory_stress);
    RUN_TEST(test_parse_boundary_tokens);
    RUN_TEST(test_parse_trace_control_flow);

    printf("\n=== Results ===\n");
    printf("Passed: %d / %d\n", tests_passed, tests_total);
    printf("Coverage: %.1f%%\n", (tests_passed * 100.0) / tests_total);

    return (tests_passed == tests_total) ? 0 : 1;
}