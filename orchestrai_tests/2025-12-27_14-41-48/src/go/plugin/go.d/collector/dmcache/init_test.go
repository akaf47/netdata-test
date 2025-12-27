package dmcache

import (
	"context"
	"errors"
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/stm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock Proc interface for testing
type MockProc struct {
	mock.Mock
}

func (m *MockProc) KernelModules() (map[string]bool, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]bool), args.Error(1)
}

func (m *MockProc) Slabinfo() (stm.Slabinfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(stm.Slabinfo), args.Error(1)
}

func (m *MockProc) Meminfo() (stm.Meminfo, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(stm.Meminfo), args.Error(1)
}

// Mock file reader for testing
type MockFileReader struct {
	mock.Mock
}

func (m *MockFileReader) ReadFile(path string) ([]byte, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// TestNew tests the New function
func TestNew(t *testing.T) {
	dc := New()
	assert.NotNil(t, dc)
	assert.Equal(t, "dmcache", dc.Name())
}

// TestDmCache_Init_Success tests successful initialization
func TestDmCache_Init_Success(t *testing.T) {
	dc := New()
	
	assert.NoError(t, dc.Init(context.Background()))
	assert.NotNil(t, dc.charts)
}

// TestDmCache_Init_WithoutModuleLoaded tests initialization when dmcache module not loaded
func TestDmCache_Init_WithoutModuleLoaded(t *testing.T) {
	dc := New()
	
	// Mock proc to return empty modules
	mockProc := &MockProc{}
	mockProc.On("KernelModules").Return(map[string]bool{}, nil)
	dc.proc = mockProc
	
	err := dc.Init(context.Background())
	// Should handle gracefully - initialize but may skip module checks
	assert.NoError(t, err)
}

// TestDmCache_Init_InvalidPath tests initialization with invalid path
func TestDmCache_Init_InvalidPath(t *testing.T) {
	dc := New()
	dc.Source = "/nonexistent/path"
	
	err := dc.Init(context.Background())
	// Should handle invalid path gracefully
	assert.NoError(t, err) // Or specific error type based on implementation
}

// TestDmCache_Check tests the Check function
func TestDmCache_Check_Success(t *testing.T) {
	dc := New()
	err := dc.Init(context.Background())
	require.NoError(t, err)
	
	ok, err := dc.Check(context.Background())
	assert.True(t, ok)
	assert.NoError(t, err)
}

// TestDmCache_Check_WithoutInit tests Check before Init
func TestDmCache_Check_WithoutInit(t *testing.T) {
	dc := New()
	
	ok, err := dc.Check(context.Background())
	assert.False(t, ok)
	assert.Error(t, err)
}

// TestDmCache_Check_ContextCancelled tests Check with cancelled context
func TestDmCache_Check_ContextCancelled(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	ok, err := dc.Check(ctx)
	assert.False(t, ok)
	assert.Error(t, err)
}

// TestDmCache_Collect tests the Collect function
func TestDmCache_Collect_Success(t *testing.T) {
	dc := New()
	err := dc.Init(context.Background())
	require.NoError(t, err)
	
	mx := dc.Collect(context.Background())
	assert.NotNil(t, mx)
}

// TestDmCache_Collect_WithoutInit tests Collect before Init
func TestDmCache_Collect_WithoutInit(t *testing.T) {
	dc := New()
	
	mx := dc.Collect(context.Background())
	assert.Nil(t, mx)
}

// TestDmCache_Collect_WithContextCancelled tests Collect with cancelled context
func TestDmCache_Collect_WithContextCancelled(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	
	mx := dc.Collect(ctx)
	// Should return nil or empty result on cancelled context
	assert.Nil(t, mx)
}

// TestDmCache_Charts tests the Charts function
func TestDmCache_Charts(t *testing.T) {
	dc := New()
	err := dc.Init(context.Background())
	require.NoError(t, err)
	
	charts := dc.Charts()
	assert.NotNil(t, charts)
	assert.Greater(t, len(charts.(*module.Charts)), 0)
}

// TestDmCache_Name tests the Name function
func TestDmCache_Name(t *testing.T) {
	dc := New()
	assert.Equal(t, "dmcache", dc.Name())
}

// TestDmCache_ModuleDescriptor tests the ModuleDescriptor function
func TestDmCache_ModuleDescriptor(t *testing.T) {
	desc := ModuleDescriptor()
	assert.NotNil(t, desc)
	assert.Equal(t, "dmcache", desc.Name)
}

// TestDmCache_GetUpdateEvery tests update frequency configuration
func TestDmCache_GetUpdateEvery(t *testing.T) {
	dc := New()
	dc.UpdateEvery = 10
	
	assert.Equal(t, 10, dc.UpdateEvery)
}

// TestDmCache_GetUpdateEvery_Default tests default update frequency
func TestDmCache_GetUpdateEvery_Default(t *testing.T) {
	dc := New()
	assert.Greater(t, dc.UpdateEvery, 0)
}

// TestDmCache_Priority tests priority configuration
func TestDmCache_Priority(t *testing.T) {
	dc := New()
	dc.Priority = 100
	
	assert.Equal(t, 100, dc.Priority)
}

// TestDmCache_Cleanup tests cleanup on close
func TestDmCache_Cleanup(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	// Should not panic
	assert.NotPanics(t, func() {
		dc.Close()
	})
}

// TestDmCache_Init_Idempotent tests that Init can be called multiple times
func TestDmCache_Init_Idempotent(t *testing.T) {
	dc := New()
	
	err1 := dc.Init(context.Background())
	err2 := dc.Init(context.Background())
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

// TestDmCache_Collect_EmptyMetrics tests Collect returns empty metrics on failure
func TestDmCache_Collect_EmptyMetrics(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	mockProc := &MockProc{}
	mockProc.On("Slabinfo").Return(nil, errors.New("read error"))
	dc.proc = mockProc
	
	mx := dc.Collect(context.Background())
	// Should handle errors and return nil or empty metrics
	assert.Nil(t, mx)
}

// TestDmCache_Check_Slabinfo_Error tests Check with slabinfo read error
func TestDmCache_Check_Slabinfo_Error(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	mockProc := &MockProc{}
	mockProc.On("Slabinfo").Return(nil, errors.New("read error"))
	dc.proc = mockProc
	
	ok, err := dc.Check(context.Background())
	assert.False(t, ok)
	assert.Error(t, err)
}

// TestDmCache_Collect_PartialMetrics tests Collect with partial data
func TestDmCache_Collect_PartialMetrics(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	mockProc := &MockProc{}
	mockSlabinfo := stm.Slabinfo{
		"kmalloc-256": {
			Name:     "kmalloc-256",
			ObjSize:  256,
			NumObjs:  1000,
			NumActive: 500,
		},
	}
	mockProc.On("Slabinfo").Return(mockSlabinfo, nil)
	dc.proc = mockProc
	
	mx := dc.Collect(context.Background())
	assert.NotNil(t, mx)
}

// TestDmCache_Source_CustomPath tests with custom source path
func TestDmCache_Source_CustomPath(t *testing.T) {
	dc := New()
	dc.Source = "/custom/proc/path"
	
	err := dc.Init(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "/custom/proc/path", dc.Source)
}

// TestDmCache_Timeout_Configuration tests timeout configuration
func TestDmCache_Timeout(t *testing.T) {
	dc := New()
	dc.Timeout = 5
	
	assert.Equal(t, 5, dc.Timeout)
}

// TestDmCache_Init_WithTimeout tests initialization with timeout
func TestDmCache_Init_WithTimeout(t *testing.T) {
	dc := New()
	dc.Timeout = 1
	
	err := dc.Init(context.Background())
	assert.NoError(t, err)
}

// TestDmCache_Collect_MultipleMetrics tests collecting multiple metrics
func TestDmCache_Collect_MultipleMetrics(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	mockProc := &MockProc{}
	mockSlabinfo := stm.Slabinfo{
		"kmalloc-256": {
			Name:        "kmalloc-256",
			ObjSize:     256,
			NumObjs:     1000,
			NumActive:   500,
			ObjPerSlab:  16,
			PagesPerSlab: 1,
			TunedObjs:   0,
			TunedPages:  0,
			NumSlabs:    63,
			NumActiveSlabs: 62,
		},
		"kmalloc-512": {
			Name:        "kmalloc-512",
			ObjSize:     512,
			NumObjs:     2000,
			NumActive:   1000,
			ObjPerSlab:  8,
			PagesPerSlab: 1,
			TunedObjs:   0,
			TunedPages:  0,
			NumSlabs:    250,
			NumActiveSlabs: 248,
		},
	}
	mockProc.On("Slabinfo").Return(mockSlabinfo, nil)
	dc.proc = mockProc
	
	mx := dc.Collect(context.Background())
	assert.NotNil(t, mx)
}

// TestDmCache_Charts_NotNil tests that Charts returns non-nil value
func TestDmCache_Charts_NotNil(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	charts := dc.Charts()
	assert.NotNil(t, charts)
}

// TestDmCache_Check_NoError tests Check succeeds without error
func TestDmCache_Check_NoError(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	ok, err := dc.Check(context.Background())
	assert.NoError(t, err)
}

// TestDmCache_Collect_NotNil tests Collect returns non-nil
func TestDmCache_Collect_NotNil(t *testing.T) {
	dc := New()
	dc.Init(context.Background())
	
	mx := dc.Collect(context.Background())
	assert.NotNil(t, mx)
}