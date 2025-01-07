package aws_component_helper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestState contains tests for the State struct and its methods
func TestState(t *testing.T) {
	// Test the basic functionality of the State struct
	t.Run("basic", func(t *testing.T) {
		*preserveStates = false // Disable state preservation for this test

		// Copy a folder to a temporary test folder
		testFolder, err := files.CopyFolderToTemp("../../../", strings.Replace(t.Name(), "/", "-", -1), func(path string) bool { return true })
		require.NoError(t, err)
		defer func() { // Ensure test folder is removed after the test
			if err := os.RemoveAll(testFolder); err != nil {
				t.Errorf("failed to cleanup test folder: %v", err)
			}
		}()

		// Create a new State instance
		state := NewState("default", testFolder)

		// Validate initial state attributes
		assert.Equal(t, state.basepath, testFolder)
		assert.Equal(t, state.namespace, "default")

		// Set up the state (create necessary directories)
		err = state.SetUp()
		assert.NoError(t, err)

		// Ensure directories exist
		assert.DirExists(t, state.basepath)
		assert.DirExists(t, state.CurrentDir())

		// Teardown the state (remove namespace-specific directory)
		state.Teardown()
		assert.NoDirExists(t, state.CurrentDir()) // Namespace directory should no longer exist
		assert.DirExists(t, state.basepath)       // Base directory should still exist
	})

	// Test the Fork method of the State struct
	t.Run("fork", func(t *testing.T) {
		*preserveStates = false // Disable state preservation for this test

		// Copy a folder to a temporary test folder
		testFolder, err := files.CopyFolderToTemp("../../../", strings.Replace(t.Name(), "/", "-", -1), func(path string) bool { return true })
		require.NoError(t, err)
		defer func() { // Ensure test folder is removed after the test
			if err := os.RemoveAll(testFolder); err != nil {
				t.Errorf("failed to cleanup test folder: %v", err)
			}
		}()

		// Create a new State instance
		state := NewState("default", testFolder)

		// Validate initial state attributes
		assert.Equal(t, state.basepath, testFolder)
		assert.Equal(t, state.namespace, "default")

		// Set up the state
		err = state.SetUp()
		assert.NoError(t, err)

		// Ensure directories exist
		assert.DirExists(t, state.basepath)
		assert.DirExists(t, state.CurrentDir())

		// Create a temporary file in the current state directory
		file, err := os.CreateTemp(state.CurrentDir(), "tmpfile")
		assert.FileExists(t, file.Name())
		fileName := filepath.Base(file.Name()) // Store the file name for verification later

		// Fork the state to create a child state
		forkState, err := state.Fork("fork")
		require.NoError(t, err)

		// Validate forked state attributes
		assert.NotEqual(t, forkState.CurrentDir(), state.CurrentDir())        // Forked state should have a different directory
		assert.DirExists(t, forkState.CurrentDir())                           // Forked directory should exist
		assert.FileExists(t, filepath.Join(forkState.CurrentDir(), fileName)) // Forked directory should contain the copied file

		// Validate parent state directories
		assert.DirExists(t, state.basepath)
		assert.DirExists(t, state.CurrentDir())

		// Teardown the parent state
		state.Teardown()
		assert.NoDirExists(t, state.CurrentDir()) // Parent namespace directory should no longer exist
		assert.DirExists(t, state.basepath)       // Base directory should still exist

		// Validate forked state directories
		assert.DirExists(t, forkState.CurrentDir())                           // Forked directory should still exist
		assert.FileExists(t, filepath.Join(forkState.CurrentDir(), fileName)) // File should still exist in the forked directory

		// Teardown the forked state
		forkState.Teardown()
		assert.NoDirExists(t, forkState.CurrentDir()) // Forked namespace directory should no longer exist
		assert.DirExists(t, forkState.basepath)       // Base directory should still exist
	})
}
