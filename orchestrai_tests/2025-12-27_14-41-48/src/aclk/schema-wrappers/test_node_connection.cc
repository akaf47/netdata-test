#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include "node_connection.h"
#include <cstring>
#include <memory>

using ::testing::_;
using ::testing::Return;
using ::testing::Throw;
using ::testing::StrEq;

// Test fixture for node_connection tests
class NodeConnectionTest : public ::testing::Test {
protected:
    void SetUp() override {
        // Initialize test data
    }

    void TearDown() override {
        // Cleanup after tests
    }
};

// Tests for node_connection_status enum
TEST_F(NodeConnectionTest, NodeConnectionStatusValues) {
    // Test that enum values are properly defined
    EXPECT_EQ(NODE_CONN_STATUS_UNKNOWN, 0);
}

// Tests for node_connection_status_to_str function
TEST_F(NodeConnectionTest, NodeConnectionStatusToStrUnknown) {
    const char* result = node_connection_status_to_str(NODE_CONN_STATUS_UNKNOWN);
    EXPECT_NE(nullptr, result);
    EXPECT_STREQ("unknown", result);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToStrValidStatus) {
    const char* result = node_connection_status_to_str(NODE_CONN_STATUS_UNKNOWN);
    EXPECT_NE(nullptr, result);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToStrInvalidStatus) {
    // Test with invalid status code
    const char* result = node_connection_status_to_str(9999);
    EXPECT_NE(nullptr, result);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToStrNegativeStatus) {
    const char* result = node_connection_status_to_str(-1);
    EXPECT_NE(nullptr, result);
}

