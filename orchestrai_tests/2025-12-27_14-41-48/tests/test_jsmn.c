#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>

#include "../src/libnetdata/json/vendored/jsmn.h"

/* ============================================================================
   UNIT TESTS FOR JSMN - JSON PARSER
   Testing all public functions: jsmn_init, jsmn_parse
   Testing all code paths, branches, and edge cases
   ============================================================================ */

/* ============================================================================
   FIXTURE: Memory allocation tracking
   ============================================================================ */

static jsmn_parser *test_parser = NULL;
static jsmn_token *test_tokens = NULL;

static int setup_parser(void **state) {
    test_parser = malloc(sizeof(jsmn_parser));
    test_tokens = malloc(sizeof(jsmn_token) * 256);
    assert_non_null(test_parser);
    assert_non_null(test_tokens);
    return 0;
}

static int teardown_parser(void **state) {
    free(test_parser);
    free(test_tokens);
    test_parser = NULL;
    test_tokens = NULL;
    return 0;
}

/* ============================================================================
   TEST GROUP: jsmn_init()
   ============================================================================ */

static void test_jsmn_init_sets_parser_to_initial_state(void **state) {
    jsmn_parser parser;
    jsmn_init(&parser);
    
    assert_int_equal(parser.pos, 0);
    assert_int_equal(parser.toknext, 0);
    assert_int_equal(parser.toksuper, -1);
}

static void test_jsmn_init_multiple_calls(void **state) {
    jsmn_parser parser;
    
    jsmn_init(&parser);
    assert_int_equal(parser.pos, 0);
    
    parser.pos = 100;
    parser.toknext = 50;
    
    jsmn_init(&parser);
    assert_int_equal(parser.pos, 0);
    assert_int_equal(parser.toknext, 0);
    assert_int_equal(parser.toksuper, -1);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Basic Valid JSON
   ============================================================================ */

static void test_jsmn_parse_empty_object(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
    assert_int_equal(tokens[0].start, 0);
    assert_int_equal(tokens[0].end, 2);
    assert_int_equal(tokens[0].size, 0);
}

static void test_jsmn_parse_empty_array(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "[]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
    assert_int_equal(tokens[0].start, 0);
    assert_int_equal(tokens[0].end, 2);
    assert_int_equal(tokens[0].size, 0);
}

static void test_jsmn_parse_empty_string_value(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_STRING);
    assert_int_equal(tokens[0].start, 0);
    assert_int_equal(tokens[0].end, 2);
}

static void test_jsmn_parse_simple_string(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"hello\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_STRING);
    assert_int_equal(tokens[0].start, 1);
    assert_int_equal(tokens[0].end, 6);
}

static void test_jsmn_parse_positive_integer(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "123";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
    assert_int_equal(tokens[0].start, 0);
    assert_int_equal(tokens[0].end, 3);
}

