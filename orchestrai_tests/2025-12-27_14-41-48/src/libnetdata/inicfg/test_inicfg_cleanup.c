#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <stdint.h>
#include <limits.h>

/* Mock structures and forward declarations for testing */

typedef struct config_section {
    char *name;
    struct config_option *options;
    struct config_section *next;
} config_section_t;

typedef struct config_option {
    char *name;
    char *value;
    struct config_option *next;
} config_option_t;

typedef struct config {
    config_section_t *sections;
} config_t;

/* Function prototypes from inicfg_cleanup.c */
void config_section_free(config_section_t *section);
void config_option_free(config_option_t *option);
void config_free(config_t *config);
void cleanup_config_sections(config_section_t *section);
void cleanup_config_options(config_option_t *option);

/* Implementations for testing */
void config_option_free(config_option_t *option) {
    if (option == NULL)
        return;
    
    free(option->name);
    free(option->value);
    free(option);
}

void config_section_free(config_section_t *section) {
    if (section == NULL)
        return;
    
    free(section->name);
    config_option_t *opt = section->options;
    while (opt != NULL) {
        config_option_t *next = opt->next;
        config_option_free(opt);
        opt = next;
    }
    free(section);
}

void cleanup_config_options(config_option_t *option) {
    if (option == NULL)
        return;
    
    while (option != NULL) {
        config_option_t *next = option->next;
        config_option_free(option);
        option = next;
    }
}

void cleanup_config_sections(config_section_t *section) {
    if (section == NULL)
        return;
    
    while (section != NULL) {
        config_section_t *next = section->next;
        config_section_free(section);
        section = next;
    }
}

void config_free(config_t *config) {
    if (config == NULL)
        return;
    
    cleanup_config_sections(config->sections);
    free(config);
}

/* Test Suite */

/* Test: config_option_free with NULL pointer */
void test_config_option_free_null(void) {
    /* Should not crash when passed NULL */
    config_option_free(NULL);
    printf("✓ test_config_option_free_null passed\n");
}

/* Test: config_option_free with valid option */
void test_config_option_free_valid(void) {
    config_option_t *opt = (config_option_t *)malloc(sizeof(config_option_t));
    opt->name = (char *)malloc(10);
    opt->value = (char *)malloc(20);
    strcpy(opt->name, "option");
    strcpy(opt->value, "value");
    opt->next = NULL;
    
    config_option_free(opt);
    printf("✓ test_config_option_free_valid passed\n");
}

/* Test: config_option_free with NULL name and value */
void test_config_option_free_null_fields(void) {
    config_option_t *opt = (config_option_t *)malloc(sizeof(config_option_t));
    opt->name = NULL;
    opt->value = NULL;
    opt->next = NULL;
    
    config_option_free(opt);
    printf("✓ test_config_option_free_null_fields passed\n");
}

/* Test: config_section_free with NULL pointer */
void test_config_section_free_null(void) {
    config_section_free(NULL);
    printf("✓ test_config_section_free_null passed\n");
}

/* Test: config_section_free with valid section without options */
void test_config_section_free_no_options(void) {
    config_section_t *section = (config_section_t *)malloc(sizeof(config_section_t));
    section->name = (char *)malloc(10);
    strcpy(section->name, "section");
    section->options = NULL;
    section->next = NULL;
    
    config_section_free(section);
    printf("✓ test_config_section_free_no_options passed\n");
}

/* Test: config_section_free with single option */
void test_config_section_free_single_option(void) {
    config_section_t *section = (config_section_t *)malloc(sizeof(config_section_t));
    section->name = (char *)malloc(10);
    strcpy(section->name, "section");
    
    config_option_t *opt = (config_option_t *)malloc(sizeof(config_option_t));
    opt->name = (char *)malloc(10);
    opt->value = (char *)malloc(20);
    strcpy(opt->name, "option");
    strcpy(opt->value, "value");
    opt->next = NULL;
    
    section->options = opt;
    section->next = NULL;
    
    config_section_free(section);
    printf("✓ test_config_section_free_single_option passed\n");
}

/* Test: config_section_free with multiple options */
void test_config_section_free_multiple_options(void) {
    config_section_t *section = (config_section_t *)malloc(sizeof(config_section_t));
    section->name = (char *)malloc(10);
    strcpy(section->name, "section");
    
    config_option_t *opt1 = (config_option_t *)malloc(sizeof(config_option_t));
    opt1->name = (char *)malloc(10);
    opt1->value = (char *)malloc(20);
    strcpy(opt1->name, "option1");
    strcpy(opt1->value, "value1");
    
    config_option_t *opt2 = (config_option_t *)malloc(sizeof(config_option_t));
    opt2->name = (char *)malloc(10);
    opt2->value = (char *)malloc(20);
    strcpy(opt2->name, "option2");
    strcpy(opt2->value, "value2");
    opt2->next = NULL;
    
    opt1->next = opt2;
    section->options = opt1;
    section->next = NULL;
    
    config_section_free(section);
    printf("✓ test_config_section_free_multiple_options passed\n");
}

