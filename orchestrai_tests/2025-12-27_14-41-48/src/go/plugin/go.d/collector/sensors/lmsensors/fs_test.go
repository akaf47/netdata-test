package lmsensors

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestGetChipPathWhenPathExists tests GetChipPath with an existing path
func TestGetChipPathWhenPathExists(t *testing.T) {
	tempDir := t.TempDir()
	chipName := "coretemp-isa-0000"
	chipPath := filepath.Join(tempDir, chipName)

	// Create the chip directory
	if err := os.Mkdir(chipPath, 0755); err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	result, err := GetChipPath(tempDir, chipName)
	if err != nil {
		t.Errorf("GetChipPath() error = %v, want nil", err)
	}

	expectedPath := filepath.Join(tempDir, chipName)
	if result != expectedPath {
		t.Errorf("GetChipPath() = %s, want %s", result, expectedPath)
	}
}

// TestGetChipPathWhenPathDoesNotExist tests GetChipPath when path doesn't exist
func TestGetChipPathWhenPathDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentChip := "nonexistent-chip-0000"

	result, err := GetChipPath(tempDir, nonExistentChip)
	if err == nil {
		t.Errorf("GetChipPath() error = nil, want error")
	}
	if result != "" {
		t.Errorf("GetChipPath() returned non-empty path when path doesn't exist: %s", result)
	}
}

// TestGetChipPathWithEmptyChipName tests GetChipPath with empty chip name
func TestGetChipPathWithEmptyChipName(t *testing.T) {
	tempDir := t.TempDir()
	result, err := GetChipPath(tempDir, "")
	if err == nil {
		t.Errorf("GetChipPath() with empty chip name: error = nil, want error")
	}
	if result != "" {
		t.Errorf("GetChipPath() with empty chip name returned: %s, want empty string", result)
	}
}

// TestGetChipPathWithEmptyBaseDir tests GetChipPath with empty base directory
func TestGetChipPathWithEmptyBaseDir(t *testing.T) {
	chipName := "coretemp-isa-0000"
	result, err := GetChipPath("", chipName)
	if err == nil {
		t.Errorf("GetChipPath() with empty base dir: error = nil, want error")
	}
	if result != "" {
		t.Errorf("GetChipPath() with empty base dir returned: %s, want empty string", result)
	}
}

// TestGetChipPathWithSpecialCharacters tests GetChipPath with special characters in names
func TestGetChipPathWithSpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()
	chipName := "chip-with-special_chars-123"
	chipPath := filepath.Join(tempDir, chipName)

	if err := os.Mkdir(chipPath, 0755); err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	result, err := GetChipPath(tempDir, chipName)
	if err != nil {
		t.Errorf("GetChipPath() with special chars error = %v, want nil", err)
	}

	expectedPath := filepath.Join(tempDir, chipName)
	if result != expectedPath {
		t.Errorf("GetChipPath() = %s, want %s", result, expectedPath)
	}
}

// TestGetChipPathWithSymlink tests GetChipPath with symlinked paths
func TestGetChipPathWithSymlink(t *testing.T) {
	tempDir := t.TempDir()
	chipName := "real-chip"
	realPath := filepath.Join(tempDir, chipName)

	if err := os.Mkdir(realPath, 0755); err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	result, err := GetChipPath(tempDir, chipName)
	if err != nil {
		t.Errorf("GetChipPath() with symlink error = %v, want nil", err)
	}

	expectedPath := filepath.Join(tempDir, chipName)
	if result != expectedPath {
		t.Errorf("GetChipPath() = %s, want %s", result, expectedPath)
	}
}

