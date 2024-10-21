package component

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

func scaffoldTempFolder(t *testing.T, testFolder string, componentPath []string) string {
	subFolderPath := strings.Join(append([]string{testFolder, "components", "terraform"}, componentPath...), "/")
	t.Logf("Subfolder path: %s", subFolderPath)
	err := os.MkdirAll(subFolderPath, 0777)
	require.NoError(t, err)

	return subFolderPath
}

func SetupTempDir(t *testing.T, testFolder string, componentPath []string) string,string {
	t.Helper()

	test_structure.RunTestStage(t, "copy_fixtures_to_temp_folder", func() {
		t.Log("Copying fixtures and component (src/) to temp folder...")
		testFolder, err := files.CopyTerraformFolderToTemp(opts.FixturesPath, t.Name())
		require.NoError(t, err)
		fmt.Printf("running in %s\n", testFolder)

		// Copy the component to the test folder
		componentFolderPath := scaffoldTempFolder(t, testFolder, []string{opts.ComponentName})
		err = files.CopyFolderContents("../src", commponentFolderPath)
		require.NoError(t, err)
	})

	defer test_structure.RunTestStage(t, "cleanup_temp_folder", func() {
		t.Log("Cleaning up temp folder...")
		os.RemoveAll(testFolder)
	})

	return testFolder,componentFolderPath
}
