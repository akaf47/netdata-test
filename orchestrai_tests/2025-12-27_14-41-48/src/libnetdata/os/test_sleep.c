#include <stdarg.h>
#include <stddef.h>
#include <setjmp.h>
#include <cmocka.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <errno.h>
#include <signal.h>
#include <unistd.h>

// Mock declarations
int __wrap_nanosleep(const struct timespec *req, struct timespec *rem);
unsigned int __wrap_sleep(unsigned int seconds);
int __wrap_usleep(useconds_t usec);
int __wrap_clock_nanosleep(clockid_t clockid, int flags,
                           const struct timespec *request,
                           struct timespec *remain);

// Forward declaration of sleep functions to be tested
void netdata_sleep(unsigned long long milliseconds);
int netdata_nanosleep(unsigned long long nanoseconds);
void netdata_usleep(unsigned long long microseconds);

// Test: netdata_sleep with 0 milliseconds
static void test_netdata_sleep_zero_milliseconds(void **state) {
    (void) state;
    
    // Zero sleep should return immediately
    will_return(__wrap_nanosleep, 0);
    netdata_sleep(0);
}

// Test: netdata_sleep with 1 millisecond
static void test_netdata_sleep_one_millisecond(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    netdata_sleep(1);
}

// Test: netdata_sleep with 1000 milliseconds (1 second)
static void test_netdata_sleep_one_second(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    netdata_sleep(1000);
}

// Test: netdata_sleep with large millisecond value
static void test_netdata_sleep_large_value(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    netdata_sleep(86400000); // 24 hours in milliseconds
}

// Test: netdata_sleep with max millisecond value
static void test_netdata_sleep_max_value(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    netdata_sleep(0xFFFFFFFFFFFFFFFFULL); // Max uint64_t
}

// Test: netdata_sleep when nanosleep returns EINTR
static void test_netdata_sleep_interrupted(void **state) {
    (void) state;
    
    // First call interrupted
    will_return(__wrap_nanosleep, -1);
    errno = EINTR;
    
    // Retry after interruption
    will_return(__wrap_nanosleep, 0);
    
    netdata_sleep(1000);
}

// Test: netdata_sleep with multiple interruptions
static void test_netdata_sleep_multiple_interruptions(void **state) {
    (void) state;
    
    // First interruption
    will_return(__wrap_nanosleep, -1);
    errno = EINTR;
    
    // Second interruption
    will_return(__wrap_nanosleep, -1);
    errno = EINTR;
    
    // Third interruption
    will_return(__wrap_nanosleep, -1);
    errno = EINTR;
    
    // Finally succeeds
    will_return(__wrap_nanosleep, 0);
    
    netdata_sleep(1000);
}

// Test: netdata_nanosleep with 0 nanoseconds
static void test_netdata_nanosleep_zero(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    int result = netdata_nanosleep(0);
    assert_int_equal(result, 0);
}

// Test: netdata_nanosleep with 1 nanosecond
static void test_netdata_nanosleep_one_nanosecond(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    int result = netdata_nanosleep(1);
    assert_int_equal(result, 0);
}

// Test: netdata_nanosleep with 1 second worth of nanoseconds
static void test_netdata_nanosleep_one_second_ns(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    int result = netdata_nanosleep(1000000000);
    assert_int_equal(result, 0);
}

// Test: netdata_nanosleep with large nanosecond value
static void test_netdata_nanosleep_large_value(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, 0);
    int result = netdata_nanosleep(86400000000000); // 24 hours in nanoseconds
    assert_int_equal(result, 0);
}

// Test: netdata_nanosleep when interrupted
static void test_netdata_nanosleep_interrupted(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, -1);
    errno = EINTR;
    will_return(__wrap_nanosleep, 0);
    
    int result = netdata_nanosleep(1000000);
    assert_int_equal(result, 0);
}

// Test: netdata_nanosleep with EINVAL error
static void test_netdata_nanosleep_einval(void **state) {
    (void) state;
    
    will_return(__wrap_nanosleep, -1);
    errno = EINVAL;
    
    int result = netdata_nanosleep(1000000);
    assert_int_equal(result, -1);
}

