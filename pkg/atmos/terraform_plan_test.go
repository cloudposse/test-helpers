package atmos

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestPlanWithExitCodeWithNoChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-no-error",
		Stack:         testStack,
	}

	Apply(t, options)
	_, err = GetExitCodeForAtmosCommandE(t, options, "version")
	require.NoError(t, err)

	exitCode := PlanExitCode(t, options)
	require.Equal(t, DefaultSuccessExitCode, exitCode)
}

func TestPlanWithExitCodeWithChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-basic-configuration",
		Stack:         testStack,
		Vars: map[string]interface{}{
			"cnt": 1,
		},
	}
	exitCode := PlanExitCode(t, options)
	require.Equal(t, TerraformPlanChangesPresentExitCode, exitCode)
}

func TestPlanWithExitCodeWithFailure(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	options := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-with-plan-error",
		Stack:         testStack,
	}

	exitCode, getExitCodeErr := PlanExitCodeE(t, options)
	require.NoError(t, getExitCodeErr)
	require.Equal(t, exitCode, 1)
}
