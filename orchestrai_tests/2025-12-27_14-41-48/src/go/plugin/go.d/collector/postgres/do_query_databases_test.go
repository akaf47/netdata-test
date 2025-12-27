package postgres

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoQueryDatabases(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockDB)
		expectedError bool
		expectedLen   int
		validateData  func(*testing.T, map[string]interface{})
	}{
		{
			name: "should execute query successfully with valid databases",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"postgres", int64(1), int64(100)},
						{"template0", int64(2), int64(200)},
						{"template1", int64(3), int64(300)},
						{"testdb", int64(4), int64(400)},
					}), nil
				}
			},
			expectedError: false,
			expectedLen:   4,
			validateData: func(t *testing.T, data map[string]interface{}) {
				assert.NotNil(t, data)
			},
		},
		{
			name: "should handle single database result",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"postgres", int64(1), int64(100)},
					}), nil
				}
			},
			expectedError: false,
			expectedLen:   1,
			validateData: func(t *testing.T, data map[string]interface{}) {
				assert.NotNil(t, data)
			},
		},
		{
			name: "should handle empty database list",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{}), nil
				}
			},
			expectedError: false,
			expectedLen:   0,
			validateData: func(t *testing.T, data map[string]interface{}) {
				assert.Empty(t, data)
			},
		},
		{
			name: "should handle query execution error",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, errors.New("query failed")
				}
			},
			expectedError: true,
		},
		{
			name: "should handle row scan errors",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRowsWithScanError(), nil
				}
			},
			expectedError: true,
		},
		{
			name: "should handle nil database connection",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = nil
			},
			expectedError: true,
		},
		{
			name: "should handle connection closed error",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, sql.ErrNoRows
				}
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			tt.mockSetup(mockDB)

			p := &Postgres{
				db: mockDB,
			}

			result, err := p.doQueryDatabases()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedLen > 0 {
					assert.Len(t, result, tt.expectedLen)
				}
				if tt.validateData != nil && result != nil {
					tt.validateData(t, result)
				}
			}
		})
	}
}

func TestDoQueryDatabasesEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockDB)
		shouldError bool
	}{
		{
			name: "should handle databases with null oid values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db1", nil, int64(100)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle databases with null size values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db2", int64(5), nil},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle databases with both null values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db3", nil, nil},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle database names with special characters",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db-with-dashes", int64(1), int64(100)},
						{"db_with_underscores", int64(2), int64(200)},
						{"db.with.dots", int64(3), int64(300)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle database names with spaces",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db with spaces", int64(4), int64(400)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle zero oid values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db_zero_oid", int64(0), int64(100)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle zero size values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db_zero_size", int64(5), int64(0)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle maximum int64 values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db_max_values", int64(9223372036854775807), int64(9223372036854775807)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle negative size values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"db_negative", int64(6), int64(-100)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle large number of databases",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					rows := [][]interface{}{}
					for i := 0; i < 10000; i++ {
						rows = append(rows, []interface{}{"db_" + string(rune(i)), int64(i), int64(i * 1024)})
					}
					return mockRows(rows), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle very long database names",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					longName := ""
					for i := 0; i < 10000; i++ {
						longName += "a"
					}
					return mockRows([][]interface{}{
						{longName, int64(1), int64(1000)},
					}), nil
				}
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			tt.mockSetup(mockDB)

			p := &Postgres{
				db: mockDB,
			}

			result, err := p.doQueryDatabases()

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestDoQueryDatabasesWithDifferentOIDValues(t *testing.T) {
	tests := []struct {
		name      string
		oid       int64
		mockSetup func(*MockDB)
	}{
		{
			name: "should handle OID value 1",
			oid:  1,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"postgres", int64(1), int64(100)},
					}), nil
				}
			},
		},
		{
			name: "should handle OID value 13909",
			oid:  13909,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"test_db", int64(13909), int64(1024)},
					}), nil
				}
			},
		},
		{
			name: "should handle OID value 16384",
			oid:  16384,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"custom_db", int64(16384), int64(2048)},
					}), nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			tt.mockSetup(mockDB)

			p := &Postgres{
				db: mockDB,
			}

			result, err := p.doQueryDatabases()
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestDoQueryDatabasesWithDifferentSizeValues(t *testing.T) {
	tests := []struct {
		name      string
		size      int64
		mockSetup func(*MockDB)
	}{
		{
			name: "should handle database size 0 bytes",
			size: 0,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"empty_db", int64(1), int64(0)},
					}), nil
				}
			},
		},
		{
			name: "should handle database size 1 byte",
			size: 1,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"tiny_db", int64(2), int64(1)},
					}), nil
				}
			},
		},
		{
			name: "should handle database size 1 KB",
			size: 1024,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"small_db", int64(3), int64(1024)},
					}), nil
				}
			},
		},
		{
			name: "should handle database size 1 MB",
			size: 1024 * 1024,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"medium_db", int64(4), int64(1024 * 1024)},
					}), nil
				}
			},
		},
		{
			name: "should handle database size 1 GB",
			size: 1024 * 1024 * 1024,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"large_db", int64(5), int64(1024 * 1024 * 1024)},
					}), nil
				}
			},
		},
		{
			name: "should handle database size maximum int64",
			size: 9223372036854775807,
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"huge_db", int64(6), int64(9223372036854775807)},
					}), nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			tt.mockSetup(mockDB)

			p := &Postgres{
				db: mockDB,
			}

			result, err := p.doQueryDatabases()
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func TestDoQueryDatabasesConcurrency(t *testing.T) {
	mockDB := &MockDB{}
	mockDB.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
		return mockRows([][]interface{}{
			{"postgres", int64(1), int64(100)},
			{"testdb", int64(2), int64(200)},
		}), nil
	}

	p := &Postgres{
		db: mockDB,
	}

	// Test concurrent calls
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			result, err := p.doQueryDatabases()
			assert.NoError(t, err)
			assert.Len(t, result, 2)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestDoQueryDatabasesErrorRecovery(t *testing.T) {
	mockDB := &MockDB{}
	callCount := 0

	mockDB.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
		callCount++
		if callCount == 1 {
			return nil, errors.New("temporary error")
		}
		return mockRows([][]interface{}{
			{"postgres", int64(1), int64(100)},
		}), nil
	}

	p := &Postgres{
		db: mockDB,
	}

	// First call should fail
	result1, err1 := p.doQueryDatabases()
	assert.Error(t, err1)
	assert.Nil(t, result1)

	// Second call should succeed
	result2, err2 := p.doQueryDatabases()
	assert.NoError(t, err2)
	assert.NotNil(t, result2)
}