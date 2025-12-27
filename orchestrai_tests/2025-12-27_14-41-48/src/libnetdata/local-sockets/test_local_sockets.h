#ifndef TEST_LOCAL_SOCKETS_H
#define TEST_LOCAL_SOCKETS_H

#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <fcntl.h>
#include <errno.h>
#include <time.h>

/* Test framework macros */
#define TEST_ASSERT(condition, message) do { \
    if (!(condition)) { \
        fprintf(stderr, "FAIL: %s\n", message); \
        return 0; \
    } \
} while(0)

#define TEST_ASSERT_EQ(actual, expected, message) do { \
    if ((actual) != (expected)) { \
        fprintf(stderr, "FAIL: %s (expected %ld, got %ld)\n", message, (long)(expected), (long)(actual)); \
        return 0; \
    } \
} while(0)

#define TEST_ASSERT_NULL(ptr, message) do { \
    if ((ptr) != NULL) { \
        fprintf(stderr, "FAIL: %s (expected NULL)\n", message); \
        return 0; \
    } \
} while(0)

#define TEST_ASSERT_NOT_NULL(ptr, message) do { \
    if ((ptr) == NULL) { \
        fprintf(stderr, "FAIL: %s (expected non-NULL)\n", message); \
        return 0; \
    } \
} while(0)

#define TEST_ASSERT_STR_EQ(actual, expected, message) do { \
    if (strcmp((actual), (expected)) != 0) { \
        fprintf(stderr, "FAIL: %s (expected '%s', got '%s')\n", message, (expected), (actual)); \
        return 0; \
    } \
} while(0)

/* Test function prototypes for local_sockets.h coverage */

/* Test socket structure initialization */
int test_local_socket_initialization(void);

/* Test socket creation and configuration */
int test_local_socket_creation(void);
int test_local_socket_creation_invalid_family(void);
int test_local_socket_creation_invalid_type(void);

/* Test socket binding */
int test_local_socket_bind_unix(void);
int test_local_socket_bind_inet(void);
int test_local_socket_bind_invalid_address(void);
int test_local_socket_bind_already_in_use(void);

/* Test socket connection */
int test_local_socket_connect_unix(void);
int test_local_socket_connect_inet(void);
int test_local_socket_connect_nonexistent(void);
int test_local_socket_connect_timeout(void);

/* Test socket listening */
int test_local_socket_listen(void);
int test_local_socket_listen_invalid_backlog(void);

/* Test socket acceptance */
int test_local_socket_accept(void);
int test_local_socket_accept_no_connections(void);

/* Test socket reading and writing */
int test_local_socket_send_data(void);
int test_local_socket_send_empty_buffer(void);
int test_local_socket_send_large_buffer(void);
int test_local_socket_recv_data(void);
int test_local_socket_recv_empty_socket(void);
int test_local_socket_recv_timeout(void);

/* Test socket options */
int test_local_socket_set_nonblocking(void);
int test_local_socket_set_blocking(void);
int test_local_socket_set_reuse_addr(void);
int test_local_socket_set_timeout(void);

/* Test socket closing and cleanup */
int test_local_socket_close(void);
int test_local_socket_close_invalid_fd(void);
int test_local_socket_cleanup(void);

/* Test error handling */
int test_local_socket_errno_handling(void);
int test_local_socket_permission_denied(void);

/* Test edge cases */
int test_local_socket_null_pointer_handling(void);
int test_local_socket_zero_length_data(void);
int test_local_socket_negative_fd(void);
int test_local_socket_max_connections(void);

#endif /* TEST_LOCAL_SOCKETS_H */