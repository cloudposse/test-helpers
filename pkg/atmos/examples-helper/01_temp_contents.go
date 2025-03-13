package examples_helper

import (
	"github.com/charmbracelet/log"
	c "github.com/cloudposse/test-helpers/pkg/atmos/examples-helper/config"
	"github.com/stretchr/testify/require"

	"testing"
)

func (s *TestSuite) CreateTempContents(t *testing.T, config *c.Config) {
	const phaseName = "run prehook"
	if config.SkipTempContents {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	cmd := s.SetupConfiguration.TempContentsCmd
	if s.Config.SkipTempContents || cmd == nil {
		return
	}
	t.Setenv("TEST_TEMP_DIR", config.TempDir)
	log.WithPrefix(t.Name()).Info("Running Temp Contents Command", "command", cmd.String(), "dir", cmd.Dir)

	//var stdout, stderr bytes.Buffer
	//cmd.Stdout = &stdout
	//cmd.Stderr = &stderr

	err := cmd.Run()

	//if stdout.Len() > 0 {
	//	log.WithPrefix(t.Name()).Info("Command:", "stdout", stdout.String())
	//}
	//if stderr.Len() > 0 {
	//	log.WithPrefix(t.Name()).Warn("Command:", "stderr", stderr.String())
	//}
	require.NoError(t, err)

	s.logPhaseStatus(phaseName, "completed")
}
