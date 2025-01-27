package component_helper

import (
	"testing"

	log "github.com/charmbracelet/log"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/component-helper/config"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) DeployDependencies(t *testing.T, config *c.Config) {
	const phaseName = "deploy dependencies"
	if config.SkipDeployDependencies {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	if len(s.Dependencies) == 0 {
		log.Info("no dependencies to deploy")
		s.logPhaseStatus(phaseName, "completed")
		return
	}

	for _, dependency := range s.Dependencies {
		log.Info("deploying dependency", "component", dependency.ComponentName, "stack", dependency.StackName)
		atmosOptions := getAtmosOptions(t, config, dependency.ComponentName, dependency.StackName, dependency.AdditionalVars)
		_, err := atmos.ApplyE(t, atmosOptions)
		if err != nil {
			s.logPhaseStatus(phaseName, "failed")
			require.NoError(t, err)
		}
	}
	s.logPhaseStatus(phaseName, "completed")
}
