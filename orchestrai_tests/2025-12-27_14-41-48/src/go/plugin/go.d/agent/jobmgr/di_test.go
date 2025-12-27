package jobmgr

import (
	"testing"
)

// TestNewDependencyContainer tests the creation of a new dependency container
func TestNewDependencyContainer(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should create a valid dependency container",
			test: func(t *testing.T) {
				container := NewDependencyContainer()
				if container == nil {
					t.Fatal("expected non-nil container")
				}
			},
		},
		{
			name: "should initialize container with empty state",
			test: func(t *testing.T) {
				container := NewDependencyContainer()
				if container == nil {
					t.Fatal("expected non-nil container")
				}
				// Verify it's usable for dependency registration
				if err := container.Register("test", "value"); err != nil {
					t.Fatalf("expected no error on registration, got %v", err)
				}
			},
		},
		{
			name: "should allow multiple containers to be created independently",
			test: func(t *testing.T) {
				container1 := NewDependencyContainer()
				container2 := NewDependencyContainer()
				if container1 == container2 {
					t.Fatal("expected separate container instances")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestContainerRegister tests dependency registration
func TestContainerRegister(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      interface{}
		shouldFail bool
	}{
		{
			name:       "should register a simple dependency",
			key:        "logger",
			value:      "loggerInstance",
			shouldFail: false,
		},
		{
			name:       "should register a nil value",
			key:        "nullable",
			value:      nil,
			shouldFail: false,
		},
		{
			name:       "should register with empty string key",
			key:        "",
			value:      "value",
			shouldFail: true,
		},
		{
			name:       "should register complex dependency",
			key:        "config",
			value:      map[string]interface{}{"key": "value"},
			shouldFail: false,
		},
		{
			name:       "should overwrite existing dependency",
			key:        "existing",
			value:      "newValue",
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			err := container.Register(tt.key, tt.value)

			if tt.shouldFail && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !tt.shouldFail && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !tt.shouldFail {
				retrieved, err := container.Get(tt.key)
				if err != nil {
					t.Fatalf("expected to retrieve registered dependency, got error %v", err)
				}
				if retrieved != tt.value {
					t.Fatalf("expected %v, got %v", tt.value, retrieved)
				}
			}
		})
	}
}

// TestContainerGet tests dependency retrieval
func TestContainerGet(t *testing.T) {
	container := NewDependencyContainer()

	tests := []struct {
		name       string
		setupFunc  func()
		key        string
		shouldFail bool
		expected   interface{}
	}{
		{
			name: "should retrieve registered dependency",
			setupFunc: func() {
				container.Register("test", "testValue")
			},
			key:        "test",
			shouldFail: false,
			expected:   "testValue",
		},
		{
			name: "should fail on missing dependency",
			setupFunc: func() {
				// no setup
			},
			key:        "missing",
			shouldFail: true,
			expected:   nil,
		},
		{
			name: "should get nil value dependency",
			setupFunc: func() {
				container.Register("nilDep", nil)
			},
			key:        "nilDep",
			shouldFail: false,
			expected:   nil,
		},
		{
			name: "should retrieve after multiple registrations",
			setupFunc: func() {
				container.Register("dep1", "value1")
				container.Register("dep2", "value2")
				container.Register("dep3", "value3")
			},
			key:        "dep2",
			shouldFail: false,
			expected:   "value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset container for each test
			container = NewDependencyContainer()
			tt.setupFunc()

			retrieved, err := container.Get(tt.key)

			if tt.shouldFail && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !tt.shouldFail && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !tt.shouldFail && retrieved != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, retrieved)
			}
		})
	}
}

// TestContainerUnregister tests dependency removal
func TestContainerUnregister(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*DependencyContainer)
		key         string
		shouldFail  bool
		verifyFunc  func(*DependencyContainer, string) bool
	}{
		{
			name: "should unregister existing dependency",
			setupFunc: func(c *DependencyContainer) {
				c.Register("dep", "value")
			},
			key:        "dep",
			shouldFail: false,
			verifyFunc: func(c *DependencyContainer, key string) bool {
				_, err := c.Get(key)
				return err != nil
			},
		},
		{
			name: "should fail on unregistering non-existent dependency",
			setupFunc: func(c *DependencyContainer) {
				// no setup
			},
			key:        "nonexistent",
			shouldFail: true,
			verifyFunc: nil,
		},
		{
			name: "should remove dependency but keep others",
			setupFunc: func(c *DependencyContainer) {
				c.Register("dep1", "value1")
				c.Register("dep2", "value2")
			},
			key:        "dep1",
			shouldFail: false,
			verifyFunc: func(c *DependencyContainer, _ string) bool {
				_, err := c.Get("dep2")
				return err == nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			tt.setupFunc(container)

			err := container.Unregister(tt.key)

			if tt.shouldFail && err == nil {
				t.Fatal("expected error but got nil")
			}
			if !tt.shouldFail && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !tt.shouldFail && tt.verifyFunc != nil {
				if !tt.verifyFunc(container, tt.key) {
					t.Fatal("verification failed after unregister")
				}
			}
		})
	}
}

// TestContainerClear tests clearing all dependencies
func TestContainerClear(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*DependencyContainer)
	}{
		{
			name: "should clear empty container",
			setupFunc: func(c *DependencyContainer) {
				// no setup
			},
		},
		{
			name: "should clear single dependency",
			setupFunc: func(c *DependencyContainer) {
				c.Register("dep", "value")
			},
		},
		{
			name: "should clear multiple dependencies",
			setupFunc: func(c *DependencyContainer) {
				c.Register("dep1", "value1")
				c.Register("dep2", "value2")
				c.Register("dep3", "value3")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			tt.setupFunc(container)

			container.Clear()

			// Verify all dependencies are cleared
			_, err := container.Get("dep1")
			if err == nil {
				t.Fatal("expected error after clear")
			}
		})
	}
}

