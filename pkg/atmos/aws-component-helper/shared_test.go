package aws_component_helper

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMatchFilter tests the matchFilter function under various input scenarios
func TestMatchFilter(t *testing.T) {
	// Case: Invalid test name without a '/' (should return an error)
	t.Run("wrong", func(t *testing.T) {
		matchRegexp := ""
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance")
		assert.Error(t, err)    // Expect an error due to invalid format
		assert.False(t, result) // Result should be false
	})

	// Case: Match the default suite
	t.Run("default", func(t *testing.T) {
		matchRegexp := ""
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default")
		assert.NoError(t, err) // No error expected
		assert.True(t, result) // Should match the default suite
	})

	// Case: Match a specific suite
	t.Run("suite", func(t *testing.T) {
		matchRegexp := "default"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default")
		assert.NoError(t, err) // No error expected
		assert.True(t, result) // Should match the suite
	})

	// Case: Match all tests within a suite
	t.Run("all tests", func(t *testing.T) {
		matchRegexp := "default"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default/test1")
		assert.NoError(t, err) // No error expected
		assert.True(t, result) // Should match all tests under "default"
	})

	// Case: Match a specific test within a suite
	t.Run("one test tests", func(t *testing.T) {
		matchRegexp := "default/test1"
		matchSuiteAndTest = &matchRegexp
		result, err := matchFilter("TestAcceptance/default/test1")
		assert.NoError(t, err) // No error expected
		assert.True(t, result) // Should match "test1"

		result, err = matchFilter("TestAcceptance/default/test2")
		assert.NoError(t, err)  // No error expected
		assert.False(t, result) // Should not match "test2"
	})

	// Case: Match multilevel subtests
	t.Run("multilevel subtests", func(t *testing.T) {
		matchRegexp := "default/test1"
		matchSuiteAndTest = &matchRegexp

		// Matches subtests under "test1"
		result, err := matchFilter("TestAcceptance/default/test1/subtest1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test1/subtest2")
		assert.NoError(t, err)
		assert.True(t, result)

		// Does not match subtests under "test2"
		result, err = matchFilter("TestAcceptance/default/test2/subtest1")
		assert.NoError(t, err)
		assert.False(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	// Case: Filter longer than the test hierarchy
	t.Run("filter longer then test", func(t *testing.T) {
		matchRegexp := "default/test1/subtest1"
		matchSuiteAndTest = &matchRegexp

		// Matches as filter is more specific than the test path
		result, err := matchFilter("TestAcceptance/default/test1")
		assert.NoError(t, err)
		assert.True(t, result)

		// Does not match unrelated paths
		result, err = matchFilter("TestAcceptance/default/test2/subtest1")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	// Case: Use regex for matching
	t.Run("regexp", func(t *testing.T) {
		matchRegexp := "default/.*/subtest1"
		matchSuiteAndTest = &matchRegexp

		// Matches "subtest1" in various hierarchies
		result, err := matchFilter("TestAcceptance/default/test1/subtest1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test1/subtest2")
		assert.NoError(t, err)
		assert.False(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest1")
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = matchFilter("TestAcceptance/default/test2/subtest2")
		assert.NoError(t, err)
		assert.False(t, result)
	})
}

// TestCreateStateDir tests the creation of a state directory
func TestCreateStateDir(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Attempt to create a "state" subdirectory
	err := createDir(tempDir, "state")
	assert.NoError(t, err) // Ensure no errors occurred

	// Verify that the directory exists
	stateDir := filepath.Join(tempDir, "state")
	assert.DirExists(t, stateDir)
}

// TestCreateTerraformComponentsDir tests the creation of Terraform components directory structure
func TestCreateTerraformComponentsDir(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Attempt to create the Terraform components directory structure
	err := createTerraformComponentsDir(tempDir)
	assert.NoError(t, err) // Ensure no errors occurred

	// Verify that the directory structure exists
	componentsDir := filepath.Join(tempDir, "components", "terraform")
	assert.DirExists(t, componentsDir)
}

// TestGetTestName tests the retrieval of the current test's name
func TestGetTestName(t *testing.T) {
	// Attempt to retrieve the test name
	testName, err := getTestName()
	assert.NoError(t, err)                  // Ensure no errors occurred
	assert.Contains(t, testName, "shared_") // Validate that the name includes the expected pattern
}
