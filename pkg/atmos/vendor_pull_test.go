package atmos

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
)

func setupVendorTest(t *testing.T) (string, func()) {
	testFolder, err := files.CopyTerraformFolderToTemp(atmosExamplePath, t.Name())
	require.NoError(t, err, "Failed to copy Terraform folder to temp directory")

	cleanup := func() {
		os.RemoveAll(testFolder)
	}

	return testFolder, cleanup
}

func TestVendorPullBasic(t *testing.T) {
	t.Parallel()

	testFolder, cleanup := setupVendorTest(t)
	defer cleanup()

	fmt.Printf("running in %s\n", testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
	})

	_, err := VendorPullE(t, options)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(testFolder, "components", "terraform", "vpc"))
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(testFolder, "components", "terraform", "vpc-flow-logs-bucket"))
	require.NoError(t, err)
}

func TestVendorPullSingleComponent(t *testing.T) {
	t.Parallel()

	testFolder, cleanup := setupVendorTest(t)
	defer cleanup()

	fmt.Printf("running in %s\n", testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath:   testFolder,
		VendorComponent: "vpc",
	})

	_, err := VendorPullE(t, options)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(testFolder, "components", "terraform", "vpc"))
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(testFolder, "components", "terraform", "vpc-flow-logs-bucket"))
	require.Error(t, err)
}

func TestVendorPullByTag(t *testing.T) {
	t.Parallel()

	testFolder, cleanup := setupVendorTest(t)
	defer cleanup()

	fmt.Printf("running in %s\n", testFolder)

	options := WithDefaultRetryableErrors(t, &Options{
		AtmosBasePath: testFolder,
		VendorTags:    []string{"storage"},
	})

	_, err := VendorPullE(t, options)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(testFolder, "components", "terraform", "vpc"))
	require.Error(t, err)

	_, err = os.Stat(filepath.Join(testFolder, "components", "terraform", "vpc-flow-logs-bucket"))
	require.NoError(t, err)
}
