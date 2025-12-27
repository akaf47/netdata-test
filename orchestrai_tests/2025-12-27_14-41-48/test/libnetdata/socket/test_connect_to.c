#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <errno.h>
#include <fcntl.h>
#include <time.h>
#include <pthread.h>

/* Mock and test framework */
#include <setjmp.h>
#include <stdarg.h>
#include <cmocka.h>

/* Include the header under test */
#include "src/libnetdata/socket/connect-to.h"

/* Mock socket functions */
static int mock_socket_fd = 42;

int __wrap_socket(int domain, int type, int protocol) {
    check_expected(domain);
    check_expected(type);
    check_expected(protocol);
    return mock_socket_return_value();
}

int __wrap_connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen) {
    check_expected(sockfd);
    check_expected_ptr(addr);
    check_expected(addrlen);
    return mock_connect_return_value();
}

int __wrap_getaddrinfo(const char *node, const char *service,
                       const struct addrinfo *hints,
                       struct addrinfo **res) {
    check_expected_ptr(node);
    check_expected_ptr(service);
    check_expected_ptr(hints);
    check_expected_ptr(res);
    return mock_getaddrinfo_return_value();
}

void __wrap_freeaddrinfo(struct addrinfo *res) {
    check_expected_ptr(res);
}

int __wrap_fcntl(int fd, int cmd, ...) {
    va_list args;
    va_start(args, cmd);
    check_expected(fd);
    check_expected(cmd);
    if (cmd == F_SETFL) {
        int flags = va_arg(args, int);
        check_expected(flags);
    }
    va_end(args);
    return mock_fcntl_return_value();
}

/* Test helper functions */
static int mock_socket_return_value(void) {
    return (intptr_t)mock();
}

static int mock_connect_return_value(void) {
    return (intptr_t)mock();
}

static int mock_getaddrinfo_return_value(void) {
    return (intptr_t)mock();
}

static int mock_fcntl_return_value(void) {
    return (intptr_t)mock();
}

/* ============================================================================
 * TEST: connect_to_socket_basic_ipv4_success
 * ============================================================================
 */
static void test_connect_to_socket_basic_ipv4_success(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_basic_ipv6_success
 * ============================================================================
 */
static void test_connect_to_socket_basic_ipv6_success(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET6);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in6));
    will_return(__wrap_connect, 0);
    
    /* Act & Assert */
    int result = socket_connect("::1", 8080, AF_INET6);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_null_hostname
 * ============================================================================
 */
