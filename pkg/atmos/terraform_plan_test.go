package atmos

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)



func TestPlanWithExitCodeWithNoChanges(t *testing.T) {
	t.Parallel()
	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)

	options := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-no-error",
		Stack:         testStack,
	}

	Apply(t, options)
	exitCode := PlanExitCode(t, options)
	require.Equal(t, DefaultSuccessExitCode, exitCode)
}

func TestPlanWithExitCodeWithChanges(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())

	require.NoError(t, err)

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

	options := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-with-plan-error",
		Stack:         testStack,
	}

	exitCode, getExitCodeErr := PlanExitCodeE(t, options)
	require.NoError(t, getExitCodeErr)
	require.Equal(t, exitCode, 1)
}
