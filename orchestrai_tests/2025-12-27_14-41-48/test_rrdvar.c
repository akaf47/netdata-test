#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>

/* Mock structures and function declarations */
typedef struct {
    char *name;
    double value;
    int updated;
    struct rrdvar *next;
} RRDVAR;

typedef struct {
    RRDVAR *variables;
    int count;
} RRDHOST;

typedef struct {
    RRDVAR *variables;
    int count;
} RRDDIM;

/* Forward declarations of functions to test */
extern RRDVAR *rrdvar_create(RRDHOST *host, RRDDIM *dimension, const char *name, double value);
extern RRDVAR *rrdvar_find_and_create(RRDHOST *host, const char *name, double value);
extern RRDVAR *rrdvar_find(RRDHOST *host, const char *name);
extern void rrdvar_free(RRDHOST *host, RRDVAR *rv);
extern void rrdvar_rename(RRDHOST *host, RRDVAR *rv, const char *new_name);
extern double rrdvar_get_value(RRDHOST *host, RRDVAR *rv);
extern void rrdvar_set_value(RRDHOST *host, RRDVAR *rv, double value);

/* Test fixtures */
static RRDHOST* test_host = NULL;
static RRDDIM* test_dimension = NULL;

static int setup(void **state) {
    test_host = malloc(sizeof(RRDHOST));
    if (!test_host) return -1;
    
    test_host->variables = NULL;
    test_host->count = 0;
    
    test_dimension = malloc(sizeof(RRDDIM));
    if (!test_dimension) {
        free(test_host);
        return -1;
    }
    
    test_dimension->variables = NULL;
    test_dimension->count = 0;
    
    *state = test_host;
    return 0;
}

static int teardown(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    if (host) {
        RRDVAR *current = host->variables;
        while (current) {
            RRDVAR *next = current->next;
            if (current->name) free(current->name);
            free(current);
            current = next;
        }
        free(host);
    }
    
    if (test_dimension) {
        RRDVAR *current = test_dimension->variables;
        while (current) {
            RRDVAR *next = current->next;
            if (current->name) free(current->name);
            free(current);
            current = next;
        }
        free(test_dimension);
    }
    
    test_host = NULL;
    test_dimension = NULL;
    return 0;
}

/* Tests for rrdvar_create */
static void test_rrdvar_create_with_valid_params(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "test_variable";
    double test_value = 42.5;
    
    RRDVAR *result = rrdvar_create(host, test_dimension, var_name, test_value);
    
    assert_non_null(result);
    assert_string_equal(result->name, var_name);
    assert_true(result->value == test_value);
    assert_non_null(result->next);
}

static void test_rrdvar_create_with_null_host(void **state) {
    const char *var_name = "test_variable";
    double test_value = 42.5;
    
    RRDVAR *result = rrdvar_create(NULL, test_dimension, var_name, test_value);
    
    assert_null(result);
}

static void test_rrdvar_create_with_null_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    double test_value = 42.5;
    
    RRDVAR *result = rrdvar_create(host, test_dimension, NULL, test_value);
    
    assert_null(result);
}

static void test_rrdvar_create_with_empty_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    double test_value = 42.5;
    
    RRDVAR *result = rrdvar_create(host, test_dimension, "", test_value);
    
    assert_null(result);
}

static void test_rrdvar_create_with_zero_value(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "zero_var";
    
    RRDVAR *result = rrdvar_create(host, test_dimension, var_name, 0.0);
    
    assert_non_null(result);
    assert_true(result->value == 0.0);
}

static void test_rrdvar_create_with_negative_value(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "negative_var";
    double test_value = -123.456;
    
    RRDVAR *result = rrdvar_create(host, test_dimension, var_name, test_value);
    
    assert_non_null(result);
    assert_true(result->value == test_value);
}

static void test_rrdvar_create_with_large_value(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "large_var";
    double test_value = 1.7976931348623157e+308; /* near DBL_MAX */
    
    RRDVAR *result = rrdvar_create(host, test_dimension, var_name, test_value);
    
    assert_non_null(result);
    assert_true(result->value == test_value);
}

static void test_rrdvar_create_with_small_value(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "small_var";
    double test_value = 2.2250738585072014e-308; /* near DBL_MIN */
    
    RRDVAR *result = rrdvar_create(host, test_dimension, var_name, test_value);
    
    assert_non_null(result);
    assert_true(result->value == test_value);
}

