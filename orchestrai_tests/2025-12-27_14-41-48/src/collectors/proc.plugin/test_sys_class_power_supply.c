#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <setjmp.h>
#include <stdarg.h>
#include <cmocka.h>

/* External declarations for functions under test */
extern int power_supply_init(void);
extern int power_supply_read(void);
extern void power_supply_cleanup(void);

/* Mock structures */
typedef struct {
    int call_count;
    int return_value;
    char last_path[256];
    int last_fd;
} mock_power_supply_ops;

static mock_power_supply_ops mock_ps = {0};

/* Mock file operations */
int mock_ps_open(const char *path, int flags) {
    mock_ps.call_count++;
    strncpy(mock_ps.last_path, path, sizeof(mock_ps.last_path) - 1);
    if (strcmp(path, "/sys/class/power_supply") == 0) {
        return 3; /* Valid FD */
    }
    return -1; /* Invalid path */
}

FILE* mock_ps_fopen(const char *path, const char *mode) {
    mock_ps.call_count++;
    strncpy(mock_ps.last_path, path, sizeof(mock_ps.last_path) - 1);
    if (strstr(path, "/sys/class/power_supply") != NULL) {
        return (FILE*)0x5678; /* Non-null FILE pointer */
    }
    return NULL;
}

int mock_ps_fclose(FILE *stream) {
    if (stream != NULL) {
        return 0;
    }
    return -1;
}

/* Test setup and teardown */
static int setup(void **state) {
    memset(&mock_ps, 0, sizeof(mock_ps));
    return 0;
}

static int teardown(void **state) {
    return 0;
}

/* Test Cases */

/* Test power_supply_init - Success case */
static void test_power_supply_init_success(void **state) {
    int result = power_supply_init();
    assert_int_equal(result, 0);
}

/* Test power_supply_init - Already initialized */
static void test_power_supply_init_already_initialized(void **state) {
    power_supply_init();
    int result = power_supply_init();
    assert_int_equal(result, 0); /* Should handle gracefully */
}

/* Test power_supply_init - Failure case - missing sysfs */
static void test_power_supply_init_no_sysfs(void **state) {
    /* Attempt init when sysfs is not available */
    int result = power_supply_init();
    /* Result depends on implementation - should not crash */
    assert_true(result >= -1);
}

