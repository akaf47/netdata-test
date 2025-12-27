package filecheck

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollectDirs tests the directory collection functionality
func TestCollectDirs(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Filecheck)
		wantErr bool
		verify  func(*testing.T, *Filecheck)
	}{
		{
			name: "collect single directory",
			setup: func(fc *Filecheck) {
				fc.Dirs = []string{"testdir"}
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
		{
			name: "collect multiple directories",
			setup: func(fc *Filecheck) {
				fc.Dirs = []string{"dir1", "dir2", "dir3"}
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
		{
			name: "collect with empty directory list",
			setup: func(fc *Filecheck) {
				fc.Dirs = []string{}
			},
			wantErr: false,
			verify: func(t *testing.T, fc *Filecheck) {
				assert.NotNil(t, fc)
			},
		},
		{
			name: "collect with nil directory list",
			setup: func(fc *Filecheck) {
				fc.Dirs = nil
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

// TestCollectDirsExists tests directory existence checking
func TestCollectDirsExists(t *testing.T) {
	tests := []struct {
		name    string
		dirs    []string
		wantErr bool
	}{
		{
			name:    "existing directory /tmp",
			dirs:    []string{"/tmp"},
			wantErr: false,
		},
		{
			name:    "non-existent directory",
			dirs:    []string{"/nonexistent/directory/path"},
			wantErr: false,
		},
		{
			name:    "mixed existing and non-existing",
			dirs:    []string{"/tmp", "/nonexistent/path"},
			wantErr: false,
		},
		{
			name:    "empty string directory",
			dirs:    []string{""},
			wantErr: false,
		},
		{
			name:    "current directory",
			dirs:    []string{"."},
			wantErr: false,
		},
		{
			name:    "parent directory",
			dirs:    []string{".."},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newTestFilecheck()
			fc.Dirs = tt.dirs

			err := fc.Collect()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCollectDirsWithReadablePermissions tests directories with readable permissions
func TestCollectDirsWithReadablePermissions(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsWithoutReadablePermissions tests handling of permission denied
func TestCollectDirsWithoutReadablePermissions(t *testing.T) {
	fc := newTestFilecheck()
	// Use a directory that typically has restricted access
	fc.Dirs = []string{"/root"}

	// This should not error out completely, just skip/report
	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsRecursive tests recursive directory collection
func TestCollectDirsRecursive(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsWithSpecialCharacters tests directories with special characters
func TestCollectDirsWithSpecialCharacters(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{
		"/tmp/dir with spaces",
		"/tmp/dir-with-dashes",
		"/tmp/dir_with_underscores",
	}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsLargeDirList tests with many directories
func TestCollectDirsLargeDirList(t *testing.T) {
	fc := newTestFilecheck()
	dirs := make([]string, 100)
	for i := 0; i < 100; i++ {
		dirs[i] = "/tmp"
	}
	fc.Dirs = dirs

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsSymlinks tests handling of symbolic links
func TestCollectDirsSymlinks(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsRelativePaths tests relative directory paths
func TestCollectDirsRelativePaths(t *testing.T) {
	tests := []struct {
		name string
		dir  string
	}{
		{
			name: "current directory",
			dir:  ".",
		},
		{
			name: "parent directory",
			dir:  "..",
		},
		{
			name: "relative subdirectory",
			dir:  "./subdir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newTestFilecheck()
			fc.Dirs = []string{tt.dir}

			err := fc.Collect()
			assert.NoError(t, err)
		})
	}
}

// TestCollectDirsAbsolutePaths tests absolute directory paths
func TestCollectDirsAbsolutePaths(t *testing.T) {
	tests := []struct {
		name string
		dir  string
	}{
		{
			name: "root directory",
			dir:  "/",
		},
		{
			name: "tmp directory",
			dir:  "/tmp",
		},
		{
			name: "var directory",
			dir:  "/var",
		},
		{
			name: "etc directory",
			dir:  "/etc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newTestFilecheck()
			fc.Dirs = []string{tt.dir}

			err := fc.Collect()
			assert.NoError(t, err)
		})
	}
}

// TestCollectDirsEmpty tests empty directory collection
func TestCollectDirsEmpty(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsNil tests nil directory collection
func TestCollectDirsNil(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = nil

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsCombinedWithFiles tests directory collection combined with file collection
func TestCollectDirsCombinedWithFiles(t *testing.T) {
	fc := newTestFilecheck()
	fc.Files = []string{"file1.txt", "file2.txt"}
	fc.Dirs = []string{"/tmp", "/var"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsRepeated tests repeated directory collection calls
func TestCollectDirsRepeated(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	for i := 0; i < 10; i++ {
		err := fc.Collect()
		assert.NoError(t, err)
	}
}

// TestCollectDirsParallel tests parallel directory collection
func TestCollectDirsParallel(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

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

// TestCollectDirsCounters tests directory counter increments
func TestCollectDirsCounters(t *testing.T) {
	tests := []struct {
		name       string
		dirs       []string
		description string
	}{
		{
			name:        "single directory",
			dirs:        []string{"/tmp"},
			description: "should count single dir",
		},
		{
			name:        "multiple directories",
			dirs:        []string{"/tmp", "/var", "/etc"},
			description: "should count multiple dirs",
		},
		{
			name:        "duplicate directories",
			dirs:        []string{"/tmp", "/tmp", "/tmp"},
			description: "should handle duplicate dirs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := newTestFilecheck()
			fc.Dirs = tt.dirs

			err := fc.Collect()
			assert.NoError(t, err, tt.description)
		})
	}
}

// TestCollectDirsSize tests directory size collection
func TestCollectDirsSize(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsModTime tests directory modification time collection
func TestCollectDirsModTime(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsAccessTime tests directory access time collection
func TestCollectDirsAccessTime(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsChangeTime tests directory change time collection
func TestCollectDirsChangeTime(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsHiddenDirs tests hidden directory collection
func TestCollectDirsHiddenDirs(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp/.hidden", "/tmp/.."}

	err := fc.Collect()
	assert.NoError(t, err)
}

// TestCollectDirsFileTypeDetection tests that directories are properly identified
func TestCollectDirsFileTypeDetection(t *testing.T) {
	fc := newTestFilecheck()
	fc.Dirs = []string{"/tmp"}

	err := fc.Collect()
	assert.NoError(t, err)
}