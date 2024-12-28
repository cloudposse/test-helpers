package aws_component_helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

var (
	atmosApply         = atmos.Apply
	atmosDestroy       = atmos.Destroy
	atmosPlanExitCodeE = atmos.PlanExitCodeE
	atmosVendorPull    = atmos.VendorPull
)

type Atmos struct {
	t       *testing.T
	options *atmos.Options
}

func NewAtmos(t *testing.T, options *atmos.Options) *Atmos {
	return &Atmos{
		t:       t,
		options: options,
	}
}

func (ts *Atmos) GetAndDeploy(t *testing.T, componentName string, stackName string) *AtmosComponent {
	component := NewAtmosComponent(componentName, stackName)
	ts.Deploy(t, component)
	return component
}

func (ts *Atmos) Deploy(t *testing.T, component *AtmosComponent) {
	options := ts.getAtmosOptions()
	options.Component = component.ComponentName
	options.Stack = component.StackName
	defer os.RemoveAll(options.AtmosBasePath)
	copyDirectoryRecursively(ts.options.AtmosBasePath, options.AtmosBasePath)
	atmosApply(t, options)
}

func (ts *Atmos) Destroy(t *testing.T, component *AtmosComponent) {
	options := ts.getAtmosOptions()
	options.Component = component.ComponentName
	options.Stack = component.StackName
	defer os.RemoveAll(options.AtmosBasePath)
	copyDirectoryRecursively(ts.options.AtmosBasePath, options.AtmosBasePath)
	atmosDestroy(t, options)
}

func (ts *Atmos) getAtmosOptions() *atmos.Options {
	result, err := ts.options.Clone()
	require.NoError(ts.t, err)

	randID := random.UniqueId()
	randomId := strings.ToLower(randID)

	basePath := filepath.Dir(filepath.Clean(ts.options.AtmosBasePath))
	dirName := filepath.Base(ts.options.AtmosBasePath)
	tmpDir := filepath.Join(basePath, fmt.Sprintf(".%s-%s", dirName, randomId))

	result.AtmosBasePath = tmpDir
	resultEnvVars := result.EnvVars
	envvars := map[string]string{
		"ATMOS_BASE_PATH":       result.AtmosBasePath,
		"ATMOS_CLI_CONFIG_PATH": result.AtmosBasePath,
	}

	err = mergo.Merge(&envvars, resultEnvVars)
	require.NoError(t, err)

	result.EnvVars = envvars

	return result
}

//func GetAtmosOptions(t *testing.T, suite *TestSuite, componentName string, stackName string, vars map[string]interface{}) *atmos.Options {
//	mergedVars := map[string]interface{}{
//		"attributes": []string{suite.RandomIdentifier},
//		"region":     suite.AwsRegion,
//	}
//
//	// If we are not skipping the nuking of the test account, add the default tags
//	if !suite.SkipNukeTestAccount {
//		nukeVars := map[string]interface{}{
//			"default_tags": map[string]string{
//				"CreatedByTerratestRun": suite.RandomIdentifier,
//			},
//		}
//
//		err := mergo.Merge(&mergedVars, nukeVars)
//		require.NoError(t, err)
//	}
//
//	// Merge in any additional vars passed in
//	err := mergo.Merge(&mergedVars, vars)
//	require.NoError(t, err)
//
//	atmosOptions := atmos.WithDefaultRetryableErrors(t, &atmos.Options{
//		AtmosBasePath: suite.TempDir,
//		ComponentName:     componentName,
//		StackName:         stackName,
//		NoColor:       true,
//		BackendConfig: map[string]interface{}{
//			"workspace_key_prefix": strings.Join([]string{suite.RandomIdentifier, stackName}, "-"),
//		},
//		Vars: mergedVars,
//	})
//	return atmosOptions
//}
//
//func deployDependencies(t *testing.T, suite *TestSuite) error {
//	for _, dependency := range suite.Dependencies {
//		_, _, err := DeployComponent(t, suite, dependency.ComponentName, dependency.StackName, map[string]interface{}{})
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func destroyDependencies(t *testing.T, suite *TestSuite) error {
//	// iterate over dependencies in reverse order and destroy them
//	for i := len(suite.Dependencies) - 1; i >= 0; i-- {
//		_, _, err := DestroyComponent(t, suite, suite.Dependencies[i].ComponentName, suite.Dependencies[i].StackName, map[string]interface{}{})
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func DeployComponent(t *testing.T, suite *TestSuite, componentName string, stackName string, vars map[string]interface{}) (*atmos.Options, string, error) {
//	options := GetAtmosOptions(t, suite, componentName, stackName, vars)
//	out, err := atmos.ApplyE(t, options)
//
//	return options, out, err
//}

//func verifyEnabledFlag(t *testing.T, suite *TestSuite, componentName string, stackName string) (*atmos.Options, error) {
//	vars := map[string]interface{}{
//		"enabled": false,
//	}
//	options := GetAtmosOptions(t, suite, componentName, stackName, vars)
//
//	exitCode, err := atmos.PlanExitCodeE(t, options)
//
//	if err != nil {
//		return options, err
//	}
//
//	if exitCode != 0 {
//		return options, fmt.Errorf("running atmos terraform plan with enabled flag set to false resulted in resource changes")
//	}
//
//	return options, nil
//}

//func DestroyComponent(t *testing.T, suite *TestSuite, componentName string, stackName string, vars map[string]interface{}) (*atmos.Options, string, error) {
//	options := GetAtmosOptions(t, suite, componentName, stackName, vars)
//	out, err := atmos.DestroyE(t, options)
//
//	return options, out, err
//}
//
//func vendorDependencies(t *testing.T, suite *TestSuite) error {
//	options := GetAtmosOptions(t, suite, "", "", map[string]interface{}{})
//	_, err := atmos.VendorPullE(t, options)
//
//	return err
//}
