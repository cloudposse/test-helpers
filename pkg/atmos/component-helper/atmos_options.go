package component_helper

import (
	"strings"
	"testing"

	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/component-helper/config"
	"github.com/stretchr/testify/require"
)

func getAtmosOptions(t *testing.T, config *c.Config, componentName string, stackName string, vars *map[string]interface{}) *atmos.Options {
	mergedVars := map[string]interface{}{
		"attributes": []string{config.RandomIdentifier},
	}

	// If we are not skipping the nuking of the test account, add the default tags
	// if !suite.SkipNukeTestAccount {
	// 	nukeVars := map[string]interface{}{
	// 		"default_tags": map[string]string{
	// 			"CreatedByAtmosTestHelpers": config.RandomIdentifier,
	// 		},
	// 	}

	// 	err := mergo.Merge(&mergedVars, nukeVars)
	// 	require.NoError(t, err)
	// }

	// Merge in any additional vars passed in
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
	}
	return atmosOptions
}
