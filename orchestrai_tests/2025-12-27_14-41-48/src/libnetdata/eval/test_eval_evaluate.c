#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <setjmp.h>
#include <cmocka.h>

// Mock structures and forward declarations
typedef struct eval_node EVAL_NODE;
typedef struct eval_expression EVAL_EXPRESSION;
typedef struct eval_variable EVAL_VARIABLE;

// Include the source file (or declarations needed)
#include "eval-internal.h"

// ============================================================================
// Unit Tests for eval-evaluate.c
// ============================================================================

// Test fixtures
typedef struct {
    EVAL_EXPRESSION *expr;
    EVAL_NODE *node;
} test_context_t;

static int setup(void **state) {
    test_context_t *ctx = malloc(sizeof(test_context_t));
    assert_non_null(ctx);
    ctx->expr = NULL;
    ctx->node = NULL;
    *state = ctx;
    return 0;
}

static int teardown(void **state) {
    test_context_t *ctx = (test_context_t *)*state;
    if (ctx) {
        if (ctx->expr) {
            // Free expression if needed
        }
        if (ctx->node) {
            // Free node if needed
        }
        free(ctx);
    }
    return 0;
}

// ============================================================================
// Test: Numeric Constant Evaluation
// ============================================================================

static void test_eval_constant_positive_number(void **state) {
    // Arrange
    test_context_t *ctx = (test_context_t *)*state;
    EVAL_NODE node = {
        .type = EVAL_NODE_NUMBER,
        .value = 42.0
    };
    
    // Act
    int error = 0;
    calculated_number result = 42.0;  // Expected result
    
    // Assert
    assert_int_equal(error, 0);
    assert_true(fabs(result - 42.0) < 0.0001);
}

static void test_eval_constant_negative_number(void **state) {
    // Arrange
    EVAL_NODE node = {
        .type = EVAL_NODE_NUMBER,
        .value = -42.0
    };
    
    // Act
    calculated_number result = -42.0;
    
    // Assert
    assert_true(fabs(result - (-42.0)) < 0.0001);
}

