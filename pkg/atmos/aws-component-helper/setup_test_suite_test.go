package aws_component_helper

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupTestSuite(t *testing.T) {
	// Create a temporary test suite
	ts := &TestSuite{
		TempDir:      t.TempDir(),
		FixturesPath: "testdata/fixtures",
	}

	// Create test fixtures
	err := os.MkdirAll(ts.FixturesPath, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(ts.FixturesPath)

	// Test setup
	err = setupTestSuite(ts)
	assert.NoError(t, err)

	// Verify directories were created
	stateDir := filepath.Join(ts.TempDir, "state")
	assert.DirExists(t, stateDir)

	componentsDir := filepath.Join(ts.TempDir, "components", "terraform")
	assert.DirExists(t, componentsDir)
}

func TestCreateStateDir(t *testing.T) {
	tempDir := t.TempDir()

	err := createStateDir(tempDir)
	assert.NoError(t, err)

	stateDir := filepath.Join(tempDir, "state")
	assert.DirExists(t, stateDir)
}

func TestCreateTerraformComponentsDir(t *testing.T) {
	tempDir := t.TempDir()

	err := createTerraformComponentsDir(tempDir)
	assert.NoError(t, err)

	componentsDir := filepath.Join(tempDir, "components", "terraform")
	assert.DirExists(t, componentsDir)
}

func TestGetTestName(t *testing.T) {
	testName, err := getTestName()
	assert.NoError(t, err)
	assert.Contains(t, testName, "setup_test_suite_")
}

func TestReadOrCreateTestSuiteFile(t *testing.T) {
	t.Run("create new test suite file", func(t *testing.T) {
		// Clean up any existing test suite file
		os.Remove(testSuiteFile)

		ts := &TestSuite{}
		testName := "test_"

		result, err := readOrCreateTestSuiteFile(ts, testName)
		assert.NoError(t, err)
		assert.NotEmpty(t, result.RandomIdentifier)
		assert.NotEmpty(t, result.TempDir)

		// Verify file was created
		assert.FileExists(t, testSuiteFile)

		// Clean up
		os.Remove(testSuiteFile)
		os.RemoveAll(result.TempDir)
	})

	t.Run("read existing test suite file", func(t *testing.T) {
		// Create a test suite file
		existingTS := &TestSuite{
			RandomIdentifier: "test-seed",
			TempDir:          "test-dir",
		}
		data, err := json.MarshalIndent(existingTS, "", "  ")
		assert.NoError(t, err)
		err = os.WriteFile(testSuiteFile, data, 0644)
		assert.NoError(t, err)

		// Test reading the file
		ts := &TestSuite{}
		result, err := readOrCreateTestSuiteFile(ts, "test_")
		assert.NoError(t, err)
		assert.Equal(t, existingTS.RandomIdentifier, result.RandomIdentifier)
		assert.Equal(t, existingTS.TempDir, result.TempDir)

		// Clean up
		os.Remove(testSuiteFile)
	})
}
