package component_helper

import (
	"dario.cat/mergo"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) VerifyEnabledFlag(componentName, stackName string, additionalVars *map[string]interface{}) {
	if s.Config.SkipEnabledFlagTest {
		s.logPhaseStatus("verifiy enabled flag", "skipped")
		s.T().Skip()
	}

	mergedVars := map[string]interface{}{
		"enabled":      false,
	}

	// Merge in any additional vars passed in
	if additionalVars != nil {
		err := mergo.Merge(&mergedVars, additionalVars)
		require.NoError(s.T(), err)
	}

	s.DriftTest(componentName, stackName, &mergedVars)
}