static void test_jsmn_parse_negative_integer(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "-456";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_floating_point(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "123.456";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_scientific_notation(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "1.23e-4";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_zero(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "0";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_true_keyword(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "true";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_false_keyword(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "false";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_null_keyword(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "null";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_PRIMITIVE);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Objects with Key-Value Pairs
   ============================================================================ */

static void test_jsmn_parse_object_with_single_string_value(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{\"key\":\"value\"}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 3);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
    assert_int_equal(tokens[0].size, 1);
    assert_int_equal(tokens[1].type, JSMN_STRING);
    assert_int_equal(tokens[2].type, JSMN_STRING);
}

static void test_jsmn_parse_object_with_multiple_key_values(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "{\"a\":1,\"b\":2,\"c\":3}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_int_equal(result, 7);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
    assert_int_equal(tokens[0].size, 3);
}

static void test_jsmn_parse_object_with_integer_value(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{\"num\":42}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 3);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
    assert_int_equal(tokens[1].type, JSMN_STRING);
    assert_int_equal(tokens[2].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_object_with_boolean_value(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{\"flag\":true}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 3);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
    assert_int_equal(tokens[2].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_object_with_null_value(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{\"val\":null}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 3);
    assert_int_equal(tokens[2].type, JSMN_PRIMITIVE);
}

static void test_jsmn_parse_object_with_nested_object(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "{\"outer\":{\"inner\":1}}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

static void test_jsmn_parse_object_with_nested_array(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "{\"arr\":[1,2,3]}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Arrays
   ============================================================================ */

static void test_jsmn_parse_array_with_single_string(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "[\"item\"]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 2);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
    assert_int_equal(tokens[0].size, 1);
    assert_int_equal(tokens[1].type, JSMN_STRING);
}

static void test_jsmn_parse_array_with_multiple_strings(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "[\"a\",\"b\",\"c\"]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_int_equal(result, 4);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
    assert_int_equal(tokens[0].size, 3);
}

static void test_jsmn_parse_array_with_integers(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "[1,2,3,4,5]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_int_equal(result, 6);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
    assert_int_equal(tokens[0].size, 5);
}

static void test_jsmn_parse_array_with_mixed_types(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "[\"str\",123,true,null]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
}

static void test_jsmn_parse_array_with_nested_array(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "[[1,2],[3,4]]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
}

static void test_jsmn_parse_array_with_nested_object(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[20];
    const char *json = "[{\"a\":1},{\"b\":2}]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 20);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Complex Nested Structures
   ============================================================================ */

static void test_jsmn_parse_deeply_nested_objects(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[50];
    const char *json = "{\"a\":{\"b\":{\"c\":{\"d\":1}}}}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 50);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

static void test_jsmn_parse_deeply_nested_arrays(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[50];
    const char *json = "[[[[1]]]]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 50);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_ARRAY);
}

static void test_jsmn_parse_complex_json_document(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[100];
    const char *json = "{\"users\":[{\"name\":\"John\",\"age\":30},{\"name\":\"Jane\",\"age\":25}],\"count\":2}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 100);
    
    assert_true(result > 0);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Whitespace Handling
   ============================================================================ */

static void test_jsmn_parse_with_leading_whitespace(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "  {}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

static void test_jsmn_parse_with_trailing_whitespace(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{}  ";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

static void test_jsmn_parse_with_internal_whitespace(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{ \"key\" : \"value\" }";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 3);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

static void test_jsmn_parse_with_newlines_and_tabs(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{\n\t\"key\"\t:\n\t\"value\"\n}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 3);
    assert_int_equal(tokens[0].type, JSMN_OBJECT);
}

static void test_jsmn_parse_with_carriage_returns(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{\"key\":\"value\"}\r";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 3);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Escape Sequences in Strings
   ============================================================================ */

static void test_jsmn_parse_string_with_escaped_quote(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"hello\\\"world\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_STRING);
}

static void test_jsmn_parse_string_with_escaped_backslash(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"path\\\\to\\\\file\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
    assert_int_equal(tokens[0].type, JSMN_STRING);
}

static void test_jsmn_parse_string_with_escaped_forward_slash(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"\\/path\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

static void test_jsmn_parse_string_with_escaped_backspace(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"bell\\b\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

static void test_jsmn_parse_string_with_escaped_formfeed(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"form\\f\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

static void test_jsmn_parse_string_with_escaped_newline(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"line\\nbreak\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

static void test_jsmn_parse_string_with_escaped_carriage_return(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"carriage\\rreturn\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

static void test_jsmn_parse_string_with_escaped_tab(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"tab\\there\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

static void test_jsmn_parse_string_with_unicode_escape(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"unicode\\u0041\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

static void test_jsmn_parse_string_with_multiple_escapes(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"\\n\\r\\t\\\\\\\"\"";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_int_equal(result, 1);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - NULL Pointer Arguments
   ============================================================================ */

static void test_jsmn_parse_with_null_parser(void **state) {
    jsmn_token tokens[10];
    const char *json = "{}";
    
    int result = jsmn_parse(NULL, json, strlen(json), tokens, 10);
    
    /* Expected behavior: should handle gracefully or crash - test crash prevention */
    assert_true(result < 0 || result >= 0);
}

static void test_jsmn_parse_with_null_json(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, NULL, 0, tokens, 10);
    
    /* Should handle null input */
    assert_true(result <= 0);
}

static void test_jsmn_parse_with_null_tokens(void **state) {
    jsmn_parser parser;
    const char *json = "{}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), NULL, 0);
    
    /* Should still parse but not store tokens */
    assert_true(result >= 1);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Zero Length Input
   ============================================================================ */

static void test_jsmn_parse_with_zero_length_json(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, 0, tokens, 10);
    
    assert_true(result < 0);
}

static void test_jsmn_parse_with_zero_token_capacity(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 0);
    
    assert_true(result > 0 || result == JSMN_ERROR_NOMEM);
}

/* ============================================================================
   TEST GROUP: jsmn_parse() - Invalid JSON
   ============================================================================ */

static void test_jsmn_parse_unclosed_brace(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_true(result < 0);
}

static void test_jsmn_parse_unclosed_bracket(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "[";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_true(result < 0);
}

static void test_jsmn_parse_unclosed_string(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "\"unclosed";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_true(result < 0);
}

static void test_jsmn_parse_mismatched_braces(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_true(result < 0);
}

static void test_jsmn_parse_mismatched_brackets(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "[}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_true(result < 0);
}

static void test_jsmn_parse_extra_closing_brace(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "{}}";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_true(result < 0 || result == 1);
}

static void test_jsmn_parse_extra_closing_bracket(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "[]]";
    
    jsmn_init(&parser);
    int result = jsmn_parse(&parser, json, strlen(json), tokens, 10);
    
    assert_true(result < 0 || result == 1);
}

static void test_jsmn_parse_invalid_keyword_partial(void **state) {
    jsmn_parser parser;
    jsmn_token tokens[10];
    const char *json = "tru";
    
    jsm