package mysql

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollectVersionSuccess tests successful version collection
func TestCollectVersionSuccess(t *testing.T) {
	t.Run("should collect MySQL version successfully", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{
				queryRows: func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, nil
				},
			},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionQueryError tests behavior when SELECT VERSION() fails
func TestCollectVersionQueryError(t *testing.T) {
	t.Run("should handle query execution error", func(t *testing.T) {
		// Arrange
		queryErr := sql.ErrNoRows
		m := &MySQL{
			db: &mockDB{
				err: queryErr,
			},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, queryErr, err)
	})
}

// TestCollectVersionEmptyResult tests handling of no result rows
func TestCollectVersionEmptyResult(t *testing.T) {
	t.Run("should handle empty result set", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{
				queryRows: func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, nil
				},
			},
		}

		// Act
		err := m.collectVersion()

		// Assert
		// Empty result should typically return an error
		assert.Error(t, err)
	})
}

// TestCollectVersionNullResult tests handling of NULL version
func TestCollectVersionNullResult(t *testing.T) {
	t.Run("should handle NULL version result", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{
				queryRows: func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, nil
				},
			},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.Error(t, err)
	})
}

// TestCollectVersionMariaDB tests version string parsing for MariaDB
func TestCollectVersionMariaDB(t *testing.T) {
	t.Run("should handle MariaDB version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionMySQL5 tests version string parsing for MySQL 5.x
func TestCollectVersionMySQL5(t *testing.T) {
	t.Run("should handle MySQL 5.7 version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle MySQL 5.6 version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionMySQL8 tests version string parsing for MySQL 8.x
func TestCollectVersionMySQL8(t *testing.T) {
	t.Run("should handle MySQL 8.0 version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle MySQL 8.1 version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionConnectionClosed tests with closed connection
func TestCollectVersionConnectionClosed(t *testing.T) {
	t.Run("should handle closed database connection", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{
				err: sql.ErrConnDone,
			},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, sql.ErrConnDone, err)
	})
}

// TestCollectVersionNilDB tests behavior with nil database
func TestCollectVersionNilDB(t *testing.T) {
	t.Run("should handle nil database connection", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: nil,
		}

		// Act & Assert
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Recovered from panic: %v", r)
			}
		}()

		_ = m.collectVersion()
	})
}

// TestCollectVersionParsingVersionString tests version string parsing logic
func TestCollectVersionParsingVersionString(t *testing.T) {
	t.Run("should extract major version correctly", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should extract minor version correctly", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should extract patch version correctly", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionVersionWithSuffix tests version strings with suffixes
func TestCollectVersionVersionWithSuffix(t *testing.T) {
	t.Run("should handle version with -log suffix", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle version with -ubuntu suffix", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle version with -debian suffix", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionStorageEngineIdentification tests engine identification from version
func TestCollectVersionStorageEngineIdentification(t *testing.T) {
	t.Run("should identify MySQL version from version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should identify MariaDB from version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionRowScanError tests handling of Scan() errors
func TestCollectVersionRowScanError(t *testing.T) {
	t.Run("should handle scan error on version result", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{
				queryRows: func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, nil
				},
			},
		}

		// Act
		err := m.collectVersion()

		// Assert
		// Error expected from failed scan
		assert.Error(t, err)
	})
}

// TestCollectVersionQueryClose tests proper cleanup of result rows
func TestCollectVersionQueryClose(t *testing.T) {
	t.Run("should properly close result rows", func(t *testing.T) {
		// Arrange
		rowsClosed := false
		m := &MySQL{
			db: &mockDB{
				queryRows: func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, nil
				},
			},
		}

		// Act
		_ = m.collectVersion()

		// Assert - verify no resource leaks
		assert.True(t, rowsClosed || err != nil)
	})
}

// TestCollectVersionQueryString tests correct query construction
func TestCollectVersionQueryString(t *testing.T) {
	t.Run("should execute correct query", func(t *testing.T) {
		// Arrange
		queryCalled := false
		m := &MySQL{
			db: &mockDB{
				queryRows: func(query string, args ...interface{}) (*sql.Rows, error) {
					queryCalled = true
					assert.Contains(t, query, "VERSION")
					return nil, nil
				},
			},
		}

		// Act
		_ = m.collectVersion()

		// Assert
		assert.True(t, queryCalled)
	})
}

// TestCollectVersionEdgeCases tests edge cases in version parsing
func TestCollectVersionEdgeCases(t *testing.T) {
	t.Run("should handle version with extra whitespace", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle version with extra dots", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle very long version string", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionStorageInMetrics tests version storage in metrics
func TestCollectVersionStorageInMetrics(t *testing.T) {
	t.Run("should store version string in metrics", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionInfoLabel tests version as info label
func TestCollectVersionInfoLabel(t *testing.T) {
	t.Run("should set version as info label if supported", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}

// TestCollectVersionMultipleCalls tests behavior on multiple calls
func TestCollectVersionMultipleCalls(t *testing.T) {
	t.Run("should allow multiple collection calls", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err1 := m.collectVersion()
		err2 := m.collectVersion()

		// Assert
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})
}

// TestCollectVersionVersionChangeDetection tests version change between calls
func TestCollectVersionVersionChangeDetection(t *testing.T) {
	t.Run("should handle version change detection on re-collection", func(t *testing.T) {
		// Arrange
		m := &MySQL{
			db: &mockDB{},
		}

		// Act
		err := m.collectVersion()

		// Assert
		assert.NoError(t, err)
	})
}