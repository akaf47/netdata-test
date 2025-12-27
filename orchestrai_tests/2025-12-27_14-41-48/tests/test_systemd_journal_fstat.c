#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <sys/stat.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>

/* Mock structures and prototypes from systemd-journal-fstat.c */
struct inode_entry {
    ino_t inode;
    dev_t device;
};

/* Forward declarations of functions being tested */
int systemd_journal_fstat_get_size(const char *path);
int systemd_journal_fstat_get_inode(const char *path, ino_t *inode, dev_t *device);
int systemd_journal_fstat_compare_files(const char *path1, const char *path2);
int systemd_journal_fstat_is_regular_file(const char *path);

/* ============================================================================
   TEST SUITE: systemd_journal_fstat_get_size
   ============================================================================ */

static int test_get_size_setup(void **state) {
    return 0;
}

static int test_get_size_teardown(void **state) {
    return 0;
}

/* Test: Get file size successfully */
static void test_systemd_journal_fstat_get_size_success(void **state) {
    const char *test_file = "/tmp/test_journal_size_123456.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    
    /* Write 1024 bytes */
    for (int i = 0; i < 1024; i++) {
        fputc('A', fp);
    }
    fclose(fp);
    
    int size = systemd_journal_fstat_get_size(test_file);
    assert_int_equal(size, 1024);
    
    unlink(test_file);
}

/* Test: Get size of empty file */
static void test_systemd_journal_fstat_get_size_empty_file(void **state) {
    const char *test_file = "/tmp/test_journal_empty_123456.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    fclose(fp);
    
    int size = systemd_journal_fstat_get_size(test_file);
    assert_int_equal(size, 0);
    
    unlink(test_file);
}

/* Test: Get size of non-existent file */
static void test_systemd_journal_fstat_get_size_nonexistent(void **state) {
    const char *nonexistent = "/tmp/this_file_does_not_exist_987654321.tmp";
    int size = systemd_journal_fstat_get_size(nonexistent);
    assert_int_equal(size, -1);
}

/* Test: Get size with NULL path */
static void test_systemd_journal_fstat_get_size_null_path(void **state) {
    int size = systemd_journal_fstat_get_size(NULL);
    assert_int_equal(size, -1);
}

/* Test: Get size with empty path string */
static void test_systemd_journal_fstat_get_size_empty_path(void **state) {
    int size = systemd_journal_fstat_get_size("");
    assert_int_equal(size, -1);
}

/* Test: Get size of large file */
static void test_systemd_journal_fstat_get_size_large_file(void **state) {
    const char *test_file = "/tmp/test_journal_large_123456.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    
    int large_size = 1024 * 1024; /* 1MB */
    char buffer[4096];
    memset(buffer, 'X', sizeof(buffer));
    
    for (int i = 0; i < large_size / sizeof(buffer); i++) {
        fwrite(buffer, 1, sizeof(buffer), fp);
    }
    fclose(fp);
    
    int size = systemd_journal_fstat_get_size(test_file);
    assert_int_equal(size, large_size);
    
    unlink(test_file);
}

/* ============================================================================
   TEST SUITE: systemd_journal_fstat_get_inode
   ============================================================================ */

static int test_get_inode_setup(void **state) {
    return 0;
}

static int test_get_inode_teardown(void **state) {
    return 0;
}

/* Test: Get inode successfully */
static void test_systemd_journal_fstat_get_inode_success(void **state) {
    const char *test_file = "/tmp/test_journal_inode_123456.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    fprintf(fp, "test content");
    fclose(fp);
    
    ino_t inode = 0;
    dev_t device = 0;
    int result = systemd_journal_fstat_get_inode(test_file, &inode, &device);
    
    assert_int_equal(result, 0);
    assert_true(inode > 0);
    assert_true(device > 0);
    
    unlink(test_file);
}

