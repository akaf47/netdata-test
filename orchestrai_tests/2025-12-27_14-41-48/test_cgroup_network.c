#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <dirent.h>
#include <time.h>
#include <assert.h>

/* Mock structures and functions */
typedef struct {
    char *name;
    unsigned long long bytes_in;
    unsigned long long bytes_out;
    unsigned long long packets_in;
    unsigned long long packets_out;
    unsigned long long dropped_in;
    unsigned long long dropped_out;
} cgroup_network_interface;

typedef struct {
    char *cgroup_name;
    cgroup_network_interface *interfaces;
    int interface_count;
    unsigned long long total_bytes_in;
    unsigned long long total_bytes_out;
} cgroup_network_data;

/* Test data structures */
typedef struct {
    int call_count;
    int return_value;
    char *last_called_with;
} mock_call_tracker;

/* Global mock tracker */
mock_call_tracker mock_tracker = {0, 0, NULL};

/* Mock function implementations */
int mock_read_cgroup_file(const char *path, char *buffer, size_t bufsize) {
    mock_tracker.call_count++;
    mock_tracker.last_called_with = (char *)path;
    
    if (strstr(path, "error") != NULL) {
        return -1;
    }
    
    if (strstr(path, "empty") != NULL) {
        strcpy(buffer, "");
        return 0;
    }
    
    if (strstr(path, "network") != NULL) {
        strcpy(buffer, "eth0: 1000 2000\neth1: 3000 4000\n");
        return strlen(buffer);
    }
    
    strcpy(buffer, "1000 2000");
    return strlen(buffer);
}

int mock_parse_network_stats(const char *stats_str, cgroup_network_interface *iface) {
    if (!stats_str || !iface) {
        return -1;
    }
    
    if (strlen(stats_str) == 0) {
        return -1;
    }
    
    sscanf(stats_str, "%llu %llu", &iface->bytes_in, &iface->bytes_out);
    return 0;
}

cgroup_network_data *create_network_data(const char *cgroup_name) {
    cgroup_network_data *data = malloc(sizeof(cgroup_network_data));
    if (!data) {
        return NULL;
    }
    
    data->cgroup_name = malloc(strlen(cgroup_name) + 1);
    if (!data->cgroup_name) {
        free(data);
        return NULL;
    }
    
    strcpy(data->cgroup_name, cgroup_name);
    data->interfaces = NULL;
    data->interface_count = 0;
    data->total_bytes_in = 0;
    data->total_bytes_out = 0;
    
    return data;
}

void free_network_data(cgroup_network_data *data) {
    if (!data) {
        return;
    }
    
    if (data->cgroup_name) {
        free(data->cgroup_name);
    }
    
    if (data->interfaces) {
        for (int i = 0; i < data->interface_count; i++) {
            if (data->interfaces[i].name) {
                free(data->interfaces[i].name);
            }
        }
        free(data->interfaces);
    }
    
    free(data);
}

int add_network_interface(cgroup_network_data *data, const char *ifname, 
                         unsigned long long bytes_in, unsigned long long bytes_out) {
    if (!data || !ifname) {
        return -1;
    }
    
    cgroup_network_interface *new_interfaces = realloc(data->interfaces,
                                                       (data->interface_count + 1) * sizeof(cgroup_network_interface));
    if (!new_interfaces) {
        return -1;
    }
    
    data->interfaces = new_interfaces;
    data->interfaces[data->interface_count].name = malloc(strlen(ifname) + 1);
    
    if (!data->interfaces[data->interface_count].name) {
        return -1;
    }
    
    strcpy(data->interfaces[data->interface_count].name, ifname);
    data->interfaces[data->interface_count].bytes_in = bytes_in;
    data->interfaces[data->interface_count].bytes_out = bytes_out;
    data->interfaces[data->interface_count].packets_in = 0;
    data->interfaces[data->interface_count].packets_out = 0;
    data->interfaces[data->interface_count].dropped_in = 0;
    data->interfaces[data->interface_count].dropped_out = 0;
    
    data->total_bytes_in += bytes_in;
    data->total_bytes_out += bytes_out;
    data->interface_count++;
    
    return 0;
}

int read_cgroup_network_stats(const char *cgroup_path, cgroup_network_data *data) {
    if (!cgroup_path || !data) {
        return -1;
    }
    
    char buffer[4096];
    int ret = mock_read_cgroup_file(cgroup_path, buffer, sizeof(buffer));
    
    if (ret < 0) {
        return -1;
    }
    
    if (ret == 0 || strlen(buffer) == 0) {
        return 0;
    }
    
    return 0;
}

unsigned long long get_total_bytes(cgroup_network_data *data) {
    if (!data) {
        return 0;
    }
    
    return data->total_bytes_in + data->total_bytes_out;
}

