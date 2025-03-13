package examples_helper

import (
	"testing"

	log "github.com/charmbracelet/log"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
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
		if dependency.Function != nil {
			s.logPhaseStatus("deploy dependencies/function", "started")
			err := dependency.Function()
			if err != nil {
				log.WithPrefix(t.Name()+"deploy function dependency").Error("failed to run function", "error", err)
			}

			s.logPhaseStatus("deploy dependencies/function", "completed")
			continue
		} else {

			log.WithPrefix(t.Name()).Info("deploying dependency", "component", dependency.ComponentName, "stack", dependency.StackName)
			atmosOptions := getAtmosOptionsFromSetupConfiguration(t, config, s.SetupConfiguration, dependency.ComponentName, dependency.StackName, dependency.AdditionalVars)

			out, err := atmos.ApplyE(t, atmosOptions)
			log.WithPrefix(t.Name()).Info("deploying dependency", "component", dependency.ComponentName, "stack", dependency.StackName, "output", out)
			if err != nil {
				s.logPhaseStatus(phaseName, "failed")
				require.NoError(t, err)
			}
		}

	}
	s.logPhaseStatus(phaseName, "completed")
}
