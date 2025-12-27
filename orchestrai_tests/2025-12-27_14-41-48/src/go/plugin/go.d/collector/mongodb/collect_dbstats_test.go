package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCollectDBStats_SuccessfulCollection(t *testing.T) {
	// Arrange
	m := &MongoDB{
		databases: []string{"admin", "test"},
	}
	
	dbStats := bson.M{
		"db":        "test",
		"collections": int32(5),
		"dataSize":  int64(1024000),
		"indexes":   int32(10),
		"indexSize": int64(512000),
		"ok":        float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
	assert.Equal(t, "test", result["database"])
}

func TestCollectDBStats_EmptyDatabase(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{}
	
	// Act
	result := collectDBStats(m, dbStats, "")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_NilDBStats(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	
	// Act
	result := collectDBStats(m, nil, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_WithAllFields(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":           "testdb",
		"collections":  int32(3),
		"dataSize":     int64(2048000),
		"indexes":      int32(5),
		"indexSize":    int64(1024000),
		"avgObjSize":   int32(512),
		"storageSize":  int64(4096000),
		"views":        int32(2),
		"ok":           float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "testdb")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "testdb", result["database"])
}

func TestCollectDBStats_MissingOkField(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":        "test",
		"dataSize":  int64(1024000),
		"ok":        float64(0),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_ZeroValues(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":           "test",
		"collections":  int32(0),
		"dataSize":     int64(0),
		"indexes":      int32(0),
		"indexSize":    int64(0),
		"storageSize":  int64(0),
		"avgObjSize":   int32(0),
		"ok":           float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	require.NotNil(t, result)
}

func TestCollectDBStats_LargeValues(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":          "test",
		"collections": int32(1000),
		"dataSize":    int64(9223372036854775807), // Max int64
		"indexes":     int32(500),
		"indexSize":   int64(4611686018427387903),
		"ok":          float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	require.NotNil(t, result)
}

func TestCollectDBStats_NegativeValues(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":        "test",
		"dataSize":  int64(-1),
		"indexSize": int64(-1),
		"ok":        float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_StringValues(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":        "test",
		"dataSize":  "1024000",
		"indexSize": "512000",
		"ok":        float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_FloatValues(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":        "test",
		"dataSize":  float64(1024000.5),
		"indexSize": float64(512000.25),
		"ok":        float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_Decimal128Values(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	decimal, _ := primitive.ParseDecimal128("1024000.00")
	dbStats := bson.M{
		"db":        "test",
		"dataSize":  decimal,
		"indexSize": decimal,
		"ok":        float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_MultipleCollections(t *testing.T) {
	// Arrange
	m := &MongoDB{
		databases: []string{"admin", "config", "local", "test", "mydb"},
	}
	dbStats := bson.M{
		"db":           "mydb",
		"collections":  int32(100),
		"dataSize":     int64(10240000),
		"indexes":      int32(50),
		"indexSize":    int64(5120000),
		"storageSize":  int64(20480000),
		"avgObjSize":   int32(1024),
		"ok":           float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "mydb")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "mydb", result["database"])
}

func TestCollectDBStats_DatabaseNameWithSpecialChars(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":       "test-db_2024",
		"dataSize": int64(1024000),
		"ok":       float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test-db_2024")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "test-db_2024", result["database"])
}

func TestCollectDBStats_ExtraUnexpectedFields(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":            "test",
		"dataSize":      int64(1024000),
		"ok":            float64(1),
		"extraField1":   "value1",
		"extraField2":   int32(999),
		"extraField3":   bson.M{"nested": "object"},
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_OkFieldFloat(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":   "test",
		"ok":   float64(0.5),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_OkFieldInteger(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db": "test",
		"ok": int32(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_EmptyDatabaseName(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":       "",
		"dataSize": int64(1024000),
		"ok":       float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_SystemDatabase(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":           "admin",
		"collections":  int32(10),
		"dataSize":     int64(2048000),
		"indexes":      int32(20),
		"indexSize":    int64(1024000),
		"storageSize":  int64(4096000),
		"ok":           float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "admin")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "admin", result["database"])
}

func TestCollectDBStats_LocalDatabase(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":       "local",
		"dataSize": int64(512000),
		"ok":       float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "local")
	
	// Assert
	require.NotNil(t, result)
}

func TestCollectDBStats_ConfigDatabase(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":           "config",
		"collections":  int32(15),
		"dataSize":     int64(3072000),
		"indexes":      int32(25),
		"indexSize":    int64(1536000),
		"storageSize":  int64(6144000),
		"ok":           float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "config")
	
	// Assert
	require.NotNil(t, result)
}

func TestCollectDBStats_AllNumericTypesInt32(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":          "test",
		"collections": int32(5),
		"indexes":     int32(10),
		"ok":          float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	require.NotNil(t, result)
}

func TestCollectDBStats_AllNumericTypesInt64(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db":          "test",
		"dataSize":    int64(1024000),
		"indexSize":   int64(512000),
		"storageSize": int64(2048000),
		"ok":          float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	require.NotNil(t, result)
}

func TestCollectDBStats_BooleanValue(t *testing.T) {
	// Arrange
	m := &MongoDB{}
	dbStats := bson.M{
		"db": "test",
		"ok": true,
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_NilCollections(t *testing.T) {
	// Arrange
	m := &MongoDB{
		databases: nil,
	}
	dbStats := bson.M{
		"db": "test",
		"ok": float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "test")
	
	// Assert
	assert.NotNil(t, result)
}

func TestCollectDBStats_SingleDatabase(t *testing.T) {
	// Arrange
	m := &MongoDB{
		databases: []string{"single"},
	}
	dbStats := bson.M{
		"db":       "single",
		"dataSize": int64(1024000),
		"ok":       float64(1),
	}
	
	// Act
	result := collectDBStats(m, dbStats, "single")
	
	// Assert
	require.NotNil(t, result)
	assert.Equal(t, "single", result["database"])
}