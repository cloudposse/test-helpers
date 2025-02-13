package component_helper

import (
	"strings"
	"dario.cat/mergo"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) DriftTest(componentName, stackName string, additionalVars *map[string]interface{}) {

	mergedVars := map[string]interface{}{
		"attributes": []string{s.Config.RandomIdentifier},
	}

	// Merge in any additional vars passed in
	if additionalVars != nil {
		err := mergo.Merge(&mergedVars, additionalVars)
		require.NoError(s.T(), err)
	}

	atmosOptions := getAtmosOptions(s.T(), s.Config, componentName, stackName, &mergedVars)

	outputs, err := atmos.PlanE(s.T(), atmosOptions)
	require.NoError(s.T(), err)
	noChanges := strings.Contains(outputs, "No changes. Your infrastructure matches the configuration.") || strings.Contains(outputs, "without changing any real infrastructure.")
	require.True(s.T(), noChanges)
}
