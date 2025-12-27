package filecheck

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollect tests the main Collect function
func TestCollect(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Filecheck)
		wantErr bool
		verify  func(*testing.T, *Filecheck)
	}{
		{
			name: "successful collection with files and directories",
			setup: func(fc *Filecheck) {
				fc.Files = []string{"testfile.txt"}
				fc.Dirs = []string{"testdir"}
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
		{
			name: "collection with empty files and dirs",
			setup: func(fc *Filecheck) {
				fc.Files = []string{}
				fc.Dirs = []string{}
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
		{
			name: "collection with nil files and dirs",
			setup: func(fc *Filecheck) {
				fc.Files = nil
				fc.Dirs = nil
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
		{
			name: "collection with multiple files",
			setup: func(fc *Filecheck) {
				fc.Files = []string{"file1.txt", "file2.txt", "file3.txt"}
				fc.Dirs = []string{}
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
		{
			name: "collection with multiple directories",
			setup: func(fc *Filecheck) {
				fc.Files = []string{}
				fc.Dirs = []string{"dir1", "dir2", "dir3"}
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newTestFilecheck()
			tt.setup(fc)

			err := fc.Collect()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			tt.verify(t, fc)
		})
	}
}

// TestCollectWithFileCheckState tests collection with various file states
func TestCollectWithFileCheckState(t *testing.T) {
	tests := []struct {
		name      string
		fileState map[string]bool // true if file exists
		wantErr   bool
	}{
		{
			name:      "all files exist",
			fileState: map[string]bool{"file1.txt": true, "file2.txt": true},
			wantErr:   false,
		},
		{
			name:      "no files exist",
			fileState: map[string]bool{"file1.txt": false, "file2.txt": false},
			wantErr:   false,
		},
		{
			name:      "mixed file existence",
			fileState: map[string]bool{"file1.txt": true, "file2.txt": false},
			wantErr:   false,
		},
		{
			name:      "single file exists",
			fileState: map[string]bool{"file1.txt": true},
			wantErr:   false,
		},
		{
			name:      "single file missing",
			fileState: map[string]bool{"file1.txt": false},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newTestFilecheck()
			files := make([]string, 0, len(tt.fileState))
			for f := range tt.fileState {
				files = append(files, f)
			}
			fc.Files = files

			err := fc.Collect()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCollectParallelCalls tests concurrent collect calls
func TestCollectParallelCalls(t *testing.T) {
	fc := newTestFilecheck()
	fc.Files = []string{"file1.txt"}
	fc.Dirs = []string{"dir1"}

	done := make(chan bool)

	for i := 0; i < 5; i++ {
		go func() {
			err := fc.Collect()
			assert.NoError(t, err)
			done <- true
		}()
	}

	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestCollectWithLargeFileList tests collection with many files
func TestCollectWithLargeFileList(t *testing.T) {
	fc := newTestFilecheck()
	files := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		files[i] = "file" + string(rune(i))
	}
	fc.Files = files

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectWithLargeDirList tests collection with many directories
func TestCollectWithLargeDirList(t *testing.T) {
	fc := newTestFilecheck()
	dirs := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		dirs[i] = "dir" + string(rune(i))
	}
	fc.Dirs = dirs

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectTiming tests that collect completes in reasonable time
func TestCollectTiming(t *testing.T) {
	fc := newTestFilecheck()
	fc.Files = []string{"file1.txt", "file2.txt"}
	fc.Dirs = []string{"dir1", "dir2"}

	start := time.Now()
	err := fc.Collect()
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.True(t, elapsed < 5*time.Second, "Collect took too long: %v", elapsed)
}

// TestCollectRepeated tests multiple successive collect calls
func TestCollectRepeated(t *testing.T) {
	fc := newTestFilecheck()
	fc.Files = []string{"file1.txt"}
	fc.Dirs = []string{"dir1"}

	for i := 0; i < 10; i++ {
		err := fc.Collect()
		assert.NoError(t, err)
	}
}

// TestCollectWithSpecialCharacters tests files with special characters
func TestCollectWithSpecialCharacters(t *testing.T) {
	fc := newTestFilecheck()
	fc.Files = []string{
		"file with spaces.txt",
		"file-with-dashes.txt",
		"file_with_underscores.txt",
		"file.multiple.dots.txt",
	}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectWithPaths tests files with path elements
func TestCollectWithPaths(t *testing.T) {
	fc := newTestFilecheck()
	fc.Files = []string{
		"/etc/hosts",
		"/var/log/syslog",
		"./relative/path/file.txt",
		"../parent/file.txt",
	}

	err := fc.Collect()
	assert.NoError(t, err)
}

// Helper function to create test Filecheck instance
func newTestFilecheck() *Filecheck {
	fc := &Filecheck{
		Files: []string{},
		Dirs:  []string{},
	}
	return fc
}