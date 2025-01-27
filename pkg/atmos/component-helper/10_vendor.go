package component_helper

import (
	"os"
	"path/filepath"
	"testing"

	log "github.com/charmbracelet/log"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	c "github.com/cloudposse/test-helpers/pkg/atmos/component-helper/config"
	"github.com/stretchr/testify/require"
)

func vendorFileExists(t *testing.T, config *c.Config) bool {
	vendorFile := filepath.Join(config.TempDir, "vendor.yaml")
	_, err := os.Stat(vendorFile)
	if os.IsNotExist(err) {
		log.WithPrefix(t.Name()).Warn("skipping test phase because vendor file does not exist", "phase", "setup test suite/vendor dependencies", "vendorFile", vendorFile)
		return false
	}

	require.NoError(t, err)

	return true
}

func (s *TestSuite) VendorDependencies(t *testing.T, config *c.Config) {
	const phaseName = "vendor dependencies"
	if !vendorFileExists(t, config) {
		return
	}

	if !vendorFileExists(t, config) || config.SkipVendorDependencies {
		s.logPhaseStatus(phaseName, "skipped")
		return
	}

	s.logPhaseStatus(phaseName, "started")

	log.Debug("running atmos vendor pull in tempdir")
	atmosOptions := getAtmosOptions(t, config, "", "", nil)
	_, err := atmos.VendorPullE(t, atmosOptions)
	if err != nil {
		s.logPhaseStatus(phaseName, "error")
		require.NoError(t, err)
	}

	s.logPhaseStatus(phaseName, "completed")
}
