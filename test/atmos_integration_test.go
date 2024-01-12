package atmos

import (
	"testing"

	"github.com/cloudposse/terratest-helpers/pkg/atmos"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestAtmos(t *testing.T) {
	t.Parallel()

	testFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/atmos-example")

	atmosOptions := &atmos.Options{
		AtmosBasePath: testFolder,
		Component:     "test",
		Stack:         "test-test-test",
		Vars:          map[string]interface{}{},
		NoColor:       true,
		MaxRetries:    0,
	}

	atmos.Plan(t, atmosOptions)
}