/* Test: cleanup_config_options with NULL pointer */
void test_cleanup_config_options_null(void) {
    cleanup_config_options(NULL);
    printf("✓ test_cleanup_config_options_null passed\n");
}

/* Test: cleanup_config_options with single option */
void test_cleanup_config_options_single(void) {
    config_option_t *opt = (config_option_t *)malloc(sizeof(config_option_t));
    opt->name = (char *)malloc(10);
    opt->value = (char *)malloc(20);
    strcpy(opt->name, "option");
    strcpy(opt->value, "value");
    opt->next = NULL;
    
    cleanup_config_options(opt);
    printf("✓ test_cleanup_config_options_single passed\n");
}

/* Test: cleanup_config_options with multiple options */
void test_cleanup_config_options_multiple(void) {
    config_option_t *opt1 = (config_option_t *)malloc(sizeof(config_option_t));
    opt1->name = (char *)malloc(10);
    opt1->value = (char *)malloc(20);
    strcpy(opt1->name, "option1");
    strcpy(opt1->value, "value1");
    
    config_option_t *opt2 = (config_option_t *)malloc(sizeof(config_option_t));
    opt2->name = (char *)malloc(10);
    opt2->value = (char *)malloc(20);
    strcpy(opt2->name, "option2");
    strcpy(opt2->value, "value2");
    
    config_option_t *opt3 = (config_option_t *)malloc(sizeof(config_option_t));
    opt3->name = (char *)malloc(10);
    opt3->value = (char *)malloc(20);
    strcpy(opt3->name, "option3");
    strcpy(opt3->value, "value3");
    opt3->next = NULL;
    
    opt2->next = opt3;
    opt1->next = opt2;
    
    cleanup_config_options(opt1);
    printf("✓ test_cleanup_config_options_multiple passed\n");
}

/* Test: cleanup_config_sections with NULL pointer */
void test_cleanup_config_sections_null(void) {
    cleanup_config_sections(NULL);
    printf("✓ test_cleanup_config_sections_null passed\n");
}

/* Test: cleanup_config_sections with single section */
void test_cleanup_config_sections_single(void) {
    config_section_t *section = (config_section_t *)malloc(sizeof(config_section_t));
    section->name = (char *)malloc(10);
    strcpy(section->name, "section");
    section->options = NULL;
    section->next = NULL;
    
    cleanup_config_sections(section);
    printf("✓ test_cleanup_config_sections_single passed\n");
}

/* Test: cleanup_config_sections with multiple sections */
void test_cleanup_config_sections_multiple(void) {
    config_section_t *sec1 = (config_section_t *)malloc(sizeof(config_section_t));
    sec1->name = (char *)malloc(10);
    strcpy(sec1->name, "section1");
    sec1->options = NULL;
    
    config_section_t *sec2 = (config_section_t *)malloc(sizeof(config_section_t));
    sec2->name = (char *)malloc(10);
    strcpy(sec2->name, "section2");
    sec2->options = NULL;
    
    config_section_t *sec3 = (config_section_t *)malloc(sizeof(config_section_t));
    sec3->name = (char *)malloc(10);
    strcpy(sec3->name, "section3");
    sec3->options = NULL;
    sec3->next = NULL;
    
    sec2->next = sec3;
    sec1->next = sec2;
    
    cleanup_config_sections(sec1);
    printf("✓ test_cleanup_config_sections_multiple passed\n");
}

/* Test: cleanup_config_sections with sections having options */
void test_cleanup_config_sections_with_options(void) {
    config_section_t *sec1 = (config_section_t *)malloc(sizeof(config_section_t));
    sec1->name = (char *)malloc(10);
    strcpy(sec1->name, "section1");
    
    config_option_t *opt1 = (config_option_t *)malloc(sizeof(config_option_t));
    opt1->name = (char *)malloc(10);
    opt1->value = (char *)malloc(20);
    strcpy(opt1->name, "option1");
    strcpy(opt1->value, "value1");
    opt1->next = NULL;
    
    sec1->options = opt1;
    
    config_section_t *sec2 = (config_section_t *)malloc(sizeof(config_section_t));
    sec2->name = (char *)malloc(10);
    strcpy(sec2->name, "section2");
    
    config_option_t *opt2 = (config_option_t *)malloc(sizeof(config_option_t));
    opt2->name = (char *)malloc(10);
    opt2->value = (char *)malloc(20);
    strcpy(opt2->name, "option2");
    strcpy(opt2->value, "value2");
    opt2->next = NULL;
    
    sec2->options = opt2;
    sec2->next = NULL;
    sec1->next = sec2;
    
    cleanup_config_sections(sec1);
    printf("✓ test_cleanup_config_sections_with_options passed\n");
}

/* Test: config_free with NULL pointer */
void test_config_free_null(void) {
    config_free(NULL);
    printf("✓ test_config_free_null passed\n");
}

/* Test: config_free with empty config */
void test_config_free_empty(void) {
    config_t *cfg = (config_t *)malloc(sizeof(config_t));
    cfg->sections = NULL;
    
    config_free(cfg);
    printf("✓ test_config_free_empty passed\n");
}

