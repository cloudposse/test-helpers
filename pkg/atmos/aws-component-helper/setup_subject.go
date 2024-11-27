package aws_component_helper

import (
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
)

func awsComponentTestCleanup(t *testing.T, opts *atmos.Options, destroy bool, tmpDir string, protectDir string) {
	if destroy {
		out, err := atmos.DestroyE(t, opts)
		if err == nil {
			// if the destroy was successful, remove the temp directory
			if tmpDir == protectDir {
				t.Logf("Not removing protected directory %s", protectDir)
			} else {
				t.Log("Cleaning out temp folder...")
				err = os.RemoveAll(tmpDir)
				if err == nil {
					t.Logf("Removed temp directory %s", tmpDir)
				} else {
					assert.NoError(t, err, "Failed to remove temp directory %s", tmpDir)
				}
			}
		} else {
			// if the destroy failed, leave the temp directory in place
			assert.NoError(t, err, "Failed to destroy subject. Leaving source and state in %v\n\nDestroy output:\n%v\n\n", tmpDir, out)
		}
	}
}

type ComponentTestResults struct {
	Output string
	TestID string
}

func IsLocalAwsComponentTest(t *testing.T) bool {
	atmosPlainOpts := &atmos.Options{}
	existingTestID, err := atmos.RunAtmosCommandE(t, atmosPlainOpts, "test", "get-test-id")

	return err == nil && len(existingTestID) > 0
}

func AwsComponentTestHelper(t *testing.T, opts *atmos.Options, callback func(t *testing.T, opts *atmos.Options, results ComponentTestResults)) {
	testSrcRoot := os.Getenv("ATMOS_BASE_PATH")
	testRoot := testSrcRoot
	if testSrcRoot == "" {
		assert.FailNow(t, "ATMOS_BASE_PATH must be set, but is empty")
	}

	atmosPlainOpts := &atmos.Options{}
	doApply := false

	var testID string

	existingTestID, err := atmos.RunAtmosCommandE(t, atmosPlainOpts, "test", "get-test-id")
	if err != nil || len(existingTestID) == 0 {
		doApply = true
		// Copy test source to a temp directory and create a new test ID
		t.Log("Copying files to temp folder...")
		testRoot, err = files.CopyTerraformFolderToTemp(testSrcRoot, t.Name())
		require.NoError(t, err)
		err = copyDirectoryRecursively(filepath.Join(testSrcRoot, "state"), filepath.Join(testRoot, "state"))
		require.NoError(t, err)
		atmosPlainOpts.AtmosBasePath = testRoot
		testID, err = atmos.RunAtmosCommandE(t, atmosPlainOpts, "test", "make-test-id")
		require.NoError(t, err)
		testID = strings.TrimSpace(testID)
	} else {
		testID = strings.TrimSpace(existingTestID)
	}

	t.Logf("Running test \"%s\" with test ID \"%s\" in directory %s", t.Name(), testID, testRoot)

	options := atmos.WithDefaultRetryableErrors(t, opts)
	options.AtmosBasePath = testRoot
	// Keep the output quiet
	if !testing.Verbose() {
		options.Logger = logger.Discard
	}

	defer awsComponentTestCleanup(t, options, doApply, testRoot, testSrcRoot)

	// Apply the deployment
	out := ""
	if doApply {
		out, err = atmos.ApplyE(t, options)
		require.NoError(t, err, "Failed to deploy component, skipping other tests.")
	}
	// Call the callback function for assertions
	callback(t, options, ComponentTestResults{
		Output: out,
		TestID: testID,
	})

	if !doApply {
		t.Logf("\n\n\nTests complete in %s\n\n", testRoot)
	}
}
