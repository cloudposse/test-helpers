package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/aws-component-helper"
	"github.com/stretchr/testify/require"
	//"github.com/stretchr/testify/require"
)

var suite *helper.Suite

// TestMain is the entry point for the test suite. It initializes the test
// suite and runs the tests.
func TestMain(m *testing.M) {
	var err error

	// Configure the test suite
	suite, err = helper.NewTestSuite("us-east-2", "basic", "test-use2-sandbox")
	if err != nil {
		panic(err)
	}

	// Add dependencies for the component under test in the same stack. If you
	// want to add dependencies in different stacks, use AddDependenciesWithStacks.
	//
	// Dependencies are deployed in serial in the order they are added.
	suite.AddDependencies([]string{"dep1", "dep2"})

	// Create a new testing object since TestMain doesn't have one and we need
	// one to call the Setup and Teardown functions
	t := &testing.T{}

	defer suite.TearDown(t)
	err = suite.Setup(t)
	if err != nil {
		panic(err)
	}

	if !suite.SkipTests {
		m.Run()
	}
}

func TestBasic(t *testing.T) {
	additionalVars := map[string]interface{}{
		"revision": "fromTest",
	}
	defer suite.DestroyComponentUnderTest(t, additionalVars)

	_, err := suite.DeployComponentUnderTest(t, additionalVars)
	require.NoError(t, err)

	revision := atmos.Output(t, suite.AtmosOptions, "revision")
	expected := fmt.Sprintf("%s-%s", strings.ToLower(suite.RandomIdentifier), additionalVars["revision"])
	require.Equal(t, expected, revision)
}
