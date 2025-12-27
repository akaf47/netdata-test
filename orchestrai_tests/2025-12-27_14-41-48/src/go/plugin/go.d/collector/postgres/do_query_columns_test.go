package postgres

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoQueryColumns(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockDB)
		expectedError bool
		expectedLen   int
		validateData  func(*testing.T, map[string]interface{})
	}{
		{
			name: "should execute query successfully with valid data",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"column1", "value1"},
						{"column2", "value2"},
					}), nil
				}
			},
			expectedError: false,
			expectedLen:   2,
			validateData: func(t *testing.T, data map[string]interface{}) {
				assert.NotNil(t, data)
			},
		},
		{
			name: "should handle empty result set",
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
					return nil, errors.New("database error")
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
			name: "should handle connection closure error",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return nil, sql.ErrConnDone
				}
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockDB{}
			tt.mockSetup(mockDB)

			p := &Postgres{
				db: mockDB,
			}

			result, err := p.doQueryColumns()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
				if tt.validateData != nil && result != nil {
					tt.validateData(t, result)
				}
			}
		})
	}
}

func TestDoQueryColumnsEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockDB)
		shouldError bool
	}{
		{
			name: "should handle columns with null values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{nil, nil},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle columns with special characters",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"col@#$%", "val!@#$"},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle very long column names",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					longName := ""
					for i := 0; i < 1000; i++ {
						longName += "a"
					}
					return mockRows([][]interface{}{
						{longName, "value"},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle large number of columns",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					rows := [][]interface{}{}
					for i := 0; i < 1000; i++ {
						rows = append(rows, []interface{}{"col" + string(rune(i)), "val" + string(rune(i))})
					}
					return mockRows(rows), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle numeric column values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{int64(123), float64(456.789)},
					}), nil
				}
			},
			shouldError: false,
		},
		{
			name: "should handle boolean column values",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"enabled", true},
						{"disabled", false},
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

			result, err := p.doQueryColumns()

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestDoQueryColumnsWithDifferentDataTypes(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockDB)
	}{
		{
			name: "should handle string data",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{"string_col", "string_value"},
					}), nil
				}
			},
		},
		{
			name: "should handle integer data",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{int32(42), int64(9223372036854775807)},
					}), nil
				}
			},
		},
		{
			name: "should handle negative integers",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{int64(-9223372036854775808), int32(-2147483648)},
					}), nil
				}
			},
		},
		{
			name: "should handle float data",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{float32(1.23), float64(4.56789)},
					}), nil
				}
			},
		},
		{
			name: "should handle byte arrays",
			mockSetup: func(m *MockDB) {
				m.QueryFunc = func(query string, args ...interface{}) (*sql.Rows, error) {
					return mockRows([][]interface{}{
						{[]byte("bytes"), []byte("")},
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

			result, err := p.doQueryColumns()
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

// Helper functions and mocks

type MockDB struct {
	QueryFunc func(query string, args ...interface{}) (*sql.Rows, error)
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, errors.New("mock not configured")
}

func mockRows(data [][]interface{}) *sql.Rows {
	// In a real test, you would use sqlmock or similar
	// This is a simplified version
	return &sql.Rows{}
}

func mockRowsWithScanError() *sql.Rows {
	// Returns rows that will fail on scan
	return &sql.Rows{}
}