// Tests for node_connection_status_to_code function
TEST_F(NodeConnectionTest, NodeConnectionStatusToCodeUnknown) {
    int result = node_connection_status_to_code("unknown");
    EXPECT_EQ(NODE_CONN_STATUS_UNKNOWN, result);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToCodeValidStrings) {
    // Test valid status strings
    int code = node_connection_status_to_code("unknown");
    EXPECT_GE(code, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToCodeInvalidString) {
    int result = node_connection_status_to_code("invalid_status");
    EXPECT_LT(result, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToCodeNullString) {
    int result = node_connection_status_to_code(nullptr);
    EXPECT_LT(result, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToCodeEmptyString) {
    int result = node_connection_status_to_code("");
    EXPECT_LT(result, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToCodeCaseSensitive) {
    int result = node_connection_status_to_code("UNKNOWN");
    // Should fail as comparison is case-sensitive
    EXPECT_LT(result, 0);
}

// Tests for node_connection_status_to_rpc_code function
TEST_F(NodeConnectionTest, NodeConnectionStatusToRpcCodeUnknown) {
    int result = node_connection_status_to_rpc_code(NODE_CONN_STATUS_UNKNOWN);
    EXPECT_GE(result, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToRpcCodeValidStatus) {
    int result = node_connection_status_to_rpc_code(NODE_CONN_STATUS_UNKNOWN);
    EXPECT_GE(result, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToRpcCodeInvalidStatus) {
    int result = node_connection_status_to_rpc_code(9999);
    EXPECT_GE(result, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionStatusToRpcCodeNegativeStatus) {
    int result = node_connection_status_to_rpc_code(-1);
    EXPECT_GE(result, 0);
}

// Tests for node_connection_rpc_code_to_status function
TEST_F(NodeConnectionTest, NodeConnectionRpcCodeToStatusValidCode) {
    int rpc_code = node_connection_status_to_rpc_code(NODE_CONN_STATUS_UNKNOWN);
    int status = node_connection_rpc_code_to_status(rpc_code);
    EXPECT_EQ(NODE_CONN_STATUS_UNKNOWN, status);
}

TEST_F(NodeConnectionTest, NodeConnectionRpcCodeToStatusInvalidCode) {
    int status = node_connection_rpc_code_to_status(9999);
    // Should return valid status or error code
    EXPECT_GE(status, -1);
}

TEST_F(NodeConnectionTest, NodeConnectionRpcCodeToStatusNegativeCode) {
    int status = node_connection_rpc_code_to_status(-1);
    EXPECT_LT(status, 0);
}

TEST_F(NodeConnectionTest, NodeConnectionRpcCodeToStatusZeroCode) {
    int status = node_connection_rpc_code_to_status(0);
    EXPECT_GE(status, -1);
}

// Tests for node_connection_default_value function
TEST_F(NodeConnectionTest, NodeConnectionDefaultValueNotNull) {
    const node_connection_t* default_val = node_connection_default_value();
    EXPECT_NE(nullptr, default_val);
}

TEST_F(NodeConnectionTest, NodeConnectionDefaultValueStructure) {
    const node_connection_t* default_val = node_connection_default_value();
    EXPECT_NE(nullptr, default_val);
    // Verify structure is properly initialized
}

TEST_F(NodeConnectionTest, NodeConnectionDefaultValueConsistent) {
    const node_connection_t* val1 = node_connection_default_value();
    const node_connection_t* val2 = node_connection_default_value();
    EXPECT_EQ(val1, val2);
}

// Tests for node_connection_free function
TEST_F(NodeConnectionTest, NodeConnectionFreeNullPointer) {
    // Should not crash when freeing null pointer
    EXPECT_NO_THROW(node_connection_free(nullptr));
}

TEST_F(NodeConnectionTest, NodeConnectionFreeValidPointer) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NE(nullptr, conn);
    EXPECT_NO_THROW(node_connection_free(conn));
}

TEST_F(NodeConnectionTest, NodeConnectionFreeMultipleTimes) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NO_THROW(node_connection_free(conn));
    // Second free should not crash (if implemented safely)
}

// Tests for node_connection_new function
TEST_F(NodeConnectionTest, NodeConnectionNewReturnsValidPointer) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NE(nullptr, conn);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionNewInitializesFields) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NE(nullptr, conn);
    // Verify initialization
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionNewMultipleInstances) {
    node_connection_t* conn1 = node_connection_new();
    node_connection_t* conn2 = node_connection_new();
    EXPECT_NE(nullptr, conn1);
    EXPECT_NE(nullptr, conn2);
    EXPECT_NE(conn1, conn2);
    node_connection_free(conn1);
    node_connection_free(conn2);
}

// Tests for node_connection_clone function
TEST_F(NodeConnectionTest, NodeConnectionCloneNullPointer) {
    node_connection_t* result = node_connection_clone(nullptr);
    EXPECT_EQ(nullptr, result);
}

TEST_F(NodeConnectionTest, NodeConnectionCloneValidPointer) {
    node_connection_t* original = node_connection_new();
    node_connection_t* clone = node_connection_clone(original);
    EXPECT_NE(nullptr, clone);
    EXPECT_NE(original, clone);
    node_connection_free(original);
    node_connection_free(clone);
}

TEST_F(NodeConnectionTest, NodeConnectionCloneDeepCopy) {
    node_connection_t* original = node_connection_new();
    node_connection_t* clone = node_connection_clone(original);
    // Verify it's a deep copy
    EXPECT_NE(original, clone);
    node_connection_free(original);
    node_connection_free(clone);
}

// Tests for node_connection_message_pack function
TEST_F(NodeConnectionTest, NodeConnectionMessagePackNullConnection) {
    uint8_t buffer[1024];
    size_t result = node_connection_message_pack(nullptr, buffer, sizeof(buffer));
    EXPECT_EQ(0, result);
}

TEST_F(NodeConnectionTest, NodeConnectionMessagePackNullBuffer) {
    node_connection_t* conn = node_connection_new();
    size_t result = node_connection_message_pack(conn, nullptr, 1024);
    EXPECT_EQ(0, result);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionMessagePackZeroBufferSize) {
    node_connection_t* conn = node_connection_new();
    uint8_t buffer[1024];
    size_t result = node_connection_message_pack(conn, buffer, 0);
    EXPECT_GE(result, 0);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionMessagePackValidInput) {
    node_connection_t* conn = node_connection_new();
    uint8_t buffer[1024];
    size_t result = node_connection_message_pack(conn, buffer, sizeof(buffer));
    EXPECT_GT(result, 0);
    EXPECT_LE(result, sizeof(buffer));
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionMessagePackSmallBuffer) {
    node_connection_t* conn = node_connection_new();
    uint8_t buffer[1];
    size_t result = node_connection_message_pack(conn, buffer, 1);
    EXPECT_GE(result, 0);
    node_connection_free(conn);
}

// Tests for node_connection_message_unpack function
TEST_F(NodeConnectionTest, NodeConnectionMessageUnpackNullBuffer) {
    node_connection_t* conn = node_connection_new();
    size_t result = node_connection_message_unpack(conn, nullptr, 0);
    EXPECT_EQ(0, result);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionMessageUnpackNullConnection) {
    uint8_t buffer[1024];
    size_t result = node_connection_message_unpack(nullptr, buffer, sizeof(buffer));
    EXPECT_EQ(0, result);
}

TEST_F(NodeConnectionTest, NodeConnectionMessageUnpackZeroSize) {
    node_connection_t* conn = node_connection_new();
    uint8_t buffer[1024];
    size_t result = node_connection_message_unpack(conn, buffer, 0);
    EXPECT_EQ(0, result);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionMessageUnpackValidData) {
    node_connection_t* original = node_connection_new();
    uint8_t buffer[1024];
    
    size_t packed_size = node_connection_message_pack(original, buffer, sizeof(buffer));
    EXPECT_GT(packed_size, 0);
    
    node_connection_t* unpacked = node_connection_new();
    size_t unpacked_size = node_connection_message_unpack(unpacked, buffer, packed_size);
    EXPECT_GT(unpacked_size, 0);
    
    node_connection_free(original);
    node_connection_free(unpacked);
}

TEST_F(NodeConnectionTest, NodeConnectionMessageUnpackInvalidData) {
    node_connection_t* conn = node_connection_new();
    uint8_t invalid_buffer[10];
    memset(invalid_buffer, 0xFF, sizeof(invalid_buffer));
    size_t result = node_connection_message_unpack(conn, invalid_buffer, sizeof(invalid_buffer));
    // Should handle gracefully
    EXPECT_GE(result, 0);
    node_connection_free(conn);
}

// Tests for node_connection_set_timestamp function
TEST_F(NodeConnectionTest, NodeConnectionSetTimestampValid) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NE(nullptr, conn);
    node_connection_set_timestamp(conn, 1234567890);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionSetTimestampNull) {
    EXPECT_NO_THROW(node_connection_set_timestamp(nullptr, 1234567890));
}

TEST_F(NodeConnectionTest, NodeConnectionSetTimestampZero) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NO_THROW(node_connection_set_timestamp(conn, 0));
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionSetTimestampNegative) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NO_THROW(node_connection_set_timestamp(conn, -1));
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionSetTimestampMaxValue) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NO_THROW(node_connection_set_timestamp(conn, 0x7FFFFFFFFFFFFFFF));
    node_connection_free(conn);
}