int update_network_interface(cgroup_network_data *data, const char *ifname,
                            unsigned long long new_bytes_in, unsigned long long new_bytes_out) {
    if (!data || !ifname) {
        return -1;
    }
    
    for (int i = 0; i < data->interface_count; i++) {
        if (strcmp(data->interfaces[i].name, ifname) == 0) {
            data->total_bytes_in -= data->interfaces[i].bytes_in;
            data->total_bytes_out -= data->interfaces[i].bytes_out;
            
            data->interfaces[i].bytes_in = new_bytes_in;
            data->interfaces[i].bytes_out = new_bytes_out;
            
            data->total_bytes_in += new_bytes_in;
            data->total_bytes_out += new_bytes_out;
            
            return 0;
        }
    }
    
    return -1;
}

cgroup_network_interface *get_interface(cgroup_network_data *data, const char *ifname) {
    if (!data || !ifname) {
        return NULL;
    }
    
    for (int i = 0; i < data->interface_count; i++) {
        if (strcmp(data->interfaces[i].name, ifname) == 0) {
            return &data->interfaces[i];
        }
    }
    
    return NULL;
}

/* Test suite */
void test_create_network_data_success() {
    cgroup_network_data *data = create_network_data("docker_container");
    assert(data != NULL);
    assert(data->cgroup_name != NULL);
    assert(strcmp(data->cgroup_name, "docker_container") == 0);
    assert(data->interface_count == 0);
    assert(data->total_bytes_in == 0);
    assert(data->total_bytes_out == 0);
    assert(data->interfaces == NULL);
    free_network_data(data);
    printf("✓ test_create_network_data_success\n");
}

void test_create_network_data_empty_name() {
    cgroup_network_data *data = create_network_data("");
    assert(data != NULL);
    assert(strlen(data->cgroup_name) == 0);
    free_network_data(data);
    printf("✓ test_create_network_data_empty_name\n");
}

void test_free_network_data_null() {
    free_network_data(NULL);
    printf("✓ test_free_network_data_null\n");
}

void test_free_network_data_with_interfaces() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    add_network_interface(data, "eth1", 3000, 4000);
    free_network_data(data);
    printf("✓ test_free_network_data_with_interfaces\n");
}

void test_add_network_interface_success() {
    cgroup_network_data *data = create_network_data("test");
    int ret = add_network_interface(data, "eth0", 1000, 2000);
    assert(ret == 0);
    assert(data->interface_count == 1);
    assert(strcmp(data->interfaces[0].name, "eth0") == 0);
    assert(data->interfaces[0].bytes_in == 1000);
    assert(data->interfaces[0].bytes_out == 2000);
    assert(data->total_bytes_in == 1000);
    assert(data->total_bytes_out == 2000);
    free_network_data(data);
    printf("✓ test_add_network_interface_success\n");
}

void test_add_network_interface_multiple() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    add_network_interface(data, "eth1", 3000, 4000);
    add_network_interface(data, "lo", 100, 100);
    
    assert(data->interface_count == 3);
    assert(data->total_bytes_in == 4100);
    assert(data->total_bytes_out == 6100);
    free_network_data(data);
    printf("✓ test_add_network_interface_multiple\n");
}

void test_add_network_interface_null_data() {
    int ret = add_network_interface(NULL, "eth0", 1000, 2000);
    assert(ret == -1);
    printf("✓ test_add_network_interface_null_data\n");
}

void test_add_network_interface_null_ifname() {
    cgroup_network_data *data = create_network_data("test");
    int ret = add_network_interface(data, NULL, 1000, 2000);
    assert(ret == -1);
    free_network_data(data);
    printf("✓ test_add_network_interface_null_ifname\n");
}

void test_add_network_interface_empty_ifname() {
    cgroup_network_data *data = create_network_data("test");
    int ret = add_network_interface(data, "", 1000, 2000);
    assert(ret == 0);
    assert(data->interface_count == 1);
    free_network_data(data);
    printf("✓ test_add_network_interface_empty_ifname\n");
}

void test_add_network_interface_zero_bytes() {
    cgroup_network_data *data = create_network_data("test");
    int ret = add_network_interface(data, "eth0", 0, 0);
    assert(ret == 0);
    assert(data->interfaces[0].bytes_in == 0);
    assert(data->interfaces[0].bytes_out == 0);
    free_network_data(data);
    printf("✓ test_add_network_interface_zero_bytes\n");
}

void test_add_network_interface_large_values() {
    cgroup_network_data *data = create_network_data("test");
    unsigned long long large_val = 18446744073709551615ULL; // Max uint64
    int ret = add_network_interface(data, "eth0", large_val, large_val);
    assert(ret == 0);
    assert(data->interfaces[0].bytes_in == large_val);
    free_network_data(data);
    printf("✓ test_add_network_interface_large_values\n");
}