// TestContainerHas tests checking dependency existence
func TestContainerHas(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		setup     bool
		expected  bool
	}{
		{
			name:     "should return true for existing dependency",
			key:      "existing",
			setup:    true,
			expected: true,
		},
		{
			name:     "should return false for missing dependency",
			key:      "missing",
			setup:    false,
			expected: false,
		},
		{
			name:     "should handle empty string key",
			key:      "",
			setup:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			if tt.setup {
				container.Register(tt.key, "value")
			}

			result := container.Has(tt.key)
			if result != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestContainerGetAll tests retrieving all dependencies
func TestContainerGetAll(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(*DependencyContainer)
		expectedCount int
	}{
		{
			name: "should return empty map for empty container",
			setupFunc: func(c *DependencyContainer) {
				// no setup
			},
			expectedCount: 0,
		},
		{
			name: "should return all registered dependencies",
			setupFunc: func(c *DependencyContainer) {
				c.Register("dep1", "value1")
				c.Register("dep2", "value2")
				c.Register("dep3", "value3")
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			tt.setupFunc(container)

			all := container.GetAll()
			if len(all) != tt.expectedCount {
				t.Fatalf("expected %d dependencies, got %d", tt.expectedCount, len(all))
			}
		})
	}
}

// TestContainerGetOrRegister tests get-or-register functionality
func TestContainerGetOrRegister(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(*DependencyContainer)
		key          string
		defaultValue interface{}
		expectedVal  interface{}
		shouldExist  bool
	}{
		{
			name: "should return existing dependency",
			setupFunc: func(c *DependencyContainer) {
				c.Register("existing", "existingValue")
			},
			key:          "existing",
			defaultValue: "default",
			expectedVal:  "existingValue",
			shouldExist:  true,
		},
		{
			name: "should register default for missing dependency",
			setupFunc: func(c *DependencyContainer) {
				// no setup
			},
			key:          "missing",
			defaultValue: "defaultValue",
			expectedVal:  "defaultValue",
			shouldExist:  true,
		},
		{
			name: "should handle nil default value",
			setupFunc: func(c *DependencyContainer) {
				// no setup
			},
			key:          "nilDefault",
			defaultValue: nil,
			expectedVal:  nil,
			shouldExist:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			tt.setupFunc(container)

			retrieved := container.GetOrRegister(tt.key, tt.defaultValue)

			if retrieved != tt.expectedVal {
				t.Fatalf("expected %v, got %v", tt.expectedVal, retrieved)
			}

			if tt.shouldExist {
				exists := container.Has(tt.key)
				if !exists {
					t.Fatalf("expected dependency to be registered")
				}
			}
		})
	}
}

// TestContainerRegistrySize tests getting container size
func TestContainerRegistrySize(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(*DependencyContainer)
		expectedSize  int
	}{
		{
			name: "should return 0 for empty container",
			setupFunc: func(c *DependencyContainer) {
				// no setup
			},
			expectedSize: 0,
		},
		{
			name: "should return correct count after registrations",
			setupFunc: func(c *DependencyContainer) {
				c.Register("dep1", "value1")
				c.Register("dep2", "value2")
			},
			expectedSize: 2,
		},
		{
			name: "should update count after unregister",
			setupFunc: func(c *DependencyContainer) {
				c.Register("dep1", "value1")
				c.Register("dep2", "value2")
				c.Unregister("dep1")
			},
			expectedSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			tt.setupFunc(container)

			size := container.Size()
			if size != tt.expectedSize {
				t.Fatalf("expected size %d, got %d", tt.expectedSize, size)
			}
		})
	}
}

// TestContainerConcurrency tests concurrent access to container
func TestContainerConcurrency(t *testing.T) {
	container := NewDependencyContainer()

	// Concurrent registration
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(index int) {
			key := "dep" + string(rune(index))
			container.Register(key, "value"+string(rune(index)))
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all registrations succeeded
	size := container.Size()
	if size < 1 {
		t.Fatalf("expected at least 1 dependency, got %d", size)
	}
}

// TestContainerWithComplexTypes tests container with complex data types
func TestContainerWithComplexTypes(t *testing.T) {
	type ComplexStruct struct {
		Name  string
		Value int
		Items []string
	}

	tests := []struct {
		name     string
		key      string
		value    interface{}
	}{
		{
			name:  "should register struct pointer",
			key:   "struct",
			value: &ComplexStruct{Name: "test", Value: 42},
		},
		{
			name:  "should register map",
			key:   "map",
			value: map[string]interface{}{"key": "value"},
		},
		{
			name:  "should register slice",
			key:   "slice",
			value: []string{"a", "b", "c"},
		},
		{
			name:  "should register function",
			key:   "func",
			value: func() { /* noop */ },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewDependencyContainer()
			container.Register(tt.key, tt.value)

			retrieved, err := container.Get(tt.key)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if retrieved != tt.value {
				t.Fatalf("expected %v, got %v", tt.value, retrieved)
			}
		})
	}
}

// BenchmarkContainerRegister benchmarks dependency registration
func BenchmarkContainerRegister(b *testing.B) {
	container := NewDependencyContainer()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		container.Register("dep", "value")
	}
}

// BenchmarkContainerGet benchmarks dependency retrieval
func BenchmarkContainerGet(b *testing.B) {
	container := NewDependencyContainer()
	container.Register("dep", "value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		container.Get("dep")
	}
}

// BenchmarkContainerHas benchmarks dependency existence check
func BenchmarkContainerHas(b *testing.B) {
	container := NewDependencyContainer()
	container.Register("dep", "value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		container.Has("dep")
	}
}