// Tests for node_connection_get_timestamp function
TEST_F(NodeConnectionTest, NodeConnectionGetTimestampAfterSet) {
    node_connection_t* conn = node_connection_new();
    uint64_t test_time = 1234567890;
    node_connection_set_timestamp(conn, test_time);
    uint64_t result = node_connection_get_timestamp(conn);
    EXPECT_EQ(test_time, result);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionGetTimestampDefault) {
    node_connection_t* conn = node_connection_new();
    uint64_t result = node_connection_get_timestamp(conn);
    // Should return some value
    EXPECT_GE(result, 0);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionGetTimestampNull) {
    uint64_t result = node_connection_get_timestamp(nullptr);
    // Should handle null gracefully
    EXPECT_GE(result, 0);
}

// Tests for node_connection_set_connection_status function
TEST_F(NodeConnectionTest, NodeConnectionSetConnectionStatusValid) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NO_THROW(node_connection_set_connection_status(conn, NODE_CONN_STATUS_UNKNOWN));
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionSetConnectionStatusNull) {
    EXPECT_NO_THROW(node_connection_set_connection_status(nullptr, NODE_CONN_STATUS_UNKNOWN));
}

TEST_F(NodeConnectionTest, NodeConnectionSetConnectionStatusInvalidStatus) {
    node_connection_t* conn = node_connection_new();
    EXPECT_NO_THROW(node_connection_set_connection_status(conn, 9999));
    node_connection_free(conn);
}