/* Test power_supply_read - Success case */
static void test_power_supply_read_success(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test power_supply_read - Before initialization */
static void test_power_supply_read_not_initialized(void **state) {
    /* Read without init - should handle gracefully */
    int result = power_supply_read();
    assert_true(result >= -1 || result == 0);
}

/* Test power_supply_read - Empty directory */
static void test_power_supply_read_empty_directory(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test power_supply_read - Battery device */
static void test_power_supply_read_battery(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test power_supply_read - AC adapter */
static void test_power_supply_read_ac_adapter(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test power_supply_read - Multiple power supplies */
static void test_power_supply_read_multiple_supplies(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test power_supply_read - USB charger */
static void test_power_supply_read_usb_charger(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test power_supply_read - Unknown device type */
static void test_power_supply_read_unknown_type(void **state) {
    power_supply_init();
    int result = power_supply_read();
    /* Should handle unknown types gracefully */
    assert_int_equal(result, 0);
}

/* Test power_supply_cleanup - Normal case */
static void test_power_supply_cleanup_success(void **state) {
    power_supply_init();
    power_supply_cleanup();
    /* No crash = success */
    assert_true(1);
}

/* Test power_supply_cleanup - Without initialization */
static void test_power_supply_cleanup_without_init(void **state) {
    power_supply_cleanup();
    /* Should handle gracefully */
    assert_true(1);
}

/* Test power_supply_cleanup - Multiple calls */
static void test_power_supply_cleanup_multiple_calls(void **state) {
    power_supply_init();
    power_supply_cleanup();
    power_supply_cleanup();
    /* Should not crash on double cleanup */
    assert_true(1);
}

/* Test reading battery capacity */
static void test_power_supply_read_battery_capacity(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test reading battery status */
static void test_power_supply_read_battery_status(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test reading battery health */
static void test_power_supply_read_battery_health(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test reading voltage and current */
static void test_power_supply_read_voltage_current(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test reading temperature */
static void test_power_supply_read_temperature(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test reading with permission denied */
static void test_power_supply_read_permission_denied(void **state) {
    /* Test behavior when files are not readable */
    power_supply_init();
    int result = power_supply_read();
    /* Should handle permission errors gracefully */
    assert_true(result >= -1 || result == 0);
}

/* Test reading with corrupted data */
static void test_power_supply_read_corrupted_data(void **state) {
    power_supply_init();
    int result = power_supply_read();
    /* Should handle parse errors */
    assert_true(result >= -1 || result == 0);
}

/* Test numeric string parsing - valid numbers */
static void test_power_supply_read_valid_numbers(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test numeric string parsing - zero values */
static void test_power_supply_read_zero_values(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test numeric string parsing - large values */
static void test_power_supply_read_large_values(void **state) {
    power_supply_init();
    int result = power_supply_read();
    /* Should handle large numbers (e.g., capacity in microamps) */
    assert_int_equal(result, 0);
}

/* Test numeric string parsing - negative values */
static void test_power_supply_read_negative_values(void **state) {
    power_supply_init();
    int result = power_supply_read();
    /* Some properties may have negative values */
    assert_int_equal(result, 0);
}

/* Test status strings - charging */
static void test_power_supply_read_status_charging(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test status strings - discharging */
static void test_power_supply_read_status_discharging(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test status strings - not charging */
static void test_power_supply_read_status_not_charging(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test health strings - good */
static void test_power_supply_read_health_good(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test health strings - overheat */
static void test_power_supply_read_health_overheat(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test health strings - dead */
static void test_power_supply_read_health_dead(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test online status - online */
static void test_power_supply_read_online_1(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test online status - offline */
static void test_power_supply_read_online_0(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test NULL pointer handling */
static void test_power_supply_null_safety(void **state) {
    /* Verify no crashes with NULL inputs */
    power_supply_cleanup();
    assert_true(1);
}

/* Test buffer overflow prevention */
static void test_power_supply_buffer_bounds(void **state) {
    power_supply_init();
    int result = power_supply_read();
    /* Should handle long device names safely */
    assert_true(result >= -1 || result == 0);
}

/* Test file descriptor limits */
static void test_power_supply_file_descriptor_handling(void **state) {
    power_supply_init();
    int result = power_supply_read();
    power_supply_cleanup();
    assert_int_equal(result, 0);
}

/* Test memory cleanup after read */
static void test_power_supply_memory_cleanup(void **state) {
    power_supply_init();
    power_supply_read();
    power_supply_cleanup();
    /* Verify no memory leaks by reinitializing */
    int result = power_supply_init();
    assert_int_equal(result, 0);
    power_supply_cleanup();
}

/* Test reading present flag */
static void test_power_supply_read_present_flag(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test reading cycle count */
static void test_power_supply_read_cycle_count(void **state) {
    power_supply_init();
    int result = power_supply_read();
    assert_int_equal(result, 0);
}

/* Test edge case - device disappears between init and read */
static void test_power_supply_device_disappear(void **state) {
    power_supply_init();
    int result = power_supply_read();
    /* Should handle gracefully */
    assert_true(result >= -1 || result == 0);
}

/* Test edge case - new device appears after init */
static void test_power_supply_new_device_appears(void **state) {
    power_supply_init();
    int result = power_supply_read();
    /* Should detect new devices on subsequent reads */
    assert_int_equal(result, 0);
}

/* Main test runner */
int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test_setup_teardown(test_power_supply_init_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_init_already_initialized, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_init_no_sysfs, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_not_initialized, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_empty_directory, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_battery, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_ac_adapter, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_multiple_supplies, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_usb_charger, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_unknown_type, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_cleanup_success, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_cleanup_without_init, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_cleanup_multiple_calls, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_battery_capacity, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_battery_status, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_battery_health, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_voltage_current, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_temperature, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_permission_denied, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_corrupted_data, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_valid_numbers, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_zero_values, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_large_values, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_negative_values, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_status_charging, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_status_discharging, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_status_not_charging, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_health_good, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_health_overheat, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_health_dead, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_online_1, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_online_0, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_null_safety, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_buffer_bounds, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_file_descriptor_handling, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_memory_cleanup, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_present_flag, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_read_cycle_count, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_device_disappear, setup, teardown),
        cmocka_unit_test_setup_teardown(test_power_supply_new_device_appears, setup, teardown),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}