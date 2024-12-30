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

func TestState(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		*preserveStates = false

		testFolder, err := files.CopyFolderToTemp("../../../", strings.Replace(t.Name(), "/", "-", -1), func(path string) bool { return true })
		require.NoError(t, err)
		defer os.RemoveAll(testFolder)

		state := NewState("default", testFolder)

		assert.Equal(t, state.basepath, testFolder)
		assert.Equal(t, state.namespace, "default")

		err = state.SetUp()
		assert.NoError(t, err)

		assert.DirExists(t, state.basepath)
		assert.DirExists(t, state.CurrentDir())

		state.Teardown()
		assert.NoDirExists(t, state.CurrentDir())
		assert.DirExists(t, state.basepath)
	})

	t.Run("fork", func(t *testing.T) {
		*preserveStates = false

		testFolder, err := files.CopyFolderToTemp("../../../", strings.Replace(t.Name(), "/", "-", -1), func(path string) bool { return true })
		require.NoError(t, err)
		defer os.RemoveAll(testFolder)

		state := NewState("default", testFolder)

		assert.Equal(t, state.basepath, testFolder)
		assert.Equal(t, state.namespace, "default")

		err = state.SetUp()
		assert.NoError(t, err)

		assert.DirExists(t, state.basepath)
		assert.DirExists(t, state.CurrentDir())

		file, err := os.CreateTemp(state.CurrentDir(), "tmpfile")
		assert.FileExists(t, file.Name())

		fileName := filepath.Base(file.Name())

		forkState, err := state.Fork("fork")
		require.NoError(t, err)
		assert.NotEqual(t, forkState.CurrentDir(), state.CurrentDir())
		assert.DirExists(t, forkState.CurrentDir())
		assert.FileExists(t, filepath.Join(forkState.CurrentDir(), fileName))

		assert.DirExists(t, state.basepath)
		assert.DirExists(t, state.CurrentDir())
		state.Teardown()
		assert.NoDirExists(t, state.CurrentDir())
		assert.DirExists(t, state.basepath)

		assert.DirExists(t, forkState.CurrentDir())
		assert.FileExists(t, filepath.Join(forkState.CurrentDir(), fileName))
		forkState.Teardown()
		assert.NoDirExists(t, forkState.CurrentDir())
		assert.DirExists(t, forkState.basepath)
	})

}
