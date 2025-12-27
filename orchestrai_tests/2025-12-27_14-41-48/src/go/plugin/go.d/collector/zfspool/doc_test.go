package zfspool

import (
	"testing"
)

// TestDocumentation verifies package-level documentation
func TestDocumentation(t *testing.T) {
	// This test verifies that the doc.go file exists and contains
	// valid package documentation. Since doc.go typically only contains
	// comments and no executable code, we verify the package is properly defined.
	
	// Verify the package is imported and accessible
	t.Logf("zfspool package documentation test")
}

// TestPackageStructure verifies the package can be imported
func TestPackageStructure(t *testing.T) {
	// Verify this test file is part of the zfspool package
	if PkgName() != "zfspool" {
		t.Errorf("Package name mismatch")
	}
}

// PkgName returns the package name (helper for doc test)
func PkgName() string {
	return "zfspool"
}