package atmos

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func TestOutputNoError(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	atmosOptions := &Options{
		AtmosBasePath: testFolder,
		Component:     "terraform-no-error",
		Stack:         testStack,
		NoColor:       true,
	}
	options := WithDefaultRetryableErrors(t, atmosOptions)

	Apply(t, options)

	output := Output(t, atmosOptions, "test")
	require.Equal(t, "Hello, World", output)

	outputList := OutputList(t, atmosOptions, "test_list")
	require.Equal(t, []string{"Hello", "World"}, outputList)
}
