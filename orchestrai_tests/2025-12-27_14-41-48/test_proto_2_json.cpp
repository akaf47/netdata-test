#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <cstring>
#include <cstdlib>
#include <cstdint>
#include "../proto_2_json.h"

using ::testing::Return;
using ::testing::SetArgPointee;
using ::testing::StrEq;
using ::testing::NotNull;
using ::testing::IsNull;
using ::testing::Eq;

// Mock for external dependencies
class MockProtobufMessage {
public:
    virtual ~MockProtobufMessage() = default;
    virtual int serialize(char* buffer, int buffer_size) = 0;
    virtual int deserialize(const char* buffer, int buffer_size) = 0;
};

class Proto2JsonTest : public ::testing::Test {
protected:
    void SetUp() override {
        // Initialize test fixtures
        test_buffer = (char*)malloc(4096);
        memset(test_buffer, 0, 4096);
    }

    void TearDown() override {
        if (test_buffer) {
            free(test_buffer);
            test_buffer = nullptr;
        }
    }

    char* test_buffer;
    static const int BUFFER_SIZE = 4096;
};

// Test proto2json_parse_simple with valid input
TEST_F(Proto2JsonTest, ParseSimpleValidInput) {
    const char* json_input = "{\"type\": \"test\", \"value\": 123}";
    int result = proto2json_parse_simple(json_input);
    EXPECT_NE(result, -1);
}

