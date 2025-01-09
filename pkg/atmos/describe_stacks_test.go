package atmos

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestDescribeStacksNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		NoColor:       true,
	})

	out := DescribeStacks(t, options)

	workspace := out.Stacks["test-test-test"].Components.Terraform["terraform-no-error"].Workspace
	require.Equal(t, workspace, "test-test-test")
}