/* Test: Get inode with NULL inode pointer */
static void test_systemd_journal_fstat_get_inode_null_inode_ptr(void **state) {
    const char *test_file = "/tmp/test_journal_inode_null.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    fclose(fp);
    
    dev_t device = 0;
    int result = systemd_journal_fstat_get_inode(test_file, NULL, &device);
    
    assert_int_equal(result, -1);
    
    unlink(test_file);
}

/* Test: Get inode with NULL device pointer */
static void test_systemd_journal_fstat_get_inode_null_device_ptr(void **state) {
    const char *test_file = "/tmp/test_journal_device_null.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    fclose(fp);
    
    ino_t inode = 0;
    int result = systemd_journal_fstat_get_inode(test_file, &inode, NULL);
    
    assert_int_equal(result, -1);
    
    unlink(test_file);
}

/* Test: Get inode with NULL path */
static void test_systemd_journal_fstat_get_inode_null_path(void **state) {
    ino_t inode = 0;
    dev_t device = 0;
    int result = systemd_journal_fstat_get_inode(NULL, &inode, &device);
    
    assert_int_equal(result, -1);
}

/* Test: Get inode with non-existent file */
static void test_systemd_journal_fstat_get_inode_nonexistent(void **state) {
    const char *nonexistent = "/tmp/nonexistent_inode_file_987654321.tmp";
    ino_t inode = 0;
    dev_t device = 0;
    int result = systemd_journal_fstat_get_inode(nonexistent, &inode, &device);
    
    assert_int_equal(result, -1);
}

/* Test: Get inode with both pointers NULL */
static void test_systemd_journal_fstat_get_inode_all_null_ptrs(void **state) {
    const char *test_file = "/tmp/test_journal_all_null.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    fclose(fp);
    
    int result = systemd_journal_fstat_get_inode(test_file, NULL, NULL);
    assert_int_equal(result, -1);
    
    unlink(test_file);
}

/* ============================================================================
   TEST SUITE: systemd_journal_fstat_compare_files
   ============================================================================ */

static int test_compare_files_setup(void **state) {
    return 0;
}

static int test_compare_files_teardown(void **state) {
    return 0;
}

/* Test: Compare identical files */
static void test_systemd_journal_fstat_compare_files_identical(void **state) {
    const char *file1 = "/tmp/test_journal_compare1_123456.tmp";
    const char *file2 = "/tmp/test_journal_compare2_123456.tmp";
    
    FILE *fp1 = fopen(file1, "w");
    assert_non_null(fp1);
    fprintf(fp1, "identical content");
    fclose(fp1);
    
    FILE *fp2 = fopen(file2, "w");
    assert_non_null(fp2);
    fprintf(fp2, "identical content");
    fclose(fp2);
    
    int result = systemd_journal_fstat_compare_files(file1, file2);
    assert_int_equal(result, 1);
    
    unlink(file1);
    unlink(file2);
}

/* Test: Compare different files */
static void test_systemd_journal_fstat_compare_files_different(void **state) {
    const char *file1 = "/tmp/test_journal_diff1_123456.tmp";
    const char *file2 = "/tmp/test_journal_diff2_123456.tmp";
    
    FILE *fp1 = fopen(file1, "w");
    assert_non_null(fp1);
    fprintf(fp1, "content1");
    fclose(fp1);
    
    FILE *fp2 = fopen(file2, "w");
    assert_non_null(fp2);
    fprintf(fp2, "content2");
    fclose(fp2);
    
    int result = systemd_journal_fstat_compare_files(file1, file2);
    assert_int_equal(result, 0);
    
    unlink(file1);
    unlink(file2);
}

/* Test: Compare with non-existent first file */
static void test_systemd_journal_fstat_compare_files_first_nonexistent(void **state) {
    const char *file1 = "/tmp/nonexistent_comp1_987654321.tmp";
    const char *file2 = "/tmp/test_journal_comp2_123456.tmp";
    
    FILE *fp2 = fopen(file2, "w");
    assert_non_null(fp2);
    fprintf(fp2, "content");
    fclose(fp2);
    
    int result = systemd_journal_fstat_compare_files(file1, file2);
    assert_int_equal(result, -1);
    
    unlink(file2);
}

