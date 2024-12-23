package aws_component_helper

import (
	"fmt"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	tt "github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const atmosExamplePath = "test/fixtures/aws-component-helper"

func mockAtmos() {
	atmosApply = func(t tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "apply", "-input=false", "-auto-approve")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}

	atmosDestroy = func(t tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "destroy", "-input=false", "-auto-approve")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}

	atmosPlanExitCodeE = func(t tt.TestingT, options *atmos.Options) (int, error) {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "plan", "-input=false", "-detailed-exitcode")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return 0, nil
	}

	atmosVendorPull = func(t tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "vendor", "pull")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}
}

func TestComponentTestSuitesMinimum(t *testing.T) {
	componentTestSuites := &XTestSuites{
		FixturesPath: "testdata/fixtures",
		TempDir:      "testdata/tmp",
		AwsAccountId: "123456789012",
		AwsRegion:    "us-west-2",
		suites: map[string]*XTestSuite{
			"default": {
				tests: map[string]*ComponentTest{
					"test": {
						Subject: &AtmosComponent{
							Component: "vpc",
							Stack:     "default-test",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, componentTestSuites.FixturesPath, "testdata/fixtures")
	assert.Equal(t, componentTestSuites.TempDir, "testdata/tmp")
	assert.Equal(t, componentTestSuites.AwsAccountId, "123456789012")
	assert.Equal(t, componentTestSuites.AwsRegion, "us-west-2")
	// assert.Equal(t, componentTestSuites.AtmosOptions, "us-west-2")
	assert.Equal(t, componentTestSuites.suites["default"].tests["test"].Subject.Component, "vpc")
	assert.Equal(t, componentTestSuites.suites["default"].tests["test"].Subject.Stack, "default-test")
}

func TestComponentTestSuitesCreate(t *testing.T) {
	getAwsAaccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	testFolder, err := files.CopyTerraformFolderToTemp("../../../", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	componentTestSuites := NewTestSuites(t, testFolder, "us-west-2", atmosExamplePath)

	assert.Equal(t, componentTestSuites.SourceDir, testFolder)

	// There is one testSuite
	assert.NotEmpty(t, componentTestSuites.suites)
	assert.NotNil(t, componentTestSuites.suites["default"])

	// testSuite does not have setup steps
	assert.Empty(t, componentTestSuites.suites["default"].setup)

	// testSuite has one test
	assert.NotEmpty(t, componentTestSuites.suites["default"].tests)
	assert.NotNil(t, componentTestSuites.suites["default"].tests["two-private-subnets"])

	// test has no setup step
	assert.Empty(t, componentTestSuites.suites["default"].tests["two-private-subnets"].setup)

	// test has subject
	assert.NotNil(t, componentTestSuites.suites["default"].tests["two-private-subnets"].Subject)
	assert.Equal(t, componentTestSuites.suites["default"].tests["two-private-subnets"].Subject.Component, "vpc")
	assert.Equal(t, componentTestSuites.suites["default"].tests["two-private-subnets"].Subject.Stack, "default-test")

	// test has one assert
	assert.NotEmpty(t, componentTestSuites.suites["default"].tests["two-private-subnets"].assert)
	assert.NotNil(t, componentTestSuites.suites["default"].tests["two-private-subnets"].assert[0])
	assert.Equal(t, componentTestSuites.suites["default"].tests["two-private-subnets"].assert[0].Component, "assert")
	assert.Equal(t, componentTestSuites.suites["default"].tests["two-private-subnets"].assert[0].Stack, "default-test")
}

func TestComponentTestSuitesRun(t *testing.T) {
	getAwsAaccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	mockAtmos()

	testFolder, err := files.CopyFolderToTemp("../../../", t.Name(), func(path string) bool { return true })
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	componentTestSuites := NewTestSuites(t, testFolder, "us-west-2", atmosExamplePath)

	componentTestSuites.Run(t, &atmos.Options{})
}
