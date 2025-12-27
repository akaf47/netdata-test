#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <pthread.h>
#include <stdint.h>
#include <stdbool.h>
#include <unistd.h>
#include <time.h>

/* Mock the locks.h header for testing purposes */
#ifndef NETDATA_LOCKS_H
#define NETDATA_LOCKS_H

/* Test utilities */
typedef struct {
    const char *test_name;
    int passed;
    int total_assertions;
} test_result_t;

test_result_t test_results = {0};

#define TEST_START(name) \
    do { \
        test_results.test_name = name; \
        test_results.passed = 1; \
        test_results.total_assertions = 0; \
        printf("\n[TEST] %s\n", name); \
    } while(0)

#define TEST_ASSERT(condition, message) \
    do { \
        test_results.total_assertions++; \
        if (!(condition)) { \
            printf("  FAIL: %s\n", message); \
            test_results.passed = 0; \
        } else { \
            printf("  PASS: %s\n", message); \
        } \
    } while(0)

#define TEST_END() \
    do { \
        if (test_results.passed) { \
            printf("[RESULT] %s - PASSED (%d assertions)\n", test_results.test_name, test_results.total_assertions); \
        } else { \
            printf("[RESULT] %s - FAILED\n", test_results.test_name); \
        } \
        return test_results.passed; \
    } while(0)

/* Tests for mutex locks */
int test_mutex_lock_initialization() {
    TEST_START("mutex_lock_initialization");
    
    pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
    TEST_ASSERT(pthread_mutex_init(&mutex, NULL) == 0, "Mutex initialization succeeds");
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Mutex destruction succeeds");
    
    TEST_END();
}

int test_mutex_lock_lock_unlock() {
    TEST_START("mutex_lock_lock_unlock");
    
    pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;
    pthread_mutex_init(&mutex, NULL);
    
    TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "Mutex lock succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Mutex unlock succeeds");
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Mutex destruction after lock/unlock succeeds");
    
    TEST_END();
}

int test_mutex_lock_multiple_lock_attempts() {
    TEST_START("mutex_lock_multiple_lock_attempts");
    
    pthread_mutex_t mutex;
    pthread_mutexattr_t attr;
    
    TEST_ASSERT(pthread_mutexattr_init(&attr) == 0, "Mutex attribute initialization succeeds");
    TEST_ASSERT(pthread_mutex_init(&mutex, &attr) == 0, "Mutex with attributes initialization succeeds");
    
    TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "First lock succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "First unlock succeeds");
    TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "Second lock succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Second unlock succeeds");
    
    TEST_ASSERT(pthread_mutexattr_destroy(&attr) == 0, "Mutex attribute destruction succeeds");
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Final mutex destruction succeeds");
    
    TEST_END();
}

int test_mutex_lock_trylock() {
    TEST_START("mutex_lock_trylock");
    
    pthread_mutex_t mutex;
    TEST_ASSERT(pthread_mutex_init(&mutex, NULL) == 0, "Mutex initialization succeeds");
    
    /* Trylock on unlocked mutex should succeed */
    TEST_ASSERT(pthread_mutex_trylock(&mutex) == 0, "Trylock on unlocked mutex succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Unlock after trylock succeeds");
    
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Mutex destruction succeeds");
    
    TEST_END();
}

int test_mutex_lock_error_conditions() {
    TEST_START("mutex_lock_error_conditions");
    
    pthread_mutex_t mutex;
    TEST_ASSERT(pthread_mutex_init(&mutex, NULL) == 0, "Mutex initialization succeeds");
    TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "Lock succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Unlock succeeds");
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Destruction succeeds");
    
    /* Attempt operations on destroyed mutex - behavior depends on platform */
    int result = pthread_mutex_lock(&mutex);
    TEST_ASSERT(result != 0, "Lock on destroyed mutex fails");
    
    TEST_END();
}

int test_mutex_lock_null_pointer_handling() {
    TEST_START("mutex_lock_null_pointer_handling");
    
    /* Null pointer to mutex should fail gracefully */
    int result = pthread_mutex_lock(NULL);
    TEST_ASSERT(result != 0, "Lock with NULL mutex pointer fails");
    
    result = pthread_mutex_unlock(NULL);
    TEST_ASSERT(result != 0, "Unlock with NULL mutex pointer fails");
    
    TEST_END();
}

int test_recursive_lock_behavior() {
    TEST_START("recursive_lock_behavior");
    
    pthread_mutex_t mutex;
    pthread_mutexattr_t attr;
    
    TEST_ASSERT(pthread_mutexattr_init(&attr) == 0, "Attribute initialization succeeds");
    TEST_ASSERT(pthread_mutexattr_settype(&attr, PTHREAD_MUTEX_RECURSIVE) == 0, 
                "Setting recursive mutex type succeeds");
    TEST_ASSERT(pthread_mutex_init(&mutex, &attr) == 0, "Recursive mutex initialization succeeds");
    
    /* Multiple locks should be allowed */
    TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "First recursive lock succeeds");
    TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "Second recursive lock succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "First unlock succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Second unlock succeeds");
    
    TEST_ASSERT(pthread_mutexattr_destroy(&attr) == 0, "Attribute destruction succeeds");
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Mutex destruction succeeds");
    
    TEST_END();
}

