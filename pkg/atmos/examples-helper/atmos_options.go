package examples_helper

import (
	"github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/dependency"
	"path/filepath"
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)

func GetAtmosOptions(t *testing.T, config *c.Config, componentName string, stackName string, vars *map[string]interface{}) *atmos.Options {
	mergedVars := map[string]interface{}{
		"attributes": []string{config.RandomIdentifier},
	}

	if vars != nil {
		err := mergo.Merge(&mergedVars, vars)
		require.NoError(t, err)
	}

	accountID := aws.GetAccountId(t)
	require.NotEmpty(t, accountID)

	atmosOptions := &atmos.Options{
		AtmosBasePath: config.TempDir,
		Component:     componentName,
		Stack:         stackName,

		NoColor: true,
		BackendConfig: map[string]interface{}{
			"workspace_key_prefix": strings.Join([]string{config.RandomIdentifier, stackName}, "-"),
		},
		Vars: mergedVars,
		EnvVars: map[string]string{
			"ATMOS_BASE_PATH":            config.TempDir,
			"ATMOS_CLI_CONFIG_PATH":      config.TempDir,
			"COMPONENT_HELPER_STATE_DIR": config.StateDir,
			"TEST_ACCOUNT_ID":            accountID,
		},
	}
	return atmosOptions
}

func getAtmosOptionsFromSetupConfiguration(t *testing.T, config *c.Config, configuration SetupConfiguration, componentName string, stackName string, vars *map[string]interface{}, targets []string) *atmos.Options {
	mergedVars := map[string]interface{}{
		"attributes": []string{config.RandomIdentifier},
	}

	if vars != nil {
		err := mergo.Merge(&mergedVars, vars)
		require.NoError(t, err)
	}

	atmosOptions := &atmos.Options{
		AtmosBasePath: filepath.Join(config.TempDir, configuration.AtmosBaseDir),
		Component:     componentName,
		Stack:         stackName,
		NoColor:       true,
		BackendConfig: map[string]interface{}{
			"workspace_key_prefix": strings.Join([]string{config.RandomIdentifier, stackName}, "-"),
		},
		Vars: mergedVars,
		EnvVars: map[string]string{
			"ATMOS_BASE_PATH":            filepath.Join(config.TempDir, configuration.AtmosBaseDir),
			"ATMOS_CLI_CONFIG_PATH":      filepath.Join(config.TempDir, configuration.AtmosBaseDir),
			"COMPONENT_HELPER_STATE_DIR": config.StateDir,
		},
		Targets: targets,
	}
	return atmosOptions
}

func getAtmosOptions(t *testing.T, config *c.Config, s *TestSuite, d *dependency.Dependency) *atmos.Options {
	mergedVars := map[string]interface{}{}
	if d.AddRandomAttribute {
		mergedVars = map[string]interface{}{
			"attributes": []string{config.RandomIdentifier},
		}
	}

	if d.AdditionalVars != nil {
		err := mergo.Merge(&mergedVars, d.AdditionalVars)
		require.NoError(t, err)
	}

	atmosOptions := &atmos.Options{
		AtmosBasePath: filepath.Join(config.TempDir, s.SetupConfiguration.AtmosBaseDir),
		Component:     d.ComponentName,
		Stack:         d.StackName,
		NoColor:       true,
		BackendConfig: map[string]interface{}{
			"workspace_key_prefix": strings.Join([]string{config.RandomIdentifier, d.StackName}, "-"),
		},
		Vars: mergedVars,
		EnvVars: map[string]string{
			"ATMOS_BASE_PATH":            filepath.Join(config.TempDir, s.SetupConfiguration.AtmosBaseDir),
			"ATMOS_CLI_CONFIG_PATH":      filepath.Join(config.TempDir, s.SetupConfiguration.AtmosBaseDir),
			"COMPONENT_HELPER_STATE_DIR": config.StateDir,
		},
		Targets: d.Targets,
	}
	return atmosOptions
}