/* Test: Compare with non-existent second file */
static void test_systemd_journal_fstat_compare_files_second_nonexistent(void **state) {
    const char *file1 = "/tmp/test_journal_comp1_123456.tmp";
    const char *file2 = "/tmp/nonexistent_comp2_987654321.tmp";
    
    FILE *fp1 = fopen(file1, "w");
    assert_non_null(fp1);
    fprintf(fp1, "content");
    fclose(fp1);
    
    int result = systemd_journal_fstat_compare_files(file1, file2);
    assert_int_equal(result, -1);
    
    unlink(file1);
}

/* Test: Compare with NULL first path */
static void test_systemd_journal_fstat_compare_files_null_first_path(void **state) {
    const char *file2 = "/tmp/test_journal_nullfirst_123456.tmp";
    
    FILE *fp2 = fopen(file2, "w");
    assert_non_null(fp2);
    fclose(fp2);
    
    int result = systemd_journal_fstat_compare_files(NULL, file2);
    assert_int_equal(result, -1);
    
    unlink(file2);
}

/* Test: Compare with NULL second path */
static void test_systemd_journal_fstat_compare_files_null_second_path(void **state) {
    const char *file1 = "/tmp/test_journal_nullsecond_123456.tmp";
    
    FILE *fp1 = fopen(file1, "w");
    assert_non_null(fp1);
    fclose(fp1);
    
    int result = systemd_journal_fstat_compare_files(file1, NULL);
    assert_int_equal(result, -1);
    
    unlink(file1);
}

/* Test: Compare with both NULL paths */
static void test_systemd_journal_fstat_compare_files_both_null_paths(void **state) {
    int result = systemd_journal_fstat_compare_files(NULL, NULL);
    assert_int_equal(result, -1);
}

/* Test: Compare same file with itself */
static void test_systemd_journal_fstat_compare_files_same_file(void **state) {
    const char *file = "/tmp/test_journal_samefile_123456.tmp";
    
    FILE *fp = fopen(file, "w");
    assert_non_null(fp);
    fprintf(fp, "same file");
    fclose(fp);
    
    int result = systemd_journal_fstat_compare_files(file, file);
    assert_int_equal(result, 1);
    
    unlink(file);
}

/* Test: Compare empty files */
static void test_systemd_journal_fstat_compare_files_empty_files(void **state) {
    const char *file1 = "/tmp/test_journal_empty1_123456.tmp";
    const char *file2 = "/tmp/test_journal_empty2_123456.tmp";
    
    FILE *fp1 = fopen(file1, "w");
    assert_non_null(fp1);
    fclose(fp1);
    
    FILE *fp2 = fopen(file2, "w");
    assert_non_null(fp2);
    fclose(fp2);
    
    int result = systemd_journal_fstat_compare_files(file1, file2);
    assert_int_equal(result, 1);
    
    unlink(file1);
    unlink(file2);
}

/* ============================================================================
   TEST SUITE: systemd_journal_fstat_is_regular_file
   ============================================================================ */

static int test_is_regular_file_setup(void **state) {
    return 0;
}

static int test_is_regular_file_teardown(void **state) {
    return 0;
}

/* Test: Regular file returns true */
static void test_systemd_journal_fstat_is_regular_file_regular(void **state) {
    const char *test_file = "/tmp/test_journal_regular_123456.tmp";
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    fprintf(fp, "regular file");
    fclose(fp);
    
    int result = systemd_journal_fstat_is_regular_file(test_file);
    assert_int_equal(result, 1);
    
    unlink(test_file);
}