// Test proto2json_parse_simple with null input
TEST_F(Proto2JsonTest, ParseSimpleNullInput) {
    int result = proto2json_parse_simple(nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_simple with empty string
TEST_F(Proto2JsonTest, ParseSimpleEmptyString) {
    const char* empty = "";
    int result = proto2json_parse_simple(empty);
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_simple with malformed JSON
TEST_F(Proto2JsonTest, ParseSimpleMalformedJson) {
    const char* malformed = "{invalid json}";
    int result = proto2json_parse_simple(malformed);
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_simple with large input
TEST_F(Proto2JsonTest, ParseSimpleLargeInput) {
    std::string large_input(3000, '{');
    large_input += "\"key\": \"value\"";
    large_input += std::string(3000, '}');
    int result = proto2json_parse_simple(large_input.c_str());
    EXPECT_GE(result, -1);
}

// Test proto2json_parse_complex with valid nested structure
TEST_F(Proto2JsonTest, ParseComplexValidNested) {
    const char* nested_json = "{\"outer\": {\"inner\": {\"value\": 42}}}";
    int result = proto2json_parse_complex(nested_json);
    EXPECT_NE(result, -1);
}

// Test proto2json_parse_complex with null input
TEST_F(Proto2JsonTest, ParseComplexNullInput) {
    int result = proto2json_parse_complex(nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_complex with empty object
TEST_F(Proto2JsonTest, ParseComplexEmptyObject) {
    const char* empty_obj = "{}";
    int result = proto2json_parse_complex(empty_obj);
    EXPECT_GE(result, 0);
}

// Test proto2json_parse_complex with array in JSON
TEST_F(Proto2JsonTest, ParseComplexWithArray) {
    const char* array_json = "{\"items\": [1, 2, 3, 4, 5]}";
    int result = proto2json_parse_complex(array_json);
    EXPECT_NE(result, -1);
}

// Test proto2json_serialize with valid message
TEST_F(Proto2JsonTest, SerializeValidMessage) {
    const char* test_msg = "test_message";
    int result = proto2json_serialize(test_msg, test_buffer, BUFFER_SIZE);
    EXPECT_NE(result, -1);
    if (result > 0) {
        EXPECT_NE(test_buffer, nullptr);
    }
}

// Test proto2json_serialize with null message
TEST_F(Proto2JsonTest, SerializeNullMessage) {
    int result = proto2json_serialize(nullptr, test_buffer, BUFFER_SIZE);
    EXPECT_EQ(result, -1);
}

// Test proto2json_serialize with null buffer
TEST_F(Proto2JsonTest, SerializeNullBuffer) {
    const char* test_msg = "message";
    int result = proto2json_serialize(test_msg, nullptr, BUFFER_SIZE);
    EXPECT_EQ(result, -1);
}

// Test proto2json_serialize with zero buffer size
TEST_F(Proto2JsonTest, SerializeZeroBufferSize) {
    const char* test_msg = "message";
    int result = proto2json_serialize(test_msg, test_buffer, 0);
    EXPECT_EQ(result, -1);
}

// Test proto2json_serialize with small buffer
TEST_F(Proto2JsonTest, SerializeSmallBuffer) {
    const char* long_msg = "this_is_a_very_long_message_that_exceeds_buffer_capacity";
    int result = proto2json_serialize(long_msg, test_buffer, 10);
    EXPECT_EQ(result, -1);
}

// Test proto2json_deserialize with valid JSON
TEST_F(Proto2JsonTest, DeserializeValidJson) {
    const char* valid_json = "{\"type\": \"test\", \"id\": 123}";
    int result = proto2json_deserialize(valid_json);
    EXPECT_NE(result, -1);
}

// Test proto2json_deserialize with null input
TEST_F(Proto2JsonTest, DeserializeNullInput) {
    int result = proto2json_deserialize(nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_deserialize with empty string
TEST_F(Proto2JsonTest, DeserializeEmptyString) {
    const char* empty = "";
    int result = proto2json_deserialize(empty);
    EXPECT_EQ(result, -1);
}

// Test proto2json_convert_field with string field
TEST_F(Proto2JsonTest, ConvertFieldString) {
    const char* field_name = "name";
    const char* field_value = "test_value";
    int result = proto2json_convert_field(field_name, field_value, nullptr);
    EXPECT_NE(result, -1);
}

// Test proto2json_convert_field with numeric field
TEST_F(Proto2JsonTest, ConvertFieldNumeric) {
    const char* field_name = "count";
    const char* field_value = "42";
    int result = proto2json_convert_field(field_name, field_value, nullptr);
    EXPECT_NE(result, -1);
}

// Test proto2json_convert_field with null field name
TEST_F(Proto2JsonTest, ConvertFieldNullName) {
    int result = proto2json_convert_field(nullptr, "value", nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_convert_field with null field value
TEST_F(Proto2JsonTest, ConvertFieldNullValue) {
    const char* field_name = "field";
    int result = proto2json_convert_field(field_name, nullptr, nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_convert_field with empty field name
TEST_F(Proto2JsonTest, ConvertFieldEmptyName) {
    int result = proto2json_convert_field("", "value", nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_validate_schema with valid schema
TEST_F(Proto2JsonTest, ValidateSchemaValid) {
    const char* valid_schema = "{\"version\": \"1.0\", \"type\": \"object\"}";
    int result = proto2json_validate_schema(valid_schema);
    EXPECT_EQ(result, 1);
}

// Test proto2json_validate_schema with invalid schema
TEST_F(Proto2JsonTest, ValidateSchemainvalid) {
    const char* invalid_schema = "not_a_schema";
    int result = proto2json_validate_schema(invalid_schema);
    EXPECT_EQ(result, 0);
}

// Test proto2json_validate_schema with null
TEST_F(Proto2JsonTest, ValidateSchemaNullInput) {
    int result = proto2json_validate_schema(nullptr);
    EXPECT_EQ(result, 0);
}

// Test proto2json_get_error_message after error
TEST_F(Proto2JsonTest, GetErrorMessageAfterError) {
    proto2json_parse_simple(nullptr);
    const char* error_msg = proto2json_get_error_message();
    EXPECT_NE(error_msg, nullptr);
    EXPECT_NE(strlen(error_msg), 0);
}

// Test proto2json_get_error_message on clean state
TEST_F(Proto2JsonTest, GetErrorMessageCleanState) {
    const char* error_msg = proto2json_get_error_message();
    if (error_msg != nullptr) {
        EXPECT_TRUE(strlen(error_msg) == 0 || strcmp(error_msg, "OK") == 0);
    }
}

// Test proto2json_clear_error
TEST_F(Proto2JsonTest, ClearError) {
    proto2json_parse_simple(nullptr);
    proto2json_clear_error();
    const char* error_msg = proto2json_get_error_message();
    if (error_msg != nullptr) {
        EXPECT_TRUE(strlen(error_msg) == 0 || strcmp(error_msg, "OK") == 0);
    }
}

// Test proto2json_init with null
TEST_F(Proto2JsonTest, InitNull) {
    int result = proto2json_init(nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_init with valid config
TEST_F(Proto2JsonTest, InitValidConfig) {
    const char* config = "{\"strict\": true}";
    int result = proto2json_init(config);
    EXPECT_NE(result, -1);
}

// Test proto2json_cleanup
TEST_F(Proto2JsonTest, Cleanup) {
    proto2json_init("{\"strict\": true}");
    int result = proto2json_cleanup();
    EXPECT_EQ(result, 0);
}

// Test proto2json_reset_state
TEST_F(Proto2JsonTest, ResetState) {
    int result = proto2json_reset_state();
    EXPECT_EQ(result, 0);
}

// Test proto2json_set_strict_mode
TEST_F(Proto2JsonTest, SetStrictModeTrue) {
    int result = proto2json_set_strict_mode(1);
    EXPECT_EQ(result, 0);
}

// Test proto2json_set_strict_mode false
TEST_F(Proto2JsonTest, SetStrictModeFalse) {
    int result = proto2json_set_strict_mode(0);
    EXPECT_EQ(result, 0);
}

// Test proto2json_get_strict_mode
TEST_F(Proto2JsonTest, GetStrictMode) {
    proto2json_set_strict_mode(1);
    int mode = proto2json_get_strict_mode();
    EXPECT_EQ(mode, 1);
}

// Test proto2json_parse_array with valid array
TEST_F(Proto2JsonTest, ParseArrayValid) {
    const char* array_json = "[1, 2, 3, 4, 5]";
    int result = proto2json_parse_array(array_json);
    EXPECT_NE(result, -1);
}

// Test proto2json_parse_array with empty array
TEST_F(Proto2JsonTest, ParseArrayEmpty) {
    const char* empty_array = "[]";
    int result = proto2json_parse_array(empty_array);
    EXPECT_NE(result, -1);
}

// Test proto2json_parse_array with null
TEST_F(Proto2JsonTest, ParseArrayNull) {
    int result = proto2json_parse_array(nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_array with nested arrays
TEST_F(Proto2JsonTest, ParseArrayNested) {
    const char* nested = "[[1, 2], [3, 4], [5, 6]]";
    int result = proto2json_parse_array(nested);
    EXPECT_NE(result, -1);
}

// Test proto2json_convert_type with boolean
TEST_F(Proto2JsonTest, ConvertTypeBool) {
    const char* bool_val = "true";
    int result = proto2json_convert_type("bool", bool_val);
    EXPECT_NE(result, -1);
}

// Test proto2json_convert_type with integer
TEST_F(Proto2JsonTest, ConvertTypeInt) {
    const char* int_val = "123";
    int result = proto2json_convert_type("int", int_val);
    EXPECT_NE(result, -1);
}

// Test proto2json_convert_type with float
TEST_F(Proto2JsonTest, ConvertTypeFloat) {
    const char* float_val = "3.14";
    int result = proto2json_convert_type("float", float_val);
    EXPECT_NE(result, -1);
}

// Test proto2json_convert_type with string
TEST_F(Proto2JsonTest, ConvertTypeString) {
    const char* str_val = "hello";
    int result = proto2json_convert_type("string", str_val);
    EXPECT_NE(result, -1);
}

// Test proto2json_convert_type with unknown type
TEST_F(Proto2JsonTest, ConvertTypeUnknown) {
    int result = proto2json_convert_type("unknown_type", "value");
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_with_schema with matching schema
TEST_F(Proto2JsonTest, ParseWithSchemaMatching) {
    const char* schema = "{\"properties\": {\"name\": {\"type\": \"string\"}}}";
    const char* data = "{\"name\": \"test\"}";
    int result = proto2json_parse_with_schema(schema, data);
    EXPECT_NE(result, -1);
}

// Test proto2json_parse_with_schema with non-matching schema
TEST_F(Proto2JsonTest, ParseWithSchemaMismatch) {
    const char* schema = "{\"properties\": {\"name\": {\"type\": \"string\"}}}";
    const char* data = "{\"name\": 123}";
    int result = proto2json_parse_with_schema(schema, data);
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_with_schema with null schema
TEST_F(Proto2JsonTest, ParseWithSchemaNullSchema) {
    const char* data = "{\"name\": \"test\"}";
    int result = proto2json_parse_with_schema(nullptr, data);
    EXPECT_EQ(result, -1);
}

// Test proto2json_parse_with_schema with null data
TEST_F(Proto2JsonTest, ParseWithSchemaNullData) {
    const char* schema = "{\"properties\": {\"name\": {\"type\": \"string\"}}}";
    int result = proto2json_parse_with_schema(schema, nullptr);
    EXPECT_EQ(result, -1);
}

// Test proto2json_get_buffer_size
TEST_F(Proto2JsonTest, GetBufferSize) {
    int size = proto2json_get_buffer_size();
    EXPECT_GT(size, 0);
}

// Test proto2json_set_buffer_size with valid size
TEST_F(Proto2JsonTest, SetBufferSizeValid) {
    int result = proto2json_set_buffer_size(8192);
    EXPECT_EQ(result, 0);
}

// Test proto2json_set_buffer_size with invalid size
TEST_F(Proto2JsonTest, SetBufferSizeInvalid) {
    int result = proto2json_set_buffer_size(0);
    EXPECT_EQ(result, -1);
}

// Test proto2json_get_version
TEST_F(Proto2JsonTest, GetVersion) {
    const char* version = proto2json_get_version();
    EXPECT_NE(version, nullptr);
    EXPECT_GT(strlen(version), 0);
}

// Boundary test: maximum integer value
TEST_F(Proto2JsonTest, BoundaryMaxInt) {
    char int_str[20];
    snprintf(int_str, sizeof(int_str), "%d", INT_MAX);
    int result = proto2json_convert_type("int", int_str);
    EXPECT_NE(result, -1);
}

// Boundary test: minimum integer value
TEST_F(Proto2JsonTest, BoundaryMinInt) {
    char int_str[20];
    snprintf(int_str, sizeof(int_str), "%d", INT_MIN);
    int result = proto2json_convert_type("int", int_str);
    EXPECT_NE(result, -1);
}

// Boundary test: very large JSON string
TEST_F(Proto2JsonTest, BoundaryLargeJson) {
    std::string large_json = "{";
    for (int i = 0; i < 1000; ++i) {
        large_json += "\"key" + std::to_string(i) + "\": \"value" + std::to_string(i) + "\",";
    }
    large_json.pop_back();
    large_json += "}";
    int result = proto2json_parse_simple(large_json.c_str());
    EXPECT_NE(result, -1);
}

// Test error recovery after multiple errors
TEST_F(Proto2JsonTest, ErrorRecoveryMultiple) {
    proto2json_parse_simple(nullptr);
    proto2json_deserialize(nullptr);
    proto2json_clear_error();
    const char* test_json = "{\"key\": \"value\"}";
    int result = proto2json_parse_simple(test_json);
    EXPECT_NE(result, -1);
}

// Test concurrent initialization (sequential)
TEST_F(Proto2JsonTest, SequentialInit) {
    int result1 = proto2json_init("{\"mode\": \"test\"}");
    EXPECT_NE(result1, -1);
    int result2 = proto2json_init("{\"mode\": \"test2\"}");
    EXPECT_NE(result2, -1);
    proto2json_cleanup();
}

// Test state persistence
TEST_F(Proto2JsonTest, StatePersistence) {
    proto2json_set_strict_mode(1);
    const char* test_json = "{\"key\": \"value\"}";
    proto2json_parse_simple(test_json);
    int mode = proto2json_get_strict_mode();
    EXPECT_EQ(mode, 1);
}

// Test special characters in JSON
TEST_F(Proto2JsonTest, SpecialCharactersInJson) {
    const char* special_json = "{\"text\": \"Special chars: \\n \\t \\r \\\\ \\\" \"}";
    int result = proto2json_parse_simple(special_json);
    EXPECT_NE(result, -1);
}

// Test unicode handling
TEST_F(Proto2JsonTest, UnicodeInJson) {
    const char* unicode_json = "{\"emoji\": \"\\u0048\\u0065\\u006c\\u006c\\u006f\"}";
    int result = proto2json_parse_simple(unicode_json);
    EXPECT_NE(result, -1);
}

// Test floating point precision
TEST_F(Proto2JsonTest, FloatPrecision) {
    const char* float_json = "{\"value\": 3.141592653589793}";
    int result = proto2json_parse_simple(float_json);
    EXPECT_NE(result, -1);
}