void test_read_cgroup_network_stats_null_path() {
    cgroup_network_data *data = create_network_data("test");
    int ret = read_cgroup_network_stats(NULL, data);
    assert(ret == -1);
    free_network_data(data);
    printf("✓ test_read_cgroup_network_stats_null_path\n");
}

void test_read_cgroup_network_stats_null_data() {
    int ret = read_cgroup_network_stats("/path/to/cgroup", NULL);
    assert(ret == -1);
    printf("✓ test_read_cgroup_network_stats_null_data\n");
}

void test_read_cgroup_network_stats_both_null() {
    int ret = read_cgroup_network_stats(NULL, NULL);
    assert(ret == -1);
    printf("✓ test_read_cgroup_network_stats_both_null\n");
}

void test_read_cgroup_network_stats_file_error() {
    cgroup_network_data *data = create_network_data("test");
    int ret = read_cgroup_network_stats("/path/with/error", data);
    assert(ret == -1);
    free_network_data(data);
    printf("✓ test_read_cgroup_network_stats_file_error\n");
}

void test_read_cgroup_network_stats_empty_file() {
    cgroup_network_data *data = create_network_data("test");
    int ret = read_cgroup_network_stats("/path/with/empty", data);
    assert(ret == 0);
    free_network_data(data);
    printf("✓ test_read_cgroup_network_stats_empty_file\n");
}

void test_get_total_bytes_null() {
    unsigned long long total = get_total_bytes(NULL);
    assert(total == 0);
    printf("✓ test_get_total_bytes_null\n");
}

void test_get_total_bytes_empty() {
    cgroup_network_data *data = create_network_data("test");
    unsigned long long total = get_total_bytes(data);
    assert(total == 0);
    free_network_data(data);
    printf("✓ test_get_total_bytes_empty\n");
}

void test_get_total_bytes_with_data() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    add_network_interface(data, "eth1", 3000, 4000);
    unsigned long long total = get_total_bytes(data);
    assert(total == 10000);
    free_network_data(data);
    printf("✓ test_get_total_bytes_with_data\n");
}

void test_update_network_interface_success() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    int ret = update_network_interface(data, "eth0", 5000, 6000);
    assert(ret == 0);
    assert(data->interfaces[0].bytes_in == 5000);
    assert(data->interfaces[0].bytes_out == 6000);
    assert(data->total_bytes_in == 5000);
    assert(data->total_bytes_out == 6000);
    free_network_data(data);
    printf("✓ test_update_network_interface_success\n");
}

void test_update_network_interface_null_data() {
    int ret = update_network_interface(NULL, "eth0", 5000, 6000);
    assert(ret == -1);
    printf("✓ test_update_network_interface_null_data\n");
}

void test_update_network_interface_null_ifname() {
    cgroup_network_data *data = create_network_data("test");
    int ret = update_network_interface(data, NULL, 5000, 6000);
    assert(ret == -1);
    free_network_data(data);
    printf("✓ test_update_network_interface_null_ifname\n");
}

void test_update_network_interface_not_found() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    int ret = update_network_interface(data, "eth1", 5000, 6000);
    assert(ret == -1);
    assert(data->interfaces[0].bytes_in == 1000);
    free_network_data(data);
    printf("✓ test_update_network_interface_not_found\n");
}

void test_update_network_interface_zero_values() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    int ret = update_network_interface(data, "eth0", 0, 0);
    assert(ret == 0);
    assert(data->total_bytes_in == 0);
    assert(data->total_bytes_out == 0);
    free_network_data(data);
    printf("✓ test_update_network_interface_zero_values\n");
}

void test_update_network_interface_multiple_interfaces() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    add_network_interface(data, "eth1", 3000, 4000);
    add_network_interface(data, "lo", 100, 100);
    
    int ret = update_network_interface(data, "eth1", 5000, 6000);
    assert(ret == 0);
    assert(data->total_bytes_in == 6100);
    assert(data->total_bytes_out == 8100);
    free_network_data(data);
    printf("✓ test_update_network_interface_multiple_interfaces\n");
}

void test_get_interface_success() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    
    cgroup_network_interface *iface = get_interface(data, "eth0");
    assert(iface != NULL);
    assert(strcmp(iface->name, "eth0") == 0);
    assert(iface->bytes_in == 1000);
    assert(iface->bytes_out == 2000);
    free_network_data(data);
    printf("✓ test_get_interface_success\n");
}

void test_get_interface_null_data() {
    cgroup_network_interface *iface = get_interface(NULL, "eth0");
    assert(iface == NULL);
    printf("✓ test_get_interface_null_data\n");
}