/* Test: Directory returns false */
static void test_systemd_journal_fstat_is_regular_file_directory(void **state) {
    const char *test_dir = "/tmp/test_journal_dir_123456";
    mkdir(test_dir, 0755);
    
    int result = systemd_journal_fstat_is_regular_file(test_dir);
    assert_int_equal(result, 0);
    
    rmdir(test_dir);
}

/* Test: Non-existent file returns error */
static void test_systemd_journal_fstat_is_regular_file_nonexistent(void **state) {
    const char *nonexistent = "/tmp/nonexistent_regular_987654321.tmp";
    int result = systemd_journal_fstat_is_regular_file(nonexistent);
    assert_int_equal(result, -1);
}

/* Test: NULL path returns error */
static void test_systemd_journal_fstat_is_regular_file_null_path(void **state) {
    int result = systemd_journal_fstat_is_regular_file(NULL);
    assert_int_equal(result, -1);
}

/* Test: Empty path string returns error */
static void test_systemd_journal_fstat_is_regular_file_empty_path(void **state) {
    int result = systemd_journal_fstat_is_regular_file("");
    assert_int_equal(result, -1);
}

/* Test: Symlink to file */
static void test_systemd_journal_fstat_is_regular_file_symlink(void **state) {
    const char *test_file = "/tmp/test_journal_target_123456.tmp";
    const char *test_link = "/tmp/test_journal_link_123456.tmp";
    
    FILE *fp = fopen(test_file, "w");
    assert_non_null(fp);
    fprintf(fp, "target file");
    fclose(fp);
    
    symlink(test_file, test_link);
    
    int result = systemd_journal_fstat_is_regular_file(test_link);
    /* Result depends on implementation - could be 1 for followed symlink or 0 for symlink itself */
    assert_true(result == 0 || result == 1);
    
    unlink(test_link);
    unlink(test_file);
}

/* ============================================================================
   TEST RUNNER SETUP
   ============================================================================ */

int main(void) {
    const struct CMUnitTest tests[] = {
        /* systemd_journal_fstat_get_size tests */
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_size_success,
            test_get_size_setup, test_get_size_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_size_empty_file,
            test_get_size_setup, test_get_size_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_size_nonexistent,
            test_get_size_setup, test_get_size_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_size_null_path,
            test_get_size_setup, test_get_size_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_size_empty_path,
            test_get_size_setup, test_get_size_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_size_large_file,
            test_get_size_setup, test_get_size_teardown),
        
        /* systemd_journal_fstat_get_inode tests */
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_inode_success,
            test_get_inode_setup, test_get_inode_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_inode_null_inode_ptr,
            test_get_inode_setup, test_get_inode_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_inode_null_device_ptr,
            test_get_inode_setup, test_get_inode_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_inode_null_path,
            test_get_inode_setup, test_get_inode_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_inode_nonexistent,
            test_get_inode_setup, test_get_inode_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_get_inode_all_null_ptrs,
            test_get_inode_setup, test_get_inode_teardown),
        
        /* systemd_journal_fstat_compare_files tests */
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_identical,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_different,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_first_nonexistent,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_second_nonexistent,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_null_first_path,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_null_second_path,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_both_null_paths,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_same_file,
            test_compare_files_setup, test_compare_files_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_compare_files_empty_files,
            test_compare_files_setup, test_compare_files_teardown),
        
        /* systemd_journal_fstat_is_regular_file tests */
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_is_regular_file_regular,
            test_is_regular_file_setup, test_is_regular_file_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_is_regular_file_directory,
            test_is_regular_file_setup, test_is_regular_file_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_is_regular_file_nonexistent,
            test_is_regular_file_setup, test_is_regular_file_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_is_regular_file_null_path,
            test_is_regular_file_setup, test_is_regular_file_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_is_regular_file_empty_path,
            test_is_regular_file_setup, test_is_regular_file_teardown),
        cmocka_unit_test_setup_teardown(
            test_systemd_journal_fstat_is_regular_file_symlink,
            test_is_regular_file_setup, test_is_regular_file_teardown),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}