static void test_rrdvar_create_multiple_variables(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var1 = rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *var2 = rrdvar_create(host, test_dimension, "var2", 2.0);
    RRDVAR *var3 = rrdvar_create(host, test_dimension, "var3", 3.0);
    
    assert_non_null(var1);
    assert_non_null(var2);
    assert_non_null(var3);
    assert_int_equal(host->count, 3);
}

/* Tests for rrdvar_find */
static void test_rrdvar_find_existing_variable(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "existing_var";
    
    rrdvar_create(host, test_dimension, var_name, 42.5);
    RRDVAR *result = rrdvar_find(host, var_name);
    
    assert_non_null(result);
    assert_string_equal(result->name, var_name);
    assert_true(result->value == 42.5);
}

static void test_rrdvar_find_non_existing_variable(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *result = rrdvar_find(host, "non_existing");
    
    assert_null(result);
}

static void test_rrdvar_find_with_null_host(void **state) {
    RRDVAR *result = rrdvar_find(NULL, "var_name");
    
    assert_null(result);
}

static void test_rrdvar_find_with_null_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *result = rrdvar_find(host, NULL);
    
    assert_null(result);
}

static void test_rrdvar_find_with_empty_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *result = rrdvar_find(host, "");
    
    assert_null(result);
}

static void test_rrdvar_find_in_empty_host(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *result = rrdvar_find(host, "any_var");
    
    assert_null(result);
}

static void test_rrdvar_find_case_sensitive(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    rrdvar_create(host, test_dimension, "MyVar", 42.0);
    RRDVAR *result = rrdvar_find(host, "myvar");
    
    assert_null(result);
}

static void test_rrdvar_find_multiple_variables(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    rrdvar_create(host, test_dimension, "var1", 1.0);
    rrdvar_create(host, test_dimension, "var2", 2.0);
    rrdvar_create(host, test_dimension, "var3", 3.0);
    
    RRDVAR *result = rrdvar_find(host, "var2");
    
    assert_non_null(result);
    assert_string_equal(result->name, "var2");
    assert_true(result->value == 2.0);
}

/* Tests for rrdvar_find_and_create */
static void test_rrdvar_find_and_create_existing(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "test_var";
    
    RRDVAR *created = rrdvar_create(host, test_dimension, var_name, 10.0);
    RRDVAR *found = rrdvar_find_and_create(host, var_name, 20.0);
    
    assert_ptr_equal(found, created);
    assert_true(found->value == 10.0); /* Original value unchanged */
}

static void test_rrdvar_find_and_create_new(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *var_name = "new_var";
    double test_value = 55.5;
    
    RRDVAR *result = rrdvar_find_and_create(host, var_name, test_value);
    
    assert_non_null(result);
    assert_string_equal(result->name, var_name);
    assert_true(result->value == test_value);
}

static void test_rrdvar_find_and_create_with_null_host(void **state) {
    RRDVAR *result = rrdvar_find_and_create(NULL, "var_name", 10.0);
    
    assert_null(result);
}

static void test_rrdvar_find_and_create_with_null_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *result = rrdvar_find_and_create(host, NULL, 10.0);
    
    assert_null(result);
}

static void test_rrdvar_find_and_create_with_empty_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *result = rrdvar_find_and_create(host, "", 10.0);
    
    assert_null(result);
}

/* Tests for rrdvar_free */
static void test_rrdvar_free_single_variable(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "var_to_free", 42.0);
    int count_before = host->count;
    
    rrdvar_free(host, var);
    
    assert_int_equal(host->count, count_before - 1);
    assert_null(rrdvar_find(host, "var_to_free"));
}

static void test_rrdvar_free_from_multiple(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var1 = rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *var2 = rrdvar_create(host, test_dimension, "var2", 2.0);
    RRDVAR *var3 = rrdvar_create(host, test_dimension, "var3", 3.0);
    
    rrdvar_free(host, var2);
    
    assert_non_null(rrdvar_find(host, "var1"));
    assert_null(rrdvar_find(host, "var2"));
    assert_non_null(rrdvar_find(host, "var3"));
}

static void test_rrdvar_free_first_variable(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var1 = rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *var2 = rrdvar_create(host, test_dimension, "var2", 2.0);
    
    rrdvar_free(host, var1);
    
    assert_null(rrdvar_find(host, "var1"));
    assert_non_null(rrdvar_find(host, "var2"));
}

static void test_rrdvar_free_last_variable(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var1 = rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *var2 = rrdvar_create(host, test_dimension, "var2", 2.0);
    
    rrdvar_free(host, var2);
    
    assert_non_null(rrdvar_find(host, "var1"));
    assert_null(rrdvar_find(host, "var2"));
}

