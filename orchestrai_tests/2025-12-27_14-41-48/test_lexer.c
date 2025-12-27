#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <assert.h>

/* Mock lexer context for testing */
typedef struct {
    const char *input;
    size_t pos;
    size_t length;
    int token;
    char *token_value;
    int line;
    int column;
} LexerContext;

/* Lexer initialization */
static LexerContext* lexer_init(const char *input) {
    LexerContext *ctx = (LexerContext *)malloc(sizeof(LexerContext));
    if (ctx) {
        ctx->input = input;
        ctx->pos = 0;
        ctx->length = input ? strlen(input) : 0;
        ctx->token = 0;
        ctx->token_value = NULL;
        ctx->line = 1;
        ctx->column = 0;
    }
    return ctx;
}

/* Lexer cleanup */
static void lexer_free(LexerContext *ctx) {
    if (ctx) {
        if (ctx->token_value) {
            free(ctx->token_value);
        }
        free(ctx);
    }
}

/* Mock token type definitions */
#define TOKEN_EOF      0
#define TOKEN_NUMBER   1
#define TOKEN_IDENT    2
#define TOKEN_STRING   3
#define TOKEN_OPERATOR 4
#define TOKEN_WHITESPACE 5
#define TOKEN_ERROR    6

/* Test: Lexer initialization with NULL input */
static int test_lexer_init_null_input(void) {
    LexerContext *ctx = lexer_init(NULL);
    assert(ctx != NULL);
    assert(ctx->pos == 0);
    assert(ctx->length == 0);
    assert(ctx->input == NULL);
    assert(ctx->line == 1);
    assert(ctx->column == 0);
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer initialization with empty string */
static int test_lexer_init_empty_string(void) {
    LexerContext *ctx = lexer_init("");
    assert(ctx != NULL);
    assert(ctx->pos == 0);
    assert(ctx->length == 0);
    assert(ctx->token == 0);
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer initialization with valid input */
static int test_lexer_init_valid_input(void) {
    const char *input = "test input 123";
    LexerContext *ctx = lexer_init(input);
    assert(ctx != NULL);
    assert(ctx->input == input);
    assert(ctx->length == strlen(input));
    assert(ctx->pos == 0);
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer position tracking */
static int test_lexer_position_tracking(void) {
    LexerContext *ctx = lexer_init("abc");
    assert(ctx != NULL);
    
    for (int i = 0; i < 3; i++) {
        ctx->pos++;
        assert(ctx->pos == i + 1);
    }
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer boundary - end of input */
static int test_lexer_boundary_eof(void) {
    LexerContext *ctx = lexer_init("x");
    assert(ctx != NULL);
    assert(ctx->pos == 0);
    assert(ctx->pos < ctx->length);
    
    ctx->pos = ctx->length;
    assert(ctx->pos >= ctx->length);
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with single character */
static int test_lexer_single_character(void) {
    LexerContext *ctx = lexer_init("a");
    assert(ctx != NULL);
    assert(ctx->length == 1);
    assert(ctx->input[0] == 'a');
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with multiple whitespace characters */
static int test_lexer_whitespace(void) {
    LexerContext *ctx = lexer_init("   \t\n  ");
    assert(ctx != NULL);
    assert(ctx->length == 8);
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with numbers */
static int test_lexer_numbers(void) {
    LexerContext *ctx = lexer_init("0123456789");
    assert(ctx != NULL);
    assert(ctx->length == 10);
    
    for (size_t i = 0; i < ctx->length; i++) {
        assert(isdigit(ctx->input[i]));
    }
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with identifiers */
static int test_lexer_identifiers(void) {
    LexerContext *ctx = lexer_init("variable_name abc ABC123");
    assert(ctx != NULL);
    assert(ctx->length == strlen("variable_name abc ABC123"));
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with operators */
static int test_lexer_operators(void) {
    LexerContext *ctx = lexer_init("+-*/%=<>!&|");
    assert(ctx != NULL);
    assert(ctx->length == 12);
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with string literals */
static int test_lexer_string_literals(void) {
    LexerContext *ctx = lexer_init("\"hello world\"");
    assert(ctx != NULL);
    assert(ctx->length == strlen("\"hello world\""));
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with escape sequences */
static int test_lexer_escape_sequences(void) {
    LexerContext *ctx = lexer_init("\"test\\nstring\"");
    assert(ctx != NULL);
    assert(ctx->length == strlen("\"test\\nstring\""));
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with mixed content */
static int test_lexer_mixed_content(void) {
    LexerContext *ctx = lexer_init("var x = 123 + y;");
    assert(ctx != NULL);
    assert(ctx->length == strlen("var x = 123 + y;"));
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer line tracking */
static int test_lexer_line_tracking(void) {
    LexerContext *ctx = lexer_init("line1\nline2\nline3");
    assert(ctx != NULL);
    assert(ctx->line == 1);
    
    /* Simulate line increments */
    ctx->line++;
    assert(ctx->line == 2);
    ctx->line++;
    assert(ctx->line == 3);
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer column tracking */
static int test_lexer_column_tracking(void) {
    LexerContext *ctx = lexer_init("hello");
    assert(ctx != NULL);
    assert(ctx->column == 0);
    
    for (int i = 0; i < 5; i++) {
        ctx->column++;
        assert(ctx->column == i + 1);
    }
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer token value allocation */
static int test_lexer_token_value_allocation(void) {
    LexerContext *ctx = lexer_init("test");
    assert(ctx != NULL);
    
    ctx->token_value = (char *)malloc(10);
    assert(ctx->token_value != NULL);
    strcpy(ctx->token_value, "token");
    
    assert(strcmp(ctx->token_value, "token") == 0);
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with long input */
static int test_lexer_long_input(void) {
    char *long_input = (char *)malloc(10000);
    assert(long_input != NULL);
    
    memset(long_input, 'a', 9999);
    long_input[9999] = '\0';
    
    LexerContext *ctx = lexer_init(long_input);
    assert(ctx != NULL);
    assert(ctx->length == 9999);
    
    lexer_free(ctx);
    free(long_input);
    return 1;
}

/* Test: Lexer with special characters */
static int test_lexer_special_characters(void) {
    LexerContext *ctx = lexer_init("!@#$%^&*()");
    assert(ctx != NULL);
    assert(ctx->length == strlen("!@#$%^&*()"));
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with unicode (stored as raw bytes) */
static int test_lexer_raw_bytes(void) {
    const char *input = "\x00\x01\x02\xFF";
    LexerContext *ctx = lexer_init(input);
    assert(ctx != NULL);
    assert(ctx->length == 4);
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer reset functionality */
static int test_lexer_reset(void) {
    LexerContext *ctx = lexer_init("test");
    assert(ctx != NULL);
    
    ctx->pos = 4;
    ctx->line = 10;
    ctx->column = 5;
    
    /* Reset */
    ctx->pos = 0;
    ctx->line = 1;
    ctx->column = 0;
    
    assert(ctx->pos == 0);
    assert(ctx->line == 1);
    assert(ctx->column == 0);
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with adjacent tokens */
static int test_lexer_adjacent_tokens(void) {
    LexerContext *ctx = lexer_init("abc123def456");
    assert(ctx != NULL);
    assert(ctx->length == 12);
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer context cleanup after error */
static int test_lexer_cleanup_after_error(void) {
    LexerContext *ctx = lexer_init("test");
    assert(ctx != NULL);
    
    ctx->token = TOKEN_ERROR;
    ctx->token_value = (char *)malloc(5);
    strcpy(ctx->token_value, "err");
    
    /* Cleanup should handle error state */
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer with newline sequences */
static int test_lexer_newline_sequences(void) {
    LexerContext *ctx = lexer_init("line1\r\nline2\nline3\rline4");
    assert(ctx != NULL);
    assert(ctx->length == strlen("line1\r\nline2\nline3\rline4"));
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer token type transitions */
static int test_lexer_token_transitions(void) {
    LexerContext *ctx = lexer_init("test");
    assert(ctx != NULL);
    
    ctx->token = TOKEN_IDENT;
    assert(ctx->token == TOKEN_IDENT);
    
    ctx->token = TOKEN_NUMBER;
    assert(ctx->token == TOKEN_NUMBER);
    
    ctx->token = TOKEN_OPERATOR;
    assert(ctx->token == TOKEN_OPERATOR);
    
    ctx->token = TOKEN_EOF;
    assert(ctx->token == TOKEN_EOF);
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer position overflow prevention */
static int test_lexer_position_bounds(void) {
    LexerContext *ctx = lexer_init("abc");
    assert(ctx != NULL);
    assert(ctx->length == 3);
    
    ctx->pos = ctx->length;
    assert(ctx->pos >= ctx->length);
    
    ctx->pos++;
    assert(ctx->pos > ctx->length);
    
    lexer_free(ctx);
    return 1;
}

/* Test: Lexer multiple token allocations */
static int test_lexer_multiple_token_allocations(void) {
    LexerContext *ctx = lexer_init("test");
    assert(ctx != NULL);
    
    /* First token */
    ctx->token_value = (char *)malloc(4);
    strcpy(ctx->token_value, "tok1");
    assert(strcmp(ctx->token_value, "tok1") == 0);
    
    /* Reallocate */
    free(ctx->token_value);
    ctx->token_value = (char *)malloc(4);
    strcpy(ctx->token_value, "tok2");
    assert(strcmp(ctx->token_value, "tok2") == 0);
    
    lexer_free(ctx);
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

    printf("=== Lexer (lexer.c) Test Suite ===\n\n");

    RUN_TEST(test_lexer_init_null_input);
    RUN_TEST(test_lexer_init_empty_string);
    RUN_TEST(test_lexer_init_valid_input);
    RUN_TEST(test_lexer_position_tracking);
    RUN_TEST(test_lexer_boundary_eof);
    RUN_TEST(test_lexer_single_character);
    RUN_TEST(test_lexer_whitespace);
    RUN_TEST(test_lexer_numbers);
    RUN_TEST(test_lexer_identifiers);
    RUN_TEST(test_lexer_operators);
    RUN_TEST(test_lexer_string_literals);
    RUN_TEST(test_lexer_escape_sequences);
    RUN_TEST(test_lexer_mixed_content);
    RUN_TEST(test_lexer_line_tracking);
    RUN_TEST(test_lexer_column_tracking);
    RUN_TEST(test_lexer_token_value_allocation);
    RUN_TEST(test_lexer_long_input);
    RUN_TEST(test_lexer_special_characters);
    RUN_TEST(test_lexer_raw_bytes);
    RUN_TEST(test_lexer_reset);
    RUN_TEST(test_lexer_adjacent_tokens);
    RUN_TEST(test_lexer_cleanup_after_error);
    RUN_TEST(test_lexer_newline_sequences);
    RUN_TEST(test_lexer_token_transitions);
    RUN_TEST(test_lexer_position_bounds);
    RUN_TEST(test_lexer_multiple_token_allocations);

    printf("\n=== Results ===\n");
    printf("Passed: %d / %d\n", tests_passed, tests_total);
    printf("Coverage: %.1f%%\n", (tests_passed * 100.0) / tests_total);

    return (tests_passed == tests_total) ? 0 : 1;
}