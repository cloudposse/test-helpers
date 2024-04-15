package atmos

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestApplyNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-no-error",
		Stack:         testStack,
		NoColor:       true,
	})

	out := Apply(t, options)

	require.Contains(t, out, "Hello, World")
}

func TestApplyNoColor(t *testing.T) {
	t.Skip("atmos doesn't support running with --no-color yet")
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)

	require.NoError(t, err)
	fmt.Printf("running in %s\n", testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-no-error",
		Stack:         testStack,
		NoColor:       true,
	})

	out := Apply(t, options)

	// Check that NoColor correctly doesn't output the colour escape codes which look like [0m,[1m or [32m
	require.NotRegexp(t, `\[\d*m`, out, "Output should not contain color escape codes")
}

func TestApplyWithErrorNoRetry(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-with-error",
		Stack:         testStack,
	})

	out, err := ApplyE(t, options)

	require.Error(t, err)
	require.Contains(t, out, "This is the first run, exiting with an error")
}

func TestApplyWithErrorWithRetry(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-with-error",
		Stack:         testStack,
		MaxRetries:    1,
		RetryableAtmosErrors: map[string]string{
			"but this error was expected and warrants a retry": "Intentional failure in test fixture",
		},
	})

	out, err := ApplyE(t, options)

	require.NotNil(t, err)
	require.Contains(t, out, "This is the first run, exiting with an error")
}

func TestIdempotentNoChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-no-error",
		Stack:         testStack,
		NoColor:       true,
	})

	_, err = ApplyAndIdempotentE(t, options)
	require.Equal(t, nil, err)
}

func TestIdempotentWithChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-not-idempotent",
		Stack:         testStack,
		NoColor:       true,
	})

	out, err := ApplyAndIdempotentE(t, options)

	require.NotEmpty(t, out)
	require.Error(t, err)
	require.EqualError(t, err, "terraform configuration not idempotent")
}

func TestParallelism(t *testing.T) {
	// This test depends on precise timing of the concurrent parallel calls in terraform, so we need to run this test
	// serially by itself so that other concurrent test runs won't influence the timing.

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-parallelism",
		Stack:         testStack,
		NoColor:       true,
	})

	// Run the first time with parallelism set to 5 and it should take about 5 seconds (plus or minus 10 seconds to
	// account for other CPU hogging stuff)
	options.Parallelism = 5
	start := time.Now()
	Apply(t, options)
	end := time.Now()
	require.WithinDuration(t, end, start, 15*time.Second)

	// Run the second time with parallelism set to 1 and it should take at least 25 seconds
	options.Parallelism = 1
	start = time.Now()
	Apply(t, options)
	end = time.Now()
	duration := end.Sub(start)
	require.GreaterOrEqual(t, int64(duration.Seconds()), int64(25))
}

func TestApplyWithPlanFile(t *testing.T) {
	t.Skip(("atmos doesn't support running with apply options with plan file yet"))
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	planFilePath := filepath.Join(testFolder, "plan.out")

	options := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-basic-configuration",
		Stack:         testStack,
		Vars: map[string]interface{}{
			"cnt": 1,
		},
		NoColor:      true,
		PlanFilePath: planFilePath,
	}
	_, err = PlanE(t, options)
	require.NoError(t, err)
	require.FileExists(t, planFilePath, "Plan file was not saved to expected location:", planFilePath)

	applyOptions := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-basic-configuration",
		Stack:         testStack,
		PlanFilePath:  planFilePath,
	}

	out, err := ApplyE(t, applyOptions)
	require.NoError(t, err)
	require.Contains(t, out, "1 added, 0 changed, 0 destroyed.")
	require.NotRegexp(t, `\[\d*m`, out, "Output should not contain color escape codes")
}
