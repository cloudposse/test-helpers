package aws_component_helper

import (
	"github.com/cloudposse/test-helpers/pkg/atmos"
)

var (
	atmosApply         = atmos.Apply
	atmosDestroy       = atmos.Destroy
	atmosPlanExitCodeE = atmos.PlanExitCodeE
	atmosVendorPull    = atmos.VendorPull
)

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
