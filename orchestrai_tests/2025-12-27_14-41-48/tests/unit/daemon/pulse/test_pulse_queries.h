#ifndef TEST_PULSE_QUERIES_H
#define TEST_PULSE_QUERIES_H

#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>
#include <stddef.h>

/* Mock declarations for external dependencies */
int mock_database_query(const char *query);
int mock_prepare_statement(const char *query);
int mock_bind_parameter(int index, const char *value);
int mock_execute_statement(void);
void mock_finalize_statement(void);

/* Test fixtures */
typedef struct {
    char *test_query;
    int expected_result;
} test_fixture_t;

/* ============================================================================
 * Test Cases for Query Preparation and Validation
 * ============================================================================ */

static void test_pulse_query_initialization_success(void **state)
{
    /* Test successful query initialization */
    const char *query = "SELECT * FROM metrics";
    int result = mock_prepare_statement(query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_initialization_null_query(void **state)
{
    /* Test with NULL query parameter */
    int result = mock_prepare_statement(NULL);
    assert_int_not_equal(result, 0);
}

static void test_pulse_query_initialization_empty_query(void **state)
{
    /* Test with empty query string */
    const char *query = "";
    int result = mock_prepare_statement(query);
    assert_int_not_equal(result, 0);
}

static void test_pulse_query_parameter_binding_success(void **state)
{
    /* Test successful parameter binding */
    int result = mock_bind_parameter(0, "value1");
    assert_int_equal(result, 0);
}

static void test_pulse_query_parameter_binding_null_value(void **state)
{
    /* Test parameter binding with NULL value */
    int result = mock_bind_parameter(0, NULL);
    assert_int_equal(result, 0);  /* NULL should be handled gracefully */
}

static void test_pulse_query_parameter_binding_empty_value(void **state)
{
    /* Test parameter binding with empty string */
    int result = mock_bind_parameter(0, "");
    assert_int_equal(result, 0);
}

static void test_pulse_query_parameter_binding_invalid_index(void **state)
{
    /* Test parameter binding with invalid index */
    int result = mock_bind_parameter(-1, "value");
    assert_int_not_equal(result, 0);
}

static void test_pulse_query_parameter_binding_large_value(void **state)
{
    /* Test parameter binding with large string value */
    char large_value[10000];
    memset(large_value, 'A', sizeof(large_value) - 1);
    large_value[sizeof(large_value) - 1] = '\0';
    
    int result = mock_bind_parameter(0, large_value);
    assert_int_equal(result, 0);
}

static void test_pulse_query_execution_success(void **state)
{
    /* Test successful query execution */
    int result = mock_execute_statement();
    assert_int_equal(result, 0);
}

static void test_pulse_query_finalization_success(void **state)
{
    /* Test successful statement finalization */
    mock_finalize_statement();
    /* If no assertion fails, finalization was successful */
}

static void test_pulse_query_multiple_parameters_success(void **state)
{
    /* Test binding multiple parameters in sequence */
    assert_int_equal(mock_bind_parameter(0, "param1"), 0);
    assert_int_equal(mock_bind_parameter(1, "param2"), 0);
    assert_int_equal(mock_bind_parameter(2, "param3"), 0);
    assert_int_equal(mock_execute_statement(), 0);
}

static void test_pulse_query_special_characters_in_params(void **state)
{
    /* Test parameter binding with special characters */
    const char *special_value = "'; DROP TABLE metrics; --";
    int result = mock_bind_parameter(0, special_value);
    assert_int_equal(result, 0);  /* Should handle safely */
}

static void test_pulse_query_unicode_characters(void **state)
{
    /* Test parameter binding with unicode characters */
    const char *unicode_value = "测试数据中文";
    int result = mock_bind_parameter(0, unicode_value);
    assert_int_equal(result, 0);
}

static void test_pulse_query_with_whitespace(void **state)
{
    /* Test query with leading/trailing whitespace */
    const char *query = "  \n  SELECT * FROM metrics  \n  ";
    int result = mock_prepare_statement(query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_case_sensitivity(void **state)
{
    /* Test query handling with mixed case */
    const char *query1 = "SELECT * FROM metrics";
    const char *query2 = "select * from metrics";
    
    int result1 = mock_prepare_statement(query1);
    int result2 = mock_prepare_statement(query2);
    assert_int_equal(result1, result2);
}

static void test_pulse_query_complex_select(void **state)
{
    /* Test complex SELECT query */
    const char *complex_query = 
        "SELECT id, name, value FROM metrics WHERE type = ? "
        "AND timestamp >= ? AND timestamp <= ? ORDER BY timestamp DESC LIMIT ?";
    
    int result = mock_prepare_statement(complex_query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_with_joins(void **state)
{
    /* Test query with JOIN clauses */
    const char *join_query = 
        "SELECT m.id, m.value, d.name FROM metrics m "
        "INNER JOIN devices d ON m.device_id = d.id";
    
    int result = mock_prepare_statement(join_query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_with_aggregations(void **state)
{
    /* Test query with aggregation functions */
    const char *agg_query = 
        "SELECT device_id, COUNT(*), AVG(value), MAX(value), MIN(value) "
        "FROM metrics GROUP BY device_id HAVING COUNT(*) > ?";
    
    int result = mock_prepare_statement(agg_query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_insert_statement(void **state)
{
    /* Test INSERT query */
    const char *insert_query = 
        "INSERT INTO metrics (device_id, timestamp, value) VALUES (?, ?, ?)";
    
    int result = mock_prepare_statement(insert_query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_update_statement(void **state)
{
    /* Test UPDATE query */
    const char *update_query = 
        "UPDATE metrics SET value = ? WHERE id = ?";
    
    int result = mock_prepare_statement(update_query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_delete_statement(void **state)
{
    /* Test DELETE query */
    const char *delete_query = 
        "DELETE FROM metrics WHERE timestamp < ?";
    
    int result = mock_prepare_statement(delete_query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_numeric_parameters(void **state)
{
    /* Test binding numeric parameters as strings */
    assert_int_equal(mock_bind_parameter(0, "12345"), 0);
    assert_int_equal(mock_bind_parameter(1, "-999"), 0);
    assert_int_equal(mock_bind_parameter(2, "0"), 0);
    assert_int_equal(mock_bind_parameter(3, "3.14159"), 0);
}

static void test_pulse_query_boolean_parameters(void **state)
{
    /* Test binding boolean parameters */
    assert_int_equal(mock_bind_parameter(0, "true"), 0);
    assert_int_equal(mock_bind_parameter(1, "false"), 0);
    assert_int_equal(mock_bind_parameter(2, "1"), 0);
    assert_int_equal(mock_bind_parameter(3, "0"), 0);
}

static void test_pulse_query_very_long_query(void **state)
{
    /* Test with a very long query string */
    char long_query[5000];
    strcpy(long_query, "SELECT ");
    for (int i = 0; i < 100; i++) {
        strcat(long_query, "col");
        strcat(long_query, i % 10 == 0 ? " FROM metrics WHERE id = ? " : ", ");
    }
    
    int result = mock_prepare_statement(long_query);
    assert_int_equal(result, 0);
}

static void test_pulse_query_sequential_execution(void **state)
{
    /* Test multiple sequential query executions */
    for (int i = 0; i < 10; i++) {
        assert_int_equal(mock_prepare_statement("SELECT * FROM metrics"), 0);
        assert_int_equal(mock_execute_statement(), 0);
        mock_finalize_statement();
    }
}

/* ============================================================================
 * Test Suite Configuration
 * ============================================================================ */

const struct CMUnitTest pulse_queries_tests[] = {
    cmocka_unit_test(test_pulse_query_initialization_success),
    cmocka_unit_test(test_pulse_query_initialization_null_query),
    cmocka_unit_test(test_pulse_query_initialization_empty_query),
    cmocka_unit_test(test_pulse_query_parameter_binding_success),
    cmocka_unit_test(test_pulse_query_parameter_binding_null_value),
    cmocka_unit_test(test_pulse_query_parameter_binding_empty_value),
    cmocka_unit_test(test_pulse_query_parameter_binding_invalid_index),
    cmocka_unit_test(test_pulse_query_parameter_binding_large_value),
    cmocka_unit_test(test_pulse_query_execution_success),
    cmocka_unit_test(test_pulse_query_finalization_success),
    cmocka_unit_test(test_pulse_query_multiple_parameters_success),
    cmocka_unit_test(test_pulse_query_special_characters_in_params),
    cmocka_unit_test(test_pulse_query_unicode_characters),
    cmocka_unit_test(test_pulse_query_with_whitespace),
    cmocka_unit_test(test_pulse_query_case_sensitivity),
    cmocka_unit_test(test_pulse_query_complex_select),
    cmocka_unit_test(test_pulse_query_with_joins),
    cmocka_unit_test(test_pulse_query_with_aggregations),
    cmocka_unit_test(test_pulse_query_insert_statement),
    cmocka_unit_test(test_pulse_query_update_statement),
    cmocka_unit_test(test_pulse_query_delete_statement),
    cmocka_unit_test(test_pulse_query_numeric_parameters),
    cmocka_unit_test(test_pulse_query_boolean_parameters),
    cmocka_unit_test(test_pulse_query_very_long_query),
    cmocka_unit_test(test_pulse_query_sequential_execution),
};

int main(void)
{
    return cmocka_run_group_tests(
        pulse_queries_tests,
        sizeof(pulse_queries_tests) / sizeof(pulse_queries_tests[0])
    );
}

#endif /* TEST_PULSE_QUERIES_H */