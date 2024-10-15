package atmos

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/awsnuke"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func MakeComponentFolder(t *testing.T, testFolder string, componentPath []string) string {
	subFolderPath := strings.Join(append([]string{testFolder, "components", "terraform"}, componentPath...), "/")
	t.Logf("Subfolder path: %s", subFolderPath)
	err := os.MkdirAll(subFolderPath, 0777)
	require.NoError(t, err)

	return subFolderPath
}

type AwsComponentTestOptions struct {
	AwsRegion     string
	ComponentName string
	FixturesPath  string
	SkipAwsNuke   bool
	StackName     string
}

// Option type represents a configuration option
type AwsComponentTestOption func(*AwsComponentTestOptions)

// WithFixturesPath is an option for setting the FixturesPath
func WithFixturesPath(fixturesPath string) AwsComponentTestOption {
	return func(a *AwsComponentTestOptions) {
		a.FixturesPath = fixturesPath
	}
}

// WithSkipAwsNuke is an option for setting SkipAwsNuke
func WithSkipAwsNuke(skip bool) AwsComponentTestOption {
	return func(a *AwsComponentTestOptions) {
		a.SkipAwsNuke = skip
	}
}

// NewAwsComponentTestOptions creates a new AwsComponentTestOptions with required fields and optional configuration
func NewAwsComponentTestOptions(awsRegion, componentName, stackName string, opts ...AwsComponentTestOption) AwsComponentTestOptions {
	options := &AwsComponentTestOptions{
		AwsRegion:     awsRegion,
		ComponentName: componentName,
		StackName:     stackName,
		FixturesPath:  "./fixtures",
		SkipAwsNuke:   false,
	}

	// Apply optional configurations
	for _, opt := range opts {
		opt(options)
	}

	return *options
}

func AwsComponentTestHelper(t *testing.T, opts AwsComponentTestOptions, callback func(t *testing.T, opts *Options, output string)) {
	t.Helper() // Marks this function as a test helper

	t.Log("Copying fixtures and component (src/) to temp folder...")
	testFolder, err := files.CopyTerraformFolderToTemp(opts.FixturesPath, t.Name())
	require.NoError(t, err)
	fmt.Printf("running in %s\n", testFolder)
	defer os.RemoveAll(testFolder)

	// Copy the component to the test folder
	commponentFolderPath := MakeComponentFolder(t, testFolder, []string{opts.ComponentName})
	err = files.CopyFolderContents("../src", commponentFolderPath)
	require.NoError(t, err)

	// Perform setup tasks here
	t.Log("Performing test setup...")
	randID := strings.ToLower(random.UniqueId())
	atmosOptions := WithDefaultRetryableErrors(t, &Options{
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
	options := WithDefaultRetryableErrors(t, atmosOptions)

	// Clean up after the test with deferred functions
	if !opts.SkipAwsNuke {
		defer awsnuke.NukeTestAccountByTag(t, "CreatedByTerratestRun", randID, []string{opts.AwsRegion}, false)
	}
	defer Destroy(t, options)

	// Apply the deployment
	out := Apply(t, options)

	// Call the callback function for assertions
	callback(t, options, out)
}
