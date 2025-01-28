package component_helper

import (
	"os"
	"testing"

	log "github.com/charmbracelet/log"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/component-helper/config"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) DestroyConfigFile(t *testing.T, config *c.Config) {
	const phaseName = "teardown/destroy config file"

	if anyPhasesSkipped(config) {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	err := os.Remove(config.ConfigFilePath)
	if err != nil {
		s.logPhaseStatus(phaseName, "failed")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")
}

func (s *TestSuite) DestroyDependencies(t *testing.T, config *c.Config) {
	const phaseName = "destroy dependencies"

	if config.SkipDestroyDependencies {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	if len(s.Dependencies) == 0 {
		log.Info("no dependencies to destroy")
		s.logPhaseStatus(phaseName, "completed")
		return
	}

	for i := len(s.Dependencies) - 1; i >= 0; i-- {
		dependency := s.Dependencies[i]
		log.Info("destroying dependency", "component", dependency.ComponentName, "stack", dependency.StackName)
		atmosOptions := getAtmosOptions(t, config, dependency.ComponentName, dependency.StackName, dependency.AdditionalVars)
		_, err := atmos.DestroyE(t, atmosOptions)
		if err != nil {
			s.logPhaseStatus(phaseName, "failed")
			require.NoError(t, err)
		}
	}

	s.logPhaseStatus(phaseName, "completed")
}

func (s *TestSuite) DestroyTempDir(t *testing.T, config *c.Config) {
	const phaseName = "teardown/destroy temp dir"

	if anyPhasesSkipped(config) {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	log.WithPrefix(t.Name()).Info("removing terraform state directory", "path", config.StateDir)
	err := os.RemoveAll(config.StateDir)
	require.NoError(t, err)

	log.WithPrefix(t.Name()).Info("removing temp directory", "path", config.TempDir)
	err = os.RemoveAll(config.TempDir)
	require.NoError(t, err)

	s.logPhaseStatus(phaseName, "completed")
}
