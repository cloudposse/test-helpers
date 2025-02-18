package component_helper

import (
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/component-helper/config"
	"github.com/stretchr/testify/require"
	"github.com/gruntwork-io/terratest/modules/aws"
)

func getAtmosOptions(t *testing.T, config *c.Config, componentName string, stackName string, vars *map[string]interface{}) *atmos.Options {
	mergedVars := map[string]interface{}{
		"attributes": []string{config.RandomIdentifier},
	}

	if vars != nil {
		err := mergo.Merge(&mergedVars, vars)
		require.NoError(t, err)
	}

	atmosOptions := &atmos.Options{
		AtmosBasePath: config.TempDir,
		Component:     componentName,
		Stack:         stackName,
		NoColor:       true,
		BackendConfig: map[string]interface{}{
			"workspace_key_prefix": strings.Join([]string{config.RandomIdentifier, stackName}, "-"),
		},
		Vars: mergedVars,
		EnvVars: map[string]string{
			"ATMOS_BASE_PATH":            config.TempDir,
			"ATMOS_CLI_CONFIG_PATH":      config.TempDir,
			"COMPONENT_HELPER_STATE_DIR": config.StateDir,
func getAtmosOptions(t *testing.T, config *c.Config, componentName string, stackName string, vars *map[string]interface{}) *atmos.Options {
	accountID, err := aws.GetAccountId(t)
	require.NoError(t, err, "Failed to get AWS Account ID")

	mergedVars := map[string]interface{}{
		"attributes": []string{config.RandomIdentifier},
	}

	// ... rest of the function ...

	return &atmos.Options{
		ATMOS_BASE_PATH:            config.TempDir,
		ATMOS_CLI_CONFIG_PATH:      config.TempDir,
		COMPONENT_HELPER_STATE_DIR: config.StateDir,
		TEST_ACCOUNT_ID:            accountID,
	}
}
		},
	}
	return atmosOptions
}
