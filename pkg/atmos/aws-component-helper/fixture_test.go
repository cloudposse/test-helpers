package aws_component_helper

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	tt "github.com/cloudposse/test-helpers/pkg/testing"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const atmosExamplePath = "test/fixtures/aws-component-helper"

// mockAtmos provides mock implementations for Atmos commands used in the tests
func mockAtmos() {
	// Mock the "apply" command
	atmosApply = func(_ tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "apply", "-input=false", "-auto-approve")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}

	// Mock the "destroy" command
	atmosDestroy = func(_ tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "destroy", "-input=false", "-auto-approve")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}

	// Mock the "plan" command with exit code
	atmosPlanExitCodeE = func(_ tt.TestingT, options *atmos.Options) (int, error) {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "plan", "-input=false", "-detailed-exitcode")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return 0, nil
	}

	// Mock the "vendor pull" command
	atmosVendorPull = func(_ tt.TestingT, options *atmos.Options) string {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "vendor", "pull")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return ""
	}

	// Mock the "output all" command
	atmosOutputAllE = func(_ tt.TestingT, options *atmos.Options, _ string, _ interface{}) error {
		options, args := atmos.GetCommonOptions(options, atmos.FormatArgs(options, "terraform", "output", "--skip-init", "--json")...)
		description := fmt.Sprintf("%s %v", options.AtmosBinary, args)
		fmt.Println(description)
		return nil
	}
}

// TestFixtureMinimum validates the minimal initialization of a Fixture
func TestFixtureMinimum(t *testing.T) {
	componentTestSuites := &Fixture{
		FixturesPath: "testdata/fixtures",
		TempDir:      "testdata/tmp",
		AwsAccountId: "123456789012",
		AwsRegion:    "us-west-2",
	}

	// Assert that the Fixture fields are correctly initialized
	assert.Equal(t, componentTestSuites.FixturesPath, "testdata/fixtures")
	assert.Equal(t, componentTestSuites.TempDir, "testdata/tmp")
	assert.Equal(t, componentTestSuites.AwsAccountId, "123456789012")
	assert.Equal(t, componentTestSuites.AwsRegion, "us-west-2")
}

// TestFixtureCreate validates the setup of a Fixture with a temporary folder
func TestFixtureCreate(t *testing.T) {
	// Mock AWS Account ID retrieval
	getAwsAccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	// Create a temporary folder for the test
	testFolder, err := files.CopyTerraformFolderToTemp("../../../", t.Name())
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	// Create a new Fixture
	componentTestSuites := NewFixture(t, testFolder, "us-west-2", atmosExamplePath)

	assert.Equal(t, componentTestSuites.SourceDir, testFolder)

	defer componentTestSuites.TearDown() // Ensure cleanup
	componentTestSuites.SetUp(&atmos.Options{})
}

// TestFixtureSuitesRun validates the execution of test suites within a Fixture
func TestFixtureSuitesRun(t *testing.T) {
	// Mock AWS Account ID retrieval
	getAwsAccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	// Mock Atmos commands
	mockAtmos()

	// Create a temporary folder for the test
	testFolder, err := files.CopyFolderToTemp("../../../", t.Name(), func(path string) bool { return true })
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	// Initialize the Fixture
	fixture := NewFixture(t, testFolder, "us-west-2", atmosExamplePath)

	fixture.SetUp(&atmos.Options{})
	defer fixture.TearDown() // Ensure cleanup

	// Define a test suite
	fixture.Suite("default", func(t *testing.T, suite *Suite) {
		// Add a dependency
		suite.AddDependency("vpc/deps", "default-test")

		// Setup step
		suite.Setup(t, func(t *testing.T, atm *Atmos) {
			atm.GetAndDeploy("vpc/manual-deps", "default-test", nil)
		})

		// First test case
		suite.Test(t, "two-private-subnets", func(t *testing.T, atm *Atmos) {
			inputs := map[string]interface{}{
				"name":                    "vpc-terraform",
				"availability_zones":      []string{"a", "b"},
				"public_subnets_enabled":  false,
				"nat_gateway_enabled":     false,
				"nat_instance_enabled":    false,
				"subnet_type_tag_key":     "eg.cptest.co/subnet/type",
				"max_subnet_count":        3,
				"vpc_flow_logs_enabled":   false,
				"ipv4_primary_cidr_block": "172.16.0.0/16",
			}

			component := atm.GetAndDeploy("vpc/private-only", "default-test", inputs)
			defer atm.Destroy(component)
		})

		// Second test case
		suite.Test(t, "public-subnets", func(t *testing.T, atm *Atmos) {
			component := atm.GetAndDeploy("vpc/full", "default-test", nil)
			defer atm.Destroy(component)
		})

		// Teardown step
		suite.TearDown(t, func(t *testing.T, atm *Atmos) {
			atm.GetAndDestroy("vpc/manual-deps", "default-test", nil)
		})
	})
}