static void test_rrdvar_free_with_null_var(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    rrdvar_create(host, test_dimension, "var1", 1.0);
    int count_before = host->count;
    
    rrdvar_free(host, NULL);
    
    assert_int_equal(host->count, count_before); /* No change */
}

static void test_rrdvar_free_with_null_host(void **state) {
    /* Should not crash */
    rrdvar_free(NULL, NULL);
}

/* Tests for rrdvar_rename */
static void test_rrdvar_rename_success(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const char *old_name = "old_name";
    const char *new_name = "new_name";
    
    RRDVAR *var = rrdvar_create(host, test_dimension, old_name, 42.0);
    rrdvar_rename(host, var, new_name);
    
    assert_string_equal(var->name, new_name);
    assert_null(rrdvar_find(host, old_name));
    assert_non_null(rrdvar_find(host, new_name));
}

static void test_rrdvar_rename_with_null_new_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "old_name", 42.0);
    rrdvar_rename(host, var, NULL);
    
    assert_string_equal(var->name, "old_name"); /* Name unchanged */
}

static void test_rrdvar_rename_with_empty_new_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "old_name", 42.0);
    rrdvar_rename(host, var, "");
    
    assert_string_equal(var->name, "old_name"); /* Name unchanged */
}

static void test_rrdvar_rename_with_null_var(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    /* Should not crash */
    rrdvar_rename(host, NULL, "new_name");
}

static void test_rrdvar_rename_with_null_host(void **state) {
    RRDVAR *var = malloc(sizeof(RRDVAR));
    var->name = malloc(16);
    strcpy(var->name, "test");
    var->next = NULL;
    
    /* Should not crash */
    rrdvar_rename(NULL, var, "new_name");
    
    free(var->name);
    free(var);
}

static void test_rrdvar_rename_to_duplicate_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var1 = rrdvar_create(host, test_dimension, "var1", 1.0);
    RRDVAR *var2 = rrdvar_create(host, test_dimension, "var2", 2.0);
    
    rrdvar_rename(host, var2, "var1");
    
    /* Both variables should still exist (rename overwrites) */
    assert_non_null(rrdvar_find(host, "var1"));
}

static void test_rrdvar_rename_same_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 42.0);
    rrdvar_rename(host, var, "test_var");
    
    assert_string_equal(var->name, "test_var");
    assert_true(var->value == 42.0);
}

/* Tests for rrdvar_get_value */
static void test_rrdvar_get_value_valid(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    double test_value = 123.456;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", test_value);
    double result = rrdvar_get_value(host, var);
    
    assert_true(result == test_value);
}

static void test_rrdvar_get_value_zero(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 0.0);
    double result = rrdvar_get_value(host, var);
    
    assert_true(result == 0.0);
}

static void test_rrdvar_get_value_negative(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", -999.999);
    double result = rrdvar_get_value(host, var);
    
    assert_true(result == -999.999);
}

static void test_rrdvar_get_value_null_var(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    double result = rrdvar_get_value(host, NULL);
    
    assert_true(result == 0.0 || isnan(result)); /* Expect 0 or NaN */
}

static void test_rrdvar_get_value_null_host(void **state) {
    RRDVAR *var = malloc(sizeof(RRDVAR));
    var->name = malloc(16);
    strcpy(var->name, "test");
    var->value = 42.0;
    var->next = NULL;
    
    double result = rrdvar_get_value(NULL, var);
    
    assert_true(result == 42.0);
    
    free(var->name);
    free(var);
}

/* Tests for rrdvar_set_value */
static void test_rrdvar_set_value_valid(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 10.0);
    double new_value = 99.5;
    
    rrdvar_set_value(host, var, new_value);
    
    assert_true(var->value == new_value);
    assert_int_equal(var->updated, 1);
}

static void test_rrdvar_set_value_to_zero(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 42.0);
    rrdvar_set_value(host, var, 0.0);
    
    assert_true(var->value == 0.0);
    assert_int_equal(var->updated, 1);
}

static void test_rrdvar_set_value_negative(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 100.0);
    rrdvar_set_value(host, var, -50.5);
    
    assert_true(var->value == -50.5);
    assert_int_equal(var->updated, 1);
}

static void test_rrdvar_set_value_large(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 1.0);
    double large_value = 1.7976931348623157e+308;
    
    rrdvar_set_value(host, var, large_value);
    
    assert_true(var->value == large_value);
}