void test_get_interface_null_ifname() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    
    cgroup_network_interface *iface = get_interface(data, NULL);
    assert(iface == NULL);
    free_network_data(data);
    printf("✓ test_get_interface_null_ifname\n");
}

void test_get_interface_not_found() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    
    cgroup_network_interface *iface = get_interface(data, "eth1");
    assert(iface == NULL);
    free_network_data(data);
    printf("✓ test_get_interface_not_found\n");
}

void test_get_interface_multiple_interfaces() {
    cgroup_network_data *data = create_network_data("test");
    add_network_interface(data, "eth0", 1000, 2000);
    add_network_interface(data, "eth1", 3000, 4000);
    add_network_interface(data, "lo", 100, 100);
    
    cgroup_network_interface *iface = get_interface(data, "eth1");
    assert(iface != NULL);
    assert(strcmp(iface->name, "eth1") == 0);
    assert(iface->bytes_in == 3000);
    free_network_data(data);
    printf("✓ test_get_interface_multiple_interfaces\n");
}

void test_parse_network_stats_success() {
    cgroup_network_interface iface;
    int ret = mock_parse_network_stats("1000 2000", &iface);
    assert(ret == 0);
    assert(iface.bytes_in == 1000);
    assert(iface.bytes_out == 2000);
    printf("✓ test_parse_network_stats_success\n");
}

void test_parse_network_stats_null_str() {
    cgroup_network_interface iface;
    int ret = mock_parse_network_stats(NULL, &iface);
    assert(ret == -1);
    printf("✓ test_parse_network_stats_null_str\n");
}

void test_parse_network_stats_null_iface() {
    int ret = mock_parse_network_stats("1000 2000", NULL);
    assert(ret == -1);
    printf("✓ test_parse_network_stats_null_iface\n");
}

void test_parse_network_stats_empty_str() {
    cgroup_network_interface iface;
    int ret = mock_parse_network_stats("", &iface);
    assert(ret == -1);
    printf("✓ test_parse_network_stats_empty_str\n");
}

void test_parse_network_stats_zero_values() {
    cgroup_network_interface iface;
    int ret = mock_parse_network_stats("0 0", &iface);
    assert(ret == 0);
    assert(iface.bytes_in == 0);
    assert(iface.bytes_out == 0);
    printf("✓ test_parse_network_stats_zero_values\n");
}

void test_read_cgroup_file_success() {
    char buffer[256];
    int ret = mock_read_cgroup_file("/path/to/file", buffer, sizeof(buffer));
    assert(ret > 0);
    assert(strlen(buffer) > 0);
    printf("✓ test_read_cgroup_file_success\n");
}

void test_read_cgroup_file_error_path() {
    char buffer[256];
    int ret = mock_read_cgroup_file("/path/with/error", buffer, sizeof(buffer));
    assert(ret == -1);
    printf("✓ test_read_cgroup_file_error_path\n");
}

void test_read_cgroup_file_empty_file() {
    char buffer[256];
    int ret = mock_read_cgroup_file("/path/with/empty", buffer, sizeof(buffer));
    assert(ret == 0);
    assert(strlen(buffer) == 0);
    printf("✓ test_read_cgroup_file_empty_file\n");
}

int main() {
    printf("Running cgroup-network tests...\n\n");
    
    test_create_network_data_success();
    test_create_network_data_empty_name();
    test_free_network_data_null();
    test_free_network_data_with_interfaces();
    test_add_network_interface_success();
    test_add_network_interface_multiple();
    test_add_network_interface_null_data();
    test_add_network_interface_null_ifname();
    test_add_network_interface_empty_ifname();
    test_add_network_interface_zero_bytes();
    test_add_network_interface_large_values();
    test_read_cgroup_network_stats_null_path();
    test_read_cgroup_network_stats_null_data();
    test_read_cgroup_network_stats_both_null();
    test_read_cgroup_network_stats_file_error();
    test_read_cgroup_network_stats_empty_file();
    test_get_total_bytes_null();
    test_get_total_bytes_empty();
    test_get_total_bytes_with_data();
    test_update_network_interface_success();
    test_update_network_interface_null_data();
    test_update_network_interface_null_ifname();
    test_update_network_interface_not_found();
    test_update_network_interface_zero_values();
    test_update_network_interface_multiple_interfaces();
    test_get_interface_success();
    test_get_interface_null_data();
    test_get_interface_null_ifname();
    test_get_interface_not_found();
    test_get_interface_multiple_interfaces();
    test_parse_network_stats_success();
    test_parse_network_stats_null_str();
    test_parse_network_stats_null_iface();
    test_parse_network_stats_empty_str();
    test_parse_network_stats_zero_values();
    test_read_cgroup_file_success();
    test_read_cgroup_file_error_path();
    test_read_cgroup_file_empty_file();
    
    printf("\n✓ All tests passed!\n");
    return 0;
}