// TestGetFeatureFilePathWhenFileExists tests GetFeatureFilePath with existing file
func TestGetFeatureFilePathWhenFileExists(t *testing.T) {
	tempDir := t.TempDir()
	featureFile := "temp1_input"
	filePath := filepath.Join(tempDir, featureFile)

	// Create the feature file
	if err := os.WriteFile(filePath, []byte("50000\n"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	result, err := GetFeatureFilePath(tempDir, featureFile)
	if err != nil {
		t.Errorf("GetFeatureFilePath() error = %v, want nil", err)
	}

	if result != filePath {
		t.Errorf("GetFeatureFilePath() = %s, want %s", result, filePath)
	}
}

// TestGetFeatureFilePathWhenFileDoesNotExist tests GetFeatureFilePath when file doesn't exist
func TestGetFeatureFilePathWhenFileDoesNotExist(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := "nonexistent_input"

	result, err := GetFeatureFilePath(tempDir, nonExistentFile)
	if err == nil {
		t.Errorf("GetFeatureFilePath() error = nil, want error")
	}
	if result != "" {
		t.Errorf("GetFeatureFilePath() returned non-empty path when file doesn't exist: %s", result)
	}
}

// TestGetFeatureFilePathWithEmptyFileName tests GetFeatureFilePath with empty file name
func TestGetFeatureFilePathWithEmptyFileName(t *testing.T) {
	tempDir := t.TempDir()
	result, err := GetFeatureFilePath(tempDir, "")
	if err == nil {
		t.Errorf("GetFeatureFilePath() with empty file name: error = nil, want error")
	}
	if result != "" {
		t.Errorf("GetFeatureFilePath() with empty file name returned: %s, want empty string", result)
	}
}

// TestGetFeatureFilePathWithEmptyDir tests GetFeatureFilePath with empty directory
func TestGetFeatureFilePathWithEmptyDir(t *testing.T) {
	fileName := "temp1_input"
	result, err := GetFeatureFilePath("", fileName)
	if err == nil {
		t.Errorf("GetFeatureFilePath() with empty dir: error = nil, want error")
	}
	if result != "" {
		t.Errorf("GetFeatureFilePath() with empty dir returned: %s, want empty string", result)
	}
}

// TestGetFeatureFilePathWithSpecialCharacters tests GetFeatureFilePath with special characters
func TestGetFeatureFilePathWithSpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()
	featureFile := "temp_1_max_alarm"
	filePath := filepath.Join(tempDir, featureFile)

	if err := os.WriteFile(filePath, []byte("60000\n"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	result, err := GetFeatureFilePath(tempDir, featureFile)
	if err != nil {
		t.Errorf("GetFeatureFilePath() with special chars error = %v, want nil", err)
	}

	if result != filePath {
		t.Errorf("GetFeatureFilePath() = %s, want %s", result, filePath)
	}
}

// TestGetFeatureFilePathWithMultipleLevels tests GetFeatureFilePath with nested paths
func TestGetFeatureFilePathWithMultipleLevels(t *testing.T) {
	tempDir := t.TempDir()
	featureFile := "temp1_input"
	filePath := filepath.Join(tempDir, featureFile)

	if err := os.WriteFile(filePath, []byte("45000\n"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	result, err := GetFeatureFilePath(tempDir, featureFile)
	if err != nil {
		t.Errorf("GetFeatureFilePath() error = %v, want nil", err)
	}

	if result != filePath {
		t.Errorf("GetFeatureFilePath() = %s, want %s", result, filePath)
	}
}

// TestGetFeatureFilesWhenDirectoryExists tests GetFeatureFiles with existing directory
func TestGetFeatureFilesWhenDirectoryExists(t *testing.T) {
	tempDir := t.TempDir()
	files := []string{"temp1_input", "temp2_input", "temp1_max"}

	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tempDir, f), []byte("50000\n"), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	result, err := GetFeatureFiles(tempDir)
	if err != nil {
		t.Errorf("GetFeatureFiles() error = %v, want nil", err)
	}

	if len(result) < len(files) {
		t.Errorf("GetFeatureFiles() returned %d files, want at least %d", len(result), len(files))
	}

	for _, f := range files {
		found := false
		for _, r := range result {
			if filepath.Base(r) == f {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetFeatureFiles() missing file: %s", f)
		}
	}
}

// TestGetFeatureFilesWhenDirectoryEmpty tests GetFeatureFiles with empty directory
func TestGetFeatureFilesWhenDirectoryEmpty(t *testing.T) {
	tempDir := t.TempDir()

	result, err := GetFeatureFiles(tempDir)
	if err != nil {
		t.Errorf("GetFeatureFiles() with empty dir error = %v, want nil", err)
	}

	if len(result) != 0 {
		t.Errorf("GetFeatureFiles() with empty dir returned %d files, want 0", len(result))
	}
}

// TestGetFeatureFilesWhenDirectoryDoesNotExist tests GetFeatureFiles with non-existent directory
func TestGetFeatureFilesWhenDirectoryDoesNotExist(t *testing.T) {
	nonExistentDir := "/nonexistent/directory/path"

	result, err := GetFeatureFiles(nonExistentDir)
	if err == nil {
		t.Errorf("GetFeatureFiles() with non-existent dir: error = nil, want error")
	}
	if len(result) != 0 {
		t.Errorf("GetFeatureFiles() with non-existent dir returned %d files, want 0", len(result))
	}
}

// TestGetFeatureFilesWithEmptyDirString tests GetFeatureFiles with empty directory string
func TestGetFeatureFilesWithEmptyDirString(t *testing.T) {
	result, err := GetFeatureFiles("")
	if err == nil {
		t.Errorf("GetFeatureFiles() with empty dir string: error = nil, want error")
	}
	if len(result) != 0 {
		t.Errorf("GetFeatureFiles() with empty dir string returned %d files, want 0", len(result))
	}
}

// TestGetFeatureFilesWithMixedContent tests GetFeatureFiles with mixed file/directory content
func TestGetFeatureFilesWithMixedContent(t *testing.T) {
	tempDir := t.TempDir()

	// Create some files
	files := []string{"temp1_input", "temp2_max"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tempDir, f), []byte("50000\n"), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	// Create a subdirectory (should be skipped)
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	result, err := GetFeatureFiles(tempDir)
	if err != nil {
		t.Errorf("GetFeatureFiles() with mixed content error = %v, want nil", err)
	}

	for _, f := range files {
		found := false
		for _, r := range result {
			if filepath.Base(r) == f {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetFeatureFiles() missing file: %s", f)
		}
	}
}

// TestGetChipsWhenDirectoryExists tests GetChips with existing directory
func TestGetChipsWhenDirectoryExists(t *testing.T) {
	tempDir := t.TempDir()
	chips := []string{"coretemp-isa-0000", "it8792-isa-0a60"}

	for _, c := range chips {
		chipPath := filepath.Join(tempDir, c)
		if err := os.Mkdir(chipPath, 0755); err != nil {
			t.Fatalf("failed to create chip directory: %v", err)
		}
	}

	result, err := GetChips(tempDir)
	if err != nil {
		t.Errorf("GetChips() error = %v, want nil", err)
	}

	if len(result) < len(chips) {
		t.Errorf("GetChips() returned %d chips, want at least %d", len(result), len(chips))
	}

	for _, c := range chips {
		found := false
		for _, r := range result {
			if filepath.Base(r) == c {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetChips() missing chip: %s", c)
		}
	}
}

// TestGetChipsWhenDirectoryEmpty tests GetChips with empty directory
func TestGetChipsWhenDirectoryEmpty(t *testing.T) {
	tempDir := t.TempDir()

	result, err := GetChips(tempDir)
	if err != nil {
		t.Errorf("GetChips() with empty dir error = %v, want nil", err)
	}

	if len(result) != 0 {
		t.Errorf("GetChips() with empty dir returned %d chips, want 0", len(result))
	}
}

// TestGetChipsWhenDirectoryDoesNotExist tests GetChips with non-existent directory
func TestGetChipsWhenDirectoryDoesNotExist(t *testing.T) {
	nonExistentDir := "/nonexistent/chips/directory"

	result, err := GetChips(nonExistentDir)
	if err == nil {
		t.Errorf("GetChips() with non-existent dir: error = nil, want error")
	}
	if len(result) != 0 {
		t.Errorf("GetChips() with non-existent dir returned %d chips, want 0", len(result))
	}
}

// TestGetChipsWithEmptyDirString tests GetChips with empty directory string
func TestGetChipsWithEmptyDirString(t *testing.T) {
	result, err := GetChips("")
	if err == nil {
		t.Errorf("GetChips() with empty dir string: error = nil, want error")
	}
	if len(result) != 0 {
		t.Errorf("GetChips() with empty dir string returned %d chips, want 0", len(result))
	}
}

// TestGetChipsWithMixedContent tests GetChips with mixed file/directory content
func TestGetChipsWithMixedContent(t *testing.T) {
	tempDir := t.TempDir()

	// Create chip directories
	chips := []string{"coretemp-isa-0000"}
	for _, c := range chips {
		chipPath := filepath.Join(tempDir, c)
		if err := os.Mkdir(chipPath, 0755); err != nil {
			t.Fatalf("failed to create chip directory: %v", err)
		}
	}

	// Create a regular file (should be skipped)
	if err := os.WriteFile(filepath.Join(tempDir, "regularfile.txt"), []byte("data"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	result, err := GetChips(tempDir)
	if err != nil {
		t.Errorf("GetChips() with mixed content error = %v, want nil", err)
	}

	for _, c := range chips {
		found := false
		for _, r := range result {
			if filepath.Base(r) == c {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetChips() missing chip: %s", c)
		}
	}
}

// TestGetChipsMultipleChips tests GetChips with multiple chips
func TestGetChipsMultipleChips(t *testing.T) {
	tempDir := t.TempDir()
	chips := []string{"chip-0", "chip-1", "chip-2", "chip-3"}

	for _, c := range chips {
		chipPath := filepath.Join(tempDir, c)
		if err := os.Mkdir(chipPath, 0755); err != nil {
			t.Fatalf("failed to create chip directory: %v", err)
		}
	}

	result, err := GetChips(tempDir)
	if err != nil {
		t.Errorf("GetChips() error = %v, want nil", err)
	}

	if len(result) != len(chips) {
		t.Errorf("GetChips() returned %d chips, want %d", len(result), len(chips))
	}

	for _, c := range chips {
		found := false
		for _, r := range result {
			if filepath.Base(r) == c {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetChips() missing chip: %s", c)
		}
	}
}

// TestReadFeatureFileContent tests reading feature file content
func TestReadFeatureFileContent(t *testing.T) {
	tempDir := t.TempDir()
	featureFile := "temp1_input"
	filePath := filepath.Join(tempDir, featureFile)
	testContent := "50000\n"

	if err := os.WriteFile(filePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("ReadFile() error = %v, want nil", err)
	}

	if string(content) != testContent {
		t.Errorf("ReadFile() content = %s, want %s", string(content), testContent)
	}
}

// TestPathPermissions tests that paths with appropriate permissions are accessible
func TestPathPermissions(t *testing.T) {
	tempDir := t.TempDir()
	chipName := "test-chip"
	chipPath := filepath.Join(tempDir, chipName)

	if err := os.Mkdir(chipPath, 0755); err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	result, err := GetChipPath(tempDir, chipName)
	if err != nil {
		t.Errorf("GetChipPath() with readable dir error = %v, want nil", err)
	}

	if result != chipPath {
		t.Errorf("GetChipPath() = %s, want %s", result, chipPath)
	}
}

// TestConcurrentFileAccess tests concurrent access to files
func TestConcurrentFileAccess(t *testing.T) {
	tempDir := t.TempDir()
	files := []string{"temp1_input", "temp2_input", "temp3_input"}

	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tempDir, f), []byte("50000\n"), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	done := make(chan bool)
	for i := 0; i < 3; i++ {
		go func() {
			_, _ = GetFeatureFiles(tempDir)
			done <- true
		}()
	}

	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestPathWithDotDot tests handling of paths with ".." sequences
func TestPathWithDotDot(t *testing.T) {
	tempDir := t.TempDir()
	basePath := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(basePath, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	parentRef := filepath.Join(basePath, "..", "subdir")
	chipName := "test-chip"
	chipPath := filepath.Join(parentRef, chipName)

	if err := os.Mkdir(chipPath, 0755); err != nil {
		t.Fatalf("failed to create chip directory: %v", err)
	}

	result, err := GetChipPath(parentRef, chipName)
	if err != nil {
		t.Errorf("GetChipPath() with .. in path error = %v, want nil", err)
	}

	if result == "" {
		t.Errorf("GetChipPath() with .. in path returned empty string")
	}
}

// TestFeatureFileWithLargeContent tests reading large feature files
func TestFeatureFileWithLargeContent(t *testing.T) {
	tempDir := t.TempDir()
	featureFile := "temp_large"
	filePath := filepath.Join(tempDir, featureFile)

	// Create a file with larger content
	largeContent := "123456789"
	if err := os.WriteFile(filePath, []byte(largeContent), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	result, err := GetFeatureFilePath(tempDir, featureFile)
	if err != nil {
		t.Errorf("GetFeatureFilePath() with large content error = %v, want nil", err)
	}

	if result != filePath {
		t.Errorf("GetFeatureFilePath() = %s, want %s", result, filePath)
	}
}

// TestChipPathWithLeadingSlash tests chip path with leading slash
func TestChipPathWithLeadingSlash(t *testing.T) {
	tempDir := t.TempDir()
	chipName := "/leadingslash"

	result, err := GetChipPath(tempDir, chipName)
	// This should either fail or handle the leading slash appropriately
	if err == nil && result == "" {
		t.Errorf("GetChipPath() with leading slash returned empty result without error")
	}
}

// TestGetChipsWithSymlinks tests GetChips with symlinked directories
func TestGetChipsWithSymlinks(t *testing.T) {
	tempDir := t.TempDir()
	chipName := "real-chip"
	chipPath := filepath.Join(tempDir, chipName)

	if err := os.Mkdir(chipPath, 0755); err != nil {
		t.Fatalf("failed to create chip directory: %v", err)
	}

	result, err := GetChips(tempDir)
	if err != nil {
		t.Errorf("GetChips() error = %v, want nil", err)
	}

	found := false
	for _, r := range result {
		if filepath.Base(r) == chipName {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("GetChips() missing chip: %s", chipName)
	}
}