static void test_rrdvar_set_value_small(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 100.0);
    double small_value = 2.2250738585072014e-308;
    
    rrdvar_set_value(host, var, small_value);
    
    assert_true(var->value == small_value);
}

static void test_rrdvar_set_value_null_var(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    /* Should not crash */
    rrdvar_set_value(host, NULL, 42.0);
}

static void test_rrdvar_set_value_null_host(void **state) {
    RRDVAR *var = malloc(sizeof(RRDVAR));
    var->name = malloc(16);
    strcpy(var->name, "test");
    var->value = 10.0;
    var->updated = 0;
    var->next = NULL;
    
    /* Should still update value */
    rrdvar_set_value(NULL, var, 99.0);
    
    assert_true(var->value == 99.0);
    
    free(var->name);
    free(var);
}

static void test_rrdvar_set_value_multiple_times(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "test_var", 1.0);
    
    rrdvar_set_value(host, var, 2.0);
    assert_true(var->value == 2.0);
    
    rrdvar_set_value(host, var, 3.0);
    assert_true(var->value == 3.0);
    
    rrdvar_set_value(host, var, 4.0);
    assert_true(var->value == 4.0);
}

/* Edge cases and integration tests */
static void test_rrdvar_lifecycle(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    /* Create */
    RRDVAR *var = rrdvar_create(host, test_dimension, "lifecycle_var", 10.0);
    assert_non_null(var);
    
    /* Find */
    RRDVAR *found = rrdvar_find(host, "lifecycle_var");
    assert_ptr_equal(found, var);
    
    /* Set value */
    rrdvar_set_value(host, var, 20.0);
    assert_true(rrdvar_get_value(host, var) == 20.0);
    
    /* Rename */
    rrdvar_rename(host, var, "renamed_var");
    assert_string_equal(var->name, "renamed_var");
    assert_non_null(rrdvar_find(host, "renamed_var"));
    
    /* Free */
    rrdvar_free(host, var);
    assert_null(rrdvar_find(host, "renamed_var"));
}

static void test_rrdvar_stress_many_variables(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    const int num_vars = 100;
    
    for (int i = 0; i < num_vars; i++) {
        char name[32];
        snprintf(name, sizeof(name), "var_%d", i);
        rrdvar_create(host, test_dimension, name, (double)i);
    }
    
    assert_int_equal(host->count, num_vars);
    
    /* Verify all can be found */
    for (int i = 0; i < num_vars; i++) {
        char name[32];
        snprintf(name, sizeof(name), "var_%d", i);
        RRDVAR *found = rrdvar_find(host, name);
        assert_non_null(found);
        assert_true(found->value == (double)i);
    }
}

static void test_rrdvar_special_characters_in_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    RRDVAR *var = rrdvar_create(host, test_dimension, "var-with_special.chars$", 42.0);
    assert_non_null(var);
    assert_string_equal(var->name, "var-with_special.chars$");
    
    RRDVAR *found = rrdvar_find(host, "var-with_special.chars$");
    assert_non_null(found);
}

static void test_rrdvar_long_name(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    char long_name[1024];
    memset(long_name, 'a', 1023);
    long_name[1023] = '\0';
    
    RRDVAR *var = rrdvar_create(host, test_dimension, long_name, 42.0);
    assert_non_null(var);
    assert_string_equal(var->name, long_name);
    
    RRDVAR *found = rrdvar_find(host, long_name);
    assert_non_null(found);
}

static void test_rrdvar_float_precision(void **state) {
    RRDHOST *host = (RRDHOST *)*state;
    
    double value = 3.14159265358979323846;
    RRDVAR *var = rrdvar_create(host, test_dimension, "pi", value);
    
    double retrieved = rrdvar_get_value(host, var);
    assert_true(retrieved == value);
}

/* Test suite main */
int main(void) {
    const struct CMUnitTest tests[] = {
        /* rrdvar_create tests */
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_valid_params, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_null_host, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_null_name, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_empty_name, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_zero_value, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_negative_value, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_large_value, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_with_small_value, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_create_multiple_variables, setup, teardown),
        
        /* rrdvar_find tests */
        cmocka_unit_test_setup_teardown(test_rrdvar_find_existing_variable, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_find_non_existing_variable, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_find_with_null_host, setup, teardown),
        cmocka_unit_test_setup_teardown(test_rrdvar_find_with_null_name, setup, teardown),
        cmoc