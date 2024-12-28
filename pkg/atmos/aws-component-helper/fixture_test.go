package aws_component_helper

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	tt "github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

const atmosExamplePath = "test/fixtures/aws-component-helper"

func mockAtmos() {
	atmosApply = func(_ tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "apply", "-input=false", "-auto-approve")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}

	atmosDestroy = func(_ tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "destroy", "-input=false", "-auto-approve")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}

	atmosPlanExitCodeE = func(_ tt.TestingT, options *atmos.Options) (int, error) {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "plan", "-input=false", "-detailed-exitcode")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return 0, nil
	}

	atmosVendorPull = func(_ tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "vendor", "pull")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}
}

func TestFixtureMinimum(t *testing.T) {
	componentTestSuites := &Fixture{
		FixturesPath: "testdata/fixtures",
		TempDir:      "testdata/tmp",
		AwsAccountId: "123456789012",
		AwsRegion:    "us-west-2",
	}
	assert.Equal(t, componentTestSuites.FixturesPath, "testdata/fixtures")
	assert.Equal(t, componentTestSuites.TempDir, "testdata/tmp")
	assert.Equal(t, componentTestSuites.AwsAccountId, "123456789012")
	assert.Equal(t, componentTestSuites.AwsRegion, "us-west-2")
}

func TestFixtureCreate(t *testing.T) {
	getAwsAaccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	testFolder, err := files.CopyTerraformFolderToTemp("../../../", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	componentTestSuites := NewFixture(t, testFolder, "us-west-2", atmosExamplePath)

	assert.Equal(t, componentTestSuites.SourceDir, testFolder)
}

func TestFixtureSuitesRun(t *testing.T) {
	getAwsAaccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	mockAtmos()

	testFolder, err := files.CopyFolderToTemp("../../../", t.Name(), func(path string) bool { return true })
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	fixture := NewFixture(t, testFolder, "us-west-2", atmosExamplePath)

	fixture.SetUp(&atmos.Options{})
	defer fixture.TearDown()

	fixture.Suite("default", func(t *testing.T, suite *Suite) {
		// suite.AddDependency(t, "vpc", "default-test")

		// var deps *AtmosComponent

		// suite.Setup(t, func(t *testing.T, atm *Atmos) {
		// 	deps = atm.GetAndDeploy(t, "vpc/deps", "default-test")
		// })

		// suite.Test(t, "two-private-subnets", func(t *testing.T, atm *Atmos) {
		// 	component := atm.GetAndDeploy(t, "vpc/private-only", "default-test")
		// 	defer atm.Destroy(t, component)
		// })

		// suite.Test(t, "public-subnets", func(t *testing.T, atm *Atmos) {
		// 	component := atm.GetAndDeploy(t, "vpc/full", "default-test")
		// 	defer atm.Destroy(t, component)
		// })

		// suite.TearDown(t, func(t *testing.T, atm *Atmos) {
		// 	atm.Destroy(t, deps)
		// })
	})

}
