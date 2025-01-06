package aws_component_helper

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAcceptance is an end-to-end acceptance test for AWS component helper functionality.
func TestAcceptance(t *testing.T) {
	// Mock AWS account ID retrieval callback
	getAwsAaccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	// Copy the test folder to a temporary location for isolation
	testFolder, err := files.CopyFolderToTemp("../../../", t.Name(), func(path string) bool { return true })
	require.NoError(t, err)
	defer func() {
		// Clean up the temporary test folder
		if err := os.RemoveAll(testFolder); err != nil {
			t.Errorf("failed to cleanup test folder: %v", err)
		}
	}()

	fmt.Printf("running in %s\n", testFolder)

	// Initialize a new test fixture for the acceptance test
	fixture := NewFixture(t, testFolder, "us-west-2", "test/fixtures/aws-component-helper")

	// Set up and tear down the test environment
	fixture.SetUp(&atmos.Options{})
	defer fixture.TearDown()

	// Run a suite of tests within the "default" group
	fixture.Suite("default", func(t *testing.T, suite *Suite) {
		// Test case: "basic"
		suite.Test(t, "basic", func(t *testing.T, atm *Atmos) {
			inputs := map[string]interface{}{
				"cnt": 2, // Input variable for the Terraform configuration
			}
			// Ensure proper cleanup by destroying resources after the test
			defer atm.GetAndDestroy("terraform-basic-configuration", "default-test", inputs)
			// Deploy the Terraform component
			atm.GetAndDeploy("terraform-basic-configuration", "default-test", inputs)
		})

		// Test case: "no-error"
		suite.Test(t, "no-error", func(t *testing.T, atm *Atmos) {
			// Ensure proper cleanup by destroying resources after the test
			defer atm.GetAndDestroy("terraform-no-error", "default-test", nil)
			// Deploy the Terraform component
			component := atm.GetAndDeploy("terraform-no-error", "default-test", nil)

			// Expected outputs from the Terraform component
			mapOfObjects := map[string]interface{}{
				"a": map[string]interface{}{"b": "c"},
				"d": map[string]interface{}{"e": "f"},
			}

			// Assert outputs against expected values
			assert.Equal(t, "Hello, World", atm.Output(component, "test"))
			assert.Equal(t, []string{"a", "b", "c"}, atm.OutputList(component, "test_list"))
			assert.Equal(t, mapOfObjects, atm.OutputMapOfObjects(component, "test_map_of_objects"))

			// Struct definitions for complex outputs
			type structValue1 struct {
				B string
			}

			type structValue2 struct {
				E string
			}

			type structValue struct {
				A structValue1
				D structValue2
			}

			// Parse and validate structured output
			structResult := structValue{}
			structExpected := structValue{
				A: structValue1{B: "c"},
				D: structValue2{E: "f"},
			}

			// Assert that the parsed struct matches the expected structure
			atm.OutputStruct(component, "test_map_of_objects", &structResult)
			assert.Equal(t, structExpected, structResult)
		})
	})
}