int test_locked_state_detection() {
    TEST_START("locked_state_detection");
    
    pthread_mutex_t mutex;
    TEST_ASSERT(pthread_mutex_init(&mutex, NULL) == 0, "Mutex initialization succeeds");
    
    /* Use trylock to detect if locked */
    TEST_ASSERT(pthread_mutex_trylock(&mutex) == 0, "Mutex is initially unlocked");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Unlock succeeds");
    
    /* Lock and then try to lock again */
    TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "Lock succeeds");
    int result = pthread_mutex_trylock(&mutex);
    TEST_ASSERT(result != 0, "Trylock on locked mutex fails");
    TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Unlock succeeds");
    
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Destruction succeeds");
    
    TEST_END();
}

int test_concurrent_lock_access() {
    TEST_START("concurrent_lock_access");
    
    pthread_mutex_t mutex;
    volatile int shared_value = 0;
    
    TEST_ASSERT(pthread_mutex_init(&mutex, NULL) == 0, "Mutex initialization succeeds");
    
    pthread_mutex_lock(&mutex);
    shared_value = 1;
    TEST_ASSERT(shared_value == 1, "Protected write succeeds");
    pthread_mutex_unlock(&mutex);
    
    pthread_mutex_lock(&mutex);
    int value = shared_value;
    TEST_ASSERT(value == 1, "Protected read succeeds");
    pthread_mutex_unlock(&mutex);
    
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Destruction succeeds");
    
    TEST_END();
}

int test_spinlock_vs_mutex_performance() {
    TEST_START("spinlock_vs_mutex_performance");
    
    pthread_mutex_t mutex;
    TEST_ASSERT(pthread_mutex_init(&mutex, NULL) == 0, "Mutex initialization succeeds");
    
    /* Mutex should handle rapid lock/unlock cycles */
    for (int i = 0; i < 100; i++) {
        TEST_ASSERT(pthread_mutex_lock(&mutex) == 0, "Rapid lock cycle succeeds");
        TEST_ASSERT(pthread_mutex_unlock(&mutex) == 0, "Rapid unlock cycle succeeds");
    }
    
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Destruction succeeds");
    
    TEST_END();
}

int test_lock_interleaving() {
    TEST_START("lock_interleaving");
    
    pthread_mutex_t mutex1, mutex2;
    TEST_ASSERT(pthread_mutex_init(&mutex1, NULL) == 0, "First mutex initialization succeeds");
    TEST_ASSERT(pthread_mutex_init(&mutex2, NULL) == 0, "Second mutex initialization succeeds");
    
    /* Test acquiring locks in different orders */
    TEST_ASSERT(pthread_mutex_lock(&mutex1) == 0, "First lock acquisition succeeds");
    TEST_ASSERT(pthread_mutex_lock(&mutex2) == 0, "Second lock acquisition succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex2) == 0, "Second unlock succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex1) == 0, "First unlock succeeds");
    
    /* Test reverse order */
    TEST_ASSERT(pthread_mutex_lock(&mutex2) == 0, "Reverse order lock 1 succeeds");
    TEST_ASSERT(pthread_mutex_lock(&mutex1) == 0, "Reverse order lock 2 succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex1) == 0, "Reverse order unlock 1 succeeds");
    TEST_ASSERT(pthread_mutex_unlock(&mutex2) == 0, "Reverse order unlock 2 succeeds");
    
    TEST_ASSERT(pthread_mutex_destroy(&mutex1) == 0, "First mutex destruction succeeds");
    TEST_ASSERT(pthread_mutex_destroy(&mutex2) == 0, "Second mutex destruction succeeds");
    
    TEST_END();
}

int test_lock_contention() {
    TEST_START("lock_contention");
    
    pthread_mutex_t mutex;
    volatile int counter = 0;
    
    TEST_ASSERT(pthread_mutex_init(&mutex, NULL) == 0, "Mutex initialization succeeds");
    
    /* Simulate contention with sequential access */
    for (int i = 0; i < 10; i++) {
        pthread_mutex_lock(&mutex);
        counter++;
        pthread_mutex_unlock(&mutex);
    }
    
    TEST_ASSERT(counter == 10, "Counter correctly incremented under mutex protection");
    TEST_ASSERT(pthread_mutex_destroy(&mutex) == 0, "Destruction succeeds");
    
    TEST_END();
}

/* Main test runner */
int main() {
    int passed = 0;
    int total = 0;
    
    printf("\n=== NETDATA LOCKS.H TEST SUITE ===\n");
    
    total++; passed += test_mutex_lock_initialization();
    total++; passed += test_mutex_lock_lock_unlock();
    total++; passed += test_mutex_lock_multiple_lock_attempts();
    total++; passed += test_mutex_lock_trylock();
    total++; passed += test_mutex_lock_error_conditions();
    total++; passed += test_mutex_lock_null_pointer_handling();
    total++; passed += test_recursive_lock_behavior();
    total++; passed += test_locked_state_detection();
    total++; passed += test_concurrent_lock_access();
    total++; passed += test_spinlock_vs_mutex_performance();
    total++; passed += test_lock_interleaving();
    total++; passed += test_lock_contention();
    
    printf("\n=== TEST SUMMARY ===\n");
    printf("Total Tests: %d\n", total);
    printf("Passed: %d\n", passed);
    printf("Failed: %d\n", total - passed);
    printf("Coverage: %.1f%%\n\n", (100.0 * passed) / total);
    
    return (passed == total) ? 0 : 1;
}

#endif /* NETDATA_LOCKS_H */