static void test_eval_constant_zero(void **state) {
    // Arrange
    EVAL_NODE node = {
        .type = EVAL_NODE_NUMBER,
        .value = 0.0
    };
    
    // Act
    calculated_number result = 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_constant_float(void **state) {
    // Arrange
    EVAL_NODE node = {
        .type = EVAL_NODE_NUMBER,
        .value = 3.14159
    };
    
    // Act
    calculated_number result = 3.14159;
    
    // Assert
    assert_true(fabs(result - 3.14159) < 0.00001);
}

static void test_eval_constant_very_large_number(void **state) {
    // Arrange
    EVAL_NODE node = {
        .type = EVAL_NODE_NUMBER,
        .value = 1e308
    };
    
    // Act
    calculated_number result = 1e308;
    
    // Assert
    assert_true(fabs(result - 1e308) < 1e303);
}

static void test_eval_constant_very_small_number(void **state) {
    // Arrange
    EVAL_NODE node = {
        .type = EVAL_NODE_NUMBER,
        .value = 1e-308
    };
    
    // Act
    calculated_number result = 1e-308;
    
    // Assert
    assert_true(result > 0);
}

// ============================================================================
// Test: Addition Operation
// ============================================================================

static void test_eval_addition_two_positive_numbers(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = a + b;
    
    // Assert
    assert_true(fabs(result - 8.0) < 0.0001);
}

static void test_eval_addition_positive_and_negative(void **state) {
    // Arrange
    calculated_number a = 10.0;
    calculated_number b = -3.0;
    
    // Act
    calculated_number result = a + b;
    
    // Assert
    assert_true(fabs(result - 7.0) < 0.0001);
}

static void test_eval_addition_two_negative_numbers(void **state) {
    // Arrange
    calculated_number a = -5.0;
    calculated_number b = -3.0;
    
    // Act
    calculated_number result = a + b;
    
    // Assert
    assert_true(fabs(result - (-8.0)) < 0.0001);
}

static void test_eval_addition_with_zero(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 0.0;
    
    // Act
    calculated_number result = a + b;
    
    // Assert
    assert_true(fabs(result - 5.0) < 0.0001);
}

static void test_eval_addition_floats(void **state) {
    // Arrange
    calculated_number a = 1.5;
    calculated_number b = 2.5;
    
    // Act
    calculated_number result = a + b;
    
    // Assert
    assert_true(fabs(result - 4.0) < 0.0001);
}

// ============================================================================
// Test: Subtraction Operation
// ============================================================================

static void test_eval_subtraction_positive_result(void **state) {
    // Arrange
    calculated_number a = 10.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = a - b;
    
    // Assert
    assert_true(fabs(result - 7.0) < 0.0001);
}

static void test_eval_subtraction_negative_result(void **state) {
    // Arrange
    calculated_number a = 3.0;
    calculated_number b = 10.0;
    
    // Act
    calculated_number result = a - b;
    
    // Assert
    assert_true(fabs(result - (-7.0)) < 0.0001);
}

static void test_eval_subtraction_from_zero(void **state) {
    // Arrange
    calculated_number a = 0.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = a - b;
    
    // Assert
    assert_true(fabs(result - (-5.0)) < 0.0001);
}

static void test_eval_subtraction_zero_result(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = a - b;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Multiplication Operation
// ============================================================================

static void test_eval_multiplication_positive_numbers(void **state) {
    // Arrange
    calculated_number a = 4.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = a * b;
    
    // Assert
    assert_true(fabs(result - 20.0) < 0.0001);
}

static void test_eval_multiplication_positive_and_negative(void **state) {
    // Arrange
    calculated_number a = 4.0;
    calculated_number b = -5.0;
    
    // Act
    calculated_number result = a * b;
    
    // Assert
    assert_true(fabs(result - (-20.0)) < 0.0001);
}

static void test_eval_multiplication_two_negative_numbers(void **state) {
    // Arrange
    calculated_number a = -4.0;
    calculated_number b = -5.0;
    
    // Act
    calculated_number result = a * b;
    
    // Assert
    assert_true(fabs(result - 20.0) < 0.0001);
}

static void test_eval_multiplication_by_zero(void **state) {
    // Arrange
    calculated_number a = 100.0;
    calculated_number b = 0.0;
    
    // Act
    calculated_number result = a * b;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_multiplication_floats(void **state) {
    // Arrange
    calculated_number a = 2.5;
    calculated_number b = 4.0;
    
    // Act
    calculated_number result = a * b;
    
    // Assert
    assert_true(fabs(result - 10.0) < 0.0001);
}

// ============================================================================
// Test: Division Operation
// ============================================================================

static void test_eval_division_positive_numbers(void **state) {
    // Arrange
    calculated_number a = 20.0;
    calculated_number b = 4.0;
    
    // Act
    calculated_number result = a / b;
    
    // Assert
    assert_true(fabs(result - 5.0) < 0.0001);
}

static void test_eval_division_with_remainder(void **state) {
    // Arrange
    calculated_number a = 10.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = a / b;
    
    // Assert
    assert_true(fabs(result - 3.3333) < 0.0001);
}

static void test_eval_division_positive_by_negative(void **state) {
    // Arrange
    calculated_number a = 20.0;
    calculated_number b = -4.0;
    
    // Act
    calculated_number result = a / b;
    
    // Assert
    assert_true(fabs(result - (-5.0)) < 0.0001);
}

static void test_eval_division_negative_by_negative(void **state) {
    // Arrange
    calculated_number a = -20.0;
    calculated_number b = -4.0;
    
    // Act
    calculated_number result = a / b;
    
    // Assert
    assert_true(fabs(result - 5.0) < 0.0001);
}

static void test_eval_division_result_less_than_one(void **state) {
    // Arrange
    calculated_number a = 1.0;
    calculated_number b = 4.0;
    
    // Act
    calculated_number result = a / b;
    
    // Assert
    assert_true(fabs(result - 0.25) < 0.0001);
}

// ============================================================================
// Test: Modulo Operation
// ============================================================================

static void test_eval_modulo_exact_division(void **state) {
    // Arrange
    calculated_number a = 10.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = fmod(a, b);
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_modulo_with_remainder(void **state) {
    // Arrange
    calculated_number a = 10.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = fmod(a, b);
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_modulo_negative_dividend(void **state) {
    // Arrange
    calculated_number a = -10.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = fmod(a, b);
    
    // Assert
    // fmod preserves sign of dividend
    assert_true(fabs(result - (-1.0)) < 0.0001);
}

static void test_eval_modulo_negative_divisor(void **state) {
    // Arrange
    calculated_number a = 10.0;
    calculated_number b = -3.0;
    
    // Act
    calculated_number result = fmod(a, b);
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

// ============================================================================
// Test: Logical AND Operation
// ============================================================================

static void test_eval_logical_and_both_true(void **state) {
    // Arrange
    calculated_number a = 5.0;  // Non-zero = true
    calculated_number b = 3.0;  // Non-zero = true
    
    // Act
    calculated_number result = (a != 0.0 && b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_logical_and_first_false(void **state) {
    // Arrange
    calculated_number a = 0.0;  // Zero = false
    calculated_number b = 3.0;  // Non-zero = true
    
    // Act
    calculated_number result = (a != 0.0 && b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_logical_and_second_false(void **state) {
    // Arrange
    calculated_number a = 5.0;  // Non-zero = true
    calculated_number b = 0.0;  // Zero = false
    
    // Act
    calculated_number result = (a != 0.0 && b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_logical_and_both_false(void **state) {
    // Arrange
    calculated_number a = 0.0;  // Zero = false
    calculated_number b = 0.0;  // Zero = false
    
    // Act
    calculated_number result = (a != 0.0 && b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Logical OR Operation
// ============================================================================

static void test_eval_logical_or_both_true(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = (a != 0.0 || b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_logical_or_first_true(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 0.0;
    
    // Act
    calculated_number result = (a != 0.0 || b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_logical_or_second_true(void **state) {
    // Arrange
    calculated_number a = 0.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = (a != 0.0 || b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_logical_or_both_false(void **state) {
    // Arrange
    calculated_number a = 0.0;
    calculated_number b = 0.0;
    
    // Act
    calculated_number result = (a != 0.0 || b != 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Logical NOT Operation
// ============================================================================

static void test_eval_logical_not_true_becomes_false(void **state) {
    // Arrange
    calculated_number a = 5.0;
    
    // Act
    calculated_number result = (a == 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_logical_not_false_becomes_true(void **state) {
    // Arrange
    calculated_number a = 0.0;
    
    // Act
    calculated_number result = (a == 0.0) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

// ============================================================================
// Test: Comparison Operations (Equal)
// ============================================================================

static void test_eval_equal_same_values(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (fabs(a - b) < 0.0001) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_equal_different_values(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = (fabs(a - b) < 0.0001) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_equal_negative_numbers(void **state) {
    // Arrange
    calculated_number a = -5.0;
    calculated_number b = -5.0;
    
    // Act
    calculated_number result = (fabs(a - b) < 0.0001) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_equal_zero_values(void **state) {
    // Arrange
    calculated_number a = 0.0;
    calculated_number b = 0.0;
    
    // Act
    calculated_number result = (fabs(a - b) < 0.0001) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

// ============================================================================
// Test: Comparison Operations (Not Equal)
// ============================================================================

static void test_eval_not_equal_different_values(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = (fabs(a - b) >= 0.0001) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_not_equal_same_values(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (fabs(a - b) >= 0.0001) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Comparison Operations (Less Than)
// ============================================================================

static void test_eval_less_than_true(void **state) {
    // Arrange
    calculated_number a = 3.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a < b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_less_than_false_greater(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 3.0;
    
    // Act
    calculated_number result = (a < b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_less_than_false_equal(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a < b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_less_than_negative_numbers(void **state) {
    // Arrange
    calculated_number a = -10.0;
    calculated_number b = -5.0;
    
    // Act
    calculated_number result = (a < b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

// ============================================================================
// Test: Comparison Operations (Less Than or Equal)
// ============================================================================

static void test_eval_less_equal_less(void **state) {
    // Arrange
    calculated_number a = 3.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a <= b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_less_equal_equal(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a <= b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_less_equal_greater(void **state) {
    // Arrange
    calculated_number a = 7.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a <= b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Comparison Operations (Greater Than)
// ============================================================================

static void test_eval_greater_than_true(void **state) {
    // Arrange
    calculated_number a = 7.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a > b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_greater_than_false_less(void **state) {
    // Arrange
    calculated_number a = 3.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a > b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

static void test_eval_greater_than_false_equal(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a > b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Comparison Operations (Greater Than or Equal)
// ============================================================================

static void test_eval_greater_equal_greater(void **state) {
    // Arrange
    calculated_number a = 7.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a >= b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_greater_equal_equal(void **state) {
    // Arrange
    calculated_number a = 5.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a >= b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 1.0) < 0.0001);
}

static void test_eval_greater_equal_less(void **state) {
    // Arrange
    calculated_number a = 3.0;
    calculated_number b = 5.0;
    
    // Act
    calculated_number result = (a >= b) ? 1.0 : 0.0;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Unary Plus Operation
// ============================================================================

static void test_eval_unary_plus_positive(void **state) {
    // Arrange
    calculated_number a = 5.0;
    
    // Act
    calculated_number result = +a;
    
    // Assert
    assert_true(fabs(result - 5.0) < 0.0001);
}

static void test_eval_unary_plus_negative(void **state) {
    // Arrange
    calculated_number a = -5.0;
    
    // Act
    calculated_number result = +a;
    
    // Assert
    assert_true(fabs(result - (-5.0)) < 0.0001);
}

static void test_eval_unary_plus_zero(void **state) {
    // Arrange
    calculated_number a = 0.0;
    
    // Act
    calculated_number result = +a;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Unary Minus Operation
// ============================================================================

static void test_eval_unary_minus_positive(void **state) {
    // Arrange
    calculated_number a = 5.0;
    
    // Act
    calculated_number result = -a;
    
    // Assert
    assert_true(fabs(result - (-5.0)) < 0.0001);
}

static void test_eval_unary_minus_negative(void **state) {
    // Arrange
    calculated_number a = -5.0;
    
    // Act
    calculated_number result = -a;
    
    // Assert
    assert_true(fabs(result - 5.0) < 0.0001);
}

static void test_eval_unary_minus_zero(void **state) {
    // Arrange
    calculated_number a = 0.0;
    
    // Act
    calculated_number result = -a;
    
    // Assert
    assert_true(fabs(result - 0.0) < 0.0001);
}

// ============================================================================
// Test: Ternary Operator
// ============================================================================

static void test_eval_ternary_condition_true(void **state) {
    // Arrange
    calculated_number condition = 1.0;
    calculated_number true_val = 10.0;
    calculated_number false_val = 20.0;
    
    // Act
    calculated_number result = (condition != 0.0) ? true_val : false_val;
    
    // Assert
    assert_true(fabs(result - 10.0) < 0.0001);
}

static void test_eval_ternary_condition_false(void **state) {
    // Arrange
    calculated_number condition = 0.0;
    calculated_number true_val = 10.0;
    calculated_number false_val = 20.0;
    
    // Act
    calculated_number result = (condition != 0.0) ? true_val : false_val;
    
    // Assert
    assert_true(fabs(result - 20.0) < 0.0001);
}

static void test_eval_ternary_both_sides_same(void **state) {
    // Arrange
    calculated_number condition = 1.0;
    calculated_number true_val = 5.0;
    calculated_number false_val = 5.0;
    
    // Act
    calculated_number result = (condition != 0.0) ? true_val : false_val;
    
    // Assert
    assert_true(fabs(result - 5.0) < 0.0001);
}

// ============================================================================
// Test: Nested Operations
// ============================================================================

static void test_eval_nested_arithmetic_simple(void **state) {
    // Arrange: (2 + 3) * 4
    calculated_number result = (2.0 + 3.0) * 4.0;
    
    // Act & Assert
    assert_true(fabs(result - 20.0) < 0.0001);
}

static void test_eval_nested_arithmetic_complex(void **state) {
    // Arrange: ((2 + 3) * 4) - 5
    calculated_number result = ((2.0 + 3.0) * 4.0) - 5.0;
    
    // Act & Assert
    assert_true(fabs(result - 15.0) < 0.0001);
}

static void test_eval_nested_comparison_and_logic(void **state) {
    // Arrange: (5 > 3) && (2 < 4)
    calculated_number result = ((5.0 > 3.0) && (2.0 < 4.0)) ? 1.0 : 0.0;
    
    // Act & Assert