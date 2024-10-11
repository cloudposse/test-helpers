package component

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cloudposse/terratest-helpers/pkg/atmos"
	"github.com/cloudposse/terratest-helpers/pkg/awsnuke"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/random"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/require"
)

func MakeComponentFolder(t *testing.T, testFolder string, componentPath []string) string {
	subFolderPath := strings.Join(append([]string{testFolder, "components", "terraform"}, componentPath...), "/")
	t.Logf("Subfolder path: %s", subFolderPath)
	err := os.MkdirAll(subFolderPath, 0777)
	require.NoError(t, err)

	return subFolderPath
}

func AwsComponentTestHelper(t *testing.T, opts AwsComponentTestOptions, callback func(t *testing.T, opts *atmos.Options, output string)) {
	t.Helper()

	var options *atmos.Options
	var out string
	var randID string
	var testFolder string

	test_structure.RunTestStage(t, "copy_fixtures_to_temp_folder", func() {
		t.Log("Copying fixtures and component (src/) to temp folder...")
		testFolder, err := files.CopyTerraformFolderToTemp(opts.FixturesPath, t.Name())
		require.NoError(t, err)
		fmt.Printf("running in %s\n", testFolder)

		// Copy the component to the test folder
		commponentFolderPath := MakeComponentFolder(t, testFolder, []string{opts.ComponentName})
		err = files.CopyFolderContents("../src", commponentFolderPath)
		require.NoError(t, err)
	})

	defer test_structure.RunTestStage(t, "cleanup_temp_folder", func() {
		t.Log("Cleaning up temp folder...")
		os.RemoveAll(testFolder)
	})

	test_structure.RunTestStage(t, "terraform_apply_dependencies", func() {
		for _, dependency := range opts.StackDependencies {
			t.Logf("Applying dependency: %s in stack %s in region %s", dependency.Component, dependency.StackName, dependency.Region)
			_, err := atmos.ApplyE(t, atmos.WithDefaultRetryableErrors(t, &atmos.Options{
				AtmosBasePath: testFolder,
				Component:     dependency.Component,
				Stack:         dependency.StackName,
				NoColor:       true,
				BackendConfig: map[string]interface{}{
					"workspace_key_prefix": strings.Join([]string{randID, opts.StackName}, "-"),
				},
				Vars: map[string]interface{}{
					"attributes": []string{randID},
					"default_tags": map[string]string{
						"CreatedByTerratestRun": randID,
					},
					"region": dependency.Region,
				},
			}))
			require.NoError(t, err)
		}
	})

	test_structure.RunTestStage(t, "terraform_setup", func() {
		// Perform setup tasks here
		t.Log("Performing test setup...")
		randID = strings.ToLower(random.UniqueId())
		atmosOptions := atmos.WithDefaultRetryableErrors(t, &atmos.Options{
			AtmosBasePath: testFolder,
			Component:     opts.ComponentName,
			Stack:         opts.StackName,
			NoColor:       true,
			BackendConfig: map[string]interface{}{
				"workspace_key_prefix": strings.Join([]string{randID, opts.StackName}, "-"),
			},
			Vars: map[string]interface{}{
				"attributes": []string{randID},
				"default_tags": map[string]string{
					"CreatedByTerratestRun": randID,
				},
				"region": opts.AwsRegion,
			},
		})
		options = atmos.WithDefaultRetryableErrors(t, atmosOptions)
	})

	defer test_structure.RunTestStage(t, "aws_nuke", func() {
		if !opts.SkipAwsNuke {
			awsnuke.NukeTestAccountByTag(t, "CreatedByTerratestRun", randID, []string{opts.AwsRegion}, false)
		}
	})

	test_structure.RunTestStage(t, "terraform_destroy_dependencies", func() {
		for _, dependency := range opts.StackDependencies {
			t.Logf("Destroying dependency: %s in stack %s in region %s", dependency.Component, dependency.StackName, dependency.Region)
			_, err := atmos.ApplyE(t, atmos.WithDefaultRetryableErrors(t, &atmos.Options{
				AtmosBasePath: testFolder,
				Component:     dependency.Component,
				Stack:         dependency.StackName,
				NoColor:       true,
				BackendConfig: map[string]interface{}{
					"workspace_key_prefix": strings.Join([]string{randID, opts.StackName}, "-"),
				},
				Vars: map[string]interface{}{
					"attributes": []string{randID},
					"default_tags": map[string]string{
						"CreatedByTerratestRun": randID,
					},
					"region": dependency.Region,
				},
			}))
			require.NoError(t, err)
		}
	})

	defer test_structure.RunTestStage(t, "terraform_destroy", func() {
		atmos.Destroy(t, options)
	})

	test_structure.RunTestStage(t, "terraform_apply", func() {
		out = atmos.Apply(t, options)
	})

	test_structure.RunTestStage(t, "callback", func() {
		callback(t, options, out)
	})
}
