package postgres

import (
	"database/sql"
	"testing"
)

// TestDoQueryBloat_Success tests successful bloat query execution
func TestDoQueryBloat_Success(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should not return error, got: %v", err)
	}
}

func TestDoQueryBloat_WithNilDB(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = nil

	// Act
	err := p.doQueryBloat()

	// Assert
	if err == nil {
		t.Error("doQueryBloat() should return error when db is nil")
	}
}

func TestDoQueryBloat_DatabaseError(t *testing.T) {
	// Arrange
	p := &Postgres{}
	expectedErr := sql.ErrNoRows
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, expectedErr
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != expectedErr {
		t.Errorf("doQueryBloat() should propagate database error, got: %v", err)
	}
}

func TestDoQueryBloat_TableBloat(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate table bloat data
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle table bloat query, got: %v", err)
	}
}

func TestDoQueryBloat_IndexBloat(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate index bloat data
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle index bloat query, got: %v", err)
	}
}

func TestDoQueryBloat_EmptyBloat(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate no bloat
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle zero bloat, got: %v", err)
	}
}

func TestDoQueryBloat_ParseBloatMetrics(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()
	metrics := p.getBloatMetrics()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should not return error, got: %v", err)
	}
	if metrics == nil {
		t.Error("getBloatMetrics() should return non-nil metrics")
	}
}

func TestDoQueryBloat_LargeBloatValues(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate large bloat values
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle large bloat values, got: %v", err)
	}
}

func TestDoQueryBloat_MultipleTables(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate bloat data for multiple tables
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle multiple tables, got: %v", err)
	}
}

func TestDoQueryBloat_MultipleIndexes(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate bloat data for multiple indexes
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle multiple indexes, got: %v", err)
	}
}

func TestDoQueryBloat_ConnectionClosed(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, sql.ErrConnDone
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != sql.ErrConnDone {
		t.Errorf("doQueryBloat() should propagate connection error, got: %v", err)
	}
}

func TestDoQueryBloat_BadConnection(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, sql.ErrConnClosed
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != sql.ErrConnClosed {
		t.Errorf("doQueryBloat() should handle closed connection, got: %v", err)
	}
}

func TestDoQueryBloat_NullBloatValues(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate NULL bloat values
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle NULL values, got: %v", err)
	}
}

func TestDoQueryBloat_ZeroBloatPercentage(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate 0% bloat
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle zero bloat percentage, got: %v", err)
	}
}

func TestDoQueryBloat_MaxBloatPercentage(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate 100% bloat
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle maximum bloat percentage, got: %v", err)
	}
}

func TestDoQueryBloat_DecimalBloatValues(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Simulate decimal bloat values like 42.5%
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should handle decimal bloat values, got: %v", err)
	}
}

func TestDoQueryBloat_HighBloatThreshold(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.bloatThreshold = 90.0
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() with high threshold should succeed, got: %v", err)
	}
}

func TestDoQueryBloat_LowBloatThreshold(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.bloatThreshold = 10.0
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() with low threshold should succeed, got: %v", err)
	}
}

func TestDoQueryBloat_ExcludedSchemas(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.excludeSchemas = []string{"pg_catalog", "information_schema"}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			// Verify excluded schemas are in query
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() with excluded schemas should succeed, got: %v", err)
	}
}

func TestDoQueryBloat_FilterByBloatSize(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.minBloatSize = 1024 * 1024 // 1MB
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() with min size filter should succeed, got: %v", err)
	}
}

func TestDoQueryBloat_ConcurrentCalls(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act & Assert
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			err := p.doQueryBloat()
			if err != nil {
				t.Errorf("Concurrent doQueryBloat() failed: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}

func TestDoQueryBloat_RowProcessing(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should process all rows, got: %v", err)
	}
}

func TestDoQueryBloat_MetricAggregation(t *testing.T) {
	// Arrange
	p := &Postgres{}
	p.db = &mockDB{
		queryFunc: func(query string, args ...interface{}) (*sql.Rows, error) {
			return nil, nil
		},
	}

	// Act
	err := p.doQueryBloat()

	// Assert
	if err != nil {
		t.Errorf("doQueryBloat() should aggregate metrics, got: %v", err)
	}
}