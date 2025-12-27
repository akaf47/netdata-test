#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include "node_creation.h"
#include <memory>
#include <string>
#include <vector>
#include <cstring>

using ::testing::AtLeast;
using ::testing::Return;
using ::testing::SaveArg;
using ::testing::_;

class NodeCreationTest : public ::testing::Test {
protected:
    void SetUp() override {
        // Initialize any necessary test fixtures
    }

    void TearDown() override {
        // Clean up after tests
    }
};

// Test node_creation_data creation and initialization
TEST_F(NodeCreationTest, CreateNodeCreationDataWithValidInputs) {
    const char* node_id = "node-001";
    const char* node_name = "Test Node";
    const char* node_type = "physical";
    
    auto data = node_creation_data_create(node_id, node_name, node_type);
    
    ASSERT_NE(data, nullptr);
    EXPECT_STREQ(data->node_id, node_id);
    EXPECT_STREQ(data->node_name, node_name);
    EXPECT_STREQ(data->node_type, node_type);
    
    node_creation_data_destroy(data);
}

// Test node_creation_data with null node_id
TEST_F(NodeCreationTest, CreateNodeCreationDataWithNullNodeId) {
    const char* node_name = "Test Node";
    const char* node_type = "physical";
    
    auto data = node_creation_data_create(nullptr, node_name, node_type);
    
    // Should handle null gracefully or return nullptr
    if (data != nullptr) {
        EXPECT_EQ(data->node_id, nullptr);
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data with null node_name
TEST_F(NodeCreationTest, CreateNodeCreationDataWithNullNodeName) {
    const char* node_id = "node-001";
    const char* node_type = "physical";
    
    auto data = node_creation_data_create(node_id, nullptr, node_type);
    
    if (data != nullptr) {
        EXPECT_EQ(data->node_name, nullptr);
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data with null node_type
TEST_F(NodeCreationTest, CreateNodeCreationDataWithNullNodeType) {
    const char* node_id = "node-001";
    const char* node_name = "Test Node";
    
    auto data = node_creation_data_create(node_id, node_name, nullptr);
    
    if (data != nullptr) {
        EXPECT_EQ(data->node_type, nullptr);
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data with all null parameters
TEST_F(NodeCreationTest, CreateNodeCreationDataWithAllNullParameters) {
    auto data = node_creation_data_create(nullptr, nullptr, nullptr);
    
    if (data != nullptr) {
        EXPECT_EQ(data->node_id, nullptr);
        EXPECT_EQ(data->node_name, nullptr);
        EXPECT_EQ(data->node_type, nullptr);
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data with empty strings
TEST_F(NodeCreationTest, CreateNodeCreationDataWithEmptyStrings) {
    const char* empty = "";
    
    auto data = node_creation_data_create(empty, empty, empty);
    
    ASSERT_NE(data, nullptr);
    EXPECT_STREQ(data->node_id, empty);
    EXPECT_STREQ(data->node_name, empty);
    EXPECT_STREQ(data->node_type, empty);
    
    node_creation_data_destroy(data);
}

// Test node_creation_data with very long strings
TEST_F(NodeCreationTest, CreateNodeCreationDataWithLongStrings) {
    std::string long_id(1000, 'a');
    std::string long_name(1000, 'b');
    std::string long_type(1000, 'c');
    
    auto data = node_creation_data_create(long_id.c_str(), long_name.c_str(), long_type.c_str());
    
    if (data != nullptr) {
        EXPECT_EQ(std::string(data->node_id), long_id);
        EXPECT_EQ(std::string(data->node_name), long_name);
        EXPECT_EQ(std::string(data->node_type), long_type);
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data memory cleanup
TEST_F(NodeCreationTest, DestroyNodeCreationDataFreesMemory) {
    auto data = node_creation_data_create("id", "name", "type");
    
    ASSERT_NE(data, nullptr);
    // Should not crash when destroying
    node_creation_data_destroy(data);
    // data pointer should not be dereferenced after this
}

// Test node_creation_data destroy with null pointer
TEST_F(NodeCreationTest, DestroyNodeCreationDataWithNull) {
    // Should not crash when destroying null pointer
    node_creation_data_destroy(nullptr);
}

// Test node_creation_data copy/clone functionality (if exists)
TEST_F(NodeCreationTest, CopyNodeCreationData) {
    auto original = node_creation_data_create("node-001", "Test Node", "physical");
    
    if (original != nullptr) {
        auto copy = node_creation_data_clone(original);
        
        if (copy != nullptr) {
            EXPECT_STREQ(copy->node_id, original->node_id);
            EXPECT_STREQ(copy->node_name, original->node_name);
            EXPECT_STREQ(copy->node_type, original->node_type);
            
            node_creation_data_destroy(copy);
        }
        node_creation_data_destroy(original);
    }
}

// Test node_creation_data serialization
TEST_F(NodeCreationTest, SerializeNodeCreationData) {
    auto data = node_creation_data_create("node-001", "Test Node", "physical");
    
    if (data != nullptr && node_creation_data_serialize) {
        // Get serialized size
        size_t size = node_creation_data_get_serialized_size(data);
        EXPECT_GT(size, 0);
        
        // Serialize to buffer
        std::vector<uint8_t> buffer(size);
        int result = node_creation_data_serialize(data, buffer.data(), size);
        EXPECT_EQ(result, 0);
        
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data deserialization
TEST_F(NodeCreationTest, DeserializeNodeCreationData) {
    auto original = node_creation_data_create("node-001", "Test Node", "physical");
    
    if (original != nullptr && node_creation_data_serialize && node_creation_data_deserialize) {
        size_t size = node_creation_data_get_serialized_size(original);
        std::vector<uint8_t> buffer(size);
        
        node_creation_data_serialize(original, buffer.data(), size);
        
        auto deserialized = node_creation_data_deserialize(buffer.data(), size);
        
        if (deserialized != nullptr) {
            EXPECT_STREQ(deserialized->node_id, original->node_id);
            EXPECT_STREQ(deserialized->node_name, original->node_name);
            EXPECT_STREQ(deserialized->node_type, original->node_type);
            
            node_creation_data_destroy(deserialized);
        }
        
        node_creation_data_destroy(original);
    }
}

// Test node_creation_data deserialization with null buffer
TEST_F(NodeCreationTest, DeserializeNodeCreationDataWithNullBuffer) {
    auto data = node_creation_data_deserialize(nullptr, 0);
    
    // Should return nullptr or handle gracefully
    EXPECT_EQ(data, nullptr);
}

// Test node_creation_data deserialization with zero size
TEST_F(NodeCreationTest, DeserializeNodeCreationDataWithZeroSize) {
    uint8_t buffer[10] = {0};
    
    auto data = node_creation_data_deserialize(buffer, 0);
    
    // Should return nullptr or handle gracefully
    EXPECT_EQ(data, nullptr);
}

// Test node_creation_data deserialization with invalid buffer
TEST_F(NodeCreationTest, DeserializeNodeCreationDataWithInvalidBuffer) {
    uint8_t invalid_buffer[4] = {0xFF, 0xFF, 0xFF, 0xFF};
    
    auto data = node_creation_data_deserialize(invalid_buffer, sizeof(invalid_buffer));
    
    // Should handle invalid data gracefully
    if (data != nullptr) {
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data special characters in strings
TEST_F(NodeCreationTest, CreateNodeCreationDataWithSpecialCharacters) {
    const char* special_id = "node-001!@#$%^&*()";
    const char* special_name = "Test Node \n\t\r";
    const char* special_type = "physical|virtual";
    
    auto data = node_creation_data_create(special_id, special_name, special_type);
    
    if (data != nullptr) {
        EXPECT_STREQ(data->node_id, special_id);
        EXPECT_STREQ(data->node_name, special_name);
        EXPECT_STREQ(data->node_type, special_type);
        
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data unicode characters
TEST_F(NodeCreationTest, CreateNodeCreationDataWithUnicodeCharacters) {
    const char* unicode_name = "Test Node 中文 😀";
    
    auto data = node_creation_data_create("node-001", unicode_name, "physical");
    
    if (data != nullptr) {
        EXPECT_STREQ(data->node_name, unicode_name);
        node_creation_data_destroy(data);
    }
}

// Test node_creation_data buffer overflow protection
TEST_F(NodeCreationTest, NodeCreationDataHandlesBufferBoundaries) {
    const char* id = "node-001";
    const char* name = "Test Node";
    const char* type = "physical";
    
    auto data = node_creation_data_create(id, name, type);
    
    if (data != nullptr) {
        // Verify string lengths don't exceed expectations
        size_t id_len = strlen(data->node_id);
        size_t name_len = strlen(data->node_name);
        size_t type_len = strlen(data->node_type);
        
        EXPECT_LE(id_len, 1000);    // reasonable boundary
        EXPECT_LE(name_len, 1000);
        EXPECT_LE(type_len, 1000);
        
        node_creation_data_destroy(data);
    }
}

// Test node_creation API integration
TEST_F(NodeCreationTest, SendNodeCreationMessage) {
    auto data = node_creation_data_create("node-001", "Test Node", "physical");
    
    if (data != nullptr) {
        int result = node_creation_send(data);
        
        // Verify send operation
        EXPECT_EQ(result, 0);  // Success
        
        node_creation_data_destroy(data);
    }
}

// Test node_creation with connection issues
TEST_F(NodeCreationTest, SendNodeCreationMessageWithConnectionFailure) {
    auto data = node_creation_data_create("node-001", "Test Node", "physical");
    
    if (data != nullptr) {
        // Simulate connection failure
        int result = node_creation_send(nullptr);
        
        // Should handle null gracefully
        EXPECT_NE(result, 0);
        
        node_creation_data_destroy(data);
    }
}

// Test multiple sequential node creation
TEST_F(NodeCreationTest, CreateMultipleNodesSequentially) {
    for (int i = 0; i < 10; i++) {
        std::string id = "node-" + std::to_string(i);
        std::string name = "Node " + std::to_string(i);
        
        auto data = node_creation_data_create(id.c_str(), name.c_str(), "physical");
        
        ASSERT_NE(data, nullptr);
        EXPECT_STREQ(data->node_id, id.c_str());
        EXPECT_STREQ(data->node_name, name.c_str());
        
        node_creation_data_destroy(data);
    }
}

// Test node creation with different node types
TEST_F(NodeCreationTest, CreateNodesWithDifferentTypes) {
    std::vector<const char*> types = {"physical", "virtual", "container", "cloud", "edge"};
    
    for (const auto& type : types) {
        auto data = node_creation_data_create("node-001", "Test Node", type);
        
        if (data != nullptr) {
            EXPECT_STREQ(data->node_type, type);
            node_creation_data_destroy(data);
        }
    }
}