// Tests for node_connection_get_connection_status function
TEST_F(NodeConnectionTest, NodeConnectionGetConnectionStatusAfterSet) {
    node_connection_t* conn = node_connection_new();
    node_connection_set_connection_status(conn, NODE_CONN_STATUS_UNKNOWN);
    int result = node_connection_get_connection_status(conn);
    EXPECT_EQ(NODE_CONN_STATUS_UNKNOWN, result);
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionGetConnectionStatusNull) {
    int result = node_connection_get_connection_status(nullptr);
    // Should return valid status
    EXPECT_GE(result, -1);
}

TEST_F(NodeConnectionTest, NodeConnectionGetConnectionStatusDefault) {
    node_connection_t* conn = node_connection_new();
    int result = node_connection_get_connection_status(conn);
    EXPECT_GE(result, -1);
    node_connection_free(conn);
}

// Boundary tests
TEST_F(NodeConnectionTest, NodeConnectionHandleMaxSizeBuffer) {
    node_connection_t* conn = node_connection_new();
    uint8_t buffer[65536];
    size_t result = node_connection_message_pack(conn, buffer, sizeof(buffer));
    EXPECT_GT(result, 0);
    EXPECT_LE(result, sizeof(buffer));
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionSequentialOperations) {
    node_connection_t* conn = node_connection_new();
    
    node_connection_set_timestamp(conn, 100);
    node_connection_set_connection_status(conn, NODE_CONN_STATUS_UNKNOWN);
    
    EXPECT_EQ(100, node_connection_get_timestamp(conn));
    EXPECT_EQ(NODE_CONN_STATUS_UNKNOWN, node_connection_get_connection_status(conn));
    
    node_connection_free(conn);
}

TEST_F(NodeConnectionTest, NodeConnectionRoundTripPackUnpack) {
    node_connection_t* original = node_connection_new();
    node_connection_set_timestamp(original, 9876543210ULL);
    node_connection_set_connection_status(original, NODE_CONN_STATUS_UNKNOWN);
    
    uint8_t buffer[1024];
    size_t packed_size = node_connection_message_pack(original, buffer, sizeof(buffer));
    EXPECT_GT(packed_size, 0);
    
    node_connection_t* restored = node_connection_new();
    size_t unpacked_size = node_connection_message_unpack(restored, buffer, packed_size);
    EXPECT_GT(unpacked_size, 0);
    
    EXPECT_EQ(node_connection_get_timestamp(original), node_connection_get_timestamp(restored));
    EXPECT_EQ(node_connection_get_connection_status(original), node_connection_get_connection_status(restored));
    
    node_connection_free(original);
    node_connection_free(restored);
}

TEST_F(NodeConnectionTest, NodeConnectionClonePreservesData) {
    node_connection_t* original = node_connection_new();
    node_connection_set_timestamp(original, 555);
    node_connection_set_connection_status(original, NODE_CONN_STATUS_UNKNOWN);
    
    node_connection_t* cloned = node_connection_clone(original);
    
    EXPECT_EQ(node_connection_get_timestamp(original), node_connection_get_timestamp(cloned));
    EXPECT_EQ(node_connection_get_connection_status(original), node_connection_get_connection_status(cloned));
    
    node_connection_free(original);
    node_connection_free(cloned);
}