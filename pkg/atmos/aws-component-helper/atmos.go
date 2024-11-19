package aws_component_helper

import (
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/stretchr/testify/require"
)

func GetAtmosOptions(t *testing.T, suite *TestSuite, componentName string, stackName string, vars map[string]interface{}) *atmos.Options {
	mergedVars := map[string]interface{}{
		"attributes": []string{suite.RandomSeed},
		"default_tags": map[string]string{
			"CreatedByTerratestRun": suite.RandomSeed,
		},
		"region": suite.AwsRegion,
		"environment": suite.RandomSeed,
	}

	err := mergo.Merge(&mergedVars, vars)
	require.NoError(t, err)

	atmosOptions := atmos.WithDefaultRetryableErrors(t, &atmos.Options{
		AtmosBasePath: suite.TempDir,
		Component:     componentName,
		Stack:         stackName,
		NoColor:       true,
		BackendConfig: map[string]interface{}{
			"workspace_key_prefix": strings.Join([]string{suite.RandomSeed, stackName}, "-"),
		},
		Vars: mergedVars,
	})
	return atmosOptions
}

func deployDependencies(t *testing.T, suite *TestSuite) error {
	for _, dependency := range suite.Dependencies {
		_, _, err := deployComponent(t, suite, dependency.ComponentName, dependency.StackName, map[string]interface{}{})
		if err != nil {
			return err
		}
	}

	return nil
}

func destroyDependencies(t *testing.T, suite *TestSuite) error {
	// iterate over dependencies in reverse order and destroy them
	for i := len(suite.Dependencies) - 1; i >= 0; i-- {
		_, _, err := destroyComponent(t, suite, suite.Dependencies[i].ComponentName, suite.Dependencies[i].StackName, map[string]interface{}{})
		if err != nil {
			return err
		}
	}
	return nil
}

func deployComponent(t *testing.T, suite *TestSuite, componentName string, stackName string, vars map[string]interface{}) (*atmos.Options, string, error) {
	options := GetAtmosOptions(t, suite, componentName, stackName, vars)
	out, err := atmos.ApplyE(t, options)

	return options, out, err
}

func destroyComponent(t *testing.T, suite *TestSuite, componentName string, stackName string, vars map[string]interface{}) (*atmos.Options, string, error) {
	options := GetAtmosOptions(t, suite, componentName, stackName, vars)
	out, err := atmos.DestroyE(t, options)

	return options, out, err
}

func vendorDependencies(t *testing.T, suite *TestSuite) error {
	options := GetAtmosOptions(t, suite, "", "", map[string]interface{}{})
	_, err := atmos.VendorPullE(t, options)

	return err
}
