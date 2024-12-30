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

func TestAcceptance(t *testing.T) {
	getAwsAaccountIdCallback = func() (string, error) {
		return "123456789012", nil
	}

	testFolder, err := files.CopyFolderToTemp("../../../", t.Name(), func(path string) bool { return true })
	require.NoError(t, err)
	defer os.RemoveAll(testFolder)

	fmt.Printf("running in %s\n", testFolder)

	fixture := NewFixture(t, testFolder, "us-west-2", "test/fixtures/aws-component-helper")

	fixture.SetUp(&atmos.Options{})
	defer fixture.TearDown()

	fixture.Suite("default", func(t *testing.T, suite *Suite) {
		suite.Test(t, "basic", func(t *testing.T, atm *Atmos) {
			inputs := map[string]interface{}{
				"cnt": 2,
			}
			defer atm.GetAndDestroy("terraform-basic-configuration", "default-test", inputs)
			atm.GetAndDeploy("terraform-basic-configuration", "default-test", inputs)

		})

		suite.Test(t, "no-error", func(t *testing.T, atm *Atmos) {
			defer atm.GetAndDestroy("terraform-no-error", "default-test", nil)
			component := atm.GetAndDeploy("terraform-no-error", "default-test", nil)

			assert.Equal(t, atm.Output(component, "test"), "Hello, World")
			assert.Equal(t, atm.OutputList(component, "test_list"), []string{"a", "b", "c"})
		})
	})
}