// Test: netdata_usleep with 0 microseconds
static void test_netdata_usleep_zero(void **state) {
    (void) state;
    
    will_return(__wrap_usleep, 0);
    netdata_usleep(0);
}

// Test: netdata_usleep with 1 microsecond
static void test_netdata_usleep_one_microsecond(void **state) {
    (void) state;
    
    will_return(__wrap_usleep, 0);
    netdata_usleep(1);
}

// Test: netdata_usleep with 1 millisecond (1000 microseconds)
static void test_netdata_usleep_one_millisecond(void **state) {
    (void) state;
    
    will_return(__wrap_usleep, 0);
    netdata_usleep(1000);
}

// Test: netdata_usleep with 1 second (1000000 microseconds)
static void test_netdata_usleep_one_second(void **state) {
    (void) state;
    
    will_return(__wrap_usleep, 0);
    netdata_usleep(1000000);
}

// Test: netdata_usleep with large microsecond value
static void test_netdata_usleep_large_value(void **state) {
    (void) state;
    
    will_return(__wrap_usleep, 0);
    netdata_usleep(86400000000); // 24 hours in microseconds
}

// Test: netdata_usleep when interrupted
static void test_netdata_usleep_interrupted(void **state) {
    (void) state;
    
    will_return(__wrap_usleep, 1); // EINTR
    will_return(__wrap_usleep, 0);
    
    netdata_usleep(1000);
}

// Test: netdata_usleep with multiple interruptions
static void test_netdata_usleep_multiple_interruptions(void **state) {
    (void) state;
    
    will_return(__wrap_usleep, 1); // EINTR
    will_return(__wrap_usleep, 1); // EINTR
    will_return(__wrap_usleep, 0); // Success
    
    netdata_usleep(1000);
}

// Mock implementations
int __wrap_nanosleep(const struct timespec *req, struct timespec *rem) {
    (void) req;
    (void) rem;
    return mock_type(int);
}

unsigned int __wrap_sleep(unsigned int seconds) {
    (void) seconds;
    return mock_type(unsigned int);
}

int __wrap_usleep(useconds_t usec) {
    (void) usec;
    return mock_type(int);
}

int __wrap_clock_nanosleep(clockid_t clockid, int flags,
                           const struct timespec *request,
                           struct timespec *remain) {
    (void) clockid;
    (void) flags;
    (void) request;
    (void) remain;
    return mock_type(int);
}

// Test suite
int main(void) {
    const struct CMUnitTest tests[] = {
        cmocka_unit_test(test_netdata_sleep_zero_milliseconds),
        cmocka_unit_test(test_netdata_sleep_one_millisecond),
        cmocka_unit_test(test_netdata_sleep_one_second),
        cmocka_unit_test(test_netdata_sleep_large_value),
        cmocka_unit_test(test_netdata_sleep_max_value),
        cmocka_unit_test(test_netdata_sleep_interrupted),
        cmocka_unit_test(test_netdata_sleep_multiple_interruptions),
        cmocka_unit_test(test_netdata_nanosleep_zero),
        cmocka_unit_test(test_netdata_nanosleep_one_nanosecond),
        cmocka_unit_test(test_netdata_nanosleep_one_second_ns),
        cmocka_unit_test(test_netdata_nanosleep_large_value),
        cmocka_unit_test(test_netdata_nanosleep_interrupted),
        cmocka_unit_test(test_netdata_nanosleep_einval),
        cmocka_unit_test(test_netdata_usleep_zero),
        cmocka_unit_test(test_netdata_usleep_one_microsecond),
        cmocka_unit_test(test_netdata_usleep_one_millisecond),
        cmocka_unit_test(test_netdata_usleep_one_second),
        cmocka_unit_test(test_netdata_usleep_large_value),
        cmocka_unit_test(test_netdata_usleep_interrupted),
        cmocka_unit_test(test_netdata_usleep_multiple_interruptions),
    };

    return cmocka_run_group_tests(tests, NULL, NULL);
}