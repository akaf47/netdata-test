#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdarg.h>
#include <stddef.h>

/* Forward declarations for functions being tested */
extern void *eval_parse_expression(const char *expression);
extern void eval_free_expression(void *expression);
extern int eval_evaluate_expression(void *expression);

/* Test setup and teardown */
static int setup(void **state) {
    return 0;
}

static int teardown(void **state) {
    return 0;
}

/* ============= NULL/EMPTY INPUT TESTS ============= */
static void test_parse_expression_with_null_input(void **state) {
    void *result = eval_parse_expression(NULL);
    assert_null(result);
}

static void test_parse_expression_with_empty_string(void **state) {
    void *result = eval_parse_expression("");
    assert_null(result);
}

static void test_parse_expression_with_whitespace_only(void **state) {
    void *result = eval_parse_expression("   ");
    assert_null(result);
}

/* ============= VALID EXPRESSION PARSING TESTS ============= */
static void test_parse_simple_number(void **state) {
    void *expr = eval_parse_expression("42");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_decimal_number(void **state) {
    void *expr = eval_parse_expression("3.14");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_negative_number(void **state) {
    void *expr = eval_parse_expression("-42");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_variable_reference(void **state) {
    void *expr = eval_parse_expression("myvar");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= ARITHMETIC OPERATIONS TESTS ============= */
static void test_parse_addition(void **state) {
    void *expr = eval_parse_expression("2 + 3");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_subtraction(void **state) {
    void *expr = eval_parse_expression("5 - 2");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_multiplication(void **state) {
    void *expr = eval_parse_expression("4 * 3");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_division(void **state) {
    void *expr = eval_parse_expression("10 / 2");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_modulo(void **state) {
    void *expr = eval_parse_expression("10 % 3");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= OPERATOR PRECEDENCE TESTS ============= */
static void test_parse_precedence_mult_before_add(void **state) {
    void *expr = eval_parse_expression("2 + 3 * 4");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_precedence_division_before_subtraction(void **state) {
    void *expr = eval_parse_expression("10 - 8 / 2");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_mixed_operators(void **state) {
    void *expr = eval_parse_expression("2 + 3 * 4 - 5 / 2");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= PARENTHESES TESTS ============= */
static void test_parse_parentheses_simple(void **state) {
    void *expr = eval_parse_expression("(2 + 3)");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_parentheses_nested(void **state) {
    void *expr = eval_parse_expression("((2 + 3) * (4 - 1))");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_parentheses_override_precedence(void **state) {
    void *expr = eval_parse_expression("(2 + 3) * 4");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= COMPARISON OPERATORS TESTS ============= */
static void test_parse_equals(void **state) {
    void *expr = eval_parse_expression("5 == 5");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_not_equals(void **state) {
    void *expr = eval_parse_expression("5 != 3");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_less_than(void **state) {
    void *expr = eval_parse_expression("3 < 5");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_greater_than(void **state) {
    void *expr = eval_parse_expression("5 > 3");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_less_than_or_equal(void **state) {
    void *expr = eval_parse_expression("5 <= 5");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_greater_than_or_equal(void **state) {
    void *expr = eval_parse_expression("5 >= 3");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= LOGICAL OPERATORS TESTS ============= */
static void test_parse_logical_and(void **state) {
    void *expr = eval_parse_expression("1 && 1");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_logical_or(void **state) {
    void *expr = eval_parse_expression("0 || 1");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_logical_not(void **state) {
    void *expr = eval_parse_expression("!0");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_complex_boolean_expression(void **state) {
    void *expr = eval_parse_expression("(5 > 3) && (2 < 4)");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= TERNARY OPERATOR TESTS ============= */
static void test_parse_ternary_operator(void **state) {
    void *expr = eval_parse_expression("5 > 3 ? 1 : 0");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_nested_ternary(void **state) {
    void *expr = eval_parse_expression("5 > 3 ? (2 < 4 ? 1 : 0) : -1");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= UNARY OPERATORS TESTS ============= */
static void test_parse_unary_minus(void **state) {
    void *expr = eval_parse_expression("-5");
    assert_non_null(expr);
    eval_free_expression(expr);
}

static void test_parse_unary_plus(void **state) {
    void *expr = eval_parse_expression("+5");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= MALFORMED EXPRESSION TESTS ============= */
static void test_parse_unmatched_left_paren(void **state) {
    void *expr = eval_parse_expression("(2 + 3");
    assert_null(expr);
}

static void test_parse_unmatched_right_paren(void **state) {
    void *expr = eval_parse_expression("2 + 3)");
    assert_null(expr);
}

static void test_parse_missing_operand(void **state) {
    void *expr = eval_parse_expression("2 +");
    assert_null(expr);
}

static void test_parse_consecutive_operators(void **state) {
    void *expr = eval_parse_expression("2 + * 3");
    assert_null(expr);
}

static void test_parse_invalid_characters(void **state) {
    void *expr = eval_parse_expression("2 + @ 3");
    assert_null(expr);
}

static void test_parse_division_by_zero_expression(void **state) {
    void *expr = eval_parse_expression("1 / 0");
    assert_non_null(expr);
    eval_free_expression(expr);
}

/* ============= EVALUATION TESTS ============= */
static void test_evaluate_simple_number(void **state) {
    void *expr = eval_parse_expression("42");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 42);
    eval_free_expression(expr);
}

static void test_evaluate_addition(void **state) {
    void *expr = eval_parse_expression("2 + 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

static void test_evaluate_subtraction(void **state) {
    void *expr = eval_parse_expression("5 - 2");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 3);
    eval_free_expression(expr);
}

static void test_evaluate_multiplication(void **state) {
    void *expr = eval_parse_expression("4 * 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 12);
    eval_free_expression(expr);
}

static void test_evaluate_division(void **state) {
    void *expr = eval_parse_expression("10 / 2");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

static void test_evaluate_modulo(void **state) {
    void *expr = eval_parse_expression("10 % 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_precedence_mult_before_add(void **state) {
    void *expr = eval_parse_expression("2 + 3 * 4");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 14);
    eval_free_expression(expr);
}

static void test_evaluate_parentheses_simple(void **state) {
    void *expr = eval_parse_expression("(2 + 3) * 4");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 20);
    eval_free_expression(expr);
}

static void test_evaluate_equals_true(void **state) {
    void *expr = eval_parse_expression("5 == 5");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_equals_false(void **state) {
    void *expr = eval_parse_expression("5 == 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_not_equals_true(void **state) {
    void *expr = eval_parse_expression("5 != 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_not_equals_false(void **state) {
    void *expr = eval_parse_expression("5 != 5");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_less_than_true(void **state) {
    void *expr = eval_parse_expression("3 < 5");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_less_than_false(void **state) {
    void *expr = eval_parse_expression("5 < 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_greater_than_true(void **state) {
    void *expr = eval_parse_expression("5 > 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_greater_than_false(void **state) {
    void *expr = eval_parse_expression("3 > 5");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_logical_and_true(void **state) {
    void *expr = eval_parse_expression("1 && 1");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_logical_and_false(void **state) {
    void *expr = eval_parse_expression("1 && 0");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_logical_or_true(void **state) {
    void *expr = eval_parse_expression("0 || 1");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_logical_or_false(void **state) {
    void *expr = eval_parse_expression("0 || 0");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_logical_not_true(void **state) {
    void *expr = eval_parse_expression("!0");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 1);
    eval_free_expression(expr);
}

static void test_evaluate_logical_not_false(void **state) {
    void *expr = eval_parse_expression("!1");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_ternary_true_branch(void **state) {
    void *expr = eval_parse_expression("1 ? 5 : 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

static void test_evaluate_ternary_false_branch(void **state) {
    void *expr = eval_parse_expression("0 ? 5 : 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 3);
    eval_free_expression(expr);
}

static void test_evaluate_unary_minus(void **state) {
    void *expr = eval_parse_expression("-5");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, -5);
    eval_free_expression(expr);
}

static void test_evaluate_unary_plus(void **state) {
    void *expr = eval_parse_expression("+5");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

/* ============= EDGE CASES AND BOUNDARY VALUES ============= */
static void test_evaluate_zero(void **state) {
    void *expr = eval_parse_expression("0");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 0);
    eval_free_expression(expr);
}

static void test_evaluate_large_number(void **state) {
    void *expr = eval_parse_expression("999999");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 999999);
    eval_free_expression(expr);
}

static void test_evaluate_negative_number(void **state) {
    void *expr = eval_parse_expression("-999");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, -999);
    eval_free_expression(expr);
}

static void test_free_null_expression(void **state) {
    eval_free_expression(NULL);
}

static void test_evaluate_complex_nested_expression(void **state) {
    void *expr = eval_parse_expression("((5 + 3) * 2 - 4) / 2");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 6);
    eval_free_expression(expr);
}

static void test_evaluate_all_operators_combined(void **state) {
    void *expr = eval_parse_expression("((5 > 3 && 2 < 4) ? 10 : 5) + 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 13);
    eval_free_expression(expr);
}

/* ============= WHITESPACE HANDLING TESTS ============= */
static void test_parse_expression_with_extra_spaces(void **state) {
    void *expr = eval_parse_expression("2   +   3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

static void test_parse_expression_with_tabs(void **state) {
    void *expr = eval_parse_expression("2\t+\t3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

static void test_parse_expression_with_leading_whitespace(void **state) {
    void *expr = eval_parse_expression("  2 + 3");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

static void test_parse_expression_with_trailing_whitespace(void **state) {
    void *expr = eval_parse_expression("2 + 3  ");
    assert_non_null(expr);
    int result = eval_evaluate_expression(expr);
    assert_int_equal(result, 5);
    eval_free_expression(expr);
}

/* ============= TEST SUITE CONFIGURATION ============= */
int main(void) {
    const struct CMUnitTest tests[] = {
        /* NULL/Empty input tests */
        cmocka_unit_test(test_parse_expression_with_null_input),
        cmocka_unit_test(test_parse_expression_with_empty_string),
        cmocka_unit_test(test_parse_expression_with_whitespace_only),
        
        /* Valid expression parsing tests */
        cmocka_unit_test(test_parse_simple_number),
        cmocka_unit_test(test_parse_decimal_number),
        cmocka_unit_test(test_parse_negative_number),
        cmocka_unit_test(test_parse_variable_reference),
        
        /* Arithmetic operations tests */
        cmocka_unit_test(test_parse_addition),
        cmocka_unit_test(test_parse_subtraction),
        cmocka_unit_test(test_parse_multiplication),
        cmocka_unit_test(test_parse_division),
        cmocka_unit_test(test_parse_modulo),
        
        /* Operator precedence tests */
        cmocka_unit_test(test_parse_precedence_mult_before_add),
        cmocka_unit_test(test_parse_precedence_division_before_subtraction),
        cmocka_unit_test(test_parse_mixed_operators),
        
        /* Parentheses tests */
        cmocka_unit_test(test_parse_parentheses_simple),
        cmocka_unit_test(test_parse_parentheses_nested),
        cmocka_unit_test(test_parse_parentheses_override_precedence),
        
        /* Comparison operators tests */
        cmocka_unit_test(test_parse_equals),
        cmocka_unit_test(test_parse_not_equals),
        cmocka_unit_test(test_parse_less_than),
        cmocka_unit_test(test_parse_greater_than),
        cmocka_unit_test(test_parse_less_than_or_equal),
        cmocka_unit_test(test_parse_greater_than_or_equal),
        
        /* Logical operators tests */
        cmocka_unit_test(test_parse_logical_and),
        cmocka_unit_test(test_parse_logical_or),
        cmocka_unit_test(test_parse_logical_not),
        cmocka_unit_test(test_parse_complex_boolean_expression),
        
        /* Ternary operator tests */
        cmocka_unit_test(test_parse_ternary_operator),
        cmocka_unit_test(test_parse_nested_ternary),
        
        /* Unary operators tests */
        cmocka_unit_test(test_parse_unary_minus),
        cmocka_unit_test(test_parse_unary_plus),
        
        /* Malformed expression tests */
        cmocka_unit_test(test_parse_unmatched_left_paren),
        cmocka_unit_test(test_parse_unmatched_right_paren),
        cmocka_unit_test(test_parse_missing_operand),
        cmocka_unit_test(test_parse_consecutive_operators),
        cmocka_unit_test(test_parse_invalid_characters),
        cmocka_unit_test(test_parse_division_by_zero_expression),
        
        /* Evaluation tests */
        cmocka_unit_test(test_evaluate_simple_number),
        cmocka_unit_test(test_evaluate_addition),
        cmocka_unit_test(test_evaluate_subtraction),
        cmocka_unit_test(test_evaluate_multiplication),
        cmocka_unit_test(test_evaluate_division),
        cmocka_unit_test(test_evaluate_modulo),
        cmocka_unit_test(test_evaluate_precedence_mult_before_add),
        cmocka_unit_test(test_evaluate_parentheses_simple),
        cmocka_unit_test(test_evaluate_equals_true),
        cmocka_unit_test(test_evaluate_equals_false),
        cmocka_unit_test(test_evaluate_not_equals_true),
        cmocka_unit_test(test_evaluate_not_equals_false),
        cmocka_unit_test(test_evaluate_less_than_true),
        cmocka_unit_test(test_evaluate_less_than_false),
        cmocka_unit_test(test_evaluate_greater_than_true),
        cmocka_unit_test(test_evaluate_greater_than_false),
        cmocka_unit_test(test_evaluate_logical_and_true),
        cmocka_unit_test(test_evaluate_logical_and_false),
        cmocka_unit_test(test_evaluate_logical_or_true),
        cmocka_unit_test(test_evaluate_logical_or_false),
        cmocka_unit_test(test_evaluate_logical_not_true),
        cmocka_unit_test(test_evaluate_logical_not_false),
        cmocka_unit_test(test_evaluate_ternary_true_branch),
        cmocka_unit_test(test_evaluate_ternary_false_branch),
        cmocka_unit_test(test_evaluate_unary_minus),
        cmocka_unit_test(test_evaluate_unary_plus),
        
        /* Edge cases and boundary value tests */
        cmocka_unit_test(test_evaluate_zero),
        cmocka_unit_test(test_evaluate_large_number),
        cmocka_unit_test(test_evaluate_negative_number),
        cmocka_unit_test(test_free_null_expression),
        cmocka_unit_test(test_evaluate_complex_nested_expression),
        cmocka_unit_test(test_evaluate_all_operators_combined),
        
        /* Whitespace handling tests */
        cmocka_unit_test(test_parse_expression_with_extra_spaces),
        cmocka_unit_test(test_parse_expression_with_tabs),
        cmocka_unit_test(test_parse_expression_with_leading_whitespace),
        cmocka_unit_test(test_parse_expression_with_trailing_whitespace),
    };

    return cmocka_run_group_tests(tests, setup, teardown);
}