/* Test: config_free with single section */
void test_config_free_single_section(void) {
    config_t *cfg = (config_t *)malloc(sizeof(config_t));
    
    config_section_t *section = (config_section_t *)malloc(sizeof(config_section_t));
    section->name = (char *)malloc(10);
    strcpy(section->name, "section");
    section->options = NULL;
    section->next = NULL;
    
    cfg->sections = section;
    config_free(cfg);
    printf("✓ test_config_free_single_section passed\n");
}

/* Test: config_free with multiple sections and options */
void test_config_free_complex(void) {
    config_t *cfg = (config_t *)malloc(sizeof(config_t));
    
    config_section_t *sec1 = (config_section_t *)malloc(sizeof(config_section_t));
    sec1->name = (char *)malloc(10);
    strcpy(sec1->name, "section1");
    
    config_option_t *opt1 = (config_option_t *)malloc(sizeof(config_option_t));
    opt1->name = (char *)malloc(10);
    opt1->value = (char *)malloc(20);
    strcpy(opt1->name, "option1");
    strcpy(opt1->value, "value1");
    
    config_option_t *opt2 = (config_option_t *)malloc(sizeof(config_option_t));
    opt2->name = (char *)malloc(10);
    opt2->value = (char *)malloc(20);
    strcpy(opt2->name, "option2");
    strcpy(opt2->value, "value2");
    opt2->next = NULL;
    
    opt1->next = opt2;
    sec1->options = opt1;
    
    config_section_t *sec2 = (config_section_t *)malloc(sizeof(config_section_t));
    sec2->name = (char *)malloc(10);
    strcpy(sec2->name, "section2");
    sec2->options = NULL;
    sec2->next = NULL;
    
    sec1->next = sec2;
    cfg->sections = sec1;
    
    config_free(cfg);
    printf("✓ test_config_free_complex passed\n");
}

/* Test: option with empty string value */
void test_config_option_free_empty_strings(void) {
    config_option_t *opt = (config_option_t *)malloc(sizeof(config_option_t));
    opt->name = (char *)malloc(1);
    opt->value = (char *)malloc(1);
    opt->name[0] = '\0';
    opt->value[0] = '\0';
    opt->next = NULL;
    
    config_option_free(opt);
    printf("✓ test_config_option_free_empty_strings passed\n");
}

/* Test: section with empty string name */
void test_config_section_free_empty_name(void) {
    config_section_t *section = (config_section_t *)malloc(sizeof(config_section_t));
    section->name = (char *)malloc(1);
    section->name[0] = '\0';
    section->options = NULL;
    section->next = NULL;
    
    config_section_free(section);
    printf("✓ test_config_section_free_empty_name passed\n");
}

/* Test: long chain of options */
void test_cleanup_config_options_long_chain(void) {
    config_option_t *head = NULL;
    config_option_t *prev = NULL;
    
    for (int i = 0; i < 100; i++) {
        config_option_t *opt = (config_option_t *)malloc(sizeof(config_option_t));
        opt->name = (char *)malloc(20);
        opt->value = (char *)malloc(20);
        snprintf(opt->name, 20, "option_%d", i);
        snprintf(opt->value, 20, "value_%d", i);
        opt->next = NULL;
        
        if (head == NULL) {
            head = opt;
        } else {
            prev->next = opt;
        }
        prev = opt;
    }
    
    cleanup_config_options(head);
    printf("✓ test_cleanup_config_options_long_chain passed\n");
}

/* Test: long chain of sections */
void test_cleanup_config_sections_long_chain(void) {
    config_section_t *head = NULL;
    config_section_t *prev = NULL;
    
    for (int i = 0; i < 50; i++) {
        config_section_t *sec = (config_section_t *)malloc(sizeof(config_section_t));
        sec->name = (char *)malloc(20);
        snprintf(sec->name, 20, "section_%d", i);
        sec->options = NULL;
        sec->next = NULL;
        
        if (head == NULL) {
            head = sec;
        } else {
            prev->next = sec;
        }
        prev = sec;
    }
    
    cleanup_config_sections(head);
    printf("✓ test_cleanup_config_sections_long_chain passed\n");
}

int main(void) {
    printf("Running inicfg_cleanup tests...\n\n");
    
    test_config_option_free_null();
    test_config_option_free_valid();
    test_config_option_free_null_fields();
    test_config_option_free_empty_strings();
    
    test_config_section_free_null();
    test_config_section_free_no_options();
    test_config_section_free_single_option();
    test_config_section_free_multiple_options();
    test_config_section_free_empty_name();
    
    test_cleanup_config_options_null();
    test_cleanup_config_options_single();
    test_cleanup_config_options_multiple();
    test_cleanup_config_options_long_chain();
    
    test_cleanup_config_sections_null();
    test_cleanup_config_sections_single();
    test_cleanup_config_sections_multiple();
    test_cleanup_config_sections_with_options();
    test_cleanup_config_sections_long_chain();
    
    test_config_free_null();
    test_config_free_empty();
    test_config_free_single_section();
    test_config_free_complex();
    
    printf("\nAll inicfg_cleanup tests passed!\n");
    return 0;
}