static void test_connect_to_socket_null_hostname(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect(NULL, 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_empty_hostname
 * ============================================================================
 */
static void test_connect_to_socket_empty_hostname(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect("", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_invalid_port_zero
 * ============================================================================
 */
static void test_connect_to_socket_invalid_port_zero(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 0, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_invalid_port_negative
 * ============================================================================
 */
static void test_connect_to_socket_invalid_port_negative(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", -1, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_invalid_port_above_max
 * ============================================================================
 */
static void test_connect_to_socket_invalid_port_above_max(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 65536, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_socket_creation_fails
 * ============================================================================
 */
static void test_connect_to_socket_socket_creation_fails(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, -1);
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_connection_refused
 * ============================================================================
 */
static void test_connect_to_socket_connection_refused(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = ECONNREFUSED;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_connection_timeout
 * ============================================================================
 */
static void test_connect_to_socket_connection_timeout(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = ETIMEDOUT;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_invalid_address_family
 * ============================================================================
 */
static void test_connect_to_socket_invalid_address_family(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_UNSPEC);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_hostname_not_found
 * ============================================================================
 */
static void test_connect_to_socket_hostname_not_found(void **state) {
    (void) state;
    
    /* Arrange */
    expect_string(__wrap_getaddrinfo, node, "nonexistent.invalid");
    expect_not_null(__wrap_getaddrinfo, service);
    expect_not_null(__wrap_getaddrinfo, hints);
    will_return(__wrap_getaddrinfo, EAI_NONAME);
    
    /* Act & Assert */
    int result = socket_connect_by_hostname("nonexistent.invalid", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_by_hostname_success
 * ============================================================================
 */
static void test_connect_to_socket_by_hostname_success(void **state) {
    (void) state;
    
    /* Arrange */
    struct addrinfo *res = malloc(sizeof(struct addrinfo));
    struct sockaddr_in *sin = malloc(sizeof(struct sockaddr_in));
    
    sin->sin_family = AF_INET;
    sin->sin_port = htons(8080);
    inet_pton(AF_INET, "127.0.0.1", &sin->sin_addr);
    
    res->ai_family = AF_INET;
    res->ai_socktype = SOCK_STREAM;
    res->ai_protocol = IPPROTO_TCP;
    res->ai_addr = (struct sockaddr *)sin;
    res->ai_addrlen = sizeof(struct sockaddr_in);
    res->ai_next = NULL;
    
    expect_string(__wrap_getaddrinfo, node, "localhost");
    expect_not_null(__wrap_getaddrinfo, service);
    expect_not_null(__wrap_getaddrinfo, hints);
    will_return(__wrap_getaddrinfo, 0);
    will_set_parameter(__wrap_getaddrinfo, res, res);
    
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    expect_value(__wrap_freeaddrinfo, res, res);
    
    /* Act & Assert */
    int result = socket_connect_by_hostname("localhost", 8080, AF_INET);
    assert_int_equal(result, mock_socket_fd);
    
    free(sin);
    free(res);
}

/* ============================================================================
 * TEST: connect_to_socket_non_blocking_success
 * ============================================================================
 */
static void test_connect_to_socket_non_blocking_success(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_fcntl, fd, mock_socket_fd);
    expect_value(__wrap_fcntl, cmd, F_SETFL);
    expect_value(__wrap_fcntl, flags, O_NONBLOCK);
    will_return(__wrap_fcntl, 0);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = EINPROGRESS;
    
    /* Act & Assert */
    int result = socket_connect_non_blocking("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_non_blocking_fcntl_fails
 * ============================================================================
 */
static void test_connect_to_socket_non_blocking_fcntl_fails(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_fcntl, fd, mock_socket_fd);
    expect_value(__wrap_fcntl, cmd, F_SETFL);
    expect_value(__wrap_fcntl, flags, O_NONBLOCK);
    will_return(__wrap_fcntl, -1);
    
    /* Act & Assert */
    int result = socket_connect_non_blocking("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_with_timeout_success
 * ============================================================================
 */
static void test_connect_to_socket_with_timeout_success(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    /* Act & Assert */
    int result = socket_connect_with_timeout("127.0.0.1", 8080, AF_INET, 5);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_with_timeout_zero
 * ============================================================================
 */
static void test_connect_to_socket_with_timeout_zero(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect_with_timeout("127.0.0.1", 8080, AF_INET, 0);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_with_timeout_negative
 * ============================================================================
 */
static void test_connect_to_socket_with_timeout_negative(void **state) {
    (void) state;
    
    /* Act & Assert */
    int result = socket_connect_with_timeout("127.0.0.1", 8080, AF_INET, -5);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_connection_in_progress_timeout
 * ============================================================================
 */
static void test_connect_to_socket_connection_in_progress_timeout(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_fcntl, fd, mock_socket_fd);
    expect_value(__wrap_fcntl, cmd, F_SETFL);
    expect_value(__wrap_fcntl, flags, O_NONBLOCK);
    will_return(__wrap_fcntl, 0);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = EINPROGRESS;
    
    /* Act & Assert */
    int result = socket_connect_with_timeout("127.0.0.1", 8080, AF_INET, 1);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_long_hostname
 * ============================================================================
 */
static void test_connect_to_socket_long_hostname(void **state) {
    (void) state;
    
    /* Arrange */
    char long_hostname[1024];
    memset(long_hostname, 'a', sizeof(long_hostname) - 1);
    long_hostname[sizeof(long_hostname) - 1] = '\0';
    
    /* Act & Assert */
    int result = socket_connect(long_hostname, 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_boundary_port_max
 * ============================================================================
 */
static void test_connect_to_socket_boundary_port_max(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 65535, AF_INET);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_boundary_port_min
 * ============================================================================
 */
static void test_connect_to_socket_boundary_port_min(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 1, AF_INET);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_network_unreachable
 * ============================================================================
 */
static void test_connect_to_socket_network_unreachable(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = ENETUNREACH;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_host_unreachable
 * ============================================================================
 */
static void test_connect_to_socket_host_unreachable(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = EHOSTUNREACH;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_permission_denied
 * ============================================================================
 */
static void test_connect_to_socket_permission_denied(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = EACCES;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_address_in_use
 * ============================================================================
 */
static void test_connect_to_socket_address_in_use(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, -1);
    errno = EADDRINUSE;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_no_memory
 * ============================================================================
 */
static void test_connect_to_socket_no_memory(void **state) {
    (void) state;
    
    /* Arrange */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, -1);
    errno = ENOMEM;
    
    /* Act & Assert */
    int result = socket_connect("127.0.0.1", 8080, AF_INET);
    assert_int_equal(result, -1);
}

/* ============================================================================
 * TEST: connect_to_socket_with_special_ipv4_addresses
 * ============================================================================
 */
static void test_connect_to_socket_with_special_ipv4_addresses(void **state) {
    (void) state;
    
    /* Test with 0.0.0.0 */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    int result = socket_connect("0.0.0.0", 8080, AF_INET);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_with_broadcast_address
 * ============================================================================
 */
static void test_connect_to_socket_with_broadcast_address(void **state) {
    (void) state;
    
    /* Test with 255.255.255.255 */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    int result = socket_connect("255.255.255.255", 8080, AF_INET);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_ipv6_all_zeros
 * ============================================================================
 */
static void test_connect_to_socket_ipv6_all_zeros(void **state) {
    (void) state;
    
    /* Test with :: */
    expect_value(__wrap_socket, domain, AF_INET6);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in6));
    will_return(__wrap_connect, 0);
    
    int result = socket_connect("::", 8080, AF_INET6);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_ipv6_full_address
 * ============================================================================
 */
static void test_connect_to_socket_ipv6_full_address(void **state) {
    (void) state;
    
    /* Test with full IPv6 address */
    expect_value(__wrap_socket, domain, AF_INET6);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in6));
    will_return(__wrap_connect, 0);
    
    int result = socket_connect("2001:0db8:85a3:0000:0000:8a2e:0370:7334", 8080, AF_INET6);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_port_1024
 * ============================================================================
 */
static void test_connect_to_socket_port_1024(void **state) {
    (void) state;
    
    /* Test boundary at 1024 (first non-privileged port) */
    expect_value(__wrap_socket, domain, AF_INET);
    expect_value(__wrap_socket, type, SOCK_STREAM);
    expect_value(__wrap_socket, protocol, IPPROTO_TCP);
    will_return(__wrap_socket, mock_socket_fd);
    
    expect_value(__wrap_connect, sockfd, mock_socket_fd);
    expect_not_null(__wrap_connect, addr);
    expect_value(__wrap_connect, addrlen, sizeof(struct sockaddr_in));
    will_return(__wrap_connect, 0);
    
    int result = socket_connect("127.0.0.1", 1024, AF_INET);
    assert_int_equal(result, mock_socket_fd);
}

/* ============================================================================
 * TEST: connect_to_socket_hostname_with_getaddrinfo_multiple_addresses
 * ============================================================================
 */
static void test_connect_to_socket_hostname_multiple_results(void **state) {
    (void) state;
    
    /* Arrange - simulate multiple