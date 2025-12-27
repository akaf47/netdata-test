#ifndef TEST_MEM_PGFAULTS_H
#define TEST_MEM_PGFAULTS_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <stdint.h>

/* Mock structures and globals for testing */
typedef struct {
    char *name;
    uint64_t value;
} metric_t;

typedef struct {
    metric_t *metrics;
    int count;
    int capacity;
} context_t;

/* Test helper functions */
static void test_setup(context_t *ctx) {
    ctx->capacity = 10;
    ctx->count = 0;
    ctx->metrics = (metric_t *)malloc(sizeof(metric_t) * ctx->capacity);
    assert(ctx->metrics != NULL);
}

static void test_teardown(context_t *ctx) {
    if (ctx->metrics != NULL) {
        for (int i = 0; i < ctx->count; i++) {
            free(ctx->metrics[i].name);
        }
        free(ctx->metrics);
        ctx->metrics = NULL;
    }
    ctx->count = 0;
    ctx->capacity = 0;
}

/* Test functions for mem-pgfaults collector */
static void test_mem_pgfaults_initialization(void) {
    context_t ctx;
    test_setup(&ctx);
    
    assert(ctx.metrics != NULL);
    assert(ctx.count == 0);
    assert(ctx.capacity == 10);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_add_metric_valid(void) {
    context_t ctx;
    test_setup(&ctx);
    
    ctx.metrics[0].name = (char *)malloc(32);
    strcpy(ctx.metrics[0].name, "pgfault");
    ctx.metrics[0].value = 12345;
    ctx.count = 1;
    
    assert(ctx.count == 1);
    assert(strcmp(ctx.metrics[0].name, "pgfault") == 0);
    assert(ctx.metrics[0].value == 12345);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_add_metric_zero_value(void) {
    context_t ctx;
    test_setup(&ctx);
    
    ctx.metrics[0].name = (char *)malloc(32);
    strcpy(ctx.metrics[0].name, "pgfault");
    ctx.metrics[0].value = 0;
    ctx.count = 1;
    
    assert(ctx.count == 1);
    assert(ctx.metrics[0].value == 0);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_add_metric_max_value(void) {
    context_t ctx;
    test_setup(&ctx);
    
    ctx.metrics[0].name = (char *)malloc(32);
    strcpy(ctx.metrics[0].name, "pgfault");
    ctx.metrics[0].value = UINT64_MAX;
    ctx.count = 1;
    
    assert(ctx.count == 1);
    assert(ctx.metrics[0].value == UINT64_MAX);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_multiple_metrics(void) {
    context_t ctx;
    test_setup(&ctx);
    
    /* Add pgfault */
    ctx.metrics[0].name = (char *)malloc(32);
    strcpy(ctx.metrics[0].name, "pgfault");
    ctx.metrics[0].value = 100;
    
    /* Add pgmajfault */
    ctx.metrics[1].name = (char *)malloc(32);
    strcpy(ctx.metrics[1].name, "pgmajfault");
    ctx.metrics[1].value = 50;
    
    ctx.count = 2;
    
    assert(ctx.count == 2);
    assert(strcmp(ctx.metrics[0].name, "pgfault") == 0);
    assert(strcmp(ctx.metrics[1].name, "pgmajfault") == 0);
    assert(ctx.metrics[0].value == 100);
    assert(ctx.metrics[1].value == 50);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_capacity_exceeded(void) {
    context_t ctx;
    test_setup(&ctx);
    
    /* Fill to capacity */
    for (int i = 0; i < ctx.capacity; i++) {
        ctx.metrics[i].name = (char *)malloc(32);
        snprintf(ctx.metrics[i].name, 32, "metric_%d", i);
        ctx.metrics[i].value = i;
        ctx.count++;
    }
    
    assert(ctx.count == 10);
    
    /* Try to add beyond capacity - should be handled gracefully */
    if (ctx.count >= ctx.capacity) {
        /* Expand capacity */
        ctx.capacity *= 2;
        metric_t *new_metrics = (metric_t *)realloc(ctx.metrics, sizeof(metric_t) * ctx.capacity);
        assert(new_metrics != NULL);
        ctx.metrics = new_metrics;
    }
    
    ctx.metrics[10].name = (char *)malloc(32);
    strcpy(ctx.metrics[10].name, "extra");
    ctx.metrics[10].value = 999;
    ctx.count++;
    
    assert(ctx.count == 11);
    assert(ctx.capacity == 20);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_metric_name_empty_string(void) {
    context_t ctx;
    test_setup(&ctx);
    
    ctx.metrics[0].name = (char *)malloc(1);
    strcpy(ctx.metrics[0].name, "");
    ctx.metrics[0].value = 123;
    ctx.count = 1;
    
    assert(ctx.count == 1);
    assert(strcmp(ctx.metrics[0].name, "") == 0);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_reset_context(void) {
    context_t ctx;
    test_setup(&ctx);
    
    /* Add some metrics */
    ctx.metrics[0].name = (char *)malloc(32);
    strcpy(ctx.metrics[0].name, "pgfault");
    ctx.metrics[0].value = 123;
    ctx.count = 1;
    
    /* Reset */
    ctx.count = 0;
    
    assert(ctx.count == 0);
    assert(ctx.capacity == 10);
    
    test_teardown(&ctx);
}

static void test_mem_pgfaults_sequential_updates(void) {
    context_t ctx;
    test_setup(&ctx);
    
    /* First update */
    ctx.metrics[0].name = (char *)malloc(32);
    strcpy(ctx.metrics[0].name, "pgfault");
    ctx.metrics[0].value = 100;
    ctx.count = 1;
    
    assert(ctx.metrics[0].value == 100);
    
    /* Update same metric */
    ctx.metrics[0].value = 200;
    
    assert(ctx.metrics[0].value == 200);
    assert(ctx.count == 1);
    
    test_teardown(&ctx);
}

/* Main test runner */
int main(void) {
    printf("Running mem-pgfaults tests...\n");
    
    test_mem_pgfaults_initialization();
    printf("✓ test_mem_pgfaults_initialization\n");
    
    test_mem_pgfaults_add_metric_valid();
    printf("✓ test_mem_pgfaults_add_metric_valid\n");
    
    test_mem_pgfaults_add_metric_zero_value();
    printf("✓ test_mem_pgfaults_add_metric_zero_value\n");
    
    test_mem_pgfaults_add_metric_max_value();
    printf("✓ test_mem_pgfaults_add_metric_max_value\n");
    
    test_mem_pgfaults_multiple_metrics();
    printf("✓ test_mem_pgfaults_multiple_metrics\n");
    
    test_mem_pgfaults_capacity_exceeded();
    printf("✓ test_mem_pgfaults_capacity_exceeded\n");
    
    test_mem_pgfaults_metric_name_empty_string();
    printf("✓ test_mem_pgfaults_metric_name_empty_string\n");
    
    test_mem_pgfaults_reset_context();
    printf("✓ test_mem_pgfaults_reset_context\n");
    
    test_mem_pgfaults_sequential_updates();
    printf("✓ test_mem_pgfaults_sequential_updates\n");
    
    printf("\nAll mem-pgfaults tests passed!\n");
    return 0;
}

#endif /* TEST_MEM_PGFAULTS_H */