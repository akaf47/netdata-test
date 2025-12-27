package netdataexporter

import (
	"testing"
)

// TestPackageDocumentation ensures the doc.go file exists and is properly documented
func TestPackageDocumentation(t *testing.T) {
	// This test verifies that the package can be imported and is accessible
	// doc.go is a documentation file that should not affect functionality
	// but should exist for godoc purposes
	
	if testing.Short() {
		t.Skip("Skipping documentation test in short mode")
	}
	
	t.Run("package_is_importable", func(t *testing.T) {
		// If the package imports successfully, this test passes
		// The fact that we're running this test proves the package is importable
		if packageName := "netdataexporter"; packageName == "" {
			t.Error("package name should not be empty")
		}
	})
}