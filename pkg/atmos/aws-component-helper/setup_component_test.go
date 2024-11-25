package aws_component_helper

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupComponentUnderTest(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test source directory with some content
	srcDir, err := os.MkdirTemp("", "src-*")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Create a test file in the source directory
	testFile := filepath.Join(srcDir, "test.tf")
	if err := os.WriteFile(testFile, []byte("test content"), 0666); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ts := &TestSuite{
		TempDir:         tempDir,
		ComponentName:   "test-component",
		ComponentSrcPath: srcDir,
	}

	if err := setupComponentUnderTest(ts); err != nil {
		t.Errorf("setupComponentUnderTest failed: %v", err)
	}

	// Verify the component directory was created
	expectedDir := filepath.Join(tempDir, "components", "terraform", "test-component")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Error("Component directory was not created")
	}

	// Verify the test file was copied
	expectedFile := filepath.Join(expectedDir, "test.tf")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Error("Test file was not copied")
	}
}

func TestCreateComponentDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ts := &TestSuite{
		TempDir:       tempDir,
		ComponentName: "test-component",
	}

	if err := createComponentDir(ts); err != nil {
		t.Errorf("createComponentDir failed: %v", err)
	}

	expectedDir := filepath.Join(tempDir, "components", "terraform", "test-component")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Error("Component directory was not created")
	}
}

func TestCopyComponentSrcToTempDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test source directory with some content
	srcDir, err := os.MkdirTemp("", "src-*")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Create a test file in the source directory
	testFile := filepath.Join(srcDir, "test.tf")
	if err := os.WriteFile(testFile, []byte("test content"), 0666); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create the test suite before trying to use it
	ts := &TestSuite{
		TempDir:         tempDir,
		ComponentName:   "test-component",
		ComponentSrcPath: srcDir,
	}

	// Create the component directory first
	componentDir := filepath.Join(tempDir, "components", "terraform", "test-component")
	if err := os.MkdirAll(componentDir, 0777); err != nil {
		t.Fatalf("Failed to create component dir: %v", err)
	}

	if err := copyComponentSrcToTempDir(ts); err != nil {
		t.Errorf("copyComponentSrcToTempDir failed: %v", err)
	}

	// Verify the test file was copied
	expectedFile := filepath.Join(componentDir, "test.tf")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Error("Test file was not copied")
	}
} 
